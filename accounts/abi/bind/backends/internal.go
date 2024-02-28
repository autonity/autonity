package backends

import (
	"context"

	ethereum "github.com/autonity/autonity"
	"github.com/autonity/autonity/accounts/abi/bind"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/eth/filters"
	"github.com/autonity/autonity/ethdb"
)

// InternalBackend implements the contract.Backend interface to interact with the
// protocol contracts. This is used internally by the accountability module and by the autonity cache.
type InternalBackend struct {
	SimulatedBackend
	TxSender *func(signedTx *types.Transaction) error
}

func NewInternalBackend(txSender *func(signedTx *types.Transaction) error) func(*core.BlockChain, ethdb.Database) bind.ContractBackend {
	return func(blockchain *core.BlockChain, db ethdb.Database) bind.ContractBackend {
		backend := &InternalBackend{
			SimulatedBackend{
				database:   db,
				blockchain: blockchain,
				events:     filters.NewEventSystem(&filterBackend{db, blockchain}, false),
				config:     blockchain.Config(),
			}, txSender}
		return backend
	}
}

// PendingCodeAt is used for gas estimation but it really doesn't matter, we can use the current state
func (b *InternalBackend) PendingCodeAt(ctx context.Context, contract common.Address) ([]byte, error) {
	return b.CodeAt(ctx, contract, nil)
}

func (b *InternalBackend) EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error) {
	// Increase the current block number by 1, we want to simulate the tx as if it will be included in the next block
	// Otherwise precompiles which check blockNumber timing conditions (e.g. accusation verifier) might wrongly fail
	// NOTE: Theoretically also header.GasLimit should be changed since GasLimit adapts
	// based on parent GasLimit and desired GasCeil. However it is not a big deal.
	block := types.NewBlockWithHeader(b.blockchain.CurrentHeader())
	height := block.Number().Add(block.Number(), common.Big1)
	block.SetHeaderNumber(height)

	// very hacky but that will avoid us further complications
	b.SimulatedBackend.pendingBlock = block
	b.SimulatedBackend.pendingState, _ = b.blockchain.State()

	return b.SimulatedBackend.EstimateGas(ctx, call)
}

func (b *InternalBackend) SendTransaction(_ context.Context, tx *types.Transaction) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	return (*b.TxSender)(tx)
}
