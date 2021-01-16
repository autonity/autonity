package test

import (
	"fmt"
	"testing"
	"time"
)

func TestTendermintStopUpToFNodes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cases := []*testCase{
		{
			name:          "one node stops at block 1",
			numValidators: 5,
			numBlocks:     100,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VE": hookStopNode("VE", 1),
			},
			stopTime: make(map[string]time.Time),
			maliciousPeers: map[string]injectors{
				"VE": {},
			},
		},
		{
			name:          "one node stops at block 5",
			numValidators: 5,
			numBlocks:     100,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VE": hookStopNode("VE", 5),
			},
			stopTime: make(map[string]time.Time),
			maliciousPeers: map[string]injectors{
				"VE": {},
			},
		},
		{
			name:          "F nodes stop at block 1",
			numValidators: 7,
			numBlocks:     100,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VD": hookStopNode("VD", 1),
				"VE": hookStopNode("VE", 1),
			},
			stopTime: make(map[string]time.Time),
			maliciousPeers: map[string]injectors{
				"VD": {},
				"VE": {},
			},
		},
		{
			name:          "F nodes stop at block 5",
			numValidators: 7,
			numBlocks:     100,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 5),
			},
			stopTime: make(map[string]time.Time),
			maliciousPeers: map[string]injectors{
				"VD": {},
				"VE": {},
			},
		},
		{
			name:          "F nodes stop at blocks 4,5",
			numValidators: 7,
			numBlocks:     100,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VD": hookStopNode("VD", 4),
				"VE": hookStopNode("VE", 5),
			},
			stopTime: make(map[string]time.Time),
			maliciousPeers: map[string]injectors{
				"VD": {},
				"VE": {},
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

func TestTendermintStartStopSingleNode(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cases := []*testCase{
		{
			name:          "one node stops for 5 seconds",
			numValidators: 5,
			numBlocks:     200,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VE": hookStopNode("VE", 5),
			},
			afterHooks: map[string]hook{
				"VE": hookStartNode("VE", 0.5),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "one node stops for 10 seconds",
			numValidators: 5,
			numBlocks:     200,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VE": hookStopNode("VE", 5),
			},
			afterHooks: map[string]hook{
				"VE": hookStartNode("VE", 1),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "one node stops for 20 seconds",
			numValidators: 5,
			numBlocks:     300,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VE": hookStopNode("VE", 5),
			},
			afterHooks: map[string]hook{
				"VE": hookStartNode("VE", 2),
			},
			stopTime: make(map[string]time.Time),
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
			runTest(t, testCase)
		})
	}
}

func TestTendermintStartStopFNodes(t *testing.T) {
	// to be tracked by https://github.com/clearmatics/autonity/issues/604
	t.Skip("skipping test since the upstream update cause local e2e test framework go routine leak.")
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cases := []*testCase{
		{
			name:          "f nodes stop for 5 seconds at the same block",
			numValidators: 7,
			numBlocks:     150,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 5),
			},
			afterHooks: map[string]hook{
				"VD": hookStartNode("VD", 0.5),
				"VE": hookStartNode("VE", 0.5),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f nodes stop for 5 seconds at different blocks",
			numValidators: 7,
			numBlocks:     150,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 6),
			},
			afterHooks: map[string]hook{
				"VD": hookStartNode("VD", 0.5),
				"VE": hookStartNode("VE", 0.5),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f nodes stop for 10 seconds at the same block",
			numValidators: 7,
			numBlocks:     200,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 5),
			},
			afterHooks: map[string]hook{
				"VD": hookStartNode("VD", 1),
				"VE": hookStartNode("VE", 1),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f nodes stop for 10 seconds at different blocks",
			numValidators: 7,
			numBlocks:     200,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 6),
			},
			afterHooks: map[string]hook{
				"VD": hookStartNode("VD", 1),
				"VE": hookStartNode("VE", 1),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f nodes stop for 10 seconds at the same block",
			numValidators: 7,
			numBlocks:     200,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 5),
			},
			afterHooks: map[string]hook{
				"VD": hookStartNode("VD", 1),
				"VE": hookStartNode("VE", 1),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f nodes stop for 10 seconds at different blocks",
			numValidators: 7,
			numBlocks:     200,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 6),
			},
			afterHooks: map[string]hook{
				"VD": hookStartNode("VD", 1),
				"VE": hookStartNode("VE", 1),
			},
			stopTime: make(map[string]time.Time),
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
			runTest(t, testCase)
		})
	}
}

