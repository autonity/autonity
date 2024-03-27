package backend

import (
	"github.com/autonity/autonity/trie"
	"math/big"
	"testing"
	"time"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/core/state"
	"github.com/autonity/autonity/core/types"
)

type ModifyCommitteeEngine struct {
	*testing.T
	*Backend
	Modifier
	marker
}

type Modifier interface {
	ModifyHeader(header *types.Header) *types.Header
}

type marker interface {
	markMalicious(*types.Block)
}

func (m *ModifyCommitteeEngine) VerifyProposal(block *types.Block) (time.Duration, error) {
	if block.Number().Uint64() < 2 {
		// skip genesis and the first block
		return m.Backend.VerifyProposal(block)
	}

	m.marker.markMalicious(block)

	return 0, nil
}

func (m *ModifyCommitteeEngine) VerifyHeader(chain consensus.ChainHeaderReader, header *types.Header, seal bool) error {
	if header.Number.Uint64() < 2 {
		// skip genesis and the first block
		return m.Backend.VerifyHeader(chain, header, seal)
	}
	return nil
}

func (m *ModifyCommitteeEngine) FinalizeAndAssemble(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header, receipts *[]*types.Receipt) (*types.Block, error) {
	// create a normal block and check for errors
	block, err := m.Backend.FinalizeAndAssemble(chain, header, state, txs, uncles, receipts)
	if err != nil {
		m.T.Error("m.core.FinalizeAndAssemble returned error:", err, "Expected nil")
	}

	if header.Number.Cmp(big.NewInt(0)) == 0 {
		// skip genesis
		return block, nil
	}

	if m.address != header.Coinbase {
		// if the malicious validator is a proposer
		return block, nil
	}

	lastMinedBlock := m.Backend.HeadBlock()
	if lastMinedBlock.Number().Cmp(header.Number) != 0 {
		return block, nil
	}

	header = m.Modifier.ModifyHeader(block.Header())

	// create a new block with the modified header
	newBlock := types.NewBlock(header, block.Transactions(), block.Uncles(), *receipts, new(trie.Trie))

	newBlock, err = m.Backend.AddSeal(newBlock)
	if err != nil {
		m.Errorf("cant seal the block: %v - %v", err, newBlock)
	}

	// we want be sure that the block is modified but not broken
	switch _, err = m.Backend.VerifyProposal(newBlock); err {
	case consensus.ErrInconsistentCommitteeSet:
	// nothing to do
	default:
		m.Error("Mock FinalizeAndAssemble created incorrect block:", err, newBlock)
	}

	return newBlock, nil
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

type changeValidator struct {
	blocksFromRemoved map[common.Hash]uint64
	blocksFromAdded   map[common.Hash]uint64
	brokenValidators  Changes
}

func newChangeValidator(blocksFromRemoved map[common.Hash]uint64, blocksFromAdded map[common.Hash]uint64, changedValidators Changes) changeValidator {
	return changeValidator{
		blocksFromRemoved: blocksFromRemoved,
		blocksFromAdded:   blocksFromAdded,
		brokenValidators:  changedValidators,
	}
}

func (o changeValidator) markMalicious(block *types.Block) {
	if _, ok := o.brokenValidators.added[block.Coinbase()]; ok {
		o.blocksFromAdded[block.Hash()] = block.Number().Uint64()
	}

	if _, ok := o.brokenValidators.removed[block.Coinbase()]; ok {
		o.blocksFromRemoved[block.Hash()] = block.Number().Uint64()
	}
}

type addValidator Changes

func NewAddValidator(t *testing.T, engine consensus.Engine, changedValidators Changes, blocksFromRemoved, blocksFromAdded map[common.Hash]uint64) *ModifyCommitteeEngine {
	basicEngine, ok := engine.(*Backend)
	if !ok {
		panic("*Backend type is expected")
	}
	return &ModifyCommitteeEngine{
		T:        t,
		Backend:  basicEngine,
		Modifier: addValidator(changedValidators),
		marker:   newChangeValidator(blocksFromRemoved, blocksFromAdded, changedValidators),
	}
}

func (p addValidator) ModifyHeader(header *types.Header) *types.Header {
	additionalValidator := &types.CommitteeMember{
		Address:     common.Address{3},
		VotingPower: new(big.Int).SetUint64(1),
	}
	p.added[additionalValidator.Address] = struct{}{}

	header.Committee.Members = append(header.Committee.Members, additionalValidator)

	return header
}

type removeValidator Changes

func NewRemoveValidator(t *testing.T, engine consensus.Engine, changedValidators Changes, blocksFromRemoved, blocksFromAdded map[common.Hash]uint64) *ModifyCommitteeEngine {
	basicEngine, ok := engine.(*Backend)
	if !ok {
		panic("*Backend type is expected")
	}
	return &ModifyCommitteeEngine{
		T:        t,
		Backend:  basicEngine,
		Modifier: removeValidator(changedValidators),
		marker:   newChangeValidator(blocksFromRemoved, blocksFromAdded, changedValidators),
	}
}

func (p removeValidator) ModifyHeader(header *types.Header) *types.Header {
	p.removed[header.Committee.Members[len(header.Committee.Members)-1].Address] = struct{}{}
	header.Committee.Members = header.Committee.Members[:len(header.Committee.Members)-1]
	return header
}

type replaceValidator Changes

func NewReplaceValidator(t *testing.T, engine consensus.Engine, changedValidators Changes, blocksFromRemoved, blocksFromAdded map[common.Hash]uint64) *ModifyCommitteeEngine {
	basicEngine, ok := engine.(*Backend)
	if !ok {
		panic("*Backend type is expected")
	}
	return &ModifyCommitteeEngine{
		T:        t,
		Backend:  basicEngine,
		Modifier: replaceValidator(changedValidators),
		marker:   newChangeValidator(blocksFromRemoved, blocksFromAdded, changedValidators),
	}
}

func (p replaceValidator) ModifyHeader(header *types.Header) *types.Header {
	maliciousValidator := types.CommitteeMember{
		Address:     common.Address{3},
		VotingPower: new(big.Int).SetUint64(1),
	}
	p.added[maliciousValidator.Address] = struct{}{}
	p.removed[header.Committee.Members[len(header.Committee.Members)-1].Address] = struct{}{}

	header.Committee.Members = append(header.Committee.Members[:len(header.Committee.Members)-1], &maliciousValidator)

	return header
}
