package message

import (
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
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

func (s *Map) DumpMsgView() []*RoundMsgView {
	s.RLock()         // nolint
	defer s.RUnlock() // nolint
	views := make([]*RoundMsgView, 0, 32)
	for r, state := range s.internal {
		views = append(views, state.DumpMsgView(r))
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

// returns whether the new vote brought some power contribution or the vote was redundant
func (s *RoundMessages) AddPrevote(prevote *Prevote) bool {
	s.Lock()
	defer s.Unlock()
	voteContributed := s.prevotes.Add(prevote)
	// update round power cache
	for index, power := range prevote.Signers().Powers() {
		s.power.Set(index, power)
	}
	return voteContributed
}

func (s *RoundMessages) AllPrevotes() []Msg {
	return s.prevotes.Messages()
}

func (s *RoundMessages) AllPrecommits() []Msg {
	return s.precommits.Messages()
}

// returns whether the new vote brought some power contribution or the vote was redundant
func (s *RoundMessages) AddPrecommit(precommit *Precommit) bool {
	s.Lock()
	defer s.Unlock()
	voteContributed := s.precommits.Add(precommit)
	// update round power cache
	for index, power := range precommit.Signers().Powers() {
		s.power.Set(index, power)
	}
	return voteContributed
}

// used to gossip quorum of prevotes
func (s *RoundMessages) PrevoteFor(hash common.Hash) []*Prevote {
	prevotes := s.prevotes.VotesFor(hash)
	return AggregatePrevotes(prevotes) // we allow complex aggregate here
}

// used to create the quorum certificate when we managed to finalize a block and to gossip quorum of precommits
func (s *RoundMessages) PrecommitFor(hash common.Hash) []*Precommit {
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

func (s *RoundMessages) DumpMsgView(round int64) *RoundMsgView {
	s.RLock()         // nolint
	defer s.RUnlock() // nolint
	view := &RoundMsgView{}
	view.Round = uint64(round)

	if s.proposal != nil {
		view.HaveProposal = true
	}

	if s.prevotes != nil {
		values, signers := s.prevotes.DumpMsgView()
		view.Prevotes = values
		view.PrevotesSigners = signers
	}

	if s.precommits != nil {
		values, signers := s.precommits.DumpMsgView()
		view.Precommits = values
		view.PrecommitsSigners = signers
	}
	return view
}

type RoundMsgView struct {
	Round uint64

	HaveProposal bool

	// prevoted values in the round
	Prevotes []common.Hash
	// signers for each prevoted values listed in the Prevotes slice.
	PrevotesSigners []*big.Int

	// precommitted values in the round
	Precommits []common.Hash
	// signders for each precommitted values listed in the Precommit slice.
	PrecommitsSigners []*big.Int
}

var errInvalidLostSyncMsg = errors.New("invalid ask sync message")

// AskSyncMsg carries all the msgs' views, include future rounds of current consensus engine for tendermint state recovery.
type AskSyncMsg struct {
	Height        uint64
	KnownMessages []*RoundMsgView
	// following fields are ignored when rlp/json encoding/decoding.
	// They are built locally when validating the message
	validated        bool                                `rlp:"-"`
	prevoteSigners   map[uint64]map[common.Hash]*big.Int `rlp:"-"` // maps point to the signers bitmap for that specific value
	precommitSigners map[uint64]map[common.Hash]*big.Int `rlp:"-"` // maps point to the signers bitmap for that specific value
	hasProposal      map[uint64]struct{}                 `rlp:"-"` // marks round where remote node does have a proposal
}

func (m *AskSyncMsg) Validate() error {
	// cannot have more than `MaxRound` + 1 distinct rounds
	if len(m.KnownMessages) > constants.MaxRound+1 {
		return errInvalidLostSyncMsg
	}

	rounds := make(map[uint64]struct{})
	prevoteSigners := make(map[uint64]map[common.Hash]*big.Int)
	precommitSigners := make(map[uint64]map[common.Hash]*big.Int)
	hasProposal := make(map[uint64]struct{})
	for _, v := range m.KnownMessages {
		// view cannot be nil
		if v == nil {
			return errInvalidLostSyncMsg
		}
		// round number cannot be > `MaxRound`
		if v.Round > constants.MaxRound {
			return errInvalidLostSyncMsg
		}
		// rounds should not repeat
		if _, ok := rounds[v.Round]; ok {
			return errInvalidLostSyncMsg
		}
		rounds[v.Round] = struct{}{}

		// if the remote peer does have a proposal for this round, mark it
		if v.HaveProposal {
			hasProposal[v.Round] = struct{}{}
		}

		// sanity check prevotes of this round
		currentPrevoteSigners, err := validateVotes(v.Prevotes, v.PrevotesSigners)
		if err != nil {
			return fmt.Errorf("error while sanity checking prevotes: %w", err)
		}

		// sanity check precommits of this round
		currentPrecommitSigners, err := validateVotes(v.Precommits, v.PrecommitsSigners)
		if err != nil {
			return fmt.Errorf("error while sanity checking precommits: %w", err)
		}
		prevoteSigners[v.Round] = currentPrevoteSigners
		precommitSigners[v.Round] = currentPrecommitSigners
	}

	// populate local fields if valid
	m.validated = true
	m.prevoteSigners = prevoteSigners
	m.precommitSigners = precommitSigners
	m.hasProposal = hasProposal
	return nil
}

func validateVotes(values []common.Hash, signers []*big.Int) (map[common.Hash]*big.Int, error) {
	// number of values and number of signers should be coherent
	if len(values) != len(signers) {
		return nil, errInvalidLostSyncMsg
	}

	voteSigners := make(map[common.Hash]*big.Int)
	for i, value := range values {
		// values of same round votes shouldn't repeat
		if _, ok := voteSigners[value]; ok {
			return nil, errInvalidLostSyncMsg
		}
		// signers shouldn't be nil or empty
		if signers[i] == nil || signers[i] == common.Big0 {
			return nil, errInvalidLostSyncMsg
		}
		voteSigners[value] = signers[i]
	}
	return voteSigners, nil
}

func (m *AskSyncMsg) HasProposal() map[uint64]struct{} {
	if !m.validated {
		panic("AskSyncMsg.HasProposal() called before validated")
	}
	return m.hasProposal
}

func (m *AskSyncMsg) Prevotes() map[uint64]map[common.Hash]*big.Int {
	if !m.validated {
		panic("AskSyncMsg.Prevotes() called before validated")
	}
	return m.prevoteSigners
}

func (m *AskSyncMsg) Precommits() map[uint64]map[common.Hash]*big.Int {
	if !m.validated {
		panic("AskSyncMsg.Precommits() called before validated")
	}
	return m.precommitSigners
}
