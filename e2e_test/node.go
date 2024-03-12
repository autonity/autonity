package e2e

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"net"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/consul/sdk/freeport"

	"github.com/autonity/autonity/consensus/acn"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/eth/downloader"
	"github.com/autonity/autonity/p2p/enode"

	ethereum "github.com/autonity/autonity"
	"github.com/autonity/autonity/cmd/gengen/gengen"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/eth"
	"github.com/autonity/autonity/eth/ethconfig"
	"github.com/autonity/autonity/ethclient"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/node"
	"github.com/autonity/autonity/p2p"
	"github.com/autonity/autonity/params"
)

const (
	localhost = "127.0.0.1"
)

var (
	baseNodeConfig = &node.Config{
		Name:    "autonity",
		Version: params.Version,
		ExecutionP2P: p2p.Config{
			MaxPeers: 50,
		},
		ConsensusP2P: p2p.Config{
			MaxPeers: 100000,
		},
		HTTPHost: localhost,
		WSHost:   localhost,
	}

	terminalColors = []struct {
		foreground string
		background string
	}{{
		foreground: log.White,
		background: log.BackgroundBlack,
	}, {
		foreground: log.Black,
		background: log.BackgroundLightCyan,
	}, {
		foreground: log.Black,
		background: log.BackgroundLightYellow,
	}, {
		foreground: log.Black,
		background: log.BackgroundCyan,
	}, {
		foreground: log.Black,
		background: log.BackgroundLightGreen,
	}}
)

// Node provides an enhanced interface to node.Node with useful additions, the
// *node.Node is embedded so that its api is available through Node.
type Node struct {
	*node.Node
	isRunning    bool
	Config       *node.Config
	Eth          *eth.Ethereum
	EthConfig    *ethconfig.Config
	WsClient     *ethclient.Client
	Nonce        uint64
	Key          *ecdsa.PrivateKey
	ConsensusKey blst.SecretKey
	Address      common.Address
	Tracker      *TransactionTracker
	// The transactions that this node has sent.
	SentTxs     []*types.Transaction
	CustHandler *interfaces.Services
	ID          int
}

