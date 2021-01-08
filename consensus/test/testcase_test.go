package test

import (
	"crypto/ecdsa"
	"errors"
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
	"github.com/clearmatics/autonity/common/keygenerator"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/metrics"
	"github.com/clearmatics/autonity/p2p"
	"github.com/clearmatics/autonity/p2p/enode"
	"github.com/davecgh/go-spew/spew"
	"go.uber.org/goleak"
	"golang.org/x/sync/errgroup"
)

const (
	ValidatorPrefix   = "V"
	StakeholderPrefix = "S"
	ParticipantPrefix = "P"
	ExternalPrefix    = "E"
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

	networkRates         map[string]networkRate //map[validatorIndex]networkRate
	beforeHooks          map[string]hook        //map[validatorIndex]beforeHook
	afterHooks           map[string]hook        //map[validatorIndex]afterHook
	sendTransactionHooks map[string]sendTransactionHook
	finalAssert          func(t *testing.T, validators map[string]*testNode)
	stopTime             map[string]time.Time
	genesisHook          func(g *core.Genesis) *core.Genesis
	mu                   sync.RWMutex
	noQuorumAfterBlock   uint64
	noQuorumTimeout      time.Duration
	topology             *Topology
	skipNoLeakCheck      bool
}

type injectors struct {
	cons func(basic consensus.Engine) consensus.Engine
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
type sendTransactionHook func(validator *testNode, fromAddr common.Address, toAddr common.Address) (bool, *types.Transaction, error)

func runTest(t *testing.T, test *testCase) {
	if test.isSkipped {
		t.SkipNow()
	}

	// fixme the noQuorum tests does not close nodes properly
	if !test.skipNoLeakCheck {
		// TODO: (screwyprof) Fix the following gorotine leaks
		defer goleak.VerifyNone(t,
			goleak.IgnoreTopFunction("github.com/JekaMas/notify._Cfunc_CFRunLoopRun"),
			goleak.IgnoreTopFunction("github.com/JekaMas/notify.(*nonrecursiveTree).internal"),
			goleak.IgnoreTopFunction("internal/poll.runtime_pollWait"),
			goleak.IgnoreTopFunction("net/http.(*persistConn).writeLoop"),
			goleak.IgnoreTopFunction("github.com/clearmatics/autonity/miner.(*worker).loop"),
			goleak.IgnoreTopFunction("github.com/clearmatics/autonity/miner.(*worker).updater"),
			goleak.IgnoreTopFunction("github.com/clearmatics/autonity/miner.(*worker).newWorkLoop.func1"),
		)
	}

	// needed to prevent go-routine leak at github.com/clearmatics/autonity/metrics.(*meterArbiter).tick
	// see: metrics/meter.go:55
	defer metrics.DefaultRegistry.UnregisterAll()

	log.Root().SetHandler(log.LvlFilterHandler(log.LvlCrit, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))
	_, err := fdlimit.Raise(512 * uint64(test.numValidators))
	if err != nil {
		t.Log("can't rise file description limit. errors are possible")
	}

	// default topology if not set anything
	nodeNames := getNodeNames()[:test.numValidators]
	stakeholderName := nodeNames[len(nodeNames)-1]

	if test.topology != nil {
		err := test.topology.Validate()
		if err != nil {
			t.Fatal(err)
		}
		nodeNames = getNodeNamesByPrefix(test.topology.graph.GetNames(), ValidatorPrefix)
		test.numValidators = len(nodeNames)

		stakeholderNames := getNodeNamesByPrefix(test.topology.graph.GetNames(), StakeholderPrefix)
		if len(stakeholderNames) > 0 {
			stakeholderName = stakeholderNames[0]
		}

		participantNames := getNodeNamesByPrefix(test.topology.graph.GetNames(), ParticipantPrefix)
		externalNames := getNodeNamesByPrefix(test.topology.graph.GetNames(), ExternalPrefix)

		nodeNames = append(nodeNames, stakeholderNames...)
		nodeNames = append(nodeNames, participantNames...)
		nodeNames = append(nodeNames, externalNames...)
	}

	nodesNum := len(nodeNames)
	// Generate a batch of accounts to seal and fund with
	nodes := make(map[string]*testNode, nodesNum)

	// Replace normal resolver with resolver that can resolve our test hostnames
	enode.V4ResolveFunc = func(host string) (ips []net.IP, e error) {
		if len(host) > 4 || !(strings.HasPrefix(host, ValidatorPrefix) ||
			strings.HasPrefix(host, StakeholderPrefix) ||
			strings.HasPrefix(host, ParticipantPrefix) ||
			strings.HasPrefix(host, ExternalPrefix)) {
			return nil, &net.DNSError{Err: "not found", Name: host, IsNotFound: true}
		}

		return []net.IP{
			net.ParseIP("127.0.0.1"),
		}, nil
	}

	generateNodesPrivateKey(t, nodes, nodeNames, nodesNum)
	setNodesPortAndEnode(t, nodes)

	genesis := makeGenesis(nodes, stakeholderName)

	if test.genesisHook != nil {
		genesis = test.genesisHook(genesis)
	}

	for i, peer := range nodes {
		var engineConstructor func(basic consensus.Engine) consensus.Engine
		if test.maliciousPeers != nil {
			engineConstructor = test.maliciousPeers[i].cons
		}

		peer.listener[0].Close()
		peer.listener[1].Close()

		rates := test.networkRates[i]
		peer.node, err = makePeer(genesis, peer.privateKey, fmt.Sprintf("127.0.0.1:%d", peer.port), peer.rpcPort, rates.in, rates.out, engineConstructor)
		if err != nil {
			t.Fatal("cant make a node", i, err)
		}
	}

	wg := &errgroup.Group{}
	for _, peer := range nodes {
		peer := peer

		wg.Go(func() error {
			return peer.startNode()
		})
	}
	err = wg.Wait()
	if err != nil {
		t.Fatal(err)
	}

	s := " "
	for i, p := range nodes {
		s += fmt.Sprintf("%s %s === %s  -- %s\n", s, i, p.enode.URLv4(), crypto.PubkeyToAddress(p.privateKey.PublicKey).String())

	}
	fmt.Println(s)

	if test.topology != nil && !test.topology.WithChanges() {
		err := test.topology.ConnectNodes(nodes)
		if err != nil {
			t.Fatal(err)
		}
	}

	defer func() {
		wgClose := &errgroup.Group{}
		for _, peer := range nodes {
			peer := peer

			wgClose.Go(func() error {
				if !peer.isRunning {
					return nil
				}

				errInner := peer.node.Close()
				if errInner != nil {
					return fmt.Errorf("error on node close %v", err)
				}

				peer.node.Wait()

				return nil
			})
		}

		err = wgClose.Wait()
		if err != nil {
			t.Fatal(err)
		}

		// level DB needs a second to close
		time.Sleep(time.Second)
	}()

	wg = &errgroup.Group{}
	for _, peer := range nodes {
		peer := peer

		wg.Go(func() error {
			return peer.startService()
		})
	}
	err = wg.Wait()
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		for _, peer := range nodes {
			peer.subscription.Unsubscribe()
		}
	}()

	// each peer sends one tx per block
	sendTransactions(t, test, nodes, test.txPerPeer, true, nodeNames)
	if test.finalAssert != nil {
		test.finalAssert(t, nodes)
	}

	// check topology
	if test.topology != nil {
		err := test.topology.CheckTopology(nodes)
		if err != nil {
			t.Fatal(err)
		}
	}

	if len(test.maliciousPeers) != 0 {
		maliciousTest(t, test, nodes)
	}
}

