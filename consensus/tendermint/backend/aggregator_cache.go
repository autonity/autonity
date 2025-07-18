package backend

import (
	"math/big"
	"sync"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/events"
	"github.com/autonity/autonity/core/types"
)

type cacheStep int

const (
	stepReceived cacheStep = iota
	stepDispatched
)

type voteCache struct {
	mu       sync.RWMutex
	internal map[uint64]map[int64]map[common.Hash]*types.Bitmap
}

func newVoteCache() *voteCache {
	return &voteCache{
		internal: make(map[uint64]map[int64]map[common.Hash]*types.Bitmap),
	}
}

func (c *voteCache) contains(height uint64, round int64, value common.Hash, bm *types.Bitmap) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if _, ok := c.internal[height]; !ok {
		return false
	}
	if _, ok := c.internal[height][round]; !ok {
		return false
	}
	if _, ok := c.internal[height][round][value]; !ok {
		return false
	}
	known := c.internal[height][round][value]
	return known.Contains(bm)
}

func (c *voteCache) merge(height uint64, round int64, value common.Hash, bm *types.Bitmap) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.internal[height]; !ok {
		c.internal[height] = make(map[int64]map[common.Hash]*types.Bitmap)
	}
	if _, ok := c.internal[height][round]; !ok {
		c.internal[height][round] = make(map[common.Hash]*types.Bitmap)
	}

	if existing, ok := c.internal[height][round][value]; ok {
		c.internal[height][round][value] = existing.Merge(bm)
	} else {
		c.internal[height][round][value] = bm.Copy()
	}
}

func (c *voteCache) invalidate(height uint64, round int64, vote message.Vote) {
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

	toInvalidate := vote.Signers().Bitmap
	known := c.internal[height][round][vote.Value()]

	c.internal[height][round][vote.Value()] = known.AndNot(toInvalidate)
}

func (c *voteCache) pruneToHeight(height uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for h := range c.internal {
		if h < height {
			delete(c.internal, h)
		}
	}
}

func (c *voteCache) mergedVoters(height uint64, round int64) *types.Bitmap {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if _, ok := c.internal[height]; !ok {
		return nil
	}
	if _, ok := c.internal[height][round]; !ok {
		return nil
	}
	var merged *types.Bitmap
	for _, bm := range c.internal[height][round] {
		if merged == nil {
			merged = types.NewBitmap()
		}
		merged = merged.Merge(bm)
	}
	return merged
}

func (c *voteCache) presentPowerForRound(height uint64, round int64, committeePowers []*big.Int) *big.Int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	merged := c.mergedVoters(height, round)
	if merged == nil || merged.Count() == 0 {
		return big.NewInt(0)
	}
	power := big.NewInt(0)
	for i, p := range committeePowers {
		if merged.IsSet(i) {
			power.Add(power, p)
		}
	}
	return power
}

