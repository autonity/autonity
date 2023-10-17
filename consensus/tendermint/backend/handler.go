package backend

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/events"
	"github.com/autonity/autonity/p2p"
	lru "github.com/hashicorp/golang-lru"
	"io"
)

const (
	ProposeMsg        = 0x11
	PrevoteMsg        = 0x12
	PrecommitMsg      = 0x13
	SyncMsg           = 0x14
	AccountabilityMsg = 0x15
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
		if _, err := sb.HandleMsg(addr, msg, nil); err != nil {
			sb.logger.Error("Could not handle cached message", "err", err)
		}
	}
}

// HandleMsg implements consensus.Handler.HandleMsg
func (sb *Backend) HandleMsg(addr common.Address, msg p2p.Msg, errCh chan<- error) (bool, error) {
	if msg.Code < ProposeMsg || msg.Code > AccountabilityMsg {
		return false, nil
	}

	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()

	switch msg.Code {
	case ProposeMsg:
		return handleConsensusMsg[message.Propose](sb, addr, msg, errCh)
	case PrevoteMsg:
		return handleConsensusMsg[message.Prevote](sb, addr, msg, errCh)
	case PrecommitMsg:
		return handleConsensusMsg[message.Precommit](sb, addr, msg, errCh)
	case SyncMsg:
		if !sb.coreStarted {
			sb.logger.Debug("Sync message received but core not running")
			return true, nil // we return nil as we don't want to shut down the connection if core is stopped
		}
		sb.logger.Debug("Received sync message", "from", addr)
		go sb.Post(events.SyncEvent{Addr: addr})
	case AccountabilityMsg:
		if !sb.coreStarted {
			sb.logger.Debug("Accountability Msg received but core not running")
			return true, nil // we return nil as we don't want to shut down the connection if core is stopped
		}
		var data []byte
		if err := msg.Decode(&data); err != nil {
			// this error will freeze peer for 30 seconds by according to dev p2p protocol.
			return true, errDecodeFailed
		}

		// post the off chain accountability msg to the event handler, let the event handler to handle DoS attack vectors.
		sb.logger.Debug("Received Accountability Msg", "from", addr)
		go sb.Post(events.AccountabilityEvent{Sender: addr, Payload: data, ErrCh: errCh})
	default:
		return false, nil
	}

	return true, nil
}

func handleConsensusMsg[M message.Consensus](sb *Backend, addr common.Address, msg p2p.Msg, errCh chan<- error) (bool, error) {
	if !sb.coreStarted {
		buffer := new(bytes.Buffer)
		if _, err := io.Copy(buffer, msg.Payload); err != nil {
			return true, errDecodeFailed
		}
		savedMsg := msg
		savedMsg.Payload = buffer
		sb.pendingMessages.Enqueue(UnhandledMsg{addr: addr, msg: savedMsg})
		return true, nil // return nil to avoid shutting down connection during block sync.
	}
	var consensusMsg M
	if err := msg.Decode(&consensusMsg); err != nil {
		return true, err
	}
	// If reading was fine then cache the original payload to avoid
	// encoding work during gossip
	if _, err := msg.Payload.(*bytes.Reader).Seek(0, io.SeekStart); err != nil {
		return true, err
	}
	payload := make([]byte, msg.Size)
	if _, err := msg.Payload.Read(payload); err != nil {
		return true, err
	}
	hash := sha256.Sum256(payload)
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
	go sb.Post(events.MessageEvent[M]{
		Message: consensusMsg,
		ErrCh:   errCh,
	})
	return true, nil
}

// SetBroadcaster implements consensus.Handler.SetBroadcaster
func (sb *Backend) SetBroadcaster(broadcaster consensus.Broadcaster) {
	sb.Broadcaster = broadcaster
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
