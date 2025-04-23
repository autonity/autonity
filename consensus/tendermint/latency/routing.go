package latency

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/autonity/autonity/consensus/tendermint/bft"

	"github.com/autonity/autonity/autonity"
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

var errTooOldMessage = errors.New("too old message")
var errUnknownClusters = errors.New("unknown clustering")

// ScaleThresholdForClustering is the minimum number of validators required to do network clustering
var ScaleThresholdForClustering = 10 // by according to the simulation and testing, there was minimal difference in performance when the number of validators was < 32.

// VerticalRelayingRedundancy is the number of relayers of each cluster to receive the original sender's message.
var VerticalRelayingRedundancy = 2

// HorizontalRelayingRedundancy is the number of relayers of other clusters to receive the relayer's message.
var HorizontalRelayingRedundancy = 1

var MeasurementWindow = 2000 // The time window in Millisecond to measure the latency of peers at the beginning of an epoch.

type PeerSelector interface {
	SelectPeers(committee *types.Committee, msg message.Msg, from common.Address) ([]types.CommitteeMember, error)
}

type Router struct {
	self    common.Address
	nodeKey *ecdsa.PrivateKey

	mu                     sync.RWMutex
	epochDefaultClusters   *Clusters
	epochOptimizedClusters *Clusters
	lastEpochClusters      *Clusters

	broadcaster consensus.Broadcaster
	contracts   *autonity.ProtocolContracts
	reporter    *Reporter

	lastMeasuredEpoch *big.Int

	optimizationEventChan chan *autonity.LatencyKMOptimization
	optimizationEventSub  event.Subscription

	epochEventChan chan core.EpochHeadEvent
	epochEventSub  event.Subscription

	curEpochInfo *types.EpochInfo

	pinger       ping.Pinger
	peerSelector PeerSelector

	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewRouter(
	broadcaster consensus.Broadcaster,
	nodeKey *ecdsa.PrivateKey,
	pinger ping.Pinger,
	selector PeerSelector,
) *Router {
	r := &Router{
		broadcaster:           broadcaster,
		nodeKey:               nodeKey,
		epochEventChan:        make(chan core.EpochHeadEvent, 2),
		optimizationEventChan: make(chan *autonity.LatencyKMOptimization, 2),
		pinger:                ping.NewPinger(ping.TCP),
		lastMeasuredEpoch:     new(big.Int).SetInt64(-1),
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

	return r.PeerSelector().SelectPeers(committee, msg, from)
}

func (r *Router) Forward(committee *types.Committee, m message.Msg, sender common.Address) {
	recipients, err := r.Route(committee, m, sender)
	if err != nil {
		//if !errors.Is(err, consensus.ErrFutureEpochMessage) {
		log.Debug("Forward: No recipients for message from router, broadcast", "error", err, "height", m.H(), "message type", m.Code())
		//	return
		//}
		// forward to all the committee members if the router cannot resolve recipients.
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

func (r *Router) Start(ctx context.Context, chain *core.BlockChain) {
	log.Info("Router: starting latency router")
	optimizationEventSub, err := chain.ProtocolContracts().Latency.WatchKMOptimization(nil, r.optimizationEventChan)
	if err != nil {
		log.Error("Error starting latency router for clustering view optimization", err)
		return
	}

	curEpoch, err := chain.LatestEpoch()
	if err != nil {
		log.Error("Error fetching latest epoch", "err", err)
		return
	}
	r.curEpochInfo = curEpoch

	r.optimizationEventSub = optimizationEventSub
	r.epochEventSub = chain.SubscribeEpochHeadEvent(r.epochEventChan)
	r.contracts = chain.ProtocolContracts()
	r.reporter, err = NewReporter(chain.Config().ChainID, r.nodeKey, r.contracts)
	if err != nil {
		log.Error("failed to create reporter", "err", err)
		return
	}
	r.self = r.reporter.txOpts.From

	// set default clusters for current epoch.
	r.setDefaultCluster(r.buildDefaultClusters(func() []common.Address {
		result := make([]common.Address, r.curEpochInfo.Committee.Len())
		for i, member := range r.curEpochInfo.Committee.Members {
			result[i] = member.Address
		}
		return result
	}()))

	// As from here, we already subscribe the optimization event, however if the optimization was already happened,
	// we'd need to set optimized clusters for current epoch if it was happened.
	optimizationHeight, err := r.contracts.GetMetricsStatus(nil)
	if err != nil {
		log.Error("failed to get optimized clusters height", "err", err)
		return
	}

	if optimizationHeight.Cmp(common.Big0) > 0 && optimizationHeight.Cmp(r.curEpochInfo.NextEpochBlock) < 0 {
		err = r.optimizeCluster(optimizationHeight.Uint64())
		if err != nil {
			log.Error("Router: failed to optimize the clustering", "err", err)
			return
		}
	}

	ctx, r.cancel = context.WithCancel(ctx)
	r.wg.Add(1)
	go r.loop(ctx)
}

func (r *Router) Stop() {
	r.cancel()
	r.epochEventSub.Unsubscribe()
	r.optimizationEventSub.Unsubscribe()
	r.wg.Wait()
}

func (r *Router) SetBroadcaster(broadcaster consensus.Broadcaster) {
	r.broadcaster = broadcaster
}

// Internal package functions
func (r *Router) setDefaultCluster(c *Clusters) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.epochDefaultClusters = c
}

func (r *Router) setOptimizedCluster(c *Clusters) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.epochOptimizedClusters = c
}

func (r *Router) viewRotation(newDefaultCluster *Clusters) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.epochOptimizedClusters != nil {
		r.lastEpochClusters = r.epochOptimizedClusters
	} else {
		r.lastEpochClusters = r.epochDefaultClusters
	}

	r.epochDefaultClusters = newDefaultCluster
	r.epochOptimizedClusters = nil
}

func (r *Router) resolveClusters(h uint64) (*Clusters, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	//todo: in a single epoch there could be multiple cluster views, we need to use the view according to the height.
	if r.epochDefaultClusters != nil && h >= r.epochDefaultClusters.nextEpochHeight {
		log.Warn("returning future epoch message", "height", h, "nextEpochHeight", r.epochDefaultClusters.nextEpochHeight, "cluster", r.epochDefaultClusters)
		return nil, consensus.ErrFutureEpochMessage
	}

	// always try to pick the optimized one 1st
	if r.epochOptimizedClusters != nil && h >= r.epochOptimizedClusters.activatedHeight &&
		r.epochDefaultClusters != nil && h < r.epochDefaultClusters.nextEpochHeight {
		return r.epochOptimizedClusters, nil
	}

	// otherwise, try to pick the default clusters.
	if r.epochDefaultClusters != nil && h >= r.epochDefaultClusters.activatedHeight && h < r.epochDefaultClusters.nextEpochHeight {
		return r.epochDefaultClusters, nil
	}

	// edge case, around epoch rotation, some message might be from past epoch.
	if r.lastEpochClusters != nil && h >= r.lastEpochClusters.activatedHeight && h < r.lastEpochClusters.nextEpochHeight {
		return r.lastEpochClusters, nil
	}

	if r.lastEpochClusters != nil && h < r.lastEpochClusters.activatedHeight {
		return nil, errTooOldMessage
	}

	return nil, errUnknownClusters
}

// buildDefaultClusters partitions the committee into default clusters
func (r *Router) buildDefaultClusters(committee []common.Address) *Clusters {
	if len(committee) <= ScaleThresholdForClustering {
		return nil
	}

	numClusters := int(math.Floor(math.Sqrt(float64(len(committee)))))
	clusters := make([][]common.Address, numClusters)
	for i, addr := range committee {
		//k := int(math.Min(float64(i/numClusters), float64(numClusters-1)))
		k := i % numClusters
		clusters[k] = append(clusters[k], addr)
	}

	defaultClusters := NewCluster(r.curEpochInfo.EpochBlock.Uint64(), r.curEpochInfo.NextEpochBlock.Uint64(), clusters)

	//todo: use the latest(last epoch) optimized to cluster to build the default cluster, starting from epoch 1
	log.Debug("Router: set default clusters", "clusters", func() [][]int {
		clusterInts := make([][]int, len(clusters))
		for i, cluster := range clusters {
			for _, member := range cluster {
				clusterInts[i] = append(clusterInts[i], slices.Index(committee, member))
			}
		}
		return clusterInts
	}(), "height", r.curEpochInfo.EpochBlock.Uint64())
	return defaultClusters
}

func (r *Router) startMeasurementTask(ctx context.Context) (cancel context.CancelFunc) {
	ctx, cancel = context.WithCancel(ctx)
	go func() {
		// todo: the random delay is used to distribute the load of ping messages into a certain period.
		//  2 seconds distribution period might be too short for a large scale network.
		// this is probably not required, but we reduce the measurement window to 2 seconds
		delay := time.Duration(rand.Intn(MeasurementWindow)) * time.Millisecond
		select {
		case <-time.After(delay):
			if err := r.measureToReport(); err != nil {
				log.Warn("measureToReport failed", "err", err)
			}
		case <-ctx.Done():
			log.Info("Measurement task stopped by context cancellation")
		}
	}()
	return cancel
}

func (r *Router) measureToReport() error {
	committee, err := r.contracts.Latency.GetCommittee(nil)
	if err != nil {
		return err
	}
	latencyVec, err := r.fetchLatency(committee)
	if err != nil {
		return err
	}

	err = r.reporter.ReportLatency(committee, latencyVec)
	if err == nil {
		var sb strings.Builder
		sb.WriteString("\nRouter: latency reported!!\n")
		for addr, lat := range latencyVec {
			sb.WriteString(fmt.Sprintf("[%s → %dms]\n", addr.Hex(), lat))
		}
		sb.WriteString("\n")
		log.Info(sb.String())
	}
	return err
}

func (r *Router) fetchLatency(validators []common.Address) (map[common.Address]uint8, error) {
	if r.broadcaster == nil {
		return nil, errors.New("broadcaster not set, can't fetch latency")
	}

	committeeEnodes := r.broadcaster.CommitteeEnodes()

	latency := make(map[common.Address]uint8)
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
		// set self latency to 1, as 0 is the default value from the storage DB.
		// those node who did not measure the latencies, their data will be filled with 0 by default.
		if addr == r.self {
			latency[addr] = 1
		}
		latency[addr] = latencyArray[i]
	}
	return latency, nil
}

