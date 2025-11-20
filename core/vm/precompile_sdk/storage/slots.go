package storage

import (
	"encoding/binary"
	"fmt"
	"reflect"

	"github.com/autonity/autonity/common"
)

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
}

func AssignSlots(stateType reflect.Type) map[string]SlotInfo {
	if stateType.Kind() != reflect.Struct {
		panic("stateType must be a struct type")
	}

	slotMap := make(map[string]SlotInfo)
	currentSlotIndex := uint64(0)
	currentOffset := int(0)
	// start recursive assignement
	err := assignSlotsRecursive(stateType, "", &currentSlotIndex, &currentOffset, slotMap)
	if err != nil {
		panic(err)
	}

	return slotMap
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
	info := SlotInfo{
		Slot:      computeSlotHash(*currentSlot),
		Offset:    *currentOffset,
		Size:      32, // take full slot
		IsDynamic: true,
	}
	if fType.Kind() == reflect.Map {
		info.KeyType = fType.Key()
		info.ValueType = fType.Elem()
	} else if fType.Kind() == reflect.Slice {
		info.ValueType = fType.Elem()
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

		if err := assignSlotsRecursive(typ, nextPrefix, currentSlot, currentOffset, slotMap); err != nil {
			return err
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
