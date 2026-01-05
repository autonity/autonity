package admin_balance

import (
	"math/big"
	"testing"

	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/autonity/tests"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/core/vm/pdk"
	"github.com/autonity/autonity/core/vm/pdk/abiselector"
	"github.com/autonity/autonity/core/vm/pdk/examples/trade_engine"
	"github.com/autonity/autonity/core/vm/pdk/storage"
	"github.com/autonity/autonity/crypto"
)

// TestAdminBalanceContract_Init verifies defaults are set.
func TestAdminBalanceContract_Init(t *testing.T) {
	r := tests.Setup(t, nil)
	SetupAdminBalanceContract(r.Evm, ContractAddress)

	bc := vm.PrecompiledContractsIstanbul[ContractAddress]
	bcTyped, _ := bc.(*pdk.BaseContract)
	// Check slots via storage.
	st := storage.NewStorage(bcTyped.Address, r.Evm.StateDB, bcTyped.Slots)
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

	// assign slots
	SetupAdminBalanceContract(runner.Evm, ContractAddress)

	bc := vm.PrecompiledContractsIstanbul[ContractAddress]
	bcTyped, _ := bc.(*pdk.BaseContract)
	// Check slots via storage.
	st := storage.NewStorage(bcTyped.Address, runner.Evm.StateDB, bcTyped.Slots)
	newAdmin := common.BytesToAddress([]byte("0xalive"))
	input := buildInput(t, bcTyped.Dispatcher, "UpdateAdmin", newAdmin)
	_, err := bc.Run(input, runner.Evm.Context.BlockNumber.Uint64(), runner.Evm, contractOwner)
	require.NoError(t, err, "UpdateAdmin failed")

	address, err := storage.Get[common.Address](st.Field("Admin"))
	require.NoError(t, err, "GetAddress failed")
	require.Equal(t, newAdmin, address, "GetAddress mismatch")
}

func TestCrossContractCallInPrecompile_SubmitOrder(t *testing.T) {
	r := tests.Setup(t, nil)
	SetupAdminBalanceContract(r.Evm, ContractAddress)

	bc := vm.PrecompiledContractsIstanbul[ContractAddress]
	bcTyped, _ := bc.(*pdk.BaseContract)

	trade_engine.SetupTradingEngineContract(r.Evm, trade_engine.ContractAddress)

	pair := "NTN/USDC"
	side := uint8(0) // Bid
	price := big.NewInt(100)
	qty := big.NewInt(10)
	input := buildInput(t, bcTyped.Dispatcher, "CallSubmitOrder", trade_engine.ContractAddress,
		pair, side, price, qty)

	sender := common.HexToAddress("0xdummyuser")
	result, err := bc.Run(input, r.Evm.Context.BlockNumber.Uint64(), r.Evm, sender)
	require.NoError(t, err)
	require.Greater(t, len(result), 0) // Return: ABI-packed orderID hash (32B)

	orderID := common.BytesToHash(result[:32])
	t.Log("orderID", orderID)
}

/*
// commented out as it is dependent on trader object which is a temporary solidity contract
func TestCrossContractCallInSolidity_SubmitOrder(t *testing.T) {
	r := tests.Setup(t, nil)
	bc := NewAdminBalanceContract(r.Evm, ContractAddress)
	pdk.AddToPrecompiles(ContractAddress, bc)

	trade_engine.SetupTradingEngineContract(r.Evm, trade_engine.ContractAddress)

	pair := "NTN/USDC"
	side := uint8(0) // Bid
	price := big.NewInt(100)
	qty := big.NewInt(10)

	usedGas, err := r.Trader.Trade(nil, pair, side, price, qty)
	require.NoError(t, err)
	t.Log("used gas", usedGas)

	// get the logs to fetch the order ID
	logs := r.Evm.StateDB.GetLogs(common.Hash{}, common.Hash{})
	require.Greater(t, len(logs), 0, "no logs found")

	// event thrown by trading engine solidity wrapper
	eventSig := crypto.Keccak256Hash([]byte("OrderSubmitted(bytes32,string,uint8,uint256,uint256)"))

	found := false
	for _, log := range logs {
		if log.Topics[0] == eventSig {
			orderID := log.Topics[1]
			t.Log("orderID", orderID.Hex())
			found = true
		}
	}
	require.True(t, found)
}
*/

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