// NewNode creates a new running node as the given user with the provided
// genesis.
//
// Unfortunately we need to provide a genesis file here to be able to set the
// ethereum service on the node before starting but we can only find out the
// port the node bound on till after starting if using the 0 port. This means
// that we have to predefine ports in the genesis, which could cause problems
// if anything is already bound on that port.
func NewNode(validator *gengen.Validator, genesis *core.Genesis, id int) (*Node, error) {
	address := crypto.PubkeyToAddress(validator.NodeKey.PublicKey)

	// Copy the base node config, so we can modify it without damaging the
	// original.
	nodeConfig := copyNodeConfig(baseNodeConfig)

	// p2p key and address
	nodeConfig.ExecutionP2P.PrivateKey = validator.NodeKey
	nodeConfig.ExecutionP2P.ListenAddr = fmt.Sprintf("%s:%d", localhost, validator.NodePort)

	// consensus key used by consensus engine.
	nodeConfig.ConsensusKey = validator.ConsensusKey
	nodeConfig.ConsensusP2P.PrivateKey = validator.NodeKey
	nodeConfig.ConsensusP2P.ListenAddr = fmt.Sprintf("%s:%d", localhost, validator.AcnPort)

	// Set rpc ports
	//c.HTTPPort = freeport.GetOne(t)
	//c.WSPort = freeport.GetOne(t)

	nodeConfig.DataDir = ""

	// copy the base eth config, so we can modify it without damaging the
	// original.
	ethConfig := &ethconfig.Config{}
	ethconfig.Defaults.SyncMode = downloader.FullSync
	ethconfig.Defaults.Miner.Recommit = time.Second
	if err := copyConfig(&ethconfig.Defaults, ethConfig); err != nil {
		return nil, err
	}
	// Set the min gas price on the mining pool config, otherwise the miner
	// starts with a default min gas price. Which causes transactions to be
	// dropped.
	ethConfig.Miner.GasPrice = (&big.Int{}).SetUint64(genesis.Config.AutonityContractConfig.MinBaseFee)
	ethConfig.Miner.Etherbase = crypto.PubkeyToAddress(validator.NodeKey.PublicKey)
	ethConfig.Genesis = genesis
	ethConfig.NetworkID = genesis.Config.ChainID.Uint64()

	// Give this logger context based on the node address so that we can easily
	// trace single node execution in the logs. We set the logger only on the
	// copy, since it is not useful for black box testing and it is also not
	// marshalable since the implementation contains unexported fields.
	logger := log.NewGlogHandler(log.StreamHandler(os.Stderr, log.FormatFunc(func(record *log.Record) []byte {
		b := log.TerminalFormat(false).Format(record)
		if id < len(terminalColors) {
			prefix := []byte(terminalColors[id].background + terminalColors[id].foreground)
			suffix := []byte("\x1b[0;K\033[0m\n")
			return append(append(prefix, b[:len(b)-1]...), suffix...)
		}
		return b
	})))

	logger.Verbosity(log.DefaultVerbosity)
	nodeConfig.Logger = log.New()
	nodeConfig.Logger.SetHandler(logger)

	// set custom tendermint services
	nodeConfig.SetTendermintServices(validator.TendermintServices)

	n := &Node{
		Config:       nodeConfig,
		EthConfig:    ethConfig,
		Key:          validator.NodeKey,
		ConsensusKey: validator.ConsensusKey,
		Address:      address,
		Tracker:      NewTransactionTracker(),
		CustHandler:  validator.TendermintServices,
		ID:           id,
	}

	return n, nil
}

func (n *Node) Running() bool {
	return n.isRunning
}

func (n *Node) Start() error {
	if n.isRunning {
		return nil
	}
	var err error
	defer func() {
		if err == nil {
			n.isRunning = true
		}
	}()

	if n.Node, err = node.New(n.Config); err != nil {
		return err
	}
	if n.Eth, err = eth.New(n.Node, n.EthConfig); err != nil {
		return fmt.Errorf("cannot create new eth: %w", err)
	}

	acn.New(n.Node, n.Eth, ethconfig.Defaults.NetworkID)
	if err = n.Node.Start(); err != nil {
		return fmt.Errorf("failed to start a node: %w", err)
	}
	if n.WsClient, err = ethclient.Dial(n.WSEndpoint()); err != nil {
		return err
	}
	if n.Nonce, err = n.WsClient.PendingNonceAt(context.Background(), n.Address); err != nil {
		return err
	}
	err = n.Tracker.StartTracking(n.WsClient)
	return err
}

// Close shuts down the node and releases all resources and removes the datadir
// unless an error is returned, in which case there is no guarantee that all
// resources are released.
func (n *Node) Close(deleteDataDir bool) error {
	if !n.isRunning {
		return nil
	}
	var err error
	defer func() {
		if err == nil {
			n.isRunning = false
		}
	}()
	err = n.Tracker.StopTracking()
	if err != nil {
		return err
	}
	n.WsClient.Close()
	if n.Node != nil {
		err = n.Node.Close() // This also shuts down the Eth service
	}
	if deleteDataDir {
		os.RemoveAll(n.Config.DataDir)
	}
	return err
}

// SendAUTtracked functions like SendAUT but also waits for the transaction to be processed.
func (n *Node) SendAUTtracked(ctx context.Context, recipient common.Address, value int64) error {
	tx, err := n.SendAUT(ctx, recipient, value)
	if err != nil {
		return err
	}
	return n.AwaitTransactions(ctx, tx)
}

