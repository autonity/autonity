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
	"errors"
	"math/big"
	"time"

	"github.com/clearmatics/autonity/core"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/rpc"
)

const (
	inmemorySnapshots = 128 // Number of recent vote snapshots to keep in memory
	inmemoryPeers     = 40
	inmemoryMessages  = 1024
)

// ErrStartedEngine is returned if the engine is already started
var ErrStartedEngine = errors.New("started engine")

var (
	// errUnknownBlock is returned when the list of committee is requested for a block
	// that is not part of the local blockchain.
	errUnknownBlock = errors.New("unknown block")
	// errUnauthorized is returned if a header is signed by a non authorized entity.
	errUnauthorized = errors.New("unauthorized")
)
var (
	defaultDifficulty = big.NewInt(1)
	nilUncleHash      = types.CalcUncleHash(nil) // Always Keccak256(RLP([])) as uncles are meaningless outside of PoW.
	emptyNonce        = types.BlockNonce{}
	now               = time.Now
)

// Author retrieves the Ethereum address of the account that minted the given
// block, which may be different from the header's coinbase if a consensus
// engine is based on signatures.
func (sb *Backend) Author(header *types.Header) (common.Address, error) {
	return types.Ecrecover(header)
}

// Prepare initializes the consensus fields of a block header according to the
// rules of a particular engine. The changes are executed inline.
func (sb *Backend) Prepare(chain consensus.ChainReader, header *types.Header) error {
	// unused fields, force to set to empty
	header.Coinbase = sb.address
	header.Nonce = emptyNonce
	header.MixDigest = types.BFTDigest

	// copy the parent extra data as the header extra data
	number := header.Number.Uint64()
	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	// use the same difficulty for all blocks
	header.Difficulty = defaultDifficulty

	// set header's timestamp
	header.Time = new(big.Int).Add(big.NewInt(int64(parent.Time)), new(big.Int).SetUint64(sb.config.BlockPeriod)).Uint64()
	if int64(header.Time) < time.Now().Unix() {
		header.Time = uint64(time.Now().Unix())
	}
	return nil
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
	return defaultDifficulty
}

func (sb *Backend) SetProposedBlockHash(hash common.Hash) {
	sb.proposedBlockHash = hash
}

// APIs returns the RPC APIs this consensus engine provides.
func (sb *Backend) APIs(chain consensus.ChainReader) []rpc.API {
	return []rpc.API{{
		Namespace: "tendermint",
		Version:   "1.0",
		Service:   &API{chain: chain, tendermint: sb, getCommittee: getCommittee},
		Public:    true,
	}}
}

// getCommittee retrieves the committee for the given header.
func getCommittee(header *types.Header, chain consensus.ChainReader) (types.Committee, error) {
	parent := chain.GetHeaderByHash(header.ParentHash)
	if parent == nil {
		return nil, errUnknownBlock
	}
	return parent.Committee, nil
}

// Start implements consensus.Start
func (sb *Backend) Start(ctx context.Context, blockchain *core.BlockChain) error {
	// the mutex along with coreStarted should prevent double start
	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()
	if sb.coreStarted {
		return ErrStartedEngine
	}

	// Set blockchain fields
	sb.blockchain = blockchain

	sb.stopped = make(chan struct{})

	// clear previous data
	sb.proposedBlockHash = common.Hash{}

	// Start Tendermint
	sb.core.Start(ctx, sb.autonityContract, sb.blockchain)
	sb.coreStarted = true

	return nil
}

// Stop implements consensus.Stop
func (sb *Backend) Close() error {
	// the mutex along with coreStarted should prevent double stop
	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()

	if !sb.coreStarted {
		return ErrStoppedEngine
	}

	// Stop Tendermint
	sb.core.Stop()
	sb.coreStarted = false
	close(sb.stopped)

	return nil
}

func (sb *Backend) SealHash(header *types.Header) common.Hash {
	return types.SigHash(header)
}
