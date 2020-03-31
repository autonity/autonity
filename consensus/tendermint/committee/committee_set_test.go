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

var (
	testAddress  = "70524d664ffe731100208a0154e556f9bb679ae6"
	testAddress2 = "b37866a925bccd69cfa98d43b510f1d23d78a851"
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
		{size: 3, round: 0},
		{size: 3, round: 1},
		{size: 3, round: 2},
		{size: 3, round: 3},
		{size: 3, round: 10},
		{size: 10, round: 0},
		{size: 10, round: 1},
		{size: 10, round: 2},
		{size: 10, round: 8},
		{size: 10, round: 7},
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

func TestValidatorSet(t *testing.T) {
	//testNewValidatorSet(t)
	//testNormalValSet(t)
	//testEmptyValSet(t)
}

func testNewValidatorSet(t *testing.T) {
	var validators types.Committee
	const ValCnt = 100

	// Create 100 members with random addresses
	for i := 0; i < ValCnt; i++ {
		key, _ := crypto.GenerateKey()
		addr := crypto.PubkeyToAddress(key.PublicKey)
		val := types.CommitteeMember{Address: addr, VotingPower: new(big.Int).SetUint64(1)}
		validators = append(validators, val)
	}

	// Create Set
	valSet, err := NewSet(validators, validators[0].Address)
	if err != nil || valSet == nil {
		t.Error("the validator byte array cannot be parsed")
		t.FailNow()
	}

	if valSet.Size() != ValCnt {
		t.Errorf("validator set has %d elements instead of %d", valSet.Size(), ValCnt)
	}

	valsMap := make(map[string]struct{})
	for _, val := range validators {
		valsMap[val.String()] = struct{}{}
	}

	// Check members sorting: should be in ascending order
	for i := 0; i < ValCnt-1; i++ {
		val, err := valSet.GetByIndex(i)
		if err != nil {
			t.Error("unexpected error")
		}
		nextVal, err := valSet.GetByIndex(i + 1)
		if err != nil {
			t.Error("unexpected error")
		}
		if strings.Compare(val.String(), nextVal.String()) >= 0 {
			t.Error("validator set is not sorted in ascending order")
		}

		if _, ok := valsMap[val.String()]; !ok {
			t.Errorf("validator set has unexpected element %s. Original members %v, given %v",
				val.String(), validators, valSet.Committee())
		}
	}
}

func testNormalValSet(t *testing.T) {
	b1 := common.Hex2Bytes(testAddress)
	b2 := common.Hex2Bytes(testAddress2)
	addr1 := common.BytesToAddress(b1)
	addr2 := common.BytesToAddress(b2)
	val1 := types.CommitteeMember{Address: addr1, VotingPower: new(big.Int).SetUint64(1)}
	val2 := types.CommitteeMember{Address: addr2, VotingPower: new(big.Int).SetUint64(1)}

	committeeSet, err := NewSet(types.Committee{val1, val2}, val1.Address)
	if committeeSet == nil || err != nil {
		t.Errorf("the format of validator set is invalid")
		t.FailNow()
	}

	// check size
	if size := committeeSet.Size(); size != 2 {
		t.Errorf("the size of validator set is wrong: have %v, want 2", size)
	}
	// test get by index
	if val, err := committeeSet.GetByIndex(0); err != nil || !reflect.DeepEqual(val, val1) {
		t.Errorf("validator mismatch: have %v, want %v", val, val1)
	}
	// test get by invalid index
	if _, err := committeeSet.GetByIndex(2); err != consensus.ErrCommitteeMemberNotFound {
		t.Errorf("validator mismatch: have %s, want nil", err)
	}
	// test get by address
	if _, val, err := committeeSet.GetByAddress(addr2); err != nil || !reflect.DeepEqual(val, val2) {
		t.Errorf("validator mismatch: have %v, want %v", val, val2)
	}
	// test get by invalid address
	invalidAddr := common.HexToAddress("0x9535b2e7faaba5288511d89341d94a38063a349b")
	if _, _, err := committeeSet.GetByAddress(invalidAddr); err != consensus.ErrCommitteeMemberNotFound {
		t.Errorf("validator mismatch: have %s, want error", err)
	}
	// test get proposers
	if val := committeeSet.GetProposer(0); !reflect.DeepEqual(val, val2) {
		t.Errorf("proposers mismatch: have %v, want %v", val, val1)
	}
	// test calculate proposers
	lastProposer := addr1
	committeeSet, _ = NewSet(types.Committee{val1, val2}, lastProposer)
	if val := committeeSet.GetProposer(0); !reflect.DeepEqual(val, val2) {
		t.Errorf("proposers mismatch: have %v, want %v", val, val2)
	}
	if val := committeeSet.GetProposer(3); !reflect.DeepEqual(val, val1) {
		t.Errorf("proposers mismatch: have %v, want %v", val, val1)
	}
	// test empty last proposers
	lastProposer = common.Address{}
	committeeSet, _ = NewSet(types.Committee{val1, val2}, lastProposer)
	if val := committeeSet.GetProposer(3); !reflect.DeepEqual(val, val2) {
		t.Errorf("proposers mismatch: have %v, want %v", val, val2)
	}
}

func testEmptyValSet(t *testing.T) {
	valSet, err := NewSet(types.Committee{}, common.Address{})
	if valSet != nil || err != ErrEmptyCommitteeSet {
		t.Errorf("validator set should be nil and error returned")
	}
}
