package clustering

import (
	"github.com/autonity/autonity/e2e_test"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClusteringHappyCase(t *testing.T) {

	// todo: write helper functions to create network with customized pinger.

	network, err := e2e.NewNetwork(t, 24, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	defer network.Shutdown(t)

	// runs for about 10 epoches period.
	_ = network.WaitToMineNBlocks(500, 500, false)
}
