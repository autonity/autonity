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

// TODO: Decide whether to send an event deal with timeouts
type timeoutEvent struct{}

type timeout struct {
	timer   *time.Timer
	started bool
	sync.RWMutex
}

// runAfterTimeout() will be run in a separate go routine, so values used inside the function needs to be managed separately
func (t *timeout) scheduleTimeout(stepTimeout time.Duration, height int64, round int64, runAfterTimeout func(h int64, r int64)) *time.Timer {
	t.Lock()
	defer t.Unlock()
	t.started = true
	t.timer = time.AfterFunc(stepTimeout, func() {
		runAfterTimeout(height, round)
	})
	return t.timer
}

func (t *timeout) stopTimer() bool {
	t.RLock()
	defer t.RUnlock()
	return t.timer.Stop()
}

func (c *core) onTimeoutPropose(h int64, r int64) {
	if h == c.currentRoundState.Height().Int64() && r == c.currentRoundState.Round().Int64() && c.step == StepAcceptProposal {
		c.sendPrevote(true)
		c.setStep(StepProposeDone)
	}
}

func (c *core) onTimeoutPrevote(h int64, r int64) {
	if h == c.currentRoundState.Height().Int64() && r == c.currentRoundState.Round().Int64() && c.step == StepProposeDone {
		c.sendPrecommit(true)
		c.setStep(StepPrevoteDone)
	}
}

func (c *core) onTimeoutPrecommit(h int64, r int64) {
	if h == c.currentRoundState.Height().Int64() && r == c.currentRoundState.Round().Int64() {
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
