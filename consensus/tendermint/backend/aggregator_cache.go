package backend

import (
	"math/big"
	"sync"

	"golang.org/x/exp/slices"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/events"
	"github.com/autonity/autonity/core/types"
)

type oneBitmap []byte

func (o oneBitmap) Contains(a oneBitmap) bool {
	if len(o) != len(a) {
		return false
	}
	for i := range o {
		if (o[i] & a[i]) != a[i] {
			return false
		}
	}
	return true
}

func (o oneBitmap) Merge(a oneBitmap) oneBitmap {
	if len(o) != len(a) {
		return nil
	}
	result := make(oneBitmap, len(o))
	for i := range o {
		result[i] = o[i] | a[i]
	}
	return result
}

func (o oneBitmap) Present(index int) bool {
	if index < 0 || index >= len(o)*8 {
		return false
	}
	byteIndex := index / 8
	bitIndex := index % 8
	return (o[byteIndex] & (1 << bitIndex)) != 0
}

type voteCache struct {
	mu       sync.RWMutex
	internal map[uint64]map[int64]map[common.Hash]oneBitmap
}

func newVoteCache() *voteCache {
	return &voteCache{
		internal: make(map[uint64]map[int64]map[common.Hash]oneBitmap),
	}
}

func (c *voteCache) Contains(height uint64, round int64, committeeSize int, vote message.Vote) bool {
	signers := vote.Signers()
	c.mu.RLock()
	defer c.mu.RUnlock()

	if _, ok := c.internal[height]; !ok {
		return false
	}
	if _, ok := c.internal[height][round]; !ok {
		return false
	}
	if _, ok := c.internal[height][round][vote.Value()]; !ok {
		return false
	}
	bm := oneBitmap(signers.Bits.ToSingleBitmap(committeeSize))
	known := c.internal[height][round][vote.Value()]
	return known.Contains(bm)
}

func (c *voteCache) Merge(height uint64, round int64, committeeSize int, vote message.Vote) {
	signers := vote.Signers()
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.internal[height]; !ok {
		c.internal[height] = make(map[int64]map[common.Hash]oneBitmap)
	}
	if _, ok := c.internal[height][round]; !ok {
		c.internal[height][round] = make(map[common.Hash]oneBitmap)
	}

	bm := oneBitmap(signers.Bits.ToSingleBitmap(committeeSize))
	if existing, ok := c.internal[height][round][vote.Value()]; ok {
		c.internal[height][round][vote.Value()] = existing.Merge(bm)
	} else {
		c.internal[height][round][vote.Value()] = bm
	}
}

func (c *voteCache) PruneToHeight(height uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for h := range c.internal {
		if h < height {
			delete(c.internal, h)
		}
	}
}

func (c *voteCache) PresentPower(height uint64, round int64, committee *types.Committee) *big.Int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if _, ok := c.internal[height]; !ok {
		return big.NewInt(0)
	}
	if _, ok := c.internal[height][round]; !ok {
		return big.NewInt(0)
	}
	known := c.internal[height][round]
	powers := make([]*big.Int, len(known))
	i := 0
	for _, bm := range known {
		for _, member := range committee.Members {
			if bm.Present(int(member.Index)) {
				if powers[i] == nil {
					powers[i] = big.NewInt(0)
				}
				powers[i].Add(powers[i], member.VotingPower)
			}
		}
		i++
	}
	return slices.MaxFunc(powers, func(a, b *big.Int) int {
		return a.Cmp(b)
	})
}

type aggregatorCache struct {
	committeeSize  map[uint64]int // the committee size used to create the bitmaps
	precommitCache *voteCache
	prevoteCache   *voteCache
}

func newAggregatorCache() *aggregatorCache {
	return &aggregatorCache{
		committeeSize:  make(map[uint64]int),
		precommitCache: newVoteCache(),
		prevoteCache:   newVoteCache(),
	}
}

func (c *aggregatorCache) Contains(height uint64, round int64, committeeSize int, event events.UnverifiedMessageEvent) bool {
	msg := event.Message
	switch msg.(type) {
	case *message.Precommit:
		return c.precommitCache.Contains(height, round, committeeSize, msg.(message.Vote))
	case *message.Prevote:
		return c.prevoteCache.Contains(height, round, committeeSize, msg.(message.Vote))
	default:
		return false
	}
}

func (c *aggregatorCache) MarkCommitteeSize(height uint64, size int) {
	c.committeeSize[height] = size
}

func (c *aggregatorCache) AddPrevote(vote *message.Prevote) {
	c.prevoteCache.Merge(vote.H(), vote.R(), c.committeeSize[vote.H()], vote)
}

func (c *aggregatorCache) AddPrecommit(vote *message.Precommit) {
	c.precommitCache.Merge(vote.H(), vote.R(), c.committeeSize[vote.H()], vote)
}

func (c *aggregatorCache) PruneToHeight(height uint64) {
	c.precommitCache.PruneToHeight(height)
	c.prevoteCache.PruneToHeight(height)
	for h := range c.committeeSize {
		if h < height {
			delete(c.committeeSize, h)
		}
	}
}
