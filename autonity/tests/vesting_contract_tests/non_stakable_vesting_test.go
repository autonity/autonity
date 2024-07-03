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
	cliff := start + amount/2
	createSchedule(r, amount, start, cliff, end)
	user := tests.User
	subscribeAmount := big.NewInt(amount)
	subscribeToSchedule(r, user, subscribeAmount, common.Big0)

	r.Run("cannot unlock before cliff", func(r *tests.Runner) {
		_, _, err := r.NonStakableVesting.UnlockedFunds(nil, user, common.Big0)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: cliff period not reached yet", err.Error())

		_, err = r.NonStakableVesting.ReleaseAllFunds(tests.FromSender(user, nil), common.Big0)
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
			r.NonStakableVesting.ReleaseAllFunds(tests.FromSender(user, nil), common.Big0),
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
			r.NonStakableVesting.ReleaseAllFunds(tests.FromSender(user, nil), common.Big0),
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
			r.NonStakableVesting.ReleaseFund(tests.FromSender(user, nil), common.Big0, releaseAmount),
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
			r.NonStakableVesting.ReleaseAllFunds(tests.FromSender(user, nil), common.Big0),
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
			r.NonStakableVesting.ReleaseFund(tests.FromSender(user, nil), common.Big0, releaseAmount),
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
			r.NonStakableVesting.ReleaseAllFunds(tests.FromSender(user, nil), common.Big0),
		)
		newBalance, _, err = r.Autonity.BalanceOf(nil, user)
		require.NoError(r.T, err)
		require.Equal(r.T, new(big.Int).Add(balance, unlockAmount), newBalance)
	})
}

func TestNonStakableAccessRestriction(t *testing.T) {
	r := tests.Setup(t, nil)

	r.Run("only operator can set total nominal value", func(r *tests.Runner) {
		_, err := r.NonStakableVesting.SetTotalNominal(nil, common.Big0)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: caller is not the operator", err.Error())
	})

	r.Run("only operator can set max allowed duration", func(r *tests.Runner) {
		_, err := r.NonStakableVesting.SetMaxAllowedDuration(nil, common.Big0)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: caller is not the operator", err.Error())
	})

	user := tests.User

	r.Run("only operator can create schedule", func(r *tests.Runner) {
		_, err := r.NonStakableVesting.CreateSchedule(nil, common.Big1, common.Big0, common.Big0, common.Big1)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: caller is not the operator", err.Error())

		_, err = r.NonStakableVesting.CreateSchedule(tests.FromSender(user, nil), common.Big1, common.Big0, common.Big0, common.Big1)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: caller is not the operator", err.Error())
	})

	r.NoError(
		r.NonStakableVesting.CreateSchedule(operator, common.Big1, common.Big0, common.Big0, common.Big1),
	)

	r.Run("only operator can create new contract", func(r *tests.Runner) {
		_, err := r.NonStakableVesting.NewContract(nil, user, common.Big0, common.Big0)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: caller is not the operator", err.Error())
	})

	r.Run("only operator can change contract beneficiary", func(r *tests.Runner) {
		_, err := r.NonStakableVesting.ChangeContractBeneficiary(nil, user, common.Big0, user)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: caller is not the operator", err.Error())
	})

	r.Run("only autonity can unlock tokens", func(r *tests.Runner) {
		_, err := r.NonStakableVesting.UnlockTokens(nil)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: function restricted to Autonity contract", err.Error())
	})
}

