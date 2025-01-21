package asm

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/autonity/tests"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/params"
)

func TestAuctionAmounts(t *testing.T) {
	var stabilizationMath *tests.StabilizationMathTest
	setup := func() *tests.Runner {
		r := tests.Setup(t, nil)
		var err error
		_, _, stabilizationMath, err = r.DeployStabilizationMathTest(tests.FromSender(params.DeployerAddress, common.Big0))
		require.NoError(t, err)
		return r
	}

	tests.RunWithSetup("Sqrt increase auction amount performs correctly", setup, func(r *tests.Runner) {
		duration := big.NewInt(1000)
		maxAmount := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
		minAmount := new(big.Int).Div(maxAmount, big.NewInt(2))

		// function is max - (max-min)*(1 - sqrt(currentTime - startTime)/sqrt(duration))
		expected := map[string]*big.Int{
			"0":    minAmount, // 0.5
			"50":   toBase("0.611803398874989485", 18),
			"100":  toBase("0.658113883008418967", 18),
			"200":  toBase("0.723606797749978970", 18),
			"500":  toBase("0.853553390593273762", 18),
			"800":  toBase("0.947213595499957939", 18),
			"900":  toBase("0.974341649025256900", 18),
			"1000": maxAmount, // 1.0
			"1500": maxAmount, // 1.0 stays at max amount
		}

		for at, expectedAmount := range expected {
			currentTime, _ := new(big.Int).SetString(at, 10)
			result, _, err := stabilizationMath.SqrtIncreaseAuctionAmount(
				nil,
				common.Big0, // start time
				currentTime, // current time
				maxAmount,
				minAmount,
				duration,
			)
			require.NoError(t, err)
			require.Equal(t, expectedAmount, result)
		}

	})
}

func toBase(s string, decimals int64) *big.Int {
	f, _, _ := new(big.Float).Parse(s, 10)
	result, _ := new(big.Float).Mul(
		f,
		new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(decimals), nil)),
	).Int(nil)
	return result
}
