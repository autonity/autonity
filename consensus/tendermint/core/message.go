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
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/validator"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/rlp"
)

const (
	msgProposal uint64 = iota
	msgPrevote
	msgPrecommit
)

type Message struct {
	Code          uint64
	Msg           []byte
	Address       common.Address
	Signature     []byte
	CommittedSeal []byte
}

// ==============================================
//
// define the functions that needs to be provided for rlp Encoder/Decoder.

// EncodeRLP serializes m into the Ethereum RLP format.
func (m *Message) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.Code, m.Msg, m.Address, m.Signature, m.CommittedSeal})
}

func (m *Message) GetCode() uint64 {
	return m.Code
}

func (m *Message) GetSignature() []byte {
	return m.Signature
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (m *Message) DecodeRLP(s *rlp.Stream) error {
	var msg struct {
		Code          uint64
		Msg           []byte
		Address       common.Address
		Signature     []byte
		CommittedSeal []byte
	}

	if err := s.Decode(&msg); err != nil {
		return err
	}
	m.Code, m.Msg, m.Address, m.Signature, m.CommittedSeal = msg.Code, msg.Msg, msg.Address, msg.Signature, msg.CommittedSeal
	return nil
}

var ErrUnauthorizedAddress = errors.New("unauthorized address")

// ==============================================
//
// define the functions that needs to be provided for core.

func (m *Message) FromPayload(b []byte, valSet validator.Set, validateFn func(validator.Set, []byte, []byte) (common.Address, error)) (*validator.Validator, error) {
	// Decode message
	err := rlp.DecodeBytes(b, m)
	if err != nil {
		return nil, err
	}

	// Validate message (on a message without Signature)
	if validateFn == nil {
		log.Error("validateFn is not set")
		return nil, nil
	}

	// Still return the message even the err is not nil
	var payload []byte
	payload, err = m.PayloadNoSig()
	if err != nil {
		return nil, err
	}

	addr, err := validateFn(valSet, payload, m.Signature)
	if err != nil {
		return nil, err
	}

	//ensure message was singed by the sender
	if !bytes.Equal(m.Address.Bytes(), addr.Bytes()) {
		return nil, ErrUnauthorizedAddress
	}

	_, v := valSet.GetByAddress(addr)
	return &v, nil
}

func (m *Message) Payload() ([]byte, error) {
	return rlp.EncodeToBytes(m)
}

func (m *Message) PayloadNoSig() ([]byte, error) {
	return rlp.EncodeToBytes(&Message{
		Code:          m.Code,
		Msg:           m.Msg,
		Address:       m.Address,
		Signature:     []byte{},
		CommittedSeal: m.CommittedSeal,
	})
}

func (m *Message) Decode(val interface{}) error {
	return rlp.DecodeBytes(m.Msg, val)
}

func (m *Message) String() string {
	return fmt.Sprintf("{Code: %v, Address: %v}", m.Code, m.Address.String())
}

// ==============================================
//
// helper functions

func Encode(val interface{}) ([]byte, error) {
	return rlp.EncodeToBytes(val)
}
