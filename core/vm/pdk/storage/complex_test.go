package storage_test

import (
	"math/big"
	"testing"

	"github.com/autonity/autonity/autonity/tests"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm/pdk/events"
	"github.com/autonity/autonity/core/vm/pdk/storage"
	"github.com/stretchr/testify/require"
)

// --- Complex Storage Structures ---

type SliceInMap struct {
	// Map: Address -> Slice[uint64]
	Data storage.Map[common.Address, storage.Slice[storage.Var[uint64]]]
}

type MapInSlice struct {
	// Slice of Map: uint64 -> uint64
	History storage.Slice[storage.Map[uint64, storage.Var[uint64]]]
}

type NestedArrays struct {
	// Array of Array: [2][3]uint64
	Matrix storage.Array[storage.Array[storage.Var[uint64], [3]any], [2]any]
}

// --- Event Structures ---

type EventWithArray struct {
	User   common.Address `indexed:"true"`
	Scores [3]uint64
}

type EventWithBigIntArray struct {
	ID    uint64 `indexed:"true"`
	Rates [2]*big.Int
}

// --- Tests ---

func TestStorage_SliceInMap(t *testing.T) {
	r := tests.Setup(t, nil)
	addr := common.HexToAddress("0x1111")
	st := storage.NewStorage(addr, r.Evm.StateDB)
	
	state := new(SliceInMap)
	storage.BindState(st, common.Hash{}, state)

	user := common.HexToAddress("0xAAAA")

	// Append to slice inside map
	slice := state.Data.Get(user)
	slice.Append(func(v *storage.Var[uint64]) {
		v.Set(100)
	})
	slice.Append(func(v *storage.Var[uint64]) {
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
	st := storage.NewStorage(addr, r.Evm.StateDB)

	state := new(MapInSlice)
	storage.BindState(st, common.Hash{}, state)

	// Create a new map in the slice
	state.History.Append(func(m *storage.Map[uint64, storage.Var[uint64]]) {
		m.Get(1).Set(10)
		m.Get(2).Set(20)
	})

	// Create another map
	state.History.Append(func(m *storage.Map[uint64, storage.Var[uint64]]) {
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

func TestEvent_WithArrays(t *testing.T) {
	r := tests.Setup(t, nil)
	addr := common.HexToAddress("0x3333")
	st := storage.NewStorage(addr, r.Evm.StateDB)

	// Test EventWithArray ([3]uint64)
	evt1 := EventWithArray{
		User:   common.HexToAddress("0xUSER"),
		Scores: [3]uint64{10, 20, 30},
	}
	
	// This relies on events.Emit using abiselector.ResolveABIType which now supports arrays
	err := events.Emit(st, evt1)
	if err != nil {
		t.Fatalf("Failed to emit event with array: %v", err)
	}
	
	// Test EventWithBigIntArray ([2]*big.Int)
	evt2 := EventWithBigIntArray{
		ID:    1,
		Rates: [2]*big.Int{big.NewInt(1000), big.NewInt(2000)},
	}
	
	err = events.Emit(st, evt2)
	if err != nil {
		t.Fatalf("Failed to emit event with big.Int array: %v", err)
	}
}
