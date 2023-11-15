package core

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/autonity/autonity/consensus"
	tdmcommittee "github.com/autonity/autonity/consensus/tendermint/core/committee"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/helpers"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	tctypes "github.com/autonity/autonity/consensus/tendermint/core/types"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/rlp"
	"github.com/autonity/autonity/trie"
	"go.uber.org/mock/gomock"

	"github.com/autonity/autonity/common"
	tcrypto "github.com/autonity/autonity/consensus/tendermint/crypto"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const minSize, maxSize = 4, 100
const timeoutDuration, sleepDuration = 1 * time.Microsecond, 1 * time.Millisecond

func setCommitteeAndSealOnBlock(t *testing.T, b *types.Block, c interfaces.Committee, keys map[common.Address]*ecdsa.PrivateKey, signerIndex int) {
	h := b.Header()
	h.Committee = c.Committee()
	err := tcrypto.SignHeader(h, keys[c.Committee()[signerIndex].Address])
	require.NoError(t, err)
	*b = *b.WithSeal(h)
}

// The following tests aim to test lines 1 - 21 & 57 - 60 of Tendermint Algorithm described on page 6 of
// https://arxiv.org/pdf/1807.04938.pdf.
func TestStartRoundVariables(t *testing.T) {
	prevHeight := big.NewInt(int64(rand.Intn(100) + 1))
	prevBlock := generateBlock(prevHeight)
	currentHeight := big.NewInt(prevHeight.Int64() + 1)
	currentBlock := generateBlock(currentHeight)
	currentRound := int64(0)

	// We don't care who is the next proposer so for simplicity we ensure that clientAddress is not the next
	// proposer by setting clientAddress to be the last proposer. This will ensure that the test will not run the
	// broadcast method.
	committeeSizeAndMaxRound := rand.Intn(maxSize-minSize) + minSize
	committeeSet, keys := prepareCommittee(t, committeeSizeAndMaxRound)

	// This header now needs to be signed  and have a committee to be able to construct a round robin committeeSet.
	setCommitteeAndSealOnBlock(t, prevBlock, committeeSet, keys, 0)
	members := committeeSet.Committee()
	clientAddress := members[len(members)-1].Address

	overrideDefaultCoreValues := func(core *Core) {
		core.height = big.NewInt(-1)
		core.round = int64(-1)
		core.committee = committeeSet
		core.step = tctypes.PrecommitDone
		core.lockedValue = &types.Block{}
		core.lockedRound = 0
		core.validValue = &types.Block{}
		core.validRound = 0
	}

	checkConsensusState := func(t *testing.T, h *big.Int, r int64, s tctypes.Step, lv *types.Block, lr int64, vv *types.Block, vr int64, core *Core) {
		assert.Equal(t, h, core.Height())
		assert.Equal(t, r, core.Round())
		assert.Equal(t, s, core.step)
		assert.Equal(t, lv, core.lockedValue)
		assert.Equal(t, lr, core.lockedRound)
		assert.Equal(t, vv, core.validValue)
		assert.Equal(t, vr, core.validRound)
	}

	t.Run("ensure round 0 state variables are set correctly", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddress)
		backendMock.EXPECT().HeadBlock().Return(prevBlock, clientAddress)
		backendMock.EXPECT().Logger().AnyTimes().Return(log.Root())

		core := New(backendMock, nil)

		overrideDefaultCoreValues(core)
		core.StartRound(context.Background(), currentRound)

		// Check the initial consensus state
		checkConsensusState(t, currentHeight, currentRound, tctypes.Propose, nil, int64(-1), nil, int64(-1), core)

		// stop the timer to clean up
		err := core.proposeTimeout.StopTimer()
		assert.NoError(t, err)
	})
	t.Run("ensure round x state variables are updated correctly", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// In this test we are interested in making sure that that values which change in the current round that may
		// have an impact on the actions performed in the following round (in case of round change) are persisted
		// through to the subsequent round.
		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddress)
		backendMock.EXPECT().HeadBlock().Return(prevBlock, clientAddress).MaxTimes(2)
		backendMock.EXPECT().Logger().AnyTimes().Return(log.Root())

		core := New(backendMock, nil)
		overrideDefaultCoreValues(core)
		core.StartRound(context.Background(), currentRound)

		// Check the initial consensus state
		checkConsensusState(t, currentHeight, currentRound, tctypes.Propose, nil, int64(-1), nil, int64(-1), core)

		// Update locked and valid Value (if locked value changes then valid value also changes, ie quorum(prevotes)
		// delivered in prevote step)
		core.lockedValue = currentBlock
		core.lockedRound = currentRound
		core.validValue = currentBlock
		core.validRound = currentRound

		// Move to next round and check the expected state
		core.StartRound(context.Background(), currentRound+1)

		checkConsensusState(t, currentHeight, currentRound+1, tctypes.Propose, currentBlock, currentRound, currentBlock, currentRound, core)

		// Update valid value (we didn't receive quorum prevote in prevote step, also the block changed, ie, locked
		// value and valid value are different)
		currentBlock2 := generateBlock(currentHeight)
		core.validValue = currentBlock2
		core.validRound = currentRound + 1

		// Move to next round and check the expected state
		core.StartRound(context.Background(), currentRound+2)

		checkConsensusState(t, currentHeight, currentRound+2, tctypes.Propose, currentBlock, currentRound, currentBlock2, currentRound+1, core)

		// stop the timer to clean up
		err := core.proposeTimeout.StopTimer()
		assert.NoError(t, err)
	})
}

func TestStartRound(t *testing.T) {
	// Committee will be ordered such that the proposer for round(n) == committeeSet.members[n % len(committeeSet.members)]
	committeeSizeAndMaxRound := rand.Intn(maxSize-minSize) + minSize
	committeeSet, privateKeys := prepareCommittee(t, committeeSizeAndMaxRound)
	members := committeeSet.Committee()
	clientAddr := members[0].Address

	t.Run("client is the proposer and valid value is nil", func(t *testing.T) {

		lastBlockProposer := members[len(members)-1].Address
		prevHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		prevBlock := generateBlock(prevHeight)
		setCommitteeAndSealOnBlock(t, prevBlock, committeeSet, privateKeys, len(members)-1)

		proposalHeight := big.NewInt(prevHeight.Int64() + 1)
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		for currentRound%int64(len(members)) != 0 {
			currentRound = int64(rand.Intn(committeeSizeAndMaxRound))
		}

		msg, proposal := generateBlockProposal(t, currentRound, proposalHeight, int64(-1), clientAddr, false, privateKeys[clientAddr])

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)
		backendMock.EXPECT().Logger().AnyTimes().Return(log.Root())

		core := New(backendMock, nil)
		core.committee = committeeSet
		// We assume that round 0 can only happen when we move to a new height, therefore, height is
		// incremented by 1 in start round when round = 0, However, in test case where
		// round is more than 0, then we need to explicitly update the height.
		if currentRound > 0 {
			core.height = proposalHeight
		}
		core.pendingCandidateBlocks[proposalHeight.Uint64()] = proposal.ProposalBlock

		if currentRound == 0 {
			// We expect the following extra calls when round = 0
			backendMock.EXPECT().HeadBlock().Return(prevBlock, lastBlockProposer)
		}

		backendMock.EXPECT().SetProposedBlockHash(proposal.ProposalBlock.Hash())
		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().DoAndReturn(signer(privateKeys[clientAddr]))
		backendMock.EXPECT().Broadcast(committeeSet.Committee(), msg.Bytes).Return(nil)

		core.StartRound(context.Background(), currentRound)
	})
	t.Run("client is the proposer and valid value is not nil", func(t *testing.T) {

		proposalHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		// Valid round can only be set after round 0, hence the smallest value the round can have is 1 for the valid
		// value to have the smallest value which is 0
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound) + 1)
		for currentRound%int64(len(members)) != 0 {
			currentRound = int64(rand.Intn(committeeSizeAndMaxRound) + 1)
		}
		validR := currentRound - 1
		msg, proposal := generateBlockProposal(t, currentRound, proposalHeight, validR, clientAddr, false, privateKeys[clientAddr])

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)
		backendMock.EXPECT().Logger().AnyTimes().Return(log.Root())

		core := New(backendMock, nil)
		core.committee = committeeSet
		core.height = proposalHeight
		core.validRound = validR
		core.validValue = proposal.ProposalBlock

		backendMock.EXPECT().SetProposedBlockHash(proposal.ProposalBlock.Hash())
		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().DoAndReturn(signer(privateKeys[clientAddr]))
		backendMock.EXPECT().Broadcast(committeeSet.Committee(), msg.Bytes).Return(nil)

		core.StartRound(context.Background(), currentRound)
	})
	t.Run("client is not the proposer", func(t *testing.T) {
		clientIndex := len(members) - 1
		clientAddr = members[clientIndex].Address

		prevHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		prevBlock := generateBlock(prevHeight)
		// ensure the client is not the proposer for current round
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		for currentRound%int64(clientIndex) == 0 {
			currentRound = int64(rand.Intn(committeeSizeAndMaxRound))
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)
		backendMock.EXPECT().Logger().AnyTimes().Return(log.Root())

		core := New(backendMock, nil)

		if currentRound > 0 {
			core.committee = committeeSet
		}

		if currentRound == 0 {
			backendMock.EXPECT().HeadBlock().Return(prevBlock, clientAddr)
		}

		core.StartRound(context.Background(), currentRound)

		assert.Equal(t, currentRound, core.Round())
		assert.True(t, core.proposeTimeout.TimerStarted())

		// stop the timer to clean up
		err := core.proposeTimeout.StopTimer()
		assert.NoError(t, err)
	})
	t.Run("at proposal Timeout expiry Timeout event is sent", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)
		backendMock.EXPECT().Logger().AnyTimes().Return(log.Root())

		c := New(backendMock, nil)
		c.setCommitteeSet(committeeSet)
		c.setHeight(currentHeight)
		c.setRound(currentRound)

		assert.False(t, c.proposeTimeout.TimerStarted())
		backendMock.EXPECT().Post(tctypes.TimeoutEvent{RoundWhenCalled: currentRound, HeightWhenCalled: currentHeight, Step: consensus.MsgProposal})
		c.prevoteTimeout.ScheduleTimeout(timeoutDuration, c.Round(), c.Height(), c.onTimeoutPropose)
		assert.True(t, c.prevoteTimeout.TimerStarted())
		time.Sleep(sleepDuration)
	})
	t.Run("at reception of proposal Timeout event prevote nil is sent", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		timeoutE := tctypes.TimeoutEvent{RoundWhenCalled: currentRound, HeightWhenCalled: currentHeight, Step: consensus.MsgProposal}
		prevoteMsg, prevoteMsgRLPNoSig, prevoteMsgRLPWithSig := prepareVote(t, consensus.MsgPrevote, currentRound, currentHeight, common.Hash{}, clientAddr, privateKeys[clientAddr])

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)
		backendMock.EXPECT().Logger().AnyTimes().Return(log.Root())

		c := New(backendMock, nil)
		c.setCommitteeSet(committeeSet)
		c.setHeight(currentHeight)
		c.setRound(currentRound)

		backendMock.EXPECT().Sign(prevoteMsgRLPNoSig).Return(prevoteMsg.Signature, nil)
		backendMock.EXPECT().Broadcast(committeeSet.Committee(), prevoteMsgRLPWithSig).Return(nil)

		c.handleTimeoutPropose(context.Background(), timeoutE)
		assert.Equal(t, tctypes.Prevote, c.step)
	})
}

