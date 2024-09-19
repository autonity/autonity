package core

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"sync"
	"time"

	probing "github.com/prometheus-community/pro-bing"

	"github.com/autonity/autonity/cmd/netdiag/strats"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/p2p"
	"github.com/autonity/autonity/rlp"
)

// `ProtocolMessages` should be last
const (
	PingMsg = uint64(iota)
	PongMsg
	DataMsg
	AckDataMsg
	UpdateTCPSocket
	DisseminateRequest
	DisseminateReport
	LatencyArrayMsg
	TriggerRequest
	GraphReady
	ProtocolMessages
)

var protocolHandlers = map[uint64]func(e *Engine, p *Peer, data io.Reader) error{
	PingMsg:            handlePing,
	PongMsg:            handlePong,
	DataMsg:            handleData, // random string of bytes
	AckDataMsg:         handleAckData,
	UpdateTCPSocket:    handleUpdateTcpSocket,
	DisseminateRequest: handleDisseminatePacket,
	DisseminateReport:  handleDisseminateReport,
	LatencyArrayMsg:    handleLatencyArray,
	TriggerRequest:     handleTriggerRequest,
	GraphReady:         handleGraphReady,
	//BlockMsg //serialized block
	//AckBlockMsg
}

var (
	errInvalidMsgCode   = errors.New("invalid message code")
	errUnknownRequest   = errors.New("no matching request id")
	errPeerNotConnected = errors.New("peer not connected")
	errPacketAlreadyRcv = errors.New("packet already received")
	errTimeout          = errors.New("request timed out")
)

type Peer struct {
	*p2p.Peer
	p2p.MsgReadWriter
	// ICMP stats
	// Delay on TIME
	Ip        string
	Rtt       time.Duration
	requests  map[uint64]chan any
	Connected bool
	id        int
	sync.RWMutex
}

type DisseminatePacket struct {
	StrategyCode   uint64
	RequestId      uint64
	OriginalSender uint64
	MaxPeers       uint64
	Hop            uint8
	Partial        bool
	Seq            uint16
	Total          uint16
	Data           []byte
}

func (p *Peer) DisseminateRequest(code uint64, requestId uint64, hop uint8, originalSender uint64, maxPeers uint64, data []byte, partial bool, seqNum, total uint16) error {
	var chunks []DisseminatePacket
	maxSize := 20_000
	log.Info("Disseminate Packet", "requestId", requestId, "originalSender", originalSender)
	if (p.IsUDP() || p.IsQuic()) && len(data) > maxSize {
		snum := 0
		total := len(data) / maxSize
		if len(data)%maxSize > 0 {
			total++
		}
		for i := 0; i < len(data); i += maxSize {
			packet := DisseminatePacket{
				StrategyCode:   code,
				RequestId:      requestId,
				OriginalSender: originalSender,
				MaxPeers:       maxPeers,
				Hop:            hop,
				Data:           data[i:min(i+maxSize, len(data))],
				Partial:        true,
				Seq:            uint16(snum),
				Total:          uint16(total),
			}
			snum++
			chunks = append(chunks, packet)
		}
	} else {
		packet := DisseminatePacket{
			StrategyCode:   code,
			RequestId:      requestId,
			OriginalSender: originalSender,
			MaxPeers:       maxPeers,
			Hop:            hop,
			Data:           data,
			Partial:        partial,
			Seq:            seqNum,
			Total:          total,
		}
		chunks = append(chunks, packet)
	}

	if p == nil {
		log.Error("p is nil ??")
	}

	now := time.Now()
	for _, chunk := range chunks {
		packet := chunk
		err := p2p.Send(p.MsgReadWriter, DisseminateRequest, packet)
		if err != nil {
			log.Error("Disseminate Request", "err", err)
		}
	}
	log.Info("Total time to send", "num chunks", len(chunks), "time taken", time.Since(now).Microseconds())

	return nil
}

func (p *Peer) dispatchResponse(requestId uint64, packet any) error {
	p.RLock()
	req, ok := p.requests[requestId]
	p.RUnlock()
	if !ok {
		log.Error("Unknown request id", "id", requestId)
		return nil
	}
	req <- packet // what if timeout here?
	p.Lock()
	delete(p.requests, requestId)
	p.Unlock()
	return nil
}

