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
		newInfo.ValueType = valType.Elem() // Peel off slice wrapper
		newInfo.IsDynamic = true
		newInfo.Resolver = NewResolverForType(valType, newSlot)
	} else if valType.Kind() == reflect.Array {
		// Arrays inside maps need a resolver to support .Index()
		newInfo.Resolver = NewResolverForType(valType, newSlot)
	}
	return &Path{
		s:    p.s,
		slot: newSlot,
		info: newInfo,
	}
}

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
	info := SlotInfo{
		Offset:    int(offset),
		Size:      int(p.info.ValueType.Elem().Size()),
		ValueType: p.info.ValueType.Elem(),
		SubSlots:  p.info.SubSlots,
	}
	// always create a new instance
	return &Path{
		s:    p.s,
		slot: slot,
		info: info,
	}
}

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

func Get[T any](p *Path) (T, error) {
	var zero T // zero value to return in case of error
	if p.err != nil {
		return zero, p.err
	}
	accessor, ok := getAccessor(reflect.TypeOf(zero))
	if !ok {
		return zero, fmt.Errorf("no accessor for type %T", zero)
	}
	data := p.s.stateDB.GetState(p.s.address, p.slot)
	val, err := accessor.ReadAt(data, p.info.Offset, p.s)
	return val.(T), err
}

func Set[T any](p *Path, v T) error {
	if p.err != nil {
		return p.err
	}
	valueType := reflect.TypeOf(v)
	accessor, ok := getAccessor(valueType)
	if !ok {
		return fmt.Errorf("no accessor for type %T", v)
	}
	// full slot, there is no need of extra read, just set
	if p.info.Offset == 0 && p.info.Size == 32 {
		var data common.Hash
		newSlot := accessor.WriteAt(data, 0, v, p.s)
		// normal full 32 byte write
		p.s.stateDB.SetState(p.s.address, p.slot, newSlot)
		return nil
	}

	// fetch current slot value
	current := p.s.stateDB.GetState(p.s.address, p.slot)
	// update the specific part
	newSlot := accessor.WriteAt(current, p.info.Offset, v, p.s)
	p.s.stateDB.SetState(p.s.address, p.slot, newSlot)
	return nil
}

func (p *Path) Len() (uint64, error) {
	if p.err != nil {
		return 0, p.err
	}
	if !p.info.IsDynamic || p.info.ValueType.Kind() != reflect.Slice {
		return 0, fmt.Errorf("len() called on non-dynamic array")
	}
	data := p.s.stateDB.GetState(p.s.address, p.slot)
	return binary.BigEndian.Uint64(data[:8]), nil
}

func (p *Path) setDynamicLength(newLen uint64) error {
	var headData common.Hash
	binary.BigEndian.PutUint64(headData[:8], newLen)
	p.s.stateDB.SetState(p.s.address, p.slot, headData)
	return nil
}
