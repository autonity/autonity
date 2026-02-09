package storage

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/autonity/tests"
	"github.com/autonity/autonity/common"
)

// TestArrayGetOutOfBounds
func TestArrays(t *testing.T) {
	r := tests.Setup(t, nil)
	contractAddr := common.HexToAddress("0x1234")
	baseSlot := common.Hash{}

	st := NewStorage(contractAddr, r.Evm.StateDB)

	// Array of uint64

	// Create an array with 3 elements (using [3]uint64 as shape)
	var arr Array[Var[uint64], [3]uint64]
	arr.Bind(st, baseSlot, 0)

	// Set valid elements
	arr.Get(0).Set(100)
	arr.Get(1).Set(200)
	arr.Get(2).Set(300)

	// Access valid indices - should work
	require.Equal(t, uint64(100), arr.Get(0).Get())
	require.Equal(t, uint64(200), arr.Get(1).Get())
	require.Equal(t, uint64(300), arr.Get(2).Get())

	// Access out-of-bounds index - should panic with clear message
	require.PanicsWithValue(t, "array index out of bounds", func() {
		arr.Get(5)
	})
}
