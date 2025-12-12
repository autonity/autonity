package admin_balance

import (
	"math/big"
	"testing"

	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/autonity/tests"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm/pdk"
	"github.com/autonity/autonity/core/vm/pdk/abiselector"
	"github.com/autonity/autonity/core/vm/pdk/storage"
	"github.com/autonity/autonity/crypto"
)

// TestAdminBalanceContract_Init verifies defaults are set.
func TestAdminBalanceContract_Init(t *testing.T) {
	r := tests.Setup(t, nil)
	precompiledAddr := common.HexToAddress("0x1")
	c := NewAdminBalanceContract(r.Evm, precompiledAddr)

	pdk.AddToPrecompiles(precompiledAddr, c)

	// Check slots via storage.
	st := storage.NewStorage(c.Address, r.Evm.StateDB, c.Slots)
	admin, err := storage.Get[common.Address](st.Field("Admin"))
	if err != nil {
		t.Fatal(err)
	}
	expectedAdmin := common.HexToAddress("0x000000000000000000000000000000000000dead")
	if admin != expectedAdmin {
		t.Errorf("expected admin %v, got %v", expectedAdmin, admin)
	}

	balance, err := storage.Get[storage.Uint256](st.Field("Balance"))
	if err != nil {
		t.Fatal(err)
	}
	expBal := &uint256.Int{100}
	if balance.Eq(expBal) == false {
		t.Errorf("expected balance %v, got %v", storage.NewUint256FromInt(100).Int, balance)
	}
}

func TestAdminBalanceContract_FullFlow(t *testing.T) {
	runner := tests.Setup(t, nil)
	precompiledAddr := common.BigToAddress(big.NewInt(1))
	c := NewAdminBalanceContract(runner.Evm, precompiledAddr)
	// register
	pdk.AddToPrecompiles(precompiledAddr, c)

	// assign slots
	newAdmin := common.BytesToAddress([]byte("0xalive"))
	st := storage.NewStorage(c.Address, runner.Evm.StateDB, c.Slots)
	input := buildInput(t, c.Dispatcher, "UpdateAdmin", newAdmin)
	_, err := c.Run(input, runner.Evm.Context.BlockNumber.Uint64(), runner.Evm, contractOwner)
	require.NoError(t, err, "UpdateAdmin failed")

	address, err := storage.Get[common.Address](st.Field("Admin"))
	require.NoError(t, err, "GetAddress failed")
	require.Equal(t, newAdmin, address, "GetAddress mismatch")
}

func buildInput(t *testing.T, d *abiselector.Dispatcher, methodName string, args ...interface{}) []byte {
	abiMethod, ok := d.ABI.Methods[methodName]
	if !ok {
		t.Error("Method not found")
	}
	sel := crypto.Keccak256Hash([]byte(abiMethod.Sig)).Bytes()[:4]

	argsPacked, err := abiMethod.Inputs.Pack(args...)
	require.NoError(t, err)
	input := append(sel, argsPacked...)
	return input
}
