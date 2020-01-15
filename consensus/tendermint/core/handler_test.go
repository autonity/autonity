package core

import (
	"context"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/validator"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/rlp"
	"github.com/golang/mock/gomock"
	"golang.org/x/sync/errgroup"
	"gopkg.in/karalabe/cookiejar.v2/collections/prque"
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

	setCurrentRoundState := func(round int64, height int64, step Step) *roundMessages {
		currentState := NewRoundMessages(big.NewInt(round), big.NewInt(height))
		currentState.SetStep(step)
		return currentState
	}

	cases := []struct {
		currentState *roundMessages
		message      *Message
		outcome      error
	}{
		{
			setCurrentRoundState(1, 2, propose),
			createPrevote(1, 2),
			errFutureStepMessage,
		},
		{
			setCurrentRoundState(1, 2, propose),
			createPrevote(2, 2),
			errFutureRoundMessage,
		},
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
		{
			setCurrentRoundState(5, 2, precommit),
			createPrecommit(20, 2),
			errFutureRoundMessage,
		},
	}

	for _, testCase := range cases {
		logger := log.New("backend", "test", "id", 0)
		engine := core{
			logger:            logger,
			address:           currentValidator.GetAddress(),
			backlogs:          make(map[committee.Validator]*prque.Prque),
			curRoundMessages:  testCase.currentState,
			futureRoundChange: make(map[int64]int64),
			committeeSet:      &validatorSet{Set: validators},
			proposeTimeout:    newTimeout(propose, logger),
			prevoteTimeout:    newTimeout(prevote, logger),
			precommitTimeout:  newTimeout(precommit, logger),
		}

		err := engine.handleCheckedMsg(context.Background(), testCase.message, sender)

		if err != testCase.outcome {
			t.Fatal("unexpected handlecheckedmsg returning ",
				"err=", err, ", expecting=", testCase.outcome, " with msgCode=", testCase.message.Code)
		}

		if err != nil {
			backlogValue, _ := engine.backlogs[sender].Pop()
			msg := backlogValue.(*Message)
			if msg != testCase.message {
				t.Fatal("unexpected backlog message")
			}
		}
	}
}

func TestCoreStopDoesntPanic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	addr := common.HexToAddress("0x0123456789")

	backendMock := NewMockBackend(ctrl)
	backendMock.EXPECT().Address().AnyTimes().Return(addr)

	c := New(backendMock, nil)
	if err := c.Stop(); err != nil {
		t.Fatal(err)
	}
}

func TestCoreMultipleStopsDontPanic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	addr := common.HexToAddress("0x0123456789")

	backendMock := NewMockBackend(ctrl)
	backendMock.EXPECT().Address().AnyTimes().Return(addr)

	c := New(backendMock, nil)
	if err := c.Stop(); err != nil {
		t.Fatal(err)
	}

	if err := c.Stop(); err != nil {
		t.Fatal(err)
	}
}

func TestCoreMultipleConcurrentStopsDontPanic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	addr := common.HexToAddress("0x0123456789")

	backendMock := NewMockBackend(ctrl)
	backendMock.EXPECT().Address().AnyTimes().Return(addr)

	c := New(backendMock, nil)

	wg := errgroup.Group{}
	for i := 0; i < 10; i++ {
		wg.Go(c.Stop)
	}

	if err := wg.Wait(); err != nil {
		t.Fatal(err)
	}
}
