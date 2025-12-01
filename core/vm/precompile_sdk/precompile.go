package precompile_sdk

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm"
)

func AddToPrecompiles(address common.Address, c vm.PrecompiledContract) {
	vm.PrecompiledAddressesIstanbul = append(vm.PrecompiledAddressesIstanbul, address)
	if vm.PrecompiledContractsIstanbul == nil {
		vm.PrecompiledContractsIstanbul = make(map[common.Address]vm.PrecompiledContract)
	}
	vm.PrecompiledContractsIstanbul[address] = c
}
