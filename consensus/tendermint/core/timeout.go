package core

import (
	"context"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/log"
	"math/big"
	"sync"
	"time"

	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/types"
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
	now := time.Now()
	switch t.Step {
	case Propose:
		ProposeTimer.Update(now.Sub(t.Start))
		ProposeBg.Add(now.Sub(t.Start).Nanoseconds())
	case Prevote:
		PrevoteTimer.UpdateSince(t.Start)
		PrevoteBg.Add(now.Sub(t.Start).Nanoseconds())
	case Precommit:
		PrecommitTimer.UpdateSince(t.Start)
		PrecommitBg.Add(now.Sub(t.Start).Nanoseconds())
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

// ///////////// On Timeout Functions ///////////////
func (c *Core) measureMetricsOnTimeOut(step uint8, r int64) {
	switch step {
	case consensus.MsgProposal:
		duration := c.timeoutPropose(r)
		ProposeTimer.Update(duration)
		ProposeBg.Add(duration.Nanoseconds())
		return
	case consensus.MsgPrevote:
		duration := c.timeoutPrevote(r)
		PrevoteTimer.Update(duration)
		PrevoteBg.Add(duration.Nanoseconds())
		return
	case consensus.MsgPrecommit:
		duration := c.timeoutPrecommit(r)
		PrecommitTimer.Update(duration)
		PrecommitBg.Add(duration.Nanoseconds())
		return
	}
}

func (c *Core) onTimeoutPropose(r int64, h *big.Int) {
	msg := types.TimeoutEvent{
		RoundWhenCalled:  r,
		HeightWhenCalled: h,
		Step:             consensus.MsgProposal,
	}
	// It's unsafe to call logTimeoutEvent here !
	c.logger.Debug("TimeoutEvent(Propose): Sent", "round", r, "height", h)
	c.measureMetricsOnTimeOut(msg.Step, r)
	c.SendEvent(msg)
}

func (c *Core) onTimeoutPrevote(r int64, h *big.Int) {
	msg := types.TimeoutEvent{
		RoundWhenCalled:  r,
		HeightWhenCalled: h,
		Step:             consensus.MsgPrevote,
	}
	c.logger.Debug("TimeoutEvent(Prevote): Sent", "round", r, "height", h)
	c.measureMetricsOnTimeOut(msg.Step, r)
	c.SendEvent(msg)
}

func (c *Core) onTimeoutPrecommit(r int64, h *big.Int) {
	msg := types.TimeoutEvent{
		RoundWhenCalled:  r,
		HeightWhenCalled: h,
		Step:             consensus.MsgPrecommit,
	}
	c.logger.Debug("TimeoutEvent(Precommit): Sent", "round", r, "height", h)
	c.measureMetricsOnTimeOut(msg.Step, r)
	c.SendEvent(msg)
}

// ///////////// Handle Timeout Functions ///////////////
func (c *Core) handleTimeoutPropose(ctx context.Context, msg types.TimeoutEvent) {
	if msg.HeightWhenCalled.Cmp(c.Height()) == 0 && msg.RoundWhenCalled == c.Round() && c.step == types.Propose {
		c.logTimeoutEvent("TimeoutEvent(Propose): Received", "Propose", msg)
		c.prevoter.SendPrevote(ctx, true)
		c.SetStep(types.Prevote)
	}
}

func (c *Core) handleTimeoutPrevote(ctx context.Context, msg types.TimeoutEvent) {
	if msg.HeightWhenCalled.Cmp(c.Height()) == 0 && msg.RoundWhenCalled == c.Round() && c.step == types.Prevote {
		c.logTimeoutEvent("TimeoutEvent(Prevote): Received", "Prevote", msg)
		c.precommiter.SendPrecommit(ctx, true)
		c.SetStep(types.Precommit)
	}
}

func (c *Core) handleTimeoutPrecommit(ctx context.Context, msg types.TimeoutEvent) {
	if msg.HeightWhenCalled.Cmp(c.Height()) == 0 && msg.RoundWhenCalled == c.Round() {
		c.logTimeoutEvent("TimeoutEvent(Precommit): Received", "Precommit", msg)
		c.StartRound(ctx, c.Round()+1)
	}
}

// ///////////// Calculate Timeout Duration Functions ///////////////
// The Timeout may need to be changed depending on the Step
func (c *Core) timeoutPropose(round int64) time.Duration {
	return types.InitialProposeTimeout + time.Duration(c.blockPeriod)*time.Second + time.Duration(round)*types.ProposeTimeoutDelta
}

func (c *Core) timeoutPrevote(round int64) time.Duration {
	return types.InitialPrevoteTimeout + time.Duration(round)*types.PrevoteTimeoutDelta
}

func (c *Core) timeoutPrecommit(round int64) time.Duration {
	return types.InitialPrecommitTimeout + time.Duration(round)*types.PrecommitTimeoutDelta
}

func (c *Core) logTimeoutEvent(message string, msgType string, timeout types.TimeoutEvent) {
	c.logger.Debug(message,
		"from", c.address.String(),
		"type", msgType,
		"currentHeight", c.Height(),
		"msgHeight", timeout.HeightWhenCalled,
		"currentRound", c.Round(),
		"msgRound", timeout.RoundWhenCalled,
		"currentStep", c.step,
		"msgStep", timeout.Step,
	)
}
