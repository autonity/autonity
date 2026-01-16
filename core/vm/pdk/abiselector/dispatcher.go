package abiselector

import (
	"fmt"
	"math/big"
	"reflect"
	"strings"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/core/vm/pdk/storage"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/log"
)

type selector [4]byte

var (
	GoTypeToABI = map[reflect.Type]func() (abi.Type, error){
		reflect.TypeOf(common.Address{}): func() (abi.Type, error) {
			return abi.NewType("address", "address", nil)
		},
		reflect.TypeOf((*big.Int)(nil)): func() (abi.Type, error) {
			return abi.NewType("uint256", "uint256", nil)
		},
		reflect.TypeOf(false): func() (abi.Type, error) {
			return abi.NewType("bool", "bool", nil)
		},
		reflect.TypeOf([]byte{}): func() (abi.Type, error) {
			return abi.NewType("bytes", "bytes", nil)
		},
		reflect.TypeOf(""): func() (abi.Type, error) {
			return abi.NewType("string", "string", nil)
		},
		reflect.TypeOf(uint8(0)): func() (abi.Type, error) {
			return abi.NewType("uint8", "", nil)
		},
		reflect.TypeOf(uint16(0)): func() (abi.Type, error) {
			return abi.NewType("uint16", "", nil)
		},
		reflect.TypeOf(uint32(0)): func() (abi.Type, error) {
			return abi.NewType("uint32", "", nil)
		},
		reflect.TypeOf(uint64(0)): func() (abi.Type, error) {
			return abi.NewType("uint64", "", nil)
		},
		reflect.TypeOf(uint(0)): func() (abi.Type, error) {
			return abi.NewType("uint256", "", nil)
		},
		reflect.TypeOf(int8(0)): func() (abi.Type, error) {
			return abi.NewType("int8", "", nil)
		},
		reflect.TypeOf(int16(0)): func() (abi.Type, error) {
			return abi.NewType("int16", "", nil)
		},
		reflect.TypeOf(int32(0)): func() (abi.Type, error) {
			return abi.NewType("int32", "", nil)
		},
		reflect.TypeOf(int64(0)): func() (abi.Type, error) {
			return abi.NewType("int64", "", nil)
		},
		reflect.TypeOf(int(0)): func() (abi.Type, error) {
			return abi.NewType("int256", "", nil)
		},
		reflect.TypeOf(common.Hash{}): func() (abi.Type, error) {
			return abi.NewType("bytes32", "", nil)
		},
		reflect.TypeOf([32]byte{}): func() (abi.Type, error) {
			return abi.NewType("bytes32", "", nil)
		},
		reflect.TypeOf([20]byte{}): func() (abi.Type, error) {
			return abi.NewType("bytes20", "", nil)
		},
		reflect.TypeOf([4]byte{}): func() (abi.Type, error) {
			return abi.NewType("bytes4", "", nil)
		},
	}
)

type Dispatcher struct {
	ABI            abi.ABI
	Methods        map[selector]reflect.Value // selector registry
	SelectorToName map[selector]string
}

func newDispatcher() *Dispatcher {
	return &Dispatcher{
		ABI:            abi.ABI{Methods: make(map[string]abi.Method)},
		Methods:        map[selector]reflect.Value{},
		SelectorToName: map[selector]string{},
	}
}

func (d *Dispatcher) AddMethod(method abi.Method, goMethod reflect.Value) {
	d.ABI.Methods[method.Name] = method
	var sel = crypto.Keccak256Hash([]byte(method.Sig)).Bytes()[:4]
	var selArray [4]byte
	copy(selArray[:], sel)
	d.Methods[selArray] = goMethod
	d.SelectorToName[selArray] = method.Name
}

func ResolveABIType(goType reflect.Type) (abi.Type, error) {
	mapped, ok := GoTypeToABI[goType]
	if ok {
		return mapped()
	}

	if goType.Kind() == reflect.Ptr {
		return ResolveABIType(goType.Elem())
	}

	// Handle dynamic arrays (slices only)
	if goType.Kind() == reflect.Slice {
		elemType := goType.Elem()
		elemABIType, err := ResolveABIType(elemType)
		if err != nil {
			return abi.Type{}, fmt.Errorf("failed to resolve slice element type %s: %w", elemType, err)
		}
		arrayTypeStr := elemABIType.String() + "[]"
		return abi.NewType(arrayTypeStr, arrayTypeStr, nil)
	}

	if goType.Kind() == reflect.Struct {
		return resolveStructABIType(goType)
	}

	return abi.Type{}, fmt.Errorf("unsupported type %s", goType.String())
}

