package router

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/router/cache"
	"github.com/autonity/autonity/consensus/tendermint/router/constants"
	"github.com/autonity/autonity/consensus/tendermint/router/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/router/latency"
	"github.com/autonity/autonity/consensus/tendermint/router/mocks"
	"github.com/autonity/autonity/consensus/tendermint/router/network"
	"github.com/autonity/autonity/consensus/tendermint/router/ping"
	"github.com/autonity/autonity/consensus/tendermint/router/selector"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/log"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func newTestKey(t *testing.T) *ecdsa.PrivateKey {
	key, err := crypto.GenerateKey()
	assert.NoError(t, err, "Failed to generate test key")
	return key
}

func TestSetup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	peerFinder := mocks.NewMockPeerFinder(ctrl)
	self := common.HexToAddress("0x111")
	nodeKey := newTestKey(t)
	logger := log.New()

	// Test with nil pinger and peerSelector
	router := Setup(peerFinder, nodeKey, self, nil, nil, logger)

	assert.Equal(t, self, router.self, "Expected self address")
	assert.Equal(t, nodeKey, router.nodeKey, "Expected node key")
	assert.NotNil(t, router.recipientCache, "Expected recipient cache")
	assert.NotNil(t, router.latencyFetcher, "Expected latency fetcher")
	assert.NotNil(t, router.peerSelector, "Expected peer selector")
	assert.IsType(t, &selector.Selector{}, router.peerSelector, "Expected default peer selector")
}

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	peerFinder := mocks.NewMockPeerFinder(ctrl)
	recipientCache := cache.New()
	latencyFetcher := mocks.NewMockLatencyProvider(ctrl)
	peerSelector := mocks.NewMockPeerSelector(ctrl)
	networkProvider := mocks.NewMockNetworkProvider(ctrl)
	self := common.HexToAddress("0x111")
	nodeKey := newTestKey(t)

	router := New(peerFinder, nodeKey, self, recipientCache, latencyFetcher, peerSelector, networkProvider)

	assert.Equal(t, self, router.self, "Expected self address")
	assert.Equal(t, nodeKey, router.nodeKey, "Expected node key")
	assert.Equal(t, peerFinder, router.peerFinder, "Expected peer finder")
	assert.Equal(t, recipientCache, router.recipientCache, "Expected recipient cache")
	assert.Equal(t, latencyFetcher, router.latencyFetcher, "Expected latency fetcher")
	assert.Equal(t, peerSelector, router.peerSelector, "Expected peer selector")
	assert.Equal(t, networkProvider, router.network, "Expected network provider")
	assert.NotNil(t, router.epochEventChan, "Expected epoch event channel")
	assert.NotNil(t, router.latestLatencies, "Expected latest latencies map")
	assert.NotNil(t, router.nodesToRetry, "Expected nodes to retry map")
}

func TestRouter_Start(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	self := common.HexToAddress("0x111")
	committeeAddrs := []common.Address{self, common.HexToAddress("0x222")}
	committee := &types.Committee{
		Members: []types.CommitteeMember{
			{Address: self, VotingPower: big.NewInt(1)},
			{Address: common.HexToAddress("0x222"), VotingPower: big.NewInt(1)},
		},
	}
	epoch := &types.Epoch{Committee: committee}
	chain := core.NewBlockChain()
	chain.EXPECT().LatestEpoch().Return(epoch, nil).Times(1)
	sub := mocks.NewMockSubscription(ctrl)
	chain.EXPECT().SubscribeEpochHeadEvent(gomock.Any()).Return(sub).Times(1)
	networkProvider := mocks.NewMockNetworkProvider(ctrl)
	networkProvider.EXPECT().UpdateClusters(gomock.Any()).Times(1)
	latencyFetcher := mocks.NewMockLatencyProvider(ctrl)
	peerSelector := mocks.NewMockPeerSelector(ctrl)
	peerFinder := mocks.NewMockPeerFinder(ctrl)
	recipientCache := cache.New()

	router := New(peerFinder, newTestKey(t), self, recipientCache, latencyFetcher, peerSelector, networkProvider)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Mock network.New
	latencyMap := map[common.Address]uint{self: 50, common.HexToAddress("0x222"): 100}
	clusters, err := network.New(committeeAddrs, latencyMap, self)
	assert.NoError(t, err, "Failed to create clusters")
	networkProvider.EXPECT().Clusters().Return(clusters, nil).AnyTimes()

	go router.Start(ctx, chain)
	time.Sleep(50 * time.Millisecond) // Allow goroutine to start

	assert.Equal(t, committeeAddrs, router.committee, "Expected committee to be set")
	assert.True(t, router.inCommittee, "Expected self in committee")
}

