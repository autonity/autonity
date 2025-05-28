package network

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/router/constants"
)

func TestNew_ValidCommittee(t *testing.T) {
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
	self := common.HexToAddress("0x111")
	numClusters := 2 // sqrt(4) = 2

	clusters, err := New(committee, latencyMap, self)
	assert.NoError(t, err, "Expected no error for valid committee")

	// Verify cluster structure
	assert.Equal(t, numClusters, len(clusters.Base()), "Expected 2 clusters")
	assert.Equal(t, 0, clusters.ID(), "Expected self in cluster 0")
	assert.Equal(t, self, clusters.Self(), "Expected self address")
	assert.Equal(t, 50, int(clusters.minLatency), "Expected min latency 50")
	assert.Equal(t, uint(200*constants.MaxLatencyCapFactor), clusters.maxLatency, "Expected max latency capped")
	// Cluster 0: 0x333 (self removed by Prepare)
	// Cluster 1: 0x222, 0x444 (sorted by latency)
	assert.Equal(t, []Node{{Addr: common.HexToAddress("0x333"), Lat: 150, ClusterID: 0}}, clusters.base[0], "Expected cluster 0 nodes")
	assert.Equal(t, []Node{
		{Addr: common.HexToAddress("0x222"), Lat: 100, ClusterID: 1},
		{Addr: common.HexToAddress("0x444"), Lat: 200, ClusterID: 1},
	}, clusters.base[1], "Expected cluster 1 nodes")
	// Verify address-to-cluster mapping
	assert.Equal(t, 0, clusters.IDByAddress(common.HexToAddress("0x111")), "Expected 0x111 in cluster 0")
	assert.Equal(t, 1, clusters.IDByAddress(common.HexToAddress("0x222")), "Expected 0x222 in cluster 1")
	// Verify bucket assignments
	assert.NotEmpty(t, clusters.BucketNodes(), "Expected remote bucket nodes assigned")
	assert.NotEmpty(t, clusters.LocalBucketNodes(), "Expected local bucket nodes assigned")
}

func TestNew_SelfNotInCommittee(t *testing.T) {
	committee := []common.Address{
		common.HexToAddress("0x222"),
		common.HexToAddress("0x333"),
	}
	latencyMap := map[common.Address]uint{
		common.HexToAddress("0x222"): 100,
		common.HexToAddress("0x333"): 150,
	}
	self := common.HexToAddress("0x111")

	_, err := New(committee, latencyMap, self)
	assert.Error(t, err, "Expected error for self not in committee")
	assert.Equal(t, "self address not in committee", err.Error(), "Expected correct error message")
}

func TestNew_EmptyCommittee(t *testing.T) {
	committee := []common.Address{}
	latencyMap := map[common.Address]uint{}
	self := common.HexToAddress("0x111")

	_, err := New(committee, latencyMap, self)
	assert.Error(t, err, "Expected error for empty committee")
	assert.Equal(t, "committee cannot be empty", err.Error(), "Expected correct error message")
}

func TestNew_SingleMemberCommittee(t *testing.T) {
	committee := []common.Address{common.HexToAddress("0x111")}
	latencyMap := map[common.Address]uint{common.HexToAddress("0x111"): 50}
	self := common.HexToAddress("0x111")
	numClusters := 1 // sqrt(1) = 1

	clusters, err := New(committee, latencyMap, self)
	assert.NoError(t, err, "Expected no error for single-member committee")

	assert.Equal(t, numClusters, len(clusters.Base()), "Expected 1 cluster")
	assert.Equal(t, 0, clusters.ID(), "Expected self in cluster 0")
	assert.Equal(t, self, clusters.Self(), "Expected self address")
	assert.Equal(t, 50, int(clusters.minLatency), "Expected min latency 50")
	// no capping because min and max are the same
	assert.Equal(t, 50, int(clusters.maxLatency), "Expected max latency 50")
	assert.Empty(t, clusters.base[0], "Expected empty cluster 0 (self removed)")
	assert.Empty(t, clusters.BucketNodes(), "Expected no remote bucket nodes")
	assert.Empty(t, clusters.LocalBucketNodes(), "Expected no local bucket nodes")
}

