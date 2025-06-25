package latency

import (
	"context"
	"errors"
	"sync"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/router/cluster"
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

func NewFetcher(pinger ping.Pinger) *Fetcher {
	return &Fetcher{
		pinger: pinger,
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
			latency[member] = cluster.DefaultLatency
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
			latency[member] = cluster.DefaultLatency
		}
	}

	latencyArray := f.pingPeers(context.Background(), pingTargets)
	for i, addr := range validators {
		if addr == self || pingTargets[i].IP == "" {
			continue
		}
		if latencyArray[i].Err != nil {
			failedNodes = append(failedNodes, addr)
			latency[addr] = cluster.DefaultLatency
			continue
		}
		latency[addr] = uint(latencyArray[i].Latency.Milliseconds()) // #nosec
	}

	return latency, failedNodes, nil
}

func (f *Fetcher) SetBroadcaster(broadcaster interfaces.PeerFinder) {
	f.peerFinder = broadcaster
}

func (f *Fetcher) pingPeers(ctx context.Context, targets []ping.Target) []ping.Result {
	results := make([]ping.Result, len(targets))
	var wg sync.WaitGroup
	for i, t := range targets {
		if t.IP == "" {
			results[i] = ping.Result{Err: errors.New("skipped: empty target")}
			continue
		}
		wg.Add(1)
		go func(index int, target ping.Target) {
			defer wg.Done()
			results[index] = f.pinger.Ping(ctx, target)
		}(i, t)
	}
	wg.Wait()

	return results
}

func (f *Fetcher) Pinger() ping.Pinger {
	return f.pinger
}

func (f *Fetcher) SetPinger(pinger ping.Pinger) {
	f.pinger = pinger
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
