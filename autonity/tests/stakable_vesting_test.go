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
	r.NoError(r.autonity.Mint(operator, r.stakableVesting.address, amountBigInt))
	// require.NoError(r.t, err)
	r.NoError(
		r.stakableVesting.NewSchedule(
			operator, beneficiary, big.NewInt(amount), big.NewInt(startTime),
			big.NewInt(cliffTime), big.NewInt(endTime),
		),
	)
	// require.NoError(r.t, err)
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
	r.NoError(r.stakableVesting.Bond(fromSender(user, bondingGas), scheduleID, validatorAddress, bondingAmount))
	// require.NoError(r.t, err)
	abi, err := LiquidMetaData.GetAbi()
	require.NoError(r.t, err)
	liquidContract := &Liquid{&contract{validator.LiquidContract, abi, r}}
	r.NoError(liquidContract.Redistribute(fromSender(r.autonity.address, reward)))
	// require.NoError(r.t, err)
	bondedValidators := make([]common.Address, 1)
	bondedValidators[0] = validatorAddress
	gasUsedDistribute := r.NoError(r.stakableVesting.RewardsDistributed(fromAutonity, bondedValidators))
	// require.NoError(r.t, err)
	if rejected == false {
		r.NoError(liquidContract.Mint(fromAutonity, r.stakableVesting.address, bondingAmount))
		// require.NoError(r.t, err)
		liquid = liquid.Add(liquid, bondingAmount)
	}
	gasUsedBond := r.NoError(
		r.stakableVesting.BondingApplied(
			fromAutonity, big.NewInt(int64(bondingID)), validatorAddress, bondingAmount, false, rejected,
		),
	)
	// require.NoError(r.t, err)
	newLiquid, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, scheduleID, validatorAddress)
	require.NoError(r.t, err)
	require.Equal(r.t, liquid, newLiquid)
	return gasUsedDistribute, gasUsedBond
}

func waitSomeEpoch(r *runner, endTime int64) int64 {
	require.True(r.t, r.evm.Context.Time.IsInt64(), "invalid data")
	currentTime := r.evm.Context.Time.Int64()
	for currentTime < endTime {
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
	r.NoError(r.stakableVesting.Unbond(fromSender(user, unbondingGas), scheduleID, validatorAddress, unbondingAmount))
	// require.NoError(r.t, err)
	abi, err := LiquidMetaData.GetAbi()
	require.NoError(r.t, err)
	validator, _, err := r.autonity.GetValidator(nil, validatorAddress)
	require.NoError(r.t, err)
	liquidContract := &Liquid{&contract{validator.LiquidContract, abi, r}}
	r.NoError(liquidContract.Redistribute(fromSender(r.autonity.address, reward)))
	// require.NoError(r.t, err)
	bondedValidators := make([]common.Address, 1)
	bondedValidators[0] = validatorAddress
	gasUsedDistribute := r.NoError(r.stakableVesting.RewardsDistributed(fromAutonity, bondedValidators))
	// require.NoError(r.t, err)
	if rejected == false {
		r.NoError(liquidContract.Unlock(fromAutonity, r.stakableVesting.address, unbondingAmount))
		// require.NoError(r.t, err)
		r.NoError(liquidContract.Burn(fromAutonity, r.stakableVesting.address, unbondingAmount))
		// require.NoError(r.t, err)
		liquid = liquid.Sub(liquid, unbondingAmount)
	}
	gasUsedUnbond := r.NoError(r.stakableVesting.UnbondingApplied(fromAutonity, big.NewInt(int64(unbondingID)), validatorAddress, rejected))
	// require.NoError(r.t, err)
	newLiquid, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, scheduleID, validatorAddress)
	require.NoError(r.t, err)
	require.Equal(r.t, liquid, newLiquid)
	gasUsedRelease := r.NoError(r.stakableVesting.UnbondingReleased(fromAutonity, big.NewInt(int64(unbondingID)), unbondingAmount, rejected))
	// require.NoError(r.t, err)
	return gasUsedDistribute, gasUsedUnbond, gasUsedRelease
}