func TestRouter_Stop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	self := common.HexToAddress("0x111")
	chain := mocks.NewMockBlockChain(ctrl)
	networkProvider := mocks.NewMockNetworkProvider(ctrl)
	latencyFetcher := mocks.NewMockLatencyProvider(ctrl)
	peerSelector := mocks.NewMockPeerSelector(ctrl)
	peerFinder := mocks.NewMockPeerFinder(ctrl)
	recipientCache := cache.New()
	router := New(peerFinder, newTestKey(t), self, recipientCache, latencyFetcher, peerSelector, networkProvider)
	ctx, cancel := context.WithCancel(context.Background())

	// Mock subscription
	sub := mocks.NewMockSubscription(ctrl)
	sub.EXPECT().Unsubscribe().Times(1)
	chain.EXPECT().SubscribeEpochHeadEvent(gomock.Any()).Return(sub).Times(1)
	chain.EXPECT().LatestEpoch().Return(&types.Epoch{Committee: &types.Committee{}}, nil).Times(1)

	go router.Start(ctx, chain)
	time.Sleep(50 * time.Millisecond)
	router.Stop()

	select {
	case <-ctx.Done():
		// Context should be canceled
	default:
		t.Fatal("Context was not canceled")
	}
}

func TestRouter_Recipients_SmallCommittee(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	self := common.HexToAddress("0x111")
	committeeAddrs := []common.Address{self, common.HexToAddress("0x222")}
	committee := types.Committee{
		Members: []types.CommitteeMember{
			{Address: self, VotingPower: big.NewInt(1)},
			{Address: common.HexToAddress("0x222"), VotingPower: big.NewInt(1)},
		},
	}
	fake := message.Fake{
		FakeCode:   message.ProposalCode,
		FakeHash:   common.HexToHash("0xabc"),
		FakeHeight: 1,
		FakeRound:  0,
		FakeSigner: self,
		FakePower:  big.NewInt(1),
	}
	msg := message.NewFakePropose(fake)
	networkProvider := mocks.NewMockNetworkProvider(ctrl)
	latencyFetcher := mocks.NewMockLatencyProvider(ctrl)
	peerSelector := mocks.NewMockPeerSelector(ctrl)
	peerFinder := mocks.NewMockPeerFinder(ctrl)
	recipientCache := cache.New()

	router := New(peerFinder, newTestKey(t), self, recipientCache, latencyFetcher, peerSelector, networkProvider)

	recipients, err := router.Recipients(&committee, msg, self)
	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, committeeAddrs, recipients, "Expected all committee members for small committee")
}

