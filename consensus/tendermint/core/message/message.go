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
	"io"
	"math/big"
	"sync"

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
type memberLookup func(address common.Address) *types.CommitteeMember                 //TODO(lorenzo) maybe not needed

// Msg can represent both:
// 1. an individual message as sent by a validator
// 2. an aggregate message, resulting from bls signature aggregation of multiple individual messages
type Msg interface {
	// Code returns the message code, it must always match the concrete type.
	Code() uint8

	// R returns the message round.
	R() int64

	// H returns the message height.
	H() uint64

	// Value returns the block hash being voted for.
	Value() common.Hash

	// Returns the sender address. This is not available until the message has been validated.
	// the sender is actually populated at decoding, but it cannot be relied upon until after signature verification.
	//Senders() []common.Address

	// Power returns the message voting power.
	Power() *big.Int

	// String returns a string description of the message.
	String() string

	// Hash returns the hash of the message. This is not available if the underlying payload
	// hasn't be assigned.
	Hash() common.Hash

	// Payload returns the rlp-encoded payload ready to be broadcasted.
	Payload() []byte

	// Signature returns the signature of this message
	Signature() blst.Signature

	// Validate execute the message's signature verification, cryptographically verifying the sender and assigning the power value.
	PreValidate(header *types.Header) error

	SignatureHash() common.Hash
	SenderKey() blst.PublicKey
}

// non-aggregated
type IndividualMsg interface {
	Sender() common.Address
	Index() uint64
	Msg
}

type AggregateMsg interface {
	Senders() []common.Address
	SendersCoeff() Coefficients //TODO(lorenzo) name sucks
	Powers() []*big.Int
	Validate() error //TODO(lorenzo) aggregated votes are verified right away in aggregatro
	Msg
}

type base struct {
	// attributes are left private to avoid direct modification
	round     int64
	height    uint64
	signature blst.Signature

	payload        []byte
	signatureInput []any
	signatureHash  common.Hash
	power          *big.Int
	hash           common.Hash
	verified       bool           //TODO(lorenzo) do we need this?
	senderKey      blst.PublicKey // populated at prevalidate step
	sync.RWMutex                  //TODO(lorenzo) this might not be needed once we have the bls aggregator
}

/* //TODO(lorenzo) for later
// used by tests to simulate unverified messages
func (b *base) unvalidate() {
	b.verified = false
	b.power = new(big.Int)
	b.sender = common.Address{}
}*/

type individualMsg struct {
	sender common.Address
	index  uint64 //TODO(lorenzo) check (index in the committee in header)
	base
}

// sender is populated at decoding time, however we cannot rely on it until signature verification
func (im *individualMsg) Sender() common.Address {
	im.RLock()
	defer im.RUnlock()
	if !im.verified {
		panic("unverified message")
	}
	return im.sender
}

// sender is populated at decoding time, however we cannot rely on it until signature verification
func (im *individualMsg) Index() uint64 {
	im.RLock()
	defer im.RUnlock()
	if !im.verified {
		panic("unverified message")
	}
	return im.index
}

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
	p.RLock()
	defer p.RUnlock()
	return fmt.Sprintf("{Round: %v, Height: %v, ValidRound: %v, ProposedBlockHash: %v}",
		p.round, p.H(), p.validRound, p.block.Hash().String())
}

