package bls

import (
	"bytes"
	"crypto/ecdsa"
	"github.com/autonity/autonity/crypto"
	"math/big"
	"testing"

	"github.com/autonity/autonity/crypto/secp256k1"
	"github.com/stretchr/testify/require"

	"crypto/rand"
	mrand "math/rand"

	"github.com/autonity/autonity/crypto/bls/common"
)

func TestVerifyOwnershipProof(t *testing.T) {
	privKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	address := crypto.PubkeyToAddress(privKey.PublicKey)

	validatorKey, err := SecretKeyFromECDSAKey(privKey)
	require.NoError(t, err)

	proof, err := GenerateValidatorKeyProof(validatorKey, address.Bytes())
	require.NoError(t, err)

	sig, err := SignatureFromBytes(proof)
	require.NoError(t, err)

	err = ValidateValidatorKeyProof(validatorKey.PublicKey(), sig, address.Bytes())
	require.NoError(t, err)
}

func TestDisallowZeroSecretKeys(t *testing.T) {
	t.Run("blst", func(t *testing.T) {
		// Blst does a zero check on the key during deserialization.
		_, err := SecretKeyFromBytes(common.ZeroSecretKey[:])
		require.Equal(t, common.ErrSecretUnmarshal, err)
	})
}

func TestDisallowZeroPublicKeys(t *testing.T) {
	t.Run("blst", func(t *testing.T) {
		_, err := PublicKeyFromBytes(common.InfinitePublicKey[:])
		require.Equal(t, common.ErrInfinitePubKey, err)
	})
}

func TestDisallowZeroPublicKeys_AggregatePubkeys(t *testing.T) {
	t.Run("blst", func(t *testing.T) {
		_, err := AggregatePublicKeys([][]byte{common.InfinitePublicKey[:], common.InfinitePublicKey[:]})
		require.Equal(t, common.ErrInfinitePubKey, err)
	})
}

func TestValidateSecretKeyString(t *testing.T) {
	t.Run("blst", func(t *testing.T) {
		zeroNum := new(big.Int).SetUint64(0)
		_, err := SecretKeyFromBigNum(zeroNum.String())
		require.NotNil(t, err)

		randBytes := make([]byte, 40)
		n, err := rand.Read(randBytes)
		require.NoError(t, err)
		require.Equal(t, n, len(randBytes))
		rBigNum := new(big.Int).SetBytes(randBytes)

		// Expect larger than expected key size to fail.
		_, err = SecretKeyFromBigNum(rBigNum.String())
		require.NotNil(t, err)

		key, err := RandKey()
		require.NoError(t, err)
		rBigNum = new(big.Int).SetBytes(key.Marshal())

		// Expect correct size to pass.
		_, err = SecretKeyFromBigNum(rBigNum.String())
		require.NoError(t, err)
	})
}

func TestReuseECDSAKeyForBLS(t *testing.T) {
	// new ecdsa key, it is a 32 bytes random number.
	ecdsaKey, err := ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
	require.NoError(t, err)

	blsPrivateKey, err := SecretKeyFromECDSAKey(ecdsaKey)
	require.NoError(t, err)

	for i := 0; i < 100; i++ {
		// use ecdsa key, to elicit a deterministic private key of bls.
		skNewGenerated, err := SecretKeyFromECDSAKey(ecdsaKey)
		require.NoError(t, err)

		// check the deterministic property.
		require.Equal(t, true, bytes.Equal(blsPrivateKey.Marshal(), skNewGenerated.Marshal()))

		// check basic features of bls signature signing and verification.
		pk := skNewGenerated.PublicKey()

		require.Equal(t, 98, len(pk.Hex()))

		msg := []byte("hello world!")
		sig := skNewGenerated.Sign(msg)

		require.Equal(t, 194, len(sig.Hex()))
		require.Equal(t, true, sig.Verify(pk, msg))
	}
}

// benchmarks

// benchmark bls public key aggregation
func benchmarkAggregatePK(seed int64, n int, b *testing.B) {
	// print out parameters
	b.Logf("seed: %d, n: %d\n", seed, n)

	// initialize deterministic randomness
	rand := mrand.New(mrand.NewSource(seed)) //nolint

	// setup public keys
	var pks [][]byte
	for i := 0; i < n; i++ {
		ecdsaKey, err := ecdsa.GenerateKey(secp256k1.S256(), rand)
		if err != nil {
			b.Fatal("Failed to generate random ecdsa key. error: ", err)
		}
		sk, err := SecretKeyFromECDSAKey(ecdsaKey)
		if err != nil {
			b.Fatal("Failed to generate random bls key: ", err)
		}
		pk := sk.PublicKey().Marshal()
		pks = append(pks, pk)
	}

	// start the actual aggregation benchmarking
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := AggregatePublicKeys(pks); err != nil {
			b.Fatal("failed pks aggregation: ", err)
		}
	}

}

