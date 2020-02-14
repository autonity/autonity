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
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/ethash"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/rawdb"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/core/vm"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/ethdb"
	"github.com/clearmatics/autonity/event"
	"github.com/clearmatics/autonity/params"
)

const (
	// testCode is the testing contract binary code which will initialises some
	// variables in constructor
	testCode = "0x60806040527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0060005534801561003457600080fd5b5060fc806100436000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c80630c4dae8814603757806398a213cf146053575b600080fd5b603d607e565b6040518082815260200191505060405180910390f35b607c60048036036020811015606757600080fd5b81019080803590602001909291905050506084565b005b60005481565b806000819055507fe9e44f9f7da8c559de847a3232b57364adc0354f15a2cd8dc636d54396f9587a6000546040518082815260200191505060405180910390a15056fea265627a7a723058208ae31d9424f2d0bc2a3da1a5dd659db2d71ec322a17db8f87e19e209e3a1ff4a64736f6c634300050a0032"

	// testGas is the gas required for contract deployment.
	testGas = 144109
)

type testCase struct {
	testTxPoolConfig  core.TxPoolConfig
	ethashChainConfig *params.ChainConfig

	testBankKey     *ecdsa.PrivateKey
	testBankAddress common.Address
	testBankFunds   *big.Int

	testUserKey     *ecdsa.PrivateKey
	testUserAddress common.Address
	testUserFunds   *big.Int

	// Test transactions
	pendingTxs []*types.Transaction
	newTxs     []*types.Transaction

	testConfig *Config
}

func getTestCase() (*testCase, error) {
	testBankKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	testUserKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	t := &testCase{
		testBankKey:     testBankKey,
		testBankAddress: crypto.PubkeyToAddress(testBankKey.PublicKey),
		testBankFunds:   big.NewInt(1000000000000000000),

		testUserKey:     testUserKey,
		testUserAddress: crypto.PubkeyToAddress(testUserKey.PublicKey),
		testUserFunds:   big.NewInt(1000),

		testConfig: &Config{
			Recommit: time.Second,
			GasFloor: params.GenesisGasLimit,
			GasCeil:  params.GenesisGasLimit,
		},
	}

	t.testTxPoolConfig = core.DefaultTxPoolConfig
	t.testTxPoolConfig.Journal = ""

	t.ethashChainConfig = &params.ChainConfig{}
	*t.ethashChainConfig = *params.TestChainConfig

	tx1, err := types.SignTx(types.NewTransaction(0, t.testUserAddress, t.testUserFunds, params.TxGas, nil, nil), types.HomesteadSigner{}, testBankKey)
	if err != nil {
		return nil, err
	}
	t.pendingTxs = append(t.pendingTxs, tx1)

	tx2, err := types.SignTx(types.NewTransaction(1, t.testUserAddress, t.testUserFunds, params.TxGas, nil, nil), types.HomesteadSigner{}, testBankKey)
	if err != nil {
		return nil, err
	}
	t.newTxs = append(t.newTxs, tx2)

	rand.Seed(time.Now().UnixNano())

	return t, nil
}

// testWorkerBackend implements worker.Backend interfaces and wraps all information needed during the testing.
type testWorkerBackend struct {
	db         ethdb.Database
	txPool     *core.TxPool
	chain      *core.BlockChain
	testTxFeed event.Feed
	genesis    *core.Genesis
	uncleBlock *types.Block
}

func newTestWorkerBackend(t *testing.T, testCase *testCase, chainConfig *params.ChainConfig, engine consensus.Engine, db ethdb.Database, n int) *testWorkerBackend {
	gspec := core.Genesis{
		Config: chainConfig,
		Alloc: core.GenesisAlloc{
			testCase.testBankAddress: {Balance: testCase.testBankFunds},
		},
	}

	switch engine.(type) {
	case *ethash.Ethash:
	default:
		t.Fatalf("unexpected consensus engine type: %T", engine)
	}
	genesis := gspec.MustCommit(db)

	cacher := core.NewTxSenderCacher()
	chain, _ := core.NewBlockChain(db, &core.CacheConfig{TrieDirtyDisabled: true}, gspec.Config, engine, vm.Config{}, nil, cacher)
	txpool := core.NewTxPool(testCase.testTxPoolConfig, chainConfig, chain, cacher)

	// Generate a small n-block chain and an uncle block for it
	if n > 0 {
		blocks, _ := core.GenerateChain(chainConfig, genesis, engine, db, n, func(i int, gen *core.BlockGen) {
			gen.SetCoinbase(testCase.testBankAddress)
		})
		if _, err := chain.InsertChain(blocks); err != nil {
			t.Fatalf("failed to insert origin chain: %v", err)
		}
	}
	parent := genesis
	if n > 0 {
		parent = chain.GetBlockByHash(chain.CurrentBlock().ParentHash())
	}
	blocks, _ := core.GenerateChain(chainConfig, parent, engine, db, 1, func(i int, gen *core.BlockGen) {
		gen.SetCoinbase(testCase.testUserAddress)
	})

	return &testWorkerBackend{
		db:         db,
		chain:      chain,
		txPool:     txpool,
		genesis:    &gspec,
		uncleBlock: blocks[0],
	}
}

