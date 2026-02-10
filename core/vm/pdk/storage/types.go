package storage

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/holiman/uint256"

	"github.com/autonity/autonity/common"
)

// Var represents a single value in storage.
type Var[T any] struct {
	st       *Storage
	baseSlot common.Hash
	offset   uint64
}

// Bind binds the variable to a storage instance and a base slot.
func (v *Var[T]) Bind(st *Storage, baseSlot common.Hash, offset uint64) (common.Hash, uint64) {
	var zero T
	acc, ok := getAccessor(reflect.TypeOf(zero))
	if !ok {
		panic(fmt.Errorf("no accessor for type %T", zero))
	}
	size := acc.Size()

	uSize := uint64(size) // #nosec G115
	if uSize+offset > 32 {
		baseSlot = addSlot(baseSlot, 1)
		offset = 0
	}
	v.st = st
	v.offset = offset
	v.baseSlot = baseSlot
	if offset+uSize == 32 { // exactly filled the slot
		// return new slot
		return addSlot(baseSlot, 1), 0
	}
	return baseSlot, offset + uSize
}

// Get retrieves the value of the variable.
func (v *Var[T]) Get() T {
	var zero T
	acc, ok := getAccessor(reflect.TypeOf(zero))
	if !ok {
		panic(fmt.Errorf("no accessor found for type %T", zero))
	}
	val, err := acc.ReadAt(v.baseSlot, int(v.offset), v.st) // #nosec G115
	if err != nil {
		panic(err)
	}
	return val.(T)
}

// Set sets the value of the variable.
func (v *Var[T]) Set(val T) {
	typ := reflect.TypeOf(val)
	acc, ok := getAccessor(typ)
	if !ok {
		panic(fmt.Errorf("no accessor found for type %T", val))
	}
	err := acc.WriteAt(v.baseSlot, int(v.offset), val, v.st) // #nosec G115
	if err != nil {
		panic(err)
	}
}

// Clear removes the value of the variable from storage.
func (v *Var[T]) Clear() {
	var zero T
	typ := reflect.TypeOf(zero)
	acc, ok := getAccessor(typ)
	if !ok {
		panic(fmt.Errorf("no accessor for type %T", zero))
	}

	if clearer, ok := acc.(ClearerAccessor); ok {
		// #nosec G115
		if err := clearer.Clear(v.baseSlot, int(v.offset), v.st); err != nil {
			panic(err)
		}
		return
	}

	// #nosec G115
	if err := acc.WriteAt(v.baseSlot, int(v.offset), zero, v.st); err != nil {
		panic(err)
	}
}

// Uint256 represents a 256-bit unsigned integer.
type Uint256 struct {
	uint256.Int
}

// NewUint256FromInt creates a new Uint256 from a uint64 value.
func NewUint256FromInt(val uint64) Uint256 {
	return Uint256{*uint256.NewInt(val)}
}

// NewUint256FromBig creates a new Uint256 from a *big.Int value.
func NewUint256FromBig(val *big.Int) Uint256 {
	u := Uint256{*uint256.NewInt(0)}
	u.SetFromBig(val)
	return u
}

// Get returns the uint64 representation of the Uint256.
func (u *Uint256) Get() uint64 {
	return u.Uint64()
}
