package asm

import (
	"math/big"
	"strings"
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

	tests.RunWithSetup("sqrt increase auction amount performs correctly", setup, func(r *tests.Runner) {
		duration := big.NewInt(1000)
		collateral := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
		liquidationRatio := toBase("1.5", 18)

		// function is C/(1 - (L-1)*(1 - sqrt((currentTime - startTime)/duration)))
		expected := map[string]*big.Int{
			"0":    toBase("0.666666666666666666", 18), // C/L
			"50":   toBase("0.720359060949715970", 18),
			"100":  toBase("0.745219722700413018", 18),
			"200":  toBase("0.783457635340899531", 18),
			"500":  toBase("0.872260419102717064", 18),
			"800":  toBase("0.949860290487784360", 18),
			"900":  toBase("0.974983530382842913", 18),
			"1000": collateral, // 1.0
			"1500": collateral, // 1.0 stays at max amount
		}

		for at, expectedAmount := range expected {
			currentTime, _ := new(big.Int).SetString(at, 10)
			result, _, err := stabilizationMath.SqrtIncreaseAuctionAmount(
				nil,
				common.Big0, // start time
				currentTime, // current time
				collateral,
				liquidationRatio,
				duration,
			)
			require.NoError(t, err)
			require.Equal(t, expectedAmount, result)
		}
	})

	tests.RunWithSetup("linear decrease auction amount performs correctly", setup, func(r *tests.Runner) {
		duration := big.NewInt(1000)
		maximumOffer := toBase("2.0", 18)
		minimumOffer := big.NewInt(0)

		// function is max - (max - min) * (currentTime - startTime) / duration
		expected := map[string]*big.Int{
			"0":    maximumOffer,
			"50":   toBase("1.900000000000000000", 18),
			"100":  toBase("1.800000000000000000", 18),
			"200":  toBase("1.600000000000000000", 18),
			"533":  toBase("0.934000000000000000", 18),
			"800":  toBase("0.400000000000000000", 18),
			"900":  toBase("0.200000000000000000", 18),
			"1000": minimumOffer, // 1.0
			"1500": minimumOffer, // 1.0 stays at max amount
		}

		for at, expectedAmount := range expected {
			currentTime, _ := new(big.Int).SetString(at, 10)
			result, _, err := stabilizationMath.LinearDecreaseAuctionAmount(
				nil,
				common.Big0, // start time
				currentTime, // current time
				minimumOffer,
				maximumOffer,
				duration,
			)
			require.NoError(t, err)
			require.Equal(t, expectedAmount.String(), result.String())
		}
	})
}

func toBase(s string, decimals int64) *big.Int {
	parts := strings.Split(s, ".")
	before := parts[0]
	after := ""
	if len(parts) > 1 {
		after = parts[1]
	}
	combined := before + after
	result := new(big.Int)
	result.SetString(combined, 10)
	exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(decimals-int64(len(after))), nil)
	result.Mul(result, exp)
	return result
}
