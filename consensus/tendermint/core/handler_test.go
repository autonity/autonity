package core

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/events"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/log"
	"go.uber.org/mock/gomock"
)

func TestHandleCheckedMessage(t *testing.T) {
	committeeSet, keysMap := NewTestCommitteeSetWithKeys(4)
	header := types.Header{Committee: committeeSet.Committee(), Number: common.Big1}
	currentValidator, _ := committeeSet.GetByIndex(0)
	sender, _ := committeeSet.GetByIndex(1)
	senderKey := keysMap[sender.Address]

	createPrevote := func(round int64, height int64) message.Msg {
		return message.NewPrevote(round, uint64(height), common.BytesToHash([]byte{0x1}), makeSigner(senderKey, sender.Address))
	}

	createPrecommit := func(round int64, height int64) message.Msg {
		return message.NewPrecommit(round, uint64(height), common.BytesToHash([]byte{0x1}), makeSigner(senderKey, sender.Address))
	}

	cases := []struct {
		round   int64
		height  *big.Int
		step    Step
		message message.Msg
		outcome error
		panic   bool
	}{
		{
			1,
			big.NewInt(2),
			Propose,
			createPrevote(1, 2),
			constants.ErrFutureStepMessage,
			false,
		},
		{
			1,
			big.NewInt(2),
			Propose,
			createPrevote(2, 2),
			constants.ErrFutureRoundMessage,
			false,
		},
		{
			0,
			big.NewInt(2),
			Propose,
			createPrevote(0, 3),
			//TODO(lorenzo) check if fine
			//constants.ErrFutureHeightMessage,
			nil,
			true,
		},
		{
			0,
			big.NewInt(2),
			Prevote,
			createPrevote(0, 2),
			nil,
			false,
		},
		{
			0,
			big.NewInt(2),
			Precommit,
			createPrecommit(0, 2),
			nil,
			false,
		},
		{
			0,
			big.NewInt(5),
			Precommit,
			createPrecommit(0, 10),
			//TODO(lorenzo) check if fine
			//			constants.ErrFutureHeightMessage,
			nil,
			true,
		},
		{
			5,
			big.NewInt(2),
			Precommit,
			createPrecommit(20, 2),
			constants.ErrFutureRoundMessage,
			false,
		},
	}

	for _, testCase := range cases {
		logger := log.New("backend", "test", "id", 0)
		messageMap := message.NewMap()
		engine := Core{
			logger:            logger,
			address:           currentValidator.Address,
			backlogs:          make(map[common.Address][]message.Msg),
			round:             testCase.round,
			height:            testCase.height,
			step:              testCase.step,
			futureRoundChange: make(map[int64]map[common.Address]*big.Int),
			messages:          messageMap,
			curRoundMessages:  messageMap.GetOrCreate(0),
			committee:         committeeSet,
			proposeTimeout:    NewTimeout(Propose, logger),
			prevoteTimeout:    NewTimeout(Prevote, logger),
			precommitTimeout:  NewTimeout(Precommit, logger),
		}
		engine.SetDefaultHandlers()

		func() {
			defer func() {
				r := recover()
				if r == nil && testCase.panic {
					t.Errorf("The code did not panic")
				}
				if r != nil && !testCase.panic {
					t.Errorf("Unexpected panic")
				}
			}()
			testCase.message.Validate(header.CommitteeMember)
			err := engine.handleMsg(context.Background(), testCase.message)

			if !errors.Is(err, testCase.outcome) {
				t.Fatal("unexpected handlecheckedmsg returning ",
					"err=", err, ", expecting=", testCase.outcome, " with msgCode=", testCase.message.Code())
			}

			if err != nil {
				backlogValue := engine.backlogs[sender.Address][0]
				if backlogValue != testCase.message {
					t.Fatal("unexpected backlog message")
				}
			}
		}()
	}
}

func TestHandleMsg(t *testing.T) {
	t.Run("old height message return error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		backendMock := interfaces.NewMockBackend(ctrl)
		c := &Core{
			logger:   log.New("backend", "test", "id", 0),
			backend:  backendMock,
			address:  common.HexToAddress("0x1234567890"),
			backlogs: make(map[common.Address][]message.Msg),
			step:     Propose,
			round:    1,
			height:   big.NewInt(2),
		}
		c.SetDefaultHandlers()

		prevote := message.NewPrevote(2, 1, common.BytesToHash([]byte{0x1}), defaultSigner)
		if err := c.handleMsg(context.Background(), prevote); !errors.Is(err, constants.ErrOldHeightMessage) {
			t.Fatal("errOldHeightMessage not returned")
		}
	})
}

func TestCoreStopDoesntPanic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	addr := common.HexToAddress("0x0123456789")

	backendMock := interfaces.NewMockBackend(ctrl)
	backendMock.EXPECT().Logger().AnyTimes().Return(log.Root())
	backendMock.EXPECT().Address().AnyTimes().Return(addr)

	logger := log.New("testAddress", "0x0000")
	eMux := event.NewTypeMuxSilent(nil, logger)
	sub := eMux.Subscribe(events.MessageEvent{})

	backendMock.EXPECT().Subscribe(gomock.Any()).Return(sub).MaxTimes(5)

	c := New(backendMock, nil)
	_, c.cancel = context.WithCancel(context.Background())
	c.subscribeEvents()
	c.stopped <- struct{}{}
	c.stopped <- struct{}{}
	c.stopped <- struct{}{}

	c.Stop()
}
