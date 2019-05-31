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
	"bytes"
	"math/big"
	"testing"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint"
	"github.com/clearmatics/autonity/consensus/tendermint/validator"
	"github.com/clearmatics/autonity/crypto"
)

func TestHandlePrecommit(t *testing.T) {
	N := uint64(4)

	proposal := newTestProposalBlock()
	expectedSubject := &tendermint.Subject{
		View: &tendermint.View{
			Round:    big.NewInt(0),
			Sequence: proposal.Number(),
		},
		Digest: proposal.Hash(),
	}

	testCases := []struct {
		system      *testSystem
		expectedErr error
	}{
		{
			// normal case
			func() *testSystem {
				sys := newTestSystemWithBackend(N)

				for i, backend := range sys.backends {
					c := backend.engine.(*core)
					c.valSet = backend.peers
					c.current = newTestRoundState(
						&tendermint.View{
							Round:    big.NewInt(0),
							Sequence: big.NewInt(1),
						},
						c.valSet,
					)

					if i == 0 {
						// replica 0 is the proposer
						c.state = StatePrevoteDone
					}
				}
				return sys
			}(),
			nil,
		},
		{
			// future message
			func() *testSystem {
				sys := newTestSystemWithBackend(N)

				for i, backend := range sys.backends {
					c := backend.engine.(*core)
					c.valSet = backend.peers
					if i == 0 {
						// replica 0 is the proposer
						c.current = newTestRoundState(
							expectedSubject.View,
							c.valSet,
						)
						c.state = StateProposeDone
					} else {
						c.current = newTestRoundState(
							&tendermint.View{
								Round:    big.NewInt(2),
								Sequence: big.NewInt(3),
							},
							c.valSet,
						)
					}
				}
				return sys
			}(),
			errFutureMessage,
		},
		{
			// subject not match
			func() *testSystem {
				sys := newTestSystemWithBackend(N)

				for i, backend := range sys.backends {
					c := backend.engine.(*core)
					c.valSet = backend.peers
					if i == 0 {
						// replica 0 is the proposer
						c.current = newTestRoundState(
							expectedSubject.View,
							c.valSet,
						)
						c.state = StateProposeDone
					} else {
						c.current = newTestRoundState(
							&tendermint.View{
								Round:    big.NewInt(0),
								Sequence: big.NewInt(0),
							},
							c.valSet,
						)
					}
				}
				return sys
			}(),
			errOldMessage,
		},
		{
			// jump state
			func() *testSystem {
				sys := newTestSystemWithBackend(N)

				for i, backend := range sys.backends {
					c := backend.engine.(*core)
					c.valSet = backend.peers
					c.current = newTestRoundState(
						&tendermint.View{
							Round:    big.NewInt(0),
							Sequence: proposal.Number(),
						},
						c.valSet,
					)

					// only replica0 stays at StateProposeDone
					// other replicas are at StatePrevoteDone
					if i != 0 {
						c.state = StatePrevoteDone
					} else {
						c.state = StateProposeDone
					}
				}
				return sys
			}(),
			nil,
		},
		// TODO: double send message
	}

OUTER:
	for _, test := range testCases {
		test.system.Run(false)

		v0 := test.system.backends[0]
		r0 := v0.engine.(*core)

		for i, v := range test.system.backends {
			validator := r0.valSet.GetByIndex(uint64(i))
			m, _ := Encode(v.engine.(*core).current.Subject())
			if err := r0.handlePrecommit(&message{
				Code:          msgPrecommit,
				Msg:           m,
				Address:       validator.Address(),
				Signature:     []byte{},
				CommittedSeal: validator.Address().Bytes(), // small hack
			}, validator); err != nil {
				if err != test.expectedErr {
					t.Errorf("error mismatch: have %v, want %v", err, test.expectedErr)
				}
				if r0.current.IsHashLocked() {
					t.Errorf("block should not be locked")
				}
				continue OUTER
			}
		}

		// prepared is normal case
		if r0.state != StatePrecommitDone {
			// There are not enough commit messages in core
			if r0.state != StatePrevoteDone {
				t.Errorf("state mismatch: have %v, want %v", r0.state, StatePrevoteDone)
			}
			if r0.current.Precommits.Size() > 2*r0.valSet.F() {
				t.Errorf("the size of commit messages should be less than %v", 2*r0.valSet.F()+1)
			}
			if r0.current.IsHashLocked() {
				t.Errorf("block should not be locked")
			}
			continue
		}

		// core should have 2F+1 prepare messages
		if r0.current.Precommits.Size() <= 2*r0.valSet.F() {
			t.Errorf("the size of commit messages should be larger than 2F+1: size %v", r0.current.Precommits.Size())
		}

		// check signatures large than 2F+1
		signedCount := 0
		committedSeals := v0.GetPrecommittedMsg(0).committedSeals
		for _, validator := range r0.valSet.List() {
			for _, seal := range committedSeals {
				if bytes.Equal(validator.Address().Bytes(), seal[:common.AddressLength]) {
					signedCount++
					break
				}
			}
		}
		if signedCount <= 2*r0.valSet.F() {
			t.Errorf("the expected signed count should be larger than %v, but got %v", 2*r0.valSet.F(), signedCount)
		}
		if !r0.current.IsHashLocked() {
			t.Errorf("block should be locked")
		}
	}
}

