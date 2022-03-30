package core

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/rlp"
	"io"
	"math/big"
	"reflect"
)

const (
	msgProposal uint64 = iota
	msgPrevote
	msgPrecommit
)

var (
	errMsgPayloadNotDecoded = errors.New("message not decoded")
	ErrUnauthorizedAddress  = errors.New("unauthorized address")
)

type Message struct {
	Code          uint64
	Msg           []byte
	Address       common.Address
	Signature     []byte
	CommittedSeal []byte

	power      uint64
	decodedMsg ConsensusMsg // cached decoded Msg
	payload    []byte       // rlp encoded Message
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

// ==============================================
//
// define the functions that needs to be provided for core.

func (m *Message) FromPayload(b []byte) error {
	m.payload = b
	// Decode message
	err := rlp.DecodeBytes(b, m)
	if err != nil {
		return err
	}
	// Decode the payload, this will cache the decoded msg payload.
	switch m.Code {
	case msgProposal:
		var proposal Proposal
		return m.Decode(&proposal)
	case msgPrevote, msgPrecommit:
		var vote Vote
		return m.Decode(&vote)
	default:
		return errMsgPayloadNotDecoded
	}
}

func (m *Message) Validate(validateFn func(*types.Header, []byte, []byte) (common.Address, error), previousHeader *types.Header) (*types.CommitteeMember, error) {
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

	m.power = v.VotingPower.Uint64()
	return v, nil
}

func (m *Message) Payload() []byte {
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
	var msg string
	if m.Code == msgProposal {
		var proposal Proposal
		err := m.Decode(&proposal)
		if err != nil {
			return ""
		}
		msg = proposal.String()
	}

	if m.Code == msgPrevote || m.Code == msgPrecommit {
		var vote Vote
		err := m.Decode(&vote)
		if err != nil {
			return ""
		}
		msg = vote.String()
	}
	return fmt.Sprintf("{sender: %v, power: %v, msgCode: %v, msg: %v}", m.Address.String(), m.power, m.Code, msg)
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
