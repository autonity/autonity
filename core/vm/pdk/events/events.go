package events

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm/pdk/abiselector"
	"github.com/autonity/autonity/core/vm/pdk/storage"
	"github.com/autonity/autonity/crypto"
)

var eventCache = sync.Map{}

type eventSchema struct {
	topicHash     common.Hash
	indexedFields []int
	dataArgs      abi.Arguments
}

// Emit emits the given event to the storage's state log, caching the event schema for future use.
func Emit(st *storage.Storage, event interface{}) error {
	var schema *eventSchema
	var err error
	typ := reflect.TypeOf(event)
	if typ.Kind() != reflect.Struct {
		return fmt.Errorf("event must be a struct")
	}
	schemaVal, found := eventCache.Load(typ)
	if !found {
		schema, err = inferEventSchema(typ)
		if err != nil {
			return err
		}
		eventCache.Store(typ, schema)
	} else {
		schema = schemaVal.(*eventSchema)
	}
	// prepare indexed topics
	topics := make([]common.Hash, 0, len(schema.indexedFields)+1)
	val := reflect.ValueOf(event)
	for _, idx := range schema.indexedFields {
		fieldVal := val.Field(idx)
		indexedHash, err := encodeTopic(fieldVal)
		if err != nil {
			return err
		}
		topics = append(topics, indexedHash)
	}

	// data values
	var dataValues []interface{}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !field.IsExported() {
			continue
		}
		if isIndexed(field) {
			continue
		}
		dataValues = append(dataValues, val.Field(i).Interface())
	}
	data, err := schema.dataArgs.PackValues(dataValues)
	if err != nil {
		return err
	}
	st.AddLog(topics, data)
	return nil
}

func encodeTopic(fieldValue reflect.Value) (common.Hash, error) {
	abiType, err := abiselector.ResolveABIType(fieldValue.Type())
	if err != nil {
		return common.Hash{}, fmt.Errorf("unsupported indexed field type %s", fieldValue.Type().String())
	}
	if abiType.T == abi.StringTy {
		return crypto.Keccak256Hash([]byte(fieldValue.String())), nil
	}
	if abiType.T == abi.BytesTy {
		return crypto.Keccak256Hash(fieldValue.Bytes()), nil
	}

	packed, err := abi.Arguments{{Type: abiType}}.Pack(fieldValue.Interface())
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to pack indexed field: %v", err)
	}
	// large and complex types needs to be hashed as per ABI spec
	if abiType.T == abi.ArrayTy || abiType.T == abi.SliceTy || abiType.T == abi.TupleTy {
		return crypto.Keccak256Hash(packed), nil
	}
	return common.BytesToHash(packed), nil
}

func isIndexed(structField reflect.StructField) bool {
	indexedTag := structField.Tag.Get("indexed")
	return indexedTag == "true"
}

func inferEventSchema(typ reflect.Type) (*eventSchema, error) {
	var indexedFields []int
	var argTypes []string
	var dataArgs abi.Arguments

	eventName := typ.Name()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !field.IsExported() {
			continue
		}

		abiType, err := abiselector.ResolveABIType(field.Type)
		if err != nil {
			return nil, fmt.Errorf("unsupported field type %s for field %s", field.Type.String(), field.Name)
		}
		argTypes = append(argTypes, abiType.String())

		if isIndexed(field) {
			indexedFields = append(indexedFields, i)
			continue
		}
		dataArgs = append(dataArgs, abi.Argument{
			Name: field.Name, // abi uses field names as argument names
			Type: abiType,
		})
	}
	eventSignature := fmt.Sprintf("%s(%s)", eventName, strings.Join(argTypes, ","))
	topicHash := crypto.Keccak256Hash([]byte(eventSignature))

	return &eventSchema{
		topicHash:     topicHash,
		indexedFields: indexedFields,
		dataArgs:      dataArgs,
	}, nil
}
