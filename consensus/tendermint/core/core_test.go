package core

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/metrics"
	"math/big"
	"testing"
)

func TestCore_MeasureHeightRoundMetrics(t *testing.T) {
	t.Run("measure metrics of new height", func(t *testing.T) {
		c := &core{
			address:          common.Address{},
			logger:           log.New("core", "test", "id", 0),
			proposeTimeout:   newTimeout(propose, log.New("core", "test", "id", 0)),
			prevoteTimeout:   newTimeout(prevote, log.New("core", "test", "id", 0)),
			precommitTimeout: newTimeout(precommit, log.New("core", "test", "id", 0)),
			round:            0,
			height:           big.NewInt(1),
		}
		c.measureHeightRoundMetrics(0)
		if m := metrics.Get("tendermint/height/change"); m == nil {
			t.Fatalf("test case failed.")
		}
	})

	t.Run("measure metrics of new round", func(t *testing.T) {
		c := &core{
			address:          common.Address{},
			logger:           log.New("core", "test", "id", 0),
			proposeTimeout:   newTimeout(propose, log.New("core", "test", "id", 0)),
			prevoteTimeout:   newTimeout(prevote, log.New("core", "test", "id", 0)),
			precommitTimeout: newTimeout(precommit, log.New("core", "test", "id", 0)),
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
	t.Run("measure metrics on timeout of propose", func(t *testing.T) {
		c := &core{
			address:          common.Address{},
			logger:           log.New("core", "test", "id", 0),
			proposeTimeout:   newTimeout(propose, log.New("core", "test", "id", 0)),
			prevoteTimeout:   newTimeout(prevote, log.New("core", "test", "id", 0)),
			precommitTimeout: newTimeout(precommit, log.New("core", "test", "id", 0)),
			round:            0,
			height:           big.NewInt(1),
		}
		c.measureMetricsOnTimeOut(msgProposal, 2)
		if m := metrics.Get("tendermint/timer/propose"); m == nil {
			t.Fatalf("test case failed.")
		}
	})

	t.Run("measure metrics on timeout of prevote", func(t *testing.T) {
		c := &core{
			address:          common.Address{},
			logger:           log.New("core", "test", "id", 0),
			proposeTimeout:   newTimeout(propose, log.New("core", "test", "id", 0)),
			prevoteTimeout:   newTimeout(prevote, log.New("core", "test", "id", 0)),
			precommitTimeout: newTimeout(precommit, log.New("core", "test", "id", 0)),
			round:            0,
			height:           big.NewInt(1),
		}
		c.measureMetricsOnTimeOut(msgPrevote, 2)
		if m := metrics.Get("tendermint/timer/prevote"); m == nil {
			t.Fatalf("test case failed.")
		}
	})

	t.Run("measure metrics on timeout of precommit", func(t *testing.T) {
		c := &core{
			address:          common.Address{},
			logger:           log.New("core", "test", "id", 0),
			proposeTimeout:   newTimeout(propose, log.New("core", "test", "id", 0)),
			prevoteTimeout:   newTimeout(prevote, log.New("core", "test", "id", 0)),
			precommitTimeout: newTimeout(precommit, log.New("core", "test", "id", 0)),
			round:            0,
			height:           big.NewInt(1),
		}
		c.measureMetricsOnTimeOut(msgPrecommit, 2)
		if m := metrics.Get("tendermint/timer/precommit"); m == nil {
			t.Fatalf("test case failed.")
		}
	})
}
