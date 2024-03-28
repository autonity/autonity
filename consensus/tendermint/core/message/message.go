// Package message implements an interface and the three underlying consensus messages types that
// tendermint is using: Propose, Prevote and Precommit.
// In addition to that, we have a special type, the "Light Proposal" which is being used for
// accountability purposes. Light proposals are never directly brodcasted
// over the network but always part of a proof object, defined in the accountability package.
// Messages can exist in two states: unverified and verified depending on their signature verification.
// When verified, calling `Validate` the voting power information becomes available, and the sender can be relied upon.
// There are three ways that a consensus message can be instantiated:
//   - using a "New" constructor, e.g. NewPrevote :
//     The object is fully created, with signature and final payload already
//     pre-computed. Internal state is unverified as voting power information is not available.
//   - decoding a RLP-encoded message from the wire. State unverified.
//   - using a Fake constructor.
package message

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/rlp"
)

var (
	ErrBadSignature        = errors.New("bad signature")
	ErrUnauthorizedAddress = errors.New("unauthorized address")
	ErrInvalidSenders      = errors.New("invalid senders information")
)

const (
	ProposalCode uint8 = iota
	PrevoteCode
	PrecommitCode
	LightProposalCode
	AggregatePrevoteCode
	AggregatePrecommitCode
)

type Signer func(hash common.Hash) (signature blst.Signature, address common.Address) //TODO(lorenzo) add power

type Propose struct {
	block      *types.Block
	validRound int64
	individualMsg
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
	// since we do not have ecrecover with BLS signatures, we need to also send the sender in the message.
	// It is sent not-signed to facilitate aggregation.
	// If tampered with, the signature will fail anyways because we will fetch the wrong key.
	Sender    common.Address
	Signature *blst.BlsSignature
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
	return fmt.Sprintf("{%s, ValidRound: %v, ProposedBlockHash: %v}",
		p.base.String(), p.validRound, p.block.Hash().String())
}

/*
func (p *Propose) MustVerify(header *types.Header) *Propose {
	if err := p.Validate(header); err != nil {
		panic("validation failed")
	}
	return p
}*/

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
	signaturePayload, _ := rlp.EncodeToBytes([]any{ProposalCode, uint64(r), h, validRound, isValidRoundNil, block.Hash()})
	signatureInput := crypto.Hash(signaturePayload)
	signature, validator := signer(signatureInput)

	payload, _ := rlp.EncodeToBytes(&extPropose{
		Code:            ProposalCode,
		Round:           uint64(r),
		Height:          h,
		ValidRound:      validRound,
		IsValidRoundNil: isValidRoundNil,
		ProposalBlock:   block,
		Sender:          validator,
		Signature:       signature.(*blst.BlsSignature),
	})

	// we don't need to assign here the voting power as it is going to be retrieved
	// after a PreValidate() call during processing.
	return &Propose{
		block:      block,
		validRound: vr,
		individualMsg: individualMsg{
			sender: validator,
			base: base{
				round:          r,
				height:         h,
				signatureInput: signatureInput,
				signature:      signature,
				payload:        payload,
				hash:           crypto.Hash(payload),
				verified:       false,
			},
		},
	}
}

/* //TODO(lorenzo) for later
// used in tests. Simulates an unverified proposal coming from a remote peer
func NewUnverifiedPropose(r int64, h uint64, vr int64, block *types.Block, signer Signer) *Propose {
	propose := NewPropose(r, h, vr, block, signer)
	propose.unvalidate()
	return propose
}*/

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
	if ext.Signature == nil {
		return constants.ErrInvalidMessage
	}
	if ext.Round > constants.MaxRound || ext.ValidRound > constants.MaxRound {
		return constants.ErrInvalidMessage
	}
	if ext.Height == 0 {
		return constants.ErrInvalidMessage
	}
	if ext.Height != ext.ProposalBlock.NumberU64() {
		return constants.ErrInvalidMessage
	}
	if ext.IsValidRoundNil {
		if ext.ValidRound != 0 {
			return constants.ErrInvalidMessage
		}
		p.validRound = -1
	} else {
		if ext.ValidRound >= ext.Round {
			return constants.ErrInvalidMessage
		}
		p.validRound = int64(ext.ValidRound)
	}

	p.round = int64(ext.Round)
	p.height = ext.Height
	p.block = ext.ProposalBlock
	p.sender = ext.Sender
	p.signature = ext.Signature
	p.payload = payload
	// precompute hash and signature hash
	signaturePayload, _ := rlp.EncodeToBytes([]any{ProposalCode, ext.Round, ext.Height, ext.ValidRound, ext.IsValidRoundNil, p.block.Hash()})
	p.signatureInput = crypto.Hash(signaturePayload)
	p.hash = crypto.Hash(payload)
	return nil
}

