package core

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/committee"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math/big"
	"math/rand"
	"sort"
	"testing"
	"time"
)

// The following tests are not specific to proposal messages but rather apply to all messages
func TestHandleMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	currentCommittee, keys := prepareCommittee(t)
	committeeSet, err := committee.NewSet(currentCommittee, currentCommittee[rand.Intn(len(currentCommittee))].Address)
	assert.Nil(t, err)

	backendMock := NewMockBackend(ctrl)
	// We don't care what address the client has
	backendMock.EXPECT().Address().Return(currentCommittee[0].Address).MaxTimes(2)
	core := New(backendMock)

	t.Run("message sender is not in the committee set", func(t *testing.T) {
		// Prepare message
		key, err := crypto.GenerateKey()
		assert.Nil(t, err)

		msg := &Message{Address: crypto.PubkeyToAddress(key.PublicKey), Code: uint64(rand.Intn(3)), Msg: []byte("random message1")}

		msgRlpNoSig, err := msg.PayloadNoSig()
		assert.Nil(t, err)

		msg.Signature, err = crypto.Sign(crypto.Keccak256(msgRlpNoSig), key)
		assert.Nil(t, err)

		msgRlpWithSig, err := msg.Payload()
		assert.Nil(t, err)

		core.setCommitteeSet(committeeSet)
		err = core.handleMsg(context.Background(), msgRlpWithSig)

		assert.Error(t, err, "unauthorised sender, sender is not is committees set")
	})

	t.Run("message sender is not the message siger", func(t *testing.T) {
		msg := &Message{Address: crypto.PubkeyToAddress(keys[0].PublicKey), Code: uint64(rand.Intn(3)), Msg: []byte("random message2")}

		msgRlpNoSig, err := msg.PayloadNoSig()
		assert.Nil(t, err)

		msg.Signature, err = crypto.Sign(crypto.Keccak256(msgRlpNoSig), keys[1])
		assert.Nil(t, err)

		msgRlpWithSig, err := msg.Payload()
		assert.Nil(t, err)

		core.setCommitteeSet(committeeSet)
		err = core.handleMsg(context.Background(), msgRlpWithSig)

		assert.Error(t, err, "unauthorised sender, sender is not the signer of the message")
	})

	t.Run("malicious sender sends incorrect signature", func(t *testing.T) {
		sig, err := crypto.Sign(crypto.Keccak256([]byte("random bytes")), keys[0])
		assert.Nil(t, err)

		msg := &Message{Address: crypto.PubkeyToAddress(keys[0].PublicKey), Code: uint64(rand.Intn(3)), Msg: []byte("random message2"), Signature: sig}
		msgRlpWithSig, err := msg.Payload()

		core.setCommitteeSet(committeeSet)
		err = core.handleMsg(context.Background(), msgRlpWithSig)

		assert.Error(t, err, "malicious sender sends different signature to signature of message")
	})
}

