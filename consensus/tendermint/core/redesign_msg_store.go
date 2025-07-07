package core

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/log"
)

// need committee hight to make the pre-allocations - optional
type committeeProvider interface {
	CommitteeByHeight(height uint64) (*types.Committee, error)
}
type roundPowerCache struct {
	votedHash     common.Hash
	hashVotePower *message.AggregatedPower
	nilVotePower  *message.AggregatedPower
	// this is possible but not expected normally, it should be empty in most cases
	otherVotePower map[common.Hash]*message.AggregatedPower
}

type roundStore struct {
	sync.RWMutex
	proposals  []*message.Propose // should this be an array?
	prevotes   []*message.Prevote
	precommits []*message.Precommit
	powerCache *roundPowerCache

	// todo: (review) this is optional but adding it because it is most accessed, number of prevotes are generally high for large committees
	prevoteByValue map[common.Hash][]*message.Prevote
	nilPrevotes    []*message.Prevote

	votesBySigner map[int]map[common.Hash]message.Vote // prevotes and precommits
	seenHashes    map[common.Hash]struct{}             // common for prevote/precommit/proposal
}

func newRoundStore(committeeSize int) *roundStore {
	return &roundStore{
		RWMutex:    sync.RWMutex{},
		proposals:  make([]*message.Propose, 0, 1),
		prevotes:   make([]*message.Prevote, 0, int(float64(committeeSize)*(2.0/3.0))), // pre allocate minimum committee size
		precommits: make([]*message.Precommit, 0, committeeSize/10),
		powerCache: &roundPowerCache{
			hashVotePower:  message.NewAggregatedPower(),
			nilVotePower:   message.NewAggregatedPower(),
			otherVotePower: make(map[common.Hash]*message.AggregatedPower),
		},
		prevoteByValue: make(map[common.Hash][]*message.Prevote),
		nilPrevotes:    make([]*message.Prevote, 0, int(float64(committeeSize)/3.0)),

		votesBySigner: make(map[int]map[common.Hash]message.Vote),
		seenHashes:    make(map[common.Hash]struct{}),
	}

}

type heightStore struct {
	sync.RWMutex
	rounds        []*roundStore
	committeeSize int // the committee size at this height
	maxRoundSeen  int64
}

func newHeightStore(committeeSize int) *heightStore {
	return &heightStore{
		RWMutex:       sync.RWMutex{},
		rounds:        make([]*roundStore, 0, 1), // 1 round by default
		committeeSize: committeeSize,
		maxRoundSeen:  -1,
	}
}

func (hs *heightStore) getOrCreateRoundStore(round int64) *roundStore {
	hs.RLock()
	if int(round) < len(hs.rounds) && hs.rounds[round] != nil {
		rs := hs.rounds[round]
		hs.RUnlock()
		return rs
	}
	hs.RUnlock()

	// re lock the height store to check and create round
	hs.Lock()
	defer hs.Unlock()

	if round < int64(len(hs.rounds)) && hs.rounds[round] != nil {
		return hs.rounds[round]
	}
	if int(round) >= len(hs.rounds) {
		newRounds := make([]*roundStore, round+1)
		copy(newRounds, hs.rounds)
		hs.rounds = newRounds
	}
	rs := newRoundStore(hs.committeeSize)
	hs.rounds[round] = rs
	if round > hs.maxRoundSeen {
		hs.maxRoundSeen = round
	}
	return rs
}

type MsgStore struct {
	storeLock sync.RWMutex

	// the first height that msg are buffered from after node is start.
	firstHeight uint64
	store       map[uint64]*heightStore
	committeeProvider
}

func NewMsgStore() *MsgStore {
	return &MsgStore{
		storeLock:         sync.RWMutex{},
		firstHeight:       uint64(0),
		store:             make(map[uint64]*heightStore),
		committeeProvider: nil,
	}
}

func (ms *MsgStore) GetMaxRoundSeen(height uint64) int64 {
	ms.storeLock.RLock()
	defer ms.storeLock.RUnlock()
	hs, ok := ms.store[height]
	if !ok {
		return -1
	}
	return hs.maxRoundSeen
}

