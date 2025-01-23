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

var e18 = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
var e12 = new(big.Int).Exp(big.NewInt(10), big.NewInt(12), nil)

// newtonPrice = 1234567 * 10^12 = 1.234567 * 10^18
var newtonPrice = new(big.Int).Mul(big.NewInt(1234567), e12)
var basicConfig = tests.IStabilizationConfig{
	BorrowInterestRate:        new(big.Int).Div(e18, big.NewInt(2)),
	LiquidationRatio:          new(big.Int).Mul(big.NewInt(15), new(big.Int).Exp(big.NewInt(10), big.NewInt(17), nil)),
	MinCollateralizationRatio: new(big.Int).Mul(big.NewInt(25), new(big.Int).Exp(big.NewInt(10), big.NewInt(17), nil)),
	MinDebtRequirement:        new(big.Int).Exp(big.NewInt(10), big.NewInt(6), nil),
	TargetPrice:               new(big.Int).Set(e18),
}

func TestStabilizationConstructor(t *testing.T) {
	setup := func() *tests.Runner {
		return tests.Setup(t, nil)
	}

	tests.RunWithSetup("Test constructor zero mcr", setup, func(r *tests.Runner) {
		_, _, _, err := r.DeployStabilization(
			nil,
			tests.IStabilizationConfig{
				BorrowInterestRate:        basicConfig.BorrowInterestRate,
				LiquidationRatio:          basicConfig.LiquidationRatio,
				MinCollateralizationRatio: big.NewInt(0),
				MinDebtRequirement:        basicConfig.MinDebtRequirement,
				TargetPrice:               basicConfig.TargetPrice,
			},
			r.Autonity.Address(),
			params.TestAutonityContractConfig.Operator,
			r.Oracle.Address(),
			r.SupplyControl.Address(),
			r.Auctioneer.Address(),
			r.Acu.Address(),
			r.Autonity.Address(),
		)
		require.ErrorAs(r.T, err, &tests.StabilizationInvalidParameterError{})
	})

	tests.RunWithSetup("Test constructor invalid ratios", setup, func(r *tests.Runner) {
		// liquidation ratio == min collateralization ratio
		_, _, _, err := r.DeployStabilization(
			nil,
			tests.IStabilizationConfig{
				BorrowInterestRate:        basicConfig.BorrowInterestRate,
				LiquidationRatio:          e18,
				MinCollateralizationRatio: e18,
				MinDebtRequirement:        basicConfig.MinDebtRequirement,
				TargetPrice:               basicConfig.TargetPrice,
			},
			r.Autonity.Address(),
			params.TestAutonityContractConfig.Operator,
			r.Oracle.Address(),
			r.SupplyControl.Address(),
			r.Auctioneer.Address(),
			r.Acu.Address(),
			r.Autonity.Address(),
		)
		require.ErrorAs(r.T, err, &tests.StabilizationInvalidParameterError{})

		// liquidation ratio > min collateralization ratio
		_, _, _, err = r.DeployStabilization(
			nil,
			tests.IStabilizationConfig{
				BorrowInterestRate:        basicConfig.BorrowInterestRate,
				LiquidationRatio:          new(big.Int).Add(basicConfig.MinCollateralizationRatio, big.NewInt(1)),
				MinCollateralizationRatio: basicConfig.MinCollateralizationRatio,
				MinDebtRequirement:        basicConfig.MinDebtRequirement,
				TargetPrice:               basicConfig.TargetPrice,
			},
			r.Autonity.Address(),
			params.TestAutonityContractConfig.Operator,
			r.Oracle.Address(),
			r.SupplyControl.Address(),
			r.Auctioneer.Address(),
			r.Acu.Address(),
			r.Autonity.Address(),
		)
		require.ErrorAs(r.T, err, &tests.StabilizationInvalidParameterError{})
	})
}

