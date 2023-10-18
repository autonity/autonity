package message

import (
	"errors"
	"fmt"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/rlp"
	"io"
	"math/big"
)

var (
	ErrMsgPayloadNotDecoded = errors.New("payload not decoded")
	ErrBadSignature         = errors.New("bad signature")
	ErrUnauthorizedAddress  = errors.New("unauthorized address")
)

const (
	ProposalCode uint8 = iota
	PrevoteCode
	PrecommitCode
	// MsgLightProposal is only used by accountability that it converts full proposal to a lite one
	// which contains just meta-data of a proposal for a sustainable on-chain proof mechanism.
	MsgLightProposal
)

type SigVerifier func(*types.Header, []byte, []byte) (common.Address, error)

/*type Message struct {
	Signature     []byte
	CommittedSeal []byte // todo(youssef): this should be moved in the precommit object

	// todo:(youssef) this might be rlp encoded too in the message structure even if nil
	Power        *big.Int
	ConsensusMsg ConsensusMsg // cached decoded Msg
	Bytes        []byte       // rlp encoded bytes with only the 1st 5 fields of this Message struct.
}*/

type Message interface {
	R() int64
	H() uint64
	Code() uint8
	Sender() common.Address
	Power() *big.Int
	String() string
	// Hash() common.Hash

}

type baseMessage struct {
	// attributes are left private to avoid direct modification
	round     int64
	height    uint64
	signature []byte

	bytes  []byte
	power  *big.Int
	sender common.Address
	hash   common.Hash
}

type Propose struct {
	block      *types.Block
	validRound int64
	baseMessage
}

type LightProposal struct {
	round      int64
	height     *big.Int
	validRound int64
	value      common.Hash
	signature  []byte // the signature computes from the tuple: (Round, Height, ValidRound, ProposalBlock.Hash())
}

// extPropose is the proposal object being exchanged at the networking level
type extPropose struct {
	round           uint64
	height          uint64
	validRound      uint64
	isValidRoundNil bool
	proposalBlock   *types.Block
	signature       []byte
}

func (p Propose) Code() uint8 {
	return ProposalCode
}

func (p Propose) ValidRound() int64 {
	return p.validRound
}

func (p Propose) String() string {
	return fmt.Sprintf("{Round: %v, Height: %v, ValidRound: %v, ProposedBlockHash: %v}",
		p.round, p.H(), p.validRound, p.block.Hash().String())
}

func (p Propose) EncodeRLP(w io.Writer) error {
	if p.block == nil {
		log.Crit("encoding proposal with empty block")
	}
	// RLP encoding doesn't support negative big.Int, so we have to pass one additional field to represents validRound = -1.
	// Note that we could have as well indexed rounds starting by 1, but we want to stay close as possible to the spec.
	isValidRoundNil := false
	var validRound uint64
	if p.validRound == -1 {
		validRound = 0
		isValidRoundNil = true
	} else {
		validRound = uint64(p.validRound)
	}
	ext := extPropose{
		round:           uint64(p.round),
		height:          p.height,
		validRound:      validRound,
		isValidRoundNil: isValidRoundNil,
		proposalBlock:   p.block,
		signature:       p.signature,
	}

	return rlp.Encode(w, ext)
}

func NewPropose(r int64, h uint64, vr int64, p *types.Block, signer func([]byte) ([]byte, error)) *Propose {
	isValidRoundNil := false
	validRound := uint64(0)
	if vr == -1 {
		isValidRoundNil = true
	} else {
		validRound = uint64(vr)
	}
	rawProposal, _ := rlp.EncodeToBytes([]any{ProposalCode, uint64(r), h, validRound, isValidRoundNil, p.Hash()})
	signature, _ := signer(rawProposal)
	return &Propose{
		block:      p,
		validRound: vr,
		baseMessage: baseMessage{
			round:     r,
			height:    h,
			signature: signature,
		},
	}
}

type Prevote struct {
	Value common.Hash
	baseMessage
}

func (p Prevote) Code() uint8 {
	return PrevoteCode
}

type Precommit struct {
	Value common.Hash
	baseMessage
}

func (p Precommit) Code() uint8 {
	return ProposalCode
}

func (m baseMessage) Sender() common.Address {
	return m.sender
}

func (m baseMessage) H() uint64 {
	return m.Height
}

func (m baseMessage) R() int64 {
	return m.Round
}

func (m baseMessage) Power() *big.Int {
	return m.power
}