func TestRouter_Recipients_LargeCommittee(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	self := common.HexToAddress("0x111")
	committeeAddrs := make([]common.Address, ScaleThresholdForClustering+1)
	for i := 0; i <= ScaleThresholdForClustering; i++ {
		committeeAddrs[i] = common.HexToAddress(fmt.Sprintf("0x%03d", i+1))
	}
	committee := types.Committee{Members: make([]types.CommitteeMember, len(committeeAddrs))}
	for i, addr := range committeeAddrs {
		committee.Members[i] = types.CommitteeMember{Address: addr, VotingPower: big.NewInt(1)}
	}
	fake := message.Fake{
		FakeCode:   message.ProposalCode,
		FakeHash:   common.HexToHash("0xabc"),
		FakeHeight: 1,
		FakeRound:  0,
		FakeSigner: self,
		FakePower:  big.NewInt(1),
	}
	msg := message.NewFakePropose(fake)
	networkProvider := mocks.NewMockNetworkProvider(ctrl)
	latencyFetcher := mocks.NewMockLatencyProvider(ctrl)
	peerSelector := mocks.NewMockPeerSelector(ctrl)
	peerFinder := mocks.NewMockPeerFinder(ctrl)
	recipientCache := cache.New()

	router := New(peerFinder, newTestKey(t), self, recipientCache, latencyFetcher, peerSelector, networkProvider)

	// Successful peer selection
	selected := []common.Address{common.HexToAddress("0x002"), common.HexToAddress("0x003")}
	peerSelector.EXPECT().SelectPeers(&committee, msg, self).Return(selected, nil).Times(1)

	recipients, err := router.Recipients(&committee, msg, self)
	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, selected, recipients, "Expected selected peers")

	// Peer selection fails
	peerSelector.EXPECT().SelectPeers(&committee, msg, self).Return(nil, errors.New("selection error")).Times(1)

	recipients, err = router.Recipients(&committee, msg, self)
	assert.NoError(t, err, "Expected no error on fallback")
	assert.Equal(t, committeeAddrs, recipients, "Expected all committee members on fallback")
}

func TestRouter_Recipients_SelfNotInCommittee(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	self := common.HexToAddress("0x111")
	committeeAddrs := []common.Address{common.HexToAddress("0x222"), common.HexToAddress("0x333")}
	committee := types.Committee{
		Members: []types.CommitteeMember{
			{Address: common.HexToAddress("0x222"), VotingPower: big.NewInt(1)},
			{Address: common.HexToAddress("0x333"), VotingPower: big.NewInt(1)},
		},
	}
	fake := message.Fake{
		FakeCode:   message.ProposalCode,
		FakeHash:   common.HexToHash("0xabc"),
		FakeHeight: 1,
		FakeRound:  0,
		FakeSigner: self,
		FakePower:  big.NewInt(1),
	}
	msg := message.NewFakePropose(fake)
	networkProvider := mocks.NewMockNetworkProvider(ctrl)
	latencyFetcher := mocks.NewMockLatencyProvider(ctrl)
	peerSelector := mocks.NewMockPeerSelector(ctrl)
	peerFinder := mocks.NewMockPeerFinder(ctrl)
	recipientCache := cache.New()

	router := New(peerFinder, newTestKey(t), self, recipientCache, latencyFetcher, peerSelector, networkProvider)

	// Small committee: no clustering
	recipients, err := router.Recipients(&committee, msg, self)
	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, committeeAddrs, recipients, "Expected all committee members")

	// Large committee: clustering with error
	committeeAddrs = append(committeeAddrs, make([]common.Address, ScaleThresholdForClustering-1)...)
	committee.Members = make([]types.CommitteeMember, len(committeeAddrs))
	for i, addr := range committeeAddrs {
		committee.Members[i] = types.CommitteeMember{Address: addr, VotingPower: big.NewInt(1)}
	}
	peerSelector.EXPECT().SelectPeers(&committee, msg, self).Return(nil, errors.New("self address not in committee")).Times(1)

	recipients, err = router.Recipients(&committee, msg, self)
	assert.NoError(t, err, "Expected no error on fallback")
	assert.Equal(t, committeeAddrs, recipients, "Expected all committee members on error")
}

