package core

import (
	"errors"
	"math/big"

	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
)

type backlogMessageEvent struct {
	msg message.Msg
}

func (c *Core) checkMessage(round int64, height uint64) error {
	h := new(big.Int).SetUint64(height)
	switch {
	// future height messages get buffered at the peer handler context, they shouldn't arrive at tendermint. panic.
	case h.Cmp(c.Height()) > 0:
		panic("Received future height message in tendermint routine")
	case h.Cmp(c.Height()) < 0:
		return constants.ErrOldHeightMessage
	case round > c.Round():
		return constants.ErrFutureRoundMessage
	case round < c.Round():
		return constants.ErrOldRoundMessage
	}
	return nil
}

// TODO(lorenzo) do something more smart
// func (c *Core) storeBacklog(msg message.Msg, src common.Address) {
func (c *Core) storeBacklog(msg message.Msg) {
	//logger := c.logger.New("from", src, "step", c.step)
	//logger := c.logger.New("step", c.step)

	/*
		if src == c.address {
			logger.Warn("Rejected backloging message, coming from local", "msg", msg)
			return
		}*/

	//c.backlogs[src] = append(c.backlogs[src], msg)
	c.backlogs = append(c.backlogs, msg)
}

func (c *Core) processBacklog() {
	initialLen := len(c.backlogs)
	if initialLen > 0 {
		// For loop will change the size for backlog therefore we need to keep track of the initial length and
		// adjust for index change. This is done by keeping track of how many elements have been removed and
		// subtracting it from the for-loop iterator, since each removed element will cause the index to change for
		// each element after the removed element.
		totalElemRemoved := 0
		for i := 0; i < initialLen; i++ {
			offset := i - totalElemRemoved
			curMsg := c.backlogs[offset]

			r := curMsg.R()
			h := curMsg.H()
			err := c.checkMessage(r, h)
			if errors.Is(err, constants.ErrFutureRoundMessage) {
				//logger.Debug("Future message in backlog", "msg", curMsg, "err", err)
				continue

			}
			//logger.Debug("Post backlog event", "msg", curMsg)

			go c.SendEvent(backlogMessageEvent{
				msg: curMsg,
			})

			c.backlogs = append(c.backlogs[:offset], c.backlogs[offset+1:]...)
			totalElemRemoved++
		}
	}
	//TODO(lorenzo do we need parallelization)
	// process future height messages
	go c.backend.ProcessFutureMsgs(c.Height().Uint64())
	/*
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
					err := c.checkMessage(r, h)
					if errors.Is(err, constants.ErrFutureRoundMessage) {
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
	*/
}
