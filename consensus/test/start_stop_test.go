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
				"E": hookStopNode("E", 1),
			},
			stopTime: make(map[string]time.Time),
			maliciousPeers: map[string]injectors{
				"E": {},
			},
		},
		{
			name:          "one node stops at block 5",
			numValidators: 5,
			numBlocks:     10,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"E": hookStopNode("E", 5),
			},
			stopTime: make(map[string]time.Time),
			maliciousPeers: map[string]injectors{
				"E": {},
			},
		},
		{
			name:          "F nodes stop at block 1",
			numValidators: 7,
			numBlocks:     10,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"D": hookStopNode("D", 1),
				"E": hookStopNode("E", 1),
			},
			stopTime: make(map[string]time.Time),
			maliciousPeers: map[string]injectors{
				"D": {},
				"E": {},
			},
		},
		{
			name:          "F nodes stop at block 5",
			numValidators: 7,
			numBlocks:     10,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"D": hookStopNode("D", 5),
				"E": hookStopNode("E", 5),
			},
			stopTime: make(map[string]time.Time),
			maliciousPeers: map[string]injectors{
				"D": {},
				"E": {},
			},
		},
		{
			name:          "F nodes stop at blocks 4,5",
			numValidators: 7,
			numBlocks:     10,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"D": hookStopNode("D", 4),
				"E": hookStopNode("E", 5),
			},
			stopTime: make(map[string]time.Time),
			maliciousPeers: map[string]injectors{
				"D": {},
				"E": {},
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
				"E": hookStopNode("E", 5),
			},
			afterHooks: map[string]hook{
				"E": hookStartNode("E", 5),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "one node stops for 10 seconds",
			numValidators: 5,
			numBlocks:     20,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"E": hookStopNode("E", 5),
			},
			afterHooks: map[string]hook{
				"E": hookStartNode("E", 10),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "one node stops for 20 seconds",
			numValidators: 5,
			numBlocks:     30,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"E": hookStopNode("E", 5),
			},
			afterHooks: map[string]hook{
				"E": hookStartNode("E", 20),
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
				"D": hookStopNode("D", 5),
				"E": hookStopNode("E", 5),
			},
			afterHooks: map[string]hook{
				"D": hookStartNode("D", 5),
				"E": hookStartNode("E", 5),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f nodes stop for 5 seconds at different blocks",
			numValidators: 7,
			numBlocks:     15,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"D": hookStopNode("D", 5),
				"E": hookStopNode("E", 6),
			},
			afterHooks: map[string]hook{
				"D": hookStartNode("D", 5),
				"E": hookStartNode("E", 5),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f nodes stop for 10 seconds at the same block",
			numValidators: 7,
			numBlocks:     20,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"D": hookStopNode("D", 5),
				"E": hookStopNode("E", 5),
			},
			afterHooks: map[string]hook{
				"D": hookStartNode("D", 10),
				"E": hookStartNode("E", 10),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f nodes stop for 10 seconds at different blocks",
			numValidators: 7,
			numBlocks:     20,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"D": hookStopNode("D", 5),
				"E": hookStopNode("E", 6),
			},
			afterHooks: map[string]hook{
				"D": hookStartNode("D", 10),
				"E": hookStartNode("E", 10),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f nodes stop for 10 seconds at the same block",
			numValidators: 7,
			numBlocks:     20,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"D": hookStopNode("D", 5),
				"E": hookStopNode("E", 5),
			},
			afterHooks: map[string]hook{
				"D": hookStartNode("D", 10),
				"E": hookStartNode("E", 10),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f nodes stop for 10 seconds at different blocks",
			numValidators: 7,
			numBlocks:     20,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"D": hookStopNode("D", 5),
				"E": hookStopNode("E", 6),
			},
			afterHooks: map[string]hook{
				"D": hookStartNode("D", 10),
				"E": hookStartNode("E", 10),
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
				"D": hookStopNode("D", 5),
				"E": hookStopNode("E", 5),
			},
			afterHooks: map[string]hook{
				"D": hookStartNode("D", 5),
				"E": hookStartNode("E", 5),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f+1 nodes stop for 5 seconds at different blocks",
			numValidators: 5,
			numBlocks:     15,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"D": hookStopNode("D", 5),
				"E": hookStopNode("E", 6),
			},
			afterHooks: map[string]hook{
				"D": hookStartNode("D", 5),
				"E": hookStartNode("E", 5),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f+1 nodes stop for 10 seconds at the same block",
			numValidators: 5,
			numBlocks:     20,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"D": hookStopNode("D", 5),
				"E": hookStopNode("E", 5),
			},
			afterHooks: map[string]hook{
				"D": hookStartNode("D", 10),
				"E": hookStartNode("E", 10),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f+1 nodes stop for 10 seconds at different blocks",
			numValidators: 5,
			numBlocks:     20,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"D": hookStopNode("D", 5),
				"E": hookStopNode("E", 6),
			},
			afterHooks: map[string]hook{
				"D": hookStartNode("D", 10),
				"E": hookStartNode("E", 10),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f+1 nodes stop for 20 seconds at the same block",
			numValidators: 5,
			numBlocks:     20,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"D": hookStopNode("D", 5),
				"E": hookStopNode("E", 5),
			},
			afterHooks: map[string]hook{
				"D": hookStartNode("D", 20),
				"E": hookStartNode("E", 20),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f+1 nodes stop for 20 seconds at different blocks",
			numValidators: 5,
			numBlocks:     20,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"D": hookStopNode("D", 5),
				"E": hookStopNode("E", 6),
			},
			afterHooks: map[string]hook{
				"D": hookStartNode("D", 20),
				"E": hookStartNode("E", 20),
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
				"C": hookStopNode("C", 5),
				"D": hookStopNode("D", 5),
				"E": hookStopNode("E", 5),
			},
			afterHooks: map[string]hook{
				"C": hookStartNode("C", 5),
				"D": hookStartNode("D", 5),
				"E": hookStartNode("E", 5),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f+2 nodes stop for 5 seconds at different blocks",
			numValidators: 5,
			numBlocks:     20,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"C": hookStopNode("C", 4),
				"D": hookStopNode("D", 5),
				"E": hookStopNode("E", 7),
			},
			afterHooks: map[string]hook{
				"C": hookStartNode("C", 5),
				"D": hookStartNode("D", 5),
				"E": hookStartNode("E", 5),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f+2 nodes stop for 10 seconds at the same block",
			numValidators: 5,
			numBlocks:     20,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"C": hookStopNode("C", 5),
				"D": hookStopNode("D", 5),
				"E": hookStopNode("E", 5),
			},
			afterHooks: map[string]hook{
				"C": hookStartNode("C", 10),
				"D": hookStartNode("D", 10),
				"E": hookStartNode("E", 10),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f+2 nodes stop for 10 seconds at different blocks",
			numValidators: 5,
			numBlocks:     20,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"C": hookStopNode("C", 4),
				"D": hookStopNode("D", 5),
				"E": hookStopNode("E", 7),
			},
			afterHooks: map[string]hook{
				"C": hookStartNode("C", 10),
				"D": hookStartNode("D", 10),
				"E": hookStartNode("E", 10),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f+2 nodes stop for 20 seconds at the same block",
			numValidators: 5,
			numBlocks:     20,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"C": hookStopNode("C", 5),
				"D": hookStopNode("D", 5),
				"E": hookStopNode("E", 5),
			},
			afterHooks: map[string]hook{
				"C": hookStartNode("C", 20),
				"D": hookStartNode("D", 20),
				"E": hookStartNode("E", 20),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "f+2 nodes stop for 20 seconds at different blocks",
			numValidators: 5,
			numBlocks:     20,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"C": hookStopNode("C", 4),
				"D": hookStopNode("D", 5),
				"E": hookStopNode("E", 7),
			},
			afterHooks: map[string]hook{
				"C": hookStartNode("C", 20),
				"D": hookStartNode("D", 20),
				"E": hookStartNode("E", 20),
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
				"A": hookStopNode("A", 3),
				"B": hookStopNode("B", 3),
				"C": hookStopNode("C", 5),
				"D": hookStopNode("D", 5),
				"E": hookStopNode("E", 7),
			},
			afterHooks: map[string]hook{
				"A": hookStartNode("A", 60),
				"B": hookStartNode("B", 60),
				"C": hookStartNode("C", 60),
				"D": hookStartNode("D", 60),
				"E": hookStartNode("E", 60),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "all nodes stop for 60 seconds at different blocks (2+3)",
			numValidators: 5,
			numBlocks:     50,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"A": hookStopNode("A", 3),
				"B": hookStopNode("B", 3),
				"C": hookStopNode("C", 5),
				"D": hookStopNode("D", 5),
				"E": hookStopNode("E", 5),
			},
			afterHooks: map[string]hook{
				"A": hookStartNode("A", 60),
				"B": hookStartNode("B", 60),
				"C": hookStartNode("C", 60),
				"D": hookStartNode("D", 60),
				"E": hookStartNode("E", 60),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "all nodes stop for 30 seconds at the same block",
			numValidators: 5,
			numBlocks:     50,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"A": hookStopNode("A", 3),
				"B": hookStopNode("B", 3),
				"C": hookStopNode("C", 3),
				"D": hookStopNode("D", 3),
				"E": hookStopNode("E", 3),
			},
			afterHooks: map[string]hook{
				"A": hookStartNode("A", 30),
				"B": hookStartNode("B", 30),
				"C": hookStartNode("C", 30),
				"D": hookStartNode("D", 30),
				"E": hookStartNode("E", 30),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:          "all nodes stop for 60 seconds at the same block",
			numValidators: 5,
			numBlocks:     50,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"A": hookStopNode("A", 3),
				"B": hookStopNode("B", 3),
				"C": hookStopNode("C", 3),
				"D": hookStopNode("D", 3),
				"E": hookStopNode("E", 3),
			},
			afterHooks: map[string]hook{
				"A": hookStartNode("A", 60),
				"B": hookStartNode("B", 60),
				"C": hookStartNode("C", 60),
				"D": hookStartNode("D", 60),
				"E": hookStartNode("E", 60),
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