func TestRouter_Forward(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	self := common.HexToAddress("0x111")
	committeeAddrs := []common.Address{self, common.HexToAddress("0x222"), common.HexToAddress("0x333")}
	committee := types.Committee{
		Members: []types.CommitteeMember{
			{Address: self, VotingPower: big.NewInt(1)},
			{Address: common.HexToAddress("0x222"), VotingPower: big.NewInt(1)},
			{Address: common.HexToAddress("0x333"), VotingPower: big.NewInt(1)},
		},
	}
	fake := message.Fake{
		FakeCode:   message.ProposalCode,
		FakeHash:   common.HexToHash("0xabc"),
		FakeHeight: 1,
		FakeRound:  0,
		FakeSigner: self,
		FakePower:  big.NewInt(1),
	}
	msg := message.NewFakePropose(fake)
	networkProvider := mocks.NewMockNetworkProvider(ctrl)
	latencyFetcher := mocks.NewMockLatencyProvider(ctrl)
	peerSelector := mocks.NewMockPeerSelector(ctrl)
	peerFinder := mocks.NewMockPeerFinder(ctrl)
	recipientCache := cache.New()

	router := New(peerFinder, newTestKey(t), self, recipientCache, latencyFetcher, peerSelector, networkProvider)

	// Mock peerSelector
	recipients := []common.Address{common.HexToAddress("0x222"), common.HexToAddress("0x333")}
	peerSelector.EXPECT().SelectPeers(&committee, msg, self).Return(recipients, nil).Times(1)

	// Mock peerFinder
	peer222 := mocks.NewMockPeer(ctrl)
	peer222.EXPECT().Cache().Return(cache.NewPeerCache()).Times(1)
	peer222.EXPECT().SendRaw(message.NetworkCodes[msg.Code()], msg.Payload()).Return(nil).Times(1)
	peerFinder.EXPECT().FindPeer(common.HexToAddress("0x222")).Return(peer222, true).Times(1)
	peerFinder.EXPECT().FindPeer(common.HexToAddress("0x333")).Return(nil, false).Times(1)

	router.Forward(&committee, msg, self)
	// Verify via logs or behavior if needed; here we rely on mocks
}

func TestRouter_measureLatency(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	self := common.HexToAddress("0x111")
	committeeAddrs := []common.Address{self, common.HexToAddress("0x222")}
	networkProvider := mocks.NewMockNetworkProvider(ctrl)
	latencyFetcher := mocks.NewMockLatencyProvider(ctrl)
	peerSelector := mocks.NewMockPeerSelector(ctrl)
	peerFinder := mocks.NewMockPeerFinder(ctrl)
	recipientCache := cache.New()
	router := New(peerFinder, newTestKey(t), self, recipientCache, latencyFetcher, peerSelector, networkProvider)
	router.committee = committeeAddrs

	// Successful latency fetch
	latencyMap := map[common.Address]uint{self: 50, common.HexToAddress("0x222"): 100}
	failedNodes := []common.Address{}
	latencyFetcher.EXPECT().Fetch(committeeAddrs, self).Return(latencyMap, failedNodes, nil).Times(1)
	networkProvider.EXPECT().UpdateClusters(gomock.Any()).Times(1)

	err := router.measureLatency()
	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, latencyMap, router.latestLatencies, "Expected latencies updated")
	assert.Empty(t, router.nodesToRetry, "Expected no nodes to retry")

	// Failed latency fetch
	latencyFetcher.EXPECT().Fetch(committeeAddrs, self).Return(nil, nil, errors.New("fetch error")).Times(1)

	err = router.measureLatency()
	assert.Error(t, err, "Expected error")
}