func TestTendermintStartStopFPlusOneNodes(t *testing.T) {
	t.Skip("Times out intermittently but quite often")
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cases := []*testCase{
		{
			name:          "f+1 nodes stop for 5 seconds at the same block",
			numValidators: 5,
			numBlocks:     150,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 5),
			},
			afterHooks: map[string]hook{
				"VD": hookStartNode("VD", 0.5),
				"VE": hookStartNode("VE", 0.5),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f+1 nodes stop for 5 seconds at different blocks",
			numValidators: 5,
			numBlocks:     150,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 6),
			},
			afterHooks: map[string]hook{
				"VD": hookStartNode("VD", 0.5),
				"VE": hookStartNode("VE", 0.5),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f+1 nodes stop for 10 seconds at the same block",
			numValidators: 5,
			numBlocks:     200,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 5),
			},
			afterHooks: map[string]hook{
				"VD": hookStartNode("VD", 1),
				"VE": hookStartNode("VE", 1),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f+1 nodes stop for 10 seconds at different blocks",
			numValidators: 5,
			numBlocks:     200,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 6),
			},
			afterHooks: map[string]hook{
				"VD": hookStartNode("VD", 1),
				"VE": hookStartNode("VE", 1),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f+1 nodes stop for 20 seconds at the same block",
			numValidators: 5,
			numBlocks:     200,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 5),
			},
			afterHooks: map[string]hook{
				"VD": hookStartNode("VD", 2),
				"VE": hookStartNode("VE", 2),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f+1 nodes stop for 20 seconds at different blocks",
			numValidators: 5,
			numBlocks:     200,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 6),
			},
			afterHooks: map[string]hook{
				"VD": hookStartNode("VD", 2),
				"VE": hookStartNode("VE", 2),
			},
			stopTime: make(map[string]time.Time),
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
			runTest(t, testCase)
		})
	}
}

func TestTendermintStartStopFPlusTwoNodes(t *testing.T) {
	t.Skip("This test fails intermittently see https://github.com/clearmatics/autonity/issues/624")
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cases := []*testCase{
		{
			name:          "f+2 nodes stop for 5 seconds at the same block",
			numValidators: 5,
			numBlocks:     200,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VC": hookStopNode("VC", 5),
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 5),
			},
			afterHooks: map[string]hook{
				"VC": hookStartNode("VC", 0.5),
				"VD": hookStartNode("VD", 0.5),
				"VE": hookStartNode("VE", 0.5),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f+2 nodes stop for 5 seconds at different blocks",
			numValidators: 5,
			numBlocks:     200,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VC": hookStopNode("VC", 4),
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 7),
			},
			afterHooks: map[string]hook{
				"VC": hookStartNode("VC", 0.5),
				"VD": hookStartNode("VD", 0.5),
				"VE": hookStartNode("VE", 0.5),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f+2 nodes stop for 10 seconds at the same block",
			numValidators: 5,
			numBlocks:     300,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VC": hookStopNode("VC", 5),
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 5),
			},
			afterHooks: map[string]hook{
				"VC": hookStartNode("VC", 1),
				"VD": hookStartNode("VD", 1),
				"VE": hookStartNode("VE", 1),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f+2 nodes stop for 10 seconds at different blocks",
			numValidators: 5,
			numBlocks:     300,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VC": hookStopNode("VC", 4),
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 7),
			},
			afterHooks: map[string]hook{
				"VC": hookStartNode("VC", 1),
				"VD": hookStartNode("VD", 1),
				"VE": hookStartNode("VE", 1),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f+2 nodes stop for 20 seconds at the same block",
			numValidators: 5,
			numBlocks:     300,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VC": hookStopNode("VC", 5),
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 5),
			},
			afterHooks: map[string]hook{
				"VC": hookStartNode("VC", 2),
				"VD": hookStartNode("VD", 2),
				"VE": hookStartNode("VE", 2),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f+2 nodes stop for 20 seconds at different blocks",
			numValidators: 5,
			numBlocks:     300,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VC": hookStopNode("VC", 4),
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 7),
			},
			afterHooks: map[string]hook{
				"VC": hookStartNode("VC", 2),
				"VD": hookStartNode("VD", 2),
				"VE": hookStartNode("VE", 2),
			},
			stopTime: make(map[string]time.Time),
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
			runTest(t, testCase)
		})
	}
}

