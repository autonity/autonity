package core

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/config"
	tcrypto "github.com/clearmatics/autonity/consensus/tendermint/crypto"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const minSize, maxSize = 4, 100
const timeoutDuration, sleepDuration = 1 * time.Microsecond, 1 * time.Millisecond

func setCommitteeAndSealOnBlock(t *testing.T, b *types.Block, c committee, keys map[common.Address]*ecdsa.PrivateKey, signerIndex int) {
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

	overrideDefaultCoreValues := func(core *core) {
		core.height = big.NewInt(-1)
		core.round = int64(-1)
		core.step = precommitDone
		core.lockedValue = &types.Block{}
		core.lockedRound = 0
		core.validValue = &types.Block{}
		core.validRound = 0
	}

	checkConsensusState := func(t *testing.T, h *big.Int, r int64, s Step, lv *types.Block, lr int64, vv *types.Block, vr int64, core *core) {
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

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddress)
		backendMock.EXPECT().LastCommittedProposal().Return(prevBlock, clientAddress)

		core := New(backendMock, config.RoundRobinConfig())

		overrideDefaultCoreValues(core)
		core.startRound(context.Background(), currentRound)

		// Check the initial consensus state
		checkConsensusState(t, currentHeight, currentRound, propose, nil, int64(-1), nil, int64(-1), core)
	})
	t.Run("ensure round x state variables are updated correctly", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// In this test we are interested in making sure that that values which change in the current round that may
		// have an impact on the actions performed in the following round (in case of round change) are persisted
		// through to the subsequent round.
		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddress)
		backendMock.EXPECT().LastCommittedProposal().Return(prevBlock, clientAddress).MaxTimes(2)

		core := New(backendMock, config.RoundRobinConfig())
		overrideDefaultCoreValues(core)
		core.startRound(context.Background(), currentRound)

		// Check the initial consensus state
		checkConsensusState(t, currentHeight, currentRound, propose, nil, int64(-1), nil, int64(-1), core)

		// Update locked and valid Value (if locked value changes then valid value also changes, ie quorum(prevotes)
		// delivered in prevote step)
		core.lockedValue = currentBlock
		core.lockedRound = currentRound
		core.validValue = currentBlock
		core.validRound = currentRound

		// Move to next round and check the expected state
		core.startRound(context.Background(), currentRound+1)

		checkConsensusState(t, currentHeight, currentRound+1, propose, currentBlock, currentRound, currentBlock, currentRound, core)

		// Update valid value (we didn't receive quorum prevote in prevote step, also the block changed, ie, locked
		// value and valid value are different)
		currentBlock2 := generateBlock(currentHeight)
		core.validValue = currentBlock2
		core.validRound = currentRound + 1

		// Move to next round and check the expected state
		core.startRound(context.Background(), currentRound+2)

		checkConsensusState(t, currentHeight, currentRound+2, propose, currentBlock, currentRound, currentBlock2, currentRound+1, core)
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
		proposalBlock := generateBlock(proposalHeight)
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		for currentRound%int64(len(members)) != 0 {
			currentRound = int64(rand.Intn(committeeSizeAndMaxRound))
		}

		proposalMsg, proposalMsgRLPNoSig, proposalMsgRLPWithSig := prepareProposal(t, currentRound, proposalHeight, int64(-1), proposalBlock, clientAddr, privateKeys[clientAddr])

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)

		core := New(backendMock, config.RoundRobinConfig())
		// We assume that round 0 can only happen when we move to a new height, therefore, height is
		// incremented by 1 in start round when round = 0, and the committee set is updated. However, in test case where
		// round is more than 0, then we need to explicitly update the committee set and height.
		if currentRound > 0 {
			core.committee = committeeSet
			core.height = proposalHeight
		}
		core.pendingUnminedBlocks[proposalHeight.Uint64()] = proposalBlock

		if currentRound == 0 {
			// We expect the following extra calls when round = 0
			backendMock.EXPECT().LastCommittedProposal().Return(prevBlock, lastBlockProposer)
		}
		backendMock.EXPECT().SetProposedBlockHash(proposalBlock.Hash())
		backendMock.EXPECT().Sign(proposalMsgRLPNoSig).Return(proposalMsg.Signature, nil)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet.Committee(), proposalMsgRLPWithSig).Return(nil)

		core.startRound(context.Background(), currentRound)
	})
	t.Run("client is the proposer and valid value is not nil", func(t *testing.T) {
		proposalHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		proposalBlock := generateBlock(proposalHeight)
		// Valid round can only be set after round 0, hence the smallest value the the round can have is 1 for the valid
		// value to have the smallest value which is 0
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound) + 1)
		for currentRound%int64(len(members)) != 0 {
			currentRound = int64(rand.Intn(committeeSizeAndMaxRound) + 1)
		}
		validR := currentRound - 1

		proposalMsg, proposalMsgRLPNoSig, proposalMsgRLPWithSig := prepareProposal(t, currentRound, proposalHeight, validR, proposalBlock, clientAddr, privateKeys[clientAddr])

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)

		core := New(backendMock, config.DefaultConfig())
		core.committee = committeeSet
		core.height = proposalHeight
		core.validRound = validR
		core.validValue = proposalBlock

		backendMock.EXPECT().SetProposedBlockHash(proposalBlock.Hash())
		backendMock.EXPECT().Sign(proposalMsgRLPNoSig).Return(proposalMsg.Signature, nil)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet.Committee(), proposalMsgRLPWithSig).Return(nil)

		core.startRound(context.Background(), currentRound)
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

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)

		core := New(backendMock, config.DefaultConfig())

		if currentRound > 0 {
			core.committee = committeeSet
		}

		if currentRound == 0 {
			backendMock.EXPECT().LastCommittedProposal().Return(prevBlock, clientAddr)
		}

		core.startRound(context.Background(), currentRound)

		assert.Equal(t, currentRound, core.Round())
		assert.True(t, core.proposeTimeout.timerStarted())
	})
	t.Run("at proposal timeout expiry timeout event is sent", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)
		c := New(backendMock, config.DefaultConfig())
		c.setCommitteeSet(committeeSet)
		c.setHeight(currentHeight)
		c.setRound(currentRound)

		assert.False(t, c.proposeTimeout.timerStarted())
		backendMock.EXPECT().Post(TimeoutEvent{currentRound, currentHeight, msgProposal})
		c.prevoteTimeout.scheduleTimeout(timeoutDuration, c.Round(), c.Height(), c.onTimeoutPropose)
		assert.True(t, c.prevoteTimeout.timerStarted())
		time.Sleep(sleepDuration)
	})
	t.Run("at reception of proposal timeout event prevote nil is sent", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		timeoutE := TimeoutEvent{currentRound, currentHeight, msgProposal}
		prevoteMsg, prevoteMsgRLPNoSig, prevoteMsgRLPWithSig := prepareVote(t, msgPrevote, currentRound, currentHeight, common.Hash{}, clientAddr, privateKeys[clientAddr])

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)
		c := New(backendMock, config.DefaultConfig())
		c.setCommitteeSet(committeeSet)
		c.setHeight(currentHeight)
		c.setRound(currentRound)

		backendMock.EXPECT().Sign(prevoteMsgRLPNoSig).Return(prevoteMsg.Signature, nil)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet.Committee(), prevoteMsgRLPWithSig).Return(nil)

		c.handleTimeoutPropose(context.Background(), timeoutE)
		assert.Equal(t, prevote, c.step)
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

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)

		c := New(backendMock, config.DefaultConfig())
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setStep(propose)
		c.setCommitteeSet(committeeSet)

		// members[currentRound] means that the sender is the proposer for the current round
		// assume that the message is from a member of committee set and the signature is signing the contents, however,
		// the proposal block inside the message is invalid
		invalidMsg, invalidProposal := generateBlockProposal(t, currentRound, currentHeight, -1, members[currentRound].Address, true)

		// prepare prevote nil
		prevoteMsg, prevoteMsgRLPNoSig, prevoteMsgRLPWithSig := prepareVote(t, msgPrevote, currentRound, currentHeight, common.Hash{}, clientAddr, privateKeys[clientAddr])

		backendMock.EXPECT().VerifyProposal(*invalidProposal.ProposalBlock).Return(time.Duration(1), errors.New("invalid proposal"))
		backendMock.EXPECT().Sign(prevoteMsgRLPNoSig).Return(prevoteMsg.Signature, nil)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet.Committee(), prevoteMsgRLPWithSig).Return(nil)

		err := c.handleCheckedMsg(context.Background(), invalidMsg, members[currentRound])
		assert.Error(t, err, "expected an error for invalid proposal")
		assert.Equal(t, prevote, c.step)
	})
	t.Run("receive proposal with validRound = -1 and client's lockedRound = -1", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		clientLockedRound := int64(-1)
		proposalMsg, proposal := generateBlockProposal(t, currentRound, currentHeight, -1, members[currentRound].Address, false)
		prevoteMsg, prevoteMsgRLPNoSig, prevoteMsgRLPWithSig := prepareVote(t, msgPrevote, currentRound, currentHeight, proposal.ProposalBlock.Hash(), clientAddr, privateKeys[clientAddr])

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)

		c := New(backendMock, config.DefaultConfig())
		// if lockedRround = - 1 then lockedValue = nil
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setStep(propose)
		c.setCommitteeSet(committeeSet)
		c.lockedRound = clientLockedRound
		c.lockedValue = nil

		backendMock.EXPECT().VerifyProposal(*proposal.ProposalBlock).Return(time.Duration(1), nil)
		backendMock.EXPECT().Sign(prevoteMsgRLPNoSig).Return(prevoteMsg.Signature, nil)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet.Committee(), prevoteMsgRLPWithSig).Return(nil)

		err := c.handleCheckedMsg(context.Background(), proposalMsg, members[currentRound])
		assert.Nil(t, err)
		assert.Equal(t, prevote, c.step)
		assert.Nil(t, c.lockedValue)
		assert.Equal(t, clientLockedRound, c.lockedRound)
	})
	t.Run("receive proposal with validRound = -1 and client's lockedValue is same as proposal block", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		clientLockedRound := int64(0)
		proposalMsg, proposal := generateBlockProposal(t, currentRound, currentHeight, -1, members[currentRound].Address, false)
		prevoteMsg, prevoteMsgRLPNoSig, prevoteMsgRLPWithSig := prepareVote(t, msgPrevote, currentRound, currentHeight, proposal.ProposalBlock.Hash(), clientAddr, privateKeys[clientAddr])

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)

		c := New(backendMock, config.DefaultConfig())
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setStep(propose)
		c.setCommitteeSet(committeeSet)
		c.lockedRound = clientLockedRound
		c.lockedValue = proposal.ProposalBlock
		c.validRound = clientLockedRound
		c.validValue = proposal.ProposalBlock

		backendMock.EXPECT().VerifyProposal(*proposal.ProposalBlock).Return(time.Duration(1), nil)
		backendMock.EXPECT().Sign(prevoteMsgRLPNoSig).Return(prevoteMsg.Signature, nil)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet.Committee(), prevoteMsgRLPWithSig).Return(nil)

		err := c.handleCheckedMsg(context.Background(), proposalMsg, members[currentRound])
		assert.Nil(t, err)
		assert.Equal(t, prevote, c.step)
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
		proposalMsg, proposal := generateBlockProposal(t, currentRound, currentHeight, -1, members[currentRound].Address, false)
		prevoteMsg, prevoteMsgRLPNoSig, prevoteMsgRLPWithSig := prepareVote(t, msgPrevote, currentRound, currentHeight, common.Hash{}, clientAddr, privateKeys[clientAddr])

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)

		c := New(backendMock, config.DefaultConfig())
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setStep(propose)
		c.setCommitteeSet(committeeSet)
		c.lockedRound = clientLockedRound
		c.lockedValue = clientLockedValue
		c.validRound = clientLockedRound
		c.validValue = clientLockedValue

		backendMock.EXPECT().VerifyProposal(*proposal.ProposalBlock).Return(time.Duration(1), nil)
		backendMock.EXPECT().Sign(prevoteMsgRLPNoSig).Return(prevoteMsg.Signature, nil)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet.Committee(), prevoteMsgRLPWithSig).Return(nil)

		err := c.handleCheckedMsg(context.Background(), proposalMsg, members[currentRound])
		assert.Nil(t, err)
		assert.Equal(t, prevote, c.step)
		assert.Equal(t, clientLockedValue, c.lockedValue)
		assert.Equal(t, clientLockedRound, c.lockedRound)
		assert.Equal(t, clientLockedValue, c.validValue)
		assert.Equal(t, clientLockedRound, c.validRound)
	})
}

