package core

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/committee"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/log"
)

func TestSendPropose(t *testing.T) {
	t.Run("valid block given, proposal is broadcast", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		// setup proposer cryptographic material
		proposerKey, err := crypto.GenerateKey()
		require.NoError(t, err)
		proposer := crypto.PubkeyToAddress(proposerKey.PublicKey)
		proposerConsensusKey, err := blst.RandKey()
		require.NoError(t, err)

		messages := message.NewMap()
		height := new(big.Int).SetUint64(1)
		curRoundMessages := messages.GetOrCreate(1)
		validRound := int64(1)
		logger := log.New("backend", "test", "id", 0)

		testCommittee := types.Committee{
			Members: []types.CommitteeMember{{
				Address:           proposer,
				VotingPower:       big.NewInt(1),
				ConsensusKey:      proposerConsensusKey.PublicKey(),
				ConsensusKeyBytes: proposerConsensusKey.PublicKey().Marshal(),
				Index:             0,
			}},
		}

		preBlock := types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(0),
		})
		proposal := generateBlockProposal(1, height, validRound, true, makeSigner(proposerConsensusKey), &testCommittee.Members[0], preBlock.Header())

		valSet, err := committee.NewRoundRobinSet(&testCommittee, testCommittee.Members[0].Address)
		if err != nil {
			t.Error(err)
		}

		valSet.SetLastHeader(preBlock.Header())
		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().SetProposedBlockHash(proposal.Block().Hash())
		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().DoAndReturn(makeSigner(proposerConsensusKey))
		backendMock.EXPECT().Broadcast(gomock.Any(), proposal)
		backendMock.EXPECT().HeadBlock().Return(preBlock)

		c := &Core{
			address:          proposer,
			backend:          backendMock,
			curRoundMessages: curRoundMessages,
			logger:           logger,
			messages:         messages,
			round:            1,
			height:           big.NewInt(1),
			validRound:       validRound,
			committee:        valSet,
		}

		c.SetDefaultHandlers()
		c.proposer.SendProposal(context.Background(), proposal.Block())
	})
}

