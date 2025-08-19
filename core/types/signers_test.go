package types

import (
	"bytes"
	"errors"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
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

func newCommitte(members []CommitteeMember) *Committee {
	return &Committee{
		Members: members,
	}
}

func signersAggregation(
	t *testing.T,
	members []CommitteeMember,
	coeffs []int,
	signature blst.Signature,
	msgInput common.Hash,
) {
	lc := newCommitte(members)
	signers := NewSigners(lc)
	keys := make([]blst.PublicKey, 0, lc.Len())
	for i, c := range coeffs {
		for c > 0 {
			signers.AddSigner(uint64(i)) // nolint
			c--
		}
		keys = append(keys, lc.MemberByIndex(i).ConsensusKey)
	}

	publicKey := signers.AggregatePublicKey(keys)
	require.True(t, publicKey.Validate())
	require.True(t, signature.Verify(publicKey, msgInput[:], blst.DefaultAssumeZeroValid))
}

func TestPublicKeyAggregation(t *testing.T) {
	csize := 33 // at least 4
	members := make([]CommitteeMember, 0, csize)
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

		members = append(members, CommitteeMember{
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
	signersAggregation(t, members, coeffs, aggregateSign, msg)
	signersAggregation(t, members, coeffs, aggregateSign, msg)
}

func TestSerialization(t *testing.T) {
	err := committee.Enrich()
	require.NoError(t, err)

	t.Run("auxiliary data structures are not serialized on the wire", func(t *testing.T) {
		// increment some sender info
		s := NewSigners(committee)
		// 00001111
		s.AddSigner(0)
		s.AddSigner(1)
		s.AddSigner(1)
		s.AddSigner(2)
		s.AddSigner(2)
		s.AddSigner(2)
		s.AddSigner(3)

		// bits, coefficients + auxiliary data structures should be set
		require.Equal(t, (*big.Int)(s.Bitmap).Bytes()[0], byte(0xf))
		require.Equal(t, len(s.Coefficients), 4)
		require.Equal(t, s.Coefficients[0].String(), common.Big1.String())
		require.Equal(t, s.Coefficients[1].String(), common.Big2.String())
		require.Equal(t, s.Coefficients[2].String(), common.Big3.String())
		require.Equal(t, s.Coefficients[3].String(), common.Big1.String())

		// one validator is missing from the aggregate
		require.Equal(t, s.Len(), committee.Len()-1)
		expectedVotingPower := committee.Members[0].VotingPower.Uint64() + committee.Members[1].VotingPower.Uint64() +
			committee.Members[2].VotingPower.Uint64() + committee.Members[3].VotingPower.Uint64()
		require.Equal(t, s.Power().Uint64(), expectedVotingPower)
		require.Equal(t, s.CommitteeSize(), committee.Len())
		require.True(t, s.validated)

		// encode and decode
		payload, err := rlp.EncodeToBytes(s)
		require.NoError(t, err)
		decoded := &Signers{}
		err = rlp.Decode(bytes.NewBuffer(payload), decoded)
		require.NoError(t, err)

		// only bits and coefficients should be set, auxiliary data structures should be empty
		require.Equal(t, (*big.Int)(decoded.Bitmap).Bytes()[0], byte(0xf))
		require.Equal(t, len(decoded.Coefficients), 4)
		require.Equal(t, s.Coefficients[0].String(), common.Big1.String())
		require.Equal(t, s.Coefficients[1].String(), common.Big2.String())
		require.Equal(t, s.Coefficients[2].String(), common.Big3.String())
		require.Equal(t, s.Coefficients[3].String(), common.Big1.String())

		require.False(t, decoded.validated)

		// accessing auxiliary data structure should panic
		defer expectPanic(t)
		require.Nil(t, decoded.Power())
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
		s := NewSigners(committee)
		require.Equal(t, s.CommitteeSize(), committee.Len())

		require.False(t, s.Bitmap.IsSet(0))

		s.AddSigner(0)
		require.True(t, s.Bitmap.IsSet(0))
		require.Equal(t, common.Big1.String(), s.Coefficients[0].String())
		require.Equal(t, s.Power().Uint64(), committee.Members[0].VotingPower.Uint64())
		require.Equal(t, s.Len(), 1)

		s.AddSigner(0)
		require.True(t, s.Bitmap.IsSet(0))
		require.Equal(t, s.Coefficients[0].String(), common.Big2.String())
		s.AddSigner(0)
		require.True(t, s.Bitmap.IsSet(0))
		require.Equal(t, s.Coefficients[0].String(), common.Big3.String())

		require.Equal(t, s.Power().Uint64(), committee.Members[0].VotingPower.Uint64())
		require.Equal(t, s.Len(), 1)

		s.AddSigner(1)
		require.True(t, s.Bitmap.IsSet(1))
		require.Equal(t, s.Coefficients[1].String(), common.Big1.String())

		require.Equal(t, s.Power().Uint64(), committee.Members[0].VotingPower.Uint64()+committee.Members[1].VotingPower.Uint64())
		require.Equal(t, s.Len(), 2)

		s.AddSigner(2)
		require.True(t, s.Bitmap.IsSet(2))
		require.Equal(t, s.Coefficients[2].String(), common.Big1.String())
		s.AddSigner(2)
		require.True(t, s.Bitmap.IsSet(2))
		require.Equal(t, s.Coefficients[2].String(), common.Big2.String())
		s.AddSigner(2)
		require.True(t, s.Bitmap.IsSet(2))
		require.Equal(t, s.Coefficients[2].String(), common.Big3.String())

		s.AddSigner(2)
		s.AddSigner(2)
		s.AddSigner(2)
		require.Equal(t, s.Coefficients[2].String(), big.NewInt(6).String())

		s.AddSigner(0)
		require.Equal(t, s.Coefficients[0].String(), common.Big4.String())

		s.AddSigner(1)
		require.True(t, s.Bitmap.IsSet(1))
		require.Equal(t, s.Coefficients[1].String(), common.Big2.String())
		s.AddSigner(1)
		require.True(t, s.Bitmap.IsSet(1))
		require.Equal(t, s.Coefficients[1].String(), common.Big3.String())
		s.AddSigner(1)
		require.True(t, s.Bitmap.IsSet(1))
		require.Equal(t, s.Coefficients[1].String(), common.Big4.String())

		s.AddSigner(1)
		require.True(t, s.Bitmap.IsSet(1))
		require.Equal(t, s.Coefficients[1].String(), big.NewInt(5).String())
	})
	t.Run("Merge correctly merges two senders info", func(t *testing.T) {
		// +10 to avoid hitting the panic in `increment` related to the max allowed coefficient
		s1 := NewSigners(committee)
		s1.AddSigner(0)
		s1.AddSigner(1)
		s2 := NewSigners(committee)
		s2.AddSigner(2)
		s2.AddSigner(3)
		s2.AddSigner(4)

		s1.Merge(s2)

		require.Equal(t, s1.CommitteeSize(), committee.Len())
		require.Equal(t, s1.Len(), committee.Len())
		for i := range committee.Members {
			require.True(t, s1.Bitmap.IsSet(i))
			require.Equal(t, s1.Coefficients[i].String(), common.Big1.String())
		}
		require.Equal(t, s1.Power().Uint64(), totalPower.Uint64())

		s3 := NewSigners(committee)
		s3.AddSigner(0)
		s3.AddSigner(1)
		s3.AddSigner(2)
		s3.AddSigner(3)
		s3.AddSigner(4)

		s1.Merge(s3)

		require.Equal(t, s1.CommitteeSize(), committee.Len())
		require.Equal(t, s1.Len(), committee.Len())
		for i := range committee.Members {
			require.True(t, s1.Bitmap.IsSet(i))
			require.Equal(t, s1.Coefficients[i].String(), common.Big2.String())
		}
		require.Equal(t, s1.Power().Uint64(), totalPower.Uint64())

		s4 := NewSigners(committee)
		s4.AddSigner(0)
		s4.AddSigner(0)
		s4.AddSigner(0)
		s4.AddSigner(2)
		s4.AddSigner(2)
		s4.AddSigner(2)
		s4.AddSigner(2)

		s1.Merge(s4)

		require.True(t, s1.Bitmap.IsSet(0))
		require.Equal(t, s1.Coefficients[0].String(), big.NewInt(5).String())
		require.True(t, s1.Bitmap.IsSet(2))
		require.Equal(t, s1.Coefficients[2].String(), big.NewInt(6).String())

		s5 := NewSigners(committee)
		s5.AddSigner(1)

		s1.Merge(s5)
		require.Equal(t, true, s1.Bitmap.IsSet(0))
		require.Equal(t, true, s1.Bitmap.IsSet(1))
		require.Equal(t, true, s1.Bitmap.IsSet(2))
		require.Equal(t, s1.Coefficients[0].String(), big.NewInt(5).String())
		require.Equal(t, s1.Coefficients[1].String(), common.Big3.String())
		require.Equal(t, s1.Coefficients[2].String(), big.NewInt(6).String())
	})
	t.Run("Power returns the aggregated power of the senders", func(t *testing.T) {
		s := NewSigners(committee)
		s.AddSigner(0)
		s.AddSigner(1)

		require.Equal(t, s.Power(), new(big.Int).Add(committee.Members[0].VotingPower, committee.Members[1].VotingPower))

		s.AddSigner(2)
		s.AddSigner(3)
		s.AddSigner(4)

		require.Equal(t, s.Power(), totalPower)

		// duplicated power shouldn't be counted
		s.AddSigner(2)
		s.AddSigner(3)

		require.Equal(t, s.Power(), totalPower)
	})
	t.Run("FlattenUniq returns the indexes of the senders (de-duplicated)", func(t *testing.T) {
		s := NewSigners(committee)
		s.AddSigner(0)
		s.AddSigner(1)

		require.Equal(t, s.FlattenUniq(), []int{0, 1})

		s.AddSigner(3)
		s.AddSigner(3)
		s.AddSigner(0)
		s.AddSigner(0)
		s.AddSigner(0)
		s.AddSigner(1)
		s.AddSigner(2)

		require.Equal(t, s.FlattenUniq(), []int{0, 1, 2, 3})
	})
	t.Run("Contains returns expected result", func(t *testing.T) {
		s := NewSigners(committee)

		require.False(t, s.Contains(0))
		require.False(t, s.Contains(1))
		require.False(t, s.Contains(2))

		s.AddSigner(0)
		s.AddSigner(2)

		require.True(t, s.Contains(0))
		require.False(t, s.Contains(1))
		require.True(t, s.Contains(2))
	})
	t.Run("AddsInformation, RespectBoundaries return expected results", func(t *testing.T) {
		s1 := NewSigners(committee)
		s2 := NewSigners(committee)

		require.False(t, s1.AddsInformation(s2))

		s1.increment(0) //0

		require.False(t, s1.AddsInformation(s2))

		s2.increment(0) //0

		require.False(t, s1.AddsInformation(s2))

		s2.increment(1) //0,1

		require.True(t, s1.AddsInformation(s2))

		s2.increment(2) //0,1,2

		require.True(t, s1.AddsInformation(s2))

		s1.increment(0)
		s1.increment(0) //0,0,0
		s2.increment(0) //0,0,1,2

		require.True(t, s1.AddsInformation(s2))

		s2.increment(0) //0,0,0,1,2

		require.True(t, s1.AddsInformation(s2))

		s1.increment(1)
		s1.increment(2)

		require.False(t, s1.AddsInformation(s2))

	})
}

func TestValidation(t *testing.T) {

	nilSigner := &Signers{
		Bitmap:       nil,
		Coefficients: nil,
	}
	err := nilSigner.SanityCheck()
	t.Log(err)
	require.Error(t, err)

	csize := committee.Len()
	wrongSizeSigner := NewSigners(committee)
	wrongSizeSigner.Coefficients = make([]*big.Int, csize+10)
	require.True(t, errors.Is(wrongSizeSigner.Validate(committee), ErrWrongSizeSigners))

	s := NewSigners(committee)
	require.True(t, errors.Is(s.Validate(committee), ErrEmptySigners))

	// A
	s.increment(0)
	t.Log(s.String())
	require.Nil(t, s.Validate(committee))

	// A + B
	s.increment(1)
	t.Log(s.String())
	require.Nil(t, s.Validate(committee))

	// A + B + C
	s.increment(2)
	t.Log(s.String())
	require.Nil(t, s.Validate(committee))

	// A + B + C + D
	s.increment(3)
	t.Log(s.String())
	require.Nil(t, s.Validate(committee))

	// 2A + B + C + D
	s.increment(0)
	t.Log(s.String())
	require.Nil(t, s.Validate(committee))

	// 3A + B + C + D
	s.increment(0)
	t.Log(s.String())
	require.Nil(t, s.Validate(committee))

	// 4A + 2B + C + D
	s.increment(0)
	s.increment(1)
	t.Log(s.String())
	require.Nil(t, s.Validate(committee))

	s.Coefficients = append(s.Coefficients, big.NewInt(0xca))
	require.True(t, errors.Is(s.Validate(committee), ErrWrongCoefficientLen))

	s = NewSigners(committee)
	// 2A
	s.increment(0)
	s.increment(0)
	require.True(t, errors.Is(s.Validate(committee), ErrInvalidSingleSig))

	// 4A + B
	s.increment(0)
	s.increment(0)
	s.increment(1)
	require.True(t, errors.Is(s.Validate(committee), ErrInvalidCoefficient))

	s = NewSigners(committee)
	s.increment(0)
	s.Validate(committee)
	require.True(t, s.validated)
	require.Equal(t, csize, s.CommitteeSize())
	require.Equal(t, 1, s.Len())
}

func TestRightmostSigners(t *testing.T) {
	err := committee.Enrich()
	require.NoError(t, err)

	s := NewSigners(committee)
	require.Equal(t, committee.Len(), s.RightmostSigner())
	s.AddSigner(3)
	t.Log(s.String())
	require.Equal(t, 3, s.RightmostSigner())
	s.AddSigner(4)
	t.Log(s.String())
	require.Equal(t, 3, s.RightmostSigner())
	s.AddSigner(2)
	t.Log(s.String())
	require.Equal(t, 2, s.RightmostSigner())

	other := NewSigners(committee)
	other.AddSigner(4)
	other.AddSigner(2)
	other.AddSigner(1)
	require.Equal(t, 1, other.RightmostSigner())

	s.Merge(other)
	require.Equal(t, 1, s.RightmostSigner())

	signerCopy := s.Copy()
	require.Equal(t, 1, signerCopy.RightmostSigner())

	payload, err := rlp.EncodeToBytes(signerCopy)
	require.NoError(t, err)
	decoded := &Signers{}
	err = rlp.Decode(bytes.NewBuffer(payload), decoded)
	require.NoError(t, err)

	t.Log(decoded.String())
	err = decoded.Validate(committee)
	require.NoError(t, err)

	require.Equal(t, 1, decoded.RightmostSigner())

	decoded2 := &Signers{}
	err = rlp.Decode(bytes.NewBuffer(payload), decoded2)
	require.NoError(t, err)

	defer expectPanic(t)
	decoded2.RightmostSigner()
}
