package abiselector

import (
	"fmt"
	"math/big"
	"reflect"
	"strings"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/core/vm/precompile_sdk/storage"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/log"
)

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
		reflect.TypeOf([32]byte{}): func() (abi.Type, error) {
			return abi.NewType("bytes32", "", nil)
		},
		//todo: more types and structs
	}
)

type Dispatcher struct {
	ABI            abi.ABI
	Methods        map[[4]byte]reflect.Value // selector registry
	SelectorToName map[[4]byte]string
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		ABI:            abi.ABI{Methods: make(map[string]abi.Method)},
		Methods:        map[[4]byte]reflect.Value{},
		SelectorToName: map[[4]byte]string{},
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
	if !ok {
		// check if this is a dynamic type
		if goType.Kind() == reflect.Slice {
			elemType := goType.Elem()
			elemMapped, elemOk := GoTypeToABI[elemType]
			if elemOk {
				abiElemType, err := elemMapped()
				if err != nil {
					log.Info(err.Error())
					return abiElemType, err
				}
				// construct dynamic array type
				arrayTypeStr := abiElemType.String() + "[]"
				// udpate mapped
				mapped = func() (abi.Type, error) {
					return abi.NewType(arrayTypeStr, arrayTypeStr, nil)
				}
				return mapped()
			}
		}
		err := fmt.Errorf("unsupported input type %s", goType.String())
		return abi.Type{}, err
	}
	return mapped()
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
		//todo(piyush): merge checks
		if mt.NumIn() < 3 || mt.NumOut() < 1 {
			log.Info("Method %s has incompatible signature, skipping", m.Name)
			continue
		}
		if mt.In(0) != contractType {
			log.Info("Method %s has incompatible signature, skipping", m.Name)
			continue
		} // first arg is receiver - Contract
		if mt.In(1) != reflect.TypeOf((*vm.EVM)(nil)) {
			log.Info("Method %s has incompatible signature, skipping", m.Name)
			continue
		} // second arg is *vm.EVM

		if mt.In(2) != reflect.TypeOf(common.Address{}) {
			log.Info("Method %s has incompatible signature, skipping", m.Name)
			continue
		} // third arg is caller address

		if mt.In(3) != reflect.TypeOf((*storage.Storage)(nil)) {
			log.Info("Method %s has incompatible signature, skipping", m.Name)
			continue
		} // 4th arg is storage

		// last out argument must be error
		if mt.Out(mt.NumOut()-1) != reflect.TypeOf((*error)(nil)).Elem() {
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
		// review:
		abiMethod := abi.NewMethod(m.Name, sig, abi.Function, "nonpayable", false, false, inputs, outputs)
		goMethod := contractVal.Method(i)
		log.Info("Go Method name", "name", goMethod.Type().Name())
		d.AddMethod(abiMethod, goMethod)
	}
	return nil
}

func (d *Dispatcher) Dispatch(input []byte, evm interface{}, caller interface{}, storage interface{}) ([]byte, error) {
	if len(input) < 4 {
		return nil, fmt.Errorf("Input too short")
	}
	var selector [4]byte
	copy(selector[:], input[:4])

	goMethod, ok := d.Methods[selector]
	if !ok {
		return nil, fmt.Errorf("Method not found")
	}

	methodName, ok := d.SelectorToName[selector]
	if !ok {
		return nil, fmt.Errorf("Method not found")
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
