package storage

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"reflect"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/core/vm/precompile_sdk/types"
	"github.com/autonity/autonity/crypto"
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

func AssignSlots(stateType reflect.Type) map[string]SlotInfo {
	if stateType.Kind() != reflect.Struct {
		panic("stateType must be a struct type")
	}

	slotMap := make(map[string]SlotInfo)

	currentSlotIndex := uint64(0)
	currentOffset := int(0)
	// start recursive assignement
	err := assignSlotsRecursive(stateType, "", &currentSlotIndex, &currentOffset, slotMap)
	if err != nil {
		panic(err)
	}

	return slotMap
}

func getPrimitiveSize(typ reflect.Type) (int, bool) {
	if reflect.TypeOf(types.Address{}) == typ {
		return 20, true
	}
	if reflect.TypeOf(types.Uint256{}) == typ {
		return 32, true
	}
	switch typ.Kind() {
	case reflect.Bool:
		return 1, true
	case reflect.Int8, reflect.Uint8:
	case reflect.Int16, reflect.Uint16:
		return 2, true
	case reflect.Int32, reflect.Uint32:
		return 4, true
	case reflect.Int64, reflect.Uint64:
		return 8, true
	case reflect.Array:
		if typ.Elem().Kind() == reflect.Uint8 {
			if typ.Len() > 0 && typ.Len() <= 32 {
				return typ.Len(), true
			}
		}
	}
	return 0, false
}

func isDynamicType(typ reflect.Type) bool {
	if typ.Kind() == reflect.Map || typ.Kind() == reflect.Slice {
		return true
	}
	return false
}

func computeSlotHash(slotIndex uint64) common.Hash {
	var slot [32]byte
	binary.BigEndian.PutUint64(slot[24:], slotIndex)
	return common.BytesToHash(slot[:])
}

func assignSlotsRecursive(typ reflect.Type, prefix string, currentSlot *uint64,
	currentOffset *int, slotMap map[string]SlotInfo) error {

	for fIndex := range typ.NumField() {
		// naive implementation assuming basic types
		field := typ.Field(fIndex)
		fType := field.Type
		fName := prefix + field.Name

		if size, ok := getPrimitiveSize(fType); ok {
			if *currentOffset+size > 32 { // doesn't fit in current slot
				// move to next slot
				*currentSlot++
				*currentOffset = 0
			}
			// fits, assign
			info := SlotInfo{
				Slot:      computeSlotHash(*currentSlot),
				Offset:    *currentOffset,
				Size:      size,
				IsDynamic: false,
			}
			slotMap[fName] = info
			*currentOffset += size
		}

		if isDynamicType(fType) {
			// dynamic type will start from the new slot, increment slot and reset offset
			if *currentOffset > 0 {
				*currentSlot++
				*currentOffset = 0
			}
			info := SlotInfo{
				Slot:      computeSlotHash(*currentSlot),
				Offset:    *currentOffset,
				Size:      32, // take full slot
				IsDynamic: true,
			}
			if fType.Kind() == reflect.Map {
				info.KeyType = fType.Key()
				info.ValueType = fType.Elem()
			} else if fType.Kind() == reflect.Slice {
				info.ValueType = fType.Elem()
			} else {
				panic("unsupported dynamic type")
			}
			slotMap[fName] = info
			*currentSlot++
			*currentOffset = 0
			//todo
			continue
		}

		if fType.Kind() == reflect.Struct {
			err := assignSlotsRecursive(fType, fName+".", currentSlot, currentOffset, slotMap)
			if err != nil {
				panic(err)
			}
			// continue with next field
			continue
		}
	}
	return nil
}

func (s *Storage) GetUint256(field string) (types.Uint256, error) {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return types.Uint256{}, fmt.Errorf("no such slot for %s", field)
	}

	if slotInfo.Size != 32 || slotInfo.Offset != 0 {
		return types.Uint256{}, fmt.Errorf("invalid slot layout, field %s is not a uint256", field)
	}

	data := s.stateDB.GetState(s.address, slotInfo.Slot)
	bi := big.NewInt(0)
	bi.SetBytes(data[:])
	return types.NewUint256FromBig(bi), nil
}

