package router

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/log"
)

type NodeLatency struct {
	Addr      common.Address
	Lat       uint
	ClusterID int
}

type ClusterView struct {
	Members []NodeLatency
}

type Clusters struct {
	base                 []ClusterView
	ownClusterID         int
	minLatency           uint
	maxLatency           uint
	bucketSize           float64
	clusterToBucket      map[int][]int         // clusterID -> list of bucket indices it can fill
	bucketNodes          map[int]NodeLatency   // bucketIdx -> preselected node (remote clusters)
	bucketFallbacks      map[int][]NodeLatency // bucketIdx -> fallback nodes (remote clusters)
	localBucketNodes     []NodeLatency         // preselected nodes for local cluster
	localBucketFallbacks []NodeLatency         // fallback nodes for local cluster
	addressToCluster     map[common.Address]int
}

// Base returns the base cluster views
func (c *Clusters) Base() []ClusterView {
	return c.base
}

// createClusters initializes the Clusters struct with committee members and latency data
func createClusters(
	committee []common.Address,
	latencyMap map[common.Address]uint,
	self common.Address,
	numClusters int,
) Clusters {
	c := Clusters{
		base:             make([]ClusterView, numClusters),
		addressToCluster: make(map[common.Address]int),
		minLatency:       uint(math.MaxUint),
		maxLatency:       uint(0),
		clusterToBucket:  make(map[int][]int),
		bucketNodes:      make(map[int]NodeLatency),
		bucketFallbacks:  make(map[int][]NodeLatency),
		ownClusterID:     -1,
	}

	for i, addr := range committee {
		clusterID := i % numClusters
		latency := DefaultLatency
		if lat, ok := latencyMap[addr]; ok {
			latency = lat
		}
		if latency < c.minLatency {
			c.minLatency = latency
		}
		if latency > c.maxLatency {
			c.maxLatency = latency
		}
		c.base[clusterID].Members = append(c.base[clusterID].Members, NodeLatency{Addr: addr, Lat: uint(latency), ClusterID: clusterID})
		c.addressToCluster[addr] = clusterID
		if addr == self {
			c.ownClusterID = clusterID
		}
	}
	c.maxLatency = uint(float64(c.maxLatency) * MaxLatencyCapFactor)

	return c
}

func (c *Clusters) ComputeLatencyBuckets() (remoteBuckets, localBuckets [][]NodeLatency) {
	BucketCount := len(c.base)
	c.bucketSize = float64(c.maxLatency-c.minLatency) / float64(BucketCount)
	if c.bucketSize < 1 {
		c.bucketSize = 1
	}
	remoteBuckets = make([][]NodeLatency, BucketCount)
	localBuckets = make([][]NodeLatency, BucketCount)

	for clusterID, cluster := range c.base {
		bucketSet := make(map[int]bool)
		isLocal := clusterID == c.ownClusterID
		targetBuckets := remoteBuckets
		if isLocal {
			targetBuckets = localBuckets
		} else {
			c.clusterToBucket[clusterID] = []int{}
		}
		for _, node := range cluster.Members {
			bucketIdx := int(float64(node.Lat-c.minLatency) / c.bucketSize)
			if bucketIdx >= BucketCount {
				bucketIdx = BucketCount - 1
			}
			if !isLocal {
				bucketSet[bucketIdx] = true
			}
			targetBuckets[bucketIdx] = append(targetBuckets[bucketIdx], node)
		}
		if !isLocal {
			for bucketIdx := range bucketSet {
				c.clusterToBucket[clusterID] = append(c.clusterToBucket[clusterID], bucketIdx)
			}
		}
	}

	return remoteBuckets, localBuckets
}

func (c *Clusters) AssignRemoteNodes(buckets [][]NodeLatency, numClusters int) {
	BucketCount := len(c.base)
	filledBuckets := make([]bool, BucketCount)
	clusterAssigned := make(map[int]bool)

	// except own cluster
	for len(c.bucketNodes) < numClusters-1 && len(c.bucketNodes) < BucketCount {
		// count cluster candidates for each bucket
		clusterCandidatesForBucket := make([][]int, BucketCount)
		for clusterID, bucketIndices := range c.clusterToBucket {
			if clusterID == c.ownClusterID || clusterAssigned[clusterID] {
				continue
			}
			for _, bucketIdx := range bucketIndices {
				if !filledBuckets[bucketIdx] {
					clusterCandidatesForBucket[bucketIdx] = append(clusterCandidatesForBucket[bucketIdx], clusterID)
				}
			}
		}

		// Find the bucket with the fewest candidates (but at least one)
		// we want to prioritize the buckets with lowest candidates
		minCandidates := math.MaxInt32
		bestBucket := -1
		for i, clusterIDs := range clusterCandidatesForBucket {
			count := len(clusterIDs)
			if count > 0 && count < minCandidates && !filledBuckets[i] {
				minCandidates = count
				bestBucket = i
			}
		}
		// there was no cluster available for this bucket (latency range)
		if bestBucket == -1 {
			break
		}

		var bestNode *NodeLatency
		for _, clusterID := range clusterCandidatesForBucket[bestBucket] {
			for _, node := range buckets[bestBucket] {
				if node.ClusterID == clusterID {
					bestNode = &node
					clusterAssigned[clusterID] = true
					break
				}
			}
			if bestNode != nil {
				break
			}
		}

		if bestNode != nil {
			c.bucketNodes[bestBucket] = *bestNode
			filledBuckets[bestBucket] = true
		} else {
			break
		}
	}

}

