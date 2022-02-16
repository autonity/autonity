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
			committee:        committeeSet,
			height:           big.NewInt(3),
		}

		c.sendPrevote(context.Background(), true)
	})

	t.Run("valid proposal given, non nil prevote", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		committeSet := newTestCommitteeSet(4)
		member := committeSet.Committee()[0]
		logger := log.New("backend", "test", "id", 0)

		proposal := NewProposal(
			1,
			big.NewInt(2),
			1,
			types.NewBlockWithHeader(&types.Header{Number: big.NewInt(2)}))

		messages := newMessagesMap()
		curMessages := messages.getOrCreate(2)
		curMessages.SetProposal(proposal, nil, true)

		expectedMsg := createPrevote(t, curMessages.GetProposalHash(), 1, big.NewInt(2), member)

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).Return([]byte{0x1}, nil)

		payload := expectedMsg.Payload()

		backendMock.EXPECT().Broadcast(gomock.Any(), gomock.Any(), payload)

		c := &core{
			backend:          backendMock,
			address:          member.Address,
			logger:           logger,
			height:           big.NewInt(2),
			committee:        committeSet,
			messages:         messages,
			round:            1,
			step:             prevote,
			curRoundMessages: curMessages,
		}

		c.sendPrevote(context.Background(), false)
	})
}

