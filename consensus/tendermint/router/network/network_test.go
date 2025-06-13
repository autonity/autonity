package network

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/router/constants"
)

func Test_SameClusteringViewForValidators(_ *testing.T) {
	//todo
}
func TestNew_ValidCommittee(t *testing.T) {
	committee := []common.Address{
		common.HexToAddress("0x111"), // Cluster 0
		common.HexToAddress("0x222"), // Cluster 1
		common.HexToAddress("0x333"), // Cluster 0
		common.HexToAddress("0x444"), // Cluster 1
	}
	latencyMap := map[common.Address]uint{
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
	// self latency not considered
	assert.Equal(t, uint(100), clusters.minLatency, "Expected min latency 100")
	assert.Equal(t, uint(200*constants.MaxLatencyCapFactor), clusters.maxLatency, "Expected max latency capped")
	// Cluster 0: 0x111, 0x333
	// Cluster 1: 0x222, 0x444 (sorted by latency)
	assert.Equal(t, []Node{{Addr: common.HexToAddress("0x111"), Lat: 132, ClusterID: 0}, {Addr: common.HexToAddress("0x333"), Lat: 150, ClusterID: 0}}, clusters.base[0], "Expected cluster 0 nodes")
	assert.Equal(t, []Node{
		{Addr: common.HexToAddress("0x222"), Lat: 100, ClusterID: 1},
		{Addr: common.HexToAddress("0x444"), Lat: 200, ClusterID: 1},
	}, clusters.base[1], "Expected cluster 1 nodes")
	// Verify address-to-cluster mapping
	assert.Equal(t, 0, clusters.IDByAddress(common.HexToAddress("0x111")), "Expected 0x111 in cluster 0")
	assert.Equal(t, 1, clusters.IDByAddress(common.HexToAddress("0x222")), "Expected 0x222 in cluster 1")
	// Verify bucket assignments
	assert.NotEmpty(t, clusters.BucketNodes(), "Expected remote bucket nodes assigned")
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
	clusters.ComputeLatencyBuckets()

	var maxLatency uint = 150
	maxLatency = uint(float64(maxLatency) * constants.MaxLatencyCapFactor)
	expectedBucketSize := float64(maxLatency-100) / float64(numClusters)
	assert.Equal(t, expectedBucketSize, clusters.bucketSize, "Expected correct bucket size")
	assert.Equal(t, numClusters, len(clusters.bucketNodes), "Expected 2 buckets")
	assert.Contains(t, clusters.bucketNodes[0], Node{Addr: common.HexToAddress("0x222"), Lat: 100, ClusterID: 1}, "Expected node 0x222 in remote bucket 0")
}

func TestClusters_ComputeLatencyBuckets_EqualLatencies(t *testing.T) {
	committee := []common.Address{
		common.HexToAddress("0x111"), // Cluster 0
		common.HexToAddress("0x222"), // Cluster 1
		common.HexToAddress("0x333"), // Cluster 1
	}
	latencyMap := map[common.Address]uint{
		common.HexToAddress("0x111"): 100,
		common.HexToAddress("0x222"): 100,
	}
	self := common.HexToAddress("0x333")
	numClusters := 2

	clusters, err := createClusters(committee, latencyMap, self, numClusters)
	assert.NoError(t, err, "Expected no error for create cluster")
	clusters.ComputeLatencyBuckets()

	assert.Equal(t, 1.0, clusters.bucketSize, "Expected bucket size 1 for equal latencies")
	assert.Contains(t, clusters.BucketNodes()[0], Node{Addr: common.HexToAddress("0x111"), Lat: 100, ClusterID: 0}, "Expected node 0x222 in bucket 0")
}

func TestClusters_AssignNodesToLatencyBuckets(t *testing.T) {

	t.Run("Basic Assignment", func(t *testing.T) {
		committee := []common.Address{
			common.HexToAddress("0x111"), // self
			common.HexToAddress("0x222"),
			common.HexToAddress("0x333"),
		}
		latencyMap := map[common.Address]uint{
			common.HexToAddress("0x111"): 50,
			common.HexToAddress("0x222"): 100,
			common.HexToAddress("0x333"): 150,
		}
		self := common.HexToAddress("0x111")
		numClusters := 3

		clusters, err := createClusters(committee, latencyMap, self, numClusters)
		assert.NoError(t, err)

		clusters.ComputeLatencyBuckets()

		assert.Len(t, clusters.bucketNodes, 2, "Expected two remote clusters assigned")
		// Find the bucket for node 0x222
		found := false
		for _, nodes := range clusters.bucketNodes {
			if len(nodes) > 0 && nodes[0].Addr == common.HexToAddress("0x222") {
				assert.Equal(t, 1, nodes[0].ClusterID, "Expected node from cluster 1")
				found = true
			}
		}
		assert.True(t, found, "Node 0x222 should have been assigned")
	})

	t.Run("Fills Skipped Buckets in Second Pass", func(t *testing.T) {
		committee := []common.Address{
			common.HexToAddress("0x111"), // self
			common.HexToAddress("0x222"), // C1
			common.HexToAddress("0x333"), // C2, but will be assigned same cluster as 0x222
		}
		latencyMap := map[common.Address]uint{
			common.HexToAddress("0x111"): 50,
			common.HexToAddress("0x222"): 100,
			common.HexToAddress("0x333"): 150,
		}
		self := common.HexToAddress("0x111")
		numClusters := 2 // Force 0x222 and 0x333 into the same remote cluster (C1)

		clusters, err := createClusters(committee, latencyMap, self, numClusters)
		assert.NoError(t, err)

		clusters.ComputeLatencyBuckets()

		assert.Len(t, clusters.bucketNodes, 2, "Expected both buckets to be filled")

		bucketForNode222 := clusters.getBucketIndex(100)
		assert.Equal(t, common.HexToAddress("0x222"), clusters.bucketNodes[bucketForNode222][0].Addr)

		bucketForNode333 := clusters.getBucketIndex(150)
		assert.NotNil(t, clusters.bucketNodes[bucketForNode333], "Expected bucket for 0x333 to be filled")
		assert.Equal(t, common.HexToAddress("0x333"), clusters.bucketNodes[bucketForNode333][0].Addr)
	})

	t.Run("Assigns Fallbacks Correctly", func(t *testing.T) {
		committee := []common.Address{
			common.HexToAddress("0xSELF"), // self
			common.HexToAddress("0xAAA"),
			common.HexToAddress("0xBBB"),
			common.HexToAddress("0xCCC"),
		}
		latencyMap := map[common.Address]uint{
			common.HexToAddress("0xSELF"): 50,
			common.HexToAddress("0xAAA"):  90,
			common.HexToAddress("0xBBB"):  100,
			common.HexToAddress("0xCCC"):  110,
		}
		self := common.HexToAddress("0xSELF")
		numClusters := 2 // All remote nodes go into C1

		clusters, err := createClusters(committee, latencyMap, self, numClusters)
		assert.NoError(t, err)
		clusters.ComputeLatencyBuckets()

		bucketIdx := clusters.getBucketIndex(90)
		assignedNodes := clusters.bucketNodes[bucketIdx]

		assert.Len(t, assignedNodes, 2, "Expected 1 primary and 1 fallback node")
		assert.Equal(t, common.HexToAddress("0xAAA"), assignedNodes[0].Addr, "Expected 0xAAA to be primary")
		fallbackAddrs := []common.Address{assignedNodes[1].Addr}
		assert.Contains(t, fallbackAddrs, common.HexToAddress("0xCCC"))
	})

	t.Run("Complex Scenario with Both Passes", func(t *testing.T) {
		committee := []common.Address{
			common.HexToAddress("0xSELF"), // C0
			common.HexToAddress("0xAAA"),  // C1
			common.HexToAddress("0xBBB"),  // C2
			common.HexToAddress("0xCCC"),  // C0
			common.HexToAddress("0xDDD"),  // C1
		}
		latencyMap := map[common.Address]uint{
			common.HexToAddress("0xSELF"): 50,
			common.HexToAddress("0xAAA"):  100, // For bucket 0, from C1
			common.HexToAddress("0xBBB"):  200, // For bucket 2, from C2
			common.HexToAddress("0xCCC"):  250, // Ignored (remote node in own cluster)
			common.HexToAddress("0xDDD"):  150, // For bucket 1, from C1
		}
		self := common.HexToAddress("0xSELF")
		numClusters := 3

		clusters, err := createClusters(committee, latencyMap, self, numClusters)
		assert.NoError(t, err)
		clusters.ComputeLatencyBuckets()

		assert.Len(t, clusters.bucketNodes, 3, "Expected all 3 buckets to be filled")

		// Check pass 1 assignments
		assert.Equal(t, common.HexToAddress("0xAAA"), clusters.bucketNodes[clusters.getBucketIndex(100)][0].Addr)
		assert.Equal(t, common.HexToAddress("0xBBB"), clusters.bucketNodes[clusters.getBucketIndex(200)][0].Addr)

		// Check pass 2 assignment
		assert.Equal(t, common.HexToAddress("0xDDD"), clusters.bucketNodes[clusters.getBucketIndex(150)][0].Addr)
	})

	t.Run("Handles No Remote Nodes", func(t *testing.T) {
		committee := []common.Address{
			common.HexToAddress("0x111"), // self
		}
		latencyMap := map[common.Address]uint{
			common.HexToAddress("0x111"): 50,
		}
		self := common.HexToAddress("0x111")
		numClusters := 1

		clusters, err := createClusters(committee, latencyMap, self, numClusters)
		assert.NoError(t, err)
		clusters.ComputeLatencyBuckets()

		assert.Len(t, clusters.bucketNodes, 0, "Expected no buckets to be assigned")
	})
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
		common.HexToAddress("0x222"), // Cluster 0
	}
	latencyMap := map[common.Address]uint{
		common.HexToAddress("0x111"): 50,
		common.HexToAddress("0x222"): 50,
	}
	self := common.HexToAddress("0x222")
	numClusters := 1

	clusters, err := createClusters(committee, latencyMap, self, numClusters)
	assert.NoError(t, err, "Expected no error for create cluster")

	node, err := clusters.GetNode(0, common.HexToAddress("0x111"))
	assert.NoError(t, err, "Expected no error for existing node")
	assert.Equal(t, Node{Addr: common.HexToAddress("0x111"), Lat: 50, ClusterID: 0}, node, "Expected correct node")
}

func TestClusters_LatencyByAddress(t *testing.T) {
	committee := []common.Address{
		common.HexToAddress("0x111"), // Cluster 0
		common.HexToAddress("0x222"), // Cluster 0
	}
	latencyMap := map[common.Address]uint{
		common.HexToAddress("0x111"): 50,
		common.HexToAddress("0x222"): 0, // Cluster 0
	}
	self := common.HexToAddress("0x222")
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
	self := common.HexToAddress("0x222")
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
			latencyMap := map[common.Address]uint{common.HexToAddress("0x111"): uint(50 + i%10)} // #nosec
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
