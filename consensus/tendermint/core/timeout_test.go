package core

import (
	"context"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/helpers"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	tcmessage "github.com/autonity/autonity/consensus/tendermint/core/message"
	tctypes "github.com/autonity/autonity/consensus/tendermint/core/types"

	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/metrics"
	"github.com/autonity/autonity/rlp"
	"github.com/golang/mock/gomock"
)

func TestCore_measureMetricsOnStopTimer(t *testing.T) {

	t.Run("measure metric on stop timer of propose", func(t *testing.T) {
		tm := &tctypes.Timeout{
			Timer:   nil,
			Started: true,
			Step:    tctypes.Propose,
			Start:   time.Now(),
			Mutex:   sync.Mutex{},
		}
		tm.MeasureMetricsOnStopTimer()
		if m := metrics.Get("tendermint/timer/propose"); m == nil {
			t.Fatalf("test case failed.")
		}
	})

	t.Run("measure metric on stop timer of prevote", func(t *testing.T) {
		tm := &tctypes.Timeout{
			Timer:   nil,
			Started: true,
			Step:    tctypes.Prevote,
			Start:   time.Now(),
			Mutex:   sync.Mutex{},
		}
		tm.MeasureMetricsOnStopTimer()
		if m := metrics.Get("tendermint/timer/prevote"); m == nil {
			t.Fatalf("test case failed.")
		}
	})

	t.Run("measure metric on stop timer of precommit", func(t *testing.T) {
		tm := &tctypes.Timeout{
			Timer:   nil,
			Started: true,
			Step:    tctypes.Precommit,
			Start:   time.Now(),
			Mutex:   sync.Mutex{},
		}
		tm.MeasureMetricsOnStopTimer()
		if m := metrics.Get("tendermint/timer/precommit"); m == nil {
			t.Fatalf("test case failed.")
		}
	})
}

func TestHandleTimeoutPrevote(t *testing.T) {
	t.Run("on Timeout received, send precommit nil and switch step", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		committeeSet, _ := helpers.NewTestCommitteeSetWithKeys(4)
		currentValidator, _ := committeeSet.GetByIndex(0)
		logger := log.New("backend", "test", "id", 0)
		messages := tcmessage.NewMessagesMap()
		curRoundMessages := messages.GetOrCreate(1)
		mockBackend := interfaces.NewMockBackend(ctrl)
		engine := Core{
			logger:           logger,
			backend:          mockBackend,
			address:          currentValidator.Address,
			curRoundMessages: curRoundMessages,
			messages:         messages,
			round:            1,
			height:           big.NewInt(2),
			committee:        committeeSet,
			step:             tctypes.Prevote,
			proposeTimeout:   tctypes.NewTimeout(tctypes.Propose, logger),
			prevoteTimeout:   tctypes.NewTimeout(tctypes.Prevote, logger),
			precommitTimeout: tctypes.NewTimeout(tctypes.Precommit, logger),
		}
		engine.SetDefaultHandlers()
		timeoutEvent := tctypes.TimeoutEvent{
			RoundWhenCalled:  1,
			HeightWhenCalled: big.NewInt(2),
			Step:             consensus.MsgPrevote,
		}
		// should send precommit nil
		mockBackend.EXPECT().Sign(gomock.Any()).Times(2)
		mockBackend.EXPECT().Broadcast(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Do(
			func(ctx context.Context, c types.Committee, payload []byte) {
				message := new(tcmessage.Message)
				if err := rlp.DecodeBytes(payload, message); err != nil {
					t.Fatalf("could not decode payload")
				}
				if message.Code != consensus.MsgPrecommit {
					t.Fatalf("unexpected message code, should be precommit")
				}
				precommit := new(tcmessage.Vote)
				if err := rlp.DecodeBytes(message.Payload, precommit); err != nil {
					t.Fatalf("could not decode precommit")
				}
				if precommit.ProposedBlockHash != (common.Hash{}) {
					t.Fatalf("not a nil vote")
				}
				if precommit.Round != 1 || precommit.Height.Uint64() != 2 {
					t.Fatalf("bad message view")
				}
			})

		engine.handleTimeoutPrevote(context.Background(), timeoutEvent)

		if engine.step != tctypes.Precommit {
			t.Fatalf("should be precommit step now")
		}
	})
}

