package e2e

import (
	"testing"

	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/node"
	"github.com/stretchr/testify/require"
)

func newCustomGossiper(b interfaces.Backend) interfaces.Gossiper {
	return &customGossiper{b.Gossiper()}
}

type customGossiper struct {
	interfaces.Gossiper
}

func (cg *customGossiper) Gossip(_ types.Committee, _ []byte) {
	log.Warn("No gossip happening!")
}

// this test just has the purpose of verifying that the customGossiper works as intended
func TestCustomGossiper(t *testing.T) {
	vals, err := Validators(t, 6, "10e18,v,100,0.0.0.0:%s,%s")
	require.NoError(t, err)

	for _, val := range vals {
		val.TendermintServices = &node.TendermintServices{Gossiper: newCustomGossiper}
	}
	// creates a network of 6 vals and starts all the nodes in it
	network, err := NewNetworkFromValidators(t, vals, true)
	require.NoError(t, err)
	defer network.Shutdown()

	// network should be stuck since no peer is gossiping messages
	err = network.WaitForHeight(1, 10)
	require.Error(t, err)
}
