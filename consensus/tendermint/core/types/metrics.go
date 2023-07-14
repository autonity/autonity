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

package types

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

	ProposalSent          = metrics.NewRegisteredTimer("tendermint/proposal/sent", nil)     // time between round start and proposal sent
	ProposalReceivedTimer = metrics.NewRegisteredTimer("tendermint/proposal/received", nil) // time between round start and proposal received
	ProposalVerifiedTimer = metrics.NewRegisteredTimer("tendermint/proposal/verified", nil) // time to verify proposal
	CommitTimer           = metrics.NewRegisteredTimer("tendermint/commit", nil)            // time between round start and commit (--> block queued for insertion)
)
