package admin_balance

import (
	"math/big"
	"reflect"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/core/vm/pdk"
	"github.com/autonity/autonity/core/vm/pdk/storage"
)

type AdminBalanceContract struct {
	*pdk.BaseContract
}

func NewAdminBalanceContract(evm *vm.EVM, address common.Address) *AdminBalanceContract {
	c := &AdminBalanceContract{}
	c.BaseContract = pdk.SetupContract(c, reflect.TypeOf(AdminBalanceState{}), address)
	// set default values
	st := storage.NewStorage(c.Address, evm.StateDB, c.Slots)
	protocolAdmin := common.HexToAddress("0x000000000000000000000000000000000000dead")
	err := storage.Set[common.Address](st.Field("Admin"), protocolAdmin)
	if err != nil {
		panic(err)
	}
	err = storage.Set[storage.Uint256](st.Field("Balance"), storage.NewUint256FromInt(100))
	if err != nil {
		panic(err)
	}
	return c
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