func (b *testWorkerBackend) BlockChain() *core.BlockChain { return b.chain }
func (b *testWorkerBackend) TxPool() *core.TxPool         { return b.txPool }

func newTestWorker(t *testCase, chainConfig *params.ChainConfig, engine consensus.Engine, backend Backend, h hooks, waitInit bool) *worker {
	w := newWorker(t.testConfig, chainConfig, engine, backend, new(event.TypeMux), h, true)
	w.setEtherbase(t.testBankAddress)
	if waitInit {
		w.init()

		// Ensure worker has finished initialization
		timer := time.NewTicker(10 * time.Millisecond)
		defer timer.Stop()
		for range timer.C {
			b := w.pendingBlock()
			if b != nil && b.NumberU64() >= 1 {
				break
			}
		}
	}
	return w
}

func newTestBackend(t *testing.T, testCase *testCase, chainConfig *params.ChainConfig, engine consensus.Engine, db ethdb.Database, blocks int) *testWorkerBackend {
	backend := newTestWorkerBackend(t, testCase, chainConfig, engine, db, blocks)

	errs := backend.txPool.AddLocals(testCase.pendingTxs)
	for _, err := range errs {
		if err != nil {
			t.Fatal(errs)
		}
	}

	return backend
}

func (b *testWorkerBackend) newRandomUncle() *types.Block {
	var parent *types.Block
	cur := b.chain.CurrentBlock()
	if cur.NumberU64() == 0 {
		parent = b.chain.Genesis()
	} else {
		parent = b.chain.GetBlockByHash(b.chain.CurrentBlock().ParentHash())
	}
	blocks, _ := core.GenerateChain(b.chain.Config(), parent, b.chain.Engine(), b.db, 1, func(i int, gen *core.BlockGen) {
		var addr = make([]byte, common.AddressLength)
		rand.Read(addr)
		gen.SetCoinbase(common.BytesToAddress(addr))
	})
	return blocks[0]
}

func (b *testWorkerBackend) newRandomTx(testCase *testCase, creation bool) *types.Transaction {
	var tx *types.Transaction
	if creation {
		tx, _ = types.SignTx(types.NewContractCreation(b.txPool.Nonce(testCase.testBankAddress), big.NewInt(0), testGas, nil, common.FromHex(testCode)), types.HomesteadSigner{}, testCase.testBankKey)
	} else {
		tx, _ = types.SignTx(types.NewTransaction(b.txPool.Nonce(testCase.testBankAddress), testCase.testUserAddress, big.NewInt(1000), params.TxGas, nil, nil), types.HomesteadSigner{}, testCase.testBankKey)
	}
	return tx
}

func TestGenerateBlockAndImportEthash(t *testing.T) {
	testCase, err := getTestCase()
	if err != nil {
		t.Error(err)
	}
	testGenerateBlockAndImport(t, testCase)
}

