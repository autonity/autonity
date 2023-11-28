package crypto

import (
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/autonity/autonity/common"
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
