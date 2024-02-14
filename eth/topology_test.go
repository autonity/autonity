package eth

import (
	"crypto/ecdsa"
	"math"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/p2p/enode"
	"github.com/autonity/autonity/params"
)

// currently we only support graph diameter = 2, to support diameter > 2, floydWarshall needs to be modified
func TestEthExecutionLayerGraph(t *testing.T) {
	const targetDiameter = 2
	nodeCount := int(max(1000, params.TestAutonityContractConfig.MaxCommitteeSize))
	graph := NewBulkGraphTester(targetDiameter, nodeCount, t)
	for n := 1; n <= nodeCount; n++ {
		graph.AddNewNode()
		graph.TestGraph()
	}
	t.Logf("heavy tests done %v", graph.heavyTestCount)
}

// benchmark construction of the list of connection of a single node
func BenchmarkEdgeConstruction(b *testing.B) {
	const targetDiameter = 2
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
	base           int
	topology       networkTopology
	nodes          []*enode.Node
	privateKeys    map[*ecdsa.PrivateKey]bool
	nodesIdx       map[*enode.Node]int
	connections    [][]*enode.Node
	distance       [][]int
	localNodes     []*enode.LocalNode
	graphChanged   bool
	heavyTestCount int
	// if bulkTest = true, the graph is tested for each node added, and some optimimzation must be applied
	// set bulkTest = false if a single graph is to be tested
	bulkTest bool
}

