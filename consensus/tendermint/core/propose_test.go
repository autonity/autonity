package core

import (
	"context"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/clearmatics/autonity/consensus"
	"github.com/golang/mock/gomock"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/log"
)

func TestSendPropose(t *testing.T) {
	t.Run("valid block given, proposal is broadcast", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		messages := newMessagesMap()
		addr := common.HexToAddress("0x0123456789")
		block := types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(1),
		})

		curRoundMessages := messages.getOrCreate(1)
		validRound := int64(1)
		logger := log.New("backend", "test", "id", 0)
		proposalBlock := NewProposal(1, big.NewInt(1), validRound, block)
		proposal, err := Encode(proposalBlock)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		expectedMsg := &Message{
			Code:          msgProposal,
			Msg:           proposal,
			Address:       addr,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
		}

		payloadNoSig, err := expectedMsg.PayloadNoSig()
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		payload, err := expectedMsg.Payload()
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		testCommittee := types.Committee{
			types.CommitteeMember{
				Address:     addr,
				VotingPower: big.NewInt(1)},
		}

		valSet, err := newRoundRobinSet(testCommittee, testCommittee[0].Address)
		if err != nil {
			t.Error(err)
		}

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().SetProposedBlockHash(block.Hash())
		backendMock.EXPECT().Sign(payloadNoSig).Return([]byte{0x1}, nil)
		backendMock.EXPECT().Broadcast(gomock.Any(), gomock.Any(), payload)

		c := &core{
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

		c.sendProposal(context.Background(), block)
	})
}

