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
	"math/big"
	"sync"

	"github.com/clearmatics/autonity/common"
)

// NewRoundState creates a new roundState instance with the given view and validatorSet
// we need to keep a reference of proposal in order to propose locked proposal when there is a lock and itself is the proposer
func NewRoundState(r *big.Int, h *big.Int) *roundState {
	return &roundState{
		round:      r,
		height:     h,
		step:       propose,
		proposal:   new(Proposal),
		Prevotes:   newMessageSet(),
		Precommits: newMessageSet(),
	}
}

// roundState stores the consensus step
type roundState struct {
	round  *big.Int
	height *big.Int
	step   Step

	proposal    *Proposal
	proposalMsg *Message
	Prevotes    messageSet
	Precommits  messageSet
	mu          sync.RWMutex
}

func (s *roundState) Update(r *big.Int, h *big.Int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.round = r
	s.height = h
	s.proposal = new(Proposal)
	s.proposalMsg = nil
	s.Prevotes = newMessageSet()
	s.Precommits = newMessageSet()
}

func (s *roundState) SetProposal(proposal *Proposal, msg *Message) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.proposalMsg = msg
	s.proposal = proposal
}

func (s *roundState) Proposal() *Proposal {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.proposal != nil {
		return s.proposal
	}

	return nil
}

func (s *roundState) SetRound(r *big.Int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.round = big.NewInt(r.Int64())
}

func (s *roundState) Round() *big.Int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.round
}

func (s *roundState) SetHeight(height *big.Int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.height = height
}

func (s *roundState) Height() *big.Int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.height
}

func (s *roundState) SetStep(step Step) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.step = step
}

func (s *roundState) Step() Step {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.step
}

func (s *roundState) GetCurrentProposalHash() common.Hash {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.proposal.ProposalBlock != nil {
		return s.proposal.ProposalBlock.Hash()
	}

	return common.Hash{}
}

func (s *roundState) GetMessages() []*Message {
	s.mu.RLock()
	defer s.mu.RUnlock()

	prevoteMsgs := s.Prevotes.GetMessages()
	precommitMsgs := s.Precommits.GetMessages()

	result := make([]*Message, 0, len(prevoteMsgs)+len(precommitMsgs)+1)
	if s.proposalMsg != nil {
		result = append(result, s.proposalMsg)
	}

	result = append(result, prevoteMsgs...)
	result = append(result, precommitMsgs...)
	return result
}
