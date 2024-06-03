package types

import (
	"bytes"
	"errors"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/hexutil"
	"github.com/autonity/autonity/rlp"
)

var (
	committee = Committee{{
		Address:           common.HexToAddress("0x76a685e4bf8cbcd25d7d3b6c342f64a30b503380"),
		ConsensusKeyBytes: hexutil.MustDecode("0x951f3f7ab473eb0d00eaaa569ba1a0be2877b794e29e0cbf504b7f00cb879a824b0b913397e0071a87cebaae2740002b"),
		VotingPower:       hexutil.MustDecodeBig("0x3029"),
	}, {
		Address:           common.HexToAddress("0xc44276975a6c2d12e62e18d814b507c38fc3646f"),
		ConsensusKeyBytes: hexutil.MustDecode("0x8bddc21fca7f3a920064729547605c73e55c17e20917eddc8788b97990c0d7e9420e51a97ea400fb58a5c28fa63984eb"),
		VotingPower:       hexutil.MustDecodeBig("0x3139"),
	}, {
		Address:           common.HexToAddress("0x1a72cb9d17c9e7acad03b4d3505f160e3782f2d5"),
		ConsensusKeyBytes: hexutil.MustDecode("0x9679c8ebd47d18b93acd90cd380debdcfdb140f38eca207c61463a47be85398ec3082a66f7f30635c11470f5c8e5cf6b"),
		VotingPower:       hexutil.MustDecodeBig("0x3056"),
	}, {
		Address:           common.HexToAddress("0xbaa58a01e5ca81dc288e2c46a8a467776bdb81c6"),
		ConsensusKeyBytes: hexutil.MustDecode("0xa460c204c407b6272f7731b0d15daca8f2564cf7ace301769e3b42de2482fc3bf8116dd13c0545e806441d074d02dcc2"),
		VotingPower:       hexutil.MustDecodeBig("0x39"),
	}, {
		Address:           common.HexToAddress("0xca72cb9d17c9e7acad03b4d3505f160e3782f2d5"),
		ConsensusKeyBytes: hexutil.MustDecode("0x8b36875624f4cc6adc4e39dfd059fb0490288542c0b63dd16ffa10e2d5b47d437d6cf8325ec2d96d8255e46d561d20d0"),
		VotingPower:       hexutil.MustDecodeBig("0x7033"),
	}}
)

func expectPanic(t *testing.T) {
	if r := recover(); r == nil {
		t.Fatalf("The code did not panic")
	} else {
		t.Log(r)
	}
}

func TestValidatorBitmap(t *testing.T) {
	t.Run("Simple bitmap with 4 validators (1 byte)", func(t *testing.T) {
		// 00 00 00 00
		n := 4
		bitmap := NewValidatorBitmap(n)
		require.True(t, bitmap.Valid(n))

		for i := 0; i < n; i++ {
			require.Equal(t, byte(0), bitmap.Get(i))
		}

		// 10 00 00 00
		bitmap.Set(0, 2)
		require.Equal(t, byte(2), bitmap.Get(0))

		// 11 00 00 00
		bitmap.Set(0, 3)
		require.Equal(t, byte(3), bitmap.Get(0))

		// 11 01 00 00
		bitmap.Set(1, 1)
		require.Equal(t, byte(1), bitmap.Get(1))

		// 11 01 00 00
		bitmap.Set(2, 0)
		require.Equal(t, byte(0), bitmap.Get(2))

		// 11 01 10 00
		bitmap.Set(2, 2)
		require.Equal(t, byte(2), bitmap.Get(2))

		// 11 01 10 11
		bitmap.Set(3, 3)
		require.Equal(t, byte(3), bitmap.Get(3))

		// 11 00 10 11
		bitmap.Set(1, 0)
		require.Equal(t, byte(0), bitmap.Get(1))

		// 11 00 10 00
		bitmap.Set(3, 0)
		require.Equal(t, byte(0), bitmap.Get(3))
	})
	t.Run("bitmap with 6 validators (2 byte)", func(t *testing.T) {
		// 00 00 00 00 | 00 00 00 00
		n := 6
		bitmap := NewValidatorBitmap(n)
		require.True(t, bitmap.Valid(n))

		for i := 0; i < n; i++ {
			require.Equal(t, byte(0), bitmap.Get(i))
		}

		// 10 00 00 00 | 00 00 00 00
		bitmap.Set(0, 2)
		require.Equal(t, byte(2), bitmap.Get(0))

		// 11 00 00 00 | 00 00 00 00
		bitmap.Set(0, 3)
		require.Equal(t, byte(3), bitmap.Get(0))

		// 11 00 00 00 | 01 00 00 00
		bitmap.Set(4, 1)
		require.Equal(t, byte(1), bitmap.Get(4))

		// 11 00 00 00 | 01 11 00 00
		bitmap.Set(5, 3)
		require.Equal(t, byte(3), bitmap.Get(5))
	})
	t.Run("Setting out of bound bit should panic", func(t *testing.T) {
		defer expectPanic(t)
		bitmap := NewValidatorBitmap(0)
		bitmap.Set(1, 1)
	})
	t.Run("Setting value >= 3 should panic", func(t *testing.T) {
		defer expectPanic(t)
		bitmap := NewValidatorBitmap(10)
		bitmap.Set(1, 4)
	})
	t.Run("Getting out of bound bit should panic", func(t *testing.T) {
		defer expectPanic(t)
		bitmap := NewValidatorBitmap(0)
		bitmap.Get(1)
	})
	t.Run("Valid method", func(t *testing.T) {
		n := 4
		bitmap := NewValidatorBitmap(n)
		require.True(t, bitmap.Valid(n))
		bitmap = append(bitmap, []byte{0xca, 0xfe}...)
		require.False(t, bitmap.Valid(n))
		bitmap = make([]byte, 0)
		require.False(t, bitmap.Valid(n))
	})
}

