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

	catMap := NewMap[common.Address, uint64](st.Field("Category"))
	err := catMap.Set(testValidator, testVal)
	require.NoError(t, err)

	val, err := catMap.Get(testValidator)
	require.NoError(t, err)
	require.Equal(t, testVal, val)
}

func TestStorageMapWithPrimitiveSlice(t *testing.T) {
	st := newStorage(t)
	testVal := uint32(95)

	// Set length to 1 (full slot for dynamic array)
	scoreMapPath := st.Field("Scores").Map(testValidator)
	scoresSlice := NewSlice[uint32](scoreMapPath)
	err := scoresSlice.Append(testVal)
	require.NoError(t, err)

	length, err := scoresSlice.Len()
	require.NoError(t, err)
	require.Equal(t, uint64(1), length)

	// Get element at index 0
	val, err := scoresSlice.Get(0)
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
	profileSlice := NewSlice[Profile](mapPath) // Wraps for push

	// append new Profile
	err := profileSlice.Append(testProfile)
	require.NoError(t, err)

	// verify len
	length, err := profileSlice.Len()
	require.NoError(t, err)
	require.Equal(t, uint64(1), length)
	prof, err := profileSlice.Get(0)
	require.NoError(t, err)

	prof.Config.IsEnabled = false
	err = profileSlice.Set(0, prof)
	require.NoError(t, err)

	// verify update
	prof, err = profileSlice.Get(0)
	require.NoError(t, err)
	require.Equal(t, uint64(42), prof.Level)
	require.Equal(t, false, prof.Config.IsEnabled)

	// pop
	err = profileSlice.Pop()
	require.NoError(t, err)
	length, err = profileSlice.Len()
	require.NoError(t, err)
	require.Equal(t, uint64(0), length)
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

	boolMap := NewMap[bool, uint64](st.Field("MapBool"))

	require.NoError(t, boolMap.Set(true, 100))
	require.NoError(t, boolMap.Set(false, 200))

	valTrue, err := boolMap.Get(true)
	require.NoError(t, err)
	require.Equal(t, uint64(100), valTrue)

	valFalse, err := boolMap.Get(false)
	require.NoError(t, err)
	require.Equal(t, uint64(200), valFalse)

	profileMap := NewMap[uint64, Profile](st.Field("MapUint64"))
	key := uint64(9999)
	p := Profile{
		Level: 10,
	}
	require.NoError(t, profileMap.Set(key, p))
	p1, err := profileMap.Get(key)
	require.NoError(t, err)
	require.Equal(t, uint64(10), p1.Level)
}

func TestNestedMaps(t *testing.T) {
	st := newComplexStorage(t)
	rootMap := st.Field("UserPermissions")
	user := testUserAddr
	resourceID := uint64(42)

	// Construct Path: Root -> Map(Addr) -> Map(ID) -> Value
	innerMapPath := rootMap.Map(user)
	innerMap := NewMap[uint64, bool](innerMapPath)
	require.NoError(t, innerMap.Set(resourceID, true))

	gotPerm, err := innerMap.Get(resourceID)
	require.NoError(t, err)
	require.True(t, gotPerm)

	// Verify distinct key (different ID) is still false (default)
	gotOther, err := innerMap.Get(99)
	require.NoError(t, err)
	require.False(t, gotOther)
}

func TestSliceOfMaps(t *testing.T) {
	st := newComplexStorage(t)

	snapshotsArr := NewSlice[Map[uint64, uint64]](st.Field("Snapshots"))

	// 1. Add a new Snapshot (Index 0)
	map0, err := snapshotsArr.Grow()
	require.NoError(t, err)

	err = map0.Set(100, 500)
	require.NoError(t, err)

	// 1. Add a new Snapshot (Index 0)
	map1, err := snapshotsArr.Grow()
	require.NoError(t, err)
	// 4. Set values in Snapshot 1
	// Snapshots[1][Key: 100] = 900 (Different value, same key, diff map)
	err = map1.Set(100, 900)
	require.NoError(t, err)

	// 5. Verify Isolation
	val0Map, err := snapshotsArr.Get(0)
	require.NoError(t, err)
	val0, err := val0Map.Get(100)
	require.NoError(t, err)
	require.Equal(t, uint64(500), val0)

	val1Map, err := snapshotsArr.Get(1)
	require.NoError(t, err)
	val1, err := val1Map.Get(100)
	require.NoError(t, err)
	require.Equal(t, uint64(900), val1)
}

