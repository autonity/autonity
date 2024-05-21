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

var reward = big.NewInt(1000000000)

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
	var start int64 = 100 + r.evm.Context.Time.Int64()
	var cliff int64 = 500 + start
	// by making (end - start == scheduleTotalAmount) we have (totalUnlocked = currentTime - start)
	end := scheduleTotalAmount + start
	createSchedule(r, user, scheduleTotalAmount, start, cliff, end)
	scheduleID := common.Big0
	// do not modify userBalance
	userBalance, _, err := r.autonity.BalanceOf(nil, user)
	require.NoError(r.t, err)

	r.run("cannot release before cliff", func(r *runner) {
		r.waitSomeBlock(cliff)
		require.Equal(r.t, big.NewInt(cliff), r.evm.Context.Time, "time mismatch")
		_, _, err := r.stakableVesting.UnlockedFunds(nil, user, scheduleID)
		require.Equal(r.t, "execution reverted: cliff period not reached yet", err.Error())
		_, err = r.stakableVesting.ReleaseFunds(fromSender(user, nil), scheduleID)
		require.Equal(r.t, "execution reverted: cliff period not reached yet", err.Error())
		userNewBalance, _, err := r.autonity.BalanceOf(nil, user)
		require.NoError(r.t, err)
		require.Equal(r.t, userBalance, userNewBalance, "funds released before cliff period")
	})

	r.run("release calculation follows epoch based linear function in time", func(r *runner) {
		currentTime := r.waitSomeEpoch(cliff + 1)
		require.True(r.t, currentTime <= end, "release is not linear after end")
		// contract has the context of last block, so time is 1s less than currentTime
		unlocked := currentTime - 1 - start
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
		currentTime := r.waitSomeEpoch(cliff + 1)
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
		r.waitSomeEpoch(end + 1)
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

func TestBonding(t *testing.T) {
	r := setup(t, nil)
	var scheduleTotalAmount int64 = 1000
	var start int64 = 100 + r.evm.Context.Time.Int64()
	var cliff int64 = 500 + start
	// by making (end - start == scheduleTotalAmount) we have (totalUnlocked = currentTime - start)
	end := scheduleTotalAmount + start
	users, validators, liquidContracts := setupSchedules(r, 2, 2, scheduleTotalAmount, start, cliff, end)

	bondingGas, _, err := r.stakableVesting.RequiredBondingGasCost(nil)
	require.NoError(r.t, err)

	user := users[0]
	scheduleID := common.Big0
	validator := validators[0]
	liquidContract := liquidContracts[0]

	r.run("can bond all funds before cliff", func(r *runner) {
		require.True(r.t, r.evm.Context.Time.Cmp(big.NewInt(start+1)) < 0, "schedule started already")
		bondingAmount := big.NewInt(scheduleTotalAmount / 2)
		_, err := r.stakableVesting.Bond(fromSender(user, bondingGas), scheduleID, validator, bondingAmount)
		require.Equal(r.t, "execution reverted: schedule not started yet", err.Error())
		r.waitSomeBlock(start + 1)
		require.True(r.t, r.evm.Context.Time.Cmp(big.NewInt(cliff+1)) < 0, "schedule cliff finished already")
		bondAndFinalize(r, user, validator, liquidContract, scheduleID, bondingAmount, bondingGas)
	})

	// start schedule for bonding for all the tests remaining
	r.waitSomeBlock(start + 1)

	r.run("cannot bond more than total", func(r *runner) {
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
		r.NoError(
			r.stakableVesting.Bond(fromSender(user, bondingGas), scheduleID, validator, bondingAmount),
		)
		// let bonding apply
		r.waitNextEpoch()
		currentTime := r.waitSomeEpoch(cliff + 1)
		// contract has context of last block
		unlocked := currentTime - 1 - start
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
		r.waitSomeEpoch(end + 1)
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

	r.run("track liquids when bonding from multiple schedules to multiple validators", func(r *runner) {
		// TODO: complete
	})

	r.run("when bonded, release NTN first", func(r *runner) {
		liquidBalance, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, scheduleID, validator)
		require.NoError(r.t, err)
		require.True(r.t, scheduleTotalAmount > 10, "cannot test")
		bondingAmount := big.NewInt(scheduleTotalAmount / 10)
		bondAndFinalize(r, user, validator, liquidContract, scheduleID, bondingAmount, bondingGas)
		remaining := new(big.Int).Sub(big.NewInt(scheduleTotalAmount), bondingAmount)
		require.True(r.t, remaining.Cmp(common.Big0) > 0, "no NTN remains")
		r.waitSomeEpoch(cliff + 1)
		unlocked, _, err := r.stakableVesting.UnlockedFunds(nil, user, scheduleID)
		require.NoError(r.t, err)
		require.True(r.t, unlocked.Cmp(remaining) < 0, "don't want to release all NTN in the test")
		balance, _, err := r.autonity.BalanceOf(nil, user)
		require.NoError(r.t, err)
		r.NoError(
			r.stakableVesting.ReleaseFunds(fromSender(user, nil), scheduleID),
		)
		newLiquidBalance, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, scheduleID, validator)
		require.NoError(r.t, err)
		require.Equal(r.t, new(big.Int).Add(liquidBalance, bondingAmount), newLiquidBalance, "lquid released")
		newBalance, _, err := r.autonity.BalanceOf(nil, user)
		require.NoError(r.t, err)
		require.Equal(r.t, new(big.Int).Add(balance, unlocked), newBalance, "balance not updated")
	})

	r.run("test release when bonding to multiple validator", func(r *runner) {
		// TODO: complete
	})
}

func TestUnbonding(t *testing.T) {
	r := setup(t, nil)
	var scheduleTotalAmount int64 = 1000
	var start int64 = 100 + r.evm.Context.Time.Int64()
	var cliff int64 = 500 + start
	// by making (end - start == scheduleTotalAmount) we have (totalUnlocked = currentTime - start)
	end := scheduleTotalAmount + start
	validatorCount := 2
	scheduleCount := 2
	users, validators, liquidContracts := setupSchedules(r, scheduleCount, validatorCount, scheduleTotalAmount, start, cliff, end)

	bondingGas, _, err := r.stakableVesting.RequiredBondingGasCost(nil)
	require.NoError(r.t, err)
	unbondingGas, _, err := r.stakableVesting.RequiredUnbondingGasCost(nil)
	require.NoError(r.t, err)

	// bond from all schedules to all validators
	r.waitSomeBlock(start + 1)
	bondingAmount := big.NewInt(scheduleTotalAmount / int64(validatorCount))
	require.True(r.t, bondingAmount.Cmp(common.Big0) > 0, "not enough to bond")
	for _, user := range users {
		for i := 0; i < scheduleCount; i++ {
			for _, validator := range validators {
				r.NoError(
					r.stakableVesting.Bond(fromSender(user, bondingGas), big.NewInt(int64(i)), validator, bondingAmount),
				)
			}
		}
	}

	r.waitNextEpoch()
	for _, user := range users {
		for i := 0; i < scheduleCount; i++ {
			totalLiquid := big.NewInt(0)
			for _, validator := range validators {
				liquid, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, big.NewInt(int64(i)), validator)
				require.NoError(r.t, err)
				require.Equal(r.t, bondingAmount, liquid)
				totalLiquid.Add(totalLiquid, liquid)
			}
			require.Equal(r.t, big.NewInt(scheduleTotalAmount), totalLiquid)
		}
	}

	// for testing single unbonding
	user := users[0]
	scheduleID := common.Big0
	validator := validators[0]
	liquidContract := liquidContracts[0]

	r.run("can unbond", func(r *runner) {
		liquid, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, scheduleID, validator)
		require.NoError(r.t, err)
		require.Equal(r.t, bondingAmount, liquid, "liquid not minted properly")
		unbondAndRelease(r, user, validator, liquidContract, scheduleID, liquid, unbondingGas)
	})

	r.run("cannot unbond more than total liquid", func(r *runner) {
		unbondingAmount := new(big.Int).Add(bondingAmount, big.NewInt(10))
		_, err = r.stakableVesting.Unbond(fromSender(user, unbondingGas), scheduleID, validator, unbondingAmount)
		require.Equal(r.t, "execution reverted: not enough unlocked liquid tokens", err.Error())
		unbondingAmount = big.NewInt(10)
		remaining := new(big.Int).Sub(bondingAmount, unbondingAmount)
		require.True(r.t, remaining.Cmp(common.Big0) > 0, "cannot test if no liquid remains")
		r.NoError(
			r.stakableVesting.Unbond(fromSender(user, unbondingGas), scheduleID, validator, unbondingAmount),
		)
		lockedLiquid, _, err := r.stakableVesting.LockedLiquidBalanceOf(nil, user, scheduleID, validator)
		require.NoError(r.t, err)
		require.Equal(r.t, unbondingAmount, lockedLiquid)
		unbondingAmount = new(big.Int).Add(remaining, big.NewInt(10))
		_, err = r.stakableVesting.Unbond(fromSender(user, unbondingGas), scheduleID, validator, unbondingAmount)
		require.Equal(r.t, "execution reverted: not enough unlocked liquid tokens", err.Error())
		r.waitNextEpoch()
		_, err = r.stakableVesting.Unbond(fromSender(user, unbondingGas), scheduleID, validator, unbondingAmount)
		require.Equal(r.t, "execution reverted: not enough unlocked liquid tokens", err.Error())
		r.NoError(
			r.stakableVesting.Unbond(fromSender(user, unbondingGas), scheduleID, validator, remaining),
		)
	})

	r.run("cannot unbond if released", func(r *runner) {
		liquid, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, scheduleID, validator)
		require.NoError(r.t, err)
		totalToRelease := liquid.Int64() + 10
		currentTime := r.waitSomeEpoch(totalToRelease + start + 1)
		totalToRelease = currentTime - 1 - start
		r.NoError(
			r.stakableVesting.ReleaseAllLNTN(fromSender(user, nil), scheduleID),
		)
		_, err = r.stakableVesting.Unbond(fromSender(user, unbondingGas), scheduleID, validator, liquid)
		require.Equal(r.t, "execution reverted: not enough unlocked liquid tokens", err.Error())
		// LNTN will be released from then first validator in the list
		newLiquid, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, scheduleID, validator)
		require.NoError(r.t, err)
		require.True(r.t, newLiquid.Cmp(common.Big0) == 0, "liquid remains after unbonding")
		// if more unlocked funds remain, then LNTN will be released from 2nd validator
		validator1 := validators[1]
		_, err = r.stakableVesting.Unbond(fromSender(user, unbondingGas), scheduleID, validator1, liquid)
		require.Equal(r.t, "execution reverted: not enough unlocked liquid tokens", err.Error())
		releasedFromAnotherValidator1 := totalToRelease - liquid.Int64()
		remainingLiquid := new(big.Int).Sub(liquid, big.NewInt(releasedFromAnotherValidator1))
		liquid, _, err = r.stakableVesting.LiquidBalanceOf(nil, user, scheduleID, validator1)
		require.NoError(r.t, err)
		require.Equal(r.t, remainingLiquid, liquid, "liquid balance mismatch")
		unbondAndRelease(r, user, validator1, liquidContracts[1], scheduleID, remainingLiquid, unbondingGas)
	})

	r.run("track liquid when unbonding from multiple schedules to multiple validators", func(r *runner) {
		// TODO: complete
	})
}

