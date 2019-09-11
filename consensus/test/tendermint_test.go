package test

import (
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/clearmatics/autonity/accounts"
	"github.com/clearmatics/autonity/accounts/keystore"
	"github.com/clearmatics/autonity/common/fdlimit"
	"github.com/clearmatics/autonity/consensus"
	tendermintCore "github.com/clearmatics/autonity/consensus/tendermint/core"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/eth"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/p2p/enode"
	"golang.org/x/sync/errgroup"
)

func TestTendermintSuccess(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cases := []*testCase{
		{
			name:      "no malicious",
			numPeers:  5,
			numBlocks: 5,
			txPerPeer: 1,
		},
		{
			name:      "one node - always accepts blocks",
			numPeers:  5,
			numBlocks: 5,
			txPerPeer: 1,
			maliciousPeers: map[int]func(basic consensus.Engine) consensus.Engine{
				4: func(basic consensus.Engine) consensus.Engine {
					return tendermintCore.NewVerifyHeaderAlwaysTrueEngine(basic)
				},
			},
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
			runTest(t, testCase)
		})
	}
}

func TestTendermintSlowConnections(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cases := []*testCase{
		{
			name:      "no malicious, one slow node",
			numPeers:  5,
			numBlocks: 5,
			txPerPeer: 1,
			networkRates: map[int]networkRate{
				4: {50 * 1024, 50 * 1024},
			},
		},
		{
			name:      "no malicious, all nodes are slow",
			numPeers:  5,
			numBlocks: 5,
			txPerPeer: 1,
			networkRates: map[int]networkRate{
				0: {50 * 1024, 50 * 1024},
				1: {50 * 1024, 50 * 1024},
				2: {50 * 1024, 50 * 1024},
				3: {50 * 1024, 50 * 1024},
				4: {50 * 1024, 50 * 1024},
			},
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
			runTest(t, testCase)
		})
	}
}

func TestTendermintLongRun(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cases := []*testCase{
		{
			name:      "no malicious - 30 tx per second",
			numPeers:  5,
			numBlocks: 10,
			txPerPeer: 30,
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
			runTest(t, testCase)
		})
	}
}

