package core

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
	"math/big"
)

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

// TendermintState save an instant status for the tendermint consensus engine.
type TendermintState struct {
	// return error code, 0 for okay, -1 for timeout.
	Code int64
	// validator address
	Client common.Address

	// core state of tendermint
	Height      big.Int
	Round       int64
	Step        uint64
	Proposal    *common.Hash
	LockedValue *common.Hash
	LockedRound int64
	ValidValue  *common.Hash
	ValidRound  int64

	// committee state
	ParentCommittee types.Committee
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

	// current height messages and known message in case of gossip.
	CurHeightMessages []string
	KnownMsgHash      []common.Hash
}
