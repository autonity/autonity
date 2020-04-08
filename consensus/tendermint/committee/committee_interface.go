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

package committee

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
)

type Set interface {
	// Return the validator size
	Size() int
	// Return the underlying types.Committee
	Committee() types.Committee
	// Get validator by index
	GetByIndex(i int) (types.CommitteeMember, error)
	// Get validator by given address
	GetByAddress(addr common.Address) (int, types.CommitteeMember, error)
	// Get the round proposer
	GetProposer(round int64) types.CommitteeMember
	// Check whether the validator with given address is the round proposer
	IsProposer(round int64, address common.Address) bool
	// Copy validator set
	Copy() Set
	// Get the maximum number of faulty nodes
	F() uint64
	// Get the optimal quorum size
	Quorum() uint64
}

// ----------------------------------------------------------------------------

type ProposalSelector func(Set, common.Address, int64) types.CommitteeMember
