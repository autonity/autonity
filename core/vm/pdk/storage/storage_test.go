package storage

import (
	"encoding/binary"
	"math/big"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/autonity/tests"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/crypto"
)

type Config struct {
	TimeoutSeconds uint32
	IsEnabled      bool
	Padding        [3]byte // to fill the complete 8 bytes
}

type Profile struct {
	NameHash common.Hash
	Level    uint64
	Config   Config
}

type DeepCommission struct {
	MaxRate      uint64
	LastModified int64
	Info         Config
}

// test struct with various field types
type TestState struct {
	// primitive types
	Version uint8
	Address common.Address
	Balance Uint256

	// nested struct
	AdminProfile Profile

	// dynamic array/slice
	History []common.Hash

	// map with struct
	UserStats map[common.Address]DeepCommission

	// map with primitive arrays
	Scores map[common.Address][]uint32

	Category   map[common.Address]uint64
	UserArrays map[common.Address][]Profile
}

var (
	testContractAddr = common.HexToAddress("0xCAFECAFECAFECAFECAFECAFECAFECAFECAFECAFE")
	testUserAddr     = common.HexToAddress("0x1134567890123456789012345678901234567890")
	testValidator    = common.HexToAddress("0x1234567890123456789012345678901234567890")
)

func newStorage(t *testing.T) *Storage {
	r := tests.Setup(t, nil)
	slots := AssignSlots(reflect.TypeOf(TestState{}))
	return NewStorage(testContractAddr, r.Evm.StateDB, slots)
}

func TestStoragePrimitives(t *testing.T) {
	st := newStorage(t)
	err := Set(st.Field("Version"), uint8(42))
	require.NoError(t, err)
	v, err := Get[uint8](st.Field("Version"))
	require.NoError(t, err)
	require.Equal(t, uint8(42), v)

	err = Set[common.Address](st.Field("Address"), testUserAddr)
	require.NoError(t, err)
	addr, err := Get[common.Address](st.Field("Address"))
	require.NoError(t, err)
	require.Equal(t, testUserAddr, addr)

	testBalance := NewUint256FromInt(1000000)
	err = Set(st.Field("Balance"), testBalance)
	require.NoError(t, err)
	bal, err := Get[Uint256](st.Field("Balance"))
	require.NoError(t, err)
	require.Equal(t, testBalance, bal)
}

func TestStorageNestedStruct(t *testing.T) {
	st := newStorage(t)
	testTimeout := uint32(300)

	// chaining for deep nested
	err := Set(st.Field("AdminProfile").Field("Config").Field("TimeoutSeconds"), testTimeout)
	require.NoError(t, err)
	timeout, err := Get[uint32](st.Field("AdminProfile").Field("Config").Field("TimeoutSeconds"))
	require.NoError(t, err)
	require.Equal(t, testTimeout, timeout)

	// set structs field (chaining)
	testLevel := uint64(300)
	err = Set(st.Field("AdminProfile").Field("Level"), testLevel)
	require.NoError(t, err)
	level, err := Get[uint64](st.Field("AdminProfile").Field("Level"))
	require.NoError(t, err)
	require.Equal(t, testLevel, level)
}

func TestStorageMapWithPrimitives(t *testing.T) {
	st := newStorage(t)
	testVal := uint64(95)
	err := Set(st.Field("Category").Map(testValidator), testVal)
	require.NoError(t, err)
	val, err := Get[uint64](st.Field("Category").Map(testValidator))
	require.NoError(t, err)
	require.Equal(t, testVal, val)
}

func TestStorageMapWithPrimitiveSlice(t *testing.T) {
	st := newStorage(t)
	testVal := uint32(95)

	// Set length to 1 (full slot for dynamic array)
	err := Set(st.Field("Scores").Map(testValidator), uint64(1))
	require.NoError(t, err)

	// Set element at index 0
	err = Set(st.Field("Scores").Map(testValidator).Index(0), testVal)
	require.NoError(t, err)

	// Get element at index 0
	val, err := Get[uint32](st.Field("Scores").Map(testValidator).Index(0))
	require.NoError(t, err)
	require.Equal(t, testVal, val)
}

