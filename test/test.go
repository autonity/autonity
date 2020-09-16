package test

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	ethereum "github.com/clearmatics/autonity"
	"github.com/clearmatics/autonity/accounts/abi/bind"
	"github.com/clearmatics/autonity/cmd/gengen/gengen"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/contracts/autonity"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/eth"
	"github.com/clearmatics/autonity/eth/downloader"
	"github.com/clearmatics/autonity/ethclient"
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

type TransactionTracker struct {
	client *ethclient.Client
	heads  chan *types.Header
	sub    ethereum.Subscription
	wg     sync.WaitGroup
}

func TrackTransactions(client *ethclient.Client) (*TransactionTracker, error) {
	heads := make(chan *types.Header)
	// The subscription client will buffer 20000 notifications before closing the subscription, if that happens the Err() chan will return ErrSubscriptionQueueOverflow
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

func (tr *TransactionTracker) AwaitTransactions(hashes []common.Hash) error {
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
				return err
			}

			for _, t := range b.Transactions() {
				delete(hashmap, t.Hash())
				if len(hashmap) == 0 {
					return nil
				}
			}
		case err := <-tr.sub.Err():
			// Will be nil if closed by calling Unsubscribe()
			return err

		}
	}
}

func (tr *TransactionTracker) Close() { // How do I wait for stuff to finish here
	tr.sub.Unsubscribe()
	tr.wg.Wait()
}

// type TransactionTracker struct {
// 	c                     *ethclient.Client
// 	completedTransactions map[common.Hash]struct{}
// 	mu                    sync.Mutex
// 	waiting               func(common.Hash, error)
// 	err                   error
// }

// func TrackTransactions(client *ethclient.Client) (*TransactionTracker, error) {
// 	tr := &TransactionTracker{
// 		c:                     client,
// 		completedTransactions: make(map[common.Hash]struct{}),
// 	}
// 	heads := make(chan *types.Header)
// 	sub, err := client.SubscribeNewHead(context.Background(), heads)
// 	if err != nil {
// 		return nil, err
// 	}
// 	go func() {
// 	Finished:
// 		for {
// 			select {
// 			case h := <-heads:
// 				b, err := client.BlockByHash(context.Background(), h.Hash())
// 				if err != nil {
// 					tr.mu.Lock()
// 					tr.err = err
// 					tr.mu.Unlock()
// 					break Finished
// 				}

// 				tr.mu.Lock()
// 				for _, t := range b.Transactions() {
// 					tr.completedTransactions[t.Hash()] = struct{}{}
// 					if tr.waiting != nil {
// 						tr.waiting(t.Hash(), nil) // notify waiting
// 					}
// 				}
// 				tr.mu.Unlock()
// 			case err, ok := <-sub.Err():
// 				if !ok {
// 					// Unsubscribe was called, the subscription is over
// 					return
// 				}
// 				if err != nil {
// 					tr.mu.Lock()
// 					tr.err = err
// 					tr.mu.Unlock()
// 					break Finished
// 				}

// 			}
// 		}
// 		if tr.err != nil {
// 			tr.mu.Lock()
// 			tr.waiting(common.Hash{}, tr.err)
// 			tr.mu.Unlock()
// 		}
// 	}()
// 	return tr, nil

// }

// func (tr *TransactionTracker) AwaitTransactions(hashes []common.Hash) error {
// 	hashmap := make(map[common.Hash]struct{}, len(hashes))
// 	for i := range hashes {
// 		hashmap[hashes[i]] = struct{}{}
// 	}

// 	var wg sync.WaitGroup
// 	wg.Add(1)

// 	tr.mu.Lock()
// 	// Remove already completed transactions
// 	for h := range tr.completedTransactions {
// 		delete(hashmap, h)
// 	}
// 	// Then register the waiter to wait for the remaining transactions
// 	tr.waiting = func(h common.Hash, err error) {
// 		if err != nil {
// 			wg.Done()
// 			return
// 		}
// 		delete(hashmap, h)
// 		if len(hashmap) == 0 {
// 			wg.Done()
// 		}
// 	}
// 	tr.mu.Unlock()
// 	wg.Wait()
// 	return tr.err
// }

// func AwaitTransactions(ctx context.Context, client *ethclient.Client, hashes []common.Hash) (<-chan common.Hash, <-chan error, error) {
// 	heads := make(chan *types.Header)
// 	sub, err := client.SubscribeNewHead(context.Background(), heads)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	hashes := make(chan common.Hash)
// 	errors := make(chan error)
// 	go func() {
// 		for {
// 			select {
// 			case h := <-heads:
// 				b, err := client.BlockByHash(context.Background(), h.Hash())
// 				if err != nil {
// 					errors <- err
// 				}
// 				for _, t := range b.Transactions() {
// 					hashes <- t.Hash()
// 				}
// 			case err := <-sub.Err():
// 				errors <- err
// 			case <-ctx.Done():
// 				return
// 			}
// 		}

// 	}()
// 	return hashes, errors, nil
// }
// func MinedTransactions(ctx context.Context, client *ethclient.Client) (<-chan common.Hash, <-chan error, error) {
// 	heads := make(chan *types.Header)
// 	sub, err := client.SubscribeNewHead(context.Background(), heads)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	hashes := make(chan common.Hash)
// 	errors := make(chan error)
// 	go func() {
// 		for {
// 			select {
// 			case h := <-heads:
// 				b, err := client.BlockByHash(context.Background(), h.Hash())
// 				if err != nil {
// 					errors <- err
// 				}
// 				for _, t := range b.Transactions() {
// 					hashes <- t.Hash()
// 				}
// 			case err := <-sub.Err():
// 				errors <- err
// 			case <-ctx.Done():
// 				return
// 			}
// 		}

// 	}()
// 	return hashes, errors, nil
// }

func sendTx(nodeAddr string, senderKey string, action func(*ethclient.Client, *bind.TransactOpts, *autonitybindings.Autonity) (*types.Transaction, error)) error {
	client, err := ethclient.Dial("ws://" + nodeAddr)
	if err != nil {
		return err
	}
	defer client.Close()
	a, err := autonitybindings.NewAutonity(autonity.ContractAddress, client)
	if err != nil {
		return err
	}
	k, err := crypto.HexToECDSA(senderKey)
	if err != nil {
		return fmt.Errorf("failed to decode sender key: %v", err)
	}

	senderAddress := crypto.PubkeyToAddress(k.PublicKey)
	opts := bind.NewKeyedTransactor(k)
	opts.From = senderAddress
	opts.Context = context.Background()

	heads := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), heads)
	if err != nil {
		return err
	}

	tx, err := action(client, opts, a)
	if err != nil {
		return err
	}

	txHash := tx.Hash()
	marshalled, err := json.MarshalIndent(tx, "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("Sent transaction:\n%s\n", string(marshalled))

	fmt.Print("\nAwaiting mining ...")

	for {
		select {
		case h := <-heads:
			b, err := client.BlockByHash(context.Background(), h.Hash())
			if err != nil {
				return err
			}
			for _, t := range b.Transactions() {
				if t.Hash() == txHash {
					r, err := client.TransactionReceipt(context.Background(), txHash)
					if err != nil {
						return err
					}
					marshalled, err := json.MarshalIndent(r, "", "  ")
					if err != nil {
						return err
					}
					fmt.Printf("\n\nTransaction receipt:\n%v\n", string(marshalled))
					return nil
				}
			}
			fmt.Print(".")

		case err := <-sub.Err():
			return err
		}
	}
}

func Genesis(users []*gengen.User) (*core.Genesis, error) {
	return gengen.NewGenesis(1, users)
}
