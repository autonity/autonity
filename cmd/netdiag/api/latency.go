package api

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/autonity/autonity/cmd/netdiag/core"
)

func (p *P2POp) TriggerLatencyBroadcast(arg *ArgStrategy, _ *ArgEmpty) error {
	// first set our own latency matrix
	latency := core.PingPeers(p.Engine)
	p.Engine.State.LatencyMatrix[p.Engine.Id] = core.FilterAveRtt(latency)
	if err := core.BroadcastLatency(p.Engine, uint64(arg.Strategy), latency); err != nil {
		return fmt.Errorf("error in broadcast latency: %s", err.Error())
	}

	// then broadcast to all peers
	errs := make([]<-chan error, len(p.Engine.Peers))
	for id, peer := range p.Engine.Peers {
		ch := make(chan error, 1)
		errs[id] = ch
		if id == p.Engine.Id {
			ch <- nil
			continue
		}
		if peer == nil || !peer.Connected {
			ch <- errTargetNotConnected
		} else {
			go func(peer *core.Peer, ch chan error) {
				err := peer.SendTriggerRequest(uint64(arg.Strategy))
				ch <- err
			}(peer, ch)
		}
	}
	for _, ch := range errs {
		err := <-ch
		if err != nil {
			return fmt.Errorf("error in send trigger request: %s", err.Error())
		}
	}
	return nil
}

func (p *P2POp) BroadcastLatencyArray(strat *ArgStrategy, reply *ResultBroadcastLatencyArray) error {
	result := ResultBroadcastLatencyArray{}
	errs := make([]string, len(p.Engine.Peers))
	acks := make([]bool, len(p.Engine.Peers))
	var hasError atomic.Bool
	var wg sync.WaitGroup

	for i, peer := range p.Engine.Peers {
		if i == p.Engine.Id {
			errs[i] = ""
			continue
		}
		if peer == nil || !peer.Connected {
			errs[i] = "peer not connected or nil"
			hasError.Store(true)
			continue
		}

		wg.Add(1)
		go func(id int, peer *core.Peer) {
			_, _, err := peer.SendLatencyArray(uint64(strat.Strategy), p.Engine.State.LatencyMatrix[p.Engine.Id])
			if err != nil {
				hasError.Store(true)
				errs[id] = err.Error()
			} else {
				acks[id] = true
				errs[id] = ""
			}
			wg.Done()
		}(i, peer)
	}

	wg.Wait()
	result.HasError = hasError.Load()
	result.Errors = errs
	result.AckReceived = acks
	*reply = result
	return nil
}

type ResultLatencyMatrix [][]time.Duration

func (r *ResultLatencyMatrix) String() string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "Latency matrix global view:\n")
	valid := true
	for id, latencyArray := range *r {
		fmt.Fprintf(&builder, "peer %d\t", id)
		for peerID, l := range latencyArray {
			if peerID == id {
				if l != 0 {
					valid = false
				}
				fmt.Fprintf(&builder, "0\t")
				continue
			}
			if l == 0 {
				fmt.Fprintf(&builder, "inf\t")
			} else {
				fmt.Fprintf(&builder, "%d\t", l)
			}
		}
		fmt.Fprintf(&builder, "\n")
	}
	fmt.Fprintf(&builder, "Is valid? %v\n", valid)
	return builder.String()
}

func (p *P2POp) LatencyMatrix(_ *ArgEmpty, reply *ResultLatencyMatrix) error {
	result := make([][]time.Duration, len(p.Engine.Enodes))
	for i, latencyArray := range p.Engine.State.LatencyMatrix {
		result[i] = make([]time.Duration, len(p.Engine.Enodes))
		copy(result[i], latencyArray)
	}
	*reply = (ResultLatencyMatrix)(result)
	return nil
}
