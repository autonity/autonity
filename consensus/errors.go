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

package consensus

import "errors"

var (
	// ErrUnknownAncestor is returned when validating a block requires an ancestor
	// that is unknown.
	ErrUnknownAncestor = errors.New("unknown ancestor")

	// ErrPrunedAncestor is returned when validating a block requires an ancestor
	// that is known, but the state of which is not available.
	ErrPrunedAncestor = errors.New("pruned ancestor")

	// ErrFutureTimestampBlock is returned when a block's timestamp is in the future according
	// to the current node.
	ErrFutureTimestampBlock = errors.New("block in the future")

	// ErrInvalidNumber is returned if a block's number doesn't equal its parent's
	// plus one.
	ErrInvalidNumber = errors.New("invalid block number")

	// ErrInconsistentCommitteeSet is returned if the committee set is inconsistent
	ErrInconsistentCommitteeSet = errors.New("inconsistent committee set")

	// ErrInconsistentEpochInfo is returned if the epoch info is inconsistent
	ErrInconsistentEpochInfo = errors.New("inconsistent epoch info")

	// ErrCommitteeMemberNotFound is returned if the committee member is missing from
	// the committee set.
	ErrCommitteeMemberNotFound = errors.New("committee member not found")

	// ErrInvalidEpochInfo is returned if the epoch info is wrong
	ErrInvalidEpochInfo = errors.New("invalid epoch info")

	// ErrOutOfEpochRange is return if the inserting height is out of epoch range.
	ErrOutOfEpochRange = errors.New("the inserting height is out of epoch range")

	// ErrInvalidEpochBoundary is return if the bi-direction link of epoch header is broken.
	ErrInvalidEpochBoundary = errors.New("invalid epoch boundary")
)
