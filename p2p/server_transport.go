package p2p

import (
	"crypto/ecdsa"
	"net"

	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/crypto/sha3"
)

type testTransport struct {
	rpub *ecdsa.PublicKey
	*rlpx

	closeErr error
}

func NewTestTransport(rpub *ecdsa.PublicKey, fd net.Conn) transport {
	wrapped := newRLPX(fd).(*rlpx)
	wrapped.rw = newRLPXFrameRW(fd, secrets{
		MAC:        zero16,
		AES:        zero16,
		IngressMAC: sha3.NewKeccak256(),
		EgressMAC:  sha3.NewKeccak256(),
	})
	return &testTransport{rpub: rpub, rlpx: wrapped}
}

func (c *testTransport) doEncHandshake(prv *ecdsa.PrivateKey, dialDest *ecdsa.PublicKey) (*ecdsa.PublicKey, error) {
	return c.rpub, nil
}

func (c *testTransport) doProtoHandshake(our *protoHandshake) (*protoHandshake, error) {
	pubkey := crypto.FromECDSAPub(c.rpub)[1:]
	return &protoHandshake{ID: pubkey, Name: "test"}, nil
}

func (c *testTransport) close(err error) {
	c.rlpx.fd.Close()
	c.closeErr = err
}