func TestSerialization(t *testing.T) {
	err := committee.Enrich()
	require.NoError(t, err)

	t.Run("auxiliary data structures are not serialized on the wire", func(t *testing.T) {
		// increment some sender info
		s := NewSigners(len(committee))
		// 11 01 01 10
		s.Increment(&committee[0])
		s.Increment(&committee[0])
		s.Increment(&committee[0])
		s.Increment(&committee[0])
		s.Increment(&committee[1])
		s.Increment(&committee[2])
		s.Increment(&committee[3])
		s.Increment(&committee[3])

		// bits, coefficients + auxiliary data structures should be set
		require.Equal(t, s.Bits[0], byte(0xd6))
		require.Equal(t, len(s.Coefficients), 1)
		require.Equal(t, s.Coefficients[0], uint16(4))
		// one validator is missing from the aggregate
		require.Equal(t, len(s.Powers()), len(committee)-1)
		require.Equal(t, s.Len(), len(committee)-1)
		expectedVotingPower := committee[0].VotingPower.Uint64() + committee[1].VotingPower.Uint64() + committee[2].VotingPower.Uint64() + committee[3].VotingPower.Uint64()
		require.Equal(t, s.Power().Uint64(), expectedVotingPower)
		require.Equal(t, s.CommitteeSize(), len(committee))
		require.True(t, s.validated)
		require.True(t, s.powerAssigned)

		// encode and decode
		payload, err := rlp.EncodeToBytes(s)
		require.NoError(t, err)
		decoded := &Signers{}
		err = rlp.Decode(bytes.NewBuffer(payload), decoded)
		require.NoError(t, err)

		// only bits and coefficients should be set, auxiliary data structures should be empty
		require.Equal(t, decoded.Bits[0], byte(0xd6))
		require.Equal(t, len(decoded.Coefficients), 1)
		require.Equal(t, decoded.Coefficients[0], uint16(4))

		require.False(t, decoded.validated)
		require.False(t, decoded.powerAssigned)

		// accessing auxiliary data structure should panic
		defer expectPanic(t)
		require.Nil(t, decoded.Powers())
	})
}

