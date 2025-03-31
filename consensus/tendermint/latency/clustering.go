package latency

import (
	"bytes"
	"github.com/autonity/autonity/consensus"
	"math"
	"math/rand"
	"sort"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/latency/kmeans"
)

// KmeansClusterSeed is the seed to init the initial center position for each cluster.
var KmeansClusterSeed = int64(12345)

type Clusters struct {
	activatedHeight  uint64
	nextEpochHeight  uint64
	base             [][]common.Address
	addressToCluster map[common.Address]int
}

func NewCluster(activationHeight uint64, nextEpochHeight uint64, base [][]common.Address) *Clusters {
	clusters := &Clusters{
		base:            base,
		activatedHeight: activationHeight,
		nextEpochHeight: nextEpochHeight,
	}

	clusters.buildAddressIndex()
	return clusters
}

func (c *Clusters) clusterByID(id int) []common.Address {
	if id >= len(c.base) || id < 0 {
		return nil
	}
	return c.base[id]
}

// buildAddressIndex builds the index for quick querying of node in clusters, it should be called on the setup phase.
func (c *Clusters) buildAddressIndex() {
	c.addressToCluster = make(map[common.Address]int)
	for i, cluster := range c.base {
		for _, member := range cluster {
			c.addressToCluster[member] = i
		}
	}
}

// clusterContaining returns the ID of the cluster which containing the given address
func (c *Clusters) clusterContaining(address common.Address) int {
	if idx, exists := c.addressToCluster[address]; exists {
		return idx
	}
	return -1
}

// selectK selects pseudo random k members from each cluster exclude the selected cluster.
func (c *Clusters) selectK(broadcaster consensus.Broadcaster, k int, seed int64, excepted int) []common.Address {
	var result []common.Address
	r := rand.New(rand.NewSource(seed))

	for i, nodes := range c.base {
		if i == excepted {
			continue
		}

		if len(nodes) <= k {
			result = append(result, nodes...)
		} else {
			selectedIndices := make(map[int]struct{})
			count := 0
			var selectedNodes []common.Address
			for len(selectedNodes) < k && count < len(nodes) {
				index := r.Intn(len(nodes))
				if _, ok := selectedIndices[index]; !ok {
					count++
					_, connected := broadcaster.FindPeer(nodes[index])
					if connected {
						selectedIndices[index] = struct{}{}
						selectedNodes = append(selectedNodes, nodes[index])
					}
				}
			}

			if len(selectedNodes) > 0 {
				result = append(result, selectedNodes...)
			}
		}

	}

	return result
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

func AssignClusters(h uint64, nextEpochHeight uint64, latencyMat map[common.Address][]uint8, k int) (*Clusters, error) {
	nodes := fromLatencyMat(latencyMat)
	km := kmeans.New()
	cstrs, err := km.Partition(nodes, k, KmeansClusterSeed)
	if err != nil {
		return nil, err
	}

	result := make([][]common.Address, k)
	for i, cluster := range cstrs {
		// Collect addresses from cluster observations
		var addresses []common.Address
		for _, obs := range cluster.Observations {
			addresses = append(addresses, obs.(*node).address)
		}

		// Sort addresses lexicographically to make it deterministic.
		sort.Slice(addresses, func(a, b int) bool {
			return bytes.Compare(addresses[a][:], addresses[b][:]) < 0
		})

		result[i] = addresses
	}

	return NewCluster(h, nextEpochHeight, result), nil
}
