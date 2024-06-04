package mining

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core"
	e2e "github.com/autonity/autonity/e2e_test"
)

func TestMiningStartAfterGenesisTime(t *testing.T) {
	delay := 2 * 60
	genesisStart := uint64(time.Now().Unix()) + uint64(delay)
	validators, _ := e2e.Validators(t, 4, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	network, err := e2e.NewNetworkFromValidators(t, validators, true, func(genesis *core.Genesis) {
		genesis.Timestamp = genesisStart
	})
	require.NoError(t, err)
	ticker := time.NewTicker(500 * time.Millisecond)
	for range ticker.C {
		if uint64(time.Now().Unix()) < genesisStart {
			continue
		}
		if uint64(time.Now().Unix()) > genesisStart+1 {
			if network[0].Eth.BlockChain().CurrentHeader().Number.Uint64() == 0 {
				t.Error("Genesis launch failure")
			}
			break
		}
	}
	ticker.Stop()
}

// TestMiningManagementOfValidators, shrink and extend the committee size, and check the mining state for validator and
// non validator nodes.
func TestMiningManagementOfValidators(t *testing.T) {
	numOfNodes := 4
	network, err := e2e.NewNetwork(t, numOfNodes, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	defer network.Shutdown(t)
	// wait for the consensus engine to work.
	network.WaitToMineNBlocks(2, 10, false)

	// all validators should be mining.
	for i := 0; i < numOfNodes; i++ {
		isMining, err := network[i].WsClient.IsMining(context.Background())
		require.NoError(t, err)
		require.True(t, isMining)
	}

	client := network[0]
	optKey := client.Key

	// shrink committee size to less than numOfNodes, some validators shouldn't be mining if they
	// are no longer in the committee.
	newSize := new(big.Int).SetUint64(uint64(numOfNodes - 1))
	tm := 5 * time.Second
	err = client.AwaitSetCommitteeSize(optKey, newSize, tm)
	require.NoError(t, err)

	// wait for epoch rotation
	epochPeriod := client.EthConfig.Genesis.Config.AutonityContractConfig.EpochPeriod
	for {
		if client.Eth.BlockChain().CurrentHeader().Number.Uint64()%epochPeriod == 0 {
			break
		}
		time.Sleep(time.Second)
	}

	// get new committee, and check the new size.
	shrunkCommittee, err := client.Interactor.Call(nil).GetCommittee()
	require.NoError(t, err)
	require.Equal(t, newSize.Uint64(), uint64(len(shrunkCommittee)))
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
		mining, err := network[i].WsClient.IsMining(context.Background())
		require.NoError(t, err)
		require.Equal(t, isMining, mining)
	}

	// now extend the committee size, after to epoch rotation, new validator should start ming again.
	err = client.AwaitSetCommitteeSize(optKey, new(big.Int).SetUint64(uint64(numOfNodes)), tm)
	require.NoError(t, err)

	// wait for epoch rotation
	for {
		if client.Eth.BlockChain().CurrentHeader().Number.Uint64()%epochPeriod == 0 {
			break
		}
		time.Sleep(time.Second)
	}

	// get new committee, and check the new size.
	extendedCommittee, err := client.Interactor.Call(nil).GetCommittee()
	require.NoError(t, err)
	require.Equal(t, numOfNodes, len(extendedCommittee))

	// wait for a while to get the new validator's mining worker be started
	network.WaitToMineNBlocks(10, 10, false)

	// all validators should be mining.
	for i := 0; i < numOfNodes; i++ {
		isMining, err := network[i].WsClient.IsMining(context.Background())
		require.NoError(t, err)
		require.True(t, isMining)
	}
}
