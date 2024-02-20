// Copyright 2020 The go-ethereum Authors
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

// Package miner implements Ethereum block creation and mining.
package miner

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/accounts/abi/bind/backends"
	"github.com/autonity/autonity/log"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/ethash"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/rawdb"
	"github.com/autonity/autonity/core/state"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/eth/downloader"
	"github.com/autonity/autonity/ethdb/memorydb"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/trie"
)

type mockBackend struct {
	bc     *core.BlockChain
	txPool *core.TxPool
}

func (m *mockBackend) Logger() log.Logger {
	return log.Root()
}

func NewMockBackend(bc *core.BlockChain, txPool *core.TxPool) *mockBackend {
	return &mockBackend{
		bc:     bc,
		txPool: txPool,
	}
}
func (m *mockBackend) BlockChain() *core.BlockChain {
	return m.bc
}

func (m *mockBackend) StateAtBlock(block *types.Block, reexec uint64, base *state.StateDB, checkLive bool, preferDisk bool) (statedb *state.StateDB, err error) {
	return m.bc.StateAt(block.Hash())
}

func (m *mockBackend) TxPool() *core.TxPool {
	return m.txPool
}

type testBlockChain struct {
	statedb       *state.StateDB
	gasLimit      uint64
	chainHeadFeed *event.Feed
}

func (bc *testBlockChain) MinBaseFee() *big.Int {
	return new(big.Int)
}

func (bc *testBlockChain) Config() *params.ChainConfig {
	return nil
}

func (bc *testBlockChain) CurrentBlock() *types.Block {
	return types.NewBlock(&types.Header{
		GasLimit: bc.gasLimit,
	}, nil, nil, nil, new(trie.Trie))
}

func (bc *testBlockChain) GetBlock(hash common.Hash, number uint64) *types.Block {
	return bc.CurrentBlock()
}

func (bc *testBlockChain) StateAt(common.Hash) (*state.StateDB, error) {
	return bc.statedb, nil
}

func (bc *testBlockChain) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	return bc.chainHeadFeed.Subscribe(ch)
}

// consensus should be started only after `DoneEvent` or `SyncedEvent` are broadcasted
// then it can continue running regardless of what happens in downloader
func TestMiner(t *testing.T) {
	// create miner and start it (as if we are in the committee)
	miner, mux := createMiner(t)
	miner.Start()
	waitForMiningState(t, miner, false)
	require.True(t, miner.shouldStart)
	require.False(t, miner.canStart)

	// Start the downloader (simulate syncing from another peer)
	mux.Post(downloader.StartEvent{})
	waitForMiningState(t, miner, false)

	// Stop the downloader and wait for the update loop to run
	// syncing has failed, so we do not start mining
	mux.Post(downloader.FailedEvent{})
	waitForMiningState(t, miner, false)

	// start syncing again
	mux.Post(downloader.StartEvent{})
	waitForMiningState(t, miner, false)
	require.True(t, miner.shouldStart)
	require.False(t, miner.canStart)

	// Successfully stop the downloader, mining should start
	mux.Post(downloader.DoneEvent{})
	waitForMiningState(t, miner, true)
	require.True(t, miner.shouldStart)
	require.True(t, miner.canStart)

	// Subsequent downloader events after mining is started should not cause the
	// miner to start or stop. This prevents a security vulnerability
	// that would allow entities to present fake high blocks that would
	// stop mining operations by causing a downloader sync
	// until it was discovered they were invalid, whereon mining would resume.
	mux.Post(downloader.StartEvent{})
	waitForMiningState(t, miner, true)
	mux.Post(downloader.FailedEvent{})
	waitForMiningState(t, miner, true)
	mux.Post(downloader.SyncedEvent{})
	waitForMiningState(t, miner, true)
	mux.Post(downloader.DoneEvent{})
	waitForMiningState(t, miner, true)
}

func TestMinerStartStopAfterDownloaderEvents(t *testing.T) {
	miner, mux := createMiner(t)
	miner.Start()
	waitForMiningState(t, miner, false)

	// Start the downloader
	mux.Post(downloader.StartEvent{})
	waitForMiningState(t, miner, false)

	// Downloader finally succeeds.
	mux.Post(downloader.DoneEvent{})
	waitForMiningState(t, miner, true)

	miner.Stop()
	waitForMiningState(t, miner, false)
	require.False(t, miner.shouldStart)
	require.True(t, miner.canStart)

	miner.Start()
	waitForMiningState(t, miner, true)

	miner.Stop()
	waitForMiningState(t, miner, false)
}

func TestStartWhileDownload(t *testing.T) {
	miner, mux := createMiner(t)
	waitForMiningState(t, miner, false)
	miner.Start()
	waitForMiningState(t, miner, false)

	// Start the downloader and wait for the update loop to run
	mux.Post(downloader.StartEvent{})
	waitForMiningState(t, miner, false)

	// Starting the miner while syncing should not work
	miner.Start()
	waitForMiningState(t, miner, false)
	require.True(t, miner.shouldStart)
	require.False(t, miner.canStart)
}

