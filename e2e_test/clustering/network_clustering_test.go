package clustering

import (
	"context"
	"math/big"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/autonity/autonity/autonity/bindings"
	ccore "github.com/autonity/autonity/core"
	crypto2 "github.com/autonity/autonity/crypto"
	enodes "github.com/autonity/autonity/p2p/enode"
	"github.com/autonity/autonity/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	routerInterfaces "github.com/autonity/autonity/consensus/tendermint/router/interfaces"
	"github.com/autonity/autonity/core/types"
	e2e "github.com/autonity/autonity/e2e_test"
)

//todo: test to verify clustering view for all validators is same

// TestClusteringHappyCase is a happy case to test 5 clusters with each of them contains 5 nodes. The latency measurement
// is a base on an local simulator which generates [0, 500) ms delays.
func TestClusteringHappyCase(t *testing.T) {
	// mocked service with a local ping simulator which generates [0, 500) ms latency.
	mockedService := &interfaces.Services{Pinger: newSimulatedPinger}

	validators, err := e2e.Validators(t, 10, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	for _, validator := range validators {
		validator.TendermintServices = mockedService
	}

	// creates a network of 25 validators and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, validators, true)
	require.NoError(t, err)
	defer network.Shutdown(t)

	// runs for about 10 epoches period.
	network.WaitToMineNBlocks(200, 250, false)
}

// TestClusteringResetAllNodes, it stops all nodes one by one, and start them again one by one. The network should recover to
// mining.
func TestClusteringResetAllNodes(t *testing.T) {
	numOfNodes := 36
	mockedService := &interfaces.Services{Pinger: newSimulatedPinger}

	validators, err := e2e.Validators(t, numOfNodes, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	for _, validator := range validators {
		validator.TendermintServices = mockedService
	}

	// creates a network of 25 validators and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, validators, true)
	require.NoError(t, err)
	defer network.Shutdown(t)

	// wait for the consensus engine to work.
	network.WaitToMineNBlocks(10, 10, false)

	// reset all nodes concurrently
	for _, n := range network {
		go resetNode(t, n)
	}

	// network should be up and continue to mine blocks
	network.WaitToMineNBlocks(300, 300, false)
}

// TestClusteringResetFNodes, it stops random selected F nodes one by one, and observe if the net is still mining, then it recover
// F nodes one by one, the network should keep mining all the time.
func TestClusteringResetFNodes(t *testing.T) {
	numOfNodes := 36
	mockedService := &interfaces.Services{Pinger: newSimulatedPinger}

	validators, err := e2e.Validators(t, numOfNodes, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	for _, validator := range validators {
		validator.TendermintServices = mockedService
	}
	// creates a network of 25 validators and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, validators, true)
	require.NoError(t, err)
	defer network.Shutdown(t)

	// wait for the consensus engine to work.
	network.WaitToMineNBlocks(60, 60, false)

	// stop random selected F nodes.
	f := bft.F(new(big.Int).SetInt64(int64(numOfNodes))).Int64()
	fNodes := make(map[int]struct{})
	for i := int64(0); i < f; {
		selectedID := rand.Intn(len(network))
		if _, ok := fNodes[selectedID]; ok {
			continue
		}
		i++

		go stopNode(t, network, selectedID)
		fNodes[selectedID] = struct{}{}
	}

	// network should be up and continue to mine blocks
	network.WaitToMineNBlocks(300, 300, false)

	// recover nodes
	for id := range fNodes {
		go startNode(t, network, id)
	}
	// network should be up and continue to mine blocks
	network.WaitToMineNBlocks(10, 10, false)
}

// NoRelayingSelector is used for not to relay proposal in the network for Faulty nodes.
type NoRelayingSelector struct{}

var newNoRelayingSelector = func(_ interfaces.Router) routerInterfaces.PeerSelector { return &NoRelayingSelector{} }

func (r *NoRelayingSelector) SetBroadcaster(_ routerInterfaces.PeerFinder) {
}

func committeeAddresses(committee *types.Committee) []common.Address {
	addresses := make([]common.Address, 0)
	if committee == nil || len(committee.Members) == 0 {
		return []common.Address{}
	}
	for _, member := range committee.Members {
		addresses = append(addresses, member.Address)
	}
	return addresses
}

func (r *NoRelayingSelector) SelectPeers(committee *types.Committee, msg message.Msg, from common.Address) ([]common.Address, error) {
	// if not part of the committee return
	if member := committee.MemberByAddress(from); member == nil {
		return []common.Address{}, nil
	}

	if msg.Code() != message.ProposalCode {
		return committeeAddresses(committee), nil
	}

	return nil, nil
}

func TestFFaultyRelayers(t *testing.T) {
	numOfNodes := 36
	pinger := newSimulatedPinger

	validators, err := e2e.Validators(t, numOfNodes, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	// Create the network with F num of nodes which does not relay proposal.
	f := bft.F(new(big.Int).SetInt64(int64(numOfNodes))).Int64()
	for i, validator := range validators {
		mockedService := &interfaces.Services{Pinger: pinger}
		if int64(i) < f {
			mockedService.Selector = newNoRelayingSelector
		}
		validator.TendermintServices = mockedService
	}

	network, err := e2e.NewNetworkFromValidators(t, validators, true)
	require.NoError(t, err)
	defer network.Shutdown(t)

	// wait for the consensus engine to work.
	network.WaitToMineNBlocks(200, 250, false)
}

func Test2FFaultyRelayers(t *testing.T) {
	numOfNodes := 36
	pinger := newSimulatedPinger

	validators, err := e2e.Validators(t, numOfNodes, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	// Create the network with F num of nodes which does not relay proposal.
	f := bft.F(new(big.Int).SetInt64(int64(numOfNodes))).Int64()
	overF := f * 2
	for i, validator := range validators {
		mockedService := &interfaces.Services{Pinger: pinger}
		if int64(i) < overF {
			mockedService.Selector = newNoRelayingSelector
		}
		validator.TendermintServices = mockedService
	}

	network, err := e2e.NewNetworkFromValidators(t, validators, true)
	require.NoError(t, err)
	defer network.Shutdown(t)

	// wait for the consensus engine to work.
	network.WaitToMineNBlocks(250, 200, false)
}

func Test3FFaultyRelayers(t *testing.T) {
	numOfNodes := 36
	pinger := newSimulatedPinger

	validators, err := e2e.Validators(t, numOfNodes, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	// Create the network with F num of nodes which does not relay proposal.
	f := bft.F(new(big.Int).SetInt64(int64(numOfNodes))).Int64()
	overF := f * 3
	for i, validator := range validators {
		mockedService := &interfaces.Services{Pinger: pinger}
		if int64(i) < overF {
			mockedService.Selector = newNoRelayingSelector
		}
		validator.TendermintServices = mockedService
	}

	network, err := e2e.NewNetworkFromValidators(t, validators, true)
	require.NoError(t, err)
	defer network.Shutdown(t)

	// wait for the consensus engine to work.
	network.WaitToMineNBlocks(60, 60, false)
}

func TestCrossingClusteringThreshold(t *testing.T) {
	t.Run("from below to above", func(t *testing.T) {
		vals, err := e2e.Validators(t, 3, "10e18,v,100,127.0.0.1:%s,%s,%s,%s")
		require.NoError(t, err)
		customClusteringThreshold := 10
		network, err := e2e.NewNetworkFromValidators(t, vals, true, func(genesis *ccore.Genesis) {
			genesis.Config.AutonityContractConfig.ClusteringThreshold = uint64(customClusteringThreshold)
			// make the epoch a little bit longer to allow for registering all the new vals and for them to sync properly
			genesis.Config.AutonityContractConfig.EpochPeriod = uint64(150)
			// reduce lookback period for omission so that it is more reactive
			genesis.Config.OmissionAccountabilityConfig.LookbackWindow = 5
			// set oracle vote period low so that we can lower the epoch period later
			genesis.Config.OracleContractConfig.VotePeriod = uint64(5)
		})
		require.NoError(t, err)
		defer network.Shutdown(t)

		require.NoError(t, network.WaitForNetworkToStartMining())

		// node 0 is set as the operator
		operatorKey := network[0].Key
		genesis := network[0].EthConfig.Genesis
		originalValidators := []common.Address{network[0].Address, network[1].Address, network[2].Address}
		var originalEnodes []*enodes.Node
		for i, node := range network {
			originalEnodes = append(originalEnodes, enodes.NewV4(&node.Key.PublicKey, vals[i].NodeIP, vals[i].NodePort, vals[i].NodePort))
		}

		// fetch voting power of initial nodes
		autonityContract, err := bindings.NewAutonity(params.AutonityContractAddress, network[0].WsClient)
		require.NoError(t, err)
		originalCommittee, err := autonityContract.GetCommittee(nil)
		require.NoError(t, err)
		additionalNodesVotingPower := originalCommittee[0].VotingPower.Uint64() - 1 // keep original vals in the first spots of the committee

		// register 16 additional validators to go above clustering threshold
		n := 16
		validators, err := e2e.Validators(t, n, "0,v,100,127.0.0.1:%s,%s,%s,%s")
		require.NoError(t, err)
		nodes := make([]*e2e.Node, 0, n)
		for i := 0; i < n; i++ {
			node, err := e2e.NewValidatorNode(validators[i], genesis, i+3, false)
			require.NoError(t, err)

			node.Config.ExecutionP2P.BootstrapNodes = originalEnodes

			err = node.Start()
			require.NoError(t, err)
			defer node.Close(true)
			nodes = append(nodes, node)

			treasuryAddress := crypto2.PubkeyToAddress(node.TreasuryKey.PublicKey)
			enode := enodes.NewV4(&node.Key.PublicKey, validators[i].NodeIP, validators[i].NodePort, validators[i].NodePort)
			ens := enodes.AppendConsensusEndpoint(validators[i].AcnIP.String(), strconv.Itoa(validators[i].AcnPort), enode.String())
			oracleAddress := crypto2.PubkeyToAddress(validators[i].OracleKey.PublicKey)
			nodeAddress := crypto2.PubkeyToAddress(validators[i].NodeKey.PublicKey)

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			err = network[2].SendAUTtracked(ctx, treasuryAddress, 1e17)
			cancel()
			require.NoError(t, err)

			proof, err := crypto2.AutonityPOPProof(node.Key, validators[i].OracleKey, treasuryAddress.Hex(), node.ConsensusKey)
			require.NoError(t, err)
			err = network[0].AwaitRegisterValidator(node.TreasuryKey, ens, oracleAddress, node.ConsensusKey.PublicKey().Marshal(), proof, 3*time.Second)
			require.NoError(t, err)

			err = network[0].AwaitMintNTN(operatorKey, treasuryAddress, new(big.Int).SetUint64(params.Ether), 3*time.Second)
			require.NoError(t, err)
			err = network[0].AwaitBondStake(node.TreasuryKey, nodeAddress, new(big.Int).SetUint64(additionalNodesVotingPower), 3*time.Second)
			require.NoError(t, err)
		}

		// we should still be in epoch 0
		epochID, err := autonityContract.GetEpochID(nil)
		require.NoError(t, err)
		epochPeriod, err := autonityContract.GetEpochPeriod(nil)
		require.NoError(t, err)

		require.Equal(t, uint64(0), epochID.Uint64())

		// check if last started nodes is ready for consensus
		lastNode := nodes[len(nodes)-1]
		currentChainHeight := network[0].GetChainHeight()
		time.Sleep(1 * time.Second) // don't care about <1 s discrepancies
		lastNodeChainHeight := lastNode.GetChainHeight()
		t.Logf("current chain height: %d, last node chain height: %d", currentChainHeight, lastNodeChainHeight)
		require.GreaterOrEqual(t, lastNodeChainHeight, currentChainHeight)

		// shorten next epoch to speed up test
		newEpochPeriod := uint64(30)
		err = network[0].AwaitSetEpochPeriod(operatorKey, new(big.Int).SetUint64(newEpochPeriod), 3*time.Second)
		require.NoError(t, err)

		// wait for all new validators to join the committee
		t.Logf("epoch period %d", epochPeriod.Uint64())
		err = network.WaitForHeight(epochPeriod.Uint64(), int(epochPeriod.Uint64()))
		require.NoError(t, err)

		// committee should be bigger now and should be == to the clustering threshold + 1
		epochBlock := network[0].Eth.BlockChain().GetBlockByNumber(epochPeriod.Uint64())
		committee := epochBlock.Header().Epoch.Committee
		err = committee.Enrich()
		require.NoError(t, err)
		require.Greater(t, committee.Len(), customClusteringThreshold)

		for i, m := range committee.Members {
			t.Logf("index %d index %d voting power %v", i, m.Index, m.VotingPower.String())
		}

		// original members should be in the lowest spot of the committee
		require.True(t, slices.Contains(originalValidators, committee.Members[0].Address))
		require.True(t, slices.Contains(originalValidators, committee.Members[1].Address))
		require.True(t, slices.Contains(originalValidators, committee.Members[2].Address))

		// close the epoch
		err = network.WaitForHeight(epochPeriod.Uint64()+newEpochPeriod, int(epochPeriod.Uint64()*3))
		require.NoError(t, err)

		// we should be in epoch 2 now
		epochID, err = autonityContract.GetEpochID(nil)
		require.NoError(t, err)
		require.Equal(t, uint64(2), epochID.Uint64())

		// everyone should have 0 omission score
		omissionContract, err := bindings.NewOmissionAccountability(params.OmissionAccountabilityContractAddress, network[0].WsClient)
		require.NoError(t, err)
		for _, member := range committee.Members {
			score, err := omissionContract.GetInactivityScore(nil, member.Address)
			require.NoError(t, err)
			t.Logf("address %s score %d original validator %v", member.Address.String(), score.Uint64(), slices.Contains(originalValidators, member.Address))
			assert.Equal(t, uint64(0), score.Uint64())
		}
	})
}

func resetNode(t *testing.T, node *e2e.Node) {
	err := node.Close(false)
	node.Wait()
	require.NoError(t, err)
	err = node.Start()
	require.NoError(t, err)
}

func stopNode(t *testing.T, net e2e.Network, id int) {
	err := net[id].Close(false)
	net[id].Wait()
	require.NoError(t, err)
	time.Sleep(time.Second * 10)
}

func startNode(t *testing.T, net e2e.Network, id int) {
	err := net[id].Start()
	require.NoError(t, err)
}
