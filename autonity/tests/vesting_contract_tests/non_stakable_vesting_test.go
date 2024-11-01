package vestingtests

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/autonity/tests"
	"github.com/autonity/autonity/common"
)

func TestReleaseFromNonStakeableContract(t *testing.T) {
	var amount int64 = 100
	start := time.Now().Unix() + 10
	// having (amount = end - start) makes (unlockedFunds = time - start)
	end := amount + start
	cliffDuration := big.NewInt(amount / 2)
	cliff := start + cliffDuration.Int64()
	user := tests.User
	subscribedAmount := big.NewInt(amount)

	setup := func() *tests.Runner {
		r := tests.Setup(t, nil)
		createSchedule(r, amount, start, end)
		subscribeToSchedule(r, user, subscribedAmount, common.Big0, cliffDuration)
		return r
	}

	tests.RunWithSetup("vested and withdrawale vested funds are 0 before start", setup, func(r *tests.Runner) {
		vestedFunds, _, err := r.NonStakeableVesting.VestedFunds(nil, user, common.Big0)
		require.NoError(r.T, err)
		require.True(r.T, vestedFunds.Cmp(common.Big0) == 0)

		withdrawable, _, err := r.NonStakeableVesting.WithdrawableVestedFunds(nil, user, common.Big0)
		require.NoError(r.T, err)
		require.True(r.T, withdrawable.Cmp(common.Big0) == 0)

		releaseAllNTN(r, user, common.Big0)
	})

	tests.RunWithSetup("vested funds increase but withdrawale vested funds are 0 after start and before cliff", setup, func(r *tests.Runner) {
		currentTime := r.WaitForEpochsUntil(start + 2)
		require.True(r.T, currentTime < cliff, "cannot test, cliff reached")
		vestedFunds, _, err := r.NonStakeableVesting.VestedFunds(nil, user, common.Big0)
		require.NoError(r.T, err)
		require.Equal(r.T, big.NewInt(currentTime-start-1), vestedFunds)

		withdrawable, _, err := r.NonStakeableVesting.WithdrawableVestedFunds(nil, user, common.Big0)
		require.NoError(r.T, err)
		require.True(r.T, withdrawable.Cmp(common.Big0) == 0)

		releaseAllNTN(r, user, common.Big0)
	})

	tests.RunWithSetup("vested and withdrawale vested funds are equal after cliff", setup, func(r *tests.Runner) {
		r.WaitForEpochsUntil(cliff + 1)
		vestedFunds, _, err := r.NonStakeableVesting.VestedFunds(nil, user, common.Big0)
		require.NoError(r.T, err)

		withdrawable, _, err := r.NonStakeableVesting.WithdrawableVestedFunds(nil, user, common.Big0)
		require.NoError(r.T, err)
		require.Equal(r.T, vestedFunds, withdrawable)
	})

	tests.RunWithSetup("unlocks linearly (epoch based) (after cliff) between start and end", setup, func(r *tests.Runner) {
		currentTime := r.WaitForEpochsUntil(cliff + 1)
		unlockAmount := big.NewInt(currentTime - start - 1)

		epochID, _, err := r.Autonity.EpochID(nil)
		require.NoError(r.T, err)
		// mine some more blocks, shouldn't matter because unlocking is epoch based
		r.WaitNBlocks(10)
		newEpochID, _, err := r.Autonity.EpochID(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, epochID, newEpochID, "cannot test if epoch progresses")

		unlockedFunds, _, err := r.NonStakeableVesting.VestedFunds(nil, user, common.Big0)
		require.NoError(r.T, err)
		require.Equal(r.T, unlockAmount, unlockedFunds)

		releaseNTN(r, user, new(big.Int).Add(unlockAmount, common.Big1), false, false)
		releaseNTN(r, user, unlockedFunds, true, true)
		releaseAllNTN(r, user, unlockedFunds)

		// unlocked funds should be 0 now
		unlockedFunds, _, err = r.NonStakeableVesting.VestedFunds(nil, user, common.Big0)
		require.NoError(r.T, err)
		require.True(r.T, unlockedFunds.Cmp(common.Big0) == 0)

		releaseNTN(r, user, new(big.Int).Add(unlockAmount, common.Big1), false, false)
		releaseNTN(r, user, unlockedFunds, true, true)
		releaseAllNTN(r, user, unlockedFunds)
	})

	tests.RunWithSetup("can release in chunks", setup, func(r *tests.Runner) {
		currentTime := r.WaitForEpochsUntil(cliff + 1)
		unlockAmount := big.NewInt(currentTime - start - 1)
		require.True(r.T, unlockAmount.Cmp(common.Big2) >= 0, "cannot test")
		unlockedFunds, _, err := r.NonStakeableVesting.VestedFunds(nil, user, common.Big0)
		require.NoError(r.T, err)
		require.Equal(r.T, unlockAmount, unlockedFunds)

		releaseAmount := new(big.Int).Sub(unlockAmount, common.Big1)
		releaseNTN(r, user, releaseAmount, true, false)

		unlockAmount = new(big.Int).Sub(unlockAmount, releaseAmount)
		unlockedFunds, _, err = r.NonStakeableVesting.VestedFunds(nil, user, common.Big0)
		require.NoError(r.T, err)
		require.Equal(r.T, unlockAmount, unlockedFunds)

		releaseNTN(r, user, new(big.Int).Add(unlockAmount, common.Big1), false, false)
		releaseNTN(r, user, unlockAmount, true, true)
		releaseAllNTN(r, user, unlockAmount)
	})

	tests.RunWithSetup("cannot unlock more than total amount", setup, func(r *tests.Runner) {
		r.WaitForEpochsUntil(end + 1)

		unlockAmount := big.NewInt(amount)
		unlockedFunds, _, err := r.NonStakeableVesting.VestedFunds(nil, user, common.Big0)
		require.NoError(r.T, err)
		require.Equal(r.T, unlockAmount, unlockedFunds)

		// wait some more, shouldn't unlock more
		r.WaitNextEpoch()
		unlockedFunds, _, err = r.NonStakeableVesting.VestedFunds(nil, user, common.Big0)
		require.NoError(r.T, err)
		require.Equal(r.T, unlockAmount, unlockedFunds)

		// withdraw and wait some more, shouldn't unlock anymore
		releaseAmount := common.Big1
		releaseNTN(r, user, releaseAmount, true, false)
		unlockAmount = new(big.Int).Sub(unlockAmount, releaseAmount)
		unlockedFunds, _, err = r.NonStakeableVesting.VestedFunds(nil, user, common.Big0)
		require.NoError(r.T, err)
		require.Equal(r.T, unlockAmount, unlockedFunds)

		r.WaitNextEpoch()
		unlockedFunds, _, err = r.NonStakeableVesting.VestedFunds(nil, user, common.Big0)
		require.NoError(r.T, err)
		require.Equal(r.T, unlockAmount, unlockedFunds)
		releaseNTN(r, user, new(big.Int).Add(unlockAmount, common.Big1), false, false)
		releaseNTN(r, user, unlockAmount, true, true)
		releaseAllNTN(r, user, unlockAmount)
	})
}

