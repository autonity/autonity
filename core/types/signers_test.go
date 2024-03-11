package types

import (
	"bytes"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/hexutil"
	"github.com/autonity/autonity/rlp"
)

var (
	header = &Header{
		ParentHash:  common.HexToHash("0a5843ac1cb04865017cb35a57b50b07084e5fcee39b5acadade33149f4fff9e"),
		UncleHash:   common.HexToHash("0a5843ac1c124732472342342387423897431293123020912dade33149f4fffe"),
		Coinbase:    common.HexToAddress("8888f1f195afa192cfee860698584c030f4c9db1"),
		Root:        common.HexToHash("0a5843ac1cb0486345235234564778768967856745645654649f4fff9321321e"),
		TxHash:      common.HexToHash("0a58213121cb0486345235234564778768967856745645654649f4fff932132e"),
		ReceiptHash: common.HexToHash("9a58213121cb0486345235234564778768967856745645654649f4fff932132e"),
		Bloom:       BytesToBloom(bytes.Repeat([]byte("a"), 128)),
		Difficulty:  big.NewInt(199),
		Number:      big.NewInt(239),
		GasLimit:    uint64(1000),
		GasUsed:     uint64(400),
		Time:        uint64(12343),
		MixDigest:   common.HexToHash("0a58213121cb0486345235234564778768967853123645654649f4fff932132e"),
		Nonce:       [8]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
		BaseFee:     big.NewInt(20000),
		Committee: []CommitteeMember{{
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
			ConsensusKeyBytes: hexutil.MustDecode("0xa679c8ebd47d18b93acd90cd380debdcfdb140f38eca207c61463a47be85398ec3082a66f7f30635c11470f5c8e5cf6b"),
			VotingPower:       hexutil.MustDecodeBig("0x7033"),
		}},
		ProposerSeal:      bytes.Repeat([]byte("c"), 65),
		Round:             uint64(3),
		QuorumCertificate: new(AggregateSignature),
	}
)

// TODO(lorenzo) add tests for malformed bit and coefficient length (e.g. longer and shorter than supposed)

