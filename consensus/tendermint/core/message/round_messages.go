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

func (s *Map) Snapshot() []*RoundMsgView {
	s.RLock()         // nolint
	defer s.RUnlock() // nolint
	var views []*RoundMsgView
	for r, v := range s.internal {
		views = append(views, v.Snapshot(r))
	}
	return views
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
	prevotes         *Set
	precommits       *Set
	power            *AggregatedPower // power for all messages
	sync.RWMutex
}

// we need a reference to proposal also for proposing a proposal with vr!=0 if needed
func NewRoundMessages() *RoundMessages {
	return &RoundMessages{
		prevotes:         NewSet(),
		precommits:       NewSet(),
		power:            NewAggregatedPower(),
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
func (s *RoundMessages) Power() *AggregatedPower {
	s.RLock()
	defer s.RUnlock()
	return s.power.Copy()
}

func (s *RoundMessages) PrevotesPower(hash common.Hash) *big.Int {
	return s.prevotes.PowerFor(hash).Power()
}

func (s *RoundMessages) PrevotesAggregatedPower(hash common.Hash) *AggregatedPower {
	return s.prevotes.PowerFor(hash)
}

func (s *RoundMessages) PrevotesTotalPower() *big.Int {
	return s.prevotes.TotalPower().Power()
}

func (s *RoundMessages) PrevotesTotalAggregatedPower() *AggregatedPower {
	return s.prevotes.TotalPower()
}

func (s *RoundMessages) PrecommitsPower(hash common.Hash) *big.Int {
	return s.precommits.PowerFor(hash).Power()
}

func (s *RoundMessages) PrecommitsAggregatedPower(hash common.Hash) *AggregatedPower {
	return s.precommits.PowerFor(hash)
}

func (s *RoundMessages) PrecommitsTotalPower() *big.Int {
	return s.precommits.TotalPower().Power()
}

func (s *RoundMessages) PrecommitsTotalAggregatedPower() *AggregatedPower {
	return s.precommits.TotalPower()
}

func (s *RoundMessages) AddPrevote(prevote *Prevote) {
	s.Lock()
	defer s.Unlock()
	s.prevotes.Add(prevote)
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
	s.Lock()
	defer s.Unlock()
	s.precommits.Add(precommit)
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

func (s *RoundMessages) Snapshot(round int64) *RoundMsgView {
	s.RLock()         // nolint
	defer s.RUnlock() // nolint
	view := &RoundMsgView{}
	view.Round = uint64(round)

	if s.proposal != nil {
		view.Proposal = s.proposal.Value()
	}

	if s.prevotes != nil {
		values, signers := s.prevotes.Snapshot()
		view.Prevotes = values
		view.PrevotesSigners = signers
	}

	if s.precommits != nil {
		values, signers := s.precommits.Snapshot()
		view.Precommits = values
		view.PrecommitsSigners = signers
	}
	return view
}

type RoundMsgView struct {
	Round uint64

	// the 1st received proposal that the client prevoted, normally only one, could be multiple if proposer equivocated.
	// different proposals will be exchanged if one find there is an equivocated one.
	Proposal common.Hash

	// prevoted values in the round
	Prevotes []common.Hash
	// signers for each prevoted values listed in the Prevotes slice.
	PrevotesSigners []*big.Int

	// precommitted values in the round
	Precommits []common.Hash
	// signders for each precommitted values listed in the Precommit slice.
	PrecommitsSigners []*big.Int
}

// LostSyncMsg carries all the msgs' views, include future rounds of current consensus engine for tendermint state recovery.
type LostSyncMsg struct {
	Height      uint64
	RoundsViews []*RoundMsgView
}
