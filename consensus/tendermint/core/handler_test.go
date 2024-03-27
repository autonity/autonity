package core

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/events"

	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/log"
)

type testCase struct {
	id               uint64
	round            int64
	height           *big.Int
	step             Step
	message          message.Msg
	outcome          error
	panic            bool
	shouldDisconnect bool
	jailed           bool // signals if the sender should be considered as jailed
}

func (tc *testCase) String() string {
	return fmt.Sprintf("%#v", tc)
}

func TestHandleCheckedMessage(t *testing.T) {
	committeeSet, keysMap := NewTestCommitteeSetWithKeys(4)
	currentValidator, _ := committeeSet.GetByIndex(0)
	sender, _ := committeeSet.GetByIndex(1)
	senderKey := keysMap[sender.Address]

	createPrevote := func(round int64, height int64) message.Msg {
		return message.NewPrevote(round, uint64(height), common.BytesToHash([]byte{0x1}), makeSigner(senderKey, sender.Address))
	}

	createPrecommit := func(round int64, height int64) message.Msg {
		return message.NewPrecommit(round, uint64(height), common.BytesToHash([]byte{0x1}), makeSigner(senderKey, sender.Address))
	}

	// NOTE: jailed is ignored in this test case, it is useful on for the test of HandleMessage
	cases := []testCase{
		{
			0,
			1,
			big.NewInt(2),
			Propose,
			createPrevote(1, 2),
			nil,
			false,
			false,
			false,
		},
		{
			1,
			1,
			big.NewInt(2),
			Propose,
			createPrevote(2, 2),
			constants.ErrFutureRoundMessage,
			false,
			false,
			false,
		},
		{
			2,
			0,
			big.NewInt(2),
			Propose,
			createPrevote(0, 3),
			constants.ErrFutureHeightMessage,
			true,
			false,
			false,
		},
		{
			3,
			0,
			big.NewInt(2),
			Prevote,
			createPrevote(0, 2),
			nil,
			false,
			false,
			false,
		},
		{
			4,
			0,
			big.NewInt(2),
			Precommit,
			createPrecommit(0, 2),
			nil,
			false,
			false,
			false,
		},
		{
			5,
			0,
			big.NewInt(5),
			Precommit,
			createPrecommit(0, 10),
			constants.ErrFutureHeightMessage,
			true,
			false,
			false,
		},
		{
			6,
			5,
			big.NewInt(2),
			Precommit,
			createPrecommit(20, 2),
			constants.ErrFutureRoundMessage,
			false,
			false,
			false,
		},
		{
			7,
			2,
			big.NewInt(2),
			Precommit,
			createPrecommit(1, 1),
			constants.ErrOldHeightMessage,
			false,
			false,
			false,
		},

		{
			8,
			2,
			big.NewInt(2),
			PrecommitDone,
			createPrecommit(2, 2),
			constants.ErrHeightClosed,
			false,
			false,
			false,
		},
		{
			9,
			2,
			big.NewInt(2),
			Precommit,
			createPrecommit(1, 2),
			constants.ErrOldRoundMessage,
			false,
			false,
			false,
		},
	}

	for _, tc := range cases {
		logger := log.New("backend", "test", "id", 0)
		messageMap := message.NewMap()
		engine := Core{
			logger:            logger,
			address:           currentValidator.Address,
			backlogs:          make(map[common.Address][]message.Msg),
			round:             tc.round,
			height:            tc.height,
			step:              tc.step,
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
				if r == nil && tc.panic {
					t.Log(tc.String())
					t.Errorf("The code did not panic")
				}
				if r != nil && !tc.panic {
					t.Log(tc.String())
					t.Errorf("Unexpected panic")
				}
			}()
			tc.message.Validate(committeeSet.CommitteeMember)
			err := engine.handleValidMsg(context.Background(), tc.message)

			if !errors.Is(err, tc.outcome) {
				t.Log(tc.String())
				t.Fatal("unexpected behaviour, handleValidMsg returning", "err=", err, ", expecting=", tc.outcome)
			}

			if err != nil {
				// check if disconnection is required
				disconnect := shouldDisconnectSender(err)
				if tc.shouldDisconnect != disconnect {
					t.Log(tc.String())
					t.Fatal("unexpected behaviour, shouldDisconnectSender returning", "disconnect=", disconnect, ", expecting=", tc.shouldDisconnect)
				}

				if err == constants.ErrFutureRoundMessage {
					// check backlog
					backlogValue := engine.backlogs[sender.Address][0]
					if backlogValue != tc.message {
						t.Fatal("unexpected backlog message")
					}
				}
			}
		}()
	}
}

