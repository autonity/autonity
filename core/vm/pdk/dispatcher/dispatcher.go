// Package dispatcher handles method dispatching for PDK contracts.
package dispatcher

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
	// GoTypeToABI maps Go types to ABI types.
	GoTypeToABI = map[reflect.Type]func() (abi.Type, error){
		reflect.TypeOf(common.Address{}): func() (abi.Type, error) {
			return abi.NewType("address", "address", nil)
		},
		reflect.TypeOf((*big.Int)(nil)): func() (abi.Type, error) {
			return abi.NewType("int256", "int256", nil)
		},
		reflect.TypeOf(storage.Uint256{}): func() (abi.Type, error) {
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
		reflect.TypeOf(common.Hash{}): func() (abi.Type, error) {
			return abi.NewType("bytes32", "", nil)
		},
	}
)

// Dispatcher handles the mapping of method selectors to Go methods and dispatching calls.
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

// AddMethod registers a new ABI method and its corresponding Go method.
func (d *Dispatcher) AddMethod(method abi.Method, goMethod reflect.Value) {
	d.ABI.Methods[method.Name] = method
	var sel = crypto.Keccak256Hash([]byte(method.Sig)).Bytes()[:4]
	var selArray [4]byte
	copy(selArray[:], sel)
	d.Methods[selArray] = goMethod
	d.SelectorToName[selArray] = method.Name
}

