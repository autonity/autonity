package graph

import (
	"sort"
)

type view string

const (
	TB view = "TB" // top bottom
	BT      = "BT" // bottom top
	RL      = "RL" // right left
	LR      = "LR" // left right
	TD      = "TD" // same as TB
)

type Graph struct {
	graph       *graph
	names       []string
	initialized bool

	View      view        `"graph" @Ident`
	Edges     []*Edge     `@@*`
	SubGraphs []*SubGraph `@@*`
}

// GetEdges assumes that SetNodeName was called for all nodes in the graphs
func (gr *Graph) GetEdges(index int) []int {
	if !gr.initialized {
		gr.setEdges()
	}

	nodeName, ok := gr.graph.nodeIndexes[index]
	if !ok {
		return nil
	}

	edges, ok := gr.graph.edges[nodeName]
	if !ok {
		return nil
	}

	nodesNames := make([]int, len(edges))
	for i, name := range edges {
		nodesNames[i] = gr.graph.nodeNames[name]
	}

	return nodesNames
}

func (gr *Graph) GetNames() []string {
	if gr.names != nil {
		return gr.names
	}

	names := make(map[string]struct{})

	for _, edge := range gr.Edges {
		if edge == nil {
			continue
		}

		names[edge.LeftNode] = struct{}{}
		names[edge.RightNode] = struct{}{}
	}

	for _, sub := range gr.SubGraphs {
		for _, edge := range sub.Edges {
			if edge == nil {
				continue
			}

			names[edge.LeftNode] = struct{}{}
			names[edge.RightNode] = struct{}{}
		}
	}

	namesSlice := make([]string, 0, len(names))
	for name := range names {
		namesSlice = append(namesSlice, name)
	}

	sort.Strings(namesSlice)

	gr.names = namesSlice
	return namesSlice
}

func (gr *Graph) SetNodeName(name string, index int) {
	if gr.graph == nil {
		gr.graph = &graph{}
	}
	gr.graph.SetNodeName(name, index)
}

func (gr *Graph) setEdges() {
	if gr.graph == nil {
		gr.graph = &graph{}
	}

	for _, edge := range gr.Edges {
		if edge == nil {
			continue
		}

		gr.graph.setEdges(edge.LeftNode, edge.RightNode)
		if !edge.Directed {
			gr.graph.setEdges(edge.RightNode, edge.LeftNode)
		}
	}

	for _, sub := range gr.SubGraphs {
		for _, edge := range sub.Edges {
			if edge == nil {
				continue
			}

			gr.graph.setEdges(edge.LeftNode, edge.RightNode)
			if !edge.Directed {
				gr.graph.setEdges(edge.RightNode, edge.LeftNode)
			}
		}
	}

	gr.initialized = true
}

type SubGraph struct {
	Name  string  `"subgraph" @Ident`
	Edges []*Edge `@@*"end"`
}

type Edge struct {
	LeftNode  string `@Ident[" "|"\t"]"-""-"["-"]`
	Directed  bool   `[@">"][" "|"\t"]`
	RightNode string `@Ident[";"]`
}

type graph struct {
	edges       map[string][]string
	nodeIndexes map[int]string
	nodeNames   map[string]int
}

func (g *graph) SetNodeName(name string, index int) {
	if g.nodeIndexes == nil {
		g.nodeIndexes = make(map[int]string)
	}
	if g.nodeNames == nil {
		g.nodeNames = make(map[string]int)
	}
	g.nodeIndexes[index] = name
	g.nodeNames[name] = index
}

func (g *graph) setEdges(nodeA string, nodes ...string) {
	if g.edges == nil {
		g.edges = make(map[string][]string, len(nodes))
	}
	nodeEdges, _ := g.edges[nodeA]
	nodeEdges = append(nodeEdges, nodes...)

	g.edges[nodeA] = nodeEdges
}
