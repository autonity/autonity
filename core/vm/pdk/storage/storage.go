package storage

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/core/vm"
)

type Storage struct {
	address common.Address // storage scope, should be contract address
	stateDB vm.StateDB     // state accessor

	// cache
	cache map[common.Hash]common.Hash
	dirty map[common.Hash]struct{}
}

func NewStorage(address common.Address, stateDB vm.StateDB) *Storage {
	return &Storage{
		address: address,
		stateDB: stateDB,
		cache:   make(map[common.Hash]common.Hash),
		dirty:   make(map[common.Hash]struct{}),
	}
}

func (s *Storage) GetState(slot common.Hash) common.Hash {
	if val, ok := s.cache[slot]; ok {
		return val
	}
	val := s.stateDB.GetState(s.address, slot)
	s.cache[slot] = val
	return val
}

func (s *Storage) SetState(slot common.Hash, value common.Hash) {
	s.cache[slot] = value
	s.dirty[slot] = struct{}{}
	return
}

func (s *Storage) Commit() {
	for slot := range s.dirty {
		val := s.cache[slot]
		s.stateDB.SetState(s.address, slot, val)
	}
	s.dirty = make(map[common.Hash]struct{})
	return
}

func (s *Storage) AddLog(topics []common.Hash, data []byte) {
	log := &types.Log{
		Address: s.address,
		Topics:  topics,
		Data:    data,
	}

	s.stateDB.AddLog(log)
	return
}
