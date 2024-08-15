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
	errLatencyMatrixNotReady = errors.New("latency matrix not ready")
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
	strategy := &LatencyMatrixOptimize{base, Graph{id: int(base.State.Id)}, peerSetUpperBound}
	// go strategy.start()
	return strategy
}

func (l *LatencyMatrixOptimize) start() {
	for {
		ready, err := l.isLatencyMatrixReady()
		if err != nil {
			log.Crit("error in latency matrix", "err", err)
		}
		if ready {
			break
		}
	}
	l.constructGraph(len(l.State.LatencyMatrix))
}

func (l *LatencyMatrixOptimize) isLatencyMatrixReady() (bool, error) {
	for id, array := range l.State.LatencyMatrix {
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

func (l *LatencyMatrixOptimize) constructGraph(peers int) {
	if l.graph.initiated && l.graph.peers == peers {
		return
	}
	if !l.graph.constructing.CompareAndSwap(false, true) {
		return
	}
	l.graph.peers = peers
	timeMatrix := make([][]int64, peers)
	for i := 0; i < peers; i++ {
		timeMatrix[i] = make([]int64, peers)
	}
	l.graph.rootedConnection = make([][]int, peers)
	maxConnections := int(math.Ceil(float64(l.peerSetUpperBound*peers) / 100))
	minConnections := int(math.Ceil(math.Sqrt(float64(peers))))
	maxConnections = max(maxConnections, minConnections)

	for root := 0; root < peers; root++ {
		for u := 0; u < peers; u++ {
			for v := 0; v < peers; v++ {
				timeMatrix[u][v] = (l.State.LatencyMatrix[root][u] + l.State.LatencyMatrix[u][root] +
					l.State.LatencyMatrix[u][v] + l.State.LatencyMatrix[v][u]).Microseconds()
			}
		}

		var low int64 = 0
		high := 10 * time.Second.Microseconds()
		for low < high {
			mid := (low + high) / 2
			if !l.graph.constructConnection(root, maxConnections, mid, timeMatrix) {
				low = mid + 1
			} else {
				high = mid
			}
		}
		l.graph.constructConnection(root, maxConnections, low, timeMatrix)
	}
	l.graph.initiated = true
	l.graph.constructing.Store(false)
}

func (l *LatencyMatrixOptimize) ConstructGraph(maxPeers int) error {
	ready, err := l.isLatencyMatrixReady()
	if err != nil {
		return err
	}
	if !ready {
		return errLatencyMatrixNotReady
	}
	l.constructGraph(maxPeers)
	return nil
}

func (l *LatencyMatrixOptimize) Execute(packetId uint64, data []byte, maxPeers int) error {
	if !l.graph.initiated || maxPeers != l.graph.peers {
		return errGraphNotConstructed
	}
	if l.graph.rootedConnection[l.State.Id] == nil {
		return errGraphConstruction
	}

	return l.send(l.State.Id, packetId, uint64(maxPeers), 1, data)
}

func (l *LatencyMatrixOptimize) HandlePacket(packetId uint64, hop uint8, originalSender uint64, maxPeers uint64, data []byte, partial bool, seqNum, total uint16) error {
	if !l.graph.initiated || maxPeers != uint64(l.graph.peers) {
		return errGraphNotConstructed
	}
	if hop == 0 || l.graph.rootedConnection[originalSender] == nil {
		return nil
	}
	if hop > 1 {
		return errInvalidArgumentHop
	}

	return l.send(originalSender, packetId, maxPeers, 0, data)
}

func (l *LatencyMatrixOptimize) send(root, packetId, maxPeers uint64, hop uint8, data []byte) error {
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
			err := peer.DisseminateRequest(l.Code, packetId, hop, root, uint64(maxPeers), data, false, 0, 0)
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
}

// Returns `false` if construction is not possible otherwise returns `true`.
// It is assumed here that `timeMatrix` is a symmetric matrix.
func (g *Graph) constructConnection(
	root, maxConnection int,
	maxTime int64,
	timeMatrix [][]int64,
) bool {
	degreeInRemain := make([]int, g.peers)
	degreeOut := make([]int, g.peers)
	maxDegreeOut := make([]int, g.peers)
	isConnected := make([]bool, g.peers)
	isRootChild := make([]bool, g.peers)
	isConnected[root] = true
	isRootChild[root] = true
	for u := 0; u < g.peers; u++ {
		for v := 0; v < g.peers; v++ {
			if v == root || v == u || timeMatrix[u][v] > maxTime {
				continue
			}
			degreeInRemain[v]++
			maxDegreeOut[u]++
		}
	}

	connection := make([]int, 0, g.peers-1)
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
			return false
		}
		// connect `node` with `root` or someone connected to `root`
		parent := -1
		for u := 0; u < g.peers; u++ {
			if u == node || degreeOut[u] >= maxConnection || timeMatrix[u][node] > maxTime {
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
			return false // Sure?
		}

		connectNode := func(node, parent int) {
			degreeOut[parent]++
			for u := 0; u < g.peers; u++ {
				if u == node || timeMatrix[u][node] > maxTime {
					continue
				}
				if u != parent {
					maxDegreeOut[u]--
				}
				if parent != root {
					degreeInRemain[u]--
				}
			}
			isConnected[node] = true
			if parent == root {
				isRootChild[node] = true
			}
			if parent == g.id {
				connection = append(connection, node)
			}
		}
		if !isRootChild[parent] {
			// Need to connect `parent` with `root`
			connectNode(parent, root)
		}
		connectNode(node, parent)
	}
	g.rootedConnection[root] = connection
	return true
}
