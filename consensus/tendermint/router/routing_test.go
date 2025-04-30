package router

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/internal/testrand"
)

func TestRouter(t *testing.T) {
	// ToDo: test routing
	t.Run("Test constructor should set self", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		broadcaster := consensus.NewMockBroadcaster(ctrl)
		key, err := crypto.GenerateKey()
		require.NoError(t, err)

		router := NewRouter(broadcaster, key, nil, nil)
		require.NotNil(t, router.self)
	})

	t.Run("Test cluster seed does not panic", func(t *testing.T) {
		msg := message.NewFakePropose(message.Fake{
			FakeRound:  0,
			FakeHeight: 0,
			FakeHash:   testrand.Hash(),
		})

		require.NotPanics(t, func() {
			s := seed(msg)
			require.NotZero(t, s)
		})
	})

	t.Run("Test default clustering", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		broadcaster := consensus.NewMockBroadcaster(ctrl)
		key, err := crypto.GenerateKey()
		require.NoError(t, err)

		router := NewRouter(broadcaster, key, nil, nil)
		committee := make([]common.Address, 200)
		for i := 0; i < 200; i++ {
			committee[i] = testrand.Address()
		}
		router.curEpochInfo = &types.EpochInfo{
			Epoch: types.Epoch{
				PreviousEpochBlock: big.NewInt(0),
				NextEpochBlock:     big.NewInt(120),
				Delta:              nil,
			},
			EpochBlock: big.NewInt(0),
		}
		router.buildClusters(committee)
	})
}
