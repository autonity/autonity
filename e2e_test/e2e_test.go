package test

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/autonity/autonity/accounts/abi/bind"
	"github.com/autonity/autonity/common"

	"github.com/autonity/autonity/test"
	"github.com/stretchr/testify/require"
)

// TODO: move node resetting(start&stop) tests from ./consensus/test to this new framework since the new framework is
//  simple and stable than the legacy one.

// This test checks that we can process transactions that transfer value from
// one participant to another.
func TestSendingValue(t *testing.T) {
	network, err := test.NewNetwork(t, 2, "10e18,v,1,0.0.0.0:%s,%s")
	require.NoError(t, err)
	defer network.Shutdown()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = network[0].SendAUTtracked(ctx, network[1].Address, 10)
	require.NoError(t, err)
}

// This test checks that when a transaction is processed the fees are divided
// between validators and stakeholders.
func TestFeeRedistributionValidatorsAndDelegators(t *testing.T) {
	vals, err := test.Validators(t, 3, "10e18,v,10000,0.0.0.0:%s,%s")
	require.NoError(t, err)

	vals[2].Stake = 25000

	network, err := test.NewNetworkFromValidators(t, vals, true)
	require.NoError(t, err)
	defer network.Shutdown()

	n := network[0]

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	// retrieve current balance
	// send liquid newton to some random address
	// check balance - shouldnt have increased
	// wait for epoch
	// check claimable fees
	// redeem fees

	// Setup Bindings
	autonityContract, err := NewAutonity(autonityContractAddr, n.WsClient)
	require.NoError(t, err)
	valAddrs, err := autonityContract.GetValidators(nil)
	require.NoError(t, err)
	liquidContracts := make([]*Liquid, len(valAddrs))
	validators := make([]AutonityValidator, len(valAddrs))
	for i, valAddr := range valAddrs {
		validators[i], err = autonityContract.GetValidator(nil, valAddr)
		require.NoError(t, err)
		liquidContracts[i], err = NewLiquid(validators[i].LiquidContract, n.WsClient)
		require.NoError(t, err)
	}
	transactor, _ := bind.NewKeyedTransactorWithChainID(vals[0].TreasuryKey, big.NewInt(1234))
	tx, err := liquidContracts[0].Transfer(
		transactor,
		common.Address{66, 66}, big.NewInt(1337))

	require.NoError(t, err)
	_ = network.WaitToMineNBlocks(2, 20)
	tx2, err := n.SendAUT(ctx, network[1].Address, 10)
	require.NoError(t, err)
	err = network.AwaitTransactions(ctx, tx, tx2)
	require.NoError(t, err)
	// claimable fees should be 0 before epoch
	for i := range liquidContracts {
		unclaimed, _ := liquidContracts[i].UnclaimedRewards(&bind.CallOpts{}, validators[i].Treasury)
		require.Equal(t, big.NewInt(0).Bytes(), unclaimed.Bytes())
	}

	// wait for epoch

	// calculate reward pool
	r1, _ := n.WsClient.TransactionReceipt(context.Background(), tx.Hash())
	r2, _ := n.WsClient.TransactionReceipt(context.Background(), tx2.Hash())

	b1, _ := n.WsClient.BlockByNumber(context.Background(), r1.BlockNumber)
	b2, _ := n.WsClient.BlockByNumber(context.Background(), r2.BlockNumber)

	rewardT1 := new(big.Int).Mul(new(big.Int).SetUint64(r1.CumulativeGasUsed), b1.BaseFee())
	rewardT2 := new(big.Int).Mul(new(big.Int).SetUint64(r2.CumulativeGasUsed), b2.BaseFee())

	totalFees := new(big.Int).Add(rewardT1, rewardT2)
	treasuryRewards := new(big.Int).Div(new(big.Int).Mul(totalFees, new(big.Int).SetUint64(15)), big.NewInt(10000))
	totalRewards := new(big.Int).Sub(totalFees, treasuryRewards)

	balanceBeforeEpoch, _ := n.WsClient.BalanceAt(context.Background(), autonityContractAddr, nil)
	require.Equal(t, totalFees, balanceBeforeEpoch)

	err = network.WaitToMineNBlocks(30, 90)
	require.NoError(t, err)

	fmt.Println("total rewards", totalRewards)
	balanceGlobalTreasury, _ := n.WsClient.BalanceAt(context.Background(), common.Address{120}, nil)
	cfg, _ := autonityContract.Config(nil)
	fmt.Println(cfg)
	require.Equal(t, treasuryRewards, balanceGlobalTreasury)

	stake := []int64{10000 - 1337, 10000, 25000}
	epochStake := []int64{10000, 10000, 25000}
	totalStake := int64(45000)
	for i := range liquidContracts {
		unclaimed, _ := liquidContracts[i].UnclaimedRewards(&bind.CallOpts{}, validators[i].Treasury)
		totalValRewards := new(big.Int).Div(new(big.Int).Mul(totalRewards, big.NewInt(epochStake[i])), big.NewInt(totalStake))
		valCommission := new(big.Int).Div(new(big.Int).Mul(totalValRewards, big.NewInt(12)), big.NewInt(100))
		stakerReward := new(big.Int).Sub(totalValRewards, valCommission)
		require.Equal(t, new(big.Int).Div(new(big.Int).Mul(stakerReward, big.NewInt(stake[i])), big.NewInt(epochStake[i])), unclaimed)
	}

}

