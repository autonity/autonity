package constants

import "errors"

var (
	// ErrNotFromProposer is returned when received message is supposed to be from
	// proposer.
	ErrNotFromProposer = errors.New("message does not come from proposer")
	// ErrFutureHeightMessage is returned when curRoundMessages view is earlier than the
	// view of the received message.
	ErrFutureHeightMessage = errors.New("future height message")
	// ErrOldHeightMessage is returned when the received message's view is earlier
	// than curRoundMessages view.
	ErrOldHeightMessage = errors.New("old height message")
	// ErrOldRoundMessage message is returned when message is of the same Height but form a smaller round
	ErrOldRoundMessage = errors.New("same height but old round message")
	// ErrFutureRoundMessage message is returned when message is of the same Height but form a newer round
	ErrFutureRoundMessage = errors.New("same height but future round message")
	// ErrFutureStepMessage message is returned when it's a prevote or precommit message of the same Height same round
	// while the current step is propose.
	ErrFutureStepMessage = errors.New("same round but future step message")
	// ErrInvalidMessage is returned when the message is malformed.
	ErrInvalidMessage = errors.New("invalid message")
	// ErrInvalidSenderOfCommittedSeal is returned when the committed seal is not from the sender of the message.
	ErrInvalidSenderOfCommittedSeal = errors.New("invalid sender of committed seal")
	// ErrFailedDecodeProposal is returned when the PROPOSAL message is malformed.
	ErrFailedDecodeProposal = errors.New("failed to decode PROPOSAL")
	// ErrFailedDecodePrevote is returned when the PREVOTE message is malformed.
	ErrFailedDecodePrevote = errors.New("failed to decode PREVOTE")
	// ErrFailedDecodePrecommit is returned when the PRECOMMIT message is malformed.
	ErrFailedDecodePrecommit = errors.New("failed to decode PRECOMMIT")
	// ErrFailedDecodeVote is returned for when PREVOTE or PRECOMMIT is malformed.
	ErrFailedDecodeVote = errors.New("failed to decode vote")
	// ErrNilPrevoteSent is returned when timer could be stopped in time
	ErrNilPrevoteSent = errors.New("timer expired and nil prevote sent")
	// ErrNilPrecommitSent is returned when timer could be stopped in time
	ErrNilPrecommitSent = errors.New("timer expired and nil precommit sent")
	// ErrMovedToNewRound is returned when timer could be stopped in time
	ErrMovedToNewRound = errors.New("timer expired and new round started")
)