func TestNew_InvalidLatencyMap(t *testing.T) {
	committee := []common.Address{
		common.HexToAddress("0x111"),
		common.HexToAddress("0x222"),
	}
	latencyMap := map[common.Address]uint{
		common.HexToAddress("0x222"): 100, // Missing 0x111
	}
	self := common.HexToAddress("0x111")

	clusters, err := New(committee, latencyMap, self)
	assert.NoError(t, err, "Expected no error for partial latency map")

	_, err = clusters.LatencyByAddress(common.HexToAddress("0x111"))
	assert.Equal(t, "address not found in any cluster", err.Error(), "Expected correct error message")
}

func TestCreateClusters_SelfNotInCommittee(t *testing.T) {
	committee := []common.Address{
		common.HexToAddress("0x222"),
		common.HexToAddress("0x333"),
	}
	latencyMap := map[common.Address]uint{
		common.HexToAddress("0x222"): 100,
		common.HexToAddress("0x333"): 150,
	}
	self := common.HexToAddress("0x111")
	numClusters := 2

	_, err := createClusters(committee, latencyMap, self, numClusters)
	assert.Error(t, err, "Expected error for self not in committee")
	assert.Equal(t, "self address not in committee", err.Error(), "Expected correct error message")
}

func TestCreateClusters_EmptyCommittee(t *testing.T) {
	committee := []common.Address{}
	latencyMap := map[common.Address]uint{}
	self := common.HexToAddress("0x111")
	numClusters := 1

	_, err := createClusters(committee, latencyMap, self, numClusters)
	assert.Error(t, err, "Expected error for empty committee")
	assert.Equal(t, "committee cannot be empty", err.Error(), "Expected correct error message")
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

	clusters, err := createClusters(committee, latencyMap, self, numClusters)
	assert.NoError(t, err, "Expected no error for create cluster")
	clusters.Prepare(self)
	remoteBuckets, localBuckets := clusters.ComputeLatencyBuckets()

	var maxLatency uint = 150
	maxLatency = uint(float64(maxLatency) * constants.MaxLatencyCapFactor)
	expectedBucketSize := float64(maxLatency-50) / float64(numClusters)
	assert.Equal(t, expectedBucketSize, clusters.bucketSize, "Expected correct bucket size")
	assert.Equal(t, numClusters, len(remoteBuckets), "Expected 2 remote buckets")
	assert.Equal(t, numClusters, len(localBuckets), "Expected 2 local buckets")
	assert.Contains(t, remoteBuckets[1], Node{Addr: common.HexToAddress("0x222"), Lat: 100, ClusterID: 1}, "Expected node 0x222 in remote bucket 1")
	assert.Contains(t, localBuckets[1], Node{Addr: common.HexToAddress("0x333"), Lat: 150, ClusterID: 0}, "Expected node 0x333 in local bucket 1")
	assert.Empty(t, localBuckets[0], "Expected empty local bucket 0")
	assert.Contains(t, clusters.clusterToBucket[1], 1, "Expected cluster 1 mapped to bucket 1")
}