func TestStartingAndStoppingNodes(t *testing.T) {
	network, err := test.NewNetwork(t, 5, "10e18,v,1,0.0.0.0:%s,%s")
	require.NoError(t, err)
	defer network.Shutdown()
	n := network[0]

	// Send a tx to see that the network is working
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = n.SendAUTtracked(ctx, network[1].Address, 10)
	require.NoError(t, err)

	// Stop a node
	err = network[1].Close()
	network[1].Wait()
	require.NoError(t, err)

	// Send a tx to see that the network is working
	ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = n.SendAUTtracked(ctx, network[1].Address, 10)
	require.NoError(t, err)

	// Stop a node
	err = network[2].Close()
	network[2].Wait()
	require.NoError(t, err)

	// We have now stopped more than F nodes, so we expect tx processing to time out.
	// Well wait 5 times the avgTransactionDuration before we assume the tx is not being processed.
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err = n.SendAUTtracked(ctx, network[1].Address, 10)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expecting %q, instead got: %v ", context.DeadlineExceeded.Error(), err)
	}

	// We start a node again and expect the previously unprocessed transaction to be processed
	err = network[2].Start()
	require.NoError(t, err)

	// Ensure that the previously sent transaction is now processed
	ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = n.AwaitSentTransactions(ctx)
	require.NoError(t, err)

	// Send a tx to see that the network is still working
	ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = n.SendAUTtracked(ctx, network[1].Address, 10)
	require.NoError(t, err)

	// Start the last stopped node
	err = network[1].Start()
	require.NoError(t, err)

	// Send a tx to see that the network is still working
	err = n.SendAUTtracked(context.Background(), network[1].Address, 10)
	require.NoError(t, err)
}

// Test details
//a.setup 7 validators with 100 voting power on each, and keep 7 committee seats as well.
//b.start the network with 1st 3 nodes only, the network should be on-hold since the online voting power is less than 2/3 of 7
//c.after the on-holding for a while, start the 4th node, then the network should start to produce blocks without any on-holding.
func TestTendermintQuorum(t *testing.T) {
	users, err := test.Validators(t, 6, "10e18,v,100,0.0.0.0:%s,%s")
	require.NoError(t, err)
	network, err := test.NewNetworkFromValidators(t, users, false)
	require.NoError(t, err)
	defer network.Shutdown()
	for i, n := range network {
		// start 3 nodes
		if i < 3 {
			err = n.Start()
			require.NoError(t, err)
		}
	}
	// check if network on hold
	err = network.WaitToMineNBlocks(3, 60)
	require.Error(t, err, "Network is not supposed to be mining blocks at this point")
	// start 4th node
	err = network[3].Start()
	require.NoError(t, err)

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(3, 60)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}

//a. setup 6 validators with 100 voting power on each, and keep 6 committee seats as well.
//b. start all the 6 validators for a while, as the network is producing blocks for a while.
//c.stop 3 nodes one by one, then the network should on-hold when the online voting power is less than 2/3 of 6.
//d.after the on-holding for a while, recover the stopped nodes, then the network should start to produce blocks without any on-holding.
func TestTendermintQuorum2(t *testing.T) {
	users, err := test.Validators(t, 6, "10e18,v,100,0.0.0.0:%s,%s")
	require.NoError(t, err)
	// creates a network of 6 users and starts all the nodes in it
	network, err := test.NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(3, 60)
	require.NoError(t, err)
	defer network.Shutdown()
	// stop 3 nodes and verify if transaction processing is halted
	for i, n := range network {
		// stop last 3 nodes
		if i > 2 {
			err = n.Close()
			n.Wait()
			require.NoError(t, err)
		}
	}
	// check if network on hold
	err = network.WaitToMineNBlocks(3, 60)
	require.Error(t, err, "Network is not supposed to be mining blocks at this point")

	// start the nodes back again
	for i, n := range network {
		// start back 4,5 & 6th node
		if i > 2 {
			err = n.Start()
			require.NoError(t, err)

		}
	}
	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(3, 60)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}

