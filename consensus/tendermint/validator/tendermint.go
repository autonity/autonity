package validator

import (
	"bytes"
	"fmt"
	"math"
	"math/big"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/common/heap"
	"github.com/clearmatics/autonity/consensus/tendermint"
)

// The maximum allowed total voting power.
// We set the ProposerPriority of freshly added validators to -1.125*totalVotingPower.
// To compute 1.125*totalVotingPower efficiently, we do:
// totalVotingPower + (totalVotingPower >> 3) because
// x + (x >> 3) = x + x/8 = x * (1 + 0.125).
// MaxTotalVotingPower is the largest int64 `x` with the property that `x + (x >> 3)` is
// still in the bounds of int64.
const MaxTotalVotingPower = int64(8198552921648689607)

// tendermintProposer returns the current proposer. If the validator set is empty, nil
// is returned.
func tendermintProposer(valSet tendermint.ValidatorSet, _ common.Address, oldround, round uint64) tendermint.Validator {
	size := valSet.Size()
	if size == 0 {
		return nil
	}

	if round == 0 || oldround < round {
		valSet.IncrementProposerPriority(int(round - oldround))
	}

	return valSet.GetHighest()
}

func (vals *defaultSet) GetHighest() tendermint.Validator {
	var proposer tendermint.Validator
	for _, val := range vals.List() {
		if proposer == nil || val.Address() != proposer.Address() {
			proposer = CompareProposerPriority(proposer, val)
		}
	}

	return proposer.Copy()
}

// Returns the one with higher ProposerPriority.
func CompareProposerPriority(v tendermint.Validator, other tendermint.Validator) tendermint.Validator {
	if v == nil {
		return other
	}

	if v.ProposerPriority() > other.ProposerPriority() {
		return v
	} else if v.ProposerPriority() < other.ProposerPriority() {
		return other
	} else {
		result := bytes.Compare(v.Address().Bytes(), other.Address().Bytes())
		if result < 0 {
			return v
		} else if result > 0 {
			return other
		} else {
			panic("Cannot compare identical validators")
			return nil
		}
	}
}

// IncrementProposerPriority increments ProposerPriority of each validator and updates the
// proposer. Panics if validator set is empty.
// `times` must be positive.
func (vals *defaultSet) IncrementProposerPriority(times int) {
	if times <= 0 {
		panic("Cannot call IncrementProposerPriority with non-positive times")
	}

	var proposer tendermint.Validator
	const shiftEveryNthIter = 10
	// call IncrementProposerPriority(1) times times:
	for i := 0; i < times; i++ {
		shiftByAvgProposerPriority := i%shiftEveryNthIter == 0
		proposer = vals.incrementProposerPriority(shiftByAvgProposerPriority)
	}
	isShiftedAvgOnLastIter := (times-1)%shiftEveryNthIter == 0
	if !isShiftedAvgOnLastIter {
		validatorsHeap := heap.New()
		vals.shiftByAvgProposerPriority(validatorsHeap)
	}
	vals.proposer = proposer.Copy()
}

func (vals *defaultSet) incrementProposerPriority(subAvg bool) tendermint.Validator {
	for _, val := range vals.validators {
		// Check for overflow for sum.
		val.SetProposerPriority(safeAddClip(val.ProposerPriority(), val.VotingPower()))
	}

	validatorsHeap := heap.New()
	if subAvg { // shift by avg ProposerPriority
		vals.shiftByAvgProposerPriority(validatorsHeap)
	} else { // just update the heap
		for _, val := range vals.validators {
			validatorsHeap.PushComparable(val, proposerPriorityComparable{val})
		}
	}

	// Decrement the validator with most ProposerPriority:
	mostest := validatorsHeap.Peek().(tendermint.Validator)
	// mind underflow
	mostest.SetProposerPriority(safeSubClip(mostest.ProposerPriority(), vals.TotalVotingPower()))

	return mostest
}

// TotalVotingPower returns the sum of the voting powers of all validators.
func (vals *defaultSet) TotalVotingPower() int64 {
	if vals.totalVotingPower == 0 {
		sum := int64(0)
		for _, val := range vals.validators {
			// mind overflow
			sum = safeAddClip(sum, val.VotingPower())
		}
		if sum > MaxTotalVotingPower {
			panic(fmt.Sprintf(
				"Total voting power should be guarded to not exceed %v; got: %v",
				MaxTotalVotingPower,
				sum))
		}
		vals.totalVotingPower = sum
	}
	return vals.totalVotingPower
}

func (vals *defaultSet) computeAvgProposerPriority() int64 {
	n := int64(len(vals.validators))
	sum := big.NewInt(0)
	for _, val := range vals.validators {
		sum.Add(sum, big.NewInt(val.ProposerPriority()))
	}
	avg := sum.Div(sum, big.NewInt(n))
	if avg.IsInt64() {
		return avg.Int64()
	}

	// this should never happen: each val.ProposerPriority is in bounds of int64
	panic(fmt.Sprintf("Cannot represent avg ProposerPriority as an int64 %v", avg))
}

func (vals *defaultSet) shiftByAvgProposerPriority(validatorsHeap *heap.Heap) {
	avgProposerPriority := vals.computeAvgProposerPriority()
	for _, val := range vals.validators {
		val.SetProposerPriority(safeSubClip(val.ProposerPriority(), avgProposerPriority))
		validatorsHeap.PushComparable(val, proposerPriorityComparable{val})
	}
}

//------------------
// Use with Heap for sorting validators by ProposerPriority

type proposerPriorityComparable struct {
	tendermint.Validator
}

// We want to find the validator with the greatest ProposerPriority.
func (ac proposerPriorityComparable) Less(o interface{}) bool {
	other := o.(proposerPriorityComparable).Validator
	larger := CompareProposerPriority(ac, other)
	return bytes.Equal(larger.Address().Bytes(), ac.Address().Bytes())
}

///////////////////////////////////////////////////////////////////////////////
// Safe addition/subtraction

func safeAdd(a, b int64) (int64, bool) {
	if b > 0 && a > math.MaxInt64-b {
		return -1, true
	} else if b < 0 && a < math.MinInt64-b {
		return -1, true
	}
	return a + b, false
}

func safeSub(a, b int64) (int64, bool) {
	if b > 0 && a < math.MinInt64+b {
		return -1, true
	} else if b < 0 && a > math.MaxInt64+b {
		return -1, true
	}
	return a - b, false
}

func safeAddClip(a, b int64) int64 {
	c, overflow := safeAdd(a, b)
	if overflow {
		if b < 0 {
			return math.MinInt64
		}
		return math.MaxInt64
	}
	return c
}

func safeSubClip(a, b int64) int64 {
	c, overflow := safeSub(a, b)
	if overflow {
		if b > 0 {
			return math.MinInt64
		}
		return math.MaxInt64
	}
	return c
}
