package backend

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
	"go.uber.org/mock/gomock"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/fixsizecache"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/events"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/internal/testrand"
	"github.com/autonity/autonity/log"
)

var (
	committee = &types.Committee{
		Members: []types.CommitteeMember{*makeBogusMember(0), *makeBogusMember(1), *makeBogusMember(2), *makeBogusMember(3), *makeBogusMember(4), *makeBogusMember(5), *makeBogusMember(6)},
	}

	quorum = bft.Quorum(committee.TotalVotingPower())
	csize  = committee.Len()
)

func makePropose(chain *core.BlockChain, backend *Backend, r int64, h uint64) *message.Propose {
	// don't care that it has empty proposer seal, we just want to check that aggregator sends it to Core
	epoch, err := chain.LatestEpoch()
	if err != nil {
		panic(err)
	}
	block, err := makeBlockWithoutSeal(chain, backend, chain.CurrentBlock())
	if err != nil {
		panic("cannot create block")
	}
	propose := message.NewPropose(r, h, -1, block, backend.Sign, &epoch.Committee.Members[0])
	return propose
}

func makeBogusPropose(r int64, h uint64, senderIndex uint64) *message.Propose {
	header := &types.Header{Number: new(big.Int).SetUint64(h)}
	propose := message.NewPropose(r, h, -1, types.NewBlockWithHeader(header), testSigner, makeBogusMember(senderIndex))
	return propose
}

func makeBogusMember(index uint64) *types.CommitteeMember {
	return &types.CommitteeMember{Index: index, VotingPower: common.Big1, ConsensusKey: testKey.PublicKey(), ConsensusKeyBytes: testKey.PublicKey().Marshal(), Address: testAddress}
}
func mineOneBlock(t *testing.T, chain *core.BlockChain, backend *Backend) {
	oldHeight := backend.core.Height().Uint64()

	block, err := makeBlock(chain, backend, chain.CurrentBlock())
	_, err = chain.InsertChain(types.Blocks{block})
	require.NoError(t, err)
	err = backend.NewChainHead()
	require.NoError(t, err)

	// wait for start of new height in core
	waitFor(t, func() bool {
		return backend.core.Height().Uint64() == oldHeight+1
	}, 100*time.Millisecond, 1*time.Second, "cannot mine a block")
}

// changes the signer key exploiting the Fake object. Makes it possible to arbitrarily control the outcome of the verification.
func tweakPrevote(prevote *message.Prevote, key blst.PublicKey) *message.Prevote {
	return message.NewFakePrevote(message.Fake{
		FakeValue:          prevote.Value(),
		FakeSigners:        prevote.Signers(),
		FakeRound:          uint64(prevote.R()),
		FakeHeight:         prevote.H(),
		FakeSignatureInput: prevote.SignatureInput(),
		FakeSignature:      prevote.Signature(),
		FakePayload:        prevote.Payload(),
		FakeHash:           prevote.Hash(),
		FakeSignerKey:      key,
	})
}

// waits for the condition function to be true, re-checking based on a ticker duration.
// if timeout expires before the condition is true, the test is considered failed.
// condition should be a non-blocking function
func waitFor(t *testing.T, condition func() bool, tickerDuration time.Duration, timeoutDuration time.Duration, failMessage string) {
	ticker := time.NewTicker(tickerDuration)
	timeout := time.NewTimer(timeoutDuration)

	for {
		select {
		case <-ticker.C:
			if condition() {
				return
			}
		case <-timeout.C:
			t.Fatal("Timeout expired: " + failMessage)
		}
	}
}

func waitForExpects(t *testing.T, ctrl *gomock.Controller) {
	// wait for all EXPECTS to be satisfied before calling `Finish()`
	waitFor(t, func() bool {
		return ctrl.Satisfied()
	}, 10*time.Millisecond, 1*time.Second, "mock EXPECTS() are not satisfied")
	ctrl.Finish()
}

// if the condition becomes true at any point, fail the test
// condition should be non-blocking
// returns a function to be deferred by the main testing function
func failIf(t *testing.T, condition func() (bool, error)) func() {
	wg := new(sync.WaitGroup)
	done := make(chan struct{})
	toBeDeferred := func() {
		close(done)
		wg.Wait()
	}

	wg.Add(1)
	go failIfInner(t, condition, done, wg)
	return toBeDeferred
}

func failIfInner(t *testing.T, condition func() (bool, error), done <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-done:
			return
		default:
			result, err := condition()
			if result {
				t.Error("condition evaluated to true: " + err.Error())
				return
			}
			time.Sleep(1 * time.Millisecond)
		}
	}
}

func makeBogusEvent(msg message.Msg) events.UnverifiedMessageEvent {
	return events.UnverifiedMessageEvent{Message: msg, ErrCh: nil, Sender: common.Address{}, Posted: time.Now()}
}

func TestAggregatorStartAndStop(t *testing.T) {
	_, backend := newBlockChain(1)
	require.NotNil(t, backend.aggregator)
	require.NotNil(t, backend.aggregator.internalCoreCh)
	require.NotNil(t, backend.aggregator.internalFdCh)
	require.NotNil(t, backend.aggregator.internalBacklogCh)
	time.Sleep(1 * time.Second)
	require.NoError(t, backend.Close())
	ctx, cancel := context.WithCancel(context.Background())
	require.NoError(t, backend.Start(ctx))
	time.Sleep(1 * time.Second)
	require.NoError(t, backend.Close())
	cancel()
}

