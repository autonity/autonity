package tests

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
)

func TestRestrictedFunctionAccess(t *testing.T) {
	testUser := common.Address{200}

	setup := func() *Runner {
		r := Setup(t, nil)
		r.GiveMeSomeMoney(testUser, new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
		return r
	}

	RunWithSetup("Set ATN supply operator restricted to operator", setup, func(r *Runner) {
		_, err := r.Stabilization.SetAtnSupplyOperator(FromSender(testUser, common.Big0), testUser)
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted")

		_, err = r.Stabilization.SetAtnSupplyOperator(r.Operator, testUser)
		require.NoError(t, err)
	})

	RunWithSetup("Deposit restricted to ATN supply operator", setup, func(r *Runner) {
		_, err := r.Autonity.Mint(r.Operator, testUser, big.NewInt(1000))
		require.NoError(t, err)

		_, err = r.Autonity.Approve(FromSender(testUser, common.Big0), r.Stabilization.Address(), common.Big256)
		require.NoError(t, err)

		// test user is not the ATN supply operator
		_, err = r.Stabilization.Deposit(FromSender(testUser, common.Big0), big.NewInt(1000))
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted")

		// set test user as the ATN supply operator
		_, err = r.Stabilization.SetAtnSupplyOperator(r.Operator, testUser)
		require.NoError(t, err)

		// test user is now the ATN supply operator
		_, err = r.Stabilization.Deposit(FromSender(testUser, common.Big0), big.NewInt(10))
		require.NoError(t, err)
	})

	RunWithSetup("removeCDPRestrictions restricted to operator", setup, func(r *Runner) {
		_, err := r.Stabilization.RemoveCDPRestrictions(FromSender(testUser, common.Big0))
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted")

		_, err = r.Stabilization.RemoveCDPRestrictions(r.Operator)
		require.NoError(t, err)
	})
}

func TestInterestRate(t *testing.T) {
	atnSupplyOperator := common.Address{200}
	setup := func() *Runner {
		r := Setup(t, nil)
		r.GiveMeSomeMoney(atnSupplyOperator, big.NewInt(100000000000000000))
		_, err := r.Stabilization.SetAtnSupplyOperator(r.Operator, atnSupplyOperator)
		require.NoError(t, err)
		return r
	}

	RunWithSetup("Interest rate should be zero before restrictions are removed", setup, func(r *Runner) {
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