// The following tests aim to test lines 22 - 27 of Tendermint Algorithm described on page 6 of
// https://arxiv.org/pdf/1807.04938.pdf.
func TestNewProposal(t *testing.T) {
	committeeSizeAndMaxRound := rand.Intn(maxSize-minSize) + minSize
	committeeSet, privateKeys := prepareCommittee(t, committeeSizeAndMaxRound)
	members := committeeSet.Committee()
	clientAddr := members[0].Address

	t.Run("receive invalid proposal for current round", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)
		backendMock.EXPECT().Logger().AnyTimes().Return(log.Root())

		c := New(backendMock, nil)
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.SetStep(tctypes.Propose)
		c.setCommitteeSet(committeeSet)

		// members[currentRound] means that the sender is the proposer for the current round
		// assume that the message is from a member of committee set and the signature is signing the contents, however,
		// the proposal block inside the message is invalid
		invalidMsg, invalidProposal := generateBlockProposal(t, currentRound, currentHeight, -1, members[currentRound].Address, true, privateKeys[members[currentRound].Address])

		// prepare prevote nil and target the malicious proposer and the corresponding value.
		prevoteMsg, prevoteMsgRLPNoSig, prevoteMsgRLPWithSig := voteForBadProposal(t, consensus.MsgPrevote, currentRound, currentHeight, clientAddr, privateKeys[clientAddr])

		backendMock.EXPECT().VerifyProposal(invalidProposal.ProposalBlock).Return(time.Duration(1), errors.New("invalid proposal"))
		backendMock.EXPECT().Sign(prevoteMsgRLPNoSig).Return(prevoteMsg.Signature, nil)
		backendMock.EXPECT().Broadcast(committeeSet.Committee(), prevoteMsgRLPWithSig).Return(nil)

		err := c.handleValidMsg(context.Background(), invalidMsg)
		assert.Error(t, err, "expected an error for invalid proposal")
		assert.Equal(t, currentHeight, c.Height())
		assert.Equal(t, currentRound, c.Round())
		assert.Equal(t, tctypes.Prevote, c.step)
	})
	t.Run("receive proposal with validRound = -1 and client's lockedRound = -1", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		clientLockedRound := int64(-1)
		proposalMsg, proposal := generateBlockProposal(t, currentRound, currentHeight, -1, members[currentRound].Address, false, privateKeys[members[currentRound].Address])
		prevoteMsg, prevoteMsgRLPNoSig, prevoteMsgRLPWithSig := prepareVote(t, consensus.MsgPrevote, currentRound, currentHeight, proposal.ProposalBlock.Hash(), clientAddr, privateKeys[clientAddr])

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)
		backendMock.EXPECT().Logger().AnyTimes().Return(log.Root())

		c := New(backendMock, nil)
		// if lockedRround = - 1 then lockedValue = nil
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.SetStep(tctypes.Propose)
		c.setCommitteeSet(committeeSet)
		c.lockedRound = clientLockedRound
		c.lockedValue = nil

		backendMock.EXPECT().VerifyProposal(proposal.ProposalBlock).Return(time.Duration(1), nil)
		backendMock.EXPECT().Sign(prevoteMsgRLPNoSig).Return(prevoteMsg.Signature, nil)
		backendMock.EXPECT().Broadcast(committeeSet.Committee(), prevoteMsgRLPWithSig).Return(nil)

		err := c.handleValidMsg(context.Background(), proposalMsg)
		assert.NoError(t, err)
		assert.Equal(t, currentHeight, c.Height())
		assert.Equal(t, currentRound, c.Round())
		assert.Equal(t, tctypes.Prevote, c.step)
		assert.Nil(t, c.lockedValue)
		assert.Equal(t, clientLockedRound, c.lockedRound)
	})
	t.Run("receive proposal with validRound = -1 and client's lockedValue is same as proposal block", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		clientLockedRound := int64(0)
		proposalMsg, proposal := generateBlockProposal(t, currentRound, currentHeight, -1, members[currentRound].Address, false, privateKeys[members[currentRound].Address])
		prevoteMsg, prevoteMsgRLPNoSig, prevoteMsgRLPWithSig := prepareVote(t, consensus.MsgPrevote, currentRound, currentHeight, proposal.ProposalBlock.Hash(), clientAddr, privateKeys[clientAddr])

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)
		backendMock.EXPECT().Logger().AnyTimes().Return(log.Root())

		c := New(backendMock, nil)
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.SetStep(tctypes.Propose)
		c.setCommitteeSet(committeeSet)
		c.lockedRound = clientLockedRound
		c.lockedValue = proposal.ProposalBlock
		c.validRound = clientLockedRound
		c.validValue = proposal.ProposalBlock

		backendMock.EXPECT().VerifyProposal(proposal.ProposalBlock).Return(time.Duration(1), nil)
		backendMock.EXPECT().Sign(prevoteMsgRLPNoSig).Return(prevoteMsg.Signature, nil)
		backendMock.EXPECT().Broadcast(committeeSet.Committee(), prevoteMsgRLPWithSig).Return(nil)

		err := c.handleValidMsg(context.Background(), proposalMsg)
		assert.NoError(t, err)
		assert.Equal(t, currentHeight, c.Height())
		assert.Equal(t, currentRound, c.Round())
		assert.Equal(t, tctypes.Prevote, c.step)
		assert.Equal(t, proposal.ProposalBlock, c.lockedValue)
		assert.Equal(t, clientLockedRound, c.lockedRound)
		assert.Equal(t, proposal.ProposalBlock, c.validValue)
		assert.Equal(t, clientLockedRound, c.validRound)
	})
	t.Run("receive proposal with validRound = -1 and client's lockedValue is different from proposal block", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		clientLockedRound := int64(0)
		clientLockedValue := generateBlock(currentHeight)
		proposalMsg, proposal := generateBlockProposal(t, currentRound, currentHeight, -1, members[currentRound].Address, false, privateKeys[members[currentRound].Address])
		prevoteMsg, prevoteMsgRLPNoSig, prevoteMsgRLPWithSig := prepareVote(t, consensus.MsgPrevote, currentRound, currentHeight, common.Hash{}, clientAddr, privateKeys[clientAddr])

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)
		backendMock.EXPECT().Logger().AnyTimes().Return(log.Root())

		c := New(backendMock, nil)
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.SetStep(tctypes.Propose)
		c.setCommitteeSet(committeeSet)
		c.lockedRound = clientLockedRound
		c.lockedValue = clientLockedValue
		c.validRound = clientLockedRound
		c.validValue = clientLockedValue

		backendMock.EXPECT().VerifyProposal(proposal.ProposalBlock).Return(time.Duration(1), nil)
		backendMock.EXPECT().Sign(prevoteMsgRLPNoSig).Return(prevoteMsg.Signature, nil)
		backendMock.EXPECT().Broadcast(committeeSet.Committee(), prevoteMsgRLPWithSig).Return(nil)

		err := c.handleValidMsg(context.Background(), proposalMsg)
		assert.NoError(t, err)
		assert.Equal(t, currentHeight, c.Height())
		assert.Equal(t, currentRound, c.Round())
		assert.Equal(t, tctypes.Prevote, c.step)
		assert.Equal(t, clientLockedValue, c.lockedValue)
		assert.Equal(t, clientLockedRound, c.lockedRound)
		assert.Equal(t, clientLockedValue, c.validValue)
		assert.Equal(t, clientLockedRound, c.validRound)
	})
}

