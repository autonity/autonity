// Package common provides the BLS interfaces that are implemented by the various BLS wrappers.
package common

// BLSSecretKey represents a BLS secret or private key.
type BLSSecretKey interface {
	PublicKey() BLSPublicKey
	Sign(msg []byte) BLSSignature
	Marshal() []byte
	IsZero() bool
	Hex() string
}

// BLSPublicKey represents a BLS public key.
type BLSPublicKey interface {
	Marshal() []byte
	Copy() BLSPublicKey
	Aggregate(p2 BLSPublicKey) (BLSPublicKey, error)
	IsInfinite() bool
	Hex() string
}

// BLSSignature represents a BLS signature.
type BLSSignature interface {
	Verify(pubKey BLSPublicKey, msg []byte) bool
	AggregateVerify(pubKeys []BLSPublicKey, msgs [][32]byte) bool
	FastAggregateVerify(pubKeys []BLSPublicKey, msg [32]byte) bool
	Marshal() []byte
	Copy() BLSSignature
	Hex() string
}
