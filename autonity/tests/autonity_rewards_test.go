package tests

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/core"
)

func TestAutonityRewardsDistribution(t *testing.T) {
	r := setup(t, func(genesis *core.Genesis) *core.Genesis {
		// genesis.Config.AutonityContractConfig.EpochPeriod = 20
		return genesis
	})

	r.run("Test finalize with non-deployer account fails", func(rr *runner) {
		account := rr.randomAccount()
		_, err := r.autonity.Finalize(&runOptions{origin: account})
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted: function restricted to the protocol")
	})

	r.run("Test rewards distribution with only self bonded stake", func(rr *runner) {
		treasuryAccount := rr.genesis.Config.AutonityContractConfig.Treasury
		treasuryFee := new(big.Int).SetUint64(rr.genesis.Config.AutonityContractConfig.TreasuryFee)

		rewards := big.NewInt(1000000000000000)
		initFunds := rr.getBalanceOf(rr.autonity.address)
		require.Equal(t, int64(0), initFunds.Int64())

		// fund autonity contract
		rr.giveMeSomeMoney(rr.autonity.address, rewards)
		loadedBalance := rr.getBalanceOf(rr.autonity.address)
		require.Equal(t, rewards, loadedBalance)

		// get validators and initial ATN balance
		var initValidatorBalance []*big.Int
		for _, v := range rr.committee.validators {
			initValidatorBalance = append(initValidatorBalance, rr.getBalanceOf(v.Treasury))
		}

		initBalanceTreasury := rr.getBalanceOf(treasuryAccount)

		// finalize
		rr.waitNextEpoch()

		// check rewards
		expectedTreasuryRewards := new(big.Int).Div(
			new(big.Int).Mul(treasuryFee, rewards),
			big.NewInt(1e18),
		)
		afterTreasuryBalance := rr.getBalanceOf(treasuryAccount)
		require.Equal(t, expectedTreasuryRewards, new(big.Int).Sub(afterTreasuryBalance, initBalanceTreasury))

		// check validators rewards
		expectedValidatorRewards := new(big.Int).Sub(
			rewards,
			expectedTreasuryRewards,
		)
		totalStake := big.NewInt(0)
		for _, v := range rr.committee.validators {
			totalStake = new(big.Int).Add(totalStake, v.BondedStake)
		}

		for i, v := range rr.committee.validators {
			afterBalance := rr.getBalanceOf(v.Treasury)
			expectedReward := new(big.Int).Div(
				new(big.Int).Mul(v.BondedStake, expectedValidatorRewards),
				totalStake,
			)
			require.Equal(t, expectedReward, new(big.Int).Sub(afterBalance, initValidatorBalance[i]))
		}

	})

}
