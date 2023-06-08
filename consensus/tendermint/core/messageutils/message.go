package messageutils

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/rlp"
	"io"
	"math/big"
	"reflect"
)

var (
	ErrMsgPayloadNotDecoded = errors.New("message not decoded")
	ErrUnauthorizedAddress  = errors.New("unauthorized address")
)

type Message struct {
	Code          uint8
	TbftMsgBytes  []byte // rlp encoded tendermint msgs: proposal, prevote, precommit and lite proposal only for accountability.
	Address       common.Address
	Signature     []byte
	CommittedSeal []byte

	Power      *big.Int
	DecodedMsg ConsensusMsg // cached decoded Msg
	Bytes      []byte       // rlp encoded bytes with only the 1st 5 fields of this Message struct.
}

// ==============================================
//
// define the functions that needs to be provided for rlp Encoder/Decoder.

// EncodeRLP serializes m into the Ethereum RLP format.
func (m *Message) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.Code, m.TbftMsgBytes, m.Address, m.Signature, m.CommittedSeal})
}

// MsgHash Unified the Hash calculation of consensus msg. RLPHash(msg) hashes both public fields and private fields of
// msg, while the rlp.EncodeToBytes(AccountabilityEvent) function, it calls interface EncodeRLP() that is implemented
// by Message struct for only public fields. To keep away the in-consistent of hashing between AFD and precompiled
// contract, we unified the consensus msg hashing in this MsgHash() function.
func (m *Message) MsgHash() common.Hash {
	return types.RLPHash(&Message{
		Code:         m.Code,
		TbftMsgBytes: m.TbftMsgBytes,
		Address:      m.Address,
		Signature:    m.Signature,
		// BLSSignature:  m.BLSSignature, leave it at D4 merge.
	})
}

