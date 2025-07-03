package selector

import (
	"math/big"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/router/cache"
	"github.com/autonity/autonity/consensus/tendermint/router/cluster"
	"github.com/autonity/autonity/consensus/tendermint/router/mocks"
	"github.com/autonity/autonity/core/types"
)

func TestSelector_SelectPeers_Proposal_CacheHit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	np := mocks.NewMockClustersProvider(ctrl)
	recipients := mocks.NewMockRecipients(ctrl)
	peerFinder := mocks.NewMockPeerFinder(ctrl)
	selector := New(np, recipients)
	selector.SetBroadcaster(peerFinder)

	self := common.HexToAddress("0x111")
	committeeAddrs := []common.Address{
		self,
		common.HexToAddress("0x222"),
		common.HexToAddress("0x333"),
	}
	committee := types.Committee{
		Members: []types.CommitteeMember{
			{Address: self, VotingPower: big.NewInt(1)},
			{Address: common.HexToAddress("0x222"), VotingPower: big.NewInt(1)},
			{Address: common.HexToAddress("0x333"), VotingPower: big.NewInt(1)},
		},
	}
	latencyMap := map[common.Address]uint{
		self:                         50,
		common.HexToAddress("0x222"): 100,
		common.HexToAddress("0x333"): 150,
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
	from := self

	clusters, err := cluster.New(committeeAddrs, latencyMap, self)
	assert.NoError(t, err, "Failed to create clusters")
	np.EXPECT().Clusters().Return(clusters).Times(1)

	cacheKey := cache.GenerateKey(int(originator), true)
	cacheEntry := cache.Entry{Recipients: []common.Address{common.HexToAddress("0x222"), common.HexToAddress("0x333")}}
	recipients.EXPECT().Get(cacheKey).Return(cacheEntry, true).Times(1)
	peerFinder.EXPECT().FindPeer(common.HexToAddress("0x222")).Return(consensus.NewMockPeer(ctrl), true).Times(2) // allConnected + clusterStatus
	peerFinder.EXPECT().FindPeer(common.HexToAddress("0x333")).Return(consensus.NewMockPeer(ctrl), true).Times(2) // allConnected + clusterStatus
	recipients.EXPECT().UpdateLastUsed(cacheKey).Times(1)

	result, err := selector.SelectPeers(&committee, msg, from)
	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, []common.Address{common.HexToAddress("0x222"), common.HexToAddress("0x333")}, result, "Expected cached recipients")
}

func TestSelector_SelectPeers_SelfNotInCommittee(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	np := mocks.NewMockClustersProvider(ctrl)
	recipients := mocks.NewMockRecipients(ctrl)
	peerFinder := mocks.NewMockPeerFinder(ctrl)
	selector := New(np, recipients)
	selector.SetBroadcaster(peerFinder)

	self := common.HexToAddress("0x111")
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
	from := self

	np.EXPECT().Clusters().Return(cluster.Clusters{}).Times(1)

	_, err := selector.SelectPeers(&committee, msg, from)
	assert.Error(t, err, "Expected error for self not in committee")
	assert.Contains(t, err.Error(), "no clusters", "Expected correct error message")
}

func TestSelector_SelectPeers_EmptyCommittee(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	np := mocks.NewMockClustersProvider(ctrl)
	recipients := mocks.NewMockRecipients(ctrl)
	peerFinder := mocks.NewMockPeerFinder(ctrl)
	selector := New(np, recipients)
	selector.SetBroadcaster(peerFinder)

	self := common.HexToAddress("0x111")
	committee := types.Committee{Members: []types.CommitteeMember{}}
	fake := message.Fake{
		FakeCode:   message.ProposalCode,
		FakeHash:   common.HexToHash("0xabc"),
		FakeHeight: 1,
		FakeRound:  0,
		FakeSigner: self,
		FakePower:  big.NewInt(1),
	}
	msg := message.NewFakePropose(fake)
	from := self

	np.EXPECT().Clusters().Return(cluster.Clusters{}).Times(1)

	_, err := selector.SelectPeers(&committee, msg, from)
	assert.Error(t, err, "Expected error for empty committee")
	assert.Contains(t, err.Error(), "no clusters", "Expected correct error message")
}

