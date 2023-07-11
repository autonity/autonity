package message

import (
	"bytes"

	"github.com/stretchr/testify/require"
	"math/big"
	"reflect"
	"testing"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/rlp"
)

var skey, _ = crypto.GenerateKey()

func Signer(data []byte) ([]byte, error) {
	hashData := crypto.Keccak256(data)
	return crypto.Sign(hashData, skey)
}

func TestProposalEncodeDecode(t *testing.T) {
	t.Run("Valid round is positive", func(t *testing.T) {
		proposal := NewProposal(
			1,
			big.NewInt(2),
			3,
			types.NewBlockWithHeader(&types.Header{}),
			Signer)

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
			types.NewBlockWithHeader(&types.Header{}),
			Signer)

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

func TestLiteProposal(t *testing.T) {
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	proposer := crypto.PubkeyToAddress(key.PublicKey)
	h := big.NewInt(100)
	r := int64(1)
	vr := int64(-1)
	v := common.Hash{}

	liteP := &LightProposal{
		Round:      r,
		Height:     h,
		ValidRound: vr,
		Value:      v,
	}

	payload := liteP.BytesNoSignature()
	require.NotEmpty(t, payload)
	hashData := crypto.Keccak256(payload)
	sig, err := crypto.Sign(hashData, key)
	require.NoError(t, err)
	liteP.Signature = sig

	err = liteP.VerifySignature(proposer)
	require.NoError(t, err)

	require.Equal(t, h, liteP.H())
	require.Equal(t, r, liteP.R())
	require.Equal(t, vr, liteP.ValidRound)
	require.Equal(t, v, liteP.V())
	require.Equal(t, sig, liteP.Signature)

	buf := &bytes.Buffer{}
	err = liteP.EncodeRLP(buf)
	require.NoError(t, err)
	s := rlp.NewStream(buf, 0)
	decLiteP := &LightProposal{}
	err = decLiteP.DecodeRLP(s)
	require.NoError(t, err)
	require.Equal(t, liteP.H(), decLiteP.H())
	require.Equal(t, liteP.R(), decLiteP.R())
	require.Equal(t, liteP.ValidRound, decLiteP.ValidRound)
	require.Equal(t, liteP.V(), decLiteP.V())
	require.Equal(t, liteP.Signature, decLiteP.Signature)
}
