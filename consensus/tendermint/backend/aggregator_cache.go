package backend

import (
	"math/big"
	"sync"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/events"
)

type CacheStep int

const (
	StepReceived CacheStep = iota
	StepDispatched
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
	bitIndex = 7 - bitIndex
	return (o[byteIndex] & (1 << bitIndex)) != 0
}

func (o oneBitmap) Invalidate(b oneBitmap) oneBitmap {
	if len(o) != len(b) {
		return nil
	}
	result := make(oneBitmap, len(o))
	for i := range o {
		result[i] = o[i] & (^b[i])
	}
	return result
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

func (c *voteCache) Invalidate(height uint64, round int64, committeeSize int, vote message.Vote) {
	signers := vote.Signers()
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.internal[height]; !ok {
		return
	}
	if _, ok := c.internal[height][round]; !ok {
		return
	}
	if _, ok := c.internal[height][round][vote.Value()]; !ok {
		return
	}

	toInvalidate := oneBitmap(signers.Bits.ToSingleBitmap(committeeSize))
	known := c.internal[height][round][vote.Value()]

	c.internal[height][round][vote.Value()] = known.Invalidate(toInvalidate)
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

func (c *voteCache) PresentPower(height uint64, round int64, committeePowers []*big.Int) *big.Int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if _, ok := c.internal[height]; !ok {
		return big.NewInt(0)
	}
	if _, ok := c.internal[height][round]; !ok {
		return big.NewInt(0)
	}
	known := c.internal[height][round]
	power := big.NewInt(0)
	for _, bm := range known {
		for i, memberPower := range committeePowers {
			if bm.Present(i) {
				power.Add(power, memberPower)
			}
		}
	}
	return power
}

func (c *voteCache) PresentPowerFor(height uint64, round int64, value common.Hash, committeePower []*big.Int) *big.Int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if _, ok := c.internal[height]; !ok {
		return big.NewInt(0)
	}
	if _, ok := c.internal[height][round]; !ok {
		return big.NewInt(0)
	}
	if _, ok := c.internal[height][round][value]; !ok {
		return big.NewInt(0)
	}

	bm := c.internal[height][round][value]
	power := big.NewInt(0)
	for i, memberPower := range committeePower {
		if bm.Present(i) {
			power = power.Add(power, memberPower)
		}
	}
	return power
}

type filteredCacheKey struct {
	code  uint8
	round int64
	value common.Hash
}

type aggregatorCache struct {
	committeePowers map[uint64][]*big.Int              // the committee size used to create the bitmaps
	voteCaches      map[uint8]map[CacheStep]*voteCache // code -> step -> vote cache

	filterMu sync.RWMutex
	filtered map[uint64]map[filteredCacheKey][]events.UnverifiedMessageEvent
}

func newAggregatorCache() *aggregatorCache {
	return &aggregatorCache{
		committeePowers: make(map[uint64][]*big.Int),
		voteCaches: map[uint8]map[CacheStep]*voteCache{
			message.PrecommitCode: {
				StepReceived:   newVoteCache(),
				StepDispatched: newVoteCache(),
			},
			message.PrevoteCode: {
				StepReceived:   newVoteCache(),
				StepDispatched: newVoteCache(),
			},
		},

		filtered: make(map[uint64]map[filteredCacheKey][]events.UnverifiedMessageEvent),
	}
}

func (c *aggregatorCache) Contains(height uint64, round int64, committeeSize int, event events.UnverifiedMessageEvent) (received bool, dispatched bool) {
	msg := event.Message
	switch msg.(type) {
	case *message.Precommit, *message.Prevote:
		received = c.voteCaches[msg.Code()][StepReceived].Contains(height, round, committeeSize, msg.(message.Vote))
		dispatched = c.voteCaches[msg.Code()][StepDispatched].Contains(height, round, committeeSize, msg.(message.Vote))
		return received, dispatched
	default:
		return false, false
	}
}

