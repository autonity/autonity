package core

import (
	"context"
	"errors"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/committee"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/helpers"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	tctypes "github.com/autonity/autonity/consensus/tendermint/core/types"
	"github.com/autonity/autonity/rlp"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/crypto/secp256k1"
	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/log"
	"github.com/golang/mock/gomock"
)

func TestSendPrecommit(t *testing.T) {
	t.Run("proposal is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Broadcast(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		messages := message.NewMessagesMap()
		c := &Core{
			logger:           log.New("backend", "test", "id", 0),
			backend:          backendMock,
			messages:         messages,
			curRoundMessages: messages.GetOrCreate(0),
			round:            2,
			height:           big.NewInt(3),
		}
		c.SetDefaultHandlers()

		c.precommiter.SendPrecommit(context.Background(), false)
	})

	t.Run("valid proposal given, non nil pre-commit", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		committeeSet, keys := helpers.NewTestCommitteeSetWithKeys(7)
		me, _ := committeeSet.GetByIndex(0)
		logger := log.New("backend", "test", "id", 0)
		val, _ := committeeSet.GetByIndex(2)
		addr := val.Address

		proposal := message.NewProposal(
			1,
			big.NewInt(2),
			1,
			types.NewBlockWithHeader(&types.Header{}),
			signer(keys[me.Address]),
		)

		messages := message.NewMessagesMap()
		curRoundMessages := messages.GetOrCreate(1)
		curRoundMessages.SetProposal(proposal, nil, false)

		var preCommit = message.Vote{
			Round:             1,
			Height:            big.NewInt(2),
			ProposedBlockHash: curRoundMessages.GetProposalHash(),
		}

		encodedVote, err := rlp.EncodeToBytes(&preCommit)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		expectedMsg := &message.Message{
			Code:          consensus.MsgPrecommit,
			Payload:       encodedVote,
			Address:       addr,
			CommittedSeal: []byte{0x1},
			Signature:     []byte{0x1},
		}

		payloadNoSig, err := expectedMsg.BytesNoSignature()
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).Return([]byte{0x1}, nil)
		backendMock.EXPECT().Sign(gomock.Eq(payloadNoSig)).Return([]byte{0x1}, nil)

		payload := expectedMsg.GetBytes()

		backendMock.EXPECT().Broadcast(gomock.Any(), gomock.Any(), payload)

		c := &Core{
			backend:          backendMock,
			address:          addr,
			logger:           logger,
			committee:        committeeSet,
			messages:         messages,
			curRoundMessages: curRoundMessages,
			round:            1,
			height:           big.NewInt(2),
		}
		c.SetDefaultHandlers()
		c.precommiter.SendPrecommit(context.Background(), false)
	})

	t.Run("valid proposal given, nil pre-commit", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		committeeSet, keys := helpers.NewTestCommitteeSetWithKeys(7)
		me, _ := committeeSet.GetByIndex(0)
		logger := log.New("backend", "test", "id", 0)
		val, _ := committeeSet.GetByIndex(2)
		addr := val.Address

		proposal := message.NewProposal(
			1,
			big.NewInt(2),
			1,
			types.NewBlockWithHeader(&types.Header{}),
			signer(keys[me.Address]))

		messages := message.NewMessagesMap()
		curRoundMessages := messages.GetOrCreate(1)
		curRoundMessages.SetProposal(proposal, nil, true)

		var preCommit = message.Vote{
			Round:             1,
			Height:            big.NewInt(2),
			ProposedBlockHash: common.Hash{},
		}

		encodedVote, err := rlp.EncodeToBytes(&preCommit)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		expectedMsg := &message.Message{
			Code:          consensus.MsgPrecommit,
			Payload:       encodedVote,
			Address:       addr,
			CommittedSeal: []byte{0x1},
			Signature:     []byte{0x1},
		}

		payloadNoSig, err := expectedMsg.BytesNoSignature()
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).Return([]byte{0x1}, errors.New("seal sign error"))
		backendMock.EXPECT().Sign(payloadNoSig).Return([]byte{0x1}, nil)

		payload := expectedMsg.GetBytes()

		backendMock.EXPECT().Broadcast(gomock.Any(), gomock.Any(), payload)

		c := &Core{
			backend:          backendMock,
			address:          addr,
			logger:           logger,
			curRoundMessages: curRoundMessages,
			messages:         messages,
			committee:        committeeSet,
			height:           big.NewInt(2),
			round:            1,
		}

		c.SetDefaultHandlers()
		c.precommiter.SendPrecommit(context.Background(), true)
	})

}