// The following tests aim to test lines 28 - 33 of Tendermint Algorithm described on page 6 of
// https://arxiv.org/pdf/1807.04938.pdf.
func TestOldProposal(t *testing.T) {
	committeeSizeAndMaxRound := rand.Intn(maxSize-minSize) + minSize
	committeeSet, privateKeys := prepareCommittee(t, committeeSizeAndMaxRound)
	members := committeeSet.Committee()
	clientAddr := members[0].Address

	t.Run("receive proposal with vr >= 0 and client's lockedRound <= vr", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		// vr >= 0 && vr < round_p
		proposalValidRound := int64(rand.Intn(int(currentRound)))
		// -1 <= c.lockedRound <= vr
		clientLockedRound := int64(rand.Intn(int(proposalValidRound+2) - 1))
		proposalMsg, proposal := generateBlockProposal(t, currentRound, currentHeight, proposalValidRound, members[currentRound].Address, false)
		prevoteMsg, prevoteMsgRLPNoSig, prevoteMsgRLPWithSig := prepareVote(t, msgPrevote, currentRound, currentHeight, proposal.ProposalBlock.Hash(), clientAddr, privateKeys[clientAddr])

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)

		c := New(backendMock, config.DefaultConfig())
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setStep(propose)
		c.setCommitteeSet(committeeSet)
		c.lockedRound = clientLockedRound
		c.validRound = clientLockedRound
		// Although the following is not possible it is required to ensure that c.lockRound <= proposalValidRound is
		// responsible for sending the prevote for the incoming proposal
		c.lockedValue = nil
		c.validValue = nil
		c.messages.getOrCreate(proposalValidRound).AddPrevote(proposal.ProposalBlock.Hash(), Message{Code: msgPrevote, power: c.committeeSet().Quorum()})

		backendMock.EXPECT().VerifyProposal(*proposal.ProposalBlock).Return(time.Duration(1), nil)
		backendMock.EXPECT().Sign(prevoteMsgRLPNoSig).Return(prevoteMsg.Signature, nil)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet.Committee(), prevoteMsgRLPWithSig).Return(nil)

		err := c.handleCheckedMsg(context.Background(), proposalMsg, members[currentRound])
		assert.Nil(t, err)
		assert.Equal(t, prevote, c.step)
		assert.Nil(t, c.lockedValue)
		assert.Equal(t, clientLockedRound, c.lockedRound)
		assert.Nil(t, c.validValue)
		assert.Equal(t, clientLockedRound, c.validRound)
	})
	t.Run("receive proposal with vr >= 0 and client's lockedValue is same as proposal block", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		// vr >= 0 && vr < round_p
		proposalValidRound := int64(rand.Intn(int(currentRound)))
		proposalMsg, proposal := generateBlockProposal(t, currentRound, currentHeight, proposalValidRound, members[currentRound].Address, false)
		prevoteMsg, prevoteMsgRLPNoSig, prevoteMsgRLPWithSig := prepareVote(t, msgPrevote, currentRound, currentHeight, proposal.ProposalBlock.Hash(), clientAddr, privateKeys[clientAddr])

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)

		c := New(backendMock, config.DefaultConfig())
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setStep(propose)
		c.setCommitteeSet(committeeSet)
		// Although the following is not possible it is required to ensure that c.lockedValue = proposal is responsible
		// for sending the prevote for the incoming proposal
		c.lockedRound = proposalValidRound + 1
		c.validRound = proposalValidRound + 1
		c.lockedValue = proposal.ProposalBlock
		c.validValue = proposal.ProposalBlock
		c.messages.getOrCreate(proposalValidRound).AddPrevote(proposal.ProposalBlock.Hash(), Message{Code: msgPrevote, power: c.committeeSet().Quorum()})

		backendMock.EXPECT().VerifyProposal(*proposal.ProposalBlock).Return(time.Duration(1), nil)
		backendMock.EXPECT().Sign(prevoteMsgRLPNoSig).Return(prevoteMsg.Signature, nil)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet.Committee(), prevoteMsgRLPWithSig).Return(nil)

		err := c.handleCheckedMsg(context.Background(), proposalMsg, members[currentRound])
		assert.Nil(t, err)
		assert.Equal(t, prevote, c.step)
		assert.Equal(t, proposalValidRound+1, c.lockedRound)
		assert.Equal(t, proposalValidRound+1, c.validRound)
		assert.Equal(t, proposal.ProposalBlock, c.lockedValue)
		assert.Equal(t, proposal.ProposalBlock, c.validValue)
	})
	t.Run("receive proposal with vr >= 0 and clients is lockedRound > vr with a different value", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		clientLockedValue := generateBlock(currentHeight)
		// vr >= 0 && vr < round_p
		proposalValidRound := int64(rand.Intn(int(currentRound)))
		proposalMsg, proposal := generateBlockProposal(t, currentRound, currentHeight, proposalValidRound, members[currentRound].Address, false)
		prevoteMsg, prevoteMsgRLPNoSig, prevoteMsgRLPWithSig := prepareVote(t, msgPrevote, currentRound, currentHeight, common.Hash{}, clientAddr, privateKeys[clientAddr])

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)

		c := New(backendMock, config.DefaultConfig())
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setStep(propose)
		c.setCommitteeSet(committeeSet)
		// Although the following is not possible it is required to ensure that c.lockedValue = proposal is responsible
		// for sending the prevote for the incoming proposal
		c.lockedRound = proposalValidRound + 1
		c.validRound = proposalValidRound + 1
		c.lockedValue = clientLockedValue
		c.validValue = clientLockedValue
		c.messages.getOrCreate(proposalValidRound).AddPrevote(proposal.ProposalBlock.Hash(), Message{Code: msgPrevote, power: c.committeeSet().Quorum()})

		backendMock.EXPECT().VerifyProposal(*proposal.ProposalBlock).Return(time.Duration(1), nil)
		backendMock.EXPECT().Sign(prevoteMsgRLPNoSig).Return(prevoteMsg.Signature, nil)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet.Committee(), prevoteMsgRLPWithSig).Return(nil)

		err := c.handleCheckedMsg(context.Background(), proposalMsg, members[currentRound])
		assert.Nil(t, err)
		assert.Equal(t, prevote, c.step)
		assert.Equal(t, proposalValidRound+1, c.lockedRound)
		assert.Equal(t, proposalValidRound+1, c.validRound)
		assert.Equal(t, clientLockedValue, c.lockedValue)
		assert.Equal(t, clientLockedValue, c.validValue)
	})

	// line 28 check upon condition on prevote handler.
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
		proposalValidRound := int64(rand.Intn(int(currentRound)))

		// -1 <= c.lockedRound < vr, if the client lockedValue = vr then the client had received the prevotes in a
		// timely manner thus there are no old prevote yet to arrive
		clientLockedRound := int64(rand.Intn(int(proposalValidRound)) - 1)

		// the new round proposal
		proposalMsg, proposal := generateBlockProposal(t, currentRound, currentHeight, proposalValidRound, members[currentRound].Address, false)

		// old proposal some random block
		clientLockedValue := generateBlock(currentHeight)

		// the old round prevote msg to be handled to get the full quorum prevote on old round vr with value v.
		prevoteMsg, _, _ := prepareVote(t, msgPrevote, proposalValidRound, currentHeight, proposal.ProposalBlock.Hash(), clientAddr, privateKeys[clientAddr])

		// the expected prevote msg to be broadcast for the new round with <currentHeight, currentRound, proposal.ProposalBlock.Hash()>
		prevoteMsgToBroadcast, prevoteMsgRLPNoSig, prevoteMsgRLPWithSig := prepareVote(t, msgPrevote, currentRound, currentHeight, proposal.ProposalBlock.Hash(), clientAddr, privateKeys[clientAddr])

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)

		c := New(backendMock, config.RoundRobin)
		c.setCommitteeSet(committeeSet)
		// construct round state with: old round's quorum-1 prevote for v on valid round.
		c.messages.getOrCreate(proposalValidRound).AddPrevote(proposal.ProposalBlock.Hash(), Message{Code: msgPrevote, power: c.committeeSet().Quorum() - 1})

		// client on new round's step propose.
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setStep(propose)
		c.lockedRound = clientLockedRound
		c.validRound = clientLockedRound
		c.lockedValue = clientLockedValue
		c.validValue = clientLockedValue

		//schedule the proposer timeout since the client is not the proposer for this round
		c.proposeTimeout.scheduleTimeout(1*time.Second, c.Round(), c.Height(), c.onTimeoutPropose)

		backendMock.EXPECT().VerifyProposal(*proposal.ProposalBlock).Return(time.Duration(1), nil)
		backendMock.EXPECT().Sign(prevoteMsgRLPNoSig).Return(prevoteMsgToBroadcast.Signature, nil)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet.Committee(), prevoteMsgRLPWithSig).Return(nil)

		// now we handle new round's proposal with round_p > vr on value v.
		err := c.handleCheckedMsg(context.Background(), proposalMsg, members[currentRound])
		assert.Nil(t, err)

		// check timer was stopped after receiving the proposal
		assert.False(t, c.proposeTimeout.timerStarted())

		// now we receive the last old round's prevote MSG to get quorum prevote on vr for value v.
		// the old round's prevote is accepted into the round state which now have the line 28 condition satisfied.
		// now to take the action of line 28 which was not align with pseudo code before.
		sender := 0
		err = c.handleCheckedMsg(context.Background(), prevoteMsg, members[sender])
		assert.Nil(t, err)

		assert.Equal(t, prevote, c.step)
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

	t.Run("prevote timeout started after quorum of prevotes with different hashes", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		sender := 1
		prevoteMsg, _, _ := prepareVote(t, msgPrevote, currentRound, currentHeight, generateBlock(currentHeight).Hash(), members[sender].Address, privateKeys[members[sender].Address])

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)

		c := New(backendMock, config.DefaultConfig())
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setStep(prevote)
		c.setCommitteeSet(committeeSet)
		// create quorum prevote messages however there is no quorum on a specific hash
		c.curRoundMessages.AddPrevote(common.Hash{}, Message{Address: members[2].Address, Code: msgPrevote, power: c.committeeSet().Quorum() - 2})
		c.curRoundMessages.AddPrevote(generateBlock(currentHeight).Hash(), Message{Address: members[3].Address, Code: msgPrevote, power: 1})

		assert.False(t, c.prevoteTimeout.timerStarted())
		err := c.handleCheckedMsg(context.Background(), prevoteMsg, members[sender])
		assert.Nil(t, err)
		assert.True(t, c.prevoteTimeout.timerStarted())
	})
	t.Run("prevote timeout is not started multiple times", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		sender1 := 1
		prevote1Msg, _, _ := prepareVote(t, msgPrevote, currentRound, currentHeight, generateBlock(currentHeight).Hash(), members[sender1].Address, privateKeys[members[sender1].Address])
		sender2 := 2
		prevote2Msg, _, _ := prepareVote(t, msgPrevote, currentRound, currentHeight, generateBlock(currentHeight).Hash(), members[sender2].Address, privateKeys[members[sender2].Address])

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)

		c := New(backendMock, config.DefaultConfig())
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setStep(prevote)
		c.setCommitteeSet(committeeSet)
		// create quorum prevote messages however there is no quorum on a specific hash
		c.curRoundMessages.AddPrevote(common.Hash{}, Message{Address: members[3].Address, Code: msgPrevote, power: c.committeeSet().Quorum() - 2})
		c.curRoundMessages.AddPrevote(generateBlock(currentHeight).Hash(), Message{Address: members[0].Address, Code: msgPrevote, power: 1})

		assert.False(t, c.prevoteTimeout.timerStarted())

		err := c.handleCheckedMsg(context.Background(), prevote1Msg, members[sender1])
		assert.Nil(t, err)
		assert.True(t, c.prevoteTimeout.timerStarted())

		timeNow := time.Now()

		err = c.handleCheckedMsg(context.Background(), prevote2Msg, members[sender2])
		assert.Nil(t, err)
		assert.True(t, c.prevoteTimeout.timerStarted())
		assert.True(t, c.prevoteTimeout.start.Before(timeNow))

	})
	t.Run("at prevote timeout expiry timeout event is sent", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)

		c := New(backendMock, config.DefaultConfig())
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setStep(prevote)
		c.setCommitteeSet(committeeSet)

		assert.False(t, c.prevoteTimeout.timerStarted())
		backendMock.EXPECT().Post(TimeoutEvent{currentRound, currentHeight, msgPrevote})
		c.prevoteTimeout.scheduleTimeout(timeoutDuration, c.Round(), c.Height(), c.onTimeoutPrevote)
		assert.True(t, c.prevoteTimeout.timerStarted())
		time.Sleep(sleepDuration)
	})
	t.Run("at reception of prevote timeout event precommit nil is sent", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		timeoutE := TimeoutEvent{currentRound, currentHeight, msgPrevote}
		precommitMsg, precommitMsgRLPNoSig, precommitMsgRLPWithSig := prepareVote(t, msgPrecommit, currentRound, currentHeight, common.Hash{}, clientAddr, privateKeys[clientAddr])
		committedSeal := PrepareCommittedSeal(common.Hash{}, currentRound, currentHeight)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)

		c := New(backendMock, config.DefaultConfig())
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setStep(prevote)
		c.setCommitteeSet(committeeSet)

		backendMock.EXPECT().Sign(committedSeal).Return(precommitMsg.CommittedSeal, nil)
		backendMock.EXPECT().Sign(precommitMsgRLPNoSig).Return(precommitMsg.Signature, nil)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet.Committee(), precommitMsgRLPWithSig).Return(nil)

		c.handleTimeoutPrevote(context.Background(), timeoutE)
		assert.Equal(t, precommit, c.step)
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
		currentStep := Step(rand.Intn(2) + 1)
		proposalMsg, proposal := generateBlockProposal(t, currentRound, currentHeight, int64(rand.Intn(int(currentRound+1)-1)), members[currentRound].Address, false)
		sender := 1
		prevoteMsg, _, _ := prepareVote(t, msgPrevote, currentRound, currentHeight, proposal.ProposalBlock.Hash(), members[sender].Address, privateKeys[members[sender].Address])
		precommitMsg, precommitMsgRLPNoSig, precommitMsgRLPWithSig := prepareVote(t, msgPrecommit, currentRound, currentHeight, proposal.ProposalBlock.Hash(), clientAddr, privateKeys[clientAddr])

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)

		c := New(backendMock, config.DefaultConfig())
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setStep(currentStep)
		c.setCommitteeSet(committeeSet)
		c.curRoundMessages.SetProposal(&proposal, proposalMsg, true)
		c.curRoundMessages.AddPrevote(proposal.ProposalBlock.Hash(), Message{Address: members[2].Address, Code: msgPrevote, power: c.committeeSet().Quorum() - 1})

		if currentStep == prevote {
			committedSeal := PrepareCommittedSeal(proposal.ProposalBlock.Hash(), currentRound, currentHeight)

			backendMock.EXPECT().Sign(committedSeal).Return(precommitMsg.CommittedSeal, nil)
			backendMock.EXPECT().Sign(precommitMsgRLPNoSig).Return(precommitMsg.Signature, nil)
			backendMock.EXPECT().Broadcast(context.Background(), committeeSet.Committee(), precommitMsgRLPWithSig).Return(nil)

			err := c.handleCheckedMsg(context.Background(), prevoteMsg, members[sender])
			assert.Nil(t, err)

			assert.Equal(t, proposal.ProposalBlock, c.lockedValue)
			assert.Equal(t, currentRound, c.lockedRound)
			assert.Equal(t, precommit, c.step)

		} else if currentStep == precommit {
			err := c.handleCheckedMsg(context.Background(), prevoteMsg, members[sender])
			assert.Nil(t, err)

			assert.Equal(t, proposal.ProposalBlock, c.validValue)
			assert.Equal(t, currentRound, c.validRound)
		}
	})

	t.Run("receive more than quorum prevote for proposal block when in step >= prevote", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		//randomly choose prevote or precommit step
		currentStep := Step(rand.Intn(2) + 1)
		proposalMsg, proposal := generateBlockProposal(t, currentRound, currentHeight, int64(rand.Intn(int(currentRound+1)-1)), members[currentRound].Address, false)
		sender1 := 1
		prevoteMsg1, _, _ := prepareVote(t, msgPrevote, currentRound, currentHeight, proposal.ProposalBlock.Hash(), members[sender1].Address, privateKeys[members[sender1].Address])
		sender2 := 2
		prevoteMsg2, _, _ := prepareVote(t, msgPrevote, currentRound, currentHeight, proposal.ProposalBlock.Hash(), members[sender2].Address, privateKeys[members[sender2].Address])
		precommitMsg, precommitMsgRLPNoSig, precommitMsgRLPWithSig := prepareVote(t, msgPrecommit, currentRound, currentHeight, proposal.ProposalBlock.Hash(), clientAddr, privateKeys[clientAddr])

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)

		c := New(backendMock, config.DefaultConfig())
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setStep(currentStep)
		c.setCommitteeSet(committeeSet)
		c.curRoundMessages.SetProposal(&proposal, proposalMsg, true)
		c.curRoundMessages.AddPrevote(proposal.ProposalBlock.Hash(), Message{Address: members[3].Address, Code: msgPrevote, power: c.committeeSet().Quorum() - 1})

		// receive first prevote to increase the total to quorum
		if currentStep == prevote {
			committedSeal := PrepareCommittedSeal(proposal.ProposalBlock.Hash(), currentRound, currentHeight)

			backendMock.EXPECT().Sign(committedSeal).Return(precommitMsg.CommittedSeal, nil)
			backendMock.EXPECT().Sign(precommitMsgRLPNoSig).Return(precommitMsg.Signature, nil)
			backendMock.EXPECT().Broadcast(context.Background(), committeeSet.Committee(), precommitMsgRLPWithSig).Return(nil)

			err := c.handleCheckedMsg(context.Background(), prevoteMsg1, members[sender1])
			assert.Nil(t, err)

			assert.Equal(t, proposal.ProposalBlock, c.lockedValue)
			assert.Equal(t, currentRound, c.lockedRound)
			assert.Equal(t, precommit, c.step)

		} else if currentStep == precommit {
			err := c.handleCheckedMsg(context.Background(), prevoteMsg1, members[sender1])
			assert.Nil(t, err)

			assert.Equal(t, proposal.ProposalBlock, c.validValue)
			assert.Equal(t, currentRound, c.validRound)
		}

		// receive second prevote to increase the total to more than quorum
		lockedValueBefore := c.lockedValue
		validValueBefore := c.validValue
		lockedRoundBefore := c.lockedRound
		validRoundBefore := c.validRound

		err := c.handleCheckedMsg(context.Background(), prevoteMsg2, members[sender2])
		assert.Nil(t, err)

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
	prevoteMsg, _, _ := prepareVote(t, msgPrevote, currentRound, currentHeight, common.Hash{}, members[sender].Address, privateKeys[members[sender].Address])
	precommitMsg, precommitMsgRLPNoSig, precommitMsgRLPWithSig := prepareVote(t, msgPrecommit, currentRound, currentHeight, common.Hash{}, clientAddr, privateKeys[clientAddr])
	committedSeal := PrepareCommittedSeal(common.Hash{}, currentRound, currentHeight)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	backendMock := NewMockBackend(ctrl)
	backendMock.EXPECT().Address().Return(clientAddr)

	c := New(backendMock, config.DefaultConfig())
	c.setHeight(currentHeight)
	c.setRound(currentRound)
	c.setStep(prevote)
	c.setCommitteeSet(committeeSet)
	c.curRoundMessages.AddPrevote(common.Hash{}, Message{Address: members[2].Address, Code: msgPrevote, power: c.committeeSet().Quorum() - 1})

	backendMock.EXPECT().Sign(committedSeal).Return(precommitMsg.CommittedSeal, nil)
	backendMock.EXPECT().Sign(precommitMsgRLPNoSig).Return(precommitMsg.Signature, nil)
	backendMock.EXPECT().Broadcast(context.Background(), committeeSet.Committee(), precommitMsgRLPWithSig).Return(nil)

	err := c.handleCheckedMsg(context.Background(), prevoteMsg, members[sender])
	assert.Nil(t, err)

	assert.Equal(t, precommit, c.step)
}

