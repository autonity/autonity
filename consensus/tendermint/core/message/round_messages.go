package message

import (
	"math/big"
	"sync"

	"github.com/autonity/autonity/common"
)

type Map struct {
	internal map[int64]*RoundMessages
	sync.RWMutex
}

func NewMap() *Map {
	return &Map{
		internal: make(map[int64]*RoundMessages),
	}
}

func (s *Map) Reset() {
	s.Lock()
	defer s.Unlock()
	s.internal = make(map[int64]*RoundMessages)
}

func (s *Map) GetOrCreate(round int64) *RoundMessages {
	s.Lock()
	defer s.Unlock()
	state, ok := s.internal[round]
	if ok {
		return state
	}
	state = NewRoundMessages()
	s.internal[round] = state
	return state
}

func (s *Map) All() []Msg {
	s.RLock()
	defer s.RUnlock()

	messages := make([][]Msg, len(s.internal))
	var totalLen int
	i := 0
	for _, state := range s.internal {
		messages[i] = state.AllMessages()
		totalLen += len(messages[i])
		i++
	}
	result := make([]Msg, 0, totalLen)
	for _, ms := range messages {
		result = append(result, ms...)
	}

	return result
}

func (s *Map) GetRounds() []int64 {
	s.RLock()
	defer s.RUnlock()

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
	sync.RWMutex
}

// we need a reference to proposal also for proposing a proposal with vr!=0 if needed
func NewRoundMessages() *RoundMessages {
	return &RoundMessages{
		prevotes:         NewSet[*Prevote](),
		precommits:       NewSet[*Precommit](),
		verifiedProposal: false,
	}
}

func (s *RoundMessages) SetProposal(proposal *Propose, verified bool) {
	s.Lock()
	defer s.Unlock()
	s.proposal = proposal
	s.verifiedProposal = verified
}

func (s *RoundMessages) PrevotesPower(hash common.Hash) *big.Int {
	return s.prevotes.PowerFor(hash)
}

func (s *RoundMessages) PrevotesTotalPower() *big.Int {
	return s.prevotes.TotalPower()
}

func (s *RoundMessages) PrecommitsPower(hash common.Hash) *big.Int {
	return s.precommits.PowerFor(hash)
}

func (s *RoundMessages) PrecommitsTotalPower() *big.Int {
	return s.precommits.TotalPower()
}

func (s *RoundMessages) AddPrevote(prevote *Prevote) {
	s.prevotes.Add(prevote)
}

func (s *RoundMessages) AllPrevotes() []Msg {
	return s.prevotes.Messages()
}

func (s *RoundMessages) AllPrecommits() []Msg {
	return s.precommits.Messages()
}

func (s *RoundMessages) AddPrecommit(precommit *Precommit) {
	s.precommits.Add(precommit)
}

func (s *RoundMessages) PrecommitsFor(hash common.Hash) []*Precommit {
	return s.precommits.VotesFor(hash)
}

func (s *RoundMessages) Proposal() *Propose {
	s.RLock()
	defer s.RUnlock()
	return s.proposal
}

func (s *RoundMessages) ProposalHash() common.Hash {
	s.RLock()
	defer s.RUnlock()
	if s.proposal == nil {
		return common.Hash{}
	}
	return s.proposal.block.Hash()
}

func (s *RoundMessages) IsProposalVerified() bool {
	s.RLock()
	defer s.RUnlock()
	return s.verifiedProposal
}

func (s *RoundMessages) AllMessages() []Msg {
	s.RLock()
	defer s.RUnlock()

	prevotes := s.prevotes.Messages()
	precommits := s.precommits.Messages()

	result := make([]Msg, 0, len(prevotes)+len(precommits)+1)
	if s.proposal != nil {
		result = append(result, Msg(s.proposal))
	}

	result = append(result, prevotes...)
	result = append(result, precommits...)
	return result
}
