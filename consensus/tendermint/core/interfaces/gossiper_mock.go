// Code generated by MockGen. DO NOT EDIT.
// Source: consensus/tendermint/core/interfaces/gossiper.go
//
// Generated by this command:
//
//	mockgen -source=consensus/tendermint/core/interfaces/gossiper.go -package=interfaces -destination=consensus/tendermint/core/interfaces/gossiper_mock.go
//
// Package interfaces is a generated GoMock package.
package interfaces

import (
	reflect "reflect"

	common "github.com/autonity/autonity/common"
	consensus "github.com/autonity/autonity/consensus"
	message "github.com/autonity/autonity/consensus/tendermint/core/message"
	types "github.com/autonity/autonity/core/types"
	lru "github.com/hashicorp/golang-lru"
	gomock "go.uber.org/mock/gomock"
)

// MockGossiper is a mock of Gossiper interface.
type MockGossiper struct {
	ctrl     *gomock.Controller
	recorder *MockGossiperMockRecorder
}

// MockGossiperMockRecorder is the mock recorder for MockGossiper.
type MockGossiperMockRecorder struct {
	mock *MockGossiper
}

// NewMockGossiper creates a new mock instance.
func NewMockGossiper(ctrl *gomock.Controller) *MockGossiper {
	mock := &MockGossiper{ctrl: ctrl}
	mock.recorder = &MockGossiperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGossiper) EXPECT() *MockGossiperMockRecorder {
	return m.recorder
}

// Address mocks base method.
func (m *MockGossiper) Address() common.Address {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Address")
	ret0, _ := ret[0].(common.Address)
	return ret0
}

// Address indicates an expected call of Address.
func (mr *MockGossiperMockRecorder) Address() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Address", reflect.TypeOf((*MockGossiper)(nil).Address))
}

// AskSync mocks base method.
func (m *MockGossiper) AskSync(header *types.Header) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AskSync", header)
}

// AskSync indicates an expected call of AskSync.
func (mr *MockGossiperMockRecorder) AskSync(header any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AskSync", reflect.TypeOf((*MockGossiper)(nil).AskSync), header)
}

// Broadcaster mocks base method.
func (m *MockGossiper) Broadcaster() consensus.Broadcaster {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Broadcaster")
	ret0, _ := ret[0].(consensus.Broadcaster)
	return ret0
}

// Broadcaster indicates an expected call of Broadcaster.
func (mr *MockGossiperMockRecorder) Broadcaster() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Broadcaster", reflect.TypeOf((*MockGossiper)(nil).Broadcaster))
}

// Gossip mocks base method.
func (m *MockGossiper) Gossip(committee types.Committee, message message.Msg) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Gossip", committee, message)
}

// Gossip indicates an expected call of Gossip.
func (mr *MockGossiperMockRecorder) Gossip(committee, message any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Gossip", reflect.TypeOf((*MockGossiper)(nil).Gossip), committee, message)
}

// KnownMessages mocks base method.
func (m *MockGossiper) KnownMessages() *lru.ARCCache {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "KnownMessages")
	ret0, _ := ret[0].(*lru.ARCCache)
	return ret0
}

// KnownMessages indicates an expected call of KnownMessages.
func (mr *MockGossiperMockRecorder) KnownMessages() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "KnownMessages", reflect.TypeOf((*MockGossiper)(nil).KnownMessages))
}

// RecentMessages mocks base method.
func (m *MockGossiper) RecentMessages() *lru.ARCCache {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecentMessages")
	ret0, _ := ret[0].(*lru.ARCCache)
	return ret0
}

// RecentMessages indicates an expected call of RecentMessages.
func (mr *MockGossiperMockRecorder) RecentMessages() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecentMessages", reflect.TypeOf((*MockGossiper)(nil).RecentMessages))
}

// SetBroadcaster mocks base method.
func (m *MockGossiper) SetBroadcaster(broadcaster consensus.Broadcaster) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetBroadcaster", broadcaster)
}

// SetBroadcaster indicates an expected call of SetBroadcaster.
func (mr *MockGossiperMockRecorder) SetBroadcaster(broadcaster any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetBroadcaster", reflect.TypeOf((*MockGossiper)(nil).SetBroadcaster), broadcaster)
}

// UpdateStopChannel mocks base method.
func (m *MockGossiper) UpdateStopChannel(arg0 chan struct{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateStopChannel", arg0)
}

// UpdateStopChannel indicates an expected call of UpdateStopChannel.
func (mr *MockGossiperMockRecorder) UpdateStopChannel(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStopChannel", reflect.TypeOf((*MockGossiper)(nil).UpdateStopChannel), arg0)
}