// The following tests aim to test lines 47 - 48 & 65 - 67 of Tendermint Algorithm described on page 6 of
// https://arxiv.org/pdf/1807.04938.pdf.
func TestPrecommitTimeout(t *testing.T) {
	committeeSizeAndMaxRound := rand.Intn(maxSize-minSize) + minSize
	committeeSet, privateKeys := prepareCommittee(t, committeeSizeAndMaxRound)
	members := committeeSet.Committee()
	clientAddr := members[0].Address

	t.Run("precommit timeout started after quorum of precommits with different hashes", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		sender := 1
		precommitMsg, _, _ := prepareVote(t, msgPrecommit, currentRound, currentHeight, generateBlock(currentHeight).Hash(), members[sender].Address, privateKeys[members[sender].Address])

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)

		c := New(backendMock, config.DefaultConfig())
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		//TODO: this should be changed to Step(rand.Intn(3)) to make sure precommit timeout can be started from any step
		c.setStep(precommit)
		c.setCommitteeSet(committeeSet)
		// create quorum precommit messages however there is no quorum on a specific hash
		c.curRoundMessages.AddPrecommit(common.Hash{}, Message{Address: members[2].Address, Code: msgPrecommit, power: c.committeeSet().Quorum() - 2})
		c.curRoundMessages.AddPrecommit(generateBlock(currentHeight).Hash(), Message{Address: members[3].Address, Code: msgPrecommit, power: 1})

		assert.False(t, c.precommitTimeout.timerStarted())
		err := c.handleCheckedMsg(context.Background(), precommitMsg, members[sender])
		assert.Nil(t, err)
		assert.True(t, c.precommitTimeout.timerStarted())
	})
	t.Run("precommit timeout is not started multiple times", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		sender1 := 1
		precommit1Msg, _, _ := prepareVote(t, msgPrecommit, currentRound, currentHeight, generateBlock(currentHeight).Hash(), members[sender1].Address, privateKeys[members[sender1].Address])
		sender2 := 2
		precommit2Msg, _, _ := prepareVote(t, msgPrecommit, currentRound, currentHeight, generateBlock(currentHeight).Hash(), members[sender2].Address, privateKeys[members[sender2].Address])

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)

		c := New(backendMock, config.DefaultConfig())
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		//TODO: this should be changed to Step(rand.Intn(3)) to make sure precommit timeout can be started from any step
		c.setStep(precommit)
		c.setCommitteeSet(committeeSet)
		// create quorum prevote messages however there is no quorum on a specific hash
		c.curRoundMessages.AddPrecommit(common.Hash{}, Message{Address: members[3].Address, Code: msgPrecommit, power: c.committeeSet().Quorum() - 2})
		c.curRoundMessages.AddPrecommit(generateBlock(currentHeight).Hash(), Message{Address: members[0].Address, Code: msgPrecommit, power: 1})

		assert.False(t, c.precommitTimeout.timerStarted())

		err := c.handleCheckedMsg(context.Background(), precommit1Msg, members[sender1])
		assert.Nil(t, err)
		assert.True(t, c.precommitTimeout.timerStarted())

		timeNow := time.Now()

		err = c.handleCheckedMsg(context.Background(), precommit2Msg, members[sender2])
		assert.Nil(t, err)
		assert.True(t, c.precommitTimeout.timerStarted())
		assert.True(t, c.precommitTimeout.start.Before(timeNow))

	})
	t.Run("at precommit timeout expiry timeout event is sent", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)

		c := New(backendMock, config.DefaultConfig())
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		//TODO: this should be changed to Step(rand.Intn(3)) to make sure precommit timeout can be started from any step
		c.setStep(precommit)
		c.setCommitteeSet(committeeSet)

		assert.False(t, c.precommitTimeout.timerStarted())
		backendMock.EXPECT().Post(TimeoutEvent{currentRound, currentHeight, msgPrecommit})
		c.precommitTimeout.scheduleTimeout(timeoutDuration, c.Round(), c.Height(), c.onTimeoutPrecommit)
		assert.True(t, c.precommitTimeout.timerStarted())
		time.Sleep(sleepDuration)
	})
	t.Run("at reception of precommit timeout event next round will be started", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		// ensure client is not the proposer for next round
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		for (currentRound+1)%int64(len(members)) == 0 {
			currentRound = int64(rand.Intn(committeeSizeAndMaxRound))
		}
		timeoutE := TimeoutEvent{currentRound, currentHeight, msgPrecommit}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)

		c := New(backendMock, config.DefaultConfig())
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		//TODO: this should be changed to Step(rand.Intn(3)) to make sure precommit timeout can be started from any step
		c.setStep(precommit)
		c.setCommitteeSet(committeeSet)

		c.handleTimeoutPrecommit(context.Background(), timeoutE)

		assert.Equal(t, currentRound+1, c.Round())
		assert.Equal(t, propose, c.step)
	})
}

