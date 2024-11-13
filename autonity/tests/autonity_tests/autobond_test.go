package autonitytests

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/autonity/tests"
	"github.com/autonity/autonity/common"
)

func waitForEndEpochBlock(rr *tests.Runner) {
	info, _, err := rr.Autonity.GetEpochInfo(nil)
	require.NoError(rr.T, err)
	diff := new(big.Int).Sub(info.NextEpochBlock, rr.Evm.Context.BlockNumber)
	rr.WaitNBlocks(int(diff.Uint64()))
}

func TestAutobondPermissions(t *testing.T) {
	setup := func() *tests.Runner {
		return tests.Setup(t, nil)
	}

	tests.RunWithSetup("autobond cannot be called by non-rewards distributor", setup, func(r *tests.Runner) {
		validator := r.Committee.Validators[0]
		waitForEndEpochBlock(r)
		_, err := r.Autonity.Autobond(
			tests.FromSender(validator.NodeAddress, nil),
			validator.NodeAddress,
			common.Big5,
			common.Big0,
		)
		require.Error(r.T, err)
		require.Contains(r.T, err.Error(), "caller is not a reward distributor")
	})

	tests.RunWithSetup("autobond cannot be called until the end of an epoch", setup, func(r *tests.Runner) {
		validator := r.Committee.Validators[0]
		_, err := r.Autonity.Mint(r.Operator, r.Accountability.Address(), common.Big32)
		require.NoError(r.T, err)

		_, err = r.Autonity.Autobond(
			tests.FromSender(r.Accountability.Address(), nil),
			validator.NodeAddress,
			common.Big5,
			common.Big0,
		)
		require.Error(r.T, err)
		require.Contains(r.T, err.Error(), "epoch is not ended")
	})

	tests.RunWithSetup("autobond can be called by all rewards distributing contracts at epoch end", setup, func(r *tests.Runner) {
		validator := r.Committee.Validators[0]
		liquidContract := r.LiquidStateContract(validator.NodeAddress)

		_, err := r.Autonity.Mint(r.Operator, r.Accountability.Address(), common.Big32)
		require.NoError(r.T, err)
		_, err = r.Autonity.Mint(r.Operator, r.OmissionAccountability.Address(), common.Big32)
		require.NoError(r.T, err)
		_, err = r.Autonity.Mint(r.Operator, liquidContract.Address(), common.Big32)
		require.NoError(r.T, err)

		waitForEndEpochBlock(r)
		_, err = r.Autonity.Autobond(
			tests.FromSender(r.Accountability.Address(), nil),
			validator.NodeAddress,
			common.Big5,
			common.Big0,
		)
		require.NoError(r.T, err)

		_, err = r.Autonity.Autobond(
			tests.FromSender(r.OmissionAccountability.Address(), nil),
			validator.NodeAddress,
			common.Big5,
			common.Big1,
		)
		require.NoError(r.T, err)

		_, err = r.Autonity.Autobond(
			tests.FromSender(liquidContract.Address(), nil),
			validator.NodeAddress,
			common.Big5,
			common.Big2,
		)
		require.NoError(r.T, err)
	})
}

func TestBondedIncrease(t *testing.T) {
	setup := func() *tests.Runner {
		return tests.Setup(t, nil)
	}
	tests.RunWithSetup("self-bonded increases after autobond", setup, func(r *tests.Runner) {
		validator := r.Committee.Validators[0]

		val, _, err := r.Autonity.GetValidator(nil, validator.NodeAddress)
		require.NoError(r.T, err)
		oldSelfBonded := val.SelfBondedStake
		oldTotalBonded := val.BondedStake

		_, err = r.Autonity.Mint(r.Operator, r.OmissionAccountability.Address(), common.Big32)
		require.NoError(r.T, err)

		waitForEndEpochBlock(r)
		_, err = r.Autonity.Autobond(
			tests.FromSender(r.OmissionAccountability.Address(), nil),
			validator.NodeAddress,
			common.Big5,
			common.Big0,
		)
		require.NoError(r.T, err)

		val, _, err = r.Autonity.GetValidator(nil, validator.NodeAddress)
		require.NoError(r.T, err)
		newSelfBonded := val.SelfBondedStake
		newTotalBonded := val.BondedStake

		require.Equal(r.T, new(big.Int).Add(oldSelfBonded, common.Big5), newSelfBonded)
		require.Equal(r.T, new(big.Int).Add(oldTotalBonded, common.Big5), newTotalBonded)
	})

	tests.RunWithSetup("total-bonded increases after autobond", setup, func(r *tests.Runner) {
		validator := r.Committee.Validators[0]

		val, _, err := r.Autonity.GetValidator(nil, validator.NodeAddress)
		require.NoError(r.T, err)
		oldSelfBonded := val.SelfBondedStake
		oldTotalBonded := val.BondedStake

		_, err = r.Autonity.Mint(r.Operator, r.OmissionAccountability.Address(), common.Big32)
		require.NoError(r.T, err)

		waitForEndEpochBlock(r)
		_, err = r.Autonity.Autobond(
			tests.FromSender(r.OmissionAccountability.Address(), nil),
			validator.NodeAddress,
			common.Big0,
			common.Big5,
		)
		require.NoError(r.T, err)

		val, _, err = r.Autonity.GetValidator(nil, validator.NodeAddress)
		require.NoError(r.T, err)
		newSelfBonded := val.SelfBondedStake
		newTotalBonded := val.BondedStake

		require.Equal(r.T, oldSelfBonded, newSelfBonded)
		require.Equal(r.T, new(big.Int).Add(oldTotalBonded, common.Big5), newTotalBonded)
	})
}

func TestOnlyAutobondedRewards(t *testing.T) {
	setup := func() *tests.Runner {
		return tests.Setup(t, nil)
	}

	tests.RunWithSetup("autobonded rewards are the only rewards", setup, func(r *tests.Runner) {
		validator := r.Committee.Validators[0]
		balanceBefore, _, err := r.Autonity.BalanceOf(nil, validator.NodeAddress)
		require.NoError(r.T, err)

		val, _, err := r.Autonity.GetValidator(nil, validator.NodeAddress)
		require.NoError(r.T, err)
		stakeBefore := val.BondedStake

		// wait some epochs
		r.WaitNextEpoch()
		r.WaitNextEpoch()
		r.WaitNextEpoch()

		val, _, err = r.Autonity.GetValidator(nil, validator.NodeAddress)
		require.NoError(r.T, err)
		stakeAfter := val.BondedStake

		balanceAfter, _, err := r.Autonity.BalanceOf(nil, validator.NodeAddress)
		require.NoError(r.T, err)

		require.True(r.T, stakeAfter.Cmp(stakeBefore) > 0, "stake should increase for validators")
		require.Equal(r.T, balanceBefore, balanceAfter, "balance should not increase for validators")
	})
}
