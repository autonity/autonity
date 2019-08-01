package tendermint

import (
	"crypto/ecdsa"
	"io/ioutil"
	"math/big"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
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
	//t.SkipNow()

	log.Root().SetHandler(log.LvlFilterHandler(log.LvlInfo, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))
	_ = fdlimit.Raise(2048)

	var err error
	// Generate a batch of accounts to seal and fund with
	sealers := make([]*ecdsa.PrivateKey, 4)
	for i := 0; i < len(sealers); i++ {
		sealers[i], err = crypto.GenerateKey()
		if err != nil {
			t.Fatal("cant make pk", err)
		}
	}
	// Create a Clique network based off of the Rinkeby config

	genesis := makeGenesis(sealers)

	// get enode addresses before node.Start()
	addresses := make([]string, len(sealers))
	ports := make([]int, len(sealers))
	urls := make([]string, len(sealers))

	listeners := make([]net.Listener, len(sealers))
	for i := range listeners {
		var err error
		listeners[i], err = net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			panic(err)
		}
	}

	for i, sealer := range sealers {
		listener := listeners[i]
		addresses[i] = listener.Addr().String()
		port := strings.Split(listener.Addr().String(), ":")[1]
		ports[i], _ = strconv.Atoi(port)
		urls[i] = enode.V4URL(sealer.PublicKey, net.IPv4(127, 0, 0, 1), ports[i], ports[i])
	}

	genesis.Config.EnodeWhitelist = urls

	nodes := make([]*node.Node, len(sealers))
	enodes := make([]*enode.Node, len(sealers))
	for i := range sealers {
		// Start the node and wait until it's up
		// We want only one mock node
		var engineConstructor func(basic consensus.Engine) consensus.Engine
		if i == 0 {
			engineConstructor = func(basic consensus.Engine) consensus.Engine {
				return tendermintCore.NewCoreQuorumAlwaysFalse(basic)
			}
		}

		listeners[i].Close()
		nodes[i], err = makeSealer(genesis, sealers[i], addresses[i], engineConstructor)
		if err != nil {
			t.Fatal("cant make a validator", i, err)
		}
	}

	for _, sealer := range sealers {
		//todo: me
		crypto.PubkeyToAddress(sealer.PublicKey)
	}

	for i, node := range nodes {
		// Inject the signer key and start sealing with it
		store := node.AccountManager().Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)
		signer, err := store.ImportECDSA(sealers[i], "")
		if err != nil {
			t.Fatal("import pk", i, err)
		}

		if err := store.Unlock(signer, ""); err != nil {
			t.Fatal("cant unlock", i, err)
		}
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(nodes))

	for i := range nodes {
		go func(node *node.Node, i int) {
			err = node.Start()
			if err != nil {
				t.Fatal("cannot start a node", i, err)
			}

			// Start tracking the node and it's enode
			enodes[i] = node.Server().Self()

			for node.Server().NodeInfo().Ports.Listener == 0 {
				time.Sleep(250 * time.Millisecond)
			}
			wg.Done()
		}(nodes[i], i)
	}

	wg.Wait()

	defer func() {
		for _, node := range nodes {
			_ = node.Stop()
		}
	}()

	wg = &sync.WaitGroup{}
	wg.Add(len(nodes))
	for i := range nodes {
		go func(node *node.Node) {
			for _, n := range enodes {
				node.Server().AddPeer(n)
			}
			wg.Done()
		}(nodes[i])
	}
	wg.Wait()

	wg = &sync.WaitGroup{}
	wg.Add(len(nodes))
	for i := range nodes {
		go func(node *node.Node) {
			var ethereum *eth.Ethereum
			if err = node.Service(&ethereum); err != nil {
				t.Fatal("cant start a node", i, err)
			}

			// todo: check if we need it
			if err = ethereum.StartMining(1); err != nil {
				t.Fatal("cant start mining", i, err)
			}
			wg.Done()
		}(nodes[i])
	}

	time.Sleep(20 * time.Second)

	// Start injecting transactions from the faucet like crazy
	nonces := make([]uint64, len(sealers))
	for {
		index := rand.Intn(len(sealers))

		// Fetch the accessor for the relevant signer
		var ethereum *eth.Ethereum
		if err := nodes[index%len(nodes)].Service(&ethereum); err != nil {
			panic(err)
		}
		// Create a self transaction and inject into the pool
		tx, err := types.SignTx(types.NewTransaction(nonces[index], crypto.PubkeyToAddress(sealers[index].PublicKey), new(big.Int), 21000, big.NewInt(100000000000), nil), types.HomesteadSigner{}, sealers[index])
		if err != nil {
			panic(err)
		}
		if err := ethereum.TxPool().AddLocal(tx); err != nil {
			panic(err)
		}
		nonces[index]++

		// Wait if we're too saturated
		if pend, _ := ethereum.TxPool().Stats(); pend > 2048 {
			time.Sleep(100 * time.Millisecond)
		}
	}
}

var (
	defaultDifficulty = big.NewInt(1)
	nilUncleHash      = types.CalcUncleHash(nil) // Always Keccak256(RLP([])) as uncles are meaningless outside of PoW.
	emptyNonce        = types.BlockNonce{}
	now               = time.Now
)

func makeGenesis(sealers []*ecdsa.PrivateKey) *core.Genesis {
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
	for _, sealer := range sealers {
		genesis.Alloc[crypto.PubkeyToAddress(sealer.PublicKey)] = core.GenesisAccount{
		Balance: new(big.Int).Exp(big.NewInt(2), big.NewInt(128), nil),
	}
	}

	validatorsAddresses := make([]string, len(sealers))
	for i, validator := range sealers {
		validatorsAddresses[i] = crypto.PubkeyToAddress(validator.PublicKey).String()
		genesis.Config.EnodeWhitelist = append(genesis.Config.EnodeWhitelist, enode.PubkeyToIDV4(&validator.PublicKey).String())
	}

	genesis.Validators = validatorsAddresses
	genesis.SetBFT()

	return genesis
}

func makeSealer(genesis *core.Genesis, nodekey *ecdsa.PrivateKey, listenAddr string, cons func(basic consensus.Engine) consensus.Engine) (*node.Node, error) {
	// Define the basic configurations for the Ethereum node
	datadir, _ := ioutil.TempDir("", "")

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
			PrivateKey: nodekey,
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