func (s *Storage) SetUint256(field string, value types.Uint256) error {
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

func (s *Storage) GetUint64(field string) (uint64, error) {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return 0, fmt.Errorf("no such slot for %s", field)
	}
	if slotInfo.Size != 8 {
		return 0, fmt.Errorf("invalid slot layout, field %s is not a uint64", field)
	}
	data := s.stateDB.GetState(s.address, slotInfo.Slot)
	return binary.BigEndian.Uint64(data[slotInfo.Offset : slotInfo.Offset+slotInfo.Size]), nil
}

func (s *Storage) SetUint64(field string, value uint64) error {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return fmt.Errorf("no such slot for %s", field)
	}
	if slotInfo.Size != 8 {
		return fmt.Errorf("invalid slot layout, field %s is not a uint64", field)
	}
	data := s.stateDB.GetState(s.address, slotInfo.Slot)
	binary.BigEndian.PutUint64(data[slotInfo.Offset:slotInfo.Offset+slotInfo.Size], value)
	s.stateDB.SetState(s.address, slotInfo.Slot, data)
	return nil
}

func (s *Storage) GetUint32(field string) (uint32, error) {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return 0, fmt.Errorf("no such slot for %s", field)
	}
	if slotInfo.Size != 4 {
		return 0, fmt.Errorf("invalid slot layout, field %s is not a uint32", field)
	}
	data := s.stateDB.GetState(s.address, slotInfo.Slot)
	return binary.BigEndian.Uint32(data[slotInfo.Offset : slotInfo.Offset+slotInfo.Size]), nil
}

func (s *Storage) SetUint32(field string, value uint32) error {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return fmt.Errorf("no such slot for %s", field)
	}
	if slotInfo.Size != 4 {
		return fmt.Errorf("invalid slot layout, field %s is not a uint32", field)
	}
	data := s.stateDB.GetState(s.address, slotInfo.Slot)
	binary.BigEndian.PutUint32(data[slotInfo.Offset:slotInfo.Offset+slotInfo.Size], value)
	s.stateDB.SetState(s.address, slotInfo.Slot, data)
	return nil
}

func (s *Storage) GetUint16(field string) (uint16, error) {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return 0, fmt.Errorf("no such slot for %s", field)
	}
	if slotInfo.Size != 2 {
		return 0, fmt.Errorf("invalid slot layout, field %s is not a uint16", field)
	}
	data := s.stateDB.GetState(s.address, slotInfo.Slot)
	return binary.BigEndian.Uint16(data[slotInfo.Offset : slotInfo.Offset+slotInfo.Size]), nil
}

func (s *Storage) SetUint16(field string, value uint16) error {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return fmt.Errorf("no such slot for %s", field)
	}
	if slotInfo.Size != 2 {
		return fmt.Errorf("invalid slot layout, field %s is not a uint16", field)
	}
	data := s.stateDB.GetState(s.address, slotInfo.Slot)
	binary.BigEndian.PutUint16(data[slotInfo.Offset:slotInfo.Offset+slotInfo.Size], value)
	s.stateDB.SetState(s.address, slotInfo.Slot, data)
	return nil
}

func (s *Storage) GetUint8(field string) (uint8, error) {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return 0, fmt.Errorf("no such slot for %s", field)
	}
	if slotInfo.Size != 1 {
		return 0, fmt.Errorf("invalid slot layout, field %s is not a uint8", field)
	}
	data := s.stateDB.GetState(s.address, slotInfo.Slot)
	return data[slotInfo.Offset], nil
}

func (s *Storage) SetUint8(field string, value uint8) error {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return fmt.Errorf("no such slot for %s", field)
	}
	if slotInfo.Size != 1 {
		return fmt.Errorf("invalid slot layout, field %s is not a uint8", field)
	}
	data := s.stateDB.GetState(s.address, slotInfo.Slot)
	data[slotInfo.Offset] = value
	s.stateDB.SetState(s.address, slotInfo.Slot, data)
	return nil
}

