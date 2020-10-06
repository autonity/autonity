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
	//t.Skip("This test is intermittently failing, see - https://github.com/clearmatics/autonity/issues/619")
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
	//t.Skip("This test is intermittently failing, see - https://github.com/clearmatics/autonity/issues/619")
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
