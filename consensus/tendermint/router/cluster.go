package router

import (
	"bytes"
	"errors"
	"math"
	"sort"
	"sync"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/router/kmeans"
	"github.com/autonity/autonity/log"
)

type NodeLatency struct {
	Addr common.Address
	Lat  uint
}

type ClusterView struct {
	Members []NodeLatency
}

type Clusters struct {
	base             []ClusterView
	addressToCluster map[common.Address]int
}

type node struct {
	address     common.Address
	latencyView []float64
}

func (n *node) Coordinates() kmeans.Coordinates {
	return n.latencyView
}

func (n *node) Distance(p2 kmeans.Coordinates) float64 {
	var r float64
	for i, v := range n.latencyView {
		r += math.Pow(v-p2[i], 2)
	}
	return r
}

func fromLatencyMat(committee []common.Address, latencyMat [][]uint8) []kmeans.Observation {
	nodes := make([]kmeans.Observation, len(latencyMat))
	i := 0
	for index, row := range latencyMat {
		floatRow := make([]float64, len(row))
		for j, val := range row {
			// a value zero means the corresponding measurer did not measure the latency
			// to the target node, in this case we set it to the median of uint8.
			if val == 0 {
				val = math.MaxUint8 / 2
			}
			floatRow[j] = float64(val)
		}
		nodes[i] = &node{
			address:     committee[index],
			latencyView: floatRow,
		}
		i++
	}
	return nodes
}

func assignClusters(committee []common.Address, latencyMat [][]uint8, k int) ([][]common.Address, error) {
	nodes := fromLatencyMat(committee, latencyMat)
	km := kmeans.New()
	kmClusters, err := km.Partition(nodes, k)
	if err != nil {
		return nil, err
	}

	optimalSize := k
	var result [][]common.Address
	for _, cluster := range kmClusters {
		// Collect addresses from cluster observations
		var addresses []common.Address
		for _, obs := range cluster.Observations {
			addresses = append(addresses, obs.(*node).address)
		}

		// Sort addresses lexicographically to make it deterministic.
		sort.Slice(addresses, func(a, b int) bool {
			return bytes.Compare(addresses[a][:], addresses[b][:]) < 0
		})

		// split large clusters into multiple ones.
		if len(addresses) >= optimalSize*2 {
			numSlices := len(addresses) / optimalSize
			for i := 0; i < numSlices; i++ {
				start := i * optimalSize
				end := start + optimalSize
				if i == numSlices-1 {
					end = len(addresses)
				}
				result = append(result, addresses[start:end])
			}
		} else {
			result = append(result, addresses)
		}
	}
	return result, nil
}

