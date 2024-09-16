package strats

import (
	"fmt"
	"math/rand"
	"sort"

	"golang.org/x/exp/slices"

	"github.com/autonity/autonity/cmd/netdiag/strats/kmeans"
	"github.com/autonity/autonity/log"
)

type KmeansNLeadersGraphConstructor struct {
	BaseStrategy
	ByzantineChance float64
	NumClusters     int
	NLeaders        int

	graph           Graph
	localLeaders    [][]int
	clusters        [][]int
	assignedCluster int
}

func init() {
	registerStrategy("K-Means Leaf Forward - (k=6, byz=30%, n=3)", func(base BaseStrategy) Strategy {
		return createKMeansNLeaders(base, 6, 3, 0.3)
	})
	registerStrategy("K-Means NLeaders - (k=6, byz=0%, n=3)", func(base BaseStrategy) Strategy {
		return createKMeansNLeaders(base, 6, 3, 0)
	})
	registerStrategy("K-Means Leaf Forward - (k=6, byz=0%, leaders=6)", func(base BaseStrategy) Strategy {
		return createKMeansNLeaders(base, 6, 6, 0)
	})
}

func createKMeansNLeaders(base BaseStrategy, numClusters int, nleaders int, byzantineChance float64) *GraphStrategy {
	graph := Graph{
		id:             int(base.State.Id),
		peerGraphReady: make([]bool, base.State.Peers),
	}
	constructor := &KmeansNLeadersGraphConstructor{
		BaseStrategy:    base,
		ByzantineChance: byzantineChance,
		NumClusters:     numClusters,
		NLeaders:        nleaders,
		graph:           graph,
		localLeaders:    make([][]int, numClusters),
		clusters:        make([][]int, numClusters),
		assignedCluster: -1,
	}
	return &GraphStrategy{
		BaseStrategy:     base,
		GraphConstructor: constructor,
		peerGraphReady:   make([]bool, base.State.Peers),
	}
}

func (k *KmeansNLeadersGraphConstructor) ConstructGraph(_ int) error {
	ready, err := k.isLatencyMatrixReady()
	if err != nil {
		return err
	}
	if !ready {
		return ErrLatencyMatrixNotReady
	}
	return k.constructGraph()
}

func (k *KmeansNLeadersGraphConstructor) LatencyType() (LatencyType, int) {
	return LatencyTypeRelative, k.State.Peers
}

func (k *KmeansNLeadersGraphConstructor) RouteBroadcast(originalSender int, from int) ([]int, error) {
	log.Debug("Routing packet", "originalSender", originalSender, "localId", k.State.Id)
	// we should only cut the packet if we are not the originator of the packet
	if rand.Float64() < k.ByzantineChance && originalSender != int(k.State.Id) {
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
		destinationPeers = filterArray(flatten(k.localLeaders), func(i int) bool {
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

func (k *KmeansNLeadersGraphConstructor) isLatencyMatrixReady() (bool, error) {
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

func (k *KmeansNLeadersGraphConstructor) constructGraph() error {
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
	for i := 0; i < k.State.Peers; i++ {
		rootedConnection[i] = make([]int, 0)
	}
	for i, c := range clusters {
		clusterLeaders := minNItems(c, k.NLeaders)
		log.Debug("Cluster leader elected", "iCluster", i, "clusterLeader", clusterLeaders)
		k.localLeaders[i] = minNItems(c, k.NLeaders)
		for _, clusterLeader := range clusterLeaders {
			rootedConnection[clusterLeader] = filterArray(c, func(i int) bool {
				return i != clusterLeader
			})
		}
		if slices.Contains(clusterLeaders, int(k.State.Id)) {
			k.assignedCluster = i
		}
	}

	k.graph.rootedConnection = rootedConnection
	k.graph.initiated = true
	k.graph.peerGraphReady[k.State.Id] = true
	k.graph.constructing.Store(false)
	k.clusters = clusters
	return nil
}

func minNItems(slice []int, n int) []int {
	if n >= len(slice) {
		return slice
	}
	sliceCopy := append([]int(nil), slice...)
	sort.Ints(sliceCopy)
	return sliceCopy[:n]
}

func flatten(slice [][]int) []int {
	var result []int
	for _, subSlice := range slice {
		result = append(result, subSlice...)
	}
	return result
}