//a. setup 7 validators (A, B, C, D, E, F, G) with 100 voting power on each, and keep 7 committee seats as well.
//b. start the network for a while(the network is producing blocks for a while).
//c. Shut down nodes A and B, now the network keeps liveness with C, D, E, F online, TXs should be mined.
//d. After a while, shut down Node C and D, with only E and F online, the network should on-hold, no transaction should be mined.
//e. Then start up node A and B, they should be synchronized with node E and F, and network liveness should recovered, new TXs should be mined after the recover.
//f. Then start up node C and D, both C and D should get synchronized to the latest chain height.
//g. Then shut down node E and F, the network should still keep liveness, TXs are mined.
//h. Recover E and F, they should get synchronized finally.
func TestTendermintQuorum4(t *testing.T) {
	users, err := test.Validators(t, 7, "10e18,v,100,0.0.0.0:%s,%s")
	require.NoError(t, err)
	// creates a network of 7 users and starts all the nodes in it
	network, err := test.NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown()
	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(3, 60)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
	i := 0
	for i < 2 {
		// stop 1st and 2nd node
		err = network[i].Close()
		require.NoError(t, err)
		network[i].Wait()
		i++
	}
	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(3, 60)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
	// shutting down 3rd and 4th node
	for i < 4 {
		// stop 1st and 2nd node
		err = network[i].Close()
		require.NoError(t, err)
		network[i].Wait()
		i++
	}

	// network should be on hold
	err = network.WaitToMineNBlocks(3, 60)
	require.Error(t, err, "Network is not supposed to be mining blocks at this point")

	// start back 1st and 2nd node
	i = 0
	for i < 2 {
		err = network[i].Start()
		require.NoError(t, err)
		i++
	}

	err = network.WaitForSyncComplete()
	require.NoError(t, err)

	// we are restoring liveliness we would wait for
	// network should be back up and continue to mine blocks
	err = network.WaitToMineNBlocks(3, 60)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")

	// start back 3rd and 4th node
	for i < 4 {
		err = network[i].Start()
		require.NoError(t, err)
		i++
	}

	// wait for sync completion
	err = network.WaitForSyncComplete()
	require.NoError(t, err)

	// network should be back up and continue to mine blocks
	err = network.WaitToMineNBlocks(3, 60)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")

	// bring down 5th and 6th node
	for i < 6 {
		// stop 1st and 2nd node
		err = network[i].Close()
		require.NoError(t, err)
		network[i].Wait()
		i++
	}

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(3, 60)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
	i = 4
	// start back 5th and 6th node
	for i < 4 {
		err = network[i].Start()
		require.NoError(t, err)
		i++
	}

	// wait for sync completion
	err = network.WaitForSyncComplete()
	require.NoError(t, err)

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(3, 60)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}

// setup up a network of 12 nodes
// ensure the newtwork is running and blocks are getting mined
// start/stop nodes in parallel
// wait for network to go on hold
// restart all nodes and ensures that network resumes mining new blocks
func TestStartStopAllNodesInParallel(t *testing.T) {
	const nodeCount = 12
	users, err := test.Validators(t, nodeCount, "10e18,v,100,0.0.0.0:%s,%s")
	require.NoError(t, err)
	// creates a network of 6 users and starts all the nodes in it
	network, err := test.NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown()
	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(3, 60)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")

	type nodeStatus struct {
		lock   sync.Mutex
		status bool
	}
	var wg sync.WaitGroup
	mlock := sync.RWMutex{}
	m := make(map[int]*nodeStatus, nodeCount) // map to maintain start status of all nodes
	for i := 0; i < nodeCount; i++ {
		m[i] = &nodeStatus{sync.Mutex{}, true}
	}

	// randomly start stop nodes
	for j := 0; j < 40; j++ {
		i := rand.Intn(nodeCount) //nolint:gosec
		switch i % 4 {
		case 0, 1, 2:
			wg.Add(1)
			go func() {
				defer wg.Done()
				mlock.RLock()
				nodeStatus := m[i]
				mlock.RUnlock()
				nodeStatus.lock.Lock()
				defer nodeStatus.lock.Unlock()
				if !nodeStatus.status {
					return
				}
				e := network[i].Close()
				require.NoError(t, e)
				network[i].Wait()
				nodeStatus.status = false
				mlock.Lock()
				m[i] = nodeStatus
				mlock.Unlock()
			}()
		case 3:
			wg.Add(1)
			go func() {
				defer wg.Done()
				mlock.RLock()
				nodeStatus := m[i]
				mlock.RUnlock()
				nodeStatus.lock.Lock()
				defer nodeStatus.lock.Unlock()
				if nodeStatus.status {
					return
				}
				e := network[i].Start()
				require.NoError(t, e)
				nodeStatus.status = true
				mlock.Lock()
				m[i] = nodeStatus
				mlock.Unlock()
			}()
		}
	}
	// waiting for all go routines to be over
	wg.Wait()
	// start back all nodes
	for _, n := range network {
		err = n.Start()
		require.NoError(t, err)
	}

	// Verify network is not on hold anymore and producing blocks
	err = network.WaitToMineNBlocks(3, 60)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}