func (p *Peer) RTT() time.Duration {
	return p.Rtt
}

func (p *Peer) dispatchRequest(requestId uint64, code uint64, packet any) (chan any, error) {
	responseCh := make(chan any, 1)
	p.Lock()
	p.requests[requestId] = responseCh
	p.Unlock()
	return responseCh, p2p.Send(p, code, packet)
}

// **** PING *****

type PingPacket struct {
	RequestId uint64
}

type PongPacket struct {
	RequestId uint64
	Time      uint64
}

func (p *Peer) SendPing() (uint64, error) {
	id := rand.Uint64()
	req, err := p.dispatchRequest(id, PingMsg, PingPacket{id})
	if err != nil {
		return 0, err
	}
	fmt.Println("[PING] >>", "requestId", id)
	timer := time.NewTimer(5 * time.Second)
	select {
	case ans := <-req:
		return ans.(PongPacket).Time, nil
	case <-timer.C:
		return 0, errTimeout
	}
}

func handlePing(_ *Engine, p *Peer, data io.Reader) error {
	now := uint64(time.Now().UnixNano())
	var ping PingPacket
	if err := rlp.Decode(data, &ping); err != nil {
		return err
	}
	fmt.Println("[PING] << , [PONG] >> ", "requestId", ping.RequestId)
	return p2p.Send(p, PongMsg, PongPacket{ping.RequestId, now})
}

func handlePong(_ *Engine, p *Peer, msg io.Reader) error {
	var pong PongPacket
	if err := rlp.Decode(msg, &pong); err != nil {
		return err
	}
	log.Debug("[PONG] <<", "requestId", pong.RequestId)
	return p.dispatchResponse(pong.RequestId, pong)
}

// ********************
// ***** SENDDATA *****
// ********************

type DataPacket struct {
	RequestId uint64
	Data      []byte
}

type AckDataPacket struct {
	RequestId uint64
	Time      uint64
}

func (p *Peer) SendData(data []byte) (uint64, time.Duration, error) {
	startTime := time.Now()
	id := rand.Uint64()
	log.Debug("Sending data to peer", "requestId", id, "peerId", p.id)
	req, err := p.dispatchRequest(id, DataMsg, DataPacket{id, data})
	if err != nil {
		log.Error("Unable to dispatch request", "error", err)
		return 0, 0, err
	}
	dispatchDuration := time.Since(startTime)
	timer := time.NewTimer(5 * time.Second)
	select {
	case ans := <-req:
		return ans.(AckDataPacket).Time, dispatchDuration, nil
	case <-timer.C:
		return 0, dispatchDuration, errTimeout
	}
}

func (p *Peer) sendDataAsync(data []byte) (chan any, error) {
	id := rand.Uint64()
	log.Debug("Sending data to peer (async) ", "requestId", id, "peerId", p.id)
	return p.dispatchRequest(id, DataMsg, DataPacket{id, data})
}

func handleData(_ *Engine, p *Peer, data io.Reader) error {
	var dataPacket DataPacket
	if err := rlp.Decode(data, &dataPacket); err != nil {
		return err
	}
	now := uint64(time.Now().UnixNano()) // <-- We could add a timestamp before decoding too ?
	log.Debug("Received data packet", "requestId", dataPacket.RequestId, "fromPeer", p.id)
	return p2p.Send(p, AckDataMsg, AckDataPacket{dataPacket.RequestId, now})
}

func handleAckData(_ *Engine, p *Peer, msg io.Reader) error {
	var ack AckDataPacket
	if err := rlp.Decode(msg, &ack); err != nil {
		return err
	}
	log.Debug("[ACKDATA] << ", "requestId", ack.RequestId, "fromPeer", p.id)
	return p.dispatchResponse(ack.RequestId, ack)
}

// ***********************
// ***** SendLatency *****
// ***********************

type LatencyArrayPacket struct {
	RequestId    uint64
	Strategy     uint64
	LatencyArray []uint64
}