// The following tests aim to test lines 28 - 33 of Tendermint Algorithm described on page 6 of
// https://arxiv.org/pdf/1807.04938.pdf.
func TestOldProposal(t *testing.T) {
	//t.Skip("Broken for some random values https://github.com/autonity/autonity/issues/715")
	committeeSizeAndMaxRound := rand.Intn(maxSize-minSize) + minSize
	committeeSet, privateKeys := prepareCommittee(t, committeeSizeAndMaxRound)
	members := committeeSet.Committee()
	clientAddr := members[0].Address

	t.Run("receive proposal with vr >= 0 and client's lockedRound <= vr", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		// vr >= 0 && vr < round_p
		proposalValidRound := int64(0)
		if currentRound > 0 {
			proposalValidRound = int64(rand.Intn(int(currentRound)))
		}
		// -1 <= c.lockedRound <= vr
		clientLockedRound := int64(rand.Intn(int(proposalValidRound+2) - 1))
		proposalMsg, proposal := generateBlockProposal(t, currentRound, currentHeight, proposalValidRound, members[currentRound].Address, false, privateKeys[members[currentRound].Address])
		prevoteMsg, prevoteMsgRLPNoSig, prevoteMsgRLPWithSig := prepareVote(t, consensus.MsgPrevote, currentRound, currentHeight, proposal.ProposalBlock.Hash(), clientAddr, privateKeys[clientAddr])

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)
		backendMock.EXPECT().Logger().AnyTimes().Return(log.Root())

		c := New(backendMock, nil)
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.SetStep(tctypes.Propose)
		c.setCommitteeSet(committeeSet)
		c.lockedRound = clientLockedRound
		c.validRound = clientLockedRound
		// Although the following is not possible it is required to ensure that c.lockRound <= proposalValidRound is
		// responsible for sending the prevote for the incoming proposal
		c.lockedValue = nil
		c.validValue = nil
		c.messages.GetOrCreate(proposalValidRound).AddPrevote(proposal.ProposalBlock.Hash(), message.Message{Code: consensus.MsgPrevote, Power: c.CommitteeSet().Quorum()})

		backendMock.EXPECT().VerifyProposal(proposal.ProposalBlock).Return(time.Duration(1), nil)
		backendMock.EXPECT().Sign(prevoteMsgRLPNoSig).Return(prevoteMsg.Signature, nil)
		backendMock.EXPECT().Broadcast(committeeSet.Committee(), prevoteMsgRLPWithSig).Return(nil)

		err := c.handleValidMsg(context.Background(), proposalMsg)
		assert.NoError(t, err)
		assert.Equal(t, currentHeight, c.Height())
		assert.Equal(t, currentRound, c.Round())
		assert.Equal(t, tctypes.Prevote, c.step)
		assert.Nil(t, c.lockedValue)
		assert.Equal(t, clientLockedRound, c.lockedRound)
		assert.Nil(t, c.validValue)
		assert.Equal(t, clientLockedRound, c.validRound)
	})
	t.Run("receive proposal with vr >= 0 and client's lockedValue is same as proposal block", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		// vr >= 0 && vr < round_p
		proposalValidRound := int64(0)
		if currentRound != 0 {
			proposalValidRound = int64(rand.Intn(int(currentRound)))
		}
		t.Log("currentHeight", currentHeight, "currentRound", currentRound, "proposalValidRound", proposalValidRound)
		proposalMsg, proposal := generateBlockProposal(t, currentRound, currentHeight, proposalValidRound, members[currentRound].Address, false, privateKeys[members[currentRound].Address])
		prevoteMsg, prevoteMsgRLPNoSig, prevoteMsgRLPWithSig := prepareVote(t, consensus.MsgPrevote, currentRound, currentHeight, proposal.ProposalBlock.Hash(), clientAddr, privateKeys[clientAddr])

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)
		backendMock.EXPECT().Logger().AnyTimes().Return(log.Root())

		c := New(backendMock, nil)
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.SetStep(tctypes.Propose)
		c.setCommitteeSet(committeeSet)
		// Although the following is not possible it is required to ensure that c.lockedValue = proposal is responsible
		// for sending the prevote for the incoming proposal
		c.lockedRound = proposalValidRound + 1
		c.validRound = proposalValidRound + 1
		c.lockedValue = proposal.ProposalBlock
		c.validValue = proposal.ProposalBlock
		c.messages.GetOrCreate(proposalValidRound).AddPrevote(proposal.ProposalBlock.Hash(), message.Message{Code: consensus.MsgPrevote, Power: c.CommitteeSet().Quorum()})

		backendMock.EXPECT().VerifyProposal(proposal.ProposalBlock).Return(time.Duration(1), nil)
		backendMock.EXPECT().Sign(prevoteMsgRLPNoSig).Return(prevoteMsg.Signature, nil)
		backendMock.EXPECT().Broadcast(committeeSet.Committee(), prevoteMsgRLPWithSig).Return(nil)

		err := c.handleValidMsg(context.Background(), proposalMsg)
		assert.NoError(t, err)
		assert.Equal(t, currentHeight, c.Height())
		assert.Equal(t, currentRound, c.Round())
		assert.Equal(t, tctypes.Prevote, c.step)
		assert.Equal(t, proposalValidRound+1, c.lockedRound)
		assert.Equal(t, proposalValidRound+1, c.validRound)
		assert.Equal(t, proposal.ProposalBlock, c.lockedValue)
		assert.Equal(t, proposal.ProposalBlock, c.validValue)
	})
	t.Run("receive proposal with vr >= 0 and clients is lockedRound > vr with a different value", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound)) + 1 //+1 to prevent 0 passed to randoms
		clientLockedValue := generateBlock(currentHeight)
		// vr >= 0 && vr < round_p
		proposalValidRound := int64(rand.Intn(int(currentRound)))
		proposalMsg, proposal := generateBlockProposal(t, currentRound, currentHeight, proposalValidRound, members[currentRound].Address, false, privateKeys[members[currentRound].Address])
		prevoteMsg, prevoteMsgRLPNoSig, prevoteMsgRLPWithSig := prepareVote(t, consensus.MsgPrevote, currentRound, currentHeight, common.Hash{}, clientAddr, privateKeys[clientAddr])

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)
		backendMock.EXPECT().Logger().AnyTimes().Return(log.Root())

		c := New(backendMock, nil)
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.SetStep(tctypes.Propose)
		c.setCommitteeSet(committeeSet)
		// Although the following is not possible it is required to ensure that c.lockedValue = proposal is responsible
		// for sending the prevote for the incoming proposal
		c.lockedRound = proposalValidRound + 1
		c.validRound = proposalValidRound + 1
		c.lockedValue = clientLockedValue
		c.validValue = clientLockedValue
		c.messages.GetOrCreate(proposalValidRound).AddPrevote(proposal.ProposalBlock.Hash(), message.Message{Code: consensus.MsgPrevote, Power: c.CommitteeSet().Quorum()})

		backendMock.EXPECT().VerifyProposal(proposal.ProposalBlock).Return(time.Duration(1), nil)
		backendMock.EXPECT().Sign(prevoteMsgRLPNoSig).Return(prevoteMsg.Signature, nil)
		backendMock.EXPECT().Broadcast(committeeSet.Committee(), prevoteMsgRLPWithSig).Return(nil)

		err := c.handleValidMsg(context.Background(), proposalMsg)
		assert.NoError(t, err)
		assert.Equal(t, currentHeight, c.Height())
		assert.Equal(t, currentRound, c.Round())
		assert.Equal(t, tctypes.Prevote, c.step)
		assert.Equal(t, proposalValidRound+1, c.lockedRound)
		assert.Equal(t, proposalValidRound+1, c.validRound)
		assert.Equal(t, clientLockedValue, c.lockedValue)
		assert.Equal(t, clientLockedValue, c.validValue)
	})

	// line 28 check upon condition on prevote handler.
	/*
		Please refer to the discussion history: https://github.com/autonity/autonity/pull/615
		This test case validates the need for the upon condition defined in line 28 of the tendermint
		white paper, which addresses the scenario laid out below.

		Round 0:
			Quorum clients sent precommit nil hence there was a round change
		Round 1:
			The proposer of round 1 had received quorum prevotes for the proposal of round 0, thus, the proposer of round 1
		re-proposes the proposal with vr = 0, More than quorum of the network is yet to receive the prevotes from the
		previous round, thus there are their lockedValue = nil and lockedRound = -1. The proposer of round 1 has a really good
		connection with the rest of the network, thus it is able to send the proposal to its peers before they receive enough
		prevotes from round 0 to form a quorum. Now the proposal received by the peers would only be able to satisfy line
		28 of the Tendermint pseudo code, however, the prevotes from the previous round are yet to arrive. Since the
		proposal of the current round has been received the timer would be stopped.

		autonity/consensus/tendermint/Core/propose.go, Lines 131 to 133 at 78f199d

		 if err := c.proposeTimeout.stopTimer(); err != nil {
		 	return err
		 }

		A quorum prevote for round 0 finally arrive, however, these will be added to message set and without
		the existence of the line 28 upon condition nothing would happen, even though enough  messages
		are present in the message set to send a prevote for the old proposal.

		This was previously the case in:
		autonity/consensus/tendermint/Core/prevote.go, Lines 69 to 74 at 78f199d

		 if err == errOldRoundMessage {
		 	// We only process old rounds while future rounds messages are pushed on to the backlog
		 	oldRoundMessages := c.messages.getOrCreate(preVote.Round)
		 	c.acceptVote(oldRoundMessages, prevote, preVote.ProposedBlockHash, *msg)
		 }
		 return err

		Without the line 28 upon condition the client is stuck since the timer has been stopped, thus a prevote nil
		cannot be sent and the timer cannot be restarted until startRound() is called for a new round. The
		resending of the message set will also not help because it would only send messages to peers which they
		haven't seen and since there were no new messages the peers will not be able to make progress. This can
		also happen where the client's lockedRound < vr, it cannot happen for lockedRound = vr because that means
		the client had received enough prevote in a timely manner and there are no old prevote to arrive.

		Therefore we had a liveness bug in implementations of Tendermint in commits prior to this one.
	*/
	t.Run("handle proposal before full quorum prevote on valid round is satisfied, exe action by applying old round prevote into round state", func(t *testing.T) {
		clientIndex := len(members) - 1
		clientAddr = members[clientIndex].Address

		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))

		// ensure the client is not the proposer for current round
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		for currentRound%int64(clientIndex) == 0 {
			currentRound = int64(rand.Intn(committeeSizeAndMaxRound))
		}

		// vr >= 0 && vr < round_p
		proposalValidRound := int64(0)
		if currentRound > 0 {
			proposalValidRound = int64(rand.Intn(int(currentRound)))
		}

		// -1 <= c.lockedRound < vr, if the client lockedValue = vr then the client had received the prevotes in a
		// timely manner thus there are no old prevote yet to arrive
		clientLockedRound := int64(-1)
		if proposalValidRound > 0 {
			clientLockedRound = int64(rand.Intn(int(proposalValidRound)) - 1)
		}

		// the new round proposal
		proposalMsg, proposal := generateBlockProposal(t, currentRound, currentHeight, proposalValidRound, members[currentRound].Address, false, privateKeys[members[currentRound].Address])

		// old proposal some random block
		clientLockedValue := generateBlock(currentHeight)

		// the old round prevote msg to be handled to get the full quorum prevote on old round vr with value v.
		prevoteMsg, _, _ := prepareVote(t, consensus.MsgPrevote, proposalValidRound, currentHeight, proposal.ProposalBlock.Hash(), clientAddr, privateKeys[clientAddr])

		// the expected prevote msg to be broadcast for the new round with <currentHeight, currentRound, proposal.ProposalBlock.Hash()>
		prevoteMsgToBroadcast, prevoteMsgRLPNoSig, prevoteMsgRLPWithSig := prepareVote(t, consensus.MsgPrevote, currentRound, currentHeight, proposal.ProposalBlock.Hash(), clientAddr, privateKeys[clientAddr])

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)
		backendMock.EXPECT().Logger().AnyTimes().Return(log.Root())

		c := New(backendMock, nil)
		c.setCommitteeSet(committeeSet)
		// construct round state with: old round's quorum-1 prevote for v on valid round.
		c.messages.GetOrCreate(proposalValidRound).AddPrevote(proposal.ProposalBlock.Hash(), message.Message{Code: consensus.MsgPrevote, Power: new(big.Int).Sub(c.CommitteeSet().Quorum(), common.Big1)})

		// client on new round's step propose.
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.SetStep(tctypes.Propose)
		c.lockedRound = clientLockedRound
		c.validRound = clientLockedRound
		c.lockedValue = clientLockedValue
		c.validValue = clientLockedValue

		//schedule the proposer Timeout since the client is not the proposer for this round
		c.proposeTimeout.ScheduleTimeout(1*time.Second, c.Round(), c.Height(), c.onTimeoutPropose)

		backendMock.EXPECT().VerifyProposal(proposal.ProposalBlock).Return(time.Duration(1), nil)
		backendMock.EXPECT().Sign(prevoteMsgRLPNoSig).Return(prevoteMsgToBroadcast.Signature, nil)
		backendMock.EXPECT().Broadcast(committeeSet.Committee(), prevoteMsgRLPWithSig).Return(nil)

		// now we handle new round's proposal with round_p > vr on value v.
		err := c.handleValidMsg(context.Background(), proposalMsg)
		assert.NoError(t, err)

		// check timer was stopped after receiving the proposal
		assert.False(t, c.proposeTimeout.TimerStarted())

		// now we receive the last old round's prevote MSG to get quorum prevote on vr for value v.
		// the old round's prevote is accepted into the round state which now have the line 28 condition satisfied.
		// now to take the action of line 28 which was not align with pseudo code before.

		err = c.handleValidMsg(context.Background(), prevoteMsg)
		assert.NoError(t, err)
		assert.Equal(t, currentHeight, c.Height())
		assert.Equal(t, currentRound, c.Round())
		assert.Equal(t, tctypes.Prevote, c.step)
		assert.Equal(t, clientLockedValue, c.lockedValue)
		assert.Equal(t, clientLockedRound, c.lockedRound)
		assert.Equal(t, clientLockedValue, c.validValue)
		assert.Equal(t, clientLockedRound, c.validRound)
	})
}