func (c *Clusters) AssignRemoteFallbacks() {
	BucketCount := len(c.base)
	for i := 0; i < BucketCount; i++ {
		targetLat := uint(float64(i)*c.bucketSize + c.bucketSize/2 + float64(c.minLatency))
		var clusterID int
		// For buckets with a primary node, use the same cluster
		if primary, exists := c.bucketNodes[i]; exists {
			clusterID = primary.ClusterID
			if clusterID == c.ownClusterID {
				continue // Skip local cluster
			}
			// Select fallback nodes from the same cluster, excluding the primary
			for _, node := range c.base[clusterID].Members {
				if node.Addr == primary.Addr {
					continue // Skip the primary node
				}
				diff := uint(math.Abs(float64(node.Lat - targetLat)))
				if diff < uint(math.MaxUint32) { // Arbitrary threshold to include close nodes
					c.bucketFallbacks[i] = append(c.bucketFallbacks[i], node)
				}
			}
		} else {
			// For buckets without a primary, use unassigned clusters
			usedClusters := make(map[int]bool)
			for _, node := range c.bucketNodes {
				usedClusters[node.ClusterID] = true
			}
			for clusterID := range c.base {
				if clusterID == c.ownClusterID || usedClusters[clusterID] {
					continue
				}
				// Select the closest node from this cluster
				var closest *NodeLatency
				minDiff := uint(math.MaxUint32)
				for _, node := range c.base[clusterID].Members {
					diff := uint(math.Abs(float64(node.Lat - targetLat)))
					if diff < minDiff {
						minDiff = diff
						closest = &node
					}
				}
				if closest != nil {
					c.bucketFallbacks[i] = append(c.bucketFallbacks[i], *closest)
					usedClusters[clusterID] = true
				}
			}
		}
	}
}

// PreselectLocalNodes preselects sqrt(n) nodes and fallbacks for the local cluster
func (c *Clusters) PreselectLocalNodes(localNodes []NodeLatency, localBuckets [][]NodeLatency, self common.Address) {
	BucketCount := len(c.base)
	targetLocalNodes := int(math.Sqrt(float64(len(localNodes))))
	if targetLocalNodes == 0 {
		return
	}

	// Select up to sqrt(n) nodes, prioritizing bucket coverage
	usedNodes := make(map[common.Address]bool)
	for i := 0; i < BucketCount && len(c.localBucketNodes) < targetLocalNodes; i++ {
		for _, node := range localBuckets[i] {
			if !usedNodes[node.Addr] && len(c.localBucketNodes) < targetLocalNodes {
				c.localBucketNodes = append(c.localBucketNodes, node)
				usedNodes[node.Addr] = true
				break
			}
		}
	}

	// Add fallback nodes for local cluster
	for i := 0; i < BucketCount && len(c.localBucketFallbacks) < targetLocalNodes; i++ {
		for _, node := range localBuckets[i] {
			if !usedNodes[node.Addr] && len(c.localBucketFallbacks) < targetLocalNodes {
				c.localBucketFallbacks = append(c.localBucketFallbacks, node)
				usedNodes[node.Addr] = true
			}
		}
	}
}

// PrepareClusters filters self from clusters and sorts members by latency
func (c *Clusters) PrepareClusters(self common.Address) {
	for clusterID := range c.base {
		var peers []NodeLatency
		for _, node := range c.base[clusterID].Members {
			if node.Addr == self {
				continue
			}
			peers = append(peers, node)
		}
		if len(peers) == 0 {
			c.base[clusterID] = ClusterView{}
			continue
		}
		sort.Slice(peers, func(i, j int) bool { return peers[i].Lat < peers[j].Lat })
		c.base[clusterID] = ClusterView{Members: peers}
	}
}

func PrintLatencyBuckets(remoteBuckets, localBuckets [][]NodeLatency, ownClusterID int, bucketSize float64, minLatency uint) {
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

	// Log local buckets
	sb.WriteString(fmt.Sprintf("Local Buckets (Cluster #%d):\n", ownClusterID))
	for bucketIdx, nodes := range localBuckets {
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

// NewClusters creates and fully initializes a Clusters object from committee addresses and latency data.
func NewClusters(
	committee []common.Address,
	latencyMap map[common.Address]uint,
	self common.Address,
) Clusters {
	numClusters := int(math.Floor(math.Sqrt(float64(len(committee)))))

	c := createClusters(committee, latencyMap, self, numClusters)

	c.PrepareClusters(self)

	remoteBuckets, localBuckets := c.ComputeLatencyBuckets()

	PrintLatencyBuckets(remoteBuckets, localBuckets, c.ownClusterID, c.bucketSize, c.minLatency)

	c.AssignRemoteNodes(remoteBuckets, numClusters)

	c.AssignRemoteFallbacks()

	c.PreselectLocalNodes(c.base[c.ownClusterID].Members, localBuckets, self)

	return c
}

// clusterContaining returns the ID of the cluster which containing the given address
func (c *Clusters) clusterContaining(address common.Address) int {
	if idx, exists := c.addressToCluster[address]; exists {
		return idx
	}
	return -1
}

func (c *Clusters) addressToMember(id int, address common.Address) (NodeLatency, error) {
	for _, m := range c.base[id].Members {
		if m.Addr == address {
			return m, nil
		}
	}
	return NodeLatency{}, errors.New("address not found in any cluster")
}

// todo: use this function later
func (c *Clusters) LatencyFromNode(address common.Address) (NodeLatency, error) {
	id := c.clusterContaining(address)
	for _, m := range c.base[id].Members {
		if m.Addr == address {
			return m, nil
		}
	}
	return NodeLatency{}, errors.New("address not found in any cluster")
}
