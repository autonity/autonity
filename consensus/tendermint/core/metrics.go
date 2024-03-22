// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"github.com/autonity/autonity/metrics"
)

var (
	HeightChangeMeter = metrics.NewRegisteredMeter("tendermint/height/change", nil)
	RoundChangeMeter  = metrics.NewRegisteredMeter("tendermint/round/change", nil)
	ProposeTimer      = metrics.NewRegisteredTimer("tendermint/timer/propose", nil)
	PrevoteTimer      = metrics.NewRegisteredTimer("tendermint/timer/prevote", nil)
	PrecommitTimer    = metrics.NewRegisteredTimer("tendermint/timer/precommit", nil)

	// metrics to measure duration of tendermint phases
	HeightTimer            = metrics.NewRegisteredTimer("tendermint/height", nil)             // duration of a height
	RoundTimer             = metrics.NewRegisteredTimer("tendermint/round", nil)              // duration of a round
	ProposeStepTimer       = metrics.NewRegisteredTimer("tendermint/step/propose", nil)       // duration of propose phase
	PrevoteStepTimer       = metrics.NewRegisteredTimer("tendermint/step/prevote", nil)       // duration of prevote phase
	PrecommitStepTimer     = metrics.NewRegisteredTimer("tendermint/step/precommit", nil)     // duration of precommit phase
	PrecommitDoneStepTimer = metrics.NewRegisteredTimer("tendermint/step/precommitDone", nil) // duration of precommit done phase

	ProposalSentTimer     = metrics.NewRegisteredTimer("tendermint/proposal/sent", nil)     // time between round start and proposal sent
	ProposalReceivedTimer = metrics.NewRegisteredTimer("tendermint/proposal/received", nil) // time between round start and proposal received
	ProposalVerifiedTimer = metrics.NewRegisteredTimer("tendermint/proposal/verified", nil) // time to verify proposal
	CommitTimer           = metrics.NewRegisteredTimer("tendermint/commit", nil)            // time between round start and commit (--> block queued for insertion)

	// Instant metrics

	ProposeBg   = metrics.NewRegisteredBufferedGauge("tendermint/bg/propose", nil, nil)
	PrevoteBg   = metrics.NewRegisteredBufferedGauge("tendermint/bg/prevote", nil, nil)
	PrecommitBg = metrics.NewRegisteredBufferedGauge("tendermint/bg/precommit", nil, nil)

	// metrics to measure duration of tendermint phases
	HeightBg            = metrics.NewRegisteredBufferedGauge("tendermint/height.bg", nil, nil)             // duration of a height
	RoundBg             = metrics.NewRegisteredBufferedGauge("tendermint/round.bg", nil, nil)              // duration of a round
	ProposeStepBg       = metrics.NewRegisteredBufferedGauge("tendermint/step/propose.bg", nil, nil)       // duration of propose phase
	PrevoteStepBg       = metrics.NewRegisteredBufferedGauge("tendermint/step/prevote.bg", nil, nil)       // duration of prevote phase
	PrecommitStepBg     = metrics.NewRegisteredBufferedGauge("tendermint/step/precommit.bg", nil, nil)     // duration of precommit phase
	PrecommitDoneStepBg = metrics.NewRegisteredBufferedGauge("tendermint/step/precommitDone.bg", nil, nil) // duration of precommit done phase

	ProposalSentBg                 = metrics.NewRegisteredBufferedGauge("tendermint/proposal/sent.bg", nil, getIntPointer(1024))                      // time between round start and proposal sent
	ProposalSentBlockTSDeltaBg     = metrics.NewRegisteredBufferedGauge("tendermint/proposal/relative/blockTS/sent.bg", nil, getIntPointer(1024))     // time between round start and proposal sent
	ProposalReceivedBg             = metrics.NewRegisteredBufferedGauge("tendermint/proposal/received.bg", nil, getIntPointer(1024))                  // time between round start and proposal received
	ProposalReceivedBlockTSDeltaBg = metrics.NewRegisteredBufferedGauge("tendermint/proposal/relative/blockTS/received.bg", nil, getIntPointer(1024)) // time between round start and proposal sent
	ProposalVerifiedBg             = metrics.NewRegisteredBufferedGauge("tendermint/proposal/verified.bg", nil, getIntPointer(1024))                  // time to verify proposal
	PrevoteSentBg                  = metrics.NewRegisteredBufferedGauge("tendermint/prevote/sent.bg", nil, getIntPointer(1024))                       // time between round start and prevote sent
	PrevoteSentBlockTSDeltaBg      = metrics.NewRegisteredBufferedGauge("tendermint/prevote/relative/blockTS/sent.bg", nil, getIntPointer(1024))      // time between round start and prevote sent
	PrevoteQuorumReceivedBg        = metrics.NewRegisteredBufferedGauge("tendermint/prevote/quorum/received.bg", nil, nil)                            // time between round start and prevote quorum received
	PrevoteQuorumBlockTSDeltaBg    = metrics.NewRegisteredBufferedGauge("tendermint/prevote/relative/blockTS/quorum/received.bg", nil, nil)           // time between round start and prevote quorum received
	PrecommitSentBg                = metrics.NewRegisteredBufferedGauge("tendermint/precommit/sent.bg", nil, getIntPointer(1024))                     // time between round start and precommit sent
	PrecommitSentBlockTSDeltaBg    = metrics.NewRegisteredBufferedGauge("tendermint/precommit/relative/blockTS/sent.bg", nil, getIntPointer(1024))    // time between round start and precommit sent
	PrecommitQuorumReceivedBg      = metrics.NewRegisteredBufferedGauge("tendermint/precommit/quorum/received.bg", nil, nil)                          // time between round start and precommit quorum received
	PrecommitQuorumBlockTSDeltaBg  = metrics.NewRegisteredBufferedGauge("tendermint/precommit/relative/blockTS/quorum/received.bg", nil, nil)         // time between round start and precommit quorum received
	CommitBg                       = metrics.NewRegisteredBufferedGauge("tendermint/commit.bg", nil, nil)                                             // time between round start and commit (--> block queued for insertion)
)

func getIntPointer(val int) *int {
	return &val
}
