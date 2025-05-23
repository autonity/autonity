package backend

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/fixsizecache"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/rlp"
)

type msgRouter interface {
	Recipients(committee *types.Committee, msg message.Msg, from common.Address) ([]common.Address, error)
	SetBroadcaster(broadcaster consensus.Broadcaster)
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
	if numTargets > 10 { // todo: minimum nodes to start slow gossip
		numTargets = int(math.Sqrt(float64(len(committee.Members))))
	} else if numTargets == 0 {
		log.Error("no target to slow gossip", "num", len(committee.Members), "numTargets", numTargets, "committee", committee.Members)
		return
	}
	targetIndices := rand.Perm(numTargets)
	recipients := make([]common.Address, numTargets)
	for i := 0; i < numTargets; i++ {
		recipients[i] = committee.Members[targetIndices[i]].Address
	}
	g.gossip(msg, recipients)
}

func (g *Gossiper) Gossip(committee *types.Committee, msg message.Msg) {
	recipients, err := g.router.Recipients(committee, msg, g.address)
	if err != nil {
		log.Debug("no recipients for message from msgRouter, broadcast", "error", err, "height", msg.H(), "message type", msg.Code())
		// forward future epoch proposal to all the committee members, as most of them are still in the committee.
		recipients := make([]common.Address, 0, committee.Len())
		for _, val := range committee.Members {
			if val.Address != g.address {
				recipients = append(recipients, val.Address)
			}
		}
	}
	g.gossip(msg, recipients)
}

func (g *Gossiper) gossip(msg message.Msg, recipients []common.Address) {
	hash := msg.Hash()
	if msg.Originator() == g.address {
		g.knownMessages.Add(hash, true)
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
			lostPeers = append(lostPeers, addr)
		}
	}
	if len(lostPeers) > 0 {
		g.logger.Debug("peers not found", "len", len(lostPeers), "peers", lostPeers)
	}
}

func (g *Gossiper) Gossip(committee *types.Committee, msg message.Msg) {
	recipients, err := g.router.Route(committee, msg, g.address)
	if err != nil {
		log.Debug("Gossiper: No recipients for message from msgRouter, broadcast", "error", err, "height", msg.H(), "message type", msg.Code())
		// forward future epoch proposal to all the committee members, as most of them are still in the committee.
		recipients := make([]common.Address, 0, committee.Len())
		for _, val := range committee.Members {
			if val.Address != g.address {
				recipients = append(recipients, val.Address)
			}
		}
	}
	g.gossip(msg, recipients)
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

	// send to everyone except ourselves
	targets := make([]common.Address, 0, committee.Len())
	for _, val := range committee.Members {
		if val.Address != g.address {
			targets = append(targets, val.Address)
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
