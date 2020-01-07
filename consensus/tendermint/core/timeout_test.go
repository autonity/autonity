package core

import (
	"context"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/validator"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/metrics"
	"github.com/clearmatics/autonity/rlp"
	"github.com/golang/mock/gomock"
	"math/big"
	"sync"
	"testing"
	"time"
)

func TestCore_measureMetricsOnStopTimer(t *testing.T) {

	t.Run("measure metric on stop timer of propose", func(t *testing.T) {
		tm := &timeout{
			timer:   nil,
			started: true,
			step:    propose,
			start:   time.Now(),
			Mutex:   sync.Mutex{},
		}
		tm.measureMetricsOnStopTimer()
		if m := metrics.Get("tendermint/timer/propose"); m == nil {
			t.Fatalf("test case failed.")
		}
	})

	t.Run("measure metric on stop timer of prevote", func(t *testing.T) {
		tm := &timeout{
			timer:   nil,
			started: true,
			step:    prevote,
			start:   time.Now(),
			Mutex:   sync.Mutex{},
		}
		tm.measureMetricsOnStopTimer()
		if m := metrics.Get("tendermint/timer/prevote"); m == nil {
			t.Fatalf("test case failed.")
		}
	})

	t.Run("measure metric on stop timer of precommit", func(t *testing.T) {
		tm := &timeout{
			timer:   nil,
			started: true,
			step:    precommit,
			start:   time.Now(),
			Mutex:   sync.Mutex{},
		}
		tm.measureMetricsOnStopTimer()
		if m := metrics.Get("tendermint/timer/precommit"); m == nil {
			t.Fatalf("test case failed.")
		}
	})
}

func TestHandleTimeoutPrevote(t *testing.T) {
	t.Run("on timeout received, send precommit nil and switch step", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		validators, _ := newTestValidatorSetWithKeys(4)
		currentValidator := validators.GetByIndex(0)
		logger := log.New("backend", "test", "id", 0)
		currentState := NewRoundState(new(big.Int).SetUint64(1), new(big.Int).SetUint64(2))
		currentState.SetStep(prevote)
		mockBackend := NewMockBackend(ctrl)
		engine := core{
			logger:             logger,
			backend:            mockBackend,
			address:            currentValidator.GetAddress(),
			roundState:         currentState,
			futureRoundsChange: make(map[int64]int64),
			valSet:             &validatorSet{Set: validators},
			proposeTimeout:     newTimeout(propose, logger),
			prevoteTimeout:     newTimeout(prevote, logger),
			precommitTimeout:   newTimeout(precommit, logger),
		}
		timeoutEvent := TimeoutEvent{
			roundWhenCalled:  1,
			heightWhenCalled: 2,
			step:             msgPrevote,
		}
		// should send precommit nil
		mockBackend.EXPECT().Sign(gomock.Any()).Times(2)
		mockBackend.EXPECT().Broadcast(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Do(
			func(ctx context.Context, valSet validator.Set, payload []byte) {
				message := new(Message)
				if err := rlp.DecodeBytes(payload, message); err != nil {
					t.Fatalf("could not decode payload")
				}
				if message.Code != msgPrecommit {
					t.Fatalf("unexpected message code, should be precommit")
				}
				precommit := new(Vote)
				if err := rlp.DecodeBytes(message.Msg, precommit); err != nil {
					t.Fatalf("could not decode precommit")
				}
				if precommit.ProposedBlockHash != (common.Hash{}) {
					t.Fatalf("not a nil vote")
				}
				if precommit.Round.Uint64() != 1 || precommit.Height.Uint64() != 2 {
					t.Fatalf("bad message view")
				}
			})

		engine.handleTimeoutPrevote(context.Background(), timeoutEvent)

		if engine.roundState.step != precommit {
			t.Fatalf("should be precommit step now")
		}
	})
}

func TestHandleTimeoutPrecommit(t *testing.T) {
	t.Run("on timeout precommit received, start new round", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		validators, _ := newTestValidatorSetWithKeys(4)
		currentValidator := validators.GetByIndex(0)
		logger := log.New("backend", "test", "id", 0)
		currentState := NewRoundState(new(big.Int).SetUint64(1), new(big.Int).SetUint64(2))
		currentState.SetStep(prevote)
		mockBackend := NewMockBackend(ctrl)
		engine := core{
			logger:                       logger,
			backend:                      mockBackend,
			address:                      currentValidator.GetAddress(),
			roundState:                   currentState,
			currentHeightOldRoundsStates: make(map[int64]*roundState),
			futureRoundsChange:           make(map[int64]int64),
			valSet:                       &validatorSet{Set: validators},
			proposeTimeout:               newTimeout(propose, logger),
			prevoteTimeout:               newTimeout(prevote, logger),
			precommitTimeout:             newTimeout(precommit, logger),
		}
		timeoutEvent := TimeoutEvent{
			roundWhenCalled:  1,
			heightWhenCalled: 2,
			step:             msgPrecommit,
		}

		block := types.NewBlockWithHeader(&types.Header{OriginalHeader: types.OriginalHeader{Number: big.NewInt(1)}})
		mockBackend.EXPECT().LastCommittedProposal().Return(block, currentValidator.GetAddress())
		engine.handleTimeoutPrecommit(context.Background(), timeoutEvent)

		if engine.roundState.height.Uint64() != 2 || engine.roundState.round.Uint64() != 2 {
			t.Fatalf("should be next round")
		}

		if engine.roundState.step != propose {
			t.Fatalf("should be propose step")
		}
	})
}

func TestOnTimeoutPrevote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBackend := NewMockBackend(ctrl)
	engine := core{
		backend:    mockBackend,
		logger:     log.New("backend", "test", "id", 0),
		roundState: NewRoundState(new(big.Int).SetUint64(2), new(big.Int).SetUint64(4)),
	}
	mockBackend.EXPECT().Post(gomock.Any()).Times(1).Do(func(ev interface{}) {
		timeoutEvent, ok := ev.(TimeoutEvent)
		if !ok {
			t.Fatalf("could not cast to timeoutevent")
		}
		if timeoutEvent.roundWhenCalled != 2 || timeoutEvent.heightWhenCalled != 4 {
			t.Fatalf("bad view")
		}
		if timeoutEvent.step != msgPrevote {
			t.Fatalf("bad step")
		}
	})
	engine.onTimeoutPrevote(2, 4)
}

func TestOnTimeoutPrecommit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBackend := NewMockBackend(ctrl)
	engine := core{
		backend:    mockBackend,
		logger:     log.New("backend", "test", "id", 0),
		roundState: NewRoundState(new(big.Int).SetUint64(2), new(big.Int).SetUint64(4)),
	}
	mockBackend.EXPECT().Post(gomock.Any()).Times(1).Do(func(ev interface{}) {
		timeoutEvent, ok := ev.(TimeoutEvent)
		if !ok {
			t.Fatalf("could not cast to timeoutevent")
		}
		if timeoutEvent.roundWhenCalled != 2 || timeoutEvent.heightWhenCalled != 4 {
			t.Fatalf("bad view")
		}
		if timeoutEvent.step != msgPrecommit {
			t.Fatalf("bad step")
		}
	})
	engine.onTimeoutPrecommit(2, 4)
}
