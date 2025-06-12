package backend

import (
	"fmt"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/fixsizecache"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/rlp"
)

type Gossiper struct {
	knownMessages *fixsizecache.Cache[common.Hash, bool] // the cache of self messages
	address       common.Address                         // address of the local peer
	broadcaster   consensus.Broadcaster
	logger        log.Logger
	stopped       chan struct{}
}

func NewGossiper(knownMessages *fixsizecache.Cache[common.Hash, bool], address common.Address, logger log.Logger, stopped chan struct{}) *Gossiper {
	return &Gossiper{
		knownMessages: knownMessages,
		address:       address,
		logger:        logger,
		stopped:       stopped,
	}
}

func (g *Gossiper) SetBroadcaster(broadcaster consensus.Broadcaster) {
	g.broadcaster = broadcaster
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

func (g *Gossiper) Gossip(committee *types.Committee, msg message.Msg) {
	hash := msg.Hash()
	if !g.knownMessages.Contains(hash) {
		g.knownMessages.Add(hash, true)
	}
	if g.broadcaster == nil {
		return
	}

	// forward future epoch proposal to all the committee members, as most of them are still in the committee.
	recipients := committee.Members
	code := message.NetworkCodes[msg.Code()]
	for _, val := range recipients {
		if val.Address == g.address {
			continue
		}
		if p, ok := g.broadcaster.FindPeer(val.Address); ok {
			if p.Cache().Contains(hash) {
				// This peer had this event, skip it
				continue
			}
			p.Cache().Add(hash, true)
			go p.SendRaw(code, msg.Payload()) //nolint
		}
	}
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
