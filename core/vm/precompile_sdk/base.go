package precompile_sdk

import (
	"reflect"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/core/vm/precompile_sdk/abiselector"
	"github.com/autonity/autonity/core/vm/precompile_sdk/storage"
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

func (b *BaseContract) Run(input []byte, blockNumber uint64, evm *vm.EVM, caller common.Address) ([]byte, error) {
	st := storage.NewStorage(b.Address, evm.StateDB, b.Slots)
	return b.Dispatcher.Dispatch(input, evm, caller, st)
}

func (p *BaseContract) RequiredGas(input []byte) uint64 {
	// todo: we can keep it dynamic, for now a fixed value
	return 1000
}
