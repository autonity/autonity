package main

import (
	"errors"
	"io"
	"math/rand"
	"time"

	"github.com/autonity/autonity/p2p"
	"github.com/autonity/autonity/rlp"
)

const (
	PingMsg = uint64(iota)
	PongMsg
	DataMsg
	AckDataMsg
)

var protocolHandlers = map[uint64]func(p *peer, data io.Reader) error{
	PingMsg: handlePing,
	PongMsg: handlePong,
	//DataMsg:    handleData, // random string of bytes
	//AckDataMsg: handleAckData,
	//BlockMsg //serialized block
	//AckBlockMsg
}

var (
	errUnknownRequest = errors.New("no matching request id")
)

type peer struct {
	*p2p.Peer
	p2p.MsgReadWriter
	// ICMP stats
	// Delay on TIME
	address   string
	requests  map[uint64]chan any
	connected bool
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

func (p *peer) reply(code uint64, data any) error {
	return p2p.Send(p, code, data)
}

func (p *peer) dispatchResponse(requestId uint64, packet any) error {
	req, ok := p.requests[requestId]
	if !ok {
		return errUnknownRequest
	}
	req <- packet
	return nil
}

func (p *peer) dispatchRequest(requestId uint64, code uint64, packet any) (chan any, error) {
	responseCh := make(chan any)
	p.requests[requestId] = responseCh

	return responseCh, p2p.Send(p, code, packet)
}

func (p *peer) sendPing() (uint64, error) {
	id := rand.Uint64()
	req, err := p.dispatchRequest(id, PingMsg, PingPacket{id})
	if err != nil {
		return 0, err
	}
	// we should check for timeout here
	return (<-req).(PongPacket).Time, nil
}

type PingPacket struct {
	RequestId uint64
}

type PongPacket struct {
	RequestId uint64
	Time      uint64
}

func handlePing(p *peer, data io.Reader) error {
	now := uint64(time.Now().UnixNano())
	var ping PingPacket
	if err := rlp.Decode(data, &ping); err != nil {
		return err
	}
	return p.reply(PongMsg, PongPacket{ping.RequestId, now})
}

func handlePong(p *peer, msg io.Reader) error {
	var pong PongPacket
	if err := rlp.Decode(msg, &pong); err != nil {
		return err
	}
	return p.dispatchResponse(pong.RequestId, pong)
}
