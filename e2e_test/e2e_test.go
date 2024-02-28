package e2e

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"text/tabwriter"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/accounts/abi/bind"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	ccore "github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/params"
)

// TODO: move node resetting(start&stop) tests from ./consensus/test to this new framework since the new framework is
//  simple and stable than the legacy one.

// This test checks that we can process transactions that transfer value from
// one participant to another.
func TestSendingValue(t *testing.T) {
	network, err := NewNetwork(t, 7, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	defer network.Shutdown()
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	err = network[0].SendAUTtracked(ctx, network[1].Address, 10)
	require.NoError(t, err)
	_ = network.WaitToMineNBlocks(20, 30, false)
}

func TestProtocolContractsDeployment(t *testing.T) {
	network, err := NewNetwork(t, 2, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	defer network.Shutdown()
	// Autonity Contract
	autonityContract, _ := autonity.NewAutonity(params.AutonityContractAddress, network[0].WsClient)
	autonityConfig, err := autonityContract.Config(nil)
	require.NoError(t, err)
	require.Equal(t, params.TestAutonityContractConfig.Operator, autonityConfig.Protocol.OperatorAccount)
	require.Equal(t, params.TestAutonityContractConfig.BlockPeriod, autonityConfig.Protocol.BlockPeriod.Uint64())
	require.Equal(t, params.TestAutonityContractConfig.EpochPeriod, autonityConfig.Protocol.EpochPeriod.Uint64())
	require.Equal(t, params.TestAutonityContractConfig.MaxCommitteeSize, autonityConfig.Protocol.CommitteeSize.Uint64())
	require.Equal(t, params.TestAutonityContractConfig.DelegationRate, autonityConfig.Policy.DelegationRate.Uint64())
	require.Equal(t, params.TestAutonityContractConfig.MinBaseFee, autonityConfig.Policy.MinBaseFee.Uint64())
	require.Equal(t, params.TestAutonityContractConfig.Treasury, autonityConfig.Policy.TreasuryAccount)
	require.Equal(t, params.TestAutonityContractConfig.TreasuryFee, autonityConfig.Policy.TreasuryFee.Uint64())
	require.Equal(t, params.TestAutonityContractConfig.UnbondingPeriod, autonityConfig.Policy.UnbondingPeriod.Uint64())
	require.Equal(t, params.UpgradeManagerContractAddress, autonityConfig.Contracts.UpgradeManagerContract)
	require.Equal(t, params.AccountabilityContractAddress, autonityConfig.Contracts.AccountabilityContract)
	require.Equal(t, params.ACUContractAddress, autonityConfig.Contracts.AcuContract)
	require.Equal(t, params.StabilizationContractAddress, autonityConfig.Contracts.StabilizationContract)
	require.Equal(t, params.OracleContractAddress, autonityConfig.Contracts.OracleContract)
	require.Equal(t, params.SupplyControlContractAddress, autonityConfig.Contracts.SupplyControlContract)
	// Accountability Contract
	accountabilityContract, _ := autonity.NewAccountability(params.AccountabilityContractAddress, network[0].WsClient)
	accountabilityConfig, err := accountabilityContract.Config(nil)
	require.NoError(t, err)
	require.Equal(t, params.DefaultAccountabilityConfig.SlashingRatePrecision, accountabilityConfig.SlashingRatePrecision.Uint64())
	require.Equal(t, params.DefaultAccountabilityConfig.HistoryFactor, accountabilityConfig.HistoryFactor.Uint64())
	require.Equal(t, params.DefaultAccountabilityConfig.CollusionFactor, accountabilityConfig.CollusionFactor.Uint64())
	require.Equal(t, params.DefaultAccountabilityConfig.JailFactor, accountabilityConfig.JailFactor.Uint64())
	require.Equal(t, params.DefaultAccountabilityConfig.BaseSlashingRateMid, accountabilityConfig.BaseSlashingRateMid.Uint64())
	require.Equal(t, params.DefaultAccountabilityConfig.BaseSlashingRateLow, accountabilityConfig.BaseSlashingRateLow.Uint64())
	require.Equal(t, params.DefaultAccountabilityConfig.InnocenceProofSubmissionWindow, accountabilityConfig.InnocenceProofSubmissionWindow.Uint64())
	// Oracle Contract -- todo
	// ACU Contract -- todo
	// Supply Control Contract -- todo
	// Stabilization Contract -- todo
	// Upgrade Manager Contract
	upgradeManagerContract, _ := autonity.NewUpgradeManager(params.UpgradeManagerContractAddress, network[0].WsClient)
	upgradeManagerAutonityAddress, err := upgradeManagerContract.Autonity(nil)
	require.NoError(t, err)
	require.Equal(t, params.AutonityContractAddress, upgradeManagerAutonityAddress)
	err = network.WaitToMineNBlocks(2, 15, false)
	require.NoError(t, err)
}

func TestProtocolContractCache(t *testing.T) {
	t.Run("If minimum base fee is updated, cached value is updated as well", func(t *testing.T) {
		network, err := NewNetwork(t, 2, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
		require.NoError(t, err)
		defer network.Shutdown()
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		initialMinBaseFee := new(big.Int).SetUint64(uint64(params.InitialBaseFee))
		require.Equal(t, initialMinBaseFee.Bytes(), network[0].Eth.BlockChain().MinBaseFee().Bytes())
		require.Equal(t, initialMinBaseFee.Bytes(), network[1].Eth.BlockChain().MinBaseFee().Bytes())

		// update min base fee
		updatedMinBaseFee, _ := new(big.Int).SetString("30000000000", 10)
		autonityContract, _ := autonity.NewAutonity(params.AutonityContractAddress, network[0].WsClient)
		transactOpts, _ := bind.NewKeyedTransactorWithChainID(network[0].Key, params.TestChainConfig.ChainID)
		tx, err := autonityContract.SetMinimumBaseFee(transactOpts, updatedMinBaseFee)
		require.NoError(t, err)
		err = network.AwaitTransactions(ctx, tx)
		require.NoError(t, err)

		// contract should be updated
		minBaseFee, err := autonityContract.GetMinimumBaseFee(new(bind.CallOpts))
		require.NoError(t, err)
		require.Equal(t, updatedMinBaseFee.Bytes(), minBaseFee.Bytes())

		// caches should be updated too
		require.Equal(t, updatedMinBaseFee.Bytes(), network[0].Eth.BlockChain().MinBaseFee().Bytes())
		require.Equal(t, updatedMinBaseFee.Bytes(), network[1].Eth.BlockChain().MinBaseFee().Bytes())
	})
}

// a node is verifying a proposal, but while he is verifying the finalized block is injected from p2p layer
func TestNodeAlreadyHasProposedBlock(t *testing.T) {
	vals, err := Validators(t, 2, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	// trick to get a handle to Proposer of node 0
	var proposer interfaces.Proposer
	vals[0].TendermintServices = &interfaces.Services{Proposer: func(c interfaces.Core) interfaces.Proposer {
		proposer = &core.Proposer{Core: c.(*core.Core)}
		return proposer
	}}

	// start the network
	network, err := NewNetworkFromValidators(t, vals, true)
	require.NoError(t, err)
	defer network.Shutdown()

	// mine a couple blocks
	err = network.WaitToMineNBlocks(3, 60, false)
	require.NoError(t, err)

	// stop consensus engine of node 0, this will halt the network
	node := network[0]
	err = node.Eth.Engine().Close()
	require.NoError(t, err)

	// get latest inserted block and generate proposal out of it
	block := node.Eth.BlockChain().CurrentBlock()
	proposal := message.NewPropose(0, block.NumberU64(), -1, block, func(hash common.Hash) ([]byte, common.Address) {
		out, _ := crypto.Sign(hash[:], node.Key)
		return out, common.Address{}
	}).MustVerify(func(address common.Address) *types.CommitteeMember {
		return &types.CommitteeMember{
			Address:     common.Address{},
			VotingPower: common.Big1,
		}
	})

	// handle the proposal
	err = proposer.HandleProposal(context.TODO(), proposal)
	require.True(t, errors.Is(err, constants.ErrAlreadyHaveBlock))
}

func TestStartingAndStoppingNodes(t *testing.T) {
	network, err := NewNetwork(t, 5, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	defer network.Shutdown()
	n := network[0]
	// Send a TX to see that the network is working
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = n.SendAUTtracked(ctx, network[1].Address, 10)
	require.NoError(t, err)
	// Stop a node
	err = network[1].Close(true)
	network[1].Wait()
	require.NoError(t, err)
	// Send a TX to see that the network is working
	ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = n.SendAUTtracked(ctx, network[1].Address, 10)
	require.NoError(t, err)

	// Stop a node
	err = network[2].Close(true)
	network[2].Wait()
	require.NoError(t, err)
	// We have now stopped more than F nodes, so we expect TX processing to time out.
	// Well wait 5 times the avgTransactionDuration before we assume the TX is not being processed.
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
	// Send a TX to see that the network is still working
	ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = n.SendAUTtracked(ctx, network[1].Address, 10)
	require.NoError(t, err)

	// Start the last stopped node
	err = network[1].Start()
	require.NoError(t, err)
	// Send a TX to see that the network is still working
	err = n.SendAUTtracked(context.Background(), network[1].Address, 10)
	require.NoError(t, err)
}

func NewBroadcasterWithCheck(currentHeight uint64, stopped bool) func(c interfaces.Core) interfaces.Broadcaster {
	return func(c interfaces.Core) interfaces.Broadcaster {
		return &broadcasterWithCheck{c.(*core.Core), currentHeight, stopped}
	}
}

type broadcasterWithCheck struct {
	*core.Core
	currentHeight uint64
	stopped       bool
}

func (s *broadcasterWithCheck) Broadcast(msg message.Msg) {
	logger := s.Logger().New("step", s.Step())
	logger.Debug("Broadcasting", "message", msg.String())

	if s.stopped && msg.H() <= s.currentHeight {
		panic("stopped node sent old consensus message once he came back up")
	}

	s.BroadcastAll(msg)
}

// Tests that a stopped node, once restarted, will sync up to chain head before sending consensus messages
func TestWaitForChainSyncAfterStop(t *testing.T) {
	network, err := NewNetwork(t, 5, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	defer network.Shutdown()

	err = network.WaitToMineNBlocks(10, 60, false)
	require.NoError(t, err)

	// Stop node 0, he will need to sync up once restarted
	err = network[0].Close(false)
	require.NoError(t, err)
	network[0].Wait()

	err = network.WaitToMineNBlocks(10, 60, false)
	require.NoError(t, err)

	// Stop node 1. This will make the network lose liveness and halt
	err = network[1].Close(false)
	require.NoError(t, err)
	network[1].Wait()

	// network should be stalled now
	err = network.WaitToMineNBlocks(1, 5, false)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expecting %q, instead got: %v ", context.DeadlineExceeded.Error(), err)
	}

	// restart node 1. He will rightfully send a consensus message for the consensus instance he stopped at (chainHeight+1)
	chainHeight := network[2].GetChainHeight()
	network[1].CustHandler = &interfaces.Services{Broadcaster: NewBroadcasterWithCheck(chainHeight, true)}
	err = network[1].Start()
	require.NoError(t, err)

	// give time for node 1 to come back up
	err = network.WaitToMineNBlocks(15, 60, false)
	require.NoError(t, err)

	// restart node 0. He should sync up and not send old consensus messages
	chainHeight = network[2].GetChainHeight()
	network[0].CustHandler = &interfaces.Services{Broadcaster: NewBroadcasterWithCheck(chainHeight, true)}
	err = network[0].Start()
	require.NoError(t, err)

	// give time for node 0 to come back up
	err = network.WaitToMineNBlocks(15, 60, false)
	require.NoError(t, err)

}

// Test details
// a.setup 7 validators with 100 voting power on each, and keep 7 committee seats as well.
// b.start the network with 1st 3 nodes only, the network should be on-hold since the online voting power is less than 2/3 of 7
// c.after the on-holding for a while, start the 4th node, then the network should start to produce blocks without any on-holding.
func TestTendermintQuorum(t *testing.T) {
	users, err := Validators(t, 6, "10e18,v,100,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	network, err := NewNetworkFromValidators(t, users, false)
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
	err = network.WaitToMineNBlocks(3, 60, false)
	require.Error(t, err, "Network is not supposed to be mining blocks at this point")
	// start 4th node
	err = network[3].Start()
	require.NoError(t, err)

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(3, 60, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}

// a. setup 6 validators with 100 voting power on each, and keep 6 committee seats as well.
// b. start all the 6 validators for a while, as the network is producing blocks for a while.
// c.stop 3 nodes one by one, then the network should on-hold when the online voting power is less than 2/3 of 6.
// d.after the on-holding for a while, recover the stopped nodes, then the network should start to produce blocks without any on-holding.
func TestTendermintQuorum2(t *testing.T) {
	users, err := Validators(t, 6, "10e18,v,100,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	// creates a network of 6 users and starts all the nodes in it
	network, err := NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(3, 60, false)
	require.NoError(t, err)
	defer network.Shutdown()
	// stop 3 nodes and verify if transaction processing is halted
	for i, n := range network {
		// stop last 3 nodes
		if i > 2 {
			err = n.Close(true)
			n.Wait()
			require.NoError(t, err)
		}
	}
	// check if network on hold
	err = network.WaitToMineNBlocks(3, 60, false)
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
	err = network.WaitToMineNBlocks(3, 60, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}

// a. setup 7 validators (A, B, C, D, E, F, G) with 100 voting power on each, and keep 7 committee seats as well.
// b. start the network for a while(the network is producing blocks for a while).
// c. Shut down nodes A and B, now the network keeps liveness with C, D, E, F online, TXs should be mined.
// d. After a while, shut down Node C and D, with only E and F online, the network should on-hold, no transaction should be mined.
// e. Then start up node A and B, they should be synchronized with node E and F, and network liveness should recovered, new TXs should be mined after the recover.
// f. Then start up node C and D, both C and D should get synchronized to the latest chain height.
// g. Then shut down node E and F, the network should still keep liveness, TXs are mined.
// h. Recover E and F, they should get synchronized finally.
func TestTendermintQuorum4(t *testing.T) {
	users, err := Validators(t, 7, "10e18,v,100,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	// creates a network of 7 users and starts all the nodes in it
	network, err := NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown()
	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(3, 60, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
	i := 0
	for i < 2 {
		// stop 1st and 2nd node
		err = network[i].Close(true)
		require.NoError(t, err)
		network[i].Wait()
		i++
	}
	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(3, 60, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
	// shutting down 3rd and 4th node
	for i < 4 {
		// stop 1st and 2nd node
		err = network[i].Close(true)
		require.NoError(t, err)
		network[i].Wait()
		i++
	}

	// network should be on hold
	err = network.WaitToMineNBlocks(3, 60, false)
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
	err = network.WaitToMineNBlocks(3, 60, false)
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
	err = network.WaitToMineNBlocks(3, 60, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")

	// bring down 5th and 6th node
	for i < 6 {
		// stop 1st and 2nd node
		err = network[i].Close(true)
		require.NoError(t, err)
		network[i].Wait()
		i++
	}

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(3, 60, false)
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
	err = network.WaitToMineNBlocks(3, 60, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}

// setup up a network of 12 nodes
// ensure the newtwork is running and blocks are getting mined
// start/stop nodes in parallel
// wait for network to go on hold
// restart all nodes and ensures that network resumes mining new blocks
func TestStartStopAllNodesInParallel(t *testing.T) {
	const nodeCount = 12
	users, err := Validators(t, nodeCount, "10e18,v,100,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	// creates a network of 6 users and starts all the nodes in it
	network, err := NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown()
	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(3, 60, false)
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
				e := network[i].Close(true)
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
	err = network.WaitToMineNBlocks(3, 60, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}

/*
// UNSUPPORTED ON NON UNIX DEV ENV
func updateRlimit() {
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		fmt.Println("Error Getting Rlimit ", err)
	}
	fmt.Println(rLimit)
	rLimit.Max = 65536
	rLimit.Cur = 65536
	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		fmt.Println("Error Setting Rlimit ", err)
	}
	err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		fmt.Println("Error Getting Rlimit ", err)
	}
	fmt.Println("Rlimit Final", rLimit)
}
*/

// TestLargeNetwork test internally +100 nodes setups. This will be moved later
// in its own cmd package.
func TestLargeNetwork(t *testing.T) {
	t.Skip("only on demand")
	//
	//------ Config section -------
	//

	// DefaultVerbosity will set the log levels for the main components: consensus, eth, blockchain..
	log.DefaultVerbosity = log.LvlError
	//Set the root logger level for everything else.
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlError, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))
	// Fast epoch to see changes in committee reflected fast
	params.TestAutonityContractConfig.EpochPeriod = 5
	// total peers to be deployed
	const peerCount = 150
	// peers above max committee are participants
	const maxCommittee = 50

	//
	//----End config -------------
	//

	validators, _ := Validators(t, peerCount, "10e18,v,1000,127.0.0.1:%s,%s,%s,%s")
	network, err := NewInMemoryNetwork(t, validators, true, func(genesis *ccore.Genesis) {
		genesis.Config.AutonityContractConfig.MaxCommitteeSize = maxCommittee
	})
	require.NoError(t, err)
	defer network.Shutdown()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	autonityContract, err := autonity.NewAutonity(params.AutonityContractAddress, network[0].WsClient)
	require.NoError(t, err)
	transactOpts, err := bind.NewKeyedTransactorWithChainID(network[0].Key, params.TestChainConfig.ChainID)
	require.NoError(t, err)

	getNetworkState := func() (height uint64, committee types.Committee) {
		for _, n := range network {
			nodeHeight := n.Eth.BlockChain().CurrentHeader().Number.Uint64()
			if nodeHeight > height {
				height = nodeHeight
				committee = n.Eth.BlockChain().CurrentHeader().Committee
			}
		}
		return
	}

	inCommittee := func(address common.Address, committee types.Committee) bool {
		for i := range committee {
			if committee[i].Address == address {
				return true
			}
		}
		return false
	}

	w := tabwriter.NewWriter(os.Stdout, 1, 2, 1, ' ', 0)
	tickCount := 0
	go func() {
		for range ticker.C {
			tickCount++
			fmt.Println("-------------", tickCount, "---------------")

			if tickCount%10 == 4 {
				newMaxCommittee := rand.Intn(50) + 1
				autonityContract.SetCommitteeSize(transactOpts, big.NewInt(int64(newMaxCommittee)))
				fmt.Println("new committee size:", newMaxCommittee)
			}
			maxHeight, committee := getNetworkState()
			fmt.Println("max height:", maxHeight)
			fmt.Println("go routines:", runtime.NumGoroutine())
			var peersBuf strings.Builder
			peersBuf.WriteString("n\t")
			var committeeBuf strings.Builder
			committeeBuf.WriteString("c\t")
			var execCountBuf strings.Builder
			execCountBuf.WriteString("en\t")
			var consCountBuf strings.Builder
			consCountBuf.WriteString("cn\t")
			var stateBuf strings.Builder
			stateBuf.WriteString("b\t")
			for i, node := range network {
				peersBuf.WriteString(strconv.Itoa(i) + "\t")
				execCountBuf.WriteString(strconv.Itoa(node.ExecutionServer().PeerCount()) + "\t")
				consCountBuf.WriteString(strconv.Itoa(node.ConsensusServer().PeerCount()) + "\t")
				stateBuf.WriteString(strconv.Itoa(int(node.Eth.BlockChain().CurrentHeader().Number.Uint64())) + "\t")
				if inCommittee(node.Address, committee) {
					committeeBuf.WriteString("X" + "\t")
				} else {
					committeeBuf.WriteString(" " + "\t")
				}
			}
			fmt.Fprintln(w, peersBuf.String())
			fmt.Fprintln(w, committeeBuf.String())
			fmt.Fprintln(w, execCountBuf.String())
			fmt.Fprintln(w, consCountBuf.String())
			fmt.Fprintln(w, stateBuf.String())
			w.Flush()
		}
	}()

	err = network.WaitToMineNBlocks(20, 30, false)
	require.NoError(t, err)
}

func TestLoad(t *testing.T) {
	t.Skip("only on demand")
	peerCount := 7
	targetTPS := 1000
	validators, _ := Validators(t, peerCount, "10e18,v,1000,127.0.0.1:%s,%s,%s,%s")
	network, err := NewInMemoryNetwork(t, validators, true)
	require.NoError(t, err)
	gasFeeCap := new(big.Int).SetUint64(network[0].EthConfig.Genesis.Config.AutonityContractConfig.MinBaseFee)
	//err = network[0].Eth.StartMining(1)
	require.NoError(t, err)
	fmt.Println(network[0].Eth.BlockChain().CurrentHeader().Number)
	closeCh := make(chan struct{})
	for j := 0; j < peerCount; j++ {
		go func(id int) {
			timer := time.NewTicker(time.Second)
			defer timer.Stop()
			nonce := 0
			for {
				select {
				case <-timer.C:
					for i := 0; i < targetTPS/peerCount; i++ {
						rawTx := types.NewTx(&types.DynamicFeeTx{
							Nonce:     uint64(nonce),
							GasTipCap: common.Big1,
							GasFeeCap: new(big.Int).Mul(gasFeeCap, common.Big2),
							Gas:       21000,
							To:        &common.Address{},
							Value:     common.Big1,
							Data:      nil,
						})
						signed, err := types.SignTx(rawTx, types.LatestSigner(network[0].EthConfig.Genesis.Config), network[id].Key)
						require.NoError(t, err)
						network[id].Eth.TxPool().AddLocal(signed)
						nonce++
					}
				case <-closeCh:
					return
				}
			}
		}(j)
	}
	err = network.WaitToMineNBlocks(1800, 2300, false)
	require.NoError(t, err)
	close(closeCh)
}
