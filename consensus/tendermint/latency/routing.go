package latency

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/latency/ping"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/p2p/enode"
)

type SenderType int

const (
	originator SenderType = iota + 1
	localRelayerOriginCluster
	firstRelayerRemoteCluster
	localRelayerRemoteCluster
)

const (
	// ScaleThresholdForClustering is the minimum number of validators required to do network clustering
	ScaleThresholdForClustering = 10
	DefaultLatency              = 128 // assumed default RTT in ms
	farThreshold                = 100 // indicates cluster is far
	nearThreshold               = 24  // indicates cluster is far
	diversityThreshold          = 50  // ~50ms in uint8 scale, indicates significant distance from lowest latency node

	maxSelectedNodes = 50 // Total cap on selected nodes
)

var errUnknownClusters = errors.New("unknown clustering")

var MeasurementWindow = 2000 // The time window in Millisecond to measure the latency of peers at the beginning of an epoch.

type PeerSelector interface {
	SelectPeersByLatency(committee *types.Committee, msg message.Msg, from common.Address) ([]types.CommitteeMember, error)
}

type Router struct {
	self    common.Address
	nodeKey *ecdsa.PrivateKey

	clusterLock sync.RWMutex
	clusters    Clusters

	broadcaster consensus.Broadcaster

	epochEventChan     chan core.EpochHeadEvent
	epochEventSub      event.Subscription
	chainHeadEventSub  event.Subscription
	chainHeadEventChan chan core.ChainHeadEvent

	committee   []common.Address
	inCommittee bool

	pinger       ping.Pinger
	peerSelector PeerSelector

	cancel context.CancelFunc
	wg     sync.WaitGroup

	latestLatencies map[common.Address]uint
	latencyMu       sync.RWMutex
}

func NewRouter(
	broadcaster consensus.Broadcaster,
	nodeKey *ecdsa.PrivateKey,
	pinger ping.Pinger,
	selector PeerSelector,
) *Router {
	r := &Router{
		broadcaster:        broadcaster,
		nodeKey:            nodeKey,
		epochEventChan:     make(chan core.EpochHeadEvent, 2),
		chainHeadEventChan: make(chan core.ChainHeadEvent, 2),
		pinger:             ping.NewPinger(ping.TCP),
	}
	r.SetDefaultHandlers()

	if pinger != nil {
		r.pinger = pinger
	}

	if selector != nil {
		r.peerSelector = selector
	}
	return r
}

func (r *Router) SetDefaultHandlers() {
	r.peerSelector = &Selector{Router: r, LoggedHR: make(map[string]uint64), RecentHeights: [50]uint64{}, HeightIndex: 0, HeightLock: sync.Mutex{}}
}

func (r *Router) PeerSelector() PeerSelector {
	return r.peerSelector
}

// Exported functions

// Route just select recipients from the clusters, it does not do the message sending.
func (r *Router) Route(committee *types.Committee, msg message.Msg, from common.Address) ([]types.CommitteeMember, error) {
	// no route for small scale network.
	if committee.Len() <= ScaleThresholdForClustering {
		return committee.Members, nil
	}

	return r.PeerSelector().SelectPeersByLatency(committee, msg, from)
}

func (r *Router) Forward(committee *types.Committee, m message.Msg, sender common.Address) {
	recipients, err := r.Route(committee, m, sender)
	if err != nil {
		log.Debug("Forward: No recipients for message from router, broadcast", "error", err, "height", m.H(), "message type", m.Code())
		recipients = committee.Members
	}
	lostPeers := make([]common.Address, 0)
	for _, recipient := range recipients {
		if recipient.Address == sender {
			continue
		}
		if p, ok := r.broadcaster.FindPeer(recipient.Address); ok {
			if p.Cache().Contains(m.Hash()) {
				// This peer had this event, skip it
				continue
			}
			p.Cache().Add(m.Hash(), true)
			go p.SendRaw(message.NetworkCodes[m.Code()], m.Payload()) //nolint
		} else {
			lostPeers = append(lostPeers, recipient.Address)
		}
	}
	if len(lostPeers) > 0 {
		log.Debug("Router: peers not found", "len", len(lostPeers), "peers", lostPeers)
	}
}

func (r *Router) Start(ctx context.Context, chain *core.BlockChain, address common.Address) {
	log.Info("Router: starting latency router")

	curEpoch, err := chain.LatestEpoch()
	if err != nil {
		log.Error("Error fetching latest epoch", "err", err)
		return
	}

	r.epochEventSub = chain.SubscribeEpochHeadEvent(r.epochEventChan)
	r.chainHeadEventSub = chain.SubscribeChainHeadEvent(r.chainHeadEventChan)
	r.self = address

	result := make([]common.Address, curEpoch.Committee.Len())
	for i, member := range curEpoch.Committee.Members {
		result[i] = member.Address
	}
	r.committee = result
	r.updateClusters(NewClusters(result, r.latestLatencies, r.broadcaster, r.self))

	ctx, r.cancel = context.WithCancel(ctx)
	r.wg.Add(1)
	go r.loop(ctx)
}

