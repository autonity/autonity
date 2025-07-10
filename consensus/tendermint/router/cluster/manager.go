package cluster

import (
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/log"
)

type Manager struct {
	clusters Clusters
	sync.RWMutex
}

func New(committee []common.Address, latencyMap map[common.Address]uint, self common.Address) (Clusters, error) {
	numClusters := int(math.Floor(math.Sqrt(float64(len(committee)))))
	c, err := createClusters(committee, latencyMap, self, numClusters)
	if err != nil {
		return Clusters{}, err
	}
	c.computeLatencyBuckets()

	return c, nil
}

func (m *Manager) Clusters() Clusters {
	m.RLock()
	defer m.RUnlock()
	return m.clusters
}

func (m *Manager) UpdateClusters(clusters Clusters) {
	m.Lock()
	defer m.Unlock()
	m.clusters = clusters
	var sb strings.Builder
	sb.WriteString("Updating cluster, new cluster view: [")
	for i, cv := range m.clusters.Base() {
		if i > 0 {
			sb.WriteString("; ")
		}
		sb.WriteString(fmt.Sprintf("\nC%d:[", i))
		for j, m := range cv {
			if j > 0 {
				sb.WriteString(",")
			}
			addr := m.Addr.Hex()
			sb.WriteString(fmt.Sprintf("%s:%d", addr, m.Lat))
		}
		sb.WriteString("]\n")
	}
	sb.WriteString("]")
	log.Debug(sb.String())
}