func TestMapOfAddressToStructArrays(t *testing.T) {
	st := newStorage(t) // AssignSlots(reflect.TypeOf(TestState{}))

	testAddr := testUserAddr // Address wrapper
	testProfile := Profile{
		Level: 42,
		Config: Config{
			TimeoutSeconds: 300,
			IsEnabled:      true,
		},
	}

	// Set: Chain to map value (dynamic array), append profile, set sub-fields
	mapPath := st.Field("UserArrays").Map(testAddr)
	arr, err := NewArray[Profile](mapPath) // Wraps for push
	require.NoError(t, err)

	profileRef, err := arr.Grow()
	require.NoError(t, err)
	err = Set[uint64](profileRef, testProfile.Level)
	require.NoError(t, err)

	err = Set(profileRef.Field("Level"), testProfile.Level)
	require.NoError(t, err)

	err = Set[uint32](profileRef.Field("Config").Field("TimeoutSeconds"), testProfile.Config.TimeoutSeconds)
	require.NoError(t, err)
	err = Set[bool](profileRef.Field("Config").Field("IsEnabled"), testProfile.Config.IsEnabled)
	require.NoError(t, err)

	// Update nested: Index(0).Field("Config").Field("Enabled").Set(false)
	updatedPath := mapPath.Index(0).Field("Config").Field("IsEnabled")
	err = Set[bool](updatedPath, false)
	require.NoError(t, err)

	_, err = arr.Len()
	require.NoError(t, err)

	// Nested get: Direct field chain
	levelPath := st.Field("UserArrays").Map(testAddr).Index(0).Field("Level")
	gotLevel, err := Get[uint64](levelPath) // uintAccessor.ReadAt(struct_slot, level_offset=0, st)
	require.NoError(t, err)
	require.Equal(t, uint64(42), gotLevel)

	err = arr.Shrink()
	require.NoError(t, err)

}

type ComplexState struct {
	MapBool   map[bool]uint64
	MapUint8  map[uint8]string // Using []byte/string if supported, else uint64
	MapUint64 map[uint64]Profile

	// Mapping Address -> (ID -> Active Status)
	UserPermissions map[common.Address]map[uint64]bool

	// A list of historical snapshots, each snapshot is a map of ID->Value
	Snapshots []map[uint64]uint64

	SnapshotsFixedMapping [10]map[uint64]uint64
	SnapshotsFixedStruct  [10]Profile
	// 4 * 8 bytes = 32 bytes (Fits in 1 slot)
	PackedRates [3]uint64
	UserData    []byte // dynamic byte fields
}

// Helper to create storage for ComplexState
func newComplexStorage(t *testing.T) *Storage {
	r := tests.Setup(t, nil)
	slots := AssignSlots(reflect.TypeOf(ComplexState{}))
	return NewStorage(testContractAddr, r.Evm.StateDB, slots)
}

func TestMapKeyVariations(t *testing.T) {
	st := newComplexStorage(t)

	truePath := st.Field("MapBool").Map(true)
	falsePath := st.Field("MapBool").Map(false)

	require.NoError(t, Set(truePath, uint64(100)))
	require.NoError(t, Set(falsePath, uint64(200)))

	valTrue, err := Get[uint64](truePath)
	require.NoError(t, err)
	require.Equal(t, uint64(100), valTrue)

	valFalse, err := Get[uint64](falsePath)
	require.NoError(t, err)
	require.Equal(t, uint64(200), valFalse)

	profileMap := st.Field("MapUint64")
	key := uint64(9999)

	// Set Field in the struct inside map
	err = Set(profileMap.Map(key).Field("Level"), uint64(55))
	require.NoError(t, err)

	// Get Field back
	gotLevel, err := Get[uint64](profileMap.Map(key).Field("Level"))
	require.NoError(t, err)
	require.Equal(t, uint64(55), gotLevel)
}

