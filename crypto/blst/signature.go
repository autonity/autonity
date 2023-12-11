package blst

import (
	"encoding/hex"
	"fmt"
	"sync"

	"crypto/rand"
	"github.com/pkg/errors"
	blst "github.com/supranational/blst/bindings/go"
)

// In BLS (Boneh-Lynn-Shacham) signature scheme, to generate a proof of possession, you
// should use the G2 (Group 2) setting. The reason for this is that G2 is commonly used
// for proving knowledge of secret keys in BLS signatures, and it is specifically designed
// for this purpose. G1 (Group 1) in BLS is typically used for generating public keys and
// signatures. Therefore, when creating a proof of possession in BLS, it is recommended to
// use the G2 setting.
// the domain separation tag used by the hash_to_field operation, called by hash_to_curve
var dst = []byte("BLS_SIG_BLS12381G2_XMD:SHA-256_SSWU_RO_POP_")

const scalarBytes = 32
const randBitsEntropy = 64

// BlsSignature used in the BLS signature scheme.
type BlsSignature struct {
	s *blstSignature
}

// SignatureFromBytes creates a BLS signature from a LittleEndian byte slice.
func SignatureFromBytes(sig []byte) (Signature, error) {
	if len(sig) != BLSSignatureLength {
		return nil, fmt.Errorf("signature must be %d bytes", BLSSignatureLength)
	}
	signature := new(blstSignature).Uncompress(sig)
	if signature == nil {
		return nil, errors.New("could not unmarshal bytes into signature")
	}
	// Group check signature. Do not check for infinity since an aggregated signature
	// could be infinite.
	if !signature.SigValidate(false) {
		return nil, errors.New("signature not in group")
	}
	return &BlsSignature{s: signature}, nil
}

// Verify a blst signature given a public key, a message.
//
// In IETF draft BLS specification:
// Verify(PK, message, signature) -> VALID or INVALID: a verification
//
//	algorithm that outputs VALID if signature is a valid signature of
//	message under public key PK, and INVALID otherwise.
func (s *BlsSignature) Verify(pubKey PublicKey, msg []byte) bool {
	// Signature and PKs are assumed to have been validated upon decompression!
	return s.s.Verify(false, pubKey.(*BlsPublicKey).p, false, msg, dst)
}

// AggregateVerify verifies each public key against its respective message. This is vulnerable to
// rogue public-key attack. Each user must provide a proof-of-knowledge of the public key.
//
// Note: The msgs must be distinct. For maximum performance, this method does not ensure distinct
// messages.
//
// In IETF draft BLS specification:
// AggregateVerify((PK_1, message_1), ..., (PK_n, message_n),
//
//	signature) -> VALID or INVALID: an aggregate verification
//	algorithm that outputs VALID if signature is a valid aggregated
//	signature for a collection of public keys and messages, and
//	outputs INVALID otherwise.
func (s *BlsSignature) AggregateVerify(pubKeys []PublicKey, msgs [][32]byte) bool {
	size := len(pubKeys)
	if size == 0 {
		return false
	}
	if size != len(msgs) {
		return false
	}
	msgSlices := make([][]byte, len(msgs))
	rawKeys := make([]*blstPublicKey, len(msgs))
	for i := 0; i < size; i++ {
		msgSlices[i] = msgs[i][:]
		rawKeys[i] = pubKeys[i].(*BlsPublicKey).p
	}
	// Signature and PKs are assumed to have been validated upon decompression!
	return s.s.AggregateVerify(false, rawKeys, false, msgSlices, dst)
}

// FastAggregateVerify verifies all the provided public keys with their aggregated signature.
//
// In IETF draft BLS specification:
// FastAggregateVerify(PK_1, ..., PK_n, message, signature) -> VALID
//
//	or INVALID: a verification algorithm for the aggregate of multiple
//	signatures on the same message.  This function is faster than
//	AggregateVerify.
func (s *BlsSignature) FastAggregateVerify(pubKeys []PublicKey, msg [32]byte) bool {
	if len(pubKeys) == 0 {
		return false
	}
	rawKeys := make([]*blstPublicKey, len(pubKeys))
	for i := 0; i < len(pubKeys); i++ {
		rawKeys[i] = pubKeys[i].(*BlsPublicKey).p
	}

	return s.s.FastAggregateVerify(true, rawKeys, msg[:], dst)
}

