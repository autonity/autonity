package main

import (
	"bufio"
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"text/tabwriter"
	"time"

	"github.com/olekukonko/tablewriter"
	probing "github.com/prometheus-community/pro-bing"

	"github.com/autonity/autonity/log"
)

const (
	timeFormat = "15:04:05.000"
)

var (
	errInvalidRpcArg      = errors.New("invalid RPC argument")
	errTargetNotConnected = errors.New("target peer not connected")
)

type Argument interface {
	AskUserInput() error
}

type ArgTarget struct {
	Target int
}

func (a *ArgTarget) AskUserInput() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter target peer index: ")
	input, _ := reader.ReadString('\n')
	targetIndex, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		fmt.Println("Invalid target index.")
		return err
	}
	a.Target = targetIndex
	return nil
}

type ArgTargetDataSize struct {
	Target      int
	PacketCount int
	Size        int
}

func (a *ArgTargetDataSize) AskUserInput() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter target index: ")
	input, _ := reader.ReadString('\n')
	targetIndex, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		fmt.Println("Invalid target index.")
		return err
	}
	a.Target = targetIndex
	fmt.Print("Enter number of DevP2P packets: ")
	input, _ = reader.ReadString('\n')
	packetCount, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		fmt.Println("Invalid number.")
		return err
	}
	a.PacketCount = packetCount
	fmt.Print("Enter size (kB) - max 15000: ")
	input, _ = reader.ReadString('\n')
	size, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		fmt.Println("Invalid size.")
		return err
	}
	a.Target = targetIndex
	a.Size = size * 1000
	return nil
}

type ArgEmpty struct {
}

// P2POp represents p2p operation commands
type P2POp struct {
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

func (p *P2POp) ConnectedPeers(_ *ArgEmpty, reply *ResultConnectedPeers) error {
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
func (p *P2POp) PingIcmp(args *ArgTarget, reply *ResultPingIcmp) error {
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

func (p *P2POp) PingIcmpBroadcast(_ *ArgEmpty, reply *ResultIcmpAll) error {
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

type ResultPing struct {
	Id                    int
	RequestTime           time.Time
	ReceiverReceptionTime time.Time
	PongReceivedTime      time.Time
}

func (r *ResultPing) String() string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "DevP2P PingDevP2P results for target %d:\n", r.Id)
	fmt.Fprintf(&builder, "Request time: %s\n", r.RequestTime.Format(timeFormat))
	fmt.Fprintf(&builder, "Receiver Reception Time (RCT): %s\n", r.ReceiverReceptionTime.Format(timeFormat))
	fmt.Fprintf(&builder, "Pong Received: %s\n", r.PongReceivedTime.Format(timeFormat))
	RTT := r.PongReceivedTime.Sub(r.RequestTime)
	fmt.Fprintf(&builder, "RTT: %s\n", RTT)
	theoryReceptionTimestamp := r.RequestTime.Add(RTT / 2)
	fmt.Fprintf(&builder, "Theoretical Reception Timestamp (TRT): %s\n", theoryReceptionTimestamp.Format(timeFormat))
	if theoryReceptionTimestamp.After(r.ReceiverReceptionTime) {
		fmt.Fprintf(&builder, "Delta TRT/RCT: %s\n", theoryReceptionTimestamp.Sub(r.ReceiverReceptionTime))
	} else {
		fmt.Fprintf(&builder, "Delta RCT/TRT: %s\n", r.ReceiverReceptionTime.Sub(theoryReceptionTimestamp))
	}
	return builder.String()
}

func (p *P2POp) PingDevP2P(args *ArgTarget, reply *ResultPing) error {
	if args.Target < 0 || args.Target >= len(p.engine.peers) {
		return errInvalidRpcArg
	}
	peer := p.engine.peers[args.Target]
	if peer == nil || !peer.connected {
		return errTargetNotConnected
	}
	result := ResultPing{
		Id:          args.Target,
		RequestTime: time.Now(),
	}
	timeReceived, err := peer.sendPing()
	if err != nil {
		return err
	}
	result.ReceiverReceptionTime = time.Unix(int64(timeReceived)/int64(time.Second), int64(timeReceived)%int64(time.Second))
	result.PongReceivedTime = time.Now()
	*reply = result
	return nil
}

type ResultSendRandomData struct {
	Id                    int
	Size                  int
	PacketCount           int
	RequestTime           time.Time
	ReceiverReceptionTime time.Time // The one in the ACK
	AckReceivedTime       time.Time // locally
	TotalSyscallDuration  time.Duration
}

func (r *ResultSendRandomData) String() string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "DevP2P Send Rand Data results for target %d:\n", r.Id)
	fmt.Fprintf(&builder, "Packet Size: %dkB\n", r.Size/1000)
	fmt.Fprintf(&builder, "Packet Count: %d\n", r.PacketCount)
	fmt.Fprintf(&builder, "Total Data Size: %dkB\n", (r.PacketCount*r.Size)/1000)
	fmt.Fprintf(&builder, "Request time (RT): %s\n", r.RequestTime.Format(timeFormat))
	fmt.Fprintf(&builder, "Last Reception Time (RCT): %s\n", r.ReceiverReceptionTime.Format(timeFormat))
	fmt.Fprintf(&builder, "All ACK Received: %s\n", r.AckReceivedTime.Format(timeFormat))
	duration := r.ReceiverReceptionTime.Sub(r.RequestTime)
	fmt.Fprintf(&builder, "Duration RCT-RT: %s\n", duration)
	fmt.Fprintf(&builder, "Total syscall wait: %s\n", r.TotalSyscallDuration)

	bandwithWithLatency := float64(r.PacketCount*r.Size) / (duration.Seconds() * 1e6)
	// duration is here TimeOfAdvertisedReception
	fmt.Fprintf(&builder, "Bandwidth with latency : %.6f MB/s\n", bandwithWithLatency)
	durationWithoutLatency := duration - (r.AckReceivedTime.Sub(r.ReceiverReceptionTime))
	bandwithWithoutLatency := float64(r.PacketCount*r.Size) / (durationWithoutLatency.Seconds() * 1e6)
	fmt.Fprintf(&builder, "Bandwidth without latency : %.6f MB/s\n", bandwithWithoutLatency)
	return builder.String()
}

