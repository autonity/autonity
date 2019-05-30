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
	"reflect"
	"testing"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint"
	"github.com/clearmatics/autonity/consensus/tendermint/validator"
	"github.com/clearmatics/autonity/crypto"
)

func TestHandlePrevote(t *testing.T) {
	N := uint64(4)
	F := uint64(1)

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
						c.state = StateProposeDone
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
								Sequence: big.NewInt(1)},
							c.valSet,
						)
					}
				}
				return sys
			}(),
			errInconsistentSubject,
		},
		{
			// less than 2F+1
			func() *testSystem {
				sys := newTestSystemWithBackend(N)

				// save less than 2*F+1 replica
				sys.backends = sys.backends[2*int(F)+1:]

				for i, backend := range sys.backends {
					c := backend.engine.(*core)
					c.valSet = backend.peers
					c.current = newTestRoundState(
						expectedSubject.View,
						c.valSet,
					)

					if i == 0 {
						// replica 0 is the proposer
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
			if err := r0.handlePrevote(&message{
				Code:    msgPrevote,
				Msg:     m,
				Address: validator.Address(),
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
		if r0.state != StatePrevoteDone {
			// There are not enough PREPARE messages in core
			if r0.state != StateProposeDone {
				t.Errorf("state mismatch: have %v, want %v", r0.state, StateProposeDone)
			}
			if r0.current.Prevotes.Size() > 2*r0.valSet.F() {
				t.Errorf("the size of PREPARE messages should be less than %v", 2*r0.valSet.F()+1)
			}
			if r0.current.IsHashLocked() {
				t.Errorf("block should not be locked")
			}

			continue
		}

		// core should have 2F+1 PREPARE messages
		if r0.current.Prevotes.Size() <= 2*r0.valSet.F() {
			t.Errorf("the size of PREPARE messages should be larger than 2F+1: size %v", r0.current.Precommits.Size())
		}

		// a message will be delivered to backend if 2F+1
		if int64(len(v0.sentMsgs)) != 1 {
			t.Errorf("the Send() should be called once: times %v", len(test.system.backends[0].sentMsgs))
		}

		// verify COMMIT messages
		decodedMsg := new(message)
		err := decodedMsg.FromPayload(v0.sentMsgs[0], nil)
		if err != nil {
			t.Errorf("error mismatch: have %v, want nil", err)
		}

		if decodedMsg.Code != msgPrecommit {
			t.Errorf("message code mismatch: have %v, want %v", decodedMsg.Code, msgPrecommit)
		}
		var m *tendermint.Subject
		err = decodedMsg.Decode(&m)
		if err != nil {
			t.Errorf("error mismatch: have %v, want nil", err)
		}
		if !reflect.DeepEqual(m, expectedSubject) {
			t.Errorf("subject mismatch: have %v, want %v", m, expectedSubject)
		}
		if !r0.current.IsHashLocked() {
			t.Errorf("block should be locked")
		}
	}
}

// round is not checked for now
func TestVerifyPrevote(t *testing.T) {
	// for log purpose
	privateKey, _ := crypto.GenerateKey()
	peer := validator.New(getPublicKeyAddress(privateKey))
	valSet := validator.NewSet(tendermint.RoundRobin, []common.Address{peer.Address()}, tendermint.RoundRobin)

	sys := newTestSystemWithBackend(1)

	testCases := []struct {
		expected error

		prepare    *tendermint.Subject
		roundState *roundState
	}{
		{
			// normal case
			expected: nil,
			prepare: &tendermint.Subject{
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
			prepare: &tendermint.Subject{
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
			prepare: &tendermint.Subject{
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
			prepare: &tendermint.Subject{
				View:   &tendermint.View{Round: big.NewInt(0), Sequence: nil},
				Digest: newTestProposalBlock().Hash(),
			},
			roundState: newTestRoundState(
				&tendermint.View{Round: big.NewInt(1), Sequence: big.NewInt(1)},
				valSet,
			),
		},
		{
			// wrong PREPARE message with same sequence but different round
			expected: errInconsistentSubject,
			prepare: &tendermint.Subject{
				View:   &tendermint.View{Round: big.NewInt(1), Sequence: big.NewInt(0)},
				Digest: newTestProposalBlock().Hash(),
			},
			roundState: newTestRoundState(
				&tendermint.View{Round: big.NewInt(0), Sequence: big.NewInt(0)},
				valSet,
			),
		},
		{
			// wrong PREPARE message with same round but different sequence
			expected: errInconsistentSubject,
			prepare: &tendermint.Subject{
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

		if err := c.verifyPrevote(test.prepare, peer); err != nil {
			if err != test.expected {
				t.Errorf("result %d: error mismatch: have %v, want %v", i, err, test.expected)
			}
		}
	}
}
