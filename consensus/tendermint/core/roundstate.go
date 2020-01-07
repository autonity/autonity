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
)

// NewRoundState creates a new roundState instance with the given view and validatorSet
// we need to keep a reference of proposal in order to propose locked proposal when there is a lock and itself is the proposer
func NewRoundState(r *big.Int, h *big.Int) *roundState {
	return &roundState{
		round:  r,
		height: h,
		step:   propose,
	}
}

// roundState stores the consensus step
type roundState struct {
	round  *big.Int
	height *big.Int
	step   Step

	// TODO: potentially add getters and setters for allRoundMessages
	allRoundMessages map[int64]roundMessageSet
	mu               sync.RWMutex
}

type roundMessageSet struct {
	proposal    *Proposal
	proposalMsg *Message
	prevotes    messageSet
	precommits  messageSet
}

func newRoundMessageSet() roundMessageSet {
	return roundMessageSet{
		prevotes:   newMessageSet(),
		precommits: newMessageSet(),
	}
}

func (s *roundState) Update(r *big.Int, h *big.Int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.round = r
	s.height = h
}

func (s *roundState) SetProposal(round int64, proposal *Proposal, msg *Message) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rms := s.allRoundMessages[round]
	rms.proposalMsg = msg
	rms.proposal = proposal
}

func (s *roundState) Proposal(round int64) *Proposal {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rms := s.allRoundMessages[round]
	if rms.proposal != nil {
		return rms.proposal
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

func (s *roundState) CurrentState() (*big.Int, *big.Int, uint64) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.height, s.round, uint64(s.step)
}

func (s *roundState) GetMessages(round int64) []*Message {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rms := s.allRoundMessages[round]

	prevoteMsgs := rms.prevotes.GetMessages()
	precommitMsgs := rms.precommits.GetMessages()

	result := make([]*Message, 0, len(prevoteMsgs)+len(precommitMsgs)+1)
	if rms.proposalMsg != nil {
		result = append(result, rms.proposalMsg)
	}
	result = append(result, prevoteMsgs...)
	result = append(result, precommitMsgs...)

	return result
}

func (s *roundState) GetAllRoundMessages() []*Message {
	var messages []*Message

	for roundNumber, _ := range s.allRoundMessages {
		messages = append(messages, s.GetMessages(roundNumber)...)
	}

	return messages
}
