package types

import (
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/log"
	"math/big"
	"sync"
	"time"
)

const (
	InitialProposeTimeout   = 500 * time.Millisecond
	ProposeTimeoutDelta     = 200 * time.Millisecond
	InitialPrevoteTimeout   = 500 * time.Millisecond
	PrevoteTimeoutDelta     = 200 * time.Millisecond
	InitialPrecommitTimeout = 500 * time.Millisecond
	PrecommitTimeoutDelta   = 200 * time.Millisecond
)

type TimeoutEvent struct {
	RoundWhenCalled  int64
	HeightWhenCalled *big.Int
	// message type: MsgProposal MsgPrevote	MsgPrecommit
	Step uint8
}

type Timeout struct {
	Timer   *time.Timer
	Started bool
	Step    Step
	// Start will be refreshed on each new schedule, it is used for metric collection of tendermint Timeout.
	Start  time.Time
	Logger log.Logger
	sync.Mutex
}

func NewTimeout(s Step, logger log.Logger) *Timeout {
	return &Timeout{
		Started: false,
		Step:    s,
		Start:   time.Now(),
		Logger:  logger,
	}
}

// runAfterTimeout() will be run in a separate go routine, so values used inside the function needs to be managed separately
func (t *Timeout) ScheduleTimeout(stepTimeout time.Duration, round int64, height *big.Int, runAfterTimeout func(r int64, h *big.Int)) {
	t.Lock()
	defer t.Unlock()
	t.Started = true
	t.Start = time.Now()
	t.Timer = time.AfterFunc(stepTimeout, func() {
		runAfterTimeout(round, height)
	})
}

func (t *Timeout) TimerStarted() bool {
	t.Lock()
	defer t.Unlock()
	return t.Started
}

func (t *Timeout) StopTimer() error {
	t.Lock()
	defer t.Unlock()
	if t.Started {
		if t.Started = !t.Timer.Stop(); t.Started {
			switch t.Step {
			case Propose:
				return constants.ErrNilPrevoteSent
			case Prevote:
				return constants.ErrNilPrecommitSent
			case Precommit:
				return constants.ErrMovedToNewRound
			}
		}
		t.MeasureMetricsOnStopTimer()
	}
	return nil
}

func (t *Timeout) MeasureMetricsOnStopTimer() {
	switch t.Step {
	case Propose:
		TendermintProposeTimer.UpdateSince(t.Start)
	case Prevote:
		TendermintPrevoteTimer.UpdateSince(t.Start)
	case Precommit:
		TendermintPrecommitTimer.UpdateSince(t.Start)
	}
}

func (t *Timeout) Reset(s Step) {
	err := t.StopTimer()
	if err != nil {
		t.Logger.Debug("Can't stop consensus timer", "err", err)
	}

	t.Lock()
	defer t.Unlock()
	t.Timer = nil
	t.Started = false
	t.Step = s
	t.Start = time.Time{}
}
