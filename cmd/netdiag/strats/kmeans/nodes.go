package kmeans

import (
	"math"
	"time"

	"github.com/parallelo-ai/kmeans"
)

type Node struct {
	Id              int
	AssignedCluster int
	LatencyView     []float64
}

func (n *Node) Coordinates() kmeans.Coordinates {
	return n.LatencyView
}

func (n *Node) Distance(p2 kmeans.Coordinates) float64 {
	var r float64
	for i, v := range n.LatencyView {
		r += math.Pow(v-p2[i], 2)
	}
	return r
}

func fromLatencyMatrix(latencyMatrix [][]time.Duration) []*Node {
	nodes := make([]*Node, len(latencyMatrix))
	for i, row := range latencyMatrix {
		nodes[i] = &Node{
			Id:          i,
			LatencyView: make([]float64, len(row)),
		}
		for j, latency := range row {
			nodes[i].LatencyView[j] = float64(latency.Nanoseconds()) / 1e6
		}
	}
	return nodes
}
