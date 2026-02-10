package tradeengine

import (
	"crypto/sha256"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/autonity/tests"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/core/vm/pdk"
	"github.com/autonity/autonity/core/vm/pdk/dispatcher"
	"github.com/autonity/autonity/crypto"
)

func TestSubmitOrder(t *testing.T) {
	r := tests.Setup(t, nil)

	te := SetupTradingEngineContract(r.Evm)

	pair := "NTN/USDC"
	side := uint8(0) // Bid
	price := big.NewInt(100)
	qty := big.NewInt(10)

	bc := vm.PrecompiledContractsIstanbul[ContractAddress]
	bcTyped, _ := bc.(*pdk.BaseContract)
	input := buildInput(t, bcTyped.Dispatcher, "SubmitOrder", pair, side, price, qty)
	sender := common.HexToAddress("0xdummyuser")

	//todo: use evm call to truly test precompile execution context
	result, err := bc.Run(input, r.Evm.Context.BlockNumber.Uint64(), r.Evm, sender)
	require.NoError(t, err)
	require.Greater(t, len(result), 0)

	orderID := common.BytesToHash(result[:32])
	require.NotEqual(t, common.Hash{}, orderID) // Non-zero ID

	ord := te.Orders.Get(orderID)
	require.Equal(t, sender, ord.User.Get())
	require.Equal(t, side, ord.Side.Get())
	ordPrice := ord.Price.Get()
	require.Equal(t, price.Uint64(), ordPrice.Uint64())
	ordQty := ord.Qty.Get()
	require.Equal(t, qty.Uint64(), ordQty.Uint64())
	require.Equal(t, uint8(0), ord.Status.Get()) // Open

	updatedBook := te.Books.Get(sha256.Sum256([]byte(pair)))
	uNextID := updatedBook.NextID.Get()
	require.Equal(t, uint64(2), uNextID.Uint64())

	gQty := updatedBook.Bids.Get(0).TotalQty.Get()
	require.Equal(t, qty.Uint64(), gQty.Uint64())
	require.Equal(t, uint64(1), updatedBook.Bids.Get(0).OrderIDs.Len()) // Grown append
	require.Equal(t, orderID, updatedBook.Bids.Get(0).OrderIDs.Get(0).Get())
}

func buildInput(t *testing.T, d *dispatcher.Dispatcher, methodName string, args ...interface{}) []byte {
	abiMethod, ok := d.ABI.Methods[methodName]
	require.True(t, ok, "Method not found: %s", methodName)

	sel := crypto.Keccak256Hash([]byte(abiMethod.Sig)).Bytes()[:4]
	argsPacked, err := abiMethod.Inputs.Pack(args...)
	require.NoError(t, err)
	input := append(sel, argsPacked...)
	return input
}

func TestDeleteOrder(t *testing.T) {
	r := tests.Setup(t, nil)
	te := SetupTradingEngineContract(r.Evm)

	pair := "NTN/USDC"
	side := uint8(0) // Bid
	price := big.NewInt(100)
	qty := big.NewInt(10)
	sender := common.HexToAddress("0xdummyuser")

	bc := vm.PrecompiledContractsIstanbul[ContractAddress]
	bcTyped, _ := bc.(*pdk.BaseContract)
	input := buildInput(t, bcTyped.Dispatcher, "SubmitOrder", pair, side, price, qty)

	result, err := bc.Run(input, r.Evm.Context.BlockNumber.Uint64(), r.Evm, sender)
	require.NoError(t, err)
	orderID := common.BytesToHash(result[:32])

	// Verify order exists
	ord := te.Orders.Get(orderID)
	require.Equal(t, sender, ord.User.Get())

	// Force Delete via PDK Map.Delete logic
	te.Orders.Delete(orderID)
	// We must verify it's gone.
	// Since Orders.Get(id) will return a bound wrapper for zero state if deleted,
	// checking primitive fields for zero values confirms deletion.
	ordDeleted := te.Orders.Get(orderID)

	require.Equal(t, common.Address{}, ordDeleted.User.Get())
	require.Equal(t, uint8(0), ordDeleted.Side.Get())
	// Uint256 zero check
	p := ordDeleted.Price.Get()
	require.Equal(t, uint64(0), p.Uint64())
}
