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
			numValidators:      2,
			numBlocks:          5,
			txPerPeer:          1,
			noQuorumAfterBlock: 3,
			beforeHooks: map[string]hook{
				"VB": hookForceStopNode("VB", 3),
			},
			stopTime: make(map[string]time.Time),
		},
		{
			name:               "3 validators, two go down after block 3",
			numValidators:      3,
			numBlocks:          5,
			txPerPeer:          1,
			noQuorumAfterBlock: 3,
			noQuorumTimeout:    time.Second * 3,
			beforeHooks: map[string]hook{
				"VB": hookForceStopNode("VB", 3),
				"VC": hookForceStopNode("VC", 3),
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
