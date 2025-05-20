// Copyright 2015 The go-ethereum Authors
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

package core

import (
	"fmt"
	"math/big"

	"github.com/autonity/autonity/triedb/pathdb"
	"github.com/holiman/uint256"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/misc"
	"github.com/autonity/autonity/core/rawdb"
	"github.com/autonity/autonity/core/state"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/ethdb"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/triedb"
)

// BlockGen creates blocks for testing.
// See GenerateChain for a detailed explanation.
type BlockGen struct {
	i       int
	cm      *chainMaker
	parent  *types.Block
	header  *types.Header
	statedb *state.StateDB

	gasPool  *GasPool
	txs      []*types.Transaction
	receipts []*types.Receipt
	uncles   []*types.Header

	engine consensus.Engine
}

// SetCoinbase sets the coinbase of the generated block.
// It can be called at most once.
func (b *BlockGen) SetCoinbase(addr common.Address) {
	if b.gasPool != nil {
		if len(b.txs) > 0 {
			panic("coinbase must be set before adding transactions")
		}
		panic("coinbase can only be set once")
	}
	b.header.Coinbase = addr
	b.gasPool = new(GasPool).AddGas(b.header.GasLimit)
}

// SetExtra sets the extra data field of the generated block.
func (b *BlockGen) SetExtra(data []byte) {
	b.header.Extra = data
}

// SetNonce sets the nonce field of the generated block.
func (b *BlockGen) SetNonce(nonce types.BlockNonce) {
	b.header.Nonce = nonce
}

// SetDifficulty sets the difficulty field of the generated block. This method is
// useful for Clique tests where the difficulty does not depend on time. For the
// ethash tests, please use OffsetTime, which implicitly recalculates the diff.
func (b *BlockGen) SetDifficulty(diff *big.Int) {
	b.header.Difficulty = diff
}

// addTx adds a transaction to the generated block. If no coinbase has
// been set, the block's coinbase is set to the zero address.
//
// There are a few options can be passed as well in order to run some
// customized rules.
// - bc:       enables the ability to query historical block hashes for BLOCKHASH
// - vmConfig: extends the flexibility for customizing evm rules, e.g. enable extra EIPs
func (b *BlockGen) addTx(bc *BlockChain, vmConfig vm.Config, tx *types.Transaction) {
	if b.gasPool == nil {
		b.SetCoinbase(common.Address{})
	}
	var (
		blockContext = NewEVMBlockContext(b.header, bc, &b.header.Coinbase)
		evm          = vm.NewEVM(blockContext, b.statedb, b.cm.config, vmConfig)
	)
	b.statedb.SetTxContext(tx.Hash(), len(b.txs))
	receipt, err := ApplyTransaction(evm, b.gasPool, b.statedb, b.header, tx, &b.header.GasUsed)
	if err != nil {
		panic(err)
	}
	// Merge the tx-local access event into the "block-local" one, in order to collect
	// all values, so that the witness can be built.
	if b.statedb.GetTrie().IsVerkle() {
		b.statedb.AccessEvents().Merge(evm.AccessEvents)
	}
	b.txs = append(b.txs, tx)
	b.receipts = append(b.receipts, receipt)
}

// AddTx adds a transaction to the generated block. If no coinbase has
// been set, the block's coinbase is set to the zero address.
//
// AddTx panics if the transaction cannot be executed. In addition to the protocol-imposed
// limitations (gas limit, etc.), there are some further limitations on the content of
// transactions that can be added. Notably, contract code relying on the BLOCKHASH
// instruction will panic during execution if it attempts to access a block number outside
// of the range created by GenerateChain.
func (b *BlockGen) AddTx(tx *types.Transaction) {
	// Wrap the chain config in an empty BlockChain object to satisfy ChainContext.
	bc := &BlockChain{chainConfig: b.cm.config}
	b.addTx(bc, vm.Config{}, tx)
}

// AddTxWithChain adds a transaction to the generated block. If no coinbase has
// been set, the block's coinbase is set to the zero address.
//
// AddTxWithChain panics if the transaction cannot be executed. In addition to the
// protocol-imposed limitations (gas limit, etc.), there are some further limitations on
// the content of transactions that can be added. If contract code relies on the BLOCKHASH
// instruction, the block in chain will be returned.
func (b *BlockGen) AddTxWithChain(bc *BlockChain, tx *types.Transaction) {
	b.addTx(bc, vm.Config{}, tx)
}

