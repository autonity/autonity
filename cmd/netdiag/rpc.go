package main

import (
	"bufio"
	"crypto/rand"
	"errors"
	"fmt"
	"math"
	rand2 "math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
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

type ArgSize struct {
	Size int
}

func (a *ArgSize) AskUserInput() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter size (kB) - max 15000: ")
	input, _ := reader.ReadString('\n')
	size, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		fmt.Println("Invalid size.")
		return err
	}
	a.Size = size * 1000
	return nil
}

type ArgCount struct {
	PacketCount int
}

func (a *ArgCount) AskUserInput() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter number of DevP2P packets: ")
	input, _ := reader.ReadString('\n')
	packetCount, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		fmt.Println("Invalid number.")
		return err
	}
	a.PacketCount = packetCount
	return nil
}

type ArgTargetSize struct {
	ArgTarget
	ArgSize
}

func (a *ArgTargetSize) AskUserInput() error {
	if err := a.ArgTarget.AskUserInput(); err != nil {
		return err
	}
	return a.ArgSize.AskUserInput()
}

type ArgSizeCount struct {
	ArgSize
	ArgCount
}

func (a *ArgSizeCount) AskUserInput() error {
	if err := a.ArgSize.AskUserInput(); err != nil {
		return err
	}
	return a.ArgCount.AskUserInput()
}

type ArgTargetSizeCount struct {
	ArgTarget
	ArgSize
	ArgCount
}

func (a *ArgTargetSizeCount) AskUserInput() error {
	if err := a.ArgTarget.AskUserInput(); err != nil {
		return err
	}
	if err := a.ArgCount.AskUserInput(); err != nil {
		return err
	}
	return a.ArgSize.AskUserInput()
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
	headers := make([]string, len(r.ConnectedPeers))
	for i := range r.ConnectedPeers {
		headers[i] = strconv.Itoa(i)
	}
	table.SetHeader(headers)
	link := make([]string, len(r.ConnectedPeers))
	for i, connected := range r.ConnectedPeers {
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

type PacketResult struct {
	TimeReqReceived time.Time
	SyscallDuration time.Duration
	Err             string
}
type ResultSendRandomData struct {
	Id                    int
	Size                  int
	PacketCount           int
	RequestTime           time.Time
	ReceiverReceptionTime time.Time // The one in the ACK
	AckReceivedTime       time.Time // locally
	TotalSyscallDuration  time.Duration
	HasErrors             bool
	PacketResults         []PacketResult
}

func (r *ResultSendRandomData) String() string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "DevP2P Send Rand Data results for target %d:\n", r.Id)
	fmt.Fprintf(&builder, "Packet Size: %dkB\n", r.Size/1000)
	fmt.Fprintf(&builder, "Packet Count: %d\n", r.PacketCount)
	fmt.Fprintf(&builder, "Total Data Size: %dkB\n", (r.PacketCount*r.Size)/1000)
	fmt.Fprintf(&builder, "Request time (RT): %s\n", r.RequestTime.Format(timeFormat))
	if !r.HasErrors {
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
	}
	errorsCount := 0
	for i, res := range r.PacketResults {
		if res.Err != "" {
			fmt.Fprintf(&builder, "Packet #%d: %s\n", i, res.Err)
			errorsCount += 1
		} else {
			fmt.Fprintf(&builder, "Packet #%d: %s\n", i, res.TimeReqReceived.Sub(r.RequestTime))
		}
	}
	fmt.Fprintf(&builder, "\nPacket Loss: %d/%d - %d%% \n", errorsCount, len(r.PacketResults), (errorsCount*100)/len(r.PacketResults))
	return builder.String()
}

