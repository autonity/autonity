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
