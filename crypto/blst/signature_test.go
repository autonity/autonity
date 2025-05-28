package blst

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	mrand "math/rand"
	"testing"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/rlp"
	blst "github.com/supranational/blst/bindings/go"

	"github.com/stretchr/testify/require"
)

func TestSignVerify(t *testing.T) {
	priv, err := RandKey()
	require.NoError(t, err)
	pub := priv.PublicKey()
	msg := []byte("hello")
	sig := priv.Sign(msg)
	require.Equal(t, true, sig.Verify(pub, msg, DefaultAssumeZeroValid), "Signature did not verify")
}

func TestPOPVerify(t *testing.T) {
	priv, err := RandKey()
	require.NoError(t, err)
	pub := priv.PublicKey()
	popProof := priv.POPProof(pub.Marshal())
	require.True(t, popProof.POPVerify(pub, pub.Marshal()))
}

// since the keys are different for each signature, the order of pubkey for aggregation verification matters.
func TestAggregateVerifyWithDifferentKeys(t *testing.T) {
	pubkeys := make([]PublicKey, 0, 100)
	sigs := make([]Signature, 0, 100)
	var msgs [][32]byte
	for i := 0; i < 100; i++ {
		// with each different key signed different msg.
		msg := [32]byte{'h', 'e', 'l', 'l', 'o', byte(i)}
		priv, err := RandKey()
		require.NoError(t, err)
		pub := priv.PublicKey()
		sig := priv.Sign(msg[:])
		pubkeys = append(pubkeys, pub)
		sigs = append(sigs, sig)
		msgs = append(msgs, msg)
	}
	aggSig := AggregateSignatures(sigs)
	require.Equal(t, true, aggSig.AggregateVerify(pubkeys, msgs), "Signature did not verify")
}

// we want to test and verify if aggregated signatures are deterministic when they are aggregated with different orders
// of signatures.
func TestSignatureAggregationWithDifferentOrdersOfSignatures(t *testing.T) {
	msgHello := [32]byte{'h', 'e', 'l', 'l', 'o'}
	//msgMsg := [32]byte{'m', 'y'}
	//msgWorld := [32]byte{'w', 'o', 'r', 'l', 'd'}

	key1, err := RandKey()
	require.NoError(t, err)
	key2, err := RandKey()
	require.NoError(t, err)
	key3, err := RandKey()
	require.NoError(t, err)

	helloSig1 := key1.Sign(msgHello[:])
	helloSig2 := key2.Sign(msgHello[:])
	helloSig3 := key3.Sign(msgHello[:])

	helloSigs1 := make([]Signature, 0, 3)
	helloSigs1 = append(helloSigs1, helloSig1)
	helloSigs1 = append(helloSigs1, helloSig2)
	helloSigs1 = append(helloSigs1, helloSig3)
	helloAggSig1 := AggregateSignatures(helloSigs1)

	helloSigs2 := make([]Signature, 0, 3)
	helloSigs2 = append(helloSigs2, helloSig2)
	helloSigs2 = append(helloSigs2, helloSig1)
	helloSigs2 = append(helloSigs2, helloSig3)
	helloAggSig2 := AggregateSignatures(helloSigs2)

	t.Log("helloAggSig1: ", helloAggSig1.Hex())
	t.Log("helloAggSig2: ", helloAggSig2.Hex())
	require.Equal(t, helloAggSig1, helloAggSig2)
}

// since the key is the same for all signatures, the order of pubkey for aggregation verification does not matter in this test.
func TestAggregateVerifyWithSameKey(t *testing.T) {
	pubkeys := make([]PublicKey, 0, 100)
	sigs := make([]Signature, 0, 100)
	var msgs [][32]byte
	// generate unique blst key.
	priv, err := RandKey()
	require.NoError(t, err)
	pub := priv.PublicKey()
	// sign different msgs with the key.
	for i := 0; i < 100; i++ {
		// with each different key signed different msg.
		msg := [32]byte{'h', 'e', 'l', 'l', 'o', byte(i)}
		sig := priv.Sign(msg[:])

		// ordering the keys, signatures, and msgs in memory.
		pubkeys = append(pubkeys, pub)
		sigs = append(sigs, sig)
		msgs = append(msgs, msg)
	}
	aggSig := AggregateSignatures(sigs)

	// verification of aggregated signature.
	require.Equal(t, true, aggSig.AggregateVerify(pubkeys, msgs), "Signature did not verify")
}