func TestPackedFixedArray(t *testing.T) {
	st := newComplexStorage(t)

	ratesArr, err := NewArray[uint64](st.Field("PackedRates"))
	require.NoError(t, err)

	// Set Index 0 (Offset 0)
	err = ratesArr.Set(0, uint64(10))
	require.NoError(t, err)

	// Set Index 1 (Offset 8)
	err = ratesArr.Set(1, uint64(20))
	require.NoError(t, err)

	// Set Index 3 - index out of bounds
	err = ratesArr.Set(3, uint64(40))
	require.Error(t, err)

	// Verify
	v0, _ := ratesArr.Get(0)
	v1, _ := ratesArr.Get(1)

	require.Equal(t, uint64(10), v0)
	require.Equal(t, uint64(20), v1)
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
		headData := st.GetState(headSlot)
		require.Equal(t, uint64(0), binary.BigEndian.Uint64(headData[24:]))
		baseSlot := crypto.Keccak256Hash(headSlot.Bytes())
		require.Equal(t, common.Hash{}, st.stateDB.GetState(st.address, baseSlot)) // Untouched chunk 0
	})

	t.Run("SmallBytesSingleChunk", func(t *testing.T) {
		data := []byte("hello world") // 11 bytes <32
		require.NoError(t, Set(bytesPath, data))

		got, err := Get[[]byte](bytesPath)
		require.NoError(t, err)
		require.Equal(t, data, got)

		headData := st.GetState(bytesPath.slot)
		require.Equal(t, uint64(11), binary.BigEndian.Uint64(headData[24:]))

		baseSlot := crypto.Keccak256Hash(bytesPath.slot.Bytes())
		chunk0 := st.GetState(baseSlot)
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
		headData := st.GetState(bytesPath.slot)
		require.Equal(t, uint64(40), binary.BigEndian.Uint64(headData[24:]))

		baseSlot := crypto.Keccak256Hash(bytesPath.slot.Bytes())
		chunk0 := st.GetState(baseSlot) // Full 32B
		chunk1 := st.GetState(common.BigToHash(new(big.Int).Add(new(big.Int).SetBytes(baseSlot.Bytes()), big.NewInt(1))))

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

		headData := st.GetState(bytesPath.slot)
		require.Equal(t, uint64(100), binary.BigEndian.Uint64(headData[24:]))

		baseSlot := crypto.Keccak256Hash(bytesPath.slot.Bytes())
		// Chunk 0: bytes 0-31
		chunk0 := st.GetState(baseSlot)
		require.Equal(t, data[:32], chunk0[:])
		// Chunk 3: bytes 96-100 (partial)
		chunk3Slot := common.BigToHash(new(big.Int).Add(new(big.Int).SetBytes(baseSlot.Bytes()), big.NewInt(3)))
		chunk3 := st.GetState(chunk3Slot)
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
		require.NotEqual(t, initial, got)

		headData := st.GetState(bytesPath.slot)
		require.Equal(t, uint64(len(updated)), binary.BigEndian.Uint64(headData[24:])) // New len

		baseSlot := crypto.Keccak256Hash(bytesPath.slot.Bytes())
		chunk1 := st.GetState(common.BigToHash(new(big.Int).Add(new(big.Int).SetBytes(baseSlot.Bytes()), big.NewInt(1))))
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

		// Grow outer to len=2 (sets length at base)
		outerArr := NewSlice[Slice[uint64]](st.Field("Data"))
		row0, err := outerArr.Grow() // Row 0
		require.NoError(t, err)

		// For row 0: Grow inner to len=3, set [0]=10, [1]=20
		require.NoError(t, row0.Append(10))
		require.NoError(t, row0.Append(20))

		val00, err := row0.Get(0)
		require.Equal(t, uint64(10), val00)
		val01, err := row0.Get(1)
		require.Equal(t, uint64(20), val01)

		row1, err := outerArr.Grow() // Row 1
		require.NoError(t, err)

		require.NoError(t, row1.Append(30))
		require.NoError(t, row1.Append(40))

		val10, err := row1.Get(0)
		require.Equal(t, uint64(30), val10)
		val11, err := row1.Get(1)
		require.Equal(t, uint64(40), val11)
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

		outerArr, err := NewArray[[]uint64](path)
		require.NoError(t, err)
		require.Equal(t, uint64(2), outerArr.Len())

		row0Path, err := outerArr.PathAt(0)
		require.NoError(t, err)

		inner0Slice := NewSlice[uint64](row0Path)
		_, err = inner0Slice.Grow()
		require.NoError(t, err)
		err = inner0Slice.Set(0, 123)
		require.NoError(t, err)

		val, err := inner0Slice.Get(0)
		require.NoError(t, err)
		require.Equal(t, uint64(123), val)

	})
}

