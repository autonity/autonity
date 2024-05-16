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
	"sort"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/rlp"
)

var (
	ErrBadSignature            = errors.New("bad signature")
	ErrUnauthorizedAddress     = errors.New("unauthorized address")
	ErrInvalidComplexAggregate = errors.New("complex aggregate does not carry quorum")
)

const (
	ProposalCode uint8 = iota
	PrevoteCode
	PrecommitCode
	LightProposalCode
)

type Signer func(hash common.Hash) blst.Signature

// TODO(lorenzo) refinements, maybe we can just send the sender index instead of the sender address
type Propose struct {
	block      *types.Block
	validRound int64
	// node address of the sender, populated at decoding phase
	sender common.Address
	// populated at PreValidate phase
	senderIndex int      // index of the sender in the committee
	power       *big.Int // power of sender
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

func (p *Propose) Power() *big.Int {
	if !p.preverified {
		panic("Trying to access power on not preverified proposal")
	}
	return p.power
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

func (p *Propose) ToLight() *LightProposal {
	return NewLightProposal(p)
}

func NewPropose(r int64, h uint64, vr int64, block *types.Block, signer Signer, self *types.CommitteeMember) *Propose {
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
	signature := signer(signatureInput)

	validator := self.Address

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

	return &Propose{
		block:       block,
		validRound:  vr,
		sender:      validator,
		senderIndex: int(self.Index),
		power:       new(big.Int).Set(self.VotingPower),
		base: base{
			height:         h,
			round:          r,
			signatureInput: signatureInput,
			signature:      signature,
			payload:        payload,
			hash:           crypto.Hash(payload),
			verified:       true,
			preverified:    true,
			senderKey:      self.ConsensusKey,
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
	p.verified = false
	p.preverified = false
	return nil
}

func (p *Propose) Sender() common.Address {
	return p.sender
}

func (p *Propose) SenderIndex() int {
	if !p.preverified {
		panic("Trying to access sender index on not preverified proposal")
	}
	return p.senderIndex
}

func (p *Propose) PreValidate(header *types.Header) error {
	validator := header.CommitteeMember(p.sender)
	if validator == nil {
		return ErrUnauthorizedAddress
	}

	p.senderKey = validator.ConsensusKey
	p.senderIndex = int(validator.Index)
	p.power = validator.VotingPower //TODO(lorenzo) do I need a copy here? same for lightproposal
	p.preverified = true
	return nil
}

type LightProposal struct {
	blockHash  common.Hash
	validRound int64
	// node address of the sender, populated at decoding phase
	sender common.Address
	// populated at PreValidate phase
	senderIndex int      // index of the sender in the committee
	power       *big.Int // power of sender
	base
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

func (p *LightProposal) Power() *big.Int {
	return p.power
}

// TODO(lorenzo) refinements, would be useful to print also sender and power, but we need to make sure they are trsuted (verified)
// same goes for the other message types
func (p *LightProposal) String() string {
	return fmt.Sprintf("{%s, Code: %v, value: %v}", p.base.String(), p.Code(), p.blockHash)
}

func NewLightProposal(proposal *Propose) *LightProposal {
	if !proposal.verified || !proposal.preverified {
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
		blockHash:   proposal.Block().Hash(),
		validRound:  proposal.validRound,
		sender:      proposal.sender,
		senderIndex: proposal.senderIndex,
		power:       proposal.power,
		base: base{
			round:          proposal.round,
			height:         proposal.height,
			signature:      proposal.signature,
			signatureInput: proposal.signatureInput,
			payload:        payload,
			hash:           crypto.Hash(payload),
			verified:       true,
			preverified:    true,
			senderKey:      proposal.senderKey,
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
	p.verified = false
	p.preverified = false
	return nil
}

func (p *LightProposal) Sender() common.Address {
	return p.sender
}

func (p *LightProposal) SenderIndex() int {
	return p.senderIndex
}

func (p *LightProposal) PreValidate(header *types.Header) error {
	validator := header.CommitteeMember(p.sender)
	if validator == nil {
		return ErrUnauthorizedAddress
	}

	p.senderKey = validator.ConsensusKey
	p.senderIndex = int(validator.Index)
	p.power = validator.VotingPower
	p.preverified = true
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
	Senders   *types.SendersInfo
	Signature *blst.BlsSignature
}

// TODO: would be good to do the same thing for proposal and lightproposal (to avoid code repetition)
type vote struct {
	senders *types.SendersInfo
	base
}

func (v *vote) Senders() *types.SendersInfo {
	return v.senders
}

func (v *vote) Power() *big.Int {
	return v.senders.Power()
}

func (v *vote) PreValidate(header *types.Header) error {
	if err := v.senders.Valid(len(header.Committee)); err != nil {
		return fmt.Errorf("Invalid senders information: %w", err)
	}

	// compute aggregated key and auxiliary data structures
	indexes := v.senders.Flatten()
	keys := make([][]byte, len(indexes))
	powers := make(map[int]*big.Int)
	power := new(big.Int)

	for i, index := range indexes {
		member := header.Committee[index]

		keys[i] = member.ConsensusKeyBytes

		_, alreadyPresent := powers[index]
		if !alreadyPresent {
			powers[index] = new(big.Int).Set(member.VotingPower)
			power.Add(power, member.VotingPower)
		}
	}
	aggregatedKey, err := blst.AggregatePublicKeys(keys)
	if err != nil {
		panic("Error while aggregating keys from committee: " + err.Error())
	}

	v.senders.AssignPower(powers, power)
	v.senderKey = aggregatedKey
	v.preverified = true

	// if the aggregate is a complex aggregate, it needs to carry quorum
	if v.senders.IsComplex() && v.Power().Cmp(bft.Quorum(header.TotalVotingPower())) < 0 {
		return ErrInvalidComplexAggregate
	}
	return nil
}

type Prevote struct {
	value common.Hash
	vote
}

func (p *Prevote) Code() uint8 {
	return PrevoteCode
}

func (p *Prevote) Value() common.Hash {
	return p.value
}

func (p *Prevote) String() string {
	return fmt.Sprintf("{r:  %v, h: %v , Code: %v, value: %v}",
		p.round, p.height, p.Code(), p.value)
}

type Precommit struct {
	value common.Hash
	vote
}

func (p *Precommit) Code() uint8 {
	return PrecommitCode
}

func (p *Precommit) Value() common.Hash {
	return p.value
}

func (p *Precommit) String() string {
	return fmt.Sprintf("{r:  %v, h: %v , Code: %v, value: %v}",
		p.round, p.height, p.Code(), p.value)
}

func newVote[
	E Prevote | Precommit,
	PE interface {
		*E
		Msg
	}](r int64, h uint64, value common.Hash, signer Signer, self *types.CommitteeMember, csize int) *E {
	code := PE(new(E)).Code()

	// Pay attention that we're adding the message Code to the signature input data.
	signaturePayload, _ := rlp.EncodeToBytes([]any{code, uint64(r), h, value})
	signatureInput := crypto.Hash(signaturePayload)
	signature := signer(signatureInput)

	// TODO(lorenzo) refinements, aggregates for different heights might have different len(senders.bits). Is that a problem?
	senders := types.NewSendersInfo(csize)
	senders.Increment(self)

	payload, _ := rlp.EncodeToBytes(extVote{
		Code:      code,
		Round:     uint64(r),
		Height:    h,
		Value:     value,
		Senders:   senders,
		Signature: signature.(*blst.BlsSignature),
	})
	vote := E{
		value: value,
		vote: vote{
			senders: senders,
			base: base{
				round:          r,
				height:         h,
				signature:      signature,
				payload:        payload,
				hash:           crypto.Hash(payload),
				signatureInput: signatureInput,
				verified:       true,
				preverified:    true,
				senderKey:      self.ConsensusKey,
			},
		},
	}
	return &vote
}

func NewPrevote(r int64, h uint64, value common.Hash, signer Signer, self *types.CommitteeMember, csize int) *Prevote {
	return newVote[Prevote](r, h, value, signer, self, csize)
}

func NewPrecommit(r int64, h uint64, value common.Hash, signer Signer, self *types.CommitteeMember, csize int) *Precommit {
	return newVote[Precommit](r, h, value, signer, self, csize)
}

// NOTE: these functions allow for the creation of complex aggregates
func AggregatePrevotes(votes []Vote) *Prevote {
	return AggregateVotes[Prevote](votes)
}

func AggregatePrecommits(votes []Vote) *Precommit {
	return AggregateVotes[Precommit](votes)
}

//TODO(lorenzo) refinements2, Instead of creating a new message object I can re-use one of the votes being aggregated

// NOTE: this function assumes that:
// 1. all votes are for the same signature input (code,h,r,value)
// 2. all votes have previously been preverified and cryptographically verified
func AggregateVotes[
	E Prevote | Precommit,
	PE interface {
		*E
		Msg
	}](votes []Vote) *E {
	// TODO(lorenzo) what if len(votes) == 0? is it possible?
	code := PE(new(E)).Code()

	// use votes[0] as a set representative
	representative := votes[0]

	// TODO(lorenzo) refinements, aggregates for different heights might have different len(senders.bits). Is that a problem?
	senders := types.NewSendersInfo(representative.Senders().CommitteeSize())

	// order votes by decreasing number of distinct signers.
	// This ensures that we reduce as much as possible the number of duplicated signatures for the same validator
	sort.Slice(votes, func(i, j int) bool {
		return votes[i].Senders().Len() > votes[j].Senders().Len()
	})

	// compute new aggregated signature and related sender information
	var signatures []blst.Signature
	var publicKeys [][]byte
	for _, vote := range votes {
		// do not aggregate votes if they do not add any useful information
		// e.g. senders contains already at least 1 signature for all signers of vote.Senders()
		// we would just create and gossip new aggregates that would uselessly flood the network
		// additionally, we also check if the resulting aggregate respects the coefficient boundaries.
		// this avoids that we aggregate two complex aggregates together, which can lead to coefficient breaching.
		if senders.AddsInformation(vote.Senders()) && senders.RespectsBoundaries(vote.Senders()) {
			senders.Merge(vote.Senders())
			signatures = append(signatures, vote.Signature())
			publicKeys = append(publicKeys, vote.SenderKey().Marshal())
		}
		//TODO(lorenzo) here we are silently dropping votes. This is not good + we should probably do it before signature verification
	}
	aggregatedSignature := blst.Aggregate(signatures)
	aggregatedPublicKey, err := blst.AggregatePublicKeys(publicKeys)
	if err != nil {
		panic("Cannot generate aggregate public key from valid votes: " + err.Error())
	}

	h := representative.H()
	r := representative.R()
	value := representative.Value()
	signatureInput := representative.SignatureInput()

	payload, _ := rlp.EncodeToBytes(extVote{
		Code:      code,
		Round:     uint64(r),
		Height:    h,
		Value:     value,
		Senders:   senders,
		Signature: aggregatedSignature.(*blst.BlsSignature),
	})

	aggregateVote := E{
		value: value,
		vote: vote{
			senders: senders,
			base: base{
				height:         h,
				round:          r,
				signatureInput: signatureInput,
				signature:      aggregatedSignature,
				payload:        payload,
				hash:           crypto.Hash(payload),
				verified:       true, // verified due to all votes being verified
				preverified:    true,
				senderKey:      aggregatedPublicKey, // this is not strictly necessary since the vote is already verified
			},
		},
	}
	return &aggregateVote
}

// NOTE: these functions will aggregate votes as much as possible without creating complex aggregates
func AggregatePrevotesSimple(votes []Vote) []*Prevote {
	return AggregateVotesSimple[Prevote](votes)
}

func AggregatePrecommitsSimple(votes []Vote) []*Precommit {
	return AggregateVotesSimple[Precommit](votes)
}

// NOTE: this function assumes that:
// 1. all votes are for the same signature input (code,h,r,value)
// 2. all votes have previously been cryptographically verified
func AggregateVotesSimple[
	E Prevote | Precommit,
	PE interface {
		*E
		Msg
	}](votes []Vote) []*E {
	// TODO(lorenzo) what if len(votes) == 0? is it possible?
	code := PE(new(E)).Code()

	skip := make([]bool, len(votes))
	var sendersList []*types.SendersInfo
	var signaturesList [][]blst.Signature
	var publicKeysList [][][]byte

	// order votes by decreasing number of distinct signers.
	// This ensures that we reduce as much as possible the number of duplicated signatures for the same validator
	sort.Slice(votes, func(i, j int) bool {
		return votes[i].Senders().Len() > votes[j].Senders().Len()
	})

	//TODO(lorenzo) in the following we are silently dropping votes. This is not good + we should probably do it before signature verification

	// TODO(lorenzo) this is not the most optimized version I believe
	// at least add an early termination
	for i, vote := range votes {
		if skip[i] {
			continue
		}
		senders := vote.Senders().Copy()
		signatures := []blst.Signature{vote.Signature()}
		publicKeys := [][]byte{vote.SenderKey().Marshal()}
		for j := i + 1; j < len(votes); j++ {
			if skip[j] {
				continue
			}
			other := votes[j]
			if !senders.AddsInformation(other.Senders()) {
				skip[j] = true // TODO(lorenzo) consider keeping it, it might aggregate with other votes
				continue
			}
			if !senders.CanMergeSimple(other.Senders()) {
				continue
			}
			senders.Merge(other.Senders())
			signatures = append(signatures, other.Signature())
			publicKeys = append(publicKeys, other.SenderKey().Marshal())
			skip[j] = true
		}
		sendersList = append(sendersList, senders)
		signaturesList = append(signaturesList, signatures)
		publicKeysList = append(publicKeysList, publicKeys)
	}

	// build aggregates
	representative := votes[0]
	h := representative.H()
	r := representative.R()
	value := representative.Value()
	signatureInput := representative.SignatureInput()

	n := len(sendersList)
	aggregateVotes := make([]*E, n)
	for i := 0; i < n; i++ {
		var aggregatedSignature blst.Signature
		var aggregatedPublicKey blst.PublicKey
		var err error
		if len(signaturesList[i]) == 1 {
			aggregatedSignature = signaturesList[i][0]
			aggregatedPublicKey, _ = blst.PublicKeyFromBytes(publicKeysList[i][0])
		} else {
			aggregatedSignature = blst.Aggregate(signaturesList[i])
			aggregatedPublicKey, err = blst.AggregatePublicKeys(publicKeysList[i])
			if err != nil {
				panic("Cannot generate aggregate public key from valid votes: " + err.Error())
			}
		}

		payload, _ := rlp.EncodeToBytes(extVote{
			Code:      code,
			Round:     uint64(r),
			Height:    h,
			Value:     value,
			Senders:   sendersList[i],
			Signature: aggregatedSignature.(*blst.BlsSignature),
		})

		aggregateVote := E{
			value: value,
			vote: vote{
				senders: sendersList[i],
				base: base{
					height:         h,
					round:          r,
					signatureInput: signatureInput,
					signature:      aggregatedSignature,
					payload:        payload,
					hash:           crypto.Hash(payload),
					verified:       true, // verified due to all votes being verified
					preverified:    true,
					senderKey:      aggregatedPublicKey, // this is not strictly necessary since the vote is already verified
				},
			},
		}
		aggregateVotes[i] = &aggregateVote
	}
	return aggregateVotes
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
	if encoded.Senders == nil || encoded.Senders.Bits == nil || len(encoded.Senders.Bits) == 0 {
		return constants.ErrInvalidMessage
	}
	p.height = encoded.Height
	p.round = int64(encoded.Round)
	p.value = encoded.Value
	p.signature = encoded.Signature
	p.senders = encoded.Senders
	p.payload = payload
	// precompute hash and signature hash
	signaturePayload, _ := rlp.EncodeToBytes([]any{PrevoteCode, encoded.Round, encoded.Height, encoded.Value})
	p.signatureInput = crypto.Hash(signaturePayload)
	p.hash = crypto.Hash(payload)
	p.verified = false
	p.preverified = false
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
	if encoded.Senders == nil || encoded.Senders.Bits == nil || len(encoded.Senders.Bits) == 0 {
		return constants.ErrInvalidMessage
	}
	p.height = encoded.Height
	p.round = int64(encoded.Round)
	p.value = encoded.Value
	p.signature = encoded.Signature
	p.senders = encoded.Senders
	p.payload = payload
	// precompute hash and signature hash
	signaturePayload, _ := rlp.EncodeToBytes([]any{PrecommitCode, encoded.Round, encoded.Height, encoded.Value})
	p.signatureInput = crypto.Hash(signaturePayload)
	p.hash = crypto.Hash(payload)
	p.verified = false
	p.preverified = false
	return nil
}

// PrepareCommittedSeal returns the input data to compute the committed seal for a given block hash.
func PrepareCommittedSeal(hash common.Hash, round int64, height *big.Int) common.Hash {
	// this is matching the signature input that we get from the committed messages.
	buf, _ := rlp.EncodeToBytes([]any{PrecommitCode, uint64(round), height.Uint64(), hash})
	return crypto.Hash(buf)
}

// TODO(lorenzo) refinements, update fake msg

/*
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
}*/
