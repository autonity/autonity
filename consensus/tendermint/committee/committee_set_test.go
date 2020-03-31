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
	"fmt"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/core/types"
	"math/big"
	"math/rand"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/clearmatics/autonity/crypto"
)

func TestNewSet(t *testing.T) {
	var committeeSetSizes = []int64{1, 2, 10, 100}
	var assertSet = func(t *testing.T, n int64) {
		t.Helper()

		committeeMembers := createTestCommitteeMembers(t, n)
		// Ensure last block proposer is chosen at random to test next proposer is chosen via round-robin
		lastBlockProposer := committeeMembers[rand.Intn(int(n))].Address

		// create copy since slice are pass by references
		// need to ensure a different copy of the committeMemebers is passed otherwise the sorting will affect the
		// committeeMembers and would not give any meaningful tests
		copyCommitteeMembers := copyMembers(committeeMembers)

		// next proposer is chosen after sorting
		sort.Sort(committeeMembers)
		// test the next proposer is chosen through round-robin

		roundRobinOffset := getMemberIndex(committeeMembers, lastBlockProposer)
		if len(committeeMembers) > 1 {
			roundRobinOffset += 1
		}
		proposers := map[int64]types.CommitteeMember{0: committeeMembers[nextProposerIndex(roundRobinOffset, 0, int64(len(committeeMembers)))]}

		set, err := NewSet(copyCommitteeMembers, lastBlockProposer)

		assertNilError(t, err)

		if lastBlockProposer != set.lastBlockProposer {
			t.Fatalf("lastBlockProposer not set properly, expected: %v and got: %v", lastBlockProposer, set.lastBlockProposer)
		}

		if roundRobinOffset != set.roundRobinOffset {
			t.Fatalf("roundRobinOffset not set properly, expected: %v and got: %v", roundRobinOffset, set.roundRobinOffset)
		}

		// This will also check sorting
		if !reflect.DeepEqual(committeeMembers, set.members) {
			t.Fatalf("committee memebers are not set properly, expected: %v and got: %v", committeeMembers, set.members)
		}

		if !reflect.DeepEqual(proposers, set.proposers) {
			t.Fatalf("initial round proposer not set properly, expected: %v and got: %v ", proposers, set.proposers)
		}
	}

	for _, size := range committeeSetSizes {
		t.Run(fmt.Sprintf("committee set of %v member/s", size), func(t *testing.T) {
			assertSet(t, size)
		})
	}

	t.Run("cannot create set with lastBlockProposer not in members", func(t *testing.T) {
		committeeMembers := createTestCommitteeMembers(t, 2)
		lastBlockProposer := committeeMembers[1]
		committeeMembers = committeeMembers[:1]
		_, err := NewSet(committeeMembers, lastBlockProposer.Address)
		assertError(t, ErrLastBlockProposerNotInCommitteeSet, err)

	})

	t.Run("cannot create empty set with members as nil", func(t *testing.T) {
		_, err := NewSet(nil, common.Address{})
		assertError(t, ErrEmptyCommitteeSet, err)
	})

	t.Run("cannot create empty set with members as types.Committee{}", func(t *testing.T) {
		_, err := NewSet(types.Committee{}, common.Address{})
		assertError(t, ErrEmptyCommitteeSet, err)
	})
}

func TestSet_Size(t *testing.T) {
	var committeeSetSizes = []int64{1, 2, 10, 100}
	var assertSetSize = func(t *testing.T, n int64) {
		t.Helper()

		committeeMembers := createTestCommitteeMembers(t, n)
		// only testing size so don't care about sorting
		set, err := NewSet(committeeMembers, committeeMembers[0].Address)
		assertNilError(t, err)

		setSize := set.Size()
		if int64(setSize) != n {
			t.Fatalf("expected committee set size: %v and got: %v", n, setSize)
		}
	}

	for _, size := range committeeSetSizes {
		t.Run(fmt.Sprintf("committee size of %v member/s", size), func(t *testing.T) {
			assertSetSize(t, size)
		})
	}

}