func TestAggregatorMessageHandling(t *testing.T) {
	t.Run("current height, current round proposal should be processed right away", func(t *testing.T) {
		// the chain will not autonomously mine as there is no miner module providing candidate blocks
		chain, backend := newBlockChain(1)
		defer func() { require.NoError(t, backend.Close()) }()

		// don't care that it has empty proposer seal, we just want to check that aggregator sends it to Core
		propose := makePropose(chain, backend, 0, 1)
		errCh := make(chan error)
		ctrl := gomock.NewController(t)
		defer waitForExpects(t, ctrl)

		mc := interfaces.NewMockEventDispatcher(ctrl)
		called := atomic.NewBool(false)
		backend.coreEventDispatcher = mc

		mc.EXPECT().Post(gomock.Cond(func(ev any) bool {
			event := ev.(events.MessageEvent)
			if propose.Hash() == event.Message().Hash() {
				return true
			}
			return false
		})).Do(func(ev any) {
			called.Store(true)
		}).Times(1)

		backend.aggregatorMessageCh <- events.UnverifiedMessageEvent{Message: propose, ErrCh: errCh, Sender: common.Address{}, Posted: time.Now()}
		defer failIf(t, func() (bool, error) {
			select {
			case err := <-errCh:
				return true, fmt.Errorf("error while validating the propose: %w", err)
			default:
				// do nothing
			}
			return false, nil
		})()

		waitFor(t, func() bool {
			return called.Load()
		}, 10*time.Millisecond, 1*time.Second, "proposal was not processed by the aggregator")
	})
	t.Run("current height, future round proposal should be buffered", func(t *testing.T) {
		h := uint64(1)
		r := int64(10)

		ctrl := gomock.NewController(t)
		defer waitForExpects(t, ctrl)

		coreMock := interfaces.NewMockCore(ctrl)
		backendMock := interfaces.NewMockBackend(ctrl)

		backendMock.EXPECT().CommitteeByHeight(gomock.Any()).Return(committee, nil).AnyTimes()
		coreMock.EXPECT().Height().Return(new(big.Int).SetUint64(h)).Times(1)
		coreMock.EXPECT().Round().Return(r - 1).Times(1)

		a := &aggregator{
			messages:       make(map[uint64]map[int64]*RoundInfo),
			messagesFrom:   make(map[common.Address][]common.Hash),
			core:           coreMock,
			backend:        backendMock,
			logger:         log.Root(),
			signerSetCache: newAggregatorCache(),
			internalCoreCh: make(chan events.MessageEventer, 1),
			internalFdCh:   make(chan events.MessageEventer, 1),
		}

		propose := makeBogusPropose(r, h, 0)

		a.handleEvent(makeBogusEvent(propose))

		roundInfo := a.messages[h][r]
		require.Equal(t, 1, len(roundInfo.proposals))
		require.Equal(t, propose.Hash(), roundInfo.proposals[0].Message.Hash())
	})
	t.Run("current height, current round prevote should be buffered", func(t *testing.T) {
		h := uint64(1)
		r := int64(0)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		coreMock := interfaces.NewMockCore(ctrl)
		backendMock := interfaces.NewMockBackend(ctrl)

		backendMock.EXPECT().CommitteeByHeight(gomock.Any()).Return(committee, nil).Times(1)
		coreMock.EXPECT().Height().Return(new(big.Int).SetUint64(h)).Times(1)
		coreMock.EXPECT().Round().Return(r).Times(1)

		a := &aggregator{
			messages:       make(map[uint64]map[int64]*RoundInfo),
			messagesFrom:   make(map[common.Address][]common.Hash),
			core:           coreMock,
			backend:        backendMock,
			logger:         log.Root(),
			signerSetCache: newAggregatorCache(),
		}

		value := common.Hash{0xca, 0xfe}
		prevote := message.NewPrevote(r, h, value, testSigner, testCommitteeMember, committee.Len())

		a.handleEvent(makeBogusEvent(prevote))

		roundInfo := a.messages[h][r]
		require.Equal(t, 1, len(roundInfo.prevotes[value]))
		require.Equal(t, prevote.Hash(), roundInfo.prevotes[value][0].Message.Hash())
	})
	t.Run("current height, current round prevote should be processed by time-based aggregation", func(t *testing.T) {
		committeeSize := 4
		chain, backend := newBlockChain(committeeSize)
		defer func() { require.NoError(t, backend.Close()) }()

		genesis := chain.Genesis()
		genesisCommittee := genesis.Header().Epoch.Committee
		h := uint64(1)
		r := int64(0)

		prevote := message.NewPrevote(r, h, common.Hash{0xca, 0xfe}, backend.Sign, &genesisCommittee.Members[0], committeeSize)

		ctrl := gomock.NewController(t)
		defer waitForExpects(t, ctrl)

		mc := interfaces.NewMockEventDispatcher(ctrl)
		called := atomic.NewBool(false)
		mc.EXPECT().Post(gomock.Cond(func(ev any) bool {
			event := ev.(events.MessageEvent)
			if prevote.Hash() == event.Message().Hash() {
				return true
			}
			return false
		})).Do(func(ev any) {
			called.Store(true)
		}).Times(1)
		backend.coreEventDispatcher = mc
		errCh := make(chan error)

		failIf(t, func() (bool, error) {
			select {
			case err := <-errCh:
				return true, fmt.Errorf("error while validating the prevote: %w", err)
			default:
				// do nothing
			}
			return false, nil
		})()

		backend.aggregatorMessageCh <- events.UnverifiedMessageEvent{Message: prevote, ErrCh: errCh, Sender: genesisCommittee.Members[0].Address, Posted: time.Now()}

		waitFor(t, func() bool {
			return called.Load()
		}, 20*time.Millisecond, 200*time.Millisecond, "prevote was not processed by the time-based aggregation")
	})
	t.Run("current height, future round complex aggregate carrying quorum should trigger processing", func(t *testing.T) {
		committeeSize := 4
		chain, backend := newBlockChain(committeeSize)
		defer func() { require.NoError(t, backend.Close()) }()
		genesis := chain.Genesis()
		genesisCommittee := genesis.Header().Epoch.Committee

		h := uint64(1)
		r := int64(10)

		value := common.Hash{0xca, 0xfe}
		prevote := message.NewPrevote(r, h, value, backend.Sign, &genesisCommittee.Members[0], committeeSize)
		prevote.Signers().Increment(&genesisCommittee.Members[1])
		prevote.Signers().Increment(&genesisCommittee.Members[2])
		prevote.Signers().Increment(&genesisCommittee.Members[3])

		errCh := make(chan error)

		defer failIf(t, func() (bool, error) {
			select {
			case err := <-errCh:
				return true, fmt.Errorf("error while validating the prevote: %w", err)
			default:
				// do nothing
			}
			return false, nil
		})()

		backend.aggregatorMessageCh <- events.UnverifiedMessageEvent{Message: prevote, ErrCh: errCh, Sender: genesisCommittee.Members[0].Address, Posted: time.Now()}

		// core should switch to round 10 if message gets processed by it
		waitFor(t, func() bool {
			return backend.core.Round() == r
		}, 1*time.Millisecond, 30*time.Millisecond, "future round messages did not cause round change in core")
	})
	t.Run("current height, future round prevote should be processed if F voting power is reached", func(t *testing.T) {
		committeeSize := 4
		chain, backend := newBlockChain(committeeSize)
		defer func() { require.NoError(t, backend.Close()) }()
		genesis := chain.Genesis()
		genesisCommittee := genesis.Header().Epoch.Committee

		h := uint64(1)
		r := int64(10)

		ctrl := gomock.NewController(t)
		defer waitForExpects(t, ctrl)

		// send message to the aggregator and wait for time based aggregation to send it to Core
		value := testrand.Hash()
		prevote := message.NewPrevote(r, h, value, backend.Sign, &genesisCommittee.Members[0], committeeSize)
		errCh := make(chan error)

		defer failIf(t, func() (bool, error) {
			select {
			case err := <-errCh:
				return true, fmt.Errorf("error while validating the prevote: %w", err)
			default:
				// do nothing
			}
			return false, nil
		})()

		called := atomic.NewBool(false)

		coreEventDispatcherMock := interfaces.NewMockEventDispatcher(ctrl)
		coreEventDispatcher := backend.coreEventDispatcher
		backend.coreEventDispatcher = coreEventDispatcherMock

		coreEventDispatcherMock.EXPECT().Post(gomock.Any()).Do(func(ev any) {
			if event, ok := ev.(events.MessageEvent); ok {
				if event.Message().Hash() == prevote.Hash() {
					called.Store(true)
				}
				switch event.Message().(type) {
				case *message.Prevote:
					coreEventDispatcher.Post(ev) // re-post to the original dispatcher
				default:
					// do nothing, we are only interested in prevote processing
				}
			}
		}).AnyTimes()

		backend.aggregatorMessageCh <- events.UnverifiedMessageEvent{Message: prevote, ErrCh: errCh, Sender: genesisCommittee.Members[0].Address, Posted: time.Now()}
		waitFor(t, func() bool {
			return called.Load()
		}, 20*time.Millisecond, 200*time.Millisecond, "future round prevote has not been processed by time-based aggregation")
		require.Equal(t, uint64(100), backend.aggregator.signerSetCache.totalPowerForRound(h, r, stepDispatched).Uint64())

		// now send message that will reach quorum (together with the previous msg in Core)
		prevote = tweakPrevote(message.NewPrevote(r, h, value, backend.Sign, &genesisCommittee.Members[1], committeeSize), backend.consensusKey.PublicKey())
		backend.aggregatorMessageCh <- events.UnverifiedMessageEvent{Message: prevote, ErrCh: errCh, Sender: genesisCommittee.Members[0].Address, Posted: time.Now()}

		// core should switch to round 10 if message gets processed by it
		waitFor(t, func() bool {
			return backend.core.Round() == r
		}, 1*time.Millisecond, 30*time.Millisecond, "future round messages did not cause round change in core")
	})
}

