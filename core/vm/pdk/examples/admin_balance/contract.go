package admin_balance

import (
	"math/big"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/core/vm/pdk"
	"github.com/autonity/autonity/core/vm/pdk/examples/trade_engine"
	"github.com/autonity/autonity/core/vm/pdk/storage"
)

var (
	contractOwner   = common.BytesToAddress([]byte("owner"))
	ContractAddress = common.HexToAddress("0x1") // example address for precompiled contract
)

type AdminBalanceContract struct {
	Admin   common.Address
	Balance storage.Uint256
}

func SetupAdminBalanceContract(vm *vm.EVM, address common.Address) *AdminBalanceContract {
	// todo: define address
	// address := common.HexToAddress("0x23")
	c := &AdminBalanceContract{}
	pdk.AddToPrecompiles(address, c, vm, initialize)
	return c
}

func initialize(st *storage.Storage) {
	// set default values
	protocolAdmin := common.HexToAddress("0x000000000000000000000000000000000000dead")
	err := storage.Set[common.Address](st.Field("Admin"), protocolAdmin)
	if err != nil {
		panic(err)
	}
	err = storage.Set[storage.Uint256](st.Field("Balance"), storage.NewUint256FromInt(100))
	if err != nil {
		panic(err)
	}
	st.Commit()
}

func (c *AdminBalanceContract) UpdateBalance(evm *vm.EVM, caller common.Address, st *storage.Storage, newBalance *big.Int) error {
	// authorization checks
	newbal := storage.NewUint256FromBig(newBalance)
	return storage.Set[storage.Uint256](st.Field("Balance"), newbal)
}

func (c *AdminBalanceContract) UpdateAdmin(evm *vm.EVM, caller common.Address, st *storage.Storage, adminAddress common.Address) error {
	// authorization checks

	return storage.Set[common.Address](st.Field("Admin"), adminAddress)
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
