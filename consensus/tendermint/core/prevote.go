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

	"github.com/clearmatics/autonity/consensus/tendermint"
)

func (c *core) sendPrevote() {
	logger := c.logger.New("state", c.state)

	sub := c.current.Subject()
	encodedSubject, err := Encode(sub)
	if err != nil {
		logger.Error("Failed to encode", "subject", sub)
		return
	}
	c.broadcast(&message{
		Code: msgPrevote,
		Msg:  encodedSubject,
	})
}

func (c *core) handlePrevote(msg *message, src tendermint.Validator) error {
	// Decode PREPARE message
	var prepare *tendermint.Subject
	err := msg.Decode(&prepare)
	if err != nil {
		return errFailedDecodePrevote
	}

	if err := c.checkMessage(msgPrevote, prepare.View); err != nil {
		return err
	}

	// If it is locked, it can only process on the locked block.
	// Passing verifyPrevote and checkMessage implies it is processing on the locked block since it was verified in the Proposald state.
	if err := c.verifyPrevote(prepare, src); err != nil {
		return err
	}

	c.acceptPrevote(msg, src)

	// Change to Prevoted state if we've received enough PREPARE messages or it is locked
	// and we are in earlier state before Prevoted state.
	if ((c.current.IsHashLocked() && prepare.Digest == c.current.GetLockedHash()) || c.current.GetPrevoteOrPrecommitSize() > 2*c.valSet.F()) &&
		c.state.Cmp(StatePrevoteDone) < 0 {
		c.current.LockHash()
		c.setState(StatePrevoteDone)
		c.sendPrecommit()
	}

	return nil
}

// verifyPrevote verifies if the received PREPARE message is equivalent to our subject
func (c *core) verifyPrevote(prepare *tendermint.Subject, src tendermint.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	sub := c.current.Subject()
	if !reflect.DeepEqual(prepare, sub) {
		logger.Warn("Inconsistent subjects between PREPARE and proposal", "expected", sub, "got", prepare)
		return errInconsistentSubject
	}

	return nil
}

func (c *core) acceptPrevote(msg *message, src tendermint.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	// Add the PREPARE message to current round state
	if err := c.current.Prevotes.Add(msg); err != nil {
		logger.Error("Failed to add PREPARE message to round state", "msg", msg, "err", err)
		return err
	}

	return nil
}