// old height messages should be buffered and processed periodically
func TestAggregatorOldHeightMessage(t *testing.T) {
	t.Run("Old height messages are buffered as stale", func(t *testing.T) {
		h := uint64(5) // currentHeight

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		coreMock := interfaces.NewMockCore(ctrl)
		coreMock.EXPECT().Height().Return(new(big.Int).SetUint64(h)).Times(1)

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().CommitteeByHeight(gomock.Any()).Return(committee, nil).AnyTimes()

		a := &aggregator{
			staleMessages:  make(map[common.Hash][]events.UnverifiedMessageEvent),
			messagesFrom:   make(map[common.Address][]common.Hash),
			core:           coreMock,
			backend:        backendMock,
			logger:         log.Root(),
			signerSetCache: newAggregatorCache(),
		}
		prevote := message.NewPrevote(0, h-2, common.Hash{0xca, 0xfe}, testSigner, testCommitteeMember, committee.Len())

		a.handleEvent(makeBogusEvent(prevote))

		events := a.staleMessages[prevote.SignatureInput()]
		require.Equal(t, 1, len(events))
		require.Equal(t, prevote.Hash(), events[0].Message.Hash())
	})
	t.Run("Old height messages are processed by the stale messages time-based aggregation", func(t *testing.T) {
		chain, backend := newBlockChain(1)
		genesis := chain.Genesis()
		afdChan := make(chan events.MessageEventer, 100)
		backend.afdDispatchCh = afdChan

		mineOneBlock(t, chain, backend)

		genesisCommittee := genesis.Header().Epoch.Committee
		prevote := message.NewPrevote(0, 1, common.Hash{0xca, 0xfe}, backend.Sign, &genesisCommittee.Members[0], genesisCommittee.Len())
		errCh := make(chan error)

		defer failIf(t, func() (bool, error) {
			select {
			case <-errCh:
				return true, errors.New("message was not buffered")
			default:
				// do nothing
			}
			return false, nil
		})()

		backend.aggregatorMessageCh <- events.UnverifiedMessageEvent{Message: prevote, ErrCh: errCh, Sender: genesisCommittee.Members[0].Address, Posted: time.Now()}

		// check that old message has been processed by stale messages time-based aggregation
		waitFor(t, func() bool {
			select {
			case ev := <-afdChan:
				if oldMsg, ok := ev.(events.OldMessageEvent); ok {
					return oldMsg.Message().Hash() == prevote.Hash()
				}
			default:
				// do nothing
			}
			return false
		}, 100*time.Millisecond, 4*time.Second, "old height message was not processed")
	})
}

func TestAggregatorSaveMessage(t *testing.T) {
	t.Run("Save proposal", func(t *testing.T) {
		a := &aggregator{
			messages:       make(map[uint64]map[int64]*RoundInfo),
			signerSetCache: newAggregatorCache(),
		}

		r := int64(4)
		h := uint64(2)

		propose := makeBogusPropose(r, h, 0)

		proposeEvent := makeBogusEvent(propose)

		a.saveMessage(proposeEvent)

		roundInfo := a.messages[h][r]

		require.Equal(t, propose.Hash(), roundInfo.proposals[0].Message.Hash())
	})
	t.Run("Save prevote", func(t *testing.T) {
		a := &aggregator{messages: make(map[uint64]map[int64]*RoundInfo)}

		r := int64(5)
		h := uint64(0)

		value := common.Hash{0xca, 0xfe}
		vote := message.NewPrevote(r, h, value, testSigner, testCommitteeMember, 1)

		voteEvent := makeBogusEvent(vote)

		a.saveMessage(voteEvent)

		roundInfo := a.messages[h][r]

		require.Equal(t, vote.Hash(), roundInfo.prevotes[value][0].Message.Hash())
	})
	t.Run("Save precommit", func(t *testing.T) {
		a := &aggregator{messages: make(map[uint64]map[int64]*RoundInfo)}

		r := int64(5)
		h := uint64(0)

		value := common.Hash{0xca, 0xfe}
		vote := message.NewPrecommit(r, h, value, testSigner, testCommitteeMember, 1)

		voteEvent := makeBogusEvent(vote)

		a.saveMessage(voteEvent)

		roundInfo := a.messages[h][r]

		require.Equal(t, vote.Hash(), roundInfo.precommits[value][0].Message.Hash())
	})
	t.Run("Save multiple message (individual and aggregates)", func(t *testing.T) {
		a := &aggregator{messages: make(map[uint64]map[int64]*RoundInfo)}
		committeeSize := 4

		r := int64(4)
		h := uint64(2)

		value := common.Hash{0xca, 0xfe}

		propose := makeBogusPropose(r, h, 0)
		proposeEvent := makeBogusEvent(propose)
		a.saveMessage(proposeEvent)

		roundInfo := a.messages[h][r]
		require.Equal(t, propose.Hash(), roundInfo.proposals[0].Message.Hash())

		// save again the same proposal
		a.saveMessage(proposeEvent)

		roundInfo = a.messages[h][r]
		require.Equal(t, propose.Hash(), roundInfo.proposals[1].Message.Hash())

		// save precommit for same (h,r) from the same guy
		precommit := message.NewPrecommit(r, h, value, testSigner, testCommitteeMember, committeeSize)
		voteEvent := events.UnverifiedMessageEvent{Message: precommit, ErrCh: nil, Sender: common.Address{}, Posted: time.Now()}

		a.saveMessage(voteEvent)

		roundInfo = a.messages[h][r]

		require.Equal(t, precommit.Hash(), roundInfo.precommits[value][0].Message.Hash())

		// save vote for same (h,r,v) from different validator
		prevote := message.NewPrevote(r, h, value, testSigner, makeBogusMember(1), committeeSize)
		voteEvent = events.UnverifiedMessageEvent{Message: prevote, ErrCh: nil, Sender: common.Address{}, Posted: time.Now()}

		a.saveMessage(voteEvent)

		roundInfo = a.messages[h][r]

		require.Equal(t, prevote.Hash(), roundInfo.prevotes[value][0].Message.Hash())

		// save aggregated vote for same (h,r), different value
		otherValue := common.Hash{0x13, 0x37}
		prevote = message.NewPrevote(r, h, otherValue, testSigner, makeBogusMember(0), committeeSize)
		prevote.Signers().Increment(makeBogusMember(1))
		prevote.Signers().Increment(makeBogusMember(2))
		voteEvent = events.UnverifiedMessageEvent{Message: prevote, ErrCh: nil, Sender: common.Address{}, Posted: time.Now()}

		a.saveMessage(voteEvent)

		roundInfo = a.messages[h][r]

		require.Equal(t, prevote.Hash(), roundInfo.prevotes[otherValue][0].Message.Hash())

		// save proposal from index 3

		propose = makeBogusPropose(r, h, 3)
		proposeEvent = events.UnverifiedMessageEvent{Message: propose, ErrCh: nil, Sender: common.Address{}, Posted: time.Now()}
		a.saveMessage(proposeEvent)

		roundInfo = a.messages[h][r]
		require.Equal(t, propose.Hash(), roundInfo.proposals[2].Message.Hash())

		// save aggregated vote for same (h,r), different value
		precommit = message.NewPrecommit(r, h, otherValue, testSigner, makeBogusMember(0), committeeSize)
		precommit.Signers().Increment(makeBogusMember(1))
		precommit.Signers().Increment(makeBogusMember(2))
		precommit.Signers().Increment(makeBogusMember(3))
		voteEvent = events.UnverifiedMessageEvent{Message: precommit, ErrCh: nil, Sender: common.Address{}, Posted: time.Now()}

		a.saveMessage(voteEvent)

		roundInfo = a.messages[h][r]

		require.Equal(t, precommit.Hash(), roundInfo.precommits[otherValue][0].Message.Hash())

		// save aggregated vote for same (h,r), different value
		precommit = message.NewPrecommit(r, h, otherValue, testSigner, makeBogusMember(0), committeeSize)
		precommit.Signers().Increment(makeBogusMember(1))
		voteEvent = events.UnverifiedMessageEvent{Message: precommit, ErrCh: nil, Sender: common.Address{}, Posted: time.Now()}

		a.saveMessage(voteEvent)

		roundInfo = a.messages[h][r]

		require.Equal(t, 2, len(roundInfo.precommits[otherValue]))
		require.Equal(t, precommit.Hash(), roundInfo.precommits[otherValue][1].Message.Hash())

		// save vote for different round
		otherRound := r + 10
		precommit = message.NewPrecommit(otherRound, h, value, testSigner, testCommitteeMember, committeeSize)
		voteEvent = events.UnverifiedMessageEvent{Message: precommit, ErrCh: nil, Sender: common.Address{}, Posted: time.Now()}

		a.saveMessage(voteEvent)

		roundInfo = a.messages[h][otherRound]

		require.Equal(t, precommit.Hash(), roundInfo.precommits[value][0].Message.Hash())

		// save vote for different height
		otherHeight := h + 3
		precommit = message.NewPrecommit(otherRound, otherHeight, value, testSigner, testCommitteeMember, committeeSize)
		voteEvent = events.UnverifiedMessageEvent{Message: precommit, ErrCh: nil, Sender: common.Address{}, Posted: time.Now()}

		a.saveMessage(voteEvent)

		roundInfo = a.messages[otherHeight][otherRound]

		require.Equal(t, precommit.Hash(), roundInfo.precommits[value][0].Message.Hash())
	})
}