func (c *aggregatorCache) Filter(committeeSize int, event events.UnverifiedMessageEvent) bool {
	msg := event.Message
	switch msg.(type) {
	case *message.Precommit, *message.Prevote:
		received, dispatched := c.Contains(msg.H(), msg.R(), committeeSize, event)
		if dispatched {
			return true
		} else if received && !dispatched {
			c.filterMu.Lock()
			defer c.filterMu.Unlock()
			if _, ok := c.filtered[msg.H()]; !ok {
				c.filtered[msg.H()] = make(map[filteredCacheKey][]events.UnverifiedMessageEvent)
			}

			key := filteredCacheKey{
				code:  msg.Code(),
				round: msg.R(),
				value: msg.Value(),
			}
			c.filtered[msg.H()][key] = append(c.filtered[msg.H()][key], event)
			return true
		}
		return received && dispatched
	default:
		return false
	}
}

func (c *aggregatorCache) InvalidateVote(msg message.Vote) {
	c.voteCaches[msg.Code()][StepReceived].Invalidate(msg.H(), msg.R(), len(c.committeePowers[msg.H()]), msg)
}

func (c *aggregatorCache) EmptyFiltered(h uint64, r int64, code uint8, value common.Hash) []events.UnverifiedMessageEvent {
	key := filteredCacheKey{
		code:  code,
		round: r,
		value: value,
	}
	c.filterMu.Lock()
	defer c.filterMu.Unlock()
	if _, ok := c.filtered[h]; !ok {
		return nil
	}
	filtered := c.filtered[h][key]
	// remove the key from the cache
	delete(c.filtered[h], key)
	return filtered
}

func (c *aggregatorCache) MarkCommittee(height uint64, committee *types.Committee) {
	if len(c.committeePowers) == committee.Len() {
		return // already marked
	}
	c.committeePowers[height] = func() []*big.Int {
		powers := make([]*big.Int, committee.Len())
		for i, member := range committee.Members {
			powers[i] = new(big.Int).Set(member.VotingPower)
		}
		return powers
	}()
}

func (c *aggregatorCache) AddVote(msg message.Vote, step CacheStep) {
	if _, ok := c.committeePowers[msg.H()]; !ok {
		panic("aggregatorCache: committee powers not set for height")
	}
	c.voteCaches[msg.Code()][step].Merge(msg.H(), msg.R(), len(c.committeePowers[msg.H()]), msg)
}

func (c *aggregatorCache) PresentPower(height uint64, round int64, code uint8, step CacheStep) *big.Int {
	if _, ok := c.committeePowers[height]; !ok {
		panic("aggregatorCache: committee powers not set for height")
	}
	committeePowers := c.committeePowers[height]
	return c.voteCaches[code][step].PresentPower(height, round, committeePowers)
}

func (c *aggregatorCache) PresentPowerFor(height uint64, round int64, value common.Hash, code uint8, step CacheStep) *big.Int {
	if _, ok := c.committeePowers[height]; !ok {
		panic("aggregatorCache: committee powers not set for height")
	}
	committeePowers := c.committeePowers[height]
	return c.voteCaches[code][step].PresentPowerFor(height, round, value, committeePowers)
}

func (c *aggregatorCache) PruneToHeight(height uint64) {
	c.voteCaches[message.PrecommitCode][StepReceived].PruneToHeight(height)
	c.voteCaches[message.PrecommitCode][StepDispatched].PruneToHeight(height)

	c.voteCaches[message.PrevoteCode][StepReceived].PruneToHeight(height)
	c.voteCaches[message.PrevoteCode][StepDispatched].PruneToHeight(height)

	c.filterMu.Lock()
	defer c.filterMu.Unlock()

	for h := range c.committeePowers {
		if h < height {
			delete(c.committeePowers, h)
			delete(c.filtered, h)
		}
	}
}
