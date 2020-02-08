package core

import (
	"context"
	"math/big"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/log"
)

func TestSendPrevote(t *testing.T) {
	t.Run("proposal is empty and send prevote nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		messages := newMessagesMap()
		curRoundMessages := messages.getOrCreate(2)
		backendMock := NewMockBackend(ctrl)
		committeeSet := newTestCommitteeSet(4)
		backendMock.EXPECT().Broadcast(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)
		backendMock.EXPECT().Sign(gomock.Any()).Times(1)
		c := &core{
			logger:           log.New("backend", "test", "id", 0),
			backend:          backendMock,
			messages:         messages,
			curRoundMessages: curRoundMessages,
			round:            2,
			committeeSet:     committeeSet,
			height:           big.NewInt(3),
		}

		c.sendPrevote(context.Background(), true)
	})

	t.Run("valid proposal given, non nil prevote", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		committeSet := newTestCommitteeSet(4)
		logger := log.New("backend", "test", "id", 0)

		proposal := NewProposal(
			1,
			big.NewInt(2),
			1,
			types.NewBlockWithHeader(&types.Header{OriginalHeader: big.NewInt(2)}))

		messages := newMessagesMap()
		curMessages := messages.getOrCreate(2)
		curMessages.SetProposal(proposal, nil, true)

		addr := common.HexToAddress("0x0123456789")

		var preVote = Vote{
			Round:             1,
			Height:            big.NewInt(2),
			ProposedBlockHash: curMessages.GetProposalHash(),
		}

		encodedVote, err := Encode(&preVote)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		expectedMsg := &Message{
			Code:          msgPrevote,
			Msg:           encodedVote,
			Address:       addr,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
		}

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).Return([]byte{0x1}, nil)

		payload, err := expectedMsg.Payload()
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		backendMock.EXPECT().Broadcast(gomock.Any(), gomock.Any(), payload)

		c := &core{
			backend:          backendMock,
			address:          addr,
			logger:           logger,
			height:           big.NewInt(3),
			committeeSet:     committeSet,
			messages:         messages,
			round:            1,
			step:             propose,
			curRoundMessages: curMessages,
		}

		c.sendPrevote(context.Background(), false)
	})
}

