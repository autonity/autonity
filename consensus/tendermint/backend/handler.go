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
	tendermintMsg     = 0x11
	tendermintSyncMsg = 0x12
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
func (sb *Backend) HandleMsg(addr common.Address, msg p2p.Msg) (bool, error) {
	if msg.Code != tendermintMsg && msg.Code != tendermintSyncMsg {
		return false, nil
	}

	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()

	switch msg.Code {
	case tendermintMsg:
		if !sb.coreStarted {
			buffer := new(bytes.Buffer)
			if _, err := io.Copy(buffer, msg.Payload); err != nil {
				return true, errDecodeFailed
			}
			savedMsg := msg
			savedMsg.Payload = buffer
			sb.pendingMessages.Enqueue(UnhandledMsg{addr: addr, msg: savedMsg})
			return true, nil //return nil to avoid shutting down connection during block sync.
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
	case tendermintSyncMsg:
		if !sb.coreStarted {
			sb.logger.Info("Sync message received but core not running")
			return true, nil // we return nil as we don't want to shutdown the connection if core is stopped
		}
		sb.logger.Info("Received sync message", "from", addr)
		sb.postEvent(events.SyncEvent{Addr: addr})
	default:
		return false, nil
	}

	return true, nil
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