func (r *Router) loop(ctx context.Context) {
	defer r.wg.Done()

	ticker := time.NewTicker(10 * time.Second)
	var cancel context.CancelFunc
	defer func() {
		if cancel != nil {
			cancel()
		}
	}()

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			// skip measurement if we have a small scale network or the broadcaster were not be created.
			if r.curEpochInfo.Committee.Len() <= ScaleThresholdForClustering || r.broadcaster == nil {
				continue
			}

			// skip measurement if node is not in the committee.
			member := r.curEpochInfo.Committee.MemberByAddress(r.self)
			if member == nil {
				continue
			}

			// if node reported, skip.
			reported, err := r.contracts.Latency.ClientReported(nil, new(big.Int).SetUint64(member.Index))
			if err != nil || reported {
				continue
			}

			// otherwise, we trigger the measurement once we have quorum peers connected.
			connectedPeers := r.broadcaster.FindPeers(func() []common.Address {
				result := make([]common.Address, r.curEpochInfo.Committee.Len())
				for i, member := range r.curEpochInfo.Committee.Members {
					result[i] = member.Address
				}
				return result
			}())

			quorum := bft.Quorum(new(big.Int).SetInt64(int64(r.curEpochInfo.Committee.Len())))
			// count local client itself with connected peers.
			if int64(len(connectedPeers)+1) >= quorum.Int64() {
				cancel = r.startMeasurementTask(ctx)
			}

		case optimizationEv := <-r.optimizationEventChan:
			log.Info("Router: ready to optimize the clustering", "new view will be applied at height", optimizationEv.Height)
			if optimizationEv.Height.Cmp(r.curEpochInfo.NextEpochBlock) >= 0 {
				log.Info("Router: skip to optimize the clustering", "activation height over epoch", optimizationEv.Height.Uint64())
				continue
			}

			err := r.optimizeCluster(optimizationEv.Height.Uint64())
			if err != nil {
				log.Error("Router: failed to optimize the clustering", "err", err)
			}

		case epochEv := <-r.epochEventChan:
			log.Info("Router: new epoch detected", "height", epochEv.Header.Number.String())
			r.curEpochInfo = &types.EpochInfo{
				Epoch:      *epochEv.Header.Epoch.Copy(),
				EpochBlock: epochEv.Header.Number,
			}

			r.viewRotation(r.buildDefaultClusters(func() []common.Address {
				result := make([]common.Address, r.curEpochInfo.Committee.Len())
				for i, member := range r.curEpochInfo.Committee.Members {
					result[i] = member.Address
				}
				return result
			}()))
			// we cannot trigger measurement at epoch rotation immediately since members need time
			// to create connections with new members, and for new members they need more time to create full mesh
			// connectivity with other members.
		}
	}
}