func TestSelector_SelectPeers_NonProposal_NoCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	np := mocks.NewMockClustersProvider(ctrl)
	recipients := mocks.NewMockRecipients(ctrl)
	peerFinder := mocks.NewMockPeerFinder(ctrl)
	selector := New(np, recipients)
	selector.SetBroadcaster(peerFinder)

	self := common.HexToAddress("0x111")
	committeeAddrs := []common.Address{
		self,
		common.HexToAddress("0x222"),
		common.HexToAddress("0x333"),
	}
	committee := types.Committee{
		Members: []types.CommitteeMember{
			{Address: self, VotingPower: big.NewInt(1)},
			{Address: common.HexToAddress("0x222"), VotingPower: big.NewInt(1)},
			{Address: common.HexToAddress("0x333"), VotingPower: big.NewInt(1)},
		},
	}
	latencyMap := map[common.Address]uint{
		self:                         50,
		common.HexToAddress("0x222"): 100,
		common.HexToAddress("0x333"): 150,
	}
	fake := message.Fake{
		FakeCode:   message.PrevoteCode,
		FakeHash:   common.HexToHash("0xabc"),
		FakeHeight: 1,
		FakeRound:  0,
		FakeSigner: self,
		FakePower:  big.NewInt(1),
	}
	msg := message.NewFakePrevote(fake)
	from := self

	clusters, err := cluster.New(committeeAddrs, latencyMap, self)
	assert.NoError(t, err, "Failed to create clusters")
	np.EXPECT().Clusters().Return(clusters).AnyTimes()

	cacheKey := cache.GenerateKey(int(originator), false)
	recipients.EXPECT().Get(cacheKey).Return(cache.Entry{}, false).Times(1)
	peerFinder.EXPECT().FindPeer(common.HexToAddress("0x222")).Return(consensus.NewMockPeer(ctrl), true).AnyTimes()
	peerFinder.EXPECT().FindPeer(common.HexToAddress("0x333")).Return(consensus.NewMockPeer(ctrl), true).AnyTimes()
	recipients.EXPECT().Set(cacheKey, gomock.Any()).Times(1)

	result, err := selector.SelectPeers(&committee, msg, from)
	assert.NoError(t, err, "Expected no error")
	assert.Contains(t, result, common.HexToAddress("0x222"), "Expected node from remote cluster")
}

func TestSelector_selectNodesByLatencySpread(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	np := mocks.NewMockClustersProvider(ctrl)
	recipients := mocks.NewMockRecipients(ctrl)
	peerFinder := mocks.NewMockPeerFinder(ctrl)
	selector := &selector{clustersProvider: np, recipientCache: recipients}
	selector.SetBroadcaster(peerFinder)

	self := common.HexToAddress("0x111")
	committeeAddrs := []common.Address{
		self,
		common.HexToAddress("0x222"),
		common.HexToAddress("0x333"),
		common.HexToAddress("0x444"),
	}
	latencyMap := map[common.Address]uint{
		self:                         50,
		common.HexToAddress("0x222"): 100,
		common.HexToAddress("0x333"): 150,
		common.HexToAddress("0x444"): 200,
	}
	clusters, err := cluster.New(committeeAddrs, latencyMap, self)
	assert.NoError(t, err, "Failed to create clusters")
	np.EXPECT().Clusters().Return(clusters).Times(1)
	peerFinder.EXPECT().FindPeer(common.HexToAddress("0x222")).Return(consensus.NewMockPeer(ctrl), true).Times(1)

	nodes := selector.selectNodesByLatencySpread()
	assert.Len(t, nodes, 1, "Expected one node per remote cluster")
	assert.Contains(t, nodes, cluster.Node{Addr: common.HexToAddress("0x222"), Lat: 100, ClusterID: 1}, "Expected node 0x222")
}

