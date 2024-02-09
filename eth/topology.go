package eth

import (
	"math"

	"github.com/autonity/autonity/p2p/enode"
)

type networkTopology struct {
	diameter uint
	minNodes int
}

func NewGraphTopology(diameter uint, minNodes int) networkTopology {
	// Only diameter = 2, to support diameter > 2, ComputeBase function has to be modified
	if diameter != 2 {
		panic("diameter value must be 2")
	}
	return networkTopology{
		diameter: diameter,
		minNodes: minNodes,
	}
}

func (g *networkTopology) SetDiameter(d uint) {
	if d != 2 {
		panic("diameter value must be 2")
	}
	g.diameter = d
}

func (g *networkTopology) SetMinNodes(n int) {
	g.minNodes = n
}

func (g *networkTopology) computeSquareRoot(n uint) uint {
	return uint(math.Ceil(math.Sqrt(float64(n))))
}

// Returns the number of matching digits in i and j for g.diameter least significant digits
// Both i and j are considered to be in b-base number system
func (g *networkTopology) countMatchingDigits(i, j, b uint) uint {
	var count uint
	digitCount := g.diameter
	for digitCount > 0 {
		if i%b == j%b {
			count++
		}
		i /= b
		j /= b
		digitCount--
	}
	return count
}

// compute b such that b^d >= n and (b-1)^d < n where d = g.diameter
// for now only g.diameter = 2 is supported
func (g *networkTopology) ComputeBase(n uint) uint {
	return g.computeSquareRoot(n)
}

// Returns the list of adjacentNodes to connect with localNode. Given that the order of the input array nodes is same
// for everyone, connecting to only adjacentNodes will create a connected graph with diameter = g.diameter
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
	b := g.ComputeBase(uint(len(nodes)))
	connections := make([]*enode.Node, 0, len(nodes))
	for i, node := range nodes {
		if g.countMatchingDigits(uint(i), uint(myIdx), b) == 1 {
			connections = append(connections, node)
		}
	}
	return connections
}
