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
	context "context"

	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/core/types"
)

func (c *core) storeUnminedBlockMsg(ctx context.Context, unminedBlock *types.Block) {
	// c.logNewUnminedBlockEvent(unminedBlock) NOT SAFE !
	if err := c.checkUnminedBlockMsg(unminedBlock); err != nil {
		if err == errInvalidMessage {
			c.logger.Error("NewUnminedBlockEvent: invalid unminedBlock", "err", err)
			return
		}

		if err == errOldHeightMessage {
			c.logger.Error("NewUnminedBlockEvent: old height unminedBlock", "err", err)
			return
		}
	}
	c.logger.Debug("NewUnminedBlockEvent: Storing unmined block", "number", unminedBlock.NumberU64(), "hash", unminedBlock.Hash())
	c.updatePendingUnminedBlocks(ctx, unminedBlock)
}

func (c *core) updatePendingUnminedBlocks(ctx context.Context, unminedBlock *types.Block) {
	c.pendingUnminedBlocksMu.Lock()
	defer c.pendingUnminedBlocksMu.Unlock()

	// Get all heights from c.pendingUnminedBlocks and remove previous height unmined blocks
	var heights = make([]uint64, 0)
	for h := range c.pendingUnminedBlocks {
		heights = append(heights, h)
	}
	for _, ub := range heights {
		if ub < c.Height().Uint64() {
			delete(c.pendingUnminedBlocks, ub)
		}
	}

	if c.isWaitingForUnminedBlock {
		select {
		case c.pendingUnminedBlockCh <- unminedBlock:
		case <-ctx.Done():
		}
		c.isWaitingForUnminedBlock = false
	}
	c.pendingUnminedBlocks[unminedBlock.NumberU64()] = unminedBlock
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