func TestAggregatorHandleVote(t *testing.T) {
	quorumMinusOne := new(big.Int).Set(quorum)
	quorumMinusOne.Sub(quorumMinusOne, common.Big1)

	a := &aggregator{
		messages:          make(map[uint64]map[int64]*RoundInfo),
		knownMessages:     fixsizecache.New[common.Hash, bool](numBuckets, numEntries, fixsizecache.HashKey[common.Hash]),
		logger:            log.Root(),
		signerSetCache:    newAggregatorCache(),
		messagesFrom:      make(map[common.Address][]common.Hash),
		internalCoreCh:    make(chan events.MessageEventer, 10),
		internalFdCh:      make(chan events.MessageEventer, 10),
		internalBacklogCh: make(chan events.UnverifiedMessageEvent, 10),
	}

	t.Run("equivocated votes dont trigger processing", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer waitForExpects(t, ctrl)
		backendMock := interfaces.NewMockBackend(ctrl)
		coreMock := interfaces.NewMockCore(ctrl)
		a.core = coreMock
		a.backend = backendMock
		a.signerSetCache = newAggregatorCache()

		r := int64(5)
		h := uint64(10)
		value := testrand.Hash()
		a.signerSetCache.markCommittee(h, committee)

		voteA := message.NewPrevote(r, h, value, testSigner, &committee.Members[0], csize)
		voteA.Signers().Increment(&committee.Members[0])
		voteA.Signers().Increment(&committee.Members[1])
		voteAEvent := makeBogusEvent(voteA)

		voteB := message.NewPrevote(r, h, common.Hash{}, testSigner, &committee.Members[0], csize)
		voteB.Signers().Increment(&committee.Members[0])
		voteB.Signers().Increment(&committee.Members[1])
		voteB.Signers().Increment(&committee.Members[2])
		voteBEvent := makeBogusEvent(voteB)

		backendMock.EXPECT().DispatchToCore(gomock.Any()).Times(0)

		a.signerSetCache.addVote(voteA, stepReceived)
		a.handleVote(voteAEvent, quorum)
		require.Equal(t, 1, len(a.messages[h][r].prevotes[value]))

		a.signerSetCache.addVote(voteB, stepReceived)
		a.handleVote(voteBEvent, quorum)
		require.Equal(t, 1, len(a.messages[h][r].prevotes[common.Hash{}]))
	})
	t.Run("quorum for v triggers processing", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer waitForExpects(t, ctrl)

		backendMock := interfaces.NewMockBackend(ctrl)
		coreMock := interfaces.NewMockCore(ctrl)
		a.core = coreMock
		a.backend = backendMock
		a.signerSetCache = newAggregatorCache()

		r := int64(5)
		h := uint64(0)
		value := testrand.Hash()
		a.signerSetCache.markCommittee(h, committee)

		backendMock.EXPECT().Address().Return(testAddress).AnyTimes()

		singleVote := message.NewPrecommit(r, h, value, testSigner, &committee.Members[0], csize)
		singleVote.Signers().Increment(&committee.Members[0])
		singleVoteEvent := makeBogusEvent(singleVote)

		// no quorum reached, vote should be buffered
		a.signerSetCache.addEvent(singleVoteEvent, stepReceived)
		a.handleVote(singleVoteEvent, quorum)
		require.Equal(t, singleVote.Hash(), a.messages[h][r].precommits[value][0].Message.Hash())

		// simple aggregate with quorum should trigger processing
		vote := message.NewPrecommit(r, h, value, testSigner, &committee.Members[1], csize)
		vote.Signers().Increment(&committee.Members[1])
		vote.Signers().Increment(&committee.Members[2])
		vote.Signers().Increment(&committee.Members[3])
		vote.Signers().Increment(&committee.Members[4])
		voteEvent := makeBogusEvent(vote)

		a.signerSetCache.markCommittee(h, committee)
		a.signerSetCache.addVote(vote, stepReceived)
		a.handleVote(voteEvent, quorum)

		require.Nil(t, a.messages[h][r].precommits[value])
	})
	t.Run("quorum for * triggers processing", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer waitForExpects(t, ctrl)

		backendMock := interfaces.NewMockBackend(ctrl)
		coreMock := interfaces.NewMockCore(ctrl)
		a.core = coreMock
		a.backend = backendMock
		a.signerSetCache = newAggregatorCache()

		r := int64(5)
		h := uint64(0)
		value := testrand.Hash()
		a.signerSetCache.markCommittee(h, committee)

		backendMock.EXPECT().Address().Return(testAddress).AnyTimes()

		voteA := message.NewPrecommit(r, h, value, testSigner, &committee.Members[0], csize)
		voteA.Signers().Increment(&committee.Members[0])
		voteA.Signers().Increment(&committee.Members[1])
		voteA.Signers().Increment(&committee.Members[2])
		voteEventA := makeBogusEvent(voteA)

		a.signerSetCache.addEvent(voteEventA, stepReceived)
		a.handleVote(voteEventA, quorum)

		// quorum for * is not reached, vote should be buffered
		require.Equal(t, voteA.Hash(), a.messages[h][r].precommits[value][0].Message.Hash())

		voteB := message.NewPrecommit(r, h, common.Hash{}, testSigner, &committee.Members[3], csize)
		voteB.Signers().Increment(&committee.Members[3])
		voteB.Signers().Increment(&committee.Members[4])
		voteEventB := makeBogusEvent(voteB)
		// quorum for * is reached, vote should be processed
		a.signerSetCache.addEvent(voteEventB, stepReceived)
		a.handleVote(voteEventB, quorum)
		require.Nil(t, a.messages[h][r].precommits[value])
	})
	t.Run("quorum for v doesnt trigger processing if already dispatched", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer waitForExpects(t, ctrl)
		backendMock := interfaces.NewMockBackend(ctrl)
		coreMock := interfaces.NewMockCore(ctrl)
		a.core = coreMock
		a.backend = backendMock
		a.signerSetCache = newAggregatorCache()

		r := int64(5)
		h := uint64(0)
		v := testrand.Hash()
		a.signerSetCache.markCommittee(h, committee)

		backendMock.EXPECT().Address().Return(testAddress).AnyTimes()
		voteA := message.NewPrecommit(r, h, v, testSigner, &committee.Members[0], csize)
		voteA.Signers().Increment(&committee.Members[0])
		voteA.Signers().Increment(&committee.Members[1])
		voteA.Signers().Increment(&committee.Members[2])
		voteA.Signers().Increment(&committee.Members[3])
		voteA.Signers().Increment(&committee.Members[4])

		voteEventA := makeBogusEvent(voteA)
		// this signals that the event has already been processed for [0-4]
		a.signerSetCache.addEvent(voteEventA, stepReceived)
		a.signerSetCache.addEvent(voteEventA, stepDispatched)

		// this vote has new signers (mainly committee[5]) and quorum
		voteB := message.NewPrecommit(r, h, v, testSigner, &committee.Members[5], csize)
		voteB.Signers().Increment(&committee.Members[5])
		voteB.Signers().Increment(&committee.Members[4])
		voteB.Signers().Increment(&committee.Members[3])
		voteB.Signers().Increment(&committee.Members[2])
		voteB.Signers().Increment(&committee.Members[1])

		a.signerSetCache.addEvent(makeBogusEvent(voteB), stepReceived)
		a.handleVote(makeBogusEvent(voteB), quorum)
		// should be pending
		require.Equal(t, voteB.Hash(), a.messages[h][r].precommits[v][0].Message.Hash())
	})
	t.Run("quorum for * doesnt trigger processing if already dispatched", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer waitForExpects(t, ctrl)
		backendMock := interfaces.NewMockBackend(ctrl)
		coreMock := interfaces.NewMockCore(ctrl)
		a.core = coreMock
		a.backend = backendMock
		a.signerSetCache = newAggregatorCache()

		r := int64(5)
		h := uint64(0)
		v := testrand.Hash()
		a.signerSetCache.markCommittee(h, committee)

		backendMock.EXPECT().Address().Return(testAddress).AnyTimes()
		voteA := message.NewPrecommit(r, h, v, testSigner, &committee.Members[0], csize)
		voteA.Signers().Increment(&committee.Members[0])
		voteA.Signers().Increment(&committee.Members[1])
		voteA.Signers().Increment(&committee.Members[2])
		voteA.Signers().Increment(&committee.Members[3])
		voteA.Signers().Increment(&committee.Members[4])

		voteEventA := makeBogusEvent(voteA)
		a.signerSetCache.addEvent(voteEventA, stepReceived)
		a.signerSetCache.addEvent(voteEventA, stepDispatched)

		voteB := message.NewPrecommit(r, h, common.Hash{}, testSigner, &committee.Members[5], csize)
		voteB.Signers().Increment(&committee.Members[5])
		voteB.Signers().Increment(&committee.Members[4])

		a.signerSetCache.addEvent(makeBogusEvent(voteB), stepReceived)
		a.handleVote(makeBogusEvent(voteB), quorum)

		require.Equal(t, voteB.Hash(), a.messages[h][r].precommits[common.Hash{}][0].Message.Hash())
	})
}

