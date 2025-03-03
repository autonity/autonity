package latency

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math"
	"math/big"
	"sync"
	"time"

	"golang.org/x/exp/slices"

	"github.com/autonity/autonity/accounts/abi/bind"
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
var ScaleThresholdForClustering = 10 // by according to the simulation and testing, there was minimal difference in performance when the number of validators was < 32.
// ClusterRedundancyParameter is the number of members of each cluster to send a proposal to
var ClusterRedundancyParameter = 3

type PeerSelector interface {
	SelectPeers(committee *types.Committee, msg message.Msg, from common.Address) []types.CommitteeMember
}

type Router struct {
	self     common.Address
	nodeKey  *ecdsa.PrivateKey
	cache    *latencyCache
	clusters *clusterCache

	broadcaster consensus.Broadcaster
	contracts   *autonity.ProtocolContracts
	reporter    *Reporter

	epochEventChan chan core.EpochHeadEvent
	epochEventSub  event.Subscription

	chainEventChan chan core.ChainEvent
	chainEventSub  event.Subscription

	curEpochInfo *types.EpochInfo
	measured     bool

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
		broadcaster:    broadcaster,
		nodeKey:        nodeKey,
		epochEventChan: make(chan core.EpochHeadEvent),
		chainEventChan: make(chan core.ChainEvent),
		pinger:         ping.NewPinger(ping.TCP),
		cache:          newLatencyCache(),
		clusters:       newClusterCache(),
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
func (r *Router) Route(committee *types.Committee, msg message.Msg, from common.Address) []types.CommitteeMember {
	return r.PeerSelector().SelectPeers(committee, msg, from)
}

func (r *Router) ClusteringActive(height uint64) bool {
	_, ok := r.clusters.clustersAt(height)
	return ok
}

func (r *Router) Start(ctx context.Context, chain *core.BlockChain) {
	log.Info("Router: starting latency router")
	curEpoch, err := chain.LatestEpoch()
	if err != nil {
		log.Error("Error fetching latest epoch", "err", err)
		return
	}
	r.curEpochInfo = curEpoch
	r.measured = false // measure on startup

	r.epochEventSub = chain.SubscribeEpochHeadEvent(r.epochEventChan)
	r.chainEventSub = chain.SubscribeChainEvent(r.chainEventChan)
	r.contracts = chain.ProtocolContracts()
	r.reporter, err = NewReporter(chain.Config().ChainID, r.nodeKey, r.contracts)
	if err != nil {
		log.Error("failed to create reporter", "err", err)
		return
	}
	r.self = r.reporter.txOpts.From

	r.setDefaultClusters(func() []common.Address {
		result := make([]common.Address, r.curEpochInfo.Committee.Len())
		for i, member := range r.curEpochInfo.Committee.Members {
			result[i] = member.Address
		}
		return result
	}())

	ctx, r.cancel = context.WithCancel(ctx)
	r.wg.Add(1)
	go r.loop(ctx)
}

func (r *Router) Stop() {
	r.cancel()
	r.chainEventSub.Unsubscribe()
	r.epochEventSub.Unsubscribe()
	r.wg.Wait()
}

func (r *Router) SetBroadcaster(broadcaster consensus.Broadcaster) {
	r.broadcaster = broadcaster
}

// Internal package functions

// setDefaultClusters partitions the committee into default clusters
func (r *Router) setDefaultClusters(committee []common.Address) {
	if len(committee) <= ScaleThresholdForClustering {
		return
	}
	numClusters := numClustersFor(committee)
	clusters := make([][]common.Address, numClusters)
	for i, addr := range committee {
		k := int(math.Min(float64(i/numClusters), float64(numClusters-1)))
		clusters[k] = append(clusters[k], addr)
	}
	r.clusters.insertClustering(r.curEpochInfo.EpochBlock.Uint64(), clusters)
}

// reset clusters, it is used to merge the cluster when we have a small scale of network.
func (r *Router) resetClusters() {
	r.clusters = newClusterCache()
}

func (r *Router) report() error {
	if r.contracts == nil {
		return errors.New("contracts not set, can't report latency")
	}
	committee, err := r.contracts.Latency.GetCommittee(nil)
	if err != nil {
		return err
	}
	missing := r.cache.missingMeasurements(committee)
	missingLatencies, err := r.fetchLatency(missing)
	if err != nil {
		return err
	}
	r.cache.insertMeasurements(missing, missingLatencies)

	// this should not have the full measurements
	latencyVec, err := r.cache.latencyView(committee)
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
		case epochEv := <-r.epochEventChan:
			log.Info("Router: new epoch detected", "height", epochEv.Header.Number.String())
			r.cache.markNewEpoch()
			// new epoch, prune height to last epoch block
			r.clusters.pruneTo(r.curEpochInfo.EpochBlock.Uint64())
			r.curEpochInfo = &types.EpochInfo{
				Epoch:      *epochEv.Header.Epoch.Copy(),
				EpochBlock: epochEv.Header.Number,
			}
			r.measured = false
			// on new epoch, we have a small scale of network, clustering does not benefit anymore.
			if epochEv.Header.Epoch.Committee.Len() <= ScaleThresholdForClustering {
				log.Warn("Router: new epoch detected, committee too small resetting clusters")
				r.resetClusters()
			}
		case ev := <-r.chainEventChan:
			if r.curEpochInfo.Committee.Len() <= ScaleThresholdForClustering {
				log.Info("Router: not going to measure latency within a small network")
				continue
			}

			height := ev.Block.Number()
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
			}
			r.processBlock(ev.Block)
		}
	}
}

