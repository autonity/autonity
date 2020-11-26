package test

import (
	"fmt"
	"github.com/clearmatics/autonity/autonity"
	"github.com/clearmatics/autonity/metrics"
	"github.com/stretchr/testify/require"
	"testing"
)

// This test case verified if network economic metrics are exported from autonity contract during the runtime and they
// should be presented at the global metric store.
func TestExportEconomicMetrics(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	// if economic metrics were exported successfully, then they should be presented at metric store.
	metricChecker := func(t *testing.T, validators map[string]*testNode) {
		require.NotNil(t, metrics.Get(autonity.GlobalMetricIDGasPrice))
		require.NotNil(t, metrics.Get(autonity.GlobalOperatorBalanceMetricID))
		require.NotNil(t, metrics.Get(autonity.GlobalMetricIDStakeSupply))
	}
	cases := []*testCase{
		{
			name:          "Network economic metrics should be exported and store at metric store successfully",
			numValidators: 6,
			numBlocks:     5,
			txPerPeer:     1,
			finalAssert:   metricChecker,
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
			runTest(t, testCase)
		})
	}
}
