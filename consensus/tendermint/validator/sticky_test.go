package validator

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/clearmatics/autonity/common"
	"github.com/golang/mock/gomock"
)

func TestCalcSeedNotFoundProposer(t *testing.T) {
	proposerAddress := common.BytesToAddress(bytes.Repeat([]byte{1}, common.AddressLength))

	testCases := []struct {
		validatorIndex int
		round          uint64

		resultOffset uint64
	}{
		{
			round:        0,
			resultOffset: 0,
		},
		{
			round:        1,
			resultOffset: 1,
		},
		{
			round:        2,
			resultOffset: 2,
		},
		{
			round:        10,
			resultOffset: 10,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(fmt.Sprintf("validator index %d, validator is nil %v, round %d", testCase.validatorIndex, true, testCase.round), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			validatorSet := NewMockSet(ctrl)
			validatorSet.EXPECT().
				GetByAddress(gomock.Eq(proposerAddress)).
				Return(testCase.validatorIndex, nil)

			res := calcSeed(validatorSet, proposerAddress, testCase.round)
			if res != testCase.resultOffset {
				t.Errorf("got %d, expected %d", res, testCase.resultOffset)
			}
		})
	}
}

func TestCalcSeedWithProposer(t *testing.T) {
	proposerAddress := common.BytesToAddress(bytes.Repeat([]byte{1}, common.AddressLength))

	testCases := []struct {
		validatorIndex int
		round          uint64

		resultOffset uint64
	}{
		{
			validatorIndex: 0,
			round:          0,
			resultOffset:   0,
		},
		{
			validatorIndex: 1,
			round:          0,
			resultOffset:   1,
		},
		{
			validatorIndex: 2,
			round:          0,
			resultOffset:   2,
		},
		{
			validatorIndex: 10,
			round:          0,
			resultOffset:   10,
		},

		{
			validatorIndex: 0,
			round:          1,
			resultOffset:   1,
		},
		{
			validatorIndex: 1,
			round:          1,
			resultOffset:   2,
		},
		{
			validatorIndex: 2,
			round:          1,
			resultOffset:   3,
		},
		{
			validatorIndex: 10,
			round:          1,
			resultOffset:   11,
		},

		{
			validatorIndex: 0,
			round:          2,
			resultOffset:   2,
		},
		{
			validatorIndex: 1,
			round:          2,
			resultOffset:   3,
		},
		{
			validatorIndex: 2,
			round:          2,
			resultOffset:   4,
		},
		{
			validatorIndex: 10,
			round:          2,
			resultOffset:   12,
		},

		{
			validatorIndex: 0,
			round:          10,
			resultOffset:   10,
		},
		{
			validatorIndex: 1,
			round:          10,
			resultOffset:   11,
		},
		{
			validatorIndex: 2,
			round:          10,
			resultOffset:   12,
		},
		{
			validatorIndex: 10,
			round:          10,
			resultOffset:   20,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(fmt.Sprintf("validator index %d, validator is nil %v, round %d", testCase.validatorIndex, false, testCase.round), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			validator := NewMockValidator(ctrl)
			validatorSet := NewMockSet(ctrl)
			validatorSet.EXPECT().
				GetByAddress(gomock.Eq(proposerAddress)).
				Return(testCase.validatorIndex, validator)

			res := calcSeed(validatorSet, proposerAddress, testCase.round)
			if res != testCase.resultOffset {
				t.Errorf("got %d, expected %d", res, testCase.resultOffset)
			}
		})
	}
}

