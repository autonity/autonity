// Code generated by MockGen. DO NOT EDIT.
// Source: consensus/tendermint/core/core_backend.go

// Package core is a generated GoMock package.
package core

import (
	context "context"
	common "github.com/clearmatics/autonity/common"
	consensus "github.com/clearmatics/autonity/consensus"
	validator "github.com/clearmatics/autonity/consensus/tendermint/validator"
	state "github.com/clearmatics/autonity/core/state"
	types "github.com/clearmatics/autonity/core/types"
	event "github.com/clearmatics/autonity/event"
	p2p "github.com/clearmatics/autonity/p2p"
	rpc "github.com/clearmatics/autonity/rpc"
	gomock "github.com/golang/mock/gomock"
	big "math/big"
	reflect "reflect"
	time "time"
)

// MockBackend is a mock of Backend interface
type MockBackend struct {
	ctrl     *gomock.Controller
	recorder *MockBackendMockRecorder
}

// MockBackendMockRecorder is the mock recorder for MockBackend
type MockBackendMockRecorder struct {
	mock *MockBackend
}

// NewMockBackend creates a new mock instance
func NewMockBackend(ctrl *gomock.Controller) *MockBackend {
	mock := &MockBackend{ctrl: ctrl}
	mock.recorder = &MockBackendMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBackend) EXPECT() *MockBackendMockRecorder {
	return m.recorder
}

// Author mocks base method
func (m *MockBackend) Author(header *types.Header) (common.Address, error) {
	ret := m.ctrl.Call(m, "Author", header)
	ret0, _ := ret[0].(common.Address)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Author indicates an expected call of Author
func (mr *MockBackendMockRecorder) Author(header interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Author", reflect.TypeOf((*MockBackend)(nil).Author), header)
}

// VerifyHeader mocks base method
func (m *MockBackend) VerifyHeader(chain consensus.ChainReader, header *types.Header, seal bool) error {
	ret := m.ctrl.Call(m, "VerifyHeader", chain, header, seal)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifyHeader indicates an expected call of VerifyHeader
func (mr *MockBackendMockRecorder) VerifyHeader(chain, header, seal interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyHeader", reflect.TypeOf((*MockBackend)(nil).VerifyHeader), chain, header, seal)
}

// VerifyHeaders mocks base method
func (m *MockBackend) VerifyHeaders(chain consensus.ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	ret := m.ctrl.Call(m, "VerifyHeaders", chain, headers, seals)
	ret0, _ := ret[0].(chan<- struct{})
	ret1, _ := ret[1].(<-chan error)
	return ret0, ret1
}

// VerifyHeaders indicates an expected call of VerifyHeaders
func (mr *MockBackendMockRecorder) VerifyHeaders(chain, headers, seals interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyHeaders", reflect.TypeOf((*MockBackend)(nil).VerifyHeaders), chain, headers, seals)
}

// VerifyUncles mocks base method
func (m *MockBackend) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	ret := m.ctrl.Call(m, "VerifyUncles", chain, block)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifyUncles indicates an expected call of VerifyUncles
func (mr *MockBackendMockRecorder) VerifyUncles(chain, block interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyUncles", reflect.TypeOf((*MockBackend)(nil).VerifyUncles), chain, block)
}

// VerifySeal mocks base method
func (m *MockBackend) VerifySeal(chain consensus.ChainReader, header *types.Header) error {
	ret := m.ctrl.Call(m, "VerifySeal", chain, header)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifySeal indicates an expected call of VerifySeal
func (mr *MockBackendMockRecorder) VerifySeal(chain, header interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifySeal", reflect.TypeOf((*MockBackend)(nil).VerifySeal), chain, header)
}

// Prepare mocks base method
func (m *MockBackend) Prepare(chain consensus.ChainReader, header *types.Header) error {
	ret := m.ctrl.Call(m, "Prepare", chain, header)
	ret0, _ := ret[0].(error)
	return ret0
}

// Prepare indicates an expected call of Prepare
func (mr *MockBackendMockRecorder) Prepare(chain, header interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Prepare", reflect.TypeOf((*MockBackend)(nil).Prepare), chain, header)
}

// Finalize mocks base method
func (m *MockBackend) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header) {
	m.ctrl.Call(m, "Finalize", chain, header, state, txs, uncles)
}

// Finalize indicates an expected call of Finalize
func (mr *MockBackendMockRecorder) Finalize(chain, header, state, txs, uncles interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Finalize", reflect.TypeOf((*MockBackend)(nil).Finalize), chain, header, state, txs, uncles)
}

// FinalizeAndAssemble mocks base method
func (m *MockBackend) FinalizeAndAssemble(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {
	ret := m.ctrl.Call(m, "FinalizeAndAssemble", chain, header, state, txs, uncles, receipts)
	ret0, _ := ret[0].(*types.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FinalizeAndAssemble indicates an expected call of FinalizeAndAssemble
func (mr *MockBackendMockRecorder) FinalizeAndAssemble(chain, header, state, txs, uncles, receipts interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FinalizeAndAssemble", reflect.TypeOf((*MockBackend)(nil).FinalizeAndAssemble), chain, header, state, txs, uncles, receipts)
}

// Seal mocks base method
func (m *MockBackend) Seal(chain consensus.ChainReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {
	ret := m.ctrl.Call(m, "Seal", chain, block, results, stop)
	ret0, _ := ret[0].(error)
	return ret0
}

// Seal indicates an expected call of Seal
func (mr *MockBackendMockRecorder) Seal(chain, block, results, stop interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Seal", reflect.TypeOf((*MockBackend)(nil).Seal), chain, block, results, stop)
}

// SealHash mocks base method
func (m *MockBackend) SealHash(header *types.Header) common.Hash {
	ret := m.ctrl.Call(m, "SealHash", header)
	ret0, _ := ret[0].(common.Hash)
	return ret0
}

// SealHash indicates an expected call of SealHash
func (mr *MockBackendMockRecorder) SealHash(header interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SealHash", reflect.TypeOf((*MockBackend)(nil).SealHash), header)
}

// CalcDifficulty mocks base method
func (m *MockBackend) CalcDifficulty(chain consensus.ChainReader, time uint64, parent *types.Header) *big.Int {
	ret := m.ctrl.Call(m, "CalcDifficulty", chain, time, parent)
	ret0, _ := ret[0].(*big.Int)
	return ret0
}

// CalcDifficulty indicates an expected call of CalcDifficulty
func (mr *MockBackendMockRecorder) CalcDifficulty(chain, time, parent interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CalcDifficulty", reflect.TypeOf((*MockBackend)(nil).CalcDifficulty), chain, time, parent)
}

// APIs mocks base method
func (m *MockBackend) APIs(chain consensus.ChainReader) []rpc.API {
	ret := m.ctrl.Call(m, "APIs", chain)
	ret0, _ := ret[0].([]rpc.API)
	return ret0
}

// APIs indicates an expected call of APIs
func (mr *MockBackendMockRecorder) APIs(chain interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "APIs", reflect.TypeOf((*MockBackend)(nil).APIs), chain)
}

// Close mocks base method
func (m *MockBackend) Close() error {
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockBackendMockRecorder) Close() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockBackend)(nil).Close))
}

// NewChainHead mocks base method
func (m *MockBackend) NewChainHead() error {
	ret := m.ctrl.Call(m, "NewChainHead")
	ret0, _ := ret[0].(error)
	return ret0
}

// NewChainHead indicates an expected call of NewChainHead
func (mr *MockBackendMockRecorder) NewChainHead() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewChainHead", reflect.TypeOf((*MockBackend)(nil).NewChainHead))
}

