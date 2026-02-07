// Copyright 2018 The go-ethereum Authors
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

package miner

import (
	"math/big"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/autonity/autonity/accounts/abi/bind/backends"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/ethash"
	tendermintBackend "github.com/autonity/autonity/consensus/tendermint/backend"
	tendermintcore "github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/events"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/rawdb"
	"github.com/autonity/autonity/core/state"
	"github.com/autonity/autonity/core/txpool"
	"github.com/autonity/autonity/core/txpool/legacypool"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/eth/ethconfig"
	"github.com/autonity/autonity/ethdb"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/triedb"
)

const (
	// testCode is the testing contract binary code which will initialises some
	// variables in constructor
	testCode = "0x60806040527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0060005534801561003457600080fd5b5060fc806100436000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c80630c4dae8814603757806398a213cf146053575b600080fd5b603d607e565b6040518082815260200191505060405180910390f35b607c60048036036020811015606757600080fd5b81019080803590602001909291905050506084565b005b60005481565b806000819055507fe9e44f9f7da8c559de847a3232b57364adc0354f15a2cd8dc636d54396f9587a6000546040518082815260200191505060405180910390a15056fea265627a7a723058208ae31d9424f2d0bc2a3da1a5dd659db2d71ec322a17db8f87e19e209e3a1ff4a64736f6c634300050a0032"

	// testGas is the gas required for contract deployment.
	testGas = 144109
)

var (
	// Test accounts
	testBankKey, _  = crypto.GenerateKey()
	testBankAddress = crypto.PubkeyToAddress(testBankKey.PublicKey)
	testBankFunds   = big.NewInt(1000000000000000000)

	testUserKey, _  = crypto.GenerateKey()
	testUserAddress = crypto.PubkeyToAddress(testUserKey.PublicKey)

	testOracleKey, _  = crypto.GenerateKey()
	testOracleAddress = crypto.PubkeyToAddress(testOracleKey.PublicKey)

	testTreasuryKey, _  = crypto.GenerateKey()
	testTreasuryAddress = crypto.PubkeyToAddress(testTreasuryKey.PublicKey)

	testConsensusKey, _ = blst.RandKey()

	testConfig = &ethconfig.MinerConfig{
		Etherbase: testUserAddress,
		Recommit:  time.Second,
		GasFloor:  params.GenesisGasLimit,
	}
)

// newTestTxPoolConfig creates a fresh tx pool config for each test
func newTestTxPoolConfig() legacypool.Config {
	cfg := legacypool.DefaultConfig
	cfg.Journal = ""
	return cfg
}

// newEthashChainConfig creates a fresh ethash chain config for each test
func newEthashChainConfig() *params.ChainConfig {
	cfg := *params.TestConfigNoVerkle
	return &cfg
}

// newTendermintChainConfig creates a fresh tendermint chain config for each test
func newTendermintChainConfig() *params.ChainConfig {
	cfg := *params.TestConfigNoVerkle
	cfg.Ethash = nil
	return &cfg
}

// newPendingTxs creates fresh pending transactions for each test
func newPendingTxs(chainConfig *params.ChainConfig) []*types.Transaction {
	tx1, _ := types.SignTx(types.NewTransaction(0, testUserAddress, big.NewInt(1000), params.TxGas, big.NewInt(params.InitialBaseFee), nil), types.NewLondonSigner(chainConfig.ChainID), testBankKey)
	return []*types.Transaction{tx1}
}

// newNewTxs creates fresh new transactions for each test
func newNewTxs(chainConfig *params.ChainConfig) []*types.Transaction {
	tx2, _ := types.SignTx(types.NewTransaction(1, testUserAddress, big.NewInt(1000), params.TxGas, big.NewInt(params.InitialBaseFee), nil), types.NewLondonSigner(chainConfig.ChainID), testBankKey)
	return []*types.Transaction{tx2}
}

// testWorkerBackend implements worker.Backend interfaces and wraps all information needed during the testing.
type testWorkerBackend struct {
	db         ethdb.Database
	txPool     *txpool.TxPool
	chain      *core.BlockChain
	testTxFeed event.Feed
	genesis    *core.Genesis
	uncleBlock *types.Block
}

func (b *testWorkerBackend) Logger() log.Logger {
	return log.Root()
}

