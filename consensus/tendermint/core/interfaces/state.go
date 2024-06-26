package interfaces

import (
	"math/big"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
)

// VoteState save the prevote or precommit voting status for a specific value.
type VoteState struct {
	Value            common.Hash
	ProposalVerified bool
	VotePower        *big.Int
}

// RoundState save the voting status for a specific round.
type RoundState struct {
	Round          int64
	Proposal       common.Hash
	PrevoteState   []VoteState
	PrecommitState []VoteState
}

// MsgForDump publicly exports the key fields of a message, to allow for json serialization
type MsgForDump struct {
	Code           uint8
	Hash           common.Hash
	Payload        []byte
	Height         uint64
	Round          int64
	SignatureInput common.Hash
	Signature      []byte
	Power          *big.Int

	// only for votes
	Signers *types.Signers
	Value   common.Hash

	// only for proposals
	Block      *types.Block
	ValidRound int64
	Signer     common.Address
}

// TendermintState save an instant status for the tendermint consensus engine.
type CoreState struct {
	// validator address
	Client common.Address

	// Core state of tendermint
	Height      *big.Int
	Round       int64
	Step        uint64
	Proposal    *common.Hash
	LockedValue *common.Hash
	LockedRound int64
	ValidValue  *common.Hash
	ValidRound  int64

	// committee state
	Committee       types.Committee
	Proposer        common.Address
	IsProposer      bool
	QuorumVotePower *big.Int
	RoundStates     []RoundState
	ProposerPolicy  uint64

	// extra state
	SentProposal          bool
	SentPrevote           bool
	SentPrecommit         bool
	SetValidRoundAndValue bool

	// timer state
	BlockPeriod           uint64
	ProposeTimerStarted   bool
	PrevoteTimerStarted   bool
	PrecommitTimerStarted bool

	// current height messages payloads
	CurHeightMessages []MsgForDump
	// future round messages payloads
	FutureRoundMessages []MsgForDump
	// Known msg of gossip.
	KnownMsgHash []common.Hash
}
