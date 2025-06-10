package asm

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/autonity/tests"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/internal/testrand"
)

var scaleFactor = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)

func TestAuctioneerInterestAuction(t *testing.T) {
	var oracleScaleFactor *big.Int
	setup := func() *tests.Runner {
		r := tests.Setup(t, nil)
		decimals, _, err := r.Oracle.GetDecimals(nil)
		require.NoError(t, err)
		oracleScaleFactor = new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
		primePrices(r, newtonAutonPrice, newtonUSDPrice)
		r.NoError(r.Stabilization.UseFixedGenesisPrices(r.Operator, false))
		return r
	}

	tests.RunWithSetup("Paying interest above threshold starts an auction", setup, func(r *tests.Runner) {
		// Verify no initial open auctions
		auctions, _, err := r.Auctioneer.OpenAuctions(nil)
		require.NoError(t, err)
		require.Len(t, auctions, 0)

		config, _, err := r.Auctioneer.GetConfig(nil)
		require.NoError(t, err)

		interestPayment := new(big.Int).Add(config.InterestAuctionThreshold, big.NewInt(1000))
		r.GiveMeSomeMoney(
			r.Stabilization.Address(),
			interestPayment,
		)

		_, err = r.Auctioneer.PaidInterest(tests.FromSender(r.Stabilization.Address(), interestPayment))
		require.NoError(t, err)

		// Check that an auction has been started
		auctions, _, err = r.Auctioneer.OpenAuctions(nil)
		require.NoError(t, err)
		require.Len(t, auctions, 1)

		// Check that the auction is for the correct amount
		require.Equal(t, interestPayment, auctions[0].Amount)
	})

	tests.RunWithSetup("Paying interest below threshold does not start an auction", setup, func(r *tests.Runner) {
		auctions, _, err := r.Auctioneer.OpenAuctions(nil)
		require.NoError(t, err)
		require.Len(t, auctions, 0)

		config, _, err := r.Auctioneer.GetConfig(nil)
		require.NoError(t, err)

		interestPayment := new(big.Int).Sub(config.InterestAuctionThreshold, big.NewInt(1000))
		r.GiveMeSomeMoney(r.Stabilization.Address(), interestPayment)

		_, err = r.Auctioneer.PaidInterest(tests.FromSender(r.Stabilization.Address(), interestPayment))
		require.NoError(t, err)

		auctions, _, err = r.Auctioneer.OpenAuctions(nil)
		require.NoError(t, err)
		require.Len(t, auctions, 0)
	})

	tests.RunWithSetup("Interest auction price starts at market price plus discount", setup, func(r *tests.Runner) {
		auctionAmount := setupInterestAuction(r, big.NewInt(1000))
		round, _, err := r.Oracle.LatestRoundData(nil, "NTN-ATN")
		require.NoError(t, err)

		// Check that the auction is for the correct amount
		auctions, _, err := r.Auctioneer.OpenAuctions(nil)
		require.NoError(t, err)
		require.Len(t, auctions, 1)

		// Check that the auction start info is correct
		auction := auctions[0]
		require.NoError(t, err)
		require.Equal(t, round.Price.Int64(), auction.StartPrice.Int64())
		require.Equal(t, auction.StartTimestamp, r.Evm.Context.Time)
		require.Equal(t, auctionAmount, auction.Amount)
		require.Equal(t, newtonAutonPrice, auction.StartPrice)

		ntnCost, _, err := r.Auctioneer.MinInterestPayment(nil, auction.Id)
		require.NoError(t, err)

		config, _, err := r.Auctioneer.GetConfig(nil)
		require.NoError(t, err)

		priceDiscount := newFloat0().Sub(
			newFloat(auction.StartPrice),
			newFloat0().Quo(
				newFloat0().Mul(newFloat(auction.StartPrice), newFloat(config.InterestAuctionDiscount)),
				newFloat(scaleFactor),
			),
		)
		expectedCost := newFloat0().Quo(
			newFloat0().Mul(newFloat(auctionAmount), newFloat(oracleScaleFactor)),
			priceDiscount,
		)
		expectedCostInt, _ := expectedCost.Int(nil)
		require.Equal(t, expectedCostInt, ntnCost)
	})

	tests.RunWithSetup("Interest auction price ends at 0", setup, func(r *tests.Runner) {
		auctionAmount := setupInterestAuction(r, big.NewInt(1000))
		config, _, err := r.Auctioneer.GetConfig(nil)
		require.NoError(t, err)

		// Check that the auction is for the correct amount
		auctions, _, err := r.Auctioneer.OpenAuctions(nil)
		require.NoError(t, err)
		require.Len(t, auctions, 1)
		auction := auctions[0]
		require.Equal(t, auctionAmount, auction.Amount)

		//wait until the auction ends
		auctionEnd := new(big.Int).Add(auction.StartTimestamp, config.InterestAuctionDuration)

		for r.Evm.Context.Time.Cmp(auctionEnd) < 0 {
			r.WaitNBlocks(1)
		}

		// check the auction price
		ntnCost, _, err := r.Auctioneer.MinInterestPayment(nil, auction.Id)
		require.NoError(t, err)
		require.Equal(t, uint64(0), ntnCost.Uint64())
	})

	tests.RunWithSetup("Test bid is accepted if it's at the minimum", setup, func(r *tests.Runner) {
		auctionAmount := setupInterestAuction(r, big.NewInt(1000))

		auctions, _, err := r.Auctioneer.OpenAuctions(nil)
		require.NoError(t, err)
		require.Len(t, auctions, 1)
		auction := auctions[0]

		ntnCost, _, err := r.Auctioneer.MinInterestPayment(nil, auction.Id)
		require.NoError(t, err)

		// bid at the minimum
		userAccount := testrand.Address()

		_, err = r.Autonity.Mint(r.Operator, userAccount, ntnCost)
		require.NoError(t, err)

		_, err = r.Autonity.Approve(tests.FromSender(userAccount, common.Big0), r.Auctioneer.Address(), ntnCost)
		require.NoError(t, err)

		ntnBalanceBefore, _, err := r.Autonity.BalanceOf(nil, userAccount)
		require.NoError(t, err)
		atnBalanceBefore := r.GetBalanceOf(userAccount)

		_, err = r.Auctioneer.BidInterest(tests.FromSender(userAccount, common.Big0), auction.Id, ntnCost)
		require.NoError(t, err)

		ntnBalanceAfter, _, err := r.Autonity.BalanceOf(nil, userAccount)
		require.NoError(t, err)
		atnBalanceAfter := r.GetBalanceOf(userAccount)

		require.Equal(t, new(big.Int).Sub(ntnBalanceBefore, ntnCost), ntnBalanceAfter)
		require.Equal(t, new(big.Int).Add(atnBalanceBefore, auctionAmount), atnBalanceAfter)

		// Check that the auction is closed
		auctions, _, err = r.Auctioneer.OpenAuctions(nil)
		require.NoError(t, err)
		require.Len(t, auctions, 0)
	})

	tests.RunWithSetup("Test bid is accepted if it's above the minimum", setup, func(r *tests.Runner) {
		auctionAmount := setupInterestAuction(r, big.NewInt(9000))

		auctions, _, err := r.Auctioneer.OpenAuctions(nil)
		require.NoError(t, err)
		require.Len(t, auctions, 1)
		auction := auctions[0]

		ntnCost, _, err := r.Auctioneer.MinInterestPayment(nil, auction.Id)
		require.NoError(t, err)

		// bid above the minimum
		bidder := testrand.Address()
		// bid twice the amount
		bidAmount := new(big.Int).Add(ntnCost, ntnCost)

		_, err = r.Autonity.Mint(r.Operator, bidder, bidAmount)
		require.NoError(t, err)

		atnBalanceBefore := r.GetBalanceOf(bidder)
		ntnBalanceBefore, _, err := r.Autonity.BalanceOf(nil, bidder)
		require.NoError(t, err)

		_, err = r.Autonity.Approve(tests.FromSender(bidder, common.Big0), r.Auctioneer.Address(), bidAmount)
		require.NoError(t, err)

		_, err = r.Auctioneer.BidInterest(tests.FromSender(bidder, common.Big0), auction.Id, bidAmount)
		require.NoError(t, err)

		ntnBalanceAfter, _, err := r.Autonity.BalanceOf(nil, bidder)
		require.NoError(t, err)
		atnBalanceAfter := r.GetBalanceOf(bidder)

		require.Equal(t, bidAmount, new(big.Int).Sub(ntnBalanceBefore, ntnBalanceAfter))
		require.Equal(t, auctionAmount, new(big.Int).Sub(atnBalanceAfter, atnBalanceBefore))
	})

	tests.RunWithSetup("Test bid is rejected if it's below the minimum", setup, func(r *tests.Runner) {
		setupInterestAuction(r, big.NewInt(1000))

		auctions, _, err := r.Auctioneer.OpenAuctions(nil)
		require.NoError(t, err)
		require.Len(t, auctions, 1)
		auction := auctions[0]

		ntnCost, _, err := r.Auctioneer.MinInterestPayment(nil, auction.Id)
		require.NoError(t, err)

		// bid below the minimum
		bidder := testrand.Address()
		// bid half the amount
		bidAmount := new(big.Int).Div(ntnCost, big.NewInt(2))

		_, err = r.Autonity.Mint(r.Operator, bidder, bidAmount)
		require.NoError(t, err)

		atnBalanceBefore := r.GetBalanceOf(bidder)
		ntnBalanceBefore, _, err := r.Autonity.BalanceOf(nil, bidder)
		require.NoError(t, err)

		_, err = r.Autonity.Approve(tests.FromSender(bidder, common.Big0), r.Auctioneer.Address(), bidAmount)
		require.NoError(t, err)

		_, err = r.Auctioneer.BidInterest(tests.FromSender(bidder, common.Big0), auction.Id, bidAmount)
		require.Error(t, err)
		require.ErrorAs(t, err, &tests.AuctioneerBidTooLowError{})

		ntnBalanceAfter, _, err := r.Autonity.BalanceOf(nil, bidder)
		require.NoError(t, err)
		atnBalanceAfter := r.GetBalanceOf(bidder)

		require.Equal(t, uint64(0), new(big.Int).Sub(ntnBalanceBefore, ntnBalanceAfter).Uint64())
		require.Equal(t, uint64(0), new(big.Int).Sub(atnBalanceAfter, atnBalanceBefore).Uint64())
	})
}