func (m *Message) GetSignature() []byte {
	return m.Signature
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (m *Message) DecodeRLP(s *rlp.Stream) error {
	var msg struct {
		Code          uint8
		Msg           []byte
		Address       common.Address
		Signature     []byte
		CommittedSeal []byte
	}

	if err := s.Decode(&msg); err != nil {
		return err
	}
	m.Code, m.TbftMsgBytes, m.Address, m.Signature, m.CommittedSeal = msg.Code, msg.Msg, msg.Address, msg.Signature, msg.CommittedSeal
	return nil
}

// ==============================================
//
// define the functions that needs to be provided for core.

func (m *Message) FromPayload(b []byte) error {
	m.Bytes = b
	// Decode message
	err := rlp.DecodeBytes(b, m)
	if err != nil {
		return err
	}
	// Decode the payload, this will cache the decoded msg payload.
	switch m.Code {
	case consensus.MsgProposal:
		var proposal Proposal
		return m.Decode(&proposal)
	case consensus.MsgPrevote, consensus.MsgPrecommit:
		var vote Vote
		return m.Decode(&vote)
	default:
		return ErrMsgPayloadNotDecoded
	}
}

func (m *Message) Validate(validateFn func(*types.Header, []byte, []byte) (common.Address, error), previousHeader *types.Header) (*types.CommitteeMember, error) {
	// Validate message (on a message without Signature)
	msgHeight, err := m.Height()
	if err != nil {
		return nil, err
	}
	if previousHeader.Number.Uint64()+1 != msgHeight.Uint64() {
		return nil, fmt.Errorf("inconsistent message verification")
		// don't know why the legacy code panic here, it introduces live-ness issue of the network.
		// panic("inconsistent message verification")
	}

	// Still return the message even the err is not nil
	var payload []byte
	var signature []byte

	if m.Type() != consensus.MsgLiteProposal {
		signature = m.Signature
		if payload, err = m.PayloadNoSig(); err != nil {
			return nil, err
		}
	} else {
		var lite = &LiteProposal{
			Round:      m.R(),
			Height:     msgHeight,
			ValidRound: m.ValidRound(),
			Value:      m.Value(),
		}
		signature = m.LiteSig()
		if payload, err = lite.PayloadNoSig(); err != nil {
			return nil, err
		}
	}

	addr, err := validateFn(previousHeader, payload, signature)
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

	m.Power = v.VotingPower
	return v, nil
}

func (m *Message) GetPayload() []byte {
	if m.Bytes == nil {
		payload, err := rlp.EncodeToBytes(m)
		if err != nil {
			// We panic if there is an error, reasons:
			// Either we received the message and we managed to decode it, hence it must be possible to encode it.
			// If we can't encode the payload for our own generated messages, that's a programming error.
			panic("could not decode message payload")
		}
		m.Bytes = payload
	}
	return m.Bytes
}

func (m *Message) GetPower() *big.Int {
	return m.Power
}

func (m *Message) PayloadNoSig() ([]byte, error) {
	return rlp.EncodeToBytes(&Message{
		Code:          m.Code,
		TbftMsgBytes:  m.TbftMsgBytes,
		Address:       m.Address,
		Signature:     []byte{},
		CommittedSeal: m.CommittedSeal,
	})
}

func (m *Message) Decode(val interface{}) error {
	//Decode is responsible to rlp-decode m.TbftMsgBytes. It is meant to only perform the actual decoding once,
	//saving a cached value in m.decodedMsg.

	rval := reflect.ValueOf(val)
	if rval.Kind() != reflect.Ptr {
		return errors.New("decode arg must be a pointer")
	}

	// check if we already have a cached value decoded
	if m.DecodedMsg != nil {
		if !rval.Type().AssignableTo(reflect.TypeOf(m.DecodedMsg)) {
			return errors.New("type mismatch with decoded value")
		}
		rval.Elem().Set(reflect.ValueOf(m.DecodedMsg).Elem())
		return nil
	}

	err := rlp.DecodeBytes(m.TbftMsgBytes, val)
	if err != nil {
		return err
	}

	// copy the result via Set (shallow)
	nval := reflect.New(rval.Elem().Type()) // we need first to allocate memory
	nval.Elem().Set(rval.Elem())
	m.DecodedMsg = nval.Interface().(ConsensusMsg)
	return nil
}

func (m *Message) String() string {
	var msg string
	if m.Code == consensus.MsgProposal {
		var proposal Proposal
		err := m.Decode(&proposal)
		if err != nil {
			return ""
		}
		msg = proposal.String()
	}

	if m.Code == consensus.MsgPrevote || m.Code == consensus.MsgPrecommit {
		var vote Vote
		err := m.Decode(&vote)
		if err != nil {
			return ""
		}
		msg = vote.String()
	}
	return fmt.Sprintf("{sender: %v, power: %v, msgCode: %v, msg: %v}", m.Address.String(), m.Power, m.Code, msg)
}

func (m *Message) Round() (int64, error) {
	if m.DecodedMsg == nil {
		return 0, ErrMsgPayloadNotDecoded
	}
	return m.DecodedMsg.R(), nil
}

func (m *Message) Height() (*big.Int, error) {
	if m.DecodedMsg == nil {
		return nil, ErrMsgPayloadNotDecoded
	}
	return m.DecodedMsg.H(), nil
}

func (m *Message) R() int64 {
	r, err := m.Round()
	// msg should be decoded, it shouldn't be an error.
	if err != nil {
		panic(err)
	}
	return r
}

func (m *Message) H() uint64 {
	h, err := m.Height()
	if err != nil {
		panic(err)
	}
	return h.Uint64()
}

func (m *Message) Sender() common.Address {
	return m.Address
}

func (m *Message) Type() uint8 {
	return m.Code
}

func (m *Message) Value() common.Hash {
	return m.DecodedMsg.V()
}

func (m *Message) ValidRound() int64 {
	if m.Code == consensus.MsgProposal {
		proposal, ok := m.DecodedMsg.(*Proposal)
		if !ok {
			panic("Only proposal message has valid round")
		}

		return proposal.VR()
	}

	proposal, ok := m.DecodedMsg.(*LiteProposal)
	if !ok {
		panic("Only proposal message has valid round")
	}

	return proposal.VR()
}

// ToLiteProposal convert a decoded proposal into a lite proposal for accountability proof, only used by AFD.
func (m *Message) ToLiteProposal() *Message {
	if m.DecodedMsg == nil {
		return nil
	}

	var liteProposal = &LiteProposal{
		Round:      m.R(),
		Height:     new(big.Int).SetUint64(m.H()),
		ValidRound: m.ValidRound(),
		Value:      m.Value(),
		Signature:  m.LiteSig(),
	}

	encoded, err := rlp.EncodeToBytes(liteProposal)
	if err != nil {
		return nil
	}

	var liteMsg = &Message{
		Code:         consensus.MsgLiteProposal,
		TbftMsgBytes: encoded,
		Address:      m.Address,
	}
	// decode it to init the decodedMsg interface.
	var lp LiteProposal
	if err = liteMsg.Decode(&lp); err != nil {
		return nil
	}
	return liteMsg
}

func (m *Message) LiteSig() []byte {
	switch msg := m.DecodedMsg.(type) {
	case *LiteProposal:
		return msg.Sig()
	case *Proposal:
		return msg.LiteSignature()
	default:
		panic("wrong type casting for lite signature")
	}
}

func (m *Message) BadProposer() common.Address {
	vote, ok := m.DecodedMsg.(*Vote)
	if !ok {
		panic("Only vote message vote for bad proposal")
	}
	return vote.MaliciousProposer
}

func (m *Message) BadValue() common.Hash {
	vote, ok := m.DecodedMsg.(*Vote)
	if !ok {
		panic("Only vote message vote for bad proposal")
	}
	return vote.MaliciousValue
}

// ==============================================
//
// helper functions

func Encode(val interface{}) ([]byte, error) {
	return rlp.EncodeToBytes(val)
}
