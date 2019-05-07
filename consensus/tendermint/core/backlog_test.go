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
	"sync"
	"testing"
	"time"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint"
	"github.com/clearmatics/autonity/consensus/tendermint/validator"
	"github.com/clearmatics/autonity/event"
	"github.com/clearmatics/autonity/log"
	"gopkg.in/karalabe/cookiejar.v2/collections/prque"
)

func TestCheckMessage(t *testing.T) {
	c := &core{
		state: StateAcceptRequest,
		current: newRoundState(&tendermint.View{
			Sequence: big.NewInt(1),
			Round:    big.NewInt(0),
		}, newTestValidatorSet(4), common.Hash{}, nil, nil, nil),
	}

	// invalid view format
	err := c.checkMessage(msgProposal, nil)
	if err != errInvalidMessage {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidMessage)
	}

	testStates := []State{StateAcceptRequest, StateProposeDone, StatePrevoteDone, StatePrecommiteDone}
	testCode := []uint64{msgProposal, msgPrevote, msgPrecommit, msgRoundChange}

	// future sequence
	v := &tendermint.View{
		Sequence: big.NewInt(2),
		Round:    big.NewInt(0),
	}
	for i := 0; i < len(testStates); i++ {
		c.state = testStates[i]
		for j := 0; j < len(testCode); j++ {
			err := c.checkMessage(testCode[j], v)
			if err != errFutureMessage {
				t.Errorf("error mismatch: have %v, want %v", err, errFutureMessage)
			}
		}
	}

	// future round
	v = &tendermint.View{
		Sequence: big.NewInt(1),
		Round:    big.NewInt(1),
	}
	for i := 0; i < len(testStates); i++ {
		c.state = testStates[i]
		for j := 0; j < len(testCode); j++ {
			err := c.checkMessage(testCode[j], v)
			if testCode[j] == msgRoundChange {
				if err != nil {
					t.Errorf("error mismatch: have %v, want nil", err)
				}
			} else if err != errFutureMessage {
				t.Errorf("error mismatch: have %v, want %v", err, errFutureMessage)
			}
		}
	}

	// current view but waiting for round change
	v = &tendermint.View{
		Sequence: big.NewInt(1),
		Round:    big.NewInt(0),
	}
	c.waitingForRoundChange = true
	for i := 0; i < len(testStates); i++ {
		c.state = testStates[i]
		for j := 0; j < len(testCode); j++ {
			err := c.checkMessage(testCode[j], v)
			if testCode[j] == msgRoundChange {
				if err != nil {
					t.Errorf("error mismatch: have %v, want nil", err)
				}
			} else if err != errFutureMessage {
				t.Errorf("error mismatch: have %v, want %v", err, errFutureMessage)
			}
		}
	}
	c.waitingForRoundChange = false

	v = c.currentView()
	// current view, state = StateAcceptRequest
	c.state = StateAcceptRequest
	for i := 0; i < len(testCode); i++ {
		err = c.checkMessage(testCode[i], v)
		if testCode[i] == msgRoundChange {
			if err != nil {
				t.Errorf("error mismatch: have %v, want nil", err)
			}
		} else if testCode[i] == msgProposal {
			if err != nil {
				t.Errorf("error mismatch: have %v, want nil", err)
			}
		} else {
			if err != errFutureMessage {
				t.Errorf("error mismatch: have %v, want %v", err, errFutureMessage)
			}
		}
	}

	// current view, state = StateProposeDone
	c.state = StateProposeDone
	for i := 0; i < len(testCode); i++ {
		err = c.checkMessage(testCode[i], v)
		if testCode[i] == msgRoundChange {
			if err != nil {
				t.Errorf("error mismatch: have %v, want nil", err)
			}
		} else if err != nil {
			t.Errorf("error mismatch: have %v, want nil", err)
		}
	}

	// current view, state = StatePrevoteDone
	c.state = StatePrevoteDone
	for i := 0; i < len(testCode); i++ {
		err = c.checkMessage(testCode[i], v)
		if testCode[i] == msgRoundChange {
			if err != nil {
				t.Errorf("error mismatch: have %v, want nil", err)
			}
		} else if err != nil {
			t.Errorf("error mismatch: have %v, want nil", err)
		}
	}

	// current view, state = StateCommitted
	c.state = StatePrecommiteDone
	for i := 0; i < len(testCode); i++ {
		err = c.checkMessage(testCode[i], v)
		if testCode[i] == msgRoundChange {
			if err != nil {
				t.Errorf("error mismatch: have %v, want nil", err)
			}
		} else if err != nil {
			t.Errorf("error mismatch: have %v, want nil", err)
		}
	}

}

