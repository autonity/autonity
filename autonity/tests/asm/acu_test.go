package asm

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/autonity/tests"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/internal/testrand"
	"github.com/autonity/autonity/params"
)

var Int256Max, _ = new(big.Int).SetString("57896044618658097711785492504343953926634992332820282019728792003956564819967", 10)
var OracleScaleFactor = new(big.Int).SetUint64(10000000)

const BigFloatPrecision = 256 // will allow full representation of the solidity uint256 range
func newFloat(value *big.Int) *big.Float {
	return new(big.Float).SetPrec(BigFloatPrecision).SetInt(value)
}

func TestACUConstructor(t *testing.T) {
	setup := func() *tests.Runner {
		return tests.Setup(t, nil)
	}

	tests.RunWithSetup("Test ACU constructor defaults", setup, func(r *tests.Runner) {
		round, _, err := r.Acu.Round(nil)
		require.NoError(t, err)

		require.Equal(t, uint64(0), round.Uint64())
		scaleFactor, _, err := r.Acu.ScaleFactor(nil)
		require.NoError(t, err)
		require.Equal(t, OracleScaleFactor, scaleFactor)
	})

	tests.RunWithSetup("Test ACU constructor errors", setup, func(r *tests.Runner) {
		// symbols and quantities must have the same length
		_, _, _, err := r.DeployACU(
			nil,
			[]string{"FOO"},
			nil,
			big.NewInt(1),
			common.Address{},
			common.Address{},
			common.Address{},
		)
		require.Error(t, err)
		require.ErrorAs(t, err, &tests.ACUInvalidBasketError{})

		// symbols and quantities must have the same length
		_, _, _, err = r.DeployACU(
			nil,
			[]string{},
			[]*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)},
			big.NewInt(1),
			common.Address{},
			common.Address{},
			common.Address{},
		)
		require.Error(t, err)
		require.ErrorAs(t, err, &tests.ACUInvalidBasketError{})

		// bad quantity
		_, _, _, err = r.DeployACU(
			nil,
			[]string{"FOO"},
			[]*big.Int{new(big.Int).Add(Int256Max, big.NewInt(1))},
			big.NewInt(1),
			common.Address{},
			common.Address{},
			common.Address{},
		)
		require.Error(t, err)
		require.ErrorAs(t, err, &tests.ACUInvalidBasketError{})
	})
}

func TestModifyACU(t *testing.T) {
	setup := func() *tests.Runner {
		return tests.Setup(t, nil)
	}

	tests.RunWithSetup("Test modify basket", setup, func(r *tests.Runner) {
		newSymbols := []string{"FOO", "BAR"}
		newQuantities := []*big.Int{big.NewInt(1), big.NewInt(2)}
		newScale := big.NewInt(18)
		newScaleFactor := new(big.Int).Exp(big.NewInt(10), newScale, nil)

		_, err := r.Acu.ModifyBasket(r.Operator, newSymbols, newQuantities, newScale)
		require.NoError(t, err)

		symbols, _, err := r.Acu.Symbols(nil)
		require.NoError(t, err)
		require.Equal(t, newSymbols, symbols)

		quantities, _, err := r.Acu.Quantities(nil)
		require.NoError(t, err)
		require.Equal(t, newQuantities, quantities)

		scale, _, err := r.Acu.Scale(nil)
		require.NoError(t, err)
		require.Equal(t, newScale, scale)

		scaleFactor, _, err := r.Acu.ScaleFactor(nil)
		require.NoError(t, err)
		require.Equal(t, newScaleFactor, scaleFactor)

		// ToDo(scott): verify emitted BasketModified event
	})

	tests.RunWithSetup("Test modify basket errors", setup, func(r *tests.Runner) {
		// unauthorized
		_, err := r.Acu.ModifyBasket(
			tests.FromSender(testrand.Address(), big.NewInt(0)),
			[]string{"FOO"},
			[]*big.Int{big.NewInt(1)},
			big.NewInt(1),
		)
		require.ErrorAs(t, err, &tests.ACUUnauthorizedError{})

		// invalid basket size
		_, err = r.Acu.ModifyBasket(r.Operator, []string{"FOO"}, []*big.Int{}, big.NewInt(1))
		require.ErrorAs(t, err, &tests.ACUInvalidBasketError{})

		// bad basket quantity
		_, err = r.Acu.ModifyBasket(
			r.Operator,
			[]string{"FOO"},
			[]*big.Int{new(big.Int).Add(Int256Max, big.NewInt(1))},
			big.NewInt(1),
		)
		require.ErrorAs(t, err, &tests.ACUInvalidBasketError{})
	})
}

