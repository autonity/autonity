package storage

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/holiman/uint256"

	"github.com/autonity/autonity/common"
)

type Var[T any] struct {
	st       *Storage
	baseSlot common.Hash
	offset   uint64
}

func (v *Var[T]) Bind(st *Storage, baseSlot common.Hash, offset uint64) (common.Hash, uint64) {
	var zero T
	acc, ok := getAccessor(reflect.TypeOf(zero))
	if !ok {
		panic(fmt.Errorf("no accessor for type %T", zero))
	}
	size := acc.Size()

	if uint64(size)+offset > 32 {
		baseSlot = addSlot(baseSlot, 1)
		offset = 0
	}
	v.st = st
	v.offset = offset
	v.baseSlot = baseSlot
	if offset+uint64(size) == 32 { // exactly filled the slot
		// return new slot
		return addSlot(baseSlot, 1), 0
	}
	return baseSlot, offset + uint64(size)
}

func (v *Var[T]) Get() T {
	var zero T
	acc, ok := getAccessor(reflect.TypeOf(zero))
	if !ok {
		panic(fmt.Errorf("no accessor found for type %T", zero))
	}
	val, err := acc.ReadAt(v.baseSlot, int(v.offset), v.st)
	if err != nil {
		panic(err)
	}
	return val.(T)
}

func (v *Var[T]) Set(val T) {
	typ := reflect.TypeOf(val)
	acc, ok := getAccessor(typ)
	if !ok {
		panic(fmt.Errorf("no accessor found for type %T", val))
	}
	err := acc.WriteAt(v.baseSlot, int(v.offset), val, v.st)
	if err != nil {
		panic(err)
	}
	return
}

func (v *Var[T]) Clear() {
	var zero T
	typ := reflect.TypeOf(zero)
	acc, ok := getAccessor(typ)
	if !ok {
		panic(fmt.Errorf("no accessor for type %T", zero))
	}

	if clearer, ok := acc.(ClearerAccessor); ok {
		if err := clearer.Clear(v.baseSlot, int(v.offset), v.st); err != nil {
			panic(err)
		}
		return
	}

	if err := acc.WriteAt(v.baseSlot, int(v.offset), zero, v.st); err != nil {
		panic(err)
	}
}

type Uint256 struct {
	uint256.Int
}

func NewUint256FromInt(val uint64) Uint256 {
	return Uint256{*uint256.NewInt(val)}
}

func NewUint256FromBig(val *big.Int) Uint256 {
	u := Uint256{*uint256.NewInt(0)}
	u.SetFromBig(val)
	return u
}

func (u *Uint256) Get() uint64 {
	return u.Uint64()
}
