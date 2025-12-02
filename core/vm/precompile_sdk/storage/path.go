package storage

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"reflect"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/crypto"
)

// Path represents a path to a specific slot in the storage.
type Path struct {
	s    *Storage
	slot common.Hash
	info SlotInfo
	err  error
}

func (p *Path) Error() error {
	return p.err
}

// Map navigates to the map value for the given key, returns the path to that value.
func (p *Path) Map(key any) *Path {
	if p.err != nil {
		return p
	}

	if p.info.KeyType == nil {
		p.err = fmt.Errorf("field is not a map")
		return p
	}

	keyBytes, err := encodeTo32Bytes(key, p.info.KeyType)
	if err != nil {
		p.err = err
		return p
	}

	newSlot := crypto.Keccak256Hash(append(keyBytes, p.slot.Bytes()...))
	valType := p.info.ValueType
	newInfo := SlotInfo{
		IsDynamic:   true,
		ValueType:   p.info.ValueType,
		SubSlots:    p.info.SubSlots,
		numElements: p.info.numElements,
		Size:        0, // Since it'a dynamic root
	}

	if valType.Kind() == reflect.Map {
		newInfo.KeyType = valType.Key()    // Setup Key for the next .Map() call
		newInfo.ValueType = valType.Elem() // Peel off the map wrapper
		newInfo.IsDynamic = true
	} else if valType.Kind() == reflect.Slice {
		newInfo.ValueType = valType
		newInfo.IsDynamic = true
		newInfo.Resolver = NewResolverForType(valType, newSlot)
	} else if valType.Kind() == reflect.Array {
		newInfo.Resolver = NewResolverForType(valType, newSlot)
	}
	return &Path{
		s:    p.s,
		slot: newSlot,
		info: newInfo,
	}
}

// Index navigates to the array/slice element at the given index, returns the path to that element.
func (p *Path) Index(idx uint64) *Path {
	if p.err != nil {
		return p
	}
	// indexing is only for arrays/slices
	kind := p.info.ValueType.Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		p.err = fmt.Errorf("path is not a slice or array")
		return p
	}

	if p.info.Resolver == nil {
		p.err = fmt.Errorf("no resolver for slice/array")
		return p
	}
	slot, offset, err := p.info.Resolver.ResolveElementSlot(idx)
	if err != nil {
		p.err = err
		return p
	}
	elemType := p.info.ValueType.Elem()
	info := SlotInfo{
		Offset:    int(offset),
		Size:      int(p.info.ValueType.Elem().Size()),
		ValueType: elemType,
		SubSlots:  p.info.SubSlots,
	}

	if elemType.Kind() == reflect.Map { // if this is a map, setup key type
		info.IsDynamic = true
		info.KeyType = elemType.Key()
	}

	// If the element is a Slice (Slice of Slices), we need a new resolver
	if elemType.Kind() == reflect.Slice {
		info.IsDynamic = true
		info.Resolver = NewDynamicResolver(slot)
	}
	// always create a new instance
	return &Path{
		s:    p.s,
		slot: slot,
		info: info,
	}
}

// Field navigates to the struct field with the given name, returns the path to that field.
func (p *Path) Field(name string) *Path {
	if p.err != nil {
		return p
	}

	subInfo, ok := p.info.SubSlots[name]
	if !ok {
		p.err = fmt.Errorf("field %s not found", name)
		return p
	}
	// Resolve Relative Slot from SubSlots
	relIdx := binary.BigEndian.Uint64(subInfo.Slot[24:])

	baseBig := new(big.Int).SetBytes(p.slot.Bytes())
	newBig := new(big.Int).Add(baseBig, new(big.Int).SetUint64(relIdx))

	return &Path{
		s:    p.s,
		slot: common.BigToHash(newBig),
		info: subInfo,
	}
}

// Get retrieves the value at the given path.
func Get[T any](p *Path) (T, error) {
	var zero T // zero value to return in case of error
	if p.err != nil {
		return zero, p.err
	}
	accessor, ok := getAccessor(reflect.TypeOf(zero))
	if !ok {
		return zero, fmt.Errorf("no accessor for type %T", zero)
	}
	val, err := accessor.ReadAt(p.slot, p.info.Offset, p.s)
	return val.(T), err
}

// Set sets the value at the given path.
func Set[T any](p *Path, v T) error {
	if p.err != nil {
		return p.err
	}
	valueType := reflect.TypeOf(v)
	accessor, ok := getAccessor(valueType)
	if !ok {
		return fmt.Errorf("no accessor for type %T", v)
	}

	// update the specific part
	return accessor.WriteAt(p.slot, p.info.Offset, v, p.s)
}