func TestTendermintStartStopAllNodes(t *testing.T) {
	// Track by https://github.com/clearmatics/autonity/issues/711
	// failed due to 40m timeout happens sometimes.
	// https://github.com/clearmatics/autonity/runs/1315225586?check_suite_focus=true
	t.Skip("case is flaky due to 40m timeouts.")
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cases := []*testCase{
		{
			name:          "all nodes stop for 60 seconds at different blocks(2+2+1)",
			numValidators: 5,
			numBlocks:     500,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VA": hookStopNode("VA", 3),
				"VB": hookStopNode("VB", 3),
				"VC": hookStopNode("VC", 5),
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 7),
			},
			afterHooks: map[string]hook{
				"VA": hookStartNode("VA", 6),
				"VB": hookStartNode("VB", 6),
				"VC": hookStartNode("VC", 6),
				"VD": hookStartNode("VD", 6),
				"VE": hookStartNode("VE", 6),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "all nodes stop for 60 seconds at different blocks (2+3)",
			numValidators: 5,
			numBlocks:     500,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VA": hookStopNode("VA", 3),
				"VB": hookStopNode("VB", 3),
				"VC": hookStopNode("VC", 5),
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 5),
			},
			afterHooks: map[string]hook{
				"VA": hookStartNode("VA", 6),
				"VB": hookStartNode("VB", 6),
				"VC": hookStartNode("VC", 6),
				"VD": hookStartNode("VD", 6),
				"VE": hookStartNode("VE", 6),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "all nodes stop for 30 seconds at the same block",
			numValidators: 5,
			numBlocks:     500,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VA": hookStopNode("VA", 3),
				"VB": hookStopNode("VB", 3),
				"VC": hookStopNode("VC", 3),
				"VD": hookStopNode("VD", 3),
				"VE": hookStopNode("VE", 3),
			},
			afterHooks: map[string]hook{
				"VA": hookStartNode("VA", 3),
				"VB": hookStartNode("VB", 3),
				"VC": hookStartNode("VC", 3),
				"VD": hookStartNode("VD", 3),
				"VE": hookStartNode("VE", 3),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "all nodes stop for 60 seconds at the same block",
			numValidators: 5,
			numBlocks:     500,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VA": hookStopNode("VA", 3),
				"VB": hookStopNode("VB", 3),
				"VC": hookStopNode("VC", 3),
				"VD": hookStopNode("VD", 3),
				"VE": hookStopNode("VE", 3),
			},
			afterHooks: map[string]hook{
				"VA": hookStartNode("VA", 6),
				"VB": hookStartNode("VB", 6),
				"VC": hookStartNode("VC", 6),
				"VD": hookStartNode("VD", 6),
				"VE": hookStartNode("VE", 6),
			},
			stopTime: make(map[string]time.Time),
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
			runTest(t, testCase)
		})
	}
}
