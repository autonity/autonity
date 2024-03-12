package core

import (
	"context"
	"math/big"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/autonity/autonity/common"
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
		backendMock.EXPECT().Broadcast(gomock.Any(), gomock.Any()).Times(1)
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
		signer := makeSigner(keys[member.Address].consensus, member.Address)
		logger := log.New("backend", "test", "id", 0)

		proposal := message.NewPropose(
			1,
			2,
			1,
			types.NewBlockWithHeader(&types.Header{Number: big.NewInt(2)}),
			signer)

		messages := message.NewMap()
		curMessages := messages.GetOrCreate(2)
		curMessages.SetProposal(proposal, true)

		expectedMsg := message.NewPrevote(1, 2, curMessages.ProposalHash(), signer)

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).DoAndReturn(signer)
		backendMock.EXPECT().Broadcast(gomock.Any(), expectedMsg)

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
	committeeSet, keys := NewTestCommitteeSetWithKeys(4)
	member := committeeSet.Committee()[0]
	signer := makeSigner(keys[member.Address].consensus, member.Address)

	t.Run("pre-vote given with no errors, pre-vote added", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		messages := message.NewMap()
		curRoundMessages := messages.GetOrCreate(2)
		proposal := message.NewPropose(
			1,
			2,
			1,
			types.NewBlockWithHeader(&types.Header{}),
			signer)

		curRoundMessages.SetProposal(proposal, true)
		prevote := message.NewPrevote(1, 2, curRoundMessages.ProposalHash(), signer).MustVerify(stubVerifier)

		backendMock := interfaces.NewMockBackend(ctrl)
		c := &Core{
			address:          member.Address,
			messages:         messages,
			curRoundMessages: curRoundMessages,
			logger:           log.Root(),
			round:            1,
			height:           big.NewInt(2),
			committee:        committeeSet,
			prevoteTimeout:   NewTimeout(Prevote, log.Root()),
			backend:          backendMock,
			step:             Prevote,
		}

		c.SetDefaultHandlers()
		err := c.prevoter.HandlePrevote(context.Background(), prevote)
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

		proposal := message.NewPropose(
			2,
			3,
			1,
			types.NewBlockWithHeader(&types.Header{Number: big.NewInt(3)}),
			signer)

		messagesMap := message.NewMap()
		curRoundMessage := messagesMap.GetOrCreate(2)
		curRoundMessage.SetProposal(proposal, true)

		prevote := message.NewPrevote(2, 3, curRoundMessage.ProposalHash(), signer).MustVerify(stubVerifierWithPower(3))
		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).DoAndReturn(signer).AnyTimes()

		precommit := message.NewPrecommit(2, 3, curRoundMessage.ProposalHash(), signer)

		backendMock.EXPECT().Broadcast(gomock.Any(), precommit)

		c := &Core{
			address:          member.Address,
			backend:          backendMock,
			curRoundMessages: curRoundMessage,
			logger:           log.Root(),
			proposeTimeout:   NewTimeout(Propose, log.Root()),
			prevoteTimeout:   NewTimeout(Prevote, log.Root()),
			precommitTimeout: NewTimeout(Precommit, log.Root()),
			committee:        committeeSet,
			round:            2,
			height:           big.NewInt(3),
			step:             Prevote,
		}
		c.SetDefaultHandlers()
		if err := c.prevoter.HandlePrevote(context.Background(), prevote); err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}
		if s := c.curRoundMessages.PrevotesPower(curRoundMessage.ProposalHash()); s.Cmp(common.Big3) != 0 {
			t.Fatalf("Expected 1 prevote, but got %d", s)
		}

		if !reflect.DeepEqual(c.validValue, c.curRoundMessages.Proposal().Block()) {
			t.Fatalf("Expected %v, got %v", c.curRoundMessages.Proposal().Block(), c.validValue)
		}
	})

	t.Run("pre-vote given at pre-vote step, nil pre-commit sent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		messages := message.NewMap()

		member2 := committeeSet.Committee()[1]
		curRoundMessage := messages.GetOrCreate(2)

		expectedMsg := message.NewPrevote(2, 3, common.Hash{}, makeSigner(keys[member2.Address].consensus, member2.Address)).MustVerify(stubVerifierWithPower(3))
		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).DoAndReturn(makeSigner(keys[member2.Address].consensus, member2.Address)).AnyTimes()

		precommit := message.NewPrecommit(2, 3, common.Hash{}, makeSigner(keys[member2.Address].consensus, member2.Address))

		backendMock.EXPECT().Broadcast(gomock.Any(), precommit)

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
			committee:        committeeSet,
			proposeTimeout:   NewTimeout(Propose, logger),
			prevoteTimeout:   NewTimeout(Prevote, logger),
			precommitTimeout: NewTimeout(Precommit, logger),
		}
		c.SetDefaultHandlers()
		if err := c.prevoter.HandlePrevote(context.Background(), expectedMsg); err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}
	})

	// This test hasn't been implemented yet !
	/*
		t.Run("pre-vote given at pre-vote step, pre-vote Timeout triggered", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			member2 := committeeSet.Committee()[1]
			proposal := message.NewPropose(
				2,
				2,
				1,
				types.NewBlockWithHeader(&types.Header{}),
				signer)

			curRoundMessages.SetProposal(proposal, true)

			prevote := message.NewPrevote(2, 2, curRoundMessages.ProposalHash(), makeSigner(keys[member2.Address], member2.Address))
			backendMock := interfaces.NewMockBackend(ctrl)
			backendMock.EXPECT().Address().AnyTimes().Return(member2.Address)
			backendMock.EXPECT().Logger().AnyTimes().Return(log.Root())
			backendMock.EXPECT().Sign(gomock.Any()).DoAndReturn(makeSigner(keys[member2.Address], member2.Address)).AnyTimes()

			c := New(backendMock, nil)
			c.curRoundMessages = curRoundMessages
			c.height = big.NewInt(2)
			c.round = 2
			c.step = Prevote
			c.prevoteTimeout = NewTimeout(Prevote, log.Root())
			c.committee = committeeSet

			err := c.prevoter.HandlePrevote(context.Background(), prevote.MustVerify(stubVerifier))
			if err != nil {
				t.Fatalf("Expected nil, got %v", err)
			}
		})
	*/
}
