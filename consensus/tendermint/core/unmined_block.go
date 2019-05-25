package core

import (
	"github.com/clearmatics/autonity/consensus/tendermint"
	"github.com/clearmatics/autonity/core/types"
)

func (c *core) handleUnminedBlock(unminedBlock *types.Block) error {
	logger := c.logger.New("step", c.step, "height", c.currentRoundState.height)

	if err := c.checkUnminedBlockMsg(unminedBlock); err != nil {
		if err == errInvalidMessage {
			logger.Warn("invalid unminedBlock")
			return err
		}
		logger.Warn("unexpected unminedBlock", "err", err, "number", unminedBlock.Number(), "hash", unminedBlock.Hash())
		return err
	}

	logger.Trace("handleUnminedBlock", "number", unminedBlock.Number(), "hash", unminedBlock.Hash())

	c.latestPendingUnminedBlock = unminedBlock
	// TODO: remove, we should not be sending a proposal from handleUnminedBlock
	if c.step == StepAcceptProposal {
		c.sendProposal(unminedBlock)
	}
	return nil
}

// check request step
// return errInvalidMessage if the message is invalid
// return errFutureHeightMessage if the height of proposal is larger than currentRoundState height
// return errOldHeightMessage if the height of proposal is smaller than currentRoundState height
func (c *core) checkUnminedBlockMsg(unminedBlock *types.Block) error {
	if unminedBlock == nil {
		return errInvalidMessage
	}

	//TODO: make the err message more specific to block maybe use consensus.ErrFutureBlock?
	if c := c.currentRoundState.height.Cmp(unminedBlock.Number()); c > 0 {
		return errOldHeightMessage
	} else if c < 0 {
		return errFutureHeightMessage
	} else {
		return nil
	}
}

func (c *core) storeUnminedBlockMsg(unminedBlock *types.Block) {
	logger := c.logger.New("step", c.step)

	logger.Trace("Store future unminedBlock", "number", unminedBlock.Number(), "hash", unminedBlock.Hash())

	c.pendingUnminedBlocksMu.Lock()
	defer c.pendingUnminedBlocksMu.Unlock()

	c.pendingUnminedBlocks.Push(unminedBlock, float32(-unminedBlock.Number().Int64()))
}

func (c *core) processPendingRequests() {
	c.pendingUnminedBlocksMu.Lock()
	defer c.pendingUnminedBlocksMu.Unlock()

	for !(c.pendingUnminedBlocks.Empty()) {
		m, prio := c.pendingUnminedBlocks.Pop()
		ub, ok := m.(*types.Block)
		if !ok {
			c.logger.Warn("Malformed request, skip", "msg", m)
			continue
		}
		// Push back if it's a future message
		err := c.checkUnminedBlockMsg(ub)
		if err != nil {
			if err == errFutureHeightMessage {
				c.logger.Trace("Stop processing request", "number", ub.Number(), "hash", ub.Hash())
				c.pendingUnminedBlocks.Push(m, prio)
				break
			}
			c.logger.Trace("Skip the pending request", "number", ub.Number(), "hash", ub.Hash(), "err", err)
			continue
		}
		c.logger.Trace("Post pending request", "number", ub.Number(), "hash", ub.Hash())

		go c.sendEvent(tendermint.NewUnminedBlockEvent{
			NewUnminedBlock: *ub,
		})
	}
}