func Test20MinEpoch1ValidatorSigAggVerification(t *testing.T) {
	epochLength := 1200
	pubkeys := make([]PublicKey, 0, epochLength)
	sigs := make([]Signature, 0, epochLength)
	var msgs [][32]byte
	// generate unique blst key.
	priv, err := RandKey()
	pub := priv.PublicKey()
	require.NoError(t, err)
	// sign different msgs with the key.
	for i := 0; i < epochLength; i++ {
		// with the key of same validator to sign different msgs.
		msg := [32]byte{'h', 'e', 'l', 'l', 'o', byte(i)}
		sig := priv.Sign(msg[:])

		// ordering the keys, signatures, and msgs in memory.
		pubkeys = append(pubkeys, pub)
		sigs = append(sigs, sig)
		msgs = append(msgs, msg)
	}
	// aggregate the signatures of single validator into one aggregated signature.
	aggSig := AggregateSignatures(sigs)

	// verification of aggregated signatures of the validator.
	require.Equal(t, true, aggSig.AggregateVerify(pubkeys, msgs), "Signature did not verify")
}

func Test10MinEpoch1ValidatorSigAggVerification(t *testing.T) {
	epochLength := 600
	pubkeys := make([]PublicKey, 0, epochLength)
	sigs := make([]Signature, 0, epochLength)
	var msgs [][32]byte
	// generate unique blst key.
	priv, err := RandKey()
	pub := priv.PublicKey()
	require.NoError(t, err)
	// sign different msgs with the key.
	for i := 0; i < epochLength; i++ {
		// with the key of same validator to sign different msgs.
		msg := [32]byte{'h', 'e', 'l', 'l', 'o', byte(i)}
		sig := priv.Sign(msg[:])

		// ordering the keys, signatures, and msgs in memory.
		pubkeys = append(pubkeys, pub)
		sigs = append(sigs, sig)
		msgs = append(msgs, msg)
	}
	// aggregate the signatures of single validator into one aggregated signature.
	aggSig := AggregateSignatures(sigs)

	// verification of aggregated signature.
	require.Equal(t, true, aggSig.AggregateVerify(pubkeys, msgs), "Signature did not verify")
}

func Test5MinEpoch1ValidatorSigAggVerification(t *testing.T) {
	epochLength := 300
	pubkeys := make([]PublicKey, 0, epochLength)
	sigs := make([]Signature, 0, epochLength)
	var msgs [][32]byte
	// generate unique blst key.
	priv, err := RandKey()
	pub := priv.PublicKey()
	require.NoError(t, err)
	// sign different msgs with the key.
	for i := 0; i < epochLength; i++ {
		// with the key of same validator to sign different msgs.
		msg := [32]byte{'h', 'e', 'l', 'l', 'o', byte(i)}
		sig := priv.Sign(msg[:])

		// ordering the keys, signatures, and msgs in memory.
		pubkeys = append(pubkeys, pub)
		sigs = append(sigs, sig)
		msgs = append(msgs, msg)
	}
	// aggregate the signatures of single validator into one aggregated signature.
	aggSig := AggregateSignatures(sigs)

	// verification of aggregated signature.
	require.Equal(t, true, aggSig.AggregateVerify(pubkeys, msgs), "Signature did not verify")
}

// if the msg is distinct, then the order of public key does not impact the aggregation verification.
func Test20MinEpoch21ValidatorsOnSameMsgFastAggregateVerifyBatch(t *testing.T) {
	numOfValidator := 21
	epochLength := 1200
	skeys := make([]SecretKey, 0, numOfValidator)
	pubkeys := make([]PublicKey, 0, numOfValidator)
	for i := 0; i < numOfValidator; i++ {
		priv, err := RandKey()
		require.NoError(t, err)
		skeys = append(skeys, priv)
		pubkeys = append(pubkeys, priv.PublicKey())
	}

	// assume we have 2 rounds for each height, we only count step: preVote, and preCommit.
	for i := 0; i < epochLength*2*2; i++ {
		msg := [32]byte{'h', 'e', 'l', 'l', 'o', byte(i)}
		sigs := make([]Signature, 0, numOfValidator)
		for i := 0; i < numOfValidator; i++ {
			// with different key to sign a distinct msg for per voting step..
			sig := skeys[i].Sign(msg[:])
			sigs = append(sigs, sig)
		}
		require.Equal(t, true, FastAggregateVerifyBatch(sigs, pubkeys, msg), "Signature did not verify")
	}
}

// if the msg is distinct, then the order of public key does not impact the aggregation verification.
func TestFastAggregateVerifyBatch2(t *testing.T) {
	pubkeys := make([]PublicKey, 0, 100)
	sigs := make([]Signature, 0, 100)
	msg := [32]byte{'h', 'e', 'l', 'l', 'o'}
	for i := 0; i < 100; i++ {
		// with different key to sign a distinct msg.
		priv, err := RandKey()
		require.NoError(t, err)
		pub := priv.PublicKey()
		sig := priv.Sign(msg[:])
		pubkeys = append(pubkeys, pub)
		sigs = append(sigs, sig)
	}

	require.Equal(t, true, FastAggregateVerifyBatch(sigs, pubkeys, msg), "Signature did not verify")
}

