package backend

import (
	"context"
	"reflect"
	"sort"
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
	return &aggregator{
		backend:   backend,
		messages:  make(map[common.Hash][]events.UnverifiedMessageEvent),
		logger:    logger,
		votesFrom: make(map[common.Address][]common.Hash),
		toIgnore:  make(map[common.Hash]struct{}),
	}
}

type aggregator struct {
	backend  interfaces.Backend
	messages map[common.Hash][]events.UnverifiedMessageEvent
	//TODO(lorenzo) write tests for the ignoring behaviour
	votesFrom map[common.Address][]common.Hash
	toIgnore  map[common.Hash]struct{}
	sub       *event.TypeMuxSubscription
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	logger    log.Logger
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
			p2pSender := event.P2pSender

			a.votesFrom[p2pSender] = append(a.votesFrom[p2pSender], msg.Hash())

			// if proposal or aggregatedVote, verify right away
			switch msg.(type) {
			case *message.Propose:
				propose := msg.(*message.Propose)
				if err := propose.Validate(); err != nil {
					tryDisconnect(event.ErrCh, err)
					for _, hash := range a.votesFrom[p2pSender] {
						a.toIgnore[hash] = struct{}{}
					}
					break
				}
				go a.backend.Post(events.MessageEvent{
					Message: msg,
					ErrCh:   event.ErrCh,
				})
			case *message.AggregatePrevote, *message.AggregatePrecommit:
				if err := msg.(message.AggregateMsg).Validate(); err != nil {
					tryDisconnect(event.ErrCh, err)
					for _, hash := range a.votesFrom[p2pSender] {
						a.toIgnore[hash] = struct{}{}
					}
					break
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
				var publicKeys []blst.PublicKey
				var signatures []blst.Signature
				var messages []message.Msg
				var p2pSenders []common.Address
				var errChs []chan<- error

				for _, e := range batch {
					m := e.Message
					// skip messages to be ignored
					_, ignore := a.toIgnore[m.Hash()]
					if ignore {
						continue
					}
					messages = append(messages, m)
					publicKeys = append(publicKeys, m.SenderKey())
					signatures = append(signatures, m.Signature())
					p2pSenders = append(p2pSenders, e.P2pSender)
					errChs = append(errChs, e.ErrCh)
				}

				aggregateSignature := blst.Aggregate(signatures)
				valid := aggregateSignature.FastAggregateVerify(publicKeys, hash)
				if valid {
					//TODO(lorenzo) wrap in function since used also later
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
					//TODO(lorenzo) write tests that end up here

					// at least one of the signatures is invalid
					invalids, err := blst.FindFastInvalidSignatures(signatures, publicKeys, hash)
					if err != nil {
						//TODO(lorenzo) double check implementation, but I don't think we should end up here
						// also I am not sure we even need an error in the return values
						panic(err)
					}
					// remove invalid messages and sent the rest of the batch
					//TODO(lorenzo) check if function already returns sorted
					sort.Slice(invalids, func(i, j int) bool {
						return invalids[i] < invalids[j]
					})
					validVotes := make([]message.Msg, len(messages)-len(invalids))
					//TODO(lorenzo) should be good but test
					j := 0
					for i, msg := range messages {
						if j < len(invalids) && uint(i) == invalids[j] {
							j++
							continue
						}
						validVotes[i-j] = msg
					}

					// we discarded invalid votes, send the valid ones to core and FD

					// same batch --> same signature input --> same height
					height := validVotes[0].H()

					parent := a.backend.BlockChain().GetHeaderByNumber(height - 1)
					if parent == nil {
						// shouldn't happen due to future msgs being buffered before
						a.logger.Crit("Cannot fetch parent header for non-future consensus message")
					}

					var aggregateVote message.Msg
					switch validVotes[0].(type) {
					case *message.Prevote:
						aggregateVote = message.NewAggregatePrevote(validVotes, parent)
					case *message.Precommit:
						aggregateVote = message.NewAggregatePrecommit(validVotes, parent)
					default:
						a.logger.Crit("messages being aggregated are not individual votes", "type", reflect.TypeOf(validVotes[0]))
					}

					go a.backend.Post(events.MessageEvent{
						Message: aggregateVote,
						ErrCh:   nil, //TODO(lorenzo) do we add an errCh here?
					})

					// disconnect validators who sent us invalid votes at p2p layer and ignore the msgs coming from them
					for _, index := range invalids {
						tryDisconnect(errChs[index], message.ErrBadSignature)
						for _, hash := range a.votesFrom[p2pSenders[index]] {
							a.toIgnore[hash] = struct{}{}
						}
					}
				}
			}
			a.messages = make(map[common.Hash][]events.UnverifiedMessageEvent)
			a.votesFrom = make(map[common.Address][]common.Hash)
			a.toIgnore = make(map[common.Hash]struct{})
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
