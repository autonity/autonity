package economic

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/math"
	"github.com/shopspring/decimal"
	"math/big"
	"math/rand"
	"time"
)

const gasPerTxn = 21000

type TXNPacker interface {
	packTXNs(num uint64, gasLimit uint64) (uint64, uint64)
}

type RatedPacker struct {
	ratio decimal.Decimal
}

func (f *RatedPacker) packTXNs(_ uint64, gasLimit uint64) (uint64, uint64) {
	// Calculate the total gas to be used based on the ratio
	totalGasToUse := decimal.NewFromInt(int64(gasLimit)).Mul(f.ratio)

	// Each transaction costs 21,000 gas
	gasPerTXN := decimal.NewFromInt(gasPerTxn)

	// Calculate the number of transactions that can fit
	numTXNs := totalGasToUse.Div(gasPerTXN).Floor()

	// Calculate the actual gas used
	gasUsed := numTXNs.Mul(gasPerTXN)

	// Return the number of transactions and the gas used
	return numTXNs.BigInt().Uint64(), gasUsed.BigInt().Uint64()
}

type dynamicPacker struct {
	increasingInterval int // in blocks
	decreasingInterval int // in blocks
}

// dynamic packer just pack TXN with an increasing interval in blocks, within which, the packer pack TXN in an increasing
// way, which means new block always has more TXNs than its parent. While after that interval, it enters into the
// decreasing interval in blocks, within which, the packer start to packing TXNs in a linear decreasing way.
func (d *dynamicPacker) packTXNs(number uint64, gasLimit uint64) (uint64, uint64) {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Determine the current block's position in the interval
	intervalPosition := (number + 1) % uint64(d.increasingInterval+d.decreasingInterval)

	// Calculate the number of transactions based on the interval
	var numTxns uint64
	if intervalPosition < uint64(d.increasingInterval) {
		// Increasing interval: pack more transactions than the parent block
		additionalTxns := rand.Intn(5) + 1 // Randomly add 1 to 5 transactions
		numTxns = gasLimit/gasPerTxn + uint64(additionalTxns)
	} else {
		// Decreasing interval: pack half of the parent block's transactions
		// drop to dust block?
		numTxns = 1
		/*
			numTxns = parent.gasLimit / gasPerTxn / 2
			if numTxns < 1 {
				numTxns = 1 // Ensure at least 1 transaction is packed
			}*/
	}

	// Calculate the total gas used
	usedGas := numTxns * gasPerTxn

	return numTxns, usedGas
}

// CalcBaseFee calculates the basefee of the header.
func CalcBaseFee(parent *block, params *systemParams) *big.Int {
	var (
		parentGasTarget          = parent.gasLimit / params.ElasticityMultiplier
		parentGasTargetBig       = new(big.Int).SetUint64(parentGasTarget)
		baseFeeChangeDenominator = new(big.Int).SetUint64(params.BaseFeeChangeDenominator)
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

		minBaseFee := big.NewInt(0).SetUint64(params.MinBaseFee.Uint64())
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
