package latency

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math"
	"net"
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

// ScaleThresholdForClustering is the minimum number of validators required to do network clustering
var ScaleThresholdForClustering = 1 // by according to the simulation and testing, there was minimal difference in performance when the number of validators was < 32.
// ClusterRedundancyParameter is the number of members of each cluster to send a proposal to
var ClusterRedundancyParameter = 3
var ErrInvalidPeerType = errors.New("invalid peer type")

type peerLatency interface {
	consensus.Peer
	RemoteAddr() net.Addr
	Node() *enode.Node
}

type Router struct {
	self        common.Address
	clusterLock sync.RWMutex
	clusters    Clusters
	nodeKey     *ecdsa.PrivateKey

	broadcaster consensus.Broadcaster
	contracts   *autonity.ProtocolContracts
	reporter    *Reporter

	reportedEventChan chan *autonity.LatencyReported
	reportEventSub    event.Subscription

	epochEventChan chan core.EpochHeadEvent
	epochEventSub  event.Subscription

	chainEventChan chan core.ChainEvent
	chainEventSub  event.Subscription

	curEpochInfo *types.EpochInfo
	// measurementWindow uint64
	lastReportedHeight  uint64
	lastRefreshedHeight uint64
	measured            bool

	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewRouter(
	broadcaster consensus.Broadcaster,
	nodeKey *ecdsa.PrivateKey,
) *Router {
	r := &Router{
		broadcaster:       broadcaster,
		nodeKey:           nodeKey,
		reportedEventChan: make(chan *autonity.LatencyReported),
		epochEventChan:    make(chan core.EpochHeadEvent),
		chainEventChan:    make(chan core.ChainEvent),
	}
	return r
}

// Route just select recipients from the clusters, it does not do the message sending.
func (r *Router) Route(committee *types.Committee, msg message.Msg, from common.Address) []types.CommitteeMember {
	// if not part of the committee return
	if member := committee.MemberByAddress(from); member == nil {
		return nil
	}

	r.clusterLock.RLock()
	defer r.clusterLock.RUnlock()

	// currently only proposals are routed through clustering
	// if the clusters are not yet formed, or there is no clusters at all, we should default to the full committee
	if msg.Code() != message.ProposalCode || r.clusters == nil {
		return committee.Members
	}

	// if we are sending the proposal, we should send it to every cluster
	var recipients []types.CommitteeMember
	if from == r.self {
		for _, addr := range r.clusters.selectK(ClusterRedundancyParameter) {
			if member := committee.MemberByAddress(addr); member != nil {
				recipients = append(recipients, *member)
			}
		}
	}
	// if we are receiving the proposal from outside our own cluster, we should send it to our own cluster
	if ownCluster := r.clusters.clusterContaining(r.self); ownCluster != r.clusters.clusterContaining(from) && ownCluster >= 0 {
		for _, addr := range r.clusters[ownCluster] {
			if member := committee.MemberByAddress(addr); member != nil {
				recipients = append(recipients, *member)
			}
		}
	}

	return recipients
}

func (r *Router) ClusteringActive() bool {
	r.clusterLock.RLock()
	defer r.clusterLock.RUnlock()
	return r.clusters != nil
}

func (r *Router) Start(ctx context.Context, chain *core.BlockChain) {
	log.Info("Router: starting latency router")
	reportEventSub, err := chain.ProtocolContracts().Latency.WatchReported(nil, r.reportedEventChan, nil)
	if err != nil {
		log.Error("Error starting reported event subscription", "err", err)
		return
	}

	curEpoch, err := chain.LatestEpoch()
	if err != nil {
		log.Error("Error fetching latest epoch", "err", err)
		return
	}
	r.curEpochInfo = curEpoch
	r.measured = false // measure on startup

	r.epochEventSub = chain.SubscribeEpochHeadEvent(r.epochEventChan)
	r.chainEventSub = chain.SubscribeChainEvent(r.chainEventChan)
	r.reportEventSub = reportEventSub
	r.contracts = chain.ProtocolContracts()
	r.reporter, err = NewReporter(chain.Config().ChainID, r.nodeKey, r.contracts)
	if err != nil {
		log.Error("failed to create reporter", "err", err)
		return
	}
	r.self = r.reporter.txOpts.From

	ctx, r.cancel = context.WithCancel(ctx)
	r.wg.Add(1)
	go r.loop(ctx)
}

func (r *Router) loop(ctx context.Context) {
	defer r.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case epochEv := <-r.epochEventChan:
			r.curEpochInfo = &types.EpochInfo{
				Epoch:      *epochEv.Header.Epoch.Copy(),
				EpochBlock: epochEv.Header.Number,
			}
			r.measured = false

			// on new epoch, we have a small scale of network, clustering does not benefit anymore.
			if r.curEpochInfo.Committee.Len() <= ScaleThresholdForClustering {
				r.resetClusters()
			}

		case ev := <-r.chainEventChan:
			if r.curEpochInfo.Committee.Len() <= ScaleThresholdForClustering {
				log.Info("not going to measure latency within a small network")
				continue
			}

			height := ev.Block.NumberU64()

			if !r.measured {
				log.Info(
					"Router: new epoch reporting latency",
					"height",
					height,
					"reporter",
					r.self,
				)
				if err := r.report(); err != nil {
					log.Error("Router: failed to report latency", "err", err)
				} else {
					r.measured = true
				}
			} else {
				log.Info(
					"Router: already reported, skipping",
					"height",
					height,
				)
			}

			if r.lastReportedHeight > r.lastRefreshedHeight {
				if err := r.refreshClusters(); err != nil {
					log.Error("Router: failed to refresh clusters", "err", err)
				} else {
					r.lastRefreshedHeight = height
				}
			}

		case ev := <-r.reportedEventChan:
			if r.curEpochInfo.Committee.Len() <= ScaleThresholdForClustering {
				log.Info("Router: not going to cluster a small scale network")
				continue
			}
			log.Info("Router: latency report detected, scheduling network clustering", "reporter", ev.Reporter)
			r.lastReportedHeight = max(r.lastReportedHeight, ev.Raw.BlockNumber)
		}
	}
}

