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
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"sync"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/metrics"
	"github.com/autonity/autonity/rlp"
)

var (
	ErrBadSignature          = errors.New("bad signature")
	ErrUnauthorizedAddress   = errors.New("unauthorized address")
	ErrInvalidIndividualVote = errors.New("individual vote has 0 signature")
)

// Internal message codes used by tendermint consensus engine and its accountability module.
const (
	ProposalCode uint8 = iota
	PrevoteCode
	PrecommitCode
	LightProposalCode
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

func (p *Propose) SignatureInput() common.Hash {
	return p.signatureInput
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

func (p *Propose) DecodeRLPPayload(payload []byte, hash common.Hash) error {
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
	p.hash = hash
	p.verified = false
	p.preverified = false
	return nil
}

func (p *Propose) DecodeRLP(s *rlp.Stream) error {
	payload, err := s.Raw()
	if err != nil {
		return err
	}
	return p.DecodeRLPPayload(payload, crypto.Hash(payload))
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

func (p *Propose) SignerKey() blst.PublicKey {
	if !p.preverified {
		panic("Trying to access signer key on not preverified message")
	}
	return p.base.signerKey
}

func (p *Propose) Signature() (blst.Signature, error) {
	return p.base.signature, nil
}

func (p *Propose) PreValidate(committee *types.Committee, _ bool) error {
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

func (p *LightProposal) SignatureInput() common.Hash {
	return p.signatureInput
}

func (p *LightProposal) Power() *big.Int {
	return p.power
}

func (p *LightProposal) String() string {
	return fmt.Sprintf("{code: %v, %s, ValidRound: %v, BlockHash: %v, signer: %v, power: %v}", p.Code(), p.base.String(), p.validRound, p.blockHash, p.signer, p.power)
}

func (p *LightProposal) SignerKey() blst.PublicKey {
	if !p.preverified {
		panic("Trying to access signer key on not preverified message")
	}
	return p.base.signerKey
}

func (p *LightProposal) Signature() (blst.Signature, error) {
	return p.base.signature, nil
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

func (p *LightProposal) DecodeRLPPayload(payload []byte, hash common.Hash) error {

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
	p.hash = hash
	p.verified = false
	p.preverified = false
	return nil
}

func (p *LightProposal) DecodeRLP(s *rlp.Stream) error {
	payload, err := s.Raw()
	if err != nil {
		return err
	}
	return p.DecodeRLPPayload(payload, crypto.Hash(payload))
}

func (p *LightProposal) Signer() common.Address {
	return p.signer
}

func (p *LightProposal) SignerIndex() int {
	return p.signerIndex
}

func (p *LightProposal) PreValidate(committee *types.Committee, _ bool) error {
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
	value   common.Hash
	code    uint8
	base

	// signature caching
	signerKeyOnce      sync.Once `rlp:"-"` // ensures that the signerKey is computed only once
	sigOnce            sync.Once `rlp:"-"`
	sigErr             error     `rlp:"-"`
	signatureBytes     []byte    `rlp:"-"` // used for caching the raw sig bytes at decoding phase
	signatureInputOnce sync.Once `rlp:"-"`
}

func (v *vote) SignatureInput() common.Hash {
	v.signatureInputOnce.Do(func() {
		v.signatureInput = VoteSignatureInput(v.H(), uint64(v.R()), v.code, v.Value()) // #nosec
	})
	return v.signatureInput
}

func (v *vote) Value() common.Hash {
	return v.value
}

func (v *vote) Signature() (blst.Signature, error) {
	v.sigOnce.Do(func() {
		if len(v.signatureBytes) == 0 {
			return
		}
		v.base.signature, v.sigErr = blst.SignatureFromBytes(v.signatureBytes)
		if v.sigErr == nil && v.signers.Len() == 1 && v.base.signature.IsZero() {
			v.sigErr = ErrInvalidIndividualVote
		}
		v.signatureBytes = nil
	})
	return v.base.signature, v.sigErr
}

func (v *vote) SignerKey() blst.PublicKey {
	if !v.preverified {
		panic("Trying to access signer key on not preverified message")
	}
	v.signerKeyOnce.Do(func() {
		if v.base.signerKey != nil {
			return
		}
		keys := make([]blst.PublicKey, 0, v.signers.Len())
		it := v.signers.NewIterator()
		for it.Next() {
			member := v.signers.Committee().MemberByIndex(it.Index())
			keys = append(keys, member.ConsensusKey)
		}
		// keys can't be zero as they are validated in PreValidate
		if len(keys) == 1 {
			v.base.signerKey = keys[0]
		} else {
			v.base.signerKey = v.signers.AggregatePublicKey(keys)
		}
	})
	return v.base.signerKey
}

func (v *vote) Signers() *types.Signers {
	return v.signers
}

func (v *vote) Power() *big.Int {
	return v.signers.Power()
}

var maxCoefficientBg = metrics.GetOrRegisterBufferedGauge("votes/maxcoefficient", nil)

func (v *vote) PreValidate(committee *types.Committee, shouldRespectCap bool) error {
	if v.preverified {
		return nil
	}

	if err := v.signers.Validate(committee); err != nil {
		return fmt.Errorf("invalid signers information: %w", err)
	}

	maxCoefficient := v.signers.MaxCoefficient()
	if shouldRespectCap && maxCoefficient.BitLen() > common.VoteCap {
		return types.ErrInvalidCoefficient
	}

	if metrics.Enabled {
		maxCoefficientBg.Add(maxCoefficient.Int64())
	}

	v.preverified = true
	return nil
}

func (v *vote) Validate() error {
	// loads Signer key if not already done
	v.SignerKey()
	// computes Signature if not already done
	_, err := v.Signature()
	if err != nil {
		return err
	}
	// compute signature input
	v.SignatureInput()
	return v.base.Validate()
}

func (v *vote) String() string {
	return fmt.Sprintf("%s, signers: {%s}",
		v.base.String(), v.signers.String())
}

type Prevote struct {
	vote
}

func (p *Prevote) Code() uint8 {
	return PrevoteCode
}

func (p *Prevote) String() string {
	return fmt.Sprintf("{code: %v, %s, value: %v}",
		p.Code(), p.vote.String(), p.value)
}

type Precommit struct {
	vote
}

func (p *Precommit) Code() uint8 {
	return PrecommitCode
}

func (p *Precommit) String() string {
	return fmt.Sprintf("{code: %v, %s, value: %v}",
		p.Code(), p.vote.String(), p.value)
}

func newVote[
	E Prevote | Precommit,
	PE interface {
		*E
		Msg
	}](r int64, h uint64, value common.Hash, signer Signer, self *types.CommitteeMember, committee *types.Committee) *E {
	code := PE(new(E)).Code()

	// Pay attention that we're adding the message Code to the signature input data.
	signaturePayload, _ := rlp.EncodeToBytes([]any{code, uint64(r), h, value})
	signatureInput := crypto.Hash(signaturePayload)
	signature := signer(signatureInput)

	signers := types.NewSigners(committee)
	signers.AddSigner(self.Index)

	payload, _ := rlp.EncodeToBytes(extVote{
		Code:      code,
		Round:     uint64(r), // #nosec
		Height:    h,
		Value:     value,
		Signers:   signers,
		Signature: signature.(*blst.BlsSignature),
	})
	vote := E{
		vote: vote{
			value:   value,
			code:    code,
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

func NewPrevote(r int64, h uint64, value common.Hash, signer Signer, self *types.CommitteeMember, committee *types.Committee) *Prevote {
	return newVote[Prevote](r, h, value, signer, self, committee)
}

func NewPrecommit(r int64, h uint64, value common.Hash, signer Signer, self *types.CommitteeMember, committee *types.Committee) *Precommit {
	return newVote[Precommit](r, h, value, signer, self, committee)
}

// NOTE: this function assumes that:
// 1. all votes are for the same signature input (code,h,r,value)
// 2. all votes have previously been preverified and cryptographically verified
func AggregatePrevotes(votes []Vote) []*Prevote {
	return AggregateVotes[Prevote](votes, false)
}

// NOTE: this function assumes that:
// 1. all votes are for the same signature input (code,h,r,value)
// 2. all votes have previously been preverified and cryptographically verified
func AggregatePrecommits(votes []Vote) []*Precommit {
	return AggregateVotes[Precommit](votes, false)
}

func AggregatePrevotesSingle(votes []Vote) *Prevote {
	aggregatedVote := AggregateVotes[Prevote](votes, true)
	if len(aggregatedVote) != 1 {
		panic(fmt.Sprintf("aggregation without boundary checking produced %d votes instead of 1", len(aggregatedVote)))
	}
	return aggregatedVote[0]
}

func AggregatePrecommitsSingle(votes []Vote) *Precommit {
	aggregatedVote := AggregateVotes[Precommit](votes, true)
	if len(aggregatedVote) != 1 {
		panic(fmt.Sprintf("aggregation without boundary checking produced %d votes instead of 1", len(aggregatedVote)))
	}
	return aggregatedVote[0]
}

var (
	votesIn  = metrics.GetOrRegisterCounter("aggregator/votesin", nil)
	votesOut = metrics.GetOrRegisterCounter("aggregator/votesout", nil)
)

// NOTE: this function assumes that:
// 1. all votes are for the same signature input (code,h,r,value)
// 2. all votes have previously been cryptographically verified
func AggregateVotes[E Prevote | Precommit](votes []Vote, ignoreBoundaries bool) []*E {
	// length safety checks
	if len(votes) == 0 {
		panic("Trying to aggregate empty set of votes")
	}

	if metrics.Enabled {
		votesIn.Inc(int64(len(votes)))
	}

	// order votes by decreasing number of distinct signers.
	// This ensures that we reduce as much as possible the number of duplicated signatures for the same validator
	sort.Slice(votes, func(i, j int) bool {
		return votes[i].Signers().Len() > votes[j].Signers().Len()
	})

	results := make([]*E, 0)
	votesToProcess := votes

	for len(votesToProcess) > 0 { // this loop should never run more than once, if there is no maxCoefficient cap breach
		representative := votesToProcess[0]
		aggregateSigner := representative.Signers().Copy()
		sig, err := representative.Signature()
		if err != nil {
			continue
		}
		signaturesToAggregate := []blst.Signature{sig}

		var nextVotesToProcess []Vote
		for _, otherVote := range votesToProcess[1:] {
			if !aggregateSigner.AddsInformation(otherVote.Signers()) {
				// aggregate signer is super set of otherVote
				continue
			}
			// since we sort first we should never have other vote as super set of aggregateSignerW
			// overlapping signers, we can aggregate them
			if ignoreBoundaries || aggregateSigner.RespectsBoundaries(otherVote.Signers()) {
				aggregateSigner.Merge(otherVote.Signers())
				otherVoteSig, err := otherVote.Signature()
				if err != nil {
					continue
				}
				signaturesToAggregate = append(signaturesToAggregate, otherVoteSig)
			} else {
				// unmergeable due to coefficient cap breach, Rare case
				nextVotesToProcess = append(nextVotesToProcess, otherVote)
			}
		}
		var aggregatedSignature blst.Signature
		if len(signaturesToAggregate) == 1 {
			aggregatedSignature = signaturesToAggregate[0]
		} else {
			aggregatedSignature = blst.AggregateSignatures(signaturesToAggregate)
		}

		payload, _ := rlp.EncodeToBytes(extVote{
			Code:      representative.Code(),
			Round:     uint64(representative.R()), // nolint
			Height:    representative.H(),
			Value:     representative.Value(),
			Signers:   aggregateSigner,
			Signature: aggregatedSignature.(*blst.BlsSignature),
		})
		mainAggregate := E{
			vote: vote{
				value:   representative.Value(),
				signers: aggregateSigner,
				code:    representative.Code(),
				base: base{
					height:         representative.H(),
					round:          representative.R(),
					signatureInput: representative.SignatureInput(),
					signature:      aggregatedSignature,
					payload:        payload,
					hash:           crypto.Hash(payload),
					verified:       true,
					preverified:    true,
				},
			},
		}

		results = append(results, &mainAggregate)
		votesToProcess = nextVotesToProcess
	}

	if metrics.Enabled {
		votesOut.Inc(int64(len(results)))
	}

	return results
}

func (p *Prevote) DecodeRLPPayload(payload []byte, hash common.Hash) error {
	var s rlp.Stream
	s.Reset(bytes.NewReader(payload), uint64(len(payload)))

	if _, err := s.List(); err != nil { // Begin decoding list
		return constants.ErrInvalidMessage
	}

	// Field 1: Code
	code, err := s.Uint()
	if err != nil {
		return err
	}
	if code != uint64(PrevoteCode) {
		return constants.ErrInvalidMessage
	}

	// Field 2: Round
	round, err := s.Uint()
	if err != nil {
		return err
	}
	if round > constants.MaxRound {
		return constants.ErrInvalidMessage
	}
	p.round = int64(round)

	// Field 3: Height (uint64)
	height, err := s.Uint()
	if err != nil {
		return err
	}
	if height == 0 {
		return constants.ErrInvalidMessage
	}
	p.height = height

	// Field 4: Value (32 bytes)
	valueBytes, err := s.Bytes()
	if err != nil {
		return err
	}
	if len(valueBytes) != common.HashLength {
		return constants.ErrInvalidMessage
	}
	p.value = common.BytesToHash(valueBytes)

	// Field 5: Signers (nested list [Bitmap bytes, Coefficients list])
	if _, err := s.List(); err != nil { // Begin signers list
		return err
	}

	// Bitmap ([]byte)
	bitmap, err := s.Bytes()
	if err != nil {
		return err
	}

	// Coefficients (list of []byte)
	if _, err := s.List(); err != nil {
		return err
	}

	var coefficients []*big.Int
	for { // coefficients list
		coeffBytes, err := s.Bytes()
		if errors.Is(err, rlp.EOL) {
			break
		}
		if err != nil {
			return err
		}
		coefficients = append(coefficients, new(big.Int).SetBytes(coeffBytes))
	}

	if err := s.ListEnd(); err != nil { // End coefficients list
		return err
	}
	if err := s.ListEnd(); err != nil { // End signers list
		return err
	}

	p.signers = &types.Signers{
		Bitmap:       (*types.Bitmap)(new(big.Int).SetBytes(bitmap)),
		Coefficients: coefficients,
	}
	if err := p.signers.SanityCheck(); err != nil {
		return constants.ErrInvalidMessage
	}

	// Field 6: Signature ([]byte)
	sigBytes, err := s.Bytes()
	if err != nil {
		return err
	}
	p.signatureBytes = sigBytes

	if err := s.ListEnd(); err != nil { // End outer list
		return err
	}

	p.hash = hash
	p.code = PrevoteCode
	p.verified = false
	p.preverified = false
	p.payload = payload
	return nil
}

func (p *Prevote) DecodeRLP(s *rlp.Stream) error {
	// Read the raw RLP payload from the stream.
	payload, err := s.Raw()
	if err != nil {
		return err
	}
	return p.DecodeRLPPayload(payload, crypto.Hash(payload))
}

func (p *Precommit) DecodeRLPPayload(payload []byte, hash common.Hash) error {
	var s rlp.Stream
	s.Reset(bytes.NewReader(payload), uint64(len(payload)))

	if _, err := s.List(); err != nil { // Begin decoding list
		return constants.ErrInvalidMessage
	}

	// Field 1: Code
	code, err := s.Uint()
	if err != nil {
		return err
	}
	if code != uint64(PrecommitCode) {
		return constants.ErrInvalidMessage
	}

	// Field 2: Round
	round, err := s.Uint()
	if err != nil {
		return err
	}
	if round > constants.MaxRound {
		return constants.ErrInvalidMessage
	}
	p.round = int64(round)

	// Field 3: Height (uint64)
	height, err := s.Uint()
	if err != nil {
		return err
	}
	if height == 0 {
		return constants.ErrInvalidMessage
	}
	p.height = height

	// Field 4: Value (32 bytes)
	valueBytes, err := s.Bytes()
	if err != nil {
		return err
	}
	if len(valueBytes) != common.HashLength {
		return constants.ErrInvalidMessage
	}
	p.value = common.BytesToHash(valueBytes)

	// Field 5: Signers (nested list [Bitmap bytes, Coefficients list])
	if _, err := s.List(); err != nil { // Begin signers list
		return err
	}

	// Bitmap ([]byte)
	bitmap, err := s.Bytes()
	if err != nil {
		return err
	}

	// Coefficients (list of []byte)
	if _, err := s.List(); err != nil {
		return err
	}

	var coefficients []*big.Int
	for { // coefficients list
		coeffBytes, err := s.Bytes()
		if errors.Is(err, rlp.EOL) {
			break
		}
		if err != nil {
			return err
		}
		coefficients = append(coefficients, new(big.Int).SetBytes(coeffBytes))
	}

	if err := s.ListEnd(); err != nil { // End coefficients list
		return err
	}
	if err := s.ListEnd(); err != nil { // End signers list
		return err
	}

	p.signers = &types.Signers{
		Bitmap:       (*types.Bitmap)(new(big.Int).SetBytes(bitmap)),
		Coefficients: coefficients,
	}
	if err := p.signers.SanityCheck(); err != nil {
		return constants.ErrInvalidMessage
	}

	// Field 6: Signature ([]byte)
	sigBytes, err := s.Bytes()
	if err != nil {
		return err
	}
	p.signatureBytes = sigBytes

	if err := s.ListEnd(); err != nil { // End outer list
		return err
	}

	p.hash = hash
	p.code = PrecommitCode
	p.verified = false
	p.preverified = false
	p.payload = payload
	return nil
}

func (p *Precommit) DecodeRLP(s *rlp.Stream) error {
	// Read the raw RLP payload from the stream.
	payload, err := s.Raw()
	if err != nil {
		return err
	}
	return p.DecodeRLPPayload(payload, crypto.Hash(payload))
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
			signers := m.(Vote).Signers()
			it := signers.NewIterator()
			for it.Next() {
				power.Set(it.Index(), signers.PowerByIndex(it.Index()))
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
	FakeSigners        *types.Signers
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

func (f Fake) Code() uint8                                  { return f.FakeCode }
func (f Fake) R() int64                                     { return int64(f.FakeRound) } //nolint:gosec
func (f Fake) H() uint64                                    { return f.FakeHeight }
func (f Fake) Value() common.Hash                           { return f.FakeValue }
func (f Fake) Power() *big.Int                              { return f.FakePower }
func (f Fake) String() string                               { return "{fake}" }
func (f Fake) Hash() common.Hash                            { return f.FakeHash }
func (f Fake) Payload() []byte                              { return f.FakePayload }
func (f Fake) Signature() (blst.Signature, error)           { return f.FakeSignature, nil }
func (f Fake) PreValidate(_ *types.Committee, _ bool) error { return nil }
func (f Fake) Validate() error                              { return nil }
func (f Fake) SignatureInput() common.Hash                  { return f.FakeSignatureInput }
func (f Fake) SignerKey() blst.PublicKey                    { return f.FakeSignerKey }
func (f Fake) Verified() bool                               { return true }
func (f Fake) PreVerified() bool                            { return true }
func (f Fake) DecodeRLPPayload(payload []byte, _ common.Hash) error {
	fake := &Fake{}
	return rlp.DecodeBytes(payload, fake)
}

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
	prevote := &Prevote{
		vote: vote{
			value:   f.FakeValue,
			signers: f.FakeSigners,
			code:    PrevoteCode,
			base: base{
				round:          int64(f.FakeRound),
				height:         f.FakeHeight,
				signatureInput: f.FakeSignatureInput,
				signature:      f.FakeSignature,
				signerKey:      f.FakeSignerKey,
				payload:        f.FakePayload,
				hash:           f.FakeHash,
				preverified:    true,
				verified:       true,
			},
		},
	}
	return prevote
}

func NewFakePrecommit(f Fake) *Precommit {
	precommit := &Precommit{
		vote: vote{
			signers: f.FakeSigners,
			value:   f.FakeValue,
			code:    PrecommitCode,
			base: base{
				round:          int64(f.FakeRound),
				height:         f.FakeHeight,
				signatureInput: f.FakeSignatureInput,
				signature:      f.FakeSignature,
				signerKey:      f.FakeSignerKey,
				payload:        f.FakePayload,
				hash:           f.FakeHash,
				preverified:    true,
				verified:       true,
			},
		},
	}
	return precommit
}
