package precompile_sdk

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm"
)

type PrecompileWrapper struct {
	Contract
}

func (p *PrecompileWrapper) RequiredGas(input []byte) uint64 {
	// todo: we can keep it dynamic, for now a fixed value
	return 1000
}
func (p *PrecompileWrapper) Run(input []byte, blockNumber uint64, evm *vm.EVM, caller common.Address) ([]byte, error) {
	return p.Run(input, blockNumber, evm, caller)
}

func (c *PrecompileWrapper) ContractAddress() common.Address {
	if base, ok := c.Contract.(*BaseContract); ok {
		return base.Address
	}
	return common.Address{}
}

func AddToPrecompiles(c PrecompileWrapper) {
	vm.PrecompiledAddressesIstanbul = append(vm.PrecompiledAddressesIstanbul, c.ContractAddress())
}