func BenchmarkAggregatePK0_100(b *testing.B) { benchmarkAggregatePK(0, 100, b) }
func BenchmarkAggregatePK1_100(b *testing.B) { benchmarkAggregatePK(1, 100, b) }
func BenchmarkAggregatePK2_100(b *testing.B) { benchmarkAggregatePK(2, 100, b) }
func BenchmarkAggregatePK0_200(b *testing.B) { benchmarkAggregatePK(0, 200, b) }
func BenchmarkAggregatePK1_200(b *testing.B) { benchmarkAggregatePK(1, 200, b) }
func BenchmarkAggregatePK2_200(b *testing.B) { benchmarkAggregatePK(2, 200, b) }
func BenchmarkAggregatePK0_300(b *testing.B) { benchmarkAggregatePK(0, 300, b) }
func BenchmarkAggregatePK1_300(b *testing.B) { benchmarkAggregatePK(1, 300, b) }
func BenchmarkAggregatePK2_300(b *testing.B) { benchmarkAggregatePK(2, 300, b) }

// benchmark bls signature aggregation
func benchmarkAggregateSigs(seed int64, n int, b *testing.B) {
	// print out parameters
	b.Logf("seed: %d, n: %d\n", seed, n)

	// initialize deterministic randomness
	rand := mrand.New(mrand.NewSource(seed)) //nolint

	// generate msg
	var msg [32]byte
	rand.Read(msg[:])

	// generate signatures over constant msg
	var sigs []common.BLSSignature
	for i := 0; i < n; i++ {
		ecdsaKey, err := ecdsa.GenerateKey(secp256k1.S256(), rand)
		if err != nil {
			b.Fatal("Failed to generate random ecdsa key. error: ", err)
		}
		sk, err := SecretKeyFromECDSAKey(ecdsaKey)
		if err != nil {
			b.Fatal("Failed to generate random bls key: ", err)
		}
		sig := sk.Sign(msg[:])
		sigs = append(sigs, sig)
	}

	// start the actual aggregation benchmarking
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = AggregateSignatures(sigs)
	}

}

func BenchmarkAggregateSigs0_3(b *testing.B)   { benchmarkAggregateSigs(0, 3, b) }
func BenchmarkAggregateSigs0_100(b *testing.B) { benchmarkAggregateSigs(0, 100, b) }
func BenchmarkAggregateSigs1_100(b *testing.B) { benchmarkAggregateSigs(1, 100, b) }
func BenchmarkAggregateSigs2_100(b *testing.B) { benchmarkAggregateSigs(2, 100, b) }
func BenchmarkAggregateSigs0_200(b *testing.B) { benchmarkAggregateSigs(0, 200, b) }
func BenchmarkAggregateSigs1_200(b *testing.B) { benchmarkAggregateSigs(1, 200, b) }
func BenchmarkAggregateSigs2_200(b *testing.B) { benchmarkAggregateSigs(2, 200, b) }
func BenchmarkAggregateSigs0_300(b *testing.B) { benchmarkAggregateSigs(0, 300, b) }
func BenchmarkAggregateSigs1_300(b *testing.B) { benchmarkAggregateSigs(1, 300, b) }
func BenchmarkAggregateSigs2_300(b *testing.B) { benchmarkAggregateSigs(2, 300, b) }

// benchmark simple signature verification
func BenchmarkSigVerify(b *testing.B) {
	// initialize deterministic randomness
	rand := mrand.New(mrand.NewSource(0)) //nolint

	// generate msg
	var msg [32]byte
	rand.Read(msg[:])

	// generate signature over constant msg
	ecdsaKey, err := ecdsa.GenerateKey(secp256k1.S256(), rand)
	if err != nil {
		b.Fatal("Failed to generate random ecdsa key. error: ", err)
	}
	sk, err := SecretKeyFromECDSAKey(ecdsaKey)
	if err != nil {
		b.Fatal("Failed to generate random bls key: ", err)
	}
	sig := sk.Sign(msg[:])

	// start the actual verification benchmarking
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res := sig.Verify(sk.PublicKey(), msg[:])
		if !res {
			b.Fatal("failed signature verification")
		}
	}

}

// benchmark aggregated signature verification on 1 message
func BenchmarkSigVerifyAgg(b *testing.B) {
	// initialize deterministic randomness
	rand := mrand.New(mrand.NewSource(0)) //nolint

	// generate msg
	var msg [32]byte
	rand.Read(msg[:])

	// generate signatures over constant msg
	var sigs []common.BLSSignature
	var pks [][]byte
	for i := 0; i < 100; i++ {
		ecdsaKey, err := ecdsa.GenerateKey(secp256k1.S256(), rand)
		if err != nil {
			b.Fatal("Failed to generate random ecdsa key. error: ", err)
		}
		sk, err := SecretKeyFromECDSAKey(ecdsaKey)
		if err != nil {
			b.Fatal("Failed to generate random bls key: ", err)
		}
		pk := sk.PublicKey().Marshal()
		pks = append(pks, pk)
		sig := sk.Sign(msg[:])
		sigs = append(sigs, sig)
	}

	sig := AggregateSignatures(sigs)
	aggPk, err := AggregatePublicKeys(pks)
	if err != nil {
		b.Fatal("failed pks aggregation: ", err)
	}

	// start the actual aggregation benchmarking
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res := sig.Verify(aggPk, msg[:])
		if !res {
			b.Fatal("failed signature verification")
		}
	}
}