func (ms *MsgStore) HasVote(height uint64, round int64, hash common.Hash) bool {
	if round < 0 || round >= constants.MaxRound {
		return false
	}
	hs, err := ms.getOrCreateHeightStore(height)
	if err != nil {
		return false
	}
	rs := hs.getOrCreateRoundStore(round)
	rs.RLock()
	defer rs.RUnlock()
	_, ok := rs.seenHashes[hash]
	return ok
}

func (ms *MsgStore) GetSignerPrevotesByRound(height uint64, round int64, signerIndex int) map[common.Hash]*message.Prevote {
	if round < 0 || round >= constants.MaxRound {
		return nil
	}
	hs, err := ms.getOrCreateHeightStore(height)
	if err != nil {
		return nil
	}
	rs := hs.getOrCreateRoundStore(round)
	rs.RLock()
	defer rs.RUnlock()
	if votes, ok := rs.votesBySigner[signerIndex]; ok {
		result := make(map[common.Hash]*message.Prevote, len(votes))
		for k, v := range votes {
			if _, isVote := v.(*message.Prevote); isVote {
				result[k] = v.(*message.Prevote)
			}
		}
		return result
	}
	return nil
}

func (ms *MsgStore) GetSignerPrecommitsByRound(height uint64, round int64, signerIndex int) map[common.Hash]*message.Precommit {
	if round < 0 || round >= constants.MaxRound {
		return nil
	}
	hs, err := ms.getOrCreateHeightStore(height)
	if err != nil {
		return nil
	}
	rs := hs.getOrCreateRoundStore(round)
	rs.RLock()
	defer rs.RUnlock()
	if votes, ok := rs.votesBySigner[signerIndex]; ok {
		result := make(map[common.Hash]*message.Precommit, len(votes))
		for k, v := range votes {
			if _, isVote := v.(*message.Precommit); isVote {
				result[k] = v.(*message.Precommit)
			}
		}
		return result
	}
	return nil
}

func (ms *MsgStore) SetCommitteeProvider(provider committeeProvider) {
	if provider == nil {
		panic("committee provider cannot be nil")
	}
	ms.committeeProvider = provider
}

func (ms *MsgStore) getOrCreateHeightStore(height uint64) (*heightStore, error) {
	ms.storeLock.RLock()
	store, ok := ms.store[height]
	ms.storeLock.RUnlock()
	if ok {
		return store, nil
	}

	// relock, check and create
	ms.storeLock.Lock()
	defer ms.storeLock.Unlock()
	store, ok = ms.store[height]
	if ok {
		return store, nil
	}
	// Fetch the committee for this specific height to get its size.
	committee, err := ms.committeeProvider.CommitteeByHeight(height)
	if err != nil {
		return nil, fmt.Errorf("committee provider error: %s", err)
	}

	store = newHeightStore(committee.Len())
	ms.store[height] = store

	return store, nil
}

// Save messages
func (ms *MsgStore) Save(m message.Msg) error {
	height, round := m.H(), m.R()
	if round < 0 || round > constants.MaxRound {
		panic("round is out of bounds")
	}

	ms.storeLock.Lock()
	if ms.firstHeight == 0 {
		ms.firstHeight = height
	}
	ms.storeLock.Unlock()

	hs, err := ms.getOrCreateHeightStore(height)
	if err != nil {
		return err
	}
	rs := hs.getOrCreateRoundStore(round)

	rs.Lock()
	defer rs.Unlock()
	rs.seenHashes[m.Hash()] = struct{}{}

	switch msg := m.(type) {
	case *message.Propose:
		rs.proposals = append(rs.proposals, msg)
	case *message.Prevote:
		rs.prevotes = append(rs.prevotes, msg)
		ms.updatePrevotePower(rs, msg)
		ms.updateSignerIndex(rs, msg)
	case *message.Precommit:
		rs.precommits = append(rs.precommits, msg)
		ms.updateSignerIndex(rs, msg)
	}
	return nil
}

