package tenclient

import (
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"net"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/clearmatics/autonity/accounts"
	"github.com/clearmatics/autonity/accounts/keystore"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/common/keygenerator"
	"github.com/clearmatics/autonity/consensus/tendermint/config"
	tendermintCore "github.com/clearmatics/autonity/consensus/tendermint/core"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/eth"
	"github.com/clearmatics/autonity/eth/downloader"
	"github.com/clearmatics/autonity/event"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/node"
	"github.com/clearmatics/autonity/p2p"
	"github.com/clearmatics/autonity/p2p/enode"
	"github.com/clearmatics/autonity/params"
	"github.com/stretchr/testify/assert"
)

const (
	ValidatorPrefix   = "V"
	StakeholderPrefix = "S"
	ParticipantPrefix = "P"
)

type testNode struct {
	isRunning      bool
	isInited       bool
	wasStopped     bool //fixme should be removed
	privateKey     *ecdsa.PrivateKey
	address        string
	port           int
	url            string
	listener       []net.Listener
	rpcPort        int
	node           *node.Node
	enode          *enode.Node
	service        *eth.Ethereum
	eventChan      chan core.ChainEvent
	subscription   event.Subscription
	transactions   map[common.Hash]struct{}
	transactionsMu sync.Mutex
	blocks         map[uint64]block
	lastBlock      uint64
	txsSendCount   *int64
	txsChainCount  map[uint64]int64
	isMalicious    bool
}

func (validator *testNode) startNode() error {
	// Inject the signer key and start sealing with it
	store := validator.node.AccountManager().Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)

	var (
		err    error
		signer accounts.Account
	)

	if !validator.isInited {
		signer, err = store.ImportECDSA(validator.privateKey, "")
		if err != nil {
			return fmt.Errorf("import pk: %s", err)
		}

		for {
			// wait until the private key is imported
			_, err = validator.node.AccountManager().Find(signer)
			if err == nil {
				break
			}
			time.Sleep(50 * time.Microsecond)
		}

		validator.isInited = true
	} else {
		signer = store.Accounts()[0]
	}

	if err = store.Unlock(signer, ""); err != nil {
		return fmt.Errorf("cant unlock: %s", err)
	}

	validator.node.ResetEventMux()

	if err := validator.node.Start(); err != nil {
		return fmt.Errorf("cannot start a node %s", err)
	}

	// Start tracking the node and it's enode
	validator.enode = validator.node.Server().Self()
	return nil
}

type block struct {
	hash common.Hash
	txs  int
}

type testCase struct {
	name                   string
	isSkipped              bool
	numValidators          int
	numBlocks              int
	txPerPeer              int
	validatorsCanBeStopped *int64

	removedPeers            map[common.Address]uint64
	addedValidatorsBlocks   map[common.Hash]uint64
	removedValidatorsBlocks map[common.Hash]uint64 //nolint: unused, structcheck
	changedValidators       tendermintCore.Changes //nolint: unused,structcheck

	sendTransactionHooks map[string]func(validator *testNode, fromAddr common.Address, toAddr common.Address) (bool, *types.Transaction, error)
	finalAssert          func(t *testing.T, validators map[string]*testNode)
	stopTime             map[string]time.Time
	mu                   sync.RWMutex
	noQuorumAfterBlock   uint64
	noQuorumTimeout      time.Duration

	//use this to add a peer nodes[nodeKey].node.Server().AddPeer(nodes[k].node.Server().Self())
	//topology        *Topology
	skipNoLeakCheck bool
}

func setupNodes(t *testing.T, test *testCase) map[string]*testNode {

	nodeNames := []string{"VA", "VB", "VC"}
	nodesNum := len(nodeNames)
	nodes := make(map[string]*testNode, nodesNum)

	// This looks kinda horrible
	enode.SetResolveFunc(func(host string) (ips []net.IP, e error) {
		if len(host) > 4 || !(strings.HasPrefix(host, ValidatorPrefix) ||
			strings.HasPrefix(host, StakeholderPrefix) ||
			strings.HasPrefix(host, ParticipantPrefix)) {
			return nil, &net.DNSError{Err: "not found", Name: host, IsNotFound: true}
		}

		return []net.IP{
			net.ParseIP("127.0.0.1"),
		}, nil
	})
	generateNodesPrivateKey(t, nodes, nodeNames, nodesNum)
	setNodesPortAndEnode(t, nodes)

	genesis := makeGenesis(nodes)
	for i, validator := range nodes {
		var err error

		validator.listener[0].Close()
		validator.listener[1].Close()

		validator.node, err = makeValidator(genesis, validator.privateKey, fmt.Sprintf("127.0.0.1:%d", validator.port), validator.rpcPort)
		if err != nil {
			t.Fatal("cant make a node", i, err)
		}
	}

	// start the nodes
	for _, validator := range nodes {
		err := validator.startNode()
		assert.NoError(t, err)
	}
	// Connect everyone to everyone here
	for _, n := range nodes {
		for _, nInner := range nodes {
			n.node.Server().AddPeer(nInner.node.Server().Self())
		}
	}
	return nodes
}

