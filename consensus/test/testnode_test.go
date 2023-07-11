package test

import (
	"crypto/ecdsa"
	"fmt"
	"net"
	"sync"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/eth"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/node"
	"github.com/autonity/autonity/p2p/enode"
)

type networkRate struct {
	in  int64
	out int64
}

type testNode struct {
	netNode
	isRunning      bool
	node           *node.Node
	nodeConfig     *node.Config
	ethConfig      *eth.Config
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

type netNode struct {
	listener   []net.Listener
	privateKey *ecdsa.PrivateKey
	host       string
	address    common.Address
	port       int
	url        string
	rpcPort    int
}

func (n *netNode) EthAddress() common.Address {
	return crypto.PubkeyToAddress(n.privateKey.PublicKey)
}

type block struct {
	hash common.Hash
	txs  int
}

func (validator *testNode) startNode() error {
	// Start the node and configure a full Ethereum node on it
	var err error
	validator.node, err = node.New(validator.nodeConfig)
	if err != nil {
		return err
	}

	validator.service, err = eth.New(validator.node, validator.ethConfig)
	if err != nil {
		return err
	}

	if err := validator.node.Start(); err != nil {
		return fmt.Errorf("cannot start a node %s", err)
	}

	// Start tracking the node and it's enode
	validator.enode = validator.node.Server().Self()
	return nil
}

func (validator *testNode) startService() error {
	if validator.eventChan == nil {
		validator.eventChan = make(chan core.ChainEvent, 1024)
		validator.transactions = make(map[common.Hash]struct{})
		validator.blocks = make(map[uint64]block)
		validator.txsSendCount = new(int64)
		validator.txsChainCount = make(map[uint64]int64)
	} else {
		// validator is restarting
		// we need to retrieve missed block events since last stop as we're not subscribing fast enough
		curBlock := validator.service.BlockChain().CurrentBlock().Number().Uint64()
		for blockNum := validator.lastBlock + 1; blockNum <= curBlock; blockNum++ {
			block := validator.service.BlockChain().GetBlockByNumber(blockNum)
			event := core.ChainEvent{
				Block: block,
				Hash:  block.Hash(),
				Logs:  nil,
			}
			validator.eventChan <- event
		}
	}

	validator.subscription = validator.service.BlockChain().SubscribeChainEvent(validator.eventChan)
	validator.isRunning = true
	return nil
}
