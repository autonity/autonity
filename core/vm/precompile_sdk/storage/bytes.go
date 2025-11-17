package storage

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"reflect"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/crypto"
)

func (s *Storage) SetBytes(field string, value []byte) error {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return fmt.Errorf("no such slot for %s", field)
	}
	// for byte slices only value is expected
	if !slotInfo.IsDynamic || slotInfo.ValueType == nil || slotInfo.KeyType != nil {
		return fmt.Errorf("field %s is not a dynamic bytes type", field)
	}

	if slotInfo.ValueType.Kind() != reflect.Uint8 {
		return fmt.Errorf("field %s is not a bytes slice (value type not uint8)", field)
	}
	// first store the data length
	length := uint64(len(value))
	// right aligned
	var lenData common.Hash
	binary.BigEndian.PutUint64(lenData[24:], length)
	s.stateDB.SetState(s.address, slotInfo.Slot, lenData)
	if length == 0 {
		return nil
	}
	// calculate base slot for data
	base := crypto.Keccak256Hash(slotInfo.Slot.Bytes())
	baseBig := new(big.Int).SetBytes(base.Bytes())
	numChunks := (length + 31) / 32
	for i := range numChunks {
		chunkSlotBig := new(big.Int).Add(baseBig, big.NewInt(int64(i)))
		// calculate slot for this chunk
		chunkSlotHash := common.BigToHash(chunkSlotBig)
		var chunkData common.Hash
		start := i * 32
		end := start + 32
		if end > length {
			end = length
		}
		// write chunk
		copy(chunkData[:end-start], value[start:end])
		s.stateDB.SetState(s.address, chunkSlotHash, chunkData)
	}
	return nil
}

func (s *Storage) GetBytes(field string) ([]byte, error) {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return nil, fmt.Errorf("no such slot for %s", field)
	}
	// for byte slices only value is expected
	if !slotInfo.IsDynamic || slotInfo.ValueType == nil || slotInfo.KeyType != nil {
		return nil, fmt.Errorf("field %s is not a dynamic bytes type", field)
	}

	if slotInfo.ValueType.Kind() != reflect.Uint8 {
		return nil, fmt.Errorf("field %s is not a bytes slice (value type not uint8)", field)
	}

	lenData := s.stateDB.GetState(s.address, slotInfo.Slot)
	length := binary.BigEndian.Uint64(lenData[24:])
	if length == 0 {
		return []byte{}, nil
	}
	result := make([]byte, length)
	base := crypto.Keccak256Hash(slotInfo.Slot.Bytes())
	baseBig := new(big.Int).SetBytes(base.Bytes())
	numChunks := (length + 31) / 32
	for i := range numChunks {
		chunkSlotBig := new(big.Int).Add(baseBig, big.NewInt(int64(i)))
		// calculate slot for this chunk
		chunkSlotHash := common.BigToHash(chunkSlotBig)
		chunkData := s.stateDB.GetState(s.address, chunkSlotHash)
		start := i * 32
		end := start + 32
		if end > length {
			end = length
		}
		// read chunk
		copy(result[start:end], chunkData[:end-start])
	}
	return result, nil
}
