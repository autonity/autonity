// Package message implements an interface and the three underlying consensus messages types that
// tendermint is using: Propose, Prevote and Precommit.
// In addition to that, we have a special type, the "Light Proposal" which is being used for
// accountability purposes. Light proposals are never directly brodcasted
// over the network but always part of a proof object, defined in the accountability package.
// There are three ways that a consensus message can be instantiated:
//   - using a "New" constructor, e.g. NewPrevote :
//     The object is fully created, with signature and final payload already pre-computed. Ready for use.
//   - decoding a RLP-encoded message from the wire. Needs to pass a two step verification process.
//   - using a Fake constructor. Used in tests.
//
// Messages received from the wire needs to pass two verification steps before they can be trusted:
// - Preverification, which attaches some auxiliary data to the message, such as bls keys and power information.
// - Verification, which validates the actual BLS signature.
//
// The two flags `preverified` and `verified` provide information on the status of the message.
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
	"github.com/autonity/autonity/metrics"
	"github.com/autonity/autonity/rlp"
)

var (
	ErrBadSignature            = errors.New("bad signature")
	ErrUnauthorizedAddress     = errors.New("unauthorized address")
	ErrInvalidComplexAggregate = errors.New("complex aggregate does not carry quorum")
	ErrInvalidIndividualVote   = errors.New("individual vote has 0 signature")
)

// Internal message codes used by tendermint consensus engine and its accountability module.
const (
	ProposalCode uint8 = iota
	PrevoteCode
	PrecommitCode
	LightProposalCode
	EvidenceVoteCode
)

// Message IDs used by the ACN p2p network layer to deliver raw messages of the upper layer.
const (
	ProposeNetworkMsg        uint64 = 0x11
	PrevoteNetworkMsg        uint64 = 0x12
	PrecommitNetworkMsg      uint64 = 0x13
	SyncNetworkMsg           uint64 = 0x14
	AccountabilityNetworkMsg uint64 = 0x15
)

