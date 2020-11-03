package test

import (
	"crypto/ecdsa"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/clearmatics/autonity/consensus"

	"github.com/clearmatics/autonity/event"
	"github.com/clearmatics/autonity/node"
	"github.com/clearmatics/autonity/p2p/enode"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/eth"
)

type networkRate struct {
	in  int64
	out int64
}

type testNode struct {
	netNode
	isRunning               bool
	isInited                bool
	wasStopped              bool //fixme should be removed
	node                    *node.Node
	nodeConfig              *node.Config
	ethConfig               *eth.Config
	engineConstructor       func(basic consensus.Engine) consensus.Engine
	enode                   *enode.Node
	service                 *eth.Ethereum
	eventChan               chan core.ChainEvent
	subscription            event.Subscription
	transactions            map[common.Hash]struct{}
	transactionsMu          sync.Mutex
	untrustedTransactions   map[common.Hash]struct{}
	untrustedTransactionsMu sync.Mutex
	blocks                  map[uint64]block
	lastBlock               uint64
	txsSendCount            *int64
	txsChainCount           map[uint64]int64
	isMalicious             bool
}

type netNode struct {
	listener   []net.Listener
	privateKey *ecdsa.PrivateKey
	address    string
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

	validator.service, err = eth.New(validator.node, validator.ethConfig, validator.engineConstructor)
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

func (validator *testNode) stopNode() error {
	//remove pending transactions
	addr := crypto.PubkeyToAddress(validator.privateKey.PublicKey)

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		pendingTxsMap, queuedTxsMap := validator.service.TxPool().Content()
		if len(pendingTxsMap) == 0 && len(queuedTxsMap) == 0 {
			break
		}

		canBreak := true
		for txAddr, txs := range pendingTxsMap {
			if addr != txAddr {
				continue
			}
			if len(txs) != 0 {
				canBreak = false
			}
		}
		for txAddr, txs := range queuedTxsMap {
			if addr != txAddr {
				continue
			}
			if len(txs) != 0 {
				canBreak = false
			}
		}
		if canBreak {
			break
		}
	}

	return validator.forceStopNode()
}

func (validator *testNode) forceStopNode() error {
	if err := validator.node.Close(); err != nil {
		return fmt.Errorf("cannot stop a node on block %d: %q", validator.lastBlock, err)
	}
	validator.node.Wait()
	validator.isRunning = false
	validator.wasStopped = true

	return nil
}

func (validator *testNode) startService() error {
	if validator.eventChan == nil {
		validator.eventChan = make(chan core.ChainEvent, 1024)
		validator.transactions = make(map[common.Hash]struct{})
		validator.untrustedTransactions = make(map[common.Hash]struct{})
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

	if err := validator.service.StartMining(1); err != nil {
		return fmt.Errorf("cant start mining %s", err)
	}

	for !validator.service.IsMining() {
		time.Sleep(50 * time.Millisecond)
	}

	validator.isRunning = true

	return nil
}
