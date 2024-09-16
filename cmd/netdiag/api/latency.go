package api

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	probing "github.com/prometheus-community/pro-bing"

	"github.com/autonity/autonity/cmd/netdiag/core"
	"github.com/autonity/autonity/cmd/netdiag/strats"
	"github.com/autonity/autonity/log"
)

func (p *P2POp) TriggerLatencyBroadcast(arg *ArgStrategy, _ *ArgEmpty) error {
	// first set our own latency matrix
	var latency []probing.Statistics
	t, _ := p.Engine.Strategies[arg.Strategy].LatencyType()
	if t == strats.LatencyTypeRelative {
		log.Debug("Pinging all peers", "strategy", arg.Strategy, "latencyType", "relative")
		latency = core.PingPeers(p.Engine)
	} else {
		log.Debug("Pinging NTP servers", "strategy", arg.Strategy, "latencyType", "fixed")
		latency = core.PingFixedNTP()
	}
	log.Debug("Got latency results", "aveRTT", core.FilterAveRtt(latency, t))

	p.Engine.State.LatencyMatrix[p.Engine.Id] = core.FilterAveRtt(latency, t)
	if err := core.BroadcastLatency(p.Engine, uint64(arg.Strategy), latency); err != nil {
		return fmt.Errorf("error in broadcast latency: %s", err.Error())
	}
	return core.TriggerLatencyBroadcast(p.Engine, uint64(arg.Strategy))
}

type ResultBroadcastLatencyArray struct {
	HasError    bool
	AckReceived []bool
	Errors      []string
}

func (r *ResultBroadcastLatencyArray) String() string {
	var builder strings.Builder

	if r.HasError {
		errCount := 0
		target := make([]int, 0, len(r.Errors))
		for i, err := range r.Errors {
			if err != "" {
				errCount++
				target = append(target, i)
			}
		}
		fmt.Fprintf(&builder, "Got %d errors from the following targets\n", errCount)
		for _, id := range target {
			fmt.Fprintf(&builder, "%d, ", id)
		}
		fmt.Fprintf(&builder, "\n")
	}

	ackReceived := make([]int, 0, len(r.AckReceived))
	ackNotReceived := make([]int, 0, len(r.AckReceived))
	for i, received := range r.AckReceived {
		if received {
			ackReceived = append(ackReceived, i)
		} else {
			ackNotReceived = append(ackNotReceived, i)
		}
	}

	fmt.Fprintf(&builder, "\nReceived responses from the following %d targets:\n", len(ackReceived))
	for _, id := range ackReceived {
		if r.Errors[id] == "" {
			fmt.Fprintf(&builder, "%d, ", id)
		}
	}

	fmt.Fprintf(&builder, "\nFollowing %d targets did not respond:\n", len(ackNotReceived)-1)
	for _, id := range ackNotReceived {
		if r.Errors[id] != "" {
			fmt.Fprintf(&builder, "%d with error: %s\n", id, r.Errors[id])
		}
	}
	fmt.Fprintf(&builder, "\n")
	return builder.String()
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
				fmt.Fprintf(&builder, "%d\t", l)
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
