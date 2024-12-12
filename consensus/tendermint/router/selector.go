package router

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
	"sync"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/log"
)

var errUnknownClusters = errors.New("unknown clustering")

type SenderType int

const (
	originator SenderType = iota + 1
	firstRelayerOriginCluster
	localRelayerOriginCluster
	firstRelayerRemoteCluster
	localRelayerRemoteCluster
)

const (
	ScaleThresholdForClustering = 6
	DefaultLatency              = uint(132)
	farThreshold                = 180
	DefaultNearThreshold        = 50
	diversityThreshold          = 40
	MeasurementWindow           = 2000
	MaxLatencyCapFactor         = 0.75
)

type PeerSelector interface {
	SelectPeers(committee *types.Committee, msg message.Msg, from common.Address) ([]common.Address, error)
}

type Selector struct {
	*Router
	heightLock    sync.Mutex
	loggedHR      map[string]uint64
	recentHeights [50]uint64
	heightIndex   int
}

func NewSelector(router *Router) *Selector {
	s := &Selector{Router: router}
	s.loggedHR = make(map[string]uint64)
	s.recentHeights = [50]uint64{}
	s.heightIndex = 0
	return s
}

func (r *Selector) containsAddress(addrs []common.Address, addr common.Address) bool {
	for _, a := range addrs {
		if a == addr {
			return true
		}
	}
	return false
}

func (r *Selector) routingCandidatesFromCluster(clusterID int, exclude []common.Address, committee *types.Committee) []NodeLatency {
	members := r.clusters.base[clusterID].Members
	candidates := make([]NodeLatency, 0, len(members))
	for _, node := range members {
		if r.containsAddress(exclude, node.Addr) {
			continue
		}
		if _, ok := r.broadcaster.FindPeer(node.Addr); ok {
			if member := committee.MemberByAddress(node.Addr); member != nil {
				candidates = append(candidates, node)
			}
		}
	}
	return candidates
}

func (r *Selector) allConnected(recipients []common.Address) bool {
	for _, recipient := range recipients {
		if _, ok := r.broadcaster.FindPeer(recipient); !ok {
			return false
		}
	}
	return true
}

func (r *Selector) SelectPeers(committee *types.Committee, msg message.Msg, from common.Address) ([]common.Address, error) {
	switch msg.Code() {
	case message.ProposalCode:
		return r.selectProposalPeers(committee, msg, from)
	default:
		return r.selectNonProposalPeers(committee, msg, from)
	}
}

func (r *Selector) selectProposalPeers(committee *types.Committee, msg message.Msg, from common.Address) ([]common.Address, error) {
	return r.selectPeersWithBuckets(committee, msg, from, true)
}

func (r *Selector) selectNonProposalPeers(committee *types.Committee, msg message.Msg, from common.Address) ([]common.Address, error) {
	return r.selectPeersWithBuckets(committee, msg, from, false)
}

func (r *Selector) selectPeersWithBuckets(committee *types.Committee, msg message.Msg, from common.Address, isProposal bool) ([]common.Address, error) {
	clusters := r.Clusters()
	if len(clusters.base) == 0 {
		fmt.Println("Selector: no clusters, falling back to all committee members")
		return r.committeeAddresses(committee), nil
	}

	senderClusterID := clusters.clusterContaining(from)
	originClusterID := clusters.clusterContaining(msg.Originator())
	ownClusterID := clusters.clusterContaining(r.self)

	if senderClusterID == -1 || originClusterID == -1 || ownClusterID == -1 {
		fmt.Println("Selector: unknown clusters", "sender", from.Hex(), "originator", msg.Originator().Hex(), "msg hash", msg.Hash().Hex(), "self", r.self.Hex())
		return nil, errors.New("unknown clusters")
	}

	senderType := determineSenderType(from, r.self, msg, originClusterID, ownClusterID, senderClusterID)
	cacheKey := GenerateCacheKey(from, senderType, msg.Code())
	if cached, exists := r.cache.Get(cacheKey); exists {
		if r.allConnected(cached.Recipients) {
			r.cache.UpdateLastUsed(cacheKey)
			r.clusterStatus(r.buildResultFromRecipients(cached.Recipients, clusters), msg, from, senderType, ownClusterID, originClusterID)
			return cached.Recipients, nil
		}
	}

	recipients := r.selectBucketBasedNodes(clusters, committee, senderType, ownClusterID, from, isProposal)
	selected := make([]common.Address, 0, len(recipients))
	for _, r := range recipients {
		selected = append(selected, r.Addr)
	}

	r.cache.Set(cacheKey, selected)
	r.clusterStatus(r.buildResultFromRecipients(selected, clusters), msg, from, senderType, ownClusterID, originClusterID)
	return selected, nil
}

