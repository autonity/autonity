package storage

import (
	"encoding/binary"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/autonity/tests"
	"github.com/autonity/autonity/common"
)

type ConfigState struct {
	TimeoutSeconds Var[uint32]
	IsEnabled      Var[bool]
	Padding        Var[[3]byte]
}

type ProfileState struct {
	NameHash Var[common.Hash]
	Level    Var[uint64]
	Config   ConfigState // Nested Struct Wrapper
}

type DeepCommissionState struct {
	MaxRate      Var[uint64]
	LastModified Var[int64]
	Info         ConfigState
}

// TestState mirrors the old struct but with Wrappers
type TestState struct {
	// Primitive types
	Version Var[uint8]
	Address Var[common.Address]
	Balance Var[Uint256]

	// Nested struct
	AdminProfile ProfileState

	// Dynamic array/slice of Primitives
	History Slice[Var[common.Hash]]

	// Map with Struct Value
	UserStats Map[common.Address, DeepCommissionState]

	// Map with Slice Value
	Scores Map[common.Address, Slice[Var[uint32]]]

	// Map with Primitive Value
	Category Map[common.Address, Var[uint64]]

	// Map with Slice of Structs
	UserArrays Map[common.Address, Slice[ProfileState]]
}

var (
	testContractAddr = common.HexToAddress("0xCAFECAFECAFECAFECAFECAFECAFECAFECAFECAFE")
	testUserAddr     = common.HexToAddress("0x1134567770123456777012345677701234567770")
	testValidator    = common.HexToAddress("0x1234567770123456777012345677701234567770")
)

func newStorage(t *testing.T) (*Storage, *TestState) {
	r := tests.Setup(t, nil)
	st := NewStorage(testContractAddr, r.Evm.StateDB)

	state := new(TestState)
	BindState(st, common.Hash{}, state)
	return st, state
}

func TestStoragePrimitives(t *testing.T) {
	_, state := newStorage(t)

	state.Version.Set(42)
	v := state.Version.Get()
	require.Equal(t, uint8(42), v)

	state.Address.Set(testUserAddr)
	addr := state.Address.Get()
	require.Equal(t, testUserAddr, addr)

	testBalance := NewUint256FromInt(1000000)
	state.Balance.Set(testBalance)
	bal := state.Balance.Get()
	require.Equal(t, testBalance, bal)
}

func TestStorageNestedStruct(t *testing.T) {
	_, state := newStorage(t)
	testTimeout := uint32(300)

	state.AdminProfile.Config.TimeoutSeconds.Set(testTimeout)

	timeout := state.AdminProfile.Config.TimeoutSeconds.Get()
	require.Equal(t, testTimeout, timeout)

	testLevel := uint64(300)
	state.AdminProfile.Level.Set(testLevel)

	level := state.AdminProfile.Level.Get()
	require.Equal(t, testLevel, level)
}

func TestStorageMapWithPrimitives(t *testing.T) {
	_, state := newStorage(t)
	testVal := uint64(95)

	state.Category.Get(testValidator).Set(testVal)

	val := state.Category.Get(testValidator).Get()
	require.Equal(t, testVal, val)
}

func TestStorageMapWithPrimitiveSlice(t *testing.T) {
	_, state := newStorage(t)
	testVal := uint32(95)

	scoresSlice := state.Scores.Get(testValidator)

	// Append using the Callback Pattern
	scoresSlice.Append(func(item *Var[uint32]) {
		item.Set(testVal)
	})

	length := scoresSlice.Len()
	require.Equal(t, uint64(1), length)

	// Get element at index 0
	val := scoresSlice.Get(0).Get()
	require.Equal(t, testVal, val)
}

func TestMapOfAddressToStructArrays(t *testing.T) {
	_, state := newStorage(t)

	profileSlice := state.UserArrays.Get(testUserAddr)

	profileSlice.Append(func(p *ProfileState) {
		p.Level.Set(42)
		p.Config.TimeoutSeconds.Set(300)
		p.Config.IsEnabled.Set(true)
	})

	length := profileSlice.Len()
	require.Equal(t, uint64(1), length)

	prof := profileSlice.Get(0)
	prof.Config.IsEnabled.Set(false)

	profCheck := profileSlice.Get(0) // Re-bind (cheap)

	lvl := profCheck.Level.Get()
	require.Equal(t, uint64(42), lvl)

	enabled := profCheck.Config.IsEnabled.Get()
	require.Equal(t, false, enabled)
}

type ComplexState struct {
	MapBool   Map[bool, Var[uint64]]
	MapUint8  Map[uint8, Var[string]]
	MapUint64 Map[uint64, ProfileState]

	// Nested Maps: Address -> ID -> Bool
	UserPermissions Map[common.Address, Map[uint64, Var[bool]]]

	// Slice of Maps
	Snapshots Slice[Map[uint64, Var[uint64]]]

	// Fixed Array of Maps [10]...
	SnapshotsFixedMapping Array[Map[uint64, Var[uint64]], [10]any]

	// Fixed Array of Structs [10]...
	SnapshotsFixedStruct Array[ProfileState, [10]any]

	// Packed Rates [3]uint64
	PackedRates Array[Var[uint64], [3]any]

	UserData Var[[]byte]
}