func TestCheckFeeRedirection(t *testing.T) {
	hookGenerator := func() (hook, hook) {
		prevBlockBalance := uint64(0)
		fBefore := func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error {
			addr, err := validator.service.BlockChain().Config().AutonityContractConfig.GetContractAddress()
			if err != nil {
				t.Fatal(err)
			}
			st, _ := validator.service.BlockChain().State()
			if block.NumberU64() == 1 && st.GetBalance(addr).Uint64() != 0 {
				t.Fatal("incorrect balance on the first block")
			}
			return nil
		}
		fAfter := func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error {
			autonityContractAddress, err := validator.service.BlockChain().Config().AutonityContractConfig.GetContractAddress()
			if err != nil {
				t.Fatal(err)
			}
			st, _ := validator.service.BlockChain().State()

			if block.NumberU64() == 1 && prevBlockBalance != 0 {
				t.Fatal("incorrect balance on the first block")
			}
			contractBalance := st.GetBalance(autonityContractAddress)
			if block.NumberU64() > 1 && block.NumberU64() <= uint64(tCase.numBlocks) {
				if contractBalance.Uint64() < prevBlockBalance {
					t.Fatal("Balance must be increased")
				}
			}
			prevBlockBalance = contractBalance.Uint64()
			return nil
		}
		return fBefore, fAfter
	}

	case1Before, case1After := hookGenerator()
	case2Before, case2After := hookGenerator()
	case3Before, case3After := hookGenerator()
	cases := []*testCase{
		{
			name:      "no malicious - 1 tx per second",
			numPeers:  5,
			numBlocks: 5,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				3: case1Before,
			},
			afterHooks: map[int]hook{
				3: case1After,
			},
		},
		{
			name:      "no malicious - 10 tx per second",
			numPeers:  6,
			numBlocks: 10,
			txPerPeer: 10,
			beforeHooks: map[int]hook{
				5: case2Before,
			},
			afterHooks: map[int]hook{
				5: case2After,
			},
		},
		{
			name:      "no malicious - 5 tx per second 4 peers",
			numPeers:  4,
			numBlocks: 5,
			txPerPeer: 5,
			beforeHooks: map[int]hook{
				1: case3Before,
			},
			afterHooks: map[int]hook{
				1: case3After,
			},
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
			runTest(t, testCase)

		})
	}
}
func TestCheckBlockWithSmallFee(t *testing.T) {
	hookGenerator := func() (hook, hook) {
		prevBlockBalance := uint64(0)
		fBefore := func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error {
			addr, err := validator.service.BlockChain().Config().AutonityContractConfig.GetContractAddress()
			if err != nil {
				t.Fatal(err)
			}
			st, _ := validator.service.BlockChain().State()
			if block.NumberU64() == 1 && st.GetBalance(addr).Uint64() != 0 {
				t.Fatal("incorrect balance on the first block")
			}
			return nil
		}
		fAfter := func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error {
			autonityContractAddress, err := validator.service.BlockChain().Config().AutonityContractConfig.GetContractAddress()
			if err != nil {
				t.Fatal(err)
			}
			st, _ := validator.service.BlockChain().State()

			if block.NumberU64() == 1 && prevBlockBalance != 0 {
				t.Fatal("incorrect balance on the first block")
			}
			contractBalance := st.GetBalance(autonityContractAddress)
			if block.NumberU64() > 1 && block.NumberU64() <= uint64(tCase.numBlocks) {
				if contractBalance.Uint64() < prevBlockBalance {
					t.Fatal("Balance must be increased")
				}
			}
			prevBlockBalance = contractBalance.Uint64()
			return nil
		}
		return fBefore, fAfter
	}

	case1Before, case1After := hookGenerator()
	cases := []*testCase{
		{
			name:      "no malicious - 1 tx per second",
			numPeers:  5,
			numBlocks: 5,
			txPerPeer: 3,
			beforeHooks: map[int]hook{
				3: case1Before,
			},
			afterHooks: map[int]hook{
				3: case1After,
			},
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
			runTest(t, testCase)

		})
	}
}