type LightProposal struct {
	blockHash  common.Hash
	validRound int64
	individualMsg
}
type extLightProposal struct {
	Code            uint8
	Round           uint64
	Height          uint64
	ValidRound      uint64
	IsValidRoundNil bool
	ProposalBlock   common.Hash
	Sender          common.Address
	Signature       *blst.BlsSignature
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

// TODO(lorenzo) would be useful to print also sender and power, but we need to make sure they are trsuted (verified)
func (p *LightProposal) String() string {
	return fmt.Sprintf("{%s, Code: %v, value: %v}", p.base.String(), p.Code(), p.blockHash)
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
		Sender:          proposal.sender,
		Signature:       proposal.signature.(*blst.BlsSignature),
	})
	return &LightProposal{
		blockHash:  proposal.Block().Hash(),
		validRound: proposal.validRound,
		individualMsg: individualMsg{
			sender: proposal.sender,
			base: base{
				round:          proposal.round,
				height:         proposal.height,
				signature:      proposal.signature,
				signatureInput: proposal.signatureInput,
				payload:        payload,
				power:          proposal.Power(),
				hash:           crypto.Hash(payload),
				verified:       true,
			},
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
	if ext.Signature == nil {
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
		if ext.ValidRound >= ext.Round {
			return constants.ErrInvalidMessage
		}
		p.validRound = int64(ext.ValidRound)
	}
	p.round = int64(ext.Round)
	p.height = ext.Height
	p.blockHash = ext.ProposalBlock
	p.sender = ext.Sender
	p.signature = ext.Signature
	p.payload = payload
	// precompute hash and signature hash
	signaturePayload, _ := rlp.EncodeToBytes([]any{ProposalCode, ext.Round, ext.Height, ext.ValidRound, ext.IsValidRoundNil, p.blockHash})
	p.signatureInput = crypto.Hash(signaturePayload)
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
	Sender    common.Address
	Signature *blst.BlsSignature
}

type Prevote struct {
	value common.Hash
	individualMsg
}

func (p *Prevote) Code() uint8 {
	return PrevoteCode
}

func (p *Prevote) Value() common.Hash {
	return p.value
}

/*
func (p *Prevote) MustVerify(header *types.Header) *Prevote {
	if err := p.Validate(header); err != nil {
		panic("verification failed")
	}
	return p
}*/

func (p *Prevote) String() string {
	return fmt.Sprintf("{r:  %v, h: %v , sender: %v, power: %v, Code: %v, value: %v}",
		p.round, p.height, p.sender, p.power, p.Code(), p.value)
}

type Precommit struct {
	value common.Hash
	individualMsg
}

func (p *Precommit) Code() uint8 {
	return PrecommitCode
}

func (p *Precommit) Value() common.Hash {
	return p.value
}

/*
func (p *Precommit) MustVerify(header *types.Header) *Precommit {
	if err := p.Validate(header); err != nil {
		panic("verification failed")
	}
	return p
}*/

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
	signaturePayload, _ := rlp.EncodeToBytes([]any{code, uint64(r), h, value})
	signatureInput := crypto.Hash(signaturePayload)
	signature, validator := signer(signatureInput)

	payload, _ := rlp.EncodeToBytes(extVote{
		Code:      code,
		Round:     uint64(r),
		Height:    h,
		Value:     value,
		Sender:    validator,
		Signature: signature.(*blst.BlsSignature),
	})
	vote := E{
		value: value,
		individualMsg: individualMsg{
			sender: validator,
			base: base{
				round:          r,
				height:         h,
				signature:      signature,
				payload:        payload,
				hash:           crypto.Hash(payload),
				signatureInput: signatureInput,
				verified:       false,
			},
		},
	}
	return &vote
}

func NewPrevote(r int64, h uint64, value common.Hash, signer Signer) *Prevote {
	return newVote[Prevote](r, h, value, signer)
}

/* //TODO(lorenzo) for later
// used in tests. Simulates an unverified prevote coming from a remote peer
func NewUnverifiedPrevote(r int64, h uint64, value common.Hash, signer Signer) *Prevote {
	prevote := newVote[Prevote](r, h, value, signer)
	prevote.unvalidate()
	return prevote
}*/

