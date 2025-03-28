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
	Message message.Msg
	ErrCh   chan<- error
	Sender  common.Address
	Posted  time.Time
}

// MessageEvent is posted from the aggregator to core and the fault detector
type MessageEvent struct {
	Message message.Msg
	ErrCh   chan<- error
	Posted  time.Time
}

// old messages are posted only to the fault detector
type OldMessageEvent struct {
	Message message.Msg
	ErrCh   chan<- error
}

type Poster interface {
	Post(interface{}) error
}

// CommitEvent is posted when a proposal is committed
type CommitEvent struct{}

type CoreEvent interface {
	Height() uint64
	Round() int64
	Code() uint8
	Value() common.Hash
}

type RoundChangeEvent struct {
	height uint64
	round  int64
}

func NewRoundChangeEvent(height uint64, round int64) RoundChangeEvent {
	return RoundChangeEvent{height: height, round: round}
}

func (r RoundChangeEvent) Height() uint64 {
	return r.height
}

func (r RoundChangeEvent) Round() int64 {
	return r.round
}

func (r RoundChangeEvent) Code() uint8 {
	panic("not implemented")
}

func (r RoundChangeEvent) Value() common.Hash {
	panic("not implemented")
}

// change in voting power
type PowerChangeEvent struct {
	height uint64
	round  int64
	code   uint8
	value  common.Hash
}

func NewPowerChangeEvent(code uint8, height uint64, round int64, value common.Hash) PowerChangeEvent {
	return PowerChangeEvent{code: code, height: height, round: round, value: value}
}

func (p PowerChangeEvent) Height() uint64 {
	return p.height
}

func (p PowerChangeEvent) Round() int64 {
	return p.round
}

func (p PowerChangeEvent) Code() uint8 {
	return p.code
}

func (p PowerChangeEvent) Value() common.Hash {
	return p.value
}

// change in future round voting power
type FuturePowerChangeEvent struct {
	height uint64
	round  int64
}

func NewFuturePowerChangeEvent(height uint64, round int64) FuturePowerChangeEvent {
	return FuturePowerChangeEvent{height: height, round: round}
}

func (f FuturePowerChangeEvent) Height() uint64 {
	return f.height
}

func (f FuturePowerChangeEvent) Round() int64 {
	return f.round
}

func (f FuturePowerChangeEvent) Code() uint8 {
	panic("not implemented")
}

func (f FuturePowerChangeEvent) Value() common.Hash {
	panic("not implemented")
}

type LostSyncEvent struct {
	Sender  common.Address
	Payload []byte
	ErrCh   chan<- error
}

type AccountabilityEvent struct {
	Sender  common.Address
	Payload []byte
	ErrCh   chan<- error
}