func (r *Selector) selectNodesByLatencySpread() ([]NodeLatency, error) {

	recipients := make([]NodeLatency, 0)
	usedClusters := make(map[int]bool)

	// Select one node from each remote cluster
	for clusterID := range r.clusters.base {
		if clusterID == r.clusters.ownClusterID || usedClusters[clusterID] {
			continue
		}
		var selectedNode *NodeLatency
		// Try primary nodes from this cluster
		for _, primary := range r.clusters.bucketNodes {
			if primary.ClusterID == clusterID && !usedClusters[clusterID] {
				if _, ok := r.broadcaster.FindPeer(primary.Addr); ok {
					selectedNode = &primary
					break
				}
			}
		}
		// Try fallback nodes from this cluster
		if selectedNode == nil {
			for _, fallbacks := range r.clusters.bucketFallbacks {
				for _, node := range fallbacks {
					if node.ClusterID == clusterID && !usedClusters[clusterID] {
						if _, ok := r.broadcaster.FindPeer(node.Addr); ok {
							selectedNode = &node
							break
						}
					}
				}
				if selectedNode != nil {
					break
				}
			}
		}
		// Try any node from the cluster
		if selectedNode == nil {
			for _, node := range r.clusters.base[clusterID].Members {
				if usedClusters[node.ClusterID] || node.Addr == r.self {
					continue
				}
				if _, ok := r.broadcaster.FindPeer(node.Addr); ok {
					selectedNode = &node
					break
				}
			}
		}
		if selectedNode != nil {
			recipients = append(recipients, NodeLatency{selectedNode.Addr, selectedNode.Lat, selectedNode.ClusterID})
			usedClusters[clusterID] = true
		}
	}

	// Select local cluster nodes (up to sqrt(n))
	if r.clusters.ownClusterID != -1 && !usedClusters[r.clusters.ownClusterID] {
		localTarget := int(math.Sqrt(float64(len(r.clusters.base[r.clusters.ownClusterID].Members))))
		localSelected := 0
		// Try local primary nodes
		for _, node := range r.clusters.localBucketNodes {
			if localSelected >= localTarget || usedClusters[node.ClusterID] || node.Addr == r.self {
				continue
			}
			if _, ok := r.broadcaster.FindPeer(node.Addr); ok {
				recipients = append(recipients, node)
				usedClusters[node.ClusterID] = true
				localSelected++
			}
		}
		// Try local fallback nodes
		for _, node := range r.clusters.localBucketFallbacks {
			if localSelected >= localTarget || usedClusters[node.ClusterID] || node.Addr == r.self {
				continue
			}
			if _, ok := r.broadcaster.FindPeer(node.Addr); ok {
				recipients = append(recipients, node)
				usedClusters[node.ClusterID] = true
				localSelected++
			}
		}
		// Try any node from local cluster
		if localSelected < localTarget {
			for _, node := range r.clusters.base[r.clusters.ownClusterID].Members {
				if localSelected >= localTarget || usedClusters[node.ClusterID] || node.Addr == r.self {
					continue
				}
				if _, ok := r.broadcaster.FindPeer(node.Addr); ok {
					recipients = append(recipients, node)
					usedClusters[node.ClusterID] = true
					localSelected++
				}
			}
		}
	}

	return recipients, nil
}

func determineSenderType(from, self common.Address, msg message.Msg, originClusterID, ownClusterID, senderClusterID int) SenderType {
	switch {
	case from == self:
		return originator
	case from == msg.Originator() && originClusterID == ownClusterID:
		return firstRelayerOriginCluster
	case originClusterID == ownClusterID:
		return localRelayerOriginCluster
	case originClusterID != ownClusterID && originClusterID == senderClusterID:
		return firstRelayerRemoteCluster
	case originClusterID != ownClusterID:
		return localRelayerRemoteCluster
	default:
		return localRelayerRemoteCluster
	}
}

