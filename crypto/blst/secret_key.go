package blst

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	common2 "github.com/autonity/autonity/common"
	blst "github.com/supranational/blst/bindings/go"
)

var zeroSecretKey = new(blst.SecretKey)

// bls12SecretKey used in the BLS signature scheme.
type bls12SecretKey struct {
	p *blst.SecretKey
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
	if secKey.IsZero() {
		return nil, ErrZeroKey
	}
	return secKey, nil
}

// SecretKeyFromECDSAKey creates a deterministic BLS private key from an ecdsa secret source.
// todo: (Jason) remove this function on top of node key management since currently there is no place to save BLS key, thus we still derive it from ECDSA.
func SecretKeyFromECDSAKey(ecdsaKey []byte) (SecretKey, error) {
	ecdsaKey = common2.LeftPadBytes(ecdsaKey, 32)

	blsSK := blst.KeyGen(ecdsaKey)
	if blsSK == nil {
		return nil, ErrSecretConvert
	}

	wrappedKey := &bls12SecretKey{p: blsSK}
	if wrappedKey.IsZero() {
		return nil, ErrZeroKey
	}
	return wrappedKey, nil
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
	if wrappedKey.IsZero() {
		return nil, ErrZeroKey
	}
	return wrappedKey, nil
}

// BlsPublicKey obtains the public key corresponding to the BLS secret key.
func (s *bls12SecretKey) PublicKey() PublicKey {
	return &BlsPublicKey{p: new(blstPublicKey).From(s.p)}
}

// IsZero checks if the secret key is a zero key.
func (s *bls12SecretKey) IsZero() bool {
	return s.p.Equals(zeroSecretKey)
}

// Sign a message using a secret key - in a Validator client.
//
// In IETF draft BLS specification:
// Sign(SK, message) -> signature: a signing algorithm that generates
//
//	a deterministic signature given a secret key SK and a message.
func (s *bls12SecretKey) Sign(msg []byte) Signature {
	signature := new(blstSignature).Sign(s.p, msg, dst)
	return &BlsSignature{s: signature}
}

// Marshal a secret key into a LittleEndian byte slice.
func (s *bls12SecretKey) Marshal() []byte {
	keyBytes := s.p.Serialize()
	if len(keyBytes) < BLSSecretKeyLength {
		emptyBytes := make([]byte, BLSSecretKeyLength-len(keyBytes))
		keyBytes = append(emptyBytes, keyBytes...)
	}
	return keyBytes
}

func (s *bls12SecretKey) Hex() string {
	return HexPrefix + hex.EncodeToString(s.Marshal())
}
