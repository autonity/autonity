package backends

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	ethereum "github.com/autonity/autonity"
	"github.com/autonity/autonity/accounts/abi/bind"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/state"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/eth/filters"
	"github.com/autonity/autonity/ethdb"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/rpc"
)

// InternalBackend implements the contract.Backend interface to interact with the
// protocol contracts. This is used internally by the accountability module and by the autonity cache.
type InternalBackend struct {
	mu           sync.Mutex
	database     ethdb.Database
	blockchain   *core.BlockChain
	filterSystem *filters.FilterSystem
	config       *params.ChainConfig
	pendingBlock *types.Block
	pendingState *state.StateDB
	TxSender     func(signedTx *types.Transaction) error
}

func NewInternalBackend(txSender func(signedTx *types.Transaction) error) func(*core.BlockChain, ethdb.Database) bind.ContractBackend {
	return func(blockchain *core.BlockChain, db ethdb.Database) bind.ContractBackend {
		filterSystem := filters.NewFilterSystem(blockchain, filters.Config{})

		backend := &InternalBackend{
			database:     db,
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
	if number == nil || number.Cmp(b.blockchain.CurrentHeader().Number) == 0 {
		return b.blockchain.CurrentBlock(), nil
	}
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
	return statedb.GetBalance(address), nil
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
	receipt, _, _, _ := core.GetReceipt(b.database, txHash)
	return receipt, nil
}

// TransactionByHash returns the transaction for a given hash.
func (b *InternalBackend) TransactionByHash(ctx context.Context, txHash common.Hash) (*types.Transaction, common.Hash, uint64, uint64, error) {
	tx, blockHash, blockNumber, index := core.GetTransaction(b.database, txHash)
	return tx, blockHash, blockNumber, index, nil
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
	b.mu.Lock()
	defer b.mu.Unlock()

	// Increase the current block number by 1, we want to simulate the tx as if it will be included in the next block
	// Otherwise precompiles which check blockNumber timing conditions (e.g. accusation verifier) might wrongly fail
	block := types.NewBlockWithHeader(b.blockchain.CurrentHeader())
	height := block.Number().Add(block.Number(), common.Big1)
	block.SetHeaderNumber(height)

	b.pendingBlock = block
	b.pendingState, _ = b.blockchain.State()

	// Create a helper to check if a gas allowance results in an executable transaction
	executable := func(gas uint64) (bool, *core.ExecutionResult, error) {
		call.Gas = gas

		currentState, err := b.blockchain.State()
		if err != nil {
			return false, nil, err
		}

		// Create a new environment which holds all relevant information
		// about the transaction and calling mechanisms.
		msg := core.Message{
			From:       call.From,
			To:         call.To,
			Value:      call.Value,
			Data:       call.Data,
			GasLimit:   call.Gas,
			GasPrice:   call.GasPrice,
			GasFeeCap:  call.GasFeeCap,
			GasTipCap:  call.GasTipCap,
			AccessList: call.AccessList,
		}

		// Setup context so it may be cancelled the call has completed
		// or, in case of unmetered gas, setup a context with a timeout.
		var cancel context.CancelFunc
		if ctx == nil {
			ctx, cancel = context.WithTimeout(context.Background(), vm.Timeout)
		} else {
			ctx, cancel = context.WithTimeout(ctx, vm.Timeout)
		}
		defer cancel()

		// Get a new instance of the EVM.
		evm, vmError, err := b.blockchain.GetEVM(ctx, msg, currentState, block.Header())
		if err != nil {
			return false, nil, err
		}

		// Execute the message.
		res, err := core.ApplyMessage(evm, msg, new(core.GasPool).AddGas(call.Gas))
		if err != nil {
			return false, nil, err
		}
		return res.Failed(), res, nil
	}

	// Execute the binary search and hone in on an executable gas limit
	lo := params.TxGas - 1
	hi := call.Gas
	if hi == 0 {
		hi = b.blockchain.CurrentHeader().GasLimit
	}
	cap := hi

	// Binary search the gas limit
	for lo+1 < hi {
		mid := (hi + lo) / 2
		failed, _, err := executable(mid)
		if err != nil {
			return 0, err
		}
		if failed {
			lo = mid
		} else {
			hi = mid
		}
	}

	// If the transaction still failed with the highest gas limit, return the error
	if hi == cap {
		failed, result, err := executable(hi)
		if err != nil {
			return 0, err
		}
		if failed {
			if result != nil && result.Err != vm.ErrOutOfGas {
				return 0, result.Err
			}
			return 0, fmt.Errorf("gas required exceeds allowance (%d)", cap)
		}
	}
	return hi, nil
}

// SendTransaction sends a transaction to the network.
func (b *InternalBackend) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	return (*b.TxSender)(tx)
}

// SubscribeFilterLogs creates a subscription that will write all logs matching the
// given criteria to the given logs channel.
func (b *InternalBackend) SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	// Create a filter for the given criteria
	filter := b.filterSystem.NewLogFilter(filters.FilterCriteria{
		FromBlock: query.FromBlock,
		ToBlock:   query.ToBlock,
		Addresses: query.Addresses,
		Topics:    query.Topics,
	}, nil)

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
	filter := b.filterSystem.NewLogFilter(filters.FilterCriteria{
		FromBlock: query.FromBlock,
		ToBlock:   query.ToBlock,
		Addresses: query.Addresses,
		Topics:    query.Topics,
	}, nil)

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
