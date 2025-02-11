package latency

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/crypto"
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
}
