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

type oneBitmap []byte

func newOneBitmap(size int) oneBitmap {
	if size <= 0 {
		return nil
	}
	bytes := (size + 7) / 8 // Calculate the number of bytes needed
	return make(oneBitmap, bytes)
}

func newOneBitmapProposal(size int, proposerIndex int) oneBitmap {
	if size <= 0 || proposerIndex < 0 || proposerIndex >= size {
		return nil
	}
	bm := newOneBitmap(size)
	byteIndex := proposerIndex / 8
	bitIndex := 7 - (proposerIndex % 8) // Convert to big-endian bit index
	bm[byteIndex] = 1 << bitIndex
	return bm
}

func (o oneBitmap) contains(a oneBitmap) bool {
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

func (o oneBitmap) merge(a oneBitmap) oneBitmap {
	if len(o) != len(a) {
		return nil
	}
	result := make(oneBitmap, len(o))
	for i := range o {
		result[i] = o[i] | a[i]
	}
	return result
}

func (o oneBitmap) invalidate(b oneBitmap) oneBitmap {
	if len(o) != len(b) {
		return nil
	}
	result := make(oneBitmap, len(o))
	for i := range o {
		result[i] = o[i] & (^b[i])
	}
	return result
}

func (o oneBitmap) presentIndexes() []int {
	var indexes []int
	for byteIndex, b := range o {
		for bitIndex := 0; b != 0; bitIndex++ {
			if b&0x80 != 0 { // Check if the most significant bit is set
				indexes = append(indexes, byteIndex*8+bitIndex)
			}
			b <<= 1 // Shift left to check the next bit
		}
	}
	return indexes
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

func (c *voteCache) contains(height uint64, round int64, value common.Hash, bm oneBitmap) bool {
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
	return known.contains(bm)
}

func (c *voteCache) merge(height uint64, round int64, value common.Hash, bm oneBitmap) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.internal[height]; !ok {
		c.internal[height] = make(map[int64]map[common.Hash]oneBitmap)
	}
	if _, ok := c.internal[height][round]; !ok {
		c.internal[height][round] = make(map[common.Hash]oneBitmap)
	}

	if existing, ok := c.internal[height][round][value]; ok {
		c.internal[height][round][value] = existing.merge(bm)
	} else {
		c.internal[height][round][value] = bm
	}
}

func (c *voteCache) invalidate(height uint64, round int64, committeeSize int, vote message.Vote) {
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

	c.internal[height][round][vote.Value()] = known.invalidate(toInvalidate)
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

func (c *voteCache) mergedVoters(height uint64, round int64) oneBitmap {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if _, ok := c.internal[height]; !ok {
		return nil
	}
	if _, ok := c.internal[height][round]; !ok {
		return nil
	}
	var merged oneBitmap
	for _, bm := range c.internal[height][round] {
		if merged == nil {
			merged = make(oneBitmap, len(bm))
		}
		merged = merged.merge(bm)
	}
	return merged
}

func (c *voteCache) presentPowerForRound(height uint64, round int64, committeePowers []*big.Int) *big.Int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	merged := c.mergedVoters(height, round)
	if len(merged) == 0 {
		return big.NewInt(0)
	}
	power := big.NewInt(0)
	for _, index := range merged.presentIndexes() {
		power.Add(power, committeePowers[index])
	}
	return power
}

func (c *voteCache) presentPowerForValue(height uint64, round int64, value common.Hash, committeePower []*big.Int) *big.Int {
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
	for _, index := range bm.presentIndexes() {
		power.Add(power, committeePower[index])
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

func (c *aggregatorCache) contains(height uint64, round int64, committeeSize int, event events.UnverifiedMessageEvent) (received bool, dispatched bool) {
	msg := event.Message
	switch msg.(type) {
	case *message.Propose:
		bm := newOneBitmapProposal(committeeSize, msg.(*message.Propose).SignerIndex())
		received = c.voteCaches[message.ProposalCode][stepReceived].contains(height, round, msg.Value(), bm)
		dispatched = c.voteCaches[message.ProposalCode][stepDispatched].contains(height, round, msg.Value(), bm)
		return received, dispatched
	case *message.Precommit, *message.Prevote:
		bm := msg.(message.Vote).Signers().Bits.ToSingleBitmap(committeeSize)
		received = c.voteCaches[msg.Code()][stepReceived].contains(height, round, msg.Value(), bm)
		dispatched = c.voteCaches[msg.Code()][stepDispatched].contains(height, round, msg.Value(), bm)
		return received, dispatched
	default:
		return false, false
	}
}

func (c *aggregatorCache) filter(committeeSize int, event events.UnverifiedMessageEvent) bool {
	msg := event.Message
	switch msg.(type) {
	case *message.Precommit, *message.Prevote:
		received, dispatched := c.contains(msg.H(), msg.R(), committeeSize, event)
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
	c.voteCaches[msg.Code()][stepReceived].invalidate(msg.H(), msg.R(), len(c.committeePowers[msg.H()]), msg)
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
	c.voteCaches[message.ProposalCode][step].merge(
		msg.H(),
		msg.R(),
		msg.Value(),
		newOneBitmapProposal(len(c.committeePowers[msg.H()]), msg.SignerIndex()),
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
		msg.Signers().Bits.ToSingleBitmap(len(c.committeePowers[msg.H()])),
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
	merged := newOneBitmap(len(c.committeePowers[height]))
	for _, code := range []uint8{message.PrecommitCode, message.PrevoteCode, message.ProposalCode} {
		voters := c.voteCaches[code][step].mergedVoters(height, round)
		if len(voters) != 0 {
			merged = merged.merge(voters)
		}
	}
	if len(merged) == 0 {
		return big.NewInt(0)
	}
	power := big.NewInt(0)
	for _, index := range merged.presentIndexes() {
		power.Add(power, c.committeePowers[height][index])
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
