package backend

import (
	"math/big"
	"time"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/fixsizecache"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/log"
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

func (g *Gossiper) Gossip(committee *types.Committee, message message.Msg) {
	hash := message.Hash()
	if !g.knownMessages.Contains(hash) {
		g.knownMessages.Add(hash, true)
	}
	if g.broadcaster == nil {
		return
	}
	code := NetworkCodes[message.Code()]
	payload := message.Payload()
	for _, val := range committee.Members {
		if val.Address == g.address {
			continue
		}
		if p, ok := g.broadcaster.FindPeer(val.Address); ok {
			if p.Cache().Contains(hash) {
				// This peer had this event, skip it
				continue
			}
			p.Cache().Add(hash, true)
			go p.SendRaw(code, payload) //nolint
		}
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
				g.logger.Debug("Asking sync to", "addr", addr)
				go p.Send(SyncNetworkMsg, []byte{}) //nolint

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
