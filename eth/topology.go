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
	minNodes int
}

func NewGraphTopology(minNodes int) networkTopology {
	return networkTopology{
		minNodes: minNodes,
	}
}

func (g *networkTopology) SetMinNodes(n int) {
	g.minNodes = n
}

func (g *networkTopology) computeSquareRoot(n int) int {
	return int(math.Ceil(math.Sqrt(float64(n))))
}

// base = b such that b*b >= n
func (g *networkTopology) ComputeBase(n int) int {
	return g.computeSquareRoot(n)
}

// Construction mechanism: each node is represented as a number in b-base number system with 2 digits, i.e. each node = {i,j} where 0 <= i,j < b
// b is chosen such that totalNodes <= b*b. Two nodes {a,b} and {c,d} are connected if (a = c and b != d) or (a != c and b = d)
// It constructs array of the edges by keeping one digit of myIndex fix and changing all the rest
func (g *networkTopology) edges(myIndex, totalNodes int) []int {
	b := g.ComputeBase(totalNodes)

	lsb := myIndex % b
	msb := (myIndex / b) * b
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

func (g *networkTopology) componentSize(componentEndIndex []int, componentIndex int) int {
	if componentIndex > 0 {
		return componentEndIndex[componentIndex] - componentEndIndex[componentIndex-1]
	}
	return componentEndIndex[componentIndex]
}

// The nodes in a single component are numbered from 0 to (componentSize-1) from big to small
func (g *networkTopology) componentRelativeIndex(componentEndIndex, actualIndex int) int {
	return componentEndIndex - 1 - actualIndex
}

func (g *networkTopology) indexFromRelativeIndex(componentEndIndex, relativeIndex int) int {
	return componentEndIndex - 1 - relativeIndex
}

func (g *networkTopology) componentCount(nodeCount int) int {
	return (nodeCount + MaxGraphSize - 1) / MaxGraphSize // components = math.Ceil(totalNodes/MaxGraphSize)
}

// divides the graph with totalNodes in one or more components where each component size <= MaxGraphSize
// returns the array of (end index + 1) of the components
func (g *networkTopology) componentEndIndex(totalNodes int) []int {
	components := g.componentCount(totalNodes)
	componentEndIndex := make([]int, components)
	for i := 0; i < components-1; i++ {
		componentEndIndex[i] = (i + 1) * MaxGraphSize
	}
	componentEndIndex[components-1] = totalNodes
	// For any two components, a and b, 2*size(a) >= size(b) is followed when creating the component
	// Otherwise the nodes from the smaller component will have too many edges
	if components > 1 && g.componentSize(componentEndIndex, components-1) < (MaxGraphSize+1)/2 {
		componentEndIndex[components-2] -= MaxGraphSize / 2
	}
	return componentEndIndex
}

// returns the index of the component which nodeIndex belongs to
func (g *networkTopology) componentIndex(componentEndIndex []int, nodeIndex int) int {
	low := 0
	high := len(componentEndIndex) - 1
	for low < high {
		mid := (low + high) >> 1
		if componentEndIndex[mid] > nodeIndex {
			high = mid
		} else {
			low = mid + 1
		}
	}
	return low
}

// If totalNodes <= MaxGraphSize, it uses 'The Construction Mechanism' (via g.edges(int,int))
// Otherwise it divides the graph in several components where each component has nodes <= MaxGraphSize.
// So each component can apply 'The Construction Mechanism' to connect the nodes inside this component.
// After that, each component can be considered a single node
// If the number of components <= MaxGraphSize, 'The Construction Mechanism' is applied to connect the components,
// otherwise they are divided into more components recursively and eventually will be connected.
// Consider two components a and b such that they are connected directly by an edge where size(a) >= size(b)
// But to connect the component, we need to create edge between some nodes from component a and b
// Lets number all the nodes in component a from 0 to (a-1) and all the nodes in component b from 0 to (b-1)
// Each node, c from component a will be connected to another node, d from compoent b such that d = c % size(b).
// For example, for two compoents a and b with size = 3, all the ndoes in each component will be numbered from 0 to 2
// Then node i = {0,1,2} from component a will be connected to node i from component b
// For two components a with size 3 and b with size 2, nodes will be numbered form 0 to 2 (component a) and from 0 to 1 (component b)
// Then node i = {0,1} from component a will be connected to node i from component b, and node = 2 from component a will be connected to node 0 from component b
func (g *networkTopology) adjacentNodesIndex(myIndex, totalNodes int) []int {
	if totalNodes <= MaxGraphSize {
		return g.edges(myIndex, totalNodes)
	}
	componentCount := g.componentCount(totalNodes)
	componentEndIndex := g.componentEndIndex(totalNodes)
	// Index of the component in which myIndex belongs to
	myComponentIndex := g.componentIndex(componentEndIndex, myIndex)
	relativeIndex := g.componentRelativeIndex(componentEndIndex[myComponentIndex], myIndex)
	myComponentSize := g.componentSize(componentEndIndex, myComponentIndex)
	connections := g.edges(relativeIndex, myComponentSize)
	for i := 0; i < len(connections); i++ {
		connections[i] = g.indexFromRelativeIndex(componentEndIndex[myComponentIndex], connections[i])
	}
	componentConnections := g.adjacentNodesIndex(myComponentIndex, componentCount)
	for _, componentIndex := range componentConnections {
		componentSize := g.componentSize(componentEndIndex, componentIndex)
		peerRelativeIndex := relativeIndex
		if myComponentSize >= componentSize {
			if peerRelativeIndex >= componentSize {
				// for a < 2*b and a >= b, (a % b) can be written as (a - b)
				// components are divided in g.componentEndIndex(int) in such a way that we have
				// peerRelativeIndex < 2*componentSize
				peerRelativeIndex -= componentSize
			}
			peerIndex := g.indexFromRelativeIndex(componentEndIndex[componentIndex], peerRelativeIndex)
			connections = append(connections, peerIndex)
		} else {
			factor := (componentSize + myComponentSize - 1) / myComponentSize // factor = math.Ceil(componentSize/myComponentSize)
			for factor > 0 && peerRelativeIndex < componentSize {
				peerIndex := g.indexFromRelativeIndex(componentEndIndex[componentIndex], peerRelativeIndex)
				connections = append(connections, peerIndex)
				peerRelativeIndex += myComponentSize
				factor--
			}
		}
	}
	return connections
}

func (g *networkTopology) MyIndex(nodes []*enode.Node, localNode *enode.LocalNode) int {
	for i, node := range nodes {
		if node.ID() == localNode.ID() {
			return i
		}
	}
	return -1
}

// the input array (nodes []*enode.Node) must be same for everyone in order to create a connected graph
// Returns the list of adjacentNodes to connect with localNode. Given that the order of the input array nodes is same
// for everyone, connecting to only adjacentNodes will create a connected graph with diameter <= 4
func (g *networkTopology) RequestSubset(nodes []*enode.Node, myIndex int) []*enode.Node {
	if len(nodes) < g.minNodes {
		// connect to all nodes
		return nodes
	}
	// If the node is not in committee, it has all slots available, so connect to all committee nodes
	if myIndex == -1 {
		return nodes
	}
	adjacentNodes := g.adjacentNodesIndex(myIndex, len(nodes))
	connections := make([]*enode.Node, 0, len(adjacentNodes))
	for _, index := range adjacentNodes {
		connections = append(connections, nodes[index])
	}
	return connections
}
