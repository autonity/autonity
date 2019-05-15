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
	"github.com/clearmatics/autonity/event"
	"github.com/clearmatics/autonity/log"
	"gopkg.in/karalabe/cookiejar.v2/collections/prque"
)

func TestCheckRequestMsg(t *testing.T) {
	c := &core{
		state: StateAcceptRequest,
		current: newRoundState(&tendermint.View{
			Sequence: big.NewInt(1),
			Round:    big.NewInt(0),
		}, newTestValidatorSet(4), common.Hash{}, nil, nil, nil),
	}

	// invalid request
	err := c.checkRequestMsg(nil)
	if err != errInvalidMessage {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidMessage)
	}
	r := &tendermint.Request{
		ProposalBlock: nil,
	}
	err = c.checkRequestMsg(r)
	if err != errInvalidMessage {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidMessage)
	}

	// old request
	r = &tendermint.Request{
		ProposalBlock: makeBlock(0),
	}
	err = c.checkRequestMsg(r)
	if err != errOldMessage {
		t.Errorf("error mismatch: have %v, want %v", err, errOldMessage)
	}

	// future request
	r = &tendermint.Request{
		ProposalBlock: makeBlock(2),
	}
	err = c.checkRequestMsg(r)
	if err != errFutureMessage {
		t.Errorf("error mismatch: have %v, want %v", err, errFutureMessage)
	}

	// current request
	r = &tendermint.Request{
		ProposalBlock: makeBlock(1),
	}
	err = c.checkRequestMsg(r)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
}

func TestStoreRequestMsg(t *testing.T) {
	logger := log.New("backend", "test", "id", 0)
	backend := &testSystemBackend{
		events: event.NewTypeMuxSilent(logger),
	}
	c := &core{
		logger:  logger,
		backend: backend,
		state:   StateAcceptRequest,
		current: newRoundState(&tendermint.View{
			Sequence: big.NewInt(0),
			Round:    big.NewInt(0),
		}, newTestValidatorSet(4), common.Hash{}, nil, nil, nil),
		pendingRequests:   prque.New(),
		pendingRequestsMu: new(sync.Mutex),
	}
	requests := []tendermint.Request{
		{
			ProposalBlock: makeBlock(1),
		},
		{
			ProposalBlock: makeBlock(2),
		},
		{
			ProposalBlock: makeBlock(3),
		},
	}

	c.storeRequestMsg(&requests[1])
	c.storeRequestMsg(&requests[0])
	c.storeRequestMsg(&requests[2])
	if c.pendingRequests.Size() != len(requests) {
		t.Errorf("the size of pending requests mismatch: have %v, want %v", c.pendingRequests.Size(), len(requests))
	}

	c.current.sequence = big.NewInt(3)

	c.subscribeEvents()
	defer c.unsubscribeEvents()

	c.processPendingRequests()

	const timeoutDura = 2 * time.Second
	timeout := time.NewTimer(timeoutDura)
	select {
	case ev := <-c.events.Chan():
		e, ok := ev.Data.(tendermint.RequestEvent)
		if !ok {
			t.Errorf("unexpected event comes: %v", reflect.TypeOf(ev.Data))
		}
		if e.ProposalBlock.Number().Cmp(requests[2].ProposalBlock.Number()) != 0 {
			t.Errorf("the number of proposal mismatch: have %v, want %v", e.ProposalBlock.Number(), requests[2].ProposalBlock.Number())
		}
	case <-timeout.C:
		t.Error("unexpected timeout occurs")
	}
}