func TestStakingRevert(t *testing.T) {
	r := setup(t, nil)
	// TODO: complete setup

	r.run("fails to notify reward distribution", func(r *runner) {
		// TODO: complete
	})

	r.run("reject bonding request and notify rejection", func(r *runner) {
		// TODO: complete
	})

	r.run("reject bonding request but fails to notify", func(r *runner) {
		// TODO: complete
	})

	r.run("revert applied bonding", func(r *runner) {
		// TODO: complete
	})

	r.run("reject unbonding request and notify rejection", func(r *runner) {
		// TODO: complete
	})

	r.run("reject unbonding request but fails to notify", func(r *runner) {
		// TODO: complete
	})

	r.run("revert applied unbonding", func(r *runner) {
		// TODO: complete
	})

	r.run("revert released unbonding", func(r *runner) {
		// TODO: complete
	})
}

func TestRwardTracking(t *testing.T) {
	r := setup(t, nil)
	gas := big.NewInt(10)
	balance := r.getBalanceOf(user)
	fmt.Printf("user balance %v\n", balance)
	balance = r.getBalanceOf(r.autonity.address)
	fmt.Printf("autonity balance %v\n", balance)
	r.NoError(
		r.autonity.Receive(fromSender(user, gas)),
	)
	balance = r.getBalanceOf(user)
	fmt.Printf("user balance %v\n", balance)
	balance = r.getBalanceOf(r.autonity.address)
	fmt.Printf("autonity balance %v\n", balance)
	// var scheduleTotalAmount int64 = 1000
	// var start int64 = 100 + r.evm.Context.Time.Int64()
	// var cliff int64 = 500 + start
	// // by making (end - start == scheduleTotalAmount) we have (totalUnlocked = currentTime - start)
	// end := scheduleTotalAmount + start
	// validatorCount := 2
	// scheduleCount := 2
	// users, validators, liquidContracts := setupSchedules(r, scheduleCount, validatorCount, scheduleTotalAmount, start, cliff, end)

	// bondingGas, _, err := r.stakableVesting.RequiredBondingGasCost(nil)
	// require.NoError(r.t, err)
	// // unbondingGas, _, err := r.stakableVesting.RequiredUnbondingGasCost(nil)
	// // require.NoError(r.t, err)

	// // for testing single unbonding
	// user := users[0]
	// scheduleID := common.Big0
	// validator := validators[0]
	// liquidContract := liquidContracts[0]

	// start schedule to bond
	// r.waitSomeBlock(start + 1)

	r.run("bond and get reward", func(r *runner) {
		// bondingAmount := big.NewInt(scheduleTotalAmount)
		// r.NoError(
		// 	r.stakableVesting.Bond(fromSender(user, bondingGas), scheduleID, validator, bondingAmount),
		// )
		// r.waitNextEpoch()
		// lastEpochTime := new(big.Int).Sub(r.evm.Context.Time, common.Big1)
		// stakeSupply, _, err := r.autonity.TotalSupply(nil)
		// require.NoError(r.t, err)
		// inflationReserve, _, err := r.autonity.InflationReserve(nil)
		// require.NoError(r.t, err)

		// r.waitNextEpoch()
		// currentEpochTime := new(big.Int).Sub(r.evm.Context.Time, common.Big1)
	})

	r.run("bond in differenet epoch and track reward", func(r *runner) {
		// TODO: complete
	})

	r.run("release liquid and track reward", func(r *runner) {
		// TODO: complete
	})

	r.run("unbond and track reward", func(r *runner) {
		// TODO: complete
	})

	r.run("unbond in different epoch and track reward", func(r *runner) {
		// TODO: complete
	})

	r.run("bond from multiple schedules to multiple validators and track reward", func(r *runner) {
		// TODO: complete
	})
}

