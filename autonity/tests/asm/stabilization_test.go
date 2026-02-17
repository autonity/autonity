package asm

import (
	"math/big"
	"testing"

	"github.com/ALTree/bigfloat"
	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/autonity/tests"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/internal/testrand"
	"github.com/autonity/autonity/params"
)

var (
	e18 = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
)

// newtonAutonPrice = 1234567 * 10^12 = 1.234567 * 10^18
var newtonAutonPrice = toBase("1.234567", 18)
var newtonUSDPrice = toBase("2.03626873", 18)
var basicConfig = tests.IStabilizationConfig{
	BorrowInterestRate:        toBase("0.5", 18),
	AnnouncementWindow:        big.NewInt(30),
	LiquidationRatio:          toBase("1.5", 18),
	MinCollateralizationRatio: toBase("2.5", 18),
	MinDebtRequirement:        new(big.Int).Exp(big.NewInt(10), big.NewInt(6), nil),
	TargetPrice:               toBase("1.0", 18),
	DefaultNTNATNPrice:        toBase("1.0", 18),
	DefaultNTNUSDPrice:        toBase("1.0", 18),
	DefaultACUUSDPrice:        big.NewInt(0),
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
				AnnouncementWindow:        basicConfig.AnnouncementWindow,
				LiquidationRatio:          basicConfig.LiquidationRatio,
				MinCollateralizationRatio: big.NewInt(0),
				MinDebtRequirement:        basicConfig.MinDebtRequirement,
				TargetPrice:               basicConfig.TargetPrice,
				DefaultNTNATNPrice:        basicConfig.DefaultNTNATNPrice,
				DefaultNTNUSDPrice:        basicConfig.DefaultNTNUSDPrice,
				DefaultACUUSDPrice:        basicConfig.DefaultACUUSDPrice,
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
				AnnouncementWindow:        basicConfig.AnnouncementWindow,
				LiquidationRatio:          e18,
				MinCollateralizationRatio: e18,
				MinDebtRequirement:        basicConfig.MinDebtRequirement,
				TargetPrice:               basicConfig.TargetPrice,
				DefaultNTNATNPrice:        basicConfig.DefaultNTNATNPrice,
				DefaultNTNUSDPrice:        basicConfig.DefaultNTNUSDPrice,
				DefaultACUUSDPrice:        basicConfig.DefaultACUUSDPrice,
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
				AnnouncementWindow:        basicConfig.AnnouncementWindow,
				LiquidationRatio:          new(big.Int).Add(basicConfig.MinCollateralizationRatio, big.NewInt(1)),
				MinCollateralizationRatio: basicConfig.MinCollateralizationRatio,
				MinDebtRequirement:        basicConfig.MinDebtRequirement,
				TargetPrice:               basicConfig.TargetPrice,
				DefaultNTNATNPrice:        basicConfig.DefaultNTNATNPrice,
				DefaultNTNUSDPrice:        basicConfig.DefaultNTNUSDPrice,
				DefaultACUUSDPrice:        basicConfig.DefaultACUUSDPrice,
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

	tests.RunWithSetup("Test constructor zero announcement window", setup, func(r *tests.Runner) {
		_, _, _, err := r.DeployStabilization(
			nil,
			tests.IStabilizationConfig{
				BorrowInterestRate:        basicConfig.BorrowInterestRate,
				AnnouncementWindow:        common.Big0,
				LiquidationRatio:          basicConfig.LiquidationRatio,
				MinCollateralizationRatio: basicConfig.MinCollateralizationRatio,
				MinDebtRequirement:        basicConfig.MinDebtRequirement,
				TargetPrice:               basicConfig.TargetPrice,
				DefaultNTNATNPrice:        basicConfig.DefaultNTNATNPrice,
				DefaultNTNUSDPrice:        basicConfig.DefaultNTNUSDPrice,
				DefaultACUUSDPrice:        basicConfig.DefaultACUUSDPrice,
			},
			common.Address{},
			common.Address{},
			common.Address{},
			common.Address{},
			common.Address{},
			common.Address{},
			common.Address{},
		)
		require.ErrorAs(r.T, err, &tests.StabilizationZeroValueError{})
	})
}

func TestStabilizationDeposit(t *testing.T) {
	userAccount := testrand.Address()
	secondUserAccount := testrand.Address()
	fundedAmount := new(big.Int).Mul(e18, big.NewInt(100))
	setup := func() *tests.Runner {
		r := tests.Setup(t, nil)
		r.NoError(r.Stabilization.UseFixedGenesisPrices(r.Operator, false))
		primePrices(r, newtonAutonPrice, newtonUSDPrice)
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
		r.NoError(r.Stabilization.UseFixedGenesisPrices(r.Operator, false))
		setBasicConfig(r)
		primePrices(r, newtonAutonPrice, newtonUSDPrice)
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

		liquidatable, _, err := r.Stabilization.IsLiquidatable(nil, userAccount)
		require.NoError(t, err)
		require.False(t, liquidatable)

		// make liquidatable
		r.NoError(r.Stabilization.UpdateRatios(r.Operator, cfg.MinCollateralizationRatio, new(big.Int).Add(cfg.MinCollateralizationRatio, e18)))
		progressTime(r, cfg.AnnouncementWindow.Int64())

		liquidatable, _, err = r.Stabilization.IsLiquidatable(nil, userAccount)
		require.NoError(t, err)
		require.True(t, liquidatable)

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
		r.NoError(r.Stabilization.UpdateRatios(r.Operator, cfg.LiquidationRatio, new(big.Int).Add(cfg.MinCollateralizationRatio, e18)))
		progressTime(r, cfg.AnnouncementWindow.Int64())

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

		newtonAcuPrice, _, err := r.Stabilization.CollateralPriceACU(nil)
		require.NoError(t, err)

		// borrow
		borrowAmount := new(big.Int).Div(borrowLimit, big.NewInt(2))
		collateralRequired, _, err := r.Stabilization.MinimumCollateral(
			nil,
			borrowAmount,
			newtonAcuPrice,
			cfg.TargetPrice,
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
		r.NoError(r.Stabilization.UseFixedGenesisPrices(r.Operator, false))
		setBasicConfig(r)
		primePrices(r, newtonAutonPrice, newtonUSDPrice)
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
		pendingTimestamp := new(big.Int).SetUint64(r.Evm.Context.Time)
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
				borrowLimit,
			),
		)
		isLiquidatable, _, err := r.Stabilization.IsLiquidatable(nil, userAccount)
		require.NoError(t, err)
		require.False(t, isLiquidatable)

		mcr := cfg.MinCollateralizationRatio

		//make liquidatable
		r.NoError(r.Stabilization.UpdateRatios(r.Operator, mcr, new(big.Int).Add(mcr, e18)))
		progressTime(r, cfg.AnnouncementWindow.Int64())

		cfg, _, err = r.Stabilization.Config(nil)
		require.NoError(t, err)
		require.Equal(t, mcr, cfg.LiquidationRatio)

		isLiquidatable, _, err = r.Stabilization.IsLiquidatable(nil, userAccount)
		require.NoError(t, err)
		require.True(t, isLiquidatable)

		_, err = r.Stabilization.Borrow(tests.FromSender(userAccount, nil), common.Big1)
		require.ErrorAs(r.T, err, &tests.StabilizationLiquidatableError{})
	})

	tests.RunWithSetup("Test cannot borrow over limit", setup, func(r *tests.Runner) {
		announceWindow, _, err := r.Stabilization.GetAnnouncementWindow(nil)
		require.NoError(t, err)

		// set correct parameters to avoid liquidation error
		r.NoError(r.Stabilization.UpdateRatios(r.Operator, basicConfig.LiquidationRatio, basicConfig.MinCollateralizationRatio))
		progressTime(r, announceWindow.Int64())

		cdp, _, err := r.Stabilization.Cdps(nil, userAccount)
		require.NoError(t, err)

		borrowLimit, _, err := r.Stabilization.MaxBorrow(nil, cdp.Collateral)
		require.NoError(t, err)

		_, err = r.Stabilization.Borrow(
			tests.FromSender(userAccount, nil),
			new(big.Int).Add(borrowLimit, common.Big1),
		)
		require.ErrorAs(r.T, err, &tests.StabilizationInsufficientCollateralError{})
	})
}

