package blst

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	blst "github.com/supranational/blst/bindings/go"
	"testing"
)

func TestSignVerify(t *testing.T) {
	priv, err := RandKey()
	require.NoError(t, err)
	pub := priv.PublicKey()
	msg := []byte("hello")
	sig := priv.Sign(msg)
	require.Equal(t, true, sig.Verify(pub, msg), "Signature did not verify")
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

	fmt.Println("helloAggSig1: ", helloAggSig1.Hex())
	fmt.Println("helloAggSig2: ", helloAggSig2.Hex())
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
	t.Skip("skip time cost full unit test")
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
	t.Skip("skip time cost full unit test")
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
	t.Skip("skip time cost full unit test")
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
func Test20MinEpoch21ValidatorsOnSameMsgFastAggregateVerify(t *testing.T) {
	t.Skip("skip time cost full unit test")
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
		aggSig := AggregateSignatures(sigs)
		require.Equal(t, true, aggSig.FastAggregateVerify(pubkeys, msg), "Signature did not verify")
	}
}

// if the msg is distinct, then the order of public key does not impact the aggregation verification.
func TestFastAggregateVerify(t *testing.T) {
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
	aggSig := AggregateSignatures(sigs)

	require.Equal(t, true, aggSig.FastAggregateVerify(pubkeys, msg), "Signature did not verify")
}

func TestVerifyCompressed(t *testing.T) {
	priv, err := RandKey()
	require.NoError(t, err)
	pub := priv.PublicKey()
	msg := []byte("hello")
	sig := priv.Sign(msg)
	require.Equal(t, true, sig.Verify(pub, msg), "Non compressed signature did not verify")
	require.Equal(t, true, VerifyCompressed(sig.Marshal(), pub.Marshal(), msg), "Compressed signatures and pubkeys did not verify")
}

func TestMultipleSignatureVerification(t *testing.T) {
	pubkeys := make([]PublicKey, 0, 100)
	sigs := make([][]byte, 0, 100)
	var msgs [][32]byte
	for i := 0; i < 100; i++ {
		msg := [32]byte{'h', 'e', 'l', 'l', 'o', byte(i)}
		priv, err := RandKey()
		require.NoError(t, err)
		pub := priv.PublicKey()
		sig := priv.Sign(msg[:]).Marshal()
		pubkeys = append(pubkeys, pub)
		sigs = append(sigs, sig)
		msgs = append(msgs, msg)
	}
	verify, err := VerifyMultipleSignatures(sigs, msgs, pubkeys)
	require.NoError(t, err, "Signature did not verify")
	require.Equal(t, true, verify, "Signature did not verify")
}

// with same key to sign different msgs, and do the signature aggregation.
// in such case, a same validator can form a single aggregation of signatures for the entire Epoch.
func TestMultipleSignatureByDistinctKeyVerification(t *testing.T) {
	pubkeys := make([]PublicKey, 0, 100)
	sigs := make([][]byte, 0, 100)
	priv, err := RandKey()
	var msgs [][32]byte
	for i := 0; i < 100; i++ {
		msg := [32]byte{'h', 'e', 'l', 'l', 'o', byte(i)}
		require.NoError(t, err)
		pub := priv.PublicKey()
		sig := priv.Sign(msg[:]).Marshal()
		pubkeys = append(pubkeys, pub)
		sigs = append(sigs, sig)
		msgs = append(msgs, msg)
	}
	verify, err := VerifyMultipleSignatures(sigs, msgs, pubkeys)
	require.NoError(t, err, "Signature did not verify")
	require.Equal(t, true, verify, "Signature did not verify")
}

func TestFastAggregateVerify_ReturnsFalseOnEmptyPubKeyList(t *testing.T) {
	var pubkeys []PublicKey
	msg := [32]byte{'h', 'e', 'l', 'l', 'o'}

	aggSig := newAggregateSignature()
	require.Equal(t, false, aggSig.FastAggregateVerify(pubkeys, msg), "Expected FastAggregateVerify to return false with empty input ")
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
	signatureB, ok := signatureA.Copy().(*BlsSignature)
	require.Equal(t, true, ok)
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

// newAggregateSignature creates a blank aggregate signature.
func newAggregateSignature() Signature {
	sig := blst.HashToG2([]byte{'m', 'o', 'c', 'k'}, generalDST).ToAffine()
	return &BlsSignature{s: sig}
}
