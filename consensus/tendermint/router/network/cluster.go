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
	base                 [][]Node
	ownClusterID         int
	self                 common.Address
	minLatency           uint
	maxLatency           uint
	bucketSize           float64
	clusterToBucket      map[int][]int  // clusterID -> list of bucket indices it can fill
	bucketNodes          map[int]Node   // bucketIdx -> preselected node (remote clusters)
	bucketFallbacks      map[int][]Node // bucketIdx -> fallback nodes (remote clusters)
	localBucketNodes     []Node         // preselected nodes for local cluster
	localBucketFallbacks []Node         // fallback nodes for local cluster
	addressToCluster     map[common.Address]int
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
		clusterToBucket:  make(map[int][]int),
		bucketNodes:      make(map[int]Node),
		bucketFallbacks:  make(map[int][]Node),
		ownClusterID:     -1,
	}

	for i, addr := range committee {
		clusterID := i % numClusters
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
		if addr == self {
			c.self = addr
			c.ownClusterID = clusterID
		}
	}
	if c.ownClusterID == -1 {
		return Clusters{}, errors.New("self address not in committee")
	}
	c.maxLatency = uint(float64(c.maxLatency) * constants.MaxLatencyCapFactor)
	if c.maxLatency < c.minLatency {
		c.maxLatency = c.minLatency
	}
	return c, nil
}

// Base returns the base cluster views
func (c *Clusters) Base() [][]Node {
	return c.base
}

func (c *Clusters) ComputeLatencyBuckets() (remoteBuckets, localBuckets [][]Node) {
	BucketCount := len(c.base)
	c.bucketSize = float64(c.maxLatency-c.minLatency) / float64(BucketCount)
	if c.bucketSize < 1 {
		c.bucketSize = 1
	}
	// latency buckets are used to group nodes by latency ranges, we create two sets:
	// one for remote nodes and one for the local nodes
	remoteBuckets = make([][]Node, BucketCount)
	localBuckets = make([][]Node, BucketCount)

	for clusterID, cluster := range c.base {
		bucketSet := make(map[int]bool)
		isLocal := clusterID == c.ownClusterID
		targetBuckets := remoteBuckets
		if isLocal {
			targetBuckets = localBuckets
		} else {
			c.clusterToBucket[clusterID] = []int{}
		}
		for _, node := range cluster {
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

func (c *Clusters) AssignRemoteNodes(buckets [][]Node, numClusters int) {
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

		var bestNode *Node
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
			for _, node := range c.base[clusterID] {
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
			for clusterID = range c.base {
				if clusterID == c.ownClusterID || usedClusters[clusterID] {
					continue
				}
				// Select the closest node from this cluster
				var closest *Node
				minDiff := uint(math.MaxUint32)
				for _, node := range c.base[clusterID] {
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
func (c *Clusters) PreselectLocalNodes(localNodes []Node, localBuckets [][]Node) {
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

// Prepare filters self from clusters and sorts members by latency
func (c *Clusters) Prepare(self common.Address) {
	for clusterID := range c.base {
		var peers []Node
		for _, node := range c.base[clusterID] {
			if node.Addr == self {
				continue
			}
			peers = append(peers, node)
		}
		if len(peers) == 0 {
			c.base[clusterID] = []Node{}
			continue
		}
		sort.Slice(peers, func(i, j int) bool { return peers[i].Lat < peers[j].Lat })
		c.base[clusterID] = peers
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

func (c *Clusters) BucketNodes() map[int]Node {
	return c.bucketNodes
}

func (c *Clusters) BucketFallBacks() map[int][]Node {
	return c.bucketFallbacks
}

func (c *Clusters) LocalBucketNodes() []Node {
	return c.localBucketNodes
}

func (c *Clusters) LocalBucketFallBacks() []Node {
	return c.localBucketFallbacks
}
