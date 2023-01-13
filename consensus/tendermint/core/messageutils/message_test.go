package messageutils

import (
	"bytes"
	"errors"
	"github.com/autonity/autonity/core/types"
	"math/big"
	"reflect"
	"testing"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/rlp"
)

func TestMessageEncodeDecode(t *testing.T) {
	msg := &Message{
		Code:          MsgProposal,
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

func TestMessageValidate(t *testing.T) {

	t.Run("validate function fails, error returned", func(t *testing.T) {
		msg := CreatePrevote(t, common.Hash{}, 1, new(big.Int).SetUint64(26), types.CommitteeMember{VotingPower: big.NewInt(0)})
		lastHeader := &types.Header{Number: new(big.Int).SetUint64(25)}
		payload := msg.GetPayload()
		wantErr := errors.New("some error")

		validateFn := func(_ *types.Header, _ []byte, _ []byte) (common.Address, error) {
			return common.Address{}, wantErr
		}

		decMsg := &Message{}
		if err := decMsg.FromPayload(payload); err != nil {
			t.Fatalf("have %v, want nil", err)
		}
		_, err := decMsg.Validate(validateFn, lastHeader)
		if err == nil {
			t.Fatalf("want error, nil returned")
		}
	})

	t.Run("not a committee member, error returned", func(t *testing.T) {
		member := types.CommitteeMember{Address: common.HexToAddress("0x1234567890"), VotingPower: big.NewInt(1)}

		msg := CreatePrevote(t, common.Hash{}, 1, new(big.Int).SetUint64(26), member)
		payload := msg.GetPayload()

		validateFn := func(_ *types.Header, _ []byte, _ []byte) (common.Address, error) { //nolint
			return member.Address, nil
		}
		lastHeader := &types.Header{Number: new(big.Int).SetUint64(25), Committee: []types.CommitteeMember{
			{
				Address:     common.HexToAddress("0x1234567899"),
				VotingPower: big.NewInt(2),
			},
		}}
		decMsg := &Message{}
		if err := decMsg.FromPayload(payload); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		_, err := decMsg.Validate(validateFn, lastHeader)

		if err == nil {
			t.Fatalf("want error, nil returned")
		}
	})

	t.Run("valid params given, valid validator returned", func(t *testing.T) {
		authorizedAddress := common.HexToAddress("0x1234567890")
		msg := CreatePrevote(t, common.Hash{}, 1, new(big.Int).SetUint64(26),
			types.CommitteeMember{
				Address:     authorizedAddress,
				VotingPower: big.NewInt(1),
			})
		payload := msg.GetPayload()

		val := types.CommitteeMember{
			Address:     authorizedAddress,
			VotingPower: new(big.Int).SetUint64(1),
		}

		h := types.Header{
			Committee: types.Committee{val},
			Number:    big.NewInt(25),
		}
		validateFn := func(_ *types.Header, _ []byte, _ []byte) (common.Address, error) { //nolint
			return authorizedAddress, nil
		}

		decMsg := &Message{}

		if err := decMsg.FromPayload(payload); err != nil {
			t.Fatalf("have %v, want nil", err)
		}
		newVal, err := decMsg.Validate(validateFn, &h)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}
		if !reflect.DeepEqual(val.String(), (*newVal).String()) {
			t.Errorf("Validators are not the same: have %v, want %v", newVal, val)
		}
	})

	t.Run("incorrect previous block given, panic", func(t *testing.T) {
		count := 0
		for i := uint64(0); i < 50; i++ {
			if i == 23 {
				continue
			}
			member := types.CommitteeMember{Address: common.HexToAddress("0x1234567890"), VotingPower: big.NewInt(1)}

			msg := CreatePrevote(t, common.Hash{}, 1, new(big.Int).SetUint64(24), member)
			payload := msg.GetPayload()

			validateFn := func(_ *types.Header, _ []byte, _ []byte) (common.Address, error) {
				return member.Address, nil
			}
			lastHeader := &types.Header{Number: new(big.Int).SetUint64(i), Committee: []types.CommitteeMember{member}}
			decMsg := &Message{}
			if err := decMsg.FromPayload(payload); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			func() {
				defer func() {
					if r := recover(); r != nil {
						count++
					}
				}()
				_, _ = decMsg.Validate(validateFn, lastHeader)
			}()
		}
		if count != 49 {
			t.Fatal("panic was expected")
		}
	})
}

func TestMessageDecode(t *testing.T) {
	vote := &Vote{
		Round:             1,
		Height:            big.NewInt(2),
		ProposedBlockHash: common.BytesToHash([]byte{0x1}),
	}

	payload, err := Encode(vote)
	if err != nil {
		t.Fatalf("have %v, want nil", err)
	}

	msg := &Message{
		Code:    MsgProposal,
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

func FuzzFromPayload(f *testing.F) {
	authorizedAddress := common.HexToAddress("0x1234567890")
	var preVote = Vote{
		Round:             1,
		Height:            new(big.Int).SetUint64(26),
		ProposedBlockHash: common.Hash{},
	}

	encodedVote, err := Encode(&preVote)
	if err != nil {
		return
	}

	msg := &Message{
		Code:          MsgPrevote,
		Msg:           encodedVote,
		Address:       authorizedAddress,
		CommittedSeal: []byte{},
		Signature:     []byte{0x1},
		Power:         big.NewInt(1),
	}
	f.Add(msg.GetPayload()) // Use f.Add to provide a seed corpus
	f.Fuzz(func(t *testing.T, seed []byte) {
		decMsg := &Message{}
		err := decMsg.FromPayload(seed)
		if err != nil {
			return
		}
	})
}
