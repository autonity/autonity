package storage

import (
	"math/big"
	"testing"

	"github.com/autonity/autonity/autonity/tests"
	"github.com/autonity/autonity/common"
	"github.com/stretchr/testify/require"
)

var testAddr = common.HexToAddress("0x1234")

// TestBigIntPacking tests if multiple big.Ints can be packed and if they corrupt each other
func TestBigIntPacking(t *testing.T) {
	// Setup test state
	r := tests.Setup(t, nil)
	st := NewStorage(testAddr, r.Evm.StateDB)

	// Try to pack two big.Ints in a struct
	type TwoBigInts struct {
		First  Var[*big.Int]
		Second Var[*big.Int]
	}

	var s TwoBigInts
	BindState(st, common.Hash{}, &s)

	// Set first value
	val1 := big.NewInt(12345)
	s.First.Set(val1)
	st.Commit()

	// Set second value
	val2 := big.NewInt(67890)
	s.Second.Set(val2)
	st.Commit()

	// Read back first value - should not be corrupted
	readVal1 := s.First.Get()
	require.Equal(t, val1.String(), readVal1.String(), "First big.Int was corrupted by second write")

	// Read back second value
	readVal2 := s.Second.Get()
	require.Equal(t, val2.String(), readVal2.String(), "Second big.Int has wrong value")

	t.Logf("First.baseSlot: %s, offset: %d", s.First.baseSlot.Hex(), s.First.offset)
	t.Logf("Second.baseSlot: %s, offset: %d", s.Second.baseSlot.Hex(), s.Second.offset)
}

// TestMixedTypesPacking tests packing small types with big.Int
func TestMixedTypesPacking(t *testing.T) {
	r := tests.Setup(t, nil)
	st := NewStorage(testAddr, r.Evm.StateDB)

	type MixedStruct struct {
		SmallInt Var[uint8]
		BigInt   Var[*big.Int]
		AnotherSmallInt Var[uint16]
	}

	var s MixedStruct
	BindState(st, common.Hash{}, &s)

	// Log slot assignments
	t.Logf("SmallInt: slot=%s, offset=%d", s.SmallInt.baseSlot.Hex(), s.SmallInt.offset)
	t.Logf("BigInt: slot=%s, offset=%d", s.BigInt.baseSlot.Hex(), s.BigInt.offset)
	t.Logf("AnotherSmallInt: slot=%s, offset=%d", s.AnotherSmallInt.baseSlot.Hex(), s.AnotherSmallInt.offset)

	// Set all values
	s.SmallInt.Set(uint8(42))
	st.Commit()

	s.BigInt.Set(big.NewInt(999999))
	st.Commit()

	s.AnotherSmallInt.Set(uint16(1234))
	st.Commit()

	// Verify all values are correct
	require.Equal(t, uint8(42), s.SmallInt.Get(), "SmallInt corrupted")
	require.Equal(t, "999999", s.BigInt.Get().String(), "BigInt corrupted")
	require.Equal(t, uint16(1234), s.AnotherSmallInt.Get(), "AnotherSmallInt corrupted")
}

