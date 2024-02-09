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

	"github.com/autonity/autonity/crypto/blst"

	"github.com/davecgh/go-spew/spew"
	"go.uber.org/goleak"
	"golang.org/x/sync/errgroup"

	"github.com/autonity/autonity/common/fdlimit"
	"github.com/autonity/autonity/common/graph"
	"github.com/autonity/autonity/common/keygenerator"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/metrics"
	"github.com/autonity/autonity/p2p"
	"github.com/autonity/autonity/p2p/enode"
)

const (
	ValidatorPrefix = "V"
	ExternalPrefix  = "E"
)

type testCase struct {
	name                   string
	numValidators          int
	numBlocks              int
	validatorsCanBeStopped *int64

	networkRates    map[string]networkRate //map[validatorIndex]networkRate
	beforeHooks     map[string]hook        //map[validatorIndex]beforeHook
	afterHooks      map[string]hook        //map[validatorIndex]afterHook
	finalAssert     func(t *testing.T, validators map[string]*testNode)
	genesisHook     func(g *core.Genesis) *core.Genesis
	mu              sync.RWMutex
	topology        *Topology
	skipNoLeakCheck bool
}

type hook func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error

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

func runTest(t *testing.T, test *testCase) {
	if !test.skipNoLeakCheck {
		// TODO: (screwyprof) Fix the following gorotine leaks
		defer goleak.VerifyNone(t,
			goleak.IgnoreTopFunction("github.com/JekaMas/notify._Cfunc_CFRunLoopRun"),
			goleak.IgnoreTopFunction("github.com/JekaMas/notify.(*nonrecursiveTree).internal"),
			goleak.IgnoreTopFunction("internal/poll.runtime_pollWait"),
			goleak.IgnoreTopFunction("net/http.(*persistConn).writeLoop"),
			goleak.IgnoreTopFunction("github.com/autonity/autonity/miner.(*worker).loop"),
			goleak.IgnoreTopFunction("github.com/autonity/autonity/miner.(*worker).updater"),
			goleak.IgnoreTopFunction("github.com/autonity/autonity/miner.(*worker).newWorkLoop.func1"),
		)
	}

	// needed to prevent go-routine leak at github.com/autonity/autonity/metrics.(*meterArbiter).tick
	// see: metrics/meter.go:55
	defer metrics.DefaultRegistry.UnregisterAll()

	log.Root().SetHandler(log.LvlFilterHandler(log.LvlInfo, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))
	_, err := fdlimit.Raise(512 * uint64(test.numValidators))
	if err != nil {
		t.Log("can't rise file description limit. errors are possible")
	}

	// prepare node names for validators
	nodeNames := getNodeNames(test.numValidators)

	// if topology test is enabled, then prepare node names from topologies.
	if test.topology != nil {
		err := test.topology.Validate()
		if err != nil {
			t.Fatal(err)
		}
		nodeNames = getNodeNamesByPrefix(test.topology.graph.GetNames(), ValidatorPrefix)
		test.numValidators = len(nodeNames)
		externalNames := getNodeNamesByPrefix(test.topology.graph.GetNames(), ExternalPrefix)
		nodeNames = append(nodeNames, externalNames...)
	}

	nodesNum := len(nodeNames)
	// Generate a batch of accounts to seal and fund with
	nodes := make(map[string]*testNode, nodesNum)
	generateNodesPrivateKeys(t, nodes, nodeNames, nodesNum)

	// Replace normal DNS resolver with the resolver for this test framework.
	enode.V4ResolveFunc = func(host string) (ips []net.IP, e error) {
		if len(host) > 4 || !(strings.HasPrefix(host, ValidatorPrefix) ||
			strings.HasPrefix(host, ExternalPrefix)) {
			return nil, &net.DNSError{Err: "not found", Name: host, IsNotFound: true}
		}

		return []net.IP{
			net.ParseIP("127.0.0.1"),
		}, nil
	}

	setNodesPortAndEnode(t, nodes)

	// Make genesis and apply customized genesis configurations.
	genesis := makeGenesis(t, nodes, nodeNames)
	if test.genesisHook != nil {
		genesis = test.genesisHook(genesis)
	}

	// Start the node as an application container for ethereum service.
	wg := &errgroup.Group{}
	for i, peer := range nodes {
		peer := peer
		peer.listener[0].Close()
		peer.listener[1].Close()
		peer.listener[2].Close()

		rates := test.networkRates[i]
		peer.nodeConfig, peer.ethConfig = makeNodeConfig(t, genesis, peer.nodeKey, peer.consensusKey,
			fmt.Sprintf("127.0.0.1:%d", peer.port),
			fmt.Sprintf("127.0.0.1:%d", peer.acnPort),
			peer.rpcPort, rates.in, rates.out)

		wg.Go(func() error {
			// if we have only a single validator, force mining start to bypass sync check
			return peer.startNode(nodesNum == 1)
		})
	}

	err = wg.Wait()
	if err != nil {
		t.Fatal(err)
	}

	for nodeName, node := range nodes {
		fmt.Printf("%s === %s  -- %s\n", nodeName, node.enode.URLv4(), crypto.PubkeyToAddress(node.nodeKey.PublicKey).String())
	}

	// apply topology changes over test network.
	if test.topology != nil && !test.topology.WithChanges() {
		err := test.topology.ConnectNodes(nodes)
		if err != nil {
			t.Fatal(err)
		}
	}

	// cleaners for node on shutdown.
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
				return os.RemoveAll(peer.nodeConfig.DataDir)
			})
		}

		err = wgClose.Wait()
		if err != nil {
			t.Fatal(err)
		}

		// level DB needs a second to close
		time.Sleep(time.Second)
	}()

	// init test controller and start mining.
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

	// start test controllers.
	startTestControllers(t, test, nodes, true)
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
		m := t.transformPeerListToMap(nodes[nodeKey].node.ExecutionServer().Peers(), nodes)
		for k := range connectionsList {
			if _, ok := m[k]; ok {
				continue
			}
			nodes[nodeKey].node.ExecutionServer().AddPeer(nodes[k].node.ExecutionServer().Self())
			nodes[nodeKey].node.ConsensusServer().AddPeer(nodes[k].node.ConsensusServer().Self())
		}
		for k := range m {
			if _, ok := connectionsList[k]; ok {
				continue
			}
			nodes[nodeKey].node.ExecutionServer().RemovePeer(nodes[k].node.ExecutionServer().Self())
			nodes[nodeKey].node.ConsensusServer().RemovePeer(nodes[k].node.ConsensusServer().Self())
		}
	}

	return nil
}

