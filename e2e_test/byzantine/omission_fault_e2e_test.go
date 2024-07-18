package byzantine

import (
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core"
	e2e "github.com/autonity/autonity/e2e_test"
	"github.com/autonity/autonity/params"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOmissionFaultSingleFaultyNode(t *testing.T) {
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
	require.NoError(t, err)
	network[0].Wait()

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

	for _, n := range network[1:] {
		score, err := omissionContract.GetInactivityScore(nil, n.Address)
		require.NoError(t, err)
		require.Equal(t, uint64(0), score.Uint64())
	}
}

// Make 1/3 members be omission faulty out of the committee, make the net proposer effort and total net effort to zero
// during the period. The network should keep live-ness and faulty nodes should be detected.
func TestOmissionFaultFFaultyNodes(t *testing.T) {
	numOfNodes := 6
	numOfFaultyNodes := 2
	validators, err := e2e.Validators(t, numOfNodes, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	network, err := e2e.NewNetworkFromValidators(t, validators, false, func(genesis *core.Genesis) {
		genesis.Config.AutonityContractConfig.EpochPeriod = 150
		genesis.Config.OmissionAccountabilityConfig.LookbackWindow = 10
	})
	require.NoError(t, err)
	defer network.Shutdown(t)

	// start quorum nodes only, let F nodes be faulty from the very beginning.
	faultyNodes := make([]common.Address, numOfFaultyNodes)
	for i, n := range network {
		if i >= numOfFaultyNodes {
			err = n.Start()
			require.NoError(t, err)
		} else {
			faultyNodes[i] = n.Address
		}
	}

	// network should be up and continue to mine blocks
	network.WaitToMineNBlocks(200, 200, false)

	endPoint := network[numOfFaultyNodes].WsClient
	omissionContract, err := autonity.NewOmissionAccountability(params.OmissionAccountabilityContractAddress, endPoint)
	require.NoError(t, err)

	for _, node := range faultyNodes {
		inActivityScore, err := omissionContract.GetInactivityScore(nil, node)
		require.NoError(t, err)
		require.Greater(t, inActivityScore.Uint64(), uint64(0))
	}

	// start faulty node, the inactivity score shouldn't increase.
	for i, n := range network {
		if i < numOfNodes {
			err = n.Start()
			require.NoError(t, err)
		}
	}

	// network should be up and continue to mine blocks
	network.WaitToMineNBlocks(150, 150, false)
	for _, n := range network[numOfFaultyNodes:] {
		score, err := omissionContract.GetInactivityScore(nil, n.Address)
		require.NoError(t, err)
		require.Equal(t, uint64(0), score.Uint64())
	}
}
