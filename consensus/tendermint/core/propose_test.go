package core

import (
	"context"
	"errors"
	"github.com/autonity/autonity/consensus/tendermint/core/committee"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/helpers"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/messageutils"
	tctypes "github.com/autonity/autonity/consensus/tendermint/core/types"
	"github.com/stretchr/testify/require"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/autonity/autonity/consensus"
	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/log"
)

func TestSendPropose(t *testing.T) {
	t.Run("valid block given, proposal is broadcast", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		messages := messageutils.NewMessagesMap()
		addr := common.HexToAddress("0x0123456789")
		block := types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(1),
		})

		curRoundMessages := messages.GetOrCreate(1)
		validRound := int64(1)
		logger := log.New("backend", "test", "id", 0)
		proposalBlock := messageutils.NewProposal(1, big.NewInt(1), validRound, block)
		proposal, err := messageutils.Encode(proposalBlock)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		expectedMsg := &messageutils.Message{
			Code:          messageutils.MsgProposal,
			Msg:           proposal,
			Address:       addr,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
		}

		payloadNoSig, err := expectedMsg.PayloadNoSig()
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		payload := expectedMsg.GetPayload()

		testCommittee := types.Committee{
			types.CommitteeMember{
				Address:     addr,
				VotingPower: big.NewInt(1)},
		}

		valSet, err := committee.NewRoundRobinSet(testCommittee, testCommittee[0].Address)
		if err != nil {
			t.Error(err)
		}

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().SetProposedBlockHash(block.Hash())
		backendMock.EXPECT().Sign(payloadNoSig).Return([]byte{0x1}, nil)
		backendMock.EXPECT().Broadcast(gomock.Any(), gomock.Any(), payload)

		c := &Core{
			address:          addr,
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
		c.proposer.SendProposal(context.Background(), block)
	})
}

