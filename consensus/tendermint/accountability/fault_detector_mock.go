// Code generated by MockGen. DO NOT EDIT.
// Source: consensus/tendermint/accountability/fault_detector.go
//
// Generated by this command:
//
//	mockgen -source=consensus/tendermint/accountability/fault_detector.go -package=accountability -destination=consensus/tendermint/accountability/fault_detector_mock.go
//

// Package accountability is a generated GoMock package.
package accountability

import (
	big "math/big"
	reflect "reflect"

	autonity "github.com/autonity/autonity/autonity"
	common "github.com/autonity/autonity/common"
	consensus "github.com/autonity/autonity/consensus"
	core "github.com/autonity/autonity/core"
	state "github.com/autonity/autonity/core/state"
	types "github.com/autonity/autonity/core/types"
	event "github.com/autonity/autonity/event"
	params "github.com/autonity/autonity/params"
	gomock "go.uber.org/mock/gomock"
)

// MockChainContext is a mock of ChainContext interface.
type MockChainContext struct {
	ctrl     *gomock.Controller
	recorder *MockChainContextMockRecorder
}

// MockChainContextMockRecorder is the mock recorder for MockChainContext.
type MockChainContextMockRecorder struct {
	mock *MockChainContext
}

// NewMockChainContext creates a new mock instance.
func NewMockChainContext(ctrl *gomock.Controller) *MockChainContext {
	mock := &MockChainContext{ctrl: ctrl}
	mock.recorder = &MockChainContextMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChainContext) EXPECT() *MockChainContextMockRecorder {
	return m.recorder
}

// CommitteeOfHeight mocks base method.
func (m *MockChainContext) CommitteeOfHeight(height uint64) (*types.Committee, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CommitteeOfHeight", height)
	ret0, _ := ret[0].(*types.Committee)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CommitteeOfHeight indicates an expected call of CommitteeOfHeight.
func (mr *MockChainContextMockRecorder) CommitteeOfHeight(height any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CommitteeOfHeight", reflect.TypeOf((*MockChainContext)(nil).CommitteeOfHeight), height)
}

// Config mocks base method.
func (m *MockChainContext) Config() *params.ChainConfig {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Config")
	ret0, _ := ret[0].(*params.ChainConfig)
	return ret0
}

// Config indicates an expected call of Config.
func (mr *MockChainContextMockRecorder) Config() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Config", reflect.TypeOf((*MockChainContext)(nil).Config))
}

// CurrentBlock mocks base method.
func (m *MockChainContext) CurrentBlock() *types.Block {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CurrentBlock")
	ret0, _ := ret[0].(*types.Block)
	return ret0
}

// CurrentBlock indicates an expected call of CurrentBlock.
func (mr *MockChainContextMockRecorder) CurrentBlock() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CurrentBlock", reflect.TypeOf((*MockChainContext)(nil).CurrentBlock))
}

// CurrentHeader mocks base method.
func (m *MockChainContext) CurrentHeader() *types.Header {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CurrentHeader")
	ret0, _ := ret[0].(*types.Header)
	return ret0
}

// CurrentHeader indicates an expected call of CurrentHeader.
func (mr *MockChainContextMockRecorder) CurrentHeader() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CurrentHeader", reflect.TypeOf((*MockChainContext)(nil).CurrentHeader))
}

// Engine mocks base method.
func (m *MockChainContext) Engine() consensus.Engine {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Engine")
	ret0, _ := ret[0].(consensus.Engine)
	return ret0
}

// Engine indicates an expected call of Engine.
func (mr *MockChainContextMockRecorder) Engine() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Engine", reflect.TypeOf((*MockChainContext)(nil).Engine))
}

