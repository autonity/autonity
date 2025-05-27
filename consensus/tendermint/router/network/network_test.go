package network

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/router/constants"
)

func TestNew(t *testing.T) {
	committee := []common.Address{
		common.HexToAddress("0x111"), // Cluster 0
		common.HexToAddress("0x222"), // Cluster 1
		common.HexToAddress("0x333"), // Cluster 0
		common.HexToAddress("0x444"), // Cluster 1
	}
	latencyMap := map[common.Address]uint{
		common.HexToAddress("0x111"): 50,
		common.HexToAddress("0x222"): 100,
		common.HexToAddress("0x333"): 150,
		common.HexToAddress("0x444"): 200,
	}
	self := common.HexToAddress("0x111") // In Cluster 0
	numClusters := 2                     // sqrt(4) = 2

	clusters := New(committee, latencyMap, self)

	// Verify cluster structure
	assert.Equal(t, numClusters, len(clusters.Base()), "Expected 2 clusters")
	assert.Equal(t, 0, clusters.ID(), "Expected self in cluster 0")
	assert.Equal(t, self, clusters.Self(), "Expected self address")
	assert.Equal(t, 50, int(clusters.minLatency), "Expected min latency 50")
	assert.Equal(t, uint(200*constants.MaxLatencyCapFactor), clusters.maxLatency, "Expected max latency capped")
	// Cluster 0: 0x111 (removed by Prepare), 0x333
	// Cluster 1: 0x222, 0x444
	// Nodes are sorted by latency after Prepare
	assert.Equal(t, []Node{{Addr: common.HexToAddress("0x333"), Lat: 150, ClusterID: 0}}, clusters.base[0], "Expected cluster 0 nodes (self removed)")
	assert.Equal(t, []Node{
		{Addr: common.HexToAddress("0x222"), Lat: 100, ClusterID: 1},
		{Addr: common.HexToAddress("0x444"), Lat: 200, ClusterID: 1},
	}, clusters.base[1], "Expected cluster 1 nodes")
}

func TestClusters_ComputeLatencyBuckets(t *testing.T) {
	committee := []common.Address{
		common.HexToAddress("0x111"), // Cluster 0
		common.HexToAddress("0x222"), // Cluster 1
		common.HexToAddress("0x333"), // Cluster 0
	}
	latencyMap := map[common.Address]uint{
		common.HexToAddress("0x111"): 50,
		common.HexToAddress("0x222"): 100,
		common.HexToAddress("0x333"): 150,
	}
	self := common.HexToAddress("0x111")
	numClusters := 2

	clusters := createClusters(committee, latencyMap, self, numClusters)
	clusters.Prepare(self)
	remoteBuckets, localBuckets := clusters.ComputeLatencyBuckets()

	var maxLatency uint = 150
	maxLatency = uint(float64(maxLatency) * constants.MaxLatencyCapFactor)
	// Verify bucket size and assignments
	expectedBucketSize := float64(maxLatency-50) / float64(numClusters) // (max - min) / numClusters
	assert.Equal(t, expectedBucketSize, clusters.bucketSize, "Expected correct bucket size")
	assert.Equal(t, numClusters, len(remoteBuckets), "Expected 2 remote buckets")
	assert.Equal(t, numClusters, len(localBuckets), "Expected 2 local buckets")
	assert.Contains(t, remoteBuckets[1], Node{Addr: common.HexToAddress("0x222"), Lat: 100, ClusterID: 1}, "Expected node 0x222 in remote bucket 1")
	assert.Contains(t, localBuckets[1], Node{Addr: common.HexToAddress("0x333"), Lat: 150, ClusterID: 0}, "Expected node 0x333 in local bucket 1")
	assert.Nil(t, localBuckets[0], nil, "Expected nil in local bucket 0")
}

func TestClusters_AssignRemoteNodes(t *testing.T) {
	committee := []common.Address{
		common.HexToAddress("0x111"), // Cluster 0
		common.HexToAddress("0x222"), // Cluster 1
		common.HexToAddress("0x333"), // Cluster 0
	}
	latencyMap := map[common.Address]uint{
		common.HexToAddress("0x111"): 50,
		common.HexToAddress("0x222"): 100,
		common.HexToAddress("0x333"): 150,
	}
	self := common.HexToAddress("0x111")
	numClusters := 2

	clusters := createClusters(committee, latencyMap, self, numClusters)
	clusters.Prepare(self)
	remoteBuckets, _ := clusters.ComputeLatencyBuckets()
	clusters.AssignRemoteNodes(remoteBuckets, numClusters)

	// Verify bucket nodes
	assert.Len(t, clusters.BucketNodes(), 1, "Expected one remote cluster assigned")
	for _, node := range clusters.BucketNodes() {
		assert.Equal(t, 1, node.ClusterID, "Expected node from cluster 1")
		assert.Equal(t, common.HexToAddress("0x222"), node.Addr, "Expected node 0x222")
	}
}

