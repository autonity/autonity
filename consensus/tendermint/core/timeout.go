package core

import (
	"math/big"
	"sync"
	"time"

	"github.com/clearmatics/autonity/common"
)

const (
	initialProposeTimeout   = 10 * time.Second
	initialPrevoteTimeout   = 10 * time.Second
	initialPrecommitTimeout = 10 * time.Second
)

type timeoutEvent struct {
	roundWhenCalled  int64
	heightWhenCalled int64
	// message type: msgProposal msgPrevote	msgPrecommit
	step uint64
}

type timeout struct {
	timer   *time.Timer
	started bool
	sync.RWMutex
}

// runAfterTimeout() will be run in a separate go routine, so values used inside the function needs to be managed separately
func (t *timeout) scheduleTimeout(stepTimeout time.Duration, round int64, height int64, runAfterTimeout func(r int64, h int64)) *time.Timer {
	t.Lock()
	defer t.Unlock()
	t.started = true
	t.timer = time.AfterFunc(stepTimeout, func() {
		runAfterTimeout(round, height)
	})
	return t.timer
}

func (t *timeout) stopTimer() bool {
	t.RLock()
	defer t.RUnlock()
	return t.timer.Stop()
}

func (c *core) logTimeoutEvent(message string, msgType string, timeout timeoutEvent) {
	c.logger.Debug(message,
		"from", c.address.String(),
		"type", msgType,
		"currentHeight", c.currentRoundState.Height(),
		"msgHeight", timeout.heightWhenCalled,
		"currentRound", c.currentRoundState.Round(),
		"msgRound", timeout.roundWhenCalled,
		"currentStep", c.step,
		"msgStep", timeout.step,
	)
}

func (c *core) onTimeoutPropose(r int64, h int64) {
	msg := timeoutEvent{
		roundWhenCalled:  r,
		heightWhenCalled: h,
		step:             msgProposal,
	}

	c.logTimeoutEvent("TimeoutEvent(Propose): Sent", "Propose", msg)

	c.sendEvent(msg)
}

func (c *core) handleTimeoutPropose(msg timeoutEvent) {
	if msg.heightWhenCalled == c.currentRoundState.Height().Int64() && msg.roundWhenCalled == c.currentRoundState.Round().Int64() && c.step == StepAcceptProposal {
		c.logTimeoutEvent("TimeoutEvent(Propose): Received", "Propose", msg)

		c.sendPrevote(true)
		c.setStep(StepProposeDone)
	}
}

func (c *core) onTimeoutPrevote(r int64, h int64) {
	msg := timeoutEvent{
		roundWhenCalled:  r,
		heightWhenCalled: h,
		step:             msgPrevote,
	}

	c.logTimeoutEvent("TimeoutEvent(Prevote): Sent", "Prevote", msg)

	c.sendEvent(msg)

}

func (c *core) handleTimeoutPrevote(msg timeoutEvent) {
	if msg.heightWhenCalled == c.currentRoundState.Height().Int64() && msg.roundWhenCalled == c.currentRoundState.Round().Int64() && c.step == StepProposeDone {
		c.logTimeoutEvent("TimeoutEvent(Prevote): Received", "Prevote", msg)

		c.sendPrecommit(true)
		c.setStep(StepPrevoteDone)
	}
}

func (c *core) onTimeoutPrecommit(r int64, h int64) {
	msg := timeoutEvent{
		roundWhenCalled:  r,
		heightWhenCalled: h,
		step:             msgPrecommit,
	}

	c.logTimeoutEvent("TimeoutEvent(Precommit): Sent", "Precommit", msg)

	c.sendEvent(msg)
}

func (c *core) handleTimeoutPrecommit(msg timeoutEvent) {
	if msg.heightWhenCalled == c.currentRoundState.Height().Int64() && msg.roundWhenCalled == c.currentRoundState.Round().Int64() {
		c.logTimeoutEvent("TimeoutEvent(Precommit): Received", "Precommit", msg)

		c.startRound(new(big.Int).Add(c.currentRoundState.Round(), common.Big1))
	}
}

// The timeout may need to be changed depending on the Step
func timeoutPropose(round int64) time.Duration {
	return initialProposeTimeout + time.Duration(round)*time.Second
}

func timeoutPrevote(round int64) time.Duration {
	return initialPrevoteTimeout + time.Duration(round)*time.Second
}

func timeoutPrecommit(round int64) time.Duration {
	return initialPrecommitTimeout + time.Duration(round)*time.Second
}
