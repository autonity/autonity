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
			name:      "one node stops at block 1",
			numPeers:  5,
			numBlocks: 10,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				4: hookStopNode(4, 1),
			},
			stopTime: make(map[int]time.Time),
			maliciousPeers: map[int]injectors{
				4: {},
			},
		},
		{
			name:      "one node stops at block 5",
			numPeers:  5,
			numBlocks: 10,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				4: hookStopNode(4, 5),
			},
			stopTime: make(map[int]time.Time),
			maliciousPeers: map[int]injectors{
				4: {},
			},
		},
		{
			name:      "F nodes stop at block 1",
			numPeers:  7,
			numBlocks: 10,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				3: hookStopNode(3, 1),
				4: hookStopNode(4, 1),
			},
			stopTime: make(map[int]time.Time),
			maliciousPeers: map[int]injectors{
				3: {},
				4: {},
			},
		},
		{
			name:      "F nodes stop at block 5",
			numPeers:  7,
			numBlocks: 10,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 5),
			},
			stopTime: make(map[int]time.Time),
			maliciousPeers: map[int]injectors{
				3: {},
				4: {},
			},
		},
		{
			name:      "F nodes stop at blocks 4,5",
			numPeers:  7,
			numBlocks: 10,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				3: hookStopNode(3, 4),
				4: hookStopNode(4, 5),
			},
			stopTime: make(map[int]time.Time),
			maliciousPeers: map[int]injectors{
				3: {},
				4: {},
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
			name:      "one node stops for 5 seconds",
			numPeers:  5,
			numBlocks: 20,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				4: hookStopNode(4, 5),
			},
			afterHooks: map[int]hook{
				4: hookStartNode(4, 5),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "one node stops for 10 seconds",
			numPeers:  5,
			numBlocks: 20,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				4: hookStopNode(4, 5),
			},
			afterHooks: map[int]hook{
				4: hookStartNode(4, 10),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "one node stops for 20 seconds",
			numPeers:  5,
			numBlocks: 30,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				4: hookStopNode(4, 5),
			},
			afterHooks: map[int]hook{
				4: hookStartNode(4, 20),
			},
			stopTime: make(map[int]time.Time),
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
			name:      "f nodes stop for 5 seconds at the same block",
			numPeers:  7,
			numBlocks: 15,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 5),
			},
			afterHooks: map[int]hook{
				3: hookStartNode(3, 5),
				4: hookStartNode(4, 5),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f nodes stop for 5 seconds at different blocks",
			numPeers:  7,
			numBlocks: 15,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 6),
			},
			afterHooks: map[int]hook{
				3: hookStartNode(3, 5),
				4: hookStartNode(4, 5),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f nodes stop for 10 seconds at the same block",
			numPeers:  7,
			numBlocks: 20,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 5),
			},
			afterHooks: map[int]hook{
				3: hookStartNode(3, 10),
				4: hookStartNode(4, 10),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f nodes stop for 10 seconds at different blocks",
			numPeers:  7,
			numBlocks: 20,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 6),
			},
			afterHooks: map[int]hook{
				3: hookStartNode(3, 10),
				4: hookStartNode(4, 10),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f nodes stop for 10 seconds at the same block",
			numPeers:  7,
			numBlocks: 20,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 5),
			},
			afterHooks: map[int]hook{
				3: hookStartNode(3, 10),
				4: hookStartNode(4, 10),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f nodes stop for 10 seconds at different blocks",
			numPeers:  7,
			numBlocks: 20,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 6),
			},
			afterHooks: map[int]hook{
				3: hookStartNode(3, 10),
				4: hookStartNode(4, 10),
			},
			stopTime: make(map[int]time.Time),
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
			name:      "f+1 nodes stop for 5 seconds at the same block",
			numPeers:  5,
			numBlocks: 15,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 5),
			},
			afterHooks: map[int]hook{
				3: hookStartNode(3, 5),
				4: hookStartNode(4, 5),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f+1 nodes stop for 5 seconds at different blocks",
			numPeers:  5,
			numBlocks: 15,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 6),
			},
			afterHooks: map[int]hook{
				3: hookStartNode(3, 5),
				4: hookStartNode(4, 5),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f+1 nodes stop for 10 seconds at the same block",
			numPeers:  5,
			numBlocks: 20,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 5),
			},
			afterHooks: map[int]hook{
				3: hookStartNode(3, 10),
				4: hookStartNode(4, 10),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f+1 nodes stop for 10 seconds at different blocks",
			numPeers:  5,
			numBlocks: 20,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 6),
			},
			afterHooks: map[int]hook{
				3: hookStartNode(3, 10),
				4: hookStartNode(4, 10),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f+1 nodes stop for 20 seconds at the same block",
			numPeers:  5,
			numBlocks: 20,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 5),
			},
			afterHooks: map[int]hook{
				3: hookStartNode(3, 20),
				4: hookStartNode(4, 20),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f+1 nodes stop for 20 seconds at different blocks",
			numPeers:  5,
			numBlocks: 20,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 6),
			},
			afterHooks: map[int]hook{
				3: hookStartNode(3, 20),
				4: hookStartNode(4, 20),
			},
			stopTime: make(map[int]time.Time),
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
			name:      "f+2 nodes stop for 5 seconds at the same block",
			numPeers:  5,
			numBlocks: 20,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				2: hookStopNode(2, 5),
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 5),
			},
			afterHooks: map[int]hook{
				2: hookStartNode(2, 5),
				3: hookStartNode(3, 5),
				4: hookStartNode(4, 5),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f+2 nodes stop for 5 seconds at different blocks",
			numPeers:  5,
			numBlocks: 20,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				2: hookStopNode(2, 4),
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 7),
			},
			afterHooks: map[int]hook{
				2: hookStartNode(2, 5),
				3: hookStartNode(3, 5),
				4: hookStartNode(4, 5),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f+2 nodes stop for 10 seconds at the same block",
			numPeers:  5,
			numBlocks: 20,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				2: hookStopNode(2, 5),
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 5),
			},
			afterHooks: map[int]hook{
				2: hookStartNode(2, 10),
				3: hookStartNode(3, 10),
				4: hookStartNode(4, 10),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f+2 nodes stop for 10 seconds at different blocks",
			numPeers:  5,
			numBlocks: 20,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				2: hookStopNode(2, 4),
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 7),
			},
			afterHooks: map[int]hook{
				2: hookStartNode(2, 10),
				3: hookStartNode(3, 10),
				4: hookStartNode(4, 10),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f+2 nodes stop for 20 seconds at the same block",
			numPeers:  5,
			numBlocks: 20,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				2: hookStopNode(2, 5),
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 5),
			},
			afterHooks: map[int]hook{
				2: hookStartNode(2, 20),
				3: hookStartNode(3, 20),
				4: hookStartNode(4, 20),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f+2 nodes stop for 20 seconds at different blocks",
			numPeers:  5,
			numBlocks: 20,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				2: hookStopNode(2, 4),
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 7),
			},
			afterHooks: map[int]hook{
				2: hookStartNode(2, 20),
				3: hookStartNode(3, 20),
				4: hookStartNode(4, 20),
			},
			stopTime: make(map[int]time.Time),
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
			name:      "all nodes stop for 60 seconds at different blocks(2+2+1)",
			numPeers:  5,
			numBlocks: 50,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				0: hookStopNode(0, 3),
				1: hookStopNode(1, 3),
				2: hookStopNode(2, 5),
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 7),
			},
			afterHooks: map[int]hook{
				0: hookStartNode(0, 60),
				1: hookStartNode(1, 60),
				2: hookStartNode(2, 60),
				3: hookStartNode(3, 60),
				4: hookStartNode(4, 60),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "all nodes stop for 60 seconds at different blocks (2+3)",
			numPeers:  5,
			numBlocks: 50,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				0: hookStopNode(0, 3),
				1: hookStopNode(1, 3),
				2: hookStopNode(2, 5),
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 5),
			},
			afterHooks: map[int]hook{
				0: hookStartNode(0, 60),
				1: hookStartNode(1, 60),
				2: hookStartNode(2, 60),
				3: hookStartNode(3, 60),
				4: hookStartNode(4, 60),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "all nodes stop for 30 seconds at the same block",
			numPeers:  5,
			numBlocks: 50,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				0: hookStopNode(0, 3),
				1: hookStopNode(1, 3),
				2: hookStopNode(2, 3),
				3: hookStopNode(3, 3),
				4: hookStopNode(4, 3),
			},
			afterHooks: map[int]hook{
				0: hookStartNode(0, 30),
				1: hookStartNode(1, 30),
				2: hookStartNode(2, 30),
				3: hookStartNode(3, 30),
				4: hookStartNode(4, 30),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "all nodes stop for 60 seconds at the same block",
			numPeers:  5,
			numBlocks: 50,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				0: hookStopNode(0, 3),
				1: hookStopNode(1, 3),
				2: hookStopNode(2, 3),
				3: hookStopNode(3, 3),
				4: hookStopNode(4, 3),
			},
			afterHooks: map[int]hook{
				0: hookStartNode(0, 60),
				1: hookStartNode(1, 60),
				2: hookStartNode(2, 60),
				3: hookStartNode(3, 60),
				4: hookStartNode(4, 60),
			},
			stopTime: make(map[int]time.Time),
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
			runTest(t, testCase)
		})
	}
}
