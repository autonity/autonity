package autonitytests

import (
	"testing"

	"github.com/autonity/autonity/autonity/tests"
	"github.com/autonity/autonity/common"
	"github.com/stretchr/testify/require"
)

func TestNonStakableAccessRestriction(t *testing.T) {
	r := tests.Setup(t, nil)

	r.Run("only operator can set max allowed duration", func(r *tests.Runner) {
		_, err := r.Autonity.SetMaxAllowedDuration(nil, common.Big0)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: caller is not the operator", err.Error())
	})

	user := tests.User

	r.Run("only operator can create schedule", func(r *tests.Runner) {
		_, err := r.Autonity.CreateSchedule(nil, common.Address{}, common.Big1, common.Big0, common.Big0)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: caller is not the operator", err.Error())

		_, err = r.Autonity.CreateSchedule(tests.FromSender(user, nil), common.Address{}, common.Big1, common.Big0, common.Big0)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: caller is not the operator", err.Error())
	})
}
