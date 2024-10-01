package vestingtests

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/autonity/tests"
	"github.com/autonity/autonity/common"
)

var operator = tests.Operator

func TestReleaseFromNonStakableContract(t *testing.T) {
	r := tests.Setup(t, nil)
	var amount int64 = 100
	start := r.Evm.Context.Time.Int64() + 1
	// having (amount = end - start) makes (unlockedFunds = time - start)
	end := amount + start
	cliffDuration := big.NewInt(amount / 2)
	cliff := start + cliffDuration.Int64()
	createSchedule(r, amount, start, end)
	user := tests.User
	subscribeAmount := big.NewInt(amount)
	subscribeToSchedule(r, user, subscribeAmount, common.Big0, cliffDuration)

	r.Run("cannot unlock before cliff", func(r *tests.Runner) {
		_, _, err := r.NonStakableVesting.UnlockedFunds(nil, user, common.Big0)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: cliff period not reached yet", err.Error())

		_, err = r.NonStakableVesting.ReleaseAllNTN(tests.FromSender(user, nil), common.Big0)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: cliff period not reached yet", err.Error())
	})

	r.Run("unlocks linearly (after cliff) between start and end", func(r *tests.Runner) {
		currentTime := r.WaitSomeEpoch(cliff + 1)
		unlockAmount := big.NewInt(currentTime - start - 1)
		unlockedFunds, _, err := r.NonStakableVesting.UnlockedFunds(nil, user, common.Big0)
		require.NoError(r.T, err)
		require.Equal(r.T, unlockAmount, unlockedFunds)

		balance, _, err := r.Autonity.BalanceOf(nil, user)
		require.NoError(r.T, err)
		r.NoError(
			r.NonStakableVesting.ReleaseAllNTN(tests.FromSender(user, nil), common.Big0),
		)
		newBalance, _, err := r.Autonity.BalanceOf(nil, user)
		require.NoError(r.T, err)
		require.Equal(r.T, new(big.Int).Add(balance, unlockedFunds), newBalance)
		balance = newBalance

		// unlocked funds should be 0 now
		unlockedFunds, _, err = r.NonStakableVesting.UnlockedFunds(nil, user, common.Big0)
		require.NoError(r.T, err)
		require.True(r.T, unlockedFunds.Cmp(common.Big0) == 0)

		r.NoError(
			r.NonStakableVesting.ReleaseAllNTN(tests.FromSender(user, nil), common.Big0),
		)
		newBalance, _, err = r.Autonity.BalanceOf(nil, user)
		require.NoError(r.T, err)
		require.Equal(r.T, balance, newBalance)

		// unlock some more
		r.WaitNextEpoch()
		unlockAmount = big.NewInt(r.Evm.Context.Time.Int64() - currentTime)
		require.True(r.T, unlockAmount.Cmp(common.Big2) >= 0, "cannot test")
		unlockedFunds, _, err = r.NonStakableVesting.UnlockedFunds(nil, user, common.Big0)
		require.NoError(r.T, err)
		require.Equal(r.T, unlockAmount, unlockedFunds)

		releaseAmount := new(big.Int).Sub(unlockAmount, common.Big1)
		r.NoError(
			r.NonStakableVesting.ReleaseNTN(tests.FromSender(user, nil), common.Big0, releaseAmount),
		)
		newBalance, _, err = r.Autonity.BalanceOf(nil, user)
		require.NoError(r.T, err)
		require.Equal(r.T, new(big.Int).Add(balance, releaseAmount), newBalance)
		balance = newBalance

		unlockAmount = new(big.Int).Sub(unlockAmount, releaseAmount)
		unlockedFunds, _, err = r.NonStakableVesting.UnlockedFunds(nil, user, common.Big0)
		require.NoError(r.T, err)
		require.Equal(r.T, unlockAmount, unlockedFunds)

		r.NoError(
			r.NonStakableVesting.ReleaseAllNTN(tests.FromSender(user, nil), common.Big0),
		)
		newBalance, _, err = r.Autonity.BalanceOf(nil, user)
		require.NoError(r.T, err)
		require.Equal(r.T, new(big.Int).Add(balance, unlockAmount), newBalance)
	})

	r.Run("cannot unlock more than total amount", func(r *tests.Runner) {
		r.WaitSomeEpoch(end + 1)

		unlockAmount := big.NewInt(amount)
		unlockedFunds, _, err := r.NonStakableVesting.UnlockedFunds(nil, user, common.Big0)
		require.NoError(r.T, err)
		require.Equal(r.T, unlockAmount, unlockedFunds)

		// wait some more, shouldn't unlock more
		r.WaitNextEpoch()
		unlockedFunds, _, err = r.NonStakableVesting.UnlockedFunds(nil, user, common.Big0)
		require.NoError(r.T, err)
		require.Equal(r.T, unlockAmount, unlockedFunds)

		// withdraw and wait some more, shouldn't unlock anymore
		balance, _, err := r.Autonity.BalanceOf(nil, user)
		require.NoError(r.T, err)
		releaseAmount := common.Big1
		r.NoError(
			r.NonStakableVesting.ReleaseNTN(tests.FromSender(user, nil), common.Big0, releaseAmount),
		)
		newBalance, _, err := r.Autonity.BalanceOf(nil, user)
		require.NoError(r.T, err)
		require.Equal(r.T, new(big.Int).Add(balance, releaseAmount), newBalance)
		balance = newBalance
		unlockAmount = new(big.Int).Sub(unlockAmount, releaseAmount)

		r.WaitNextEpoch()
		unlockedFunds, _, err = r.NonStakableVesting.UnlockedFunds(nil, user, common.Big0)
		require.NoError(r.T, err)
		require.Equal(r.T, unlockAmount, unlockedFunds)
		r.NoError(
			r.NonStakableVesting.ReleaseAllNTN(tests.FromSender(user, nil), common.Big0),
		)
		newBalance, _, err = r.Autonity.BalanceOf(nil, user)
		require.NoError(r.T, err)
		require.Equal(r.T, new(big.Int).Add(balance, unlockAmount), newBalance)
	})
}

