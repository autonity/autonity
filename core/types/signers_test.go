package types

import (
	"bytes"
	"errors"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/bitutil"
	"github.com/autonity/autonity/common/hexutil"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/rlp"
)

var (
	committee = &Committee{
		Members: []CommitteeMember{
			{
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
			},
		},
	}
)

func expectPanic(t *testing.T) {
	if r := recover(); r == nil {
		t.Fatalf("The code did not panic")
	} else {
		t.Log(r)
	}
}

func signersAggregation[T uint16 | uint32](
	t *testing.T,
	committee []*CommitteeMember,
	coeffs []int,
	signature blst.Signature,
	msgInput common.Hash,
) {
	signers := NewSigners[T](len(committee))
	keys := make([]blst.PublicKey, 0, len(committee))
	for i, c := range coeffs {
		for c > 0 {
			signers.AddMember(committee[i])
			c--
		}
		keys = append(keys, committee[i].ConsensusKey)
	}

	publicKey := signers.AggregatePublicKey(keys)
	require.True(t, publicKey.Validate())
	require.True(t, signature.Verify(publicKey, msgInput[:], blst.DefaultAssumeZeroValid))
}

func TestPublicKeyAggregation(t *testing.T) {
	csize := 10 // at least 4
	members := make([]*CommitteeMember, 0, csize)
	privateKeyMap := make(map[blst.SecretKey]struct{})
	coeffs := make([]int, 0, csize)

	signatures := make([]blst.Signature, 0)

	msg := common.BigToHash(big.NewInt(1234345677))

	for i := 0; i < csize; i++ {
		var privateKey blst.SecretKey
		var err error
		for {
			privateKey, err = blst.RandKey()
			require.NoError(t, err)
			if _, ok := privateKeyMap[privateKey]; ok {
				continue
			}
			privateKeyMap[privateKey] = struct{}{}
			break
		}

		members = append(members, &CommitteeMember{
			Index:        uint64(i),
			VotingPower:  big.NewInt(1),
			Address:      common.BigToAddress(big.NewInt(int64(i + 1))),
			ConsensusKey: privateKey.PublicKey(),
		})
		coeffs = append(coeffs, i+1)

		signature := privateKey.Sign(msg[:])

		for c := 0; c < coeffs[i]; c++ {
			signatures = append(signatures, signature)
		}
	}

	aggregateSign := blst.AggregateSignatures(signatures)
	signersAggregation[uint16](t, members, coeffs, aggregateSign, msg)
	signersAggregation[uint32](t, members, coeffs, aggregateSign, msg)
}