func TestStoreBacklog(t *testing.T) {
	c := &core{
		logger:     log.New("backend", "test", "id", 0),
		backlogs:   make(map[tendermint.Validator]*prque.Prque),
		backlogsMu: new(sync.Mutex),
	}
	v := &tendermint.View{
		Round:    big.NewInt(10),
		Sequence: big.NewInt(10),
	}
	p := validator.New(common.BytesToAddress([]byte("12345667890")))
	// push proposal msg
	proposal := &tendermint.Proposal{
		View:          v,
		ProposalBlock: makeBlock(1),
	}
	proposalPayload, _ := Encode(proposal)
	m := &message{
		Code: msgProposal,
		Msg:  proposalPayload,
	}
	c.storeBacklog(m, p)
	msg := c.backlogs[p].PopItem()
	if !reflect.DeepEqual(msg, m) {
		t.Errorf("message mismatch: have %v, want %v", msg, m)
	}

	// push prepare msg
	subject := &tendermint.Subject{
		View:   v,
		Digest: common.BytesToHash([]byte("1234567890")),
	}
	subjectPayload, _ := Encode(subject)

	m = &message{
		Code: msgPrevote,
		Msg:  subjectPayload,
	}
	c.storeBacklog(m, p)
	msg = c.backlogs[p].PopItem()
	if !reflect.DeepEqual(msg, m) {
		t.Errorf("message mismatch: have %v, want %v", msg, m)
	}

	// push commit msg
	m = &message{
		Code: msgPrecommit,
		Msg:  subjectPayload,
	}
	c.storeBacklog(m, p)
	msg = c.backlogs[p].PopItem()
	if !reflect.DeepEqual(msg, m) {
		t.Errorf("message mismatch: have %v, want %v", msg, m)
	}

	// push roundChange msg
	m = &message{
		Code: msgRoundChange,
		Msg:  subjectPayload,
	}
	c.storeBacklog(m, p)
	msg = c.backlogs[p].PopItem()
	if !reflect.DeepEqual(msg, m) {
		t.Errorf("message mismatch: have %v, want %v", msg, m)
	}
}

func TestProcessFutureBacklog(t *testing.T) {
	backend := &testSystemBackend{
		events: new(event.TypeMux),
	}
	c := &core{
		logger:     log.New("backend", "test", "id", 0),
		backlogs:   make(map[tendermint.Validator]*prque.Prque),
		backlogsMu: new(sync.Mutex),
		backend:    backend,
		current: newRoundState(&tendermint.View{
			Sequence: big.NewInt(1),
			Round:    big.NewInt(0),
		}, newTestValidatorSet(4), common.Hash{}, nil, nil, nil),
		state: StateAcceptRequest,
	}
	c.subscribeEvents()
	defer c.unsubscribeEvents()

	v := &tendermint.View{
		Round:    big.NewInt(10),
		Sequence: big.NewInt(10),
	}
	p := validator.New(common.BytesToAddress([]byte("12345667890")))
	// push a future msg
	subject := &tendermint.Subject{
		View:   v,
		Digest: common.BytesToHash([]byte("1234567890")),
	}
	subjectPayload, _ := Encode(subject)
	m := &message{
		Code: msgPrecommit,
		Msg:  subjectPayload,
	}
	c.storeBacklog(m, p)
	c.processBacklog()

	const timeoutDura = 2 * time.Second
	timeout := time.NewTimer(timeoutDura)
	select {
	case e, ok := <-c.events.Chan():
		if !ok {
			return
		}
		t.Errorf("unexpected events comes: %v", e)
	case <-timeout.C:
		// success
	}
}

func TestProcessBacklog(t *testing.T) {
	v := &tendermint.View{
		Round:    big.NewInt(0),
		Sequence: big.NewInt(1),
	}
	proposal := &tendermint.Proposal{
		View:          v,
		ProposalBlock: makeBlock(1),
	}
	proposalPayload, _ := Encode(proposal)

	subject := &tendermint.Subject{
		View:   v,
		Digest: common.BytesToHash([]byte("1234567890")),
	}
	subjectPayload, _ := Encode(subject)

	msgs := []*message{
		{
			Code: msgProposal,
			Msg:  proposalPayload,
		},
		{
			Code: msgPrevote,
			Msg:  subjectPayload,
		},
		{
			Code: msgPrecommit,
			Msg:  subjectPayload,
		},
		{
			Code: msgRoundChange,
			Msg:  subjectPayload,
		},
	}
	for i := 0; i < len(msgs); i++ {
		testProcessBacklog(t, msgs[i])
	}
}

func testProcessBacklog(t *testing.T, msg *message) {
	vset := newTestValidatorSet(1)
	backend := &testSystemBackend{
		events: new(event.TypeMux),
		peers:  vset,
	}
	c := &core{
		logger:     log.New("backend", "test", "id", 0),
		backlogs:   make(map[tendermint.Validator]*prque.Prque),
		backlogsMu: new(sync.Mutex),
		backend:    backend,
		state:      State(msg.Code),
		current: newRoundState(&tendermint.View{
			Sequence: big.NewInt(1),
			Round:    big.NewInt(0),
		}, newTestValidatorSet(4), common.Hash{}, nil, nil, nil),
	}
	c.subscribeEvents()
	defer c.unsubscribeEvents()

	c.storeBacklog(msg, vset.GetByIndex(0))
	c.processBacklog()

	const timeoutDura = 2 * time.Second
	timeout := time.NewTimer(timeoutDura)
	select {
	case ev := <-c.events.Chan():
		e, ok := ev.Data.(backlogEvent)
		if !ok {
			t.Errorf("unexpected event comes: %v", reflect.TypeOf(ev.Data))
		}
		if e.msg.Code != msg.Code {
			t.Errorf("message code mismatch: have %v, want %v", e.msg.Code, msg.Code)
		}
		// success
	case <-timeout.C:
		t.Error("unexpected timeout occurs")
	}
}