func TestTreasuryFunds(t *testing.T) {
	var amount int64 = 100
	start := time.Now().Unix() + 10
	// having (amount = end - start) makes (unlockedFunds = time - start)
	end := amount + start
	setup := func() *tests.Runner {
		r := tests.Setup(t, nil)
		createSchedule(r, amount, start, end)
		return r
	}

	tests.RunWithSetup("unsubscribed funds go to treasury", setup, func(r *tests.Runner) {
		r.WaitForEpochsUntil(end + 1)
		releaseTreasuryFunds(r, big.NewInt(amount), true)
	})

	tests.RunWithSetup("all treasury funds cannot be withdrawn before total duration has passed", setup, func(r *tests.Runner) {
		releaseTreasuryFunds(r, big.NewInt(amount), false, "schedule total duration not expired yet")
	})

	tests.RunWithSetup("cannot withdraw more than once from the same schedule", setup, func(r *tests.Runner) {
		r.WaitForEpochsUntil(end + 1)
		releaseTreasuryFunds(r, big.NewInt(amount), true)
		releaseTreasuryFunds(r, common.Big0, true)
	})

	subscribedAmount := amount / 2
	unsubscribedAmount := amount - subscribedAmount
	user := tests.User

	newSetup := func() *tests.Runner {
		r := setup()
		currentTime := r.WaitForEpochsUntil(start + 10)
		expiredFunds := (currentTime - start - 1) * subscribedAmount / amount
		require.True(r.T, expiredFunds > 0, "cannot test")

		subscribeToSchedule(r, user, big.NewInt(subscribedAmount), common.Big0, common.Big0)
		expiredFromContract, _, err := r.NonStakeableVesting.GetExpiredFunds(nil, user, common.Big0)
		require.NoError(r.T, err)
		require.Equal(r.T, big.NewInt(expiredFunds), expiredFromContract)
		return r
	}

	tests.RunWithSetup("expired funds go to treasury", newSetup, func(r *tests.Runner) {
		expiredFromContract, _, err := r.NonStakeableVesting.GetExpiredFunds(nil, user, common.Big0)
		require.NoError(r.T, err)
		r.WaitForEpochsUntil(end + 1)

		userFunds := new(big.Int).Sub(big.NewInt(subscribedAmount), expiredFromContract)
		treasuryFunds := new(big.Int).Add(expiredFromContract, big.NewInt(unsubscribedAmount))
		scheduleTracker, _, err := r.NonStakeableVesting.GetScheduleTracker(nil, common.Big0)
		require.NoError(r.T, err)
		require.Equal(r.T, big.NewInt(unsubscribedAmount), scheduleTracker.UnsubscribedAmount)
		require.Equal(r.T, expiredFromContract, scheduleTracker.ExpiredFromContract)

		releaseNTN(r, user, new(big.Int).Add(userFunds, common.Big1), false, true)
		releaseNTN(r, user, userFunds, true, true)
		releaseTreasuryFunds(r, treasuryFunds, true)
		scheduleTracker, _, err = r.NonStakeableVesting.GetScheduleTracker(nil, common.Big0)
		require.NoError(r.T, err)
		require.True(r.T, new(big.Int).Add(scheduleTracker.UnsubscribedAmount, scheduleTracker.ExpiredFromContract).Cmp(common.Big0) == 0)
		releaseAllNTN(r, user, userFunds)
	})

	tests.RunWithSetup("expired funds can be withdrawn by treasury any time after they are expired", newSetup, func(r *tests.Runner) {
		expiredFromContract, _, err := r.NonStakeableVesting.GetExpiredFunds(nil, user, common.Big0)
		require.NoError(r.T, err)
		require.True(r.T, end+1 > r.Evm.Context.Time.Int64(), "cannot test, contract ended already")

		treasury, _, err := r.Autonity.GetTreasuryAccount(nil)
		require.NoError(r.T, err)
		balance := r.GetNewtonBalanceOf(treasury)
		r.NoError(
			r.NonStakeableVesting.ReleaseExpiredFundsForTreasury(
				tests.FromSender(treasury, nil),
				common.Big0,
			),
		)
		newBalance := r.GetNewtonBalanceOf(treasury)
		require.Equal(r.T, new(big.Int).Add(balance, expiredFromContract), newBalance)
		scheduleTracker, _, err := r.NonStakeableVesting.GetScheduleTracker(nil, common.Big0)
		require.NoError(r.T, err)
		require.True(r.T, scheduleTracker.ExpiredFromContract.Cmp(common.Big0) == 0)

		balance = newBalance
		r.NoError(
			r.NonStakeableVesting.ReleaseExpiredFundsForTreasury(
				tests.FromSender(treasury, nil),
				common.Big0,
			),
		)
		newBalance = r.GetNewtonBalanceOf(treasury)
		require.Equal(r.T, balance, newBalance)
	})
}

