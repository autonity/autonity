package misbehaviourdetector

import (
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/state"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/params"
	"github.com/golang/mock/gomock"
	"math/big"
	"reflect"
)

// MockBlockChainContext is a mock of ChainReader interface
type MockBlockChainContext struct {
	ctrl     *gomock.Controller
	recorder *MockBlockChainContextMockRecorder
}

func (m *MockBlockChainContext) GetTd(hash common.Hash, number uint64) *big.Int {
	//TODO implement me
	panic("implement me")
}

func (m *MockBlockChainContext) GetMinBaseFee(header *types.Header) (*big.Int, error) {
	//TODO implement me
	panic("implement me")
}

// MockBlockChainContextMockRecorder is the mock recorder for MockBlockChainContext
type MockBlockChainContextMockRecorder struct {
	mock *MockBlockChainContext
}

// NewMockBlockChainContext creates a new mock instance
func NewMockBlockChainContext(ctrl *gomock.Controller) *MockBlockChainContext {
	mock := &MockBlockChainContext{ctrl: ctrl}
	mock.recorder = &MockBlockChainContextMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBlockChainContext) EXPECT() *MockBlockChainContextMockRecorder {
	return m.recorder
}

// Config mocks base method
func (m *MockBlockChainContext) Config() *params.ChainConfig {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Config")
	ret0, _ := ret[0].(*params.ChainConfig)
	return ret0
}

// Config indicates an expected call of Config
func (mr *MockBlockChainContextMockRecorder) Config() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Config", reflect.TypeOf((*MockBlockChainContext)(nil).Config))
}

// CurrentHeader mocks base method
func (m *MockBlockChainContext) CurrentHeader() *types.Header {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CurrentHeader")
	ret0, _ := ret[0].(*types.Header)
	return ret0
}

// CurrentHeader indicates an expected call of CurrentHeader
func (mr *MockBlockChainContextMockRecorder) CurrentHeader() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CurrentHeader", reflect.TypeOf((*MockBlockChainContext)(nil).CurrentHeader))
}

// GetHeader mocks base method
func (m *MockBlockChainContext) GetHeader(hash common.Hash, number uint64) *types.Header {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHeader", hash, number)
	ret0, _ := ret[0].(*types.Header)
	return ret0
}

// GetHeader indicates an expected call of GetHeader
func (mr *MockBlockChainContextMockRecorder) GetHeader(hash, number interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHeader", reflect.TypeOf((*MockBlockChainContext)(nil).GetHeader), hash, number)
}

// GetHeaderByNumber mocks base method
func (m *MockBlockChainContext) GetHeaderByNumber(number uint64) *types.Header {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHeaderByNumber", number)
	ret0, _ := ret[0].(*types.Header)
	return ret0
}

// GetHeaderByNumber indicates an expected call of GetHeaderByNumber
func (mr *MockBlockChainContextMockRecorder) GetHeaderByNumber(number interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHeaderByNumber", reflect.TypeOf((*MockBlockChainContext)(nil).GetHeaderByNumber), number)
}

// GetHeaderByHash mocks base method
func (m *MockBlockChainContext) GetHeaderByHash(hash common.Hash) *types.Header {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHeaderByHash", hash)
	ret0, _ := ret[0].(*types.Header)
	return ret0
}

// GetHeaderByHash indicates an expected call of GetHeaderByHash
func (mr *MockBlockChainContextMockRecorder) GetHeaderByHash(hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHeaderByHash", reflect.TypeOf((*MockBlockChainContext)(nil).GetHeaderByHash), hash)
}

// GetBlock mocks base method
func (m *MockBlockChainContext) GetBlock(hash common.Hash, number uint64) *types.Block {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlock", hash, number)
	ret0, _ := ret[0].(*types.Block)
	return ret0
}

// GetBlock indicates an expected call of GetBlock
func (mr *MockBlockChainContextMockRecorder) GetBlock(hash, number interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlock", reflect.TypeOf((*MockBlockChainContext)(nil).GetBlock), hash, number)
}

// Engine mocks base method
func (m *MockBlockChainContext) Engine() consensus.Engine {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Engine")
	ret0, _ := ret[0].(consensus.Engine)
	return ret0
}

// Engine indicates an expected call of Engine
func (mr *MockBlockChainContextMockRecorder) Engine() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Engine", reflect.TypeOf((*MockBlockChainContext)(nil).Engine))
}

func (m *MockBlockChainContext) CurrentBlock() *types.Block {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CurrentBlock")
	ret0, _ := ret[0].(*types.Block)
	return ret0
}

// CurrentBlock Engine indicates an expected call of Engine
func (mr *MockBlockChainContextMockRecorder) CurrentBlock() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CurrentBlock", reflect.TypeOf((*MockBlockChainContext)(nil).CurrentBlock))
}

func (m *MockBlockChainContext) SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubscribeChainEvent")
	ret0, _ := ret[0].(event.Subscription)
	return ret0
}

func (mr *MockBlockChainContextMockRecorder) SubscribeChainEvent() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeChainEvent", reflect.TypeOf((*MockBlockChainContext)(nil).SubscribeChainEvent))
}

func (m *MockBlockChainContext) State() (*state.StateDB, error) {
	panic("implement me")
}

func (m *MockBlockChainContext) GetAutonityContract() *autonity.Contract {
	panic("implement me")
}

func (m *MockBlockChainContext) StateAt(root common.Hash) (*state.StateDB, error) {
	panic("implement me")
}

func (m *MockBlockChainContext) HasBadBlock(hash common.Hash) bool {
	panic("implement me")
}

func (m *MockBlockChainContext) Validator() core.Validator {
	panic("implement me")
}