func (c *voteCache) presentPowerForValue(height uint64, round int64, value common.Hash, committeePowers []*big.Int) *big.Int {
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
	for i, p := range committeePowers {
		if bm.IsSet(i) {
			power.Add(power, p)
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
	voteCaches      map[uint8]map[cacheStep]*voteCache // code -> step -> vote cache

	filterMu sync.RWMutex
	filtered map[uint64]map[filteredCacheKey][]events.UnverifiedMessageEvent
}

func newAggregatorCache() *aggregatorCache {
	return &aggregatorCache{
		committeePowers: make(map[uint64][]*big.Int),
		voteCaches: map[uint8]map[cacheStep]*voteCache{
			message.PrecommitCode: {
				stepReceived:   newVoteCache(),
				stepDispatched: newVoteCache(),
			},
			message.PrevoteCode: {
				stepReceived:   newVoteCache(),
				stepDispatched: newVoteCache(),
			},
			message.ProposalCode: {
				stepReceived:   newVoteCache(),
				stepDispatched: newVoteCache(),
			},
		},

		filtered: make(map[uint64]map[filteredCacheKey][]events.UnverifiedMessageEvent),
	}
}

func (c *aggregatorCache) contains(height uint64, round int64, event events.UnverifiedMessageEvent) (received bool, dispatched bool) {
	msg := event.Message
	switch msg.(type) {
	case *message.Propose:
		bm := types.NewBitmap()
		bm.Set(msg.(*message.Propose).SignerIndex())
		received = c.voteCaches[message.ProposalCode][stepReceived].contains(height, round, msg.Value(), bm)
		dispatched = c.voteCaches[message.ProposalCode][stepDispatched].contains(height, round, msg.Value(), bm)
		return received, dispatched
	case *message.Precommit, *message.Prevote:
		bm := msg.(message.Vote).Signers().Bitmap
		received = c.voteCaches[msg.Code()][stepReceived].contains(height, round, msg.Value(), bm)
		dispatched = c.voteCaches[msg.Code()][stepDispatched].contains(height, round, msg.Value(), bm)
		return received, dispatched
	default:
		return false, false
	}
}

func (c *aggregatorCache) filter(event events.UnverifiedMessageEvent) bool {
	msg := event.Message
	switch msg.(type) {
	case *message.Precommit, *message.Prevote:
		received, dispatched := c.contains(msg.H(), msg.R(), event)
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

func (c *aggregatorCache) invalidateVotes(msg message.Vote) {
	c.voteCaches[msg.Code()][stepReceived].invalidate(msg.H(), msg.R(), msg)
}

func (c *aggregatorCache) emptyFiltered(h uint64, r int64, code uint8, value common.Hash) []events.UnverifiedMessageEvent {
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

func (c *aggregatorCache) markCommittee(height uint64, committee *types.Committee) {
	if len(c.committeePowers[height]) == committee.Len() {
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

func (c *aggregatorCache) addEvent(event events.UnverifiedMessageEvent, step cacheStep) {
	msg := event.Message
	switch msg.(type) {
	case *message.Propose:
		c.addProposal(msg.(*message.Propose), step)
	case *message.Precommit, *message.Prevote:
		c.addVote(msg.(message.Vote), step)
	default:
		panic("aggregatorCache: unsupported message type for addEvent")
	}
}

func (c *aggregatorCache) addProposal(msg *message.Propose, step cacheStep) {
	if _, ok := c.committeePowers[msg.H()]; !ok {
		panic("aggregatorCache: committee powers not set for height")
	}
	bm := types.NewBitmap()
	bm.Set(msg.SignerIndex())
	c.voteCaches[message.ProposalCode][step].merge(
		msg.H(),
		msg.R(),
		msg.Value(),
		bm,
	)
}

func (c *aggregatorCache) addVote(msg message.Vote, step cacheStep) {
	if _, ok := c.committeePowers[msg.H()]; !ok {
		panic("aggregatorCache: committee powers not set for height")
	}
	c.voteCaches[msg.Code()][step].merge(
		msg.H(),
		msg.R(),
		msg.Value(),
		msg.Signers().Bitmap,
	)
}

func (c *aggregatorCache) presentPowerForValue(height uint64, round int64, value common.Hash, code uint8, step cacheStep) *big.Int {
	if _, ok := c.committeePowers[height]; !ok {
		panic("aggregatorCache: committee powers not set for height")
	}
	committeePowers := c.committeePowers[height]
	return c.voteCaches[code][step].presentPowerForValue(height, round, value, committeePowers)
}

func (c *aggregatorCache) totalPowerForCode(height uint64, round int64, code uint8, step cacheStep) *big.Int {
	if _, ok := c.committeePowers[height]; !ok {
		panic("aggregatorCache: committee powers not set for height")
	}
	committeePowers := c.committeePowers[height]
	return c.voteCaches[code][step].presentPowerForRound(height, round, committeePowers)
}

// / Total power for a round, considering all vote types (precommit, prevote, proposal)
func (c *aggregatorCache) totalPowerForRound(height uint64, round int64, step cacheStep) *big.Int {
	if _, ok := c.committeePowers[height]; !ok {
		panic("aggregatorCache: committee powers not set for height")
	}
	merged := types.NewBitmap()
	for _, code := range []uint8{message.PrecommitCode, message.PrevoteCode, message.ProposalCode} {
		voters := c.voteCaches[code][step].mergedVoters(height, round)
		if voters != nil && voters.Count() != 0 {
			merged = merged.Merge(voters)
		}
	}
	if merged.Count() == 0 {
		return big.NewInt(0)
	}
	power := big.NewInt(0)
	for i, p := range c.committeePowers[height] {
		if merged.IsSet(i) {
			power.Add(power, p)
		}
	}
	return power
}

func (c *aggregatorCache) pruneToHeight(height uint64) {
	c.voteCaches[message.PrecommitCode][stepReceived].pruneToHeight(height)
	c.voteCaches[message.PrecommitCode][stepDispatched].pruneToHeight(height)

	c.voteCaches[message.PrevoteCode][stepReceived].pruneToHeight(height)
	c.voteCaches[message.PrevoteCode][stepDispatched].pruneToHeight(height)

	c.filterMu.Lock()
	defer c.filterMu.Unlock()

	for h := range c.committeePowers {
		if h < height {
			delete(c.committeePowers, h)
			delete(c.filtered, h)
		}
	}
}
