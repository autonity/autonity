package network

import (
	"errors"
	"math"
	"sort"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/router/constants"
)

type Node struct {
	Addr      common.Address
	Lat       uint
	ClusterID int
}

type Clusters struct {
	base             [][]Node
	ownClusterID     int
	self             common.Address
	minLatency       uint
	maxLatency       uint
	bucketSize       float64
	bucketNodes      map[int][]Node // bucketIdx -> preselected node (remote clusters)
	bucketFallbacks  map[int][]Node // bucketIdx -> fallback nodes (remote clusters)
	addressToCluster map[common.Address]int
}

// createClusters initializes the Clusters struct with committee members and latency data
func createClusters(
	committee []common.Address,
	latencyMap map[common.Address]uint,
	self common.Address,
	numClusters int,
) (Clusters, error) {
	if len(committee) == 0 {
		return Clusters{}, errors.New("committee cannot be empty")
	}
	c := Clusters{
		base:             make([][]Node, numClusters),
		addressToCluster: make(map[common.Address]int),
		minLatency:       uint(math.MaxUint),
		maxLatency:       uint(0),
		bucketNodes:      make(map[int][]Node),
		bucketFallbacks:  make(map[int][]Node),
		ownClusterID:     -1,
	}

	for i, addr := range committee {
		clusterID := i % numClusters
		if addr == self {
			c.self = addr
			c.ownClusterID = clusterID
			c.addressToCluster[addr] = clusterID
			continue // Skip adding self to base
		}
		latency := constants.DefaultLatency
		if lat, ok := latencyMap[addr]; ok {
			latency = lat
		}
		if latency < c.minLatency {
			c.minLatency = latency
		}
		if latency > c.maxLatency {
			c.maxLatency = latency
		}
		c.base[clusterID] = append(c.base[clusterID], Node{Addr: addr, Lat: latency, ClusterID: clusterID})
		c.addressToCluster[addr] = clusterID
	}
	if c.ownClusterID == -1 {
		return Clusters{}, errors.New("self address not in committee")
	}
	c.maxLatency = uint(float64(c.maxLatency) * constants.MaxLatencyCapFactor)
	if c.maxLatency < c.minLatency {
		c.maxLatency = c.minLatency
	}

	// Sort each cluster by latency
	for clusterID := range c.base {
		sort.Slice(c.base[clusterID], func(i, j int) bool {
			return c.base[clusterID][i].Lat < c.base[clusterID][j].Lat
		})
	}
	return c, nil
}

// Base returns the base cluster views
func (c *Clusters) Base() [][]Node {
	return c.base
}

// getBucketIndex calculates the bucket index for a given latency
func (c *Clusters) getBucketIndex(latency uint) int {
	if c.bucketSize < 1 {
		return 0
	}
	bucketIdx := int(float64(latency-c.minLatency) / c.bucketSize)
	if bucketIdx >= len(c.base) {
		bucketIdx = len(c.base) - 1
	}
	return bucketIdx
}

func (c *Clusters) ComputeLatencyBuckets() [][]Node {
	BucketCount := len(c.base)
	c.bucketSize = float64(c.maxLatency-c.minLatency) / float64(BucketCount)
	if c.bucketSize < 1 {
		c.bucketSize = 1
	}
	buckets := make([][]Node, BucketCount)

	for clusterID, cluster := range c.base {
		if clusterID == c.ownClusterID {
			continue
		}
		bucketSet := make(map[int]bool)
		for _, node := range cluster {
			bucketIdx := c.getBucketIndex(node.Lat)
			buckets[bucketIdx] = append(buckets[bucketIdx], node)
			bucketSet[bucketIdx] = true
		}
	}

	return buckets
}

// AssignRemoteNodes assigns one primary node and up to 2 fallback nodes per remote cluster to buckets
func (c *Clusters) AssignRemoteNodes(buckets [][]Node, numClusters int) {
	BucketCount := len(buckets)
	filledBuckets := make([]bool, BucketCount)
	clusterAssigned := make(map[int]bool)

	for len(c.bucketNodes) < numClusters-1 && len(c.bucketNodes) < BucketCount {
		// Find bucket with lowest-latency unassigned node
		bestBucket := -1
		lowestLatency := uint(math.MaxUint)
		var bestNode *Node
		var bestNodeIndex int
		for bucketIdx, nodes := range buckets {
			if filledBuckets[bucketIdx] || len(nodes) == 0 {
				continue
			}
			for _, node := range nodes {
				if node.Lat < lowestLatency && !clusterAssigned[node.ClusterID] {
					lowestLatency = node.Lat
					bestBucket = bucketIdx
					bestNode = &node
					// Find index of this node in c.base[clusterID]
					for i, n := range c.base[node.ClusterID] {
						if n.Addr == node.Addr {
							bestNodeIndex = i
							break
						}
					}
				}
			}
		}
		if bestBucket == -1 || bestNode == nil {
			break
		}

		// Assign primary node
		c.bucketNodes[bestBucket] = append(c.bucketNodes[bestBucket], *bestNode)
		filledBuckets[bestBucket] = true
		clusterAssigned[bestNode.ClusterID] = true

		// Assign up to 2 fallback nodes from the same cluster, using subsequent indices
		clusterNodes := c.base[bestNode.ClusterID]
		if bestNodeIndex > 0 {
			// Add previous node as fallback if it exists
			c.bucketFallbacks[bestBucket] = append(c.bucketFallbacks[bestBucket], clusterNodes[bestNodeIndex-1])
		}
		if bestNodeIndex < len(clusterNodes)-1 {
			// add next node as fallback if it exists
			c.bucketFallbacks[bestBucket] = append(c.bucketFallbacks[bestBucket], clusterNodes[bestNodeIndex+1])
		}

	}
}

// IDByAddress returns the ID of the cluster which containing the given address
func (c *Clusters) IDByAddress(address common.Address) int {
	if idx, exists := c.addressToCluster[address]; exists {
		return idx
	}
	return -1
}

func (c *Clusters) GetNode(id int, address common.Address) (Node, error) {
	for _, m := range c.base[id] {
		if m.Addr == address {
			return m, nil
		}
	}
	return Node{}, errors.New("address not found in any cluster")
}

func (c *Clusters) LatencyByAddress(address common.Address) (Node, error) {
	id := c.IDByAddress(address)
	if id == -1 {
		return Node{}, errors.New("address not found in any cluster")
	}
	for _, m := range c.base[id] {
		if m.Addr == address {
			return m, nil
		}
	}
	return Node{}, errors.New("address not found in any cluster")
}

func (c *Clusters) MembersByID(clusterID int) []Node {
	return c.base[clusterID]
}

func (c *Clusters) ID() int {
	return c.ownClusterID
}

func (c *Clusters) Self() common.Address {
	return c.self
}

func (c *Clusters) BucketNodes() map[int][]Node {
	return c.bucketNodes
}