func (p *Peer) SendLatencyArray(strategy uint64, latency []time.Duration) (uint64, time.Duration, error) {
	startTime := time.Now()
	id := rand.Uint64()
	log.Debug("[LatencyArrayPacket] >> ", "id", id)
	buff := make([]uint64, len(latency))
	for i, l := range latency {
		buff[i] = uint64(l)
	}
	req, err := p.dispatchRequest(id, LatencyArrayMsg, LatencyArrayPacket{id, strategy, buff})
	if err != nil {
		log.Error("Unable to dispatch request", "error", err)
		return 0, 0, err
	}
	dispatchDuration := time.Since(startTime)
	timer := time.NewTimer(5 * time.Second)
	select {
	case ans := <-req:
		return ans.(AckDataPacket).Time, dispatchDuration, nil
	case <-timer.C:
		return 0, dispatchDuration, errTimeout
	}
}

// handleLatencyArray is called when a peer receives a latency array
// It will store the latency array in the state and construct the graph if possible
func handleLatencyArray(e *Engine, p *Peer, data io.Reader) error {
	var latencyPacket LatencyArrayPacket
	if err := rlp.Decode(data, &latencyPacket); err != nil {
		return err
	}
	now := uint64(time.Now().UnixNano()) // <-- We could add a timestamp before decoding too ?
	log.Debug("[LatencyArrayPacket] << ", "id", latencyPacket.RequestId)

	latencyType, _ := e.Strategies[latencyPacket.Strategy].LatencyType()
	var latencyLen int
	if latencyType == strats.LatencyTypeRelative {
		latencyLen = len(e.Peers)
	} else {
		latencyLen = len(NtpServers)
	}

	e.State.LatencyMatrix[p.id] = make([]time.Duration, latencyLen)

	for i, l := range latencyPacket.LatencyArray {
		if i >= latencyLen {
			log.Error("Invalid latency array received from peer", "peer", p.id, "len", len(e.Peers), "index", i)
			break
		}
		e.State.LatencyMatrix[p.id][i] = time.Duration(l)
	}

	go func(e *Engine, strategy uint64) {
		if err := e.Strategies[int(strategy)].ConstructGraph(len(e.Peers)); err != nil {
			if !errors.Is(err, strats.ErrLatencyMatrixNotReady) {
				log.Error("Error constructing graph", "error", err)
			} else {
				log.Debug("Received latency array, but latency matrix not complete")
			}
		} else {
			if err := BroadcastGraphReady(e, strategy); err != nil {
				log.Error("Error in broadcast graph ready", "error", err)
			}
		}
	}(e, latencyPacket.Strategy)
	log.Debug("Got latency array from peer", "peerId", p.ID())
	return p2p.Send(p, AckDataMsg, AckDataPacket{latencyPacket.RequestId, now})
}

// ***********************
// ***** SendTrigger *****
// ***********************

type TriggerPacket struct {
	RequestId uint64
	Strategy  uint64
}

func (p *Peer) SendTriggerRequest(strategy uint64) error {
	id := rand.Uint64()
	log.Debug("[TriggerPacket][TriggerRequest] >> ", "id", id)
	req, err := p.dispatchRequest(
		id,
		TriggerRequest,
		TriggerPacket{RequestId: id, Strategy: strategy},
	)
	if err != nil {
		log.Error("Unable to dispatch request", "error", err)
		return err
	}
	timer := time.NewTimer(5 * time.Second)
	select {
	case <-req:
		return nil
	case <-timer.C:
		return errTimeout
	}
}

// handleTriggerRequest is called when a peer receives a trigger request
// It will ping all peers and broadcast it's own latency array
func handleTriggerRequest(e *Engine, p *Peer, data io.Reader) error {
	var trigger TriggerPacket
	if err := rlp.Decode(data, &trigger); err != nil {
		return err
	}
	now := uint64(time.Now().UnixNano()) // <-- We could add a timestamp before decoding too ?

	latencyType, _ := e.Strategies[trigger.Strategy].LatencyType()
	go func(e *Engine, p *Peer, trigger *TriggerPacket) {
		// get latency for all peers or fixed server set
		var latency []probing.Statistics
		if latencyType == strats.LatencyTypeRelative {
			latency = PingPeers(e)
		} else {
			latency = PingFixedNTP()
		}

		// set our own latency in the matrix
		e.State.LatencyMatrix[e.Id] = FilterAveRtt(latency, latencyType)

		if err := BroadcastLatency(e, trigger.Strategy, latency); err != nil {
			log.Error("Error in broadcast latency", "error", err)
		}
	}(e, p, &trigger)
	log.Debug("[TriggerPacket][TriggerRequest] << ", trigger.RequestId)
	return p2p.Send(p, AckDataMsg, AckDataPacket{trigger.RequestId, now})
}

