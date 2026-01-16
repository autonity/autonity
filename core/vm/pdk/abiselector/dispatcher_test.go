package abiselector

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/core/vm/pdk/storage"
	"github.com/autonity/autonity/crypto"
)

// Test struct types for ABI tuple mapping
type TestStruct struct {
	ID      *big.Int
	Address common.Address
	Amount  *big.Int
	Active  bool
	Data    []byte
}

type NestedStruct struct {
	Inner  TestStruct
	Values []*big.Int
	Name   string
}

type ByteArrayStruct struct {
	Hash32   [32]byte
	Address  [20]byte
	Selector [4]byte
}

type mockContract struct {
	called      bool
	param1      *big.Int
	param2      common.Address
	ret1        *big.Int
	ret2        bool
	structParam TestStruct
	nestedParam NestedStruct
	byteParam   ByteArrayStruct
	structRet   TestStruct
	sliceParam  []TestStruct
}

func (m *mockContract) NoParamVoid(_ *vm.EVM, _ common.Address, _ *storage.Storage) error {
	m.called = true
	return nil
}

func (m *mockContract) WithUintParam(_ *vm.EVM, _ common.Address, _ *storage.Storage, param *big.Int) error {
	m.param1 = param
	return nil
}

func (m *mockContract) WithReturn(_ *vm.EVM, _ common.Address, _ *storage.Storage) (*big.Int, error) {
	return m.ret1, nil
}

func (m *mockContract) MultiParamReturn(_ *vm.EVM, _ common.Address, _ *storage.Storage, addr common.Address, val *big.Int) (bool, *big.Int, error) {
	m.param2 = addr
	m.param1 = val
	return m.ret2, m.ret1, nil
}

func (m *mockContract) UnsupportedParam(_ *vm.EVM, _ common.Address, _ *storage.Storage, _ float64) error {
	return nil
}

func (m *mockContract) WithStructParam(_ *vm.EVM, _ common.Address, _ *storage.Storage, param TestStruct) error {
	m.structParam = param
	return nil
}

func (m *mockContract) WithNestedStructParam(_ *vm.EVM, _ common.Address, _ *storage.Storage, param NestedStruct) error {
	m.nestedParam = param
	return nil
}

func (m *mockContract) WithByteArrayStruct(_ *vm.EVM, _ common.Address, _ *storage.Storage, param ByteArrayStruct) error {
	m.byteParam = param
	return nil
}

func (m *mockContract) WithStructReturn(_ *vm.EVM, _ common.Address, _ *storage.Storage) (TestStruct, error) {
	return m.structRet, nil
}

func (m *mockContract) WithStructSlice(_ *vm.EVM, _ common.Address, _ *storage.Storage, params []TestStruct) ([]TestStruct, error) {
	m.sliceParam = params
	return params, nil
}

func mockEVM() *vm.EVM {
	return &vm.EVM{Context: vm.BlockContext{}}
}

func mockStorage() *storage.Storage {
	return &storage.Storage{}
}

func TestInferABIMethods(t *testing.T) {
	mc := &mockContract{
		ret1: big.NewInt(42),
		ret2: true,
	}
	dispatcher := newDispatcher()
	_ = InferABIMethods(dispatcher, reflect.ValueOf(mc))
	//t.Log("Error:", err)

	// Check registered methods.
	expectedMethods := []string{
		"NoParamVoid", "WithUintParam", "WithReturn", "MultiParamReturn",
		"WithStructParam", "WithNestedStructParam", "WithByteArrayStruct",
		"WithStructReturn", "WithStructSlice",
	} // Unsupported skipped.
	if len(dispatcher.ABI.Methods) != len(expectedMethods) {
		t.Errorf("expected %d methods, got %d", len(expectedMethods), len(dispatcher.ABI.Methods))
	}

	for _, name := range expectedMethods {
		if _, ok := dispatcher.ABI.Methods[name]; !ok {
			t.Errorf("missing method %s", name)
		}
	}

	// Check unsupported skipped.
	if _, ok := dispatcher.ABI.Methods["UnsupportedParam"]; ok {
		t.Error("unsupported method should not register")
	}
}

