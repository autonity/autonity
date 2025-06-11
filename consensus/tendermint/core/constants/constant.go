package constants

import "time"

const (
	MaxRound                      = 999              // consequence of backlog priority
	ProposerElectionCycleConstant = 99               // seed for controlling proposer election repetition
	AskSyncInterval               = 5 * time.Second  // the interval to check the liveness and the asking for sync.
	DefaultSyncTimeout            = 30 * time.Second // the default sync timeout of the livenessTracker.
)