func (s *Storage) GetInt64(field string) (int64, error) {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return 0, fmt.Errorf("no such slot for %s", field)
	}
	if slotInfo.Size != 8 {
		return 0, fmt.Errorf("invalid slot layout, field %s is not a int64", field)
	}
	data := s.stateDB.GetState(s.address, slotInfo.Slot)
	return int64(binary.BigEndian.Uint64(data[slotInfo.Offset : slotInfo.Offset+slotInfo.Size])), nil
}

func (s *Storage) SetInt64(field string, value int64) error {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return fmt.Errorf("no such slot for %s", field)
	}
	if slotInfo.Size != 8 {
		return fmt.Errorf("invalid slot layout, field %s is not a int64", field)
	}
	data := s.stateDB.GetState(s.address, slotInfo.Slot)
	binary.BigEndian.PutUint64(data[slotInfo.Offset:slotInfo.Offset+slotInfo.Size], uint64(value))
	s.stateDB.SetState(s.address, slotInfo.Slot, data)
	return nil
}

func (s *Storage) GetInt32(field string) (int32, error) {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return 0, fmt.Errorf("no such slot for %s", field)
	}
	if slotInfo.Size != 4 {
		return 0, fmt.Errorf("invalid slot layout, field %s is not a int32", field)
	}
	data := s.stateDB.GetState(s.address, slotInfo.Slot)
	return int32(binary.BigEndian.Uint32(data[slotInfo.Offset : slotInfo.Offset+slotInfo.Size])), nil
}

func (s *Storage) SetInt32(field string, value int32) error {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return fmt.Errorf("no such slot for %s", field)
	}
	if slotInfo.Size != 4 {
		return fmt.Errorf("invalid slot layout, field %s is not a int32", field)
	}
	data := s.stateDB.GetState(s.address, slotInfo.Slot)
	binary.BigEndian.PutUint32(data[slotInfo.Offset:slotInfo.Offset+slotInfo.Size], uint32(value))
	s.stateDB.SetState(s.address, slotInfo.Slot, data)
	return nil
}

func (s *Storage) GetAddress(field string) (types.Address, error) {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return types.Address{}, fmt.Errorf("no such slot for %s", field)
	}
	if slotInfo.Size != 20 {
		return types.Address{}, fmt.Errorf("invalid slot layout, field %s is not a uint256", field)
	}
	data := s.stateDB.GetState(s.address, slotInfo.Slot)
	return types.NewAddressFromBytes(data[slotInfo.Offset : slotInfo.Offset+slotInfo.Size]), nil
}

func (s *Storage) SetAddress(field string, address types.Address) error {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return fmt.Errorf("no such slot for %s", field)
	}

	if slotInfo.Size != 20 {
		return fmt.Errorf("invalid slot layout, field %s is not a uint256", field)
	}
	addr := address.ToCommonAddress()
	// read full
	data := s.stateDB.GetState(addr, slotInfo.Slot)
	// write after the offset in 32 byte slot
	copy(data[slotInfo.Offset:], addr.Bytes())
	// write full
	s.stateDB.SetState(s.address, slotInfo.Slot, data)
	return nil
}

func (s *Storage) GetBool(field string) (bool, error) {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return false, fmt.Errorf("no such slot for %s", field)
	}
	if slotInfo.Size != 1 {
		return false, fmt.Errorf("invalid slot layout, field %s is not a bool", field)
	}
	data := s.stateDB.GetState(s.address, slotInfo.Slot)
	return data[slotInfo.Offset] == 1, nil
}

func (s *Storage) SetBool(field string, value bool) error {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return fmt.Errorf("no such slot for %s", field)
	}
	if slotInfo.Size != 1 {
		return fmt.Errorf("invalid slot layout, field %s is not a bool", field)
	}
	data := s.stateDB.GetState(s.address, slotInfo.Slot)
	if value && data[slotInfo.Offset] == 1 ||
		!value && data[slotInfo.Offset] == 0 {
		return nil // no change
	}
	if value {
		data[slotInfo.Offset] = 1
	} else {
		data[slotInfo.Offset] = 0
	}
	s.stateDB.SetState(s.address, slotInfo.Slot, data)
	return nil
}

