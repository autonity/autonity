package eth

import (
	"bytes"

	"github.com/autonity/autonity/p2p/enode"
)

type NetworkTopology struct {
	diameter uint
	minNodes int
}

func NewGraphTopology(diameter uint, minNodes int) *NetworkTopology {
	return &NetworkTopology{
		diameter: diameter,
		minNodes: minNodes,
	}
}

func (g *NetworkTopology) SetDiameter(d uint) {
	g.diameter = d
}

func (g *NetworkTopology) SetMinNodes(n int) {
	g.minNodes = n
}

func (g *NetworkTopology) computeSquareRoot(n uint) uint {
	if n == 0 {
		return 0
	}
	var low uint = 1
	high := n
	for low < high {
		mid := (low + high) / 2
		if mid >= n/mid {
			high = mid
		} else {
			low = mid + 1
		}
	}
	return high
}

func (g *NetworkTopology) countMatchingDigits(a, b, base uint) uint {
	var count uint
	for g.diameter > 0 {
		if a%base == b%base {
			count++
		}
		a /= base
		b /= base
		g.diameter--
	}
	return count
}

func (g *NetworkTopology) RequestSubset(nodes []*enode.Node, localNode *enode.LocalNode) []*enode.Node {
	if len(nodes) < g.minNodes {
		return nodes
	}
	myIdx := -1
	for i, node := range nodes {
		if bytes.Equal(node.ID().Bytes(), localNode.ID().Bytes()) {
			myIdx = i
			break
		}
	}
	if myIdx == -1 {
		return nil
	}
	b := g.computeSquareRoot(uint(len(nodes)))
	connections := make([]*enode.Node, 0, len(nodes))
	for i, node := range nodes {
		if g.countMatchingDigits(uint(i), uint(myIdx), b) == 1 {
			connections = append(connections, node)
		}
	}
	return connections
}
