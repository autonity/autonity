package backend

import (
	"math"
	"math/big"
	"math/rand"
	"time"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/fixsizecache"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/router/constants"
	"github.com/autonity/autonity/consensus/tendermint/router/interfaces"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/log"
)

type msgRouter interface {
	Recipients(committee *types.Committee, msg message.Msg, from common.Address) ([]common.Address, error)
	SetBroadcaster(broadcaster interfaces.PeerFinder)
}

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

func (g *Gossiper) SlowGossip(committee *types.Committee, msg message.Msg) {
	// only gossip to very small committee
	numTargets := len(committee.Members)
	targetIndices := rand.Perm(numTargets) // target indices to select from the full committee
	if numTargets > constants.ScaleThresholdForClustering {
		numTargets = int(math.Sqrt(float64(numTargets)))
	}
	recipients := make([]common.Address, numTargets)
	for i := 0; i < numTargets; i++ {
		recipients[i] = committee.Members[targetIndices[i]].Address
	}
	g.gossip(msg, recipients)
}

func (g *Gossiper) Gossip(committee *types.Committee, msg message.Msg) {
	recipients, _ := g.router.Recipients(committee, msg, g.address)
	g.gossip(msg, recipients)
}

func (g *Gossiper) gossip(msg message.Msg, recipients []common.Address) {
	hash := msg.Hash()
	if msg.Originator() == g.address {
		g.knownMessages.Add(hash, true)
	}
	// if it's an aggregate the originator is the representative of the signers, so check in cache first and then add
	switch m := msg.(type) {
	case *message.Prevote:
	case *message.Precommit:
		if m.Signers().Len() > 1 && !g.knownMessages.Contains(hash) {
			g.knownMessages.Add(hash, true)
		}
	}

	if g.broadcaster == nil {
		return
	}
	code := message.NetworkCodes[msg.Code()]
	payload := msg.Payload()
	lostPeers := make([]common.Address, 0)
	for _, addr := range recipients {
		if addr == g.address {
			continue
		}
		if p, ok := g.broadcaster.FindPeer(addr); ok {
			if p.Cache().Contains(hash) {
				// This peer had this event, skip it
				continue
			}
			p.Cache().Add(hash, true)
			go p.SendRaw(code, payload) //nolint
		} else {
			// todo: Jason, shall we select other backups for liveness?
			lostPeers = append(lostPeers, addr)
		}
	}
	if len(lostPeers) > 0 {
		g.logger.Debug("peers not found", "len", len(lostPeers), "peers", lostPeers)
	}
}

func (g *Gossiper) AskSync(committee *types.Committee) {

	targets := make([]common.Address, 0, committee.Len())
	for _, val := range committee.Members {
		if val.Address != g.address {
			targets = append(targets, val.Address)
		}
	}

	if g.broadcaster != nil && len(targets) > 0 {
		for {
			ps := g.broadcaster.FindPeers(targets)
			// If we didn't find any peers try again in 10ms or exit if we have
			// been stopped.
			if len(ps) == 0 {
				t := time.NewTimer(retryPeriod * time.Millisecond)
				select {
				case <-t.C:
					continue
				case <-g.stopped:
					return
				}
			}
			count := new(big.Int)
			for addr, p := range ps {
				//ask to a quorum nodes to sync, 1 must then be honest and updated
				if count.Cmp(bft.Quorum(committee.TotalVotingPower())) >= 0 {
					break
				}
				g.logger.Debug("asking sync to", "addr", addr)
				go p.Send(message.SyncNetworkMsg, []byte{}) //nolint

				member := committee.MemberByAddress(addr)
				if member == nil {
					g.logger.Error("could not retrieve member from address")
					continue
				}
				count.Add(count, member.VotingPower)
			}
			break
		}
	}
}
