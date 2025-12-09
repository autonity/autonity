package storage

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"reflect"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/crypto"
)

var (
	zeroHashBytes    = common.Hash{}.Bytes()
	zeroAddressBytes = common.Address{}.Bytes()
)

type ValueAccessor interface {
	//ReadAt reads the value at the given slot and offset
	ReadAt(slot common.Hash, offset int, st *Storage) (any, error)
	// WriteAt returns the full 32 bytes encoded value
	WriteAt(slot common.Hash, offset int, value any, st *Storage) error
}

var registry = make(map[reflect.Type]ValueAccessor)

func getAccessor(t reflect.Type) (ValueAccessor, bool) {
	accessor, ok := registry[t]
	return accessor, ok
}

func init() {
	registry[reflect.TypeOf(uint64(0))] = uintAccessor{size: 8}
	registry[reflect.TypeOf(uint32(0))] = uintAccessor{size: 4}
	registry[reflect.TypeOf(uint16(0))] = uintAccessor{size: 2}
	registry[reflect.TypeOf(uint8(0))] = uintAccessor{size: 1}
	registry[reflect.TypeOf(int64(0))] = intAccessor{size: 8}
	registry[reflect.TypeOf(int32(0))] = intAccessor{size: 4}
	registry[reflect.TypeOf(false)] = boolAccessor{}
	registry[reflect.TypeOf(Uint256{})] = Uint256Accessor{}
	registry[reflect.TypeOf(common.Address{})] = AddressAccessor{}
	registry[reflect.TypeOf([]byte{})] = ByteAccessor{}
	registry[reflect.TypeOf(common.Hash{})] = HashAccessor{}
}

type HashAccessor struct{}

func (u HashAccessor) ReadAt(slot common.Hash, offset int, st *Storage) (any, error) {
	if offset != 0 {
		return nil, fmt.Errorf("hash values must start from zero offset")
	}
	return st.stateDB.GetState(st.address, slot), nil
}

func (u HashAccessor) WriteAt(slot common.Hash, offset int, value any, st *Storage) error {
	if offset != 0 {
		return fmt.Errorf("hash values must start from zero offset")
	}
	st.stateDB.SetState(st.address, slot, value.(common.Hash))
	return nil
}

type Uint256Accessor struct{}

func (u Uint256Accessor) ReadAt(slot common.Hash, offset int, st *Storage) (any, error) {
	if offset != 0 {
		return nil, fmt.Errorf("uint256 values must start from zero offset")
	}
	data := st.stateDB.GetState(st.address, slot)
	var u256 Uint256
	u256.SetBytes(data[:])
	return u256, nil
}

func (u Uint256Accessor) WriteAt(slot common.Hash, _ int, value any, st *Storage) error {
	var u256 = value.(Uint256)
	var data common.Hash
	u256.WriteToSlice(data[:])
	st.stateDB.SetState(st.address, slot, data)
	return nil
}

type AddressAccessor struct{}

func (a AddressAccessor) ReadAt(slot common.Hash, offset int, st *Storage) (any, error) {
	data := st.stateDB.GetState(st.address, slot)
	return common.BytesToAddress(data[offset : offset+20]), nil
}

func (a AddressAccessor) WriteAt(slot common.Hash, offset int, value any, st *Storage) error {
	data := st.stateDB.GetState(st.address, slot)
	var addr = value.(common.Address)
	copy(data[offset:offset+20], zeroAddressBytes)
	copy(data[offset:offset+20], addr.Bytes())

	st.stateDB.SetState(st.address, slot, data)
	return nil
}

type uintAccessor struct{ size int }

func (u uintAccessor) ReadAt(slot common.Hash, offset int, st *Storage) (any, error) {
	if offset+u.size > 32 {
		return nil, fmt.Errorf("offset and size exceed slot boundary")
	}
	data := st.stateDB.GetState(st.address, slot)
	fieldBytes := data[offset : offset+u.size]
	switch u.size {
	case 1:
		return fieldBytes[0], nil
	case 2:
		return binary.BigEndian.Uint16(fieldBytes), nil
	case 4:
		return binary.BigEndian.Uint32(fieldBytes), nil
	case 8:
		return binary.BigEndian.Uint64(fieldBytes), nil
	default:
		return nil, fmt.Errorf("invalid size: %d", u.size)
	}
}

