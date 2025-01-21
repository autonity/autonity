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
		return r
	}

	tests.RunWithSetup("Paying interest above threshold starts an auction", setup, func(r *tests.Runner) {
		// Verify no initial open auctions
		auctions, _, err := r.Auctioneer.OpenAuctions(nil)
		require.NoError(t, err)
		require.Len(t, auctions, 0)

		config, _, err := r.Auctioneer.Config(nil)
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

		config, _, err := r.Auctioneer.Config(nil)
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
		require.Equal(t, round.Round.Int64(), auction.StartRound.Int64())
		require.Equal(t, auction.StartTimestamp, r.Evm.Context.Time)
		require.Equal(t, auctionAmount, auction.Amount)

		oracleData, _, err := r.Oracle.GetRoundData(nil, auction.StartRound, "NTN-ATN")
		require.NoError(t, err)
		require.Equal(t, newtonPrice, oracleData.Price)

		ntnCost, _, err := r.Auctioneer.MinInterestPayment(nil, auction.Id)
		require.NoError(t, err)

		config, _, err := r.Auctioneer.Config(nil)
		require.NoError(t, err)

		priceDiscount := newFloat0().Sub(
			newFloat(oracleData.Price),
			newFloat0().Quo(
				newFloat0().Mul(newFloat(oracleData.Price), newFloat(config.InterestAuctionDiscount)),
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
		config, _, err := r.Auctioneer.Config(nil)
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

func setupInterestAuction(r *tests.Runner, atnAboveMininum *big.Int) (auctionAmount *big.Int) {
	or := newOracleTestRounds([]*big.Int{newtonPrice, newtonPrice, newtonPrice, newtonPrice})
	or.initialize(r)
	or.increment(r)
	or.increment(r)

	// start an auction
	config, _, err := r.Auctioneer.Config(nil)
	require.NoError(r.T, err)

	auctionAmount = new(big.Int).Add(config.InterestAuctionThreshold, atnAboveMininum)
	// start an interest auction
	r.GiveMeSomeMoney(r.Stabilization.Address(), auctionAmount)
	_, err = r.Auctioneer.PaidInterest(tests.FromSender(r.Stabilization.Address(), auctionAmount))
	require.NoError(r.T, err)

	return auctionAmount
}

// testOracleRounds is a helper struct for testing the auctioneer with an oracle
type testOracleRounds struct {
	currentRound int
	symbolPrices map[string][]*big.Int
}

func newOracleTestRounds(ntnPrices []*big.Int) *testOracleRounds {
	or := &testOracleRounds{
		currentRound: 0,
		symbolPrices: make(map[string][]*big.Int),
	}
	or.symbolPrices["NTN-ATN"] = ntnPrices
	return or
}

func (o *testOracleRounds) initialize(r *tests.Runner) {
	config, _, err := r.Oracle.Config(nil)
	require.NoError(r.T, err)
	var symbols []string
	for symbol := range o.symbolPrices {
		symbols = append(symbols, symbol)
	}
	_, err = r.Oracle.SetSymbols(r.Operator, symbols)
	require.NoError(r.T, err)
	r.WaitNBlocks(2 * int(config.VotePeriod.Int64()))
}

func (o *testOracleRounds) increment(r *tests.Runner) {
	config, _, err := r.Oracle.Config(nil)
	require.NoError(r.T, err)
	voter := r.Committee.Validators[0].OracleAddress

	var nextReports []tests.IOracleReport
	var currentReports []tests.IOracleReport

	for _, prices := range o.symbolPrices {
		if o.currentRound > 0 {
			currentReports = append(
				currentReports,
				tests.IOracleReport{
					Price:      prices[o.currentRound%len(prices)],
					Confidence: 100,
				},
			)
		}
		nextReports = append(
			nextReports,
			tests.IOracleReport{
				Price:      prices[(o.currentRound+1)%len(prices)],
				Confidence: 100,
			},
		)
	}

	salt := big.NewInt(1234)
	commit := tests.MakeOracleCommit(r.T, salt, voter, nextReports)

	_, err = r.Oracle.Vote(tests.FromSender(voter, common.Big0), commit, currentReports, salt, 0)
	require.NoError(r.T, err)

	r.WaitNBlocks(int(config.VotePeriod.Int64()))
	o.currentRound++
}
