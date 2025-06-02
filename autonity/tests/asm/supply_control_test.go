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

func TestSupplyControlConstructor(t *testing.T) {
	totalSupply := e18
	setup := func() *tests.Runner {
		r := tests.Setup(t, nil)
		r.GiveMeSomeMoney(params.DeployerAddress, totalSupply)
		return r
	}

	tests.RunWithSetup("Test supply control deploys with total supply", setup, func(r *tests.Runner) {
		_, _, supplyControl, err := r.DeploySupplyControl(
			tests.FromSender(params.DeployerAddress, totalSupply),
			r.Autonity.Address(),
			params.TestAutonityContractConfig.Operator,
			r.Stabilization.Address(),
		)
		require.NoError(t, err)

		// Check total supply
		actualTotalSupply, _, err := supplyControl.GetTotalSupply(nil)
		require.NoError(t, err)

		require.Equal(t, totalSupply, actualTotalSupply)

		actualAvailableSupply, _, err := supplyControl.AvailableSupply(nil)
		require.NoError(t, err)

		require.Equal(t, totalSupply, actualAvailableSupply)
	})

	tests.RunWithSetup("Test supply control deploy fails on zero total supply", setup, func(r *tests.Runner) {
		_, _, _, err := r.DeploySupplyControl(
			tests.FromSender(params.DeployerAddress, common.Big0),
			r.Autonity.Address(),
			params.TestAutonityContractConfig.Operator,
			r.Stabilization.Address(),
		)
		require.Error(t, err)
		require.ErrorAs(t, err, &tests.SupplyControlZeroValueError{})
	})
}

func TestSupplyControlAuthorization(t *testing.T) {
	totalSupply := e18
	setup := func() *tests.Runner {
		r := tests.Setup(t, nil)
		r.GiveMeSomeMoney(params.DeployerAddress, totalSupply)
		return r
	}

	tests.RunWithSetup("Test mint cannot be called from non-stabilizer addresses", setup, func(r *tests.Runner) {
		unauthorizedAccounts := []common.Address{
			params.DeployerAddress,
			params.TestAutonityContractConfig.Operator,
			testrand.Address(),
		}
		for _, unauthorizedAccount := range unauthorizedAccounts {
			_, err := r.SupplyControl.Mint(tests.FromSender(unauthorizedAccount, nil), testrand.Address(), common.Big1)
			require.ErrorAs(t, err, &tests.SupplyControlUnauthorizedError{})
		}
	})

	tests.RunWithSetup("Test mint can be called from stabilizer address", setup, func(r *tests.Runner) {
		_, err := r.SupplyControl.Mint(tests.FromSender(r.Stabilization.Address(), nil), testrand.Address(), common.Big1)
		require.NoError(t, err)
	})

	tests.RunWithSetup("Test burn cannot be called from non-stabilizer addresses", setup, func(r *tests.Runner) {
		unauthorizedAccounts := []common.Address{
			params.DeployerAddress,
			params.TestAutonityContractConfig.Operator,
			testrand.Address(),
		}
		for _, unauthorizedAccount := range unauthorizedAccounts {
			r.GiveMeSomeMoney(unauthorizedAccount, common.Big32)
			_, err := r.SupplyControl.Burn(tests.FromSender(unauthorizedAccount, common.Big1))
			require.ErrorAs(t, err, &tests.SupplyControlUnauthorizedError{})
		}
	})

	tests.RunWithSetup("Test burn can be called from stabilizer address", setup, func(r *tests.Runner) {
		r.GiveMeSomeMoney(r.Stabilization.Address(), common.Big32)
		_, err := r.SupplyControl.Burn(tests.FromSender(r.Stabilization.Address(), common.Big1))
		require.NoError(t, err)
	})

	tests.RunWithSetup("Test set stabilizer can be called by autonity", setup, func(r *tests.Runner) {
		_, err := r.SupplyControl.SetStabilizer(
			tests.FromAutonity,
			testrand.Address(),
		)
		require.NoError(t, err)
	})

	tests.RunWithSetup("Test set stabilizer cannot be called by non-autonity", setup, func(r *tests.Runner) {
		unauthorizedAccounts := []common.Address{
			params.DeployerAddress,
			params.TestAutonityContractConfig.Operator,
			testrand.Address(),
		}
		for _, unauthorizedAccount := range unauthorizedAccounts {
			_, err := r.SupplyControl.SetStabilizer(
				tests.FromSender(unauthorizedAccount, nil),
				testrand.Address(),
			)
			require.ErrorAs(t, err, &tests.SupplyControlUnauthorizedError{})
		}
	})

	tests.RunWithSetup("Test set operator can be called by autonity", setup, func(r *tests.Runner) {
		_, err := r.SupplyControl.SetOperator(
			tests.FromSender(r.Autonity.Address(), nil),
			testrand.Address(),
		)
		require.NoError(t, err)
	})

	tests.RunWithSetup("Test set operator cannot be called by non-autonity", setup, func(r *tests.Runner) {
		unauthorizedAccounts := []common.Address{
			params.DeployerAddress,
			r.Stabilization.Address(),
			testrand.Address(),
		}
		for _, unauthorizedAccount := range unauthorizedAccounts {
			_, err := r.SupplyControl.SetOperator(
				tests.FromSender(unauthorizedAccount, nil),
				testrand.Address(),
			)
			require.ErrorAs(t, err, &tests.SupplyControlUnauthorizedError{})
		}
	})
}

