package latency

import (
	"errors"
	"sync"

	"github.com/autonity/autonity/common"
)

var ErrMissingLatencyMeasurements = errors.New("missing latency measurements")

type latencyCache struct {
	mu           sync.RWMutex
	measurements map[common.Address]uint8
	matrix       map[common.Address]map[common.Address]uint8
}

func newLatencyCache() *latencyCache {
	return &latencyCache{
		measurements: make(map[common.Address]uint8),
		matrix:       make(map[common.Address]map[common.Address]uint8),
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

func (l *latencyCache) updateMatrix(validators []common.Address, latMat [][]uint8) {
	l.mu.Lock()
	defer l.mu.Unlock()
	for i := range validators {
		for j := range validators {
			if _, ok := l.matrix[validators[i]]; !ok {
				l.matrix[validators[i]] = make(map[common.Address]uint8)
			}
			if i == j {
				l.matrix[validators[i]][validators[j]] = 0
			} else if latMat[i][j] == 0 {
				l.matrix[validators[i]][validators[j]] = ^uint8(0)
			} else {
				l.matrix[validators[i]][validators[j]] = latMat[i][j]
			}
		}
	}
}

func (l *latencyCache) insertMatrixLine(reporter common.Address, validators []common.Address, latencies []uint8) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, ok := l.matrix[reporter]; !ok {
		l.matrix[reporter] = make(map[common.Address]uint8)
	}
	for i, validator := range validators {
		if validator == reporter {
			l.matrix[reporter][validator] = 0
		} else if latencies[i] == 0 {
			l.matrix[reporter][validator] = ^uint8(0)
		} else {
			l.matrix[reporter][validator] = latencies[i]
		}
	}
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
	l.mu.RLock()
	defer l.mu.RUnlock()

	for _, va := range validators {
		for i, vb := range validators {
			if i == 0 {
				result[va] = make([]uint8, len(validators))
			}
			if val, ok := l.matrix[va][vb]; ok {
				result[va][i] = val
			} else {
				result[va][i] = ^uint8(0)
			}
		}
	}
	return result
}
