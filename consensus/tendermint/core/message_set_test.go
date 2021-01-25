package core

import (
	"github.com/clearmatics/autonity/core/types"
	"testing"

	"github.com/clearmatics/autonity/common"
)

func TestMessageSetAddVote(t *testing.T) {
	blockHash := common.BytesToHash([]byte("123456789"))
	msg := types.ConsensusMessage{Address: common.BytesToAddress([]byte("987654321")), power: 1}

	ms := newMessageSet()
	ms.AddVote(blockHash, msg)
	ms.AddVote(blockHash, msg)

	if got := ms.VotePower(blockHash); got != 1 {
		t.Fatalf("Expected 1 vote, got %v", got)
	}
}

func TestMessageSetVotesSize(t *testing.T) {
	blockHash := common.BytesToHash([]byte("123456789"))

	ms := newMessageSet()
	if got := ms.VotePower(blockHash); got != 0 {
		t.Fatalf("Expected 0, got %v", got)
	}
}

func TestMessageSetAddNilVote(t *testing.T) {
	msg := types.ConsensusMessage{Address: common.BytesToAddress([]byte("987654321")), power: 1}

	ms := newMessageSet()
	ms.AddVote(common.Hash{}, msg)
	ms.AddVote(common.Hash{}, msg)
	if got := ms.VotePower(common.Hash{}); got != 1 {
		t.Fatalf("Expected 1 nil vote, got %v", got)
	}
}

func TestMessageSetTotalSize(t *testing.T) {
	blockHash := common.BytesToHash([]byte("123456789"))
	blockHash2 := common.BytesToHash([]byte("7890"))
	nilHash := common.Hash{}
	type vote struct {
		msg  types.ConsensusMessage
		hash common.Hash
	}
	testCases := []struct {
		voteList      []vote
		expectedPower uint64
	}{{
		[]vote{
			{types.ConsensusMessage{Address: common.BytesToAddress([]byte("1")), power: 1}, blockHash},
			{types.ConsensusMessage{Address: common.BytesToAddress([]byte("2")), power: 1}, blockHash},
		},
		2,
	}, {
		[]vote{
			{types.ConsensusMessage{Address: common.BytesToAddress([]byte("1")), power: 1}, blockHash},
			{types.ConsensusMessage{Address: common.BytesToAddress([]byte("2")), power: 3}, blockHash2},
		},
		4,
	}, {
		[]vote{
			{types.ConsensusMessage{Address: common.BytesToAddress([]byte("1")), power: 1}, blockHash},
			{types.ConsensusMessage{Address: common.BytesToAddress([]byte("2")), power: 1}, blockHash},
			{types.ConsensusMessage{Address: common.BytesToAddress([]byte("3")), power: 5}, blockHash},
			{types.ConsensusMessage{Address: common.BytesToAddress([]byte("4")), power: 1}, nilHash},
		},
		8,
	}, {
		[]vote{
			{types.ConsensusMessage{Address: common.BytesToAddress([]byte("1")), power: 1}, blockHash},
			{types.ConsensusMessage{Address: common.BytesToAddress([]byte("2")), power: 0}, blockHash},
		},
		1,
	}, {
		[]vote{
			{types.ConsensusMessage{Address: common.BytesToAddress([]byte("1")), power: 1}, blockHash},
			{types.ConsensusMessage{Address: common.BytesToAddress([]byte("2")), power: 1}, blockHash2},
		},
		2,
	}, {
		[]vote{
			{types.ConsensusMessage{Address: common.BytesToAddress([]byte("1")), power: 3}, blockHash},
			{types.ConsensusMessage{Address: common.BytesToAddress([]byte("1")), power: 5}, blockHash2}, // should be discarded
		},
		3,
	}}

	for _, test := range testCases {

		ms := newMessageSet()
		for _, msg := range test.voteList {
			ms.AddVote(msg.hash, msg.msg)
		}
		if got := ms.TotalVotePower(); got != test.expectedPower {
			t.Fatalf("Expected %v total voting power, got %v", test.expectedPower, got)
		}
	}
}

func TestMessageSetValues(t *testing.T) {
	t.Run("not known hash given, nil returned", func(t *testing.T) {
		blockHash := common.BytesToHash([]byte("123456789"))
		ms := newMessageSet()

		if got := ms.Values(blockHash); got != nil {
			t.Fatalf("Expected nils, got %v", got)
		}
	})

	t.Run("known hash given, message returned", func(t *testing.T) {
		blockHash := common.BytesToHash([]byte("123456789"))
		msg := types.ConsensusMessage{Address: common.BytesToAddress([]byte("987654321"))}

		ms := newMessageSet()
		ms.AddVote(blockHash, msg)

		if got := len(ms.Values(blockHash)); got != 1 {
			t.Fatalf("Expected 1 message, got %v", got)
		}
	})
}