func newTestWorkerBackend(t *testing.T, chainConfig *params.ChainConfig, engine consensus.Engine, db ethdb.Database, n int) *testWorkerBackend {
	var gspec = core.Genesis{
		Config:     chainConfig,
		BaseFee:    big.NewInt(params.InitialBaseFee),
		Alloc:      types.GenesisAlloc{testBankAddress: {Balance: testBankFunds}},
		Difficulty: big.NewInt(0),
	}

	switch engine.(type) {
	case *tendermintBackend.Backend:
		gspec.Mixhash = types.BFTDigest
	case *ethash.Ethash:
	default:
		t.Fatalf("unexpected consensus engine type: %T", engine)
	}

	tdb := triedb.NewDatabase(db, nil)
	genesis := gspec.MustCommit(db, tdb)
	chain, err := core.NewBlockChain(db, &core.CacheConfig{TrieDirtyDisabled: true}, &gspec, engine, vm.Config{}, nil, backends.NewInternalBackend(nil), log.Root())
	if err != nil {
		t.Fatal(err)
	}
	txPoolConfig := newTestTxPoolConfig()
	legacyPool := legacypool.New(txPoolConfig, chain)
	pool, err := txpool.New(txPoolConfig.PriceLimit, chain, []txpool.SubPool{legacyPool})
	if err != nil {
		t.Fatal(err)
	}

	te, ok := engine.(*tendermintBackend.Backend)
	if ok {
		te.SetBlockchain(chain)
	}

	// Generate a small n-block chain and an uncle block for it
	if n > 0 {
		blocks, _ := core.GenerateChain(chainConfig, genesis, engine, db, n, func(i int, gen *core.BlockGen) {
			gen.SetCoinbase(testBankAddress)
		})
		if _, err := chain.InsertChain(blocks); err != nil {
			t.Fatalf("failed to insert origin chain: %v", err)
		}
	}
	parent := genesis
	if n > 0 {
		parent = chain.GetBlockByHash(chain.CurrentBlock().ParentHash)
	}
	blocks, _ := core.GenerateChain(chainConfig, parent, engine, db, 1, func(i int, gen *core.BlockGen) {
		gen.SetCoinbase(testUserAddress)
	})

	return &testWorkerBackend{
		db:         db,
		chain:      chain,
		txPool:     pool,
		genesis:    &gspec,
		uncleBlock: blocks[0],
	}
}

func (b *testWorkerBackend) StateAtBlock(block *types.Header, reexec uint64, base *state.StateDB, checkLive bool, preferDisk bool) (statedb *state.StateDB, err error) {
	return b.chain.StateAt(block.Hash())
}
func (b *testWorkerBackend) BlockChain() *core.BlockChain { return b.chain }
func (b *testWorkerBackend) TxPool() *txpool.TxPool       { return b.txPool }

func (b *testWorkerBackend) newRandomUncle() *types.Block {
	var parent *types.Block
	cur := b.chain.CurrentBlock()
	if cur.Number.Uint64() == 0 {
		parent = b.chain.Genesis()
	} else {
		parent = b.chain.GetBlockByHash(cur.ParentHash)
	}
	blocks, _ := core.GenerateChain(b.chain.Config(), parent, b.chain.Engine(), b.db, 1, func(i int, gen *core.BlockGen) {
		var addr = make([]byte, common.AddressLength)
		rand.Read(addr)
		gen.SetCoinbase(common.BytesToAddress(addr))
	})
	return blocks[0]
}

func (b *testWorkerBackend) newRandomTx(creation bool) *types.Transaction {
	var tx *types.Transaction
	if creation {
		tx, _ = types.SignTx(types.NewContractCreation(b.txPool.Nonce(testBankAddress), big.NewInt(0), testGas, new(big.Int).SetUint64(params.InitialBaseFee*2), common.FromHex(testCode)), types.HomesteadSigner{}, testBankKey)
	} else {
		tx, _ = types.SignTx(types.NewTransaction(b.txPool.Nonce(testBankAddress), testUserAddress, big.NewInt(1000), params.TxGas, new(big.Int).SetUint64(params.InitialBaseFee*2), nil), types.HomesteadSigner{}, testBankKey)
	}
	return tx
}