func (r *Router) Stop() {
	r.cancel()
	r.epochEventSub.Unsubscribe()
	r.wg.Wait()
}

func (r *Router) SetBroadcaster(broadcaster consensus.Broadcaster) {
	r.broadcaster = broadcaster
}

func (r *Router) buildClusters(committee []common.Address) Clusters {
	if len(committee) <= ScaleThresholdForClustering {
		return Clusters{}
	}
	numClusters := int(math.Floor(math.Sqrt(float64(len(committee)))))
	clusterViews := make([]ClusterView, numClusters)
	addressToCluster := make(map[common.Address]int)

	for i, addr := range committee {
		k := i % numClusters
		latency := DefaultLatency
		if lat, ok := r.latestLatencies[addr]; ok {
			latency = int(lat)
		}
		clusterViews[k].Members = append(clusterViews[k].Members, nodeLatency{Addr: addr, Lat: uint(latency)})
		addressToCluster[addr] = k
	}

	return Clusters{
		base:             clusterViews,
		addressToCluster: addressToCluster,
	}
}

func (c *Clusters) ClusterContaining(addr common.Address) int {
	return c.addressToCluster[addr]
}

func (r *Router) refreshClustersLatencies(latMap map[common.Address]uint) {
	r.updateClusters(NewClusters(r.committee, latMap, r.broadcaster, r.self))
}

func (r *Router) measureLatency() error {
	log.Info("Router: measure latency")
	latencyMap, err := r.fetchLatency(r.committee)
	if err != nil {
		log.Error("Router: failed to fetch latency", "err", err)
		return err
	}
	r.refreshClustersLatencies(latencyMap)
	r.latencyMu.Lock()
	defer r.latencyMu.Unlock()
	r.latestLatencies = latencyMap
	return nil
}

func (r *Router) fetchLatency(validators []common.Address) (map[common.Address]uint, error) {
	if r.broadcaster == nil {
		return nil, errors.New("broadcaster not set, can't fetch latency")
	}

	committeeEnodes := r.broadcaster.CommitteeEnodes()

	latency := make(map[common.Address]uint)
	pingTargets := make([]ping.Target, len(validators))

	for i, member := range validators {
		if member == r.self {
			pingTargets[i] = ping.Target{}
			continue
		}
		if memberNode, ok := findByAddress(committeeEnodes, member); ok {
			ip := memberNode.IP()
			port := memberNode.TCP()
			pingTargets[i] = ping.Target{IP: ip.String(), Port: port}
			log.Debug("Router: fetching latency", "targetIP", ip, "targetPort", port)
		} else {
			log.Error("Router: peer not found in broadcaster", "peer", member)
			pingTargets[i] = ping.Target{}
		}
	}

	latencyArray := r.pingPeers(pingTargets)
	for i, addr := range validators {
		// set self latency to 0
		if addr == r.self {
			latency[addr] = 0
		}
		latency[addr] = latencyArray[i]
	}
	return latency, nil
}

func (r *Router) loop(ctx context.Context) {
	defer r.wg.Done()

	ticker := time.NewTicker(5 * time.Minute)
	var cancel context.CancelFunc
	defer func() {
		if cancel != nil {
			cancel()
		}
	}()

	if err := r.measureLatency(); err != nil {
		log.Warn("latency measurement failed", "err", err)
	}
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			// skip measurement if node is not in the committee.
			// there is no broadcaster for unit test context.
			if !r.inCommittee || r.broadcaster == nil {
				continue
			}
			delay := time.Duration(rand.Intn(MeasurementWindow)) * time.Millisecond
			time.Sleep(delay)
			if err := r.measureLatency(); err != nil {
				log.Warn("measureToReport failed", "err", err)
			}
		case epochEv := <-r.epochEventChan:
			log.Info("Router: new epoch detected", "height", epochEv.Header.Number.String())
			epoch := epochEv.Header.Epoch
			r.inCommittee = epoch.Committee.MemberByAddress(r.self) != nil
			if !r.inCommittee {
				log.Info("Router: not in committee, skipping measurement")
				continue
			}
			r.updateCommittee(epoch)
			r.updateClusters(r.buildClusters(r.committee))
			if err := r.measureLatency(); err != nil {
				log.Warn("measureToReport failed", "err", err)
			}
		}
	}
}

