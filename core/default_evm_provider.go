package core

import (
	"math/big"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/state"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/core/vm"
	"github.com/clearmatics/autonity/params"
)

func NewDefaultEVMProvider(hg *headerGetter, vmConfig vm.Config, chainConfig *params.ChainConfig) *defaultEVMProvider {
	return &defaultEVMProvider{
		hg:          hg,
		vmConfig:    vmConfig,
		chainConfig: chainConfig,
	}
}

// defaultEVMProvider implements autonity.EVMProvider
type defaultEVMProvider struct {
	hg          *headerGetter
	vmConfig    vm.Config
	chainConfig *params.ChainConfig
}

func (p *defaultEVMProvider) EVM(header *types.Header, origin common.Address, statedb *state.StateDB) *vm.EVM {
	coinbase, _ := types.Ecrecover(header)
	evmContext := vm.Context{
		CanTransfer: CanTransfer,
		Transfer:    Transfer,
		GetHash:     GetHashFn(header, p.hg),
		Origin:      origin,
		Coinbase:    coinbase,
		BlockNumber: header.Number,
		Time:        new(big.Int).SetUint64(header.Time),
		GasLimit:    header.GasLimit,
		Difficulty:  header.Difficulty,
		GasPrice:    new(big.Int).SetUint64(0x0),
	}
	evm := vm.NewEVM(evmContext, statedb, p.chainConfig, p.vmConfig)
	return evm
}
