package backends

import (
	"context"
	"errors"
	"math/big"
	"sync"
	"time"

	ethereum "github.com/autonity/autonity"
	"github.com/autonity/autonity/accounts/abi/bind"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/rawdb"
	"github.com/autonity/autonity/core/state"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/eth/filters"
	"github.com/autonity/autonity/eth/gasestimator"
	"github.com/autonity/autonity/ethdb"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/internal/ethapi"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/rpc"
)

type APIBackend interface {
	filters.Backend
	GetEVM(ctx context.Context, state *state.StateDB, header *types.Header, vmConfig *vm.Config, blockCtx *vm.BlockContext) *vm.EVM
}
type TxSender func(signedTx *types.Transaction) error

// This nil assignment ensures at compile time that SimulatedBackend implements bind.ContractBackend.
var _ bind.ContractBackend = (*InternalBackend)(nil)

// InternalBackend implements the contract.Backend interface to interact with the
// protocol contracts. This is used internally by the accountability module and by the autonity cache.
type InternalBackend struct {
	mu           sync.Mutex
	apiBackend   ethapi.Backend
	database     ethdb.Database
	blockchain   *core.BlockChain
	filterSystem *filters.FilterSystem
	config       *params.ChainConfig
	pendingBlock *types.Block
	pendingState *state.StateDB
	TxSender     TxSender
}

var (
	errBlockNumberUnsupported  = errors.New("simulatedBackend cannot access blocks other than the latest block")
	errBlockDoesNotExist       = errors.New("block does not exist in blockchain")
	errTransactionDoesNotExist = errors.New("transaction does not exist")
)

func NewInternalBackend(txSender TxSender, ethAPIBackend ethapi.Backend) func(*core.BlockChain, ethdb.Database) bind.ContractBackend {
	return func(blockchain *core.BlockChain, db ethdb.Database) bind.ContractBackend {
		filterSystem := filters.NewFilterSystem(ethAPIBackend, filters.Config{})
		backend := &InternalBackend{
			database:     db,
			apiBackend:   ethAPIBackend,
			blockchain:   blockchain,
			filterSystem: filterSystem,
			config:       blockchain.Config(),
			TxSender:     txSender,
		}
		return backend
	}
}

// ChainDb returns the backing database
func (b *InternalBackend) ChainDb() ethdb.Database {
	return b.database
}

// HeaderByNumber returns a block header from the current canonical chain.
func (b *InternalBackend) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	if number == nil || number.Cmp(b.blockchain.CurrentHeader().Number) == 0 {
		return b.blockchain.CurrentHeader(), nil
	}
	return b.blockchain.GetHeaderByNumber(number.Uint64()), nil
}

// HeaderByHash returns a block header from the current canonical chain.
func (b *InternalBackend) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	return b.blockchain.GetHeaderByHash(hash), nil
}

// BlockByNumber returns a block from the current canonical chain.
func (b *InternalBackend) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	return b.blockchain.GetBlockByNumber(number.Uint64()), nil
}

// BlockByHash returns a block from the current canonical chain.
func (b *InternalBackend) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	return b.blockchain.GetBlockByHash(hash), nil
}

// PendingBlockAndReceipts returns the pending block and associated receipts.
func (b *InternalBackend) PendingBlockAndReceipts() (*types.Block, types.Receipts) {
	return b.pendingBlock, nil
}

// StateAndHeaderByNumber returns a state and header from the current canonical chain.
func (b *InternalBackend) StateAndHeaderByNumber(ctx context.Context, number *big.Int) (*state.StateDB, *types.Header, error) {
	if number == nil || number.Cmp(b.blockchain.CurrentHeader().Number) == 0 {
		header := b.blockchain.CurrentHeader()
		stateDb, err := b.blockchain.StateAt(header.Root)
		return stateDb, header, err
	}
	header := b.blockchain.GetHeaderByNumber(number.Uint64())
	if header == nil {
		return nil, nil, ethereum.NotFound
	}
	stateDb, err := b.blockchain.StateAt(header.Root)
	return stateDb, header, err
}

// StateAndHeaderByNumberOrHash returns a state and header from the current canonical chain.
func (b *InternalBackend) StateAndHeaderByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*state.StateDB, *types.Header, error) {
	if blockNr, ok := blockNrOrHash.Number(); ok {
		return b.StateAndHeaderByNumber(ctx, big.NewInt(blockNr.Int64()))
	}
	if hash, ok := blockNrOrHash.Hash(); ok {
		header := b.blockchain.GetHeaderByHash(hash)
		if header == nil {
			return nil, nil, ethereum.NotFound
		}
		stateDb, err := b.blockchain.StateAt(header.Root)
		return stateDb, header, err
	}
	return nil, nil, ethereum.NotFound
}

// PendingCodeAt is used for gas estimation but it really doesn't matter, we can use the current state
func (b *InternalBackend) PendingCodeAt(ctx context.Context, contract common.Address) ([]byte, error) {
	return b.CodeAt(ctx, contract, nil)
}

// CodeAt returns the code associated with a certain account in the blockchain.
func (b *InternalBackend) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	statedb, _, err := b.StateAndHeaderByNumber(ctx, blockNumber)
	if err != nil {
		return nil, err
	}
	return statedb.GetCode(contract), nil
}

// BalanceAt returns the balance of the given account at the given block.
func (b *InternalBackend) BalanceAt(ctx context.Context, address common.Address, blockNumber *big.Int) (*big.Int, error) {
	statedb, _, err := b.StateAndHeaderByNumber(ctx, blockNumber)
	if err != nil {
		return nil, err
	}
	return statedb.GetBalance(address).ToBig(), nil
}

