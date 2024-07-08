package main

import (
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/autonity/autonity/cmd/netdiag/strats"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/p2p"
	"github.com/autonity/autonity/rlp"
)

const (
	PingMsg = uint64(iota)
	PongMsg
	DataMsg
	AckDataMsg
	UpdateTCPSocket
	DisseminateRequest
	DisseminateReport
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
	//BlockMsg //serialized block
	//AckBlockMsg
}

var (
	errUnknownRequest = errors.New("no matching request id")
	errTimeout        = errors.New("request timed out")
)

type Peer struct {
	*p2p.Peer
	p2p.MsgReadWriter
	// ICMP stats
	// Delay on TIME
	ip        string
	rtt       time.Duration
	requests  map[uint64]chan any
	connected bool
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
	maxSize := 50_000
	if p.IsUDP() && len(data) > maxSize {
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

	for _, chunk := range chunks {
		packet := chunk
		fmt.Println("[DISSEMINATE] >>", p.ip, "|", "maxPeers", maxPeers, "originalSender", originalSender, "chunkID", chunk.Seq)
		go func() {
			err := p2p.Send(p.MsgReadWriter, DisseminateRequest, packet)
			if err != nil {
				log.Error("Disseminate Request", "err", err)
			}
		}()
	}

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
	return p.rtt
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

func (p *Peer) sendPing() (uint64, error) {
	id := rand.Uint64()
	req, err := p.dispatchRequest(id, PingMsg, PingPacket{id})
	if err != nil {
		return 0, err
	}
	fmt.Println("[PING] >>", id)
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
	fmt.Println("[PING] << , [PONG] >> ", ping.RequestId)
	return p2p.Send(p, PongMsg, PongPacket{ping.RequestId, now})
}

func handlePong(_ *Engine, p *Peer, msg io.Reader) error {
	var pong PongPacket
	if err := rlp.Decode(msg, &pong); err != nil {
		return err
	}
	fmt.Println("[PONG] <<", pong.RequestId)
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

func (p *Peer) sendData(data []byte) (uint64, time.Duration, error) {
	startTime := time.Now()
	id := rand.Uint64()
	fmt.Println("[DATAPACKET] >> ", id)
	req, err := p.dispatchRequest(id, DataMsg, DataPacket{id, data})
	if err != nil {
		log.Error("Unable to dispatch request", "error", err)
		return 0, 0, err
	}
	dispatchDuration := time.Now().Sub(startTime)
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
	fmt.Println("[DATAPACKET] >> ", id)
	return p.dispatchRequest(id, DataMsg, DataPacket{id, data})
}

func handleData(_ *Engine, p *Peer, data io.Reader) error {
	var dataPacket DataPacket
	if err := rlp.Decode(data, &dataPacket); err != nil {
		return err
	}
	now := uint64(time.Now().UnixNano()) // <-- We could add a timestamp before decoding too ?
	fmt.Println("[DATAPACKET] << ", dataPacket.RequestId)
	return p2p.Send(p, AckDataMsg, AckDataPacket{dataPacket.RequestId, now})
}

func handleAckData(_ *Engine, p *Peer, msg io.Reader) error {
	var ack AckDataPacket
	if err := rlp.Decode(msg, &ack); err != nil {
		return err
	}
	fmt.Println("[ACKDATA] << ", ack.RequestId)
	return p.dispatchResponse(ack.RequestId, ack)
}

// ***************************
// ***** UpdateTCPSocket **** *
// ***************************

type TCPOptionsPacket struct {
	BufferSize uint64
	Reset      bool
}

func (p *Peer) sendUpdateTcpSocket(bufferSize int, reset bool) error {
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
	for i, _ := range e.peers {
		peer, err := checkPeer(i, e.peers)
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

	chunkInfo, ok := e.state.ReceivedPackets[packet.RequestId]
	if ok {
		chunkInfo.SeqNum[packet.Seq] = math.MaxInt
		chunkInfo.Partial = false
		for _, num := range chunkInfo.SeqNum {
			if num != math.MaxInt { // none of the slots in seqNum array should be empty, for a complete packet reception
				chunkInfo.Partial = true
				break
			}
		}
		e.state.ReceivedPackets[packet.RequestId] = chunkInfo
	} else {
		chunkInfo := strats.ChunkInfo{Total: int(packet.Total),
			Partial: packet.Partial,
			SeqNum:  make([]int, max(packet.Total, 1)),
		}
		// use max.Int to mark the seq as received
		chunkInfo.SeqNum[packet.Seq] = math.MaxInt
		e.state.ReceivedPackets[packet.RequestId] = chunkInfo
	}
	return nil
}

func handleDisseminatePacket(e *Engine, p *Peer, data io.Reader) error {
	var packet DisseminatePacket
	if err := rlp.Decode(data, &packet); err != nil {
		return err
	}

	now := uint64(time.Now().UnixNano()) // <-- We could add a timestamp before decoding too ?
	fmt.Println("[DisseminatePacket] << ", packet.RequestId, "FROM:", p.ID(), "ORIGIN", packet.OriginalSender, "HOP", packet.Hop)
	// check if first time received.
	if pktInfo, ok := e.state.ReceivedPackets[packet.RequestId]; ok {
		// check if the seqNum is already received
		if len(pktInfo.SeqNum) <= int(packet.Seq) {
			log.Crit("invalid seqNum", "packet Info", pktInfo, "received packet", packet)
		}

		if !packet.Partial || (packet.Partial && pktInfo.SeqNum[packet.Seq] == math.MaxInt) {
			// do nothing
			return nil
		}
	}
	cacheDisseminatePacket(e, &packet)

	if err := e.strategies[packet.StrategyCode].HandlePacket(packet.RequestId, packet.Hop, packet.OriginalSender, packet.MaxPeers, packet.Data, packet.Partial, packet.Seq, packet.Total); err != nil {
		log.Error("Error handling packet: ", err)
	}
	if e.peers[packet.OriginalSender] == nil {
		fmt.Println("ERROR ORIGINAL SENDER NOT FOUND", packet.OriginalSender)
		return nil
	}
	pktInfo := e.state.ReceivedPackets[packet.RequestId]
	// all packets related to this requestID has been received
	if !pktInfo.Partial {
		log.Info("complete packet received, sending report", "total chunks", len(pktInfo.SeqNum))
		return p2p.Send(e.peers[packet.OriginalSender], DisseminateReport, DisseminateReportPacket{
			RequestId: packet.RequestId,
			Sender:    uint64(e.peerToId(p)),
			Hop:       packet.Hop,
			Time:      now,
		}) // should we ask for ACK?
	} else {
		log.Info("partial packet received", "chunks", pktInfo.SeqNum, "current chunk", packet.Seq)
	}
	return nil
}

func handleDisseminateReport(e *Engine, p *Peer, data io.Reader) error {
	var packet DisseminateReportPacket
	if err := rlp.Decode(data, &packet); err != nil {
		return err
	}
	channel, ok := e.state.ReceivedReports[packet.RequestId]
	if !ok {
		log.Error("Dissemination report id not found!")
		return nil // or error maybe
	}
	channel <- &strats.IndividualDisseminateResult{
		Sender:        e.peerToId(p),
		Relay:         int(packet.Sender),
		Hop:           int(packet.Hop),
		ReceptionTime: time.Unix(int64(packet.Time)/int64(time.Second), int64(packet.Time)%int64(time.Second)),
	}
	return nil
}