func TestTendermintStartStop(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cases := []*testCase{
		{
			name:      "one node stops for 5 seconds",
			isSkipped: true,
			numPeers:  5,
			numBlocks: 10,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				4: hookStopNode(4, 5),
			},
			afterHooks: map[int]hook{
				4: hookStartNode(4, 5),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "one node stops for 10 seconds",
			isSkipped: true,
			numPeers:  5,
			numBlocks: 10,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				4: hookStopNode(4, 5),
			},
			afterHooks: map[int]hook{
				4: hookStartNode(4, 10),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "one node stops for 20 seconds",
			isSkipped: true,
			numPeers:  5,
			numBlocks: 10,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				4: hookStopNode(4, 5),
			},
			afterHooks: map[int]hook{
				4: hookStartNode(4, 20),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f node stop for 5 seconds at the same block",
			isSkipped: true,
			numPeers:  5,
			numBlocks: 10,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 5),
			},
			afterHooks: map[int]hook{
				3: hookStartNode(3, 5),
				4: hookStartNode(4, 5),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f node stop for 5 seconds at different blocks",
			isSkipped: true,
			numPeers:  5,
			numBlocks: 10,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 6),
			},
			afterHooks: map[int]hook{
				3: hookStartNode(3, 5),
				4: hookStartNode(4, 5),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f+1 node stop for 5 seconds at the same block",
			isSkipped: true,
			numPeers:  5,
			numBlocks: 10,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				2: hookStopNode(2, 5),
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 5),
			},
			afterHooks: map[int]hook{
				2: hookStartNode(2, 5),
				3: hookStartNode(3, 5),
				4: hookStartNode(4, 5),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f+1 node stop for 5 seconds at different blocks",
			isSkipped: true,
			numPeers:  5,
			numBlocks: 10,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				2: hookStopNode(2, 4),
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 7),
			},
			afterHooks: map[int]hook{
				2: hookStartNode(2, 5),
				3: hookStartNode(3, 5),
				4: hookStartNode(4, 5),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f node stop for 10 seconds at the same block",
			isSkipped: true,
			numPeers:  5,
			numBlocks: 10,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 5),
			},
			afterHooks: map[int]hook{
				3: hookStartNode(3, 10),
				4: hookStartNode(4, 10),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f node stop for 10 seconds at different blocks",
			isSkipped: true,
			numPeers:  5,
			numBlocks: 10,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 6),
			},
			afterHooks: map[int]hook{
				3: hookStartNode(3, 10),
				4: hookStartNode(4, 10),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f+1 node stop for 10 seconds at the same block",
			isSkipped: true,
			numPeers:  5,
			numBlocks: 10,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				2: hookStopNode(2, 5),
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 5),
			},
			afterHooks: map[int]hook{
				2: hookStartNode(2, 10),
				3: hookStartNode(3, 10),
				4: hookStartNode(4, 10),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f+1 node stop for 10 seconds at different blocks",
			isSkipped: true,
			numPeers:  5,
			numBlocks: 10,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				2: hookStopNode(2, 4),
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 7),
			},
			afterHooks: map[int]hook{
				2: hookStartNode(2, 10),
				3: hookStartNode(3, 10),
				4: hookStartNode(4, 10),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f node stop for 20 seconds at the same block",
			isSkipped: true,
			numPeers:  5,
			numBlocks: 10,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 5),
			},
			afterHooks: map[int]hook{
				3: hookStartNode(3, 20),
				4: hookStartNode(4, 20),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f node stop for 20 seconds at different blocks",
			isSkipped: true,
			numPeers:  5,
			numBlocks: 10,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 6),
			},
			afterHooks: map[int]hook{
				3: hookStartNode(3, 20),
				4: hookStartNode(4, 20),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f+1 node stop for 20 seconds at the same block",
			isSkipped: true,
			numPeers:  5,
			numBlocks: 10,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				2: hookStopNode(2, 5),
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 5),
			},
			afterHooks: map[int]hook{
				2: hookStartNode(2, 20),
				3: hookStartNode(3, 20),
				4: hookStartNode(4, 20),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f+1 node stop for 20 seconds at different blocks",
			isSkipped: true,
			numPeers:  5,
			numBlocks: 10,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				2: hookStopNode(2, 4),
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 7),
			},
			afterHooks: map[int]hook{
				2: hookStartNode(2, 20),
				3: hookStartNode(3, 20),
				4: hookStartNode(4, 20),
			},
			stopTime: make(map[int]time.Time),
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
			runTest(t, testCase)
		})
	}
}

type testCase struct {
	name           string
	isSkipped      bool
	numPeers       int
	numBlocks      int
	txPerPeer      int
	maliciousPeers map[int]func(basic consensus.Engine) consensus.Engine //map[validatorIndex]consensusConstructor
	networkRates   map[int]networkRate                                   //map[validatorIndex]networkRate
	beforeHooks    map[int]hook                                          //map[validatorIndex]beforeHook
	afterHooks     map[int]hook                                          //map[validatorIndex]afterHook
	stopTime       map[int]time.Time
	mu             sync.RWMutex
}