func TestHandleMsg(t *testing.T) {
	committeeSet, keysMap := NewTestCommitteeSetWithKeys(4)
	currentValidator, _ := committeeSet.GetByIndex(0)
	sender, _ := committeeSet.GetByIndex(1)
	senderKey := keysMap[sender.Address]

	cases := []testCase{
		{
			0,
			1,
			big.NewInt(2),
			Propose,
			message.NewPrevote(2, 1, common.BytesToHash([]byte{0x1}), makeSigner(senderKey, sender.Address)),
			constants.ErrOldHeightMessage,
			false,
			false,
			false,
		},
		{
			1,
			1,
			big.NewInt(2),
			Propose,
			message.NewPrevote(2, 3, common.BytesToHash([]byte{0x1}), makeSigner(senderKey, sender.Address)),
			constants.ErrFutureHeightMessage,
			false,
			false,
			false,
		},
		{
			2,
			1,
			big.NewInt(2),
			Propose,
			message.NewPrevote(1, 2, common.BytesToHash([]byte{0x1}), makeSigner(senderKey, sender.Address)),
			ErrValidatorJailed,
			false,
			false,
			true,
		},
		{
			3,
			1,
			big.NewInt(2),
			Propose,
			message.NewPrevote(1, 2, common.BytesToHash([]byte{0x1}), func(hash common.Hash) ([]byte, common.Address) {
				out, _ := crypto.Sign(append(hash[:], []byte{0xca, 0xfe}...), senderKey)
				return out, sender.Address
			}),
			message.ErrBadSignature,
			false,
			true,
			false,
		},
		{
			4,
			1,
			big.NewInt(2),
			Propose,
			message.NewPropose(1, 2, -1, types.NewBlockWithHeader(&types.Header{}), makeSigner(senderKey, sender.Address)),
			constants.ErrNotFromProposer,
			false,
			true,
			false,
		},
	}

	for _, tc := range cases {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().IsJailed(sender.Address).Return(tc.jailed).MaxTimes(1)

		logger := log.New("backend", "test", "id", 3)
		c := New(backendMock, nil, currentValidator.Address, logger)

		c.height = tc.height
		c.round = tc.round
		c.step = tc.step
		c.committee = committeeSet

		func() {
			defer func() {
				r := recover()
				if r == nil && tc.panic {
					t.Log(tc.String())
					t.Errorf("The code did not panic")
				}
				if r != nil && !tc.panic {
					t.Log(tc.String())
					t.Errorf("Unexpected panic")
				}
			}()
			tc.message.Validate(committeeSet.CommitteeMember)
			err := c.handleMsg(context.Background(), tc.message)

			if !errors.Is(err, tc.outcome) {
				t.Log(tc.String())
				t.Fatal("unexpected behaviour, handleMsg returning", "err=", err, ", expecting=", tc.outcome)
			}

			if err != nil {
				// check if disconnection is required
				disconnect := shouldDisconnectSender(err)
				if tc.shouldDisconnect != disconnect {
					t.Log(tc.String())
					t.Fatal("unexpected behaviour, shouldDisconnectSender returning", "disconnect=", disconnect, ", expecting=", tc.shouldDisconnect)
				}

				if err == constants.ErrFutureRoundMessage {
					backlogValue := c.backlogs[sender.Address][0]
					if backlogValue != tc.message {
						t.Fatal("unexpected backlog message")
					}
				}

				if err == constants.ErrFutureHeightMessage {
					backlogValue := c.backlogUntrusted[tc.message.H()][0]
					if backlogValue != tc.message {
						t.Fatal("unexpected untrusted backlog message")
					}
				}
			}
		}()
	}
}

func TestCoreStopDoesntPanic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	backendMock := interfaces.NewMockBackend(ctrl)
	eMux := event.NewTypeMuxSilent(nil, log.Root())
	sub := eMux.Subscribe(events.MessageEvent{})

	backendMock.EXPECT().Subscribe(gomock.Any()).Return(sub).MaxTimes(5)

	c := New(backendMock, nil, common.HexToAddress("0x0123456789"), log.Root())
	_, c.cancel = context.WithCancel(context.Background())
	c.subscribeEvents()
	c.stopped <- struct{}{}
	c.stopped <- struct{}{}
	c.stopped <- struct{}{}

	c.Stop()
}