func TestRouter_retryLatency(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	self := common.HexToAddress("0x111")
	committeeAddrs := []common.Address{self, common.HexToAddress("0x222")}
	networkProvider := mocks.NewMockNetworkProvider(ctrl)
	latencyFetcher := mocks.NewMockLatencyProvider(ctrl)
	peerSelector := mocks.NewMockPeerSelector(ctrl)
	peerFinder := mocks.NewMockPeerFinder(ctrl)
	recipientCache := cache.New()
	router := New(peerFinder, newTestKey(t), self, recipientCache, latencyFetcher, peerSelector, networkProvider)
	router.committee = committeeAddrs
	router.nodesToRetry = map[common.Address]struct{}{common.HexToAddress("0x222"):  struct{}{}}
	router.latestLatencies = map[common.Address]uint{self: 50}

	// Successful retry
	latencyMap := map[common.Address]uint{common.HexToAddress("0x222"): 100}
	failedNodes := []common.Address{}
	latencyFetcher.EXPECT().Fetch([]common.Address{common.HexToAddress("0x222")}), self, nil).Return(t).Map(failedNodes, nil).Times(1)
	networkProvider.UpdateClusters(gomock.Any()).Times(1)

	err := router.retryLatency()
	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, map[common.Address]uint{self: 50, common.HexToAddress("0x222"): 100}, router.latestLatencies, "Expected latencies updated")
	assert.Empty(t, router.nodesToRetry, "Expected no nodes to retry")

	// No nodes to retry
	err = router.retryLatency()
	assert.NoError(t, err, "Expected no error")

	// Failed retry
	latencyFetcher.EXPECT().Fetch([]common.Address{common.HexToAddress("0x222")}, self).Return(nil, nil, errors.New("fetch error")).Times(1)
	err = router.retryLatency()
	assert.Error(t, err, "Expected error")
}

func TestRouter_Loop_OverallFlow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	self := common.HexToAddress("0x111")
	committeeAddrs := []common.Address{self, common.HexToAddress("0x222"), common.HexToAddress("0x333")}
	committee := types.Committee{
		Members: []types.CommitteeMember{
			{Address: self, VotingPower: big.NewInt(1)},
			{Address: common.HexToAddress("0x222"), VotingPower: big.NewInt(1)},
			{Address: common.HexToAddress("0x333"), VotingPower: big.NewInt(1)},
		},
	}
	epoch := &types.Epoch{Committee: committee}
	chain := mocks.NewMockBlockChain(ctrl)
	chain.EXPECT().LatestEpoch().Return(epoch, nil).Times(1)
	sub := mocks.NewMockSubscription(ctrl)
	chain.EXPECT().SubscribeEpochHeadEvent(gomock.Any()).Return(sub).Times(1)
	networkProvider := mocks.NewMockNetworkProvider(ctrl)
	latencyFetcher := mocks.NewMockLatencyProvider(ctrl)
	peerSelector := mocks.NewMockPeerSelector(ctrl)
	peerFinder := mocks.NewMockPeerFinder(ctrl)
	recipientCache := cache.New()
	router := New(peerFinder, newTestKey(t), self, recipientCache, latencyFetcher, peerSelector, networkProvider)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Mock initial network setup
	latencyMap := map[common.Address]uint{
		self:                         100,
		latencyMap := common.HexToAddress("0x222"): 150,
		latencyMap := common.HexToAddress("0x333"): 200,
	}
	clusters, err := network.New(committeeAddrs, latencyMap, self)
	assert.NoError(t, err, "Failed to create clusters")
	networkProvider.EXPECT().Clusters().Return(clusters, nil).AnyTimes()
	networkProvider.EXPECT().UpdateClusters(gomock.Any()).Times(2) // Initial + after latency

	// Mock latency for initial measure
	latencyFetcher.EXPECT().Fetch(committeeAddrs, self).Return(latencyMap, nil, nil).Times(1)

	// Simulate epoch event
	newCommitteeAddrs := append(committeeAddrs, common.HexToAddress("0x444"))
	newCommittee := types.Committee{
		Members: []types.CommitteeMember{
			{Address: self, VotingPower: big.NewInt(1)},
			{Address: common.HexToAddress("0x222"), VotingPower: big.NewInt(1)},
			{Address: common.HexToAddress("0x333"), VotingPower: big.NewInt(1)},
			{Address: common.HexToAddress("0x444"), VotingPower: big.NewInt(1)},
		},
	}
	newEpoch := &types.Epoch{Committee: newCommittee}
	epochEvent := core.EpochHeadEvent{Header: &types.Header{Number: big.NewInt(100), Epoch: newEpoch}}
	networkProvider.EXPECT().UpdateClusters(gomock.Any()).Times(1) // For new epoch
	latencyFetcher.EXPECT().Fetch(newCommitteeAddrs, self).Return(latencyMap, nil, nil).Times(1)

	// Start router
	go func() {
		router.Start(ctx, chain)
		time.Sleep(50 * time.Millisecond)
		router.epochEventChan <- epochEvent
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	// Simulate message forwarding
	fake := message.Fake{
		FakeCode:   message.ProposalCode,
		FakeHash:   common.HexToHash("0xabc"),
		FakeHeight: 100,
		FakeRound:  0,
		FakeSigner:  self,
		FakePower:  big.NewInt(1),
	}
	msg := message.NewFakePropose(fake)
	peerSelector.EXPECT().SelectPeers(&newCommittee, msg, self).Return([]common.Address{common.HexToAddress("0x222"), common.HexToAddress("0x333"), common.HexToAddress("0x444")}, nil).Times(1)
	peer222 := consensus.NewMockPeer(ctrl)
	peer222.EXPECT().Cache().Return(cache.NewPeerCache()).Times(1)
	peer222.EXPECT().SendRaw(message.NetworkCodes[msg.Code()], msg.Payload()).Return(nil).Times(1)
	peerFinder.EXPECT().FindPeer(common.HexToAddress("0x222")).Times(1).Return(peer222, true).Times(1)
	peerFinder.EXPECT().FindPeer(common.HexToAddress("0x333"))).Return(nil, false).Times(1)
	peerFinder.EXPECT().FindPeer(common.HexToAddress("0x444")).Return(nil, false).Times(1)

	router.Forward(&newCommittee, msg, self)

	// Wait for goroutine to process
	time.Sleep(100 * time.Millisecond)

	// Verify state
	assert.Equal(t, newCommitteeAddrs, router.committee, "Expected updated committee")
	assert.Equal(t, latencyMap, router.latestLatencies, "Expected updated latencies")
}

