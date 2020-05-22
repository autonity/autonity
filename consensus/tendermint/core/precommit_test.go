package core

import (
	"context"
	"errors"
	"github.com/clearmatics/autonity/consensus/tendermint/committee"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/crypto/secp256k1"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/log"
	"github.com/golang/mock/gomock"
)

func TestSendPrecommit(t *testing.T) {
	t.Run("proposal is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Broadcast(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		messages := newMessagesMap()
		c := &core{
			logger:           log.New("backend", "test", "id", 0),
			backend:          backendMock,
			messages:         messages,
			curRoundMessages: messages.getOrCreate(0),
			round:            2,
			height:           big.NewInt(3),
		}

		c.sendPrecommit(context.Background(), false)
	})

	t.Run("valid proposal given, non nil pre-commit", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		committeeSet := newTestCommitteeSet(4)
		logger := log.New("backend", "test", "id", 0)
		val, _ := committeeSet.GetByIndex(2)
		addr := val.Address

		proposal := NewProposal(
			1,
			big.NewInt(2),
			1,
			types.NewBlockWithHeader(&types.Header{}))

		messages := newMessagesMap()
		curRoundMessages := messages.getOrCreate(1)
		curRoundMessages.SetProposal(proposal, nil, false)

		var preCommit = Vote{
			Round:             1,
			Height:            big.NewInt(2),
			ProposedBlockHash: curRoundMessages.GetProposalHash(),
		}

		encodedVote, err := Encode(&preCommit)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		expectedMsg := &Message{
			Code:          msgPrecommit,
			Msg:           encodedVote,
			Address:       addr,
			CommittedSeal: []byte{0x1},
			Signature:     []byte{0x1},
		}

		payloadNoSig, err := expectedMsg.PayloadNoSig()
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).Return([]byte{0x1}, nil)
		backendMock.EXPECT().Sign(gomock.Eq(payloadNoSig)).Return([]byte{0x1}, nil)

		payload, err := expectedMsg.Payload()
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		backendMock.EXPECT().Broadcast(gomock.Any(), gomock.Any(), payload)

		c := &core{
			backend:          backendMock,
			address:          addr,
			logger:           logger,
			committeeSet:     committeeSet,
			messages:         messages,
			curRoundMessages: curRoundMessages,
			round:            1,
			height:           big.NewInt(2),
		}

		c.sendPrecommit(context.Background(), false)
	})

	t.Run("valid proposal given, nil pre-commit", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		committeeSet := newTestCommitteeSet(4)
		logger := log.New("backend", "test", "id", 0)
		val, _ := committeeSet.GetByIndex(2)
		addr := val.Address

		proposal := NewProposal(
			1,
			big.NewInt(2),
			1,
			types.NewBlockWithHeader(&types.Header{}))

		messages := newMessagesMap()
		curRoundMessages := messages.getOrCreate(1)
		curRoundMessages.SetProposal(proposal, nil, true)

		var preCommit = Vote{
			Round:             1,
			Height:            big.NewInt(2),
			ProposedBlockHash: common.Hash{},
		}

		encodedVote, err := Encode(&preCommit)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		expectedMsg := &Message{
			Code:          msgPrecommit,
			Msg:           encodedVote,
			Address:       addr,
			CommittedSeal: []byte{0x1},
			Signature:     []byte{0x1},
		}

		payloadNoSig, err := expectedMsg.PayloadNoSig()
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).Return([]byte{0x1}, errors.New("seal sign error"))
		backendMock.EXPECT().Sign(payloadNoSig).Return([]byte{0x1}, nil)

		payload, err := expectedMsg.Payload()
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		backendMock.EXPECT().Broadcast(gomock.Any(), gomock.Any(), payload)

		c := &core{
			backend:          backendMock,
			address:          addr,
			logger:           logger,
			curRoundMessages: curRoundMessages,
			messages:         messages,
			committeeSet:     committeeSet,
			height:           big.NewInt(2),
			round:            1,
		}

		c.sendPrecommit(context.Background(), true)
	})

}

