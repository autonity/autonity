package storage

import (
	"encoding/binary"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/crypto"
)

type Slice[T any] struct {
	st           *Storage
	lenSlot      common.Hash
	elemBaseSlot common.Hash
}

func (s *Slice[T]) Bind(st *Storage, baseSlot common.Hash, offset uint64) (common.Hash, uint64) {
	s.st = st
	if offset > 0 {
		baseSlot = addSlot(baseSlot, 1)
	}
	s.lenSlot = baseSlot
	s.elemBaseSlot = crypto.Keccak256Hash(baseSlot.Bytes())
	return addSlot(baseSlot, 1), 0
}

func (s *Slice[T]) Len() uint64 {
	if s.st == nil {
		return 0
	}
	data := s.st.GetState(s.lenSlot)
	return binary.BigEndian.Uint64(data[:8])
}

func (s *Slice[T]) Get(i uint64) *T {
	itemSlot := addSlot(s.elemBaseSlot, i)
	val := new(T)
	BindState(s.st, itemSlot, val)
	return val
}

func (s *Slice[T]) Append(setFunc func(*T)) {
	length := s.Len()
	// set the length
	if err := s.setLength(length + 1); err != nil {
		panic(err)
	}

	elem := s.Get(length) // bound element
	setFunc(elem)
}

func (s *Slice[T]) setLength(newLen uint64) error {
	var headData common.Hash
	binary.BigEndian.PutUint64(headData[:8], newLen)
	s.st.SetState(s.lenSlot, headData)
	return nil
}
