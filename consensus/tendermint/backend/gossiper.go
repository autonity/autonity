package backend

import (
	"context"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/core/types"
	lru "github.com/hashicorp/golang-lru"
)

type Gossiper struct {
	recentMessages *lru.ARCCache  // the cache of peer's messages
	knownMessages  *lru.ARCCache  // the cache of self messages
	address        common.Address // address of the local peer
	broadcaster    consensus.Broadcaster
}

func NewGossiper(recentMessages *lru.ARCCache, knownMessages *lru.ARCCache, address common.Address) *Gossiper {
	return &Gossiper{
		recentMessages: recentMessages,
		knownMessages:  knownMessages,
		address:        address,
	}
}

func (g *Gossiper) SetBroadcaster(broadcaster consensus.Broadcaster) {
	g.broadcaster = broadcaster
}

func (g *Gossiper) Gossip(_ context.Context, committee types.Committee, payload []byte) {
	hash := types.RLPHash(payload)
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

			go p.Send(TendermintMsg, payload) //nolint
		}
	}
}
