package pdk

import (
	"reflect"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/core/vm/pdk/abiselector"
	"github.com/autonity/autonity/core/vm/pdk/storage"
)

type BaseContract struct {
	contract interface{} // app contract instance

	Address    common.Address
	Slots      map[string]storage.SlotInfo
	Dispatcher *abiselector.Dispatcher
}

func (b *BaseContract) GetAppContract() interface{} {
	return b.contract
}

func (b *BaseContract) Run(input []byte, _ uint64, evm *vm.EVM, caller common.Address) ([]byte, error) {
	st := storage.NewStorage(b.Address, evm.StateDB, b.Slots)
	out, err := b.Dispatcher.Dispatch(input, evm, caller, st)
	if err != nil {
		return nil, err
	}
	st.Commit()
	return out, nil
}

func (b *BaseContract) RequiredGas(_ []byte) uint64 {
	// todo: we can keep it dynamic, for now a fixed value
	return 1000
}

func AddToPrecompiles(address common.Address, contractPtr interface{}, evm *vm.EVM, initFunc func(st *storage.Storage)) {
	elem := reflect.ValueOf(contractPtr).Elem()
	base := &BaseContract{
		contract:   contractPtr,
		Address:    address,
		Slots:      storage.AssignSlots(elem.Type()),
		Dispatcher: abiselector.GetOrRegisterDispatcher(contractPtr),
	}
	if initFunc != nil {
		st := storage.NewStorage(address, evm.StateDB, base.Slots)
		initFunc(st)
	}

	addToPrecompile := func(registry map[common.Address]vm.PrecompiledContract) {
		if registry == nil {
			registry = make(map[common.Address]vm.PrecompiledContract)
		}
		registry[address] = base
	}

	addToPrecompile(vm.PrecompiledContractsByzantium)
	addToPrecompile(vm.PrecompiledContractsHomestead)
	addToPrecompile(vm.PrecompiledContractsIstanbul)
	addToPrecompile(vm.PrecompiledContractsBerlin)
	addToPrecompile(vm.PrecompiledContractsBLS)
}
