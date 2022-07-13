package types

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/messageutils"
	"github.com/autonity/autonity/core/types"
	"math/big"
)

type CoreStateRequestEvent struct {
	StateChan chan TendermintState
}

// VoteState save the prevote or precommit voting status for a specific value.
type VoteState struct {
	Value            common.Hash
	ProposalVerified bool
	VotePower        uint64
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
	messageutils.Message
	Hash   common.Hash
	Power  uint64
	Height *big.Int
	Round  int64
}

// TendermintState save an instant status for the tendermint consensus engine.
type TendermintState struct {
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
	QuorumVotePower uint64
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
	// backlog unchecked msgs.
	UncheckedMsgs []*MsgForDump
	// Known msg of gossip.
	KnownMsgHash []common.Hash
}