func checkReleaseAllNTN(r *runner, user common.Address, scheduleID, unlockAmount *big.Int) {
	schedule, _, err := r.stakableVesting.GetSchedule(nil, user, scheduleID)
	require.NoError(r.t, err)
	scheduleNTN := schedule.CurrentNTNAmount
	initBalance, _, err := r.autonity.BalanceOf(nil, user)
	require.NoError(r.t, err)
	totalUnlocked, _, err := r.stakableVesting.UnlockedFunds(nil, user, scheduleID)
	require.NoError(r.t, err)
	require.True(r.t, unlockAmount.Cmp(totalUnlocked) == 0, "unlocked amount mismatch")
	r.NoError(r.stakableVesting.ReleaseAllNTN(fromSender(user, nil), scheduleID))
	// require.NoError(r.t, err)
	newBalance, _, err := r.autonity.BalanceOf(nil, user)
	require.NoError(r.t, err)
	require.Equal(r.t, new(big.Int).Add(initBalance, totalUnlocked), newBalance, "balance mismatch")
	schedule, _, err = r.stakableVesting.GetSchedule(nil, user, scheduleID)
	require.NoError(r.t, err)
	require.True(r.t, new(big.Int).Sub(scheduleNTN, unlockAmount).Cmp(schedule.CurrentNTNAmount) == 0, "schedule not updated properly")
}

func bondAndFinalize(r *runner, user, validator common.Address, liquidContract *Liquid, scheduleID, bondingAmount, bondingGas *big.Int) {
	liquidOfContract, _, err := liquidContract.BalanceOf(nil, r.stakableVesting.address)
	require.NoError(r.t, err)
	liquidOfUser, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, scheduleID, validator)
	require.NoError(r.t, err)
	schedule, _, err := r.stakableVesting.GetSchedule(nil, user, scheduleID)
	require.NoError(r.t, err)
	remaining := new(big.Int).Sub(schedule.CurrentNTNAmount, bondingAmount)
	r.NoError(
		r.stakableVesting.Bond(fromSender(user, bondingGas), scheduleID, validator, bondingAmount),
	)
	schedule, _, err = r.stakableVesting.GetSchedule(nil, user, scheduleID)
	require.NoError(r.t, err)
	require.Equal(r.t, remaining, schedule.CurrentNTNAmount, "schedule not updated properly")
	// let bonding apply
	r.waitNextEpoch()
	liquid, _, err := liquidContract.BalanceOf(nil, r.stakableVesting.address)
	require.NoError(r.t, err)
	require.Equal(r.t, new(big.Int).Add(bondingAmount, liquidOfContract), liquid, "liquid not minted")
	liquid, _, err = r.stakableVesting.LiquidBalanceOf(nil, user, scheduleID, validator)
	require.NoError(r.t, err)
	require.Equal(r.t, new(big.Int).Add(bondingAmount, liquidOfUser), liquid, "vesting contract cannot track liquid balance")
}

