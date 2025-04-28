package latency

import (
	"math"
	"sort"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
)

type nodeLatency struct {
	Addr common.Address
	Lat  uint
}

type ClusterView struct {
	Members []nodeLatency
	//IQR     uint
	Avg float64
	//Q1      uint
	//Q3      uint
}

type Clusters struct {
	base             []ClusterView
	addressToCluster map[common.Address]int
}

//	func computeIQR(values []uint) (q1, q2, q3, iqr uint) {
//		n := len(values)
//		if n == 0 {
//			return 0, 0, 0, 0
//		}
//		sort.Slice(values, func(i, j int) bool { return values[i] < values[j] })
//		q1 = values[n/4] // 1,2,3,4,5,6,7,8,9,10
//		q2 = values[n/2]
//		q3 = values[(3*n)/4]
//		iqr = q3 - q1
//		return
//	}
//
// NewClusters creates and fully initializes a Clusters object from committee addresses and latency data.
func NewClusters(
	committee []common.Address,
	latencyMap map[common.Address]uint,
	broadcaster consensus.Broadcaster,
	self common.Address,
) *Clusters {
	if len(committee) <= ScaleThresholdForClustering {
		return nil
	}

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
		clusterViews[clusterID].Members = append(clusterViews[clusterID].Members, nodeLatency{Addr: addr, Lat: uint(latency)})
		addressToCluster[addr] = clusterID
	}

	// Step 2: Prepare each cluster (filter, sort, compute stats)
	for clusterID := range clusterViews {
		var peers []nodeLatency
		var latencies []uint
		sum := 0

		for _, node := range clusterViews[clusterID].Members {
			if node.Addr == self {
				continue
			}
			if _, ok := broadcaster.FindPeer(node.Addr); ok {
				lat := latencyMap[node.Addr]
				peers = append(peers, nodeLatency{node.Addr, lat})
				latencies = append(latencies, lat)
				sum += int(lat)
			}
		}

		if len(peers) == 0 {
			clusterViews[clusterID] = ClusterView{}
			continue
		}

		sort.Slice(peers, func(i, j int) bool { return peers[i].Lat < peers[j].Lat })
		//q1, _, q3, iqr := computeIQR(latencies)
		avg := float64(sum) / float64(len(latencies))

		clusterViews[clusterID] = ClusterView{
			Members: peers,
			//IQR:     iqr,
			Avg: avg,
			//Q1:      q1,
			//Q3:      q3,
		}
	}

	return &Clusters{
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