// EpochOfHeight mocks base method.
func (m *MockChainContext) EpochOfHeight(height uint64) (*types.EpochInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EpochOfHeight", height)
	ret0, _ := ret[0].(*types.EpochInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EpochOfHeight indicates an expected call of EpochOfHeight.
func (mr *MockChainContextMockRecorder) EpochOfHeight(height any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EpochOfHeight", reflect.TypeOf((*MockChainContext)(nil).EpochOfHeight), height)
}

// GetBlock mocks base method.
func (m *MockChainContext) GetBlock(hash common.Hash, number uint64) *types.Block {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlock", hash, number)
	ret0, _ := ret[0].(*types.Block)
	return ret0
}

// GetBlock indicates an expected call of GetBlock.
func (mr *MockChainContextMockRecorder) GetBlock(hash, number any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlock", reflect.TypeOf((*MockChainContext)(nil).GetBlock), hash, number)
}

// GetHeader mocks base method.
func (m *MockChainContext) GetHeader(hash common.Hash, number uint64) *types.Header {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHeader", hash, number)
	ret0, _ := ret[0].(*types.Header)
	return ret0
}

// GetHeader indicates an expected call of GetHeader.
func (mr *MockChainContextMockRecorder) GetHeader(hash, number any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHeader", reflect.TypeOf((*MockChainContext)(nil).GetHeader), hash, number)
}

// GetHeaderByHash mocks base method.
func (m *MockChainContext) GetHeaderByHash(hash common.Hash) *types.Header {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHeaderByHash", hash)
	ret0, _ := ret[0].(*types.Header)
	return ret0
}

// GetHeaderByHash indicates an expected call of GetHeaderByHash.
func (mr *MockChainContextMockRecorder) GetHeaderByHash(hash any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHeaderByHash", reflect.TypeOf((*MockChainContext)(nil).GetHeaderByHash), hash)
}

// GetHeaderByNumber mocks base method.
func (m *MockChainContext) GetHeaderByNumber(number uint64) *types.Header {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHeaderByNumber", number)
	ret0, _ := ret[0].(*types.Header)
	return ret0
}

// GetHeaderByNumber indicates an expected call of GetHeaderByNumber.
func (mr *MockChainContextMockRecorder) GetHeaderByNumber(number any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHeaderByNumber", reflect.TypeOf((*MockChainContext)(nil).GetHeaderByNumber), number)
}

// GetTd mocks base method.
func (m *MockChainContext) GetTd(hash common.Hash, number uint64) *big.Int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTd", hash, number)
	ret0, _ := ret[0].(*big.Int)
	return ret0
}

// GetTd indicates an expected call of GetTd.
func (mr *MockChainContextMockRecorder) GetTd(hash, number any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTd", reflect.TypeOf((*MockChainContext)(nil).GetTd), hash, number)
}

// HasBadBlock mocks base method.
func (m *MockChainContext) HasBadBlock(hash common.Hash) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasBadBlock", hash)
	ret0, _ := ret[0].(bool)
	return ret0
}

// HasBadBlock indicates an expected call of HasBadBlock.
func (mr *MockChainContextMockRecorder) HasBadBlock(hash any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasBadBlock", reflect.TypeOf((*MockChainContext)(nil).HasBadBlock), hash)
}

// MinBaseFee mocks base method.
func (m *MockChainContext) MinBaseFee() *big.Int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MinBaseFee")
	ret0, _ := ret[0].(*big.Int)
	return ret0
}

// MinBaseFee indicates an expected call of MinBaseFee.
func (mr *MockChainContextMockRecorder) MinBaseFee() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MinBaseFee", reflect.TypeOf((*MockChainContext)(nil).MinBaseFee))
}

// ProtocolContracts mocks base method.
func (m *MockChainContext) ProtocolContracts() *autonity.ProtocolContracts {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProtocolContracts")
	ret0, _ := ret[0].(*autonity.ProtocolContracts)
	return ret0
}

// ProtocolContracts indicates an expected call of ProtocolContracts.
func (mr *MockChainContextMockRecorder) ProtocolContracts() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProtocolContracts", reflect.TypeOf((*MockChainContext)(nil).ProtocolContracts))
}

// State mocks base method.
func (m *MockChainContext) State() (*state.StateDB, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "State")
	ret0, _ := ret[0].(*state.StateDB)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// State indicates an expected call of State.
func (mr *MockChainContextMockRecorder) State() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "State", reflect.TypeOf((*MockChainContext)(nil).State))
}

// StateAt mocks base method.
func (m *MockChainContext) StateAt(root common.Hash) (*state.StateDB, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StateAt", root)
	ret0, _ := ret[0].(*state.StateDB)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StateAt indicates an expected call of StateAt.
func (mr *MockChainContextMockRecorder) StateAt(root any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StateAt", reflect.TypeOf((*MockChainContext)(nil).StateAt), root)
}

// SubscribeChainEvent mocks base method.
func (m *MockChainContext) SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubscribeChainEvent", ch)
	ret0, _ := ret[0].(event.Subscription)
	return ret0
}

// SubscribeChainEvent indicates an expected call of SubscribeChainEvent.
func (mr *MockChainContextMockRecorder) SubscribeChainEvent(ch any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeChainEvent", reflect.TypeOf((*MockChainContext)(nil).SubscribeChainEvent), ch)
}

// Validator mocks base method.
func (m *MockChainContext) Validator() core.Validator {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Validator")
	ret0, _ := ret[0].(core.Validator)
	return ret0
}

// Validator indicates an expected call of Validator.
func (mr *MockChainContextMockRecorder) Validator() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validator", reflect.TypeOf((*MockChainContext)(nil).Validator))
}