func TestSet_Committee(t *testing.T) {
	var committeeSetSizes = []int64{1, 2, 10, 100}
	var assertSetCommittee = func(t *testing.T, n int64) {
		t.Helper()

		committeeMembers := createTestCommitteeMembers(t, n)
		set, err := NewSet(copyMembers(committeeMembers), committeeMembers[0].Address)
		sort.Sort(committeeMembers)
		assertNilError(t, err)

		gotCommittee := set.Committee()

		if !reflect.DeepEqual(committeeMembers, gotCommittee) {
			t.Fatalf("expected committee: %v and got: %v", committeeMembers, gotCommittee)
		}
	}

	for _, size := range committeeSetSizes {
		t.Run(fmt.Sprintf("get committee of %v member/s", size), func(t *testing.T) {
			assertSetCommittee(t, size)
		})
	}
}

func TestSet_GetByIndex(t *testing.T) {
	committeeMembers := createTestCommitteeMembers(t, 4)
	sort.Sort(committeeMembers)
	set, err := NewSet(copyMembers(committeeMembers), committeeMembers[0].Address)
	assertNilError(t, err)

	t.Run("can get member by index", func(t *testing.T) {
		expectedMember := committeeMembers[1]
		gotMember, err := set.GetByIndex(1)
		assertNilError(t, err)

		if !reflect.DeepEqual(expectedMember, gotMember) {
			t.Fatalf("expected member: %v and got %v", expectedMember, gotMember)
		}
	})

	t.Run("error on accessing member index not in committee", func(t *testing.T) {
		_, err := set.GetByIndex(6)
		assertError(t, consensus.ErrCommitteeMemberNotFound, err)
	})
}

func TestSet_GetByAddress(t *testing.T) {
	committeeMembers := createTestCommitteeMembers(t, 4)
	sort.Sort(committeeMembers)
	set, err := NewSet(copyMembers(committeeMembers), committeeMembers[0].Address)
	assertNilError(t, err)

	t.Run("can get member by Address", func(t *testing.T) {
		expectedMember := committeeMembers[1]
		index, gotMember, err := set.GetByAddress(expectedMember.Address)
		assertNilError(t, err)

		if index != 1 {
			t.Fatalf("incorrect index of member expected: %v and got %v", 1, index)
		}

		if !reflect.DeepEqual(expectedMember, gotMember) {
			t.Fatalf("expected member: %v and got %v", expectedMember, gotMember)
		}
	})

	t.Run("error on accessing member address not in committee", func(t *testing.T) {
		_, _, err := set.GetByAddress(common.HexToAddress("testaddress"))
		assertError(t, consensus.ErrCommitteeMemberNotFound, err)
	})
}

func TestSet_GetProposer(t *testing.T) {
	testCases := []struct {
		size  int64
		round int64
	}{
		{size: 3, round: 0}, {size: 3, round: 1}, {size: 3, round: 2}, {size: 3, round: 3}, {size: 3, round: 10},
		{size: 10, round: 0}, {size: 10, round: 1}, {size: 10, round: 2}, {size: 10, round: 8}, {size: 10, round: 7},
		{size: 10, round: 10},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("validator set size %v and round %v", testCase.size, testCase.round), func(t *testing.T) {
			committeeMembers := createTestCommitteeMembers(t, testCase.size)
			sort.Sort(committeeMembers)
			lastBlockProposer := committeeMembers[rand.Intn(int(testCase.size))].Address
			roundRobinOffset := getMemberIndex(committeeMembers, lastBlockProposer)
			if len(committeeMembers) > 1 {
				roundRobinOffset += 1
			}
			expectedProposerIndex := (roundRobinOffset + testCase.round) % testCase.size
			expectedProposer := committeeMembers[expectedProposerIndex]

			set, err := NewSet(copyMembers(committeeMembers), lastBlockProposer)
			assertNilError(t, err)

			gotProposer := set.GetProposer(int64(testCase.round))

			if expectedProposer != gotProposer {
				t.Fatalf("expected proposer: %v and got: %v", expectedProposer, gotProposer)
			}
		})
	}
}