func TestNonStakeableAccessRestriction(t *testing.T) {
	user := tests.User

	setup := func() *tests.Runner {
		r := tests.Setup(t, nil)
		r.NoError(
			r.Autonity.CreateSchedule(r.Operator, r.NonStakeableVesting.Address(), common.Big1, r.Evm.Context.Time, common.Big1),
		)
		return r
	}

	tests.RunWithSetup("only operator can create new contract", setup, func(r *tests.Runner) {
		_, err := r.NonStakeableVesting.NewContract(nil, user, common.Big0, common.Big0, common.Big0)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: caller is not the operator", err.Error())
	})

	tests.RunWithSetup("only operator can change contract beneficiary", setup, func(r *tests.Runner) {
		_, err := r.NonStakeableVesting.ChangeContractBeneficiary(nil, user, common.Big0, user)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: caller is not the operator", err.Error())
	})

	tests.RunWithSetup("only treasury account can claim treasury funds", setup, func(r *tests.Runner) {
		_, err := r.NonStakeableVesting.ReleaseAllFundsForTreasury(
			tests.FromSender(tests.User, nil),
			common.Big0,
		)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: caller is not treasury account", err.Error())

		_, err = r.NonStakeableVesting.ReleaseAllFundsForTreasury(nil, common.Big0)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: caller is not treasury account", err.Error())

		_, err = r.NonStakeableVesting.ReleaseExpiredFundsForTreasury(
			tests.FromSender(tests.User, nil),
			common.Big0,
		)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: caller is not treasury account", err.Error())

		_, err = r.NonStakeableVesting.ReleaseExpiredFundsForTreasury(nil, common.Big0)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: caller is not treasury account", err.Error())
	})
}