// AddTxWithVMConfig adds a transaction to the generated block. If no coinbase has
// been set, the block's coinbase is set to the zero address.
// The evm interpreter can be customized with the provided vm config.
func (b *BlockGen) AddTxWithVMConfig(tx *types.Transaction, config vm.Config) {
	b.addTx(nil, config, tx)
}

// GetBalance returns the balance of the given address at the generated block.
func (b *BlockGen) GetBalance(addr common.Address) *uint256.Int {
	return b.statedb.GetBalance(addr)
}

// AddUncheckedTx forcefully adds a transaction to the block without any
// validation.
//
// AddUncheckedTx will cause consensus failures when used during real
// chain processing. This is best used in conjunction with raw block insertion.
func (b *BlockGen) AddUncheckedTx(tx *types.Transaction) {
	b.txs = append(b.txs, tx)
}

// Number returns the block number of the block being generated.
func (b *BlockGen) Number() *big.Int {
	return new(big.Int).Set(b.header.Number)
}

// Timestamp returns the timestamp of the block being generated.
func (b *BlockGen) Timestamp() uint64 {
	return b.header.Time
}

// BaseFee returns the EIP-1559 base fee of the block being generated.
func (b *BlockGen) BaseFee() *big.Int {
	return new(big.Int).Set(b.header.BaseFee)
}

// Gas returns the amount of gas left in the current block.
func (b *BlockGen) Gas() uint64 {
	return b.header.GasLimit - b.header.GasUsed
}

// Signer returns a valid signer instance for the current block.
func (b *BlockGen) Signer() types.Signer {
	return types.MakeSigner(b.cm.config, b.header.Number)
}

// AddUncheckedReceipt forcefully adds a receipts to the block without a
// backing transaction.
//
// AddUncheckedReceipt will cause consensus failures when used during real
// chain processing. This is best used in conjunction with raw block insertion.
func (b *BlockGen) AddUncheckedReceipt(receipt *types.Receipt) {
	b.receipts = append(b.receipts, receipt)
}

// TxNonce returns the next valid transaction nonce for the
// account at addr. It panics if the account does not exist.
func (b *BlockGen) TxNonce(addr common.Address) uint64 {
	if !b.statedb.Exist(addr) {
		panic("account does not exist")
	}
	return b.statedb.GetNonce(addr)
}

// AddUncle adds an uncle header to the generated block.
func (b *BlockGen) AddUncle(h *types.Header) {
	// The uncle will have the same timestamp and auto-generated difficulty
	h.Time = b.header.Time

	var parent *types.Header
	for i := b.i - 1; i >= 0; i-- {
		if b.cm.chain[i].Hash() == h.ParentHash {
			parent = b.cm.chain[i].Header()
			break
		}
	}
	h.Difficulty = b.engine.CalcDifficulty(b.cm, b.header.Time, parent)

	// The gas limit and price should be derived from the parent
	h.GasLimit = parent.GasLimit
	if b.cm.config.IsLondon(h.Number) {
		h.BaseFee = misc.CalcBaseFee(b.cm.config, parent, b.cm)
		if !b.cm.config.IsLondon(parent.Number) {
			parentGasLimit := parent.GasLimit * b.cm.config.ElasticityMultiplier()
			h.GasLimit = CalcGasLimit(parentGasLimit, parentGasLimit)
		}
	}
	b.uncles = append(b.uncles, h)
}

// PrevBlock returns a previously generated block by number. It panics if
// num is greater or equal to the number of the block being generated.
// For index -1, PrevBlock returns the parent block given to GenerateChain.
func (b *BlockGen) PrevBlock(index int) *types.Block {
	if index >= b.i {
		panic(fmt.Errorf("block index %d out of range (%d,%d)", index, -1, b.i))
	}
	if index == -1 {
		return b.cm.bottom
	}
	return b.cm.chain[index]
}

// OffsetTime modifies the time instance of a block, implicitly changing its
// associated difficulty. It's useful to test scenarios where forking is not
// tied to chain length directly.
func (b *BlockGen) OffsetTime(seconds int64) {
	b.header.Time += uint64(seconds)
	if b.header.Time <= b.cm.bottom.Header().Time {
		panic("block time out of range")
	}
	//b.header.Difficulty = b.engine.CalcDifficulty(b.cm, b.header.Time, b.parent.Header())
}

