package core

import (
	context "context"

	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/core/types"
)

func (c *core) storeUnminedBlockMsg(ctx context.Context, unminedBlock *types.Block) {
	// c.logNewUnminedBlockEvent(unminedBlock) NOT SAFE !
	if err := c.checkUnminedBlockMsg(unminedBlock); err != nil {
		if err == errInvalidMessage {
			c.logger.Error("NewUnminedBlockEvent: invalid unminedBlock", "err", err)
			return
		}

		if err == errOldHeightMessage {
			c.logger.Debug("NewUnminedBlockEvent: old height unminedBlock", "err", err)
			return
		}
	}
	c.logger.Debug("NewUnminedBlockEvent: Storing unmined block", "number", unminedBlock.NumberU64(), "hash", unminedBlock.Hash())
	c.updatePendingUnminedBlocks(ctx, unminedBlock)
}

func (c *core) updatePendingUnminedBlocks(ctx context.Context, unminedBlock *types.Block) {
	c.pendingUnminedBlocksMu.Lock()
	defer c.pendingUnminedBlocksMu.Unlock()
	c.pendingUnminedBlocks[unminedBlock.NumberU64()] = unminedBlock
	if c.isWaitingForUnminedBlock {
		// if mainEventLoop is blocking at startRound() as a proposer waiting for the unMinedBlock of a height,
		//get the corresponding block from buffer, and forward the one that startRound() is waiting for.
		waitForBlock := c.pendingUnminedBlocks[c.Height().Uint64()]
		if waitForBlock != nil {
			select {
			case c.pendingUnminedBlockCh <- unminedBlock:
			case <-ctx.Done():
			}
			c.isWaitingForUnminedBlock = false
		}
	}

	// release buffered un-mined blocks before the height of waitForBlock
	for height := range c.pendingUnminedBlocks {
		if height < c.Height().Uint64() {
			delete(c.pendingUnminedBlocks, height)
		}
	}
}

func (c *core) getUnminedBlock() *types.Block {
	c.pendingUnminedBlocksMu.Lock()
	defer c.pendingUnminedBlocksMu.Unlock()

	ub, ok := c.pendingUnminedBlocks[c.Height().Uint64()]

	if ok {
		return ub
	}

	c.isWaitingForUnminedBlock = true
	return nil
}

// check request step
// return errInvalidMessage if the message is invalid
// return errFutureHeightMessage if the height of proposal is larger than curRoundMessages height
// return errOldHeightMessage if the height of proposal is smaller than curRoundMessages height
func (c *core) checkUnminedBlockMsg(unminedBlock *types.Block) error {
	if unminedBlock == nil {
		return errInvalidMessage
	}

	number := unminedBlock.Number()
	if currentIsHigher := c.Height().Cmp(number); currentIsHigher > 0 {
		return errOldHeightMessage
	} else if currentIsHigher < 0 {
		return consensus.ErrFutureBlock
	} else {
		return nil
	}
}