func TestBondingGasConsumption(t *testing.T) {
	r := setup(t, nil)
	var scheduleTotalAmount int64 = 1000
	scheduleCount := 10
	for i := 0; i < scheduleCount; i++ {
		createSchedule(r, user, scheduleTotalAmount, 0, 0, 1000)
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
	bondingAmount := big.NewInt(scheduleTotalAmount)
	r.NoError(
		r.autonity.Mint(operator, user, bondingAmount),
	)
	r.NoError(
		r.autonity.Bond(fromSender(user, nil), validator, bondingAmount),
	)
	r.waitNextEpoch()

	r.run("single bond", func(r *runner) {
		bondingID := len(validators) + 1
		var iteration int64 = 10
		bondingAmount := big.NewInt(scheduleTotalAmount / iteration)
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
		bondingAmount := big.NewInt(scheduleTotalAmount)
		for i := 1; i < scheduleCount; i++ {
			r.NoError(
				r.stakableVesting.Bond(fromSender(user, bondingGas), big.NewInt(int64(i)), validator, bondingAmount),
			)
			bondingID++
		}
		r.waitNextEpoch()
		delegatedStake := big.NewInt(scheduleTotalAmount * int64(scheduleCount-1))
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
		bondingAmount := big.NewInt(scheduleTotalAmount)
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
	var scheduleTotalAmount int64 = 1000
	scheduleCount := 10
	for i := 0; i < scheduleCount; i++ {
		createSchedule(r, user, scheduleTotalAmount, 0, 0, 1000)
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

	bondingAmount := big.NewInt(scheduleTotalAmount)
	for i := 0; i < scheduleCount; i++ {
		r.NoError(
			r.stakableVesting.Bond(fromSender(user, bondingGas), big.NewInt(int64(i)), validator, bondingAmount),
		)
	}
	r.waitNextEpoch()

	r.run("single unbond", func(r *runner) {
		var iteration int64 = 10
		unbondingAmount := big.NewInt(scheduleTotalAmount / iteration)
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
		unbondingAmount := big.NewInt(scheduleTotalAmount)
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
	var scheduleTotalAmount int64 = 1000
	var start int64 = 100
	var cliff int64 = 500
	// by making (end - start == scheduleTotalAmount) we have (totalUnlocked = currentTime - start)
	var end int64 = scheduleTotalAmount + start
	createSchedule(r, user, scheduleTotalAmount, start, cliff, end)
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
		currentTime := waitSomeEpoch(r, cliff)
		require.True(r.t, currentTime <= end, "release is not linear after end")
		unlocked := currentTime - start
		require.True(r.t, scheduleTotalAmount > unlocked, "cannot test if all funds unlocked")
		epochID, _, err := r.autonity.EpochID(nil)
		require.NoError(r.t, err)
		// mine some more blocks, release should be epoch based
		r.waitNBlocks(10)
		currentTime += 10
		checkReleaseAllNTN(r, user, scheduleID, big.NewInt(unlocked))

		r.waitNBlocks(10)
		currentTime += 10
		require.Equal(r.t, big.NewInt(currentTime), r.evm.Context.Time, "time mismatch, release won't work")
		// no more should be released as epoch did not change
		newEpochID, _, err := r.autonity.EpochID(nil)
		require.NoError(r.t, err)
		require.Equal(r.t, epochID, newEpochID, "cannot test if epoch progresses")
		checkReleaseAllNTN(r, user, scheduleID, common.Big0)
	})

	r.run("can release in chunks", func(r *runner) {
		currentTime := waitSomeEpoch(r, cliff)
		require.True(r.t, currentTime <= end, "cannot test, release is not linear after end")
		totalUnlocked, _, err := r.stakableVesting.UnlockedFunds(nil, user, scheduleID)
		require.NoError(r.t, err)
		require.True(r.t, totalUnlocked.IsInt64(), "invalid data")
		require.True(r.t, totalUnlocked.Int64() > 1, "cannot test chunks")
		unlockFraction := big.NewInt(totalUnlocked.Int64() / 2)
		// release only a chunk of total unlocked
		r.NoError(
			r.stakableVesting.ReleaseNTN(fromSender(user, nil), scheduleID, unlockFraction),
		)
		userNewBalance, _, err := r.autonity.BalanceOf(nil, user)
		require.NoError(r.t, err)
		require.Equal(r.t, new(big.Int).Add(userBalance, unlockFraction), userNewBalance, "balance mismatch")
		data, _, err := r.autonity.EpochID(nil)
		require.NoError(r.t, err)
		require.True(r.t, data.IsInt64(), "invalid data")
		epochID := data.Int64()
		r.waitNBlocks(10)
		data, _, err = r.autonity.EpochID(nil)
		require.NoError(r.t, err)
		require.True(r.t, data.IsInt64(), "invalid data")
		require.Equal(r.t, epochID, data.Int64(), "epoch progressed, more funds will release")
		require.True(r.t, r.evm.Context.Time.Cmp(big.NewInt(currentTime)) > 0, "time did not progress")
		checkReleaseAllNTN(r, user, scheduleID, new(big.Int).Sub(totalUnlocked, unlockFraction))
	})

	r.run("cannot release more than total", func(r *runner) {
		waitSomeEpoch(r, end)
		// progress some more epoch, should not matter after end
		r.waitNextEpoch()
		currentTime := r.evm.Context.Time
		checkReleaseAllNTN(r, user, scheduleID, big.NewInt(scheduleTotalAmount))
		r.waitNextEpoch()
		require.True(r.t, r.evm.Context.Time.Cmp(currentTime) > 0, "time did not progress")
		// cannot release more
		checkReleaseAllNTN(r, user, scheduleID, common.Big0)
	})
}

func TestStaking(t *testing.T) {
	r := setup(t, nil)
	var scheduleTotalAmount int64 = 1000
	var start int64 = 100
	var cliff int64 = 500
	// by making (end - start == scheduleTotalAmount) we have (totalUnlocked = currentTime - start)
	var end int64 = scheduleTotalAmount + start

	// create multiple schedules for multiple user
	users := make([]common.Address, 2)
	users[0] = user
	users[1] = common.HexToAddress("0x88")
	require.NotEqual(r.t, users[0], users[1], "same user")
	totalSchedules := 2
	for _, user := range users {
		for i := 0; i < totalSchedules; i++ {
			createSchedule(r, user, scheduleTotalAmount, start, cliff, end)
		}
	}

	// use multiple validators
	bondingGas, _, err := r.stakableVesting.RequiredBondingGasCost(nil)
	require.NoError(r.t, err)
	committee, _, err := r.autonity.GetCommittee(nil)
	require.NoError(r.t, err)
	validatorCount := 2
	validators := make([]common.Address, validatorCount)
	liquidContracts := make([]*Liquid, validatorCount)
	abi, err := LiquidMetaData.GetAbi()
	require.NoError(r.t, err)
	for i := 0; i < validatorCount; i++ {
		validators[i] = committee[i].Addr
		validatorInfo, _, err := r.autonity.GetValidator(nil, validators[i])
		require.NoError(r.t, err)
		liquidContracts[i] = &Liquid{&contract{validatorInfo.LiquidContract, abi, r}}
	}

	r.run("can bond all funds before cliff or start", func(r *runner) {
		require.True(r.t, r.evm.Context.Time.Cmp(big.NewInt(start)) < 0, "schedule started already")
		bondingAmount := big.NewInt(scheduleTotalAmount / 2)
		remaining := new(big.Int).Sub(big.NewInt(scheduleTotalAmount), bondingAmount)
		validator := validators[0]
		liquidContract := liquidContracts[0]
		user := users[0]
		scheduleID := common.Big0
		bondAndFinalize(r, user, validator, liquidContract, scheduleID, bondingAmount, bondingGas)
		require.True(r.t, r.evm.Context.Time.Cmp(big.NewInt(start)) < 0, "schedule started already")
		bondAndFinalize(r, user, validator, liquidContract, scheduleID, remaining, bondingGas)
	})

	r.run("cannot bond more than total", func(r *runner) {
		validator := validators[0]
		liquidContract := liquidContracts[0]
		user := users[0]
		scheduleID := common.Big0
		bondingAmount := big.NewInt(scheduleTotalAmount + 10)
		_, err := r.stakableVesting.Bond(fromSender(user, bondingGas), scheduleID, validator, bondingAmount)
		require.Equal(r.t, "execution reverted: not enough tokens", err.Error())
		bondingAmount = big.NewInt(scheduleTotalAmount / 2)
		remaining := new(big.Int).Sub(big.NewInt(scheduleTotalAmount), bondingAmount)
		r.NoError(
			r.stakableVesting.Bond(fromSender(user, bondingGas), scheduleID, validator, bondingAmount),
		)
		schedule, _, err := r.stakableVesting.GetSchedule(nil, user, scheduleID)
		require.NoError(r.t, err)
		require.Equal(r.t, remaining, schedule.CurrentNTNAmount, "schedule not updated properly")
		bondingAmount = new(big.Int).Add(big.NewInt(10), remaining)
		_, err = r.stakableVesting.Bond(fromSender(user, bondingGas), scheduleID, validator, bondingAmount)
		require.Equal(r.t, "execution reverted: not enough tokens", err.Error())
		// let bonding apply
		r.waitNextEpoch()
		_, err = r.stakableVesting.Bond(fromSender(user, bondingGas), scheduleID, validator, bondingAmount)
		require.Equal(r.t, "execution reverted: not enough tokens", err.Error())
		bondAndFinalize(r, user, validator, liquidContract, scheduleID, remaining, bondingGas)
	})

	r.run("can release liquid tokens", func(r *runner) {
		bondingAmount := big.NewInt(scheduleTotalAmount)
		validator := validators[0]
		liquidContract := liquidContracts[0]
		scheduleID := common.Big0
		r.NoError(
			r.stakableVesting.Bond(fromSender(user, bondingGas), scheduleID, validator, bondingAmount),
		)
		// let bonding apply
		r.waitNextEpoch()
		currentTime := waitSomeEpoch(r, cliff)
		unlocked := currentTime - start
		// mine some more block, release should be epoch based
		r.waitNBlocks(10)
		r.NoError(
			r.stakableVesting.ReleaseAllLNTN(fromSender(user, nil), scheduleID),
		)
		liquid, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, scheduleID, validator)
		require.NoError(r.t, err)
		require.Equal(
			r.t, big.NewInt(scheduleTotalAmount-unlocked), liquid,
			"liquid release don't follow epoch based linear function",
		)
		liquid, _, err = liquidContract.BalanceOf(nil, r.stakableVesting.address)
		require.NoError(r.t, err)
		require.Equal(r.t, big.NewInt(scheduleTotalAmount-unlocked), liquid, "liquid not transferred")
		liquid, _, err = liquidContract.BalanceOf(nil, user)
		require.NoError(r.t, err)
		require.Equal(r.t, big.NewInt(unlocked), liquid, "liquid not received")
		waitSomeEpoch(r, end)
		// progress more epoch, shouldn't matter
		r.waitNextEpoch()
		r.NoError(
			r.stakableVesting.ReleaseAllLNTN(fromSender(user, nil), scheduleID),
		)
		liquid, _, err = r.stakableVesting.LiquidBalanceOf(nil, user, scheduleID, validator)
		require.NoError(r.t, err)
		require.True(r.t, liquid.Cmp(common.Big0) == 0, "all liquid tokens not released")
		liquid, _, err = liquidContract.BalanceOf(nil, r.stakableVesting.address)
		require.NoError(r.t, err)
		require.True(r.t, liquid.Cmp(common.Big0) == 0, "liquid not transferred")
		liquid, _, err = liquidContract.BalanceOf(nil, user)
		require.NoError(r.t, err)
		require.Equal(r.t, big.NewInt(scheduleTotalAmount), liquid, "liquid not received")
	})

	r.run("track liquids when bonding from mutliple schedule to multiple validators", func(r *runner) {
		// TODO: complete
	})

	r.run("can unbond", func(r *runner) {
		// TODO: complete
	})

	r.run("cannot unbond more than total liquid", func(r *runner) {
		// TODO: complete
	})

	r.run("cannot unbond if released", func(r *runner) {
		// TODO: complete
	})

	r.run("when bonded, release NTN first", func(r *runner) {
		// TODO: complete
	})

	r.run("test release when bonding to multiple validator", func(r *runner) {
		// TODO: complete
	})
}

func TestStakingRevert(t *testing.T) {

}

func TestAUTRwardTracking(t *testing.T) {

}

func TestScheduleCancel(t *testing.T) {

}
