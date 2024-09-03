package core

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"sync"
	"time"

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
	log.Info("Disseminate Packet")
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
		//go func() {
		err := p2p.Send(p.MsgReadWriter, DisseminateRequest, packet)
		if err != nil {
			log.Error("Disseminate Request", "err", err)
		} else {
			//log.Info("[DISSEMINATE]", " ip ", p.ip, "maxPeers", maxPeers, "originalSender", originalSender, "chunkID", packet.Seq)
		}
		//}()
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
	log.Debug("[DATAPACKET] >> ", "id", id)
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
	log.Debug("[DATAPACKET] >> ", "id", id)
	return p.dispatchRequest(id, DataMsg, DataPacket{id, data})
}

func handleData(_ *Engine, p *Peer, data io.Reader) error {
	var dataPacket DataPacket
	if err := rlp.Decode(data, &dataPacket); err != nil {
		return err
	}
	now := uint64(time.Now().UnixNano()) // <-- We could add a timestamp before decoding too ?
	log.Debug("[DATAPACKET] << ", "requestId", dataPacket.RequestId)
	return p2p.Send(p, AckDataMsg, AckDataPacket{dataPacket.RequestId, now})
}

func handleAckData(_ *Engine, p *Peer, msg io.Reader) error {
	var ack AckDataPacket
	if err := rlp.Decode(msg, &ack); err != nil {
		return err
	}
	log.Debug("[ACKDATA] << ", "requestId", ack.RequestId)
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
	e.State.LatencyMatrix[p.id] = make([]time.Duration, len(e.Peers))
	for i, l := range latencyPacket.LatencyArray {
		if i >= len(e.Peers) {
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
		id, TriggerRequest,
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
	go func(e *Engine, p *Peer, trigger *TriggerPacket) {
		if !e.State.PingReceived {
			// get latency for all peers
			latency := PingPeers(e)

			// set our own latency in the matrix
			e.State.LatencyMatrix[e.Id] = FilterAveRtt(latency)
			e.State.PingReceived = true
			for peerId, l := range e.State.LatencyMatrix[e.Id] {
				if e.Id != peerId && l == 0 {
					e.State.PingReceived = false
					break
				}
			}
		}

		if err := BroadcastLatency(e, trigger.Strategy, e.State.LatencyMatrix[e.Id]); err != nil {
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
}

func cacheDisseminatePacket(e *Engine, packet *DisseminatePacket) error {

	chunkInfo, ok := e.State.ReceivedPackets[packet.RequestId]
	if ok {
		chunkInfo.SeqReceived[packet.Seq] = true
		chunkInfo.TotalReceived++
		if chunkInfo.TotalReceived == len(chunkInfo.SeqReceived) {
			chunkInfo.Partial = false
		}
		e.State.ReceivedPackets[packet.RequestId] = chunkInfo
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
		e.State.ReceivedPackets[packet.RequestId] = chunkInfo
	}
	return nil
}

func handleDisseminatePacket(e *Engine, p *Peer, data io.Reader) error {
	var packet DisseminatePacket
	if err := rlp.Decode(data, &packet); err != nil {
		return err
	}

	now := uint64(time.Now().UnixNano()) // <-- We could add a timestamp before decoding too ?
	//fmt.Println("[DisseminatePacket] << ", packet.RequestId, "FROM:", p.ID(), "ORIGIN", packet.OriginalSender, "HOP", packet.Hop)
	// check if first time received.
	if pktInfo, ok := e.State.ReceivedPackets[packet.RequestId]; ok {
		// check if the seqNum is already received
		if len(pktInfo.SeqReceived) <= int(packet.Seq) {
			log.Crit("invalid seqNum", "packet Info", pktInfo, "received packet", packet)
		}

		if !packet.Partial || (packet.Partial && pktInfo.SeqReceived[packet.Seq]) {
			//log.Info("packet has already arrived", "packet", packet.Seq)
			// do nothing
			return nil
		}
	}
	cacheDisseminatePacket(e, &packet)

	if err := e.Strategies[packet.StrategyCode].HandlePacket(packet.RequestId, packet.Hop, packet.OriginalSender, packet.MaxPeers, packet.Data, packet.Partial, packet.Seq, packet.Total); err != nil {
		log.Error("Error handling packet: ", err)
	}
	if e.Peers[packet.OriginalSender] == nil {
		fmt.Println("ERROR ORIGINAL SENDER NOT FOUND", packet.OriginalSender)
		return nil
	}
	pktInfo := e.State.ReceivedPackets[packet.RequestId]
	// all packets related to this requestID has been received
	if !pktInfo.Partial {
		log.Info("complete packet received, sending report", "total chunks", len(pktInfo.SeqReceived))
		return p2p.Send(e.Peers[packet.OriginalSender], DisseminateReport, DisseminateReportPacket{
			RequestId: packet.RequestId,
			Sender:    uint64(p.id),
			Hop:       packet.Hop,
			Time:      now,
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
