package tests

import (
	"math/big"
	"testing"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/params"
	"github.com/stretchr/testify/require"
)

var fromAutonity = &runOptions{origin: params.AutonityContractAddress}

var reward = big.NewInt(1000000000)

func createSchedule(r *runner, beneficiary common.Address, amount, startTime, cliffTime, endTime int64) {
	amountBigInt := big.NewInt(amount)
	_, err := r.autonity.Mint(operator, r.stakableVesting.address, amountBigInt)
	require.NoError(r.t, err)
	_, err = r.stakableVesting.NewSchedule(
		operator, beneficiary, big.NewInt(amount), big.NewInt(startTime),
		big.NewInt(cliffTime), big.NewInt(endTime),
	)
	require.NoError(r.t, err)
}

func fromSender(sender common.Address, value *big.Int) *runOptions {
	return &runOptions{origin: sender, value: value}
}

func bondAndApply(
	r *runner, validatorAddress, user common.Address, bondingID int, scheduleID, bondingAmount, bondingGas *big.Int, rejected bool,
) (uint64, uint64) {
	validator, _, err := r.autonity.GetValidator(nil, validatorAddress)
	require.NoError(r.t, err)
	liquid, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, scheduleID, validatorAddress)
	require.NoError(r.t, err)
	_, err = r.stakableVesting.Bond(fromSender(user, bondingGas), scheduleID, validatorAddress, bondingAmount)
	require.NoError(r.t, err)
	abi, err := LiquidMetaData.GetAbi()
	require.NoError(r.t, err)
	liquidContract := &Liquid{&contract{validator.LiquidContract, abi, r}}
	_, err = liquidContract.Redistribute(fromSender(r.autonity.address, reward))
	require.NoError(r.t, err)
	bondedValidators := make([]common.Address, 1)
	bondedValidators[0] = validatorAddress
	gasUsedDistribute, err := r.stakableVesting.RewardsDistributed(fromAutonity, bondedValidators)
	require.NoError(r.t, err)
	if rejected == false {
		_, err = liquidContract.Mint(fromAutonity, r.stakableVesting.address, bondingAmount)
		require.NoError(r.t, err)
		liquid = liquid.Add(liquid, bondingAmount)
	}
	gasUsedBond, err := r.stakableVesting.BondingApplied(
		fromAutonity, big.NewInt(int64(bondingID)), validatorAddress, bondingAmount, false, rejected,
	)
	require.NoError(r.t, err)
	newLiquid, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, scheduleID, validatorAddress)
	require.NoError(r.t, err)
	require.Equal(r.t, liquid, newLiquid)
	return gasUsedDistribute, gasUsedBond
}

func endCliffPeriod(r *runner, cliff int64) int64 {
	require.True(r.t, r.evm.Context.Time.IsInt64(), "invalid data")
	currentTime := r.evm.Context.Time.Int64()
	for currentTime < cliff {
		r.waitNextEpoch()
		currentTime = r.evm.Context.Time.Int64()
	}
	return currentTime
}

func unbondAndApply(
	r *runner, validatorAddress, user common.Address, unbondingID int, scheduleID, unbondingAmount, unbondingGas *big.Int, rejected bool,
) (uint64, uint64, uint64) {
	liquid, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, scheduleID, validatorAddress)
	require.NoError(r.t, err)
	_, err = r.stakableVesting.Unbond(fromSender(user, unbondingGas), scheduleID, validatorAddress, unbondingAmount)
	require.NoError(r.t, err)
	abi, err := LiquidMetaData.GetAbi()
	require.NoError(r.t, err)
	validator, _, err := r.autonity.GetValidator(nil, validatorAddress)
	require.NoError(r.t, err)
	liquidContract := &Liquid{&contract{validator.LiquidContract, abi, r}}
	_, err = liquidContract.Redistribute(fromSender(r.autonity.address, reward))
	require.NoError(r.t, err)
	bondedValidators := make([]common.Address, 1)
	bondedValidators[0] = validatorAddress
	gasUsedDistribute, err := r.stakableVesting.RewardsDistributed(fromAutonity, bondedValidators)
	require.NoError(r.t, err)
	if rejected == false {
		_, err = liquidContract.Unlock(fromAutonity, r.stakableVesting.address, unbondingAmount)
		require.NoError(r.t, err)
		_, err = liquidContract.Burn(fromAutonity, r.stakableVesting.address, unbondingAmount)
		require.NoError(r.t, err)
		liquid = liquid.Sub(liquid, unbondingAmount)
	}
	gasUsedUnbond, err := r.stakableVesting.UnbondingApplied(fromAutonity, big.NewInt(int64(unbondingID)), validatorAddress, rejected)
	require.NoError(r.t, err)
	newLiquid, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, scheduleID, validatorAddress)
	require.NoError(r.t, err)
	require.Equal(r.t, liquid, newLiquid)
	gasUsedRelease, err := r.stakableVesting.UnbondingReleased(fromAutonity, big.NewInt(int64(unbondingID)), unbondingAmount, rejected)
	require.NoError(r.t, err)
	return gasUsedDistribute, gasUsedUnbond, gasUsedRelease
}

