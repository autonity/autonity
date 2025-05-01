package latency

import (
	"bytes"
	"math"
	"math/rand"
	"sort"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/latency/kmeans"
)

// KmeansClusterSeed is the seed to init the initial center position for each cluster.
var KmeansClusterSeed = int64(12345)

type Clusters struct {
	curEpochHeight   uint64
	nextEpochHeight  uint64
	base             [][]common.Address
	addressToCluster map[common.Address]int
}

// DoTransition remove those nodes which are removed from cluster, add those new ones into an individual cluster, and rebuild the index.
func (c *Clusters) DoTransition(removed map[common.Address]struct{}, added []common.Address) {

	for i := 0; i < len(c.base); i++ {
		var newCluster []common.Address
		for _, addr := range c.base[i] {
			if _, ok := removed[addr]; !ok {
				newCluster = append(newCluster, addr)
			}
		}
		c.base[i] = newCluster
	}

	// check to remove empty clusters after the removal.
	var newBase [][]common.Address
	for _, cluster := range c.base {
		if len(cluster) > 0 {
			newBase = append(newBase, cluster)
		}
	}
	c.base = newBase

	// append the newly added ones into an individual cluster.
	if len(added) > 0 {
		c.base = append(c.base, added)
	}

	// rebuild the cluster.
	c.buildAddressIndex()
}

func NewCluster(curEpochHeight uint64, nextEpochHeight uint64, base [][]common.Address) *Clusters {
	clusters := &Clusters{
		base:            base,
		curEpochHeight:  curEpochHeight,
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
func (c *Clusters) selectK(k int, seed int64, excepted int, broadcaster consensus.Broadcaster) [][]common.Address {
	result := make([][]common.Address, len(c.base))
	r := rand.New(rand.NewSource(seed))

	maxTrials := k * 5
	for i, cluster := range c.base {
		if i == excepted {
			result[i] = []common.Address{}
			continue
		}

		var selectedCluster []common.Address
		if len(cluster) <= k {
			for _, addr := range cluster {
				if _, ok := broadcaster.FindPeer(addr); ok {
					selectedCluster = append(selectedCluster, addr)
				}
			}
		} else {
			selectedIndices := make(map[int]struct{})
			tries := 0
			for len(selectedCluster) < k && tries < maxTrials {
				index := r.Intn(len(cluster))
				if _, ok := selectedIndices[index]; !ok {
					if _, ok := broadcaster.FindPeer(cluster[index]); ok {
						selectedIndices[index] = struct{}{}
						selectedCluster = append(selectedCluster, cluster[index])
					}
					tries++
				}
			}
		}

		result[i] = selectedCluster
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

func AssignClusters(curEpochBlock uint64, nextEpochHeight uint64, committee []common.Address, latencyMat [][]uint8, k int) (*Clusters, error) {
	nodes := fromLatencyMat(committee, latencyMat)
	km := kmeans.New()
	kmClusters, err := km.Partition(nodes, k, KmeansClusterSeed)
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

	return NewCluster(curEpochBlock, nextEpochHeight, result), nil
}
