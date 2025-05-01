package router

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/router/ping"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/log"
)

const (
	ScaleThresholdForClustering = 6
	DefaultLatency              = 128
	farThreshold                = 130
	nearThreshold               = 32
	diversityThreshold          = 50
	MeasurementWindow           = 2000
)

type Router struct {
	self            common.Address
	nodeKey         *ecdsa.PrivateKey
	clusterMu       sync.RWMutex
	clusters        Clusters
	broadcaster     consensus.Broadcaster
	epochEventChan  chan core.EpochHeadEvent
	epochEventSub   event.Subscription
	committee       []common.Address
	inCommittee     bool
	fetcher         *LatencyFetcher
	peerSelector    PeerSelector
	cache           *PeerSelectionCache
	cancel          context.CancelFunc
	wg              sync.WaitGroup
	latestLatencies map[common.Address]uint
	latencyMu       sync.RWMutex
	nodesToRetry    map[common.Address]struct{}
	retryMu         sync.RWMutex
}

func New(
	broadcaster consensus.Broadcaster,
	nodeKey *ecdsa.PrivateKey,
	pinger ping.Pinger,
	selector PeerSelector,
	self common.Address,
) *Router {
	cache := NewPeerSelectionCache()
	router := &Router{
		broadcaster:     broadcaster,
		nodeKey:         nodeKey,
		epochEventChan:  make(chan core.EpochHeadEvent, 2),
		fetcher:         NewLatencyFetcher(pinger, broadcaster),
		cache:           cache,
		latestLatencies: make(map[common.Address]uint),
		nodesToRetry:    make(map[common.Address]struct{}),
		self:            self,
	}
	if selector == nil {
		router.peerSelector = NewSelector(router)
	} else {
		router.peerSelector = selector
	}
	return router
}

func (m *Router) committeeAddresses(committee *types.Committee) []common.Address {
	result := make([]common.Address, committee.Len())
	for i, member := range committee.Members {
		result[i] = member.Address
	}
	return result
}

func (m *Router) Route(committee *types.Committee, msg message.Msg, from common.Address) ([]common.Address, error) {
	if committee.Len() <= ScaleThresholdForClustering {
		return m.committeeAddresses(committee), nil
	}
	return m.peerSelector.SelectPeers(committee, msg, from)
}

func (m *Router) Forward(committee *types.Committee, msg message.Msg, sender common.Address) {
	recipients, err := m.Route(committee, msg, sender)
	if err != nil {
		log.Debug("Forward: No recipients for message, broadcast", "error", err, "height", msg.H(), "message type", msg.Code())
		recipients = m.committeeAddresses(committee)
	}
	lostPeers := make([]common.Address, 0)
	for _, recipient := range recipients {
		if recipient == sender {
			continue
		}
		if p, ok := m.broadcaster.FindPeer(recipient); ok {
			if p.Cache().Contains(msg.Hash()) {
				continue
			}
			p.Cache().Add(msg.Hash(), true)
			go p.SendRaw(message.NetworkCodes[msg.Code()], msg.Payload())
		} else {
			lostPeers = append(lostPeers, recipient)
		}
	}
	if len(lostPeers) > 0 {
		log.Debug("Router: peers not found", "len", len(lostPeers), "peers", lostPeers)
	}
}

func (m *Router) Start(ctx context.Context, chain *core.BlockChain, address common.Address) {
	log.Info("Router: starting routing manager")

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
	m.updateClusters(NewClusters(result, m.latestLatencies, m.self))

	ctx, m.cancel = context.WithCancel(ctx)
	m.wg.Add(1)
	go m.loop(ctx)
}

func (m *Router) Stop() {
	m.cancel()
	m.epochEventSub.Unsubscribe()
	m.wg.Wait()
}

func (m *Router) SetBroadcaster(broadcaster consensus.Broadcaster) {
	m.broadcaster = broadcaster
	m.fetcher.broadcaster = broadcaster
}

func (m *Router) refreshClustersLatencies(latMap map[common.Address]uint) {
	m.updateClusters(NewClusters(m.committee, latMap, m.self))
}

