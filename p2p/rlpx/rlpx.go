// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package rlpx implements the RLPx transport protocol.
package rlpx

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"hash"
	"io"
	mrand "math/rand"
	"net"
	"time"

	"github.com/golang/snappy"
	"golang.org/x/crypto/sha3"

	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/crypto/ecies"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/rlp"
)

type session interface {
	WriteFrame(conn io.Writer, code uint64, snappyByte byte, data []byte) error
	ReadFrame(conn io.Reader) ([]byte, byte, error)
	Init(rbuf ReadBuffer, wbuf WriteBuffer)
}

// Conn is an RLPx network connection. It wraps a low-level network connection. The
// underlying connection should not be used for other activity when it is wrapped by Conn.
//
// Before sending messages, a handshake must be performed by calling the Handshake method.
// This type is not generally safe for concurrent use, but reading and writing of messages
// may happen concurrently after the handshake.
type Conn struct {
	dialDest *ecdsa.PublicKey
	conn     net.Conn
	session  session

	// These are the buffers for snappy compression.
	// Compression is enabled if they are non-nil.
	snappyReadBuffer  []byte
	snappyWriteBuffer []byte
}

// NewConn wraps the given network connection. If dialDest is non-nil, the connection
// behaves as the initiator during the handshake.
func NewConn(conn net.Conn, dialDest *ecdsa.PublicKey) *Conn {
	return &Conn{
		dialDest: dialDest,
		conn:     conn,
	}
}

// SetSnappy enables or disables snappy compression of messages. This is usually called
// after the devp2p Hello message exchange when the negotiated version indicates that
// compression is available on both ends of the connection.
func (c *Conn) SetSnappy(snappy bool) {
	if snappy {
		c.snappyReadBuffer = []byte{}
		c.snappyWriteBuffer = []byte{}
	} else {
		c.snappyReadBuffer = nil
		c.snappyWriteBuffer = nil
	}
}

// SetReadDeadline sets the deadline for all future read operations.
func (c *Conn) SetReadDeadline(time time.Time) error {
	return c.conn.SetReadDeadline(time)
}

// SetWriteDeadline sets the deadline for all future write operations.
func (c *Conn) SetWriteDeadline(time time.Time) error {
	return c.conn.SetWriteDeadline(time)
}

// SetDeadline sets the deadline for all future read and write operations.
func (c *Conn) SetDeadline(time time.Time) error {
	return c.conn.SetDeadline(time)
}

// Read reads a message from the connection.
// The returned data buffer is valid until the next call to Read.
func (c *Conn) Read() (code uint64, data []byte, wireSize int, err error) {
	if c.session == nil {
		panic("can't ReadMsg before handshake")
	}
	frame, snappyByte, err := c.session.ReadFrame(c.conn)
	if err != nil {
		return 0, nil, 0, err
	}
	code, data, err = rlp.SplitUint64(frame)
	if err != nil {
		return 0, nil, 0, fmt.Errorf("invalid message code: %v", err)
	}
	wireSize = len(data)

	// If snappy is enabled, verify and decompress message.
	if c.snappyReadBuffer != nil && snappyByte == 0xFF {
		var actualSize int
		actualSize, err = snappy.DecodedLen(data)
		if err != nil {
			log.Error("Snappy decode length failed")
			return code, nil, 0, err
		}
		if actualSize > maxUint24 {
			return code, nil, 0, errPlainMessageTooLarge
		}
		c.snappyReadBuffer = growslice(c.snappyReadBuffer, actualSize)
		data, err = snappy.Decode(c.snappyReadBuffer, data)
		if err != nil {
			log.Error("Snappy decode failed")
		}
	}
	return code, data, wireSize, err
}

// Write writes a message to the connection.
//
// Write returns the written size of the message data. This may be less than or equal to
// len(data) depending on whether snappy compression is enabled.
func (c *Conn) Write(code uint64, data []byte) (uint32, error) {
	if c.session == nil {
		panic("can't WriteMsg before handshake")
	}
	if len(data) > maxUint24 {
		return 0, errPlainMessageTooLarge
	}
	var snappyByte byte = 0x00
	if c.snappyWriteBuffer != nil && len(data) > minSizeToCompress {
		// Ensure the buffer has sufficient size.
		// Package snappy will allocate its own buffer if the provided
		// one is smaller than MaxEncodedLen.
		c.snappyWriteBuffer = growslice(c.snappyWriteBuffer, snappy.MaxEncodedLen(len(data)))
		data = snappy.Encode(c.snappyWriteBuffer, data)
		snappyByte = 0xFF
	}

	wireSize := uint32(len(data))
	err := c.session.WriteFrame(c.conn, code, snappyByte, data)
	return wireSize, err
}

