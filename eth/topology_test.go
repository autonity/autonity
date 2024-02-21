package eth

import (
	"crypto/ecdsa"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/p2p/enode"
	"github.com/autonity/autonity/params"
)

func TestEthExecutionLayerGraph(t *testing.T) {
	const targetDiameter = 4
	nodeCount := int(max(1000, params.TestAutonityContractConfig.MaxCommitteeSize))
	graph := NewBulkGraphTester(targetDiameter, nodeCount, t)
	for n := 1; n <= nodeCount; n++ {
		graph.AddNewNode()
		graph.TestGraph()
	}
}

// benchmark construction of the list of connection of a single node
func BenchmarkEdgeConstruction(b *testing.B) {
	const targetDiameter = 4
	nodeCount := int(max(1000, params.TestAutonityContractConfig.MaxCommitteeSize))
	graph := NewGraphTester(targetDiameter, nodeCount, b)
	// create the graph
	for i := 0; i < nodeCount; i++ {
		graph.AddNewNode()
	}
	// benchmark on edge construction for the last node
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		graph.topology.RequestSubset(graph.nodes, graph.localNodes[i%nodeCount])
	}
}

type graphTester struct {
	t              require.TestingT
	totalNodeCount int
	targetDiameter int
	topology       networkTopology
	nodes          []*enode.Node
	privateKeys    map[*ecdsa.PrivateKey]bool
	nodesIdx       map[*enode.Node]int
	connections    [][]*enode.Node
	distance       [][]int
	localNodes     []*enode.LocalNode
	// to check if edges are bidirectional
	edgeChecker [][]int
	// if bulkTest = true, the graph is tested for each node added, and some optimimzation must be applied
	// set bulkTest = false if a single graph is to be tested
	bulkTest bool
}

func (graph *graphTester) initiateGraph(targetDiameter, totalNodeCount int) {
	graph.totalNodeCount = totalNodeCount
	graph.targetDiameter = targetDiameter
	graph.topology = NewGraphTopology(0)
	graph.nodes = make([]*enode.Node, 0, totalNodeCount)
	graph.privateKeys = make(map[*ecdsa.PrivateKey]bool)
	graph.nodesIdx = make(map[*enode.Node]int)
	graph.localNodes = make([]*enode.LocalNode, 0, totalNodeCount)
	graph.edgeChecker = make([][]int, totalNodeCount)
	for i := 0; i < totalNodeCount; i++ {
		graph.edgeChecker[i] = make([]int, totalNodeCount)
	}
	graph.distance = make([][]int, totalNodeCount)
	for i := 0; i < totalNodeCount; i++ {
		graph.distance[i] = make([]int, totalNodeCount)
		for j := 0; j < totalNodeCount; j++ {
			graph.distance[i][j] = -1
		}
	}
	graph.connections = make([][]*enode.Node, totalNodeCount)
}

func NewBulkGraphTester(targetDiameter int, totalNodeCount int, t require.TestingT) graphTester {
	graph := graphTester{
		t:        t,
		bulkTest: true,
	}
	graph.initiateGraph(targetDiameter, totalNodeCount)
	return graph
}

func NewGraphTester(targetDiameter int, totalNodeCount int, t require.TestingT) graphTester {
	graph := graphTester{
		t:        t,
		bulkTest: false,
	}
	graph.initiateGraph(targetDiameter, totalNodeCount)
	return graph
}

func (graph *graphTester) AddNewNode() {
	for {
		privateKey, err := crypto.GenerateKey()
		require.NoError(graph.t, err)
		if _, ok := graph.privateKeys[privateKey]; !ok {
			newEnode := "enode://" + string(crypto.PubECDSAToHex(&privateKey.PublicKey)[2:]) + "@3.209.45.79:30303"
			newNode, err := enode.ParseV4(newEnode)
			require.NoError(graph.t, err)
			require.NotEqual(graph.t, nil, newNode)
			graph.nodesIdx[newNode] = len(graph.nodes)
			graph.nodes = append(graph.nodes, newNode)
			// related localNode
			db, err := enode.OpenDB("")
			require.NoError(graph.t, err)
			localNode := enode.NewLocalNode(db, privateKey, nil)
			require.Equal(graph.t, newNode.ID(), localNode.ID())
			graph.localNodes = append(graph.localNodes, localNode)
			graph.privateKeys[privateKey] = true
			break
		}
	}
	if !graph.bulkTest && len(graph.nodes) < graph.totalNodeCount {
		return
	}
	testID := len(graph.nodes)
	task := sync.WaitGroup{}
	for i := 0; i < len(graph.nodes); i++ {
		task.Add(1)
		go func(idx int) {
			edges := graph.topology.RequestSubset(graph.nodes, graph.localNodes[idx])
			for _, peer := range edges {
				peerIdx := graph.nodesIdx[peer]
				// put unidirectional edge, i.e. idx -> peerIdx
				// if edges are bidirectional, we will have peerIdx -> idx
				graph.edgeChecker[idx][peerIdx] = testID
			}
			graph.connections[idx] = edges
			task.Done()
		}(i)
	}
	task.Wait()
	// check if edges are bidirectional
	for i := 0; i < len(graph.nodes); i++ {
		for _, peer := range graph.connections[i] {
			peerIdx := graph.nodesIdx[peer]
			require.Equal(graph.t, testID, graph.edgeChecker[peerIdx][i])
		}
		require.Equal(graph.t, 0, graph.edgeChecker[i][i], "self loop detected")
	}
}

func (graph *graphTester) testGraphDegree() {
	// check if the degree properties hold
	for i := 0; i < len(graph.nodes); i++ {
		require.True(graph.t, len(graph.connections[i]) <= MaxDegree)
	}
}

func (graph *graphTester) TestGraph() {
	graph.testGraphDegree()
	if !graph.bulkTest {
		graph.testGraphDiamter()
	} else if len(graph.nodes)%100 == 0 {
		graph.testGraphDiamter()
	} else {
		// check if graph is connected
		visited := make([]bool, len(graph.nodes))
		graph.dfs(0, visited)
		for _, check := range visited {
			require.True(graph.t, check, "graph disconnected")
		}
	}
}

func (graph *graphTester) dfs(nodeIdx int, visited []bool) {
	if visited[nodeIdx] {
		return
	}
	visited[nodeIdx] = true
	for _, peer := range graph.connections[nodeIdx] {
		graph.dfs(graph.nodesIdx[peer], visited)
	}
}

func (graph *graphTester) testGraphDiamter() {
	for i := 0; i < len(graph.nodes); i++ {
		graph.bfs(i, graph.distance[i])
		for j := 0; j < len(graph.nodes); j++ {
			d := graph.distance[i][j]
			require.True(graph.t, d >= 0, "graph disconnected")
			require.True(graph.t, d <= graph.targetDiameter, "graph diameter more than expected")
		}
	}
}

func (graph *graphTester) bfs(sourceIdx int, dis []int) {
	for i := 0; i < len(graph.nodes); i++ {
		dis[i] = -1
	}
	// enque source
	queue := make([]int, 0, len(graph.nodes))
	queue = append(queue, sourceIdx)
	dis[sourceIdx] = 0
	for len(queue) > 0 {
		// pop
		nodeIdx := queue[0]
		queue = queue[1:]
		for _, peer := range graph.connections[nodeIdx] {
			peerIdx := graph.nodesIdx[peer]
			if dis[peerIdx] < 0 {
				// enque adjacent nodes
				queue = append(queue, peerIdx)
				dis[peerIdx] = dis[nodeIdx] + 1
			}
		}
	}
}
