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
	"errors"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/core/types"
	"math"
	"reflect"
	"sort"
	"sync"
)

var ErrEmptyCommitteeSet = errors.New("committee set can't be empty")

type Set struct {
	members           types.Committee
	lastBlockProposer common.Address
	roundRobinOffset  int64
	proposers         map[int64]types.CommitteeMember // cached computed values

	mu sync.RWMutex // members doesn't need to be protected as it is read-only
}

func NewSet(members types.Committee, lastBlockProposer common.Address) (*Set, error) {

	// Ensure non empty set
	if len(members) == 0 {
		return nil, ErrEmptyCommitteeSet
	}

	// Create new Set
	committee := &Set{
		members:           members,
		lastBlockProposer: lastBlockProposer,
		proposers:         make(map[int64]types.CommitteeMember),
	}

	// sort validator
	sort.Sort(committee.members)

	// calculate offset for round robin selection of next proposer
	committee.roundRobinOffset = getMemberIndex(committee.members, lastBlockProposer)
	if len(members) > 1 {
		committee.roundRobinOffset++
	}
	committee.proposers[0] = committee.getNextProposer(0)

	return committee, nil
}

func (set *Set) Size() int {
	set.mu.RLock()
	defer set.mu.RUnlock()
	return len(set.members)
}

func (set *Set) Committee() types.Committee {
	set.mu.RLock()
	defer set.mu.RUnlock()
	return copyMembers(set.members)
}

func (set *Set) GetByIndex(i int) (types.CommitteeMember, error) {
	set.mu.RLock()
	defer set.mu.RUnlock()
	if i < 0 || i >= len(set.members) {
		return types.CommitteeMember{}, consensus.ErrCommitteeMemberNotFound
	}
	return set.members[i], nil
}

func (set *Set) GetByAddress(addr common.Address) (int, types.CommitteeMember, error) {
	set.mu.RLock()
	defer set.mu.RUnlock()
	for i, member := range set.members {
		if addr == member.Address {
			return i, member, nil
		}
	}
	return -1, types.CommitteeMember{}, consensus.ErrCommitteeMemberNotFound
}

func (set *Set) GetProposer(round int64) types.CommitteeMember {
	set.mu.Lock()
	defer set.mu.Unlock()

	v, ok := set.proposers[round]
	if !ok {
		v = set.getNextProposer(round)
		set.proposers[round] = v
	}

	return v
}

func (set *Set) IsProposer(round int64, address common.Address) bool {
	_, val, err := set.GetByAddress(address)
	if err != nil {
		return false
	}
	curProposer := set.GetProposer(round)
	return reflect.DeepEqual(curProposer, val)
}

func (set *Set) Copy() *Set {
	newSet, _ := NewSet(copyMembers(set.members), set.lastBlockProposer)
	return newSet
}

func (set *Set) F() int { return int(math.Ceil(float64(set.Size())/3.)) - 1 }

func (set *Set) Quorum() int { return int(math.Ceil((2 * float64(set.Size())) / 3.)) }

func (set *Set) getNextProposer(round int64) types.CommitteeMember {
	return set.members[nextProposerIndex(set.roundRobinOffset, round, int64(len(set.members)))]
}

func nextProposerIndex(offset, round, committeeSize int64) int64 {
	// Round-Robin
	return (offset + round) % committeeSize
}

func getMemberIndex(members types.Committee, memberAddr common.Address) int64 {
	var index = -1
	for i, member := range members {
		if memberAddr == member.Address {
			index = i
		}
	}
	return int64(index)
}

func copyMembers(members types.Committee) types.Committee {
	membersCopy := make(types.Committee, len(members))
	for i, val := range members {
		membersCopy[i] = val
	}
	return membersCopy
}