// The following tests aim to test lines 49 - 54 of Tendermint Algorithm described on page 6 of
// https://arxiv.org/pdf/1807.04938.pdf.
func TestQuorumPrecommit(t *testing.T) {
	committeeSizeAndMaxRound := rand.Intn(maxSize-minSize) + minSize
	committeeSet, privateKeys := prepareCommittee(t, committeeSizeAndMaxRound)
	members := committeeSet.Committee()
	clientAddr := members[0].Address
	currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
	nextHeight := currentHeight.Uint64() + 1
	nextProposal := generateBlock(big.NewInt(int64(nextHeight)))
	nextProposalMsg, nextProposalMsgRLPNoSig, nextProposalMsgRLPWithSig := prepareProposal(t, 0, big.NewInt(int64(nextHeight)), int64(-1), nextProposal, clientAddr, privateKeys[clientAddr])
	currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
	proposalMsg, proposal := generateBlockProposal(t, currentRound, currentHeight, int64(rand.Intn(int(currentRound+1)-1)), members[currentRound].Address, false)
	sender := 1
	precommitMsg, _, _ := prepareVote(t, msgPrecommit, currentRound, currentHeight, proposal.ProposalBlock.Hash(), members[sender].Address, privateKeys[members[sender].Address])
	setCommitteeAndSealOnBlock(t, proposal.ProposalBlock, committeeSet, privateKeys, 1)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	backendMock := NewMockBackend(ctrl)
	backendMock.EXPECT().Address().Return(clientAddr)

	c := New(backendMock, config.RoundRobinConfig())
	c.setHeight(currentHeight)
	c.setRound(currentRound)
	c.setStep(precommit)
	c.setCommitteeSet(committeeSet)
	c.curRoundMessages.SetProposal(&proposal, proposalMsg, true)
	quorumPrecommitMsg := Message{Address: members[2].Address, Code: msgPrevote, power: c.committeeSet().Quorum() - 1}
	c.curRoundMessages.AddPrecommit(proposal.ProposalBlock.Hash(), quorumPrecommitMsg)

	// The committed seal order is unpredictable, therefore, using gomock.Any()
	// TODO: investigate what order should be on committed seals
	backendMock.EXPECT().Commit(proposal.ProposalBlock, currentRound, gomock.Any())

	err := c.handleCheckedMsg(context.Background(), precommitMsg, members[sender])
	assert.Nil(t, err)

	newCommitteeSet, err := newRoundRobinSet(committeeSet.Committee(), members[currentRound].Address)
	assert.Nil(t, err)
	backendMock.EXPECT().LastCommittedProposal().Return(proposal.ProposalBlock, members[currentRound].Address).MaxTimes(2)

	// if the client is the next proposer
	if newCommitteeSet.GetProposer(0).Address == clientAddr {
		c.pendingUnminedBlocks[nextHeight] = nextProposal
		backendMock.EXPECT().SetProposedBlockHash(nextProposal.Hash())
		backendMock.EXPECT().Sign(nextProposalMsgRLPNoSig).Return(nextProposalMsg.Signature, nil)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, nextProposalMsgRLPWithSig).Return(nil)
	}

	// It is hard to control tendermint's state machine if we construct the full backend since it overwrites the
	// state we simulated on this test context again and again. So we assume the CommitEvent is sent from miner/worker
	// thread via backend's interface, and it is handled to start new round here:
	c.handleCommit(context.Background())

	assert.Equal(t, big.NewInt(int64(nextHeight)), c.Height())
	assert.Equal(t, int64(0), c.Round())
	assert.Equal(t, propose, c.step)
}