// NetworkCodes maps msg code to its raw msg ID for msg relaying in ACN network.
var NetworkCodes = map[uint8]uint64{
	ProposalCode:  ProposeNetworkMsg,
	PrevoteCode:   PrevoteNetworkMsg,
	PrecommitCode: PrecommitNetworkMsg,
}

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
	return fmt.Sprintf("{code: %v, %s, ValidRound: %v, ProposedBlockHash: %v, signer: %v, power: %v}",
		p.Code(), p.base.String(), p.validRound, p.block.Hash().String(), p.signer, p.power)
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
	// proposals are signed by only one validator, signature cannot be 0
	if ext.Signature.IsZero() {
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
	// quorum certificate and commit round of a proposal should always be set to the default value
	if ext.ProposalBlock.HasCommitInformation() {
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

func (p *Propose) PreValidate(committee *types.Committee) error {
	if p.preverified {
		return nil
	}

	validator := committee.MemberByAddress(p.signer)
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

func (p *LightProposal) String() string {
	return fmt.Sprintf("{code: %v, %s, ValidRound: %v, BlockHash: %v, signer: %v, power: %v}", p.Code(), p.base.String(), p.validRound, p.blockHash, p.signer, p.power)
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
	// proposals are signed by only one validator, signature cannot be 0
	if ext.Signature.IsZero() {
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

func (p *LightProposal) PreValidate(committee *types.Committee) error {
	if p.preverified {
		return nil
	}

	validator := committee.MemberByAddress(p.signer)
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
type extVote[T uint16 | uint32] struct {
	// Code is redundant with the p2p.msg code however it is required
	// because we don't want to re-serialize the message again in order
	// to compute the hash value.
	Code      uint8
	Round     uint64
	Height    uint64
	Value     common.Hash
	Signers   *types.SignersBase[T]
	Signature *blst.BlsSignature
}

// TODO: would be good to do the same thing for proposal and lightproposal (to avoid code repetition)
type vote[T uint16 | uint32] struct {
	signers *types.SignersBase[T]
	base
}

func (v *vote[T]) Signers() *types.SignersBase[T] {
	return v.signers
}

func (v *vote[T]) Power() *big.Int {
	return v.signers.Power()
}

func (v *vote[T]) PreValidate(committee *types.Committee) error {
	if v.preverified {
		return nil
	}

	if err := v.signers.Validate(committee.Len()); err != nil {
		return fmt.Errorf("invalid signers information: %w", err)
	}

	// if it is an individual signature, it cannot be 0
	if v.signers.Len() == 1 && v.signature.IsZero() {
		return ErrInvalidIndividualVote
	}

	// compute aggregated key and auxiliary data structures
	indexes := v.signers.FlattenUniq()
	keys := make([]blst.PublicKey, len(indexes))
	powers := make(map[int]*big.Int)
	power := new(big.Int)

	for i, index := range indexes {
		member := committee.Members[index]
		keys[i] = member.ConsensusKey
		powers[index] = member.VotingPower
		power.Add(power, member.VotingPower)
	}

	v.signers.AssignPower(powers, power)
	v.signerKey = v.signers.AggregatePublicKey(keys)
	v.preverified = true
	return nil
}

func (v *vote[T]) String() string {
	return fmt.Sprintf("%s, signers: {%s}",
		v.base.String(), v.signers.String())
}

type Prevote struct {
	value common.Hash
	vote[uint16]
}

func (p *Prevote) Code() uint8 {
	return PrevoteCode
}

func (p *Prevote) Value() common.Hash {
	return p.value
}

func (p *Prevote) String() string {
	return fmt.Sprintf("{code: %v, %s, value: %v}",
		p.Code(), p.vote.String(), p.value)
}

type Precommit struct {
	value common.Hash
	vote[uint16]
}

func (p *Precommit) Code() uint8 {
	return PrecommitCode
}

func (p *Precommit) Value() common.Hash {
	return p.value
}

func (p *Precommit) String() string {
	return fmt.Sprintf("{code: %v, %s, value: %v}",
		p.Code(), p.vote.String(), p.value)
}

type EvidenceVote struct {
	value common.Hash
	vote[uint32]
}

func (e *EvidenceVote) Code() uint8 {
	return EvidenceVoteCode
}

func (e *EvidenceVote) Value() common.Hash {
	return e.value
}

func (e *EvidenceVote) String() string {
	return fmt.Sprintf("{code: %v, %s, value: %v}",
		e.Code(), e.vote.String(), e.value)
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

	signers := types.NewVoteSigners(csize)
	signers.AddMember(self)

	payload, _ := rlp.EncodeToBytes(extVote[uint16]{
		Code:      code,
		Round:     uint64(r), // #nosec
		Height:    h,
		Value:     value,
		Signers:   signers.SignersBase,
		Signature: signature.(*blst.BlsSignature),
	})
	vote := E{
		value: value,
		vote: vote[uint16]{
			signers: signers.SignersBase,
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

// NOTE: this function assumes that:
// 1. all votes are for the same signature input (code,h,r,value)
// 2. all votes have previously been preverified and cryptographically verified
func AggregatePrevotes(votes []Vote) []*Prevote {
	return AggregateVotes[Prevote](votes)
}

// NOTE: this function assumes that:
// 1. all votes are for the same signature input (code,h,r,value)
// 2. all votes have previously been preverified and cryptographically verified
func AggregatePrecommits(votes []Vote) []*Precommit {
	return AggregateVotes[Precommit](votes)
}

// NOTE: this function assumes that:
// 1. all votes are for the same signature input (code,h,r,value)
// 2. all votes have previously been preverified and cryptographically verified
func AggregatePrevotesToEvidence(votes []Vote) *EvidenceVote {
	aggregateSignature, aggregateKey := AggregateVotesToQuorum(votes)

	representative := votes[0]
	h := representative.H()
	r := representative.R()
	value := representative.Value()
	signatureInput := representative.SignatureInput()

	payload, _ := rlp.EncodeToBytes(extVote[uint32]{
		Code:      EvidenceVoteCode,
		Round:     uint64(r), // #nosec
		Height:    h,
		Value:     value,
		Signers:   aggregateSignature.Signers.SignersBase,
		Signature: aggregateSignature.Signature,
	})

	return &EvidenceVote{
		value: value,
		vote: vote[uint32]{
			signers: aggregateSignature.Signers.SignersBase,
			base: base{
				height:         h,
				round:          r,
				signatureInput: signatureInput,
				signature:      aggregateSignature.Signature,
				payload:        payload,
				hash:           crypto.Hash(payload),
				verified:       true, // verified due to all votes being verified
				preverified:    true,
				signerKey:      aggregateKey, // this is not strictly necessary since the vote is already verified
			},
		},
	}
}

// NOTE: this function assumes that:
// 1. all votes are for the same signature input (code,h,r,value)
// 2. all votes have previously been preverified and cryptographically verified
func AggregatePrecommitsToQuorum(votes []Vote) *types.AggregateSignature {
	signature, _ := AggregateVotesToQuorum(votes)
	return signature
}

// NOTE: this function assumes that:
// 1. all votes are for the same signature input (code,h,r,value)
// 2. all votes have previously been preverified and cryptographically verified
func AggregateVotesToQuorum(votes []Vote) (*types.AggregateSignature, blst.PublicKey) {
	// length safety checks
	if len(votes) == 0 {
		panic("Trying to aggregate empty set of votes")
	}

	quorumSigners := types.NewQuorumSigners(votes[0].Signers().CommitteeSize())
	signatures := make([]blst.Signature, 0, len(votes))
	publicKeys := make([]blst.PublicKey, 0, len(votes))

	for _, vote := range votes {
		// merge the vote only if it contributes to the quorumSigners
		voteToQuorum := vote.Signers().ToQuorumSigners()
		if !quorumSigners.AddsInformation(voteToQuorum.SignersBase) {
			continue
		}

		signatures = append(signatures, vote.Signature())
		publicKeys = append(publicKeys, vote.SignerKey())

		// `signers.Merge` will iterate over the other object. So megre the smaller object into the larger one.
		// This will give an approximate total runtime complexity of `O(n * logn)`, where n = number of total elements
		// in all the signers combined
		if quorumSigners.Len() < voteToQuorum.Len() {
			voteToQuorum.Merge(quorumSigners.SignersBase)
			quorumSigners = voteToQuorum
		} else {
			quorumSigners.Merge(voteToQuorum.SignersBase)
		}
	}

	aggregatedSignature := blst.AggregateSignatures(signatures)
	aggregateKey, err := blst.AggregatePublicKeys(publicKeys) // remove?
	if err != nil {
		panic("Cannot generate aggregate public key from valid votes: " + err.Error()) //nolint
	}

	return &types.AggregateSignature{
		Signature: aggregatedSignature.(*blst.BlsSignature),
		Signers:   quorumSigners,
	}, aggregateKey
}

var (
	validVotesCounter     = metrics.GetOrRegisterCounter("aggregator/backend/valid", nil)     // measures time for message passing from backend to aggregator
	aggregateVotesCounter = metrics.GetOrRegisterCounter("aggregator/backend/aggregate", nil) // measures time for message passing from backend to aggregator
)

// NOTE: this function assumes that:
// 1. all votes are for the same signature input (code,h,r,value)
// 2. all votes have previously been preverified and cryptographically verified
//
// set `skipBoundaryCheck[0] = true` only if the boundary checks are done previously
// and the caller is sure that no boundary check is needed anymore.
func AggregateVotes[E Prevote | Precommit](votes []Vote, skipBoundaryCheck ...bool) []*E {
	// TODO: return for len = 1
	// length safety checks
	if len(votes) == 0 {
		panic("Trying to aggregate empty set of votes")
	}

	if metrics.Enabled {
		validVotesCounter.Inc(int64(len(votes)))
	}

	doBoundaryCheck := true
	if len(skipBoundaryCheck) > 0 {
		doBoundaryCheck = !skipBoundaryCheck[0]
	}

	// bitmap to track if the each vote has any contribution
	totalBitMap := types.NewValidatorBitmap(votes[0].Signers().CommitteeSize())

	// compute new aggregated signature and related signers information
	voteDistributed := make([][]Vote, 0)
	signatures := make([][]blst.Signature, 0)
	publicKeys := make([][]blst.PublicKey, 0)
	// Both `Prevote` and `Precommit` have `SignersBase[uint16]`
	aggregateSigners := make([]*types.SignersBase[uint16], 0)
	for _, vote := range votes {
		// do not aggregate votes if they do not add any useful information
		// e.g. signers contains already at least 1 signature for all signers of vote.Signers()
		// we would just create and gossip new aggregates that would uselessly flood the network
		if !totalBitMap.Merge(vote.Signers().Bits) {
			continue
		}
		index := -1
		for i, signers := range aggregateSigners {
			if !doBoundaryCheck {
				index = i
				break
			}
			// we check if the resulting aggregate respects the coefficient boundaries.
			// this avoids that we aggregate two complex aggregates together, which can lead to coefficient breaching.
			if signers.RespectsBoundaries(vote.Signers()) {
				index = i
				break
			}
		}

		if index == -1 {
			index = len(aggregateSigners)
			aggregateSigners = append(aggregateSigners, types.NewSigners[uint16](vote.Signers().CommitteeSize()))
			signatures = append(signatures, make([]blst.Signature, 0))
			publicKeys = append(publicKeys, make([]blst.PublicKey, 0))
			voteDistributed = append(voteDistributed, make([]Vote, 0))
		}

		signatures[index] = append(signatures[index], vote.Signature())
		publicKeys[index] = append(publicKeys[index], vote.SignerKey())
		voteDistributed[index] = append(voteDistributed[index], vote)

		aggregateSigners[index].Merge(vote.Signers())
	}

	aggregates := make([]*E, 0, len(aggregateSigners))

	representative := votes[0]
	c := representative.Code()
	h := representative.H()
	r := representative.R()
	value := representative.Value()
	signatureInput := representative.SignatureInput()

	for i, signers := range aggregateSigners {

		aggregatedSignature := blst.AggregateSignatures(signatures[i])
		aggregatedPublicKey, err := blst.AggregatePublicKeys(publicKeys[i]) // TODO: remove?
		if err != nil {
			panic("Cannot generate aggregate public key from valid votes: " + err.Error()) //nolint
		}

		payload, _ := rlp.EncodeToBytes(extVote[uint16]{
			Code:      c,
			Round:     uint64(r), // #nosec
			Height:    h,
			Value:     value,
			Signers:   signers,
			Signature: aggregatedSignature.(*blst.BlsSignature),
		})

		aggregateVote := E{
			value: value,
			vote: vote[uint16]{
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
		aggregates = append(aggregates, &aggregateVote)
	}

	if metrics.Enabled {
		aggregateVotesCounter.Inc(int64(len(aggregates)))
	}
	return aggregates
}

// NOTE: this function assumes that:
// 1. all votes are for the same signature input (code,h,r,value)
// 2. all votes have previously been preverified and cryptographically verified
//
// returns `true` if the `newVote` can be merged with any of the `nonMergeable` votes.
// if the `newVote` can be merged, then it will also return the index of the `nonMergeable` which
// was merged with `newVote`. `nonMergeable` votes cannot be merged with one another.
//
// NOTE: DO NOT MODIFY `nonMergeable`
func AggregateLastVote[E Prevote | Precommit](nonMergeable []Vote, newVote Vote) (bool, int, *E) {
	for i, vote := range nonMergeable {
		if vote.Signers().RespectsBoundaries(newVote.Signers()) {
			return true, i, AggregateVotes[E]([]Vote{vote, newVote}, true)[0]
		}
	}
	return false, 0, nil
}

func (e *EvidenceVote) DecodeRLP(s *rlp.Stream) error {
	payload, err := s.Raw()
	if err != nil {
		return err
	}

	encoded := &extVote[uint32]{}
	if err := rlp.DecodeBytes(payload, encoded); err != nil {
		return err
	}
	if encoded.Code != EvidenceVoteCode {
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
	if encoded.Signers == nil || encoded.Signers.Bits == nil || encoded.Signers.Coefficients == nil {
		return constants.ErrInvalidMessage
	}
	if encoded.Signers.SanityCheck() != nil {
		return constants.ErrInvalidMessage
	}
	e.height = encoded.Height
	e.round = int64(encoded.Round)
	e.value = encoded.Value
	e.signature = encoded.Signature
	e.signers = encoded.Signers
	e.payload = payload
	// precompute hash and signature hash
	e.signatureInput = VoteSignatureInput(encoded.Height, encoded.Round, EvidenceVoteCode, encoded.Value)
	e.hash = crypto.Hash(payload)
	e.verified = false
	e.preverified = false
	return nil
}

func (p *Prevote) DecodeRLP(s *rlp.Stream) error {
	payload, err := s.Raw()
	if err != nil {
		return err
	}

	encoded := &extVote[uint16]{}
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
	if encoded.Signers == nil || encoded.Signers.Bits == nil || encoded.Signers.Coefficients == nil {
		return constants.ErrInvalidMessage
	}
	if encoded.Signers.SanityCheck() != nil {
		return constants.ErrInvalidMessage
	}
	p.height = encoded.Height
	p.round = int64(encoded.Round)
	p.value = encoded.Value
	p.signature = encoded.Signature
	p.signers = encoded.Signers
	p.payload = payload
	// precompute hash and signature hash
	p.signatureInput = VoteSignatureInput(encoded.Height, encoded.Round, PrevoteCode, encoded.Value)
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
	encoded := &extVote[uint16]{}
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
	if encoded.Signers == nil || encoded.Signers.Bits == nil || encoded.Signers.Coefficients == nil {
		return constants.ErrInvalidMessage
	}
	if encoded.Signers.SanityCheck() != nil {
		return constants.ErrInvalidMessage
	}
	p.height = encoded.Height
	p.round = int64(encoded.Round)
	p.value = encoded.Value
	p.signature = encoded.Signature
	p.signers = encoded.Signers
	p.payload = payload
	// precompute hash and signature hash
	p.signatureInput = VoteSignatureInput(encoded.Height, encoded.Round, PrecommitCode, encoded.Value)
	p.hash = crypto.Hash(payload)
	p.verified = false
	p.preverified = false
	return nil
}

func VoteSignatureInput(h uint64, r uint64, code uint8, v common.Hash) common.Hash {
	signaturePayload, _ := rlp.EncodeToBytes([]any{code, r, h, v})
	return crypto.Hash(signaturePayload)
}

// PrepareCommittedSeal returns the input data to compute the committed seal for a given block hash.
func PrepareCommittedSeal(hash common.Hash, round int64, height *big.Int) common.Hash {
	// this is matching the signature input that we get from the committed messages.
	buf, _ := rlp.EncodeToBytes([]any{PrecommitCode, uint64(round), height.Uint64(), hash})
	return crypto.Hash(buf)
}

// computes the power of a set of messages. Every sender's power is counted only once
func Power(messages []Msg) *big.Int {
	power := NewAggregatedPower()

	for _, msg := range messages {
		switch m := msg.(type) {
		case *Propose:
			power.Set(m.SignerIndex(), m.Power())
		case *Prevote, *Precommit:
			vote := m.(Vote)
			for index, signerPower := range vote.Signers().Powers() {
				power.Set(index, signerPower)
			}
		case *EvidenceVote:
			vote := Msg(m).(Evidence)
			for index, signerPower := range vote.Signers().Powers() {
				power.Set(index, signerPower)
			}
		default:
			panic("unknown message type")
		}
	}
	return power.Power()
}

// OverQuorumVotes compute voting power out from a set of prevotes or precommits of a certain round and height, the caller
// should make sure that the votes belong to a certain round and height, it returns a set of votes that the corresponding
// voting power is over quorum, otherwise it returns nil.
func OverQuorumVotes(msgs []Msg, quorum *big.Int) (overQuorumVotes []Msg) {
	if Power(msgs).Cmp(quorum) >= 0 {
		return msgs
	}
	return nil
}

type Fake struct {
	FakeCode           uint8
	FakeRound          uint64
	FakeHeight         uint64
	FakeValue          common.Hash
	FakePayload        []byte
	FakeHash           common.Hash
	FakeSigners        *types.VoteSigners
	FakeSignature      blst.Signature
	FakeSignatureInput common.Hash
	FakeSignerKey      blst.PublicKey

	// used only for proposal
	FakeBlock         *types.Block
	FakeValidRound    uint64
	FakeValidRoundNil bool
	FakeSigner        common.Address
	FakeSignerIndex   uint64
	FakePower         *big.Int
	FakeVerified      bool // for prevote and precommits this is set to true by default for now
}

func (f Fake) Code() uint8                          { return f.FakeCode }
func (f Fake) R() int64                             { return int64(f.FakeRound) }
func (f Fake) H() uint64                            { return f.FakeHeight }
func (f Fake) Value() common.Hash                   { return f.FakeValue }
func (f Fake) Power() *big.Int                      { return f.FakePower }
func (f Fake) String() string                       { return "{fake}" }
func (f Fake) Hash() common.Hash                    { return f.FakeHash }
func (f Fake) Payload() []byte                      { return f.FakePayload }
func (f Fake) Signature() blst.Signature            { return f.FakeSignature }
func (f Fake) PreValidate(_ *types.Committee) error { return nil }
func (f Fake) Validate() error                      { return nil }
func (f Fake) SignatureInput() common.Hash          { return f.FakeSignatureInput }
func (f Fake) SignerKey() blst.PublicKey            { return f.FakeSignerKey }
func (f Fake) Verified() bool                       { return true }
func (f Fake) PreVerified() bool                    { return true }

func NewFakePropose(f Fake) *Propose {
	var vr int64
	if f.FakeValidRoundNil {
		vr = -1
	} else {
		vr = int64(f.FakeValidRound)
	}
	return &Propose{
		block:       f.FakeBlock,
		validRound:  vr,
		signer:      f.FakeSigner,
		signerIndex: int(f.FakeSignerIndex),
		power:       f.FakePower,
		base: base{
			round:          int64(f.FakeRound),
			height:         f.FakeHeight,
			signatureInput: f.FakeSignatureInput,
			signature:      f.FakeSignature,
			payload:        f.FakePayload,
			hash:           f.FakeHash,
			signerKey:      f.FakeSignerKey,
			preverified:    true,
			verified:       f.FakeVerified,
		},
	}
}

func NewFakePrevote(f Fake) *Prevote {
	return &Prevote{
		value: f.FakeValue,
		vote: vote[uint16]{
			signers: f.FakeSigners.SignersBase,
			base: base{
				round:          int64(f.FakeRound),
				height:         f.FakeHeight,
				signatureInput: f.FakeSignatureInput,
				signature:      f.FakeSignature,
				payload:        f.FakePayload,
				hash:           f.FakeHash,
				signerKey:      f.FakeSignerKey,
				preverified:    true,
				verified:       true,
			},
		},
	}
}

func NewFakePrecommit(f Fake) *Precommit {
	return &Precommit{
		value: f.FakeValue,
		vote: vote[uint16]{
			signers: f.FakeSigners.SignersBase,
			base: base{
				round:          int64(f.FakeRound),
				height:         f.FakeHeight,
				signatureInput: f.FakeSignatureInput,
				signature:      f.FakeSignature,
				payload:        f.FakePayload,
				hash:           f.FakeHash,
				signerKey:      f.FakeSignerKey,
				preverified:    true,
				verified:       true,
			},
		},
	}
}
