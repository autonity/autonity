package byzantine

import (
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core"
	e2e "github.com/autonity/autonity/e2e_test"
	"github.com/autonity/autonity/params"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// Let a single node be off, and then recover it.
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
	require.LessOrEqual(t, curScore.Uint64(), inActivityScore.Uint64())

	for _, n := range network[1:] {
		score, err := omissionContract.GetInactivityScore(nil, n.Address)
		require.NoError(t, err)
		require.Equal(t, uint64(0), score.Uint64())
	}
}

// Edge case: make 1/3 members be omission faulty out of the committee, make the net proposer effort and total net
// effort to zero during the period. The network should keep up and faulty nodes should be detected.
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
	network.WaitToMineNBlocks(300, 300, false)
	endPoint := network[numOfFaultyNodes].WsClient
	omissionContract, err := autonity.NewOmissionAccountability(params.OmissionAccountabilityContractAddress, endPoint)
	require.NoError(t, err)

	autonity, err := autonity.NewAutonity(params.AutonityContractAddress, endPoint)
	require.NoError(t, err)
	epochID, err := autonity.EpochID(nil)
	require.NoError(t, err)
	t.Log("current epoch id", "id", epochID.Uint64())

	for _, n := range network[numOfFaultyNodes:] {
		score, err := omissionContract.GetInactivityScore(nil, n.Address)
		require.NoError(t, err)
		require.Equal(t, uint64(0), score.Uint64())
	}

	for _, n := range faultyNodes {
		score, err := omissionContract.GetInactivityScore(nil, n)
		require.NoError(t, err)
		require.Greater(t, score.Uint64(), uint64(0))
	}
}

// Edge case, a node keeps resetting, as the resetting proposer couldn't provide activity
// proof, it will be addressed as faulty.
func TestNodeKeepRestart(t *testing.T) {
	numOfNodes := 5
	validators, err := e2e.Validators(t, numOfNodes, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	network, err := e2e.NewNetworkFromValidators(t, validators, true, func(genesis *core.Genesis) {
		genesis.Config.AutonityContractConfig.EpochPeriod = 200
	})
	require.NoError(t, err)
	defer network.Shutdown(t)

	for network[1].Eth.BlockChain().CurrentHeader().Number.Uint64() <= 200 {
		err := network[0].Restart()
		require.NoError(t, err)
		time.Sleep(30 * time.Second)
	}

	endPoint := network[1].WsClient
	omissionContract, err := autonity.NewOmissionAccountability(params.OmissionAccountabilityContractAddress, endPoint)
	require.NoError(t, err)

	score, err := omissionContract.GetInactivityScore(nil, network[0].Address)
	require.NoError(t, err)
	require.Greater(t, score.Uint64(), uint64(0))
}

// Edge case, there are over quorum node keeps resetting, the resetting node is faulty as it shouldn't be able to
// provide proof.
func TestNetworkOnholdingWithNodeKeepsResetting(t *testing.T) {
	numOfNodes := 4
	validators, err := e2e.Validators(t, numOfNodes, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	network, err := e2e.NewNetworkFromValidators(t, validators, true, func(genesis *core.Genesis) {
		genesis.Config.AutonityContractConfig.EpochPeriod = 200
	})
	require.NoError(t, err)
	defer network.Shutdown(t)

	for network[1].Eth.BlockChain().CurrentHeader().Number.Uint64() <= 200 {
		err := network[0].Restart()
		require.NoError(t, err)
		time.Sleep(20 * time.Second)
	}

	endPoint := network[1].WsClient
	omissionContract, err := autonity.NewOmissionAccountability(params.OmissionAccountabilityContractAddress, endPoint)
	require.NoError(t, err)

	score, err := omissionContract.GetInactivityScore(nil, network[0].Address)
	require.NoError(t, err)
	require.Greater(t, score.Uint64(), uint64(0))
}

// Node has an accident resetting during epoch, we will need to resolve a reasonable range of the look-back window to
// allow those node operators to finish the node without being addressed as faulty.
func TestNodeHasOneResetDuringEpoch(t *testing.T) {
	numOfNodes := 5
	validators, err := e2e.Validators(t, numOfNodes, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	// todo: (Jason) to resolve minimum requirement of look back window for production network.
	// resolve lookBack window to 120 blocks/seconds base on my local e2e test feedback, otherwise the node still be
	// addressed as faulty.
	lookBackWindow := uint64(120)
	network, err := e2e.NewNetworkFromValidators(t, validators, true, func(genesis *core.Genesis) {
		genesis.Config.AutonityContractConfig.EpochPeriod = 300
		genesis.Config.OmissionAccountabilityConfig.LookbackWindow = lookBackWindow
	})
	require.NoError(t, err)
	defer network.Shutdown(t)

	resetHeight := uint64(60)
	for network[1].Eth.BlockChain().CurrentHeader().Number.Uint64() <= 300 {
		height := network[1].Eth.BlockChain().CurrentHeader().Number.Uint64()
		if height == resetHeight {
			err := network[0].Restart()
			require.NoError(t, err)
		}
		time.Sleep(1 * time.Second)
	}

	endPoint := network[1].WsClient
	omissionContract, err := autonity.NewOmissionAccountability(params.OmissionAccountabilityContractAddress, endPoint)
	require.NoError(t, err)

	score, err := omissionContract.GetInactivityScore(nil, network[0].Address)
	require.NoError(t, err)
	require.Equal(t, uint64(0), score.Uint64())
}
