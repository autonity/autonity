// Package message implements an interface and the three underlying consensus messages types that
// tendermint is using: Propose, Prevote and Precommit.
// In addition to that, we have a special type, the "Light Proposal" which is being used for
// accountability purposes. Light proposals are never directly brodcasted
// over the network but always part of a proof object, defined in the accountability package.
// Messages can exist in two states: unverified and verified depending on their signature verification.
// When verified, calling `Validate` the voting power information and the sender information become available
// There are three ways that a consensus message can be instantiated:
//   - using a "New" constructor, e.g. NewPrevote :
//     The following created object is then fully created, with signature and final payload already
//     pre-computed. Internal state is unverified as voting power information is not available.
//   - using "FromWire": signature and payload object available. State unverified.
//   - using a Fake constructor.
package message

import (
	"errors"
	"fmt"
	"io"
	"math/big"
	"sync"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/rlp"
)

var (
	ErrBadSignature        = errors.New("bad signature")
	ErrUnauthorizedAddress = errors.New("unauthorized address")
)

const (
	ProposalCode uint8 = iota
	PrevoteCode
	PrecommitCode
	LightProposalCode
)

type Signer func(hash common.Hash) (signature []byte, address common.Address)

type Msg interface {
	// Code returns the message code, it must always matching the concrete type.
	Code() uint8

	// R returns the message round.
	R() int64

	// H returns the mesage height.
	H() uint64

	// Value returns the block hash being voted for.
	Value() common.Hash

	// Returns the sender address. This is not available until the message has been validated
	Sender() common.Address

	// Power returns the message voting power.
	Power() *big.Int

	// String returns a string description of the message.
	String() string

	// Hash returns the hash of the messagel. This is not available if the underlying payload
	// hasn't be assigned.
	Hash() common.Hash

	// Payload returns the rlp-encoded payload ready to be broadcasted.
	Payload() []byte

	// Signature returns the signature of this message
	Signature() []byte

	// Validate execute the message's signature verification, assign the power and the sender.
	Validate(func(address common.Address) *types.CommitteeMember) error
}

type base struct {
	// attributes are left private to avoid direct modification
	round     int64
	height    uint64
	signature []byte

	payload        []byte
	signatureInput []any
	power          *big.Int
	sender         common.Address
	hash           common.Hash
	verified       bool
	sync.RWMutex   // To remove once we can merge the parralel signature verification work.
}

