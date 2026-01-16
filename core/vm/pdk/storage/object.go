package storage

import (
	"reflect"

	"github.com/autonity/autonity/common"
)

type Binder interface {
	Bind(st *Storage, baseSlot common.Hash, startOffset uint64) (nextSlot common.Hash, nextOffset uint64)
}

func BindState(st *Storage, baseSlot common.Hash, target any) uint64 {
	nextSlot, _ := bindStateRecursive(st, baseSlot, 0, target)
	return slotDiff(baseSlot, nextSlot)
}

func bindStateRecursive(st *Storage, baseSlot common.Hash, offset uint64, target any) (common.Hash, uint64) {

	if binder, ok := target.(Binder); ok {
		return binder.Bind(st, baseSlot, offset)
	}

	val := reflect.ValueOf(target)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	// always expecting a struct here
	if val.Kind() != reflect.Struct {
		return baseSlot, uint64(0)
	}
	currentSlot := baseSlot
	currentOffset := offset

	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldType := val.Type().Field(i)
		if !fieldType.IsExported() {
			continue
		}
		if binder, ok := fieldVal.Addr().Interface().(Binder); ok {
			currentSlot, currentOffset = binder.Bind(st, currentSlot, currentOffset)
			continue
		}
		if fieldVal.Kind() == reflect.Struct {
			currentSlot, currentOffset = bindStateRecursive(st, currentSlot, currentOffset, fieldVal.Addr().Interface())
		} else {
			panic("unsupported field type for binding: " + fieldVal.Type().String())
		}
	}
	return currentSlot, currentOffset
}