// GenerateChain creates a chain of n blocks. The first block's
// parent will be the provided parent. db is used to store
// intermediate states and should contain the parent's state trie.
//
// The generator function is called with a new block generator for
// every block. Any transactions and uncles added to the generator
// become part of the block. If gen is nil, the blocks will be empty
// and their coinbase will be the zero address.
//
// Blocks created by GenerateChain do not contain valid proof of work
// values. Inserting them into BlockChain requires use of FakePow or
// a similar non-validating proof of work implementation.
func GenerateChain(config *params.ChainConfig, parent *types.Block, engine consensus.Engine, db ethdb.Database, n int, gen func(int, *BlockGen)) ([]*types.Block, []types.Receipts) {
	if config == nil {
		config = params.TestChainConfig
	}

	blocks, receipts := make(types.Blocks, n), make([]types.Receipts, n)

	cm := newChainMaker(parent, config, engine)

	genblock := func(i int, parent *types.Block, triedb *triedb.Database, statedb *state.StateDB) (*types.Block, types.Receipts) {
		b := &BlockGen{i: i, cm: cm, parent: parent, statedb: statedb, engine: engine}
		b.header = cm.makeHeader(parent, statedb, b.engine)

		// Set the difficulty for clique block. The chain maker doesn't have access
		// to a chain, so the difficulty will be left unset (nil). Set it here to the
		// correct value.
		if b.header.Difficulty == nil {
			if config.TerminalTotalDifficulty == nil {
				// Clique chain
				b.header.Difficulty = big.NewInt(2)
			} else {
				// Post-merge chain
				b.header.Difficulty = big.NewInt(0)
			}
		}

		// Mutate the state and block according to any hard-fork specs
		if daoBlock := config.DAOForkBlock; daoBlock != nil {
			limit := new(big.Int).Add(daoBlock, params.DAOForkExtraRange)
			if b.header.Number.Cmp(daoBlock) >= 0 && b.header.Number.Cmp(limit) < 0 {
				if config.DAOForkSupport {
					b.header.Extra = common.CopyBytes(params.DAOForkBlockExtra)
				}
			}
		}
		if config.DAOForkSupport && config.DAOForkBlock != nil && config.DAOForkBlock.Cmp(b.header.Number) == 0 {
			misc.ApplyDAOHardFork(statedb)
		}

		if config.IsPrague(b.header.Number) || config.IsVerkle(b.header.Number) {
			// EIP-2935
			/*
				blockContext := NewEVMBlockContext(b.header, cm, &b.header.Coinbase)
				blockContext.Random = &common.Hash{} // enable post-merge instruction set
				evm := vm.NewEVM(blockContext, statedb, cm.config, vm.Config{})
				ProcessParentBlockHash(b.header.ParentHash, evm)

			*/
		}

		// Execute any user modifications to the block
		if gen != nil {
			gen(i, b)
		}
		body := types.Body{Transactions: b.txs, Uncles: b.uncles}
		block, err := b.engine.FinalizeAndAssemble(cm, b.header, statedb, &body, &b.receipts)
		if err != nil {
			panic(err)
		}

		// Write state changes to db
		root, err := statedb.Commit(b.header.Number.Uint64(), config.IsEIP158(b.header.Number), config.IsCancun(b.header.Number))
		if err != nil {
			panic(fmt.Sprintf("state write error: %v", err))
		}
		if err = triedb.Commit(root, false); err != nil {
			panic(fmt.Sprintf("trie write error: %v", err))
		}
		cm.add(block, b.receipts)
		return block, b.receipts
	}

	triedCfg := *triedb.HashDefaults
	if config.VerkleBlock != nil {
		triedCfg.PathDB = pathdb.Defaults
		triedCfg.IsVerkle = true
		triedCfg.HashDB = nil
	}
	tdb := triedb.NewDatabase(db, &triedCfg)
	defer tdb.Close()

	for i := 0; i < n; i++ {
		statedb, err := state.New(parent.Root(), state.NewDatabase(tdb, nil))
		if err != nil {
			panic(err)
		}
		block, receipt := genblock(i, parent, tdb, statedb)
		blocks[i] = block
		receipts[i] = receipt
		parent = block
	}
	return blocks, receipts
}

// GenerateChainWithGenesis is a wrapper of GenerateChain which will initialize
// genesis block to database first according to the provided genesis specification
// then generate chain on top.
func GenerateChainWithGenesis(genesis *Genesis, engine consensus.Engine, n int, gen func(int, *BlockGen)) (ethdb.Database, []*types.Block, []types.Receipts) {
	db := rawdb.NewMemoryDatabase()
	triedb := triedb.NewDatabase(db, triedb.HashDefaults)
	defer triedb.Close()
	g, err := genesis.Commit(db, triedb)
	if err != nil {
		panic(err)
	}
	blocks, receipts := GenerateChain(genesis.Config, g, engine, db, n, gen)
	return db, blocks, receipts
}

