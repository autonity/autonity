package storage

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/core/vm"
)

// Storage handles the interaction with the underlying state database and provides a cache.
type Storage struct {
	address common.Address // storage scope, should be contract address
	stateDB vm.StateDB     // state accessor

	// cache
	cache map[common.Hash]common.Hash
	dirty map[common.Hash]struct{}
}

// NewStorage creates a new storage instance.
func NewStorage(address common.Address, stateDB vm.StateDB) *Storage {
	return &Storage{
		address: address,
		stateDB: stateDB,
		cache:   make(map[common.Hash]common.Hash),
		dirty:   make(map[common.Hash]struct{}),
	}
}

// GetState retrieves the state value for a given slot.
func (s *Storage) GetState(slot common.Hash) common.Hash {
	if val, ok := s.cache[slot]; ok {
		return val
	}
	val := s.stateDB.GetState(s.address, slot)
	s.cache[slot] = val
	return val
}

// SetState sets the state value for a given slot.
func (s *Storage) SetState(slot common.Hash, value common.Hash) {
	s.cache[slot] = value
	s.dirty[slot] = struct{}{}
}

// Commit flushes all dirty state changes to the state database.
func (s *Storage) Commit() {
	for slot := range s.dirty {
		val := s.cache[slot]
		s.stateDB.SetState(s.address, slot, val)
	}
	s.dirty = make(map[common.Hash]struct{})
}

// AddLog adds a log to the state database.
func (s *Storage) AddLog(topics []common.Hash, data []byte) {
	log := &types.Log{
		Address: s.address,
		Topics:  topics,
		Data:    data,
	}

	s.stateDB.AddLog(log)
}