func TestAuctioneerDebtAuction(t *testing.T) {
	setup := func() *tests.Runner {
		r := tests.Setup(t, nil)
		r.NoError(r.Stabilization.RemoveCDPRestrictions(r.Operator))
		r.NoError(r.Stabilization.UseFixedGenesisPrices(r.Operator, false))
		return r
	}

	tests.RunWithSetup("Debt auction is not callable on a non-liquidatable cdp", setup, func(r *tests.Runner) {
		user, _ := setupCDP(r, toBase("100.00", 18), []*big.Int{newtonAutonPrice}, false)

		cdp, _, err := r.Stabilization.Cdps(nil, user)
		require.NoError(t, err)
		require.True(t, cdp.Principal.Cmp(common.Big0) > 0, "CDP should have a non-zero principal")

		debtAmount, _, err := r.Stabilization.DebtAmount(nil, user)
		require.NoError(t, err)

		liquidatable, _, err := r.Stabilization.IsLiquidatable(nil, user)
		require.NoError(t, err)
		require.False(t, liquidatable, "CDP should not be liquidatable")

		oracleRound, _, err := r.Oracle.GetRound(nil)
		require.NoError(t, err)

		liquidator := testrand.Address()
		r.GiveMeSomeMoney(liquidator, debtAmount)

		_, err = r.Auctioneer.BidDebt(
			tests.FromSender(liquidator, debtAmount),
			user,
			new(big.Int).Sub(oracleRound, common.Big1),
			big.NewInt(1000),
		)
		require.Error(t, err)
		require.ErrorAs(t, err, &tests.AuctioneerInvalidRoundError{})
	})

	tests.RunWithSetup("Cannot bid on a debt auction that with a round before the cdp was created", setup, func(r *tests.Runner) {
		or := newOracleTestRounds([]*big.Int{
			newtonAutonPrice,
			newtonAutonPrice,
			newtonAutonPrice,
			new(big.Int).Div(newtonAutonPrice, common.Big2),
			new(big.Int).Div(newtonAutonPrice, common.Big2),
			new(big.Int).Div(newtonAutonPrice, common.Big2),
		})

		or.initialize(r)
		or.increment(r)
		or.increment(r)

		invalidOracleRound, _, err := r.Oracle.GetRound(nil)
		require.NoError(t, err)

		or.increment(r)

		depositAmount := toBase("100.00", 18)
		borrowLimit, _, err := r.Stabilization.MaxBorrow(nil, depositAmount)
		require.NoError(t, err)

		user := testrand.Address()
		_, err = r.Autonity.Mint(r.Operator, user, depositAmount)
		require.NoError(t, err)

		_, err = r.Autonity.Approve(tests.FromSender(user, common.Big0), r.Stabilization.Address(), depositAmount)
		require.NoError(t, err)

		_, err = r.Stabilization.Deposit(tests.FromSender(user, common.Big0), depositAmount)
		require.NoError(t, err)

		_, err = r.Stabilization.Borrow(tests.FromSender(user, nil), borrowLimit)
		require.NoError(t, err)

		or.increment(r)

		// cdp should be liquidatable from interest accumulation
		liquidatable, _, err := r.Stabilization.IsLiquidatable(nil, user)
		require.NoError(t, err)
		require.True(t, liquidatable, "CDP should be liquidatable")

		// try to bid on a debt auction with an invalid round
		liquidator := testrand.Address()
		debtAmount, _, err := r.Stabilization.DebtAmount(nil, user)
		require.NoError(t, err)
		r.GiveMeSomeMoney(liquidator, debtAmount)

		_, err = r.Auctioneer.BidDebt(
			tests.FromSender(liquidator, debtAmount),
			user,
			new(big.Int).Sub(invalidOracleRound, common.Big1),
			big.NewInt(1000),
		)
		require.Error(t, err)
		require.ErrorAs(t, err, &tests.AuctioneerInvalidRoundError{})
	})

	tests.RunWithSetup("Cannot bid on a debt auction that is no longer liquidatable", setup, func(r *tests.Runner) {
		or := newOracleTestRounds([]*big.Int{
			newtonAutonPrice,
			newtonAutonPrice,
			new(big.Int).Div(newtonAutonPrice, common.Big2),
			newtonAutonPrice,
			new(big.Int).Mul(newtonAutonPrice, common.Big2),
			new(big.Int).Mul(newtonAutonPrice, common.Big2),
		})

		or.initialize(r)
		or.increment(r)
		or.increment(r)

		depositAmount := toBase("100.00", 18)
		borrowLimit, _, err := r.Stabilization.MaxBorrow(nil, depositAmount)
		require.NoError(t, err)

		user := testrand.Address()
		_, err = r.Autonity.Mint(r.Operator, user, depositAmount)
		require.NoError(t, err)

		_, err = r.Autonity.Approve(tests.FromSender(user, common.Big0), r.Stabilization.Address(), depositAmount)
		require.NoError(t, err)

		_, err = r.Stabilization.Deposit(tests.FromSender(user, common.Big0), depositAmount)
		require.NoError(t, err)

		_, err = r.Stabilization.Borrow(tests.FromSender(user, nil), borrowLimit)
		require.NoError(t, err)

		or.increment(r)

		// cdp should be liquidatable from price decrease
		liquidatable, _, err := r.Stabilization.IsLiquidatable(nil, user)
		require.NoError(t, err)
		require.True(t, liquidatable, "CDP should be liquidatable")

		liquidatableRound, _, err := r.Oracle.GetRound(nil)
		require.NoError(t, err)
		liquidatableRound = new(big.Int).Sub(liquidatableRound, common.Big1)

		// increment the oracle rounds to make the cdp no longer liquidatable
		or.increment(r)
		or.increment(r)
		or.increment(r)

		// cdp should no longer be liquidatable
		liquidatable, _, err = r.Stabilization.IsLiquidatable(nil, user)
		require.NoError(t, err)
		require.False(t, liquidatable, "CDP should not be liquidatable")

		liquidator := testrand.Address()
		debtAmount, _, err := r.Stabilization.DebtAmount(nil, user)
		require.NoError(t, err)
		r.GiveMeSomeMoney(liquidator, debtAmount)

		_, err = r.Auctioneer.BidDebt(
			tests.FromSender(liquidator, debtAmount),
			user,
			liquidatableRound,
			big.NewInt(1000),
		)

		require.Error(t, err)
		require.ErrorAs(t, err, &tests.AuctioneerNotLiquidatableError{})
	})

	tests.RunWithSetup("Cannot bid on a debt auction with a round before the cdp was last updated", setup, func(r *tests.Runner) {
		or := newOracleTestRounds([]*big.Int{
			newtonAutonPrice,
			newtonAutonPrice,
			new(big.Int).Div(newtonAutonPrice, common.Big2),
			new(big.Int).Div(newtonAutonPrice, common.Big2),
			new(big.Int).Div(newtonAutonPrice, common.Big4),
			new(big.Int).Div(newtonAutonPrice, common.Big4),
			new(big.Int).Div(newtonAutonPrice, common.Big4),
		})
		or.initialize(r)
		or.increment(r)
		or.increment(r)

		depositAmount := toBase("100.00", 18)
		borrowLimit, _, err := r.Stabilization.MaxBorrow(nil, depositAmount)
		require.NoError(t, err)

		user := testrand.Address()
		_, err = r.Autonity.Mint(r.Operator, user, depositAmount)
		require.NoError(t, err)

		_, err = r.Autonity.Approve(tests.FromSender(user, common.Big0), r.Stabilization.Address(), depositAmount)
		require.NoError(t, err)

		_, err = r.Stabilization.Deposit(tests.FromSender(user, common.Big0), depositAmount)
		require.NoError(t, err)

		_, err = r.Stabilization.Borrow(tests.FromSender(user, nil), borrowLimit)
		require.NoError(t, err)

		or.increment(r)

		// cdp should be liquidatable from price decrease
		liquidatable, _, err := r.Stabilization.IsLiquidatable(nil, user)
		require.NoError(t, err)
		require.True(t, liquidatable, "CDP should be liquidatable")

		liquidatableRound, _, err := r.Oracle.GetRound(nil)
		require.NoError(t, err)
		liquidatableRound = new(big.Int).Sub(liquidatableRound, common.Big1)

		or.increment(r)
		// increase the cdp's deposit
		_, err = r.Autonity.Mint(r.Operator, user, depositAmount)
		require.NoError(t, err)

		_, err = r.Autonity.Approve(tests.FromSender(user, common.Big0), r.Stabilization.Address(), depositAmount)
		require.NoError(t, err)

		_, err = r.Stabilization.Deposit(tests.FromSender(user, common.Big0), depositAmount)
		require.NoError(t, err)

		// should no longer be liquidatable
		liquidatable, _, err = r.Stabilization.IsLiquidatable(nil, user)
		require.NoError(t, err)
		require.False(t, liquidatable, "CDP should not be liquidatable")

		// wait for liquidatable again
		or.increment(r)
		or.increment(r)

		liquidatable, _, err = r.Stabilization.IsLiquidatable(nil, user)
		require.NoError(t, err)
		require.True(t, liquidatable, "CDP should be liquidatable")

		liquidator := testrand.Address()
		debtAmount, _, err := r.Stabilization.DebtAmount(nil, user)
		require.NoError(t, err)

		r.GiveMeSomeMoney(liquidator, debtAmount)

		_, err = r.Auctioneer.BidDebt(
			tests.FromSender(liquidator, debtAmount),
			user,
			liquidatableRound,
			big.NewInt(1000),
		)
		require.Error(t, err)
		require.ErrorAs(t, err, &tests.AuctioneerInvalidRoundError{})
	})

	tests.RunWithSetup("Cannot use an oracle round before the last liquidation ratio update", setup, func(r *tests.Runner) {
		user, or := setupCDP(r, toBase("100.00", 18), []*big.Int{newtonAutonPrice}, true)
		makeLiquidatable(r, user)

		liquidatable, _, err := r.Stabilization.IsLiquidatable(nil, user)
		require.NoError(t, err)
		require.True(t, liquidatable, "CDP should be liquidatable")

		liquidatableRound, _, err := r.Oracle.GetRound(nil)
		require.NoError(t, err)

		or.increment(r)
		or.increment(r)

		cfg, _, err := r.Stabilization.Config(nil)
		r.NoError(r.Stabilization.UpdateRatios(
			r.Operator,
			new(big.Int).Add(cfg.LiquidationRatio, common.Big1),
			cfg.MinCollateralizationRatio,
		))
		progressTime(r, cfg.AnnouncementWindow.Int64()+1)

		updated, _, err := r.Stabilization.LastUpdated(nil)
		require.NoError(t, err)
		roundData, _, err := r.Oracle.GetRoundData(nil, liquidatableRound, "NTN-ATN")
		require.NoError(t, err)
		require.True(t, updated.LiquidationRatioTimestamp.Cmp(roundData.Timestamp) > 0)

		liquidator := testrand.Address()
		debtAmount, _, err := r.Stabilization.DebtAmount(nil, user)
		r.GiveMeSomeMoney(liquidator, new(big.Int).Mul(debtAmount, common.Big2))

		_, err = r.Auctioneer.BidDebt(
			tests.FromSender(liquidator, debtAmount),
			user,
			liquidatableRound,
			big.NewInt(1000),
		)
		require.ErrorAs(t, err, &tests.AuctioneerInvalidRoundError{})
	})

	tests.RunWithSetup("Can bid with maximum current amount", setup, func(r *tests.Runner) {
		user, or := setupCDP(
			r,
			toBase("100.00", 18),
			[]*big.Int{
				newtonAutonPrice,
				newtonAutonPrice,
				new(big.Int).Div(newtonAutonPrice, common.Big2),
				new(big.Int).Div(newtonAutonPrice, common.Big2),
				new(big.Int).Div(newtonAutonPrice, common.Big2),
			},
			true,
		)
		or.increment(r)

		// cdp should be liquidatable from price decrease
		liquidatable, _, err := r.Stabilization.IsLiquidatable(nil, user)
		require.NoError(t, err)
		require.True(t, liquidatable, "CDP should be liquidatable")

		liquidatableRound, _, err := r.Oracle.GetRound(nil)
		require.NoError(t, err)
		liquidatableRound = new(big.Int).Sub(liquidatableRound, common.Big1)

		or.increment(r)

		liquidator := testrand.Address()
		debtAmount, _, err := r.Stabilization.DebtAmount(nil, user)
		require.NoError(t, err)
		r.GiveMeSomeMoney(liquidator, debtAmount)

		maxReturn, _, err := r.Auctioneer.MaxLiquidationReturn(nil, user, liquidatableRound)
		require.NoError(t, err)

		balanceBefore, _, err := r.Autonity.BalanceOf(nil, liquidator)
		require.NoError(t, err)

		atnBalanceBefore := r.GetBalanceOf(liquidator)

		_, err = r.Auctioneer.BidDebt(
			tests.FromSender(liquidator, debtAmount),
			user,
			liquidatableRound,
			maxReturn,
		)
		require.NoError(t, err)

		balanceAfter, _, err := r.Autonity.BalanceOf(nil, liquidator)
		require.NoError(t, err)

		atnBalanceAfter := r.GetBalanceOf(liquidator)

		require.Equal(t, new(big.Int).Add(balanceBefore, maxReturn), balanceAfter)
		require.Equal(t, new(big.Int).Sub(atnBalanceBefore, debtAmount), atnBalanceAfter)
	})

	tests.RunWithSetup("Cannot bid with above maximum current amount", setup, func(r *tests.Runner) {
		user, or := setupCDP(
			r,
			toBase("100.00", 18),
			[]*big.Int{
				newtonAutonPrice,
				newtonAutonPrice,
				new(big.Int).Div(newtonAutonPrice, common.Big2),
				new(big.Int).Div(newtonAutonPrice, common.Big2),
				new(big.Int).Div(newtonAutonPrice, common.Big2),
			},
			true,
		)
		or.increment(r)

		// cdp should be liquidatable from price decrease
		liquidatable, _, err := r.Stabilization.IsLiquidatable(nil, user)
		require.NoError(t, err)
		require.True(t, liquidatable, "CDP should be liquidatable")

		liquidatableRound, _, err := r.Oracle.GetRound(nil)
		require.NoError(t, err)
		liquidatableRound = new(big.Int).Sub(liquidatableRound, common.Big1)

		or.increment(r)

		liquidator := testrand.Address()
		debtAmount, _, err := r.Stabilization.DebtAmount(nil, user)
		require.NoError(t, err)
		r.GiveMeSomeMoney(liquidator, debtAmount)

		maxReturn, _, err := r.Auctioneer.MaxLiquidationReturn(nil, user, liquidatableRound)
		require.NoError(t, err)

		// bid above the maximum
		_, err = r.Auctioneer.BidDebt(
			tests.FromSender(liquidator, debtAmount),
			user,
			liquidatableRound,
			new(big.Int).Add(maxReturn, common.Big1),
		)
		require.Error(t, err)
		require.ErrorAs(t, err, &tests.AuctioneerBidTooLowError{})
	})

	tests.RunWithSetup("Can bid below maximum current amount", setup, func(r *tests.Runner) {
		user, or := setupCDP(
			r,
			toBase("100.00", 18),
			[]*big.Int{
				newtonAutonPrice,
				newtonAutonPrice,
				new(big.Int).Div(newtonAutonPrice, common.Big2),
				new(big.Int).Div(newtonAutonPrice, common.Big2),
				new(big.Int).Div(newtonAutonPrice, common.Big2),
			},
			true,
		)
		or.increment(r)

		// cdp should be liquidatable from price decrease
		liquidatable, _, err := r.Stabilization.IsLiquidatable(nil, user)
		require.NoError(t, err)
		require.True(t, liquidatable, "CDP should be liquidatable")

		liquidatableRound, _, err := r.Oracle.GetRound(nil)
		require.NoError(t, err)
		liquidatableRound = new(big.Int).Sub(liquidatableRound, common.Big1)

		or.increment(r)

		liquidator := testrand.Address()
		debtAmount, _, err := r.Stabilization.DebtAmount(nil, user)
		require.NoError(t, err)

		r.GiveMeSomeMoney(liquidator, debtAmount)

		lowBid := big.NewInt(1000)
		maxReturn, _, err := r.Auctioneer.MaxLiquidationReturn(nil, user, liquidatableRound)
		require.NoError(t, err)
		require.True(t, maxReturn.Cmp(lowBid) > 0)

		balanceBefore, _, err := r.Autonity.BalanceOf(nil, liquidator)
		require.NoError(t, err)
		atnBalanceBefore := r.GetBalanceOf(liquidator)

		_, err = r.Auctioneer.BidDebt(
			tests.FromSender(liquidator, debtAmount),
			user,
			liquidatableRound,
			lowBid,
		)

		require.NoError(t, err)

		balanceAfter, _, err := r.Autonity.BalanceOf(nil, liquidator)
		require.NoError(t, err)
		atnBalanceAfter := r.GetBalanceOf(liquidator)

		require.Equal(t, new(big.Int).Add(balanceBefore, lowBid), balanceAfter)
		require.Equal(t, new(big.Int).Sub(atnBalanceBefore, debtAmount), atnBalanceAfter)
	})
}

