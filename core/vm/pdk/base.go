package pdk

import (
	"fmt"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/core/vm/pdk/abiselector"
	"github.com/autonity/autonity/core/vm/pdk/gas"
	"github.com/autonity/autonity/core/vm/pdk/storage"
)

type BaseContract struct {
	contract interface{}

	Address    common.Address
	Dispatcher *abiselector.Dispatcher
	gasConfig  *gas.Config
}

func (b *BaseContract) GetAppContract() interface{} {
	return b.contract
}

func (b *BaseContract) Run(input []byte, _ uint64, evm *vm.EVM, caller common.Address) (ret []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("pdk precompile panic: %v", r)
		}
	}()

	st := storage.NewStorage(b.Address, evm.StateDB)
	storage.BindState(st, common.Hash{}, b.contract)
	out, err := b.Dispatcher.Dispatch(input, evm, caller, st)
	if err != nil {
		return nil, err
	}
	st.Commit()
	return out, nil
}

func (b *BaseContract) RequiredGas(input []byte) uint64 {
	if b.gasConfig == nil {
		return 100_000
	}
	return b.gasConfig.GetGas(input)
}

// AddToPrecompiles registers a PDK contract as a precompile at the given address.
func AddToPrecompiles(
	address common.Address,
	contractPtr interface{},
	evm *vm.EVM,
	initFunc func(st *storage.Storage, bc *BaseContract),
	gasConfig *gas.Config) {
	if gasConfig == nil {
		panic("gas configuration is required for PDK contracts")
	}

	base := &BaseContract{
		contract:   contractPtr,
		Address:    address,
		Dispatcher: abiselector.RegisterDispatcher(contractPtr),
		gasConfig:  gasConfig,
	}

	if err := gasConfig.Finalize(base.Dispatcher.ABI); err != nil {
		panic(fmt.Errorf("failed to finalize gas config for contract at %s: %v", address.Hex(), err))
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
