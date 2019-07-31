package core

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/core/state"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/rpc"
	"math/big"
)

func (c *core) Author(header *types.Header) (common.Address, error) {
	return c.backend.Author(header)
}

func (c *core) VerifyHeader(chain consensus.ChainReader, header *types.Header, seal bool) error {
	return c.backend.VerifyHeader(chain, header, seal)
}

func (c *core) VerifyHeaders(chain consensus.ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	return c.backend.VerifyHeaders(chain, headers, seals)
}

func (c *core) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	return c.backend.VerifyUncles(chain, block)
}

func (c *core) VerifySeal(chain consensus.ChainReader, header *types.Header) error {
	return c.backend.VerifySeal(chain, header)
}

func (c *core) Prepare(chain consensus.ChainReader, header *types.Header) error {
	return c.backend.Prepare(chain, header)
}

func (c *core) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {
	return c.backend.Finalize(chain, header, state, txs, uncles, receipts)
}

func (c *core) Seal(chain consensus.ChainReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {
	return c.backend.Seal(chain, block, results, stop)
}

func (c *core) SealHash(header *types.Header) common.Hash {
	return c.backend.SealHash(header)
}

func (c *core) CalcDifficulty(chain consensus.ChainReader, time uint64, parent *types.Header) *big.Int {
	return c.backend.CalcDifficulty(chain, time, parent)
}

func (c *core) APIs(chain consensus.ChainReader) []rpc.API {
	return c.backend.APIs(chain)
}

func (c *core) Close() error {
	return c.Stop()
}
