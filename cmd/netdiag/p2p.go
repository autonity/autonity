package main

import (
	"errors"
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
	ProtocolMessages
)

var protocolHandlers = map[uint64]func(p *Peer, data io.Reader) error{
	PingMsg:         handlePing,
	PongMsg:         handlePong,
	DataMsg:         handleData, // random string of bytes
	AckDataMsg:      handleAckData,
	UpdateTCPSocket: handleUpdateTcpSocket,
	//BlockMsg //serialized block
	//AckBlockMsg
}

var (
	errUnknownRequest = errors.New("no matching request id")
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

func (p *Peer) reply(code uint64, data any) error {
	return p2p.Send(p, code, data)
}

func (p *Peer) dispatchResponse(requestId uint64, packet any) error {
	p.RLock()
	req, ok := p.requests[requestId]
	p.RUnlock()
	if !ok {
		return errUnknownRequest
	}
	req <- packet
	return nil
}

func (p *Peer) dispatchRequest(requestId uint64, code uint64, packet any) (chan any, error) {
	responseCh := make(chan any)
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
	// we should check for timeout here
	return (<-req).(PongPacket).Time, nil
}

func handlePing(p *Peer, data io.Reader) error {
	now := uint64(time.Now().UnixNano())
	var ping PingPacket
	if err := rlp.Decode(data, &ping); err != nil {
		return err
	}
	return p.reply(PongMsg, PongPacket{ping.RequestId, now})
}

func handlePong(p *Peer, msg io.Reader) error {
	var pong PongPacket
	if err := rlp.Decode(msg, &pong); err != nil {
		return err
	}
	return p.dispatchResponse(pong.RequestId, pong)
}

// **** SENDDATA *****
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
	req, err := p.dispatchRequest(id, DataMsg, DataPacket{id, data})
	if err != nil {
		return 0, 0, err
	}
	dispatchDuration := time.Now().Sub(startTime)
	ackTime := (<-req).(AckDataPacket).Time
	return ackTime, dispatchDuration, nil
}

func handleData(p *Peer, data io.Reader) error {
	var dataPacket DataPacket
	if err := rlp.Decode(data, &dataPacket); err != nil {
		return err
	}
	now := uint64(time.Now().UnixNano()) // <-- We could add a timestamp before decoding too ?
	return p.reply(AckDataMsg, AckDataPacket{dataPacket.RequestId, now})
}

func handleAckData(p *Peer, msg io.Reader) error {
	var ack AckDataPacket
	if err := rlp.Decode(msg, &ack); err != nil {
		return err
	}
	return p.dispatchResponse(ack.RequestId, ack)
}

// **** TCPSOCKETOPTIONS *****
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

func handleUpdateTcpSocket(p *Peer, msg io.Reader) error {
	log.Info("received handle update")
	var opts TCPOptionsPacket
	if err := rlp.Decode(msg, &opts); err != nil {

		return err
	}
	p.UpdateSocketOptions(int(opts.BufferSize), opts.NoDelay)
	return nil
}
