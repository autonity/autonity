package blst

import (
	mrand "math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindInvalid(t *testing.T) {

	// initialize deterministic randomness
	rand := mrand.New(mrand.NewSource(0)) //nolint

	n := 30
	failedN := 5

	invalidSk, err := RandKey()
	require.NoError(t, err)

	// create messages
	var msgs [][32]byte
	var msg [32]byte
	for i := 0; i < n; i++ {
		rand.Read(msg[:])
		msgs = append(msgs, msg)
	}

	invalidSig := invalidSk.Sign(msgs[failedN][:])

	// generate signatures over msgs
	sigs := make([]Signature, 0, n)

	var pks []PublicKey
	for i := 0; i < n; i++ {
		sk, err := RandKey()
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
		invalid, err := FindInvalidSignatures([]Signature{invalidSig}, []PublicKey{pks[failedN]}, [][32]byte{msgs[failedN]})
		require.NoError(t, err)
		require.Len(t, invalid, 1)
		require.Equal(t, uint(0), invalid[0])
	})

	t.Run("single valid", func(t *testing.T) {
		invalid, err := FindInvalidSignatures([]Signature{sigs[0]}, []PublicKey{pks[0]}, [][32]byte{msgs[0]})
		require.NoError(t, err)
		require.Empty(t, invalid)
	})

	t.Run("last invalid", func(t *testing.T) {
		invalid, err := FindInvalidSignatures([]Signature{sigs[0], invalidSig}, []PublicKey{pks[0], pks[1]}, [][32]byte{msgs[0], msgs[1]})
		require.NoError(t, err)
		require.Len(t, invalid, 1)
	})

	t.Run("large valid", func(t *testing.T) {
		invalid, err := FindInvalidSignatures(sigs, pks, msgs)
		require.NoError(t, err)
		require.Empty(t, invalid)
	})

	t.Run("large invalid", func(t *testing.T) {

		makeInvalid := func(n ...int) []Signature {

			newSigs := make([]Signature, len(sigs))

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

			require.Contains(t, invalid, uint(0))
			require.Contains(t, invalid, uint(3))
			require.Contains(t, invalid, uint(5))
			require.Contains(t, invalid, uint(7))
			require.Contains(t, invalid, uint(n-1))
		})
	})
}

func TestFastFindInvalid(t *testing.T) {

	// initialize deterministic randomness
	rand := mrand.New(mrand.NewSource(0)) //nolint

	n := 30
	failedN := 5

	invalidSk, err := RandKey()
	require.NoError(t, err)

	// create messages
	var msg [32]byte
	rand.Read(msg[:])

	invalidSig := invalidSk.Sign(msg[:])

	// generate signatures over msgs
	sigs := make([]Signature, 0, n)

	var pks []PublicKey
	for i := 0; i < n; i++ {
		sk, err := RandKey()
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
		invalid, err := FindFastInvalidSignatures([]Signature{invalidSig}, []PublicKey{pks[failedN]}, msg)
		require.NoError(t, err)
		require.Len(t, invalid, 1)
		require.Equal(t, uint(0), invalid[0])
	})

	t.Run("last invalid", func(t *testing.T) {
		invalid, err := FindFastInvalidSignatures([]Signature{sigs[0], invalidSig}, []PublicKey{pks[0], pks[1]}, msg)
		require.NoError(t, err)
		require.Len(t, invalid, 1)
	})

	t.Run("single valid", func(t *testing.T) {
		invalid, err := FindFastInvalidSignatures([]Signature{sigs[0]}, []PublicKey{pks[0]}, msg)
		require.NoError(t, err)
		require.Empty(t, invalid)
	})

	t.Run("large valid", func(t *testing.T) {
		invalid, err := FindFastInvalidSignatures(sigs, pks, msg)
		require.NoError(t, err)
		require.Empty(t, invalid)
	})

	t.Run("large invalid", func(t *testing.T) {

		makeInvalid := func(n ...int) []Signature {

			newSigs := make([]Signature, len(sigs))

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

			require.Contains(t, invalid, uint(0))
			require.Contains(t, invalid, uint(3))
			require.Contains(t, invalid, uint(5))
			require.Contains(t, invalid, uint(7))
			require.Contains(t, invalid, uint(n-1))
		})
	})
}
