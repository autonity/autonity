package topology

import (
	"math"
	"math/big"
	"testing"

	"github.com/autonity/autonity/consensus/tendermint/bft"
	"github.com/stretchr/testify/require"
)

const (
	MinGraphSize = 1
	MaxGraphSize = 1000
)

func TestDegreeIsMoreThanOneThirdForFaultTolerant(t *testing.T) {
	tester := func(nodeCount int) {
		smallestDegreeChecker := func(graph Topology, smallestDegree int) {
			require.True(t, graph.isNetworkValid())
			// smallest possible degree count
			require.True(t, graph.DegreeCount() >= smallestDegree)
			require.True(t, graph.DegreeCount()-2 < smallestDegree)

			graph.SetDegree(graph.DegreeCount() - 2)

			require.True(t, graph.isNetworkValid())
			// smallest possible degree count
			require.True(t, graph.DegreeCount() >= smallestDegree)
			require.True(t, graph.DegreeCount()-2 < smallestDegree)
		}

		graph, err := NewFaultTolerantTopology(nodeCount, 0, 0)
		require.NoError(t, err)
		smallestDegreeChecker(graph, (nodeCount-1)/3+1)

		graph, err = NewOptimizedTopology(nodeCount, 0, 0)
		require.NoError(t, err)
		smallestDegreeChecker(graph, min(nodeCount-1, int(math.Ceil(math.Sqrt(float64(nodeCount))))))
	}
	testForMultipleGraph(2, MaxGraphSize, tester)
}

func TestDegreeCountNotMoreThanNetworkSize(t *testing.T) {
	tester := func(graph Topology) {
		err := graph.SetDegree(10)
		require.Error(t, err)
		require.Equal(t, errDegreeCount, err)

		require.Equal(t, 9, graph.DegreeCount())
	}

	_, err := NewFaultTolerantTopology(10, 10, 1)
	require.Error(t, err)
	require.Equal(t, errDegreeCount, err)

	_, err = NewOptimizedTopology(10, 10, 1)
	require.Error(t, err)
	require.Equal(t, errDegreeCount, err)

	graph, err := NewFaultTolerantTopology(10, 9, 1)
	require.NoError(t, err)
	tester(graph)

	graph, err = NewOptimizedTopology(10, 9, 1)
	require.NoError(t, err)
	tester(graph)
}

func TestAdjacentNodesAreConnected(t *testing.T) {
	tester := func(nodeCount int) {
		adjacencyCheck := func(graph Topology) {
			for node := 0; node < nodeCount; node++ {
				edges, err := graph.RequestSubset(node)
				require.NoError(t, err)
				for _, peer := range edges {
					require.True(t, graph.isAdjacent(node, peer))
				}
			}
		}

		graph, err := NewFaultTolerantTopology(nodeCount, 0, 0)
		require.NoError(t, err)
		adjacencyCheck(graph)

		graph, err = NewOptimizedTopology(nodeCount, 0, 0)
		require.NoError(t, err)
		adjacencyCheck(graph)
	}

	testForMultipleGraph(MinGraphSize, MaxGraphSize, tester)
}

func TestFaultTolerantGraphIsBidirectional(t *testing.T) {
	tester := func(nodeCount int) {
		graph, err := NewFaultTolerantTopology(nodeCount, 0, 0)
		require.NoError(t, err)

		edgesMap := make(map[int]bool)
		for node := 0; node < nodeCount; node++ {
			edges, err := graph.RequestSubset(node)
			require.NoError(t, err)
			for _, peer := range edges {
				key := combinedIndex(node, peer, nodeCount)
				if _, ok := edgesMap[key]; ok {
					delete(edgesMap, key)
				} else {
					edgesMap[key] = true
				}
			}
		}

		require.Equal(t, 0, len(edgesMap))
	}

	testForMultipleGraph(MinGraphSize, MaxGraphSize, tester)
}

func TestDegreeCount(t *testing.T) {
	tester := func(nodeCount int) {
		degreeChecker := func(graph Topology, nodeCount int) {
			for node := 0; node < nodeCount; node++ {
				edges, err := graph.RequestSubset(node)
				require.NoError(t, err)
				if graph.isAdjacent(node, node) {
					require.Equal(t, graph.DegreeCount(), len(edges)+1)
				} else {
					require.Equal(t, graph.DegreeCount(), len(edges))
				}
			}
		}

		graph, err := NewFaultTolerantTopology(nodeCount, 0, 0)
		require.NoError(t, err)
		degreeChecker(graph, nodeCount)

		graph, err = NewOptimizedTopology(nodeCount, 0, 0)
		require.NoError(t, err)
		degreeChecker(graph, graph.DegreeCount()+1)
		for node := graph.DegreeCount() + 1; node < nodeCount; node++ {
			edges, err := graph.RequestSubset(node)
			require.NoError(t, err)
			require.Equal(t, 0, len(edges))
		}
	}

	testForMultipleGraph(MinGraphSize, MaxGraphSize, tester)
}