func TestDispatch_SimpleVoid(t *testing.T) {
	d := newDispatcher()
	mock := &mockContract{}
	_ = InferABIMethods(d, reflect.ValueOf(mock))

	// Compute selector for NoParamVoid().
	method := d.ABI.Methods["NoParamVoid"]
	sel := crypto.Keccak256Hash([]byte(method.Sig)).Bytes()[:4]
	input := sel // No args.

	caller := common.Address{}
	out, err := d.Dispatch(input, mockEVM(), caller, mockStorage())
	if err != nil {
		t.Fatalf("dispatch error: %v", err)
	}
	if len(out) != 0 {
		t.Error("expected empty output for void")
	}
	if !mock.called {
		t.Error("method not called")
	}
}

func TestDispatch_WithParam(t *testing.T) {
	d := newDispatcher()
	mock := &mockContract{}
	err := InferABIMethods(d, reflect.ValueOf(mock))
	require.NoError(t, err)

	method := d.ABI.Methods["WithUintParam"]
	sel := crypto.Keccak256Hash([]byte(method.Sig)).Bytes()[:4]
	param := big.NewInt(42)
	args, _ := method.Inputs.Pack(param)
	input := append(sel, args...)

	caller := common.Address{}
	_, err = d.Dispatch(input, mockEVM(), caller, mockStorage())
	require.NoError(t, err)
	if mock.param1.Cmp(param) != 0 {
		t.Errorf("expected param %v, got %v", param, mock.param1)
	}
}

func TestDispatch_WithReturn(t *testing.T) {
	d := newDispatcher()
	mock := &mockContract{ret1: big.NewInt(100)}
	InferABIMethods(d, reflect.ValueOf(mock))

	method := d.ABI.Methods["WithReturn"]
	sel := crypto.Keccak256Hash([]byte(method.Sig)).Bytes()[:4]
	input := sel

	caller := common.Address{}
	out, err := d.Dispatch(input, mockEVM(), caller, mockStorage())
	if err != nil {
		t.Fatalf("dispatch error: %v", err)
	}

	ret, err := method.Outputs.Unpack(out)
	if err != nil || len(ret) != 1 {
		t.Fatal("unpack failed")
	}
	if retBig, ok := ret[0].(*big.Int); !ok || retBig.Cmp(mock.ret1) != 0 {
		t.Errorf("expected return %v, got %v", mock.ret1, retBig)
	}
}

func TestDispatch_MultiParamReturn(t *testing.T) {
	d := newDispatcher()
	mock := &mockContract{ret1: big.NewInt(200), ret2: true}
	InferABIMethods(d, reflect.ValueOf(mock))

	method := d.ABI.Methods["MultiParamReturn"]
	sel := crypto.Keccak256Hash([]byte(method.Sig)).Bytes()[:4]
	addr := common.HexToAddress("0x123")
	val := big.NewInt(300)
	args, _ := method.Inputs.Pack(addr, val)
	input := append(sel, args...)

	out, err := d.Dispatch(input, mockEVM(), addr, mockStorage())
	if err != nil {
		t.Fatalf("dispatch error: %v", err)
	}

	if mock.param2 != addr || mock.param1.Cmp(val) != 0 {
		t.Error("params not set")
	}

	ret, err := method.Outputs.Unpack(out)
	if err != nil || len(ret) != 2 {
		t.Fatal("unpack failed")
	}
	if retBool, ok := ret[0].(bool); !ok || retBool != true {
		t.Error("bool return mismatch")
	}
	if retBig, ok := ret[1].(*big.Int); !ok || retBig.Cmp(big.NewInt(200)) != 0 {
		t.Error("big.Int return mismatch")
	}
}