func TestScheduleUpdateWhenSlashed(t *testing.T) {
	r := setup(t, nil)
	// TODO: complete setup

	r.run("schedule total value update when bonded validator slashed", func(r *runner) {
		// TODO: complete
	})
}

func TestScheduleCancel(t *testing.T) {
	r := setup(t, nil)
	// TODO: complete setup

	r.run("beneficiary changes when schedule is canceled", func(r *runner) {
		// TODO: complete
	})
}

func TestAccessRestriction(t *testing.T) {
	r := setup(t, nil)
	// TODO: complete setup

	r.run("only operator can create schedule", func(r *runner) {
		// TODO: complete
	})

	r.run("only operator can cancel schedule", func(r *runner) {
		// TODO: complete
	})

	r.run("only autonity can notify staking operations", func(r *runner) {
		// TODO: complete
		/*
			1. rewards distribution
			2. bonding applied
			3. unbonding applied
			4. unbonding released
		*/
	})

	r.run("only operator can set gas cost", func(r *runner) {
		// TODO: complete
	})

	r.run("cannot request bonding or unbonding without enough gas", func(r *runner) {
		// TODO: complete
	})
}

func setupSchedules(
	r *runner, scheduleCount, validatorCount int, scheduleTotalAmount, start, cliff, end int64,
) (users []common.Address, validators []common.Address, liquidContracts []*Liquid) {
	users = make([]common.Address, 2)
	users[0] = user
	users[1] = common.HexToAddress("0x88")
	require.NotEqual(r.t, users[0], users[1], "same user")
	for _, user := range users {
		for i := 0; i < scheduleCount; i++ {
			createSchedule(r, user, scheduleTotalAmount, start, cliff, end)
		}
	}

	// use multiple validators
	committee, _, err := r.autonity.GetCommittee(nil)
	require.NoError(r.t, err)
	validators = make([]common.Address, validatorCount)
	liquidContracts = make([]*Liquid, validatorCount)
	abi, err := LiquidMetaData.GetAbi()
	require.NoError(r.t, err)
	for i := 0; i < validatorCount; i++ {
		validators[i] = committee[i].Addr
		validatorInfo, _, err := r.autonity.GetValidator(nil, validators[i])
		require.NoError(r.t, err)
		liquidContracts[i] = &Liquid{&contract{validatorInfo.LiquidContract, abi, r}}
	}
	return
}

