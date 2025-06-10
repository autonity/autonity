package constants

const (
	MaxRound                      = 999 // consequence of backlog priority
	ProposerElectionCycleConstant = 99  // seed for controlling proposer election repetition
	AskSyncInterval               = 5   // the interval in seconds to check the liveness and rise AskSync request.
	SyncTimeout                   = 30  //
)
