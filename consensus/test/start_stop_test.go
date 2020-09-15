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
			numBlocks:     10,
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
			numBlocks:     10,
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
			numBlocks:     10,
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
			numBlocks:     10,
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
			numBlocks:     10,
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
			numBlocks:     20,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VE": hookStopNode("VE", 5),
			},
			afterHooks: map[string]hook{
				"VE": hookStartNode("VE", 5),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "one node stops for 10 seconds",
			numValidators: 5,
			numBlocks:     20,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VE": hookStopNode("VE", 5),
			},
			afterHooks: map[string]hook{
				"VE": hookStartNode("VE", 10),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "one node stops for 20 seconds",
			numValidators: 5,
			numBlocks:     30,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VE": hookStopNode("VE", 5),
			},
			afterHooks: map[string]hook{
				"VE": hookStartNode("VE", 20),
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
			numBlocks:     15,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 5),
			},
			afterHooks: map[string]hook{
				"VD": hookStartNode("VD", 5),
				"VE": hookStartNode("VE", 5),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f nodes stop for 5 seconds at different blocks",
			numValidators: 7,
			numBlocks:     15,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 6),
			},
			afterHooks: map[string]hook{
				"VD": hookStartNode("VD", 5),
				"VE": hookStartNode("VE", 5),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f nodes stop for 10 seconds at the same block",
			numValidators: 7,
			numBlocks:     20,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 5),
			},
			afterHooks: map[string]hook{
				"VD": hookStartNode("VD", 10),
				"VE": hookStartNode("VE", 10),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f nodes stop for 10 seconds at different blocks",
			numValidators: 7,
			numBlocks:     20,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 6),
			},
			afterHooks: map[string]hook{
				"VD": hookStartNode("VD", 10),
				"VE": hookStartNode("VE", 10),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f nodes stop for 10 seconds at the same block",
			numValidators: 7,
			numBlocks:     20,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 5),
			},
			afterHooks: map[string]hook{
				"VD": hookStartNode("VD", 10),
				"VE": hookStartNode("VE", 10),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f nodes stop for 10 seconds at different blocks",
			numValidators: 7,
			numBlocks:     20,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 6),
			},
			afterHooks: map[string]hook{
				"VD": hookStartNode("VD", 10),
				"VE": hookStartNode("VE", 10),
			},
			stopTime: make(map[string]time.Time),
		},
	}

	for i := 0; i < 100; i ++ {
		for _, testCase := range cases {
			testCase := testCase
			t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
				runTest(t, testCase)
			})
		}
	}
}