func createSchedule(r *runner, beneficiary common.Address, amount, startTime, cliffTime, endTime int64) {
	amountBigInt := big.NewInt(amount)
	r.NoError(
		r.autonity.Mint(operator, r.stakableVesting.address, amountBigInt),
	)
	r.NoError(
		r.stakableVesting.NewSchedule(
			operator, beneficiary, big.NewInt(amount), big.NewInt(startTime),
			big.NewInt(cliffTime), big.NewInt(endTime),
		),
	)
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
	r.NoError(
		r.stakableVesting.Bond(fromSender(user, bondingGas), scheduleID, validatorAddress, bondingAmount),
	)
	abi, err := LiquidMetaData.GetAbi()
	require.NoError(r.t, err)
	liquidContract := &Liquid{&contract{validator.LiquidContract, abi, r}}
	r.NoError(
		liquidContract.Redistribute(fromSender(r.autonity.address, reward), common.Big0),
	)
	bondedValidators := make([]common.Address, 1)
	bondedValidators[0] = validatorAddress
	gasUsedDistribute := r.NoError(
		r.stakableVesting.RewardsDistributed(fromAutonity, bondedValidators),
	)
	if rejected == false {
		r.NoError(
			liquidContract.Mint(fromAutonity, r.stakableVesting.address, bondingAmount),
		)
		liquid = liquid.Add(liquid, bondingAmount)
	}
	gasUsedBond := r.NoError(
		r.stakableVesting.BondingApplied(
			fromAutonity, big.NewInt(int64(bondingID)), validatorAddress, bondingAmount, false, rejected,
		),
	)
	newLiquid, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, scheduleID, validatorAddress)
	require.NoError(r.t, err)
	require.Equal(r.t, liquid, newLiquid)
	return gasUsedDistribute, gasUsedBond
}

