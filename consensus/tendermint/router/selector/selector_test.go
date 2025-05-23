package selector

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/fixsizecache"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/router"
	"github.com/autonity/autonity/consensus/tendermint/router/cache"
	"github.com/autonity/autonity/consensus/tendermint/router/network"
	"github.com/autonity/autonity/core/types"
)

const (
	seed = 12345
)

type MockCommittee struct {
	members []common.Address
}

func (mc *MockCommittee) MemberByAddress(addr common.Address) *common.Address {
	for _, member := range mc.members {
		if member == addr {
			return &member
		}
	}
	return nil
}

type MockPeer struct {
	cache *fixsizecache.Cache[common.Hash, bool]
	lat   uint
}

func NewMockPeer(lat uint) *MockPeer {
	return &MockPeer{
		cache: fixsizecache.New[common.Hash, bool](1000, 10, fixsizecache.HashKey[common.Hash]),
		lat:   lat,
	}
}

func (m *MockPeer) Send(msgcode uint64, data interface{}) error {
	time.Sleep(time.Duration(m.lat) * time.Millisecond)
	return nil
}

func (m *MockPeer) SendRaw(msgcode uint64, data []byte) error {
	return nil
}

func (m *MockPeer) Cache() *fixsizecache.Cache[common.Hash, bool] {
	return m.cache
}

// MockCache simulates caching
type MockCache struct {
	data sync.Map
}

func (mc *MockCache) Get(key string) (cache.Entry, bool) {
	if val, ok := mc.data.Load(key); ok {
		return val.(cache.Entry), true
	}
	return cache.Entry{}, false
}

func (mc *MockCache) Set(key string, recipients []common.Address) {
	mc.data.Store(key, CacheEntry{Recipients: recipients})
}

func (mc *MockCache) UpdateLastUsed(key string) {}
func (mc *MockCache) Invalidate()               {}
func (mc *MockCache) Cleanup()                  {}

// MockMessage simulates a message
type MockMessage struct {
	code       uint64
	originator common.Address
	hash       common.Hash
}

func (m *MockMessage) Code() uint64               { return m.code }
func (m *MockMessage) Originator() common.Address { return m.originator }
func (m *MockMessage) Hash() common.Hash          { return m.hash }

// generateClusters creates clusters using router.NewClusters
func generateClusters(nodeCount int, maxLatency uint) (network.Clusters, *types.Committee, common.Address) {
	latencyMap := make(map[common.Address]uint)
	addresses := make([]common.Address, nodeCount)
	committee := &types.Committee{Members: make([]types.CommitteeMember, nodeCount)}

	// Generate committee members and latency map
	for i := 0; i < nodeCount; i++ {
		addr := common.HexToAddress(fmt.Sprintf("0x%d", i))
		addresses[i] = addr
		committee.Members[i] = types.CommitteeMember{Address: addr}
		latencyMap[addr] = uint(rand.Intn(int(maxLatency)))
	}
	self := common.HexToAddress("0x0")

	// Create clusters using router.NewClusters
	clusters := network.New(addresses, latencyMap, self)

	return clusters, committee, self

}

// simulatePropagation simulates message propagation with latency
func simulatePropagation(selector *Selector, committee *types.Committee, msg message.Msg, from common.Address) time.Duration {
	start := time.Now()
	visited := sync.Map{}
	var wg sync.WaitGroup

	var propagate func(addr common.Address, depth int)
	propagate = func(addr common.Address, depth int) {
		defer wg.Done()
		if _, loaded := visited.LoadOrStore(addr, true); loaded {
			return
		}

		// Simulate latency
		peerWg := sync.WaitGroup{}
		peers, _ := selector.SelectPeers(committee, msg, addr)
		for _, peer := range peers {
			p, ok := selector.broadcaster.FindPeer(peer)
			if ok {
				peerWg.Add(1)
				go func() {
					p.Send(uint64(msg.Code()), new(interface{}))
					peerWg.Done()
				}()
			}
		}
		peerWg.Wait()

		// Propagate to next peers
		for _, peer := range peers {
			wg.Add(1)
			go propagate(peer, depth+1)
		}
	}

	wg.Add(1)
	propagate(from, 0)
	wg.Wait()
	return time.Since(start)
}

