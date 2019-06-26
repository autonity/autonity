package core

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/tendermint"
	"github.com/clearmatics/autonity/core/types"
)

func (c *core) handleUnminedBlock(unminedBlock *types.Block) error {
	if err := c.checkUnminedBlockMsg(unminedBlock); err != nil {
		if err == errInvalidMessage {
			c.logger.Warn("invalid unminedBlock")
			return err
		}
		c.logger.Warn("unexpected unminedBlock", "err", err, "number", unminedBlock.Number(), "hash", unminedBlock.Hash())
		return err
	}

	c.logNewUnminedBlockEvent(unminedBlock)

	c.latestPendingUnminedBlockMu.Lock()
	wasNilOrDiffHeight := c.latestPendingUnminedBlock == nil || c.latestPendingUnminedBlock.Number() != c.currentRoundState.Height()
	c.latestPendingUnminedBlock = unminedBlock
	c.latestPendingUnminedBlockMu.Unlock()

	if wasNilOrDiffHeight {
		c.unminedBlockCh <- struct{}{}
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

	number := unminedBlock.Number()
	if c := c.currentRoundState.Height().Cmp(number); c > 0 {
		// TODO: probably delete this case?
		return errOldHeightMessage
	} else if c < 0 {
		return consensus.ErrFutureBlock
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

	for !c.pendingUnminedBlocks.Empty() {
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

		c.sendEvent(tendermint.NewUnminedBlockEvent{
			NewUnminedBlock: *ub,
		})
	}
}

func (c *core) logNewUnminedBlockEvent(ub *types.Block) {
	c.logger.Debug("NewUnminedBlockEvent: Received",
		"from", c.address.String(),
		"type", "New Unmined Block",
		"hash", ub.Hash(),
		"currentHeight", c.currentRoundState.Height(),
		"currentRound", c.currentRoundState.Round(),
		"currentStep", c.step,
		"currentProposer", c.isProposer(),
		"msgHeight", ub.Header().Number.Uint64(),
		"isNilMsg", ub.Hash() == common.Hash{},
	)
}