func unbondAndApply(
	r *runner, validatorAddress, user common.Address, unbondingID int, scheduleID, unbondingAmount, unbondingGas *big.Int, rejected bool,
) (uint64, uint64, uint64) {
	liquid, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, scheduleID, validatorAddress)
	require.NoError(r.t, err)
	r.NoError(
		r.stakableVesting.Unbond(fromSender(user, unbondingGas), scheduleID, validatorAddress, unbondingAmount),
	)
	abi, err := LiquidMetaData.GetAbi()
	require.NoError(r.t, err)
	validator, _, err := r.autonity.GetValidator(nil, validatorAddress)
	require.NoError(r.t, err)
	liquidContract := &Liquid{&contract{validator.LiquidContract, abi, r}}
	r.NoError(
		liquidContract.Redistribute(fromSender(r.autonity.address, reward), common.Big0),
	)
	bondedValidators := make([]common.Address, 1)
	bondedValidators[0] = validatorAddress
	gasUsedDistribute := r.NoError(
		r.stakableVesting.RewardsDistributed(fromAutonity, bondedValidators),
	)
	if rejected == false {
		r.NoError(
			liquidContract.Unlock(fromAutonity, r.stakableVesting.address, unbondingAmount),
		)
		r.NoError(
			liquidContract.Burn(fromAutonity, r.stakableVesting.address, unbondingAmount),
		)
		liquid = liquid.Sub(liquid, unbondingAmount)
	}
	gasUsedUnbond := r.NoError(
		r.stakableVesting.UnbondingApplied(fromAutonity, big.NewInt(int64(unbondingID)), validatorAddress, rejected),
	)
	newLiquid, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, scheduleID, validatorAddress)
	require.NoError(r.t, err)
	require.Equal(r.t, liquid, newLiquid)
	gasUsedRelease := r.NoError(
		r.stakableVesting.UnbondingReleased(fromAutonity, big.NewInt(int64(unbondingID)), unbondingAmount, rejected),
	)
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
	r.NoError(
		r.stakableVesting.ReleaseAllNTN(fromSender(user, nil), scheduleID),
	)
	newBalance, _, err := r.autonity.BalanceOf(nil, user)
	require.NoError(r.t, err)
	require.Equal(r.t, new(big.Int).Add(initBalance, totalUnlocked), newBalance, "balance mismatch")
	schedule, _, err = r.stakableVesting.GetSchedule(nil, user, scheduleID)
	require.NoError(r.t, err)
	require.True(r.t, new(big.Int).Sub(scheduleNTN, unlockAmount).Cmp(schedule.CurrentNTNAmount) == 0, "schedule not updated properly")
}