func TestClusterSelection(t *testing.T) {
	nodeCount := 100
	maxLatency := uint(300)

	clusters, committee, self := generateClusters(nodeCount, maxLatency)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	broadcaster := consensus.NewMockBroadcaster(ctrl)
	for _, cluster := range clusters.base {
		for _, node := range cluster.Members {
			mockedPeer := NewMockPeer(node.Lat)
			broadcaster.EXPECT().FindPeer(node.Addr).Return(mockedPeer, true).AnyTimes()
		}
	}
	router := &router.Router{
		self:        self,
		clusters:    clusters,
		broadcaster: broadcaster,
		cache:       &MockCache{},
	}
	selector := New(router)
	router.peerSelector = selector

	msg := &message.Fake{
		FakeCode:       0, // Proposal message
		FakeOriginator: selector.self,
		FakeHash:       common.HexToHash("0x1"),
	}

	peers, err := selector.SelectPeers(committee, msg, selector.self)
	if err != nil {
		t.Fatalf("Failed to select peers: %v", err)
	}

	clusterAssigned := make(map[int]bool)
	for _, peer := range peers {
		for _, cluster := range clusters.base {
			for _, node := range cluster.Members {
				if node.Addr == peer {
					clusterAssigned[node.ClusterID] = true
				}
			}
		}
	}

	// Check if one node per remote cluster is selected
	remoteClusters := len(clusters.base) - 1 // Exclude own cluster
	assignedClusters := len(clusterAssigned)
	if clusterAssigned[clusters.clusterContaining(selector.self)] {
		assignedClusters-- // Adjust for local cluster
	}
	if assignedClusters < remoteClusters {
		t.Errorf("Expected %d remote clusters assigned, got %d", remoteClusters, assignedClusters)
	}
}