/*
// ==============================================
//
// define the functions that needs to be provided for rlp Encoder/Decoder.

// EncodeRLP serializes m into the Ethereum RLP format.
func (m *Message) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []any{m.Code, m.Payload, m.Address, m.Signature, m.CommittedSeal})
}

// Hash Unified the Hash calculation of consensus msg. RLPHash(msg) hashes both public fields and private fields of
// msg, while the rlp.EncodeToBytes(AccountabilityEvent) function, it calls interface EncodeRLP() that is implemented
// by Message struct for only public fields. To keep away the in-consistent of hashing between AFD and precompiled
// contract, we unified the consensus msg hashing in this Hash() function.
func (m *Message) Hash() common.Hash {
	return types.RLPHash(&Message{
		Code:      m.Code,
		Payload:   m.Payload,
		Address:   m.Address,
		Signature: m.Signature,
		// BLSSignature:  m.BLSSignature, leave it at D4 merge.
	})
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
	m.Code, m.Payload, m.Address, m.Signature, m.CommittedSeal = msg.Code, msg.Msg, msg.Address, msg.Signature, msg.CommittedSeal
	return nil
}

func FromBytes(b []byte) (*Message, error) {
	msg := &Message{Bytes: b}
	// Decode message
	if err := rlp.DecodeBytes(b, msg); err != nil {
		return nil, err
	}
	// Decode the payload, this will cache the decoded msg payload.
	return msg, msg.DecodePayload()
}

func (m *Message) DecodePayload() error {
	switch m.Code {
	case consensus.MsgProposal:
		return m.Decode(new(Proposal))
	case consensus.MsgPrevote, consensus.MsgPrecommit:
		return m.Decode(new(Vote))
	default:
		return ErrMsgPayloadNotDecoded
	}
}

func (m *Message) Validate(validateSig SigVerifier, previousHeader *types.Header) error {
	if previousHeader.Number.Uint64()+1 != m.H() {
		// don't know why the legacy code panic here, it introduces live-ness issue of the network.
		// youssef: that is really bad and should never happen, could be because of a race-condition
		// I'm reintroducing the panic to check if this scenario happens in the wild. We must never fail silently.
		panic("Autonity has encountered a problem which led to an inconsistent state, please report this issue.")
		//return fmt.Errorf("inconsistent message verification")
	}
	signature := m.Signature
	payload, err := m.BytesNoSignature()
	if err != nil {
		return err
	}
	if lp, ok := m.ConsensusMsg.(*LightProposal); ok {
		// in the case of a light proposal, the signature that matters is the inner-one.
		payload = lp.BytesNoSignature()
		signature = lp.Signature
	}

	recoveredAddress, err := validateSig(previousHeader, payload, signature)
	if err != nil {
		return err
	}
	// ensure message was signed by the sender
	if m.Address != recoveredAddress {
		return ErrBadSignature
	}
	validator := previousHeader.CommitteeMember(recoveredAddress)
	// validateSig check as well if the header is in the committee, so this seems unnecessary
	if validator == nil {
		return ErrUnauthorizedAddress
	}

	// check if the lite proposal signature inside the proposal is correct or not.
	if proposal, ok := m.ConsensusMsg.(*Proposal); ok {
		if err := proposal.VerifyLightProposalSignature(m.Address); err != nil {
			return err
		}
	}

	m.Power = validator.VotingPower
	return nil
}

func (m *Message) GetBytes() []byte {
	if m.Bytes == nil {
		data, err := rlp.EncodeToBytes(&Message{Code: m.Code, Payload: m.Payload, Address: m.Address, Signature: m.Signature, CommittedSeal: m.CommittedSeal})
		if err != nil {
			// We panic if there is an error, reasons:
			// Either we received the message and we managed to decode it, hence it must be possible to encode it.
			// If we can't encode the payload for our own generated messages, that's a programming error.
			panic("could not decode message payload")
		}
		m.Bytes = data
	}
	return m.Bytes
}

func (m *Message) GetPower() *big.Int {
	return m.Power
}

func (m *Message) BytesNoSignature() ([]byte, error) {
	// youssef: not sure if the returned error is necessary here as we are in control of the object.
	return rlp.EncodeToBytes(&Message{
		Code:          m.Code,
		Payload:       m.Payload,
		Address:       m.Address,
		Signature:     []byte{},
		CommittedSeal: m.CommittedSeal,
	})
}

// Todo(youssef): this function is called from too many places
func (m *Message) Decode(val any) error {
	//Decode is responsible to rlp-decode m.Payload. It is meant to only perform the actual decoding once,
	//saving a cached value in m.decodedMsg.

	rval := reflect.ValueOf(val)
	if rval.Kind() != reflect.Ptr {
		return errors.New("decode arg must be a pointer")
	}

	// check if we already have a cached value decoded
	if m.ConsensusMsg != nil {
		if !rval.Type().AssignableTo(reflect.TypeOf(m.ConsensusMsg)) {
			return errors.New("type mismatch with decoded value")
		}
		rval.Elem().Set(reflect.ValueOf(m.ConsensusMsg).Elem())
		return nil
	}

	err := rlp.DecodeBytes(m.Payload, val)
	if err != nil {
		return err
	}

	// copy the result via Set (shallow)
	nval := reflect.New(rval.Elem().Type()) // we need first to allocate memory
	nval.Elem().Set(rval.Elem())
	m.ConsensusMsg = nval.Interface().(ConsensusMsg)
	return nil
}

// ToLightProposal convert a decoded proposal into a lite proposal for accountability proof, only used by AFD.
func (m *Message) ToLightProposal() *Message {
	var proposal *Proposal
	var ok bool
	if proposal, ok = m.ConsensusMsg.(*Proposal); !ok {
		log.Crit("error creating a light proposal")
	}
	lightProposal := &LightProposal{
		Round:      m.R(),
		Height:     new(big.Int).SetUint64(m.H()),
		ValidRound: proposal.ValidRound,
		Value:      m.Value(),
		Signature:  proposal.LightSignature,
	}
	encoded, _ := rlp.EncodeToBytes(lightProposal)
	message := &Message{
		Code:         consensus.MsgLightProposal,
		Payload:      encoded,
		ConsensusMsg: lightProposal,
		Address:      m.Address,
		// We don't really care about the outer message signature
	}
	return message
}

func (m *Message) String() string {
	return fmt.Sprintf("{sender: %v, power: %v, code: %v, inner: %v}", m.Address.String(), m.Power, m.Code, m.ConsensusMsg)
}
*/
