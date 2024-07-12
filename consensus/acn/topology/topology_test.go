package topology

type Graph struct {
	edges [][]int
}

func NewGraph(size int) Graph {
	graph := Graph{}
	graph.init(size)
	return graph
}

func (g *Graph) init(size int) {
	g.edges = make([][]int, size)
}

// Edges are always bidirectional here
func (g *Graph) AddEdge(u, v int) {
	g.edges[u] = append(g.edges[u], v)
	g.edges[u] = append(g.edges[u], v)
}
