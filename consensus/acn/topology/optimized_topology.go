package topology

import (
	"math"
)

type optimizedTopology struct {
	*topologyBase
}

func NewOptimizedTopology(totalNodes, degreeCount, minNodes int) (Topology, error) {
	if degreeCount >= totalNodes {
		return &optimizedTopology{}, errDegreeCount
	}

	network := &optimizedTopology{&topologyBase{}}
	network.init(totalNodes, degreeCount, minNodes, network.boundaryNodes, network.makeNetworkValid, network.minDegree)
	return network, nil
}

func (g *optimizedTopology) minDegree() int {
	return min(g.nodeCount-1, int(math.Ceil(math.Sqrt(float64(g.nodeCount)))))
}

func (g *optimizedTopology) isNetworkValid() bool {
	if g.nodeCount == 1 {
		return true
	}
	return g.degreeCount >= g.minDegree()
}

func (g *optimizedTopology) makeNetworkValid() {
	g.checkDegreeLowerBound()
}

func (g *optimizedTopology) boundaryNodes(nodeIndex int) (firstActiveNode, lastActiveNode int) {
	if nodeIndex > g.degreeCount {
		// no edges coming out of these nodes
		return -1, -1
	}

	// node `i` connects to all nodes from `i*d+1` to `(i+1)*d` where `d = degree`
	firstActiveNode = (nodeIndex*g.degreeCount + 1) % g.nodeCount
	lastActiveNode = (firstActiveNode + g.degreeCount - 1) % g.nodeCount
	return
}
