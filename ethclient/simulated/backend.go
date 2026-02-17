// Copyright 2023 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package simulated

import (
	"errors"
	"fmt"
	"time"

	"github.com/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/eth"
	"github.com/autonity/autonity/eth/ethconfig"
	"github.com/autonity/autonity/eth/filters"
	"github.com/autonity/autonity/ethclient"
	"github.com/autonity/autonity/node"
	"github.com/autonity/autonity/p2p"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/rpc"
)

// Client exposes the methods provided by the Ethereum RPC client.
type Client interface {
	ethereum.BlockNumberReader
	ethereum.ChainReader
	ethereum.ChainStateReader
	ethereum.ContractCaller
	ethereum.GasEstimator
	ethereum.GasPricer
	ethereum.GasPricer1559
	ethereum.FeeHistoryReader
	ethereum.LogFilterer
	ethereum.PendingStateReader
	ethereum.PendingContractCaller
	ethereum.TransactionReader
	ethereum.TransactionSender
	ethereum.ChainIDReader
}

// simClient wraps ethclient. This exists to prevent extracting ethclient.Client
// from the Client interface returned by Backend.
type simClient struct {
	*ethclient.Client
}

// Backend is a simulated blockchain. You can use it to test your contracts or
// other code that interacts with the Ethereum chain.
type Backend struct {
	node      *node.Node
	client    simClient
	eth       *eth.Ethereum
	simulator *tendermintSimulator
}

// NewBackend creates a new simulated blockchain that can be used as a backend for
// contract bindings in unit tests.
//
// A simulated backend always uses chainID 1337.
func NewBackend(alloc types.GenesisAlloc, options ...func(nodeConf *node.Config, ethConf *ethconfig.Config)) *Backend {
	// Create simulator validator with ECDSA and BLS keys
	validator, err := newSimulatorValidator()
	if err != nil {
		panic(err) // this should never happen
	}

	// Create the default configurations for the outer node shell and the Ethereum
	// service to mutate with the options afterwards
	nodeConf := node.DefaultConfig
	nodeConf.DataDir = ""
	nodeConf.ExecutionP2P = p2p.Config{NoDiscovery: true}

	ethConf := ethconfig.Defaults
	ethConf.Genesis = &core.Genesis{
		Config:    params.TestChainConfig,
		Mixhash:   types.BFTDigest,
		GasLimit:  ethconfig.Defaults.Miner.GasCeil,
		Alloc:     alloc,
		Timestamp: 1,
	}
	ethConf.SyncMode = ethconfig.FullSync
	ethConf.TxPool.NoLocals = true

	// Inject validator into genesis
	appendSimulatorValidator(ethConf.Genesis, validator)

	// Prepare the genesis config
	if err := ethConf.Genesis.Config.Prepare(); err != nil {
		panic(err) // this should never happen
	}

	for _, option := range options {
		option(&nodeConf, &ethConf)
	}
	// Assemble the Ethereum stack to run the chain with
	stack, err := node.New(&nodeConf)
	if err != nil {
		panic(err) // this should never happen
	}
	sim, err := newWithNode(stack, &ethConf, 0, validator)
	if err != nil {
		panic(err) // this should never happen
	}
	return sim
}

// newWithNode sets up a simulated backend on an existing node. The provided node
// must not be started and will be started by this method.
func newWithNode(stack *node.Node, conf *eth.Config, blockPeriod uint64, validator *simulatorValidator) (*Backend, error) {
	ethBackend, err := eth.New(stack, conf)
	if err != nil {
		return nil, err
	}
	// Register the filter system
	filterSystem := filters.NewFilterSystem(ethBackend.APIBackend, filters.Config{})
	stack.RegisterAPIs([]rpc.API{{
		Namespace: "eth",
		Service:   filters.NewFilterAPI(filterSystem),
	}})
	// Start the node
	if err := stack.Start(); err != nil {
		return nil, err
	}
	n := stack.Attach()

	// Create the Tendermint simulator
	simulator := newTendermintSimulator(validator, ethBackend)

	return &Backend{
		node:      stack,
		client:    simClient{ethclient.NewClient(n)},
		eth:       ethBackend,
		simulator: simulator,
	}, nil
}

// Close shuts down the simBackend.
// The simulated backend can't be used afterwards.
func (n *Backend) Close() error {
	if n.client.Client != nil {
		n.client.Close()
		n.client = simClient{}
	}
	var err error
	if n.node != nil {
		err = errors.Join(err, n.node.Close())
		n.node = nil
	}
	return err
}

// Commit seals a block and moves the chain forward to a new empty block.
func (n *Backend) Commit() common.Hash {
	block, err := n.simulator.createBlock()
	if err != nil {
		panic(fmt.Sprintf("failed to create block: %v", err))
	}

	chain := n.eth.BlockChain()

	// Insert the block using InsertChain
	// Note: InsertChain may return 0 inserted if the block triggers a sidechain
	// or is being processed asynchronously
	_, err = chain.InsertChain(types.Blocks{block})
	if err != nil {
		panic(fmt.Sprintf("failed to insert block %d: %v", block.NumberU64(), err))
	}

	// Check if block is in the chain
	if chain.GetBlockByHash(block.Hash()) == nil {
		panic(fmt.Sprintf("block %d not found in chain after InsertChain", block.NumberU64()))
	}

	// Set this block as the canonical head if not already
	currentHead := chain.CurrentBlock()
	if currentHead.Hash() != block.Hash() {
		_, err = chain.SetCanonical(block)
		if err != nil {
			panic(fmt.Sprintf("failed to set canonical head to block %d: %v", block.NumberU64(), err))
		}
	}

	return block.Hash()
}

// Rollback removes all pending transactions, reverting to the last committed state.
func (n *Backend) Rollback() {
	// Clear any fork/time state in the simulator
	n.simulator.rollback()
}

// Fork creates a side-chain that can be used to simulate reorgs.
//
// This function should be called with the ancestor block where the new side
// chain should be started. Transactions (old and new) can then be applied on
// top and Commit-ed.
//
// Note, the side-chain will only become canonical (and trigger the events) when
// it becomes longer. Until then CallContract will still operate on the current
// canonical chain.
//
// There is a % chance that the side chain becomes canonical at the same length
// to simulate live network behavior.
func (n *Backend) Fork(parentHash common.Hash) error {
	return n.simulator.setForkParent(parentHash)
}

// AdjustTime changes the block timestamp and creates a new block.
// It can only be called on empty blocks.
func (n *Backend) AdjustTime(adjustment time.Duration) error {
	n.simulator.adjustTime(adjustment)
	n.Commit()
	return nil
}

// Client returns a client that accesses the simulated chain.
func (n *Backend) Client() Client {
	return n.client
}

// EthClient returns the underlying ethclient.Client for direct access to all RPC methods.
// Use this when you need methods not exposed by the Client interface.
func (n *Backend) EthClient() *ethclient.Client {
	return n.client.Client
}

// Node returns the underlying node for direct access to RPC endpoints.
func (n *Backend) Node() *node.Node {
	return n.node
}
