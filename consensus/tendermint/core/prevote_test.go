package core

import (
	"context"
	"math/big"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/stretchr/testify/require"

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
		// return random signature just to allow prevote encoding
		backendMock.EXPECT().Sign(gomock.Any()).Times(1).Return(testSignature)

		c := &Core{
			logger:           log.New("backend", "test", "id", 0),
			backend:          backendMock,
			messages:         messages,
			curRoundMessages: curRoundMessages,
			round:            2,
			committee:        committeeSet,
			height:           big.NewInt(3),
			address:          committeeSet.Committee().Members[0].Address,
		}

		c.SetDefaultHandlers()
		c.prevoter.SendPrevote(context.Background(), true)
	})

	t.Run("valid proposal given, non nil prevote", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		committeeSet, keys := NewTestCommitteeSetWithKeys(4)
		member := committeeSet.Committee().Members[0]
		signer := makeSigner(keys[member.Address].consensus)
		logger := log.New("backend", "test", "id", 0)
		csize := committeeSet.Committee().Len()

		proposal := message.NewPropose(
			1,
			2,
			1,
			types.NewBlockWithHeader(&types.Header{Number: big.NewInt(2)}),
			signer,
			&member)

		messages := message.NewMap()
		curMessages := messages.GetOrCreate(2)
		curMessages.SetProposal(proposal, true)

		expectedMsg := message.NewPrevote(1, 2, curMessages.ProposalHash(), signer, &member, csize)

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).DoAndReturn(signer)
		backendMock.EXPECT().Broadcast(gomock.Any(), expectedMsg)

		c := &Core{
			backend:          backendMock,
			address:          member.Address,
			logger:           logger,
			height:           big.NewInt(2),
			committee:        committeeSet,
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
	member := committeeSet.Committee().Members[0]
	signer := makeSigner(keys[member.Address].consensus)
	csize := committeeSet.Committee().Len()

	t.Run("pre-vote given with no errors, pre-vote added", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer waitForExpects(ctrl)

		messages := message.NewMap()
		curRoundMessages := messages.GetOrCreate(2)
		proposal := message.NewPropose(
			1,
			2,
			1,
			types.NewBlockWithHeader(&types.Header{}),
			signer,
			&member)

		curRoundMessages.SetProposal(proposal, true)
		prevote := message.NewPrevote(1, 2, curRoundMessages.ProposalHash(), signer, &member, csize)

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(gomock.Any())
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
			signer,
			&member)

		messagesMap := message.NewMap()
		curRoundMessage := messagesMap.GetOrCreate(2)
		curRoundMessage.SetProposal(proposal, true)

		// quorum of prevotes for v
		var prevotes []*message.Prevote
		for i := 0; i < 3; i++ {
			val := committeeSet.Committee().Members[i]
			prevote := message.NewPrevote(2, 3, curRoundMessage.ProposalHash(), makeSigner(keys[val.Address].consensus), &val, csize)
			prevotes = append(prevotes, prevote)
		}
		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).DoAndReturn(signer).AnyTimes()

		precommit := message.NewPrecommit(2, 3, curRoundMessage.ProposalHash(), signer, &member, csize)

		backendMock.EXPECT().Broadcast(gomock.Any(), precommit)
		backendMock.EXPECT().Post(gomock.Any()).MaxTimes(3)

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
		for _, prevote := range prevotes {
			if err := c.prevoter.HandlePrevote(context.Background(), prevote); err != nil {
				t.Fatalf("Expected nil, got %v", err)
			}
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

		member2 := committeeSet.Committee().Members[1]
		curRoundMessage := messages.GetOrCreate(2)

		// quorum of prevotes for nil
		var prevotes []*message.Prevote
		for i := 0; i < 3; i++ {
			val := committeeSet.Committee().Members[i]
			prevote := message.NewPrevote(2, 3, common.Hash{}, makeSigner(keys[val.Address].consensus), &val, csize)
			prevotes = append(prevotes, prevote)
		}
		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).DoAndReturn(makeSigner(keys[member2.Address].consensus)).AnyTimes()
		backendMock.EXPECT().Post(gomock.Any()).MaxTimes(3)

		precommit := message.NewPrecommit(2, 3, common.Hash{}, makeSigner(keys[member2.Address].consensus), &member2, csize)

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
		for _, prevote := range prevotes {
			if err := c.prevoter.HandlePrevote(context.Background(), prevote); err != nil {
				t.Fatalf("Expected nil, got %v", err)
			}
		}
	})

	t.Run("quorum pre-vote given at pre-vote step, pre-vote Timeout triggered", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		messages := message.NewMap()
		curRoundMessages := messages.GetOrCreate(2)
		member2 := committeeSet.Committee().Members[1]

		proposal := message.NewPropose(
			2,
			2,
			1,
			types.NewBlockWithHeader(&types.Header{}),
			signer,
			&member)

		curRoundMessages.SetProposal(proposal, true)

		// 2 prevotes for nil, 1 for v
		var prevotes []*message.Prevote
		val := committeeSet.Committee().Members[0]
		prevotes = append(prevotes, message.NewPrevote(2, 2, curRoundMessages.ProposalHash(), makeSigner(keys[val.Address].consensus), &val, csize))
		val2 := committeeSet.Committee().Members[1]
		prevotes = append(prevotes, message.NewPrevote(2, 2, common.Hash{}, makeSigner(keys[val2.Address].consensus), &val2, csize))
		val3 := committeeSet.Committee().Members[2]
		prevotes = append(prevotes, message.NewPrevote(2, 2, common.Hash{}, makeSigner(keys[val3.Address].consensus), &val3, csize))

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Address().AnyTimes().Return(member2.Address)
		backendMock.EXPECT().Logger().AnyTimes().Return(log.Root())
		backendMock.EXPECT().Post(gomock.Any()).MaxTimes(4)
		backendMock.EXPECT().Sign(gomock.Any()).DoAndReturn(makeSigner(keys[member2.Address].consensus)).AnyTimes()

		c := &Core{
			address:          member2.Address,
			backend:          backendMock,
			messages:         messages,
			curRoundMessages: curRoundMessages,
			logger:           log.Root(),
			round:            2,
			height:           big.NewInt(2),
			step:             Prevote,
			committee:        committeeSet,
			proposeTimeout:   NewTimeout(Propose, log.Root()),
			prevoteTimeout:   NewTimeout(Prevote, log.Root()),
			precommitTimeout: NewTimeout(Precommit, log.Root()),
		}
		c.SetDefaultHandlers()

		for _, prevote := range prevotes {
			err := c.prevoter.HandlePrevote(context.Background(), prevote)
			if err != nil {
				t.Fatalf("Expected nil, got %v", err)
			}
		}

		require.True(t, c.prevoteTimeout.TimerStarted())
	})
}
