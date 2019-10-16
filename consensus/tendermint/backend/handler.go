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

package backend

import (
	"errors"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/tendermint/events"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/p2p"
	"github.com/hashicorp/golang-lru"
)

const (
	tendermintMsg     = 0x11
	tendermintSyncMsg = 0x12
)

var (
	// errDecodeFailed is returned when decode message fails
	errDecodeFailed = errors.New("fail to decode tendermint message")
)

// Protocol implements consensus.Handler.Protocol
func (sb *Backend) Protocol() (protocolName string, extraMsgCodes uint64) {
	return "tendermint", 2
}

// HandleMsg implements consensus.Handler.HandleMsg
func (sb *Backend) HandleMsg(addr common.Address, msg p2p.Msg) (bool, error) {
	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()

	if msg.Code == tendermintMsg {
		if !sb.coreStarted {
			return true, ErrStoppedEngine
		}

		var data []byte
		if err := msg.Decode(&data); err != nil {
			return true, errDecodeFailed
		}

		hash := types.RLPHash(data)

		// Mark peer's message
		ms, ok := sb.recentMessages.Get(addr)
		var m *lru.ARCCache
		if ok {
			m, _ = ms.(*lru.ARCCache)
		} else {
			m, _ = lru.NewARC(inmemoryMessages)
			sb.recentMessages.Add(addr, m)
		}
		m.Add(hash, true)

		// Mark self known message
		if _, ok := sb.knownMessages.Get(hash); ok {
			return true, nil
		}
		sb.knownMessages.Add(hash, true)

		sb.postEvent(events.MessageEvent{
			Payload: data,
		})

		return true, nil
	}

	if msg.Code == tendermintSyncMsg {
		if !sb.coreStarted {
			sb.logger.Error("Consensus core stopped - can't process sync messages")
			return true, ErrStoppedEngine
		}
		sb.logger.Info("Received sync message", "from", addr)
		sb.postEvent(events.SyncEvent{Addr: addr})
		return true, nil
	}

	return false, nil
}

// SetBroadcaster implements consensus.Handler.SetBroadcaster
func (sb *Backend) SetBroadcaster(broadcaster consensus.Broadcaster) {
	sb.broadcaster = broadcaster
}

func (sb *Backend) NewChainHead() error {
	sb.coreMu.RLock()
	defer sb.coreMu.RUnlock()
	if !sb.coreStarted {
		return ErrStoppedEngine
	}
	sb.postEvent(events.CommitEvent{})
	return nil
}
