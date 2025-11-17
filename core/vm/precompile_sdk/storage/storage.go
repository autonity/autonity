package storage

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm"
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
}

type Storage struct {
	address common.Address // storage scope, generally should be the contract address
	stateDB vm.StateDB     // state accessor
	slotMap map[string]SlotInfo
}

func NewStorage(address common.Address, stateDB vm.StateDB, slots map[string]SlotInfo) *Storage {
	return &Storage{address: address, stateDB: stateDB, slotMap: slots}
}

func (s *Storage) GetUint256(field string) (Uint256, error) {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return Uint256{}, fmt.Errorf("no such slot for %s", field)
	}

	if slotInfo.Size != 32 || slotInfo.Offset != 0 {
		return Uint256{}, fmt.Errorf("invalid slot layout, field %s is not a uint256", field)
	}

	data := s.stateDB.GetState(s.address, slotInfo.Slot)
	bi := big.NewInt(0)
	bi.SetBytes(data[:])
	return NewUint256FromBig(bi), nil
}

func (s *Storage) SetUint256(field string, value Uint256) error {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return fmt.Errorf("no such slot for %s", field)
	}
	if slotInfo.Size != 32 || slotInfo.Offset != 0 {
		return fmt.Errorf("invalid slot layout, field %s is not a uint256", field)
	}
	var data common.Hash
	value.WriteToSlice(data[:])
	s.stateDB.SetState(s.address, slotInfo.Slot, data)
	return nil
}

func (s *Storage) GetAddress(field string) (Address, error) {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return Address{}, fmt.Errorf("no such slot for %s", field)
	}
	if slotInfo.Size != 20 {
		return Address{}, fmt.Errorf("invalid slot layout, field %s is not an address", field)
	}
	data := s.stateDB.GetState(s.address, slotInfo.Slot)
	return NewAddressFromBytes(data[slotInfo.Offset : slotInfo.Offset+slotInfo.Size]), nil
}

func (s *Storage) SetAddress(field string, address Address) error {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return fmt.Errorf("no such slot for %s", field)
	}

	if slotInfo.Size != 20 {
		return fmt.Errorf("invalid slot layout, field %s is not an address", field)
	}
	// read full
	data := s.stateDB.GetState(s.address, slotInfo.Slot)
	// write after the offset in 32 byte slot
	addr := address.ToCommonAddress()
	copy(data[slotInfo.Offset:], addr.Bytes())
	// write full
	s.stateDB.SetState(s.address, slotInfo.Slot, data)
	return nil
}
