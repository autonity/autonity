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
	"github.com/autonity/autonity/core/vm/pdk/examples/trade_engine"
	"github.com/autonity/autonity/core/vm/pdk/storage"
	"github.com/autonity/autonity/crypto"
)

// TestAdminBalanceContract_Init verifies defaults are set.
func TestAdminBalanceContract_Init(t *testing.T) {
	r := tests.Setup(t, nil)
	c := NewAdminBalanceContract(r.Evm, ContractAddress)

	pdk.AddToPrecompiles(ContractAddress, c)

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
	c := NewAdminBalanceContract(runner.Evm, ContractAddress)
	// register
	pdk.AddToPrecompiles(ContractAddress, c)

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

func TestCrossContractCallInPrecompile_SubmitOrder(t *testing.T) {
	r := tests.Setup(t, nil)
	bc := NewAdminBalanceContract(r.Evm, ContractAddress)
	pdk.AddToPrecompiles(ContractAddress, bc)

	trade_engine.SetupTradingEngineContract(r.Evm, trade_engine.ContractAddress)

	pair := "NTN/USDC"
	side := uint8(0) // Bid
	price := big.NewInt(100)
	qty := big.NewInt(10)
	input := buildInput(t, bc.BaseContract.Dispatcher, "CallSubmitOrder", trade_engine.ContractAddress,
		pair, side, price, qty)

	sender := common.HexToAddress("0xdummyuser")
	result, err := bc.Run(input, r.Evm.Context.BlockNumber.Uint64(), r.Evm, sender)
	require.NoError(t, err)
	require.Greater(t, len(result), 0) // Return: ABI-packed orderID hash (32B)

	orderID := common.BytesToHash(result[:32])
	require.NotEqual(t, common.Hash{}, orderID) // Non-zero ID
	t.Log("orderID", orderID)
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
