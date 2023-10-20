package core

import (
	"context"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"math/big"
	"time"

	"github.com/autonity/autonity/consensus/tendermint/core/types"
)

// ///////////// On Timeout Functions ///////////////
func (c *Core) measureMetricsOnTimeOut(step uint8, r int64) {
	switch step {
	case message.MsgProposal:
		duration := c.timeoutPropose(r)
		types.ProposeTimer.Update(duration)
		types.ProposeBg.Add(duration.Nanoseconds())
		return
	case message.MsgPrevote:
		duration := c.timeoutPrevote(r)
		types.PrevoteTimer.Update(duration)
		types.PrevoteBg.Add(duration.Nanoseconds())
		return
	case message.MsgPrecommit:
		duration := c.timeoutPrecommit(r)
		types.PrecommitTimer.Update(duration)
		types.PrecommitBg.Add(duration.Nanoseconds())
		return
	}
}

func (c *Core) onTimeoutPropose(r int64, h *big.Int) {
	msg := types.TimeoutEvent{
		RoundWhenCalled:  r,
		HeightWhenCalled: h,
		Step:             message.MsgProposal,
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
		Step:             message.MsgPrevote,
	}
	c.logger.Debug("TimeoutEvent(Prevote): Sent", "round", r, "height", h)
	c.measureMetricsOnTimeOut(msg.Step, r)
	c.SendEvent(msg)
}

func (c *Core) onTimeoutPrecommit(r int64, h *big.Int) {
	msg := types.TimeoutEvent{
		RoundWhenCalled:  r,
		HeightWhenCalled: h,
		Step:             message.MsgPrecommit,
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
