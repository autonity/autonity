package economic

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/math"
	"github.com/shopspring/decimal"
	"math/big"
	"math/rand"
	"time"
)

type TXNPacker interface {
	packTXNs(parent *block) (uint64, uint64)
}

type fullPacker struct {
	params *systemParams
}

func (f *fullPacker) packTXNs(parent *block) (uint64, uint64) {
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

type twoThirdPacker struct {
	params *systemParams
}

func (f *twoThirdPacker) packTXNs(parent *block) (uint64, uint64) {
	// Define the gas cost per transaction
	const gasPerTxn = 21000

	// Generate a random number of additional transactions (1 to 5)
	rand.Seed(time.Now().UnixNano())
	additionalTxns := rand.Intn(5) + 1

	// Calculate the minimum gas required to exceed the parent's gas limit
	minGasRequired := parent.gasLimit*2/3 + uint64(additionalTxns*gasPerTxn)

	// Calculate the number of transactions needed to exceed the gas limit
	numTxns := minGasRequired / gasPerTxn

	// Calculate the total gas used
	usedGas := numTxns * gasPerTxn

	return numTxns, usedGas
}

type halfPacker struct {
	params *systemParams
}

func (f *halfPacker) packTXNs(parent *block) (uint64, uint64) {
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

type oneThirdPacker struct {
	params *systemParams
}

func (f *oneThirdPacker) packTXNs(parent *block) (uint64, uint64) {
	// Define the gas cost per transaction
	const gasPerTxn = 21000

	// Generate a random number of additional transactions (1 to 5)
	rand.Seed(time.Now().UnixNano())
	additionalTxns := rand.Intn(5) + 1

	// Calculate the minimum gas required to exceed the parent's gas limit
	minGasRequired := parent.gasLimit/3 + uint64(additionalTxns*gasPerTxn)

	// Calculate the number of transactions needed to exceed the gas limit
	numTxns := minGasRequired / gasPerTxn

	// Calculate the total gas used
	usedGas := numTxns * gasPerTxn

	return numTxns, usedGas
}

type dustFiller struct {
	params *systemParams
}

func (f *dustFiller) packTXNs(_ *block) (uint64, uint64) {
	// Define the gas cost per transaction
	const gasPerTxn = 21000

	// Generate a random number of additional transactions (1 to 5)
	rand.Seed(time.Now().UnixNano())
	numTxns := uint64(rand.Intn(5)) + 1

	// Calculate the total gas used
	usedGas := numTxns * gasPerTxn

	return numTxns, usedGas
}

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