func (r *Selector) selectBucketBasedNodes(clusters Clusters, committee *types.Committee, senderType SenderType, ownClusterID int, from common.Address, isProposal bool) []NodeLatency {
	var recipients []NodeLatency
	var minNodes, lowLatencyNodes int

	switch senderType {
	case originator:
		var err error
		recipients, err = r.selectNodesByLatencySpread()
		if err != nil {
			log.Error("Error selecting nodes for originator", "error", err)
			break
		}
		// additional nodes
		if isProposal {
			minNodes = 1
			lowLatencyNodes = 0
		} else {
			minNodes = 2
			lowLatencyNodes = 6
		}
		// 1 closest node from each cluster
		for clusterID, cluster := range clusters.base {
			recipients = append(recipients, r.selectCloseNodes(cluster, committee, clusterID, minNodes, lowLatencyNodes, []common.Address{r.self, from})...)
		}

	case firstRelayerOriginCluster:
		// remote clusters
		minNodes = 1
		lowLatencyNodes = 4
		for clusterID, cluster := range clusters.base {
			if clusterID == ownClusterID {
				continue
			}
			recipients = append(recipients, r.selectCloseNodes(cluster, committee, clusterID, minNodes, lowLatencyNodes, []common.Address{r.self, from})...)
		}
		// local cluster
		minNodes = len(clusters.base[ownClusterID].Members)
		recipients = append(recipients, r.selectCloseNodes(clusters.base[ownClusterID], committee, ownClusterID, minNodes, 0, []common.Address{r.self, from})...)

	case firstRelayerRemoteCluster: // now also includes messages from the first Relayer in origin cluster
		// remote clusters
		minNodes = 0
		lowLatencyNodes = 4
		for clusterID, cluster := range clusters.base {
			if clusterID == ownClusterID {
				continue
			}
			recipients = append(recipients, r.selectCloseNodes(cluster, committee, clusterID, minNodes, lowLatencyNodes, []common.Address{r.self, from})...)
		}
		// local cluster
		minNodes = len(clusters.base[ownClusterID].Members)
		recipients = append(recipients, r.selectCloseNodes(clusters.base[ownClusterID], committee, ownClusterID, minNodes, 0, []common.Address{r.self, from})...)

	case localRelayerOriginCluster, localRelayerRemoteCluster:
		localNodes := len(clusters.base[ownClusterID].Members)
		targetLocalNodes := localNodes
		if isProposal {
			// Select sqrt(n) nodes from local cluster for proposal
			targetLocalNodes = int(math.Sqrt(float64(localNodes)))
		}
		localCandidates := r.routingCandidatesFromCluster(ownClusterID, []common.Address{r.self, from}, committee)
		rand.Shuffle(len(localCandidates), func(i, j int) {
			localCandidates[i], localCandidates[j] = localCandidates[j], localCandidates[i]
		})
		for i := 0; i < targetLocalNodes && i < len(localCandidates); i++ {
			recipients = append(recipients, localCandidates[i])
		}
	}

	recipients = r.deduplicate(recipients)
	// Sort by latency for consistent ordering
	sort.Slice(recipients, func(i, j int) bool { return recipients[i].Lat < recipients[j].Lat })
	return recipients
}

func (r *Selector) deduplicate(recipients []NodeLatency) []NodeLatency {
	seen := make(map[common.Address]struct{}, len(recipients))
	var selected []NodeLatency
	for _, node := range recipients {
		if _, exists := seen[node.Addr]; !exists {
			seen[node.Addr] = struct{}{}
			selected = append(selected, node)
		}
	}
	return selected
}

func (r *Selector) selectCloseNodes(cluster ClusterView, committee *types.Committee, clusterID, minNodes, lowLatencyNodes int, exclude []common.Address) []NodeLatency {
	var selected, candidates []NodeLatency
	candidates = r.routingCandidatesFromCluster(clusterID, exclude, committee)
	if len(candidates) == 0 {
		return nil
	}
	seen := make(map[common.Address]struct{}, len(candidates))
	i := 0
	for ; i < len(candidates) && len(selected) < minNodes; i++ {
		c := candidates[i]
		if _, exists := seen[c.Addr]; !exists {
			seen[c.Addr] = struct{}{}
			selected = append(selected, candidates[i])
		}
	}

	for ; i < len(candidates) && candidates[i].Lat < uint(DefaultNearThreshold) && len(selected) < lowLatencyNodes; i++ {
		c := candidates[i]
		if _, exists := seen[c.Addr]; !exists {
			seen[c.Addr] = struct{}{}
			selected = append(selected, candidates[i])
		}
	}

	return selected
}

