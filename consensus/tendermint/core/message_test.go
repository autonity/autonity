package core

import (
	"bytes"
	"errors"
	"github.com/clearmatics/autonity/core/types"
	"math/big"
	"reflect"
	"testing"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/config"
	"github.com/clearmatics/autonity/consensus/tendermint/validator"
	"github.com/clearmatics/autonity/rlp"
)

func TestMessageEncodeDecode(t *testing.T) {
	msg := &Message{
		Code:          msgProposal,
		Msg:           []byte{0x1},
		Address:       common.HexToAddress("0x1234567890"),
		Signature:     []byte{0x2},
		CommittedSeal: []byte{},
	}

	buf := &bytes.Buffer{}
	err := msg.EncodeRLP(buf)
	if err != nil {
		t.Fatalf("have %v, want nil", err)
	}

	s := rlp.NewStream(buf, 0)

	decMsg := &Message{}
	err = decMsg.DecodeRLP(s)
	if err != nil {
		t.Fatalf("have %v, want nil", err)
	}

	if !reflect.DeepEqual(decMsg, msg) {
		t.Errorf("Messages are not the same: have %v, want %v", decMsg, msg)
	}
}

func TestMessageString(t *testing.T) {
	msg := &Message{
		Code:    msgProposal,
		Address: common.HexToAddress("0x1234567890"),
	}

	want := "{Code: 0, Address: 0x0000000000000000000000000000001234567890}"
	if got := msg.String(); got != want {
		t.Errorf("Expected %v, got %v", want, got)
	}
}

func TestMessageFromPayload(t *testing.T) {
	t.Run("nil validator function given, nil validator returned", func(t *testing.T) {
		msg := &Message{
			Code:    msgProposal,
			Address: common.HexToAddress("0x1234567890"),
		}

		payload, _ := msg.Payload()

		decMsg := &Message{}
		val, err := decMsg.FromPayload(payload, nil, nil)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}

		if val != nil {
			t.Fatalf("validator must be nil, but got %v", val)
		}
	})

	t.Run("validator function fails, nil returned", func(t *testing.T) {
		msg := &Message{
			Code:    msgProposal,
			Address: common.HexToAddress("0x1234567890"),
		}

		payload, _ := msg.Payload()
		wantErr := errors.New("some error")

		validateFn := func(set validator.Set, data []byte, sig []byte) (common.Address, error) {
			return common.Address{}, wantErr
		}

		decMsg := &Message{}
		_, err := decMsg.FromPayload(payload, nil, validateFn)
		if err != wantErr {
			t.Fatalf("want error %v, got %v", wantErr, err)
		}
	})

	t.Run("unauthorized address given, error returned", func(t *testing.T) {
		msg := &Message{
			Code:    msgProposal,
			Address: common.HexToAddress("0x1234567890"),
		}

		payload, _ := msg.Payload()

		validateFn := func(set validator.Set, data []byte, sig []byte) (common.Address, error) {
			return common.Address{}, nil
		}

		decMsg := &Message{}
		_, err := decMsg.FromPayload(payload, nil, validateFn)
		if err != ErrUnauthorizedAddress {
			t.Fatalf("have %v, want %v", err, ErrUnauthorizedAddress)
		}
	})

	t.Run("valid params given, valid validator returned", func(t *testing.T) {
		authorizedAddress := common.HexToAddress("0x1234567890")
		msg := &Message{
			Code:    msgProposal,
			Address: authorizedAddress,
		}

		payload, _ := msg.Payload()

		val := validator.New(authorizedAddress, new(big.Int).SetUint64(1))
		valSet := validator.NewSet([]types.CommitteeMember{{Address: authorizedAddress, VotingPower: new(big.Int).SetUint64(1)}}, config.RoundRobin)
		validateFn := func(set validator.Set, data []byte, sig []byte) (common.Address, error) {
			return authorizedAddress, nil
		}

		decMsg := &Message{}
		newVal, err := decMsg.FromPayload(payload, valSet, validateFn)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}

		if !reflect.DeepEqual(val.String(), (*newVal).String()) {
			t.Errorf("Validators are not the same: have %v, want %v", newVal, val)
		}
	})
}

func TestMessageDecode(t *testing.T) {
	vote := &Vote{
		Round:             big.NewInt(1),
		Height:            big.NewInt(2),
		ProposedBlockHash: common.BytesToHash([]byte{0x1}),
	}

	payload, err := Encode(vote)
	if err != nil {
		t.Fatalf("have %v, want nil", err)
	}

	msg := &Message{
		Code:    msgProposal,
		Msg:     payload,
		Address: common.HexToAddress("0x1234567890"),
	}

	var decVote Vote
	err = msg.Decode(&decVote)
	if err != nil {
		t.Fatalf("have %v, want nil", err)
	}

	if !reflect.DeepEqual(vote.String(), decVote.String()) {
		t.Errorf("Votes are not the same: have %v, want %v", decVote, vote)
	}
}
