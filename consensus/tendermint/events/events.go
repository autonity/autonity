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
	"time"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
)

// NewCandidateBlockEvent is posted to propose a proposal
type NewCandidateBlockEvent struct {
	NewCandidateBlock types.Block
	CreatedAt         time.Time
}

// UnverifiedMessageEvent is posted from the peer handlers to the aggregator
type UnverifiedMessageEvent struct {
	Message      message.Msg
	ErrCh        chan<- error
	Sender       common.Address
	Posted       time.Time
	Disseminated bool // proposals can be disseminated early
}

// MessageEvent is posted from the aggregator to core and the fault detector
type MessageEvent struct {
	message      message.Msg
	errCh        chan<- error
	posted       time.Time
	sender       common.Address
	disseminated bool
}

func NewMessageEvent(message message.Msg, errCh chan<- error, sender common.Address, posted time.Time, disseminated bool) MessageEvent {
	return MessageEvent{
		message:      message,
		errCh:        errCh,
		sender:       sender,
		posted:       posted,
		disseminated: disseminated,
	}
}

func (m MessageEvent) Message() message.Msg {
	return m.message
}

func (m MessageEvent) Sender() common.Address {
	return m.sender
}

func (m MessageEvent) Posted() time.Time {
	return m.posted
}

func (m MessageEvent) ErrCh() chan<- error {
	return m.errCh
}

func (m MessageEvent) Disseminated() bool {
	return m.disseminated
}

func (m *MessageEvent) SetDisseminated(disseminated bool) {
	m.disseminated = disseminated
}

// old messages are posted only to the fault detector
type OldMessageEvent struct {
	message message.Msg
	errCh   chan<- error
	sender  common.Address
	posted  time.Time
}

func NewOldMessageEvent(message message.Msg, errCh chan<- error, sender common.Address, posted time.Time) OldMessageEvent {
	return OldMessageEvent{
		message: message,
		errCh:   errCh,
		sender:  sender,
		posted:  posted,
	}
}

func (o OldMessageEvent) Message() message.Msg {
	return o.message
}

func (o OldMessageEvent) Sender() common.Address {
	return o.sender
}

func (o OldMessageEvent) Posted() time.Time {
	return o.posted
}

func (o OldMessageEvent) ErrCh() chan<- error {
	return o.errCh
}

type MessageEventer interface {
	Message() message.Msg
	Sender() common.Address
	Posted() time.Time
	ErrCh() chan<- error
}

type Poster interface {
	Post(interface{}) error
}

// CommitEvent is posted when a proposal is committed
type CommitEvent struct{}

type AccountabilityEvent struct {
	Sender  common.Address
	Payload []byte
	ErrCh   chan<- error
}
