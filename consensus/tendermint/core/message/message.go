package message

import (
	"errors"
	"fmt"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/p2p"
	"github.com/autonity/autonity/rlp"
	"io"
	"math/big"
)

var (
	ErrBadSignature        = errors.New("bad signature")
	ErrUnauthorizedAddress = errors.New("unauthorized address")
)

const (
	ProposalCode uint8 = iota
	PrevoteCode
	PrecommitCode
	// LightProposalCode is only used by accountability that it converts full proposal to a lite one
	// which contains just meta-data of a proposal for a sustainable on-chain proof mechanism.
	LightProposalCode
)

type Signer func(hash common.Hash) (sig []byte, err error)

type Msg interface {
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

func NewPropose(r int64, h uint64, vr int64, block *types.Block, signer func(hash common.Hash) ([]byte, error)) *Propose {
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
	signature, _ := signer(crypto.Hash(signatureInputEncoded))

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
			signature:      signature,
			payload:        payload,
			hash:           crypto.Hash(payload),
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
		blockHash:  proposal.hash,
		validRound: proposal.validRound,
		base: base{
			round:     proposal.round,
			height:    proposal.height,
			signature: proposal.signature,
			payload:   payload,
			power:     proposal.Power(),
			sender:    proposal.sender,
			hash:      crypto.Hash(payload),
		},
	}
}

func (p *LightProposal) DecodeRLP(s *rlp.Stream) error {
	ext := &extLightProposal{}
	if err := s.Decode(ext); err != nil {
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
	return nil
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
	base
}

func (p *Prevote) Code() uint8 {
	return PrevoteCode
}

func (p *Prevote) Value() common.Hash {
	return p.value
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

func (p *Precommit) String() string {
	return fmt.Sprintf("{r:  %v, h: %v , sender: %v, power: %v, Code: %v, value: %v}",
		p.round, p.height, p.sender, p.power, p.Code(), p.value)
}

func newVote[
	E Prevote | Precommit,
	PE interface {
		*E
		Msg
	}](r int64, h uint64, value common.Hash, signer func(hash common.Hash) ([]byte, error)) *E {
	code := PE(new(E)).Code()
	// Pay attention that we're adding the message Code to the signature input data.
	signatureInput := []any{code, uint64(r), h, value}
	signatureEncodedInput, _ := rlp.EncodeToBytes(signatureInput)
	signature, _ := signer(crypto.Hash(signatureEncodedInput))
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
			payload:        payload,
			hash:           crypto.Hash(payload),
			signatureInput: signatureInput,
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
	Msg
}](p2pMsg p2p.Msg) (PT, error) {
	message := PT(new(T))
	if err := p2pMsg.Decode(message); err != nil {
		return message, err
	}
	if _, err := p2pMsg.Payload.(io.Seeker).Seek(0, io.SeekStart); err != nil {
		return message, err
	}
	payload := make([]byte, p2pMsg.Size)
	if _, err := p2pMsg.Payload.Read(payload); err != nil {
		return message, err
	}
	message.setPayload(payload)
	return message, nil
}

func (b *base) Sender() common.Address {
	if b.sender == (common.Address{}) {
		panic("sender is not set")
	}
	return b.sender
}

func (b *base) H() uint64 {
	return b.height
}

func (b *base) setPayload(payload []byte) {
	b.payload = payload
	b.hash = crypto.Hash(payload)
}

func (b *base) EncodeRLP(w io.Writer) error {
	_, err := w.Write(b.payload)
	return err
}

func (b *base) R() int64 {
	return b.round
}

func (b *base) Power() *big.Int {
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
func (f Fake) setPayload(i []byte)                                            {}
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