// Handshake performs the handshake. This must be called before any data is written
// or read from the connection.
func (c *Conn) Handshake(prv *ecdsa.PrivateKey) (*ecdsa.PublicKey, error) {
	var (
		sec Secrets
		err error
		h   handshakeState
	)
	if c.dialDest != nil {
		sec, err = h.runInitiator(c.conn, prv, c.dialDest)
	} else {
		sec, err = h.runRecipient(c.conn, prv)
	}
	if err != nil {
		log.Error("handshake failed", "error", err)
		return nil, err
	}
	_, ok := c.conn.(*QuicConn)
	if ok {
		c.session = NewQuicSession()
	} else {
		c.InitWithSecrets(sec)
	}
	c.session.Init(h.rbuf, h.wbuf)
	return sec.remote, err
}

// InitWithSecrets injects connection secrets as if a handshake had
// been performed. This cannot be called after the handshake.
func (c *Conn) InitWithSecrets(sec Secrets) {
	if c.session != nil {
		panic("can't handshake twice")
	}
	macc, err := aes.NewCipher(sec.MAC)
	if err != nil {
		panic("invalid MAC secret: " + err.Error())
	}
	encc, err := aes.NewCipher(sec.AES)
	if err != nil {
		panic("invalid AES secret: " + err.Error())
	}
	// we use an all-zeroes IV for AES because the key used
	// for encryption is ephemeral.
	iv := make([]byte, encc.BlockSize())

	c.session = &sessionState{
		enc:        cipher.NewCTR(encc, iv),
		dec:        cipher.NewCTR(encc, iv),
		egressMAC:  newHashMAC(macc, sec.EgressMAC),
		ingressMAC: newHashMAC(macc, sec.IngressMAC),
	}
}

// Close closes the underlying network connection.
func (c *Conn) Close() error {
	return c.conn.Close()
}

// Constants for the handshake.
const (
	sskLen = 16                     // ecies.MaxSharedKeyLength(pubKey) / 2
	sigLen = crypto.SignatureLength // elliptic S256
	pubLen = 64                     // 512 bit pubkey in uncompressed representation without format byte
	shaLen = 32                     // hash length (for nonce etc)

	eciesOverhead = 65 /* pubkey */ + 16 /* IV */ + 32 /* MAC */

	minSizeToCompress = 512 // for smaller packets the compression is ineffective
)

var (
	// this is used in place of actual frame header data.
	// TODO: replace this when Msg contains the protocol type code.
	zeroHeader = []byte{0xC2, 0x80, 0x80}

	// errPlainMessageTooLarge is returned if a decompressed message length exceeds
	// the allowed 24 bits (i.e. length >= 16MB).
	errPlainMessageTooLarge = errors.New("message length >= 16MB")
)

// Secrets represents the connection secrets which are negotiated during the handshake.
type Secrets struct {
	AES, MAC              []byte
	EgressMAC, IngressMAC hash.Hash
	remote                *ecdsa.PublicKey
}

// handshakeState contains the state of the encryption handshake.
type handshakeState struct {
	initiator            bool
	remote               *ecies.PublicKey  // remote-pubk
	initNonce, respNonce []byte            // nonce
	randomPrivKey        *ecies.PrivateKey // ecdhe-random
	remoteRandomPub      *ecies.PublicKey  // ecdhe-random-pubk

	rbuf ReadBuffer
	wbuf WriteBuffer
}

// RLPx v4 handshake auth (defined in EIP-8).
type authMsgV4 struct {
	Signature       [sigLen]byte
	InitiatorPubkey [pubLen]byte
	Nonce           [shaLen]byte
	Version         uint

	// Ignore additional fields (forward-compatibility)
	Rest []rlp.RawValue `rlp:"tail"`
}

// RLPx v4 handshake response (defined in EIP-8).
type authRespV4 struct {
	RandomPubkey [pubLen]byte
	Nonce        [shaLen]byte
	Version      uint

	// Ignore additional fields (forward-compatibility)
	Rest []rlp.RawValue `rlp:"tail"`
}

