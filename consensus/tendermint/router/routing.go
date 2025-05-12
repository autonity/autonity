package router

import (
	"context"
	"crypto/ecdsa"
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
	clusters    *ClusterRotation
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
	reportedThisEpoch bool
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
		broadcaster:       broadcaster,
		nodeKey:           nodeKey,
		epochEventChan:    make(chan core.EpochHeadEvent, 2),
		fetcher:           NewLatencyFetcher(pinger, broadcaster),
		cache:             cache,
		latestLatencies:   make(map[common.Address]uint),
		nodesToRetry:      make(map[common.Address]struct{}),
		self:              self,
		reportedThisEpoch: false,
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
	m.contracts = chain.ProtocolContracts()
	m.epochEventSub = chain.SubscribeEpochHeadEvent(m.epochEventChan)

	m.committee = m.committeeAddresses(curEpoch.Committee)
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
	m.initializeClusters(curEpoch)

	ctx, m.cancel = context.WithCancel(ctx)
	m.wg.Add(2)
	go m.watchReported(ctx)
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
	m.clusters.UpdateLatencies(latMap, m.self)
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

func (m *Router) report() error {
	log.Info("Router: checking latency report status")
	if m.reportedThisEpoch {
		log.Info("Router: already reported this epoch, skipping")
		return nil
	}
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
	m.reportedThisEpoch = true
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
		m.wg.Add(1)
		go func() {
			defer m.wg.Done()
			delay := time.Duration(rand.Intn(MeasurementWindow)) * time.Millisecond
			time.Sleep(delay)
			if err := m.report(); err != nil {
				log.Warn("Router: initial latency report failed", "err", err)
			}
		}()
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
				log.Info("Router: not in committee clustering not needed, skipping measurement")
				continue
			}
			m.updateCommittee(epoch)
			if epochEv.Header.Number.Uint64() == m.clusters.latestEpochBlock {
				log.Info("Router: clusters for this epoch already established, skipping clustering")
				continue
			}

			log.Info("Router: new epoch, constructing default clusters")
			clusters, err := NewClusters(m.committee, m.latestLatencies, nil, m.self)
			if err != nil {
				// this should never happen
				panic("Router: failed to create default clusters")
			} else {
				log.Info("Router: created default clusters")
			}

			m.clusters.EpochStart(epochEv.Header.Number.Uint64(), clusters)
			m.logUpdateClusters(epochEv.Header.Number.Uint64(), clusters)
			// new epoch, reset the reported status
			m.reportedThisEpoch = false
			// reset lat mat
			m.latencyMat = nil
			if err := m.measureLatency(); err != nil {
				log.Warn("measureToReport failed", "err", err)
			}
			m.wg.Add(1)
			go func() {
				delay := time.Duration(rand.Intn(MeasurementWindow)) * time.Millisecond
				time.Sleep(delay)
				defer m.wg.Done()
				if err := m.report(); err != nil {
					log.Warn("Router: initial latency report failed", "err", err)
				}
			}()
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

func (m *Router) readLatencyMatrix(previousEpoch bool) ([][]uint8, []common.Address, error) {
	var committee []common.Address
	var err error
	if previousEpoch {
		committee, err = m.contracts.Latency.GetLastCommittee(nil)
	} else {
		committee, err = m.contracts.Latency.GetCommittee(nil)
	}

	if err != nil {
		log.Error("Router: optimizeCluster fetch committee", "err", err)
		return nil, nil, err
	}

	// as reading the entire matrix could be reverted due to too much gas consumption, we have to read row by row.
	latencyMat := make([][]uint8, len(committee))
	index := new(big.Int).SetUint64(0)
	for i, _ := range committee {
		var latency []uint8
		var err error
		if previousEpoch {
			latency, err = m.contracts.Latency.ReadLastEpochReport(nil, index.SetInt64(int64(i)))
		} else {
			latency, err = m.contracts.Latency.ReadReport(nil, index.SetInt64(int64(i)))
		}
		if err != nil {
			log.Error("Router: optimizeCluster failed to read latency", "err", err)
			return nil, nil, err
		}
		latencyMat[i] = latency
	}
	return latencyMat, committee, nil
}

func (m *Router) logUpdateClusters(height uint64, c Clusters) {
	var sb strings.Builder
	log.Info("Router: updating clusters", "height", height)
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
	sb.WriteString("]\n")
	log.Info(sb.String())
}

func (m *Router) initializeClusters(epoch *types.EpochInfo) {
	log.Info("Router: initializing clusters", "epoch", epoch.EpochBlock.Uint64())
	if epoch.EpochBlock.Cmp(common.Big0) == 0 {
		// initialize with default clusters
		log.Info("Router: initializing with default clusters")
		clusters, err := NewClusters(m.committeeAddresses(epoch.Committee), m.latestLatencies, nil, m.self)
		if err != nil {
			// should never error on default clusters
			panic(err)
		}
		m.clusters = &ClusterRotation{
			previousEpochClusters: Clusters{},
			transitionalClusters:  clusters,
			latestEpochClusters:   Clusters{},
		}
		return
	}
	// get the previous clusters
	lockInBlock, err := m.contracts.Latency.MatrixLockInBlock(nil)
	if err != nil {
		log.Error("Router: failed to get lock in block", "err", err)
		return
	}
	prevLatMat, prevCommittee, err := m.readLatencyMatrix(true)
	if err != nil {
		log.Error("Router: init - failed to read latency matrix", "err", err)
		return
	}
	currentCommittee := m.committeeAddresses(epoch.Committee)
	transitionalClusters, err := NewClusters(
		currentCommittee,
		m.latestLatencies,
		getForCommittee(latencyReports{prevLatMat, prevCommittee}, currentCommittee),
		m.self,
	)
	if err != nil {
		log.Error("Router: init - failed to create transitional transitionalClusters", "err", err)
		return
	}
	prevClusters, err := NewClusters(
		prevCommittee,
		m.latestLatencies,
		prevLatMat,
		m.self,
	)
	if err != nil {
		log.Error("Router: init - failed to create previous transitionalClusters", "err", err)
		return
	}
	prevEpochLockIn, err := m.contracts.Latency.LastMatrixLockInBlock(nil)
	if err != nil {
		log.Error("Router: init -failed to get last matrix lock in block", "err", err)
		return
	}

	if lockInBlock.Cmp(common.Big0) == 0 {
		log.Info("Router: init - lock in block is 0, using transitional clusters")
		// we are in a new epoch but have not yet locked in
		// the transitionalClusters, so we need to use the transitional transitionalClusters
		m.clusters = &ClusterRotation{
			previousEpochBlock:      epoch.PreviousEpochBlock.Uint64(),
			latestEpochBlock:        epoch.EpochBlock.Uint64(),
			lastMatrixLockInBlock:   prevEpochLockIn.Uint64(),
			latestMatrixLockInBlock: lockInBlock.Uint64(),

			previousEpochClusters: prevClusters,
			transitionalClusters:  transitionalClusters,
			latestEpochClusters:   Clusters{},
		}
	} else {
		log.Info("Router: init - lock in block is not 0, using current clusters and previous clusters")
		// we are in a new epoch, and we have locked in
		currentLatMat, _, err := m.readLatencyMatrix(false)
		if err != nil {
			log.Error("Router: failed to read current latency matrix", "err", err)
			return
		}
		currentClusters, err := NewClusters(
			currentCommittee,
			m.latestLatencies,
			currentLatMat,
			m.self,
		)
		if err != nil {
			log.Error("Router: failed to create current clusters", "err", err)
			return
		}
		m.clusters = &ClusterRotation{
			previousEpochBlock:      epoch.PreviousEpochBlock.Uint64(),
			latestEpochBlock:        epoch.EpochBlock.Uint64(),
			lastMatrixLockInBlock:   prevEpochLockIn.Uint64(),
			latestMatrixLockInBlock: lockInBlock.Uint64(),
			previousEpochClusters:   prevClusters,
			transitionalClusters:    transitionalClusters,
			latestEpochClusters:     currentClusters,
		}
	}
}

func (m *Router) watchReported(ctx context.Context) {
	defer m.wg.Done()

	reported := make(chan *autonity.LatencyReported, 2)
	reportedSub, err := m.contracts.Latency.WatchReported(nil, reported, nil)
	if err != nil {
		log.Error("Router: failed to subscribe to reported event", "err", err)
		return
	}
	optimization := make(chan *autonity.LatencyKMOptimization, 2)
	kmOptimizationSub, err := m.contracts.Latency.WatchKMOptimization(nil, optimization)
	if err != nil {
		log.Error("Router: failed to subscribe to km optimization event", "err", err)
		return
	}
	defer reportedSub.Unsubscribe()
	defer kmOptimizationSub.Unsubscribe()
	for {
		select {
		case <-ctx.Done():
			return
		case ev := <-reported:
			log.Info(
				"Router: reported event received",
				"totalReported",
				ev.TotalReported,
				"nextEpoch",
				m.epoch.NextEpochBlock.Uint64(),
				"previousEpoch",
				m.epoch.PreviousEpochBlock.Uint64(),
				"block",
				ev.Raw.BlockNumber,
			)
		case optimizationEv := <-optimization:
			log.Info(
				"Router: km optimization event received",
				"effectiveHeight",
				optimizationEv.Height.Uint64(),
				"nextEpoch",
				m.epoch.NextEpochBlock.Uint64(),
				"previousEpoch",
				m.epoch.PreviousEpochBlock.Uint64(),
				"block",
				optimizationEv.Raw.BlockNumber,
			)
			if err := m.OptimizeClusters(optimizationEv.Height.Uint64()); err != nil {
				log.Error("Could not optimize clusters", "err", err)
			}
		}
	}
}

func (m *Router) OptimizeClusters(lockInBlock uint64) error {
	log.Info("Router: ready to optimize the clustering", "effectiveHeight", lockInBlock)
	if m.clusters.latestMatrixLockInBlock == lockInBlock {
		log.Info("Router: clusters already optimized for this lock in block, skipping")
		return nil
	}
	latMat, committee, err := m.readLatencyMatrix(false)
	if err != nil {
		log.Error("Router: failed to read latency matrix", "err", err)
		return err
	} else {
		log.Info("Router: read latency matrix successfully", "len(latMat)", len(latMat))
	}
	m.latencyMat = latMat
	clusters, err := NewClusters(committee, m.latestLatencies, latMat, m.self)
	if err != nil {
		log.Error("Router: failed to create new clusters", "err", err)
		return err
	}
	m.clusters.LockIn(lockInBlock, clusters)
	m.logUpdateClusters(lockInBlock, clusters)
	return nil
}

func (m *Router) Clusters(height uint64) Clusters {
	if clusters := m.clusters.GetClusters(height); len(clusters.base) > 0 {
		return clusters
	}
	log.Error(
		"Router: no clusters found",
		"height",
		height,
		"previousEpochBlock",
		m.epoch.PreviousEpochBlock.Uint64(),
		"nextEpochBlock",
		m.epoch.NextEpochBlock.Uint64(),
	)
	return Clusters{}
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