func TestValidatorBitmap(t *testing.T) {
	t.Run("Simple bitmap with 4 validators (1 byte)", func(t *testing.T) {
		// 00 00 00 00
		n := 4
		bitmap := NewValidatorBitmap(n)

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
	//TODO(lorenzo) refinements2, boundaries tests
}

func TestSerialization(t *testing.T) {
	t.Run("auxiliary data structures are not serialized on the wire", func(t *testing.T) {
		// increment some sender info
		s := NewSigners(len(header.Committee))
		// 11 01 01 10
		s.Increment(header, 0)
		s.Increment(header, 0)
		s.Increment(header, 0)
		s.Increment(header, 0)
		s.Increment(header, 1)
		s.Increment(header, 2)
		s.Increment(header, 3)
		s.Increment(header, 3)

		// bits, bytes + auxiliary data structures should be set
		require.Equal(t, s.Bits[0], byte(0xd6))
		require.Equal(t, len(s.Coefficients), 1)
		require.Equal(t, s.Coefficients[0], uint16(4))
		// one validator is missing from the aggregate
		require.Equal(t, len(s.Addresses()), len(header.Committee)-1)
		require.Equal(t, len(s.Powers()), len(header.Committee)-1)
		require.Equal(t, s.CommitteeSize(), len(header.Committee))

		// encode and decode
		payload, err := rlp.EncodeToBytes(s)
		require.NoError(t, err)
		decoded := &Signers{}
		err = rlp.Decode(bytes.NewBuffer(payload), decoded)
		require.NoError(t, err)

		// only bits and bytes should be set, auxiliary data structures should be empty
		require.Equal(t, s.Bits[0], byte(0xd6))
		require.Equal(t, len(s.Coefficients), 1)
		require.Equal(t, s.Coefficients[0], uint16(4))
		require.Nil(t, decoded.Addresses())
		require.Nil(t, decoded.Powers())
		require.Equal(t, decoded.CommitteeSize(), 0)

	})
}

func TestSigners(t *testing.T) {
	//TODO(lorenzo) refinements2, boundaries tests
	t.Run("Increment should update the bitmap and the auxiliary maps correctly", func(t *testing.T) {
		// +10 to avoid hitting the panic in `increment` related to the max allowed coefficient
		s := NewSigners(len(header.Committee) + 10)

		require.Equal(t, s.Bits.Get(0), byte(0))
		s.Increment(header, 0)

		require.Equal(t, len(s.Addresses()), 1)
		require.Equal(t, len(s.Powers()), 1)
		require.Equal(t, s.Bits.Get(0), byte(1))
		require.Equal(t, s.Powers()[0], header.Committee[0].VotingPower)
		require.Equal(t, s.Addresses()[0], header.Committee[0].Address)

		s.Increment(header, 0)
		require.Equal(t, s.Bits.Get(0), byte(2))
		s.Increment(header, 0)
		require.Equal(t, s.Bits.Get(0), byte(3))
		require.Equal(t, s.Coefficients[0], uint16(3))
		s.Increment(header, 0)
		require.Equal(t, s.Bits.Get(0), byte(3))
		require.Equal(t, s.Coefficients[0], uint16(4))
		s.Increment(header, 0)
		require.Equal(t, s.Bits.Get(0), byte(3))
		require.Equal(t, s.Coefficients[0], uint16(5))

		require.Equal(t, len(s.Addresses()), 1)
		require.Equal(t, len(s.Powers()), 1)
		require.Equal(t, s.Powers()[0], header.Committee[0].VotingPower)
		require.Equal(t, s.Addresses()[0], header.Committee[0].Address)

		s.Increment(header, 1)

		require.Equal(t, len(s.Addresses()), 2)
		require.Equal(t, len(s.Powers()), 2)
		require.Equal(t, s.Bits.Get(1), uint8(1))
		require.Equal(t, s.Powers()[1], header.Committee[1].VotingPower)
		require.Equal(t, s.Addresses()[1], header.Committee[1].Address)

		s.Increment(header, 2)
		require.Equal(t, s.Bits.Get(2), uint8(1))
		s.Increment(header, 2)
		require.Equal(t, s.Bits.Get(2), uint8(2))
		s.Increment(header, 2)
		require.Equal(t, s.Bits.Get(2), uint8(3))
		require.Equal(t, s.Coefficients[1], uint16(3))
		s.Increment(header, 2)
		require.Equal(t, s.Bits.Get(2), uint8(3))
		require.Equal(t, s.Coefficients[1], uint16(4))
		s.Increment(header, 2)
		require.Equal(t, s.Bits.Get(2), uint8(3))
		require.Equal(t, s.Coefficients[1], uint16(5))

		s.Increment(header, 2)
		s.Increment(header, 2)
		s.Increment(header, 2)
		require.Equal(t, s.Coefficients[1], uint16(8))

		s.Increment(header, 0)
		require.Equal(t, s.Coefficients[0], uint16(6))

		s.Increment(header, 1)
		require.Equal(t, s.Bits.Get(1), uint8(2))
		s.Increment(header, 1)
		require.Equal(t, s.Bits.Get(1), uint8(3))
		require.Equal(t, s.Coefficients[1], uint16(3))
		s.Increment(header, 1)
		require.Equal(t, s.Bits.Get(1), uint8(3))
		require.Equal(t, s.Coefficients[1], uint16(4))
		s.Increment(header, 1)
		require.Equal(t, s.Bits.Get(1), uint8(3))
		require.Equal(t, s.Coefficients[1], uint16(5))

		s.Increment(header, 1)
		require.Equal(t, s.Bits.Get(1), uint8(3))
		require.Equal(t, s.Coefficients[2], uint16(8))
	})
	t.Run("Merge correctly merges two senders info", func(t *testing.T) {
		// +10 to avoid hitting the panic in `increment` related to the max allowed coefficient
		s1 := NewSigners(len(header.Committee) + 10)
		s1.Increment(header, 0)
		s1.Increment(header, 1)
		s2 := NewSigners(len(header.Committee) + 10)
		s2.Increment(header, 2)
		s2.Increment(header, 3)
		s2.Increment(header, 4)

		s1.Merge(s2)

		require.Equal(t, len(s1.Addresses()), len(header.Committee))
		require.Equal(t, len(s1.Powers()), len(header.Committee))
		for i, member := range header.Committee {
			require.Equal(t, s1.Bits.Get(i), uint8(1))
			require.Equal(t, s1.Powers()[i], member.VotingPower)
			require.Equal(t, s1.Addresses()[i], member.Address)
		}

		s3 := NewSigners(len(header.Committee) + 10)
		s3.Increment(header, 0)
		s3.Increment(header, 1)
		s3.Increment(header, 2)
		s3.Increment(header, 3)
		s3.Increment(header, 4)

		s1.Merge(s3)

		require.Equal(t, len(s1.Addresses()), len(header.Committee))
		require.Equal(t, len(s1.Powers()), len(header.Committee))
		for i, member := range header.Committee {
			require.Equal(t, s1.Bits.Get(i), uint8(2))
			require.Equal(t, s1.Powers()[i], member.VotingPower)
			require.Equal(t, s1.Addresses()[i], member.Address)
		}

		s4 := NewSigners(len(header.Committee) + 10)
		s4.Increment(header, 0)
		s4.Increment(header, 0)
		s4.Increment(header, 0)
		s4.Increment(header, 2)
		s4.Increment(header, 2)
		s4.Increment(header, 2)
		s4.Increment(header, 2)

		s1.Merge(s4)

		require.Equal(t, s1.Bits.Get(0), uint8(3))
		require.Equal(t, s1.Coefficients[0], uint16(5))
		require.Equal(t, s1.Bits.Get(2), uint8(3))
		require.Equal(t, s1.Coefficients[1], uint16(6))

		s5 := NewSigners(len(header.Committee) + 10)
		s5.Increment(header, 1)

		s1.Merge(s5)
		require.Equal(t, s1.Bits.Get(0), uint8(3))
		require.Equal(t, s1.Bits.Get(1), uint8(3))
		require.Equal(t, s1.Bits.Get(2), uint8(3))
		require.Equal(t, s1.Coefficients[0], uint16(5))
		require.Equal(t, s1.Coefficients[1], uint16(3))
		require.Equal(t, s1.Coefficients[2], uint16(6))

	})
	t.Run("Power returns the aggregated power of the senders", func(t *testing.T) {
		s := NewSigners(len(header.Committee))
		s.Increment(header, 0)
		s.Increment(header, 1)

		require.Equal(t, s.Power(), new(big.Int).Add(header.Committee[0].VotingPower, header.Committee[1].VotingPower))

		s.Increment(header, 2)
		s.Increment(header, 3)
		s.Increment(header, 4)

		require.Equal(t, s.Power(), header.TotalVotingPower())

		// duplicated power shouldn't be counted
		s.Increment(header, 2)
		s.Increment(header, 3)

		require.Equal(t, s.Power(), header.TotalVotingPower())
	})
	t.Run("Flatten returns the indexes of the senders (repeated in case of multiple contribution to the signature)", func(t *testing.T) {
		// +10 to avoid hitting the panic related to maximum allowed coefficient in `increment`
		s := NewSigners(len(header.Committee) + 10)
		s.Increment(header, 0)
		s.Increment(header, 1)

		require.Equal(t, s.Flatten(), []int{0, 1})

		s.Increment(header, 3)
		s.Increment(header, 3)
		s.Increment(header, 0)
		s.Increment(header, 0)
		s.Increment(header, 0)
		s.Increment(header, 1)
		s.Increment(header, 2)

		require.Equal(t, s.Flatten(), []int{0, 0, 0, 0, 1, 1, 2, 3, 3})

		s.Increment(header, 1)
		s.Increment(header, 1)
		s.Increment(header, 1)
		s.Increment(header, 1)
		s.Increment(header, 1)

		require.Equal(t, s.Flatten(), []int{0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 2, 3, 3})
	})
	t.Run("FlattenUniq returns the indexes of the senders (de-duplicated)", func(t *testing.T) {
		s := NewSigners(len(header.Committee))
		s.Increment(header, 0)
		s.Increment(header, 1)

		require.Equal(t, s.FlattenUniq(), []int{0, 1})

		s.Increment(header, 3)
		s.Increment(header, 3)
		s.Increment(header, 0)
		s.Increment(header, 0)
		s.Increment(header, 0)
		s.Increment(header, 1)
		s.Increment(header, 2)

		require.Equal(t, s.FlattenUniq(), []int{0, 1, 2, 3})
	})
}

func TestValidation(t *testing.T) {
	csize := 10
	t.Run("Validate aggregates should be valid", func(t *testing.T) {
		s := NewSigners(csize)

		// A
		s.increment(0)
		fmt.Println(s.String())
		require.True(t, s.Validate(csize))

		// A + B
		s.increment(1)
		fmt.Println(s.String())
		require.True(t, s.Validate(csize))

		// A + B + C
		s.increment(2)
		fmt.Println(s.String())
		require.True(t, s.Validate(csize))

		// A + B + C + D
		s.increment(3)
		fmt.Println(s.String())
		require.True(t, s.Validate(csize))

		// 2A + B + C + D
		s.increment(0)
		fmt.Println(s.String())
		require.True(t, s.Validate(csize))

		// 3A + B + C + D
		s.increment(0)
		fmt.Println(s.String())
		require.True(t, s.Validate(csize))

		// 4A + 2B + C + D
		s.increment(0)
		s.increment(1)
		fmt.Println(s.String())
		require.True(t, s.Validate(csize))

		// 3A + 3B + C + D + E
		s = NewSigners(csize)
		s.increment(0)
		s.increment(0)
		s.increment(0)
		s.increment(1)
		s.increment(1)
		s.increment(1)
		s.increment(2)
		s.increment(3)
		s.increment(4)
		fmt.Println(s.String())
		require.True(t, s.Validate(csize))
	})
	t.Run("Invalid aggregates should be invalid", func(t *testing.T) {
		s := NewSigners(csize)

		// 2A
		s.increment(0)
		s.increment(0)
		fmt.Println(s.String())

		// 2A + B
		s.increment(1)
		fmt.Println(s.String())
		require.False(t, s.Validate(csize))

		// 3A + B + C
		s.increment(0)
		s.increment(2)
		fmt.Println(s.String())
		require.False(t, s.Validate(csize))

		// 4A + 2B + C
		s.increment(0)
		s.increment(1)
		fmt.Println(s.String())
		require.False(t, s.Validate(csize))

		// 5A + 2B + C + D
		s.increment(0)
		s.increment(3)
		fmt.Println(s.String())
		require.False(t, s.Validate(csize))

		// 5A + 3B + C + D
		s.increment(1)
		fmt.Println(s.String())
		require.False(t, s.Validate(csize))

		// 5A + 3B + C + 2D
		s.increment(3)
		fmt.Println(s.String())
		require.False(t, s.Validate(csize))

		// 9A + 3B + C + 2D
		s.increment(0)
		s.increment(0)
		s.increment(0)
		s.increment(0)
		fmt.Println(s.String())
		require.False(t, s.Validate(csize))

		s = NewSigners(csize)

		// 2A + 2B + C
		s.increment(0)
		s.increment(0)
		s.increment(1)
		s.increment(1)
		s.increment(2)
		fmt.Println(s.String())
		require.False(t, s.Validate(csize))

		// 2A + 2B + 4C
		s.increment(2)
		s.increment(2)
		s.increment(2)
		fmt.Println(s.String())
		require.False(t, s.Validate(csize))

		s = NewSigners(csize)

		// 3A + B + C + 8D + E
		s.increment(0)
		s.increment(0)
		s.increment(0)
		s.increment(1)
		s.increment(2)
		s.increment(3)
		s.increment(3)
		s.increment(3)
		s.increment(3)
		s.increment(3)
		s.increment(3)
		s.increment(3)
		s.increment(3)
		s.increment(4)
		fmt.Println(s.String())
		require.False(t, s.Validate(csize))

		// TODO(lorenzo) add some more cases with greater committesize
		// It would be good to do some sort of fuzz testing of different combinations of signatures

	})
}
