package router

import (
	"errors"
	"math"
	"sort"

	"github.com/autonity/autonity/common"
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

// NewClusters creates and fully initializes a Clusters object from committee addresses and latency data.
func NewClusters(
	committee []common.Address,
	latencyMap map[common.Address]uint,
	self common.Address,
) Clusters {
	numClusters := int(math.Floor(math.Sqrt(float64(len(committee)))))
	clusterViews := make([]ClusterView, numClusters)
	addressToCluster := make(map[common.Address]int)

	// Step 1: Distribute committee members into clusters
	for i, addr := range committee {
		clusterID := i % numClusters
		latency := DefaultLatency
		if lat, ok := latencyMap[addr]; ok {
			latency = int(lat)
		}
		clusterViews[clusterID].Members = append(clusterViews[clusterID].Members, NodeLatency{Addr: addr, Lat: uint(latency)})
		addressToCluster[addr] = clusterID
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
	}
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
	return  NodeLatency{}, errors.New("address not found in any cluster")
}
