package abiselector

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/core/vm/precompile_sdk/storage"
	"github.com/autonity/autonity/crypto"
)

type mockContract struct {
	called bool
	param1 *big.Int
	param2 common.Address
	ret1   *big.Int
	ret2   bool
}

// Simple void method.
func (m *mockContract) NoParamVoid(_ *vm.EVM, _ common.Address, _ *storage.Storage) error {
	m.called = true
	return nil
}

// Method with uint256 param.
func (m *mockContract) WithUintParam(_ *vm.EVM, _ common.Address, _ *storage.Storage, param *big.Int) error {
	m.param1 = param
	return nil
}

// Method with return.
func (m *mockContract) WithReturn(_ *vm.EVM, _ common.Address, _ *storage.Storage) (*big.Int, error) {
	return m.ret1, nil
}

// Complex: Multi-param, multi-return.
func (m *mockContract) MultiParamReturn(_ *vm.EVM, _ common.Address, _ *storage.Storage, addr common.Address, val *big.Int) (bool, *big.Int, error) {
	m.param2 = addr
	m.param1 = val
	return m.ret2, m.ret1, nil
}

// Unsupported type for error test (e.g., float64 not in map).
func (m *mockContract) UnsupportedParam(_ *vm.EVM, _ common.Address, _ *storage.Storage, f float64) error {
	return nil
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
	dispatcher := NewDispatcher()
	_ = InferABIMethods(dispatcher, reflect.ValueOf(mc))
	//t.Log("Error:", err)

	// Check registered methods.
	expectedMethods := []string{"NoParamVoid", "WithUintParam", "WithReturn", "MultiParamReturn"} // Unsupported skipped.
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

// TestDispatch_SimpleVoid: Call no-param void.
func TestDispatch_SimpleVoid(t *testing.T) {
	d := NewDispatcher()
	mock := &mockContract{}
	_ = InferABIMethods(d, reflect.ValueOf(mock))

	// Compute selector for NoParamVoid().
	method := d.ABI.Methods["NoParamVoid"]
	sel := crypto.Keccak256Hash([]byte(method.Sig)).Bytes()[:4]
	input := sel // No args.

	caller := common.Address{}
	out, err := d.Dispatch(input, caller, mockEVM(), mockStorage())
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

// TestDispatch_WithParam: Uint256 input.
func TestDispatch_WithParam(t *testing.T) {
	d := NewDispatcher()
	mock := &mockContract{}
	InferABIMethods(d, reflect.ValueOf(mock))

	method := d.ABI.Methods["WithUintParam"]
	sel := crypto.Keccak256Hash([]byte(method.Sig)).Bytes()[:4]
	param := big.NewInt(42)
	args, _ := method.Inputs.Pack(param)
	input := append(sel, args...)

	caller := common.Address{}
	_, err := d.Dispatch(input, caller, mockEVM(), mockStorage())
	if err != nil {
		t.Fatalf("dispatch error: %v", err)
	}
	if mock.param1.Cmp(param) != 0 {
		t.Errorf("expected param %v, got %v", param, mock.param1)
	}
}

// TestDispatch_WithReturn: Output packing.
func TestDispatch_WithReturn(t *testing.T) {
	d := NewDispatcher()
	mock := &mockContract{ret1: big.NewInt(100)}
	InferABIMethods(d, reflect.ValueOf(mock))

	method := d.ABI.Methods["WithReturn"]
	sel := crypto.Keccak256Hash([]byte(method.Sig)).Bytes()[:4]
	input := sel

	caller := common.Address{}
	out, err := d.Dispatch(input, caller, mockEVM(), mockStorage())
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

// TestDispatch_MultiParamReturn: Complex case.
func TestDispatch_MultiParamReturn(t *testing.T) {
	d := NewDispatcher()
	mock := &mockContract{ret1: big.NewInt(200), ret2: true}
	InferABIMethods(d, reflect.ValueOf(mock))

	method := d.ABI.Methods["MultiParamReturn"]
	sel := crypto.Keccak256Hash([]byte(method.Sig)).Bytes()[:4]
	addr := common.HexToAddress("0x123")
	val := big.NewInt(300)
	args, _ := method.Inputs.Pack(addr, val)
	input := append(sel, args...)

	out, err := d.Dispatch(input, addr, mockEVM(), mockStorage())
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
