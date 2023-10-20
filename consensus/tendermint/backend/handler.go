package backend

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/backend/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/events"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/p2p"
	"github.com/autonity/autonity/rlp"
	lru "github.com/hashicorp/golang-lru"
	"io"
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
	if msg.Code != constants.TendermintMsgLightProposal &&
		msg.Code != constants.TendermintMsgVote &&
		msg.Code != constants.TendermintMsgProposal &&
		msg.Code != constants.SyncMsg &&
		msg.Code != constants.AccountabilityMsg {
		return false, nil
	}

	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()

	switch msg.Code {
	//case TendermintMsg:
	//	if !sb.coreStarted {
	//		buffer := new(bytes.Buffer)
	//		if _, err := io.Copy(buffer, msg.Payload); err != nil {
	//			return true, errDecodeFailed
	//		}
	//		savedMsg := msg
	//		savedMsg.Payload = buffer
	//		sb.pendingMessages.Enqueue(UnhandledMsg{addr: addr, msg: savedMsg})
	//		return true, nil //return nil to avoid shutting down connection during block sync.
	//	}
	//	var data []byte
	//	// todo(youssef): this will be decoded again in core, why?
	//	if err := msg.Decode(&data); err != nil {
	//		return true, errDecodeFailed
	//	}
	//	hash := types.RLPHash(data)
	//	// Mark peer's message
	//	ms, ok := sb.recentMessages.Get(addr)
	//	var m *lru.ARCCache
	//	if ok {
	//		m, _ = ms.(*lru.ARCCache)
	//	} else {
	//		m, _ = lru.NewARC(inmemoryMessages)
	//		sb.recentMessages.Add(addr, m)
	//	}
	//	m.Add(hash, true)
	//	// Mark self known message
	//	if _, ok := sb.knownMessages.Get(hash); ok {
	//		return true, nil
	//	}
	//	sb.knownMessages.Add(hash, true)
	//	sb.postEvent(events.MessageEvent{
	//		Payload: data,
	//		ErrCh:   errCh,
	//	})
	case constants.TendermintMsgProposal:
		fallthrough
	case constants.TendermintMsgVote:
		fallthrough
	case constants.TendermintMsgLightProposal:
		if !sb.coreStarted {
			buffer := new(bytes.Buffer)
			if _, err := io.Copy(buffer, msg.Payload); err != nil {
				sb.logger.Info("msg.Payload copy failed", "err", err)
				return true, errDecodeFailed
			}
			savedMsg := msg
			savedMsg.Payload = buffer
			sb.pendingMessages.Enqueue(UnhandledMsg{addr: addr, msg: savedMsg})
			return true, nil //return nil to avoid shutting down connection during block sync.
		}

		payload, err := io.ReadAll(msg.Payload)
		if err != nil {
			sb.logger.Info("msg.Payload read failed", "err", err)
			return true, errDecodeFailed
		}

		var decodedMessage *message.Message

		//var data message.ConsensusMsg

		//var err error
		//
		if msg.Code == constants.TendermintMsgLightProposal {

			messageToDecode := &message.MessageLightProposal{}
			if err := rlp.DecodeBytes(payload, &messageToDecode); err != nil {
				sb.logger.Info("payload decode failed", "err", err)
				return true, errDecodeFailed
			}
			decodedMessage = messageToDecode.ToMessage()

		} else if msg.Code == constants.TendermintMsgProposal {

			messageToDecode := &message.MessageProposal{}
			if err := rlp.DecodeBytes(payload, &messageToDecode); err != nil {
				sb.logger.Info("payload decode failed", "err", err)
				return true, errDecodeFailed
			}
			decodedMessage = messageToDecode.ToMessage()

			//var decoded message.Proposal
			//
			////var data []byte
			//// todo(youssef): this will be decoded again in core, why?
			//if err := rlp.DecodeBytes(payload, &decoded); err != nil {
			//	return true, errDecodeFailed
			//}
			//data = &decoded
		} else if msg.Code == constants.TendermintMsgVote {

			messageToDecode := &message.MessageVote{}
			if err := rlp.DecodeBytes(payload, &messageToDecode); err != nil {
				sb.logger.Info("payload decode failed", "err", err)
				return true, errDecodeFailed
			}
			decodedMessage = messageToDecode.ToMessage()

			//var decoded message.Vote
			//
			//if err := rlp.DecodeBytes(payload, &decoded); err != nil {
			//	return true, errDecodeFailed
			//}
			//data = &decoded
		}

		sb.logger.Debug("payload", "hex", hex.EncodeToString(payload))

		//if err := rlp.DecodeBytes(payload, &messageToDecode); err != nil {
		//	sb.logger.Info("payload decode failed", "err", err)
		//	return true, errDecodeFailed
		//}

		//data = &decoded

		// todo(maks): do we need re-encode this as RLP for this hash?
		hash := types.RLPHash(payload)

		fmt.Printf("handleMsg hash %x\n", hash)

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
		sb.postEvent(events.NewMessageEvent{
			Payload: payload,
			Message: decodedMessage,
			ErrCh:   errCh,
		})

	case constants.SyncMsg:
		if !sb.coreStarted {
			sb.logger.Debug("Sync message received but core not running")
			return true, nil // we return nil as we don't want to shut down the connection if core is stopped
		}
		sb.logger.Debug("Received sync message", "from", addr)
		sb.postEvent(events.SyncEvent{Addr: addr})
	case constants.AccountabilityMsg:
		if !sb.coreStarted {
			sb.logger.Debug("Accountability Msg received but core not running")
			return true, nil // we return nil as we don't want to shut down the connection if core is stopped
		}
		var data []byte
		if err := msg.Decode(&data); err != nil {
			// this error will freeze peer for 30 seconds by according to dev p2p protocol.
			sb.logger.Info("AccountabilityMsg decode failed", "err", err)
			return true, errDecodeFailed
		}

		// post the off chain accountability msg to the event handler, let the event handler to handle DoS attack vectors.
		sb.logger.Debug("Received Accountability Msg", "from", addr)
		sb.postEvent(events.AccountabilityEvent{Sender: addr, Payload: data, ErrCh: errCh})

	default:
		return false, nil
	}

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
