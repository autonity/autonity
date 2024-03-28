// Code generated by MockGen. DO NOT EDIT.
// Source: consensus/tendermint/core/interfaces/core_backend.go
//
// Generated by this command:
//
//	mockgen -source=consensus/tendermint/core/interfaces/core_backend.go -package=interfaces -destination=consensus/tendermint/core/interfaces/core_backend_mock.go
//
// Package interfaces is a generated GoMock package.
package interfaces

import (
	context "context"
	big "math/big"
	reflect "reflect"
	time "time"

	abi "github.com/autonity/autonity/accounts/abi"
	autonity "github.com/autonity/autonity/autonity"
	common "github.com/autonity/autonity/common"
	message "github.com/autonity/autonity/consensus/tendermint/core/message"
	core "github.com/autonity/autonity/core"
	types "github.com/autonity/autonity/core/types"
	blst "github.com/autonity/autonity/crypto/blst"
	event "github.com/autonity/autonity/event"
	log "github.com/autonity/autonity/log"
	gomock "go.uber.org/mock/gomock"
)

// MockBackend is a mock of Backend interface.
type MockBackend struct {
	ctrl     *gomock.Controller
	recorder *MockBackendMockRecorder
}

// MockBackendMockRecorder is the mock recorder for MockBackend.
type MockBackendMockRecorder struct {
	mock *MockBackend
}

// NewMockBackend creates a new mock instance.
func NewMockBackend(ctrl *gomock.Controller) *MockBackend {
	mock := &MockBackend{ctrl: ctrl}
	mock.recorder = &MockBackendMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBackend) EXPECT() *MockBackendMockRecorder {
	return m.recorder
}

// AddSeal mocks base method.
func (m *MockBackend) AddSeal(block *types.Block) (*types.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddSeal", block)
	ret0, _ := ret[0].(*types.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddSeal indicates an expected call of AddSeal.
func (mr *MockBackendMockRecorder) AddSeal(block any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddSeal", reflect.TypeOf((*MockBackend)(nil).AddSeal), block)
}

// Address mocks base method.
func (m *MockBackend) Address() common.Address {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Address")
	ret0, _ := ret[0].(common.Address)
	return ret0
}

// Address indicates an expected call of Address.
func (mr *MockBackendMockRecorder) Address() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Address", reflect.TypeOf((*MockBackend)(nil).Address))
}

// AskSync mocks base method.
func (m *MockBackend) AskSync(header *types.Header) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AskSync", header)
}

// AskSync indicates an expected call of AskSync.
func (mr *MockBackendMockRecorder) AskSync(header any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AskSync", reflect.TypeOf((*MockBackend)(nil).AskSync), header)
}

// BlockChain mocks base method.
func (m *MockBackend) BlockChain() *core.BlockChain {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BlockChain")
	ret0, _ := ret[0].(*core.BlockChain)
	return ret0
}

// BlockChain indicates an expected call of BlockChain.
func (mr *MockBackendMockRecorder) BlockChain() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BlockChain", reflect.TypeOf((*MockBackend)(nil).BlockChain))
}

// Broadcast mocks base method.
func (m *MockBackend) Broadcast(committee types.Committee, message message.Msg) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Broadcast", committee, message)
}

// Broadcast indicates an expected call of Broadcast.
func (mr *MockBackendMockRecorder) Broadcast(committee, message any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Broadcast", reflect.TypeOf((*MockBackend)(nil).Broadcast), committee, message)
}

// Commit mocks base method.
func (m *MockBackend) Commit(proposalBlock *types.Block, round int64, seals types.AggregateSignature) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Commit", proposalBlock, round, seals)
	ret0, _ := ret[0].(error)
	return ret0
}

// Commit indicates an expected call of Commit.
func (mr *MockBackendMockRecorder) Commit(proposalBlock, round, seals any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Commit", reflect.TypeOf((*MockBackend)(nil).Commit), proposalBlock, round, seals)
}

// FutureMsgs mocks base method.
func (m *MockBackend) FutureMsgs() []message.Msg {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FutureMsgs")
	ret0, _ := ret[0].([]message.Msg)
	return ret0
}

