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

type mockContract struct {
	called bool
	param1 *big.Int
	param2 common.Address
	ret1   *big.Int
	ret2   bool
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
