// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package backend

import (
	"context"
	"math/big"

	"github.com/clearmatics/autonity/core"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/rpc"
)

// Author retrieves the Ethereum address of the account that minted the given
// block, which may be different from the header's coinbase if a consensus
// engine is based on signatures.
func (sb *Backend) Author(header *types.Header) (common.Address, error) {
	return sb.core.Author(header)
}

// Prepare initializes the consensus fields of a block header according to the
// rules of a particular engine. The changes are executed inline.
func (sb *Backend) Prepare(chain consensus.ChainReader, header *types.Header) error {
	return sb.core.Prepare(chain, header)
}

//
// So this method is meant to allow interrupting of mining a block to start on
// a new block, it doesn't make sense for autonity though because if we are not
// the proposer then we don't need this unsigned block, and if we are the
// proposer we only want the one unsigned block per round since we can't send
// multiple differing proposals.
//
// So we want to have just the latest block available to be taken from here when this node becomes the proposer.
func (sb *Backend) Seal(chain consensus.ChainReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {
	return sb.core.Seal(chain, block, results, stop)
}

// CalcDifficulty is the difficulty adjustment algorithm. It returns the difficulty
// that a new block should have based on the previous blocks in the blockchain and the
// current signer.
func (sb *Backend) CalcDifficulty(chain consensus.ChainReader, time uint64, parent *types.Header) *big.Int {
	return sb.core.CalcDifficulty(chain, time, parent)
}

// APIs returns the RPC APIs this consensus engine provides.
func (sb *Backend) APIs(chain consensus.ChainReader) []rpc.API {
	return sb.core.APIs(chain)
}

// Start implements consensus.Start
func (sb *Backend) Start(ctx context.Context, blockchain *core.BlockChain) error {
	return sb.core.Start(ctx, blockchain)
}

// Stop implements consensus.Stop
func (sb *Backend) Close() error {
	return sb.core.Close()
}

func (sb *Backend) SealHash(header *types.Header) common.Hash {
	return types.SigHash(header)
}