// FutureMsgs indicates an expected call of FutureMsgs.
func (mr *MockBackendMockRecorder) FutureMsgs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FutureMsgs", reflect.TypeOf((*MockBackend)(nil).FutureMsgs))
}

// GetContractABI mocks base method.
func (m *MockBackend) GetContractABI() *abi.ABI {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContractABI")
	ret0, _ := ret[0].(*abi.ABI)
	return ret0
}

// GetContractABI indicates an expected call of GetContractABI.
func (mr *MockBackendMockRecorder) GetContractABI() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContractABI", reflect.TypeOf((*MockBackend)(nil).GetContractABI))
}

// Gossip mocks base method.
func (m *MockBackend) Gossip(committee types.Committee, message message.Msg) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Gossip", committee, message)
}

// Gossip indicates an expected call of Gossip.
func (mr *MockBackendMockRecorder) Gossip(committee, message any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Gossip", reflect.TypeOf((*MockBackend)(nil).Gossip), committee, message)
}

// Gossiper mocks base method.
func (m *MockBackend) Gossiper() Gossiper {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Gossiper")
	ret0, _ := ret[0].(Gossiper)
	return ret0
}

// Gossiper indicates an expected call of Gossiper.
func (mr *MockBackendMockRecorder) Gossiper() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Gossiper", reflect.TypeOf((*MockBackend)(nil).Gossiper))
}

// HandleUnhandledMsgs mocks base method.
func (m *MockBackend) HandleUnhandledMsgs(ctx context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "HandleUnhandledMsgs", ctx)
}

// HandleUnhandledMsgs indicates an expected call of HandleUnhandledMsgs.
func (mr *MockBackendMockRecorder) HandleUnhandledMsgs(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleUnhandledMsgs", reflect.TypeOf((*MockBackend)(nil).HandleUnhandledMsgs), ctx)
}

// HeadBlock mocks base method.
func (m *MockBackend) HeadBlock() *types.Block {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HeadBlock")
	ret0, _ := ret[0].(*types.Block)
	return ret0
}

// HeadBlock indicates an expected call of HeadBlock.
func (mr *MockBackendMockRecorder) HeadBlock() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HeadBlock", reflect.TypeOf((*MockBackend)(nil).HeadBlock))
}

// IsJailed mocks base method.
func (m *MockBackend) IsJailed(address common.Address) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsJailed", address)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsJailed indicates an expected call of IsJailed.
func (mr *MockBackendMockRecorder) IsJailed(address any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsJailed", reflect.TypeOf((*MockBackend)(nil).IsJailed), address)
}

// KnownMsgHash mocks base method.
func (m *MockBackend) KnownMsgHash() []common.Hash {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "KnownMsgHash")
	ret0, _ := ret[0].([]common.Hash)
	return ret0
}

// KnownMsgHash indicates an expected call of KnownMsgHash.
func (mr *MockBackendMockRecorder) KnownMsgHash() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "KnownMsgHash", reflect.TypeOf((*MockBackend)(nil).KnownMsgHash))
}

// Logger mocks base method.
func (m *MockBackend) Logger() log.Logger {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Logger")
	ret0, _ := ret[0].(log.Logger)
	return ret0
}

// Logger indicates an expected call of Logger.
func (mr *MockBackendMockRecorder) Logger() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Logger", reflect.TypeOf((*MockBackend)(nil).Logger))
}

// Post mocks base method.
func (m *MockBackend) Post(ev any) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Post", ev)
}

// Post indicates an expected call of Post.
func (mr *MockBackendMockRecorder) Post(ev any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Post", reflect.TypeOf((*MockBackend)(nil).Post), ev)
}

// ProcessFutureMsgs mocks base method.
func (m *MockBackend) ProcessFutureMsgs(height uint64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ProcessFutureMsgs", height)
}

