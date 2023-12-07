package crypto

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/crypto/blst"
	"golang.org/x/crypto/blake2b"
)

func Hash(data []byte) common.Hash {
	return blake2b.Sum256(data)
}

func PrivECDSAToHex(k *ecdsa.PrivateKey) []byte {
	return hexEncode(FromECDSA(k))
}

func PubECDSAToHex(k *ecdsa.PublicKey) []byte {
	return hexEncode(FromECDSAPub(k))
}

func PrivECDSAFromHex(k []byte) (*ecdsa.PrivateKey, error) {
	data, err := hexDecode(k)
	if err != nil {
		return nil, err
	}
	return ToECDSA(data)
}

func PubECDSAFromHex(k []byte) (*ecdsa.PublicKey, error) {
	data, err := hexDecode(k)
	if err != nil {
		return nil, err
	}
	return UnmarshalPubkey(data)
}

func hexEncode(src []byte) []byte {
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return dst
}

func hexDecode(src []byte) ([]byte, error) {
	dst := make([]byte, hex.DecodedLen(len(src)))
	_, err := hex.Decode(dst, src)
	return dst, err
}

// PopProof generate the proof of possession of the validator key and the treasury account in Autonity when onboarding a
// validator. The native input msg just contains validator's treasury address. Note: DST is already specificed at
// blst/blst/signatures.go line: 21.
func PopProof(priKey blst.SecretKey, msg []byte) ([]byte, error) {
	// the msg contains treasury address and the public key of private key.
	m := append(msg, priKey.PublicKey().Marshal()...)
	proof := priKey.Sign(Hash(m).Bytes())

	err := PopVerify(priKey.PublicKey(), proof, msg)
	if err != nil {
		return nil, err
	}
	return proof.Marshal(), nil
}

func PopVerify(pubKey blst.PublicKey, sig blst.Signature, msg []byte) error {
	m := append(msg, pubKey.Marshal()...)
	if !sig.Verify(pubKey, Hash(m).Bytes()) {
		return fmt.Errorf("cannot verify proof")
	}

	return nil
}
