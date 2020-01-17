package core

import (
	"bytes"
	"math/big"
	"reflect"
	"testing"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/rlp"
)

func TestProposalEncodeDecode(t *testing.T) {
	t.Run("Valid round is positive", func(t *testing.T) {
		proposal := NewProposal(
			1,
			big.NewInt(2),
			3,
			types.NewBlockWithHeader(&types.Header{}))

		buf := &bytes.Buffer{}
		err := proposal.EncodeRLP(buf)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}

		s := rlp.NewStream(buf, 0)

		decProposal := &Proposal{}
		err = decProposal.DecodeRLP(s)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}

		if decProposal.Round != proposal.Round {
			t.Errorf("Rounds are not the same: have %v, want %v", decProposal.Round, proposal.Round)
		}

		if decProposal.Height.Uint64() != proposal.Height.Uint64() {
			t.Errorf("Heights are not the same: have %v, want %v", decProposal.Height.Uint64(), proposal.Height.Uint64())
		}

		if decProposal.ValidRound != proposal.ValidRound {
			t.Errorf("Valid Rounds are not the same: have %v, want %v", decProposal.ValidRound, proposal.ValidRound)
		}
	})

	t.Run("Valid round is negative", func(t *testing.T) {
		proposal := NewProposal(
			1,
			big.NewInt(2),
			-1,
			types.NewBlockWithHeader(&types.Header{}))

		buf := &bytes.Buffer{}
		err := proposal.EncodeRLP(buf)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}

		s := rlp.NewStream(buf, 0)

		decProposal := &Proposal{}
		err = decProposal.DecodeRLP(s)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}

		if decProposal.Round != proposal.Round {
			t.Errorf("Rounds are not the same: have %v, want %v", decProposal.Round, proposal.Round)
		}

		if decProposal.Height.Int64() != proposal.Height.Int64() {
			t.Errorf("Heights are not the same: have %v, want %v", decProposal.Height.Int64(), proposal.Height.Int64())
		}

		if decProposal.ValidRound != -1 {
			t.Errorf("Valid Rounds are not the same: have %v, want %v", decProposal.ValidRound, proposal.ValidRound)
		}
	})

}

func TestVoteEncodeDecode(t *testing.T) {
	vote := &Vote{
		Round:             1,
		Height:            big.NewInt(2),
		ProposedBlockHash: common.BytesToHash([]byte("1234567890")),
	}

	buf := &bytes.Buffer{}
	err := vote.EncodeRLP(buf)
	if err != nil {
		t.Fatalf("have %v, want nil", err)
	}

	s := rlp.NewStream(buf, 0)

	decVote := &Vote{}
	err = decVote.DecodeRLP(s)
	if err != nil {
		t.Fatalf("Expected nil, got %v", err)
	}

	if !reflect.DeepEqual(decVote, vote) {
		t.Errorf("Votes are not the same: have %v, want %v", decVote, vote)
	}
}

func TestVoteString(t *testing.T) {
	vote := &Vote{
		Round:             1,
		Height:            big.NewInt(2),
		ProposedBlockHash: common.BytesToHash([]byte("1")),
	}

	want := "{Round: 1, Height: 2 ProposedBlockHash: 0x0000000000000000000000000000000000000000000000000000000000000031}"
	has := vote.String()
	if has != want {
		t.Errorf("Vote is not stringified correctly: have %v, want %v", has, want)
	}
}