func TestSelector_selectBucketBasedNodes_Originator(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	np := mocks.NewMockClustersProvider(ctrl)
	recipients := mocks.NewMockRecipients(ctrl)
	peerFinder := mocks.NewMockPeerFinder(ctrl)
	selector := &selector{clustersProvider: np, recipientCache: recipients}
	selector.SetBroadcaster(peerFinder)

	self := common.HexToAddress("0x111")
	committeeAddrs := []common.Address{
		self,
		common.HexToAddress("0x222"),
		common.HexToAddress("0x333"),
		common.HexToAddress("0x444"),
	}
	committee := types.Committee{
		Members: []types.CommitteeMember{
			{Address: self, VotingPower: big.NewInt(1)},
			{Address: common.HexToAddress("0x222"), VotingPower: big.NewInt(1)},
			{Address: common.HexToAddress("0x333"), VotingPower: big.NewInt(1)},
			{Address: common.HexToAddress("0x444"), VotingPower: big.NewInt(1)},
		},
	}
	latencyMap := map[common.Address]uint{
		self:                         50,
		common.HexToAddress("0x222"): 100,
		common.HexToAddress("0x333"): 150,
		common.HexToAddress("0x444"): 200,
	}
	clusters, err := cluster.New(committeeAddrs, latencyMap, self)
	assert.NoError(t, err, "Failed to create clusters")
	np.EXPECT().Clusters().Return(clusters).AnyTimes()
	peerFinder.EXPECT().FindPeer(common.HexToAddress("0x222")).Return(consensus.NewMockPeer(ctrl), true).AnyTimes()
	peerFinder.EXPECT().FindPeer(common.HexToAddress("0x333")).Return(consensus.NewMockPeer(ctrl), true).AnyTimes()
	peerFinder.EXPECT().FindPeer(common.HexToAddress("0x444")).Return(consensus.NewMockPeer(ctrl), true).AnyTimes()

	nodes := selector.selectNodesForProposal( clusters, &committee, originator, 0)
	assert.Contains(t, nodes, cluster.Node{Addr: common.HexToAddress("0x222"), Lat: 100, ClusterID: 1}, "Expected node from cluster 1")
	assert.Contains(t, nodes, cluster.Node{Addr: common.HexToAddress("0x333"), Lat: 150, ClusterID: 0}, "Expected node from cluster 0")
}

func TestSelector_selectBucketBasedNodes_FirstRelayerOriginCluster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	np := mocks.NewMockClustersProvider(ctrl)
	recipients := mocks.NewMockRecipients(ctrl)
	peerFinder := mocks.NewMockPeerFinder(ctrl)
	selector := &selector{clustersProvider: np, recipientCache: recipients}
	selector.SetBroadcaster(peerFinder)

	self := common.HexToAddress("0x111")
	committeeAddrs := []common.Address{
		self,
		common.HexToAddress("0x222"),
		common.HexToAddress("0x333"),
		common.HexToAddress("0x444"),
	}
	committee := types.Committee{
		Members: []types.CommitteeMember{
			{Address: self, VotingPower: big.NewInt(1)},
			{Address: common.HexToAddress("0x222"), VotingPower: big.NewInt(1)},
			{Address: common.HexToAddress("0x333"), VotingPower: big.NewInt(1)},
			{Address: common.HexToAddress("0x444"), VotingPower: big.NewInt(1)},
		},
	}
	latencyMap := map[common.Address]uint{
		self:                         50,
		common.HexToAddress("0x222"): 100,
		common.HexToAddress("0x333"): 150,
		common.HexToAddress("0x444"): 200,
	}
	clusters, err := cluster.New(committeeAddrs, latencyMap, self)
	assert.NoError(t, err, "Failed to create clusters")
	np.EXPECT().Clusters().Return(clusters).AnyTimes()
	peerFinder.EXPECT().FindPeer(common.HexToAddress("0x222")).Return(consensus.NewMockPeer(ctrl), true).AnyTimes()
	peerFinder.EXPECT().FindPeer(common.HexToAddress("0x333")).Return(consensus.NewMockPeer(ctrl), true).AnyTimes()
	peerFinder.EXPECT().FindPeer(common.HexToAddress("0x444")).Return(consensus.NewMockPeer(ctrl), true).AnyTimes()

	nodes := selector.selectNodesForProposal(clusters, &committee, firstRelayerOriginCluster, 0)
	assert.Contains(t, nodes, cluster.Node{Addr: common.HexToAddress("0x222"), Lat: 100, ClusterID: 1}, "Expected node from remote cluster")
}

func TestSelector_deduplicate(t *testing.T) {
	selector := &selector{}
	nodes := []cluster.Node{
		{Addr: common.HexToAddress("0x111"), Lat: 50, ClusterID: 0},
		{Addr: common.HexToAddress("0x111"), Lat: 50, ClusterID: 0},
		{Addr: common.HexToAddress("0x222"), Lat: 100, ClusterID: 1},
	}
	result := selector.deduplicate(nodes)
	assert.Len(t, result, 2, "Expected duplicates removed")
	assert.Contains(t, result, cluster.Node{Addr: common.HexToAddress("0x111"), Lat: 50, ClusterID: 0}, "Expected node 0x111")
	assert.Contains(t, result, cluster.Node{Addr: common.HexToAddress("0x222"), Lat: 100, ClusterID: 1}, "Expected node 0x222")
}

