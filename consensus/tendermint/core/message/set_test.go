package message

import (
	"math/big"
	"testing"

	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/crypto/blst"

	"github.com/autonity/autonity/common"
)

var (
	testKey, _          = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	testConsensusKey, _ = blst.SecretKeyFromHex("667e85b8b64622c4b8deadf59964e4c6ae38768a54dbbbc8bbd926777b896584")
	testAddr            = crypto.PubkeyToAddress(testKey.PublicKey)
)

func defaultSigner(h common.Hash) (blst.Signature, common.Address) {
	signature := testConsensusKey.Sign(h[:])
	return signature, testAddr
}
func stubVerifier(address common.Address) *types.CommitteeMember {
	return &types.CommitteeMember{
		Address:           address,
		VotingPower:       common.Big1,
		ConsensusKey:      testConsensusKey.PublicKey(),
		ConsensusKeyBytes: testConsensusKey.PublicKey().Marshal(),
	}
}

func TestMessageSetAddVote(t *testing.T) {
	blockHash := common.BytesToHash([]byte("123456789"))
	msg := newVote[Prevote](1, 1, blockHash, defaultSigner).MustVerify(stubVerifier)
	msg.power = common.Big1
	ms := NewSet[*Prevote]()
	ms.Add(msg)
	ms.Add(msg)
	if got := ms.PowerFor(blockHash); got.Cmp(common.Big1) != 0 {
		t.Fatalf("Expected 1 vote, got %v", got)
	}
}

func TestMessageSetVotesSize(t *testing.T) {
	blockHash := common.BytesToHash([]byte("123456789"))
	ms := NewSet[*Prevote]()
	if got := ms.PowerFor(blockHash); got.Cmp(common.Big0) != 0 {
		t.Fatalf("Expected 0, got %v", got)
	}
}

func TestMessageSetAddNilVote(t *testing.T) {
	msg := newVote[Prevote](1, 1, common.Hash{}, defaultSigner).MustVerify(stubVerifier)
	ms := NewSet[*Prevote]()
	ms.Add(msg)
	ms.Add(msg)
	if got := ms.PowerFor(common.Hash{}); got.Cmp(common.Big1) != 0 {
		t.Fatalf("Expected 1 nil vote, got %v", got)
	}
}

func TestMessageSetTotalSize(t *testing.T) {
	blockHash := common.BytesToHash([]byte("123456789"))
	blockHash2 := common.BytesToHash([]byte("7890"))
	nilHash := common.Hash{}

	testCases := []struct {
		voteList      []Fake
		expectedPower *big.Int
	}{{
		[]Fake{
			{FakeSender: common.BytesToAddress([]byte("1")), FakePower: common.Big1, FakeValue: blockHash},
			{FakeSender: common.BytesToAddress([]byte("2")), FakePower: common.Big1, FakeValue: blockHash},
		},
		common.Big2,
	}, {
		[]Fake{
			{FakeSender: common.BytesToAddress([]byte("1")), FakePower: common.Big1, FakeValue: blockHash},
			{FakeSender: common.BytesToAddress([]byte("2")), FakePower: common.Big3, FakeValue: blockHash2},
		},
		big.NewInt(4),
	}, {
		[]Fake{
			{FakeSender: common.BytesToAddress([]byte("1")), FakePower: common.Big1, FakeValue: blockHash},
			{FakeSender: common.BytesToAddress([]byte("2")), FakePower: common.Big1, FakeValue: blockHash},
			{FakeSender: common.BytesToAddress([]byte("3")), FakePower: big.NewInt(5), FakeValue: blockHash},
			{FakeSender: common.BytesToAddress([]byte("4")), FakePower: common.Big1, FakeValue: nilHash},
		},
		big.NewInt(8),
	}, {
		[]Fake{
			{FakeSender: common.BytesToAddress([]byte("1")), FakePower: common.Big1, FakeValue: blockHash},
			{FakeSender: common.BytesToAddress([]byte("2")), FakePower: common.Big0, FakeValue: blockHash},
		},
		common.Big1,
	}, {
		[]Fake{
			{FakeSender: common.BytesToAddress([]byte("1")), FakePower: common.Big1, FakeValue: blockHash},
			{FakeSender: common.BytesToAddress([]byte("2")), FakePower: common.Big1, FakeValue: blockHash2},
		},
		common.Big2,
	}, {
		[]Fake{
			{FakeSender: common.BytesToAddress([]byte("1")), FakePower: common.Big3, FakeValue: blockHash},
			{FakeSender: common.BytesToAddress([]byte("1")), FakePower: big.NewInt(5), FakeValue: blockHash2}, // should be discarded
		},
		common.Big3,
	}}

	for _, test := range testCases {
		ms := NewSet[*Prevote]()
		for _, msg := range test.voteList {
			ms.Add(NewFakePrevote(msg))
		}
		if got := ms.TotalPower(); got.Cmp(test.expectedPower) != 0 {
			t.Fatalf("Expected %v total voting power, got %v", test.expectedPower, got)
		}
	}
}

func TestMessageSetValues(t *testing.T) {
	t.Run("not known hash given, nil returned", func(t *testing.T) {
		blockHash := common.BytesToHash([]byte("123456789"))
		ms := NewSet[*Prevote]()
		if got := ms.VotesFor(blockHash); got != nil {
			t.Fatalf("Expected nils, got %v", got)
		}
	})

	t.Run("known hash given, message returned", func(t *testing.T) {
		blockHash := common.BytesToHash([]byte("123456789"))
		msg := Fake{FakeSender: common.BytesToAddress([]byte("987654321")), FakeValue: blockHash}

		ms := NewSet[*Prevote]()
		ms.Add(NewFakePrevote(msg))

		if got := len(ms.VotesFor(blockHash)); got != 1 {
			t.Fatalf("Expected 1 message, got %v", got)
		}
	})
}