func (s *Storage) SetBytes(field string, value []byte) error {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return fmt.Errorf("no such slot for %s", field)
	}
	// for byte slices only value is expected
	if !slotInfo.IsDynamic || slotInfo.ValueType == nil || slotInfo.KeyType != nil {
		return fmt.Errorf("field %s is not a dynamic bytes type", field)
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

// todo: encode to 32 Bytes, key being passed as interface and keytype passed as reflect type, this keyType is actual
// map keyType and the key is the user value what use wants to pass
func encodeTo32Bytes(key interface{}, keyType reflect.Type) ([]byte, error) {
	// zero-initialized 32-byte buffer
	encodedBytes := make([]byte, 32)
	keyValue := reflect.ValueOf(key)

	if keyValue.Type() != keyType {
		return nil, fmt.Errorf("key type mismatch, expected %v got %v", keyType, keyValue.Type())
	}

	if keyType == reflect.TypeOf(types.Address{}) {
		addr := key.(types.Address)
		keyBytes := addr.ToCommonAddress().Bytes() // 20 bytes
		// Right-align the 20-byte address in the 32-byte buffer
		copy(encodedBytes[12:], keyBytes)
		return encodedBytes, nil
	}

	if keyType == reflect.TypeOf(types.Uint256{}) {
		value := key.(types.Uint256)
		value.WriteToSlice(encodedBytes)
		return encodedBytes, nil
	}

	if keyType.Kind() == reflect.Bool {
		if keyValue.Bool() {
			encodedBytes[31] = 1
		}
		// else, buffer is already all zeros
		return encodedBytes, nil
	}

	var isNegative bool
	size := int(keyType.Size())
	keyBytes := make([]byte, size)
	switch keyType.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val := keyValue.Int()
		isNegative = val < 0
		switch size {
		case 1:
			keyBytes[0] = byte(int8(val))
		case 2:
			binary.BigEndian.PutUint16(keyBytes, uint16(int16(val)))
		case 4:
			binary.BigEndian.PutUint32(keyBytes, uint32(int32(val)))
		case 8:
			binary.BigEndian.PutUint64(keyBytes, uint64(val))
		}
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uval := keyValue.Uint()
		switch size {
		case 1:
			keyBytes[0] = byte(uint8(uval))
		case 2:
			binary.BigEndian.PutUint16(keyBytes, uint16(uval))
		case 4:
			binary.BigEndian.PutUint32(keyBytes, uint32(uval))
		case 8:
			binary.BigEndian.PutUint64(keyBytes, uval)
		}
	default:
		return nil, fmt.Errorf("unsupported key type: %v", keyType)
	}

	if isNegative { // sign extension
		for i := range encodedBytes {
			encodedBytes[i] = 0xff
		}
	}
	// Right-align the key bytes.
	copy(encodedBytes[32-len(keyBytes):], keyBytes)
	return encodedBytes, nil
}

func (s *Storage) GetAddressFromMap(field string, key interface{}) (types.Address, error) {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return types.Address{}, fmt.Errorf("no such slot for %s", field)
	}
	if !slotInfo.IsDynamic || slotInfo.ValueType == nil || slotInfo.KeyType != nil {
		return types.Address{}, fmt.Errorf("field %s is not a dynamic bytes type", field)
	}
	if slotInfo.ValueType != reflect.TypeOf(types.Address{}) {
		return types.Address{}, fmt.Errorf("field %s is not a Address type", field)
	}
	keyBytes, err := encodeTo32Bytes(key, reflect.TypeOf(types.Address{}))
	if err != nil {
		return types.Address{}, err
	}
	valueSlot := crypto.Keccak256Hash(append(keyBytes, slotInfo.Slot.Bytes()...))
	data := s.stateDB.GetState(s.address, valueSlot)
	// right aligned address bytes
	return types.NewAddressFromBytes(data[12:]), nil
}

func (s *Storage) SetAddressInMap(field string, key interface{}, value types.Address) error {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return fmt.Errorf("no such slot for %s", field)
	}
	if !slotInfo.IsDynamic || slotInfo.ValueType == nil || slotInfo.KeyType != nil {
		return fmt.Errorf("field %s is not a dynamic bytes type", field)
	}
	if slotInfo.ValueType != reflect.TypeOf(types.Address{}) {
		return fmt.Errorf("field %s is not a Address type", field)
	}
	keyBytes, err := encodeTo32Bytes(key, reflect.TypeOf(types.Address{}))
	if err != nil {
		return err
	}
	valueSlot := crypto.Keccak256Hash(append(keyBytes, slotInfo.Slot.Bytes()...))
	var data common.Hash
	copy(data[12:], value.ToCommonAddress().Bytes())
	s.stateDB.SetState(s.address, valueSlot, data)
	return nil
}

