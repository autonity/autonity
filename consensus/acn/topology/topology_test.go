package topology

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

const MinNodes = 2

func TestDegreeIsMoreThanOneThird(t *testing.T) {
	tester := func(nodeCount int) {
		graph, err := NewGraphTopology(nodeCount, 0, MinNodes)
		require.NoError(t, err)
		require.True(t, graph.isNetworkValid())
		// smallest possible degree count
		require.True(t, graph.degreeCount > (nodeCount-1)/3)
		require.True(t, graph.degreeCount-2 <= (nodeCount-1)/3)

		graph.SetDegreeCount(graph.degreeCount - 2)

		require.True(t, graph.isNetworkValid())
		// smallest possible degree count
		require.True(t, graph.degreeCount > (nodeCount-1)/3)
		require.True(t, graph.degreeCount-2 <= (nodeCount-1)/3)
	}
	testForMultipleGraph(2, 1000, tester)
}

func TestDegreeCountNotMoreThanNetworkSize(t *testing.T) {
	_, err := NewGraphTopology(10, 10, 1)
	require.Error(t, err)
	require.Equal(t, errDegreeCount, err)

	graph, err := NewGraphTopology(10, 9, 1)
	require.NoError(t, err)
	err = graph.SetDegreeCount(10)
	require.Error(t, err)
	require.Equal(t, errDegreeCount, err)

	require.Equal(t, 9, graph.degreeCount)
}

func TestAdjacentNodesAreConnected(t *testing.T) {
	tester := func(nodeCount int) {
		graph, err := NewGraphTopology(nodeCount, 0, 0)
		require.NoError(t, err)
		for node := 0; node < nodeCount; node++ {
			edges := graph.RequestSubset(node)
			for _, peer := range edges {
				require.True(t, graph.isAdjacent(node, peer))
			}
		}
	}

	testForMultipleGraph(2, 1000, tester)
}

func TestGraphIsBidirectional(t *testing.T) {
	tester := func(nodeCount int) {
		graph, err := NewGraphTopology(nodeCount, 0, 0)
		require.NoError(t, err)
		edgesMap := make(map[int]bool)
		for node := 0; node < nodeCount; node++ {
			edges := graph.RequestSubset(node)
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

	testForMultipleGraph(2, 1000, tester)
}

func TestDegreeCount(t *testing.T) {
	tester := func(nodeCount int) {
		graph, err := NewGraphTopology(nodeCount, 0, 0)
		require.NoError(t, err)
		for node := 0; node < nodeCount; node++ {
			edges := graph.RequestSubset(node)
			require.Equal(t, graph.degreeCount, len(edges))
		}
	}

	testForMultipleGraph(2, 1000, tester)
}

func TestGraphIsConnected(t *testing.T) {
	tester := func(nodeCount int) {
		graphTopology, err := NewGraphTopology(nodeCount, 0, 0)
		require.NoError(t, err)
		graph := newGraph(nodeCount)
		for node := 0; node < nodeCount; node++ {
			graph.edges[node] = graphTopology.RequestSubset(node)
		}
		require.Equal(t, nodeCount, graph.dfs(0))
	}

	testForMultipleGraph(1, 1000, tester)
}

func TestGraphDiameter(t *testing.T) {
	// TODO (tariq): complete
}

func TestGraphConnectivity(t *testing.T) {
	// TODO (tariq): need to optimize
	tester := func(nodeCount int) {
		fmt.Printf("node %v\n", nodeCount)
		graphTopology, err := NewGraphTopology(nodeCount, 0, 0)
		require.NoError(t, err)
		require.Equal(t, graphTopology.degreeCount, graphTopology.connectivity())
	}
	// TODO (tariq): increase upto 1000
	testForMultipleGraph(1, 100, tester)
}

func BenchmarkEdgeConstructionForDegreeOneThird(b *testing.B) {
	graph, err := NewGraphTopology(1000, 0, 0)
	require.NoError(b, err)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		graph.RequestSubset(0)
	}
}

func BenchmarkEdgeConstructionForDegreeHalf(b *testing.B) {
	graph, err := NewGraphTopology(1000, 500, 0)
	require.NoError(b, err)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		graph.RequestSubset(0)
	}
}

func BenchmarkEdgeConstructionForDegreeTwoThird(b *testing.B) {
	graph, err := NewGraphTopology(1000, 2002/3, 0)
	require.NoError(b, err)
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
	g.distance = make([][]int, size)
	g.visited = make([]bool, size)

	for i := 0; i < size; i++ {
		g.distance[i] = make([]int, size)
		for j := 0; j < size; j++ {
			g.distance[i][j] = g.targetDiameter + 1
		}
	}
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