func (r *Selector) buildResultFromRecipients(recipients []common.Address, clusters Clusters) []NodeLatency {
	result := make([]NodeLatency, 0, len(recipients))
	for _, recipient := range recipients {
		id := clusters.clusterContaining(recipient)
		member, err := clusters.addressToMember(id, recipient)
		if err == nil {
			result = append(result, member)
		}
	}
	return result
}

func (r *Selector) clusterStatus(recipients []NodeLatency, msg message.Msg, from common.Address, senderType SenderType, ownClusterID int, originClusterID int) {
	logKey := fmt.Sprintf("%d-%d-%d", msg.H(), msg.R(), msg.Code())

	r.heightLock.Lock()
	if _, logged := r.loggedHR[logKey]; logged {
		r.heightLock.Unlock()
		return
	}

	oldHeight := r.recentHeights[r.heightIndex]
	if oldHeight != 0 {
		for k, h := range r.loggedHR {
			if h == oldHeight {
				delete(r.loggedHR, k)
			}
		}
	}

	r.recentHeights[r.heightIndex] = msg.H()
	r.loggedHR[logKey] = msg.H()
	r.heightIndex = (r.heightIndex + 1) % 50
	r.heightLock.Unlock()

	var sb strings.Builder
	totalDisconnected := 0
	totalSelected := 0
	var fullyConnectedClusters []string
	var totalConnected int
	sender := "originator"
	switch senderType {
	case originator:
		sender = "originator"
	case firstRelayerOriginCluster:
		sender = "first relayer origin cluster"
	case localRelayerOriginCluster:
		sender = "local relayer origin cluster"
	case firstRelayerRemoteCluster:
		sender = "first relayer remote cluster"
	case localRelayerRemoteCluster:
		sender = "local relayer remote cluster"
	}

	msgType := "proposal"
	switch msg.Code() {
	case message.ProposalCode:
		msgType = "Proposal"
	case message.PrevoteCode:
		msgType = "Prevote"
	case message.PrecommitCode:
		msgType = "Precommit"
	case message.LightProposalCode:
		msgType = "Light Proposal"
	default:
		msgType = "Unknown"
	}

	sb.WriteString(fmt.Sprintf("\nCluster routing status:\t Height=%d, Round=%d, From=%s Message=%s MessageHash=%s SenderType=%s localCluster=%d originCluster=%d\n",
		msg.H(), msg.R(), from.Hex(), msgType, msg.Hash().Hex(), sender, ownClusterID, originClusterID))

	clusterMap := make(map[int][]NodeLatency)
	for _, node := range recipients {
		clusterMap[node.ClusterID] = append(clusterMap[node.ClusterID], node)
	}

	for clusterID, cluster := range clusterMap {
		var lostPeers []string
		connectedCount := 0
		latencyList := []string{}

		for _, peer := range cluster {
			_, ok := r.broadcaster.FindPeer(peer.Addr)
			if ok {
				connectedCount++
				totalConnected++
				latencyList = append(latencyList, fmt.Sprintf("%s-%d", peer.Addr.Hex(), peer.Lat))
			} else {
				lostPeers = append(lostPeers, peer.Addr.Hex())
				totalDisconnected++
			}
			totalSelected++
		}

		if len(lostPeers) == 0 && len(cluster) > 0 {
			fullyConnectedClusters = append(fullyConnectedClusters,
				fmt.Sprintf("C%d:%d L:%s\n", clusterID, len(cluster), latencyList))
		} else {
			sb.WriteString(fmt.Sprintf("Cluster #%d: selected:%d connected:%d\nL:%s\n", clusterID, len(cluster), connectedCount, latencyList))
			sb.WriteString("  X disconnected:")
			for _, peerHex := range lostPeers {
				sb.WriteString(" ")
				sb.WriteString(peerHex)
			}
			sb.WriteByte('\n')
		}
	}

	if len(fullyConnectedClusters) > 0 {
		sb.WriteString("Fully connected: ")
		for i, clusterInfo := range fullyConnectedClusters {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(clusterInfo)
		}
		sb.WriteByte('\n')
	}

	sb.WriteString(fmt.Sprintf("Total: selected:%d connected:%d disconnected:%d\n", totalSelected, totalConnected, totalDisconnected))

	log.Info(sb.String())
}
