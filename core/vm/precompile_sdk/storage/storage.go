package storage

import (
	"fmt"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm"
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