// resolveStructABIType converts Go struct types to ABI tuple types
func resolveStructABIType(goType reflect.Type) (abi.Type, error) {
	if goType.Kind() != reflect.Struct {
		return abi.Type{}, fmt.Errorf("expected struct type, got %s", goType.Kind())
	}

	var components []abi.ArgumentMarshaling
	for i := 0; i < goType.NumField(); i++ {
		field := goType.Field(i)
		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		fieldABIType, err := ResolveABIType(field.Type)
		if err != nil {
			return abi.Type{}, fmt.Errorf("failed to resolve field %s of type %s: %w", field.Name, field.Type, err)
		}

		fieldName := strings.ToLower(field.Name[0:1]) + field.Name[1:]
		components = append(components, abi.ArgumentMarshaling{
			Name:         fieldName,
			Type:         fieldABIType.String(),
			InternalType: fieldABIType.String(),
			Components:   convertTupleElemsToArgumentMarshaling(fieldABIType.TupleElems),
		})
	}

	if len(components) == 0 {
		return abi.Type{}, fmt.Errorf("struct %s has no ABI-serializable fields", goType.Name())
	}

	tupleType, err := abi.NewType("tuple", "", components)
	if err != nil {
		return abi.Type{}, fmt.Errorf("failed to create tuple type for struct %s: %w", goType.Name(), err)
	}

	return tupleType, nil
}

// convertTupleElemsToArgumentMarshaling converts []*abi.Type to []abi.ArgumentMarshaling for nested tuples
func convertTupleElemsToArgumentMarshaling(tupleElems []*abi.Type) []abi.ArgumentMarshaling {
	if tupleElems == nil {
		return nil
	}

	var components []abi.ArgumentMarshaling
	for i, elem := range tupleElems {
		components = append(components, abi.ArgumentMarshaling{
			Name:         fmt.Sprintf("field%d", i),
			Type:         elem.String(),
			InternalType: elem.String(),
			Components:   convertTupleElemsToArgumentMarshaling(elem.TupleElems), // Recursive for deeply nested
		})
	}
	return components
}

func InferABIMethods(d *Dispatcher, contractVal reflect.Value) error {
	contractType := contractVal.Type()
loop:
	for i := range contractType.NumMethod() {
		m := contractType.Method(i)
		if !m.IsExported() {
			// no registration for exported methods
			continue
		}
		mt := m.Type
		// method signature ==> func (evm *vm.EVM, caller common.Address, storage *storage.Storage, args...)
		if mt.NumIn() < 3 || mt.NumOut() < 1 ||
			mt.In(0) != contractType || // first arg is receiver itself (contract)
			mt.In(1) != reflect.TypeOf((*vm.EVM)(nil)) || // second arg is *vm.EVM
			mt.In(2) != reflect.TypeOf(common.Address{}) || // third arg is caller address
			mt.In(3) != reflect.TypeOf((*storage.Storage)(nil)) || // 4th arg is storage
			mt.Out(mt.NumOut()-1) != reflect.TypeOf((*error)(nil)).Elem() {
			log.Info("Method %s has incompatible signature, skipping", m.Name)
			continue
		}

		// structure inputs
		inputs := abi.Arguments{}
		for j := 4; j < mt.NumIn(); j++ {
			abiType, err := ResolveABIType(mt.In(j))
			if err != nil {
				log.Info("Skipping method due to input type", "method", m.Name, "err", err)
				continue loop
			}
			// 1-based indexing(j-3) for input arguments, as 0 is reserved for the method receiver
			inputs = append(inputs, abi.Argument{
				Name: fmt.Sprintf("arg%d", j-3),
				Type: abiType,
			})
		}

		// structure outputs
		outputs := abi.Arguments{}
		for j := 0; j < mt.NumOut()-1; j++ {
			abiType, err := ResolveABIType(mt.Out(j))
			if err != nil {
				log.Info(err.Error())
				continue loop
			}
			outputs = append(outputs, abi.Argument{Name: fmt.Sprintf("arg%d", j), Type: abiType, Indexed: false})
		}
		// build signature
		sig := strings.ToLower(m.Name[0:1]) + m.Name[1:] // camelCase
		abiMethod := abi.NewMethod(m.Name, sig, abi.Function, "nonpayable", false, false, inputs, outputs)
		goMethod := contractVal.Method(i)
		log.Info("Go Method name", "name", goMethod.Type().Name())
		d.AddMethod(abiMethod, goMethod)
	}
	return nil
}

func (d *Dispatcher) Dispatch(input []byte, evm interface{}, caller interface{}, storage interface{}) ([]byte, error) {
	if len(input) < 4 {
		return nil, fmt.Errorf("input too short")
	}
	var sel selector
	copy(sel[:], input[:4])

	goMethod, ok := d.Methods[sel]
	if !ok {
		return nil, fmt.Errorf("method not found")
	}

	methodName, ok := d.SelectorToName[sel]
	if !ok {
		return nil, fmt.Errorf("method not found")
	}
	method := d.ABI.Methods[methodName]
	args, err := method.Inputs.Unpack(input[4:])
	if err != nil {
		return nil, err
	}
	in := []reflect.Value{reflect.ValueOf(evm), reflect.ValueOf(caller), reflect.ValueOf(storage)}
	for _, arg := range args {
		in = append(in, reflect.ValueOf(arg))
	}
	out := goMethod.Call(in)

	// error would be last
	errVal := out[len(out)-1]
	if !errVal.IsNil() {
		return nil, errVal.Interface().(error)
	}

	// more outputs apart from error
	if len(out) > 1 {
		retVals := out[:len(out)-1]
		retInterfaces := make([]interface{}, len(retVals))
		for i, rv := range retVals {
			retInterfaces[i] = rv.Interface()
		}
		return method.Outputs.PackValues(retInterfaces)
	}
	// void returns
	return nil, nil
}
