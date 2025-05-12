package router

import (
	"math"
	"math/big"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/internal/testrand"
)

func TestRouter(t *testing.T) {
	// ToDo: test routing
	t.Run("Test default clustering", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		broadcaster := consensus.NewMockBroadcaster(ctrl)
		key, err := crypto.GenerateKey()
		require.NoError(t, err)

		router := New(broadcaster, key, nil, nil, crypto.PubkeyToAddress(key.PublicKey))
		committee := make([]common.Address, 200)
		for i := 0; i < 200; i++ {
			committee[i] = testrand.Address()
		}
		// Set up the mock to return the committee
		router.epoch = &types.Epoch{
			PreviousEpochBlock: big.NewInt(0),
			NextEpochBlock:     big.NewInt(120),
			Delta:              nil,
		}
	})

	t.Run("Test clustering for even number", func(t *testing.T) {
		n := 100
		committee := make([]common.Address, n)
		latMat := make([][]uint8, n)
		for i := 0; i < n; i++ {
			committee[i] = testrand.Address()
			latMat[i] = make([]uint8, n)
			for j := 0; j < n; j++ {
				latMat[i][j] = uint8(rand.Intn(256))
			}
		}
		// Set up the mock to return the committee
		clusters, err := assignClusters(committee, latMat, int(math.Sqrt(float64(n))))
		require.NoError(t, err)
		require.Equal(t, len(clusters), int(math.Sqrt(float64(n))))
		for _, cluster := range clusters {
			require.Len(t, cluster, n/int(math.Sqrt(float64(n))))
		}
	})

	t.Run("Test clustering for odd number", func(t *testing.T) {
		n := 135
		committee := make([]common.Address, n)
		latMat := make([][]uint8, n)
		for i := 0; i < n; i++ {
			committee[i] = testrand.Address()
			latMat[i] = make([]uint8, n)
			for j := 0; j < n; j++ {
				latMat[i][j] = uint8(rand.Intn(256))
			}
		}
		clusters, err := assignClusters(committee, latMat, int(math.Sqrt(float64(n))))
		require.NoError(t, err)
		require.Equal(t, len(clusters), int(math.Sqrt(float64(n))))
		for _, cluster := range clusters {
			if len(cluster) == n/int(math.Sqrt(float64(n))) {
				continue
			}
			require.True(t, len(clusters) <= n/int(math.Sqrt(float64(n)))+n%int(math.Sqrt(float64(n))))
		}
	})

}
