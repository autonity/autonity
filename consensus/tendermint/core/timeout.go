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
	c.sendEvent(timeoutEvent{
		roundWhenCalled:  r,
		heightWhenCalled: h,
		step:             msgProposal,
	})
}

func (c *core) handleTimeoutPropose(timeoutE timeoutEvent) {
	if timeoutE.heightWhenCalled == c.currentRoundState.Height().Int64() && timeoutE.roundWhenCalled == c.currentRoundState.Round().Int64() && c.step == StepAcceptProposal {
		c.sendPrevote(true)
		c.setStep(StepProposeDone)
	}
}

func (c *core) onTimeoutPrevote(r int64, h int64) {
	c.sendEvent(timeoutEvent{
		roundWhenCalled:  r,
		heightWhenCalled: h,
		step:             msgPrevote,
	})

}

func (c *core) handleTimeoutPrevote(timeoutE timeoutEvent) {
	if timeoutE.heightWhenCalled == c.currentRoundState.Height().Int64() && timeoutE.roundWhenCalled == c.currentRoundState.Round().Int64() && c.step == StepProposeDone {
		c.sendPrecommit(true)
		c.setStep(StepPrevoteDone)
	}
}

func (c *core) onTimeoutPrecommit(r int64, h int64) {
	c.sendEvent(timeoutEvent{
		roundWhenCalled:  r,
		heightWhenCalled: h,
		step:             msgPrecommit,
	})
}

func (c *core) handleTimeoutPrecommit(timeoutE timeoutEvent) {
	if timeoutE.heightWhenCalled == c.currentRoundState.Height().Int64() && timeoutE.roundWhenCalled == c.currentRoundState.Round().Int64() {
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