// The following tests aim to test lines 34 - 35 & 61 - 64 of Tendermint Algorithm described on page 6 of
// https://arxiv.org/pdf/1807.04938.pdf.
func TestPrevoteTimeout(t *testing.T) {
	committeeSizeAndMaxRound := rand.Intn(maxSize-minSize) + minSize
	committeeSet, privateKeys := prepareCommittee(t, committeeSizeAndMaxRound)
	members := committeeSet.Committee()
	clientAddr := members[0].Address

	t.Run("prevote Timeout started after quorum of prevotes with different hashes", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		sender := 1
		prevoteMsg, _, _ := prepareVote(t, consensus.MsgPrevote, currentRound, currentHeight, generateBlock(currentHeight).Hash(), members[sender].Address, privateKeys[members[sender].Address])

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)
		backendMock.EXPECT().Logger().AnyTimes().Return(log.Root())

		c := New(backendMock, nil)
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.SetStep(tctypes.Prevote)
		c.setCommitteeSet(committeeSet)
		// create quorum prevote messages however there is no quorum on a specific hash
		c.curRoundMessages.AddPrevote(common.Hash{}, message.Message{Address: members[2].Address, Code: consensus.MsgPrevote, Power: new(big.Int).Sub(c.CommitteeSet().Quorum(), common.Big2)})
		c.curRoundMessages.AddPrevote(generateBlock(currentHeight).Hash(), message.Message{Address: members[3].Address, Code: consensus.MsgPrevote, Power: common.Big1})

		assert.False(t, c.prevoteTimeout.TimerStarted())
		err := c.handleValidMsg(context.Background(), prevoteMsg)
		assert.NoError(t, err)
		assert.True(t, c.prevoteTimeout.TimerStarted())

		// stop the timer to clean up
		err = c.prevoteTimeout.StopTimer()
		assert.NoError(t, err)

		assert.Equal(t, currentHeight, c.Height())
		assert.Equal(t, currentRound, c.Round())
		assert.Equal(t, tctypes.Prevote, c.step)
	})
	t.Run("prevote Timeout is not started multiple times", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		sender1 := 1
		prevote1Msg, _, _ := prepareVote(t, consensus.MsgPrevote, currentRound, currentHeight, generateBlock(currentHeight).Hash(), members[sender1].Address, privateKeys[members[sender1].Address])
		sender2 := 2
		prevote2Msg, _, _ := prepareVote(t, consensus.MsgPrevote, currentRound, currentHeight, generateBlock(currentHeight).Hash(), members[sender2].Address, privateKeys[members[sender2].Address])

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)
		backendMock.EXPECT().Logger().AnyTimes().Return(log.Root())

		c := New(backendMock, nil)
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.SetStep(tctypes.Prevote)
		c.setCommitteeSet(committeeSet)
		// create quorum prevote messages however there is no quorum on a specific hash
		c.curRoundMessages.AddPrevote(common.Hash{}, message.Message{Address: members[3].Address, Code: consensus.MsgPrevote, Power: new(big.Int).Sub(c.CommitteeSet().Quorum(), common.Big2)})
		c.curRoundMessages.AddPrevote(generateBlock(currentHeight).Hash(), message.Message{Address: members[0].Address, Code: consensus.MsgPrevote, Power: common.Big1})

		assert.False(t, c.prevoteTimeout.TimerStarted())

		err := c.handleValidMsg(context.Background(), prevote1Msg)
		assert.NoError(t, err)
		assert.True(t, c.prevoteTimeout.TimerStarted())

		timeNow := time.Now()

		err = c.handleValidMsg(context.Background(), prevote2Msg)
		assert.NoError(t, err)
		assert.True(t, c.prevoteTimeout.TimerStarted())
		assert.True(t, c.prevoteTimeout.Start.Before(timeNow))

		// stop the timer to clean up
		err = c.prevoteTimeout.StopTimer()
		assert.NoError(t, err)

		assert.Equal(t, currentHeight, c.Height())
		assert.Equal(t, currentRound, c.Round())
		assert.Equal(t, tctypes.Prevote, c.step)
	})
	t.Run("at prevote Timeout expiry Timeout event is sent", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)
		backendMock.EXPECT().Logger().AnyTimes().Return(log.Root())

		c := New(backendMock, nil)
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.SetStep(tctypes.Prevote)
		c.setCommitteeSet(committeeSet)

		assert.False(t, c.prevoteTimeout.TimerStarted())
		backendMock.EXPECT().Post(tctypes.TimeoutEvent{RoundWhenCalled: currentRound, HeightWhenCalled: currentHeight, Step: consensus.MsgPrevote})
		c.prevoteTimeout.ScheduleTimeout(timeoutDuration, c.Round(), c.Height(), c.onTimeoutPrevote)
		assert.True(t, c.prevoteTimeout.TimerStarted())
		assert.Equal(t, currentHeight, c.Height())
		assert.Equal(t, currentRound, c.Round())
		assert.Equal(t, tctypes.Prevote, c.step)
		time.Sleep(sleepDuration)
	})
	t.Run("at reception of prevote Timeout event precommit nil is sent", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		timeoutE := tctypes.TimeoutEvent{RoundWhenCalled: currentRound, HeightWhenCalled: currentHeight, Step: consensus.MsgPrevote}
		precommitMsg, precommitMsgRLPNoSig, precommitMsgRLPWithSig := prepareVote(t, consensus.MsgPrecommit, currentRound, currentHeight, common.Hash{}, clientAddr, privateKeys[clientAddr])
		committedSeal := helpers.PrepareCommittedSeal(common.Hash{}, currentRound, currentHeight)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)
		backendMock.EXPECT().Logger().AnyTimes().Return(log.Root())

		c := New(backendMock, nil)
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.SetStep(tctypes.Prevote)
		c.setCommitteeSet(committeeSet)

		backendMock.EXPECT().Sign(committedSeal).Return(precommitMsg.CommittedSeal, nil)
		backendMock.EXPECT().Sign(precommitMsgRLPNoSig).Return(precommitMsg.Signature, nil)
		backendMock.EXPECT().Broadcast(committeeSet.Committee(), precommitMsgRLPWithSig).Return(nil)

		c.handleTimeoutPrevote(context.Background(), timeoutE)
		assert.Equal(t, currentHeight, c.Height())
		assert.Equal(t, currentRound, c.Round())
		assert.Equal(t, tctypes.Precommit, c.step)
	})
}