func (p *P2POp) SendRandomData(args *ArgTargetSizeCount, reply *ResultSendRandomData) error {
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
	packetResults := make([]PacketResult, args.PacketCount)
	var lastReceived atomic.Uint64
	var hasError atomic.Bool

	for i := 0; i < args.PacketCount; i++ {
		go func(id int) {
			timeReqReceived, syscallDuration, err := peer.sendData(buff)
			packetResults[id] = PacketResult{
				TimeReqReceived: time.Unix(int64(timeReqReceived)/int64(time.Second), int64(timeReqReceived)%int64(time.Second)),
				SyscallDuration: syscallDuration,
				Err:             "",
			}
			if err != nil {
				hasError.Store(true)
				packetResults[id].Err = err.Error()
			}
			lastReceived.Store(timeReqReceived)
			finishedCh <- uint64(syscallDuration.Nanoseconds())
		}(i)
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
	result.HasErrors = hasError.Load()
	result.PacketResults = packetResults
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
			_ = p.SendRandomData(&ArgTargetSizeCount{
				ArgTarget: ArgTarget{args.Target},
				ArgCount:  ArgCount{10},
				ArgSize:   ArgSize{1024},
			}, res)
			time.Sleep(2 * time.Second)
			// measure
			res2 := &ResultSendRandomData{}
			_ = p.SendRandomData(&ArgTargetSizeCount{
				ArgTarget: ArgTarget{args.Target},
				ArgCount:  ArgCount{1},
				ArgSize:   ArgSize{200000},
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
		_ = p.SendRandomData(&ArgTargetSizeCount{
			ArgTarget: ArgTarget{args.Target},
			ArgCount:  ArgCount{1},
			ArgSize:   ArgSize{200000},
		}, res)
		reply.Durations[i] = res.ReceiverReceptionTime.Sub(res.RequestTime)
		time.Sleep(time.Duration(i) * 500 * time.Millisecond)
	}
	return nil
}

type ResultSendRandomDataFrequencyAnalysis struct {
	Target        int
	Size          int
	Delay         []time.Duration
	Duration      []time.Duration
	ReplyDuration []time.Duration
}

func (r *ResultSendRandomDataFrequencyAnalysis) String() string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "Transmissions Analysis Report\n")
	fmt.Fprintf(&builder, "Target Peer: %d\n", r.Target)
	fmt.Fprintf(&builder, "Data size (kb): %d\n", r.Size/1000)
	for i := range r.Delay {
		bandwidth := float64(r.Size) / ((r.Duration[i] - r.ReplyDuration[i]).Seconds() * 1e6)
		fmt.Fprintf(&builder, "%d\tDelay: %s \tDuration: %s\tReplyDuration: %s\tBandwidth: %.6fMB/s\n", i, r.Delay[i], r.Duration[i], r.ReplyDuration[i], bandwidth)
	}
	return builder.String()
}

func (p *P2POp) SendRandomDataFrequencyAnalysis(args *ArgTargetSize, reply *ResultSendRandomDataFrequencyAnalysis) error {
	// TODO: do it across multiple peers, see how it goes
	reply.Delay = make([]time.Duration, 30)
	reply.Duration = make([]time.Duration, 30)
	reply.ReplyDuration = make([]time.Duration, 30)
	k := 0
	for i := 0; i < 10; i++ {
		for j := 0; j < 3; j++ {
			reply.Delay[k] = time.Duration(i) * 250 * time.Millisecond
			time.Sleep(reply.Delay[k])
			res := &ResultSendRandomData{}
			_ = p.SendRandomData(&ArgTargetSizeCount{
				ArgTarget: args.ArgTarget,
				ArgCount:  ArgCount{1},
				ArgSize:   args.ArgSize,
			}, res)
			reply.Duration[k] = res.ReceiverReceptionTime.Sub(res.RequestTime)
			reply.ReplyDuration[k] = time.Now().Sub(res.ReceiverReceptionTime)
			k++
		}
	}
	reply.Target = args.Target
	reply.Size = args.Size
	return nil
}

type ResultSimpleBroadcast struct {
	Size          int
	Count         int
	StartTime     time.Time
	PacketResults [][]PacketResult
}

func (r *ResultSimpleBroadcast) String() string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "Simple Broadcasting Results \n")
	fmt.Fprintf(&builder, "Size: %d\n", r.Size)
	fmt.Fprintf(&builder, "Count: %d\n", r.Count)
	for i := 0; i < r.Count; i++ {
		fmt.Fprintf(&builder, "--- Packet #%d \n", i)
		receivedCount := 0
		aggregateDuration := time.Duration(0)
		peerCount := 0
		for p := 0; p < len(r.PacketResults); p++ {
			if r.PacketResults[p] == nil {
				continue
			}
			res := r.PacketResults[p][i]
			peerCount++
			if res.Err == "" {
				aggregateDuration += res.TimeReqReceived.Sub(r.StartTime)
				receivedCount++
				fmt.Fprintf(&builder, "Peer %d Duration: %s\n", p, res.TimeReqReceived.Sub(r.StartTime))
			} else {
				fmt.Fprintf(&builder, "Peer %d ERROR: %s\n", p, res.Err)
			}
		}
		fmt.Fprintf(&builder, "ACKed by %d/%d, Average Duration %s \n", receivedCount, peerCount, aggregateDuration/time.Duration(receivedCount))
	}

	return builder.String()
}

