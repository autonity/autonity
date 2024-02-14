package eth

import (
	"math"

	"github.com/autonity/autonity/p2p/enode"
)

type networkTopology struct {
	diameter int
	minNodes int
}

func NewGraphTopology(diameter int, minNodes int) networkTopology {
	// Only diameter = 2, to support diameter > 2, ComputeBase and adjacentNodesIdx function has to be modified
	if diameter != 2 {
		panic("diameter value must be 2")
	}
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

// Returns the number of matching digits in i and j for g.diameter least significant digits
// Both i and j are considered to be in b-base number system
func (g *networkTopology) countMatchingDigits(i, j, b int) int {
	count := 0
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
func (g *networkTopology) ComputeBase(n int) int {
	return g.computeSquareRoot(n)
}

// maximum number of degree of each node = d*(b-1)^(d-1), where b = base, d = diameter
func (g *networkTopology) MaxDegree(totalNodeCount int) int {
	b := g.ComputeBase(totalNodeCount)
	d := g.diameter
	return d * int(math.Pow(float64(b-1), float64(d-1)))
}

// maximum number of degree for each node = d*(b-1)^(d-1) where d = g.diameter
// it constructs array of them by keeping one digit of myIdx fix and changing all the rest
func (g *networkTopology) adjacentNodesIdx(myIdx, totalNodes int) []int {
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
	adjacentNodes := g.adjacentNodesIdx(myIdx, len(nodes))
	connections := make([]*enode.Node, 0, len(adjacentNodes))
	for _, idx := range adjacentNodes {
		connections = append(connections, nodes[idx])
	}
	return connections
}
