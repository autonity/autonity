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
			nil,
			nil,
		},
		{
			"no malicious, one slow node",
			5,
			5,
			nil,
			map[int]networkRate{
				4: {50 * 1024, 50 * 1024},
			},
		},
		{
			"one node - always accepts blocks",
			5,
			5,
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
	maliciousPeers map[int]func(basic consensus.Engine) consensus.Engine
	networkRates   map[int]networkRate
}

type networkRate struct {
	in  int64
	out int64
}

type testNode struct {
	privateKey   *ecdsa.PrivateKey
	address      string
	port         int
	url          string
	listener     net.Listener
	node         *node.Node
	enode        *enode.Node
	service      *eth.Ethereum
	eventChan    chan core.ChainEvent
	subscription event.Subscription
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

	//time.Sleep(2 * time.Second)

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
	sendTransactions(t, test, validators)
}

func sendTransactions(t *testing.T, test testCase, validators []*testNode) {
	var err error

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

			nextValidatorIndex := (index + 1) % len(validators)
			toAddr := crypto.PubkeyToAddress(validators[nextValidatorIndex].privateKey.PublicKey)

		wgLoop:
			for {
				select {
				case ev := <-validator.eventChan:
					currentBlock := ev.Block.Number().Uint64()
					if currentBlock <= lastBlock {
						return fmt.Errorf("expected next block %d got %d. Block %v", lastBlock+1, currentBlock, ev.Block)
					}
					lastBlock = currentBlock

					if blocksPassed < test.numBlocks {
						if innerErr := sendTx(ethereum, validator.privateKey, fromAddr, toAddr); innerErr != nil {
							return innerErr
						}
					}

					blocksPassed++
					if blocksPassed >= test.numBlocks+5 {
						pending, queued := ethereum.TxPool().Stats()
						if pending != 0 {
							t.Fatal("after a new block it should be 0 pending transactions got", pending)
						}
						if queued != 0 {
							t.Fatal("after a new block it should be 0 queued transactions got", pending)
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
