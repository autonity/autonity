package tendermint

import (
	common "github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	autonity "github.com/clearmatics/autonity/contracts/autonity"
	"github.com/clearmatics/autonity/core/state"
	types "github.com/clearmatics/autonity/core/types"
)

type Finalizer interface {

	// Finalize runs any post-transaction state modifications (e.g. block rewards)
	// but does not assemble the block.
	//
	// Note: The block header and state database might be updated to reflect any
	// consensus rules that happen at finalization (e.g. block rewards).
	Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
		uncles []*types.Header, receipts []*types.Receipt) (types.Committee, *types.Receipt, error)

	// FinalizeAndAssemble runs any post-transaction state modifications (e.g. block
	// rewards) and assembles the final block.
	//
	// Note: The block header and state database might be updated to reflect any
	// consensus rules that happen at finalization (e.g. block rewards).
	FinalizeAndAssemble(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
		uncles []*types.Header, receipts *[]*types.Receipt) (*types.Block, error)
}

type DefaultFinalizer struct {
	autonityContract *autonity.Contract
}

func NewFinalizer(ac *autonity.Contract) *DefaultFinalizer {
	return &DefaultFinalizer{
		autonityContract: ac,
	}
}

// Finalize runs any post-transaction state modifications (e.g. block rewards)
// Finaize doesn't modify the passed header.
func (f *DefaultFinalizer) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header, receipts []*types.Receipt) (types.Committee, *types.Receipt, error) {

	// AutonityContractFinalize is called to deploy the Autonity Contract at block #1. it returns as well the
	// committee field containaining the list of committee members allowed to participate in consensus for the next block.
	return f.autonityContract.FinalizeAndGetCommittee(txs, receipts, header, state)
}

// FinalizeAndAssemble call Finaize to compute post transacation state modifications
// and assembles the final block.
func (f *DefaultFinalizer) FinalizeAndAssemble(chain consensus.ChainReader, header *types.Header, statedb *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header, receipts *[]*types.Receipt) (*types.Block, error) {

	statedb.Prepare(common.ACHash(header.Number), common.Hash{}, len(txs))
	committeeSet, receipt, err := f.Finalize(chain, header, statedb, txs, uncles, *receipts)
	if err != nil {
		return nil, err
	}
	*receipts = append(*receipts, receipt)
	// No block rewards in BFT, so the state remains as is and uncles are dropped
	header.Root = statedb.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	header.UncleHash = nilUncleHash

	// add committee to extraData's committee section
	header.Committee = committeeSet
	return types.NewBlock(header, txs, nil, *receipts), nil
}
