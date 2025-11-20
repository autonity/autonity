package storage

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"reflect"

	"github.com/autonity/autonity/common"
)

func (s *Storage) getDynamicArrayLength(field string) (uint64, error) {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return 0, fmt.Errorf("field %s not found", field)
	}
	if !slotInfo.IsDynamic || slotInfo.ValueType.Kind() != reflect.Slice {
		return 0, fmt.Errorf("field %s is not a dynamic array", field)
	}
	data := s.stateDB.GetState(s.address, slotInfo.Slot)
	return binary.BigEndian.Uint64(data[:8]), nil
}

func (s *Storage) setDynamicArrayLength(field string, length uint64) error {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return fmt.Errorf("field %s not found", field)
	}
	if !slotInfo.IsDynamic || slotInfo.ValueType.Kind() != reflect.Slice {
		return fmt.Errorf("field %s is not a dynamic array", field)
	}
	var data common.Hash
	binary.BigEndian.PutUint64(data[:8], length)
	s.stateDB.SetState(s.address, slotInfo.Slot, data)
	return nil
}

func (s *Storage) PushUint256ToSlice(field string, value Uint256) error {
	len, err := s.getDynamicArrayLength(field)
	if err != nil {
		return err
	}
	if err := s.SetUint256ArrayElement(field, len, value); err != nil {
		return err
	}
	return s.setDynamicArrayLength(field, len+1)
}

func (s *Storage) PushAddressToSlice(field string, value Address) error {
	len, err := s.getDynamicArrayLength(field)
	if err != nil {
		return err
	}
	if err := s.SetAddressArrayElement(field, len, value); err != nil {
		return err
	}
	return s.setDynamicArrayLength(field, len+1)
}

func (s *Storage) PopUint256FromSlice(field string) (Uint256, error) {
	len, err := s.getDynamicArrayLength(field)
	if err != nil {
		return Uint256{}, err
	}
	if len == 0 {
		return Uint256{}, fmt.Errorf("pop from empty array")
	}
	value, err := s.GetUint256ArrayElement(field, len-1)
	if err != nil {
		return Uint256{}, err
	}
	// Clear slot (security + gas refund)
	slot, _ := s.arrayElementSlotRaw(field, len-1)
	s.stateDB.SetState(s.address, slot, common.Hash{})
	err = s.setDynamicArrayLength(field, len-1)
	if err != nil {
		return Uint256{}, err
	}
	return value, nil
}

// todo(review):

// Internal raw slot resolver (used by Pop & auto-extend)
func (s *Storage) arrayElementSlotRaw(field string, index uint64) (common.Hash, error) {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return common.Hash{}, fmt.Errorf("field %s not found", field)
	}
	_, elemInfo, err := resolveSliceElemSlot(slotInfo.Slot, slotInfo, index)
	if err != nil {
		return common.Hash{}, err
	}
	return elemInfo.Slot, nil
}

// GetUint256ArrayElement element from dynamic OR fixed-size array (top-level or nested)
func (s *Storage) GetUint256ArrayElement(field string, index uint64) (Uint256, error) {
	slot, err := s.arrayElementRawSlot(field, index)
	if err != nil {
		return Uint256{}, err
	}
	data := s.stateDB.GetState(s.address, slot)
	return NewUint256FromBig(new(big.Int).SetBytes(data[:])), nil
}

func (s *Storage) SetUint256ArrayElement(field string, index uint64, value Uint256) error {
	slot, err := s.arrayElementRawSlot(field, index)
	if err != nil {
		return err
	}
	var data common.Hash
	value.WriteToSlice(data[:])
	s.stateDB.SetState(s.address, slot, data)
	return nil
}

func (s *Storage) GetAddressArrayElement(field string, index uint64) (Address, error) {
	slot, err := s.arrayElementRawSlot(field, index)
	if err != nil {
		return Address{}, err
	}
	data := s.stateDB.GetState(s.address, slot)
	// address in array element = left-padded
	return NewAddressFromBytes(data[:20]), nil
}

func (s *Storage) SetAddressArrayElement(field string, index uint64, value Address) error {
	slot, err := s.arrayElementRawSlot(field, index)
	if err != nil {
		return err
	}
	var data common.Hash
	copy(data[:20], value.ToCommonAddress().Bytes()) // left-padded
	s.stateDB.SetState(s.address, slot, data)
	return nil
}

func (s *Storage) GetUint64ArrayElement(field string, index uint64) (uint64, error) {
	slot, err := s.arrayElementRawSlot(field, index)
	if err != nil {
		return 0, err
	}
	data := s.stateDB.GetState(s.address, slot)
	return binary.BigEndian.Uint64(data[:8]), nil
}

func (s *Storage) SetUint64ArrayElement(field string, index uint64, value uint64) error {
	slot, err := s.arrayElementRawSlot(field, index)
	if err != nil {
		return err
	}
	var data common.Hash
	binary.BigEndian.PutUint64(data[:8], value)
	s.stateDB.SetState(s.address, slot, data)
	return nil
}

// arrayElementRawSlot returns the exact storage slot for any array element
func (s *Storage) arrayElementRawSlot(field string, index uint64) (common.Hash, error) {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return common.Hash{}, fmt.Errorf("field %s not found", field)
	}

	// Accept both dynamic arrays and fixed-size arrays
	if !slotInfo.IsDynamic && slotInfo.ValueType.Kind() != reflect.Array {
		return common.Hash{}, fmt.Errorf("field %s is not an array", field)
	}

	if slotInfo.KeyType != nil {
		return common.Hash{}, fmt.Errorf("use Get/SetFromMapArray for mapping-value arrays")
	}

	_, elemInfo, err := resolveSliceElemSlot(slotInfo.Slot, slotInfo, index)
	if err != nil {
		return common.Hash{}, err
	}
	return elemInfo.Slot, nil
}