func TestStabilizationDeposit(t *testing.T) {
	userAccount := testrand.Address()
	secondUserAccount := testrand.Address()
	fundedAmount := new(big.Int).Mul(e18, big.NewInt(100))
	setup := func() *tests.Runner {
		r := tests.Setup(t, nil)
		for _, account := range []common.Address{userAccount, secondUserAccount} {
			r.GiveMeSomeMoney(account, new(big.Int).Mul(e18, big.NewInt(100)))
			_, err := r.Autonity.Mint(r.Operator, account, fundedAmount)
			require.NoError(t, err)

			// approve stabilization
			_, err = r.Autonity.Approve(tests.FromSender(account, nil), r.Stabilization.Address(), fundedAmount)
			require.NoError(t, err)
		}
		// remove cdp restrictions for this teat
		_, err := r.Stabilization.RemoveCDPRestrictions(r.Operator)
		require.NoError(t, err)
		return r
	}

	tests.RunWithSetup("Test deposit zero not allowed", setup, func(r *tests.Runner) {
		_, err := r.Stabilization.Deposit(tests.FromSender(userAccount, nil), big.NewInt(0))
		require.ErrorAs(r.T, err, &tests.StabilizationInvalidAmountError{})
	})

	tests.RunWithSetup("Test deposit initial", setup, func(r *tests.Runner) {
		depositAmount := new(big.Int).Div(fundedAmount, big.NewInt(2))
		deposit(r, userAccount, depositAmount)
		// check cdp
		cdp, _, err := r.Stabilization.Cdps(nil, userAccount)
		require.NoError(t, err)
		require.Greater(t, cdp.Timestamp.Int64(), int64(0))
		require.Equal(t, depositAmount, cdp.Collateral)

		accounts, _, err := r.Stabilization.Accounts(nil)
		require.NoError(t, err)
		require.Equal(t, []common.Address{userAccount}, accounts)
		// ToDo(scott): verify deposit events
	})

	tests.RunWithSetup("Test deposit subsequent", setup, func(r *tests.Runner) {
		depositAmount := new(big.Int).Div(fundedAmount, big.NewInt(10))
		deposit(r, userAccount, depositAmount)
		deposit(r, userAccount, depositAmount)
		// check cdp
		cdp, _, err := r.Stabilization.Cdps(nil, userAccount)
		require.NoError(t, err)
		require.Greater(t, cdp.Timestamp.Int64(), int64(0))
		require.Equal(t, new(big.Int).Mul(depositAmount, big.NewInt(2)), cdp.Collateral)

		accounts, _, err := r.Stabilization.Accounts(nil)
		require.NoError(t, err)
		require.Equal(t, []common.Address{userAccount}, accounts) // no duplicates
		// ToDo(scott): verify deposit events
	})

	tests.RunWithSetup("Test insufficient allowance", setup, func(r *tests.Runner) {
		_, err := r.Autonity.Approve(tests.FromSender(userAccount, nil), r.Stabilization.Address(), big.NewInt(0))
		require.NoError(t, err)

		_, err = r.Stabilization.Deposit(tests.FromSender(userAccount, nil), big.NewInt(1))
		require.ErrorAs(r.T, err, &tests.StabilizationInsufficientAllowanceError{})
	})

	tests.RunWithSetup("Test deposit with a second user", setup, func(r *tests.Runner) {
		depositAmount := new(big.Int).Div(fundedAmount, big.NewInt(10))
		deposit(r, userAccount, depositAmount)
		deposit(r, secondUserAccount, depositAmount)
		// check cdp
		cdp, _, err := r.Stabilization.Cdps(nil, secondUserAccount)
		require.NoError(t, err)
		require.Greater(t, cdp.Timestamp.Int64(), int64(0))
		require.Equal(t, depositAmount, cdp.Collateral)

		accounts, _, err := r.Stabilization.Accounts(nil)
		require.NoError(t, err)
		require.ElementsMatch(t, []common.Address{userAccount, secondUserAccount}, accounts)
		// ToDo(scott): verify deposit events
	})
}

