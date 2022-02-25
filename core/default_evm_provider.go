package core

import (
	"math/big"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/state"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/core/vm"
)

// defaultEVMProvider implements autonity.EVMProvider
type defaultEVMProvider struct {
	bc *BlockChain
}

func (p *defaultEVMProvider) EVM(header *types.Header, origin common.Address, statedb *state.StateDB) *vm.EVM {
	evmContext := vm.BlockContext{
		CanTransfer: CanTransfer,
		Transfer:    Transfer,
		GetHash:     GetHashFn(header, p.bc),
		Coinbase:    header.Coinbase,
		BlockNumber: header.Number,
		Time:        new(big.Int).SetUint64(header.Time),
		GasLimit:    header.GasLimit,
		Difficulty:  header.Difficulty,
		BaseFee:     header.BaseFee,
	}
	txContext := vm.TxContext{
		Origin:   origin,
		GasPrice: new(big.Int).SetUint64(0x0),
	}
	vmConfig := *p.bc.GetVMConfig()
	evm := vm.NewEVM(evmContext, txContext, statedb, p.bc.Config(), vmConfig)
	return evm
}
