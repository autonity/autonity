package tests

import (
	"context"
	"math/big"

	autonity "github.com/autonity/autonity"
	"github.com/autonity/autonity/accounts/abi/bind"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/core/vm"
)

type RunnerBackend struct {
	Runner *Runner
}

var _ bind.ContractBackend = (*RunnerBackend)(nil)

func (b *RunnerBackend) CallContract(ctx context.Context, call autonity.CallMsg, blockNumber *big.Int) ([]byte, error) {
	if call.From == (common.Address{}) {
		call.From = User
	}

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
	signer := types.LatestSignerForChainID(b.Runner.Evm.ChainConfig().ChainID)
	sender, err := types.Sender(signer, tx)
	if err != nil {
		return err
	}

	_, _, err = b.Runner.call(
		&runOptions{origin: sender, value: tx.Value()},
		*tx.To(),
		tx.Data(),
	)
	return err
}

func (b *RunnerBackend) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	return b.Runner.Evm.StateDB.GetCode(contract), nil
}

func (b *RunnerBackend) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	return b.Runner.Evm.StateDB.GetCode(account), nil
}

func (b *RunnerBackend) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	return &types.Header{
		Number: b.Runner.Evm.Context.BlockNumber,
		Time:   uint64(b.Runner.Evm.Context.Time.Int64()),
	}, nil
}

func (b *RunnerBackend) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return common.Big0, nil
}
func (b *RunnerBackend) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	return common.Big0, nil
}

func (b *RunnerBackend) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	return 0, nil
}

func (b *RunnerBackend) EstimateGas(ctx context.Context, call autonity.CallMsg) (gas uint64, err error) {
	return 0, nil
}

func (b *RunnerBackend) FilterLogs(ctx context.Context, query autonity.FilterQuery) ([]types.Log, error) {
	return make([]types.Log, 0), nil
}

func (b *RunnerBackend) SubscribeFilterLogs(ctx context.Context, query autonity.FilterQuery, ch chan<- types.Log) (autonity.Subscription, error) {
	return nil, nil
}