func TestBucketSelection(t *testing.T) {
	rand.Seed(seed)
	committeeSize := 100
	clusters, committee, self := generateClusters(committeeSize, 300)
	numClusters := len(clusters.base)

	t.Run("verify cluster assignments", func(t *testing.T) {
		// Verify cluster assignments
		assert.Equal(t, numClusters, len(clusters.base), "Incorrect number of clusters")
		for clusterID := 0; clusterID < numClusters; clusterID++ {
			assert.NotEmpty(t, clusters.base[clusterID].Members, "Cluster %d should have members", clusterID)
			for _, node := range clusters.base[clusterID].Members {
				assert.Equal(t, clusterID, node.ClusterID, "Node %v has incorrect ClusterID", node.Addr)
				assert.True(t, node.Lat >= clusters.minLatency && node.Lat <= clusters.maxLatency, "Node %v latency out of range", node.Addr)
			}
		}
	})

	// Verify ownClusterID
	ownClusterID := clusters.clusterContaining(self)
	assert.NotEqual(t, -1, ownClusterID, "Self should be in a cluster")
	assert.Equal(t, ownClusterID, clusters.ownClusterID, "ownClusterID mismatch")

	// Verify bucket assignments (remote clusters only)
	assert.NotContains(t, clusters.clusterToBucket, ownClusterID, "Local cluster should not have bucket assignments")
	for clusterID, bucketIndices := range clusters.clusterToBucket {
		assert.NotEqual(t, ownClusterID, clusterID, "Remote cluster %d should not be local", clusterID)
		assert.NotEmpty(t, bucketIndices, "Cluster %d should have bucket assignments", clusterID)
		for _, bucketIdx := range bucketIndices {
			assert.True(t, bucketIdx >= 0 && bucketIdx < numClusters, "Invalid bucket index %d for cluster %d", bucketIdx, clusterID)
		}
	}

	// Verify remote primary nodes
	usedClusters := make(map[int]bool)
	for bucketIdx, node := range clusters.bucketNodes {
		assert.NotEqual(t, ownClusterID, node.ClusterID, "Primary node in bucket %d should not be local", bucketIdx)
		assert.False(t, usedClusters[node.ClusterID], "Cluster %d used multiple times in bucketNodes", node.ClusterID)
		usedClusters[node.ClusterID] = true
		bucketLat := uint(float64(bucketIdx)*clusters.bucketSize + float64(clusters.minLatency))
		assert.True(t, node.Lat >= bucketLat && node.Lat < bucketLat+uint(clusters.bucketSize), "Primary node %v in bucket %d has incorrect latency", node.Addr, bucketIdx)
	}
	assert.LessOrEqual(t, len(clusters.bucketNodes), numClusters-1, "Too many primary nodes assigned")

	// Verify remote fallback nodes
	for bucketIdx, fallbacks := range clusters.bucketFallbacks {
		for _, node := range fallbacks {
			assert.NotEqual(t, ownClusterID, node.ClusterID, "Fallback node in bucket %d should not be local", bucketIdx)
			if primary, exists := clusters.bucketNodes[bucketIdx]; exists {
				assert.Equal(t, primary.ClusterID, node.ClusterID, "Fallback node %v in bucket %d should be from same cluster as primary", node.Addr, bucketIdx)
				assert.NotEqual(t, primary.Addr, node.Addr, "Fallback node %v should not be primary", node.Addr)
			}
			targetLat := uint(float64(bucketIdx)*clusters.bucketSize + clusters.bucketSize/2 + float64(clusters.minLatency))
			diff := uint(math.Abs(float64(node.Lat - targetLat)))
			assert.LessOrEqual(t, diff, uint(clusters.bucketSize), "Fallback node %v in bucket %d too far from target latency %d", node.Addr, bucketIdx, targetLat)
		}
	}

	// Verify local nodes
	if ownClusterID != -1 {
		targetLocalNodes := int(math.Sqrt(float64(len(clusters.base[ownClusterID].Members))))
		assert.LessOrEqual(t, len(clusters.localBucketNodes), targetLocalNodes, "Too many local primary nodes")
		assert.LessOrEqual(t, len(clusters.localBucketFallbacks), targetLocalNodes, "Too many local fallback nodes")
		usedAddrs := make(map[common.Address]bool)
		for _, node := range clusters.localBucketNodes {
			assert.Equal(t, ownClusterID, node.ClusterID, "Local primary node %v should be in local cluster", node.Addr)
			assert.False(t, usedAddrs[node.Addr], "Duplicate local primary node %v", node.Addr)
			usedAddrs[node.Addr] = true
		}
		for _, node := range clusters.localBucketFallbacks {
			assert.Equal(t, ownClusterID, node.ClusterID, "Local fallback node %v should be in local cluster", node.Addr)
			assert.False(t, usedAddrs[node.Addr], "Duplicate local fallback node %v", node.Addr)
			usedAddrs[node.Addr] = true
		}
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	broadcaster := consensus.NewMockBroadcaster(ctrl)
	for _, cluster := range clusters.base {
		for _, node := range cluster.Members {
			mockedPeer := NewMockPeer(node.Lat)
			broadcaster.EXPECT().FindPeer(node.Addr).Return(mockedPeer, true).AnyTimes()
		}
	}
	router := &router.Router{
		self:        self,
		clusters:    clusters,
		broadcaster: broadcaster,
		cache:       &MockCache{},
	}
	selector := New(router)
	router.peerSelector = selector
	// Disconnect half of primary nodes
	for bucketIdx, node := range clusters.bucketNodes {
		if bucketIdx%2 == 0 {
			broadcaster.EXPECT().FindPeer(node.Addr).Return(nil, false).AnyTimes()
		}
	}
	msg := &message.Fake{
		FakeCode:       0, // Proposal message
		FakeOriginator: selector.self,
		FakeHash:       common.HexToHash("0x1"),
	}
	peers, err := selector.SelectPeers(committee, msg, self)
	assert.NoError(t, err, "Peer selection should succeed")
	assert.NotEmpty(t, peers, "Should select peers")

	// Verify one node per cluster
	minPeerLat, maxPeerLat := uint(math.MaxUint), uint(0)
	for _, peer := range peers {
		clusterID := clusters.clusterContaining(peer)
		assert.NotEqual(t, -1, clusterID, "Selected peer %v not in any cluster", peer)
		node, err := clusters.addressToMember(clusterID, peer)
		assert.NoError(t, err, "Peer %v not found in cluster %d", peer, clusterID)
		if node.Lat < minPeerLat {
			minPeerLat = node.Lat
		}
		if node.Lat > maxPeerLat {
			maxPeerLat = node.Lat
		}
	}
	// Verify all clusters are represented (except possibly empty ones)
	for clusterID := range clusters.base {
		if len(clusters.base[clusterID].Members) == 0 {
			continue
		}
		hasConnectedNode := false
		for _, node := range clusters.base[clusterID].Members {
			if _, ok := broadcaster.FindPeer(node.Addr); ok && committee.MemberByAddress(node.Addr) != nil {
				hasConnectedNode = true
				break
			}
		}
		if hasConnectedNode {
			assert.True(t, usedClusters[clusterID], "Cluster %d not represented despite having connected nodes", clusterID)
		}
	}
	latencySpread := maxPeerLat - minPeerLat
	expectedSpread := uint(float64(clusters.maxLatency-clusters.minLatency) * 0.5)
	assert.GreaterOrEqual(t, latencySpread, expectedSpread, "Latency spread %d too narrow, expected at least %d", latencySpread, expectedSpread)
}

func TestSmallCommittee(t *testing.T) {
	committeeSize := 4
	clusters, _, self := generateClusters(committeeSize, 300)

	numClusters := int(math.Floor(math.Sqrt(float64(committeeSize))))
	assert.Equal(t, numClusters, len(clusters.base), "Incorrect number of clusters")
	ownClusterID := clusters.clusterContaining(self)
	assert.NotEqual(t, -1, ownClusterID, "Self should be in a cluster")

	// Verify at most numClusters-1 remote nodes
	assert.LessOrEqual(t, len(clusters.bucketNodes), numClusters-1, "Too many primary nodes")
	for _, node := range clusters.bucketNodes {
		assert.NotEqual(t, ownClusterID, node.ClusterID, "Primary node %v should not be local", node.Addr)
	}
}

func TestNoLocalCluster(t *testing.T) {
	committeeSize := 10
	clusters, _, _ := generateClusters(committeeSize, 300)

	assert.Equal(t, -1, clusters.ownClusterID, "ownClusterID should be -1")
	assert.Empty(t, clusters.localBucketNodes, "No local primary nodes expected")
	assert.Empty(t, clusters.localBucketFallbacks, "No local fallback nodes expected")
	numClusters := int(math.Floor(math.Sqrt(float64(committeeSize))))
	assert.LessOrEqual(t, len(clusters.bucketNodes), numClusters, "Too many primary nodes")
}

func TestPeerSelectionStrategies(t *testing.T) {
	nodeCounts := []int{100}
	maxLatency := uint(300)

	for _, nodeCount := range nodeCounts {
		t.Run(fmt.Sprintf("Nodes=%d", nodeCount), func(t *testing.T) {
			// Setup selector
			clusters, committee, self := generateClusters(nodeCount, maxLatency)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			broadcaster := consensus.NewMockBroadcaster(ctrl)
			for _, cluster := range clusters.base {
				for _, node := range cluster.Members {
					mockedPeer := NewMockPeer(node.Lat)
					broadcaster.EXPECT().FindPeer(node.Addr).Return(mockedPeer, true).AnyTimes()
				}
			}
			router := &router.Router{
				self:        self,
				clusters:    clusters,
				broadcaster: broadcaster,
				cache:       &MockCache{},
			}
			selector := New(router)
			router.peerSelector = selector

			// Test for proposal and non-proposal messages
			for _, msgType := range []struct {
				name string
				code uint8
			}{
				{name: "Proposal", code: 0},
				{name: "NonProposal", code: 1},
			} {
				t.Run(msgType.name, func(t *testing.T) {
					msg := &message.Fake{
						FakeCode:       msgType.code,
						FakeOriginator: selector.self,
						FakeHash:       common.HexToHash(fmt.Sprintf("0x%d", rand.Int())),
					}
					// Run propagation
					duration := simulatePropagation(selector, committee, msg, selector.self)
					t.Logf("%s Propagation Time: %v", msgType.name, duration)
				})
			}
		})
	}
}