// TestPackedSmallTypes tests that small types can be safely packed
func TestPackedSmallTypes(t *testing.T) {
	r := tests.Setup(t, nil)
	st := NewStorage(testAddr, r.Evm.StateDB)

	type PackedSmall struct {
		A Var[uint8]
		B Var[uint16]
		C Var[uint32]
		D Var[uint64]
		E Var[bool]
	}

	var s PackedSmall
	BindState(st, common.Hash{}, &s)

	// Log slot assignments
	t.Logf("A: slot=%s, offset=%d", s.A.baseSlot.Hex(), s.A.offset)
	t.Logf("B: slot=%s, offset=%d", s.B.baseSlot.Hex(), s.B.offset)
	t.Logf("C: slot=%s, offset=%d", s.C.baseSlot.Hex(), s.C.offset)
	t.Logf("D: slot=%s, offset=%d", s.D.baseSlot.Hex(), s.D.offset)
	t.Logf("E: slot=%s, offset=%d", s.E.baseSlot.Hex(), s.E.offset)

	// Check if they're in the same slot
	sameSlot := s.A.baseSlot == s.B.baseSlot &&
		s.B.baseSlot == s.C.baseSlot &&
		s.C.baseSlot == s.D.baseSlot &&
		s.D.baseSlot == s.E.baseSlot

	if sameSlot {
		t.Log("All fields are packed in the same slot")
	} else {
		t.Log("Fields are spread across multiple slots")
	}

	// Set all values
	s.A.Set(uint8(10))
	s.B.Set(uint16(1000))
	s.C.Set(uint32(100000))
	s.D.Set(uint64(10000000))
	s.E.Set(true)
	st.Commit()

	// Verify all values
	require.Equal(t, uint8(10), s.A.Get(), "A corrupted")
	require.Equal(t, uint16(1000), s.B.Get(), "B corrupted")
	require.Equal(t, uint32(100000), s.C.Get(), "C corrupted")
	require.Equal(t, uint64(10000000), s.D.Get(), "D corrupted")
	require.Equal(t, true, s.E.Get(), "E corrupted")
}

// TestUint256AlwaysFullSlot verifies that Uint256 always gets a full slot
func TestUint256AlwaysFullSlot(t *testing.T) {
	r := tests.Setup(t, nil)
	st := NewStorage(testAddr, r.Evm.StateDB)

	type WithUint256 struct {
		Small1  Var[uint8]
		BigNum  Var[Uint256]
		Small2  Var[uint8]
	}

	var s WithUint256
	BindState(st, common.Hash{}, &s)

	t.Logf("Small1: slot=%s, offset=%d", s.Small1.baseSlot.Hex(), s.Small1.offset)
	t.Logf("BigNum: slot=%s, offset=%d", s.BigNum.baseSlot.Hex(), s.BigNum.offset)
	t.Logf("Small2: slot=%s, offset=%d", s.Small2.baseSlot.Hex(), s.Small2.offset)

	// Uint256 should always be at offset 0 in its own slot
	require.Equal(t, uint64(0), s.BigNum.offset, "Uint256 should always be at offset 0")

	// Set values
	s.Small1.Set(uint8(1))
	st.Commit()

	s.BigNum.Set(NewUint256FromInt(123456789))
	st.Commit()

	s.Small2.Set(uint8(2))
	st.Commit()

	// Verify no corruption
	require.Equal(t, uint8(1), s.Small1.Get())
	bigNumVal := s.BigNum.Get()
	require.Equal(t, uint64(123456789), bigNumVal.Uint64())
	require.Equal(t, uint8(2), s.Small2.Get())
}

// TestBytesAlwaysFullSlot verifies that []byte always gets its own slot
func TestBytesAlwaysFullSlot(t *testing.T) {
	r := tests.Setup(t, nil)
	st := NewStorage(testAddr, r.Evm.StateDB)

	type WithBytes struct {
		Small1 Var[uint8]
		Data   Var[[]byte]
		Small2 Var[uint8]
	}

	var s WithBytes
	BindState(st, common.Hash{}, &s)

	t.Logf("Small1: slot=%s, offset=%d", s.Small1.baseSlot.Hex(), s.Small1.offset)
	t.Logf("Data: slot=%s, offset=%d", s.Data.baseSlot.Hex(), s.Data.offset)
	t.Logf("Small2: slot=%s, offset=%d", s.Small2.baseSlot.Hex(), s.Small2.offset)

	// []byte should always be at offset 0 in its own slot
	require.Equal(t, uint64(0), s.Data.offset, "[]byte should always be at offset 0")

	// Set values
	s.Small1.Set(uint8(99))
	st.Commit()

	testData := []byte("Hello, PDK!")
	s.Data.Set(testData)
	st.Commit()

	s.Small2.Set(uint8(88))
	st.Commit()

	// Verify no corruption
	require.Equal(t, uint8(99), s.Small1.Get())
	require.Equal(t, testData, s.Data.Get())
	require.Equal(t, uint8(88), s.Small2.Get())
}

