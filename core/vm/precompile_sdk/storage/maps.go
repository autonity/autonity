package storage

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"reflect"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/crypto"
)

func (s *Storage) GetAddressFromMap(field string, key interface{}) (Address, error) {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return Address{}, fmt.Errorf("no such slot for %s", field)
	}
	if !slotInfo.IsDynamic || slotInfo.ValueType == nil || slotInfo.KeyType == nil {
		return Address{}, fmt.Errorf("field %s is not a dynamic bytes type", field)
	}
	if slotInfo.ValueType != reflect.TypeOf(Address{}) {
		return Address{}, fmt.Errorf("field %s is not a Address type", field)
	}
	keyBytes, err := encodeTo32Bytes(key, slotInfo.KeyType)
	if err != nil {
		return Address{}, err
	}
	valueSlot := crypto.Keccak256Hash(append(keyBytes, slotInfo.Slot.Bytes()...))
	data := s.stateDB.GetState(s.address, valueSlot)
	// right aligned address bytes
	return NewAddressFromBytes(data[12:]), nil
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
	copy(data[12:], value.ToCommonAddress().Bytes())
	s.stateDB.SetState(s.address, valueSlot, data)
	return nil
}

func (s *Storage) GetUint256FromMap(field string, key interface{}) (Uint256, error) {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return Uint256{}, fmt.Errorf("no such slot for %s", field)
	}
	if !slotInfo.IsDynamic || slotInfo.ValueType == nil || slotInfo.KeyType == nil {
		return Uint256{}, fmt.Errorf("field %s is not a dynamic bytes type", field)
	}
	if slotInfo.ValueType != reflect.TypeOf(Uint256{}) {
		return Uint256{}, fmt.Errorf("field %s is not a uint256 type", field)
	}
	keyBytes, err := encodeTo32Bytes(key, slotInfo.KeyType)
	if err != nil {
		return Uint256{}, err
	}
	valueSlot := crypto.Keccak256Hash(append(keyBytes, slotInfo.Slot.Bytes()...))
	data := s.stateDB.GetState(s.address, valueSlot)
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
	binary.BigEndian.PutUint64(data[24:], value) // Right-align in slot
	s.stateDB.SetState(s.address, valueSlot, data)
	return nil
}

func (s *Storage) GetUint64FromMap(field string, key interface{}) (uint64, error) {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return 0, fmt.Errorf("no such slot for %s", field)
	}
	if !slotInfo.IsDynamic || slotInfo.KeyType == nil || slotInfo.ValueType == nil {
		return 0, fmt.Errorf("field %s is not a map", field)
	}
	if slotInfo.ValueType.Kind() != reflect.Uint64 {
		return 0, fmt.Errorf("map value type for %s is not uint64", field)
	}
	keyBytes, err := encodeTo32Bytes(key, slotInfo.KeyType)
	if err != nil {
		return 0, err
	}
	valueSlot := crypto.Keccak256Hash(append(keyBytes, slotInfo.Slot.Bytes()...))
	data := s.stateDB.GetState(s.address, valueSlot)
	return binary.BigEndian.Uint64(data[24:]), nil
}