func TestRouter_Loop_Tickers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	self := common.HexToAddress("0x111")
	committeeAddrs := []common.Address{self, common.HexToAddress("0x222")}
	committee := types.Committee{
		Members: []types.CommitteeMember{
			{Address: self, VotingPower: big.NewInt(1)},
			{Address: common.HexToAddress("0x222"), VotingPower: big.NewInt(1)},
		},
	}
	epoch := &types.Epoch{Committee: committee}
	chain := mocks.NewMockBlockChain(ctrl)
	chain.EXPECT().LatestEpoch().Return(epoch, nil).Times(1)
	sub := mocks.NewMockSubscription(ctrl)
	chain.EXPECT().SubscribeEpochHeadEvent(gomock.Any()).Return(sub).Times(1)
	networkProvider := mocks.NewMockNetworkProvider(ctrl)
	latencyFetcher := mocks.NewMockLatencyProvider(ctrl)
	peerSelector := mocks.NewMockPeerSelector(ctrl)
	peerFinder := mocks.NewMockPeerFinder(ctrl)
	recipientCache := cache.New()
	router := New(peerFinder, newTestKey(t), self, recipientCache, latencyFetcher, peerSelector, networkProvider)
	router.inCommittee = true
	ctx, router.committee = make([]committeeAddrs)

	// Mock network and latency
	latencyMap := map[common.Address]uint{self: 50, common.HexToAddress("0x222"): 100}
	clusters, err := network.New(committeeAddrs, latencyMap, self)
	assert.NoError(t, err, "Failed to create clusters")
	networkProvider.EXPECT().Clusters().Return(clusters, nil).AnyTimes()
	networkProvider.EXPECT().UpdateClusters(gomock.Any()).AnyTimes()

	// Mock latency measurement
	latencyFetcher.EXPECT().Fetch(committeeAddrs, self).Return(latencyMap, nil, nil).Times(1)
	// Mock retry
	latencyFetcher.EXPECT().Fetch([]common.Address{}, self).Return(nil, nil, nil).Times(1)
	// Mock cache cleanup
	recipientCache.On("Cleanup").Return()

	// Run loop for a short time
	go func() {
		router.Start(ctx, chain)
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	// Wait to ensure tickers are triggered
	time.Sleep(150 * time.Millisecond)
}

