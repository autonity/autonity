package core

import (
	"context"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/helpers"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/messageutils"
	"github.com/autonity/autonity/consensus/tendermint/core/types"
	"github.com/influxdata/influxdb/pkg/deep"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/events"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/rlp"
	"github.com/golang/mock/gomock"
)

func TestHandleCheckedMessage(t *testing.T) {
	committeeSet, keysMap := helpers.NewTestCommitteeSetWithKeys(4)
	currentValidator, _ := committeeSet.GetByIndex(0)
	sender, _ := committeeSet.GetByIndex(1)
	senderKey := keysMap[sender.Address]

	createPrevote := func(round int64, height int64) *messageutils.Message {
		vote := &messageutils.Vote{
			Round:             round,
			Height:            big.NewInt(height),
			ProposedBlockHash: common.BytesToHash([]byte{0x1}),
		}
		encoded, err := rlp.EncodeToBytes(&vote)
		if err != nil {
			t.Fatalf("could not encode vote")
		}
		return &messageutils.Message{
			Code:         consensus.MsgPrevote,
			TbftMsgBytes: encoded,
			Address:      sender.Address,
			Power:        common.Big1,
		}
	}

	createPrecommit := func(round int64, height int64) *messageutils.Message {
		vote := &messageutils.Vote{
			Round:             round,
			Height:            big.NewInt(height),
			ProposedBlockHash: common.BytesToHash([]byte{0x1}),
		}
		encoded, err := rlp.EncodeToBytes(&vote)
		if err != nil {
			t.Fatalf("could not encode vote")
		}
		data := helpers.PrepareCommittedSeal(common.BytesToHash([]byte{0x1}), vote.Round, vote.Height)
		hashData := crypto.Keccak256(data)
		commitSign, err := crypto.Sign(hashData, senderKey)
		if err != nil {
			t.Fatalf("error signing")
		}
		return &messageutils.Message{
			Code:          consensus.MsgPrecommit,
			TbftMsgBytes:  encoded,
			Address:       sender.Address,
			CommittedSeal: commitSign,
			Power:         common.Big1,
		}
	}

	cases := []struct {
		round   int64
		height  *big.Int
		step    types.Step
		message *messageutils.Message
		outcome error
		panic   bool
	}{
		{
			1,
			big.NewInt(2),
			types.Propose,
			createPrevote(1, 2),
			constants.ErrFutureStepMessage,
			false,
		},
		{
			1,
			big.NewInt(2),
			types.Propose,
			createPrevote(2, 2),
			constants.ErrFutureRoundMessage,
			false,
		},
		{
			0,
			big.NewInt(2),
			types.Propose,
			createPrevote(0, 3),
			constants.ErrFutureHeightMessage,
			true,
		},
		{
			0,
			big.NewInt(2),
			types.Prevote,
			createPrevote(0, 2),
			nil,
			false,
		},
		{
			0,
			big.NewInt(2),
			types.Precommit,
			createPrecommit(0, 2),
			nil,
			false,
		},
		{
			0,
			big.NewInt(5),
			types.Precommit,
			createPrecommit(0, 10),
			constants.ErrFutureHeightMessage,
			true,
		},
		{
			5,
			big.NewInt(2),
			types.Precommit,
			createPrecommit(20, 2),
			constants.ErrFutureRoundMessage,
			false,
		},
	}

	for _, testCase := range cases {
		logger := log.New("backend", "test", "id", 0)
		messageMap := messageutils.NewMessagesMap()
		engine := Core{
			logger:            logger,
			address:           currentValidator.Address,
			backlogs:          make(map[common.Address][]*messageutils.Message),
			round:             testCase.round,
			height:            testCase.height,
			step:              testCase.step,
			futureRoundChange: make(map[int64]map[common.Address]*big.Int),
			messages:          messageMap,
			curRoundMessages:  messageMap.GetOrCreate(0),
			committee:         committeeSet,
			proposeTimeout:    types.NewTimeout(types.Propose, logger),
			prevoteTimeout:    types.NewTimeout(types.Prevote, logger),
			precommitTimeout:  types.NewTimeout(types.Precommit, logger),
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

			err := engine.handleCheckedMsg(context.Background(), testCase.message)

			if err != testCase.outcome {
				t.Fatal("unexpected handlecheckedmsg returning ",
					"err=", err, ", expecting=", testCase.outcome, " with msgCode=", testCase.message.Code)
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
			backlogs: make(map[common.Address][]*messageutils.Message),
			step:     types.Propose,
			round:    1,
			height:   big.NewInt(2),
		}
		c.SetDefaultHandlers()
		vote := &messageutils.Vote{
			Round:             2,
			Height:            big.NewInt(1),
			ProposedBlockHash: common.BytesToHash([]byte{0x1}),
		}
		payload, err := rlp.EncodeToBytes(vote)
		require.NoError(t, err)
		msg := &messageutils.Message{
			Code:         consensus.MsgPrevote,
			TbftMsgBytes: payload,
			DecodedMsg:   vote,
			Address:      common.Address{},
		}

		if err := c.handleMsg(context.Background(), msg); err != constants.ErrOldHeightMessage {
			t.Fatal("errOldHeightMessage not returned")
		}
	})

	t.Run("future height message return error but are saved", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		backendMock := interfaces.NewMockBackend(ctrl)
		c := &Core{
			logger:           log.New("backend", "test", "id", 0),
			backend:          backendMock,
			address:          common.HexToAddress("0x1234567890"),
			backlogs:         make(map[common.Address][]*messageutils.Message),
			backlogUnchecked: map[uint64][]*messageutils.Message{},
			step:             types.Propose,
			round:            1,
			height:           big.NewInt(2),
		}
		c.SetDefaultHandlers()
		vote := &messageutils.Vote{
			Round:             2,
			Height:            big.NewInt(3),
			ProposedBlockHash: common.BytesToHash([]byte{0x1}),
		}
		payload, err := rlp.EncodeToBytes(vote)
		require.NoError(t, err)
		msg := &messageutils.Message{
			Code:         consensus.MsgPrevote,
			TbftMsgBytes: payload,
			DecodedMsg:   vote,
			Address:      common.Address{},
		}

		if err := c.handleMsg(context.Background(), msg); err != constants.ErrFutureHeightMessage {
			t.Fatal("errFutureHeightMessage not returned")
		}
		if backlog, ok := c.backlogUnchecked[3]; !(ok && len(backlog) > 0 && deep.Equal(backlog[0], msg)) {
			t.Fatal("future message not saved in the untrusted buffer")
		}
	})
}

func TestCoreStopDoesntPanic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	addr := common.HexToAddress("0x0123456789")

	backendMock := interfaces.NewMockBackend(ctrl)
	backendMock.EXPECT().Address().AnyTimes().Return(addr)

	logger := log.New("testAddress", "0x0000")
	eMux := event.NewTypeMuxSilent(nil, logger)
	sub := eMux.Subscribe(events.MessageEvent{})

	backendMock.EXPECT().Subscribe(gomock.Any()).Return(sub).MaxTimes(5)

	c := New(backendMock)
	_, c.cancel = context.WithCancel(context.Background())
	c.subscribeEvents()
	c.stopped <- struct{}{}
	c.stopped <- struct{}{}
	c.stopped <- struct{}{}

	c.Stop()
}
