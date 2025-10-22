package latency

import (
	"context"
	"errors"
	"math/big"
	"sort"
	"sync"

	"github.com/autonity/autonity/common"
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

	validators, committeeEnodes := sortByAddress(validators, f.peerFinder.CommitteeEnodes())

	latency := make(map[common.Address]uint)
	pingTargets := make([]ping.Target, len(validators))
	failedNodes := make([]common.Address, 0)

	enodeAt := 0
	for i, member := range validators {
		if member == self {
			pingTargets[i] = ping.Target{}
			continue
		}
		if _, ok := f.peerFinder.FindPeer(member); !ok {
			log.Debug("Node not connected, skipping ping", "peer", member.Hex())
			pingTargets[i] = ping.Target{}
			failedNodes = append(failedNodes, member)
			continue
		}
		// both `validators` and `committeeEnodes` are sorted by their address in ascending order
		// the total complexity of the loop becomes: `O(len(validators) + len(committeeEnodes))`
		if memberNode, ok, biggerEnodeIndex := findEnode(enodeAt, committeeEnodes, member); ok {
			ip := memberNode.IP()
			port := memberNode.TCP()
			pingTargets[i] = ping.Target{IP: ip.String(), Port: port}
			log.Debug("Fetching latency", "targetIP", ip, "targetPort", port)
			enodeAt = biggerEnodeIndex
		} else {
			log.Error("Peer not found in broadcaster enodes", "peer", member.Hex())
			pingTargets[i] = ping.Target{}
			failedNodes = append(failedNodes, member)
			enodeAt = biggerEnodeIndex
		}
	}

	latencyArray := f.pingPeers(context.Background(), pingTargets)
	for i, addr := range validators {
		if addr == self || pingTargets[i].IP == "" {
			continue
		}
		if latencyArray[i].Err != nil {
			failedNodes = append(failedNodes, addr)
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

func sortByAddress(validators []common.Address, committeeEnodes []*enode.Node) ([]common.Address, []*enode.Node) {

	sort.Slice(committeeEnodes, func(i, j int) bool {
		address1 := enodeAddress(committeeEnodes[i])
		address2 := enodeAddress(committeeEnodes[j])
		return new(big.Int).SetBytes(address1[:]).Cmp(new(big.Int).SetBytes(address2[:])) < 0
	})

	sort.Slice(validators, func(i, j int) bool {
		return new(big.Int).SetBytes(validators[i][:]).Cmp(new(big.Int).SetBytes(validators[j][:])) < 0
	})

	return validators, committeeEnodes
}

func findEnode(enodeAt int, enodes []*enode.Node, member common.Address) (*enode.Node, bool, int) {
	// this function assumes that the slice `endoes[enodeAt:]` is sorted by address in ascending order
	// and `member` is not present in the slice `endoes[:enodeAt+1]`
	for enodeAt < len(enodes) {
		enodeAddress := enodeAddress(enodes[enodeAt])
		diff := new(big.Int).SetBytes(enodeAddress[:]).Cmp(new(big.Int).SetBytes(member[:]))
		if diff == 0 {
			// all the remaining enode addresses are bigger than `member`
			return enodes[enodeAt], true, enodeAt + 1
		} else if diff == 1 {
			// all the remaining enode addresses are bigger than `member`
			break
		}
		enodeAt++
	}
	return nil, false, enodeAt
}

func enodeAddress(enode *enode.Node) common.Address {
	return crypto.PubkeyToAddress(*enode.Pubkey())
}