// round is not checked for now
func TestVerifyPrecommit(t *testing.T) {
	// for log purpose
	privateKey, _ := crypto.GenerateKey()
	peer := validator.New(getPublicKeyAddress(privateKey), 1)
	valSet := validator.NewSet(tendermint.RoundRobin, peer)

	sys := newTestSystemWithBackend(1)

	testCases := []struct {
		expected   error
		commit     *tendermint.Subject
		roundState *roundState
	}{
		{
			// normal case
			expected: nil,
			commit: &tendermint.Subject{
				View:   &tendermint.View{Round: big.NewInt(0), Sequence: big.NewInt(0)},
				Digest: newTestProposalBlock().Hash(),
			},
			roundState: newTestRoundState(
				&tendermint.View{Round: big.NewInt(0), Sequence: big.NewInt(0)},
				valSet,
			),
		},
		{
			// old message
			expected: errInconsistentSubject,
			commit: &tendermint.Subject{
				View:   &tendermint.View{Round: big.NewInt(0), Sequence: big.NewInt(0)},
				Digest: newTestProposalBlock().Hash(),
			},
			roundState: newTestRoundState(
				&tendermint.View{Round: big.NewInt(1), Sequence: big.NewInt(1)},
				valSet,
			),
		},
		{
			// different digest
			expected: errInconsistentSubject,
			commit: &tendermint.Subject{
				View:   &tendermint.View{Round: big.NewInt(0), Sequence: big.NewInt(0)},
				Digest: common.BytesToHash([]byte("1234567890")),
			},
			roundState: newTestRoundState(
				&tendermint.View{Round: big.NewInt(1), Sequence: big.NewInt(1)},
				valSet,
			),
		},
		{
			// malicious package(lack of sequence)
			expected: errInconsistentSubject,
			commit: &tendermint.Subject{
				View:   &tendermint.View{Round: big.NewInt(0), Sequence: nil},
				Digest: newTestProposalBlock().Hash(),
			},
			roundState: newTestRoundState(
				&tendermint.View{Round: big.NewInt(1), Sequence: big.NewInt(1)},
				valSet,
			),
		},
		{
			// wrong prepare message with same sequence but different round
			expected: errInconsistentSubject,
			commit: &tendermint.Subject{
				View:   &tendermint.View{Round: big.NewInt(1), Sequence: big.NewInt(0)},
				Digest: newTestProposalBlock().Hash(),
			},
			roundState: newTestRoundState(
				&tendermint.View{Round: big.NewInt(0), Sequence: big.NewInt(0)},
				valSet,
			),
		},
		{
			// wrong prepare message with same round but different sequence
			expected: errInconsistentSubject,
			commit: &tendermint.Subject{
				View:   &tendermint.View{Round: big.NewInt(0), Sequence: big.NewInt(1)},
				Digest: newTestProposalBlock().Hash(),
			},
			roundState: newTestRoundState(
				&tendermint.View{Round: big.NewInt(0), Sequence: big.NewInt(0)},
				valSet,
			),
		},
	}
	for i, test := range testCases {
		c := sys.backends[0].engine.(*core)
		c.current = test.roundState

		if err := c.verifyPrecommit(test.commit, peer); err != nil {
			if err != test.expected {
				t.Errorf("result %d: error mismatch: have %v, want %v", i, err, test.expected)
			}
		}
	}
}
