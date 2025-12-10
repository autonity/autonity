package storage

import (
	"encoding/binary"
	"fmt"
	"reflect"

	"github.com/autonity/autonity/common"
)

// encodeTo32Bytes, key is passed as interface and keytype passed as reflect type, this keyType is actual
// map keyType and the key is the user value what use wants to pass
func encodeTo32Bytes(key interface{}, keyType reflect.Type) ([]byte, error) {
	// zero-initialized 32-byte buffer
	encodedBytes := make([]byte, 32)
	keyValue := reflect.ValueOf(key)

	if keyValue.Type() != keyType {
		return nil, fmt.Errorf("key type mismatch, expected %v got %v", keyType, keyValue.Type())
	}

	if keyType == reflect.TypeOf(common.Address{}) {
		addr := key.(common.Address)
		keyBytes := addr.Bytes() // 20 bytes
		// Right-align the 20-byte address in the 32-byte buffer
		copy(encodedBytes[12:], keyBytes)
		return encodedBytes, nil
	}
	if keyType == reflect.TypeOf(common.Hash{}) {
		hash := key.(common.Hash)
		keyBytes := hash.Bytes() // 32 bytes full slot
		copy(encodedBytes[:], keyBytes)
		return encodedBytes, nil
	}

	if keyType == reflect.TypeOf(Uint256{}) {
		value := key.(Uint256)
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
			keyBytes[0] = uint8(uval)
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

func getPrimitiveSize(typ reflect.Type) (int, bool) {
	if reflect.TypeOf(common.Address{}) == typ {
		return 20, true
	}
	if reflect.TypeOf(common.Hash{}) == typ {
		return 32, true
	}
	if reflect.TypeOf(Uint256{}) == typ {
		return 32, true
	}
	switch typ.Kind() {
	case reflect.Bool, reflect.Int8, reflect.Uint8:
		return 1, true
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

// getElemSize calculates the storage size of a single element
func getElemSize(typ reflect.Type) (int, bool) {
	if size, ok := getPrimitiveSize(typ); ok {
		return size, true
	}
	if typ.Kind() == reflect.Struct {
		_, slots := compileTypeLayout(typ)
		return int(slots * 32), true
	}
	if typ.Kind() == reflect.Array {
		elemSize, ok := getElemSize(typ.Elem())
		if !ok {
			return 0, false
		}
		return elemSize * typ.Len(), true
	}
	if typ.Kind() == reflect.Map || typ.Kind() == reflect.Slice {
		return 32, true
	}
	return 0, false
}

func computeSlotHash(slotIndex uint64) common.Hash {
	var slot [32]byte
	binary.BigEndian.PutUint64(slot[24:], slotIndex)
	return common.BytesToHash(slot[:])
}
