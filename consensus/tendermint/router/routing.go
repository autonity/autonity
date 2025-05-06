package router

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/autonity/autonity/autonity"
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
	self        common.Address
	nodeKey     *ecdsa.PrivateKey
	clusterMu   sync.RWMutex
	clusters    *ClusterMap
	broadcaster consensus.Broadcaster

	committee    []common.Address
	inCommittee  bool
	fetcher      *LatencyFetcher
	peerSelector PeerSelector
	cache        *PeerSelectionCache
	cancel       context.CancelFunc
	wg           sync.WaitGroup

	latestLatencies map[common.Address]uint
	latencyMu       sync.RWMutex
	nodesToRetry    map[common.Address]struct{}
	retryMu         sync.RWMutex
	latencyMap      map[uint64]map[common.Address]uint
	latencyMat      [][]uint8
	epoch           *types.Epoch

	optimizationEventChan chan *autonity.LatencyKMOptimization
	optimizationEventSub  event.Subscription

	epochEventChan chan core.EpochHeadEvent
	epochEventSub  event.Subscription

	lastMeasuredEpoch *big.Int
	contracts         *autonity.ProtocolContracts
	reporter          *Reporter
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
		broadcaster:           broadcaster,
		nodeKey:               nodeKey,
		epochEventChan:        make(chan core.EpochHeadEvent, 2),
		optimizationEventChan: make(chan *autonity.LatencyKMOptimization, 2),
		fetcher:               NewLatencyFetcher(pinger, broadcaster),
		cache:                 cache,
		latestLatencies:       make(map[common.Address]uint),
		nodesToRetry:          make(map[common.Address]struct{}),
		self:                  self,
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
	optimizationEventSub, err := chain.ProtocolContracts().WatchKMOptimization(nil, m.optimizationEventChan)
	if err != nil {
		log.Error("Router: error subscribing to km optimization eventsq", err)
		return
	}
	m.optimizationEventSub = optimizationEventSub
	m.epochEventSub = chain.SubscribeEpochHeadEvent(m.epochEventChan)
	m.contracts = chain.ProtocolContracts()

	m.epochEventSub = chain.SubscribeEpochHeadEvent(m.epochEventChan)
	result := make([]common.Address, curEpoch.Committee.Len())
	for i, member := range curEpoch.Committee.Members {
		result[i] = member.Address
	}
	m.committee = result
	m.epoch = &curEpoch.Epoch
	m.reporter, err = NewReporter(
		chain.Config().ChainID,
		m.nodeKey,
		m.contracts,
	)
	if err != nil {
		log.Error("Router: failed to create reporter", "err", err)
		return
	}

	m.clusters = NewClusterMap(result, m.latestLatencies, m.self)
	m.updateClusters(m.epoch.PreviousEpochBlock.Uint64(), m.clusters.LatestCluster())

	ctx, m.cancel = context.WithCancel(ctx)
	m.wg.Add(1)
	go m.loop(ctx)
}

func (m *Router) Stop() {
	m.cancel()
	m.epochEventSub.Unsubscribe()
	m.optimizationEventSub.Unsubscribe()
	m.wg.Wait()
}

func (m *Router) SetBroadcaster(broadcaster consensus.Broadcaster) {
	m.broadcaster = broadcaster
	m.fetcher.broadcaster = broadcaster
}

