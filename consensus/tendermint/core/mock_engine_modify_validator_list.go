package core

import (
	"math/big"
	"testing"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/core/state"
	"github.com/clearmatics/autonity/core/types"
)

type Modifier interface {
	ModifyHeader(header *types.Header) *types.Header
}

type ModifyCommitteeEngine struct {
	*testing.T
	*core
	Modifier
}

func (*ModifyCommitteeEngine) VerifyHeader(_ consensus.ChainReader, _ *types.Header, _ bool) error {
	return nil
}

func (c *ModifyCommitteeEngine) FinalizeAndAssemble(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header, receipts *[]*types.Receipt) (*types.Block, error) {
	// create a normal block and check for errors
	block, err := c.core.backend.FinalizeAndAssemble(chain, header, state, txs, uncles, receipts)
	if err != nil {
		c.T.Error("c.core.FinalizeAndAssemble returned error:", err, "Expected nil")
	}

	if header.Number.Cmp(big.NewInt(0)) == 0 {
		// skip genesis
		return block, nil
	}

	if c.address != header.Coinbase {
		// if the malicious validator is a proposer
		return block, nil
	}

	if c.curRoundMessages == nil || c.curRoundMessages != nil || c.Height().Cmp(header.Number) != 0 {
		return block, nil
	}

	header = c.Modifier.ModifyHeader(block.Header())

	// create a new block with the modified header
	newBlock := types.NewBlock(header, block.Transactions(), block.Uncles(), *receipts)

	newBlock, err = c.backend.AddSeal(newBlock)
	if err != nil {
		c.Errorf("cant seal the block: %v - %v", err, newBlock)
	}

	type originalBackend interface {
		GetOriginal() Backend
	}

	back := c.core.backend
	original, ok := c.core.backend.(originalBackend)
	if ok {
		back = original.GetOriginal()
	}

	// we want be sure that the block is modified but not broken
	switch _, err = back.VerifyProposal(*newBlock); err {
	case consensus.ErrInconsistentCommitteeSet:
	// nothing to do
	default:
		c.Error("Mock FinalizeAndAssemble created incorrect block:", err, newBlock)
	}

	return newBlock, nil
}

type addValidatorCore Changes

func NewAddValidatorCore(c consensus.Engine, changedValidators Changes) *ModifyCommitteeEngine {
	basicCore, ok := c.(*Core)
	if !ok {
		panic("*core type is expected")
	}
	return &ModifyCommitteeEngine{
		core:     basicCore,
		Modifier: addValidatorCore(changedValidators),
	}
}

func (p addValidatorCore) ModifyHeader(header *types.Header) *types.Header {
	additionalValidator := types.CommitteeMember{
		Address:     common.Address{3},
		VotingPower: new(big.Int).SetUint64(1),
	}
	p.added[additionalValidator.Address] = struct{}{}

	header.Committee = append(header.Committee, additionalValidator)

	return header
}

type removeValidatorCore Changes

func NewRemoveValidatorCore(c consensus.Engine, changedValidators Changes) *ModifyCommitteeEngine {
	basicCore, ok := c.(*core)
	if !ok {
		panic("*core type is expected")
	}
	return &ModifyCommitteeEngine{
		core:     basicCore,
		Modifier: removeValidatorCore(changedValidators),
	}
}

func (p removeValidatorCore) ModifyHeader(header *types.Header) *types.Header {
	p.removed[header.Committee[len(header.Committee)-1].Address] = struct{}{}
	header.Committee = header.Committee[:len(header.Committee)-1]
	return header
}

type Changes struct {
	added   map[common.Address]struct{}
	removed map[common.Address]struct{}
}

func NewChanges() Changes {
	return Changes{
		added:   make(map[common.Address]struct{}),
		removed: make(map[common.Address]struct{}),
	}
}

type replaceValidatorCore Changes

func NewReplaceValidatorCore(c consensus.Engine, changedValidators Changes) *ModifyCommitteeEngine {
	basicCore, ok := c.(*core)
	if !ok {
		panic("*core type is expected")
	}
	return &ModifyCommitteeEngine{
		core:     basicCore,
		Modifier: replaceValidatorCore(changedValidators),
	}
}

func (p replaceValidatorCore) ModifyHeader(header *types.Header) *types.Header {
	maliciousValidator := types.CommitteeMember{
		Address:     common.Address{3},
		VotingPower: new(big.Int).SetUint64(1),
	}
	p.added[maliciousValidator.Address] = struct{}{}
	p.removed[header.Committee[len(header.Committee)-1].Address] = struct{}{}

	header.Committee = append(header.Committee[:len(header.Committee)-1], maliciousValidator)

	return header
}
