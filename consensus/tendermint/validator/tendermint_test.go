package validator

import (
	"bytes"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"reflect"
	"strings"
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

		{
			vals.Copy(),
			[]int64{
				// Acumm+VotingPower-Avg:
				0 + 10 - 12 - 4, // mostest will be subtracted by total voting power (12)
				0 + 1 - 4,
				0 + 1 - 4},
			0,
			vals.GetByIndex(0)},
		{
			vals.Copy(),
			[]int64{
				0 + 10 - 12 - 4, // this will be mostest on 2nd iter, too
				0 + 1 - 4,
				0 + 1 - 4},
			1,
			vals.GetByIndex(0)}, // increment twice -> expect average to be subtracted twice
		{
			vals.Copy(),
			[]int64{
				((0 + 10 - 12 - 4) + 10 - 12) + 10 - 12 + 4, // still mostest
				((0 + 1 - 4) + 1) + 1 + 4,
				((0 + 1 - 4) + 1) + 1 + 4},
			2,
			vals.GetByIndex(0)},
		{
			vals.Copy(),
			[]int64{
				0 + 4*(10-12) + 4 - 4, // still mostest
				0 + 4*1 + 4 - 4,
				0 + 4*1 + 4 - 4},
			3,
			vals.GetByIndex(0)},
		{
			vals.Copy(),
			[]int64{
				0 + 4*(10-12) + 10 + 4 - 4, // 4 iters was mostest
				0 + 5*1 - 12 + 4 - 4,       // now this val is mostest for the 1st time (hence -12==totalVotingPower)
				0 + 5*1 + 4 - 4},
			4,
			vals.GetByIndex(1)},
		{
			vals.Copy(),
			[]int64{
				0 + 6*10 - 5*12 + 4 - 4, // mostest again
				0 + 6*1 - 12 + 4 - 4,    // mostest once up to here
				0 + 6*1 + 4 - 4},
			5,
			vals.GetByIndex(0)},
		{
			vals.Copy(),
			[]int64{
				0 + 7*10 - 6*12 + 4 - 4, // in 7 iters this val is mostest 6 times
				0 + 7*1 - 12 + 4 - 4,    // in 7 iters this val is mostest 1 time
				0 + 7*1 + 4 - 4},
			6,
			vals.GetByIndex(0)},
		{
			vals.Copy(),
			[]int64{
				0 + 8*10 - 7*12 + 4 - 4, // mostest
				0 + 8*1 - 12 + 4 - 4,
				0 + 8*1 + 4 - 4},
			7,
			vals.GetByIndex(0)},
		{
			vals.Copy(),
			[]int64{
				0 + 9*10 - 7*12 + 4 - 4,
				0 + 9*1 - 12 + 4 - 4,
				0 + 9*1 - 12 + 4 - 4}, // mostest
			8,
			vals.GetByIndex(2)},
		{
			vals.Copy(),
			[]int64{
				0 + 10*10 - 8*12 + 4 - 4, // after 10 iters this is mostest again
				0 + 10*1 - 12 + 4 - 4,    // after 6 iters this val is "mostest" once and not in between
				0 + 10*1 - 12 + 4 - 4},   // in between 10 iters this val is "mostest" once
			9,
			vals.GetByIndex(0)},
		{
			vals.Copy(),
			[]int64{
				0 + 11*10 - 9*12 + 4 - 4,
				0 + 11*1 - 12 + 4 - 4,  // after 6 iters this val is "mostest" once and not in between
				0 + 11*1 - 12 + 4 - 4}, // after 10 iters this val is "mostest" once
			10,
			vals.GetByIndex(0),
		},
		{
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
		t.Run(fmt.Sprintf("case %d, times %d", i, tc.times), func(t *testing.T) {
			if i > 1 {
				// newDefaultSet call did the first IncrementProposerPriority
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
	cases := []struct{
		a, b int64
		result int64
	}{
		{math.MaxInt64, 10, math.MaxInt64},
		{math.MaxInt64, math.MaxInt64, math.MaxInt64},
		{math.MinInt64, -10, math.MinInt64},
	}

	for _, testCase := range cases {
		if res := safeAddClip(testCase.a, testCase.b); res != testCase.result {
			t.Fatalf("test case a=%d b=%d, got %d expected %d", testCase.a, testCase.b, res, testCase.result)
		}
	}
}

func TestSafeSubClip(t *testing.T) {
	cases := []struct{
		a, b int64
		result int64
	}{
		{math.MinInt64, 10, math.MinInt64},
		{math.MinInt64, math.MinInt64, 0},
		{math.MinInt64, math.MaxInt64, math.MinInt64},
		{math.MaxInt64, -10, math.MaxInt64},
	}

	for _, testCase := range cases {
		if res := safeSubClip(testCase.a, testCase.b); res != testCase.result {
			t.Fatalf("test case a=%d b=%d, got %d expected %d", testCase.a, testCase.b, res, testCase.result)
		}
	}
}


func TestProposerSelectionManyRounds(t *testing.T) {
	validatorStorage := newValidatorStorage()
	vset := newDefaultSet(tendermint.Tendermint,
		validatorStorage.getValidator(0, 1000, 0),
		validatorStorage.getValidator(1, 300, 0),
		validatorStorage.getValidator(2, 330, 0),
	)

	var proposers []string
	for i := 0; i < 99; i++ {
		val := vset.GetProposer()
		proposers = append(proposers, val.Address().String())
		vset.IncrementProposerPriority(1)
	}
	expected := `0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000002 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000001 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000002 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000001 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000002 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000001 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000002 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000001 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000002 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000001 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000002 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000001 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000002 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000001 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000002 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000001 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000002 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000001 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000002 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000001 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000002 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000002 0x0000000000000000000000000000000000000001 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000002 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000001 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000002 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000001 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000002 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000001 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000002 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000001 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000002 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000001 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000002 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000001 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000002 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000001 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000002 0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000000`
	if expected != strings.Join(proposers, " ") {
		t.Errorf("Expected sequence of proposers was\n%v\nbut got \n%v", expected, strings.Join(proposers, " "))
	}
}

func TestProposerSelection2(t *testing.T) {
	// when all voting power is same, we go in order of addresses
	validatorStorage := newValidatorStorage()
	valList := []tendermint.Validator{
		validatorStorage.getValidator(0, 100, 0),
		validatorStorage.getValidator(1, 100, 0),
		validatorStorage.getValidator(2, 100, 0),
	}
	vset := newDefaultSet(tendermint.Tendermint, valList...)

	for i := 0; i < len(valList)*5; i++ {
		ii := (i) % len(valList)
		prop := vset.GetProposer()
		if !bytes.Equal(prop.Address().Bytes(), valList[ii].Address().Bytes()) {
			t.Fatalf("(%d): Expected %X. Got %X", i, valList[ii].Address(), prop.Address())
		}
		vset.IncrementProposerPriority(1)
	}
}

func TestProposerSelectionOneValidatorOnce(t *testing.T) {
	// One validator has more than the others, but not enough to propose twice in a row
	validatorStorage := newValidatorStorage()
	valList := []tendermint.Validator{
		validatorStorage.getValidator(0, 100, 0),
		validatorStorage.getValidator(1, 100, 0),
		validatorStorage.getValidator(2, 400, 0),
	}

	vset := newDefaultSet(tendermint.Tendermint, valList...)
	prop := vset.GetProposer()
	if !bytes.Equal(prop.Address().Bytes(), valList[2].Address().Bytes()) {
		t.Fatalf("Expected address with highest voting power to be first proposer. Got %X", prop.Address())
	}
	vset.IncrementProposerPriority(1)
	prop = vset.GetProposer()
	if !bytes.Equal(prop.Address().Bytes(), valList[0].Address().Bytes()) {
		t.Fatalf("Expected smallest address to be validator. Got %X", prop.Address())
	}
}

func TestProposerSelectionOneValidatorTwice(t *testing.T) {
	// One validator has more than the others, and enough to be proposer twice in a row
	validatorStorage := newValidatorStorage()
	valList := []tendermint.Validator{
		validatorStorage.getValidator(0, 100, 0),
		validatorStorage.getValidator(1, 100, 0),
		validatorStorage.getValidator(2, 401, 0),
	}
	vset := newDefaultSet(tendermint.Tendermint, valList...)
	prop := vset.GetProposer()
	if !bytes.Equal(prop.Address().Bytes(), valList[2].Address().Bytes()) {
		t.Fatalf("Expected address with highest voting power to be first proposer. Got %X", prop.Address())
	}
	vset.IncrementProposerPriority(1)
	prop = vset.GetProposer()
	if !bytes.Equal(prop.Address().Bytes(), valList[2].Address().Bytes()) {
		t.Fatalf("Expected address with highest voting power to be second proposer. Got %X", prop.Address())
	}
	vset.IncrementProposerPriority(1)
	prop = vset.GetProposer()
	if !bytes.Equal(prop.Address().Bytes(), valList[0].Address().Bytes()) {
		t.Fatalf("Expected smallest address to be validator. Got %X", prop.Address())
	}
}

func TestProposerSelectionDistribution(t *testing.T) {
	// each validator should be the proposer a proportional number of times
	validatorStorage := newValidatorStorage()
	valList := []tendermint.Validator{
		validatorStorage.getValidator(0, 4, 0),
		validatorStorage.getValidator(1, 5, 0),
		validatorStorage.getValidator(2, 3, 0),
	}
	vset := newDefaultSet(tendermint.Tendermint, valList...)
	propCount := make([]int, 3)
	N := 1
	for i := 0; i < 120*N; i++ {
		prop := vset.GetProposer()
		bytesAddress := prop.Address().Bytes()
		ii := bytesAddress[len(bytesAddress)-1]
		propCount[ii]++
		vset.IncrementProposerPriority(1)
	}

	if propCount[0] != 40*N {
		t.Fatalf("Expected prop count for validator with 4/12 of voting power to be %d/%d. Got %d/%d", 40*N, 120*N, propCount[0], 120*N)
	}
	if propCount[1] != 50*N {
		t.Fatalf("Expected prop count for validator with 5/12 of voting power to be %d/%d. Got %d/%d", 50*N, 120*N, propCount[1], 120*N)
	}
	if propCount[2] != 30*N {
		t.Fatalf("Expected prop count for validator with 3/12 of voting power to be %d/%d. Got %d/%d", 30*N, 120*N, propCount[2], 120*N)
	}
}

func TestProposerSelection3(t *testing.T) {
	validatorStorage := newValidatorStorage()
	valList := []tendermint.Validator{
		validatorStorage.getValidator(0, 1, 0),
		validatorStorage.getValidator(1, 1, 0),
		validatorStorage.getValidator(2, 1, 0),
		validatorStorage.getValidator(3, 1, 0),
	}
	vset := newDefaultSet(tendermint.Tendermint, valList...)
	proposerOrder := make([]tendermint.Validator, 4)
	for i := 0; i < 4; i++ {
		proposerOrder[i] = vset.GetProposer()
		vset.IncrementProposerPriority(1)
	}

	// i for the loop
	// j for the times
	// we should go in order for ever, despite some IncrementProposerPriority with times > 1
	var i, j int
	for ; i < 10000; i++ {
		got := vset.GetProposer().Address().Bytes()
		expected := proposerOrder[j%4].Address().Bytes()
		if !bytes.Equal(got, expected) {
			t.Fatalf(fmt.Sprintf("vset.Proposer (%X) does not match expected proposer (%X) for (%d, %d)", got, expected, i, j))
		}

		computed := vset.GetProposer() // findGetProposer()
		if i != 0 {
			if !bytes.Equal(got, computed.Address().Bytes()) {
				t.Fatalf(fmt.Sprintf("vset.Proposer (%X) does not match computed proposer (%X) for (%d, %d)", got, computed.Address(), i, j))
			}
		}

		// times is usually 1
		times := 1
		mod := (rand.Int() % 5) + 1
		if rand.Int()%mod > 0 {
			// sometimes its up to 5
			times = (rand.Int() % 4) + 1
		}
		vset.IncrementProposerPriority(times)

		j += times
	}
}