func (p *Peer) sendGraphReady(strategy uint64) error {
	id := rand.Uint64()
	log.Debug("[TriggerPacket][GraphReady] >> ", "id", id)
	req, err := p.dispatchRequest(
		id, GraphReady,
		TriggerPacket{RequestId: id, Strategy: strategy},
	)
	if err != nil {
		log.Error("Unable to dispatch request", "error", err)
		return err
	}
	timer := time.NewTimer(5 * time.Second)
	select {
	case <-req:
		return nil
	case <-timer.C:
		return errTimeout
	}
}

func handleGraphReady(e *Engine, p *Peer, data io.Reader) error {
	var trigger TriggerPacket
	if err := rlp.Decode(data, &trigger); err != nil {
		return err
	}
	e.Strategies[trigger.Strategy].GraphReadyForPeer(p.id)
	now := uint64(time.Now().UnixNano()) // <-- We could add a timestamp before decoding too ?
	log.Debug("[TriggerPacket][GraphReady] << ", "id", trigger.RequestId)
	return p2p.Send(p, AckDataMsg, AckDataPacket{trigger.RequestId, now})
}

// ***************************
// ***** UpdateTCPSocket *****
// ***************************

type TCPOptionsPacket struct {
	BufferSize uint64
	Reset      bool
}

func (p *Peer) SendUpdateTcpSocket(bufferSize int, reset bool) error {
	id := rand.Uint64()
	log.Info("sending update")
	_, err := p.dispatchRequest(id, UpdateTCPSocket, TCPOptionsPacket{uint64(bufferSize), reset}) // this will leak for now
	return err
}

func handleUpdateTcpSocket(e *Engine, p *Peer, msg io.Reader) error {
	log.Info("received handle update")
	var opts TCPOptionsPacket
	if err := rlp.Decode(msg, &opts); err != nil {
		log.Error("update tcp socket failure", "error", err)
		return err
	}
	if opts.Reset {
		p2p.ResetSocketOptions()
	} else {
		p2p.UpdateSystemSocketOptions(int(opts.BufferSize))
	}
	for i, _ := range e.Peers {
		peer, err := checkPeer(i, e.Peers)
		if err != nil {
			continue
		}
		peer.UpdateAppSocketBuffers(int(opts.BufferSize))
	}
	return nil
}

type DisseminateReportPacket struct {
	RequestId uint64
	Sender    uint64
	Hop       uint8
	Time      uint64
	Full      bool
}

// Returns `true` if the packet is a duplicate
func readDisseminateChunk(e *Engine, packet *DisseminatePacket) bool {
	e.RLock()
	defer e.RUnlock()

	if pktInfo, ok := e.State.ReceivedPacketsFor(packet.RequestId); ok {
		// check if the seqNum is already received
		if len(pktInfo.SeqReceived) <= int(packet.Seq) {
			log.Crit("invalid seqNum", "packet Info", pktInfo, "received packet", packet)
		}

		if !packet.Partial || (packet.Partial && pktInfo.SeqReceived[packet.Seq]) {
			//log.Info("packet has already arrived", "packet", packet.Seq)
			// do nothing
			return true
		}
	}
	return false
}