func TestStabilizationRepay(t *testing.T) {
	userAccount := testrand.Address()
	fundedAmount := new(big.Int).Mul(e18, big.NewInt(100))
	totalDeposit := new(big.Int).Div(fundedAmount, big.NewInt(10))

	setup := func() *tests.Runner {
		r := tests.Setup(t, nil)
		r.NoError(r.Stabilization.UseFixedGenesisPrices(r.Operator, false))
		setBasicConfig(r)
		primePrices(r, newtonAutonPrice, toBase("0.97", 18))
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

		borrowLimit, _, err := r.Stabilization.MaxBorrow(nil, totalDeposit)
		require.NoError(t, err)

		// borrow
		r.NoError(r.Stabilization.Borrow(
			tests.FromSender(userAccount, nil),
			new(big.Int).Div(borrowLimit, big.NewInt(2)),
		))

		r.WaitNBlocks(1) // add block after borrow
		return r
	}

	tests.RunWithSetup("Test repay zero not allowed", setup, func(r *tests.Runner) {
		_, err := r.Stabilization.Repay(tests.FromSender(userAccount, common.Big0))
		require.ErrorAs(r.T, err, &tests.StabilizationZeroValueError{})
	})

	tests.RunWithSetup("Test repay to below minimum debt requirement", setup, func(r *tests.Runner) {
		timestamp := new(big.Int).SetUint64(r.Evm.Context.Time)
		cfg, _, err := r.Stabilization.Config(nil)
		require.NoError(t, err)
		debtAmount, _, err := r.Stabilization.DebtAmountAtTime(nil, userAccount, timestamp)
		require.NoError(t, err)
		tooMuch := new(big.Int).Sub(new(big.Int).Add(debtAmount, common.Big1), cfg.MinDebtRequirement)
		_, err = r.Stabilization.Repay(tests.FromSender(userAccount, tooMuch))
		require.ErrorAs(r.T, err, &tests.StabilizationInvalidDebtPositionError{})
	})

	tests.RunWithSetup("Test repay to minimum debt", setup, func(r *tests.Runner) {
		timestamp := new(big.Int).SetUint64(r.Evm.Context.Time)
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
		timestamp := new(big.Int).SetUint64(r.Evm.Context.Time)
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
		timestamp := new(big.Int).SetUint64(r.Evm.Context.Time)
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
		timestamp := new(big.Int).SetUint64(r.Evm.Context.Time)
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

	tests.RunWithSetup("Test repay interest forwarded to auctioneer", setup, func(r *tests.Runner) {
		debtAmount, _, err := r.Stabilization.DebtAmount(nil, userAccount)
		require.NoError(t, err)

		cdp, _, err := r.Stabilization.Cdps(nil, userAccount)
		require.NoError(t, err)

		interest := new(big.Int).Sub(debtAmount, cdp.Principal)
		require.True(t, interest.Cmp(common.Big0) > 0)

		auctioneerBalanceBefore := r.GetBalanceOf(r.Auctioneer.Address())

		r.NoError(r.Stabilization.Repay(tests.FromSender(userAccount, interest)))

		auctioneerBalanceAfter := r.GetBalanceOf(r.Auctioneer.Address())

		require.Equal(t, interest, new(big.Int).Sub(auctioneerBalanceAfter, auctioneerBalanceBefore))
	})

}

func TestStabilizationLiquidate(t *testing.T) {
	userAccount := testrand.Address()
	setup := func() *tests.Runner {
		r := tests.Setup(t, nil)
		r.NoError(r.Stabilization.UseFixedGenesisPrices(r.Operator, false))
		setBasicConfig(r)
		primePrices(r, newtonAutonPrice, newtonUSDPrice)

		fundedAmount := new(big.Int).Mul(e18, big.NewInt(100))
		totalDeposit := new(big.Int).Div(fundedAmount, big.NewInt(10))

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

		borrowLimit := calcBorrowLimit(r, userAccount)
		r.NoError(r.Stabilization.Borrow(tests.FromSender(userAccount, nil), borrowLimit))

		return r
	}

	tests.RunWithSetup("Test liquidate can only be called by auctioneer", setup, func(r *tests.Runner) {
		_, err := r.Stabilization.Liquidate(
			tests.FromSender(testrand.Address(), nil),
			testrand.Address(),
			common.Big1,
			testrand.Address(),
		)
		require.ErrorAs(r.T, err, &tests.StabilizationUnauthorizedError{})
	})

	tests.RunWithSetup("Test liquidate returns overpayment to bidder", setup, func(r *tests.Runner) {
		makeLiquidatable(r, userAccount)

		debtAmount, _, err := r.Stabilization.DebtAmount(nil, userAccount)
		require.NoError(t, err)

		liquidator := testrand.Address()
		overpay := new(big.Int).Mul(debtAmount, big.NewInt(2))

		r.GiveMeSomeMoney(r.Auctioneer.Address(), overpay)

		liquidatorBalanceBefore := r.GetBalanceOf(liquidator)

		_, err = r.Stabilization.Liquidate(
			tests.FromSender(r.Auctioneer.Address(), overpay),
			userAccount,
			big.NewInt(1000),
			liquidator,
		)
		require.NoError(t, err)

		liquidatorBalanceAfter := r.GetBalanceOf(liquidator)

		require.Equal(t, new(big.Int).Sub(overpay, debtAmount), new(big.Int).Sub(liquidatorBalanceAfter, liquidatorBalanceBefore))
	})

	tests.RunWithSetup("Test liquidate reduces debt by full amount", setup, func(r *tests.Runner) {
		makeLiquidatable(r, userAccount)

		cdpBefore, _, err := r.Stabilization.Cdps(nil, userAccount)
		require.NoError(t, err)
		collateralBefore := cdpBefore.Collateral

		liquidatable, _, err := r.Stabilization.IsLiquidatable(nil, userAccount)
		require.NoError(t, err)
		require.True(t, liquidatable)

		debtAmount, _, err := r.Stabilization.DebtAmount(nil, userAccount)
		require.NoError(t, err)

		liquidator := testrand.Address()
		r.GiveMeSomeMoney(r.Auctioneer.Address(), debtAmount)

		// not liquidating full collateral
		collateralAmount := big.NewInt(1000)

		r.NoError(r.Stabilization.Liquidate(
			tests.FromSender(r.Auctioneer.Address(), debtAmount),
			userAccount,
			collateralAmount,
			liquidator,
		))

		cdp, _, err := r.Stabilization.Cdps(nil, userAccount)
		collateralAfter := cdp.Collateral

		require.NoError(t, err)
		require.Equal(t, int64(0), cdp.Principal.Int64())
		require.Equal(t, uint64(0), cdp.Interest.Uint64())
		require.Equal(t, new(big.Int).Sub(collateralBefore, collateralAmount), collateralAfter)
	})

	tests.RunWithSetup("Test liquidate cannot be called on a non-liquidatable cdp", setup, func(r *tests.Runner) {
		liquidatable, _, err := r.Stabilization.IsLiquidatable(nil, userAccount)
		require.NoError(t, err)
		require.False(t, liquidatable)

		debtAmount, _, err := r.Stabilization.DebtAmount(nil, userAccount)
		require.NoError(t, err)
		r.GiveMeSomeMoney(r.Auctioneer.Address(), debtAmount)

		_, err = r.Stabilization.Liquidate(
			tests.FromSender(r.Auctioneer.Address(), debtAmount),
			userAccount,
			common.Big1,
			testrand.Address(),
		)
		require.ErrorAs(r.T, err, &tests.StabilizationNotLiquidatableError{})
	})
}

func TestStabilizationCalculations(t *testing.T) {
	setup := func() *tests.Runner {
		r := tests.Setup(t, nil)
		r.NoError(r.Stabilization.UseFixedGenesisPrices(r.Operator, false))
		return r
	}

	tests.RunWithSetup("Test borrow limit", setup, func(r *tests.Runner) {
		// borrow limit = (collateral * price / mcr * target atn-acu)
		testCases := [][]*big.Int{
			{
				toBase("100", 18),                   // collateral
				toBase("1.2", 18),                   // price ntn-acu
				toBase("1.5", 18),                   // mcr
				toBase("1.20", 18),                  // target price atn-acu
				toBase("66.666666666666666666", 18), // expected
			},
			{
				toBase("100", 18),
				toBase("0.8", 18),
				toBase("1.5", 18),
				toBase("1.20", 18),
				toBase("44.444444444444444444", 18),
			},
			{
				toBase("100", 18),
				toBase("1.2", 18),
				toBase("1.2", 18),
				toBase("1.47", 18),
				toBase("68.027210884353741496", 18),
			},
			{
				toBase("120", 18),
				toBase("0.8", 18),
				toBase("1.5", 18),
				toBase("1.23", 18),
				toBase("52.032520325203252032", 18),
			},
		}

		calculated := func(collateral, priceNtnACU, mcr, targetPrice *big.Int) *big.Int {
			num := newFloat0().Mul(
				newFloat0().Mul(newFloat(collateral), newFloat(priceNtnACU)),
				newFloat(e18),
			)
			den := newFloat0().Mul(
				newFloat(mcr),
				newFloat(targetPrice),
			)
			result, _ := newFloat0().Quo(num, den).Int(nil)
			return result
		}

		for _, tc := range testCases {
			collateral := tc[0]
			priceNtnACU := tc[1]
			mcr := tc[2]
			targetAtnACU := tc[3]

			expected := tc[4]

			actual, _, err := r.Stabilization.BorrowLimit(
				nil,
				collateral,   // collateral
				priceNtnACU,  // collateralPrice
				targetAtnACU, // targetDebtPrice
				mcr,          // mcr
			)
			require.NoError(t, err)
			require.Equal(t, expected, actual)

			// check that the calculated value is the same
			require.Equal(t, expected, calculated(collateral, priceNtnACU, mcr, targetAtnACU))
		}
	})

	tests.RunWithSetup("Test borrow limit overflow", setup, func(r *tests.Runner) {
		ntnMaxSupply := new(big.Int).Mul(
			big.NewInt(100_000_000),
			e18,
		) // 100 million
		ntnACU := toBase("0.0001", 18)
		_, _, err := r.Stabilization.BorrowLimit(
			nil,
			ntnMaxSupply,       // collateral
			ntnACU,             // price ntn-acu
			toBase("1.23", 18), // target price atn-acu
			toBase("1.5", 18),  // mcr
		)
		require.NoError(r.T, err)
	})

	tests.RunWithSetup("Test minimum collateral", setup, func(r *tests.Runner) {
		// minimum collateral = (principal * mcr * target price) / price
		testCases := [][]*big.Int{
			{
				toBase("80", 18),  // collateral
				toBase("1.2", 18), // price ntn-acu
				toBase("1.0", 18), // target price atn-acu
				toBase("2.5", 18), // mcr
				toBase("166.666666666666666666", 18),
			},
			{
				toBase("53.333333333333333333", 18),
				toBase("0.8", 18),
				toBase("1.0", 18),
				toBase("1.5", 18),
				toBase("99.999999999999999999", 18),
			},
			{
				toBase("100", 18),
				toBase("1.2", 18),
				toBase("1.618", 18),
				toBase("1.5", 18),
				toBase("202.25", 18),
			},
			{
				toBase("66.666666666666666666", 18),
				toBase("0.8", 18),
				toBase("2.0", 18),
				toBase("1.2", 18),
				toBase("199.999999999999999998", 18),
			},
		}

		calculated := func(principal, mcr, price, targetPrice *big.Int) *big.Int {
			return new(big.Int).Div(
				new(big.Int).Mul(principal, new(big.Int).Mul(mcr, targetPrice)),
				new(big.Int).Mul(price, e18),
			)
		}

		for _, tc := range testCases {
			principal := tc[0]
			price := tc[1]
			targetPrice := tc[2]
			mcr := tc[3]
			expected := tc[4]

			actual, _, err := r.Stabilization.MinimumCollateral(nil, principal, price, targetPrice, mcr)
			require.NoError(t, err)
			require.Equal(t, expected, actual)

			// check that the calculated value is the same
			require.Equal(t, expected, calculated(principal, mcr, price, targetPrice))
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

			rateExponent, _, err := r.Stabilization.InterestExponent(nil, rate, tStart, tEnd)
			require.NoError(t, err)
			actual, _, err := r.Stabilization.InterestDue(nil, principal, rateExponent)
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
		primeOracle(r, []string{"NTN-ATN"}, []*big.Int{newtonAutonPrice})
		price, _, err := r.Stabilization.CollateralPrice(nil)
		require.NoError(t, err)
		require.Equal(t, price, newtonAutonPrice)
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

func TestInterestCalculation(t *testing.T) {
	user := tests.User
	depositAmmount := new(big.Int).Mul(
		big.NewInt(1000_000),
		e18,
	)

	borrowAmount := new(big.Int).Div(
		e18,
		big.NewInt(100),
	)
	setup := func() *tests.Runner {
		r := tests.Setup(t, nil)
		r.NoError(
			r.Stabilization.SetAtnSupplyOperator(
				r.Operator,
				user,
			),
		)
		setBasicConfig(r)
		r.NoError(
			r.Stabilization.SetMinDebtRequirement(
				r.Operator,
				common.Big0,
			),
		)
		primePrices(r, newtonAutonPrice, toBase("0.97", 18))
		r.GiveMeSomeMoney(user, new(big.Int).Mul(e18, big.NewInt(100)))
		_, err := r.Autonity.Mint(r.Operator, user, depositAmmount)
		require.NoError(t, err)

		// approve stabilization
		_, err = r.Autonity.Approve(tests.FromSender(user, nil), r.Stabilization.Address(), depositAmmount)
		require.NoError(t, err)

		deposit(r, user, depositAmmount)
		borrow(r, user, borrowAmount)
		return r
	}

	year := params.SecondsInYear

	tests.RunWithSetup("0 rate", setup, func(r *tests.Runner) {
		progressTime(r, year)
		require.Equal(r.T, borrowAmount, getDebt(r, user))
	})

	tests.RunWithSetup("single positive rate", setup, func(r *tests.Runner) {
		rateActive := new(big.Int).SetUint64(r.Evm.Context.Time)
		r.NoError(
			r.Stabilization.RemoveCDPRestrictions(r.Operator),
		)
		currentRate := getCurrentRate(r)
		require.True(r.T, currentRate.Cmp(common.Big0) > 0)
		progressTime(r, year)
		debt := calculateDebt(r, borrowAmount, common.Big0, []interestRateParam{
			{
				rate:      currentRate,
				startTime: rateActive,
				endTime:   new(big.Int).SetUint64(r.Evm.Context.Time),
			},
		})
		require.Equal(r.T, debt, getDebt(r, user))
	})

	tests.RunWithSetup("multiple rate", setup, func(r *tests.Runner) {
		rateActive := new(big.Int).SetUint64(r.Evm.Context.Time)
		r.NoError(
			r.Stabilization.RemoveCDPRestrictions(r.Operator),
		)
		firstRate := getCurrentRate(r)
		progressTime(r, year)

		rates := []interestRateParam{
			{
				rate:    common.Big0,
				endTime: big.NewInt(2 * year),
			},
			{
				rate: new(big.Int).Div(
					new(big.Int).Mul(e18, big.NewInt(5)),
					big.NewInt(100),
				), // 5%
				endTime: big.NewInt(year / 2),
			},
			{
				rate: new(big.Int).Div(
					new(big.Int).Mul(e18, big.NewInt(5)),
					big.NewInt(1000),
				), // 0.5%
				endTime: big.NewInt(5 * year),
			},
			{
				rate: new(big.Int).Div(
					new(big.Int).Mul(e18, big.NewInt(1)),
					big.NewInt(100),
				), // 1%
				endTime: big.NewInt(3 * year),
			},
		}

		// apply the rates
		window := getAnnouncementWindow(r)
		for i, rate := range rates {
			rates[i].startTime = new(big.Int).Add(window, new(big.Int).SetUint64(r.Evm.Context.Time))
			if i > 0 {
				rates[i-1].endTime = rates[i].startTime
			}
			r.NoError(
				r.Stabilization.UpdateBorrowInterestRate(
					r.Operator,
					rate.rate,
				),
			)
			progressTime(r, window.Int64())
			progressTime(r, rate.endTime.Int64())
		}
		rates[len(rates)-1].endTime = new(big.Int).SetUint64(r.Evm.Context.Time)
		// include the first one
		rates = append(
			[]interestRateParam{
				{
					rate:      firstRate,
					startTime: rateActive,
					endTime:   rates[0].startTime,
				},
			},
			rates...,
		)

		debt := calculateDebt(r, borrowAmount, common.Big0, rates)
		require.Equal(r.T, debt, getDebt(r, user))
	})

	tests.RunWithSetup("single rate with fractional window", setup, func(r *tests.Runner) {
		rateActive := new(big.Int).SetUint64(r.Evm.Context.Time)
		r.NoError(
			r.Stabilization.RemoveCDPRestrictions(r.Operator),
		)
		progressTime(r, 2*year)
		rate := getCurrentRate(r)
		debt := calculateDebt(r, borrowAmount, common.Big0, []interestRateParam{{
			rate:      rate,
			startTime: rateActive,
			endTime:   new(big.Int).SetUint64(r.Evm.Context.Time),
		}})

		rateActive = new(big.Int).SetUint64(r.Evm.Context.Time)
		borrow(r, user, borrowAmount)
		debt = new(big.Int).Add(debt, borrowAmount)
		cdp := getCdp(r, user)
		require.Equal(r.T, debt, new(big.Int).Add(cdp.Interest, cdp.Principal))

		progressTime(r, 3*year)
		debt = calculateDebt(r, debt, common.Big0, []interestRateParam{{
			rate:      rate,
			startTime: rateActive,
			endTime:   new(big.Int).SetUint64(r.Evm.Context.Time),
		}})
		require.Equal(r.T, debt, getDebt(r, user))
	})

	tests.RunWithSetup("multiple rates with fractional window", setup, func(r *tests.Runner) {
		rateActive := new(big.Int).SetUint64(r.Evm.Context.Time)
		r.NoError(
			r.Stabilization.RemoveCDPRestrictions(r.Operator),
		)
		rates := []interestRateParam{
			{
				rate:      getCurrentRate(r),
				startTime: rateActive,
				endTime:   big.NewInt(year), // duration for now
			},
			{
				rate:    common.Big0,
				endTime: big.NewInt(2 * year),
			},
			{
				rate: new(big.Int).Div(
					new(big.Int).Mul(e18, big.NewInt(5)),
					big.NewInt(100),
				), // 5%
				endTime: big.NewInt(year / 2),
			},
			{
				rate: new(big.Int).Div(
					new(big.Int).Mul(e18, big.NewInt(5)),
					big.NewInt(1000),
				), // 0.5%
				endTime: big.NewInt(5 * year),
			},
			{
				rate: new(big.Int).Div(
					new(big.Int).Mul(e18, big.NewInt(1)),
					big.NewInt(100),
				), // 1%
				endTime: big.NewInt(3 * year),
			},
		}

		window := getAnnouncementWindow(r)
		debt := borrowAmount
		debtUpdateTime := common.Big0
		for i, rate := range rates {
			rateDuration := rate.endTime
			if i > 0 {
				rates[i].startTime = new(big.Int).Add(new(big.Int).SetUint64(r.Evm.Context.Time), window)
				rates[i-1].endTime = rates[i].startTime
				r.NoError(
					r.Stabilization.UpdateBorrowInterestRate(
						r.Operator,
						rate.rate,
					),
				)
				// let it be active
				progressTime(r, window.Int64())
			}

			progressTime(r, rateDuration.Int64())
			rates[i].endTime = new(big.Int).SetUint64(r.Evm.Context.Time)
			debt = calculateDebt(r, debt, debtUpdateTime, rates[:i+1])

			debtUpdateTime = new(big.Int).SetUint64(r.Evm.Context.Time)
			borrow(r, user, borrowAmount)
			debt = new(big.Int).Add(debt, borrowAmount)
			cdp := getCdp(r, user)
			require.Equal(r.T, debt, new(big.Int).Add(cdp.Principal, cdp.Interest))
			progressTime(r, rateDuration.Int64()-window.Int64())
		}

		rates[len(rates)-1].endTime = new(big.Int).SetUint64(r.Evm.Context.Time)
		require.Equal(
			r.T,
			calculateDebt(r, debt, debtUpdateTime, rates),
			getDebt(r, user),
		)
	})

	tests.RunWithSetup("deposit followed by borrow", setup, func(r *tests.Runner) {
		rateActive := new(big.Int).SetUint64(r.Evm.Context.Time)
		r.NoError(
			r.Stabilization.RemoveCDPRestrictions(r.Operator),
		)
		currentRate := getCurrentRate(r)
		require.True(r.T, currentRate.Cmp(common.Big0) > 0)
		progressTime(r, year)
		debt := calculateDebt(r, borrowAmount, common.Big0, []interestRateParam{
			{
				rate:      currentRate,
				startTime: rateActive,
				endTime:   new(big.Int).SetUint64(r.Evm.Context.Time),
			},
		})
		require.Equal(r.T, debt, getDebt(r, user))

		r.NoError(
			r.Autonity.Approve(tests.FromSender(user, nil), r.Stabilization.Address(), depositAmmount),
		)
		r.NoError(
			r.Autonity.Mint(r.Operator, user, depositAmmount),
		)

		deposit(r, user, depositAmmount)
		borrow(r, user, borrowAmount)
		require.Equal(r.T, new(big.Int).Add(debt, borrowAmount), getDebt(r, user))

	})
}

func TestUpdateBorrowInterestRate(t *testing.T) {
	setup := func() *tests.Runner {
		return tests.Setup(t, nil)
	}

	tests.RunWithSetup("cannot update borrow interest rate until cdp restriction removed", setup, func(r *tests.Runner) {
		_, err := r.Stabilization.UpdateBorrowInterestRate(
			r.Operator,
			common.Big2,
		)
		require.ErrorAs(r.T, err, &tests.StabilizationUnauthorizedError{})
	})

	newSetup := func() *tests.Runner {
		r := tests.Setup(t, nil)
		r.NoError(
			r.Stabilization.RemoveCDPRestrictions(r.Operator),
		)
		return r
	}

	tests.RunWithSetup("onply operator can update borrow interest rate", newSetup, func(r *tests.Runner) {
		_, err := r.Stabilization.UpdateBorrowInterestRate(
			nil,
			common.Big1,
		)
		require.ErrorAs(r.T, err, &tests.StabilizationUnauthorizedError{})

		_, err = r.Stabilization.UpdateBorrowInterestRate(
			tests.FromSender(tests.User, nil),
			common.Big1,
		)
		require.ErrorAs(r.T, err, &tests.StabilizationUnauthorizedError{})
	})

	tests.RunWithSetup("updated interest rate activates after window time", newSetup, func(r *tests.Runner) {
		window := getAnnouncementWindow(r)
		realTime := new(big.Int).Add(
			window,
			new(big.Int).SetUint64(r.Evm.Context.Time),
		)
		currentRate := getCurrentRate(r)
		newRate := common.Big2
		require.NotEqual(r.T, newRate, currentRate, "cannot test")
		r.NoError(
			r.Stabilization.UpdateBorrowInterestRate(
				r.Operator,
				newRate,
			),
		)
		require.Equal(r.T, currentRate, getCurrentRate(r))
		pendingRate, activeSince, _, err := r.Stabilization.GetPendingInterestRateInfo(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, newRate, pendingRate)
		require.Equal(r.T, realTime, activeSince)
		progressTime(r, new(big.Int).Sub(realTime, new(big.Int).SetUint64(r.Evm.Context.Time)).Int64())
		require.Equal(r.T, newRate, getCurrentRate(r))
	})

	tests.RunWithSetup("pending rate is overridden by new rate", newSetup, func(r *tests.Runner) {
		window := getAnnouncementWindow(r)
		currentRate := getCurrentRate(r)
		newRate := common.Big2
		require.NotEqual(r.T, newRate, currentRate, "cannot test")
		r.NoError(
			r.Stabilization.UpdateBorrowInterestRate(
				r.Operator,
				newRate,
			),
		)

		_, activeSince, _, err := r.Stabilization.GetPendingInterestRateInfo(nil)
		require.NoError(r.T, err)
		progressTime(r, new(big.Int).Sub(activeSince, new(big.Int).SetUint64(r.Evm.Context.Time)).Int64()-1)
		// not updated yet
		require.Equal(r.T, currentRate, getCurrentRate(r))
		newRate2 := common.Big1
		require.NotEqual(r.T, newRate2, currentRate, "cannot test")
		realTime := new(big.Int).Add(window, new(big.Int).SetUint64(r.Evm.Context.Time))
		r.NoError(
			r.Stabilization.UpdateBorrowInterestRate(
				r.Operator,
				newRate2,
			),
		)
		pendingRate, activeSince, _, err := r.Stabilization.GetPendingInterestRateInfo(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, newRate2, pendingRate)
		require.Equal(r.T, realTime, activeSince)
	})
}

func TestUpdateAnnouncementWindow(t *testing.T) {
	setup := func() *tests.Runner {
		return tests.Setup(t, nil)
	}

	tests.RunWithSetup("only operator can update announcement window", setup, func(r *tests.Runner) {
		_, err := r.Stabilization.UpdateAnnouncementWindow(nil, common.Big1)
		require.ErrorAs(r.T, err, &tests.StabilizationUnauthorizedError{})

		_, err = r.Stabilization.UpdateAnnouncementWindow(
			tests.FromSender(tests.User, nil),
			common.Big1,
		)
		require.ErrorAs(r.T, err, &tests.StabilizationUnauthorizedError{})
	})

	tests.RunWithSetup("window cannot be zero", setup, func(r *tests.Runner) {
		_, err := r.Stabilization.UpdateAnnouncementWindow(r.Operator, common.Big0)
		require.ErrorAs(r.T, err, &tests.StabilizationZeroValueError{})
	})

	tests.RunWithSetup("pending window takes affect after current window", setup, func(r *tests.Runner) {
		testWindowUpdate := func(newWindow *big.Int) {
			currentWindow := getAnnouncementWindow(r)
			realTime := new(big.Int).Add(new(big.Int).SetUint64(r.Evm.Context.Time), currentWindow)
			r.NoError(
				r.Stabilization.UpdateAnnouncementWindow(
					r.Operator,
					newWindow,
				),
			)
			pendingAnnouncementWindow, activeSince, _, err := r.Stabilization.GetPendingAnnouncementWindowInfo(nil)
			require.NoError(r.T, err)
			require.Equal(r.T, newWindow, pendingAnnouncementWindow)
			require.Equal(r.T, realTime, activeSince)
			progressTime(r, new(big.Int).Sub(realTime, new(big.Int).SetUint64(r.Evm.Context.Time)).Int64()-1)
			require.Equal(r.T, currentWindow, getAnnouncementWindow(r))
			progressTime(r, 1)
			require.Equal(r.T, newWindow, getAnnouncementWindow(r))
		}

		testWindowUpdate(new(big.Int).Add(getAnnouncementWindow(r), common.Big1))
		testWindowUpdate(new(big.Int).Sub(getAnnouncementWindow(r), common.Big1))
	})

	tests.RunWithSetup("cannot update window while there is a pending one", setup, func(r *tests.Runner) {
		currentWindow := getAnnouncementWindow(r)
		activeSince := new(big.Int).Add(new(big.Int).SetUint64(r.Evm.Context.Time), currentWindow)
		r.NoError(
			r.Stabilization.UpdateAnnouncementWindow(
				r.Operator,
				currentWindow,
			),
		)
		// no update is allowed
		_, err := r.Stabilization.UpdateAnnouncementWindow(
			r.Operator,
			currentWindow,
		)
		require.ErrorAs(r.T, err, &tests.StabilizationAnnouncementWindowPendingError{})
		progressTime(r, new(big.Int).Sub(activeSince, new(big.Int).SetUint64(r.Evm.Context.Time)).Int64()-1)
		_, err = r.Stabilization.UpdateAnnouncementWindow(
			r.Operator,
			currentWindow,
		)
		require.ErrorAs(r.T, err, &tests.StabilizationAnnouncementWindowPendingError{})

		progressTime(r, 1)
		// update is allowed
		r.NoError(
			r.Stabilization.UpdateAnnouncementWindow(
				r.Operator,
				currentWindow,
			),
		)
	})
}

func TestUpdateRatios(t *testing.T) {
	setup := func() *tests.Runner {
		return tests.Setup(t, nil)
	}
	tests.RunWithSetup("only operator can update ratios", setup, func(r *tests.Runner) {
		_, err := r.Stabilization.UpdateRatios(nil, toBase("1.0", 18), toBase("1.5", 18))
		require.ErrorAs(r.T, err, &tests.StabilizationUnauthorizedError{})

		_, err = r.Stabilization.UpdateRatios(
			tests.FromSender(tests.User, nil),
			toBase("1.0", 18),
			toBase("1.5", 18),
		)
		require.ErrorAs(r.T, err, &tests.StabilizationUnauthorizedError{})
	})

	tests.RunWithSetup("ratios cannot be zero", setup, func(r *tests.Runner) {
		_, err := r.Stabilization.UpdateRatios(r.Operator, common.Big0, common.Big1)
		require.ErrorAs(r.T, err, &tests.StabilizationInvalidParameterError{})
	})

	tests.RunWithSetup("liquidation ratio cannot be less than SCALE_FACTOR", setup, func(r *tests.Runner) {
		// since lr < mcr this also implies mcr cannot be less than SCALE_FACTOR
		_, err := r.Stabilization.UpdateRatios(r.Operator, toBase("0.99", 18), toBase("1.5", 18))
		require.ErrorAs(r.T, err, &tests.StabilizationInvalidParameterError{})
	})

	tests.RunWithSetup("liquidation ratios can be updated", setup, func(r *tests.Runner) {
		newRatio := toBase("1.75", 18)
		cfg, _, err := r.Stabilization.Config(nil)
		require.NoError(r.T, err)
		require.NotEqual(r.T, newRatio, cfg.LiquidationRatio)

		r.NoError(
			r.Stabilization.UpdateRatios(
				r.Operator,
				newRatio,
				new(big.Int).Add(newRatio, toBase("1", 18)),
			),
		)

		progressTime(r, cfg.AnnouncementWindow.Int64())

		readRatio, _, err := r.Stabilization.LiquidationRatio(nil)
		require.NoError(t, err)
		cfg, _, err = r.Stabilization.Config(nil)
		require.NoError(r.T, err)
		require.Equal(t, cfg.LiquidationRatio, readRatio)
		require.Equal(t, newRatio, readRatio)

		require.Equal(t, new(big.Int).Add(newRatio, toBase("1", 18)), cfg.MinCollateralizationRatio)
	})

	tests.RunWithSetup("min collateralization ratio cannot be less than liquidation ratio", setup, func(r *tests.Runner) {
		_, err := r.Stabilization.UpdateRatios(r.Operator, toBase("1.5", 18), toBase("1.4", 18))
		require.ErrorAs(r.T, err, &tests.StabilizationInvalidParameterError{})
	})
}

func TestRestrictedFunctionAccess(t *testing.T) {
	testUser := common.Address{200}

	setup := func() *tests.Runner {
		r := tests.Setup(t, nil)
		primePrices(r, newtonAutonPrice, newtonUSDPrice)
		r.GiveMeSomeMoney(testUser, new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
		return r
	}

	tests.RunWithSetup("Set ATN supply operator restricted to operator", setup, func(r *tests.Runner) {
		_, err := r.Stabilization.SetAtnSupplyOperator(tests.FromSender(testUser, common.Big0), testUser)
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted")

		_, err = r.Stabilization.SetAtnSupplyOperator(r.Operator, testUser)
		require.NoError(t, err)
	})

	tests.RunWithSetup("Deposit restricted to ATN supply operator", setup, func(r *tests.Runner) {
		_, err := r.Autonity.Mint(r.Operator, testUser, big.NewInt(1000))
		require.NoError(t, err)

		_, err = r.Autonity.Approve(tests.FromSender(testUser, common.Big0), r.Stabilization.Address(), common.Big256)
		require.NoError(t, err)

		// test user is not the ATN supply operator
		_, err = r.Stabilization.Deposit(tests.FromSender(testUser, common.Big0), big.NewInt(1000))
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted")

		// set test user as the ATN supply operator
		_, err = r.Stabilization.SetAtnSupplyOperator(r.Operator, testUser)
		require.NoError(t, err)

		// test user is now the ATN supply operator
		_, err = r.Stabilization.Deposit(tests.FromSender(testUser, common.Big0), big.NewInt(10))
		require.NoError(t, err)
	})

	tests.RunWithSetup("removeCDPRestrictions restricted to operator", setup, func(r *tests.Runner) {
		_, err := r.Stabilization.RemoveCDPRestrictions(tests.FromSender(testUser, common.Big0))
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted")

		_, err = r.Stabilization.RemoveCDPRestrictions(r.Operator)
		require.NoError(t, err)
	})

	tests.RunWithSetup("removeFixedGenesisPrices restricted to operator", setup, func(r *tests.Runner) {
		unauthorizedUsers := []common.Address{
			testUser,
			params.DeployerAddress,
			testrand.Address(),
		}
		for _, user := range unauthorizedUsers {
			_, err := r.Stabilization.UseFixedGenesisPrices(tests.FromSender(user, common.Big0), false)
			require.Error(t, err)
			require.ErrorAs(t, err, &tests.StabilizationUnauthorizedError{})
		}
		r.NoError(r.Stabilization.UseFixedGenesisPrices(r.Operator, false))
	})
}

func TestInterestRate(t *testing.T) {
	atnSupplyOperator := common.Address{200}
	setup := func() *tests.Runner {
		r := tests.Setup(t, nil)
		r.GiveMeSomeMoney(atnSupplyOperator, big.NewInt(100000000000000000))
		_, err := r.Stabilization.SetAtnSupplyOperator(r.Operator, atnSupplyOperator)
		require.NoError(t, err)
		return r
	}

	tests.RunWithSetup("Interest rate should be zero before restrictions are removed", setup, func(r *tests.Runner) {
		// check interest rate
		config, _, err := r.Stabilization.Config(nil)
		require.NoError(t, err)
		require.Equal(t, uint64(0), config.BorrowInterestRate.Uint64())

		// remove restrictions
		_, err = r.Stabilization.RemoveCDPRestrictions(r.Operator)
		require.NoError(t, err)

		// check interest rate
		config, _, err = r.Stabilization.Config(nil)
		require.NoError(t, err)
		expectedInterestRate, _ := new(big.Int).SetString("50000000000000000", 10)
		require.Equal(t, expectedInterestRate, config.BorrowInterestRate)
	})
}

func TestUpdatableConfigParams(t *testing.T) {
	setup := func() *tests.Runner {
		r := tests.Setup(t, nil)
		r.NoError(r.Stabilization.RemoveCDPRestrictions(r.Operator))
		return r
	}
	tests.RunWithSetup("lastUpdated should return the correct timestamps", setup, func(r *tests.Runner) {
		// wait some time
		progressTime(r, 1000)
		window := getAnnouncementWindow(r)
		expectedTime := int64(r.Evm.Context.Time) + window.Int64()

		r.NoError(
			r.Stabilization.UpdateAnnouncementWindow(
				r.Operator,
				new(big.Int).Add(window, common.Big1),
			),
		)
		r.NoError(
			r.Stabilization.UpdateRatios(
				r.Operator,
				new(big.Int).Add(basicConfig.LiquidationRatio, common.Big1),
				new(big.Int).Add(basicConfig.MinCollateralizationRatio, common.Big1),
			),
		)
		r.NoError(
			r.Stabilization.UpdateBorrowInterestRate(
				r.Operator,
				new(big.Int).Add(basicConfig.BorrowInterestRate, common.Big1),
			),
		)

		progressTime(r, window.Int64()+100)
		// check last updated

		lastUpdated, _, err := r.Stabilization.LastUpdated(nil)
		require.NoError(t, err)

		require.Equal(t, expectedTime, lastUpdated.AnnouncementWindowTimestamp.Int64())
		require.Equal(t, expectedTime, lastUpdated.LiquidationRatioTimestamp.Int64())
		require.Equal(t, expectedTime, lastUpdated.MinCollateralizationRatioTimestamp.Int64())
		require.Equal(t, expectedTime, lastUpdated.BorrowInterestRateTimestamp.Int64())
	})
}

func TestFixedGenesisPrices(t *testing.T) {
	setup := func() *tests.Runner {
		return tests.Setup(t, nil)
	}

	tests.RunWithSetup("Genesis prices are set at fixed rates", setup, func(r *tests.Runner) {
		cfg, _, err := r.Stabilization.Config(nil)
		require.NoError(t, err)

		require.Equal(t, cfg.DefaultNTNATNPrice, (*big.Int)(params.DefaultStabilizationGenesis.DefaultNTNATNPrice))
		require.Equal(t, cfg.DefaultNTNUSDPrice, (*big.Int)(params.DefaultStabilizationGenesis.DefaultNTNUSDPrice))

		// initialize oracle
		or := newOracleTestRounds([]*big.Int{toBase("3.0", 18)})
		or.initialize(r)
		or.increment(r)
		or.increment(r)

		acuDecimals, _, err := r.Acu.GetScale(nil)
		require.NoError(t, err)

		acuPrice, _, err := r.Acu.Value(nil)
		require.NoError(t, err)

		// acu price scaled to oracle decimals
		scaledAcuPrice := new(big.Int).Div(
			new(big.Int).Mul(acuPrice, e18),
			toBase("1.0", acuDecimals.Int64()),
		)

		// $1.0193722 is the default ACU ntnAcuPrice set up by oracleTestRounds
		// this is not dependent on the fixed genesis prices, only the FX prices
		require.Equal(t, acuPrice, toBase("1.0193722", acuDecimals.Int64()))

		acuPriceRead, _, err := r.Stabilization.AcuPrice(nil)
		require.NoError(t, err)
		require.Equal(t, scaledAcuPrice, acuPriceRead)

		// this collateral ntnAcuPrice should be the ntnAcuPrice of NTN in ACU assuming a variable ACU ntnAcuPrice
		ntnAcuPrice, _, err := r.Stabilization.CollateralPriceACU(nil)
		require.NoError(t, err)
		require.Equal(t, new(big.Int).Div(
			new(big.Int).Mul(cfg.DefaultNTNUSDPrice, e18),
			scaledAcuPrice,
		), ntnAcuPrice)

		ntnAtnPrice, _, err := r.Stabilization.CollateralPrice(nil)
		require.NoError(t, err)
		require.Equal(t, cfg.DefaultNTNATNPrice, ntnAtnPrice)

		// should revert to oracle prices once fixed genesis prices are removed
		r.NoError(r.Stabilization.UseFixedGenesisPrices(r.Operator, false))

		// check that the oracle price is set
		oracleNtnAtnPrice, _, err := r.Oracle.LatestRoundData(nil, "NTN-ATN")
		require.NoError(t, err)

		ntnAtnPrice, _, err = r.Stabilization.CollateralPrice(nil)
		require.NoError(t, err)

		require.Equal(t, oracleNtnAtnPrice.Price, ntnAtnPrice)

		oracleNtnUsdPrice, _, err := r.Oracle.LatestRoundData(nil, "NTN-USD")
		require.NoError(t, err)

		ntnAcuPrice, _, err = r.Stabilization.CollateralPriceACU(nil)
		require.NoError(t, err)

		require.NotEqual(t, oracleNtnUsdPrice.Price, cfg.DefaultNTNUSDPrice)
		require.Equal(t, new(big.Int).Div(
			new(big.Int).Mul(oracleNtnUsdPrice.Price, e18),
			scaledAcuPrice,
		), ntnAcuPrice)
	})
}

// test helpers functions

func progressTime(r *tests.Runner, timeToAdd int64) {
	require.True(r.T, timeToAdd >= 0)
	r.Evm.Context.Time += uint64(timeToAdd)
}

func getCdp(r *tests.Runner, user common.Address) tests.IStabilizationCDP {
	cdp, _, err := r.Stabilization.Cdps(nil, user)
	require.NoError(r.T, err)
	return cdp
}

func getAnnouncementWindow(r *tests.Runner) *big.Int {
	window, _, err := r.Stabilization.GetAnnouncementWindow(nil)
	require.NoError(r.T, err)
	return window
}

func getCurrentRate(r *tests.Runner) *big.Int {
	currentRate, _, err := r.Stabilization.GetCurrentRate(nil)
	require.NoError(r.T, err)
	return currentRate
}

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

func borrow(r *tests.Runner, user common.Address, amount *big.Int) {
	stabilizationBalance := r.GetBalanceOf(r.Stabilization.Address())
	supplyCOntrollerBalance := r.GetBalanceOf(r.SupplyControl.Address())
	userBalance := r.GetBalanceOf(user)

	r.NoError(
		r.Stabilization.Borrow(
			tests.FromSender(user, nil),
			amount,
		),
	)
	require.Equal(
		r.T,
		new(big.Int).Add(userBalance, amount),
		r.GetBalanceOf(user),
	)
	require.Equal(
		r.T,
		amount,
		new(big.Int).Sub(supplyCOntrollerBalance, r.GetBalanceOf(r.SupplyControl.Address())),
	)
	require.True(
		r.T,
		r.GetBalanceOf(r.Stabilization.Address()).Cmp(stabilizationBalance) == 0,
	)
}

func calcBorrowLimit(r *tests.Runner, userAccount common.Address) *big.Int {
	cdp, _, err := r.Stabilization.Cdps(nil, userAccount)
	require.NoError(r.T, err)

	limit, _, err := r.Stabilization.MaxBorrow(nil, cdp.Collateral)
	require.NoError(r.T, err)
	return limit
}

func getDebt(r *tests.Runner, user common.Address) *big.Int {
	debt, _, err := r.Stabilization.DebtAmount(nil, user)
	require.NoError(r.T, err)
	return debt
}

func setBasicConfig(r *tests.Runner) {
	r.NoError(r.Stabilization.UpdateRatios(r.Operator, basicConfig.LiquidationRatio, basicConfig.MinCollateralizationRatio))
	progressTime(r, basicConfig.AnnouncementWindow.Int64())
	r.NoError(r.Stabilization.SetMinDebtRequirement(r.Operator, basicConfig.MinDebtRequirement))
}

func bigInt(s string) *big.Int {
	i, _ := new(big.Int).SetString(s, 10)
	return i
}

func primePrices(r *tests.Runner, ntnAtnPrice *big.Int, ntnUsdPrice *big.Int) {
	oracleScale, _, err := r.Oracle.GetDecimals(nil)
	require.NoError(r.T, err)
	oracleDecimals := int64(oracleScale)

	acuScale, _, err := r.Acu.GetScale(nil)
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
	primeOracle(r, append(symbols, "NTN-ATN", "NTN-USD"), append(prices, ntnAtnPrice, ntnUsdPrice))
	_, err = r.Acu.Update(tests.FromAutonity)
	require.NoError(r.T, err)
}

type interestRateParam struct {
	rate, startTime, endTime *big.Int
}

func rateExponent(rateParam interestRateParam) *big.Float {
	return new(big.Float).Quo(
		newFloat(
			new(big.Int).Mul(
				rateParam.rate,
				new(big.Int).Sub(rateParam.endTime, rateParam.startTime),
			),
		),
		new(big.Float).SetPrec(BigFloatPrecision).SetInt64(params.SecondsInYear),
	)
}

func calculateDebt(r *tests.Runner, debtInt *big.Int, debtStartTime *big.Int, rates []interestRateParam) *big.Int {
	// assuming rates are sorted according to their `startTime` in ascending order and we have `endTime[i] = startTime[i+1]`
	// verify the above assumption
	for i, rateParam := range rates {
		require.True(r.T, rateParam.startTime.Cmp(rateParam.endTime) <= 0, "fails", i, rateParam.startTime, rateParam.endTime)
		if i > 0 {
			require.True(r.T, rateParam.startTime.Cmp(rates[i-1].endTime) == 0)
		}
	}

	e18Float := newFloat(e18)
	debt := newFloat(debtInt)
	// debt = new(big.Float).Quo(
	// 	debt,
	// 	e18Float,
	// )
	for i, rateParam := range rates {
		// exclude rates that are not applicable
		if rateParam.endTime.Cmp(debtStartTime) <= 0 {
			continue
		}
		if rateParam.startTime.Cmp(debtStartTime) < 0 {
			// the current rate will not be applied fully
			rt := new(big.Float).Quo(
				rateExponent(
					interestRateParam{
						rate:      rateParam.rate,
						startTime: debtStartTime,
						endTime:   rateParam.endTime,
					},
				),
				e18Float,
			)
			ert := bigfloat.Exp(
				rt,
			)
			debt = new(big.Float).Mul(
				debt,
				ert,
			)
			rates = rates[i+1:]
		} else {
			rates = rates[i:]
		}
		break
	}

	// apply the rest
	for _, rateParam := range rates {
		rt := new(big.Float).Quo(
			rateExponent(rateParam),
			e18Float,
		)
		ert := bigfloat.Exp(
			rt,
		)
		debt = new(big.Float).Mul(
			debt,
			ert,
		)
	}
	debtInt, _ = debt.Int(nil)
	return debtInt
}

func makeLiquidatable(r *tests.Runner, user common.Address) {
	cfg, _, err := r.Stabilization.Config(nil)
	require.NoError(r.T, err)

	r.NoError(r.Stabilization.UpdateRatios(
		r.Operator,
		new(big.Int).Mul(cfg.MinCollateralizationRatio, big.NewInt(2)),
		new(big.Int).Mul(cfg.MinCollateralizationRatio, big.NewInt(3)),
	))
	progressTime(r, cfg.AnnouncementWindow.Int64())

	require.NoError(r.T, err)
	progressTime(r, cfg.AnnouncementWindow.Int64())
	liquidatable, _, err := r.Stabilization.IsLiquidatable(nil, user)
	require.NoError(r.T, err)
	require.True(r.T, liquidatable)
}
