package latency

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/internal/testrand"
)

func TestLatencyCache(t *testing.T) {
	t.Run("Test cache should insert matrix line", func(t *testing.T) {
		cache := newLatencyCache()
		validators := []common.Address{
			testrand.Address(),
			testrand.Address(),
			testrand.Address(),
		}

		validator1LatencyVec := []uint8{0, 2, 3}
		cache.insertMatrixLine(validators[0], validators, validator1LatencyVec)

		mat := cache.readMatrix(validators)

		// own latency should be set
		require.Equal(t, uint8(0), mat[validators[0]][0])
		require.Equal(t, uint8(2), mat[validators[0]][1])
		require.Equal(t, uint8(3), mat[validators[0]][2])
	})

	t.Run("Test should return missing measurements", func(t *testing.T) {
		cache := newLatencyCache()
		oldCommittee := []common.Address{
			testrand.Address(),
			testrand.Address(),
			testrand.Address(),
		}

		measurements := map[common.Address]uint8{
			oldCommittee[0]: 0,
			oldCommittee[1]: 1,
			oldCommittee[2]: 2,
		}
		cache.insertMeasurements(oldCommittee, measurements)

		newCommittee := append(oldCommittee, testrand.Address())

		missing := cache.missingMeasurements(newCommittee)
		require.Len(t, missing, 1)
		require.Equal(t, newCommittee[3], missing[0])
	})
}
