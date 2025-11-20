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

	//extract relative slot index from subSlot.Slot
	relIndex := binary.BigEndian.Uint64(subSlot.Slot[24:])

	// calculate final slot as baseSlot + relIndex
	baseBig := new(big.Int).SetBytes(baseSlot.Bytes())
	finalBig := new(big.Int).Add(baseBig, new(big.Int).SetUint64(relIndex))
	return common.BigToHash(finalBig), subSlot, nil
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
	numElements := sliceInfo.numElements
	if numElements == 0 {
		numElements = 1 // should not happen, but just in case
	}
	offset := new(big.Int).Mul(new(big.Int).SetUint64(index), new(big.Int).SetUint64(numElements))
	finalBig := new(big.Int).Add(dataStartBig, offset)

	elemInfo := SlotInfo{
		Slot:      common.BigToHash(finalBig),
		Offset:    0, // base offset is 0 relative to element slot
		Size:      0, // we don't know the size here
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

func (s *Storage) GetUint64FromMapStruct(mapField string, key interface{}, structFieldName string) (uint64, error) {
	mapValSlot, mapInfo, err := s.resolveMapSlot(mapField, key)
	if err != nil {
		return 0, err
	}
	fieldSlot, fieldInfo, err := resolveStructFieldSlot(structFieldName, mapInfo, mapValSlot)
	if err != nil {
		return 0, err
	}
	data := s.stateDB.GetState(s.address, fieldSlot)
	return binary.BigEndian.Uint64(data[fieldInfo.Offset : fieldInfo.Offset+fieldInfo.Size]), nil
}

func (s *Storage) SetUint64InMapStruct(mapField string, key interface{}, structFieldName string, value uint64) error {
	mapValSlot, mapInfo, err := s.resolveMapSlot(mapField, key)
	if err != nil {
		return err
	}
	fieldSlot, fieldInfo, err := resolveStructFieldSlot(structFieldName, mapInfo, mapValSlot)
	if err != nil {
		return err
	}
	data := s.stateDB.GetState(s.address, fieldSlot)
	binary.BigEndian.PutUint64(data[fieldInfo.Offset:fieldInfo.Offset+fieldInfo.Size], value)
	s.stateDB.SetState(s.address, fieldSlot, data)
	return nil
}

func (s *Storage) GetUint256FromMapSlice(mapField string, key interface{}, index uint64) (Uint256, error) {
	sliceHeadSlot, mapInfo, err := s.resolveMapSlot(mapField, key)
	if err != nil {
		return Uint256{}, err
	}
	elemSlot, _, err := resolveSliceElemSlot(sliceHeadSlot, mapInfo, index)
	if err != nil {
		return Uint256{}, err
	}

	data := s.stateDB.GetState(s.address, elemSlot)
	bi := big.NewInt(0)
	bi.SetBytes(data[:])
	return NewUint256FromBig(bi), nil
}

func (s *Storage) GetAddressFromMapSlice(mapField string, key interface{}, index uint64) (Address, error) {
	sliceHeadSlot, mapInfo, err := s.resolveMapSlot(mapField, key)
	if err != nil {
		return Address{}, err
	}
	elemSlot, _, err := resolveSliceElemSlot(sliceHeadSlot, mapInfo, index)
	if err != nil {
		return Address{}, err
	}
	data := s.stateDB.GetState(s.address, elemSlot)
	return NewAddressFromBytes(data[12:]), nil
}

func (s *Storage) GetUint64FromMapSliceStruct(mapField string, key interface{}, index uint64, structFieldName string) (uint64, error) {
	sliceHeadSlot, mapInfo, err := s.resolveMapSlot(mapField, key)
	if err != nil {
		return 0, err
	}
	structBaseSlot, _, err := resolveSliceElemSlot(sliceHeadSlot, mapInfo, index)
	if err != nil {
		return 0, err
	}

	fieldSlot, fieldInfo, err := resolveStructFieldSlot(structFieldName, mapInfo, structBaseSlot)
	if err != nil {
		return 0, err
	}
	data := s.stateDB.GetState(s.address, fieldSlot)
	return binary.BigEndian.Uint64(data[fieldInfo.Offset : fieldInfo.Offset+fieldInfo.Size]), nil
}

func (s *Storage) GetAddressFromMapSliceStruct(mapField string, key interface{}, index uint64, structFieldName string) (Address, error) {
	sliceHeadSlot, mapInfo, err := s.resolveMapSlot(mapField, key)
	if err != nil {
		return Address{}, err
	}
	structBaseSlot, _, err := resolveSliceElemSlot(sliceHeadSlot, mapInfo, index)
	if err != nil {
		return Address{}, err
	}

	fieldSlot, subSlot, err := resolveStructFieldSlot(structFieldName, mapInfo, structBaseSlot)
	if err != nil {
		return Address{}, err
	}
	data := s.stateDB.GetState(s.address, fieldSlot)
	return NewAddressFromBytes(data[subSlot.Offset : subSlot.Offset+subSlot.Size]), nil
}

func (s *Storage) GetAddressFromMap(field string, key interface{}) (Address, error) {
	slot, slotInfo, err := s.resolveMapSlot(field, key)
	if err != nil {
		return Address{}, fmt.Errorf("no such slot for %s", field)
	}
	if !slotInfo.IsDynamic || slotInfo.ValueType == nil || slotInfo.KeyType == nil {
		return Address{}, fmt.Errorf("field %s is not an Address type", field)
	}
	if slotInfo.ValueType != reflect.TypeOf(Address{}) {
		return Address{}, fmt.Errorf("field %s is not a Address type", field)
	}
	data := s.stateDB.GetState(s.address, slot)
	// left padded address bytes
	return NewAddressFromBytes(data[:20]), nil
}

func (s *Storage) SetAddressInMap(field string, key interface{}, value Address) error {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return fmt.Errorf("no such slot for %s", field)
	}
	if !slotInfo.IsDynamic || slotInfo.ValueType == nil || slotInfo.KeyType == nil {
		return fmt.Errorf("field %s is not a dynamic bytes type", field)
	}
	if slotInfo.ValueType != reflect.TypeOf(Address{}) {
		return fmt.Errorf("field %s is not a Address type", field)
	}
	keyBytes, err := encodeTo32Bytes(key, slotInfo.KeyType)
	if err != nil {
		return err
	}
	valueSlot := crypto.Keccak256Hash(append(keyBytes, slotInfo.Slot.Bytes()...))
	var data common.Hash
	// left align address bytes (20 bytes in the beginning)
	copy(data[:20], value.ToCommonAddress().Bytes())
	s.stateDB.SetState(s.address, valueSlot, data)
	return nil
}

