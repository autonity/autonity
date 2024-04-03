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
	errInvalidRpcArg = errors.New("invalid RPC argument")
)

type ArgTarget struct {
	Target int
}
type ArgEmpty struct {
}

// P2P represents p2p operation commands
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
	*reply = *(*ResultPingIcmp)(<-pingIcmp(p.engine.peers[args.Target].ip))
	return nil
}

func (p *P2pRpc) PingIcmpAll(args *ArgTarget, reply *probing.Statistics) error {
	for i := 0; i < len(p.engine.peers); i++ {
		*reply = *<-pingIcmp(p.engine.peers[args.Target].ip)
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