func TestNonStakableAccessRestriction(t *testing.T) {
	r := tests.Setup(t, nil)
	user := tests.User

	r.NoError(
		r.Autonity.CreateSchedule(operator, r.NonStakableVesting.Address(), common.Big1, common.Big0, common.Big0),
	)

	r.Run("only operator can create new contract", func(r *tests.Runner) {
		_, err := r.NonStakableVesting.NewContract(nil, user, common.Big0, common.Big0, common.Big0)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: caller is not the operator", err.Error())
	})

	r.Run("only operator can change contract beneficiary", func(r *tests.Runner) {
		_, err := r.NonStakableVesting.ChangeContractBeneficiary(nil, user, common.Big0, user)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: caller is not the operator", err.Error())
	})
}

func TestContractCreation(t *testing.T) {
	r := tests.Setup(t, nil)
	user := tests.User

	r.Run("contract needs to subsribe to schedule", func(r *tests.Runner) {
		_, err := r.NonStakableVesting.NewContract(operator, user, common.Big1, common.Big0, common.Big0)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: schedule does not exist", err.Error())
	})

	var amount int64 = 100
	start := r.Evm.Context.Time.Int64() + 1
	cliffDuration := big.NewInt(50)
	end := amount + start
	createSchedule(r, amount, start, end)

	r.Run("contract nominal amount cannot exceed schedule nominal amount", func(r *tests.Runner) {
		_, err := r.NonStakableVesting.NewContract(operator, user, big.NewInt(amount+1), common.Big0, cliffDuration)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: not enough funds to create a new contract under schedule", err.Error())
		subscribeToSchedule(r, user, common.Big1, common.Big0, cliffDuration)

		_, err = r.NonStakableVesting.NewContract(operator, user, big.NewInt(amount), common.Big0, cliffDuration)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: not enough funds to create a new contract under schedule", err.Error())

		newUser := common.HexToAddress("0x88")
		require.NotEqual(r.T, newUser, user)
		_, err = r.NonStakableVesting.NewContract(operator, newUser, big.NewInt(amount), common.Big0, cliffDuration)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: not enough funds to create a new contract under schedule", err.Error())
	})

	r.Run("contract creation after start loses unlocked funds", func(r *tests.Runner) {
		currentTime := r.WaitSomeEpoch(start + 10)
		unlockAmount := currentTime - start - 1
		// progress some more blocks, as unlocking should be epoch based
		// only unlocked funds should be expired
		r.WaitNBlocks(10)
		subscribeAmount := big.NewInt(amount)
		subscribeToSchedule(r, user, subscribeAmount, common.Big0, cliffDuration)

		// unlockAmount of contract should be 0
		// unlockAmount from schedule shows as expired funds in contract
		// total amount of contract is reduced by the amount of expired funds

		unlockedFunds, _, err := r.NonStakableVesting.UnlockedFunds(nil, user, common.Big0)
		require.NoError(r.T, err)
		require.True(r.T, unlockedFunds.Cmp(common.Big0) == 0)
		contract, _, err := r.NonStakableVesting.GetContract(nil, user, common.Big0)
		require.NoError(r.T, err)
		expiredCalculated := new(big.Int).Sub(big.NewInt(amount), contract.CurrentNTNAmount)
		require.True(r.T, expiredCalculated.Cmp(common.Big0) == 1, "nothing expired")
		require.Equal(r.T, big.NewInt(unlockAmount), expiredCalculated)

		expiredFunds, _, err := r.NonStakableVesting.GetExpiredFunds(nil, user, common.Big0)
		require.NoError(r.T, err)
		require.Equal(r.T, expiredCalculated, expiredFunds)
	})

	r.Run("contract creation before start has full funds claimable as unlocks", func(r *tests.Runner) {
		subscribeAmount := big.NewInt(amount)
		subscribeToSchedule(r, user, subscribeAmount, common.Big0, cliffDuration)
		r.WaitSomeEpoch(end + 1)
		// all should unlock, user should be able to claim everything
		unlocked, _, err := r.NonStakableVesting.UnlockedFunds(nil, user, common.Big0)
		require.NoError(r.T, err)
		require.Equal(r.T, subscribeAmount, unlocked)
		balance, _, err := r.Autonity.BalanceOf(nil, user)
		require.NoError(r.T, err)
		r.NoError(
			r.NonStakableVesting.ReleaseAllNTN(
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
		r.NonStakableVesting.NewContract(
			operator, beneficiary, amount, scheduleID, cliffDuration,
		),
	)
}

func createSchedule(r *tests.Runner, amount, startTime, endTime int64) {
	startBig := big.NewInt(startTime)
	endBig := big.NewInt(endTime)
	r.NoError(
		r.Autonity.CreateSchedule(
			operator, r.NonStakableVesting.Address(), big.NewInt(amount),
			big.NewInt(startTime), new(big.Int).Sub(endBig, startBig),
		),
	)
}
