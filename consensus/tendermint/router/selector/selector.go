package selector

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
	"github.com/autonity/autonity/consensus/tendermint/router/cache"
	"github.com/autonity/autonity/consensus/tendermint/router/constants"
	"github.com/autonity/autonity/consensus/tendermint/router/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/router/network"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/log"
)

type SenderType int

const (
	originator SenderType = iota + 1
	firstRelayerOriginCluster
	localRelayerOriginCluster
	firstRelayerRemoteCluster
	localRelayerRemoteCluster
)

type Selector struct {
	networkProvider interfaces.NetworkProvider
	recipientCache  cache.Recipients
	peerFinder      interfaces.PeerFinder
	heightLock      sync.Mutex
	loggedHR        map[string]uint64
	recentHeights   [50]uint64
	heightIndex     int
}

func New(np interfaces.NetworkProvider, cache cache.Recipients) *Selector {
	s := &Selector{
		networkProvider: np,
		recipientCache:  cache,
		loggedHR:        make(map[string]uint64),
		recentHeights:   [50]uint64{},
		heightIndex:     0,
	}
	return s
}

func (s *Selector) SetBroadcaster(broadcaster interfaces.PeerFinder) {
	s.peerFinder = broadcaster
}

func (s *Selector) containsAddress(addrs []common.Address, addr common.Address) bool {
	for _, a := range addrs {
		if a == addr {
			return true
		}
	}
	return false
}

func (s *Selector) routingCandidatesFromCluster(clusterID int, exclude []common.Address, committee *types.Committee) []network.Node {
	clusters := s.networkProvider.Clusters()
	members := clusters.MembersByID(clusterID)
	candidates := make([]network.Node, 0, len(members))
	for _, node := range members {
		if s.containsAddress(exclude, node.Addr) {
			continue
		}
		if _, ok := s.peerFinder.FindPeer(node.Addr); ok {
			if member := committee.MemberByAddress(node.Addr); member != nil {
				candidates = append(candidates, node)
			}
		}
	}
	return candidates
}

func (s *Selector) allConnected(recipients []common.Address) bool {
	for _, recipient := range recipients {
		if _, ok := s.peerFinder.FindPeer(recipient); !ok {
			return false
		}
	}
	return true
}

func (s *Selector) SelectPeers(committee *types.Committee, msg message.Msg, from common.Address) ([]common.Address, error) {
	switch msg.Code() {
	case message.ProposalCode:
		return s.selectProposalPeers(committee, msg, from)
	default:
		return s.selectNonProposalPeers(committee, msg, from)
	}
}

func (s *Selector) selectProposalPeers(committee *types.Committee, msg message.Msg, from common.Address) ([]common.Address, error) {
	return s.selectPeersWithBuckets(committee, msg, from, true)
}

func (s *Selector) selectNonProposalPeers(committee *types.Committee, msg message.Msg, from common.Address) ([]common.Address, error) {
	return s.selectPeersWithBuckets(committee, msg, from, false)
}

func (s *Selector) selectPeersWithBuckets(committee *types.Committee, msg message.Msg, from common.Address, isProposal bool) ([]common.Address, error) {
	clusters := s.networkProvider.Clusters()
	if len(clusters.Base()) == 0 {
		return []common.Address{}, errors.New("no clusters")
	}

	senderClusterID := clusters.IDByAddress(from)
	originClusterID := clusters.IDByAddress(msg.Originator())
	ownClusterID := clusters.ID()

	if senderClusterID == -1 || originClusterID == -1 || ownClusterID == -1 {
		fmt.Println("Selector: unknown clusters", "sender", from.Hex(), "originator", msg.Originator().Hex(), "msg hash", msg.Hash().Hex(), "self", clusters.Self().Hex())
		return nil, errors.New("unknown clusters")
	}

	senderType := determineSenderType(from, clusters.Self(), msg, originClusterID, ownClusterID, senderClusterID)
	cacheKey := cache.GenerateKey(from, int(senderType), msg.Code())
	if cached, exists := s.recipientCache.Get(cacheKey); exists {
		if s.allConnected(cached.Recipients) {
			s.recipientCache.UpdateLastUsed(cacheKey)
			s.clusterStatus(s.buildResultFromRecipients(cached.Recipients, clusters), msg, from, senderType, ownClusterID, originClusterID)
			return cached.Recipients, nil
		}
	}

	recipients := s.selectBucketBasedNodes(clusters, committee, senderType, ownClusterID, from, isProposal)
	selected := make([]common.Address, 0, len(recipients))
	for _, r := range recipients {
		selected = append(selected, r.Addr)
	}

	s.recipientCache.Set(cacheKey, selected)
	s.clusterStatus(s.buildResultFromRecipients(selected, clusters), msg, from, senderType, ownClusterID, originClusterID)
	return selected, nil
}

