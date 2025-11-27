package storage

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"reflect"
	"strings"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/crypto"
)

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

func AssignSlots(stateType reflect.Type) map[string]SlotInfo {
	if stateType.Kind() != reflect.Struct {
		panic("stateType must be a struct type")
	}

	rootLayout, _ := compileTypeLayout(stateType)
	globalSlotMap := make(map[string]SlotInfo)
	flattenLayout("", rootLayout, globalSlotMap)
	return globalSlotMap
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

func handlePrimitiveType(fName string, size int, currentSlot *uint64,
	currentOffset *int, slotMap map[string]SlotInfo) error {
	if *currentOffset+size > 32 { // doesn't fit in current slot
		*currentSlot++
		*currentOffset = 0
	}
	// always right align primitive types
	rightAlignedOffset := 32 - (*currentOffset + size)
	slotInfo := SlotInfo{
		Slot:        computeSlotHash(*currentSlot),
		Offset:      rightAlignedOffset,
		Size:        size,
		IsDynamic:   false,
		numElements: 1,
	}
	slotMap[fName] = slotInfo
	*currentOffset += size
	return nil
}

func handleDynamicType(fName string, fType reflect.Type, currentSlot *uint64, currentOffset *int, slotMap map[string]SlotInfo) error {
	// dynamic type will start from the new slot, increment slot and reset offset
	if *currentOffset > 0 {
		*currentSlot++
		*currentOffset = 0
	}
	baseSlot := computeSlotHash(*currentSlot)
	info := SlotInfo{
		Slot:      baseSlot,
		Offset:    *currentOffset,
		Size:      32, // take full slot
		IsDynamic: true,
	}
	if fType.Kind() == reflect.Map {
		info.KeyType = fType.Key()
		info.ValueType = fType.Elem()
	} else if fType.Kind() == reflect.Slice {
		info.ValueType = fType.Elem()
		info.Resolver = NewDynamicResolver(baseSlot)
		if info.Resolver == nil {
			panic("failed to create dynamic resolver")
		}
	}

	if info.ValueType.Kind() == reflect.Struct || info.ValueType.Kind() == reflect.Array {
		info.SubSlots, info.numElements = computeSubSlots(info.ValueType)
	}
	slotMap[fName] = info
	*currentSlot++
	*currentOffset = 0
	return nil
}

func computeSubSlots(typ reflect.Type) (map[string]SlotInfo, uint64) {
	subSlotMap := make(map[string]SlotInfo)
	relSlot := uint64(0)
	relSlotOffset := int(0)
	err := assignSlotsRecursive(typ, "", &relSlot, &relSlotOffset, subSlotMap)
	if err != nil {
		panic(err)
	}
	numElement := relSlot
	if relSlotOffset > 0 {
		// partially filled slot counts as one
		numElement++
	}
	// edge case: empty struct/arrays
	if numElement == 0 && (typ.Kind() == reflect.Struct || typ.Kind() == reflect.Array) {
		numElement = 1
	}
	return subSlotMap, numElement
}

func processField(name string, typ reflect.Type, currentSlot *uint64, currentOffset *int, slotMap map[string]SlotInfo) error {
	// primitives first
	if size, ok := getPrimitiveSize(typ); ok {
		return handlePrimitiveType(name, size, currentSlot, currentOffset, slotMap)
	}
	// dynamic types (map/slices)
	if isDynamicType(typ) {
		return handleDynamicType(name, typ, currentSlot, currentOffset, slotMap)
	}
	// nested structs or arrays
	if typ.Kind() == reflect.Struct || typ.Kind() == reflect.Array {
		// reset offset for new struct/array
		if *currentOffset > 0 {
			*currentSlot++
			*currentOffset = 0
		}
		// add prefix for the recursion
		nextPrefix := name
		if typ.Kind() == reflect.Struct { // struct fields separated by dot
			nextPrefix += "."
		}
		// add base slot info for struct
		if typ.Kind() == reflect.Struct {
			baseSlot := computeSlotHash(*currentSlot) // base is current slot
			baseInfo := SlotInfo{
				Slot:      baseSlot,
				Offset:    0,
				Size:      0, // struct has no size
				IsDynamic: false,
				ValueType: typ,
				SubSlots:  make(map[string]SlotInfo),
			}
			slotMap[name] = baseInfo
		}

		if err := assignSlotsRecursive(typ, nextPrefix, currentSlot, currentOffset, slotMap); err != nil {
			return err
		}
		if typ.Kind() == reflect.Array {
			baseInfo := slotMap[name]
			elemType := typ.Elem()
			elemSize, fixed := getElemSize(elemType)
			if !fixed || elemSize == 0 {
				return fmt.Errorf("element size must be non-zero fixed size for array")
			}
			baseInfo.Resolver = NewFixedResolver(baseInfo.Slot, uint64(elemSize))
			slotMap[name] = baseInfo
		}
		// reset offset post processing
		if *currentOffset > 0 {
			*currentSlot++
			*currentOffset = 0
		}
		return nil
	}
	return fmt.Errorf("unsupported type: %s", typ)
}

func assignSlotsRecursive(typ reflect.Type, prefix string, currentSlot *uint64,
	currentOffset *int, slotMap map[string]SlotInfo) error {

	// struct handling
	if typ.Kind() == reflect.Struct {
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			fName := prefix + field.Name
			if err := processField(fName, field.Type, currentSlot, currentOffset, slotMap); err != nil {
				return err
			}

			//link sub-slot for parent child relation
			// link sub-slot to base struct slot info
			if strings.Contains(fName, ".") {
				baseKey := prefix[:strings.LastIndex(prefix, ".")]
				if baseKey != "" {
					if baseInfo, ok := slotMap[baseKey]; ok && baseInfo.SubSlots != nil {
						baseInfo.SubSlots[field.Name] = slotMap[fName] // relative SlotInfo (rel slot index)
						slotMap[baseKey] = baseInfo
					}
				}
			}
		}
		return nil
	}

	if typ.Kind() == reflect.Array {
		for i := 0; i < typ.Len(); i++ {
			// for arrays, construct the indexing name ==> arr[0]
			elemName := fmt.Sprintf("%s[%d]", prefix, i)
			if err := processField(elemName, typ.Elem(), currentSlot, currentOffset, slotMap); err != nil {
				return err
			}
		}
		return nil
	}
	return fmt.Errorf("assignSlotRecursive: unsupported type: %s", typ)
}

func computeSlotHash(slotIndex uint64) common.Hash {
	var slot [32]byte
	binary.BigEndian.PutUint64(slot[24:], slotIndex)
	return common.BytesToHash(slot[:])
}