func TestHandleProposal(t *testing.T) {
	t.Run("old proposal given, error returned", func(t *testing.T) {
		addr := common.HexToAddress("0x0123456789")
		block := types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(1),
		})

		messages := messageutils.NewMessagesMap()
		curRoundMessages := messages.GetOrCreate(2)

		logger := log.New("backend", "test", "id", 0)

		proposalBlock := messageutils.NewProposal(1, big.NewInt(1), 1, block)
		proposal, err := messageutils.Encode(proposalBlock)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		msg := &messageutils.Message{
			Code:          messageutils.MsgProposal,
			Msg:           proposal,
			Address:       addr,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
		}

		c := &Core{
			address:          addr,
			messages:         messages,
			curRoundMessages: curRoundMessages,
			logger:           logger,
			round:            2,
			height:           big.NewInt(1),
		}

		c.SetDefaultHandlers()
		err = c.proposer.HandleProposal(context.Background(), msg)
		if err != constants.ErrOldRoundMessage {
			t.Fatalf("Expected %v, got %v", constants.ErrOldRoundMessage, err)
		}
	})

	t.Run("msg from non-proposer given, error returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		messages := messageutils.NewMessagesMap()
		addr := common.HexToAddress("0x0123456789")
		block := types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(1),
		})

		curRoundMessages := messages.GetOrCreate(2)

		logger := log.New("backend", "test", "id", 0)
		proposalBlock := messageutils.NewProposal(2, big.NewInt(1), 1, block)
		proposal, err := messageutils.Encode(proposalBlock)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		msg := &messageutils.Message{
			Code:          messageutils.MsgProposal,
			Msg:           proposal,
			Address:       addr,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
		}

		testCommittee, _ := helpers.GenerateCommittee(3)
		testCommittee = append(testCommittee, types.CommitteeMember{Address: addr, VotingPower: big.NewInt(1)})

		valSet, err := committee.NewRoundRobinSet(testCommittee, testCommittee[1].Address)
		if err != nil {
			t.Error(err)
		}

		c := &Core{
			address:          addr,
			messages:         messages,
			curRoundMessages: curRoundMessages,
			logger:           logger,
			round:            2,
			height:           big.NewInt(1),
			committee:        valSet,
		}

		c.SetDefaultHandlers()
		err = c.proposer.HandleProposal(context.Background(), msg)
		if err != constants.ErrNotFromProposer {
			t.Fatalf("Expected %v, got %v", constants.ErrNotFromProposer, err)
		}
	})

	t.Run("unverified proposal given, error returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		addr := common.HexToAddress("0x0123456789")

		block := types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(1),
		})

		messageMap := messageutils.NewMessagesMap()
		curRoundMessages := messageMap.GetOrCreate(2)

		logger := log.New("backend", "test", "id", 0)
		proposalBlock := messageutils.NewProposal(2, big.NewInt(1), 1, block)
		proposal, err := messageutils.Encode(proposalBlock)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		msg := &messageutils.Message{
			Code:          messageutils.MsgProposal,
			Msg:           proposal,
			Address:       addr,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
			Power:         common.Big1,
		}

		testCommittee := types.Committee{
			types.CommitteeMember{Address: addr, VotingPower: big.NewInt(1)},
		}

		valSet, err := committee.NewRoundRobinSet(testCommittee, testCommittee[0].Address)
		if err != nil {
			t.Error(err)
		}

		var decProposal messageutils.Proposal
		if decErr := msg.Decode(&decProposal); decErr != nil {
			t.Fatalf("Expected <nil>, got %v", decErr)
		}

		var prevote = messageutils.Vote{
			Round:             2,
			Height:            big.NewInt(1),
			ProposedBlockHash: common.Hash{},
		}

		encodedVote, err := messageutils.Encode(&prevote)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		preVoteMsg := &messageutils.Message{
			Code:          messageutils.MsgPrevote,
			Msg:           encodedVote,
			Address:       addr,
			CommittedSeal: []byte{},
			Power:         common.Big1,
		}

		payloadNoSig, err := preVoteMsg.PayloadNoSig()
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		payload := preVoteMsg.GetPayload()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().VerifyProposal(gomock.Any()).Return(time.Nanosecond, errors.New("bad block"))
		backendMock.EXPECT().Sign(payloadNoSig)
		backendMock.EXPECT().Broadcast(gomock.Any(), gomock.Any(), payload)
		backendMock.EXPECT().Post(gomock.Any()).Times(0)

		c := &Core{
			address:          addr,
			backend:          backendMock,
			messages:         messageMap,
			curRoundMessages: curRoundMessages,
			logger:           logger,
			proposeTimeout:   tctypes.NewTimeout(tctypes.Propose, logger),
			committee:        valSet,
			round:            2,
			height:           big.NewInt(1),
		}

		c.SetDefaultHandlers()
		err = c.proposer.HandleProposal(context.Background(), msg)
		if err == nil {
			t.Fatalf("Expected non nil error, got %v", err)
		}
	})

	t.Run("future proposal given, backlog event posted", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		addr := common.HexToAddress("0x0123456789")

		block := types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(1),
		})

		messageMap := messageutils.NewMessagesMap()
		curRoundMessages := messageMap.GetOrCreate(2)

		logger := log.New("backend", "test", "id", 0)
		proposalBlock := messageutils.NewProposal(2, big.NewInt(1), 1, block)
		proposal, err := messageutils.Encode(proposalBlock)
		assert.NoError(t, err)

		msg := &messageutils.Message{
			Code:          messageutils.MsgProposal,
			Msg:           proposal,
			Address:       addr,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
			Power:         common.Big1,
		}

		testCommittee := types.Committee{
			types.CommitteeMember{Address: addr, VotingPower: big.NewInt(1)},
		}

		valSet, err := committee.NewRoundRobinSet(testCommittee, testCommittee[0].Address)
		assert.NoError(t, err)
		backendMock := interfaces.NewMockBackend(ctrl)
		const eventPostingDelay = time.Second
		backendMock.EXPECT().VerifyProposal(gomock.Any()).Return(eventPostingDelay, consensus.ErrFutureBlock)
		event := backlogEvent{
			msg: msg,
		}

		backendMock.EXPECT().Post(event).Times(1)
		c := &Core{
			address:          addr,
			backend:          backendMock,
			messages:         messageMap,
			curRoundMessages: curRoundMessages,
			logger:           logger,
			proposeTimeout:   tctypes.NewTimeout(tctypes.Propose, logger),
			committee:        valSet,
			round:            2,
			height:           big.NewInt(1),
		}

		c.SetDefaultHandlers()
		err = c.proposer.HandleProposal(context.Background(), msg)
		assert.Error(t, err)
		// We wait here for at least the delay "eventPostingDelay" returned by VerifyProposal :
		// We expect above that a backlog event containing the future proposal message will be posted
		// after this amount of time. This being done asynchrounously it is necessary to pause the main thread.
		<-time.NewTimer(2 * eventPostingDelay).C
	})

	t.Run("valid proposal given, no error returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		addr := common.HexToAddress("0x0123456789")
		block := types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(1),
		})

		messages := messageutils.NewMessagesMap()
		curRoundMessages := messages.GetOrCreate(2)

		logger := log.New("backend", "test", "id", 0)
		proposalBlock := messageutils.NewProposal(2, big.NewInt(1), 2, block)
		proposal, err := messageutils.Encode(proposalBlock)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		msg := &messageutils.Message{
			Code:          messageutils.MsgProposal,
			Msg:           proposal,
			Address:       addr,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
			Power:         common.Big1,
		}

		testCommittee := types.Committee{
			types.CommitteeMember{Address: addr, VotingPower: big.NewInt(1)},
		}

		valSet, err := committee.NewRoundRobinSet(testCommittee, testCommittee[0].Address)
		if err != nil {
			t.Error(err)
		}

		var decProposal messageutils.Proposal
		if decErr := msg.Decode(&decProposal); decErr != nil {
			t.Fatalf("Expected <nil>, got %v", decErr)
		}

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().VerifyProposal(*decProposal.ProposalBlock)

		c := &Core{
			address:          addr,
			backend:          backendMock,
			messages:         messages,
			curRoundMessages: curRoundMessages,
			logger:           logger,
			round:            2,
			height:           big.NewInt(1),
			proposeTimeout:   tctypes.NewTimeout(tctypes.Propose, logger),
			committee:        valSet,
		}

		c.SetDefaultHandlers()
		err = c.proposer.HandleProposal(context.Background(), msg)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		if !reflect.DeepEqual(curRoundMessages.ProposalMsg, msg) {
			t.Fatalf("%v not equal to  %v", curRoundMessages.ProposalMsg, msg)
		}
	})

	t.Run("valid proposal given and already a quorum of precommits received for it, commit", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		committeeSet, keys := helpers.NewTestCommitteeSetWithKeys(4)
		logger := log.New("backend", "test", "id", 0)
		proposer, err := committeeSet.GetByIndex(3)
		assert.NoError(t, err)

		proposalBlock := types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(1),
		})

		messages := messageutils.NewMessagesMap()
		curRoundMessages := messages.GetOrCreate(2)

		proposalMsg := messageutils.NewProposal(2, big.NewInt(1), 2, proposalBlock)
		proposal, err := messageutils.Encode(proposalMsg)
		assert.NoError(t, err)

		backendMock := interfaces.NewMockBackend(ctrl)

		c := &Core{
			address:          common.HexToAddress("0x0123456789"),
			backend:          backendMock,
			messages:         messages,
			curRoundMessages: curRoundMessages,
			logger:           logger,
			round:            2,
			height:           big.NewInt(1),
			proposeTimeout:   tctypes.NewTimeout(tctypes.Propose, logger),
			precommitTimeout: tctypes.NewTimeout(tctypes.Precommit, logger),
			committee:        committeeSet,
			step:             tctypes.Precommit,
		}
		c.SetDefaultHandlers()
		defer c.proposeTimeout.StopTimer()   // nolint: errcheck
		defer c.precommitTimeout.StopTimer() // nolint: errcheck

		// Handle a quorum of precommits for this proposal
		for i := 0; i < 3; i++ {
			val, _ := committeeSet.GetByIndex(i)
			precommitMsg, err := preparePrecommitMsg(proposalBlock.Hash(), 2, 1, keys, val)
			assert.NoError(t, err)
			err = c.precommiter.HandlePrecommit(context.Background(), precommitMsg)
			assert.NoError(t, err)
		}

		msg := &messageutils.Message{
			Code:          messageutils.MsgProposal,
			Msg:           proposal,
			Address:       proposer.Address,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
			Power:         common.Big1,
		}
		var decProposal messageutils.Proposal
		err = msg.Decode(&decProposal)
		assert.NoError(t, err)
		backendMock.EXPECT().VerifyProposal(*decProposal.ProposalBlock)
		backendMock.EXPECT().Commit(gomock.Any(), int64(2), gomock.Any()).Times(1).Do(func(committedBlock *types.Block, _ int64, _ [][]byte) {
			assert.Equal(t, proposalBlock.Hash(), committedBlock.Hash())
		})

		err = c.proposer.HandleProposal(context.Background(), msg)
		assert.NoError(t, err)
	})

	t.Run("valid proposal given, valid round -1, pre-vote is sent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		addr := common.HexToAddress("0x0123456789")
		block := types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(1),
		})

		messages := messageutils.NewMessagesMap()
		curRoundMessages := messages.GetOrCreate(2)

		logger := log.New("backend", "test", "id", 0)
		proposalBlock := messageutils.NewProposal(2, big.NewInt(1), -1, block)
		proposal, err := messageutils.Encode(proposalBlock)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		msg := &messageutils.Message{
			Code:          messageutils.MsgProposal,
			Msg:           proposal,
			Address:       addr,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
			Power:         common.Big1,
		}

		testCommittee := types.Committee{
			types.CommitteeMember{Address: addr, VotingPower: big.NewInt(1)},
		}

		valSet, err := committee.NewRoundRobinSet(testCommittee, testCommittee[0].Address)
		if err != nil {
			t.Error(err)
		}

		var decProposal messageutils.Proposal
		if decErr := msg.Decode(&decProposal); decErr != nil {
			t.Fatalf("Expected <nil>, got %v", decErr)
		}

		var prevote = messageutils.Vote{
			Round:             2,
			Height:            big.NewInt(1),
			ProposedBlockHash: block.Hash(),
		}

		encodedVote, err := messageutils.Encode(&prevote)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		preVoteMsg := &messageutils.Message{
			Code:          messageutils.MsgPrevote,
			Msg:           encodedVote,
			Address:       addr,
			CommittedSeal: []byte{},
			Power:         common.Big1,
		}

		payloadNoSig, err := preVoteMsg.PayloadNoSig()
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		payload := preVoteMsg.GetPayload()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().VerifyProposal(*decProposal.ProposalBlock)
		backendMock.EXPECT().Sign(payloadNoSig)
		backendMock.EXPECT().Broadcast(gomock.Any(), gomock.Any(), payload)

		c := &Core{
			address:          addr,
			backend:          backendMock,
			messages:         messages,
			curRoundMessages: curRoundMessages,
			round:            2,
			height:           big.NewInt(1),
			lockedValue:      types.NewBlockWithHeader(&types.Header{}),
			lockedRound:      -1,
			logger:           logger,
			proposeTimeout:   tctypes.NewTimeout(tctypes.Propose, logger),
			validRound:       -1,
			committee:        valSet,
		}

		c.SetDefaultHandlers()
		err = c.proposer.HandleProposal(context.Background(), msg)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		if !reflect.DeepEqual(curRoundMessages.ProposalMsg, msg) {
			t.Fatalf("%v not equal to  %v", curRoundMessages.ProposalMsg, msg)
		}
	})

	t.Run("valid proposal given, vr < curR, pre-vote is sent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		addr := common.HexToAddress("0x0123456789")
		block := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(1)})
		messages := messageutils.NewMessagesMap()
		curRoundMessage := messages.GetOrCreate(2)

		logger := log.New("backend", "test", "id", 0)
		proposalBlock := messageutils.NewProposal(2, big.NewInt(1), 1, block)
		proposal, err := messageutils.Encode(proposalBlock)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		msg := &messageutils.Message{
			Code:          messageutils.MsgProposal,
			Msg:           proposal,
			Address:       addr,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
			Power:         common.Big1,
		}

		testCommittee := types.Committee{
			types.CommitteeMember{Address: addr, VotingPower: big.NewInt(1)},
		}

		valSet, err := committee.NewRoundRobinSet(testCommittee, testCommittee[0].Address)
		if err != nil {
			t.Error(err)
		}

		var decProposal messageutils.Proposal
		if decErr := msg.Decode(&decProposal); decErr != nil {
			t.Fatalf("Expected <nil>, got %v", decErr)
		}

		var prevote = messageutils.Vote{
			Round:             2,
			Height:            big.NewInt(1),
			ProposedBlockHash: block.Hash(),
		}

		encodedVote, err := messageutils.Encode(&prevote)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		preVoteMsg := &messageutils.Message{
			Code:          messageutils.MsgPrevote,
			Msg:           encodedVote,
			Address:       addr,
			CommittedSeal: []byte{},
			Power:         common.Big1,
		}

		payloadNoSig, err := preVoteMsg.PayloadNoSig()
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		payload := preVoteMsg.GetPayload()

		messages.GetOrCreate(1).AddPrevote(block.Hash(), *preVoteMsg)

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().VerifyProposal(*decProposal.ProposalBlock)
		backendMock.EXPECT().Sign(payloadNoSig)
		backendMock.EXPECT().Broadcast(gomock.Any(), gomock.Any(), payload)

		c := &Core{
			address:          addr,
			backend:          backendMock,
			curRoundMessages: curRoundMessage,
			messages:         messages,
			lockedRound:      -1,
			round:            2,
			height:           big.NewInt(1),
			lockedValue:      nil,
			logger:           logger,
			proposeTimeout:   tctypes.NewTimeout(tctypes.Propose, logger),
			validRound:       0,
			committee:        valSet,
		}

		c.SetDefaultHandlers()
		err = c.proposer.HandleProposal(context.Background(), msg)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		if !reflect.DeepEqual(curRoundMessage.ProposalMsg, msg) {
			t.Fatalf("%v not equal to  %v", curRoundMessage.ProposalMsg, msg)
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
		defer ctrl.Finish()

		messages := messageutils.NewMessagesMap()
		addr := common.HexToAddress("0x0123456789")
		preBlock := types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(0),
		})
		block := types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(1),
		})

		curRoundMessages := messages.GetOrCreate(1)
		validRound := int64(1)
		logger := log.New("backend", "test", "id", 0)
		proposalBlock := messageutils.NewProposal(1, big.NewInt(1), validRound, block)
		proposal, err := messageutils.Encode(proposalBlock)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		expectedMsg := &messageutils.Message{
			Code:          messageutils.MsgProposal,
			Msg:           proposal,
			Address:       addr,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
		}

		payloadNoSig, err := expectedMsg.PayloadNoSig()
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		payload := expectedMsg.GetPayload()

		testCommittee := types.Committee{
			types.CommitteeMember{
				Address:     addr,
				VotingPower: big.NewInt(1)},
		}

		valSet, err := committee.NewRoundRobinSet(testCommittee, testCommittee[0].Address)
		if err != nil {
			t.Error(err)
		}

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().SetProposedBlockHash(block.Hash())
		backendMock.EXPECT().Sign(payloadNoSig).Return([]byte{0x1}, nil)
		backendMock.EXPECT().Broadcast(gomock.Any(), gomock.Any(), payload)

		c := &Core{
			pendingCandidateBlocks: make(map[uint64]*types.Block),
			address:                addr,
			backend:                backendMock,
			curRoundMessages:       curRoundMessages,
			logger:                 logger,
			messages:               messages,
			round:                  1,
			height:                 big.NewInt(1),
			validRound:             validRound,
			committee:              valSet,
		}
		c.SetDefaultHandlers()
		c.pendingCandidateBlocks[uint64(0)] = preBlock
		c.proposer.HandleNewCandidateBlockMsg(context.Background(), block)
		require.Equal(t, 1, len(c.pendingCandidateBlocks))
		require.Equal(t, uint64(1), c.pendingCandidateBlocks[uint64(1)].Number().Uint64())
	})
}
