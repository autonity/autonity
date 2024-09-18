package message

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/hexutil"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/crypto/blst"
)

var (
	testKey, _          = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	testConsensusKey, _ = blst.SecretKeyFromHex("667e85b8b64622c4b8deadf59964e4c6ae38768a54dbbbc8bbd926777b896584")
	testAddr            = crypto.PubkeyToAddress(testKey.PublicKey)
	testCommitteeMember = &types.CommitteeMember{
		Address:           testAddr,
		VotingPower:       common.Big1,
		ConsensusKey:      testConsensusKey.PublicKey(),
		ConsensusKeyBytes: testConsensusKey.PublicKey().Marshal(),
		Index:             0,
	}

	testCommittee = types.Committee{{
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
	}}
	blockHash  = common.BytesToHash([]byte("123456789"))
	blockHash2 = common.BytesToHash([]byte("7890"))
)

// don't care about Address and ConsensusKey while testing the set.Add
func makeCommitteeMember(power int64, index uint64) *types.CommitteeMember {
	return &types.CommitteeMember{
		Address:           testAddr,
		VotingPower:       big.NewInt(power),
		ConsensusKey:      testConsensusKey.PublicKey(),
		ConsensusKeyBytes: testConsensusKey.PublicKey().Marshal(),
		Index:             index,
	}
}

func defaultSigner(h common.Hash) blst.Signature {
	return testConsensusKey.Sign(h[:])
}

func TestMessageSetAggregationAndPower(t *testing.T) {
	r := int64(1)
	h := uint64(1)
	csize := len(testCommittee)

	ms := NewSet()

	vote := NewPrevote(r, h, blockHash, defaultSigner, makeCommitteeMember(1, 0), csize)
	ms.Add(vote)

	require.Equal(t, common.Big1, ms.TotalPower().Power())
	require.Equal(t, common.Big1, ms.PowerFor(blockHash).Power())
	require.Equal(t, vote.Hash(), ms.Messages()[0].Hash())
	require.Equal(t, vote.Hash(), ms.VotesFor(blockHash)[0].Hash())

	// duplicated vote has no influence on power and is not saved two times
	ms.Add(vote)

	require.Equal(t, common.Big1, ms.TotalPower().Power())
	require.Equal(t, common.Big1, ms.PowerFor(blockHash).Power())
	require.Equal(t, vote.Hash(), ms.Messages()[0].Hash())
	require.Equal(t, vote.Hash(), ms.VotesFor(blockHash)[0].Hash())
	require.Equal(t, 1, len(ms.Messages()))
	require.Equal(t, 1, len(ms.VotesFor(blockHash)))

	// equivocated vote has no influence on power and it is saved
	equivocatedVote := NewPrevote(r, h, blockHash2, defaultSigner, makeCommitteeMember(1, 0), csize)
	ms.Add(equivocatedVote)

	require.Equal(t, common.Big1, ms.TotalPower().Power())
	require.Equal(t, common.Big1, ms.PowerFor(blockHash).Power())
	require.Equal(t, common.Big1, ms.PowerFor(blockHash2).Power())
	require.Equal(t, 2, len(ms.Messages()))
	require.Equal(t, 1, len(ms.VotesFor(blockHash)))
	require.Equal(t, 1, len(ms.VotesFor(blockHash2)))

	// add vote from another validator, it should get aggregated with the first one
	vote2 := NewPrevote(r, h, blockHash, defaultSigner, makeCommitteeMember(2, 1), csize)
	ms.Add(vote2)

	require.Equal(t, common.Big3, ms.TotalPower().Power())
	require.Equal(t, common.Big3, ms.PowerFor(blockHash).Power())
	require.Equal(t, 1, len(ms.VotesFor(blockHash)))
	require.True(t, ms.VotesFor(blockHash)[0].Signers().Contains(0))
	require.True(t, ms.VotesFor(blockHash)[0].Signers().Contains(1))

	// add an aggregate that cannot be merged with the previous one

	aggregate := AggregatePrevotesSimple([]Vote{NewPrevote(r, h, blockHash, defaultSigner, makeCommitteeMember(1, 0), csize), NewPrevote(r, h, blockHash, defaultSigner, makeCommitteeMember(2, 2), csize)})
	ms.Add(aggregate[0])

	require.Equal(t, common.Big5, ms.TotalPower().Power())
	require.Equal(t, common.Big5, ms.PowerFor(blockHash).Power())
	require.Equal(t, 2, len(ms.VotesFor(blockHash)))
	require.True(t, ms.VotesFor(blockHash)[0].Signers().Contains(0))
	require.True(t, ms.VotesFor(blockHash)[0].Signers().Contains(1))
	require.True(t, ms.VotesFor(blockHash)[1].Signers().Contains(0))
	require.True(t, ms.VotesFor(blockHash)[1].Signers().Contains(2))

}