func (t *Topology) transformPeerListToMap(peers []*p2p.Peer, nodes map[string]*testNode) map[string]struct{} {
	m := make(map[string]struct{})
	mapper := make(map[enode.ID]string, len(nodes))
	for index, n := range nodes {
		mapper[n.node.ExecutionServer().Self().ID()] = index
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
		peers := nodes[i].node.ExecutionServer().Peers()
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
		peers := v.node.ExecutionServer().Peers()
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
	peers := node.node.ExecutionServer().Peers()
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
	peers := nodes[index].node.ExecutionServer().Peers()
	currentConnections := t.transformPeerListToMap(peers, nodes)
	for k := range currentConnections {
		if _, ok := graphConnections[k]; ok {
			continue
		}
		fmt.Println("node", index, "removes to", k)
		nodes[index].node.ExecutionServer().RemovePeer(nodes[k].node.ExecutionServer().Self())
		nodes[index].node.ExecutionServer().RemoveTrustedPeer(nodes[k].node.ExecutionServer().Self())
		nodes[index].node.ConsensusServer().RemovePeer(nodes[k].node.ConsensusServer().Self())
		nodes[index].node.ConsensusServer().RemoveTrustedPeer(nodes[k].node.ConsensusServer().Self())
	}

	for k := range graphConnections {
		if _, ok := currentConnections[k]; ok {
			continue
		}
		fmt.Println("node", index, "connects to", k)
		nodes[index].node.ExecutionServer().AddPeer(nodes[k].node.ExecutionServer().Self())
		nodes[index].node.ExecutionServer().AddTrustedPeer(nodes[k].node.ExecutionServer().Self())
		nodes[index].node.ConsensusServer().AddPeer(nodes[k].node.ConsensusServer().Self())
		nodes[index].node.ConsensusServer().AddTrustedPeer(nodes[k].node.ConsensusServer().Self())
	}

	return nil
}
func getNodeNames(max int) []string {
	res := make([]string, max)
	for i := range res {
		res[i] = "V" + strconv.Itoa(i)
	}
	return res
}

func generateNodesPrivateKeys(t *testing.T, nodes map[string]*testNode, nodeNames []string, nodesNum int) {
	var err error
	for i := 0; i < nodesNum; i++ {
		nodes[nodeNames[i]] = new(testNode)
		nodes[nodeNames[i]].consensusKey, err = blst.RandKey()
		if err != nil {
			t.Fatal("cannot make consensus key")
		}
		nodes[nodeNames[i]].nodeKey, err = keygenerator.Next()
		if err != nil {
			t.Fatal("cant make node pk", err)
		}
		nodes[nodeNames[i]].oracleKey, err = keygenerator.Next()
		if err != nil {
			t.Fatal("cant make oracle pk", err)
		}
	}
}

func setNodesPortAndEnode(t *testing.T, nodes map[string]*testNode) {
	for addr, node := range nodes {
		if n, err := newNode(node.nodeKey, addr); err != nil {
			t.Fatal(err)
		} else {
			node.netNode = n
		}
	}
}

func newNode(privateKey *ecdsa.PrivateKey, addr string) (netNode, error) {
	n := netNode{
		nodeKey: privateKey,
		address: crypto.PubkeyToAddress(privateKey.PublicKey),
	}

	// atc listener
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return netNode{}, err
	}
	n.listener = append(n.listener, listener)

	// eth listener
	listener, err = net.Listen("tcp", "127.0.0.1:0")
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

	acnPort := strings.Split(n.listener[0].Addr().String(), ":")[1]
	n.acnhost = fmt.Sprintf("%s:%s", addr, acnPort)
	n.acnPort, _ = strconv.Atoi(acnPort)

	port := strings.Split(n.listener[1].Addr().String(), ":")[1]
	n.host = fmt.Sprintf("%s:%s", addr, port)
	n.port, _ = strconv.Atoi(port)

	rpcListener := n.listener[2]
	rpcPort, innerErr := strconv.Atoi(strings.Split(rpcListener.Addr().String(), ":")[1])
	if innerErr != nil {
		return netNode{}, fmt.Errorf("incorrect rpc port %w", innerErr)
	}

	n.rpcPort = rpcPort

	if n.port == 0 || n.rpcPort == 0 || n.acnPort == 0 {
		return netNode{}, fmt.Errorf("on node %s port equals 0", addr)
	}

	n.url = enode.V4DNSUrl(
		n.nodeKey.PublicKey,
		n.host,
		n.port,
		n.port,
	)
	n.url = enode.AppendConsensusEndpoint(addr, strconv.Itoa(n.acnPort), n.url)

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