func (p *P2POp) SimpleBroadcast(args *ArgSizeCount, reply *ResultSimpleBroadcast) error {
	buff := make([]byte, args.Size)
	if _, err := rand.Read(buff); err != nil {
		return err
	}

	results := make([][]PacketResult, len(p.engine.peers))
	startTime := time.Now()
	var wg sync.WaitGroup
	for i := range p.engine.peers {
		if p.engine.peers[i] == nil {
			continue
		}
		results[i] = make([]PacketResult, args.PacketCount)
		wg.Add(1)
		go func(id int) {
			resultsCh := make([]chan any, args.PacketCount)
			for j := 0; j < args.PacketCount; j++ {
				var err error
				resultsCh[j], err = p.engine.peers[id].sendDataAsync(buff)
				if err != nil {
					log.Error("error sending data async", "err", err)
				}
			}
			timer := time.NewTimer(5 * time.Second)
			for j := 0; j < args.PacketCount; j++ {
				select {
				case ans := <-resultsCh[j]:
					replyTime := ans.(AckDataPacket).Time
					results[id][j] = PacketResult{
						TimeReqReceived: time.Unix(int64(replyTime)/int64(time.Second), int64(replyTime)%int64(time.Second)),
						SyscallDuration: 0,
						Err:             "",
					}

				case <-timer.C:
					results[id][j] = PacketResult{Err: "TIMEOUT"}
					timer.Reset(5 * time.Millisecond)
				}
			}
			wg.Done()
		}(i)
	}

	wg.Wait()
	*reply = ResultSimpleBroadcast{
		Size:          args.Size,
		Count:         args.PacketCount,
		PacketResults: results,
		StartTime:     startTime,
	}
	return nil
}

// ***************************
// ***** DISSEMINATE *********
// ***************************
// Explanations:
// Assume p = sqrt(n) where n is the network size
// This mechanism works by splitting N into m groups of m  nodes.
// From each group we pick-up a random node, we call it the group leader.
// We broadcast initially only to the group leaders our message.
// Upon reception, the group leader is responsible to broadcast to the other members of his group the message.
// For added redundancy, after this first phase, we can randomly select some other nodes for a second round.

func (p *P2POp) Disseminate(args *ArgSizeCount, reply *ResultSimpleBroadcast) error {
	buff := make([]byte, args.Size)
	if _, err := rand.Read(buff); err != nil {
		return err
	}
	packetId := rand2.Uint64()
	p.engine.receivedReports[packetId] = make(chan *DisseminateReportPacket, len(p.engine.peers))
	groupSize := int(math.Sqrt(float64(len(p.engine.peers))))
	groupCount := groupSize
	if groupSize*groupCount < len(p.engine.peers) {
		groupCount++
	}
	for i := 0; i < groupCount; i++ {
		var target *Peer
		for target == nil {
			l := rand2.Intn(groupSize)
			target = p.engine.peers[i*groupSize+l]
			// edge cases:
			// no suitable target found in the group to deal with
			// last group size
		}
		target.sendDisseminate(packetId, buff, uint64(p.engine.id), 1)
	}

	// phase 2. select some other random peers
	for {
		select {
		case packet := <-p.engine.receivedReports[packetId]:

		}
	}
	// Wait for reports
	return nil
}

func disseminationGroup(id int, peers []*Peer) []*Peer {
	// todo: create a special object for each propagation strategy to not overload state
	groupSize := int(math.Sqrt(float64(len(peers))))
	groupCount := groupSize
	if groupSize*groupCount < len(peers) {
		groupCount++
	}
	group := make([]*Peer, groupCount)
	myGroup := id % groupSize
	for i := range group {
		group[i] = peers[myGroup*groupSize+i]
	}
	return group // we are returning ourselves so caller be aware
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