func TestClusters_AssignRemoteFallbacks(t *testing.T) {
	committee := []common.Address{
		common.HexToAddress("0x111"), // Cluster 0
		common.HexToAddress("0x222"), // Cluster 1
		common.HexToAddress("0x333"), // Cluster 0
		common.HexToAddress("0x444"), // Cluster 1
	}
	latencyMap := map[common.Address]uint{
		common.HexToAddress("0x111"): 50,
		common.HexToAddress("0x222"): 100,
		common.HexToAddress("0x333"): 150,
		common.HexToAddress("0x444"): 160,
	}
	self := common.HexToAddress("0x111")
	numClusters := 2

	clusters := createClusters(committee, latencyMap, self, numClusters)
	clusters.Prepare(self)
	remoteBuckets, _ := clusters.ComputeLatencyBuckets()
	clusters.AssignRemoteNodes(remoteBuckets, numClusters)
	clusters.AssignRemoteFallbacks()

	// Verify fallback nodes
	fallbacks := clusters.BucketFallBacks()
	assert.NotEmpty(t, fallbacks, "Expected fallback nodes")
	for bucketIdx, nodes := range fallbacks {
		for _, node := range nodes {
			if primary, exists := clusters.BucketNodes()[bucketIdx]; exists {
				assert.NotEqual(t, primary.Addr, node.Addr, "Fallback node should not be primary")
			}
			assert.Equal(t, 1, node.ClusterID, "Expected fallback from cluster 1")
		}
	}
}

func TestClusters_PreselectLocalNodes(t *testing.T) {
	committee := []common.Address{
		common.HexToAddress("0x111"), // Cluster 0
		common.HexToAddress("0x222"), // Cluster 1
		common.HexToAddress("0x333"), // Cluster 0
	}
	latencyMap := map[common.Address]uint{
		common.HexToAddress("0x111"): 50,
		common.HexToAddress("0x222"): 100,
		common.HexToAddress("0x333"): 150,
	}
	self := common.HexToAddress("0x111")
	numClusters := 2

	clusters := createClusters(committee, latencyMap, self, numClusters)
	clusters.Prepare(self)
	_, localBuckets := clusters.ComputeLatencyBuckets()
	clusters.PreselectLocalNodes(clusters.base[clusters.ownClusterID], localBuckets, self)

	// Verify local nodes and fallbacks
	assert.Len(t, clusters.LocalBucketNodes(), 1, "Expected sqrt(2) = 1 local node")
	assert.Empty(t, clusters.LocalBucketFallBacks(), "Expected no local fallbacks (only 2 nodes in cluster)")
	assert.Equal(t, common.HexToAddress("0x333"), clusters.LocalBucketNodes()[0].Addr, "Expected node 0x333 in local bucket")
}

func TestClusters_Prepare(t *testing.T) {
	committee := []common.Address{
		common.HexToAddress("0x111"), // Cluster 0
		common.HexToAddress("0x222"), // Cluster 0
	}
	latencyMap := map[common.Address]uint{
		common.HexToAddress("0x111"): 50,
		common.HexToAddress("0x222"): 100,
	}
	self := common.HexToAddress("0x111")
	numClusters := 1

	clusters := createClusters(committee, latencyMap, self, numClusters)
	clusters.Prepare(self)

	// Verify self removed and sorted by latency
	assert.Equal(t, []Node{{Addr: common.HexToAddress("0x222"), Lat: 100, ClusterID: 0}}, clusters.base[0], "Expected self removed and sorted")
}

