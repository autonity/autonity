package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/zimmski/go-leak"
	"gonum.org/v1/gonum/stat"
)

func TestTendermintSuccess(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cases := []*testCase{
		{
			name:          "no malicious",
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

func TestTendermintSlowConnections(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cases := []*testCase{
		{
			name:          "no malicious, one slow node",
			numValidators: 5,
			numBlocks:     5,
			txPerPeer:     1,
			networkRates: map[string]networkRate{
				"VE": {50 * 1024, 50 * 1024},
			},
		},
		{
			name:          "no malicious, all nodes are slow",
			numValidators: 5,
			numBlocks:     5,
			txPerPeer:     1,
			networkRates: map[string]networkRate{
				"VA": {50 * 1024, 50 * 1024},
				"VB": {50 * 1024, 50 * 1024},
				"VC": {50 * 1024, 50 * 1024},
				"VD": {50 * 1024, 50 * 1024},
				"VE": {50 * 1024, 50 * 1024},
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

type stats struct {
	mean   float64
	std    float64
	stdErr float64
	n      int
}

func TestTendermintMemoryLeak(t *testing.T) {
	t.Skip("Fails")

	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	const thresholdPerBlock = 1024 // bytes

	cases := []*testCase{
		{
			name:          "5 nodes, 10 blocks, 30 tx per peer per block",
			numValidators: 5,
			numBlocks:     10,
			txPerPeer:     30,
		},

		{
			name:          "10 nodes, 40 blocks, 20 tx per peer per block",
			numValidators: 10,
			numBlocks:     40,
			txPerPeer:     20,
		},
	}

	const repeats = 10

	for _, testCase := range cases {
		testCase := testCase

		leaks := make([]float64, repeats)
		for n := 0; n < repeats; n++ {
			n := n
			t.Run(fmt.Sprintf("test case %s, try %d", testCase.name, n), func(t *testing.T) {
				m := leak.MarkMemory()
				runTest(t, testCase)
				leaks[n] = float64(m.Release()) / float64(testCase.numBlocks)
			})
		}

		if err := checkLeaks(leaks, thresholdPerBlock); err != nil {
			t.Error(err)
		}
	}
}

func checkLeaks(leakStats []float64, threshold float64) error {
	mean, std := stat.MeanStdDev(leakStats, nil)
	stdErr := stat.StdErr(std, float64(len(leakStats)))

	st := stats{
		mean:   mean,
		std:    std,
		stdErr: stdErr,
		n:      len(leakStats),
	}

	if threshold < st.mean+st.stdErr {
		return fmt.Errorf("mean %v; std %v; stdError %v; threshold %v", st.mean, st.std, st.stdErr, threshold)
	}

	return nil
}

func TestTendermintLongRun(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cases := []*testCase{
		{
			name:          "no malicious - 30 tx per second",
			numValidators: 5,
			numBlocks:     10,
			txPerPeer:     30,
		},
		{
			name:          "no malicious - 100 blocks",
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

func TestTendermintTC7(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	test := &testCase{
		name:          "3 nodes stop, 1 recover and sync blocks and state",
		numValidators: 6,
		numBlocks:     40,
		txPerPeer:     1,
		beforeHooks: map[string]hook{
			"VD": hookStopNode("VD", 10),
			"VE": hookStopNode("VE", 15),
			"VF": hookStopNode("VF", 20),
		},
		afterHooks: map[string]hook{
			"VD": hookStartNode("VD", 40),
		},
		maliciousPeers: map[string]injectors{
			"VE": {},
			"VF": {},
		},
		stopTime: make(map[string]time.Time),
	}

	for i := 0; i < 2; i++ {
		t.Run(fmt.Sprintf("test case %s - %d", test.name, i), func(t *testing.T) {
			runTest(t, test)
		})
	}
}