func TestMessageSetAddVote(t *testing.T) {
	msg := NewPrevote(1, 1, blockHash, defaultSigner, testCommitteeMember, 1)
	ms := NewSet()
	ms.Add(msg)
	ms.Add(msg)

	require.Equal(t, common.Big1, ms.PowerFor(blockHash).Power())
	require.Equal(t, common.Big1, ms.TotalPower().Power())
}

func TestMessageSetEmpty(t *testing.T) {
	ms := NewSet()

	require.Equal(t, common.Big0, ms.PowerFor(blockHash).Power())
	require.Equal(t, common.Big0, ms.TotalPower().Power())
}

func TestMessageSetAddNilVote(t *testing.T) {
	msg := NewPrevote(1, 1, common.Hash{}, defaultSigner, testCommitteeMember, 1)
	ms := NewSet()
	ms.Add(msg)
	ms.Add(msg)
	require.Equal(t, common.Big1, ms.PowerFor(common.Hash{}).Power())
	require.Equal(t, common.Big1, ms.TotalPower().Power())
}

func TestMessageSetTotalSize(t *testing.T) {
	nilHash := common.Hash{}
	csize := len(testCommittee)

	testCases := []struct {
		voteList      []Vote
		expectedPower *big.Int
	}{{
		[]Vote{
			NewPrevote(1, 1, blockHash, defaultSigner, makeCommitteeMember(1, 0), csize),
			NewPrevote(1, 1, blockHash, defaultSigner, makeCommitteeMember(1, 1), csize),
		},
		common.Big2,
	}, {
		[]Vote{
			NewPrevote(1, 1, blockHash, defaultSigner, makeCommitteeMember(1, 0), csize),
			NewPrevote(1, 1, blockHash2, defaultSigner, makeCommitteeMember(3, 1), csize),
		},
		big.NewInt(4),
	}, {
		[]Vote{
			NewPrevote(1, 1, blockHash, defaultSigner, makeCommitteeMember(1, 0), csize),
			NewPrevote(1, 1, blockHash, defaultSigner, makeCommitteeMember(1, 1), csize),
			NewPrevote(1, 1, blockHash, defaultSigner, makeCommitteeMember(5, 2), csize),
			NewPrevote(1, 1, nilHash, defaultSigner, makeCommitteeMember(1, 3), csize),
		},
		big.NewInt(8),
	}, {
		[]Vote{
			NewPrevote(1, 1, blockHash, defaultSigner, makeCommitteeMember(0, 0), csize),
			NewPrevote(1, 1, blockHash, defaultSigner, makeCommitteeMember(1, 1), csize),
		},
		common.Big1,
	}, {
		[]Vote{
			NewPrevote(1, 1, blockHash, defaultSigner, makeCommitteeMember(1, 0), csize),
			NewPrevote(1, 1, blockHash2, defaultSigner, makeCommitteeMember(1, 1), csize),
		},
		common.Big2,
	}, {
		[]Vote{
			NewPrevote(1, 1, blockHash, defaultSigner, makeCommitteeMember(3, 0), csize),
			NewPrevote(1, 1, blockHash2, defaultSigner, makeCommitteeMember(5, 0), csize), // should be discarded
		},
		common.Big3,
	}}

	for _, test := range testCases {
		ms := NewSet()
		for _, msg := range test.voteList {
			ms.Add(msg)
		}
		if got := ms.TotalPower().Power(); got.Cmp(test.expectedPower) != 0 {
			t.Fatalf("Expected %v total voting power, got %v", test.expectedPower, got)
		}
	}
}

func TestMessageSetValues(t *testing.T) {
	t.Run("not known hash given, nil returned", func(t *testing.T) {
		ms := NewSet()
		if got := ms.VotesFor(blockHash); got != nil {
			t.Fatalf("Expected nils, got %v", got)
		}
	})

	t.Run("known hash given, message returned", func(t *testing.T) {
		msg := NewPrevote(1, 1, blockHash, defaultSigner, testCommitteeMember, 1)

		ms := NewSet()
		ms.Add(msg)

		if got := len(ms.VotesFor(blockHash)); got != 1 {
			t.Fatalf("Expected 1 message, got %v", got)
		}
	})
}
