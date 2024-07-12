package topology

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

func (g *networkTopology) RequestSubset(nodeCount, myIndex int) []int {
	edges := make([]int, 0, nodeCount)
	if nodeCount <= g.minNodes {
		for i := 0; i < nodeCount; i++ {
			if i == myIndex {
				continue
			}
			edges = append(edges, i)
		}
		return edges
	}

	// TODO (tariq): need to consider voting power of each node
	return nil
}
