package storage

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"reflect"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/crypto"
)

type Storage struct {
	address common.Address // storage scope, generally should be the contract address
	stateDB vm.StateDB     // state accessor
	slotMap map[string]SlotInfo
}

func NewStorage(address common.Address, stateDB vm.StateDB, slots map[string]SlotInfo) *Storage {
	return &Storage{address: address, stateDB: stateDB, slotMap: slots}
}

func (s *Storage) Field(name string) *Path {
	slotInfo, ok := s.slotMap[name]
	if !ok {
		return &Path{err: fmt.Errorf("no such slot for %s", name)}
	}
	return &Path{s: s, slot: slotInfo.Slot, info: slotInfo}
}

type ElementResolver interface {
	ResolveElementSlot(index uint64) (slot common.Hash, offset uint64, err error)
}

func NewResolverForType(typ reflect.Type, baseSlot common.Hash) ElementResolver {
	kind := typ.Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		return nil
	}
	elemSize, _ := getElemSize(typ.Elem())
	if kind == reflect.Slice {
		return NewDynamicResolver(baseSlot)
	}
	if elemSize == 0 {
		return nil // Err in caller
	}
	return NewFixedResolver(baseSlot, uint64(elemSize))
}

type FixedResolver struct {
	baseSlot common.Hash
	elemSize uint64
}

func NewFixedResolver(baseSlot common.Hash, elemSize uint64) *FixedResolver {
	return &FixedResolver{baseSlot: baseSlot, elemSize: elemSize}
}

func (fr *FixedResolver) ResolveElementSlot(index uint64) (common.Hash, uint64, error) {
	if fr.elemSize == 0 {
		return common.Hash{}, 0, fmt.Errorf("zero element size")
	}
	baseBig := new(big.Int).SetBytes(fr.baseSlot.Bytes())
	// calculate slot offset
	cumSize := index * fr.elemSize
	slotIdx := cumSize / 32
	slotBig := new(big.Int).Add(baseBig, new(big.Int).SetUint64(slotIdx))

	offset := cumSize % 32
	return common.BigToHash(slotBig), offset, nil
}

type DynamicResolver struct {
	baseSlot common.Hash
}

func NewDynamicResolver(headSlot common.Hash) *DynamicResolver {
	return &DynamicResolver{baseSlot: crypto.Keccak256Hash(headSlot.Bytes())}
}

func (dr *DynamicResolver) ResolveElementSlot(index uint64) (common.Hash, uint64, error) {
	baseBig := new(big.Int).SetBytes(dr.baseSlot.Bytes())
	slotBig := new(big.Int).Add(baseBig, new(big.Int).SetUint64(index))
	return common.BigToHash(slotBig), 0, nil // Full slot, offset 0
}

type SlotInfo struct {
	Slot      common.Hash
	Offset    int
	Size      int
	IsDynamic bool
	// for maps
	KeyType reflect.Type
	// for maps/arrays
	ValueType reflect.Type
	// for nested structs/maps/arrays
	SubSlots map[string]SlotInfo
	// number of elements for map values that are arrays
	numElements uint64
	// for arrays and slices
	Resolver ElementResolver
}

func AssignSlots(stateType reflect.Type) map[string]SlotInfo {
	if stateType.Kind() != reflect.Struct {
		panic("stateType must be a struct type")
	}

	rootLayout, _ := compileTypeLayout(stateType)
	globalSlotMap := make(map[string]SlotInfo)
	flattenLayout("", rootLayout, globalSlotMap)
	return globalSlotMap
}

// getSubSlotsForType drills down through Arrays/Slices to find the underlying
// Struct layout. This ensures map[Addr][]Profile gets the layout of Profile.
func getSubSlotsForType(typ reflect.Type) map[string]SlotInfo {
	if typ.Kind() == reflect.Struct {
		layout, _ := compileTypeLayout(typ)
		return layout
	}
	if typ.Kind() == reflect.Array || typ.Kind() == reflect.Slice || typ.Kind() == reflect.Map {
		return getSubSlotsForType(typ.Elem())
	}
	return nil
}

