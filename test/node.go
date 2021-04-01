package test

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/clearmatics/autonity/consensus/tendermint/config"
	"io/ioutil"
	"math/big"
	"os"
	"strconv"
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

	baseEthConfig = &eth.Config{
		SyncMode:        downloader.FullSync,
		DatabaseCache:   256,
		DatabaseHandles: 256,
		TxPool:          core.DefaultTxPoolConfig,
	}
)

// Node provides an enhanced interface to node.Node with useful additions, the
// *node.Node is embedded so that its api is available through Node.
type Node struct {
	*node.Node
	Config    *node.Config
	Eth       *eth.Ethereum
	EthConfig *eth.Config
	WsClient  *ethclient.Client
	Nonce     uint64
	Key       *ecdsa.PrivateKey
	Address   common.Address
	Tracker   *TransactionTracker
	// The transactions that this node has sent.
	SentTxs []*types.Transaction
}

// NewNode creates a new running node as the given user with the provided
// genesis.
//
// Unfortunately we need to provide a genesis file here to be able to set the
// ethereum service on the node before starting but we can only find out the
// port the node bound on till after starting if using the 0 port. This means
// that we have to predefine ports in the genesis, which could cause problems
// if anything is already bound on that port.
func NewNode(u *gengen.User, genesis *core.Genesis) (*Node, error) {

	k := u.Key.(*ecdsa.PrivateKey)
	address := crypto.PubkeyToAddress(k.PublicKey)

	// Copy the base node config, so we can modify it without damaging the
	// original.
	c := &node.Config{}
	err := copyObject(baseNodeConfig, c)
	if err != nil {
		return nil, err
	}

	// p2p key and address
	c.P2P.PrivateKey = u.Key.(*ecdsa.PrivateKey)
	c.P2P.ListenAddr = "0.0.0.0:" + strconv.Itoa(u.NodePort)

	// Set rpc ports
	userCount := len(genesis.Config.AutonityContractConfig.Users)
	c.HTTPPort = u.NodePort + userCount
	c.WSPort = u.NodePort + userCount*2

	datadir, err := ioutil.TempDir("", "autonity_datadir")
	if err != nil {
		return nil, err
	}
	c.DataDir = datadir

	// copy the base eth config, so we can modify it without damaging the
	// original.
	ec := &eth.Config{}
	err = copyObject(baseEthConfig, ec)
	if err != nil {
		return nil, err
	}
	// Set the min gas price on the mining pool config, otherwise the miner
	// starts with a default min gas price. Which causes transactions to be
	// dropped.
	ec.Miner.GasPrice = (&big.Int{}).SetUint64(genesis.Config.AutonityContractConfig.MinGasPrice)
	ec.Genesis = genesis
	ec.NetworkId = genesis.Config.ChainID.Uint64()
	ec.Tendermint = *genesis.Config.Tendermint

	node := &Node{
		Config:    c,
		EthConfig: ec,
		Key:       k,
		Address:   address,
		Tracker:   NewTransactionTracker(),
	}

	return node, nil
}

// This creates the node.Node and eth.Ethereum and starts the node.Node and
// starts eth.Ethereum mining.
func (n *Node) Start() error {
	// Provide a copy of the config to node.New, so that we can rely on
	// Node.Config field not being manipulated by node and hence use our copy
	// for black box testing.
	nodeConfigCopy := &node.Config{}
	err := copyNodeConfig(n.Config, nodeConfigCopy)
	if err != nil {
		return err
	}

	// Give this logger context based on the node address so that we can easily
	// trace single node execution in the logs. We set the logger only on the
	// copy, since it is not useful for black box testing and it is also not
	// marshalable since the implementation contains unexported fields.
	nodeConfigCopy.Logger = log.New("node", n.Address.String()[2:7])
	// n.Config.P2P.PrivateKey = n.ConfigCopy.P2P.PrivateKey
	n.Node, err = node.New(nodeConfigCopy)
	if err != nil {
		return err
	}

	// This registers the ethereum service on the n.Node, so that calling
	// n.Node.Stop will also close the eth service. Again we provide a copy of
	// the EthConfig so that we can use our copy for black box testing.
	ethConfigCopy := &eth.Config{}
	err = copyObject(n.EthConfig, ethConfigCopy)
	if err != nil {
		return err
	}
	n.Eth, err = eth.New(n.Node, ethConfigCopy)
	if err != nil {
		return err
	}
	_, _, err = core.SetupGenesisBlock(n.Eth.ChainDb(), n.EthConfig.Genesis)
	if err != nil {
		return err
	}
	err = n.Node.Start()
	if err != nil {
		return err
	}
	n.WsClient, err = ethclient.Dial(n.WSEndpoint())
	if err != nil {
		return err
	}
	n.Nonce, err = n.WsClient.PendingNonceAt(context.Background(), n.Address)
	if err != nil {
		return err
	}
	err = n.Tracker.StartTracking(n.WsClient)
	if err != nil {
		return err
	}
	return n.Eth.StartMining(1)
}