// newAggregateSignature creates a blank aggregate signature.
func newAggregateSignature() Signature {
	sig := blst.HashToG2([]byte{'m', 'o', 'c', 'k'}, generalDST).ToAffine()
	return &BlsSignature{s: sig}
}

func TestFastAggregateVerifyBatch_ReturnsFalseOnEmptyPubKeyList(t *testing.T) {
	var pubkeys []PublicKey
	msg := [32]byte{'h', 'e', 'l', 'l', 'o'}

	aggSig := newAggregateSignature()
	require.Equal(t, false, FastAggregateVerifyBatch([]Signature{aggSig}, pubkeys, msg), "Expected FastAggregateVerify to return false with empty input ")
}

func TestSignatureFromBytes(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		err   error
	}{
		{
			name: "Nil",
			err:  errors.New("signature must be 96 bytes"),
		},
		{
			name:  "Empty",
			input: []byte{},
			err:   errors.New("signature must be 96 bytes"),
		},
		{
			name:  "Short",
			input: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			err:   errors.New("signature must be 96 bytes"),
		},
		{
			name:  "Long",
			input: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			err:   errors.New("signature must be 96 bytes"),
		},
		{
			name:  "Bad",
			input: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			err:   errors.New("could not unmarshal bytes into signature"),
		},
		{
			name:  "Good",
			input: []byte{0xab, 0xb0, 0x12, 0x4c, 0x75, 0x74, 0xf2, 0x81, 0xa2, 0x93, 0xf4, 0x18, 0x5c, 0xad, 0x3c, 0xb2, 0x26, 0x81, 0xd5, 0x20, 0x91, 0x7c, 0xe4, 0x66, 0x65, 0x24, 0x3e, 0xac, 0xb0, 0x51, 0x00, 0x0d, 0x8b, 0xac, 0xf7, 0x5e, 0x14, 0x51, 0x87, 0x0c, 0xa6, 0xb3, 0xb9, 0xe6, 0xc9, 0xd4, 0x1a, 0x7b, 0x02, 0xea, 0xd2, 0x68, 0x5a, 0x84, 0x18, 0x8a, 0x4f, 0xaf, 0xd3, 0x82, 0x5d, 0xaf, 0x6a, 0x98, 0x96, 0x25, 0xd7, 0x19, 0xcc, 0xd2, 0xd8, 0x3a, 0x40, 0x10, 0x1f, 0x4a, 0x45, 0x3f, 0xca, 0x62, 0x87, 0x8c, 0x89, 0x0e, 0xca, 0x62, 0x23, 0x63, 0xf9, 0xdd, 0xb8, 0xf3, 0x67, 0xa9, 0x1e, 0x84},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			res, err := SignatureFromBytes(test.input)
			if test.err != nil {
				require.NotEqual(t, nil, err, "No error returned")
				require.Error(t, test.err, err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, 0, bytes.Compare(res.Marshal(), test.input))
			}
		})
	}
}

func TestCopy(t *testing.T) {
	priv, err := RandKey()
	require.NoError(t, err)
	key, ok := priv.(*bls12SecretKey)
	require.Equal(t, true, ok)

	signatureA := &BlsSignature{s: new(blstSignature).Sign(key.p, []byte("foo"), generalDST)}
	signatureB := signatureA.Copy()
	require.Equal(t, true, bytes.Equal(signatureA.Marshal(), signatureB.Marshal()))

	signatureA.s.Sign(key.p, []byte("bar"), generalDST)
	require.Equal(t, false, bytes.Equal(signatureA.Marshal(), signatureB.Marshal()))
}

func TestSignature_MarshalUnMarshal(t *testing.T) {
	priv, err := RandKey()
	require.NoError(t, err)
	key, ok := priv.(*bls12SecretKey)
	require.Equal(t, true, ok)

	signatureA := &BlsSignature{s: new(blstSignature).Sign(key.p, []byte("foo"), generalDST)}
	signatureBytes := signatureA.Marshal()

	signatureB, err := SignatureFromBytes(signatureBytes)
	require.NoError(t, err)
	require.Equal(t, true, bytes.Equal(signatureA.Marshal(), signatureB.Marshal()))
}

func TestSignature_Hex(t *testing.T) {
	priv, err := RandKey()
	require.NoError(t, err)
	key, ok := priv.(*bls12SecretKey)
	require.Equal(t, true, ok)

	signatureA := &BlsSignature{s: new(blstSignature).Sign(key.p, []byte("foo"), generalDST)}
	str := signatureA.Hex()
	b, err := hex.DecodeString(str[2:])
	require.NoError(t, err)

	signatureB, err := SignatureFromBytes(b)
	require.NoError(t, err)
	require.Equal(t, true, bytes.Equal(signatureA.Marshal(), signatureB.Marshal()))
}