// ProcessFutureMsgs indicates an expected call of ProcessFutureMsgs.
func (mr *MockBackendMockRecorder) ProcessFutureMsgs(height any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessFutureMsgs", reflect.TypeOf((*MockBackend)(nil).ProcessFutureMsgs), height)
}

// SetBlockchain mocks base method.
func (m *MockBackend) SetBlockchain(bc *core.BlockChain) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetBlockchain", bc)
}

// SetBlockchain indicates an expected call of SetBlockchain.
func (mr *MockBackendMockRecorder) SetBlockchain(bc any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetBlockchain", reflect.TypeOf((*MockBackend)(nil).SetBlockchain), bc)
}

// SetProposedBlockHash mocks base method.
func (m *MockBackend) SetProposedBlockHash(hash common.Hash) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetProposedBlockHash", hash)
}

// SetProposedBlockHash indicates an expected call of SetProposedBlockHash.
func (mr *MockBackendMockRecorder) SetProposedBlockHash(hash any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetProposedBlockHash", reflect.TypeOf((*MockBackend)(nil).SetProposedBlockHash), hash)
}

// Sign mocks base method.
func (m *MockBackend) Sign(hash common.Hash) (blst.Signature, common.Address) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sign", hash)
	ret0, _ := ret[0].(blst.Signature)
	ret1, _ := ret[1].(common.Address)
	return ret0, ret1
}

// Sign indicates an expected call of Sign.
func (mr *MockBackendMockRecorder) Sign(hash any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sign", reflect.TypeOf((*MockBackend)(nil).Sign), hash)
}

// Subscribe mocks base method.
func (m *MockBackend) Subscribe(types ...any) *event.TypeMuxSubscription {
	m.ctrl.T.Helper()
	varargs := []any{}
	for _, a := range types {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Subscribe", varargs...)
	ret0, _ := ret[0].(*event.TypeMuxSubscription)
	return ret0
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockBackendMockRecorder) Subscribe(types ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockBackend)(nil).Subscribe), types...)
}

// SyncPeer mocks base method.
func (m *MockBackend) SyncPeer(address common.Address) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SyncPeer", address)
}

// SyncPeer indicates an expected call of SyncPeer.
func (mr *MockBackendMockRecorder) SyncPeer(address any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SyncPeer", reflect.TypeOf((*MockBackend)(nil).SyncPeer), address)
}

// VerifyProposal mocks base method.
func (m *MockBackend) VerifyProposal(arg0 *types.Block) (time.Duration, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyProposal", arg0)
	ret0, _ := ret[0].(time.Duration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifyProposal indicates an expected call of VerifyProposal.
func (mr *MockBackendMockRecorder) VerifyProposal(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyProposal", reflect.TypeOf((*MockBackend)(nil).VerifyProposal), arg0)
}

// MockCore is a mock of Core interface.
type MockCore struct {
	ctrl     *gomock.Controller
	recorder *MockCoreMockRecorder
}

// MockCoreMockRecorder is the mock recorder for MockCore.
type MockCoreMockRecorder struct {
	mock *MockCore
}

// NewMockCore creates a new mock instance.
func NewMockCore(ctrl *gomock.Controller) *MockCore {
	mock := &MockCore{ctrl: ctrl}
	mock.recorder = &MockCoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCore) EXPECT() *MockCoreMockRecorder {
	return m.recorder
}

// Broadcaster mocks base method.
func (m *MockCore) Broadcaster() Broadcaster {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Broadcaster")
	ret0, _ := ret[0].(Broadcaster)
	return ret0
}

// Broadcaster indicates an expected call of Broadcaster.
func (mr *MockCoreMockRecorder) Broadcaster() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Broadcaster", reflect.TypeOf((*MockCore)(nil).Broadcaster))
}

// CoreState mocks base method.
func (m *MockCore) CoreState() CoreState {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CoreState")
	ret0, _ := ret[0].(CoreState)
	return ret0
}

// CoreState indicates an expected call of CoreState.
func (mr *MockCoreMockRecorder) CoreState() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CoreState", reflect.TypeOf((*MockCore)(nil).CoreState))
}

