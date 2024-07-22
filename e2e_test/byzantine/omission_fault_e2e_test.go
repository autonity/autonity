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

// No omission fault, no validator is addressed as faulty.
func TestHappyCase(t *testing.T) {

	network := createNetwork(t, 6, uint64(150), uint64(10), true)
	defer network.Shutdown(t)
	network.WaitToMineNBlocks(200, 200, false)

	endPoint := network[0].WsClient
	omissionContract, err := autonity.NewOmissionAccountability(params.OmissionAccountabilityContractAddress, endPoint)
	require.NoError(t, err)
	for _, n := range network {
		score, err := omissionContract.GetInactivityScore(nil, n.Address)
		require.NoError(t, err)
		require.Equal(t, uint64(0), score.Uint64())
	}
}

// Let a single node be omission faulty, and then recover it.
func TestSingleFaultyNode(t *testing.T) {
	network := createNetwork(t, 6, uint64(150), uint64(10), true)
	defer network.Shutdown(t)

	// wait for the consensus engine to work.
	network.WaitToMineNBlocks(20, 20, false)

	// stop the 1st node in the network for about 200 blocks/seconds.
	faultyNode := network[0].Address
	err := network[0].Close(false)
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
func TestFFaultyNodes(t *testing.T) {

	numOfNodes := 6
	numOfFaultyNodes := 2
	network := createNetwork(t, numOfNodes, uint64(150), uint64(10), false)
	defer network.Shutdown(t)

	// start quorum nodes only, let F nodes be faulty from the very beginning.
	faultyNodes := make([]common.Address, numOfFaultyNodes)
	for i, n := range network {
		if i >= numOfFaultyNodes {
			err := n.Start()
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

// Edge case, a node keeps resetting for every 20 seconds.
func TestNodeKeepResetting(t *testing.T) {

	network := createNetwork(t, 5, uint64(200), uint64(30), true)
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

// Edge case, there are over quorum node keeps resetting, the resetting nodes are faulty.
func TestNoQuorumWithNodeKeepsResetting(t *testing.T) {
	network := createNetwork(t, 4, uint64(200), uint64(30), true)
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

// Node has an accident resetting within epoch, we will need to resolve a reasonable range of the look-back window to
// allow those node operators to reset the node without being addressed as faulty. A reasonable range of look-back window
// depends on many factors:
// 1. The amount of chain data to be downloaded to get client synced again.
// 2. The HW ability of the client, computing power, memory, etc...
// 3. The delays in between ACN and execution peers.
// Base on my local test bed, 20 blocks at least is okay for the context without any TXNs volume during the test.
func TestNodeHasOneResetWithinEpoch(t *testing.T) {
	numOfNodes := 5
	epochPeriod := uint64(300)
	// set lookBack window to 120 blocks/seconds base on my local e2e test feedback,
	// otherwise the node still be addressed as faulty.
	lookBackWindow := uint64(20)

	network := createNetwork(t, numOfNodes, epochPeriod, lookBackWindow, true)
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

func createNetwork(t *testing.T, nodes int, epochPeriod uint64, lookBackWindow uint64, start bool) e2e.Network {
	validators, err := e2e.Validators(t, nodes, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	network, err := e2e.NewNetworkFromValidators(t, validators, start, func(genesis *core.Genesis) {
		genesis.Config.AutonityContractConfig.EpochPeriod = epochPeriod
		genesis.Config.OmissionAccountabilityConfig.LookbackWindow = lookBackWindow
	})
	require.NoError(t, err)
	return network
}
