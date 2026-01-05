package storage

import (
	"encoding/binary"
	"fmt"

	"github.com/autonity/autonity/common"
)

type Slice[T any] struct {
	path *Path
}

func NewSlice[T any](p *Path) *Slice[T] {
	return &Slice[T]{path: p}
}

func (s *Slice[T]) BindPath(p *Path) {
	s.path = p
}

func (s *Slice[T]) Len() (uint64, error) {
	if s.path == nil {
		return 0, nil
	}
	data := s.path.st.GetState(s.path.slot)
	return binary.BigEndian.Uint64(data[:8]), nil
}

func (s *Slice[T]) Get(i uint64) (T, error) {
	var zero T
	length, err := s.Len()
	if err != nil || i >= length {
		return zero, fmt.Errorf("index out of range")
	}
	var val T
	err = Load(s.path.Index(i), &val)
	return val, err
}

func (s *Slice[T]) Set(i uint64, value T) error {
	length, err := s.Len()
	if err != nil || i >= length {
		return fmt.Errorf("index out of range")
	}
	return Save(s.path.Index(i), value)
}

func (s *Slice[T]) Append(value T) error {
	length, err := s.Len()
	if err != nil {
		return err
	}
	// set the length
	if err := s.setLength(length + 1); err != nil {
		return err
	}
	// Set the new value at the end
	if err := Save(s.path.Index(length), value); err != nil {
		return err
	}
	return nil
}

func (s *Slice[T]) Grow() (T, error) {
	var zero T
	length, err := s.Len()
	if err != nil {
		return zero, err
	}
	// set the length
	if err = s.setLength(length + 1); err != nil {
		return zero, err
	}
	return s.Get(length)
}

func (s *Slice[T]) Pop() error {
	length, err := s.Len()
	if err != nil {
		return err
	}
	if length == 0 {
		return fmt.Errorf("empty slice")
	}
	// Clear slot
	tailP := s.path.Index(length - 1)
	s.path.st.SetState(tailP.slot, common.Hash{})

	return s.setLength(length - 1)
}

func (s *Slice[T]) setLength(newLen uint64) error {
	var headData common.Hash
	binary.BigEndian.PutUint64(headData[:8], newLen)
	s.path.st.SetState(s.path.slot, headData)
	return nil
}
