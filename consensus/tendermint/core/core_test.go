package core

import (
	"context"
	"math/big"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
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