func (r *Router) optimizeCluster(h uint64) error {
	committee, err := r.contracts.Latency.GetCommittee(nil)
	if err != nil {
		log.Error("Router: optimizeCluster fetch committee", "err", err)
		return err
	}

	// todo: it would cause coherence issue on the clustering view if the reorg are too frequent. But eventually,
	//  they should be consistent.
	// as reading the entire matrix could be reverted due to too much gas consumption, we have to read row by row.
	latencyMat := make([][]uint8, len(committee))
	index := new(big.Int).SetUint64(0)
	for i, _ := range committee {
		_, latency, err := r.contracts.Latency.ReadReport(nil, index.SetInt64(int64(i)))
		if err != nil {
			log.Error("Router: optimizeCluster failed to read latency", "err", err)
			return err
		}

		latencyMat[i] = latency
	}

	optimizedClusters, err := AssignClusters(h, r.curEpochInfo.NextEpochBlock.Uint64(), committee, latencyMat, int(math.Floor(math.Sqrt(float64(len(committee))))))
	if err != nil {
		return err
	}

	r.setOptimizedCluster(optimizedClusters)

	log.Debug(
		"Router: optimizeClusters",
		"optimizedClusters",
		func() [][]int {
			clusterInts := make([][]int, len(optimizedClusters.base))
			for i, cluster := range optimizedClusters.base {
				clusterInts[i] = make([]int, len(cluster))
				for j, member := range cluster {
					clusterInts[i][j] = slices.Index(committee, member)
				}
			}
			return clusterInts
		}(),
		"height",
		h,
	)
	return nil
}

