package autonitytests

import (
	"math/big"
	"testing"

	"github.com/autonity/autonity/autonity/tests"
	"github.com/stretchr/testify/require"
)

func TestCirculatingSupply(t *testing.T) {
	setup := func() *tests.Runner {
		return tests.Setup(t, nil)
	}

	tests.RunWithSetup("total supply and circulating supply are same", setup, func(r *tests.Runner) {
		totalSupply, _, err := r.Autonity.TotalSupply(nil)
		require.NoError(r.T, err)
		circulatingSupply, _, err := r.Autonity.CirculatingSupply(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, totalSupply, circulatingSupply)
	})

	tests.RunWithSetup("minting increases both circulating and total supply", setup, func(r *tests.Runner) {
		totalSupply, _, err := r.Autonity.TotalSupply(nil)
		require.NoError(r.T, err)
		circulatingSupply, _, err := r.Autonity.CirculatingSupply(nil)
		require.NoError(r.T, err)

		amount := big.NewInt(100)
		r.NoError(
			r.Autonity.Mint(r.Operator, tests.User, amount),
		)
		newTotalSupply, _, err := r.Autonity.TotalSupply(nil)
		require.NoError(r.T, err)
		newCirculatingSupply, _, err := r.Autonity.CirculatingSupply(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, new(big.Int).Add(totalSupply, amount), newTotalSupply)
		require.Equal(r.T, new(big.Int).Add(circulatingSupply, amount), newCirculatingSupply)
	})

	tests.RunWithSetup("burning decreases both circulating and total supply", setup, func(r *tests.Runner) {
		mintAmount := big.NewInt(100)
		r.NoError(
			r.Autonity.Mint(r.Operator, tests.User, mintAmount),
		)
		totalSupply, _, err := r.Autonity.TotalSupply(nil)
		require.NoError(r.T, err)
		circulatingSupply, _, err := r.Autonity.CirculatingSupply(nil)
		require.NoError(r.T, err)
		burnAmount := big.NewInt(60)
		r.NoError(
			r.Autonity.Burn(r.Operator, tests.User, burnAmount),
		)
		newTotalSupply, _, err := r.Autonity.TotalSupply(nil)
		require.NoError(r.T, err)
		newCirculatingSupply, _, err := r.Autonity.CirculatingSupply(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, new(big.Int).Sub(totalSupply, burnAmount), newTotalSupply)
		require.Equal(r.T, new(big.Int).Sub(circulatingSupply, burnAmount), newCirculatingSupply)
	})

}
