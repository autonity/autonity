package backend

import (
	"context"
	"reflect"
	"sync"
	"time"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/events"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/log"
)

func newAggregator(backend interfaces.Backend, logger log.Logger) *aggregator {
	return &aggregator{backend: backend, messages: make(map[common.Hash][]events.UnverifiedMessageEvent), logger: logger}
}

type aggregator struct {
	backend  interfaces.Backend
	messages map[common.Hash][]events.UnverifiedMessageEvent
	sub      *event.TypeMuxSubscription
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	logger   log.Logger
}

func (a *aggregator) start(ctx context.Context) {
	a.logger.Info("Starting the aggregator routine")
	a.sub = a.backend.Subscribe(events.UnverifiedMessageEvent{})
	ctx, a.cancel = context.WithCancel(ctx)
	a.wg.Add(1)
	go a.loop(ctx)
}

func tryDisconnect(errorCh chan<- error, err error) {
	select {
	case errorCh <- err:
	default: // do nothing
	}
}

func (a *aggregator) loop(ctx context.Context) {
	defer a.wg.Done()

	timer := time.NewTimer(20 * time.Millisecond)

loop:
	for {
		select {
		case ev, ok := <-a.sub.Chan():
			if !ok {
				break loop
			}
			event := ev.Data.(events.UnverifiedMessageEvent)
			msg := event.Message
			// if proposal or aggregatedVote, verify right away
			switch msg.(type) {
			case *message.Propose:
				propose := msg.(*message.Propose)
				if err := propose.Validate(); err != nil {
					//TODO(lorenzo) also remove msgs sent by same p2p peer?
					tryDisconnect(event.ErrCh, err)
				}
				go a.backend.Post(events.MessageEvent{
					Message: msg,
					ErrCh:   event.ErrCh,
				})
			case *message.AggregatePrevote, *message.AggregatePrecommit:
				if err := msg.(message.AggregateMsg).Validate(); err != nil {
					//TODO(lorenzo) also remove msgs sent by same p2p peer?
					tryDisconnect(event.ErrCh, err)
				}
				go a.backend.Post(events.MessageEvent{
					Message: msg,
					ErrCh:   event.ErrCh,
				})

			case *message.Prevote, *message.Precommit:
				// batch
				signatureInput := event.Message.SignatureInput()
				a.messages[signatureInput] = append(a.messages[signatureInput], event) //TODO(lorenzo) does this work + optimize allocation
			default:
				a.logger.Crit("unknown message type arrived in aggregator")
			}
		case <-timer.C:
			for hash, batch := range a.messages {
				var pubkeys []blst.PublicKey
				var signatures []blst.Signature
				var messages []message.Msg
				for _, e := range batch {
					m := e.Message
					messages = append(messages, m)
					pubkeys = append(pubkeys, m.SenderKey())
					signatures = append(signatures, m.Signature())
				}
				aggregateSignature := blst.Aggregate(signatures)
				valid := aggregateSignature.FastAggregateVerify(pubkeys, hash)
				if valid {
					// all votes are valid, send to core and FD

					// same batch --> same signature input --> same height
					height := messages[0].H()

					parent := a.backend.BlockChain().GetHeaderByNumber(height - 1)
					if parent == nil {
						// shouldn't happen due to future msgs being buffered before
						a.logger.Crit("Cannot fetch parent header for non-future consensus message")
					}

					var aggregateVote message.Msg
					switch messages[0].(type) {
					case *message.Prevote:
						aggregateVote = message.NewAggregatePrevote(messages, parent)
					case *message.Precommit:
						aggregateVote = message.NewAggregatePrecommit(messages, parent)
					default:
						a.logger.Crit("messages being aggregated are not individual votes", "type", reflect.TypeOf(messages[0]))
					}

					go a.backend.Post(events.MessageEvent{
						Message: aggregateVote,
						ErrCh:   nil, //TODO(lorenzo) do we add an errCh here?
					})
				} else {
					panic("TODO") //TODO(lorenzo) deal with unhappy case
				}
			}
			a.messages = make(map[common.Hash][]events.UnverifiedMessageEvent)
			timer = time.NewTimer(20 * time.Millisecond)
		case <-ctx.Done():
			break loop
		}
	}
}

func (a *aggregator) stop() {
	a.logger.Info("Stopping the aggregator routine")
	a.cancel()
	a.wg.Wait()
}