func TestBondingGasConsumption(t *testing.T) {
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
	bondingGas, _, err := r.stakableVesting.RequiredBondingGasCost(nil)
	require.NoError(r.t, err)
	stakingGas, _, err := r.autonity.StakingGasPrice(nil)
	require.NoError(r.t, err)
	bondingAmount := big.NewInt(amount)
	_, err = r.autonity.Mint(operator, user, bondingAmount)
	require.NoError(r.t, err)
	_, err = r.autonity.Bond(fromSender(user, nil), validator, bondingAmount)
	require.NoError(r.t, err)
	r.waitNextEpoch()

	r.run("single bond", func(r *runner) {
		bondingID := len(validators) + 1
		var iteration int64 = 10
		bondingAmount := big.NewInt(amount / iteration)
		for ; iteration > 0; iteration-- {
			gasUsedDistribute, gasUsedBond := bondAndApply(r, validator, user, bondingID, common.Big0, bondingAmount, bondingGas, false)
			totalGasUsed := new(big.Int).Mul(big.NewInt(int64(gasUsedDistribute+gasUsedBond)), stakingGas)
			require.True(
				r.t,
				bondingGas.Cmp(totalGasUsed) >= 0,
				"need more gas to notify bonding operations",
			)
			bondingID++
		}
	})

	r.run("multiple bonding", func(r *runner) {
		validatorInfo, _, err := r.autonity.GetValidator(nil, validator)
		require.NoError(r.t, err)
		abi, err := LiquidMetaData.GetAbi()
		require.NoError(r.t, err)
		liquidContract := &Liquid{&contract{validatorInfo.LiquidContract, abi, r}}
		oldLiquidBalance, _, err := liquidContract.BalanceOf(nil, r.stakableVesting.address)
		require.NoError(r.t, err)
		bondingID := len(validators) + 1
		bondingAmount := big.NewInt(amount)
		for i := 1; i < scheduleCount; i++ {
			_, err := r.stakableVesting.Bond(fromSender(user, bondingGas), big.NewInt(int64(i)), validator, bondingAmount)
			require.NoError(r.t, err)
			bondingID++
		}
		r.waitNextEpoch()
		delegatedStake := big.NewInt(amount * int64(scheduleCount-1))
		liquidBalance, _, err := liquidContract.BalanceOf(nil, r.stakableVesting.address)
		require.NoError(r.t, err)
		require.Equal(r.t, oldLiquidBalance.Add(oldLiquidBalance, delegatedStake), liquidBalance)
		for i := 1; i < scheduleCount; i++ {
			liquidBalance, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, big.NewInt(int64(i)), validator)
			require.NoError(r.t, err)
			require.Equal(r.t, bondingAmount, liquidBalance)
		}
		gasUsedDistribute, gasUsedBond := bondAndApply(r, validator, user, bondingID, common.Big0, bondingAmount, bondingGas, false)
		totalGasUsed := new(big.Int).Mul(big.NewInt(int64(gasUsedDistribute+gasUsedBond)), stakingGas)
		require.True(
			r.t,
			bondingGas.Cmp(totalGasUsed) >= 0,
			"need more gas to notify bonding operations",
		)
	})

	r.run("bonding rejected", func(r *runner) {
		bondingID := len(validators) + 1
		bondingAmount := big.NewInt(amount)
		gasUsedDistribute, gasUsedBond := bondAndApply(r, validator, user, bondingID, common.Big0, bondingAmount, bondingGas, true)
		totalGasUsed := new(big.Int).Mul(big.NewInt(int64(gasUsedDistribute+gasUsedBond)), stakingGas)
		require.True(
			r.t,
			bondingGas.Cmp(totalGasUsed) >= 0,
			"need more gas to notify bonding operations",
		)
	})
}

