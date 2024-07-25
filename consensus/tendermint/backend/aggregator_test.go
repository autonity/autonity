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
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/params"
)

var (
	committee  = types.Committee{*makeBogusMember(0), *makeBogusMember(1), *makeBogusMember(2), *makeBogusMember(3), *makeBogusMember(4), *makeBogusMember(5), *makeBogusMember(6)}
	totalPower = big.NewInt(7)
	quorum     = bft.Quorum(totalPower)
	csize      = len(committee)
)

func makePropose(chain *core.BlockChain, backend *Backend, r int64, h uint64) *message.Propose {
	// don't care that it has empty proposer seal, we just want to check that aggregator sends it to Core
	currentCommittee := chain.CurrentBlock().Header().Committee
	block, err := makeBlockWithoutSeal(chain, backend, chain.CurrentBlock())
	if err != nil {
		panic("cannot create block")
	}
	propose := message.NewPropose(r, h, -1, block, backend.Sign, &currentCommittee[0])
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
	require.NoError(t, err)
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
		sub := backend.Subscribe(events.MessageEvent{})

		// don't care that it has empty proposer seal, we just want to check that aggregator sends it to Core
		propose := makePropose(chain, backend, 0, 1)

		errCh := make(chan error)

		backend.messageCh <- events.UnverifiedMessageEvent{Message: propose, ErrCh: errCh, Sender: common.Address{}, Posted: time.Now()}

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
			select {
			case ev := <-sub.Chan():
				event := ev.Data.(events.MessageEvent)
				if propose.Hash() == event.Message.Hash() {
					return true
				}
			default:
				// do nothing
			}
			return false
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
		coreMock.EXPECT().Power(gomock.Any(), gomock.Any()).Return(message.NewAggregatedPower()).Times(1)

		a := &aggregator{
			messages:     make(map[uint64]map[int64]*RoundInfo),
			messagesFrom: make(map[common.Address][]common.Hash),
			core:         coreMock,
			backend:      backendMock,
			logger:       log.Root(),
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
		chain := newTestBlockchain()

		backendMock.EXPECT().BlockChain().Return(chain).AnyTimes()
		coreMock.EXPECT().Height().Return(new(big.Int).SetUint64(h)).Times(1)
		coreMock.EXPECT().Round().Return(r).Times(1)
		coreMock.EXPECT().VotesPower(gomock.Any(), gomock.Any(), gomock.Any()).Return(message.NewAggregatedPower()).Times(2)
		coreMock.EXPECT().VotesPowerFor(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(message.NewAggregatedPower()).Times(2)

		a := &aggregator{
			messages:     make(map[uint64]map[int64]*RoundInfo),
			messagesFrom: make(map[common.Address][]common.Hash),
			core:         coreMock,
			backend:      backendMock,
			logger:       log.Root(),
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
		sub := backend.Subscribe(events.MessageEvent{})
		genesis := chain.Genesis()
		genesisCommittee := genesis.Header().Committee
		h := uint64(1)
		r := int64(0)

		prevote := message.NewPrevote(r, h, common.Hash{0xca, 0xfe}, backend.Sign, &genesisCommittee[0], committeeSize)

		errCh := make(chan error)

		backend.messageCh <- events.UnverifiedMessageEvent{Message: prevote, ErrCh: errCh, Sender: genesisCommittee[0].Address, Posted: time.Now()}

		// check that it is processed by the time-based aggr
		waitFor(t, func() bool {
			select {
			case ev := <-sub.Chan():
				event := ev.Data.(events.MessageEvent)
				if prevote.Hash() == event.Message.Hash() {
					return true
				}
			default:
				// do nothing
			}
			return false
		}, 20*time.Millisecond, 200*time.Millisecond, "prevote was not processed by the time-based aggregation")
	})
	t.Run("current height, future round prevote should be processed if F voting power is reached", func(t *testing.T) {
		committeeSize := 4
		chain, backend := newBlockChain(committeeSize)
		sub := backend.Subscribe(events.MessageEvent{})
		genesis := chain.Genesis()
		genesisCommittee := genesis.Header().Committee

		h := uint64(1)
		r := int64(10)

		// send message to the aggregator and wait for time based aggregation to send it to Core
		value := common.Hash{0xca, 0xfe}
		prevote := message.NewPrevote(r, h, value, backend.Sign, &genesisCommittee[0], committeeSize)

		errCh := make(chan error)

		backend.messageCh <- events.UnverifiedMessageEvent{Message: prevote, ErrCh: errCh, Sender: genesisCommittee[0].Address, Posted: time.Now()}
		waitFor(t, func() bool {
			select {
			case ev := <-sub.Chan():
				event := ev.Data.(events.MessageEvent)
				if prevote.Hash() == event.Message.Hash() {
					return true
				}
			default:
				// do nothing
			}
			return false
		}, 20*time.Millisecond, 200*time.Millisecond, "future round prevote has not been processed by time-based aggregation")
		require.Equal(t, uint64(100), backend.core.Power(h, r).Power().Uint64())

		// now send message that will reach quorum (together with the previous msg in Core)
		prevote = tweakPrevote(message.NewPrevote(r, h, value, backend.Sign, &genesisCommittee[1], committeeSize), backend.consensusKey.PublicKey())

		backend.messageCh <- events.UnverifiedMessageEvent{Message: prevote, ErrCh: errCh, Sender: genesisCommittee[0].Address, Posted: time.Now()}

		// core should switch to round 10 if message gets processed by it
		waitFor(t, func() bool {
			return backend.core.Round() == r
		}, 1*time.Millisecond, 30*time.Millisecond, "future round messages did not cause round change in core")
	})
	t.Run("current height, future round complex aggregate carrying quorum should trigger processing", func(t *testing.T) {
		committeeSize := 4
		chain, backend := newBlockChain(committeeSize)
		genesis := chain.Genesis()
		genesisCommittee := genesis.Header().Committee

		h := uint64(1)
		r := int64(10)

		value := common.Hash{0xca, 0xfe}
		prevote := message.NewPrevote(r, h, value, backend.Sign, &genesisCommittee[0], committeeSize)
		prevote.Signers().Increment(&genesisCommittee[1])
		prevote.Signers().Increment(&genesisCommittee[2])
		prevote.Signers().Increment(&genesisCommittee[3])

		errCh := make(chan error)

		backend.messageCh <- events.UnverifiedMessageEvent{Message: prevote, ErrCh: errCh, Sender: genesisCommittee[0].Address, Posted: time.Now()}

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

		a := &aggregator{
			staleMessages: make(map[common.Hash][]events.UnverifiedMessageEvent),
			messagesFrom:  make(map[common.Address][]common.Hash),
			core:          coreMock,
			logger:        log.Root(),
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
		subMessageEvent := backend.Subscribe(events.MessageEvent{})
		subOldMessageEvent := backend.Subscribe(events.OldMessageEvent{})

		mineOneBlock(t, chain, backend)

		genesisCommittee := genesis.Header().Committee
		prevote := message.NewPrevote(0, 1, common.Hash{0xca, 0xfe}, backend.Sign, &genesisCommittee[0], len(genesisCommittee))
		errCh := make(chan error)

		defer failIf(t, func() (bool, error) {
			select {
			case ev := <-subMessageEvent.Chan():
				event := ev.Data.(events.MessageEvent)
				if event.Message.Hash() == prevote.Hash() {
					return true, errors.New("Message was processed as current/old height")
				}
			case <-errCh:
				return true, errors.New("Message was not buffered")
			default:
				// do nothing
			}
			return false, nil
		})()

		backend.messageCh <- events.UnverifiedMessageEvent{Message: prevote, ErrCh: errCh, Sender: genesisCommittee[0].Address, Posted: time.Now()}

		// check that old message has been processed by stale messages time-based aggregation
		waitFor(t, func() bool {
			select {
			case ev := <-subOldMessageEvent.Chan():
				event := ev.Data.(events.OldMessageEvent)
				if prevote.Hash() == event.Message.Hash() {
					return true
				}
			default:
				// do nothing
			}
			return false
		}, 100*time.Millisecond, 3*time.Second, "old height message was not processed")
	})
}

func TestAggregatorSaveMessage(t *testing.T) {
	t.Run("Save proposal", func(t *testing.T) {
		a := &aggregator{messages: make(map[uint64]map[int64]*RoundInfo)}

		r := int64(4)
		h := uint64(2)

		propose := makeBogusPropose(r, h, 0)

		proposeEvent := makeBogusEvent(propose)

		a.saveMessage(proposeEvent)

		roundInfo := a.messages[h][r]

		require.Equal(t, propose.Hash(), roundInfo.proposals[0].Message.Hash())
		require.Equal(t, common.Big1, roundInfo.power.Power())
		require.Equal(t, uint(1), roundInfo.power.Signers().Bit(0))
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
		require.Equal(t, common.Big1, roundInfo.prevotesPowerFor[value].Power())
		require.Equal(t, uint(1), roundInfo.prevotesPowerFor[value].Signers().Bit(0))
		require.Equal(t, common.Big1, roundInfo.prevotesPower.Power())
		require.Equal(t, uint(1), roundInfo.prevotesPower.Signers().Bit(0))
		require.Equal(t, common.Big1, roundInfo.power.Power())
		require.Equal(t, uint(1), roundInfo.power.Signers().Bit(0))
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
		require.Equal(t, common.Big1, roundInfo.precommitsPowerFor[value].Power())
		require.Equal(t, uint(1), roundInfo.precommitsPowerFor[value].Signers().Bit(0))
		require.Equal(t, common.Big1, roundInfo.precommitsPower.Power())
		require.Equal(t, uint(1), roundInfo.precommitsPower.Signers().Bit(0))
		require.Equal(t, common.Big1, roundInfo.power.Power())
		require.Equal(t, uint(1), roundInfo.power.Signers().Bit(0))
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
		require.Equal(t, common.Big1, roundInfo.power.Power())
		require.Equal(t, uint(1), roundInfo.power.Signers().Bit(0))

		// save again the same proposal
		a.saveMessage(proposeEvent)

		roundInfo = a.messages[h][r]
		require.Equal(t, propose.Hash(), roundInfo.proposals[1].Message.Hash())
		require.Equal(t, common.Big1, roundInfo.power.Power())
		require.Equal(t, uint(1), roundInfo.power.Signers().Bit(0))

		// save precommit for same (h,r) from the same guy
		precommit := message.NewPrecommit(r, h, value, testSigner, testCommitteeMember, committeeSize)
		voteEvent := events.UnverifiedMessageEvent{Message: precommit, ErrCh: nil, Sender: common.Address{}, Posted: time.Now()}

		a.saveMessage(voteEvent)

		roundInfo = a.messages[h][r]

		require.Equal(t, precommit.Hash(), roundInfo.precommits[value][0].Message.Hash())
		require.Equal(t, common.Big1, roundInfo.precommitsPowerFor[value].Power())
		require.Equal(t, uint(1), roundInfo.precommitsPowerFor[value].Signers().Bit(0))
		require.Equal(t, common.Big1, roundInfo.precommitsPower.Power())
		require.Equal(t, uint(1), roundInfo.precommitsPower.Signers().Bit(0))
		require.Equal(t, common.Big1, roundInfo.power.Power())
		require.Equal(t, uint(1), roundInfo.power.Signers().Bit(0))

		// save vote for same (h,r,v) from different validator
		prevote := message.NewPrevote(r, h, value, testSigner, makeBogusMember(1), committeeSize)
		voteEvent = events.UnverifiedMessageEvent{Message: prevote, ErrCh: nil, Sender: common.Address{}, Posted: time.Now()}

		a.saveMessage(voteEvent)

		roundInfo = a.messages[h][r]

		require.Equal(t, prevote.Hash(), roundInfo.prevotes[value][0].Message.Hash())
		require.Equal(t, common.Big1, roundInfo.prevotesPowerFor[value].Power())
		require.Equal(t, uint(1), roundInfo.prevotesPowerFor[value].Signers().Bit(1))
		require.Equal(t, common.Big1, roundInfo.prevotesPower.Power())
		require.Equal(t, uint(1), roundInfo.prevotesPower.Signers().Bit(1))
		require.Equal(t, common.Big2, roundInfo.power.Power())
		require.Equal(t, uint(1), roundInfo.power.Signers().Bit(0))
		require.Equal(t, uint(1), roundInfo.power.Signers().Bit(1))

		// save aggregated vote for same (h,r), different value
		otherValue := common.Hash{0x13, 0x37}
		prevote = message.NewPrevote(r, h, otherValue, testSigner, makeBogusMember(0), committeeSize)
		prevote.Signers().Increment(makeBogusMember(1))
		prevote.Signers().Increment(makeBogusMember(2))
		voteEvent = events.UnverifiedMessageEvent{Message: prevote, ErrCh: nil, Sender: common.Address{}, Posted: time.Now()}

		a.saveMessage(voteEvent)

		roundInfo = a.messages[h][r]

		require.Equal(t, prevote.Hash(), roundInfo.prevotes[otherValue][0].Message.Hash())
		require.Equal(t, common.Big3, roundInfo.prevotesPowerFor[otherValue].Power())
		require.Equal(t, uint(1), roundInfo.prevotesPowerFor[otherValue].Signers().Bit(0))
		require.Equal(t, uint(1), roundInfo.prevotesPowerFor[otherValue].Signers().Bit(1))
		require.Equal(t, uint(1), roundInfo.prevotesPowerFor[otherValue].Signers().Bit(2))
		require.Equal(t, common.Big3, roundInfo.prevotesPower.Power())
		require.Equal(t, uint(1), roundInfo.prevotesPower.Signers().Bit(0))
		require.Equal(t, uint(1), roundInfo.prevotesPower.Signers().Bit(1))
		require.Equal(t, uint(1), roundInfo.prevotesPower.Signers().Bit(2))
		require.Equal(t, common.Big3, roundInfo.power.Power())
		require.Equal(t, uint(1), roundInfo.power.Signers().Bit(0))
		require.Equal(t, uint(1), roundInfo.power.Signers().Bit(1))
		require.Equal(t, uint(1), roundInfo.power.Signers().Bit(2))

		// save proposal from index 3

		propose = makeBogusPropose(r, h, 3)
		proposeEvent = events.UnverifiedMessageEvent{Message: propose, ErrCh: nil, Sender: common.Address{}, Posted: time.Now()}
		a.saveMessage(proposeEvent)

		roundInfo = a.messages[h][r]
		require.Equal(t, propose.Hash(), roundInfo.proposals[2].Message.Hash())
		require.Equal(t, common.Big4, roundInfo.power.Power())
		require.Equal(t, uint(1), roundInfo.power.Signers().Bit(0))
		require.Equal(t, uint(1), roundInfo.power.Signers().Bit(1))
		require.Equal(t, uint(1), roundInfo.power.Signers().Bit(2))
		require.Equal(t, uint(1), roundInfo.power.Signers().Bit(3))

		// save aggregated vote for same (h,r), different value
		precommit = message.NewPrecommit(r, h, otherValue, testSigner, makeBogusMember(0), committeeSize)
		precommit.Signers().Increment(makeBogusMember(1))
		precommit.Signers().Increment(makeBogusMember(2))
		precommit.Signers().Increment(makeBogusMember(3))
		voteEvent = events.UnverifiedMessageEvent{Message: precommit, ErrCh: nil, Sender: common.Address{}, Posted: time.Now()}

		a.saveMessage(voteEvent)

		roundInfo = a.messages[h][r]

		require.Equal(t, precommit.Hash(), roundInfo.precommits[otherValue][0].Message.Hash())
		require.Equal(t, common.Big4, roundInfo.precommitsPowerFor[otherValue].Power())
		require.Equal(t, uint(1), roundInfo.precommitsPowerFor[otherValue].Signers().Bit(0))
		require.Equal(t, uint(1), roundInfo.precommitsPowerFor[otherValue].Signers().Bit(1))
		require.Equal(t, uint(1), roundInfo.precommitsPowerFor[otherValue].Signers().Bit(2))
		require.Equal(t, uint(1), roundInfo.precommitsPowerFor[otherValue].Signers().Bit(3))
		require.Equal(t, common.Big4, roundInfo.precommitsPower.Power())
		require.Equal(t, uint(1), roundInfo.precommitsPower.Signers().Bit(0))
		require.Equal(t, uint(1), roundInfo.precommitsPower.Signers().Bit(1))
		require.Equal(t, uint(1), roundInfo.precommitsPower.Signers().Bit(2))
		require.Equal(t, uint(1), roundInfo.precommitsPower.Signers().Bit(3))
		require.Equal(t, common.Big4, roundInfo.power.Power())
		require.Equal(t, uint(1), roundInfo.power.Signers().Bit(0))
		require.Equal(t, uint(1), roundInfo.power.Signers().Bit(1))
		require.Equal(t, uint(1), roundInfo.power.Signers().Bit(2))
		require.Equal(t, uint(1), roundInfo.power.Signers().Bit(3))

		// save aggregated vote for same (h,r), different value
		precommit = message.NewPrecommit(r, h, otherValue, testSigner, makeBogusMember(0), committeeSize)
		precommit.Signers().Increment(makeBogusMember(1))
		voteEvent = events.UnverifiedMessageEvent{Message: precommit, ErrCh: nil, Sender: common.Address{}, Posted: time.Now()}

		a.saveMessage(voteEvent)

		roundInfo = a.messages[h][r]

		require.Equal(t, 2, len(roundInfo.precommits[otherValue]))
		require.Equal(t, precommit.Hash(), roundInfo.precommits[otherValue][1].Message.Hash())
		require.Equal(t, common.Big4, roundInfo.precommitsPowerFor[otherValue].Power())
		require.Equal(t, uint(1), roundInfo.precommitsPowerFor[otherValue].Signers().Bit(0))
		require.Equal(t, uint(1), roundInfo.precommitsPowerFor[otherValue].Signers().Bit(1))
		require.Equal(t, uint(1), roundInfo.precommitsPowerFor[otherValue].Signers().Bit(2))
		require.Equal(t, uint(1), roundInfo.precommitsPowerFor[otherValue].Signers().Bit(3))
		require.Equal(t, common.Big4, roundInfo.precommitsPower.Power())
		require.Equal(t, uint(1), roundInfo.precommitsPower.Signers().Bit(0))
		require.Equal(t, uint(1), roundInfo.precommitsPower.Signers().Bit(1))
		require.Equal(t, uint(1), roundInfo.precommitsPower.Signers().Bit(2))
		require.Equal(t, uint(1), roundInfo.precommitsPower.Signers().Bit(3))
		require.Equal(t, common.Big4, roundInfo.power.Power())
		require.Equal(t, uint(1), roundInfo.power.Signers().Bit(0))
		require.Equal(t, uint(1), roundInfo.power.Signers().Bit(1))
		require.Equal(t, uint(1), roundInfo.power.Signers().Bit(2))
		require.Equal(t, uint(1), roundInfo.power.Signers().Bit(3))

		// save vote for different round
		otherRound := r + 10
		precommit = message.NewPrecommit(otherRound, h, value, testSigner, testCommitteeMember, committeeSize)
		voteEvent = events.UnverifiedMessageEvent{Message: precommit, ErrCh: nil, Sender: common.Address{}, Posted: time.Now()}

		a.saveMessage(voteEvent)

		roundInfo = a.messages[h][otherRound]

		require.Equal(t, precommit.Hash(), roundInfo.precommits[value][0].Message.Hash())
		require.Equal(t, common.Big1, roundInfo.precommitsPowerFor[value].Power())
		require.Equal(t, uint(1), roundInfo.precommitsPowerFor[value].Signers().Bit(0))
		require.Equal(t, common.Big1, roundInfo.precommitsPower.Power())
		require.Equal(t, uint(1), roundInfo.precommitsPower.Signers().Bit(0))
		require.Equal(t, common.Big1, roundInfo.power.Power())
		require.Equal(t, uint(1), roundInfo.power.Signers().Bit(0))

		// save vote for different height
		otherHeight := h + 3
		precommit = message.NewPrecommit(otherRound, otherHeight, value, testSigner, testCommitteeMember, committeeSize)
		voteEvent = events.UnverifiedMessageEvent{Message: precommit, ErrCh: nil, Sender: common.Address{}, Posted: time.Now()}

		a.saveMessage(voteEvent)

		roundInfo = a.messages[otherHeight][otherRound]

		require.Equal(t, precommit.Hash(), roundInfo.precommits[value][0].Message.Hash())
		require.Equal(t, common.Big1, roundInfo.precommitsPowerFor[value].Power())
		require.Equal(t, uint(1), roundInfo.precommitsPowerFor[value].Signers().Bit(0))
		require.Equal(t, common.Big1, roundInfo.precommitsPower.Power())
		require.Equal(t, uint(1), roundInfo.precommitsPower.Signers().Bit(0))
		require.Equal(t, common.Big1, roundInfo.power.Power())
		require.Equal(t, uint(1), roundInfo.power.Signers().Bit(0))
	})
}

func TestAggregatorHandleVote(t *testing.T) {
	quorumMinusOne := new(big.Int).Set(quorum)
	quorumMinusOne.Sub(quorumMinusOne, common.Big1)

	t.Run("complex aggregate handling", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer waitForExpects(t, ctrl)

		backendMock := interfaces.NewMockBackend(ctrl)
		coreMock := interfaces.NewMockCore(ctrl)

		a := &aggregator{
			messages:      make(map[uint64]map[int64]*RoundInfo),
			core:          coreMock,
			backend:       backendMock,
			knownMessages: fixsizecache.New[common.Hash, bool](numBuckets, numEntries, fixsizecache.HashKey[common.Hash]),
			logger:        log.Root(),
		}

		r := int64(5)
		h := uint64(0)
		value := common.Hash{0xca, 0xfe}

		// complex aggregate is processed right away if core doesn't have quorum
		vote := message.NewPrecommit(r, h, value, testSigner, &committee[0], csize)
		vote.Signers().Increment(&committee[0])
		voteEvent := makeBogusEvent(vote)

		coreMock.EXPECT().VotesPowerFor(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(message.NewAggregatedPower()).Times(1)
		coreMock.EXPECT().VotesPower(gomock.Any(), gomock.Any(), gomock.Any()).Return(message.NewAggregatedPower()).Times(1)
		backendMock.EXPECT().Post(gomock.Any()).Times(1)
		a.handleVote(voteEvent, committee, quorum, true)
		require.Nil(t, a.messages[h])

		// complex aggregate is processed right away if core:
		// - does have quorum for *
		// - does not have quorum for v
		corePower := message.NewAggregatedPower()
		corePower.Set(0, quorum)
		coreMock.EXPECT().VotesPowerFor(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(message.NewAggregatedPower()).Times(1)
		coreMock.EXPECT().VotesPower(gomock.Any(), gomock.Any(), gomock.Any()).Return(corePower).Times(1)
		backendMock.EXPECT().Post(gomock.Any()).Times(1)
		a.handleVote(voteEvent, committee, quorum, true)
		require.Nil(t, a.messages[h])

		// complex aggregate is buffered if core:
		// - does have quorum for v
		// - does have quorum for *
		corePowerForV := corePower.Copy()
		coreMock.EXPECT().VotesPowerFor(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(corePowerForV).Times(2)
		coreMock.EXPECT().VotesPower(gomock.Any(), gomock.Any(), gomock.Any()).Return(corePower).Times(2)
		a.handleVote(voteEvent, committee, quorum, true)
		require.Equal(t, vote.Hash(), a.messages[h][r].precommits[value][0].Message.Hash())
	})
	t.Run("individual and simple aggregates handling", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer func() {
			// wait for all EXPECTS to be satisfied before calling `Finish()`
			waitFor(t, func() bool {
				return ctrl.Satisfied()
			}, 10*time.Millisecond, 1*time.Second, "mock EXPECTS() are not satisfied")
			ctrl.Finish()
		}()
		backendMock := interfaces.NewMockBackend(ctrl)
		coreMock := interfaces.NewMockCore(ctrl)

		a := &aggregator{
			messages:      make(map[uint64]map[int64]*RoundInfo),
			core:          coreMock,
			backend:       backendMock,
			knownMessages: fixsizecache.New[common.Hash, bool](numBuckets, numEntries, fixsizecache.HashKey[common.Hash]),
			logger:        log.Root(),
		}

		r := int64(5)
		h := uint64(0)

		value := common.Hash{0xca, 0xfe}
		vote := message.NewPrecommit(r, h, value, testSigner, &committee[0], csize)
		voteEvent := makeBogusEvent(vote)

		// no quorum reached, vote should be buffered
		coreMock.EXPECT().VotesPowerFor(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(message.NewAggregatedPower()).Times(2)
		coreMock.EXPECT().VotesPower(gomock.Any(), gomock.Any(), gomock.Any()).Return(message.NewAggregatedPower()).Times(2)
		a.handleVote(voteEvent, committee, quorum, true)
		require.Equal(t, vote.Hash(), a.messages[h][r].precommits[value][0].Message.Hash())

		// simple aggregate with quorum should trigger processing
		vote = message.NewPrecommit(r, h, value, testSigner, &committee[0], csize)
		vote.Signers().Increment(&committee[1])
		vote.Signers().Increment(&committee[2])
		vote.Signers().Increment(&committee[3])
		vote.Signers().Increment(&committee[4])
		voteEvent = makeBogusEvent(vote)

		coreMock.EXPECT().VotesPowerFor(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(message.NewAggregatedPower()).Times(2)
		coreMock.EXPECT().VotesPower(gomock.Any(), gomock.Any(), gomock.Any()).Return(message.NewAggregatedPower()).Times(1)
		backendMock.EXPECT().Post(gomock.Any()).Times(1)
		a.handleVote(voteEvent, committee, quorum, true)
		require.Equal(t, 0, len(a.messages[h][r].precommits[value]))
		require.Nil(t, a.messages[h][r].precommitsPowerFor[value])
		require.Equal(t, common.Big0, a.messages[h][r].precommitsPower.Power())
		require.Equal(t, common.Big0, a.messages[h][r].power.Power())

		// half voting power for v, half for nil should trigger processing (to trigger timeouts in core)
		// but first let's process some equivocated votes. These should not trigger processing
		voteForV := message.NewPrevote(r, h, value, testSigner, &committee[0], csize)
		voteForV.Signers().Increment(&committee[1])
		voteForV.Signers().Increment(&committee[2])
		eventForV := makeBogusEvent(voteForV)
		equivocatedVoteForNil := message.NewPrevote(r, h, common.Hash{}, testSigner, &committee[0], csize)
		equivocatedVoteForNil.Signers().Increment(&committee[1])
		equivocatedEvent := makeBogusEvent(equivocatedVoteForNil)

		coreMock.EXPECT().VotesPowerFor(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(message.NewAggregatedPower()).Times(4)
		coreMock.EXPECT().VotesPower(gomock.Any(), gomock.Any(), gomock.Any()).Return(message.NewAggregatedPower()).Times(4)
		a.handleVote(eventForV, committee, quorum, true)
		a.handleVote(equivocatedEvent, committee, quorum, true)

		// now votes for Nil from other validators will trigger the processing
		voteForNil := message.NewPrevote(r, h, common.Hash{}, testSigner, &committee[3], csize)
		voteForNil.Signers().Increment(&committee[4])
		eventForNil := makeBogusEvent(voteForNil)

		coreMock.EXPECT().VotesPowerFor(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(message.NewAggregatedPower()).Times(2)
		coreMock.EXPECT().VotesPower(gomock.Any(), gomock.Any(), gomock.Any()).Return(message.NewAggregatedPower()).Times(2)
		backendMock.EXPECT().Post(gomock.Any()).Times(2)
		a.handleVote(eventForNil, committee, quorum, true)

		require.Equal(t, 0, len(a.messages[h][r].prevotes))
		require.Equal(t, common.Big0, a.messages[h][r].prevotesPower.Power())
		require.Equal(t, common.Big0, a.messages[h][r].power.Power())

		// if we reach quorum for V with the voting power in core, message is processed
		vote = message.NewPrecommit(r, h, value, testSigner, &committee[0], csize)
		voteEvent = makeBogusEvent(vote)

		corePowerForV := message.NewAggregatedPower()
		corePowerForV.Set(1, quorumMinusOne)

		coreMock.EXPECT().VotesPowerFor(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(corePowerForV).Times(2)
		coreMock.EXPECT().VotesPower(gomock.Any(), gomock.Any(), gomock.Any()).Return(message.NewAggregatedPower()).Times(1)
		backendMock.EXPECT().Post(gomock.Any()).Times(1)
		a.handleVote(voteEvent, committee, quorum, true)
		require.Equal(t, 0, len(a.messages[h][r].precommits))

		// if we reach quorum for * with the voting power in core, message is processed
		vote = message.NewPrecommit(r, h, value, testSigner, &committee[0], csize)
		voteEvent = makeBogusEvent(vote)

		corePower := message.NewAggregatedPower()
		corePower.Set(1, quorumMinusOne)

		coreMock.EXPECT().VotesPowerFor(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(message.NewAggregatedPower()).Times(2)
		coreMock.EXPECT().VotesPower(gomock.Any(), gomock.Any(), gomock.Any()).Return(corePower).Times(2)
		backendMock.EXPECT().Post(gomock.Any()).Times(1)
		a.handleVote(voteEvent, committee, quorum, true)
		require.Equal(t, 0, len(a.messages[h][r].precommits))

		// if voting power from the vote is already in core, msg gets buffered
		vote = message.NewPrecommit(r, h, value, testSigner, &committee[0], csize)
		vote.Signers().Increment(&committee[1])
		vote.Signers().Increment(&committee[2])
		voteEvent = makeBogusEvent(vote)

		corePower = message.NewAggregatedPower()
		corePower.Set(0, common.Big1)
		corePower.Set(1, common.Big1)
		corePower.Set(2, common.Big1)
		corePower.Set(3, common.Big1)

		coreMock.EXPECT().VotesPowerFor(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(message.NewAggregatedPower()).Times(2)
		coreMock.EXPECT().VotesPower(gomock.Any(), gomock.Any(), gomock.Any()).Return(corePower).Times(2)
		a.handleVote(voteEvent, committee, quorum, true)
		require.Equal(t, 1, len(a.messages[h][r].precommits))

		// message makes us reach quorum for v, gets processed
		vote = message.NewPrecommit(r, h, value, testSigner, &committee[0], csize)
		vote.Signers().Increment(&committee[1])
		vote.Signers().Increment(&committee[2])
		voteEvent = makeBogusEvent(vote)

		corePowerForV = message.NewAggregatedPower()
		corePowerForV.Set(0, common.Big1)
		corePowerForV.Set(3, common.Big1)
		corePowerForV.Set(4, common.Big1)

		coreMock.EXPECT().VotesPowerFor(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(corePowerForV).Times(2)
		coreMock.EXPECT().VotesPower(gomock.Any(), gomock.Any(), gomock.Any()).Return(message.NewAggregatedPower()).Times(1)
		backendMock.EXPECT().Post(gomock.Any()).Times(1)
		a.handleVote(voteEvent, committee, quorum, true)
		require.Equal(t, 0, len(a.messages[h][r].precommits))
	})
}

func TestAggregatorPowerContribution(t *testing.T) {
	aggregatorPower := new(big.Int)
	corePower := new(big.Int)

	require.Equal(t, common.Big0, powerContribution(aggregatorPower, corePower, committee))

	corePower.SetBit(corePower, 0, 1)
	require.Equal(t, common.Big0, powerContribution(aggregatorPower, corePower, committee))

	aggregatorPower.SetBit(aggregatorPower, 0, 1)
	require.Equal(t, common.Big0, powerContribution(aggregatorPower, corePower, committee))

	aggregatorPower.SetBit(aggregatorPower, 1, 1)
	require.Equal(t, common.Big1, powerContribution(aggregatorPower, corePower, committee))

	aggregatorPower.SetBit(aggregatorPower, 2, 1)
	aggregatorPower.SetBit(aggregatorPower, 3, 1)
	require.Equal(t, common.Big3, powerContribution(aggregatorPower, corePower, committee))

	corePower.SetBit(corePower, 2, 1)
	corePower.SetBit(corePower, 5, 1)
	corePower.SetBit(corePower, 6, 1)
	require.Equal(t, common.Big2, powerContribution(aggregatorPower, corePower, committee))
}

func TestAggregatorProcess(t *testing.T) {
	r := int64(4)
	h := uint64(23)
	var messages []message.Msg
	value := common.Hash{0xca, 0xfe}
	otherValue := common.Hash{0xff, 0xff}
	messages = append(messages, makeBogusPropose(r, h, 0))
	messages = append(messages, makeBogusPropose(r, h, 1))
	messages = append(messages, message.NewPrevote(r, h, value, testSigner, &committee[0], csize))
	messages = append(messages, message.NewPrecommit(r, h, value, testSigner, &committee[0], csize))
	messages = append(messages, message.NewPrecommit(r+2, h, value, testSigner, &committee[0], csize))
	messages = append(messages, message.NewPrevote(r, h+3, value, testSigner, &committee[0], csize))
	messages = append(messages, message.NewPrevote(r, h, otherValue, testSigner, &committee[0], csize))

	t.Run("processProposal, valid proposal is posted", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer waitForExpects(t, ctrl)
		backendMock := interfaces.NewMockBackend(ctrl)

		a := &aggregator{backend: backendMock}

		backendMock.EXPECT().Post(gomock.Any()).Times(1)
		// signature is invalid but proposal is created with `verified`=true, so it is considered valid
		propose := makeBogusPropose(0, 1, 0)
		proposeEvent := makeBogusEvent(propose)
		a.processProposal(proposeEvent, func(_ message.Msg, _ chan<- error) interface{} { return struct{}{} })
	})
	t.Run("processProposal, invalid proposal is rejected", func(t *testing.T) {
		a := &aggregator{
			messagesFrom: make(map[common.Address][]common.Hash),
		}

		propose := message.NewFakePropose(message.Fake{
			FakeSignatureInput: common.Hash{0xca, 0xfe},
			FakeSignerKey:      testKey.PublicKey(),
			FakeSignature:      testKey.Sign([]byte{0xff, 0xff}), // signature is not on FakeSignatureInput --> invalid
			FakeVerified:       false,
			FakeHash:           common.Hash{0xee, 0xee},
		})

		sender := common.Address{0xaa, 0xaa}
		errCh := make(chan error)
		proposeEvent := events.UnverifiedMessageEvent{Message: propose, ErrCh: errCh, Sender: sender, Posted: time.Now()}

		wg := new(sync.WaitGroup)
		wg.Add(1)
		go func() {
			require.Equal(t, message.ErrBadSignature, <-errCh)
			wg.Done()
		}()
		a.processProposal(proposeEvent, func(_ message.Msg, _ chan<- error) interface{} { return struct{}{} })
		wg.Wait()
	})
	t.Run("processRound processes all the messages for a round", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(gomock.Any()).AnyTimes()

		a := &aggregator{
			messages:      make(map[uint64]map[int64]*RoundInfo),
			messagesFrom:  make(map[common.Address][]common.Hash),
			backend:       backendMock,
			knownMessages: fixsizecache.New[common.Hash, bool](numBuckets, numEntries, fixsizecache.HashKey[common.Hash]),
			logger:        log.Root(),
		}

		for _, message := range messages {
			a.saveMessage(makeBogusEvent(message))
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

		a := &aggregator{
			messages:      make(map[uint64]map[int64]*RoundInfo),
			messagesFrom:  make(map[common.Address][]common.Hash),
			backend:       backendMock,
			knownMessages: fixsizecache.New[common.Hash, bool](numBuckets, numEntries, fixsizecache.HashKey[common.Hash]),
			logger:        log.Root(),
		}

		for _, message := range messages {
			a.saveMessage(makeBogusEvent(message))
		}

		require.Equal(t, 2, len(a.messages[h][r].prevotes))

		a.processVotes(h, r, message.PrevoteCode)

		require.Equal(t, 0, len(a.messages[h][r].prevotes))
		require.Equal(t, common.Big0, a.messages[h][r].prevotesPower.Power())
		require.Equal(t, 0, len(a.messages[h][r].prevotesPowerFor))
		require.NotEqual(t, common.Big0, a.messages[h][r].power.Power())
	})
	t.Run("processVotesFor processes all prevotes OR precommits for a value for a round", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(gomock.Any()).AnyTimes()

		a := &aggregator{
			messages:      make(map[uint64]map[int64]*RoundInfo),
			messagesFrom:  make(map[common.Address][]common.Hash),
			backend:       backendMock,
			knownMessages: fixsizecache.New[common.Hash, bool](numBuckets, numEntries, fixsizecache.HashKey[common.Hash]),
			logger:        log.Root(),
		}

		for _, message := range messages {
			a.saveMessage(makeBogusEvent(message))
		}

		require.Equal(t, 1, len(a.messages[h][r].prevotes[value]))

		a.processVotesFor(h, r, message.PrevoteCode, value)

		require.Equal(t, 0, len(a.messages[h][r].prevotes[value]))
		require.Equal(t, 1, len(a.messages[h][r].prevotes[otherValue]))
		require.Nil(t, a.messages[h][r].prevotesPowerFor[value])
		require.NotEqual(t, common.Big0, a.messages[h][r].prevotesPower.Power())
		require.NotEqual(t, common.Big0, a.messages[h][r].power.Power())
	})
	t.Run("ProcessBatch posts events when batches are valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer waitForExpects(t, ctrl)

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(gomock.Any()).Times(4)

		a := &aggregator{
			backend:       backendMock,
			knownMessages: fixsizecache.New[common.Hash, bool](numBuckets, numEntries, fixsizecache.HashKey[common.Hash]),
			logger:        log.Root(),
		}

		var batches [][]events.UnverifiedMessageEvent

		batches = append(batches, []events.UnverifiedMessageEvent{
			makeBogusEvent(message.NewPrevote(r, h, value, testSigner, &committee[0], csize)),
			makeBogusEvent(message.NewPrevote(r, h, value, testSigner, &committee[3], csize)),
			makeBogusEvent(message.NewPrevote(r, h, value, testSigner, &committee[4], csize)),
			makeBogusEvent(message.NewPrevote(r, h, value, testSigner, &committee[5], csize)),
		})
		batches = append(batches, []events.UnverifiedMessageEvent{
			makeBogusEvent(message.NewPrecommit(r, h, otherValue, testSigner, &committee[0], csize)),
			makeBogusEvent(message.NewPrecommit(r, h, otherValue, testSigner, &committee[3], csize)),
			makeBogusEvent(message.NewPrecommit(r, h, otherValue, testSigner, &committee[4], csize)),
			makeBogusEvent(message.NewPrecommit(r, h, otherValue, testSigner, &committee[5], csize)),
		})
		// NOTE: this will trigger two calls to Post because the votes cannot be merged in a simple aggregate
		aggregate1 := message.AggregatePrecommitsSimple([]message.Vote{message.NewPrecommit(r, h, value, testSigner, &committee[0], csize), message.NewPrecommit(r, h, value, testSigner, &committee[3], csize)})
		aggregate2 := message.AggregatePrecommitsSimple([]message.Vote{message.NewPrecommit(r, h, value, testSigner, &committee[0], csize), message.NewPrecommit(r, h, value, testSigner, &committee[4], csize)})
		batches = append(batches, []events.UnverifiedMessageEvent{
			makeBogusEvent(aggregate1[0]),
			makeBogusEvent(aggregate2[0]),
		})

		a.processBatches(batches, func(_ message.Msg, _ chan<- error) interface{} { return struct{}{} })
	})
	t.Run("ProcessBatch successfully detects and discard invalid signatures", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer waitForExpects(t, ctrl)

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(gomock.Any()).Times(3)

		a := &aggregator{
			backend:       backendMock,
			knownMessages: fixsizecache.New[common.Hash, bool](numBuckets, numEntries, fixsizecache.HashKey[common.Hash]),
			logger:        log.Root(),
		}

		var batches [][]events.UnverifiedMessageEvent

		// committee[3] and committee[6] are sending invalid sigs

		batches = append(batches, []events.UnverifiedMessageEvent{
			makeBogusEvent(message.NewPrevote(r, h, value, testInvalidSigner, &committee[3], csize)), // INVALID
		})
		// message is correctly signed only by committee[3], but signers include also committee[0]
		invalidSignersPrevote := message.NewPrevote(r, h, value, testSigner, &committee[3], csize)
		invalidSignersPrevote.Signers().Increment(&committee[0])
		aggKey, err := blst.AggregatePublicKeys([]blst.PublicKey{committee[3].ConsensusKey, committee[0].ConsensusKey})
		require.NoError(t, err)
		invalidSignersPrevote = tweakPrevote(invalidSignersPrevote, aggKey)
		batches = append(batches, []events.UnverifiedMessageEvent{
			makeBogusEvent(message.NewPrevote(r, h, value, testSigner, &committee[0], csize)),
			makeBogusEvent(message.NewPrevote(r, h, value, testInvalidSigner, &committee[3], csize)), // INVALID
			makeBogusEvent(invalidSignersPrevote),                                                    // INVALID
			makeBogusEvent(message.NewPrevote(r, h, value, testSigner, &committee[4], csize)),
			makeBogusEvent(message.NewPrevote(r, h, value, testSigner, &committee[5], csize)),
		})
		batches = append(batches, []events.UnverifiedMessageEvent{
			makeBogusEvent(message.NewPrecommit(r, h, otherValue, testSigner, &committee[0], csize)),
			makeBogusEvent(message.NewPrecommit(r, h, otherValue, testInvalidSigner, &committee[3], csize)), //INVALID
			makeBogusEvent(message.NewPrecommit(r, h, otherValue, testSigner, &committee[4], csize)),
			makeBogusEvent(message.NewPrecommit(r, h, otherValue, testSigner, &committee[5], csize)),
			makeBogusEvent(message.NewPrecommit(r, h, otherValue, testInvalidSigner, &committee[6], csize)), //INVALID
		})
		// NOTE: this will trigger two calls to Post because the votes cannot be merged in a simple aggregate
		aggregate1 := message.AggregatePrecommitsSimple([]message.Vote{message.NewPrecommit(r, h, value, testSigner, &committee[0], csize), message.NewPrecommit(r, h, value, testSigner, &committee[4], csize)})
		aggregate2 := message.AggregatePrecommitsSimple([]message.Vote{message.NewPrecommit(r, h, value, testSigner, &committee[0], csize), message.NewPrecommit(r, h, value, testInvalidSigner, &committee[3], csize)}) // INVALID
		batches = append(batches, []events.UnverifiedMessageEvent{
			makeBogusEvent(aggregate1[0]),
			makeBogusEvent(aggregate2[0]), //INVALID
			makeBogusEvent(message.NewPrecommit(r, h, value, testInvalidSigner, &committee[6], csize)), // INVALID
			makeBogusEvent(message.NewPrecommit(r, h, value, testSigner, &committee[5], csize)),
		})

		a.processBatches(batches, func(m message.Msg, _ chan<- error) interface{} {
			vote, ok := m.(message.Vote)
			require.True(t, ok)
			if vote.Signers().Contains(3) || vote.Signers().Contains(6) {
				t.Fatalf("Invalid message has been posted")
			}
			return struct{}{}
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
	backendMock.EXPECT().Post(gomock.Any()).Times(1)

	a := &aggregator{
		backend:       backendMock,
		knownMessages: fixsizecache.New[common.Hash, bool](numBuckets, numEntries, fixsizecache.HashKey[common.Hash]),
		logger:        log.Root(),
		messagesFrom:  make(map[common.Address][]common.Hash),
		messages:      make(map[uint64]map[int64]*RoundInfo),
		toIgnore:      make(map[common.Hash]struct{}),
	}

	// suppose committee[0] is sending invalid sigs

	zeroAddress := common.Address{0xff, 0xff, 0xff}
	var votesFromZero []message.Msg
	// mix of valid and invalid votes, from different validators
	votesFromZero = append(votesFromZero, message.NewPrevote(r+1, h, value, testSigner, &committee[0], csize), message.NewPrevote(r, h+3, value, testSigner, &committee[1], csize), message.NewPrevote(r, h, value, testInvalidSigner, &committee[2], csize))

	for _, vote := range votesFromZero {
		a.messagesFrom[zeroAddress] = append(a.messagesFrom[zeroAddress], vote.Hash())
		a.saveMessage(events.UnverifiedMessageEvent{Message: vote, ErrCh: nil, Sender: zeroAddress, Posted: time.Now()})
	}

	// add another vote, coming from an honest validator (at p2p layer)
	vote := message.NewPrevote(r, h, value, testSigner, &committee[1], csize)
	a.saveMessage(makeBogusEvent(vote))
	a.messagesFrom[common.Address{}] = append(a.messagesFrom[common.Address{}], vote.Hash())

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

func TestAggregatorCoreEvents(t *testing.T) {
	t.Run("RoundChangeEvent triggers rules re-evaluation - proposal", func(t *testing.T) {
		futureHeight := uint64(1)
		futureRound := int64(2)
		eventMux := event.NewTypeMuxSilent(nil, log.Root())

		ctrl := gomock.NewController(t)
		defer waitForExpects(t, ctrl)

		coreMock := interfaces.NewMockCore(ctrl)
		backendMock := interfaces.NewMockBackend(ctrl)
		chain := newTestBlockchain()

		backendMock.EXPECT().MessageCh().Return(make(chan events.UnverifiedMessageEvent)).Times(1)
		backendMock.EXPECT().BlockChain().Return(chain).AnyTimes()
		backendMock.EXPECT().Post(gomock.Any()).Times(1)
		backendMock.EXPECT().Subscribe(gomock.Any(), gomock.Any(), gomock.Any()).Return(eventMux.Subscribe(events.RoundChangeEvent{})).Times(1)
		coreMock.EXPECT().Height().Return(common.Big1).AnyTimes()

		a := &aggregator{
			core:     coreMock,
			backend:  backendMock,
			messages: make(map[uint64]map[int64]*RoundInfo),
			logger:   log.Root(),
		}

		// save a proposal for a the future round
		a.saveMessage(makeBogusEvent(makeBogusPropose(futureRound, futureHeight, 0)))

		a.start(context.TODO())

		eventMux.Post(events.RoundChangeEvent{Height: futureHeight, Round: futureRound})

		// give time to the aggregator loop to process the proposal
		time.Sleep(10 * time.Millisecond)

		a.stop()

		// proposal should have been processed
		require.Equal(t, 0, len(a.messages[futureHeight][futureRound].proposals))
	})
	t.Run("RoundChangeEvent triggers rules re-evaluation - votes", func(t *testing.T) {
		futureHeight := uint64(1)
		futureRound := int64(2)
		eventMux := event.NewTypeMuxSilent(nil, log.Root())

		ctrl := gomock.NewController(t)
		defer waitForExpects(t, ctrl)

		coreMock := interfaces.NewMockCore(ctrl)
		backendMock := interfaces.NewMockBackend(ctrl)
		chain := newTestBlockchain()

		backendMock.EXPECT().MessageCh().Return(make(chan events.UnverifiedMessageEvent)).Times(1)
		backendMock.EXPECT().BlockChain().Return(chain).AnyTimes()
		backendMock.EXPECT().Post(gomock.Any()).Times(4)
		backendMock.EXPECT().Subscribe(gomock.Any(), gomock.Any(), gomock.Any()).Return(eventMux.Subscribe(events.RoundChangeEvent{})).Times(1)
		coreMock.EXPECT().Height().Return(common.Big1).AnyTimes()
		coreMock.EXPECT().VotesPowerFor(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(message.NewAggregatedPower()).Times(2)
		coreMock.EXPECT().VotesPower(gomock.Any(), gomock.Any(), gomock.Any()).Return(message.NewAggregatedPower()).Times(2)

		a := &aggregator{
			core:          coreMock,
			backend:       backendMock,
			messages:      make(map[uint64]map[int64]*RoundInfo),
			logger:        log.Root(),
			knownMessages: fixsizecache.New[common.Hash, bool](numBuckets, numEntries, fixsizecache.HashKey[common.Hash]),
		}

		// need to pass them through a fake prevote to make signature valid
		// save messages that will trigger quorum for *
		genesisCommittee := chain.GetHeaderByNumber(0).Committee
		vote := tweakPrevote(message.NewPrevote(futureRound, futureHeight, common.Hash{0x00}, testSigner, &genesisCommittee[0], csize), testKey.PublicKey())
		a.saveMessage(makeBogusEvent(vote))
		vote = tweakPrevote(message.NewPrevote(futureRound, futureHeight, common.Hash{0x01}, testSigner, &genesisCommittee[1], csize), testKey.PublicKey())
		a.saveMessage(makeBogusEvent(vote))
		vote = tweakPrevote(message.NewPrevote(futureRound, futureHeight, common.Hash{0x02}, testSigner, &genesisCommittee[2], csize), testKey.PublicKey())
		a.saveMessage(makeBogusEvent(vote))
		vote = tweakPrevote(message.NewPrevote(futureRound, futureHeight, common.Hash{0x03}, testSigner, &genesisCommittee[3], csize), testKey.PublicKey())
		a.saveMessage(makeBogusEvent(vote))
		require.Equal(t, 4, len(a.messages[futureHeight][futureRound].prevotes))

		a.start(context.TODO())

		eventMux.Post(events.RoundChangeEvent{Height: futureHeight, Round: futureRound})

		// give time to the aggregator loop to process the votes
		time.Sleep(10 * time.Millisecond)

		a.stop()

		// votes should have been processed
		require.Equal(t, 0, len(a.messages[futureHeight][futureRound].prevotes))

	})
	t.Run("PowerChangeEvent triggers rules re-evaluation", func(t *testing.T) {
		height := uint64(1)
		round := int64(2)
		code := message.PrevoteCode
		value := common.Hash{0xca, 0xfe}
		eventMux := event.NewTypeMuxSilent(nil, log.Root())

		ctrl := gomock.NewController(t)
		defer waitForExpects(t, ctrl)

		coreMock := interfaces.NewMockCore(ctrl)
		backendMock := interfaces.NewMockBackend(ctrl)
		chain := newTestBlockchain()

		backendMock.EXPECT().MessageCh().Return(make(chan events.UnverifiedMessageEvent)).Times(1)
		coreMock.EXPECT().VotesPowerFor(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(message.NewAggregatedPower()).Times(2)
		coreMock.EXPECT().VotesPower(gomock.Any(), gomock.Any(), gomock.Any()).Return(message.NewAggregatedPower()).Times(1)
		coreMock.EXPECT().Height().Return(common.Big1).AnyTimes()
		backendMock.EXPECT().Subscribe(gomock.Any(), gomock.Any(), gomock.Any()).Return(eventMux.Subscribe(events.PowerChangeEvent{})).Times(1)
		backendMock.EXPECT().BlockChain().Return(chain).AnyTimes()

		a := &aggregator{
			core:          coreMock,
			backend:       backendMock,
			messages:      make(map[uint64]map[int64]*RoundInfo),
			logger:        log.Root(),
			knownMessages: fixsizecache.New[common.Hash, bool](numBuckets, numEntries, fixsizecache.HashKey[common.Hash]),
		}

		// save a prevote carrying quorum for same (h,r,c,v) as the power change. It should get processed due to the power change
		genesisCommittee := chain.GetHeaderByNumber(0).Committee
		vote := message.NewPrevote(round, height, value, testSigner, &genesisCommittee[0], csize)
		for i := 1; i < len(genesisCommittee); i++ {
			vote.Signers().Increment(&genesisCommittee[i])
		}
		a.saveMessage(makeBogusEvent(vote))
		require.Equal(t, 1, len(a.messages[height][round].prevotes[value]))

		a.start(context.TODO())

		eventMux.Post(events.PowerChangeEvent{Height: height, Round: round, Code: code, Value: value})

		// give time to the aggregator loop to process the prevote
		time.Sleep(10 * time.Millisecond)

		a.stop()

		require.Equal(t, 0, len(a.messages[height][round].prevotes[value]))
	})
	t.Run("FuturePowerChangeEvent triggers re-evaluation of possible round skip", func(t *testing.T) {
		height := uint64(1)
		round := int64(2)
		value := common.Hash{0xca, 0xfe}
		eventMux := event.NewTypeMuxSilent(nil, log.Root())

		ctrl := gomock.NewController(t)
		defer waitForExpects(t, ctrl)

		coreMock := interfaces.NewMockCore(ctrl)
		backendMock := interfaces.NewMockBackend(ctrl)
		chain := newTestBlockchain()

		backendMock.EXPECT().MessageCh().Return(make(chan events.UnverifiedMessageEvent)).Times(1)
		coreMock.EXPECT().Height().Return(common.Big1).AnyTimes()
		backendMock.EXPECT().BlockChain().Return(chain).AnyTimes()
		aggregatedPower := message.NewAggregatedPower()
		aggregatedPower.Set(0, chain.GetHeaderByNumber(0).TotalVotingPower())
		coreMock.EXPECT().Power(gomock.Any(), gomock.Any()).Return(aggregatedPower).Times(1)
		backendMock.EXPECT().Post(gomock.Any()).Times(1)
		backendMock.EXPECT().Subscribe(gomock.Any(), gomock.Any(), gomock.Any()).Return(eventMux.Subscribe(events.FuturePowerChangeEvent{})).Times(1)

		a := &aggregator{
			core:          coreMock,
			backend:       backendMock,
			messages:      make(map[uint64]map[int64]*RoundInfo),
			logger:        log.Root(),
			knownMessages: fixsizecache.New[common.Hash, bool](numBuckets, numEntries, fixsizecache.HashKey[common.Hash]),
		}

		// should get processed thanks to the FuturePowerChangeEvent
		vote := message.NewPrevote(round, height, value, testSigner, &committee[0], csize)
		a.saveMessage(makeBogusEvent(vote))
		require.Equal(t, 1, len(a.messages[height][round].prevotes[value]))

		a.start(context.TODO())

		eventMux.Post(events.FuturePowerChangeEvent{Height: height, Round: round})

		// give time to the aggregator loop to process the prevote
		time.Sleep(10 * time.Millisecond)

		a.stop()

		require.Nil(t, a.messages[height][round])
	})
}

func newTestBlockchain() *core.BlockChain {
	db := rawdb.NewMemoryDatabase()
	genesis := core.GenesisBlockForTesting(db, common.Address{}, common.Big0)

	chain, err := core.NewBlockChain(db, nil, params.TestChainConfig, ethash.NewFaker(), vm.Config{}, nil, &core.TxSenderCacher{}, nil, backends.NewInternalBackend(nil), log.Root())
	if err != nil {
		panic(err)
	}

	//TODO: figure out why this is needed. If removed the header fetched via `GetHeaderByNumber` doesn't have committee information
	chain.GetHeaderByNumber(0).Committee = genesis.Header().Committee
	return chain
}
