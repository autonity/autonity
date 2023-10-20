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

// Package consensus implements different Ethereum consensus engines.
package consensus

import (
	ethereum "github.com/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
)

//// Place consensus message codes in protocol to break dependency circle in between tendermint and fault detectors.
//// Todo(youssef): consider putting them in the message package.
//const (
//	MsgProposal uint8 = iota
//	MsgPrevote
//	MsgPrecommit
//	// MsgLightProposal is only used by accountability that it converts full proposal to a lite one
//	// which contains just meta-data of a proposal for a sustainable on-chain proof mechanism.
//	MsgLightProposal
//)

// setting for Autonity accountability protocol, they are a part of consensus.
const (
	ReportingSlotPeriod       = 20  // Each AFD reporting slot holds 20 blocks, each validator response for a slot.
	DeltaBlocks               = 10  // Wait until the GST + delta blocks to start accounting.
	AccountabilityHeightRange = 256 // Default msg buffer range for AFD.
)

// Broadcaster defines the interface to enqueue blocks to fetcher and find peer
type Broadcaster interface {
	// Enqueue add a block into fetcher queue
	Enqueue(id string, block *types.Block)
	// FindPeers retrives connected peers by addresses
	FindPeers(map[common.Address]struct{}) map[common.Address]ethereum.Peer
}