// makeHeaderChainWithGenesis creates a deterministic chain of headers from genesis.
func makeHeaderChainWithGenesis(genesis *Genesis, n int, engine consensus.Engine, seed int) (ethdb.Database, []*types.Header) {
	db, blocks := makeBlockChainWithGenesis(genesis, n, engine, seed)
	headers := make([]*types.Header, len(blocks))
	for i, block := range blocks {
		headers[i] = block.Header()
	}
	return db, headers
}

// makeBlockChainWithGenesis creates a deterministic chain of blocks from genesis
func makeBlockChainWithGenesis(genesis *Genesis, n int, engine consensus.Engine, seed int) (ethdb.Database, []*types.Block) {
	db, blocks, _ := GenerateChainWithGenesis(genesis, engine, n, func(i int, b *BlockGen) {
		b.SetCoinbase(common.Address{0: byte(seed), 19: byte(i)})
	})
	return db, blocks
}

func makeHeader(chain consensus.ChainReader, parent *types.Block, state *state.StateDB, engine consensus.Engine) *types.Header {
	var time uint64
	if parent.Time() == 0 {
		time = 10
	} else {
		time = parent.Time() + 10 // block time is fixed at 10 seconds
	}
	header := &types.Header{
		Root:       state.IntermediateRoot(chain.Config().IsEIP158(parent.Number())),
		ParentHash: parent.Hash(),
		Coinbase:   parent.Coinbase(),
		BaseFee:    parent.BaseFee(),
		Difficulty: engine.CalcDifficulty(chain, time, &types.Header{
			Number:     parent.Number(),
			Time:       time - 10,
			Difficulty: parent.Difficulty(),
			UncleHash:  parent.UncleHash(),
		}),
		GasLimit: parent.GasLimit(),
		Number:   new(big.Int).Add(parent.Number(), common.Big1),
		Time:     time,
	}
	if chain.Config().IsLondon(header.Number) {
		header.BaseFee = misc.CalcBaseFee(chain.Config(), parent.Header(), chain)
		if !chain.Config().IsLondon(parent.Number()) {
			parentGasLimit := parent.GasLimit() * params.ElasticityMultiplier
			header.GasLimit = CalcGasLimit(parentGasLimit, parentGasLimit)
		}
	}
	return header
}

// makeHeaderChain creates a deterministic chain of headers rooted at parent.
func makeHeaderChain(chainConfig *params.ChainConfig, parent *types.Header, n int, engine consensus.Engine, db ethdb.Database, seed int) []*types.Header {
	blocks := makeBlockChain(chainConfig, types.NewBlockWithHeader(parent), n, engine, db, seed)
	headers := make([]*types.Header, len(blocks))
	for i, block := range blocks {
		headers[i] = block.Header()
	}
	return headers
}

// makeBlockChain creates a deterministic chain of blocks rooted at parent.
func makeBlockChain(chainConfig *params.ChainConfig, parent *types.Block, n int, engine consensus.Engine, db ethdb.Database, seed int) []*types.Block {
	blocks, _ := GenerateChain(chainConfig, parent, engine, db, n, func(i int, b *BlockGen) {
		b.SetCoinbase(common.Address{0: byte(seed), 19: byte(i)})
	})
	return blocks
}

type fakeChainReader struct {
	config    *params.ChainConfig
	committee *types.Committee
}

// Config returns the chain configuration.
func (cr *fakeChainReader) Config() *params.ChainConfig {
	return cr.config
}

func (cr *fakeChainReader) CurrentHeader() *types.Header { return nil }
func (cr *fakeChainReader) GetHeaderByNumber(number uint64) *types.Header {
	if cr.committee == nil {
		return nil
	}
	header := &types.Header{Number: new(big.Int).SetUint64(number)}
	var epoch types.Epoch
	epoch.Committee = cr.committee.Copy()
	epoch.PreviousEpochBlock = common.Big0
	epoch.NextEpochBlock = new(big.Int).SetUint64(number + 30)
	epoch.Delta = new(big.Int).SetUint64(10)
	header.Epoch = &epoch
	return header
}
func (cr *fakeChainReader) GetHeaderByHash(hash common.Hash) *types.Header          { return nil }
func (cr *fakeChainReader) GetHeader(hash common.Hash, number uint64) *types.Header { return nil }
func (cr *fakeChainReader) GetBlock(hash common.Hash, number uint64) *types.Block   { return nil }
func (cr *fakeChainReader) Engine() consensus.Engine                                { return nil }
func (cr *fakeChainReader) GetTd(hash common.Hash, number uint64) *big.Int          { return nil }
func (cr *fakeChainReader) MinBaseFee() *big.Int {
	return big.NewInt(0)
}
func (cr *fakeChainReader) EpochByHeight(_ uint64) (*types.EpochInfo, error) {
	return nil, nil
}