func TestACUSetters(t *testing.T) {
	setup := func() *tests.Runner {
		return tests.Setup(t, nil)
	}

	tests.RunWithSetup("Test set operator", setup, func(r *tests.Runner) {
		newOperator := testrand.Address()
		_, err := r.Acu.SetOperator(tests.FromAutonity, newOperator)
		require.NoError(t, err)

		r.GiveMeSomeMoney(newOperator, big.NewInt(100000000000))
		_, err = r.Acu.ModifyBasket(
			tests.FromSender(newOperator, nil),
			[]string{"FOO"},
			[]*big.Int{big.NewInt(1)},
			big.NewInt(1),
		)
		require.NoError(t, err)
	})

	tests.RunWithSetup("Test set operator unauthorized", setup, func(r *tests.Runner) {
		users := []common.Address{params.DeployerAddress, params.TestAutonityContractConfig.Operator, testrand.Address()}
		for _, user := range users {
			_, err := r.Acu.SetOperator(tests.FromSender(user, nil), testrand.Address())
			require.ErrorAs(r.T, err, &tests.ACUUnauthorizedError{})
		}
	})

	tests.RunWithSetup("Test set oracle", setup, func(r *tests.Runner) {
		newOracle := testrand.Address()
		_, err := r.Acu.SetOracle(tests.FromAutonity, newOracle)
		require.NoError(t, err)
	})

	tests.RunWithSetup("Test set oracle unauthorized", setup, func(r *tests.Runner) {
		users := []common.Address{params.DeployerAddress, params.TestAutonityContractConfig.Operator, testrand.Address()}
		for _, user := range users {
			_, err := r.Acu.SetOracle(tests.FromSender(user, nil), testrand.Address())
			require.ErrorAs(r.T, err, &tests.ACUUnauthorizedError{})
		}
	})

}

func TestACUValue(t *testing.T) {
	setup := func() *tests.Runner {
		return tests.Setup(t, nil)
	}

	tests.RunWithSetup("Test no ACU value", setup, func(r *tests.Runner) {
		_, err := r.Acu.ModifyBasket(r.Operator, []string{"FOO"}, []*big.Int{big.NewInt(1)}, big.NewInt(1))
		require.NoError(t, err)

		_, _, err = r.Acu.Value(nil)
		require.ErrorAs(t, err, &tests.ACUNoACUValueError{})
	})

	tests.RunWithSetup("Test update", setup, func(r *tests.Runner) {
		symbols, prices, quantities, scaleFactor := primeACU(r)
		_, err := r.Acu.Update(tests.FromAutonity)
		require.NoError(t, err)

		value := computeAcu(symbols, prices, quantities, scaleFactor)

		acuValue, _, err := r.Acu.Value(nil)
		require.NoError(t, err)
		require.Equal(t, value, acuValue)

		acuRound, _, err := r.Acu.Round(nil)
		require.NoError(t, err)

		oracleRound, _, err := r.Oracle.GetRound(nil)
		require.NoError(t, err)
		require.Equal(t, new(big.Int).Sub(oracleRound, common.Big1), acuRound)

		// ToDo(scott): verify emitted ACUUpdated event
	})

	tests.RunWithSetup("Test update unauthorized", setup, func(r *tests.Runner) {
		primeACU(r)
		users := []common.Address{params.DeployerAddress, params.TestAutonityContractConfig.Operator, testrand.Address()}
		for _, user := range users {
			_, err := r.Acu.Update(tests.FromSender(user, nil))
			require.ErrorAs(t, err, &tests.ACUUnauthorizedError{})
		}
	})

	tests.RunWithSetup("Test update same round", setup, func(r *tests.Runner) {
		primeACU(r)
		_, err := r.Acu.Update(tests.FromAutonity)
		require.NoError(t, err)

		roundBefore, _, err := r.Acu.Round(nil)
		require.NoError(t, err)

		_, err = r.Acu.Update(tests.FromAutonity)
		require.NoError(t, err)

		roundAfter, _, err := r.Acu.Round(nil)
		require.NoError(t, err)
		require.Equal(t, roundBefore, roundAfter)

		// ToDo(scott): verify no ACUUpdated event, also is there a way to check the return value is false? I think not.
	})

	tests.RunWithSetup("Test update missing price", setup, func(r *tests.Runner) {
		symbols := params.DefaultAcuContractGenesis.Symbols
		oracleConfig, _, err := r.Oracle.Config(nil)
		require.NoError(t, err)
		r.WaitNBlocks(int(oracleConfig.VotePeriod.Int64()))
		// don't submit any votes to the oracle
		_, err = r.Acu.Update(tests.FromAutonity)
		require.NoError(t, err)

		for _, symbol := range symbols {
			report, _, err := r.Oracle.LatestRoundData(nil, symbol)
			require.NoError(t, err)
			require.Equal(t, false, report.Success)
		}

		_, _, err = r.Acu.Value(nil)
		require.ErrorAs(t, err, &tests.ACUNoACUValueError{})
	})
}

