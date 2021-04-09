package core

import (
	"context"
	"github.com/clearmatics/autonity/log"
	"math/big"
	"sync"
	"time"
)

const (
	initialProposeTimeout   = 500 * time.Millisecond
	proposeTimeoutDelta     = 200 * time.Millisecond
	initialPrevoteTimeout   = 500 * time.Millisecond
	prevoteTimeoutDelta     = 200 * time.Millisecond
	initialPrecommitTimeout = 500 * time.Millisecond
	precommitTimeoutDelta   = 200 * time.Millisecond
)

type TimeoutEvent struct {
	roundWhenCalled  int64
	heightWhenCalled *big.Int
	// message type: msgProposal msgPrevote	msgPrecommit
	step uint8
}

type timeout struct {
	timer   *time.Timer
	started bool
	step    Step
	// start will be refreshed on each new schedule, it is used for metric collection of tendermint timeout.
	start  time.Time
	logger log.Logger
	sync.Mutex
}

func newTimeout(s Step, logger log.Logger) *timeout {
	return &timeout{
		started: false,
		step:    s,
		start:   time.Now(),
		logger:  logger,
	}
}

// runAfterTimeout() will be run in a separate go routine, so values used inside the function needs to be managed separately
func (t *timeout) scheduleTimeout(stepTimeout time.Duration, round int64, height *big.Int, runAfterTimeout func(r int64, h *big.Int)) {
	t.Lock()
	defer t.Unlock()
	t.started = true
	t.start = time.Now()
	t.timer = time.AfterFunc(stepTimeout, func() {
		runAfterTimeout(round, height)
	})
}

func (t *timeout) timerStarted() bool {
	t.Lock()
	defer t.Unlock()
	return t.started
}

func (t *timeout) stopTimer() error {
	t.Lock()
	defer t.Unlock()
	if t.started {
		if t.started = !t.timer.Stop(); t.started {
			switch t.step {
			case propose:
				return errNilPrevoteSent
			case prevote:
				return errNilPrecommitSent
			case precommit:
				return errMovedToNewRound
			}
		}
		t.measureMetricsOnStopTimer()
	}
	return nil
}

func (t *timeout) measureMetricsOnStopTimer() {
	switch t.step {
	case propose:
		tendermintProposeTimer.UpdateSince(t.start)
	case prevote:
		tendermintPrevoteTimer.UpdateSince(t.start)
	case precommit:
		tendermintPrecommitTimer.UpdateSince(t.start)
	}
}

func (t *timeout) reset(s Step) {
	err := t.stopTimer()
	if err != nil {
		t.logger.Info("cant stop timer", "err", err)
	}

	t.Lock()
	defer t.Unlock()
	t.timer = nil
	t.started = false
	t.step = s
	t.start = time.Time{}
}

/////////////// On Timeout Functions ///////////////
func (c *core) measureMetricsOnTimeOut(step uint8, r int64) {
	switch step {
	case msgProposal:
		duration := c.timeoutPropose(r)
		tendermintProposeTimer.Update(duration)
		return
	case msgPrevote:
		duration := c.timeoutPrevote(r)
		tendermintPrevoteTimer.Update(duration)
		return
	case msgPrecommit:
		duration := c.timeoutPrecommit(r)
		tendermintPrecommitTimer.Update(duration)
		return
	}
}

func (c *core) onTimeoutPropose(r int64, h *big.Int) {
	msg := TimeoutEvent{
		roundWhenCalled:  r,
		heightWhenCalled: h,
		step:             msgProposal,
	}
	// It's unsafe to call logTimeoutEvent here !
	c.logger.Debug("TimeoutEvent(Propose): Sent", "round", r, "height", h)
	c.measureMetricsOnTimeOut(msg.step, r)
	c.sendEvent(msg)
}

func (c *core) onTimeoutPrevote(r int64, h *big.Int) {
	msg := TimeoutEvent{
		roundWhenCalled:  r,
		heightWhenCalled: h,
		step:             msgPrevote,
	}
	c.logger.Debug("TimeoutEvent(Prevote): Sent", "round", r, "height", h)
	c.measureMetricsOnTimeOut(msg.step, r)
	c.sendEvent(msg)
}

func (c *core) onTimeoutPrecommit(r int64, h *big.Int) {
	msg := TimeoutEvent{
		roundWhenCalled:  r,
		heightWhenCalled: h,
		step:             msgPrecommit,
	}
	c.logger.Debug("TimeoutEvent(Precommit): Sent", "round", r, "height", h)
	c.measureMetricsOnTimeOut(msg.step, r)
	c.sendEvent(msg)
}

/////////////// Handle Timeout Functions ///////////////
func (c *core) handleTimeoutPropose(ctx context.Context, msg TimeoutEvent) {
	if msg.heightWhenCalled.Cmp(c.Height()) == 0 && msg.roundWhenCalled == c.Round() && c.step == propose {
		c.logTimeoutEvent("TimeoutEvent(Propose): Received", "Propose", msg)
		c.sendPrevote(ctx, true)
		c.setStep(prevote)
	}
}

func (c *core) handleTimeoutPrevote(ctx context.Context, msg TimeoutEvent) {
	if msg.heightWhenCalled.Cmp(c.Height()) == 0 && msg.roundWhenCalled == c.Round() && c.step == prevote {
		c.logTimeoutEvent("TimeoutEvent(Prevote): Received", "Prevote", msg)
		c.sendPrecommit(ctx, true)
		c.setStep(precommit)
	}
}

func (c *core) handleTimeoutPrecommit(ctx context.Context, msg TimeoutEvent) {

	if msg.heightWhenCalled.Cmp(c.Height()) == 0 && msg.roundWhenCalled == c.Round() {
		c.logTimeoutEvent("TimeoutEvent(Precommit): Received", "Precommit", msg)
		c.startRound(ctx, c.Round()+1)
	}
}

/////////////// Calculate Timeout Duration Functions ///////////////
// The timeout may need to be changed depending on the Step
func (c *core) timeoutPropose(round int64) time.Duration {
	return initialProposeTimeout + time.Duration(c.blockPeriod)*time.Second + time.Duration(round)*proposeTimeoutDelta
}

func (c *core) timeoutPrevote(round int64) time.Duration {
	return initialPrevoteTimeout + time.Duration(round)*prevoteTimeoutDelta
}

func (c *core) timeoutPrecommit(round int64) time.Duration {
	return initialPrecommitTimeout + time.Duration(round)*precommitTimeoutDelta
}

func (c *core) logTimeoutEvent(message string, msgType string, timeout TimeoutEvent) {

	c.logger.Debug(message,
		"from", c.address.String(),
		"type", msgType,
		"currentHeight", c.Height(),
		"msgHeight", timeout.heightWhenCalled,
		"currentRound", c.Round(),
		"msgRound", timeout.roundWhenCalled,
		"currentStep", c.step,
		"msgStep", timeout.step,
	)
}