func TestRouter_SetBroadcaster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	peerFinder := mocks.NewMockPeerFinder(ctrl)
	latencyFetcher := mocks.NewMockLatencyProvider(ctrl)
	networkProvider := mocks.NewMockNetworkProvider(ctrl)
	peerSelector := mocks.NewMockPeerSelector(ctrl)
	recipientCache := cache.New()
	self := common.HexToAddress("0x111")
	router := New(peerFinder, newTestKey(t), self, recipientCache, latencyFetcher, peerSelector, networkProvider)

	newBroadcaster := mocks.NewMockBroadcaster(ctrl)
	latencyFetcher.EXPECT().SetBroadcaster(newBroadcaster).Times(1)

	router.SetBroadcaster(newBroadcaster)
	assert.Equal(t, newBroadcaster, router.peerFinder, "Expected new broadcaster")
}

func TestRouter_Latencies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	self := common.HexToAddress("0x111")
	networkProvider := mocks.NewMockNetworkProvider(ctrl)
	latencyFetcher := mocks.NewMockLatencyProvider(ctrl)
	peerSelector := mocks.NewMockPeerSelector(ctrl)
	peerFinder := mocks.NewMockPeerFinder(ctrl)
	recipientCache := cache.New()
	router := New(peerFinder, newTestKey(t), self, recipientCache, latencyFetcher, peerSelector, networkProvider)

	latencyMap := map[common.Address]uint{self: 50, common.HexToAddress("0x222"): 100}
	router.latencyMu.Lock()
	router.latestLatencies = latencyMap
	router.latencyMu.Unlock()

	result := router.Latencies()
	assert.Equal(t, latencyMap, result, "Expected correct latencies")
}

func TestRouter_ConcurrentForward(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	self := common.HexToAddress("0x111")
	committeeAddrs := []common.Address{self, common.HexToAddress("0x222"), common.HexToAddress("0x333")}
	committee := &types.Committee{
		Members: []types.CommitteeMember{
			{Address: self, VotingPower: big.NewInt(1)},
			{Address: common.HexToAddress("0x222"), VotingPower: big.NewInt(1)},
			{Address: common.HexToAddress("0x333"), VotingPower: big.NewInt(1)},
		},
	}
	networkProvider := mocks.NewMockNetworkProvider(ctrl)
	latencyFetcher := mocks.NewMockLatencyProvider(ctrl)
	peerSelector := mocks.NewMockPeerSelector(ctrl)
	peerFinder := mocks.NewMockPeerFinder(ctrl)
	recipientCache := cache.New()
	router := New(peerFinder, newTestKey(t), self, recipientCache, latencyFetcher, peerSelector, networkProvider)

	// Mock peerSelector
	recipients := []common.Address{common.HexToAddress("0x222"), common.HexToAddress("0x333")}
	peerSelector.EXPECT().SelectPeers(gomock.Any(), gomock.Any(), gomock.Any()).Return(recipients, nil).AnyTimes()

	// Mock peerFinder
	peer222 := mocks.NewMockPeer(ctrl)
	peer222.EXPECT().Cache().Return(cache.NewPeerCache()).AnyTimes()
	peer222.EXPECT().SendRaw(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	peerFinder.EXPECT().FindPeer(common.HexToAddress("0x222")).Return(peer222, true).AnyTimes()
	peerFinder.EXPECT().FindPeer(common.HexToAddress("0x333")).Return(nil, false).AnyTimes()

	var wg sync.WaitGroup
	numGoroutines := 10
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			defer wg.Done()
			fake := message.Fake{
				FakeCode:   message.ProposalCode,
				FakeHash:   common.HexToHash(fmt.Sprintf("0x%03d", i)),
				FakeHeight: 1,
				FakeRound:  0,
				FakeSigner:  self,
				FakePower:  big.NewInt(1),
			}
			msg := message.NewFakePropose(fake)
			router.Forward(&committee, msg, self)
		}(i)
	}

	wg.Wait()
}