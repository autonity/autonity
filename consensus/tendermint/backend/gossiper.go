package backend

import (
	"math/big"
	"time"

	lru "github.com/hashicorp/golang-lru"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/log"
)

type Gossiper struct {
	recentMessages *lru.ARCCache  // the cache of peer's messages
	knownMessages  *lru.ARCCache  // the cache of self messages
	address        common.Address // address of the local peer
	broadcaster    consensus.Broadcaster
	logger         log.Logger
	stopped        chan struct{}
}

func NewGossiper(recentMessages *lru.ARCCache, knownMessages *lru.ARCCache, address common.Address, logger log.Logger, stopped chan struct{}) *Gossiper {
	return &Gossiper{
		recentMessages: recentMessages,
		knownMessages:  knownMessages,
		address:        address,
		logger:         logger,
		stopped:        stopped,
	}
}

func (g *Gossiper) SetBroadcaster(broadcaster consensus.Broadcaster) {
	g.broadcaster = broadcaster
}

func (g *Gossiper) Broadcaster() consensus.Broadcaster {
	return g.broadcaster
}

func (g *Gossiper) RecentMessages() *lru.ARCCache {
	return g.recentMessages
}

func (g *Gossiper) KnownMessages() *lru.ARCCache {
	return g.knownMessages
}

func (g *Gossiper) Address() common.Address {
	return g.address
}

func (g *Gossiper) Gossip(committee types.Committee, message message.Msg) {
	hash := message.Hash()
	g.knownMessages.Add(hash, true)
	targets := make(map[common.Address]struct{})
	for _, val := range committee {
		if val.Address != g.address {
			targets[val.Address] = struct{}{}
		}
	}
	if g.broadcaster != nil && len(targets) > 0 {
		ps := g.broadcaster.FindPeers(targets)
		for addr, p := range ps {
			ms, ok := g.recentMessages.Get(addr)
			var m *lru.ARCCache
			if ok {
				m, _ = ms.(*lru.ARCCache)
				if _, k := m.Get(hash); k {
					// This peer had this event, skip it
					continue
				}
			} else {
				m, _ = lru.NewARC(inmemoryMessages)
			}

			m.Add(hash, true)
			g.recentMessages.Add(addr, m)

			go p.SendRaw(NetworkCodes[message.Code()], message.Payload()) //nolint
		}
	}
}

func (g *Gossiper) AskSync(header *types.Header) {

	targets := make(map[common.Address]struct{})
	for _, val := range header.Committee {
		if val.Address != g.address {
			targets[val.Address] = struct{}{}
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
				if count.Cmp(bft.Quorum(header.TotalVotingPower())) >= 0 {
					break
				}
				g.logger.Debug("Asking sync to", "addr", addr)
				go p.Send(SyncNetworkMsg, []byte{}) //nolint

				member := header.CommitteeMember(addr)
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
