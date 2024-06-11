package vestingtests

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/autonity/tests"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/params"
)

var operator = tests.Operator

func TestReleaseFromNonStakableContract(_ *testing.T) {
	// TODO (tariq): complete
}

func TestNonStakableAccessRestriction(_ *testing.T) {
	// TODO (tariq): complete
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

	scheduleNominals := []int64{1, 2, 3, 4}
	start := new(big.Int).Add(r.Evm.Context.Time, common.Big1)
	cliffDuration := big.NewInt(100)
	for _, nominal := range scheduleNominals {
		amount := new(big.Int).Mul(big.NewInt(nominal), big.NewInt(1000_000)) // million
		amount.Mul(amount, params.DecimalFactor)
		totalDuration := new(big.Int).Div(amount, params.DecimalFactor)
		totalDuration.Div(totalDuration, big.NewInt(1000)) // so each second generates 1000 NTN
		r.NoError(
			r.NonStakableVesting.CreateSchedule(
				operator, amount, start, cliffDuration, totalDuration,
			),
		)
	}
	r.Run("contract nominal amount cannot exceed schedule nominal amount", func(r *tests.Runner) {
		// TODO (tariq): complete
	})

	r.Run("contract creation after cliff loses claimable funds", func(r *tests.Runner) {
		// TODO (tariq): complete
	})

	r.Run("contract creation before cliff has full funds claimable as unlocks", func(r *tests.Runner) {
		// TODO (tariq): complete
	})
}

func TestUnlockTokens(t *testing.T) {
	r := tests.Setup(t, nil)

	r.Run("unsubscribed (subscribed) unlocked amount goes to protocol treasury account (non-stakable-vault)", func(r *tests.Runner) {
		// TODO (tariq): complete
	})

	r.Run("cannot unlock more than total nominal", func(r *tests.Runner) {
		// TODO (tariq): complete
	})
}