func TestSet_IsProposer(t *testing.T) {
	rounds := []int64{0, 1, 2, 3, 4, 5, 6, 7, 8}
	committeeMembers := createTestCommitteeMembers(t, 4)
	sort.Sort(committeeMembers)
	lastBlockProposerIndex := 2
	lastBlockProposer := committeeMembers[lastBlockProposerIndex].Address
	roundRobinOffset := int64(lastBlockProposerIndex + 1)

	set, err := NewSet(copyMembers(committeeMembers), lastBlockProposer)
	assertNilError(t, err)

	for _, r := range rounds {
		t.Run(fmt.Sprintf("correct proposer for round %v", r), func(t *testing.T) {
			testAddr := committeeMembers[(roundRobinOffset+r)%4].Address
			isProposer := set.IsProposer(r, testAddr)
			if !isProposer {
				t.Fatalf("expected IsProposer(0, %v) to return true", testAddr)
			}
		})
	}
	t.Run("false if addres is in committe set but is not the proposer for round", func(t *testing.T) {
		// committeeMembers[0].Address cannot be the proposer of round 0
		isProposer := set.IsProposer(0, lastBlockProposer)
		if isProposer {
			t.Fatalf("did not expect IsProposer(0, %v) to return true", lastBlockProposer)
		}
	})
	t.Run("false if address is not in committe set", func(t *testing.T) {
		testAddr := common.HexToAddress("testaddress")
		isProposer := set.IsProposer(0, common.HexToAddress("testaddress"))
		if isProposer {
			t.Fatalf("did not expect IsProposer(0, %v) to return true", testAddr)
		}
	})
}

func TestSet_Copy(t *testing.T) {
	committeeMembers := createTestCommitteeMembers(t, 4)
	set, err := NewSet(copyMembers(committeeMembers), committeeMembers[0].Address)
	assertNilError(t, err)

	copiedSet := set.Copy()
	if !reflect.DeepEqual(set, copiedSet) {
		t.Fatalf("failed to correctly copy set, expected: %v and got: %v", set, copiedSet)
	}
}

