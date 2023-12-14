package core

import (
	"errors"
	"math/big"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
)

type backlogMessageEvent struct {
	msg message.Msg
}

// checkMessageStep checks the message step
func (c *Core) checkMessageStep(round int64, height uint64, step Step) error {
	h := new(big.Int).SetUint64(height)
	switch {
	// invalid messages should be catched by the peer handler, before the message gets posted to tendermint. panic.
	case round < 0 || round > constants.MaxRound:
		c.logger.Crit("Received invalid message in tendermint routine")
	// future height messages get buffered at the peer handler context, they shouldn't arrive at tendermint. panic.
	case h.Cmp(c.Height()) > 0:
		panic("Received future height message in tendermint routine")
	case h.Cmp(c.Height()) < 0:
		return constants.ErrOldHeightMessage
	case round > c.Round():
		return constants.ErrFutureRoundMessage
	case round < c.Round():
		return constants.ErrOldRoundMessage
	case c.step == Propose && step > Propose:
		return constants.ErrFutureStepMessage
	}
	return nil
}

func (c *Core) storeBacklog(msg message.Msg, src common.Address) {
	logger := c.logger.New("from", src, "step", c.step)

	if src == c.address {
		logger.Warn("Rejected backloging message, coming from local", "msg", msg)
		return
	}

	c.backlogs[src] = append(c.backlogs[src], msg)
}

func (c *Core) processBacklog() {
	var capToLenRatio = 5

	// process future round and future step msgs
	for src, backlog := range c.backlogs {
		logger := c.logger.New("from", src, "step", c.step)

		initialLen := len(backlog)
		if initialLen > 0 {
			// For loop will change the size for backlog therefore we need to keep track of the initial length and
			// adjust for index change. This is done by keeping track of how many elements have been removed and
			// subtracting it from the for-loop iterator, since each removed element will cause the index to change for
			// each element after the removed element.
			totalElemRemoved := 0
			for i := 0; i < initialLen; i++ {
				offset := i - totalElemRemoved
				curMsg := backlog[offset]

				r := curMsg.R()
				h := curMsg.H()
				err := c.checkMessageStep(r, h, Step(curMsg.Code()))
				if errors.Is(err, constants.ErrFutureRoundMessage) || errors.Is(err, constants.ErrFutureStepMessage) {
					logger.Debug("Future message in backlog", "msg", curMsg, "err", err)
					continue
				}
				logger.Debug("Post backlog event", "msg", curMsg)

				go c.SendEvent(backlogMessageEvent{
					msg: curMsg,
				})

				backlog = append(backlog[:offset], backlog[offset+1:]...)
				totalElemRemoved++
			}
			// We need to ensure that there is no memory leak by reallocating new memory if the original underlying
			// array become very large and only a small part of it is being used by the slice.
			if cap(backlog)/capToLenRatio > len(backlog) {
				tmp := make([]message.Msg, len(backlog))
				copy(tmp, backlog)
				backlog = tmp
			}
		}
		c.backlogs[src] = backlog

	}
	//TODO(lorenzo do we need parallelization)
	// process future height messages
	go c.backend.ProcessFutureMsgs(c.Height().Uint64())
}