func TestHandlePrevote(t *testing.T) {
	t.Run("pre-vote with future height given, error returned", func(t *testing.T) {
		committeeSet := newTestCommitteeSet(4)
		messages := newMessagesMap()
		curRoundMessages := messages.getOrCreate(2)
		addr := common.HexToAddress("0x0123456789")

		var preVote = Vote{
			Round:  2,
			Height: big.NewInt(3),
		}

		encodedVote, err := Encode(&preVote)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		expectedMsg := &Message{
			Code:          msgPrevote,
			Msg:           encodedVote,
			Address:       addr,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
		}

		c := &core{
			address:          addr,
			round:            2,
			height:           big.NewInt(3),
			curRoundMessages: curRoundMessages,
			messages:         messages,
			committeeSet:     committeeSet,
			logger:           log.New("backend", "test", "id", 0),
		}

		err = c.handlePrevote(context.Background(), expectedMsg)
		if err != errFutureHeightMessage {
			t.Fatalf("Expected %v, got %v", errFutureHeightMessage, err)
		}
	})

	t.Run("pre-vote with old height given, pre-vote not added", func(t *testing.T) {
		committeeSet := newTestCommitteeSet(4)
		messages := newMessagesMap()
		curRoundMessages := messages.getOrCreate(2)
		addr := common.HexToAddress("0x0123456789")

		var preVote = Vote{
			Round:             1,
			Height:            big.NewInt(3),
			ProposedBlockHash: common.Hash{},
		}

		encodedVote, err := Encode(&preVote)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		expectedMsg := &Message{
			Code:          msgPrevote,
			Msg:           encodedVote,
			Address:       addr,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
		}

		c := &core{
			address:          addr,
			curRoundMessages: curRoundMessages,
			messages:         messages,
			logger:           log.New("backend", "test", "id", 0),
			committeeSet:     committeeSet,
			round:            1,
			height:           big.NewInt(3),
		}

		err = c.handlePrevote(context.Background(), expectedMsg)
		if err != errOldRoundMessage {
			t.Fatalf("Expected %v, got %v", errOldRoundMessage, err)
		}

		if s := curRoundMessages.PrevotesCount(common.Hash{}); s != 1 {
			t.Fatalf("Expected 1 nil-prevote, but got %d", s)
		}
	})

	t.Run("pre-vote given with no errors, pre-vote added", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		messages := newMessagesMap()
		committeeSet := newTestCommitteeSet(4)
		curRoundMessages := messages.getOrCreate(2)
		logger := log.New("backend", "test", "id", 0)

		proposal := NewProposal(
			1,
			big.NewInt(2),
			1,
			types.NewBlockWithHeader(&types.Header{}))

		curRoundMessages.SetProposal(proposal, nil, true)
		addr := common.HexToAddress("0x0123456789")

		var preVote = Vote{
			Round:             1,
			Height:            big.NewInt(2),
			ProposedBlockHash: curRoundMessages.GetProposalHash(),
		}

		encodedVote, err := Encode(&preVote)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		expectedMsg := &Message{
			Code:          msgPrevote,
			Msg:           encodedVote,
			Address:       addr,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
		}
		backendMock := NewMockBackend(ctrl)
		c := &core{
			address:          addr,
			messages:         messages,
			curRoundMessages: curRoundMessages,
			logger:           logger,
			round:            1,
			height:           big.NewInt(2),
			committeeSet:     committeeSet,
			prevoteTimeout:   newTimeout(prevote, logger),
			backend:          backendMock,
			step:             prevote,
		}

		err = c.handlePrevote(context.Background(), expectedMsg)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		if s := c.curRoundMessages.PrevotesCount(curRoundMessages.GetProposalHash()); s != 1 {
			t.Fatalf("Expected 1 prevote, but got %d", s)
		}
	})

	t.Run("pre-vote given at pre-vote step, non-nil pre-commit sent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		committeeSet := newTestCommitteeSet(4)
		logger := log.New("backend", "test", "id", 0)

		proposal := NewProposal(
			2,
			big.NewInt(3),
			1,
			types.NewBlockWithHeader(&types.Header{}))

		message := newMessagesMap()
		curRoundMessage := message.getOrCreate(2)
		curRoundMessage.SetProposal(proposal, nil, true)

		addr := common.HexToAddress("0x0123456789")

		var preVote = Vote{
			Round:             2,
			Height:            big.NewInt(3),
			ProposedBlockHash: curRoundMessage.GetProposalHash(),
		}

		encodedVote, err := Encode(&preVote)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		expectedMsg := &Message{
			Code:          msgPrevote,
			Msg:           encodedVote,
			Address:       addr,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
		}

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).Return([]byte{0x1}, nil).AnyTimes()

		var precommit = Vote{
			Round:             2,
			Height:            big.NewInt(3),
			ProposedBlockHash: curRoundMessage.GetProposalHash(),
		}

		encodedVote, err = Encode(&precommit)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		msg := &Message{
			Code:          msgPrecommit,
			Msg:           encodedVote,
			Address:       addr,
			CommittedSeal: []byte{0x1},
			Signature:     []byte{0x1},
		}

		payload, err := msg.Payload()
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		backendMock.EXPECT().Broadcast(context.Background(), gomock.Any(), payload)

		c := &core{
			address:          addr,
			backend:          backendMock,
			curRoundMessages: curRoundMessage,
			logger:           logger,
			prevoteTimeout:   newTimeout(prevote, logger),
			committeeSet:     committeeSet,
			round:            2,
			height:           big.NewInt(3),
			step:             prevote,
		}

		err = c.handlePrevote(context.Background(), expectedMsg)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		if s := c.curRoundMessages.PrevotesCount(curRoundMessage.GetProposalHash()); s != 1 {
			t.Fatalf("Expected 1 prevote, but got %d", s)
		}

		if !reflect.DeepEqual(c.validValue, c.curRoundMessages.Proposal().ProposalBlock) {
			t.Fatalf("Expected %v, got %v", c.curRoundMessages.Proposal().ProposalBlock, c.validValue)
		}
	})
	/*
		t.Run("pre-vote given at pre-vote step, nil pre-commit sent", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			curRoundState := NewRoundMessages(big.NewInt(2), big.NewInt(3))
			curRoundState.SetStep(prevote)

			addr := common.HexToAddress("0x0123456789")

			var preVote = Vote{
				Round:             big.NewInt(curRoundState.Round().Int64()),
				Height:            big.NewInt(curRoundState.Height().Int64()),
				ProposedBlockHash: curRoundState.GetProposalHash(),
			}

			encodedVote, err := Encode(&preVote)
			if err != nil {
				t.Fatalf("Expected nil, got %v", err)
			}

			expectedMsg := &Message{
				Code:          msgPrevote,
				Msg:           encodedVote,
				Address:       addr,
				CommittedSeal: []byte{},
				Signature:     []byte{0x1},
			}

			backendMock := NewMockBackend(ctrl)
			backendMock.EXPECT().Sign(gomock.Any()).Return([]byte{0x1}, nil).AnyTimes()

			var precommit = Vote{
				Round:             big.NewInt(curRoundState.Round().Int64()),
				Height:            big.NewInt(curRoundState.Height().Int64()),
				ProposedBlockHash: common.Hash{},
			}

			encodedVote, err = Encode(&precommit)
			if err != nil {
				t.Fatalf("Expected nil, got %v", err)
			}

			msg := &Message{
				Code:          msgPrecommit,
				Msg:           encodedVote,
				Address:       addr,
				CommittedSeal: []byte{0x1},
				Signature:     []byte{0x1},
			}

			payload, err := msg.Payload()
			if err != nil {
				t.Fatalf("Expected nil, got %v", err)
			}

			backendMock.EXPECT().Broadcast(context.Background(), gomock.Any(), payload)

			logger := log.New("backend", "test", "id", 0)

			c := &core{
				address:          addr,
				backend:          backendMock,
				curRoundMessages: curRoundState,
				logger:           logger,
				prevoteTimeout:   newTimeout(prevote, logger),
				committeeSet:     new(validatorSet),
			}

			err = c.handlePrevote(context.Background(), expectedMsg)
			if err != nil {
				t.Fatalf("Expected nil, got %v", err)
			}
		})

		t.Run("pre-vote given at pre-vote step, pre-vote timeout triggered", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			logger := log.New("backend", "test", "id", 0)

			proposal := NewProposal(
				big.NewInt(1),
				big.NewInt(2),
				big.NewInt(1),
				types.NewBlockWithHeader(&types.Header{}))

			addr := common.HexToAddress("0x0123456789")

			curRoundState := NewRoundMessages(big.NewInt(2), big.NewInt(3))
			curRoundState.SetProposal(proposal, nil)
			curRoundState.SetStep(prevote)
			curRoundState.Prevotes.AddVote(addr.Hash(), Message{})

			var preVote = Vote{
				Round:             big.NewInt(curRoundState.Round().Int64()),
				Height:            big.NewInt(curRoundState.Height().Int64()),
				ProposedBlockHash: curRoundState.GetProposalHash(),
			}

			encodedVote, err := Encode(&preVote)
			if err != nil {
				t.Fatalf("Expected nil, got %v", err)
			}

			expectedMsg := &Message{
				Code:          msgPrevote,
				Msg:           encodedVote,
				Address:       addr,
				CommittedSeal: []byte{},
				Signature:     []byte{0x1},
			}

			backendMock := NewMockBackend(ctrl)
			backendMock.EXPECT().Address().AnyTimes().Return(addr)

			c := New(backendMock, nil)
			c.curRoundMessages = curRoundState
			c.prevoteTimeout = newTimeout(prevote, logger)
			c.CommitteeSet() = &validatorSet{
				Set: newTestValidatorSet(2),
			}

			err = c.handlePrevote(context.Background(), expectedMsg)
			if err != nil {
				t.Fatalf("Expected nil, got %v", err)
			}
		})*/
}