func (s *Selector) selectNodesByLatencySpread() ([]network.Node, error) {

	recipients := make([]network.Node, 0)
	usedClusters := make(map[int]bool)

	clusters := s.networkProvider.Clusters()
	// Select one node from each remote cluster
	for clusterID := range clusters.Base() {
		if clusterID == clusters.ID() || usedClusters[clusterID] {
			continue
		}
		var selectedNode *network.Node
		// Try primary nodes from this cluster
		for _, primary := range clusters.BucketNodes() {
			if primary.ClusterID == clusterID && !usedClusters[clusterID] {
				if _, ok := s.peerFinder.FindPeer(primary.Addr); ok {
					selectedNode = &primary
					break
				}
			}
		}
		// Try fallback nodes from this cluster
		if selectedNode == nil {
			for _, fallbacks := range clusters.BucketFallBacks() {
				for _, node := range fallbacks {
					if node.ClusterID == clusterID && !usedClusters[clusterID] {
						if _, ok := s.peerFinder.FindPeer(node.Addr); ok {
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
			for _, node := range clusters.MembersByID(clusterID) {
				if usedClusters[node.ClusterID] {
					continue
				}
				if _, ok := s.peerFinder.FindPeer(node.Addr); ok {
					selectedNode = &node
					break
				}
			}
		}
		if selectedNode != nil {
			recipients = append(recipients, network.Node{selectedNode.Addr, selectedNode.Lat, selectedNode.ClusterID})
			usedClusters[clusterID] = true
		}
	}

	// Select local cluster nodes (up to sqrt(n))
	if clusters.ID() != -1 && !usedClusters[clusters.ID()] {
		localTarget := int(math.Sqrt(float64(len(clusters.MembersByID(clusters.ID())))))
		localSelected := 0
		// Try local primary nodes
		for _, node := range clusters.LocalBucketNodes() {
			if localSelected >= localTarget || usedClusters[node.ClusterID] {
				continue
			}
			if _, ok := s.peerFinder.FindPeer(node.Addr); ok {
				recipients = append(recipients, node)
				usedClusters[node.ClusterID] = true
				localSelected++
			}
		}
		// Try local fallback nodes
		for _, node := range clusters.LocalBucketFallBacks() {
			if localSelected >= localTarget || usedClusters[node.ClusterID] {
				continue
			}
			if _, ok := s.peerFinder.FindPeer(node.Addr); ok {
				recipients = append(recipients, node)
				usedClusters[node.ClusterID] = true
				localSelected++
			}
		}
		// Try any node from local cluster
		if localSelected < localTarget {
			for _, node := range clusters.MembersByID(clusters.ID()) {
				if localSelected >= localTarget || usedClusters[node.ClusterID] {
					continue
				}
				if _, ok := s.peerFinder.FindPeer(node.Addr); ok {
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

func (s *Selector) selectBucketBasedNodes(clusters network.Clusters, committee *types.Committee, senderType SenderType, ownClusterID int, from common.Address, isProposal bool) []network.Node {
	var recipients []network.Node
	var minNodes, lowLatencyNodes int

	switch senderType {
	case originator:
		var err error
		recipients, err = s.selectNodesByLatencySpread()
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
		for clusterID := range clusters.Base() {
			recipients = append(recipients, s.selectCloseNodes(committee, clusterID, minNodes, lowLatencyNodes, []common.Address{from})...)
		}

	case firstRelayerOriginCluster:
		// remote clusters
		minNodes = 1
		lowLatencyNodes = 4
		for clusterID := range clusters.Base() {
			if clusterID == ownClusterID {
				continue
			}
			recipients = append(recipients, s.selectCloseNodes(committee, clusterID, minNodes, lowLatencyNodes, []common.Address{from})...)
		}
		// local cluster
		minNodes = len(clusters.MembersByID(clusters.ID()))
		recipients = append(recipients, s.selectCloseNodes(committee, ownClusterID, minNodes, 0, []common.Address{from})...)

	case firstRelayerRemoteCluster: // now also includes messages from the first Relayer in origin cluster
		// remote clusters
		minNodes = 0
		lowLatencyNodes = 4
		for clusterID := range clusters.Base() {
			if clusterID == ownClusterID {
				continue
			}
			recipients = append(recipients, s.selectCloseNodes(committee, clusterID, minNodes, lowLatencyNodes, []common.Address{from})...)
		}
		// local cluster
		minNodes = len(clusters.Base()[ownClusterID])
		recipients = append(recipients, s.selectCloseNodes(committee, ownClusterID, minNodes, 0, []common.Address{from})...)

	case localRelayerOriginCluster, localRelayerRemoteCluster:
		localNodes := len(clusters.Base()[ownClusterID])
		targetLocalNodes := localNodes
		if isProposal {
			// Select sqrt(n) nodes from local cluster for proposal
			targetLocalNodes = int(math.Sqrt(float64(localNodes)))
		}
		localCandidates := s.routingCandidatesFromCluster(ownClusterID, []common.Address{from}, committee)
		rand.Shuffle(len(localCandidates), func(i, j int) {
			localCandidates[i], localCandidates[j] = localCandidates[j], localCandidates[i]
		})
		for i := 0; i < targetLocalNodes && i < len(localCandidates); i++ {
			recipients = append(recipients, localCandidates[i])
		}
	}

	recipients = s.deduplicate(recipients)
	// Sort by latency for consistent ordering
	sort.Slice(recipients, func(i, j int) bool { return recipients[i].Lat < recipients[j].Lat })
	return recipients
}

func (s *Selector) deduplicate(recipients []network.Node) []network.Node {
	seen := make(map[common.Address]struct{}, len(recipients))
	var selected []network.Node
	for _, node := range recipients {
		if _, exists := seen[node.Addr]; !exists {
			seen[node.Addr] = struct{}{}
			selected = append(selected, node)
		}
	}
	return selected
}

func (s *Selector) selectCloseNodes(committee *types.Committee, clusterID, minNodes, lowLatencyNodes int, exclude []common.Address) []network.Node {
	var selected, candidates []network.Node
	candidates = s.routingCandidatesFromCluster(clusterID, exclude, committee)
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

	for ; i < len(candidates) && candidates[i].Lat < uint(constants.DefaultNearThreshold) && len(selected) < lowLatencyNodes; i++ {
		c := candidates[i]
		if _, exists := seen[c.Addr]; !exists {
			seen[c.Addr] = struct{}{}
			selected = append(selected, candidates[i])
		}
	}

	return selected
}

func (s *Selector) buildResultFromRecipients(recipients []common.Address, clusters network.Clusters) []network.Node {
	result := make([]network.Node, 0, len(recipients))
	for _, recipient := range recipients {
		id := clusters.IDByAddress(recipient)
		member, err := clusters.GetNode(id, recipient)
		if err == nil {
			result = append(result, member)
		}
	}
	return result
}

func (s *Selector) clusterStatus(recipients []network.Node, msg message.Msg, from common.Address, senderType SenderType, ownClusterID int, originClusterID int) {
	logKey := fmt.Sprintf("%d-%d-%d", msg.H(), msg.R(), msg.Code())

	s.heightLock.Lock()
	if _, logged := s.loggedHR[logKey]; logged {
		s.heightLock.Unlock()
		return
	}

	oldHeight := s.recentHeights[s.heightIndex]
	if oldHeight != 0 {
		for k, h := range s.loggedHR {
			if h == oldHeight {
				delete(s.loggedHR, k)
			}
		}
	}

	s.recentHeights[s.heightIndex] = msg.H()
	s.loggedHR[logKey] = msg.H()
	s.heightIndex = (s.heightIndex + 1) % 50
	s.heightLock.Unlock()

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

	clusterMap := make(map[int][]network.Node)
	for _, node := range recipients {
		clusterMap[node.ClusterID] = append(clusterMap[node.ClusterID], node)
	}

	for clusterID, cluster := range clusterMap {
		var lostPeers []string
		connectedCount := 0
		latencyList := []string{}

		for _, peer := range cluster {
			_, ok := s.peerFinder.FindPeer(peer.Addr)
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
