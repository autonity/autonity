package crypto

import (
	"bufio"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	c "github.com/autonity/autonity/common"
	"github.com/autonity/autonity/crypto/bls"
	"github.com/autonity/autonity/crypto/bls/common"
	"golang.org/x/crypto/blake2b"
	"os"
)

const MsgPrefix = "Auton"
const AutonityNodeKeyLenInChar = 128
const AutonityNodeKeyLen = 64
const ECDSAKeyLen = 32

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

// SaveNodeKey saves a secp256k1 private key and its derived validator BLS
// private key to the given file with restrictive permissions. The key data is saved hex-encoded.
func SaveNodeKey(file string, ecdsaKey *ecdsa.PrivateKey, validatorKey bls.SecretKey) error {
	k := hex.EncodeToString(FromECDSA(ecdsaKey))
	d := hex.EncodeToString(validatorKey.Marshal())
	return os.WriteFile(file, []byte(k+d), 0600)
}

// LoadNodeKey loads a secp256k1 private key and a derived validator BLS private key from the given file.
func LoadNodeKey(file string) (*ecdsa.PrivateKey, bls.SecretKey, error) {
	fd, err := os.Open(file)
	if err != nil {
		return nil, nil, err
	}
	defer fd.Close()

	r := bufio.NewReader(fd)
	buf := make([]byte, AutonityNodeKeyLenInChar)
	n, err := readASCII(buf, r)
	if err != nil {
		return nil, nil, err
	} else if n != len(buf) {
		return nil, nil, fmt.Errorf("key file too short, want 128 hex characters")
	}

	if err = checkKeyFileEnd(r, AutonityNodeKeyLenInChar); err != nil {
		return nil, nil, err
	}

	ecdsaKey, err := HexToECDSA(string(buf[0 : AutonityNodeKeyLenInChar/2]))
	if err != nil {
		return nil, nil, err
	}

	validatorKeyBytes, err := hex.DecodeString(string(buf[AutonityNodeKeyLenInChar/2 : AutonityNodeKeyLenInChar]))
	if err != nil {
		return nil, nil, err
	}
	validatorKey, err := bls.SecretKeyFromBytes(validatorKeyBytes)
	if err != nil {
		return nil, nil, err
	}
	return ecdsaKey, validatorKey, nil
}

// HexToNodeKey parse the hex string into a secp256k1 private key and a derived BLS private key.
func HexToNodeKey(hexKeys string) (*ecdsa.PrivateKey, bls.SecretKey, error) {
	b, err := hex.DecodeString(hexKeys)
	if byteErr, ok := err.(hex.InvalidByteError); ok {
		return nil, nil, fmt.Errorf("invalid hex character %q in node key", byte(byteErr))
	} else if err != nil {
		return nil, nil, errors.New("invalid hex data for node key")
	}

	if len(b) != AutonityNodeKeyLen {
		return nil, nil, errors.New("invalid length of hex data for node key")
	}

	ecdsaKey, err := ToECDSA(b[0:ECDSAKeyLen])
	if err != nil {
		return nil, nil, err
	}

	validatorKey, err := bls.SecretKeyFromBytes(b[ECDSAKeyLen:AutonityNodeKeyLen])
	if err != nil {
		return nil, nil, err
	}
	return ecdsaKey, validatorKey, nil
}

func Hash(data []byte) c.Hash {
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