func TestHandlePrecommit(t *testing.T) {
	t.Run("pre-commit with future height given, error returned", func(t *testing.T) {

		addr := common.HexToAddress("0x0123456789")
		messages := message.NewMessagesMap()
		curRoundMessages := messages.GetOrCreate(1)
		var preCommit = message.Vote{
			Round:  2,
			Height: big.NewInt(3),
		}

		encodedVote, err := rlp.EncodeToBytes(&preCommit)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		expectedMsg := &message.Message{
			Code:          consensus.MsgPrecommit,
			Payload:       encodedVote,
			ConsensusMsg:  &preCommit,
			Address:       addr,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
		}

		c := &Core{
			address:          addr,
			round:            1,
			height:           big.NewInt(2),
			curRoundMessages: curRoundMessages,
			messages:         messages,
			logger:           log.New("backend", "test", "id", 0),
			//committeeSet:     new(validatorSet),
		}

		c.SetDefaultHandlers()
		err = c.precommiter.HandlePrecommit(context.Background(), expectedMsg)
		if err != constants.ErrFutureHeightMessage {
			t.Fatalf("Expected %v, got %v", constants.ErrFutureHeightMessage, err)
		}
	})

	t.Run("pre-commit with invalid signature given, error returned", func(t *testing.T) {
		committeeSet := helpers.NewTestCommitteeSet(4)
		member, _ := committeeSet.GetByIndex(1)
		messages := message.NewMessagesMap()
		curRoundMessages := messages.GetOrCreate(2)
		var preCommit = message.Vote{
			Round:             2,
			Height:            big.NewInt(3),
			ProposedBlockHash: curRoundMessages.GetProposalHash(),
		}

		encodedVote, err := rlp.EncodeToBytes(&preCommit)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		expectedMsg := &message.Message{
			Code:          consensus.MsgPrecommit,
			Payload:       encodedVote,
			ConsensusMsg:  &preCommit,
			Address:       member.Address,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
		}

		c := &Core{
			address:          member.Address,
			round:            2,
			height:           big.NewInt(3),
			curRoundMessages: curRoundMessages,
			logger:           log.New("backend", "test", "id", 0),
		}
		c.SetDefaultHandlers()
		c.SetStep(tctypes.Precommit)
		err = c.precommiter.HandlePrecommit(context.Background(), expectedMsg)
		if err != secp256k1.ErrInvalidSignatureLen {
			t.Fatalf("Expected %v, got %v", secp256k1.ErrInvalidSignatureLen, err)
		}
	})

	t.Run("pre-commit given with no errors, commit called", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		committeeSet, keys := helpers.NewTestCommitteeSetWithKeys(1)
		member, _ := committeeSet.GetByIndex(0)
		logger := log.New("backend", "test", "id", 0)

		proposal := message.NewProposal(
			2,
			big.NewInt(3),
			1,
			types.NewBlockWithHeader(&types.Header{}),
			signer(keys[member.Address]))

		messages := message.NewMessagesMap()
		curRoundMessages := messages.GetOrCreate(2)
		curRoundMessages.SetProposal(proposal, nil, true)

		msg, err := preparePrecommitMsg(proposal.ProposalBlock.Hash(), 2, 3, keys, member)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Commit(proposal.ProposalBlock, gomock.Any(), gomock.Any()).Return(nil).Do(
			func(proposalBlock *types.Block, round int64, seals [][]byte) {
				if round != 2 {
					t.Fatal("Commit called with round different than precommit seal")
				}
				if !reflect.DeepEqual([][]byte{msg.CommittedSeal}, seals) {
					t.Fatal("Commit called with wrong seal")
				}
			})

		c := &Core{
			address:          member.Address,
			backend:          backendMock,
			messages:         messages,
			curRoundMessages: curRoundMessages,
			logger:           logger,
			round:            2,
			height:           big.NewInt(3),
			step:             tctypes.Precommit,
			committee:        committeeSet,
			precommitTimeout: tctypes.NewTimeout(tctypes.Precommit, logger),
		}

		c.SetDefaultHandlers()
		err = c.precommiter.HandlePrecommit(context.Background(), msg)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}
	})

	t.Run("pre-commit given with no errors, commit cancelled", func(t *testing.T) {
		logger := log.New("backend", "test", "id", 0)
		committeeSet, keys := helpers.NewTestCommitteeSetWithKeys(1)
		member, _ := committeeSet.GetByIndex(0)
		proposal := message.NewProposal(
			1,
			big.NewInt(2),
			1,
			types.NewBlockWithHeader(&types.Header{}),
			signer(keys[member.Address]))

		messages := message.NewMessagesMap()
		curRoundMessages := messages.GetOrCreate(2)
		curRoundMessages.SetProposal(proposal, nil, true)

		msg, err := preparePrecommitMsg(proposal.ProposalBlock.Hash(), 1, 2, keys, member)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		c := &Core{
			address:          member.Address,
			curRoundMessages: curRoundMessages,
			messages:         messages,
			logger:           logger,
			round:            1,
			height:           big.NewInt(2),
			committee:        committeeSet,
			step:             tctypes.Precommit,
			precommitTimeout: tctypes.NewTimeout(tctypes.Precommit, logger),
		}
		c.SetDefaultHandlers()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err = c.precommiter.HandlePrecommit(ctx, msg)
		if err != context.Canceled {
			t.Fatalf("Expected %v, got %v", context.Canceled, err)
		}
	})

	t.Run("pre-commit given with no errors, pre-commit Timeout triggered", func(t *testing.T) {
		logger := log.New("backend", "test", "id", 0)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		backendMock := interfaces.NewMockBackend(ctrl)
		committeeSet, keys := helpers.NewTestCommitteeSetWithKeys(7)
		me, _ := committeeSet.GetByIndex(0)
		proposal := message.NewProposal(
			2,
			big.NewInt(3),
			1,
			types.NewBlockWithHeader(&types.Header{}),
			signer(keys[me.Address]))

		messages := message.NewMessagesMap()
		curRoundMessages := messages.GetOrCreate(2)
		curRoundMessages.SetProposal(proposal, nil, true)

		c := &Core{
			address:          me.Address,
			backend:          backendMock,
			curRoundMessages: curRoundMessages,
			messages:         messages,
			logger:           logger,
			round:            2,
			height:           big.NewInt(3),
			step:             tctypes.Precommit,
			committee:        committeeSet,
			precommitTimeout: tctypes.NewTimeout(tctypes.Precommit, logger),
		}
		c.SetDefaultHandlers()
		backendMock.EXPECT().Post(gomock.Any()).Times(1)

		for _, member := range committeeSet.Committee()[1:5] {

			msg, err := preparePrecommitMsg(proposal.ProposalBlock.Hash(), 2, 3, keys, member)
			if err != nil {
				t.Fatalf("Expected nil, got %v", err)
			}
			err = c.precommiter.HandlePrecommit(context.Background(), msg)
			if err != nil {
				t.Fatalf("Expected nil, got %v", err)
			}
		}
		msg, err := preparePrecommitMsg(common.Hash{}, 2, 3, keys, me)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}
		if err = c.precommiter.HandlePrecommit(context.Background(), msg); err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		<-time.NewTimer(5 * time.Second).C

	})
}

