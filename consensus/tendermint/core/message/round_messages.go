package message

import (
	"github.com/autonity/autonity/common"
	"math/big"
	"sync"
)

type Map struct {
	internal map[int64]*RoundMessages
	mu       sync.RWMutex
}

func NewMap() *Map {
	return &Map{
		internal: make(map[int64]*RoundMessages),
	}
}

func (s *Map) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.internal = make(map[int64]*RoundMessages)
}

func (s *Map) GetOrCreate(round int64) *RoundMessages {
	s.mu.Lock()
	defer s.mu.Unlock()
	state, ok := s.internal[round]
	if ok {
		return state
	}
	state = NewRoundMessages()
	s.internal[round] = state
	return state
}

func (s *Map) All() []Message {
	s.mu.RLock()
	defer s.mu.RUnlock()

	messages := make([][]Message, len(s.internal))
	var totalLen int
	i := 0
	for _, state := range s.internal {
		messages[i] = state.AllMessages()
		totalLen += len(messages[i])
		i++
	}
	result := make([]Message, 0, totalLen)
	for _, ms := range messages {
		result = append(result, ms...)
	}

	return result
}

func (s *Map) GetRounds() []int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rounds := make([]int64, 0, len(s.internal))
	for r := range s.internal {
		rounds = append(rounds, r)
	}

	return rounds
}

// RoundMessages stores all message received for a specific round.
type RoundMessages struct {
	verifiedProposal bool
	proposal         *Propose
	prevotes         *Set[*Prevote]
	precommits       *Set[*Precommit]
	mu               sync.RWMutex
}

// NewRoundMessages creates a new messages instance with the given view and validatorSet
// we need to keep a reference of proposal in order to propose locked proposal when there is a lock and itself is the proposer
func NewRoundMessages() *RoundMessages {
	return &RoundMessages{
		prevotes:         NewSet[*Prevote](),
		precommits:       NewSet[*Precommit](),
		verifiedProposal: false,
	}
}

func (s *RoundMessages) SetProposal(proposal *Propose, verified bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.proposal = proposal
	s.verifiedProposal = verified
}

func (s *RoundMessages) PrevotesPower(hash common.Hash) *big.Int {
	return s.prevotes.VotePower(hash)
}

func (s *RoundMessages) PrevotesTotalPower() *big.Int {
	return s.prevotes.TotalVotePower()
}

func (s *RoundMessages) PrecommitsPower(hash common.Hash) *big.Int {
	return s.precommits.VotePower(hash)
}

func (s *RoundMessages) PrecommitsTotalPower() *big.Int {
	return s.precommits.TotalVotePower()
}

func (s *RoundMessages) AddPrevote(prevote *Prevote) {
	s.prevotes.AddVote(prevote)
}

func (s *RoundMessages) AllPrevotes() []Message {
	return s.prevotes.Messages()
}

func (s *RoundMessages) AllPrecommits() []Message {
	return s.precommits.Messages()
}

func (s *RoundMessages) AddPrecommit(precommit *Precommit) {
	s.precommits.AddVote(precommit)
}

func (s *RoundMessages) PrecommitsFor(hash common.Hash) []*Precommit {
	return s.precommits.VotesFor(hash)
}

func (s *RoundMessages) Proposal() *Propose {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.proposal
}

func (s *RoundMessages) ProposalHash() common.Hash {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.proposal == nil {
		return common.Hash{}
	}
	return s.proposal.block.Hash()
}

func (s *RoundMessages) IsProposalVerified() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.verifiedProposal
}

func (s *RoundMessages) AllMessages() []Message {
	s.mu.RLock()
	defer s.mu.RUnlock()

	prevotes := s.prevotes.Messages()
	precommits := s.precommits.Messages()

	result := make([]Message, 0, len(prevotes)+len(precommits)+1)
	if s.proposal != nil {
		result = append(result, Message(s.proposal))
	}

	result = append(result, prevotes...)
	result = append(result, precommits...)
	return result
}