func TestNestedMaps(t *testing.T) {
	st := newComplexStorage(t)

	rootMap := st.Field("UserPermissions")

	user := testUserAddr
	resourceID := uint64(42)

	// Construct Path: Root -> Map(Addr) -> Map(ID) -> Value
	targetPath := rootMap.Map(user).Map(resourceID)

	// 1. Set Permission to True
	err := Set(targetPath, true)
	require.NoError(t, err)

	// 2. Verify
	gotPerm, err := Get[bool](targetPath)
	require.NoError(t, err)
	require.True(t, gotPerm)

	// 3. Verify distinct key (different ID) is still false (default)
	otherPath := rootMap.Map(user).Map(uint64(99))
	gotOther, err := Get[bool](otherPath)
	require.NoError(t, err)
	require.False(t, gotOther)
}

func TestArrayOfMaps(t *testing.T) {
	st := newComplexStorage(t)

	snapshotsArr, err := NewArray[map[uint64]uint64](st.Field("Snapshots"))
	require.NoError(t, err)

	// 1. Add a new Snapshot (Index 0)
	snap0 := snapshotsArr.ReferenceAt(0)
	require.NoError(t, snap0.err)

	// 3. Set values in Snapshot 0
	// Snapshots[0][Key: 100] = 500
	err = Set(snap0.Map(uint64(100)), uint64(500))
	require.NoError(t, err)

	// 1. Add a new Snapshot (Index 0)
	snap1 := snapshotsArr.ReferenceAt(1)
	require.NoError(t, snap0.err)
	// 4. Set values in Snapshot 1
	// Snapshots[1][Key: 100] = 900 (Different value, same key, diff map)
	err = Set(snap1.Map(uint64(100)), uint64(900))
	require.NoError(t, err)

	// 5. Verify Isolation
	val0, err := Get[uint64](snapshotsArr.ReferenceAt(0).Map(uint64(100)))
	require.NoError(t, err)
	require.Equal(t, uint64(500), val0)

	val1, err := Get[uint64](snapshotsArr.ReferenceAt(1).Map(uint64(100)))
	require.NoError(t, err)
	require.Equal(t, uint64(900), val1)
}

func TestSliceOfMaps(t *testing.T) {
	st := newComplexStorage(t)

	snapshotsArr, err := NewArray[map[uint64]uint64](st.Field("Snapshots")) // Type T doesn't matter much for Next/Map calls
	require.NoError(t, err)

	// 1. Add a new Snapshot (Index 0)
	snap0, err := snapshotsArr.Grow()
	require.NoError(t, snap0.err)

	// 2. Add another Snapshot (Index 1)
	snap1, err := snapshotsArr.Grow()
	require.NoError(t, snap1.err)

	// 3. Set values in Snapshot 0
	// Snapshots[0][Key: 100] = 500
	err = Set(snap0.Map(uint64(100)), uint64(500))
	require.NoError(t, err)

	// 4. Set values in Snapshot 1
	// Snapshots[1][Key: 100] = 900 (Different value, same key, diff map)
	err = Set(snap1.Map(uint64(100)), uint64(900))
	require.NoError(t, err)

	// 5. Verify Isolation
	val0, err := Get[uint64](snapshotsArr.ReferenceAt(0).Map(uint64(100)))
	require.NoError(t, err)
	require.Equal(t, uint64(500), val0)

	val1, err := Get[uint64](snapshotsArr.ReferenceAt(1).Map(uint64(100)))
	require.NoError(t, err)
	require.Equal(t, uint64(900), val1)
}

