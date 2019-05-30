package validator

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"math/big"
	"reflect"
	"testing"
	"testing/quick"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint"
	"github.com/golang/mock/gomock"
)

func TestTendermintProposerZeroSize(t *testing.T) {
	testCases := []struct {
		size     int
		oldRound uint64
		round    uint64
	}{
		{
			size:     0,
			oldRound: 0,
			round:    0,
		},
		{
			size:     0,
			oldRound: 0,
			round:    1,
		},
		{
			size:     0,
			oldRound: 1,
			round:    2,
		},
		{
			size:     0,
			oldRound: 9,
			round:    10,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(fmt.Sprintf("validator is zero address, round %d", testCase.round), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			validatorSet := tendermint.NewMockValidatorSet(ctrl)

			validatorSet.EXPECT().
				Size().
				Times(1).
				Return(testCase.size)

			val := tendermintProposer(validatorSet, common.Address{}, testCase.oldRound, testCase.round)
			if val != nil {
				t.Errorf("got wrong validator %v, expected nil", val)
			}
		})
	}
}

func TestTendermintProposer(t *testing.T) {
	testCases := []struct {
		size     int
		oldRound uint64
		round    uint64
	}{
		{
			size:     5,
			oldRound: 0,
			round:    0,
		},
		{
			size:     5,
			oldRound: 0,
			round:    1,
		},
		{
			size:     5,
			oldRound: 1,
			round:    2,
		},
		{
			size:     5,
			oldRound: 9,
			round:    10,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(fmt.Sprintf("validator is zero address, round %d", testCase.round), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			validatorSet := tendermint.NewMockValidatorSet(ctrl)
			validatorStorage := newValidatorStorage()

			validatorSet.EXPECT().
				Size().
				Times(1).
				Return(testCase.size)

			validatorSet.EXPECT().
				IncrementProposerPriority(int(testCase.round - testCase.oldRound)).
				Times(1).
				Return()

			expectedValidator := validatorStorage.getValidator(0, 100, 100)
			validatorSet.EXPECT().
				GetHighest().
				Times(1).
				Return(expectedValidator)

			val := tendermintProposer(validatorSet, common.Address{}, testCase.oldRound, testCase.round)
			if !reflect.DeepEqual(val, expectedValidator) {
				t.Errorf("got wrong validator %v, expected %v", val, expectedValidator)
			}
		})
	}
}

type validatorStorage struct {
	m map[int]tendermint.Validator
}

func newValidatorStorage() *validatorStorage {
	return &validatorStorage{make(map[int]tendermint.Validator)}
}

func (v *validatorStorage) getValidators(n int, votingPower int64) []tendermint.Validator {
	validators := make([]tendermint.Validator, n)
	for i := 0; i < n; i++ {
		val, ok := v.m[i]
		if ok {
			validators[i] = val.Copy()
			v.m[i] = val.Copy()
			continue
		}

		val = &defaultValidator{
			address:     common.BigToAddress(big.NewInt(int64(i))),
			votingPower: votingPower,
		}

		v.m[i] = val.Copy()
		validators[i] = val.Copy()
	}

	return validators
}

func (v *validatorStorage) getValidator(i int, votingPower, proposerPriority int64) tendermint.Validator {
	val, ok := v.m[i]
	if ok {
		return val.Copy()
	}

	val = &defaultValidator{
		address:     common.BigToAddress(big.NewInt(int64(i))),
		votingPower: votingPower,
	}

	v.m[i] = val

	return val.Copy()
}