// NonceAt returns the nonce of the given account at the given block.
func (b *InternalBackend) NonceAt(ctx context.Context, address common.Address, blockNumber *big.Int) (uint64, error) {
	statedb, _, err := b.StateAndHeaderByNumber(ctx, blockNumber)
	if err != nil {
		return 0, err
	}
	return statedb.GetNonce(address), nil
}

// StorageAt returns the value of key in the storage of an account in the blockchain.
func (b *InternalBackend) StorageAt(ctx context.Context, address common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error) {
	statedb, _, err := b.StateAndHeaderByNumber(ctx, blockNumber)
	if err != nil {
		return nil, err
	}
	val := statedb.GetState(address, key)
	return val[:], nil
}

// TransactionReceipt returns the receipt for a given transaction hash.
func (b *InternalBackend) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	receipt, _, _, _ := rawdb.ReadReceipt(b.database, txHash, b.config)
	if receipt == nil {
		return nil, ethereum.NotFound
	}
	return receipt, nil
}

// PendingNonceAt returns the nonce of the given account in the pending state.
func (b *InternalBackend) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	statedb, err := b.blockchain.State()
	if err != nil {
		return 0, err
	}
	return statedb.GetNonce(account), nil
}

// SuggestGasPrice returns a gas price for a new transaction.
func (b *InternalBackend) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return big.NewInt(1), nil
}

// SuggestGasTipCap returns a gas tip cap for a new transaction.
func (b *InternalBackend) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	return big.NewInt(0), nil
}

// EstimateGas returns an estimate of the gas needed for a transaction.
func (b *InternalBackend) EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error) {
	header := b.blockchain.CurrentHeader()
	stateDb, _ := b.blockchain.StateAt(header.Root)
	opts := &gasestimator.Options{
		Config:     b.ChainConfig(),
		Chain:      b.blockchain,
		Header:     header,
		State:      stateDb,
		ErrorRatio: ethapi.EstimateGasErrorRatio,
	}
	coreMsg := &core.Message{
		To:               call.To,
		From:             call.From,
		Value:            call.Value,
		GasLimit:         header.GasLimit,
		GasPrice:         call.GasPrice,
		SkipNonceChecks:  false,
		SkipFromEOACheck: false,
	}
	cost, _, err := gasestimator.Estimate(ctx, coreMsg, opts, header.GasLimit)
	return cost, err
}

// CallContract executes a contract call.
func (b *InternalBackend) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	if blockNumber != nil && blockNumber.Cmp(b.blockchain.CurrentBlock().Number) != 0 {
		return nil, errBlockNumberUnsupported
	}
	statedb, header, err := b.StateAndHeaderByNumber(ctx, blockNumber)
	if err != nil {
		return nil, err
	}
	blockCtx := core.NewEVMBlockContext(header, ethapi.NewChainContext(ctx, b.apiBackend), nil)
	evm := b.apiBackend.GetEVM(ctx, statedb, header, b.blockchain.GetVMConfig(), &blockCtx)
	gp := core.GasPool(call.Gas)
	res, err := core.ApplyMessage(evm, nil, &gp)
	if err != nil {
		return nil, err
	}
	// If the result contains a revert reason, try to unpack and return it.
	if len(res.Revert()) > 0 {
		return nil, ethapi.NewRevertError(res.Revert())
	}
	return res.Return(), res.Err
}

// SendTransaction sends a transaction to the network.
func (b *InternalBackend) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return b.TxSender(tx)
}

// SubscribeFilterLogs creates a subscription that will write all logs matching the
// given criteria to the given logs channel.
func (b *InternalBackend) SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	from, to := int64(0), int64(rpc.LatestBlockNumber)
	if query.FromBlock != nil {
		from = query.FromBlock.Int64()
	}
	if query.ToBlock != nil {
		to = query.ToBlock.Int64()
	}
	// Create a filter for the given criteria
	filter := b.filterSystem.NewRangeFilter(from, to, query.Addresses, query.Topics)

	// Create a subscription that forwards logs to the channel
	subscription := event.NewSubscription(func(quit <-chan struct{}) error {
		for {
			select {
			case <-quit:
				return nil
			default:
				logs, err := filter.Logs(ctx)
				if err != nil {
					return err
				}
				for _, log := range logs {
					select {
					case ch <- *log:
					case <-quit:
						return nil
					}
				}
				// Sleep briefly before checking for new logs
				time.Sleep(time.Second)
			}
		}
	})

	return subscription, nil
}

// FilterLogs executes a log filter operation, blocking during execution and
// returning all the results in one batch.
func (b *InternalBackend) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	// Create a filter for the given criteria
	filter := b.filterSystem.NewRangeFilter(query.FromBlock.Int64(), query.ToBlock.Int64(), query.Addresses, query.Topics)

	// Get logs from the filter
	logs, err := filter.Logs(ctx)
	if err != nil {
		return nil, err
	}

	// Convert from []*types.Log to []types.Log
	result := make([]types.Log, len(logs))
	for i, log := range logs {
		result[i] = *log
	}

	return result, nil
}

// SubscribeNewHead subscribes to notifications about changes of the head block of
// the canonical chain.
func (b *InternalBackend) SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error) {
	sub := event.NewSubscription(func(quit <-chan struct{}) error {
		return nil
	})
	return sub, nil
}

// ChainConfig returns the chain configuration.
func (b *InternalBackend) ChainConfig() *params.ChainConfig {
	return b.config
}

// CurrentBlock returns the current block.
func (b *InternalBackend) CurrentBlock() *types.Header {
	return b.blockchain.CurrentHeader()
}
