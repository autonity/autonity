package core

import (
	"context"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/committee"
	"github.com/clearmatics/autonity/core/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math/big"
	"math/rand"
	"strconv"
	"testing"
)

func prepareCommittee() types.Committee {
	// prepare committee.
	minSize := 4
	maxSize := 15
	committeeSize := rand.Intn(maxSize-minSize) + minSize
	committeeSet := types.Committee{}
	for i := 1; i <= committeeSize; i++ {
		hexString := "0x01234567890" + strconv.Itoa(i)
		member := types.CommitteeMember{
			Address:     common.HexToAddress(hexString),
			VotingPower: new(big.Int).SetInt64(1),
		}
		committeeSet = append(committeeSet, member)
	}
	return committeeSet
}

func generateBlock(height *big.Int) *types.Block {
	header := &types.Header{Number: height}
	block := types.NewBlock(header, nil, nil, nil)
	return block
}

func TestTendermintOn(t *testing.T) {
	t.Run("Line44 to Line46, step from prevote to precommit by voting for nil.", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// prepare a random size of committee, and the proposer at last committed block.
		currentCommittee := prepareCommittee()
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
		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().Return(nil, nil)

		validRoundProposed := int64(0)
		roundProposed := int64(0)

		// create consensus core and conditions.
		c := New(backendMock)
		c.committeeSet = committeeSet
		c.height = currentHeight
		c.round = roundProposed
		c.lockedRound = -1
		c.lockedValue = nil
		c.validRound = -1
		c.validValue = nil
		c.setValidRoundAndValue = false
		c.step = prevote

		// condition 2f+1 <PREVOTE, h_p, round_p, id(v)>, power of pre-vote. line 44.
		receivedPrevoteMsg := Message{
			Code:    msgPrevote,
			Address: currentCommittee[2].Address,
			power:   c.CommitteeSet().Quorum() - 1,
		}
		proposal := NewProposal(roundProposed, currentHeight, validRoundProposed, proposalBlock)
		encodedProposal, err := Encode(proposal)
		if err != nil {
			t.Error(err)
		}
		proposalMsg := &Message{
			Code:    msgProposal,
			Msg:     encodedProposal,
			Address: currentCommittee[1].Address,
		}
		c.curRoundMessages.SetProposal(proposal, proposalMsg, true)

		// prepare the wanted msg which will be broadcast.
		vote := Vote{
			Round:             roundProposed,
			Height:            currentHeight,
			ProposedBlockHash: common.Hash{},
		}
		preCommitMsg, err := Encode(&vote)
		if err != nil {
			t.Error("err")
		}
		wantedMsg, err := c.finalizeMessage(&Message{
			Code:          msgPrecommit,
			Msg:           preCommitMsg,
			Address:       clientAddr,
			CommittedSeal: []byte{},
		})
		if err != nil {
			t.Error(err)
		}
		// should check if broadcast to wanted committeeSet with wanted prevote msg.
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, wantedMsg).Return(nil)


		// received 2f Prevote for nil already.
		c.curRoundMessages.AddPrevote(common.Hash{}, receivedPrevoteMsg)
		c.messages.getOrCreate(validRoundProposed).AddPrevote(proposalBlock.Hash(), receivedPrevoteMsg)
		// on prevote of nil
		preVoteMsg := createPrevote(t, common.Hash{}, roundProposed, currentHeight, committeeSet.Committee()[0])
		err = c.handlePrevote(context.Background(), preVoteMsg)
		if err != nil {
			t.Error(err)
		}

		// checking internal state of tendermint.
		assert.True(t, c.sentPrecommit)
		assert.Equal(t, c.step, precommit)
		assert.Nil(t, c.lockedValue)
		assert.Nil(t, c.validValue)
		assert.Equal(t, c.lockedRound, int64(-1))
		assert.Equal(t, c.validRound, int64(-1))
	})
}