func TestAveragingInIncrementProposerPriorityWithVotingPower(t *testing.T) {
	// Other than TestAveragingInIncrementProposerPriority this is a more complete test showing
	// how each ProposerPriority changes in relation to the validator's voting power respectively.
	vals := newDefaultSet(
		tendermint.Tendermint,
		New(common.BigToAddress(big.NewInt(0)), 10),
		New(common.BigToAddress(big.NewInt(1)), 1),
		New(common.BigToAddress(big.NewInt(2)), 1),
	)
	tcs := []struct {
		vals                  tendermint.ValidatorSet
		wantProposerPrioritys []int64
		times                 int
		wantProposer          tendermint.Validator
	}{

		0: {
			vals.Copy(),
			[]int64{
				// Acumm+VotingPower-Avg:
				0 + 10 - 12 - 4, // mostest will be subtracted by total voting power (12)
				0 + 1 - 4,
				0 + 1 - 4},
			0,
			vals.GetByIndex(0)},
		1: {
			vals.Copy(),
			[]int64{
				0 + 10 - 12 - 4, // this will be mostest on 2nd iter, too
				0 + 1 - 4,
				0 + 1 - 4},
			1,
			vals.GetByIndex(0)}, // increment twice -> expect average to be subtracted twice
		2: {
			vals.Copy(),
			[]int64{
				((0 + 10 - 12 - 4) + 10 - 12) + 10 - 12 + 4, // still mostest
				((0 + 1 - 4) + 1) + 1 + 4,
				((0 + 1 - 4) + 1) + 1 + 4},
			2,
			vals.GetByIndex(0)},
		3: {
			vals.Copy(),
			[]int64{
				0 + 4*(10-12) + 4 - 4, // still mostest
				0 + 4*1 + 4 - 4,
				0 + 4*1 + 4 - 4},
			3,
			vals.GetByIndex(0)},
		4: {
			vals.Copy(),
			[]int64{
				0 + 4*(10-12) + 10 + 4 - 4, // 4 iters was mostest
				0 + 5*1 - 12 + 4 - 4,       // now this val is mostest for the 1st time (hence -12==totalVotingPower)
				0 + 5*1 + 4 - 4},
			4,
			vals.GetByIndex(1)},
		5: {
			vals.Copy(),
			[]int64{
				0 + 6*10 - 5*12 + 4 - 4, // mostest again
				0 + 6*1 - 12 + 4 - 4,    // mostest once up to here
				0 + 6*1 + 4 - 4},
			5,
			vals.GetByIndex(0)},
		6: {
			vals.Copy(),
			[]int64{
				0 + 7*10 - 6*12 + 4 - 4, // in 7 iters this val is mostest 6 times
				0 + 7*1 - 12 + 4 - 4,    // in 7 iters this val is mostest 1 time
				0 + 7*1 + 4 - 4},
			6,
			vals.GetByIndex(0)},
		7: {
			vals.Copy(),
			[]int64{
				0 + 8*10 - 7*12 + 4 - 4, // mostest
				0 + 8*1 - 12 + 4 - 4,
				0 + 8*1 + 4 - 4},
			7,
			vals.GetByIndex(0)},
		8: {
			vals.Copy(),
			[]int64{
				0 + 9*10 - 7*12 + 4 - 4,
				0 + 9*1 - 12 + 4 - 4,
				0 + 9*1 - 12 + 4 - 4}, // mostest
			8,
			vals.GetByIndex(2)},
		9: {
			vals.Copy(),
			[]int64{
				0 + 10*10 - 8*12 + 4 - 4, // after 10 iters this is mostest again
				0 + 10*1 - 12 + 4 - 4,    // after 6 iters this val is "mostest" once and not in between
				0 + 10*1 - 12 + 4 - 4},   // in between 10 iters this val is "mostest" once
			9,
			vals.GetByIndex(0)},
		10: {
			vals.Copy(),
			[]int64{
				0 + 11*10 - 9*12 + 4 - 4,
				0 + 11*1 - 12 + 4 - 4,  // after 6 iters this val is "mostest" once and not in between
				0 + 11*1 - 12 + 4 - 4}, // after 10 iters this val is "mostest" once
			10,
			vals.GetByIndex(0),
		},
		11: {
			vals.Copy(),
			[]int64{
				// shift twice inside incrementProposerPriority (shift every 10th iter);
				// don't shift at the end of IncremenctProposerPriority
				// last avg should be zero because
				// ProposerPriority of validator 0: (0 + 11*10 - 8*12 - 4) == 10
				// ProposerPriority of validator 1 and 2: (0 + 11*1 - 12 - 4) == -5
				// and (10 + 5 - 5) / 3 == 0
				0 + 12*10 - 10*12 - 4 - 0,
				0 + 12*1 - 12 - 4 - 0,  // after 6 iters this val is "mostest" once and not in between
				0 + 12*1 - 12 - 4 - 0}, // after 10 iters this val is "mostest" once
			11,
			vals.GetByIndex(0),
		},
	}

	for i, tc := range tcs {
		i := i
		tc := tc

		t.Run(fmt.Sprintf("case %d, times %d", i, tc.times), func(t *testing.T) {
			if i > 1 {
				tc.vals.IncrementProposerPriority(tc.times)
			}

			if tc.wantProposer.Address().String() != tc.vals.GetProposer().Address().String() {
				t.Fatalf("got wrong proposer %v, expected %v", tc.vals.GetProposer().Address().String(), tc.wantProposer.Address().String())
			}

			for valIdx, val := range tc.vals.List() {
				if tc.wantProposerPrioritys[valIdx] != val.ProposerPriority() {
					t.Fatalf("got wrong validator proposer priority %v(index %d), expected %v. List: %v",
						val.ProposerPriority(),
						valIdx,
						tc.wantProposerPrioritys[valIdx],
						tc.vals.List(),
					)
				}
			}

			fmt.Println("              ")
		})
	}
}

