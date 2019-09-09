package core

import (
	"context"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
	"github.com/golang/mock/gomock"
	"math/big"
	"testing"

	"github.com/clearmatics/autonity/log"
)

func TestSendPrevote(t *testing.T) {
	t.Run("proposal is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Broadcast(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		c := &core{
			logger:            log.New("backend", "test", "id", 0),
			backend:           backendMock,
			currentRoundState: NewRoundState(big.NewInt(2), big.NewInt(3)),
		}

		c.sendPrevote(context.Background(), false)
	})

	t.Run("valid proposal given, non nil prevote", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		proposal := NewProposal(
			big.NewInt(1),
			big.NewInt(2),
			big.NewInt(1),
			types.NewBlockWithHeader(&types.Header{}))

		curRoundState := NewRoundState(big.NewInt(2), big.NewInt(3))
		curRoundState.SetProposal(proposal)

		addr := common.HexToAddress("0x0123456789")

		var preVote = Vote{
			Round:             big.NewInt(curRoundState.Round().Int64()),
			Height:            big.NewInt(curRoundState.Height().Int64()),
			ProposedBlockHash: curRoundState.GetCurrentProposalHash(),
		}

		encodedVote, err := Encode(&preVote)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		expectedMsg := &message{
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
			backend:           backendMock,
			address:           addr,
			logger:            log.New("backend", "test", "id", 0),
			valSet:            new(validatorSet),
			currentRoundState: curRoundState,
		}

		c.sendPrevote(context.Background(), false)
	})
}