func TestHandleTimeoutPrecommit(t *testing.T) {
	t.Run("on Timeout precommit received, start new round", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		committeeSet, _ := helpers.NewTestCommitteeSetWithKeys(4)
		currentValidator, _ := committeeSet.GetByIndex(0)
		logger := log.New("backend", "test", "id", 0)
		messages := tcmessage.NewMessagesMap()
		curRoundMessages := messages.GetOrCreate(1)
		mockBackend := interfaces.NewMockBackend(ctrl)
		mockBackend.EXPECT().Post(gomock.Any()).AnyTimes()
		engine := Core{
			logger:           logger,
			backend:          mockBackend,
			address:          currentValidator.Address,
			curRoundMessages: curRoundMessages,
			messages:         messages,
			step:             tctypes.Prevote,
			round:            1,
			height:           big.NewInt(2),
			committee:        committeeSet,
			proposeTimeout:   tctypes.NewTimeout(tctypes.Propose, logger),
			prevoteTimeout:   tctypes.NewTimeout(tctypes.Prevote, logger),
			precommitTimeout: tctypes.NewTimeout(tctypes.Precommit, logger),
		}
		engine.SetDefaultHandlers()
		timeoutEvent := tctypes.TimeoutEvent{
			RoundWhenCalled:  1,
			HeightWhenCalled: big.NewInt(2),
			Step:             consensus.MsgPrecommit,
		}

		engine.handleTimeoutPrecommit(context.Background(), timeoutEvent)

		if engine.height.Uint64() != 2 || engine.round != 2 {
			t.Fatalf("should be next round")
		}

		if engine.step != tctypes.Propose {
			t.Fatalf("should be propose step")
		}
	})

}

func TestOnTimeoutPrevote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBackend := interfaces.NewMockBackend(ctrl)
	messages := tcmessage.NewMessagesMap()
	curRoundMessages := messages.GetOrCreate(2)
	engine := Core{
		backend:          mockBackend,
		logger:           log.New("backend", "test", "id", 0),
		round:            2,
		height:           big.NewInt(4),
		curRoundMessages: curRoundMessages,
		messages:         messages,
		step:             tctypes.Prevote,
	}
	engine.SetDefaultHandlers()
	mockBackend.EXPECT().Post(gomock.Any()).Times(1).Do(func(ev interface{}) {
		timeoutEvent, ok := ev.(tctypes.TimeoutEvent)
		if !ok {
			t.Fatalf("could not cast to timeoutevent")
		}
		if timeoutEvent.RoundWhenCalled != 2 || timeoutEvent.HeightWhenCalled.Uint64() != 4 {
			t.Fatalf("bad view")
		}
		if timeoutEvent.Step != consensus.MsgPrevote {
			t.Fatalf("bad step")
		}
	})
	engine.onTimeoutPrevote(2, big.NewInt(4))
}

func TestOnTimeoutPrecommit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBackend := interfaces.NewMockBackend(ctrl)
	messages := tcmessage.NewMessagesMap()
	curRoundMessages := messages.GetOrCreate(2)
	engine := Core{
		backend:          mockBackend,
		logger:           log.New("backend", "test", "id", 0),
		round:            2,
		height:           big.NewInt(4),
		step:             tctypes.Precommit,
		curRoundMessages: curRoundMessages,
		messages:         messages,
	}
	engine.SetDefaultHandlers()
	mockBackend.EXPECT().Post(gomock.Any()).Times(1).Do(func(ev interface{}) {
		timeoutEvent, ok := ev.(tctypes.TimeoutEvent)
		if !ok {
			t.Fatalf("could not cast to timeoutevent")
		}
		if timeoutEvent.RoundWhenCalled != 2 || timeoutEvent.HeightWhenCalled.Uint64() != 4 {
			t.Fatalf("bad view")
		}
		if timeoutEvent.Step != consensus.MsgPrecommit {
			t.Fatalf("bad step")
		}
	})
	engine.onTimeoutPrecommit(2, big.NewInt(4))
}
