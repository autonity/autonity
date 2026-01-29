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
	maxInt256        = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 255), big.NewInt(1))
	minInt256        = new(big.Int).Neg(new(big.Int).Lsh(big.NewInt(1), 255))
	tt255            = new(big.Int).Lsh(big.NewInt(1), 255) // two to the power of 255 to check the -ve number
	tt256            = new(big.Int).Lsh(big.NewInt(1), 256) // two to the power of 256 for modulo operation
)

type Var[T any] struct {
	st       *Storage
	baseSlot common.Hash
	offset   uint64
}

func (v *Var[T]) Bind(st *Storage, baseSlot common.Hash, offset uint64) (common.Hash, uint64) {
	var zero T
	acc, ok := getAccessor(reflect.TypeOf(zero))
	if !ok {
		panic(fmt.Errorf("no accessor for type %T", zero))
	}
	size := acc.Size()

	if uint64(size)+offset > 32 {
		baseSlot = addSlot(baseSlot, 1)
		offset = 0
	}
	v.st = st
	v.offset = offset
	v.baseSlot = baseSlot
	if offset+uint64(size) == 32 { // exactly filled the slot
		// return new slot
		return addSlot(baseSlot, 1), 0
	}
	return baseSlot, offset + uint64(size)
}

func (v *Var[T]) Get() T {
	var zero T
	acc, ok := getAccessor(reflect.TypeOf(zero))
	if !ok {
		panic(fmt.Errorf("no accessor found for type %T", zero))
	}
	val, err := acc.ReadAt(v.baseSlot, int(v.offset), v.st)
	if err != nil {
		panic(err)
	}
	return val.(T)
}

func (v *Var[T]) Set(val T) {
	typ := reflect.TypeOf(val)
	acc, ok := getAccessor(typ)
	if !ok {
		panic(fmt.Errorf("no accessor found for type %T", val))
	}
	err := acc.WriteAt(v.baseSlot, int(v.offset), val, v.st)
	if err != nil {
		panic(err)
	}
	return
}

func (v *Var[T]) Clear() {
	var zero T
	typ := reflect.TypeOf(zero)
	acc, ok := getAccessor(typ)
	if !ok {
		panic(fmt.Errorf("no accessor for type %T", zero))
	}

	if clearer, ok := acc.(ClearerAccessor); ok {
		if err := clearer.Clear(v.baseSlot, int(v.offset), v.st); err != nil {
			panic(err)
		}
		return
	}

	if err := acc.WriteAt(v.baseSlot, int(v.offset), zero, v.st); err != nil {
		panic(err)
	}
}

// ClearerAccessor is an optional interface for accessors that need special cleanup logic
type ClearerAccessor interface {
	Clear(slot common.Hash, offset int, st *Storage) error
}

type ValueAccessor interface {
	//ReadAt reads the value at the given slot and offset
	ReadAt(slot common.Hash, offset int, st *Storage) (any, error)
	//WriteAt writes the value at the given slot and offset
	WriteAt(slot common.Hash, offset int, value any, st *Storage) error

	Size() int
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
	registry[reflect.TypeOf(&big.Int{})] = BigIntAccessor{}

	// Register fixed byte arrays [1]byte ... [32]byte
	byteType := reflect.TypeOf(uint8(0))
	for i := 1; i <= 32; i++ {
		arrType := reflect.ArrayOf(i, byteType)
		registry[arrType] = FixedByteAccessor{
			typ:  arrType,
			size: i,
		}
	}
}

type FixedByteAccessor struct {
	typ  reflect.Type
	size int
}

func (f FixedByteAccessor) Size() int { return f.size }

func (f FixedByteAccessor) ReadAt(slot common.Hash, offset int, st *Storage) (any, error) {
	if offset+f.size > 32 {
		return nil, fmt.Errorf("offset and size exceed slot boundary")
	}
	data := st.GetState(slot)
	slice := data[offset : offset+f.size]

	val := reflect.New(f.typ).Elem()
	for i := 0; i < f.size; i++ {
		val.Index(i).Set(reflect.ValueOf(slice[i]))
	}
	return val.Interface(), nil
}

func (f FixedByteAccessor) WriteAt(slot common.Hash, offset int, value any, st *Storage) error {
	if offset+f.size > 32 {
		panic(fmt.Errorf("offset and size exceed slot boundary"))
	}
	// Value is expected to be [N]byte. We need to convert it to a slice to copy it.
	valVal := reflect.ValueOf(value)
	if valVal.Kind() != reflect.Array || valVal.Type().Elem().Kind() != reflect.Uint8 {
		return fmt.Errorf("expected [N]byte, got %T", value)
	}

	data := st.GetState(slot)
	copy(data[offset:offset+f.size], zeroHashBytes[:f.size]) // Clear previous

	for i := 0; i < f.size; i++ {
		data[offset+i] = uint8(valVal.Index(i).Uint())
	}

	st.SetState(slot, data)
	return nil
}

// BigIntAccessor handles reading and writing big.Int values in the range of int256. for uint256 use Uint256Accessor
type BigIntAccessor struct{}

func (b BigIntAccessor) Size() int { return 32 }

func (b BigIntAccessor) ReadAt(slot common.Hash, offset int, st *Storage) (any, error) {
	if offset != 0 {
		return nil, fmt.Errorf("big.Int values must start from zero offset")
	}
	data := st.GetState(slot)
	ret := new(big.Int).SetBytes(data[:])
	// check if negative

	if ret.Cmp(tt255) >= 0 { // negative number
		ret.Sub(ret, tt256) // convert to negative value
	}
	return ret, nil
}

