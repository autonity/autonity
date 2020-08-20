package core

import (
	"github.com/clearmatics/autonity/core/types"
	"math/big"
	"testing"

	"github.com/clearmatics/autonity/common"
)

func TestMessageSetAddVote(t *testing.T) {
	blockHash := common.BytesToHash([]byte("123456789"))
	addr := common.BytesToAddress([]byte("987654321"))
	cm := map[common.Address]*types.CommitteeMember{addr: {addr, big.NewInt(1)}}
	msg := Message{Address: addr}

	ms := newMessageSet()
	ms.AddVote(blockHash, msg)
	ms.AddVote(blockHash, msg)

	if got := ms.VotePower(blockHash, cm); got != 1 {
		t.Fatalf("Expected 1 vote, got %v", got)
	}
}

func TestMessageSetVotesSize(t *testing.T) {
	blockHash := common.BytesToHash([]byte("123456789"))

	ms := newMessageSet()
	if got := ms.VotePower(blockHash, map[common.Address]*types.CommitteeMember{}); got != 0 {
		t.Fatalf("Expected 0, got %v", got)
	}
}

func TestMessageSetAddNilVote(t *testing.T) {
	addr := common.BytesToAddress([]byte("987654321"))
	cm := map[common.Address]*types.CommitteeMember{addr: {addr, big.NewInt(1)}}
	msg := Message{Address: addr}

	ms := newMessageSet()
	ms.AddVote(common.Hash{}, msg)
	ms.AddVote(common.Hash{}, msg)
	if got := ms.VotePower(common.Hash{}, cm); got != 1 {
		t.Fatalf("Expected 1 nil vote, got %v", got)
	}
}

func TestMessageSetTotalSize(t *testing.T) {
	blockHash := common.BytesToHash([]byte("123456789"))
	blockHash2 := common.BytesToHash([]byte("7890"))
	nilHash := common.Hash{}

	addr1, addr2, addr3, addr4 := common.BytesToAddress([]byte("1")), common.BytesToAddress([]byte("2")), common.BytesToAddress([]byte("3")), common.BytesToAddress([]byte("4"))
	type vote struct {
		msg  Message
		hash common.Hash
	}
	testCases := []struct {
		cm            map[common.Address]*types.CommitteeMember
		voteList      []vote
		expectedPower uint64
	}{{
		map[common.Address]*types.CommitteeMember{
			addr1: {addr1, big.NewInt(1)},
			addr2: {addr2, big.NewInt(1)},
		},
		[]vote{
			{Message{Address: addr1}, blockHash},
			{Message{Address: addr2}, blockHash},
		},
		2,
	}, {
		map[common.Address]*types.CommitteeMember{
			addr1: {addr1, big.NewInt(1)},
			addr2: {addr2, big.NewInt(3)},
		},
		[]vote{
			{Message{Address: addr1}, blockHash},
			{Message{Address: addr2}, blockHash2},
		},
		4,
	}, {
		map[common.Address]*types.CommitteeMember{
			addr1: {addr1, big.NewInt(1)},
			addr2: {addr2, big.NewInt(1)},
			addr3: {addr3, big.NewInt(5)},
			addr4: {addr4, big.NewInt(1)},
		},
		[]vote{
			{Message{Address: addr1}, blockHash},
			{Message{Address: addr2}, blockHash},
			{Message{Address: addr3}, blockHash},
			{Message{Address: addr4}, nilHash},
		},
		8,
	}, {
		map[common.Address]*types.CommitteeMember{
			addr1: {addr1, big.NewInt(1)},
			addr2: {addr2, big.NewInt(0)},
		},
		[]vote{
			{Message{Address: addr1}, blockHash},
			{Message{Address: addr2}, blockHash},
		},
		1,
	}, {
		map[common.Address]*types.CommitteeMember{
			addr1: {addr1, big.NewInt(1)},
			addr2: {addr2, big.NewInt(1)},
		},
		[]vote{
			{Message{Address: addr1}, blockHash},
			{Message{Address: addr2}, blockHash2},
		},
		2,
	}, {
		map[common.Address]*types.CommitteeMember{
			addr1: {addr1, big.NewInt(3)},
			addr2: {addr2, big.NewInt(5)},
		},
		[]vote{
			{Message{Address: addr1}, blockHash},
			{Message{Address: addr1}, blockHash2}, // should be discarded
		},
		3,
	}}

	for _, test := range testCases {

		ms := newMessageSet()
		for _, msg := range test.voteList {
			ms.AddVote(msg.hash, msg.msg)
		}
		if got := ms.TotalVotePower(test.cm); got != test.expectedPower {
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
		msg := Message{Address: common.BytesToAddress([]byte("987654321"))}

		ms := newMessageSet()
		ms.AddVote(blockHash, msg)

		if got := len(ms.Values(blockHash)); got != 1 {
			t.Fatalf("Expected 1 message, got %v", got)
		}
	})
}
