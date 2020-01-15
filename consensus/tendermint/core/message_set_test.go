package core

import (
	"testing"

	"github.com/clearmatics/autonity/common"
)

func TestMessageSetAddVote(t *testing.T) {
	blockHash := common.BytesToHash([]byte("123456789"))
	msg := Message{Address: common.BytesToAddress([]byte("987654321"))}

	ms := newMessageSet()
	ms.AddVote(blockHash, msg)
	ms.AddVote(blockHash, msg)

	if got := ms.VoteCount(blockHash); got != 1 {
		t.Fatalf("Expected 1 vote, got %v", got)
	}
}

func TestMessageSetVotesSize(t *testing.T) {
	blockHash := common.BytesToHash([]byte("123456789"))

	ms := newMessageSet()
	if got := ms.VoteCount(blockHash); got != 0 {
		t.Fatalf("Expected 0, got %v", got)
	}
}

func TestMessageSetAddNilVote(t *testing.T) {
	msg := Message{Address: common.BytesToAddress([]byte("987654321"))}

	ms := newMessageSet()
	ms.AddNilVote(msg)
	ms.AddNilVote(msg)

	if got := ms.NilVotesSize(); got != 1 {
		t.Fatalf("Expected 1 nil vote, got %v", got)
	}
}

func TestMessageSetTotalSize(t *testing.T) {
	blockHash := common.BytesToHash([]byte("123456789"))
	msg := Message{Address: common.BytesToAddress([]byte("987654321"))}

	ms := newMessageSet()
	ms.AddVote(blockHash, msg)
	ms.AddVote(blockHash, msg)

	ms.AddNilVote(msg)
	ms.AddNilVote(msg)

	if got := ms.TotalVoteCount(); got != 2 {
		t.Fatalf("Expected 2 total votes, got %v", got)
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
		msg := Message{Address: common.BytesToAddress([]byte("987654321"))}

		ms := newMessageSet()
		ms.AddVote(blockHash, msg)

		if got := len(ms.Values(blockHash)); got != 1 {
			t.Fatalf("Expected 1 message, got %v", got)
		}
	})
}
