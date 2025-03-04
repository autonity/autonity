package latency

import (
	"errors"
	"sync"

	"golang.org/x/exp/slices"

	"github.com/autonity/autonity/common"
)

var ErrMissingLatencyMeasurements = errors.New("missing latency measurements")

type latencyCache struct {
	mu               sync.RWMutex
	measurements     map[common.Address]uint8
	matrix           map[common.Address]map[common.Address]int16
	reportsThisEpoch uint64
}

func newLatencyCache() *latencyCache {
	return &latencyCache{
		measurements:     make(map[common.Address]uint8),
		matrix:           make(map[common.Address]map[common.Address]int16),
		reportsThisEpoch: 0,
	}
}

func (l *latencyCache) missingMeasurements(committee []common.Address) []common.Address {
	l.mu.RLock()
	defer l.mu.RUnlock()
	var missing []common.Address
	for _, validator := range committee {
		if _, ok := l.measurements[validator]; !ok {
			missing = append(missing, validator)
		}
	}
	return missing
}

func (l *latencyCache) insertMatrixLine(
	reporter common.Address,
	validators []common.Address,
	latencies []uint8,
) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, ok := l.matrix[reporter]; !ok {
		l.matrix[reporter] = make(map[common.Address]int16)
	}
	for i, validator := range validators {
		if validator == reporter {
			l.matrix[reporter][validator] = 0
		} else if latencies[i] == 0 {
			l.matrix[reporter][validator] = -1
		} else {
			l.matrix[reporter][validator] = int16(latencies[i])
		}
	}
	l.reportsThisEpoch++
}

func (l *latencyCache) reportsInEpoch() uint64 {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.reportsThisEpoch
}

func (l *latencyCache) markNewEpoch() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.reportsThisEpoch = 0
}

func (l *latencyCache) insertMeasurements(validators []common.Address, latencies map[common.Address]uint8) {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, validator := range validators {
		l.measurements[validator] = latencies[validator]
	}
}

func (l *latencyCache) latencyView(validators []common.Address) (map[common.Address]uint8, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	view := make(map[common.Address]uint8)
	for _, validator := range validators {
		if _, ok := l.measurements[validator]; !ok {
			return nil, ErrMissingLatencyMeasurements
		}
		view[validator] = l.measurements[validator]
	}
	return view, nil
}

func (l *latencyCache) readMatrix(validators []common.Address) map[common.Address][]uint8 {
	result := make(map[common.Address][]uint8)

	for _, va := range validators {
		for i, vb := range validators {
			if i == 0 {
				result[va] = make([]uint8, len(validators))
			}
			if val, ok := l.matrix[va][vb]; ok && val >= 0 {
				result[va][i] = uint8(val)
			} else {
				// check if we have the opposite direction
				if opposite, ok := l.matrix[vb][va]; ok && opposite >= 0 {
					result[va][i] = uint8(opposite)
				} else {
					// default to largest latency
					result[va][i] = ^uint8(0)
				}
			}
		}
	}
	return result
}

func (l *latencyCache) readMatrixWithOutliers(validators []common.Address) (map[common.Address][]uint8, []common.Address) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	matrix := make(map[common.Address][]int16)
	for _, va := range validators {
		for i, vb := range validators {
			if i == 0 {
				matrix[va] = make([]int16, len(validators))
			}
			if val, ok := l.matrix[va][vb]; ok && val >= 0 {
				matrix[va][i] = val
			} else {
				// check if we have the opposite direction
				if opposite, ok := l.matrix[vb][va]; ok && opposite >= 0 {
					matrix[va][i] = opposite
				} else {
					// unknown measurement
					matrix[va][i] = -1
				}
			}
		}
	}
	outliers := make([]common.Address, 0)

	for validator, latencies := range matrix {
		if countMax(latencies) >= 2*len(validators)/3 {
			outliers = append(outliers, validator)
		}
	}
	validatorSubset := setExclusion(validators, outliers)
	return l.readMatrix(validatorSubset), outliers
}

func countMax(latencies []int16) int {
	var maximum int
	for _, latency := range latencies {
		if latency == -1 {
			maximum++
		}
	}
	return maximum
}

func setExclusion(committee []common.Address, outliers []common.Address) []common.Address {
	committeeSubset := make([]common.Address, len(committee)-len(outliers))
	i := 0
	for _, validator := range committee {
		if !contains(outliers, validator) {
			committeeSubset[i] = validator
			i++
		}
	}
	return committeeSubset
}

func contains(validators []common.Address, validator common.Address) bool {
	for _, v := range validators {
		if v == validator {
			return true
		}
	}
	return false
}

type clusterCache struct {
	mu       sync.RWMutex
	clusters map[uint64]*Clusters
}

func newClusterCache() *clusterCache {
	return &clusterCache{
		clusters: make(map[uint64]*Clusters),
	}
}

func (c *clusterCache) blocks() []uint64 {
	blocks := make([]uint64, len(c.clusters))
	i := 0
	for block := range c.clusters {
		blocks[i] = block
		i++
	}
	slices.Sort(blocks)
	return blocks
}

func (c *clusterCache) insertClustering(block uint64, cluster *Clusters) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.clusters[block] = cluster
}

// clustersAt returns the clusters at a given block height, which will be the clusters at
// the last block height before the given block.
func (c *clusterCache) clustersAt(block uint64) (*Clusters, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	keys := c.blocks()
	for i := len(keys) - 1; i >= 0; i-- {
		if keys[i] < block {
			return c.clusters[keys[i]], true
		}
	}
	return nil, false
}

func (c *clusterCache) pruneTo(block uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	keys := c.blocks()
	for _, key := range keys {
		if key < block {
			delete(c.clusters, key)
		}
	}
}