// The following tests aim to test lines 49 - 54 of Tendermint Algorithm described on page 6 of
// https://arxiv.org/pdf/1807.04938.pdf.
func TestFutureRoundChange(t *testing.T) {
	// In the following tests we are assuming that no committee member has voting power more than or equal to F()
	committeeSizeAndMaxRound := rand.Intn(maxSize-minSize) + minSize
	committeeSet, privateKeys := prepareCommittee(t, committeeSizeAndMaxRound)
	members := committeeSet.Committee()
	clientAddr := members[0].Address
	roundChangeThreshold := committeeSet.F()
	sender1, sender2 := members[1], members[2]
	sender1.VotingPower = big.NewInt(int64(roundChangeThreshold - 1))
	sender2.VotingPower = big.NewInt(int64(roundChangeThreshold - 1))

	t.Run("move to future round after receiving more than F voting power messages", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		// ensure client is not the proposer for next round
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		for (currentRound+1)%int64(len(members)) == 0 {
			currentRound = int64(rand.Intn(committeeSizeAndMaxRound))
		}
		currentStep := Step(rand.Intn(3))
		// create random prevote or precommit from 2 different
		msg1, _, _ := prepareVote(t, uint64(rand.Intn(2)+1), currentRound+1, currentHeight, common.Hash{}, sender1.Address, privateKeys[sender1.Address])
		msg2, _, _ := prepareVote(t, uint64(rand.Intn(2)+1), currentRound+1, currentHeight, common.Hash{}, sender2.Address, privateKeys[sender2.Address])
		msg1.power = roundChangeThreshold - 1
		msg2.power = roundChangeThreshold - 1

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)

		c := New(backendMock, config.DefaultConfig())
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setStep(currentStep)
		c.setCommitteeSet(committeeSet)

		err := c.handleCheckedMsg(context.Background(), msg1, sender1)
		assert.Equal(t, errFutureRoundMessage, err)

		err = c.handleCheckedMsg(context.Background(), msg2, sender2)
		assert.Equal(t, errFutureRoundMessage, err)

		assert.Equal(t, currentHeight, c.Height())
		assert.Equal(t, currentRound+1, c.Round())
		assert.Equal(t, propose, c.step)
		assert.Equal(t, 2, c.backlogs[sender1].Size()+c.backlogs[sender2].Size())
	})

	t.Run("different messages from the same sender cannot cause round change", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		currentStep := Step(rand.Intn(3))
		prevoteMsg, _, _ := prepareVote(t, msgPrevote, currentRound+1, currentHeight, common.Hash{}, sender1.Address, privateKeys[sender1.Address])
		precommitMsg, _, _ := prepareVote(t, msgPrecommit, currentRound+1, currentHeight, common.Hash{}, sender1.Address, privateKeys[sender1.Address])
		// The collective power of the 2 messages  is more than roundChangeThreshold
		prevoteMsg.power = roundChangeThreshold - 1
		precommitMsg.power = roundChangeThreshold - 1

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddr)

		c := New(backendMock, config.DefaultConfig())
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setStep(currentStep)
		c.setCommitteeSet(committeeSet)

		err := c.handleCheckedMsg(context.Background(), prevoteMsg, sender1)
		assert.Equal(t, errFutureRoundMessage, err)

		err = c.handleCheckedMsg(context.Background(), precommitMsg, sender1)
		assert.Equal(t, errFutureRoundMessage, err)

		assert.Equal(t, currentHeight, c.Height())
		assert.Equal(t, currentRound, c.Round())
		assert.Equal(t, currentStep, c.step)
		assert.Equal(t, 2, c.backlogs[sender1].Size())
	})
}

