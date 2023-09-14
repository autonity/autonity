package coretest

import (
	"context"
	"math/big"

	"github.com/autonity/autonity/accounts/abi/bind"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/ethdb"
	"github.com/ethereum/go-ethereum"
)

type FakeContractBackend struct{}

func FakeContractBackendProvider(_ *BlockChain, _ ethdb.Database) bind.ContractBackend {
	return new(FakeContractBackend)
}
func (*FakeContractBackend) CodeAt(_ context.Context, _ common.Address, _ *big.Int) ([]byte, error) {
	return make([]byte, 0), nil
}
func (*FakeContractBackend) CallContract(_ context.Context, _ ethereum.CallMsg, _ *big.Int) ([]byte, error) {
	return make([]byte, 0), nil
}

func (*FakeContractBackend) HeaderByNumber(_ context.Context, _ *big.Int) (*types.Header, error) {
	return &types.Header{}, nil
}

func (*FakeContractBackend) PendingCodeAt(_ context.Context, _ common.Address) ([]byte, error) {
	return make([]byte, 0), nil
}

func (*FakeContractBackend) PendingNonceAt(_ context.Context, _ common.Address) (uint64, error) {
	return 0, nil
}

func (*FakeContractBackend) SuggestGasPrice(_ context.Context) (*big.Int, error) {
	return new(big.Int), nil
}

func (*FakeContractBackend) SuggestGasTipCap(_ context.Context) (*big.Int, error) {
	return new(big.Int), nil
}

func (*FakeContractBackend) EstimateGas(_ context.Context, _ ethereum.CallMsg) (gas uint64, err error) {
	return 0, nil
}

func (*FakeContractBackend) SendTransaction(_ context.Context, _ *types.Transaction) error {
	return nil
}

func (*FakeContractBackend) FilterLogs(_ context.Context, _ ethereum.FilterQuery) ([]types.Log, error) {
	return make([]types.Log, 0), nil
}

func (*FakeContractBackend) SubscribeFilterLogs(_ context.Context, _ ethereum.FilterQuery, _ chan<- types.Log) (ethereum.Subscription, error) {
	return new(FakeSubscription), nil
}

type FakeSubscription struct{}

func (*FakeSubscription) Unsubscribe() {}

func (*FakeSubscription) Err() <-chan error {
	return make(chan error)
}
