package test

import (
	"crypto/ecdsa"
	"crypto/rand"
	"io/ioutil"
	"math/big"
	"net"
	"sync"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/common/math"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/tendermint/config"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/eth"
	"github.com/clearmatics/autonity/eth/downloader"
	"github.com/clearmatics/autonity/event"
	"github.com/clearmatics/autonity/node"
	"github.com/clearmatics/autonity/p2p"
	"github.com/clearmatics/autonity/p2p/enode"
	"github.com/clearmatics/autonity/params"
)

type networkRate struct {
	in  int64
	out int64
}

type testNode struct {
	isRunning      bool
	isInited       bool
	wasStopped     bool  //fixme should be removed
	privateKey     *ecdsa.PrivateKey
	address        string
	port           int
	url            string
	listener       net.Listener
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
}

type block struct {
	hash common.Hash
	txs int
}

func sendTx(service *eth.Ethereum, fromValidator *ecdsa.PrivateKey, fromAddr common.Address, toAddr common.Address) (*types.Transaction, error) {
	nonce := service.TxPool().Nonce(fromAddr)

	tx, err := txWithNonce(fromAddr, nonce, toAddr, fromValidator, service)
	if err != nil {
		return txWithNonce(fromAddr, nonce+1, toAddr, fromValidator, service)
	}
	return tx, nil
}

func txWithNonce(_ common.Address, nonce uint64, toAddr common.Address, fromValidator *ecdsa.PrivateKey, service *eth.Ethereum) (*types.Transaction, error) {
	randEth, err := rand.Int(rand.Reader, big.NewInt(10000000))
	if err != nil {
		return nil, err
	}
	tx, err := types.SignTx(
		types.NewTransaction(
			nonce,
			toAddr,
			big.NewInt(1),
			210000000,
			big.NewInt(100000000000+int64(randEth.Uint64())),
			nil,
		),
		types.HomesteadSigner{}, fromValidator)
	if err != nil {
		return nil, err
	}
	return tx, service.TxPool().AddLocal(tx)
}

func makeGenesis(validators []*testNode) *core.Genesis {
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

	genesis.Alloc = core.GenesisAlloc{}
	for _, validator := range validators {
		genesis.Alloc[crypto.PubkeyToAddress(validator.privateKey.PublicKey)] = core.GenesisAccount{
			Balance: new(big.Int).Exp(big.NewInt(2), big.NewInt(128), nil),
		}
	}

	validatorsAddresses := make([]string, len(validators))
	for i, validator := range validators {
		validatorsAddresses[i] = crypto.PubkeyToAddress(validator.privateKey.PublicKey).String()
	}

	enodes := make([]string, len(validators))
	for i, validator := range validators {
		enodes[i] = validator.url
	}

	genesis.Config.EnodeWhitelist = enodes

	genesis.Validators = validatorsAddresses
	err := genesis.SetBFT()
	if err != nil {
		panic(err)
	}

	return genesis
}

func makeValidator(genesis *core.Genesis, nodekey *ecdsa.PrivateKey, listenAddr string, inRate, outRate int64, cons func(basic consensus.Engine) consensus.Engine) (*node.Node, error) {
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

	if inRate != 0 || outRate != 0 {
		configNode.P2P.IsRated = true
		configNode.P2P.InRate = inRate
		configNode.P2P.OutRate = outRate
	}

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
		}, cons)
	}); err != nil {
		return nil, err
	}

	// Start the node and return if successful
	return stack, nil
}
