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
	"math/big"
	"reflect"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/rlp"
)

const (
	msgProposal uint64 = iota
	msgPrevote
	msgPrecommit
)

var (
	errMsgPayloadNotDecoded = errors.New("msg not decoded")
)

type Message struct {
	Code          uint64
	Msg           []byte
	Address       common.Address
	Signature     []byte
	CommittedSeal []byte

	power      uint64
	decodedMsg ConsensusMsg // cached decoded Msg
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

func (m *Message) FromPayload(b []byte, previousHeader *types.Header, validateFn func(*types.Header, []byte, []byte) (common.Address, error)) (*types.CommitteeMember, error) {
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
		return nil, fmt.Errorf("validator was not a committee member %q", v)
	}

	m.power = v.VotingPower.Uint64()
	return v, nil
}

func (m *Message) Payload() ([]byte, error) {
	return rlp.EncodeToBytes(m)
}

func (m *Message) GetPower() uint64 {
	return m.power
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

func (m *Message) String() string {
	return fmt.Sprintf("{Code: %v, Address: %v, Msg: %v, Signature: %v, CommittedSeal: %v, Power: %v}",
		m.Code, m.Address.String(), m.Msg, m.Signature, m.CommittedSeal, m.power)
}

func (m *Message) Round() (int64, error) {
	if m.decodedMsg == nil {
		return 0, errMsgPayloadNotDecoded
	}
	return m.decodedMsg.GetRound(), nil
}

func (m *Message) Height() (*big.Int, error) {
	if m.decodedMsg == nil {
		return nil, errMsgPayloadNotDecoded
	}
	return m.decodedMsg.GetHeight(), nil
}

// ==============================================
//
// helper functions

func Encode(val interface{}) ([]byte, error) {
	return rlp.EncodeToBytes(val)
}
