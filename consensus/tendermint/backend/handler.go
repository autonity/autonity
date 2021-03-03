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
	"bytes"
	"context"
	"errors"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/tendermint/events"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/p2p"
	"github.com/hashicorp/golang-lru"
	"io"
)

const (
	TendermintMsg     = 0x11
	TendermintSyncMsg = 0x12
)

type UnhandledMsg struct {
	addr common.Address
	msg  p2p.Msg
}

var (
	// errDecodeFailed is returned when decode message fails
	errDecodeFailed = errors.New("fail to decode tendermint message")
)

// Protocol implements consensus.Handler.Protocol
func (sb *Backend) Protocol() (protocolName string, extraMsgCodes uint64) {
	return "tendermint", 2 //nolint
}

func (sb *Backend) HandleUnhandledMsgs(ctx context.Context) {
	for unhandled := sb.pendingMessages.Dequeue(); unhandled != nil; unhandled = sb.pendingMessages.Dequeue() {
		select {
		case <-ctx.Done():
			return
		default:
			// nothing to do
		}

		addr := unhandled.(UnhandledMsg).addr
		msg := unhandled.(UnhandledMsg).msg
		if _, err := sb.HandleMsg(addr, msg); err != nil {
			sb.logger.Error("could not handle cached message", "err", err)
		}
	}
}

// HandleMsg implements consensus.Handler.HandleMsg
func (sb *Backend) HandleMsg(addr common.Address, msg p2p.Msg) ([]byte, error) {
	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()

	if msg.Code == TendermintSyncMsg {
		if !sb.coreStarted {
			sb.logger.Info("Sync message received but core not running")
			// we return nil as we don't want to shutdown the connection if core is stopped
			return nil, nil
		}
		sb.logger.Info("Received sync message", "from", addr)
		sb.postEvent(events.SyncEvent{Addr: addr})
	}

	if msg.Code == TendermintMsg {

		b := new(bytes.Buffer)
		if _, err := io.Copy(b, msg.Payload); err != nil {
			return nil, errDecodeFailed
		}
		copyPayload := make([]byte, len(b.Bytes()))
		copy(copyPayload, b.Bytes())

		if !sb.coreStarted {
			savedMsg := msg
			savedMsg.Payload = b
			sb.pendingMessages.Enqueue(UnhandledMsg{addr: addr, msg: savedMsg})
			return copyPayload, nil //return nil to avoid shutting down connection during block sync.
		}

		var data []byte
		copyMsg := msg
		buf := new(bytes.Buffer)
		buf.Write(copyPayload)
		copyMsg.Payload = buf
		if err := copyMsg.Decode(&data); err != nil {
			return copyPayload, errDecodeFailed
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
			return copyPayload, nil
		}
		sb.knownMessages.Add(hash, true)

		sb.postEvent(events.MessageEvent{
			Payload: data,
		})
		return copyPayload, nil
	}

	return nil, nil
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
