package storage

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"reflect"

	"github.com/autonity/autonity/common"
)

func addSlot(baseSlot common.Hash, numSlots uint64) common.Hash {
	baseBig := new(big.Int).SetBytes(baseSlot.Bytes())
	numSlotBig := new(big.Int).SetUint64(numSlots)
	sumBig := new(big.Int).Add(baseBig, numSlotBig)
	return common.BigToHash(sumBig)
}

func slotDiff(baseSlot, nextSlot common.Hash) uint64 {
	baseBig := new(big.Int).SetBytes(baseSlot.Bytes())
	nextBig := new(big.Int).SetBytes(nextSlot.Bytes())
	return new(big.Int).Sub(nextBig, baseBig).Uint64()
}

func getSlotConsumption[T any]() uint64 {
	var zero T
	// simulating bindstate to get the consumed slots
	consumed := BindState(nil, common.Hash{}, &zero)
	if consumed == 0 {
		return 1
	}
	return consumed
}

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

type Clearer interface {
	Clear()
}

func RecursiveClear(target any) {
	if c, ok := target.(Clearer); ok {
		c.Clear()
		return
	}

	val := reflect.ValueOf(target)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() == reflect.Struct {
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			if field.CanAddr() {
				RecursiveClear(field.Addr().Interface())
			}
		}
	}
}
