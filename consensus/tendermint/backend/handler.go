package backend

import (
	"bytes"
	"context"
	"errors"
	"io"
	"time"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/events"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/metrics"
	"github.com/autonity/autonity/p2p"
)

const (
	ProposeNetworkMsg        uint64 = 0x11
	PrevoteNetworkMsg        uint64 = 0x12
	PrecommitNetworkMsg      uint64 = 0x13
	SyncNetworkMsg           uint64 = 0x14
	AccountabilityNetworkMsg uint64 = 0x15
)

type UnhandledMsg struct {
	addr common.Address
	msg  p2p.Msg
}

var (
	// errDecodeFailed is returned when decode message fails
	errDecodeFailed = errors.New("fail to decode tendermint message")
	NetworkCodes    = map[uint8]uint64{
		message.ProposalCode:  ProposeNetworkMsg,
		message.PrevoteCode:   PrevoteNetworkMsg,
		message.PrecommitCode: PrecommitNetworkMsg,
	}
	ProposalProcessBg  = metrics.NewRegisteredBufferedGauge("acn/proposal/process", nil, nil)                  // time between round start and proposal sent
	PrevoteProcessBg   = metrics.NewRegisteredBufferedGauge("acn/prevote/process", nil, getIntPointer(1024))   // time between round start and proposal receiv
	PrecommitProcessBg = metrics.NewRegisteredBufferedGauge("acn/precommit/process", nil, getIntPointer(1024)) // time to verify proposal
	DefaultProcessBg   = metrics.NewRegisteredBufferedGauge("acn/any/process", nil, nil)                       // time to verify proposal

	TotalMessageReceivedBg = metrics.NewRegisteredMeter("acn/handler/message/received", nil)  // total message received
	MessageProcessedBg     = metrics.NewRegisteredMeter("acn/handler/message/processed", nil) // total message processed
)

func getIntPointer(val int) *int {
	return &val
}

func getProcessMetric(msgCode uint64) metrics.BufferedGauge {
	switch msgCode {
	case 0x11:
		return ProposalProcessBg
	case 0x12:
		return PrevoteProcessBg
	case 0x13:
		return PrecommitProcessBg
	}
	return DefaultProcessBg
}

// Protocol implements consensus.Handler.Protocol
func (sb *Backend) Protocol() (protocolName string, extraMsgCodes uint64) {
	return "tendermint", 5 //nolint
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
	if msg.Code < ProposeNetworkMsg || msg.Code > AccountabilityNetworkMsg {
		return false, nil
	}

	switch msg.Code {
	case ProposeNetworkMsg:
		return handleConsensusMsg[message.Propose](sb, addr, msg, errCh)
	case PrevoteNetworkMsg:
		return handleConsensusMsg[message.Prevote](sb, addr, msg, errCh)
	case PrecommitNetworkMsg:
		return handleConsensusMsg[message.Precommit](sb, addr, msg, errCh)
	case SyncNetworkMsg:
		if !sb.coreRunning.Load() {
			sb.logger.Debug("Sync message received but core not running")
			return true, nil // we return nil as we don't want to shut down the connection if core is stopped
		}
		sb.logger.Debug("Received sync message", "from", addr)
		go sb.Post(events.SyncEvent{Addr: addr})
	case AccountabilityNetworkMsg:
		if !sb.coreRunning.Load() {
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

func handleConsensusMsg[T any, PT interface {
	*T
	message.Msg
}](sb *Backend, sender common.Address, p2pMsg p2p.Msg, errCh chan<- error) (bool, error) {
	// we type cast it to byte.Reader because that's the only reader
	// type we expect here
	bReader := p2pMsg.Payload.(*bytes.Reader)
	hash, err := crypto.HashFromReader(bReader)
	if err != nil {
		log.Error("Failed to hash payload", "error", err)
		return true, err
	}
	TotalMessageReceivedBg.Mark(1)
	if sb.knownMessages.Contains(hash) {
		return true, nil
	}
	MessageProcessedBg.Mark(1)
	bReader.Seek(0, io.SeekStart)
	p2pMsg.Payload = bReader
	if !sb.coreRunning.Load() {
		sb.pendingMessages.Enqueue(UnhandledMsg{addr: sender, msg: p2pMsg})
		return true, nil // return nil to avoid shutting down connection during block sync.
	}

	if metrics.Enabled {
		defer func(start time.Time) {
			getProcessMetric(p2pMsg.Code).Add(time.Since(start).Nanoseconds())
		}(time.Now())
	}

	// Mark peer's message as known.
	peer, ok := sb.Broadcaster.FindPeer(sender)
	if !ok {
		sb.logger.Error("message received from unknown peer", "sender", sender)
		return false, nil
	}
	if !peer.Cache().Contains(hash) {
		peer.Cache().Add(hash, true)
	}

	sb.knownMessages.Add(hash, true)
	msg := PT(new(T))
	if err := p2pMsg.Decode(msg); err != nil {
		sb.logger.Error("Error decoding consensus message", "err", err)
		return true, err
	}
	sb.Post(events.MessageEvent{
		Message: msg,
		ErrCh:   errCh,
	})
	return true, nil
}

// SetBroadcaster implements consensus.Handler.SetBroadcaster
func (sb *Backend) SetBroadcaster(broadcaster consensus.Broadcaster) {
	sb.Broadcaster = broadcaster
	sb.gossiper.SetBroadcaster(broadcaster)
}

// SetEnqueuer implements consensus.Handler.SetEnqueuer
func (sb *Backend) SetEnqueuer(enqueuer consensus.Enqueuer) {
	sb.Enqueuer = enqueuer
}

func (sb *Backend) NewChainHead() error {
	if !sb.coreRunning.Load() {
		return ErrStoppedEngine
	}
	go sb.Post(events.CommitEvent{})
	return nil
}
