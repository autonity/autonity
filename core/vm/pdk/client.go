package pdk

import (
	"fmt"
	"math/big"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/core/vm/pdk/abiselector"
)

// Client is a PDK contract client which can call methods on a PDK contract deployed at a specific address.
type Client struct {
	address    common.Address
	dispatcher *abiselector.Dispatcher
}

func NewClient(address common.Address, logic interface{}) *Client {
	return &Client{
		address:    address,
		dispatcher: abiselector.GetOrRegisterDispatcher(logic),
	}
}

func (c *Client) Call(evm *vm.EVM, caller common.Address, method string, args ...interface{}) ([]byte, error) {
	abiMethod, ok := c.dispatcher.ABI.Methods[method]
	if !ok {
		return nil, fmt.Errorf("method not found: %s", method)
	}
	packedArgs, err := abiMethod.Inputs.Pack(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to pack arguments: %v", err)
	}
	input := append(abiMethod.ID, packedArgs...)

	//todo: Gas management
	gas := uint64(300_000)
	ret, _, err := evm.Call(vm.AccountRef(caller), c.address, input, gas, big.NewInt(0))
	if err != nil {
		return nil, err
	}
	return ret, nil
}