func TestResolve(t *testing.T) {
	customResolve := func(host string) (ips []net.IP, e error) {
		return []net.IP{
			net.ParseIP("127.0.0.1"),
		}, nil
	}

	// Set the resolver function for the enode package
	enode.V4ResolveFunc = customResolve

	en := "enode://57fa76dc95ef02461ce1a38d70181c27384f628a23a98fa801933ac2a45709b847d4ab42ed0fe0ebd03df5d464c064585a85a154e4443fb9143bfb6c369d5544@VD:45736"
	_, err := enode.ParseV4(en)
	if err != nil {
		t.Fatal(err)
	}
}

type Topology struct {
	graph graph.Graph
}

func (t *Topology) Validate() error {
	if len(t.graph.Edges) == 0 && len(t.graph.SubGraphs) == 0 {
		spew.Dump(t.graph)
		return errors.New("empty topology")
	}
	if len(t.graph.Edges) > 0 && len(t.graph.SubGraphs) > 0 {
		return errors.New("conflicting topologies")
	}
	for _, v := range t.graph.SubGraphs {
		if _, err := strconv.ParseUint(strings.TrimPrefix(v.Name, "b"), 10, 64); err != nil {
			return errors.New("incorrect block number")
		}
	}
	return nil
}

func (t *Topology) WithChanges() bool {
	return len(t.graph.SubGraphs) > 0
}
func (t *Topology) ConnectNodes(nodes map[string]*testNode) error {
	edges := t.getEdges(maxNumOfBlockMum(nodes))
	connections := t.getPeerConnections(edges)
	for nodeKey, connectionsList := range connections {
		m := t.transformPeerListToMap(nodes[nodeKey].node.Server().Peers(), nodes)
		for k := range connectionsList {
			if _, ok := m[k]; ok {
				continue
			}
			nodes[nodeKey].node.Server().AddPeer(nodes[k].node.Server().Self())
		}
		for k := range m {
			if _, ok := connectionsList[k]; ok {
				continue
			}
			nodes[nodeKey].node.Server().RemovePeer(nodes[k].node.Server().Self())
		}
	}

	return nil
}

