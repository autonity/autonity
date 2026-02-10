package storage

import (
	"reflect"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/crypto"
)

// Map represents a key-value mapping in storage.
type Map[K any, V any] struct {
	st       *Storage
	baseSlot common.Hash
}

// Bind binds the map to a storage instance and a base slot.
func (m *Map[K, V]) Bind(st *Storage, baseSlot common.Hash, offset uint64) (common.Hash, uint64) {
	m.st = st
	if offset > 0 {
		baseSlot = addSlot(baseSlot, 1)
	}
	m.baseSlot = baseSlot
	// return next usable slot
	return addSlot(m.baseSlot, 1), 0
}

// Delete removes the key from the map and recursively clears all storage used by the value.
func (m *Map[K, V]) Delete(key K) {
	keyBytes, err := encodeTo32Bytes(key, reflect.TypeOf(key))
	if err != nil {
		panic(err)
	}
	keyHash := crypto.Keccak256Hash(append(keyBytes, m.baseSlot.Bytes()...))

	// We must bind it to the state so it knows WHERE to clear data from.
	val := new(V)
	BindState(m.st, keyHash, val)

	RecursiveClear(val)
}

// Get obtains the pointer to value Wrapper, which is bound to the key slot
// there is no need to set, because the V can be used to V.Set(...) directly
func (m *Map[K, V]) Get(key K) *V {
	keyBytes, err := encodeTo32Bytes(key, reflect.TypeOf(key))
	if err != nil {
		panic(err)
	}
	keyHash := crypto.Keccak256Hash(append(keyBytes, m.baseSlot.Bytes()...))
	val := new(V)
	BindState(m.st, keyHash, val)
	return val
}
