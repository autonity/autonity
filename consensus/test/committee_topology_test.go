package test

import (
	"strings"
	"testing"

	"github.com/autonity/autonity/common/graph"
)

// Committee core network should keep liveness on different topology to make sure the gossiping works.
func TestStarTopology(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	topologyStr := `graph TB
    V0---V1
    V0---V2
    V0---V3
    V0-->V4`
	topology, err := graph.Parse(strings.NewReader(topologyStr))
	if err != nil {
		t.Fatal("parse error")
	}
	testCase := &testCase{
		name:          "test committee star topology",
		numValidators: 5,
		numBlocks:     5,
		topology: &Topology{
			graph: *topology,
		},
	}

	runTest(t, testCase)
}

func TestBusTopology(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	topologyStr := `graph TB
    V0---V1
    V2---V1
    V3---V2
    V4---V3
`
	topology, err := graph.Parse(strings.NewReader(topologyStr))
	if err != nil {
		t.Fatal("parse error")
	}
	testCase := &testCase{
		name:          "test committee bus topology",
		numValidators: 5,
		numBlocks:     5,
		topology: &Topology{
			graph: *topology,
		},
	}

	runTest(t, testCase)
}

func TestRingTopology(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	topologyStr := `graph TB
    V0---V1
    V2---V1
    V3---V2
    V4---V3
	V0---V4
`
	topology, err := graph.Parse(strings.NewReader(topologyStr))
	if err != nil {
		t.Fatal("parse error")
	}
	testCase := &testCase{
		name:          "test committee ring topology",
		numValidators: 5,
		numBlocks:     5,
		topology: &Topology{
			graph: *topology,
		},
	}

	runTest(t, testCase)
}