func (r *Router) Stop() {
	r.cancel()
	r.reportEventSub.Unsubscribe()
	r.chainEventSub.Unsubscribe()
	r.epochEventSub.Unsubscribe()
	r.wg.Wait()
}

func (r *Router) SetBroadcaster(broadcaster consensus.Broadcaster) {
	r.broadcaster = broadcaster
}

func (r *Router) refreshClusters() error {
	committee, err := r.contracts.Latency.GetCommittee(nil)
	if err != nil {
		return err
	}

	latency, err := r.contracts.Latency.Read(nil)
	if err != nil {
		return err
	}

	// TODO: validate and fill latency matrix with default values
	latencyMat := make(map[common.Address][]uint8)
	for i, validator := range committee {
		latencyVec := latency[i]
		for j, peer := range committee {
			// 0 is reserved for self, ^uint8(0) is reserved for non-connected peers
			if validator == peer {
				latencyVec[j] = 0
			} else if latencyVec[j] == 0 {
				latencyVec[j] = ^uint8(0)
			}
		}
		latencyMat[validator] = latencyVec
	}

	clusters, err := AssignClusters(latencyMat, int(math.Floor(math.Sqrt(float64(len(committee))))))
	if err != nil {
		return err
	}

	r.clusterLock.Lock()
	r.clusters = clusters
	r.clusterLock.Unlock()
	return nil
}

// reset clusters, it is used to merge the cluster when we have a small scale of network.
func (r *Router) resetClusters() {
	r.clusterLock.Lock()
	defer r.clusterLock.Unlock()
	r.clusters = nil
}

func (r *Router) report() error {
	latency, err := r.fetchLatency()
	if err != nil {
		return err
	}

	return r.reporter.ReportLatency(latency)
}

func (r *Router) fetchLatency() (map[common.Address]uint8, error) {
	if r.broadcaster == nil {
		return nil, errors.New("broadcaster not set, can't fetch latency")
	}

	committee, err := r.contracts.Latency.GetCommittee(nil)
	if err != nil {
		return nil, err
	}

	committeeEnodes := r.broadcaster.CommitteeEnodes()

	latency := make(map[common.Address]uint8)
	pingTargets := make([]ping.Target, len(committee))

	for i, member := range committee {
		if member == r.self {
			pingTargets[i] = ping.Target{}
			continue
		}
		if memberNode, ok := findByAddress(committeeEnodes, member); ok {
			// todo: this is a less hacky, but we should probably have a better way to get the peer ip
			ip := memberNode.IP()
			port := memberNode.TCP()
			pingTargets[i] = ping.Target{IP: ip.String(), Port: port}
			log.Info("Router: fetching latency", "targetIP", ip, "targetPort", port)
		} else {
			log.Error("Router: peer not found in broadcaster", "peer", member)
			pingTargets[i] = ping.Target{}
		}
	}

	latencyArray := PingPeers(pingTargets)
	for i, addr := range committee {
		// set self latency to 0
		if addr == r.self {
			latency[addr] = 0
		}
		latency[addr] = latencyArray[i]
	}
	return latency, nil
}

func PingPeers(targets []ping.Target) []uint8 {
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
		// icmp pinger to compare results
		ping.NewPinger(ping.TCP).Ping(t, resultCh)
		channelArray[i] = resultCh
	}
	results := make([]uint8, len(targets))
	for i, resultCh := range channelArray {
		results[i] = mapDurationToUint8(<-resultCh)
	}
	return results
}

// mapDurationToUint8 maps a duration to a uint8 value
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
