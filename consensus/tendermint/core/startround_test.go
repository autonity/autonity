package core

import (
	"context"
	"github.com/clearmatics/autonity/crypto"
	"sort"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/committee"
	"github.com/clearmatics/autonity/core/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math/big"
	"math/rand"
	"testing"
)

// The following tests aim to test lines 1 - 21 of Tendermint Algorithm described on page 6 of
// https://arxiv.org/pdf/1807.04938.pdf.

func TestStartRoundVariables(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	backendMock := NewMockBackend(ctrl)

	prevHeight := big.NewInt(int64(rand.Intn(100) + 1))
	prevBlock := generateBlock(prevHeight, types.BlockNonce{1, 2, 3, 4, 5, 6, 7, 8})
	currentHeight := big.NewInt(prevHeight.Int64() + 1)
	currentBlock := generateBlock(currentHeight, types.BlockNonce{11, 22, 33, 44, 55, 66, 77, 88})
	currentRound := int64(0)

	// We don't care who is the next proposer so for simplicity we ensure that clientAddress is not the next
	// proposer by setting clientAddress to be the last proposer. This will ensure that the test will not run the
	// broadcast method.
	committeeSet := prepareCommittee(t)
	members := committeeSet.Committee()
	clientAddress := members[len(members)-1].Address

	overrideDefaultCoreValues := func(core *core) {
		core.height = big.NewInt(-1)
		core.round = int64(-1)
		core.step = precommitDone
		core.lockedValue = &types.Block{}
		core.lockedRound = 0
		core.validValue = &types.Block{}
		core.validRound = 0
	}

	t.Run("ensure round 0 state variables are set correctly", func(t *testing.T) {
		backendMock.EXPECT().Address().Return(clientAddress)
		backendMock.EXPECT().LastCommittedProposal().Return(prevBlock, clientAddress)
		backendMock.EXPECT().Committee(currentHeight.Uint64()).Return(committeeSet, nil)

		core := New(backendMock)
		overrideDefaultCoreValues(core)
		core.startRound(context.Background(), currentRound)

		// Check the initial consensus state
		assert.Equal(t, currentHeight, core.Height())
		assert.Equal(t, currentRound, core.Round())
		assert.Equal(t, propose, core.step)
		assert.Nil(t, core.lockedValue)
		assert.Equal(t, int64(-1), core.lockedRound)
		assert.Nil(t, core.validValue)
		assert.Equal(t, int64(-1), core.validRound)
	})

	t.Run("ensure round x state variables are updated correctly", func(t *testing.T) {
		// In this test we are interested in making sure that that values which change in the current round that may
		// have an impact on the actions performed in the following round (in case of round change) are persisted
		// through to the subsequent round.
		backendMock.EXPECT().Address().Return(clientAddress)
		backendMock.EXPECT().LastCommittedProposal().Return(prevBlock, clientAddress).MaxTimes(2)
		backendMock.EXPECT().Committee(currentHeight.Uint64()).Return(committeeSet, nil).MaxTimes(2)

		core := New(backendMock)
		overrideDefaultCoreValues(core)
		core.startRound(context.Background(), currentRound)

		// Check the initial consensus state
		assert.Equal(t, currentHeight, core.Height())
		assert.Equal(t, currentRound, core.Round())
		assert.Equal(t, propose, core.step)
		assert.Nil(t, core.lockedValue)
		assert.Equal(t, int64(-1), core.lockedRound)
		assert.Nil(t, core.validValue)
		assert.Equal(t, int64(-1), core.validRound)

		// Update locked and valid Value (if locked value changes then valid value also changes, ie quorum(prevotes)
		// delivered in prevote step)
		core.lockedValue = currentBlock
		core.lockedRound = currentRound
		core.validValue = currentBlock
		core.validRound = currentRound

		// Move to next round and check the expected state
		core.startRound(context.Background(), currentRound+1)

		assert.Equal(t, core.Height(), currentHeight)
		assert.Equal(t, currentRound+1, core.Round())
		assert.Equal(t, propose, core.step)
		assert.Equal(t, currentBlock, core.lockedValue)
		assert.Equal(t, currentRound, core.lockedRound)
		assert.Equal(t, currentBlock, core.validValue)
		assert.Equal(t, currentRound, core.validRound)

		// Update valid value (we didn't receive quorum prevote in prevote step, also the block changed, ie, locked
		// value and valid value are different)
		currentBlock2 := generateBlock(currentHeight, types.BlockNonce{12, 23, 34, 45, 56, 67, 78, 89})
		core.validValue = currentBlock2
		core.validRound = currentRound + 1

		// Move to next round and check the expected state
		core.startRound(context.Background(), currentRound+2)

		assert.Equal(t, core.Height(), currentHeight)
		assert.Equal(t, currentRound+2, core.Round())
		assert.Equal(t, propose, core.step)
		assert.Equal(t, currentBlock, core.lockedValue)
		assert.Equal(t, currentRound, core.lockedRound)
		assert.Equal(t, currentBlock2, core.validValue)
		assert.Equal(t, currentRound+1, core.validRound)
	})
}