func testGenerateBlockAndImport(t *testing.T, testCase *testCase) {
	var (
		engine consensus.Engine
		db     = rawdb.NewMemoryDatabase()
	)
	testCase.ethashChainConfig = params.AllEthashProtocolChanges
	engine = ethash.NewFaker()

	b := newTestBackend(t, testCase, testCase.ethashChainConfig, engine, db, 0)
	w := newTestWorker(testCase, testCase.ethashChainConfig, engine, b, hooks{}, false)
	defer w.close()

	db2 := rawdb.NewMemoryDatabase()
	b.genesis.MustCommit(db2)
	chain, _ := core.NewBlockChain(db2, nil, b.chain.Config(), engine, vm.Config{}, nil, core.NewTxSenderCacher())
	defer chain.Stop()

	loopErr := make(chan error)
	newBlock := make(chan struct{})
	listenNewBlock := func() {
		sub := w.mux.Subscribe(core.NewMinedBlockEvent{})
		defer sub.Unsubscribe()

		for item := range sub.Chan() {
			block := item.Data.(core.NewMinedBlockEvent).Block
			_, err := chain.InsertChain([]*types.Block{block})
			if err != nil {
				loopErr <- fmt.Errorf("failed to insert new mined block:%d, error:%v", block.NumberU64(), err)
			}
			newBlock <- struct{}{}
		}
	}
	// Ignore empty commit here for less noise
	w.skipSealHook = func(task *task) bool {
		return len(task.receipts) == 0
	}
	w.start() // Start mining!
	go listenNewBlock()

	for i := 0; i < 5; i++ {
		if err := b.txPool.AddLocal(b.newRandomTx(testCase, true)); err != nil {
			t.Fatal(err)
		}
		if err := b.txPool.AddLocal(b.newRandomTx(testCase, false)); err != nil {
			t.Fatal(err)
		}

		w.postSideBlock(core.ChainSideEvent{Block: b.newRandomUncle()})
		w.postSideBlock(core.ChainSideEvent{Block: b.newRandomUncle()})

		select {
		case e := <-loopErr:
			t.Fatal(e)
		case <-newBlock:
		case <-time.NewTimer(3 * time.Second).C: // Worker needs 1s to include new changes.
			t.Fatalf("timeout")
		}
	}
}

func TestEmptyWorkEthash(t *testing.T) {
	testCase, err := getTestCase()
	if err != nil {
		t.Error(err)
	}
	testEmptyWork(t, testCase, testCase.ethashChainConfig, ethash.NewFaker())
}

func testEmptyWork(t *testing.T, testCase *testCase, chainConfig *params.ChainConfig, engine consensus.Engine) {
	defer engine.Close()

	backend := newTestBackend(t, testCase, chainConfig, engine, rawdb.NewMemoryDatabase(), 0)

	var (
		taskIndex int
		taskCh    = make(chan struct{}, 2)

		h hooks
	)
	checkEqual := func(t *testing.T, task *task, index int) {
		// The first empty work without any txs included
		receiptLen, balance := 0, big.NewInt(0)
		if index == 1 {
			// The second full work with 1 tx included
			receiptLen, balance = 1, big.NewInt(1000)
		}
		if len(task.receipts) != receiptLen {
			t.Fatalf("receipt number mismatch: have %d, want %d", len(task.receipts), receiptLen)
		}
		if task.state.GetBalance(testCase.testUserAddress).Cmp(balance) != 0 {
			t.Fatalf("account balance mismatch: have %d, want %d", task.state.GetBalance(testCase.testUserAddress), balance)
		}
	}
	h.newTaskHook = func(task *task) {
		if task.block.NumberU64() == 1 {
			checkEqual(t, task, taskIndex)
			taskIndex += 1
			taskCh <- struct{}{}
		}
	}
	h.skipSealHook = func(task *task) bool { return true }
	h.fullTaskHook = func() {
		// Arch64 unit tests are running in a VM on travis, they must
		// be given more time to execute.
		time.Sleep(time.Second)
	}

	w := newTestWorker(testCase, chainConfig, engine, backend, h, true)
	defer w.close()

	w.start() // Start mining!
	for i := 0; i < 2; i += 1 {
		select {
		case <-taskCh:
		case <-time.NewTimer(3 * time.Second).C:
			t.Error("new task timeout")
		}
	}
}