func TestValidatorBitmap(t *testing.T) {
	t.Run("Simple bitmap with 8 validators (1 byte)", func(t *testing.T) {
		// 00000000
		n := 8
		bitmap := bitutil.NewBitmap(n)
		require.True(t, bitmap.Valid(n))

		for i := 0; i < n; i++ {
			require.Equal(t, false, bitmap.HasItem(i))
		}

		for i := 0; i < n; i++ {
			bitmap.AddItem(i)
			require.Equal(t, true, bitmap.HasItem(i))
		}
	})
	t.Run("bitmap with 12 validators (2 byte)", func(t *testing.T) {
		// 00000000 | 00000000
		n := 12
		bitmap := bitutil.NewBitmap(n)
		require.True(t, bitmap.Valid(n))

		for i := 0; i < n; i++ {
			require.Equal(t, false, bitmap.HasItem(i))
		}

		// 00000001 | 00000000
		bitmap.AddItem(0)
		require.Equal(t, true, bitmap.HasItem(0))

		// 00000001 | 00000001
		bitmap.AddItem(8)
		require.Equal(t, true, bitmap.HasItem(8))

		// 00000001 | 00000011
		bitmap.AddItem(9)
		require.Equal(t, true, bitmap.HasItem(9))
	})
	t.Run("Setting out of bound bit should panic", func(t *testing.T) {
		defer expectPanic(t)
		bitmap := bitutil.NewBitmap(0)
		bitmap.AddItem(1)
	})
	t.Run("Getting out of bound bit should panic", func(t *testing.T) {
		defer expectPanic(t)
		bitmap := bitutil.NewBitmap(0)
		bitmap.HasItem(1)
	})
	t.Run("Valid method", func(t *testing.T) {
		n := 8
		bitmap := bitutil.NewBitmap(n)
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
		s := NewVoteSigners(committee.Len())
		// 00001111
		s.AddMember(&committee.Members[0])
		s.AddMember(&committee.Members[1])
		s.AddMember(&committee.Members[1])
		s.AddMember(&committee.Members[2])
		s.AddMember(&committee.Members[2])
		s.AddMember(&committee.Members[2])
		s.AddMember(&committee.Members[3])

		// bits, coefficients + auxiliary data structures should be set
		require.Equal(t, s.Bits[0], byte(0xf))
		require.Equal(t, len(s.Coefficients), 4)
		require.Equal(t, s.Coefficients[0], uint16(1))
		require.Equal(t, s.Coefficients[1], uint16(2))
		require.Equal(t, s.Coefficients[2], uint16(3))
		require.Equal(t, s.Coefficients[3], uint16(1))

		// one validator is missing from the aggregate
		require.Equal(t, len(s.Powers()), committee.Len()-1)
		require.Equal(t, s.Len(), committee.Len()-1)
		expectedVotingPower := committee.Members[0].VotingPower.Uint64() + committee.Members[1].VotingPower.Uint64() +
			committee.Members[2].VotingPower.Uint64() + committee.Members[3].VotingPower.Uint64()
		require.Equal(t, s.Power().Uint64(), expectedVotingPower)
		require.Equal(t, s.CommitteeSize(), committee.Len())
		require.True(t, s.validated)
		require.True(t, s.powerAssigned)

		// encode and decode
		payload, err := rlp.EncodeToBytes(s)
		require.NoError(t, err)
		decoded := &VoteSigners{&SignersBase[uint16]{}}
		err = rlp.Decode(bytes.NewBuffer(payload), decoded)
		require.NoError(t, err)

		// only bits and coefficients should be set, auxiliary data structures should be empty
		require.Equal(t, decoded.Bits[0], byte(0xf))
		require.Equal(t, len(decoded.Coefficients), 4)
		require.Equal(t, s.Coefficients[0], uint16(1))
		require.Equal(t, s.Coefficients[1], uint16(2))
		require.Equal(t, s.Coefficients[2], uint16(3))
		require.Equal(t, s.Coefficients[3], uint16(1))

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
	for _, member := range committee.Members {
		totalPower.Add(totalPower, member.VotingPower)
	}

	t.Run("Increment should update the bitmap and the auxiliary maps correctly", func(t *testing.T) {
		// +10 to avoid hitting the panic in `increment` related to the max allowed coefficient
		s := NewVoteSigners(committee.Len() + 10)
		require.Equal(t, s.CommitteeSize(), committee.Len()+10)

		require.Equal(t, s.Bits.HasItem(0), false)

		s.AddMember(&committee.Members[0])
		require.Equal(t, s.Bits.HasItem(0), true)
		require.Equal(t, s.Coefficients[0], uint16(1))
		require.Equal(t, len(s.Powers()), 1)
		require.Equal(t, s.Power().Uint64(), committee.Members[0].VotingPower.Uint64())
		require.Equal(t, s.Powers()[0].Uint64(), committee.Members[0].VotingPower.Uint64())
		require.Equal(t, s.Len(), 1)

		s.AddMember(&committee.Members[0])
		require.Equal(t, s.Bits.HasItem(0), true)
		require.Equal(t, s.Coefficients[0], uint16(2))
		s.AddMember(&committee.Members[0])
		require.Equal(t, s.Bits.HasItem(0), true)
		require.Equal(t, s.Coefficients[0], uint16(3))

		require.Equal(t, len(s.Powers()), 1)
		require.Equal(t, s.Power().Uint64(), committee.Members[0].VotingPower.Uint64())
		require.Equal(t, s.Powers()[0].Uint64(), committee.Members[0].VotingPower.Uint64())
		require.Equal(t, s.Len(), 1)

		s.AddMember(&committee.Members[1])
		require.Equal(t, s.Bits.HasItem(1), true)
		require.Equal(t, s.Coefficients[1], uint16(1))

		require.Equal(t, len(s.Powers()), 2)
		require.Equal(t, s.Power().Uint64(), committee.Members[0].VotingPower.Uint64()+committee.Members[1].VotingPower.Uint64())
		require.Equal(t, s.Powers()[0].Uint64(), committee.Members[0].VotingPower.Uint64())
		require.Equal(t, s.Len(), 2)

		s.AddMember(&committee.Members[2])
		require.Equal(t, s.Bits.HasItem(2), true)
		require.Equal(t, s.Coefficients[2], uint16(1))
		s.AddMember(&committee.Members[2])
		require.Equal(t, s.Bits.HasItem(2), true)
		require.Equal(t, s.Coefficients[2], uint16(2))
		s.AddMember(&committee.Members[2])
		require.Equal(t, s.Bits.HasItem(2), true)
		require.Equal(t, s.Coefficients[2], uint16(3))

		s.AddMember(&committee.Members[2])
		s.AddMember(&committee.Members[2])
		s.AddMember(&committee.Members[2])
		require.Equal(t, s.Coefficients[2], uint16(6))

		s.AddMember(&committee.Members[0])
		require.Equal(t, s.Coefficients[0], uint16(4))

		s.AddMember(&committee.Members[1])
		require.Equal(t, s.Bits.HasItem(1), true)
		require.Equal(t, s.Coefficients[1], uint16(2))
		s.AddMember(&committee.Members[1])
		require.Equal(t, s.Bits.HasItem(1), true)
		require.Equal(t, s.Coefficients[1], uint16(3))
		s.AddMember(&committee.Members[1])
		require.Equal(t, s.Bits.HasItem(1), true)
		require.Equal(t, s.Coefficients[1], uint16(4))

		s.AddMember(&committee.Members[1])
		require.Equal(t, s.Bits.HasItem(1), true)
		require.Equal(t, s.Coefficients[1], uint16(5))
	})
	t.Run("Merge correctly merges two senders info", func(t *testing.T) {
		// +10 to avoid hitting the panic in `increment` related to the max allowed coefficient
		s1 := NewVoteSigners(committee.Len() + 10)
		s1.AddMember(&committee.Members[0])
		s1.AddMember(&committee.Members[1])
		s2 := NewVoteSigners(committee.Len() + 10)
		s2.AddMember(&committee.Members[2])
		s2.AddMember(&committee.Members[3])
		s2.AddMember(&committee.Members[4])

		s1.Merge(s2.SignersBase)

		require.Equal(t, s1.CommitteeSize(), committee.Len()+10)
		require.Equal(t, len(s1.Powers()), committee.Len())
		require.Equal(t, s1.Len(), committee.Len())
		for i, member := range committee.Members {
			require.Equal(t, s1.Bits.HasItem(i), true)
			require.Equal(t, s1.Coefficients[i], uint16(1))
			require.Equal(t, s1.Powers()[i], member.VotingPower)
		}
		require.Equal(t, s1.Power().Uint64(), totalPower.Uint64())

		s3 := NewVoteSigners(committee.Len() + 10)
		s3.AddMember(&committee.Members[0])
		s3.AddMember(&committee.Members[1])
		s3.AddMember(&committee.Members[2])
		s3.AddMember(&committee.Members[3])
		s3.AddMember(&committee.Members[4])

		s1.Merge(s3.SignersBase)

		require.Equal(t, s1.CommitteeSize(), committee.Len()+10)
		require.Equal(t, len(s1.Powers()), committee.Len())
		require.Equal(t, s1.Len(), committee.Len())
		for i, member := range committee.Members {
			require.Equal(t, s1.Bits.HasItem(i), true)
			require.Equal(t, s1.Coefficients[i], uint16(2))
			require.Equal(t, s1.Powers()[i], member.VotingPower)
		}
		require.Equal(t, s1.Power().Uint64(), totalPower.Uint64())

		s4 := NewVoteSigners(committee.Len() + 10)
		s4.AddMember(&committee.Members[0])
		s4.AddMember(&committee.Members[0])
		s4.AddMember(&committee.Members[0])
		s4.AddMember(&committee.Members[2])
		s4.AddMember(&committee.Members[2])
		s4.AddMember(&committee.Members[2])
		s4.AddMember(&committee.Members[2])

		s1.Merge(s4.SignersBase)

		require.Equal(t, s1.Bits.HasItem(0), true)
		require.Equal(t, s1.Coefficients[0], uint16(5))
		require.Equal(t, s1.Bits.HasItem(2), true)
		require.Equal(t, s1.Coefficients[2], uint16(6))

		s5 := NewVoteSigners(committee.Len() + 10)
		s5.AddMember(&committee.Members[1])

		s1.Merge(s5.SignersBase)
		require.Equal(t, s1.Bits.HasItem(0), true)
		require.Equal(t, s1.Bits.HasItem(1), true)
		require.Equal(t, s1.Bits.HasItem(2), true)
		require.Equal(t, s1.Coefficients[0], uint16(5))
		require.Equal(t, s1.Coefficients[1], uint16(3))
		require.Equal(t, s1.Coefficients[2], uint16(6))
	})
	t.Run("Power returns the aggregated power of the senders", func(t *testing.T) {
		s := NewVoteSigners(committee.Len())
		s.AddMember(&committee.Members[0])
		s.AddMember(&committee.Members[1])

		require.Equal(t, s.Power(), new(big.Int).Add(committee.Members[0].VotingPower, committee.Members[1].VotingPower))

		s.AddMember(&committee.Members[2])
		s.AddMember(&committee.Members[3])
		s.AddMember(&committee.Members[4])

		require.Equal(t, s.Power(), totalPower)

		// duplicated power shouldn't be counted
		s.AddMember(&committee.Members[2])
		s.AddMember(&committee.Members[3])

		require.Equal(t, s.Power(), totalPower)
	})
	t.Run("FlattenUniq returns the indexes of the senders (de-duplicated)", func(t *testing.T) {
		s := NewVoteSigners(committee.Len())
		s.AddMember(&committee.Members[0])
		s.AddMember(&committee.Members[1])

		require.Equal(t, s.FlattenUniq(), []int{0, 1})

		s.AddMember(&committee.Members[3])
		s.AddMember(&committee.Members[3])
		s.AddMember(&committee.Members[0])
		s.AddMember(&committee.Members[0])
		s.AddMember(&committee.Members[0])
		s.AddMember(&committee.Members[1])
		s.AddMember(&committee.Members[2])

		require.Equal(t, s.FlattenUniq(), []int{0, 1, 2, 3})
	})
	t.Run("Contains returns expected result", func(t *testing.T) {
		s := NewVoteSigners(committee.Len())

		require.False(t, s.Contains(0))
		require.False(t, s.Contains(1))
		require.False(t, s.Contains(2))

		s.AddMember(&committee.Members[0])
		s.AddMember(&committee.Members[2])

		require.True(t, s.Contains(0))
		require.False(t, s.Contains(1))
		require.True(t, s.Contains(2))
	})
	t.Run("AddsInformation, RespectBoundaries return expected results", func(t *testing.T) {
		s1 := NewVoteSigners(committee.Len())
		s2 := NewVoteSigners(committee.Len())

		require.False(t, s1.AddsInformation(s2.SignersBase))

		s1.increment(0, 1, common.Big0) //0

		require.False(t, s1.AddsInformation(s2.SignersBase))

		s2.increment(0, 1, common.Big0) //0

		require.False(t, s1.AddsInformation(s2.SignersBase))

		s2.increment(1, 1, common.Big0) //0,1

		require.True(t, s1.AddsInformation(s2.SignersBase))

		s2.increment(2, 1, common.Big0) //0,1,2

		require.True(t, s1.AddsInformation(s2.SignersBase))

		s1.increment(0, 1, common.Big0)
		s1.increment(0, 1, common.Big0) //0,0,0
		s2.increment(0, 1, common.Big0) //0,0,1,2

		require.True(t, s1.AddsInformation(s2.SignersBase))

		s2.increment(0, 1, common.Big0) //0,0,0,1,2

		require.True(t, s1.AddsInformation(s2.SignersBase))

		s1.increment(1, 1, common.Big0)
		s1.increment(2, 1, common.Big0)

		require.False(t, s1.AddsInformation(s2.SignersBase))

	})
}

func TestValidation(t *testing.T) {
	csize := 10

	nilSigner := &SignersBase[uint16]{
		Bits:         nil,
		Coefficients: nil,
	}
	require.True(t, errors.Is(nilSigner.Validate(csize), ErrNilSigners))

	wrongSizeSigner := NewVoteSigners(csize)
	wrongSizeSigner.Coefficients = make([]uint16, csize+10)
	require.True(t, errors.Is(wrongSizeSigner.Validate(csize), ErrWrongSizeSigners))

	s := NewVoteSigners(csize)
	require.True(t, errors.Is(s.Validate(csize), ErrEmptySigners))

	// A
	s.increment(0, 1, common.Big0)
	t.Log(s.String())
	require.Nil(t, s.Validate(csize))

	// A + B
	s.increment(1, 1, common.Big0)
	t.Log(s.String())
	require.Nil(t, s.Validate(csize))

	// A + B + C
	s.increment(2, 1, common.Big0)
	t.Log(s.String())
	require.Nil(t, s.Validate(csize))

	// A + B + C + D
	s.increment(3, 1, common.Big0)
	t.Log(s.String())
	require.Nil(t, s.Validate(csize))

	// 2A + B + C + D
	s.increment(0, 1, common.Big0)
	t.Log(s.String())
	require.Nil(t, s.Validate(csize))

	// 3A + B + C + D
	s.increment(0, 1, common.Big0)
	t.Log(s.String())
	require.Nil(t, s.Validate(csize))

	// 4A + 2B + C + D
	s.increment(0, 1, common.Big0)
	s.increment(1, 1, common.Big0)
	t.Log(s.String())
	require.Nil(t, s.Validate(csize))

	s.Coefficients = append(s.Coefficients, []uint16{0xca, 0xfe}...)
	require.True(t, errors.Is(s.Validate(csize), ErrWrongCoefficientLen))

	s = NewVoteSigners(csize)
	s.increment(0, 1, common.Big0)
	s.increment(0, 1, common.Big0)
	require.True(t, errors.Is(s.Validate(csize), ErrInvalidSingleSig))

	s.increment(0, 1, common.Big0)
	s.increment(0, 1, common.Big0)
	s.increment(1, 1, common.Big0)
	require.True(t, errors.Is(s.Validate(csize), ErrInvalidCoefficient))

	s = NewVoteSigners(csize)
	s.increment(0, 1, common.Big0)
	s.Validate(csize)
	require.True(t, s.validated)
	require.Equal(t, csize, s.CommitteeSize())
	require.Equal(t, 1, s.Len())
}