func TestStabilizationWithdraw(t *testing.T) {
	userAccount := testrand.Address()
	fundedAmount := new(big.Int).Mul(e18, big.NewInt(100))
	totalDeposit := new(big.Int).Div(fundedAmount, big.NewInt(10))

	setup := func() *tests.Runner {
		r := tests.Setup(t, nil)
		setBasicConfig(r)
		primePrices(r, newtonPrice, toBase("0.97", 18))
		r.GiveMeSomeMoney(userAccount, new(big.Int).Mul(e18, big.NewInt(100)))
		_, err := r.Autonity.Mint(r.Operator, userAccount, fundedAmount)
		require.NoError(t, err)

		// approve stabilization
		_, err = r.Autonity.Approve(tests.FromSender(userAccount, nil), r.Stabilization.Address(), fundedAmount)
		require.NoError(t, err)

		// remove cdp restrictions for this teat
		_, err = r.Stabilization.RemoveCDPRestrictions(r.Operator)
		require.NoError(t, err)
		deposit(r, userAccount, totalDeposit)
		return r
	}

	tests.RunWithSetup("Test withdraw zero not allowed", setup, func(r *tests.Runner) {
		_, err := r.Stabilization.Withdraw(tests.FromSender(userAccount, nil), big.NewInt(0))
		require.ErrorAs(r.T, err, &tests.StabilizationInvalidAmountError{})
	})

	tests.RunWithSetup("Test withdraw full deposit", setup, func(r *tests.Runner) {
		balanceBefore, _, err := r.Autonity.BalanceOf(nil, userAccount)
		require.NoError(t, err)

		_, err = r.Stabilization.Withdraw(tests.FromSender(userAccount, nil), totalDeposit)
		require.NoError(t, err)

		balanceAfter, _, err := r.Autonity.BalanceOf(nil, userAccount)
		require.NoError(t, err)

		require.Equal(t, new(big.Int).Add(balanceBefore, totalDeposit), balanceAfter)

		// check cdp
		cdp, _, err := r.Stabilization.Cdps(nil, userAccount)
		require.NoError(t, err)
		require.Equal(t, int64(0), cdp.Collateral.Int64())
		// ToDo(scott): verify withdraw events
	})

	tests.RunWithSetup("Test cannot withdraw more than collateral balance", setup, func(r *tests.Runner) {
		_, err := r.Stabilization.Withdraw(tests.FromSender(userAccount, nil), new(big.Int).Add(totalDeposit, big.NewInt(1)))
		require.ErrorAs(r.T, err, &tests.StabilizationInvalidAmountError{})
	})

	tests.RunWithSetup("Test cannot withdraw while liquidatable", setup, func(r *tests.Runner) {
		cfg, _, err := r.Stabilization.Config(nil)
		require.NoError(t, err)
		borrowLimit := calcBorrowLimit(r, userAccount)

		_, err = r.Stabilization.Borrow(tests.FromSender(userAccount, nil), borrowLimit)
		require.NoError(t, err)

		liqidatable, _, err := r.Stabilization.IsLiquidatable(nil, userAccount)
		require.NoError(t, err)
		require.False(t, liqidatable)

		_, err = r.Stabilization.SetMinCollateralizationRatio(r.Operator, new(big.Int).Add(cfg.MinCollateralizationRatio, e18))
		require.NoError(t, err)
		_, err = r.Stabilization.SetLiquidationRatio(r.Operator, cfg.MinCollateralizationRatio)
		require.NoError(t, err)

		liqidatable, _, err = r.Stabilization.IsLiquidatable(nil, userAccount)
		require.NoError(t, err)
		require.True(t, liqidatable)

		_, err = r.Stabilization.Withdraw(tests.FromSender(userAccount, nil), big.NewInt(1))
		require.ErrorAs(r.T, err, &tests.StabilizationLiquidatableError{})
	})

	tests.RunWithSetup("Test cannot withdraw with insufficient collateral (change in mcr)", setup, func(r *tests.Runner) {
		cfg, _, err := r.Stabilization.Config(nil)
		require.NoError(t, err)

		borrowLimit := calcBorrowLimit(r, userAccount)

		// borrow
		r.NoError(r.Stabilization.Borrow(tests.FromSender(userAccount, nil), borrowLimit))

		// increase collateralization ratio
		r.NoError(r.Stabilization.SetMinCollateralizationRatio(r.Operator, new(big.Int).Add(cfg.MinCollateralizationRatio, e18)))

		// withdraw
		_, err = r.Stabilization.Withdraw(tests.FromSender(userAccount, nil), big.NewInt(1))
		require.ErrorAs(r.T, err, &tests.StabilizationInsufficientCollateralError{})
	})

	tests.RunWithSetup("Test withdraw withdraw with insufficient collateral (above current mcr) ", setup, func(r *tests.Runner) {
		cfg, _, err := r.Stabilization.Config(nil)
		require.NoError(t, err)

		cdp, _, err := r.Stabilization.Cdps(nil, userAccount)
		require.NoError(t, err)

		borrowLimit := calcBorrowLimit(r, userAccount)

		// borrow
		borrowAmount := new(big.Int).Div(borrowLimit, big.NewInt(2))
		collateralRequired, _, err := r.Stabilization.MinimumCollateral(
			nil,
			borrowAmount,
			newtonPrice,
			cfg.MinCollateralizationRatio,
		)
		require.NoError(t, err)

		withdrawMax := new(big.Int).Sub(cdp.Collateral, collateralRequired)
		r.NoError(r.Stabilization.Borrow(tests.FromSender(userAccount, nil), borrowAmount))

		_, err = r.Stabilization.Withdraw(tests.FromSender(userAccount, nil), new(big.Int).Add(withdrawMax, common.Big1))
		require.ErrorAs(r.T, err, &tests.StabilizationInsufficientCollateralError{})
	})
}