func (s *Storage) GetUint256FromMap(field string, key interface{}) (types.Uint256, error) {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return types.Uint256{}, fmt.Errorf("no such slot for %s", field)
	}
	if !slotInfo.IsDynamic || slotInfo.ValueType == nil || slotInfo.KeyType != nil {
		return types.Uint256{}, fmt.Errorf("field %s is not a dynamic bytes type", field)
	}
	if slotInfo.ValueType != reflect.TypeOf(types.Uint256{}) {
		return types.Uint256{}, fmt.Errorf("field %s is not a uint256 type", field)
	}
	keyBytes, err := encodeTo32Bytes(key, reflect.TypeOf(types.Uint256{}))
	if err != nil {
		return types.Uint256{}, err
	}
	valueSlot := crypto.Keccak256Hash(append(keyBytes, slotInfo.Slot.Bytes()...))
	data := s.stateDB.GetState(s.address, valueSlot)
	bi := big.NewInt(0)
	bi.SetBytes(data[:])
	return types.NewUint256FromBig(bi), nil
}

func (s *Storage) SetUint256InMap(field string, key interface{}, value types.Uint256) error {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return fmt.Errorf("no such slot for %s", field)
	}
	if !slotInfo.IsDynamic || slotInfo.ValueType == nil || slotInfo.KeyType != nil {
		return fmt.Errorf("field %s is not a dynamic bytes type", field)
	}
	if slotInfo.ValueType != reflect.TypeOf(types.Uint256{}) {
		return fmt.Errorf("field %s is not a uint256 type", field)
	}
	keyBytes, err := encodeTo32Bytes(key, reflect.TypeOf(types.Uint256{}))
	if err != nil {
		return err
	}
	valueSlot := crypto.Keccak256Hash(append(keyBytes, slotInfo.Slot.Bytes()...))
	var data common.Hash
	value.WriteToSlice(data[:])
	s.stateDB.SetState(s.address, valueSlot, data)
	return nil
}

func (s *Storage) SetUint64InMap(field string, key interface{}, value uint64) error {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return fmt.Errorf("no such slot for %s", field)
	}
	if !slotInfo.IsDynamic || slotInfo.KeyType == nil || slotInfo.ValueType == nil {
		return fmt.Errorf("field %s is not a map", field)
	}
	if slotInfo.ValueType.Kind() != reflect.Uint64 {
		return fmt.Errorf("map value type for %s is not uint64", field)
	}
	keyBytes, err := encodeTo32Bytes(key, slotInfo.KeyType)
	if err != nil {
		return err
	}
	valueSlot := crypto.Keccak256Hash(append(keyBytes, slotInfo.Slot.Bytes()...))
	var data common.Hash
	binary.BigEndian.PutUint64(data[24:], value) // Right-align in slot
	s.stateDB.SetState(s.address, valueSlot, data)
	return nil
}

func (s *Storage) GetUint64FromMap(field string, key interface{}) (uint64, error) {
	slotInfo, ok := s.slotMap[field]
	if !ok {
		return 0, fmt.Errorf("no such slot for %s", field)
	}
	if !slotInfo.IsDynamic || slotInfo.KeyType == nil || slotInfo.ValueType == nil {
		return 0, fmt.Errorf("field %s is not a map", field)
	}
	if slotInfo.ValueType.Kind() != reflect.Uint64 {
		return 0, fmt.Errorf("map value type for %s is not uint64", field)
	}
	keyBytes, err := encodeTo32Bytes(key, slotInfo.KeyType)
	if err != nil {
		return 0, err
	}
	valueSlot := crypto.Keccak256Hash(append(keyBytes, slotInfo.Slot.Bytes()...))
	data := s.stateDB.GetState(s.address, valueSlot)
	return binary.BigEndian.Uint64(data[24:]), nil
}