func TestValidatorSetTotalVotingPowerPanicsOnOverflow(t *testing.T) {
	// NewValidatorSet calls IncrementProposerPriority which calls TotalVotingPower()
	// which should panic on overflows:

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("the code should panic")
		}
	}()

	validatorStorage := newValidatorStorage()
	newDefaultSet(tendermint.Tendermint,
		validatorStorage.getValidator(0, math.MaxInt64, 0),
		validatorStorage.getValidator(1, math.MaxInt64, 0),
		validatorStorage.getValidator(2, math.MaxInt64, 0),
	)
}

/*
func TestAvgProposerPriority(t *testing.T) {
	// Create Validator set without calling IncrementProposerPriority:
	tcs := []struct {
		vs   ValidatorSet
		want int64
	}{
		0: {ValidatorSet{Validators: []*Validator{{ProposerPriority: 0}, {ProposerPriority: 0}, {ProposerPriority: 0}}}, 0},
		1: {ValidatorSet{Validators: []*Validator{{ProposerPriority: math.MaxInt64}, {ProposerPriority: 0}, {ProposerPriority: 0}}}, math.MaxInt64 / 3},
		2: {ValidatorSet{Validators: []*Validator{{ProposerPriority: math.MaxInt64}, {ProposerPriority: 0}}}, math.MaxInt64 / 2},
		3: {ValidatorSet{Validators: []*Validator{{ProposerPriority: math.MaxInt64}, {ProposerPriority: math.MaxInt64}}}, math.MaxInt64},
		4: {ValidatorSet{Validators: []*Validator{{ProposerPriority: math.MinInt64}, {ProposerPriority: math.MinInt64}}}, math.MinInt64},
	}
	for i, tc := range tcs {
		got := tc.vs.computeAvgProposerPriority()
		assert.Equal(t, tc.want, got, "test case: %v", i)
	}
}

func TestAveragingInIncrementProposerPriority(t *testing.T) {
	// Test that the averaging works as expected inside of IncrementProposerPriority.
	// Each validator comes with zero voting power which simplifies reasoning about
	// the expected ProposerPriority.
	tcs := []struct {
		vs    ValidatorSet
		times int
		avg   int64
	}{
		0: {ValidatorSet{
			Validators: []*Validator{
				{Address: []byte("a"), ProposerPriority: 1},
				{Address: []byte("b"), ProposerPriority: 2},
				{Address: []byte("c"), ProposerPriority: 3}}},
			1, 2},
		1: {ValidatorSet{
			Validators: []*Validator{
				{Address: []byte("a"), ProposerPriority: 10},
				{Address: []byte("b"), ProposerPriority: -10},
				{Address: []byte("c"), ProposerPriority: 1}}},
			// this should average twice but the average should be 0 after the first iteration
			// (voting power is 0 -> no changes)
			11, 1 / 3},
		2: {ValidatorSet{
			Validators: []*Validator{
				{Address: []byte("a"), ProposerPriority: 100},
				{Address: []byte("b"), ProposerPriority: -10},
				{Address: []byte("c"), ProposerPriority: 1}}},
			1, 91 / 3},
	}
	for i, tc := range tcs {
		// work on copy to have the old ProposerPriorities:
		newVset := tc.vs.CopyIncrementProposerPriority(tc.times)
		for _, val := range tc.vs.Validators {
			_, updatedVal := newVset.GetByAddress(val.Address)
			assert.Equal(t, updatedVal.ProposerPriority, val.ProposerPriority-tc.avg, "test case: %v", i)
		}
	}
}
*/

func TestSafeAdd(t *testing.T) {
	f := func(a, b int64) bool {
		c, overflow := safeAdd(a, b)
		return overflow || (!overflow && c == a+b)
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestSafeAddClip(t *testing.T) {
	assert.EqualValues(t, math.MaxInt64, safeAddClip(math.MaxInt64, 10))
	assert.EqualValues(t, math.MaxInt64, safeAddClip(math.MaxInt64, math.MaxInt64))
	assert.EqualValues(t, math.MinInt64, safeAddClip(math.MinInt64, -10))
}

func TestSafeSubClip(t *testing.T) {
	assert.EqualValues(t, math.MinInt64, safeSubClip(math.MinInt64, 10))
	assert.EqualValues(t, 0, safeSubClip(math.MinInt64, math.MinInt64))
	assert.EqualValues(t, math.MinInt64, safeSubClip(math.MinInt64, math.MaxInt64))
	assert.EqualValues(t, math.MaxInt64, safeSubClip(math.MaxInt64, -10))
}
