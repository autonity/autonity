package blst

import (
	"errors"
	"io"

	"github.com/autonity/autonity/rlp"
)

const (
	BLSSecretKeyLength       = 32
	BLSPubkeyLength          = 48
	BLSSignatureLength       = 96
	BLSPubKeyHexStringLength = 98
)

// ErrSecretHex describes an error on the wrong hex string of secrete key
var ErrSecretHex = errors.New("could not generate a blst secret key from hex")

// ErrZeroKey describes an error due to a zero secret key.
var ErrZeroKey = errors.New("generated secret key is zero")

// ErrSecretUnmarshal describes an error which happens during unmarshalling
// a secret key.
var ErrSecretUnmarshal = errors.New("could not unmarshal bytes into secret key")

// ErrInfinitePubKey describes an error due to an infinite public key.
var ErrInfinitePubKey = errors.New("received an infinite public key")

// ErrInvalidSecretKey describes an invalid length of bytes of secret blst key
var ErrInvalidSecretKey = errors.New("secret key must be 32 bytes")

// SecretKey represents a BLS secret or private key.
type SecretKey interface {
	PublicKey() PublicKey
	Sign(msg []byte) Signature
	POPProof(msg []byte) Signature
	Marshal() []byte
	Hex() string
}

// PublicKey represents a BLS public key.
type PublicKey interface {
	Marshal() []byte
	Copy() PublicKey
	Aggregate(p2 PublicKey) (PublicKey, error)
	IsInfinite() bool
	Hex() string
}

// Signature represents a BLS signature.
type Signature interface {
	Verify(pubKey PublicKey, msg []byte) bool
	POPVerify(pubKey PublicKey, msg []byte) bool
	IsZero() bool
	AggregateVerify(pubKeys []PublicKey, msgs [][32]byte) bool
	FastAggregateVerify(pubKeys []PublicKey, msg [32]byte) bool
	Marshal() []byte
	Copy() Signature
	Hex() string
	EncodeRLP(w io.Writer) error
	DecodeRLP(stream *rlp.Stream) error
	MarshalText() ([]byte, error)
	UnmarshalText(input []byte) error
}