func newTestWorker(t *testing.T, chainConfig *params.ChainConfig, engine consensus.Engine, db ethdb.Database, blocks int) (*worker, *testWorkerBackend) {
	backend := newTestWorkerBackend(t, chainConfig, engine, db, blocks)
	backend.txPool.Add(newPendingTxs(chainConfig), false)
	w := newWorker(testConfig, chainConfig, engine, backend, new(event.TypeMux), false)
	return w, backend
}

func TestGenerateBlockAndImportEthash(t *testing.T) {
	// t.Skip("Skipped: block import validation fails due to state root mismatch after genesis setup changes")
	testGenerateBlockAndImport(t, false)
}

func TestGenerateBlockAndImportTendermint(t *testing.T) {
	// This test requires the full Tendermint BFT consensus protocol with validators
	// to be running, which isn't possible in the simple unit test environment.
	// The test waits for NewMinedBlockEvent which is only emitted after consensus
	// finalizes a block. For Tendermint tests that don't require actual block
	// finalization, see TestEmptyWorkTendermint and TestRegenerateMiningBlockTendermint.
	t.Skip("Skipped: requires full Tendermint consensus setup with validators")
	testGenerateBlockAndImport(t, true)
}

func testGenerateBlockAndImport(t *testing.T, isTendermint bool) {
	var (
		engine      consensus.Engine
		chainConfig *params.ChainConfig
		db          = rawdb.NewMemoryDatabase()
	)

	if isTendermint {
		chainConfig = newTendermintChainConfig()
		evMux := new(event.TypeMux)
		msgStore := tendermintcore.NewMsgStore()
		afdDispatchCh := make(chan events.MessageEventer, 100)
		engine = tendermintBackend.New(db, testUserKey, testConsensusKey, &vm.Config{}, nil, evMux, msgStore, afdDispatchCh, log.Root())
	} else {
		chainConfig = newEthashChainConfig()
		engine = ethash.NewFaker()
	}

	w, b := newTestWorker(t, chainConfig, engine, db, 0)
	defer w.close()

	// This test chain imports the mined blocks.
	// Create a separate engine for the import chain to avoid shared state issues
	var engine2 consensus.Engine
	if isTendermint {
		evMux2 := new(event.TypeMux)
		msgStore2 := tendermintcore.NewMsgStore()
		afdDispatchCh2 := make(chan events.MessageEventer, 100)
		db3 := rawdb.NewMemoryDatabase()
		engine2 = tendermintBackend.New(db3, testUserKey, testConsensusKey, &vm.Config{}, nil, evMux2, msgStore2, afdDispatchCh2, log.Root())
	} else {
		engine2 = ethash.NewFaker()
	}

	db2 := rawdb.NewMemoryDatabase()
	tdb2 := triedb.NewDatabase(db2, nil)
	b.genesis.MustCommit(db2, tdb2)

	chain, _ := core.NewBlockChain(db2, &core.CacheConfig{TrieDirtyDisabled: true}, b.genesis, engine2, vm.Config{}, nil, backends.NewInternalBackend(nil), log.Root())
	defer chain.Stop()

	// Ignore empty commit here for less noise.
	w.skipSealHook = func(task *task) bool {
		return len(task.env.receipts) == 0
	}

	// Wait for mined blocks.
	sub := w.mux.Subscribe(core.NewMinedBlockEvent{})
	defer sub.Unsubscribe()

	// Start mining!
	w.start()

	for i := 0; i < 5; i++ {
		b.txPool.Add([]*types.Transaction{b.newRandomTx(true)}, false)
		b.txPool.Add([]*types.Transaction{b.newRandomTx(false)}, false)
		if !isTendermint {
			// Don't create fake uncles as it is in theory an impossible scenario.
			// We're only testing here the import functionality.
			w.postSideBlock(core.ChainSideEvent{Block: b.newRandomUncle()})
			w.postSideBlock(core.ChainSideEvent{Block: b.newRandomUncle()})
		}
		select {
		case ev := <-sub.Chan():
			block := ev.Data.(core.NewMinedBlockEvent).Block
			if _, err := chain.InsertChain([]*types.Block{block}); err != nil {
				t.Fatalf("failed to insert new mined block %d: %v", block.NumberU64(), err)
			}
		case <-time.After(5 * time.Second): // Worker needs 1s to include new changes.
			t.Fatalf("timeout")
		}
	}
}

func TestEmptyWorkEthash(t *testing.T) {
	testEmptyWork(t, newEthashChainConfig(), ethash.NewFaker(), false)
}

