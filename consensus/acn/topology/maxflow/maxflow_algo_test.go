package maxflow

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMaxFlowAlgoOnBidirectionalGraph(t *testing.T) {
	// make a graph
	/*
		Edges:
		0 - 1 with cap = 10
		1 - 2 with cap = 1
		0 - 2 with cap = 10
		2 - 4 with cap = 2
		2 - 3 with cap = 2
		0 - 3 with cap = 10
		3 - 4 with cap = 5
		4 - 5 with cap = 6
		5 - 6 with cap = 10
		6 - 7 with cap = 4
		5 - 7 with cap = 1

		flow from 0 to 7 should be 5
	*/

	edges := make([]Edge, 0)
	edges = append(edges, Edge{0, 1, 10, 0, 10})
	edges = append(edges, Edge{2, 1, 1, 0, 1})
	edges = append(edges, Edge{0, 2, 10, 0, 10})
	edges = append(edges, Edge{2, 4, 2, 0, 2})
	edges = append(edges, Edge{2, 3, 2, 0, 2})
	edges = append(edges, Edge{0, 3, 10, 0, 10})
	edges = append(edges, Edge{3, 4, 5, 0, 5})
	edges = append(edges, Edge{4, 5, 6, 0, 6})
	edges = append(edges, Edge{5, 6, 10, 0, 10})
	edges = append(edges, Edge{6, 7, 4, 0, 4})
	edges = append(edges, Edge{5, 7, 1, 0, 1})
	nodes := 8

	graph := NewGraph(nodes)
	for _, edge := range edges {
		graph.AddEdge(edge.from, edge.to, edge.cap, edge.cap)
	}

	flow := graph.DinitzMaxFlow(0, nodes-1)
	require.Equal(t, 5, flow)
}

func TestMaxFlowAlgoOnBidirectionalGraph1(t *testing.T) {
	// make a graph
	/*
		Edges:
		0 - 1 with cap = 10
		1 - 2 with cap = 1
		0 - 2 with cap = 2
		2 - 4 with cap = 5
		2 - 3 with cap = 6
		0 - 3 with cap = 10
		3 - 4 with cap = 5
		4 - 5 with cap = 100
		5 - 6 with cap = 10
		6 - 7 with cap = 11
		5 - 7 with cap = 1

		flow from 0 to 7 should be 10
	*/

	edges := make([]Edge, 0)
	edges = append(edges, Edge{0, 1, 10, 0, 10})
	edges = append(edges, Edge{2, 1, 1, 0, 1})
	edges = append(edges, Edge{0, 2, 2, 0, 2})
	edges = append(edges, Edge{2, 4, 5, 0, 5})
	edges = append(edges, Edge{2, 3, 6, 0, 6})
	edges = append(edges, Edge{0, 3, 10, 0, 10})
	edges = append(edges, Edge{3, 4, 5, 0, 5})
	edges = append(edges, Edge{4, 5, 100, 0, 100})
	edges = append(edges, Edge{5, 6, 10, 0, 10})
	edges = append(edges, Edge{6, 7, 11, 0, 11})
	edges = append(edges, Edge{5, 7, 1, 0, 1})
	nodes := 8

	graph := NewGraph(nodes)
	for _, edge := range edges {
		graph.AddEdge(edge.from, edge.to, edge.cap, edge.cap)
	}

	flow := graph.DinitzMaxFlow(0, nodes-1)
	require.Equal(t, 10, flow)
}

func TestMaxFlowAlgoOnBidirectionalGraph2(t *testing.T) {
	// make a graph
	/*
		Edges:
		0 - 1 with cap = 10
		1 - 2 with cap = 1
		0 - 2 with cap = 2
		2 - 4 with cap = 5
		2 - 3 with cap = 6
		0 - 3 with cap = 10
		3 - 4 with cap = 8
		4 - 5 with cap = 100
		5 - 6 with cap = 13
		6 - 7 with cap = 12
		5 - 7 with cap = 3

		flow from 0 to 7 should be 13
	*/

	edges := make([]Edge, 0)
	edges = append(edges, Edge{0, 1, 10, 0, 10})
	edges = append(edges, Edge{2, 1, 1, 0, 1})
	edges = append(edges, Edge{0, 2, 2, 0, 2})
	edges = append(edges, Edge{2, 4, 5, 0, 5})
	edges = append(edges, Edge{2, 3, 6, 0, 6})
	edges = append(edges, Edge{0, 3, 10, 0, 10})
	edges = append(edges, Edge{3, 4, 8, 0, 8})
	edges = append(edges, Edge{4, 5, 100, 0, 100})
	edges = append(edges, Edge{5, 6, 13, 0, 13})
	edges = append(edges, Edge{6, 7, 12, 0, 12})
	edges = append(edges, Edge{5, 7, 3, 0, 3})
	nodes := 8

	graph := NewGraph(nodes)
	for _, edge := range edges {
		graph.AddEdge(edge.from, edge.to, edge.cap, edge.cap)
	}

	flow := graph.DinitzMaxFlow(0, nodes-1)
	require.Equal(t, 13, flow)
}

