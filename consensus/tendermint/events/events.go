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

package events

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
)

// NewCandidateBlockEvent is posted to propose a proposal
type NewCandidateBlockEvent struct {
	NewCandidateBlock types.Block
}

//TODO(lorenzo) can we save the errch into the message ? or anyways avoid always having to bring it along the message

// UnverifiedMessageEvent is posted from the peer handlers to the aggregator
type UnverifiedMessageEvent struct {
	Message message.Msg
	ErrCh   chan<- error //error channel
}

// MessageEvent is posted from the aggregator to core and the fault detector
type MessageEvent struct {
	Message message.Msg
	ErrCh   chan<- error //error channel
}

/*  //TODO(lorenzo) for later
// old messages are posted only to the fault detector
type OldMessageEvent struct {
	Message message.Msg
	ErrCh   chan<- error
}*/

type Poster interface {
	Post(interface{}) error
}

// CommitEvent is posted when a proposal is committed
type CommitEvent struct{}

type SyncEvent struct {
	Addr common.Address
}

type AccountabilityEvent struct {
	Sender  common.Address
	Payload []byte
	ErrCh   chan<- error
}
