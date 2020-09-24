package test

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strconv"
	"sync"
	"time"

	ethereum "github.com/clearmatics/autonity"
	"github.com/clearmatics/autonity/cmd/gengen/gengen"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/eth"
	"github.com/clearmatics/autonity/eth/downloader"
	"github.com/clearmatics/autonity/ethclient"
	"github.com/clearmatics/autonity/log"
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

// Users returns 'count' users using the given formatString and starting port.
// The format string should have a string placeholder for the port and the key.
// E.G. for a validator '10e18,v,1,0.0.0.0:%s,%s'.
func Users(count int, formatString string, startingPort int) ([]*gengen.User, error) {
	var users []*gengen.User
	for i := startingPort; i < startingPort+count; i++ {

		portString := strconv.Itoa(i)
		u, err := gengen.ParseUser(fmt.Sprintf(formatString, portString, "key"+portString))
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// Node provides an enhanced interface to node.Node with useful additions, the
// *node.Node is embedded so that its api is available through Node.
type Node struct {
	*node.Node
	WsClient     *ethclient.Client
	Nonce        uint64
	Key          *ecdsa.PrivateKey
	Address      common.Address
	Tracker      *TransactionTracker
	SentTxs      []common.Hash
	ProcessedTxs []common.Hash
}

// Unfortunately we need to provide a genesis file here to be able to set the
// ethereum service on the node before starting but we can only find out the
// address port the node bound on till after starting if using the 0 port. This
// means that we have to predefine ports in the genesis, which could cause
// problems if anything is already bound on that port.
func NewNode(u *gengen.User, genesis *core.Genesis) (*Node, func(), error) {

	k := u.Key.(*ecdsa.PrivateKey)
	address := crypto.PubkeyToAddress(k.PublicKey)

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

	// Give this logger context based on the node address so that we can easily
	// trace single node execution in the logs.
	c.Logger = log.New("node", address.String()[2:7])

	n, err := node.New(&c)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	cleanup = func() {
		os.RemoveAll(datadir)
		n.Close()
	}

	// copy the base eth config
	ec := *baseEthConfig
	// Set the min gas price on the mining pool config, otherwise the miner
	// starts with a defalt min gas price. Which causes transactions to be
	// dropped.
	ec.Miner.GasPrice = (&big.Int{}).SetUint64(genesis.Config.AutonityContractConfig.MinGasPrice)
	ec.Genesis = genesis
	ec.NetworkId = genesis.Config.ChainID.Uint64()
	ec.Tendermint = *genesis.Config.Tendermint

	// Register an injector on the node to provide the ethereum service.
	err = n.Register(func(ctx *node.ServiceContext) (node.Service, error) {
		return eth.New(ctx, &ec, nil)
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

	node := &Node{
		Node:    n,
		Key:     k,
		Address: address,
	}
	return node, cleanup, nil
}

func (n *Node) Start() error {
	err := n.Node.Start()
	if err != nil {
		return err
	}
	n.WsClient, err = ethclient.Dial("ws://" + n.WSEndpoint())
	if err != nil {
		return err
	}
	n.Nonce, err = n.WsClient.PendingNonceAt(context.Background(), n.Address)
	return err
}

func (n *Node) SendETracked(ctx context.Context, recipient common.Address, value int64) error {
	err := n.SendE(ctx, recipient, value)
	if err != nil {
		return err
	}
	err = n.AwaitSentTransactions(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (n *Node) SendE(ctx context.Context, recipient common.Address, value int64) error {
	if n.Tracker == nil {
		t, err := TrackTransactions(n.WsClient)
		if err != nil {
			return err
		}
		n.Tracker = t
	}

	tx, err := ValueTransferTransaction(
		n.WsClient,
		n.Key,
		n.Address,
		recipient,
		n.Nonce,
		big.NewInt(value))

	if err != nil {
		return err
	}
	err = n.WsClient.SendTransaction(ctx, tx)
	if err != nil {
		return err
	}
	n.Nonce++
	n.SentTxs = append(n.SentTxs, tx.Hash())
	return nil
}

func (n *Node) AwaitSentTransactions(ctx context.Context) error {
	// Cleanup
	defer func() {
		n.Tracker.Close()
		n.Tracker = nil
	}()
	processed, err := n.Tracker.AwaitTransactions(ctx, n.SentTxs)
	for _, p := range processed {
		for i, s := range n.SentTxs {
			if p == s {
				// Remove from sent txs
				n.SentTxs = (append(n.SentTxs[:i], n.SentTxs[i+1:]...))
			}
		}
	}
	n.ProcessedTxs = append(n.ProcessedTxs, processed...)
	return err
}

type TransactionTracker struct {
	client *ethclient.Client
	heads  chan *types.Header
	sub    ethereum.Subscription
	wg     sync.WaitGroup
}

func TrackTransactions(client *ethclient.Client) (*TransactionTracker, error) {
	heads := make(chan *types.Header)
	// The subscription client will buffer 20000 notifications before closing
	// the subscription, if that happens the Err() chan will return
	// ErrSubscriptionQueueOverflow
	sub, err := client.SubscribeNewHead(context.Background(), heads)
	if err != nil {
		return nil, err
	}
	return &TransactionTracker{
		client: client,
		sub:    sub,
		heads:  heads,
	}, nil

}

func (tr *TransactionTracker) AwaitTransactions(ctx context.Context, hashes []common.Hash) (processed []common.Hash, err error) {
	hashmap := make(map[common.Hash]struct{}, len(hashes))
	for i := range hashes {
		hashmap[hashes[i]] = struct{}{}
	}
	tr.wg.Add(1)
	defer tr.wg.Done()
	for {
		select {
		case h := <-tr.heads:
			b, err := tr.client.BlockByHash(context.Background(), h.Hash())
			if err != nil {
				return processed, err
			}
			println("blocknumbrrrr", h.Number.String())

			for _, t := range b.Transactions() {
				h := t.Hash()
				if _, ok := hashmap[h]; ok {
					processed = append(processed, h)
				}
				if len(processed) == len(hashes) {
					return processed, nil
				}
			}
		case err := <-tr.sub.Err():
			// Will be nil if closed by calling Unsubscribe()
			if err == nil {
				err = errors.New("subscription closed")
			}
			return processed, err
		case <-ctx.Done():
			println("context done")
			return processed, ctx.Err()
		}
	}
}

func (tr *TransactionTracker) Close() { // How do I wait for stuff to finish here
	tr.sub.Unsubscribe()
	tr.wg.Wait()
}

func Genesis(users []*gengen.User) (*core.Genesis, error) {
	g, err := gengen.NewGenesis(1, users)
	if err != nil {
		return nil, err
	}
	// Make the tests fast
	g.Config.Tendermint.BlockPeriod = 0
	return g, nil
}

func ValueTransferTransaction(client *ethclient.Client, senderKey *ecdsa.PrivateKey, sender, recipient common.Address, nonce uint64, value *big.Int) (*types.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	// Figure out the gas allowance and gas price values
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to suggest gas price: %v", err)
	}

	msg := ethereum.CallMsg{From: sender, To: &recipient, GasPrice: gasPrice, Value: value}
	gasLimit, err := client.EstimateGas(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to estimate gas needed: %v", err)
	}

	// Create the transaction and sign it
	rawTx := types.NewTransaction(nonce, recipient, value, gasLimit, gasPrice, nil)
	signed, err := types.SignTx(rawTx, types.HomesteadSigner{}, senderKey)
	if err != nil {
		return nil, err
	}
	return signed, nil
}