func TestCache_BufferWrites(t *testing.T) {
	r := tests.Setup(t, nil)
	addr := common.HexToAddress("0x123")
	slots := AssignSlots(reflect.TypeOf(TestState{}))
	st := NewStorage(addr, r.Evm.StateDB, slots)

	// write to storage cache
	verPath := st.Field("Version")
	err := Set(verPath, uint8(99))
	require.NoError(t, err)

	// read from storage cache
	val, err := Get[uint8](verPath)
	require.NoError(t, err)
	require.Equal(t, uint8(99), val)

	// Verify StateDB is NOT updated yet
	slot := verPath.slot
	directVal := r.Evm.StateDB.GetState(addr, slot)
	require.Equal(t, common.Hash{}, directVal, "StateDB should be empty before commit")

	// commit and check stateDB
	st.Commit()
	directValAfter := r.Evm.StateDB.GetState(addr, slot)
	// We expect the value 99 (0x63) right-aligned
	expected := common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000063")
	require.Equal(t, expected, directValAfter, "StateDB should have value after commit")
}

func TestCache_ReadThrough(t *testing.T) {
	r := tests.Setup(t, nil)
	addr := common.HexToAddress("0x456")
	slots := AssignSlots(reflect.TypeOf(TestState{}))
	st := NewStorage(addr, r.Evm.StateDB, slots)

	path := st.Field("Balance")
	slot := path.slot

	// Manually write to StateDB
	balanceVal := common.BigToHash(big.NewInt(5000))
	r.Evm.StateDB.SetState(addr, slot, balanceVal)

	// Read from Storage (It should miss cache, hit StateDB and populate cache)
	bal, err := Get[Uint256](path)
	require.NoError(t, err)
	require.Equal(t, uint64(5000), bal.Uint64())

	// check internal cache structure
	cachedVal, ok := st.cache[slot]
	require.True(t, ok, "Slot should be cached after read")
	require.Equal(t, balanceVal, cachedVal)
}

func TestCache_Overwrite(t *testing.T) {
	r := tests.Setup(t, nil)
	addr := common.HexToAddress("0x789")
	slots := AssignSlots(reflect.TypeOf(TestState{}))
	st := NewStorage(addr, r.Evm.StateDB, slots)

	path := st.Field("Version")
	// write
	require.NoError(t, Set(path, uint8(10)))
	// Overwrite
	require.NoError(t, Set(path, uint8(20)))
	// Verify Cache has latest
	val, _ := Get[uint8](path)
	require.Equal(t, uint8(20), val)

	// commit and verify stateDB
	st.Commit()
	directVal := r.Evm.StateDB.GetState(addr, path.slot)
	expected := common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000014") // 20
	require.Equal(t, expected, directVal)
}

// TestBytesCache verifies that multi-slot writes (Bytes) are buffered correctly
func TestCache_Bytes(t *testing.T) {
	r := tests.Setup(t, nil)
	addr := common.HexToAddress("0xABC")
	slots := AssignSlots(reflect.TypeOf(ComplexState{})) // Has UserData []byte
	st := NewStorage(addr, r.Evm.StateDB, slots)

	path := st.Field("UserData")
	data := []byte{0xAA, 0xBB, 0xCC}

	// Write Bytes
	require.NoError(t, Set(path, data))

	// Verify StateDB is empty
	headSlot := path.slot
	require.Equal(t, common.Hash{}, r.Evm.StateDB.GetState(addr, headSlot))

	// Commit and Verify StateDB has Head Length
	st.Commit()

	headVal := r.Evm.StateDB.GetState(addr, headSlot)
	require.Equal(t, uint64(3), new(big.Int).SetBytes(headVal.Bytes()).Uint64())
}