func TestSigners(t *testing.T) {
	err := committee.Enrich()
	require.NoError(t, err)

	totalPower := new(big.Int)
	for _, member := range committee {
		totalPower.Add(totalPower, member.VotingPower)
	}

	t.Run("Increment should update the bitmap and the auxiliary maps correctly", func(t *testing.T) {
		// +10 to avoid hitting the panic in `increment` related to the max allowed coefficient
		s := NewSigners(len(committee) + 10)
		require.Equal(t, s.CommitteeSize(), len(committee)+10)

		require.Equal(t, s.Bits.Get(0), byte(0))

		s.Increment(&committee[0])
		require.Equal(t, s.Bits.Get(0), byte(1))
		require.Equal(t, len(s.Powers()), 1)
		require.Equal(t, s.Power().Uint64(), committee[0].VotingPower.Uint64())
		require.Equal(t, s.Powers()[0].Uint64(), committee[0].VotingPower.Uint64())
		require.Equal(t, s.Len(), 1)

		s.Increment(&committee[0])
		require.Equal(t, s.Bits.Get(0), byte(2))
		s.Increment(&committee[0])
		require.Equal(t, s.Bits.Get(0), byte(3))
		require.Equal(t, s.Coefficients[0], uint16(3))
		s.Increment(&committee[0])
		require.Equal(t, s.Bits.Get(0), byte(3))
		require.Equal(t, s.Coefficients[0], uint16(4))
		s.Increment(&committee[0])
		require.Equal(t, s.Bits.Get(0), byte(3))
		require.Equal(t, s.Coefficients[0], uint16(5))

		require.Equal(t, len(s.Powers()), 1)
		require.Equal(t, s.Power().Uint64(), committee[0].VotingPower.Uint64())
		require.Equal(t, s.Powers()[0].Uint64(), committee[0].VotingPower.Uint64())
		require.Equal(t, s.Len(), 1)

		s.Increment(&committee[1])
		require.Equal(t, s.Bits.Get(1), byte(1))

		require.Equal(t, len(s.Powers()), 2)
		require.Equal(t, s.Power().Uint64(), committee[0].VotingPower.Uint64()+committee[1].VotingPower.Uint64())
		require.Equal(t, s.Powers()[0].Uint64(), committee[0].VotingPower.Uint64())
		require.Equal(t, s.Len(), 2)

		s.Increment(&committee[2])
		require.Equal(t, s.Bits.Get(2), byte(1))
		s.Increment(&committee[2])
		require.Equal(t, s.Bits.Get(2), byte(2))
		s.Increment(&committee[2])
		require.Equal(t, s.Bits.Get(2), byte(3))
		require.Equal(t, s.Coefficients[1], uint16(3))
		s.Increment(&committee[2])
		require.Equal(t, s.Bits.Get(2), byte(3))
		require.Equal(t, s.Coefficients[1], uint16(4))
		s.Increment(&committee[2])
		require.Equal(t, s.Bits.Get(2), byte(3))
		require.Equal(t, s.Coefficients[1], uint16(5))

		s.Increment(&committee[2])
		s.Increment(&committee[2])
		s.Increment(&committee[2])
		require.Equal(t, s.Coefficients[1], uint16(8))

		s.Increment(&committee[0])
		require.Equal(t, s.Coefficients[0], uint16(6))

		s.Increment(&committee[1])
		require.Equal(t, s.Bits.Get(1), uint8(2))
		s.Increment(&committee[1])
		require.Equal(t, s.Bits.Get(1), uint8(3))
		require.Equal(t, s.Coefficients[1], uint16(3))
		s.Increment(&committee[1])
		require.Equal(t, s.Bits.Get(1), uint8(3))
		require.Equal(t, s.Coefficients[1], uint16(4))
		s.Increment(&committee[1])
		require.Equal(t, s.Bits.Get(1), uint8(3))
		require.Equal(t, s.Coefficients[1], uint16(5))

		s.Increment(&committee[1])
		require.Equal(t, s.Bits.Get(1), uint8(3))
		require.Equal(t, s.Coefficients[2], uint16(8))
	})
	t.Run("Merge correctly merges two senders info", func(t *testing.T) {
		// +10 to avoid hitting the panic in `increment` related to the max allowed coefficient
		s1 := NewSigners(len(committee) + 10)
		s1.Increment(&committee[0])
		s1.Increment(&committee[1])
		s2 := NewSigners(len(committee) + 10)
		s2.Increment(&committee[2])
		s2.Increment(&committee[3])
		s2.Increment(&committee[4])

		s1.Merge(s2)

		require.Equal(t, s1.CommitteeSize(), len(committee)+10)
		require.Equal(t, len(s1.Powers()), len(committee))
		require.Equal(t, s1.Len(), len(committee))
		for i, member := range committee {
			require.Equal(t, s1.Bits.Get(i), uint8(1))
			require.Equal(t, s1.Powers()[i], member.VotingPower)
		}
		require.Equal(t, s1.Power().Uint64(), totalPower.Uint64())

		s3 := NewSigners(len(committee) + 10)
		s3.Increment(&committee[0])
		s3.Increment(&committee[1])
		s3.Increment(&committee[2])
		s3.Increment(&committee[3])
		s3.Increment(&committee[4])

		s1.Merge(s3)

		require.Equal(t, s1.CommitteeSize(), len(committee)+10)
		require.Equal(t, len(s1.Powers()), len(committee))
		require.Equal(t, s1.Len(), len(committee))
		for i, member := range committee {
			require.Equal(t, s1.Bits.Get(i), uint8(2))
			require.Equal(t, s1.Powers()[i], member.VotingPower)
		}
		require.Equal(t, s1.Power().Uint64(), totalPower.Uint64())

		s4 := NewSigners(len(committee) + 10)
		s4.Increment(&committee[0])
		s4.Increment(&committee[0])
		s4.Increment(&committee[0])
		s4.Increment(&committee[2])
		s4.Increment(&committee[2])
		s4.Increment(&committee[2])
		s4.Increment(&committee[2])

		s1.Merge(s4)

		require.Equal(t, s1.Bits.Get(0), uint8(3))
		require.Equal(t, s1.Coefficients[0], uint16(5))
		require.Equal(t, s1.Bits.Get(2), uint8(3))
		require.Equal(t, s1.Coefficients[1], uint16(6))

		s5 := NewSigners(len(committee) + 10)
		s5.Increment(&committee[1])

		s1.Merge(s5)
		require.Equal(t, s1.Bits.Get(0), uint8(3))
		require.Equal(t, s1.Bits.Get(1), uint8(3))
		require.Equal(t, s1.Bits.Get(2), uint8(3))
		require.Equal(t, s1.Coefficients[0], uint16(5))
		require.Equal(t, s1.Coefficients[1], uint16(3))
		require.Equal(t, s1.Coefficients[2], uint16(6))
	})
	t.Run("Power returns the aggregated power of the senders", func(t *testing.T) {
		s := NewSigners(len(committee))
		s.Increment(&committee[0])
		s.Increment(&committee[1])

		require.Equal(t, s.Power(), new(big.Int).Add(committee[0].VotingPower, committee[1].VotingPower))

		s.Increment(&committee[2])
		s.Increment(&committee[3])
		s.Increment(&committee[4])

		require.Equal(t, s.Power(), totalPower)

		// duplicated power shouldn't be counted
		s.Increment(&committee[2])
		s.Increment(&committee[3])

		require.Equal(t, s.Power(), totalPower)
	})
	t.Run("Flatten returns the indexes of the senders (repeated in case of multiple contribution to the signature)", func(t *testing.T) {
		// +10 to avoid hitting the panic related to maximum allowed coefficient in `increment`
		s := NewSigners(len(committee) + 10)
		s.Increment(&committee[0])
		s.Increment(&committee[1])

		require.Equal(t, s.Flatten(), []int{0, 1})

		s.Increment(&committee[3])
		s.Increment(&committee[3])
		s.Increment(&committee[0])
		s.Increment(&committee[0])
		s.Increment(&committee[0])
		s.Increment(&committee[1])
		s.Increment(&committee[2])

		require.Equal(t, s.Flatten(), []int{0, 0, 0, 0, 1, 1, 2, 3, 3})

		s.Increment(&committee[1])
		s.Increment(&committee[1])
		s.Increment(&committee[1])
		s.Increment(&committee[1])
		s.Increment(&committee[1])

		require.Equal(t, s.Flatten(), []int{0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 2, 3, 3})
	})
	t.Run("FlattenUniq returns the indexes of the senders (de-duplicated)", func(t *testing.T) {
		s := NewSigners(len(committee))
		s.Increment(&committee[0])
		s.Increment(&committee[1])

		require.Equal(t, s.FlattenUniq(), []int{0, 1})

		s.Increment(&committee[3])
		s.Increment(&committee[3])
		s.Increment(&committee[0])
		s.Increment(&committee[0])
		s.Increment(&committee[0])
		s.Increment(&committee[1])
		s.Increment(&committee[2])

		require.Equal(t, s.FlattenUniq(), []int{0, 1, 2, 3})
	})
	t.Run("Contains returns expected result", func(t *testing.T) {
		s := NewSigners(len(committee))

		require.False(t, s.Contains(0))
		require.False(t, s.Contains(1))
		require.False(t, s.Contains(2))

		s.Increment(&committee[0])
		s.Increment(&committee[2])

		require.True(t, s.Contains(0))
		require.False(t, s.Contains(1))
		require.True(t, s.Contains(2))
	})
	t.Run("IsComplex returns expected result", func(t *testing.T) {
		s := NewSigners(len(committee))

		require.False(t, s.IsComplex())

		s.Increment(&committee[0])
		s.Increment(&committee[1])
		s.Increment(&committee[2])

		require.False(t, s.IsComplex())

		s.Increment(&committee[0])

		require.True(t, s.IsComplex())

		s.Increment(&committee[0])
		s.Increment(&committee[1])
		s.Increment(&committee[2])
		s.Increment(&committee[0])
		s.Increment(&committee[1])
		s.Increment(&committee[2])

		require.True(t, s.IsComplex())

	})
	t.Run("AddsInformation, RespectBoundaries and CanMergeSimple return expected results", func(t *testing.T) {
		s1 := NewSigners(len(committee))
		s2 := NewSigners(len(committee))

		require.False(t, s1.AddsInformation(s2))
		require.True(t, s1.CanMergeSimple(s2))
		require.True(t, s1.RespectsBoundaries(s2))

		s1.increment(0)

		require.False(t, s1.AddsInformation(s2))
		require.True(t, s1.CanMergeSimple(s2))
		require.True(t, s1.RespectsBoundaries(s2))

		s2.increment(0)

		require.False(t, s1.AddsInformation(s2))
		require.False(t, s1.CanMergeSimple(s2))
		require.True(t, s1.RespectsBoundaries(s2))

		s2.increment(1)

		require.True(t, s1.AddsInformation(s2))
		require.False(t, s1.CanMergeSimple(s2))
		require.True(t, s1.RespectsBoundaries(s2))

		s2.increment(2)

		require.True(t, s1.AddsInformation(s2))
		require.False(t, s1.CanMergeSimple(s2))
		require.True(t, s1.RespectsBoundaries(s2))

		s1.increment(0)
		s1.increment(0)
		s2.increment(0)

		require.True(t, s1.AddsInformation(s2))
		require.False(t, s1.CanMergeSimple(s2))
		require.True(t, s1.RespectsBoundaries(s2))

		s2.increment(0)

		require.True(t, s1.AddsInformation(s2))
		require.False(t, s1.CanMergeSimple(s2))
		require.False(t, s1.RespectsBoundaries(s2))

		s1.increment(1)
		s1.increment(2)

		require.False(t, s1.AddsInformation(s2))
		require.False(t, s1.CanMergeSimple(s2))
		require.False(t, s1.RespectsBoundaries(s2))

	})
}