func TestVerifyPrecommitCommittedSeal(t *testing.T) {
	t.Run("invalid signer address given, error returned", func(t *testing.T) {
		c := &Core{
			logger: log.New("backend", "test", "id", 0),
		}
		c.SetDefaultHandlers()
		addrMsg := common.HexToAddress("0x0123456789")
		err := c.precommiter.VerifyCommittedSeal(addrMsg, nil, common.Hash{}, 1, big.NewInt(1))
		if err != secp256k1.ErrInvalidSignatureLen {
			t.Fatalf("Expected %v, got %v", secp256k1.ErrInvalidSignatureLen, err)
		}
	})

	t.Run("invalid signer address given, error returned", func(t *testing.T) {
		c := &Core{
			logger: log.New("backend", "test", "id", 0),
		}
		c.SetDefaultHandlers()

		addrMsg := common.HexToAddress("0x0123456789")
		key, err := helpers.GeneratePrivateKey()
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		data := helpers.PrepareCommittedSeal(common.Hash{}, 3, big.NewInt(28))
		hashData := crypto.Keccak256(data)
		sig, err := crypto.Sign(hashData, key)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		err = c.precommiter.VerifyCommittedSeal(addrMsg, sig, common.Hash{}, 3, big.NewInt(28))
		if err != constants.ErrInvalidSenderOfCommittedSeal {
			t.Fatalf("Expected %v, got %v", constants.ErrInvalidSenderOfCommittedSeal, err)
		}
	})

	t.Run("valid params given, no error returned", func(t *testing.T) {
		c := &Core{
			logger: log.New("backend", "test", "id", 0),
		}
		c.SetDefaultHandlers()

		addrMsg := helpers.GetAddress()
		key, err := helpers.GeneratePrivateKey()
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		data := helpers.PrepareCommittedSeal(addrMsg.Hash(), 1, big.NewInt(13))
		hashData := crypto.Keccak256(data)
		sig, err := crypto.Sign(hashData, key)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		err = c.precommiter.VerifyCommittedSeal(addrMsg, sig, addrMsg.Hash(), 1, big.NewInt(13))
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}
	})
}