// runRecipient negotiates a session token on conn.
// it should be called on the listening side of the connection.
//
// prv is the local client's private key.
func (h *handshakeState) runRecipient(conn io.ReadWriter, prv *ecdsa.PrivateKey) (s Secrets, err error) {
	authMsg := new(authMsgV4)
	authPacket, err := h.readMsg(authMsg, prv, conn)
	if err != nil {
		log.Error("failed to read auth request packet", "err", err)
		return s, err
	}
	if err := h.handleAuthMsg(authMsg, prv); err != nil {
		return s, err
	}

	authRespMsg, err := h.makeAuthResp()
	if err != nil {
		return s, err
	}
	authRespPacket, err := h.sealEIP8(authRespMsg)
	if err != nil {
		return s, err
	}
	if _, err = conn.Write(authRespPacket); err != nil {
		log.Error("failed to write auth response packet", "err", err)
		return s, err
	}

	return h.secrets(authPacket, authRespPacket)
}

func (h *handshakeState) handleAuthMsg(msg *authMsgV4, prv *ecdsa.PrivateKey) error {
	// Import the remote identity.
	rpub, err := importPublicKey(msg.InitiatorPubkey[:])
	if err != nil {
		return err
	}
	h.initNonce = msg.Nonce[:]
	h.remote = rpub

	// Generate random keypair for ECDH.
	// If a private key is already set, use it instead of generating one (for testing).
	if h.randomPrivKey == nil {
		h.randomPrivKey, err = ecies.GenerateKey(rand.Reader, crypto.S256(), nil)
		if err != nil {
			return err
		}
	}

	// Check the signature.
	token, err := h.staticSharedSecret(prv)
	if err != nil {
		return err
	}
	signedMsg := xor(token, h.initNonce)
	remoteRandomPub, err := crypto.Ecrecover(signedMsg, msg.Signature[:])
	if err != nil {
		return err
	}
	h.remoteRandomPub, _ = importPublicKey(remoteRandomPub)
	return nil
}

// secrets is called after the handshake is completed.
// It extracts the connection secrets from the handshake values.
func (h *handshakeState) secrets(auth, authResp []byte) (Secrets, error) {
	ecdheSecret, err := h.randomPrivKey.GenerateShared(h.remoteRandomPub, sskLen, sskLen)
	if err != nil {
		return Secrets{}, err
	}

	// derive base secrets from ephemeral key agreement
	sharedSecret := crypto.Keccak256(ecdheSecret, crypto.Keccak256(h.respNonce, h.initNonce))
	aesSecret := crypto.Keccak256(ecdheSecret, sharedSecret)
	s := Secrets{
		remote: h.remote.ExportECDSA(),
		AES:    aesSecret,
		MAC:    crypto.Keccak256(ecdheSecret, aesSecret),
	}

	// setup sha3 instances for the MACs
	mac1 := sha3.NewLegacyKeccak256()
	mac1.Write(xor(s.MAC, h.respNonce))
	mac1.Write(auth)
	mac2 := sha3.NewLegacyKeccak256()
	mac2.Write(xor(s.MAC, h.initNonce))
	mac2.Write(authResp)
	if h.initiator {
		s.EgressMAC, s.IngressMAC = mac1, mac2
	} else {
		s.EgressMAC, s.IngressMAC = mac2, mac1
	}

	return s, nil
}

// staticSharedSecret returns the static shared secret, the result
// of key agreement between the local and remote static node key.
func (h *handshakeState) staticSharedSecret(prv *ecdsa.PrivateKey) ([]byte, error) {
	return ecies.ImportECDSA(prv).GenerateShared(h.remote, sskLen, sskLen)
}

// runInitiator negotiates a session token on conn.
// it should be called on the dialing side of the connection.
//
// prv is the local client's private key.
func (h *handshakeState) runInitiator(conn io.ReadWriter, prv *ecdsa.PrivateKey, remote *ecdsa.PublicKey) (s Secrets, err error) {
	h.initiator = true
	h.remote = ecies.ImportECDSAPublic(remote)

	authMsg, err := h.makeAuthMsg(prv)
	if err != nil {
		return s, err
	}
	authPacket, err := h.sealEIP8(authMsg)
	if err != nil {
		return s, err
	}

	if _, err = conn.Write(authPacket); err != nil {
		log.Error("failed to write auth packet", "err", err)
		return s, err
	}

	authRespMsg := new(authRespV4)
	authRespPacket, err := h.readMsg(authRespMsg, prv, conn)
	if err != nil {
		log.Error("failed to read auth response packet", "err", err)
		return s, err
	}
	if err := h.handleAuthResp(authRespMsg); err != nil {
		return s, err
	}

	return h.secrets(authPacket, authRespPacket)
}

