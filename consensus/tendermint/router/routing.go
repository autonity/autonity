package router

import (
	"context"
	"crypto/ecdsa"
	"math/rand"
	"sync"
	"time"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/router/cache"
	"github.com/autonity/autonity/consensus/tendermint/router/constants"
	"github.com/autonity/autonity/consensus/tendermint/router/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/router/latency"
	"github.com/autonity/autonity/consensus/tendermint/router/network"
	"github.com/autonity/autonity/consensus/tendermint/router/ping"
	"github.com/autonity/autonity/consensus/tendermint/router/selector"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/log"
)

func Setup(
	nodeKey *ecdsa.PrivateKey,
	self common.Address,
	logger log.Logger,
) *Router {
	peerCache := cache.New()
	nw := &network.Network{}
	pinger, _ := ping.NewPinger(ping.ProtocolTCP, logger)
	peerSelector := selector.New(nw, peerCache)
	fetcher := latency.NewFetcher(pinger)
	return New(nodeKey, self, peerCache, fetcher, peerSelector, nw)
}

type Router struct {
	self           common.Address
	nodeKey        *ecdsa.PrivateKey
	epochEventChan chan core.EpochHeadEvent
	epochEventSub  event.Subscription
	committee      []common.Address
	inCommittee    bool
	cancel         context.CancelFunc
	wg             sync.WaitGroup

	latestLatencies map[common.Address]uint
	latencyMu       sync.RWMutex

	nodesToRetry map[common.Address]struct{}
	retryMu      sync.RWMutex

	network             interfaces.NetworkProvider
	peerFinder          interfaces.PeerFinder
	latencyFetcher      interfaces.LatencyProvider
	peerSelector        interfaces.PeerSelector
	recipientCache      cache.Recipients
	clusteringThreshold int
}

func New(
	nodeKey *ecdsa.PrivateKey,
	self common.Address,
	recipientCache cache.Recipients,
	latencyFetcher interfaces.LatencyProvider,
	peerSelector interfaces.PeerSelector,
	networkProvider interfaces.NetworkProvider,
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
		clusteringThreshold: constants.ScaleThresholdForClustering,
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
	result := make([]common.Address, committee.Len())
	for i, member := range committee.Members {
		result[i] = member.Address
	}
	return result
}

func (m *Router) Recipients(committee *types.Committee, msg message.Msg, from common.Address) ([]common.Address, error) {
	if committee.Len() <= m.clusteringThreshold {
		return m.committeeAddresses(committee), nil
	}
	recipients, err := m.peerSelector.SelectPeers(committee, msg, from)
	if err != nil {
		log.Debug("selector: no clusters, falling back to all committee members")
		return m.committeeAddresses(committee), nil
	}
	return recipients, nil
}

func (m *Router) Forward(committee *types.Committee, msg message.Msg, sender common.Address) {
	if m.peerFinder == nil {
		log.Info("Router: peer finder not set")
		return
	}

	recipients, _ := m.Recipients(committee, msg, sender)
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
					log.Error("Router: failed to send message", "recipient", recipient.Hex(), "error", err)
					return
				}
			}()
		} else {
			lostPeers = append(lostPeers, recipient)
		}
	}
	if len(lostPeers) > 0 {
		log.Debug("Router: peers not found", "len", len(lostPeers), "peers", lostPeers)
	}
}

