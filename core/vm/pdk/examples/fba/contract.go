package fba

import (
	"fmt"
	"math/big"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/core/vm/pdk"
	"github.com/autonity/autonity/core/vm/pdk/gas"
	"github.com/autonity/autonity/core/vm/pdk/storage"
)

var ClearingContractAddress = common.HexToAddress("0x0000000000000000000000000000000000009999")

type Trade struct {
	Buyer    common.Address
	Seller   common.Address
	Quantity uint64
	Price    *big.Int
}

type ClearingContract struct {
	Margins           storage.Map[common.Address, storage.Var[*big.Int]]
	MaxTradesPerBatch uint64
}

func (c *ClearingContract) GetMargin(
	_ *vm.EVM,
	_ common.Address,
	_ *storage.Storage,
	account common.Address,
) (*big.Int, error) {
	return c.Margins.Get(account).Get(), nil
}

func (c *ClearingContract) SetMargin(
	_ *vm.EVM,
	_ common.Address,
	_ *storage.Storage,
	account common.Address,
	amount *big.Int,
) error {
	c.Margins.Get(account).Set(amount)
	return nil
}

func (c *ClearingContract) ProcessBatch(
	_ *vm.EVM,
	_ common.Address,
	_ *storage.Storage,
	trades []Trade,
) error {
	if uint64(len(trades)) > c.MaxTradesPerBatch {
		return fmt.Errorf("batch too large: %d > %d", len(trades), c.MaxTradesPerBatch)
	}

	for _, trade := range trades {
		buyerMargin := c.Margins.Get(trade.Buyer).Get()
		buyerMargin.Sub(buyerMargin, trade.Price)
		c.Margins.Get(trade.Buyer).Set(buyerMargin)

		sellerMargin := c.Margins.Get(trade.Seller).Get()
		sellerMargin.Add(sellerMargin, trade.Price)
		c.Margins.Get(trade.Seller).Set(sellerMargin)
	}

	return nil
}

func calculateProcessBatchGas(params []byte) uint64 {
	if len(params) < 32 {
		return 0
	}

	numTrades := new(big.Int).SetBytes(params[:32]).Uint64()
	tradesGas := numTrades * 500
	accountsGas := numTrades * 2 * 1000

	return tradesGas + accountsGas
}

func SetupClearingContract(evm *vm.EVM) *ClearingContract {
	contract := &ClearingContract{
		MaxTradesPerBatch: 1000,
	}

	gasConfig := gas.NewConfig(50_000)
	gasConfig.SetMethodGas("GetMargin", gas.MethodGas{Base: 30_000})
	gasConfig.SetMethodGas("SetMargin", gas.MethodGas{Base: 80_000})
	gasConfig.SetMethodGas("ProcessBatch", gas.MethodGas{
		Base:       100_000,
		Calculator: calculateProcessBatchGas,
	})

	pdk.AddToPrecompiles(ClearingContractAddress, contract, evm, nil, gasConfig)
	return contract
}