func TestHandlePrecommit(t *testing.T) {
	t.Run("pre-commit with future height given, error returned", func(t *testing.T) {

		addr := common.HexToAddress("0x0123456789")
		messages := newMessagesMap()
		curRoundMessages := messages.getOrCreate(1)
		var preCommit = Vote{
			Round:  2,
			Height: big.NewInt(3),
		}

		encodedVote, err := Encode(&preCommit)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		expectedMsg := &Message{
			Code:          msgPrecommit,
			Msg:           encodedVote,
			Address:       addr,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
		}

		c := &core{
			address:          addr,
			round:            1,
			height:           big.NewInt(2),
			curRoundMessages: curRoundMessages,
			messages:         messages,
			logger:           log.New("backend", "test", "id", 0),
			//committeeSet:     new(validatorSet),
		}

		err = c.handlePrecommit(context.Background(), expectedMsg)
		if err != errFutureHeightMessage {
			t.Fatalf("Expected %v, got %v", errFutureHeightMessage, err)
		}
	})

	t.Run("pre-commit with invalid signature given, error returned", func(t *testing.T) {
		committeeSet := newTestCommitteeSet(4)
		member, _ := committeeSet.GetByIndex(1)
		messages := newMessagesMap()
		curRoundMessages := messages.getOrCreate(2)
		var preCommit = Vote{
			Round:             2,
			Height:            big.NewInt(3),
			ProposedBlockHash: curRoundMessages.GetProposalHash(),
		}

		encodedVote, err := Encode(&preCommit)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		expectedMsg := &Message{
			Code:          msgPrecommit,
			Msg:           encodedVote,
			Address:       member.Address,
			CommittedSeal: []byte{},
			Signature:     []byte{0x1},
		}

		c := &core{
			address:          member.Address,
			round:            2,
			height:           big.NewInt(3),
			curRoundMessages: curRoundMessages,
			logger:           log.New("backend", "test", "id", 0),
		}
		c.setStep(precommit)
		err = c.handlePrecommit(context.Background(), expectedMsg)
		if err != secp256k1.ErrInvalidSignatureLen {
			t.Fatalf("Expected %v, got %v", secp256k1.ErrInvalidSignatureLen, err)
		}
	})

	t.Run("pre-commit given with no errors, commit called", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		committeeSet, keys := newTestCommitteeSetWithKeys(1)
		member, _ := committeeSet.GetByIndex(0)
		logger := log.New("backend", "test", "id", 0)

		proposal := NewProposal(
			2,
			big.NewInt(3),
			1,
			types.NewBlockWithHeader(&types.Header{}))

		messages := newMessagesMap()
		curRoundMessages := messages.getOrCreate(2)
		curRoundMessages.SetProposal(proposal, nil, true)

		msg, err := preparePrecommitMsg(proposal.ProposalBlock.Hash(), 2, 3, keys, member)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Commit(proposal.ProposalBlock, gomock.Any(), gomock.Any()).Return(nil).Do(
			func(proposalBlock *types.Block, round int64, seals [][]byte) {
				if round != 2 {
					t.Fatal("Commit called with round different than precommit seal")
				}
				if !reflect.DeepEqual([][]byte{msg.CommittedSeal}, seals) {
					t.Fatal("Commit called with wrong seal")
				}
			})

		c := &core{
			address:          member.Address,
			backend:          backendMock,
			messages:         messages,
			curRoundMessages: curRoundMessages,
			logger:           logger,
			round:            2,
			height:           big.NewInt(3),
			step:             precommit,
			committeeSet:     committeeSet,
			precommitTimeout: newTimeout(precommit, logger),
		}

		err = c.handlePrecommit(context.Background(), msg)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}
	})

	t.Run("pre-commit given with no errors, commit cancelled", func(t *testing.T) {
		logger := log.New("backend", "test", "id", 0)
		committeeSet, keys := newTestCommitteeSetWithKeys(1)
		member, _ := committeeSet.GetByIndex(0)
		proposal := NewProposal(
			1,
			big.NewInt(2),
			1,
			types.NewBlockWithHeader(&types.Header{}))

		messages := newMessagesMap()
		curRoundMessages := messages.getOrCreate(2)
		curRoundMessages.SetProposal(proposal, nil, true)

		msg, err := preparePrecommitMsg(proposal.ProposalBlock.Hash(), 1, 2, keys, member)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		c := &core{
			address:          member.Address,
			curRoundMessages: curRoundMessages,
			messages:         messages,
			logger:           logger,
			round:            1,
			height:           big.NewInt(2),
			committeeSet:     committeeSet,
			step:             precommit,
			precommitTimeout: newTimeout(precommit, logger),
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err = c.handlePrecommit(ctx, msg)
		if err != context.Canceled {
			t.Fatalf("Expected %v, got %v", context.Canceled, err)
		}
	})

	t.Run("pre-commit given with no errors, pre-commit timeout triggered", func(t *testing.T) {
		logger := log.New("backend", "test", "id", 0)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		backendMock := NewMockBackend(ctrl)
		proposal := NewProposal(
			2,
			big.NewInt(3),
			1,
			types.NewBlockWithHeader(&types.Header{}))

		messages := newMessagesMap()
		curRoundMessages := messages.getOrCreate(2)
		curRoundMessages.SetProposal(proposal, nil, true)

		committeeSet, keys := newTestCommitteeSetWithKeys(7)
		me, _ := committeeSet.GetByIndex(0)

		c := &core{
			address:          me.Address,
			backend:          backendMock,
			curRoundMessages: curRoundMessages,
			messages:         messages,
			logger:           logger,
			round:            2,
			height:           big.NewInt(3),
			step:             precommit,
			committeeSet:     committeeSet,
			precommitTimeout: newTimeout(precommit, logger),
		}
		backendMock.EXPECT().Post(gomock.Any()).Times(1)

		for _, member := range committeeSet.Committee()[1:5] {

			msg, err := preparePrecommitMsg(proposal.ProposalBlock.Hash(), 2, 3, keys, member)
			if err != nil {
				t.Fatalf("Expected nil, got %v", err)
			}
			err = c.handlePrecommit(context.Background(), msg)
			if err != nil {
				t.Fatalf("Expected nil, got %v", err)
			}
		}
		msg, err := preparePrecommitMsg(common.Hash{}, 2, 3, keys, me)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}
		if err = c.handlePrecommit(context.Background(), msg); err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		<-time.NewTimer(5 * time.Second).C

	})
}