// HandleMsg mocks base method
func (m *MockBackend) HandleMsg(address common.Address, data p2p.Msg) (bool, error) {
	ret := m.ctrl.Call(m, "HandleMsg", address, data)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HandleMsg indicates an expected call of HandleMsg
func (mr *MockBackendMockRecorder) HandleMsg(address, data interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleMsg", reflect.TypeOf((*MockBackend)(nil).HandleMsg), address, data)
}

// SetBroadcaster mocks base method
func (m *MockBackend) SetBroadcaster(arg0 consensus.Broadcaster) {
	m.ctrl.Call(m, "SetBroadcaster", arg0)
}

// SetBroadcaster indicates an expected call of SetBroadcaster
func (mr *MockBackendMockRecorder) SetBroadcaster(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetBroadcaster", reflect.TypeOf((*MockBackend)(nil).SetBroadcaster), arg0)
}

// Protocol mocks base method
func (m *MockBackend) Protocol() (string, uint64) {
	ret := m.ctrl.Call(m, "Protocol")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(uint64)
	return ret0, ret1
}

// Protocol indicates an expected call of Protocol
func (mr *MockBackendMockRecorder) Protocol() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Protocol", reflect.TypeOf((*MockBackend)(nil).Protocol))
}

// Start mocks base method
func (m *MockBackend) Start(ctx context.Context, chain consensus.ChainReader, currentBlock func() *types.Block, hasBadBlock func(common.Hash) bool) error {
	ret := m.ctrl.Call(m, "Start", ctx, chain, currentBlock, hasBadBlock)
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start
func (mr *MockBackendMockRecorder) Start(ctx, chain, currentBlock, hasBadBlock interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockBackend)(nil).Start), ctx, chain, currentBlock, hasBadBlock)
}

