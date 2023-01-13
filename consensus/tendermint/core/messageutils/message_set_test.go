package messageutils

import (
	"math/big"
	"testing"

	"github.com/autonity/autonity/common"
)

func TestMessageSetAddVote(t *testing.T) {
	blockHash := common.BytesToHash([]byte("123456789"))
	msg := Message{Address: common.BytesToAddress([]byte("987654321")), Power: common.Big1}

	ms := NewMessageSet()
	ms.AddVote(blockHash, msg)
	ms.AddVote(blockHash, msg)

	if got := ms.VotePower(blockHash); got.Cmp(common.Big1) != 0 {
		t.Fatalf("Expected 1 vote, got %v", got)
	}
}

func TestMessageSetVotesSize(t *testing.T) {
	blockHash := common.BytesToHash([]byte("123456789"))

	ms := NewMessageSet()
	if got := ms.VotePower(blockHash); got.Cmp(common.Big0) != 0 {
		t.Fatalf("Expected 0, got %v", got)
	}
}

func TestMessageSetAddNilVote(t *testing.T) {
	msg := Message{Address: common.BytesToAddress([]byte("987654321")), Power: common.Big1}

	ms := NewMessageSet()
	ms.AddVote(common.Hash{}, msg)
	ms.AddVote(common.Hash{}, msg)
	if got := ms.VotePower(common.Hash{}); got.Cmp(common.Big1) != 0 {
		t.Fatalf("Expected 1 nil vote, got %v", got)
	}
}

func TestMessageSetTotalSize(t *testing.T) {
	blockHash := common.BytesToHash([]byte("123456789"))
	blockHash2 := common.BytesToHash([]byte("7890"))
	nilHash := common.Hash{}
	type vote struct {
		msg  Message
		hash common.Hash
	}
	testCases := []struct {
		voteList      []vote
		expectedPower *big.Int
	}{{
		[]vote{
			{Message{Address: common.BytesToAddress([]byte("1")), Power: common.Big1}, blockHash},
			{Message{Address: common.BytesToAddress([]byte("2")), Power: common.Big1}, blockHash},
		},
		common.Big2,
	}, {
		[]vote{
			{Message{Address: common.BytesToAddress([]byte("1")), Power: common.Big1}, blockHash},
			{Message{Address: common.BytesToAddress([]byte("2")), Power: common.Big3}, blockHash2},
		},
		big.NewInt(4),
	}, {
		[]vote{
			{Message{Address: common.BytesToAddress([]byte("1")), Power: common.Big1}, blockHash},
			{Message{Address: common.BytesToAddress([]byte("2")), Power: common.Big1}, blockHash},
			{Message{Address: common.BytesToAddress([]byte("3")), Power: big.NewInt(5)}, blockHash},
			{Message{Address: common.BytesToAddress([]byte("4")), Power: common.Big1}, nilHash},
		},
		big.NewInt(8),
	}, {
		[]vote{
			{Message{Address: common.BytesToAddress([]byte("1")), Power: common.Big1}, blockHash},
			{Message{Address: common.BytesToAddress([]byte("2")), Power: common.Big0}, blockHash},
		},
		common.Big1,
	}, {
		[]vote{
			{Message{Address: common.BytesToAddress([]byte("1")), Power: common.Big1}, blockHash},
			{Message{Address: common.BytesToAddress([]byte("2")), Power: common.Big1}, blockHash2},
		},
		common.Big2,
	}, {
		[]vote{
			{Message{Address: common.BytesToAddress([]byte("1")), Power: common.Big3}, blockHash},
			{Message{Address: common.BytesToAddress([]byte("1")), Power: big.NewInt(5)}, blockHash2}, // should be discarded
		},
		common.Big3,
	}}

	for _, test := range testCases {

		ms := NewMessageSet()
		for _, msg := range test.voteList {
			ms.AddVote(msg.hash, msg.msg)
		}
		if got := ms.TotalVotePower(); got.Cmp(test.expectedPower) != 0 {
			t.Fatalf("Expected %v total voting power, got %v", test.expectedPower, got)
		}
	}
}

func TestMessageSetValues(t *testing.T) {
	t.Run("not known hash given, nil returned", func(t *testing.T) {
		blockHash := common.BytesToHash([]byte("123456789"))
		ms := NewMessageSet()

		if got := ms.Values(blockHash); got != nil {
			t.Fatalf("Expected nils, got %v", got)
		}
	})

	t.Run("known hash given, message returned", func(t *testing.T) {
		blockHash := common.BytesToHash([]byte("123456789"))
		msg := Message{Address: common.BytesToAddress([]byte("987654321"))}

		ms := NewMessageSet()
		ms.AddVote(blockHash, msg)

		if got := len(ms.Values(blockHash)); got != 1 {
			t.Fatalf("Expected 1 message, got %v", got)
		}
	})
}