// SendAUT creates a value transfer transaction to send E to the recipient.
func (n *Node) SendAUT(ctx context.Context, recipient common.Address, value int64) (*types.Transaction, error) {
	tx, err := ValueTransferTransaction(
		n.WsClient,
		n.Key,
		n.Address,
		recipient,
		n.Nonce,
		n.EthConfig,
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

func (n *Node) GetChainHeight() uint64 {
	return n.Eth.BlockChain().CurrentHeader().Number.Uint64()
}

func (n *Node) IsSyncComplete() bool {
	syncResult := n.Eth.APIBackend.SyncProgress()
	return syncResult.CurrentBlock >= syncResult.HighestBlock
}

// Network represents a network of nodes and provides funtionality to easily
// create, start and stop a collection of nodes.
type Network []*Node

// WaitForSyncComplete waits for sync to be completed
// for all running nodes in the quorum
func (nw Network) WaitForSyncComplete() error {
	// we will wait maximum one minute for All nodes be synced completely
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	// ticker to periodically check for sync
	syncTicker := time.NewTicker(1 * time.Second)
	opCh := make(chan error)
	quit := make(chan bool)
	count := 0
	for _, n := range nw {
		if !n.isRunning {
			continue
		}
		// count of the all spawned goroutines
		count++
		go func(n *Node) {
			for {
				select {
				case <-syncTicker.C:
					if n.IsSyncComplete() {
						opCh <- nil
						return
					}
					// context expired, send error on error channel
				case <-ctx.Done():
					opCh <- ctx.Err()
					return
				case <-quit:
					return
				}
			}
		}(n)
	}

	// return if none of the nodes are running
	if count == 0 {
		return nil
	}
	for err := range opCh {
		if err != nil {
			// we will close the quit to channel to signal
			// all goroutines to exit before returning error
			close(quit)
			return err
		}
		count--
		// We have received from all goroutines
		if count == 0 {
			return nil
		}
	}
	return nil
}

func (nw Network) isNetworkLive(chainHeights []uint64) bool {
	// compare the current chain heights with the previously recorded chain height
	for i, n := range nw {
		// skipping nodes which are not running
		if !n.isRunning {
			continue
		}
		currHeight := n.Eth.BlockChain().CurrentHeader().Number.Uint64()
		if currHeight <= chainHeights[i] {
			// this node is not mining blocks with in the block period
			return false
		}
	}
	return true
}

// WaitForNetworkToStartMining waits for all nodes to advance
// their chain heights after new blocks are mined
func (nw Network) WaitForNetworkToStartMining() error {
	// we will wait maximum one minute for network to be live again
	// and start mining blocks
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	// cache current chain height for all nodes
	chainHeights := make([]uint64, len(nw))
	runningCount := 0
	for i, n := range nw {
		if n.isRunning {
			runningCount++
			chainHeights[i] = n.Eth.BlockChain().CurrentHeader().Number.Uint64()
		}
	}
	// return if none of the nodes are running
	if runningCount == 0 {
		return fmt.Errorf("can't mine new blocks, there are no running nodes in the quorum")
	}
	syncTicker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-syncTicker.C:
			if nw.isNetworkLive(chainHeights) {
				return nil
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// WaitToMineNBlocks waits for network to mine given number of
// blocks in the given time window default value for numSec can be kept 60 seconds
// if verifyRate == true --> we return an error if we cannot satisfy that 1 block/s rate
func (nw Network) WaitToMineNBlocks(numBlocks uint64, numSec int, verifyRate bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(numSec)*time.Second)
	defer cancel()
	// cache current chain height for all nodes
	chainHeights := make([]uint64, len(nw))
	lastHeights := make([]uint64, len(nw))
	for i, n := range nw {
		if n.isRunning {
			chainHeights[i] = n.Eth.BlockChain().CurrentHeader().Number.Uint64()
			lastHeights[i] = chainHeights[i]
		}
	}
	syncTicker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-syncTicker.C:
			totalRunning := 0
			syncedNodes := 0
			for i, n := range nw {
				// skipping nodes which are not running
				if !n.isRunning {
					continue
				}
				currHeader := n.Eth.BlockChain().CurrentHeader()
				currHeight := currHeader.Number.Uint64()
				if currHeight > chainHeights[i]+numBlocks {
					syncedNodes++
				}
				totalRunning++

				// verify block rate against parent if we moved forward
				// it is not bulletproof but good enough
				// (we could have moved 2 blocks from last iteration, with first block not respecting the rate and second yes)
				if verifyRate && currHeight > lastHeights[i] {
					currTime := currHeader.Time
					parentTime := n.Eth.BlockChain().GetHeaderByHash(currHeader.ParentHash).Time
					if currTime-parentTime != 1 {
						return fmt.Errorf("Block rate not respected. parentTime: %d, currTime: %d", parentTime, currTime)
					}
				}
				lastHeights[i] = currHeight
			}
			// all the running nodes should reach the required chainHeight
			if syncedNodes == totalRunning {
				return nil
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// WaitForHeight waits for all nodes in the network to mine at least a given height
func (nw Network) WaitForHeight(height uint64, numSec int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(numSec)*time.Second)
	defer cancel()
	// cache current chain height for all nodes
	syncTicker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-syncTicker.C:
			totalRunning := 0
			syncedNodes := 0
			for _, n := range nw {
				// skipping nodes which are not running
				if !n.isRunning {
					continue
				}
				currHeight := n.Eth.BlockChain().CurrentHeader().Number.Uint64()
				if currHeight >= height {
					syncedNodes++
				}
				totalRunning++
			}
			// all the running nodes should reach the required chainHeight
			if syncedNodes == totalRunning {
				return nil
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// NewNetworkFromValidators generates a network of nodes that are running and
// mining. For each provided user a corresponding node is created. If there is
// an error it will be returned immediately, meaning that some nodes may be
// running and others not.
func NewNetworkFromValidators(t *testing.T, validators []*gengen.Validator, start bool, options ...gengen.GenesisOption) (Network, error) {
	g, err := Genesis(validators, options...)
	if err != nil {
		return nil, fmt.Errorf("failed the genesis: %w", err)
	}
	network := make([]*Node, len(validators))
	for i, u := range validators {
		n, err := NewNode(u, g, i)
		if len(validators) > 21 {
			n.EthConfig.DatabaseCache = 16
			n.EthConfig.DatabaseHandles = 8
		}
		if err != nil {
			return nil, fmt.Errorf("failed to build node for network: %v", err)
		}

		if start {
			err = n.Start()
			if err != nil {
				return nil, fmt.Errorf("failed to start node for network: %v", err)
			}
		}
		network[i] = n
	}
	go communicatePort(network[0].Config.WSPort)
	// There is a race condition in miner.worker its field snapshotBlock is set
	// only when new transactions are received or commitNewWork is called. But
	// both of these happen in goroutines separate to the call to miner.Start
	// and miner.Start does not wait for snapshotBlock to be set. Therefore
	// there is currently no way to know when it is safe to call estimate gas.
	// What we do here is sleep a bit and cross our fingers.
	time.Sleep(10 * time.Millisecond)
	return network, nil
}

// NewNetwork generates a network of nodes that are running and mining.
// For an explanation of the parameters see 'Validators'.
func NewNetwork(t *testing.T, count int, formatString string) (Network, error) {
	users, err := Validators(t, count, formatString)
	if err != nil {
		return nil, fmt.Errorf("failed to build users: %v", err)
	}
	return NewNetworkFromValidators(t, users, true)
}

type pipeManager struct {
	nodes   sync.Map
	network p2p.Network
}
type pipeDialer struct {
	node    *Node
	manager *pipeManager
	count   uint64
	fail    uint64
}

func newPipeManager(net p2p.Network) *pipeManager {
	return &pipeManager{network: net}
}

func (pm *pipeManager) createPipeDialer(node *Node) *pipeDialer {
	return &pipeDialer{
		node:    node,
		manager: pm,
	}
}

// Dial implements the p2p.NodeDialer interface by connecting to the node using
// an in-memory net.Pipe
func (p *pipeDialer) Dial(_ context.Context, dest *enode.Node) (conn net.Conn, err error) {
	p.count++
	if p.node.ID == 0 {
		fmt.Println("attempt", "cs", p.count, "f", p.fail, "type", p.manager.network)
	}
	n, ok := p.manager.nodes.Load(dest.ID())
	if !ok || !n.(*Node).Running() {
		// try again a bit later, the node may not have started yet
		<-time.After(10 * time.Millisecond)
		n, ok = p.manager.nodes.Load(dest.ID())
		if !ok || !n.(*Node).Running() {
			p.fail++
			return nil, fmt.Errorf("node not running: %s", dest.ID())
		}
	}
	pipe1, pipe2 := net.Pipe()
	go func() {
		switch p.manager.network {
		case p2p.Execution:
			n.(*Node).Node.ExecutionServer().SetupConn(pipe1, 4, nil)
		case p2p.Consensus:
			n.(*Node).Node.ConsensusServer().SetupConn(pipe1, 4, nil)
		}
	}()
	return pipe2, nil
}

func NewInMemoryNetwork(t *testing.T, validators []*gengen.Validator, start bool, options ...gengen.GenesisOption) (Network, error) {
	g, err := Genesis(validators, options...)
	if err != nil {
		return nil, fmt.Errorf("failed the genesis: %w", err)
	}
	network := make([]*Node, len(validators))
	executionManager := newPipeManager(p2p.Execution)
	consensusManager := newPipeManager(p2p.Consensus)
	bootnode1, _ := enode.Parse(enode.ValidSchemes, g.Config.AutonityContractConfig.Validators[0].Enode)
	baseNodeConfig.ExecutionP2P.BootstrapNodes = []*enode.Node{bootnode1}

	wg := sync.WaitGroup{}
	for i, u := range validators {
		wg.Add(1)
		go func(id int, val *gengen.Validator) {
			n, _ := NewNode(val, g, id)
			if id == 0 {
				n.Config.WSPort = freeport.GetOne(t)
			}
			nodeID := enode.PubkeyToIDV4(&val.NodeKey.PublicKey)
			n.Config.ConsensusP2P.Dialer = consensusManager.createPipeDialer(n)
			n.Config.ExecutionP2P.Dialer = executionManager.createPipeDialer(n)
			executionManager.nodes.Store(nodeID, n)
			consensusManager.nodes.Store(nodeID, n)
			network[id] = n
			wg.Done()
		}(i, u)
	}
	wg.Wait()

	startCh := make(chan error)
	if start {
		for _, n := range network {
			n := n
			go func() {
				startCh <- n.Start()
			}()
		}
	}
	for range network {
		err := <-startCh
		if err != nil {
			t.Fatalf("failed to start node with error %v", err)
		}
	}

	go communicatePort(network[0].Config.WSPort)
	// There is a race condition in miner.worker its field snapshotBlock is set
	// only when new transactions are received or commitNewWork is called. But
	// both of these happen in goroutines separate to the call to miner.Start
	// and miner.Start does not wait for snapshotBlock to be set. Therefore
	// there is currently no way to know when it is safe to call estimate gas.
	// What we do here is sleep a bit and cross our fingers.
	time.Sleep(10 * time.Millisecond)
	return network, nil
}

// AwaitTransactions ensures that the entire network has processed the provided transactions.
func (nw Network) AwaitTransactions(ctx context.Context, txs ...*types.Transaction) error {
	for _, node := range nw {
		err := node.AwaitTransactions(ctx, txs...)
		if err != nil {
			return err
		}
	}
	return nil
}

// Shutdown closes all nodes in the network, any errors that are encounter are
// printed to stdout.
func (nw Network) Shutdown() {
	for _, node := range nw {
		if node != nil && node.isRunning {
			err := node.Close(true)
			if err != nil {
				fmt.Printf("error shutting down node %v: %v", node.Address.String(), err)
			}
		}
	}
}

// ValueTransferTransaction builds a signed value transfer transaction from the
// sender to the recipient with the given value and nonce, it uses the client
// to suggest a gas price and to estimate the gas.
func ValueTransferTransaction(client *ethclient.Client,
	senderKey *ecdsa.PrivateKey,
	sender, recipient common.Address,
	nonce uint64,
	ethConfig *ethconfig.Config,
	value *big.Int) (*types.Transaction, error) {
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
	signed, err := types.SignTx(rawTx, types.LatestSigner(ethConfig.Genesis.Config), senderKey)
	if err != nil {
		return nil, err
	}
	return signed, nil
}

// Validators returns 'count' users using the given formatString and starting port.
// The format string should have a string placeholder for the port and the key.
// The format string should follow the format defined for users in the gengen
// package see the variable 'userDescription' in the gengen package for a
// detailed description of the meaning of the format string.
// E.G. for a validator '10e18,v,1,0.0.0.0:%s,%s,%s,%s'.
func Validators(t *testing.T, count int, formatString string) ([]*gengen.Validator, error) {
	var validators []*gengen.Validator
	for i := 0; i < count; i++ {
		portString := strconv.Itoa(freeport.GetOne(t))
		u, err := gengen.ParseValidator(fmt.Sprintf(formatString, portString, "key"+portString))
		if err != nil {
			return nil, err
		}
		//add port ip for consensus channel
		u.AcnIP = u.NodeIP
		u.AcnPort = freeport.GetOne(t)
		u.TreasuryKey, _ = crypto.GenerateKey()
		validators = append(validators, u)
	}
	return validators, nil
}

// This is used by the monitor tool to retrieve a useful websocket port
func communicatePort(port int) {
	conn, err := net.Dial("tcp", "localhost:55000")
	if err != nil {
		return
	}
	conn.Write([]byte(strconv.Itoa(port)))
	conn.Close()
}

// Genesis creates a genesis instance from the provided users.
func Genesis(users []*gengen.Validator, options ...gengen.GenesisOption) (*core.Genesis, error) {
	g, err := gengen.NewGenesis(users, options...)
	if err != nil {
		return nil, err
	}
	// Make the tests fast
	if err := g.Config.AutonityContractConfig.Prepare(); err != nil {
		return nil, err
	}
	return g, nil
}

// copyConfig copies an object so that the copy shares no memory with the
// original.
func copyConfig(source, dest *ethconfig.Config) error {
	*dest = *source
	return nil
}

// copyNodeConfig deep copy a node config.
func copyNodeConfig(source *node.Config) *node.Config {
	return &node.Config{
		ConsensusKey: source.ConsensusKey,
		Name:         source.Name,
		UserIdent:    source.UserIdent,
		Version:      source.Version,
		DataDir:      source.DataDir,
		ExecutionP2P: p2p.Config{
			PrivateKey:       source.ExecutionP2P.PrivateKey,
			MaxPeers:         source.ExecutionP2P.MaxPeers,
			MaxPendingPeers:  source.ExecutionP2P.MaxPendingPeers,
			DialRatio:        source.ExecutionP2P.DialRatio,
			NoDiscovery:      source.ExecutionP2P.NoDiscovery,
			DiscoveryV5:      source.ExecutionP2P.DiscoveryV5,
			Name:             source.ExecutionP2P.Name,
			BootstrapNodes:   source.ExecutionP2P.BootstrapNodes,
			BootstrapNodesV5: source.ExecutionP2P.BootstrapNodesV5,
			StaticNodes:      source.ExecutionP2P.StaticNodes,
			TrustedNodes:     source.ExecutionP2P.TrustedNodes,
			NetRestrict:      source.ExecutionP2P.NetRestrict,
			NodeDatabase:     source.ExecutionP2P.NodeDatabase,
			Protocols:        source.ExecutionP2P.Protocols,
			ListenAddr:       source.ExecutionP2P.ListenAddr,
			NAT:              source.ExecutionP2P.NAT,
			Dialer:           source.ExecutionP2P.Dialer,
			NoDial:           source.ExecutionP2P.NoDial,
			EnableMsgEvents:  source.ExecutionP2P.EnableMsgEvents,
			Logger:           source.ExecutionP2P.Logger,
			IsRated:          source.ExecutionP2P.IsRated,
			InRate:           source.ExecutionP2P.InRate,
			OutRate:          source.ExecutionP2P.OutRate,
		},
		ConsensusP2P: p2p.Config{
			PrivateKey:       source.ConsensusP2P.PrivateKey,
			MaxPeers:         source.ConsensusP2P.MaxPeers,
			MaxPendingPeers:  source.ConsensusP2P.MaxPendingPeers,
			DialRatio:        source.ConsensusP2P.DialRatio,
			NoDiscovery:      source.ConsensusP2P.NoDiscovery,
			DiscoveryV5:      source.ConsensusP2P.DiscoveryV5,
			Name:             source.ConsensusP2P.Name,
			BootstrapNodes:   source.ConsensusP2P.BootstrapNodes,
			BootstrapNodesV5: source.ConsensusP2P.BootstrapNodesV5,
			StaticNodes:      source.ConsensusP2P.StaticNodes,
			TrustedNodes:     source.ConsensusP2P.TrustedNodes,
			NetRestrict:      source.ConsensusP2P.NetRestrict,
			NodeDatabase:     source.ConsensusP2P.NodeDatabase,
			Protocols:        source.ConsensusP2P.Protocols,
			ListenAddr:       source.ConsensusP2P.ListenAddr,
			NAT:              source.ConsensusP2P.NAT,
			Dialer:           source.ConsensusP2P.Dialer,
			NoDial:           source.ConsensusP2P.NoDial,
			EnableMsgEvents:  source.ConsensusP2P.EnableMsgEvents,
			Logger:           source.ConsensusP2P.Logger,
			IsRated:          source.ConsensusP2P.IsRated,
			InRate:           source.ConsensusP2P.InRate,
			OutRate:          source.ConsensusP2P.OutRate,
		},
		KeyStoreDir:           source.KeyStoreDir,
		ExternalSigner:        source.ExternalSigner,
		UseLightweightKDF:     source.UseLightweightKDF,
		InsecureUnlockAllowed: source.InsecureUnlockAllowed,
		NoUSB:                 source.NoUSB,
		USB:                   source.USB,
		SmartCardDaemonPath:   source.SmartCardDaemonPath,
		IPCPath:               source.IPCPath,
		HTTPHost:              source.HTTPHost,
		HTTPPort:              source.HTTPPort,
		HTTPCors:              source.HTTPCors,
		HTTPVirtualHosts:      source.HTTPVirtualHosts,
		HTTPModules:           source.HTTPModules,
		HTTPTimeouts:          source.HTTPTimeouts,
		HTTPPathPrefix:        source.HTTPPathPrefix,
		WSHost:                source.WSHost,
		WSPort:                source.WSPort,
		WSPathPrefix:          source.WSPathPrefix,
		WSOrigins:             source.WSOrigins,
		WSModules:             source.WSModules,
		WSExposeAll:           source.WSExposeAll,
		GraphQLCors:           source.GraphQLCors,
		GraphQLVirtualHosts:   source.GraphQLVirtualHosts,
		Logger:                source.Logger,
		AllowUnprotectedTxs:   source.AllowUnprotectedTxs,
	}
}
