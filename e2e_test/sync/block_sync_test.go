package sync

import (
	"context"
	e2e "github.com/autonity/autonity/e2e_test"
	"github.com/autonity/autonity/eth/downloader"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestSnapSyncMode(t *testing.T) {
	network, err := e2e.NewNetwork(t, 7, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	defer network.Shutdown(t)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = network[0].SendAUTtracked(ctx, network[1].Address, 10)
	require.NoError(t, err)

	_ = network.WaitToMineNBlocks(100, 100, false)

	// create a node which runs in snap sync mode.
	snapSyncClients, err := e2e.Validators(t, 1, "10e18,v,10000,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	snapSyncNode, err := e2e.NewNoneValidatorNode(snapSyncClients[0], network[0].EthConfig.Genesis, len(network), downloader.SnapSync)
	require.NoError(t, err)
	err = snapSyncNode.Start()
	require.NoError(t, err)
	// Snap sync might take a while since it dumps and replicates entire world state.
	_ = network.WaitToMineNBlocks(100, 100, false)
	require.Equal(t, true, snapSyncNode.IsSyncComplete())
	require.True(t, true, snapSyncNode.GetChainHeight() > 0)
	_, parentEHead, curEHead, nextEHead, err := snapSyncNode.Eth.BlockChain().LatestEpoch()
	require.NoError(t, err)
	require.Greater(t, parentEHead, uint64(0))
	require.Greater(t, curEHead, parentEHead)
	require.NotEqual(t, nextEHead, curEHead)

	err = snapSyncNode.Close(true)
	require.NoError(t, err)
}

func TestFullSyncMode(t *testing.T) {
	network, err := e2e.NewNetwork(t, 7, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	defer network.Shutdown(t)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = network[0].SendAUTtracked(ctx, network[1].Address, 10)
	require.NoError(t, err)

	_ = network.WaitToMineNBlocks(100, 100, false)

	// create a node which runs in snap sync mode.
	snapSyncClients, err := e2e.Validators(t, 1, "10e18,v,10000,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	fullSyncNode, err := e2e.NewNoneValidatorNode(snapSyncClients[0], network[0].EthConfig.Genesis, len(network), downloader.FullSync)
	require.NoError(t, err)
	err = fullSyncNode.Start()
	require.NoError(t, err)
	// Snap sync might take a while since it dumps and replicates entire world state.
	_ = network.WaitToMineNBlocks(100, 100, false)
	require.Equal(t, true, fullSyncNode.IsSyncComplete())
	require.True(t, true, fullSyncNode.GetChainHeight() > 0)
	_, parentEHead, curEHead, nextEHead, err := fullSyncNode.Eth.BlockChain().LatestEpoch()
	require.NoError(t, err)
	require.Greater(t, parentEHead, uint64(0))
	require.Greater(t, curEHead, parentEHead)
	require.NotEqual(t, nextEHead, curEHead)

	err = fullSyncNode.Close(true)
	require.NoError(t, err)
}
