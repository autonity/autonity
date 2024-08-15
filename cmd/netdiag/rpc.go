package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	rand2 "math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"text/tabwriter"
	"time"

	"github.com/olekukonko/tablewriter"
	probing "github.com/prometheus-community/pro-bing"

	"github.com/autonity/autonity/cmd/netdiag/strats"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/p2p"
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
	fmt.Print("Enter target peer index, 0 for all:")
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
	// This tabular view was a bad idea ! convert that into simple rows later ...
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
		peerStats := <-ch
		if p.engine.peers[i] != nil {
			p.engine.peers[i].rtt = peerStats.AvgRtt
			p.engine.state.LatencyMatrix[p.engine.id][i] = peerStats.AvgRtt
		}
		(*reply)[i] = peerStats
	}
	return nil
}

type ResultPing struct {
	Id                    int
	RequestTime           time.Time
	ReceiverReceptionTime time.Time
	PongReceivedTime      time.Time
}

func (r *ResultPing) rtt() time.Duration {
	return r.PongReceivedTime.Sub(r.RequestTime)
}

func (r *ResultPing) String() string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "DevP2P PingDevP2P results for target %d:\n", r.Id)
	fmt.Fprintf(&builder, "Request time: %s\n", r.RequestTime.Format(timeFormat))
	fmt.Fprintf(&builder, "Receiver Reception Time (RCT): %s\n", r.ReceiverReceptionTime.Format(timeFormat))
	fmt.Fprintf(&builder, "Pong Received: %s\n", r.PongReceivedTime.Format(timeFormat))
	RTT := r.rtt()
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

type ResultPingBroadcast struct {
	pings  []ResultPing
	errors []error
}

func (r *ResultPingBroadcast) String() string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "PingDevP2P for all peers\n")
	minRtt := 100 * time.Second
	for i, item := range r.pings {
		if r.errors[i] == nil {
			minRtt = min(minRtt, item.rtt())
		}
	}
	fmt.Fprintf(&builder, "Min RTT: %v\n", minRtt)
	for i, item := range r.pings {
		if r.errors[i] != nil {
			fmt.Fprintf(&builder, "Error for target peer %v : %v\n", item.Id, r.errors[i])
		} else {
			fmt.Fprint(&builder, item.String())
		}
	}
	return builder.String()
}

func (p *P2POp) PingDevP2PBroadcast(_ *ArgEmpty, reply *ResultPingBroadcast) error {
	result := ResultPingBroadcast{
		pings:  make([]ResultPing, 0, len(p.engine.peers)),
		errors: make([]error, 0, len(p.engine.peers)),
	}
	for id, peer := range p.engine.peers {
		if id == p.engine.id {
			continue
		}
		ping := ResultPing{
			Id:          id,
			RequestTime: time.Now(),
		}
		if peer == nil || !peer.connected {
			result.pings = append(result.pings, ping)
			result.errors = append(result.errors, errTargetNotConnected)
		} else {
			timeReceived, err := peer.sendPing()
			if err == nil {
				ping.ReceiverReceptionTime = time.Unix(int64(timeReceived)/int64(time.Second), int64(timeReceived)%int64(time.Second))
				ping.PongReceivedTime = time.Now()
				p.engine.state.LatencyMatrix[p.engine.id][id] = ping.rtt()
			}
			result.pings = append(result.pings, ping)
			result.errors = append(result.errors, err)
		}
	}
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

type ResultSendLatencyArray struct {
	HasError    bool
	AckReceived []bool
	Errors      []string
}

func (r *ResultSendLatencyArray) String() string {
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
		if r.Errors[id] != "" {
			fmt.Fprintf(&builder, "%d with error: %s\n", id, r.Errors[id])
		} else {
			fmt.Fprintf(&builder, "%d\n", id)
		}
	}

	fmt.Fprintf(&builder, "\nFollowing %d targets did not respond:\n", len(ackNotReceived))
	for _, id := range ackNotReceived {
		if r.Errors[id] != "" {
			fmt.Fprintf(&builder, "%d with error: %s\n", id, r.Errors[id])
		} else {
			fmt.Fprintf(&builder, "%d\n", id)
		}
	}
	fmt.Fprintf(&builder, "\n")
	return builder.String()
}