func TestTendermintStartStopFPlusOneNodes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cases := []*testCase{
		{
			name:          "f+1 nodes stop for 5 seconds at the same block",
			numValidators: 5,
			numBlocks:     15,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 5),
			},
			afterHooks: map[string]hook{
				"VD": hookStartNode("VD", 5),
				"VE": hookStartNode("VE", 5),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f+1 nodes stop for 5 seconds at different blocks",
			numValidators: 5,
			numBlocks:     15,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 6),
			},
			afterHooks: map[string]hook{
				"VD": hookStartNode("VD", 5),
				"VE": hookStartNode("VE", 5),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f+1 nodes stop for 10 seconds at the same block",
			numValidators: 5,
			numBlocks:     20,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 5),
			},
			afterHooks: map[string]hook{
				"VD": hookStartNode("VD", 10),
				"VE": hookStartNode("VE", 10),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f+1 nodes stop for 10 seconds at different blocks",
			numValidators: 5,
			numBlocks:     20,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 6),
			},
			afterHooks: map[string]hook{
				"VD": hookStartNode("VD", 10),
				"VE": hookStartNode("VE", 10),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f+1 nodes stop for 20 seconds at the same block",
			numValidators: 5,
			numBlocks:     20,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 5),
			},
			afterHooks: map[string]hook{
				"VD": hookStartNode("VD", 20),
				"VE": hookStartNode("VE", 20),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f+1 nodes stop for 20 seconds at different blocks",
			numValidators: 5,
			numBlocks:     20,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 6),
			},
			afterHooks: map[string]hook{
				"VD": hookStartNode("VD", 20),
				"VE": hookStartNode("VE", 20),
			},
			stopTime: make(map[string]time.Time),
		},
	}

	for i := 0; i < 100; i ++ {
		for _, testCase := range cases {
			testCase := testCase
			t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
				runTest(t, testCase)
			})
		}
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
			numBlocks:     20,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VC": hookStopNode("VC", 5),
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 5),
			},
			afterHooks: map[string]hook{
				"VC": hookStartNode("VC", 5),
				"VD": hookStartNode("VD", 5),
				"VE": hookStartNode("VE", 5),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f+2 nodes stop for 5 seconds at different blocks",
			numValidators: 5,
			numBlocks:     20,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VC": hookStopNode("VC", 4),
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 7),
			},
			afterHooks: map[string]hook{
				"VC": hookStartNode("VC", 5),
				"VD": hookStartNode("VD", 5),
				"VE": hookStartNode("VE", 5),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f+2 nodes stop for 10 seconds at the same block",
			numValidators: 5,
			numBlocks:     30,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VC": hookStopNode("VC", 5),
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 5),
			},
			afterHooks: map[string]hook{
				"VC": hookStartNode("VC", 10),
				"VD": hookStartNode("VD", 10),
				"VE": hookStartNode("VE", 10),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f+2 nodes stop for 10 seconds at different blocks",
			numValidators: 5,
			numBlocks:     30,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VC": hookStopNode("VC", 4),
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 7),
			},
			afterHooks: map[string]hook{
				"VC": hookStartNode("VC", 10),
				"VD": hookStartNode("VD", 10),
				"VE": hookStartNode("VE", 10),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f+2 nodes stop for 20 seconds at the same block",
			numValidators: 5,
			numBlocks:     30,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VC": hookStopNode("VC", 5),
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 5),
			},
			afterHooks: map[string]hook{
				"VC": hookStartNode("VC", 20),
				"VD": hookStartNode("VD", 20),
				"VE": hookStartNode("VE", 20),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f+2 nodes stop for 20 seconds at different blocks",
			numValidators: 5,
			numBlocks:     30,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VC": hookStopNode("VC", 4),
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 7),
			},
			afterHooks: map[string]hook{
				"VC": hookStartNode("VC", 20),
				"VD": hookStartNode("VD", 20),
				"VE": hookStartNode("VE", 20),
			},
			stopTime: make(map[string]time.Time),
		},
	}

	for i := 0; i < 100; i ++ {
		for _, testCase := range cases {
			testCase := testCase
			t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
				runTest(t, testCase)
			})
		}
	}
}

func TestTendermintStartStopAllNodes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cases := []*testCase{
		{
			name:          "all nodes stop for 60 seconds at different blocks(2+2+1)",
			numValidators: 5,
			numBlocks:     50,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VA": hookStopNode("VA", 3),
				"VB": hookStopNode("VB", 3),
				"VC": hookStopNode("VC", 5),
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 7),
			},
			afterHooks: map[string]hook{
				"VA": hookStartNode("VA", 60),
				"VB": hookStartNode("VB", 60),
				"VC": hookStartNode("VC", 60),
				"VD": hookStartNode("VD", 60),
				"VE": hookStartNode("VE", 60),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "all nodes stop for 60 seconds at different blocks (2+3)",
			numValidators: 5,
			numBlocks:     50,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VA": hookStopNode("VA", 3),
				"VB": hookStopNode("VB", 3),
				"VC": hookStopNode("VC", 5),
				"VD": hookStopNode("VD", 5),
				"VE": hookStopNode("VE", 5),
			},
			afterHooks: map[string]hook{
				"VA": hookStartNode("VA", 60),
				"VB": hookStartNode("VB", 60),
				"VC": hookStartNode("VC", 60),
				"VD": hookStartNode("VD", 60),
				"VE": hookStartNode("VE", 60),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "all nodes stop for 30 seconds at the same block",
			numValidators: 5,
			numBlocks:     50,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VA": hookStopNode("VA", 3),
				"VB": hookStopNode("VB", 3),
				"VC": hookStopNode("VC", 3),
				"VD": hookStopNode("VD", 3),
				"VE": hookStopNode("VE", 3),
			},
			afterHooks: map[string]hook{
				"VA": hookStartNode("VA", 30),
				"VB": hookStartNode("VB", 30),
				"VC": hookStartNode("VC", 30),
				"VD": hookStartNode("VD", 30),
				"VE": hookStartNode("VE", 30),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "all nodes stop for 60 seconds at the same block",
			numValidators: 5,
			numBlocks:     50,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VA": hookStopNode("VA", 3),
				"VB": hookStopNode("VB", 3),
				"VC": hookStopNode("VC", 3),
				"VD": hookStopNode("VD", 3),
				"VE": hookStopNode("VE", 3),
			},
			afterHooks: map[string]hook{
				"VA": hookStartNode("VA", 60),
				"VB": hookStartNode("VB", 60),
				"VC": hookStartNode("VC", 60),
				"VD": hookStartNode("VD", 60),
				"VE": hookStartNode("VE", 60),
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
