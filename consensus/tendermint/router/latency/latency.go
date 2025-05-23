package latency

import (
	"errors"
	"time"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/router/constants"
	"github.com/autonity/autonity/consensus/tendermint/router/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/router/ping"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/p2p/enode"
)

type Fetcher struct {
	pinger     ping.Pinger
	peerFinder interfaces.PeerFinder
}

func NewFetcher(pinger ping.Pinger, peerFinder interfaces.PeerFinder) *Fetcher {
	return &Fetcher{
		pinger:     pinger,
		peerFinder: peerFinder,
	}
}

// Fetch measures latency for all validators, returning the latency map and a list of nodes that failed measurement.
func (f *Fetcher) Fetch(validators []common.Address, self common.Address) (map[common.Address]uint, []common.Address, error) {
	if f.peerFinder == nil {
		return nil, nil, errors.New("broadcaster not set, can't fetch latency")
	}

	committeeEnodes := f.peerFinder.CommitteeEnodes()
	latency := make(map[common.Address]uint)
	pingTargets := make([]ping.Target, len(validators))
	failedNodes := make([]common.Address, 0)

	for i, member := range validators {
		if member == self {
			pingTargets[i] = ping.Target{}
			latency[member] = 0
			continue
		}
		if _, ok := f.peerFinder.FindPeer(member); !ok {
			log.Debug("Node not connected, skipping ping", "peer", member.Hex())
			pingTargets[i] = ping.Target{}
			failedNodes = append(failedNodes, member)
			latency[member] = constants.DefaultLatency
			continue
		}
		if memberNode, ok := enodeByAddress(committeeEnodes, member); ok {
			ip := memberNode.IP()
			port := memberNode.TCP()
			pingTargets[i] = ping.Target{IP: ip.String(), Port: port}
			log.Debug("Fetching latency", "targetIP", ip, "targetPort", port)
		} else {
			log.Error("Peer not found in broadcaster enodes", "peer", member.Hex())
			pingTargets[i] = ping.Target{}
			failedNodes = append(failedNodes, member)
			latency[member] = constants.DefaultLatency
		}
	}

	latencyArray := f.pingPeers(pingTargets)
	for i, addr := range validators {
		if addr == self || pingTargets[i].IP == "" {
			continue
		}
		if latencyArray[i].Nanoseconds() == 0 {
			failedNodes = append(failedNodes, addr)
			continue
		}
		latency[addr] = uint(latencyArray[i].Milliseconds())
	}

	return latency, failedNodes, nil
}

func (f *Fetcher) SetBroadcaster(broadcaster consensus.Broadcaster) {
	f.peerFinder = broadcaster
}

func (f *Fetcher) pingPeers(targets []ping.Target) []time.Duration {
	channelArray := make([]chan time.Duration, len(targets))
	for i, t := range targets {
		resultCh := make(chan time.Duration, 1)
		if t.IP == "" {
			resultCh <- time.Second * 5 // Default for skipped pings
			channelArray[i] = resultCh
			continue
		}
		f.pinger.Ping(t, resultCh)
		channelArray[i] = resultCh
	}
	results := make([]time.Duration, len(targets))
	for i, resultCh := range channelArray {
		results[i] = <-resultCh
	}
	return results
}

func enodeByAddress(committeeEnodes []*enode.Node, addr common.Address) (*enode.Node, bool) {
	for _, memberNode := range committeeEnodes {
		pubKey := memberNode.Pubkey()
		if crypto.PubkeyToAddress(*pubKey) == addr {
			return memberNode, true
		}
	}
	return nil, false
}