// Close shuts down the node and releases all resources and removes the datadir
// unless an error is returned, in which case there is no guarantee that all
// resources are released.
func (n *Node) Close() error {
	err := n.Tracker.StopTracking()
	if err != nil {
		return err
	}
	n.WsClient.Close()
	if n.Node != nil {
		err = n.Node.Close() // This also shuts down the Eth service
	}
	os.RemoveAll(n.Config.DataDir)
	return err
}

// SendETracked functions like SendE but also waits for the transaction to be processed.
func (n *Node) SendETracked(ctx context.Context, recipient common.Address, value int64) error {
	tx, err := n.SendE(ctx, recipient, value)
	if err != nil {
		return err
	}
	return n.AwaitTransactions(ctx, tx)
}

// SendE creates a value transfer transaction to send E to the recipient.
func (n *Node) SendE(ctx context.Context, recipient common.Address, value int64) (*types.Transaction, error) {
	tx, err := ValueTransferTransaction(
		n.WsClient,
		n.Key,
		n.Address,
		recipient,
		n.Nonce,
		big.NewInt(value))

	if err != nil {
		return nil, err
	}
	err = n.WsClient.SendTransaction(ctx, tx)
	if err != nil {
		return nil, err
	}
	n.Nonce++
	n.SentTxs = append(n.SentTxs, tx)
	return tx, nil
}

// AwaitTransactions awaits all the provided transactions.
func (n *Node) AwaitTransactions(ctx context.Context, txs ...*types.Transaction) error {
	sentHashes := make([]common.Hash, len(txs))
	for i, tx := range txs {
		sentHashes[i] = tx.Hash()
	}
	return n.Tracker.AwaitTransactions(ctx, sentHashes)
}

// AwaitSentTransactions awaits all the transactions that this node has sent.
func (n *Node) AwaitSentTransactions(ctx context.Context) error {
	return n.AwaitTransactions(ctx, n.SentTxs...)
}

func (n *Node) ProcessedTxBlock(tx *types.Transaction) *types.Block {
	return n.Tracker.GetProcessedBlock(tx.Hash())
}

// TxFee returns the gas fee for the given transaction.
func (n *Node) TxFee(ctx context.Context, tx *types.Transaction) (*big.Int, error) {
	r, err := n.WsClient.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		return nil, err
	}
	return big.NewInt(0).Mul(new(big.Int).SetUint64(r.GasUsed), tx.GasPrice()), nil
}

// Network represents a network of nodes and provides funtionality to easily
// create, start and stop a collection of nodes.
type Network []*Node

// NewNetworkWithMaliciousUser generate a network of nodes that are running and
// mining with a block period 1 second. The rule ID is ID of the malicious behaviour.
// The malicious user will do a one-time shot misbehaviour.
func NewNetworkWithMaliciousUser(users []*gengen.User, ruleID uint8) (Network, error) {
	g, err := Genesis(users)
	if err != nil {
		return nil, err
	}
	network := make([]*Node, len(users))
	for i, u := range users {
		n, err := NewNode(u, g)
		if err != nil {
			return nil, fmt.Errorf("failed to build node for network: %v", err)
		}

		// set block period to 1 second since otherwise it cause none necessary accusations.
		n.EthConfig.Tendermint.BlockPeriod = 1
		n.EthConfig.Tendermint.MisbehaveConfig = &config.MaliciousConfig{RuleID: ruleID}

		if err := n.Start(); err != nil {
			return nil, fmt.Errorf("failed to start node for network: %v", err)
		}

		network[i] = n
	}
	// There is a race condition in miner.worker its field snapshotBlock is set
	// only when new transacting are received or commitNewWork is called. But
	// both of these happen in goroutines separate to the call to miner.Start
	// and miner.Start does not wait for snapshotBlock to be set. Therefore
	// there is currently no way to know when it is safe to call estimate gas.
	// What we do here is sleep a bit and cross our fingers.
	time.Sleep(10 * time.Millisecond)
	return network, nil
}

