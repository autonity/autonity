package message

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/p2p"
	"github.com/autonity/autonity/rlp"
	"golang.org/x/crypto/blake2b"
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

	// LightProposalCode is only used by accountability that it converts full proposal to a lite one
	// which contains just meta-data of a proposal for a sustainable on-chain proof mechanism.
	LightProposalCode
)

type SigVerifier func(*types.Header, []byte, []byte) (common.Address, error)

type Message interface {
	R() int64
	H() uint64
	Code() uint8
	Sender() common.Address
	Power() *big.Int
	String() string
	Hash() common.Hash
	Value() common.Hash
	Payload() []byte
	setPayload([]byte)
	Signature() []byte
	Validate(func(address common.Address) *types.CommitteeMember) error
}

type baseMessage struct {
	// attributes are left private to avoid direct modification
	round     int64
	height    uint64
	signature []byte

	payload        []byte
	signatureInput []any
	power          *big.Int
	sender         common.Address
	hash           common.Hash
}

type Propose struct {
	block      *types.Block
	validRound int64
	baseMessage
}

// extPropose is the actual proposal object being exchanged on the network
// before RLP serialization.
type extPropose struct {
	Code            uint8
	Round           uint64
	Height          uint64
	ValidRound      uint64
	IsValidRoundNil bool
	ProposalBlock   *types.Block
	Signature       []byte
}

func (p *Propose) Code() uint8 {
	return ProposalCode
}

func (p *Propose) Block() *types.Block {
	return p.block
}

func (p *Propose) ValidRound() int64 {
	return p.validRound
}
func (p *Propose) Value() common.Hash {
	return p.block.Hash()
}

func (p *Propose) String() string {
	return fmt.Sprintf("{Round: %v, Height: %v, ValidRound: %v, ProposedBlockHash: %v}",
		p.round, p.H(), p.validRound, p.block.Hash().String())
}

func NewPropose(r int64, h uint64, vr int64, block *types.Block, signer func([]byte) ([]byte, error)) *Propose {
	isValidRoundNil := false
	validRound := uint64(0)
	if vr == -1 {
		isValidRoundNil = true
	} else {
		validRound = uint64(vr)
	}
	// Calculate signature first
	signatureInput := []any{ProposalCode, uint64(r), h, validRound, isValidRoundNil, block.Hash()}
	signatureInputEncoded, _ := rlp.EncodeToBytes(signatureInput)
	signature, _ := signer(signatureInputEncoded)

	payload, _ := rlp.EncodeToBytes(&extPropose{
		Code:            ProposalCode,
		Round:           uint64(r),
		Height:          h,
		ValidRound:      validRound,
		IsValidRoundNil: isValidRoundNil,
		ProposalBlock:   block,
		Signature:       signature,
	})
	return &Propose{
		block:      block,
		validRound: vr,
		baseMessage: baseMessage{
			round:          r,
			height:         h,
			signatureInput: signatureInput,
			signature:      signature,
			payload:        payload,
			hash:           blake2b.Sum256(payload),
		},
	}
}
func (p *Propose) DecodeRLP(s *rlp.Stream) error {
	ext := &extPropose{}
	if err := s.Decode(ext); err != nil {
		return err
	}
	if ext.Code != ProposalCode {
		return constants.ErrInvalidMessage
	}
	if ext.ProposalBlock == nil {
		return constants.ErrInvalidMessage
	}
	if ext.Round > constants.MaxRound || ext.ValidRound > constants.MaxRound {
		return constants.ErrInvalidMessage
	}
	if ext.Height == 0 {
		return constants.ErrInvalidMessage
	}
	if ext.IsValidRoundNil {
		if ext.ValidRound != 0 {
			return constants.ErrInvalidMessage
		}
		p.validRound = -1
	} else {
		p.validRound = int64(ext.ValidRound)
	}
	p.round = int64(ext.Round)
	p.height = ext.Height
	p.block = ext.ProposalBlock
	p.signature = ext.Signature
	p.signatureInput = []any{ProposalCode, ext.Round, ext.Height, ext.ValidRound, ext.IsValidRoundNil, p.block.Hash()}
	return nil
}

type LightProposal struct {
	blockHash  common.Hash
	validRound int64
	baseMessage
}

func (p *LightProposal) Code() uint8 {
	return LightProposalCode
}

func (p *LightProposal) ValidRound() int64 {
	return p.validRound
}

func (p *LightProposal) Value() common.Hash {
	return p.blockHash
}

func (p *LightProposal) String() string {
	return fmt.Sprintf("{sender: %v, power: %v, Code: %v, value: %v}", p.sender.String(), p.power, p.Code(), p.blockHash)
}

func NewLightProposal(proposal *Propose) *LightProposal {
	return &LightProposal{
		blockHash:  proposal.hash,
		validRound: proposal.validRound,
		baseMessage: baseMessage{
			round:     proposal.round,
			height:    proposal.height,
			signature: proposal.signature,
			payload:   nil,
			power:     nil,
			sender:    common.Address{},
			hash:      common.Hash{},
		},
	}
}

// extVote is object being transmitted over the network to carry votes.
type extVote struct {
	// code is redundant with the p2p.msg code however it is required
	// because we don't want to re-serialize the message again in order
	// to compute the hash value.
	// todo: remove the need to hash those values or at least try doing something
	// more efficient.
	Code      uint8
	Round     uint64
	Height    uint64
	Value     common.Hash
	Signature []byte
}

type Prevote struct {
	value common.Hash
	baseMessage
}

func (p *Prevote) Code() uint8 {
	return PrevoteCode
}

func (p *Prevote) Value() common.Hash {
	return p.value
}

