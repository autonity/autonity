package main

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"sync"
	"time"

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
	requests  map[uint64]chan any
	connected bool
	sync.RWMutex
}

type request struct {
	code       uint64
	packet     any
	responseCh <-chan any
}

type response struct {
	code   uint64
	packet any
}

func (p *Peer) Send(code uint64, data any) error {
	// to remove
	return p2p.Send(p, code, data)
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
	NoDelay    bool
}

func (p *Peer) sendUpdateTcpSocket(bufferSize int, noDelay bool) error {
	id := rand.Uint64()
	log.Info("sending update")
	_, err := p.dispatchRequest(id, UpdateTCPSocket, TCPOptionsPacket{uint64(bufferSize), noDelay}) // this will leak for now
	return err
}

func handleUpdateTcpSocket(_ *Engine, p *Peer, msg io.Reader) error {
	log.Info("received handle update")
	var opts TCPOptionsPacket
	if err := rlp.Decode(msg, &opts); err != nil {

		return err
	}
	p.UpdateSocketOptions(int(opts.BufferSize), opts.NoDelay)
	return nil
}

type DisseminateReportPacket struct {
	RequestId uint64
	Sender    uint64
	Hop       uint8
	Time      uint64
}

func handleDisseminatePacket(e *Engine, p *Peer, data io.Reader) error {
	var packet DisseminatePacket
	if err := rlp.Decode(data, &packet); err != nil {
		return err
	}
	now := uint64(time.Now().UnixNano()) // <-- We could add a timestamp before decoding too ?
	fmt.Println("[DisseminatePacket] << ", packet.RequestId, "FROM:", p.ID(), "ORIGIN", packet.OriginalSender, "HOP", packet.Hop)
	// check if first time received.
	if _, ok := e.receivedPackets[packet.RequestId]; ok {
		// do nothing
		return nil
	}
	e.receivedPackets[packet.RequestId] = struct{}{}
	if packet.Hop == 1 {
		// need to disseminate in the group
		group := disseminationGroup(e.id, e.peers)
		for i := range group {
			if group[i] != nil {
				group[i].sendDisseminate(packet.RequestId, packet.Data, packet.OriginalSender, 0)
			}
		}
	}
	if packet.Hop == 0 {
		// todo: include random peer selection logic - maybe set it as a parameter?
	}
	if e.peers[packet.OriginalSender] == nil {
		fmt.Println("ERROR ORIGINAL SENDER NOT FOUND", packet.OriginalSender)
		return nil
	}
	p2p.Send(e.peers[packet.OriginalSender], DisseminateReport, DisseminateReportPacket{
		RequestId: packet.RequestId,
		Sender:    uint64(e.peerToId(p)),
		Hop:       packet.Hop,
		Time:      now,
	}) // should we ask for ACK?
	return nil
}

func handleDisseminateReport(e *Engine, p *Peer, data io.Reader) error {
	var packet DisseminateReportPacket
	if err := rlp.Decode(data, &packet); err != nil {
		return err
	}
	channel, ok := e.receivedReports[packet.RequestId]
	if !ok {
		log.Error("Dissemination report id not found!")
		return nil // or error maybe
	}
	channel <- &IndividualDisseminateResult{
		Sender:        e.peerToId(p),
		Relay:         int(packet.Sender),
		Hop:           int(packet.Hop),
		ReceptionTime: time.Unix(int64(packet.Time)/int64(time.Second), int64(packet.Time)%int64(time.Second)),
	}
	return nil
}