func TestSelector_selectCloseNodes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	np := mocks.NewMockClustersProvider(ctrl)
	recipients := mocks.NewMockRecipients(ctrl)
	peerFinder := mocks.NewMockPeerFinder(ctrl)
	selector := &selector{clustersProvider: np, recipientCache: recipients}
	selector.SetBroadcaster(peerFinder)

	self := common.HexToAddress("0x111")
	committeeAddrs := []common.Address{
		self,
		common.HexToAddress("0x222"),
		common.HexToAddress("0x333"),
	}
	committee := types.Committee{
		Members: []types.CommitteeMember{
			{Address: self, VotingPower: big.NewInt(1)},
			{Address: common.HexToAddress("0x222"), VotingPower: big.NewInt(1)},
			{Address: common.HexToAddress("0x333"), VotingPower: big.NewInt(1)},
		},
	}
	latencyMap := map[common.Address]uint{
		self:                         0,
		common.HexToAddress("0x222"): 10,
		common.HexToAddress("0x333"): 15,
	}
	clusters, err := cluster.New(committeeAddrs, latencyMap, self)
	assert.NoError(t, err, "Failed to create clusters")
	np.EXPECT().Clusters().Return(clusters).Times(1)
	peerFinder.EXPECT().FindPeer(common.HexToAddress("0x222")).Return(consensus.NewMockPeer(ctrl), true).Times(1)
	peerFinder.EXPECT().FindPeer(common.HexToAddress("0x333")).Return(consensus.NewMockPeer(ctrl), true).Times(1)

	nodes := selector.selectCloseNodes(&committee, 0, 1, 2, self)
	assert.Len(t, nodes, 2, "Expected 2 nodes (1 min, 1 low-latency)")
	assert.Equal(t, common.HexToAddress("0x222"), nodes[0].Addr, "Expected lowest latency node first")
}

func TestSelector_routingCandidatesFromCluster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	np := mocks.NewMockClustersProvider(ctrl)
	recipients := mocks.NewMockRecipients(ctrl)
	peerFinder := mocks.NewMockPeerFinder(ctrl)
	selector := &selector{clustersProvider: np, recipientCache: recipients}
	selector.SetBroadcaster(peerFinder)

	self := common.HexToAddress("0x111")
	committeeAddrs := []common.Address{
		self,
		common.HexToAddress("0x222"),
	}
	committee := types.Committee{
		Members: []types.CommitteeMember{
			{Address: self, VotingPower: big.NewInt(1)},
			{Address: common.HexToAddress("0x222"), VotingPower: big.NewInt(1)},
		},
	}
	latencyMap := map[common.Address]uint{
		self:                         50,
		common.HexToAddress("0x222"): 100,
	}
	clusters, err := cluster.New(committeeAddrs, latencyMap, self)
	assert.NoError(t, err, "Failed to create clusters")
	np.EXPECT().Clusters().Return(clusters).Times(1)
	peerFinder.EXPECT().FindPeer(common.HexToAddress("0x222")).Return(consensus.NewMockPeer(ctrl), true).Times(1)

	candidates := selector.routingCandidatesFromCluster(clusters.ID(), self, &committee)
	assert.Len(t, candidates, 1, "Expected one candidate")
	assert.Equal(t, common.HexToAddress("0x222"), candidates[0].Addr, "Expected node 0x222")
}

func TestSelector_allConnected(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	np := mocks.NewMockClustersProvider(ctrl)
	recipients := mocks.NewMockRecipients(ctrl)
	peerFinder := mocks.NewMockPeerFinder(ctrl)
	selector := &selector{clustersProvider: np, recipientCache: recipients}
	selector.SetBroadcaster(peerFinder)

	recipientsToTest := []common.Address{common.HexToAddress("0x111"), common.HexToAddress("0x222")}
	peerFinder.EXPECT().FindPeer(common.HexToAddress("0x111")).Return(consensus.NewMockPeer(ctrl), true).Times(1)
	peerFinder.EXPECT().FindPeer(common.HexToAddress("0x222")).Return(nil, false).Times(1)

	result := selector.allConnected(recipientsToTest)
	assert.False(t, result, "Expected false due to disconnected peer")
}

