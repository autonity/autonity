package bls

import (
	"bytes"
	"crypto/ecdsa"
	"math/big"
	mrand "math/rand"
	"testing"

	"github.com/autonity/autonity/crypto/secp256k1"
	"github.com/stretchr/testify/require"

	"crypto/rand"
	"github.com/autonity/autonity/crypto/bls/common"
)

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

func TestFindInvalid(t *testing.T) {

	// initialize deterministic randomness
	rand := mrand.New(mrand.NewSource(0)) //nolint

	n := 30
	failedN := 5

	invalidEcdsaKey, err := ecdsa.GenerateKey(secp256k1.S256(), rand)
	require.NoError(t, err)
	invalidSk, err := SecretKeyFromECDSAKey(invalidEcdsaKey)

	// create messages
	var msgs [][32]byte
	var msg [32]byte
	for i := 0; i < n; i++ {
		rand.Read(msg[:])
		msgs = append(msgs, msg)
	}

	invalidSig := invalidSk.Sign(msgs[failedN][:])

	// generate signatures over msgs
	sigs := make([]common.BLSSignature, 0, n)

	var pks []common.BLSPublicKey
	for i := 0; i < n; i++ {
		ecdsaKey, err := ecdsa.GenerateKey(secp256k1.S256(), rand)
		require.NoError(t, err)
		sk, err := SecretKeyFromECDSAKey(ecdsaKey)
		require.NoError(t, err)
		pk := sk.PublicKey()
		pks = append(pks, pk)

		sig := sk.Sign(msgs[i][:])
		sigs = append(sigs, sig)
	}

	t.Run("empty", func(t *testing.T) {
		invalid, err := FindInvalidSignatures(nil, nil, nil)
		require.NoError(t, err)
		require.Empty(t, invalid)
	})

	t.Run("single invalid", func(t *testing.T) {
		invalid, err := FindInvalidSignatures([]common.BLSSignature{invalidSig}, []common.BLSPublicKey{pks[failedN]}, [][32]byte{msgs[failedN]})
		require.NoError(t, err)
		require.Len(t, invalid, 1)
		require.Equal(t, 0, invalid[0])
	})

	t.Run("single valid", func(t *testing.T) {
		invalid, err := FindInvalidSignatures([]common.BLSSignature{sigs[0]}, []common.BLSPublicKey{pks[0]}, [][32]byte{msgs[0]})
		require.NoError(t, err)
		require.Empty(t, invalid)
	})

	t.Run("last invalid", func(t *testing.T) {
		invalid, err := FindInvalidSignatures([]common.BLSSignature{sigs[0], invalidSig}, []common.BLSPublicKey{pks[0], pks[1]}, [][32]byte{msgs[0], msgs[1]})
		require.NoError(t, err)
		require.Len(t, invalid, 1)
	})

	t.Run("large valid", func(t *testing.T) {
		invalid, err := FindInvalidSignatures(sigs, pks, msgs)
		require.NoError(t, err)
		require.Empty(t, invalid)
	})

	t.Run("large invalid", func(t *testing.T) {

		makeInvalid := func(n ...int) []common.BLSSignature {

			newSigs := make([]common.BLSSignature, len(sigs))

			copy(newSigs, sigs)

			for _, i := range n {
				invalidSig := invalidSk.Sign(msgs[i][:])

				newSigs[i] = invalidSig
			}
			return newSigs
		}

		t.Run("single", func(t *testing.T) {
			sigs := makeInvalid(failedN)

			invalid, err := FindInvalidSignatures(sigs, pks, msgs)
			require.NoError(t, err)
			require.Len(t, invalid, 1)
		})

		t.Run("multiple", func(t *testing.T) {
			sigs := makeInvalid(0, 3, 5, 7, n-1)

			invalid, err := FindInvalidSignatures(sigs, pks, msgs)
			require.NoError(t, err)
			require.Len(t, invalid, 5)

			require.Contains(t, invalid, 0)
			require.Contains(t, invalid, 3)
			require.Contains(t, invalid, 5)
			require.Contains(t, invalid, 7)
			require.Contains(t, invalid, n-1)
		})
	})
}

func TestFastFindInvalid(t *testing.T) {

	// initialize deterministic randomness
	rand := mrand.New(mrand.NewSource(0)) //nolint

	n := 30
	failedN := 5

	invalidEcdsaKey, err := ecdsa.GenerateKey(secp256k1.S256(), rand)
	require.NoError(t, err)
	invalidSk, err := SecretKeyFromECDSAKey(invalidEcdsaKey)

	// create messages
	var msg [32]byte
	rand.Read(msg[:])

	invalidSig := invalidSk.Sign(msg[:])

	// generate signatures over msgs
	sigs := make([]common.BLSSignature, 0, n)

	var pks []common.BLSPublicKey
	for i := 0; i < n; i++ {
		ecdsaKey, err := ecdsa.GenerateKey(secp256k1.S256(), rand)
		require.NoError(t, err)
		sk, err := SecretKeyFromECDSAKey(ecdsaKey)
		require.NoError(t, err)
		pk := sk.PublicKey()
		pks = append(pks, pk)

		sig := sk.Sign(msg[:])
		sigs = append(sigs, sig)
	}

	t.Run("empty", func(t *testing.T) {
		invalid, err := FindFastInvalidSignatures(nil, nil, [32]byte{})
		require.NoError(t, err)
		require.Empty(t, invalid)
	})

	t.Run("single invalid", func(t *testing.T) {
		invalid, err := FindFastInvalidSignatures([]common.BLSSignature{invalidSig}, []common.BLSPublicKey{pks[failedN]}, msg)
		require.NoError(t, err)
		require.Len(t, invalid, 1)
		require.Equal(t, 0, invalid[0])
	})

	t.Run("last invalid", func(t *testing.T) {
		invalid, err := FindFastInvalidSignatures([]common.BLSSignature{sigs[0], invalidSig}, []common.BLSPublicKey{pks[0], pks[1]}, msg)
		require.NoError(t, err)
		require.Len(t, invalid, 1)
	})

	t.Run("single valid", func(t *testing.T) {
		invalid, err := FindFastInvalidSignatures([]common.BLSSignature{sigs[0]}, []common.BLSPublicKey{pks[0]}, msg)
		require.NoError(t, err)
		require.Empty(t, invalid)
	})

	t.Run("large valid", func(t *testing.T) {
		invalid, err := FindFastInvalidSignatures(sigs, pks, msg)
		require.NoError(t, err)
		require.Empty(t, invalid)
	})

	t.Run("large invalid", func(t *testing.T) {

		makeInvalid := func(n ...int) []common.BLSSignature {

			newSigs := make([]common.BLSSignature, len(sigs))

			copy(newSigs, sigs)

			for _, i := range n {
				invalidSig := invalidSk.Sign(msg[:])

				newSigs[i] = invalidSig
			}
			return newSigs
		}

		t.Run("single", func(t *testing.T) {
			sigs := makeInvalid(failedN)

			invalid, err := FindFastInvalidSignatures(sigs, pks, msg)
			require.NoError(t, err)
			require.Len(t, invalid, 1)
		})

		t.Run("multiple", func(t *testing.T) {
			sigs := makeInvalid(0, 3, 5, 7, n-1)

			invalid, err := FindFastInvalidSignatures(sigs, pks, msg)
			require.NoError(t, err)
			require.Len(t, invalid, 5)

			require.Contains(t, invalid, 0)
			require.Contains(t, invalid, 3)
			require.Contains(t, invalid, 5)
			require.Contains(t, invalid, 7)
			require.Contains(t, invalid, n-1)
		})
	})
}