func TestSignatureEncodeDecode(t *testing.T) {
	priv, err := RandKey()
	require.NoError(t, err)
	signature := priv.Sign([]byte{0xca, 0xfe})

	signatureRLP, err := rlp.EncodeToBytes(signature)
	require.NoError(t, err)

	decodedSig := &BlsSignature{}
	err = rlp.Decode(bytes.NewReader(signatureRLP), decodedSig)
	require.NoError(t, err)

	require.Equal(t, signature.Marshal(), decodedSig.Marshal())
}

func opposite(t *testing.T, key SecretKey) SecretKey {
	keySerialized := key.Marshal()
	oppositeKey, err := SecretKeyFromBytes(keySerialized)
	require.NoError(t, err)

	oppositeKey.(*bls12SecretKey).p.SubAssign(key.(*bls12SecretKey).p)
	oppositeKey.(*bls12SecretKey).p.SubAssign(key.(*bls12SecretKey).p)
	return oppositeKey
}

func generateZeroKeyPair(t *testing.T) (SecretKey, SecretKey) {
	x1, err := RandKey()
	require.NoError(t, err)

	x2 := opposite(t, x1)

	X1 := x1.PublicKey()
	X2 := x2.PublicKey()

	// individual pubkeys are valid
	require.True(t, X1.(*BlsPublicKey).p.KeyValidate())
	require.True(t, X2.(*BlsPublicKey).p.KeyValidate())

	// aggregate pubkey is not valid, since it is the infinite pubkey
	zeroKey, err := AggregatePublicKeys([]PublicKey{X1, X2})
	require.NoError(t, err)
	require.False(t, zeroKey.(*BlsPublicKey).p.KeyValidate())

	return x1, x2
}

