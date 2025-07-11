package core

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/events"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/log"
)

func makeBogusMessageEvent(msg message.Msg, alreadyDisseminated bool) events.MessageEvent {
	return events.NewMessageEvent(msg, nil, common.Address{}, time.Now(), alreadyDisseminated)
}

type testCase struct {
	id         uint64
	round      int64
	height     *big.Int
	step       Step
	message    message.Msg
	outcome    error
	panic      bool
	shouldJail bool
}

func (tc *testCase) String() string {
	return fmt.Sprintf("%#v", tc)
}

func searchForFutureMsg(engine *Core, msg message.Msg) bool {
	messages := engine.futureRound[msg.R()]
	for _, event := range messages {
		if event.Message().Hash() == msg.Hash() {
			return true
		}
	}
	return false
}

func TestHandleMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	committeeSet, keysMap := NewTestCommitteeSetWithKeys(4)
	currentValidator, _ := committeeSet.MemberByIndex(0)
	sender, _ := committeeSet.MemberByIndex(1)
	senderKey := keysMap[sender.Address].consensus

	createPrevote := func(round int64, height int64) message.Msg {
		return message.NewPrevote(round, uint64(height), common.BytesToHash([]byte{0x1}), makeSigner(senderKey), sender, 4)
	}

	createPrecommit := func(round int64, height int64) message.Msg {
		return message.NewPrecommit(round, uint64(height), common.BytesToHash([]byte{0x1}), makeSigner(senderKey), sender, 4)
	}

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
		},
		{
			2,
			1,
			big.NewInt(2),
			Propose,
			createPrevote(1, 5),
			nil,
			true, // future height should panic
			true, // doesn't matter
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
		},
		{
			5,
			5,
			big.NewInt(2),
			Precommit,
			createPrecommit(20, 2),
			constants.ErrFutureRoundMessage,
			false,
			false,
		},
		{
			6,
			2,
			big.NewInt(2),
			Precommit,
			createPrecommit(1, 2),
			constants.ErrOldRoundMessage,
			false,
			false,
		},
		{
			7,
			1,
			big.NewInt(2),
			Propose,
			message.NewPropose(1, 2, -1, types.NewBlockWithHeader(&types.Header{}), makeSigner(senderKey), sender),
			constants.ErrNotFromProposer,
			false,
			true,
		},
	}

	for _, tc := range cases {
		logger := log.New("backend", "test", "id", 0)
		messageMap := message.NewMap()
		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(gomock.Any()).AnyTimes()
		eventCh := make(chan events.CoreEvent, EventQueueSize)
		engine := Core{
			logger:           logger,
			address:          currentValidator.Address,
			round:            tc.round,
			height:           tc.height,
			step:             tc.step,
			futureRound:      make(map[int64][]events.MessageEvent),
			futurePower:      make(map[int64]*message.AggregatedPower),
			messages:         messageMap,
			curRoundMessages: messageMap.GetOrCreate(0),
			committee:        committeeSet,
			proposeTimeout:   NewTimeout(Propose, logger),
			prevoteTimeout:   NewTimeout(Prevote, logger),
			precommitTimeout: NewTimeout(Precommit, logger),
			backend:          backendMock,
			eventCh:          eventCh,
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
			err := engine.handleMsg(context.Background(), tc.message)

			if !errors.Is(err, tc.outcome) {
				t.Log(tc.String())
				t.Fatal("unexpected behaviour, handleMsg returning", "err=", err, ", expecting=", tc.outcome)
			}

			if err != nil {
				// check if jailing is required
				shouldJail := shouldJailSigner(err)
				if tc.shouldJail != shouldJail {
					t.Log(tc.String())
					t.Fatal("unexpected behaviour, shouldJailSigner returning", "shouldJail=", shouldJail, ", expecting=", tc.shouldJail)
				}
			}
		}()
	}
}

// this test differs from the previous one because we check that the future power gets updated correctly
func TestHandleFutureRound(t *testing.T) {
	// setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	committeeSet, keysMap := NewTestCommitteeSetWithKeys(10)
	sender1, _ := committeeSet.MemberByIndex(0)
	sender2, _ := committeeSet.MemberByIndex(1)

	currentHeight := big.NewInt(1)
	currentRound := int64(0)
	logger := log.New("backend", "test", "id", 0)
	messageMap := message.NewMap()
	backendMock := interfaces.NewMockBackend(ctrl)
	backendMock.EXPECT().Post(gomock.Any()).AnyTimes()
	eventCh := make(chan events.CoreEvent, EventQueueSize)
	engine := Core{
		logger:           logger,
		address:          sender1.Address,
		round:            currentRound,
		height:           currentHeight,
		step:             Propose,
		futureRound:      make(map[int64][]events.MessageEvent),
		futurePower:      make(map[int64]*message.AggregatedPower),
		messages:         messageMap,
		curRoundMessages: messageMap.GetOrCreate(0),
		committee:        committeeSet,
		proposeTimeout:   NewTimeout(Propose, logger),
		prevoteTimeout:   NewTimeout(Prevote, logger),
		precommitTimeout: NewTimeout(Precommit, logger),
		backend:          backendMock,
		eventCh:          eventCh,
		syncState:        &SyncState{},
	}
	engine.SetDefaultHandlers()

	// handling vote
	vote := message.NewPrevote(currentRound+1, currentHeight.Uint64(), common.BytesToHash([]byte{0x1}), makeSigner(keysMap[sender2.Address].consensus), sender2, 4)
	// future round messages are forwarded right away
	mockRouter := interfaces.NewMockRouter(ctrl)
	mockRouter.EXPECT().Forward(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1)
	backendMock.EXPECT().Router().Return(mockRouter).Times(1)
	engine.handleEvent(context.Background(), makeBogusMessageEvent(vote, false))

	// check that vote was saved in the future messages and power was updated accordingly
	found := searchForFutureMsg(&engine, vote)
	require.True(t, found)
	require.Equal(t, common.Big1, engine.futurePower[vote.R()].Power())

	lastHeader := &types.Header{Number: currentHeight.Sub(currentHeight, common.Big1)}
	// same thing for future round proposal
	propose := message.NewPropose(currentRound+1, currentHeight.Uint64(), -1, generateBlock(currentHeight, lastHeader), makeSigner(keysMap[sender1.Address].consensus), sender1)
	// proposals are never disseminated in Core
	engine.handleEvent(context.Background(), makeBogusMessageEvent(propose, false))

	found = searchForFutureMsg(&engine, propose)
	require.True(t, found)
	require.Equal(t, common.Big2, engine.futurePower[propose.R()].Power())
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
