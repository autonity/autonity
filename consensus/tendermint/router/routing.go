package router

import (
	"context"
	"crypto/ecdsa"
	"math/rand"
	"sync"
	"time"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/fixsizecache"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/router/cache"
	"github.com/autonity/autonity/consensus/tendermint/router/cluster"
	"github.com/autonity/autonity/consensus/tendermint/router/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/router/latency"
	"github.com/autonity/autonity/consensus/tendermint/router/ping"
	"github.com/autonity/autonity/consensus/tendermint/router/selector"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/metrics"
)

const (
	DefaultScaleThresholdForClustering = int64(64)
	latencyDataExpiry                  = 5 * time.Minute
	retryLatencyTimeout                = 30 * time.Second
	cacheCleanupInterval               = 10 * time.Minute
	latencyMeasurementDelayCap         = 2000
)

var (
	proposeHashesOut = metrics.GetOrRegisterCounter("router/propose/hash/egress", nil)   //nolint:goconst
	precommitHashOut = metrics.GetOrRegisterCounter("router/precommit/hash/egress", nil) //nolint:goconst
	prevoteHashOut   = metrics.GetOrRegisterCounter("router/prevote/hash/egress", nil)   //nolint:goconst
	forwardCounter   = metrics.GetOrRegisterCounter("router/forward", nil)               //nolint:goconst
)

func Setup(
	nodeKey *ecdsa.PrivateKey,
	self common.Address,
	logger log.Logger,
) *Router {
	peerCache := cache.New()
	cm := &cluster.Manager{}
	pinger, _ := ping.NewPinger(ping.ProtocolTCP, logger)
	peerSelector := selector.New(cm, peerCache)
	fetcher := latency.NewFetcher(pinger)
	return New(nodeKey, self, peerCache, fetcher, peerSelector, cm, logger)
}

type Router struct {
	self                common.Address
	nodeKey             *ecdsa.PrivateKey
	epochEventChan      chan core.EpochHeadEvent
	epochEventSub       event.Subscription
	chain               interfaces.BlockChainProvider
	clusteringThreshold int64

	committee   []common.Address
	inCommittee bool
	cancel      context.CancelFunc
	wg          sync.WaitGroup

	latestLatencies map[common.Address]uint
	latencyMu       sync.RWMutex

	nodesToRetry map[common.Address]struct{}
	retryMu      sync.RWMutex

	network        interfaces.ClustersProvider
	peerFinder     interfaces.PeerFinder
	latencyFetcher interfaces.LatencyProvider
	peerSelector   interfaces.PeerSelector
	recipientCache cache.Recipients
	hashCache      *fixsizecache.Cache[common.Hash, bool] // the cache of self messages

	logger log.Logger
}

func New(
	nodeKey *ecdsa.PrivateKey,
	self common.Address,
	recipientCache cache.Recipients,
	latencyFetcher interfaces.LatencyProvider,
	peerSelector interfaces.PeerSelector,
	networkProvider interfaces.ClustersProvider,
	logger log.Logger,
) *Router {
	router := &Router{
		nodeKey:             nodeKey,
		epochEventChan:      make(chan core.EpochHeadEvent, 2),
		latencyFetcher:      latencyFetcher,
		recipientCache:      recipientCache,
		latestLatencies:     make(map[common.Address]uint),
		nodesToRetry:        make(map[common.Address]struct{}),
		self:                self,
		peerSelector:        peerSelector,
		network:             networkProvider,
		logger:              logger,
		clusteringThreshold: DefaultScaleThresholdForClustering,
	}
	if metrics.Enabled {
		router.hashCache = fixsizecache.New[common.Hash, bool](5987, 5, fixsizecache.HashKey[common.Hash])
	}
	return router
}

func (m *Router) Pinger() ping.Pinger {
	return m.latencyFetcher.Pinger()
}

func (m *Router) SetPinger(pinger ping.Pinger) {
	m.latencyFetcher.SetPinger(pinger)
}

func (m *Router) Selector() interfaces.PeerSelector {
	return m.peerSelector
}

func (m *Router) SetSelector(selector interfaces.PeerSelector) {
	m.peerSelector = selector
}

func (m *Router) committeeAddresses(committee *types.Committee) []common.Address {
	addresses := make([]common.Address, committee.Len())
	for i, member := range committee.Members {
		addresses[i] = member.Address
	}
	return addresses
}

func (m *Router) clusteringThresholdByHeight(height uint64) int64 {
	threshold, err := m.chain.ClusteringThresholdByHeight(height)
	if err != nil {
		log.Error("Router: failed to get clustering threshold", "error", err)
		return DefaultScaleThresholdForClustering
	}
	return threshold.Int64()
}

func (m *Router) Recipients(committee *types.Committee, msg message.Msg, from common.Address) ([]common.Address, error) {
	if int64(committee.Len()) <= m.clusteringThreshold {
		return m.committeeAddresses(committee), nil
	}
	recipients, err := m.peerSelector.SelectPeers(committee, msg, from)
	if err != nil {
		m.logger.Debug("selector: no clusters, falling back to all committee members")
		return m.committeeAddresses(committee), nil
	}
	return recipients, nil
}

