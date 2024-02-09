package eth

import (
	"crypto/ecdsa"
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/p2p/enode"
	"github.com/autonity/autonity/params"
)

func bfs(
	source *enode.Node, dis map[*enode.Node]int, localNodes map[*enode.Node]*enode.LocalNode,
	nodes []*enode.Node, adjacentNodes func([]*enode.Node, *enode.LocalNode) []*enode.Node,
) {
	for key := range dis {
		dis[key] = -1
	}
	// enque source
	dis[source] = 0
	queue := make([]*enode.Node, 0, len(nodes))
	queue = append(queue, source)
	for len(queue) > 0 {
		// pop
		node := queue[0]
		queue = queue[1:]
		connections := adjacentNodes(nodes, localNodes[node])
		for _, peer := range connections {
			if d, ok := dis[peer]; !ok || d < 0 {
				// enque adjacent nodes
				dis[peer] = dis[node] + 1
				queue = append(queue, peer)
			}
		}
	}
}

func TestEthExecutionLayerGraph(t *testing.T) {
	const targetDiameter = 2
	tSelector := &networkTopology{}
	tSelector.SetDiameter(targetDiameter)
	tSelector.SetMinNodes(0)
	db, err := enode.OpenDB("")
	require.NoError(t, err)
	nodeCount := int(max(100, params.TestAutonityContractConfig.MaxCommitteeSize))
	nodes := make([]*enode.Node, 0, nodeCount)
	localNodes := make(map[*enode.Node]*enode.LocalNode)
	privateKeys := make(map[*ecdsa.PrivateKey]bool)
	for n := 1; n <= nodeCount; n++ {
		fmt.Printf("\ngraph %v starts\n", n)
		for {
			privateKey, err := crypto.GenerateKey()
			require.NoError(t, err)
			if _, ok := privateKeys[privateKey]; !ok {
				newEnode := "enode://" + string(crypto.PubECDSAToHex(&privateKey.PublicKey)[2:]) + "@3.209.45.79:30303"
				newNode, err := enode.ParseV4(newEnode)
				require.NoError(t, err)
				require.NotEqual(t, nil, newNode)
				nodes = append(nodes, newNode)
				// related localNode
				// db is not used here, so a single db for all nodes
				localNode := enode.NewLocalNode(db, privateKey, nil)
				localNodes[newNode] = localNode
				privateKeys[privateKey] = true
				require.Equal(t, newNode.ID(), localNode.ID())
				break
			}
		}
		// check if graph connected and the diameter and degree properties hold
		base := tSelector.ComputeBase(uint(len(nodes)))
		require.True(
			t, int(math.Pow(float64(base), targetDiameter)) >= len(nodes) &&
				int(math.Pow(float64(base-1), targetDiameter)) < len(nodes),
		)
		maxDegree := targetDiameter * int(math.Pow(float64(base-1), targetDiameter-1))
		for _, node := range nodes {
			// check if max distance from node to any other node in the graph is targetDiameter
			dis := make(map[*enode.Node]int)
			bfs(node, dis, localNodes, nodes, tSelector.RequestSubset)
			for _, peer := range nodes {
				d, ok := dis[peer]
				require.True(t, ok && d >= 0, "graph not connected")
				require.True(t, d <= targetDiameter, "graph diameter more than desired")
			}
			connections := tSelector.RequestSubset(nodes, localNodes[node])
			require.True(t, len(connections) <= maxDegree)
		}
	}
}
