package e2e

import (
	"context"
	"testing"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/node"
	lru "github.com/hashicorp/golang-lru"
	"github.com/stretchr/testify/require"
)

type customGossiper struct {
	RecentMessages *lru.ARCCache  // the cache of peer's messages
	KnownMessages  *lru.ARCCache  // the cache of self messages
	Address        common.Address // address of the local peer
	Broadcaster    consensus.Broadcaster
	interfaces.Gossiper
}

func (cg *customGossiper) Gossip(_ context.Context, _ types.Committee, _ []byte) {
}

func TestCustomGossiper(t *testing.T) {
	vals, err := Validators(t, 6, "10e18,v,100,0.0.0.0:%s,%s")
	require.NoError(t, err)

	for _, val := range vals {
		val.TendermintServices = &node.TendermintServices{Gossiper: &customGossiper{}}
	}
	// creates a network of 6 vals and starts all the nodes in it
	network, err := NewNetworkFromValidators(t, vals, true)
	require.NoError(t, err)
	defer network.Shutdown()

	// network should be stuck since no peer is gossiping messages
	err = network.WaitForHeight(1, 5)
	require.Error(t, err)
}