func (ms *MsgStore) updateSignerIndex(rs *roundStore, vote message.Vote) {
	vote.Signers().ForEachDistinctSigner(func(signerIndex int) {
		if _, ok := rs.votesBySigner[signerIndex]; !ok {
			rs.votesBySigner[signerIndex] = make(map[common.Hash]message.Vote)
		}
		rs.votesBySigner[signerIndex][vote.Value()] = vote
	}, vote.Signers().CommitteeSize())
}

func (ms *MsgStore) updatePrevotePower(rs *roundStore, msg *message.Prevote) {
	value := msg.Value()
	cache := rs.powerCache
	var targetPower *message.AggregatedPower
	if value == common.NilValue {
		targetPower = cache.nilVotePower
		rs.nilPrevotes = append(rs.nilPrevotes, msg)
	} else if cache.votedHash == (common.Hash{}) || value == cache.votedHash {
		if cache.votedHash == (common.Hash{}) { // first time we see a value, we prioritze this for power
			cache.votedHash = value
		}
		targetPower = cache.hashVotePower
		rs.prevoteByValue[value] = append(rs.prevoteByValue[value], msg)
	} else { // we have a different value, store it in the otherVotePower map
		if _, ok := cache.otherVotePower[value]; !ok {
			cache.otherVotePower[value] = message.NewAggregatedPower()
		}
		targetPower = cache.otherVotePower[value]
		rs.prevoteByValue[value] = append(rs.prevoteByValue[value], msg)
	}

	for index, power := range msg.Signers().Powers() {
		targetPower.Set(index, power)
	}
}

// todo: (review) can we eliminate the need for this method?
func (ms *MsgStore) RemoveMsg(height uint64, round int64, code uint8, hash common.Hash) {
	if round < 0 || round > constants.MaxRound {
		return
	}
	hs, _ := ms.getOrCreateHeightStore(height)
	rs := hs.getOrCreateRoundStore(round)
	rs.Lock()
	defer rs.Unlock()
	switch code {
	case message.PrevoteCode:
		remainingPrevotes := rs.prevotes[:0]
		wasRemoved := false
		for _, prevote := range rs.prevotes {
			if prevote.Hash() != hash {
				remainingPrevotes = append(remainingPrevotes, prevote)
			} else {
				wasRemoved = true
			}
		}
		if wasRemoved == false {
			return // nothing to update, no need to rebuild the round
		}
		// reset round data
		rs.prevotes = remainingPrevotes
		rs.powerCache = &roundPowerCache{
			nilVotePower:   message.NewAggregatedPower(),
			hashVotePower:  message.NewAggregatedPower(),
			otherVotePower: make(map[common.Hash]*message.AggregatedPower),
		}
		rs.prevoteByValue = make(map[common.Hash][]*message.Prevote)
		rs.nilPrevotes = rs.nilPrevotes[:0]
		for signerIndex, votes := range rs.votesBySigner {
			for value, vote := range votes {
				if _, ok := vote.(*message.Prevote); ok {
					delete(rs.votesBySigner[signerIndex], value)
				}
			}
		}
		delete(rs.seenHashes, hash)

		// rebuild the prevotes power
		for _, p := range rs.prevotes {
			ms.updatePrevotePower(rs, p)
			ms.updateSignerIndex(rs, p)
		}
	case message.PrecommitCode:
		filtered := rs.precommits[:0]
		for _, p := range rs.precommits {
			if p.Hash() != hash {
				filtered = append(filtered, p)
			}
		}
		rs.precommits = filtered
		for signerIndex, votes := range rs.votesBySigner {
			for value, vote := range votes {
				if _, ok := vote.(*message.Precommit); ok {
					delete(rs.votesBySigner[signerIndex], value)
				}
			}
		}

		for _, p := range rs.precommits {
			ms.updateSignerIndex(rs, p)
		}

	case message.ProposalCode:
		filtered := rs.proposals[:0]
		for _, p := range rs.proposals {
			if p.Hash() != hash {
				filtered = append(filtered, p)
			}
		}
		rs.proposals = filtered
	}
}