func TestUnbondingGasConsumption(t *testing.T) {
	r := setup(t, nil)
	var amount int64 = 1000
	scheduleCount := 10
	for i := 0; i < scheduleCount; i++ {
		createSchedule(r, user, amount, 0, 0, 1000)
	}
	committee, _, err := r.autonity.GetCommittee(nil)
	require.NoError(r.t, err)
	validator := committee[0].Addr
	bondingGas, _, err := r.stakableVesting.RequiredBondingGasCost(nil)
	require.NoError(r.t, err)
	unbondingGas, _, err := r.stakableVesting.RequiredUnbondingGasCost(nil)
	require.NoError(r.t, err)
	stakingGas, _, err := r.autonity.StakingGasPrice(nil)
	require.NoError(r.t, err)

	bondingAmount := big.NewInt(amount)
	for i := 0; i < scheduleCount; i++ {
		_, err := r.stakableVesting.Bond(fromSender(user, bondingGas), big.NewInt(int64(i)), validator, bondingAmount)
		require.NoError(r.t, err)
	}
	r.waitNextEpoch()

	r.run("single unbond", func(r *runner) {
		var iteration int64 = 10
		unbondingAmount := big.NewInt(amount / iteration)
		unbondingID := 0
		for ; iteration > 0; iteration-- {
			gasUsedDistribute, gasUsedUnbond, gasUsedRelease := unbondAndApply(r, validator, user, unbondingID, common.Big0, unbondingAmount, unbondingGas, false)
			totalGasUsed := new(big.Int).Mul(big.NewInt(int64(gasUsedDistribute+gasUsedUnbond+gasUsedRelease)), stakingGas)
			require.True(
				r.t,
				unbondingGas.Cmp(totalGasUsed) >= 0,
				"need more gas to notify unbonding operations",
			)
			unbondingID++
		}
	})

	r.run("multiple unbond", func(r *runner) {
		//
	})

	r.run("unbond rejected", func(r *runner) {
		unbondingID := 0
		unbondingAmount := big.NewInt(amount)
		gasUsedDistribute, gasUsedUnbond, gasUsedRelease := unbondAndApply(r, validator, user, unbondingID, common.Big0, unbondingAmount, unbondingGas, true)
		totalGasUsed := new(big.Int).Mul(big.NewInt(int64(gasUsedDistribute+gasUsedUnbond+gasUsedRelease)), stakingGas)
		require.True(
			r.t,
			unbondingGas.Cmp(totalGasUsed) >= 0,
			"need more gas to notify unbonding operations",
		)
	})
}