func TestAuctioneerSetters(t *testing.T) {
	setup := func() *tests.Runner {
		return tests.Setup(t, nil)
	}

	tests.RunWithSetup("Only operator functions", setup, func(r *tests.Runner) {
		// test setInterestAuctionThreshold
		_, err := r.Auctioneer.SetInterestAuctionThreshold(r.Operator, big.NewInt(1000))
		require.NoError(t, err)

		_, err = r.Auctioneer.SetInterestAuctionThreshold(
			tests.FromSender(testrand.Address(), common.Big0),
			big.NewInt(1000),
		)
		require.ErrorAs(t, err, &tests.AuctioneerUnauthorizedError{})

		// test setInterestAuctionDiscount
		_, err = r.Auctioneer.SetInterestAuctionDiscount(r.Operator, big.NewInt(1000))
		require.NoError(t, err)

		_, err = r.Auctioneer.SetInterestAuctionDiscount(
			tests.FromSender(testrand.Address(), common.Big0),
			big.NewInt(1000),
		)
		require.ErrorAs(t, err, &tests.AuctioneerUnauthorizedError{})

		// test setLiquidationAuctionDuration
		_, err = r.Auctioneer.SetLiquidationAuctionDuration(r.Operator, big.NewInt(1000))
		require.NoError(t, err)

		_, err = r.Auctioneer.SetLiquidationAuctionDuration(
			tests.FromSender(testrand.Address(), common.Big0),
			big.NewInt(1000),
		)
		require.ErrorAs(t, err, &tests.AuctioneerUnauthorizedError{})

		// test setInterestAuctionDuration
		_, err = r.Auctioneer.SetInterestAuctionDuration(r.Operator, big.NewInt(1000))
		require.NoError(t, err)

		_, err = r.Auctioneer.SetInterestAuctionDuration(
			tests.FromSender(testrand.Address(), common.Big0),
			big.NewInt(1000),
		)
		require.ErrorAs(t, err, &tests.AuctioneerUnauthorizedError{})

		// test set proceeds address
		_, err = r.Auctioneer.SetProceedAddress(r.Operator, testrand.Address())
		require.NoError(t, err)

		_, err = r.Auctioneer.SetProceedAddress(
			tests.FromSender(testrand.Address(), common.Big0),
			testrand.Address(),
		)
		require.ErrorAs(t, err, &tests.AuctioneerUnauthorizedError{})
	})

}

