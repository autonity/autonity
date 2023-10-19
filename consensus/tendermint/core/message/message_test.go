package message

import (
	"errors"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/core/types"
	"github.com/stretchr/testify/require"
	"math/big"
	"reflect"
	"testing"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/rlp"
)

//func TestMessageEncodeDecode(t *testing.T) {
//	msg := &Message{
//		Code:          consensus.MsgProposal,
//		Payload:       []byte{0x1},
//		Address:       common.HexToAddress("0x1234567890"),
//		Signature:     []byte{0x2},
//		CommittedSeal: []byte{},
//	}
//
//	buf := &bytes.Buffer{}
//	err := msg.EncodeRLP(buf)
//	if err != nil {
//		t.Fatalf("have %v, want nil", err)
//	}
//
//	s := rlp.NewStream(buf, 0)
//
//	decMsg := &Message{}
//	err = decMsg.DecodeRLP(s)
//	if err != nil {
//		t.Fatalf("have %v, want nil", err)
//	}
//
//	if !reflect.DeepEqual(decMsg, msg) {
//		t.Errorf("Messages are not the same: have %v, want %v", decMsg, msg)
//	}
//}

//func TestMessageVote(t *testing.T) {
//
//	payload := []byte("kokoko")
//
//	vote := Vote{
//		Round:             0,
//		Height:            big.NewInt(0x37),
//		ProposedBlockHash: common.Hash{},
//	}
//
//	msg := &MessageVote{
//		MessageBase: MessageBase{
//			Code:          0x21,
//			Payload:       payload,
//			Address:       common.Address{},
//			Signature:     nil,
//			CommittedSeal: nil,
//			Power:         big.NewInt(0x14),
//			Bytes:         nil,
//		},
//		Vote: vote,
//	}
//
//	encoded, err := rlp.EncodeToBytes(msg)
//	if err != nil {
//		panic(fmt.Sprintf("makeMsg EncodeToReader failed: %s", err))
//	}
//
//	fmt.Printf("NEW MESSAGE PAYLOAD: %s\n", hex.EncodeToString(encoded))
//
//	newMsg := &MessageVote{}
//
//	err = rlp.DecodeBytes(encoded, newMsg)
//	require.NoError(t, err)
//
//	require.Equal(t, msg, newMsg)
//
//}
//
//func TestMessageProposal(t *testing.T) {
//
//	payload := []byte("kokoko")
//
//	proposal := Proposal{
//		Round:         0,
//		Height:        big.NewInt(0x37),
//		ProposalBlock: types.NewBlockWithHeader(&types.Header{}),
//	}
//
//	msg := &MessageProposal{
//		MessageBase: MessageBase{
//			Code:          0x21,
//			Payload:       payload,
//			Address:       common.Address{},
//			Signature:     nil,
//			CommittedSeal: nil,
//			Power:         big.NewInt(0x14),
//			Bytes:         nil,
//		},
//		Proposal: proposal,
//	}
//
//	encoded, err := rlp.EncodeToBytes(msg)
//	if err != nil {
//		panic(fmt.Sprintf("makeMsg EncodeToReader failed: %s", err))
//	}
//
//	fmt.Printf("NEW MESSAGE PAYLOAD: %s\n", hex.EncodeToString(encoded))
//
//	newMsg := &MessageProposal{}
//
//	err = rlp.DecodeBytes(encoded, newMsg)
//	require.NoError(t, err)
//
//	require.Equal(t, msg, newMsg)
//
//}

func TestMessageValidate(t *testing.T) {

	t.Run("validate function fails, error returned", func(t *testing.T) {
		msg := CreatePrevote(t, common.Hash{}, 1, new(big.Int).SetUint64(26), types.CommitteeMember{VotingPower: big.NewInt(0)})
		lastHeader := &types.Header{Number: new(big.Int).SetUint64(25)}
		payload := msg.GetBytes()
		wantErr := errors.New("some error")

		validateFn := func(_ *types.Header, _ []byte, _ []byte) (common.Address, error) {
			return common.Address{}, wantErr
		}

		decMsg, err := FromBytes(payload)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}
		if err = decMsg.Validate(validateFn, lastHeader); err == nil {
			t.Fatalf("want error, nil returned")
		}
	})

	t.Run("not a committee member, error returned", func(t *testing.T) {
		member := types.CommitteeMember{Address: common.HexToAddress("0x1234567890"), VotingPower: big.NewInt(1)}

		msg := CreatePrevote(t, common.Hash{}, 1, new(big.Int).SetUint64(26), member)
		payload := msg.GetBytes()

		validateFn := func(_ *types.Header, _ []byte, _ []byte) (common.Address, error) { //nolint
			return member.Address, nil
		}
		lastHeader := &types.Header{Number: new(big.Int).SetUint64(25), Committee: []types.CommitteeMember{
			{
				Address:     common.HexToAddress("0x1234567899"),
				VotingPower: big.NewInt(2),
			},
		}}
		decMsg, err := FromBytes(payload)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := decMsg.Validate(validateFn, lastHeader); err == nil {
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
		payload := msg.GetBytes()

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

		decMsg, err := FromBytes(payload)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}

		if err := decMsg.Validate(validateFn, &h); err != nil {
			t.Fatalf("have %v, want nil", err)
		}
	})

	t.Run("incorrect previous block given, don't panic but return an error", func(t *testing.T) {
		// This test is incorrect and should be depreciated
		t.Skip("we keep panic for the time being")
		count := 0
		for i := uint64(0); i < 50; i++ {
			if i == 23 {
				continue
			}
			member := types.CommitteeMember{Address: common.HexToAddress("0x1234567890"), VotingPower: big.NewInt(1)}

			msg := CreatePrevote(t, common.Hash{}, 1, new(big.Int).SetUint64(24), member)
			payload := msg.GetBytes()

			validateFn := func(_ *types.Header, _ []byte, _ []byte) (common.Address, error) {
				return member.Address, nil
			}
			lastHeader := &types.Header{Number: new(big.Int).SetUint64(i), Committee: []types.CommitteeMember{member}}
			decMsg, err := FromBytes(payload)
			if err != nil {
				t.Fatalf("have %v, want nil", err)
			}
			func() {
				defer func() {
					if r := recover(); r != nil {
						count++
					}
				}()
				err := decMsg.Validate(validateFn, lastHeader)
				require.Error(t, err, "inconsistent message verification")
			}()
		}
		if count != 0 {
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

	payload, err := rlp.EncodeToBytes(vote)
	if err != nil {
		t.Fatalf("have %v, want nil", err)
	}

	msg := &Message{
		Code:    consensus.MsgProposal,
		Payload: payload,
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

	encodedVote, err := rlp.EncodeToBytes(&preVote)
	if err != nil {
		return
	}

	msg := &Message{
		Code:          consensus.MsgPrevote,
		Payload:       encodedVote,
		Address:       authorizedAddress,
		CommittedSeal: []byte{},
		Signature:     []byte{0x1},
		Power:         big.NewInt(1),
	}
	f.Add(msg.GetBytes()) // Use f.Add to provide a seed corpus
	f.Fuzz(func(t *testing.T, seed []byte) {
		_, err := FromBytes(seed)
		if err != nil {
			return
		}
	})
}