func TestVerifyPrecommitCommittedSeal(t *testing.T) {
	t.Run("invalid signer address given, error returned", func(t *testing.T) {
		c := &core{
			logger: log.New("backend", "test", "id", 0),
		}

		addrMsg := common.HexToAddress("0x0123456789")

		err := c.verifyCommittedSeal(addrMsg, nil, common.Hash{}, 1, big.NewInt(1))
		if err != secp256k1.ErrInvalidSignatureLen {
			t.Fatalf("Expected %v, got %v", secp256k1.ErrInvalidSignatureLen, err)
		}
	})

	t.Run("invalid signer address given, error returned", func(t *testing.T) {
		c := &core{
			logger: log.New("backend", "test", "id", 0),
		}

		addrMsg := common.HexToAddress("0x0123456789")
		key, err := generatePrivateKey()
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		data := PrepareCommittedSeal(common.Hash{}, 3, big.NewInt(28))
		hashData := crypto.Keccak256(data)
		sig, err := crypto.Sign(hashData, key)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		err = c.verifyCommittedSeal(addrMsg, sig, common.Hash{}, 3, big.NewInt(28))
		if err != errInvalidSenderOfCommittedSeal {
			t.Fatalf("Expected %v, got %v", errInvalidSenderOfCommittedSeal, err)
		}
	})

	t.Run("valid params given, no error returned", func(t *testing.T) {
		c := &core{
			logger: log.New("backend", "test", "id", 0),
		}

		addrMsg := getAddress()
		key, err := generatePrivateKey()
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		data := PrepareCommittedSeal(addrMsg.Hash(), 1, big.NewInt(13))
		hashData := crypto.Keccak256(data)
		sig, err := crypto.Sign(hashData, key)
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}

		err = c.verifyCommittedSeal(addrMsg, sig, addrMsg.Hash(), 1, big.NewInt(13))
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}
	})
}

func TestHandleCommit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := log.New("backend", "test", "id", 0)

	block := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(3)})
	addr := common.HexToAddress("0x0123456789")

	backendMock := NewMockBackend(ctrl)
	backendMock.EXPECT().LastCommittedProposal().MinTimes(1).Return(block, addr)

	testCommittee, _ := generateCommittee(3)
	testCommittee = append(testCommittee, types.CommitteeMember{Address: addr, VotingPower: big.NewInt(1)})
	committeeSet, err := committee.NewSet(testCommittee, testCommittee[0].Address)
	if err != nil {
		t.Error(err)
	}

	backendMock.EXPECT().Committee(gomock.Any()).Return(committeeSet, nil)

	c := &core{
		address:          addr,
		backend:          backendMock,
		round:            2,
		height:           big.NewInt(3),
		messages:         newMessagesMap(),
		logger:           logger,
		proposeTimeout:   newTimeout(propose, logger),
		prevoteTimeout:   newTimeout(prevote, logger),
		precommitTimeout: newTimeout(precommit, logger),
		committeeSet:     committeeSet,
	}
	c.handleCommit(context.Background())
	if c.round != 0 || c.height.Cmp(big.NewInt(4)) != 0 {
		t.Fatalf("Expected new round")
	}
	// cleanup timer otherwise it would cause panic.
	c.proposeTimeout.stopTimer()
}

func preparePrecommitMsg(proposalHash common.Hash, round int64, height int64, keys addressKeyMap, member types.CommitteeMember) (*Message, error) {
	var preCommit = Vote{
		Round:             round,
		Height:            big.NewInt(height),
		ProposedBlockHash: proposalHash,
	}

	encodedVote, err := Encode(&preCommit)
	if err != nil {
		return nil, err
	}

	data := PrepareCommittedSeal(proposalHash, preCommit.Round, preCommit.Height)
	hashData := crypto.Keccak256(data)
	sig, err := crypto.Sign(hashData, keys[member.Address])
	if err != nil {
		return nil, err
	}

	expectedMsg := &Message{
		Code:          msgPrecommit,
		Msg:           encodedVote,
		Address:       member.Address,
		CommittedSeal: sig,
		Signature:     []byte{0x1},
		power:         member.VotingPower.Uint64(),
	}
	return expectedMsg, err
}