func TestAggregatorProcess(t *testing.T) {
	r := int64(4)
	h := uint64(23)
	var messages []message.Msg
	value := common.Hash{0xca, 0xfe}
	otherValue := common.Hash{0xff, 0xff}
	messages = append(messages, makeBogusPropose(r, h, 0))
	messages = append(messages, makeBogusPropose(r, h, 1))
	messages = append(messages, message.NewPrevote(r, h, value, testSigner, &committee.Members[0], csize))
	messages = append(messages, message.NewPrecommit(r, h, value, testSigner, &committee.Members[0], csize))
	messages = append(messages, message.NewPrecommit(r+2, h, value, testSigner, &committee.Members[0], csize))
	messages = append(messages, message.NewPrevote(r, h+3, value, testSigner, &committee.Members[0], csize))
	messages = append(messages, message.NewPrevote(r, h, otherValue, testSigner, &committee.Members[0], csize))

	t.Run("processProposal, proposal is dispatched to core", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer waitForExpects(t, ctrl)
		backendMock := interfaces.NewMockBackend(ctrl)

		a := &aggregator{
			backend:        backendMock,
			internalFdCh:   make(chan events.MessageEventer, 1),
			internalCoreCh: make(chan events.MessageEventer, 1),
			signerSetCache: newAggregatorCache(),
		}
		a.DispatchCoreEvents(context.Background())
		backendMock.EXPECT().DispatchToCore(gomock.Any()).Times(1)
		backendMock.EXPECT().Address().Return(testAddress).AnyTimes()
		propose := makeBogusPropose(0, 1, 0)
		proposeEvent := makeBogusEvent(propose)
		a.signerSetCache.markCommittee(proposeEvent.Message.H(), committee)
		a.processProposal(proposeEvent, currentHeightEventBuilder)
	})
	t.Run("processRound processes all the messages for a round", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(gomock.Any()).AnyTimes()
		backendMock.EXPECT().DispatchToCore(gomock.Any()).AnyTimes()
		backendMock.EXPECT().Address().Return(testAddress).AnyTimes()

		a := &aggregator{
			messages:          make(map[uint64]map[int64]*RoundInfo),
			messagesFrom:      make(map[common.Address][]common.Hash),
			backend:           backendMock,
			knownMessages:     fixsizecache.New[common.Hash, bool](numBuckets, numEntries, fixsizecache.HashKey[common.Hash]),
			logger:            log.Root(),
			signerSetCache:    newAggregatorCache(),
			internalFdCh:      make(chan events.MessageEventer, 10),
			internalCoreCh:    make(chan events.MessageEventer, 10),
			internalBacklogCh: make(chan events.UnverifiedMessageEvent, 10),
		}

		for _, message := range messages {
			a.saveMessage(makeBogusEvent(message))
			a.signerSetCache.markCommittee(message.H(), committee)
		}

		// NOTE: checking a.messages[h][r].precommits (or prevotes) is not really semantically exact as we are checking the number of different values, rather than the number of actual votes.
		require.Equal(t, 5, len(a.messages[h][r].proposals)+len(a.messages[h][r].prevotes)+len(a.messages[h][r].precommits))

		a.processRound(h, r)

		require.Nil(t, a.messages[h][r])
		require.Equal(t, 1, len(a.messages[h][r+2].precommits)) // different round shouldn't have been touched
	})
	t.Run("processVotes processes all prevotes OR precommits for a round", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(gomock.Any()).AnyTimes()
		backendMock.EXPECT().DispatchToCore(gomock.Any()).AnyTimes()
		backendMock.EXPECT().Address().Return(testAddress).AnyTimes()

		a := &aggregator{
			messages:          make(map[uint64]map[int64]*RoundInfo),
			messagesFrom:      make(map[common.Address][]common.Hash),
			backend:           backendMock,
			knownMessages:     fixsizecache.New[common.Hash, bool](numBuckets, numEntries, fixsizecache.HashKey[common.Hash]),
			logger:            log.Root(),
			signerSetCache:    newAggregatorCache(),
			internalFdCh:      make(chan events.MessageEventer, 10),
			internalCoreCh:    make(chan events.MessageEventer, 10),
			internalBacklogCh: make(chan events.UnverifiedMessageEvent, 10),
		}

		for _, message := range messages {
			a.saveMessage(makeBogusEvent(message))
			a.signerSetCache.markCommittee(message.H(), committee)
		}

		require.Equal(t, 2, len(a.messages[h][r].prevotes))

		a.processVotes(h, r, message.PrevoteCode)

		require.Equal(t, 0, len(a.messages[h][r].prevotes))
	})
	t.Run("processVotesFor processes all prevotes OR precommits for a value for a round", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(gomock.Any()).AnyTimes()
		backendMock.EXPECT().DispatchToCore(gomock.Any()).AnyTimes()
		backendMock.EXPECT().Address().Return(testAddress).AnyTimes()

		a := &aggregator{
			messages:          make(map[uint64]map[int64]*RoundInfo),
			messagesFrom:      make(map[common.Address][]common.Hash),
			backend:           backendMock,
			knownMessages:     fixsizecache.New[common.Hash, bool](numBuckets, numEntries, fixsizecache.HashKey[common.Hash]),
			logger:            log.Root(),
			signerSetCache:    newAggregatorCache(),
			internalFdCh:      make(chan events.MessageEventer, 10),
			internalCoreCh:    make(chan events.MessageEventer, 10),
			internalBacklogCh: make(chan events.UnverifiedMessageEvent, 10),
		}

		for _, message := range messages {
			a.saveMessage(makeBogusEvent(message))
			a.signerSetCache.markCommittee(message.H(), committee)
		}

		require.Equal(t, 1, len(a.messages[h][r].prevotes[value]))

		a.processVotesFor(h, r, message.PrevoteCode, value)

		require.Equal(t, 0, len(a.messages[h][r].prevotes[value]))
		require.Equal(t, 1, len(a.messages[h][r].prevotes[otherValue]))
	})
	t.Run("ProcessBatch posts events when batches are valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer waitForExpects(t, ctrl)

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().DispatchToCore(gomock.Any()).Times(4) // two aggregates
		backendMock.EXPECT().DispatchToFD(gomock.Any()).Times(4)   // two aggregates
		backendMock.EXPECT().Address().Return(testAddress).AnyTimes()

		a := &aggregator{
			backend:           backendMock,
			knownMessages:     fixsizecache.New[common.Hash, bool](numBuckets, numEntries, fixsizecache.HashKey[common.Hash]),
			signerSetCache:    newAggregatorCache(),
			internalFdCh:      make(chan events.MessageEventer),
			internalCoreCh:    make(chan events.MessageEventer),
			internalBacklogCh: make(chan events.UnverifiedMessageEvent, 10),
			logger:            log.Root(),
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		a.DispatchCoreEvents(ctx)
		a.DispatchFaultDetectorEvents(ctx)

		var batches [][]events.UnverifiedMessageEvent

		batches = append(batches, []events.UnverifiedMessageEvent{
			makeBogusEvent(message.NewPrevote(r, h, value, testSigner, &committee.Members[0], csize)),
			makeBogusEvent(message.NewPrevote(r, h, value, testSigner, &committee.Members[3], csize)),
			makeBogusEvent(message.NewPrevote(r, h, value, testSigner, &committee.Members[4], csize)),
			makeBogusEvent(message.NewPrevote(r, h, value, testSigner, &committee.Members[5], csize)),
		})
		batches = append(batches, []events.UnverifiedMessageEvent{
			makeBogusEvent(message.NewPrecommit(r, h, otherValue, testSigner, &committee.Members[0], csize)),
			makeBogusEvent(message.NewPrecommit(r, h, otherValue, testSigner, &committee.Members[3], csize)),
			makeBogusEvent(message.NewPrecommit(r, h, otherValue, testSigner, &committee.Members[4], csize)),
			makeBogusEvent(message.NewPrecommit(r, h, otherValue, testSigner, &committee.Members[5], csize)),
		})
		// NOTE: this will trigger two calls to Post because the votes cannot be merged in a simple aggregate
		aggregate1 := message.AggregatePrecommitsSimple([]message.Vote{message.NewPrecommit(r, h, value, testSigner, &committee.Members[0], csize), message.NewPrecommit(r, h, value, testSigner, &committee.Members[3], csize)})
		aggregate2 := message.AggregatePrecommitsSimple([]message.Vote{message.NewPrecommit(r, h, value, testSigner, &committee.Members[0], csize), message.NewPrecommit(r, h, value, testSigner, &committee.Members[4], csize)})
		batches = append(batches, []events.UnverifiedMessageEvent{
			makeBogusEvent(aggregate1[0]),
			makeBogusEvent(aggregate2[0]),
		})

		a.signerSetCache.markCommittee(h, committee)
		a.processBatches(batches, currentHeightEventBuilder)

		close(a.internalCoreCh)
		close(a.internalFdCh)
		close(a.internalBacklogCh)
	})
	t.Run("ProcessBatch successfully detects and discard invalid signatures", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer waitForExpects(t, ctrl)

		backendMock := interfaces.NewMockBackend(ctrl)

		backendMock.EXPECT().DispatchToFD(gomock.Any()).AnyTimes()
		backendMock.EXPECT().DispatchToCore(gomock.Any()).AnyTimes()
		backendMock.EXPECT().Address().Return(testAddress).AnyTimes()

		a := &aggregator{
			backend:           backendMock,
			knownMessages:     fixsizecache.New[common.Hash, bool](numBuckets, numEntries, fixsizecache.HashKey[common.Hash]),
			logger:            log.Root(),
			signerSetCache:    newAggregatorCache(),
			internalFdCh:      make(chan events.MessageEventer, 10),
			internalCoreCh:    make(chan events.MessageEventer, 10),
			internalBacklogCh: make(chan events.UnverifiedMessageEvent, 10),
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		a.DispatchCoreEvents(ctx)
		a.DispatchFaultDetectorEvents(ctx)
		a.signerSetCache.markCommittee(h, committee)

		var batches [][]events.UnverifiedMessageEvent

		// committee[3] and committee[6] are sending invalid sigs

		batches = append(batches, []events.UnverifiedMessageEvent{
			makeBogusEvent(message.NewPrevote(r, h, value, testInvalidSigner, &committee.Members[3], csize)), // INVALID
		})
		// message is correctly signed only by committee[3], but signers include also committee[0]
		invalidSignersPrevote := message.NewPrevote(r, h, value, testSigner, &committee.Members[3], csize)
		invalidSignersPrevote.Signers().Increment(&committee.Members[0])
		aggKey, err := blst.AggregatePublicKeys([]blst.PublicKey{committee.Members[3].ConsensusKey, committee.Members[0].ConsensusKey})
		require.NoError(t, err)
		invalidSignersPrevote = tweakPrevote(invalidSignersPrevote, aggKey)
		batches = append(batches, []events.UnverifiedMessageEvent{
			makeBogusEvent(message.NewPrevote(r, h, value, testSigner, &committee.Members[0], csize)),
			makeBogusEvent(message.NewPrevote(r, h, value, testInvalidSigner, &committee.Members[3], csize)), // INVALID
			makeBogusEvent(invalidSignersPrevote),                                                            // INVALID
			makeBogusEvent(message.NewPrevote(r, h, value, testSigner, &committee.Members[4], csize)),
			makeBogusEvent(message.NewPrevote(r, h, value, testSigner, &committee.Members[5], csize)),
		})
		batches = append(batches, []events.UnverifiedMessageEvent{
			makeBogusEvent(message.NewPrecommit(r, h, otherValue, testSigner, &committee.Members[0], csize)),
			makeBogusEvent(message.NewPrecommit(r, h, otherValue, testInvalidSigner, &committee.Members[3], csize)), //INVALID
			makeBogusEvent(message.NewPrecommit(r, h, otherValue, testSigner, &committee.Members[4], csize)),
			makeBogusEvent(message.NewPrecommit(r, h, otherValue, testSigner, &committee.Members[5], csize)),
			makeBogusEvent(message.NewPrecommit(r, h, otherValue, testInvalidSigner, &committee.Members[6], csize)), //INVALID
		})
		// NOTE: this will trigger two calls to Post because the votes cannot be merged in a simple aggregate
		aggregate1 := message.AggregatePrecommitsSimple([]message.Vote{message.NewPrecommit(r, h, value, testSigner, &committee.Members[0], csize), message.NewPrecommit(r, h, value, testSigner, &committee.Members[4], csize)})
		aggregate2 := message.AggregatePrecommitsSimple([]message.Vote{message.NewPrecommit(r, h, value, testSigner, &committee.Members[0], csize), message.NewPrecommit(r, h, value, testInvalidSigner, &committee.Members[3], csize)}) // INVALID
		batches = append(batches, []events.UnverifiedMessageEvent{
			makeBogusEvent(aggregate1[0]),
			makeBogusEvent(aggregate2[0]), //INVALID
			makeBogusEvent(message.NewPrecommit(r, h, value, testInvalidSigner, &committee.Members[6], csize)), // INVALID
			makeBogusEvent(message.NewPrecommit(r, h, value, testSigner, &committee.Members[5], csize)),
		})

		a.processBatches(batches, func(m message.Msg, ev events.UnverifiedMessageEvent) interface{} {
			vote, ok := m.(message.Vote)
			require.True(t, ok)
			if vote.Signers().Contains(3) || vote.Signers().Contains(6) {
				t.Fatalf("Invalid message has been posted")
			}
			return currentHeightEventBuilder(m, ev)
		})
	})
}

