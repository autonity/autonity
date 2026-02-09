package pdk

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/autonity/tests"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/core/vm/pdk/dispatcher"
	"github.com/autonity/autonity/core/vm/pdk/gas"
	"github.com/autonity/autonity/core/vm/pdk/storage"
	"github.com/autonity/autonity/params"
)

// TestAddToPrecompiles_UpdatesAddressLists verifies that PDK precompiles are added to
// the PrecompiledAddresses* slices for EIP-2929 warming
func TestAddToPrecompiles_UpdatesAddressLists(t *testing.T) {
	r := tests.Setup(t, nil)
	testAddr := common.HexToAddress("0xABCDEF")

	// Get initial list sizes
	rules := params.Rules{IsBerlin: true}
	initialAddresses := vm.ActivePrecompiles(rules)
	initialCount := len(initialAddresses)

	// Define a simple contract
	type TestContract struct{}

	// Add a PDK precompile
	tc := &TestContract{}
	gasConfig := gas.NewConfig(50_000)
	AddToPrecompiles(testAddr, tc, r.Evm, nil, gasConfig)

	// Verify address was added to all lists
	require.Contains(t, vm.PrecompiledAddressesHomestead, testAddr, "Address should be in Homestead list")
	require.Contains(t, vm.PrecompiledAddressesByzantium, testAddr, "Address should be in Byzantium list")
	require.Contains(t, vm.PrecompiledAddressesIstanbul, testAddr, "Address should be in Istanbul list")
	require.Contains(t, vm.PrecompiledAddressesBerlin, testAddr, "Address should be in Berlin list")

	// Verify ActivePrecompiles returns the new address
	activeAddresses := vm.ActivePrecompiles(rules)
	require.Len(t, activeAddresses, initialCount+1, "ActivePrecompiles should include new address")
	require.Contains(t, activeAddresses, testAddr, "ActivePrecompiles should contain PDK precompile")
}

// TestAddToPrecompiles_InitCommit verifies that initialization state is committed
func TestAddToPrecompiles_InitCommit(t *testing.T) {
	r := tests.Setup(t, nil)
	testAddr := common.HexToAddress("0x123456")

	type TestContract struct {
		Value storage.Var[common.Address]
	}

	tc := &TestContract{}
	initCalled := false
	initValue := common.HexToAddress("0xDEADBEEF")

	initFunc := func(_ *storage.Storage, bc *BaseContract) {
		initCalled = true
		contract := bc.GetAppContract().(*TestContract)
		contract.Value.Set(initValue)
	}

	gasConfig := gas.NewConfig(50_000)
	AddToPrecompiles(testAddr, tc, r.Evm, initFunc, gasConfig)

	require.True(t, initCalled, "Init function should have been called")

	// Verify the value was committed (can be read back)
	storedValue := tc.Value.Get()
	require.Equal(t, initValue, storedValue, "Init state should be persisted via auto-commit")
}

type PanicContract struct{}

func (p *PanicContract) Explode(_ *vm.EVM, _ common.Address, _ *storage.Storage) error {
	panic("boom")
}

func TestRun_PanicRecovery_Execution(t *testing.T) {
	r := tests.Setup(t, nil)
	testAddr := common.HexToAddress("0xPANIC")

	pc := &PanicContract{}
	gasConfig := gas.NewConfig(50_000)

	base := &BaseContract{
		contract:   pc,
		Address:    testAddr,
		Dispatcher: dispatcher.RegisterDispatcher(pc),
		gasConfig:  gasConfig,
	}

	err := gasConfig.Finalize(base.Dispatcher.ABI)
	require.NoError(t, err)

	method, ok := base.Dispatcher.ABI.Methods["Explode"]
	require.True(t, ok, "Method Explode should be registered")

	input := method.ID // Selector

	// Call Run - this should recover from panic
	ret, err := base.Run(input, 1000, r.Evm, common.Address{})

	// Verify result
	require.Error(t, err, "Run should return error on panic")
	require.Nil(t, ret, "Return value should be nil on error")
	require.Contains(t, err.Error(), "pdk precompile panic", "Error message should mention panic")
	require.Contains(t, err.Error(), "boom", "Error message should contain panic value")
}

type TestContract struct {
	Value storage.Var[uint64]
}

func (t *TestContract) GetValue(_ *vm.EVM, _ common.Address, _ *storage.Storage) (uint64, error) {
	return t.Value.Get(), nil
}

func (t *TestContract) SetValue(_ *vm.EVM, _ common.Address, _ *storage.Storage, value uint64) error {
	t.Value.Set(value)
	return nil
}

func TestBaseContractGas(t *testing.T) {
	r := tests.Setup(t, nil)
	contractAddr := common.HexToAddress("0x1234")
	contract := &TestContract{}

	gasConfig := gas.NewConfig(50_000)
	gasConfig.SetMethodGas("GetValue", gas.MethodGas{Base: 30_000})
	gasConfig.SetMethodGas("SetValue", gas.MethodGas{Base: 80_000})

	AddToPrecompiles(contractAddr, contract, r.Evm, nil, gasConfig)

	base := vm.PrecompiledContractsByzantium[contractAddr].(*BaseContract)

	getValueMethod := base.Dispatcher.ABI.Methods["GetValue"]
	getValueInput := getValueMethod.ID

	setValueMethod := base.Dispatcher.ABI.Methods["SetValue"]
	setValueInput := setValueMethod.ID

	getValueGas := base.RequiredGas(getValueInput)
	require.Equal(t, uint64(30_000), getValueGas)

	setValueGas := base.RequiredGas(setValueInput)
	require.Equal(t, uint64(80_000), setValueGas)

	unknownInput := []byte{0xFF, 0xFF, 0xFF, 0xFF}
	unknownGas := base.RequiredGas(unknownInput)
	require.Equal(t, uint64(50_000), unknownGas)
}

type BatchContract struct {
	Items storage.Slice[storage.Var[uint64]]
}

func (b *BatchContract) ProcessBatch(_ *vm.EVM, _ common.Address, _ *storage.Storage, items []uint64) error {
	for _, item := range items {
		b.Items.Append(func(v *storage.Var[uint64]) {
			v.Set(item)
		})
	}
	return nil
}

func TestDynamicGasCalculation(t *testing.T) {
	r := tests.Setup(t, nil)
	contractAddr := common.HexToAddress("0x5678")

	batchContract := &BatchContract{}

	gasConfig := gas.NewConfig(50_000)
	gasConfig.SetMethodGas("ProcessBatch", gas.MethodGas{
		Base: 100_000,
		Calculator: func(params []byte) uint64 {
			if len(params) < 32 {
				return 0
			}
			length := new(big.Int).SetBytes(params[:32]).Uint64()
			return length * 500
		},
	})

	AddToPrecompiles(contractAddr, batchContract, r.Evm, nil, gasConfig)

	base := vm.PrecompiledContractsByzantium[contractAddr].(*BaseContract)

	input := make([]byte, 36)
	method := base.Dispatcher.ABI.Methods["ProcessBatch"]
	copy(input[:4], method.ID[:4])
	length := big.NewInt(100)
	lengthBytes := length.Bytes()
	copy(input[4+32-len(lengthBytes):36], lengthBytes)

	gasUsed := base.RequiredGas(input)
	expectedGas := uint64(100_000 + 100*500)
	require.Equal(t, expectedGas, gasUsed)
}
