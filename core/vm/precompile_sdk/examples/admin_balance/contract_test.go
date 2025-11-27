package admin_balance

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/autonity/tests"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm/precompile_sdk"
	"github.com/autonity/autonity/core/vm/precompile_sdk/abiselector"
	"github.com/autonity/autonity/core/vm/precompile_sdk/storage"
	"github.com/autonity/autonity/crypto"
)

// TestAdminBalanceContract_Init verifies defaults are set.
func TestAdminBalanceContract_Init(t *testing.T) {
	r := tests.Setup(t, nil)
	precompiledAddr := common.HexToAddress("0x1")
	c := NewAdminBalanceContract(r.Evm, precompiledAddr)

	precompile_sdk.AddToPrecompiles(precompile_sdk.PrecompileWrapper{Contract: c})

	// Check slots via storage.
	st := storage.NewStorage(c.Address, r.Evm.StateDB, c.Slots)
	admin, err := st.GetAddress("Admin")
	if err != nil {
		t.Fatal(err)
	}
	expectedAdmin := storage.Address(common.HexToAddress("0x000000000000000000000000000000000000dead"))
	if admin != expectedAdmin {
		t.Errorf("expected admin %v, got %v", expectedAdmin, admin)
	}

	balance, err := st.GetUint256("Balance")
	if err != nil {
		t.Fatal(err)
	}
	if balance.Eq(storage.NewUint256FromInt(100).Int) == false {
		t.Errorf("expected balance %v, got %v", storage.NewUint256FromInt(100).Int, balance)
	}
}

func TestAdminBalanceContract_FullFlow(t *testing.T) {
	runner := tests.Setup(t, nil)
	precompiledAddr := common.BigToAddress(big.NewInt(1))
	c := NewAdminBalanceContract(runner.Evm, precompiledAddr)
	// register
	precompile_sdk.AddToPrecompiles(precompile_sdk.PrecompileWrapper{Contract: c})

	// assign slots
	newAdmin := common.BytesToAddress([]byte("0xalive"))
	caller := common.BytesToAddress([]byte("adminCaller"))
	st := storage.NewStorage(c.Address, runner.Evm.StateDB, c.Slots)
	input := buildInput(t, c.Dispatcher, "UpdateAdmin", newAdmin)
	_, err := c.Run(input, runner.Evm.Context.BlockNumber.Uint64(), runner.Evm, caller)
	require.NoError(t, err, "UpdateAdmin failed")

	address, err := st.GetAddress("Admin")
	require.NoError(t, err, "GetAddress failed")
	require.Equal(t, newAdmin, address.ToCommonAddress(), "GetAddress mismatch")
	t.Log(address.ToCommonAddress().String())
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
