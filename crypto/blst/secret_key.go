package blst

import (
	"crypto/rand"

	"crypto/subtle"
	"encoding/hex"
	"fmt"
	blst "github.com/supranational/blst/bindings/go"
)

// bls12SecretKey used in the BLS signature scheme.
type bls12SecretKey struct {
	p *blst.SecretKey
}

func SecretKeyFromHex(key string) (SecretKey, error) {
	if len(key) != BLSSecretKeyLength*2 {
		return nil, ErrSecretHex
	}

	b, err := hex.DecodeString(key)
	if err != nil {
		return nil, err
	}

	sk, err := SecretKeyFromBytes(b)
	if err != nil {
		return nil, err
	}

	return sk, nil
}

// RandKey creates a new private key using a random method provided as an io.Reader.
func RandKey() (SecretKey, error) {
	// Generate 32 bytes of randomness
	var ikm [32]byte
	_, err := rand.Read(ikm[:])
	if err != nil {
		return nil, err
	}
	// Defensive check, that we have not generated a secret key,
	secKey := &bls12SecretKey{blst.KeyGen(ikm[:])}

	if IsZero(secKey.Marshal()) {
		return nil, ErrZeroKey
	}

	return secKey, nil
}

// SecretKeyFromBytes creates a BLS private key from a BigEndian byte slice.
func SecretKeyFromBytes(privKey []byte) (SecretKey, error) {
	if len(privKey) != BLSSecretKeyLength {
		return nil, fmt.Errorf("secret key must be %d bytes", BLSSecretKeyLength)
	}
	secKey := new(blst.SecretKey).Deserialize(privKey)
	if secKey == nil {
		return nil, ErrSecretUnmarshal
	}
	wrappedKey := &bls12SecretKey{p: secKey}

	if IsZero(privKey) {
		return nil, ErrZeroKey
	}
	return wrappedKey, nil
}

// BlsPublicKey obtains the public key corresponding to the BLS secret key.
func (s *bls12SecretKey) PublicKey() PublicKey {
	return &BlsPublicKey{p: new(blstPublicKey).From(s.p)}
}

// IsZero checks if the secret key is a zero key. We don't rely on the CGO to refer to the type of C.blst_scalar which
// is implemented in C to initialize the memory bits of C.blst_scalar to be zero. It is better for go binder to
// check if all the bytes of the secret key are zero.
func IsZero(sKey []byte) bool {
	b := byte(0)
	for _, s := range sKey {
		b |= s
	}
	return subtle.ConstantTimeByteEq(b, 0) == 1
}

// POPProof generate a proof of possession using a secret key - in a Validator client.
func (s *bls12SecretKey) POPProof(msg []byte) Signature {
	signature := new(blstSignature).Sign(s.p, msg, popDST)
	return &BlsSignature{s: signature}
}

// Sign a message using a secret key - in a Validator client.
//
// In IETF draft BLS specification:
// Sign(SK, message) -> signature: a signing algorithm that generates
//
//	a deterministic signature given a secret key SK and a message.
func (s *bls12SecretKey) Sign(msg []byte) Signature {
	signature := new(blstSignature).Sign(s.p, msg, generalDST)
	return &BlsSignature{s: signature}
}

// Marshal a secret key into a LittleEndian byte slice.
func (s *bls12SecretKey) Marshal() []byte {
	keyBytes := s.p.Serialize()
	return keyBytes
}

func (s *bls12SecretKey) Hex() string {
	return HexPrefix + hex.EncodeToString(s.Marshal())
}
