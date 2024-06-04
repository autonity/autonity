package core

import (
	"context"
	"math/big"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
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
			round:            0,
			height:           big.NewInt(1),
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
			round:            0,
			height:           big.NewInt(1),
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
			round:            0,
			height:           big.NewInt(1),
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
			round:            0,
			height:           big.NewInt(1),
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
			round:            0,
			height:           big.NewInt(1),
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
	c := New(backendMock, nil, common.Address{}, log.Root(), false)
	t.Run("SetStep", func(t *testing.T) {
		timeoutDuration := c.timeoutPropose(0)
		timeoutCallback := func(_ int64, _ *big.Int) {}
		c.proposeTimeout.ScheduleTimeout(timeoutDuration, 0, common.Big1, timeoutCallback)
		require.True(t, c.proposeTimeout.TimerStarted())

		c.SetStep(context.Background(), Propose)
		require.Equal(t, Propose, c.step)
		// set step should also stop timeouts
		require.False(t, c.proposeTimeout.TimerStarted())
	})

	t.Run("setRound", func(t *testing.T) {
		c.setRound(27)
		require.Equal(t, int64(27), c.Round())
	})

	t.Run("setHeight", func(t *testing.T) {
		c := &Core{}
		c.setHeight(new(big.Int).SetUint64(10))
		require.Equal(t, uint64(10), c.height.Uint64())
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
			step:        Propose,
			round:       1,
			height:      big.NewInt(2),
		}

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
			step:        Propose,
			round:       3,
			height:      big.NewInt(2),
		}

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
			step:        Propose,
			round:       3,
			height:      big.NewInt(2),
		}

		c.setLastHeader(&types.Header{Committee: testCommittee})

		// scenario: we just switched from round 0 --> 3. Future height messages shouldn't be processed
		c.processFuture(0, 3)
		// scenario: we just switched from round 3 --> 0 of new height. Future height messages should be processed
		c.processFuture(3, 0)
	})
	// TODO(lorenzo) need more tests for future round msgs?
	/* //TODO(lorenzo) adapt or delete this tests from backlog_test.go
	t.Run("same height, but old round", func(t *testing.T) {
		msg := message.NewPrevote(0, 1, common.Hash{}, defaultSigner)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		committeeSet := NewTestCommitteeSet(1)
		val, _ := committeeSet.GetByIndex(0)

		expected := backlogMessageEvent{
			msg: msg,
		}

		evChan := make(chan any, 1)

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(expected).Do(func(ev any) {
			evChan <- ev
		})

		c := &Core{
			logger:   log.New("backend", "test", "id", 0),
			backend:  backendMock,
			address:  common.HexToAddress("0x1234567890"),
			backlogs: make(map[common.Address][]message.Msg),
			round:    1,
			step:     Prevote,
			height:   big.NewInt(1),
		}

		c.setLastHeader(&types.Header{Committee: committeeSet.Committee()})

		c.storeBacklog(msg, val.Address)
		c.processBacklog()

		timeout := time.NewTimer(2 * time.Second)
		select {
		case ev := <-evChan:
			e, ok := ev.(backlogMessageEvent)
			if !ok {
				t.Errorf("unexpected event comes: %v", reflect.TypeOf(ev))
			}
			if e.msg.Code() != msg.Code() {
				t.Errorf("message code mismatch: have %v, want %v", e.msg.Code(), msg.Code())
			}
		case <-timeout.C:
			t.Error("unexpected Timeout occurs")
		}
	})

	t.Run("future height message are not processed", func(t *testing.T) {
		msg := message.NewPrevote(2, 4, common.Hash{}, defaultSigner)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(gomock.Any()).Times(0)

		committeeSet := NewTestCommitteeSet(2)
		val, _ := committeeSet.GetByIndex(0)

		c := &Core{
			logger:   log.New("backend", "test", "id", 0),
			backend:  backendMock,
			address:  common.HexToAddress("0x1234567890"),
			backlogs: make(map[common.Address][]message.Msg),
			round:    2,
			height:   big.NewInt(3),
		}

		c.setLastHeader(&types.Header{Committee: committeeSet.Committee()})

		c.storeBacklog(msg, val.Address)
		c.processBacklog()
	})

	t.Run("future height message are processed when height change", func(t *testing.T) {
		msg := message.NewPrevote(2, 4, common.Hash{}, defaultSigner)
		msg2 := message.NewPrecommit(2, 4, common.Hash{}, defaultSigner)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(gomock.Any()).Times(0)

		committeeSet := NewTestCommitteeSet(2)
		val, _ := committeeSet.GetByIndex(0)

		c := &Core{
			logger:           log.New("backend", "test", "id", 0),
			backend:          backendMock,
			address:          common.HexToAddress("0x1234567890"),
			backlogs:         make(map[common.Address][]message.Msg),
			round:            2,
			height:           big.NewInt(3),
			proposeTimeout:   NewTimeout(Propose, log.New("ProposeTimeout")),
			prevoteTimeout:   NewTimeout(Prevote, log.New("PrevoteTimeout")),
			precommitTimeout: NewTimeout(Precommit, log.New("PrecommitTimeout")),
			committee:        committeeSet,
			messages:         message.NewMap(),
		}
		c.curRoundMessages = c.messages.GetOrCreate(2)

		c.setLastHeader(&types.Header{Committee: committeeSet.Committee()})

		c.storeBacklog(msg, val.Address)
		c.storeBacklog(msg2, val.Address)
		c.SetStep(context.Background(), Prevote)
		c.processBacklog()
		c.setHeight(big.NewInt(4))

		backendMock.EXPECT().Post(gomock.Any()).Times(2)
		c.SetStep(context.Background(), Prevote)
		c.processBacklog()
		timeout := time.NewTimer(2 * time.Second)
		<-timeout.C
	})

	t.Run("untrusted messages are processed when height change", func(t *testing.T) {
		msg := message.NewPrevote(2, 4, common.Hash{}, defaultSigner)
		msg2 := message.NewPrecommit(2, 4, common.Hash{}, defaultSigner)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)

		committeeSet := NewTestCommitteeSet(2)

		c := &Core{
			logger:           log.New("backend", "test", "id", 0),
			backend:          backendMock,
			address:          common.HexToAddress("0x1234567890"),
			backlogs:         make(map[common.Address][]message.Msg),
			backlogUntrusted: map[uint64][]message.Msg{},
			round:            2,
			height:           big.NewInt(3),
			proposeTimeout:   NewTimeout(Propose, log.New("ProposeTimeout")),
			prevoteTimeout:   NewTimeout(Prevote, log.New("PrevoteTimeout")),
			precommitTimeout: NewTimeout(Precommit, log.New("PrecommitTimeout")),
			committee:        committeeSet,
			messages:         message.NewMap(),
		}
		c.curRoundMessages = c.messages.GetOrCreate(2)

		c.setLastHeader(&types.Header{Committee: committeeSet.Committee()})

		backendMock.EXPECT().Post(gomock.Any()).Times(0)
		c.storeFutureMessage(msg)
		c.storeFutureMessage(msg2)
		c.SetStep(context.Background(), Prevote)
		c.processBacklog()
		c.setHeight(big.NewInt(4))

		backendMock.EXPECT().Post(gomock.Any()).Times(2)
		c.SetStep(context.Background(), Prevote)

		backendMock.EXPECT().Post(gomock.Any()).Times(0)
		c.processBacklog()
		<-time.NewTimer(2 * time.Second).C
	})

	t.Run("future round message are processed when round change", func(t *testing.T) {
		msg := message.NewPrevote(2, 4, common.Hash{}, defaultSigner)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(gomock.Any()).Times(0)
		committeeSet := NewTestCommitteeSet(2)
		val, err := committeeSet.GetByIndex(0)
		require.NoError(t, err)

		c := New(backendMock, nil, common.HexToAddress("0x1234567890"), log.Root(), false)
		c.setRound(1)
		c.setHeight(big.NewInt(4))
		c.step = Prevote
		c.setLastHeader(&types.Header{Committee: committeeSet.Committee()})

		c.storeBacklog(msg, val.Address)
		c.processBacklog()
		backendMock.EXPECT().Post(gomock.Any()).Times(1)
		c.setRound(2)
		c.processBacklog()
		timeout := time.NewTimer(2 * time.Second)
		<-timeout.C
	})*/

}
