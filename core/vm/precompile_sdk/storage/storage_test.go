package storage

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/autonity/tests"
	"github.com/autonity/autonity/common"
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
	Address Address
	Balance Uint256

	// nested struct
	AdminProfile Profile

	// dynamic array/slice
	History []common.Hash

	// map with struct
	UserStats map[Address]DeepCommission

	// map with primitive arrays
	Scores map[Address][]uint32

	Category   map[Address]uint64
	UserArrays map[Address][]Profile
}

var (
	testContractAddr = common.HexToAddress("0xCAFECAFECAFECAFECAFECAFECAFECAFECAFECAFE")
	testUserAddr     = NewAddressFromBytes(common.HexToAddress("0x1134567890123456789012345678901234567890").Bytes())
	testValidator    = NewAddressFromBytes(common.HexToAddress("0x1234567890123456789012345678901234567890").Bytes())
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

	err = Set(st.Field("Address"), testUserAddr)
	require.NoError(t, err)
	addr, err := Get[Address](st.Field("Address"))
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
	mapPath := st.Field("UserArrays").Map(testAddr) // Resolves keccak(encode(addr) || base_slot)
	arr, err := NewArray[Profile](mapPath)          // Wraps for push
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
	updatedPath := mapPath.Index(0).Field("Config").Field("IsEnabled") // resolveSliceElemSlot → struct rel offset
	err = Set[bool](updatedPath, false)                                // Accessor overlays bool at offset (byte 4 in packed struct, right-aligned)
	require.NoError(t, err)

	// Get: Chain read back, verify
	_, err = arr.Len()
	require.NoError(t, err)

	// Nested get: Direct field chain
	levelPath := st.Field("UserArrays").Map(testAddr).Index(0).Field("Level")
	gotLevel, err := Get[uint64](levelPath) // uintAccessor.ReadAt(struct_slot, level_offset=0, st)
	require.NoError(t, err)
	require.Equal(t, uint64(42), gotLevel)

	// Pop: Remove last (clears slot for refund, dec len)
	err = arr.Shrink()
	require.NoError(t, err)

}

// --- New Complex State for Advanced Tests ---

type ComplexState struct {
	MapBool   map[bool]uint64
	MapUint8  map[uint8]string // Using []byte/string if supported, else uint64
	MapUint64 map[uint64]Profile

	// Mapping Address -> (ID -> Active Status)
	UserPermissions map[Address]map[uint64]bool

	// A list of historical snapshots, each snapshot is a map of ID->Value
	Snapshots []map[uint64]uint64

	SnapshotsFixedMapping [10]map[uint64]uint64
	SnapshotsFixedStruct  [10]Profile
	// 4 * 8 bytes = 32 bytes (Fits in 1 slot)
	PackedRates [3]uint64
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

// --- Test 2: Nested Maps ---
func TestNestedMaps(t *testing.T) {
	st := newComplexStorage(t)

	// Solidity: mapping(address => mapping(uint64 => bool))
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

	snapshotsArr, err := NewArray[map[uint64]uint64](st.Field("Snapshots")) // Type T doesn't matter much for Next/Map calls
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

	// Solidity: mapping(uint64 => uint64)[] snapshots;
	snapshotsArr, err := NewArray[map[uint64]uint64](st.Field("Snapshots")) // Type T doesn't matter much for Next/Map calls
	require.NoError(t, err)

	// 1. Add a new Snapshot (Index 0)
	// Next() returns path to Snapshots[0] (which is a Map root)
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

// --- Test 4: Packed Fixed Arrays ---
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