func TestMaxFlowAlgoOnUnidirectionalGraph(t *testing.T) {
	// make a graph
	/*
		Edges:
		0 -> 1 with cap = 10
		1 -> 2 with cap = 1
		0 -> 2 with cap = 10
		2 -> 4 with cap = 2
		2 -> 3 with cap = 2
		0 -> 3 with cap = 10
		3 -> 4 with cap = 5
		4 -> 5 with cap = 6
		5 -> 6 with cap = 10
		7 -> 6 with cap = 4
		5 -> 7 with cap = 1

		flow from 0 to 7 should be 1
	*/

	edges := make([]Edge, 0)
	edges = append(edges, Edge{0, 1, 10, 0, 0})
	edges = append(edges, Edge{1, 2, 1, 0, 0})
	edges = append(edges, Edge{0, 2, 10, 0, 0})
	edges = append(edges, Edge{2, 4, 2, 0, 0})
	edges = append(edges, Edge{2, 3, 2, 0, 0})
	edges = append(edges, Edge{0, 3, 10, 0, 0})
	edges = append(edges, Edge{3, 4, 5, 0, 0})
	edges = append(edges, Edge{4, 5, 6, 0, 0})
	edges = append(edges, Edge{5, 6, 10, 0, 0})
	edges = append(edges, Edge{7, 6, 4, 0, 0})
	edges = append(edges, Edge{5, 7, 1, 0, 0})
	nodes := 8

	graph := NewGraph(nodes)
	for _, edge := range edges {
		graph.AddEdge(edge.from, edge.to, edge.cap, 0)
	}

	flow := graph.DinitzMaxFlow(0, nodes-1)
	require.Equal(t, 1, flow)
}

func TestMaxFlowAlgoOnUnidirectionalGraph1(t *testing.T) {
	// make a graph
	/*
		Edges:
		0 -> 1 with cap = 10
		1 -> 2 with cap = 1
		0 -> 2 with cap = 10
		2 -> 4 with cap = 20
		2 -> 3 with cap = 2
		0 -> 3 with cap = 10
		3 -> 4 with cap = 5
		4 -> 5 with cap = 100
		5 -> 6 with cap = 10
		6 -> 7 with cap = 8
		5 -> 7 with cap = 10

		flow from 0 to 7 should be 16
	*/

	edges := make([]Edge, 0)
	edges = append(edges, Edge{0, 1, 10, 0, 0})
	edges = append(edges, Edge{1, 2, 1, 0, 0})
	edges = append(edges, Edge{0, 2, 10, 0, 0})
	edges = append(edges, Edge{2, 4, 20, 0, 0})
	edges = append(edges, Edge{2, 3, 2, 0, 0})
	edges = append(edges, Edge{0, 3, 10, 0, 0})
	edges = append(edges, Edge{3, 4, 5, 0, 0})
	edges = append(edges, Edge{4, 5, 100, 0, 0})
	edges = append(edges, Edge{5, 6, 10, 0, 0})
	edges = append(edges, Edge{6, 7, 8, 0, 0})
	edges = append(edges, Edge{5, 7, 10, 0, 0})
	nodes := 8

	graph := NewGraph(nodes)
	for _, edge := range edges {
		graph.AddEdge(edge.from, edge.to, edge.cap, 0)
	}

	flow := graph.DinitzMaxFlow(0, nodes-1)
	require.Equal(t, 16, flow)
}

func TestMaxFlow_SameSourceSink(t *testing.T) {
	graph := NewGraph(1)
	require.Equal(t, inf, graph.DinitzMaxFlow(0, 0))
}
