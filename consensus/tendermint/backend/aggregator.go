package backend

import (
	"context"
	"sync"
	"time"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/events"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/event"
)

//TODO(lorenzo) add logs to start and stop

func newAggregator(backend interfaces.Backend) *aggregator {
	return &aggregator{backend: backend, messages: make(map[common.Hash][]events.UnverifiedMessageEvent)}
}

type aggregator struct {
	backend  interfaces.Backend
	messages map[common.Hash][]events.UnverifiedMessageEvent
	sub      *event.TypeMuxSubscription
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

func (a *aggregator) start(ctx context.Context) {
	a.sub = a.backend.Subscribe(events.UnverifiedMessageEvent{})
	ctx, a.cancel = context.WithCancel(ctx)
	a.wg.Add(1)
	go a.loop(ctx)
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
			// if proposal verify right away
			if msg.Code() == message.ProposalCode {
				propose := msg.(*message.Propose)
				if propose.Validate() == nil {
					go a.backend.Post(events.MessageEvent{
						Message: msg,
						ErrCh:   event.ErrCh,
					})
				} else {
					panic("TODO: signature verification failed") //TODO(lorenzo) disconnect peer and eventually remove msgs
				}
			} else {
				// batch
				signatureHash := event.Message.SignatureHash()
				a.messages[signatureHash] = append(a.messages[signatureHash], event) //TODO(lorenzo) does this work + optimize allocation
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
					// send to core and FD

					// same batch --> same height
					height := batch[0].Message.H()

					parent := a.backend.BlockChain().GetHeaderByNumber(height - 1)
					if parent == nil {
						// shouldn't happen due to future msgs being buffered before
						panic("TODO")
					}

					var aggregateVote message.Msg
					switch messages[0].(type) {
					case *message.Prevote:
						aggregateVote = message.NewAggregatePrevote(messages, parent)
					case *message.Precommit:
						aggregateVote = message.NewAggregatePrecommit(messages, parent)
					}
					//TODO(lorenzo) is this the bestway?
					err := aggregateVote.PreValidate(parent)
					if err != nil {
						panic(err)
					}

					//TODO fix
					go a.backend.Post(events.MessageEvent{
						Message: aggregateVote,
						ErrCh:   nil, //TODO(lorenzo) fix (custom event type)
						//ErrCh:   errCh,
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
	a.cancel()
	a.wg.Wait()
}