func TestPackedFixedArray(t *testing.T) {
	st := newComplexStorage(t)

	ratesArr, err := NewArray[uint64](st.Field("PackedRates"))
	require.NoError(t, err)

	// Set Index 0 (Offset 0)
	err = Set(ratesArr.ReferenceAt(0), uint64(10))
	require.NoError(t, err)

	// Set Index 1 (Offset 8)
	err = Set(ratesArr.ReferenceAt(1), uint64(20))
	require.NoError(t, err)

	// Set Index 3 (Offset 24 - Last element in slot)
	err = Set(ratesArr.ReferenceAt(3), uint64(40))
	require.NoError(t, err)

	// Verify
	v0, _ := Get[uint64](ratesArr.ReferenceAt(0))
	v1, _ := Get[uint64](ratesArr.ReferenceAt(1))
	v3, _ := Get[uint64](ratesArr.ReferenceAt(3))

	require.Equal(t, uint64(10), v0)
	require.Equal(t, uint64(20), v1)
	require.Equal(t, uint64(40), v3)
}

func TestByteAccessor(t *testing.T) {
	st := newComplexStorage(t)

	bytesPath := st.Field("UserData")
	if bytesPath.err != nil {
		t.Fatalf("Failed to resolve UserData path: %v", bytesPath.err)
	}

	t.Run("EmptyBytes", func(t *testing.T) {
		empty := []byte{}
		require.NoError(t, Set(bytesPath, empty))

		got, err := Get[[]byte](bytesPath)
		require.NoError(t, err)
		require.Equal(t, empty, got)

		// Verify head slot: len=0, no data slots touched (should be zero)
		headSlot := bytesPath.slot
		headData := st.stateDB.GetState(st.address, headSlot)
		require.Equal(t, uint64(0), binary.BigEndian.Uint64(headData[:8]))
		baseSlot := crypto.Keccak256Hash(headSlot.Bytes())
		require.Equal(t, common.Hash{}, st.stateDB.GetState(st.address, baseSlot)) // Untouched chunk 0
	})

	t.Run("SmallBytesSingleChunk", func(t *testing.T) {
		data := []byte("hello world") // 11 bytes <32
		require.NoError(t, Set(bytesPath, data))

		got, err := Get[[]byte](bytesPath)
		require.NoError(t, err)
		require.Equal(t, data, got)

		headData := st.stateDB.GetState(st.address, bytesPath.slot)
		require.Equal(t, uint64(11), binary.BigEndian.Uint64(headData[:8]))

		baseSlot := crypto.Keccak256Hash(bytesPath.slot.Bytes())
		chunk0 := st.stateDB.GetState(st.address, baseSlot)
		require.Equal(t, data, chunk0[:11]) // First 11 bytes match
		empty21Bytes := [21]byte{}
		require.EqualValues(t, empty21Bytes[:], chunk0[11:]) // Rest zero-padded
	})

	t.Run("MediumBytesTwoChunks", func(t *testing.T) {
		data := make([]byte, 40) // 40 >32, spans two chunks
		copy(data, "this is a medium byte array that spans slots exactly")
		require.NoError(t, Set(bytesPath, data))

		got, err := Get[[]byte](bytesPath)
		require.NoError(t, err)
		require.Equal(t, data, got)

		// Verify chunks
		headData := st.stateDB.GetState(st.address, bytesPath.slot)
		require.Equal(t, uint64(40), binary.BigEndian.Uint64(headData[:8]))

		baseSlot := crypto.Keccak256Hash(bytesPath.slot.Bytes())
		chunk0 := st.stateDB.GetState(st.address, baseSlot) // Full 32B
		chunk1 := st.stateDB.GetState(st.address, common.BigToHash(new(big.Int).Add(new(big.Int).SetBytes(baseSlot.Bytes()), big.NewInt(1))))

		require.Equal(t, data[:32], chunk0[:])
		require.Equal(t, data[32:], chunk1[:8]) // Partial: 8B

		empty24Bytes := [24]byte{}
		require.Equal(t, empty24Bytes[:], chunk1[8:]) // Rest zero
	})

	t.Run("LargeBytesMultiChunk", func(t *testing.T) {
		data := make([]byte, 100) // Multi-chunk: 4 full (128B needed, but partial last)
		for i := range data {
			data[i] = byte('A' + (i % 26)) // Repeat pattern for easy verify
		}
		require.NoError(t, Set(bytesPath, data))

		got, err := Get[[]byte](bytesPath)
		require.NoError(t, err)
		require.Equal(t, data, got)

		headData := st.stateDB.GetState(st.address, bytesPath.slot)
		require.Equal(t, uint64(100), binary.BigEndian.Uint64(headData[:8]))

		baseSlot := crypto.Keccak256Hash(bytesPath.slot.Bytes())
		// Chunk 0: bytes 0-31
		chunk0 := st.stateDB.GetState(st.address, baseSlot)
		require.Equal(t, data[:32], chunk0[:])
		// Chunk 3: bytes 96-100 (partial)
		chunk3Slot := common.BigToHash(new(big.Int).Add(new(big.Int).SetBytes(baseSlot.Bytes()), big.NewInt(3)))
		chunk3 := st.stateDB.GetState(st.address, chunk3Slot)
		require.Equal(t, data[96:100], chunk3[:4])
		empty28Bytes := [28]byte{}
		require.Equal(t, empty28Bytes[:], chunk3[4:]) // Padded zero
	})

	t.Run("OverwriteResize", func(t *testing.T) {
		// Set initial large
		initial := make([]byte, 50)
		copy(initial, "initial data longer than 32")
		require.NoError(t, Set(bytesPath, initial))

		// Overwrite smaller: should clear extra chunks
		updated := []byte("shorter data")
		require.NoError(t, Set(bytesPath, updated))

		got, err := Get[[]byte](bytesPath)
		require.NoError(t, err)
		require.Equal(t, updated, got)
		require.NotEqual(t, initial, got) // Different

		headData := st.stateDB.GetState(st.address, bytesPath.slot)
		require.Equal(t, uint64(len(updated)), binary.BigEndian.Uint64(headData[:8])) // New len

		baseSlot := crypto.Keccak256Hash(bytesPath.slot.Bytes())
		chunk1 := st.stateDB.GetState(st.address, common.BigToHash(new(big.Int).Add(new(big.Int).SetBytes(baseSlot.Bytes()), big.NewInt(1))))
		require.Equal(t, common.Hash{}, chunk1) // Fully zeroed
	})

	t.Run("ErrorLargeSize", func(t *testing.T) {
		// Test panic on >1<<20 (1MB); use recover for coverage
		huge := make([]byte, 1<<20+1)
		defer func() {
			if r := recover(); r != nil {
				require.Contains(t, r.(error).Error(), "bytes too large") // Matches panic msg
			}
		}()
		// This should panic in WriteAt
		Set(bytesPath, huge)
		t.Error("Expected panic on oversized bytes") // Fail if no panic
	})
}

