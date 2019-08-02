package tendermint

import (
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/clearmatics/autonity/accounts"
	"github.com/clearmatics/autonity/accounts/keystore"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/common/fdlimit"
	"github.com/clearmatics/autonity/common/math"
	"github.com/clearmatics/autonity/consensus"
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
	"golang.org/x/sync/errgroup"
)

func TestTendermint(t *testing.T) {
	cases := []testCase{
		{
			"no malicious",
			5,
			5,
			map[int]struct{}{},
			nil,
		},
		{
			"one node - always accepts blocks",
			5,
			5,
			map[int]struct{}{
				4: {},
			},
			func(index int) func(basic consensus.Engine) consensus.Engine {
				if index == 4 {
					return func(basic consensus.Engine) consensus.Engine {
						return tendermintCore.NewVerifyHeaderAlwaysTrueEngine(basic)
					}
				}
				return nil
			},
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
	maliciousPeers map[int]struct{}
	newConsensus   consensusConstructor
}

type consensusConstructor func(index int) func(basic consensus.Engine) consensus.Engine

type testNode struct {
	validator *ecdsa.PrivateKey
	address string
	ports int
	url string
	listener net.Listener
	node *node.Node
	enode *enode.Node
	service *eth.Ethereum
	eventChan chan core.ChainEvent
	subscription event.Subscription
}

func tunTest(t *testing.T, test testCase) {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlError, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))
	_ = fdlimit.Raise(2048)

	var err error

	// Generate a batch of accounts to seal and fund with
	validators := make([]*ecdsa.PrivateKey, test.numPeers)

	for i := 0; i < len(validators); i++ {
		validators[i], err = crypto.GenerateKey()
		if err != nil {
			t.Fatal("cant make pk", err)
		}
	}

	addresses := make([]string, len(validators))
	ports := make([]int, len(validators))
	urls := make([]string, len(validators))
	listeners := make([]net.Listener, len(validators))

	for i := range listeners {
		listeners[i], err = net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			panic(err)
		}
	}

	for i, sealer := range validators {
		listener := listeners[i]
		addresses[i] = listener.Addr().String()
		port := strings.Split(listener.Addr().String(), ":")[1]
		ports[i], _ = strconv.Atoi(port)
		urls[i] = enode.V4URL(sealer.PublicKey, net.IPv4(127, 0, 0, 1), ports[i], ports[i])
	}

	genesis := makeGenesis(validators, urls)
	nodes := make([]*node.Node, len(validators))
	enodes := make([]*enode.Node, len(validators))
	for i := range validators {
		var engineConstructor func(basic consensus.Engine) consensus.Engine
		if test.newConsensus != nil {
			engineConstructor = test.newConsensus(i)
		}

		listeners[i].Close()
		nodes[i], err = makeValidator(genesis, validators[i], addresses[i], engineConstructor)
		if err != nil {
			t.Fatal("cant make a validator", i, err)
		}
	}

	for i, node := range nodes {
		// Inject the signer key and start sealing with it
		store := node.AccountManager().Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)

		var signer accounts.Account
		signer, err = store.ImportECDSA(validators[i], "")
		if err != nil {
			t.Fatal("import pk", i, err)
		}

		if err = store.Unlock(signer, ""); err != nil {
			t.Fatal("cant unlock", i, err)
		}
	}

	wg := &errgroup.Group{}
	for i := range nodes {
		node := nodes[i]
		i := i

		wg.Go(func() error {
			err = node.Start()
			if err != nil {
				return fmt.Errorf("cannot start a node %d %s", i, err)
			}

			// Start tracking the node and it's enode
			enodes[i] = node.Server().Self()
			return nil
		})
	}
	err = wg.Wait()
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		for _, node := range nodes {
			err = node.Stop()
			if err != nil {
				panic(err)
			}
		}
	}()

	//time.Sleep(2 * time.Second)

	wg = &errgroup.Group{}
	services := make([]*eth.Ethereum, len(nodes))
	for i := range nodes {
		node := nodes[i]
		i := i

		wg.Go(func() error {
			var ethereum *eth.Ethereum
			if err = node.Service(&ethereum); err != nil {
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

			services[i] = ethereum

			return nil
		})
	}
	err = wg.Wait()
	if err != nil {
		t.Fatal(err)
	}

	wg = &errgroup.Group{}
	for i := range nodes {
		node := nodes[i]
		i := i

		wg.Go(func() error {
			log.Debug("peers", "i", i,
				"peers", len(node.Server().Peers()),
				"staticPeers", len(node.Server().StaticNodes),
				"trustedPeers", len(node.Server().TrustedNodes),
				"nodes", len(nodes))
			return nil
		})
	}
	err = wg.Wait()
	if err != nil {
		t.Fatal(err)
	}

	chs := make([]chan core.ChainEvent, len(services))
	subs := make([]event.Subscription, len(services))
	for i, ethereum := range services {
		chs[i] = make(chan core.ChainEvent, 1024)
		subs[i] = ethereum.BlockChain().SubscribeChainEvent(chs[i])
	}

	defer func() {
		for _, sub := range subs {
			sub.Unsubscribe()
		}
	}()


	// each peer sends one tx per block
	sendTransactions(t, chs, test, services, validators, subs)
}

