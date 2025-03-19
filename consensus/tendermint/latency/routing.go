package latency

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math"
	"math/big"
	"math/rand"
	"slices"
	"sync"
	"time"

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

// proposeNetworkMsg is redefined here to avoid circular dependencies
var proposeNetworkMsg uint64 = 0x11

// ScaleThresholdForClustering is the minimum number of validators required to do network clustering
var ScaleThresholdForClustering = 9 // by according to the simulation and testing, there was minimal difference in performance when the number of validators was < 32.
// ClusterRedundancyParameter is the number of members of each cluster to send a proposal to
var ClusterRedundancyParameter = 5

var MeasurementWindow = 10000 // The time window in Millisecond to measure the latency of peers at the beginning of an epoch.

type PeerSelector interface {
	SelectPeers(committee *types.Committee, msg message.Msg, from common.Address) ([]types.CommitteeMember, error)
}

type Router struct {
	self    common.Address
	nodeKey *ecdsa.PrivateKey

	mu                    sync.RWMutex
	epochDefaultCluster   *Clusters
	epochOptimizedCluster *Clusters
	lastEpochCluster      *Clusters

	broadcaster consensus.Broadcaster
	contracts   *autonity.ProtocolContracts
	reporter    *Reporter

	optimizationEventChan chan *autonity.LatencyClusteringViewOptimized
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
		epochEventChan:        make(chan core.EpochHeadEvent),
		optimizationEventChan: make(chan *autonity.LatencyClusteringViewOptimized),
		pinger:                ping.NewPinger(ping.TCP),
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
	r.peerSelector = &Selector{r}
}

func (r *Router) PeerSelector() PeerSelector {
	return r.peerSelector
}

// Exported functions

// Route just select recipients from the clusters, it does not do the message sending.
func (r *Router) Route(committee *types.Committee, msg message.Msg, from common.Address) ([]types.CommitteeMember, error) {
	return r.PeerSelector().SelectPeers(committee, msg, from)
}

// Forward just forward decoded proposal from p2p msg handler, the proposal could be a future proposal within
// current epoch or from the next epoch when node is around epoch rotation, thus, the router should be able to buffer
// future proposals that the clustering haven't been done.
func (r *Router) Forward(bc *core.BlockChain, m message.Msg, sender common.Address) {
	if m.H() < r.curEpochInfo.EpochBlock.Uint64() {
		log.Info("Don't forward too old message", "height", m.H())
		return
	}

	if m.H() >= r.curEpochInfo.NextEpochBlock.Uint64() {
		log.Info("Buffer future epoch's message", "height", m.H())
		return
	}

	committee, err := bc.CommitteeByHeight(m.H())
	if err != nil {
		log.Info("Forward: Failed to query epoch", "error", err, "height", m.H())
		return
	}

	var recipients []types.CommitteeMember
	recipients, err = r.Route(committee, m, sender)
	if err != nil {
		if !errors.Is(err, consensus.ErrFutureEpochMessage) {
			log.Debug("No recipients for proposal", "error", err, "height", m.H())
			return
		}
		// forward to all the committee members, as most of them are still in the committee.
		recipients = committee.Members
	}

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
			go p.SendRaw(proposeNetworkMsg, m.Payload()) //nolint
		} else {
			//todo: shall we select other backups for live ness?
		}
	}
}

func (r *Router) Start(ctx context.Context, chain *core.BlockChain) {
	log.Info("Router: starting latency router")
	optimizationEventSub, err := chain.ProtocolContracts().Latency.WatchClusteringViewOptimized(nil, r.optimizationEventChan)
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
	optimizationHeight, err := r.contracts.GetNewViewHeight(nil)
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
	r.epochDefaultCluster = c
}

func (r *Router) setOptimizedCluster(c *Clusters) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.epochOptimizedCluster = c
}

func (r *Router) viewRotation(newDefaultCluster *Clusters) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.epochOptimizedCluster != nil {
		r.lastEpochCluster = r.epochOptimizedCluster
	} else {
		r.lastEpochCluster = r.epochDefaultCluster
	}

	r.epochDefaultCluster = newDefaultCluster
	r.epochOptimizedCluster = nil
}

func (r *Router) resolveClusters(h uint64) (*Clusters, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.epochDefaultCluster != nil && h >= r.epochDefaultCluster.nextEpochHeight {
		return nil, consensus.ErrFutureEpochMessage
	}

	// always try to pick the optimized one 1st
	if r.epochOptimizedCluster != nil && h >= r.epochOptimizedCluster.activatedHeight && h < r.epochDefaultCluster.nextEpochHeight {
		return r.epochOptimizedCluster, nil
	}
	// otherwise, try to pick the default clusters.
	if r.epochDefaultCluster != nil && h >= r.epochDefaultCluster.activatedHeight && h < r.epochDefaultCluster.nextEpochHeight {
		return r.epochDefaultCluster, nil
	}

	// edge case, around epoch rotation, some message might be from past epoch.
	if r.lastEpochCluster != nil && h >= r.lastEpochCluster.activatedHeight && h < r.lastEpochCluster.nextEpochHeight {
		return r.lastEpochCluster, nil
	}

	if r.lastEpochCluster != nil && h < r.lastEpochCluster.activatedHeight {
		return nil, errTooOldMessage
	}

	return nil, errUnknownClusters
}

