package core

import (
	"context"
	"errors"
	"github.com/autonity/autonity/consensus/tendermint/core/committee"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/helpers"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	tctypes "github.com/autonity/autonity/consensus/tendermint/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/rlp"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/autonity/autonity/consensus"
	"github.com/stretchr/testify/assert"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/log"
)

func TestSendPropose(t *testing.T) {
	t.Run("valid block given, proposal is broadcast", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		messages := message.NewMessagesMap()
		proposerKey, err := crypto.GenerateKey()
		require.NoError(t, err)
		proposer := crypto.PubkeyToAddress(proposerKey.PublicKey)
		height := new(big.Int).SetUint64(1)
		curRoundMessages := messages.GetOrCreate(1)
		validRound := int64(1)
		logger := log.New("backend", "test", "id", 0)

		msg, proposal := generateBlockProposal(t, 1, height, validRound, proposer, true, proposerKey)

		assert.NoError(t, err)

		testCommittee := types.Committee{
			types.CommitteeMember{
				Address:     proposer,
				VotingPower: big.NewInt(1)},
		}

		valSet, err := committee.NewRoundRobinSet(testCommittee, testCommittee[0].Address)
		if err != nil {
			t.Error(err)
		}

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().SetProposedBlockHash(proposal.ProposalBlock.Hash())
		backendMock.EXPECT().Sign(gomock.Any()).Times(2).DoAndReturn(signer(proposerKey))
		backendMock.EXPECT().Broadcast(gomock.Any(), msg.Bytes)

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
		c.proposer.SendProposal(context.Background(), proposal.ProposalBlock)
	})
}