func TestContractCreation(t *testing.T) {
	r := tests.Setup(t, nil)
	totalNominal, _, err := r.NonStakableVesting.TotalNominal(nil)
	require.NoError(r.T, err)

	r.Run("schedule nominal amount cannot exceed total nominal amount", func(r *tests.Runner) {
		_, err := r.NonStakableVesting.CreateSchedule(
			operator, new(big.Int).Add(totalNominal, common.Big1), common.Big0, common.Big0, common.Big1,
		)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: not enough funds to create a new schedule", err.Error())

		r.NoError(
			r.NonStakableVesting.CreateSchedule(
				operator, totalNominal, common.Big0, common.Big0, common.Big1,
			),
		)
	})

	r.Run("sum of schedule nominal amount cannot exceed total nominal amount", func(r *tests.Runner) {
		schduleCount := 4
		eachScheduleNominal := new(big.Int).Div(totalNominal, big.NewInt(int64(schduleCount)))
		for i := 1; i < schduleCount; i++ {
			r.NoError(
				r.NonStakableVesting.CreateSchedule(
					operator, eachScheduleNominal, common.Big0, common.Big0, common.Big1,
				),
			)
		}

		_, err := r.NonStakableVesting.CreateSchedule(
			operator, new(big.Int).Add(eachScheduleNominal, common.Big1), common.Big0, common.Big0, common.Big1,
		)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: not enough funds to create a new schedule", err.Error())

		r.NoError(
			r.NonStakableVesting.CreateSchedule(
				operator, eachScheduleNominal, common.Big0, common.Big0, common.Big1,
			),
		)

		_, err = r.NonStakableVesting.CreateSchedule(
			operator, common.Big1, common.Big0, common.Big0, common.Big1,
		)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: not enough funds to create a new schedule", err.Error())
	})

	user := tests.User

	r.Run("contract needs to subsribe to schedule", func(r *tests.Runner) {
		_, err := r.NonStakableVesting.NewContract(operator, user, common.Big1, common.Big0)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: invalid schedule ID", err.Error())
	})

	var amount int64 = 100
	start := r.Evm.Context.Time.Int64() + 1
	cliff := start + 50
	end := amount + start
	createSchedule(r, amount, start, cliff, end)

	r.Run("contract nominal amount cannot exceed schedule nominal amount", func(r *tests.Runner) {
		_, err := r.NonStakableVesting.NewContract(operator, user, big.NewInt(amount+1), common.Big0)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: not enough funds to create a new contract under schedule", err.Error())
		subscribeToSchedule(r, user, common.Big1, common.Big0)

		_, err = r.NonStakableVesting.NewContract(operator, user, big.NewInt(amount), common.Big0)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: not enough funds to create a new contract under schedule", err.Error())

		newUser := common.HexToAddress("0x88")
		require.NotEqual(r.T, newUser, user)
		_, err = r.NonStakableVesting.NewContract(operator, newUser, big.NewInt(amount), common.Big0)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: not enough funds to create a new contract under schedule", err.Error())
	})

	r.Run("contract creation after cliff loses claimable funds", func(r *tests.Runner) {
		treasuryAccount, _, err := r.Autonity.GetTreasuryAccount(nil)
		require.NoError(r.T, err)
		treasuryBalance, _, err := r.Autonity.BalanceOf(nil, treasuryAccount)
		require.NoError(r.T, err)
		currentTime := r.WaitSomeEpoch(cliff + 1)
		unlockAmount := currentTime - start - 1
		subscribeAmount := big.NewInt(amount)
		subscribeToSchedule(r, user, subscribeAmount, common.Big0)

		// unlockAmount of contract should be 0
		// unlockAmount from schedule goes to treasury account and shows as withdrawn in contract
		unlockedFunds, _, err := r.NonStakableVesting.UnlockedFunds(nil, user, common.Big0)
		require.NoError(r.T, err)
		require.True(r.T, unlockedFunds.Cmp(common.Big0) == 0)
		contract, _, err := r.NonStakableVesting.GetContract(nil, user, common.Big0)
		require.NoError(r.T, err)
		require.Equal(r.T, big.NewInt(unlockAmount), contract.WithdrawnValue)

		treasuryNewBalance, _, err := r.Autonity.BalanceOf(nil, treasuryAccount)
		require.NoError(r.T, err)
		require.Equal(r.T, new(big.Int).Add(treasuryBalance, big.NewInt(unlockAmount)), treasuryNewBalance)
	})

	r.Run("contract creation before cliff has full funds claimable as unlocks", func(r *tests.Runner) {
		subscribeAmount := big.NewInt(amount)
		subscribeToSchedule(r, user, subscribeAmount, common.Big0)
		r.WaitSomeEpoch(end + 1)
		// all should unlock, user should be able to claim everything
		unlocked, _, err := r.NonStakableVesting.UnlockedFunds(nil, user, common.Big0)
		require.NoError(r.T, err)
		require.Equal(r.T, subscribeAmount, unlocked)
		balance, _, err := r.Autonity.BalanceOf(nil, user)
		require.NoError(r.T, err)
		r.NoError(
			r.NonStakableVesting.ReleaseAllFunds(
				tests.FromSender(user, nil), common.Big0,
			),
		)
		newBalance, _, err := r.Autonity.BalanceOf(nil, user)
		require.NoError(r.T, err)
		require.Equal(r.T, new(big.Int).Add(balance, subscribeAmount), newBalance)
	})
}

