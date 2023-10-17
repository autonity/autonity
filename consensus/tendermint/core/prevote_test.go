package core

import (
	"context"
	"errors"

	"go.uber.org/mock/gomock"
	"math/big"
	"reflect"
	"testing"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/log"
)

func TestSendPrevote(t *testing.T) {
	t.Run("proposal is empty and send prevote nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		messages := message.NewMap()
		curRoundMessages := messages.GetOrCreate(2)
		backendMock := interfaces.NewMockBackend(ctrl)
		committeeSet := NewTestCommitteeSet(4)
		backendMock.EXPECT().Broadcast(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)
		backendMock.EXPECT().Sign(gomock.Any()).Times(1)
		c := &Core{
			logger:           log.New("backend", "test", "id", 0),
			backend:          backendMock,
			messages:         messages,
			curRoundMessages: curRoundMessages,
			round:            2,
			committee:        committeeSet,
			height:           big.NewInt(3),
		}

		c.SetDefaultHandlers()
		c.prevoter.SendPrevote(context.Background(), true)
	})

	t.Run("valid proposal given, non nil prevote", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		committeSet, keys := NewTestCommitteeSetWithKeys(4)
		member := committeSet.Committee()[0]
		logger := log.New("backend", "test", "id", 0)

		proposal := message.NewPropose(
			1,
			2,
			1,
			types.NewBlockWithHeader(&types.Header{Number: big.NewInt(2)}),
			makeSigner(keys[member.Address]))

		messages := message.NewMap()
		curMessages := messages.GetOrCreate(2)
		curMessages.SetProposal(proposal, true)

		expectedMsg := message.NewPrevote(1, 2, curMessages.ProposalHash(), makeSigner(keys[member.Address]))

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).Return([]byte{0x1}, nil)

		payload := expectedMsg.Payload()

		backendMock.EXPECT().Broadcast(gomock.Any(), payload)

		c := &Core{
			backend:          backendMock,
			address:          member.Address,
			logger:           logger,
			height:           big.NewInt(2),
			committee:        committeSet,
			messages:         messages,
			round:            1,
			step:             Prevote,
			curRoundMessages: curMessages,
		}

		c.SetDefaultHandlers()
		c.prevoter.SendPrevote(context.Background(), false)
	})
}

