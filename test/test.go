package test

import (
	"crypto/ecdsa"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/clearmatics/autonity/cmd/gengen/gengen"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/eth"
	"github.com/clearmatics/autonity/eth/downloader"
	"github.com/clearmatics/autonity/node"
	"github.com/clearmatics/autonity/p2p"
	"github.com/clearmatics/autonity/params"
)

var (
	baseNodeConfig *node.Config = &node.Config{
		Name:    "autonity",
		Version: params.Version,
		P2P: p2p.Config{
			MaxPeers:              100,
			DialHistoryExpiration: time.Millisecond,
		},
		NoUSB:    true,
		HTTPHost: "0.0.0.0",
		WSHost:   "0.0.0.0",
	}
	genesis *core.Genesis = &core.Genesis{}

	baseEthConfig = &eth.Config{
		SyncMode:        downloader.FullSync,
		DatabaseCache:   256,
		DatabaseHandles: 256,
		TxPool:          core.DefaultTxPoolConfig,
	}

	// Set up 4 validators with 10 E each and 1 stake.
	// Users = []string{
	// 	"10e18,v,1,:6780,key1",
	// 	"10e18,v,1,:6781,key2",
	// 	"10e18,v,1,:6782,key3",
	// 	"10e18,v,1,:6783,key3",
	// }
)

func Users(count int, initialE, userType, initialStake string, startingPort int) ([]*gengen.User, error) {
	var users []*gengen.User
	for i := startingPort; i < startingPort+count; i++ {

		portString := strconv.Itoa(i)
		u, err := gengen.ParseUser(strings.Join([]string{initialE, userType, initialStake, ":" + portString, "key" + portString}, ","))
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// Unfortunately we need to provide a genesis file here to be able to set the
// ethereum service on the node before starting but we can only find out the
// address port the node bound on till after starting if using the 0 port. This
// means that we have to predefine ports in the genesis, which could cause
// problems if anything is already bound on that port.
func Node(u *gengen.User, genesis *core.Genesis) (*node.Node, func(), error) {
	// Copy the base node config
	c := *baseNodeConfig

	// p2p key and address
	c.P2P.PrivateKey = u.Key.(*ecdsa.PrivateKey)
	c.P2P.ListenAddr = "0.0.0.0:" + strconv.Itoa(u.NodePort)

	// Set rpc ports
	userCount := len(genesis.Config.AutonityContractConfig.Users)
	c.HTTPPort = u.NodePort + userCount
	c.WSPort = u.NodePort + userCount*2

	datadir, err := ioutil.TempDir("", "autonity_datadir")
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		os.RemoveAll(datadir)
	}

	c.DataDir = datadir

	n, err := node.New(&c)
	if err != nil {
		cleanup()
		return nil, nil, err
	}

	// copy the base eth config
	ec := *baseEthConfig
	ec.Genesis = genesis
	ec.NetworkId = genesis.Config.ChainID.Uint64()
	ec.Tendermint = *genesis.Config.Tendermint

	// Register an injector on the node to provide the ethereum service.
	err = n.Register(func(ctx *node.ServiceContext) (node.Service, error) {
		return eth.New(ctx, baseEthConfig, nil)
	})
	if err != nil {
		cleanup()
		return nil, nil, err
	}

	// Now we need to initialise the db with the genesis
	chaindb, err := n.OpenDatabase("chaindata", 0, 0, "")
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	defer chaindb.Close()

	_, _, err = core.SetupGenesisBlock(chaindb, genesis)
	if err != nil {
		cleanup()
		return nil, nil, err
	}

	return n, cleanup, nil
}

func Genesis(users []*gengen.User) (*core.Genesis, error) {
	return gengen.NewGenesis(1, users)
}