func TestStabilizationBorrow(t *testing.T) {
	userAccount := testrand.Address()
	fundedAmount := new(big.Int).Mul(e18, big.NewInt(100))
	totalDeposit := new(big.Int).Div(fundedAmount, big.NewInt(10))

	setup := func() *tests.Runner {
		r := tests.Setup(t, nil)
		setBasicConfig(r)
		primePrices(r, newtonPrice, toBase("0.97", 18))
		r.GiveMeSomeMoney(userAccount, new(big.Int).Mul(e18, big.NewInt(100)))
		_, err := r.Autonity.Mint(r.Operator, userAccount, fundedAmount)
		require.NoError(t, err)

		// approve stabilization
		_, err = r.Autonity.Approve(tests.FromSender(userAccount, nil), r.Stabilization.Address(), fundedAmount)
		require.NoError(t, err)

		// remove cdp restrictions for this teat
		_, err = r.Stabilization.RemoveCDPRestrictions(r.Operator)
		require.NoError(t, err)
		deposit(r, userAccount, totalDeposit)
		return r
	}

	tests.RunWithSetup("Test borrow zero not allowed", setup, func(r *tests.Runner) {
		_, err := r.Stabilization.Borrow(tests.FromSender(userAccount, nil), big.NewInt(0))
		require.ErrorAs(r.T, err, &tests.StabilizationInvalidAmountError{})
	})

	tests.RunWithSetup("Test borrow to limit", setup, func(r *tests.Runner) {
		// we need to make sure we are on a different block from the setup function, otherwise
		// our timestamp checks will fail
		r.WaitNBlocks(1)

		cdpBefore, _, err := r.Stabilization.Cdps(nil, userAccount)
		require.NoError(t, err)

		atnBefore := r.GetBalanceOf(userAccount)
		borrowLimit := calcBorrowLimit(r, userAccount)

		r.NoError(r.Stabilization.Borrow(tests.FromSender(userAccount, nil), borrowLimit))

		cdpAfter, _, err := r.Stabilization.Cdps(nil, userAccount)
		require.NoError(t, err)

		atnAfter := r.GetBalanceOf(userAccount)

		require.Greater(t, cdpAfter.Timestamp.Int64(), cdpBefore.Timestamp.Int64())
		require.Equal(t, borrowLimit, cdpAfter.Principal)
		require.Equal(t, uint64(0), cdpAfter.Interest.Uint64())
		require.Equal(t, new(big.Int).Add(atnBefore, borrowLimit), atnAfter)
		// ToDo(scott): verify borrow events
	})

	setupWithBorrow := func() *tests.Runner {
		r := setup()
		borrowLimit := calcBorrowLimit(r, userAccount)
		r.NoError(
			r.Stabilization.Borrow(
				tests.FromSender(userAccount, nil),
				new(big.Int).Div(borrowLimit, big.NewInt(2)),
			),
		)
		return r
	}

	tests.RunWithSetup("Test subsequent borrows", setupWithBorrow, func(r *tests.Runner) {
		r.WaitNBlocks(1) // add block between borrows
		borrowLimit := calcBorrowLimit(r, userAccount)
		borrowAmount := new(big.Int).Div(borrowLimit, big.NewInt(2))
		cdpsBefore, _, err := r.Stabilization.Cdps(nil, userAccount)
		require.NoError(t, err)
		pendingTimestamp := new(big.Int).Set(r.Evm.Context.Time)
		debt, _, err := r.Stabilization.DebtAmountAtTime(nil, userAccount, pendingTimestamp)
		require.NoError(t, err)
		interest := new(big.Int).Sub(debt, borrowAmount)
		amount := new(big.Int).Div(new(big.Int).Sub(borrowLimit, debt), big.NewInt(2))

		atnBefore := r.GetBalanceOf(userAccount)
		r.NoError(r.Stabilization.Borrow(tests.FromSender(userAccount, nil), amount))
		atnAfter := r.GetBalanceOf(userAccount)
		require.Equal(t, new(big.Int).Add(atnBefore, amount), atnAfter)

		cdpsAfter, _, err := r.Stabilization.Cdps(nil, userAccount)
		require.NoError(t, err)

		require.Greater(t, cdpsAfter.Timestamp.Int64(), cdpsBefore.Timestamp.Int64())
		require.Equal(t, new(big.Int).Add(borrowAmount, amount), cdpsAfter.Principal)
		require.Equal(t, interest, cdpsAfter.Interest)
	})

	tests.RunWithSetup("Test borrow minimum", setup, func(r *tests.Runner) {
		cfg, _, err := r.Stabilization.Config(nil)
		require.NoError(t, err)
		minDebt := cfg.MinDebtRequirement

		atnBefore := r.GetBalanceOf(userAccount)
		r.NoError(r.Stabilization.Borrow(tests.FromSender(userAccount, nil), minDebt))
		atnAfter := r.GetBalanceOf(userAccount)
		require.Equal(t, new(big.Int).Add(atnBefore, minDebt), atnAfter)

		cdp, _, err := r.Stabilization.Cdps(nil, userAccount)
		require.NoError(t, err)
		require.Equal(t, minDebt, cdp.Principal)

		// ToDo(scott): verify borrow events
	})

	tests.RunWithSetup("Test borrow below minimum fails", setup, func(r *tests.Runner) {
		cfg, _, err := r.Stabilization.Config(nil)
		require.NoError(t, err)
		minDebt := cfg.MinDebtRequirement

		_, err = r.Stabilization.Borrow(tests.FromSender(userAccount, nil), new(big.Int).Sub(minDebt, common.Big1))
		require.ErrorAs(r.T, err, &tests.StabilizationInvalidDebtPositionError{})
	})

	tests.RunWithSetup("Test borrow reverts if position is liquidatable", setup, func(r *tests.Runner) {
		borrowLimit := calcBorrowLimit(r, userAccount)
		cfg, _, err := r.Stabilization.Config(nil)
		require.NoError(t, err)
		r.NoError(
			r.Stabilization.Borrow(
				tests.FromSender(userAccount, nil),
				new(big.Int).Sub(borrowLimit, common.Big1),
			),
		)
		isLiquidatable, _, err := r.Stabilization.IsLiquidatable(nil, userAccount)
		require.NoError(t, err)
		require.False(t, isLiquidatable)

		mcr := cfg.MinCollateralizationRatio
		_, err = r.Stabilization.SetMinCollateralizationRatio(r.Operator, new(big.Int).Add(mcr, e18))
		require.NoError(t, err)

		_, err = r.Stabilization.SetLiquidationRatio(r.Operator, mcr)
		require.NoError(t, err)

		isLiquidatable, _, err = r.Stabilization.IsLiquidatable(nil, userAccount)
		require.NoError(t, err)
		require.True(t, isLiquidatable)

		_, err = r.Stabilization.Borrow(tests.FromSender(userAccount, nil), common.Big1)
		require.ErrorAs(r.T, err, &tests.StabilizationLiquidatableError{})
	})

	tests.RunWithSetup("Test cannot borrow over limit", setup, func(r *tests.Runner) {
		// set correct parameters to avoid liquidation error
		r.NoError(r.Stabilization.SetMinCollateralizationRatio(r.Operator, basicConfig.MinCollateralizationRatio))
		r.NoError(r.Stabilization.SetLiquidationRatio(r.Operator, basicConfig.LiquidationRatio))

		cdp, _, err := r.Stabilization.Cdps(nil, userAccount)
		require.NoError(t, err)

		borrowLimit, _, err := r.Stabilization.MaxBorrow(nil, cdp.Collateral)
		require.NoError(t, err)

		_, err = r.Stabilization.Borrow(tests.FromSender(userAccount, nil), new(big.Int).Add(borrowLimit, common.Big1))
		require.ErrorAs(r.T, err, &tests.StabilizationInsufficientCollateralError{})
	})
}

