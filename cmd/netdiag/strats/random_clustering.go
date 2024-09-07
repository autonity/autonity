package strats

import (
	"fmt"
	"math/rand"

	"github.com/autonity/autonity/log"
)

type RandomClustering struct {
	BaseStrategy
	graph        Graph
	localLeaders []int
	NumClusters  int
}

type RandomClusteringGraphConstructor struct {
	stateId      uint64
	numClusters  int
	localLeaders []int
	graph        Graph
}

func init() {
	registerStrategy("Random Clustering - 6 clusters", func(base BaseStrategy) Strategy {
		return createRandomClustering(base, 6)
	})
}

func createRandomClustering(base BaseStrategy, numClusters int) *GraphStrategy {
	graph := Graph{
		id:             int(base.State.Id),
		peerGraphReady: make([]bool, base.State.Peers),
	}
	strategy := &GraphStrategy{
		base,
		&RandomClusteringGraphConstructor{stateId: base.State.Id, numClusters: numClusters, localLeaders: make([]int, 0), graph: graph},
		make([]bool, base.State.Peers),
	}
	return strategy
}

func (r *RandomClusteringGraphConstructor) ConstructGraph(maxPeers int) error {
	if r.graph.initiated {
		return nil
	}

	if !r.graph.constructing.CompareAndSwap(false, true) {
		return nil
	}

	r.graph.initiated = true
	r.graph.peers = maxPeers
	r.graph.rootedConnection = make([][]int, maxPeers)
	// randomly assign local leaders
	random := rand.New(rand.NewSource(12345))
	for i := 0; i < r.numClusters; i++ {
		added := false
		for !added {
			if proposedLeader := random.Intn(maxPeers); !alreadyIn(proposedLeader, r.localLeaders) {
				r.localLeaders = append(r.localLeaders, proposedLeader)
				added = true
			}
		}
	}
	log.Info(
		"Elected local leaders",
		"localLeaders",
		fmt.Sprintf("%v", r.localLeaders),
		"numClusters",
		r.numClusters,
	)
	clusterMap := make(map[int]int)
	for i := 0; i < maxPeers; i++ {
		if alreadyIn(i, r.localLeaders) {
			clusterMap[i] = indexOf(i, r.localLeaders)
		} else {
			clusterMap[i] = random.Intn(r.numClusters)
		}
	}

	// assign local leaders to the graph
	for rootPeer := 0; rootPeer < maxPeers; rootPeer++ {
		r.graph.rootedConnection[rootPeer] = make([]int, 0)
		if alreadyIn(rootPeer, r.localLeaders) {
			cluster := indexOf(rootPeer, r.localLeaders)
			for peer, peerCluster := range clusterMap {
				if peerCluster == cluster && peer != rootPeer {
					r.graph.rootedConnection[rootPeer] = append(r.graph.rootedConnection[rootPeer], peer)
				}
			}
		}

	}
	log.Info(
		"Constructed graph",
		"rootedConnection",
		fmt.Sprintf("%v", r.graph.rootedConnection),
		"localLeaders",
		fmt.Sprintf("%v", r.localLeaders),
	)
	r.graph.peerGraphReady[r.stateId] = true
	return nil
}

func (r *RandomClusteringGraphConstructor) RouteBroadcast(originalSender int, _ int) ([]int, error) {
	if !r.graph.initiated {
		return nil, ErrGraphConstruction
	}

	if len(r.graph.rootedConnection) == 0 {
		log.Error("Graph not initiated, rootedConnection is empty")
		return nil, ErrGraphConstruction
	}

	// first collect all peers to send to
	var destinationPeers []int

	// if we are the originator of the packet, we need to send to all local leaders
	if originalSender == int(r.stateId) {
		log.Debug(
			"Original sender, adding all local leaders",
			"localLeaders",
			fmt.Sprintf("%v", r.localLeaders),
			"root",
			originalSender,
		)
		// we don't need to send it to ourselves
		destinationPeers = filterArray(r.localLeaders, func(i int) bool {
			return i != int(r.stateId)
		})
	}

	// if we are not the originator of the packet we only need to send it if we are the local leader
	if len(r.graph.rootedConnection[r.stateId]) > 0 {
		destinationPeers = append(
			destinationPeers,
			// we don't need to send it to the originator
			filterArray(r.graph.rootedConnection[r.stateId], func(i int) bool {
				return i != int(originalSender)
			})...,
		)
	}

	return destinationPeers, nil
}

func alreadyIn(item int, list []int) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}

func indexOf(item int, list []int) int {
	for i, v := range list {
		if v == item {
			return i
		}
	}
	return -1
}