func TestHandleProposal(t *testing.T) {
	t.Run("old proposal given, error returned", func(t *testing.T) {
		addr := common.HexToAddress("0x0123456789")
		block := types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(1),
		})

		messages := newMessagesMap()
		curRoundMessages := messages.getOrCreate(2)

		logger := log.New("backend", "test", "id", 0)

		proposalBlock := NewProposal(1, big.NewInt(1), 1, block)
		proposal, err := Encode(proposalBlock)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		msg := &Message{
			Code:          msgProposal,
			Msg:           proposal,
			Address:       addr,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
		}

		c := &core{
			address:          addr,
			messages:         messages,
			curRoundMessages: curRoundMessages,
			logger:           logger,
			round:            2,
			height:           big.NewInt(1),
		}

		err = c.handleProposal(context.Background(), msg)
		if err != errOldRoundMessage {
			t.Fatalf("Expected %v, got %v", errOldRoundMessage, err)
		}
	})

	t.Run("msg from non-proposer given, error returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		messages := newMessagesMap()
		addr := common.HexToAddress("0x0123456789")
		block := types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(1),
		})

		curRoundMessages := messages.getOrCreate(2)

		logger := log.New("backend", "test", "id", 0)
		proposalBlock := NewProposal(2, big.NewInt(1), 1, block)
		proposal, err := Encode(proposalBlock)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		msg := &Message{
			Code:          msgProposal,
			Msg:           proposal,
			Address:       addr,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
		}

		testCommittee, _ := generateCommittee(3)
		testCommittee = append(testCommittee, types.CommitteeMember{Address: addr, VotingPower: big.NewInt(1)})

		valSet, err := newRoundRobinSet(testCommittee, testCommittee[1].Address)
		if err != nil {
			t.Error(err)
		}

		c := &core{
			address:          addr,
			messages:         messages,
			curRoundMessages: curRoundMessages,
			logger:           logger,
			round:            2,
			height:           big.NewInt(1),
			committee:        valSet,
		}

		err = c.handleProposal(context.Background(), msg)
		if err != errNotFromProposer {
			t.Fatalf("Expected %v, got %v", errNotFromProposer, err)
		}
	})

	t.Run("unverified proposal given, error returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		addr := common.HexToAddress("0x0123456789")
		sender := types.CommitteeMember{Address: addr, VotingPower: big.NewInt(1)}
		block := types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(1),
		})

		message := newMessagesMap()
		curRoundMessages := message.getOrCreate(2)

		logger := log.New("backend", "test", "id", 0)
		proposalBlock := NewProposal(2, big.NewInt(1), 1, block)
		proposal, err := Encode(proposalBlock)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		msg := &Message{
			Code:          msgProposal,
			Msg:           proposal,
			Address:       addr,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
		}

		testCommittee := types.Committee{
			types.CommitteeMember{Address: addr, VotingPower: big.NewInt(1)},
		}

		valSet, err := newRoundRobinSet(testCommittee, testCommittee[0].Address)
		if err != nil {
			t.Error(err)
		}

		var decProposal Proposal
		if decErr := msg.Decode(&decProposal); decErr != nil {
			t.Fatalf("Expected <nil>, got %v", decErr)
		}

		var prevote = Vote{
			Round:             2,
			Height:            big.NewInt(1),
			ProposedBlockHash: common.Hash{},
		}

		encodedVote, err := Encode(&prevote)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		preVoteMsg := &Message{
			Code:          msgPrevote,
			Msg:           encodedVote,
			Address:       addr,
			CommittedSeal: []byte{},
		}

		payloadNoSig, err := preVoteMsg.PayloadNoSig()
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		payload, err := preVoteMsg.Payload()
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		event := backlogEvent{
			src: sender,
			msg: msg,
		}

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().VerifyProposal(gomock.Any()).Return(time.Nanosecond, consensus.ErrFutureBlock)
		backendMock.EXPECT().Sign(payloadNoSig)
		backendMock.EXPECT().Broadcast(gomock.Any(), gomock.Any(), payload)
		backendMock.EXPECT().Post(event).AnyTimes()

		c := &core{
			address:          addr,
			backend:          backendMock,
			messages:         message,
			curRoundMessages: curRoundMessages,
			logger:           logger,
			proposeTimeout:   newTimeout(propose, logger),
			committee:        valSet,
			round:            2,
			height:           big.NewInt(1),
		}

		err = c.handleProposal(context.Background(), msg)
		if err != consensus.ErrFutureBlock {
			t.Fatalf("Expected %v, got %v", consensus.ErrFutureBlock, err)
		}
	})

	t.Run("valid proposal given, no error returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		addr := common.HexToAddress("0x0123456789")
		block := types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(1),
		})

		messages := newMessagesMap()
		curRoundMessages := messages.getOrCreate(2)

		logger := log.New("backend", "test", "id", 0)
		proposalBlock := NewProposal(2, big.NewInt(1), 2, block)
		proposal, err := Encode(proposalBlock)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		msg := &Message{
			Code:          msgProposal,
			Msg:           proposal,
			Address:       addr,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
		}

		testCommittee := types.Committee{
			types.CommitteeMember{Address: addr, VotingPower: big.NewInt(1)},
		}

		valSet, err := newRoundRobinSet(testCommittee, testCommittee[0].Address)
		if err != nil {
			t.Error(err)
		}

		var decProposal Proposal
		if decErr := msg.Decode(&decProposal); decErr != nil {
			t.Fatalf("Expected <nil>, got %v", decErr)
		}

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().VerifyProposal(*decProposal.ProposalBlock)

		c := &core{
			address:          addr,
			backend:          backendMock,
			messages:         messages,
			curRoundMessages: curRoundMessages,
			logger:           logger,
			round:            2,
			height:           big.NewInt(1),
			proposeTimeout:   newTimeout(propose, logger),
			committee:        valSet,
		}

		err = c.handleProposal(context.Background(), msg)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		if !reflect.DeepEqual(curRoundMessages.proposalMsg, msg) {
			t.Fatalf("%v not equal to  %v", curRoundMessages.proposalMsg, msg)
		}
	})

	t.Run("valid proposal given, valid round -1, pre-vote is sent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		addr := common.HexToAddress("0x0123456789")
		block := types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(1),
		})

		messages := newMessagesMap()
		curRoundMessages := messages.getOrCreate(2)

		logger := log.New("backend", "test", "id", 0)
		proposalBlock := NewProposal(2, big.NewInt(1), -1, block)
		proposal, err := Encode(proposalBlock)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		msg := &Message{
			Code:          msgProposal,
			Msg:           proposal,
			Address:       addr,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
		}

		testCommittee := types.Committee{
			types.CommitteeMember{Address: addr, VotingPower: big.NewInt(1)},
		}

		valSet, err := newRoundRobinSet(testCommittee, testCommittee[0].Address)
		if err != nil {
			t.Error(err)
		}

		var decProposal Proposal
		if decErr := msg.Decode(&decProposal); decErr != nil {
			t.Fatalf("Expected <nil>, got %v", decErr)
		}

		var prevote = Vote{
			Round:             2,
			Height:            big.NewInt(1),
			ProposedBlockHash: block.Hash(),
		}

		encodedVote, err := Encode(&prevote)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		preVoteMsg := &Message{
			Code:          msgPrevote,
			Msg:           encodedVote,
			Address:       addr,
			CommittedSeal: []byte{},
		}

		payloadNoSig, err := preVoteMsg.PayloadNoSig()
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		payload, err := preVoteMsg.Payload()
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().VerifyProposal(*decProposal.ProposalBlock)
		backendMock.EXPECT().Sign(payloadNoSig)
		backendMock.EXPECT().Broadcast(gomock.Any(), gomock.Any(), payload)

		c := &core{
			address:          addr,
			backend:          backendMock,
			messages:         messages,
			curRoundMessages: curRoundMessages,
			round:            2,
			height:           big.NewInt(1),
			lockedValue:      types.NewBlockWithHeader(&types.Header{}),
			lockedRound:      -1,
			logger:           logger,
			proposeTimeout:   newTimeout(propose, logger),
			validRound:       -1,
			committee:        valSet,
		}

		err = c.handleProposal(context.Background(), msg)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		if !reflect.DeepEqual(curRoundMessages.proposalMsg, msg) {
			t.Fatalf("%v not equal to  %v", curRoundMessages.proposalMsg, msg)
		}
	})

	t.Run("valid proposal given, vr < curR, pre-vote is sent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		addr := common.HexToAddress("0x0123456789")
		block := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(1)})
		messages := newMessagesMap()
		curRoundMessage := messages.getOrCreate(2)

		logger := log.New("backend", "test", "id", 0)
		proposalBlock := NewProposal(2, big.NewInt(1), 1, block)
		proposal, err := Encode(proposalBlock)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		msg := &Message{
			Code:          msgProposal,
			Msg:           proposal,
			Address:       addr,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
		}

		testCommittee := types.Committee{
			types.CommitteeMember{Address: addr, VotingPower: big.NewInt(1)},
		}

		valSet, err := newRoundRobinSet(testCommittee, testCommittee[0].Address)
		if err != nil {
			t.Error(err)
		}

		var decProposal Proposal
		if decErr := msg.Decode(&decProposal); decErr != nil {
			t.Fatalf("Expected <nil>, got %v", decErr)
		}

		var prevote = Vote{
			Round:             2,
			Height:            big.NewInt(1),
			ProposedBlockHash: block.Hash(),
		}

		encodedVote, err := Encode(&prevote)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		preVoteMsg := &Message{
			Code:          msgPrevote,
			Msg:           encodedVote,
			Address:       addr,
			CommittedSeal: []byte{},
		}

		payloadNoSig, err := preVoteMsg.PayloadNoSig()
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		payload, err := preVoteMsg.Payload()
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}
		messages.getOrCreate(1).AddPrevote(block.Hash(), *preVoteMsg)

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().VerifyProposal(*decProposal.ProposalBlock)
		backendMock.EXPECT().Sign(payloadNoSig)
		backendMock.EXPECT().Broadcast(gomock.Any(), gomock.Any(), payload)

		c := &core{
			address:          addr,
			backend:          backendMock,
			curRoundMessages: curRoundMessage,
			messages:         messages,
			lockedRound:      -1,
			round:            2,
			height:           big.NewInt(1),
			lockedValue:      nil,
			logger:           logger,
			proposeTimeout:   newTimeout(propose, logger),
			validRound:       0,
			committee:        valSet,
		}

		err = c.handleProposal(context.Background(), msg)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		if !reflect.DeepEqual(curRoundMessage.proposalMsg, msg) {
			t.Fatalf("%v not equal to  %v", curRoundMessage.proposalMsg, msg)
		}
	})
}
