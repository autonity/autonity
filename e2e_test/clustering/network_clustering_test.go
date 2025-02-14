package clustering

import (
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/e2e_test"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestClusteringHappyCase is a happy case to test 5 clusters with each of them contains 5 nodes. The latency measurement
// is a base on an local simulator which generates [0, 500) ms delays.
func TestClusteringHappyCase(t *testing.T) {
	// mocked service with a local ping simulator which generates [0, 500) ms latency.
	mockedService := &interfaces.Services{Pinger: NewSimulatedPinger()}

	validators, err := e2e.Validators(t, 25, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	for _, validator := range validators {
		validator.TendermintServices = mockedService
	}

	// creates a network of 25 validators and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, validators, true)
	require.NoError(t, err)
	defer network.Shutdown(t)

	// runs for about 10 epoches period.
	_ = network.WaitToMineNBlocks(500, 500, false)
}

// todo: test customized K nodes selectors for message routing.
func TestCustomizedSelectors(t *testing.T) {

}
