package storage

import (
	"bytes"
	"encoding/binary"
	"math/big"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/crypto"
)

// TestEncodeTo32Bytes_SignExtension proves Bug #3: Sign extension happens AFTER copy
func TestEncodeTo32Bytes_SignExtension(t *testing.T) {
	tests := []struct {
		name     string
		key      interface{}
		keyType  reflect.Type
		expected []byte // What Solidity would produce
	}{
		{
			name:    "int8(-1)",
			key:     int8(-1),
			keyType: reflect.TypeOf(int8(0)),
			// -1 in two's complement is 0xFF
			// Sign-extended to 32 bytes should be all 0xFF
			expected: bytes.Repeat([]byte{0xff}, 32),
		},
		{
			name:     "int16(-1)",
			key:      int16(-1),
			keyType:  reflect.TypeOf(int16(0)),
			expected: bytes.Repeat([]byte{0xff}, 32),
		},
		{
			name:    "int32(-256)",
			key:     int32(-256),
			keyType: reflect.TypeOf(int32(0)),
			// -256 in two's complement 32-bit = 0xFFFFFF00
			// Sign-extended should be 0xFFFFFFFF FFFFFFFF FFFFFFFF FFFFFF00
			expected: func() []byte {
				b := bytes.Repeat([]byte{0xff}, 32)
				val := int32(-256)
				binary.BigEndian.PutUint32(b[28:32], uint32(val)) // #nosec G115
				return b
			}(),
		},
		{
			name:     "int64(-1)",
			key:      int64(-1),
			keyType:  reflect.TypeOf(int64(0)),
			expected: bytes.Repeat([]byte{0xff}, 32),
		},
		{
			name:    "int8(127) positive",
			key:     int8(127),
			keyType: reflect.TypeOf(int8(0)),
			// Positive numbers should have 0x00 padding
			expected: func() []byte {
				b := make([]byte, 32)
				b[31] = 0x7f
				return b
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encodeTo32Bytes(tt.key, tt.keyType)
			require.NoError(t, err)

			t.Logf("Key: %v", tt.key)
			t.Logf("Expected: %x", tt.expected)
			t.Logf("Got:      %x", got)

			if !bytes.Equal(got, tt.expected) {
				t.Logf("BUG PROVEN: Sign extension is incorrect")
				t.Logf("Expected all 0xff for negative numbers, but got different result")
			}

			require.Equal(t, tt.expected, got, "Sign extension should match Solidity encoding")
		})
	}
}

// TestMapKeyHashing_NegativeIntegers demonstrates how sign extension bug affects map keys
func TestMapKeyHashing_NegativeIntegers(t *testing.T) {
	baseSlot := common.Hash{} // Assume map is at slot 0

	// Test that int8(-1) and int8(-2) produce different hashes
	key1 := int8(-1)
	key2 := int8(-2)

	encoded1, _ := encodeTo32Bytes(key1, reflect.TypeOf(key1))
	encoded2, _ := encodeTo32Bytes(key2, reflect.TypeOf(key2))

	hash1 := crypto.Keccak256Hash(append(encoded1, baseSlot.Bytes()...))
	hash2 := crypto.Keccak256Hash(append(encoded2, baseSlot.Bytes()...))

	t.Logf("int8(-1) encoded: %x", encoded1)
	t.Logf("int8(-2) encoded: %x", encoded2)
	t.Logf("Hash for key -1: %x", hash1)
	t.Logf("Hash for key -2: %x", hash2)

	require.NotEqual(t, hash1, hash2, "Different keys should produce different hashes")

	// Verify sign extension for -1 (should be all 0xff)
	expectedNegOne := bytes.Repeat([]byte{0xff}, 32)
	if !bytes.Equal(encoded1, expectedNegOne) {
		t.Logf("BUG PROVEN: int8(-1) should encode as all 0xff bytes")
		t.Logf("Expected: %x", expectedNegOne)
		t.Logf("Got:      %x", encoded1)
	}
}