func (p *P2POp) SendLatencyArray(_ *ArgEmpty, reply *ResultSendLatencyArray) error {
	result := ResultSendLatencyArray{}
	errs := make([]string, len(p.engine.peers))
	acks := make([]bool, len(p.engine.peers))
	var hasError atomic.Bool
	var wg sync.WaitGroup

	for i, peer := range p.engine.peers {
		if i == p.engine.id {
			errs[i] = ""
			continue
		}
		if peer == nil || !peer.connected {
			errs[i] = "peer not connected or nil"
			hasError.Store(true)
			continue
		}

		wg.Add(1)
		go func(id int, peer *Peer) {
			_, _, err := peer.sendLatencyArray(p.engine.state.LatencyMatrix[p.engine.id])
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
	symmetric := true
	for id, latencyArray := range *r {
		fmt.Fprintf(&builder, "Latency from peer %d is following:\n", id)
		for peerID, l := range latencyArray {
			if peerID == id {
				if l != 0 {
					valid = false
				}
				continue
			}
			if l == 0 {
				fmt.Fprintf(&builder, "peer %d : inf\n", peerID)
			} else {
				fmt.Fprintf(&builder, "peer %d : %d\n", peerID, l)
			}

			if (*r)[id][peerID] != (*r)[peerID][id] {
				symmetric = false
			}
		}
	}

	fmt.Fprintf(&builder, "Is valid? %v\n", valid)
	fmt.Fprintf(&builder, "Is symmetric? %v\n", symmetric)

	return builder.String()
}

func (p *P2POp) LatencyMatrix(_ *ArgEmpty, reply *ResultLatencyMatrix) error {
	result := make([][]time.Duration, len(p.engine.enodes))
	for i, latencyArray := range p.engine.state.LatencyMatrix {
		result[i] = make([]time.Duration, len(p.engine.enodes))
		copy(result[i], latencyArray)
	}
	*reply = (ResultLatencyMatrix)(result)
	return nil
}

type ResultTCPSocketUpdate struct {
	Reset bool
}

func (r *ResultTCPSocketUpdate) String() string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "Reset: %v\n", r.Reset)

	commands := []string{
		fmt.Sprintf("net.ipv4.tcp_window_scaling=1"),
		fmt.Sprintf("net.core.rmem_max"),
		fmt.Sprintf("net.core.wmem_max"),
		fmt.Sprintf("net.ipv4.tcp_rmem"),
		fmt.Sprintf("net.ipv4.tcp_wmem"),
		fmt.Sprintf("net.ipv4.tcp_slow_start_after_idle"),
	}

	for _, l := range commands {
		localCommand := l
		var out bytes.Buffer
		execCmd := exec.Command("sysctl", localCommand)
		execCmd.Stdout = &out
		err := execCmd.Run()
		if err != nil {
			log.Error(" command failure ", "err", err, "cmd", execCmd)
		}
		builder.WriteString(out.String())
	}
	return builder.String()
}

type ResultTCPSocketTuning struct {
	Target         int
	MinDuration    time.Duration
	Reset          bool
	BufferSize     int
	SizeToDuration map[int]time.Duration
	Durations      []time.Duration
}

func (r *ResultTCPSocketTuning) String() string {
	var builder strings.Builder
	//fmt.Fprintf(&builder, "DevP2P TCP tuning for target %d 200kB:\n", r.Target)
	//fmt.Fprintf(&builder, "Duration: %s\n", r.MinDuration)
	fmt.Fprintf(&builder, "Reset: %v\n", r.Reset)
	//fmt.Fprintf(&builder, "BufferSize: %d\n", r.BufferSize)
	//for size, Duration := range r.SizeToDuration {
	//	fmt.Fprintf(&builder, "Size %d: %s\n", size, Duration)
	//}
	//for i := range r.Durations {
	//	fmt.Fprintf(&builder, "Delay %s: %s\n", 500*time.Millisecond, r.Durations[i])
	//}
	return builder.String()
}

