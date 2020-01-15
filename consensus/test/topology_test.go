package test

import (
	"fmt"
	"github.com/clearmatics/autonity/common/graph"
	"strings"
	"testing"
)

func TestTendermintStarSuccess(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	topologyStr := `graph TB
    VA---VB
    VA---VC
    VA---VD
    VA-->VE`
	topology, err := graph.Parse(strings.NewReader(topologyStr))
	if err != nil {
		t.Fatal("parse error")
	}
	cases := []*testCase{
		{
			name:          "no malicious",
			numValidators: 5,
			numBlocks:     5,
			txPerPeer:     1,
			topology: &Topology{
				graph: *topology,
			},
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
			runTest(t, testCase)
		})
	}
}

func TestTendermintStarOverParticipantSuccess(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	topologyStr := `graph TB
    PF---VA
    PF---VB
    PF---VC
    PF---VD
    PF-->VE`

	topology, err := graph.Parse(strings.NewReader(topologyStr))
	if err != nil {
		t.Fatal("parse error")
	}
	cases := []*testCase{
		{
			name:          "no malicious",
			numValidators: 5,
			numBlocks:     5,
			txPerPeer:     1,
			topology: &Topology{
				graph: *topology,
			},
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
			runTest(t, testCase)
		})
	}
}

func TestTendermintBusSuccess(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	topologyStr := `graph TB
    VA---VB
    VC---VB
    VD---VC
    VE---VD
`

	topology, err := graph.Parse(strings.NewReader(topologyStr))
	if err != nil {
		t.Fatal("parse error")
	}
	cases := []*testCase{
		{
			name:          "no malicious",
			numValidators: 5,
			numBlocks:     5,
			txPerPeer:     1,
			topology: &Topology{
				graph: *topology,
			},
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
			runTest(t, testCase)
		})
	}
}
