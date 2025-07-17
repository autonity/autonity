package blst

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/pkg/errors"

	"github.com/autonity/autonity/common/hexutil"
	"github.com/autonity/autonity/log"
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

/* Dealing with the 0 signature:
*  - do not accept 0 signature and/or 0 public key when verifying an individual signature. This scenario
*    can never happen since we do group check the public keys at validator registration.
*  - accept 0 sig and 0 public key if verifying an aggregate signature. This can theoretically happen, even
* 	 by chance. Moreover it is not possible to verify if a subset of the signatures being aggregated are 0.
	 The important thing is that all nodes reach to the same result about signatures regardless of the order
*    or method used for verification.
*/

var generalDST = []byte("BLS_SIG_BLS12381G2_XMD:SHA-256_SSWU_RO_POP_")
var popDST = []byte("BLS_POP_BLS12381G2_XMD:SHA-256_SSWU_RO_POP_")

const scalarBytes = 32
const randBitsEntropy = 64

// Individual signatures MUST be checked for 0 value at decoding / preValidation phase.
// Aggregate signatures are allowed to be 0.
// Therefore, the default behavior in Autonity is to assume the (zero sig, zero key)
// pairs valid during signature verification, since individual zero signatures are
// pre-filtered at previous stages.
// IMPORTANT: PLEASE do not use this constant blindly, carefully consider whether to
// pass true or false based on the use case.
const DefaultAssumeZeroValid = true

// BlsSignature used in the BLS signature scheme.
type BlsSignature struct {
	s *blstSignature
}

// SignatureFromBytes creates a BLS signature from a LittleEndian byte slice.
// it does not check it for infinity.
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

// IsZero groups check the signature AND verifies it is not the 0 signature
func (s *BlsSignature) IsZero() bool {
	return !s.s.SigValidate(true)
}

// AggregateSignatures converts a list of signatures into a single, aggregated sig.
// does not group check or check for infinity the aggregate signature
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

// POPVerify verify a proof of possession, it assumes that the zero public key was
// checked, the group and zero signature were checked.
func (s *BlsSignature) POPVerify(pubKey PublicKey, msg []byte) bool {
	return s.s.Verify(false, pubKey.(*BlsPublicKey).p, false, msg, popDST)
}

// Verify a blst signature given a public key, a message.
// if the last parameter `assumeZeroValid` is:
//   - true: zero signature and zero key pair will be considered as valid
//   - false: zero signature and zero key pair will be considered as invalid
func (s *BlsSignature) Verify(pubKey PublicKey, msg []byte, assumeZeroValid bool) bool {
	// check for zero signature and zero key only if it is assumed valid in the first place
	if assumeZeroValid && s.IsZero() && !pubKey.Validate() {
		return true
	}
	// Signature and PKs are assumed to have been validated upon decompression!
	return s.s.Verify(false, pubKey.(*BlsPublicKey).p, false, msg, generalDST)
}

// AggregateVerify verifies each public key against its respective message. This is vulnerable to
// rogue public-key attack. Each user must provide a proof-of-knowledge of the public key.
//
// NOTE:
// - The msgs MUST BE DISTINCT!. For maximum performance, this method does not ensure distinct messages.
// - 0 signatures/public keys should be filtered out by the caller
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

// Inspired from the logic of VerifyMultipleSignatures.
// It ensures that all the individual signatures are valid if the aggregate is valid. It is a defence for the consensus attack (see TestBlsAttacks)
// e(P1, S1 * r1 + S2 * r2 + ... + Sn * rn) = e(X1 * r1 + X2 * r2 + ... + Xn * rn, H(m))
func FastAggregateVerifyBatch(sigs []Signature, pubkeys []PublicKey, msg [32]byte) bool {
	n := len(sigs)

	if n != len(pubkeys) || n == 0 {
		return false
	}

	// Secure source of RNG
	randomScalar := func() *blstScalar {
		var rbytes [scalarBytes]byte

		_, err := rand.Read(rbytes[:])
		if err != nil {
			panic("Cannot source randomness")
		}

		// Protect against the generator returning 0. Since the scalar value is
		// derived from a big endian byte slice, we take the last byte.
		rbytes[len(rbytes)-1] |= 0x01
		scalar := new(blstScalar)
		scalar.FromBEndian(rbytes[:])
		return scalar
	}

	// extract raw signatures and public keys (EC points) and generate random scalars
	var rawKeys blstPublicKeySet
	var rawSigs blstSignatureSet
	var scalars []*blstScalar
	for i := 0; i < n; i++ {
		rawKeys = append(rawKeys, *pubkeys[i].(*BlsPublicKey).p)
		rawSigs = append(rawSigs, *sigs[i].(*BlsSignature).s)
		scalars = append(scalars, randomScalar())
	}

	// multiply keys and sigs with scalar using Pippenger Algorithm
	// https://hackmd.io/@drouyang/SyYwhWIso
	aggregatedKeyAffine := rawKeys.Mult(scalars, randBitsEntropy).ToAffine()
	aggregatedSignatureAffine := rawSigs.Mult(scalars, randBitsEntropy).ToAffine()
	aggregatedSignature := &BlsSignature{s: aggregatedSignatureAffine}
	aggregatedKey := &BlsPublicKey{p: aggregatedKeyAffine}

	// I believe this should never happen, but if it does print some debugging info and fail fast
	if !aggregatedSignature.s.SigValidate(false) {
		var serializedSigs [][]byte
		var serializedKeys [][]byte
		var serializedScalars [][]byte
		for i := range rawSigs {
			serializedSigs = append(serializedSigs, rawSigs[i].Serialize())
			serializedKeys = append(serializedKeys, rawKeys[i].Serialize())
			serializedScalars = append(serializedScalars, scalars[i].Serialize())
		}
		log.Error("Unexpected error: aggregate signature is not part of the group, please report it", "sigs", serializedSigs, "rawKeys", serializedKeys, "scalars", serializedScalars)
		panic("unexpected error: aggregate signature is not part of the group")
	}

	return aggregatedSignature.Verify(aggregatedKey, msg[:], DefaultAssumeZeroValid)
}

/* ----------------
*  utility funcs
*  ----------------
 */

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

// for allowing JSON encoding/decoding of quorum certificate in the header
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
