package admin_balance

import (
	"math/big"
	"reflect"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/core/vm/precompile_sdk"
	"github.com/autonity/autonity/core/vm/precompile_sdk/storage"
	"github.com/autonity/autonity/core/vm/precompile_sdk/types"
)

type AdminBalanceContract struct {
	*precompile_sdk.BaseContract
}

func NewAdminBalanceContract(evm *vm.EVM, address common.Address) *AdminBalanceContract {
	c := &AdminBalanceContract{}
	c.BaseContract = precompile_sdk.SetupContract(c, reflect.TypeOf(AdminBalanceState{}), address)
	// set default values
	st := storage.NewStorage(c.Address, evm.StateDB, c.Slots)
	protocolAdmin := types.Address(common.HexToAddress("0x000000000000000000000000000000000000dead"))
	err := st.SetAddress("Admin", protocolAdmin)
	if err != nil {
		panic(err)
	}
	err = st.SetUint256("Balance", types.NewUint256FromInt(100))
	if err != nil {
		panic(err)
	}
	return c
}

func (c *AdminBalanceContract) UpdateBalance(evm *vm.EVM, caller common.Address, st *storage.Storage, newBalance *big.Int) error {
	// authorization checks
	newbal := types.NewUint256FromBig(newBalance)
	return st.SetUint256("Balance", newbal)
}

func (c *AdminBalanceContract) UpdateAdmin(evm *vm.EVM, caller common.Address, st *storage.Storage, adminAddress common.Address) error {
	// authorization checks
	return st.SetAddress("Admin", types.NewAddressFromBytes(adminAddress.Bytes()))
}
