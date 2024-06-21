package sync

import (
	"context"
	e2e "github.com/autonity/autonity/e2e_test"
	"github.com/autonity/autonity/eth/downloader"
	"github.com/autonity/autonity/params"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestSnapSyncMode(t *testing.T) {
	testSyncMode(t, downloader.SnapSync)
}

func TestFullSyncMode(t *testing.T) {
	testSyncMode(t, downloader.FullSync)
}

func testSyncMode(t *testing.T, mode downloader.SyncMode) {
	network, err := e2e.NewNetwork(t, 7, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	defer network.Shutdown(t)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = network[0].SendAUTtracked(ctx, network[1].Address, 10)
	require.NoError(t, err)

	_ = network.WaitToMineNBlocks(100, 100, false)

	// create a node which runs in the specified sync mode.
	identities, err := e2e.Validators(t, 1, "10e18,v,10000,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	syncNode, err := e2e.NewNoneValidatorNode(identities[0], network[0].EthConfig.Genesis, len(network), mode)
	require.NoError(t, err)
	err = syncNode.Start()
	require.NoError(t, err)
	// Snap sync might take a while since it dumps and replicates entire world state.
	_ = network.WaitToMineNBlocks(100, 100, false)
	require.Equal(t, true, syncNode.IsSyncComplete())
	require.True(t, true, syncNode.GetChainHeight() > 0)
	epoch, err := syncNode.Eth.BlockChain().LatestEpoch()
	require.NoError(t, err)
	require.Greater(t, epoch.PreviousEpochBlock.Uint64(), uint64(0))
	require.Greater(t, epoch.EpochBlock.Uint64(), epoch.PreviousEpochBlock.Uint64())
	require.NotEqual(t, epoch.NextEpochBlock.Uint64(), epoch.EpochBlock.Uint64())
	require.Equal(t, params.DefaultOmissionAccountabilityConfig.Delta, epoch.Delta.Uint64())

	err = syncNode.Close(true)
	require.NoError(t, err)
}
