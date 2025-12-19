package storage

import (
	"fmt"
	"reflect"
)

type Array[T any] struct {
	path *Path
	len  uint64
}

func NewArray[T any](p *Path) (*Array[T], error) {
	kind := p.info.ValueType.Kind()
	if kind != reflect.Array {
		return nil, fmt.Errorf("not an array: %v", kind)
	}

	return &Array[T]{
		path: p,
		len:  uint64(p.info.ValueType.Len()),
	}, nil
}

func (a *Array[T]) Len() uint64 {
	return a.len
}

func (a *Array[T]) Get(index uint64) (T, error) {
	var zero T
	if index >= a.len {
		return zero, fmt.Errorf("index %d out of range, len=%d", index, a.len)
	}
	var val T
	err := Load(a.path.Index(index), &val)
	return val, err
}

func (a *Array[T]) Set(index uint64, value T) error {
	if index >= a.len {
		return fmt.Errorf("index %d out of range, len=%d", index, a.len)
	}
	return Save(a.path.Index(index), value)
}

func (a *Array[T]) PathAt(index uint64) (*Path, error) {
	if index >= a.len {
		return nil, fmt.Errorf("index %d out of range, len=%d", index, a.len)
	}
	if a.path == nil {
		return nil, fmt.Errorf("nil path")
	}
	return a.path.Index(index), nil
}