// NewClusters creates and fully initializes a Clusters object from committee addresses and latency data.
func NewClusters(
	committee []common.Address,
	latencyMap map[common.Address]uint,
	latencyMat [][]uint8,
	self common.Address,
) (Clusters, error) {
	numClusters := int(math.Floor(math.Sqrt(float64(len(committee)))))
	clusterViews := make([]ClusterView, numClusters)
	addressToCluster := make(map[common.Address]int)

	// Step 1: Distribute committee members into clusters
	if latencyMat != nil {
		clusteredAddresses, err := assignClusters(committee, latencyMat, numClusters)
		// number of clusters can change due to merging/splitting
		if len(clusteredAddresses) != numClusters {
			clusterViews = make([]ClusterView, len(clusteredAddresses))
		}
		if err != nil {
			return Clusters{}, err
		}
		for clusterID, cl := range clusteredAddresses {
			for _, addr := range cl {
				latency := DefaultLatency
				if lat, ok := latencyMap[addr]; ok {
					latency = int(lat)
				}
				clusterViews[clusterID].Members = append(clusterViews[clusterID].Members, NodeLatency{Addr: addr, Lat: uint(latency)})
				addressToCluster[addr] = clusterID
			}
		}
	} else {
		for i, addr := range committee {
			latency := DefaultLatency
			if lat, ok := latencyMap[addr]; ok {
				latency = int(lat)
			}
			clusterViews[i%numClusters].Members = append(clusterViews[i%numClusters].Members, NodeLatency{Addr: addr, Lat: uint(latency)})
			addressToCluster[addr] = i % numClusters
		}
	}

	// Step 2: Prepare each cluster (filter, sort, compute stats)
	for clusterID, cluster := range clusterViews {
		var peers []NodeLatency
		var latencies []uint
		sum := 0

		for _, node := range cluster.Members {
			if node.Addr == self {
				continue
			}
			peers = append(peers, node)
			latencies = append(latencies, node.Lat)
			sum += int(node.Lat)
		}

		if len(peers) == 0 {
			clusterViews[clusterID] = ClusterView{}
			continue
		}

		sort.Slice(peers, func(i, j int) bool { return peers[i].Lat < peers[j].Lat })
		clusterViews[clusterID] = ClusterView{
			Members: peers,
		}
	}

	return Clusters{
		base:             clusterViews,
		addressToCluster: addressToCluster,
	}, nil
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

func UpdateClusterLatencies(c Clusters, latencyMap map[common.Address]uint, self common.Address) Clusters {
	for clusterID, cluster := range c.base {
		var peers []NodeLatency
		var latencies []uint
		sum := 0

		for _, node := range cluster.Members {
			if node.Addr == self {
				continue
			}
			node.Lat = latencyMap[node.Addr]
			peers = append(peers, node)
			latencies = append(latencies, node.Lat)
			sum += int(node.Lat)
		}

		if len(peers) == 0 {
			c.base[clusterID] = ClusterView{}
			continue
		}

		sort.Slice(peers, func(i, j int) bool { return peers[i].Lat < peers[j].Lat })
		c.base[clusterID] = ClusterView{
			Members: peers,
		}
	}
	return c
}

type ClusterRotation struct {
	mu                 sync.RWMutex
	previousEpochBlock uint64
	latestEpochBlock   uint64

	lastMatrixLockInBlock   uint64
	latestMatrixLockInBlock uint64

	previousEpochClusters Clusters
	transitionalClusters  Clusters
	latestEpochClusters   Clusters
}

func (cr *ClusterRotation) EpochStart(epochBlock uint64, transitionalClusters Clusters) {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	// new epoch has started, need to create transitional clusters
	cr.previousEpochBlock = cr.latestEpochBlock
	cr.latestEpochBlock = epochBlock

	cr.previousEpochClusters = cr.latestEpochClusters
	cr.transitionalClusters = transitionalClusters
	cr.latestEpochClusters = Clusters{}

	cr.lastMatrixLockInBlock = cr.latestMatrixLockInBlock
	cr.latestMatrixLockInBlock = 0
	log.Info(
		"ClusterRotation: new epoch started",
		"previousEpochBlock", cr.previousEpochBlock,
		"latestEpochBlock", cr.latestEpochBlock,
		"lastMatrixLockInBlock", cr.lastMatrixLockInBlock,
		"latestMatrixLockInBlock", cr.latestMatrixLockInBlock,
	)
}

func (cr *ClusterRotation) LockIn(lockInBlock uint64, cluster Clusters) {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	cr.latestEpochClusters = cluster
	cr.latestMatrixLockInBlock = lockInBlock
}

func (cr *ClusterRotation) GetClusters(height uint64) Clusters {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	if height < cr.previousEpochBlock {
		if height > cr.lastMatrixLockInBlock {
			return cr.previousEpochClusters
		} else {
			// no clusters available for before previous epoch lock in block
			return Clusters{}
		}
	}

	if height < cr.latestMatrixLockInBlock || cr.latestMatrixLockInBlock == 0 {
		return cr.transitionalClusters
	}
	return cr.latestEpochClusters
}

func (cr *ClusterRotation) UpdateLatencies(latencies map[common.Address]uint, self common.Address) {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	// update latencies based on whether we are in the transitional clusters or
	// the latest epoch clusters phase
	if len(cr.latestEpochClusters.base) != 0 {
		cr.latestEpochClusters = UpdateClusterLatencies(cr.latestEpochClusters, latencies, self)
	} else if len(cr.transitionalClusters.base) != 0 {
		cr.transitionalClusters = UpdateClusterLatencies(cr.transitionalClusters, latencies, self)
	}
}
