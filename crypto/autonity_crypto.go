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

const AutonityPOPLen = 226
const AutonityNodeKeyLenInChar = 128
const AutonityNodeKeyLen = 64
const ECDSAKeyLen = 32

var (
	ErrorInvalidPOP    = errors.New("invalid Autonity POP")
	ErrorInvalidSigner = errors.New("mismatched Autonity POP signer")
)

// SaveNodeKey saves a secp256k1 private key and its derived validator BLS
// private key to the given file with restrictive permissions. The key data is saved hex-encoded.
func SaveNodeKey(file string, ecdsaKey *ecdsa.PrivateKey, consensusKey blst.SecretKey) error {
	k := hex.EncodeToString(FromECDSA(ecdsaKey))
	d := hex.EncodeToString(consensusKey.Marshal())
	return os.WriteFile(file, []byte(k+d), 0600)
}

// LoadNodeKey loads a secp256k1 private key and a derived validator BLS private key from the given file.
func LoadNodeKey(file string) (*ecdsa.PrivateKey, blst.SecretKey, error) {
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

	consensusKeyBytes, err := hex.DecodeString(string(buf[AutonityNodeKeyLenInChar/2 : AutonityNodeKeyLenInChar]))
	if err != nil {
		return nil, nil, err
	}
	consensusKey, err := blst.SecretKeyFromBytes(consensusKeyBytes)
	if err != nil {
		return nil, nil, err
	}
	return ecdsaKey, consensusKey, nil
}

// HexToNodeKey parse the hex string into a secp256k1 private key and a derived BLS private key.
func HexToNodeKey(hexKeys string) (*ecdsa.PrivateKey, blst.SecretKey, error) {
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

	consensusKey, err := blst.SecretKeyFromBytes(b[ECDSAKeyLen:AutonityNodeKeyLen])
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

// BLSPOPProof todo: (Jason) in the native BLS spec, the msg used to generate POP is just the public key, while in Autonity,
//
//	we include validator treasury to prevent from cloning during the on-boarding TX propagation in P2P network.
//	we need to double check if this none standard implementation could introduce any issue, otherwise we need to include
//	an extra signature in the on-boarding TX.
func BLSPOPProof(priKey blst.SecretKey, msg []byte) ([]byte, error) {
	// the msg contains treasury address and the public key of private key.
	m := append(msg, priKey.PublicKey().Marshal()...)
	proof := priKey.POPProof(Hash(m).Bytes())

	err := BLSPOPVerify(priKey.PublicKey(), proof, msg)
	if err != nil {
		return nil, err
	}
	return proof.Marshal(), nil
}

func BLSPOPVerify(pubKey blst.PublicKey, sig blst.Signature, msg []byte) error {
	m := append(msg, pubKey.Marshal()...)
	if !sig.POPVerify(pubKey, Hash(m).Bytes()) {
		return fmt.Errorf("cannot verify BLS POP")
	}

	return nil
}

func AutonityPOPProof(nodeKey, oracleKey *ecdsa.PrivateKey, treasuryHex string, consensusKey blst.SecretKey) ([]byte, error) {
	msg, err := hexutil.Decode(treasuryHex)
	if err != nil {
		return nil, err
	}

	hash := POPMsgHash(msg)
	// sign the msg hash with ecdsa node key and oracle key
	nodeSignature, err := Sign(hash.Bytes(), nodeKey)
	if err != nil {
		return nil, err
	}
	oracleSignature, err := Sign(hash.Bytes(), oracleKey)
	if err != nil {
		return nil, err
	}

	// generate the BLS POP
	blsPOPProof, err := BLSPOPProof(consensusKey, msg)
	if err != nil {
		return nil, err
	}

	signatures := append(append(nodeSignature[:], oracleSignature[:]...), blsPOPProof[:]...)
	return signatures, nil
}

func AutonityPOPVerify(signatures []byte, treasuryHex string, nodeAddress, oracleAddress common.Address, consensusKey []byte) error {
	if len(signatures) != AutonityPOPLen {
		return ErrorInvalidPOP
	}

	msg, err := hexutil.Decode(treasuryHex)
	if err != nil {
		return err
	}

	hash := POPMsgHash(msg)
	if err = ECDSAPOPVerify(signatures[0:common.SealLength], hash, nodeAddress); err != nil {
		return err
	}

	blsSigOffset := common.SealLength * 2
	if err = ECDSAPOPVerify(signatures[common.SealLength:blsSigOffset], hash, oracleAddress); err != nil {
		return err
	}

	// check zero signature.
	validatorSig, err := blst.SignatureFromBytes(signatures[blsSigOffset:])
	if err != nil {
		return err
	}

	// check zero public key.
	blsPubKey, err := blst.PublicKeyFromBytes(consensusKey)
	if err != nil {
		return err
	}
	return BLSPOPVerify(blsPubKey, validatorSig, msg)
}

func ECDSAPOPVerify(sig []byte, hash common.Hash, expectedSigner common.Address) error {
	signer, err := SigToAddr(hash[:], sig)
	if err != nil {
		return err
	}

	if signer != expectedSigner {
		return ErrorInvalidSigner
	}

	return nil
}

func POPMsgHash(msg []byte) common.Hash {
	// Add ethereum signed message prefix to maintain compatibility with web3.eth.sign
	// refer here : https://web3js.readthedocs.io/en/v1.2.0/web3-eth-accounts.html#sign
	prefix := fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(msg))
	hash := Keccak256Hash([]byte(prefix), msg)
	return hash
}

func GenAutonityNodeKey(keyFile string) (*ecdsa.PrivateKey, blst.SecretKey, error) {
	nodeKey, err := GenerateKey()
	if err != nil {
		return nil, nil, err
	}

	// generate random bls key, and save both node key and the validator key in node key file.
	consensusKey, err := blst.RandKey()
	if err != nil {
		return nil, nil, err
	}

	if err = SaveNodeKey(keyFile, nodeKey, consensusKey); err != nil {
		return nil, nil, err
	}

	return nodeKey, consensusKey, nil
}