func TestStabilizationRepay(t *testing.T) {
	userAccount := testrand.Address()
	fundedAmount := new(big.Int).Mul(e18, big.NewInt(100))
	totalDeposit := new(big.Int).Div(fundedAmount, big.NewInt(10))

	setup := func() *tests.Runner {
		r := tests.Setup(t, nil)
		setBasicConfig(r)
		primePrices(r, newtonPrice, toBase("0.97", 18))
		r.GiveMeSomeMoney(userAccount, new(big.Int).Mul(e18, big.NewInt(100)))
		_, err := r.Autonity.Mint(r.Operator, userAccount, fundedAmount)
		require.NoError(t, err)

		// approve stabilization
		_, err = r.Autonity.Approve(tests.FromSender(userAccount, nil), r.Stabilization.Address(), fundedAmount)
		require.NoError(t, err)

		// remove cdp restrictions for this teat
		_, err = r.Stabilization.RemoveCDPRestrictions(r.Operator)
		require.NoError(t, err)
		deposit(r, userAccount, totalDeposit)

		// borrow
		r.NoError(r.Stabilization.Borrow(
			tests.FromSender(userAccount, nil),
			new(big.Int).Div(calcBorrowLimit(r, userAccount), big.NewInt(2)),
		))

		r.WaitNBlocks(1) // add block after borrow
		return r
	}

	tests.RunWithSetup("Test repay zero not allowed", setup, func(r *tests.Runner) {
		_, err := r.Stabilization.Repay(tests.FromSender(userAccount, common.Big0))
		require.ErrorAs(r.T, err, &tests.StabilizationZeroValueError{})
	})

	tests.RunWithSetup("Test repay invalid position", setup, func(r *tests.Runner) {
		timestamp := new(big.Int).Set(r.Evm.Context.Time)
		cfg, _, err := r.Stabilization.Config(nil)
		require.NoError(t, err)
		debtAmount, _, err := r.Stabilization.DebtAmountAtTime(nil, userAccount, timestamp)
		require.NoError(t, err)
		tooMuch := new(big.Int).Sub(new(big.Int).Add(debtAmount, common.Big1), cfg.MinDebtRequirement)
		_, err = r.Stabilization.Repay(tests.FromSender(userAccount, tooMuch))
		require.ErrorAs(r.T, err, &tests.StabilizationInvalidDebtPositionError{})
	})

	tests.RunWithSetup("Test repay to minimum debt", setup, func(r *tests.Runner) {
		timestamp := new(big.Int).Set(r.Evm.Context.Time)
		cfg, _, err := r.Stabilization.Config(nil)
		require.NoError(t, err)
		debtAmount, _, err := r.Stabilization.DebtAmountAtTime(nil, userAccount, timestamp)
		require.NoError(t, err)
		payment := new(big.Int).Sub(debtAmount, cfg.MinDebtRequirement)

		_, err = r.Stabilization.Repay(tests.FromSender(userAccount, payment))
		require.NoError(t, err)

		cdp, _, err := r.Stabilization.Cdps(nil, userAccount)
		require.NoError(t, err)
		require.Equal(t, cfg.MinDebtRequirement, cdp.Principal)
		require.Equal(t, uint64(0), cdp.Interest.Uint64())

		// ToDo(scott): verify repay events
	})

	tests.RunWithSetup("Test repay interest", setup, func(r *tests.Runner) {
		borrowAmount := new(big.Int).Div(calcBorrowLimit(r, userAccount), big.NewInt(2))
		timestamp := new(big.Int).Set(r.Evm.Context.Time)
		debtAmount, _, err := r.Stabilization.DebtAmountAtTime(nil, userAccount, timestamp)
		require.NoError(t, err)
		interest := new(big.Int).Sub(debtAmount, borrowAmount)

		r.NoError(r.Stabilization.Repay(tests.FromSender(userAccount, interest)))

		cdp, _, err := r.Stabilization.Cdps(nil, userAccount)
		require.NoError(t, err)
		require.Equal(t, borrowAmount, cdp.Principal)
		require.Equal(t, uint64(0), cdp.Interest.Uint64())

		// ToDo(scott): verify repay events
	})

	tests.RunWithSetup("Test repay full debt", setup, func(r *tests.Runner) {
		timestamp := new(big.Int).Set(r.Evm.Context.Time)
		debtAmount, _, err := r.Stabilization.DebtAmountAtTime(nil, userAccount, timestamp)
		require.NoError(t, err)

		r.NoError(r.Stabilization.Repay(tests.FromSender(userAccount, debtAmount)))

		cdp, _, err := r.Stabilization.Cdps(nil, userAccount)
		require.NoError(t, err)
		require.Equal(t, int64(0), cdp.Principal.Int64())
		require.Equal(t, uint64(0), cdp.Interest.Uint64())

		// ToDo(scott): verify repay events
	})

	tests.RunWithSetup("Test repay surplus is returned", setup, func(r *tests.Runner) {
		timestamp := new(big.Int).Set(r.Evm.Context.Time)
		debtAmount, _, err := r.Stabilization.DebtAmountAtTime(nil, userAccount, timestamp)
		require.NoError(t, err)
		surplus := new(big.Int).Set(common.Big256)

		balanceBefore := r.GetBalanceOf(userAccount)

		r.NoError(r.Stabilization.Repay(tests.FromSender(userAccount, new(big.Int).Add(debtAmount, surplus))))

		balanceAfter := r.GetBalanceOf(userAccount)

		cdp, _, err := r.Stabilization.Cdps(nil, userAccount)
		require.NoError(t, err)
		require.Equal(t, int64(0), cdp.Principal.Int64())
		require.Equal(t, uint64(0), cdp.Interest.Uint64())
		require.Equal(t, debtAmount, new(big.Int).Sub(balanceBefore, balanceAfter))

		// ToDo(scott): verify repay events
	})
}

