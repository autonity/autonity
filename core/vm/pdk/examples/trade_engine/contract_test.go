package trade_engine

import (
	"crypto/sha256"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/autonity/tests"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm/pdk"
	"github.com/autonity/autonity/core/vm/pdk/abiselector"
	"github.com/autonity/autonity/crypto"
)

func TestSubmitOrder(t *testing.T) {
	r := tests.Setup(t, nil)
	precompiledAddr := common.HexToAddress("0x1")
	c := NewTradingEngineContract(r.Evm, precompiledAddr)

	pdk.AddToPrecompiles(precompiledAddr, c)

	pair := "NTN/USDC"
	side := uint8(0) // Bid
	price := big.NewInt(100)
	qty := big.NewInt(10)

	input := buildInput(t, c.BaseContract.Dispatcher, "SubmitOrder", pair, side, price, qty)

	sender := common.HexToAddress("0xdummyuser")

	result, err := c.Run(input, r.Evm.Context.BlockNumber.Uint64(), r.Evm, sender)
	require.NoError(t, err)
	require.Greater(t, len(result), 0) // Return: ABI-packed orderID hash (32B)

	orderID := common.BytesToHash(result[:32])
	require.NotEqual(t, common.Hash{}, orderID) // Non-zero ID

	ord, err := c.repo.GetOrder(orderID)
	require.NoError(t, err)
	require.Equal(t, sender, ord.User)
	require.Equal(t, side, uint8(ord.Side))
	require.Equal(t, price.Uint64(), ord.Price.Uint64())
	require.Equal(t, qty.Uint64(), ord.Qty.Uint64())
	require.Equal(t, uint8(0), ord.Status) // Open

	updatedBook, err := c.repo.GetBook(sha256.Sum256([]byte(pair)))
	require.NoError(t, err)
	require.Equal(t, uint64(2), updatedBook.NextID.Uint64())

	// todo: verify/fix order book updated
	//require.Equal(t, qty.Uint64(), updatedBook.Bids[0].TotalQty.Uint64())
	//require.Equal(t, 1, len(updatedBook.Bids[0].OrderIDs)) // Grown append
	//require.Equal(t, orderID, updatedBook.Bids[0].OrderIDs[0])

}

func buildInput(t *testing.T, d *abiselector.Dispatcher, methodName string, args ...interface{}) []byte {
	abiMethod, ok := d.ABI.Methods[methodName]
	require.True(t, ok, "Method not found: %s", methodName)

	sel := crypto.Keccak256Hash([]byte(abiMethod.Sig)).Bytes()[:4]
	argsPacked, err := abiMethod.Inputs.Pack(args...)
	require.NoError(t, err)
	input := append(sel, argsPacked...)
	return input
}
