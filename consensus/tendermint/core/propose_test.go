package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/autonity/autonity/consensus/tendermint/core/committee"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
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

		height := new(big.Int).SetUint64(1)
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

		proposal := generateBlockProposal(1, height, validRound, true, makeSigner(proposerConsensusKey), &testCommittee.Members[0])

		valSet, err := committee.NewRoundRobinSet(&testCommittee, testCommittee.Members[0].Address)
		if err != nil {
			t.Error(err)
		}

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().SetProposedBlockHash(proposal.Block().Hash())
		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().DoAndReturn(makeSigner(proposerConsensusKey))
		backendMock.EXPECT().Broadcast(gomock.Any(), proposal)
		roundState := newTendermintState(log.New(), nil, nil)
		c := &Core{
			address:     proposer,
			backend:     backendMock,
			roundsState: roundState,
			logger:      logger,
			committee:   valSet,
		}
		c.SetHeight(common.Big1)
		c.SetRound(1)
		c.SetValidRoundAndValue(validRound, proposal.Block())

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

		proposal := message.NewPropose(round, height, 1, block, signer, signerMember)

		ctrl := gomock.NewController(t)
		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().VerifyProposal(proposal.Block())
		c := &Core{
			roundsState: newTendermintState(log.New(), nil, nil),
			address:     addr,
			committee:   committeeSet,
			logger:      log.Root(),
			backend:     backendMock,
		}
		c.SetHeight(new(big.Int).SetUint64(height))
		c.SetRound(round)
		c.roundsState.GetOrCreate(round)

		c.SetDefaultHandlers()
		err := c.proposer.HandleProposal(context.Background(), proposal)
		require.NoError(t, err)
		proposal2 := message.NewPropose(round, height, 87, block, signer, signerMember)
		err = c.proposer.HandleProposal(context.Background(), proposal2)
		if !errors.Is(err, constants.ErrAlreadyHaveProposal) {
			t.Fatalf("Expected %v, got %v", constants.ErrAlreadyHaveProposal, err)
		}
	})
	t.Run("old proposal given, error returned", func(t *testing.T) {
		block := types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(1),
		})
		proposal := message.NewPropose(round, height, 1, block, signer, signerMember)
		c := &Core{
			address:     addr,
			roundsState: newTendermintState(log.New(), nil, nil),
			committee:   committeeSet,
			logger:      log.Root(),
		}
		c.SetHeight(new(big.Int).SetUint64(height))
		c.SetRound(round + 1)
		c.roundsState.GetOrCreate(round + 1)

		c.SetDefaultHandlers()
		err := c.proposer.HandleProposal(context.Background(), proposal)
		if !errors.Is(err, constants.ErrOldRoundMessage) {
			t.Fatalf("Expected %v, got %v", constants.ErrOldRoundMessage, err)
		}
	})
	t.Run("msg from non-proposer given, error returned", func(t *testing.T) {
		block := types.NewBlockWithHeader(&types.Header{
			Number: new(big.Int).SetUint64(height),
		})

		logger := log.New("backend", "test", "id", 0)

		nonProposer := &committeeSet.Committee().Members[1]
		proposal := message.NewPropose(round, height, 1, block, makeSigner(keys[nonProposer.Address].consensus), nonProposer)

		c := &Core{
			roundsState: newTendermintState(log.New(), nil, nil),
			address:     addr,
			logger:      logger,
			committee:   committeeSet,
		}
		c.SetHeight(new(big.Int).SetUint64(height))
		c.SetRound(round + 1)
		c.roundsState.GetOrCreate(2)

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
		proposal := newUnverifiedPropose(round, height, 1, block, signer, signerMember)
		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(gomock.Any()).Times(0)

		c := &Core{
			roundsState:    newTendermintState(log.New(), nil, nil),
			address:        addr,
			backend:        backendMock,
			logger:         log.Root(),
			proposeTimeout: NewTimeout(Propose, log.Root()),
			committee:      committeeSet,
		}
		c.SetHeight(new(big.Int).SetUint64(height))
		c.SetRound(round)
		c.roundsState.GetOrCreate(2)

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
		proposal := message.NewPropose(round, height, 1, block, signer, signerMember)
		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().VerifyProposal(gomock.Any()).Return(eventPostingDelay, consensus.ErrFutureTimestampBlock)
		event := backlogMessageEvent{
			msg: proposal,
		}
		backendMock.EXPECT().Post(event).Times(1)
		c := &Core{
			roundsState:    newTendermintState(log.New(), nil, nil),
			address:        addr,
			backend:        backendMock,
			logger:         log.Root(),
			proposeTimeout: NewTimeout(Propose, log.Root()),
			committee:      committeeSet,
		}
		c.SetHeight(new(big.Int).SetUint64(height))
		c.SetRound(round)
		c.roundsState.GetOrCreate(round)

		c.SetDefaultHandlers()
		err := c.proposer.HandleProposal(context.Background(), proposal)
		require.Error(t, err)
		// We wait here for at least the delay "eventPostingDelay" returned by VerifyProposal :
		// We expect above that a backlog event containing the future proposal message will be posted
		// after this amount of time. This being done asynchrounously it is necessary to pause the main thread.
		<-time.NewTimer(2 * eventPostingDelay).C
	})

	t.Run("valid proposal given, no error returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		block := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(1)})

		proposal := message.NewPropose(round, height, 2, block, signer, signerMember)
		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().VerifyProposal(proposal.Block())

		c := &Core{
			address:        addr,
			backend:        backendMock,
			roundsState:    newTendermintState(log.New(), nil, nil),
			logger:         log.Root(),
			proposeTimeout: NewTimeout(Propose, log.Root()),
			committee:      committeeSet,
		}
		c.SetHeight(common.Big1)
		c.SetRound(round)
		curRoundMessages := c.roundsState.GetOrCreate(round)

		c.SetDefaultHandlers()
		err := c.proposer.HandleProposal(context.Background(), proposal)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

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

		proposal := message.NewPropose(2, 1, 2, proposalBlock, makeSigner(keys[proposer.Address].consensus), proposer)
		backendMock := interfaces.NewMockBackend(ctrl)
		c := &Core{
			address:          committeeSet.Committee().Members[0].Address,
			backend:          backendMock,
			roundsState:      newTendermintState(log.New(), nil, nil),
			logger:           logger,
			proposeTimeout:   NewTimeout(Propose, logger),
			prevoteTimeout:   NewTimeout(Prevote, logger),
			precommitTimeout: NewTimeout(Precommit, logger),
			committee:        committeeSet,
		}
		c.SetHeight(common.Big1)
		c.SetRound(2)
		c.UpdateStep(Precommit)
		c.roundsState.GetOrCreate(2)

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

		backendMock.EXPECT().VerifyProposal(proposal.Block())
		backendMock.EXPECT().Commit(gomock.Any(), int64(2), gomock.Any()).Times(1).Do(func(committedBlock *types.Block, _ int64, _ types.AggregateSignature) {
			require.Equal(t, proposalBlock.Hash(), committedBlock.Hash())
		})

		err = c.proposer.HandleProposal(context.Background(), proposal)
		require.NoError(t, err)
	})
	t.Run("valid proposal given, valid round -1, pre-vote is sent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		block := types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(1),
		})

		logger := log.New("backend", "test", "id", 0)
		proposal := message.NewPropose(round, height, -1, block, signer, signerMember)
		prevote := message.NewPrevote(round, height, block.Hash(), signer, signerMember, csize)
		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().VerifyProposal(proposal.Block())
		backendMock.EXPECT().Broadcast(gomock.Any(), prevote)
		backendMock.EXPECT().Sign(gomock.Any()).DoAndReturn(signer)
		c := &Core{
			address:          addr,
			backend:          backendMock,
			roundsState:      newTendermintState(log.New(), nil, nil),
			logger:           logger,
			proposeTimeout:   NewTimeout(Propose, logger),
			prevoteTimeout:   NewTimeout(Prevote, logger),
			precommitTimeout: NewTimeout(Precommit, logger),
			committee:        committeeSet,
		}
		c.SetHeight(common.Big1)
		curRoundMessages := c.roundsState.GetOrCreate(round)
		c.SetRound(round)
		//c.SetLockedRoundAndValue(-1, types.NewBlockWithHeader(&types.Header{}))

		c.SetDefaultHandlers()
		err := c.proposer.HandleProposal(context.Background(), proposal)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		if !reflect.DeepEqual(curRoundMessages.Proposal(), proposal) {
			t.Fatalf("%v not equal to  %v", curRoundMessages.Proposal(), proposal)
		}
	})
	t.Run("valid proposal given, vr < curR with quorum, pre-vote is sent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		block := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(int64(height))})

		proposal := message.NewPropose(round, height, round-1, block, signer, signerMember)

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().VerifyProposal(proposal.Block())
		backendMock.EXPECT().Broadcast(gomock.Any(), message.NewPrevote(round, height, proposal.Block().Hash(), signer, signerMember, csize))
		backendMock.EXPECT().Sign(gomock.Any()).DoAndReturn(signer)

		c := &Core{
			address:          addr,
			backend:          backendMock,
			roundsState:      newTendermintState(log.New(), nil, nil),
			logger:           log.Root(),
			proposeTimeout:   NewTimeout(Propose, log.Root()),
			prevoteTimeout:   NewTimeout(Prevote, log.Root()),
			precommitTimeout: NewTimeout(Precommit, log.Root()),
			committee:        committeeSet,
		}
		c.SetHeight(new(big.Int).SetUint64(height))
		c.SetRound(round)
		c.SetValidRoundAndValue(0, nil)
		for i := 0; i < 3; i++ {
			val, _ := committeeSet.MemberByIndex(i)
			prevote := message.NewPrevote(round-1, height, proposal.Block().Hash(), makeSigner(keys[val.Address].consensus), val, csize)
			c.roundsState.GetOrCreate(round - 1).AddPrevote(prevote)
		}
		curRoundMessage := c.roundsState.GetOrCreate(round)
		c.SetDefaultHandlers()
		err := c.proposer.HandleProposal(context.Background(), proposal)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

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
			roundsState:            newTendermintState(log.New(), nil, nil),
			pendingCandidateBlocks: make(map[uint64]*types.Block),
		}
		c.SetHeight(big.NewInt(11))
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

		preBlock := types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(0),
		})

		validRound := int64(1)
		logger := log.New("backend", "test", "id", 0)

		proposal := generateBlockProposal(1, height, validRound, false, makeSigner(proposerKey), proposer)

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().SetProposedBlockHash(proposal.Block().Hash())
		backendMock.EXPECT().Broadcast(gomock.Any(), proposal)
		backendMock.EXPECT().Sign(gomock.Any()).DoAndReturn(makeSigner(proposerKey))

		c := &Core{
			pendingCandidateBlocks: make(map[uint64]*types.Block),
			address:                proposer.Address,
			backend:                backendMock,
			roundsState:            newTendermintState(log.New(), nil, nil),
			logger:                 logger,
			committee:              committeeSet,
		}
		c.SetHeight(common.Big1)
		c.SetRound(1)
		c.SetValidRoundAndValue(validRound, nil)
		c.SetDefaultHandlers()
		c.pendingCandidateBlocks[uint64(0)] = preBlock
		c.proposer.HandleNewCandidateBlockMsg(context.Background(), proposal.Block())
		require.Equal(t, 1, len(c.pendingCandidateBlocks))
		require.Equal(t, uint64(1), c.pendingCandidateBlocks[uint64(1)].Number().Uint64())
	})
}
