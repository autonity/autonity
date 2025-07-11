package backend

import (
	"errors"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
	"go.uber.org/mock/gomock"

	"github.com/autonity/autonity/accounts/abi/bind/backends"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/fixsizecache"
	"github.com/autonity/autonity/consensus/ethash"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/events"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/rawdb"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/internal/testrand"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/params"
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

func TestAggregatorMessageHandling(t *testing.T) {
	t.Run("current height, current round proposal should be processed right away", func(t *testing.T) {
		// the chain will not autonomously mine as there is no miner module providing candidate blocks
		chain, backend := newBlockChain(1)

		// don't care that it has empty proposer seal, we just want to check that aggregator sends it to Core
		propose := makePropose(chain, backend, 0, 1)
		errCh := make(chan error)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mc := interfaces.NewMockEventDispatcher(ctrl)
		called := atomic.NewBool(false)
		mc.EXPECT().Post(gomock.Cond(func(ev any) bool {
			event := ev.(events.MessageEvent)
			if propose.Hash() == event.Message().Hash() {
				return true
			}
			return false
		})).Do(func(ev any) {
			called.Store(true)
		}).Times(1)
		backend.coreEventDispatcher = mc

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
		}, time.Millisecond, time.Second, "proposal was not processed by the aggregator")
	})
	t.Run("current height, future round proposal should be buffered", func(t *testing.T) {
		h := uint64(1)
		r := int64(10)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		coreMock := interfaces.NewMockCore(ctrl)
		backendMock := interfaces.NewMockBackend(ctrl)
		chain := newTestBlockchain()

		backendMock.EXPECT().BlockChain().Return(chain).AnyTimes()
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

		close(a.internalCoreCh)
		close(a.internalFdCh)
	})
	t.Run("current height, current round prevote should be buffered", func(t *testing.T) {
		h := uint64(1)
		r := int64(0)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		coreMock := interfaces.NewMockCore(ctrl)
		backendMock := interfaces.NewMockBackend(ctrl)
		chain := newTestBlockchain()

		backendMock.EXPECT().BlockChain().Return(chain).AnyTimes()
		coreMock.EXPECT().Height().Return(new(big.Int).SetUint64(h)).Times(2)
		coreMock.EXPECT().Round().Return(r).Times(2)

		a := &aggregator{
			messages:       make(map[uint64]map[int64]*RoundInfo),
			messagesFrom:   make(map[common.Address][]common.Hash),
			core:           coreMock,
			backend:        backendMock,
			logger:         log.Root(),
			signerSetCache: newAggregatorCache(),
		}

		value := common.Hash{0xca, 0xfe}
		prevote := message.NewPrevote(r, h, value, testSigner, testCommitteeMember, 1)

		a.handleEvent(makeBogusEvent(prevote))

		roundInfo := a.messages[h][r]
		require.Equal(t, 1, len(roundInfo.prevotes[value]))
		require.Equal(t, prevote.Hash(), roundInfo.prevotes[value][0].Message.Hash())
	})
	t.Run("current height, current round prevote should be processed by time-based aggregation", func(t *testing.T) {
		committeeSize := 4
		chain, backend := newBlockChain(committeeSize)
		genesis := chain.Genesis()
		genesisCommittee := genesis.Header().Epoch.Committee
		h := uint64(1)
		r := int64(0)

		prevote := message.NewPrevote(r, h, common.Hash{0xca, 0xfe}, backend.Sign, &genesisCommittee.Members[0], committeeSize)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
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

		backend.aggregatorMessageCh <- events.UnverifiedMessageEvent{Message: prevote, ErrCh: errCh, Sender: genesisCommittee.Members[0].Address, Posted: time.Now()}

		waitFor(t, func() bool {
			return called.Load()
		}, 20*time.Millisecond, 200*time.Millisecond, "prevote was not processed by the time-based aggregation")
	})

	t.Run("current height, future round prevote should be processed if F voting power is reached", func(t *testing.T) {
		committeeSize := 4
		chain, backend := newBlockChain(committeeSize)
		genesis := chain.Genesis()
		genesisCommittee := genesis.Header().Epoch.Committee

		h := uint64(1)
		r := int64(10)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// send message to the aggregator and wait for time based aggregation to send it to Core
		value := common.Hash{0xca, 0xfe}
		prevote := message.NewPrevote(r, h, value, backend.Sign, &genesisCommittee.Members[0], committeeSize)
		errCh := make(chan error)

		passed := make(chan struct{}, 1)
		coreEventDispatcherMock := interfaces.NewMockEventDispatcher(ctrl)
		coreEventDispatcher := backend.coreEventDispatcher
		backend.coreEventDispatcher = coreEventDispatcherMock

		coreEventDispatcherMock.EXPECT().Post(gomock.Any()).DoAndReturn(func(ev any) {
			switch ev.(type) {
			case events.MessageEvent:
				if ev.(events.MessageEvent).Message().Hash() == prevote.Hash() {
					t.Log("prevote processed by core")
					passed <- struct{}{}
				}
			}
			coreEventDispatcher.Post(ev)
		}).AnyTimes()

		backend.aggregatorMessageCh <- events.UnverifiedMessageEvent{Message: prevote, ErrCh: errCh, Sender: genesisCommittee.Members[0].Address, Posted: time.Now()}
		waitFor(t, func() bool {
			select {
			case <-passed:
				return true
			default:
				// do nothing
			}
			return false
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
	t.Run("current height, future round complex aggregate carrying quorum should trigger processing", func(t *testing.T) {
		committeeSize := 4
		chain, backend := newBlockChain(committeeSize)
		genesis := chain.Genesis()
		genesisCommittee := genesis.Header().Epoch.Committee

		h := uint64(1)
		r := int64(10)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		value := common.Hash{0xca, 0xfe}
		prevote := message.NewPrevote(r, h, value, backend.Sign, &genesisCommittee.Members[0], committeeSize)
		prevote.Signers().Increment(&genesisCommittee.Members[1])
		prevote.Signers().Increment(&genesisCommittee.Members[2])
		prevote.Signers().Increment(&genesisCommittee.Members[3])

		errCh := make(chan error)

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
		backendMock.EXPECT().BlockChain().Return(newTestBlockchain()).AnyTimes()

		a := &aggregator{
			staleMessages:  make(map[common.Hash][]events.UnverifiedMessageEvent),
			messagesFrom:   make(map[common.Address][]common.Hash),
			core:           coreMock,
			backend:        backendMock,
			logger:         log.Root(),
			signerSetCache: newAggregatorCache(),
		}
		prevote := message.NewPrevote(0, h-2, common.Hash{0xca, 0xfe}, testSigner, testCommitteeMember, 1)

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
		a.handleVote(voteAEvent, quorum, true)
		require.Equal(t, 1, len(a.messages[h][r].prevotes[value]))

		a.signerSetCache.addVote(voteB, stepReceived)
		a.handleVote(voteBEvent, quorum, true)
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
		a.handleVote(singleVoteEvent, quorum, true)
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
		a.handleVote(voteEvent, quorum, true)

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
		a.handleVote(voteEventA, quorum, true)

		// quorum for * is not reached, vote should be buffered
		require.Equal(t, voteA.Hash(), a.messages[h][r].precommits[value][0].Message.Hash())

		voteB := message.NewPrecommit(r, h, common.Hash{}, testSigner, &committee.Members[3], csize)
		voteB.Signers().Increment(&committee.Members[3])
		voteB.Signers().Increment(&committee.Members[4])
		voteEventB := makeBogusEvent(voteB)
		// quorum for * is reached, vote should be processed
		a.signerSetCache.addEvent(voteEventB, stepReceived)
		a.handleVote(voteEventB, quorum, true)
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
		a.handleVote(makeBogusEvent(voteB), quorum, true)
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
		a.handleVote(makeBogusEvent(voteB), quorum, true)

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
		a.DispatchCoreEvents()
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
			internalFdCh:      make(chan events.MessageEventer, 10),
			internalCoreCh:    make(chan events.MessageEventer, 10),
			internalBacklogCh: make(chan events.UnverifiedMessageEvent, 10),
			logger:            log.Root(),
		}
		a.DispatchCoreEvents()
		a.DispatchFaultDetectorEvents()

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
		a.DispatchCoreEvents()
		a.DispatchFaultDetectorEvents()
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
	a.DispatchCoreEvents()
	a.DispatchFaultDetectorEvents()

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

func newTestBlockchain() *core.BlockChain {
	db := rawdb.NewMemoryDatabase()
	core.GenesisBlockForTesting(db, common.Address{}, common.Big0)
	chain, err := core.NewBlockChain(
		db,
		nil,
		params.TestChainConfig,
		ethash.NewFaker(),
		vm.Config{},
		nil,
		&core.TxSenderCacher{},
		nil,
		backends.NewInternalBackend(nil),
		log.Root(),
	)
	if err != nil {
		panic(err)
	}

	return chain
}

func makeTestCommitteeWithMember(index int, power *big.Int) (*types.Committee, message.Signer) {
	cc := committee.Copy()
	cc.Members = make([]types.CommitteeMember, len(committee.Members))
	var key blst.SecretKey
	var err error
	for i, member := range committee.Members {
		if i == index {
			key, err = blst.RandKey()
			if err != nil {
				panic(err)
			}
			cc.Members[i] = types.CommitteeMember{
				Address:      member.Address,
				VotingPower:  power,
				ConsensusKey: key.PublicKey(),
				Index:        member.Index,
			}
		} else {
			cc.Members[i] = member
		}
	}
	signer := func(data common.Hash) blst.Signature {
		signature := key.Sign(data[:])
		return signature
	}
	return cc, signer
}