func (m *Router) recordDistinctHash(msg message.Msg) {
	if !metrics.Enabled {
		return
	}

	if m.hashCache.Contains(msg.Hash()) {
		return
	}
	m.hashCache.Add(msg.Hash(), true)

	switch msg.Code() {
	case message.ProposalCode:
		proposeHashesOut.Inc(1)
	case message.PrecommitCode:
		precommitHashOut.Inc(1)
	case message.PrevoteCode:
		prevoteHashOut.Inc(1)
	default:
	}
}

// TODO: rename to Send and deal with all the edge cases (full broadcast, forward, slow gossip, etc.) inside this function itself
func (m *Router) Forward(committee *types.Committee, msg message.Msg, sender common.Address, recipients []common.Address) {
	if m.peerFinder == nil {
		m.logger.Info("Router: peer finder not set")
		return
	}
	if len(recipients) == 0 {
		recipients, _ = m.Recipients(committee, msg, sender)
	}

	m.recordDistinctHash(msg)
	if sender != m.self && metrics.Enabled {
		// simple forward
		forwardCounter.Inc(1)
	}

	lostPeers := make([]common.Address, 0)
	for _, recipient := range recipients {
		if recipient == sender {
			continue
		}
		if p, ok := m.peerFinder.FindPeer(recipient); ok {
			if p.Cache().Contains(msg.Hash()) {
				continue
			}
			p.Cache().Add(msg.Hash(), true)
			go func() {
				err := p.SendRaw(message.NetworkCodes[msg.Code()], msg.Payload())
				if err != nil {
					m.logger.Error("Router: failed to send message", "recipient", recipient.Hex(), "error", err)
					return
				}
			}()
		} else {
			lostPeers = append(lostPeers, recipient)
		}
	}
	if len(lostPeers) > 0 {
		m.logger.Debug("Router: peers not found", "len", len(lostPeers), "peers", lostPeers)
	}
}

func (m *Router) Start(ctx context.Context, chain interfaces.BlockChainProvider) {
	m.logger.Info("Router: starting")

	curEpoch, err := chain.LatestEpoch()
	if err != nil {
		panic("Error fetching latest epoch: " + err.Error())
	}

	m.epochEventSub = chain.SubscribeEpochHeadEvent(m.epochEventChan)
	addresses := make([]common.Address, curEpoch.Committee.Len())
	for i, member := range curEpoch.Committee.Members {
		addresses[i] = member.Address
	}
	m.committee = addresses
	m.inCommittee = curEpoch.Committee.MemberByAddress(m.self) != nil
	nw, err := cluster.New(addresses, m.latestLatencies, m.self)
	m.chain = chain
	m.clusteringThreshold = m.clusteringThresholdByHeight(curEpoch.EpochBlock.Uint64() + 1)
	if err != nil {
		m.logger.Error("Router: failed to create network", "err", err)
	} else {
		m.updateNetwork(nw)
	}
	ctx, m.cancel = context.WithCancel(ctx)
	m.wg.Add(1)
	go m.loop(ctx)
}

func (m *Router) Stop() {
	m.cancel()
	m.epochEventSub.Unsubscribe()
	m.wg.Wait()
}

func (m *Router) SetBroadcaster(broadcaster interfaces.PeerFinder) {
	m.peerFinder = broadcaster
	m.latencyFetcher.SetBroadcaster(broadcaster)
	m.peerSelector.SetBroadcaster(broadcaster)
}

func (m *Router) refreshClustersLatencies(latMap map[common.Address]uint) {
	nw, err := cluster.New(m.committee, latMap, m.self)
	if err != nil {
		m.logger.Error("Router: failed to create network", "err", err)
		return
	}
	m.updateNetwork(nw)
}

func (m *Router) measureLatency() error {
	m.logger.Info("Router: measure latency")
	latencyMap, failedNodes, err := m.latencyFetcher.Fetch(m.committee, m.self)
	if err != nil {
		m.logger.Error("Router: failed to fetch latency", "err", err)
		return err
	}

	m.latencyMu.Lock()
	for addr, lat := range latencyMap {
		m.latestLatencies[addr] = lat
	}
	m.latencyMu.Unlock()

	m.retryMu.Lock()
	m.nodesToRetry = make(map[common.Address]struct{})
	for _, addr := range failedNodes {
		m.nodesToRetry[addr] = struct{}{}
	}
	m.retryMu.Unlock()

	m.refreshClustersLatencies(m.Latencies())
	m.logger.Debug("Router: latency measurement completed", "failed_nodes", len(failedNodes))
	return nil
}

