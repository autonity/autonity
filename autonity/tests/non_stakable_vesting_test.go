package tests

import (
	"math/big"
	"testing"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/params"
	"github.com/stretchr/testify/require"
)

func TestReleaseFromNonStakableContract(t *testing.T) {
	// TODO (tariq): complete
}

func TestNonStakableAccessRestriction(t *testing.T) {
	// TODO (tariq): complete
}

func TestContractCreation(t *testing.T) {
	r := setup(t, nil)
	totalNominal, _, err := r.nonStakableVesting.TotalNominal(nil)
	require.NoError(r.t, err)

	r.run("schedule nominal amount cannot exceed total nominal amount", func(r *runner) {
		_, err := r.nonStakableVesting.CreateSchedule(
			operator, new(big.Int).Add(totalNominal, common.Big1), common.Big0, common.Big0, common.Big1,
		)
		require.Equal(r.t, "execution reverted: not enough funds to create a new schedule", err.Error())

		r.NoError(
			r.nonStakableVesting.CreateSchedule(
				operator, totalNominal, common.Big0, common.Big0, common.Big1,
			),
		)
	})

	r.run("sum of schedule nominal amount cannot exceed total nominal amount", func(r *runner) {
		schduleCount := 4
		eachScheduleNominal := new(big.Int).Div(totalNominal, big.NewInt(int64(schduleCount)))
		for i := 1; i < schduleCount; i++ {
			r.NoError(
				r.nonStakableVesting.CreateSchedule(
					operator, eachScheduleNominal, common.Big0, common.Big0, common.Big1,
				),
			)
		}

		_, err := r.nonStakableVesting.CreateSchedule(
			operator, new(big.Int).Add(eachScheduleNominal, common.Big1), common.Big0, common.Big0, common.Big1,
		)
		require.Equal(r.t, "execution reverted: not enough funds to create a new schedule", err.Error())

		r.NoError(
			r.nonStakableVesting.CreateSchedule(
				operator, eachScheduleNominal, common.Big0, common.Big0, common.Big1,
			),
		)

		_, err = r.nonStakableVesting.CreateSchedule(
			operator, common.Big1, common.Big0, common.Big0, common.Big1,
		)
		require.Equal(r.t, "execution reverted: not enough funds to create a new schedule", err.Error())
	})

	r.run("contract needs to subsribe to schedule", func(r *runner) {
		_, err := r.nonStakableVesting.NewContract(operator, user, common.Big1, common.Big0)
		require.Equal(r.t, "execution reverted: invalid schedule ID", err.Error())
	})

	scheduleNominals := []int64{1, 2, 3, 4}
	start := new(big.Int).Add(r.evm.Context.Time, common.Big1)
	cliffDuration := big.NewInt(100)
	for _, nominal := range scheduleNominals {
		amount := new(big.Int).Mul(big.NewInt(nominal), big.NewInt(1000_000)) // million
		amount.Mul(amount, params.DecimalFactor)
		totalDuration := new(big.Int).Div(amount, params.DecimalFactor)
		totalDuration.Div(totalDuration, big.NewInt(1000)) // so each second generates 1000 NTN
		r.NoError(
			r.nonStakableVesting.CreateSchedule(
				operator, amount, start, cliffDuration, totalDuration,
			),
		)
	}
	r.run("contract nominal amount cannot exceed schedule nominal amount", func(r *runner) {
		// TODO (tariq): complete
	})

	r.run("contract creation after cliff loses claimable funds", func(r *runner) {
		// TODO (tariq): complete
	})

	r.run("contract creation before cliff has full funds claimable as unlocks", func(r *runner) {
		// TODO (tariq): complete
	})
}

func TestUnlockTokens(t *testing.T) {
	r := setup(t, nil)

	r.run("unsubscribed (subscribed) unlocked amount goes to protocol treasury account (non-stakable-vault)", func(r *runner) {
		// TODO (tariq): complete
	})

	r.run("cannot unlock more than total nominal", func(r *runner) {
		// TODO (tariq): complete
	})
}
