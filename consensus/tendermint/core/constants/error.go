package constants

import "errors"

var (
	// ErrNotFromProposer is returned when received message is supposed to be from
	// proposer.
	ErrNotFromProposer = errors.New("message does not come from proposer")
	// ErrAlreadyHaveProposal is returned when we receive a proposal but we previously already processed one.
	ErrAlreadyHaveProposal = errors.New("a proposal was already processed in the round")
	// ErrAlreadyHaveBlock is returned when we are processing a proposal but we already included the proposed block in our local chain.
	ErrAlreadyHaveBlock = errors.New("proposed block is already in our local chain")
	// ErrOldRoundMessage message is returned when message is of the same Height but form a smaller round
	ErrOldRoundMessage = errors.New("same height but old round message")
	// ErrFutureRoundMessage message is returned when message is of the same Height but form a newer round
	ErrFutureRoundMessage = errors.New("same height but future round message")
	// ErrInvalidMessage is returned when the message is malformed.
	ErrInvalidMessage = errors.New("invalid message")
	// ErrNilPrevoteSent is returned when timer could not be stopped in time
	ErrNilPrevoteSent = errors.New("timer expired and nil prevote sent")
	// ErrNilPrecommitSent is returned when timer could not be stopped in time
	ErrNilPrecommitSent = errors.New("timer expired and nil precommit sent")
	// ErrMovedToNewRound is returned when timer could not be stopped in time
	ErrMovedToNewRound = errors.New("timer expired and new round started")
	// ErrRedundantVote is returned when we process a vote that doesn't bring any voting power contribution to Core
	// equivocated votes are considered as redundant, since they don't really signal any progress in the consensus instance
	// they could eventually lead to block finalization, however it is not a guarantee
	ErrRedundantVote = errors.New("vote did not bring any meaningful power contribution to Core")
)