func newComplexStorage(t *testing.T) (*Storage, *ComplexState) {
	r := tests.Setup(t, nil)
	st := NewStorage(testContractAddr, r.Evm.StateDB)
	state := new(ComplexState)
	BindState(st, common.Hash{}, state)
	return st, state
}

func TestMapKeyVariations(t *testing.T) {
	_, state := newComplexStorage(t)

	state.MapBool.Get(true).Set(100)
	state.MapBool.Get(false).Set(200)

	valTrue := state.MapBool.Get(true).Get()
	require.Equal(t, uint64(100), valTrue)

	valFalse := state.MapBool.Get(false).Get()
	require.Equal(t, uint64(200), valFalse)

	key := uint64(9999)

	profile := state.MapUint64.Get(key)
	profile.Level.Set(10)

	p1 := state.MapUint64.Get(key)
	lvl := p1.Level.Get()
	require.Equal(t, uint64(10), lvl)
}

func TestNestedMaps(t *testing.T) {
	_, state := newComplexStorage(t)
	user := testUserAddr
	resourceID := uint64(42)

	state.UserPermissions.Get(user).Get(resourceID).Set(true)

	gotPerm := state.UserPermissions.Get(user).Get(resourceID).Get()
	require.True(t, gotPerm)

	gotOther := state.UserPermissions.Get(user).Get(99).Get()
	require.False(t, gotOther)
}

func TestSliceOfMaps(t *testing.T) {
	_, state := newComplexStorage(t)

	state.Snapshots.Append(func(m *Map[uint64, Var[uint64]]) {
		m.Get(100).Set(500)
	})

	state.Snapshots.Append(func(m *Map[uint64, Var[uint64]]) {
		m.Get(100).Set(900)
	})

	val0 := state.Snapshots.Get(0).Get(100).Get()
	require.Equal(t, uint64(500), val0)

	val1 := state.Snapshots.Get(1).Get(100).Get()
	require.Equal(t, uint64(900), val1)
}

func TestPackedFixedArray(t *testing.T) {
	_, state := newComplexStorage(t)

	state.PackedRates.Get(0).Set(10)

	state.PackedRates.Get(1).Set(20)

	v0 := state.PackedRates.Get(0).Get()
	v1 := state.PackedRates.Get(1).Get()

	require.Equal(t, uint64(10), v0)
	require.Equal(t, uint64(20), v1)

	require.Equal(t, uint64(3), state.PackedRates.Len())
}

func TestByteAccessor(t *testing.T) {
	st, state := newComplexStorage(t)

	bytesVar := state.UserData

	t.Run("EmptyBytes", func(t *testing.T) {
		empty := []byte{}
		bytesVar.Set(empty)

		got := bytesVar.Get()
		require.Equal(t, empty, got)

		headData := st.GetState(bytesVar.baseSlot)
		require.Equal(t, uint64(0), binary.BigEndian.Uint64(headData[24:]))
	})

	t.Run("MediumBytesTwoChunks", func(t *testing.T) {
		data := make([]byte, 40)
		copy(data, "this is a medium byte array that spans slots exactly")
		bytesVar.Set(data)

		got := bytesVar.Get()
		require.Equal(t, data, got)
	})

	t.Run("OverwriteResize", func(t *testing.T) {
		initial := make([]byte, 50)
		copy(initial, "initial data longer than 32")
		bytesVar.Set(initial)

		updated := []byte("shorter data")
		bytesVar.Set(updated)

		got := bytesVar.Get()
		require.Equal(t, updated, got)
	})
}

func Test2DArrays(t *testing.T) {
	t.Run("Primitive2D", func(t *testing.T) {
		// [2][3]uint64
		type Primitive2D struct {
			Data Array[Array[Var[uint64], [3]any], [2]any]
		}

		r := tests.Setup(t, nil)
		st := NewStorage(testContractAddr, r.Evm.StateDB)
		state := new(Primitive2D)
		BindState(st, common.Hash{}, state)

		state.Data.Get(0).Get(0).Set(5)
		state.Data.Get(0).Get(1).Set(2)
		state.Data.Get(1).Get(2).Set(6)

		v00 := state.Data.Get(0).Get(0).Get()
		require.Equal(t, uint64(5), v00)

		v12 := state.Data.Get(1).Get(2).Get()
		require.Equal(t, uint64(6), v12)
	})

	t.Run("Dynamic2D", func(t *testing.T) {
		type Dynamic2D struct {
			Data Slice[Slice[Var[uint64]]]
		}
		r := tests.Setup(t, nil)
		st := NewStorage(testContractAddr, r.Evm.StateDB)
		state := new(Dynamic2D)
		BindState(st, common.Hash{}, state)

		state.Data.Append(func(row *Slice[Var[uint64]]) {
			row.Append(func(v *Var[uint64]) { v.Set(10) })
			row.Append(func(v *Var[uint64]) { v.Set(20) })
		})

		state.Data.Append(func(row *Slice[Var[uint64]]) {
			row.Append(func(v *Var[uint64]) { v.Set(30) })
		})

		val01 := state.Data.Get(0).Get(1).Get()
		require.Equal(t, uint64(20), val01)

		val10 := state.Data.Get(1).Get(0).Get()
		require.Equal(t, uint64(30), val10)
	})
}