// The following tests aim to test lines 34 - 43 of Tendermint Algorithm described on page 6 of
// https://arxiv.org/pdf/1807.04938.pdf.
func TestQuorumPrevote(t *testing.T) {
	committeeSizeAndMaxRound := rand.Intn(maxSize-minSize) + minSize
	committeeSet, privateKeys := prepareCommittee(t, committeeSizeAndMaxRound)
	members := committeeSet.Committee()
	clientAddr := members[0].Address

	t.Run("receive quroum prevote for proposal block when in step >= prevote", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		//randomly choose prevote or precommit step
		currentStep := tctypes.Step(rand.Intn(2) + 1)                                                                                                                                                           //nolint:gosec
		proposalMsg, proposal := generateBlockProposal(t, currentRound, currentHeight, int64(rand.Intn(int(currentRound+1))), members[currentRound].Address, false, privateKeys[members[currentRound].Address]) //nolint:gosec
		sender := 1
		prevoteMsg, _, _ := prepareVote(t, consensus.MsgPrevote, currentRound, currentHeight, proposal.ProposalBlock.Hash(), members[sender].Address, privateKeys[members[sender].Address])
		precommitMsg, precommitMsgRLPNoSig, precommitMsgRLPWithSig := prepareVote(t, consensus.MsgPrecommit, currentRound, currentHeight, proposal.ProposalBlock.Hash(), clientAddr, privateKeys[clientAddr])

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)
		backendMock.EXPECT().Logger().AnyTimes().Return(log.Root())

		c := New(backendMock, nil)
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.SetStep(currentStep)
		c.setCommitteeSet(committeeSet)
		c.curRoundMessages.SetProposal(proposal, proposalMsg, true)
		c.curRoundMessages.AddPrevote(proposal.ProposalBlock.Hash(), message.Message{Address: members[2].Address, Code: consensus.MsgPrevote, Power: new(big.Int).Sub(c.CommitteeSet().Quorum(), common.Big1)})

		if currentStep == tctypes.Prevote {
			committedSeal := helpers.PrepareCommittedSeal(proposal.ProposalBlock.Hash(), currentRound, currentHeight)

			backendMock.EXPECT().Sign(committedSeal).Return(precommitMsg.CommittedSeal, nil)
			backendMock.EXPECT().Sign(precommitMsgRLPNoSig).Return(precommitMsg.Signature, nil)
			backendMock.EXPECT().Broadcast(committeeSet.Committee(), precommitMsgRLPWithSig).Return(nil)

			err := c.handleValidMsg(context.Background(), prevoteMsg)
			assert.NoError(t, err)

			assert.Equal(t, proposal.ProposalBlock, c.lockedValue)
			assert.Equal(t, currentRound, c.lockedRound)
			assert.Equal(t, tctypes.Precommit, c.step)

		} else if currentStep == tctypes.Precommit {
			err := c.handleValidMsg(context.Background(), prevoteMsg)
			assert.NoError(t, err)

			assert.Equal(t, proposal.ProposalBlock, c.validValue)
			assert.Equal(t, currentRound, c.validRound)
		}
		assert.Equal(t, currentHeight, c.Height())
		assert.Equal(t, currentRound, c.Round())
	})

	t.Run("receive more than quorum prevote for proposal block when in step >= prevote", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		//randomly choose prevote or precommit step
		currentStep := tctypes.Step(rand.Intn(2) + 1)                                                                                                                                    //nolint:gosec
		proposalMsg, proposal := generateBlockProposal(t, currentRound, currentHeight, currentRound-1, members[currentRound].Address, false, privateKeys[members[currentRound].Address]) //nolint:gosec
		sender1 := 1
		prevoteMsg1, _, _ := prepareVote(t, consensus.MsgPrevote, currentRound, currentHeight, proposal.ProposalBlock.Hash(), members[sender1].Address, privateKeys[members[sender1].Address])
		sender2 := 2
		prevoteMsg2, _, _ := prepareVote(t, consensus.MsgPrevote, currentRound, currentHeight, proposal.ProposalBlock.Hash(), members[sender2].Address, privateKeys[members[sender2].Address])
		precommitMsg, precommitMsgRLPNoSig, precommitMsgRLPWithSig := prepareVote(t, consensus.MsgPrecommit, currentRound, currentHeight, proposal.ProposalBlock.Hash(), clientAddr, privateKeys[clientAddr])

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)
		backendMock.EXPECT().Logger().AnyTimes().Return(log.Root())

		c := New(backendMock, nil)
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.SetStep(currentStep)
		c.setCommitteeSet(committeeSet)
		c.curRoundMessages.SetProposal(proposal, proposalMsg, true)
		c.curRoundMessages.AddPrevote(proposal.ProposalBlock.Hash(), message.Message{Address: members[3].Address, Code: consensus.MsgPrevote, Power: new(big.Int).Sub(c.CommitteeSet().Quorum(), common.Big1)})

		// receive first prevote to increase the total to quorum
		if currentStep == tctypes.Prevote {
			committedSeal := helpers.PrepareCommittedSeal(proposal.ProposalBlock.Hash(), currentRound, currentHeight)

			backendMock.EXPECT().Sign(committedSeal).Return(precommitMsg.CommittedSeal, nil)
			backendMock.EXPECT().Sign(precommitMsgRLPNoSig).Return(precommitMsg.Signature, nil)
			backendMock.EXPECT().Broadcast(committeeSet.Committee(), precommitMsgRLPWithSig).Return(nil)

			err := c.handleValidMsg(context.Background(), prevoteMsg1)
			assert.NoError(t, err)

			assert.Equal(t, proposal.ProposalBlock, c.lockedValue)
			assert.Equal(t, currentRound, c.lockedRound)
			assert.Equal(t, tctypes.Precommit, c.step)

		} else if currentStep == tctypes.Precommit {
			err := c.handleValidMsg(context.Background(), prevoteMsg1)
			assert.NoError(t, err)

			assert.Equal(t, proposal.ProposalBlock, c.validValue)
			assert.Equal(t, currentRound, c.validRound)
		}

		// receive second prevote to increase the total to more than quorum
		lockedValueBefore := c.lockedValue
		validValueBefore := c.validValue
		lockedRoundBefore := c.lockedRound
		validRoundBefore := c.validRound

		err := c.handleValidMsg(context.Background(), prevoteMsg2)
		assert.NoError(t, err)

		assert.Equal(t, currentHeight, c.Height())
		assert.Equal(t, currentRound, c.Round())
		assert.Equal(t, lockedValueBefore, c.lockedValue)
		assert.Equal(t, validValueBefore, c.validValue)
		assert.Equal(t, lockedRoundBefore, c.lockedRound)
		assert.Equal(t, validRoundBefore, c.validRound)
	})
}