// NewAggregateSignature creates a blank aggregate signature.
func NewAggregateSignature() Signature {
	sig := blst.HashToG2([]byte{'m', 'o', 'c', 'k'}, dst).ToAffine()
	return &BlsSignature{s: sig}
}

// AggregateSignatures converts a list of signatures into a single, aggregated sig.
func AggregateSignatures(sigs []Signature) Signature {
	if len(sigs) == 0 {
		return nil
	}

	rawSigs := make([]*blstSignature, len(sigs))
	for i := 0; i < len(sigs); i++ {
		rawSigs[i] = sigs[i].(*BlsSignature).s
	}

	// Signature and PKs are assumed to have been validated upon decompression!
	signature := new(blstAggregateSignature)
	signature.Aggregate(rawSigs, false)
	return &BlsSignature{s: signature.ToAffine()}
}

// Aggregate is an alias for AggregateSignatures, defined to conform to BLS specification.
//
// In IETF draft BLS specification:
// Aggregate(signature_1, ..., signature_n) -> signature: an
//
//	aggregation algorithm that compresses a collection of signatures
//	into a single signature.
func Aggregate(sigs []Signature) Signature {
	return AggregateSignatures(sigs)
}

// VerifyMultipleSignatures verifies a non-singular set of signatures and its respective pubkeys and messages.
// This method provides a safe way to verify multiple signatures at once. We pick a number randomly from 1 to max
// uint64 and then multiply the signature by it. We continue doing this for all signatures and its respective pubkeys.
// S* = S_1 * r_1 + S_2 * r_2 + ... + S_n * r_n
// P'_{i,j} = P_{i,j} * r_i
// e(S*, G) = \prod_{i=1}^n \prod_{j=1}^{m_i} e(P'_{i,j}, M_{i,j})
// Using this we can verify multiple signatures safely.
func VerifyMultipleSignatures(sigs [][]byte, msgs [][32]byte, pubKeys []PublicKey) (bool, error) {
	if len(sigs) == 0 || len(pubKeys) == 0 {
		return false, nil
	}
	rawSigs := new(blstSignature).BatchUncompress(sigs)

	length := len(sigs)
	if length != len(pubKeys) || length != len(msgs) {
		return false, errors.Errorf("provided signatures, pubkeys and messages have differing lengths. S: %d, P: %d,M %d",
			length, len(pubKeys), len(msgs))
	}
	mulP1Aff := make([]*blstPublicKey, length)
	rawMsgs := make([]blst.Message, length)

	for i := 0; i < length; i++ {
		mulP1Aff[i] = pubKeys[i].(*BlsPublicKey).p
		rawMsgs[i] = msgs[i][:]
	}
	// Secure source of RNG
	randLock := new(sync.Mutex)

	randFunc := func(scalar *blst.Scalar) {
		var rbytes [scalarBytes]byte
		randLock.Lock()
		// Ignore error as the error will always be nil in `read` in math/rand.
		_, _ = rand.Read(rbytes[:])
		randLock.Unlock()
		scalar.FromBEndian(rbytes[:])
	}
	dummySig := new(blstSignature)

	// Validate signatures since we uncompress them here. Public keys should already be validated.
	return dummySig.MultipleAggregateVerify(rawSigs, true, mulP1Aff, false, rawMsgs, dst, randFunc, randBitsEntropy), nil
}

// Marshal a signature into a LittleEndian byte slice.
func (s *BlsSignature) Marshal() []byte {
	return s.s.Compress()
}

// Copy returns a full deep copy of a signature.
func (s *BlsSignature) Copy() Signature {
	sign := *s.s
	return &BlsSignature{s: &sign}
}

func (s *BlsSignature) Hex() string {
	return HexPrefix + hex.EncodeToString(s.Marshal())
}

// VerifyCompressed verifies that the compressed signature and pubkey
// are valid from the message provided.
func VerifyCompressed(signature, pub, msg []byte) bool {
	// Validate signature and PKs since we will uncompress them here
	return new(blstSignature).VerifyCompressed(signature, true, pub, true, msg, dst)
}
