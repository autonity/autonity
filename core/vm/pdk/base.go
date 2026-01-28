package pdk

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/core/vm/pdk/abiselector"
	"github.com/autonity/autonity/core/vm/pdk/storage"
)

type BaseContract struct {
	contract interface{} // app contract instance

	Address    common.Address
	Dispatcher *abiselector.Dispatcher
}

func (b *BaseContract) GetAppContract() interface{} {
	return b.contract
}

func (b *BaseContract) Run(input []byte, _ uint64, evm *vm.EVM, caller common.Address) ([]byte, error) {
	st := storage.NewStorage(b.Address, evm.StateDB)
	storage.BindState(st, common.Hash{}, b.contract)
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

// AddToPrecompiles registers a PDK contract as a precompile at the given address.
func AddToPrecompiles(
	address common.Address,
	contractPtr interface{},
	evm *vm.EVM,
	initFunc func(st *storage.Storage, bc *BaseContract)) {
	base := &BaseContract{
		contract:   contractPtr,
		Address:    address,
		Dispatcher: abiselector.RegisterDispatcher(contractPtr),
	}
	st := storage.NewStorage(address, evm.StateDB)
	storage.BindState(st, common.Hash{}, contractPtr)

	if initFunc != nil {
		initFunc(st, base)
		st.Commit()
	}

	vm.PrecompiledContractRWMutex.Lock()
	defer vm.PrecompiledContractRWMutex.Unlock()

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

	// Update address lists for EIP-2929 warming
	vm.PrecompiledAddressesHomestead = append(vm.PrecompiledAddressesHomestead, address)
	vm.PrecompiledAddressesByzantium = append(vm.PrecompiledAddressesByzantium, address)
	vm.PrecompiledAddressesIstanbul = append(vm.PrecompiledAddressesIstanbul, address)
	vm.PrecompiledAddressesBerlin = append(vm.PrecompiledAddressesBerlin, address)
}