// some tests taken from paper "0" of Nguyen Thoi Minh Quan
func TestBlsAttacks(t *testing.T) {
	x1, x2 := generateZeroKeyPair(t)

	X1 := x1.PublicKey()
	X2 := x2.PublicKey()
	// individual pubkeys are valid
	require.True(t, X1.(*BlsPublicKey).p.KeyValidate())
	require.True(t, X2.(*BlsPublicKey).p.KeyValidate())

	aggX, err := AggregatePublicKeys([]PublicKey{X1, X2})
	require.NoError(t, err)
	// aggregate pubkey is not valid, since it is the infinite pubkey
	require.False(t, aggX.(*BlsPublicKey).p.KeyValidate())

	// The autonity policy is that a 0 signature can be considered valid if the
	// aggregated public key of the signers is also 0.
	//
	// 0 signatures should be rejected if:
	// - the signature is not an aggregate (dealt with at message decoding/prevalidation)
	// - the sum of the public key of the signers is not 0 (as it implies the signature is not valid)

	t.Run("Splitting zero - verification consistency", func(t *testing.T) {
		m := common.Hash{0xca, 0xfe}
		sig1 := x1.Sign(m[:])
		sig2 := x2.Sign(m[:])

		aggSig := AggregateSignatures([]Signature{sig1, sig2})
		require.True(t, aggSig.IsZero())

		// individual signatures are valid
		require.True(t, sig1.Verify(X1, m[:], DefaultAssumeZeroValid))
		require.True(t, sig2.Verify(X2, m[:], DefaultAssumeZeroValid))

		// 0 aggregate signature with 0 public key is considerate valid by all validation methods
		// except for AggregateVerify
		require.True(t, aggSig.Verify(aggX, m[:], DefaultAssumeZeroValid))
		require.True(t, FastAggregateVerifyBatch([]Signature{sig1, sig2}, []PublicKey{X1, X2}, m))
		require.True(t, FastAggregateVerifyBatch([]Signature{aggSig}, []PublicKey{aggX}, m))
		require.True(t, aggSig.AggregateVerify([]PublicKey{X1, X2}, [][32]byte{m, m}))

		// zero sig should be filtered out by the caller
		require.False(t, aggSig.AggregateVerify([]PublicKey{aggX}, [][32]byte{m}))

		// all these cases should fail
		require.False(t, aggSig.Verify(X1, m[:], DefaultAssumeZeroValid))
		require.False(t, aggSig.Verify(X2, m[:], DefaultAssumeZeroValid))
		require.False(t, sig1.Verify(aggX, m[:], DefaultAssumeZeroValid))

		require.False(t, FastAggregateVerifyBatch([]Signature{aggSig}, []PublicKey{X1}, m))
		require.False(t, FastAggregateVerifyBatch([]Signature{aggSig}, []PublicKey{X2}, m))
		require.False(t, FastAggregateVerifyBatch([]Signature{sig1, sig2}, []PublicKey{X2, X1}, m))

		require.False(t, aggSig.AggregateVerify([]PublicKey{X1}, [][32]byte{m}))

	})
	t.Run("Splitting zero - non message binding", func(t *testing.T) {
		// The user publishes signature sig3.
		m3 := common.Hash{0xca, 0xfe}
		x3, err := RandKey()
		require.NoError(t, err)
		X3 := x3.PublicKey()
		sig3 := x3.Sign(m3[:])

		m := common.Hash{0xff, 0xff}
		aggsig := sig3

		// aggsig = sig3 is a valid signature for (m,m3,m) if we use AggregateVerify. Note that the malicious party doesn't even have to sign m.
		// however this is true only if we allow the usage of AggregateVerify with non-distinct messages
		require.True(t, aggsig.AggregateVerify([]PublicKey{X1, X3, X2}, [][32]byte{m, m3, m}))

		// if the 0 key is aggregated, the 0 signature is detected as invalid
		X12, err := AggregatePublicKeys([]PublicKey{X1, X2})
		require.False(t, aggsig.AggregateVerify([]PublicKey{X12, X3}, [][32]byte{m, m3}))
	})
	t.Run("Splitting zero attack - non key binding", func(t *testing.T) {
		m := common.Hash{0xca, 0xfe}
		x3, err := RandKey()
		require.NoError(t, err)
		X3 := x3.PublicKey()
		sig3 := x3.Sign(m[:])

		X3PlusZero, err := AggregatePublicKeys([]PublicKey{X1, X2, X3})
		require.NoError(t, err)

		// same sig is valid both for X3 only and for {X1,X2,X3}
		require.True(t, FastAggregateVerifyBatch([]Signature{sig3}, []PublicKey{X3}, m))
		require.True(t, FastAggregateVerifyBatch([]Signature{sig3}, []PublicKey{X3PlusZero}, m))
	})
	t.Run("Consensus attack", func(t *testing.T) {
		m := common.Hash{0xca, 0xfe}

		x3, err := RandKey()
		require.NoError(t, err)
		X3 := x3.PublicKey()
		sig3 := x3.Sign(m[:])

		x4, err := RandKey()
		require.NoError(t, err)
		X4 := x4.PublicKey()
		sig4 := x4.Sign(m[:])

		// offset both signatures of the same value in opposite directions
		// we can do that by aggregating with signatures from x1 and x2 (since x1 = -x2)
		sig1 := x1.Sign(m[:])
		sig2 := x2.Sign(m[:])

		sig3Offsetted := AggregateSignatures([]Signature{sig3, sig1})
		sig4Offsetted := AggregateSignatures([]Signature{sig4, sig2})

		// individual signatures are not valid anymore
		require.False(t, sig3Offsetted.Verify(X3, m[:], DefaultAssumeZeroValid))
		require.False(t, sig4Offsetted.Verify(X4, m[:], DefaultAssumeZeroValid))

		// aggregate signature is not valid with FastAggregateVerifyBatch
		require.False(t, FastAggregateVerifyBatch([]Signature{sig3Offsetted, sig4Offsetted}, []PublicKey{X3, X4}, m))
	})
	t.Run("Consensus attack + splitting zero", func(t *testing.T) {
		m := common.Hash{0xca, 0xfe}

		x3, x4 := generateZeroKeyPair(t)

		X3 := x3.PublicKey()
		sig3 := x3.Sign(m[:])

		X4 := x4.PublicKey()
		sig4 := x4.Sign(m[:])

		sig34 := AggregateSignatures([]Signature{sig3, sig4})
		require.True(t, sig34.IsZero())

		// offset both signatures of the same value in opposite directions
		// we can do that by aggregating with signatures from x1 and x2 (since x1 = -x2)
		sig1 := x1.Sign(m[:])
		sig2 := x2.Sign(m[:])

		sig3Offsetted := AggregateSignatures([]Signature{sig3, sig1})
		require.False(t, sig3Offsetted.IsZero())
		sig4Offsetted := AggregateSignatures([]Signature{sig4, sig2})
		require.False(t, sig4Offsetted.IsZero())

		sig34Offsetted := AggregateSignatures([]Signature{sig3Offsetted, sig4Offsetted})
		require.True(t, sig34Offsetted.IsZero())

		// individual signatures are not valid anymore
		require.False(t, sig3Offsetted.Verify(X3, m[:], DefaultAssumeZeroValid))
		require.False(t, sig4Offsetted.Verify(X4, m[:], DefaultAssumeZeroValid))

		// aggregate signature is not valid with FastAggregateVerifyBatch
		require.False(t, FastAggregateVerifyBatch([]Signature{sig3Offsetted, sig4Offsetted}, []PublicKey{X3, X4}, m))
	})
	t.Run("Consensus attack + splitting zero 2", func(t *testing.T) {
		m := common.Hash{0xca, 0xfe}

		x3, err := RandKey()
		require.NoError(t, err)

		x4, err := RandKey()
		require.NoError(t, err)

		X3 := x3.PublicKey()
		sig3 := x3.Sign(m[:])

		X4 := x4.PublicKey()
		sig4 := x4.Sign(m[:])

		sig34 := AggregateSignatures([]Signature{sig3, sig4})
		require.False(t, sig34.IsZero())

		x3Opposite := opposite(t, x3)
		x4Opposite := opposite(t, x4)

		sig1 := x3Opposite.Sign(m[:])
		sig2 := x4Opposite.Sign(m[:])

		sig3Offsetted := AggregateSignatures([]Signature{sig3, sig2})
		require.False(t, sig3Offsetted.IsZero())
		sig4Offsetted := AggregateSignatures([]Signature{sig4, sig1})
		require.False(t, sig4Offsetted.IsZero())

		sig34Offsetted := AggregateSignatures([]Signature{sig3Offsetted, sig4Offsetted})
		require.True(t, sig34Offsetted.IsZero())

		// individual signatures are not valid anymore
		require.False(t, sig3Offsetted.Verify(X3, m[:], DefaultAssumeZeroValid))
		require.False(t, sig4Offsetted.Verify(X4, m[:], DefaultAssumeZeroValid))

		// aggregate signature is not valid with FastAggregateVerifyBatch
		require.False(t, FastAggregateVerifyBatch([]Signature{sig3Offsetted, sig4Offsetted}, []PublicKey{X3, X4}, m))
	})
}

