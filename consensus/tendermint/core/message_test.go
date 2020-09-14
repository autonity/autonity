package core

import (
	"bytes"
	"errors"
	"math/big"
	"reflect"
	"testing"

	"github.com/clearmatics/autonity/core/types"

	"github.com/clearmatics/autonity/common"
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

func TestMessageValidate(t *testing.T) {
	t.Run("nil validator function given, panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expect panic")
			}
		}()
		msg := createPrevote(t, common.Hash{}, 1, new(big.Int).SetUint64(26), types.CommitteeMember{})
		lastHeader := &types.Header{Number: new(big.Int).SetUint64(1)}
		payload := msg.Payload()

		decMsg := &Message{}
		err := decMsg.FromPayload(payload)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}
		val, err := decMsg.Validate(nil, lastHeader)
		if val != nil {
			t.Fatalf("validator must be nil, but got %v", val)
		}
	})

	t.Run("validate function fails, nil returned", func(t *testing.T) {
		msg := createPrevote(t, common.Hash{}, 1, new(big.Int).SetUint64(26), types.CommitteeMember{VotingPower: big.NewInt(0)})
		lastHeader := &types.Header{Number: new(big.Int).SetUint64(25)}
		payload := msg.Payload()
		wantErr := errors.New("some error")

		validateFn := func(previousHeader *types.Header, data []byte, sig []byte) (common.Address, error) {
			return common.Address{}, wantErr
		}

		decMsg := &Message{}
		if err := decMsg.FromPayload(payload); err != nil {
			t.Fatalf("have %v, want nil", err)
		}
		_, err := decMsg.Validate(validateFn, lastHeader)
		if err != wantErr {
			t.Fatalf("want error %v, got %v", wantErr, err)
		}
	})

	t.Run("not a committee member, error returned", func(t *testing.T) {
		member := types.CommitteeMember{common.HexToAddress("0x1234567890"), big.NewInt(1)}

		msg := createPrevote(t, common.Hash{}, 1, new(big.Int).SetUint64(26), member)
		payload := msg.Payload()

		validateFn := func(previousHeader *types.Header, data []byte, sig []byte) (common.Address, error) {
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
		val, err := decMsg.Validate(validateFn, lastHeader)
		if val != nil {
			t.Fatal("should return nil validator")
		}
		if err.Error() != "message received is not from a committee member: 0000000000000000000000000000001234567890" {
			t.Fatalf("bad error: %v", err)
		}
	})

	t.Run("valid params given, valid validator returned", func(t *testing.T) {
		authorizedAddress := common.HexToAddress("0x1234567890")
		msg := createPrevote(t, common.Hash{}, 1, new(big.Int).SetUint64(26),
			types.CommitteeMember{
				Address:     authorizedAddress,
				VotingPower: big.NewInt(1),
			})
		payload := msg.Payload()

		val := types.CommitteeMember{
			Address:     authorizedAddress,
			VotingPower: new(big.Int).SetUint64(1),
		}

		h := types.Header{
			Committee: types.Committee{val},
			Number:    big.NewInt(25),
		}
		validateFn := func(previousHeader *types.Header, data []byte, sig []byte) (common.Address, error) {
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
			member := types.CommitteeMember{common.HexToAddress("0x1234567890"), big.NewInt(1)}

			msg := createPrevote(t, common.Hash{}, 1, new(big.Int).SetUint64(24), member)
			payload := msg.Payload()

			validateFn := func(previousHeader *types.Header, data []byte, sig []byte) (common.Address, error) {
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
						count += 1
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