func compileTypeLayout(typ reflect.Type) (map[string]SlotInfo, uint64) {
	layout := make(map[string]SlotInfo)
	currentSlot := uint64(0)
	currentOffset := 0
	ensureNewSlot := func() {
		if currentOffset > 0 {
			currentSlot++
			currentOffset = 0
		}
	}

	addField := func(name string, fieldType reflect.Type) {
		if size, ok := getPrimitiveSize(fieldType); ok {
			if currentOffset+size > 32 {
				currentSlot++
				currentOffset = 0
			}
			layout[name] = SlotInfo{
				Slot:      computeSlotHash(currentSlot),
				Offset:    32 - (currentOffset + size), // right aligned
				Size:      size,
				ValueType: fieldType,
			}
			currentOffset += size
			return
		}

		if isDynamicType(fieldType) {
			ensureNewSlot()
			elemT := fieldType.Elem()
			info := SlotInfo{
				Slot:      computeSlotHash(currentSlot),
				Offset:    currentOffset,
				ValueType: elemT,
				IsDynamic: true,
				Size:      32, // full slot for dynamic type
			}

			if fieldType.Kind() == reflect.Map {
				info.KeyType = fieldType.Key()
				info.ValueType = elemT
			} else if fieldType.Kind() == reflect.Slice {
				info.ValueType = fieldType // keep the slice type
			}
			// recurse to get the underlying type
			info.SubSlots = getSubSlotsForType(elemT)
			// attach resovler
			info.Resolver = NewResolverForType(fieldType, info.Slot)
			layout[name] = info
			currentSlot++
			return
		}

		// static composite types structs/arrays
		if fieldType.Kind() == reflect.Struct || fieldType.Kind() == reflect.Array {
			ensureNewSlot()
			var subLayOut map[string]SlotInfo
			var slotsUsed uint64

			if fieldType.Kind() == reflect.Struct {
				subLayOut, slotsUsed = compileTypeLayout(fieldType)
			} else {
				// For Arrays, we need the Layout of the ELEMENT, not the array itself
				// But we need the Size of the ARRAY
				subLayOut, _ = compileTypeLayout(fieldType.Elem()) // Get Element Layout for SubSlots

				// Calculate Array Size
				elemSize, _ := getElemSize(fieldType.Elem())
				totalSize := elemSize * fieldType.Len()
				slotsUsed = uint64((totalSize + 31) / 32)
			}
			info := SlotInfo{
				Slot:      computeSlotHash(currentSlot),
				ValueType: fieldType,
				SubSlots:  subLayOut,
			}
			if fieldType.Kind() == reflect.Array {
				info.Resolver = NewResolverForType(fieldType, info.Slot)
			}
			layout[name] = info
			currentSlot += slotsUsed
			// empty struct edge case
			if slotsUsed == 0 && fieldType.Kind() == reflect.Struct && fieldType.NumField() > 0 {
				currentSlot++
			}
			return
		}
		panic(fmt.Sprintf("unsupported field type: %v", fieldType))
	}
	if typ.Kind() == reflect.Struct {
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			if !field.IsExported() {
				continue
			}
			addField(field.Name, field.Type)
		}
	}
	if typ.Kind() == reflect.Array {
		elem := typ.Elem()
		elemSize, _ := getElemSize(elem)
		totalSize := elemSize * typ.Len()
		slots := (totalSize + 31) / 32
		currentSlot += uint64(slots)
		if totalSize%32 != 0 {
			currentOffset = int(totalSize % 32)
		}
	}
	if currentOffset > 0 {
		currentSlot++
	}
	return layout, currentSlot
}

// todo: optimize flattening
func flattenLayout(prefix string, layout map[string]SlotInfo, target map[string]SlotInfo) {
	for name, info := range layout {
		fullName := name
		if prefix != "" {
			fullName = prefix + "." + name
		}
		target[fullName] = info
		if info.SubSlots != nil && !info.IsDynamic && info.ValueType.Kind() == reflect.Struct {
			flattenWithBase(fullName, info.Slot, info.SubSlots, target)
		}
	}
}

func flattenWithBase(prefix string, baseSlot common.Hash, subLayout map[string]SlotInfo, target map[string]SlotInfo) {
	baseBig := new(big.Int).SetBytes(baseSlot.Bytes())
	for name, relInfo := range subLayout {
		fullName := prefix + "." + name
		relVal := binary.BigEndian.Uint64(relInfo.Slot[24:])
		absBig := new(big.Int).Add(baseBig, new(big.Int).SetUint64(relVal))
		absInfo := relInfo
		absInfo.Slot = common.BigToHash(absBig)
		target[fullName] = absInfo
		if relInfo.SubSlots != nil && !relInfo.IsDynamic && relInfo.ValueType.Kind() == reflect.Struct {
			flattenWithBase(fullName, absInfo.Slot, relInfo.SubSlots, target)
		}
	}
}
