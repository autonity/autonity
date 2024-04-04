package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"text/tabwriter"

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
	fmt.Fprintf(&builder, "%s ping statistics\n", r.Addr)
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
	var (
		builder     strings.Builder
		header      = ""
		packetsSent = "PacketsSent"
		PacketsRecv = "PacketsRecv"
		PacketLoss  = "PacketLoss"
		MinRtt      = "MinRtt"
		AvgRtt      = "AvgRtt"
		MaxRtt      = "MaxRtt"
		StdDevRtt   = "StdDevRtt"
	)
	table := tabwriter.NewWriter(&builder, 0, 8, 2, '\t', 0)
	for i := range *r {
		header += "\t" + strconv.Itoa(i)
	}
	table.Write([]byte(header + "\n"))
	for _, result := range *r {
		packetsSent += "\t" + strconv.Itoa(result.PacketsSent)
		PacketsRecv += "\t" + strconv.Itoa(result.PacketsRecv)
		PacketLoss += "\t" + strconv.FormatFloat(result.PacketLoss, 'f', 2, 64)
		MinRtt += "\t" + result.MinRtt.String()
		AvgRtt += "\t" + result.AvgRtt.String()
		MaxRtt += "\t" + result.MaxRtt.String()
		StdDevRtt += "\t" + result.StdDevRtt.String()
	}
	table.Write([]byte(packetsSent + "\n"))
	table.Write([]byte(PacketsRecv + "\n"))
	table.Write([]byte(PacketLoss + "\n"))
	table.Write([]byte(MinRtt + "\n"))
	table.Write([]byte(AvgRtt + "\n"))
	table.Write([]byte(MaxRtt + "\n"))
	table.Write([]byte(StdDevRtt + "\n"))
	table.Flush()
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
