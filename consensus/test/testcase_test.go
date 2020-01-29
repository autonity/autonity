package test

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/common/fdlimit"
	"github.com/clearmatics/autonity/common/graph"
	"github.com/clearmatics/autonity/consensus"
	tendermintCore "github.com/clearmatics/autonity/consensus/tendermint/core"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/metrics"
	"github.com/clearmatics/autonity/p2p/enode"
	"github.com/davecgh/go-spew/spew"
	"go.uber.org/goleak"
	"golang.org/x/sync/errgroup"
)

const (
	ValidatorPrefix   = "V"
	StakeholderPrefix = "S"
	ParticipantPrefix = "P"
)

type testCase struct {
	name                   string
	isSkipped              bool
	numValidators          int
	numBlocks              int
	txPerPeer              int
	validatorsCanBeStopped *int64

	maliciousPeers          map[string]injectors
	removedPeers            map[common.Address]uint64
	addedValidatorsBlocks   map[common.Hash]uint64
	removedValidatorsBlocks map[common.Hash]uint64 //nolint: unused, structcheck
	changedValidators       tendermintCore.Changes //nolint: unused,structcheck

	networkRates         map[string]networkRate //map[validatorIndex]networkRate
	beforeHooks          map[string]hook        //map[validatorIndex]beforeHook
	afterHooks           map[string]hook        //map[validatorIndex]afterHook
	sendTransactionHooks map[string]func(validator *testNode, fromAddr common.Address, toAddr common.Address) (bool, *types.Transaction, error)
	finalAssert          func(t *testing.T, validators map[string]*testNode)
	stopTime             map[string]time.Time
	genesisHook          func(g *core.Genesis) *core.Genesis
	mu                   sync.RWMutex
	noQuorumAfterBlock   uint64
	noQuorumTimeout      time.Duration
	topology             *Topology
}

type injectors struct {
	cons  func(basic consensus.Engine) consensus.Engine
	backs func(basic tendermintCore.Backend) tendermintCore.Backend
}

func (test *testCase) getBeforeHook(index string) hook {
	test.mu.Lock()
	defer test.mu.Unlock()

	if test.beforeHooks == nil {
		return nil
	}

	validatorHook, ok := test.beforeHooks[index]
	if !ok || validatorHook == nil {
		return nil
	}

	return validatorHook
}

func (test *testCase) getAfterHook(index string) hook {
	test.mu.Lock()
	defer test.mu.Unlock()

	if test.afterHooks == nil {
		return nil
	}

	validatorHook, ok := test.afterHooks[index]
	if !ok || validatorHook == nil {
		return nil
	}

	return validatorHook
}

func (test *testCase) setStopTime(index string, stopTime time.Time) {
	test.mu.Lock()
	test.stopTime[index] = stopTime
	test.mu.Unlock()
}

func (test *testCase) getStopTime(index string) time.Time {
	test.mu.RLock()
	currentTime := test.stopTime[index]
	test.mu.RUnlock()

	return currentTime
}

type hook func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error