// The following tests aim to test lines 44 - 46 of Tendermint Algorithm described on page 6 of
// https://arxiv.org/pdf/1807.04938.pdf.
func TestQuorumPrevoteNil(t *testing.T) {
	committeeSizeAndMaxRound := rand.Intn(maxSize-minSize) + minSize
	committeeSet, privateKeys := prepareCommittee(t, committeeSizeAndMaxRound)
	members := committeeSet.Committee()
	clientAddr := members[0].Address

	currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
	currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
	sender := 1
	prevoteMsg, _, _ := prepareVote(t, consensus.MsgPrevote, currentRound, currentHeight, common.Hash{}, members[sender].Address, privateKeys[members[sender].Address])
	precommitMsg, precommitMsgRLPNoSig, precommitMsgRLPWithSig := prepareVote(t, consensus.MsgPrecommit, currentRound, currentHeight, common.Hash{}, clientAddr, privateKeys[clientAddr])
	committedSeal := helpers.PrepareCommittedSeal(common.Hash{}, currentRound, currentHeight)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	backendMock := interfaces.NewMockBackend(ctrl)
	backendMock.EXPECT().Address().Return(clientAddr)
	backendMock.EXPECT().Logger().AnyTimes().Return(log.Root())

	c := New(backendMock, nil)
	c.setHeight(currentHeight)
	c.setRound(currentRound)
	c.SetStep(tctypes.Prevote)
	c.setCommitteeSet(committeeSet)
	c.curRoundMessages.AddPrevote(common.Hash{}, message.Message{Address: members[2].Address, Code: consensus.MsgPrevote, Power: new(big.Int).Sub(c.CommitteeSet().Quorum(), common.Big1)})

	backendMock.EXPECT().Sign(committedSeal).Return(precommitMsg.CommittedSeal, nil)
	backendMock.EXPECT().Sign(precommitMsgRLPNoSig).Return(precommitMsg.Signature, nil)
	backendMock.EXPECT().Broadcast(committeeSet.Committee(), precommitMsgRLPWithSig).Return(nil)

	err := c.handleValidMsg(context.Background(), prevoteMsg)
	assert.NoError(t, err)

	assert.Equal(t, currentHeight, c.Height())
	assert.Equal(t, currentRound, c.Round())
	assert.Equal(t, tctypes.Precommit, c.step)
}

// The following tests aim to test lines 47 - 48 & 65 - 67 of Tendermint Algorithm described on page 6 of
// https://arxiv.org/pdf/1807.04938.pdf.
func TestPrecommitTimeout(t *testing.T) {
	committeeSizeAndMaxRound := rand.Intn(maxSize-minSize) + minSize
	committeeSet, privateKeys := prepareCommittee(t, committeeSizeAndMaxRound)
	members := committeeSet.Committee()
	clientAddr := members[0].Address

	t.Run("precommit Timeout started after quorum of precommits with different hashes", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		sender := 1
		precommitMsg, _, _ := prepareVote(t, consensus.MsgPrecommit, currentRound, currentHeight, generateBlock(currentHeight).Hash(), members[sender].Address, privateKeys[members[sender].Address])

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)
		backendMock.EXPECT().Logger().AnyTimes().Return(log.Root())

		c := New(backendMock, nil)
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		//TODO: this should be changed to Step(rand.Intn(3)) to make sure precommit Timeout can be started from any step
		c.SetStep(tctypes.Precommit)
		c.setCommitteeSet(committeeSet)
		// create quorum precommit messages however there is no quorum on a specific hash
		c.curRoundMessages.AddPrecommit(common.Hash{}, message.Message{Address: members[2].Address, Code: consensus.MsgPrecommit, Power: new(big.Int).Sub(c.CommitteeSet().Quorum(), common.Big2)})
		c.curRoundMessages.AddPrecommit(generateBlock(currentHeight).Hash(), message.Message{Address: members[3].Address, Code: consensus.MsgPrecommit, Power: common.Big1})

		assert.False(t, c.precommitTimeout.TimerStarted())
		err := c.handleValidMsg(context.Background(), precommitMsg)
		assert.NoError(t, err)
		assert.True(t, c.precommitTimeout.TimerStarted())

		// stop the timer to clean up
		err = c.precommitTimeout.StopTimer()
		assert.NoError(t, err)

		assert.Equal(t, currentHeight, c.Height())
		assert.Equal(t, currentRound, c.Round())
		assert.Equal(t, tctypes.Precommit, c.step)
	})
	t.Run("precommit Timeout is not started multiple times", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		sender1 := 1
		precommit1Msg, _, _ := prepareVote(t, consensus.MsgPrecommit, currentRound, currentHeight, generateBlock(currentHeight).Hash(), members[sender1].Address, privateKeys[members[sender1].Address])
		sender2 := 2
		precommit2Msg, _, _ := prepareVote(t, consensus.MsgPrecommit, currentRound, currentHeight, generateBlock(currentHeight).Hash(), members[sender2].Address, privateKeys[members[sender2].Address])

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)
		backendMock.EXPECT().Logger().AnyTimes().Return(log.Root())

		c := New(backendMock, nil)
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		//TODO: this should be changed to Step(rand.Intn(3)) to make sure precommit Timeout can be started from any step
		c.SetStep(tctypes.Precommit)
		c.setCommitteeSet(committeeSet)
		// create quorum prevote messages however there is no quorum on a specific hash
		c.curRoundMessages.AddPrecommit(common.Hash{}, message.Message{Address: members[3].Address, Code: consensus.MsgPrecommit, Power: new(big.Int).Sub(c.CommitteeSet().Quorum(), common.Big2)})
		c.curRoundMessages.AddPrecommit(generateBlock(currentHeight).Hash(), message.Message{Address: members[0].Address, Code: consensus.MsgPrecommit, Power: common.Big1})

		assert.False(t, c.precommitTimeout.TimerStarted())

		err := c.handleValidMsg(context.Background(), precommit1Msg)
		assert.NoError(t, err)
		assert.True(t, c.precommitTimeout.TimerStarted())

		timeNow := time.Now()

		err = c.handleValidMsg(context.Background(), precommit2Msg)
		assert.NoError(t, err)
		assert.True(t, c.precommitTimeout.TimerStarted())
		assert.True(t, c.precommitTimeout.Start.Before(timeNow))

		// stop the timer to clean up
		err = c.precommitTimeout.StopTimer()
		assert.NoError(t, err)

		assert.Equal(t, currentHeight, c.Height())
		assert.Equal(t, currentRound, c.Round())
		assert.Equal(t, tctypes.Precommit, c.step)
	})
	t.Run("at precommit Timeout expiry Timeout event is sent", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)
		backendMock.EXPECT().Logger().AnyTimes().Return(log.Root())

		c := New(backendMock, nil)
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		//TODO: this should be changed to Step(rand.Intn(3)) to make sure precommit Timeout can be started from any step
		c.SetStep(tctypes.Precommit)
		c.setCommitteeSet(committeeSet)

		assert.False(t, c.precommitTimeout.TimerStarted())
		backendMock.EXPECT().Post(tctypes.TimeoutEvent{RoundWhenCalled: currentRound, HeightWhenCalled: currentHeight, Step: consensus.MsgPrecommit})
		c.precommitTimeout.ScheduleTimeout(timeoutDuration, c.Round(), c.Height(), c.onTimeoutPrecommit)
		assert.True(t, c.precommitTimeout.TimerStarted())
		assert.Equal(t, currentHeight, c.Height())
		assert.Equal(t, currentRound, c.Round())
		assert.Equal(t, tctypes.Precommit, c.step)
		time.Sleep(sleepDuration)
	})
	t.Run("at reception of precommit Timeout event next round will be started", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		// ensure client is not the proposer for next round
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		for (currentRound+1)%int64(len(members)) == 0 {
			currentRound = int64(rand.Intn(committeeSizeAndMaxRound))
		}
		timeoutE := tctypes.TimeoutEvent{RoundWhenCalled: currentRound, HeightWhenCalled: currentHeight, Step: consensus.MsgPrecommit}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)
		backendMock.EXPECT().Logger().AnyTimes().Return(log.Root())

		c := New(backendMock, nil)
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		//TODO: this should be changed to Step(rand.Intn(3)) to make sure precommit Timeout can be started from any step
		c.SetStep(tctypes.Precommit)
		c.setCommitteeSet(committeeSet)

		c.handleTimeoutPrecommit(context.Background(), timeoutE)

		assert.Equal(t, currentHeight, c.Height())
		assert.Equal(t, currentRound+1, c.Round())
		assert.Equal(t, tctypes.Propose, c.step)

		// stop the timer to clean up, since start round can start propose Timeout
		err := c.proposeTimeout.StopTimer()
		assert.NoError(t, err)
	})
}

