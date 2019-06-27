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
	"github.com/clearmatics/autonity/consensus/tendermint"
)

// newRoundState creates a new roundState instance with the given view and validatorSet
// we need to keep a reference of proposal in order to propose locked proposal when there is a lock and itself is the proposer
func newRoundState(r *big.Int, h *big.Int, hasBadProposal func(hash common.Hash) bool) *roundState {
	return &roundState{
		round:          r,
		height:         h,
		proposal:       new(tendermint.Proposal),
		Prevotes:       newMessageSet(),
		Precommits:     newMessageSet(),
		hasBadProposal: hasBadProposal,
	}
}

// roundState stores the consensus step
type roundState struct {
	round      *big.Int
	height     *big.Int
	proposal   *tendermint.Proposal
	Prevotes   messageSet
	Precommits messageSet

	mu             *sync.RWMutex
	hasBadProposal func(hash common.Hash) bool
}

func (s *roundState) SetProposal(proposal *tendermint.Proposal) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.proposal = proposal
}

func (s *roundState) Proposal() *tendermint.Proposal {
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

func (s *roundState) GetCurrentProposalHash() common.Hash {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.proposal.ProposalBlock != nil {
		return s.proposal.ProposalBlock.Hash()
	}

	return common.Hash{}
}
