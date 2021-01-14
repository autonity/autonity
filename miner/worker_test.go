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

	"github.com/clearmatics/autonity/consensus/tendermint"
	"github.com/stretchr/testify/require"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/ethash"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/rawdb"
	"github.com/clearmatics/autonity/core/state"
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

var (
	// Test chain configurations
	testTxPoolConfig      core.TxPoolConfig
	ethashChainConfig     *params.ChainConfig
	tendermintChainConfig *params.ChainConfig

	// Test accounts
	testBankKey, _  = crypto.GenerateKey()
	testBankAddress = crypto.PubkeyToAddress(testBankKey.PublicKey)
	testBankFunds   = big.NewInt(1000000000000000000)

	testUserKey, _  = crypto.GenerateKey()
	testUserAddress = crypto.PubkeyToAddress(testUserKey.PublicKey)

	// Test transactions
	pendingTxs []*types.Transaction
	newTxs     []*types.Transaction

	testConfig = &Config{
		Recommit: time.Second,
		GasFloor: params.GenesisGasLimit,
		GasCeil:  params.GenesisGasLimit,
	}
)

func init() {
	testTxPoolConfig = core.DefaultTxPoolConfig
	testTxPoolConfig.Journal = ""
	ethashChainConfig = params.TestChainConfig
	tendermintChainConfig = params.AutonityTestChainConfig
	tendermintChainConfig.AutonityContractConfig = &params.AutonityContractGenesis{
		Users: []params.User{{
			Address: &testUserAddress,
			Enode:   "enode://d73b857969c86415c0c000371bcebd9ed3cca6c376032b3f65e58e9e2b79276fbc6f59eb1e22fcd6356ab95f42a666f70afd4985933bd8f3e05beb1a2bf8fdde@172.25.0.11:30303",
			Type:    params.UserValidator,
			Stake:   1,
		}},
	}
	tendermintChainConfig.AutonityContractConfig.Prepare()
	tx1, _ := types.SignTx(types.NewTransaction(0, testUserAddress, big.NewInt(1000), params.TxGas, nil, nil), types.HomesteadSigner{}, testBankKey)
	pendingTxs = append(pendingTxs, tx1)
	tx2, _ := types.SignTx(types.NewTransaction(1, testUserAddress, big.NewInt(1000), params.TxGas, nil, nil), types.HomesteadSigner{}, testBankKey)
	newTxs = append(newTxs, tx2)
	rand.Seed(time.Now().UnixNano())
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

func newTestWorkerBackend(t *testing.T, chainConfig *params.ChainConfig, engine consensus.Engine, db ethdb.Database, n int) *testWorkerBackend {
	var gspec = core.Genesis{
		Config:     chainConfig,
		Alloc:      core.GenesisAlloc{testBankAddress: {Balance: testBankFunds}},
		Difficulty: big.NewInt(1),
	}

	switch engine.(type) {
	case *tendermint.Bridge:
		gspec.Mixhash = types.BFTDigest
	case *ethash.Ethash:
	default:
		t.Fatalf("unexpected consensus engine type: %T", engine)
	}

	genesis := gspec.MustCommit(db)
	senderCacher := &core.TxSenderCacher{}
	chain, _ := core.NewBlockChain(db, &core.CacheConfig{TrieDirtyDisabled: true}, gspec.Config, engine, vm.Config{}, nil, senderCacher, nil)
	txpool := core.NewTxPool(testTxPoolConfig, chainConfig, chain, senderCacher)

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
		parent = chain.GetBlockByHash(chain.CurrentBlock().ParentHash())
	}
	blocks, _ := core.GenerateChain(chainConfig, parent, engine, db, 1, func(i int, gen *core.BlockGen) {
		gen.SetCoinbase(testUserAddress)
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

func (b *testWorkerBackend) newRandomTx(creation bool) *types.Transaction {
	var tx *types.Transaction
	if creation {
		tx, _ = types.SignTx(types.NewContractCreation(b.txPool.Nonce(testBankAddress), big.NewInt(0), testGas, nil, common.FromHex(testCode)), types.HomesteadSigner{}, testBankKey)
	} else {
		tx, _ = types.SignTx(types.NewTransaction(b.txPool.Nonce(testBankAddress), testUserAddress, big.NewInt(1000), params.TxGas, nil, nil), types.HomesteadSigner{}, testBankKey)
	}
	return tx
}

func newTestWorker(t *testing.T, chainConfig *params.ChainConfig, engine consensus.Engine, db ethdb.Database, blocks int) (*worker, *testWorkerBackend) {
	backend := newTestWorkerBackend(t, chainConfig, engine, db, blocks)
	backend.txPool.AddLocals(pendingTxs)
	w := newWorker(testConfig, chainConfig, engine, backend, new(event.TypeMux), nil, false)
	w.setEtherbase(testBankAddress)
	return w, backend
}

func TestGenerateBlockAndImportEthash(t *testing.T) {
	testGenerateBlockAndImport(t, false)
}

func TestGenerateBlockAndImportTendermint(t *testing.T) {
	testGenerateBlockAndImport(t, true)
}

type mockPeers struct{}

func (m *mockPeers) Peers() []consensus.Peer {
	return nil
}

func testGenerateBlockAndImport(t *testing.T, isTendermint bool) {
	var (
		engine      consensus.Engine
		chainConfig *params.ChainConfig
		db          = rawdb.NewMemoryDatabase()
		err         error
	)
	if isTendermint {
		chainConfig = tendermintChainConfig
		engine, db, err = newBackend(chainConfig)
		require.NoError(t, err)
	} else {
		chainConfig = params.AllEthashProtocolChanges
		engine = ethash.NewFaker()
	}

	w, b := newTestWorker(t, chainConfig, engine, db, 0)
	defer w.close()

	// This test chain imports the mined blocks.
	db2 := rawdb.NewMemoryDatabase()
	b.genesis.MustCommit(db2)
	chain, _ := core.NewBlockChain(db2, nil, b.chain.Config(), engine, vm.Config{}, nil, &core.TxSenderCacher{}, nil)
	defer chain.Stop()

	// We need to set the broadcaster so that the engine will process messages,
	// it should really be set to a real value but that would require a massive
	// rewrite of these tests
	if bridge, ok := engine.(*tendermint.Bridge); ok {
		bridge.SetExtraComponents(chain, nil)
	}

	// Ignore empty commit here for less noise.
	w.skipSealHook = func(task *task) bool {
		return len(task.receipts) == 0
	}

	// Wait for mined blocks.
	sub := w.mux.Subscribe(core.NewMinedBlockEvent{})
	defer sub.Unsubscribe()

	// Start mining!
	w.start()

	for i := 0; i < 5; i++ {
		b.txPool.AddLocal(b.newRandomTx(true))
		b.txPool.AddLocal(b.newRandomTx(false))
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
		case <-time.After(3 * time.Second): // Worker needs 1s to include new changes.
			t.Fatalf("timeout")
		}
	}
}

func TestEmptyWorkEthash(t *testing.T) {
	testEmptyWork(t, ethashChainConfig, ethash.NewFaker(), false)
}
func TestEmptyWorkTendermint(t *testing.T) {
	engine, _, err := newBackend(tendermintChainConfig)
	require.NoError(t, err)
	testEmptyWork(t, tendermintChainConfig, engine, true)
}

func testEmptyWork(t *testing.T, chainConfig *params.ChainConfig, engine consensus.Engine, isTendermint bool) {
	defer engine.Close()

	w, _ := newTestWorker(t, chainConfig, engine, rawdb.NewMemoryDatabase(), 0)
	defer w.close()

	var (
		taskIndex int
		taskCh    = make(chan struct{}, 2)
	)
	checkEqual := func(t *testing.T, task *task, index int) {
		// The first empty work without any txs included
		receiptLen, balance := 0, big.NewInt(0)
		if index == 1 {
			// The second full work with 1 tx included
			receiptLen, balance = 1, big.NewInt(1000)
		}
		if isTendermint {
			// With tendermint there is no immediate empty block creation after a newchainhead event.
			// The first mined block receipts will then cointain both the transaction receipt and the autonity contract
			// finalize function receipt.
			receiptLen, balance = 2, big.NewInt(1000)
			if index == 1 {
				receiptLen, balance = 1, big.NewInt(1000)
			}
		}
		if len(task.receipts) != receiptLen {
			t.Fatalf("receipt number mismatch: have %d, want %d", len(task.receipts), receiptLen)
		}
		if task.state.GetBalance(testUserAddress).Cmp(balance) != 0 {
			t.Fatalf("account balance mismatch: have %d, want %d", task.state.GetBalance(testUserAddress), balance)
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

	w, b := newTestWorker(t, ethashChainConfig, ethash, rawdb.NewMemoryDatabase(), 1)
	defer w.close()

	var taskCh = make(chan struct{})

	taskIndex := 0
	w.newTaskHook = func(task *task) {
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
	w.skipSealHook = func(task *task) bool {
		return true
	}
	w.fullTaskHook = func() {
		time.Sleep(100 * time.Millisecond)
	}
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
	testRegenerateMiningBlock(t, ethashChainConfig, ethash.NewFaker(), false)
}

func TestRegenerateMiningBlockTendermint(t *testing.T) {

	engine, _, err := newBackend(tendermintChainConfig)
	require.NoError(t, err)
	testRegenerateMiningBlock(t, tendermintChainConfig, engine, true)
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
				receiptLen, balance := 2, big.NewInt(2000)
				if isTendermint {
					receiptLen += 1 // Autonity Contract Finalize additional receipt
				}
				if len(task.receipts) != receiptLen {
					t.Errorf("receipt number mismatch: have %d, want %d", len(task.receipts), receiptLen)
				}
				if task.state.GetBalance(testUserAddress).Cmp(balance) != 0 {
					t.Errorf("account balance mismatch: have %d, want %d", task.state.GetBalance(testUserAddress), balance)
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
	// Ignore the first two works in case of pow
	// Ignore the first one for Tendermint
	maxSkippedCases := 2
	if isTendermint {
		maxSkippedCases = 1
	}
	for i := 0; i < maxSkippedCases; i += 1 {
		select {
		case <-taskCh:
		case <-time.NewTimer(time.Second).C:
			t.Error("new task timeout")
		}
	}
	b.txPool.AddLocals(newTxs)
	time.Sleep(time.Second)

	select {
	case <-taskCh:
	case <-time.NewTimer(time.Second).C:
		t.Error("new task timeout")
	}
}

func TestAdjustIntervalEthash(t *testing.T) {
	testAdjustInterval(t, ethashChainConfig, ethash.NewFaker())
}

func TestAdjustIntervalClique(t *testing.T) {
	engine, _, err := newBackend(tendermintChainConfig)
	require.NoError(t, err)
	testAdjustInterval(t, tendermintChainConfig, engine)
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

func newBackend(config *params.ChainConfig) (*tendermint.Bridge, ethdb.Database, error) {
	peers := &mockPeers{}
	broadcaster := tendermint.NewBroadcaster(common.Address{}, peers)
	syncer := tendermint.NewSyncer(peers)

	db := rawdb.NewMemoryDatabase()

	// generate genesis block
	g := core.DefaultGenesisBlock()
	g.Config = config
	g.Difficulty = big.NewInt(1)
	g.Mixhash = types.BFTDigest

	_, err := g.Commit(db)
	if err != nil {
		return nil, nil, err
	}
	hg, err := core.NewHeaderGetter(db)
	if err != nil {
		return nil, nil, err
	}

	vmConfig := vm.Config{}
	autonityContract, err := core.NewAutonityContractFromConfig(
		db,
		hg,
		core.NewDefaultEVMProvider(hg, vmConfig, config),
		config.AutonityContractConfig,
	)
	if err != nil {
		return nil, nil, err
	}
	statedb := state.NewDatabase(db)

	finalizer := tendermint.NewFinalizer(autonityContract)
	verifier := tendermint.NewVerifier(&vmConfig, finalizer, tendermintChainConfig.Tendermint.BlockPeriod)
	latestBlockRetriever := tendermint.NewLatestBlockRetriever(db, statedb)
	engine := tendermint.New(
		tendermintChainConfig.Tendermint,
		testUserKey,
		broadcaster,
		syncer,
		verifier,
		finalizer,
		latestBlockRetriever,
		autonityContract,
		statedb,
	)

	return engine, db, nil
}
