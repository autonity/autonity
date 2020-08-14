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
	"github.com/clearmatics/autonity_cookiejar/collections/prque"
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

	backlogPrque := c.backlogs[src]
	if backlogPrque == nil {
		backlogPrque = prque.New()
	}
	msgRound, errRound := msg.Round()
	msgHeight, errHeight := msg.Height()
	if errRound == nil && errHeight == nil {
		backlogPrque.Push(msg, toPriority(msg.Code, msgRound, msgHeight))
	}

	c.backlogs[src] = backlogPrque
}

func (c *core) processBacklog() {
	c.backlogsMu.Lock()
	defer c.backlogsMu.Unlock()

	for src, backlog := range c.backlogs {
		if backlog == nil {
			continue
		}

		logger := c.logger.New("from", src, "step", c.step)
		var isFuture bool

		// We stop processing if
		//   1. backlog is empty
		//   2. The first message in queue is a future message
		for !(backlog.Empty() || isFuture) {
			m, prio := backlog.Pop()
			msg := m.(*Message)
			msgRound, _ := msg.Round() // error checking done before push
			msgHeight, _ := msg.Height()

			// Push back if it's a future message
			err := c.checkMessage(msgRound, msgHeight, Step(msg.Code))
			if err != nil {
				if err == errFutureHeightMessage || err == errFutureRoundMessage || err == errFutureStepMessage {
					logger.Debug("Stop processing backlog", "msg", msg, "err", err)
					backlog.Push(msg, prio)
					isFuture = true
					break
				}
				logger.Debug("Skip the backlog event", "msg", msg, "err", err)
				continue
			}
			logger.Debug("Post backlog event", "msg", msg)

			go c.sendEvent(backlogEvent{
				src: src,
				msg: msg,
			})
		}
	}
}

func toPriority(msgCode uint64, r int64, h *big.Int) int64 {
	// 10 * Round limits the range of message code is from 0 to 9
	// 1000 * Height limits the range of round is from 0 to 99
	return -(h.Int64()*10*(MaxRound+1) + r*10 + int64(msgPriority[msgCode]))
}