// NewNetworkFromUsers generates a network of nodes that are running and
// mining. For each provided user a corresponding node is created. If there is
// an error it will be returned immediately, meaning that some nodes may be
// running and others not.
func NewNetworkFromUsers(users []*gengen.User) (Network, error) {
	g, err := Genesis(users)
	if err != nil {
		return nil, err
	}
	network := make([]*Node, len(users))
	for i, u := range users {
		n, err := NewNode(u, g)
		if err != nil {
			return nil, fmt.Errorf("failed to build node for network: %v", err)
		}
		if err := n.Start(); err != nil {
			return nil, fmt.Errorf("failed to start node for network: %v", err)
		}

		network[i] = n
	}
	// There is a race condition in miner.worker its field snapshotBlock is set
	// only when new transacting are received or commitNewWork is called. But
	// both of these happen in goroutines separate to the call to miner.Start
	// and miner.Start does not wait for snapshotBlock to be set. Therefore
	// there is currently no way to know when it is safe to call estimate gas.
	// What we do here is sleep a bit and cross our fingers.
	time.Sleep(10 * time.Millisecond)
	return network, nil
}

// NewNetwork generates a network of nodes that are running, but not mining.
// For an explanation of the parameters see 'Users'.
func NewNetwork(count int, formatString string, startingPort int) (Network, error) {
	users, err := Users(count, formatString, startingPort)
	if err != nil {
		return nil, fmt.Errorf("failed to build users: %v", err)
	}
	return NewNetworkFromUsers(users)
}

// AwaitTransactions ensures that the entire network has processed the provided transactions.
func (n Network) AwaitTransactions(ctx context.Context, txs ...*types.Transaction) error {
	for _, node := range n {
		err := node.AwaitTransactions(ctx, txs...)
		if err != nil {
			return err
		}
	}
	return nil
}

// Shutdown closes all nodes in the network, any errors that are encounter are
// printed to stdout.
func (n Network) Shutdown() {
	for _, node := range n {
		if node != nil {
			err := node.Close()
			if err != nil {
				fmt.Printf("error shutting down node %v: %v", node.Address.String(), err)
			}
		}
	}
}

// ValueTransferTransaction builds a signed value transfer transaction from the
// sender to the recipient with the given value and nonce, it uses the client
// to suggest a gas price and to estimate the gas.
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

// Users returns 'count' users using the given formatString and starting port.
// The format string should have a string placeholder for the port and the key.
// The format string should follow the format defined for users in the gengen
// package see the variable 'userDescription' in the gengen package for a
// detailed description of the meaning of the format string.
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

// Genesis creates a genesis instance from the provided users.
func Genesis(users []*gengen.User) (*core.Genesis, error) {
	g, err := gengen.NewGenesis(1, users)
	if err != nil {
		return nil, err
	}
	// Make the tests fast
	g.Config.Tendermint.BlockPeriod = 0
	return g, nil
}

// Since the node config is not marshalable by default we construct a
// marshalable struct which we marshal and unmarshal and then unpack into the
// original struct type.
func copyNodeConfig(source, dest *node.Config) error {
	s := &MarshalableNodeConfig{}
	s.Config = *source
	p := MarshalableP2PConfig{}
	p.Config = source.P2P

	crypto.FromECDSA(source.P2P.PrivateKey)

	p.PrivateKey = (*MarshalableECDSAPrivateKey)(source.P2P.PrivateKey)
	s.P2P = p
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	u := new(MarshalableNodeConfig)
	err = json.Unmarshal(data, u)
	if err != nil {
		return err
	}
	*dest = u.Config
	dest.P2P = u.P2P.Config
	dest.P2P.PrivateKey = (*ecdsa.PrivateKey)(u.P2P.PrivateKey)
	return nil
}

type MarshalableNodeConfig struct {
	node.Config
	P2P MarshalableP2PConfig
}

type MarshalableP2PConfig struct {
	p2p.Config
	PrivateKey *MarshalableECDSAPrivateKey
}

type MarshalableECDSAPrivateKey ecdsa.PrivateKey

func (k *MarshalableECDSAPrivateKey) UnmarshalJSON(b []byte) error {
	key, err := crypto.PrivECDSAFromHex(b[1 : len(b)-1])
	if err != nil {
		return err
	}
	*k = MarshalableECDSAPrivateKey(*key)
	return nil
}

func (k *MarshalableECDSAPrivateKey) MarshalJSON() ([]byte, error) {
	return []byte(`"` + hex.EncodeToString(crypto.FromECDSA((*ecdsa.PrivateKey)(k))) + `"`), nil
}

// copyObject copies an object so that the copy shares no memory with the
// original.
func copyObject(source, dest interface{}) error {
	data, err := json.Marshal(source)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}
