package core

import (
	"testing"
	"time"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/core/types"
)

type VerifySelfProposalAlwaysTrueBackend struct {
	*testing.T
	Backend
	marker
}

type marker interface {
	markMalicious(types.Block)
}

func (b *VerifySelfProposalAlwaysTrueBackend) VerifyProposal(block types.Block) (time.Duration, error) {
	if block.Number().Uint64() < 2 {
		// skip genesis and the first block
		return b.Backend.VerifyProposal(block)
	}

	b.marker.markMalicious(block)

	return 0, nil
}

func (b *VerifySelfProposalAlwaysTrueBackend) VerifyHeader(chain consensus.ChainReader, header *types.Header, seal bool) error {
	if header.Number.Uint64() < 2 {
		// skip genesis and the first block
		return b.Backend.VerifyHeader(chain, header, seal)
	}
	return nil
}

func (b *VerifySelfProposalAlwaysTrueBackend) GetOriginal() Backend {
	return b.Backend
}

type changeValidatorBackend struct {
	*testing.T
	blocksFromRemoved map[common.Hash]uint64
	blocksFromAdded   map[common.Hash]uint64
	brokenValidators  Changes
}

func NewChangeValidatorBackend(t *testing.T, basicCore Backend, brokenValidators Changes, blocksFromRemoved, blocksFromAdded map[common.Hash]uint64) *VerifySelfProposalAlwaysTrueBackend {
	return &VerifySelfProposalAlwaysTrueBackend{
		T:       t,
		Backend: basicCore,
		marker: changeValidatorBackend{
			T:                 t,
			blocksFromRemoved: blocksFromRemoved,
			blocksFromAdded:   blocksFromAdded,
			brokenValidators:  brokenValidators,
		},
	}
}

func (o changeValidatorBackend) markMalicious(block types.Block) {
	if _, ok := o.brokenValidators.added[block.Coinbase()]; ok {
		o.blocksFromAdded[block.Hash()] = block.Number().Uint64()
	}

	if _, ok := o.brokenValidators.removed[block.Coinbase()]; ok {
		o.blocksFromRemoved[block.Hash()] = block.Number().Uint64()
	}
}
