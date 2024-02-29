package topology

import (
	"fmt"
	"github.com/autonity/autonity/common/graph"
	e2e "github.com/autonity/autonity/e2e_test"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

// run network with different topology with gossiping, thus the network should have liveness at all case.
func TestTopology(t *testing.T) {
	topologies := []string{
		`graph TB
		V0---V1
		V0---V2
		V0---V3
		V0-->V4`,
		`graph TB
    	V0---V1
    	V2---V1
    	V3---V2
    	V4---V3`,
		`graph TB
		V0---V1
		V2---V1
		V3---V2
		V4---V3
		V0---V4`,
	}

	for _, tp := range topologies {
		runTopologyTest(t, tp)
	}
}

func runTopologyTest(t *testing.T, topology string) {
	network, err := e2e.NewNetwork(t, 5, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	defer network.Shutdown()

	// wait for the consensus engine to work.
	network.WaitToMineNBlocks(2, 10, false)

	graphStar, err := graph.Parse(strings.NewReader(topology))
	require.NoError(t, err)
	topologyManager := e2e.NewTopologyManager(graphStar)

	nodes := make(map[string]*e2e.Node, len(network))
	for i, node := range network {
		nodeID := fmt.Sprintf("V%d", i)
		nodes[nodeID] = node
	}
	err = topologyManager.ApplyTopologyOverNetwork(nodes)
	require.NoError(t, err)

	network.WaitToMineNBlocks(30, 30, false)
}
