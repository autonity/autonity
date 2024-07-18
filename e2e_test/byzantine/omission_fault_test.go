package byzantine

import (
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/core"
	e2e "github.com/autonity/autonity/e2e_test"
	"github.com/autonity/autonity/params"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSingleValidatorOmissionFault(t *testing.T) {
	numOfNodes := 6
	validators, err := e2e.Validators(t, numOfNodes, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	network, err := e2e.NewNetworkFromValidators(t, validators, true, func(genesis *core.Genesis) {
		genesis.Config.AutonityContractConfig.EpochPeriod = 150
		genesis.Config.OmissionAccountabilityConfig.LookbackWindow = 10
	})
	require.NoError(t, err)
	defer network.Shutdown(t)

	// wait for the consensus engine to work.
	network.WaitToMineNBlocks(20, 20, false)

	// stop the 1st node in the network for about 200 blocks/seconds.
	faultyNode := network[0].Address
	err = network[0].Close(false)
	network[0].Wait()
	require.NoError(t, err)

	// network should be up and continue to mine blocks
	network.WaitToMineNBlocks(200, 200, false)

	endPoint := network[1].WsClient
	omissionContract, err := autonity.NewOmissionAccountability(params.OmissionAccountabilityContractAddress, endPoint)
	require.NoError(t, err)

	inActivityScore, err := omissionContract.GetInactivityScore(nil, faultyNode)
	require.NoError(t, err)
	require.Greater(t, inActivityScore.Uint64(), uint64(0))

	// restart faulty node, the inactivity score shouldn't increase.
	err = network[0].Start()
	require.NoError(t, err)

	// network should be up and continue to mine blocks
	network.WaitToMineNBlocks(300, 300, false)

	curScore, err := omissionContract.GetInactivityScore(nil, faultyNode)
	require.NoError(t, err)
	require.Equal(t, inActivityScore.Uint64(), curScore.Uint64())
}
