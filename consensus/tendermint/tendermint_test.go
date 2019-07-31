package tendermint

import (
	"crypto/ecdsa"
	"github.com/clearmatics/autonity/consensus/tendermint/config"
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

	// Generate a batch of accounts to seal and fund with
	faucets := make([]*ecdsa.PrivateKey, 128)
	for i := 0; i < len(faucets); i++ {
		faucets[i], _ = crypto.GenerateKey()
	}
	sealers := make([]*ecdsa.PrivateKey, 4)
	for i := 0; i < len(sealers); i++ {
		sealers[i], _ = crypto.GenerateKey()
	}
	// Create a Clique network based off of the Rinkeby config

	genesis := makeGenesis(faucets, sealers)

	var (
		nodes  []*node.Node
		enodes []*enode.Node
	)

	// get enode addresses before node.Start()
	addresses := make([]string, len(sealers))
	ports := make([]int, len(sealers))
	urls := make([]string, len(sealers))
	for i, sealer := range sealers {
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			panic(err)
		}

		addresses[i] = strings.Split(listener.Addr().String(), ":")[1]
		ports[i], _ = strconv.Atoi(addresses[i])

		listener.Close()

		urls[i] = enode.V4URL(sealer.PublicKey, net.IPv4(127,0,0,1), ports[i], ports[i])
	}

	genesis.Config.EnodeWhitelist = urls

	for i, sealer := range sealers {
		// Start the node and wait until it's up

		// We want only one mock node
		var engineConstructor func(basic consensus.Engine) consensus.Engine
		if i == 0 {
			engineConstructor = func(basic consensus.Engine) consensus.Engine {
				return tendermintCore.NewCoreQuorumAlwaysFalse(basic)
			}
		}

		node, err := makeSealer(genesis, addresses[i],  engineConstructor)
		if err != nil {
			panic(err)
		}

		//todo: aggregate into the single Stop function
		defer func() {
			_ = node.Stop()
		}()

		for node.Server().NodeInfo().Ports.Listener == 0 {
			time.Sleep(250 * time.Millisecond)
		}
		// Connect the node to al the previous ones
		for _, n := range enodes {
			node.Server().AddPeer(n)
		}
		// Start tracking the node and it's enode
		nodes = append(nodes, node)
		enodes = append(enodes, node.Server().Self())

		// Inject the signer key and start sealing with it
		store := node.AccountManager().Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)
		signer, err := store.ImportECDSA(sealer, "")
		if err != nil {
			panic(err)
		}
		if err := store.Unlock(signer, ""); err != nil {
			panic(err)
		}
	}
	// Iterate over all the nodes and start signing with them
	time.Sleep(10 * time.Second)

	for _, node := range nodes {
		var ethereum *eth.Ethereum
		if err := node.Service(&ethereum); err != nil {
			panic(err)
		}

		// todo: check if we need it
		if err := ethereum.StartMining(1); err != nil {
			panic(err)
		}
	}

	time.Sleep(10 * time.Second)

	// Start injecting transactions from the faucet like crazy
	nonces := make([]uint64, len(faucets))
	for {
		index := rand.Intn(len(faucets))

		// Fetch the accessor for the relevant signer
		var ethereum *eth.Ethereum
		if err := nodes[index%len(nodes)].Service(&ethereum); err != nil {
			panic(err)
		}
		// Create a self transaction and inject into the pool
		tx, err := types.SignTx(types.NewTransaction(nonces[index], crypto.PubkeyToAddress(faucets[index].PublicKey), new(big.Int), 21000, big.NewInt(100000000000), nil), types.HomesteadSigner{}, faucets[index])
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

func makeGenesis(faucets []*ecdsa.PrivateKey, sealers []*ecdsa.PrivateKey) *core.Genesis {
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
	for _, faucet := range faucets {
		genesis.Alloc[crypto.PubkeyToAddress(faucet.PublicKey)] = core.GenesisAccount{
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


func makeSealer(genesis *core.Genesis, listenAddr string, cons func(basic consensus.Engine) consensus.Engine) (*node.Node, error) {
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
			ListenAddr:  "0.0.0.0:0",
			NoDiscovery: true,
			MaxPeers:    25,
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
			Tendermint: 	 *config.DefaultConfig,
		}, cons)
	}); err != nil {
		return nil, err
	}
	// Start the node and return if successful
	return stack, stack.Start()
}
