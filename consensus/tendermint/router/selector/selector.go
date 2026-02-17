package selector

import (
	"context"
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
	"github.com/autonity/autonity/consensus/tendermint/router/cluster"
	"github.com/autonity/autonity/consensus/tendermint/router/interfaces"
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

const defaultNearThreshold = 50

type selector struct {
	clustersProvider interfaces.ClustersProvider
	recipientCache   cache.Recipients
	peerFinder       interfaces.PeerFinder
	heightLock       sync.Mutex
	loggedHR         map[string]uint64
	recentHeights    [50]uint64
	heightIndex      int
}

func New(np interfaces.ClustersProvider, cache cache.Recipients) interfaces.PeerSelector {
	s := &selector{
		clustersProvider: np,
		recipientCache:   cache,
		loggedHR:         make(map[string]uint64),
		recentHeights:    [50]uint64{},
		heightIndex:      0,
	}
	return s
}

func (s *selector) SetBroadcaster(broadcaster interfaces.PeerFinder) {
	s.peerFinder = broadcaster
}

func (s *selector) routingCandidatesFromCluster(clusterID int, self common.Address, committee *types.Committee) []cluster.Node {
	clusters := s.clustersProvider.Clusters()
	members := clusters.MembersByID(clusterID)
	candidates := make([]cluster.Node, 0, len(members))
	for _, node := range members {
		if self == node.Addr {
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

func (s *selector) allConnected(recipients []common.Address) bool {
	for _, recipient := range recipients {
		if _, ok := s.peerFinder.FindPeer(recipient); !ok {
			return false
		}
	}
	return true
}

func (s *selector) SelectPeers(committee *types.Committee, msg message.Msg, from common.Address) ([]common.Address, error) {
	return s.selectPeersWithBuckets(committee, msg, from, msg.Code() == message.ProposalCode)
}

// determines the routing base of a message
func routingBase(committee *types.Committee, msg message.Msg) common.Address {
	var rb common.Address
	switch m := msg.(type) {
	case *message.Propose:
		rb = m.Signer()
	case *message.Prevote:
		routingBaseIndex := m.Signers().RightmostSigner()
		rb = committee.Members[routingBaseIndex].Address
	case *message.Precommit:
		routingBaseIndex := m.Signers().RightmostSigner()
		rb = committee.Members[routingBaseIndex].Address
	default:
		panic("unknown msg type. msg: " + msg.String())
	}
	return rb
}

func (s *selector) selectPeersWithBuckets(committee *types.Committee, msg message.Msg, from common.Address, isProposal bool) ([]common.Address, error) {
	clusters := s.clustersProvider.Clusters()
	if len(clusters.Base()) == 0 {
		return []common.Address{}, errors.New("no clusters")
	}

	rb := routingBase(committee, msg)
	senderClusterID := clusters.IDByAddress(from)
	originClusterID := clusters.IDByAddress(rb)
	ownClusterID := clusters.ID()

	if senderClusterID == -1 || originClusterID == -1 || ownClusterID == -1 {
		log.Debug("selector: unknown clusters", "sender", from.Hex(), "routingBase", rb.Hex(), "msg hash", msg.Hash().Hex(), "self", clusters.Self().Hex())
		return nil, errors.New("unknown clusters")
	}

	excludeSender := func(recipients []common.Address, sender common.Address) []common.Address {
		ret := make([]common.Address, 0, len(recipients))
		for _, peer := range recipients {
			if peer == sender {
				continue
			}
			ret = append(ret, peer)
		}
		return ret
	}

	senderType := determineSenderType(from, clusters.Self(), rb, originClusterID, ownClusterID, senderClusterID)
	cacheKey := cache.GenerateKey(int(senderType), isProposal)
	if cached, exists := s.recipientCache.Get(cacheKey); exists {
		if s.allConnected(cached.Recipients) {
			s.recipientCache.UpdateLastUsed(cacheKey)
			selected := excludeSender(cached.Recipients, from)
			s.clusterStatus(s.buildResultFromRecipients(selected, clusters), msg, from, senderType, ownClusterID, originClusterID)
			return selected, nil
		}
	}
	var recipients []cluster.Node
	if isProposal {
		recipients = s.selectNodesForProposal(clusters, committee, senderType, ownClusterID)
	} else {
		recipients = s.selectNodesForNonProposalMessages(clusters, committee, senderType, ownClusterID)
	}
	recipientsAddr := make([]common.Address, 0, len(recipients))
	for _, r := range recipients {
		recipientsAddr = append(recipientsAddr, r.Addr)
	}

	s.recipientCache.Set(cacheKey, recipientsAddr)
	selected := excludeSender(recipientsAddr, from)
	s.clusterStatus(s.buildResultFromRecipients(selected, clusters), msg, from, senderType, ownClusterID, originClusterID)
	return selected, nil
}

func (s *selector) selectNodesByLatencySpread() []cluster.Node {
	recipients := make([]cluster.Node, 0)
	usedClusters := make(map[int]bool)

	clusters := s.clustersProvider.Clusters()

	// Select one connected node from each remote cluster
	for _, nodes := range clusters.BucketNodes() {
		for _, node := range nodes {
			if usedClusters[node.ClusterID] || node.ClusterID == clusters.ID() {
				continue
			}
			if _, ok := s.peerFinder.FindPeer(node.Addr); ok {
				recipients = append(recipients, node)
				usedClusters[node.ClusterID] = true
				break // Take the first connected node (primary or fallback)
			}
		}
	}

	// fallback, optimal node by latency spread was not found pick another
	for clusterID, nodes := range clusters.Base() {
		if clusterID == clusters.ID() || usedClusters[clusterID] {
			continue // Skip local cluster and already selected clusters
		}
		for _, node := range nodes {
			if _, ok := s.peerFinder.FindPeer(node.Addr); ok {
				recipients = append(recipients, node)
				usedClusters[node.ClusterID] = true
				break // Take the first connected node
			}
		}
	}
	return recipients
}

func determineSenderType(from, self, routingBase common.Address, originClusterID, ownClusterID, senderClusterID int) SenderType {
	switch {
	case from == self:
		return originator
	case from == routingBase && originClusterID == ownClusterID:
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

func (s *selector) selectNodesForNonProposalMessages(clusters cluster.Clusters, committee *types.Committee, senderType SenderType, ownClusterID int) []cluster.Node {
	recipients := make([]cluster.Node, 0, len(clusters.Base())*2)
	var minNodes, lowLatencyNodes int

	switch senderType {
	case originator:
		lowLatencyNodes = 0
		for clusterID := range clusters.Base() {
			if clusterID == ownClusterID {
				//local cluster
				localNodes := int(math.Sqrt(float64(len(clusters.Base()[ownClusterID]))))
				recipients = append(recipients, s.selectCloseNodes(committee, clusterID, localNodes, lowLatencyNodes, clusters.Self())...)
			} else {
				// everyone in remote clusters
				minNodes = math.MaxInt
				recipients = append(recipients, s.selectCloseNodes(committee, clusterID, minNodes, lowLatencyNodes, clusters.Self())...)
			}
		}
	case firstRelayerOriginCluster:
		minNodes = 1 // 1 node from each remote cluster randomly selected
		for clusterID := range clusters.Base() {
			recipients = append(recipients, s.selectRandomNodes(committee, clusterID, minNodes, clusters.Self())...)
		}
		minNodes = len(clusters.MembersByID(clusters.ID()))
		recipients = append(recipients, s.selectCloseNodes(committee, ownClusterID, minNodes, 0, clusters.Self())...)

	case firstRelayerRemoteCluster, localRelayerOriginCluster, localRelayerRemoteCluster: // 20 nodes
		minNodes = len(clusters.Base()[ownClusterID])
		recipients = append(recipients, s.selectCloseNodes(committee, ownClusterID, minNodes, 0, clusters.Self())...)
	}

	recipients = s.deduplicate(recipients)
	// Sort by latency for consistent ordering
	sort.Slice(recipients, func(i, j int) bool { return recipients[i].Lat < recipients[j].Lat })
	return recipients
}

func (s *selector) selectNodesForProposal(clusters cluster.Clusters, committee *types.Committee, senderType SenderType, ownClusterID int) []cluster.Node {
	recipients := make([]cluster.Node, 0, len(clusters.Base())*2)
	var minNodes, lowLatencyNodes int
	switch senderType {
	case originator:
		recipients = append(recipients, s.selectNodesByLatencySpread()...)
		lowLatencyNodes = 0
		for clusterID := range clusters.Base() {
			if clusterID == ownClusterID {
				minNodes = int(math.Sqrt(float64(len(clusters.Base()[ownClusterID]))))
				recipients = append(recipients, s.selectCloseNodes(committee, clusterID, minNodes, lowLatencyNodes, clusters.Self())...)
			} else {
				minNodes = 1
				recipients = append(recipients, s.selectCloseNodes(committee, clusterID, minNodes, lowLatencyNodes, clusters.Self())...)
			}
		}
	case firstRelayerOriginCluster:
		minNodes = 1
		lowLatencyNodes = 4
		for clusterID := range clusters.Base() {
			if clusterID == ownClusterID {
				// local cluster full
				minNodes = len(clusters.MembersByID(clusters.ID()))
				recipients = append(recipients, s.selectCloseNodes(committee, ownClusterID, minNodes, lowLatencyNodes, clusters.Self())...)
			} else {
				recipients = append(recipients, s.selectCloseNodes(committee, clusterID, minNodes, lowLatencyNodes, clusters.Self())...)
			}
		}

	case firstRelayerRemoteCluster: // now also includes messages from the first Relayer in origin cluster
		lowLatencyNodes = 4
		for clusterID := range clusters.Base() {
			if clusterID == ownClusterID {
				// local cluster full
				minNodes = len(clusters.Base()[ownClusterID])
				recipients = append(recipients, s.selectCloseNodes(committee, ownClusterID, minNodes, lowLatencyNodes, clusters.Self())...)
			} else {
				minNodes = 0
				// only low latency nodes from remote clusters
				recipients = append(recipients, s.selectCloseNodes(committee, clusterID, minNodes, lowLatencyNodes, clusters.Self())...)
			}
		}
	case localRelayerOriginCluster, localRelayerRemoteCluster: // 20 nodes
		minNodes = len(clusters.Base()[ownClusterID])
		recipients = append(recipients, s.selectCloseNodes(committee, ownClusterID, minNodes, lowLatencyNodes, clusters.Self())...)
	}

	recipients = s.deduplicate(recipients)
	// Sort by latency for consistent ordering
	sort.Slice(recipients, func(i, j int) bool { return recipients[i].Lat < recipients[j].Lat })
	return recipients
}

func (s *selector) deduplicate(recipients []cluster.Node) []cluster.Node {
	seen := make(map[common.Address]struct{}, len(recipients))
	var selected []cluster.Node
	for _, node := range recipients {
		if _, exists := seen[node.Addr]; !exists {
			seen[node.Addr] = struct{}{}
			selected = append(selected, node)
		}
	}
	return selected
}

func (s *selector) selectCloseNodes(committee *types.Committee, clusterID, minNodes, lowLatencyNodes int, self common.Address) []cluster.Node {
	var selected, candidates []cluster.Node
	candidates = s.routingCandidatesFromCluster(clusterID, self, committee)
	if len(candidates) == 0 {
		return nil
	}
	i := 0
	for ; i < len(candidates) && len(selected) < minNodes; i++ {
		selected = append(selected, candidates[i])
	}

	for ; i < len(candidates) && candidates[i].Lat < uint(defaultNearThreshold) && len(selected) < lowLatencyNodes; i++ {
		selected = append(selected, candidates[i])
	}

	return selected
}

func (s *selector) selectRandomNodes(committee *types.Committee, clusterID, minNodes int, self common.Address) []cluster.Node {
	var selected, candidates []cluster.Node
	candidates = s.routingCandidatesFromCluster(clusterID, self, committee)
	if len(candidates) == 0 {
		return nil
	}

	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	for i := 0; i < len(candidates) && len(selected) < minNodes; i++ {
		selected = append(selected, candidates[i])
	}
	return selected
}

func (s *selector) buildResultFromRecipients(recipients []common.Address, clusters cluster.Clusters) []cluster.Node {
	result := make([]cluster.Node, 0, len(recipients))
	for _, recipient := range recipients {
		id := clusters.IDByAddress(recipient)
		member, err := clusters.GetNode(id, recipient)
		if err == nil {
			result = append(result, member)
		}
	}
	return result
}

func (s *selector) clusterStatus(recipients []cluster.Node, msg message.Msg, from common.Address, senderType SenderType, ownClusterID int, originClusterID int) {
	if !log.Root().Enabled(context.Background(), log.LevelDebug) {
		return
	}

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
	var sender string
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

	var msgType string
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

	clusterMap := make(map[int][]cluster.Node)
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
	log.Debug("cluster status", "stats", sb.String())
}