func sendTransactions(t *testing.T, chs []chan core.ChainEvent, test testCase, services []*eth.Ethereum, validators []*ecdsa.PrivateKey, subs []event.Subscription) {
	var err error

	wg := &errgroup.Group{}
	for index := range chs {
		index := index

		// skip malicious nodes
		if _, ok := test.maliciousPeers[index]; ok {
			continue
		}

		wg.Go(func() error {
			var blocksPassed int
			var lastBlock uint64

			chainEvents := chs[index]
			ethereum := services[index]
			fromAddr := crypto.PubkeyToAddress(validators[index].PublicKey)

			nextValidatorIndex := (index + 1) % len(validators)
			toAddr := crypto.PubkeyToAddress(validators[nextValidatorIndex].PublicKey)
			from := validators[index]

		wgLoop:
			for {
				select {
				case ev := <-chainEvents:
					currentBlock := ev.Block.Number().Uint64()
					if currentBlock <= lastBlock {
						return fmt.Errorf("expected next block %d got %d. Block %v", lastBlock+1, currentBlock, ev.Block)
					}
					lastBlock = currentBlock

					if blocksPassed < test.numBlocks {
						if innerErr := sendTx(ethereum, from, fromAddr, toAddr); innerErr != nil {
							return innerErr
						}
					}

					blocksPassed++
					if blocksPassed >= test.numBlocks+3 {
						pending, queued := ethereum.TxPool().Stats()
						if pending != 0 {
							t.Fatal("after a new block it should be 0 pending transactions got", pending)
						}
						if queued != 0 {
							t.Fatal("after a new block it should be 0 queued transactions got", pending)
						}

						break wgLoop
					}
				case innerErr := <-subs[index].Err():
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
}

func sendTx(service *eth.Ethereum, fromValidator *ecdsa.PrivateKey, fromAddr common.Address, toAddr common.Address) error {
	nonce := service.TxPool().State().GetNonce(fromAddr)

	tx, err := types.SignTx(
		types.NewTransaction(
			nonce,
			toAddr,
			big.NewInt(1),
			210000000,
			big.NewInt(100000000000),
			nil,
		),
		types.HomesteadSigner{}, fromValidator)
	if err != nil {
		return err
	}

	return service.TxPool().AddLocal(tx)
}

func makeGenesis(validators []*ecdsa.PrivateKey, enodes []string) *core.Genesis {
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
	for _, sealer := range validators {
		genesis.Alloc[crypto.PubkeyToAddress(sealer.PublicKey)] = core.GenesisAccount{
			Balance: new(big.Int).Exp(big.NewInt(2), big.NewInt(128), nil),
		}
	}

	validatorsAddresses := make([]string, len(validators))
	for i, validator := range validators {
		validatorsAddresses[i] = crypto.PubkeyToAddress(validator.PublicKey).String()
	}
	genesis.Config.EnodeWhitelist = enodes

	genesis.Validators = validatorsAddresses
	err := genesis.SetBFT()
	if err != nil {
		panic(err)
	}

	return genesis
}

func makeValidator(genesis *core.Genesis, nodekey *ecdsa.PrivateKey, listenAddr string, cons func(basic consensus.Engine) consensus.Engine) (*node.Node, error) {
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
