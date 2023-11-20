package crypto

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/autonity/autonity/crypto/bls"
	"github.com/autonity/autonity/crypto/bls/common"
)

const MsgPrefix = "Auton"

func ValidateValidatorKeyProof(pubKey common.BLSPublicKey, sig common.BLSSignature, msg []byte) error {
	m := append([]byte(MsgPrefix), msg...)
	hash := Keccak256Hash(m)
	if !sig.Verify(pubKey, hash.Bytes()) {
		return fmt.Errorf("cannot verify proof")
	}

	return nil
}

func GenerateValidatorKeyProof(priKey bls.SecretKey, msg []byte) ([]byte, error) {
	// generate bls key ownership proof
	m := append([]byte(MsgPrefix), msg...)
	hash := Keccak256Hash(m)
	proof := priKey.Sign(hash.Bytes())

	err := ValidateValidatorKeyProof(priKey.PublicKey(), proof, msg)
	if err != nil {
		return nil, err
	}
	return proof.Marshal(), nil
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
