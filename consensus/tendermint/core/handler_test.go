package core

import (
	"context"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/rlp"
	"math/big"
	"testing"
)

func TestHandleCheckedMessage(t *testing.T) {
	validators, keysMap := newTestValidatorSetWithKeys(4)
	currentValidator := validators.GetByIndex(0)
	sender := validators.GetByIndex(1)
	senderKey := keysMap[sender.GetAddress()]

	createPrevote := func(round int64, height int64) *Message {
		vote := &Vote{
			Round:             big.NewInt(round),
			Height:            big.NewInt(height),
			ProposedBlockHash: common.BytesToHash([]byte{0x1}),
		}
		encoded, err := rlp.EncodeToBytes(&vote)
		if err != nil {
			t.Fatalf("could not encode vote")
		}
		return &Message{
			Code:    msgPrevote,
			Msg:     encoded,
			Address: sender.GetAddress(),
		}
	}

	createPrecommit := func(round int64, height int64) *Message {
		vote := &Vote{
			Round:             big.NewInt(round),
			Height:            big.NewInt(height),
			ProposedBlockHash: common.BytesToHash([]byte{0x1}),
		}
		encoded, err := rlp.EncodeToBytes(&vote)
		if err != nil {
			t.Fatalf("could not encode vote")
		}
		data := PrepareCommittedSeal(common.BytesToHash([]byte{0x1}), vote.Round, vote.Height)
		hashData := crypto.Keccak256(data)
		commitSign, err := crypto.Sign(hashData, senderKey)
		if err != nil {
			t.Fatalf("error signing")
		}
		return &Message{
			Code:          msgPrecommit,
			Msg:           encoded,
			Address:       sender.GetAddress(),
			CommittedSeal: commitSign,
		}
	}

	setCurrentRoundState := func(round int64, height int64, step Step) struct {
		round  *big.Int
		height *big.Int
		step   Step
	} {

		return struct {
			round  *big.Int
			height *big.Int
			step   Step
		}{big.NewInt(round), big.NewInt(height), step}
	}

	cases := []struct {
		currentState struct {
			round  *big.Int
			height *big.Int
			step   Step
		}
		message *Message
		outcome error
	}{
		{
			setCurrentRoundState(0, 2, propose),
			createPrevote(0, 3),
			errFutureHeightMessage,
		},
		{
			setCurrentRoundState(0, 2, prevote),
			createPrevote(0, 2),
			nil,
		},
		{
			setCurrentRoundState(0, 2, precommit),
			createPrecommit(0, 2),
			nil,
		},
		{
			setCurrentRoundState(0, 5, precommit),
			createPrecommit(0, 10),
			errFutureHeightMessage,
		},
	}

	for _, testCase := range cases {
		logger := log.New("backend", "test", "id", 0)
		engine := core{
			logger:           logger,
			address:          currentValidator.GetAddress(),
			round:            testCase.currentState.round,
			height:           testCase.currentState.height,
			step:             testCase.currentState.step,
			valSet:           &validatorSet{Set: validators},
			proposeTimeout:   newTimeout(propose, logger),
			prevoteTimeout:   newTimeout(prevote, logger),
			precommitTimeout: newTimeout(precommit, logger),
		}

		err := engine.handleCheckedMsg(context.Background(), testCase.message, sender)

		if err != testCase.outcome {
			t.Fatal("unexpected handlecheckedmsg returning ",
				"err=", err, ", expecting=", testCase.outcome, " with msgCode=", testCase.message.Code)
		}
	}

}
