package storage

import (
	"reflect"

	"github.com/autonity/autonity/common"
)

// Array represents a fixed-size array in storage.
type Array[V any, S any] struct {
	st       *Storage
	len      uint64
	baseSlot common.Hash
	elemSize uint64
}

// Bind binds the array to a storage instance and a base slot.
func (a *Array[V, S]) Bind(st *Storage, baseSlot common.Hash, offset uint64) (common.Hash, uint64) {
	a.st = st
	// new slot for all arrays
	if offset > 0 {
		baseSlot = addSlot(baseSlot, 1)
	}
	a.baseSlot = baseSlot
	// determine length of the array from shape S
	var shape S
	shapeType := reflect.TypeOf(shape)
	if shapeType.Kind() != reflect.Array {
		panic("shape must be an array")
	}
	l := shapeType.Len()
	if l < 0 {
		l = 0
	}
	a.len = uint64(l) // #nosec G115
	// determine slot consumption per element V
	a.elemSize = getSlotConsumption[V]()
	totalSlots := a.len * a.elemSize

	return addSlot(baseSlot, totalSlots), 0
}

// Len returns the length of the array.
func (a *Array[V, S]) Len() uint64 {
	return a.len
}

// Get returns a pointer to the element at the given index.
func (a *Array[V, S]) Get(index uint64) *V {
	if index >= a.len {
		panic("array index out of bounds")
	}
	itemSlot := addSlot(a.baseSlot, index*a.elemSize)
	val := new(V)
	BindState(a.st, itemSlot, val)
	return val
}
