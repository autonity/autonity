package storage

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"reflect"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/crypto"
)

func resolveStructFieldSlot(fieldName string, structInfo SlotInfo, baseSlot common.Hash) (common.Hash, SlotInfo, error) {
	if structInfo.ValueType.Kind() != reflect.Struct {
		return common.Hash{}, SlotInfo{}, fmt.Errorf("struct field %s is not a struct", fieldName)
	}

	subSlot, ok := structInfo.SubSlots[fieldName]
	if !ok {
		return common.Hash{}, SlotInfo{}, fmt.Errorf("no such sub-slot %s in struct", fieldName)
	}

	//extract relative slot index from subSlot.Slot //todo
	relIndex := binary.BigEndian.Uint64(subSlot.Slot[24:])

	// calculate final slot as baseSlot + relIndex
	baseBig := new(big.Int).SetBytes(baseSlot.Bytes())
	finalBig := new(big.Int).Add(baseBig, new(big.Int).SetUint64(relIndex))
	return common.BigToHash(finalBig), subSlot, nil
}

// getElemSize calculates the storage size of a single element
// Used for Fixed Array Resolver calculations
func getElemSize(typ reflect.Type) (int, bool) {
	// Primitive
	if size, ok := getPrimitiveSize(typ); ok {
		return size, true
	}
	// Struct
	if typ.Kind() == reflect.Struct {
		// Compile layout to see total size
		_, slots := compileTypeLayout(typ)
		return int(slots * 32), true
	}
	// Array
	if typ.Kind() == reflect.Array {
		elemSize, ok := getElemSize(typ.Elem())
		if !ok {
			return 0, false
		}
		return elemSize * typ.Len(), true
	}
	if typ.Kind() == reflect.Map || typ.Kind() == reflect.Slice {
		return 32, true
	}
	return 0, false
}

func resolveSliceElemSlot(baseSlot common.Hash, sliceInfo SlotInfo, index uint64) (common.Hash, SlotInfo, error) {
	if sliceInfo.ValueType.Kind() != reflect.Slice && sliceInfo.ValueType.Kind() != reflect.Array {
		return common.Hash{}, SlotInfo{}, fmt.Errorf("slice element index %d is not a slice or array", index)
	}

	// figure out starting location
	var dataStartBig *big.Int
	if sliceInfo.ValueType.Kind() == reflect.Slice {
		dataStart := crypto.Keccak256Hash(baseSlot.Bytes())
		dataStartBig = new(big.Int).SetBytes(dataStart.Bytes())
	} else {
		// array or struct array
		dataStartBig = new(big.Int).SetBytes(baseSlot.Bytes())
	}

	// calculate element slot
	elemSize, isFixed := getElemSize(sliceInfo.ValueType)
	if !isFixed && elemSize == 0 {
		return common.Hash{}, SlotInfo{}, fmt.Errorf("cannot determine element size for slice/array")
	}
	offset := new(big.Int).Mul(new(big.Int).SetUint64(index), new(big.Int).SetUint64(uint64(elemSize/32)))
	if elemSize%32 != 0 {
		offset = offset.Add(offset, big.NewInt(1)) // round up for partial slots
	}
	finalBig := new(big.Int).Add(dataStartBig, offset)

	elemInfo := SlotInfo{
		Slot:      common.BigToHash(finalBig),
		Offset:    0, // base offset is 0 relative to element slot
		Size:      elemSize,
		SubSlots:  sliceInfo.SubSlots,
		ValueType: sliceInfo.ValueType,
	}
	return common.BigToHash(finalBig), elemInfo, nil
}

func (s *Storage) resolveMapSlot(field string, key interface{}) (common.Hash, SlotInfo, error) {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return common.Hash{}, SlotInfo{}, fmt.Errorf("no such slot for %s", field)
	}

	if !slotInfo.IsDynamic || slotInfo.KeyType == nil || slotInfo.ValueType == nil {
		return common.Hash{}, SlotInfo{}, fmt.Errorf("field %s is not a map", field)
	}

	keyBytes, err := encodeTo32Bytes(key, slotInfo.KeyType)
	if err != nil {
		return common.Hash{}, SlotInfo{}, err
	}
	valueSlot := crypto.Keccak256Hash(append(keyBytes, slotInfo.Slot.Bytes()...))
	return valueSlot, slotInfo, nil
}
