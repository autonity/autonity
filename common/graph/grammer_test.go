package graph

import (
	"bufio"
	"os"
	"reflect"
	"testing"

	"github.com/alecthomas/participle"
)

func TestGraphLexer(t *testing.T) {
	parser, err := participle.Build(&Graph{})
	if err != nil {
		t.Fatal(err)
	}

	graph := &Graph{}

	file, err := os.Open("graph.md")
	if err != nil {
		t.Fatal(err)
	}

	err = parser.Parse(bufio.NewReader(file), graph, participle.AllowTrailing(true))
	if err != nil {
		t.Fatal(err)
	}

	expected := &Graph{
		View: "TB",
		Edges: []*Edge{
			{
				LeftNode:  "B",
				Directed:  false,
				RightNode: "A",
			},
			{
				LeftNode:  "C",
				Directed:  false,
				RightNode: "A",
			},
			{
				LeftNode:  "A",
				Directed:  false,
				RightNode: "D",
			},
			{
				LeftNode:  "A",
				Directed:  true,
				RightNode: "E",
			},
		},
	}

	if !reflect.DeepEqual(expected, graph) {
		t.Errorf("got %v\n\nexpected %v", graph, expected)
	}
}

func TestGraphLexerGetNames(t *testing.T) {
	parser, err := participle.Build(&Graph{})
	if err != nil {
		t.Fatal(err)
	}

	graph := &Graph{}

	file, err := os.Open("graph.md")
	if err != nil {
		t.Fatal(err)
	}

	err = parser.Parse(bufio.NewReader(file), graph, participle.AllowTrailing(true))
	if err != nil {
		t.Fatal(err)
	}

	nodeNames := graph.GetNames()

	expected := []string{"A", "B", "C", "D", "E"}

	if !reflect.DeepEqual(nodeNames, expected) {
		t.Errorf("got %v\n\nexpected %v", nodeNames, expected)
	}
}

func TestGraphLexerGetEdges(t *testing.T) {
	parser, err := participle.Build(&Graph{})
	if err != nil {
		t.Fatal(err)
	}

	graph := &Graph{}

	file, err := os.Open("graph.md")
	if err != nil {
		t.Fatal(err)
	}

	err = parser.Parse(bufio.NewReader(file), graph, participle.AllowTrailing(true))
	if err != nil {
		t.Fatal(err)
	}

	nodeNames := graph.GetNames()
	for i, name := range nodeNames {
		graph.SetNodeName(name, i)
	}

	expected := [][]int{
		0: {1, 2, 3, 4},
		1: {0},
		2: {0},
		3: {0},
		4: nil,
	}

	for i := 0; i < len(nodeNames); i++ {
		edges := graph.GetEdges(i)
		if !reflect.DeepEqual(expected[i], edges) {
			t.Errorf("for node %q(%d)\ngot %v\n\nexpected %v", nodeNames[i], i, edges, expected[i])
		}
	}
}

func TestSubGraphLexer(t *testing.T) {
	parser, err := participle.Build(&Graph{})
	if err != nil {
		t.Fatal(err)
	}

	graph := &Graph{}

	file, err := os.Open("subgraph.md")
	if err != nil {
		t.Fatal(err)
	}

	err = parser.Parse(bufio.NewReader(file), graph, participle.AllowTrailing(true))
	if err != nil {
		t.Fatal(err)
	}

	expected := &Graph{
		View: "TB",
		Edges: []*Edge{
			{
				LeftNode:  "B",
				Directed:  false,
				RightNode: "A",
			},
			{
				LeftNode:  "C",
				Directed:  false,
				RightNode: "A",
			},
			{
				LeftNode:  "A",
				Directed:  false,
				RightNode: "D",
			},
			{
				LeftNode:  "A",
				Directed:  true,
				RightNode: "E",
			},
		},

		SubGraphs: []*SubGraph{
			{
				Name: "one",
				Edges: []*Edge{
					{
						LeftNode:  "a1",
						Directed:  true,
						RightNode: "a2",
					},
				},
			},
			{
				Name: "two",
				Edges: []*Edge{
					{
						LeftNode:  "b1",
						Directed:  true,
						RightNode: "b2",
					},
					{
						LeftNode:  "b2",
						Directed:  false,
						RightNode: "b3",
					},
				},
			},
			{
				Name: "three",
				Edges: []*Edge{
					{
						LeftNode:  "c1",
						Directed:  false,
						RightNode: "c2",
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(expected, graph) {
		t.Errorf("got %v\n\nexpected %v", graph, expected)
	}
}

func TestSubGraphLexerGetNames(t *testing.T) {
	parser, err := participle.Build(&Graph{})
	if err != nil {
		t.Fatal(err)
	}

	graph := &Graph{}

	file, err := os.Open("subgraph.md")
	if err != nil {
		t.Fatal(err)
	}

	err = parser.Parse(bufio.NewReader(file), graph, participle.AllowTrailing(true))
	if err != nil {
		t.Fatal(err)
	}

	nodeNames := graph.GetNames()

	expected := []string{"A", "B", "C", "D", "E", "a1", "a2", "b1", "b2", "b3", "c1", "c2"}

	if !reflect.DeepEqual(nodeNames, expected) {
		t.Errorf("got %v\n\nexpected %v", nodeNames, expected)
	}
}

func TestSubGraphLexerGetEdges(t *testing.T) {
	parser, err := participle.Build(&Graph{})
	if err != nil {
		t.Fatal(err)
	}

	graph := &Graph{}

	file, err := os.Open("subgraph.md")
	if err != nil {
		t.Fatal(err)
	}

	err = parser.Parse(bufio.NewReader(file), graph, participle.AllowTrailing(true))
	if err != nil {
		t.Fatal(err)
	}

	nodeNames := graph.GetNames()
	for i, name := range nodeNames {
		graph.SetNodeName(name, i)
	}

	expected := [][]int{
		0:  {1, 2, 3, 4},
		1:  {0},
		2:  {0},
		3:  {0},
		4:  nil,
		5:  {6},
		6:  nil,
		7:  {8},
		8:  {9},
		9:  {8},
		10: {11},
		11: {10},
	}

	for i := 0; i < len(expected); i++ {
		edges := graph.GetEdges(i)
		if !reflect.DeepEqual(expected[i], edges) {
			t.Errorf("for node %q(%d)\ngot %v\n\nexpected %v", nodeNames[i], i, edges, expected[i])
		}
	}
}
