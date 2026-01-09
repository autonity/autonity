package storage

import (
	"reflect"

	"github.com/autonity/autonity/common"
)

type Array[V any, S any] struct {
	st       *Storage
	len      uint64
	baseSlot common.Hash
	elemSize uint64
}

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
	a.len = uint64(shapeType.Len())
	// determine slot consumption per element V
	a.elemSize = getSlotConsumption[V]()
	totalSlots := a.len * a.elemSize

	return addSlot(baseSlot, totalSlots), 0
}

func (a *Array[V, S]) Len() uint64 {
	return a.len
}

func (a *Array[V, S]) Get(index uint64) *V {
	if index >= a.len {
		return nil
	}
	itemSlot := addSlot(a.baseSlot, index*a.elemSize)
	val := new(V)
	BindState(a.st, itemSlot, val)
	return val
}
