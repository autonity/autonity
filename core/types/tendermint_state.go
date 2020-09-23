package types

import (
	"github.com/clearmatics/autonity/common"
	"math/big"
)

// VoteState save the prevote or precommit voting status for a specific value.
type VoteState struct {
	Value            common.Hash `json:"value"            gencodec:"required"`
	ProposalVerified bool        `json:"proposalVerified" gencodec:"required"`
	VotePower        uint64      `json:"votePower"        gencodec:"required"`
}

// RoundState save the voting status for a specific round.
type RoundState struct {
	Round          int64       `json:"round"          gencodec:"required"`
	Proposal       common.Hash `json:"proposal"       gencodec:"required"`
	PrevoteState   []VoteState `json:"prevoteState"   gencodec:"required"`
	PrecommitState []VoteState `json:"precommitState" gencodec:"required"`
}

// TendermintState save an instant status for the tendermint consensus engine.
type TendermintState struct {
	// validator address
	Client common.Address `json:"client"                gencodec:"required"`

	// core state of tendermint
	Height      big.Int     `json:"height"                gencodec:"required"`
	Round       int64       `json:"round"                 gencodec:"required"`
	Step        uint64      `json:"step"                  gencodec:"required"`
	Proposal    common.Hash `json:"proposal"              gencodec:"required"`
	LockedValue common.Hash `json:"lockedValue"           gencodec:"required"`
	LockedRound int64       `json:"lockedRound"           gencodec:"required"`
	ValidValue  common.Hash `json:"validValue"            gencodec:"required"`
	ValidRound  int64       `json:"validRound"            gencodec:"required"`

	// committee state
	ParentCommittee Committee         `json:"parentCommittee"       gencodec:"required"`
	Committee       Committee         `json:"committee"             gencodec:"required"`
	Proposer        common.Address    `json:"proposer"              gencodec:"required"`
	IsProposer      bool              `json:"isProposer"            gencodec:"required"`
	QuorumVotePower uint64            `json:"quorumVotePower"       gencodec:"required"`
	RoundStates     []RoundState      `json:"roundStates"           gencodec:"required"`
	ProposerPolicy  uint64            `json:"proposerPolicy"        gencodec:"required"`

	// extra state
	SentProposal          bool `json:"sentProposal"          gencodec:"required"`
	SentPrevote           bool `json:"sentPrevote"           gencodec:"required"`
	SentPrecommit         bool `json:"sentPrecommit"         gencodec:"required"`
	SetValidRoundAndValue bool `json:"setValidRoundAndValue" gencodec:"required"`

	// timer state
	BlockPeriod           uint64 `json:"blockPeriod"           gencodec:"required"`
	ProposeTimerStarted   bool   `json:"proposeTimerStarted"   gencodec:"required"`
	PrevoteTimerStarted   bool   `json:"prevoteTimerStarted"   gencodec:"required"`
	PrecommitTimerStarted bool   `json:"precommitTimerStared"  gencodec:"required"`

	// current height messages and known message in case of gossip.
	CurHeightMessages []string        `json:"CurHeightMessages"     gencodec:"required"`
	KnownMsgHash   []common.Hash      `json:"KnownMsgHash"          gencodec:"required"`
}