// TestAddressAccessorReadModifyWrite tests if address packing works correctly
func TestAddressAccessorReadModifyWrite(t *testing.T) {
	r := tests.Setup(t, nil)
	st := NewStorage(testAddr, r.Evm.StateDB)

	type WithAddresses struct {
		Addr1 Var[common.Address]
		Flag  Var[bool]
		Addr2 Var[common.Address]
	}

	var s WithAddresses
	BindState(st, common.Hash{}, &s)

	t.Logf("Addr1: slot=%s, offset=%d", s.Addr1.baseSlot.Hex(), s.Addr1.offset)
	t.Logf("Flag: slot=%s, offset=%d", s.Flag.baseSlot.Hex(), s.Flag.offset)
	t.Logf("Addr2: slot=%s, offset=%d", s.Addr2.baseSlot.Hex(), s.Addr2.offset)

	addr1 := common.HexToAddress("0xaaaa")
	addr2 := common.HexToAddress("0xbbbb")

	// Set addr1
	s.Addr1.Set(addr1)
	st.Commit()

	// Set flag
	s.Flag.Set(true)
	st.Commit()

	// Set addr2
	s.Addr2.Set(addr2)
	st.Commit()

	// Verify all values
	require.Equal(t, addr1, s.Addr1.Get(), "Addr1 corrupted")
	require.Equal(t, true, s.Flag.Get(), "Flag corrupted")
	require.Equal(t, addr2, s.Addr2.Get(), "Addr2 corrupted")
}

// TestWriteWithoutCommitThenRead tests that uncommitted writes are visible in same transaction
func TestWriteWithoutCommitThenRead(t *testing.T) {
	r := tests.Setup(t, nil)
	st := NewStorage(testAddr, r.Evm.StateDB)

	type Simple struct {
		Value Var[uint64]
	}

	var s Simple
	BindState(st, common.Hash{}, &s)

	// Write without commit
	s.Value.Set(uint64(42))

	// Read should see the uncommitted value
	require.Equal(t, uint64(42), s.Value.Get(), "Should see uncommitted value in same storage instance")

	// Commit
	st.Commit()

	// Read again
	require.Equal(t, uint64(42), s.Value.Get(), "Should see committed value")
}

// TestSliceWithBigIntElements tests if slice of big.Int works correctly
func TestSliceWithBigIntElements(t *testing.T) {
	r := tests.Setup(t, nil)
	st := NewStorage(testAddr, r.Evm.StateDB)

	type WithSlice struct {
		Numbers Slice[Var[*big.Int]]
	}

	var s WithSlice
	BindState(st, common.Hash{}, &s)

	// Append values
	s.Numbers.Append(func(v *Var[*big.Int]) {
		v.Set(big.NewInt(100))
	})
	st.Commit()

	s.Numbers.Append(func(v *Var[*big.Int]) {
		v.Set(big.NewInt(200))
	})
	st.Commit()

	s.Numbers.Append(func(v *Var[*big.Int]) {
		v.Set(big.NewInt(300))
	})
	st.Commit()

	// Verify length
	require.Equal(t, uint64(3), s.Numbers.Len())

	// Verify values
	require.Equal(t, "100", s.Numbers.Get(0).Get().String())
	require.Equal(t, "200", s.Numbers.Get(1).Get().String())
	require.Equal(t, "300", s.Numbers.Get(2).Get().String())
}

// TestOverwriteExistingValue tests that overwriting works correctly
func TestOverwriteExistingValue(t *testing.T) {
	r := tests.Setup(t, nil)
	st := NewStorage(testAddr, r.Evm.StateDB)

	type Simple struct {
		A Var[uint64]
		B Var[*big.Int]
	}

	var s Simple
	BindState(st, common.Hash{}, &s)

	// Initial values
	s.A.Set(uint64(111))
	s.B.Set(big.NewInt(222))
	st.Commit()

	// Overwrite
	s.A.Set(uint64(333))
	s.B.Set(big.NewInt(444))
	st.Commit()

	// Verify overwritten values
	require.Equal(t, uint64(333), s.A.Get())
	require.Equal(t, "444", s.B.Get().String())
}
