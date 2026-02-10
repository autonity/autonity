// Package adminbalance provides an example contract for managing admin balances.
package adminbalance

import (
	"math/big"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/core/vm/pdk"
	"github.com/autonity/autonity/core/vm/pdk/examples/trade_engine"
	"github.com/autonity/autonity/core/vm/pdk/gas"
	"github.com/autonity/autonity/core/vm/pdk/storage"
	"github.com/autonity/autonity/log"
)

var (
	// ContractAddress is the example address for the precompiled contract.
	ContractAddress = common.HexToAddress("0x1") // example address for precompiled contract
)

// Contract is an example contract managing admin balances.
type Contract struct {
	Admin   storage.Var[common.Address]
	Balance storage.Var[storage.Uint256]
}

// SetupAdminBalanceContract initializes and registers the Contract precompile.
func SetupAdminBalanceContract(vm *vm.EVM) *Contract {
	c := &Contract{}

	gasConfig := gas.NewConfig(50_000)
	gasConfig.SetMethodGas("UpdateBalance", gas.MethodGas{Base: 80_000})
	gasConfig.SetMethodGas("UpdateAdmin", gas.MethodGas{Base: 80_000})
	gasConfig.SetMethodGas("CallSubmitOrder", gas.MethodGas{Base: 150_000})

	pdk.AddToPrecompiles(ContractAddress, c, vm, initialize, gasConfig)
	return c
}

func initialize(_ *storage.Storage, bc *pdk.BaseContract) {
	// set default values
	protocolAdmin := common.HexToAddress("0x000000000000000000000000000000000000dead")
	ac := bc.GetAppContract().(*Contract)

	ac.Admin.Set(protocolAdmin)
	ac.Balance.Set(storage.NewUint256FromInt(100))
}

// UpdateBalance updates the balance of the contract.
func (c *Contract) UpdateBalance(_ *vm.EVM, _ common.Address, _ *storage.Storage, newBalance *big.Int) error {
	// authorization checks
	newbal := storage.NewUint256FromBig(newBalance)
	c.Balance.Set(newbal)
	return nil
}

// UpdateAdmin updates the admin address of the contract.
func (c *Contract) UpdateAdmin(_ *vm.EVM, _ common.Address, _ *storage.Storage, adminAddress common.Address) error {
	c.Admin.Set(adminAddress)
	log.Info("admin updated", "newAdmin", c.Admin.Get())
	return nil
}

// CallSubmitOrder calls the SubmitOrder method on the TradingEngineContract.
func (c *Contract) CallSubmitOrder(evm *vm.EVM, caller common.Address, st *storage.Storage,
	teAddress common.Address, pair string, side uint8, price, qty *big.Int) (common.Hash, error) {
	cl := pdk.NewClient(teAddress, &tradeengine.TradingEngineContract{})
	result, err := cl.Call(evm, caller, st, "SubmitOrder", pair, side, price, qty)
	if err != nil {
		return common.Hash{}, err
	}
	orderID := common.BytesToHash(result[:32])
	return orderID, nil
}
