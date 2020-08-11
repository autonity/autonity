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

func TestTendermintMajorityExternalUsers(t *testing.T) {
	t.Skip("This test is intermittently failing, see - https://github.com/clearmatics/autonity/issues/619")
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	topologyStr := `graph TB
    VA---VB
    VA---VC
    VA---VD
    VA---VE
	VA---EE
	VA---EF
	VA---EG
	VA---EH
	VA---EI
	VA---EJ
	VA---EK
	VA---EL
	VA---EM
	VA---EN
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

func TestTendermintMajorityExternalUsersFullyConnected(t *testing.T) {
	t.Skip("This test is intermittently failing, see - https://github.com/clearmatics/autonity/issues/619")
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	topologyStr := `graph TB
    VA---VB
    VA---VC
    VA---VD
    VA---VE
	VA---EE
	VA---EF
	VA---EG
	VA---EH
	VA---EI
	VA---EJ
	VA---EK
	VA---EL
	VA---EM
	VA---EN
    VB---VC
    VB---VD
    VB---VE
	VB---EE
	VB---EF
	VB---EG
	VB---EH
	VB---EI
	VB---EJ
	VB---EK
	VB---EL
	VB---EM
	VB---EN
    VC---VD
    VC---VE
	VC---EE
	VC---EF
	VC---EG
	VC---EH
	VC---EI
	VC---EJ
	VC---EK
	VC---EL
	VC---EM
	VC---EN
    VD---VE
	VD---EE
	VD---EF
	VD---EG
	VD---EH
	VD---EI
	VD---EJ
	VD---EK
	VD---EL
	VD---EM
	VD---EN
	VE---EF
	VE---EG
	VE---EH
	VE---EI
	VE---EJ
	VE---EK
	VE---EL
	VE---EM
	VE---EN
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
	t.Skip("test is flaky - https://github.com/clearmatics/autonity/issues/496")
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
