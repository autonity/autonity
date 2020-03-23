package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/clearmatics/autonity/common/graph"
)

func TestTendermintExternalUser(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	topologyStr := `graph TB
    VA---VB
    VA---VC
    VA---VD
    VA---VE
	VA---EE
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

func TestTendermintStarOverParticipantSuccessWithExternalUser(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	topologyStr := `graph TB
    PF---VA
    PF---VB
    PF---VC
    PF---VD
    PF-->VE
	PF---EF
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

func TestTendermintBusSuccessWithExternalUser(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	topologyStr := `graph TB
    VA---VB
    VC---VB
    VD---VC
    VE---VD
    VE---EF
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

func TestTendermintChangeTopologyFromBusToStarSuccessWithExternalUser(t *testing.T) {
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
		VE---EF
    end
    subgraph b7
		VA---VB
		VA---VC
		VA---VD
		VA-->VE
		VA---EF
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

func TestTendermintChangeTopologyFromStarToBusSuccessWithExternalUser(t *testing.T) {
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
		VA---EF
	end
    subgraph b7
		VA---VB
		VC---VB
		VD---VC
		VE---VD
		VE---EF
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

func TestTendermintAddConnectionToTopologySuccessWithExternalUser(t *testing.T) {
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
		VE---EF
    end
    subgraph b20
		VA---VB
		VA---VC
		VC---VB
		VD---VC
		VE---VD
		VE---EF
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

func TestTendermintAddValidatorsToTopologySuccessWithExternalUser(t *testing.T) {
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
		VF---EH
    end
    subgraph b20
		VA---VB
		VA---VF
		VC---VB
		VD---VC
		VE---VD
		VF---VG
		VF---EH
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

func TestTendermintAddParticipantsToTopologySuccessWithExternalUser(t *testing.T) {
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
		PF---EH
    end
    subgraph b20
		VA---VB
		VA---PF
		VC---VB
		VD---VC
		VE---VD
		PF---PG
		PF---EH
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

func TestTendermintAddStakeholdersToTopologySuccessWithExternalUser(t *testing.T) {
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
		SF---EH
    end
    subgraph b20
		VA---VB
		VA---SF
		VC---VB
		VD---VC
		VE---VD
		SF---SG
		SF---EH
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
