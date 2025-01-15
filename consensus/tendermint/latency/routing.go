package latency

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math"
	"math/big"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/latency/ping"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/log"
)

// ClusterRedundancyParameter is the number of members of each cluster to send a proposal to
var ClusterRedundancyParameter = 3
var ErrInvalidPeerType = errors.New("invalid peer type")

type peerLatency interface {
	consensus.Peer
	RemoteAddr() net.Addr
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

	curEpochInfo      *types.EpochInfo
	measurementWindow uint64
	measured          bool
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

func (r *Router) Route(committee *types.Committee, msg message.Msg, from common.Address) []types.CommitteeMember {
	// if not part of the committee return
	if member := committee.MemberByAddress(from); member == nil {
		return nil
	}
	// currently only proposals are routed through clustering
	// if the clusters are not yet formed, we should default to the full committee
	if msg.Code() != message.ProposalCode || r.clusters == nil {
		return committee.Members
	}
	// if we are sending the proposal, we should send it to every cluster
	r.clusterLock.RLock()
	defer r.clusterLock.RUnlock()
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

func (r *Router) Start(ctx context.Context, chain *core.BlockChain) error {
	reportEventSub, err := chain.ProtocolContracts().Latency.WatchReported(nil, r.reportedEventChan, nil)
	if err != nil {
		return err
	}

	curEpoch, err := chain.LatestEpoch()
	if err != nil {
		return err
	}
	r.curEpochInfo = curEpoch
	epochPeriod := new(big.Int).Sub(curEpoch.NextEpochBlock, curEpoch.EpochBlock)
	r.measurementWindow = epochPeriod.Uint64() / uint64(curEpoch.Committee.Len())

	r.epochEventSub = chain.SubscribeEpochHeadEvent(r.epochEventChan)
	r.chainEventSub = chain.SubscribeChainEvent(r.chainEventChan)
	r.reportEventSub = reportEventSub
	r.contracts = chain.ProtocolContracts()
	r.reporter, err = NewReporter(chain.Config().ChainID, r.nodeKey, r.contracts)
	if err != nil {
		return err
	}
	r.self = r.reporter.txOpts.From

	for {
		select {
		case <-ctx.Done():
			return nil
		case epochEv := <-r.epochEventChan:
			r.curEpochInfo = &types.EpochInfo{
				Epoch:      *epochEv.Header.Epoch.Copy(),
				EpochBlock: epochEv.Header.Number,
			}
			epochPeriod = new(big.Int).Sub(epochEv.Header.Epoch.NextEpochBlock, epochEv.Header.Number)
			r.measurementWindow = epochPeriod.Uint64() / uint64(r.curEpochInfo.Committee.Len())
			r.measured = false

		case ev := <-r.chainEventChan:
			if r.measurementWindow == 0 {
				log.Error("invalid report window, too short epoch period?")
				continue
			}

			// todo: shall we skip clustering in a small scale network?
			if r.curEpochInfo.Committee.Len() == 1 {
				log.Debug("not going to measure latency within a small network")
				continue
			}

			height := ev.Block.NumberU64()
			committee := r.curEpochInfo.Committee
			reporterIndex := (height / r.measurementWindow) % uint64(committee.Len())
			// every validator is assigned with an independent measurement and reporting window.
			if !r.measured && committee.Members[reporterIndex].Address == r.self {
				log.Debug("Router: in reporter slot, reporting latency", "height", height, "epoch period",
					epochPeriod.Uint64(), "reporter idx", reporterIndex, "reporter", r.self)
				if err := r.report(); err != nil {
					log.Error("failed to report latency", "err", err)
				} else {
					r.measured = true
				}
			}

		case ev := <-r.reportedEventChan:
			// For every measurementWindow, there will be a unique validator assigned to measure and send the latencies.
			log.Debug("Router: latency report detected, refreshing network clustering", "reporter", ev.Reporter)
			if err := r.refreshClusters(); err != nil {
				log.Error("failed to refresh clusters", "err", err)
			}
		}
	}
}

func (r *Router) Stop() {
	r.reportEventSub.Unsubscribe()
	r.chainEventSub.Unsubscribe()
	r.epochEventSub.Unsubscribe()
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
		latencyMat[validator] = latency[i]
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

func (r *Router) report() error {
	latency, err := r.fetchLatency()
	if err != nil {
		return err
	}

	return r.reporter.ReportLatency(latency)
}

func (r *Router) fetchLatency() (map[common.Address]uint8, error) {
	committee, err := r.contracts.Latency.GetCommittee(nil)
	if err != nil {
		return nil, err
	}

	latency := make(map[common.Address]uint8)
	pingTargets := make([]ping.Target, len(committee))
	for i, member := range committee {
		if member == r.self {
			continue
		}
		if peer, ok := r.broadcaster.FindPeer(member); ok {
			// todo: this is a bit hacky, we should probably have a better way to get the p2p peer ip
			p2pPeer, ok := peer.(peerLatency)
			if !ok {
				return nil, ErrInvalidPeerType
			}

			ip, port, err := net.SplitHostPort(p2pPeer.RemoteAddr().String())
			if err != nil {
				//TODO
				continue
			}

			p, _ := strconv.Atoi(port)
			pingTargets[i] = ping.Target{IP: ip, Port: p}
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
			resultCh <- time.Duration(time.Second) * 5
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