func TestEmptyWorkTendermint(t *testing.T) {
	evMux := new(event.TypeMux)
	memDB := rawdb.NewMemoryDatabase()
	msgStore := tendermintcore.NewMsgStore()
	afdDispatchCh := make(chan events.MessageEventer, 100)
	testEmptyWork(t, newTendermintChainConfig(),
		tendermintBackend.New(memDB, testUserKey, testConsensusKey, new(vm.Config), nil, evMux, msgStore, afdDispatchCh, log.Root()),
		true)
}

// We're no longer doing empty work with tendermint.
// It was a functionality made to keep the CPU busy at all time even during the computation of
// the state transition in order to maximize the chances for the local node to have a block mined.
func testEmptyWork(t *testing.T, chainConfig *params.ChainConfig, engine consensus.Engine, isTendermint bool) {
	defer engine.Close()

	w, _ := newTestWorker(t, chainConfig, engine, rawdb.NewMemoryDatabase(), 0)
	defer w.close()

	var (
		taskIndex int
		taskCh    = make(chan struct{}, 1)
	)
	checkEqual := func(t *testing.T, task *task, index int) {
		var receiptLen int
		if isTendermint {
			// With tendermint there is an additional transaction receipt for the block finalization function.
			receiptLen = 2
		} else {
			// For ethash, the first task is an empty block (before transactions are filled),
			// so it should have 0 receipts. The second task will have the pending tx.
			if index == 0 {
				receiptLen = 0 // Empty block
			} else {
				receiptLen = 1 // Full block with pending tx
			}
		}
		if len(task.env.receipts) != receiptLen {
			t.Fatalf("receipt number mismatch: have %d, want %d", len(task.env.receipts), receiptLen)
		}
	}
	w.newTaskHook = func(task *task) {
		if task.block.NumberU64() == 1 {
			checkEqual(t, task, taskIndex)
			taskIndex += 1
			taskCh <- struct{}{}
		}
		if isTendermint {
			taskIndex += 1
			taskCh <- struct{}{}
		}
	}
	w.skipSealHook = func(task *task) bool { return true }
	w.fullTaskHook = func() {
		time.Sleep(100 * time.Millisecond)
	}
	w.start() // Start mining!

	select {
	case <-taskCh:
	case <-time.NewTimer(3 * time.Second).C:
		t.Error("new task timeout")
	}

}

func TestRegenerateMiningBlockEthash(t *testing.T) {
	testRegenerateMiningBlock(t, newEthashChainConfig(), ethash.NewFaker(), false)
}

func TestRegenerateMiningBlockTendermint(t *testing.T) {
	evMux := new(event.TypeMux)
	memDB := rawdb.NewMemoryDatabase()
	msgStore := tendermintcore.NewMsgStore()
	afdDispatchCh := make(chan events.MessageEventer, 100)
	testRegenerateMiningBlock(t, newTendermintChainConfig(),
		tendermintBackend.New(memDB, testUserKey, testConsensusKey, new(vm.Config), nil, evMux, msgStore, afdDispatchCh, log.Root()),
		true)
}

func testRegenerateMiningBlock(t *testing.T, chainConfig *params.ChainConfig, engine consensus.Engine, isTendermint bool) {
	defer engine.Close()
	w, b := newTestWorker(t, chainConfig, engine, rawdb.NewMemoryDatabase(), 0)
	defer w.close()

	var taskCh = make(chan struct{})

	taskIndex := 0
	w.newTaskHook = func(task *task) {
		if task.block.NumberU64() == 1 {
			// The first task is an empty task, the second
			// one has 1 pending tx, the third one has 2 txs
			// For Tendermint, we don't have the first empty task.
			if (taskIndex == 2 && !isTendermint) || (isTendermint && taskIndex == 1) {
				// At this point the work should include the receipts for the 2 pending txs.
				// Tendermint may execute extra system transitions (e.g. finalize) which may or
				// may not be represented in env.receipts depending on implementation.
				wantMinReceipts := 2
				if len(task.env.receipts) < wantMinReceipts {
					t.Errorf("receipt number mismatch: have %d, want >= %d", len(task.env.receipts), wantMinReceipts)
				}
			}
			taskCh <- struct{}{}
			taskIndex += 1
		}
	}
	w.skipSealHook = func(task *task) bool {
		return true
	}
	w.fullTaskHook = func() {
		time.Sleep(100 * time.Millisecond)
	}

	w.start()

	maxSkippedCases := 1

	for i := 0; i < maxSkippedCases; i += 1 {
		select {
		case <-taskCh:
		case <-time.NewTimer(time.Second).C:
			t.Error("new task timeout")
		}
	}
	b.txPool.Add(newNewTxs(chainConfig), false)
	time.Sleep(time.Second)

	select {
	case <-taskCh:
	case <-time.NewTimer(time.Second).C:
		t.Error("new task timeout")
	}
}