func TestStreamUncleBlock(t *testing.T) {
	ethash := ethash.NewFaker()
	defer ethash.Close()

	testCase, err := getTestCase()
	if err != nil {
		t.Error(err)
	}

	b := newTestBackend(t, testCase, testCase.ethashChainConfig, ethash, rawdb.NewMemoryDatabase(), 1)

	var taskCh = make(chan struct{})
	var h hooks
	taskIndex := 0
	h.newTaskHook = func(task *task) {
		if task.block.NumberU64() == 2 {
			// The first task is an empty task, the second
			// one has 1 pending tx, the third one has 1 tx
			// and 1 uncle.
			if taskIndex == 2 {
				have := task.block.Header().UncleHash
				want := types.CalcUncleHash([]*types.Header{b.uncleBlock.Header()})
				if have != want {
					t.Errorf("uncle hash mismatch: have %s, want %s", have.Hex(), want.Hex())
				}
			}
			taskCh <- struct{}{}
			taskIndex += 1
		}
	}
	h.skipSealHook = func(task *task) bool {
		return true
	}
	h.fullTaskHook = func() {
		time.Sleep(100 * time.Millisecond)
	}

	w := newTestWorker(testCase, testCase.ethashChainConfig, ethash, b, h, true)
	defer w.close()

	w.start()

	for i := 0; i < 2; i += 1 {
		select {
		case <-taskCh:
		case <-time.NewTimer(time.Second).C:
			t.Error("new task timeout")
		}
	}

	w.postSideBlock(core.ChainSideEvent{Block: b.uncleBlock})

	select {
	case <-taskCh:
	case <-time.NewTimer(time.Second).C:
		t.Error("new task timeout")
	}
}

func TestRegenerateMiningBlockEthash(t *testing.T) {
	testCase, err := getTestCase()
	if err != nil {
		t.Error(err)
	}
	testRegenerateMiningBlock(t, testCase, testCase.ethashChainConfig, ethash.NewFaker())
}

func testRegenerateMiningBlock(t *testing.T, testCase *testCase, chainConfig *params.ChainConfig, engine consensus.Engine) {
	defer engine.Close()

	b := newTestBackend(t, testCase, chainConfig, engine, rawdb.NewMemoryDatabase(), 0)

	var taskCh = make(chan struct{})
	var h hooks

	taskIndex := 0
	h.newTaskHook = func(task *task) {
		if task.block.NumberU64() == 1 {
			// The first task is an empty task, the second
			// one has 1 pending tx, the third one has 2 txs
			if taskIndex == 2 {
				receiptLen, balance := 2, big.NewInt(2000)
				if len(task.receipts) != receiptLen {
					t.Errorf("receipt number mismatch: have %d, want %d", len(task.receipts), receiptLen)
				}
				if task.state.GetBalance(testCase.testUserAddress).Cmp(balance) != 0 {
					t.Errorf("account balance mismatch: have %d, want %d", task.state.GetBalance(testCase.testUserAddress), balance)
				}
			}
			taskCh <- struct{}{}
			taskIndex += 1
		}
	}
	h.skipSealHook = func(task *task) bool {
		return true
	}
	h.fullTaskHook = func() {
		time.Sleep(100 * time.Millisecond)
	}

	w := newTestWorker(testCase, chainConfig, engine, b, h, true)
	defer w.close()

	w.start()
	// Ignore the first two works
	for i := 0; i < 2; i += 1 {
		select {
		case <-taskCh:
		case <-time.NewTimer(time.Second).C:
			t.Error("new task timeout on first 2 works")
		}
	}
	b.txPool.AddLocals(testCase.newTxs)
	time.Sleep(time.Second)

	select {
	case <-taskCh:
	case <-time.NewTimer(time.Second).C:
		t.Error("new task timeout")
	}
}

func TestAdjustIntervalEthash(t *testing.T) {
	testCase, err := getTestCase()
	if err != nil {
		t.Error(err)
	}
	testAdjustInterval(t, testCase, testCase.ethashChainConfig, ethash.NewFaker())
}

func testAdjustInterval(t *testing.T, testCase *testCase, chainConfig *params.ChainConfig, engine consensus.Engine) {
	defer engine.Close()

	backend := newTestBackend(t, testCase, chainConfig, engine, rawdb.NewMemoryDatabase(), 0)

	var h hooks
	h.skipSealHook = func(task *task) bool {
		return true
	}
	h.fullTaskHook = func() {
		time.Sleep(100 * time.Millisecond)
	}
	var (
		progress = make(chan struct{}, 10)
		result   = make([]float64, 0, 10)
		index    = 0
		start    uint32
	)
	h.resubmitHook = func(minInterval time.Duration, recommitInterval time.Duration) {
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

	w := newTestWorker(testCase, chainConfig, engine, backend, h, true)
	defer w.close()

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