// buildDefaultClusters partitions the committee into default clusters
func (r *Router) buildDefaultClusters(committee []common.Address) *Clusters {
	if len(committee) <= ScaleThresholdForClustering {
		return nil
	}
	numClusters := numClustersFor(len(committee))
	clusters := make([][]common.Address, numClusters)
	for i, addr := range committee {
		k := int(math.Min(float64(i/numClusters), float64(numClusters-1)))
		clusters[k] = append(clusters[k], addr)
	}

	defaultCluster := &Clusters{r.curEpochInfo.EpochBlock.Uint64(),
		r.curEpochInfo.NextEpochBlock.Uint64(), clusters, nil}

	log.Debug("Router: set default clusters", "clusters", func() [][]int {
		clusterInts := make([][]int, len(clusters))
		for i, cluster := range clusters {
			for _, member := range cluster {
				clusterInts[i] = append(clusterInts[i], slices.Index(committee, member))
			}
		}
		return clusterInts
	}(), "height", r.curEpochInfo.EpochBlock.Uint64())
	return defaultCluster
}

func (r *Router) startMeasurementTask() {
	// todo: clean shutdown for this go routine.
	go func() {
		rand.Seed(time.Now().UnixNano())
		delay := time.Duration(rand.Intn(MeasurementWindow)) * time.Millisecond
		time.Sleep(delay)

		if err := r.measureToReport(); err != nil {
			log.Warn("measureToReport", "err", err)
		}
	}()
}

func (r *Router) measureToReport() error {
	// todo: read committee from local cache.
	// todo: double check the data race.
	committee, err := r.contracts.Latency.GetCommittee(nil)
	if err != nil {
		return err
	}
	latencyVec, err := r.fetchLatency(committee)
	if err != nil {
		return err
	}
	return r.reporter.ReportLatency(latencyVec)
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
	for {
		select {
		case <-ctx.Done():
			return
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

			if r.curEpochInfo.Committee.MemberByAddress(r.self) == nil {
				log.Info("Router: node leaving committee, skip measurement", "height", epochEv.Header.Number.String())
				continue
			}

			// start the measurement and try to report the data.
			r.startMeasurementTask()
		}
	}
}

func (r *Router) optimizeCluster(h uint64) error {
	committee, err := r.contracts.Latency.GetCommittee(nil)
	if err != nil {
		log.Error("Router: optimizeCluster fetch committee", "err", err)
		return err
	}

	latency, err := r.contracts.Latency.Read(nil)
	if err != nil {
		log.Error("Router: optimizeCluster failed to read latency", "err", err)
		return err
	}

	// TODO: validate and fill latency matrix with default values
	latencyMat := make(map[common.Address][]uint8)
	for i, validator := range committee {
		latencyMat[validator] = latency[i]
	}

	optimizedClusters, err := AssignClusters(h, r.curEpochInfo.NextEpochBlock.Uint64(), latencyMat, int(math.Floor(math.Sqrt(float64(len(committee))))))
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
// the duration is clamped to 0-400ms and mapped linearly onto 0-255
// based on testing, we may need to adjust this mapping
func mapDurationToUint8(duration time.Duration) uint8 {
	durationMs := duration.Milliseconds()
	if durationMs < 0 {
		durationMs = 0
	} else if durationMs > 400 {
		durationMs = 400
	}
	return uint8((float64(durationMs) / 400.0) * 255.0)
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
}

func (s *Selector) SelectPeers(committee *types.Committee, msg message.Msg, from common.Address) ([]types.CommitteeMember, error) {
	if member := committee.MemberByAddress(from); member == nil {
		log.Debug("Router: from not part of committee, not routing anywhere", "address", from)
		return nil, nil
	}

	if msg.Code() != message.ProposalCode {
		return committee.Members, nil
	}

	clusters, err := s.resolveClusters(msg.H())
	if err != nil {
		return nil, err
	}

	// if we are sending the proposal, we should send it to every cluster
	var recipients []types.CommitteeMember
	if from == s.self {
		for _, addr := range clusters.selectK(ClusterRedundancyParameter, seed(msg)) {
			if member := committee.MemberByAddress(addr); member != nil {
				recipients = append(recipients, *member)
			}
		}

		// we should also send directly to every outlier
		for _, addr := range clusters.direct {
			if member := committee.MemberByAddress(addr); member != nil {
				recipients = append(recipients, *member)
			}
		}
	} else {
		log.Debug(
			"Router: not the originator of the proposal, not sending to other clusters",
			"from",
			from,
			"self",
			s.self,
		)
	}

	// if we are receiving the proposal from outside our own cluster, we should send it to our own cluster
	if ownCluster := clusters.clusterContaining(s.self); ownCluster != clusters.clusterContaining(from) && ownCluster >= 0 {
		for _, addr := range clusters.base[ownCluster] {
			if member := committee.MemberByAddress(addr); member != nil {
				recipients = append(recipients, *member)
			}
		}
	} else {
		log.Debug(
			"Router: not receiving from outside cluster, not sending to own cluster",
			"from",
			from,
			"self",
			s.self,
		)
	}

	// there could be some duplication if we are sending to ClusterRedundancyParameter members of each
	// cluster, that may include our own cluster, so we deduplicate
	return deduplicate(recipients), nil
}

func deduplicate(recipients []types.CommitteeMember) []types.CommitteeMember {
	seen := make(map[common.Address]struct{})
	var result []types.CommitteeMember
	for _, rec := range recipients {
		if _, ok := seen[rec.Address]; !ok {
			seen[rec.Address] = struct{}{}
			result = append(result, rec)
		}
	}
	return result
}

func seed(msg message.Msg) int64 {
	// this ensures we end up with a seed that is > 0 < math.MaxInt64, but is still reliant on
	// the message hash and the message height and round
	mh := int64(msg.H())*msg.R() + 1
	hash := new(big.Int).
		Mod(
			msg.Hash().Big(),
			new(big.Int).Div(big.NewInt(math.MaxInt64), big.NewInt(mh)),
		)
	return mh * hash.Int64()
}

func numClustersFor(length int) int {
	return int(math.Floor(math.Sqrt(float64(length))))
}
