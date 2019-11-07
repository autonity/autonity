package core

import (
	"github.com/clearmatics/autonity/metrics"
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
		time.Sleep(1)
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
		time.Sleep(1)
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
		time.Sleep(1)
		tm.measureMetricsOnStopTimer()
		if m := metrics.Get("tendermint/timer/precommit"); m == nil {
			t.Fatalf("test case failed.")
		}
	})
}
