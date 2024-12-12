package latency

import (
	"crypto/ecdsa"
	"errors"
	"math"
	"math/big"
	"time"

	probing "github.com/prometheus-community/pro-bing"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/acn/protocol"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/log"
)

// ClusterRedundancyParameter is the number of members of each cluster to send a proposal to
var ClusterRedundancyParameter = 3
var ErrInvalidPeerType = errors.New("invalid peer type")

type Router struct {
	self     common.Address
	clusters Clusters

	broadcaster consensus.Broadcaster
	contracts   *autonity.ProtocolContracts
	reporter    *Reporter

	reportedEventChan chan *autonity.LatencyReported
	reportEventSub    event.Subscription

	epochHeadCh  chan *autonity.AutonityNewEpoch
	epochHeadSub event.Subscription
}

func NewRouter(
	chainId *big.Int,
	contracts *autonity.ProtocolContracts,
	broadcaster consensus.Broadcaster,
	nodeKey *ecdsa.PrivateKey,
) (*Router, error) {
	reporter, err := NewReporter(chainId, nodeKey, contracts)
	if err != nil {
		return nil, err
	}

	r := &Router{
		self:              reporter.txOpts.From,
		broadcaster:       broadcaster,
		contracts:         contracts,
		reportedEventChan: make(chan *autonity.LatencyReported),
		epochHeadCh:       make(chan *autonity.AutonityNewEpoch),
		reporter:          reporter,
	}

	r.reportEventSub, err = contracts.Latency.WatchReported(nil, r.reportedEventChan, nil)
	if err != nil {
		return nil, err
	}

	r.epochHeadSub, err = contracts.AutonityContract.WatchNewEpoch(nil, r.epochHeadCh)

	return r, nil
}

func (r *Router) Route(committee *types.Committee, msg message.Msg, from common.Address) []types.CommitteeMember {
	// if not part of the committee return
	if member := committee.MemberByAddress(from); member == nil {
		return nil
	}
	if msg.Code() != message.ProposalCode {
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

func (r *Router) Start() {
	go func() {
		for {
			select {
			case <-r.reportedEventChan:
				// todo: should probably be done async
				if err := r.refreshClusters(); err != nil {
					log.Error("failed to refresh clusters", "err", err)
				}
			case <-r.epochHeadCh:
				if err := r.report(); err != nil {
					log.Error("failed to report latency", "err", err)
				}
			}
		}
	}()
}

func (r *Router) Stop() {
	r.reportEventSub.Unsubscribe()
	r.epochHeadSub.Unsubscribe()
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

	latencyMat := make(map[common.Address][]uint8)
	for i, validator := range committee {
		latencyMat[validator] = latency[i]
	}

	clusters, err := AssignClusters(latencyMat, int(math.Floor(math.Sqrt(float64(len(committee))))))
	if err != nil {
		return err
	}

	r.clusters = clusters
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
	peerIps := make([]string, len(committee))
	for i, member := range committee {
		if member == r.self {
			continue
		}
		if peer, ok := r.broadcaster.FindPeer(member); ok {
			// todo: this is a bit hacky, we should probably have a better way to get the p2p peer ip
			p2pPeer, ok := peer.(*protocol.Peer)
			if !ok {
				return nil, ErrInvalidPeerType
			}
			peerIps[i] = p2pPeer.RemoteAddr().String()
		}
	}

	results := PingPeers(peerIps)
	latencyArray := mapStatsToUint8(results)
	for i, addr := range committee {
		// set self latency to 0
		if addr == r.self {
			latency[addr] = 0
		}
		latency[addr] = latencyArray[i]
	}
	return latency, nil
}

func mapStatsToUint8(pingResults []probing.Statistics) []uint8 {
	latency := make([]uint8, len(pingResults))
	for i, result := range pingResults {
		latency[i] = mapDurationToUint8(result.AvgRtt)
	}
	return latency
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