// TestEncodeTo32Bytes_BooleanEncoding verifies boolean map key encoding
func TestEncodeTo32Bytes_BooleanEncoding(t *testing.T) {
	tests := []struct {
		name     string
		key      bool
		expected []byte
	}{
		{
			name: "false",
			key:  false,
			expected: func() []byte {
				b := make([]byte, 32)
				// All zeros
				return b
			}(),
		},
		{
			name: "true",
			key:  true,
			expected: func() []byte {
				b := make([]byte, 32)
				b[31] = 1 // Right-aligned, value 1
				return b
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encodeTo32Bytes(tt.key, reflect.TypeOf(tt.key))
			require.NoError(t, err)

			t.Logf("Key: %v", tt.key)
			t.Logf("Expected: %x", tt.expected)
			t.Logf("Got:      %x", got)

			require.Equal(t, tt.expected, got, "Boolean encoding should match Solidity")
		})
	}
}

// TestMapKeyHashing_BooleanKeys demonstrates boolean map key hashing
func TestMapKeyHashing_BooleanKeys(t *testing.T) {
	baseSlot := common.Hash{}

	encodedTrue, _ := encodeTo32Bytes(true, reflect.TypeOf(true))
	encodedFalse, _ := encodeTo32Bytes(false, reflect.TypeOf(false))

	hashTrue := crypto.Keccak256Hash(append(encodedTrue, baseSlot.Bytes()...))
	hashFalse := crypto.Keccak256Hash(append(encodedFalse, baseSlot.Bytes()...))

	t.Logf("bool(true) encoded:  %x", encodedTrue)
	t.Logf("bool(false) encoded: %x", encodedFalse)
	t.Logf("Hash for true:  %x", hashTrue)
	t.Logf("Hash for false: %x", hashFalse)

	require.NotEqual(t, hashTrue, hashFalse, "true and false should produce different hashes")
}

// TestGetSlotConsumption_NilStorage tests the getSlotConsumption helper
func TestGetSlotConsumption_NilStorage(t *testing.T) {
	// This tests that getSlotConsumption works with nil storage
	// It should only be used for types that don't access storage during binding

	type SimpleStruct struct {
		A Var[uint64]
		B Var[common.Address]
	}

	// Should not panic
	consumed := getSlotConsumption[SimpleStruct]()
	t.Logf("SimpleStruct consumes %d slots", consumed)
	require.Greater(t, consumed, uint64(0))

	// More complex type with Map - this might panic if Map.Bind accesses storage
	type ComplexStruct struct {
		M Map[common.Address, Var[uint64]]
	}

	// This should work as Map.Bind doesn't access storage (just stores the reference)
	consumed2 := getSlotConsumption[ComplexStruct]()
	t.Logf("ComplexStruct consumes %d slots", consumed2)
	require.Greater(t, consumed2, uint64(0))
}

// TestAddSlot_Overflow tests slot arithmetic
func TestAddSlot_Overflow(t *testing.T) {
	// Test adding slots works correctly
	base := common.Hash{}
	result := addSlot(base, 1)

	expected := common.BigToHash(big.NewInt(1))
	require.Equal(t, expected, result)

	// Test large slot numbers
	largeBase := common.BigToHash(new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(10)))
	result2 := addSlot(largeBase, 5)

	expected2 := common.BigToHash(new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(5)))
	require.Equal(t, expected2, result2)
}

// TestSlotDiff tests slot difference calculation
func TestSlotDiff(t *testing.T) {
	base := common.BigToHash(big.NewInt(100))
	next := common.BigToHash(big.NewInt(150))

	diff := slotDiff(base, next)
	require.Equal(t, uint64(50), diff)

	// Test with zero diff
	diff2 := slotDiff(base, base)
	require.Equal(t, uint64(0), diff2)
}
