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

import "errors"

var (
	// errInconsistentSubject is returned when received subject is different from
	// current subject.
	errInconsistentSubject = errors.New("inconsistent subjects")
	// errNotFromProposer is returned when received message is supposed to be from
	// proposer.
	errNotFromProposer = errors.New("message does not come from proposer")
	// errIgnored is returned when a message was ignored.
	errIgnored = errors.New("message is ignored")
	// errFutureMessage is returned when current view is earlier than the
	// view of the received message.
	errFutureMessage = errors.New("future message")
	// errOldMessage is returned when the received message's view is earlier
	// than current view.
	errOldMessage = errors.New("old message")
	// errInvalidMessage is returned when the message is malformed.
	errInvalidMessage = errors.New("invalid message")
	// errFailedDecodeProposal is returned when the PRE-PREPARE message is malformed.
	errFailedDecodeProposal = errors.New("failed to decode PRE-PREPARE")
	// errFailedDecodePrevote is returned when the PREPARE message is malformed.
	errFailedDecodePrevote = errors.New("failed to decode PREPARE")
	// errFailedDecodePrecommit is returned when the COMMIT message is malformed.
	errFailedDecodePrecommit = errors.New("failed to decode COMMIT")
	// errFailedDecodeMessageSet is returned when the message set is malformed.
	errFailedDecodeMessageSet = errors.New("failed to decode message set")
)
