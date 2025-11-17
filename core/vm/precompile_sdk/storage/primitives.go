package storage

import (
	"encoding/binary"
	"fmt"
)

func (s *Storage) GetUint64(field string) (uint64, error) {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return 0, fmt.Errorf("no such slot for %s", field)
	}
	if slotInfo.Size != 8 {
		return 0, fmt.Errorf("invalid slot layout, field %s is not a uint64", field)
	}
	data := s.stateDB.GetState(s.address, slotInfo.Slot)
	return binary.BigEndian.Uint64(data[slotInfo.Offset : slotInfo.Offset+slotInfo.Size]), nil
}

func (s *Storage) SetUint64(field string, value uint64) error {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return fmt.Errorf("no such slot for %s", field)
	}
	if slotInfo.Size != 8 {
		return fmt.Errorf("invalid slot layout, field %s is not a uint64", field)
	}
	data := s.stateDB.GetState(s.address, slotInfo.Slot)
	binary.BigEndian.PutUint64(data[slotInfo.Offset:slotInfo.Offset+slotInfo.Size], value)
	s.stateDB.SetState(s.address, slotInfo.Slot, data)
	return nil
}

func (s *Storage) GetUint32(field string) (uint32, error) {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return 0, fmt.Errorf("no such slot for %s", field)
	}
	if slotInfo.Size != 4 {
		return 0, fmt.Errorf("invalid slot layout, field %s is not a uint32", field)
	}
	data := s.stateDB.GetState(s.address, slotInfo.Slot)
	return binary.BigEndian.Uint32(data[slotInfo.Offset : slotInfo.Offset+slotInfo.Size]), nil
}

func (s *Storage) SetUint32(field string, value uint32) error {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return fmt.Errorf("no such slot for %s", field)
	}
	if slotInfo.Size != 4 {
		return fmt.Errorf("invalid slot layout, field %s is not a uint32", field)
	}
	data := s.stateDB.GetState(s.address, slotInfo.Slot)
	binary.BigEndian.PutUint32(data[slotInfo.Offset:slotInfo.Offset+slotInfo.Size], value)
	s.stateDB.SetState(s.address, slotInfo.Slot, data)
	return nil
}

func (s *Storage) GetUint16(field string) (uint16, error) {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return 0, fmt.Errorf("no such slot for %s", field)
	}
	if slotInfo.Size != 2 {
		return 0, fmt.Errorf("invalid slot layout, field %s is not a uint16", field)
	}
	data := s.stateDB.GetState(s.address, slotInfo.Slot)
	return binary.BigEndian.Uint16(data[slotInfo.Offset : slotInfo.Offset+slotInfo.Size]), nil
}

func (s *Storage) SetUint16(field string, value uint16) error {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return fmt.Errorf("no such slot for %s", field)
	}
	if slotInfo.Size != 2 {
		return fmt.Errorf("invalid slot layout, field %s is not a uint16", field)
	}
	data := s.stateDB.GetState(s.address, slotInfo.Slot)
	binary.BigEndian.PutUint16(data[slotInfo.Offset:slotInfo.Offset+slotInfo.Size], value)
	s.stateDB.SetState(s.address, slotInfo.Slot, data)
	return nil
}

func (s *Storage) GetUint8(field string) (uint8, error) {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return 0, fmt.Errorf("no such slot for %s", field)
	}
	if slotInfo.Size != 1 {
		return 0, fmt.Errorf("invalid slot layout, field %s is not a uint8", field)
	}
	data := s.stateDB.GetState(s.address, slotInfo.Slot)
	return data[slotInfo.Offset], nil
}

func (s *Storage) SetUint8(field string, value uint8) error {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return fmt.Errorf("no such slot for %s", field)
	}
	if slotInfo.Size != 1 {
		return fmt.Errorf("invalid slot layout, field %s is not a uint8", field)
	}
	data := s.stateDB.GetState(s.address, slotInfo.Slot)
	data[slotInfo.Offset] = value
	s.stateDB.SetState(s.address, slotInfo.Slot, data)
	return nil
}

func (s *Storage) GetInt64(field string) (int64, error) {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return 0, fmt.Errorf("no such slot for %s", field)
	}
	if slotInfo.Size != 8 {
		return 0, fmt.Errorf("invalid slot layout, field %s is not a int64", field)
	}
	data := s.stateDB.GetState(s.address, slotInfo.Slot)
	return int64(binary.BigEndian.Uint64(data[slotInfo.Offset : slotInfo.Offset+slotInfo.Size])), nil
}

func (s *Storage) SetInt64(field string, value int64) error {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return fmt.Errorf("no such slot for %s", field)
	}
	if slotInfo.Size != 8 {
		return fmt.Errorf("invalid slot layout, field %s is not a int64", field)
	}
	data := s.stateDB.GetState(s.address, slotInfo.Slot)
	binary.BigEndian.PutUint64(data[slotInfo.Offset:slotInfo.Offset+slotInfo.Size], uint64(value))
	s.stateDB.SetState(s.address, slotInfo.Slot, data)
	return nil
}

func (s *Storage) GetInt32(field string) (int32, error) {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return 0, fmt.Errorf("no such slot for %s", field)
	}
	if slotInfo.Size != 4 {
		return 0, fmt.Errorf("invalid slot layout, field %s is not a int32", field)
	}
	data := s.stateDB.GetState(s.address, slotInfo.Slot)
	return int32(binary.BigEndian.Uint32(data[slotInfo.Offset : slotInfo.Offset+slotInfo.Size])), nil
}

func (s *Storage) SetInt32(field string, value int32) error {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return fmt.Errorf("no such slot for %s", field)
	}
	if slotInfo.Size != 4 {
		return fmt.Errorf("invalid slot layout, field %s is not a int32", field)
	}
	data := s.stateDB.GetState(s.address, slotInfo.Slot)
	binary.BigEndian.PutUint32(data[slotInfo.Offset:slotInfo.Offset+slotInfo.Size], uint32(value))
	s.stateDB.SetState(s.address, slotInfo.Slot, data)
	return nil
}

func (s *Storage) GetBool(field string) (bool, error) {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return false, fmt.Errorf("no such slot for %s", field)
	}
	if slotInfo.Size != 1 {
		return false, fmt.Errorf("invalid slot layout, field %s is not a bool", field)
	}
	data := s.stateDB.GetState(s.address, slotInfo.Slot)
	return data[slotInfo.Offset] == 1, nil
}

func (s *Storage) SetBool(field string, value bool) error {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return fmt.Errorf("no such slot for %s", field)
	}
	if slotInfo.Size != 1 {
		return fmt.Errorf("invalid slot layout, field %s is not a bool", field)
	}
	data := s.stateDB.GetState(s.address, slotInfo.Slot)
	if value && data[slotInfo.Offset] == 1 ||
		!value && data[slotInfo.Offset] == 0 {
		return nil // no change
	}
	if value {
		data[slotInfo.Offset] = 1
	} else {
		data[slotInfo.Offset] = 0
	}
	s.stateDB.SetState(s.address, slotInfo.Slot, data)
	return nil
}
