// +build race

package test

import (
	"fmt"
	"testing"
	"time"
)

func TestTendermintDataRace(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cases := []*testCase{
		{
			name:      "no malicious",
			numPeers:  5,
			numBlocks: 5,
			txPerPeer: 1,
		},
		{
			name:      "no malicious - 30 tx per second",
			numPeers:  5,
			numBlocks: 10,
			txPerPeer: 30,
		},
		{
			name:      "no malicious - 30 tx per second, 60 blocks",
			numPeers:  5,
			numBlocks: 60,
			txPerPeer: 30,
		},
		{
			name:      "one node stops for 5 seconds",
			isSkipped: true,
			numPeers:  5,
			numBlocks: 10,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				4: hookStopNode(4, 5),
			},
			afterHooks: map[int]hook{
				4: hookStartNode(4, 5),
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
