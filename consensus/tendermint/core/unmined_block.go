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

	c.setLatestPendingUnminedBlock(unminedBlock)

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
	logger := c.logger.New("step", c.currentRoundState.Step())

	logger.Debug("Store future unminedBlock", "number", unminedBlock.Number(), "hash", unminedBlock.Hash())

	c.pendingUnminedBlocksMu.Lock()
	defer c.pendingUnminedBlocksMu.Unlock()

	c.pendingUnminedBlocks.Push(unminedBlock, float32(-unminedBlock.Number().Int64()))
}

func (c *core) processPendingUnminedBlock() {
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
				c.logger.Debug("Stop processing request", "number", ub.Number(), "hash", ub.Hash())
				c.pendingUnminedBlocks.Push(m, prio)
				break
			}
			c.logger.Debug("Skip the pending request", "number", ub.Number(), "hash", ub.Hash(), "err", err)
			continue
		}
		c.logger.Debug("Post pending request", "number", ub.Number(), "hash", ub.Hash())

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
		"currentStep", c.currentRoundState.Step(),
		"currentProposer", c.isProposer(),
		"msgHeight", ub.Header().Number.Uint64(),
		"isNilMsg", ub.Hash() == common.Hash{},
	)
}