func TestSupplyControlMintAndBurn(t *testing.T) {
	setup := func() *tests.Runner {
		r := tests.Setup(t, nil)
		return r
	}

	tests.RunWithSetup("Test mint to invalid recipient", setup, func(r *tests.Runner) {
		_, err := r.SupplyControl.Mint(tests.FromSender(r.Stabilization.Address(), nil), common.Address{}, common.Big1)
		require.ErrorAs(t, err, &tests.SupplyControlInvalidRecipientError{})

		_, err = r.SupplyControl.Mint(tests.FromSender(r.Stabilization.Address(), nil), r.Stabilization.Address(), common.Big1)
		require.ErrorAs(t, err, &tests.SupplyControlInvalidRecipientError{})
	})

	tests.RunWithSetup("Test mint invalid amount", setup, func(r *tests.Runner) {
		_, err := r.SupplyControl.Mint(tests.FromSender(r.Stabilization.Address(), nil), testrand.Address(), common.Big0)
		require.ErrorAs(t, err, &tests.SupplyControlInvalidAmountError{})

		totalSupply, _, err := r.SupplyControl.GetTotalSupply(nil)
		require.NoError(t, err)
		_, err = r.SupplyControl.Mint(
			tests.FromSender(r.Stabilization.Address(), nil),
			testrand.Address(),
			new(big.Int).Add(totalSupply, common.Big1),
		)
		require.ErrorAs(t, err, &tests.SupplyControlInvalidAmountError{})
	})

	tests.RunWithSetup("Test burn returns value to available supply", setup, func(r *tests.Runner) {
		toBurn := e18

		r.GiveMeSomeMoney(r.Stabilization.Address(), toBurn)
		balanceBefore := r.GetBalanceOf(r.Stabilization.Address())
		supplyControlBalanceBefore := r.GetBalanceOf(r.SupplyControl.Address())

		_, err := r.SupplyControl.Burn(tests.FromSender(r.Stabilization.Address(), toBurn))
		require.NoError(t, err)

		balanceAfter := r.GetBalanceOf(r.Stabilization.Address())
		supplyControlBalanceAfter := r.GetBalanceOf(r.SupplyControl.Address())
		require.Equal(t, toBurn, new(big.Int).Sub(balanceBefore, balanceAfter))
		require.Equal(t, toBurn, new(big.Int).Sub(supplyControlBalanceAfter, supplyControlBalanceBefore))
	})

	tests.RunWithSetup("Test mint decreases available supply", setup, func(r *tests.Runner) {
		toMint := e18
		user := testrand.Address()

		balanceBefore := r.GetBalanceOf(user)
		availableSupplyBefore, _, err := r.SupplyControl.AvailableSupply(nil)
		require.NoError(t, err)

		r.NoError(r.SupplyControl.Mint(tests.FromSender(r.Stabilization.Address(), nil), user, toMint))

		balanceAfter := r.GetBalanceOf(user)
		availableSupplyAfter, _, err := r.SupplyControl.AvailableSupply(nil)
		require.NoError(t, err)

		require.Equal(t, toMint, new(big.Int).Sub(balanceAfter, balanceBefore))
		require.Equal(t, toMint, new(big.Int).Sub(availableSupplyBefore, availableSupplyAfter))
	})

	// ToDo(scott): Test burn / mint events
}
