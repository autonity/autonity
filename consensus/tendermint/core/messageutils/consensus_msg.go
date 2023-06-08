package messageutils

import (
	"fmt"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/crypto"
	"github.com/pkg/errors"
	"io"
	"math/big"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/rlp"
)

type ConsensusMsg interface {
	R() int64
	H() *big.Int
	V() common.Hash
}

// LiteProposal is only used by accountability that it converts a Proposal to a LiteProposal for sustainable on-chain proof.
type LiteProposal struct {
	Round      int64
	Height     *big.Int
	ValidRound int64
	Value      common.Hash // the hash of the proposalBlock.
	Signature  []byte      // the signature computes from the tuple: (Round, Height, ValidRound, ProposalBlock.Hash())
}

func (lp *LiteProposal) PayloadNoSig() ([]byte, error) {
	return rlp.EncodeToBytes(&LiteProposal{
		Round:      lp.Round,
		Height:     lp.Height,
		ValidRound: lp.ValidRound,
		Value:      lp.Value,
	})
}

func (lp *LiteProposal) ValidSignature(signer common.Address) error {
	payload, err := lp.PayloadNoSig()
	if err != nil {
		return err
	}
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

func (lp *LiteProposal) R() int64 {
	return lp.Round
}

func (lp *LiteProposal) H() *big.Int {
	return lp.Height
}

func (lp *LiteProposal) V() common.Hash {
	return lp.Value
}

func (lp *LiteProposal) VR() int64 {
	return lp.ValidRound
}

func (lp *LiteProposal) Sig() []byte {
	return lp.Signature
}

// EncodeRLP serializes b into the Ethereum RLP format.
func (lp *LiteProposal) EncodeRLP(w io.Writer) error {
	isValidRoundNil := false
	var validRound uint64
	if lp.ValidRound == -1 {
		validRound = 0
		isValidRoundNil = true
	} else {
		validRound = uint64(lp.ValidRound)
	}
	return rlp.Encode(w, []interface{}{uint64(lp.Round), lp.Height, validRound, isValidRoundNil, lp.Value, lp.Signature})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (lp *LiteProposal) DecodeRLP(s *rlp.Stream) error {
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
		return errors.New("bad proposal with invalid rounds")
	}
	lp.Round = int64(lite.Round)
	lp.Height = lite.Height
	lp.ValidRound = validRound
	lp.Value = lite.Value
	lp.Signature = lite.Signature
	return nil
}

type Proposal struct {
	Round         int64
	Height        *big.Int
	ValidRound    int64
	ProposalBlock *types.Block
	LiteSig       []byte // the signature computes from the hash of tuple:(Round, Height, ValidRound, ProposalBlock.Hash())
}

func (p *Proposal) String() string {
	return fmt.Sprintf("{Round: %v, Height: %v, ValidRound: %v, ProposedBlockHash: %v}",
		p.Round, p.Height.Uint64(), p.ValidRound, p.ProposalBlock.Hash().String())
}

func (p *Proposal) V() common.Hash {
	return p.ProposalBlock.Hash()
}

func (p *Proposal) VR() int64 {
	return p.ValidRound
}

func (p *Proposal) LiteSignature() []byte {
	return p.LiteSig
}

func (p *Proposal) R() int64 {
	return p.Round
}

func (p *Proposal) H() *big.Int {
	return p.Height
}

func NewProposal(r int64, h *big.Int, vr int64, p *types.Block) *Proposal {
	return &Proposal{
		Round:         r,
		Height:        h,
		ValidRound:    vr,
		ProposalBlock: p,
	}
}

// EncodeRLP RLP encoding doesn't support negative big.Int, so we have to pass one additionnal field to represents validRound = -1.
// Note that we could have as well indexed rounds starting by 1, but we want to stay close as possible to the spec.
func (p *Proposal) EncodeRLP(w io.Writer) error {
	if p.ProposalBlock == nil {
		// Should never happen
		return errors.New("encoderlp proposal with nil block")
	}

	isValidRoundNil := false
	var validRound uint64
	if p.ValidRound == -1 {
		validRound = 0
		isValidRoundNil = true
	} else {
		validRound = uint64(p.ValidRound)
	}

	return rlp.Encode(w, []interface{}{
		uint64(p.Round),
		p.Height,
		validRound,
		isValidRoundNil,
		p.ProposalBlock,
		p.LiteSig,
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
		Signature       []byte
	}

	if err := s.Decode(&proposal); err != nil {
		return err
	}
	var validRound int64
	if proposal.IsValidRoundNil {
		if proposal.ValidRound != 0 {
			return errors.New("bad proposal validRound with isValidround nil")
		}
		validRound = -1
	} else {
		validRound = int64(proposal.ValidRound)
	}

	if !(validRound <= constants.MaxRound && proposal.Round <= constants.MaxRound) {
		return errors.New("bad proposal with invalid rounds")
	}

	if proposal.ProposalBlock == nil {
		return errors.New("bad proposal with nil decoded block")
	}

	p.Round = int64(proposal.Round)
	p.Height = proposal.Height
	p.ValidRound = validRound
	p.ProposalBlock = proposal.ProposalBlock
	p.LiteSig = proposal.Signature

	return nil
}

type BadProposalInfo struct {
	Sender common.Address
	Value  common.Hash
}

type Vote struct {
	Round             int64
	Height            *big.Int
	ProposedBlockHash common.Hash
	MaliciousProposer common.Address
	MaliciousValue    common.Hash
}

func (sub *Vote) V() common.Hash {
	return sub.ProposedBlockHash
}

func (sub *Vote) R() int64 {
	return sub.Round
}

func (sub *Vote) H() *big.Int {
	return sub.Height
}

// EncodeRLP serializes b into the Ethereum RLP format.
func (sub *Vote) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{uint64(sub.Round), sub.Height, sub.ProposedBlockHash, sub.MaliciousProposer, sub.MaliciousValue})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (sub *Vote) DecodeRLP(s *rlp.Stream) error {
	var vote struct {
		Round             uint64
		Height            *big.Int
		ProposedBlockHash common.Hash
		MaliciousProposer common.Address
		MaliciousValue    common.Hash
	}

	if err := s.Decode(&vote); err != nil {
		return err
	}
	sub.Round = int64(vote.Round)
	if sub.Round > constants.MaxRound {
		return constants.ErrInvalidMessage
	}
	sub.Height = vote.Height
	sub.ProposedBlockHash = vote.ProposedBlockHash
	sub.MaliciousProposer = vote.MaliciousProposer
	sub.MaliciousValue = vote.MaliciousValue
	return nil
}

func (sub *Vote) String() string {
	return fmt.Sprintf("{Round: %v, Height: %v ProposedBlockHash: %v}", sub.Round, sub.Height, sub.ProposedBlockHash.String())
}
