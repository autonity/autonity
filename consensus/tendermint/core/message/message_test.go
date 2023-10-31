package message

import (
	"bytes"
	"github.com/autonity/autonity/core/types"
	"math/big"
	"testing"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/rlp"
)

func TestMessageDecode(t *testing.T) {
	t.Run("prevote", func(t *testing.T) {
		vote := NewVote[Prevote](1, 2, common.HexToHash("0x1227"), stubSigner)
		decrypted := &Prevote{}
		reader := bytes.NewReader(vote.Payload())
		if err := rlp.Decode(reader, decrypted); err != nil {
			t.Fatalf("have %v, want nil", err)
		}
		if decrypted.Value() != vote.Value() {
			t.Errorf("values are not the same: have %v, want %v", decrypted, vote)
		}
		if decrypted.H() != vote.H() {
			t.Errorf("values are not the same: have %v, want %v", decrypted, vote)
		}
		if decrypted.R() != vote.R() {
			t.Errorf("values are not the same: have %v, want %v", decrypted, vote)
		}
	})
	t.Run("precommit", func(t *testing.T) {
		vote := NewVote[Precommit](1, 2, common.HexToHash("0x1227"), stubSigner)
		decrypted := &Precommit{}
		reader := bytes.NewReader(vote.Payload())
		if err := rlp.Decode(reader, decrypted); err != nil {
			t.Fatalf("have %v, want nil", err)
		}
		if decrypted.Value() != vote.Value() {
			t.Errorf("values are not the same: have %v, want %v", decrypted, vote)
		}
		if decrypted.H() != vote.H() {
			t.Errorf("values are not the same: have %v, want %v", decrypted, vote)
		}
		if decrypted.R() != vote.R() {
			t.Errorf("values are not the same: have %v, want %v", decrypted, vote)
		}
	})
	t.Run("propose", func(t *testing.T) {
		//todo: ...
	})
}

func TestValidate(t *testing.T) {
	t.Run("invalid signature, error returned", func(t *testing.T) {
		lastHeader := &types.Header{Number: new(big.Int).SetUint64(25)}
		msg := NewVote[Prevote](1, 25, lastHeader.Hash(), stubSigner)
		if err := msg.Validate(func(_ common.Address) *types.CommitteeMember {
			return &types.CommitteeMember{}
		}); err == nil {
			t.Fatalf("want error, nil returned")
		}
	})

	t.Run("not a committee member, error returned", func(t *testing.T) {
		lastHeader := &types.Header{Number: new(big.Int).SetUint64(25), Committee: []types.CommitteeMember{
			{
				Address:     common.HexToAddress("0x1234567899"),
				VotingPower: big.NewInt(2),
			},
		}}
		member := types.CommitteeMember{Address: common.HexToAddress("0x1234567890"), VotingPower: big.NewInt(1)}

		msg := NewVote[Prevote](1, 25, lastHeader.Hash(), signer)
		payload := msg.GetBytes()

		validateFn := func(_ *types.Header, _ []byte, _ []byte) (common.Address, error) { //nolint
			return member.Address, nil
		}

		decMsg, err := FromBytes(payload)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := decMsg.Validate(validateFn, lastHeader); err == nil {
			t.Fatalf("want error, nil returned")
		}
	})
	//
	//t.Run("valid params given, valid validator returned", func(t *testing.T) {
	//	authorizedAddress := common.HexToAddress("0x1234567890")
	//	msg := CreatePrevote(t, common.Hash{}, 1, new(big.Int).SetUint64(26),
	//		types.CommitteeMember{
	//			Address:     authorizedAddress,
	//			VotingPower: big.NewInt(1),
	//		})
	//	payload := msg.GetBytes()
	//
	//	val := types.CommitteeMember{
	//		Address:     authorizedAddress,
	//		VotingPower: new(big.Int).SetUint64(1),
	//	}
	//
	//	h := types.Header{
	//		Committee: types.Committee{val},
	//		Number:    big.NewInt(25),
	//	}
	//	validateFn := func(_ *types.Header, _ []byte, _ []byte) (common.Address, error) { //nolint
	//		return authorizedAddress, nil
	//	}
	//
	//	decMsg, err := FromBytes(payload)
	//	if err != nil {
	//		t.Fatalf("have %v, want nil", err)
	//	}
	//
	//	if err := decMsg.Validate(validateFn, &h); err != nil {
	//		t.Fatalf("have %v, want nil", err)
	//	}
	//})
	//
	//t.Run("incorrect previous block given, don't panic but return an error", func(t *testing.T) {
	//	// This test is incorrect and should be depreciated
	//	t.Skip("we keep panic for the time being")
	//	count := 0
	//	for i := uint64(0); i < 50; i++ {
	//		if i == 23 {
	//			continue
	//		}
	//		member := types.CommitteeMember{Address: common.HexToAddress("0x1234567890"), VotingPower: big.NewInt(1)}
	//
	//		msg := CreatePrevote(t, common.Hash{}, 1, new(big.Int).SetUint64(24), member)
	//		payload := msg.GetBytes()
	//
	//		validateFn := func(_ *types.Header, _ []byte, _ []byte) (common.Address, error) {
	//			return member.Address, nil
	//		}
	//		lastHeader := &types.Header{Number: new(big.Int).SetUint64(i), Committee: []types.CommitteeMember{member}}
	//		decMsg, err := FromBytes(payload)
	//		if err != nil {
	//			t.Fatalf("have %v, want nil", err)
	//		}
	//		func() {
	//			defer func() {
	//				if r := recover(); r != nil {
	//					count++
	//				}
	//			}()
	//			err := decMsg.Validate(validateFn, lastHeader)
	//			require.Error(t, err, "inconsistent message verification")
	//		}()
	//	}
	//	if count != 0 {
	//		t.Fatal("panic was expected")
	//	}
	//})
}

/*
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
*/