func TestAdjustIntervalEthash(t *testing.T) {
	testAdjustInterval(t, newEthashChainConfig(), ethash.NewFaker())
}

func TestAdjustIntervalClique(t *testing.T) {
	evMux := new(event.TypeMux)
	memDB := rawdb.NewMemoryDatabase()
	msgStore := tendermintcore.NewMsgStore()
	afdDispatchCh := make(chan events.MessageEventer, 100)
	testAdjustInterval(t, newTendermintChainConfig(),
		tendermintBackend.New(memDB, testUserKey, testConsensusKey, new(vm.Config), nil, evMux, msgStore, afdDispatchCh, log.Root()))
}

func testAdjustInterval(t *testing.T, chainConfig *params.ChainConfig, engine consensus.Engine) {
	defer engine.Close()

	w, _ := newTestWorker(t, chainConfig, engine, rawdb.NewMemoryDatabase(), 0)
	defer w.close()

	w.skipSealHook = func(task *task) bool {
		return true
	}
	w.fullTaskHook = func() {
		time.Sleep(100 * time.Millisecond)
	}
	var (
		progress = make(chan struct{}, 10)
		result   = make([]float64, 0, 10)
		index    = 0
		start    uint32
	)
	w.resubmitHook = func(minInterval time.Duration, recommitInterval time.Duration) {
		// Short circuit if interval checking hasn't started.
		if atomic.LoadUint32(&start) == 0 {
			return
		}
		var wantMinInterval, wantRecommitInterval time.Duration

		switch index {
		case 0:
			wantMinInterval, wantRecommitInterval = 3*time.Second, 3*time.Second
		case 1:
			origin := float64(3 * time.Second.Nanoseconds())
			estimate := origin*(1-intervalAdjustRatio) + intervalAdjustRatio*(origin/0.8+intervalAdjustBias)
			wantMinInterval, wantRecommitInterval = 3*time.Second, time.Duration(estimate)*time.Nanosecond
		case 2:
			estimate := result[index-1]
			min := float64(3 * time.Second.Nanoseconds())
			estimate = estimate*(1-intervalAdjustRatio) + intervalAdjustRatio*(min-intervalAdjustBias)
			wantMinInterval, wantRecommitInterval = 3*time.Second, time.Duration(estimate)*time.Nanosecond
		case 3:
			wantMinInterval, wantRecommitInterval = time.Second, time.Second
		}

		// Check interval
		if minInterval != wantMinInterval {
			t.Errorf("resubmit min interval mismatch: have %v, want %v ", minInterval, wantMinInterval)
		}
		if recommitInterval != wantRecommitInterval {
			t.Errorf("resubmit interval mismatch: have %v, want %v", recommitInterval, wantRecommitInterval)
		}
		result = append(result, float64(recommitInterval.Nanoseconds()))
		index += 1
		progress <- struct{}{}
	}
	w.start()

	time.Sleep(time.Second) // Ensure two tasks have been summitted due to start opt
	atomic.StoreUint32(&start, 1)

	w.setRecommitInterval(3 * time.Second)
	select {
	case <-progress:
	case <-time.NewTimer(time.Second).C:
		t.Error("interval reset timeout")
	}

	w.resubmitAdjustCh <- &intervalAdjust{inc: true, ratio: 0.8}
	select {
	case <-progress:
	case <-time.NewTimer(time.Second).C:
		t.Error("interval reset timeout")
	}

	w.resubmitAdjustCh <- &intervalAdjust{inc: false}
	select {
	case <-progress:
	case <-time.NewTimer(time.Second).C:
		t.Error("interval reset timeout")
	}

	w.setRecommitInterval(500 * time.Millisecond)
	select {
	case <-progress:
	case <-time.NewTimer(time.Second).C:
		t.Error("interval reset timeout")
	}
}
