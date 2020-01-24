package test

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"github.com/clearmatics/autonity/common/keygenerator"
	"io/ioutil"
	"math"
	"math/big"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/tendermint/config"
	tendermintCore "github.com/clearmatics/autonity/consensus/tendermint/core"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/eth"
	"github.com/clearmatics/autonity/eth/downloader"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/node"
	"github.com/clearmatics/autonity/p2p"
	"github.com/clearmatics/autonity/params"
	"golang.org/x/sync/errgroup"
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

func makeValidator(genesis *core.Genesis, nodekey *ecdsa.PrivateKey, listenAddr string, rpcPort int, inRate, outRate int64, cons func(basic consensus.Engine) consensus.Engine, backs func(basic tendermintCore.Backend) tendermintCore.Backend) (*node.Node, error) { //здесь эта переменная-функция называется cons
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
		}, cons, backs)
	}); err != nil {
		return nil, err
	}

	// Start the node and return if successful
	return stack, nil
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

func sendTransactions(t *testing.T, test *testCase, validators map[string]*testNode, txPerPeer int, errorOnTx bool, names []string) {
	const blocksToWait = 15

	txs := make(map[uint64]int) // blockNumber to count
	txsMu := &sync.Mutex{}

	test.validatorsCanBeStopped = new(int64)
	wg, ctx := errgroup.WithContext(context.Background())

	for index, validator := range validators {
		index := index
		validator := validator

		logger := log.New("addr", crypto.PubkeyToAddress(validator.privateKey.PublicKey).String(), "idx", index)

		// skip malicious nodes
		if len(test.maliciousPeers) != 0 {
			if _, ok := test.maliciousPeers[index]; ok {
				atomic.AddInt64(test.validatorsCanBeStopped, 1)
				continue
			}
		}

		wg.Go(func() error {
			return runNode(ctx, validator, test, validators, logger, index, blocksToWait, txs, txsMu, errorOnTx, txPerPeer, names)
		})
	}
	err := wg.Wait()
	keys := make([]int, 0, len(txs))
	for key := range txs {
		keys = append(keys, int(key))
	}
	sort.Ints(keys)

	fmt.Println("Transactions per block")
	for _, key := range keys {
		count := txs[uint64(key)]
		fmt.Printf("Block %d has %d transactions\n", key, count)
	}
	fmt.Println("\nPending transactions")
	for index, validator := range validators {
		validator.transactionsMu.Lock()
		fmt.Printf("Validator %s has %d transactions\n", index, len(validator.transactions))
		validator.transactionsMu.Unlock()
	}
	if err != nil {
		t.Fatal(err)
	}

	// no blocks can be mined with no quorum
	if test.noQuorumAfterBlock > 0 {
		for index, validator := range validators {
			if validator.lastBlock < test.noQuorumAfterBlock-1 {
				t.Fatalf("validator [%s] should have mined blocks. expected block number %d, but got %d",
					index, test.noQuorumAfterBlock-1, validator.lastBlock)
			}

			if validator.lastBlock > test.noQuorumAfterBlock {
				t.Fatalf("validator [%s] mined blocks without quorum. expected block number %d, but got %d",
					index, test.noQuorumAfterBlock, validator.lastBlock)
			}
		}
	}

	minHeight := checkAndReturnMinHeight(t, test, validators)
	checkNodesDontContainMaliciousBlock(t, minHeight, validators, test)
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
		return fmt.Errorf("error while executing before hook for validator index %s and block %v, err %v",
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

func runNode(ctx context.Context, validator *testNode, test *testCase, validators map[string]*testNode, logger log.Logger, index string, blocksToWait int, txs map[uint64]int, txsMu sync.Locker, errorOnTx bool, txPerPeer int, names []string) error {
	var err error
	testCanBeStopped := new(uint32)
	fromAddr := crypto.PubkeyToAddress(validator.privateKey.PublicKey)

wgLoop:
	for {
		select {
		case ev := <-validator.eventChan:
			if _, ok := validator.blocks[ev.Block.NumberU64()]; ok {
				continue
			}

			// before hook
			err = runHook(test.getBeforeHook(index), test, ev.Block, validator, index)
			if err != nil {
				return err
			}

			validator.blocks[ev.Block.NumberU64()] = block{ev.Block.Hash(), len(ev.Block.Transactions())}
			validator.lastBlock = ev.Block.NumberU64()

			logger.Error("last mined block", "validator", index,
				"num", validator.lastBlock, "hash", validator.blocks[ev.Block.NumberU64()].hash,
				"txCount", validator.blocks[ev.Block.NumberU64()].txs)

			if atomic.LoadUint32(testCanBeStopped) == 1 {
				if atomic.LoadInt64(test.validatorsCanBeStopped) == int64(len(validators)) {
					break wgLoop
				}
				if atomic.LoadInt64(test.validatorsCanBeStopped) > int64(len(validators)) {
					return fmt.Errorf("something is wrong. %d of %d validators are ready to be stopped", atomic.LoadInt64(test.validatorsCanBeStopped), uint32(len(validators)))
				}
				continue
			}

			// actual forming and sending transaction
			logger.Debug("peer", "address", crypto.PubkeyToAddress(validator.privateKey.PublicKey).String(), "block", ev.Block.Number().Uint64(), "isRunning", validator.isRunning)

			if validator.isRunning {
				txsMu.Lock()
				if _, ok := txs[validator.lastBlock]; !ok {
					txs[validator.lastBlock] = ev.Block.Transactions().Len()
				}
				txsMu.Unlock()

				for _, tx := range ev.Block.Transactions() {
					validator.transactionsMu.Lock()
					if _, ok := validator.transactions[tx.Hash()]; ok {
						validator.txsChainCount[ev.Block.NumberU64()]++
						delete(validator.transactions, tx.Hash())
					}
					validator.transactionsMu.Unlock()
				}

				if int(validator.lastBlock) <= test.numBlocks {
					err = validatorSendTransaction(
						generateToAddr(txPerPeer, names, index, validators),
						test,
						validator)
					if err != nil {
						return err
					}
				}
			}

			// after hook
			err = runHook(test.getAfterHook(index), test, ev.Block, validator, index)
			if err != nil {
				return err
			}

			if int(validator.lastBlock) > test.numBlocks {
				//all transactions were included into the chain
				if errorOnTx {
					validator.transactionsMu.Lock()
					if len(validator.transactions) == 0 {
						if atomic.CompareAndSwapUint32(testCanBeStopped, 0, 1) {
							atomic.AddInt64(test.validatorsCanBeStopped, 1)
						}
					}
					validator.transactionsMu.Unlock()
				} else {

					if atomic.CompareAndSwapUint32(testCanBeStopped, 0, 1) {
						atomic.AddInt64(test.validatorsCanBeStopped, 1)
					}
				}
			}

			if validator.isRunning && int(validator.lastBlock) >= test.numBlocks+blocksToWait {
				if errorOnTx {
					pending, queued := validator.service.TxPool().Stats()
					if pending > 0 {
						return fmt.Errorf("after a new block it should be 0 pending transactions got %d. block %d", pending, ev.Block.Number().Uint64())
					}
					if queued > 0 {
						return fmt.Errorf("after a new block it should be 0 queued transactions got %d. block %d", queued, ev.Block.Number().Uint64())
					}

					validator.transactionsMu.Lock()
					pendingTransactions := len(validator.transactions)
					havePendingTransactions := pendingTransactions != 0
					validator.transactionsMu.Unlock()

					if havePendingTransactions {
						var txsChainCount int64
						for _, txsBlockCount := range validator.txsChainCount {
							txsChainCount += txsBlockCount
						}

						if validator.wasStopped {
							//fixme an error should be returned
							logger.Error("test error!!!", "err", fmt.Errorf("a validator %s still have transactions to be mined %d. block %d. Total sent %d, total mined %d",
								index,
								pendingTransactions, ev.Block.Number().Uint64(),
								atomic.LoadInt64(validator.txsSendCount), txsChainCount))

							if atomic.CompareAndSwapUint32(testCanBeStopped, 0, 1) {
								atomic.AddInt64(test.validatorsCanBeStopped, 1)
							}
						} else {
							return fmt.Errorf("a validator %s still have transactions to be mined %d. block %d. Total sent %d, total mined %d",
								index,
								pendingTransactions, ev.Block.Number().Uint64(),
								atomic.LoadInt64(validator.txsSendCount), txsChainCount)
						}
					}
				}
			}
		case innerErr := <-validator.subscription.Err():
			if innerErr != nil {
				return fmt.Errorf("error in blockchain %q", innerErr)
			}

			time.Sleep(500 * time.Millisecond)

			// after hook
			err = runHook(test.getAfterHook(index), test, nil, validator, index)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
		// allow to exit goroutine when no quorum expected
		// check that there's no quorum within the given noQuorumTimeout
		if !hasQuorum(validators) && test.noQuorumAfterBlock > 0 {
			log.Error("No Quorum", "index", index, "last_block", validator.lastBlock)
			// wait for quorum to get restored
			time.Sleep(test.noQuorumTimeout)
			break wgLoop
		}

		// check transactions status if all blocks are passed
		txRemoveBlock, ok := test.removedPeers[fromAddr]
		if ok && (validator.lastBlock >= txRemoveBlock) {
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

func validatorSendTransaction(addresses []addressesList, test *testCase, validator *testNode) error {
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
			if tx != nil {
				atomic.AddInt64(validator.txsSendCount, 1)
			} else if skip {
				if tx, innerErr = sendTx(validator.service, validator.privateKey, fromAddr, toAddr.Address, generateRandomTx); innerErr != nil {
					return innerErr
				}
				atomic.AddInt64(validator.txsSendCount, 1)
			}

		} else {
			if tx, innerErr = sendTx(validator.service, validator.privateKey, fromAddr, toAddr.Address, generateRandomTx); innerErr != nil {
				return innerErr
			}
			atomic.AddInt64(validator.txsSendCount, 1)
		}

		validator.transactionsMu.Lock()
		if tx != nil {
			validator.transactions[tx.Hash()] = struct{}{}
		}
		validator.transactionsMu.Unlock()
	}
	return nil
}

func checkNodesDontContainMaliciousBlock(t *testing.T, minHeight uint64, validators map[string]*testNode, test *testCase) {
	// check that all nodes got the same blocks
	for i := uint64(1); i <= minHeight; i++ {
		blockHash := validators["VA"].blocks[i].hash

		for index, validator := range validators {
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