func TestRelease(t *testing.T) {
	r := setup(t, nil)
	var amount int64 = 1000
	var start int64 = 100
	var cliff int64 = 1000
	// by making (end - start == amount) we have (totalUnlocked = currentTime - start)
	var end int64 = amount + start
	createSchedule(r, user, amount, start, cliff, end)
	scheduleID := common.Big0
	// do not modify userBalance
	userBalance, _, err := r.autonity.BalanceOf(nil, user)
	require.NoError(r.t, err)

	r.run("cannot release before cliff", func(r *runner) {
		r.waitNBlocks(int(cliff) - 1)
		require.Equal(r.t, big.NewInt(cliff-1), r.evm.Context.Time, "time mismatch")
		unlocked, _, err := r.stakableVesting.UnlockedFunds(nil, user, scheduleID)
		require.NoError(r.t, err)
		require.True(r.t, unlocked.Cmp(common.Big0) == 0, "funds unlocked before cliff period")
		_, err = r.stakableVesting.ReleaseFunds(fromSender(user, nil), scheduleID)
		require.Equal(r.t, "execution reverted: cliff period not reached yet", err.Error())
		userNewBalance, _, err := r.autonity.BalanceOf(nil, user)
		require.NoError(r.t, err)
		require.Equal(r.t, userBalance, userNewBalance, "funds released before cliff period")
	})

	r.run("release calculation follows epoch based linear function in time", func(r *runner) {
		currentTime := endCliffPeriod(r, cliff)
		require.True(r.t, currentTime <= end, "release is not linear after end")
		unlocked, _, err := r.stakableVesting.UnlockedFunds(nil, user, scheduleID)
		require.NoError(r.t, err)
		require.Equal(r.t, unlocked, big.NewInt(currentTime-start), "unlocked funds calculation not linear")
		_, err = r.stakableVesting.ReleaseFunds(fromSender(user, nil), scheduleID)
		require.NoError(r.t, err)
		userNewBalance, _, err := r.autonity.BalanceOf(nil, user)
		require.NoError(r.t, err)
		require.Equal(r.t, new(big.Int).Add(userBalance, unlocked), userNewBalance, "funds release not linear")
		currentTime++
		r.waitNBlocks(1)
		require.Equal(r.t, big.NewInt(currentTime), r.evm.Context.Time, "time mismatch")
		unlocked, _, err = r.stakableVesting.UnlockedFunds(nil, user, scheduleID)
		require.NoError(r.t, err)
		require.True(r.t, unlocked.Cmp(common.Big0) == 0, "unlocked funds calculation not epoch based")
		_, err = r.stakableVesting.ReleaseFunds(fromSender(user, nil), scheduleID)
		require.NoError(r.t, err)
		userLastBalance, _, err := r.autonity.BalanceOf(nil, user)
		require.NoError(r.t, err)
		require.Equal(r.t, userNewBalance, userLastBalance, "funds release not epoch based")
	})

	r.run("can release in chunks", func(r *runner) {
		currentTime := endCliffPeriod(r, cliff)
		require.True(r.t, currentTime <= end, "release is not linear after end")
		totalUnlocked, _, err := r.stakableVesting.UnlockedFunds(nil, user, scheduleID)
		require.NoError(r.t, err)
		require.True(r.t, totalUnlocked.IsInt64(), "invalid data")
		require.True(r.t, totalUnlocked.Int64() > 1, "cannot test chunks")
		unlockFraction := big.NewInt(totalUnlocked.Int64() / 2)
		// release only a chunk of total unlocked
		_, err = r.stakableVesting.ReleaseNTN(fromSender(user, nil), scheduleID, unlockFraction)
		require.NoError(r.t, err)
		userNewBalance, _, err := r.autonity.BalanceOf(nil, user)
		require.NoError(r.t, err)
		require.Equal(r.t, new(big.Int).Add(userBalance, unlockFraction), userNewBalance, "balance mismatch")
		data, _, err := r.autonity.EpochID(nil)
		require.NoError(r.t, err)
		require.True(r.t, data.IsInt64(), "invalid data")
		epochID := data.Int64()
		r.waitNBlocks(1)
		data, _, err = r.autonity.EpochID(nil)
		require.NoError(r.t, err)
		require.True(r.t, data.IsInt64(), "invalid data")
		require.Equal(r.t, epochID, data.Int64(), "epoch progressed, more funds will release")
		newUnlockFraction, _, err := r.stakableVesting.UnlockedFunds(nil, user, scheduleID)
		require.NoError(r.t, err)
		require.Equal(r.t, totalUnlocked, new(big.Int).Add(newUnlockFraction, unlockFraction), "unlocked more than it should")
		// release all
		_, err = r.stakableVesting.ReleaseAllNTN(fromSender(user, nil), scheduleID)
		require.NoError(r.t, err)
		userNewBalance, _, err = r.autonity.BalanceOf(nil, user)
		require.NoError(r.t, err)
		require.Equal(r.t, new(big.Int).Add(userBalance, totalUnlocked), userNewBalance, "balance mismatch")
	})
}
