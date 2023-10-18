package message

import (
	"github.com/autonity/autonity/common"
	"math/big"
	"sync"
)

type MessagesMap struct {
	internal map[int64]*RoundMessages
	mu       *sync.RWMutex
}

func NewMessagesMap() *MessagesMap {
	return &MessagesMap{
		internal: make(map[int64]*RoundMessages),
		mu:       new(sync.RWMutex),
	}
}

func (s *MessagesMap) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.internal = make(map[int64]*RoundMessages)
}

func (s *MessagesMap) GetOrCreate(round int64) *RoundMessages {
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

func (s *MessagesMap) Messages() []*Message {
	s.mu.RLock()
	defer s.mu.RUnlock()

	msgs := make([][]*Message, len(s.internal))
	var totalLen int
	i := 0
	for _, state := range s.internal {
		msgs[i] = state.GetMessages()
		totalLen += len(msgs[i])
		i++
	}

	result := make([]*Message, 0, totalLen)
	for _, ms := range msgs {
		result = append(result, ms...)
	}

	return result
}

func (s *MessagesMap) GetRounds() []int64 {
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
	ProposalDetails  *Proposal
	VerifiedProposal bool
	ProposalMsg      *Message
	Prevotes         MessageSet
	Precommits       MessageSet
	mu               sync.RWMutex
}

// NewRoundMessages creates a new messages instance with the given view and validatorSet
// we need to keep a reference of proposal in order to propose locked proposal when there is a lock and itself is the proposer
func NewRoundMessages() *RoundMessages {
	return &RoundMessages{
		ProposalDetails:  new(Proposal),
		Prevotes:         NewMessageSet(),
		Precommits:       NewMessageSet(),
		VerifiedProposal: false,
	}
}

func (s *RoundMessages) SetProposal(proposal *Proposal, msg *Message, verified bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ProposalMsg = msg
	s.VerifiedProposal = verified
	s.ProposalDetails = proposal
}

//func (s *RoundMessages) SetProposal(proposal *Proposal, verified bool) {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//	//s.ProposalMsg = msg
//	s.VerifiedProposal = verified
//	s.ProposalDetails = proposal
//}

func (s *RoundMessages) PrevotesPower(hash common.Hash) *big.Int {
	return s.Prevotes.VotePower(hash)
}
func (s *RoundMessages) PrevotesTotalPower() *big.Int {
	return s.Prevotes.TotalVotePower()
}
func (s *RoundMessages) PrecommitsPower(hash common.Hash) *big.Int {
	return s.Precommits.VotePower(hash)
}
func (s *RoundMessages) PrecommitsTotalPower() *big.Int {
	return s.Precommits.TotalVotePower()
}

func (s *RoundMessages) AddPrevote(hash common.Hash, msg Message) {
	s.Prevotes.AddVote(hash, msg)
}

func (s *RoundMessages) AddPrecommit(hash common.Hash, msg Message) {
	s.Precommits.AddVote(hash, msg)
}

func (s *RoundMessages) CommitedSeals(hash common.Hash) []Message {
	return s.Precommits.Values(hash)
}

func (s *RoundMessages) Proposal() *Proposal {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.ProposalDetails != nil {
		return s.ProposalDetails
	}

	return nil
}

func (s *RoundMessages) IsProposalVerified() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.VerifiedProposal
}

func (s *RoundMessages) GetProposalHash() common.Hash {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.ProposalDetails.ProposalBlock != nil {
		return s.ProposalDetails.ProposalBlock.Hash()
	}

	return common.Hash{}
}

// func (s *RoundMessages) GetMessages() (*Proposal, []*Message, []*Message) {
func (s *RoundMessages) GetMessages() []*Message {
	s.mu.RLock()
	defer s.mu.RUnlock()

	prevoteMsgs := s.Prevotes.GetMessages()
	precommitMsgs := s.Precommits.GetMessages()

	result := make([]*Message, 0, len(prevoteMsgs)+len(precommitMsgs)+1)
	if s.ProposalMsg != nil {
		result = append(result, s.ProposalMsg)
	}

	result = append(result, prevoteMsgs...)
	result = append(result, precommitMsgs...)
	return result
	//return s.ProposalDetails, prevoteMsgs, precommitMsgs
}
