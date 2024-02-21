package eth

import (
	"math"

	"github.com/autonity/autonity/p2p/enode"
)

const (
	// max degree allowed for network in the execution layer
	MaxDegree = 25
	// if the network size exceeds MaxGraphSize, we divide the network in smaller sub-network of size MaxGraphSize
	MaxGraphSize = 64
)

type networkTopology struct {
	diameter int
	minNodes int
}

func NewGraphTopology(diameter, minNodes int) networkTopology {
	// Only diameter = 2, to support diameter > 2, ComputeBase and adjacentNodesIdx function has to be modified
	// if diameter != 2 {
	// 	panic("diameter value must be 2")
	// }
	return networkTopology{
		diameter: diameter,
		minNodes: minNodes,
	}
}

func (g *networkTopology) SetDiameter(d int) {
	if d != 2 {
		panic("diameter value must be 2")
	}
	g.diameter = d
}

func (g *networkTopology) SetMinNodes(n int) {
	g.minNodes = n
}

func (g *networkTopology) computeSquareRoot(n int) int {
	return int(math.Ceil(math.Sqrt(float64(n))))
}

// compute b such that b^d >= n and (b-1)^d < n where d = g.diameter
// for now only g.diameter = 2 is supported
func (g *networkTopology) ComputeBase(n int) int {
	return g.computeSquareRoot(n)
}

// maximum number of degree of each node = d*(b-1)^(d-1), where b = base, d = diameter
func (g *networkTopology) MaxDegree(totalNodeCount int) int {
	b := g.ComputeBase(totalNodeCount)
	d := g.diameter
	return d * int(math.Pow(float64(b-1), float64(d-1)))
}

// it constructs array of the edges by keeping one digit of myIdx fix and changing all the rest
func (g *networkTopology) edges(myIdx, totalNodes int) []int {
	b := g.ComputeBase(totalNodes)

	// the following part supports only g.diameter = 2 for ease of coding,
	// it needs to change to support g.diameter > 2 (which may not be necessary)
	lsb := myIdx % b
	msb := (myIdx / b) * b
	adjacentNodes := make([]int, 0, 2*(b-1))
	// fix msb and change lsb
	for i := 0; i < b; i++ {
		if i != lsb && msb+i < totalNodes {
			adjacentNodes = append(adjacentNodes, msb+i)
		}
	}
	// fix lsb and change msg
	for i := 0; i < b*b; i += b {
		if i != msb && i+lsb < totalNodes {
			adjacentNodes = append(adjacentNodes, i+lsb)
		}
	}
	return adjacentNodes
}

func (g *networkTopology) componentSize(componentEndIdx []int, componentIdx int) int {
	if componentIdx > 0 {
		return componentEndIdx[componentIdx] - componentEndIdx[componentIdx-1]
	}
	return componentEndIdx[componentIdx]
}

func (g *networkTopology) componentRelativeIdx(componentEndIdx, actualIdx int) int {
	return componentEndIdx - 1 - actualIdx
}

func (g *networkTopology) idxFromRelativeIdx(componentEndIdx, relativeIdx int) int {
	return componentEndIdx - 1 - relativeIdx
}

func (g *networkTopology) componentCount(nodeCount int) int {
	return (nodeCount + MaxGraphSize - 1) / MaxGraphSize // components = math.Ceil(totalNodes/MaxGraphSize)
}

func (g *networkTopology) componentEndIdx(totalNodes int) []int {
	components := g.componentCount(totalNodes)
	componentEndIdx := make([]int, components)
	for i := 0; i < components-1; i++ {
		componentEndIdx[i] = (i + 1) * MaxGraphSize
	}
	componentEndIdx[components-1] = totalNodes
	if components > 1 && g.componentSize(componentEndIdx, components-1) < (MaxGraphSize+1)/2 {
		componentEndIdx[components-2] -= MaxGraphSize / 2
	}
	return componentEndIdx
}

func (g *networkTopology) componentIdx(componentEndIdx []int, nodeIdx int) int {
	componentIdx := 0
	// scope to improve: do a binary search if len(componentEndIdx) is too big
	for nodeIdx >= componentEndIdx[componentIdx] {
		componentIdx++
	}
	return componentIdx
}

func (g *networkTopology) adjacentNodesIdx(myIdx, totalNodes int) []int {
	if totalNodes <= MaxGraphSize {
		return g.edges(myIdx, totalNodes)
	}
	components := g.componentCount(totalNodes)
	componentEndIdx := g.componentEndIdx(totalNodes)
	myComponentIdx := g.componentIdx(componentEndIdx, myIdx)
	relativeIdx := g.componentRelativeIdx(componentEndIdx[myComponentIdx], myIdx)
	myComponentSize := g.componentSize(componentEndIdx, myComponentIdx)
	connections := g.edges(relativeIdx, myComponentSize)
	for i := 0; i < len(connections); i++ {
		connections[i] = g.idxFromRelativeIdx(componentEndIdx[myComponentIdx], connections[i])
	}
	componentConnections := g.adjacentNodesIdx(myComponentIdx, components)
	for _, componentIdx := range componentConnections {
		componentSize := g.componentSize(componentEndIdx, componentIdx)
		peerRelativeIdx := relativeIdx
		if myComponentSize >= componentSize {
			if peerRelativeIdx >= componentSize {
				peerRelativeIdx -= componentSize
			}
			peerIdx := g.idxFromRelativeIdx(componentEndIdx[componentIdx], peerRelativeIdx)
			connections = append(connections, peerIdx)
		} else {
			factor := (componentSize + myComponentSize - 1) / myComponentSize // factor = math.Ceil(componentSize/myComponentSize)
			for factor > 0 && peerRelativeIdx < componentSize {
				peerIdx := g.idxFromRelativeIdx(componentEndIdx[componentIdx], peerRelativeIdx)
				connections = append(connections, peerIdx)
				peerRelativeIdx += myComponentSize
				factor--
			}
		}
	}
	return connections
}

// the input array (nodes []*enode.Node) must be same for everyone in order to create a connected graph
// Returns the list of adjacentNodes to connect with localNode. Given that the order of the input array nodes is same
// for everyone, connecting to only adjacentNodes will create a connected graph with diameter <= 4
func (g *networkTopology) RequestSubset(nodes []*enode.Node, localNode *enode.LocalNode) []*enode.Node {
	if len(nodes) < g.minNodes {
		// connect to all nodes
		return nodes
	}
	myIdx := -1
	for i, node := range nodes {
		if node.ID() == localNode.ID() {
			myIdx = i
			break
		}
	}
	// If the node is not in committee, it has all slots available, so connect to all committee nodes
	if myIdx == -1 {
		return nodes
	}
	// Graph construction mechanism for large graph
	// each node is represented as a number in b-base number system with exactly d digits
	// which requires len(nodes) <= b^d. For example, if b = 2 and d = 2
	// then we can support 4 nodes numbering 00, 01, 10, 11 (in binary).
	// Two nodes i and j is connected if they have exactly 1 digit common
	// So in above example, 00 is connected with 01 and 10; 01 is connected with 00 and 11;
	// 10 is connected with 00 and 11; 11 is connected with 01 and 10.

	// for now only diameter = 2 is supported
	// b is chosen with the following property len(nodes) <= b*b and len(nodes) > (b-1)*(b-1)
	adjacentNodes := g.adjacentNodesIdx(myIdx, len(nodes))
	connections := make([]*enode.Node, 0, len(adjacentNodes))
	for _, idx := range adjacentNodes {
		connections = append(connections, nodes[idx])
	}
	return connections
}
