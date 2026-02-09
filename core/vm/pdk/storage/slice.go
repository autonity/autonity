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
	elemSize     uint64
}

func (s *Slice[T]) Bind(st *Storage, baseSlot common.Hash, offset uint64) (common.Hash, uint64) {
	s.st = st
	if offset > 0 {
		baseSlot = addSlot(baseSlot, 1)
	}
	s.lenSlot = baseSlot
	s.elemBaseSlot = crypto.Keccak256Hash(baseSlot.Bytes())
	s.elemSize = getSlotConsumption[T]()
	return addSlot(baseSlot, 1), 0
}

func (s *Slice[T]) Len() uint64 {
	if s.st == nil {
		return 0
	}
	data := s.st.GetState(s.lenSlot)
	return binary.BigEndian.Uint64(data[24:])
}

func (s *Slice[T]) Get(i uint64) *T {
	itemSlot := addSlot(s.elemBaseSlot, i*s.elemSize)
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

func (s *Slice[T]) Clear() {
	length := s.Len()
	for i := uint64(0); i < length; i++ {
		// Get(i) creates a bound instance of the element at the correct slot
		elem := s.Get(i)
		// Recursively clear that element (handles structs, vars, nested slices)
		RecursiveClear(elem)
	}
	// Now safe to zero the header/slots
	if err := s.setLength(0); err != nil {
		panic(err)
	}
}

func (s *Slice[T]) setLength(newLen uint64) error {
	oldLen := s.Len()

	var headData common.Hash
	binary.BigEndian.PutUint64(headData[24:], newLen)
	s.st.SetState(s.lenSlot, headData)

	if newLen < oldLen {
		for i := newLen; i < oldLen; i++ {
			base := addSlot(s.elemBaseSlot, i*s.elemSize)
			for j := uint64(0); j < s.elemSize; j++ {
				s.st.SetState(addSlot(base, j), common.Hash{})
			}
		}
	}
	return nil
}
