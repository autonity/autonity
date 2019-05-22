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
	"io"
	"math/big"
	"sync"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint"
	"github.com/clearmatics/autonity/rlp"
)

// newRoundState creates a new roundState instance with the given view and validatorSet
// we need to keep a reference of proposal in order to propose locked proposal when there is a lock and itself is the proposer
// TODO: instad of using big.Int use int64 and only use big.Int when needed for rlp encoding and use big.NewInt() to initialize
func newRoundState(r *big.Int, h *big.Int, validatorSet tendermint.ValidatorSet, hasBadProposal func(hash common.Hash) bool) *roundState {
	return &roundState{
		round:          r,
		height:         h,
		Prevotes:       newMessageSet(validatorSet),
		Precommits:     newMessageSet(validatorSet),
		mu:             new(sync.RWMutex),
		hasBadProposal: hasBadProposal,
	}
}

// roundState stores the consensus step
type roundState struct {
	round    *big.Int
	height   *big.Int
	proposal *tendermint.Proposal
	// TODO: messageSet is probably not required, need to port the function from messageSet to here, will also need to introduce mutexes.
	// These should simply be map[common.Address]message
	Prevotes *messageSet
	// TODO: ensure to check the size of the committed seals as mentioned by Roberto in Correctness and Analysis of IBFT paper
	Precommits *messageSet

	mu             *sync.RWMutex
	hasBadProposal func(hash common.Hash) bool
}

func (s *roundState) GetPrevoteOrPrecommitSize() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := s.Prevotes.Size() + s.Precommits.Size()

	// find duplicate one
	for _, m := range s.Prevotes.Values() {
		if s.Precommits.Get(m.Address) != nil {
			result--
		}
	}
	return result
}

func (s *roundState) Subject() *tendermint.Subject {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.Proposal() == nil {
		return nil
	}

	return &tendermint.Subject{
		Round:  new(big.Int).Set(s.round),
		Height: new(big.Int).Set(s.height),
		Digest: s.proposal.ProposalBlock.Hash(),
	}
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

	s.round = new(big.Int).Set(r)
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

// The DecodeRLP method should read one value from the given
// Stream. It is not forbidden to read less or more, but it might
// be confusing.
func (s *roundState) DecodeRLP(stream *rlp.Stream) error {
	var rs = new(roundState)

	if err := stream.Decode(&rs); err != nil {
		return err
	}
	s.round, s.height, s.proposal, s.Prevotes, s.Precommits, s.mu = rs.round, rs.height, rs.proposal, rs.Prevotes, rs.Precommits, new(sync.RWMutex)

	return nil
}

// EncodeRLP should write the RLP encoding of its receiver to w.
// If the implementation is a pointer method, it may also be
// called for nil pointers.
//
// Implementations should generate valid RLP. The data written is
// not verified at the moment, but a future version might. It is
// recommended to write only a single value but writing multiple
// values or no value at all is also permitted.
func (s *roundState) EncodeRLP(w io.Writer) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return rlp.Encode(w, []interface{}{s.round, s.height, s.proposal, s.Prevotes, s.Precommits})
}
