package maxflow

import "math"

/**************************************************************
 * Dinitz' Max Flow w/ Shimon Even and Alon Itai optimization *
 * O(V^2E) worst. O(E\sqrt(V)) for unit cap.                  *
 **************************************************************/

const inf int = math.MaxInt

type Edge struct{ from, to, cap, flow, rev int }

type Graph struct {
	next, depth      []int
	edges            [][]Edge
	source, terminal int
}

func NewGraph(size int) Graph {
	g := Graph{}
	g.init(size)
	return g
}

func (g *Graph) init(n int) {
	g.next = make([]int, n)
	g.depth = make([]int, n)
	g.edges = make([][]Edge, n)
}

// Adds an edge in the graph between node `u` and `v`. For bidirectional edge, put `reverseCap = cap`.
// For unidirectional edge, put `reverseCap = 0`
func (g *Graph) AddEdge(u, v, cap, reverseCap int) {
	g.edges[u] = append(g.edges[u], Edge{u, v, cap, 0, len(g.edges[v])})
	g.edges[v] = append(g.edges[v], Edge{v, u, reverseCap, 0, len(g.edges[u]) - 1})
}

// Run bfs algorithm from g.source
// Returns true if a path exists from g.source to g.terminal
func (g *Graph) bfs() bool {
	for i := 0; i < len(g.depth); i++ {
		g.depth[i] = -1
	}
	g.depth[g.source] = 0
	queue := make([]int, 0, len(g.edges))
	// push
	queue = append(queue, g.source)

	for len(queue) > 0 {
		u := queue[0]
		// pop
		queue = queue[1:]
		for _, edge := range g.edges[u] {
			v := edge.to
			if edge.cap > 0 && g.depth[v] == -1 {
				g.depth[v] = g.depth[u] + 1
				if v == g.terminal {
					return true
				}
				// push
				queue = append(queue, v)
			}
		}
	}

	return g.depth[g.terminal] > -1
}

func (g *Graph) increaseFlow(u, index, flow int) {
	edge := &g.edges[u][index]
	edge.flow += flow
	edge.cap -= flow

	// decrease flow of reverse edge
	reverseEdge := &g.edges[edge.to][edge.rev]
	reverseEdge.cap += flow
	reverseEdge.flow -= flow
}

func (g *Graph) dfs(u, flowIn int) int {
	if u == g.terminal {
		return flowIn
	}

	for g.next[u] < len(g.edges[u]) {
		edge := g.edges[u][g.next[u]]
		if edge.cap > 0 && g.depth[edge.to] == g.depth[u]+1 {
			flow := g.dfs(edge.to, min(flowIn, edge.cap))
			if flow > 0 {
				g.increaseFlow(u, g.next[u], flow)
				return flow
			}
		}
		g.next[u]++
	}

	return 0
}

// Calculates maximum flow for a graph from `source` node to `sink` node.
// Reading materials:
//	1. https://en.wikipedia.org/wiki/Dinic%27s_algorithm (wiki)
//	2. https://cp-algorithms.com/graph/dinic.html (elaboration of the algo)
// 	3. http://dx.doi.org/10.1007/11685654_10 (paper)

func (g *Graph) DinitzMaxFlow(source, sink int) int {
	if source == sink {
		// should not happen
		return inf
	}
	g.source = source
	g.terminal = sink
	maxFlow := 0
	for g.bfs() {
		for i := 0; i < len(g.next); i++ {
			g.next[i] = 0
		}
		for {
			flow := g.dfs(g.source, inf)
			if flow == 0 {
				break
			}
			maxFlow += flow
		}
	}
	return maxFlow
}