func cacheDisseminatePacket(e *Engine, packet *DisseminatePacket) error {
	e.Lock()
	defer e.Unlock()
	chunkInfo, ok := e.State.ReceivedPacketsFor(packet.RequestId)
	if ok {
		if chunkInfo.SeqReceived[packet.Seq] {
			return errPacketAlreadyRcv
		}
		chunkInfo.SeqReceived[packet.Seq] = true
		chunkInfo.TotalReceived++
		if chunkInfo.TotalReceived == len(chunkInfo.SeqReceived) {
			chunkInfo.Partial = false
		}
		e.State.SetReceivedPacketsFor(packet.RequestId, chunkInfo)
	} else {
		chunkInfo := strats.ChunkInfo{
			TotalReceived: 1,
			Partial:       packet.Partial,
			SeqReceived:   make([]bool, max(packet.Total, 1)),
		}
		// use max.Int to mark the seq as received
		chunkInfo.SeqReceived[packet.Seq] = true
		if chunkInfo.TotalReceived == len(chunkInfo.SeqReceived) {
			chunkInfo.Partial = false
		}
		e.State.SetReceivedPacketsFor(packet.RequestId, chunkInfo)
	}
	return nil
}

func handleDisseminatePacket(e *Engine, p *Peer, data io.Reader) error {
	var packet DisseminatePacket
	if err := rlp.Decode(data, &packet); err != nil {
		return err
	}

	now := uint64(time.Now().UnixNano()) // <-- We could add a timestamp before decoding too ?
	// check if first time received.
	if readDisseminateChunk(e, &packet) {
		return nil
	}
	log.Info("Received packet", "from", p.id, "requestId", packet.RequestId, "seq", packet.Seq, "total", packet.Total, "partial", packet.Partial)
	if err := cacheDisseminatePacket(e, &packet); err != nil {
		log.Warn("Error caching packet", "error", err)
		if errors.Is(err, errPacketAlreadyRcv) {
			return nil
		}
	}

	if err := e.Strategies[packet.StrategyCode].HandlePacket(
		packet.RequestId,
		packet.Hop,
		packet.OriginalSender,
		uint64(p.id),
		packet.MaxPeers,
		packet.Data,
		packet.Partial,
		packet.Seq,
		packet.Total,
	); err != nil {
		log.Error("Error handling packet: ", err)
	}
	if e.Peers[packet.OriginalSender] == nil {
		fmt.Println("ERROR ORIGINAL SENDER NOT FOUND", packet.OriginalSender)
		return nil
	}
	pktInfo, _ := e.State.ReceivedPacketsFor(packet.RequestId)
	// all packets related to this requestID has been received
	if !pktInfo.Partial {
		log.Info("complete packet received, sending report", "total chunks", len(pktInfo.SeqReceived))
		return p2p.Send(e.Peers[packet.OriginalSender], DisseminateReport, DisseminateReportPacket{
			RequestId: packet.RequestId,
			Sender:    uint64(p.id),
			Hop:       packet.Hop,
			Time:      now,
			Full:      true,
		}) // should we ask for ACK?
	} else if pktInfo.TotalReceived == (len(pktInfo.SeqReceived)+1)/2 {
		log.Info("half packet received, sending report", "total chunks", pktInfo.TotalReceived)
		return p2p.Send(e.Peers[packet.OriginalSender], DisseminateReport, DisseminateReportPacket{
			RequestId: packet.RequestId,
			Sender:    uint64(p.id),
			Hop:       packet.Hop,
			Time:      now,
			Full:      false,
		}) // should we ask for ACK?
	} else {
		//log.Info("partial packet received", "chunks", pktInfo.SeqNum, "current chunk", packet.Seq)
	}
	return nil
}

func handleDisseminateReport(e *Engine, p *Peer, data io.Reader) error {
	var packet DisseminateReportPacket
	if err := rlp.Decode(data, &packet); err != nil {
		return err
	}
	channel, ok := e.State.ReceivedReports[packet.RequestId]
	if !ok {
		log.Error("Dissemination report id not found!")
		return nil // or error maybe
	}
	channel <- &strats.IndividualDisseminateResult{
		Sender:        p.id,
		Relay:         int(packet.Sender),
		Hop:           int(packet.Hop),
		ReceptionTime: time.Unix(int64(packet.Time)/int64(time.Second), int64(packet.Time)%int64(time.Second)),
		Full:          packet.Full,
	}
	return nil
}

func checkPeer(id int, peers []*Peer) (*Peer, error) {
	if id >= len(peers) {
		return nil, fmt.Errorf("peer %d not found", id)
	}

	peer := peers[id]
	if peer == nil || !peer.Connected {
		return nil, errPeerNotConnected
	}
	return peer, nil
}