// The following tests aim to test lines 49 - 54 of Tendermint Algorithm described on page 6 of
// https://arxiv.org/pdf/1807.04938.pdf.
func TestQuorumPrecommit(t *testing.T) {
	committeeSizeAndMaxRound := rand.Intn(maxSize-minSize) + minSize
	committeeSet, privateKeys := prepareCommittee(t, committeeSizeAndMaxRound)
	members := committeeSet.Committee()
	clientAddr := members[0].Address
	currentHeight := big.NewInt(int64(rand.Intn(maxSize+1) + 1))
	nextHeight := currentHeight.Uint64() + 1
	t.Log("committee size", committeeSizeAndMaxRound, "current height", currentHeight)
	nextProposalMsg, nextP := generateBlockProposal(t, 0, big.NewInt(int64(nextHeight)), int64(-1), clientAddr, false, privateKeys[clientAddr])

	currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
	proposalMsg, proposal := generateBlockProposal(t, currentRound, currentHeight, currentRound, members[currentRound].Address, false, privateKeys[members[currentRound].Address]) //nolint:gosec
	sender := 1
	precommitMsg, _, _ := prepareVote(t, consensus.MsgPrecommit, currentRound, currentHeight, proposal.ProposalBlock.Hash(), members[sender].Address, privateKeys[members[sender].Address])
	setCommitteeAndSealOnBlock(t, proposal.ProposalBlock, committeeSet, privateKeys, 1)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	backendMock := interfaces.NewMockBackend(ctrl)
	backendMock.EXPECT().Address().Return(clientAddr)
	backendMock.EXPECT().Logger().AnyTimes().Return(log.Root())

	c := New(backendMock, nil)
	c.setHeight(currentHeight)
	c.setRound(currentRound)
	c.SetStep(tctypes.Precommit)
	c.setCommitteeSet(committeeSet)
	c.curRoundMessages.SetProposal(proposal, proposalMsg, true)
	quorumPrecommitMsg := message.Message{Address: members[2].Address, Code: consensus.MsgPrevote, Power: new(big.Int).Sub(c.CommitteeSet().Quorum(), common.Big1)}

	c.curRoundMessages.AddPrecommit(proposal.ProposalBlock.Hash(), quorumPrecommitMsg)

	// The committed seal order is unpredictable, therefore, using gomock.Any()
	// TODO: investigate what order should be on committed seals
	backendMock.EXPECT().Commit(proposal.ProposalBlock, currentRound, gomock.Any())

	err := c.handleValidMsg(context.Background(), precommitMsg)
	assert.NoError(t, err)

	newCommitteeSet, err := tdmcommittee.NewRoundRobinSet(committeeSet.Committee(), members[currentRound].Address)
	c.committee = newCommitteeSet
	assert.NoError(t, err)
	backendMock.EXPECT().HeadBlock().Return(proposal.ProposalBlock, members[currentRound].Address).MaxTimes(2)

	// if the client is the next proposer
	if newCommitteeSet.GetProposer(0).Address == clientAddr {
		t.Log("is proposer")

		c.pendingCandidateBlocks[nextHeight] = nextP.ProposalBlock
		backendMock.EXPECT().SetProposedBlockHash(nextP.ProposalBlock.Hash())

		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().DoAndReturn(signer(privateKeys[clientAddr]))
		backendMock.EXPECT().Broadcast(committeeSet.Committee(), nextProposalMsg.Bytes).Return(nil)
	}

	// It is hard to control tendermint's state machine if we construct the full backend since it overwrites the
	// state we simulated on this test context again and again. So we assume the CommitEvent is sent from miner/worker
	// thread via backend's interface, and it is handled to start new round here:
	c.precommiter.HandleCommit(context.Background())

	assert.Equal(t, big.NewInt(int64(nextHeight)), c.Height())
	assert.Equal(t, int64(0), c.Round())
	assert.Equal(t, tctypes.Propose, c.step)
}

// The following tests aim to test lines 49 - 54 of Tendermint Algorithm described on page 6 of
// https://arxiv.org/pdf/1807.04938.pdf.
func TestFutureRoundChange(t *testing.T) {
	// In the following tests we are assuming that no committee member has voting power more than or equal to F()
	committeeSizeAndMaxRound := maxSize
	committeeSet, privateKeys := prepareCommittee(t, committeeSizeAndMaxRound)
	members := committeeSet.Committee()
	clientAddr := members[0].Address
	roundChangeThreshold := committeeSet.F()
	sender1, sender2 := members[1], members[2]
	sender1.VotingPower = new(big.Int).Sub(roundChangeThreshold, common.Big1)
	sender2.VotingPower = new(big.Int).Sub(roundChangeThreshold, common.Big1)

	t.Run("move to future round after receiving more than F voting power messages", func(t *testing.T) {
		currentHeight := big.NewInt(int64(maxSize))
		// ensure client is not the proposer for next round
		currentRound := int64(50)
		currentStep := tctypes.Step(rand.Intn(3)) //nolint:gosec
		// create random prevote or precommit from 2 different
		msg1, _, _ := prepareVote(t, uint8(rand.Intn(2)+1), currentRound+1, currentHeight, common.Hash{}, sender1.Address, privateKeys[sender1.Address]) //nolint:gosec
		msg2, _, _ := prepareVote(t, uint8(rand.Intn(2)+1), currentRound+1, currentHeight, common.Hash{}, sender2.Address, privateKeys[sender2.Address]) //nolint:gosec
		msg1.Power = new(big.Int).Sub(roundChangeThreshold, common.Big1)
		msg2.Power = new(big.Int).Sub(roundChangeThreshold, common.Big1)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)
		backendMock.EXPECT().Logger().AnyTimes().Return(log.Root())

		c := New(backendMock, nil)
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.SetStep(currentStep)
		c.setCommitteeSet(committeeSet)

		err := c.handleValidMsg(context.Background(), msg1)
		assert.Equal(t, constants.ErrFutureRoundMessage, err)

		err = c.handleValidMsg(context.Background(), msg2)
		assert.Equal(t, constants.ErrFutureRoundMessage, err)

		assert.Equal(t, currentHeight, c.Height())
		assert.Equal(t, currentRound+1, c.Round())
		assert.Equal(t, tctypes.Propose, c.step)
		assert.Equal(t, 2, len(c.backlogs[sender1.Address])+len(c.backlogs[sender2.Address]))
	})

	t.Run("different messages from the same sender cannot cause round change", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		currentStep := tctypes.Step(rand.Intn(3)) //nolint:gosec
		prevoteMsg, _, _ := prepareVote(t, consensus.MsgPrevote, currentRound+1, currentHeight, common.Hash{}, sender1.Address, privateKeys[sender1.Address])
		precommitMsg, _, _ := prepareVote(t, consensus.MsgPrecommit, currentRound+1, currentHeight, common.Hash{}, sender1.Address, privateKeys[sender1.Address])
		// The collective power of the 2 messages  is more than roundChangeThreshold
		prevoteMsg.Power = new(big.Int).Sub(roundChangeThreshold, common.Big1)
		precommitMsg.Power = new(big.Int).Sub(roundChangeThreshold, common.Big1)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)
		backendMock.EXPECT().Logger().AnyTimes().Return(log.Root())

		c := New(backendMock, nil)
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.SetStep(currentStep)
		c.setCommitteeSet(committeeSet)

		err := c.handleValidMsg(context.Background(), prevoteMsg)
		assert.Equal(t, constants.ErrFutureRoundMessage, err)

		err = c.handleValidMsg(context.Background(), precommitMsg)
		assert.Equal(t, constants.ErrFutureRoundMessage, err)

		assert.Equal(t, currentHeight, c.Height())
		assert.Equal(t, currentRound, c.Round())
		assert.Equal(t, currentStep, c.step)
		assert.Equal(t, 2, len(c.backlogs[sender1.Address]))
	})
}