func TestSet_QandF(t *testing.T) {
	testCases := []struct {
		N int
		Q int
		F int
	}{
		{N: 1, Q: 1, F: 0}, {N: 2, Q: 2, F: 0}, {N: 3, Q: 2, F: 0}, {N: 4, Q: 3, F: 1}, {N: 5, Q: 4, F: 1},
		{N: 6, Q: 4, F: 1}, {N: 7, Q: 5, F: 2}, {N: 8, Q: 6, F: 2}, {N: 9, Q: 6, F: 2}, {N: 10, Q: 7, F: 3},
		{N: 11, Q: 8, F: 3}, {N: 12, Q: 8, F: 3}, {N: 13, Q: 9, F: 4}, {N: 14, Q: 10, F: 4}, {N: 15, Q: 10, F: 4},
		{N: 16, Q: 11, F: 5}, {N: 17, Q: 12, F: 5}, {N: 18, Q: 12, F: 5}, {N: 19, Q: 13, F: 6}, {N: 20, Q: 14, F: 6},
		{N: 21, Q: 14, F: 6}, {N: 22, Q: 15, F: 7}, {N: 23, Q: 16, F: 7}, {N: 24, Q: 16, F: 7}, {N: 25, Q: 17, F: 8},
		{N: 26, Q: 18, F: 8}, {N: 27, Q: 18, F: 8}, {N: 28, Q: 19, F: 9}, {N: 29, Q: 20, F: 9}, {N: 30, Q: 20, F: 9},
		{N: 31, Q: 21, F: 10}, {N: 32, Q: 22, F: 10}, {N: 33, Q: 22, F: 10}, {N: 34, Q: 23, F: 11}, {N: 35, Q: 24, F: 11},
		{N: 36, Q: 24, F: 11}, {N: 37, Q: 25, F: 12}, {N: 38, Q: 26, F: 12}, {N: 39, Q: 26, F: 12}, {N: 40, Q: 27, F: 13},
		{N: 41, Q: 28, F: 13}, {N: 42, Q: 28, F: 13}, {N: 43, Q: 29, F: 14}, {N: 44, Q: 30, F: 14}, {N: 45, Q: 30, F: 14},
		{N: 46, Q: 31, F: 15}, {N: 47, Q: 32, F: 15}, {N: 48, Q: 32, F: 15}, {N: 49, Q: 33, F: 16}, {N: 50, Q: 34, F: 16},
		{N: 51, Q: 34, F: 16}, {N: 52, Q: 35, F: 17}, {N: 53, Q: 36, F: 17}, {N: 54, Q: 36, F: 17}, {N: 55, Q: 37, F: 18},
		{N: 56, Q: 38, F: 18}, {N: 57, Q: 38, F: 18}, {N: 58, Q: 39, F: 19}, {N: 59, Q: 40, F: 19}, {N: 60, Q: 40, F: 19},
		{N: 61, Q: 41, F: 20}, {N: 62, Q: 42, F: 20}, {N: 63, Q: 42, F: 20}, {N: 64, Q: 43, F: 21}, {N: 65, Q: 44, F: 21},
		{N: 66, Q: 44, F: 21}, {N: 67, Q: 45, F: 22}, {N: 68, Q: 46, F: 22}, {N: 69, Q: 46, F: 22}, {N: 70, Q: 47, F: 23},
		{N: 71, Q: 48, F: 23}, {N: 72, Q: 48, F: 23}, {N: 73, Q: 49, F: 24}, {N: 74, Q: 50, F: 24}, {N: 75, Q: 50, F: 24},
		{N: 76, Q: 51, F: 25}, {N: 77, Q: 52, F: 25}, {N: 78, Q: 52, F: 25}, {N: 79, Q: 53, F: 26}, {N: 80, Q: 54, F: 26},
		{N: 81, Q: 54, F: 26}, {N: 82, Q: 55, F: 27}, {N: 83, Q: 56, F: 27}, {N: 84, Q: 56, F: 27}, {N: 85, Q: 57, F: 28},
		{N: 86, Q: 58, F: 28}, {N: 87, Q: 58, F: 28}, {N: 88, Q: 59, F: 29}, {N: 89, Q: 60, F: 29}, {N: 90, Q: 60, F: 29},
		{N: 91, Q: 61, F: 30}, {N: 92, Q: 62, F: 30}, {N: 93, Q: 62, F: 30}, {N: 94, Q: 63, F: 31}, {N: 95, Q: 64, F: 31},
		{N: 96, Q: 64, F: 31}, {N: 97, Q: 65, F: 32}, {N: 98, Q: 66, F: 32}, {N: 99, Q: 66, F: 32}, {N: 100, Q: 67, F: 33},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("N: %v, Q: %v, F: %v", testCase.N, testCase.Q, testCase.F), func(t *testing.T) {
			committeeMembers := createTestCommitteeMembers(t, int64(testCase.N))
			set, err := NewSet(committeeMembers, committeeMembers[0].Address)
			assertNilError(t, err)

			gotF := set.F()
			gotQ := set.Quorum()

			if testCase.F != gotF {
				t.Errorf("expected F: %v and got: %v", testCase.F, gotF)
			}

			if testCase.Q != gotQ {
				t.Errorf("expected Q: %v and got: %v", testCase.Q, gotQ)
			}
		})
	}

}

func assertNilError(t *testing.T, got error) {
	t.Helper()
	if got != nil {
		t.Fatalf("unexpected error %v", got)
	}
}

func assertError(t *testing.T, want, got error) {
	t.Helper()
	if want != got {
		t.Fatalf("expected an err: %v and got: %v", want, got)
	}
}

func createTestCommitteeMembers(t *testing.T, n int64) types.Committee {
	t.Helper()
	var committee types.Committee
	for i := 0; i < int(n); i++ {
		key, err := crypto.GenerateKey()

		if err != nil {
			t.Fatal(err)
		}
		member := types.CommitteeMember{
			Address:     crypto.PubkeyToAddress(key.PublicKey),
			VotingPower: new(big.Int).SetUint64(1),
		}
		committee = append(committee, member)
	}

	if n > 0 {
		// swap 1st and last element if 1st element is less then last to ensure committee is not sorted
		firstIndex, lastIndex := 0, len(committee)-1
		comp := strings.Compare(committee[firstIndex].String(), committee[lastIndex].String())
		if comp < 0 {
			committee[firstIndex], committee[lastIndex] = committee[lastIndex], committee[firstIndex]
		}
	}

	return committee
}