func TestHandlePrevote(t *testing.T) {
	t.Run("pre-vote with future height given, error returned", func(t *testing.T) {
		committeeSet := newTestCommitteeSet(4)
		member := committeeSet.Committee()[0]
		messages := newMessagesMap()
		curRoundMessages := messages.getOrCreate(2)

		expectedMsg := createPrevote(t, common.Hash{}, 2, big.NewInt(4), member)
		c := &core{
			address:          member.Address,
			round:            2,
			height:           big.NewInt(3),
			curRoundMessages: curRoundMessages,
			messages:         messages,
			committee:        committeeSet,
			logger:           log.New("backend", "test", "id", 0),
		}

		err := c.handlePrevote(context.Background(), expectedMsg)
		if err != errFutureHeightMessage {
			t.Fatalf("Expected %v, got %v", errFutureHeightMessage, err)
		}
	})

	t.Run("pre-vote with old height given, pre-vote not added", func(t *testing.T) {
		committeeSet := newTestCommitteeSet(4)
		member := committeeSet.Committee()[0]
		messages := newMessagesMap()
		curRoundMessages := messages.getOrCreate(2)

		expectedMsg := createPrevote(t, common.Hash{}, 1, big.NewInt(1), member)

		c := &core{
			address:          member.Address,
			curRoundMessages: curRoundMessages,
			messages:         messages,
			logger:           log.New("backend", "test", "id", 0),
			committee:        committeeSet,
			round:            1,
			height:           big.NewInt(3),
		}

		err := c.handlePrevote(context.Background(), expectedMsg)
		if err != errOldHeightMessage {
			t.Fatalf("Expected %v, got %v", errOldHeightMessage, err)
		}

		if s := curRoundMessages.PrevotesPower(common.Hash{}); s != 0 {
			t.Fatalf("Expected 0 nil-prevote, but got %d", s)
		}
	})

	t.Run("pre-vote given with no errors, pre-vote added", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		messages := newMessagesMap()
		committeeSet := newTestCommitteeSet(4)
		member := committeeSet.Committee()[0]
		curRoundMessages := messages.getOrCreate(2)
		logger := log.New("backend", "test", "id", 0)

		proposal := NewProposal(
			1,
			big.NewInt(2),
			1,
			types.NewBlockWithHeader(&types.Header{}))

		curRoundMessages.SetProposal(proposal, nil, true)
		expectedMsg := createPrevote(t, curRoundMessages.GetProposalHash(), 1, big.NewInt(2), member)

		backendMock := NewMockBackend(ctrl)
		c := &core{
			address:          member.Address,
			messages:         messages,
			curRoundMessages: curRoundMessages,
			logger:           logger,
			round:            1,
			height:           big.NewInt(2),
			committee:        committeeSet,
			prevoteTimeout:   newTimeout(prevote, logger),
			backend:          backendMock,
			step:             prevote,
		}

		err := c.handlePrevote(context.Background(), expectedMsg)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		if s := c.curRoundMessages.PrevotesPower(curRoundMessages.GetProposalHash()); s != 1 {
			t.Fatalf("Expected 1 prevote, but got %d", s)
		}
	})

	t.Run("pre-vote given at pre-vote step, non-nil pre-commit sent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		committeeSet := newTestCommitteeSet(1)
		logger := log.New("backend", "test", "id", 0)
		member := committeeSet.Committee()[0]
		proposal := NewProposal(
			2,
			big.NewInt(3),
			1,
			types.NewBlockWithHeader(&types.Header{Number: big.NewInt(3)}))

		message := newMessagesMap()
		curRoundMessage := message.getOrCreate(2)
		curRoundMessage.SetProposal(proposal, nil, true)

		expectedMsg := createPrevote(t, curRoundMessage.GetProposalHash(), 2, big.NewInt(3), member)
		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).Return([]byte{0x1}, nil).AnyTimes()

		var precommit = Vote{
			Round:             2,
			Height:            big.NewInt(3),
			ProposedBlockHash: curRoundMessage.GetProposalHash(),
		}

		encodedVote, err := Encode(&precommit)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		msg := &Message{
			Code:          msgPrecommit,
			Msg:           encodedVote,
			Address:       member.Address,
			CommittedSeal: []byte{0x1},
			Signature:     []byte{0x1},
			power:         1,
		}
		payload := msg.Payload()

		backendMock.EXPECT().Broadcast(context.Background(), gomock.Any(), payload)

		c := &core{
			address:          member.Address,
			backend:          backendMock,
			curRoundMessages: curRoundMessage,
			logger:           logger,
			prevoteTimeout:   newTimeout(prevote, logger),
			committee:        committeeSet,
			round:            2,
			height:           big.NewInt(3),
			step:             prevote,
		}

		err = c.handlePrevote(context.Background(), expectedMsg)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		if s := c.curRoundMessages.PrevotesPower(curRoundMessage.GetProposalHash()); s != 1 {
			t.Fatalf("Expected 1 prevote, but got %d", s)
		}

		if !reflect.DeepEqual(c.validValue, c.curRoundMessages.Proposal().ProposalBlock) {
			t.Fatalf("Expected %v, got %v", c.curRoundMessages.Proposal().ProposalBlock, c.validValue)
		}
	})

	t.Run("pre-vote given at pre-vote step, nil pre-commit sent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		committeSet := newTestCommitteeSet(1)
		messages := newMessagesMap()
		member := committeSet.Committee()[0]
		curRoundMessage := messages.getOrCreate(2)

		addr := common.HexToAddress("0x0123456789")

		expectedMsg := createPrevote(t, common.Hash{}, 2, big.NewInt(3), member)
		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).Return([]byte{0x1}, nil).AnyTimes()

		var precommit = Vote{
			Round:             2,
			Height:            big.NewInt(3),
			ProposedBlockHash: common.Hash{},
		}

		encodedVote, err := Encode(&precommit)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		msg := &Message{
			Code:          msgPrecommit,
			Msg:           encodedVote,
			Address:       addr,
			CommittedSeal: []byte{0x1},
			Signature:     []byte{0x1},
			power:         1,
		}

		payload := msg.Payload()

		backendMock.EXPECT().Broadcast(context.Background(), gomock.Any(), payload)

		logger := log.New("backend", "test", "id", 0)

		c := &core{
			address:          addr,
			backend:          backendMock,
			messages:         messages,
			curRoundMessages: curRoundMessage,
			logger:           logger,
			round:            2,
			height:           big.NewInt(3),
			step:             prevote,
			prevoteTimeout:   newTimeout(prevote, logger),
			committee:        committeSet,
		}

		err = c.handlePrevote(context.Background(), expectedMsg)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}
	})

	// This test hasn't been implemented yet !
	t.Run("pre-vote given at pre-vote step, pre-vote timeout triggered", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		committeeSet := newTestCommitteeSet(4)
		messages := newMessagesMap()
		curRoundMessages := messages.getOrCreate(1)

		logger := log.New("backend", "test", "id", 0)

		proposal := NewProposal(
			1,
			big.NewInt(2),
			1,
			types.NewBlockWithHeader(&types.Header{}))

		addr := common.HexToAddress("0x0123456789")

		curRoundMessages.SetProposal(proposal, nil, true)

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
		backendMock.EXPECT().Address().AnyTimes().Return(addr)

		c := New(backendMock)
		c.curRoundMessages = curRoundMessages
		c.height = big.NewInt(2)
		c.round = 1
		c.step = prevote
		c.prevoteTimeout = newTimeout(prevote, logger)
		c.committee = committeeSet

		err = c.handlePrevote(context.Background(), expectedMsg)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}
	})
}

func createPrevote(t *testing.T, proposalHash common.Hash, round int64, height *big.Int, member types.CommitteeMember) *Message {
	var preVote = Vote{
		Round:             round,
		Height:            height,
		ProposedBlockHash: proposalHash,
	}

	encodedVote, err := Encode(&preVote)
	if err != nil {
		t.Fatalf("Expected nil, got %v", err)
		return nil
	}

	expectedMsg := &Message{
		Code:          msgPrevote,
		Msg:           encodedVote,
		Address:       member.Address,
		CommittedSeal: []byte{},
		Signature:     []byte{0x1},
		power:         member.VotingPower.Uint64(),
	}
	return expectedMsg
}
