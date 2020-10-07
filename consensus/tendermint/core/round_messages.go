// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"sync"

	"github.com/clearmatics/autonity/common"
)

// "messagesMap" for a lack of better name.
// could have been replaced by a Sync map

type messagesMap struct {
	internal map[int64]*roundMessages
	mu       *sync.RWMutex
}

func newMessagesMap() messagesMap {
	return messagesMap{
		internal: make(map[int64]*roundMessages),
		mu:       new(sync.RWMutex),
	}
}

func (s *messagesMap) reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.internal = make(map[int64]*roundMessages)
}

func (s *messagesMap) getOrCreate(round int64) *roundMessages {
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

func (s *messagesMap) GetMessages() []*Message {
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

func (s *messagesMap) getRounds() []int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rounds := make([]int64, 0, len(s.internal))
	for r := range s.internal {
		rounds = append(rounds, r)
	}

	return rounds
}

func (s *messagesMap) getVoteState(round int64) (common.Hash, []VoteState, []VoteState) {
	p := common.Hash{}

	if s.getOrCreate(round).Proposal() != nil && s.getOrCreate(round).Proposal().ProposalBlock != nil {
		p = s.getOrCreate(round).Proposal().ProposalBlock.Hash()
	}

	pvv := s.getOrCreate(round).GetPrevoteValues()
	pcv := s.getOrCreate(round).GetPrecommitValues()
	prevoteState := make([]VoteState, 0, len(pvv))
	precommitState := make([]VoteState, 0, len(pcv))

	for _, v := range pvv {
		var s = VoteState{
			Value:            v,
			ProposalVerified: s.getOrCreate(round).isProposalVerified(),
			VotePower:        s.getOrCreate(round).PrevotesPower(v),
		}
		prevoteState = append(prevoteState, s)
	}

	for _, v := range pcv {
		var s = VoteState{
			Value:            v,
			ProposalVerified: s.getOrCreate(round).isProposalVerified(),
			VotePower:        s.getOrCreate(round).PrecommitsPower(v),
		}
		precommitState = append(precommitState, s)
	}

	return p, prevoteState, precommitState
}

// roundMessages stores all message received for a specific round.
type roundMessages struct {
	proposal         *Proposal
	verifiedProposal bool
	proposalMsg      *Message
	prevotes         messageSet
	precommits       messageSet
	mu               sync.RWMutex
}

// NewRoundMessages creates a new messages instance with the given view and validatorSet
// we need to keep a reference of proposal in order to propose locked proposal when there is a lock and itself is the proposer
func NewRoundMessages() *roundMessages {
	return &roundMessages{
		proposal:         new(Proposal),
		prevotes:         newMessageSet(),
		precommits:       newMessageSet(),
		verifiedProposal: false,
	}
}

func (s *roundMessages) SetProposal(proposal *Proposal, msg *Message, verified bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.proposalMsg = msg
	s.verifiedProposal = verified
	s.proposal = proposal
}

func (s *roundMessages) GetPrevoteValues() []common.Hash {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.prevotes.BlockHashes()
}

func (s *roundMessages) GetPrecommitValues() []common.Hash {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.precommits.BlockHashes()
}

func (s *roundMessages) PrevotesPower(hash common.Hash) uint64 {
	return s.prevotes.VotePower(hash)
}
func (s *roundMessages) PrevotesTotalPower() uint64 {
	return s.prevotes.TotalVotePower()
}
func (s *roundMessages) PrecommitsPower(hash common.Hash) uint64 {
	return s.precommits.VotePower(hash)
}
func (s *roundMessages) PrecommitsTotalPower() uint64 {
	return s.precommits.TotalVotePower()
}

func (s *roundMessages) AddPrevote(hash common.Hash, msg Message) {
	s.prevotes.AddVote(hash, msg)
}

func (s *roundMessages) AddPrecommit(hash common.Hash, msg Message) {
	s.precommits.AddVote(hash, msg)
}

func (s *roundMessages) CommitedSeals(hash common.Hash) []Message {
	return s.precommits.Values(hash)
}

func (s *roundMessages) Proposal() *Proposal {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.proposal != nil {
		return s.proposal
	}

	return nil
}

func (s *roundMessages) isProposalVerified() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.verifiedProposal
}

func (s *roundMessages) GetProposalHash() common.Hash {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.proposal.ProposalBlock != nil {
		return s.proposal.ProposalBlock.Hash()
	}

	return common.Hash{}
}

func (s *roundMessages) GetMessages() []*Message {
	s.mu.RLock()
	defer s.mu.RUnlock()

	prevoteMsgs := s.prevotes.GetMessages()
	precommitMsgs := s.precommits.GetMessages()

	result := make([]*Message, 0, len(prevoteMsgs)+len(precommitMsgs)+1)
	if s.proposalMsg != nil {
		result = append(result, s.proposalMsg)
	}

	result = append(result, prevoteMsgs...)
	result = append(result, precommitMsgs...)
	return result
}
