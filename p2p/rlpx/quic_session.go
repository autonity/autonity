package rlpx

import (
	"io"
	"sync"

	"github.com/autonity/autonity/rlp"
)

type message struct {
	snappyByte byte
	data       []byte
	err        error
}

// quicSession contains the session keys.
type quicSession struct {
	rbuf ReadBuffer
	wbuf WriteBuffer
	once sync.Once
}

func NewQuicSession() *quicSession {
	return &quicSession{}
}

func (h *quicSession) Init(rbuf ReadBuffer, wbuf WriteBuffer) {
	h.rbuf = rbuf
	h.wbuf = wbuf
}

func (h *quicSession) ReadFrame(conn io.Reader) ([]byte, byte, error) {

	h.once.Do(func() {
		go conn.(*QuicConn).handleStreams()
	})
	msg := <-conn.(*QuicConn).readCh
	//log.Info("Read a session message", "msg", msg, "len", len(msg.data))
	return msg.data, msg.snappyByte, msg.err
}

func (h *quicSession) WriteFrame(conn io.Writer, code uint64, snappyByte byte, data []byte) error {
	h.wbuf.Reset()
	// Write header.
	fsize := rlp.IntSize(code) + len(data)
	if fsize > maxUint24 {
		return errPlainMessageTooLarge
	}
	header := h.wbuf.AppendZero(16)
	putUint24(uint32(fsize), header)
	header[3] = snappyByte
	copy(header[4:], zeroHeader)

	// Encode and encrypt the frame data.
	h.wbuf.data = rlp.AppendUint64(h.wbuf.data, code)
	h.wbuf.Write(data)
	if padding := fsize % 16; padding > 0 {
		h.wbuf.AppendZero(16 - padding)
	}
	msg := message{
		data: make([]byte, len(h.wbuf.data)),
	}
	copy(msg.data, h.wbuf.data)
	//now := time.Now()
	conn.(*QuicConn).writeCh <- msg
	//log.Info("conn write initiated", "length", len(h.wbuf.Data), "time taken", time.Since(now).Nanoseconds())
	return nil
}