func bondAndFinalize(
	r *runner, user, validator common.Address, liquidContract *Liquid, scheduleID, bondingAmount, bondingGas *big.Int,
) {
	newtonBalance, _, err := r.autonity.BalanceOf(nil, r.stakableVesting.address)
	require.NoError(r.t, err)
	liquidOfVestingContract, _, err := liquidContract.BalanceOf(nil, r.stakableVesting.address)
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
	liquidOfVestingContract = new(big.Int).Add(bondingAmount, liquidOfVestingContract)
	liquidOfUser = new(big.Int).Add(bondingAmount, liquidOfUser)
	liquid, _, err := liquidContract.BalanceOf(nil, r.stakableVesting.address)
	require.NoError(r.t, err)
	require.Equal(r.t, liquidOfVestingContract, liquid, "liquid not minted")
	liquid, _, err = r.stakableVesting.LiquidBalanceOf(nil, user, scheduleID, validator)
	require.NoError(r.t, err)
	require.Equal(r.t, liquidOfUser, liquid, "vesting contract cannot track liquid balance")
	newNewtonBalance := new(big.Int).Sub(newtonBalance, bondingAmount)
	newtonBalance, _, err = r.autonity.BalanceOf(nil, r.stakableVesting.address)
	require.NoError(r.t, err)
	require.True(r.t, newNewtonBalance.Cmp(newtonBalance) == 0, "newton balance not updated")
}

func unbondAndRelease(
	r *runner, user, validator common.Address, liquidContract *Liquid, scheduleID, unbondingAmount, unbondingGas *big.Int,
) {
	newtonBalance, _, err := r.autonity.BalanceOf(nil, r.stakableVesting.address)
	require.NoError(r.t, err)
	liquidOfUser, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, scheduleID, validator)
	require.NoError(r.t, err)
	liquidOfVestingContract, _, err := liquidContract.BalanceOf(nil, r.stakableVesting.address)
	require.NoError(r.t, err)
	schedule, _, err := r.stakableVesting.GetSchedule(nil, user, scheduleID)
	require.NoError(r.t, err)
	unbondingRequestBlock := r.evm.Context.BlockNumber
	r.NoError(
		r.stakableVesting.Unbond(fromSender(user, unbondingGas), scheduleID, validator, unbondingAmount),
	)
	r.waitNextEpoch()
	liquidOfUser = new(big.Int).Sub(liquidOfUser, unbondingAmount)
	liquidOfVestingContract = new(big.Int).Sub(liquidOfVestingContract, unbondingAmount)
	liquid, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, scheduleID, validator)
	require.NoError(r.t, err)
	require.True(r.t, liquid.Cmp(liquidOfUser) == 0, "liquid balance mismatch after unbonding")
	liquid, _, err = liquidContract.BalanceOf(nil, r.stakableVesting.address)
	require.NoError(r.t, err)
	require.True(r.t, liquid.Cmp(liquidOfVestingContract) == 0, "liquid balance mismatch after unbonding")

	// release unbonding
	unbondingPeriod, _, err := r.autonity.GetUnbondingPeriod(nil)
	require.NoError(r.t, err)
	unbondingReleaseBlock := new(big.Int).Add(unbondingRequestBlock, unbondingPeriod)
	for unbondingReleaseBlock.Cmp(r.evm.Context.BlockNumber) >= 0 {
		r.waitNextEpoch()
	}
	newNewtonAmount := new(big.Int).Add(schedule.CurrentNTNAmount, unbondingAmount)
	schedule, _, err = r.stakableVesting.GetSchedule(nil, user, scheduleID)
	require.NoError(r.t, err)
	require.Equal(r.t, newNewtonAmount, schedule.CurrentNTNAmount, "schedule not updated")
	newNewtonBalance := new(big.Int).Add(newtonBalance, unbondingAmount)
	newtonBalance, _, err = r.autonity.BalanceOf(nil, r.stakableVesting.address)
	require.NoError(r.t, err)
	require.Equal(r.t, newNewtonBalance, newtonBalance, "vesting contract balance mismatch")
}
