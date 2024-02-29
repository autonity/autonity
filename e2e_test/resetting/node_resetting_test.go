package resetting

import (
	e2e "github.com/autonity/autonity/e2e_test"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
)

// TestResetAllNodes, it stops all nodes one by one, and start them again one by one. The network should recover to
// mining.
func TestResetAllNodes(t *testing.T) {
	numOfNodes := 6
	network, err := e2e.NewNetwork(t, numOfNodes, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	defer network.Shutdown()
	// wait for the consensus engine to work.
	network.WaitToMineNBlocks(10, 60, false)

	// stop all nodes.
	for _, n := range network {
		err = n.Close(false)
		n.Wait()
		require.NoError(t, err)
	}

	// start all nodes again.
	for _, n := range network {
		err = n.Start()
		require.NoError(t, err)
	}

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(10, 60, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}

// TestResetFNodes, it stops random selected F nodes one by one, and observe if the net is still mining, then it recover
// F nodes one by one, the network should keep mining all the time.
func TestResetRandomFNodes(t *testing.T) {
	numOfNodes := 6
	network, err := e2e.NewNetwork(t, numOfNodes, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	defer network.Shutdown()
	// wait for the consensus engine to work.
	network.WaitToMineNBlocks(10, 60, false)

	// stop random selected F nodes.
	f := 2
	fNodes := make(map[int]struct{})
	//fNodes := make([]int, f)
	for i := 0; i < f; {
		selectedID := rand.Intn(len(network))
		if _, ok := fNodes[selectedID]; ok {
			continue
		}
		i++
		fNodes[selectedID] = struct{}{}
		err = network[selectedID].Close(false)
		network[selectedID].Wait()
		require.NoError(t, err)
	}

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(10, 60, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")

	for id := range fNodes {
		err = network[id].Start()
		require.NoError(t, err)
	}

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(10, 60, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}

// TestResetFPlusOneNodes, it stops random selected F+1 nodes one by one, and observe if the net is on-holding, then it
// recover anyone of the stopped node to get quorum voting power and observe if the net is mining again.
func TestResetRandomFPlusOneNodes(t *testing.T) {
	numOfNodes := 6
	network, err := e2e.NewNetwork(t, numOfNodes, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	defer network.Shutdown()
	// wait for the consensus engine to work.
	network.WaitToMineNBlocks(10, 60, false)

	// stop random selected F+1 nodes.
	f := 3
	fNodes := make(map[int]struct{})
	for i := 0; i < f; {
		selectedID := rand.Intn(len(network))
		if _, ok := fNodes[selectedID]; ok {
			continue
		}
		i++
		fNodes[selectedID] = struct{}{}
		err = network[selectedID].Close(false)
		network[selectedID].Wait()
		require.NoError(t, err)
	}

	// network should be on holding.
	err = network.WaitToMineNBlocks(10, 20, false)
	require.EqualError(t, err, "context deadline exceeded")

	// recover anyone of the stopped node, the network should produce blocks again.
	for id := range fNodes {
		err = network[id].Start()
		require.NoError(t, err)
		break
	}
	err = network.WaitToMineNBlocks(10, 60, false)
	require.NoError(t, err)
}

// TestResetFPlusTwoNodes, it stops random selected F+2 nodes one by one, and observe if the net is on-holding, then it
// recover random selected any two of the stopped node to get quorum voting power, and observe if the net is mining again.
func TestResetRandomFPlusTwoNodes(t *testing.T) {
	numOfNodes := 6
	network, err := e2e.NewNetwork(t, numOfNodes, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	defer network.Shutdown()
	// wait for the consensus engine to work.
	network.WaitToMineNBlocks(10, 60, false)

	// stop random selected F+2 nodes.
	f := 4
	fNodes := make(map[int]struct{})
	for i := 0; i < f; {
		selectedID := rand.Intn(len(network))
		if _, ok := fNodes[selectedID]; ok {
			continue
		}
		i++
		fNodes[selectedID] = struct{}{}
		err = network[selectedID].Close(false)
		network[selectedID].Wait()
		require.NoError(t, err)
	}

	// network should be on holding.
	err = network.WaitToMineNBlocks(10, 60, false)
	require.EqualError(t, err, "context deadline exceeded")

	// recover anyone two of the stopped node, the network should produce blocks again.
	count := 0
	for id := range fNodes {
		if count >= 2 {
			break
		}
		count++
		err = network[id].Start()
		require.NoError(t, err)
	}
	err = network.WaitToMineNBlocks(10, 60, false)
	require.NoError(t, err)
}

// TestKeepResettingRandomOneNode, with multiple rounds, it keeps resetting a random selected node at each round, then
// recover it at that round, after each round the network should keep mining.
func TestKeepResettingRandomOneNode(t *testing.T) {
	numOfNodes := 6
	network, err := e2e.NewNetwork(t, numOfNodes, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	defer network.Shutdown()
	// wait for the consensus engine to work.
	network.WaitToMineNBlocks(10, 60, false)

	rounds := 10
	for r := 0; r < rounds; r++ {
		// random select a node to reset.
		nodeID := rand.Intn(len(network))
		err = network[nodeID].Close(false)
		network[nodeID].Wait()
		require.NoError(t, err)
		// network should be up and continue to mine blocks
		err = network.WaitToMineNBlocks(10, 60, false)
		require.NoError(t, err, "Network should be mining new blocks now, but it's not")
		// recover that faulty node.
		err = network[nodeID].Start()
		require.NoError(t, err)
		// network should be up and continue to mine blocks
		err = network.WaitToMineNBlocks(10, 60, false)
		require.NoError(t, err, "Network should be mining new blocks now, but it's not")
	}
}

// TestKeepResettingRandomTwoNodes, with multiple rounds, it keeps resetting two random selected nodes at each round,
// then recover them at that round, after each round the network should keep mining.
func TestKeepResettingRandomTwoNodes(t *testing.T) {
	numOfNodes := 6
	network, err := e2e.NewNetwork(t, numOfNodes, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	defer network.Shutdown()
	// wait for the consensus engine to work.
	network.WaitToMineNBlocks(10, 60, false)

	rounds := 10
	for r := 0; r < rounds; r++ {
		// random select two nodes, and stop them.
		nodes := generateDistinctRandomNumbers(0, numOfNodes-1, 2)
		for _, n := range nodes {
			err = network[n].Close(false)
			network[n].Wait()
			require.NoError(t, err)
		}

		// network should be up and continue to mine blocks
		err = network.WaitToMineNBlocks(10, 60, false)
		require.NoError(t, err, "Network should be mining new blocks now, but it's not")

		// recover the two faulty nodes.
		for _, n := range nodes {
			err = network[n].Start()
			require.NoError(t, err)
		}

		// network should be up and continue to mine blocks
		err = network.WaitToMineNBlocks(10, 60, false)
		require.NoError(t, err, "Network should be mining new blocks now, but it's not")
	}
}

// TestKeepResettingRandomThreeNodes, with multiple rounds, it keeps resetting three random selected nodes at each round,
// then recover them at that round, after each round the network should keep mining.
func TestKeepResettingRandomThreeNodes(t *testing.T) {
	numOfNodes := 6
	network, err := e2e.NewNetwork(t, numOfNodes, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	defer network.Shutdown()
	// wait for the consensus engine to work.
	network.WaitToMineNBlocks(10, 60, false)

	rounds := 10
	for r := 0; r < rounds; r++ {
		// random select three nodes, and stop them.
		nodes := generateDistinctRandomNumbers(0, numOfNodes-1, 3)
		for _, n := range nodes {
			err = network[n].Close(false)
			network[n].Wait()
			require.NoError(t, err)
		}

		// network should be up and continue to mine blocks
		err = network.WaitToMineNBlocks(10, 60, false)
		require.EqualError(t, err, "context deadline exceeded")

		// recover the three faulty nodes.
		for _, n := range nodes {
			err = network[n].Start()
			require.NoError(t, err)
		}

		// network should be up and continue to mine blocks
		err = network.WaitToMineNBlocks(10, 60, false)
		require.NoError(t, err, "Network should be mining new blocks now, but it's not")
	}
}

func generateDistinctRandomNumbers(min, max, count int) []int {
	randomNumbers := make([]int, count)
	usedNumbers := make(map[int]bool)

	for i := 0; i < count; i++ {
		randomNum := rand.Intn(max-min+1) + min
		for usedNumbers[randomNum] {
			randomNum = rand.Intn(max-min+1) + min
		}
		randomNumbers[i] = randomNum
		usedNumbers[randomNum] = true
	}

	return randomNumbers
}
