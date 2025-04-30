package router

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/log"
)

var errUnknownClusters = errors.New("unknown clustering")

type PeerSelector interface {
	SelectPeers(committee *types.Committee, msg message.Msg, from common.Address) ([]types.CommitteeMember, error)
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

func (r *Selector) SelectPeers(committee *types.Committee, msg message.Msg, from common.Address) ([]types.CommitteeMember, error) {
	clusters := r.Clusters()
	if len(clusters.base) == 0 {
		log.Info("Selector: no clusters, falling back to all committee members")
		return committee.Members, nil
	}

	senderClusterID := clusters.clusterContaining(from)
	originClusterID := clusters.clusterContaining(msg.Originator())
	ownClusterID := clusters.clusterContaining(r.self)

	if senderClusterID == -1 || originClusterID == -1 || ownClusterID == -1 {
		log.Error("Selector: unknown clusters", "sender", from.Hex(), "originator", msg.Originator().Hex(), "msg hash", msg.Hash().Hex(), "self", r.self.Hex())
		return nil, errUnknownClusters
	}

	var senderType SenderType
	switch {
	case from == r.self:
		senderType = originator
	case originClusterID == ownClusterID && from != r.self:
		senderType = localRelayerOriginCluster
	case senderClusterID != ownClusterID && originClusterID != ownClusterID:
		senderType = firstRelayerRemoteCluster
	case senderClusterID == ownClusterID && originClusterID != ownClusterID:
		senderType = localRelayerRemoteCluster
	}

	cacheKey := GenerateCacheKey(from, senderType, msg.Code())
	cached, exists := r.cache.Get(cacheKey)

	if exists {
		allConnected := true
		for _, recipient := range cached.Recipients {
			if _, ok := r.broadcaster.FindPeer(recipient.Address); !ok {
				allConnected = false
				break
			}
		}
		if allConnected {
			r.cache.UpdateLastUsed(cacheKey)
			r.clusterStatus(r.buildResultFromCache(cached.Recipients, clusters), msg.H(), msg.R(), msg.Code(), from, senderType, ownClusterID, originClusterID)
			return cached.Recipients, nil
		}
	}

	result := make([][]common.Address, len(clusters.base))
	containsAddress := func(addrs []common.Address, addr common.Address) bool {
		for _, a := range addrs {
			if a == addr {
				return true
			}
		}
		return false
	}

	type memberWithLatency struct {
		node   NodeLatency
		member types.CommitteeMember
	}
	var recipients []memberWithLatency

	selectNodes := func(cluster ClusterView, clusterID int, maxNodes int, exclude []common.Address) []memberWithLatency {
		var selected, candidates []memberWithLatency
		members := cluster.Members
		for _, node := range members {
			if containsAddress(exclude, node.Addr) {
				continue
			}
			if _, ok := r.broadcaster.FindPeer(node.Addr); ok {
				if member := committee.MemberByAddress(node.Addr); member != nil {
					candidates = append(candidates, memberWithLatency{node, *member})
				}
			}
		}
		if len(candidates) == 0 {
			return nil
		}
		selected = append(selected, candidates[0])
		result[clusterID] = append(result[clusterID], candidates[0].node.Addr)
		i := 1
		//todo(piyush): do we need a cap for close nodes as well?
		for ; i < len(candidates) && candidates[i].node.Lat < uint(nearThreshold); i++ {
			selected = append(selected, candidates[i])
			result[clusterID] = append(result[clusterID], candidates[i].node.Addr)
		}
		for ; i < len(candidates) && len(selected) < maxNodes; i++ {
			selected = append(selected, candidates[i])
			result[clusterID] = append(result[clusterID], candidates[i].node.Addr)
		}
		farthestLat := selected[len(selected)-1].node.Lat
		if farthestLat < uint(farThreshold) {
			for ; i < len(candidates); i++ {
				if candidates[i].node.Lat != DefaultLatency && candidates[i].node.Lat >= farthestLat+diversityThreshold {
					selected = append(selected, candidates[i])
					result[clusterID] = append(result[clusterID], candidates[i].node.Addr)
					break
				}
			}
		}
		return selected
	}

	switch senderType {
	case originator:
		maxRemoteNodes := 2
		for clusterID, cluster := range clusters.base {
			if clusterID == ownClusterID {
				continue
			}
			recipients = append(recipients, selectNodes(cluster, clusterID, maxRemoteNodes, nil)...)
		}
		maxLocalNodes := int(math.Sqrt(float64(len(clusters.base[ownClusterID].Members))))
		recipients = append(recipients, selectNodes(clusters.base[ownClusterID], ownClusterID, maxLocalNodes, []common.Address{r.self})...)

	case localRelayerOriginCluster:
		maxRemoteNodes := 1
		for clusterID, cluster := range clusters.base {
			if clusterID == ownClusterID {
				continue
			}
			recipients = append(recipients, selectNodes(cluster, clusterID, maxRemoteNodes, nil)...)
		}
		maxLocalNodes := int(math.Sqrt(float64(len(clusters.base[ownClusterID].Members))))
		recipients = append(recipients, selectNodes(clusters.base[ownClusterID], ownClusterID, maxLocalNodes, []common.Address{r.self, from})...)

	case firstRelayerRemoteCluster:
		recipients = append(recipients, selectNodes(clusters.base[ownClusterID], ownClusterID, len(clusters.base[ownClusterID].Members), []common.Address{r.self, from})...)

	case localRelayerRemoteCluster:
		maxLocalNodes := int(math.Sqrt(float64(len(clusters.base[ownClusterID].Members))))
		recipients = append(recipients, selectNodes(clusters.base[ownClusterID], ownClusterID, maxLocalNodes, []common.Address{r.self, from})...)
	}

	if len(recipients) == 0 {
		log.Warn("No recipients selected, falling back to all committee members", "height", msg.H(), "round", msg.R(), "code", msg.Code())
		return committee.Members, nil
	}

	sort.Slice(recipients, func(i, j int) bool { return recipients[i].node.Lat < recipients[j].node.Lat })
	selected := make([]types.CommitteeMember, len(recipients))
	for i, r := range recipients {
		selected[i] = r.member
	}

	r.cache.Set(cacheKey, selected)
	r.clusterStatus(result, msg.H(), msg.R(), msg.Code(), from, senderType, ownClusterID, originClusterID)
	return selected, nil
}

func (r *Selector) buildResultFromCache(recipients []types.CommitteeMember, clusters Clusters) [][]common.Address {
	result := make([][]common.Address, len(clusters.base))
	for _, recipient := range recipients {
		clusterID := clusters.clusterContaining(recipient.Address)
		if clusterID >= 0 {
			result[clusterID] = append(result[clusterID], recipient.Address)
		}
	}
	return result
}

func (r *Selector) clusterStatus(peerCluster [][]common.Address, height uint64, round int64, code uint8, from common.Address, senderType SenderType, ownClusterID int, originClusterID int) {
	logKey := fmt.Sprintf("%d-%d-%d", height, round, code)

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

	r.recentHeights[r.heightIndex] = height
	r.loggedHR[logKey] = height
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
	case localRelayerOriginCluster:
		sender = "local relayer origin cluster"
	case firstRelayerRemoteCluster:
		sender = "first relayer remote cluster"
	case localRelayerRemoteCluster:
		sender = "local relayer remote cluster"
	}

	msgType := "proposal"
	switch code {
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

	sb.WriteString(fmt.Sprintf("\nCluster routing status:\t Height=%d, Round=%d, From=%s Message=%s SenderType=%s localCluster=%d originCluster=%d\n",
		height, round, from.Hex(), msgType, sender, ownClusterID, originClusterID))

	latencies := r.Latencies()
	for clusterID, cluster := range peerCluster {
		var lostPeers []string
		connectedCount := 0
		latencyList := []string{}

		for _, peer := range cluster {
			_, ok := r.broadcaster.FindPeer(peer)
			if ok {
				connectedCount++
				totalConnected++
				latencyList = append(latencyList, fmt.Sprintf("%s-%d", peer.Hex(), latencies[peer]))
			} else {
				lostPeers = append(lostPeers, peer.Hex())
				totalDisconnected++
			}
			totalSelected++
		}

		if len(lostPeers) == 0 {
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