func (m *Router) retryLatency() error {
	m.retryMu.Lock()
	if len(m.nodesToRetry) == 0 {
		m.retryMu.Unlock()
		return nil
	}

	nodes := make([]common.Address, 0, len(m.nodesToRetry))
	for addr := range m.nodesToRetry {
		nodes = append(nodes, addr)
	}
	m.retryMu.Unlock()

	m.logger.Debug("Router: retrying latency for nodes", "count", len(nodes))
	latencyMap, failedNodes, err := m.latencyFetcher.Fetch(nodes, m.self)
	if err != nil || len(latencyMap) == 0 {
		m.logger.Error("Router: failed to retry latency", "err", err)
		return err
	}

	updated := false
	m.latencyMu.Lock()
	for addr, lat := range latencyMap {
		if _, ok := m.latestLatencies[addr]; !ok || m.latestLatencies[addr] == cluster.DefaultLatency {
			m.latestLatencies[addr] = lat
			updated = true
		}
	}
	m.latencyMu.Unlock()

	m.retryMu.Lock()
	m.nodesToRetry = make(map[common.Address]struct{})
	for _, addr := range failedNodes {
		m.nodesToRetry[addr] = struct{}{}
	}
	m.retryMu.Unlock()

	if updated {
		m.refreshClustersLatencies(m.Latencies())
		m.logger.Debug("Router: clusters updated after latency retry")
	}

	m.logger.Debug("Router: latency retry completed", "failed_nodes", len(failedNodes))
	return nil
}

func (m *Router) loop(ctx context.Context) {
	defer m.wg.Done()

	retryTicker := time.NewTicker(retryLatencyTimeout)
	cleanupTicker := time.NewTicker(cacheCleanupInterval)
	// wait for few seconds before starting the initial latency measurement
	initialMeasurementTimer := time.NewTimer(5 * time.Second)
	defer func() {
		retryTicker.Stop()
		cleanupTicker.Stop()
		initialMeasurementTimer.Stop()
	}()

	wasClustering := false
	if m.inCommittee && int64(len(m.committee)) >= m.clusteringThreshold {
		wasClustering = true
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-initialMeasurementTimer.C:
			if !m.inCommittee || m.peerFinder == nil || int64(len(m.committee)) < m.clusteringThreshold {
				continue
			}
			m.logger.Debug("Router: initial latency measurement started", "threshold", m.clusteringThreshold)
			if err := m.measureLatency(); err != nil {
				m.logger.Warn("measureToReport failed", "err", err)
			}
			initialMeasurementTimer.C = nil
		case <-time.After(latencyDataExpiry +
			time.Duration(rand.Intn(latencyMeasurementDelayCap))*time.Millisecond):
			if !m.inCommittee || m.peerFinder == nil || int64(len(m.committee)) < m.clusteringThreshold {
				continue
			}
			if err := m.measureLatency(); err != nil {
				m.logger.Warn("measureToReport failed", "err", err)
			}
		case <-retryTicker.C:
			retryTicker = time.NewTicker(retryLatencyTimeout)
			if !m.inCommittee || m.peerFinder == nil || int64(len(m.committee)) < m.clusteringThreshold {
				continue
			}
			if err := m.retryLatency(); err != nil {
				m.logger.Warn("retryLatency failed", "err", err)
			}
		case <-cleanupTicker.C:
			m.recipientCache.Cleanup()
			m.logger.Debug("Router: cache cleanup completed")
		case epochEv := <-m.epochEventChan:
			m.logger.Info("Router: new epoch detected", "height", epochEv.Header.Number.String())
			epoch := epochEv.Header.Epoch
			m.inCommittee = epoch.Committee.MemberByAddress(m.self) != nil
			if !m.inCommittee || int64(len(m.committee)) < m.clusteringThreshold {
				m.logger.Info("Router: clustering not needed, skipping measurement")
				if wasClustering {
					// reset cluster
					m.updateNetwork(cluster.Clusters{})
					wasClustering = false
				}
				continue
			}
			wasClustering = true
			m.updateCommittee(epoch)
			m.clusteringThreshold = m.clusteringThresholdByHeight(epochEv.Header.Number.Uint64() + 1)
			// we should never fail here, as we already made sure that we are in committee
			nw, err := cluster.New(m.committee, m.Latencies(), m.self)
			if err != nil {
				panic("Router: failed to create clusters: " + err.Error())
			}
			m.updateNetwork(nw)
			if err := m.measureLatency(); err != nil {
				m.logger.Warn("measureLatencies failed", "err", err)
			}
		}
	}
}

func (m *Router) updateCommittee(epoch *types.Epoch) {
	addresses := make([]common.Address, epoch.Committee.Len())
	for i, member := range epoch.Committee.Members {
		addresses[i] = member.Address
	}
	m.committee = addresses
}

func (m *Router) updateNetwork(clusters cluster.Clusters) {
	m.network.UpdateClusters(clusters)
	m.recipientCache.Invalidate()
}

func (m *Router) Latencies() map[common.Address]uint {
	m.latencyMu.RLock()
	defer m.latencyMu.RUnlock()
	latCopy := make(map[common.Address]uint, len(m.latestLatencies))
	for k, v := range m.latestLatencies {
		latCopy[k] = v
	}
	return latCopy
}
