package test

import (
	"fmt"
	"testing"
	"time"

)

func TestTendermintNoQuorum(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cases := []*testCase{
		{
			name:               "2 validators, one goes down after block 3",
			numPeers:           2,
			numBlocks:          5,
			txPerPeer:          1,
			noQuorumAfterBlock: 3,
			beforeHooks: map[int]hook{
				1: hookForceStopNode(1, 3),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:               "3 validators, two go down after block 3",
			numPeers:           3,
			numBlocks:          5,
			txPerPeer:          1,
			noQuorumAfterBlock: 3,
			noQuorumTimeout:    time.Second * 3,
			beforeHooks: map[int]hook{
				1: hookForceStopNode(1, 3),
				2: hookForceStopNode(2, 3),
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