func (m *Router) refreshClustersLatencies(latMap map[common.Address]uint) {
	m.updateClusters(
		m.clusters.LatestHeight(),
		UpdateClusterLatencies(m.clusters.LatestCluster(), latMap, m.self),
	)
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

	log.Info("Router: checking latency report status")
	member := m.epoch.Committee.MemberByAddress(m.self)
	if member == nil {
		log.Error("Router: self not in committee")
		return nil
	}
	if reported, err := m.contracts.Latency.ClientReported(nil, new(big.Int).SetUint64(member.Index)); err != nil {
		log.Error("Router: failed to check if client reported", "err", err)
		return err
	} else if !reported {
		log.Info("Router: client not reported yet, reporting now")
		if err = m.reporter.ReportLatency(toUint8(m.latestLatencies)); err != nil {
			log.Error("Router: failed to report latency", "err", err)
			return err
		} else {
			log.Info("Router: reported latency successfully")
		}
	} else {
		log.Info("Router: client already reported")
	}
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
			prevCommittee := m.committee

			// clean up the previous epoch clusters
			m.clusters.PruneTo(m.epoch.PreviousEpochBlock.Uint64())

			m.updateCommittee(epoch)
			var clusters Clusters
			var err error
			if m.latencyMat != nil {
				clusters, err = NewClusters(
					m.committee,
					m.latestLatencies,
					getForCommittee(latencyReports{m.latencyMat, prevCommittee}, m.committee),
					m.self,
				)
			} else {
				log.Error("Router: Latency matrix nil!")
				err = errors.New("router: Latency matrix nil")
			}

			if err != nil {
				log.Error("Router: failed to create new clusters", "err", err)
				if defaultClusters, err := NewClusters(m.committee, m.latestLatencies, nil, m.self); err != nil {
					// this should never happen
					panic("Router: failed to create default clusters")
				} else {
					log.Info("Router: fallback to default clusters")
					clusters = defaultClusters
				}
			}
			m.updateClusters(epoch.PreviousEpochBlock.Uint64(), clusters)

			if err := m.measureLatency(); err != nil {
				log.Warn("measureToReport failed", "err", err)
			}

		case optimizationEv := <-m.optimizationEventChan:
			log.Info("Router: ready to optimize the clustering", "height", optimizationEv.Height)
			if optimizationEv.Height.Cmp(m.epoch.NextEpochBlock) >= 0 {
				log.Info("Router: skip to optimize the clustering", "activation height cross epoch", optimizationEv.Height.Uint64())
				continue
			} else {
				log.Info("Router: optimizing clustering", "activation height", optimizationEv.Height.Uint64())
			}
			latMat, committee, err := m.readLatencyMatrix()
			if err != nil {
				log.Error("Router: failed to read latency matrix", "err", err)
				continue
			} else {
				log.Info("Router: read latency matrix successfully", "len(latMat)", len(latMat))
			}
			m.latencyMat = latMat
			clusters, err := NewClusters(committee, m.latestLatencies, latMat, m.self)
			if err != nil {
				log.Error("Router: failed to create new clusters", "err", err)
				continue
			}
			m.updateClusters(optimizationEv.Height.Uint64(), clusters)
		}
	}
}

func (m *Router) updateCommittee(epoch *types.Epoch) {
	m.epoch = epoch
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

func (m *Router) readLatencyMatrix() ([][]uint8, []common.Address, error) {
	committee, err := m.contracts.Latency.GetCommittee(nil)
	if err != nil {
		log.Error("Router: optimizeCluster fetch committee", "err", err)
		return nil, nil, err
	}

	// as reading the entire matrix could be reverted due to too much gas consumption, we have to read row by row.
	latencyMat := make([][]uint8, len(committee))
	index := new(big.Int).SetUint64(0)
	for i, _ := range committee {
		_, latency, err := m.contracts.Latency.ReadReport(nil, index.SetInt64(int64(i)))
		if err != nil {
			log.Error("Router: optimizeCluster failed to read latency", "err", err)
			return nil, nil, err
		}
		latencyMat[i] = latency
	}
	return latencyMat, committee, nil
}

func (m *Router) updateClusters(height uint64, c Clusters) {
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
	m.clusters.AddCluster(height, c)
	m.cache.Invalidate()
}

func (m *Router) Clusters(height uint64) Clusters {
	m.clusterMu.RLock()
	defer m.clusterMu.RUnlock()
	return m.clusters.GetCluster(height)
}

func (m *Router) Latencies() map[common.Address]uint {
	m.latencyMu.RLock()
	defer m.latencyMu.RUnlock()
	return m.latestLatencies
}

func toUint8(latencies map[common.Address]uint) map[common.Address]uint8 {
	result := make(map[common.Address]uint8, len(latencies))
	for addr, lat := range latencies {
		result[addr] = mapDurationToUint8(lat)
	}
	return result
}

func mapDurationToUint8(durationMs uint) uint8 {
	if durationMs == DefaultLatency {
		return uint8(DefaultLatency)
	} else if durationMs > 400 {
		durationMs = 400
	}
	return uint8((float64(durationMs)/400.0)*254.0) + 1
}
