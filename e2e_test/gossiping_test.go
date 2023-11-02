package e2e

import (
	"testing"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/node"
	lru "github.com/hashicorp/golang-lru"
	"github.com/stretchr/testify/require"
)

func newCustomGossiper(recentMessages *lru.ARCCache, knownMessages *lru.ARCCache, address common.Address) interfaces.Gossiper {
	return &customGossiper{recentMessages: recentMessages, knownMessages: knownMessages, address: address}
}

type customGossiper struct {
	recentMessages *lru.ARCCache
	knownMessages  *lru.ARCCache
	address        common.Address
	broadcaster    consensus.Broadcaster
}

func (cg *customGossiper) Gossip(_ types.Committee, _ []byte) {
	log.Warn("No gossip happening!")
}

func (cg *customGossiper) SetBroadcaster(b consensus.Broadcaster) {
	cg.broadcaster = b
}

// this test just has the purpose of verifying that the customGossiper works as intended
func TestCustomGossiper(t *testing.T) {
	vals, err := Validators(t, 6, "10e18,v,100,0.0.0.0:%s,%s")
	require.NoError(t, err)

	for _, val := range vals {
		val.TendermintServices = &node.TendermintServices{NewCustomGossiper: newCustomGossiper}
	}
	// creates a network of 6 vals and starts all the nodes in it
	network, err := NewNetworkFromValidators(t, vals, true)
	require.NoError(t, err)
	defer network.Shutdown()

	// network should be stuck since no peer is gossiping messages
	err = network.WaitForHeight(1, 10)
	require.Error(t, err)
}
