package strats

import (
	"errors"
	"fmt"

	"github.com/autonity/autonity/cmd/netdiag/strats/kmeans"
	"github.com/autonity/autonity/log"
)

var ErrGraphNotInitiated = errors.New("graph not initiated")
var ErrGraphConstruction = errors.New("invalid graph construction")

type KmeansClusteringGraphConstructor struct {
	BaseStrategy
	graph        Graph
	localLeaders []int
	NumClusters  int
}

func init() {
	registerStrategy("K-Means Clustering - 6 clusters", func(base BaseStrategy) Strategy {
		return createKmeansClustering(base, 6)
	})
}

func createKmeansClustering(base BaseStrategy, numClusters int) *GraphStrategy {
	graph := Graph{
		id:             int(base.State.Id),
		peerGraphReady: make([]bool, base.State.Peers),
	}
	constructor := &KmeansClusteringGraphConstructor{base, graph, nil, numClusters}
	return &GraphStrategy{
		BaseStrategy:     base,
		GraphConstructor: constructor,
		peerGraphReady:   make([]bool, base.State.Peers),
	}
}

func (k *KmeansClusteringGraphConstructor) ConstructGraph(_ int) error {
	ready, err := k.isLatencyMatrixReady()
	if err != nil {
		return err
	}
	if !ready {
		return ErrLatencyMatrixNotReady
	}
	return k.constructGraph()
}

func (k *KmeansClusteringGraphConstructor) RouteBroadcast(originalSender int) ([]int, error) {
	log.Debug("Routing packet", "originalSender", originalSender, "localId", k.State.Id)
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
	return destinationPeers, nil
}

func (k *KmeansClusteringGraphConstructor) isLatencyMatrixReady() (bool, error) {
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

func (k *KmeansClusteringGraphConstructor) constructGraph() error {
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
	var localLeaders []int
	for i := 0; i < k.State.Peers; i++ {
		rootedConnection[i] = make([]int, 0)
	}
	for i, c := range clusters {
		clusterLeader := minArray(c)
		log.Debug("Cluster leader elected", "iCluster", i, "clusterLeader", clusterLeader)
		localLeaders = append(localLeaders, clusterLeader)
		for _, peer := range c {
			if peer != clusterLeader {
				rootedConnection[clusterLeader] = append(rootedConnection[clusterLeader], peer)
			}
		}
	}

	k.graph.rootedConnection = rootedConnection
	k.localLeaders = localLeaders
	k.graph.initiated = true
	k.graph.peerGraphReady[k.State.Id] = true
	k.graph.constructing.Store(false)
	return nil
}

func minArray(a []int) int {
	minimum := a[0]
	for _, v := range a {
		if v < minimum {
			minimum = v
		}
	}
	return minimum
}

func filterArray(a []int, f func(int) bool) []int {
	var result []int
	for _, v := range a {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}