func setupInterestAuction(r *tests.Runner, atnAboveMininum *big.Int) (auctionAmount *big.Int) {
	or := newOracleTestRounds([]*big.Int{newtonAutonPrice, newtonAutonPrice, newtonAutonPrice, newtonAutonPrice})
	or.initialize(r)
	or.increment(r)
	or.increment(r)

	// start an auction
	config, _, err := r.Auctioneer.GetConfig(nil)
	require.NoError(r.T, err)

	auctionAmount = new(big.Int).Add(config.InterestAuctionThreshold, atnAboveMininum)
	// start an interest auction
	r.GiveMeSomeMoney(r.Stabilization.Address(), auctionAmount)
	_, err = r.Auctioneer.PaidInterest(tests.FromSender(r.Stabilization.Address(), auctionAmount))
	require.NoError(r.T, err)

	return auctionAmount
}

func setupCDP(r *tests.Runner, cdpAmount *big.Int, ntnPrices []*big.Int, borrowMax bool) (user common.Address, or *testOracleRounds) {
	// initialize oracle
	or = newOracleTestRounds(ntnPrices)
	or.initialize(r)
	or.increment(r)
	or.increment(r)

	user = testrand.Address()
	_, err := r.Autonity.Mint(r.Operator, user, cdpAmount)
	require.NoError(r.T, err)

	_, err = r.Autonity.Approve(tests.FromSender(user, common.Big0), r.Stabilization.Address(), cdpAmount)
	require.NoError(r.T, err)

	_, err = r.Stabilization.Deposit(tests.FromSender(user, common.Big0), cdpAmount)
	require.NoError(r.T, err)

	maxBorrow, _, err := r.Stabilization.MaxBorrow(
		nil,
		cdpAmount,
	)
	require.NoError(r.T, err)
	if !borrowMax {
		maxBorrow = new(big.Int).Div(maxBorrow, big.NewInt(2))
	}

	_, err = r.Stabilization.Borrow(tests.FromSender(user, nil), maxBorrow)
	require.NoError(r.T, err)

	return user, or
}

