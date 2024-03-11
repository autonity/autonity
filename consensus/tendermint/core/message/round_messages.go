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

// TODO(lorenzo) refinements, this function has a mutex that can be taken by:
// 1. the core routine
// 2. the routine that syncs other peers
// can this be exploited by a malicious peer to slow Core down (by requesting ask sync lots of times)
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
	prevotes         *Set
	precommits       *Set
	power            *PowerInfo // power for all messages
	sync.RWMutex
}

// we need a reference to proposal also for proposing a proposal with vr!=0 if needed
func NewRoundMessages() *RoundMessages {
	return &RoundMessages{
		prevotes:         NewSet(),
		precommits:       NewSet(),
		power:            NewPowerInfo(),
		verifiedProposal: false,
	}
}

func (s *RoundMessages) SetProposal(proposal *Propose, verified bool) {
	s.Lock()
	defer s.Unlock()
	s.proposal = proposal
	s.verifiedProposal = verified
	s.power.Set(proposal.SignerIndex(), proposal.Power())
}

// total power for round (each signer counted only once, regardless of msg type)
func (s *RoundMessages) Power() *big.Int {
	return s.power.Pow()
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
	//TODO(lorenzo) can be moved in the set Add if computationally expensive
	// update round power cache
	for index, power := range prevote.Signers().Powers() {
		s.power.Set(index, power)
	}

}

func (s *RoundMessages) AllPrevotes() []Msg {
	return s.prevotes.Messages()
}

func (s *RoundMessages) AllPrecommits() []Msg {
	return s.precommits.Messages()
}

func (s *RoundMessages) AddPrecommit(precommit *Precommit) {
	s.precommits.Add(precommit)
	//TODO(lorenzo) can be moved in the set Add if computationally expensive
	// update round power cache
	for index, power := range precommit.Signers().Powers() {
		s.power.Set(index, power)
	}
}

// used to gossip quorum of prevotes
func (s *RoundMessages) PrevoteFor(hash common.Hash) *Prevote {
	prevotes := s.prevotes.VotesFor(hash)
	return AggregatePrevotes(prevotes) // we allow complex aggregate here
}

// used to create the quorum certificate when we managed to finalize a block and to gossip quorum of precommits
func (s *RoundMessages) PrecommitFor(hash common.Hash) *Precommit {
	precommits := s.precommits.VotesFor(hash)
	return AggregatePrecommits(precommits) // we allow complex aggregate here
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
