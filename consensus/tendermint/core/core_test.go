package core

import (
	"context"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/stretchr/testify/require"
	"math/big"
	"reflect"
	"time"

	"math/rand"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/metrics"
)

func TestCore_MeasureHeightRoundMetrics(t *testing.T) {
	t.Run("measure metrics of new height", func(t *testing.T) {
		c := &Core{
			address:          common.Address{},
			logger:           log.New("Core", "test", "id", 0),
			proposeTimeout:   NewTimeout(Propose, log.New("Core", "test", "id", 0)),
			prevoteTimeout:   NewTimeout(Prevote, log.New("Core", "test", "id", 0)),
			precommitTimeout: NewTimeout(Precommit, log.New("Core", "test", "id", 0)),
			roundsState:      newTendermintState(log.New(), nil, nil),
		}
		c.measureHeightRoundMetrics(0)
		if m := metrics.Get("tendermint/height/change"); m == nil {
			t.Fatalf("test case failed.")
		}
	})

	t.Run("measure metrics of new round", func(t *testing.T) {
		c := &Core{
			address:          common.Address{},
			logger:           log.New("Core", "test", "id", 0),
			proposeTimeout:   NewTimeout(Propose, log.New("Core", "test", "id", 0)),
			prevoteTimeout:   NewTimeout(Prevote, log.New("Core", "test", "id", 0)),
			precommitTimeout: NewTimeout(Precommit, log.New("Core", "test", "id", 0)),
			roundsState:      newTendermintState(log.New(), nil, nil),
		}
		c.measureHeightRoundMetrics(1)
		if m := metrics.Get("tendermint/round/change"); m == nil {
			t.Fatalf("test case failed.")
		}
	})
}

func TestCore_measureMetricsOnTimeOut(t *testing.T) {
	t.Run("measure metrics on Timeout of propose", func(t *testing.T) {
		c := &Core{
			address:          common.Address{},
			logger:           log.New("Core", "test", "id", 0),
			proposeTimeout:   NewTimeout(Propose, log.New("Core", "test", "id", 0)),
			prevoteTimeout:   NewTimeout(Prevote, log.New("Core", "test", "id", 0)),
			precommitTimeout: NewTimeout(Precommit, log.New("Core", "test", "id", 0)),
			roundsState:      newTendermintState(log.New(), nil, nil),
		}
		c.measureMetricsOnTimeOut(Propose, 2)
		if m := metrics.Get("tendermint/timer/propose"); m == nil {
			t.Fatalf("test case failed.")
		}
	})

	t.Run("measure metrics on Timeout of prevote", func(t *testing.T) {
		c := &Core{
			address:          common.Address{},
			logger:           log.New("Core", "test", "id", 0),
			proposeTimeout:   NewTimeout(Propose, log.New("Core", "test", "id", 0)),
			prevoteTimeout:   NewTimeout(Prevote, log.New("Core", "test", "id", 0)),
			precommitTimeout: NewTimeout(Precommit, log.New("Core", "test", "id", 0)),
			roundsState:      newTendermintState(log.New(), nil, nil),
		}
		c.measureMetricsOnTimeOut(Prevote, 2)
		if m := metrics.Get("tendermint/timer/prevote"); m == nil {
			t.Fatalf("test case failed.")
		}
	})

	t.Run("measure metrics on Timeout of precommit", func(t *testing.T) {
		c := &Core{
			address:          common.Address{},
			logger:           log.New("Core", "test", "id", 0),
			proposeTimeout:   NewTimeout(Propose, log.New("Core", "test", "id", 0)),
			prevoteTimeout:   NewTimeout(Prevote, log.New("Core", "test", "id", 0)),
			precommitTimeout: NewTimeout(Precommit, log.New("Core", "test", "id", 0)),
			roundsState:      newTendermintState(log.New(), nil, nil),
		}
		c.measureMetricsOnTimeOut(Precommit, 2)
		if m := metrics.Get("tendermint/timer/precommit"); m == nil {
			t.Fatalf("test case failed.")
		}
	})
}

func TestCore_Setters(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	backendMock := interfaces.NewMockBackend(ctrl)
	c := New(backendMock, nil, common.Address{}, log.Root(), false, nil)
	t.Run("SetStep", func(t *testing.T) {
		timeoutDuration := c.timeoutPropose(0)
		timeoutCallback := func(_ int64, _ *big.Int) {}
		c.proposeTimeout.ScheduleTimeout(timeoutDuration, 0, common.Big1, timeoutCallback)
		require.True(t, c.proposeTimeout.TimerStarted())

		c.roundsState = newTendermintState(log.New(), nil, nil)
		c.SetStep(context.Background(), Propose)
		require.Equal(t, Propose, c.Step())
		// set step should also stop timeouts
		require.False(t, c.proposeTimeout.TimerStarted())
	})

	t.Run("SetRound", func(t *testing.T) {
		c.SetRound(27)
		require.Equal(t, int64(27), c.Round())
	})

	t.Run("SetHeight", func(t *testing.T) {
		c := &Core{roundsState: newTendermintState(log.New(), nil, nil)}
		c.SetHeight(new(big.Int).SetUint64(10))
		require.Equal(t, uint64(10), c.Height().Uint64())
	})

	t.Run("setCommitteeSet", func(t *testing.T) {
		c := &Core{}
		committeeSizeAndMaxRound := maxSize
		committeeSet, _ := prepareCommittee(t, committeeSizeAndMaxRound)
		c.setCommitteeSet(committeeSet)
		require.Equal(t, committeeSet, c.CommitteeSet())
	})

	t.Run("setLastHeader", func(t *testing.T) {
		c := &Core{}
		prevHeight := big.NewInt(int64(rand.Intn(100) + 1)) //nolint
		prevBlock := generateBlock(prevHeight)
		c.setLastHeader(prevBlock.Header())
		require.Equal(t, prevBlock.Header(), c.LastHeader())
	})
}

