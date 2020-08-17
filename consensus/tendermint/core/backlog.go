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
	"github.com/clearmatics/autonity/core/types"
	"math/big"
)

var (
	// msgPriority is defined for calculating processing priority to speedup consensus
	// msgProposal > msgPrecommit > msgPrevote
	msgPriority = map[uint64]int{
		msgProposal:  1,
		msgPrecommit: 2,
		msgPrevote:   3,
	}
)

type backlogEvent struct {
	src types.CommitteeMember
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

func (c *core) storeBacklog(msg *Message, src types.CommitteeMember) {
	logger := c.logger.New("from", src, "step", c.step)

	if src.Address == c.address {
		logger.Warn("Backlog from self")
		return
	}

	logger.Debug("Store future message")

	c.backlogsMu.Lock()
	defer c.backlogsMu.Unlock()

	c.backlogs[src] = append(c.backlogs[src], msg)
}

func (c *core) processBacklog() {
	var capToLenRatio = 5
	c.backlogsMu.Lock()
	defer c.backlogsMu.Unlock()

	for src, backlog := range c.backlogs {
		logger := c.logger.New("from", src, "step", c.step)

		initialLen := len(backlog)
		if initialLen > 0 {
			// For loop will change the size for backlog therefore we need to keep track of the initial length and
			// adjust for index change. This is done by keeping track of how many elements have been removed and
			// subtracting it from the for-loop iterator, since each removed element will cause the index to change for
			// each element after the removed element.
			//
			// If the message is a future height, round or step message then the for-loop will move on to the next
			// iteration, however, if any other error occurs then the message is removed from the backlog and for-loop
			// moves to the next iteration.
			//
			// If there are no errors when checking the message then the message is sent to the handler via a goroutine
			// and this message is removed from the backlog.
			totalElemRemoved := 0
			for i := 0; i < initialLen; i++ {
				offset := i - totalElemRemoved
				curMsg := backlog[offset]

				r, _ := curMsg.Round()
				h, _ := curMsg.Height()
				if err := c.checkMessage(r, h, Step(curMsg.Code)); err != nil {
					if err == errFutureHeightMessage || err == errFutureRoundMessage || err == errFutureStepMessage {
						logger.Debug("Futrue message in backlog", "msg", curMsg, "err", err)
						continue
					}
					logger.Debug("Skipping the backlog message", "msg", curMsg, "err", err)
				} else {
					logger.Debug("Post backlog event", "msg", curMsg)

					go c.sendEvent(backlogEvent{
						src: src,
						msg: curMsg,
					})
				}

				backlog = append(backlog[:offset], backlog[offset+1:]...)
				totalElemRemoved++
			}
			// We need to ensure that there is no memory leak by reallocating new memory if the original underlying
			// array become very large and only a small part of it is being used by the slice. We need len(backlog) > 0
			// check again since the backlog size can change and cause division by zero errors.
			if len(backlog) > 0 && cap(backlog)/len(backlog) > capToLenRatio {
				tmp := make([]*Message, len(backlog))
				copy(tmp, backlog)
				backlog = tmp
			}
		}
		c.backlogs[src] = backlog
	}
}

func toPriority(msgCode uint64, r int64, h *big.Int) float32 {
	// TODO check for overflows!!
	// 10 * Round limits the range of message code is from 0 to 9
	// 1000 * Height limits the range of round is from 0 to 99
	return -float32(h.Uint64()*10*(MaxRound+1) + uint64(r)*10 + uint64(msgPriority[msgCode]))
}
