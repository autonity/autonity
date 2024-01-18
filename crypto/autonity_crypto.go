package crypto

import (
	"bufio"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/hexutil"
	"github.com/autonity/autonity/crypto/blst"
	"golang.org/x/crypto/blake2b"
	"os"
)

const ECDSAKeyLen = 32
const AutonityPOPLen = 226
const ECDSAKeyLenInChar = ECDSAKeyLen * 2
const AutonityKeysLen = ECDSAKeyLen + blst.BLSSecretKeyLength
const AutonityKeysLenInChar = AutonityKeysLen * 2

var (
	ErrorInvalidPOP    = errors.New("invalid Autonity POP")
	ErrorInvalidSigner = errors.New("mismatched Autonity POP signer")
)

// SaveAutonityKeys saves a secp256k1 private key and the consensus key to the given file with restrictive permissions.
// The key data is saved hex-encoded.
func SaveAutonityKeys(file string, ecdsaKey *ecdsa.PrivateKey, consensusKey blst.SecretKey) error {
	k := hex.EncodeToString(FromECDSA(ecdsaKey))
	d := hex.EncodeToString(consensusKey.Marshal())
	return os.WriteFile(file, []byte(k+d), 0600)
}

// LoadAutonityKeys loads a secp256k1 private key and a consensus private key from the given file.
func LoadAutonityKeys(file string) (*ecdsa.PrivateKey, blst.SecretKey, error) {
	fd, err := os.Open(file)
	if err != nil {
		return nil, nil, err
	}
	defer fd.Close()

	r := bufio.NewReader(fd)
	buf := make([]byte, AutonityKeysLenInChar)
	n, err := readASCII(buf, r)
	if err != nil {
		return nil, nil, err
	} else if n != len(buf) {
		return nil, nil, fmt.Errorf("key file too short, want 128 hex characters")
	}

	if err = checkKeyFileEnd(r, AutonityKeysLenInChar); err != nil {
		return nil, nil, err
	}

	ecdsaKey, err := HexToECDSA(string(buf[0:ECDSAKeyLenInChar]))
	if err != nil {
		return nil, nil, err
	}

	consensusKeyBytes, err := hex.DecodeString(string(buf[ECDSAKeyLenInChar:AutonityKeysLenInChar]))
	if err != nil {
		return nil, nil, err
	}
	consensusKey, err := blst.SecretKeyFromBytes(consensusKeyBytes)
	if err != nil {
		return nil, nil, err
	}
	return ecdsaKey, consensusKey, nil
}

// HexToAutonityKeys parse the hex string into a secp256k1 private key and a derived BLS private key.
func HexToAutonityKeys(hexKeys string) (*ecdsa.PrivateKey, blst.SecretKey, error) {
	b, err := hex.DecodeString(hexKeys)
	if byteErr, ok := err.(hex.InvalidByteError); ok {
		return nil, nil, fmt.Errorf("invalid hex character %q in node key", byte(byteErr))
	} else if err != nil {
		return nil, nil, errors.New("invalid hex data for node key")
	}

	if len(b) != AutonityKeysLen {
		return nil, nil, errors.New("invalid length of hex data for node key")
	}

	ecdsaKey, err := ToECDSA(b[0:ECDSAKeyLen])
	if err != nil {
		return nil, nil, err
	}

	consensusKey, err := blst.SecretKeyFromBytes(b[ECDSAKeyLen:AutonityKeysLen])
	if err != nil {
		return nil, nil, err
	}
	return ecdsaKey, consensusKey, nil
}

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

// BLSPOPProof generate POP of BLS private key of Autonity protocol, the hash input start with a prefix of treasury
// address and ended with the public key of the secrete key, since we don't want the POP being cloned during the
// propagation of the on-boarding TX. Thus, this POP generation is different from the spec of BLS, which means we have
// a compatibility issue with a standard POP generation/verification implementation.
func BLSPOPProof(priKey blst.SecretKey, treasury []byte) ([]byte, error) {
	// the treasury contains treasury address bytes.
	var buff []byte
	copy(buff, treasury)
	treasuryPubKey := append(buff, priKey.PublicKey().Marshal()...)
	proof := priKey.POPProof(Hash(treasuryPubKey).Bytes())

	err := BLSPOPVerify(priKey.PublicKey(), proof, treasury)
	if err != nil {
		return nil, err
	}
	return proof.Marshal(), nil
}

// BLSPOPVerify verifies the POP provided by an on-boarding validator, it assumes the public key and signature was
// checked with infinite and group.
func BLSPOPVerify(pubKey blst.PublicKey, sig blst.Signature, treasury []byte) error {
	var buff []byte
	copy(buff, treasury)
	treasuryPubKey := append(buff, pubKey.Marshal()...)
	if !sig.POPVerify(pubKey, Hash(treasuryPubKey).Bytes()) {
		return fmt.Errorf("cannot verify BLS POP")
	}

	return nil
}

func AutonityPOPProof(nodeKey, oracleKey *ecdsa.PrivateKey, treasuryHex string, consensusKey blst.SecretKey) ([]byte, error) {
	treasury, err := hexutil.Decode(treasuryHex)
	if err != nil {
		return nil, err
	}

	hash := POPMsgHash(treasury)
	// sign the treasury hash with ecdsa node key and oracle key
	nodeSignature, err := Sign(hash.Bytes(), nodeKey)
	if err != nil {
		return nil, err
	}
	oracleSignature, err := Sign(hash.Bytes(), oracleKey)
	if err != nil {
		return nil, err
	}

	// generate the BLS POP
	blsPOPProof, err := BLSPOPProof(consensusKey, treasury)
	if err != nil {
		return nil, err
	}

	signatures := append(append(nodeSignature[:], oracleSignature[:]...), blsPOPProof[:]...)
	return signatures, nil
}

func POPMsgHash(msg []byte) common.Hash {
	// Add ethereum signed message prefix to maintain compatibility with web3.eth.sign
	// refer here : https://web3js.readthedocs.io/en/v1.2.0/web3-eth-accounts.html#sign
	prefix := fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(msg))
	hash := Keccak256Hash([]byte(prefix), msg)
	return hash
}

func GenAutonityKeys() (*ecdsa.PrivateKey, blst.SecretKey, error) {
	nodeKey, err := GenerateKey()
	if err != nil {
		return nil, nil, err
	}

	// generate random bls key, and save both node key and the validator key in node key file.
	consensusKey, err := blst.RandKey()
	if err != nil {
		return nil, nil, err
	}

	return nodeKey, consensusKey, nil
}
