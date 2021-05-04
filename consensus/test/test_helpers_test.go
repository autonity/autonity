package test

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"golang.org/x/sync/errgroup"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/common/keygenerator"
	"github.com/clearmatics/autonity/consensus/tendermint/config"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/eth"
	"github.com/clearmatics/autonity/eth/downloader"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/node"
	"github.com/clearmatics/autonity/p2p"
	"github.com/clearmatics/autonity/params"
)

func sendTx(service *eth.Ethereum, key *ecdsa.PrivateKey, fromAddr common.Address, toAddr common.Address, transactionGenerator func(nonce uint64, toAddr common.Address, key *ecdsa.PrivateKey) (*types.Transaction, error)) (*types.Transaction, error) {
	nonce := service.TxPool().Nonce(fromAddr)

	var tx *types.Transaction
	var err error

	for stop := 10; stop > 0; stop-- {
		tx, err = transactionGenerator(nonce, toAddr, key)
		if err != nil {
			nonce++
			continue
		}
		err = service.TxPool().AddLocal(tx)
		if err == nil {
			break
		}
		nonce++

	}
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func generateRandomTx(nonce uint64, toAddr common.Address, key *ecdsa.PrivateKey) (*types.Transaction, error) {
	randEth, err := rand.Int(rand.Reader, big.NewInt(10000000))
	if err != nil {
		return nil, err
	}

	return types.SignTx(
		types.NewTransaction(
			nonce,
			toAddr,
			big.NewInt(1),
			210000000,
			big.NewInt(100000000000+int64(randEth.Uint64())),
			nil,
		),
		types.HomesteadSigner{}, key)
}

func makeGenesis(t *testing.T, nodes map[string]*testNode, stakeholderName string) *core.Genesis {
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
	genesis.Config.Ethash = nil
	genesis.Config.AutonityContractConfig = &params.AutonityContractGenesis{}

	genesis.Alloc = core.GenesisAlloc{}
	for _, validator := range nodes {
		genesis.Alloc[crypto.PubkeyToAddress(validator.privateKey.PublicKey)] = core.GenesisAccount{
			Balance: new(big.Int).Exp(big.NewInt(2), big.NewInt(128), nil),
		}
	}

	users := make([]params.User, 0, len(nodes))
	for n, node := range nodes {
		var nodeType params.UserType
		stake := uint64(100)

		var skip bool
		switch {
		case strings.HasPrefix(n, ValidatorPrefix):
			nodeType = params.UserValidator
		case strings.HasPrefix(n, StakeholderPrefix):
			nodeType = params.UserStakeHolder
		case strings.HasPrefix(n, ParticipantPrefix):
			nodeType = params.UserParticipant
			stake = 0
		case strings.HasPrefix(n, ExternalPrefix):
			//an unknown user
			skip = true
		default:
			require.FailNow(t, "incorrect node type")
		}

		if skip {
			continue
		}
		address := crypto.PubkeyToAddress(node.privateKey.PublicKey)
		users = append(users, params.User{
			Address: &address,
			Enode:   node.url,
			Type:    nodeType,
			Stake:   stake,
		})
	}

	//generate one sh
	shKey, err := keygenerator.Next()
	if err != nil {
		log.Error("Make genesis error", "err", err)
	}

	stakeNode, err := newNode(shKey, stakeholderName)
	if err != nil {
		log.Error("Make genesis error while adding a stakeholder", "err", err)
	}

	address := crypto.PubkeyToAddress(shKey.PublicKey)
	stakeHolder := params.User{
		Address: &address,
		Type:    params.UserStakeHolder,
		Stake:   200,
	}

	stakeHolder.Enode = stakeNode.url

	users = append(users, stakeHolder)
	genesis.Config.AutonityContractConfig.Users = users
	err = genesis.Config.AutonityContractConfig.Prepare()
	require.NoError(t, err)
	return genesis
}

func makeNodeConfig(t *testing.T, genesis *core.Genesis, nodekey *ecdsa.PrivateKey, listenAddr string, rpcPort int, inRate, outRate int64) (*node.Config, *eth.Config) {
	// Define the basic configurations for the Ethereum node
	datadir, err := ioutil.TempDir("", "")
	require.NoError(t, err)

	if listenAddr == "" {
		listenAddr = "0.0.0.0:0"
	}

	configNode := &node.Config{
		Name:    "autonity",
		Version: params.Version,
		DataDir: datadir,
		P2P: p2p.Config{
			ListenAddr:            listenAddr,
			NoDiscovery:           true,
			MaxPeers:              25,
			PrivateKey:            nodekey,
			DialHistoryExpiration: time.Millisecond,
		},
		NoUSB: true,
	}
	configNode.HTTPHost = "127.0.0.1"
	configNode.HTTPPort = rpcPort

	if inRate != 0 || outRate != 0 {
		configNode.P2P.IsRated = true
		configNode.P2P.InRate = inRate
		configNode.P2P.OutRate = outRate
	}

	ethConfig := &eth.Config{
		Genesis:         genesis,
		NetworkId:       genesis.Config.ChainID.Uint64(),
		SyncMode:        downloader.FullSync,
		DatabaseCache:   256,
		DatabaseHandles: 256,
		TxPool:          core.DefaultTxPoolConfig,
	}
	return configNode, ethConfig
}

func maliciousTest(t *testing.T, test *testCase, validators map[string]*testNode) {
	for index, validator := range validators {
		for number, block := range validator.blocks {
			if test.addedValidatorsBlocks != nil {
				if maliciousBlock, ok := test.addedValidatorsBlocks[block.hash]; ok {
					t.Errorf("a malicious block %d(%v)\nwas added to %s(%v)", number, maliciousBlock, index, validator)
				}
			}
		}
	}
}

func sendTransactions(t *testing.T, test *testCase, peers map[string]*testNode, txPerPeer int, errorOnTx bool, names []string) {
	const blocksToWait = 100

	txs := make(map[uint64]int) // blockNumber to count
	txsMu := &sync.Mutex{}

	test.validatorsCanBeStopped = new(int64)
	wg, ctx := errgroup.WithContext(context.Background())

	for index, peer := range peers {
		index := index
		peer := peer

		logger := log.New("addr", crypto.PubkeyToAddress(peer.privateKey.PublicKey).String(), "idx", index)

		// skip malicious nodes
		if len(test.maliciousPeers) != 0 {
			if _, ok := test.maliciousPeers[index]; ok {
				atomic.AddInt64(test.validatorsCanBeStopped, 1)
				continue
			}
		}

		wg.Go(func() error {
			return runNode(ctx, peer, test, peers, logger, index, blocksToWait, txs, txsMu, errorOnTx, txPerPeer, names)
		})
	}
	err := wg.Wait()
	if err != nil {
		if test.topology != nil {
			fmt.Println(test.topology.DumpTopology(peers))
		}
		t.Fatal(err)
	}

	keys := make([]int, 0, len(txs))
	for key := range txs {
		keys = append(keys, int(key))
	}
	sort.Ints(keys)

	for _, key := range keys {
		count := txs[uint64(key)]
		fmt.Printf("Block %d has %d transactions\n", key, count)
	}

	for index, peer := range peers {
		peer.transactionsMu.Lock()
		fmt.Printf("Validator %s has %d transactions\n", index, len(peer.transactions))
		peer.transactionsMu.Unlock()
	}

	// no blocks can be mined with no quorum
	if test.noQuorumAfterBlock > 0 {
		for index, peer := range peers {
			if peer.lastBlock < test.noQuorumAfterBlock-1 {
				t.Fatalf("peer [%s] should have mined blocks. expected block number %d, but got %d",
					index, test.noQuorumAfterBlock-1, peer.lastBlock)
			}

			if peer.lastBlock > test.noQuorumAfterBlock {
				t.Fatalf("peer [%s] mined blocks without quorum. expected block number %d, but got %d",
					index, test.noQuorumAfterBlock, peer.lastBlock)
			}
		}
	}

	minHeight := checkAndReturnMinHeight(t, test, peers)
	checkNodesDontContainMaliciousBlock(t, minHeight, peers, test)
	fmt.Println("\nTransactions OK")
}

func hasQuorum(validators map[string]*testNode) bool {
	active := 0
	for _, val := range validators {
		if val.isRunning {
			active++
		}
	}
	return quorum(len(validators), active)
}

func quorum(valCount, activeVals int) bool {
	return float64(activeVals) >= math.Ceil(float64(2)/float64(3)*float64(valCount))
}

func runHook(validatorHook hook, test *testCase, block *types.Block, validator *testNode, index string) error {
	if validatorHook == nil {
		return nil
	}

	err := validatorHook(block, validator, test, time.Now())
	if err != nil {
		return fmt.Errorf("error while executing hook for validator index %s and block %v, err %v",
			index, block.NumberU64(), err)
	}

	return nil
}

func hookStopNode(nodeIndex string, blockNum uint64) hook {
	return func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error {
		if block.Number().Uint64() == blockNum {
			err := validator.stopNode()
			if err != nil {
				return err
			}

			tCase.setStopTime(nodeIndex, currentTime)
		}

		return nil
	}
}

func hookForceStopNode(nodeIndex string, blockNum uint64) hook {
	return func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error {
		if block.Number().Uint64() == blockNum {
			err := validator.forceStopNode()
			if err != nil {
				return err
			}
			tCase.setStopTime(nodeIndex, currentTime)
		}
		return nil
	}
}

func hookStartNode(nodeIndex string, durationAfterStop float64) hook {
	return func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error {
		stopTime := tCase.getStopTime(nodeIndex)
		if block == nil && currentTime.Sub(stopTime).Seconds() >= durationAfterStop {
			if err := validator.startNode(); err != nil {
				return err
			}

			if err := validator.startService(); err != nil {
				return err
			}
		}

		return nil
	}
}

func runNode(ctx context.Context, peer *testNode, test *testCase, peers map[string]*testNode, logger log.Logger, index string, blocksToWait int, txs map[uint64]int, txsMu sync.Locker, errorOnTx bool, txPerPeer int, names []string) error {
	var err error
	testCanBeStopped := new(uint32)
	fromAddr := crypto.PubkeyToAddress(peer.privateKey.PublicKey)

	var noQuorumTimer *time.Timer
	if test.noQuorumAfterBlock > 0 {
		noQuorumTimer = time.NewTimer(test.noQuorumTimeout)
		defer noQuorumTimer.Stop()
	}
	periodicChecks := time.NewTicker(100 * time.Millisecond)
	defer periodicChecks.Stop()

	mux := peer.node.EventMux()
	chainEvents := mux.Subscribe(downloader.StartEvent{}, downloader.DoneEvent{})
	defer chainEvents.Unsubscribe()

	shouldSendTx := peer.service.Miner().Mining()

	isExternalUser := isExternalUser(index)
	if isExternalUser {
		atomic.AddInt64(test.validatorsCanBeStopped, 1)
	}

wgLoop:
	for {
		select {
		case ev := <-peer.eventChan:
			err = peer.service.APIBackend.IsSelfInWhitelist()
			if !isExternalUser && err != nil {
				return fmt.Errorf("a user %q should be in the whitelist: %v on block %d", index, err, ev.Block.NumberU64())
			}
			if isExternalUser && err == nil {
				return fmt.Errorf("a user %q shoulnd't be in the whitelist on block %d", index, ev.Block.NumberU64())
			}

			if test.topology != nil && test.topology.WithChanges() {
				err = test.topology.ConnectNodesForIndex(index, peers)
				if err != nil {
					return err
				}
			}

			if _, ok := peer.blocks[ev.Block.NumberU64()]; ok {
				continue
			}

			// before hook
			err = runHook(test.getBeforeHook(index), test, ev.Block, peer, index)
			if err != nil {
				return err
			}

			peer.blocks[ev.Block.NumberU64()] = block{ev.Block.Hash(), len(ev.Block.Transactions())}
			peer.lastBlock = ev.Block.NumberU64()

			logger.Error("last mined block", "peer", index,
				"num", peer.lastBlock, "hash", peer.blocks[ev.Block.NumberU64()].hash,
				"txCount", peer.blocks[ev.Block.NumberU64()].txs)

			if atomic.LoadUint32(testCanBeStopped) == 1 {
				if atomic.LoadInt64(test.validatorsCanBeStopped) == int64(len(peers)) {
					break wgLoop
				}
				if atomic.LoadInt64(test.validatorsCanBeStopped) > int64(len(peers)) {
					return fmt.Errorf("something is wrong. %d of %d peers are ready to be stopped", atomic.LoadInt64(test.validatorsCanBeStopped), uint32(len(peers)))
				}
				continue
			}

			// actual forming and sending transaction
			logger.Debug("peer", "address", crypto.PubkeyToAddress(peer.privateKey.PublicKey).String(), "block", ev.Block.Number().Uint64(), "isRunning", peer.isRunning)

			if peer.isRunning {
				txsMu.Lock()
				if _, ok := txs[peer.lastBlock]; !ok {
					txs[peer.lastBlock] = ev.Block.Transactions().Len()
				}
				txsMu.Unlock()

				for _, tx := range ev.Block.Transactions() {
					peer.transactionsMu.Lock()
					if _, ok := peer.transactions[tx.Hash()]; ok {
						peer.txsChainCount[ev.Block.NumberU64()]++
						delete(peer.transactions, tx.Hash())
					}
					peer.transactionsMu.Unlock()
				}

				currentBlock := peer.service.BlockChain().CurrentHeader().Number.Uint64()
				isBehind := currentBlock < ev.Block.NumberU64()
				if isExternalUser {
					if currentBlock > 1 {
						return fmt.Errorf("external user %v got a block %d, topology %v",
							index, currentBlock, test.topology.DumpTopology(peers))
					}
				}

				var sendTx func(addresses []addressesList, test *testCase, peer *testNode) error
				if !isExternalUser && !isBehind && shouldSendTx && int(peer.lastBlock) <= test.numBlocks {
					sendTx = peerSendTransaction
				} else if isExternalUser {
					sendTx = peerSendExternalTransaction
				}
				if sendTx != nil {
					err = sendTx(
						generateToAddr(txPerPeer, names, index, peers),
						test,
						peer)
					if err != nil {
						return err
					}
				}
			}

			// after hook
			err = runHook(test.getAfterHook(index), test, ev.Block, peer, index)
			if err != nil {
				return err
			}

			if test.topology != nil && test.topology.WithChanges() {
				err = test.topology.CheckTopologyForIndex(index, peers)
				if err != nil {
					logger.Error("check topology err", "index", index, "block", peer.lastBlock, "err", err)
					return err
				}
			}

			if int(peer.lastBlock) > test.numBlocks {
				// all transactions were included into the chain
				if errorOnTx && !isExternalUser {
					peer.transactionsMu.Lock()
					if len(peer.transactions) == 0 {
						if atomic.CompareAndSwapUint32(testCanBeStopped, 0, 1) {
							atomic.AddInt64(test.validatorsCanBeStopped, 1)
						}
					}
					peer.transactionsMu.Unlock()
				} else {
					if atomic.CompareAndSwapUint32(testCanBeStopped, 0, 1) {
						atomic.AddInt64(test.validatorsCanBeStopped, 1)
					}
				}
			}

			if !isExternalUser && peer.isRunning && int(peer.lastBlock) >= test.numBlocks+blocksToWait {
				if errorOnTx {
					pending, queued := peer.service.TxPool().Stats()
					if pending > 0 {
						return fmt.Errorf("after a new block it should be 0 pending transactions got %d. block %d", pending, ev.Block.Number().Uint64())
					}
					if queued > 0 {
						return fmt.Errorf("after a new block it should be 0 queued transactions got %d. block %d", queued, ev.Block.Number().Uint64())
					}

					peer.transactionsMu.Lock()
					pendingTransactions := len(peer.transactions)
					havePendingTransactions := pendingTransactions != 0
					peer.transactionsMu.Unlock()

					if havePendingTransactions {
						var txsChainCount int64
						for _, txsBlockCount := range peer.txsChainCount {
							txsChainCount += txsBlockCount
						}

						if peer.wasStopped {
							// fixme an error should be returned
							logger.Error("test error!!!", "err", fmt.Errorf("a peer %s still have transactions to be mined %d. block %d. Total sent %d, total mined %d",
								index,
								pendingTransactions, ev.Block.Number().Uint64(),
								atomic.LoadInt64(peer.txsSendCount), txsChainCount))

							if atomic.CompareAndSwapUint32(testCanBeStopped, 0, 1) {
								atomic.AddInt64(test.validatorsCanBeStopped, 1)
							}
						} else {
							return fmt.Errorf("a peer %s still have transactions to be mined %d. block %d. Total sent %d, total mined %d",
								index,
								pendingTransactions, ev.Block.Number().Uint64(),
								atomic.LoadInt64(peer.txsSendCount), txsChainCount)
						}
					}
				}
			}

			if isExternalUser {
				return fmt.Errorf("external user %v got a block %d, topology %v",
					index, ev.Block.NumberU64(), test.topology.DumpTopology(peers))
			}
		case innerErr := <-peer.subscription.Err():
			if innerErr != nil {
				return fmt.Errorf("error in blockchain %q", innerErr)
			}

			time.Sleep(500 * time.Millisecond)

			// after hook
			err = runHook(test.getAfterHook(index), test, nil, peer, index)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		case <-periodicChecks.C:
			if test.noQuorumAfterBlock > 0 {
				if hasQuorum(peers) {
					if !noQuorumTimer.Stop() {
						<-noQuorumTimer.C
					}
					noQuorumTimer.Reset(test.noQuorumTimeout)
				} else {
					select {
					case <-noQuorumTimer.C:
						log.Error("No Quorum", "index", index, "last_block", peer.lastBlock)
						atomic.AddInt64(test.validatorsCanBeStopped, 1)
						break wgLoop
					default:
					}
				}
			}

			if isExternalUser {
				if atomic.LoadInt64(test.validatorsCanBeStopped) == int64(len(peers)) {
					break wgLoop
				}
			}
		case ev := <-chainEvents.Chan():
			if ev == nil {
				continue
			}
			switch ev.Data.(type) {
			case downloader.StartEvent:
				shouldSendTx = false

			case downloader.DoneEvent:
				shouldSendTx = true
			}
		}

		// check transactions status if all blocks are passed
		txRemoveBlock, ok := test.removedPeers[fromAddr]
		if ok && (peer.lastBlock >= txRemoveBlock) {
			if atomic.CompareAndSwapUint32(testCanBeStopped, 0, 1) {
				atomic.AddInt64(test.validatorsCanBeStopped, 1)
				break wgLoop
			}
		}
	}
	return nil
}

func checkAndReturnMinHeight(t *testing.T, test *testCase, validators map[string]*testNode) uint64 {
	// check that all nodes reached the same minimum blockchain height
	minHeight := uint64(math.MaxUint64)
	for index, validator := range validators {
		if isExternalUser(index) {
			continue
		}
		if len(test.maliciousPeers) != 0 {
			if _, ok := test.maliciousPeers[index]; ok {
				// don't check chain for malicious peers
				continue
			}
		}

		validatorBlock := validator.lastBlock
		if minHeight > validatorBlock {
			minHeight = validatorBlock
		}

		if test.noQuorumAfterBlock > 0 {
			continue
		}

		if _, ok := test.removedPeers[crypto.PubkeyToAddress(validator.privateKey.PublicKey)]; ok {
			continue
		}

		if validatorBlock < uint64(test.numBlocks) {
			t.Fatalf("a validator is behind the network index %s and block %v - expected %d",
				index, validatorBlock, test.numBlocks)
		}
	}
	return minHeight
}

type addressesList struct {
	Address   common.Address
	NodeIndex string
}

func generateToAddr(txPerPeer int, names []string, index string, validators map[string]*testNode) []addressesList {
	addresses := make([]addressesList, 0, txPerPeer)
	for i := 0; i < txPerPeer; i++ {
		nextValidatorIndex := names[(sort.SearchStrings(names, index)+i+1)%len(names)]
		toAddr := crypto.PubkeyToAddress(validators[nextValidatorIndex].privateKey.PublicKey)
		addresses = append(addresses, addressesList{
			Address:   toAddr,
			NodeIndex: index,
		})
	}
	return addresses
}

func peerSendTransaction(addresses []addressesList, test *testCase, peer *testNode) error {
	fromAddr := crypto.PubkeyToAddress(peer.privateKey.PublicKey)
	for _, toAddr := range addresses {
		var tx *types.Transaction
		var innerErr error
		var skip bool

		if f, ok := test.sendTransactionHooks[toAddr.NodeIndex]; ok {
			skip, tx, innerErr = f(peer, fromAddr, toAddr.Address)
			if innerErr != nil {
				return innerErr
			}
			if tx != nil {
				atomic.AddInt64(peer.txsSendCount, 1)
			} else if skip {
				if tx, innerErr = sendTx(peer.service, peer.privateKey, fromAddr, toAddr.Address, generateRandomTx); innerErr != nil {
					return innerErr
				}
				atomic.AddInt64(peer.txsSendCount, 1)
			}

		} else {
			if tx, innerErr = sendTx(peer.service, peer.privateKey, fromAddr, toAddr.Address, generateRandomTx); innerErr != nil {
				return innerErr
			}
			atomic.AddInt64(peer.txsSendCount, 1)
		}

		peer.transactionsMu.Lock()
		if tx != nil {
			peer.transactions[tx.Hash()] = struct{}{}
		}
		peer.transactionsMu.Unlock()
	}
	return nil
}

func peerSendExternalTransaction(addresses []addressesList, test *testCase, validator *testNode) error {
	fromAddr := crypto.PubkeyToAddress(validator.privateKey.PublicKey)
	for _, toAddr := range addresses {
		var tx *types.Transaction
		var innerErr error
		var skip bool

		if f, ok := test.sendTransactionHooks[toAddr.NodeIndex]; ok {
			skip, tx, innerErr = f(validator, fromAddr, toAddr.Address)
			if innerErr != nil {
				return innerErr
			}
			if skip {
				if tx, innerErr = sendTx(validator.service, validator.privateKey, fromAddr, toAddr.Address, generateRandomTx); innerErr != nil {
					return innerErr
				}
			}

		} else {
			if tx, innerErr = sendTx(validator.service, validator.privateKey, fromAddr, toAddr.Address, generateRandomTx); innerErr != nil {
				return innerErr
			}
		}

		validator.untrustedTransactionsMu.Lock()
		if tx != nil {
			validator.untrustedTransactions[tx.Hash()] = struct{}{}
		}
		validator.untrustedTransactionsMu.Unlock()
	}
	return nil
}

func checkNodesDontContainMaliciousBlock(t *testing.T, minHeight uint64, validators map[string]*testNode, test *testCase) {
	// check that all nodes got the same blocks
	for i := uint64(1); i <= minHeight; i++ {
		blockHash := validators["VA"].blocks[i].hash

		for index, validator := range validators {
			if isExternalUser(index) {
				continue
			}
			if validator.isMalicious {
				continue
			}

			if len(test.maliciousPeers) != 0 {
				if _, ok := test.maliciousPeers[index]; ok {
					// don't check chain for malicious peers
					continue
				}
			}
			if validator.blocks[i].hash != blockHash {
				t.Fatalf("validators %d and %s have different blocks %d - %q vs %s",
					0, index, i, validator.blocks[i].hash.String(), blockHash.String())
			}
		}
	}
}

func isExternalUser(index string) bool {
	return strings.HasPrefix(index, "E")
}
