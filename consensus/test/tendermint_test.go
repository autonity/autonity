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
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/eth"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/p2p/enode"
	"golang.org/x/sync/errgroup"
)

func TestTendermint(t *testing.T) {
	cases := []testCase{
		{
			"no malicious",
			5,
			5,
			1,
			nil,
			nil,
		},
		{
			"no malicious - 100 tx per second",
			5,
			10,
			100,
			nil,
			nil,
		},
		{
			"no malicious, one slow node",
			5,
			5,
			1,
			nil,
			map[int]networkRate{
				4: {50 * 1024, 50 * 1024},
			},
		},
		{
			"no malicious, all nodes are slow",
			5,
			5,
			1,
			nil,
			map[int]networkRate{
				0: {50 * 1024, 50 * 1024},
				1: {50 * 1024, 50 * 1024},
				2: {50 * 1024, 50 * 1024},
				3: {50 * 1024, 50 * 1024},
				4: {50 * 1024, 50 * 1024},
			},
		},
		{
			"10 nodes, 20 blocks",
			10,
			20,
			10,
			nil,
			nil,
		},
		{
			"one node - always accepts blocks",
			5,
			5,
			1,
			map[int]func(basic consensus.Engine) consensus.Engine{
				4: func(basic consensus.Engine) consensus.Engine {
					return tendermintCore.NewVerifyHeaderAlwaysTrueEngine(basic)
				},
			},
			nil,
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
			tunTest(t, testCase)
		})
	}
}

type testCase struct {
	name           string
	numPeers       int
	numBlocks      int
	txPerPeer      int
	maliciousPeers map[int]func(basic consensus.Engine) consensus.Engine
	networkRates   map[int]networkRate
}

func tunTest(t *testing.T, test testCase) {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlError, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))
	_ = fdlimit.Raise(2048)

	var err error

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

	for i, validator := range validators {
		// Inject the signer key and start sealing with it
		store := validator.node.AccountManager().Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)

		var signer accounts.Account
		signer, err = store.ImportECDSA(validator.privateKey, "")
		if err != nil {
			t.Fatal("import pk", i, err)
		}

		if err = store.Unlock(signer, ""); err != nil {
			t.Fatal("cant unlock", i, err)
		}
	}

	wg := &errgroup.Group{}
	for i, validator := range validators {
		validator := validator
		i := i

		wg.Go(func() error {
			err = validator.node.Start()
			if err != nil {
				return fmt.Errorf("cannot start a node %d %s", i, err)
			}

			// Start tracking the node and it's enode
			validator.enode = validator.node.Server().Self()
			return nil
		})
	}
	err = wg.Wait()
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		for _, validator := range validators {
			err = validator.node.Stop()
			if err != nil {
				panic(err)
			}
		}
	}()

	wg = &errgroup.Group{}
	for i, validator := range validators {
		validator := validator
		i := i

		wg.Go(func() error {
			var ethereum *eth.Ethereum
			if err = validator.node.Service(&ethereum); err != nil {
				return fmt.Errorf("cant start a node %d %s", i, err)
			}

			for !ethereum.IsListening() {
				time.Sleep(50 * time.Millisecond)
			}

			if err = ethereum.StartMining(1); err != nil {
				return fmt.Errorf("cant start mining %d %s", i, err)
			}

			for !ethereum.IsMining() {
				time.Sleep(50 * time.Millisecond)
			}

			validator.service = ethereum

			return nil
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

	for _, validator := range validators {
		validator.eventChan = make(chan core.ChainEvent, 1024)
		validator.subscription = validator.service.BlockChain().SubscribeChainEvent(validator.eventChan)
	}

	defer func() {
		for _, validator := range validators {
			validator.subscription.Unsubscribe()
		}
	}()

	// each peer sends one tx per block
	sendTransactions(t, test, validators, test.txPerPeer, true)
}

func sendTransactions(t *testing.T, test testCase, validators []*testNode, txPerPeer int, errorOnTx bool) {
	var err error

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
			var blocksPassed int
			var lastBlock uint64

			ethereum := validator.service
			fromAddr := crypto.PubkeyToAddress(validator.privateKey.PublicKey)

		wgLoop:
			for {
				select {
				case ev := <-validator.eventChan:
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

					if blocksPassed < test.numBlocks {
						for i := 0; i < txPerPeer; i++ {
							nextValidatorIndex := (index + i + 1) % len(validators)
							toAddr := crypto.PubkeyToAddress(validators[nextValidatorIndex].privateKey.PublicKey)

							if innerErr := sendTx(ethereum, validator.privateKey, fromAddr, toAddr); innerErr != nil {
								return innerErr
							}
						}
					}

					blocksPassed++
					if blocksPassed >= test.numBlocks+5 {
						pending, queued := ethereum.TxPool().Stats()
						if errorOnTx {
							if pending != 0 {
								t.Fatal("after a new block it should be 0 pending transactions got", pending)
							}
							if queued != 0 {
								t.Fatal("after a new block it should be 0 queued transactions got", pending)
							}
						}

						break wgLoop
					}
				case innerErr := <-validator.subscription.Err():
					return fmt.Errorf("error in blockchain %q", innerErr)
				}
			}

			return nil
		})
	}
	err = wg.Wait()
	if err != nil {
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
}
