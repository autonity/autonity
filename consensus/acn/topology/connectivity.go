package topology

import (
	"github.com/autonity/autonity/consensus/acn/topology/maxflow"
)

// Calculates the connectivity of the graph using Menger's theorem.
// For more details, read: https://en.m.wikipedia.org/wiki/K-vertex-connected_graph (section: Computational complexity)
func (g *networkTopology) connectivity() int {
	// The connectivity of the graph is minimum of `P(u,v)` for non-adjacent pair of nodes `(u,v)`,
	// where `P(u,v)` = maximum vertex independent path count between non-adjacent pair of nodes `(u,v)`.
	res := g.nodeCount - 1
	for node := 0; node < g.nodeCount; node++ {
		for peer := node + 1; peer < g.nodeCount; peer++ {
			if !g.isAdjacent(node, peer) {
				res = min(res, g.maxVertexIndependentPathCount(node, peer))
			}
		}
	}
	return res
}

// Calculates the maximum count of vertex indepent path for some non-adjacent pair of nodes `(u,v)`.
func (g *networkTopology) maxVertexIndependentPathCount(u, v int) int {
	// In the graph to calculate max flow, each node `u` will be split into two separate nodes: `u_in` and `u_out`.
	// All the ingoing edges of `u` enter `u_in` where all the outgoing edges of `u` exit from `u_out` with capacity 1.
	// There will be an edge from `u_in` to `u_out` with capacity 1. For some pair of nodes `(u,v)`, max flow from
	// `u_out` to `v_in` will give us the maximum vertex indepent path count between `(u,v)` denoted as `P(u,v)`.

	// for some node `u`, we will set `u_in = u` and `u_out = u + g.nodeCount`
	graph := maxflow.NewGraph(g.nodeCount * 2)

	for node := 0; node < g.nodeCount; node++ {
		nodeOut := node + g.nodeCount
		graph.AddEdge(node, nodeOut, 1, 0)

		edges := g.RequestSubset(node)
		for _, peer := range edges {
			graph.AddEdge(nodeOut, peer, 1, 0)
		}
	}

	return graph.DinitzMaxFlow(u+g.nodeCount, v)
}