// Address mocks base method
func (m *MockBackend) Address() common.Address {
	ret := m.ctrl.Call(m, "Address")
	ret0, _ := ret[0].(common.Address)
	return ret0
}

// Address indicates an expected call of Address
func (mr *MockBackendMockRecorder) Address() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Address", reflect.TypeOf((*MockBackend)(nil).Address))
}

// Validators mocks base method
func (m *MockBackend) Validators(number uint64) validator.Set {
	ret := m.ctrl.Call(m, "Validators", number)
	ret0, _ := ret[0].(validator.Set)
	return ret0
}

// Validators indicates an expected call of Validators
func (mr *MockBackendMockRecorder) Validators(number interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validators", reflect.TypeOf((*MockBackend)(nil).Validators), number)
}

// Subscribe mocks base method
func (m *MockBackend) Subscribe(types ...interface{}) *event.TypeMuxSubscription {
	varargs := []interface{}{}
	for _, a := range types {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Subscribe", varargs...)
	ret0, _ := ret[0].(*event.TypeMuxSubscription)
	return ret0
}

// Subscribe indicates an expected call of Subscribe
func (mr *MockBackendMockRecorder) Subscribe(types ...interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockBackend)(nil).Subscribe), types...)
}

// Post mocks base method
func (m *MockBackend) Post(ev interface{}) {
	m.ctrl.Call(m, "Post", ev)
}

// Post indicates an expected call of Post
func (mr *MockBackendMockRecorder) Post(ev interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Post", reflect.TypeOf((*MockBackend)(nil).Post), ev)
}

// Broadcast mocks base method
func (m *MockBackend) Broadcast(ctx context.Context, valSet validator.Set, payload []byte) error {
	ret := m.ctrl.Call(m, "Broadcast", ctx, valSet, payload)
	ret0, _ := ret[0].(error)
	return ret0
}

// Broadcast indicates an expected call of Broadcast
func (mr *MockBackendMockRecorder) Broadcast(ctx, valSet, payload interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Broadcast", reflect.TypeOf((*MockBackend)(nil).Broadcast), ctx, valSet, payload)
}

// Gossip mocks base method
func (m *MockBackend) Gossip(ctx context.Context, valSet validator.Set, payload []byte) {
	m.ctrl.Call(m, "Gossip", ctx, valSet, payload)
}

// Gossip indicates an expected call of Gossip
func (mr *MockBackendMockRecorder) Gossip(ctx, valSet, payload interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Gossip", reflect.TypeOf((*MockBackend)(nil).Gossip), ctx, valSet, payload)
}

// Commit mocks base method
func (m *MockBackend) Commit(proposalBlock types.Block, seals [][]byte) error {
	ret := m.ctrl.Call(m, "Commit", proposalBlock, seals)
	ret0, _ := ret[0].(error)
	return ret0
}

// Commit indicates an expected call of Commit
func (mr *MockBackendMockRecorder) Commit(proposalBlock, seals interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Commit", reflect.TypeOf((*MockBackend)(nil).Commit), proposalBlock, seals)
}

