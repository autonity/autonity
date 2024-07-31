package topology

import (
	"errors"
)

var (
	errDegreeCount = errors.New("degree count is greater than network size")
	errConnection  = errors.New("degree count does not match with boundary nodes")
)

type Topology interface {
	RequestSubset(nodeIndex int) ([]int, error)
	DegreeCount() int
	SetDegree(degree int) error

	isAdjacent(u, v int) bool
	isNetworkValid() bool
}

type topologyBase struct {
	nodeCount        int
	degreeCount      int
	minNodes         int
	boundaryNodes    func(nodeIndex int) (firstActiveNode, lastActiveNode int)
	makeNetworkValid func()
	minDegree        func() int
}

func (g *topologyBase) init(
	nodeCount, degreeCount, minNodes int,
	boundaryNodes func(nodeIndex int) (firstActiveNode, lastActiveNode int),
	makeNetworkValid func(),
	minDegree func() int,
) {
	g.nodeCount = nodeCount
	g.degreeCount = degreeCount
	g.minNodes = minNodes
	g.boundaryNodes = boundaryNodes
	g.makeNetworkValid = makeNetworkValid
	g.minDegree = minDegree
	g.makeNetworkValid()
}

func (g *topologyBase) DegreeCount() int {
	return g.degreeCount
}

func (g *topologyBase) checkDegreeLowerBound() {
	if g.nodeCount < 2 {
		g.degreeCount = 0
		return
	}
	if g.nodeCount == 2 {
		g.degreeCount = 1
		return
	}
	g.degreeCount = max(g.degreeCount, g.minDegree())
}

func (g *topologyBase) isAdjacent(u, v int) bool {
	if max(u, v) >= g.nodeCount || min(u, v) < 0 || g.degreeCount == 0 {
		return false
	}

	firstActiveNode, lastActiveNode := g.boundaryNodes(u)
	if firstActiveNode <= lastActiveNode {
		return v >= firstActiveNode && v <= lastActiveNode
	}
	return v >= firstActiveNode || v <= lastActiveNode
}

func (g *topologyBase) SetDegree(degree int) error {
	if degree >= g.nodeCount {
		return errDegreeCount
	}
	g.degreeCount = degree
	g.makeNetworkValid()
	return nil
}

func (g *topologyBase) RequestSubset(nodeIndex int) ([]int, error) {
	if g.nodeCount < 2 {
		return nil, nil
	}
	if g.nodeCount <= g.minNodes {
		edges := make([]int, 0, g.nodeCount-1)
		for i := 0; i < g.nodeCount; i++ {
			if i == nodeIndex {
				continue
			}
			edges = append(edges, i)
		}
		return edges, nil
	}

	firstActiveNode, lastActiveNode := g.boundaryNodes(nodeIndex)
	if firstActiveNode == -1 {
		// no edges coming out of this node
		return nil, nil
	}

	edges := make([]int, 0, g.degreeCount)

	if firstActiveNode <= lastActiveNode {
		if g.degreeCount != lastActiveNode-firstActiveNode+1 {
			return nil, errConnection
		}
		for node := firstActiveNode; node <= lastActiveNode; node++ {
			if node == nodeIndex {
				continue
			}
			edges = append(edges, node)
		}
		return edges, nil
	}

	if g.nodeCount-(firstActiveNode-lastActiveNode-1) != g.degreeCount {
		return nil, errConnection
	}
	for node := 0; node <= lastActiveNode; node++ {
		if node == nodeIndex {
			continue
		}
		edges = append(edges, node)
	}
	for node := firstActiveNode; node < g.nodeCount; node++ {
		if node == nodeIndex {
			continue
		}
		edges = append(edges, node)
	}
	return edges, nil
}
