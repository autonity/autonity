package latency

import (
	"math"
	"math/rand"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/latency/kmeans"
)

// KmeansClusterSeed TODO(scott): should this be derivable from chain state?
var KmeansClusterSeed = int64(12345)

type Clusters [][]common.Address

// selectK selects k members from each cluster
func (c Clusters) selectK(k int, seed int64) []common.Address {
	var result []common.Address
	r := rand.New(rand.NewSource(seed))
	for _, cluster := range c {
		if len(cluster) <= k {
			result = append(result, cluster...)
		} else {
			selected := make(map[int]struct{})
			for j := 0; j < k; j++ {
				index := r.Intn(len(cluster))
				for _, ok := selected[index]; ok; {
					index = r.Intn(len(cluster))
				}
				result = append(result, cluster[index])
			}
		}
	}
	return result
}

// clusterContaining returns the cluster containing the given address
func (c Clusters) clusterContaining(address common.Address) int {
	for i, cluster := range c {
		for _, member := range cluster {
			if member == address {
				return i
			}
		}
	}
	return -1
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

func fromLatencyMat(latencyMat map[common.Address][]uint8) []kmeans.Observation {
	nodes := make([]kmeans.Observation, len(latencyMat))
	i := 0
	for address, row := range latencyMat {
		floatRow := make([]float64, len(row))
		for j, val := range row {
			floatRow[j] = float64(val)
		}
		nodes[i] = &node{
			address:     address,
			latencyView: floatRow,
		}
		i++
	}
	return nodes
}

func AssignClusters(latencyMat map[common.Address][]uint8, k int) (Clusters, error) {
	nodes := fromLatencyMat(latencyMat)
	km := kmeans.New()
	cstrs, err := km.Partition(nodes, k, KmeansClusterSeed)
	if err != nil {
		return nil, err
	}
	result := make([][]common.Address, k)
	for i, c := range cstrs {
		result[i] = make([]common.Address, len(c.Observations))
		for j, o := range c.Observations {
			result[i][j] = o.(*node).address
		}
	}
	return result, nil
}