func TestClusters_ComputeLatencyBuckets_EqualLatencies(t *testing.T) {
	committee := []common.Address{
		common.HexToAddress("0x111"), // Cluster 0
		common.HexToAddress("0x222"), // Cluster 1
	}
	latencyMap := map[common.Address]uint{
		common.HexToAddress("0x111"): 100,
		common.HexToAddress("0x222"): 100,
	}
	self := common.HexToAddress("0x111")
	numClusters := 2

	clusters, err := createClusters(committee, latencyMap, self, numClusters)
	assert.NoError(t, err, "Expected no error for create cluster")
	clusters.Prepare(self)
	remoteBuckets, _ := clusters.ComputeLatencyBuckets()

	assert.Equal(t, 1.0, clusters.bucketSize, "Expected bucket size 1 for equal latencies")
	assert.Equal(t, numClusters, len(remoteBuckets), "Expected 2 remote buckets")
	assert.Contains(t, remoteBuckets[0], Node{Addr: common.HexToAddress("0x222"), Lat: 100, ClusterID: 1}, "Expected node 0x222 in bucket 0")
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

	clusters, err := createClusters(committee, latencyMap, self, numClusters)
	assert.NoError(t, err, "Expected no error for create cluster")

	clusters.Prepare(self)
	remoteBuckets, _ := clusters.ComputeLatencyBuckets()
	clusters.AssignRemoteNodes(remoteBuckets, numClusters)

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

	clusters, err := createClusters(committee, latencyMap, self, numClusters)
	assert.NoError(t, err, "Expected no error for create cluster")
	clusters.Prepare(self)
	remoteBuckets, _ := clusters.ComputeLatencyBuckets()
	clusters.AssignRemoteNodes(remoteBuckets, numClusters)
	clusters.AssignRemoteFallbacks()

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

	clusters, err := createClusters(committee, latencyMap, self, numClusters)
	assert.NoError(t, err, "Expected no error for create cluster")
	clusters.Prepare(self)
	_, localBuckets := clusters.ComputeLatencyBuckets()
	clusters.PreselectLocalNodes(clusters.base[clusters.ownClusterID], localBuckets, self)

	assert.Len(t, clusters.LocalBucketNodes(), 1, "Expected sqrt(2) = 1 local node")
	assert.Empty(t, clusters.LocalBucketFallBacks(), "Expected no local fallbacks")
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

	clusters, err := createClusters(committee, latencyMap, self, numClusters)
	assert.NoError(t, err, "Expected no error for create cluster")
	clusters.Prepare(self)

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

	clusters, err := createClusters(committee, latencyMap, self, numClusters)
	assert.NoError(t, err, "Expected no error for create cluster")

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

	clusters, err := createClusters(committee, latencyMap, self, numClusters)
	assert.NoError(t, err, "Expected no error for create cluster")

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

	clusters, err := createClusters(committee, latencyMap, self, numClusters)
	assert.NoError(t, err, "Expected no error for create cluster")

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

	clusters, err := createClusters(committee, latencyMap, self, numClusters)
	assert.NoError(t, err, "Expected no error for create cluster")

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

	clusters, err := createClusters(committee, latencyMap, self, numClusters)
	assert.NoError(t, err, "Expected no error for create cluster")

	assert.Equal(t, self, clusters.Self(), "Expected correct self address")
}

func TestNetwork_Clusters(t *testing.T) {
	network := &Network{}
	committee := []common.Address{common.HexToAddress("0x111")}
	latencyMap := map[common.Address]uint{common.HexToAddress("0x111"): 50}
	self := common.HexToAddress("0x111")
	clusters, err := New(committee, latencyMap, self)
	assert.NoError(t, err, "Expected no error creating clusters")
	network.UpdateClusters(clusters)

	retrieved := network.Clusters()
	assert.Equal(t, clusters, retrieved, "Expected correct clusters")
}

func TestNetwork_UpdateClusters_Concurrent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	network := &Network{}
	var wg sync.WaitGroup
	numGoroutines := 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(2)
		go func(i int) {
			defer wg.Done()
			committee := []common.Address{common.HexToAddress("0x111")}
			latencyMap := map[common.Address]uint{common.HexToAddress("0x111"): uint(50 + i%10)}
			self := common.HexToAddress("0x111")
			clusters, err := New(committee, latencyMap, self)
			assert.NoError(t, err, "Expected no error for concurrency")
			network.UpdateClusters(clusters)
		}(i)
		go func() {
			defer wg.Done()
			network.Clusters()
		}()
	}

	wg.Wait()
}