func (u uintAccessor) WriteAt(slot common.Hash, offset int, value any, st *Storage) error {
	if offset+u.size > 32 {
		panic(fmt.Errorf("offset and size exceed slot boundary"))
	}
	data := st.stateDB.GetState(st.address, slot)
	copy(data[offset:offset+u.size], zeroHashBytes[:u.size])
	switch val := value.(type) {
	case uint8:
		data[offset] = val
	case uint16:
		binary.BigEndian.PutUint16(data[offset:offset+u.size], val)
	case uint32:
		binary.BigEndian.PutUint32(data[offset:offset+u.size], val)
	case uint64:
		binary.BigEndian.PutUint64(data[offset:offset+u.size], val)
	default:
		panic(fmt.Errorf("invalid type: %T", value))
	}
	st.stateDB.SetState(st.address, slot, data)
	return nil
}

type intAccessor struct{ size int }

func (i intAccessor) ReadAt(slot common.Hash, offset int, st *Storage) (any, error) {
	uVal, err := uintAccessor{size: i.size}.ReadAt(slot, offset, st)
	if err != nil {
		return nil, err
	}
	return int64(uVal.(uint64)), nil
}

func (i intAccessor) WriteAt(slot common.Hash, offset int, value any, st *Storage) error {
	uval := uint64(value.(int64))

	return uintAccessor{size: i.size}.WriteAt(slot, offset, uval, st)
}

type boolAccessor struct{}

func (b boolAccessor) ReadAt(slot common.Hash, offset int, st *Storage) (any, error) {
	data := st.stateDB.GetState(st.address, slot)
	return data[offset] == 1, nil
}

func (b boolAccessor) WriteAt(slot common.Hash, offset int, value any, st *Storage) error {
	data := st.stateDB.GetState(st.address, slot)
	data[offset] = 0
	if value.(bool) {
		data[offset] = 1
	}
	st.stateDB.SetState(st.address, slot, data)
	return nil
}

type ByteAccessor struct{}

func (b ByteAccessor) ReadAt(headSlot common.Hash, _ int, st *Storage) (any, error) {
	if st == nil {
		return nil, fmt.Errorf("storage is nil")
	}
	lenData := st.stateDB.GetState(st.address, headSlot)
	length := binary.BigEndian.Uint64(lenData[:8])
	if length == 0 {
		return []byte{}, nil
	}
	result := make([]byte, length)
	base := crypto.Keccak256Hash(headSlot.Bytes())
	baseBig := new(big.Int).SetBytes(base.Bytes())
	numChunks := (length + 31) / 32
	for i := range numChunks {
		chunkSlotBig := new(big.Int).Add(baseBig, big.NewInt(int64(i)))
		// calculate slot for this chunk
		chunkSlotHash := common.BigToHash(chunkSlotBig)
		chunkData := st.stateDB.GetState(st.address, chunkSlotHash)
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

func (b ByteAccessor) WriteAt(headSlot common.Hash, _ int, value any, st *Storage) error {
	if st == nil {
		panic("BytesAccessor requires Storage for multi-slot write")
	}
	bytesVal := value.([]byte)
	length := uint64(len(bytesVal))
	var head common.Hash
	if length > 1<<20 {
		panic(fmt.Errorf("bytes too large: %d", length))
	}
	binary.BigEndian.PutUint64(head[:8], length)
	st.stateDB.SetState(st.address, headSlot, head)
	base := crypto.Keccak256Hash(headSlot.Bytes())
	baseBig := new(big.Int).SetBytes(base.Bytes())
	numChunks := (length + 31) / 32
	for i := 0; i < int(numChunks); i++ {
		chunkSlotBig := new(big.Int).Add(baseBig, big.NewInt(int64(i)))
		var chunk common.Hash
		start := i * 32
		end := start + 32
		if end > int(length) {
			end = int(length)
		}
		copy(chunk[:end-start], bytesVal[start:end])
		chunkSlot := common.BigToHash(chunkSlotBig)
		st.stateDB.SetState(st.address, chunkSlot, chunk)
	}
	return nil
}
