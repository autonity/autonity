package tests

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/params"
	"github.com/stretchr/testify/require"
)

var fromAutonity = &runOptions{origin: params.AutonityContractAddress}

func createSchedule(r *runner, beneficiary common.Address, amount, startBlock, cliffBlock, endBlock int64) {
	amountBigInt := big.NewInt(amount)
	_, err := r.autonity.Mint(operator, r.vestingManager.address, amountBigInt)
	require.NoError(r.t, err)
	_, err = r.vestingManager.NewSchedule(
		operator, beneficiary, big.NewInt(amount), big.NewInt(startBlock),
		big.NewInt(cliffBlock), big.NewInt(endBlock), true,
	)
	require.NoError(r.t, err)
}

func fromSender(sender common.Address, value *big.Int) *runOptions {
	return &runOptions{origin: sender, value: value}
}

func bondAndApply(r *runner, validatorAddress common.Address, bondingID int, scheduleID, bondingAmount *big.Int) (uint64, uint64) {
	validator, _, err := r.autonity.GetValidator(nil, validatorAddress)
	require.NoError(r.t, err)
	value, _, err := r.vestingManager.RequiredGasCostBond(nil)
	require.NoError(r.t, err)
	liquid, _, err := r.vestingManager.LiquidBalanceOf(nil, user, common.Big0, validator.NodeAddress)
	require.NoError(r.t, err)
	_, err = r.vestingManager.Bond(fromSender(user, value), common.Big0, validator.NodeAddress, bondingAmount)
	require.NoError(r.t, err)
	abi, err := LiquidMetaData.GetAbi()
	require.NoError(r.t, err)
	liquidContract := &Liquid{&contract{validator.LiquidContract, abi, r}}
	reward := big.NewInt(1000)
	_, err = liquidContract.Redistribute(fromSender(r.autonity.address, reward))
	require.NoError(r.t, err)
	bondedValidators := make([]common.Address, 1)
	bondedValidators[0] = validator.NodeAddress
	gasUsedDistribute, err := r.vestingManager.RewardsDistributed(fromAutonity, bondedValidators)
	require.NoError(r.t, err)
	_, err = liquidContract.Mint(fromAutonity, r.vestingManager.address, bondingAmount)
	require.NoError(r.t, err)
	gasUsedBond, err := r.vestingManager.BondingApplied(
		fromAutonity, big.NewInt(int64(bondingID)), validator.NodeAddress, bondingAmount, false, false,
	)
	require.NoError(r.t, err)
	newLiquid, _, err := r.vestingManager.LiquidBalanceOf(nil, user, common.Big0, validator.NodeAddress)
	require.NoError(r.t, err)
	require.Equal(r.t, new(big.Int).Add(liquid, bondingAmount), newLiquid)
	return gasUsedDistribute, gasUsedBond
}

func TestGasConsumption(t *testing.T) {
	r := setup(t, nil)
	var amount int64 = 1000
	scheduleCount := 10
	for i := 0; i < scheduleCount; i++ {
		createSchedule(r, user, amount, 0, 0, 1000)
	}
	committee, _, err := r.autonity.GetCommittee(nil)
	require.NoError(r.t, err)
	validator := committee[0].Addr
	validators, _, err := r.autonity.GetValidators(nil)
	require.NoError(r.t, err)

	r.run("measure gas for bond", func(r *runner) {
		bondingID := len(validators)
		bondingAmount := big.NewInt(amount)
		_, err := r.autonity.Mint(operator, user, bondingAmount)
		require.NoError(r.t, err)
		_, err = r.autonity.Bond(fromSender(user, nil), validator, bondingAmount)
		require.NoError(r.t, err)
		bondingID++
		r.waitNextEpoch()
		bondingAmount = big.NewInt(amount / 10)
		for iteration := 10; iteration > 0; iteration-- {
			gasUsedDistribute, gasUsedBond := bondAndApply(r, validator, bondingID, common.Big0, bondingAmount)
			fmt.Printf("reward distribution notification gas : %v\n", gasUsedDistribute)
			fmt.Printf("bonding applied notifaction gas : %v\n", gasUsedBond)
			fmt.Printf("total gas used: %v\n", gasUsedBond+gasUsedDistribute)
			bondingID++
		}
	})

	r.run("multiple bonding", func(r *runner) {
		bondingID := len(validators)
		bondingAmount := big.NewInt(amount)
		value, _, err := r.vestingManager.RequiredGasCostBond(nil)
		require.NoError(r.t, err)
		for i := 1; i < scheduleCount; i++ {
			_, err := r.vestingManager.Bond(fromSender(user, value), big.NewInt(int64(i)), validator, bondingAmount)
			require.NoError(r.t, err)
			bondingID++
		}
		r.waitNextEpoch()
		validatorInfo, _, err := r.autonity.GetValidator(nil, validator)
		require.NoError(r.t, err)
		delegatedStake := new(big.Int).Sub(validatorInfo.BondedStake, validatorInfo.SelfBondedStake)
		require.Equal(r.t, big.NewInt(amount*(int64(scheduleCount)-1)), delegatedStake)
		abi, err := LiquidMetaData.GetAbi()
		require.NoError(r.t, err)
		liquidContract := &Liquid{&contract{validatorInfo.LiquidContract, abi, r}}
		liquidBalance, _, err := liquidContract.BalanceOf(nil, r.vestingManager.address)
		require.NoError(r.t, err)
		require.Equal(r.t, delegatedStake, liquidBalance)
		gasUsedDistribute, gasUsedBond := bondAndApply(r, validator, bondingID, common.Big0, bondingAmount)
		fmt.Printf("reward distribution notification gas : %v\n", gasUsedDistribute)
		fmt.Printf("bonding applied notifaction gas : %v\n", gasUsedBond)
		fmt.Printf("total gas used: %v\n", gasUsedBond+gasUsedDistribute)
	})
}