// testOracleRounds is a helper struct for testing the auctioneer with an oracle
type testOracleRounds struct {
	currentRound int
	symbols      []string
	symbolPrices [][]*big.Int
}

func newOracleTestRounds(ntnPrices []*big.Int) *testOracleRounds {
	or := &testOracleRounds{
		currentRound: 0,
	}
	oracleDecimals := int64(18)

	// here we are assuming the price of ATN is equal to the target price (in ACU)
	// ACU price = 1.0193722 USD from the parameters we set here
	atnUSDPrice := toBase("1.65", oracleDecimals)
	ntnUSDPrice := make([]*big.Int, len(ntnPrices))
	for i := range ntnPrices {
		ntnUSDPrice[i] = new(big.Int).Quo(
			new(big.Int).Mul(ntnPrices[i], atnUSDPrice),
			new(big.Int).SetUint64(1e18),
		)
	}
	data := map[string][]*big.Int{
		"NTN-ATN": ntnPrices,
		"NTN-USD": {newtonUSDPrice},
		"AUD-USD": {toBase("0.6757", oracleDecimals)},
		"CAD-USD": {toBase("0.75694", oracleDecimals)},
		"EUR-USD": {toBase("1.1085", oracleDecimals)},
		"GBP-USD": {toBase("1.29403", oracleDecimals)},
		"JPY-USD": {toBase("0.00713", oracleDecimals)},
		"USD-USD": {toBase("1.0", oracleDecimals)},
		"SEK-USD": {toBase("0.09597", oracleDecimals)},
		"ATN-USD": {atnUSDPrice},
	}
	for symbol, prices := range data {
		or.symbols = append(or.symbols, symbol)
		or.symbolPrices = append(or.symbolPrices, prices)
	}

	return or
}

