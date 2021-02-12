package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/clearmatics/autonity/common/graph"
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
	t.Skip("test is flaky - https://github.com/clearmatics/autonity/issues/496")
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

func TestTendermintChangeTopologyFromBusToStarSuccess(t *testing.T) {
	t.Skip("Topology tests are not stable")

	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	topologyStr := `graph TB
    subgraph b1
		VA---VB
		VC---VB
		VD---VC
		VE---VD
    end
    subgraph b7
		VA---VB
		VA---VC
		VA---VD
		VA-->VE
	end
`

	topology, err := graph.Parse(strings.NewReader(topologyStr))
	if err != nil {
		t.Fatal("parse error")
	}

	cases := []*testCase{
		{
			name:          "no malicious",
			numValidators: 5,
			numBlocks:     20,
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

func TestTendermintChangeTopologyFromStarToBusSuccess(t *testing.T) {
	t.Skip("Topology tests are not stable")

	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	topologyStr := `graph TB
    subgraph b1
		VA---VB
		VA---VC
		VA---VD
		VA-->VE
	end
    subgraph b7
		VA---VB
		VC---VB
		VD---VC
		VE---VD
    end

`

	topology, err := graph.Parse(strings.NewReader(topologyStr))
	if err != nil {
		t.Fatal("parse error")
	}

	cases := []*testCase{
		{
			name:          "no malicious",
			numValidators: 5,
			numBlocks:     20,
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

func TestTendermintAddConnectionToTopologySuccess(t *testing.T) {
	t.Skip("Topology tests are not stable")

	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	topologyStr := `graph TB
    subgraph b7
		VA---VB
		VC---VB
		VD---VC
		VE---VD
    end
    subgraph b20
		VA---VB
		VA---VC
		VC---VB
		VD---VC
		VE---VD
	end
`

	topology, err := graph.Parse(strings.NewReader(topologyStr))
	if err != nil {
		t.Fatal("parse error")
	}

	cases := []*testCase{
		{
			name:          "no malicious",
			numValidators: 5,
			numBlocks:     30,
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

func TestTendermintAddValidatorsToTopologySuccess(t *testing.T) {
	t.Skip("Topology tests are not stable")

	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	topologyStr := `graph TB
    subgraph b7
		VA---VB
		VC---VB
		VD---VC
		VE---VD
		VF---VG
    end
    subgraph b20
		VA---VB
		VA---VF
		VC---VB
		VD---VC
		VE---VD
		VF---VG
	end
`

	topology, err := graph.Parse(strings.NewReader(topologyStr))
	if err != nil {
		t.Fatal("parse error")
	}

	cases := []*testCase{
		{
			name:          "no malicious",
			numValidators: 5,
			numBlocks:     30,
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

func TestTendermintAddParticipantsToTopologySuccess(t *testing.T) {
	t.Skip("should be fixed by https://github.com/clearmatics/autonity/issues/431")

	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	topologyStr := `graph TB
    subgraph b7
		VA---VB
		VC---VB
		VD---VC
		VE---VD
		PF---PG
    end
    subgraph b20
		VA---VB
		VA---PF
		VC---VB
		VD---VC
		VE---VD
		PF---PG
	end
`

	topology, err := graph.Parse(strings.NewReader(topologyStr))
	if err != nil {
		t.Fatal("parse error")
	}

	cases := []*testCase{
		{
			name:          "no malicious",
			numValidators: 5,
			numBlocks:     30,
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

func TestTendermintAddStakeholdersToTopologySuccess(t *testing.T) {
	t.Skip("should be fixed by https://github.com/clearmatics/autonity/issues/431")

	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	topologyStr := `graph TB
    subgraph b7
		VA---VB
		VC---VB
		VD---VC
		VE---VD
		SF---SG
    end
    subgraph b20
		VA---VB
		VA---SF
		VC---VB
		VD---VC
		VE---VD
		SF---SG
	end
`

	topology, err := graph.Parse(strings.NewReader(topologyStr))
	if err != nil {
		t.Fatal("parse error")
	}

	cases := []*testCase{
		{
			name:          "no malicious",
			numValidators: 5,
			numBlocks:     30,
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