func TestAggregatorFullFlow(t *testing.T) {
	// Previous tests investigate individual aggregator functions, however, now that we
	// are relying on caching we need to make sure that messages are being properly filtered

	t.Run("aggregator should discard messages that contain redundant information", func(t *testing.T) {
		ctrl, a, backendMock, coreMock, aggregatorMsgChan := setupTestAggregator(t)
		defer waitForExpects(t, ctrl)

		h := uint64(1)
		r := int64(5)

		backendMock.EXPECT().DispatchToCore(gomock.Any()).MaxTimes(1)
		backendMock.EXPECT().DispatchToFD(gomock.Any()).MaxTimes(1)
		coreMock.EXPECT().Height().Return(big.NewInt(int64(h))).AnyTimes()
		coreMock.EXPECT().Round().Return(r).AnyTimes()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		a.start(ctx)

		// create a message
		mc := fromCommittee(committee)
		backendMock.EXPECT().CommitteeByHeight(gomock.Any()).Return(mockCommittee(mc).ToCommittee(), nil).AnyTimes()

		value := testrand.Hash()
		event := newSignedTestMsg(t, 1, 5, value, message.PrevoteCode, mc, []int{0, 1, 2})

		aggregatorMsgChan <- event

		// aggregator should process the message, and buffer it
		waitFor(t, func() bool {
			return a.messages[h] != nil && a.messages[h][r] != nil && len(a.messages[h][r].prevotes) == 1
		}, 10*time.Millisecond, 100*time.Millisecond, "Aggregator should have buffered the message")

		redundantEvent := newSignedTestMsg(t, 1, 5, value, message.PrevoteCode, mc, []int{0, 1})

		aggregatorMsgChan <- redundantEvent

		// aggregator should not process the message, as it is redundant
		waitFor(t, func() bool {
			filtered := a.signerSetCache.filtered != nil &&
				a.signerSetCache.filtered[h] != nil &&
				len(a.signerSetCache.filtered[h][filteredCacheKey{
					code:  message.PrevoteCode,
					round: r,
					value: value,
				}]) == 1
			notBuffered := a.messages[h] == nil || a.messages[h][r] == nil || len(a.messages[h][r].prevotes) == 1
			return filtered && notBuffered
		}, 10*time.Millisecond, 100*time.Millisecond, "Aggregator should not have processed the redundant message")
	})

	t.Run("aggregator should dispatch aggregated messages to core once quorum is reached", func(t *testing.T) {
		ctrl, a, backendMock, coreMock, aggregatorMsgChan := setupTestAggregator(t)
		defer waitForExpects(t, ctrl)

		h := uint64(1)
		r := int64(5)

		backendMock.EXPECT().DispatchToFD(gomock.Any()).MaxTimes(1)
		coreMock.EXPECT().Height().Return(big.NewInt(int64(h))).AnyTimes()
		coreMock.EXPECT().Round().Return(r).AnyTimes()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		a.start(ctx)

		// create a message
		mc := fromCommittee(committee)
		backendMock.EXPECT().Address().Return(mc[0].Address).AnyTimes()
		backendMock.EXPECT().CommitteeByHeight(gomock.Any()).Return(mockCommittee(mc).ToCommittee(), nil).AnyTimes()

		value := testrand.Hash()
		event := newSignedTestMsg(t, 1, 5, value, message.PrevoteCode, mc, []int{0, 1, 2})

		aggregatorMsgChan <- event

		// aggregator should process the message, and buffer it
		waitFor(t, func() bool {
			return a.messages[h] != nil && a.messages[h][r] != nil && len(a.messages[h][r].prevotes) == 1
		}, 10*time.Millisecond, 100*time.Millisecond, "Aggregator should have buffered the message")

		// send more messages to reach quorum
		event2 := newSignedTestMsg(t, 1, 5, value, message.PrevoteCode, mc, []int{3, 4, 5})
		passed := atomic.NewBool(false)
		backendMock.EXPECT().DispatchToCore(gomock.Cond(func(ev any) bool {
			msg, ok := ev.(events.MessageEventer)
			if !ok {
				return false
			}
			switch msg.Message().Code() {
			case message.PrevoteCode:
				correct := msg.Message().(message.Vote).Signers().Len() == 6 && msg.Message().Value() == value
				if correct {
					passed.Store(true)
				}
				return correct
			default:
				return false
			}
		})).Times(1)
		aggregatorMsgChan <- event2
		waitFor(t, func() bool {
			return passed.Load()
		}, 10*time.Millisecond, 100*time.Millisecond, "Aggregator should have dispatched the aggregated message to core")
	})

	t.Run("aggregator should reinject messages when an invalid signature is detected", func(t *testing.T) {
		ctrl, a, backendMock, coreMock, aggregatorMsgChan := setupTestAggregator(t)
		defer waitForExpects(t, ctrl)

		h := uint64(1)
		r := int64(5)

		backendMock.EXPECT().DispatchToFD(gomock.Any()).MaxTimes(1)
		coreMock.EXPECT().Height().Return(big.NewInt(int64(h))).AnyTimes()
		coreMock.EXPECT().Round().Return(r).AnyTimes()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		a.start(ctx)

		mc := fromCommittee(committee)
		backendMock.EXPECT().Address().Return(mc[0].Address).AnyTimes()
		backendMock.EXPECT().CommitteeByHeight(gomock.Any()).Return(mockCommittee(mc).ToCommittee(), nil).AnyTimes()

		value := testrand.Hash()
		badlySignedEvent := newSignedTestMsg(t, 1, 5, value, message.PrevoteCode, mc, []int{0, 1, 2, 3})
		badKey, err := blst.RandKey()
		require.NoError(t, err)
		badlySignedEvent.Message = tweakPrevote(badlySignedEvent.Message.(*message.Prevote), badKey.PublicKey())
		redundantEvent := newSignedTestMsg(t, 1, 5, value, message.PrevoteCode, mc, []int{0, 1})
		anotherRedundantEvent := newSignedTestMsg(t, 1, 5, value, message.PrevoteCode, mc, []int{2})

		errChan := make(chan error, 1)
		badlySignedEvent.ErrCh = errChan

		// send all events to aggregator
		aggregatorMsgChan <- badlySignedEvent
		aggregatorMsgChan <- redundantEvent
		aggregatorMsgChan <- anotherRedundantEvent

		waitFor(t, func() bool {
			return a.messages[h] != nil && a.messages[h][r] != nil && len(a.messages[h][r].prevotes) == 1
		}, 10*time.Millisecond, 100*time.Millisecond, "Aggregator should have buffered the message")

		// aggregator should have cached the filtered messages after the bad event
		require.NotNil(t, a.signerSetCache.filtered[h])
		require.Equal(t, 2, len(a.signerSetCache.filtered[h][filteredCacheKey{code: message.PrevoteCode, round: r, value: value}]))

		waitFor(t, func() bool {
			select {
			case <-errChan:
				return true
			default:
				return false
			}
		}, 10*time.Millisecond, 2*aggregationPeriod, "aggregator should process within the aggregation period")

		// aggregator should reinject the redundant messages
		passed := atomic.NewBool(false)
		backendMock.EXPECT().DispatchToCore(gomock.Cond(func(ev any) bool {
			msg, ok := ev.(events.MessageEventer)
			if !ok {
				return false
			}
			switch msg.Message().Code() {
			case message.PrevoteCode:
				correct := msg.Message().(message.Vote).Signers().Len() == 3 && msg.Message().Value() == value
				if correct {
					passed.Store(true)
				}
				return correct
			default:
				return false
			}
		})).Times(1)

		waitFor(t, func() bool {
			return passed.Load()
		}, 10*time.Millisecond, 2*aggregationPeriod, "Aggregator should have reinjected the filtered messages to core")
	})
}