func (ms *MsgStore) GetPrevotesByRoundAndValue(height uint64, round int64, value common.Hash) ([]*message.Prevote, error) {
	if round < 0 || round > constants.MaxRound {
		return nil, fmt.Errorf("round is out of bounds")
	}
	hs, err := ms.getOrCreateHeightStore(height)
	if err != nil {
		return nil, err
	}
	rs := hs.getOrCreateRoundStore(round)
	rs.RLock()
	defer rs.RUnlock()
	var result []*message.Prevote
	if value == common.NilValue {
		result = make([]*message.Prevote, len(rs.nilPrevotes))
		copy(result, rs.nilPrevotes)
	} else if p, ok := rs.prevoteByValue[value]; ok {
		result = make([]*message.Prevote, len(p))
		copy(result, p)
	}
	return result, nil
}

func (ms *MsgStore) SearchQuorum(height uint64, round int64, excludedValue common.Hash, quorum *big.Int) []*message.Prevote {
	if round < 0 || round > constants.MaxRound {
		return nil
	}
	hs, err := ms.getOrCreateHeightStore(height)
	if err != nil {
		log.Error("SearchQuorum: failed to get height store", "height", height, "error", err)
		return nil
	}
	rs := hs.getOrCreateRoundStore(round)
	rs.RLock()
	defer rs.RUnlock()

	cache := rs.powerCache
	getPrevotesCopy := func(source []*message.Prevote) []*message.Prevote {
		if len(source) == 0 { // this should never hit
			return nil
		}
		result := make([]*message.Prevote, len(source))
		copy(result, source)
		return result
	}

	// find quorum for nil value
	if excludedValue != common.NilValue && cache.nilVotePower.Power().Cmp(quorum) >= 0 {
		return getPrevotesCopy(rs.nilPrevotes)
	}

	// find quorum for the primary voted hash
	if cache.votedHash != (common.Hash{}) && cache.votedHash != excludedValue && cache.hashVotePower.Power().Cmp(quorum) >= 0 {
		return getPrevotesCopy(rs.prevoteByValue[cache.votedHash])
	}

	// check quorum for other values
	for value, power := range cache.otherVotePower {
		if value != excludedValue && power.Power().Cmp(quorum) >= 0 {
			return getPrevotesCopy(rs.prevoteByValue[value])
		}
	}
	return nil
}

func (ms *MsgStore) GetProposalsByRound(height uint64, round int64, query func(*message.Propose) bool) []*message.Propose {
	if round < 0 || round > constants.MaxRound {
		return nil
	}
	hs, err := ms.getOrCreateHeightStore(height)
	if err != nil {
		return nil
	}
	rs := hs.getOrCreateRoundStore(round)
	rs.RLock()
	defer rs.RUnlock()
	var result []*message.Propose
	if query == nil {
		result = make([]*message.Propose, len(rs.proposals))
		copy(result, rs.proposals)
		return result
	}
	for _, proposal := range rs.proposals {
		if query(proposal) {
			result = append(result, proposal)
		}
	}

	return result
}

func (ms *MsgStore) GetPrevotesByRound(height uint64, round int64, query func(prevote *message.Prevote) bool) []*message.Prevote {
	if round < 0 || round > constants.MaxRound {
		return nil
	}
	hs, err := ms.getOrCreateHeightStore(height)
	if err != nil {
		return nil
	}
	rs := hs.getOrCreateRoundStore(round)
	rs.RLock()
	defer rs.RUnlock()
	var result []*message.Prevote
	if query == nil {
		result = make([]*message.Prevote, len(rs.prevotes))
		copy(result, rs.prevotes)
		return result
	}
	for _, prevote := range rs.prevotes {
		if query(prevote) {
			result = append(result, prevote)
		}
	}
	return result
}