func makeValidator(genesis *core.Genesis, nodekey *ecdsa.PrivateKey, listenAddr string, rpcPort int) (*node.Node, error) { //здесь эта переменная-функция называется cons
	// Define the basic configurations for the Ethereum node
	datadir, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, err
	}

	if listenAddr == "" {
		listenAddr = "0.0.0.0:0"
	}

	configNode := &node.Config{
		Name:    "autonity",
		Version: params.Version,
		DataDir: datadir,
		P2P: p2p.Config{
			ListenAddr:  listenAddr,
			NoDiscovery: true,
			MaxPeers:    25,
			PrivateKey:  nodekey,
		},
		NoUSB: true,
	}
	configNode.HTTPHost = "127.0.0.1"
	configNode.HTTPPort = rpcPort

	// Start the node and configure a full Ethereum node on it
	stack, err := node.New(configNode)
	if err != nil {
		return nil, err
	}
	if err := stack.Register(func(ctx *node.ServiceContext) (node.Service, error) {
		return eth.New(ctx, &eth.Config{
			Genesis:         genesis,
			NetworkId:       genesis.Config.ChainID.Uint64(),
			SyncMode:        downloader.FullSync,
			DatabaseCache:   256,
			DatabaseHandles: 256,
			TxPool:          core.DefaultTxPoolConfig,
			Tendermint:      *config.DefaultConfig(),
		}, nil, nil)
	}); err != nil {
		return nil, err
	}

	// Start the node and return if successful
	return stack, nil
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
		port := strings.Split(listener.Addr().String(), ":")[1]
		node.address = fmt.Sprintf("%s:%s", i, port)
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

		node.url = enode.V4DNSUrl(
			node.privateKey.PublicKey,
			node.address,
			node.port,
			node.port,
		)
	}
}

func makeGenesis(nodes map[string]*testNode) *core.Genesis {
	// generate genesis block
	genesis := core.DefaultGenesisBlock()
	genesis.ExtraData = nil
	genesis.GasLimit = math.MaxUint64 - 1
	genesis.GasUsed = 0
	genesis.Difficulty = big.NewInt(1)
	genesis.Timestamp = 0
	genesis.Nonce = 0
	genesis.Mixhash = types.BFTDigest

	genesis.Config = params.TestChainConfig
	genesis.Config.Tendermint = &params.TendermintConfig{}
	genesis.Config.Ethash = nil
	genesis.Config.AutonityContractConfig = &params.AutonityContractGenesis{}

	genesis.Alloc = core.GenesisAlloc{}
	for _, validator := range nodes {
		genesis.Alloc[crypto.PubkeyToAddress(validator.privateKey.PublicKey)] = core.GenesisAccount{
			Balance: new(big.Int).Exp(big.NewInt(2), big.NewInt(128), nil),
		}
	}

	users := make([]params.User, 0, len(nodes))
	for n, validator := range nodes {
		var nodeType params.UserType
		stake := uint64(100)
		switch {
		case strings.HasPrefix(n, ValidatorPrefix):
			nodeType = params.UserValidator
		case strings.HasPrefix(n, StakeholderPrefix):
			nodeType = params.UserStakeHolder
		case strings.HasPrefix(n, ParticipantPrefix):
			nodeType = params.UserParticipant
			stake = 0
		default:
			panic("incorrect node type")

		}
		users = append(users, params.User{
			Address: crypto.PubkeyToAddress(validator.privateKey.PublicKey),
			Enode:   validator.url,
			Type:    nodeType,
			Stake:   stake,
		})
	}
	//generate one sh
	shKey, err := keygenerator.Next()
	if err != nil {
		log.Error("Make genesis error", "err", err)
	}
	users = append(users, params.User{
		Address: crypto.PubkeyToAddress(shKey.PublicKey),
		Type:    params.UserStakeHolder,
		Stake:   200,
	})
	genesis.Config.AutonityContractConfig.Users = users
	err = genesis.Config.AutonityContractConfig.AddDefault().Validate()
	if err != nil {
		panic(err)
	}

	err = genesis.SetBFT()
	if err != nil {
		panic(err)
	}

	return genesis
}