func TestStabilizationCalculations(t *testing.T) {
	setup := func() *tests.Runner {
		return tests.Setup(t, nil)
	}

	tests.RunWithSetup("Test borrow limit", setup, func(r *tests.Runner) {
		// borrow limit = (collateral * price / mcr) * (target atn-acu * acu-usd / atn-usd)
		testCases := [][]*big.Int{
			{
				toBase("100", 18),                    // collateral
				toBase("1.2", 18),                    // price
				toBase("1.5", 18),                    // mcr
				toBase("1.23", 18),                   // target price atn-acu
				toBase("1.7", 18),                    // acu-usd
				toBase("0.6", 18),                    // atn-usd
				toBase("278.800000000000000000", 18), // expected
			},
			{
				toBase("100", 18),
				toBase("0.8", 18),
				toBase("1.5", 18),
				toBase("1.23", 18),
				toBase("1.7", 18),
				toBase("0.6", 18),
				toBase("185.866666666666666666", 18),
			},
			{
				toBase("100", 18),
				toBase("1.2", 18),
				toBase("1.2", 18),
				toBase("1.23", 18),
				toBase("1.7", 18),
				toBase("0.6", 18),
				toBase("348.5", 18),
			},
			{
				toBase("100", 18),
				toBase("0.8", 18),
				toBase("1.2", 18),
				toBase("1.23", 18),
				toBase("1.7", 18),
				toBase("0.6", 18),
				toBase("232.333333333333333333", 18),
			},
		}

		calculated := func(collateral, price, mcr, targetAtnACU, acuUSD, atnUSD *big.Int) *big.Int {
			num := newFloat0().Mul(
				newFloat0().Mul(newFloat(collateral), newFloat(price)),
				newFloat0().Mul(newFloat(targetAtnACU), newFloat(acuUSD)),
			)
			den := newFloat0().Mul(
				newFloat0().Mul(newFloat(mcr), newFloat(scaleFactor)),
				newFloat(atnUSD),
			)
			result, _ := newFloat0().Quo(num, den).Int(nil)
			return result
		}

		for _, tc := range testCases {
			collateral := tc[0]
			price := tc[1]
			mcr := tc[2]
			targetAtnACU := tc[3]
			acuUSD := tc[4]
			atnUSD := tc[5]

			expected := tc[6]

			actual, _, err := r.Stabilization.BorrowLimit(
				nil,
				collateral,   // collateral
				price,        // collateralPrice
				atnUSD,       // debtPrice
				targetAtnACU, // targetDebtPrice
				acuUSD,       // acuPrice
				mcr,          // mcr
			)
			require.NoError(t, err)
			require.Equal(t, expected, actual)

			// check that the calculated value is the same
			require.Equal(t, expected, calculated(collateral, price, mcr, targetAtnACU, acuUSD, atnUSD))
		}
	})

	tests.RunWithSetup("Test minimum collateral", setup, func(r *tests.Runner) {
		testCases := [][]*big.Int{
			{
				toBase("80", 18),
				toBase("1.2", 18),
				toBase("1.5", 18),
				toBase("100", 18),
			},
			{
				bigInt("53333333333333333333"),
				toBase("0.8", 18),
				toBase("1.5", 18),
				bigInt("99999999999999999999"),
			},
			{
				toBase("100", 18),
				toBase("1.2", 18),
				toBase("1.2", 18),
				toBase("100", 18),
			},
			{
				bigInt("66666666666666666666"),
				toBase("0.8", 18),
				toBase("1.2", 18),
				bigInt("99999999999999999999"),
			},
		}

		calculated := func(principal, mcr, price *big.Int) *big.Int {
			return new(big.Int).Div(new(big.Int).Mul(principal, mcr), price)
		}

		for _, tc := range testCases {
			principal := tc[0]
			price := tc[1]
			mcr := tc[2]
			expected := tc[3]

			actual, _, err := r.Stabilization.MinimumCollateral(nil, principal, price, mcr)
			require.NoError(t, err)
			require.Equal(t, expected, actual)

			// check that the calculated value is the same
			require.Equal(t, expected, calculated(principal, mcr, price))
		}
	})

	tests.RunWithSetup("Test interest due", setup, func(r *tests.Runner) {
		testCases := [][]*big.Int{
			{
				floatStringToBigInt("100", e18),
				floatStringToBigInt("0.05", e18),
				common.Big0,
				bigInt("2628000"),
				bigInt("417535929111852800"),
			},
			{
				floatStringToBigInt("100", e18),
				floatStringToBigInt("0.05", e18),
				bigInt("2628000"),
				bigInt("2628000"),
				common.Big0,
			},
		}

		for _, tc := range testCases {
			principal := tc[0]
			rate := tc[1]
			tStart := tc[2]
			tEnd := tc[3]
			expected := tc[4]

			actual, _, err := r.Stabilization.InterestDue(nil, principal, rate, tStart, tEnd)
			require.NoError(t, err)
			require.Equal(t, expected.String(), actual.String())

			// ToDo(scott): calculate expected interest due and compare with expected
			/*
				// from python
				t = quantize(SCALING_FACTOR * Decimal(tend - tstart) / SECONDS_IN_YEAR)
				rt = quantize(SCALING_FACTOR * rate / SCALING_FACTOR * t / SCALING_FACTOR)
				exp = quantize(SCALING_FACTOR * PRB_MATH_E ** (rt / SCALING_FACTOR))
				calculated = quantize(principal / SCALING_FACTOR * (exp - SCALING_FACTOR))
				assert result == calculated
			*/
		}
	})

	tests.RunWithSetup("Test collateral price has 18 decimals", setup, func(r *tests.Runner) {
		primeOracle(r, []string{"NTN-ATN"}, []*big.Int{newtonPrice})
		price, _, err := r.Stabilization.CollateralPrice(nil)
		require.NoError(t, err)
		require.Equal(t, price, newtonPrice)
	})
}

