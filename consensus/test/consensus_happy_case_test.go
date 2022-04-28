package test

import (
	"fmt"
	"testing"
)

func TestTendermintHappyCase(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cases := []*testCase{
		{
			name:          "happy case",
			numValidators: 5,
			numBlocks:     5,
			txPerPeer:     1,
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
			runTest(t, testCase)
		})
	}
}

func TestTendermintHappyCaseLongRun(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cases := []*testCase{
		{
			name:          "Happy case - 30 tx per second",
			numValidators: 5,
			numBlocks:     10,
			txPerPeer:     30,
		},
		{
			name:          "Happy case - 100 blocks",
			numValidators: 5,
			numBlocks:     100,
			txPerPeer:     5,
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
			runTest(t, testCase)
		})
	}
}