func (m *Router) Start(ctx context.Context, chain interfaces.BlockChainProvider) {
	log.Info("Router: starting")

	curEpoch, err := chain.LatestEpoch()
	if err != nil {
		log.Error("Error fetching latest epoch", "err", err)
		return
	}

	m.epochEventSub = chain.SubscribeEpochHeadEvent(m.epochEventChan)
	result := make([]common.Address, curEpoch.Committee.Len())
	for i, member := range curEpoch.Committee.Members {
		result[i] = member.Address
	}
	m.committee = result
	m.inCommittee = curEpoch.Committee.MemberByAddress(m.self) != nil
	nw, err := network.New(result, m.latestLatencies, m.self)
	if err != nil {
		log.Error("Router: failed to create network", "err", err)
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
	nw, err := network.New(m.committee, latMap, m.self)
	if err != nil {
		log.Error("Router: failed to create network", "err", err)
		return
	}
	m.updateNetwork(nw)
}

func (m *Router) measureLatency() error {
	log.Info("Router: measure latency")
	latencyMap, failedNodes, err := m.latencyFetcher.Fetch(m.committee, m.self)
	if err != nil {
		log.Error("Router: failed to fetch latency", "err", err)
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
	log.Debug("Router: latency measurement completed", "failed_nodes", len(failedNodes))
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

	log.Debug("Router: retrying latency for nodes", "count", len(nodes))
	latencyMap, failedNodes, err := m.latencyFetcher.Fetch(nodes, m.self)
	if err != nil || len(latencyMap) == 0 {
		log.Error("Router: failed to retry latency", "err", err)
		return err
	}

	updated := false
	m.latencyMu.Lock()
	for addr, lat := range latencyMap {
		if _, ok := m.latestLatencies[addr]; !ok || m.latestLatencies[addr] == constants.DefaultLatency {
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
		m.latencyMu.Lock()
		latMap := make(map[common.Address]uint, len(m.latestLatencies))
		for addr, lat := range m.latestLatencies {
			latMap[addr] = lat
		}
		m.latencyMu.Unlock()
		m.refreshClustersLatencies(latMap)
		log.Debug("Router: clusters updated after latency retry")
	}

	log.Debug("Router: latency retry completed", "failed_nodes", len(failedNodes))
	return nil
}

func (m *Router) loop(ctx context.Context) {
	defer m.wg.Done()

	retryTicker := time.NewTicker(constants.RetryLatencyTimeout)
	cleanupTicker := time.NewTicker(constants.CacheCleanupInterval)
	// wait for few seconds before starting the initial latency measurement
	initialMeasurementTimer := time.NewTimer(5 * time.Second)
	defer func() {
		retryTicker.Stop()
		cleanupTicker.Stop()
		initialMeasurementTimer.Stop()
	}()

	wasClustering := false
	if m.inCommittee && len(m.committee) >= m.clusteringThreshold {
		wasClustering = true
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-initialMeasurementTimer.C:
			if !m.inCommittee || m.peerFinder == nil || len(m.committee) < m.clusteringThreshold {
				continue
			}
			log.Debug("Router: initial latency measurement started", "threshold", m.clusteringThreshold)
			if err := m.measureLatency(); err != nil {
				log.Warn("measureToReport failed", "err", err)
			}
			initialMeasurementTimer.C = nil
		case <-time.After(constants.LatencyDataExpiry +
			time.Duration(rand.Intn(constants.LatencyMeasurementDelayCap))*time.Millisecond):
			if !m.inCommittee || m.peerFinder == nil || len(m.committee) < m.clusteringThreshold {
				continue
			}
			if err := m.measureLatency(); err != nil {
				log.Warn("measureToReport failed", "err", err)
			}
		case <-retryTicker.C:
			retryTicker = time.NewTicker(constants.RetryLatencyTimeout)
			if !m.inCommittee || m.peerFinder == nil || len(m.committee) < m.clusteringThreshold {
				continue
			}
			if err := m.retryLatency(); err != nil {
				log.Warn("retryLatency failed", "err", err)
			}
		case <-cleanupTicker.C:
			m.recipientCache.Cleanup()
			log.Debug("Router: cache cleanup completed")
		case epochEv := <-m.epochEventChan:
			log.Info("Router: new epoch detected", "height", epochEv.Header.Number.String())
			epoch := epochEv.Header.Epoch
			m.inCommittee = epoch.Committee.MemberByAddress(m.self) != nil
			if !m.inCommittee || len(m.committee) < m.clusteringThreshold {
				log.Info("Router: clustering not needed, skipping measurement")
				if wasClustering {
					// reset cluster
					m.updateNetwork(network.Clusters{})
					wasClustering = false
				}
				continue
			}
			wasClustering = true
			m.updateCommittee(epoch)
			// we should never fail here, as we already made sure that we are in committee
			nw, _ := network.New(m.committee, m.Latencies(), m.self)
			m.updateNetwork(nw)
			if err := m.measureLatency(); err != nil {
				log.Warn("measureLatencies failed", "err", err)
			}
		}
	}
}

func (m *Router) updateCommittee(epoch *types.Epoch) {
	result := make([]common.Address, epoch.Committee.Len())
	for i, member := range epoch.Committee.Members {
		result[i] = member.Address
	}
	m.committee = result
}

func (m *Router) updateNetwork(clusters network.Clusters) {
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