// The following tests are not specific to proposal messages but rather apply to all messages
func TestHandleMessage(t *testing.T) {
	key1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	key2, err := crypto.GenerateKey()
	assert.NoError(t, err)

	key1PubAddr := crypto.PubkeyToAddress(key1.PublicKey)
	key2PubAddr := crypto.PubkeyToAddress(key2.PublicKey)

	committeeSet, err := tdmcommittee.NewRoundRobinSet(types.Committee{types.CommitteeMember{
		Address:     key1PubAddr,
		VotingPower: big.NewInt(1),
	}}, key1PubAddr)
	assert.NoError(t, err)

	t.Run("message sender is not in the committee set", func(t *testing.T) {
		prevHeight := big.NewInt(int64(rand.Intn(100) + 1))
		prevBlock := generateBlock(prevHeight)

		// Prepare message
		msg, _, _ := prepareVote(t, consensus.MsgPrevote, 1, prevHeight, common.Hash{}, key2PubAddr, key2)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(key1PubAddr)
		backendMock.EXPECT().Logger().AnyTimes().Return(log.Root())

		core := New(backendMock, nil)
		core.height = new(big.Int).Add(prevHeight, common.Big1)
		core.setCommitteeSet(committeeSet)
		core.setLastHeader(prevBlock.Header())
		err = core.handleMsg(context.Background(), msg)

		assert.Error(t, err, "unauthorised sender, sender is not is committees set")
	})

	t.Run("message sender is not the message signer", func(t *testing.T) {
		prevHeight := big.NewInt(int64(rand.Intn(100) + 1))
		prevBlock := generateBlock(prevHeight)

		msg, _ := generateBlockProposal(t, 0, new(big.Int).Add(prevHeight, common.Big1), int64(-1), key1PubAddr, false, key2)

		msgRlpNoSig, err := msg.BytesNoSignature()
		assert.NoError(t, err)

		msg.Signature, err = crypto.Sign(crypto.Keccak256(msgRlpNoSig), key2)
		assert.NoError(t, err)

		payload := msg.Bytes

		_, err = message.FromBytes(payload)
		require.NoError(t, err)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(key1PubAddr)
		backendMock.EXPECT().Logger().AnyTimes().Return(log.Root())

		core := New(backendMock, nil)
		core.setCommitteeSet(committeeSet)
		core.setLastHeader(prevBlock.Header())
		core.setHeight(prevBlock.Header().Number)
		err = core.handleMsg(context.Background(), msg)

		assert.Error(t, err, "unauthorised sender, sender is not the signer of the message")
	})

	t.Run("malicious sender sends incorrect signature", func(t *testing.T) {
		prevHeight := big.NewInt(int64(rand.Intn(100) + 1))
		prevBlock := generateBlock(prevHeight)
		sig, err := crypto.Sign(crypto.Keccak256([]byte("random bytes")), key1)
		assert.NoError(t, err)

		msg, _, _ := prepareVote(t, consensus.MsgPrevote, 1, prevHeight, common.Hash{}, key2PubAddr, key2)
		msg.Signature = sig
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(key1PubAddr)
		backendMock.EXPECT().Logger().AnyTimes().Return(log.Root())

		core := New(backendMock, nil)
		core.setCommitteeSet(committeeSet)
		core.setHeight(new(big.Int).Add(prevBlock.Header().Number, common.Big1))
		core.setLastHeader(prevBlock.Header())
		err = core.handleMsg(context.Background(), msg)

		assert.Error(t, err, "malicious sender sends different signature to signature of message")
	})
}

func voteForBadProposal(t *testing.T, step uint8, round int64, height *big.Int, voter common.Address, key *ecdsa.PrivateKey) (*message.Message, []byte, []byte) {
	// prepare the message
	voteRLP, err := rlp.EncodeToBytes(&message.Vote{Round: round, Height: height, ProposedBlockHash: common.Hash{}})
	assert.NoError(t, err)
	voteMsg := &message.Message{Code: step, Payload: voteRLP, Address: voter, Power: common.Big1}
	if step == consensus.MsgPrecommit {
		voteMsg.CommittedSeal, err = sign(helpers.PrepareCommittedSeal(common.Hash{}, round, height), key)
		assert.NoError(t, err)
	}
	voteMsgRLPNoSig, err := voteMsg.BytesNoSignature()
	assert.NoError(t, err)

	voteMsg.Signature, err = sign(voteMsgRLPNoSig, key)
	assert.NoError(t, err)

	voteMsgRLPWithSig := voteMsg.GetBytes()

	return voteMsg, voteMsgRLPNoSig, voteMsgRLPWithSig
}

func prepareVote(t *testing.T, step uint8, round int64, height *big.Int, blockHash common.Hash, clientAddr common.Address, privateKey *ecdsa.PrivateKey) (*message.Message, []byte, []byte) {
	// prepare the message
	vote := &message.Vote{Round: round, Height: height, ProposedBlockHash: blockHash}
	voteRLP, err := rlp.EncodeToBytes(vote)
	assert.NoError(t, err)
	voteMsg := &message.Message{Code: step, Payload: voteRLP, ConsensusMsg: vote, Address: clientAddr, Power: common.Big1}
	if step == consensus.MsgPrecommit {
		voteMsg.CommittedSeal, err = sign(helpers.PrepareCommittedSeal(blockHash, round, height), privateKey)
		assert.NoError(t, err)
	}
	voteMsgRLPNoSig, err := voteMsg.BytesNoSignature()
	assert.NoError(t, err)

	voteMsg.Signature, err = sign(voteMsgRLPNoSig, privateKey)
	assert.NoError(t, err)

	voteMsgRLPWithSig := voteMsg.GetBytes()

	return voteMsg, voteMsgRLPNoSig, voteMsgRLPWithSig
}

func sign(data []byte, key *ecdsa.PrivateKey) ([]byte, error) {
	hashData := crypto.Keccak256(data)
	return crypto.Sign(hashData, key)
}

func generateBlockProposal(t *testing.T, r int64, h *big.Int, vr int64, src common.Address, invalid bool, key *ecdsa.PrivateKey) (*message.Message, *message.Proposal) {
	var block *types.Block
	if invalid {
		header := &types.Header{Number: h}
		header.Difficulty = nil
		block = types.NewBlock(header, nil, nil, nil, new(trie.Trie))
	} else {
		block = generateBlock(h)
	}
	proposal := message.NewProposal(r, h, vr, block, signer(key))
	proposalRlp, err := rlp.EncodeToBytes(proposal)
	assert.NoError(t, err)

	msg := &message.Message{Code: consensus.MsgProposal, ConsensusMsg: proposal, Payload: proposalRlp, Address: src}
	raw, _ := msg.BytesNoSignature()
	msg.Signature, _ = signer(key)(raw)
	msg.Bytes, _ = rlp.EncodeToBytes(msg)
	// we have to do this because encoding and decoding changes some default values and thus same blocks are no longer equal
	return msg, proposal
}

// Committee will be ordered such that the proposer for round(n) == committeeSet.members[n % len(committeeSet.members)]
func prepareCommittee(t *testing.T, cSize int) (interfaces.Committee, helpers.AddressKeyMap) {
	committeeMembers, privateKeys := helpers.GenerateCommittee(cSize)
	committeeSet, err := tdmcommittee.NewRoundRobinSet(committeeMembers, committeeMembers[len(committeeMembers)-1].Address)
	assert.NoError(t, err)
	return committeeSet, privateKeys
}

func generateBlock(height *big.Int) *types.Block {
	// use random nonce to create different blocks
	var nonce types.BlockNonce
	for i := 0; i < len(nonce); i++ {
		nonce[0] = byte(rand.Intn(256))
	}
	header := &types.Header{Number: height, Nonce: nonce}
	block := types.NewBlockWithHeader(header)
	return block
}