func TestHandleProposal(t *testing.T) {
	committeeSet, keys := NewTestCommitteeSetWithKeys(4)
	addr := committeeSet.Committee().Members[0].Address // round 3 - height 1 proposer
	height := uint64(1)
	round := int64(3)
	signer := makeSigner(keys[addr].consensus)
	signerMember := &committeeSet.Committee().Members[0]
	csize := committeeSet.Committee().Len()

	t.Run("2 proposals received, only first one is accepted", func(t *testing.T) {
		block := types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(1),
		})

		messages := message.NewMap()
		curRoundMessages := messages.GetOrCreate(round)
		proposal := message.NewPropose(round, height, 1, block, signer, signerMember)

		wg := sync.WaitGroup{}
		wg.Add(1)
		ctrl := gomock.NewController(t)
		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().VerifyProposal(proposal.Block())
		backendMock.EXPECT().ProposedBlockHash().Return(common.Hash{})
		backendMock.EXPECT().ProposalVerified(proposal.Block()).Do(func(i any) { wg.Done() })
		backendMock.EXPECT().IsProposalStateCached(proposal.Block().Hash()).Return(false)
		c := &Core{
			address:          addr,
			messages:         messages,
			committee:        committeeSet,
			curRoundMessages: curRoundMessages,
			logger:           log.Root(),
			round:            round,
			height:           new(big.Int).SetUint64(height),
			backend:          backendMock,
		}
		c.SetDefaultHandlers()
		err := c.proposer.HandleProposal(context.Background(), proposal)
		require.NoError(t, err)
		proposal2 := message.NewPropose(round, height, 87, block, signer, signerMember)
		err = c.proposer.HandleProposal(context.Background(), proposal2)
		wg.Wait()
		if !errors.Is(err, constants.ErrAlreadyHaveProposal) {
			t.Fatalf("Expected %v, got %v", constants.ErrAlreadyHaveProposal, err)
		}
	})
	t.Run("old proposal given, error returned", func(t *testing.T) {
		block := types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(1),
		})
		messages := message.NewMap()
		curRoundMessages := messages.GetOrCreate(round + 1)
		proposal := message.NewPropose(round, height, 1, block, signer, signerMember)
		c := &Core{
			address:          addr,
			messages:         messages,
			committee:        committeeSet,
			curRoundMessages: curRoundMessages,
			logger:           log.Root(),
			round:            round + 1,
			height:           new(big.Int).SetUint64(height),
		}
		c.SetDefaultHandlers()

		err := c.proposer.HandleProposal(context.Background(), proposal)
		if !errors.Is(err, constants.ErrOldRoundMessage) {
			t.Fatalf("Expected %v, got %v", constants.ErrOldRoundMessage, err)
		}
	})
	t.Run("msg from non-proposer given, error returned", func(t *testing.T) {
		messages := message.NewMap()
		block := types.NewBlockWithHeader(&types.Header{
			Number: new(big.Int).SetUint64(height),
		})
		curRoundMessages := messages.GetOrCreate(2)

		logger := log.New("backend", "test", "id", 0)

		nonProposer := &committeeSet.Committee().Members[1]
		proposal := message.NewPropose(round, height, 1, block, makeSigner(keys[nonProposer.Address].consensus), nonProposer)

		c := &Core{
			address:          addr,
			messages:         messages,
			curRoundMessages: curRoundMessages,
			logger:           logger,
			round:            round + 1,
			height:           new(big.Int).SetUint64(height),
			committee:        committeeSet,
		}

		c.SetDefaultHandlers()
		err := c.proposer.HandleProposal(context.Background(), proposal)
		if !errors.Is(err, constants.ErrNotFromProposer) {
			t.Fatalf("Expected %v, got %v", constants.ErrNotFromProposer, err)
		}
	})
	t.Run("unverified block proposal given, panic", func(t *testing.T) {
		defer func() {
			err := recover()
			if err == nil {
				t.Errorf("expected panic")
			}
		}()
		ctrl := gomock.NewController(t)

		block := types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(1),
		})
		messageMap := message.NewMap()
		curRoundMessages := messageMap.GetOrCreate(2)
		proposal := newUnverifiedPropose(round, height, 1, block, signer, signerMember)
		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(gomock.Any()).Times(0)

		c := &Core{
			address:          addr,
			backend:          backendMock,
			messages:         messageMap,
			curRoundMessages: curRoundMessages,
			logger:           log.Root(),
			proposeTimeout:   NewTimeout(Propose, log.Root()),
			committee:        committeeSet,
			round:            round,
			height:           new(big.Int).SetUint64(height),
		}
		c.SetDefaultHandlers()
		err := c.proposer.HandleProposal(context.Background(), proposal)
		fmt.Println(err)
	})

	t.Run("future proposal given, backlog event posted", func(t *testing.T) {
		const eventPostingDelay = time.Second
		ctrl := gomock.NewController(t)
		block := types.NewBlockWithHeader(&types.Header{
			Number: new(big.Int).SetUint64(height),
		})
		messageMap := message.NewMap()
		curRoundMessages := messageMap.GetOrCreate(round)
		proposal := message.NewPropose(round, height, 1, block, signer, signerMember)
		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().ProposedBlockHash().Return(common.Hash{})
		backendMock.EXPECT().VerifyProposal(gomock.Any()).Return(eventPostingDelay, consensus.ErrFutureTimestampBlock)
		backendMock.EXPECT().IsProposalStateCached(proposal.Block().Hash()).Return(false)
		event := backlogMessageEvent{
			msg: proposal,
		}
		backendMock.EXPECT().Post(event).Times(1)
		c := &Core{
			address:          addr,
			backend:          backendMock,
			messages:         messageMap,
			curRoundMessages: curRoundMessages,
			logger:           log.Root(),
			proposeTimeout:   NewTimeout(Propose, log.Root()),
			committee:        committeeSet,
			round:            round,
			height:           new(big.Int).SetUint64(height),
		}

		c.SetDefaultHandlers()
		err := c.proposer.HandleProposal(context.Background(), proposal)
		require.Error(t, err)
		// We wait here for at least the delay "eventPostingDelay" returned by VerifyProposal :
		// We expect above that a backlog event containing the future proposal message will be posted
		// after this amount of time. This being done asynchrounously it is necessary to pause the main thread.
		<-time.NewTimer(2 * eventPostingDelay).C
	})

	t.Run("self proposal, no error, no verifyProposal invocation", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		block := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(1)})

		wg := sync.WaitGroup{}
		wg.Add(1)

		messages := message.NewMap()
		curRoundMessages := messages.GetOrCreate(round)
		proposal := message.NewPropose(round, height, 2, block, signer, signerMember)
		backendMock := interfaces.NewMockBackend(ctrl)

		backendMock.EXPECT().ProposalVerified(proposal.Block()).Do(func(i any) { wg.Done() })
		backendMock.EXPECT().ProposedBlockHash().Return(proposal.Block().Hash())

		c := &Core{
			address:          addr,
			backend:          backendMock,
			messages:         messages,
			curRoundMessages: curRoundMessages,
			logger:           log.Root(),
			round:            round,
			height:           big.NewInt(1),
			proposeTimeout:   NewTimeout(Propose, log.Root()),
			committee:        committeeSet,
		}

		c.SetDefaultHandlers()
		err := c.proposer.HandleProposal(context.Background(), proposal)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		wg.Wait()
		if !reflect.DeepEqual(curRoundMessages.Proposal(), proposal) {
			t.Fatalf("%v not equal to  %v", curRoundMessages.Proposal(), proposal)
		}
	})

	t.Run("valid proposal given, no error returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		block := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(1)})

		wg := sync.WaitGroup{}
		wg.Add(1)

		messages := message.NewMap()
		curRoundMessages := messages.GetOrCreate(round)
		proposal := message.NewPropose(round, height, 2, block, signer, signerMember)
		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().VerifyProposal(proposal.Block())
		backendMock.EXPECT().ProposalVerified(proposal.Block()).Do(func(i any) { wg.Done() })
		backendMock.EXPECT().ProposedBlockHash().Return(common.Hash{})
		backendMock.EXPECT().IsProposalStateCached(proposal.Block().Hash()).Return(false)

		c := &Core{
			address:          addr,
			backend:          backendMock,
			messages:         messages,
			curRoundMessages: curRoundMessages,
			logger:           log.Root(),
			round:            round,
			height:           big.NewInt(1),
			proposeTimeout:   NewTimeout(Propose, log.Root()),
			committee:        committeeSet,
		}

		c.SetDefaultHandlers()
		err := c.proposer.HandleProposal(context.Background(), proposal)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		wg.Wait()
		if !reflect.DeepEqual(curRoundMessages.Proposal(), proposal) {
			t.Fatalf("%v not equal to  %v", curRoundMessages.Proposal(), proposal)
		}
	})

	t.Run("valid proposal given and already a quorum of precommits received for it, commit", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		logger := log.New("backend", "test", "id", 0)
		proposer, err := committeeSet.MemberByIndex(3)
		require.NoError(t, err)
		proposalBlock := types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(1),
		})

		messages := message.NewMap()
		curRoundMessages := messages.GetOrCreate(2)

		proposal := message.NewPropose(2, 1, 2, proposalBlock, makeSigner(keys[proposer.Address].consensus), proposer)

		backendMock := interfaces.NewMockBackend(ctrl)

		c := &Core{
			address:          committeeSet.Committee().Members[0].Address,
			backend:          backendMock,
			messages:         messages,
			curRoundMessages: curRoundMessages,
			logger:           logger,
			round:            2,
			height:           big.NewInt(1),
			proposeTimeout:   NewTimeout(Propose, logger),
			prevoteTimeout:   NewTimeout(Prevote, logger),
			precommitTimeout: NewTimeout(Precommit, logger),
			committee:        committeeSet,
			step:             Precommit,
		}
		c.SetDefaultHandlers()
		defer c.proposeTimeout.StopTimer()   // nolint: errcheck
		defer c.precommitTimeout.StopTimer() // nolint: errcheck

		// Handle a quorum of precommits for this proposal
		backendMock.EXPECT().Post(gomock.Any()).MaxTimes(3)
		for i := 0; i < 3; i++ {
			val, _ := committeeSet.MemberByIndex(i)
			precommitMsg := message.NewPrecommit(2, 1, proposalBlock.Hash(), makeSigner(keys[val.Address].consensus), val, csize)
			err = c.precommiter.HandlePrecommit(context.Background(), precommitMsg)
			require.NoError(t, err)
		}
		wg := sync.WaitGroup{}
		wg.Add(1)

		backendMock.EXPECT().VerifyProposal(proposal.Block())
		backendMock.EXPECT().ProposalVerified(proposal.Block()).Do(func(i any) { wg.Done() })
		backendMock.EXPECT().ProposedBlockHash().Return(common.Hash{})
		backendMock.EXPECT().IsProposalStateCached(proposal.Block().Hash()).Return(false)
		backendMock.EXPECT().Commit(gomock.Any(), int64(2), gomock.Any()).Times(1).Do(func(committedBlock *types.Block, _ int64, _ types.AggregateSignature) {
			require.Equal(t, proposalBlock.Hash(), committedBlock.Hash())
		})

		err = c.proposer.HandleProposal(context.Background(), proposal)
		wg.Wait()
		require.NoError(t, err)
	})
	t.Run("valid proposal given, valid round -1, pre-vote is sent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		block := types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(1),
		})
		wg := sync.WaitGroup{}
		wg.Add(1)

		messages := message.NewMap()
		curRoundMessages := messages.GetOrCreate(round)
		logger := log.New("backend", "test", "id", 0)
		proposal := message.NewPropose(round, height, -1, block, signer, signerMember)
		prevote := message.NewPrevote(round, height, block.Hash(), signer, signerMember, csize)
		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().VerifyProposal(proposal.Block())
		backendMock.EXPECT().ProposedBlockHash().Return(common.Hash{})
		backendMock.EXPECT().ProposalVerified(proposal.Block()).Do(func(i any) { wg.Done() })
		backendMock.EXPECT().IsProposalStateCached(proposal.Block().Hash()).Return(false)
		backendMock.EXPECT().Broadcast(gomock.Any(), prevote)
		backendMock.EXPECT().Sign(gomock.Any()).DoAndReturn(signer)
		c := &Core{
			address:          addr,
			backend:          backendMock,
			messages:         messages,
			curRoundMessages: curRoundMessages,
			round:            round,
			height:           big.NewInt(1),
			lockedValue:      types.NewBlockWithHeader(&types.Header{}),
			lockedRound:      -1,
			logger:           logger,
			proposeTimeout:   NewTimeout(Propose, logger),
			prevoteTimeout:   NewTimeout(Prevote, logger),
			precommitTimeout: NewTimeout(Precommit, logger),
			validRound:       -1,
			committee:        committeeSet,
		}

		c.SetDefaultHandlers()
		err := c.proposer.HandleProposal(context.Background(), proposal)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		wg.Wait()
		if !reflect.DeepEqual(curRoundMessages.Proposal(), proposal) {
			t.Fatalf("%v not equal to  %v", curRoundMessages.Proposal(), proposal)
		}
	})
	t.Run("valid proposal given, vr < curR with quorum, pre-vote is sent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		block := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(int64(height))})
		messages := message.NewMap()
		curRoundMessage := messages.GetOrCreate(round)

		proposal := message.NewPropose(round, height, round-1, block, signer, signerMember)

		for i := 0; i < 3; i++ {
			val, _ := committeeSet.MemberByIndex(i)
			prevote := message.NewPrevote(round-1, height, proposal.Block().Hash(), makeSigner(keys[val.Address].consensus), val, csize)
			messages.GetOrCreate(round - 1).AddPrevote(prevote)
		}

		wg := sync.WaitGroup{}
		wg.Add(1)
		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().VerifyProposal(proposal.Block())
		backendMock.EXPECT().ProposedBlockHash().Return(common.Hash{})
		backendMock.EXPECT().ProposalVerified(proposal.Block()).Do(func(i any) { wg.Done() })
		backendMock.EXPECT().IsProposalStateCached(proposal.Block().Hash()).Return(false)
		backendMock.EXPECT().Broadcast(gomock.Any(), message.NewPrevote(round, height, proposal.Block().Hash(), signer, signerMember, csize))
		backendMock.EXPECT().Sign(gomock.Any()).DoAndReturn(signer)

		c := &Core{
			address:          addr,
			backend:          backendMock,
			curRoundMessages: curRoundMessage,
			messages:         messages,
			lockedRound:      -1,
			round:            round,
			height:           new(big.Int).SetUint64(height),
			lockedValue:      nil,
			logger:           log.Root(),
			proposeTimeout:   NewTimeout(Propose, log.Root()),
			prevoteTimeout:   NewTimeout(Prevote, log.Root()),
			precommitTimeout: NewTimeout(Precommit, log.Root()),
			validRound:       0,
			committee:        committeeSet,
		}

		c.SetDefaultHandlers()
		err := c.proposer.HandleProposal(context.Background(), proposal)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		wg.Wait()
		if !reflect.DeepEqual(curRoundMessage.Proposal(), proposal) {
			t.Fatalf("%v not equal to  %v", curRoundMessage.Proposal(), proposal)
		}
	})
}

