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

type SenderType int

const (
	originator SenderType = iota + 1
	firstRelayerOriginCluster
	localRelayerOriginCluster
	firstRelayerRemoteCluster
	localRelayerRemoteCluster
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

func (r *Selector) SelectPeers(committee *types.Committee, msg message.Msg, from common.Address) ([]common.Address, error) {
	clusters := r.Clusters(msg.H())
	if len(clusters.base) == 0 {
		log.Info("Selector: no clusters, falling back to all committee members")
		return r.committeeAddresses(committee), nil
	}

	senderClusterID := clusters.clusterContaining(from)
	originClusterID := clusters.clusterContaining(msg.Originator())
	ownClusterID := clusters.clusterContaining(r.self)

	if senderClusterID == -1 || originClusterID == -1 || ownClusterID == -1 {
		log.Error(
			"Selector: unknown clusters",
			"sender",
			from.Hex(),
			"originator",
			msg.Originator().Hex(),
			"msgHash",
			msg.Hash().Hex(),
			"self",
			r.self.Hex(),
			"senderCluster",
			senderClusterID,
			"originCluster",
			originClusterID,
			"ownCluster",
			ownClusterID,
			"height",
			msg.H(),
			"clusterId",
			clusters.id,
		)
		return nil, errUnknownClusters
	}

	var senderType SenderType
	switch {
	case from == r.self:
		senderType = originator
	case from == msg.Originator():
		senderType = firstRelayerOriginCluster
	case originClusterID == ownClusterID && from != r.self:
		senderType = localRelayerOriginCluster
	case senderClusterID != ownClusterID && originClusterID != ownClusterID:
		senderType = firstRelayerRemoteCluster
	case senderClusterID == ownClusterID && originClusterID != ownClusterID:
		senderType = localRelayerRemoteCluster
	}

	cacheKey := GenerateCacheKey(from, senderType, msg.Code(), clusters.id)
	cached, exists := r.cache.Get(cacheKey)

	if exists {
		allConnected := true
		for _, recipient := range cached.Recipients {
			if _, ok := r.broadcaster.FindPeer(recipient); !ok {
				allConnected = false
				break
			}
		}
		if allConnected {
			r.cache.UpdateLastUsed(cacheKey)
			r.clusterStatus(r.buildResultFromCache(cached.Recipients, clusters), msg, from, senderType, ownClusterID, originClusterID, clusters.id)
			return cached.Recipients, nil
		}
	}

	result := make([][]NodeLatency, len(clusters.base))
	containsAddress := func(addrs []common.Address, addr common.Address) bool {
		for _, a := range addrs {
			if a == addr {
				return true
			}
		}
		return false
	}

	var recipients []NodeLatency

	selectNodes := func(cluster ClusterView, clusterID int, maxNodes int, exclude []common.Address) []NodeLatency {
		var selected, candidates []NodeLatency
		if maxNodes < 1 {
			return selected
		}
		members := cluster.Members
		for _, node := range members {
			if containsAddress(exclude, node.Addr) {
				continue
			}
			if _, ok := r.broadcaster.FindPeer(node.Addr); ok {
				if member := committee.MemberByAddress(node.Addr); member != nil {
					candidates = append(candidates, node)
				}
			}
		}
		if len(candidates) == 0 {
			return nil
		}
		maxLowLatNodes := 3
		selected = append(selected, candidates[0])
		result[clusterID] = append(result[clusterID], candidates[0])
		i := 1
		for ; i < len(candidates) && candidates[i].Lat < uint(nearThreshold) && len(selected) < maxLowLatNodes; i++ {
			selected = append(selected, candidates[i])
			result[clusterID] = append(result[clusterID], candidates[i])
		}
		for ; i < len(candidates) && len(selected) < maxNodes-1; i++ { // leave room for one far node
			selected = append(selected, candidates[i])
			result[clusterID] = append(result[clusterID], candidates[i])
		}
		farthestLat := selected[len(selected)-1].Lat
		if farthestLat < uint(farThreshold) && len(selected) < maxNodes {
			for ; i < len(candidates); i++ {
				if candidates[i].Lat != DefaultLatency && candidates[i].Lat >= farthestLat+diversityThreshold {
					selected = append(selected, candidates[i])
					result[clusterID] = append(result[clusterID], candidates[i])
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

	case firstRelayerOriginCluster:
		maxRemoteNodes := 1
		for clusterID, cluster := range clusters.base {
			if clusterID == ownClusterID {
				continue
			}
			recipients = append(recipients, selectNodes(cluster, clusterID, maxRemoteNodes, nil)...)
		}
		maxLocalNodes := int(math.Sqrt(float64(len(clusters.base[ownClusterID].Members))))
		recipients = append(recipients, selectNodes(clusters.base[ownClusterID], ownClusterID, maxLocalNodes, []common.Address{r.self, from})...)

	case localRelayerOriginCluster:
		maxLocalNodes := int(math.Sqrt(float64(len(clusters.base[ownClusterID].Members))))
		recipients = append(recipients, selectNodes(clusters.base[ownClusterID], ownClusterID, maxLocalNodes, []common.Address{r.self, from})...)

	case firstRelayerRemoteCluster:
		recipients = append(recipients, selectNodes(clusters.base[ownClusterID], ownClusterID, len(clusters.base[ownClusterID].Members), []common.Address{r.self, from})...)

	case localRelayerRemoteCluster:
		// this is probably not needed but to complete dissemination we keep it
		maxLocalNodes := int(math.Min(2, math.Sqrt(float64(len(clusters.base[ownClusterID].Members)))))
		recipients = append(recipients, selectNodes(clusters.base[ownClusterID], ownClusterID, maxLocalNodes, []common.Address{r.self, from})...)
	}

	sort.Slice(recipients, func(i, j int) bool { return recipients[i].Lat < recipients[j].Lat })
	selected := make([]common.Address, len(recipients))
	for i, r := range recipients {
		selected[i] = r.Addr
	}

	r.cache.Set(cacheKey, selected)
	r.clusterStatus(result, msg, from, senderType, ownClusterID, originClusterID, clusters.id)
	return selected, nil
}

func (r *Selector) buildResultFromCache(recipients []common.Address, clusters Clusters) [][]NodeLatency {
	result := make([][]NodeLatency, len(clusters.base))
	for _, recipient := range recipients {
		id := clusters.clusterContaining(recipient)
		member, err := clusters.addressToMember(id, recipient)
		if err == nil {
			result[id] = append(result[id], member)
		}
	}
	return result
}

func (r *Selector) clusterStatus(
	peerCluster [][]NodeLatency,
	msg message.Msg,
	from common.Address,
	senderType SenderType,
	ownClusterID int,
	originClusterID int,
	clustersId string,
) {
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

	for clusterID, cluster := range peerCluster {
		var lostPeers []string
		connectedCount := 0
		var latencyList []string

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

		if len(lostPeers) == 0 {
			if len(cluster) > 0 {
				fullyConnectedClusters = append(fullyConnectedClusters, fmt.Sprintf("C%d:%d L:%s\n", clusterID, len(cluster), latencyList))
			}
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

	sb.WriteString(fmt.Sprintf("Total: selected:%d connected:%d disconnected:%d clusterId: %s\n", totalSelected, totalConnected, totalDisconnected, clustersId))

	log.Info(sb.String())
}