// if we detect an invalid signatures:
// 1. peer is disconnected and suspended
// 2. all the previously buffered messages we received from him are ignored
// This protects us from DoS
func TestAggregatorDosProtection(t *testing.T) {
	h := uint64(23)
	r := int64(1)
	value := common.Hash{0xca, 0xfe}

	ctrl := gomock.NewController(t)
	defer waitForExpects(t, ctrl)

	backendMock := interfaces.NewMockBackend(ctrl)
	backendMock.EXPECT().DispatchToFD(gomock.Any()).AnyTimes()
	backendMock.EXPECT().Address().Return(testAddress).AnyTimes()

	a := &aggregator{
		backend:           backendMock,
		knownMessages:     fixsizecache.New[common.Hash, bool](numBuckets, numEntries, fixsizecache.HashKey[common.Hash]),
		logger:            log.Root(),
		messagesFrom:      make(map[common.Address][]common.Hash),
		messages:          make(map[uint64]map[int64]*RoundInfo),
		toIgnore:          make(map[common.Hash]struct{}),
		internalFdCh:      make(chan events.MessageEventer, 10),
		internalCoreCh:    make(chan events.MessageEventer, 10),
		internalBacklogCh: make(chan events.UnverifiedMessageEvent, 10),
		signerSetCache:    newAggregatorCache(),
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	a.DispatchCoreEvents(ctx)
	a.DispatchFaultDetectorEvents(ctx)

	// suppose committee[0] is sending invalid sigs

	zeroAddress := common.Address{0xff, 0xff, 0xff}
	var votesFromZero []message.Msg
	// mix of valid and invalid votes, from different validators
	votesFromZero = append(
		votesFromZero,
		message.NewPrevote(r+1, h, value, testSigner, &committee.Members[0], csize),
		message.NewPrevote(r, h+3, value, testSigner, &committee.Members[1], csize),
		message.NewPrevote(r, h, value, testInvalidSigner, &committee.Members[2], csize),
	)

	for _, vote := range votesFromZero {
		a.messagesFrom[zeroAddress] = append(a.messagesFrom[zeroAddress], vote.Hash())
		a.saveMessage(events.UnverifiedMessageEvent{Message: vote, ErrCh: nil, Sender: zeroAddress, Posted: time.Now()})
		a.signerSetCache.markCommittee(vote.H(), committee)
	}

	// add another vote, coming from an honest validator (at p2p layer)
	vote := message.NewPrevote(r, h, value, testSigner, &committee.Members[1], csize)
	a.signerSetCache.markCommittee(vote.H(), committee)
	a.saveMessage(makeBogusEvent(vote))
	a.messagesFrom[common.Address{}] = append(a.messagesFrom[common.Address{}], vote.Hash())

	backendMock.EXPECT().DispatchToCore(gomock.Cond(func(ev any) bool {
		msg, ok := ev.(events.MessageEventer)
		if !ok {
			return false
		}
		return msg.Message().Hash() == vote.Hash()
	})).Times(1)

	// mark msg from 0 as invalid
	errCh := make(chan error, 1)
	a.handleInvalidMessage(errCh, message.ErrBadSignature, zeroAddress)
	require.Equal(t, message.ErrBadSignature, <-errCh)
	require.Equal(t, len(votesFromZero), len(a.toIgnore))

	// all these calls should yield a single call to Post, since only one msg came from an honest validator at p2p layer
	a.processRound(h, r+1)
	a.processRound(h+3, r)
	a.processRound(h, r)
}

func newSignedTestMsg(
	t *testing.T,
	h uint64,
	r int64,
	value common.Hash,
	c uint8,
	mc []mockCommitteeMember,
	signers []int,
) events.UnverifiedMessageEvent {
	var msg message.Msg
	switch c {
	case message.ProposalCode:
		if len(signers) != 1 {
			t.Fatalf("proposal must have exactly one signer")
		}
		header := &types.Header{Number: new(big.Int).SetUint64(h)}
		msg = message.NewPropose(r, h, -1, types.NewBlockWithHeader(header), mc[signers[0]].Sign, &mc[signers[0]].CommitteeMember)
	case message.PrevoteCode:
		var msgs []message.Vote
		for _, s := range signers {
			msgs = append(msgs, message.NewPrevote(r, h, value, mc[s].Sign, &mc[s].CommitteeMember, csize))
		}
		msg = message.AggregatePrevotes(msgs)
	case message.PrecommitCode:
		var msgs []message.Vote
		for _, s := range signers {
			msgs = append(msgs, message.NewPrecommit(r, h, value, mc[s].Sign, &mc[s].CommitteeMember, csize))
		}
		msg = message.AggregatePrecommits(msgs)
	default:
		t.Fatalf("unknown message code %d", c)
	}
	return makeBogusEvent(msg)
}

func setupTestAggregator(t *testing.T) (
	*gomock.Controller,
	*aggregator,
	*interfaces.MockBackend,
	*interfaces.MockCore,
	chan events.UnverifiedMessageEvent,
) {
	ctrl := gomock.NewController(t)
	a := &aggregator{
		messages:          make(map[uint64]map[int64]*RoundInfo),
		messagesFrom:      make(map[common.Address][]common.Hash),
		knownMessages:     fixsizecache.New[common.Hash, bool](numBuckets, numEntries, fixsizecache.HashKey[common.Hash]),
		logger:            log.Root(),
		internalFdCh:      make(chan events.MessageEventer, 10),
		internalCoreCh:    make(chan events.MessageEventer, 10),
		internalBacklogCh: make(chan events.UnverifiedMessageEvent, 10),
		signerSetCache:    newAggregatorCache(),
		toIgnore:          make(map[common.Hash]struct{}),
	}
	aggregatorMsgChan := make(chan events.UnverifiedMessageEvent, 10)

	backendMock := interfaces.NewMockBackend(ctrl)
	backendMock.EXPECT().MessageCh().Return(aggregatorMsgChan).AnyTimes()
	coreMock := interfaces.NewMockCore(ctrl)
	a.core = coreMock
	a.backend = backendMock

	return ctrl, a, backendMock, coreMock, aggregatorMsgChan
}