func scalar(value byte) *blst.Scalar {
	scalar := new(blst.Scalar)
	b := make([]byte, scalarBytes)
	b[scalarBytes-1] = value
	scalar.FromBEndian(b)
	return scalar
}

func TestFastAggregateVerifyBatch(t *testing.T) {
	t.Run("Test point multiplication", func(t *testing.T) {
		// it is a building block of the FastAggregateVerifyBatch function
		key, err := RandKey()
		require.NoError(t, err)

		pointAffine := key.PublicKey().(*BlsPublicKey).p
		point := new(blst.P1)
		point.FromAffine(pointAffine)

		one := scalar(1)
		require.Equal(t, point.Serialize(), point.Mult(one).Serialize())

		two := scalar(2)
		require.Equal(t, point.Add(point).Serialize(), point.Mult(two).Serialize())

		five := scalar(5)
		require.Equal(t, point.Add(point).Add(point).Add(point).Add(point).Serialize(), point.Mult(five).Serialize())
	})
	t.Run("Doesn't modify original sigs and keys", func(t *testing.T) {
		key1, err := RandKey()
		require.NoError(t, err)
		key2, err := RandKey()
		require.NoError(t, err)

		pubkey1 := key1.PublicKey()
		pubkey2 := key2.PublicKey()

		pubkey1bytes := pubkey1.Marshal()
		pubkey2bytes := pubkey2.Marshal()

		m := common.Hash{0xca, 0xfe}

		sig1 := key1.Sign(m[:])
		sig2 := key2.Sign(m[:])

		sig1bytes := sig1.Marshal()
		sig2bytes := sig2.Marshal()

		require.True(t, FastAggregateVerifyBatch([]Signature{sig1, sig2}, []PublicKey{pubkey1, pubkey2}, m))

		require.Equal(t, pubkey1bytes, pubkey1.Marshal())
		require.Equal(t, pubkey2bytes, pubkey2.Marshal())
		require.Equal(t, sig1bytes, sig1.Marshal())
		require.Equal(t, sig2bytes, sig2.Marshal())
	})
	t.Run("works correctly", func(t *testing.T) {
		key1, err := RandKey()
		require.NoError(t, err)
		key2, err := RandKey()
		require.NoError(t, err)

		m := common.Hash{0xca, 0xfe}
		m2 := common.Hash{0xca, 0xff}

		var signatures []Signature
		var pubkeys []PublicKey

		require.False(t, FastAggregateVerifyBatch(signatures, pubkeys, m))

		signatures = append(signatures, key1.Sign(m[:]))
		pubkeys = append(pubkeys, key1.PublicKey())

		require.True(t, FastAggregateVerifyBatch(signatures, pubkeys, m))

		signatures = append(signatures, key2.Sign(m[:]))
		require.False(t, FastAggregateVerifyBatch(signatures, pubkeys, m))
		pubkeys = append(pubkeys, key2.PublicKey())

		require.True(t, FastAggregateVerifyBatch(signatures, pubkeys, m))

		signatures = append(signatures, key1.Sign(m[:]))
		pubkeys = append(pubkeys, key1.PublicKey())

		require.True(t, FastAggregateVerifyBatch(signatures, pubkeys, m))

		// order matters
		pub := pubkeys[0]
		pubkeys[0] = pubkeys[1]
		pubkeys[1] = pub
		require.False(t, FastAggregateVerifyBatch(signatures, pubkeys, m))

		// restore correct order
		pubkeys[1] = pubkeys[0]
		pubkeys[0] = pub
		require.True(t, FastAggregateVerifyBatch(signatures, pubkeys, m))

		// invalid sig
		signatures = append(signatures, key1.Sign(m2[:]))
		pubkeys = append(pubkeys, key1.PublicKey())
		require.False(t, FastAggregateVerifyBatch(signatures, pubkeys, m))
	})
	t.Run("deals well with 0 signatures and 0 public keys", func(t *testing.T) {
		key1, err := RandKey()
		require.NoError(t, err)
		key2, err := RandKey()
		require.NoError(t, err)

		m := common.Hash{0xca, 0xfe}

		var signatures []Signature
		var pubkeys []PublicKey

		signatures = append(signatures, key1.Sign(m[:]))
		pubkeys = append(pubkeys, key1.PublicKey())
		signatures = append(signatures, key2.Sign(m[:]))
		pubkeys = append(pubkeys, key2.PublicKey())

		require.True(t, FastAggregateVerifyBatch(signatures, pubkeys, m))

		signatures = append(signatures, key1.Sign(m[:]))
		pubkeys = append(pubkeys, key1.PublicKey())

		require.True(t, FastAggregateVerifyBatch(signatures, pubkeys, m))

		x1, x2 := generateZeroKeyPair(t)
		X1 := x1.PublicKey()
		X2 := x2.PublicKey()
		zeroKey, err := AggregatePublicKeys([]PublicKey{X1, X2})
		require.NoError(t, err)
		pubkeys = append(pubkeys, zeroKey)

		// individual sigs are valid, but aggregate is 0
		sig1 := x1.Sign(m[:])
		require.False(t, sig1.IsZero())
		sig2 := x2.Sign(m[:])
		require.False(t, sig2.IsZero())
		zeroSig := AggregateSignatures([]Signature{sig1, sig2})
		require.True(t, zeroSig.IsZero())
		signatures = append(signatures, zeroSig)

		// FastAggregate doesn't crash and returns correct result
		require.True(t, FastAggregateVerifyBatch(signatures, pubkeys, m))

		// add a wrong pubkey to the mix
		x3, err := RandKey()
		require.NoError(t, err)
		pubkeys = append(pubkeys, x3.PublicKey())
		require.False(t, FastAggregateVerifyBatch(signatures, pubkeys, m))

		// other edge cases
		require.False(t, FastAggregateVerifyBatch([]Signature{zeroSig}, []PublicKey{X1}, m))
		require.False(t, FastAggregateVerifyBatch([]Signature{sig1}, []PublicKey{zeroKey}, m))
		require.True(t, FastAggregateVerifyBatch([]Signature{zeroSig}, []PublicKey{zeroKey}, m))
		require.True(t, FastAggregateVerifyBatch([]Signature{sig1, sig2}, []PublicKey{X1, X2}, m))

		randomScalar := func() *blst.Scalar {
			var rbytes [scalarBytes]byte

			_, err := rand.Read(rbytes[:])
			if err != nil {
				panic("Cannot source randomness")
			}

			// Protect against the generator returning 0. Since the scalar value is
			// derived from a big endian byte slice, we take the last byte.
			rbytes[len(rbytes)-1] |= 0x01
			scalar := new(blst.Scalar)
			scalar.FromBEndian(rbytes[:])
			return scalar
		}

		// no problem doing scalar multiplication of zero sigs
		// result is always zero
		zeroSigNotAffine := new(blst.P2)
		zeroSigNotAffine.FromAffine(zeroSig.(*BlsSignature).s)
		multSig := zeroSigNotAffine.Mult(randomScalar())
		multSigWrap := &BlsSignature{multSig.ToAffine()}
		require.True(t, multSigWrap.IsZero())

		// no problem doing scalar multiplication of zero pubkeys
		// result is always zero
		zeroKeyNotAffine := new(blst.P1)
		zeroKeyNotAffine.FromAffine(zeroKey.(*BlsPublicKey).p)
		multPub := zeroKeyNotAffine.Mult(randomScalar())
		multPubWrap := &BlsPublicKey{multPub.ToAffine()}
		require.False(t, multPubWrap.p.KeyValidate())
	})
}

