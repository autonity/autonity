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

package types

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/rlp"
	"io"
	"math/big"
	"reflect"
)

const (
	MsgProposal uint64 = iota
	MsgPrevote
	MsgPrecommit
)

var (
	errMsgPayloadNotDecoded = errors.New("message not decoded")
	ErrUnauthorizedAddress  = errors.New("unauthorized address")
)

type ConsensusMessage struct {
	Code          uint64
	Msg           []byte
	Address       common.Address
	Signature     []byte
	CommittedSeal []byte

	Power      uint64
	decodedMsg ConsensusMsg // cached decoded Msg
	payload    []byte       // rlp encoded ConsensusMessage
}

// ==============================================
//
// define the functions that needs to be provided for rlp Encoder/Decoder.

// EncodeRLP serializes m into the Ethereum RLP format.
func (m *ConsensusMessage) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.Code, m.Msg, m.Address, m.Signature, m.CommittedSeal})
}

func (m *ConsensusMessage) GetCode() uint64 {
	return m.Code
}

func (m *ConsensusMessage) GetSignature() []byte {
	return m.Signature
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (m *ConsensusMessage) DecodeRLP(s *rlp.Stream) error {
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

// ==============================================
//
// define the functions that needs to be provided for core.

func (m *ConsensusMessage) FromPayload(b []byte) error {
	m.payload = b
	// Decode message
	err := rlp.DecodeBytes(b, m)
	if err != nil {
		return err
	}
	// Decode the payload, this will cache the decoded msg payload.
	switch m.Code {
	case MsgProposal:
		var proposal Proposal
		return m.Decode(&proposal)
	case MsgPrevote, MsgPrecommit:
		var vote Vote
		return m.Decode(&vote)
	default:
		return errMsgPayloadNotDecoded
	}
}

func (m *ConsensusMessage) Validate(validateFn func(*Header, []byte, []byte) (common.Address, error), previousHeader *Header) (*CommitteeMember, error) {
	// Validate message (on a message without Signature)
	msgHeight, err := m.Height()
	if err != nil {
		return nil, err
	}
	if previousHeader.Number.Uint64()+1 != msgHeight.Uint64() {
		panic("inconsistent message verification")
	}

	// Still return the message even the err is not nil
	var payload []byte
	payload, err = m.PayloadNoSig()
	if err != nil {
		return nil, err
	}

	addr, err := validateFn(previousHeader, payload, m.Signature)
	if err != nil {
		return nil, err
	}

	//ensure message was singed by the sender
	if !bytes.Equal(m.Address.Bytes(), addr.Bytes()) {
		return nil, ErrUnauthorizedAddress
	}

	v := previousHeader.CommitteeMember(addr)
	if v == nil {
		return nil, fmt.Errorf("message received is not from a committee member: %x", addr)
	}

	m.Power = v.VotingPower.Uint64()
	return v, nil
}

func (m *ConsensusMessage) Payload() []byte {
	if m.payload == nil {
		payload, err := rlp.EncodeToBytes(m)
		if err != nil {
			// We panic if there is an error, reasons:
			// Either we received the message and we managed to decode it, hence it must be possible to encode it.
			// If we can't encode the payload for our own generated messages, that's a programming error.
			panic("could not decode message payload")
		}
		m.payload = payload
	}
	return m.payload
}

func (m *ConsensusMessage) GetPower() uint64 {
	return m.Power
}

func (m *ConsensusMessage) PayloadNoSig() ([]byte, error) {
	return rlp.EncodeToBytes(&ConsensusMessage{
		Code:          m.Code,
		Msg:           m.Msg,
		Address:       m.Address,
		Signature:     []byte{},
		CommittedSeal: m.CommittedSeal,
	})
}

func (m *ConsensusMessage) Decode(val interface{}) error {
	//Decode is responsible to rlp-decode m.Msg. It is meant to only perform the actual decoding once,
	//saving a cached value in m.decodedMsg.

	rval := reflect.ValueOf(val)
	if rval.Kind() != reflect.Ptr {
		return errors.New("decode arg must be a pointer")
	}

	// check if we already have a cached value decoded
	if m.decodedMsg != nil {
		if !rval.Type().AssignableTo(reflect.TypeOf(m.decodedMsg)) {
			return errors.New("type mismatch with decoded value")
		}
		rval.Elem().Set(reflect.ValueOf(m.decodedMsg).Elem())
		return nil
	}

	err := rlp.DecodeBytes(m.Msg, val)
	if err != nil {
		return err
	}

	// copy the result via Set (shallow)
	nval := reflect.New(rval.Elem().Type()) // we need first to allocate memory
	nval.Elem().Set(rval.Elem())
	m.decodedMsg = nval.Interface().(ConsensusMsg)
	return nil
}

func (m *ConsensusMessage) String() string {
	var msg string
	if m.Code == MsgProposal {
		var proposal Proposal
		err := m.Decode(&proposal)
		if err != nil {
			return ""
		}
		msg = proposal.String()
	}

	if m.Code == MsgPrevote || m.Code == MsgPrecommit {
		var vote Vote
		err := m.Decode(&vote)
		if err != nil {
			return ""
		}
		msg = vote.String()
	}
	return fmt.Sprintf("{sender: %v, Power: %v, msgCode: %v, msg: %v}", m.Address.String(), m.Power, m.Code, msg)
}

func (m *ConsensusMessage) Round() (int64, error) {
	if m.decodedMsg == nil {
		return 0, errMsgPayloadNotDecoded
	}
	return m.decodedMsg.GetRound(), nil
}

func (m *ConsensusMessage) Height() (*big.Int, error) {
	if m.decodedMsg == nil {
		return nil, errMsgPayloadNotDecoded
	}
	return m.decodedMsg.GetHeight(), nil
}

// used by afd for decoded msgs
func (m *ConsensusMessage) R() uint {
	r, err := m.Round()
	// msg should be decoded, it shouldn't be an error.
	if err != nil {
		panic(err)
	}
	return uint(r)
}

func (m *ConsensusMessage) H() uint64 {
	h, err := m.Height()
	if err != nil {
		panic(err)
	}
	return h.Uint64()
}

func (m *ConsensusMessage) Sender() common.Address {
	return m.Address
}

func (m *ConsensusMessage) Type() uint64 {
	return m.Code
}

func (m *ConsensusMessage) Value() common.Hash {
	return m.decodedMsg.GetValue()
}

func (m *ConsensusMessage) ValidRound() int64 {
	return m.decodedMsg.GetValidRound()
}

// ==============================================
//
// helper functions

func Encode(val interface{}) ([]byte, error) {
	return rlp.EncodeToBytes(val)
}
