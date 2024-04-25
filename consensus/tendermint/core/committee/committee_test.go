package committee

import (
	"bytes"
	"fmt"
	"github.com/autonity/autonity/crypto/blst"
	"golang.org/x/exp/slices"
	"math/big"
	"math/rand"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
)

const maxSize = 100

func TestNewRoundRobinSet(t *testing.T) {
	var committeeSetSizes = []int64{1, 2, 10, 100}
	var assertSet = func(t *testing.T, n int64) {
		c := createTestCommitteeMembers(t, n, genRandUint64(int(n), maxSize))
		c.Sort()
		// Ensure last block proposer is chosen at random to test next proposer is chosen via round-robin
		lastBlockProposer := c.Members[rand.Intn(int(n))].Address
		// create copy since slice are pass by references
		// need to ensure a different copy of the committeMemebers is passed otherwise the sorting will affect the
		// committee and would not give any meaningful tests
		copyCommitteeMembers := c.Copy()
		// next proposer is chosen after sorting
		// test the next proposer is chosen through round-robin
		roundRobinOffset := getMemberIndex(c, lastBlockProposer)
		if c.Len() > 1 {
			roundRobinOffset++
		}
		allProposers := map[int64]*types.CommitteeMember{0: c.Members[nextProposerIndex(roundRobinOffset, 0, int64(c.Len()))]}
		totalPower := c.TotalVotingPower()

		set, err := NewRoundRobinSet(copyCommitteeMembers, lastBlockProposer)
		assertNilError(t, err)
		if lastBlockProposer != set.lastBlockProposer {
			t.Fatalf("lastBlockProposer not set properly, expected: %v and got: %v", lastBlockProposer, set.lastBlockProposer)
		}
		if roundRobinOffset != set.roundRobinOffset {
			t.Fatalf("roundRobinOffset not set properly, expected: %v and got: %v", roundRobinOffset, set.roundRobinOffset)
		}

		// This will also check sorting
		if !reflect.DeepEqual(c.Members, set.committee.Members) {
			t.Fatalf("committee memebers are not set properly, expected: %v and got: %v", c, set.committee.Members)
		}

		if !reflect.DeepEqual(allProposers, set.allProposers) {
			t.Fatalf("initial round allProposers not set properly, expected: %v and got: %v ", allProposers, set.allProposers)
		}

		if totalPower.Cmp(set.committee.TotalVotingPower()) != 0 {
			t.Fatalf("totalPower not calculated properly, expected: %v and got: %v ", totalPower, set.committee.TotalVotingPower())
		}
	}

	for _, size := range committeeSetSizes {
		size := size
		t.Run(fmt.Sprintf("committee set of %v member/s", size), func(t *testing.T) {
			assertSet(t, size)
		})
	}
	t.Run("cannot create empty set with members as nil", func(t *testing.T) {
		_, err := NewRoundRobinSet(nil, common.Address{})
		assertError(t, ErrEmptyCommitteeSet, err)
	})

	t.Run("cannot create empty set with members as types.Committee{}", func(t *testing.T) {
		_, err := NewRoundRobinSet(&types.Committee{}, common.Address{})
		assertError(t, ErrEmptyCommitteeSet, err)
	})
}

// We need to ensure that the committee is sorted, so that block hashes are the same for all validators.
func TestCommitteeIsSorted(t *testing.T) {
	c := createTestCommitteeMembers(t, 10, 10)
	c.Sort()
	sorted := slices.IsSortedFunc(c.Members, func(a, b *types.CommitteeMember) bool {
		return b.VotingPower.Cmp(a.VotingPower) < 0
	})
	require.True(t, sorted)

	set, err := NewRoundRobinSet(c, c.Members[0].Address)
	require.NoError(t, err)
	assert.True(t, slices.IsSortedFunc(set.Committee().Members, func(a, b *types.CommitteeMember) bool {
		return b.VotingPower.Cmp(a.VotingPower) < 0
	}))
}

func TestSet_Committee(t *testing.T) {
	var committeeSetSizes = []int64{1, 2, 10, 100}
	var assertSetCommittee = func(t *testing.T, n int64) {
		c := createTestCommitteeMembers(t, n, genRandUint64(int(n), maxSize))
		c.Sort()
		set, err := NewRoundRobinSet(c.Copy(), c.Members[0].Address)
		assertNilError(t, err)

		gotCommittee := set.Committee()

		if !reflect.DeepEqual(c, gotCommittee) {
			t.Fatalf("expected committee: %v and got: %v", c, gotCommittee)
		}
	}

	for _, size := range committeeSetSizes {
		size := size
		t.Run(fmt.Sprintf("get committee of %v member/s", size), func(t *testing.T) {
			assertSetCommittee(t, size)
		})
	}
}