func TestHandlePrevote(t *testing.T) {
	t.Run("pre-vote with future height given, error returned", func(t *testing.T) {
		committeeSet, keys := NewTestCommitteeSetWithKeys(4)
		member := committeeSet.Committee()[0]
		messages := message.NewMap()
		curRoundMessages := messages.GetOrCreate(2)

		expectedMsg := message.NewPrevote(2, 4, common.Hash{}, makeSigner(keys[member.Address]))
		c := &Core{
			address:          member.Address,
			round:            2,
			height:           big.NewInt(3),
			curRoundMessages: curRoundMessages,
			messages:         messages,
			committee:        committeeSet,
			logger:           log.New("backend", "test", "id", 0),
		}

		c.SetDefaultHandlers()
		err := c.prevoter.HandlePrevote(context.Background(), expectedMsg)
		if !errors.Is(err, constants.ErrFutureHeightMessage) {
			t.Fatalf("Expected %v, got %v", constants.ErrFutureHeightMessage, err)
		}
	})

	t.Run("pre-vote with old height given, pre-vote not added", func(t *testing.T) {
		committeeSet, keys := NewTestCommitteeSetWithKeys(4)
		member := committeeSet.Committee()[0]
		messages := message.NewMap()
		curRoundMessages := messages.GetOrCreate(2)

		expectedMsg := message.NewPrevote(1, 1, common.Hash{}, makeSigner(keys[member.Address]))
		c := &Core{
			address:          member.Address,
			curRoundMessages: curRoundMessages,
			messages:         messages,
			logger:           log.New("backend", "test", "id", 0),
			committee:        committeeSet,
			round:            1,
			height:           big.NewInt(3),
		}

		c.SetDefaultHandlers()
		err := c.prevoter.HandlePrevote(context.Background(), expectedMsg)
		if !errors.Is(err, constants.ErrOldHeightMessage) {
			t.Fatalf("Expected %v, got %v", constants.ErrOldHeightMessage, err)
		}

		if s := curRoundMessages.PrevotesPower(common.Hash{}); s.Cmp(common.Big0) != 0 {
			t.Fatalf("Expected 0 nil-prevote, but got %d", s)
		}
	})

	t.Run("pre-vote given with no errors, pre-vote added", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		messages := message.NewMap()
		committeeSet, keys := NewTestCommitteeSetWithKeys(4)
		member := committeeSet.Committee()[0]
		curRoundMessages := messages.GetOrCreate(2)
		logger := log.New("backend", "test", "id", 0)

		proposal := message.NewPropose(
			1,
			2,
			1,
			types.NewBlockWithHeader(&types.Header{}),
			makeSigner(keys[member.Address]))

		curRoundMessages.SetProposal(proposal, true)
		expectedMsg := message.NewPrevote(1, 2, curRoundMessages.ProposalHash(), makeSigner(keys[member.Address]))

		backendMock := interfaces.NewMockBackend(ctrl)
		c := &Core{
			address:          member.Address,
			messages:         messages,
			curRoundMessages: curRoundMessages,
			logger:           logger,
			round:            1,
			height:           big.NewInt(2),
			committee:        committeeSet,
			prevoteTimeout:   NewTimeout(Prevote, logger),
			backend:          backendMock,
			step:             Prevote,
		}

		c.SetDefaultHandlers()
		err := c.prevoter.HandlePrevote(context.Background(), expectedMsg)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		if s := c.curRoundMessages.PrevotesPower(curRoundMessages.ProposalHash()); s.Cmp(common.Big1) != 0 {
			t.Fatalf("Expected 1 prevote, but got %d", s)
		}
	})

	t.Run("pre-vote given at pre-vote step, non-nil pre-commit sent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		committeeSet, keys := NewTestCommitteeSetWithKeys(1)
		logger := log.New("backend", "test", "id", 0)
		member := committeeSet.Committee()[0]
		proposal := message.NewPropose(
			2,
			3,
			1,
			types.NewBlockWithHeader(&types.Header{Number: big.NewInt(3)}),
			makeSigner(keys[member.Address]))

		messagesMap := message.NewMap()
		curRoundMessage := messagesMap.GetOrCreate(2)
		curRoundMessage.SetProposal(proposal, true)

		expectedMsg := message.NewPrevote(2, 3, curRoundMessage.ProposalHash(), makeSigner(keys[member.Address]))
		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).Return([]byte{0x1}, nil).AnyTimes()

		precommit := message.NewPrecommit(2, 3, curRoundMessage.ProposalHash(), makeSigner(keys[member.Address]))
		payload := precommit.Payload()

		backendMock.EXPECT().Broadcast(gomock.Any(), payload)

		c := &Core{
			address:          member.Address,
			backend:          backendMock,
			curRoundMessages: curRoundMessage,
			logger:           logger,
			prevoteTimeout:   NewTimeout(Prevote, logger),
			committee:        committeeSet,
			round:            2,
			height:           big.NewInt(3),
			step:             Prevote,
		}
		c.SetDefaultHandlers()
		if err := c.prevoter.HandlePrevote(context.Background(), expectedMsg); err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}
		if s := c.curRoundMessages.PrevotesPower(curRoundMessage.ProposalHash()); s.Cmp(common.Big1) != 0 {
			t.Fatalf("Expected 1 prevote, but got %d", s)
		}

		if !reflect.DeepEqual(c.validValue, c.curRoundMessages.Proposal().Block()) {
			t.Fatalf("Expected %v, got %v", c.curRoundMessages.Proposal().Block(), c.validValue)
		}
	})

	t.Run("pre-vote given at pre-vote step, nil pre-commit sent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		committeSet, keys := NewTestCommitteeSetWithKeys(2)
		messages := message.NewMap()
		member := committeSet.Committee()[0]
		member2 := committeSet.Committee()[1]
		curRoundMessage := messages.GetOrCreate(2)

		expectedMsg := message.NewPrevote(2, 3, common.Hash{}, makeSigner(keys[member.Address]))
		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).Return([]byte{0x1}, nil).AnyTimes()

		precommit := message.NewPrecommit(2, 3, common.Hash{}, makeSigner(keys[member2.Address]))
		payload := precommit.Payload()

		backendMock.EXPECT().Broadcast(gomock.Any(), payload)

		logger := log.New("backend", "test", "id", 0)

		c := &Core{
			address:          member2.Address,
			backend:          backendMock,
			messages:         messages,
			curRoundMessages: curRoundMessage,
			logger:           logger,
			round:            2,
			height:           big.NewInt(3),
			step:             Prevote,
			prevoteTimeout:   NewTimeout(Prevote, logger),
			committee:        committeSet,
		}
		c.SetDefaultHandlers()
		if err := c.prevoter.HandlePrevote(context.Background(), expectedMsg); err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}
	})

	// This test hasn't been implemented yet !
	t.Run("pre-vote given at pre-vote step, pre-vote Timeout triggered", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		committeeSet, keys := NewTestCommitteeSetWithKeys(4)
		messages := message.NewMap()
		member := committeeSet.Committee()[0]
		member2 := committeeSet.Committee()[1]
		curRoundMessages := messages.GetOrCreate(1)

		logger := log.New("backend", "test", "id", 0)

		proposal := message.NewPropose(
			1,
			2,
			1,
			types.NewBlockWithHeader(&types.Header{}),
			makeSigner(keys[member.Address]))

		curRoundMessages.SetProposal(proposal, true)

		prevote := message.NewPrevote(1, 2, curRoundMessages.ProposalHash(), makeSigner(keys[member2.Address]))
		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Address().AnyTimes().Return(member2.Address)
		backendMock.EXPECT().Logger().AnyTimes().Return(log.Root())

		c := New(backendMock, nil)
		c.curRoundMessages = curRoundMessages
		c.height = big.NewInt(2)
		c.round = 1
		c.step = Prevote
		c.prevoteTimeout = NewTimeout(Prevote, logger)
		c.committee = committeeSet

		err := c.prevoter.HandlePrevote(context.Background(), prevote)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}
	})
}