func (p *Propose) Validate() error {
	//TODO(lorenzo) check
	p.Lock()
	defer p.Unlock()
	/*	if im.verified {
		return nil
	}*/

	valid := p.signature.Verify(p.senderKey, p.signatureHash[:])
	if !valid {
		return ErrBadSignature
	}

	return nil
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
	signatureInput := []any{ProposalCode, uint64(r), h, validRound, isValidRoundNil, block.Hash()}
	signatureInputEncoded, _ := rlp.EncodeToBytes(signatureInput)
	signatureHash := crypto.Hash(signatureInputEncoded)
	signature, validator := signer(signatureHash)

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
	// after a Validate() call during processing.
	return &Propose{
		block:      block,
		validRound: vr,
		individualMsg: individualMsg{
			sender: validator,
			base: base{
				round:          r,
				height:         h,
				signatureInput: signatureInput,
				signatureHash:  signatureHash,
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
	p.signatureInput = []any{ProposalCode, ext.Round, ext.Height, ext.ValidRound, ext.IsValidRoundNil, p.block.Hash()}
	p.payload = payload
	// precompute hash and signature hash
	p.hash = crypto.Hash(payload)
	sigData, _ := rlp.EncodeToBytes(p.signatureInput) //TODO(lorenzo) fine to ignore error?
	p.signatureHash = crypto.Hash(sigData)

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

func (p *LightProposal) String() string {
	p.RLock()
	defer p.RUnlock()
	return fmt.Sprintf("{Round: %v, Height: %v, sender: %v, power: %v, Code: %v, value: %v}", p.R(), p.H(), p.sender.String(), p.power, p.Code(), p.blockHash)
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
	p.signatureInput = []any{ProposalCode, ext.Round, ext.Height, ext.ValidRound, ext.IsValidRoundNil, p.blockHash}
	p.payload = payload
	// precompute hash and signature hash
	p.hash = crypto.Hash(payload)
	sigData, _ := rlp.EncodeToBytes(p.signatureInput)
	p.signatureHash = crypto.Hash(sigData)
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
	p.RLock()
	defer p.RUnlock()
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
	p.RLock()
	defer p.RUnlock()
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
	signatureHash := crypto.Hash(signatureEncodedInput)
	signature, validator := signer(signatureHash)
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
				signatureHash:  signatureHash,
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
	p.signatureInput = []any{PrevoteCode, encoded.Round, encoded.Height, encoded.Value}
	p.payload = payload
	// precompute hash and signature hash
	p.hash = crypto.Hash(payload)
	sigData, _ := rlp.EncodeToBytes(p.signatureInput)
	p.signatureHash = crypto.Hash(sigData)
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
	p.signatureInput = []any{PrecommitCode, encoded.Round, encoded.Height, encoded.Value}
	p.payload = payload
	// precompute hash and signature hash
	p.hash = crypto.Hash(payload)
	sigData, _ := rlp.EncodeToBytes(p.signatureInput)
	p.signatureHash = crypto.Hash(sigData)
	return nil
}

type extAggregateVote struct {
	Code      uint8
	Round     uint64
	Height    uint64
	Value     common.Hash
	Senders   Coefficients
	Signature *blst.BlsSignature
}

type aggregateMsg struct {
	senders   Coefficients
	addresses []common.Address //TODO(lorenzo) maybe this can go into the coefficients
	powers    []*big.Int
	base
}

// sender is populated at decoding time, however we cannot rely on it until signature verification
func (am *aggregateMsg) Senders() []common.Address {
	am.RLock()
	defer am.RUnlock()
	if !am.verified {
		panic("unverified message")
	}
	return am.addresses
}

// sender is populated at decoding time, however we cannot rely on it until signature verification
func (am *aggregateMsg) Powers() []*big.Int {
	am.RLock()
	defer am.RUnlock()
	if !am.verified {
		panic("unverified message")
	}
	return am.powers
}

func (am *aggregateMsg) SendersCoeff() Coefficients {
	//TODO(lorenzo) check
	am.RLock()
	defer am.RUnlock()
	if !am.verified {
		panic("unverified message")
	}
	return am.senders
}

type AggregatePrevote struct {
	value common.Hash
	aggregateMsg
}

func (ap *AggregatePrevote) Code() uint8 {
	return AggregatePrevoteCode
}

func (ap *AggregatePrevote) Value() common.Hash {
	return ap.value
}
func (ap *AggregatePrevote) String() string {
	ap.RLock()
	defer ap.RUnlock()
	return fmt.Sprintf("{r:  %v, h: %v , power: %v, Code: %v, value: %v}",
		ap.round, ap.height, ap.power, ap.Code(), ap.value)
}

func (ap *AggregatePrevote) DecodeRLP(s *rlp.Stream) error {
	payload, err := s.Raw()
	if err != nil {
		return err
	}

	encoded := &extAggregateVote{}
	if err := rlp.DecodeBytes(payload, encoded); err != nil {
		return err
	}
	if encoded.Code != AggregatePrevoteCode {
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
	//TODO(lorenzo) these checks enough?
	if encoded.Senders == nil || len(encoded.Senders) == 0 {
		return constants.ErrInvalidMessage
	}
	ap.height = encoded.Height
	ap.round = int64(encoded.Round)
	ap.value = encoded.Value
	ap.signature = encoded.Signature
	ap.senders = encoded.Senders
	// note: code for signature input is still prevote (this is correct)
	ap.signatureInput = []any{PrevoteCode, encoded.Round, encoded.Height, encoded.Value}
	ap.payload = payload
	// precompute hash and signature hash
	ap.hash = crypto.Hash(payload)
	sigData, _ := rlp.EncodeToBytes(ap.signatureInput)
	ap.signatureHash = crypto.Hash(sigData)
	return nil
}

type AggregatePrecommit struct {
	value common.Hash
	aggregateMsg
}

func (ap *AggregatePrecommit) Code() uint8 {
	return AggregatePrecommitCode
}

func (ap *AggregatePrecommit) Value() common.Hash {
	return ap.value
}
func (ap *AggregatePrecommit) String() string {
	ap.RLock()
	defer ap.RUnlock()
	return fmt.Sprintf("{r:  %v, h: %v , power: %v, Code: %v, value: %v}",
		ap.round, ap.height, ap.power, ap.Code(), ap.value)
}

func (ap *AggregatePrecommit) DecodeRLP(s *rlp.Stream) error {
	payload, err := s.Raw()
	if err != nil {
		return err
	}

	encoded := &extAggregateVote{}
	if err := rlp.DecodeBytes(payload, encoded); err != nil {
		return err
	}
	if encoded.Code != AggregatePrecommitCode {
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
	//TODO(lorenzo) these checks enough?
	if encoded.Senders == nil || len(encoded.Senders) == 0 {
		return constants.ErrInvalidMessage
	}
	ap.height = encoded.Height
	ap.round = int64(encoded.Round)
	ap.value = encoded.Value
	ap.signature = encoded.Signature
	ap.senders = encoded.Senders
	// note: code for signature input is still precommit (this is correct)
	ap.signatureInput = []any{PrecommitCode, encoded.Round, encoded.Height, encoded.Value}
	ap.payload = payload
	// precompute hash and signature hash
	ap.hash = crypto.Hash(payload)
	sigData, _ := rlp.EncodeToBytes(ap.signatureInput)
	ap.signatureHash = crypto.Hash(sigData)
	return nil
}

func NewAggregatePrevote(votes []Msg, header *types.Header) *AggregatePrevote {
	return NewAggregateVote[AggregatePrevote](votes, header)
}

func NewAggregatePrecommit(votes []Msg, header *types.Header) *AggregatePrecommit {
	return NewAggregateVote[AggregatePrecommit](votes, header)
}

func NewAggregateVote[
	E AggregatePrevote | AggregatePrecommit,
	PE interface {
		*E
		Msg
	}](votes []Msg, header *types.Header) *E {
	code := PE(new(E)).Code()

	//inCommittee := header.CommitteeMember
	committeeSize := len(header.Committee) //TODO(lorenzo) right? maybe max committeesize

	// compute coefficients and aggregated signature
	var signatures []blst.Signature
	c := NewCoefficients(uint64(committeeSize))

	// TODO(lorenzo) deal with duplicated msgs power here, right now it will count them twice
	// I need something like what I do in Prevalidate

	// we have to treat diffrently single votes and aggregates
	for _, vote := range votes {
		switch m := vote.(type) {
		case *Propose, *Prevote, *Precommit:
			c.Increment(m.(IndividualMsg).Index())
		case *AggregatePrevote, *AggregatePrecommit:
			c.Merge(m.(AggregateMsg).SendersCoeff())
		}
		signatures = append(signatures, vote.Signature())
	}
	aggregatedSignature := blst.Aggregate(signatures)

	// TODO(lorenzo) this should deal with duplicated votes in the aggregate
	aggregatedPower := new(big.Int)
	for _, index := range c.FlattenUniq() {
		aggregatedPower.Add(aggregatedPower, header.Committee[index].VotingPower) //TODO(lorenzo) deal with indexes out of range
	}

	// assumes all votes are for the same signature input (code,h,r,value)
	vote := votes[0]
	h := vote.H()
	r := vote.R()
	value := vote.Value()

	signatureInput := []any{vote.Code(), uint64(r), h, value}

	payload, _ := rlp.EncodeToBytes(extAggregateVote{
		Code:      code,
		Round:     uint64(r),
		Height:    h,
		Value:     value,
		Senders:   c,
		Signature: aggregatedSignature.(*blst.BlsSignature),
	})

	//TODO(lorenzo) check power and verified
	aggregateVote := E{
		value: value,
		aggregateMsg: aggregateMsg{
			senders: c,
			base: base{
				round:          r,
				height:         h,
				payload:        payload,
				hash:           crypto.Hash(payload),
				signatureInput: signatureInput,
				//verified:       true,
				verified:  false,
				power:     aggregatedPower,
				signature: aggregatedSignature,
			},
		},
	}
	return &aggregateVote
}

// TODO(lorenzo) always returns nil
func (am *aggregateMsg) PreValidate(header *types.Header) error {
	//TODO(lorenzo) block needed?
	am.Lock()
	defer am.Unlock()
	if am.verified {
		return nil
	}

	// aggregate public keys
	indexes := am.senders.Flatten()
	var keys [][]byte //TODO(lorenzo) pre-allocate?
	aggregatedPower := new(big.Int)
	var addresses []common.Address
	var powers []*big.Int
	first := true
	var previousIndex uint64

	for _, index := range indexes {
		member := header.Committee[index]
		keys = append(keys, member.ConsensusKeyBytes)
		//ugly, I think I could use flatten uniq
		if first || index != previousIndex {
			aggregatedPower.Add(aggregatedPower, member.VotingPower)
			addresses = append(addresses, member.Address)
			powers = append(powers, member.VotingPower)
		}
		first = false
		previousIndex = index
	}

	aggregatedKey, err := blst.AggregatePublicKeys(keys)
	if err != nil {
		panic(err) //TODO(lorenzo) fix
	}

	am.power = aggregatedPower
	am.addresses = addresses //TODO(lorenzo) double check this, is this fine?
	am.powers = powers
	am.senderKey = aggregatedKey
	am.verified = true
	return nil
}

func (am *aggregateMsg) Validate() error {
	//TODO(lorenzo) check
	am.Lock()
	defer am.Unlock()
	/*	if im.verified {
		return nil
	}*/

	valid := am.signature.Verify(am.senderKey, am.signatureHash[:])
	if !valid {
		return ErrBadSignature
	}

	return nil
}

/*
// Validate verify the signature and sets the power field
func (am *aggregateMsg) Validate(header *types.Header) error {
	am.Lock()
	defer am.Unlock()
	if am.verified {
		return nil
	}

	// aggregate public keys
	indexes := am.senders.Flatten()
	aggregatedKey := new(blst.BlsPublicKey)
	aggregatedPower := new(big.Int)
	var addresses []common.Address
	first := true
	var previousIndex uint64

	for _, index := range indexes {
		member := header.Committee[index]
		aggregatedKey.Aggregate(member.ConsensusKey) //TODO(lorenzo) deal with index out of range
		//ugly
		if first || index != previousIndex {
			aggregatedPower.Add(aggregatedPower, member.VotingPower)
			addresses = append(addresses, member.Address)
		}
		first = false
		previousIndex = index
	}

	// We are not saving the rlp encoded signature input data as we want
	// to avoid this extra-serialization step if the message has already been received
	// The call to Validate() only happen after the cache check in the backend handler.
	sigData, _ := rlp.EncodeToBytes(am.signatureInput)
	hash := crypto.Hash(sigData)

	valid := am.signature.Verify(aggregatedKey, hash[:])
	if !valid {
		return ErrBadSignature
	}

	am.power = aggregatedPower
	am.addresses = addresses //TODO(lorenzo) double check this, is this fine?
	am.verified = true
	return nil
}*/

/*
// sender is populated at decoding time, however we cannot rely on it until signature verification
func (b *base) Sender() common.Address {
	b.RLock()
	defer b.RUnlock()
	if !b.verified {
		panic("unverified message")
	}
	return b.sender
}*/

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

func (b *base) Signature() blst.Signature {
	return b.signature
}

func (b *base) Payload() []byte {
	return b.payload
}

func (b *base) Hash() common.Hash {
	return b.hash
}

func (b *base) SignatureHash() common.Hash {
	return b.signatureHash
}

// TODO(lorenzo) name not really appropriate since key could be an aggregate
func (b *base) SenderKey() blst.PublicKey {
	b.RLock()
	defer b.RUnlock()
	if !b.verified {
		panic("unverified message")
	}
	return b.senderKey
}

func (im *individualMsg) PreValidate(header *types.Header) error {
	//TODO(lorenzo) do I need this block, or do I need to adapt
	im.Lock()
	defer im.Unlock()
	if im.verified {
		return nil
	}

	validator := header.CommitteeMember(im.sender)
	if validator == nil {
		return ErrUnauthorizedAddress
	}

	im.senderKey = validator.ConsensusKey
	im.index = validator.Index
	im.power = validator.VotingPower
	im.verified = true //TODO(lorenzo) not really but we can access power and consensus key
	return nil

}

/*
// Validate verify the signature and sets the power field
func (im *individualMsg) Validate(header *types.Header) error {
	im.Lock()
	defer im.Unlock()
	if im.verified {
		return nil
	}

	valid := im.signature.Verify(im.senderKey, hash[:])
	if !valid {
		return ErrBadSignature
	}

	im.power = validator.VotingPower
	im.verified = true
	return nil
}*/

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
