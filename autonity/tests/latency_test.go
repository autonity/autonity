package tests

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/internal/testrand"
	"github.com/autonity/autonity/params"
)

func TestLatency(t *testing.T) {
	setup := func() *Runner {
		// setup code
		return Setup(t, nil)
	}

	RunWithSetup("Test genesis sets the current committee on deployment", setup, func(r *Runner) {
		expectedCommittee := make([]common.Address, len(r.Committee.Validators))
		for i, v := range r.Committee.Validators {
			expectedCommittee[i] = v.NodeAddress
		}
		committee, _, err := r.Latency.GetCommittee(nil)
		require.NoError(t, err)
		require.Equal(t, expectedCommittee, committee)
	})

	RunWithSetup("Test only autonity can call setCommittee", setup, func(r *Runner) {
		users := []common.Address{r.Operator.origin, r.Committee.Validators[0].NodeAddress, testrand.Address()}
		expectedCommittee := make([]common.Address, len(r.Committee.Validators))
		for i, _ := range r.Committee.Validators {
			expectedCommittee[i] = testrand.Address()
		}
		for _, user := range users {
			_, err := r.Latency.SetCommittee(FromSender(user, common.Big0), expectedCommittee)
			require.Error(t, err)
			require.Contains(t, err.Error(), "function restricted to Autonity contract")
		}

		_, err := r.Latency.SetCommittee(FromSender(r.Autonity.address, common.Big0), expectedCommittee)
		require.NoError(t, err)

		committee, _, err := r.Latency.GetCommittee(nil)
		require.NoError(t, err)
		require.Equal(t, expectedCommittee, committee)
	})

	RunWithSetup("Test only committee member can call report latency", setup, func(r *Runner) {
		committee := genCommittee(50)
		_, err := r.Latency.SetCommittee(FromAutonity, committee)
		require.NoError(t, err)

		users := []common.Address{r.Operator.origin, params.AutonityContractAddress, testrand.Address()}
		for _, user := range users {
			_, err := r.Latency.Report(FromSender(user, common.Big0), nil)
			require.Error(t, err)
			require.Contains(t, err.Error(), "not committee")
		}

		latency := generateLatency(len(committee))
		_, err = r.Latency.Report(
			FromSender(committee[0], common.Big0),
			latency,
		)
		require.NoError(t, err)

		// get the reported latency
		reportedLatency, _, err := r.Latency.ReadReport(nil, committee[0])
		require.NoError(t, err)
		require.Equal(t, latency, reportedLatency)
	})

	RunWithSetup("Test report latency fails with invalid committee size", setup, func(r *Runner) {
		committee := genCommittee(50)
		_, err := r.Latency.SetCommittee(FromAutonity, committee)
		require.NoError(t, err)

		_, err = r.Latency.Report(
			FromSender(committee[0], common.Big0),
			generateLatency(len(committee)+1),
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid length")
	})

	RunWithSetup("Test all committee can report and read", setup, func(r *Runner) {
		committee := genCommittee(50)
		_, err := r.Latency.SetCommittee(FromAutonity, committee)
		require.NoError(t, err)

		latencyMat := make([][]uint8, len(committee))
		for i := 0; i < len(committee); i++ {
			latencyMat[i] = generateLatency(len(committee))
		}
		for i, member := range committee {
			_, err = r.Latency.Report(
				FromSender(member, common.Big0),
				latencyMat[i],
			)
			require.NoError(t, err)
		}

		readMat, _, err := r.Latency.Read(nil)
		require.NoError(t, err)
		require.Equal(t, latencyMat, readMat)
	})
}

func genCommittee(size int) []common.Address {
	committee := make([]common.Address, size)
	for i := 0; i < size; i++ {
		committee[i] = testrand.Address()
	}
	return committee
}

func generateLatency(committeeSize int) []uint8 {
	latency := make([]uint8, committeeSize)
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < committeeSize; i++ {
		latency[i] = uint8(r.Intn(255))
	}
	return latency
}
