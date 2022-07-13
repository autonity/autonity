package core

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/messageutils"
	"github.com/autonity/autonity/consensus/tendermint/core/types"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/metrics"
	"math/big"
	"testing"
)

func TestCore_MeasureHeightRoundMetrics(t *testing.T) {
	t.Run("measure metrics of new height", func(t *testing.T) {
		c := &Core{
			address:          common.Address{},
			logger:           log.New("Core", "test", "id", 0),
			proposeTimeout:   types.NewTimeout(types.Propose, log.New("Core", "test", "id", 0)),
			prevoteTimeout:   types.NewTimeout(types.Prevote, log.New("Core", "test", "id", 0)),
			precommitTimeout: types.NewTimeout(types.Precommit, log.New("Core", "test", "id", 0)),
			round:            0,
			height:           big.NewInt(1),
		}
		c.MeasureHeightRoundMetrics(0)
		if m := metrics.Get("tendermint/height/change"); m == nil {
			t.Fatalf("test case failed.")
		}
	})

	t.Run("measure metrics of new round", func(t *testing.T) {
		c := &Core{
			address:          common.Address{},
			logger:           log.New("Core", "test", "id", 0),
			proposeTimeout:   types.NewTimeout(types.Propose, log.New("Core", "test", "id", 0)),
			prevoteTimeout:   types.NewTimeout(types.Prevote, log.New("Core", "test", "id", 0)),
			precommitTimeout: types.NewTimeout(types.Precommit, log.New("Core", "test", "id", 0)),
			round:            0,
			height:           big.NewInt(1),
		}
		c.MeasureHeightRoundMetrics(1)
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
			proposeTimeout:   types.NewTimeout(types.Propose, log.New("Core", "test", "id", 0)),
			prevoteTimeout:   types.NewTimeout(types.Prevote, log.New("Core", "test", "id", 0)),
			precommitTimeout: types.NewTimeout(types.Precommit, log.New("Core", "test", "id", 0)),
			round:            0,
			height:           big.NewInt(1),
		}
		c.measureMetricsOnTimeOut(messageutils.MsgProposal, 2)
		if m := metrics.Get("tendermint/timer/propose"); m == nil {
			t.Fatalf("test case failed.")
		}
	})

	t.Run("measure metrics on Timeout of prevote", func(t *testing.T) {
		c := &Core{
			address:          common.Address{},
			logger:           log.New("Core", "test", "id", 0),
			proposeTimeout:   types.NewTimeout(types.Propose, log.New("Core", "test", "id", 0)),
			prevoteTimeout:   types.NewTimeout(types.Prevote, log.New("Core", "test", "id", 0)),
			precommitTimeout: types.NewTimeout(types.Precommit, log.New("Core", "test", "id", 0)),
			round:            0,
			height:           big.NewInt(1),
		}
		c.measureMetricsOnTimeOut(messageutils.MsgPrevote, 2)
		if m := metrics.Get("tendermint/timer/prevote"); m == nil {
			t.Fatalf("test case failed.")
		}
	})

	t.Run("measure metrics on Timeout of precommit", func(t *testing.T) {
		c := &Core{
			address:          common.Address{},
			logger:           log.New("Core", "test", "id", 0),
			proposeTimeout:   types.NewTimeout(types.Propose, log.New("Core", "test", "id", 0)),
			prevoteTimeout:   types.NewTimeout(types.Prevote, log.New("Core", "test", "id", 0)),
			precommitTimeout: types.NewTimeout(types.Precommit, log.New("Core", "test", "id", 0)),
			round:            0,
			height:           big.NewInt(1),
		}
		c.measureMetricsOnTimeOut(messageutils.MsgPrecommit, 2)
		if m := metrics.Get("tendermint/timer/precommit"); m == nil {
			t.Fatalf("test case failed.")
		}
	})
}
