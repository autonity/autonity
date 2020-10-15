package test

import (
	"crypto/ecdsa"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/clearmatics/autonity/event"
	"github.com/clearmatics/autonity/node"
	"github.com/clearmatics/autonity/p2p/enode"

	"github.com/clearmatics/autonity/accounts"
	"github.com/clearmatics/autonity/accounts/keystore"
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

type block struct {
	hash common.Hash
	txs  int
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
	var ethereum *eth.Ethereum
	if err := validator.node.Start(); err != nil {
		return fmt.Errorf("cant start a node %s", err)
	}

	time.Sleep(100 * time.Millisecond)

	validator.service = ethereum

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

	if err := ethereum.StartMining(1); err != nil {
		return fmt.Errorf("cant start mining %s", err)
	}

	for !ethereum.IsMining() {
		time.Sleep(50 * time.Millisecond)
	}

	validator.isRunning = true

	return nil
}