func NewPrecommit(r int64, h uint64, value common.Hash, signer Signer) *Precommit {
	return newVote[Precommit](r, h, value, signer)
}

/* //TODO(lorenzo) for alter
func NewUnverifiedPrecommit(r int64, h uint64, value common.Hash, signer Signer) *Precommit {
	precommit := newVote[Precommit](r, h, value, signer)
	precommit.unvalidate()
	return precommit
}*/

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
		return constants.ErrInvalidMessage
	}
	if encoded.Signature == nil {
		return constants.ErrInvalidMessage
	}
	if encoded.Height == 0 {
		return constants.ErrInvalidMessage
	}
	if encoded.Round > constants.MaxRound {
		return constants.ErrInvalidMessage
	}
	p.height = encoded.Height
	p.round = int64(encoded.Round)
	p.value = encoded.Value
	p.sender = encoded.Sender
	p.signature = encoded.Signature
	p.payload = payload
	// precompute hash and signature hash
	signaturePayload, _ := rlp.EncodeToBytes([]any{PrevoteCode, encoded.Round, encoded.Height, encoded.Value})
	p.signatureInput = crypto.Hash(signaturePayload)
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
		return constants.ErrInvalidMessage
	}
	if encoded.Signature == nil {
		return constants.ErrInvalidMessage
	}
	if encoded.Height == 0 {
		return constants.ErrInvalidMessage
	}
	if encoded.Round > constants.MaxRound {
		return constants.ErrInvalidMessage
	}
	p.height = encoded.Height
	p.round = int64(encoded.Round)
	p.value = encoded.Value
	p.sender = encoded.Sender
	p.signature = encoded.Signature
	p.payload = payload
	// precompute hash and signature hash
	signaturePayload, _ := rlp.EncodeToBytes([]any{PrecommitCode, encoded.Round, encoded.Height, encoded.Value})
	p.signatureInput = crypto.Hash(signaturePayload)
	p.hash = crypto.Hash(payload)
	return nil
}

type individualMsg struct {
	// node address of the sender, populated at decoding phase
	sender common.Address
	// index of the sender in the committee, populated at PreValidate phase
	senderIndex int
	base
}

func (im *individualMsg) Sender() common.Address {
	return im.sender
}

func (im *individualMsg) SenderIndex() int {
	return im.senderIndex
}

func (im *individualMsg) PreValidate(header *types.Header) error {
	validator := header.CommitteeMember(im.sender)
	if validator == nil {
		return ErrUnauthorizedAddress
	}

	im.senderKey = validator.ConsensusKey
	im.senderIndex = int(validator.Index)
	im.power = validator.VotingPower
	return nil
}

// PrepareCommittedSeal returns the input data to compute the committed seal for a given block hash.
func PrepareCommittedSeal(hash common.Hash, round int64, height *big.Int) common.Hash {
	// this is matching the signature input that we get from the committed messages.
	buf, _ := rlp.EncodeToBytes([]any{PrecommitCode, uint64(round), height.Uint64(), hash})
	return crypto.Hash(buf)
}

// TODO(lorenzo) update fake msg
// Fake is a dummy object used for internal testing.
type Fake struct {
	FakeCode      uint8
	FakeRound     int64
	FakeHeight    uint64
	FakeValue     common.Hash
	FakePayload   []byte
	FakeHash      common.Hash
	FakeSender    common.Address
	FakeSignature blst.Signature
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
func (f Fake) Signature() blst.Signature                                      { return f.FakeSignature }
func (f Fake) Validate(_ func(_ common.Address) *types.CommitteeMember) error { return nil }

func NewFakePrevote(f Fake) *Prevote {
	return &Prevote{
		value: f.FakeValue,
		individualMsg: individualMsg{
			sender: f.FakeSender,
			base: base{
				round:     f.FakeRound,
				height:    f.FakeHeight,
				signature: f.FakeSignature,
				payload:   f.FakePayload,
				power:     f.FakePower,
				hash:      f.FakeHash,
				verified:  true,
			},
		},
	}
}

func NewFakePrecommit(f Fake) *Precommit {
	return &Precommit{
		value: f.FakeValue,
		individualMsg: individualMsg{
			sender: f.FakeSender,
			base: base{
				round:     f.FakeRound,
				height:    f.FakeHeight,
				signature: f.FakeSignature,
				payload:   f.FakePayload,
				power:     f.FakePower,
				hash:      f.FakeHash,
				verified:  true,
			},
		},
	}
}
