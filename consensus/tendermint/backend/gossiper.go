package backend

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/fixsizecache"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/router"
	"github.com/autonity/autonity/consensus/tendermint/router/interfaces"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/metrics"
	"github.com/autonity/autonity/rlp"
)

type msgRouter interface {
	SetBroadcaster(broadcaster interfaces.PeerFinder)
	Forward(committee *types.Committee, m message.Msg, sender common.Address, recipients []common.Address)
}

var (
	slowGossipCounter = metrics.GetOrRegisterCounter("gossiper/slowGossip", nil) //nolint:goconst
	gossipCounter     = metrics.GetOrRegisterCounter("gossiper/Gossip", nil)     //nolint:goconst
)

type Gossiper struct {
	knownMessages *fixsizecache.Cache[common.Hash, bool] // the cache of self messages
	address       common.Address                         // address of the local peer
	broadcaster   consensus.Broadcaster
	logger        log.Logger
	stopped       chan struct{}
	router        msgRouter
}

func NewGossiper(
	knownMessages *fixsizecache.Cache[common.Hash, bool],
	address common.Address,
	logger log.Logger,
	stopped chan struct{},
	router msgRouter,
) *Gossiper {
	return &Gossiper{
		knownMessages: knownMessages,
		address:       address,
		logger:        logger,
		stopped:       stopped,
		router:        router,
	}
}

func (g *Gossiper) SetBroadcaster(broadcaster consensus.Broadcaster) {
	g.broadcaster = broadcaster
	g.router.SetBroadcaster(broadcaster)
}

func (g *Gossiper) Broadcaster() consensus.Broadcaster {
	return g.broadcaster
}

func (g *Gossiper) KnownMessages() *fixsizecache.Cache[common.Hash, bool] {
	return g.knownMessages
}

func (g *Gossiper) Address() common.Address {
	return g.address
}

func (g *Gossiper) UpdateStopChannel(stopCh chan struct{}) {
	g.stopped = stopCh
}

func (g *Gossiper) SlowGossip(committee *types.Committee, msg message.Msg, sender common.Address) {
	// only gossip to very small committee
	numTargets := len(committee.Members)
	targetIndices := rand.Perm(numTargets) // target indices to select from the full committee
	if numTargets > router.ScaleThresholdForClustering {
		numTargets = int(math.Sqrt(float64(numTargets)))
	}
	recipients := make([]common.Address, numTargets)
	for i := 0; i < numTargets; i++ {
		recipients[i] = committee.Members[targetIndices[i]].Address
	}
	// self message caching only if it is a message that has been locally created.
	// this can happen for proposals and votes of the local validator + aggregates created locally
	if sender == g.address {
		g.knownMessages.Add(msg.Hash(), true)
	}
	if metrics.Enabled {
		slowGossipCounter.Inc(1)
	}
	g.router.Forward(committee, msg, sender, recipients)
}

func (g *Gossiper) Gossip(committee *types.Committee, msg message.Msg, sender common.Address) {
	// self message caching only if it is a message that has been locally created.
	// this can happen for proposals and votes of the local validator + aggregates created locally
	if sender == g.address {
		g.knownMessages.Add(msg.Hash(), true)
	}
	if metrics.Enabled {
		gossipCounter.Inc(1) // increment gossip counter
	}
	g.router.Forward(committee, msg, sender, nil)
}

func (g *Gossiper) AskSync(committee *types.Committee, syncMsg *message.AskSyncMsg) error {
	// bail out early if we don't have a broadcaster
	if g.broadcaster == nil {
		return fmt.Errorf("broadcaster not initialized")
	}

	encoded, err := rlp.EncodeToBytes(syncMsg)
	if err != nil {
		log.Error("Error encoding sync msg", "err", err)
		panic("cannot encode sync message")
	}

	var numTargets int
	if committee.Len() > router.ScaleThresholdForClustering {
		numTargets = int(math.Sqrt(float64(committee.Len())))
	} else {
		numTargets = committee.Len()
	}
	indexes := rand.Perm(committee.Len())
	targets := make([]common.Address, 0, numTargets)
	for index := range indexes {
		val := committee.MemberByIndex(index)
		if val.Address != g.address {
			targets = append(targets, val.Address)
		}
		if len(targets) == numTargets {
			break
		}
	}

	// bail out if the local validator is the only one in the committee
	if len(targets) == 0 {
		return fmt.Errorf("no one to ask sync to")
	}

	ps := g.broadcaster.FindPeers(targets)
	// bail out if we cannot find any peers
	if len(ps) == 0 {
		return fmt.Errorf("cannot find any peers")
	}

	for addr, p := range ps {
		g.logger.Debug("Asking sync to", "addr", addr)
		go p.Send(message.SyncNetworkMsg, encoded) //nolint
	}
	return nil
}
