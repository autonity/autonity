// Package message implements an interface and the three underlying consensus messages types that
// tendermint is using: Propose, Prevote and Precommit.
// In addition to that, we have a special type, the "Light Proposal" which is being used for
// accountability purposes. Light proposals are never directly brodcasted
// over the network but always part of a proof object, defined in the accountability package.
// Messages can exist in two states: unverified and verified depending on their signature verification.
// When verified, calling `Validate` the voting power information becomes available, and the signer can be relied upon.
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

// TODO: To save space we could send only the signer index instead of the signer address
type Propose struct {
	block      *types.Block
	validRound int64
	// node address of the signer, populated at decoding phase
	signer common.Address
	// populated at PreValidate phase
	signerIndex int      // index of the signer in the committee
	power       *big.Int // power of signer
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
	// since we do not have ecrecover with BLS signatures, we need to also send the signer in the message.
	// It is sent not-signed to facilitate aggregation.
	// If tampered with, the signature will fail anyways because we will fetch the wrong key.
	Signer    common.Address
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
		Signer:          validator,
		Signature:       signature.(*blst.BlsSignature),
	})

	return &Propose{
		block:       block,
		validRound:  vr,
		signer:      validator,
		signerIndex: int(self.Index),
		power:       self.VotingPower,
		base: base{
			height:         h,
			round:          r,
			signatureInput: signatureInput,
			signature:      signature,
			payload:        payload,
			hash:           crypto.Hash(payload),
			verified:       true,
			preverified:    true,
			signerKey:      self.ConsensusKey,
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
	p.signer = ext.Signer
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

func (p *Propose) Signer() common.Address {
	return p.signer
}

func (p *Propose) SignerIndex() int {
	if !p.preverified {
		panic("Trying to access signer index on not preverified proposal")
	}
	return p.signerIndex
}

func (p *Propose) PreValidate(header *types.Header) error {
	validator := header.CommitteeMember(p.signer)
	if validator == nil {
		return ErrUnauthorizedAddress
	}

	p.signerKey = validator.ConsensusKey
	p.signerIndex = int(validator.Index)
	p.power = validator.VotingPower
	p.preverified = true
	return nil
}

type LightProposal struct {
	blockHash  common.Hash
	validRound int64
	// node address of the signer, populated at decoding phase
	signer common.Address
	// populated at PreValidate phase
	signerIndex int      // index of the signer in the committee
	power       *big.Int // power of signer
	base
}

type extLightProposal struct {
	Code            uint8
	Round           uint64
	Height          uint64
	ValidRound      uint64
	IsValidRoundNil bool
	ProposalBlock   common.Hash
	Signer          common.Address
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

// TODO(lorenzo) refinements, would be useful to print also signer and power, but we need to make sure they are trsuted (verified)
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
		Signer:          proposal.signer,
		Signature:       proposal.signature.(*blst.BlsSignature),
	})
	return &LightProposal{
		blockHash:   proposal.Block().Hash(),
		validRound:  proposal.validRound,
		signer:      proposal.signer,
		signerIndex: proposal.signerIndex,
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
			signerKey:      proposal.signerKey,
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
	p.signer = ext.Signer
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

func (p *LightProposal) Signer() common.Address {
	return p.signer
}

func (p *LightProposal) SignerIndex() int {
	return p.signerIndex
}

func (p *LightProposal) PreValidate(header *types.Header) error {
	validator := header.CommitteeMember(p.signer)
	if validator == nil {
		return ErrUnauthorizedAddress
	}

	p.signerKey = validator.ConsensusKey
	p.signerIndex = int(validator.Index)
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
	Signers   *types.Signers
	Signature *blst.BlsSignature
}

// TODO: would be good to do the same thing for proposal and lightproposal (to avoid code repetition)
type vote struct {
	signers *types.Signers
	base
}

func (v *vote) Signers() *types.Signers {
	return v.signers
}

func (v *vote) Power() *big.Int {
	return v.signers.Power()
}

func (v *vote) PreValidate(header *types.Header) error {
	if err := v.signers.Validate(len(header.Committee)); err != nil {
		return fmt.Errorf("Invalid signers information: %w", err)
	}

	// compute aggregated key and auxiliary data structures
	indexes := v.signers.Flatten()
	keys := make([][]byte, len(indexes))
	powers := make(map[int]*big.Int)
	power := new(big.Int)

	for i, index := range indexes {
		member := header.Committee[index]

		keys[i] = member.ConsensusKeyBytes

		_, alreadyPresent := powers[index]
		if !alreadyPresent {
			powers[index] = member.VotingPower
			power.Add(power, member.VotingPower)
		}
	}
	aggregatedKey, err := blst.AggregatePublicKeys(keys)
	if err != nil {
		panic("Error while aggregating keys from committee: " + err.Error())
	}

	v.signers.AssignPower(powers, power)
	v.signerKey = aggregatedKey
	v.preverified = true

	// if the aggregate is a complex aggregate, it needs to carry quorum
	if v.signers.IsComplex() && v.Power().Cmp(bft.Quorum(header.TotalVotingPower())) < 0 {
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

	signers := types.NewSigners(csize)
	signers.Increment(self)

	payload, _ := rlp.EncodeToBytes(extVote{
		Code:      code,
		Round:     uint64(r),
		Height:    h,
		Value:     value,
		Signers:   signers,
		Signature: signature.(*blst.BlsSignature),
	})
	vote := E{
		value: value,
		vote: vote{
			signers: signers,
			base: base{
				round:          r,
				height:         h,
				signature:      signature,
				payload:        payload,
				hash:           crypto.Hash(payload),
				signatureInput: signatureInput,
				verified:       true,
				preverified:    true,
				signerKey:      self.ConsensusKey,
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

// NOTE: this function assumes that:
// 1. all votes are for the same signature input (code,h,r,value)
// 2. all votes have previously been preverified and cryptographically verified
func AggregateVotes[
	E Prevote | Precommit,
	PE interface {
		*E
		Msg
	}](votes []Vote) *E {
	// length safety checks
	if len(votes) == 0 {
		panic("Trying to aggregate empty set of votes")
	}

	// use votes[0] as a set representative
	representative := votes[0]

	// signers of the aggregate
	signers := types.NewSigners(representative.Signers().CommitteeSize())

	// order votes by decreasing number of distinct signers.
	// This ensures that we reduce as much as possible the number of duplicated signatures for the same validator
	sort.Slice(votes, func(i, j int) bool {
		return votes[i].Signers().Len() > votes[j].Signers().Len()
	})

	// compute new aggregated signature and related signers information
	var signatures []blst.Signature
	var publicKeys [][]byte
	for _, vote := range votes {
		// do not aggregate votes if they do not add any useful information
		// e.g. signers contains already at least 1 signature for all signers of vote.Signers()
		// we would just create and gossip new aggregates that would uselessly flood the network
		// additionally, we also check if the resulting aggregate respects the coefficient boundaries.
		// this avoids that we aggregate two complex aggregates together, which can lead to coefficient breaching.
		if signers.AddsInformation(vote.Signers()) && signers.RespectsBoundaries(vote.Signers()) {
			signers.Merge(vote.Signers())
			signatures = append(signatures, vote.Signature())
			publicKeys = append(publicKeys, vote.SignerKey().Marshal())
		}
	}
	aggregatedSignature := blst.Aggregate(signatures)
	aggregatedPublicKey, err := blst.AggregatePublicKeys(publicKeys)
	if err != nil {
		panic("Cannot generate aggregate public key from valid votes: " + err.Error())
	}

	c := representative.Code()
	h := representative.H()
	r := representative.R()
	value := representative.Value()
	signatureInput := representative.SignatureInput()

	payload, _ := rlp.EncodeToBytes(extVote{
		Code:      c,
		Round:     uint64(r),
		Height:    h,
		Value:     value,
		Signers:   signers,
		Signature: aggregatedSignature.(*blst.BlsSignature),
	})

	aggregateVote := E{
		value: value,
		vote: vote{
			signers: signers,
			base: base{
				height:         h,
				round:          r,
				signatureInput: signatureInput,
				signature:      aggregatedSignature,
				payload:        payload,
				hash:           crypto.Hash(payload),
				verified:       true, // verified due to all votes being verified
				preverified:    true,
				signerKey:      aggregatedPublicKey, // this is not strictly necessary since the vote is already verified
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
	// length safety checks
	if len(votes) == 0 {
		panic("Trying to aggregate empty set of votes")
	}
	code := PE(new(E)).Code()

	skip := make([]bool, len(votes))
	var signersList []*types.Signers
	var signaturesList [][]blst.Signature
	var publicKeysList [][][]byte

	// order votes by decreasing number of distinct signers.
	// This ensures that we reduce as much as possible the number of duplicated signatures for the same validator
	sort.Slice(votes, func(i, j int) bool {
		return votes[i].Signers().Len() > votes[j].Signers().Len()
	})

	// TODO(lorenzo) this is not the most optimized version I believe. At least add an early termination
	for i, vote := range votes {
		if skip[i] {
			continue
		}
		signers := vote.Signers().Copy()
		signatures := []blst.Signature{vote.Signature()}
		publicKeys := [][]byte{vote.SignerKey().Marshal()}
		for j := i + 1; j < len(votes); j++ {
			if skip[j] {
				continue
			}
			other := votes[j]
			if !signers.AddsInformation(other.Signers()) {
				skip[j] = true // TODO(lorenzo) consider keeping it, it might aggregate with other votes
				continue
			}
			if !signers.CanMergeSimple(other.Signers()) {
				continue
			}
			signers.Merge(other.Signers())
			signatures = append(signatures, other.Signature())
			publicKeys = append(publicKeys, other.SignerKey().Marshal())
			skip[j] = true
		}
		signersList = append(signersList, signers)
		signaturesList = append(signaturesList, signatures)
		publicKeysList = append(publicKeysList, publicKeys)
	}

	// build aggregates
	representative := votes[0]
	h := representative.H()
	r := representative.R()
	value := representative.Value()
	signatureInput := representative.SignatureInput()

	n := len(signersList)
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
			Signers:   signersList[i],
			Signature: aggregatedSignature.(*blst.BlsSignature),
		})

		aggregateVote := E{
			value: value,
			vote: vote{
				signers: signersList[i],
				base: base{
					height:         h,
					round:          r,
					signatureInput: signatureInput,
					signature:      aggregatedSignature,
					payload:        payload,
					hash:           crypto.Hash(payload),
					verified:       true, // verified due to all votes being verified
					preverified:    true,
					signerKey:      aggregatedPublicKey, // this is not strictly necessary since the vote is already verified
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
	if encoded.Signers == nil || encoded.Signers.Bits == nil || len(encoded.Signers.Bits) == 0 {
		return constants.ErrInvalidMessage
	}
	p.height = encoded.Height
	p.round = int64(encoded.Round)
	p.value = encoded.Value
	p.signature = encoded.Signature
	p.signers = encoded.Signers
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
	if encoded.Signers == nil || encoded.Signers.Bits == nil || len(encoded.Signers.Bits) == 0 {
		return constants.ErrInvalidMessage
	}
	p.height = encoded.Height
	p.round = int64(encoded.Round)
	p.value = encoded.Value
	p.signature = encoded.Signature
	p.signers = encoded.Signers
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

// computes the power of a set of messages. Every sender's power is counted only once
func Power(messages []Msg) *big.Int {
	powerInfo := NewPowerInfo()

	for _, msg := range messages {
		switch m := msg.(type) {
		case *Propose:
			powerInfo.Set(m.SignerIndex(), m.Power())
		case *Prevote, *Precommit:
			vote := m.(Vote)
			for index, power := range vote.Signers().Powers() {
				powerInfo.Set(index, power)
			}
		default:
			panic("unknown message type")
		}
	}
	return powerInfo.Pow()
}

// TODO(lorenzo) refinements, update fake msg

/*
// Fake is a dummy object used for internal testing.
type Fake struct {
	FakeCode      uint8
	FakeRound     uint64
	FakeHeight    uint64
	FakeValue     common.Hash
	FakePayload   []byte
	FakeHash      common.Hash
	FakeSender    common.Address
	FakeSignature blst.Signature
	FakePower     *big.Int
}

func (f Fake) R() int64                                                       { return int64(f.FakeRound) }
func (f Fake) H() uint64                                                      { return f.FakeHeight }
func (f Fake) Code() uint8                                                    { return f.FakeCode }
func (f Fake) Signer() common.Address                                         { return f.FakeSender }
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
		base: base{
			round:     int64(f.FakeRound),
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
			round:     int64(f.FakeRound),
			height:    f.FakeHeight,
			signature: f.FakeSignature,
			payload:   f.FakePayload,
			power:     f.FakePower,
			sender:    f.FakeSender,
			hash:      f.FakeHash,
			verified:  true,
		},
	}
}*/
