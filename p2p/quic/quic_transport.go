package quic

import (
	"errors"
	"net"
	"sync"
	"time"

	quic2 "github.com/quic-go/quic-go"

	"github.com/autonity/autonity/log"
)

type QuicStream struct {
	stream quic2.Stream
}

func (s *QuicStream) Read(b []byte) (int, error) {
	n, err := s.Read(b)
	return n, err
}

func (s *QuicStream) Write(b []byte) (int, error) {
	//log.Trace("bytes wrote", "data", b)
	return s.Write(b)
}

func (s *QuicStream) Close() error {
	return s.Close()
}

func (s *QuicStream) SetReadDeadline(t time.Time) error {
	return s.SetReadDeadline(t)
}

func (s *QuicStream) SetWriteDeadline(t time.Time) error {
	return s.SetWriteDeadline(t)
}

type Conn struct {
	mu      sync.RWMutex
	Session quic2.Connection
	streams []*QuicStream
	index   int
}

func (qc *Conn) AddStream(stream quic2.Stream) {
	log.Info("Adding new stream", "conn", qc.Session.RemoteAddr().String(), "stream", stream.StreamID().StreamNum(), "Initiated By", stream.StreamID().InitiatedBy().String())
	qc.mu.Lock()
	qc.streams = append(qc.streams, &QuicStream{stream})
	qc.mu.Unlock()
}

func (qc *Conn) SelectStream() *QuicStream {
	qc.mu.Lock()
	defer qc.mu.Unlock()
	stream := qc.streams[qc.index]
	qc.index = (qc.index + 1) % len(qc.streams)
	return stream
}

// TODO: stream selection/ multi stream handling for quic connection
func (qc *Conn) Read(b []byte) (int, error) {
	n, err := qc.SelectStream().Read(b)
	return n, err
}

func (qc *Conn) Write(b []byte) (int, error) {
	//log.Trace("bytes wrote", "data", b)
	return qc.SelectStream().Write(b)
}

func (qc *Conn) Close() error {
	return qc.streams[0].Close()
}

func (qc *Conn) LocalAddr() net.Addr {
	return qc.Session.LocalAddr()
}

func (qc *Conn) RemoteAddr() net.Addr {
	return qc.Session.RemoteAddr()
}

func (qc *Conn) SetDeadline(t time.Time) error {
	return errors.New("unimplemented")
}

func (qc *Conn) SetReadDeadline(t time.Time) error {
	return qc.streams[0].SetReadDeadline(t)
}

func (qc *Conn) SetWriteDeadline(t time.Time) error {
	return qc.streams[0].SetWriteDeadline(t)
}
