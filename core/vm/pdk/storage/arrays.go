package storage

import (
	"encoding/binary"
	"fmt"
	"reflect"

	"github.com/autonity/autonity/common"
)

type Array[T any] struct {
	path *Path
	// for fixed size arrays
	isStatic  bool
	staticLen uint64
}

func NewArray[T any](p *Path) (*Array[T], error) {
	kind := p.info.ValueType.Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		return nil, fmt.Errorf("not an array/slice: %v", kind)
	}

	arr := &Array[T]{path: p}
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
	data := a.path.st.GetState(a.path.slot)
	return binary.BigEndian.Uint64(data[:8]), nil
}

// ReferenceAt is a helper to get the Path to the element at index
func (a *Array[T]) ReferenceAt(index uint64) *Path {
	return a.path.Index(index)
}

// ValueAt gets the value at index i
func (a *Array[T]) ValueAt(i uint64) (T, error) {
	var zero T
	length, err := a.Len()
	if err != nil || i >= length {
		return zero, fmt.Errorf("index %d out of range, len=%d", i, length)
	}
	return Get[T](a.path.Index(i))
}

// SetValueAt sets the value at index i
func (a *Array[T]) SetValueAt(i uint64, value T) error {
	length, err := a.Len()
	if err != nil || i >= length {
		return fmt.Errorf("index %d out of range, len=%d", i, length)
	}
	return Set[T](a.path.Index(i), value)
}

// Grow increases the array length by 1 and returns the Path to the new element. doesn't support 2D Slices.
func (a *Array[T]) Grow() (*Path, error) {
	if a.isStatic {
		return nil, fmt.Errorf("cannot grow static array")
	}
	length, err := a.Len()
	if err != nil {
		return nil, err
	}

	if err := a.setDynamicLength(length + 1); err != nil {
		return nil, err
	}

	return a.path.Index(length), nil
}

// Shrink decreases the array length by 1, doesn't support 2D Slices.
func (a *Array[T]) Shrink() error {
	if a.isStatic {
		return fmt.Errorf("cannot shrink static array")
	}
	length, err := a.Len()
	if err != nil {
		return err
	}
	// clear tail slot
	tailP := a.path.Index(length - 1)
	a.path.st.SetState(tailP.slot, common.Hash{})

	if err := a.setDynamicLength(length - 1); err != nil {
		return err
	}
	return nil
}

func (a *Array[T]) setDynamicLength(newLen uint64) error {
	var headData common.Hash
	binary.BigEndian.PutUint64(headData[:8], newLen)
	a.path.st.SetState(a.path.slot, headData)
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
