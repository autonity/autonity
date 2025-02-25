package economic

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/math"
	"github.com/autonity/autonity/core"
	"github.com/shopspring/decimal"
	"math/big"
	"math/rand"
	"time"
)

// CalcBaseFee calculates the basefee of the header.
func CalcBaseFee(parent *block, params *systemParams) *big.Int {
	var (
		parentGasTarget          = parent.gasLimit / params.elasticityMultiplier
		parentGasTargetBig       = new(big.Int).SetUint64(parentGasTarget)
		baseFeeChangeDenominator = new(big.Int).SetUint64(params.baseFeeChangeDenominator)
	)
	// If the parent gasUsed is the same as the target, the baseFee remains unchanged.
	if parent.gasUsed == parentGasTarget {
		return new(big.Int).Set(parent.baseFee)
	}
	if parent.gasUsed > parentGasTarget {
		// If the parent block used more gas than its target, the baseFee should increase.
		gasUsedDelta := new(big.Int).SetUint64(parent.gasUsed - parentGasTarget)
		x := new(big.Int).Mul(parent.baseFee, gasUsedDelta)
		y := x.Div(x, parentGasTargetBig)
		baseFeeDelta := math.BigMax(
			x.Div(y, baseFeeChangeDenominator),
			common.Big1,
		)
		return x.Add(parent.baseFee, baseFeeDelta)
	} else {
		// Otherwise if the parent block used less gas than its target, the baseFee should decrease.
		gasUsedDelta := new(big.Int).SetUint64(parentGasTarget - parent.gasUsed)
		x := new(big.Int).Mul(parent.baseFee, gasUsedDelta)
		y := x.Div(x, parentGasTargetBig)
		baseFeeDelta := x.Div(y, baseFeeChangeDenominator)

		minBaseFee := big.NewInt(0).SetUint64(params.minBaseFee.Uint64())
		return math.BigMax(
			x.Sub(parent.baseFee, baseFeeDelta),
			minBaseFee,
		)
	}
}

type blockFiller interface {
	fillBlock(parent *block) (*block, *big.Int)
}

type fullFiller struct {
	params *systemParams
}

func (f *fullFiller) genTXNs(parent *block) (uint64, uint64) {
	// Define the gas cost per transaction
	const gasPerTxn = 21000

	// Generate a random number of additional transactions (1 to 5)
	rand.Seed(time.Now().UnixNano())
	additionalTxns := rand.Intn(5) + 1

	// Calculate the minimum gas required to exceed the parent's gas limit
	minGasRequired := parent.gasLimit + uint64(additionalTxns*gasPerTxn)

	// Calculate the number of transactions needed to exceed the gas limit
	numTxns := minGasRequired / gasPerTxn

	// Calculate the total gas used
	usedGas := numTxns * gasPerTxn

	return numTxns, usedGas
}

func baseFeeChangeRate(parentBaseFee *big.Int, curBaseFee *big.Int) decimal.Decimal {
	// Convert parentBaseFee and curBaseFee to decimal.Decimal
	parentFee := decimal.NewFromBigInt(parentBaseFee, 0)
	curFee := decimal.NewFromBigInt(curBaseFee, 0)

	// Calculate the difference between curBaseFee and parentBaseFee
	diff := curFee.Sub(parentFee)

	// Calculate the change rate: (diff / parentFee) * 100
	changeRate := diff.Mul(decimal.NewFromInt(100)).Div(parentFee)

	return changeRate
}

func (f *fullFiller) fillBlock(parent *block) (*block, *big.Int) {
	baseFee := CalcBaseFee(parent, f.params)
	changeRate := baseFeeChangeRate(parent.baseFee, baseFee)

	numOfTXN, gasUsed := f.genTXNs(parent)

	atnRewards := new(big.Int).Mul(new(big.Int).SetUint64(gasUsed), baseFee)
	return &block{
		timestamp:         parent.timestamp + 1,
		number:            parent.number + 1,
		gasUsed:           gasUsed,
		numOfTXNs:         numOfTXN,
		gasLimit:          core.CalcGasLimit(parent.gasLimit, f.params.gasCeil),
		baseFee:           baseFee,
		baseFeeChangeRate: changeRate,
	}, atnRewards
}

type halfFiller struct {
	params *systemParams
}

func (f *halfFiller) genTXNs(parent *block) (uint64, uint64) {
	// Define the gas cost per transaction
	const gasPerTxn = 21000

	// Generate a random number of additional transactions (1 to 5)
	rand.Seed(time.Now().UnixNano())
	additionalTxns := rand.Intn(5) + 1

	// Calculate the minimum gas required to exceed the parent's gas limit
	minGasRequired := parent.gasLimit/2 + uint64(additionalTxns*gasPerTxn)

	// Calculate the number of transactions needed to exceed the gas limit
	numTxns := minGasRequired / gasPerTxn

	// Calculate the total gas used
	usedGas := numTxns * gasPerTxn

	return numTxns, usedGas
}

func (f *halfFiller) fillBlock(parent *block) (*block, *big.Int) {
	baseFee := CalcBaseFee(parent, f.params)
	numOfTXN, gasUsed := f.genTXNs(parent)

	atnRewards := new(big.Int).Mul(new(big.Int).SetUint64(gasUsed), baseFee)

	return &block{
		timestamp: parent.timestamp + 1,
		number:    parent.number + 1,
		gasUsed:   gasUsed,
		numOfTXNs: numOfTXN,
		gasLimit:  core.CalcGasLimit(parent.gasLimit, f.params.gasCeil),
		baseFee:   baseFee,
	}, atnRewards
}

type dustFiller struct {
	params *systemParams
}

func (f *dustFiller) genTXNs(_ *block) (uint64, uint64) {
	// Define the gas cost per transaction
	const gasPerTxn = 21000

	// Generate a random number of additional transactions (1 to 5)
	rand.Seed(time.Now().UnixNano())
	numTxns := uint64(rand.Intn(5)) + 1

	// Calculate the total gas used
	usedGas := numTxns * gasPerTxn

	return numTxns, usedGas
}

func (f *dustFiller) fillBlock(parent *block) (*block, *big.Int) {
	baseFee := CalcBaseFee(parent, f.params)
	numOfTXN, gasUsed := f.genTXNs(parent)

	atnRewards := new(big.Int).Mul(new(big.Int).SetUint64(gasUsed), baseFee)

	return &block{
		timestamp: parent.timestamp + 1,
		number:    parent.number + 1,
		gasUsed:   gasUsed,
		numOfTXNs: numOfTXN,
		gasLimit:  core.CalcGasLimit(parent.gasLimit, f.params.gasCeil),
		baseFee:   baseFee,
	}, atnRewards
}