func TestContractCreation(t *testing.T) {
	user := tests.User

	setup := func() *tests.Runner {
		return tests.Setup(t, nil)
	}

	tests.RunWithSetup("contract needs to subsribe to schedule", setup, func(r *tests.Runner) {
		_, err := r.NonStakeableVesting.NewContract(r.Operator, user, common.Big1, common.Big0, common.Big0)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: schedule does not exist", err.Error())
	})

	var amount int64 = 100
	start := time.Now().Unix() + 10
	cliffDuration := big.NewInt(0)
	end := amount + start

	newSetup := func() *tests.Runner {
		r := setup()
		createSchedule(r, amount, start, end)
		return r
	}

	tests.RunWithSetup("contract nominal amount cannot exceed schedule nominal amount", newSetup, func(r *tests.Runner) {
		_, err := r.NonStakeableVesting.NewContract(r.Operator, user, big.NewInt(amount+1), common.Big0, cliffDuration)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: not enough funds to create a new contract under schedule", err.Error())
		subscribeToSchedule(r, user, common.Big1, common.Big0, cliffDuration)

		_, err = r.NonStakeableVesting.NewContract(r.Operator, user, big.NewInt(amount), common.Big0, cliffDuration)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: not enough funds to create a new contract under schedule", err.Error())

		newUser := common.HexToAddress("0x88")
		require.NotEqual(r.T, newUser, user)
		_, err = r.NonStakeableVesting.NewContract(r.Operator, newUser, big.NewInt(amount), common.Big0, cliffDuration)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: not enough funds to create a new contract under schedule", err.Error())
	})

	tests.RunWithSetup("contract creation after start loses unlocked funds", newSetup, func(r *tests.Runner) {
		currentTime := r.WaitForEpochsUntil(start + 10)
		unlockAmount := currentTime - start - 1
		// progress some more blocks, as unlocking should be epoch based
		// only unlocked funds should be expired
		r.WaitNBlocks(10)
		subscribeAmount := big.NewInt(amount)
		subscribeToSchedule(r, user, subscribeAmount, common.Big0, cliffDuration)

		// unlockAmount of contract should be 0
		// unlockAmount from schedule shows as expired funds in contract
		// total amount of contract is reduced by the amount of expired funds

		unlockedFunds, _, err := r.NonStakeableVesting.VestedFunds(nil, user, common.Big0)
		require.NoError(r.T, err)
		require.True(r.T, unlockedFunds.Cmp(common.Big0) == 0)
		contract, _, err := r.NonStakeableVesting.GetContract(nil, user, common.Big0)
		require.NoError(r.T, err)
		expiredCalculated := new(big.Int).Sub(big.NewInt(amount), contract.CurrentNTNAmount)
		require.True(r.T, expiredCalculated.Cmp(common.Big0) == 1, "nothing expired")
		require.Equal(r.T, big.NewInt(unlockAmount), expiredCalculated)

		expiredFunds, _, err := r.NonStakeableVesting.GetExpiredFunds(nil, user, common.Big0)
		require.NoError(r.T, err)
		require.Equal(r.T, expiredCalculated, expiredFunds)

		scheduleTracker, _, err := r.NonStakeableVesting.GetScheduleTracker(nil, common.Big0)
		require.NoError(r.T, err)
		require.Equal(r.T, expiredCalculated, scheduleTracker.ExpiredFromContract)
	})

	tests.RunWithSetup("contract creation before start has full funds claimable as unlocks", newSetup, func(r *tests.Runner) {
		subscribeAmount := big.NewInt(amount)
		subscribeToSchedule(r, user, subscribeAmount, common.Big0, cliffDuration)
		r.WaitForEpochsUntil(end + 1)
		// all should unlock, user should be able to claim everything
		unlocked, _, err := r.NonStakeableVesting.VestedFunds(nil, user, common.Big0)
		require.NoError(r.T, err)
		require.Equal(r.T, subscribeAmount, unlocked)
		balance, _, err := r.Autonity.BalanceOf(nil, user)
		require.NoError(r.T, err)
		r.NoError(
			r.NonStakeableVesting.ReleaseAllNTN(
				tests.FromSender(user, nil), common.Big0,
			),
		)
		newBalance, _, err := r.Autonity.BalanceOf(nil, user)
		require.NoError(r.T, err)
		require.Equal(r.T, new(big.Int).Add(balance, subscribeAmount), newBalance)
	})
}