// ResolveABIType converts a Go reflect.Type to a corresponding abi.Type.
func ResolveABIType(goType reflect.Type) (abi.Type, error) {
	// Reject platform-dependent int/uint types
	if goType.Kind() == reflect.Int {
		return abi.Type{}, fmt.Errorf("unsupported type 'int': use *big.Int for int256, or sized types (int8, int16, int32, int64)")
	}
	if goType.Kind() == reflect.Uint {
		return abi.Type{}, fmt.Errorf("unsupported type 'uint': use storage.Uint256 for uint256, or sized types (uint8, uint16, uint32, uint64)")
	}

	mapped, ok := GoTypeToABI[goType]
	if ok {
		return mapped()
	}

	if goType.Kind() == reflect.Ptr {
		return ResolveABIType(goType.Elem())
	}

	// Handle fixed-size byte arrays: [N]byte -> bytesN (Solidity bytes1-bytes32)
	if goType.Kind() == reflect.Array {
		elemType := goType.Elem()
		arrayLen := goType.Len()

		if elemType.Kind() == reflect.Uint8 {
			if arrayLen >= 1 && arrayLen <= 32 {
				typeName := fmt.Sprintf("bytes%d", arrayLen)
				return abi.NewType(typeName, "", nil)
			}
			return abi.Type{}, fmt.Errorf("fixed byte array size must be 1-32, got [%d]byte", arrayLen)
		}

		// Generic fixed-size arrays: [N]T -> T[N]
		elemABIType, err := ResolveABIType(elemType)
		if err != nil {
			return abi.Type{}, fmt.Errorf("failed to resolve array element type %s: %w", elemType, err)
		}
		// Construct array type string: "uint256[2]", "address[10]", etc.
		// Use abi.NewType logic to parse it
		typeName := fmt.Sprintf("%s[%d]", elemABIType.String(), arrayLen)
		return abi.NewType(typeName, typeName, nil)
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

func convertTupleElemsToArgumentMarshaling(tupleElems []*abi.Type) []abi.ArgumentMarshaling {
	if tupleElems == nil {
		return nil
	}

	components := make([]abi.ArgumentMarshaling, 0, len(tupleElems))
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

// InferABIMethods inspects a contract struct and registers its exported methods to the dispatcher.
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
		if mt.NumIn() < 4 || mt.NumOut() < 1 ||
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
			paramType := mt.In(j)
			abiType, err := ResolveABIType(paramType)
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

// Dispatch executes the mapped Go method for the given input data.
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

	methodType := goMethod.Type()
	in := []reflect.Value{reflect.ValueOf(evm), reflect.ValueOf(caller), reflect.ValueOf(storage)}

	// Convert unpacked args to expected types
	// User params start at In(3): In(0)=EVM, In(1)=caller, In(2)=storage, In(3+)=user params
	for i, arg := range args {
		expectedType := methodType.In(i + 3) // +3 to skip evm, caller, storage (no receiver in bound method type)
		converted, err := convertToTargetType(arg, expectedType)
		if err != nil {
			return nil, fmt.Errorf("failed to convert argument %d: %w", i, err)
		}
		in = append(in, converted)
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

func convertToTargetType(value interface{}, targetType reflect.Type) (reflect.Value, error) {
	srcVal := reflect.ValueOf(value)
	if !srcVal.IsValid() {
		return reflect.Zero(targetType), nil
	}

	srcType := srcVal.Type()

	if srcType == targetType {
		return srcVal, nil
	}

	if srcType.AssignableTo(targetType) {
		return srcVal, nil
	}

	// Handle target pointer wrapping: method expects *T but we have T
	needsPointer := targetType.Kind() == reflect.Ptr
	actualTarget := targetType
	if needsPointer {
		actualTarget = targetType.Elem()
	}

	wrap := func(v reflect.Value) reflect.Value {
		if needsPointer {
			ptr := reflect.New(actualTarget)
			ptr.Elem().Set(v)
			return ptr
		}
		return v
	}

	// Check identity/assignability after considering pointer wrapping
	if srcType == actualTarget {
		return wrap(srcVal), nil
	}

	if srcType.AssignableTo(actualTarget) {
		return wrap(srcVal), nil
	}

	// Handle *big.Int -> storage.Uint256 conversion
	// ABI decodes 'uint256' as *big.Int. If target is storage.Uint256, convert it.
	if srcType == reflect.TypeOf((*big.Int)(nil)) && actualTarget == reflect.TypeOf(storage.Uint256{}) {
		bi := value.(*big.Int)
		if bi.Sign() < 0 {
			return reflect.Value{}, fmt.Errorf("cannot convert negative *big.Int to Uint256")
		}
		// Use storage helper to create Uint256
		u := storage.NewUint256FromBig(bi)
		return wrap(reflect.ValueOf(u)), nil
	}

	// *big.Int to sized integer conversion
	if srcType == reflect.TypeOf((*big.Int)(nil)) && canConvertBigInt(actualTarget) {
		bi := value.(*big.Int)
		converted, err := convertBigInt(bi, actualTarget)
		if err != nil {
			return reflect.Value{}, err
		}
		return wrap(converted), nil
	}

	// Direct type conversion (e.g., int32 -> int64, uint8 -> uint16)
	if srcType.ConvertibleTo(actualTarget) {
		return wrap(srcVal.Convert(actualTarget)), nil
	}

	// Struct conversion (ABI tuple -> Go struct)
	if srcType.Kind() == reflect.Struct && actualTarget.Kind() == reflect.Struct {
		converted, err := convertStruct(srcVal, actualTarget)
		if err != nil {
			return reflect.Value{}, err
		}
		return wrap(converted), nil
	}

	// Slice conversion (element-by-element)
	if srcType.Kind() == reflect.Slice && actualTarget.Kind() == reflect.Slice {
		converted, err := convertSlice(srcVal, actualTarget)
		if err != nil {
			return reflect.Value{}, err
		}
		return wrap(converted), nil
	}

	if srcType.Kind() == reflect.Slice && actualTarget.Kind() == reflect.Array {
		converted, err := convertArray(srcVal, actualTarget)
		if err != nil {
			return reflect.Value{}, err
		}
		return wrap(converted), nil
	}

	return reflect.Value{}, fmt.Errorf("cannot convert %v to %v", srcType, targetType)
}

func canConvertBigInt(t reflect.Type) bool {
	k := t.Kind()
	return k == reflect.Int8 || k == reflect.Int16 || k == reflect.Int32 || k == reflect.Int64 ||
		k == reflect.Uint8 || k == reflect.Uint16 || k == reflect.Uint32 || k == reflect.Uint64
}

func convertBigInt(bi *big.Int, t reflect.Type) (reflect.Value, error) {
	// Note: checking for overflow would be ideal here.
	// For now, we follow standard Go conversion semantics (truncation/wrapping).

	switch t.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if !bi.IsInt64() {
			return reflect.Value{}, fmt.Errorf("value %s too large for %s (max: int64)", bi, t)
		}
		return reflect.ValueOf(bi.Int64()).Convert(t), nil
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if !bi.IsUint64() {
			return reflect.Value{}, fmt.Errorf("value %s too large for %s (max: uint64)", bi, t)
		}
		return reflect.ValueOf(bi.Uint64()).Convert(t), nil
	}
	return reflect.Value{}, fmt.Errorf("unsupported integer type %s: use *big.Int or sized integer types", t)
}

func convertStruct(src reflect.Value, targetType reflect.Type) (reflect.Value, error) {
	srcType := src.Type()

	if srcType == targetType {
		return src, nil
	}

	if src.Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("source is not a struct: %v", srcType)
	}
	if targetType.Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("target is not a struct: %v", targetType)
	}

	// Field count must match exactly for tuple unpacking
	// ABI tuples unpack to structs with matching field count
	srcFieldCount := src.NumField()
	dstFieldCount := targetType.NumField()
	if srcFieldCount != dstFieldCount {
		return reflect.Value{}, fmt.Errorf(
			"struct field count mismatch: source has %d fields, target has %d",
			srcFieldCount, dstFieldCount)
	}

	// Create destination struct and convert each field
	dst := reflect.New(targetType).Elem()
	for i := 0; i < dstFieldCount; i++ {
		srcField := src.Field(i)
		dstField := dst.Field(i)

		if !dstField.CanSet() {
			return reflect.Value{}, fmt.Errorf("target field %d (%s) is not settable",
				i, targetType.Field(i).Name)
		}

		// Recursively convert field value
		converted, err := convertToTargetType(srcField.Interface(), dstField.Type())
		if err != nil {
			return reflect.Value{}, fmt.Errorf("field %d (%s -> %s): %w",
				i, srcType.Field(i).Name, targetType.Field(i).Name, err)
		}

		dstField.Set(converted)
	}

	return dst, nil
}

// convertSlice handles slice conversion element by element
func convertSlice(src reflect.Value, targetType reflect.Type) (reflect.Value, error) {
	srcType := src.Type()

	if srcType == targetType {
		return src, nil
	}

	if src.Kind() != reflect.Slice {
		return reflect.Value{}, fmt.Errorf("source is not a slice: %v", srcType)
	}
	if targetType.Kind() != reflect.Slice {
		return reflect.Value{}, fmt.Errorf("target is not a slice: %v", targetType)
	}

	srcElemType := srcType.Elem()
	dstElemType := targetType.Elem()
	if srcElemType == dstElemType {
		return src, nil
	}

	length := src.Len()
	result := reflect.MakeSlice(targetType, length, length)

	for i := 0; i < length; i++ {
		elem := src.Index(i)
		converted, err := convertToTargetType(elem.Interface(), dstElemType)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("element %d (%v -> %v): %w",
				i, srcElemType, dstElemType, err)
		}
		result.Index(i).Set(converted)
	}
	return result, nil
}

func convertArray(src reflect.Value, targetType reflect.Type) (reflect.Value, error) {
	if src.Kind() != reflect.Slice {
		return reflect.Value{}, fmt.Errorf("source is not a slice: %v", src.Type())
	}
	if targetType.Kind() != reflect.Array {
		return reflect.Value{}, fmt.Errorf("target is not an array: %v", targetType)
	}

	srcLen := src.Len()
	dstLen := targetType.Len()

	if srcLen != dstLen {
		return reflect.Value{}, fmt.Errorf("length mismatch: source slice len %d, target array len %d", srcLen, dstLen)
	}

	dst := reflect.New(targetType).Elem()
	dstElemType := targetType.Elem()

	for i := 0; i < dstLen; i++ {
		elem := src.Index(i)
		converted, err := convertToTargetType(elem.Interface(), dstElemType)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("element %d: %w", i, err)
		}
		dst.Index(i).Set(converted)
	}

	return dst, nil
}
