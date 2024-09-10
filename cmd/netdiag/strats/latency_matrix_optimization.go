package strats

import (
	"errors"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/autonity/autonity/log"
)

var (
	errGraphNotConstructed   = errors.New("graph for latrix matrix optimization is not ready yet")
	errGraphConstruction     = errors.New("invalid graph construction")
	errInvalidArgumentHop    = errors.New("hop is greater than 1")
	errInvalidLatencyMatrix  = errors.New("latency to self should be zero")
	ErrLatencyMatrixNotReady = errors.New("latency matrix not ready")
	errPeerNotFound          = errors.New("peer not found")
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
	graph := Graph{
		id:             int(base.State.Id),
		peerGraphReady: make([]bool, base.State.Peers),
	}
	strategy := &LatencyMatrixOptimize{base, graph, peerSetUpperBound}
	return strategy
}

func (l *LatencyMatrixOptimize) isLatencyMatrixReady() (bool, error) {
	for id, array := range l.State.LatencyMatrix {
		if len(array) != l.State.Peers {
			return false, nil
		}
		for peer, latency := range array {
			if id == peer {
				if latency != 0 {
					return false, errInvalidLatencyMatrix
				}
				continue
			}
			if latency == 0 {
				return false, nil
			}
		}
	}
	return true, nil
}

func (l *LatencyMatrixOptimize) GraphReadyForPeer(peerID int) {
	l.graph.peerGraphReady[peerID] = true
}

func (l *LatencyMatrixOptimize) IsGraphReadyForPeer(peerID int) bool {
	return l.graph.peerGraphReady[peerID]
}

func (l *LatencyMatrixOptimize) ConstructGraph(maxPeers int) error {
	if maxPeers <= int(l.State.Id) {
		return nil
	}
	ready, err := l.isLatencyMatrixReady()
	if err != nil {
		return err
	}
	if !ready {
		return ErrLatencyMatrixNotReady
	}
	maxConnections := int(math.Ceil(float64(l.peerSetUpperBound*maxPeers) / 100))
	connectionsNeeded := make([]bool, maxPeers)
	connectionsNeeded[l.State.Id] = true
	l.graph.constructGraph(maxPeers, maxConnections, connectionsNeeded, l.State.LatencyMatrix)
	return nil
}

func (l *LatencyMatrixOptimize) IsGraphConstructed() bool {
	return l.graph.initiated
}

func (l *LatencyMatrixOptimize) Execute(packetId uint64, data []byte, maxPeers int) error {
	if !l.graph.initiated || maxPeers != l.graph.peers {
		return errGraphNotConstructed
	}
	if l.graph.rootedConnection[l.State.Id] == nil {
		return errGraphConstruction
	}

	return l.send(l.State.Id, packetId, uint64(maxPeers), 1, data, false, 0, 0)
}

func (l *LatencyMatrixOptimize) HandlePacket(packetId uint64, hop uint8, originalSender uint64, _ uint64, maxPeers uint64, data []byte, partial bool, seqNum, total uint16) error {
	if !l.graph.initiated || maxPeers != uint64(l.graph.peers) {
		return errGraphNotConstructed
	}
	if hop == 0 || l.graph.rootedConnection[originalSender] == nil {
		return nil
	}
	if hop > 1 {
		return errInvalidArgumentHop
	}

	return l.send(originalSender, packetId, maxPeers, 0, data, partial, seqNum, total)
}