func TestResolveABIType_BasicStruct(t *testing.T) {
	structType := reflect.TypeOf(TestStruct{})
	abiType, err := ResolveABIType(structType)
	require.NoError(t, err)

	expected := "(uint256,address,uint256,bool,bytes)"
	require.Equal(t, expected, abiType.String())
}

func TestResolveABIType_NestedStruct(t *testing.T) {
	structType := reflect.TypeOf(NestedStruct{})
	abiType, err := ResolveABIType(structType)
	require.NoError(t, err)

	expected := "((uint256,address,uint256,bool,bytes),uint256[],string)"
	require.Equal(t, expected, abiType.String())
}

func TestResolveABIType_ByteArrayStruct(t *testing.T) {
	structType := reflect.TypeOf(ByteArrayStruct{})
	abiType, err := ResolveABIType(structType)
	require.NoError(t, err)

	expected := "(bytes32,bytes20,bytes4)"
	require.Equal(t, expected, abiType.String())
}

func TestResolveABIType_StructSlice(t *testing.T) {
	sliceType := reflect.TypeOf([]TestStruct{})
	abiType, err := ResolveABIType(sliceType)
	require.NoError(t, err)

	expected := "(uint256,address,uint256,bool,bytes)[]"
	require.Equal(t, expected, abiType.String())
}

func TestResolveABIType_NewByteTypes(t *testing.T) {
	tests := []struct {
		name     string
		goType   reflect.Type
		expected string
	}{
		{"bytes32", reflect.TypeOf([32]byte{}), "bytes32"},
		{"bytes20", reflect.TypeOf([20]byte{}), "bytes20"},
		{"bytes4", reflect.TypeOf([4]byte{}), "bytes4"},
		{"common.Hash", reflect.TypeOf(common.Hash{}), "bytes32"},
		{"uint", reflect.TypeOf(uint(0)), "uint256"},
		{"int", reflect.TypeOf(int(0)), "int256"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			abiType, err := ResolveABIType(tt.goType)
			require.NoError(t, err)
			require.Equal(t, tt.expected, abiType.String())
		})
	}
}

func TestResolveABIType_UnsupportedTypes(t *testing.T) {
	tests := []struct {
		name   string
		goType reflect.Type
	}{
		{"fixed array", reflect.TypeOf([5]int{})},
		{"map", reflect.TypeOf(map[string]int{})},
		{"channel", reflect.TypeOf(make(chan int))},
		{"function", reflect.TypeOf(func() {})},
		{"interface", reflect.TypeOf((*interface{})(nil)).Elem()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ResolveABIType(tt.goType)
			require.Error(t, err)
			require.Contains(t, err.Error(), "unsupported type")
		})
	}
}

// Integration test for struct ABI generation
func TestInferABIMethods_StructTypes(t *testing.T) {
	d := newDispatcher()
	mock := &mockContract{}
	err := InferABIMethods(d, reflect.ValueOf(mock))
	require.NoError(t, err)

	// Test that struct methods have correct ABI signatures
	testCases := []struct {
		methodName string
		signature  string
	}{
		{
			"WithStructParam",
			"withStructParam((uint256,address,uint256,bool,bytes))",
		},
		{
			"WithNestedStructParam",
			"withNestedStructParam(((uint256,address,uint256,bool,bytes),uint256[],string))",
		},
		{
			"WithByteArrayStruct",
			"withByteArrayStruct((bytes32,bytes20,bytes4))",
		},
		{
			"WithStructReturn",
			"withStructReturn()",
		},
		{
			"WithStructSlice",
			"withStructSlice((uint256,address,uint256,bool,bytes)[])",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.methodName, func(t *testing.T) {
			method, exists := d.ABI.Methods[tc.methodName]
			require.True(t, exists, "Method %s should exist", tc.methodName)
			require.Equal(t, tc.signature, method.Sig, "Method signature mismatch")
		})
	}
}
