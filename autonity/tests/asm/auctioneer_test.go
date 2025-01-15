package asm

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/autonity/tests"
)

func TestAuctioneerInterestAuction(t *testing.T) {
	setup := func() *tests.Runner {
		return tests.Setup(t, nil)
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
}
