package blst

import (
	"fmt"
	mrand "math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFastFindInvalid(t *testing.T) {
	// initialize deterministic randomness
	rand := mrand.New(mrand.NewSource(0)) //nolint

	// create message to be signed
	var msg [32]byte
	rand.Read(msg[:])

	// generate foreign key and foreign signature
	// these artifacts are not invalid per-se, but they are not included into the "valid set"
	foreignSk, err := RandKey()
	require.NoError(t, err)
	foreignSig := foreignSk.Sign(msg[:])

	// generate signatures over msg
	n := 30
	var sigs []Signature
	var pks []PublicKey
	for i := 0; i < n; i++ {
		sk, err := RandKey()
		require.NoError(t, err)

		sig := sk.Sign(msg[:])

		pks = append(pks, sk.PublicKey())
		sigs = append(sigs, sig)
	}

	t.Run("empty signature set", func(t *testing.T) {
		invalid := FindInvalid(nil, nil, [32]byte{})
		require.Empty(t, invalid)
	})

	t.Run("single invalid", func(t *testing.T) {
		invalid := FindInvalid([]Signature{foreignSig}, []PublicKey{pks[0]}, msg)
		require.Len(t, invalid, 1)
		require.Equal(t, uint(0), invalid[0])
	})

	t.Run("last invalid", func(t *testing.T) {
		invalid := FindInvalid([]Signature{sigs[0], foreignSig}, []PublicKey{pks[0], pks[1]}, msg)
		require.Len(t, invalid, 1)
	})

	t.Run("single valid", func(t *testing.T) {
		invalid := FindInvalid([]Signature{sigs[0]}, []PublicKey{pks[0]}, msg)
		require.Empty(t, invalid)
	})

	t.Run("large valid", func(t *testing.T) {
		invalid := FindInvalid(sigs, pks, msg)
		require.Empty(t, invalid)
	})
	t.Run("large invalid", func(t *testing.T) {

		// turns signatures at specified indexes into invalid signatures
		makeInvalid := func(valid []Signature, n ...int) []Signature {
			sigs := make([]Signature, len(valid)) //nolint
			copy(sigs, valid)

			for _, i := range n {
				sigs[i] = foreignSk.Sign(msg[:])
			}
			return sigs
		}

		t.Run("single invalid", func(t *testing.T) {
			sigs := makeInvalid(sigs, 5) //nolint

			invalid := FindInvalid(sigs, pks, msg)
			require.Len(t, invalid, 1)
			require.Equal(t, invalid[0], uint(5))
		})

		t.Run("multiple invalid", func(t *testing.T) {
			sigs := makeInvalid(sigs, 0, 3, 5, 7, n-1) //nolint

			invalid := FindInvalid(sigs, pks, msg)
			require.Len(t, invalid, 5)

			fmt.Println(invalid)

			require.Equal(t, invalid[0], uint(0))
			require.Equal(t, invalid[1], uint(3))
			require.Equal(t, invalid[2], uint(5))
			require.Equal(t, invalid[3], uint(7))
			require.Equal(t, invalid[4], uint(n-1))
		})
	})
}

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
