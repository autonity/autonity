package storage

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/autonity/tests"
	"github.com/autonity/autonity/common"
)

// TestByteAccessorSizeLimit tests the 1MB limit in ByteAccessor
func TestByteAccessorSizeLimit(t *testing.T) {
	r := tests.Setup(t, nil)
	contractAddr := common.HexToAddress("0x1234")
	slot := common.Hash{}

	st := NewStorage(contractAddr, r.Evm.StateDB)

	accessor := ByteAccessor{}

	// Test within limit
	smallData := make([]byte, 1000)
	for i := range smallData {
		smallData[i] = byte(i % 256)
	}

	err := accessor.WriteAt(slot, 0, smallData, st)
	require.NoError(t, err, "Small data should write successfully")

	// Test at limit
	limitData := make([]byte, 1<<20) // Exactly 1MB
	err = accessor.WriteAt(slot, 0, limitData, st)
	require.NoError(t, err, "Data at 1MB limit should write successfully")

	// Test over limit - should panic
	oversizeData := make([]byte, (1<<20)+1) // 1MB + 1 byte

	defer func() {
		if r := recover(); r != nil {
			t.Logf("Correctly panicked on oversized data: %v", r)
			require.Contains(t, r.(error).Error(), "bytes too large", "Should panic with size error")
		} else {
			t.Error("Should have panicked on data > 1MB")
		}
	}()

	err = accessor.WriteAt(slot, 0, oversizeData, st)
	t.Error("Should not reach here - should have panicked")
}

// TestBigIntAccessor_Int256Range tests BigIntAccessor with int256 range
func TestBigIntAccessor_Int256Range(t *testing.T) {
	r := tests.Setup(t, nil)
	contractAddr := common.HexToAddress("0x1234")
	slot := common.Hash{}

	st := NewStorage(contractAddr, r.Evm.StateDB)

	accessor := BigIntAccessor{}

	// Test maximum int256
	maxInt256 := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 255), big.NewInt(1))
	err := accessor.WriteAt(slot, 0, maxInt256, st)
	require.NoError(t, err, "Should write max int256")

	readVal, err := accessor.ReadAt(slot, 0, st)
	require.NoError(t, err)
	require.Equal(t, maxInt256, readVal.(*big.Int), "Should read back max int256")

	// Test minimum int256
	minInt256 := new(big.Int).Neg(new(big.Int).Lsh(big.NewInt(1), 255))
	err = accessor.WriteAt(slot, 0, minInt256, st)
	require.NoError(t, err, "Should write min int256")

	readVal, err = accessor.ReadAt(slot, 0, st)
	require.NoError(t, err)
	require.Equal(t, minInt256, readVal.(*big.Int), "Should read back min int256")

	// Test -1 (important for two's complement)
	negOne := big.NewInt(-1)
	err = accessor.WriteAt(slot, 0, negOne, st)
	require.NoError(t, err)

	readVal, err = accessor.ReadAt(slot, 0, st)
	require.NoError(t, err)
	require.Equal(t, negOne, readVal.(*big.Int), "Should read back -1")

	st.Commit() // Commit to flush cache to StateDB

	// Verify storage representation of -1 is all 0xff
	stored := r.Evm.StateDB.GetState(contractAddr, slot)
	expected := common.HexToHash("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	require.Equal(t, expected, stored, "-1 should be stored as all 0xff")

	// Test value exceeding max int256 - should error
	tooLarge := new(big.Int).Lsh(big.NewInt(1), 255) // 2^255 is one more than max
	err = accessor.WriteAt(slot, 0, tooLarge, st)
	require.Error(t, err, "Should error on value > max int256")
	require.Contains(t, err.Error(), "exceeds maximum", "Error should mention exceeds maximum")

	// Test value below min int256 - should error
	tooSmall := new(big.Int).Sub(minInt256, big.NewInt(1))
	err = accessor.WriteAt(slot, 0, tooSmall, st)
	require.Error(t, err, "Should error on value < min int256")
	require.Contains(t, err.Error(), "below minimum", "Error should mention below minimum")
}

