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
	"reflect"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint"
)

func (c *core) sendPrecommit() {
	sub := c.current.Subject()
	c.broadcastPrecommit(sub)
}

func (c *core) sendPrecommitForOldBlock(view *tendermint.View, digest common.Hash) {
	sub := &tendermint.Subject{
		View:   view,
		Digest: digest,
	}
	c.broadcastPrecommit(sub)
}

func (c *core) broadcastPrecommit(sub *tendermint.Subject) {
	logger := c.logger.New("state", c.state)

	encodedSubject, err := Encode(sub)
	if err != nil {
		logger.Error("Failed to encode", "subject", sub)
		return
	}
	c.broadcast(&message{
		Code: msgPrecommit,
		Msg:  encodedSubject,
	})
}

func (c *core) handlePrecommit(msg *message, src tendermint.Validator) error {
	// Decode COMMIT message
	var commit *tendermint.Subject
	err := msg.Decode(&commit)
	if err != nil {
		return errFailedDecodePrecommit
	}

	if err := c.checkMessage(msgPrecommit, commit.View); err != nil {
		return err
	}

	if err := c.verifyPrecommit(commit, src); err != nil {
		return err
	}

	if err := c.acceptPrecommit(msg); err != nil {
		c.logger.Error("Failed to record commit message", "from", src, "state", c.state, "msg", msg, "err", err)
	}

	// Precommit the proposal once we have enough COMMIT messages and we are not in the Committed state.
	//
	// If we already have a proposal, we may have chance to speed up the consensus process
	// by committing the proposal without PREPARE messages.
	if c.current.Precommits.Size() > 2*c.valSet.F() && c.state.Cmp(StatePrecommitDone) < 0 {
		// Still need to call LockHash here since state can skip Prevoted state and jump directly to the Committed state.
		c.current.LockHash()
		c.commit()
	}

	return nil
}

// verifyPrecommit verifies if the received COMMIT message is equivalent to our subject
func (c *core) verifyPrecommit(commit *tendermint.Subject, src tendermint.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	sub := c.current.Subject()
	if !reflect.DeepEqual(commit, sub) {
		logger.Warn("Inconsistent subjects between commit and proposal", "expected", sub, "got", commit)
		return errInconsistentSubject
	}

	return nil
}

// acceptPrecommit adds the COMMIT message to current round state
func (c *core) acceptPrecommit(msg *message) error {
	return c.current.Precommits.Add(msg)
}
