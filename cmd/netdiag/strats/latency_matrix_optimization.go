package strats

import (
	"math"
	"time"
)

type LatencyMatrixOptimize struct {
	BaseStrategy
	graph Graph

	peerSetUpperBound int
}

func init() {
	registerStrategy("Latency Matrix Optimization - Upper Bound min(sqrt(peers), 10%)", func(base BaseStrategy) Strategy {
		return createLatencyMatrixOpimize(base, 10)
	})
	registerStrategy("Latency Matrix Optimization - Upper Bound min(sqrt(peers), 20%)", func(base BaseStrategy) Strategy {
		return createLatencyMatrixOpimize(base, 20)
	})
}

func createLatencyMatrixOpimize(base BaseStrategy, peerSetUpperBound int) *LatencyMatrixOptimize {
	strategy := &LatencyMatrixOptimize{base, Graph{}, peerSetUpperBound}
	go strategy.start()
	return strategy
}

func (o *LatencyMatrixOptimize) start() {
	for {
		if o.isLatencyMatrixReady() {
			break
		}
	}
	o.constructGraph(len(o.State.LatencyMatrix))
}

func (o *LatencyMatrixOptimize) isLatencyMatrixReady() bool {
	for id, array := range o.State.LatencyMatrix {
		for peer, latency := range array {
			if id == peer {
				continue
			}
			if latency == 0 {
				return false
			}
		}
	}
	return true
}

func (o *LatencyMatrixOptimize) constructGraph(peers int) {
	o.graph.peers = peers
	timeMatrix := make([][]time.Duration, peers)
	for i := 0; i < peers; i++ {
		timeMatrix[i] = make([]time.Duration, peers)
	}
	o.graph.rootedConnection = make([][]int, peers)
	maxConnections := int(math.Ceil(float64(o.peerSetUpperBound * peers / 100)))
	minConnections := int(math.Ceil(math.Sqrt(float64(peers))))
	maxConnections = max(maxConnections, minConnections)
	for peer := 0; peer < peers; peer++ {
	}
	o.graph.initiated = true
}

func (o *LatencyMatrixOptimize) Execute(packetId uint64, data []byte, maxPeers int) error {
	return nil
}

func (o *LatencyMatrixOptimize) HandlePacket(packetId uint64, hop uint8, originalSender uint64, maxPeers uint64, data []byte, partial bool, seqNum, total uint16) error {
	return nil
}

type Graph struct {
	initiated bool
	peers     int
	// `rootedConnection[root]` contains the connection array of the node when
	// the original sender is `root`
	rootedConnection [][]int
}
