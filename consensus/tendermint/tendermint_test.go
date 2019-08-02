package tendermint

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/clearmatics/autonity/accounts"
	"golang.org/x/sync/errgroup"
	"io/ioutil"
	"math/big"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/clearmatics/autonity/accounts/keystore"
	"github.com/clearmatics/autonity/common/fdlimit"
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
	"github.com/clearmatics/autonity/p2p/enode"
	"github.com/clearmatics/autonity/params"
)

func TestTendermint(t *testing.T) {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlDebug, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))
	_ = fdlimit.Raise(2048)

	var err error
	// Generate a batch of accounts to seal and fund with
	validators := make([]*ecdsa.PrivateKey, 5)
	for i := 0; i < len(validators); i++ {
		validators[i], err = crypto.GenerateKey()
		if err != nil {
			t.Fatal("cant make pk", err)
		}
	}

	// get enode addresses before node.Start()
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
		// Start the node and wait until it's up
		// We want only one mock node
		var engineConstructor func(basic consensus.Engine) consensus.Engine
		if i == 10000 {
			engineConstructor = func(basic consensus.Engine) consensus.Engine {
				return tendermintCore.NewCoreQuorumAlwaysFalse(basic)
			}
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
		log.Error("STOPPING ALL NODES")
		for _, node := range nodes {
			_ = node.Stop()
		}
	}()

	time.Sleep(5*time.Second)

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
				log.Error("INIT: NODE IS WAITING - LISTENING!!!", "i", i)
				time.Sleep(250 * time.Millisecond)
			}

			if err = ethereum.StartMining(1); err != nil {
				return fmt.Errorf("cant start mining %d %s", i, err)
			}

			for !ethereum.IsMining() {
				log.Error("INIT: NODE IS WAITING - MINING!!!", "i", i)
				time.Sleep(250 * time.Millisecond)
			}
			log.Error("INIT: NODE DONE !!!", "i", i)
			services[i] =ethereum

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
			log.Error("******* peers", "i", i,
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

	log.Error("====================================================================")
	time.Sleep(30 * time.Second)
	log.Error("====================================================================")

	// Start injecting transactions from the faucet like crazy
	nonces := make([]uint64, len(validators))
	i := 0
	for {
		i++

		index := rand.Intn(len(validators))

		log.Error("++++++ 1", "i", i)

		// Fetch the accessor for the relevant signer
		var ethereum *eth.Ethereum
		if err := nodes[index%len(nodes)].Service(&ethereum); err != nil {
			panic(err)
		}

		log.Error("++++++ 2", "i", i)

		// Create a self transaction and inject into the pool
		tx, err := types.SignTx(types.NewTransaction(nonces[index], crypto.PubkeyToAddress(validators[index].PublicKey), new(big.Int), 21000, big.NewInt(100000000000), nil), types.HomesteadSigner{}, validators[index])
		if err != nil {
			panic(err)
		}
		log.Error("++++++ 3", "i", i)
		if err := ethereum.TxPool().AddLocal(tx); err != nil {
			panic(err)
		}
		log.Error("++++++ 4", "i", i)
		nonces[index]++

		// Wait if we're too saturated
		if pend, _ := ethereum.TxPool().Stats(); pend > 2048 {
			time.Sleep(100 * time.Millisecond)
		}

		log.Error("++++++ 5", "i", i)
	}
}

var (
	defaultDifficulty = big.NewInt(1)
	emptyNonce        = types.BlockNonce{}
)

func makeGenesis(validators []*ecdsa.PrivateKey, enodes []string) *core.Genesis {
	// generate genesis block
	genesis := core.DefaultGenesisBlock()
	genesis.Config = params.TestChainConfig

	// force enable Istanbul engine
	genesis.Config.Tendermint = &params.TendermintConfig{}
	genesis.Config.Ethash = nil
	genesis.Difficulty = defaultDifficulty
	genesis.Nonce = emptyNonce.Uint64()
	genesis.Mixhash = types.BFTDigest

	genesis.Alloc = core.GenesisAlloc{}
	for _, sealer := range validators {
		genesis.Alloc[crypto.PubkeyToAddress(sealer.PublicKey)] = core.GenesisAccount{
			Balance: new(big.Int).Exp(big.NewInt(2), big.NewInt(12), nil),
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
		ctx.NodeKey()
		return eth.New(ctx, &eth.Config{
			Genesis:         genesis,
			NetworkId:       genesis.Config.ChainID.Uint64(),
			SyncMode:        downloader.FullSync,
			DatabaseCache:   256,
			DatabaseHandles: 256,
			TxPool:          core.DefaultTxPoolConfig,
			Tendermint:      *config.DefaultConfig,
		}, cons)
	}); err != nil {
		return nil, err
	}

	// Start the node and return if successful
	return stack, nil
}