func (b BigIntAccessor) WriteAt(slot common.Hash, _ int, value any, st *Storage) error {
	valBig, ok := value.(*big.Int)
	if !ok {
		return fmt.Errorf("expected *big.Int type, got %T", value)
	}
	if valBig.Cmp(maxInt256) > 0 {
		return fmt.Errorf("big.Int value exceeds maximum int256")
	}
	if valBig.Cmp(minInt256) < 0 {
		return fmt.Errorf("big.Int value below minimum int256")
	}
	toWrite := new(big.Int).Set(valBig)
	if toWrite.Sign() < 0 { // negative number
		toWrite.Add(toWrite, tt256) // uint256 representation
	}
	var padded [32]byte
	toWrite.FillBytes(padded[:])
	st.SetState(slot, padded)
	return nil
}

type HashAccessor struct{}

func (b HashAccessor) Size() int { return 32 }
func (u HashAccessor) ReadAt(slot common.Hash, offset int, st *Storage) (any, error) {
	if offset != 0 {
		return nil, fmt.Errorf("hash values must start from zero offset")
	}
	return st.GetState(slot), nil
}

func (u HashAccessor) WriteAt(slot common.Hash, offset int, value any, st *Storage) error {
	if offset != 0 {
		return fmt.Errorf("hash values must start from zero offset")
	}
	st.SetState(slot, value.(common.Hash))
	return nil
}

type Uint256Accessor struct{}

func (u Uint256Accessor) Size() int { return 32 }
func (u Uint256Accessor) ReadAt(slot common.Hash, offset int, st *Storage) (any, error) {
	if offset != 0 {
		return nil, fmt.Errorf("uint256 values must start from zero offset")
	}
	data := st.GetState(slot)
	var u256 Uint256
	u256.SetBytes(data[:])
	return u256, nil
}

func (u Uint256Accessor) WriteAt(slot common.Hash, _ int, value any, st *Storage) error {
	var u256 = value.(Uint256)
	var data common.Hash
	u256.WriteToSlice(data[:])
	st.SetState(slot, data)
	return nil
}

type AddressAccessor struct{}

func (a AddressAccessor) Size() int { return 20 }
func (a AddressAccessor) ReadAt(slot common.Hash, offset int, st *Storage) (any, error) {
	data := st.GetState(slot)
	return common.BytesToAddress(data[offset : offset+20]), nil
}

func (a AddressAccessor) WriteAt(slot common.Hash, offset int, value any, st *Storage) error {
	data := st.GetState(slot)
	var addr = value.(common.Address)
	copy(data[offset:offset+20], zeroAddressBytes)
	copy(data[offset:offset+20], addr.Bytes())

	st.SetState(slot, data)
	return nil
}

type uintAccessor struct{ size int }

func (u uintAccessor) Size() int { return u.size }
func (u uintAccessor) ReadAt(slot common.Hash, offset int, st *Storage) (any, error) {
	if offset+u.size > 32 {
		return nil, fmt.Errorf("offset and size exceed slot boundary")
	}
	data := st.GetState(slot)
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
	data := st.GetState(slot)
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
	st.SetState(slot, data)
	return nil
}

type intAccessor struct{ size int }

func (i intAccessor) Size() int { return i.size }
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

func (b boolAccessor) Size() int { return 1 }

func (b boolAccessor) ReadAt(slot common.Hash, offset int, st *Storage) (any, error) {
	data := st.GetState(slot)
	return data[offset] == 1, nil
}

func (b boolAccessor) WriteAt(slot common.Hash, offset int, value any, st *Storage) error {
	data := st.GetState(slot)
	data[offset] = 0
	if value.(bool) {
		data[offset] = 1
	}
	st.SetState(slot, data)
	return nil
}

type ByteAccessor struct{}

func (b ByteAccessor) Size() int { return 32 }

func (b ByteAccessor) ReadAt(headSlot common.Hash, _ int, st *Storage) (any, error) {
	if st == nil {
		return nil, fmt.Errorf("storage is nil")
	}
	lenData := st.GetState(headSlot)
	length := binary.BigEndian.Uint64(lenData[24:])
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
		chunkData := st.GetState(chunkSlotHash)
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
	binary.BigEndian.PutUint64(head[24:], length)
	st.SetState(headSlot, head)
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
		st.SetState(chunkSlot, chunk)
	}
	return nil
}

func (b ByteAccessor) Clear(headSlot common.Hash, _ int, st *Storage) error {
	if st == nil {
		return nil
	}
	lenData := st.GetState(headSlot)
	length := binary.BigEndian.Uint64(lenData[24:])
	if length == 0 {
		return nil
	}

	// Zero out data chunks
	base := crypto.Keccak256Hash(headSlot.Bytes())
	baseBig := new(big.Int).SetBytes(base.Bytes())
	numChunks := (length + 31) / 32
	zeroHash := common.Hash{}

	for i := 0; i < int(numChunks); i++ {
		chunkSlotBig := new(big.Int).Add(baseBig, big.NewInt(int64(i)))
		chunkSlot := common.BigToHash(chunkSlotBig)
		st.SetState(chunkSlot, zeroHash)
	}

	// Zero out length
	st.SetState(headSlot, zeroHash)
	return nil
}
