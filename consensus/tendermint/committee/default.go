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
	"github.com/clearmatics/autonity/consensus/tendermint/config"
	"github.com/clearmatics/autonity/core/types"
	"math"
	"reflect"
	"sort"
	"sync"
)

var ErrEmptyCommitteeSet = errors.New("committee set can't be empty")

type defaultSet struct {
	members      types.Committee
	policy       config.ProposerPolicy
	lastProposer common.Address
	selector     ProposalSelector

	mu       sync.RWMutex                    // members doesn't need to be protected as it is read-only
	proposer map[int64]types.CommitteeMember // cached computed values
}

func NewSet(members types.Committee, policy config.ProposerPolicy, lastProposer common.Address) (*defaultSet, error) {

	if len(members) == 0 {
		return nil, ErrEmptyCommitteeSet
	}

	commitee := &defaultSet{}
	commitee.policy = policy
	commitee.members = members
	commitee.proposer = make(map[int64]types.CommitteeMember)
	// sort validator
	sort.Sort(commitee.members)

	switch policy {
	case config.Sticky:
		commitee.selector = stickyProposer
	case config.RoundRobin:
		commitee.selector = roundRobinProposer
	default:
		commitee.selector = roundRobinProposer
	}

	commitee.lastProposer = lastProposer
	commitee.proposer[0] = commitee.selector(commitee, lastProposer, 0)
	return commitee, nil
}

func copyMembers(members types.Committee) types.Committee {
	membersCopy := make(types.Committee, len(members))
	for i, val := range members {
		membersCopy[i] = val
	}
	return membersCopy
}

func (set *defaultSet) Size() int {
	return len(set.members)
}

func (set *defaultSet) Committee() types.Committee {
	return copyMembers(set.members)
}

func (set *defaultSet) GetByIndex(i int) (types.CommitteeMember, error) {
	if i < 0 || i >= len(set.members) {
		return types.CommitteeMember{}, consensus.ErrCommitteeMemberNotFound
	}
	return set.members[i], nil
}

func (set *defaultSet) GetByAddress(addr common.Address) (int, types.CommitteeMember, error) {
	for i, member := range set.members {
		if addr == member.Address {
			return i, member, nil
		}
	}
	return -1, types.CommitteeMember{}, consensus.ErrCommitteeMemberNotFound
}

func (set *defaultSet) GetProposer(round int64) types.CommitteeMember {
	set.mu.Lock()
	defer set.mu.Unlock()
	v, ok := set.proposer[round]
	if !ok {
		v = set.selector(set, set.lastProposer, round)
		set.proposer[round] = v
	}

	return v
}

func (set *defaultSet) IsProposer(round int64, address common.Address) bool {
	_, val, err := set.GetByAddress(address)
	if err != nil {
		return false
	}
	curProposer := set.GetProposer(round)
	return reflect.DeepEqual(curProposer, val)
}

func (set *defaultSet) Copy() Set {
	newSet, _ := NewSet(copyMembers(set.members), set.policy, set.lastProposer)
	return newSet
}

func (set *defaultSet) F() int { return int(math.Ceil(float64(set.Size())/3)) - 1 }

func (set *defaultSet) Quorum() int { return int(math.Ceil((2 * float64(set.Size())) / 3.)) }

func (set *defaultSet) Policy() config.ProposerPolicy { return set.policy }