func TestStabilizationOnlyAutonityFunctions(t *testing.T) {
	setup := func() *tests.Runner {
		return tests.Setup(t, nil)
	}

	tests.RunWithSetup("Test only autonity functions", setup, func(r *tests.Runner) {
		r.NoError(r.Stabilization.SetOperator(tests.FromAutonity, testrand.Address()))
		r.NoError(r.Stabilization.SetOracle(tests.FromAutonity, testrand.Address()))

		unauthorizedUsers := []common.Address{
			params.TestAutonityContractConfig.Operator,
			params.DeployerAddress,
			testrand.Address(),
		}

		for _, user := range unauthorizedUsers {
			_, err := r.Stabilization.SetOperator(tests.FromSender(user, nil), testrand.Address())
			require.ErrorAs(t, err, &tests.StabilizationUnauthorizedError{})

			_, err = r.Stabilization.SetOracle(tests.FromSender(user, nil), testrand.Address())
			require.ErrorAs(t, err, &tests.StabilizationUnauthorizedError{})
		}
	})
}

// test helpers functions

func deposit(r *tests.Runner, userAccount common.Address, amount *big.Int) {
	stabilizationBalanceBefore, _, err := r.Autonity.BalanceOf(nil, r.Stabilization.Address())
	require.NoError(r.T, err)

	balanceBefore, _, err := r.Autonity.BalanceOf(nil, userAccount)
	require.NoError(r.T, err)
	_, err = r.Stabilization.Deposit(tests.FromSender(userAccount, common.Big0), amount)
	require.NoError(r.T, err)
	balanceAfter, _, err := r.Autonity.BalanceOf(nil, userAccount)
	require.NoError(r.T, err)
	require.Equal(r.T, new(big.Int).Sub(balanceBefore, amount), balanceAfter)

	stabilizationBalanceAfter, _, err := r.Autonity.BalanceOf(nil, r.Stabilization.Address())
	require.NoError(r.T, err)

	require.Equal(r.T, new(big.Int).Add(stabilizationBalanceBefore, amount), stabilizationBalanceAfter)
}