// make sure that AggregateVerify deals well with 0 sigs
func TestAggregateVerifyWithZero(t *testing.T) {
	m1 := common.Hash{0xca, 0xfe}
	m2 := common.Hash{0xca, 0xff}

	x1, x2 := generateZeroKeyPair(t)

	x3, err := RandKey()
	require.NoError(t, err)

	x4, err := RandKey()
	require.NoError(t, err)

	t.Run("removal of 0 sigs lead to correct verification", func(t *testing.T) {
		X34, err := AggregatePublicKeys([]PublicKey{x3.PublicKey(), x4.PublicKey()})
		require.NoError(t, err)

		sig1 := x1.Sign(m1[:])
		sig2 := x2.Sign(m1[:])
		sig3 := x3.Sign(m2[:])
		sig4 := x4.Sign(m2[:])

		sig1234 := AggregateSignatures([]Signature{sig1, sig2, sig3, sig4})

		// since sig12 is zero and the related public key is zero as well, we can just verify the aggregate signature by removing those terms
		require.True(t, sig1234.AggregateVerify([]PublicKey{X34}, [][32]byte{m2}))

		sig12 := AggregateSignatures([]Signature{sig1, sig2})
		require.True(t, sig12.IsZero())

		X12, err := AggregatePublicKeys([]PublicKey{x1.PublicKey(), x2.PublicKey()})
		require.NoError(t, err)

		require.False(t, X12.Validate())

		// zero sig with 0 pubkey is rejected
		require.False(t, sig12.AggregateVerify([]PublicKey{X12}, [][32]byte{m1}))

		// 0 signature gets flagged as invalid even if the public key is not 0
		require.False(t, sig12.AggregateVerify([]PublicKey{X34}, [][32]byte{m1}))
	})
}