// The following tests are not specific to proposal messages but rather apply to all messages
func TestHandleMessage(t *testing.T) {
	key1, err := crypto.GenerateKey()
	assert.Nil(t, err)
	key2, err := crypto.GenerateKey()
	assert.Nil(t, err)

	key1PubAddr := crypto.PubkeyToAddress(key1.PublicKey)
	key2PubAddr := crypto.PubkeyToAddress(key2.PublicKey)

	committeeSet, err := newRoundRobinSet(types.Committee{types.CommitteeMember{
		Address:     key1PubAddr,
		VotingPower: big.NewInt(1),
	}}, key1PubAddr)
	assert.Nil(t, err)

	t.Run("message sender is not in the committee set", func(t *testing.T) {
		prevHeight := big.NewInt(int64(rand.Intn(100) + 1))
		prevBlock := generateBlock(prevHeight)

		// Prepare message
		msg := &Message{Address: key2PubAddr, Code: uint64(rand.Intn(3)), Msg: []byte("random message1")}

		msgRlpNoSig, err := msg.PayloadNoSig()
		assert.Nil(t, err)

		msg.Signature, err = crypto.Sign(crypto.Keccak256(msgRlpNoSig), key2)
		assert.Nil(t, err)

		msgRlpWithSig, err := msg.Payload()
		assert.Nil(t, err)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(key1PubAddr)

		core := New(backendMock, config.DefaultConfig())
		core.setCommitteeSet(committeeSet)
		core.lastHeader = prevBlock.Header()
		err = core.handleMsg(context.Background(), msgRlpWithSig)

		assert.Error(t, err, "unauthorised sender, sender is not is committees set")
	})

	t.Run("message sender is not the message signer", func(t *testing.T) {
		prevHeight := big.NewInt(int64(rand.Intn(100) + 1))
		prevBlock := generateBlock(prevHeight)
		msg := &Message{Address: key1PubAddr, Code: uint64(rand.Intn(3)), Msg: []byte("random message2")}

		msgRlpNoSig, err := msg.PayloadNoSig()
		assert.Nil(t, err)

		msg.Signature, err = crypto.Sign(crypto.Keccak256(msgRlpNoSig), key1)
		assert.Nil(t, err)

		msgRlpWithSig, err := msg.Payload()
		assert.Nil(t, err)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(key1PubAddr)

		core := New(backendMock, config.DefaultConfig())
		core.setCommitteeSet(committeeSet)
		core.lastHeader = prevBlock.Header()
		err = core.handleMsg(context.Background(), msgRlpWithSig)

		assert.Error(t, err, "unauthorised sender, sender is not the signer of the message")
	})

	t.Run("malicious sender sends incorrect signature", func(t *testing.T) {
		prevHeight := big.NewInt(int64(rand.Intn(100) + 1))
		prevBlock := generateBlock(prevHeight)
		sig, err := crypto.Sign(crypto.Keccak256([]byte("random bytes")), key1)
		assert.Nil(t, err)

		msg := &Message{Address: key1PubAddr, Code: uint64(rand.Intn(3)), Msg: []byte("random message2"), Signature: sig}
		msgRlpWithSig, err := msg.Payload()
		assert.Nil(t, err)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(key1PubAddr)

		core := New(backendMock, config.DefaultConfig())
		core.setCommitteeSet(committeeSet)
		core.lastHeader = prevBlock.Header()
		err = core.handleMsg(context.Background(), msgRlpWithSig)

		assert.Error(t, err, "malicious sender sends different signature to signature of message")
	})
}