func (p *P2POp) SendRandomData(args *ArgTargetDataSize, reply *ResultSendRandomData) error {
	if args.Target < 0 || args.Target >= len(p.engine.peers) {
		return errInvalidRpcArg
	}
	peer := p.engine.peers[args.Target]
	if peer == nil || !peer.connected {
		return errTargetNotConnected
	}
	buff := make([]byte, args.Size)
	if _, err := rand.Read(buff); err != nil {
		return err
	}
	result := ResultSendRandomData{
		Id:          args.Target,
		PacketCount: args.PacketCount,
		Size:        args.Size,
		RequestTime: time.Now(),
	}
	finishedCh := make(chan uint64, 1)
	var lastReceived atomic.Uint64
	var hasError atomic.Value
	for i := 0; i < args.PacketCount; i++ {
		go func() {
			timeReceived, syscallDuration, err := peer.sendData(buff)
			if err != nil {
				hasError.Store(true)
			}
			lastReceived.Store(timeReceived)
			finishedCh <- uint64(syscallDuration.Nanoseconds())
		}()
	}
	var totalSyscallDuration uint64
	for i := 0; i < args.PacketCount; i++ {
		finishedTime := <-finishedCh
		if finishedTime > totalSyscallDuration {
			totalSyscallDuration = finishedTime
		}
	}
	timeReceived := lastReceived.Load()
	result.TotalSyscallDuration = time.Duration(totalSyscallDuration)
	result.ReceiverReceptionTime = time.Unix(int64(timeReceived)/int64(time.Second), int64(timeReceived)%int64(time.Second))
	result.AckReceivedTime = time.Now()
	*reply = result
	return nil
}

type ResultTCPSocketTuning struct {
	Target      int
	MinDuration time.Duration
	NoDelay     bool
	BufferSize  int
	Durations   []time.Duration
}

func (r *ResultTCPSocketTuning) String() string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "DevP2P TCP tuning for target %d 200kB:\n", r.Target)
	fmt.Fprintf(&builder, "Duration: %s\n", r.MinDuration)
	fmt.Fprintf(&builder, "NoDelay: %v\n", r.NoDelay)
	fmt.Fprintf(&builder, "BufferSize: %d\n", r.BufferSize)
	for i := range r.Durations {
		fmt.Fprintf(&builder, "Delay %s: %s\n", time.Duration(i)*500*time.Millisecond, r.Durations[i])
	}
	return builder.String()
}

func (p *P2POp) TCPSocketTuning(args *ArgTarget, reply *ResultTCPSocketTuning) error {
	peer, err := checkPeer(args.Target, p.engine.peers)
	if err != nil {
		return err
	}
	bufferSizes := []int{1024, 2 * 1024, 4 * 1024, 8 * 1024, 16 * 1024, 32 * 1024, 64 * 1024, 128 * 1024, 256 * 1024, 512 * 1024, 1024 * 1024}
	minDuration := 99 * time.Second
	minNoDelay := false
	minBufferSize := 0
	for _, noDelay := range []bool{false, true} {
		for _, buffSize := range bufferSizes {
			if err := peer.sendUpdateTcpSocket(buffSize, noDelay); err != nil {
				log.Error("error sending update", "err", err)
			}
			peer.UpdateSocketOptions(buffSize, noDelay)
			// warmup
			res := &ResultSendRandomData{}
			_ = p.SendRandomData(&ArgTargetDataSize{
				Target:      args.Target,
				PacketCount: 10,
				Size:        1024,
			}, res)
			time.Sleep(2 * time.Second)
			// measure
			res2 := &ResultSendRandomData{}
			_ = p.SendRandomData(&ArgTargetDataSize{
				Target:      args.Target,
				PacketCount: 1,
				Size:        200000,
			}, res2)
			duration := res2.ReceiverReceptionTime.Sub(res2.RequestTime)
			if duration < minDuration {
				minDuration = duration
				minNoDelay = noDelay
				minBufferSize = buffSize
			}
		}
	}
	*reply = ResultTCPSocketTuning{
		Target:      args.Target,
		MinDuration: minDuration,
		NoDelay:     minNoDelay,
		BufferSize:  minBufferSize,
		Durations:   make([]time.Duration, 10),
	}
	for i := 0; i < 10; i++ {
		res := &ResultSendRandomData{}
		_ = p.SendRandomData(&ArgTargetDataSize{
			Target:      args.Target,
			PacketCount: 1,
			Size:        200000,
		}, res)
		reply.Durations[i] = res.ReceiverReceptionTime.Sub(res.RequestTime)
		time.Sleep(time.Duration(i) * 500 * time.Millisecond)
	}
	return nil
}

func checkPeer(id int, peers []*Peer) (*Peer, error) {
	if id < 0 || id >= len(peers) {
		return nil, errInvalidRpcArg
	}
	peer := peers[id]
	if peer == nil || !peer.connected {
		return nil, errTargetNotConnected
	}
	return peer, nil
}