func (p *P2POp) TCPSocketTuning(args *ArgTarget, reply *ResultTCPSocketTuning) error {
	bufferSize := 80 * 1024 * 1024
	peers := make([]*Peer, 0)
	if args.Target == 0 {
		peers = p.engine.peers
	} else {
		peer, err := checkPeer(args.Target, p.engine.peers)
		if err != nil {
			return err
		}
		peers = append(peers, peer)
	}

	p2p.UpdateSystemSocketOptions(bufferSize)
	for i, _ := range peers {
		peer, err := checkPeer(i, p.engine.peers)
		if err != nil {
			continue
		}
		peer.UpdateAppSocketBuffers(bufferSize)
		if err := peer.sendUpdateTcpSocket(bufferSize, false); err != nil {
			log.Error("error sending socket update", "err", err)
		}
	}
	reply = &ResultTCPSocketTuning{Reset: false}
	return nil
}

func (p *P2POp) ResetTCPSocketTuning(args *ArgTarget, reply *ResultTCPSocketTuning) error {
	var peers []*Peer

	if args.Target == 0 {
		peers = p.engine.peers
	} else {
		peer, err := checkPeer(args.Target, p.engine.peers)
		if err != nil {
			return err
		}
		peers = append(peers, peer)
	}
	p2p.ResetSocketOptions()
	bufferSize := 40 * 1024 * 1024
	for i, _ := range peers {
		peer, err := checkPeer(i, p.engine.peers)
		if err != nil {
			continue
		}
		if err := peer.sendUpdateTcpSocket(bufferSize, true); err != nil {
			log.Error("error sending socket update", "err", err)
		}
	}
	reply = &ResultTCPSocketTuning{Reset: true}
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

type ResultDissemination struct {
	Size              int
	MaxPeers          int
	StartTime         time.Time
	IndividualResults []strats.IndividualDisseminateResult
}

func (r *ResultDissemination) String() string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "Disseminate Results \n")
	var results []strats.IndividualDisseminateResult
	for i, res := range r.IndividualResults {
		if res.ErrorTimeout {
			fmt.Println("Error time out")
			continue
		}
		results = append(results, r.IndividualResults[i])
		fmt.Fprintf(&builder, "Peer %d Duration: %s Hops: %d Relay: %d\n", i, res.ReceptionTime.Sub(r.StartTime), res.Hop, res.Relay)
	}
	if len(results) == 0 {
		fmt.Fprintf(&builder, "Dissemination failed")
		return builder.String()
	}
	sort.Slice(results, func(a, b int) bool {
		return results[a].ReceptionTime.Before(results[b].ReceptionTime)
	})
	n := len(results)
	fmt.Fprintf(&builder, "min: %s, median:%s 2/3rd:%s max: %s\n", results[0].ReceptionTime.Sub(r.StartTime), results[n/2].ReceptionTime.Sub(r.StartTime), results[(2*n)/3].ReceptionTime.Sub(r.StartTime), results[n-1].ReceptionTime.Sub(r.StartTime))
	fmt.Fprintf(&builder, "total reports collected: %d\n ", len(results))

	return builder.String()
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
		var (
			packetResults     []PacketResult
			receivedCount     = 0
			aggregateDuration = time.Duration(0)
			peerCount         = 0
		)
		for p := 0; p < len(r.PacketResults); p++ {
			if r.PacketResults[p] == nil {
				continue
			}
			res := r.PacketResults[p][i]
			peerCount++
			if res.Err == "" {
				packetResults = append(packetResults, res)
				aggregateDuration += res.TimeReqReceived.Sub(r.StartTime)
				receivedCount++
				fmt.Fprintf(&builder, "Peer %d Duration: %s\n", p, res.TimeReqReceived.Sub(r.StartTime))
			} else {
				fmt.Fprintf(&builder, "Peer %d ERROR: %s\n", p, res.Err)
			}
		}
		fmt.Fprintf(&builder, "ACKed by %d/%d, Average Duration %s \n", receivedCount, peerCount, aggregateDuration/time.Duration(receivedCount))
		// min average 2/3rd max
		sort.Slice(packetResults, func(a, b int) bool {
			return packetResults[a].TimeReqReceived.Before(packetResults[b].TimeReqReceived)
		})
		n := len(packetResults)
		fmt.Fprintf(&builder, "min: %s, median:%s 2/3rd:%s max: %s \n", packetResults[0].TimeReqReceived.Sub(r.StartTime), packetResults[n/2].TimeReqReceived.Sub(r.StartTime), packetResults[(2*n)/3].TimeReqReceived.Sub(r.StartTime), packetResults[n-1].TimeReqReceived.Sub(r.StartTime))
	}

	return builder.String()
}

type ArgStrategy struct {
	Strategy int
}

