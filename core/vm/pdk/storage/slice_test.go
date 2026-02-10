package storage

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/autonity/tests"
	"github.com/autonity/autonity/common"
)

func TestSliceLengthAlignment(t *testing.T) {
	r := tests.Setup(t, nil)
	contractAddr := common.HexToAddress("0x1234")
	baseSlot := common.Hash{}

	st := NewStorage(contractAddr, r.Evm.StateDB)

	// Create a slice and append an element
	var slice Slice[Var[uint64]]
	slice.Bind(st, baseSlot, 0)
	slice.Append(func(v *Var[uint64]) {
		v.Set(42)
	})
	st.Commit() // Commit changes to state

	rawSlot := r.Evm.StateDB.GetState(contractAddr, baseSlot)
	leftAlignedLen := binary.BigEndian.Uint64(rawSlot[:8])
	rightAlignedLen := binary.BigEndian.Uint64(rawSlot[24:32])

	t.Logf("Raw storage slot: %x", rawSlot)
	t.Logf("Left-aligned length (bytes 0-7): %d", leftAlignedLen)
	t.Logf("Right-aligned length (bytes 24-31): %d", rightAlignedLen)

	// PDK should now use right alignment (Solidity compatible)
	require.Equal(t, uint64(0), leftAlignedLen, "Left-aligned position should be empty")
	require.Equal(t, uint64(1), rightAlignedLen, "Length should be right-aligned")

	// Simulate what Solidity would write (right-aligned)
	var soliditySlot common.Hash
	binary.BigEndian.PutUint64(soliditySlot[24:32], 5) // Solidity writes length=5 right-aligned
	r.Evm.StateDB.SetState(contractAddr, baseSlot, soliditySlot)

	// Create new storage to bypass cache
	st2 := NewStorage(contractAddr, r.Evm.StateDB)
	var slice2 Slice[Var[uint64]]
	slice2.Bind(st2, baseSlot, 0)

	// PDK should read it correctly
	pdkLength := slice2.Len()
	t.Logf("Solidity wrote length=5 (right-aligned), PDK reads: %d", pdkLength)

	require.Equal(t, uint64(5), pdkLength, "PDK should read Solidity-written slice length correctly")
}

func TestSliceTruncationLeaksData(t *testing.T) {
	r := tests.Setup(t, nil)
	contractAddr := common.HexToAddress("0x1234")
	baseSlot := common.Hash{}

	st := NewStorage(contractAddr, r.Evm.StateDB)

	var slice Slice[Var[uint64]]
	slice.Bind(st, baseSlot, 0)

	// Append 3 elements
	slice.Append(func(v *Var[uint64]) { v.Set(100) })
	slice.Append(func(v *Var[uint64]) { v.Set(200) })
	slice.Append(func(v *Var[uint64]) { v.Set(300) })

	require.Equal(t, uint64(3), slice.Len())
	require.Equal(t, uint64(100), slice.Get(0).Get())
	require.Equal(t, uint64(200), slice.Get(1).Get())
	require.Equal(t, uint64(300), slice.Get(2).Get())

	st.Commit()

	// Truncate to length 1
	err := slice.setLength(1)
	require.NoError(t, err)
	st.Commit() // Commit truncation

	require.Equal(t, uint64(1), slice.Len())

	// Elements at indices 1 and 2 should be cleared
	oldElem1Slot := addSlot(slice.elemBaseSlot, 1)
	oldElem2Slot := addSlot(slice.elemBaseSlot, 2)

	elem1Data := r.Evm.StateDB.GetState(contractAddr, oldElem1Slot)
	elem2Data := r.Evm.StateDB.GetState(contractAddr, oldElem2Slot)

	zeroHash := common.Hash{}
	require.Equal(t, zeroHash, elem1Data, "Old data at index 1 should be cleared")
	require.Equal(t, zeroHash, elem2Data, "Old data at index 2 should be cleared")
}

func TestSliceElementOverlap(t *testing.T) {
	r := tests.Setup(t, nil)
	contractAddr := common.HexToAddress("0x1234")
	baseSlot := common.Hash{}

	st := NewStorage(contractAddr, r.Evm.StateDB)

	// Define a struct that takes 2 slots
	type BigStruct struct {
		A Var[Uint256] // Slot 0
		B Var[Uint256] // Slot 1
	}

	var slice Slice[BigStruct]
	slice.Bind(st, baseSlot, 0)

	// Append 2 elements
	slice.Append(func(v *BigStruct) {
		v.A.Set(NewUint256FromInt(0xAAAA))
		v.B.Set(NewUint256FromInt(0xBBBB))
	})
	slice.Append(func(v *BigStruct) {
		v.A.Set(NewUint256FromInt(0xCCCC))
		v.B.Set(NewUint256FromInt(0xDDDD))
	})

	require.Equal(t, uint64(2), slice.Len())

	// Index 0: Slot X, X+1
	// Index 1: Slot X+2, X+3 (Corrected from overlap)

	elem0 := slice.Get(0)
	val0A := elem0.A.Get()
	val0B := elem0.B.Get()

	// Verify no corruption
	expectedA := NewUint256FromInt(0xAAAA)
	expectedB := NewUint256FromInt(0xBBBB)

	require.Equal(t, expectedA, val0A)
	require.Equal(t, expectedB, val0B, "Element 0 Field B corrupted")
}

// TestSliceClear_DynamicTypes verifies that Clear() cleans up dynamic types
func TestSliceClear_DynamicTypes(t *testing.T) {
	r := tests.Setup(t, nil)
	contractAddr := common.HexToAddress("0x1234")
	baseSlot := common.Hash{}

	st := NewStorage(contractAddr, r.Evm.StateDB)

	var slice Slice[Var[[]byte]]
	slice.Bind(st, baseSlot, 0)

	// Append data
	data := make([]byte, 100) // 4 chunks
	for i := range data {
		data[i] = 0xAA
	}

	slice.Append(func(v *Var[[]byte]) {
		v.Set(data)
	})
	st.Commit()

	require.Equal(t, uint64(1), slice.Len())
	require.Equal(t, data, slice.Get(0).Get())

	// Clear the slice
	slice.Clear()
	st.Commit()

	require.Equal(t, uint64(0), slice.Len())
}
