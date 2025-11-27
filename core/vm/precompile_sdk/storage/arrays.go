package storage

import (
	"fmt"
	"reflect"

	"github.com/autonity/autonity/common"
)

type Array[T any] struct {
	p         *Path
	elem      reflect.Type
	// for fixed size arrays
	isStatic  bool
	staticLen uint64
}

func NewArray[T any](p *Path) (*Array[T], error) {
	kind := p.info.ValueType.Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		return nil, fmt.Errorf("not an array/slice: %v", kind)
	}

	arr := &Array[T]{p: p}
	if kind == reflect.Array {
		arr.isStatic = true
		arr.staticLen = uint64(p.info.ValueType.Len())
	}
	return arr, nil
}

func (a *Array[T]) Len() (uint64, error) {
	if a.isStatic {
		return a.staticLen, nil
	}
	return a.p.Len()
}

// At is a helper to get the Path to the element at index
func (a *Array[T]) At(index uint64) *Path {
	return a.p.Index(index)
}

// GetAt gets the value at index i
func (a *Array[T]) GetAt(i uint64) (T, error) {
	var zero T
	length, err := a.Len()
	if err != nil || i >= length {
		return zero, fmt.Errorf("index %d out of range, len=%d", i, length)
	}
	return Get[T](a.p.Index(i))
}

// SetAt sets the value at index i
func (a *Array[T]) SetAt(i uint64, value T) error {
	length, err := a.Len()
	if err != nil || i >= length {
		return fmt.Errorf("index %d out of range, len=%d", i, length)
	}
	return Set[T](a.p.Index(i), value)
}

// Grow increases the array length by 1 and returns the Path to the new element.
func (a *Array[T]) Grow() (*Path, error) {
	if a.isStatic {
		return nil, fmt.Errorf("cannot grow static array")
	}
	length, err := a.Len()
	if err != nil {
		return nil, err
	}

	if err := a.p.setDynamicLength(length + 1); err != nil {
		return nil, err
	}

	return a.p.Index(length), nil
}

// Shrink decreases the array length by 1.
func (a *Array[T]) Shrink() error {
	length, err := a.Len()
	if err != nil {
		return err
	}
	// clear tail slot
	tailP := a.p.Index(length - 1)
	a.p.s.stateDB.SetState(tailP.s.address, tailP.slot, common.Hash{})

	if err := a.p.setDynamicLength(length - 1); err != nil {
		return err
	}
	return nil
}

//todo: can't support generic slices yet, we can do a slice of Paths but not of T
//// Slice - get sub range slice[start:end]
//func (a *Array[T]) Slice(start, end uint64) ([]T, error) {
//	l, err := a.Len()
//	if err != nil || end > l {
//		return nil, fmt.Errorf("slice bounds invalid (len=%d)", l)
//	}
//	res := make([]T, end-start)
//	// todo - inefficient, should be done at one shot
//	for i := start; i < end; i++ {
//		res[i-start], err = a.GetAt(i)
//		if err != nil {
//			return nil, err
//		}
//	}
//	return res, nil
//}