func TestGraphIsConnected(t *testing.T) {
	tester := func(nodeCount int) {
		graph := newGraph(nodeCount)
		connectionChecker := func(graphTopology Topology) {
			var err error
			for node := 0; node < nodeCount; node++ {
				graph.edges[node], err = graphTopology.RequestSubset(node)
				require.NoError(t, err)
			}
			graph.visited = make([]bool, nodeCount)
			require.Equal(t, nodeCount, graph.dfs(0))
			for node := 0; node < nodeCount; node++ {
				require.True(t, graph.visited[node])
			}
		}

		graphTopology, err := NewFaultTolerantTopology(nodeCount, 0, 0)
		require.NoError(t, err)
		connectionChecker(graphTopology)

		graphTopology, err = NewOptimizedTopology(nodeCount, 0, 0)
		require.NoError(t, err)
		connectionChecker(graphTopology)
	}

	testForMultipleGraph(MinGraphSize, MaxGraphSize, tester)
}

// Tests graph diameter for a bidirectional graph
func TestFaultTolerantGraphDiameter(t *testing.T) {
	tester := func(nodeCount int) {
		graph := newGraph(nodeCount)
		graphTopology, err := NewFaultTolerantTopology(nodeCount, 0, 0)
		require.NoError(t, err)
		for node := 0; node < nodeCount; node++ {
			graph.edges[node], err = graphTopology.RequestSubset(node)
			require.NoError(t, err)
		}
		graph.targetDiameter = 2
		graph.distance = make([][]int, nodeCount)
		for i := 0; i < nodeCount; i++ {
			graph.distance[i] = make([]int, nodeCount)
			for j := 0; j < nodeCount; j++ {
				graph.distance[i][j] = graph.targetDiameter + 1
			}
		}

		pairsToUpdate := nodeCount * (nodeCount - 1) / 2 // we have C(n,2) unordered pairs of nodes
		updatedPairs := make(map[int]bool)
		for graphSize := nodeCount; graphSize > 0 && len(updatedPairs) < pairsToUpdate; graphSize-- {
			source := graphSize - 1
			graph.bfs(source, source, graph.distance[source])

			// Now we need to determine shortest path distance for any pair of nodes `(u,v)` such that the path
			// between `u` and `v` includes `source`. Note that we don't need to determine the shortest path distance
			// for all pairs of nodes. We only determine the shortest path distance for some pair `(u,v)` such that
			// `distance[u][v] <= targetDiameter`, otherwise the test will fail. Which gives us opportunity to optimize here.
			for peer := 0; peer < source; peer++ {
				d := graph.distance[source][peer]
				require.True(t, d <= graph.targetDiameter, "graph diameter more than expected")
				// assuming that the graph is bidirectional, which is tested in TestGraphIsBidirectional
				graph.distance[peer][source] = d
				updatedPairs[combinedIndex(peer, source, nodeCount)] = true
			}
			// update any pair `(nodeA,nodeB)` such that the shortest path between `nodeA` and `nodeB` includes `source`
			// As we have `targetedDiameter = 2`, doing this operation is not very costly.
			for nodeAIndex, nodeA := range graph.edges[source] {
				for nodeBIndex := nodeAIndex + 1; nodeBIndex < len(graph.edges[source]); nodeBIndex++ {
					nodeB := graph.edges[source][nodeBIndex]
					if graph.distance[nodeA][nodeB] > 2 {
						graph.distance[nodeA][nodeB] = 2
						graph.distance[nodeB][nodeA] = 2
						updatedPairs[combinedIndex(nodeA, nodeB, nodeCount)] = true
					}
				}
			}
		}
	}

	testForMultipleGraph(MinGraphSize, MaxGraphSize, tester)
}

// tests diameter from `root=0` only
func TestOptimizedGraphDiameter(t *testing.T) {
	tester := func(nodeCount int) {
		graph := newGraph(nodeCount)
		network, err := NewOptimizedTopology(nodeCount, 0, 0)
		require.NoError(t, err)

		graph.targetDiameter = 2
		graph.distance = make([][]int, nodeCount)
		graph.distance[0] = make([]int, nodeCount)
		for node := 0; node < nodeCount; node++ {
			graph.edges[node], err = network.RequestSubset(node)
			require.NoError(t, err)
			graph.distance[0][node] = graph.targetDiameter + 1
		}

		graph.bfs(0, nodeCount, graph.distance[0])
		for node := 0; node < nodeCount; node++ {
			require.True(t, graph.distance[0][node] <= graph.targetDiameter)
		}
	}

	testForMultipleGraph(MinGraphSize, MaxGraphSize, tester)
}

