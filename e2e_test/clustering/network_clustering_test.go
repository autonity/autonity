package clustering

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	routerInterfaces "github.com/autonity/autonity/consensus/tendermint/router/interfaces"
	"github.com/autonity/autonity/core/types"
	e2e "github.com/autonity/autonity/e2e_test"
)

// TestClusteringHappyCase is a happy case to test 5 clusters with each of them contains 5 nodes. The latency measurement
// is a base on an local simulator which generates [0, 500) ms delays.
func TestClusteringHappyCase(t *testing.T) {
	// mocked service with a local ping simulator which generates [0, 500) ms latency.
	mockedService := &interfaces.Services{Pinger: NewSimulatedPinger()}

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
	network.WaitToMineNBlocks(250, 250, false)
}

// TestClusteringResetAllNodes, it stops all nodes one by one, and start them again one by one. The network should recover to
// mining.
func TestClusteringResetAllNodes(t *testing.T) {
	numOfNodes := 36
	mockedService := &interfaces.Services{Pinger: NewSimulatedPinger()}

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
	mockedService := &interfaces.Services{Pinger: NewSimulatedPinger()}

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

func (r *NoRelayingSelector) SetBroadcaster(broadcaster routerInterfaces.PeerFinder) {
	//TODO implement me
	panic("implement me")
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
	pinger := NewSimulatedPinger()

	validators, err := e2e.Validators(t, numOfNodes, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	// Create the network with F num of nodes which does not relay proposal.
	f := bft.F(new(big.Int).SetInt64(int64(numOfNodes))).Int64()
	for i, validator := range validators {
		mockedService := &interfaces.Services{Pinger: pinger}
		if int64(i) < f {
			mockedService.Selector = &NoRelayingSelector{}
		}
		validator.TendermintServices = mockedService
	}

	network, err := e2e.NewNetworkFromValidators(t, validators, true)
	require.NoError(t, err)
	defer network.Shutdown(t)

	// wait for the consensus engine to work.
	network.WaitToMineNBlocks(500, 500, false)
}

func Test2FFaultyRelayers(t *testing.T) {
	numOfNodes := 36
	pinger := NewSimulatedPinger()

	validators, err := e2e.Validators(t, numOfNodes, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	// Create the network with F num of nodes which does not relay proposal.
	f := bft.F(new(big.Int).SetInt64(int64(numOfNodes))).Int64()
	overF := f * 2
	for i, validator := range validators {
		mockedService := &interfaces.Services{Pinger: pinger}
		if int64(i) < overF {
			mockedService.Selector = &NoRelayingSelector{}
		}
		validator.TendermintServices = mockedService
	}

	network, err := e2e.NewNetworkFromValidators(t, validators, true)
	require.NoError(t, err)
	defer network.Shutdown(t)

	// wait for the consensus engine to work.
	network.WaitToMineNBlocks(500, 500, false)
}

func Test3FFaultyRelayers(t *testing.T) {
	numOfNodes := 36
	pinger := NewSimulatedPinger()

	validators, err := e2e.Validators(t, numOfNodes, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	// Create the network with F num of nodes which does not relay proposal.
	f := bft.F(new(big.Int).SetInt64(int64(numOfNodes))).Int64()
	overF := f * 3
	for i, validator := range validators {
		mockedService := &interfaces.Services{Pinger: pinger}
		if int64(i) < overF {
			mockedService.Selector = &NoRelayingSelector{}
		}
		validator.TendermintServices = mockedService
	}

	network, err := e2e.NewNetworkFromValidators(t, validators, true)
	require.NoError(t, err)
	defer network.Shutdown(t)

	// wait for the consensus engine to work.
	network.WaitToMineNBlocks(60, 60, false)
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
