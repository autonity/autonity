package storage

import (
	"encoding/binary"
	"fmt"
	"reflect"

	"github.com/autonity/autonity/common"
)

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

func assignSlotsRecursive(typ reflect.Type, prefix string, currentSlot *uint64,
	currentOffset *int, slotMap map[string]SlotInfo) error {

	for fIndex := range typ.NumField() {
		// naive implementation assuming basic types
		field := typ.Field(fIndex)
		fType := field.Type
		fName := prefix + field.Name

		if size, ok := getPrimitiveSize(fType); ok {
			if *currentOffset+size > 32 { // doesn't fit in current slot
				// move to next slot
				*currentSlot++
				*currentOffset = 0
			}
			// fits, assign
			info := SlotInfo{
				Slot:      computeSlotHash(*currentSlot),
				Offset:    *currentOffset,
				Size:      size,
				IsDynamic: false,
			}
			slotMap[fName] = info
			*currentOffset += size
			continue
		}

		if isDynamicType(fType) {
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
			} else {
				panic("unsupported dynamic type")
			}
			slotMap[fName] = info
			*currentSlot++
			*currentOffset = 0
			continue
		}

		if fType.Kind() == reflect.Struct {
			err := assignSlotsRecursive(fType, fName+".", currentSlot, currentOffset, slotMap)
			if err != nil {
				panic(err)
			}
			// continue with next field
			continue
		}
		return fmt.Errorf("unsupported type: %s", fType)
	}
	return nil
}

func computeSlotHash(slotIndex uint64) common.Hash {
	var slot [32]byte
	binary.BigEndian.PutUint64(slot[24:], slotIndex)
	return common.BytesToHash(slot[:])
}
