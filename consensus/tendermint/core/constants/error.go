package constants

import "errors"

var (
	// ErrNotFromProposer is returned when received message is supposed to be from
	// proposer.
	ErrNotFromProposer = errors.New("message does not come from proposer")
	// ErrAlreadyHaveProposal is returned when we receive a proposal but we previously already processed one.
	ErrAlreadyHaveProposal = errors.New("a proposal was already processed in the round")
	// ErrHeightClosed is returned when we receive a message for current height, but we already committed a proposal for it.
	ErrHeightClosed = errors.New("consensus instance already concluded")
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
	// ErrNilPrevoteSent is returned when timer could not be stopped in time
	ErrNilPrevoteSent = errors.New("timer expired and nil prevote sent")
	// ErrNilPrecommitSent is returned when timer could not be stopped in time
	ErrNilPrecommitSent = errors.New("timer expired and nil precommit sent")
	// ErrMovedToNewRound is returned when timer could not be stopped in time
	ErrMovedToNewRound = errors.New("timer expired and new round started")
)