func prepareProposal(t *testing.T, currentRound int64, proposalHeight *big.Int, validR int64, proposalBlock *types.Block, clientAddress common.Address, privateKey *ecdsa.PrivateKey) (*Message, []byte, []byte) {
	// prepare the proposal message
	proposalRLP, err := Encode(NewProposal(currentRound, proposalHeight, validR, proposalBlock))
	assert.Nil(t, err)

	proposalMsg := &Message{Code: msgProposal, Msg: proposalRLP, Address: clientAddress, power: 1}
	proposalMsgRLPNoSig, err := proposalMsg.PayloadNoSig()
	assert.Nil(t, err)

	proposalMsg.Signature, err = sign(proposalMsgRLPNoSig, privateKey)
	assert.Nil(t, err)

	proposalMsgRLPWithSig, err := proposalMsg.Payload()
	assert.Nil(t, err)

	return proposalMsg, proposalMsgRLPNoSig, proposalMsgRLPWithSig
}

func prepareVote(t *testing.T, step uint64, round int64, height *big.Int, blockHash common.Hash, clientAddr common.Address, privateKey *ecdsa.PrivateKey) (*Message, []byte, []byte) {
	// prepare the proposal message
	voteRLP, err := Encode(&Vote{Round: round, Height: height, ProposedBlockHash: blockHash})
	assert.Nil(t, err)
	voteMsg := &Message{Code: step, Msg: voteRLP, Address: clientAddr, power: 1}
	if step == msgPrecommit {
		voteMsg.CommittedSeal, err = sign(PrepareCommittedSeal(blockHash, round, height), privateKey)
		assert.Nil(t, err)
	}
	voteMsgRLPNoSig, err := voteMsg.PayloadNoSig()
	assert.Nil(t, err)

	voteMsg.Signature, err = sign(voteMsgRLPNoSig, privateKey)
	assert.Nil(t, err)

	voteMsgRLPWithSig, err := voteMsg.Payload()
	assert.Nil(t, err)

	return voteMsg, voteMsgRLPNoSig, voteMsgRLPWithSig
}

