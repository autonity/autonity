package p2p

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

type UdpTransport struct {
	readPacketCh chan Msg
	addr         *net.UDPAddr
	conn         *net.UDPConn
	maxSize      uint32
	wbuf         []byte
	wmu          sync.Mutex
}

func newUdpTransport(conn *net.UDPConn, addr *net.UDPAddr) *UdpTransport {
	return &UdpTransport{
		readPacketCh: make(chan Msg, 100),
		addr:         addr,
		conn:         conn,
		maxSize:      65000,
		wbuf:         make([]byte, 65001),
	}
}

func (u *UdpTransport) doEncHandshake(prv *ecdsa.PrivateKey) (*ecdsa.PublicKey, error) {
	return nil, nil
}

func (u *UdpTransport) doProtoHandshake(our *protoHandshake) (their *protoHandshake, err error) {
	// same logic as rlpx
	werr := make(chan error, 1)
	go func() { werr <- Send(u, handshakeMsg, our) }()
	if their, err = readProtocolHandshake(u); err != nil {
		<-werr
		return nil, err
	}
	if err := <-werr; err != nil {
		return nil, fmt.Errorf("write error: %v", err)
	}
	return their, nil
}

func (u *UdpTransport) ReadMsg() (Msg, error) {
	// we need to sign the messages
	packet := <-u.readPacketCh
	//log.Trace("low level reading packet")
	return packet, nil
}

func (u *UdpTransport) HandlePacket(packet []byte) error {
	if len(packet) < 1 {
		return errors.New("packet too short")
	}
	msg := Msg{
		Code:       uint64(packet[0]),
		Size:       uint32(len(packet) - 1),
		Payload:    bytes.NewReader(packet[1:]),
		ReceivedAt: time.Now(),
	}
	//fmt.Println("PACKET RECEIVED", "data", msg)
	//log.Trace("low level writing packet", "msg", msg)
	u.readPacketCh <- msg
	return nil
}

func (u *UdpTransport) WriteMsg(msg Msg) error {
	//log.Info("WRITING PACKET")
	u.wmu.Lock()
	defer u.wmu.Unlock()
	if msg.Size > u.maxSize {
		return errors.New("message too long")
	}
	// PACKET = [ MSG_CODE, DATA ] // MSG_CODE is 1 byte.
	u.wbuf[0] = byte(msg.Code)
	n, err := msg.Payload.Read(u.wbuf[1:])
	if n != int(msg.Size) {
		return errors.New("weird message size")
	}
	_, err = u.conn.WriteToUDP(u.wbuf[:1+n], u.addr)
	//log.Info("ON THE WIRE", "packet", u.wbuf[:1+n], "err", err)
	return err
}

func (u *UdpTransport) close(err error) {

}