func (ms *MsgStore) GetPrecommitsByRound(height uint64, round int64, query func(precommit *message.Precommit) bool) []*message.Precommit {
	if round < 0 || round > constants.MaxRound {
		return nil
	}
	hs, err := ms.getOrCreateHeightStore(height)
	if err != nil {
		return nil
	}
	rs := hs.getOrCreateRoundStore(round)
	rs.RLock()
	defer rs.RUnlock()
	var result []*message.Precommit
	if query == nil {
		result = make([]*message.Precommit, len(rs.precommits))
		copy(result, rs.precommits)
		return result
	}
	for _, precommit := range rs.precommits {
		if query(precommit) {
			result = append(result, precommit)
		}
	}
	return result
}

// GetProposals provides a generic query over all rounds for a given height.
func (ms *MsgStore) GetProposals(height uint64, query func(*message.Propose) bool) []*message.Propose {
	hs, err := ms.getOrCreateHeightStore(height)
	if err != nil {
		return nil
	}

	hs.RLock()
	defer hs.RUnlock()

	var result []*message.Propose
	for _, rs := range hs.rounds {
		if rs == nil {
			continue
		}
		rs.RLock()
		if query == nil {
			result = append(result, rs.proposals...)
		} else {
			for _, proposal := range rs.proposals {
				if query(proposal) {
					result = append(result, proposal)
				}
			}
		}
		rs.RUnlock()
	}
	return result
}

// GetPrevotes provides a generic query over all rounds for a given height.
func (ms *MsgStore) GetPrevotes(height uint64, query func(*message.Prevote) bool) []*message.Prevote {
	hs, err := ms.getOrCreateHeightStore(height)
	if err != nil {
		return nil
	}

	hs.RLock()
	defer hs.RUnlock()

	var result []*message.Prevote
	for _, rs := range hs.rounds {
		if rs == nil {
			continue
		}
		rs.RLock()
		if query == nil {
			result = append(result, rs.prevotes...)
		} else {
			for _, prevote := range rs.prevotes {
				if query(prevote) {
					result = append(result, prevote)
				}
			}
		}
		rs.RUnlock()
	}
	return result
}

// GetPrecommits provides a generic query over all rounds for a given height.
func (ms *MsgStore) GetPrecommits(height uint64, query func(*message.Precommit) bool) []*message.Precommit {
	hs, err := ms.getOrCreateHeightStore(height)
	if err != nil {
		return nil
	}

	hs.RLock()
	defer hs.RUnlock()

	var result []*message.Precommit
	for _, rs := range hs.rounds {
		if rs == nil {
			continue
		}
		rs.RLock()
		if query == nil {
			result = append(result, rs.precommits...)
		} else {
			for _, precommit := range rs.precommits {
				if query(precommit) {
					result = append(result, precommit)
				}
			}
		}
		rs.RUnlock()
	}
	return result
}

func (ms *MsgStore) PrevotesPowerFor(height uint64, round int64, value common.Hash) *big.Int {
	if round < 0 || round > constants.MaxRound {
		return new(big.Int)
	}
	hs, err := ms.getOrCreateHeightStore(height)
	if err != nil {
		return new(big.Int)
	}
	rs := hs.getOrCreateRoundStore(round)
	rs.RLock()
	defer rs.RUnlock()

	cache := rs.powerCache
	if value == common.NilValue {
		return new(big.Int).Set(cache.nilVotePower.Power())
	}
	if value == cache.votedHash {
		return new(big.Int).Set(cache.hashVotePower.Power())
	}
	if power, ok := cache.otherVotePower[value]; ok {
		return new(big.Int).Set(power.Power())
	}
	return new(big.Int) // return zero if no power found for the value
}

// FirstHeightBuffered returns the first height for which a message was stored.
func (ms *MsgStore) FirstHeightBuffered() uint64 {
	ms.storeLock.RLock()
	defer ms.storeLock.RUnlock()
	return ms.firstHeight
}

// DeleteOlds removes messages for heights less than the specified height.
func (ms *MsgStore) DeleteOlds(height uint64) {
	ms.storeLock.Lock()
	defer ms.storeLock.Unlock()
	for h := range ms.store {
		if h < height {
			delete(ms.store, h)
		}
	}
}
