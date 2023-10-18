package message

import (
	"fmt"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/log"
	"github.com/pkg/errors"
	"io"
	"math/big"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/rlp"
)

var (
	errInvalidRoundProposal = errors.New("proposal with invalid round")
	errEmptyBlockProposal   = errors.New("proposal with empty block")
	errInvalidValidRound    = errors.New("proposal with invalid isValidround")
)

type ConsensusMsg interface {
	R() int64
	H() *big.Int
	V() common.Hash
}

type NewConsensusMsg struct {
	Round  int64
	Height *big.Int
	Value  common.Hash

	Code          uint8
	Payload       []byte // rlp encoded tendermint msgs: proposal, prevote, precommit.
	Address       common.Address
	Signature     []byte
	CommittedSeal []byte

	Power *big.Int
}

//func (m *NewConsensusMsg) Validate(validateSig SigVerifier, previousHeader *types.Header) error {
//	if previousHeader.Number.Uint64()+1 != m.Height.Uint64() {
//		// don't know why the legacy code panic here, it introduces live-ness issue of the network.
//		// youssef: that is really bad and should never happen, could be because of a race-condition
//		// I'm reintroducing the panic to check if this scenario happens in the wild. We must never fail silently.
//		panic("Autonity has encountered a problem which led to an inconsistent state, please report this issue.")
//		//return fmt.Errorf("inconsistent message verification")
//	}
//	signature := m.Signature
//	payload, err := m.BytesNoSignature()
//	if err != nil {
//		return err
//	}
//
//	if lp, ok := m.ConsensusMsg.(*LightProposal); ok {
//		// in the case of a light proposal, the signature that matters is the inner-one.
//		payload = lp.BytesNoSignature()
//		signature = lp.Signature
//	}
//
//	recoveredAddress, err := validateSig(previousHeader, payload, signature)
//	if err != nil {
//		return err
//	}
//	// ensure message was signed by the sender
//	if m.Address != recoveredAddress {
//		return ErrBadSignature
//	}
//	validator := previousHeader.CommitteeMember(recoveredAddress)
//	// validateSig check as well if the header is in the committee, so this seems unnecessary
//	if validator == nil {
//		return ErrUnauthorizedAddress
//	}
//
//	// check if the lite proposal signature inside the proposal is correct or not.
//	if proposal, ok := m.ConsensusMsg.(*Proposal); ok {
//		if err := proposal.VerifyLightProposalSignature(m.Address); err != nil {
//			return err
//		}
//	}
//
//	m.Power = validator.VotingPower
//	return nil
//}

func (m *NewConsensusMsg) BytesNoSignature() ([]byte, error) {
	// youssef: not sure if the returned error is necessary here as we are in control of the object.
	return rlp.EncodeToBytes(&Message{
		Code:          m.Code,
		Payload:       m.Payload,
		Address:       m.Address,
		Signature:     []byte{},
		CommittedSeal: m.CommittedSeal,
	})
}

// LightProposal is only used by accountability that it converts a Proposal to a LightProposal for sustainable on-chain proof.
type LightProposal struct {
	//NewConsensusMsg
	Round      int64
	Height     *big.Int
	ValidRound int64
	Value      common.Hash // the hash of the proposalBlock.
	Signature  []byte      // the signature computes from the tuple: (Round, Height, ValidRound, ProposalBlock.Hash())
}

func (lp *LightProposal) BytesNoSignature() []byte {
	bytes, _ := rlp.EncodeToBytes(&LightProposal{
		Round:      lp.Round,
		Height:     lp.Height,
		ValidRound: lp.ValidRound,
		Value:      lp.Value,
	})
	return bytes
}

func (lp *LightProposal) VerifySignature(signer common.Address) error {
	payload := lp.BytesNoSignature()
	// 1. Keccak data
	hashData := crypto.Keccak256(payload)
	// 2. Recover public key
	pubkey, err := crypto.SigToPub(hashData, lp.Signature)
	if err != nil {
		return err
	}
	if crypto.PubkeyToAddress(*pubkey) != signer {
		return ErrUnauthorizedAddress
	}
	return nil
}

func (lp *LightProposal) R() int64 {
	return lp.Round
}

func (lp *LightProposal) H() *big.Int {
	return lp.Height
}

func (lp *LightProposal) V() common.Hash {
	return lp.Value
}