type Propose struct {
	block      *types.Block
	validRound int64
	base
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

func (p *Propose) MustVerify(inCommittee func(address common.Address) *types.CommitteeMember) *Propose {
	if err := p.Validate(inCommittee); err != nil {
		panic("validation failed")
	}
	return p
}

func (p *Propose) ToLight() *LightProposal {
	return NewLightProposal(p)
}

func NewPropose(r int64, h uint64, vr int64, block *types.Block, signer Signer) *Propose {
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
	signature, validator := signer(crypto.Hash(signatureInputEncoded))

	payload, _ := rlp.EncodeToBytes(&extPropose{
		Code:            ProposalCode,
		Round:           uint64(r),
		Height:          h,
		ValidRound:      validRound,
		IsValidRoundNil: isValidRoundNil,
		ProposalBlock:   block,
		Signature:       signature,
	})
	// we don't need to assign here the voting power neither the sender as they are going to be retrieved
	// after a Validate() call during processing.
	return &Propose{
		block:      block,
		validRound: vr,
		base: base{
			round:          r,
			height:         h,
			signatureInput: signatureInput,
			sender:         validator,
			signature:      signature,
			payload:        payload,
			hash:           crypto.Hash(payload),
			verified:       false,
		},
	}
}
func (p *Propose) DecodeRLP(s *rlp.Stream) error {
	payload, err := s.Raw()
	if err != nil {
		return err
	}
	ext := &extPropose{}
	if err := rlp.DecodeBytes(payload, ext); err != nil {
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
	p.payload = payload
	p.hash = crypto.Hash(payload)
	return nil
}

type LightProposal struct {
	blockHash  common.Hash
	validRound int64
	base
}
type extLightProposal struct {
	Code            uint8
	Round           uint64
	Height          uint64
	ValidRound      uint64
	IsValidRoundNil bool
	ProposalBlock   common.Hash
	Signature       []byte
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
	if !proposal.verified {
		//temporary panic to catch bugs.
		panic("unverified light-proposal creation")
	}
	isValidRoundNil := false
	validRound := uint64(0)
	if proposal.validRound == -1 {
		isValidRoundNil = true
	} else {
		validRound = uint64(proposal.validRound)
	}
	payload, _ := rlp.EncodeToBytes(extLightProposal{
		Code:            LightProposalCode,
		Round:           uint64(proposal.round),
		Height:          proposal.height,
		ValidRound:      validRound,
		IsValidRoundNil: isValidRoundNil,
		ProposalBlock:   proposal.block.Hash(),
		Signature:       proposal.signature,
	})
	return &LightProposal{
		blockHash:  proposal.Block().Hash(),
		validRound: proposal.validRound,
		base: base{
			round:          proposal.round,
			height:         proposal.height,
			signature:      proposal.signature,
			signatureInput: proposal.signatureInput,
			payload:        payload,
			power:          proposal.Power(),
			sender:         proposal.sender,
			hash:           crypto.Hash(payload),
			verified:       true,
		},
	}
}

func (p *LightProposal) DecodeRLP(s *rlp.Stream) error {
	payload, err := s.Raw()
	if err != nil {
		return err
	}
	ext := &extLightProposal{}
	if err := rlp.DecodeBytes(payload, ext); err != nil {
		return err
	}
	if ext.Code != LightProposalCode {
		return constants.ErrInvalidMessage
	}
	if ext.ProposalBlock == (common.Hash{}) {
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
	p.blockHash = ext.ProposalBlock
	p.signature = ext.Signature
	p.signatureInput = []any{ProposalCode, ext.Round, ext.Height, ext.ValidRound, ext.IsValidRoundNil, p.blockHash}
	p.payload = payload
	p.hash = crypto.Hash(payload)
	return nil
}

// extVote is object being transmitted over the network to carry votes.
type extVote struct {
	// Code is redundant with the p2p.msg code however it is required
	// because we don't want to re-serialize the message again in order
	// to compute the hash value.
	Code      uint8
	Round     uint64
	Height    uint64
	Value     common.Hash
	Signature []byte
}

type Prevote struct {
	value common.Hash
	base
}

func (p *Prevote) Code() uint8 {
	return PrevoteCode
}

func (p *Prevote) Value() common.Hash {
	return p.value
}

func (p *Prevote) MustVerify(inCommittee func(address common.Address) *types.CommitteeMember) *Prevote {
	p.Validate(inCommittee)
	return p
}

func (p *Prevote) String() string {
	return fmt.Sprintf("{r:  %v, h: %v , sender: %v, power: %v, Code: %v, value: %v}",
		p.round, p.height, p.sender, p.power, p.Code(), p.value)
}

type Precommit struct {
	value common.Hash
	base
}

func (p *Precommit) Code() uint8 {
	return PrecommitCode
}

func (p *Precommit) Value() common.Hash {
	return p.value
}

func (p *Precommit) MustVerify(inCommittee func(address common.Address) *types.CommitteeMember) *Precommit {
	p.Validate(inCommittee)
	return p
}

func (p *Precommit) String() string {
	return fmt.Sprintf("{r:  %v, h: %v , sender: %v, power: %v, Code: %v, value: %v}",
		p.round, p.height, p.sender, p.power, p.Code(), p.value)
}

func newVote[
	E Prevote | Precommit,
	PE interface {
		*E
		Msg
	}](r int64, h uint64, value common.Hash, signer Signer) *E {
	code := PE(new(E)).Code()
	// Pay attention that we're adding the message Code to the signature input data.
	signatureInput := []any{code, uint64(r), h, value}
	signatureEncodedInput, _ := rlp.EncodeToBytes(signatureInput)
	signature, validator := signer(crypto.Hash(signatureEncodedInput))
	payload, _ := rlp.EncodeToBytes(extVote{
		Code:      code,
		Round:     uint64(r),
		Height:    h,
		Value:     value,
		Signature: signature,
	})
	vote := E{
		value: value,
		base: base{
			round:          r,
			height:         h,
			signature:      signature,
			sender:         validator,
			payload:        payload,
			hash:           crypto.Hash(payload),
			signatureInput: signatureInput,
			verified:       false,
		},
	}
	return &vote
}

func NewPrevote(r int64, h uint64, value common.Hash, signer Signer) *Prevote {
	return newVote[Prevote](r, h, value, signer)
}

func NewPrecommit(r int64, h uint64, value common.Hash, signer Signer) *Precommit {
	return newVote[Precommit](r, h, value, signer)
}

func (p *Prevote) DecodeRLP(s *rlp.Stream) error {
	payload, err := s.Raw()
	if err != nil {
		return err
	}
	encoded := &extVote{}
	if err := rlp.DecodeBytes(payload, encoded); err != nil {
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
	p.payload = payload
	p.hash = crypto.Hash(payload)
	return nil
}

func (p *Precommit) DecodeRLP(s *rlp.Stream) error {
	payload, err := s.Raw()
	if err != nil {
		return err
	}
	encoded := &extVote{}
	if err := rlp.DecodeBytes(payload, encoded); err != nil {
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
	p.payload = payload
	p.hash = crypto.Hash(payload)
	return nil
}

func (b *base) Sender() common.Address {
	b.RLock()
	defer b.RUnlock()
	if !b.verified {
		panic("unverified message")
	}
	return b.sender
}

func (b *base) H() uint64 {
	return b.height
}

func (b *base) EncodeRLP(w io.Writer) error {
	_, err := w.Write(b.payload)
	return err
}

func (b *base) R() int64 {
	return b.round
}

func (b *base) Power() *big.Int {
	b.RLock()
	defer b.RUnlock()
	if !b.verified {
		panic("unverified message")
	}
	return b.power
}

func (b *base) Signature() []byte {
	return b.signature
}

func (b *base) Payload() []byte {
	return b.payload
}

func (b *base) Hash() common.Hash {
	return b.hash
}

// Validate verify the signature and set appropriate sender / power fields
func (b *base) Validate(inCommittee func(address common.Address) *types.CommitteeMember) error {
	b.Lock()
	defer b.Unlock()
	if b.verified {
		return nil
	}
	// We are not saving the rlp encoded signature input data as we want
	// to avoid this extra-serialization step if the message has already been received
	// The call to Validate() only happen after the cache check in the backend handler.
	sigData, _ := rlp.EncodeToBytes(b.signatureInput)
	hash := crypto.Hash(sigData)
	addr, err := tendermint.SigToAddr(hash, b.signature)
	if err != nil {
		return ErrBadSignature
	}
	validator := inCommittee(addr)
	if validator == nil {
		return ErrUnauthorizedAddress
	}
	b.sender = addr
	b.power = validator.VotingPower
	b.verified = true
	return nil
}

// PrepareCommittedSeal returns the input data to compute the committed seal for a given block hash.
func PrepareCommittedSeal(hash common.Hash, round int64, height *big.Int) common.Hash {
	// this is matching the signature input that we get from the committed messages.
	buf, _ := rlp.EncodeToBytes([]any{PrecommitCode, uint64(round), height.Uint64(), hash})
	return crypto.Hash(buf)
}

// Fake is a dummy object used for internal testing.
type Fake struct {
	FakeCode      uint8
	FakeRound     int64
	FakeHeight    uint64
	FakeValue     common.Hash
	FakePayload   []byte
	FakeHash      common.Hash
	FakeSender    common.Address
	FakeSignature []byte
	FakePower     *big.Int
}

func (f Fake) R() int64                                                       { return f.FakeRound }
func (f Fake) H() uint64                                                      { return f.FakeHeight }
func (f Fake) Code() uint8                                                    { return f.FakeCode }
func (f Fake) Sender() common.Address                                         { return f.FakeSender }
func (f Fake) Power() *big.Int                                                { return f.FakePower }
func (f Fake) String() string                                                 { return "{fake}" }
func (f Fake) Hash() common.Hash                                              { return f.FakeHash }
func (f Fake) Value() common.Hash                                             { return common.Hash{} }
func (f Fake) Payload() []byte                                                { return f.FakePayload }
func (f Fake) Signature() []byte                                              { return f.FakeSignature }
func (f Fake) Validate(_ func(_ common.Address) *types.CommitteeMember) error { return nil }

func NewFakePrevote(f Fake) *Prevote {
	return &Prevote{
		value: f.FakeValue,
		base: base{
			round:     f.FakeRound,
			height:    f.FakeHeight,
			signature: f.FakeSignature,
			payload:   f.FakePayload,
			power:     f.FakePower,
			sender:    f.FakeSender,
			hash:      f.FakeHash,
			verified:  true,
		},
	}
}

func NewFakePrecommit(f Fake) *Precommit {
	return &Precommit{
		value: f.FakeValue,
		base: base{
			round:     f.FakeRound,
			height:    f.FakeHeight,
			signature: f.FakeSignature,
			payload:   f.FakePayload,
			power:     f.FakePower,
			sender:    f.FakeSender,
			hash:      f.FakeHash,
		},
	}
}
