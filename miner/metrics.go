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

package miner

import (
	"github.com/autonity/autonity/metrics"
)

var (
	PrepareWorkTimer  = metrics.NewRegisteredTimer("miner/work/prepare", nil)  // time to prepare new work (block proposal)
	FillWorkTimer     = metrics.NewRegisteredTimer("miner/work/fill", nil)     // time to fill new work with txs
	CommitWorkTimer   = metrics.NewRegisteredTimer("miner/work/commit", nil)   // time to commit work (send it to the taskloop)
	FinalizeWorkTimer = metrics.NewRegisteredTimer("miner/work/finalize", nil) // time to finalize work (substep of commit)
	SealWorkTimer     = metrics.NewRegisteredTimer("miner/work/seal", nil)     // time to seal block (taskloop, waits for timestamp to be ripe and then submits to consensus engine)
	CopyWorkTimer     = metrics.NewRegisteredTimer("miner/work/copy", nil)     // time to do task deep copy (see worker ResultLoop()).
	PersistWorkTimer  = metrics.NewRegisteredTimer("miner/work/persist", nil)  // time to writeBlockAndSetHead

	// instant metrics
	PrepareWorkBg      = metrics.NewRegisteredBufferedGauge("miner/work/prepare.bg", nil)          // time to prepare new work (block proposal)
	FillWorkBg         = metrics.NewRegisteredBufferedGauge("miner/work/fill.bg", nil)             // time to fill new work with txs
	CommitWorkBg       = metrics.NewRegisteredBufferedGauge("miner/work/commit.bg", nil)           // time to commit work (send it to the taskloop)
	FinalizeWorkBg     = metrics.NewRegisteredBufferedGauge("miner/work/finalize.bg", nil)         // time to finalize work (substep of commit)
	SealWorkBg         = metrics.NewRegisteredBufferedGauge("miner/work/seal.bg", nil)             // time to seal block (taskloop, waits for timestamp to be ripe and then submits to consensus engine)
	CopyWorkBg         = metrics.NewRegisteredBufferedGauge("miner/work/copy.bg", nil)             // time to do task deep copy (see worker ResultLoop()).
	PersistWorkBg      = metrics.NewRegisteredBufferedGauge("miner/work/persist.bg", nil)          // time to writeBlockAndSetHead
	TotalTaskProcessBg = metrics.NewRegisteredBufferedGauge("miner/work/processTask.bg", nil)      // time to writeBlockAndSetHead
	TotalTaskPrepareBg = metrics.NewRegisteredBufferedGauge("miner/work/prepareTaskTotal.bg", nil) // time to writeBlockAndSetHead
)