// EncodeRLP serializes b into the Ethereum RLP format.
func (lp *LightProposal) EncodeRLP(w io.Writer) error {
	isValidRoundNil := false
	var validRound uint64
	if lp.ValidRound == -1 {
		validRound = 0
		isValidRoundNil = true
	} else {
		validRound = uint64(lp.ValidRound)
	}
	return rlp.Encode(w, []any{uint64(lp.Round), lp.Height, validRound, isValidRoundNil, lp.Value, lp.Signature})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (lp *LightProposal) DecodeRLP(s *rlp.Stream) error {
	var lite struct {
		Round           uint64
		Height          *big.Int
		ValidRound      uint64
		IsValidRoundNil bool
		Value           common.Hash
		Signature       []byte
	}
	if err := s.Decode(&lite); err != nil {
		return err
	}
	var validRound int64
	if lite.IsValidRoundNil {
		if lite.ValidRound != 0 {
			return errors.New("bad lite proposal validRound with isValidround nil")
		}
		validRound = -1
	} else {
		validRound = int64(lite.ValidRound)
	}
	if !(validRound <= constants.MaxRound && lite.Round <= constants.MaxRound) {
		return errInvalidRoundProposal
	}
	lp.Round = int64(lite.Round)
	lp.Height = lite.Height
	lp.ValidRound = validRound
	lp.Value = lite.Value
	lp.Signature = lite.Signature
	return nil
}

type Proposal struct {
	Round          int64
	Height         *big.Int
	ValidRound     int64
	ProposalBlock  *types.Block
	LightSignature []byte // the signature computes from the hash of tuple:(Round, Height, ValidRound, ProposalBlock.Hash())
}

//type NewProposal struct {
//	NewConsensusMsg
//	ValidRound     int64
//	ProposalBlock  *types.Block
//	LightSignature []byte // the signature computes from the hash of tuple:(Round, Height, ValidRound, ProposalBlock.Hash())
//}

func (p *Proposal) String() string {
	return fmt.Sprintf("{Round: %v, Height: %v, ValidRound: %v, ProposedBlockHash: %v}",
		p.Round, p.Height.Uint64(), p.ValidRound, p.ProposalBlock.Hash().String())
}

func (p *Proposal) V() common.Hash {
	return p.ProposalBlock.Hash()
}

func (p *Proposal) R() int64 {
	return p.Round
}

func (p *Proposal) H() *big.Int {
	return p.Height
}

func NewProposal(r int64, h *big.Int, vr int64, p *types.Block, signer func([]byte) ([]byte, error)) *Proposal {
	lightProposal := &LightProposal{
		Round:      r,
		Height:     h,
		ValidRound: vr,
		Value:      p.Hash(),
	}
	lightSignature, _ := signer(lightProposal.BytesNoSignature())
	return &Proposal{
		Round:         r,
		Height:        h,
		ValidRound:    vr,
		ProposalBlock: p,
		// We are adding a "light" signature in here for accountability purposes, the
		// reason is that we can't afford submitting full block proposal on-chains.
		// This could be avoided if we passed the block outside the proposal object
		LightSignature: lightSignature,
	}
}

func (p *Proposal) EncodeRLP(w io.Writer) error {
	if p.ProposalBlock == nil {
		log.Crit("encoding proposal with empty block")
	}
	// RLP encoding doesn't support negative big.Int, so we have to pass one additional field to represents validRound = -1.
	// Note that we could have as well indexed rounds starting by 1, but we want to stay close as possible to the spec.
	isValidRoundNil := false
	var validRound uint64
	if p.ValidRound == -1 {
		validRound = 0
		isValidRoundNil = true
	} else {
		validRound = uint64(p.ValidRound)
	}
	// todo(youssef): not sure how the encoder caching works with anonymous types
	return rlp.Encode(w, []any{
		uint64(p.Round),
		p.Height,
		validRound,
		isValidRoundNil,
		p.ProposalBlock,
		p.LightSignature,
	})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (p *Proposal) DecodeRLP(s *rlp.Stream) error {
	var proposal struct {
		Round           uint64
		Height          *big.Int
		ValidRound      uint64
		IsValidRoundNil bool
		ProposalBlock   *types.Block
		LightSignature  []byte
	}
	if err := s.Decode(&proposal); err != nil {
		return err
	}
	var validRound int64
	if proposal.IsValidRoundNil {
		if proposal.ValidRound != 0 {
			return errInvalidValidRound
		}
		validRound = -1
	} else {
		validRound = int64(proposal.ValidRound)
	}
	if !(validRound <= constants.MaxRound && proposal.Round <= constants.MaxRound) {
		return errInvalidRoundProposal
	}
	if proposal.ProposalBlock == nil {
		return errEmptyBlockProposal
	}
	p.Round = int64(proposal.Round)
	p.Height = proposal.Height
	p.ValidRound = validRound
	p.ProposalBlock = proposal.ProposalBlock
	p.LightSignature = proposal.LightSignature
	return nil
}

// VerifyLightProposalSignature checks that the lite proposal signature inside the proposal is correct or not
func (p *Proposal) VerifyLightProposalSignature(sender common.Address) error {
	lightProposal := &LightProposal{
		Round:      p.Round,
		Height:     p.Height,
		ValidRound: p.ValidRound,
		Value:      p.V(),
		Signature:  p.LightSignature,
	}
	return lightProposal.VerifySignature(sender)
}

type Vote struct {
	Round             uint64
	Height            *big.Int
	ProposedBlockHash common.Hash
}

func (sub *Vote) V() common.Hash {
	return sub.ProposedBlockHash
}

func (sub *Vote) R() int64 {
	return int64(sub.Round)
}

func (sub *Vote) H() *big.Int {
	return sub.Height
}

// EncodeRLP serializes b into the Ethereum RLP format.
//func (sub *Vote) EncodeRLP(w io.Writer) error {
//	return rlp.Encode(w, []any{uint64(sub.Round), sub.Height, sub.ProposedBlockHash})
//}
//
//// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
//func (sub *Vote) DecodeRLP(s *rlp.Stream) error {
//	var vote struct {
//		Round             uint64
//		Height            *big.Int
//		ProposedBlockHash common.Hash
//	}
//
//	if err := s.Decode(&vote); err != nil {
//		return err
//	}
//	sub.Round = int64(vote.Round)
//	if sub.Round > constants.MaxRound {
//		return constants.ErrInvalidMessage
//	}
//	sub.Height = vote.Height
//	sub.ProposedBlockHash = vote.ProposedBlockHash
//	return nil
//}

func (sub *Vote) String() string {
	return fmt.Sprintf("{Round: %v, Height: %v ProposedBlockHash: %v}", sub.Round, sub.Height, sub.ProposedBlockHash.String())
}
