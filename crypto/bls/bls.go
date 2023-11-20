package bls

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/autonity/autonity/crypto/bls/blst"
	"github.com/autonity/autonity/crypto/bls/common"
	"github.com/pkg/errors"
)

// SecretKeyFromECDSAKey generate a BLS private key by sourcing from an ECDSA private key.
func SecretKeyFromECDSAKey(sk *ecdsa.PrivateKey) (SecretKey, error) {
	return blst.SecretKeyFromECDSAKey(sk.D.Bytes())
}

// SecretKeyFromHex generate a BLS private key from hex string.
func SecretKeyFromHex(hexKey string) (SecretKey, error) {
	b, err := hex.DecodeString(hexKey)
	if byteErr, ok := err.(hex.InvalidByteError); ok {
		return nil, fmt.Errorf("invalid hex character %q in private key", byte(byteErr))
	} else if err != nil {
		return nil, errors.New("invalid hex data for private key")
	}
	return blst.SecretKeyFromBytes(b)
}

// SecretKeyFromBytes creates a BLS private key from a BigEndian byte slice.
func SecretKeyFromBytes(privKey []byte) (SecretKey, error) {
	return blst.SecretKeyFromBytes(privKey)
}

// SecretKeyFromBigNum takes in a big number string and creates a BLS private key.
func SecretKeyFromBigNum(s string) (SecretKey, error) {
	num := new(big.Int)
	num, ok := num.SetString(s, 10)
	if !ok {
		return nil, errors.New("could not set big int from string")
	}
	bts := num.Bytes()
	if len(bts) != 32 {
		return nil, errors.Errorf("provided big number string sets to a key unequal to 32 bytes: %d != 32", len(bts))
	}
	return SecretKeyFromBytes(bts)
}

// PublicKeyBytesFromString decode a BLS public key from a string.
func PublicKeyBytesFromString(pubKey string) ([]byte, error) {
	return blst.PublicKeyBytesFromString(pubKey)
}

// PublicKeyFromBytes creates a BLS public key from a  BigEndian byte slice.
func PublicKeyFromBytes(pubKey []byte) (PublicKey, error) {
	return blst.PublicKeyFromBytes(pubKey)
}

// SignatureFromBytes creates a BLS signature from a LittleEndian byte slice.
func SignatureFromBytes(sig []byte) (Signature, error) {
	return blst.SignatureFromBytes(sig)
}

// AggregatePublicKeys aggregates the provided raw public keys into a single key.
func AggregatePublicKeys(pubs [][]byte) (PublicKey, error) {
	return blst.AggregatePublicKeys(pubs)
}

// AggregateSignatures converts a list of signatures into a single, aggregated sig.
func AggregateSignatures(sigs []common.BLSSignature) common.BLSSignature {
	return blst.AggregateSignatures(sigs)
}

// VerifyMultipleSignatures verifies multiple signatures for distinct messages securely.
func VerifyMultipleSignatures(sigs [][]byte, msgs [][32]byte, pubKeys []common.BLSPublicKey) (bool, error) {
	return blst.VerifyMultipleSignatures(sigs, msgs, pubKeys)
}

// NewAggregateSignature creates a blank aggregate signature.
func NewAggregateSignature() common.BLSSignature {
	return blst.NewAggregateSignature()
}

// RandKey creates a new private key using a random input.
func RandKey() (common.BLSSecretKey, error) {
	return blst.RandKey()
}