// CurrentHeightMessages mocks base method.
func (m *MockCore) CurrentHeightMessages() []message.Msg {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CurrentHeightMessages")
	ret0, _ := ret[0].([]message.Msg)
	return ret0
}

// CurrentHeightMessages indicates an expected call of CurrentHeightMessages.
func (mr *MockCoreMockRecorder) CurrentHeightMessages() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CurrentHeightMessages", reflect.TypeOf((*MockCore)(nil).CurrentHeightMessages))
}

// Height mocks base method.
func (m *MockCore) Height() *big.Int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Height")
	ret0, _ := ret[0].(*big.Int)
	return ret0
}

// Height indicates an expected call of Height.
func (mr *MockCoreMockRecorder) Height() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Height", reflect.TypeOf((*MockCore)(nil).Height))
}

// Power mocks base method.
func (m *MockCore) Power(h uint64, r int64) *big.Int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Power", h, r)
	ret0, _ := ret[0].(*big.Int)
	return ret0
}

// Power indicates an expected call of Power.
func (mr *MockCoreMockRecorder) Power(h, r any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Power", reflect.TypeOf((*MockCore)(nil).Power), h, r)
}

// Precommiter mocks base method.
func (m *MockCore) Precommiter() Precommiter {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Precommiter")
	ret0, _ := ret[0].(Precommiter)
	return ret0
}

// Precommiter indicates an expected call of Precommiter.
func (mr *MockCoreMockRecorder) Precommiter() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Precommiter", reflect.TypeOf((*MockCore)(nil).Precommiter))
}

// Prevoter mocks base method.
func (m *MockCore) Prevoter() Prevoter {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Prevoter")
	ret0, _ := ret[0].(Prevoter)
	return ret0
}

// Prevoter indicates an expected call of Prevoter.
func (mr *MockCoreMockRecorder) Prevoter() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Prevoter", reflect.TypeOf((*MockCore)(nil).Prevoter))
}

// Proposer mocks base method.
func (m *MockCore) Proposer() Proposer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Proposer")
	ret0, _ := ret[0].(Proposer)
	return ret0
}

// Proposer indicates an expected call of Proposer.
func (mr *MockCoreMockRecorder) Proposer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Proposer", reflect.TypeOf((*MockCore)(nil).Proposer))
}

// Round mocks base method.
func (m *MockCore) Round() int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Round")
	ret0, _ := ret[0].(int64)
	return ret0
}

// Round indicates an expected call of Round.
func (mr *MockCoreMockRecorder) Round() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Round", reflect.TypeOf((*MockCore)(nil).Round))
}

// Start mocks base method.
func (m *MockCore) Start(ctx context.Context, contract *autonity.ProtocolContracts) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Start", ctx, contract)
}

// Start indicates an expected call of Start.
func (mr *MockCoreMockRecorder) Start(ctx, contract any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockCore)(nil).Start), ctx, contract)
}

// Stop mocks base method.
func (m *MockCore) Stop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop")
}

// Stop indicates an expected call of Stop.
func (mr *MockCoreMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockCore)(nil).Stop))
}

// VotesPower mocks base method.
func (m *MockCore) VotesPower(h uint64, r int64, code uint8) *big.Int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VotesPower", h, r, code)
	ret0, _ := ret[0].(*big.Int)
	return ret0
}

// VotesPower indicates an expected call of VotesPower.
func (mr *MockCoreMockRecorder) VotesPower(h, r, code any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VotesPower", reflect.TypeOf((*MockCore)(nil).VotesPower), h, r, code)
}

// VotesPowerFor mocks base method.
func (m *MockCore) VotesPowerFor(h uint64, r int64, code uint8, v common.Hash) *big.Int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VotesPowerFor", h, r, code, v)
	ret0, _ := ret[0].(*big.Int)
	return ret0
}

// VotesPowerFor indicates an expected call of VotesPowerFor.
func (mr *MockCoreMockRecorder) VotesPowerFor(h, r, code, v any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VotesPowerFor", reflect.TypeOf((*MockCore)(nil).VotesPowerFor), h, r, code, v)
}
