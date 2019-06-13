package core

import (
	"math/big"
	"sync"
	"time"

	"github.com/clearmatics/autonity/common"
)

const (
	initialProposeTimeout   = 5 * time.Second
	initialPrevoteTimeout   = 5 * time.Second
	initialPrecommitTimeout = 5 * time.Second
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

func (c *core) onTimeoutPropose(r int64, h int64) {
	msg := timeoutEvent{
		roundWhenCalled:  r,
		heightWhenCalled: h,
		step:             msgProposal,
	}

	c.logger.Info("MESSAGE: internal message",
		"type", "timeoutPropose",
		"currentHeight", c.currentRoundState.height,
		"currentRound", c.currentRoundState.round,
		"currentStep", c.step,
		"from", c.address.String(),
		"currentProposer", c.isProposer(),
		"msgHeight", msg.heightWhenCalled,
		"msgRound", msg.roundWhenCalled,
		"msgStep", msg.step,
		"message", msg,
	)

	c.sendEvent(msg)
}

func (c *core) handleTimeoutPropose(msg timeoutEvent) {
	if msg.heightWhenCalled == c.currentRoundState.Height().Int64() && msg.roundWhenCalled == c.currentRoundState.Round().Int64() && c.step == StepAcceptProposal {
		c.logger.Info("MESSAGE: handle internal message",
			"type", "timeoutPropose",
			"currentHeight", c.currentRoundState.height,
			"currentRound", c.currentRoundState.round,
			"currentStep", c.step,
			"from", c.address.String(),
			"currentProposer", c.isProposer(),
			"msgHeight", msg.heightWhenCalled,
			"msgRound", msg.roundWhenCalled,
			"msgStep", msg.step,
			"message", msg,
		)

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

	c.logger.Info("MESSAGE: internal message",
		"type", "timeoutPrevote",
		"currentHeight", c.currentRoundState.height,
		"currentRound", c.currentRoundState.round,
		"currentStep", c.step,
		"from", c.address.String(),
		"currentProposer", c.isProposer(),
		"msgHeight", msg.heightWhenCalled,
		"msgRound", msg.roundWhenCalled,
		"msgStep", msg.step,
		"message", msg,
	)

	c.sendEvent(msg)

}

func (c *core) handleTimeoutPrevote(msg timeoutEvent) {
	if msg.heightWhenCalled == c.currentRoundState.Height().Int64() && msg.roundWhenCalled == c.currentRoundState.Round().Int64() && c.step == StepProposeDone {
		c.logger.Info("MESSAGE: handle internal message",
			"type", "timeoutPrevote",
			"currentHeight", c.currentRoundState.height,
			"currentRound", c.currentRoundState.round,
			"currentStep", c.step,
			"from", c.address.String(),
			"currentProposer", c.isProposer(),
			"msgHeight", msg.heightWhenCalled,
			"msgRound", msg.roundWhenCalled,
			"msgStep", msg.step,
			"message", msg,
		)

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

	c.logger.Info("MESSAGE: internal message",
		"type", "timeoutPrecommit",
		"currentHeight", c.currentRoundState.height,
		"currentRound", c.currentRoundState.round,
		"currentStep", c.step,
		"from", c.address.String(),
		"currentProposer", c.isProposer(),
		"msgHeight", msg.heightWhenCalled,
		"msgRound", msg.roundWhenCalled,
		"msgStep", msg.step,
		"message", msg,
	)

	c.sendEvent(msg)
}

func (c *core) handleTimeoutPrecommit(msg timeoutEvent) {
	if msg.heightWhenCalled == c.currentRoundState.Height().Int64() && msg.roundWhenCalled == c.currentRoundState.Round().Int64() {
		c.logger.Info("MESSAGE: handle internal message",
			"type", "timeoutPrecommit",
			"currentHeight", c.currentRoundState.height,
			"currentRound", c.currentRoundState.round,
			"currentStep", c.step,
			"from", c.address.String(),
			"currentProposer", c.isProposer(),
			"msgHeight", msg.heightWhenCalled,
			"msgRound", msg.roundWhenCalled,
			"msgStep", msg.step,
			"message", msg,
		)

		c.startRound(new(big.Int).Add(c.currentRoundState.Height(), common.Big1))
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