func Test2DArrays(t *testing.T) {
	t.Run("Primitive2D", func(t *testing.T) {
		type Primitive2D struct {
			Data [2][3]uint64 // Fixed 2D array of primitives: packed contiguous
		}
		r := tests.Setup(t, nil)
		slots := AssignSlots(reflect.TypeOf(Primitive2D{}))
		st := NewStorage(testContractAddr, r.Evm.StateDB, slots)

		path := st.Field("Data")
		require.NoError(t, path.err)

		// Set values: Data[0][0]=1, Data[0][1]=2, Data[1][2]=6 (spills across slots?)
		err := Set(path.Index(0).Index(0), uint64(5))
		require.NoError(t, err)
		err = Set(path.Index(0).Index(1), uint64(2))
		require.NoError(t, err)
		err = Set(path.Index(1).Index(2), uint64(6))
		require.NoError(t, err)

		// Verify
		v00, err := Get[uint64](path.Index(0).Index(0))
		require.NoError(t, err)
		require.Equal(t, uint64(5), v00)
		v01, err := Get[uint64](path.Index(0).Index(1))
		require.NoError(t, err)
		require.Equal(t, uint64(2), v01)
		v12, err := Get[uint64](path.Index(1).Index(2))
		require.NoError(t, err)
		require.Equal(t, uint64(6), v12)

		// Check default zeros
		v10, err := Get[uint64](path.Index(1).Index(0))
		require.NoError(t, err)
		require.Zero(t, v10)
	})

	t.Run("Dynamic2D", func(t *testing.T) {
		type Dynamic2D struct {
			Data [][]uint64 // Slice of slices: outer dynamic, inner dynamic (primitives)
		}
		r := tests.Setup(t, nil)
		slots := AssignSlots(reflect.TypeOf(Dynamic2D{}))
		st := NewStorage(testContractAddr, r.Evm.StateDB, slots)

		outerPath := st.Field("Data")
		require.NoError(t, outerPath.err)

		// Grow outer to len=2 (sets length at base)
		outerArr, err := NewArray[[]uint64](outerPath)
		require.NoError(t, err)
		_, err = outerArr.Grow() // Row 0
		require.NoError(t, err)
		_, err = outerArr.Grow() // Row 1
		require.NoError(t, err)

		// For row 0: Grow inner to len=3, set [0]=10, [1]=20
		row0Path := outerArr.ReferenceAt(0)
		inner0Arr, err := NewArray[uint64](row0Path)
		require.NoError(t, err)
		for i := 0; i < 3; i++ {
			_, err = inner0Arr.Grow()
			require.NoError(t, err)
		}
		err = Set(inner0Arr.ReferenceAt(0), uint64(10))
		require.NoError(t, err)
		err = Set(inner0Arr.ReferenceAt(1), uint64(20))
		require.NoError(t, err)

		// For row 1: Grow inner to len=2, set [0]=30
		row1Path := outerArr.ReferenceAt(1)
		inner1Arr, err := NewArray[uint64](row1Path)
		require.NoError(t, err)
		for i := 0; i < 2; i++ {
			_, err = inner1Arr.Grow()
			require.NoError(t, err)
		}
		err = Set(inner1Arr.ReferenceAt(0), uint64(30))
		require.NoError(t, err)

		//// Verify isolation
		out0 := outerArr.ReferenceAt(0)
		out0Arr, err := NewArray[uint64](out0)
		require.NoError(t, err)
		got00, err := Get[uint64](out0Arr.ReferenceAt(0))
		require.NoError(t, err)
		require.Equal(t, uint64(10), got00)

		got01, err := Get[uint64](out0Arr.ReferenceAt(1))
		require.NoError(t, err)
		require.Equal(t, uint64(20), got01)

		out1 := outerArr.ReferenceAt(1)
		out1Arr, err := NewArray[uint64](out1)
		require.NoError(t, err)
		got10, err := Get[uint64](out1Arr.ReferenceAt(0))
		require.NoError(t, err)
		require.Equal(t, uint64(30), got10)

		// todo: shrink doesn't work for 2D slices
		//err = out1Arr.Shrink()
		//require.NoError(t, err)
		//
		//got10, err = Get[uint64](out1Arr.ReferenceAt(1))
		//require.NoError(t, err)
	})

	t.Run("FixedDynamic2D", func(t *testing.T) {
		// Test fixed outer, dynamic inner: [2][]uint64 – exposes layout gaps (dynamic elems in fixed)
		type FixedDynamic2D struct {
			Data [2][]uint64 // Fixed array of dynamic slices: Should err or under-resolve
		}
		r := tests.Setup(t, nil)
		slots := AssignSlots(reflect.TypeOf(FixedDynamic2D{}))
		st := NewStorage(testContractAddr, r.Evm.StateDB, slots)

		path := st.Field("Data")
		require.NoError(t, path.err)

		row0Path := path.Index(0)
		inner0Arr, err := NewArray[uint64](row0Path)
		require.NoError(t, err)
		_, err = inner0Arr.Grow()
		require.NoError(t, err)
		err = Set(inner0Arr.ReferenceAt(0), uint64(123))
		require.NoError(t, err)

		row0Path = path.Index(0)
		inner0Arr, _ = NewArray[uint64](row0Path)
		val, err := Get[uint64](inner0Arr.ReferenceAt(0))
		require.NoError(t, err)
		require.Equal(t, uint64(123), val)

	})
}