// makeAuthMsg creates the initiator handshake message.
func (h *handshakeState) makeAuthMsg(prv *ecdsa.PrivateKey) (*authMsgV4, error) {
	// Generate random initiator nonce.
	h.initNonce = make([]byte, shaLen)
	_, err := rand.Read(h.initNonce)
	if err != nil {
		return nil, err
	}
	// Generate random keypair to for ECDH.
	h.randomPrivKey, err = ecies.GenerateKey(rand.Reader, crypto.S256(), nil)
	if err != nil {
		return nil, err
	}

	// Sign known message: static-shared-secret ^ nonce
	token, err := h.staticSharedSecret(prv)
	if err != nil {
		return nil, err
	}
	signed := xor(token, h.initNonce)
	signature, err := crypto.Sign(signed, h.randomPrivKey.ExportECDSA())
	if err != nil {
		return nil, err
	}

	msg := new(authMsgV4)
	copy(msg.Signature[:], signature)
	copy(msg.InitiatorPubkey[:], crypto.FromECDSAPub(&prv.PublicKey)[1:])
	copy(msg.Nonce[:], h.initNonce)
	msg.Version = 4
	return msg, nil
}

func (h *handshakeState) handleAuthResp(msg *authRespV4) (err error) {
	h.respNonce = msg.Nonce[:]
	h.remoteRandomPub, err = importPublicKey(msg.RandomPubkey[:])
	return err
}

func (h *handshakeState) makeAuthResp() (msg *authRespV4, err error) {
	// Generate random nonce.
	h.respNonce = make([]byte, shaLen)
	if _, err = rand.Read(h.respNonce); err != nil {
		return nil, err
	}

	msg = new(authRespV4)
	copy(msg.Nonce[:], h.respNonce)
	copy(msg.RandomPubkey[:], exportPubkey(&h.randomPrivKey.PublicKey))
	msg.Version = 4
	return msg, nil
}

// readMsg reads an encrypted handshake message, decoding it into msg.
func (h *handshakeState) readMsg(msg interface{}, prv *ecdsa.PrivateKey, r io.Reader) ([]byte, error) {
	h.rbuf.Reset()
	h.rbuf.grow(512)

	// Read the size prefix.
	prefix, err := h.rbuf.Read(r, 2)
	if err != nil {
		return nil, err
	}
	size := binary.BigEndian.Uint16(prefix)

	// Read the handshake packet.
	packet, err := h.rbuf.Read(r, int(size))
	if err != nil {
		return nil, err
	}
	dec, err := ecies.ImportECDSA(prv).Decrypt(packet, nil, prefix)
	if err != nil {
		return nil, err
	}
	// Can't use rlp.DecodeBytes here because it rejects
	// trailing data (forward-compatibility).
	s := rlp.NewStream(bytes.NewReader(dec), 0)
	err = s.Decode(msg)
	h.rbuf.Lock()
	data := h.rbuf.data[:len(prefix)+len(packet)]
	h.rbuf.Unlock()
	return data, err
}

// sealEIP8 encrypts a handshake message.
func (h *handshakeState) sealEIP8(msg interface{}) ([]byte, error) {
	h.wbuf.Reset()

	// Write the message plaintext.
	if err := rlp.Encode(&h.wbuf, msg); err != nil {
		return nil, err
	}
	// Pad with random amount of data. the amount needs to be at least 100 bytes to make
	// the message distinguishable from pre-EIP-8 handshakes.
	h.wbuf.AppendZero(mrand.Intn(100) + 100)

	prefix := make([]byte, 2)
	binary.BigEndian.PutUint16(prefix, uint16(len(h.wbuf.data)+eciesOverhead))

	enc, err := ecies.Encrypt(rand.Reader, h.remote, h.wbuf.data, nil, prefix)
	return append(prefix, enc...), err
}

// importPublicKey unmarshals 512 bit public keys.
func importPublicKey(pubKey []byte) (*ecies.PublicKey, error) {
	var pubKey65 []byte
	switch len(pubKey) {
	case 64:
		// add 'uncompressed key' flag
		pubKey65 = append([]byte{0x04}, pubKey...)
	case 65:
		pubKey65 = pubKey
	default:
		return nil, fmt.Errorf("invalid public key length %v (expect 64/65)", len(pubKey))
	}
	// TODO: fewer pointless conversions
	pub, err := crypto.UnmarshalPubkey(pubKey65)
	if err != nil {
		return nil, err
	}
	return ecies.ImportECDSAPublic(pub), nil
}

func exportPubkey(pub *ecies.PublicKey) []byte {
	if pub == nil {
		panic("nil pubkey")
	}
	return elliptic.Marshal(pub.Curve, pub.X, pub.Y)[1:]
}

func xor(one, other []byte) (xor []byte) {
	xor = make([]byte, len(one))
	for i := 0; i < len(one); i++ {
		xor[i] = one[i] ^ other[i]
	}
	return xor
}
