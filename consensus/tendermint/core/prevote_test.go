package core

import (
	"context"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/helpers"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	tcmessage "github.com/autonity/autonity/consensus/tendermint/core/messageutils"
	tctypes "github.com/autonity/autonity/consensus/tendermint/core/types"
	"math/big"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/log"
)

func TestSendPrevote(t *testing.T) {
	t.Run("proposal is empty and send prevote nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		messages := tcmessage.NewMessagesMap()
		curRoundMessages := messages.GetOrCreate(2)
		backendMock := interfaces.NewMockBackend(ctrl)
		committeeSet := helpers.NewTestCommitteeSet(4)
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
		committeSet := helpers.NewTestCommitteeSet(4)
		member := committeSet.Committee()[0]
		logger := log.New("backend", "test", "id", 0)

		proposal := tcmessage.NewProposal(
			1,
			big.NewInt(2),
			1,
			types.NewBlockWithHeader(&types.Header{Number: big.NewInt(2)}))

		messages := tcmessage.NewMessagesMap()
		curMessages := messages.GetOrCreate(2)
		curMessages.SetProposal(proposal, nil, true)

		expectedMsg := tcmessage.CreatePrevote(t, curMessages.GetProposalHash(), 1, big.NewInt(2), member)

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).Return([]byte{0x1}, nil)

		payload := expectedMsg.GetPayload()

		backendMock.EXPECT().Broadcast(gomock.Any(), gomock.Any(), payload)

		c := &Core{
			backend:          backendMock,
			address:          member.Address,
			logger:           logger,
			height:           big.NewInt(2),
			committee:        committeSet,
			messages:         messages,
			round:            1,
			step:             tctypes.Prevote,
			curRoundMessages: curMessages,
		}

		c.SetDefaultHandlers()
		c.prevoter.SendPrevote(context.Background(), false)
	})
}