func TestViewFunctions(t *testing.T) {
	r := tests.Setup(t, nil)
	r.Run("Test view functions", func(r *tests.Runner) {
		symbols, _, err := r.Acu.Symbols(nil)
		require.NoError(t, err)
		require.Equal(t, params.DefaultAcuContractGenesis.Symbols, symbols)

		quantities, _, err := r.Acu.Quantities(nil)
		require.NoError(t, err)
		expectQuantities := make([]*big.Int, len(params.DefaultAcuContractGenesis.Quantities))
		for i, q := range params.DefaultAcuContractGenesis.Quantities {
			expectQuantities[i] = new(big.Int).SetUint64(q)
		}
		require.Equal(t, expectQuantities, quantities)

		scale, _, err := r.Acu.ScaleFactor(nil)
		require.NoError(t, err)
		require.Equal(t, new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(params.DefaultAcuContractGenesis.Scale)), nil), scale)
	})
}

func primeACU(r *tests.Runner) ([]string, []*big.Int, []*big.Int, *big.Int) {
	oracleScaleFactor := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	acuScale := big.NewInt(5)
	acuScaleFactor := new(big.Int).Exp(big.NewInt(10), acuScale, nil)

	symbols := []string{
		"AUD-USD",
		"CAD-USD",
		"EUR-USD",
		"GBP-USD",
		"JPY-USD",
		"USD-USD",
		"SEK-USD",
	}
	prices := []*big.Int{
		floatStringToBigInt("0.6757", oracleScaleFactor),
		floatStringToBigInt("0.75694", oracleScaleFactor),
		floatStringToBigInt("1.1085", oracleScaleFactor),
		floatStringToBigInt("1.29403", oracleScaleFactor),
		floatStringToBigInt("0.00713", oracleScaleFactor),
		floatStringToBigInt("1.0", oracleScaleFactor),
		floatStringToBigInt("0.09597", oracleScaleFactor),
	}
	quantities := []*big.Int{
		floatStringToBigInt("0.213", acuScaleFactor),
		floatStringToBigInt("0.187", acuScaleFactor),
		floatStringToBigInt("0.143", acuScaleFactor),
		floatStringToBigInt("0.104", acuScaleFactor),
		floatStringToBigInt("17.6", acuScaleFactor),
		floatStringToBigInt("0.180", acuScaleFactor),
		floatStringToBigInt("1.41", acuScaleFactor),
	}
	primeACUBasket(r, symbols, quantities, acuScale)
	primeOracle(r, symbols, prices)
	return symbols, prices, quantities, oracleScaleFactor
}

func primeOracle(r *tests.Runner, symbols []string, prices []*big.Int) {
	config, _, err := r.Oracle.Config(nil)
	voter := r.Committee.Validators[0].OracleAddress
	require.NoError(r.T, err)
	_, err = r.Oracle.SetSymbols(r.Operator, symbols)
	require.NoError(r.T, err)
	r.WaitNBlocks(2 * int(config.VotePeriod.Int64()))
	reports := make([]tests.IOracleReport, len(prices))
	for i, price := range prices {
		reports[i] = tests.IOracleReport{
			Price:      price,
			Confidence: 100,
		}
	}
	salt := big.NewInt(0)
	commit := tests.MakeOracleCommit(r.T, salt, voter, reports)
	// first vote
	_, err = r.Oracle.Vote(tests.FromSender(voter, common.Big0), commit, nil, big.NewInt(0), 0)
	require.NoError(r.T, err)

	r.WaitNBlocks(int(config.VotePeriod.Int64()))
	// second vote
	_, err = r.Oracle.Vote(tests.FromSender(voter, common.Big0), big.NewInt(0), reports, salt, 0)
	require.NoError(r.T, err)
	r.WaitNBlocks(int(config.VotePeriod.Int64()))

	// verify prices
	for i := range symbols {
		oracleRound, _, err := r.Oracle.GetRound(nil)
		require.NoError(r.T, err)
		roundData, _, err := r.Oracle.GetRoundData(nil, new(big.Int).Sub(oracleRound, big.NewInt(1)), symbols[i])
		require.NoError(r.T, err)
		require.Equal(r.T, roundData.Price, prices[i])
	}
}

func primeACUBasket(r *tests.Runner, symbols []string, quantities []*big.Int, scale *big.Int) {
	_, err := r.Acu.ModifyBasket(r.Operator, symbols, quantities, scale)
	require.NoError(r.T, err)
}

func floatStringToBigInt(s string, scaleFactor *big.Int) *big.Int {
	f, _ := new(big.Float).SetString(s)
	result, _ := new(big.Float).Mul(f, new(big.Float).SetInt(scaleFactor)).Int(nil)
	return result
}

func computeAcu(symbols []string, prices []*big.Int, quantities []*big.Int, scaleFactor *big.Int) *big.Int {
	acu := newFloat(big.NewInt(0))
	for i := range symbols {
		var price *big.Float
		if symbols[i] == "USD-USD" {
			price = newFloat(scaleFactor)
		} else {
			price = newFloat(prices[i])
		}
		acu = new(big.Float).Add(acu, new(big.Float).Mul(price, newFloat(quantities[i])))
	}
	acu = new(big.Float).Quo(acu, newFloat(scaleFactor))
	result, _ := acu.Int(nil)
	return result
}