func TestStartRound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prepareProposal := func(t *testing.T, currentRound int64, proposalHeight *big.Int, validR int64, proposalBlock *types.Block, clientAddress common.Address) (*Message, []byte, []byte) {
		// prepare the proposal message
		proposalRLP, err := Encode(NewProposal(currentRound, proposalHeight, validR, proposalBlock))
		assert.Nil(t, err)
		proposalMsg := &Message{Code: msgProposal, Msg: proposalRLP, Address: clientAddress, Signature: []byte("proposal signature")}
		proposalMsgRLPNoSig, err := proposalMsg.PayloadNoSig()
		assert.Nil(t, err)
		proposalMsgRLPWithSig, err := proposalMsg.Payload()
		assert.Nil(t, err)
		return proposalMsg, proposalMsgRLPNoSig, proposalMsgRLPWithSig
	}

	t.Run("client is the proposer and valid value is nil", func(t *testing.T) {
		committeeSet := prepareCommittee(t)
		members := committeeSet.Committee()
		lastBlockProposer := members[len(members)-1].Address
		clientAddress := members[0].Address

		prevHeight := big.NewInt(int64(rand.Intn(100) + 1))
		prevBlock := generateBlock(prevHeight, types.BlockNonce{1, 2, 3, 4, 5, 6, 7, 8})
		proposalHeight := big.NewInt(prevHeight.Int64() + 1)
		proposalBlock := generateBlock(proposalHeight, types.BlockNonce{12, 23, 34, 45, 56, 67, 78, 89})
		// Ensure clientAddress is the proposer by choosing a round such that
		// r := randomInt * len(currentCommittee)
		// r % len(currentCommittee) = 0
		currentRound := int64(len(members) * (rand.Intn(100)))

		proposalMsg, proposalMsgRLPNoSig, proposalMsgRLPWithSig := prepareProposal(t, currentRound, proposalHeight, int64(-1), proposalBlock, clientAddress)

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddress)

		core := New(backendMock)
		// We assume that round 0 can only happen when we move to a new height, therefore, height is
		// incremented by 1 in start round when round = 0, and the committee set is updated. However, in test case where
		// round is more than 0, then we need to explicitly update the committee set and height.
		if currentRound > 0 {
			core.committeeSet = committeeSet
			core.height = proposalHeight
		}
		core.pendingUnminedBlocks[proposalHeight.Uint64()] = proposalBlock

		if currentRound == 0 {
			// We expect the following extra calls when round = 0
			backendMock.EXPECT().LastCommittedProposal().Return(prevBlock, lastBlockProposer)
			backendMock.EXPECT().Committee(proposalHeight.Uint64()).Return(committeeSet, nil)
		}
		backendMock.EXPECT().SetProposedBlockHash(proposalBlock.Hash())
		backendMock.EXPECT().Sign(proposalMsgRLPNoSig).Return(proposalMsg.Signature, nil)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, proposalMsgRLPWithSig).Return(nil)

		core.startRound(context.Background(), currentRound)

		// There is no need to check for consensus state explicitly here because the broadcasting of proposal message
		// implies an implicit state.
	})

	t.Run("client is the proposer and valid value is not nil", func(t *testing.T) {
		committeeSet := prepareCommittee(t)
		members := committeeSet.Committee()
		clientAddress := members[0].Address

		proposalHeight := big.NewInt(int64(rand.Intn(100) + 1))
		proposalBlock := generateBlock(proposalHeight, types.BlockNonce{11, 22, 33, 44, 55, 66, 77, 88})
		// Valid round can only be set after round 0, hence the smallest value the the round can have is 1 for the valid
		// value to have the smallest value which is 0
		currentRound := int64(len(members) * (rand.Intn(100) + 1))
		validR := currentRound - 1

		proposalMsg, proposalMsgRLPNoSig, proposalMsgRLPWithSig := prepareProposal(t, currentRound, proposalHeight, validR, proposalBlock, clientAddress)

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddress)

		core := New(backendMock)
		core.committeeSet = committeeSet
		core.height = proposalHeight
		core.validRound = validR
		core.validValue = proposalBlock

		backendMock.EXPECT().SetProposedBlockHash(proposalBlock.Hash())
		backendMock.EXPECT().Sign(proposalMsgRLPNoSig).Return(proposalMsg.Signature, nil)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, proposalMsgRLPWithSig).Return(nil)

		core.startRound(context.Background(), currentRound)

		// There is no need to check for consensus state explicitly here because the broadcasting of proposal message
		// implies an implicit state.  Otherwise, the expected message that is to be sent will fail.
	})
	t.Run("client is not the proposer", func(t *testing.T) {
		committeeSet := prepareCommittee(t)
		members := committeeSet.Committee()
		clientIndex := len(members) - 1
		clientAddress := members[clientIndex].Address
		clientPositionInRoundRobin := clientIndex

		prevHeight := big.NewInt(int64(rand.Intn(100) + 1))
		prevBlock := generateBlock(prevHeight, types.BlockNonce{1, 2, 3, 4, 5, 6, 7, 8})
		// ensure the client is not the proposer for current round
		currentRound := int64(rand.Intn(100))
		for (currentRound >= int64(clientIndex)) && (currentRound/int64(clientPositionInRoundRobin) == 0) {
			currentRound = int64(rand.Intn(100))
		}

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddress)

		core := New(backendMock)

		if currentRound > 0 {
			core.committeeSet = committeeSet
		}

		if currentRound == 0 {
			backendMock.EXPECT().LastCommittedProposal().Return(prevBlock, clientAddress)
			backendMock.EXPECT().Committee(big.NewInt(prevHeight.Int64()+1)).Return(committeeSet, nil)
		}

		core.startRound(context.Background(), currentRound)

		assert.Equal(t, currentRound, core.Round())
		assert.True(t, core.proposeTimeout.timerStarted())
	})
}

// TODO: We should create a utility function which can we used across different test files, it can be related to this
// issue https://github.com/clearmatics/autonity/issues/525
// Committee will be ordered such that the proposer for round(n) == committeeSet.members[n % len(committeeSet.members)]
func prepareCommittee(t *testing.T) *committee.Set {
	minSize := 4
	maxSize := 100
	cSize := rand.Intn(maxSize-minSize) + minSize
	c := types.Committee{}
	for i := 1; i <= cSize; i++ {
		key, err := crypto.GenerateKey()
		assert.Nil(t, err)
		member := types.CommitteeMember{
			Address:     crypto.PubkeyToAddress(key.PublicKey),
			VotingPower: new(big.Int).SetInt64(1),
		}
		c = append(c, member)
	}

	sort.Sort(c)
	committeeSet, err := committee.NewSet(c, c[len(c)-1].Address)
	assert.Nil(t, err)
	return committeeSet
}

func generateBlock(height *big.Int, nonce types.BlockNonce) *types.Block {
	// use nonce to create different blocks
	header := &types.Header{Number: height, Nonce: nonce}
	block := types.NewBlock(header, nil, nil, nil)
	return block
}
