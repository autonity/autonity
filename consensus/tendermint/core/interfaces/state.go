package interfaces

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"math/big"
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

// MsgWithHash save the msg and extra field to be marshal to JSON.
type MsgForDump struct {
	message.Msg
	Hash   common.Hash
	Power  *big.Int
	Height *big.Int
	Round  int64
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

	// current height messages.
	CurHeightMessages []*MsgForDump
	// backlog msgs
	BacklogMessages []*MsgForDump
	// backlog of future height msgs.
	FutureMsgs []*MsgForDump
	// Known msg of gossip.
	KnownMsgHash []common.Hash
}
