package pdk

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/autonity/tests"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm"
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
	AddToPrecompiles(testAddr, tc, r.Evm, nil)

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

	initFunc := func(st *storage.Storage, bc *BaseContract) {
		initCalled = true
		contract := bc.GetAppContract().(*TestContract)
		contract.Value.Set(initValue)
		// Note: We intentionally don't call st.Commit() here to test auto-commit
	}

	AddToPrecompiles(testAddr, tc, r.Evm, initFunc)

	require.True(t, initCalled, "Init function should have been called")

	// Verify the value was committed (can be read back)
	storedValue := tc.Value.Get()
	require.Equal(t, initValue, storedValue, "Init state should be persisted via auto-commit")
}