func runTest(t *testing.T, test *testCase) {
	if test.isSkipped {
		t.SkipNow()
	}

	// TODO: (screwyprof) Fix the following gorotine leaks
	defer goleak.VerifyNone(t,
		goleak.IgnoreTopFunction("github.com/JekaMas/notify._Cfunc_CFRunLoopRun"),
		//goleak.IgnoreTopFunction("github.com/clearmatics/autonity/metrics.(*meterArbiter).tick"),
		goleak.IgnoreTopFunction("github.com/clearmatics/autonity/consensus/ethash.(*remoteSealer).loop"))

	// needed to prevent go-routine leak at github.com/clearmatics/autonity/metrics.(*meterArbiter).tick
	// see: metrics/meter.go:55
	defer metrics.DefaultRegistry.UnregisterAll()

	log.Root().SetHandler(log.LvlFilterHandler(log.LvlError, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))
	_, err := fdlimit.Raise(512 * uint64(test.numValidators))
	if err != nil {
		t.Log("can't rise file description limit. errors are possible")
	}

	nodeNames := getNodeNames()[:test.numValidators]
	if test.topology != nil {
		nodeNames = getNodeNamesByPrefix(test.topology.graph.GetNames(), ValidatorPrefix)
		test.numValidators = len(nodeNames)

		stakeholderNames := getNodeNamesByPrefix(test.topology.graph.GetNames(), StakeholderPrefix)
		participantNames := getNodeNamesByPrefix(test.topology.graph.GetNames(), ParticipantPrefix)
		nodeNames = append(nodeNames, stakeholderNames...)
		nodeNames = append(nodeNames, participantNames...)
	}
	nodesNum := len(nodeNames)
	// Generate a batch of accounts to seal and fund with
	nodes := make(map[string]*testNode, nodesNum)

	generateNodesPrivateKey(t, nodes, nodeNames, nodesNum)
	setNodesPortAndEnode(t, nodes)

	genesis := makeGenesis(nodes)
	if test.genesisHook != nil {
		genesis = test.genesisHook(genesis)
	}
	for i, validator := range nodes {
		var engineConstructor func(basic consensus.Engine) consensus.Engine
		var backendConstructor func(basic tendermintCore.Backend) tendermintCore.Backend
		if test.maliciousPeers != nil {
			engineConstructor = test.maliciousPeers[i].cons
			backendConstructor = test.maliciousPeers[i].backs
		}

		validator.listener[0].Close()
		validator.listener[1].Close()

		rates := test.networkRates[i]

		validator.node, err = makeValidator(genesis, validator.privateKey, validator.address, validator.rpcPort, rates.in, rates.out, engineConstructor, backendConstructor)
		if err != nil {
			t.Fatal("cant make a node", i, err)
		}
	}

	wg := &errgroup.Group{}
	for _, validator := range nodes {
		validator := validator

		wg.Go(func() error {
			return validator.startNode()
		})
	}
	err = wg.Wait()
	if err != nil {
		t.Fatal(err)
	}

	if test.topology != nil {
		for _, v := range test.topology.graph.Edges {
			nodes[v.LeftNode].node.Server().AddPeer(nodes[v.RightNode].node.Server().Self())
		}
	}

	defer func() {
		wgClose := &errgroup.Group{}
		for _, validator := range nodes {
			validatorInner := validator
			wgClose.Go(func() error {
				if !validatorInner.isRunning {
					return nil
				}

				errInner := validatorInner.node.Close()
				if errInner != nil {
					return fmt.Errorf("error on node close %v", err)
				}

				validatorInner.node.Wait()

				return nil
			})
		}

		err = wgClose.Wait()
		if err != nil {
			t.Fatal(err)
		}

		time.Sleep(time.Second) //level DB needs a second to close
	}()

	wg = &errgroup.Group{}
	for _, validator := range nodes {
		validator := validator

		wg.Go(func() error {
			return validator.startService()
		})
	}
	err = wg.Wait()
	if err != nil {
		t.Fatal(err)
	}

	wg = &errgroup.Group{}
	for i, validator := range nodes {
		validator := validator
		i := i

		wg.Go(func() error {
			log.Debug("peers", "i", i,
				"peers", len(validator.node.Server().Peers()),
				"nodes", len(nodes))
			return nil
		})
	}
	err = wg.Wait()
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		for _, validator := range nodes {
			validator.subscription.Unsubscribe()
		}
	}()

	// each peer sends one tx per block
	sendTransactions(t, test, nodes, test.txPerPeer, true, nodeNames)
	if test.finalAssert != nil {
		test.finalAssert(t, nodes)
	}
	//check topology
	if test.topology != nil {
		missedConnections := []graph.Edge{}
		for _, v := range test.topology.graph.Edges {
			exists := false
			for _, nd := range nodes[v.LeftNode].node.Server().Peers() {
				if nodes[v.RightNode].node.Server().Self().ID() == nd.ID() {
					exists = true
				}
			}
			if !exists {
				missedConnections = append(missedConnections, *v)
			}
		}
		if len(missedConnections) != 0 {
			spew.Dump(missedConnections)
			t.Fatal("Some connections missed")
		}
	}

	if len(test.maliciousPeers) != 0 {
		maliciousTest(t, test, nodes)
	}
}

type Topology struct {
	graph graph.Graph
}

func getNodeNames() []string {
	return []string{
		"VA", "VB", "VC", "VD", "VE", "VF", "VG", "VH", "VI", "VJ", "VK",
	}
}

func generateNodesPrivateKey(t *testing.T, nodes map[string]*testNode, nodeNames []string, nodesNum int) {
	var err error
	for i := 0; i < nodesNum; i++ {
		nodes[nodeNames[i]] = new(testNode)
		nodes[nodeNames[i]].privateKey, err = crypto.GenerateKey()
		if err != nil {
			t.Fatal("cant make pk", err)
		}
	}
}

func setNodesPortAndEnode(t *testing.T, nodes map[string]*testNode) {
	for i := range nodes {
		//port
		listener, innerErr := net.Listen("tcp", "127.0.0.1:0")
		if innerErr != nil {
			panic(innerErr)
		}
		nodes[i].listener = append(nodes[i].listener, listener)

		//rpc port
		listener, innerErr = net.Listen("tcp", "127.0.0.1:0")
		if innerErr != nil {
			panic(innerErr)
		}
		nodes[i].listener = append(nodes[i].listener, listener)
	}

	for i, node := range nodes {
		listener := node.listener[0]
		node.address = listener.Addr().String()
		port := strings.Split(listener.Addr().String(), ":")[1]
		node.port, _ = strconv.Atoi(port)

		rpcListener := node.listener[1]
		rpcPort, innerErr := strconv.Atoi(strings.Split(rpcListener.Addr().String(), ":")[1])
		if innerErr != nil {
			t.Fatal("incorrect rpc port ", innerErr)
		}

		node.rpcPort = rpcPort

		if node.port == 0 || node.rpcPort == 0 {
			t.Fatal("On node", i, "port equals 0")
		}

		node.url = enode.V4URL(
			node.privateKey.PublicKey,
			net.IPv4(127, 0, 0, 1),
			node.port,
			node.port,
		)
	}
}

func getNodeNamesByPrefix(names []string, typ string) []string {
	validators := make([]string, 0, len(names))
	for _, v := range names {
		if len(v) == 0 {
			continue
		}
		if strings.HasPrefix(v, typ) {
			validators = append(validators, v)
		}
	}
	return validators
}
