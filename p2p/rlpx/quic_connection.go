package rlpx

import (
	"context"
	"errors"
	"net"
	"runtime/debug"
	"sync"
	"time"

	quic2 "github.com/quic-go/quic-go"

	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/p2p/netutil"
)

type QuicConn struct {
	mu              sync.RWMutex
	session         quic2.Connection
	streams         []quic2.Stream
	index           int
	handlersStarted bool
	readCh          chan message
	writeCh         chan message
	closed          chan struct{}
}

func NewQuicConn(session quic2.Connection) *QuicConn {
	return &QuicConn{
		session: session,
		readCh:  make(chan message, 400),
		writeCh: make(chan message, 400),
		closed:  make(chan struct{}),
	}

}

func (qc *QuicConn) OpenStream() (quic2.Stream, error) {
	return qc.session.OpenStream()
}

func (qc *QuicConn) AcceptStream(ctx context.Context) (quic2.Stream, error) {
	return qc.session.AcceptStream(ctx)
}

func (qc *QuicConn) AddStream(stream quic2.Stream) {
	log.Info("Adding new stream", "conn", qc.session.RemoteAddr().String(), "stream", stream.StreamID(), "Initiated By", stream.StreamID().InitiatedBy().String())
	qc.mu.Lock()
	qc.streams = append(qc.streams, stream)
	if qc.handlersStarted { // streams which are added later
		go qc.reader(stream, qc.readCh)
		go qc.writer(stream, qc.writeCh)
	}
	qc.mu.Unlock()
}

func (qc *QuicConn) writer(stream quic2.Stream, writeCh <-chan message) {
	for {
		select {
		case <-qc.closed:
			log.Debug("connection closed, exiting write")
			return
		case msg := <-writeCh:

			//log.Info("new msg to write in stream", "id", stream.StreamID(), "len", len(msg.data))
			//var t time.Time
			//stream.SeWriteDeadline(t)
			now := time.Now()
			n, err := stream.Write(msg.data)
			if err != nil {
				log.Error("write error in stream", "id", stream.StreamID(), "err", err)
				return
			}
			log.Info("wrote data on stream", "id", stream.StreamID(), "num bytes", n, "time taken", time.Since(now).Microseconds())
		}
	}
}

func (qc *QuicConn) reader(stream quic2.Stream, readCh chan<- message) {
	log.Info("Reading stream", "id", stream.StreamID())
	//last := time.Now()
	r := ReadBuffer{}
	for {
		select {
		case <-qc.closed:
			log.Debug("connection closed, exiting read")
			return
		default:
			r.Reset()
			//var t time.Time
			//stream.SetReadDeadline(t)
			now := time.Now()
			//log.Info("header read - time since last read", "time taken", time.Since(last).Microseconds(), "stream ID", stream.StreamID())
			header, err := r.Read(stream, 16)
			//log.Info("length of data read in first attempt", "len", len(r.data), "capacity", cap(r.data), "stream ID", stream.StreamID(), "time taken", time.Since(now).Microseconds())
			if netutil.IsTemporaryError(err) {
				continue
			}
			if err != nil {
				msg := message{
					err: err,
				}
				readCh <- msg
				log.Error("read frame error", "err", err, "stream ID", stream.StreamID())
				return
			}

			fsize := readUint24(header[:16])
			// Frame size rounded up to 16 byte boundary for padding.
			rsize := fsize
			if padding := fsize % 16; padding > 0 {
				rsize += 16 - padding
			}
			//log.Info("Reading frame of size", "size", fsize, "rsize", rsize)
			snappyByte := header[3]

			// Read the frame content.
			now = time.Now()
			//log.Info("frame read - time since last read", "time taken", time.Since(last).Microseconds(), "stream ID", stream.StreamID())
			frame, err := r.Read(stream, int(rsize))
			if err != nil {
				log.Error("read frame content", "err", err, "stream ID", stream.StreamID())
				msg := message{
					err: err,
				}
				readCh <- msg
				return
			}
			//last = time.Now()
			msg := message{
				snappyByte: snappyByte,
				data:       make([]byte, fsize),
				err:        nil,
			}
			copy(msg.data, frame[:fsize])

			log.Info("frame read finished, new packet received", "id", stream.StreamID(), "size", rsize, "snappy byte", int(snappyByte), "time taken", time.Since(now).Microseconds())
			readCh <- msg
		}
	}

}

// TODO: Error channel
func (qc *QuicConn) handleStreams() {

	qc.mu.Lock()
	defer qc.mu.Unlock()
	for _, s := range qc.streams {
		stream := s
		go qc.reader(stream, qc.readCh)
		go qc.writer(stream, qc.writeCh)
	}
	qc.handlersStarted = true
}

func (qc *QuicConn) SelectStream() quic2.Stream {
	qc.mu.Lock()
	defer qc.mu.Unlock()
	stream := qc.streams[qc.index]
	qc.index = (qc.index + 1) % len(qc.streams)
	return stream
}

// TODO: stream selection/ multi stream handling for quic connection
func (qc *QuicConn) Read(b []byte) (int, error) {
	log.Trace("reading bytes on stream 0")
	qc.mu.RLock()
	s := qc.streams[0]
	defer qc.mu.RUnlock()
	n, err := s.Read(b)
	return n, err
}

func (qc *QuicConn) Write(b []byte) (int, error) {
	log.Trace("writing bytes on stream 0", "data", len(b))
	qc.mu.RLock()
	s := qc.streams[0]
	defer qc.mu.RUnlock()
	return s.Write(b)
}

func (qc *QuicConn) Close() error {
	log.Info("closing connection")
	debug.PrintStack()
	qc.mu.RLock()
	defer qc.mu.RUnlock()
	for _, stream := range qc.streams {
		log.Info("closing stream", "id", stream.StreamID())
		if err := stream.Close(); err != nil {
			return err
		}
	}
	close(qc.closed)
	return nil
}

func (qc *QuicConn) SetReadDeadline(t time.Time) error {
	qc.mu.RLock()
	defer qc.mu.RUnlock()
	for _, stream := range qc.streams {
		if err := stream.SetReadDeadline(time.Time{}); err != nil {
			return err
		}
	}
	return nil
}

func (qc *QuicConn) SetWriteDeadline(t time.Time) error {
	qc.mu.RLock()
	defer qc.mu.RUnlock()
	for _, stream := range qc.streams {
		if err := stream.SetWriteDeadline(time.Time{}); err != nil {
			return err
		}
	}
	return nil
}

func (qc *QuicConn) LocalAddr() net.Addr {
	return qc.session.LocalAddr()
}

func (qc *QuicConn) RemoteAddr() net.Addr {
	return qc.session.RemoteAddr()
}

func (qc *QuicConn) SetDeadline(t time.Time) error {
	return errors.New("unimplemented")
}