func (s *Storage) GetUint256FromMap(field string, key interface{}) (Uint256, error) {
	slot, slotInfo, err := s.resolveMapSlot(field, key)
	if err != nil {
		return Uint256{}, fmt.Errorf("no such slot for %s", field)
	}
	if !slotInfo.IsDynamic || slotInfo.ValueType == nil || slotInfo.KeyType == nil {
		return Uint256{}, fmt.Errorf("field %s is not a dynamic bytes type", field)
	}
	if slotInfo.ValueType != reflect.TypeOf(Uint256{}) {
		return Uint256{}, fmt.Errorf("field %s is not a uint256 type", field)
	}
	data := s.stateDB.GetState(s.address, slot)
	bi := big.NewInt(0)
	bi.SetBytes(data[:])
	return NewUint256FromBig(bi), nil
}

func (s *Storage) SetUint256InMap(field string, key interface{}, value Uint256) error {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return fmt.Errorf("no such slot for %s", field)
	}
	if !slotInfo.IsDynamic || slotInfo.ValueType == nil || slotInfo.KeyType == nil {
		return fmt.Errorf("field %s is not a dynamic bytes type", field)
	}
	if slotInfo.ValueType != reflect.TypeOf(Uint256{}) {
		return fmt.Errorf("field %s is not a uint256 type", field)
	}
	keyBytes, err := encodeTo32Bytes(key, slotInfo.KeyType)
	if err != nil {
		return err
	}
	valueSlot := crypto.Keccak256Hash(append(keyBytes, slotInfo.Slot.Bytes()...))
	var data common.Hash
	value.WriteToSlice(data[:])
	s.stateDB.SetState(s.address, valueSlot, data)
	return nil
}

func (s *Storage) SetUint64InMap(field string, key interface{}, value uint64) error {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return fmt.Errorf("no such slot for %s", field)
	}
	if !slotInfo.IsDynamic || slotInfo.KeyType == nil || slotInfo.ValueType == nil {
		return fmt.Errorf("field %s is not a map", field)
	}
	if slotInfo.ValueType.Kind() != reflect.Uint64 {
		return fmt.Errorf("map value type for %s is not uint64", field)
	}
	keyBytes, err := encodeTo32Bytes(key, slotInfo.KeyType)
	if err != nil {
		return err
	}
	valueSlot := crypto.Keccak256Hash(append(keyBytes, slotInfo.Slot.Bytes()...))
	var data common.Hash
	// left align uint64 in the slot
	binary.BigEndian.PutUint64(data[:8], value) // Right-align in slot
	s.stateDB.SetState(s.address, valueSlot, data)
	return nil
}

func (s *Storage) GetUint64FromMap(field string, key interface{}) (uint64, error) {
	slot, slotInfo, err := s.resolveMapSlot(field, key)
	if err != nil {
		return 0, fmt.Errorf("no such slot for %s", field)
	}
	if !slotInfo.IsDynamic || slotInfo.KeyType == nil || slotInfo.ValueType == nil {
		return 0, fmt.Errorf("field %s is not a map", field)
	}
	if slotInfo.ValueType.Kind() != reflect.Uint64 {
		return 0, fmt.Errorf("map value type for %s is not uint64", field)
	}
	data := s.stateDB.GetState(s.address, slot)
	return binary.BigEndian.Uint64(data[:8]), nil
}