// future round message processing
func TestProcessFuture(t *testing.T) {
	t.Run("future round msg is processed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		msg := message.NewPropose(1, 2, 1, types.NewBlockWithHeader(&types.Header{}), defaultSigner, testCommitteeMember)

		expected := backlogMessageEvent{
			msg: msg,
		}

		evChan := make(chan any, 1)

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(expected).Do(func(ev any) {
			evChan <- ev
		})

		c := &Core{
			logger:      log.New("backend", "test", "id", 0),
			backend:     backendMock,
			address:     common.HexToAddress("0x1234567890"),
			futureRound: make(map[int64][]message.Msg),
			futurePower: make(map[int64]*message.AggregatedPower),
			roundsState: newTendermintState(log.New(), nil, nil),
		}
		c.UpdateStep(Propose)
		c.SetRound(1)
		c.SetHeight(common.Big2)

		c.setLastHeader(&types.Header{Committee: testCommittee})

		c.futureRound[msg.R()] = append(c.futureRound[msg.R()], msg)
		c.processFuture(0, 1) // scenario: we just switched from round 0 --> 1

		timeout := time.NewTimer(2 * time.Second)
		select {
		case ev := <-evChan:
			e, ok := ev.(backlogMessageEvent)
			if !ok {
				t.Errorf("unexpected event comes: %v", reflect.TypeOf(ev))
			}
			if e.msg.Hash() != msg.Hash() {
				t.Errorf("message hash mismatch: have %v, want %v", e.msg.Hash(), msg.Hash())
			}
		case <-timeout.C:
			t.Error("unexpected Timeout occurs")
		}
	})
	t.Run("future round messages are processed even if we skip multiple rounds", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		msg := message.NewPropose(1, 2, 1, types.NewBlockWithHeader(&types.Header{}), defaultSigner, testCommitteeMember)

		expected := backlogMessageEvent{
			msg: msg,
		}

		evChan := make(chan any, 1)

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(expected).Do(func(ev any) {
			evChan <- ev
		})

		c := &Core{
			logger:      log.New("backend", "test", "id", 0),
			backend:     backendMock,
			address:     common.HexToAddress("0x1234567890"),
			futureRound: make(map[int64][]message.Msg),
			futurePower: make(map[int64]*message.AggregatedPower),
			roundsState: newTendermintState(log.New(), nil, nil),
		}
		c.UpdateStep(Propose)
		c.SetRound(3)
		c.SetHeight(common.Big2)
		c.setLastHeader(&types.Header{Committee: testCommittee})

		c.futureRound[msg.R()] = append(c.futureRound[msg.R()], msg)
		c.processFuture(0, 3) // scenario: we just switched from round 0 --> 3

		timeout := time.NewTimer(2 * time.Second)
		select {
		case ev := <-evChan:
			e, ok := ev.(backlogMessageEvent)
			if !ok {
				t.Errorf("unexpected event comes: %v", reflect.TypeOf(ev))
			}
			if e.msg.Hash() != msg.Hash() {
				t.Errorf("message hash mismatch: have %v, want %v", e.msg.Hash(), msg.Hash())
			}
		case <-timeout.C:
			t.Error("unexpected Timeout occurs")
		}
	})
	t.Run("future height messages are processed only when actually switching heights", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer waitForExpects(ctrl)

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().ProcessFutureMsgs(uint64(2)).Times(1) // should execute only once (when we switch height)

		c := &Core{
			logger:      log.New("backend", "test", "id", 0),
			backend:     backendMock,
			address:     common.HexToAddress("0x1234567890"),
			futureRound: make(map[int64][]message.Msg),
			futurePower: make(map[int64]*message.AggregatedPower),
			roundsState: newTendermintState(log.New(), nil, nil),
		}
		c.UpdateStep(Propose)
		c.SetRound(3)
		c.SetHeight(common.Big2)
		c.setLastHeader(&types.Header{Committee: testCommittee})

		// scenario: we just switched from round 0 --> 3. Future height messages shouldn't be processed
		c.processFuture(0, 3)
		// scenario: we just switched from round 3 --> 0 of new height. Future height messages should be processed
		c.processFuture(3, 0)
	})
}