func (o *testOracleRounds) initialize(r *tests.Runner) {
	config, _, err := r.Oracle.GetConfig(nil)
	require.NoError(r.T, err)

	_, err = r.Oracle.SetSymbols(r.Operator, o.symbols)
	require.NoError(r.T, err)
	r.WaitNBlocks(2 * int(config.VotePeriod.Int64()))

	acuScale, _, err := r.Acu.GetScale(nil)
	require.NoError(r.T, err)
	acuDecimals := acuScale.Int64()

	primeACUBasket(
		r,
		[]string{
			"AUD-USD",
			"CAD-USD",
			"EUR-USD",
			"GBP-USD",
			"JPY-USD",
			"USD-USD",
			"SEK-USD",
		},
		[]*big.Int{
			toBase("0.213", acuDecimals),
			toBase("0.187", acuDecimals),
			toBase("0.143", acuDecimals),
			toBase("0.104", acuDecimals),
			toBase("17.6", acuDecimals),
			toBase("0.180", acuDecimals),
			toBase("1.41", acuDecimals),
		},
		acuScale,
	)
}

func (o *testOracleRounds) increment(r *tests.Runner) {
	config, _, err := r.Oracle.GetConfig(nil)
	require.NoError(r.T, err)
	voter := r.Committee.Validators[0].OracleAddress

	nextReports := make([]tests.IOracleReport, len(o.symbolPrices))
	currentReports := make([]tests.IOracleReport, len(o.symbolPrices))

	for i, prices := range o.symbolPrices {
		if o.currentRound > 0 {
			currentReports[i] = tests.IOracleReport{
				Price:      prices[o.currentRound%len(prices)],
				Confidence: 100,
			}
		}
		nextReports[i] = tests.IOracleReport{
			Price:      prices[(o.currentRound+1)%len(prices)],
			Confidence: 100,
		}
	}

	salt := big.NewInt(1234)
	commit := tests.MakeOracleCommit(r.T, salt, voter, nextReports)

	if o.currentRound > 0 {
		_, err = r.Oracle.Vote(tests.FromSender(voter, common.Big0), commit, currentReports, salt, 0)
	} else {
		_, err = r.Oracle.Vote(tests.FromSender(voter, common.Big0), commit, nil, salt, 0)
	}

	require.NoError(r.T, err)
	r.WaitNBlocks(int(config.VotePeriod.Int64()))

	if o.currentRound > 0 {
		_, err = r.Acu.Update(tests.FromAutonity)
		require.NoError(r.T, err)
	}

	o.currentRound++
}