func (r *Router) processBlock(block *types.Block) {
	height := block.Number()
	reports, err := parseReports(block)
	if err != nil {
		log.Error("Router: failed to parse reports", "err", err)
		return
	}

	if len(reports) == 0 {
		log.Debug("Router: no latency reports in block", "height", height.String())
		return
	}

	log.Debug("Router: processing latency reports", "reports", len(reports))
	committee, err := r.contracts.Latency.GetCommittee(
		&bind.CallOpts{BlockNumber: height},
	)
	if err != nil {
		log.Error("Router: failed to get committee", "err", err)
		return
	}

	for _, report := range reports {
		r.cache.insertMatrixLine(report.reporter, committee, report.latencies)
	}

	// if the committee is too small, we should not cluster
	if r.curEpochInfo.Committee.Len() <= ScaleThresholdForClustering {
		log.Info("Router: not going to cluster a small scale network")
		return
	}

	// if we have enough reports, we should refresh the clusters
	if r.shouldCluster() {
		latencyMat := r.cache.readMatrix(committee)
		clusters, err := AssignClusters(latencyMat, numClustersFor(committee))
		if err != nil {
			log.Error("Router: failed to assign clusters", "err", err, "height", height)
			return
		}
		log.Debug("Router: assigned clusters", "clusters", func() [][]int {
			clusterInts := make([][]int, len(clusters))
			for i, cluster := range clusters {
				clusterInts[i] = make([]int, len(cluster))
				for j, member := range cluster {
					clusterInts[i][j] = slices.Index(committee, member)
				}
			}
			return clusterInts
		}())

		r.clusters.insertClustering(height.Uint64(), clusters)
	}
}

func (r *Router) shouldCluster() bool {
	if r.curEpochInfo.Committee.Len() <= ScaleThresholdForClustering {
		return false
	}

	// if we are in the first epoch, wait for 2/3 of the committee to report
	if r.curEpochInfo.EpochBlock.Cmp(common.Big0) == 0 {
		return r.cache.reportsInEpoch() > uint64(2*len(r.curEpochInfo.Committee.Members)/3)
	}
	return true
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

type Selector struct {
	*Router
}

func (s *Selector) SelectPeers(committee *types.Committee, msg message.Msg, from common.Address) []types.CommitteeMember {

	// if not part of the committee return
	if member := committee.MemberByAddress(from); member == nil {
		log.Debug("Router: from not part of committee, not routing anywhere", "address", from)
		return nil
	}

	// currently only proposals are routed through clustering
	// if the clusters are not yet formed, or there is no clusters at all, we should default to the full committee
	if msg.Code() != message.ProposalCode || s.clusters == nil {
		return committee.Members
	}

	clusters, ok := s.clusters.clustersAt(msg.H())
	if !ok {
		// we don't have a valid clustering for this height
		log.Debug("Router: no clusters at height, routing to everyone", "height", msg.H())
		return committee.Members
	}

	// if we are sending the proposal, we should send it to every cluster
	var recipients []types.CommitteeMember
	if from == s.self {
		for _, addr := range Clusters(clusters).selectK(ClusterRedundancyParameter, seed(msg)) {
			if member := committee.MemberByAddress(addr); member != nil {
				recipients = append(recipients, *member)
			}
		}
	}
	// if we are receiving the proposal from outside our own cluster, we should send it to our own cluster
	if ownCluster := Clusters(clusters).clusterContaining(s.self); ownCluster != Clusters(clusters).clusterContaining(from) && ownCluster >= 0 {
		for _, addr := range clusters[ownCluster] {
			if member := committee.MemberByAddress(addr); member != nil {
				recipients = append(recipients, *member)
			}
		}
	}

	return recipients
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

func numClustersFor(committee []common.Address) int {
	return int(math.Floor(math.Sqrt(float64(len(committee)))))
}
