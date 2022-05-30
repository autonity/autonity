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
			name:          "Happy case - 10 blocks",
			numValidators: 5,
			numBlocks:     10,
		},
		{
			name:          "Happy case - 100 blocks",
			numValidators: 5,
			numBlocks:     100,
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
			runTest(t, testCase)
		})
	}
}