func (test *testCase) getBeforeHook(index int) hook {
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

func (test *testCase) getAfterHook(index int) hook {
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

func (test *testCase) setStopTime(index int, stopTime time.Time) {
	test.mu.Lock()
	test.stopTime[index] = stopTime
	test.mu.Unlock()
}

func (test *testCase) getStopTime(index int) time.Time {
	test.mu.RLock()
	currentTime := test.stopTime[index]
	test.mu.RUnlock()

	return currentTime
}

type hook func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error

func runTest(t *testing.T, test *testCase) {
	if test.isSkipped {
		t.SkipNow()
	}

	log.Root().SetHandler(log.LvlFilterHandler(log.LvlError, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))
	err := fdlimit.Raise(512 * uint64(test.numPeers))
	if err != nil {
		t.Log("can't rise file description limit. errors are possible")
	}

	// Generate a batch of accounts to seal and fund with
	validators := make([]*testNode, test.numPeers)

	for i := range validators {
		validators[i] = new(testNode)
		validators[i].privateKey, err = crypto.GenerateKey()
		if err != nil {
			t.Fatal("cant make pk", err)
		}
	}

	for i := range validators {
		validators[i].listener, err = net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			panic(err)
		}
	}

	for _, validator := range validators {
		listener := validator.listener
		validator.address = listener.Addr().String()
		port := strings.Split(listener.Addr().String(), ":")[1]
		validator.port, _ = strconv.Atoi(port)

		validator.url = enode.V4URL(
			validator.privateKey.PublicKey,
			net.IPv4(127, 0, 0, 1),
			validator.port,
			validator.port,
		)
	}

	genesis := makeGenesis(validators)
	for i, validator := range validators {
		var engineConstructor func(basic consensus.Engine) consensus.Engine
		if test.maliciousPeers != nil {
			engineConstructor = test.maliciousPeers[i]
		}

		validator.listener.Close()

		rates := test.networkRates[i]

		validator.node, err = makeValidator(genesis, validator.privateKey, validator.address, rates.in, rates.out, engineConstructor)
		if err != nil {
			t.Fatal("cant make a validator", i, err)
		}
	}

	wg := &errgroup.Group{}
	for _, validator := range validators {
		validator := validator

		wg.Go(func() error {
			return validator.startNode()
		})
	}
	err = wg.Wait()
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		for _, validator := range validators {
			if validator.isRunning {
				err = validator.node.Stop()
				if err != nil {
					panic(err)
				}
			}
		}
	}()

	wg = &errgroup.Group{}
	for _, validator := range validators {
		validator := validator

		wg.Go(func() error {
			return validator.startService()
		})
	}
	err = wg.Wait()
	if err != nil {
		t.Fatal(err)
	}

	wg = &errgroup.Group{}
	for i, validator := range validators {
		validator := validator
		i := i

		wg.Go(func() error {
			log.Debug("peers", "i", i,
				"peers", len(validator.node.Server().Peers()),
				"staticPeers", len(validator.node.Server().StaticNodes),
				"trustedPeers", len(validator.node.Server().TrustedNodes),
				"nodes", len(validators))
			return nil
		})
	}
	err = wg.Wait()
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		for _, validator := range validators {
			validator.subscription.Unsubscribe()
		}
	}()

	// each peer sends one tx per block
	sendTransactions(t, test, validators, test.txPerPeer, true)
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
	if err := validator.node.Stop(); err != nil {
		return fmt.Errorf("cannot stop a node %s", err)
	}

	validator.node.Wait()

	validator.isRunning = false
	return nil
}

func (validator *testNode) startService() error {
	var ethereum *eth.Ethereum
	if err := validator.node.Service(&ethereum); err != nil {
		return fmt.Errorf("cant start a node %s", err)
	}

	time.Sleep(100 * time.Millisecond)

	if err := ethereum.StartMining(1); err != nil {
		return fmt.Errorf("cant start mining %s", err)
	}

	for !ethereum.IsMining() {
		time.Sleep(50 * time.Millisecond)
	}

	validator.service = ethereum

	if validator.eventChan == nil {
		validator.eventChan = make(chan core.ChainEvent, 1024)
	}

	validator.subscription = validator.service.BlockChain().SubscribeChainEvent(validator.eventChan)

	validator.isRunning = true

	return nil
}

