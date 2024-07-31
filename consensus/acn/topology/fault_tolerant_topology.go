package topology

import (
	"math/big"

	"github.com/autonity/autonity/consensus/tendermint/bft"
)

type faultTolerantTopology struct {
	*topologyBase
}

func NewFaultTolerantTopology(totalNodes, degreeCount, minNodes int) (Topology, error) {
	if degreeCount >= totalNodes {
		return &faultTolerantTopology{}, errDegreeCount
	}
	network := &faultTolerantTopology{&topologyBase{}}
	network.init(totalNodes, degreeCount, minNodes, network.boundaryNodes, network.makeNetworkValid, network.minDegree)
	return network, nil
}

func (g *faultTolerantTopology) minDegree() int {
	return (g.nodeCount-1)/3 + 1
}

func (g *faultTolerantTopology) makeNetworkValid() {
	g.checkDegreeLowerBound()
	for !g.isNetworkValid() {
		g.degreeCount++
	}
}

func (g *faultTolerantTopology) isNetworkValid() bool {
	if g.nodeCount == 1 {
		return true
	}
	if ((g.nodeCount - g.degreeCount - 1) & 1) == 1 { // (a%2) = (a&1)
		// to make the graph symmetric
		return false
	}

	return bft.F(big.NewInt(int64(g.nodeCount))).Int64() < int64(g.degreeCount)
}

func (g *faultTolerantTopology) boundaryNodes(nodeIndex int) (firstActiveNode, lastActiveNode int) {
	// Consider the nodes are arranged in a circular manner where the node `nodeIndex` is at the top.
	// The node `nodeIndex` has `inactiveNodeCount` inactive nodes in its both sides of the circle.
	// The active nodes are at the bottom side of the circle.
	inactiveNodeCount := (g.nodeCount - 1 - g.degreeCount) / 2
	// The first active node in the right side (node index increaments) of `nodeIndex` is `firstActiveNode`.
	firstActiveNode = (nodeIndex + inactiveNodeCount + 1) % g.nodeCount
	// The first active node in the left side (node index decreaments) of `nodeIndex` is `lastActiveNode`.
	lastActiveNode = (nodeIndex - inactiveNodeCount - 1 + g.nodeCount) % g.nodeCount
	return
}
