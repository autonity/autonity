package constants

const (
	MaxRound                      = 999 // consequence of backlog priority
	ProposerElectionCycleConstant = 99  // seed for controlling proposer election repetition
	AskSyncInterval               = 5   // the interval in seconds to check the liveness and rise AskSync request.
	DefaultSyncTimeout            = 30  // the default sync timeout of the livenessTracker.
)