func TestUnlockTokens(t *testing.T) {
	r := tests.Setup(t, nil)

	var amount int64 = 100
	start := r.Evm.Context.Time.Int64() + 29
	cliff := start + 50
	end := start + amount
	createSchedule(r, amount, start, cliff, end)
	beneficiary := tests.User

	r.Run("unsubscribed unlocked amount goes to protocol treasury account", func(r *tests.Runner) {
		// initial balance
		balanceVault, _, err := r.Autonity.BalanceOf(nil, r.NonStakableVesting.Address())
		require.NoError(r.T, err)
		treasuryAccount, _, err := r.Autonity.GetTreasuryAccount(nil)
		require.NoError(r.T, err)
		balanceTreasury, _, err := r.Autonity.BalanceOf(nil, treasuryAccount)
		require.NoError(r.T, err)

		currentTime := r.WaitSomeEpoch(cliff + 1)
		unlockAmount := min(amount, currentTime-1-start)

		// schedule unlock amount goes to treasury account
		newBalanceTreasury, _, err := r.Autonity.BalanceOf(nil, treasuryAccount)
		require.NoError(r.T, err)
		require.Equal(r.T, new(big.Int).Add(balanceTreasury, big.NewInt(unlockAmount)), newBalanceTreasury)

		newBalanceVault, _, err := r.Autonity.BalanceOf(nil, r.NonStakableVesting.Address())
		require.NoError(r.T, err)
		require.True(r.T, newBalanceVault.Cmp(balanceVault) == 0)
	})

	r.Run("subscribed unlocked amount goes to protocol non-stakable-vault", func(r *tests.Runner) {
		// initial balance
		balanceVault, _, err := r.Autonity.BalanceOf(nil, r.NonStakableVesting.Address())
		require.NoError(r.T, err)
		treasuryAccount, _, err := r.Autonity.GetTreasuryAccount(nil)
		require.NoError(r.T, err)
		balanceTreasury, _, err := r.Autonity.BalanceOf(nil, treasuryAccount)
		require.NoError(r.T, err)

		// subscribe 80 NTN (4/5 of total) to first schedule
		subscribeAmount := big.NewInt(80)
		subscribeToSchedule(r, beneficiary, subscribeAmount, common.Big0)

		currentTime := r.WaitSomeEpoch(cliff + 1)
		unlockAmount := min(amount, currentTime-1-start)
		require.True(r.T, unlockAmount%5 == 0, "cannot test")

		// 1/5 of schedule unlock amount goes to treasury account, 4/5 to non-stakable-vault
		newBalanceTreasury, _, err := r.Autonity.BalanceOf(nil, treasuryAccount)
		require.NoError(r.T, err)
		require.Equal(r.T, new(big.Int).Add(balanceTreasury, big.NewInt(unlockAmount/5)), newBalanceTreasury)

		newBalanceVault, _, err := r.Autonity.BalanceOf(nil, r.NonStakableVesting.Address())
		require.NoError(r.T, err)
		require.Equal(r.T, new(big.Int).Add(balanceVault, big.NewInt(4*unlockAmount/5)), newBalanceVault)

		// vault balance is equal to unlocked amount
		unlockedFunds, _, err := r.NonStakableVesting.UnlockedFunds(nil, beneficiary, common.Big0)
		require.NoError(r.T, err)
		require.Equal(r.T, newBalanceVault, unlockedFunds)
	})

	r.Run("cannot unlock more than total nominal", func(r *tests.Runner) {
		// initial balance
		balanceVault, _, err := r.Autonity.BalanceOf(nil, r.NonStakableVesting.Address())
		require.NoError(r.T, err)
		treasuryAccount, _, err := r.Autonity.GetTreasuryAccount(nil)
		require.NoError(r.T, err)
		balanceTreasury, _, err := r.Autonity.BalanceOf(nil, treasuryAccount)
		require.NoError(r.T, err)

		// subscribe some random amount
		subscribeAmount := big.NewInt(29)
		subscribeToSchedule(r, beneficiary, subscribeAmount, common.Big0)

		r.WaitSomeEpoch(end + 1)

		// 29 NTN goes to non-stakable-vault, 71 NTN to treasury account
		newBalanceTreasury, _, err := r.Autonity.BalanceOf(nil, treasuryAccount)
		require.NoError(r.T, err)
		require.Equal(r.T, new(big.Int).Add(balanceTreasury, new(big.Int).Sub(big.NewInt(amount), subscribeAmount)), newBalanceTreasury)
		balanceTreasury = newBalanceTreasury

		newBalanceVault, _, err := r.Autonity.BalanceOf(nil, r.NonStakableVesting.Address())
		require.NoError(r.T, err)
		require.Equal(r.T, new(big.Int).Add(balanceVault, subscribeAmount), newBalanceVault)
		balanceVault = newBalanceVault

		// vault balance is equal to unlocked amount
		unlockedFunds, _, err := r.NonStakableVesting.UnlockedFunds(nil, beneficiary, common.Big0)
		require.NoError(r.T, err)
		require.Equal(r.T, newBalanceVault, unlockedFunds)

		// progress more epoch, should not unlock anymore
		r.WaitNextEpoch()
		newBalanceVault, _, err = r.Autonity.BalanceOf(nil, r.NonStakableVesting.Address())
		require.NoError(r.T, err)
		require.Equal(r.T, balanceVault, newBalanceVault)
		newBalanceTreasury, _, err = r.Autonity.BalanceOf(nil, treasuryAccount)
		require.NoError(r.T, err)
		require.Equal(r.T, balanceTreasury, newBalanceTreasury)
	})
}

func subscribeToSchedule(r *tests.Runner, beneficiary common.Address, amount, scheduleID *big.Int) {
	r.NoError(
		r.NonStakableVesting.NewContract(
			operator, beneficiary, amount, scheduleID,
		),
	)
}

func createSchedule(r *tests.Runner, amount, startTime, cliffTime, endTime int64) {
	startBig := big.NewInt(startTime)
	cliffBig := big.NewInt(cliffTime)
	endBig := big.NewInt(endTime)
	r.NoError(
		r.NonStakableVesting.CreateSchedule(
			operator, big.NewInt(amount), big.NewInt(startTime),
			new(big.Int).Sub(cliffBig, startBig), new(big.Int).Sub(endBig, startBig),
		),
	)
}
