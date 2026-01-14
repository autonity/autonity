package tests

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"

	autonity "github.com/autonity/autonity"
	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/accounts/abi/bind"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/core/vm"
)

type RunnerBackend struct {
	Runner   *Runner
	receipts map[common.Hash]*types.Receipt
	mu       sync.Mutex
}

func NewRunnerBackend(r *Runner) *RunnerBackend {
	return &RunnerBackend{
		Runner:   r,
		receipts: make(map[common.Hash]*types.Receipt),
	}
}

var _ bind.ContractBackend = (*RunnerBackend)(nil)

func (b *RunnerBackend) CallContract(ctx context.Context, call autonity.CallMsg, blockNumber *big.Int) ([]byte, error) {
	if call.From == (common.Address{}) {
		call.From = User
	}

	snapshot := b.Runner.Evm.StateDB.Snapshot()
	defer b.Runner.Evm.StateDB.RevertToSnapshot(snapshot)

	b.Runner.Evm.TxContext.Origin = call.From

	result, _, err := b.Runner.Evm.StaticCall(
		vm.AccountRef(call.From),
		*call.To,
		call.Data,
		uint64(1<<63),
	)
	return result, err
}

func (b *RunnerBackend) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	signer := types.LatestSignerForChainID(b.Runner.Evm.ChainConfig().ChainID)
	from, err := types.Sender(signer, tx)
	if err != nil {
		return err
	}

	b.Runner.Evm.TxContext.Origin = from
	b.Runner.Evm.TxContext.GasPrice = tx.GasPrice()
	gas := tx.Gas()
	value := tx.Value()
	if value == nil {
		value = common.Big0
	}

	var (
		ret          []byte
		leftOverGas  uint64
		contractAddr common.Address
		vmerr        error
	)

	if tx.To() == nil {
		ret, contractAddr, leftOverGas, vmerr = b.Runner.Evm.Create(vm.AccountRef(from), tx.Data(), gas, value)
	} else {
		ret, leftOverGas, vmerr = b.Runner.Evm.Call(vm.AccountRef(from), *tx.To(), tx.Data(), gas, value)
	}

	// auto increment nonce for the caller
	currentNonce := b.Runner.Evm.StateDB.GetNonce(from)
	b.Runner.Evm.StateDB.SetNonce(from, currentNonce+1)
	status := uint64(1)
	if vmerr != nil {
		status = 0
		fmt.Printf("\n TX FAILED: %v\n", vmerr)

		if len(ret) > 0 {
			if reason, err := abi.UnpackRevert(ret); err == nil {
				fmt.Printf("   REASON: %s\n", reason)
			} else {
				fmt.Printf("   RAW REVERT DATA: %x\n", ret)
			}
		}
	}

	// 4. Save Receipt
	receipt := &types.Receipt{
		Type:            tx.Type(),
		Status:          status,
		TxHash:          tx.Hash(),
		ContractAddress: contractAddr,
		GasUsed:         gas - leftOverGas,
		BlockNumber:     b.Runner.Evm.Context.BlockNumber,
	}
	b.receipts[tx.Hash()] = receipt

	return vmerr
}

// TransactionReceipt returns the stored receipt.
func (b *RunnerBackend) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	receipt, ok := b.receipts[txHash]
	if !ok {
		return nil, autonity.NotFound
	}
	return receipt, nil
}

func (b *RunnerBackend) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	return b.Runner.Evm.StateDB.GetCode(contract), nil
}

func (b *RunnerBackend) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	return b.Runner.Evm.StateDB.GetCode(account), nil
}

func (b *RunnerBackend) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	if number != nil && number.Cmp(b.Runner.Evm.Context.BlockNumber) != 0 {
		return nil, errors.New("historical headers not supported")
	}
	return &types.Header{
		Number:  b.Runner.Evm.Context.BlockNumber,
		Time:    uint64(b.Runner.Evm.Context.Time.Int64()),
		BaseFee: b.Runner.Evm.Context.BaseFee,
	}, nil
}

func (b *RunnerBackend) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return common.Big0, nil
}
func (b *RunnerBackend) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	return big.NewInt(1_000_000_000), nil
}

func (b *RunnerBackend) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	return b.Runner.Evm.StateDB.GetNonce(account), nil
}

func (b *RunnerBackend) EstimateGas(ctx context.Context, call autonity.CallMsg) (gas uint64, err error) {
	snapshot := b.Runner.snapshot()
	defer b.Runner.Evm.StateDB.RevertToSnapshot(snapshot)
	if call.To == nil {
		return 10_000_000, nil
	}
	_, usedGas, err := b.Runner.call(
		&runOptions{origin: call.From, value: call.Value},
		*call.To,
		call.Data)

	return usedGas, err
}

func (b *RunnerBackend) FilterLogs(ctx context.Context, query autonity.FilterQuery) ([]types.Log, error) {
	return make([]types.Log, 0), nil
}

func (b *RunnerBackend) SubscribeFilterLogs(ctx context.Context, query autonity.FilterQuery, ch chan<- types.Log) (autonity.Subscription, error) {
	return nil, nil
}
