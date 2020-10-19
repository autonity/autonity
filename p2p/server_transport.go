package p2p

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"github.com/clearmatics/autonity/p2p/rlpx"
	"net"

	"github.com/clearmatics/autonity/crypto"
)

type testTransport struct {
	*rlpxTransport
	rpub     *ecdsa.PublicKey
	closeErr error
}

func newTestTransport(rpub *ecdsa.PublicKey, fd net.Conn, dialDest *ecdsa.PublicKey) transport {
	wrapped := newRLPX(fd, dialDest).(*rlpxTransport)
	wrapped.conn.InitWithSecrets(rlpx.Secrets{
		AES:        make([]byte, 16),
		MAC:        make([]byte, 16),
		EgressMAC:  sha256.New(),
		IngressMAC: sha256.New(),
	})
	return &testTransport{rpub: rpub, rlpxTransport: wrapped}
}

func (c *testTransport) doEncHandshake(prv *ecdsa.PrivateKey) (*ecdsa.PublicKey, error) {
	return c.rpub, nil
}

func (c *testTransport) doProtoHandshake(our *protoHandshake) (*protoHandshake, error) {
	pubkey := crypto.FromECDSAPub(c.rpub)[1:]
	return &protoHandshake{ID: pubkey, Name: "test"}, nil
}

func (c *testTransport) close(err error) {
	c.conn.Close()
	c.closeErr = err
}
