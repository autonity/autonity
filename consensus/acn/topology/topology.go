package topology

import (
	"errors"
	"math/big"

	"github.com/autonity/autonity/consensus/tendermint/bft"
)

var errDegreeCount = errors.New("degree count is greater than network size")

type networkTopology struct {
	nodeCount, degreeCount, minNodes int
}

func NewGraphTopology(totalNodes, degreeCount, minNodes int) (networkTopology, error) {
	if degreeCount >= totalNodes {
		return networkTopology{}, errDegreeCount
	}
	network := networkTopology{
		nodeCount:   totalNodes,
		degreeCount: degreeCount,
		minNodes:    minNodes,
	}

	network.makeNetworkValid()
	return network, nil
}

func (g *networkTopology) makeNetworkValid() {
	if g.nodeCount < 2 {
		return
	}
	if !g.isNetworkValid() {
		g.degreeCount = (g.nodeCount-1)/3 + 1
	}

	for !g.isNetworkValid() {
		g.degreeCount++
	}
}

func (g *networkTopology) isNetworkValid() bool {
	if ((g.nodeCount - g.degreeCount - 1) & 1) == 1 { // (a%2) = (a&1)
		// to make the graph symmetric
		return false
	}

	return bft.F(big.NewInt(int64(g.nodeCount))).Int64() < int64(g.degreeCount)
}

func (g *networkTopology) SetDegreeCount(degree int) error {
	if degree >= g.nodeCount {
		return errDegreeCount
	}
	g.degreeCount = degree
	g.makeNetworkValid()
	return nil
}

func (g *networkTopology) boundaryNodes(nodeIndex int) (rightActiveNode, leftActiveNode int) {
	// Consider the nodes are arranged in a circular manner where the node `nodeIndex` is at the top.
	// The node `nodeIndex` has `inactiveNodeCount` inactive nodes in its both sides of the circle.
	// The active nodes are at the bottom side of the circle.
	inactiveNodeCount := (g.nodeCount - 1 - g.degreeCount) / 2
	// The first active node in the right side (node index increaments) of `nodeIndex` is `rightActiveNode`.
	rightActiveNode = (nodeIndex + inactiveNodeCount + 1) % g.nodeCount
	// The first active node in the left side (node index decreaments) of `nodeIndex` is `leftActiveNode`.
	leftActiveNode = (nodeIndex - inactiveNodeCount - 1 + g.nodeCount) % g.nodeCount
	return
}

func (g *networkTopology) isAdjacent(u, v int) bool {
	rightActiveNode, leftActiveNode := g.boundaryNodes(u)
	if rightActiveNode <= leftActiveNode {
		// there is no circular crossing, v is in the middle somewhere
		return v >= rightActiveNode && v <= leftActiveNode
	}
	return v >= rightActiveNode || v <= leftActiveNode
}

func (g *networkTopology) RequestSubset(myIndex int) []int {
	if g.nodeCount < 2 {
		return nil
	}
	edges := make([]int, 0, g.nodeCount-1)
	if g.nodeCount <= g.minNodes {
		for i := 0; i < g.nodeCount; i++ {
			if i == myIndex {
				continue
			}
			edges = append(edges, i)
		}
		return edges
	}

	rightActiveNode, leftActiveNode := g.boundaryNodes(myIndex)

	if rightActiveNode <= leftActiveNode {
		// will be true if there is no circulare crossing
		for node := rightActiveNode; node <= leftActiveNode; node++ {
			edges = append(edges, node)
		}
	} else {
		for node := 0; node <= leftActiveNode; node++ {
			edges = append(edges, node)
		}
		for node := rightActiveNode; node < g.nodeCount; node++ {
			edges = append(edges, node)
		}
	}
	return edges
}
