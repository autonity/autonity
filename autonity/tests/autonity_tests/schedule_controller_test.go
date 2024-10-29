package autonitytests

import (
	"math/big"
	"testing"
	"time"

	"github.com/autonity/autonity/autonity/tests"
	"github.com/autonity/autonity/common"
	"github.com/stretchr/testify/require"
)

var operator = tests.Operator

func TestScheduleAccessRestriction(t *testing.T) {
	user := tests.User
	setup := func() *tests.Runner {
		return tests.Setup(t, nil)
	}

	tests.RunWithSetup("only operator can set max allowed duration", setup, func(r *tests.Runner) {
		_, err := r.Autonity.SetMaxScheduleDuration(nil, common.Big0)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: caller is not the operator", err.Error())

		_, err = r.Autonity.SetMaxScheduleDuration(tests.FromSender(user, nil), common.Big0)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: caller is not the operator", err.Error())
	})

	tests.RunWithSetup("only operator can create schedule", setup, func(r *tests.Runner) {
		_, err := r.Autonity.CreateSchedule(nil, user, common.Big1, common.Big1, common.Big1)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: caller is not the operator", err.Error())

		_, err = r.Autonity.CreateSchedule(tests.FromSender(user, nil), user, common.Big1, common.Big1, common.Big1)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: caller is not the operator", err.Error())
	})
}

func TestScheduleOperation(t *testing.T) {
	vaultAddress := common.HexToAddress("0x99")
	var amount int64 = 100

	setup := func() *tests.Runner {
		return tests.Setup(t, nil)
	}

	tests.RunWithSetup("schedule creation mints to vault address", setup, func(r *tests.Runner) {
		balance, _, err := r.Autonity.BalanceOf(nil, vaultAddress)
		require.NoError(r.T, err)
		createSchedule(r, vaultAddress, amount, 0, 0)
		newBalance, _, err := r.Autonity.BalanceOf(nil, vaultAddress)
		require.NoError(r.T, err)
		require.Equal(r.T, new(big.Int).Add(balance, big.NewInt(amount)), newBalance)
	})

	tests.RunWithSetup("schedule creation does not change circulating supply but changes total supply", setup, func(r *tests.Runner) {
		circulatingSupply, _, err := r.Autonity.CirculatingSupply(nil)
		require.NoError(r.T, err)
		totalSupply, _, err := r.Autonity.TotalSupply(nil)
		require.NoError(r.T, err)
		createSchedule(r, vaultAddress, amount, 0, 0)
		newCirculatingSupply, _, err := r.Autonity.CirculatingSupply(nil)
		require.NoError(r.T, err)
		require.True(r.T, circulatingSupply.Cmp(newCirculatingSupply) == 0, "circulating supply changed")
		newTotalSupply, _, err := r.Autonity.TotalSupply(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, new(big.Int).Add(totalSupply, big.NewInt(amount)), newTotalSupply)
	})

	tests.RunWithSetup("schedule total duration cannot exceed max allowed duration", setup, func(r *tests.Runner) {
		maxAllowedDuration, _, err := r.Autonity.GetMaxScheduleDuration(nil)
		require.NoError(r.T, err)
		totalDuration := new(big.Int).Add(maxAllowedDuration, common.Big1)
		_, err = r.Autonity.CreateSchedule(operator, vaultAddress, big.NewInt(amount), common.Big0, totalDuration)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: schedule total duration exceeds max allowed duration", err.Error())
	})

	start := time.Now().Unix() + 10
	// having (amount = totalDuration) makes (unlockedFunds = time - start)
	totalDuration := amount
	newSetup := func() *tests.Runner {
		r := setup()
		createSchedule(r, vaultAddress, amount, start, totalDuration)
		return r
	}

	tests.RunWithSetup("unlocking schedules change circulating supply but not total supply", newSetup, func(r *tests.Runner) {
		circulatingSupply, _, err := r.Autonity.CirculatingSupply(nil)
		require.NoError(r.T, err)
		totalSupply, _, err := r.Autonity.TotalSupply(nil)
		require.NoError(r.T, err)

		rewards := r.RewardsAfterOneEpoch()
		totalSupply.Add(totalSupply, rewards.RewardNTN)
		circulatingSupply.Add(circulatingSupply, rewards.RewardNTN)
		r.WaitNextEpoch()

		schedule, _, err := r.Autonity.GetSchedule(nil, vaultAddress, common.Big0)
		require.NoError(r.T, err)
		require.True(r.T, schedule.UnlockedAmount.Cmp(common.Big0) > 0, "schedule not unlocked")

		newCirculatingSupply, _, err := r.Autonity.CirculatingSupply(nil)
		require.NoError(r.T, err)
		newTotalSupply, _, err := r.Autonity.TotalSupply(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, totalSupply, newTotalSupply)
		require.Equal(r.T, new(big.Int).Add(circulatingSupply, schedule.UnlockedAmount), newCirculatingSupply)
	})

	tests.RunWithSetup("schedule unlocking follows epoch based linear function", newSetup, func(r *tests.Runner) {
		schedule, _, err := r.Autonity.GetSchedule(nil, vaultAddress, common.Big0)
		require.NoError(r.T, err)
		require.True(r.T, schedule.UnlockedAmount.Cmp(common.Big0) == 0)

		r.WaitNextEpoch()
		unlocked := r.Evm.Context.Time.Int64() - 1 - start
		require.True(r.T, unlocked > 0, "cannot test")
		schedule, _, err = r.Autonity.GetSchedule(nil, vaultAddress, common.Big0)
		require.NoError(r.T, err)
		require.Equal(r.T, big.NewInt(unlocked), schedule.UnlockedAmount, "unlocking mechanism not linear")

		epochID, _, err := r.Autonity.EpochID(nil)
		require.NoError(r.T, err)
		info, _, err := r.Autonity.GetEpochInfo(nil)
		require.NoError(r.T, err)
		// don't go to next epoch, but produce enough blocks that will unlock new tokens if unlocking is not epoch based
		produceBlocks := new(big.Int).Sub(info.NextEpochBlock, info.EpochBlock).Int64() - 1
		r.WaitNBlocks(int(produceBlocks))
		newEpochID, _, err := r.Autonity.EpochID(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, epochID, newEpochID, "epoch progressed, cannot test")

		schedule, _, err = r.Autonity.GetSchedule(nil, vaultAddress, common.Big0)
		require.NoError(r.T, err)
		require.Equal(r.T, big.NewInt(unlocked), schedule.UnlockedAmount, "unlocking mechanism not epoch based")
	})
}

func createSchedule(r *tests.Runner, vaultAddress common.Address, amount, startTime, totalDuration int64) {
	if startTime == 0 {
		startTime = r.Evm.Context.Time.Int64()
	}
	if totalDuration == 0 {
		totalDuration = 1
	}
	r.NoError(
		r.Autonity.CreateSchedule(
			operator, vaultAddress, big.NewInt(amount),
			big.NewInt(startTime), big.NewInt(totalDuration),
		),
	)
}
