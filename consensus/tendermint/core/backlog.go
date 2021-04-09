// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"github.com/clearmatics/autonity/common"
	"math/big"
)

const MaxSizeBacklogUnchecked = 1000

type backlogEvent struct {
	msg *Message
}
type backlogUncheckedEvent struct {
	msg *Message
}

// checkMessage checks the message step
// return errInvalidMessage if the message is invalid
// return errFutureHeightMessage if the message view is larger than curRoundMessages view
// return errOldHeightMessage if the message view is smaller than curRoundMessages view
// return errFutureStepMessage if we are at the same view but at the propose step and it's a voting message.
func (c *core) checkMessage(round int64, height *big.Int, step Step) error {
	if height == nil || round < 0 || round > MaxRound {
		return errInvalidMessage
	}

	if height.Cmp(c.Height()) > 0 {
		return errFutureHeightMessage
	} else if height.Cmp(c.Height()) < 0 {
		return errOldHeightMessage
	} else if round > c.Round() {
		return errFutureRoundMessage
	} else if round < c.Round() {
		return errOldRoundMessage
	} else if c.step == propose && step > propose {
		return errFutureStepMessage
	}

	return nil
}

func (c *core) storeBacklog(msg *Message, src common.Address) {
	logger := c.logger.New("from", src, "step", c.step)

	if src == c.address {
		logger.Warn("Backlog from self")
		return
	}

	logger.Debug("Store future message")
	c.backlogs[src] = append(c.backlogs[src], msg)
}

// storeUncheckedBacklog push to a special backlog future height consensus messages
// this is done in a way that prevents memory exhaustion in the case of a malicious peer.
func (c *core) storeUncheckedBacklog(msg *Message) {
	// future height messages of a gap wider than one block should not occur frequently as block sync should happen
	// Todo : implement a double ended priority queue (DEPQ)

	msgHeight, errHeight := msg.Height()
	if errHeight != nil {
		panic("error parsing height")
	}

	c.backlogUnchecked[msgHeight.Uint64()] = append(c.backlogUnchecked[msgHeight.Uint64()], msg)
	c.backlogUncheckedLen++
	// We discard the furthest ahead messages in priority.
	if c.backlogUncheckedLen == MaxSizeBacklogUnchecked+1 {
		maxHeight := msgHeight.Uint64()
		for k := range c.backlogUnchecked {
			if k > maxHeight && len(c.backlogUnchecked[k]) > 0 {
				maxHeight = k
			}
		}

		// Forget in the local cache that we ever received this message.
		// It's needed for it to be able to be re-received and processed later, after a consensus sync, if needed.
		c.backend.RemoveMessageFromLocalCache(c.backlogUnchecked[maxHeight][len(c.backlogUnchecked[maxHeight])-1].Payload())

		// Remove it from the backlog buffer.
		c.backlogUnchecked[maxHeight] = c.backlogUnchecked[maxHeight][:len(c.backlogUnchecked[maxHeight])-1]
		c.backlogUncheckedLen--

		if len(c.backlogUnchecked[maxHeight]) == 0 {
			delete(c.backlogUnchecked, maxHeight)
		}
	}

}

func (c *core) processBacklog() {
	var capToLenRatio = 5

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

				r, _ := curMsg.Round()
				h, _ := curMsg.Height()
				err := c.checkMessage(r, h, Step(curMsg.Code))
				if err == errFutureHeightMessage || err == errFutureRoundMessage || err == errFutureStepMessage {
					logger.Debug("Futrue message in backlog", "msg", curMsg, "err", err)
					continue

				}
				logger.Debug("Post backlog event", "msg", curMsg)

				go c.sendEvent(backlogEvent{
					msg: curMsg,
				})

				backlog = append(backlog[:offset], backlog[offset+1:]...)
				totalElemRemoved++
			}
			// We need to ensure that there is no memory leak by reallocating new memory if the original underlying
			// array become very large and only a small part of it is being used by the slice.
			if cap(backlog)/capToLenRatio > len(backlog) {
				tmp := make([]*Message, len(backlog))
				copy(tmp, backlog)
				backlog = tmp
			}
		}
		c.backlogs[src] = backlog

	}
	for height := range c.backlogUnchecked {
		if height == c.height.Uint64() {
			for _, msg := range c.backlogUnchecked[height] {
				go c.sendEvent(backlogUncheckedEvent{
					msg: msg,
				})
				c.logger.Debug("Post unchecked backlog event", "msg", msg)
			}
		}
		if height <= c.height.Uint64() {
			delete(c.backlogUnchecked, height)
		}
	}
}