// VerifyProposal mocks base method
func (m *MockBackend) VerifyProposal(arg0 types.Block) (time.Duration, error) {
	ret := m.ctrl.Call(m, "VerifyProposal", arg0)
	ret0, _ := ret[0].(time.Duration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifyProposal indicates an expected call of VerifyProposal
func (mr *MockBackendMockRecorder) VerifyProposal(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyProposal", reflect.TypeOf((*MockBackend)(nil).VerifyProposal), arg0)
}

// Sign mocks base method
func (m *MockBackend) Sign(arg0 []byte) ([]byte, error) {
	ret := m.ctrl.Call(m, "Sign", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Sign indicates an expected call of Sign
func (mr *MockBackendMockRecorder) Sign(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sign", reflect.TypeOf((*MockBackend)(nil).Sign), arg0)
}

// CheckSignature mocks base method
func (m *MockBackend) CheckSignature(data []byte, addr common.Address, sig []byte) error {
	ret := m.ctrl.Call(m, "CheckSignature", data, addr, sig)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckSignature indicates an expected call of CheckSignature
func (mr *MockBackendMockRecorder) CheckSignature(data, addr, sig interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckSignature", reflect.TypeOf((*MockBackend)(nil).CheckSignature), data, addr, sig)
}

// LastCommittedProposal mocks base method
func (m *MockBackend) LastCommittedProposal() (*types.Block, common.Address) {
	ret := m.ctrl.Call(m, "LastCommittedProposal")
	ret0, _ := ret[0].(*types.Block)
	ret1, _ := ret[1].(common.Address)
	return ret0, ret1
}

// LastCommittedProposal indicates an expected call of LastCommittedProposal
func (mr *MockBackendMockRecorder) LastCommittedProposal() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastCommittedProposal", reflect.TypeOf((*MockBackend)(nil).LastCommittedProposal))
}

// GetProposer mocks base method
func (m *MockBackend) GetProposer(number uint64) common.Address {
	ret := m.ctrl.Call(m, "GetProposer", number)
	ret0, _ := ret[0].(common.Address)
	return ret0
}

// GetProposer indicates an expected call of GetProposer
func (mr *MockBackendMockRecorder) GetProposer(number interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProposer", reflect.TypeOf((*MockBackend)(nil).GetProposer), number)
}

// HasBadProposal mocks base method
func (m *MockBackend) HasBadProposal(hash common.Hash) bool {
	ret := m.ctrl.Call(m, "HasBadProposal", hash)
	ret0, _ := ret[0].(bool)
	return ret0
}

// HasBadProposal indicates an expected call of HasBadProposal
func (mr *MockBackendMockRecorder) HasBadProposal(hash interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasBadProposal", reflect.TypeOf((*MockBackend)(nil).HasBadProposal), hash)
}

// SetProposedBlockHash mocks base method
func (m *MockBackend) SetProposedBlockHash(hash common.Hash) {
	m.ctrl.Call(m, "SetProposedBlockHash", hash)
}

// SetProposedBlockHash indicates an expected call of SetProposedBlockHash
func (mr *MockBackendMockRecorder) SetProposedBlockHash(hash interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetProposedBlockHash", reflect.TypeOf((*MockBackend)(nil).SetProposedBlockHash), hash)
}

// SyncPeer mocks base method
func (m *MockBackend) SyncPeer(address common.Address, messages []*Message) {
	m.ctrl.Call(m, "SyncPeer", address, messages)
}

// SyncPeer indicates an expected call of SyncPeer
func (mr *MockBackendMockRecorder) SyncPeer(address, messages interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SyncPeer", reflect.TypeOf((*MockBackend)(nil).SyncPeer), address, messages)
}

// ResetPeerCache mocks base method
func (m *MockBackend) ResetPeerCache(address common.Address) {
	m.ctrl.Call(m, "ResetPeerCache", address)
}

// ResetPeerCache indicates an expected call of ResetPeerCache
func (mr *MockBackendMockRecorder) ResetPeerCache(address interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetPeerCache", reflect.TypeOf((*MockBackend)(nil).ResetPeerCache), address)
}

// GetContractAddress mocks base method
func (m *MockBackend) GetContractAddress() common.Address {
	ret := m.ctrl.Call(m, "GetContractAddress")
	ret0, _ := ret[0].(common.Address)
	return ret0
}

// GetContractAddress indicates an expected call of GetContractAddress
func (mr *MockBackendMockRecorder) GetContractAddress() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContractAddress", reflect.TypeOf((*MockBackend)(nil).GetContractAddress))
}

// GetContractABI mocks base method
func (m *MockBackend) GetContractABI() string {
	ret := m.ctrl.Call(m, "GetContractABI")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetContractABI indicates an expected call of GetContractABI
func (mr *MockBackendMockRecorder) GetContractABI() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContractABI", reflect.TypeOf((*MockBackend)(nil).GetContractABI))
}

// WhiteList mocks base method
func (m *MockBackend) WhiteList() []string {
	ret := m.ctrl.Call(m, "WhiteList")
	ret0, _ := ret[0].([]string)
	return ret0
}

// WhiteList indicates an expected call of WhiteList
func (mr *MockBackendMockRecorder) WhiteList() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WhiteList", reflect.TypeOf((*MockBackend)(nil).WhiteList))
}

// AskSync mocks base method
func (m *MockBackend) AskSync(set validator.Set) {
	m.ctrl.Call(m, "AskSync", set)
}

// AskSync indicates an expected call of AskSync
func (mr *MockBackendMockRecorder) AskSync(set interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AskSync", reflect.TypeOf((*MockBackend)(nil).AskSync), set)
}
