package rlpx

import (
	"crypto/cipher"
	"crypto/hmac"
	"errors"
	"fmt"
	"hash"
	"io"
	"sync"

	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/rlp"
)

// sessionState contains the session keys.
type sessionState struct {
	enc cipher.Stream
	dec cipher.Stream

	egressMAC  hashMAC
	ingressMAC hashMAC
	rbuf       ReadBuffer
	wbuf       WriteBuffer
	once       sync.Once
}

// hashMAC holds the state of the RLPx v4 MAC contraption.
type hashMAC struct {
	cipher     cipher.Block
	hash       hash.Hash
	aesBuffer  [16]byte
	hashBuffer [32]byte
	seedBuffer [32]byte
}

func newHashMAC(cipher cipher.Block, h hash.Hash) hashMAC {
	m := hashMAC{cipher: cipher, hash: h}
	if cipher.BlockSize() != len(m.aesBuffer) {
		panic(fmt.Errorf("invalid MAC cipher block size %d", cipher.BlockSize()))
	}
	if h.Size() != len(m.hashBuffer) {
		panic(fmt.Errorf("invalid MAC digest size %d", h.Size()))
	}
	return m
}

func isQuicWriter(conn io.Writer) bool {
	_, ok := conn.(*QuicConn)
	return ok
}

func isQuicReader(conn io.Reader) bool {
	_, ok := conn.(*QuicConn)
	return ok
}

func (h *sessionState) Init(rbuf ReadBuffer, wbuf WriteBuffer) {
	h.rbuf = rbuf
	h.wbuf = wbuf
}

func (h *sessionState) ReadFrame(conn io.Reader) ([]byte, byte, error) {
	h.rbuf.Reset()

	zeroByte := byte(0)
	// Read the frame header.
	headerSize := 32
	if isQuicReader(conn) {
		headerSize = 16
	}
	header, err := h.rbuf.Read(conn, headerSize)
	if err != nil {
		log.Error("read frame error", "err", err)
		return nil, zeroByte, err
	}

	if !isQuicReader(conn) {
		// Verify header MAC.
		wantHeaderMAC := h.ingressMAC.computeHeader(header[:16])
		if !hmac.Equal(wantHeaderMAC, header[16:]) {
			return nil, zeroByte, errors.New("bad header MAC")
		}
		// Decrypt the frame header to get the frame size.
		h.dec.XORKeyStream(header[:16], header[:16])
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
	frame, err := h.rbuf.Read(conn, int(rsize))
	if err != nil {
		log.Error("read frame content", "err", err)
		return nil, snappyByte, err
	}

	if !isQuicReader(conn) {
		// Validate frame MAC.
		frameMAC, err := h.rbuf.Read(conn, 16)
		if err != nil {
			return nil, snappyByte, err
		}
		wantFrameMAC := h.ingressMAC.computeFrame(frame)
		if !hmac.Equal(wantFrameMAC, frameMAC) {
			return nil, snappyByte, errors.New("bad frame MAC")
		}

		// Decrypt the frame data.
		h.dec.XORKeyStream(frame, frame)
	}
	return frame[:fsize], snappyByte, nil
}
func (h *sessionState) WriteFrame(conn io.Writer, code uint64, snappyByte byte, data []byte) error {
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
	if !isQuicWriter(conn) {

		h.enc.XORKeyStream(header, header)

		// Write header MAC.
		h.wbuf.Write(h.egressMAC.computeHeader(header))
	}

	// Encode and encrypt the frame data.
	offset := len(h.wbuf.data)
	h.wbuf.data = rlp.AppendUint64(h.wbuf.data, code)
	h.wbuf.Write(data)
	if padding := fsize % 16; padding > 0 {
		h.wbuf.AppendZero(16 - padding)
	}
	if !isQuicWriter(conn) {
		framedata := h.wbuf.data[offset:]
		h.enc.XORKeyStream(framedata, framedata)
		// Write frame MAC.
		h.wbuf.Write(h.egressMAC.computeFrame(framedata))
	}
	//now := time.Now()
	_, err := conn.Write(h.wbuf.data)
	//log.Info("conn write", "length", len(h.wbuf.data), "time taken", time.Since(now).Nanoseconds())
	return err
}

// computeHeader computes the MAC of a frame header.
func (m *hashMAC) computeHeader(header []byte) []byte {
	sum1 := m.hash.Sum(m.hashBuffer[:0])
	return m.compute(sum1, header)
}

// computeFrame computes the MAC of framedata.
func (m *hashMAC) computeFrame(framedata []byte) []byte {
	m.hash.Write(framedata)
	seed := m.hash.Sum(m.seedBuffer[:0])
	return m.compute(seed, seed[:16])
}

// compute computes the MAC of a 16-byte 'seed'.
//
// To do this, it encrypts the current value of the hash state, then XORs the ciphertext
// with seed.conn The obtained value is written back into the hash state and hash output is
// taken again. The first 16 bytes of the resulting sum are the MAC value.
//
// This MAC construction is a horrible, legacy thing.
func (m *hashMAC) compute(sum1, seed []byte) []byte {
	if len(seed) != len(m.aesBuffer) {
		panic("invalid MAC seed")
	}

	m.cipher.Encrypt(m.aesBuffer[:], sum1)
	for i := range m.aesBuffer {
		m.aesBuffer[i] ^= seed[i]
	}
	m.hash.Write(m.aesBuffer[:])
	sum2 := m.hash.Sum(m.hashBuffer[:0])
	return sum2[:16]
}
