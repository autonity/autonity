package core

import (
	"context"
	"math/big"
	"sync"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/metrics"
)

func TestCore_measureMetricsOnStopTimer(t *testing.T) {

	t.Run("measure metric on stop timer of propose", func(t *testing.T) {
		tm := &Timeout{
			Timer:   nil,
			Started: true,
			Step:    Propose,
			Start:   time.Now(),
			Mutex:   sync.Mutex{},
		}
		tm.MeasureMetricsOnStopTimer()
		if m := metrics.Get("tendermint/timer/propose"); m == nil {
			t.Fatalf("test case failed.")
		}
	})

	t.Run("measure metric on stop timer of prevote", func(t *testing.T) {
		tm := &Timeout{
			Timer:   nil,
			Started: true,
			Step:    Prevote,
			Start:   time.Now(),
			Mutex:   sync.Mutex{},
		}
		tm.MeasureMetricsOnStopTimer()
		if m := metrics.Get("tendermint/timer/prevote"); m == nil {
			t.Fatalf("test case failed.")
		}
	})

	t.Run("measure metric on stop timer of precommit", func(t *testing.T) {
		tm := &Timeout{
			Timer:   nil,
			Started: true,
			Step:    Precommit,
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
		committeeSet, keys := NewTestCommitteeSetWithKeys(4)
		currentValidator, _ := committeeSet.GetByIndex(0)
		logger := log.New("backend", "test", "id", 0)
		messages := message.NewMap()
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
			step:             Prevote,
			proposeTimeout:   NewTimeout(Propose, logger),
			prevoteTimeout:   NewTimeout(Prevote, logger),
			precommitTimeout: NewTimeout(Precommit, logger),
		}
		engine.SetDefaultHandlers()
		timeoutEvent := TimeoutEvent{
			RoundWhenCalled:  1,
			HeightWhenCalled: big.NewInt(2),
			Step:             Prevote,
		}
		// should send precommit nil
		mockBackend.EXPECT().Sign(gomock.Any()).DoAndReturn(makeSigner(keys[currentValidator.Address].consensus, currentValidator.Address))
		mockBackend.EXPECT().Broadcast(gomock.Any(), gomock.Any()).Times(1).Do(
			func(c types.Committee, msg message.Msg) {
				if msg.Code() != message.PrecommitCode {
					t.Fatalf("unexpected message code, should be precommit")
				}
				if msg.Value() != (common.Hash{}) {
					t.Fatalf("not a nil vote")
				}
				if msg.R() != 1 || msg.H() != 2 {
					t.Fatalf("bad message view")
				}
			})

		engine.handleTimeoutPrevote(context.Background(), timeoutEvent)

		if engine.step != Precommit {
			t.Fatalf("should be precommit step now")
		}
	})
}

func TestHandleTimeoutPrecommit(t *testing.T) {
	t.Run("on Timeout precommit received, start new round", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		committeeSet, _ := NewTestCommitteeSetWithKeys(4)
		currentValidator, _ := committeeSet.GetByIndex(0)
		logger := log.New("backend", "test", "id", 0)
		messages := message.NewMap()
		curRoundMessages := messages.GetOrCreate(1)
		mockBackend := interfaces.NewMockBackend(ctrl)
		mockBackend.EXPECT().Post(gomock.Any()).AnyTimes()
		engine := Core{
			logger:           logger,
			backend:          mockBackend,
			address:          currentValidator.Address,
			curRoundMessages: curRoundMessages,
			messages:         messages,
			step:             Prevote,
			round:            1,
			height:           big.NewInt(2),
			committee:        committeeSet,
			proposeTimeout:   NewTimeout(Propose, logger),
			prevoteTimeout:   NewTimeout(Prevote, logger),
			precommitTimeout: NewTimeout(Precommit, logger),
		}
		engine.SetDefaultHandlers()
		timeoutEvent := TimeoutEvent{
			RoundWhenCalled:  1,
			HeightWhenCalled: big.NewInt(2),
			Step:             Precommit,
		}

		engine.handleTimeoutPrecommit(context.Background(), timeoutEvent)

		if engine.height.Uint64() != 2 || engine.round != 2 {
			t.Fatalf("should be next round")
		}

		if engine.step != Propose {
			t.Fatalf("should be propose step")
		}
	})

}

func TestOnTimeoutPrevote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBackend := interfaces.NewMockBackend(ctrl)
	messages := message.NewMap()
	curRoundMessages := messages.GetOrCreate(2)
	engine := Core{
		backend:          mockBackend,
		logger:           log.New("backend", "test", "id", 0),
		round:            2,
		height:           big.NewInt(4),
		curRoundMessages: curRoundMessages,
		messages:         messages,
		step:             Prevote,
	}
	engine.SetDefaultHandlers()
	mockBackend.EXPECT().Post(gomock.Any()).Times(1).Do(func(ev interface{}) {
		timeoutEvent, ok := ev.(TimeoutEvent)
		if !ok {
			t.Fatalf("could not cast to timeoutevent")
		}
		if timeoutEvent.RoundWhenCalled != 2 || timeoutEvent.HeightWhenCalled.Uint64() != 4 {
			t.Fatalf("bad view")
		}
		if timeoutEvent.Step != Prevote {
			t.Fatalf("bad step")
		}
	})
	engine.onTimeoutPrevote(2, big.NewInt(4))
}

func TestOnTimeoutPrecommit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBackend := interfaces.NewMockBackend(ctrl)
	messages := message.NewMap()
	curRoundMessages := messages.GetOrCreate(2)
	engine := Core{
		backend:          mockBackend,
		logger:           log.New("backend", "test", "id", 0),
		round:            2,
		height:           big.NewInt(4),
		step:             Precommit,
		curRoundMessages: curRoundMessages,
		messages:         messages,
	}
	engine.SetDefaultHandlers()
	mockBackend.EXPECT().Post(gomock.Any()).Times(1).Do(func(ev interface{}) {
		timeoutEvent, ok := ev.(TimeoutEvent)
		if !ok {
			t.Fatalf("could not cast to timeoutevent")
		}
		if timeoutEvent.RoundWhenCalled != 2 || timeoutEvent.HeightWhenCalled.Uint64() != 4 {
			t.Fatalf("bad view")
		}
		if timeoutEvent.Step != Precommit {
			t.Fatalf("bad step")
		}
	})
	engine.onTimeoutPrecommit(2, big.NewInt(4))
}
