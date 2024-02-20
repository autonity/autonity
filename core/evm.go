// Copyright 2016 The go-ethereum Authors
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

package core

import (
	"fmt"
	"math/big"
	"runtime/debug"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/log"
)

// ChainContext supports retrieving headers and consensus parameters from the
// current blockchain to be used during transaction processing.
type ChainContext interface {
	// GetHeader returns the hash corresponding to their hash.
	GetHeader(common.Hash, uint64) *types.Header
	// Engine retrieves the chain's consensus engine.
	Engine() consensus.Engine
	CurrentBlock() *types.Block
	Logger() log.Logger
	CurrentFastBlock() *types.Block
}

// NewEVMBlockContext creates a new context for use in the EVM.
func NewEVMBlockContext(header *types.Header, chain ChainContext, author *common.Address) vm.BlockContext {
	var (
		beneficiary common.Address
		baseFee     *big.Int
		random      *common.Hash
	)

	// If we don't have an explicit author (i.e. not mining), extract from the header
	if author == nil {
		beneficiary, _ = chain.Engine().Author(header) // Ignore error, we're past header validation
	} else {
		beneficiary = *author
	}
	if header.BaseFee != nil {
		baseFee = new(big.Int).Set(header.BaseFee)
	}
	if header.Difficulty.Cmp(common.Big0) == 0 {
		random = &header.MixDigest
	}
	blockHeight := chain.CurrentBlock().NumberU64()
	fastBlockHeight := chain.CurrentFastBlock().NumberU64()

	chain.Logger().Warn(fmt.Sprintf("setting number in blockContext. %d %d\n", blockHeight, header.Number.Uint64()))
	//chain.Logger().Warn(string(debug.Stack()[:]))

	if header.Number.Uint64() >= 1 && blockHeight < header.Number.Uint64()-1 {
		fmt.Println(blockHeight)
		fmt.Println(header.Number.Uint64())
		fmt.Println(fastBlockHeight)
		panic("test222222")
	}
	if header.Number.Uint64() == 32 {
		chain.Logger().Warn(fmt.Sprintf("32 with this node"))
		chain.Logger().Warn(string(debug.Stack()[:]))
	}

	return vm.BlockContext{
		CanTransfer: CanTransfer,
		Transfer:    Transfer,
		GetHash:     GetHashFn(header, chain),
		Coinbase:    beneficiary,
		BlockNumber: new(big.Int).Set(header.Number),
		Time:        new(big.Int).SetUint64(header.Time),
		Difficulty:  new(big.Int).Set(header.Difficulty),
		BaseFee:     baseFee,
		GasLimit:    header.GasLimit,
		Random:      random,
	}
}

// Used by the Autonity Contract
func GetDefaultEVM(chain *BlockChain) func(header *types.Header, origin common.Address, statedb vm.StateDB) *vm.EVM {
	return func(header *types.Header, origin common.Address, statedb vm.StateDB) *vm.EVM {
		evmContext := vm.BlockContext{
			CanTransfer: CanTransfer,
			Transfer:    Transfer,
			GetHash:     GetHashFn(header, chain),
			Coinbase:    header.Coinbase,
			BlockNumber: new(big.Int).Set(header.Number),
			Time:        new(big.Int).SetUint64(header.Time),
			GasLimit:    header.GasLimit,
			Difficulty:  header.Difficulty,
			BaseFee:     header.BaseFee,
		}
		txContext := vm.TxContext{
			Origin:   origin,
			GasPrice: new(big.Int).SetUint64(0x0),
		}
		evm := vm.NewEVM(evmContext, txContext, statedb, chain.chainConfig,
			vm.Config{
				//// Uncomment this to get EVM debugging logs
				//Debug: true,
				//Tracer: logger.NewMarkdownLogger(&logger.Config{
				//	EnableMemory:     true,
				//	DisableStack:     false,
				//	DisableStorage:   false,
				//	EnableReturnData: true,
				//	Debug:            true,
				//	Limit:            0,
				//	Overrides:        nil,
				//}, os.Stdout),
			},
		)
		return evm
	}
}

// NewEVMTxContext creates a new transaction context for a single transaction.
func NewEVMTxContext(msg Message) vm.TxContext {
	return vm.TxContext{
		Origin:   msg.From(),
		GasPrice: new(big.Int).Set(msg.GasPrice()),
	}
}

// GetHashFn returns a GetHashFunc which retrieves header hashes by number
func GetHashFn(ref *types.Header, chain ChainContext) func(n uint64) common.Hash {
	// Cache will initially contain [refHash.parent],
	// Then fill up with [refHash.p, refHash.pp, refHash.ppp, ...]
	var cache []common.Hash

	return func(n uint64) common.Hash {
		// If there's no hash cache yet, make one
		if len(cache) == 0 {
			cache = append(cache, ref.ParentHash)
		}
		if idx := ref.Number.Uint64() - n - 1; idx < uint64(len(cache)) {
			return cache[idx]
		}
		// No luck in the cache, but we can start iterating from the last element we already know
		lastKnownHash := cache[len(cache)-1]
		lastKnownNumber := ref.Number.Uint64() - uint64(len(cache))

		for {
			header := chain.GetHeader(lastKnownHash, lastKnownNumber)
			if header == nil {
				break
			}
			cache = append(cache, header.ParentHash)
			lastKnownHash = header.ParentHash
			lastKnownNumber = header.Number.Uint64() - 1
			if n == lastKnownNumber {
				return lastKnownHash
			}
		}
		return common.Hash{}
	}
}

// CanTransfer checks whether there are enough funds in the address' account to make a transfer.
// This does not take the necessary gas in to account to make the transfer valid.
func CanTransfer(db vm.StateDB, addr common.Address, amount *big.Int) bool {
	return db.GetBalance(addr).Cmp(amount) >= 0
}

// Transfer subtracts amount from sender and adds amount to recipient using the given Db
func Transfer(db vm.StateDB, sender, recipient common.Address, amount *big.Int) {
	db.SubBalance(sender, amount)
	db.AddBalance(recipient, amount)
}
