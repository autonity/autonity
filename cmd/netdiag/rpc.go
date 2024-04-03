package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	probing "github.com/prometheus-community/pro-bing"

	"github.com/autonity/autonity/log"
)

var (
	errInvalidRpcArg      = errors.New("invalid RPC argument")
	errTargetNotConnected = errors.New("target peer not connected")
)

type ArgTarget struct {
	Target int
}
type ArgEmpty struct {
}

// P2pRpc represents p2p operation commands
type P2pRpc struct {
	engine *Engine
}

type ResultConnectedPeers struct {
	ConnectedPeers []bool
	Total          int
}

func (r *ResultConnectedPeers) String() string {
	var builder strings.Builder
	table := tablewriter.NewWriter(&builder)
	headers := make([]string, len(r.ConnectedPeers)-1)
	for i := range r.ConnectedPeers[1:] {
		headers[i] = strconv.Itoa(i + 1)
	}
	table.SetHeader(headers)
	link := make([]string, len(r.ConnectedPeers)-1)
	for i, connected := range r.ConnectedPeers[1:] {
		if connected {
			link[i] = "X"
		} else {
			link[i] = " "
		}
	}
	table.Append(link)
	table.Render()
	builder.WriteString("Total connected: " + strconv.Itoa(r.Total) + "\n")
	return builder.String() // Print the builder's content
}

func (p *P2pRpc) ConnectedPeers(_ *ArgEmpty, reply *ResultConnectedPeers) error {
	log.Info("RPC request for connected peers") // dunno if could be generated somehow dynamically
	c := 0
	connected := make([]bool, len(p.engine.peers))
	for i, p := range p.engine.peers {
		if p != nil && p.connected {
			connected[i] = true
			c++
		}
	}
	*reply = ResultConnectedPeers{
		ConnectedPeers: connected,
		Total:          c,
	}
	return nil
}

type ResultPingIcmp probing.Statistics

func (r *ResultPingIcmp) String() string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "\n--- %s ping statistics ---\n", r.Addr)
	fmt.Fprintf(&builder, "%d packets transmitted, %d packets received, %v%% packet loss\n",
		r.PacketsSent, r.PacketsRecv, r.PacketLoss)
	fmt.Fprintf(&builder, "round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
		r.MinRtt, r.AvgRtt, r.MaxRtt, r.StdDevRtt)
	return builder.String()
}
func (p *P2pRpc) PingIcmp(args *ArgTarget, reply *ResultPingIcmp) error {
	if args.Target < 0 || args.Target >= len(p.engine.peers) {
		return errInvalidRpcArg
	}
	if p.engine.peers[args.Target] == nil {
		return errTargetNotConnected
	}
	*reply = *(*ResultPingIcmp)(<-pingIcmp(p.engine.peers[args.Target].ip))
	return nil
}

type ResultIcmpAll []*probing.Statistics

func (r *ResultIcmpAll) String() string {
	var builder strings.Builder
	table := tablewriter.NewWriter(&builder)
	headers := make([]string, len(*r))
	headers[0] = ""
	for i := range *r {
		headers[i+1] = strconv.Itoa(i + 1)
	}
	table.SetHeader(headers)
	packetsSent := make([]string, len(*r)+1)
	packetsSent[0] = "PacketsSent"
	PacketsRecv := make([]string, len(*r)+1)
	PacketsRecv[0] = "PacketsRecv"
	PacketLoss := make([]string, len(*r)+1)
	PacketLoss[0] = "PacketLoss"
	MinRtt := make([]string, len(*r)+1)
	MinRtt[0] = "MinRtt"
	AvgRtt := make([]string, len(*r)+1)
	AvgRtt[0] = "AvgRtt"
	MaxRtt := make([]string, len(*r)+1)
	MaxRtt[0] = "MaxRtt"
	StdDevRtt := make([]string, len(*r)+1)
	StdDevRtt[0] = "StdDev"
	for i, result := range *r {
		packetsSent[i+1] = strconv.Itoa(result.PacketsSent)
		PacketsRecv[i+1] = strconv.Itoa(result.PacketsSent)
		PacketLoss[i+1] = strconv.Itoa(result.PacketsSent)
		MinRtt[i+1] = strconv.Itoa(result.PacketsSent)
		AvgRtt[i+1] = strconv.Itoa(result.PacketsSent)
		MaxRtt[i+1] = strconv.Itoa(result.PacketsSent)
		StdDevRtt[i+1] = strconv.Itoa(result.PacketsSent)
	}
	table.Append(packetsSent)
	table.Append(PacketsRecv)
	table.Append(PacketLoss)
	table.Append(MinRtt)
	table.Append(AvgRtt)
	table.Append(MaxRtt)
	table.Append(StdDevRtt)
	table.Render()
	return builder.String() // Print the builder's content
}

func (p *P2pRpc) PingIcmpAll(_ *ArgEmpty, reply *ResultIcmpAll) error {
	replyChannels := make([]<-chan *probing.Statistics, len(p.engine.peers))
	*reply = make([]*probing.Statistics, len(p.engine.peers))
	for i, peer := range p.engine.peers {
		if peer == nil || !peer.connected {
			ch := make(chan *probing.Statistics, 1)
			ch <- &probing.Statistics{} // default result for non-connected peer to write
			replyChannels[i] = ch
			continue
		}
		replyChannels[i] = pingIcmp(peer.ip)
	}
	for i, ch := range replyChannels {
		(*reply)[i] = <-ch
	}
	return nil
}

func (p *P2pRpc) Ping(args *ArgTarget, reply *probing.Statistics) error {
	if args.Target < 0 || args.Target >= len(p.engine.peers) {
		return errInvalidRpcArg
	}
	*reply = *<-pingIcmp(p.engine.peers[args.Target].ip)
	return nil
}
