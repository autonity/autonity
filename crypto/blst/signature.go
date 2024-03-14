package blst

import (
	"encoding/hex"
	"fmt"
	"io"
	"sync"

	"crypto/rand"

	"github.com/pkg/errors"
	blst "github.com/supranational/blst/bindings/go"

	"github.com/autonity/autonity/common/hexutil"
	"github.com/autonity/autonity/rlp"
)

/*
Please refer to here: https://datatracker.ietf.org/doc/html/draft-irtf-cfrg-bls-signature-04#section-3
There are 3 BLS schemes that differ in handling rogue key attacks
  - basic: requires message signed by an aggregate signature to be distinct
  - message augmentation: signatures are generated over the concatenation of public key and the message
    enforcing message signed by different public key to be distinct
  - proof of possession: a separate public key called proof-of-possession is used to allow signing
    on the same message while defending against rogue key attacks
    with respective ID / domain separation tag:
  - BLS_SIG_BLS12381G2-SHA256-SSWU-RO-_NUL_
  - BLS_SIG_BLS12381G2-SHA256-SSWU-RO-_AUG_
  - BLS_SIG_BLS12381G2-SHA256-SSWU-RO-_POP_
  	- POP tag: BLS_POP_BLS12381G2-SHA256-SSWU-RO-_POP_
We implement the proof-of-possession scheme
Compared to the spec API are modified
to enforce usage of the proof-of-possession (as recommended)
*/

var generalDST = []byte("BLS_SIG_BLS12381G2_XMD:SHA-256_SSWU_RO_POP_")
var popDST = []byte("BLS_POP_BLS12381G2_XMD:SHA-256_SSWU_RO_POP_")

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

// POPVerify verify a proof of possession, it assumes that the zero public key was
// checked, the group and zero signature were checked.
func (s *BlsSignature) POPVerify(pubKey PublicKey, msg []byte) bool {
	return s.s.Verify(false, pubKey.(*BlsPublicKey).p, false, msg, popDST)
}

func (s *BlsSignature) IsZero() bool {
	return !s.s.SigValidate(true)
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
	return s.s.Verify(false, pubKey.(*BlsPublicKey).p, false, msg, generalDST)
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
	return s.s.AggregateVerify(false, rawKeys, false, msgSlices, generalDST)
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

	return s.s.FastAggregateVerify(true, rawKeys, msg[:], generalDST)
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
		// Protect against the generator returning 0. Since the scalar value is
		// derived from a big endian byte slice, we take the last byte.
		rbytes[len(rbytes)-1] |= 0x01
		scalar.FromBEndian(rbytes[:])
	}
	dummySig := new(blstSignature)

	// Validate signatures since we uncompress them here. Public keys should already be validated.
	return dummySig.MultipleAggregateVerify(rawSigs, true, mulP1Aff, false, rawMsgs, generalDST, randFunc, randBitsEntropy), nil
}

// Marshal a signature into a LittleEndian byte slice.
func (s *BlsSignature) Marshal() []byte {
	return s.s.Compress()
}

func (s *BlsSignature) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, s.Marshal())
}

func (s *BlsSignature) DecodeRLP(stream *rlp.Stream) error {
	b, err := stream.Bytes()
	if err != nil {
		return fmt.Errorf("error while decoding BLS signature: %w", err)
	}
	signature, err := SignatureFromBytes(b)
	if err != nil {
		return fmt.Errorf("error while decoding BLS signature: %w", err)
	}
	// copy inner bls signature pointer into the decoded one. Ugly but it works.
	s.s = signature.(*BlsSignature).s
	return nil
}

// for allowing JSON encoding/decoding of committed seals in the header
func (s *BlsSignature) MarshalText() ([]byte, error) {
	return hexutil.Bytes(s.Marshal()).MarshalText()
}

func (s *BlsSignature) UnmarshalText(input []byte) error {
	b := make([]byte, BLSSignatureLength)
	if err := hexutil.UnmarshalFixedText("BlsSignature", input, b); err != nil {
		return err
	}
	signature, err := SignatureFromBytes(b)
	if err != nil {
		return fmt.Errorf("error while decoding BLS signature: %w", err)
	}
	// copy inner bls signature pointer into the decoded one. Ugly but it works.
	s.s = signature.(*BlsSignature).s
	return nil
}

// Copy returns a full deep copy of a signature.
func (s *BlsSignature) Copy() *BlsSignature {
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
	return new(blstSignature).VerifyCompressed(signature, true, pub, true, msg, generalDST)
}