// benchmark aggregated signature verification on 3 messages (propose, prevote, precommit)
func BenchmarkAggregateVerify3(b *testing.B) {
	// initialize deterministic randomness
	rand := mrand.New(mrand.NewSource(0)) //nolint

	// create step messages
	var propose [32]byte
	rand.Read(propose[:])
	var prevote [32]byte
	rand.Read(prevote[:])
	var precommit [32]byte
	rand.Read(precommit[:])

	// generate signature over propose
	ecdsaKey, err := ecdsa.GenerateKey(secp256k1.S256(), rand)
	if err != nil {
		b.Fatal("Failed to generate random ecdsa key. error: ", err)
	}
	proposeSk, err := SecretKeyFromECDSAKey(ecdsaKey)
	if err != nil {
		b.Fatal("Failed to generate random bls key: ", err)
	}
	sigPropose := proposeSk.Sign(propose[:])

	// generate signatures over votes
	var sigsPrevote []common.BLSSignature
	var sigsPrecommit []common.BLSSignature
	var pks [][]byte
	for i := 0; i < 100; i++ {
		ecdsaKey, err := ecdsa.GenerateKey(secp256k1.S256(), rand)
		if err != nil {
			b.Fatal("Failed to generate random ecdsa key. error: ", err)
		}
		sk, err := SecretKeyFromECDSAKey(ecdsaKey)
		if err != nil {
			b.Fatal("Failed to generate random bls key: ", err)
		}
		pk := sk.PublicKey().Marshal()
		pks = append(pks, pk)
		sigPrevote := sk.Sign(prevote[:])
		sigsPrevote = append(sigsPrevote, sigPrevote)
		sigPrecommit := sk.Sign(precommit[:])
		sigsPrecommit = append(sigsPrecommit, sigPrecommit)
	}

	aggSigPrevote := AggregateSignatures(sigsPrevote)
	aggSigPrecommit := AggregateSignatures(sigsPrecommit)
	aggSig := AggregateSignatures([]common.BLSSignature{sigPropose, aggSigPrevote, aggSigPrecommit})
	aggPk, err := AggregatePublicKeys(pks)
	if err != nil {
		b.Fatal("failed pks aggregation: ", err)
	}

	// start the actual aggregation benchmarking
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res := aggSig.AggregateVerify([]common.BLSPublicKey{proposeSk.PublicKey(), aggPk, aggPk}, [][32]byte{propose, prevote, precommit})
		if !res {
			b.Fatal("failed signature verification")
		}
	}

}

// benchmark aggregated signature verification on 15 messages
func BenchmarkAggregateVerify15(b *testing.B) {
	// initialize deterministic randomness
	rand := mrand.New(mrand.NewSource(0)) //nolint

	n := 15

	// create messages
	var msgs [][32]byte
	var msg [32]byte
	for i := 0; i < n; i++ {
		rand.Read(msg[:])
		msgs = append(msgs, msg)
	}

	// generate signatures over msgs
	sigs := make([][]common.BLSSignature, n)
	var pks [][]byte
	for i := 0; i < 100; i++ {
		ecdsaKey, err := ecdsa.GenerateKey(secp256k1.S256(), rand)
		if err != nil {
			b.Fatal("Failed to generate random ecdsa key. error: ", err)
		}
		sk, err := SecretKeyFromECDSAKey(ecdsaKey)
		if err != nil {
			b.Fatal("Failed to generate random bls key: ", err)
		}
		pk := sk.PublicKey().Marshal()
		pks = append(pks, pk)
		for j := 0; j < len(msgs); j++ {
			sig := sk.Sign(msgs[j][:])
			sigs[j] = append(sigs[j], sig)
		}
	}

	aggPk, err := AggregatePublicKeys(pks)
	if err != nil {
		b.Fatal("failed pks aggregation: ", err)
	}
	var aggSigs []common.BLSSignature
	for i := 0; i < len(sigs); i++ {
		aggSig := AggregateSignatures(sigs[i])
		aggSigs = append(aggSigs, aggSig)
	}
	aggSigTotal := AggregateSignatures(aggSigs)

	var pksTotal []common.BLSPublicKey
	for i := 0; i < n; i++ {
		pksTotal = append(pksTotal, aggPk)
	}

	// start the actual aggregation benchmarking
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res := aggSigTotal.AggregateVerify(pksTotal, msgs)
		if !res {
			b.Fatal("failed signature verification")
		}
	}

}
