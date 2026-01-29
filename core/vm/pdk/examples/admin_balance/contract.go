package admin_balance

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
	contractOwner   = common.BytesToAddress([]byte("owner"))
	ContractAddress = common.HexToAddress("0x1") // example address for precompiled contract
)

type AdminBalanceContract struct {
	Admin   storage.Var[common.Address]
	Balance storage.Var[storage.Uint256]
}

func SetupAdminBalanceContract(vm *vm.EVM) *AdminBalanceContract {
	c := &AdminBalanceContract{}

	gasConfig := gas.NewConfig(50_000)
	gasConfig.SetMethodGas("UpdateBalance", gas.MethodGas{Base: 80_000})
	gasConfig.SetMethodGas("UpdateAdmin", gas.MethodGas{Base: 80_000})
	gasConfig.SetMethodGas("CallSubmitOrder", gas.MethodGas{Base: 150_000})

	pdk.AddToPrecompiles(ContractAddress, c, vm, initialize, gasConfig)
	return c
}

func initialize(st *storage.Storage, bc *pdk.BaseContract) {
	// set default values
	protocolAdmin := common.HexToAddress("0x000000000000000000000000000000000000dead")
	ac := bc.GetAppContract().(*AdminBalanceContract)

	ac.Admin.Set(protocolAdmin)
	ac.Balance.Set(storage.NewUint256FromInt(100))
}

func (c *AdminBalanceContract) UpdateBalance(evm *vm.EVM, caller common.Address, st *storage.Storage, newBalance *big.Int) error {
	// authorization checks
	newbal := storage.NewUint256FromBig(newBalance)
	c.Balance.Set(newbal)
	return nil
}

func (c *AdminBalanceContract) UpdateAdmin(evm *vm.EVM, caller common.Address, st *storage.Storage, adminAddress common.Address) error {
	c.Admin.Set(adminAddress)
	log.Info("admin updated", "newAdmin", c.Admin.Get())
	return nil
}

func (c *AdminBalanceContract) CallSubmitOrder(evm *vm.EVM, caller common.Address, st *storage.Storage,
	teAddress common.Address, pair string, side uint8, price, qty *big.Int) (common.Hash, error) {
	cl := pdk.NewClient(teAddress, &trade_engine.TradingEngineContract{})
	result, err := cl.Call(evm, caller, st, "SubmitOrder", pair, side, price, qty)
	if err != nil {
		return common.Hash{}, err
	}
	orderId := common.BytesToHash(result[:32])
	return orderId, nil
}