func sendTransactions(t *testing.T, test *testCase, validators []*testNode, txPerPeer int, errorOnTx bool) {
	const blocksToWait = 10

	txs := make(map[uint64]int) // blockNumber to count
	txsMu := sync.Mutex{}

	wg := &errgroup.Group{}
	for index, validator := range validators {
		index := index
		validator := validator

		// skip malicious nodes
		if test.maliciousPeers != nil {
			if _, ok := test.maliciousPeers[index]; ok {
				continue
			}
		}

		wg.Go(func() error {
			var (
				blocksPassed int
				lastBlock    uint64
				err          error
			)

			fromAddr := crypto.PubkeyToAddress(validator.privateKey.PublicKey)

		wgLoop:
			for {
				select {
				case ev := <-validator.eventChan:
					// before hook
					err = runHook(test.getBeforeHook(index), test, ev.Block, validator, index)
					if err != nil {
						return err
					}

					// actual forming and sending transaction
					log.Debug("peer", "address", crypto.PubkeyToAddress(validator.privateKey.PublicKey).String(), "block", ev.Block.Number().Uint64(), "isRunning", validator.isRunning)

					if validator.isRunning {
						currentBlock := ev.Block.Number().Uint64()
						if currentBlock <= lastBlock {
							return fmt.Errorf("expected next block %d got %d. Block %v", lastBlock+1, currentBlock, ev.Block)
						}
						lastBlock = currentBlock

						txsMu.Lock()
						if _, ok := txs[currentBlock]; !ok {
							txs[currentBlock] = ev.Block.Transactions().Len()
						}
						txsMu.Unlock()

						if blocksPassed <= test.numBlocks {
							for i := 0; i < txPerPeer; i++ {
								nextValidatorIndex := (index + i + 1) % len(validators)
								toAddr := crypto.PubkeyToAddress(validators[nextValidatorIndex].privateKey.PublicKey)

								if innerErr := sendTx(validator.service, validator.privateKey, fromAddr, toAddr); innerErr != nil {
									return innerErr
								}
							}
						}
					}

					// after hook
					err = runHook(test.getAfterHook(index), test, ev.Block, validator, index)
					if err != nil {
						return err
					}

					// check transactions status if all blocks are passed
					blocksPassed++
					if validator.isRunning && blocksPassed >= test.numBlocks+blocksToWait {
						pending, queued := validator.service.TxPool().Stats()
						if errorOnTx {
							if pending != 0 {
								return fmt.Errorf("after a new block it should be 0 pending transactions got %d. block %d", pending, ev.Block.Number().Uint64())
							}
							if queued != 0 {
								return fmt.Errorf("after a new block it should be 0 queued transactions got %d. block %d", queued, ev.Block.Number().Uint64())
							}
						}

						break wgLoop
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
				}
			}

			return nil
		})
	}
	if err := wg.Wait(); err != nil {
		t.Fatal(err)
	}

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

	//check that all nodes reached minimum height
	lastBlock := uint64(keys[len(keys)-1])
	for index, validator := range validators {
		validatorBlock := validator.service.BlockChain().CurrentBlock().Number().Uint64()

		if validatorBlock < lastBlock-blocksToWait/2 {
			t.Fatalf("a validator is behind the network index %d(%v) and block %v - expected %d",
				index, validator, validatorBlock, lastBlock)
		}
	}
}

func runHook(validatorHook hook, test *testCase, block *types.Block, validator *testNode, index int) error {
	if validatorHook == nil {
		return nil
	}

	err := validatorHook(block, validator, test, time.Now())
	if err != nil {
		return fmt.Errorf("error while executing before hook for validator index %d(%v) and block %v",
			index, validator, block)
	}

	return nil
}

func hookStopNode(nodeIndex int, blockNum uint64) hook {
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

func hookStartNode(nodeIndex int, durationAfterStop float64) hook {
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