func TestHandleCommit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := log.New("backend", "test", "id", 0)

	addr := common.HexToAddress("0x0123456789")
	testCommittee, keys := helpers.GenerateCommittee(3)

	firstKey := keys[testCommittee[0].Address]

	h := &types.Header{Number: big.NewInt(3)}

	// Sign the header so that types.Ecrecover works
	seal, err := crypto.Sign(crypto.Keccak256(types.SigHash(h).Bytes()), firstKey)
	require.NoError(t, err)

	err = types.WriteSeal(h, seal)
	require.NoError(t, err)

	h.Committee = testCommittee

	block := types.NewBlockWithHeader(h)
	testCommittee = append(testCommittee, types.CommitteeMember{Address: addr, VotingPower: big.NewInt(1)})
	committeeSet, err := committee.NewRoundRobinSet(testCommittee, testCommittee[0].Address)
	require.NoError(t, err)

	backendMock := interfaces.NewMockBackend(ctrl)
	backendMock.EXPECT().HeadBlock().MinTimes(1).Return(block, addr)

	c := &Core{
		address:          addr,
		backend:          backendMock,
		round:            2,
		height:           big.NewInt(3),
		messages:         message.NewMessagesMap(),
		logger:           logger,
		proposeTimeout:   tctypes.NewTimeout(tctypes.Propose, logger),
		prevoteTimeout:   tctypes.NewTimeout(tctypes.Prevote, logger),
		precommitTimeout: tctypes.NewTimeout(tctypes.Precommit, logger),
		committee:        committeeSet,
	}
	c.SetDefaultHandlers()
	c.precommiter.HandleCommit(context.Background())
	if c.round != 0 || c.height.Cmp(big.NewInt(4)) != 0 {
		t.Fatalf("Expected new round")
	}
	// to fix the data race detected by CI workflow.
	err = c.proposeTimeout.StopTimer()
	if err != nil {
		t.Error(err)
	}
}

func preparePrecommitMsg(proposalHash common.Hash, round int64, height int64, keys helpers.AddressKeyMap, member types.CommitteeMember) (*message.Message, error) {
	var preCommit = message.Vote{
		Round:             round,
		Height:            big.NewInt(height),
		ProposedBlockHash: proposalHash,
	}

	encodedVote, err := rlp.EncodeToBytes(&preCommit)
	if err != nil {
		return nil, err
	}

	data := helpers.PrepareCommittedSeal(proposalHash, preCommit.Round, preCommit.Height)
	hashData := crypto.Keccak256(data)
	sig, err := crypto.Sign(hashData, keys[member.Address])
	if err != nil {
		return nil, err
	}

	expectedMsg := &message.Message{
		Code:          consensus.MsgPrecommit,
		Payload:       encodedVote,
		ConsensusMsg:  &preCommit,
		Address:       member.Address,
		CommittedSeal: sig,
		Signature:     []byte{0x1},
		Power:         member.VotingPower,
	}
	return expectedMsg, err
}