func TestStickyProposerZeroSize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	proposerAddress := common.BytesToAddress(bytes.Repeat([]byte{1}, common.AddressLength))
	proposerZeroAddress := common.Address{}

	testCases := []struct {
		size     int
		round    uint64
		proposer common.Address
	}{
		{
			size:     0,
			round:    0,
			proposer: proposerZeroAddress,
		},
		{
			size:     0,
			round:    1,
			proposer: proposerZeroAddress,
		},
		{
			size:     0,
			round:    2,
			proposer: proposerZeroAddress,
		},
		{
			size:     0,
			round:    10,
			proposer: proposerZeroAddress,
		},

		{
			size:     0,
			round:    0,
			proposer: proposerAddress,
		},
		{
			size:     0,
			round:    1,
			proposer: proposerAddress,
		},
		{
			size:     0,
			round:    2,
			proposer: proposerAddress,
		},
		{
			size:     0,
			round:    10,
			proposer: proposerAddress,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(fmt.Sprintf("validator is zero address, round %d", testCase.round), func(t *testing.T) {
			validatorSet := NewMockSet(ctrl)

			validatorSet.EXPECT().
				Size().
				Return(testCase.size)

			val := stickyProposer(validatorSet, proposerAddress, testCase.round)
			if val != nil {
				t.Errorf("got wrond validator %v, expected nil", val)
			}
		})
	}
}

func TestStickyProposer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	proposerAddress := common.BytesToAddress(bytes.Repeat([]byte{1}, common.AddressLength))
	proposerZeroAddress := common.Address{}

	testCases := []struct {
		size     int
		round    uint64
		proposer common.Address
		pick     uint64
	}{
		// size is greater than pick
		{
			size:     10,
			round:    0,
			proposer: proposerZeroAddress,
			pick:     0,
		},
		{
			size:     10,
			round:    1,
			proposer: proposerZeroAddress,
			pick:     1,
		},
		{
			size:     10,
			round:    2,
			proposer: proposerZeroAddress,
			pick:     2,
		},
		{
			size:     10,
			round:    8,
			proposer: proposerZeroAddress,
			pick:     8,
		},
		// non-zero address
		{
			size:     10,
			round:    0,
			proposer: proposerAddress,
			pick:     1,
		},
		{
			size:     10,
			round:    1,
			proposer: proposerAddress,
			pick:     2,
		},
		{
			size:     10,
			round:    2,
			proposer: proposerAddress,
			pick:     3,
		},
		{
			size:     10,
			round:    8,
			proposer: proposerAddress,
			pick:     9,
		},

		// size is  less or equal to pick
		{
			size:     3,
			round:    0,
			proposer: proposerZeroAddress,
			pick:     0,
		},
		{
			size:     3,
			round:    1,
			proposer: proposerZeroAddress,
			pick:     1,
		},
		{
			size:     3,
			round:    2,
			proposer: proposerZeroAddress,
			pick:     2,
		},
		{
			size:     3,
			round:    3,
			proposer: proposerZeroAddress,
			pick:     0,
		},
		{
			size:     3,
			round:    10,
			proposer: proposerZeroAddress,
			pick:     1,
		},
		// non-zero address
		{
			size:     3,
			round:    0,
			proposer: proposerAddress,
			pick:     1,
		},
		{
			size:     3,
			round:    1,
			proposer: proposerAddress,
			pick:     2,
		},
		{
			size:     3,
			round:    2,
			proposer: proposerAddress,
			pick:     0,
		},
		{
			size:     3,
			round:    10,
			proposer: proposerAddress,
			pick:     2,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(fmt.Sprintf("validator set size %d, proposer address %s, round %d", testCase.size, testCase.proposer.String(), testCase.round), func(t *testing.T) {
			validatorSet := NewMockSet(ctrl)

			validatorSet.EXPECT().
				Size().
				Return(testCase.size)

			if testCase.proposer != proposerZeroAddress {
				index := 1
				validator := NewMockValidator(ctrl)
				validatorSet.EXPECT().
					GetByAddress(gomock.Eq(testCase.proposer)).
					Return(index, validator)
			}

			expectedValidator := NewMockValidator(ctrl)
			validatorSet.EXPECT().
				GetByIndex(gomock.Eq(testCase.pick)).
				Return(expectedValidator)

			val := stickyProposer(validatorSet, testCase.proposer, testCase.round)
			if !reflect.DeepEqual(val, expectedValidator) {
				t.Errorf("got wrond validator %v, expected %v", val, expectedValidator)
			}
		})
	}
}