func TestSelector_determineSenderType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	np := mocks.NewMockClustersProvider(ctrl)

	self := common.HexToAddress("0x111")
	committeeAddrs := []common.Address{
		self,
		common.HexToAddress("0x222"),
		common.HexToAddress("0x333"),
	}
	latencyMap := map[common.Address]uint{
		self:                         50,
		common.HexToAddress("0x222"): 100,
		common.HexToAddress("0x333"): 150,
	}
	clusters, err := cluster.New(committeeAddrs, latencyMap, self)
	assert.NoError(t, err, "Failed to create clusters")

	tests := []struct {
		name            string
		from            common.Address
		signer          common.Address
		ownClusterID    int
		senderClusterID int
		originClusterID int
		expected        SenderType
	}{
		{
			name:            "originator",
			from:            self,
			signer:          self,
			ownClusterID:    0,
			senderClusterID: 0,
			originClusterID: 0,
			expected:        originator,
		},
		{
			name:            "firstRelayerOriginCluster",
			from:            common.HexToAddress("0x222"),
			signer:          common.HexToAddress("0x222"),
			ownClusterID:    0,
			senderClusterID: 1,
			originClusterID: 0,
			expected:        firstRelayerOriginCluster,
		},
		{
			name:            "localRelayerRemoteCluster",
			from:            common.HexToAddress("0x222"),
			signer:          common.HexToAddress("0x333"),
			ownClusterID:    0,
			senderClusterID: 2,
			originClusterID: 1,
			expected:        localRelayerRemoteCluster,
		},
		{
			name:            "firstRelayerRemoteCluster",
			from:            common.HexToAddress("0x222"),
			signer:          common.HexToAddress("0x222"),
			ownClusterID:    0,
			senderClusterID: 1,
			originClusterID: 1,
			expected:        firstRelayerRemoteCluster,
		},
		{
			name:            "localRelayerRemoteCluster",
			from:            common.HexToAddress("0x222"),
			signer:          common.HexToAddress("0x333"),
			ownClusterID:    0,
			senderClusterID: 1,
			originClusterID: 2,
			expected:        localRelayerRemoteCluster,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := message.Fake{
				FakeCode:   message.PrevoteCode,
				FakeHash:   common.HexToHash("0xabc"),
				FakeHeight: 1,
				FakeRound:  0,
				FakeSigner: tt.signer,
				FakePower:  big.NewInt(1),
			}
			msg := message.NewFakePrevote(fake)
			np.EXPECT().Clusters().Return(clusters).AnyTimes()
			result := determineSenderType(tt.from, clusters.Self(), msg, tt.originClusterID, tt.ownClusterID, tt.senderClusterID)
			assert.Equal(t, tt.expected, result, "Expected correct sender type")
		})
	}
}

func TestSelector_ConcurrentSelectPeers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	np := mocks.NewMockClustersProvider(ctrl)
	recipients := mocks.NewMockRecipients(ctrl)
	peerFinder := mocks.NewMockPeerFinder(ctrl)
	selector := New(np, recipients)
	selector.SetBroadcaster(peerFinder)

	self := common.HexToAddress("0x111")
	committeeAddrs := []common.Address{
		self,
		common.HexToAddress("0x222"),
	}
	committee := types.Committee{
		Members: []types.CommitteeMember{
			{Address: self, VotingPower: big.NewInt(1)},
			{Address: common.HexToAddress("0x222"), VotingPower: big.NewInt(1)},
		},
	}
	latencyMap := map[common.Address]uint{
		self:                         50,
		common.HexToAddress("0x222"): 100,
	}
	clusters, err := cluster.New(committeeAddrs, latencyMap, self)
	assert.NoError(t, err, "Failed to create clusters")
	np.EXPECT().Clusters().Return(clusters).AnyTimes()
	recipients.EXPECT().Get(gomock.Any()).Return(cache.Entry{}, false).AnyTimes()
	peerFinder.EXPECT().FindPeer(common.HexToAddress("0x222")).Return(consensus.NewMockPeer(ctrl), true).AnyTimes()
	recipients.EXPECT().Set(gomock.Any(), gomock.Any()).AnyTimes()

	var wg sync.WaitGroup
	numGoroutines := 10
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			fake := message.Fake{
				FakeCode:   message.ProposalCode,
				FakeHash:   common.HexToHash("0xabc"),
				FakeHeight: 1,
				FakeRound:  0,
				FakeSigner: self,
				FakePower:  big.NewInt(1),
			}
			msg := message.NewFakePropose(fake)
			_, err := selector.SelectPeers(&committee, msg, self)
			assert.NoError(t, err, "Expected no error in concurrent select")
		}()
	}

	wg.Wait()
}