func calcBorrowLimit(r *tests.Runner, userAccount common.Address) *big.Int {
	cdp, _, err := r.Stabilization.Cdps(nil, userAccount)
	require.NoError(r.T, err)

	limit, _, err := r.Stabilization.MaxBorrow(nil, cdp.Collateral)
	require.NoError(r.T, err)
	return limit
}

func setBasicConfig(r *tests.Runner) {
	r.NoError(r.Stabilization.SetLiquidationRatio(r.Operator, basicConfig.LiquidationRatio))
	r.NoError(r.Stabilization.SetMinCollateralizationRatio(r.Operator, basicConfig.MinCollateralizationRatio))
	r.NoError(r.Stabilization.SetMinDebtRequirement(r.Operator, basicConfig.MinDebtRequirement))
}

func bigInt(s string) *big.Int {
	i, _ := new(big.Int).SetString(s, 10)
	return i
}

func primePrices(r *tests.Runner, ntnAtnPrice *big.Int, atnUsdPrice *big.Int) {
	oracleScale, _, err := r.Oracle.GetDecimals(nil)
	require.NoError(r.T, err)
	oracleDecimals := int64(oracleScale)

	acuScale, _, err := r.Acu.Scale(nil)
	require.NoError(r.T, err)
	acuDecimals := acuScale.Int64()

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
		toBase("0.6757", oracleDecimals),
		toBase("0.75694", oracleDecimals),
		toBase("1.1085", oracleDecimals),
		toBase("1.29403", oracleDecimals),
		toBase("0.00713", oracleDecimals),
		toBase("1.0", oracleDecimals),
		toBase("0.09597", oracleDecimals),
	}
	quantities := []*big.Int{
		toBase("0.213", acuDecimals),
		toBase("0.187", acuDecimals),
		toBase("0.143", acuDecimals),
		toBase("0.104", acuDecimals),
		toBase("17.6", acuDecimals),
		toBase("0.180", acuDecimals),
		toBase("1.41", acuDecimals),
	}
	primeACUBasket(r, symbols, quantities, big.NewInt(acuDecimals))
	primeOracle(r, append(symbols, "NTN-ATN", "ATN-USD"), append(prices, ntnAtnPrice, atnUsdPrice))
	_, err = r.Acu.Update(tests.FromAutonity)
	require.NoError(r.T, err)
}