func (m *Router) measureLatency() error {
	log.Info("Router: measure latency")
	latencyMap, failedNodes, err := m.fetcher.FetchLatency(m.committee, m.self)
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

	m.refreshClustersLatencies(m.latestLatencies)
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
	latencyMap, failedNodes, err := m.fetcher.FetchLatency(nodes, m.self)
	if err != nil || len(latencyMap) == 0 {
		log.Error("Router: failed to retry latency", "err", err)
		return err
	}

	updated := false
	m.latencyMu.Lock()
	for addr, lat := range latencyMap {
		if _, ok := m.latestLatencies[addr]; !ok || m.latestLatencies[addr] == DefaultLatency {
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

	ticker := time.NewTicker(5 * time.Minute)
	retryTicker := time.NewTicker(90 * time.Second)
	cleanupTicker := time.NewTicker(10 * time.Minute)
	defer func() {
		ticker.Stop()
		retryTicker.Stop()
		cleanupTicker.Stop()
	}()

	if len(m.committee) >= ScaleThresholdForClustering {
		if err := m.measureLatency(); err != nil {
			log.Warn("Latency measurement failed", "err", err)
		}
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if !m.inCommittee || m.broadcaster == nil || len(m.committee) < ScaleThresholdForClustering {
				continue
			}
			delay := time.Duration(rand.Intn(MeasurementWindow)) * time.Millisecond
			time.Sleep(delay)
			if err := m.measureLatency(); err != nil {
				log.Warn("measureToReport failed", "err", err)
			}
		case <-retryTicker.C:
			if !m.inCommittee || m.broadcaster == nil || len(m.committee) < ScaleThresholdForClustering {
				continue
			}
			if err := m.retryLatency(); err != nil {
				log.Warn("retryLatency failed", "err", err)
			}
		case <-cleanupTicker.C:
			m.cache.Cleanup()
			log.Debug("Router: cache cleanup completed")
		case epochEv := <-m.epochEventChan:
			log.Info("Router: new epoch detected", "height", epochEv.Header.Number.String())
			epoch := epochEv.Header.Epoch
			m.inCommittee = epoch.Committee.MemberByAddress(m.self) != nil
			if !m.inCommittee || len(m.committee) < ScaleThresholdForClustering {
				log.Info("Router: clustering not needed, skipping measurement")
				continue
			}
			m.updateCommittee(epoch)
			m.updateClusters(NewClusters(m.committee, m.latestLatencies, m.self))
			if err := m.measureLatency(); err != nil {
				log.Warn("measureToReport failed", "err", err)
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

	m.retryMu.Lock()
	m.nodesToRetry = make(map[common.Address]struct{})
	for _, addr := range result {
		if addr != m.self {
			m.nodesToRetry[addr] = struct{}{}
		}
	}
	m.retryMu.Unlock()
}

func (m *Router) updateClusters(c Clusters) {
	m.clusterMu.Lock()
	defer m.clusterMu.Unlock()
	var sb strings.Builder
	sb.WriteString("Updating cluster, new cluster view: [")
	for i, cv := range c.base {
		if i > 0 {
			sb.WriteString("; ")
		}
		sb.WriteString(fmt.Sprintf("\nC%d:[", i))
		for j, m := range cv.Members {
			if j > 0 {
				sb.WriteString(",")
			}
			addr := m.Addr.Hex()
			sb.WriteString(fmt.Sprintf("%s:%d", addr, m.Lat))
		}
		sb.WriteString("]\n")
	}
	sb.WriteString("]")
	log.Info(sb.String())
	m.clusters = c
	m.cache.Invalidate()
}

func (m *Router) Clusters() Clusters {
	m.clusterMu.RLock()
	defer m.clusterMu.RUnlock()
	return m.clusters
}

func (m *Router) Latencies() map[common.Address]uint {
	m.latencyMu.RLock()
	defer m.latencyMu.RUnlock()
	return m.latestLatencies
}