func (p *Prevote) String() string {
	return fmt.Sprintf("{sender: %v, power: %v, Code: %v, value: %v}", p.sender.String(), p.power, p.Code(), p.value)
}

type Precommit struct {
	value common.Hash
	baseMessage
}

func (p *Precommit) Code() uint8 {
	return PrecommitCode
}

func (p *Precommit) Value() common.Hash {
	return p.value
}

func (p *Precommit) String() string {
	return fmt.Sprintf("{sender: %v, power: %v, Code: %v, value: %v}", p.sender.String(), p.power, p.Code(), p.value)
}

func NewVote[
	E Prevote | Precommit,
	PE interface {
		*E
		Message
	}](r int64, h uint64, value common.Hash, signer func([]byte) ([]byte, error)) *E {
	code := PE(new(E)).Code()
	// Pay attention that we're adding the message Code to the signature input data.
	signatureInput := []any{code, uint64(r), h, value}
	signatureEncodedInput, _ := rlp.EncodeToBytes(signatureInput)
	signature, _ := signer(signatureEncodedInput)
	payload, _ := rlp.EncodeToBytes(extVote{
		Code:      code,
		Round:     uint64(r),
		Height:    h,
		Value:     value,
		Signature: signature,
	})
	vote := E{
		value: value,
		baseMessage: baseMessage{
			round:          r,
			height:         h,
			signature:      signature,
			payload:        payload,
			hash:           blake2b.Sum256(payload),
			signatureInput: signatureInput,
		},
	}
	return &vote
}

func (p *Prevote) DecodeRLP(s *rlp.Stream) error {
	encoded := extVote{}
	if err := s.Decode(&encoded); err != nil {
		return err
	}
	if encoded.Code != PrevoteCode {
		return constants.ErrFailedDecodePrevote
	}
	p.value = encoded.Value
	p.height = encoded.Height
	if p.height == 0 {
		return constants.ErrInvalidMessage
	}
	p.signature = encoded.Signature
	if encoded.Round > constants.MaxRound {
		return constants.ErrInvalidMessage
	}
	p.round = int64(encoded.Round)
	if p.round < 0 {
		return constants.ErrInvalidMessage
	}
	p.signatureInput = []any{PrevoteCode, encoded.Round, encoded.Height, encoded.Value}
	return nil
}

func (p *Precommit) DecodeRLP(s *rlp.Stream) error {
	encoded := extVote{}
	if err := s.Decode(&encoded); err != nil {
		return err
	}
	if encoded.Code != PrecommitCode {
		return constants.ErrFailedDecodePrevote
	}
	p.value = encoded.Value
	p.height = encoded.Height
	if p.height == 0 {
		return constants.ErrInvalidMessage
	}
	p.signature = encoded.Signature
	if encoded.Round > constants.MaxRound {
		return constants.ErrInvalidMessage
	}
	p.round = int64(encoded.Round)
	if p.round < 0 {
		return constants.ErrInvalidMessage
	}
	p.signatureInput = []any{PrecommitCode, encoded.Round, encoded.Height, encoded.Value}
	return nil
}

func FromWire[T any, PT interface {
	*T
	Message
}](p2pMsg p2p.Msg) (PT, error) {
	message := PT(new(T))
	if err := p2pMsg.Decode(message); err != nil {
		return message, err
	}
	if _, err := p2pMsg.Payload.(*bytes.Reader).Seek(0, io.SeekStart); err != nil {
		return message, err
	}
	payload := make([]byte, p2pMsg.Size)
	if _, err := p2pMsg.Payload.Read(payload); err != nil {
		return message, err
	}
	message.setPayload(payload)
	return message, nil
}

func (b *baseMessage) Sender() common.Address {
	if b.sender == (common.Address{}) {
		panic("coding error (to remove) ")
	}
	return b.sender
}

func (b *baseMessage) H() uint64 {
	return b.height
}

func (b *baseMessage) setPayload(payload []byte) {
	b.payload = payload
	b.hash = blake2b.Sum256(payload)
}

func (b *baseMessage) R() int64 {
	return b.round
}

func (b *baseMessage) Power() *big.Int {
	return b.power
}

func (b *baseMessage) Signature() []byte {
	return b.signature
}

func (b *baseMessage) Payload() []byte {
	return b.payload
}

func (b *baseMessage) Hash() common.Hash {
	return b.hash
}

// Validate verify the signature and set appropriate sender / power fields
func (b *baseMessage) Validate(inCommittee func(address common.Address) *types.CommitteeMember) error {
	// We are not saving the rlp encoded signature input data as we want
	// to avoid this extra-serialization step if the message has already been received
	// The call to Validate() only happen after the cache check in the backend handler.
	sigData, _ := rlp.EncodeToBytes(b.signatureInput)
	hash := blake2b.Sum256(sigData)
	addr, err := tendermint.SigToAddr(hash, b.signature)
	if err != nil {
		return err
	}
	validator := inCommittee(addr)
	if validator == nil {
		return ErrUnauthorizedAddress
	}
	b.sender = addr
	b.power = validator.VotingPower
	return nil
}

// PrepareCommittedSeal returns the input data to compute the committed seal for a given block hash.
func PrepareCommittedSeal(hash common.Hash, round int64, height *big.Int) common.Hash {
	// this is matching the signature input that we get from the committed messages.
	buf, _ := rlp.EncodeToBytes([]any{PrecommitCode, uint64(round), height.Uint64(), hash})
	return blake2b.Sum256(buf)
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
		// don't know why the legacy Code panic here, it introduces live-ness issue of the network.
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
	return fmt.Sprintf("{sender: %v, power: %v, Code: %v, inner: %v}", m.Address.String(), m.Power, m.Code, m.ConsensusMsg)
}
*/
