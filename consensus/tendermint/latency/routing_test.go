package latency

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
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

		router := NewRouter(broadcaster, key)
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
}