func (graph *graphTester) initiateGraph(targetDiameter, totalNodeCount int) {
	// following functions need to be updated to support testing diameter > 2: floydWarshall, testNodeDiameter
	require.True(graph.t, targetDiameter <= 2, "testing of diameter > 2 not supported")
	graph.totalNodeCount = totalNodeCount
	graph.targetDiameter = targetDiameter
	graph.topology = NewGraphTopology(targetDiameter, 0)
	graph.nodes = make([]*enode.Node, 0, totalNodeCount)
	graph.privateKeys = make(map[*ecdsa.PrivateKey]bool)
	graph.nodesIdx = make(map[*enode.Node]int)
	graph.localNodes = make([]*enode.LocalNode, 0, totalNodeCount)
	graph.distance = make([][]int, totalNodeCount)
	for i := 0; i < totalNodeCount; i++ {
		graph.distance[i] = make([]int, totalNodeCount)
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
	graph.graphChanged = false
	edgeAdded := make([]bool, len(graph.nodes))
	task := sync.WaitGroup{}
	for i := 0; i < len(graph.nodes); i++ {
		task.Add(1)
		go func(idx int) {
			defer task.Done()
			edges := graph.topology.RequestSubset(graph.nodes, graph.localNodes[idx])
			edgeAdded[idx] = true
			if !graph.bulkTest {
				graph.connections[idx] = edges
				return
			}
			// for bulk test, we have a new graph created each time a new node is added
			// to optimze testing of all the graphs, we track the following:
			// graph.graphChanged = true if adding the new node changes at least one edge present in the current graph (before adding the new node)
			// otherwise graph.graphChanged = false, i.e. adding the new node only adds some eges to graph, it does not change any existing edges
			if !graph.graphChanged {
				newNodes := make(map[*enode.Node]bool)
				for _, peer := range edges {
					newNodes[peer] = true
				}
				for _, peer := range graph.connections[idx] {
					if newNodes[peer] == false {
						// current egde not found in the set of new edges
						graph.graphChanged = true
						break
					}
				}
			}
			graph.connections[idx] = edges
		}(i)
	}
	task.Wait()
	for _, check := range edgeAdded {
		require.True(graph.t, check)
	}
}

func (graph *graphTester) testGraphDegree() {
	// check if the degree properties hold
	graph.base = graph.topology.ComputeBase(len(graph.nodes))
	// check if base calculation correct or not
	require.True(
		graph.t, int(math.Pow(float64(graph.base), float64(graph.targetDiameter))) >= len(graph.nodes) &&
			int(math.Pow(float64(graph.base-1), float64(graph.targetDiameter))) < len(graph.nodes),
	)
	maxDegree := graph.topology.MaxDegree(len(graph.nodes))
	for i := 0; i < len(graph.nodes); i++ {
		require.True(graph.t, len(graph.connections[i]) <= maxDegree)
	}
}

func (graph *graphTester) TestGraph() {
	graph.testGraphDegree()
	if !graph.bulkTest {
		graph.testGraphDiamter()
	} else if graph.graphChanged {
		// test the whole graph because the old edges from the graph changed
		graph.t.(*testing.T).Logf("heavy test on graph with nodes %v", len(graph.nodes))
		graph.heavyTestCount++
		graph.testGraphDiamter()
		graph.graphChanged = false
	} else {
		// test a single node, no need to test the whole graph because the old edges remain same
		graph.testLastNode()
	}
}

func (graph *graphTester) testLastNode() {
	myIdx := len(graph.nodes) - 1
	// update distance from any node to myIdx node
	for i := 0; i < len(graph.nodes); i++ {
		graph.distance[i][myIdx] = graph.targetDiameter + 1 // if dis[i][j] > graph.targetDiameter, test fails
		graph.distance[myIdx][i] = graph.targetDiameter + 1
	}
	graph.distance[myIdx][myIdx] = 0
	base := graph.topology.ComputeBase(len(graph.nodes))
	for _, peer := range graph.connections[myIdx] {
		i := graph.nodesIdx[peer]
		// check graph construction
		require.True(graph.t, graph.topology.countMatchingDigits(i, myIdx, base) == 1, "invalid graph constrcution")
		graph.distance[myIdx][i] = 1
		graph.distance[i][myIdx] = 1
		for _, distantPeer := range graph.connections[i] {
			j := graph.nodesIdx[distantPeer]
			if graph.distance[myIdx][j] > 2 {
				graph.distance[myIdx][j] = 2
				graph.distance[j][myIdx] = 2
			}
		}
		// If there is any node i remains where dis[i][myIdx] is not updated yet, then dis[i][myIdx] > 2 is true for them.
		// As we support only diameter = 2, we don't need to update those distances and the test will rightfully fail.
	}
	// Ideally we should update all pair distance becase there could be a new path including the new node that reduces their distance,
	// but it cannot be reduced less than 2. For now we only support diameter = 2, so we can do a little optimization here.
	// We can ignore all pair distance update because we already have all pair distance <= 2 (otherwise the test would fail already)
	// and it cannot be reduced any more.
	// Note that the above statement would not be true if we had diameter > 2
	for i := 0; i < len(graph.nodes); i++ {
		require.True(graph.t, graph.distance[i][myIdx] <= graph.targetDiameter, "graph diameter more than desired")
		require.True(graph.t, graph.distance[myIdx][i] <= graph.targetDiameter, "graph diameter more than desired")
	}
}

func (graph *graphTester) testGraphDiamter() {
	// check if the construction is correct
	base := graph.topology.ComputeBase(len(graph.nodes))
	for i := 0; i < len(graph.nodes); i++ {
		for _, peer := range graph.connections[i] {
			j := graph.nodesIdx[peer]
			require.True(graph.t, graph.topology.countMatchingDigits(i, j, base) == 1, "invalid graph constrcution")
		}
	}
	// check distance for each pair of nodes
	graph.floydWarshall()
	for i := 0; i < len(graph.nodes); i++ {
		for j := 0; j < len(graph.nodes); j++ {
			require.True(graph.t, graph.distance[i][j] <= graph.targetDiameter, "graph diameter more than desired")
		}
	}
}

// algorithm to measure distance between all pair of nodes in a graph: https://en.wikipedia.org/wiki/Floyd%E2%80%93Warshall_algorithm
func (graph *graphTester) floydWarshall() {
	for i := 0; i < len(graph.nodes); i++ {
		for j := 0; j < len(graph.nodes); j++ {
			graph.distance[i][j] = graph.targetDiameter + 1 // if dis[i][j] > graph.targetDiameter, test fails
		}
	}
	for i := 0; i < len(graph.nodes); i++ {
		graph.distance[i][i] = 0
		for _, peer := range graph.connections[i] {
			j := graph.nodesIdx[peer]
			graph.distance[i][j] = 1
		}
	}
	for mid := 0; mid < len(graph.nodes); mid++ {
		/*
			According to the Floyd-Warshall algorithm, the following part should iterate through all possible pair of nodes (i,j)
			and update their distance if it reduces using the midle node, i.e. dis[i][j] = min(dis[i][j], dis[i][mid] + dis[mid][j]).
			But we support only diameter = 2, so we can do the following optimization
			1. For any node i that doesn't have a direct connection to mid, dis[i][mid] > 1 is always true
			2. For any pair of nodes (i,j), if either i or j does not have a direct connection,
				then dis[i][j] = dis[i][mid] + dis[mid][j] > 2. If there is no such node mid found, then the test rightfully fails
			3. For any pair of nodes (i,j), dis[i][j] <= 2 if they have direct connection or there is a node, mid,
				which has direct connection to both i and j
			4. So we only update the pair of nodes (i,j) if both i and j have direct connection to node mid
			5. Note that if we have some pair nodes (i,j) such that in the graph the actual distance between them is greater than 2,
				then we cannot determine their actual distance via this optimization, but we don't care since the test will fail anyway
		*/
		edges := graph.connections[mid]
		for _, peer1 := range edges {
			for _, peer2 := range edges {
				i := graph.nodesIdx[peer1]
				j := graph.nodesIdx[peer2]
				if graph.distance[i][j] > 2 {
					graph.distance[i][j] = 2
				}
			}
		}
	}
}
