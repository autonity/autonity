package pdk

import (
	"reflect"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/core/vm/pdk/abiselector"
	"github.com/autonity/autonity/core/vm/pdk/storage"
)

type BaseContract struct {
	Address    common.Address
	StateType  reflect.Type
	Slots      map[string]storage.SlotInfo
	Dispatcher *abiselector.Dispatcher
}

func SetupContract(childContract interface{}, stateType reflect.Type, addr common.Address) *BaseContract {
	base := &BaseContract{}
	base.StateType = stateType
	base.Dispatcher = abiselector.NewDispatcher()
	base.Slots = storage.AssignSlots(stateType)
	base.Address = addr
	err := abiselector.InferABIMethods(base.Dispatcher, reflect.ValueOf(childContract))
	if err != nil {
		panic(err)
	}
	return base
}

func (b *BaseContract) Run(input []byte, _ uint64, evm *vm.EVM, caller common.Address) ([]byte, error) {
	st := storage.NewStorage(b.Address, evm.StateDB, b.Slots)
	return b.Dispatcher.Dispatch(input, evm, caller, st)
}

func (b *BaseContract) RequiredGas(_ []byte) uint64 {
	// todo: we can keep it dynamic, for now a fixed value
	return 1000
}

func AddToPrecompiles(address common.Address, c vm.PrecompiledContract) {
	addToPrecompile := func(registry map[common.Address]vm.PrecompiledContract) {
		if registry == nil {
			registry = make(map[common.Address]vm.PrecompiledContract)
		}
		registry[address] = c
	}

	addToPrecompile(vm.PrecompiledContractsByzantium)
	addToPrecompile(vm.PrecompiledContractsHomestead)
	addToPrecompile(vm.PrecompiledContractsIstanbul)
	addToPrecompile(vm.PrecompiledContractsBerlin)
	addToPrecompile(vm.PrecompiledContractsBLS)
}