func sign(data []byte, key *ecdsa.PrivateKey) ([]byte, error) {
	hashData := crypto.Keccak256(data)
	return crypto.Sign(hashData, key)
}

func generateBlockProposal(t *testing.T, r int64, h *big.Int, vr int64, src common.Address, invalid bool) (*Message, Proposal) {
	var block *types.Block
	if invalid {
		header := &types.Header{Number: h}
		header.Difficulty = nil
		block = types.NewBlock(header, nil, nil, nil)
	} else {
		block = generateBlock(h)
	}
	proposal := NewProposal(r, h, vr, block)
	proposalRlp, err := Encode(proposal)
	assert.Nil(t, err)

	msg := Message{Code: msgProposal, Msg: proposalRlp, Address: src}

	var p Proposal
	// we have to do this because encoding and decoding changes some default values and thus same blocks are no longer equal
	err = msg.Decode(&p)
	assert.Nil(t, err)

	return &msg, p
}

// Committee will be ordered such that the proposer for round(n) == committeeSet.members[n % len(committeeSet.members)]
func prepareCommittee(t *testing.T, cSize int) (committee, addressKeyMap) {
	committeeMembers, privateKeys := generateCommittee(cSize)
	committeeSet, err := newRoundRobinSet(committeeMembers, committeeMembers[len(committeeMembers)-1].Address)
	assert.Nil(t, err)
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