func TestSet_GetByIndex(t *testing.T) {
	c := createTestCommitteeMembers(t, 4, genRandUint64(4, maxSize))
	c.Sort()
	set, err := NewRoundRobinSet(c.Copy(), c.Members[0].Address)
	assertNilError(t, err)

	t.Run("can get member by index", func(t *testing.T) {
		expectedMember := c.Members[1]
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
	c := createTestCommitteeMembers(t, 4, genRandUint64(4, maxSize))
	c.Sort()
	set, err := NewRoundRobinSet(c.Copy(), c.Members[0].Address)
	assertNilError(t, err)

	t.Run("can get member by Address", func(t *testing.T) {
		expectedMember := c.Members[1]
		gotMember, err := set.GetByAddress(expectedMember.Address)
		assertNilError(t, err)

		if !reflect.DeepEqual(expectedMember, gotMember) {
			t.Fatalf("expected member: %v and got %v", expectedMember, gotMember)
		}
	})

	t.Run("error on accessing member address not in committee", func(t *testing.T) {
		_, err := set.GetByAddress(common.HexToAddress("testaddress"))
		assertError(t, consensus.ErrCommitteeMemberNotFound, err)
	})
}

// TestSet_GetProposer tests the round robin selection of proposers. It validates that as GetProposer is
// called with consecutive rounds, consecutive proposers are chosen in the sort order defined by
// types.Committee. The consequence of this is that proposers are selected fairly, with N-1 other proposers
// being selected between any two instances of the same proposer being selected twice. It also validates that
// the selection process starts from the committee member that follows lastBlockProposer in a sorted instance
// of types.Committee.
func TestSet_GetProposer(t *testing.T) {
	numOfPasses := 10
	setSizes := 100
	for size := 1; size <= setSizes; size++ {
		size := size
		t.Run(fmt.Sprintf("check round robin for validator set size of %v", size), func(t *testing.T) {
			c := createTestCommitteeMembers(t, int64(size), genRandUint64(size, maxSize))
			c.Sort()
			r := rand.Intn(size)
			lastBlockProposer := c.Members[r].Address
			expectedProposerAddrForRound0 := c.Members[(r+1)%size].Address

			set, err := NewRoundRobinSet(c.Copy(), lastBlockProposer)
			require.NoError(t, err)

			firstCommitteeMemberAddr := c.Members[0].Address
			var startRound, endRound int
			for i := 1; i <= numOfPasses; i++ {
				startRound = endRound
				endRound = size * i
				committeeFromCallingGetProposer := new(types.Committee)
				for j := startRound; j < endRound; j++ {
					committeeFromCallingGetProposer.Members = append(committeeFromCallingGetProposer.Members, set.GetProposer(int64(j)))
				}
				// Ensure the proposer for round % size = 0 is the following next member from the lastBlockProposer
				// in the sorted committee set.
				assert.Equal(t, expectedProposerAddrForRound0, committeeFromCallingGetProposer.Members[0].Address)

				// Determine where committeeFromCallingGetProposer and ordered c line up using
				// firstCommitteeMember.
				var startIndex int
				for k, m := range committeeFromCallingGetProposer.Members {
					if m.Address == firstCommitteeMemberAddr {
						startIndex = k
						break
					}
				}
				assert.Equal(t, c.Members, append(committeeFromCallingGetProposer.Members[startIndex:], committeeFromCallingGetProposer.Members[:startIndex]...))
			}
		})
	}
}

func TestSet_IsProposer(t *testing.T) {
	rounds := []int64{0, 1, 2, 3, 4, 5, 6, 7, 8}
	c := createTestCommitteeMembers(t, 4, genRandUint64(4, maxSize))
	c.Sort()
	lastBlockProposerIndex := 2
	lastBlockProposer := c.Members[lastBlockProposerIndex].Address
	roundRobinOffset := lastBlockProposerIndex + 1

	set, err := NewRoundRobinSet(c.Copy(), lastBlockProposer)
	assertNilError(t, err)

	for _, r := range rounds {
		r := r
		t.Run(fmt.Sprintf("correct proposer for round %v", r), func(t *testing.T) {
			testAddr := c.Members[(int64(roundRobinOffset)+r)%4].Address
			isProposer := set.GetProposer(r).Address == testAddr
			if !isProposer {
				t.Fatalf("expected IsProposer(0, %v) to return true", testAddr)
			}
		})
	}
	t.Run("false if address is in committee set but is not the proposer for round", func(t *testing.T) {
		// c[0].Address cannot be the proposer of round 0
		isProposer := set.GetProposer(0).Address == lastBlockProposer
		if isProposer {
			t.Fatalf("did not expect IsProposer(0, %v) to return true", lastBlockProposer)
		}
	})
	t.Run("false if address is not in committee set", func(t *testing.T) {
		testAddr := common.HexToAddress("testaddress")
		isProposer := set.GetProposer(0).Address == common.HexToAddress("testaddress")
		if isProposer {
			t.Fatalf("did not expect IsProposer(0, %v) to return true", testAddr)
		}
	})
}

func TestSet_QandF(t *testing.T) {
	testCases := []struct {
		TotalVP int64
		Q       uint64
		F       uint64
	}{
		{TotalVP: 1, Q: 1, F: 0}, {TotalVP: 2, Q: 2, F: 0}, {TotalVP: 3, Q: 2, F: 0}, {TotalVP: 4, Q: 3, F: 1}, {TotalVP: 5, Q: 4, F: 1},
		{TotalVP: 6, Q: 4, F: 1}, {TotalVP: 7, Q: 5, F: 2}, {TotalVP: 8, Q: 6, F: 2}, {TotalVP: 9, Q: 6, F: 2}, {TotalVP: 10, Q: 7, F: 3},
		{TotalVP: 11, Q: 8, F: 3}, {TotalVP: 12, Q: 8, F: 3}, {TotalVP: 13, Q: 9, F: 4}, {TotalVP: 14, Q: 10, F: 4}, {TotalVP: 15, Q: 10, F: 4},
		{TotalVP: 16, Q: 11, F: 5}, {TotalVP: 17, Q: 12, F: 5}, {TotalVP: 18, Q: 12, F: 5}, {TotalVP: 19, Q: 13, F: 6}, {TotalVP: 20, Q: 14, F: 6},
		{TotalVP: 21, Q: 14, F: 6}, {TotalVP: 22, Q: 15, F: 7}, {TotalVP: 23, Q: 16, F: 7}, {TotalVP: 24, Q: 16, F: 7}, {TotalVP: 25, Q: 17, F: 8},
		{TotalVP: 26, Q: 18, F: 8}, {TotalVP: 27, Q: 18, F: 8}, {TotalVP: 28, Q: 19, F: 9}, {TotalVP: 29, Q: 20, F: 9}, {TotalVP: 30, Q: 20, F: 9},
		{TotalVP: 31, Q: 21, F: 10}, {TotalVP: 32, Q: 22, F: 10}, {TotalVP: 33, Q: 22, F: 10}, {TotalVP: 34, Q: 23, F: 11}, {TotalVP: 35, Q: 24, F: 11},
		{TotalVP: 36, Q: 24, F: 11}, {TotalVP: 37, Q: 25, F: 12}, {TotalVP: 38, Q: 26, F: 12}, {TotalVP: 39, Q: 26, F: 12}, {TotalVP: 40, Q: 27, F: 13},
		{TotalVP: 41, Q: 28, F: 13}, {TotalVP: 42, Q: 28, F: 13}, {TotalVP: 43, Q: 29, F: 14}, {TotalVP: 44, Q: 30, F: 14}, {TotalVP: 45, Q: 30, F: 14},
		{TotalVP: 46, Q: 31, F: 15}, {TotalVP: 47, Q: 32, F: 15}, {TotalVP: 48, Q: 32, F: 15}, {TotalVP: 49, Q: 33, F: 16}, {TotalVP: 50, Q: 34, F: 16},
		{TotalVP: 51, Q: 34, F: 16}, {TotalVP: 52, Q: 35, F: 17}, {TotalVP: 53, Q: 36, F: 17}, {TotalVP: 54, Q: 36, F: 17}, {TotalVP: 55, Q: 37, F: 18},
		{TotalVP: 56, Q: 38, F: 18}, {TotalVP: 57, Q: 38, F: 18}, {TotalVP: 58, Q: 39, F: 19}, {TotalVP: 59, Q: 40, F: 19}, {TotalVP: 60, Q: 40, F: 19},
		{TotalVP: 61, Q: 41, F: 20}, {TotalVP: 62, Q: 42, F: 20}, {TotalVP: 63, Q: 42, F: 20}, {TotalVP: 64, Q: 43, F: 21}, {TotalVP: 65, Q: 44, F: 21},
		{TotalVP: 66, Q: 44, F: 21}, {TotalVP: 67, Q: 45, F: 22}, {TotalVP: 68, Q: 46, F: 22}, {TotalVP: 69, Q: 46, F: 22}, {TotalVP: 70, Q: 47, F: 23},
		{TotalVP: 71, Q: 48, F: 23}, {TotalVP: 72, Q: 48, F: 23}, {TotalVP: 73, Q: 49, F: 24}, {TotalVP: 74, Q: 50, F: 24}, {TotalVP: 75, Q: 50, F: 24},
		{TotalVP: 76, Q: 51, F: 25}, {TotalVP: 77, Q: 52, F: 25}, {TotalVP: 78, Q: 52, F: 25}, {TotalVP: 79, Q: 53, F: 26}, {TotalVP: 80, Q: 54, F: 26},
		{TotalVP: 81, Q: 54, F: 26}, {TotalVP: 82, Q: 55, F: 27}, {TotalVP: 83, Q: 56, F: 27}, {TotalVP: 84, Q: 56, F: 27}, {TotalVP: 85, Q: 57, F: 28},
		{TotalVP: 86, Q: 58, F: 28}, {TotalVP: 87, Q: 58, F: 28}, {TotalVP: 88, Q: 59, F: 29}, {TotalVP: 89, Q: 60, F: 29}, {TotalVP: 90, Q: 60, F: 29},
		{TotalVP: 91, Q: 61, F: 30}, {TotalVP: 92, Q: 62, F: 30}, {TotalVP: 93, Q: 62, F: 30}, {TotalVP: 94, Q: 63, F: 31}, {TotalVP: 95, Q: 64, F: 31},
		{TotalVP: 96, Q: 64, F: 31}, {TotalVP: 97, Q: 65, F: 32}, {TotalVP: 98, Q: 66, F: 32}, {TotalVP: 99, Q: 66, F: 32}, {TotalVP: 100, Q: 67, F: 33},
	}

	for _, testCase := range testCases {
		c := createTestCommitteeMembers(t, genRandUint64(1, int(testCase.TotalVP)), testCase.TotalVP)
		set, err := NewRoundRobinSet(c, c.Members[0].Address)
		assertNilError(t, err)

		gotQ := set.Quorum()
		gotF := set.F()

		if testCase.F != gotF.Uint64() {
			t.Errorf("expected F: %v and got: %v", testCase.F, gotF)
		}

		if testCase.Q != gotQ.Uint64() {
			t.Errorf("expected Q: %v and got: %v", testCase.Q, gotQ)
		}
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

// totalPower >= n
func createTestCommitteeMembers(t *testing.T, n, totalPower int64) *types.Committee {
	t.Helper()
	c := new(types.Committee)

	if n > totalPower {
		t.Fatalf("totalPower >= size of c")
	}

	q, r := getTotalPowerDistribution(totalPower, n)
	for i := 0; i < int(n); i++ {
		key, err := crypto.GenerateKey()
		if err != nil {
			t.Fatal(err)
		}
		vp := q
		if i == int(n)-1 {
			vp += r
		}

		consensusKey, err := blst.RandKey()
		require.NoError(t, err)
		member := types.CommitteeMember{
			Address:      crypto.PubkeyToAddress(key.PublicKey),
			VotingPower:  new(big.Int).SetUint64(uint64(vp)),
			ConsensusKey: consensusKey.PublicKey().Marshal(),
		}
		c.Members = append(c.Members, &member)
	}

	if n > 0 {
		// swap 1st and last element if 1st element is less then last to ensure c is not sorted
		firstIndex, lastIndex := 0, len(c.Members)-1
		comp := bytes.Compare(c.Members[firstIndex].Bytes(), c.Members[lastIndex].Bytes())
		if comp < 0 {
			c.Members[firstIndex], c.Members[lastIndex] = c.Members[lastIndex], c.Members[firstIndex]
		}
	}

	return c
}

func getTotalPowerDistribution(p, n int64) (int64, int64) {
	return p / n, p % n

}

// generate random voting power in range [min...max]
func genRandUint64(min, max int) int64 {
	return int64(rand.Intn(max-min+1) + min)
}

func getMemberIndex(committee *types.Committee, memberAddr common.Address) int64 {
	for i, m := range committee.Members {
		if memberAddr == m.Address {
			return int64(i)
		}
	}
	return int64(-1)
}