func TestStartStopMiner(t *testing.T) {
	miner, mux := createMiner(t)
	waitForMiningState(t, miner, false)
	miner.Start()
	waitForMiningState(t, miner, false)

	mux.Post(downloader.SyncedEvent{})
	waitForMiningState(t, miner, true)

	miner.Stop()
	waitForMiningState(t, miner, false)
	require.False(t, miner.shouldStart)
	require.True(t, miner.canStart)
}

func TestCloseMiner(t *testing.T) {
	miner, mux := createMiner(t)
	waitForMiningState(t, miner, false)
	miner.Start()
	waitForMiningState(t, miner, false)

	mux.Post(downloader.SyncedEvent{})
	waitForMiningState(t, miner, true)

	// Terminate the miner and wait for the update loop to run
	miner.Close()
	waitForMiningState(t, miner, false)
}

// a node syncs up, then enters the committee
func TestEnterExitCommittee(t *testing.T) {
	t.Run("Enter due to DoneEvent", func(t *testing.T) {
		miner, mux := createMiner(t)
		waitForMiningState(t, miner, false)

		// first sync attempt fails, 2nd succeeds
		mux.Post(downloader.StartEvent{})
		waitForMiningState(t, miner, false)
		mux.Post(downloader.FailedEvent{})
		waitForMiningState(t, miner, false)
		mux.Post(downloader.StartEvent{})
		waitForMiningState(t, miner, false)
		mux.Post(downloader.DoneEvent{})
		waitForMiningState(t, miner, false)
		require.False(t, miner.shouldStart)
		require.True(t, miner.canStart)

		// enter committee
		miner.Start()
		waitForMiningState(t, miner, true)
		require.True(t, miner.shouldStart)
		require.True(t, miner.canStart)

		// misc events, should be ignored
		mux.Post(downloader.StartEvent{})
		waitForMiningState(t, miner, true)
		mux.Post(downloader.FailedEvent{})
		waitForMiningState(t, miner, true)
		mux.Post(downloader.SyncedEvent{})
		waitForMiningState(t, miner, true)
		mux.Post(downloader.DoneEvent{})
		waitForMiningState(t, miner, true)

		// exit committee
		miner.Stop()
		waitForMiningState(t, miner, false)
		require.False(t, miner.shouldStart)
		require.True(t, miner.canStart)

		// misc events, should be ignored
		mux.Post(downloader.StartEvent{})
		waitForMiningState(t, miner, false)
		mux.Post(downloader.FailedEvent{})
		waitForMiningState(t, miner, false)
		mux.Post(downloader.SyncedEvent{})
		waitForMiningState(t, miner, false)
		mux.Post(downloader.DoneEvent{})
		waitForMiningState(t, miner, false)

		// enter committee again
		miner.Start()
		waitForMiningState(t, miner, true)
	})
	t.Run("Enter due to SyncedEvent", func(t *testing.T) {
		miner, mux := createMiner(t)
		waitForMiningState(t, miner, false)
		require.False(t, miner.shouldStart)
		require.False(t, miner.canStart)

		mux.Post(downloader.SyncedEvent{})
		waitForMiningState(t, miner, false)
		require.False(t, miner.shouldStart)
		require.True(t, miner.canStart)

		// enter committee
		miner.Start()
		waitForMiningState(t, miner, true)
		require.True(t, miner.shouldStart)
		require.True(t, miner.canStart)
	})
}

// waitForMiningState waits until either
// * the desired mining state was reached
// * a timeout was reached which fails the test
func waitForMiningState(t *testing.T, m *Miner, mining bool) {
	t.Helper()

	var state bool
	for i := 0; i < 100; i++ {
		time.Sleep(10 * time.Millisecond)
		if state = m.Mining(); state == mining {
			return
		}
	}
	t.Fatalf("Mining() == %t, want %t", state, mining)
}

func createMiner(t *testing.T) (*Miner, *event.TypeMux) {
	// Create Ethash config
	config := Config{
		Etherbase: common.HexToAddress("123456789"),
	}
	// Create chainConfig
	memdb := memorydb.New()
	chainDB := rawdb.NewDatabase(memdb)
	genesis := core.DefaultGenesisBlock()
	chainConfig, _, err := core.SetupGenesisBlock(chainDB, genesis)
	if err != nil {
		t.Fatalf("can't create new chain config: %v", err)
	}
	// Create event Mux
	mux := new(event.TypeMux)
	// Create consensus engine
	engine := ethash.New(ethash.Config{}, []string{}, false)
	engine.SetThreads(-1)
	// Create isLocalBlock
	isLocalBlock := func(block *types.Header) bool {
		return true
	}
	// Create Ethereum backend
	limit := uint64(1000)
	senderCacher := new(core.TxSenderCacher)
	bc, err := core.NewBlockChain(chainDB, new(core.CacheConfig), chainConfig, engine, vm.Config{}, isLocalBlock, senderCacher, &limit, backends.NewInternalBackend(nil), log.Root())
	if err != nil {
		t.Fatalf("can't create new chain %v", err)
	}
	statedb, _ := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
	blockchain := &testBlockChain{statedb, 10000000, new(event.Feed)}

	pool := core.NewTxPool(testTxPoolConfig, params.TestChainConfig, blockchain, senderCacher)
	backend := NewMockBackend(bc, pool)
	// Create Miner
	return New(backend, &config, chainConfig, mux, engine, isLocalBlock), mux
}