// TestBigIntAccessor_NonZeroOffset tests that BigInt must start at offset 0
func TestBigIntAccessor_NonZeroOffset(t *testing.T) {
	r := tests.Setup(t, nil)
	contractAddr := common.HexToAddress("0x1234")
	slot := common.Hash{}

	st := NewStorage(contractAddr, r.Evm.StateDB)

	accessor := BigIntAccessor{}

	// Try to read at non-zero offset - should error
	_, err := accessor.ReadAt(slot, 12, st)
	require.Error(t, err, "BigInt read at non-zero offset should error")
	require.Contains(t, err.Error(), "zero offset", "Error should mention offset requirement")
}

// TestUint256Accessor_NonZeroOffset tests that Uint256 must start at offset 0
func TestUint256Accessor_NonZeroOffset(t *testing.T) {
	r := tests.Setup(t, nil)
	contractAddr := common.HexToAddress("0x1234")
	slot := common.Hash{}

	st := NewStorage(contractAddr, r.Evm.StateDB)

	accessor := Uint256Accessor{}

	// Try to read at non-zero offset - should error
	_, err := accessor.ReadAt(slot, 12, st)
	require.Error(t, err, "Uint256 read at non-zero offset should error")
	require.Contains(t, err.Error(), "zero offset", "Error should mention offset requirement")
}

// TestAddressAccessor_Packing tests that addresses can be written at offsets
func TestAddressAccessor_Packing(t *testing.T) {
	r := tests.Setup(t, nil)
	contractAddr := common.HexToAddress("0x1234")
	slot := common.Hash{}

	st := NewStorage(contractAddr, r.Evm.StateDB)

	accessor := AddressAccessor{}

	// Test address at offset 12 (20 + 12 = 32, exactly fills slot)
	// We cannot pack two addresses in one slot (20+20 > 32)
	addr := common.HexToAddress("0xBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB")

	err := accessor.WriteAt(slot, 12, addr, st)
	require.NoError(t, err)

	// Read back
	val, err := accessor.ReadAt(slot, 12, st)
	require.NoError(t, err)
	require.Equal(t, addr, val.(common.Address))

	// Verify raw storage alignment
	// Offset 12 means bytes [0-12] are untouched (zero), [12-32] are address
	st.Commit()
	raw := r.Evm.StateDB.GetState(contractAddr, slot)
	
	// Bytes 0-12 should be zero
	for i := 0; i < 12; i++ {
		require.Equal(t, byte(0), raw[i], "Byte at index %d should be zero", i)
	}
	// Bytes 12-32 should match address
	addrBytes := addr.Bytes()
	for i := 0; i < 20; i++ {
		require.Equal(t, addrBytes[i], raw[12+i], "Byte at index %d should match address", 12+i)
	}

	t.Log("Addresses write correctly at offsets")
}

// TestFixedByteAccessor_AllSizes tests all supported fixed byte array sizes [1]byte to [32]byte
func TestFixedByteAccessor_AllSizes(t *testing.T) {
	// Test all sizes from 1 to 32
	for size := 1; size <= 32; size++ {
		t.Run(string(rune(size))+"bytes", func(t *testing.T) {
			// Create accessor for this size
			accessor, ok := getAccessor(arrayTypeOf(size))
			require.True(t, ok, "Should have accessor for [%d]byte", size)
			require.Equal(t, size, accessor.Size(), "Size should be %d", size)

			t.Logf("[%d]byte accessor exists with size %d", size, accessor.Size())
		})
	}
}

// Helper function to get array type
func arrayTypeOf(size int) reflect.Type {
	switch size {
	case 1:
		return reflect.TypeOf([1]byte{})
	case 2:
		return reflect.TypeOf([2]byte{})
	case 4:
		return reflect.TypeOf([4]byte{})
	case 8:
		return reflect.TypeOf([8]byte{})
	case 16:
		return reflect.TypeOf([16]byte{})
	case 20:
		return reflect.TypeOf([20]byte{})
	case 32:
		return reflect.TypeOf([32]byte{})
	default:
		// For other sizes, construct using reflection
		return reflect.ArrayOf(size, reflect.TypeOf(byte(0)))
	}
}