func (r *Router) pingPeers(targets []ping.Target) []uint8 {
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
	results := make([]uint8, len(targets))
	for i, resultCh := range channelArray {
		results[i] = mapDurationToUint8(<-resultCh)
	}
	return results
}

// mapDurationToUint8 maps a duration to an uint8 value
// the duration is clamped to 0-400ms and mapped linearly onto 1-255
// based on testing, we may need to adjust this mapping
func mapDurationToUint8(duration time.Duration) uint8 {
	durationMs := duration.Milliseconds()
	if durationMs < 0 {
		durationMs = 0
	} else if durationMs > 400 {
		durationMs = 400
	}
	return uint8((float64(durationMs)/400.0)*254.0) + 1
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

type SenderType int

const (
	originator SenderType = iota + 1
	localRelayer
	RemoteRelayer
)

type Selector struct {
	*Router
	HeightLock    sync.Mutex
	LoggedHR      map[string]uint64
	RecentHeights [50]uint64
	HeightIndex   int
}

func (s *Selector) clusterStatus(peerCluster [][]common.Address, height uint64, round int64, code uint8, from common.Address, senderType SenderType, ownClusterID int) {
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
	var zeroDisconnectClusters []string
	var totalConnected int
	sender := "originator"
	switch senderType {
	case originator:
		sender = "originator"
	case localRelayer:
		sender = "local relayer"
	case RemoteRelayer:
		sender = "remote relayer"
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

	sb.WriteString(fmt.Sprintf("\nCluster routing status:\t Height=%d, Round=%d, From=%s Message=%s SenderType=%s localCluster=%d\n", height, round, from.Hex(), msgType, sender, ownClusterID))

	for clusterID, cluster := range peerCluster {
		var lostPeers []string
		connectedCount := 0

		for _, peer := range cluster {
			_, ok := s.broadcaster.FindPeer(peer)
			if ok {
				connectedCount++
				totalConnected++
			} else {
				lostPeers = append(lostPeers, peer.Hex())
				totalDisconnected++
			}
			totalSelected++
		}

		if len(lostPeers) == 0 {
			zeroDisconnectClusters = append(zeroDisconnectClusters,
				fmt.Sprintf("C%d:%d", clusterID, len(cluster)))
		} else {
			sb.WriteString(fmt.Sprintf("Cluster #%d: selected:%d connected:%d\n", clusterID, len(cluster), connectedCount))
			sb.WriteString("  X disconnected:")
			for _, peerHex := range lostPeers {
				sb.WriteString(" ")
				sb.WriteString(peerHex)
			}
			sb.WriteByte('\n')
		}
	}

	if len(zeroDisconnectClusters) > 0 {
		sb.WriteString("Fully connected: ")
		for i, clusterInfo := range zeroDisconnectClusters {
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

func (s *Selector) SelectPeers(committee *types.Committee, msg message.Msg, from common.Address) ([]types.CommitteeMember, error) {

	// around epoch rotation, resolveClusters() can be failed since router has its own epoch synchronization context,
	// in this case, the caller should relay the message to all the members of input committee.
	clusters, err := s.resolveClusters(msg.H())
	if err != nil {
		return nil, err
	}

	seed := seed(msg)
	// if node is the original msg sender, it selects K*VerticalRelayingRedundancy relayers from every cluster vertically.
	var recipients []types.CommitteeMember
	if from == s.self {
		ownCluster := clusters.clusterContaining(from)
		// select local cluster nodes

		// select relayers from other clusters
		results := clusters.selectK(VerticalRelayingRedundancy, seed, ownCluster, s.broadcaster)
		receivers := make([]common.Address, 0)
		for _, receiver := range results {
			receivers = append(receivers, receiver...)
		}
		//numOfRelayers := len(receivers)

		// send to local cluster nodes too.
		localClusterNodes := clusters.clusterByID(ownCluster)
		if localClusterNodes != nil {
			receivers = append(receivers, localClusterNodes...)
			results[ownCluster] = localClusterNodes
		}

		for _, addr := range receivers {
			//todo: this check should move to receivers selection itself, else we run a risk of not selecting peers at all
			if member := committee.MemberByAddress(addr); member != nil && addr != s.self {
				recipients = append(recipients, *member)
			}
		}
		s.clusterStatus(results, msg.H(), msg.R(), msg.Code(), from, originator, ownCluster)

		//log.Debug(
		//	"Router: sending msg to other clusters and local cluster",
		//	"from", from,
		//	"H", msg.H(),
		//	"R", msg.R(),
		//	"C", msg.Code(),
		//	"other cluster", numOfRelayers,
		//	"local cluster", len(recipients)-numOfRelayers,
		//)
		return recipients, nil
	}

	// if node is in the same cluster of the original sender, relay the msg to local cluster nodes.
	if ownCluster := clusters.clusterContaining(s.self); ownCluster == clusters.clusterContaining(from) && ownCluster >= 0 {
		//todo: define a redundancy factor for this.
		// select and forward to 1 relayer in other clusters
		results := clusters.selectK(HorizontalRelayingRedundancy, seed, ownCluster, s.broadcaster)
		receivers := make([]common.Address, 0)
		for _, receiver := range results {
			receivers = append(receivers, receiver...)
		}

		cluster := clusters.base[ownCluster]
		// relayer in local cluster also tries to send to local members, but only to the sqrt
		k := int(math.Sqrt(float64(len(cluster))))
		if k == 0 && len(cluster) > 0 {
			k = 1
		}
		validIndices := make([]int, 0, len(cluster))
		for i, addr := range cluster {
			if addr != from && addr != s.self {
				if _, ok := s.broadcaster.FindPeer(addr); ok {
					validIndices = append(validIndices, i)
				}
			}
		}
		rand.Shuffle(len(validIndices), func(i, j int) {
			validIndices[i], validIndices[j] = validIndices[j], validIndices[i]
		})
		selectCount := k
		if selectCount > len(validIndices) {
			selectCount = len(validIndices)
		}

		localAddress := make([]common.Address, 0)
		for i := 0; i < selectCount; i++ {
			addr := cluster[validIndices[i]]
			receivers = append(receivers, addr)
			localAddress = append(localAddress, addr)
		}

		for _, addr := range receivers {
			if member := committee.MemberByAddress(addr); member != nil {
				recipients = append(recipients, *member)
			}
		}

		results[ownCluster] = localAddress
		s.clusterStatus(results, msg.H(), msg.R(), msg.Code(), from, localRelayer, ownCluster)
		//log.Debug(
		//	"Router: sending message to local cluster",
		//	"from",
		//	from,
		//	"self",
		//	s.self,
		//	"H", msg.H(),
		//	"R", msg.R(),
		//	"C", msg.Code(),
		//	"local cluster",
		//	len(recipients),
		//)
		return recipients, nil
	}

	// if we are relaying the messages from outside our own cluster, we should forward it to our own cluster vertically.
	// moreover that, we also need to forward it to the other clusters horizontally to increase the robustness of messaging.
	if ownCluster := clusters.clusterContaining(s.self); ownCluster != clusters.clusterContaining(from) && ownCluster >= 0 {

		// select local cluster nodes only
		localAddress := make([]common.Address, 0)
		for _, addr := range clusters.base[ownCluster] {
			if member := committee.MemberByAddress(addr); member != nil && addr != s.self {
				// try only connected peers
				if _, ok := s.broadcaster.FindPeer(addr); ok {
					recipients = append(recipients, *member)
					localAddress = append(localAddress, addr)
				}
			}
		}

		// to add robustness, we also relay message to other clusters horizontally.
		//results := clusters.selectK(HorizontalRelayingRedundancy, seed, ownCluster, s.broadcaster)
		//relayers := make([]common.Address, 0)
		//for _, relayer := range results {
		//	relayers = append(relayers, relayer...)
		//}
		//for _, addr := range relayers {
		//	if member := committee.MemberByAddress(addr); member != nil && addr != from {
		//		recipients = append(recipients, *member)
		//	}
		//}
		results := make([][]common.Address, len(clusters.base))
		results[ownCluster] = localAddress // add local nodes to the results
		s.clusterStatus(results, msg.H(), msg.R(), msg.Code(), from, RemoteRelayer, ownCluster)
		//log.Debug(
		//	"Router: sending message to own cluster, and horizontally relaying to other clusters",
		//	"from",
		//	from,
		//	"self",
		//	s.self,
		//	"H", msg.H(),
		//	"R", msg.R(),
		//	"C", msg.Code(),
		//	"local cluster", len(recipients)-len(relayers),
		//	"other cluster", len(relayers),
		//)
	}
	return recipients, nil
}

func seed(msg message.Msg) int64 {
	// this ensures we end up with a seed that is > 0 < math.MaxInt64, but is still reliant on
	// the message hash and the message height and round
	mh := int64(msg.H()) + msg.R() + 1
	//hash := new(big.Int).
	//	Mod(
	//		msg.Hash().Big(),
	//		new(big.Int).Div(big.NewInt(math.MaxInt64), big.NewInt(mh)),
	//	)
	//return mh * hash.Int64()
	return mh
}

func numClustersFor(length int) int {
	return int(math.Floor(math.Sqrt(float64(length))))
}