// chainMaker contains the state of chain generation.
type chainMaker struct {
	bottom      *types.Block
	engine      consensus.Engine
	config      *params.ChainConfig
	chain       []*types.Block
	chainByHash map[common.Hash]*types.Block
	receipts    []types.Receipts
}

func (cm *chainMaker) MinBaseFee() *big.Int {
	return big.NewInt(0)
}

func (cm *chainMaker) GetTd(hash common.Hash, number uint64) *big.Int {
	//TODO implement me
	panic("implement me")
}

func (cm *chainMaker) EpochByHeight(height uint64) (*types.EpochInfo, error) {
	//TODO implement me
	panic("implement me")
}

func newChainMaker(bottom *types.Block, config *params.ChainConfig, engine consensus.Engine) *chainMaker {
	return &chainMaker{
		bottom:      bottom,
		config:      config,
		engine:      engine,
		chainByHash: make(map[common.Hash]*types.Block),
	}
}

func (cm *chainMaker) add(b *types.Block, r []*types.Receipt) {
	cm.chain = append(cm.chain, b)
	cm.chainByHash[b.Hash()] = b
	cm.receipts = append(cm.receipts, r)
}

func (cm *chainMaker) blockByNumber(number uint64) *types.Block {
	if number == cm.bottom.NumberU64() {
		return cm.bottom
	}
	cur := cm.CurrentHeader().Number.Uint64()
	lowest := cm.bottom.NumberU64() + 1
	if number < lowest || number > cur {
		return nil
	}
	return cm.chain[number-lowest]
}

// ChainReader/ChainContext implementation

func (cm *chainMaker) makeHeader(parent *types.Block, state *state.StateDB, engine consensus.Engine) *types.Header {
	time := parent.Time() + 10 // block time is fixed at 10 seconds
	parentHeader := parent.Header()
	header := &types.Header{
		Root:       state.IntermediateRoot(cm.config.IsEIP158(parent.Number())),
		ParentHash: parent.Hash(),
		//MixDigest:  types.BFTDigest,
		Coinbase:   parent.Coinbase(),
		Difficulty: parent.Difficulty(),
		GasLimit:   parent.GasLimit(),
		Number:     new(big.Int).Add(parent.Number(), common.Big1),
		Time:       time,
	}

	if cm.config.IsLondon(header.Number) {
		header.BaseFee = misc.CalcBaseFee(cm.config, parentHeader, cm)
		if !cm.config.IsLondon(parent.Number()) {
			parentGasLimit := parent.GasLimit() * cm.config.ElasticityMultiplier()
			header.GasLimit = CalcGasLimit(parentGasLimit, parentGasLimit)
		}
	}
	return header
}

// Config returns the chain configuration (for consensus.ChainReader).
func (cm *chainMaker) Config() *params.ChainConfig {
	return cm.config
}

// Engine returns the consensus engine (for ChainContext).
func (cm *chainMaker) Engine() consensus.Engine {
	return cm.engine
}

func (cm *chainMaker) CurrentHeader() *types.Header {
	if len(cm.chain) == 0 {
		return cm.bottom.Header()
	}
	return cm.chain[len(cm.chain)-1].Header()
}

func (cm *chainMaker) GetHeaderByNumber(number uint64) *types.Header {
	b := cm.blockByNumber(number)
	if b == nil {
		return nil
	}
	return b.Header()
}

func (cm *chainMaker) GetHeaderByHash(hash common.Hash) *types.Header {
	b := cm.chainByHash[hash]
	if b == nil {
		return nil
	}
	return b.Header()
}

func (cm *chainMaker) GetHeader(hash common.Hash, number uint64) *types.Header {
	return cm.GetHeaderByNumber(number)
}

func (cm *chainMaker) GetBlock(hash common.Hash, number uint64) *types.Block {
	return cm.blockByNumber(number)
}