// It test the page-6, line 22 to line 27, on new proposal logic of tendermint pseudo-code.
// Please refer to the algorithm from here: https://arxiv.org/pdf/1807.04938.pdf
func TestTendermintNewProposal(t *testing.T) {
	// Below 4 test cases cover line 22 to line 27 of tendermint pseudo-code.
	// It test line 24 was executed and step was forwarded on line 27.

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// prepare a random size of committee, and the proposer at last committed block.
	currentCommittee, _ := prepareCommittee(t)
	sort.Sort(currentCommittee)
	lastProposer := currentCommittee[len(currentCommittee)-1].Address
	committeeSet, err := committee.NewSet(currentCommittee, lastProposer)
	if err != nil {
		t.Error(err)
	}

	currentHeight := big.NewInt(10)
	proposalBlock := generateBlock(currentHeight)
	clientAddr := currentCommittee[0].Address

	backendMock := NewMockBackend(ctrl)
	backendMock.EXPECT().Address().AnyTimes().Return(clientAddr)

	t.Run("on proposal with validRound as (-1) proposed and local lockedRound as (-1)", func(t *testing.T) {
		// create consensus core.
		c := New(backendMock)
		c.committeeSet = committeeSet
		c.height = currentHeight
		c.lockedRound = -1
		c.step = propose

		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().Return(nil, nil)
		backendMock.EXPECT().VerifyProposal(gomock.Any()).AnyTimes().Return(time.Duration(1), nil)

		// prepare input msg
		validRoundProposed := int64(-1)
		proposal := NewProposal(0, currentHeight, validRoundProposed, proposalBlock)
		encodedProposal, err := Encode(proposal)
		if err != nil {
			t.Error(err)
		}
		msg := &Message{
			Code:    msgProposal,
			Msg:     encodedProposal,
			Address: clientAddr,
		}

		// prepare the wanted msg which will be broadcast.
		vote := Vote{
			Round:             0,
			Height:            currentHeight,
			ProposedBlockHash: proposalBlock.Hash(),
		}
		prevoteMsg, err := Encode(&vote)
		if err != nil {
			t.Error("err")
		}
		wantedMsg, err := c.finalizeMessage(&Message{
			Code:          msgPrevote,
			Msg:           prevoteMsg,
			Address:       clientAddr,
			CommittedSeal: []byte{},
		})
		if err != nil {
			t.Error(err)
		}

		// should check if broadcast to wanted committeeSet with wanted prevote msg.
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, wantedMsg).Return(nil)

		err = c.handleProposal(context.Background(), msg)
		if err != nil {
			t.Error(err)
		}
		// checking internal state of tendermint.
		assert.True(t, c.sentPrevote)
		assert.Equal(t, c.step, prevote)
		assert.Nil(t, c.lockedValue)
		assert.Equal(t, c.lockedRound, int64(-1))
		assert.Nil(t, c.validValue)
		assert.Equal(t, c.validRound, int64(-1))
	})

	// It test line 24 was executed and step was forwarded on line 27 but with below different condition.
	t.Run("on proposal with validRound as (-1) proposed and local lockedRound as (1), but locked at the same value as proposed already.", func(t *testing.T) {
		// create consensus core.
		c := New(backendMock)
		c.committeeSet = committeeSet
		c.height = currentHeight
		c.lockedRound = 1 // set lockedRound as 1.
		c.lockedValue = proposalBlock
		c.validRound = 1
		c.validValue = proposalBlock
		c.step = propose

		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().Return(nil, nil)
		backendMock.EXPECT().VerifyProposal(gomock.Any()).AnyTimes().Return(time.Duration(1), nil)

		validRoundProposed := int64(-1)
		proposal := NewProposal(0, currentHeight, validRoundProposed, proposalBlock)
		encodedProposal, err := Encode(proposal)
		if err != nil {
			t.Error(err)
		}

		// prepare input msg
		msg := &Message{
			Code:    msgProposal,
			Msg:     encodedProposal,
			Address: clientAddr,
		}

		// prepare the wanted msg which will be broadcast.
		vote := Vote{
			Round:             0,
			Height:            currentHeight,
			ProposedBlockHash: proposalBlock.Hash(),
			//ProposedBlockHash: common.Hash{},
		}
		prevoteMsg, err := Encode(&vote)
		if err != nil {
			t.Error("err")
		}
		wantedMsg, err := c.finalizeMessage(&Message{
			Code:          msgPrevote,
			Msg:           prevoteMsg,
			Address:       clientAddr,
			CommittedSeal: []byte{},
		})
		if err != nil {
			t.Error(err)
		}

		// should check if broadcast to wanted committeeSet with wanted prevote msg.
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, wantedMsg).Return(nil)
		err = c.handleProposal(context.Background(), msg)
		if err != nil {
			t.Error(err)
		}
		// checking internal state of tendermint
		assert.True(t, c.sentPrevote)
		assert.Equal(t, c.step, prevote)
		assert.Equal(t, c.lockedValue, proposalBlock)
		assert.Equal(t, c.lockedRound, int64(1))
		assert.Equal(t, c.validValue, proposalBlock)
		assert.Equal(t, c.validRound, int64(1))
	})

	// It test line 26 was executed and step was forwarded on line 27 but with below different condition.
	t.Run("on proposal with validRound as (-1) proposed and local lockedRound as (1) and locked at different value, vote for nil", func(t *testing.T) {
		// create consensus core.
		lockedValue := generateBlock(big.NewInt(11))
		c := New(backendMock)
		c.committeeSet = committeeSet
		c.height = currentHeight
		c.lockedRound = 1
		c.lockedValue = lockedValue
		c.validRound = 1
		c.validValue = lockedValue

		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().Return(nil, nil)
		backendMock.EXPECT().VerifyProposal(gomock.Any()).AnyTimes().Return(time.Duration(1), nil)

		// prepare input proposal msg.
		validRoundProposed := int64(-1)
		proposal := NewProposal(0, currentHeight, validRoundProposed, proposalBlock)
		encodedProposal, err := Encode(proposal)
		if err != nil {
			t.Error(err)
		}

		msg := &Message{
			Code:    msgProposal,
			Msg:     encodedProposal,
			Address: clientAddr,
		}

		// prepare the wanted vote for nil msg which will be broadcast.
		vote := Vote{
			Round:             0,
			Height:            currentHeight,
			ProposedBlockHash: common.Hash{},
		}
		prevoteMsg, err := Encode(&vote)
		if err != nil {
			t.Error("err")
		}
		wantedMsg, err := c.finalizeMessage(&Message{
			Code:          msgPrevote,
			Msg:           prevoteMsg,
			Address:       clientAddr,
			CommittedSeal: []byte{},
		})
		if err != nil {
			t.Error(err)
		}

		// should check if broadcast to wanted committeeSet with wanted prevote msg.
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, wantedMsg).Return(nil)
		err = c.handleProposal(context.Background(), msg)
		if err != nil {
			t.Error(err)
		}

		// checking internal state of tendermint.
		assert.True(t, c.sentPrevote)
		assert.Equal(t, c.step, prevote)
		assert.Equal(t, c.lockedValue, lockedValue)
		assert.Equal(t, c.lockedRound, int64(1))
		assert.Equal(t, c.validValue, lockedValue)
		assert.Equal(t, c.validRound, int64(1))
	})

	// It test line 26 was executed and step was forwarded on line 27 but with invalid value proposed.
	t.Run("on proposal with invalid block, follower should step forward with voting for nil", func(t *testing.T) {
		// create consensus core.
		c := New(backendMock)
		c.committeeSet = committeeSet
		c.height = currentHeight
		c.lockedRound = -1
		c.lockedValue = nil
		c.validRound = -1
		c.validValue = nil

		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().Return(nil, nil)
		backendMock.EXPECT().VerifyProposal(gomock.Any()).AnyTimes().Return(time.Duration(1), errors.New("invalid block"))

		// prepare input proposal msg.
		validRoundProposed := int64(-1)
		proposal := NewProposal(0, currentHeight, validRoundProposed, proposalBlock)
		encodedProposal, err := Encode(proposal)
		if err != nil {
			t.Error(err)
		}

		msg := &Message{
			Code:    msgProposal,
			Msg:     encodedProposal,
			Address: clientAddr,
		}

		// prepare the wanted vote for nil msg which will be broadcast.
		vote := Vote{
			Round:             0,
			Height:            currentHeight,
			ProposedBlockHash: common.Hash{},
		}
		prevoteMsg, err := Encode(&vote)
		if err != nil {
			t.Error("err")
		}
		wantedMsg, err := c.finalizeMessage(&Message{
			Code:          msgPrevote,
			Msg:           prevoteMsg,
			Address:       clientAddr,
			CommittedSeal: []byte{},
		})
		if err != nil {
			t.Error(err)
		}

		// should check if broadcast to wanted committeeSet with wanted prevote msg.
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, wantedMsg).Return(nil)
		err = c.handleProposal(context.Background(), msg)
		if err != nil {
			assert.Equal(t, err.Error(), "invalid block")
		}

		// checking inernal state of tendermint.
		assert.True(t, c.sentPrevote)
		assert.Equal(t, c.step, prevote)
		assert.Nil(t, c.lockedValue)
		assert.Equal(t, c.lockedRound, int64(-1))
		assert.Nil(t, c.validValue)
		assert.Equal(t, c.validRound, int64(-1))
	})
}

func prepareCommittee(t *testing.T) (types.Committee, []*ecdsa.PrivateKey) {
	t.Helper()

	minSize, maxSize := 4, 15
	committeeSize := rand.Intn(maxSize-minSize) + minSize

	var committeeSet types.Committee
	var keys []*ecdsa.PrivateKey

	for i := 1; i <= committeeSize; i++ {
		key, err := crypto.GenerateKey()
		if err != nil {
			t.Fatal(err)
		}
		member := types.CommitteeMember{
			Address:     crypto.PubkeyToAddress(key.PublicKey),
			VotingPower: new(big.Int).SetInt64(1),
		}
		committeeSet = append(committeeSet, member)
		keys = append(keys, key)
	}
	return committeeSet, keys
}

func generateBlock(height *big.Int) *types.Block {
	header := &types.Header{Number: height}
	block := types.NewBlock(header, nil, nil, nil)
	return block
}
