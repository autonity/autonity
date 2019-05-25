package core

import (
	"math/big"

	"github.com/clearmatics/autonity/consensus/tendermint"
	"gopkg.in/karalabe/cookiejar.v2/collections/prque"
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
	src tendermint.Validator
	msg *message
}

// checkMessage checks the message step
// return errInvalidMessage if the message is invalid
// return errFutureHeightMessage if the message view is larger than currentRoundState view
// return errOldHeightMessage if the message view is smaller than currentRoundState view
func (c *core) checkMessage(msgCode uint64, round *big.Int, height *big.Int) error {
	if height == nil || round == nil {
		return errInvalidMessage
	}

	if height.Cmp(c.currentRoundState.Height()) > 0 {
		return errFutureHeightMessage
	} else if height.Cmp(c.currentRoundState.Height()) < 0 {
		return errOldHeightMessage
	} else if round.Cmp(c.currentRoundState.Round()) > 0 {
		return errFutureRoundMessage
	} else if round.Cmp(c.currentRoundState.Round()) < 0 {
		return errOldRoundMessage
	}

	// StepAcceptProposal only accepts msgProposal
	// other messages are future messages
	if c.step == StepAcceptProposal {
		if msgCode > msgProposal {
			return errFutureHeightMessage
		}
		return nil
	}

	// For steps(StepProposeDone, StepPrevoteDone, StepPrecommitDone),
	// can accept all message types if processing with same view
	return nil
}

func (c *core) storeBacklog(msg *message, src tendermint.Validator) {
	logger := c.logger.New("from", src, "step", c.step)

	if src.Address() == c.Address() {
		logger.Warn("Backlog from self")
		return
	}

	logger.Trace("Store future message")

	c.backlogsMu.Lock()
	defer c.backlogsMu.Unlock()

	backlog := c.backlogs[src]
	if backlog == nil {
		backlog = prque.New()
	}
	switch msg.Code {
	case msgProposal:
		var p *tendermint.Proposal
		err := msg.Decode(&p)
		if err == nil {
			backlog.Push(msg, toPriority(msg.Code, p.Round, p.Height))
		}
		// for msgRoundChange, msgPrevote and msgPrecommit cases
	default:
		var p *tendermint.Vote
		err := msg.Decode(&p)
		if err == nil {
			backlog.Push(msg, toPriority(msg.Code, p.Round, p.Height))
		}
	}
	c.backlogs[src] = backlog
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
			msg := m.(*message)
			var round, height *big.Int
			switch msg.Code {
			case msgProposal:
				var m *tendermint.Proposal
				err := msg.Decode(&m)
				if err == nil {
					round, height = m.Round, m.Height
				}
				// for msgRoundChange, msgPrevote and msgPrecommit cases
			default:
				var sub *tendermint.Vote
				err := msg.Decode(&sub)
				if err == nil {
					round, height = sub.Round, sub.Height
				}
			}
			if round == nil || height == nil {
				logger.Debug("Nil view", "msg", msg)
				continue
			}
			// Push back if it's a future message
			err := c.checkMessage(msg.Code, round, height)
			if err != nil {
				if err == errFutureHeightMessage {
					logger.Trace("Stop processing backlog", "msg", msg)
					backlog.Push(msg, prio)
					isFuture = true
					break
				}
				logger.Trace("Skip the backlog event", "msg", msg, "err", err)
				continue
			}
			logger.Trace("Post backlog event", "msg", msg)

			go c.sendEvent(backlogEvent{
				src: src,
				msg: msg,
			})
		}
	}
}

func toPriority(msgCode uint64, r *big.Int, h *big.Int) float32 {
	// FIXME: round will be reset as 0 while new height
	// 10 * Round limits the range of message code is from 0 to 9
	// 1000 * Height limits the range of round is from 0 to 99
	return -float32(h.Uint64()*1000 + r.Uint64()*10 + uint64(msgPriority[msgCode]))
}