func (t *Topology) transformPeerListToMap(peers []*p2p.Peer, nodes map[string]*testNode) map[string]struct{} {
	m := make(map[string]struct{})
	mapper := make(map[enode.ID]string, len(nodes))
	for index, n := range nodes {
		mapper[n.node.Server().Self().ID()] = index
	}
	for _, v := range peers {
		index, ok := mapper[v.Node().ID()]
		if ok {
			m[index] = struct{}{}
		} else {
			panic("Node doesn't exists")
		}
	}
	return m

}
func (t *Topology) getEdges(blockNum uint64) []*graph.Edge {
	var edges []*graph.Edge
	if t.WithChanges() {
		for _, v := range t.graph.SubGraphs {
			blockNumStr := strings.TrimPrefix(v.Name, "b")
			parsed, _ := strconv.ParseUint(blockNumStr, 10, 64)
			if blockNum >= parsed {
				edges = v.Edges
			}
		}
		if edges == nil {
			fmt.Println("empty edges")
			return nil
		}
	} else {
		edges = t.graph.Edges
	}
	return edges
}

func (t *Topology) getChangesBlocks() (map[uint64]struct{}, error) {
	m := make(map[uint64]struct{})
	if t.WithChanges() {
		for _, v := range t.graph.SubGraphs {
			blockNumStr := strings.TrimPrefix(v.Name, "b")
			parsed, err := strconv.ParseUint(blockNumStr, 10, 64)
			if err != nil {
				return nil, err
			}
			m[parsed] = struct{}{}

		}
	}
	return m, nil
}
func (t *Topology) getPeerConnections(edges []*graph.Edge) map[string]map[string]struct{} {
	res := make(map[string]map[string]struct{})
	for _, v := range edges {
		m, ok := res[v.LeftNode]
		if !ok {
			m = make(map[string]struct{})
		}
		m[v.RightNode] = struct{}{}
		res[v.LeftNode] = m

		m, ok = res[v.RightNode]
		if !ok {
			m = make(map[string]struct{})
		}
		m[v.LeftNode] = struct{}{}
		res[v.RightNode] = m
	}
	return res
}

func (t *Topology) CheckTopology(nodes map[string]*testNode) error {
	blockNum := maxNumOfBlockMum(nodes)
	edges := t.getEdges(blockNum)
	connections := t.getPeerConnections(edges)

	for i, v := range connections {
		peers := nodes[i].node.Server().Peers()
		m := t.transformPeerListToMap(peers, nodes)
		for j := range v {
			if _, ok := v[j]; !ok {
				spew.Dump(m)
				spew.Dump(v)
				spew.Dump(connections)

				return fmt.Errorf("CheckTopology incorrect topology for block %v for node %v", blockNum, i)
			}
		}

	}

	return nil
}

func (t *Topology) FullTopology(nodes map[string]*testNode) map[string]map[string]struct{} {
	m := make(map[string]map[string]struct{})
	for i, v := range nodes {
		peers := v.node.Server().Peers()
		byPeer := t.transformPeerListToMap(peers, nodes)
		m[i] = byPeer
	}

	return m
}