func (r *Router) updateCommittee(epoch *types.Epoch) {
	result := make([]common.Address, epoch.Committee.Len())
	for i, member := range epoch.Committee.Members {
		result[i] = member.Address
	}
	r.committee = result
}

func (r *Router) updateClusters(c Clusters) {
	r.clusterLock.Lock()
	defer r.clusterLock.Unlock()
	r.clusters = c
}

func (r *Router) Clusters() Clusters {
	r.clusterLock.RLock()
	defer r.clusterLock.Unlock()
	return r.clusters
}

func (r *Router) pingPeers(targets []ping.Target) []uint {
	channelArray := make([]chan time.Duration, len(targets))
	for i, t := range targets {
		resultCh := make(chan time.Duration, 1)
		if t.IP == "" {
			// default result for non-connected peer to write
			// this should be a reasonable default for max RTT
			resultCh <- time.Second * 5
			channelArray[i] = resultCh
			continue
		}

		r.pinger.Ping(t, resultCh)
		channelArray[i] = resultCh
	}
	results := make([]uint, len(targets))
	for i, resultCh := range channelArray {
		results[i] = uint(<-resultCh)
	}
	return results
}

func findByAddress(committeeEnodes []*enode.Node, addr common.Address) (*enode.Node, bool) {
	for _, memberNode := range committeeEnodes {
		pubKey := memberNode.Pubkey()
		if crypto.PubkeyToAddress(*pubKey) == addr {
			return memberNode, true
		}
	}
	return nil, false
}

type Selector struct {
	*Router
	HeightLock    sync.Mutex
	LoggedHR      map[string]uint64
	RecentHeights [50]uint64
	HeightIndex   int
}