func TestHandleNewCandidateBlockMsg(t *testing.T) {
	t.Run("invalid block send by miner", func(t *testing.T) {
		c := &Core{
			pendingCandidateBlocks: make(map[uint64]*types.Block),
		}
		c.SetDefaultHandlers()
		c.proposer.HandleNewCandidateBlockMsg(context.Background(), nil)
		require.Equal(t, 0, len(c.pendingCandidateBlocks))
	})

	t.Run("discarding old height candidate blocks", func(t *testing.T) {

		oldHeightCandidate := types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(10),
		})

		c := &Core{
			logger:                 log.New("backend", "test", "id", 0),
			height:                 big.NewInt(11),
			pendingCandidateBlocks: make(map[uint64]*types.Block),
		}
		c.SetDefaultHandlers()
		c.proposer.HandleNewCandidateBlockMsg(context.Background(), oldHeightCandidate)
		require.Equal(t, 0, len(c.pendingCandidateBlocks))
	})

	t.Run("candidate block is the one missed, proposal is broadcast", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		committeeSet, keys := NewTestCommitteeSetWithKeys(1)
		proposer, _ := committeeSet.MemberByIndex(0)
		proposerKey := keys[proposer.Address].consensus
		height := new(big.Int).SetUint64(1)

		messages := message.NewMap()
		preBlock := types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(0),
		})

		curRoundMessages := messages.GetOrCreate(1)
		validRound := int64(1)
		logger := log.New("backend", "test", "id", 0)

		proposal := generateBlockProposal(1, height, validRound, false, makeSigner(proposerKey), proposer, preBlock.Header())

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().SetProposedBlockHash(proposal.Block().Hash())
		backendMock.EXPECT().Broadcast(gomock.Any(), proposal)
		backendMock.EXPECT().Sign(gomock.Any()).DoAndReturn(makeSigner(proposerKey))
		backendMock.EXPECT().HeadBlock().Return(preBlock)

		c := &Core{
			pendingCandidateBlocks: make(map[uint64]*types.Block),
			address:                proposer.Address,
			backend:                backendMock,
			curRoundMessages:       curRoundMessages,
			logger:                 logger,
			messages:               messages,
			round:                  1,
			height:                 big.NewInt(1),
			validRound:             validRound,
			committee:              committeeSet,
		}
		c.SetDefaultHandlers()
		c.pendingCandidateBlocks[uint64(0)] = preBlock
		c.proposer.HandleNewCandidateBlockMsg(context.Background(), proposal.Block())
		require.Equal(t, 1, len(c.pendingCandidateBlocks))
		require.Equal(t, uint64(1), c.pendingCandidateBlocks[uint64(1)].Number().Uint64())
	})
}
