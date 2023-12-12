package test

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/autonity/autonity/crypto/blst"
	"math"
	"math/big"
	"os"
	"sort"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/eth/downloader"
	"github.com/autonity/autonity/eth/ethconfig"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/miner"
	"github.com/autonity/autonity/node"
	"github.com/autonity/autonity/p2p"
	"github.com/autonity/autonity/params"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

const defaultStake = 100

func makeGenesis(t *testing.T, nodes map[string]*testNode, names []string) *core.Genesis {
	// generate genesis block
	genesis := core.DefaultGenesisBlock()
	genesis.ExtraData = nil
	genesis.GasLimit = 10000000000
	genesis.GasUsed = 0
	genesis.Timestamp = 0
	genesis.Nonce = 0
	genesis.Mixhash = types.BFTDigest

	genesis.Config = params.TestChainConfig
	genesis.Config.Ethash = nil
	genesis.Config.AutonityContractConfig.Validators = nil
	genesis.Config.AutonityContractConfig.MaxCommitteeSize = 21

	genesis.Alloc = core.GenesisAlloc{}
	for _, validator := range nodes {
		genesis.Alloc[crypto.PubkeyToAddress(validator.privateKey.PublicKey)] = core.GenesisAccount{
			Balance: new(big.Int).Exp(big.NewInt(2), big.NewInt(128), nil),
		}
	}

	validators := make([]*params.Validator, 0, len(nodes))
	for _, name := range names {
		stake := big.NewInt(defaultStake)
		//stake := new(big.Int).Exp(big.NewInt(10), big.NewInt(50), nil)
		if strings.HasPrefix(name, ValidatorPrefix) {
			nodeAddr := crypto.PubkeyToAddress(nodes[name].privateKey.PublicKey)
			blsPK, err := blst.RandKey()
			require.NoError(t, err)
			treasury := nodeAddr
			oracleKey, err := crypto.GenerateKey()
			require.NoError(t, err)
			pop, err := crypto.AutonityPOPProof(nodes[name].privateKey, oracleKey, treasury.Hex(), blsPK)
			require.NoError(t, err)
			validators = append(validators, &params.Validator{
				NodeAddress:    &nodeAddr,
				OracleAddress:  crypto.PubkeyToAddress(oracleKey.PublicKey),
				POP:            pop,
				Enode:          nodes[name].url,
				Treasury:       nodeAddr,
				BondedStake:    stake,
				Key:            blsPK.PublicKey().Marshal(),
				CommissionRate: new(big.Int).SetUint64(0),
			})
		}
	}

	genesis.Config.AutonityContractConfig.Validators = validators
	err := genesis.Config.AutonityContractConfig.Prepare()
	require.NoError(t, err)
	return genesis
}

func makeNodeConfig(t *testing.T, genesis *core.Genesis, nodekey *ecdsa.PrivateKey, listenAddr string, rpcPort int, inRate, outRate int64) (*node.Config, *ethconfig.Config) {
	// Define the basic configurations for the Ethereum node
	datadir, err := os.MkdirTemp("", "")
	require.NoError(t, err)

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
	}
	configNode.HTTPHost = "127.0.0.1"
	configNode.HTTPPort = rpcPort

	if inRate != 0 || outRate != 0 {
		configNode.P2P.IsRated = true
		configNode.P2P.InRate = inRate
		configNode.P2P.OutRate = outRate
	}

	ethConfig := &ethconfig.Config{
		Genesis:         genesis,
		NetworkID:       genesis.Config.ChainID.Uint64(),
		SyncMode:        downloader.FullSync,
		DatabaseCache:   256,
		DatabaseHandles: 256,
		TxPool:          core.DefaultTxPoolConfig,
		Miner:           miner.Config{Etherbase: crypto.PubkeyToAddress(nodekey.PublicKey)},
	}
	return configNode, ethConfig
}

func startTestControllers(t *testing.T, test *testCase, peers map[string]*testNode, errorOnTx bool) {
	const blocksToWait = 100

	txs := make(map[uint64]int) // blockNumber to count

	test.validatorsCanBeStopped = new(int64)
	wg, ctx := errgroup.WithContext(context.Background())

	for index, peer := range peers {
		index := index
		peer := peer

		logger := log.New("addr", crypto.PubkeyToAddress(peer.privateKey.PublicKey).String(), "idx", index)

		wg.Go(func() error {
			return runNode(ctx, peer, test, peers, logger, index, blocksToWait, errorOnTx)
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
		fmt.Printf("block %d has %d transactions\n", key, count)
	}

	for index, peer := range peers {
		peer.transactionsMu.Lock()
		fmt.Printf("Validator %s has %d transactions\n", index, len(peer.transactions))
		peer.transactionsMu.Unlock()
	}

	minHeight := checkAndReturnMinHeight(t, test, peers)
	checkBlockConsistence(t, minHeight, peers, test)
	fmt.Println("\nTransactions OK")
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

func runNode(ctx context.Context, peer *testNode, test *testCase, peers map[string]*testNode, logger log.Logger, index string, blocksToWait int, errorOnTx bool) error {
	var err error
	testCanBeStopped := new(uint32)

	periodicChecks := time.NewTicker(100 * time.Millisecond)
	defer periodicChecks.Stop()

	isExternalUser := isExternalUser(index)
	if isExternalUser {
		atomic.AddInt64(test.validatorsCanBeStopped, 1)
	}

wgLoop:
	for {
		select {
		case ev := <-peer.eventChan:

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

			logger.Info("last mined block", "peer", index,
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
					// log warnings of leftover transactions rather than to fail the test since we have test case to verify TX get mined.
					if pending > 0 {
						log.Warn("after a new block it should be 0 pending transactions got %d. block %d", pending, ev.Block.Number().Uint64())
					}
					if queued > 0 {
						log.Warn("after a new block it should be 0 queued transactions got %d. block %d", queued, ev.Block.Number().Uint64())
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
						log.Warn("a peer %s still have transactions to be mined %d. block %d. Total sent %d, total mined %d",
							index,
							pendingTransactions, ev.Block.Number().Uint64(),
							atomic.LoadInt64(peer.txsSendCount), txsChainCount)
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
			if isExternalUser {
				if atomic.LoadInt64(test.validatorsCanBeStopped) == int64(len(peers)) {
					break wgLoop
				}
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

		validatorBlock := validator.lastBlock
		if minHeight > validatorBlock {
			minHeight = validatorBlock
		}

		if validatorBlock < uint64(test.numBlocks) {
			t.Fatalf("a validator is behind the network index %s and block %v - expected %d",
				index, validatorBlock, test.numBlocks)
		}
	}
	return minHeight
}

func checkBlockConsistence(t *testing.T, minHeight uint64, validators map[string]*testNode, test *testCase) {
	// check that all nodes got the same blocks
	for i := uint64(1); i <= minHeight; i++ {
		blockHash := validators["V0"].service.BlockChain().GetBlockByNumber(i).Hash()
		for index, validator := range validators {
			if isExternalUser(index) {
				continue
			}

			hash := validator.service.BlockChain().GetBlockByNumber(i).Hash()
			if hash != blockHash {
				t.Fatalf("validators %d and %s have different blocks %d - %q vs %s at test: %s",
					0, index, i, hash.String(), blockHash.String(), test.name)
			}
		}
	}
}

func isExternalUser(index string) bool {
	return strings.HasPrefix(index, "E")
}
