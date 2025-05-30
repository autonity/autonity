package eth

import (
	"crypto/ecdsa"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/p2p/enode"
	"github.com/autonity/autonity/params"
)

var committeeSizeToCheck = []int{50, 64, 65, 100, 200, 500, 1000}

func TestEdgeDirection(t *testing.T) {
	// test if the edges are unidirectional or bidirectional
	// test should fail if edges are not bidirectional
	nodeCount := 1000
	nodes := make([]*enode.Node, 0, nodeCount)
	privateKeys := make(map[*ecdsa.PrivateKey]bool)
	edgeChecker := make([][]int, nodeCount)
	nodesIndex := make(map[*enode.Node]int)
	for i := 0; i < nodeCount; i++ {
		edgeChecker[i] = make([]int, nodeCount)
	}
	connections := make([][]*enode.Node, nodeCount)
	topology := NewGraphTopology(0)
	for i := 0; i < nodeCount; i++ {
		testID := i + 1
		privateKey, newNode := createNewNode(t, privateKeys)
		privateKeys[privateKey] = true
		nodes = append(nodes, newNode)
		nodesIndex[newNode] = i
		for myIndex := 0; myIndex < len(nodes); myIndex++ {
			edges := topology.RequestSubset(nodes, myIndex)
			for _, peer := range edges {
				peerIndex := nodesIndex[peer]
				edgeChecker[myIndex][peerIndex] = testID
			}
			connections[myIndex] = edges
		}
		for myIndex := 0; myIndex < len(nodes); myIndex++ {
			for _, peer := range connections[myIndex] {
				peerIndex := nodesIndex[peer]
				require.True(t, edgeChecker[peerIndex][myIndex] == testID)
			}
			require.True(t, edgeChecker[myIndex][myIndex] == 0)
		}
	}
}

func TestGraphDegree(t *testing.T) {
	const targetDiameter = 4
	nodeCount := int(max(1000, params.TestAutonityContractConfig.MaxCommitteeSize))
	graph := NewBulkGraphTester(targetDiameter, nodeCount, t)
	for n := 1; n <= nodeCount; n++ {
		graph.AddNewNode()
		printDegree := false
		for _, committeeSize := range committeeSizeToCheck {
			if n == committeeSize {
				printDegree = true
				break
			}
		}
		graph.TestGraphDegree(printDegree)
	}
}

func TestGraphDiamter(t *testing.T) {
	const targetDiameter = 4
	nodeCount := int(max(1000, params.TestAutonityContractConfig.MaxCommitteeSize))
	graph := NewBulkGraphTester(targetDiameter, nodeCount, t)
	for n := 1; n <= nodeCount; n++ {
		graph.AddNewNode()
		// tests if the graph diameter <= target diameter for the whole graph
		// which ensures that the whole graph is connected
		graph.TestGraphDiamter()
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
		graph.topology.RequestSubset(graph.nodes, i%nodeCount)
	}
}

