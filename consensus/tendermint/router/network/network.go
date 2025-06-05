package network

import (
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/log"
)

type Network struct {
	clusters Clusters
	sync.RWMutex
}

func New(
	committee []common.Address,
	latencyMap map[common.Address]uint,
	self common.Address,
) (Clusters, error) {
	numClusters := int(math.Floor(math.Sqrt(float64(len(committee)))))
	c, err := createClusters(committee, latencyMap, self, numClusters)
	if err != nil {
		return Clusters{}, err
	}

	remoteBuckets := c.ComputeLatencyBuckets()

	PrintLatencyBuckets(remoteBuckets, c.bucketSize, c.minLatency)

	c.AssignRemoteNodes(remoteBuckets, numClusters)

	return c, nil
}

func (n *Network) Clusters() Clusters {
	n.RLock()
	defer n.RUnlock()
	return n.clusters
}

func (n *Network) UpdateClusters(clusters Clusters) {
	n.Lock()
	defer n.Unlock()
	var sb strings.Builder
	sb.WriteString("Updating cluster, new cluster view: [")
	for i, cv := range n.clusters.Base() {
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
	log.Info(sb.String())
	n.clusters = clusters
}

func PrintLatencyBuckets(remoteBuckets [][]Node, bucketSize float64, minLatency uint) {
	var sb strings.Builder
	sb.WriteString("\nLatency Buckets for Clusters:\n")

	// Log remote buckets
	sb.WriteString("Remote Buckets:\n")
	for bucketIdx, nodes := range remoteBuckets {
		lowerLat := uint(float64(bucketIdx)*bucketSize) + minLatency
		upperLat := uint(float64(bucketIdx+1)*bucketSize) + minLatency
		sb.WriteString(fmt.Sprintf("  Bucket #%d (Latency %d-%d ms): %d nodes\n", bucketIdx, lowerLat, upperLat, len(nodes)))
		if len(nodes) == 0 {
			sb.WriteString("    [Empty]\n")
			continue
		}
		for _, node := range nodes {
			sb.WriteString(fmt.Sprintf("    Node: %s, Latency: %d ms, ClusterID: %d\n", node.Addr.Hex(), node.Lat, node.ClusterID))
		}
	}

	log.Info(sb.String())
}
