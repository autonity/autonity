package storage

import (
	"reflect"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/crypto"
)

type Map[K any, V any] struct {
	st       *Storage
	baseSlot common.Hash
}

func (m *Map[K, V]) Bind(st *Storage, baseSlot common.Hash, offset uint64) (common.Hash, uint64) {
	m.st = st
	if offset > 0 {
		baseSlot = addSlot(baseSlot, 1)
	}
	m.baseSlot = baseSlot
	// return next usable slot
	return addSlot(m.baseSlot, 1), offset
}

// Get obtains the pointer to value Wrapper, which is bound to the key slot
// there is no need ot set, because the V can be used to V.Set(...) directly
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
