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
// lockedHash and proposal are for round change when lock exists,
// we need to keep a reference of proposal in order to propose locked proposal when there is a lock and itself is the proposer
func newRoundState(view *tendermint.View, validatorSet tendermint.ValidatorSet, lockedHash common.Hash, proposal *tendermint.Proposal, pendingRequest *tendermint.Request, hasBadProposal func(hash common.Hash) bool) *roundState {
	return &roundState{
		round:          view.Round,
		sequence:       view.Sequence,
		proposal:       proposal,
		Prevotes:       newMessageSet(validatorSet),
		Precommits:     newMessageSet(validatorSet),
		lockedHash:     lockedHash,
		mu:             new(sync.RWMutex),
		pendingRequest: pendingRequest,
		hasBadProposal: hasBadProposal,
	}
}

// roundState stores the consensus state
type roundState struct {
	round          *big.Int
	sequence       *big.Int
	proposal       *tendermint.Proposal
	Prevotes       *messageSet
	Precommits     *messageSet
	lockedHash     common.Hash
	pendingRequest *tendermint.Request

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
		View: &tendermint.View{
			Round:    new(big.Int).Set(s.round),
			Sequence: new(big.Int).Set(s.sequence),
		},
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

func (s *roundState) SetSequence(seq *big.Int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sequence = seq
}

func (s *roundState) Sequence() *big.Int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.sequence
}

func (s *roundState) LockHash() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.proposal != nil {
		s.lockedHash = s.proposal.ProposalBlock.Hash()
	}
}

func (s *roundState) UnlockHash() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.lockedHash = common.Hash{}
}

func (s *roundState) IsHashLocked() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.lockedHash == (common.Hash{}) {
		return false
	}
	return !s.hasBadProposal(s.GetLockedHash())
}

func (s *roundState) GetLockedHash() common.Hash {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.lockedHash
}

// The DecodeRLP method should read one value from the given
// Stream. It is not forbidden to read less or more, but it might
// be confusing.
func (s *roundState) DecodeRLP(stream *rlp.Stream) error {
	var ss struct {
		Round          *big.Int
		Sequence       *big.Int
		proposal       *tendermint.Proposal
		Prevotes       *messageSet
		Precommits     *messageSet
		lockedHash     common.Hash
		pendingRequest *tendermint.Request
	}

	if err := stream.Decode(&ss); err != nil {
		return err
	}
	s.round = ss.Round
	s.sequence = ss.Sequence
	s.proposal = ss.proposal
	s.Prevotes = ss.Prevotes
	s.Precommits = ss.Precommits
	s.lockedHash = ss.lockedHash
	s.pendingRequest = ss.pendingRequest
	s.mu = new(sync.RWMutex)

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

	return rlp.Encode(w, []interface{}{
		s.round,
		s.sequence,
		s.proposal,
		s.Prevotes,
		s.Precommits,
		s.lockedHash,
		s.pendingRequest,
	})
}