func TestVerifyZero(t *testing.T) {
	x1, x2 := generateZeroKeyPair(t)

	X1 := x1.PublicKey()
	X2 := x2.PublicKey()
	// individual pubkeys are valid
	require.True(t, X1.(*BlsPublicKey).p.KeyValidate())
	require.True(t, X2.(*BlsPublicKey).p.KeyValidate())

	aggX, err := AggregatePublicKeys([]PublicKey{X1, X2})
	require.NoError(t, err)
	// aggregate pubkey is not valid, since it is the infinite pubkey
	require.False(t, aggX.(*BlsPublicKey).p.KeyValidate())

	msg := []byte("hello")
	sig1 := x1.Sign(msg)
	sig2 := x2.Sign(msg)

	// individual signatures are valid whether we check for zero or not
	require.True(t, sig1.Verify(X1, msg, true), "Signature did not verify")
	require.True(t, sig2.Verify(X2, msg, true), "Signature did not verify")
	require.True(t, sig1.Verify(X1, msg, false), "Signature did not verify")
	require.True(t, sig2.Verify(X2, msg, false), "Signature did not verify")

	zeroSig := AggregateSignatures([]Signature{sig1, sig2})
	require.True(t, zeroSig.IsZero())

	// zero sig and not zero pub key should always be invalid
	require.False(t, zeroSig.Verify(X1, msg, true))
	require.False(t, zeroSig.Verify(X2, msg, false))

	// non zero sig and zero pubkey should always be invalid
	require.False(t, sig1.Verify(aggX, msg, true))
	require.False(t, sig1.Verify(aggX, msg, false))

	// zero sig and zero key should be valid if explictly allowed
	require.True(t, zeroSig.Verify(aggX, msg, true))

	// invalid if explicitly not allowed
	require.False(t, zeroSig.Verify(aggX, msg, false))
}

func benchmarkFastAggregateVerifyBatch(b *testing.B, seed int64, n int) {
	// print out parameters
	b.Logf("seed: %d, n: %d\n", seed, n)

	// initialize deterministic randomness
	rand := mrand.New(mrand.NewSource(seed)) //nolint

	// generate msg
	var msg [32]byte
	rand.Read(msg[:])

	// generate signatures over constant msg
	var sigs []Signature
	var pubkeys []PublicKey
	for i := 0; i < n; i++ {
		var ikm [32]byte
		_, err := rand.Read(ikm[:])
		require.NoError(b, err)
		innerSk := blst.KeyGen(ikm[:])
		sk, err := SecretKeyFromBytes(innerSk.Serialize())
		require.NoError(b, err)
		sig := sk.Sign(msg[:])
		sigs = append(sigs, sig)
		pubkeys = append(pubkeys, sk.PublicKey())
	}

	// start the actual aggregation benchmarking
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		valid := FastAggregateVerifyBatch(sigs, pubkeys, msg)
		if !valid {
			b.Fatal("Failed to verify signature")
		}
	}
}

func BenchmarkFastAggregateVerifyBatch0_100(b *testing.B) {
	benchmarkFastAggregateVerifyBatch(b, 0, 100)
}
func BenchmarkFastAggregateVerifyBatch1_1000(b *testing.B) {
	benchmarkFastAggregateVerifyBatch(b, 1, 1000)
}
func BenchmarkFastAggregateVerifyBatch2_200(b *testing.B) {
	benchmarkFastAggregateVerifyBatch(b, 2, 200)
}