func TestClusters_IDByAddress(t *testing.T) {
	committee := []common.Address{
		common.HexToAddress("0x111"), // Cluster 0
		common.HexToAddress("0x222"), // Cluster 1
	}
	latencyMap := map[common.Address]uint{
		common.HexToAddress("0x111"): 50,
		common.HexToAddress("0x222"): 100,
	}
	self := common.HexToAddress("0x111")
	numClusters := 2

	clusters := createClusters(committee, latencyMap, self, numClusters)

	assert.Equal(t, 0, clusters.IDByAddress(common.HexToAddress("0x111")), "Expected cluster 0 for 0x111")
	assert.Equal(t, 1, clusters.IDByAddress(common.HexToAddress("0x222")), "Expected cluster 1 for 0x222")
	assert.Equal(t, -1, clusters.IDByAddress(common.HexToAddress("0x333")), "Expected -1 for unknown address")
}

func TestClusters_GetNode(t *testing.T) {
	committee := []common.Address{
		common.HexToAddress("0x111"), // Cluster 0
	}
	latencyMap := map[common.Address]uint{
		common.HexToAddress("0x111"): 50,
	}
	self := common.HexToAddress("0x111")
	numClusters := 1

	clusters := createClusters(committee, latencyMap, self, numClusters)

	node, err := clusters.GetNode(0, common.HexToAddress("0x111"))
	assert.NoError(t, err, "Expected no error for existing node")
	assert.Equal(t, Node{Addr: common.HexToAddress("0x111"), Lat: 50, ClusterID: 0}, node, "Expected correct node")

	_, err = clusters.GetNode(0, common.HexToAddress("0x222"))
	assert.Error(t, err, "Expected error for non-existent node")
}

func TestClusters_LatencyByAddress(t *testing.T) {
	committee := []common.Address{
		common.HexToAddress("0x111"), // Cluster 0
	}
	latencyMap := map[common.Address]uint{
		common.HexToAddress("0x111"): 50,
	}
	self := common.HexToAddress("0x111")
	numClusters := 1

	clusters := createClusters(committee, latencyMap, self, numClusters)

	node, err := clusters.LatencyByAddress(common.HexToAddress("0x111"))
	assert.NoError(t, err, "Expected no error for existing address")
	assert.Equal(t, Node{Addr: common.HexToAddress("0x111"), Lat: 50, ClusterID: 0}, node, "Expected correct node")

	_, err = clusters.LatencyByAddress(common.HexToAddress("0x333"))
	assert.Error(t, err, "Expected error for non-existent address")
}

func TestClusters_MembersByID(t *testing.T) {
	committee := []common.Address{
		common.HexToAddress("0x111"), // Cluster 0
		common.HexToAddress("0x222"), // Cluster 1
	}
	latencyMap := map[common.Address]uint{
		common.HexToAddress("0x111"): 50,
		common.HexToAddress("0x222"): 100,
	}
	self := common.HexToAddress("0x111")
	numClusters := 2

	clusters := createClusters(committee, latencyMap, self, numClusters)

	members := clusters.MembersByID(0)
	assert.Equal(t, []Node{{Addr: common.HexToAddress("0x111"), Lat: 50, ClusterID: 0}}, members, "Expected correct members for cluster 0")
}

func TestClusters_Self(t *testing.T) {
	committee := []common.Address{
		common.HexToAddress("0x111"),
	}
	latencyMap := map[common.Address]uint{
		common.HexToAddress("0x111"): 50,
	}
	self := common.HexToAddress("0x111")
	numClusters := 1

	clusters := createClusters(committee, latencyMap, self, numClusters)

	assert.Equal(t, self, clusters.Self(), "Expected correct self address")
}

func TestNetwork_Clusters(t *testing.T) {
	network := &Network{}
	clusters := New(
		[]common.Address{common.HexToAddress("0x111")},
		map[common.Address]uint{common.HexToAddress("0x111"): 50},
		common.HexToAddress("0x111"),
	)
	network.UpdateClusters(clusters)

	retrieved := network.Clusters()
	assert.Equal(t, clusters, retrieved, "Expected correct clusters")
}

func TestNetwork_UpdateClusters_Concurrent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	network := &Network{}
	var wg sync.WaitGroup
	numGoroutines := 10

	for i := 0; i < numGoroutines; i++ {
		wg.Add(2)
		go func(i int) {
			defer wg.Done()
			clusters := New(
				[]common.Address{common.HexToAddress("0x111")},
				map[common.Address]uint{common.HexToAddress("0x111"): uint(50 + i)},
				common.HexToAddress("0x111"),
			)
			network.UpdateClusters(clusters)
		}(i)
		go func() {
			defer wg.Done()
			network.Clusters()
		}()
	}

	wg.Wait()
}
