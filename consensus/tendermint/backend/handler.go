package backend

import (
	"bytes"
	"context"
	"errors"
	"io"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/events"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/p2p"
	lru "github.com/hashicorp/golang-lru"
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
)

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

	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()

	switch msg.Code {
	case ProposeNetworkMsg:
		return handleConsensusMsg[message.Propose](sb, addr, msg, errCh)
	case PrevoteNetworkMsg:
		return handleConsensusMsg[message.Prevote](sb, addr, msg, errCh)
	case PrecommitNetworkMsg:
		return handleConsensusMsg[message.Precommit](sb, addr, msg, errCh)
	case SyncNetworkMsg:
		if !sb.coreStarted {
			sb.logger.Debug("Sync message received but core not running")
			return true, nil // we return nil as we don't want to shut down the connection if core is stopped
		}
		sb.logger.Debug("Received sync message", "from", addr)
		go sb.Post(events.SyncEvent{Addr: addr})
	case AccountabilityNetworkMsg:
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

func handleConsensusMsg[T any, PT interface {
	*T
	message.Msg
}](sb *Backend, sender common.Address, p2pMsg p2p.Msg, errCh chan<- error) (bool, error) {
	if !sb.coreStarted {
		// We copy the message here as it can't be saved directly due
		// to a call to Discard in the eth handler which is going to empty this buffer.
		buffer := bytes.NewBuffer(make([]byte, 0, p2pMsg.Size))
		if _, err := io.Copy(buffer, p2pMsg.Payload); err != nil {
			return true, errDecodeFailed
		}
		savedMsg := p2pMsg
		savedMsg.Payload = bytes.NewReader(buffer.Bytes())
		sb.pendingMessages.Enqueue(UnhandledMsg{addr: sender, msg: savedMsg})
		return true, nil // return nil to avoid shutting down connection during block sync.
	}
	msg := PT(new(T))
	if err := p2pMsg.Decode(msg); err != nil {
		sb.logger.Error("Error decoding consensus message", "err", err)
		return true, err
	}
	// Mark peer's message as known.
	ms, ok := sb.recentMessages.Get(sender)
	var m *lru.ARCCache
	if ok {
		m, _ = ms.(*lru.ARCCache)
	} else {
		m, _ = lru.NewARC(inmemoryMessages)
		sb.recentMessages.Add(sender, m)
	}
	m.Add(msg.Hash(), true)
	// Mark the message known for ourselves
	if _, ok := sb.knownMessages.Get(msg.Hash()); ok {
		return true, nil
	}
	sb.knownMessages.Add(msg.Hash(), true)

	// if the message is for a future height wrt to consensus engine, buffer it
	// it will be re-injected into the handleDecodedMsg function at the right height
	if msg.H() > sb.core.Height().Uint64() {
		sb.logger.Debug("Saving future height consensus message for later", "msgHeight", msg.H(), "coreHeight", sb.core.Height().Uint64())
		sb.saveFutureMsg(msg, errCh)
		return true, nil
	}
	return sb.handleDecodedMsg(msg, errCh)
}

// TODO(lorenzo) do I need generics?
func (sb *Backend) handleDecodedMsg(msg message.Msg, errCh chan<- error) (bool, error) {
	header := sb.BlockChain().GetHeaderByNumber(msg.H() - 1)
	if header == nil {
		// since this is not a future message, we should always have the header of the parent block.
		sb.logger.Crit("Missing parent header for non-future consensus message", "height", msg.H())
	}

	// verify ecdsa signature
	if err := msg.Validate(header.CommitteeMember); err != nil {
		sb.logger.Debug("Failed to verify signature for consensus msg", "hash", msg.Hash())
		return true, err
	}

	// if the sender is jailed, discard its messages
	if sb.IsJailed(msg.Sender()) {
		sb.logger.Debug("ignoring message from jailed validator", "address", msg.Sender())
		// this one is tricky. Ideally yes, we want to disconnect the sender but we can't
		// really assume that all the other committee members have the same view on the
		// jailed validator list before gossip, that is risking then to disconnect honest nodes.
		// This needs to verified though. Returning nil for the time being.
		return true, nil
	}

	// if the message is for current height, post both to tendermint core and FD
	if msg.H() == sb.core.Height().Uint64() {
		go sb.Post(events.MessageEvent{
			Message: msg,
			ErrCh:   errCh,
		})
		return true, nil
	}

	// if a message arrives here, it means it is a valid old height message.
	// this will be picked up only by the FD.
	go sb.Post(events.OldMessageEvent{
		Message: msg,
		ErrCh:   errCh,
	})
	return true, nil
}

// TODO(lorenzo) do I need generics?
func (sb *Backend) saveFutureMsg(msg message.Msg, errCh chan<- error) {
	// create event that will be re-injected in handleDecodedMsg when we reach the correct height
	e := &events.MessageEvent{
		Message: msg,
		ErrCh:   errCh,
	}
	h := msg.H()

	sb.futureLock.Lock()
	defer sb.futureLock.Unlock()

	if h < sb.futureMinHeight {
		sb.futureMinHeight = h
	}
	if h > sb.futureMaxHeight {
		sb.futureMaxHeight = h
	}
	sb.future[h] = append(sb.future[h], e)
	sb.futureSize++

	// if needed, drop heights until we are back under the threshold
	for sb.futureSize > maxFutureMsgs {
		maxHeightEvs, ok := sb.future[sb.futureMaxHeight]
		sb.logger.Debug("deleting excess future height messages", "height", sb.futureMaxHeight)
		if ok {
			sb.futureSize -= uint64(len(maxHeightEvs))
			// remove messages from knowMessages cache so they can be received again
			go func(evs []*events.MessageEvent) {
				for _, e := range evs {
					sb.knownMessages.Remove(e.Message.Hash())
				}
			}(maxHeightEvs)
			delete(sb.future, sb.futureMaxHeight)
		}
		// This value might be different wrt the actual maximum in the map (because of holes in future msg heights)
		// however it is always going to be >= actualMaximum, so it is fine
		sb.futureMaxHeight--

		// TODO(lorenzo) might want to remove this once we are sure everything works as intended
		if sb.futureMaxHeight < sb.futureMinHeight-1 {
			log.Crit("inconsistent state in future message buffer")
		}
	}
}

// TODO(lorenzo) do I need generics?
// re-inject future height messages
func (sb *Backend) ProcessFutureMsgs(height uint64) {
	sb.futureLock.Lock()
	defer sb.futureLock.Unlock()

	// shortcircuit if:
	// - we have no future messages
	// - minimum future height is greater than height
	if sb.futureSize == 0 || sb.futureMinHeight > height {
		return
	}

	// process future messages up to current height
	for h := sb.futureMinHeight; h <= height; h++ {
		evs, ok := sb.future[h]
		// there might be holes in heights in the future messages
		if ok {
			sb.logger.Debug("processing future height messages", "height", h, "n", len(sb.future[h]))
			for _, e := range evs {
				sb.handleDecodedMsg(e.Message, e.ErrCh)
				sb.futureSize--
			}
			delete(sb.future, h)
		}
	}

	// This value might be different wrt the actual minimum in the map (because of holes in future msg heights)
	// however it is always going to be <= actualMinimum, so it is fine (even though not optimal)
	sb.futureMinHeight = height + 1
}

// SetBroadcaster implements consensus.Handler.SetBroadcaster
func (sb *Backend) SetBroadcaster(broadcaster consensus.Broadcaster) {
	sb.Broadcaster = broadcaster
	sb.gossiper.SetBroadcaster(broadcaster)
}

func (sb *Backend) NewChainHead() error {
	sb.coreMu.RLock()
	defer sb.coreMu.RUnlock()
	if !sb.coreStarted {
		return ErrStoppedEngine
	}
	go sb.Post(events.CommitEvent{})
	return nil
}