func subscribeToSchedule(r *tests.Runner, beneficiary common.Address, amount, scheduleID, cliffDuration *big.Int) {
	r.NoError(
		r.NonStakeableVesting.NewContract(
			r.Operator, beneficiary, amount, scheduleID, cliffDuration,
		),
	)
}

func createSchedule(r *tests.Runner, amount, startTime, endTime int64) {
	startBig := big.NewInt(startTime)
	endBig := big.NewInt(endTime)
	r.NoError(
		r.Autonity.CreateSchedule(
			r.Operator, r.NonStakeableVesting.Address(), big.NewInt(amount),
			big.NewInt(startTime), new(big.Int).Sub(endBig, startBig),
		),
	)
}

// release NTN
func releaseNTN(r *tests.Runner, user common.Address, releaseAmount *big.Int, success, revert bool) {
	balance := r.GetNewtonBalanceOf(user)
	release := func() {
		_, err := r.NonStakeableVesting.ReleaseNTN(
			tests.FromSender(user, nil),
			common.Big0,
			releaseAmount,
		)
		newBalance := r.GetNewtonBalanceOf(user)
		if success {
			require.NoError(r.T, err)
			require.Equal(r.T, new(big.Int).Add(balance, releaseAmount), newBalance)
		} else {
			require.Error(r.T, err)
			require.Equal(r.T, "execution reverted: not enough unlocked funds", err.Error())
			require.True(r.T, balance.Cmp(newBalance) == 0)
		}
	}

	if revert {
		r.RunAndRevert(func(r *tests.Runner) {
			release()
		})
	} else {
		release()
	}
}

// release all NTN, don't revert
func releaseAllNTN(r *tests.Runner, user common.Address, unlocked *big.Int) {
	balance := r.GetNewtonBalanceOf(user)
	r.NoError(
		r.NonStakeableVesting.ReleaseAllNTN(tests.FromSender(user, nil), common.Big0),
	)
	newBalance := r.GetNewtonBalanceOf(user)
	require.True(r.T, new(big.Int).Add(balance, unlocked).Cmp(newBalance) == 0)
}

func releaseTreasuryFunds(r *tests.Runner, funds *big.Int, success bool, revertingMsg ...string) {
	treasury, _, err := r.Autonity.GetTreasuryAccount(nil)
	require.NoError(r.T, err)
	balance := r.GetNewtonBalanceOf(treasury)
	_, err = r.NonStakeableVesting.ReleaseAllFundsForTreasury(
		tests.FromSender(treasury, nil),
		common.Big0,
	)
	newBalance := r.GetNewtonBalanceOf(treasury)

	if success {
		require.NoError(r.T, err)
		require.Equal(r.T, new(big.Int).Add(balance, funds), newBalance)
	} else {
		errMsg := "execution reverted"
		if len(revertingMsg) > 0 {
			errMsg += ": " + revertingMsg[0]
		}
		require.Error(r.T, err)
		require.Equal(r.T, errMsg, err.Error())
		require.True(r.T, balance.Cmp(newBalance) == 0)
	}
}