func (s *Selector) clusterStatus(peerCluster [][]common.Address, height uint64, round int64, code uint8, from common.Address, senderType SenderType, ownClusterID int, originClusterID int) {
	logKey := fmt.Sprintf("%d-%d-%d", height, round, code)

	s.HeightLock.Lock()
	if _, logged := s.LoggedHR[logKey]; logged {
		s.HeightLock.Unlock()
		return
	}

	oldHeight := s.RecentHeights[s.HeightIndex]
	if oldHeight != 0 {
		for k, h := range s.LoggedHR {
			if h == oldHeight {
				delete(s.LoggedHR, k)
			}
		}
	}

	s.RecentHeights[s.HeightIndex] = height
	s.LoggedHR[logKey] = height
	s.HeightIndex = (s.HeightIndex + 1) % 50
	s.HeightLock.Unlock()

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

	s.latencyMu.RLock()
	for clusterID, cluster := range peerCluster {
		var lostPeers []string
		connectedCount := 0
		latencies := []string{}

		for _, peer := range cluster {
			_, ok := s.broadcaster.FindPeer(peer)
			if ok {
				connectedCount++
				totalConnected++
				latencies = append(latencies, fmt.Sprintf("%s-%d", peer.Hex(), s.latestLatencies[peer]))
			} else {
				lostPeers = append(lostPeers, peer.Hex())
				totalDisconnected++
			}
			totalSelected++
		}

		if len(lostPeers) == 0 {
			fullyConnectedClusters = append(fullyConnectedClusters,
				fmt.Sprintf("C%d:%d L:%s\n", clusterID, len(cluster), latencies))
		} else {
			sb.WriteString(fmt.Sprintf("Cluster #%d: selected:%d connected:%d\nL:%s\n", clusterID, len(cluster), connectedCount, latencies))
			sb.WriteString("  X disconnected:")
			for _, peerHex := range lostPeers {
				sb.WriteString(" ")
				sb.WriteString(peerHex)
			}
			sb.WriteByte('\n')
		}
	}
	s.latencyMu.RUnlock()

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

func (s *Selector) SelectPeersByLatency(committee *types.Committee, msg message.Msg, from common.Address) ([]types.CommitteeMember, error) {
	clusters := s.Clusters()
	if len(clusters.base) == 0 {
		return committee.Members, nil // Fallback for small networks
	}

	senderClusterID := clusters.clusterContaining(from)
	originClusterID := clusters.clusterContaining(msg.Originator())
	ownClusterID := clusters.clusterContaining(s.self)

	if senderClusterID == -1 || originClusterID == -1 || ownClusterID == -1 {
		log.Error("Router: unknown clusters", "sender", from.Hex(), "originator", msg.Originator().Hex(), "self", s.self.Hex())
		return nil, errUnknownClusters
	}

	// Result tracks selected addresses per cluster for logging
	result := make([][]common.Address, len(clusters.base))

	// Helper to check if an address is in a slice
	containsAddress := func(addrs []common.Address, addr common.Address) bool {
		for _, a := range addrs {
			if a == addr {
				return true
			}
		}
		return false
	}

	type memberWithLatency struct {
		node   nodeLatency
		member types.CommitteeMember
	}
	var recipients []memberWithLatency
	// selectDiverseNodes selects nodes from a cluster based on latency statistics, similar to SelectLatencyDiverseNodes
	selectNodes := func(cluster ClusterView, clusterID int, maxNodes int, exclude []common.Address) []memberWithLatency {
		var selected []memberWithLatency
		members := cluster.Members // Already sorted by latency

		// Filter out excluded addresses and disconnected peers
		var validMembers []memberWithLatency
		for _, node := range members {
			if containsAddress(exclude, node.Addr) {
				continue
			}
			if _, ok := s.broadcaster.FindPeer(node.Addr); ok {
				if member := committee.MemberByAddress(node.Addr); member != nil {
					validMembers = append(validMembers, memberWithLatency{node, *member})
				}
			}
		}

		if len(validMembers) == 0 {
			return nil
		}

		// select first node always - lowest latency
		selected = append(selected, validMembers[0])
		result[clusterID] = append(result[clusterID], validMembers[0].node.Addr)

		// select all close nodes
		for _, vm := range validMembers {
			if len(selected) >= maxNodes {
				return selected
			}
			if vm.node.Lat < uint(nearThreshold) {
				selected = append(selected, vm)
				result[clusterID] = append(result[clusterID], vm.node.Addr)
			}
		}

		// select one diverse node, far from closest node
		// Check if additional diverse node is needed
		if validMembers[0].node.Lat < farThreshold && len(validMembers) > 1 {
			// Find diverse node: highest latency node with Lat >= lowest + diversityThreshold
			for i := len(validMembers) - 1; i >= 1; i-- {
				if validMembers[i].node.Lat >= validMembers[0].node.Lat+diversityThreshold {
					selected = append(selected, validMembers[i])
					result[clusterID] = append(result[clusterID], validMembers[i].node.Addr)
					break
				}
			}
		}
		return selected
	}

	// Select peers based on sender type
	switch {
	case from == s.self: // Originator
		maxRemoteNodes := 3 //early burst selecting more nodes from originator
		for clusterID, cluster := range clusters.base {
			if clusterID == ownClusterID {
				continue
			}
			recipients = append(recipients, selectNodes(cluster, clusterID, maxRemoteNodes, nil)...)
		}
		// Select up to sqrt(N) peers from own cluster
		maxLocalNodes := int(math.Sqrt(float64(len(clusters.base[ownClusterID].Members))))
		recipients = append(recipients, selectNodes(clusters.base[ownClusterID], ownClusterID, maxLocalNodes, []common.Address{s.self})...)
		s.clusterStatus(result, msg.H(), msg.R(), msg.Code(), from, originator, ownClusterID, originClusterID)

	case originClusterID == ownClusterID && from != s.self: // Local relayer in originator's cluster
		maxRemoteNodes := 2 // relayer in
		for clusterID, cluster := range clusters.base {
			if clusterID == ownClusterID {
				continue
			}
			recipients = append(recipients, selectNodes(cluster, clusterID, maxRemoteNodes, nil)...)
		}
		// Select up to sqrt(N) peers from own cluster
		maxLocalNodes := int(math.Sqrt(float64(len(clusters.base[ownClusterID].Members))))
		recipients = append(recipients, selectNodes(clusters.base[ownClusterID], ownClusterID, maxLocalNodes, []common.Address{s.self, from})...)
		s.clusterStatus(result, msg.H(), msg.R(), msg.Code(), from, localRelayerOriginCluster, ownClusterID, originClusterID)

	case senderClusterID != ownClusterID && originClusterID != ownClusterID: // First relayer in remote cluster
		// Select all peers from own cluster
		recipients = append(recipients, selectNodes(clusters.base[ownClusterID], ownClusterID, len(clusters.base[ownClusterID].Members), []common.Address{s.self, from})...)
		s.clusterStatus(result, msg.H(), msg.R(), msg.Code(), from, firstRelayerRemoteCluster, ownClusterID, originClusterID)

	case senderClusterID == ownClusterID && originClusterID != ownClusterID: // Local relayer in remote cluster
		// Select up to sqrt(N) peers from own cluster
		maxLocalNodes := int(math.Sqrt(float64(len(clusters.base[ownClusterID].Members))))
		recipients = append(recipients, selectNodes(clusters.base[ownClusterID], ownClusterID, maxLocalNodes, []common.Address{s.self, from})...)
		s.clusterStatus(result, msg.H(), msg.R(), msg.Code(), from, localRelayerRemoteCluster, ownClusterID, originClusterID)
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
	return selected, nil
}