type graphTester struct {
	t              require.TestingT
	totalNodeCount int
	targetDiameter int
	topology       networkTopology
	nodes          []*enode.Node
	privateKeys    map[*ecdsa.PrivateKey]bool
	nodesIndex     map[*enode.Node]int
	connections    [][]*enode.Node
	distance       [][]int
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
	graph.nodesIndex = make(map[*enode.Node]int)
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

func createNewNode(t require.TestingT, privateKeys map[*ecdsa.PrivateKey]bool) (*ecdsa.PrivateKey, *enode.Node) {
	if privateKeys == nil {
		privateKeys = make(map[*ecdsa.PrivateKey]bool)
	}
	for {
		privateKey, err := crypto.GenerateKey()
		require.NoError(t, err)
		if _, ok := privateKeys[privateKey]; !ok {
			newEnode := "enode://" + string(crypto.PubECDSAToHex(&privateKey.PublicKey)[2:]) + "@3.209.45.79:30303"
			newNode, err := enode.ParseV4(newEnode)
			require.NoError(t, err)
			require.NotEqual(t, nil, newNode)
			return privateKey, newNode
		}
	}
}

func (graph *graphTester) AddNewNode() {
	privateKey, newNode := createNewNode(graph.t, graph.privateKeys)
	graph.nodesIndex[newNode] = len(graph.nodes)
	graph.nodes = append(graph.nodes, newNode)
	graph.privateKeys[privateKey] = true
	if !graph.bulkTest && len(graph.nodes) < graph.totalNodeCount {
		return
	}
	for i := 0; i < len(graph.nodes); i++ {
		edges := graph.topology.RequestSubset(graph.nodes, i)
		graph.connections[i] = edges
	}
}

func (graph *graphTester) TestGraphDegree(printDegree bool) {
	// check if the degree properties hold
	if printDegree {
		fmt.Printf("\n\ndegree check for committee size %v\n\n", len(graph.nodes))
	}
	for i := 0; i < len(graph.nodes); i++ {
		require.True(graph.t, len(graph.connections[i]) <= MaxDegree)

		if printDegree {
			fmt.Printf("node %v degree count : %v\n", i, len(graph.connections[i]))
		}
	}
}

// It returns a number of 2 digits in b-base number system: i*b + j where i <= j is maintained
// i <= j condition is maintained so that combinedIndex(i,j,b) = combinedIndex(j,i,b)
// This is used to convert a pair (a,b) where (0 <= a,b < n) to a single number
// So we get combinedIndex(a,b,n) = combinedIndex(b,a,n) = c which represents the pair (a,b) or (b,a)
func combinedIndex(i, j, b int) int {
	if i > j {
		return combinedIndex(j, i, b)
	}
	return i*b + j
}

func (graph *graphTester) TestGraphDiamter() {
	totalNodes := len(graph.nodes)
	for i := 0; i < totalNodes; i++ {
		for j := 0; j < i; j++ {
			graph.distance[i][j] = graph.targetDiameter + 1
			graph.distance[j][i] = graph.targetDiameter + 1
		}
		graph.distance[i][i] = 0
	}
	pairsToUpdate := totalNodes * (totalNodes - 1) / 2 // we have C(n,2) unordered pairs of nodes
	updatedPairs := make(map[int]bool)
	for nodeCount := len(graph.nodes); nodeCount > 0 && len(updatedPairs) < pairsToUpdate; nodeCount-- {
		source := nodeCount - 1
		// bfs is modified to determine shortest path distance from source only if graph diameter <= targetDiameter
		// in case graph diameter > targetDiameter, bfs will not give shortest path distance and some pair (i,j) will have
		// distance[i][j] = targetDiameter + 1 and the test will fail
		graph.bfs(source, nodeCount, graph.distance[source])
		distantNodes := make([][]int, graph.targetDiameter)
		for i := 1; i < graph.targetDiameter; i++ {
			distantNodes[i] = make([]int, 0, nodeCount)
		}
		for j := 0; j < source; j++ {
			d := graph.distance[source][j]
			require.True(graph.t, d <= graph.targetDiameter, "graph diameter more than expected")
			// assuming that the graph is bidirectional, which is tested in TestEdgeDirection
			graph.distance[j][source] = d
			updatedPairs[combinedIndex(j, source, totalNodes)] = true
			if d < graph.targetDiameter {
				distantNodes[d] = append(distantNodes[d], j)
			}
		}
		// update any pair (nodeA,nodeB) such that the shortest path between nodeA and nodeB includes source
		// As our targetedDiameter = 4, doing this operation is not very costly.
		for d := 1; d < graph.targetDiameter; d++ {
			for _, nodeAIndex := range distantNodes[d] {
				for d1 := d; d1+d <= graph.targetDiameter; d1++ {
					for _, nodeBIndex := range distantNodes[d1] {
						if nodeBIndex != nodeAIndex {
							// no need to check distance here as d+d1 <= targetDiameter
							graph.distance[nodeAIndex][nodeBIndex] = min(graph.distance[nodeAIndex][nodeBIndex], d+d1)
							graph.distance[nodeBIndex][nodeAIndex] = min(graph.distance[nodeBIndex][nodeAIndex], d+d1)
							updatedPairs[combinedIndex(nodeAIndex, nodeBIndex, totalNodes)] = true
						}
					}
				}
			}
		}
	}
}

// Here nodeCount <= totalNodes, the number of nodes we are testing
// Lets say for all nodes with ID >= nodeCount, we have distance[i][ID] and distance[ID][i] updated for all 0 <= i < totalNodes
// Now we want to update distance for all pair (i,sourceIndex) where 0 <= i,sourceIndex < nodeCount for all 0 <= i < nodeCount
// If the shortest path between i and sourceIndex includes any node j where j >= nodeCount, then distance[i][sourceIndex] can be updated
// via j, i.e. distance[i][sourceIndex] = distance[i][j] + distance[j][sourceIndex] which is already done before calling bfs.
// So in bfs we don't need to include path which has some intermediate node, x such that x >= nodeCount
func (graph *graphTester) bfs(sourceIndex, nodeCount int, dis []int) {
	// enque source
	queue := make([]int, 0, nodeCount)
	queue = append(queue, sourceIndex)
	dis[sourceIndex] = 0
	for len(queue) > 0 {
		// pop
		nodeIndex := queue[0]
		queue = queue[1:]
		for _, peer := range graph.connections[nodeIndex] {
			peerIndex := graph.nodesIndex[peer]
			if peerIndex < nodeCount && dis[peerIndex] > dis[nodeIndex]+1 {
				// enque adjacent nodes
				queue = append(queue, peerIndex)
				dis[peerIndex] = dis[nodeIndex] + 1
			}
		}
	}
}