func TestHandlePrevote(t *testing.T) {
	t.Run("pre-vote with future height given, error returned", func(t *testing.T) {
		committeeSet := helpers.NewTestCommitteeSet(4)
		member := committeeSet.Committee()[0]
		messages := tcmessage.NewMessagesMap()
		curRoundMessages := messages.GetOrCreate(2)

		expectedMsg := tcmessage.CreatePrevote(t, common.Hash{}, 2, big.NewInt(4), member)
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
		if err != constants.ErrFutureHeightMessage {
			t.Fatalf("Expected %v, got %v", constants.ErrFutureHeightMessage, err)
		}
	})

	t.Run("pre-vote with old height given, pre-vote not added", func(t *testing.T) {
		committeeSet := helpers.NewTestCommitteeSet(4)
		member := committeeSet.Committee()[0]
		messages := tcmessage.NewMessagesMap()
		curRoundMessages := messages.GetOrCreate(2)

		expectedMsg := tcmessage.CreatePrevote(t, common.Hash{}, 1, big.NewInt(1), member)

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
		if err != constants.ErrOldHeightMessage {
			t.Fatalf("Expected %v, got %v", constants.ErrOldHeightMessage, err)
		}

		if s := curRoundMessages.PrevotesPower(common.Hash{}); s.Cmp(common.Big0) != 0 {
			t.Fatalf("Expected 0 nil-prevote, but got %d", s)
		}
	})

	t.Run("pre-vote given with no errors, pre-vote added", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		messages := tcmessage.NewMessagesMap()
		committeeSet := helpers.NewTestCommitteeSet(4)
		member := committeeSet.Committee()[0]
		curRoundMessages := messages.GetOrCreate(2)
		logger := log.New("backend", "test", "id", 0)

		proposal := tcmessage.NewProposal(
			1,
			big.NewInt(2),
			1,
			types.NewBlockWithHeader(&types.Header{}))

		curRoundMessages.SetProposal(proposal, nil, true)
		expectedMsg := tcmessage.CreatePrevote(t, curRoundMessages.GetProposalHash(), 1, big.NewInt(2), member)

		backendMock := interfaces.NewMockBackend(ctrl)
		c := &Core{
			address:          member.Address,
			messages:         messages,
			curRoundMessages: curRoundMessages,
			logger:           logger,
			round:            1,
			height:           big.NewInt(2),
			committee:        committeeSet,
			prevoteTimeout:   tctypes.NewTimeout(tctypes.Prevote, logger),
			backend:          backendMock,
			step:             tctypes.Prevote,
		}

		c.SetDefaultHandlers()
		err := c.prevoter.HandlePrevote(context.Background(), expectedMsg)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		if s := c.curRoundMessages.PrevotesPower(curRoundMessages.GetProposalHash()); s.Cmp(common.Big1) != 0 {
			t.Fatalf("Expected 1 prevote, but got %d", s)
		}
	})

	t.Run("pre-vote given at pre-vote step, non-nil pre-commit sent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		committeeSet := helpers.NewTestCommitteeSet(1)
		logger := log.New("backend", "test", "id", 0)
		member := committeeSet.Committee()[0]
		proposal := tcmessage.NewProposal(
			2,
			big.NewInt(3),
			1,
			types.NewBlockWithHeader(&types.Header{Number: big.NewInt(3)}))

		message := tcmessage.NewMessagesMap()
		curRoundMessage := message.GetOrCreate(2)
		curRoundMessage.SetProposal(proposal, nil, true)

		expectedMsg := tcmessage.CreatePrevote(t, curRoundMessage.GetProposalHash(), 2, big.NewInt(3), member)
		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).Return([]byte{0x1}, nil).AnyTimes()

		var precommit = tcmessage.Vote{
			Round:             2,
			Height:            big.NewInt(3),
			ProposedBlockHash: curRoundMessage.GetProposalHash(),
		}

		encodedVote, err := tcmessage.Encode(&precommit)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		msg := &tcmessage.Message{
			Code:          tcmessage.MsgPrecommit,
			Msg:           encodedVote,
			Address:       member.Address,
			CommittedSeal: []byte{0x1},
			Signature:     []byte{0x1},
			Power:         common.Big1,
		}
		payload := msg.GetPayload()

		backendMock.EXPECT().Broadcast(context.Background(), gomock.Any(), payload)

		c := &Core{
			address:          member.Address,
			backend:          backendMock,
			curRoundMessages: curRoundMessage,
			logger:           logger,
			prevoteTimeout:   tctypes.NewTimeout(tctypes.Prevote, logger),
			committee:        committeeSet,
			round:            2,
			height:           big.NewInt(3),
			step:             tctypes.Prevote,
		}

		c.SetDefaultHandlers()
		err = c.prevoter.HandlePrevote(context.Background(), expectedMsg)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		if s := c.curRoundMessages.PrevotesPower(curRoundMessage.GetProposalHash()); s.Cmp(common.Big1) != 0 {
			t.Fatalf("Expected 1 prevote, but got %d", s)
		}

		if !reflect.DeepEqual(c.validValue, c.curRoundMessages.Proposal().ProposalBlock) {
			t.Fatalf("Expected %v, got %v", c.curRoundMessages.Proposal().ProposalBlock, c.validValue)
		}
	})

	t.Run("pre-vote given at pre-vote step, nil pre-commit sent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		committeSet := helpers.NewTestCommitteeSet(1)
		messages := tcmessage.NewMessagesMap()
		member := committeSet.Committee()[0]
		curRoundMessage := messages.GetOrCreate(2)

		addr := common.HexToAddress("0x0123456789")

		expectedMsg := tcmessage.CreatePrevote(t, common.Hash{}, 2, big.NewInt(3), member)
		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).Return([]byte{0x1}, nil).AnyTimes()

		var precommit = tcmessage.Vote{
			Round:             2,
			Height:            big.NewInt(3),
			ProposedBlockHash: common.Hash{},
		}

		encodedVote, err := tcmessage.Encode(&precommit)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		msg := &tcmessage.Message{
			Code:          tcmessage.MsgPrecommit,
			Msg:           encodedVote,
			Address:       addr,
			CommittedSeal: []byte{0x1},
			Signature:     []byte{0x1},
			Power:         common.Big1,
		}

		payload := msg.GetPayload()

		backendMock.EXPECT().Broadcast(context.Background(), gomock.Any(), payload)

		logger := log.New("backend", "test", "id", 0)

		c := &Core{
			address:          addr,
			backend:          backendMock,
			messages:         messages,
			curRoundMessages: curRoundMessage,
			logger:           logger,
			round:            2,
			height:           big.NewInt(3),
			step:             tctypes.Prevote,
			prevoteTimeout:   tctypes.NewTimeout(tctypes.Prevote, logger),
			committee:        committeSet,
		}

		c.SetDefaultHandlers()
		err = c.prevoter.HandlePrevote(context.Background(), expectedMsg)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}
	})

	// This test hasn't been implemented yet !
	t.Run("pre-vote given at pre-vote step, pre-vote Timeout triggered", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		committeeSet := helpers.NewTestCommitteeSet(4)
		messages := tcmessage.NewMessagesMap()
		curRoundMessages := messages.GetOrCreate(1)

		logger := log.New("backend", "test", "id", 0)

		proposal := tcmessage.NewProposal(
			1,
			big.NewInt(2),
			1,
			types.NewBlockWithHeader(&types.Header{}))

		addr := common.HexToAddress("0x0123456789")

		curRoundMessages.SetProposal(proposal, nil, true)

		var preVote = tcmessage.Vote{
			Round:             1,
			Height:            big.NewInt(2),
			ProposedBlockHash: curRoundMessages.GetProposalHash(),
		}

		encodedVote, err := tcmessage.Encode(&preVote)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		expectedMsg := &tcmessage.Message{
			Code:          tcmessage.MsgPrevote,
			Msg:           encodedVote,
			Address:       addr,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
			Power:         common.Big1,
		}

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Address().AnyTimes().Return(addr)

		c := New(backendMock)
		c.curRoundMessages = curRoundMessages
		c.height = big.NewInt(2)
		c.round = 1
		c.step = tctypes.Prevote
		c.prevoteTimeout = tctypes.NewTimeout(tctypes.Prevote, logger)
		c.committee = committeeSet

		err = c.prevoter.HandlePrevote(context.Background(), expectedMsg)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}
	})
}
