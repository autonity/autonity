package resetting

import (
	e2e "github.com/autonity/autonity/e2e_test"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStateRecoveryFromWAL(t *testing.T) {
	numOfNodes := 4
	network, err := e2e.NewNetwork(t, numOfNodes, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	defer network.Shutdown(t)
	// wait for the consensus engine to work.
	network.WaitToMineNBlocks(60, 60, false)

	// stop the first 2 nodes.
	err = network[0].Close(false)
	network[0].Wait()
	require.NoError(t, err)
	err = network[1].Close(false)
	network[1].Wait()
	require.NoError(t, err)

	// network should be on-hold
	network.WaitToMineNBlocks(10, 60, false)

	// recovery the 1st 2 nodes.
	err = network[0].Start()
	require.NoError(t, err)
	err = network[1].Start()
	require.NoError(t, err)

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(60, 60, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}