func TestNoDuplicateNodeInConnection(t *testing.T) {
	tester := func(nodeCount int) {
		tester2 := func(degreeCount int) {
			duplicateChecker := func(graph Topology) {
				for node := 0; node < nodeCount; node++ {
					edges, err := graph.RequestSubset(node)
					require.NoError(t, err)
					existedPeer := make(map[int]bool)
					existedPeer[node] = true
					for _, peer := range edges {
						require.False(t, existedPeer[peer])
						existedPeer[peer] = true
					}
				}
			}

			graph, err := NewFaultTolerantTopology(nodeCount, degreeCount, 0)
			require.NoError(t, err)
			duplicateChecker(graph)

			graph, err = NewOptimizedTopology(nodeCount, degreeCount, 0)
			require.NoError(t, err)
			duplicateChecker(graph)
		}

		tester2(0)
		if nodeCount > 1 {
			tester2((1 + nodeCount) / 2)
		}
		if nodeCount > 2 {
			tester2((2*nodeCount + 2) / 3)
		}
	}
	testForMultipleGraph(MinGraphSize, MaxGraphSize, tester)
}

func TestGraphConnectivityForFaultTolerant(t *testing.T) {
	tester := func(nodeCount int) {
		graphTopology, err := NewFaultTolerantTopology(nodeCount, 0, 0)
		require.NoError(t, err)
		connectivity, err := graphTopology.(*faultTolerantTopology).connectivity()
		require.NoError(t, err)
		require.True(t, connectivity > int(bft.F(big.NewInt(int64(nodeCount))).Int64()))
	}
	testForMultipleGraph(2, 100, tester)
}

func TestGraphConnectivityForLargeFaultTolerant(t *testing.T) {
	// TODO (tariq): exploit the propery of the graph that it is symmetric and has a unique shape
}

func BenchmarkEdgeConstructionForDegreeSqrt(b *testing.B) {
	graph, err := NewOptimizedTopology(1024, 0, 0)
	require.NoError(b, err)
	require.Equal(b, 32, graph.DegreeCount())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		graph.RequestSubset(0)
	}
}

func BenchmarkEdgeConstructionForDegree100(b *testing.B) {
	graph, err := NewOptimizedTopology(1024, 100, 0)
	require.NoError(b, err)
	require.Equal(b, 100, graph.DegreeCount())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		graph.RequestSubset(0)
	}
}

func BenchmarkEdgeConstructionForDegreeOneThird(b *testing.B) {
	graph, err := NewFaultTolerantTopology(1000, 0, 0)
	require.NoError(b, err)
	require.True(b, 999/3+1 <= graph.DegreeCount() && 999/3+2 >= graph.DegreeCount())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		graph.RequestSubset(0)
	}
}

func BenchmarkEdgeConstructionForDegreeHalf(b *testing.B) {
	graph, err := NewFaultTolerantTopology(1000, 500, 0)
	require.NoError(b, err)
	require.True(b, 500 <= graph.DegreeCount() && 501 >= graph.DegreeCount())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		graph.RequestSubset(0)
	}
}

func BenchmarkEdgeConstructionForDegreeTwoThird(b *testing.B) {
	graph, err := NewFaultTolerantTopology(1000, 2002/3, 0)
	require.NoError(b, err)
	require.True(b, 2002/3 <= graph.DegreeCount() && 2002/3+1 >= graph.DegreeCount())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		graph.RequestSubset(0)
	}
}

func testForMultipleGraph(
	smallestGraphSize, largestGraphSize int,
	tester func(ndoeCount int),
) {
	for size := smallestGraphSize; size <= largestGraphSize; size++ {
		tester(size)
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

type Graph struct {
	visited        []bool
	distance       [][]int
	edges          [][]int
	targetDiameter int
}

func newGraph(size int) Graph {
	graph := Graph{}
	graph.init(size)
	return graph
}

func (g *Graph) init(size int) {
	g.edges = make([][]int, size)
}

func (g *Graph) dfs(node int) int {
	if g.visited[node] {
		return 0
	}
	g.visited[node] = true
	visited := 1

	for _, peer := range g.edges[node] {
		visited += g.dfs(peer)
	}
	return visited
}

// Must sort all edges.
// The bfs is modified to determine shortest path distance from `source` to any `node < graphSize`
// considering only the subgraph of the first `graphSize` nodes. Here we are assuming that we know shortest
// path distance for any pair of nodes `(u,v)` such that the path between `u` and `v` includes some
// node `w >= graphSize`. In this case, it is enough to consider the subgraph including only the first `graphSize` nodes.
func (g *Graph) bfs(source, graphSize int, distance []int) {
	// enque source
	queue := make([]int, 0, graphSize)
	queue = append(queue, source)
	distance[source] = 0
	for len(queue) > 0 {
		// pop
		node := queue[0]
		queue = queue[1:]
		for _, peer := range g.edges[node] {
			if peer >= graphSize {
				break
			}
			if distance[peer] > distance[node]+1 {
				// enque adjacent nodes
				queue = append(queue, peer)
				distance[peer] = distance[node] + 1
			}
		}
	}
}
