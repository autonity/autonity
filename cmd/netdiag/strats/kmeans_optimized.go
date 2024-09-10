package strats

import (
	"fmt"
	"math/rand"

	"github.com/autonity/autonity/cmd/netdiag/strats/kmeans"
	"github.com/autonity/autonity/log"
)

type KmeansOptimizedGraphConstructor struct {
	BaseStrategy
	ByzantineChance float64
	NumClusters     int

	graph           Graph
	localLeaders    []int
	clusters        [][]int
	assignedCluster int
}

func init() {
	registerStrategy("K-Means Optimized - (k=6, byz=30%)", func(base BaseStrategy) Strategy {
		return createKmeansOptimized(base, 6, 0.3)
	})
	registerStrategy("K-Means Optimized - (k=6, byz=0%)", func(base BaseStrategy) Strategy {
		return createKmeansOptimized(base, 6, 0)
	})
}

func createKmeansOptimized(base BaseStrategy, numClusters int, byzantineChance float64) *GraphStrategy {
	graph := Graph{
		id:             int(base.State.Id),
		peerGraphReady: make([]bool, base.State.Peers),
	}
	constructor := &KmeansOptimizedGraphConstructor{
		BaseStrategy:    base,
		ByzantineChance: byzantineChance,
		NumClusters:     numClusters,
		graph:           graph,
		localLeaders:    nil,
		clusters:        nil,
		assignedCluster: -1,
	}
	return &GraphStrategy{
		BaseStrategy:     base,
		GraphConstructor: constructor,
		peerGraphReady:   make([]bool, base.State.Peers),
	}
}

func (k *KmeansOptimizedGraphConstructor) ConstructGraph(_ int) error {
	ready, err := k.isLatencyMatrixReady()
	if err != nil {
		return err
	}
	if !ready {
		return ErrLatencyMatrixNotReady
	}
	return k.constructGraph()
}

func (k *KmeansOptimizedGraphConstructor) LatencyType() (LatencyType, int) {
	return LatencyTypeRelative, k.State.Peers
}

func (k *KmeansOptimizedGraphConstructor) RouteBroadcast(originalSender int, from int) ([]int, error) {
	log.Debug("Routing packet", "originalSender", originalSender, "localId", k.State.Id)
	if rand.Float64() < k.ByzantineChance {
		log.Debug("Byzantine behavior, not routing packet", "localId", k.State.Id)
		return nil, nil
	}

	if len(k.graph.rootedConnection) == 0 {
		log.Error("Graph not initiated, rootedConnection is empty")
		return nil, ErrGraphConstruction
	}
	// first collect all peers to send to
	var destinationPeers []int

	// if we are the originator of the packet, we need to send to all local leaders
	if originalSender == int(k.State.Id) {
		log.Debug("Original sender, adding all local leaders", "localLeaders", fmt.Sprintf("%v", k.localLeaders), "originalSender", originalSender)
		// we don't need to send it to ourselves
		destinationPeers = filterArray(k.localLeaders, func(i int) bool {
			return i != int(k.State.Id)
		})
	}

	// if we are not the originator of the packet we only need to send it if we are the local leader
	if len(k.graph.rootedConnection[k.State.Id]) > 0 {
		destinationPeers = append(
			destinationPeers,
			// we don't need to send it to the originator
			filterArray(k.graph.rootedConnection[k.State.Id], func(i int) bool {
				return i != originalSender
			})...,
		)
	}

	// if we are a leaf node, and the original sender is in our cluster,
	// we can forward it to a random peer in another cluster
	if len(destinationPeers) == 0 && from == k.localLeaders[k.assignedCluster] {
		log.Debug("Forwarding packet to random peer in another cluster", "localId", k.State.Id)
		randCluster := randNotEqual(k.NumClusters, k.assignedCluster)
		destinationPeers = append(destinationPeers, k.clusters[randCluster][randNotEqual(len(k.clusters[randCluster]), k.localLeaders[randCluster])])
	}

	// if we receive from someone other than our local leader, we can forward it to our local leader
	if len(destinationPeers) == 0 && from != k.localLeaders[k.assignedCluster] && !inCluster(from, k.assignedCluster, k.clusters) {
		log.Debug(
			"Forwarding packet to local cluster",
			"localId",
			k.State.Id,
			"assignedCluster",
			k.assignedCluster,
			"cluster",
			k.clusters[k.assignedCluster],
		)
		destinationPeers = append(
			destinationPeers,
			filterArray(k.clusters[k.assignedCluster], func(i int) bool {
				return i != int(k.State.Id) && i != k.localLeaders[k.assignedCluster]
			})...,
		)
	}
	return destinationPeers, nil
}

func (k *KmeansOptimizedGraphConstructor) isLatencyMatrixReady() (bool, error) {
	for id, array := range k.State.LatencyMatrix {
		if len(array) != k.State.Peers {
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

func (k *KmeansOptimizedGraphConstructor) constructGraph() error {
	log.Debug("Constructing graph")
	if k.graph.initiated {
		return nil
	}

	if !k.graph.constructing.CompareAndSwap(false, true) {
		return nil
	}

	log.Debug("Assigning clusters ", "numClusters", k.NumClusters)
	clusters, err := kmeans.AssignClusters(k.State.LatencyMatrix, k.NumClusters)
	if err != nil {
		return err
	}
	log.Debug("Clusters assigned", "clusters", fmt.Sprintf("%v", clusters))

	var rootedConnection = make([][]int, k.State.Peers)
	k.localLeaders = make([]int, k.NumClusters)
	for i := 0; i < k.State.Peers; i++ {
		rootedConnection[i] = make([]int, 0)
	}
	for i, c := range clusters {
		clusterLeader := minArray(c)
		log.Debug("Cluster leader elected", "iCluster", i, "clusterLeader", clusterLeader)
		k.localLeaders[i] = clusterLeader
		for _, peer := range c {
			if peer != clusterLeader {
				rootedConnection[clusterLeader] = append(rootedConnection[clusterLeader], peer)
			}
			if peer == int(k.State.Id) {
				k.assignedCluster = i
			}
		}
	}

	k.graph.rootedConnection = rootedConnection
	k.graph.initiated = true
	k.graph.peerGraphReady[k.State.Id] = true
	k.graph.constructing.Store(false)
	k.clusters = clusters
	return nil
}

func inCluster(peer int, cluster int, clusters [][]int) bool {
	for _, c := range clusters[cluster] {
		if c == peer {
			return true
		}
	}
	return false
}

func randNotEqual(max, b int) int {
	r := rand.Intn(max)
	if r != b {
		return r
	}
	return randNotEqual(max, b)
}
