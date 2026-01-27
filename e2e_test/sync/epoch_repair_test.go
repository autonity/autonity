package sync

import (
	"testing"
	"time"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/rawdb"
	coretypes "github.com/autonity/autonity/core/types"
	e2e "github.com/autonity/autonity/e2e_test"
	"github.com/stretchr/testify/require"
)

func TestEpochHeadRepairAfterRestart(t *testing.T) {
	validators, err := e2e.Validators(t, 6, "10e18,v,1,127.0.0.1:%s,%s,%s,%s")
	require.NoError(t, err)
	genesis, err := e2e.Genesis(validators)
	require.NoError(t, err)

	enodes := make([]string, len(genesis.Config.AutonityContractConfig.Validators))
	for i, v := range genesis.Config.AutonityContractConfig.Validators {
		enodes[i] = v.Enode
	}
	staticNodes := coretypes.NewNodes(enodes, false).List

	network := make(e2e.Network, len(validators))
	for i, val := range validators {
		n, err := e2e.NewValidatorNode(t, val, genesis, i, false)
		require.NoError(t, err)
		n.Config.ExecutionP2P.StaticNodes = staticNodes
		require.NoError(t, n.Start())
		network[i] = n
	}
	defer network.Shutdown(t)

	epochPeriod := genesis.Config.AutonityContractConfig.EpochPeriod
	require.NoError(t, network.WaitToMineNBlocks(epochPeriod+5, int(epochPeriod*3), false))

	syncNode := network[len(network)-1]

	require.Eventually(t, func() bool {
		return rawdb.ReadEpochHeaderHash(syncNode.Eth.ChainDb()) != (common.Hash{})
	}, 30*time.Second, 500*time.Millisecond)

	syncNode.Eth.BlockChain().StopInsert()
	rawdb.WriteEpochHeaderHash(syncNode.Eth.ChainDb(), common.Hash{})
	require.Equal(t, common.Hash{}, rawdb.ReadEpochHeaderHash(syncNode.Eth.ChainDb()))

	require.NoError(t, syncNode.Close(false))
	syncNode.Wait()

	require.NoError(t, network.WaitToMineNBlocks(10, 60, false))

	require.NoError(t, syncNode.Start())

	require.Eventually(t, func() bool {
		return rawdb.ReadEpochHeaderHash(syncNode.Eth.ChainDb()) != (common.Hash{})
	}, 30*time.Second, 500*time.Millisecond)

	startHeight := syncNode.GetChainHeight()
	require.Eventually(t, func() bool {
		return syncNode.GetChainHeight() > startHeight
	}, 30*time.Second, 500*time.Millisecond)
}