func TestValidation(t *testing.T) {
	csize := 10

	nilSigner := &Signers{
		Bits:         nil,
		Coefficients: nil,
	}
	require.True(t, errors.Is(nilSigner.Validate(csize), ErrNilSigners))

	wrongSizeSigner := NewSigners(csize)
	wrongSizeSigner.Coefficients = make([]uint16, csize+10)
	require.True(t, errors.Is(wrongSizeSigner.Validate(csize), ErrWrongSizeSigners))

	s := NewSigners(csize)
	require.True(t, errors.Is(s.Validate(csize), ErrEmptySigners))

	// A
	s.increment(0)
	t.Logf(s.String())
	require.Nil(t, s.Validate(csize))

	// A + B
	s.increment(1)
	t.Logf(s.String())
	require.Nil(t, s.Validate(csize))

	// A + B + C
	s.increment(2)
	t.Logf(s.String())
	require.Nil(t, s.Validate(csize))

	// A + B + C + D
	s.increment(3)
	t.Logf(s.String())
	require.Nil(t, s.Validate(csize))

	// 2A + B + C + D
	s.increment(0)
	t.Logf(s.String())
	require.Nil(t, s.Validate(csize))

	// 3A + B + C + D
	s.increment(0)
	t.Logf(s.String())
	require.Nil(t, s.Validate(csize))

	// 4A + 2B + C + D
	s.increment(0)
	s.increment(1)
	t.Logf(s.String())
	require.Nil(t, s.Validate(csize))

	s.Coefficients = append(s.Coefficients, []uint16{0xca, 0xfe}...)
	require.True(t, errors.Is(s.Validate(csize), ErrWrongCoefficientLen))

	s = NewSigners(csize)
	s.increment(0)
	s.increment(0)
	require.True(t, errors.Is(s.Validate(csize), ErrInvalidSingleSig))

	s.increment(0)
	s.increment(0)
	s.increment(1)
	s.Coefficients = []uint16{uint16(csize + 10)}
	require.True(t, errors.Is(s.Validate(csize), ErrInvalidCoefficient))

	s = NewSigners(csize)
	s.increment(0)
	s.Validate(csize)
	require.True(t, s.validated)
	require.Equal(t, csize, s.CommitteeSize())
	require.Equal(t, 1, s.Len())
}