func (l *LatencyMatrixOptimize) send(root, packetId, maxPeers uint64, hop uint8, data []byte, partial bool, seqNum, total uint16) error {
	for _, peerID := range l.graph.rootedConnection[root] {
		if peer := l.Peers(peerID); peer == nil {
			return errPeerNotFound
		}
	}

	var wg sync.WaitGroup
	for _, peerID := range l.graph.rootedConnection[root] {
		peer := l.Peers(peerID)
		wg.Add(1)
		go func() {
			err := peer.DisseminateRequest(l.Code, packetId, hop, root, uint64(maxPeers), data, partial, seqNum, total)
			if err != nil {
				log.Error("DisseminateRequest err:", err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	return nil
}

type Graph struct {
	constructing atomic.Bool
	initiated    bool
	peers        int
	id           int
	// `rootedConnection[root]` contains the connection array of the node when
	// the original sender is `root`
	rootedConnection [][]int
	peerGraphReady   []bool
	// do not modify the following, this should be read-only
	// `latencyMatrix[u][v]` stores the summation of time to send signal
	// from `u` to `v` and from `v` to `u` in microseconds
	latencyMatrix [][]int64
}

func (g *Graph) constructGraph(
	peers, maxConnections int,
	connectionsNeeded []bool,
	latencyMatrix [][]time.Duration,
) {
	if g.initiated && g.peers == peers {
		return
	}
	if !g.constructing.CompareAndSwap(false, true) {
		return
	}
	g.latencyMatrix = make([][]int64, peers)
	for i := 0; i < peers; i++ {
		g.latencyMatrix[i] = make([]int64, peers)
		for j := 0; j < peers; j++ {
			g.latencyMatrix[i][j] = ((latencyMatrix[i][j] + latencyMatrix[j][i]) / 2).Microseconds()
		}
	}
	g.peers = peers
	g.rootedConnection = make([][]int, peers)
	minConnections := int(math.Ceil(math.Sqrt(float64(peers))))
	maxConnections = max(maxConnections, minConnections)

	connections := make([][][]int, peers)
	wg := sync.WaitGroup{}
	for root := 0; root < peers; root++ {

		wg.Add(1)
		go func(root int) {
			var low int64 = 0
			high := time.Second.Microseconds()
			for low < high {
				mid := (low + high) / 2
				if ready, _ := g.constructConnection(root, maxConnections, mid, connectionsNeeded); !ready {
					low = mid + 1
				} else {
					high = mid
				}
			}
			_, connections[root] = g.constructConnection(root, maxConnections, low, connectionsNeeded)
			g.rootedConnection[root] = connections[root][g.id]
			wg.Done()
		}(root)
	}
	wg.Wait()
	g.initiated = true
	g.peerGraphReady[g.id] = true
	g.constructing.Store(false)
}

// Returns `false` if construction is not possible otherwise returns `true`.
// It is assumed here that `timeMatrix` is a symmetric matrix.
func (g *Graph) constructConnection(
	root, maxConnection int,
	maxTime int64,
	connectionNeeded []bool,
) (ready bool, connections [][]int) {
	degreeInRemain := make([]int, g.peers)
	degreeOut := make([]int, g.peers)
	maxDegreeOut := make([]int, g.peers)
	isConnected := make([]bool, g.peers)
	isRootChild := make([]bool, g.peers)
	isConnected[root] = true
	isRootChild[root] = true
	for u := 0; u < g.peers; u++ {
		for v := 0; v < g.peers; v++ {
			if v == root || v == u || g.msgPropagationTime(root, u, v) > maxTime {
				continue
			}
			degreeInRemain[v]++
			maxDegreeOut[u]++
		}
	}

	connections = make([][]int, g.peers)
	for node, needed := range connectionNeeded {
		if needed {
			connections[node] = make([]int, 0, g.peers-1)
		}
	}
	for {
		node := -1
		for u := 0; u < g.peers; u++ {
			if !isConnected[u] && (node == -1 || degreeInRemain[u] < degreeInRemain[node]) {
				node = u
			}
		}
		if node == -1 {
			break
		}
		if degreeInRemain[node] == 0 {
			return false, nil
		}
		// connect `node` with `root` or someone connected to `root`
		parent := -1
		for u := 0; u < g.peers; u++ {
			if u == node || degreeOut[u] >= maxConnection || g.msgPropagationTime(root, u, node) > maxTime {
				continue
			}

			if isRootChild[u] || (degreeOut[root] < maxConnection && !isConnected[u]) {
				// either `u` is a direct child of `root` or it will be connected to `root`
				// we can travel from `root` to `node` via `u` in `maxTime`
				if parent == -1 || maxDegreeOut[u] < maxDegreeOut[parent] {
					parent = u
				}
			}
		}
		if parent == -1 {
			return false, nil // Sure?
		}

		connectNode := func(node, parent int) {
			degreeOut[parent]++
			for u := 0; u < g.peers; u++ {
				if u == parent || u == node || g.msgPropagationTime(root, u, node) > maxTime {
					continue
				}
				maxDegreeOut[u]--
			}
			isConnected[node] = true
			if parent == root {
				isRootChild[node] = true
			} else {
				for u := 0; u < g.peers; u++ {
					if isConnected[u] || g.msgPropagationTime(root, node, u) > maxTime {
						continue
					}
					degreeInRemain[u]--
				}
			}
			if connectionNeeded[parent] {
				connections[parent] = append(connections[parent], node)
			}
		}
		if !isConnected[parent] {
			// Need to connect `parent` with `root`
			connectNode(parent, root)
		}
		connectNode(node, parent)
	}
	return true, connections
}

// computes msg propagation time in microseconds from `root` to `v` via `u`
func (g *Graph) msgPropagationTime(root, u, v int) int64 {
	return g.latencyMatrix[root][u] + g.latencyMatrix[u][v]
}