func TestHandleProposal(t *testing.T) {
	t.Run("old proposal given, error returned", func(t *testing.T) {
		addr := common.HexToAddress("0x0123456789")
		block := types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(1),
		})

		messages := message.NewMessagesMap()
		curRoundMessages := messages.GetOrCreate(2)

		logger := log.New("backend", "test", "id", 0)

		proposalBlock := message.NewProposal(1, big.NewInt(1), 1, block, dummySigner)
		proposal, err := rlp.EncodeToBytes(proposalBlock)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		msg := &message.Message{
			Code:          consensus.MsgProposal,
			Payload:       proposal,
			ConsensusMsg:  proposalBlock,
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

		messages := message.NewMessagesMap()
		addr := common.HexToAddress("0x0123456789")
		block := types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(1),
		})

		curRoundMessages := messages.GetOrCreate(2)

		logger := log.New("backend", "test", "id", 0)
		proposalBlock := message.NewProposal(2, big.NewInt(1), 1, block, dummySigner)
		proposal, err := rlp.EncodeToBytes(proposalBlock)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		msg := &message.Message{
			Code:          consensus.MsgProposal,
			Payload:       proposal,
			ConsensusMsg:  proposalBlock,
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

		messageMap := message.NewMessagesMap()
		curRoundMessages := messageMap.GetOrCreate(2)

		logger := log.New("backend", "test", "id", 0)
		proposalBlock := message.NewProposal(2, big.NewInt(1), 1, block, dummySigner)
		proposal, err := rlp.EncodeToBytes(proposalBlock)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		msg := &message.Message{
			Code:          consensus.MsgProposal,
			Payload:       proposal,
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

		var decProposal message.Proposal
		if decErr := msg.Decode(&decProposal); decErr != nil {
			t.Fatalf("Expected <nil>, got %v", decErr)
		}

		var prevote = message.Vote{
			Round:             2,
			Height:            big.NewInt(1),
			ProposedBlockHash: common.Hash{},
		}

		encodedVote, err := rlp.EncodeToBytes(&prevote)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		preVoteMsg := &message.Message{
			Code:          consensus.MsgPrevote,
			Payload:       encodedVote,
			ConsensusMsg:  &prevote,
			Address:       addr,
			CommittedSeal: []byte{},
			Power:         common.Big1,
		}

		payloadNoSig, err := preVoteMsg.BytesNoSignature()
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		payload := preVoteMsg.GetBytes()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().VerifyProposal(gomock.Any()).Return(time.Nanosecond, errors.New("bad block"))
		backendMock.EXPECT().Sign(payloadNoSig)
		backendMock.EXPECT().Broadcast(gomock.Any(), payload)
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

		messageMap := message.NewMessagesMap()
		curRoundMessages := messageMap.GetOrCreate(2)

		logger := log.New("backend", "test", "id", 0)
		proposalBlock := message.NewProposal(2, big.NewInt(1), 1, block, dummySigner)
		proposal, err := rlp.EncodeToBytes(proposalBlock)
		assert.NoError(t, err)

		msg := &message.Message{
			Code:          consensus.MsgProposal,
			ConsensusMsg:  proposalBlock,
			Payload:       proposal,
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
		backendMock.EXPECT().VerifyProposal(gomock.Any()).Return(eventPostingDelay, consensus.ErrFutureTimestampBlock)
		event := backlogMessageEvent{
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

		messages := message.NewMessagesMap()
		curRoundMessages := messages.GetOrCreate(2)

		logger := log.New("backend", "test", "id", 0)
		proposalBlock := message.NewProposal(2, big.NewInt(1), 2, block, dummySigner)
		proposal, err := rlp.EncodeToBytes(proposalBlock)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		msg := &message.Message{
			Code:          consensus.MsgProposal,
			Payload:       proposal,
			ConsensusMsg:  proposalBlock,
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

		var decProposal message.Proposal
		if decErr := msg.Decode(&decProposal); decErr != nil {
			t.Fatalf("Expected <nil>, got %v", decErr)
		}

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().VerifyProposal(decProposal.ProposalBlock)

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

		messages := message.NewMessagesMap()
		curRoundMessages := messages.GetOrCreate(2)

		proposalMsg := message.NewProposal(2, big.NewInt(1), 2, proposalBlock, dummySigner)
		proposal, err := rlp.EncodeToBytes(proposalMsg)
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

		msg := &message.Message{
			Code:          consensus.MsgProposal,
			Payload:       proposal,
			ConsensusMsg:  proposalMsg,
			Address:       proposer.Address,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
			Power:         common.Big1,
		}
		var decProposal message.Proposal
		err = msg.Decode(&decProposal)
		assert.NoError(t, err)
		backendMock.EXPECT().VerifyProposal(decProposal.ProposalBlock)
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

		messages := message.NewMessagesMap()
		curRoundMessages := messages.GetOrCreate(2)

		logger := log.New("backend", "test", "id", 0)
		proposalBlock := message.NewProposal(2, big.NewInt(1), -1, block, dummySigner)
		proposal, err := rlp.EncodeToBytes(proposalBlock)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		msg := &message.Message{
			Code:          consensus.MsgProposal,
			Payload:       proposal,
			ConsensusMsg:  proposalBlock,
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

		var decProposal message.Proposal
		if decErr := msg.Decode(&decProposal); decErr != nil {
			t.Fatalf("Expected <nil>, got %v", decErr)
		}

		var prevote = message.Vote{
			Round:             2,
			Height:            big.NewInt(1),
			ProposedBlockHash: block.Hash(),
		}

		encodedVote, err := rlp.EncodeToBytes(&prevote)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		preVoteMsg := &message.Message{
			Code:          consensus.MsgPrevote,
			Payload:       encodedVote,
			ConsensusMsg:  &prevote,
			Address:       addr,
			CommittedSeal: []byte{},
			Power:         common.Big1,
		}

		payloadNoSig, err := preVoteMsg.BytesNoSignature()
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		payload := preVoteMsg.GetBytes()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().VerifyProposal(decProposal.ProposalBlock)
		backendMock.EXPECT().Sign(payloadNoSig)
		backendMock.EXPECT().Broadcast(gomock.Any(), payload)

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
		messages := message.NewMessagesMap()
		curRoundMessage := messages.GetOrCreate(2)

		logger := log.New("backend", "test", "id", 0)
		proposalBlock := message.NewProposal(2, big.NewInt(1), 1, block, dummySigner)
		proposal, err := rlp.EncodeToBytes(proposalBlock)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		msg := &message.Message{
			Code:          consensus.MsgProposal,
			Payload:       proposal,
			ConsensusMsg:  proposalBlock,
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

		var decProposal message.Proposal
		if decErr := msg.Decode(&decProposal); decErr != nil {
			t.Fatalf("Expected <nil>, got %v", decErr)
		}

		var prevote = message.Vote{
			Round:             2,
			Height:            big.NewInt(1),
			ProposedBlockHash: block.Hash(),
		}

		encodedVote, err := rlp.EncodeToBytes(&prevote)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		preVoteMsg := &message.Message{
			Code:          consensus.MsgPrevote,
			Payload:       encodedVote,
			ConsensusMsg:  &prevote,
			Address:       addr,
			CommittedSeal: []byte{},
			Power:         common.Big1,
		}

		payloadNoSig, err := preVoteMsg.BytesNoSignature()
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		payload := preVoteMsg.GetBytes()

		messages.GetOrCreate(1).AddPrevote(block.Hash(), *preVoteMsg)

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().VerifyProposal(decProposal.ProposalBlock)
		backendMock.EXPECT().Sign(payloadNoSig)
		backendMock.EXPECT().Broadcast(gomock.Any(), payload)

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

		proposerKey, err := crypto.GenerateKey()
		require.NoError(t, err)
		proposer := crypto.PubkeyToAddress(proposerKey.PublicKey)
		height := new(big.Int).SetUint64(1)

		messages := message.NewMessagesMap()
		preBlock := types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(0),
		})

		curRoundMessages := messages.GetOrCreate(1)
		validRound := int64(1)
		logger := log.New("backend", "test", "id", 0)

		msg, proposal := generateBlockProposal(t, 1, height, validRound, proposer, false, proposerKey)
		//msgPayload, err := msg.BytesNoSignature()
		require.NoError(t, err)
		msg.Signature = []byte{0x1}

		testCommittee := types.Committee{
			types.CommitteeMember{
				Address:     proposer,
				VotingPower: big.NewInt(1)},
		}

		valSet, err := committee.NewRoundRobinSet(testCommittee, testCommittee[0].Address)
		if err != nil {
			t.Error(err)
		}

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().SetProposedBlockHash(proposal.ProposalBlock.Hash())
		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().DoAndReturn(signer(proposerKey))
		backendMock.EXPECT().Broadcast(gomock.Any(), msg.Bytes)

		c := &Core{
			pendingCandidateBlocks: make(map[uint64]*types.Block),
			address:                proposer,
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
		c.proposer.HandleNewCandidateBlockMsg(context.Background(), proposal.ProposalBlock)
		require.Equal(t, 1, len(c.pendingCandidateBlocks))
		require.Equal(t, uint64(1), c.pendingCandidateBlocks[uint64(1)].Number().Uint64())
	})
}

func dummySigner(_ []byte) ([]byte, error) {
	return nil, nil
}
