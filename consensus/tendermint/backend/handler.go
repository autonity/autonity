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
	// errJailed is returned when a consensus message is discarded because the signer is jailed
	ErrJailed    = errors.New("signer is jailed")
	ErrNotFuture = errors.New("message is not future height anymore")
	NetworkCodes = map[uint8]uint64{
		message.ProposalCode:  ProposeNetworkMsg,
		message.PrevoteCode:   PrevoteNetworkMsg,
		message.PrecommitCode: PrecommitNetworkMsg,
	}
	ProposalProcessBg  = metrics.NewRegisteredBufferedGauge("acn/proposal/process", nil, nil)                          // time between round start and proposal sent
	PrevoteProcessBg   = metrics.NewRegisteredBufferedGauge("acn/prevote/process", nil, metrics.GetIntPointer(1024))   // time between round start and proposal receiv
	PrecommitProcessBg = metrics.NewRegisteredBufferedGauge("acn/precommit/process", nil, metrics.GetIntPointer(1024)) // time to verify proposal
	DefaultProcessBg   = metrics.NewRegisteredBufferedGauge("acn/any/process", nil, nil)                               // time to verify proposal

	TotalMessageReceivedBg = metrics.NewRegisteredMeter("acn/handler/message/received", nil)  // total message received
	MessageProcessedBg     = metrics.NewRegisteredMeter("acn/handler/message/processed", nil) // total message processed
)

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
func (sb *Backend) HandleMsg(sender common.Address, msg p2p.Msg, errCh chan<- error) (bool, error) {
	if msg.Code < ProposeNetworkMsg || msg.Code > AccountabilityNetworkMsg {
		return false, nil
	}

	switch msg.Code {
	case ProposeNetworkMsg:
		return handleConsensusMsg[message.Propose](sb, sender, msg, errCh)
	case PrevoteNetworkMsg:
		return handleConsensusMsg[message.Prevote](sb, sender, msg, errCh)
	case PrecommitNetworkMsg:
		return handleConsensusMsg[message.Precommit](sb, sender, msg, errCh)
	case SyncNetworkMsg:
		if !sb.coreRunning.Load() {
			sb.logger.Debug("Sync message received but core not running")
			return true, nil // we return nil as we don't want to shut down the connection if core is stopped
		}
		if sb.IsJailed(sender) {
			sb.logger.Debug("Ignoring sync message from jailed validator", "from", sender)
			return true, ErrJailed
		}
		sb.logger.Debug("Received sync message", "from", sender)
		go sb.Post(events.SyncEvent{Addr: sender})
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
		sb.logger.Debug("Received Accountability Msg", "from", sender)
		go sb.Post(events.AccountabilityEvent{Sender: sender, Payload: data, ErrCh: errCh})
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
	// if the message is for a future height wrt to consensus engine, buffer it
	// it will be re-injected into the handleDecodedMsg function at the right height
	// TODO: Due to a race condition a message that is considered as future could become current,
	// but remain stuck into the future message buffer forever
	currentHeight := sb.core.Height().Uint64()
	if msg.H() > currentHeight {
		sb.logger.Debug("Saving future height consensus message for later", "msgHeight", msg.H(), "coreHeight", currentHeight)
		err := sb.saveFutureMsg(msg, errCh, sender)
		// if we receive `ErrNotFuture`, keep processing the message as it is now a current height message
		if err == nil || !errors.Is(err, ErrNotFuture) {
			return true, err
		}
	}
	// if the height is so old that it is not useful even for accountability, discard it right away. No need to waste resources on this.
	if sb.isHeightExpired(currentHeight, msg.H()) {
		return true, nil
	}
	return sb.handleDecodedMsg(msg, errCh, sender)
}

func (sb *Backend) handleDecodedMsg(msg message.Msg, errCh chan<- error, sender common.Address) (bool, error) {
	committee, err := sb.BlockChain().CommitteeByHeight(msg.H())
	if err != nil {
		// since this is not a future message, we should always have the committee of the height.
		sb.logger.Crit("Missing committee for non-future consensus message", "height", msg.H())
	}

	// assign power and bls signer key
	if err := msg.PreValidate(committee); err != nil {
		return true, err
	}

	// if the sender is jailed, discard its messages
	switch m := msg.(type) {
	case *message.Propose:
		if sb.IsJailed(m.Signer()) {
			sb.logger.Debug("Ignoring proposal from jailed validator", "address", m.Signer())
			return true, ErrJailed
		}
	case *message.Prevote, *message.Precommit:
		vote := m.(message.Vote)
		allJailed := true
		for _, signerIndex := range vote.Signers().FlattenUniq() {
			signer := committee.Members[signerIndex].Address
			if !sb.IsJailed(signer) {
				allJailed = false
				break
			}
		}
		// unless all signers are jailed, we still process aggregates
		if allJailed {
			sb.logger.Debug("Vote message contains only signatures from jailed validators, ignoring message", "signers", vote.Signers().String())
			return true, ErrJailed
		}
	default:
		sb.logger.Crit("Tendermint backend processing unknown message")
	}

	sb.Post(events.UnverifiedMessageEvent{
		Message: msg,
		ErrCh:   errCh,
		Sender:  sender,
		Posted:  time.Now(),
	})
	return true, nil
}

func (sb *Backend) saveFutureMsg(msg message.Msg, errCh chan<- error, sender common.Address) error {
	// create event that will be re-injected in handleDecodedMsg when we reach the correct height
	e := &events.UnverifiedMessageEvent{
		Message: msg,
		ErrCh:   errCh,
		Sender:  sender,
	}
	h := msg.H()

	sb.future.Lock()
	defer sb.future.Unlock()

	// this if deals with a possible race condition where a message is considered future, however core changes height and process the future messages
	// before `msg` is stored into the future message buffer. Without this if `msg` would get stuck into the future message store until we change height again.
	if h <= sb.future.lastProcessedHeight {
		return ErrNotFuture
	}

	if h < sb.future.minHeight {
		sb.future.minHeight = h
	}
	if h > sb.future.maxHeight {
		sb.future.maxHeight = h
	}
	sb.future.messages[h] = append(sb.future.messages[h], e)
	sb.future.size++

	// if needed, drop heights until we are back under the threshold
	for sb.future.size > maxFutureMsgs {
		maxHeightEvs, ok := sb.future.messages[sb.future.maxHeight]
		sb.logger.Debug("deleting excess future height messages", "height", sb.future.maxHeight)
		if ok {
			sb.future.size -= uint64(len(maxHeightEvs))
			// remove messages from knowMessages cache so they can be received again
			go func(evs []*events.UnverifiedMessageEvent) {
				for _, e := range evs {
					sb.knownMessages.Remove(e.Message.Hash())
				}
			}(maxHeightEvs)
			delete(sb.future.messages, sb.future.maxHeight)
		}
		// This value might be different wrt the actual maximum in the map (because of holes in future msg heights)
		// however it is always going to be >= actualMaximum, so it is fine
		sb.future.maxHeight--

		// TODO: might want to remove this once we are sure everything works as intended
		if sb.future.maxHeight < sb.future.minHeight-1 {
			log.Crit("inconsistent state in future message buffer")
		}
	}
	return nil
}

// re-inject future height messages
func (sb *Backend) ProcessFutureMsgs(height uint64) {
	sb.future.Lock()
	defer sb.future.Unlock()

	// shortcircuit if:
	// - we have no future messages
	// - minimum future height is greater than height
	if sb.future.size == 0 || sb.future.minHeight > height {
		return
	}

	// process future messages up to current height
	for h := sb.future.minHeight; h <= height; h++ {
		evs, ok := sb.future.messages[h]
		// there might be holes in heights in the future messages
		if ok {
			sb.logger.Debug("processing future height messages", "height", h, "n", len(sb.future.messages[h]))
			for _, e := range evs {
				sb.handleDecodedMsg(e.Message, e.ErrCh, e.Sender)
				sb.future.size--
			}
			delete(sb.future.messages, h)
		}
	}

	// This value might be different wrt the actual minimum in the map (because of holes in future msg heights)
	// however it is always going to be <= actualMinimum, so it is fine (even though not optimal)
	sb.future.minHeight = height + 1
	sb.future.lastProcessedHeight = height
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
