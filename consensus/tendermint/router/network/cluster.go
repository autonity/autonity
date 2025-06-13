package network

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/router/constants"
	"github.com/autonity/autonity/log"
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
	bucketNodes      map[int][]Node // bucketIdx -> nodes
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
		ownClusterID:     -1,
		self:             self,
	}

	for i, addr := range committee {
		latency := constants.DefaultLatency
		clusterID := i % numClusters
		if addr == self {
			c.ownClusterID = clusterID
			c.addressToCluster[addr] = clusterID
		} else {
			// latency for self node is not considered
			if lat, ok := latencyMap[addr]; ok {
				latency = lat
			}
			if latency < c.minLatency {
				c.minLatency = latency
			}
			if latency > c.maxLatency {
				c.maxLatency = latency
			}
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
	bucketCount := len(c.base)
	c.bucketSize = float64(c.maxLatency-c.minLatency) / float64(bucketCount)
	if c.bucketSize < 1 {
		c.bucketSize = 1
	}
	buckets := make([][]Node, bucketCount)

	for _, cluster := range c.base {
		bucketSet := make(map[int]bool)
		for _, node := range cluster {
			if node.Addr == c.self {
				continue // skip self for bucket assignment
			}
			bucketIdx := c.getBucketIndex(node.Lat)
			buckets[bucketIdx] = append(buckets[bucketIdx], node)
			bucketSet[bucketIdx] = true
		}
	}

	// Sort each bucket by latency
	for _, bucket := range buckets {
		sort.Slice(bucket, func(i, j int) bool {
			return bucket[i].Lat < bucket[j].Lat
		})
	}

	printLatencyBuckets(buckets, c.bucketSize, c.minLatency)
	c.assignNodesToLatencyBuckets(buckets)
	return buckets
}

// assignNodesToLatencyBuckets attempts to assign one primary node and up to 2 fallback nodes per  cluster to buckets
// although there is no guarantee that all clusters with have a representation, but we make sure that all the buckets
// will have at least one node assigned
func (c *Clusters) assignNodesToLatencyBuckets(buckets [][]Node) {
	bucketCount := len(buckets)
	filledBuckets := make([]bool, bucketCount)
	clusterAssigned := make(map[int]bool)

	// fallback attempts to select 2 more nodes from the same cluster which
	// are adjacent to the selected node in the cluster
	selectFallbacks := func(node Node) []Node {
		fallbacks := make([]Node, 0, 2)
		var selectedNodeIndex int
		for i, n := range c.base[node.ClusterID] {
			if n.Addr == node.Addr {
				selectedNodeIndex = i
				break
			}
		}
		clusterNodes := c.base[node.ClusterID]
		if selectedNodeIndex > 0 {
			// Add previous node as fallback if it exists
			n := clusterNodes[selectedNodeIndex-1]
			if n.Addr != c.self {
				fallbacks = append(fallbacks, n)
			}
		}
		if selectedNodeIndex < len(clusterNodes)-1 {
			// add next node as fallback if it exists
			n := clusterNodes[selectedNodeIndex+1]
			if n.Addr != c.self {
				fallbacks = append(fallbacks, n)
			}
		}
		return fallbacks
	}

	// 1st pass: maximize cluster assignment to buckets
	for bucketIdx, bucketNodes := range buckets {
		for _, node := range bucketNodes {
			if clusterAssigned[node.ClusterID] {
				// Skip if this cluster is already assigned
				continue
			}
			c.bucketNodes[bucketIdx] = append(c.bucketNodes[bucketIdx], node)
			c.bucketNodes[bucketIdx] = append(c.bucketNodes[bucketIdx], selectFallbacks(node)...)
			clusterAssigned[node.ClusterID] = true
			filledBuckets[bucketIdx] = true
			break
		}
	}

	// 2nd pass: Fill remaining buckets with nodes from any of the cluster
	for bucketIdx, bucketNodes := range buckets {
		if filledBuckets[bucketIdx] || len(bucketNodes) == 0 {
			continue
		}
		bestNode := bucketNodes[0] // pick the first node in the bucket
		c.bucketNodes[bucketIdx] = append(c.bucketNodes[bucketIdx], bestNode)
		c.bucketNodes[bucketIdx] = append(c.bucketNodes[bucketIdx], selectFallbacks(bestNode)...)
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

func printLatencyBuckets(remoteBuckets [][]Node, bucketSize float64, minLatency uint) {
	var sb strings.Builder
	sb.WriteString("\nLatency Buckets for Clusters:\n")

	// Log remote buckets
	sb.WriteString("Remote Buckets:\n")
	for bucketIdx, nodes := range remoteBuckets {
		lowerLat := uint(float64(bucketIdx)*bucketSize) + minLatency
		upperLat := uint(float64(bucketIdx+1)*bucketSize) + minLatency
		sb.WriteString(fmt.Sprintf("  Bucket #%d (Latency %d-%d ms): %d nodes\n", bucketIdx, lowerLat, upperLat, len(nodes)))
		if len(nodes) == 0 {
			sb.WriteString("    [Empty]\n")
			continue
		}
		for _, node := range nodes {
			sb.WriteString(fmt.Sprintf("    Node: %s, Latency: %d ms, ClusterID: %d\n", node.Addr.Hex(), node.Lat, node.ClusterID))
		}
	}

	log.Info(sb.String())
}
