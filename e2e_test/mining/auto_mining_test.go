package mining

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/autonity/bindings"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core"
	e2e "github.com/autonity/autonity/e2e_test"
)

func TestMiningStartAfterGenesisTime(t *testing.T) {
	// Keep genesis in the near future so the test doesn't take minutes.
	// We start the nodes after the genesis timestamp to avoid relying on the
	// node-launch "pre-genesis" scheduling behavior (which can be slow/flaky under -race).
	genesisStart := uint64(time.Now().Add(15 * time.Second).Unix())
	validators, _ := e2e.Validators(t, 4, "10e18,v,1,127.0.0.1:%s,%s,%s,%s")
	network, err := e2e.NewNetworkFromValidators(t, validators, false, func(genesis *core.Genesis) {
		genesis.Timestamp = genesisStart
	})
	require.NoError(t, err)
	defer network.Shutdown(t)

	// Wait until after genesis time, then start the nodes and wait for the first block.
	if d := time.Until(time.Unix(int64(genesisStart), 0)); d > 0 {
		time.Sleep(d + 1*time.Second)
	}
	for _, n := range network {
		require.NoError(t, n.Start())
	}
	require.NoError(t, network.WaitForHeight(1, 60))
}

// TestMiningManagementOfValidators, shrink and extend the committee size, and check the mining state for validator and
// non validator nodes.
func TestMiningManagementOfValidators(t *testing.T) {
	numOfNodes := 4
	network, err := e2e.NewNetwork(t, numOfNodes, "10e18,v,1,127.0.0.1:%s,%s,%s,%s")
	require.NoError(t, err)
	defer network.Shutdown(t)
	// wait for the consensus engine to work.
	require.NoError(t, network.WaitToMineNBlocks(2, 20, false))

	// all validators should be mining.
	for i := 0; i < numOfNodes; i++ {
		i := i
		require.Eventually(t, func() bool {
			return network[i].Eth.IsMining()
		}, 20*time.Second, 250*time.Millisecond)
	}

	client := network[0]
	optKey := client.Key

	// shrink committee size to less than numOfNodes, some validators shouldn't be mining if they
	// are no longer in the committee.
	newSize := new(big.Int).SetUint64(uint64(numOfNodes - 1))
	tm := 30 * time.Second
	err = client.AwaitSetCommitteeSize(optKey, newSize, tm)
	require.NoError(t, err)

	// wait for epoch rotation
	epochPeriod := client.EthConfig.Genesis.Config.AutonityContractConfig.EpochPeriod
	before := client.Eth.BlockChain().CurrentHeader().Number.Uint64()
	require.Eventually(t, func() bool {
		h := client.Eth.BlockChain().CurrentHeader().Number.Uint64()
		return h > before && h%epochPeriod == 0
	}, time.Duration(epochPeriod*2+30)*time.Second, 500*time.Millisecond)

	// wait for a while to let the nodes get synced with epoch rotation
	require.NoError(t, network.WaitToMineNBlocks(3, 30, false))

	// get new committee, and check the new size.
	var shrunkCommittee []bindings.IAutonityCommitteeMember
	require.Eventually(t, func() bool {
		var err error
		shrunkCommittee, err = client.Interactor.Call(nil).GetCommittee()
		if err != nil {
			return false
		}
		return uint64(len(shrunkCommittee)) == newSize.Uint64()
	}, time.Duration(epochPeriod*2+30)*time.Second, 500*time.Millisecond)
	shrunkCommitteeMap := make(map[common.Address]struct{})
	for _, c := range shrunkCommittee {
		shrunkCommitteeMap[c.Addr] = struct{}{}
	}

	// check mining state after committee size shrink.
	for i := 0; i < numOfNodes; i++ {
		isMining := false
		if _, ok := shrunkCommitteeMap[network[i].Address]; ok {
			isMining = true
		}
		mining := network[i].Eth.IsMining()
		require.Equal(t, isMining, mining)
	}

	// now extend the committee size, after to epoch rotation, new validator should start ming again.
	err = client.AwaitSetCommitteeSize(optKey, new(big.Int).SetUint64(uint64(numOfNodes)), tm)
	require.NoError(t, err)

	// wait for epoch rotation
	before = client.Eth.BlockChain().CurrentHeader().Number.Uint64()
	require.Eventually(t, func() bool {
		h := client.Eth.BlockChain().CurrentHeader().Number.Uint64()
		return h > before && h%epochPeriod == 0
	}, time.Duration(epochPeriod*2+30)*time.Second, 500*time.Millisecond)

	// wait for a while to let the nodes get synced with epoch rotation
	require.NoError(t, network.WaitToMineNBlocks(3, 30, false))

	// get new committee, and check the new size.
	var extendedCommittee []bindings.IAutonityCommitteeMember
	require.Eventually(t, func() bool {
		var err error
		extendedCommittee, err = client.Interactor.Call(nil).GetCommittee()
		if err != nil {
			return false
		}
		return len(extendedCommittee) == numOfNodes
	}, time.Duration(epochPeriod*3+30)*time.Second, 500*time.Millisecond)

	// wait for a while to get the new validator's mining worker be started
	require.NoError(t, network.WaitToMineNBlocks(3, 30, false))

	// all validators should be mining.
	for i := 0; i < numOfNodes; i++ {
		i := i
		require.Eventually(t, func() bool {
			return network[i].Eth.IsMining()
		}, 20*time.Second, 250*time.Millisecond)
	}
}