func (t *Topology) DumpTopology(nodes map[string]*testNode) string {
	m := t.FullTopology(nodes)
	s := ""
	for i := range m {
		s += i + "\n"
		s += dumpConnections(i, m[i])
		s += "\n"
	}
	return s

}
func (t *Topology) CheckTopologyForIndex(index string, nodes map[string]*testNode) error {
	node := nodes[index]
	blockNum := node.lastBlock

	fmt.Println("check topology", index, blockNum)
	if t.WithChanges() {
		m, err := t.getChangesBlocks()
		if err != nil {
			return err
		}
		for i := uint64(0); i < 10; i++ {
			if _, ok := m[blockNum-i]; ok {
				fmt.Println("blocknum check exit")
				return nil
			}

		}
	}
	edges := t.getEdges(blockNum)
	if edges == nil {
		return nil
	}
	fmt.Println("check started", index, blockNum)
	allConnections := t.getPeerConnections(edges)
	indexConnections := allConnections[index]
	peers := node.node.Server().Peers()
	m := t.transformPeerListToMap(peers, nodes)
	for i := range indexConnections {
		if _, ok := m[i]; !ok {
			fmt.Println("current", dumpConnections(index, m))
			fmt.Println()
			fmt.Println("must", dumpConnections(index, indexConnections))
			return fmt.Errorf("CheckTopologyForIndex incorrect topology for %v for block %v", index, blockNum)
		}
	}
	return nil
}

func (t *Topology) ConnectNodesForIndex(index string, nodes map[string]*testNode) error {
	blockNum := nodes[index].lastBlock
	ch, err := t.getChangesBlocks()
	if err != nil {
		return err
	}
	if _, ok := ch[blockNum]; !ok {
		return nil
	}
	edges := t.getEdges(blockNum)
	if len(edges) == 0 {
		return nil
	}
	fmt.Println("+ConnectNodesForIndex", index)
	defer fmt.Println("-ConnectNodesForIndex", index)
	allConnections := t.getPeerConnections(edges)
	graphConnections := allConnections[index]
	fmt.Println(dumpConnections(index, graphConnections))
	fmt.Println()
	peers := nodes[index].node.Server().Peers()
	currentConnections := t.transformPeerListToMap(peers, nodes)
	for k := range currentConnections {
		if _, ok := graphConnections[k]; ok {
			continue
		}
		fmt.Println("node", index, "removes to", k)
		nodes[index].node.Server().RemovePeer(nodes[k].node.Server().Self())
		nodes[index].node.Server().RemoveTrustedPeer(nodes[k].node.Server().Self())
	}

	for k := range graphConnections {
		if _, ok := currentConnections[k]; ok {
			continue
		}
		fmt.Println("node", index, "connects to", k)
		nodes[index].node.Server().AddPeer(nodes[k].node.Server().Self())
		nodes[index].node.Server().AddTrustedPeer(nodes[k].node.Server().Self())
	}

	return nil
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
		nodes[nodeNames[i]].privateKey, err = keygenerator.Next()
		if err != nil {
			t.Fatal("cant make pk", err)
		}
	}
}

func setNodesPortAndEnode(t *testing.T, nodes map[string]*testNode) {
	for addr, node := range nodes {
		if n, err := newNode(node.privateKey, addr); err != nil {
			t.Fatal(err)
		} else {
			node.netNode = n
		}
	}
}

func newNode(privateKey *ecdsa.PrivateKey, addr string) (netNode, error) {
	n := netNode{
		privateKey: privateKey,
	}

	//port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return netNode{}, err
	}
	n.listener = append(n.listener, listener)

	//rpc port
	listener, err = net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return netNode{}, err
	}
	n.listener = append(n.listener, listener)

	port := strings.Split(n.listener[0].Addr().String(), ":")[1]
	n.address = fmt.Sprintf("%s:%s", addr, port)
	n.port, _ = strconv.Atoi(port)

	rpcListener := n.listener[1]
	rpcPort, innerErr := strconv.Atoi(strings.Split(rpcListener.Addr().String(), ":")[1])
	if innerErr != nil {
		return netNode{}, fmt.Errorf("incorrect rpc port %w", innerErr)
	}

	n.rpcPort = rpcPort

	if n.port == 0 || n.rpcPort == 0 {
		return netNode{}, fmt.Errorf("on node %s port equals 0", addr)
	}

	n.url = enode.V4DNSUrl(
		n.privateKey.PublicKey,
		n.address,
		n.port,
		n.port,
	)

	return n, nil
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

func maxNumOfBlockMum(nodes map[string]*testNode) uint64 {
	m := make(map[uint64]int)
	for i := range nodes {
		m[nodes[i].lastBlock]++
	}
	var max int
	var blockNum uint64
	for i, v := range m {
		if v > max {
			max = v
			blockNum = i
		}
	}
	return blockNum
}

func dumpConnections(index string, nodes map[string]struct{}) string {
	s := ""
	for i := range nodes {
		s += index + "---" + i + "\n"
	}
	return s
}