func (a *ArgStrategy) AskUserInput() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Available Dissemination Strategies: ")
	for i, s := range strats.StrategyRegistry {
		fmt.Printf("[%d] %s\n", i, s.Name)
	}
	fmt.Print("Chose strategy: ")
	input, _ := reader.ReadString('\n')
	index, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || index >= len(strats.StrategyRegistry) || index < 0 {
		fmt.Println("Invalid strategy selected")
		return err
	}
	a.Strategy = index
	return nil
}

type ArgMaxPeers struct {
	MaxPeers int
}

func (a *ArgMaxPeers) AskUserInput() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Number of peers to disseminate, 0 for all: ")
	input, _ := reader.ReadString('\n')
	count, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		fmt.Println("Invalid number")
		return err
	}
	a.MaxPeers = count
	return nil
}

type ArgDisseminate struct {
	ArgStrategy
	ArgMaxPeers
	ArgSize
}

func (a *ArgDisseminate) AskUserInput() error {
	if err := a.ArgStrategy.AskUserInput(); err != nil {
		return err
	}
	if err := a.ArgMaxPeers.AskUserInput(); err != nil {
		return err
	}
	if err := a.ArgSize.AskUserInput(); err != nil {
		return err
	}
	return nil
}

type ArgWarmUp struct {
	ArgMaxPeers
	ArgSize
}

func (a *ArgWarmUp) AskUserInput() error {
	if err := a.ArgMaxPeers.AskUserInput(); err != nil {
		return err
	}
	if err := a.ArgSize.AskUserInput(); err != nil {
		return err
	}
	return nil
}

func (p *P2POp) WarmUp(args *ArgWarmUp, reply *ResultDissemination) error {
	// Warm up - tcp slow start
	buff := make([]byte, args.Size)
	if _, err := rand.Read(buff); err != nil {
		return err
	}
	reply.StartTime = time.Now()
	reply.Size = args.Size
	recipients := args.MaxPeers
	if args.MaxPeers == 0 {
		recipients = len(p.engine.peers)
	}
	reply.MaxPeers = recipients
	rnd := rand2.New(rand2.NewSource(time.Now().UnixNano()))
	packetId := rnd.Uint64()
	p.engine.state.ReceivedReports[packetId] = make(chan *strats.IndividualDisseminateResult)
	log.Info("Started Dissemination", "size", reply.Size, "peers", recipients, "packetId", packetId)

	// broadcast strategy
	if err := p.engine.strategies[0].Execute(packetId, buff, recipients); err != nil {
		return err
	}
	reply.IndividualResults = p.engine.state.CollectReports(packetId, recipients)
	return nil
}

func (p *P2POp) Disseminate(args *ArgDisseminate, reply *ResultDissemination) error {

	buff := make([]byte, args.Size)
	if _, err := rand.Read(buff); err != nil {
		return err
	}
	reply.StartTime = time.Now()
	reply.Size = args.Size
	recipients := args.MaxPeers
	if args.MaxPeers == 0 {
		recipients = len(p.engine.peers)
	}
	reply.MaxPeers = recipients
	rnd := rand2.New(rand2.NewSource(time.Now().UnixNano()))
	packetId := rnd.Uint64()
	p.engine.state.ReceivedReports[packetId] = make(chan *strats.IndividualDisseminateResult)
	log.Info("Started Dissemination", "size", reply.Size, "peers", recipients, "packetId", packetId)

	// actual dissemination
	if err := p.engine.strategies[args.Strategy].Execute(packetId, buff, recipients); err != nil {
		return err
	}
	reply.IndividualResults = p.engine.state.CollectReports(packetId, recipients)
	return nil
}

type ArgGraphConstruct struct {
	ArgStrategy
	ArgMaxPeers
}

func (a *ArgGraphConstruct) AskUserInput() error {
	if err := a.ArgStrategy.AskUserInput(); err != nil {
		return err
	}
	if err := a.ArgMaxPeers.AskUserInput(); err != nil {
		return err
	}
	return nil
}

func (p *P2POp) ConstructGraph(arg *ArgGraphConstruct, _ *ArgEmpty) error {
	recipients := arg.MaxPeers
	if arg.MaxPeers == 0 {
		recipients = len(p.engine.peers)
	}
	return p.engine.strategies[arg.Strategy].ConstructGraph(recipients)
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