func TestCache_BufferWrites(t *testing.T) {
	st, state := newStorage(t)

	state.Version.Set(99)

	val := state.Version.Get()
	require.Equal(t, uint8(99), val)

	slot := state.Version.baseSlot
	directVal := st.stateDB.GetState(st.address, slot)
	require.Equal(t, common.Hash{}, directVal, "StateDB should be empty before commit")

	st.Commit()
	directValAfter := st.stateDB.GetState(st.address, slot)

	var expected common.Hash
	expected[0] = 0x63
	require.Equal(t, expected, directValAfter)
}

type BigIntState struct {
	Counter   Var[*big.Int]
	Balances  Map[common.Address, Var[*big.Int]]
	Historics Slice[Var[*big.Int]]
}

func TestBigIntAccessor(t *testing.T) {
	r := tests.Setup(t, nil)
	st := NewStorage(testContractAddr, r.Evm.StateDB)
	state := new(BigIntState)
	BindState(st, common.Hash{}, state)

	t.Run("Positive Values", func(t *testing.T) {
		val := big.NewInt(1000)
		state.Counter.Set(val)

		got := state.Counter.Get()
		require.Equal(t, 0, val.Cmp(got))
	})

	t.Run("Negative Values", func(t *testing.T) {
		val := big.NewInt(-1)
		state.Counter.Set(val)

		got := state.Counter.Get()
		require.Equal(t, 0, val.Cmp(got))
	})

	t.Run("Map Integration", func(t *testing.T) {
		val := big.NewInt(-5000)

		state.Balances.Get(testUserAddr).Set(val)

		got := state.Balances.Get(testUserAddr).Get()
		require.Equal(t, 0, val.Cmp(got))
	})

	t.Run("Slice Integration", func(t *testing.T) {
		val1 := big.NewInt(100)

		state.Historics.Append(func(v *Var[*big.Int]) {
			v.Set(val1)
		})

		got1 := state.Historics.Get(0).Get()
		require.Equal(t, 0, val1.Cmp(got1))
	})
}

type SliceInMap struct {
	// Map: Address -> Slice[uint64]
	Data Map[common.Address, Slice[Var[uint64]]]
}

type MapInSlice struct {
	// Slice of Map: uint64 -> uint64
	History Slice[Map[uint64, Var[uint64]]]
}


func TestStorage_SliceInMap(t *testing.T) {
	r := tests.Setup(t, nil)
	addr := common.HexToAddress("0x1111")
	st := NewStorage(addr, r.Evm.StateDB)

	state := new(SliceInMap)
	BindState(st, common.Hash{}, state)

	user := common.HexToAddress("0xAAAA")

	// Append to slice inside map
	slice := state.Data.Get(user)
	slice.Append(func(v *Var[uint64]) {
		v.Set(100)
	})
	slice.Append(func(v *Var[uint64]) {
		v.Set(200)
	})

	// Verify length
	require.Equal(t, uint64(2), slice.Len())

	// Verify elements
	require.Equal(t, uint64(100), slice.Get(0).Get())
	require.Equal(t, uint64(200), slice.Get(1).Get())

	// Verify isolation (another user)
	otherUser := common.HexToAddress("0xBBBB")
	otherSlice := state.Data.Get(otherUser)
	require.Equal(t, uint64(0), otherSlice.Len())
}

func TestStorage_MapInSlice(t *testing.T) {
	r := tests.Setup(t, nil)
	addr := common.HexToAddress("0x2222")
	st := NewStorage(addr, r.Evm.StateDB)

	state := new(MapInSlice)
	BindState(st, common.Hash{}, state)

	// Create a new map in the slice
	state.History.Append(func(m *Map[uint64, Var[uint64]]) {
		m.Get(1).Set(10)
		m.Get(2).Set(20)
	})

	// Create another map
	state.History.Append(func(m *Map[uint64, Var[uint64]]) {
		m.Get(1).Set(99) // Different value for same key
	})

	// Verify Map 0
	map0 := state.History.Get(0)
	require.Equal(t, uint64(10), map0.Get(1).Get())
	require.Equal(t, uint64(20), map0.Get(2).Get())

	// Verify Map 1
	map1 := state.History.Get(1)
	require.Equal(t, uint64(99), map1.Get(1).Get())
	require.Equal(t, uint64(0), map1.Get(2).Get()) // Should be empty
}
