package core

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/autonity/autonity/common"
	tdmcommittee "github.com/autonity/autonity/consensus/tendermint/core/committee"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/trie"
)

const minSize, maxSize = 4, 100
const timeoutDuration, sleepDuration = 1 * time.Microsecond, 1 * time.Millisecond

var testSender = common.HexToAddress("0x8605cdbbdb6d264aa742e77020dcbc58fcdce182")

// The following tests aim to test lines 1 - 21 & 57 - 60 of Tendermint Algorithm described on page 6 of
// https://arxiv.org/pdf/1807.04938.pdf.
func TestStartRoundVariables(t *testing.T) {

	t.Run("ensure round 0 state variables are set correctly", func(t *testing.T) {
		env := NewConsensusEnv(t, nil)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().HeadBlock().Return(env.previousValue)

		env.setupCore(backendMock, env.clientAddress)
		env.core.StartRound(context.Background(), env.curRound)

		// Check the initial consensus state
		env.checkState(t, env.curHeight, env.curRound, Propose, nil, int64(-1), nil, int64(-1))

		env.logging(t)
		// stop the timer to clean up
		err := env.core.proposeTimeout.StopTimer()
		assert.NoError(t, err)
	})
	t.Run("ensure round x state variables are updated correctly", func(t *testing.T) {
		env := NewConsensusEnv(t, nil)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// In this test we are interested in making sure that that values which change in the current round that may
		// have an impact on the actions performed in the following round (in case of round change) are persisted
		// through to the subsequent round.
		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().HeadBlock().Return(env.previousValue).MaxTimes(2)

		env.setupCore(backendMock, env.clientAddress)

		// Check the initial consensus state
		env.core.StartRound(context.Background(), env.curRound)
		env.checkState(t, env.curHeight, env.curRound, Propose, nil, int64(-1), nil, int64(-1))

		// Update locked and valid Value (if locked value changes then valid value also changes, ie quorum(prevotes)
		// delivered in prevote step)
		env.core.SetLockedValue(env.curBlock)
		env.core.SetLockedRound(env.curRound)
		env.core.SetValidValue(env.curBlock)
		env.core.SetValidRound(env.curRound)

		// Move to next round and check the expected state
		env.core.StartRound(context.Background(), env.curRound+1)

		env.checkState(t, env.curHeight, env.curRound+1, Propose, env.curBlock, env.curRound, env.curBlock, env.curRound)

		// Update valid value (we didn't receive quorum prevote in prevote step, also the block changed, ie, locked
		// value and valid value are different)
		currentBlock2 := generateBlock(env.curHeight)
		env.core.SetValidValue(currentBlock2)
		env.core.SetValidRound(env.curRound + 1)

		// Move to next round and check the expected state
		env.core.StartRound(context.Background(), env.curRound+2)

		env.checkState(t, env.curHeight, env.curRound+2, Propose, env.curBlock, env.curRound, currentBlock2, env.curRound+1)

		env.logging(t)
		// stop the timer to clean up
		err := env.core.proposeTimeout.StopTimer()
		assert.NoError(t, err)
	})
}

func TestStartRound(t *testing.T) {
	t.Run("client is the proposer and valid value is nil", func(t *testing.T) {
		e := NewConsensusEnv(t, nil)
		proposal := generateBlockProposal(e.curRound, e.curHeight, -1, false, e.clientSigner)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).DoAndReturn(e.clientSigner)
		backendMock.EXPECT().SetProposedBlockHash(proposal.Block().Hash())
		backendMock.EXPECT().Broadcast(e.committee.Committee(), proposal)
		if e.curRound == 0 {
			// We expect the following extra calls when round = 0
			backendMock.EXPECT().HeadBlock().Return(e.previousValue)
		}
		e.setupCore(backendMock, e.clientAddress)
		e.core.pendingCandidateBlocks[e.curHeight.Uint64()] = proposal.Block()
		e.logging(t)
		e.core.StartRound(context.Background(), e.curRound)
	})
	t.Run("client is the proposer and valid value is not nil", func(t *testing.T) {
		customizer := func(e *ConsensusENV) {
			// Valid round can only be set after round 0, hence the smallest value the round can have is 1 for the valid
			// value to have the smallest value which is 0
			currentRound := int64(rand.Intn(e.committeeSize) + 1)
			for currentRound%int64(e.committeeSize) != 0 {
				currentRound = int64(rand.Intn(e.committeeSize) + 1)
			}
			e.curRound = currentRound
			e.validRound = e.curRound - 1
		}
		env := NewConsensusEnv(t, customizer)
		proposal := generateBlockProposal(env.curRound, env.curHeight, env.validRound, false, env.clientSigner)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().DoAndReturn(env.clientSigner)
		backendMock.EXPECT().SetProposedBlockHash(proposal.Block().Hash())
		backendMock.EXPECT().Broadcast(env.committee.Committee(), proposal)

		env.setupCore(backendMock, env.clientAddress)
		env.core.validValue = proposal.Block()

		env.logging(t)
		env.core.StartRound(context.Background(), env.curRound)
	})
	t.Run("client is not the proposer", func(t *testing.T) {

		e := NewConsensusEnv(t, nil)
		clientIndex := e.committeeSize - 1
		newClientAddr := e.committee.Committee()[clientIndex].Address

		// ensure the client is not the proposer for current round
		currentRound := int64(rand.Intn(e.committeeSize))
		for currentRound%int64(clientIndex) == 0 {
			currentRound = int64(rand.Intn(e.committeeSize))
		}
		e.curRound = currentRound

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().DoAndReturn(e.clientSigner)

		e.core = New(backendMock, nil, newClientAddr, log.Root())
		e.core.setCommitteeSet(e.committee)

		if currentRound == 0 {
			backendMock.EXPECT().HeadBlock().Return(e.previousValue)
		}

		e.logging(t)
		e.core.StartRound(context.Background(), currentRound)
		assert.Equal(t, currentRound, e.core.Round())
		assert.True(t, e.core.proposeTimeout.TimerStarted())

		// stop the timer to clean up
		err := e.core.proposeTimeout.StopTimer()
		assert.NoError(t, err)
	})

	t.Run("at proposal Timeout expiry Timeout event is sent", func(t *testing.T) {
		e := NewConsensusEnv(t, nil)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().DoAndReturn(e.clientSigner)
		backendMock.EXPECT().Post(TimeoutEvent{RoundWhenCalled: e.curRound, HeightWhenCalled: e.curHeight, Step: Propose})

		e.logging(t)
		e.setupCore(backendMock, e.clientAddress)
		assert.False(t, e.core.proposeTimeout.TimerStarted())
		e.core.prevoteTimeout.ScheduleTimeout(timeoutDuration, e.core.Round(), e.core.Height(), e.core.onTimeoutPropose)
		assert.True(t, e.core.prevoteTimeout.TimerStarted())
		time.Sleep(sleepDuration)
	})
	t.Run("at reception of proposal Timeout event prevote nil is sent", func(t *testing.T) {
		customizer := func(e *ConsensusENV) {
			e.step = Propose
		}
		e := NewConsensusEnv(t, customizer)
		timeoutE := TimeoutEvent{RoundWhenCalled: e.curRound, HeightWhenCalled: e.curHeight, Step: Propose}
		prevoteMsg := message.NewPrevote(e.curRound, e.curHeight.Uint64(), common.Hash{}, e.clientSigner)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().DoAndReturn(e.clientSigner)
		backendMock.EXPECT().Broadcast(e.committee.Committee(), prevoteMsg)

		e.logging(t)
		e.setupCore(backendMock, e.clientAddress)
		e.core.handleTimeoutPropose(context.Background(), timeoutE)
		assert.Equal(t, Prevote, e.core.step)
	})
}

// The following tests aim to test lines 22 - 27 of Tendermint Algorithm described on page 6 of
// https://arxiv.org/pdf/1807.04938.pdf.
func TestNewProposal(t *testing.T) {
	committeeSizeAndMaxRound := rand.Intn(maxSize-minSize) + minSize
	committeeSet, privateKeys := prepareCommittee(t, committeeSizeAndMaxRound)
	members := committeeSet.Committee()
	clientAddr := members[0].Address
	clientSigner := makeSigner(privateKeys[clientAddr], clientAddr)
	t.Run("receive invalid proposal for current round", func(t *testing.T) {
		customizer := func(e *ConsensusENV) {
			e.step = Propose
		}
		e := NewConsensusEnv(t, customizer)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		e.setupCore(backendMock, e.clientAddress)
		// Committee()[e.state.curRound].Address means that the sender is the proposer for the current round
		// assume that the message is from a member of committee set and the signature is signing the contents, however,
		// the proposal block inside the message is invalid
		currentProposer := e.committee.Committee()[e.curRound].Address
		currentSigner := makeSigner(e.keys[currentProposer], currentProposer)
		invalidProposal := generateBlockProposal(e.curRound, e.curHeight, -1, true, currentSigner).MustVerify(stubVerifier)
		// prepare prevote nil and target the malicious proposer and the corresponding value.
		prevoteMsg := message.NewPrevote(e.curRound, e.curHeight.Uint64(), common.Hash{}, e.clientSigner)

		backendMock.EXPECT().VerifyProposal(invalidProposal.Block()).Return(time.Duration(1), errors.New("invalid proposal"))
		backendMock.EXPECT().Sign(gomock.Any()).DoAndReturn(e.clientSigner)
		backendMock.EXPECT().Broadcast(e.committee.Committee(), prevoteMsg)

		err := e.core.handleValidMsg(context.Background(), invalidProposal)
		e.logging(t)
		assert.Error(t, err, "expected an error for invalid proposal")
		assert.Equal(t, e.curHeight, e.core.Height())
		assert.Equal(t, e.curRound, e.core.Round())
		assert.Equal(t, Prevote, e.core.step)
	})
	t.Run("receive proposal with validRound = -1 and client's lockedRound = -1", func(t *testing.T) {
		customizer := func(e *ConsensusENV) {
			e.step = Propose
		}
		e := NewConsensusEnv(t, customizer)

		currentProposer := e.committee.Committee()[e.curRound].Address
		currentSigner := makeSigner(e.keys[currentProposer], currentProposer)
		proposal := generateBlockProposal(e.curRound, e.curHeight, -1, false, currentSigner).MustVerify(stubVerifier)
		prevoteMsg := message.NewPrevote(e.curRound, e.curHeight.Uint64(), proposal.Block().Hash(), e.clientSigner)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).DoAndReturn(e.clientSigner)
		backendMock.EXPECT().VerifyProposal(proposal.Block()).Return(time.Duration(1), nil)
		backendMock.EXPECT().Broadcast(e.committee.Committee(), prevoteMsg)

		e.setupCore(backendMock, e.clientAddress)
		e.logging(t)
		err := e.core.handleValidMsg(context.Background(), proposal)
		assert.NoError(t, err)
		assert.Equal(t, e.curHeight, e.core.Height())
		assert.Equal(t, e.curRound, e.core.Round())
		assert.Equal(t, Prevote, e.core.step)
		assert.Nil(t, e.core.lockedValue)
		assert.Equal(t, e.lockedRound, e.core.lockedRound)
	})
	t.Run("receive proposal with validRound = -1 and client's lockedValue is same as proposal block", func(t *testing.T) {
		customizer := func(e *ConsensusENV) {
			e.step = Propose
			e.lockedRound = 0
			e.validRound = 0

			currentProposer := e.committee.Committee()[e.curRound].Address
			currentSigner := makeSigner(e.keys[currentProposer], currentProposer)
			e.curProposal = generateBlockProposal(e.curRound, e.curHeight, -1, false, currentSigner).MustVerify(stubVerifier)
			e.lockedValue = e.curProposal.Block()
			e.validValue = e.curProposal.Block()
		}
		e := NewConsensusEnv(t, customizer)
		prevoteMsg := message.NewPrevote(e.curRound, e.curHeight.Uint64(), e.curProposal.Block().Hash(), e.clientSigner)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).DoAndReturn(e.clientSigner)
		backendMock.EXPECT().VerifyProposal(e.curProposal.Block()).Return(time.Duration(1), nil)
		backendMock.EXPECT().Broadcast(e.committee.Committee(), prevoteMsg)

		e.setupCore(backendMock, e.clientAddress)
		e.logging(t)
		err := e.core.handleValidMsg(context.Background(), e.curProposal)
		assert.NoError(t, err)
		assert.Equal(t, e.curHeight, e.core.Height())
		assert.Equal(t, e.curRound, e.core.Round())
		assert.Equal(t, Prevote, e.core.step)
		assert.Equal(t, e.curProposal.Block(), e.core.lockedValue)
		assert.Equal(t, e.lockedRound, e.core.lockedRound)
		assert.Equal(t, e.curProposal.Block(), e.core.validValue)
		assert.Equal(t, e.validRound, e.core.validRound)
	})
	t.Run("receive proposal with validRound = -1 and client's lockedValue is different from proposal block", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		clientLockedRound := int64(0)
		clientLockedValue := generateBlock(currentHeight)
		currentProposer := members[currentRound].Address
		currentSigner := makeSigner(privateKeys[currentProposer], currentProposer)
		proposal := generateBlockProposal(currentRound, currentHeight, -1, false, currentSigner).MustVerify(stubVerifier)
		prevoteMsg := message.NewPrevote(currentRound, currentHeight.Uint64(), common.Hash{}, clientSigner)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).DoAndReturn(clientSigner)

		t.Log("curRound", currentRound, "curHeight", currentHeight, "committeeSizeAndMaxRound", committeeSizeAndMaxRound)
		c := New(backendMock, nil, clientAddr, log.Root())
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setCommitteeSet(committeeSet)
		c.SetStep(context.Background(), Propose)
		c.lockedRound = clientLockedRound
		c.lockedValue = clientLockedValue
		c.validRound = clientLockedRound
		c.validValue = clientLockedValue

		backendMock.EXPECT().VerifyProposal(proposal.Block()).Return(time.Duration(1), nil)
		backendMock.EXPECT().Broadcast(committeeSet.Committee(), prevoteMsg)

		err := c.handleValidMsg(context.Background(), proposal)
		assert.NoError(t, err)
		assert.Equal(t, currentHeight, c.Height())
		assert.Equal(t, currentRound, c.Round())
		assert.Equal(t, Prevote, c.step)
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
	// extend committeeSet by 1 to prevent the out of range access in members slice with random generated currentRound.
	committeeSet, privateKeys := prepareCommittee(t, committeeSizeAndMaxRound+1)
	members := committeeSet.Committee()
	clientAddr := members[0].Address
	clientSigner := makeSigner(privateKeys[clientAddr], clientAddr)

	t.Run("receive proposal with vr >= 0 and client's lockedRound <= vr", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound-1)) + 1
		// vr >= 0 && vr < round_p
		proposalValidRound := int64(rand.Intn(int(currentRound)))
		// -1 <= c.lockedRound <= vr
		clientLockedRound := int64(rand.Intn(int(proposalValidRound+2) - 1))
		currentProposer := members[currentRound].Address
		currentSigner := makeSigner(privateKeys[currentProposer], currentProposer)
		proposal := generateBlockProposal(currentRound, currentHeight, proposalValidRound, false, currentSigner).MustVerify(stubVerifier)
		prevoteMsg := message.NewPrevote(currentRound, currentHeight.Uint64(), proposal.Block().Hash(), clientSigner)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().DoAndReturn(clientSigner)

		t.Log("curRound", currentRound, "curHeight", currentHeight, "proposalValidRound", proposalValidRound,
			"clientLockedRound", clientLockedRound,
			"committeeSizeAndMaxRound", committeeSizeAndMaxRound)

		c := New(backendMock, nil, clientAddr, log.Root())
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setCommitteeSet(committeeSet)
		c.SetStep(context.Background(), Propose)
		c.lockedRound = clientLockedRound
		c.validRound = clientLockedRound
		c.curRoundMessages = c.messages.GetOrCreate(currentRound)
		// Although the following is not possible it is required to ensure that c.lockRound <= proposalValidRound is
		// responsible for sending the prevote for the incoming proposal
		c.lockedValue = nil
		c.validValue = nil
		fakePrevote := message.Fake{
			FakePower:  c.CommitteeSet().Quorum(),
			FakeValue:  proposal.Block().Hash(),
			FakeSender: testSender,
		}
		c.messages.GetOrCreate(proposalValidRound).AddPrevote(message.NewFakePrevote(fakePrevote))

		backendMock.EXPECT().VerifyProposal(proposal.Block()).Return(time.Duration(1), nil)
		backendMock.EXPECT().Broadcast(committeeSet.Committee(), prevoteMsg)

		err := c.handleValidMsg(context.Background(), proposal)
		assert.NoError(t, err)
		assert.Equal(t, currentHeight, c.Height())
		assert.Equal(t, currentRound, c.Round())
		assert.Equal(t, Prevote, c.step)
		assert.Nil(t, c.lockedValue)
		assert.Equal(t, clientLockedRound, c.lockedRound)
		assert.Nil(t, c.validValue)
		assert.Equal(t, clientLockedRound, c.validRound)
	})
	t.Run("receive proposal with vr >= 0 and client's lockedValue is same as proposal block", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound-1)) + 1
		// vr >= 0 && vr < round_p
		proposalValidRound := int64(rand.Intn(int(currentRound)))

		t.Log("currentHeight", currentHeight, "currentRound", currentRound, "proposalValidRound", proposalValidRound)
		currentProposer := members[currentRound].Address
		currentSigner := makeSigner(privateKeys[currentProposer], currentProposer)
		proposal := generateBlockProposal(currentRound, currentHeight, proposalValidRound, false, currentSigner).MustVerify(stubVerifier)
		prevoteMsg := message.NewPrevote(currentRound, currentHeight.Uint64(), proposal.Block().Hash(), clientSigner)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).DoAndReturn(clientSigner)
		t.Log("curRound", currentRound, "curHeight", currentHeight, "proposalValidRound", proposalValidRound,
			"committeeSizeAndMaxRound", committeeSizeAndMaxRound)
		c := New(backendMock, nil, clientAddr, log.Root())
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setCommitteeSet(committeeSet)
		c.SetStep(context.Background(), Propose)
		// Although the following is not possible it is required to ensure that c.lockedValue = proposal is responsible
		// for sending the prevote for the incoming proposal
		c.lockedRound = proposalValidRound + 1
		c.validRound = proposalValidRound + 1
		c.lockedValue = proposal.Block()
		c.validValue = proposal.Block()
		c.curRoundMessages = c.messages.GetOrCreate(currentRound)
		fakePrevote := message.Fake{
			FakePower:  c.CommitteeSet().Quorum(),
			FakeValue:  proposal.Block().Hash(),
			FakeSender: testSender,
		}
		c.messages.GetOrCreate(proposalValidRound).AddPrevote(message.NewFakePrevote(fakePrevote))

		backendMock.EXPECT().VerifyProposal(proposal.Block()).Return(time.Duration(1), nil)
		backendMock.EXPECT().Broadcast(committeeSet.Committee(), prevoteMsg)

		err := c.handleValidMsg(context.Background(), proposal)
		assert.NoError(t, err)
		assert.Equal(t, currentHeight, c.Height())
		assert.Equal(t, currentRound, c.Round())
		assert.Equal(t, Prevote, c.step)
		assert.Equal(t, proposalValidRound+1, c.lockedRound)
		assert.Equal(t, proposalValidRound+1, c.validRound)
		assert.Equal(t, proposal.Block(), c.lockedValue)
		assert.Equal(t, proposal.Block(), c.validValue)
	})
	t.Run("receive proposal with vr >= 0 and client has lockedRound > vr and lockedValue != proposal", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound-1)) + 1 //+1 to prevent 0 passed to randoms
		clientLockedValue := generateBlock(currentHeight)
		// vr >= 0 && vr < round_p
		proposalValidRound := int64(rand.Intn(int(currentRound)))
		currentProposer := members[currentRound].Address
		currentSigner := makeSigner(privateKeys[currentProposer], currentProposer)
		proposal := generateBlockProposal(currentRound, currentHeight, proposalValidRound, false, currentSigner).MustVerify(stubVerifier)
		prevoteMsg := message.NewPrevote(currentRound, currentHeight.Uint64(), common.Hash{}, clientSigner)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).DoAndReturn(clientSigner)
		t.Log("curRound", currentRound, "curHeight", currentHeight, "proposalValidRound", proposalValidRound,
			"committeeSizeAndMaxRound", committeeSizeAndMaxRound)
		c := New(backendMock, nil, clientAddr, log.Root())
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setCommitteeSet(committeeSet)
		c.SetStep(context.Background(), Propose)
		c.curRoundMessages = c.messages.GetOrCreate(currentRound)

		c.lockedRound = proposalValidRound + 1
		c.validRound = proposalValidRound + 1
		c.lockedValue = clientLockedValue
		c.validValue = clientLockedValue
		fakePrevote := message.NewFakePrevote(message.Fake{FakePower: c.CommitteeSet().Quorum(), FakeValue: proposal.Block().Hash()})
		c.messages.GetOrCreate(proposalValidRound).AddPrevote(fakePrevote)

		backendMock.EXPECT().VerifyProposal(proposal.Block()).Return(time.Duration(0), nil)
		backendMock.EXPECT().Broadcast(committeeSet.Committee(), prevoteMsg)

		err := c.handleValidMsg(context.Background(), proposal)
		assert.NoError(t, err)
		assert.Equal(t, currentHeight, c.Height())
		assert.Equal(t, currentRound, c.Round())
		assert.Equal(t, Prevote, c.step)
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

	/* NOTE: We still need the check for line 28 on receival of an old prevote, HOWEVER the previous analysis is not fully accurate anymore. Indeed when the previous comment was written, the tendermint behaviour was to stop the propose timeout timer once a valid proposal was received. This was **wrong**, the timer should be stopped only when we change height,round or step. Therefore without the line 28 check in prevote.go the algorithm would still be incorrect, but it would not cause a liveness loss (clients would just prevote nil once the timer expires)*/

	t.Run("handle proposal before full quorum prevote on valid round is satisfied, exe action by applying old round prevote into round state", func(t *testing.T) {
		clientIndex := len(members) - 1
		clientAddr = members[clientIndex].Address

		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))

		// ensure the client is not the proposer for current round
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound-1)) + 1
		for currentRound%int64(clientIndex) == 0 {
			currentRound = int64(rand.Intn(committeeSizeAndMaxRound-1)) + 1
		}

		// vr >= 0 && vr < round_p
		proposalValidRound := int64(0) // 0
		if currentRound > 0 {
			proposalValidRound = int64(rand.Intn(int(currentRound)))
		}

		// -1 <= c.lockedRound < vr, if the client lockedValue = vr then the client had received the prevotes in a
		// timely manner thus there are no old prevote yet to arrive
		clientLockedRound := int64(-1) // -1
		if proposalValidRound > 0 {
			clientLockedRound = int64(rand.Intn(int(proposalValidRound)) - 1)
		}

		// the new round proposal
		currentProposer := members[currentRound].Address
		currentSigner := makeSigner(privateKeys[currentProposer], currentProposer)
		proposal := generateBlockProposal(currentRound, currentHeight, proposalValidRound, false, currentSigner).MustVerify(stubVerifier)

		// old proposal some random block
		clientLockedValue := generateBlock(currentHeight)

		// the old round prevote msg to be handled to get the full quorum prevote on old round vr with value v.
		prevoteMsg := message.NewPrevote(proposalValidRound, currentHeight.Uint64(), proposal.Block().Hash(), clientSigner).MustVerify(stubVerifier)

		// the expected prevote msg to be broadcast for the new round with <currentHeight, currentRound, proposal.Block().Hash()>
		prevoteMsgToBroadcast := message.NewPrevote(currentRound, currentHeight.Uint64(), proposal.Block().Hash(), clientSigner)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		c := New(backendMock, nil, clientAddr, log.Root())

		c.setCommitteeSet(committeeSet)
		// construct round state with: old round's quorum-1 prevote for v on valid round.
		fakePrevote := message.Fake{
			FakeRound:  currentRound,
			FakePower:  new(big.Int).Sub(c.CommitteeSet().Quorum(), common.Big1),
			FakeValue:  proposal.Block().Hash(),
			FakeSender: testSender,
		}
		c.messages.GetOrCreate(proposalValidRound).AddPrevote(message.NewFakePrevote(fakePrevote))

		// client on new round's step propose.
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.SetStep(context.Background(), Propose)
		c.lockedRound = clientLockedRound
		c.validRound = clientLockedRound
		c.lockedValue = clientLockedValue
		c.validValue = clientLockedValue

		//schedule the proposer Timeout since the client is not the proposer for this round
		c.proposeTimeout.ScheduleTimeout(1*time.Second, c.Round(), c.Height(), c.onTimeoutPropose)

		backendMock.EXPECT().VerifyProposal(proposal.Block()).Return(time.Duration(1), nil)
		backendMock.EXPECT().Broadcast(committeeSet.Committee(), prevoteMsgToBroadcast)
		backendMock.EXPECT().Sign(gomock.Any()).DoAndReturn(clientSigner)

		t.Log("curRound", currentRound, "curHeight", currentHeight, "proposalValidRound", proposalValidRound,
			"clientLockedRound", clientLockedRound,
			"committeeSizeAndMaxRound", committeeSizeAndMaxRound)

		// now we handle new round's proposal with round_p > vr on value v.
		err := c.handleValidMsg(context.Background(), proposal)
		assert.NoError(t, err)

		// check that the propose timeout is still started, as the proposal did not cause a step change
		assert.True(t, c.proposeTimeout.TimerStarted())

		// now we receive the last old round's prevote MSG to get quorum prevote on vr for value v.
		// the old round's prevote is accepted into the round state which now have the line 28 condition satisfied.
		// now to take the action of line 28 which was not align with pseudo code before.

		err = c.handleValidMsg(context.Background(), prevoteMsg)
		if !errors.Is(err, constants.ErrOldRoundMessage) {
			t.Fatalf("Expected %v, got %v", constants.ErrOldRoundMessage, err)
		}
		assert.True(t, c.sentPrevote)
		assert.Equal(t, currentHeight, c.Height())
		assert.Equal(t, currentRound, c.Round())
		assert.Equal(t, Prevote, c.step)
		assert.Equal(t, clientLockedValue, c.lockedValue)
		assert.Equal(t, clientLockedRound, c.lockedRound)
		assert.Equal(t, clientLockedValue, c.validValue)
		assert.Equal(t, clientLockedRound, c.validRound)
		// now the propose timeout should be stopped, since we moved to prevote step
		assert.False(t, c.proposeTimeout.TimerStarted())
	})
}

func TestProposeTimeout(t *testing.T) {
	committeeSizeAndMaxRound := rand.Intn(maxSize-minSize) + minSize
	committeeSet, privateKeys := prepareCommittee(t, committeeSizeAndMaxRound)
	members := committeeSet.Committee()
	clientAddr := members[0].Address

	t.Run("propose Timeout is not stopped if the proposal does not cause a step change", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))

		// proposal with vr > r
		proposal := generateBlockProposal(currentRound, currentHeight, currentRound+1, false, makeSigner(privateKeys[members[currentRound].Address], members[currentRound].Address)).MustVerify(stubVerifier) //nolint:gosec

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().VerifyProposal(proposal.Block()).Return(time.Duration(1), nil)

		c := New(backendMock, nil, clientAddr, log.Root())
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setCommitteeSet(committeeSet)
		c.SetStep(context.Background(), Propose)

		// propose timer should be started
		c.proposeTimeout.ScheduleTimeout(2*time.Second, c.Round(), c.Height(), c.onTimeoutPropose)
		assert.True(t, c.proposeTimeout.TimerStarted())
		t.Log("curRound", currentRound, "curHeight", currentHeight, "committeeSizeAndMaxRound", committeeSizeAndMaxRound)
		err := c.handleValidMsg(context.Background(), proposal)
		assert.NoError(t, err)
		// propose timer should still be running
		assert.True(t, c.proposeTimeout.TimerStarted())
		assert.Equal(t, currentHeight, c.Height())
		assert.Equal(t, currentRound, c.Round())
		assert.Equal(t, Propose, c.step)
		assert.False(t, c.sentPrevote)
	})
}

// The following tests aim to test lines 34 - 35 & 61 - 64 of Tendermint Algorithm described on page 6 of
// https://arxiv.org/pdf/1807.04938.pdf.
func TestPrevoteTimeout(t *testing.T) {
	committeeSizeAndMaxRound := rand.Intn(maxSize-minSize) + minSize
	committeeSet, privateKeys := prepareCommittee(t, committeeSizeAndMaxRound)
	members := committeeSet.Committee()
	clientAddr := members[0].Address
	clientSigner := makeSigner(privateKeys[clientAddr], clientAddr)
	t.Run("prevote Timeout started after quorum of prevotes with different hashes", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		sender := 1
		currentProposer := members[sender].Address
		currentSigner := makeSigner(privateKeys[currentProposer], currentProposer)
		prevoteMsg := message.NewPrevote(currentRound, currentHeight.Uint64(), generateBlock(currentHeight).Hash(), currentSigner).MustVerify(stubVerifier)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		t.Log("curRound", currentRound, "curHeight", currentHeight, "committeeSizeAndMaxRound", committeeSizeAndMaxRound)

		c := New(backendMock, nil, clientAddr, log.Root())
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setCommitteeSet(committeeSet)
		c.SetStep(context.Background(), Prevote)
		// create quorum prevote messages however there is no quorum on a specific hash
		prevote1 := message.Fake{
			FakeValue:  common.Hash{},
			FakeSender: members[2].Address,
			FakePower:  new(big.Int).Sub(c.CommitteeSet().Quorum(), common.Big2),
		}
		c.curRoundMessages.AddPrevote(message.NewFakePrevote(prevote1))
		prevote2 := message.Fake{
			FakeValue:  generateBlock(currentHeight).Hash(),
			FakeSender: members[3].Address,
			FakePower:  common.Big1,
		}
		c.curRoundMessages.AddPrevote(message.NewFakePrevote(prevote2))

		assert.False(t, c.prevoteTimeout.TimerStarted())
		err := c.handleValidMsg(context.Background(), prevoteMsg)
		assert.NoError(t, err)
		assert.True(t, c.prevoteTimeout.TimerStarted())

		// stop the timer to clean up
		err = c.prevoteTimeout.StopTimer()
		assert.NoError(t, err)

		assert.Equal(t, currentHeight, c.Height())
		assert.Equal(t, currentRound, c.Round())
		assert.Equal(t, Prevote, c.step)
	})
	t.Run("prevote Timeout is not started multiple times", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))

		sender1 := members[1].Address
		sender1Signer := makeSigner(privateKeys[sender1], sender1)
		prevote1Msg := message.NewPrevote(currentRound, currentHeight.Uint64(), generateBlock(currentHeight).Hash(), sender1Signer).MustVerify(stubVerifier)
		sender2 := members[2].Address
		sender2Signer := makeSigner(privateKeys[sender2], sender2)
		prevote2Msg := message.NewPrevote(currentRound, currentHeight.Uint64(), generateBlock(currentHeight).Hash(), sender2Signer).MustVerify(stubVerifier)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		t.Log("curRound", currentRound, "curHeight", currentHeight, "committeeSizeAndMaxRound", committeeSizeAndMaxRound)

		c := New(backendMock, nil, clientAddr, log.Root())
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setCommitteeSet(committeeSet)
		c.SetStep(context.Background(), Prevote)
		// create quorum prevote messages however there is no quorum on a specific hash
		prevote1 := message.Fake{
			FakeValue:  common.Hash{},
			FakeSender: members[3].Address,
			FakePower:  new(big.Int).Sub(c.CommitteeSet().Quorum(), common.Big2),
		}
		c.curRoundMessages.AddPrevote(message.NewFakePrevote(prevote1))
		prevote2 := message.Fake{
			FakeValue:  generateBlock(currentHeight).Hash(),
			FakeSender: members[0].Address,
			FakePower:  common.Big1,
		}
		c.curRoundMessages.AddPrevote(message.NewFakePrevote(prevote2))

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
		assert.Equal(t, Prevote, c.step)
	})
	t.Run("at prevote Timeout expiry Timeout event is sent", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		t.Log("curRound", currentRound, "curHeight", currentHeight, "committeeSizeAndMaxRound", committeeSizeAndMaxRound)

		c := New(backendMock, nil, clientAddr, log.Root())
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setCommitteeSet(committeeSet)
		c.SetStep(context.Background(), Prevote)

		assert.False(t, c.prevoteTimeout.TimerStarted())
		backendMock.EXPECT().Post(TimeoutEvent{RoundWhenCalled: currentRound, HeightWhenCalled: currentHeight, Step: Prevote})
		c.prevoteTimeout.ScheduleTimeout(timeoutDuration, c.Round(), c.Height(), c.onTimeoutPrevote)
		assert.True(t, c.prevoteTimeout.TimerStarted())
		assert.Equal(t, currentHeight, c.Height())
		assert.Equal(t, currentRound, c.Round())
		assert.Equal(t, Prevote, c.step)
		time.Sleep(sleepDuration)
	})
	t.Run("at reception of prevote Timeout event precommit nil is sent", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		timeoutE := TimeoutEvent{RoundWhenCalled: currentRound, HeightWhenCalled: currentHeight, Step: Prevote}
		precommitMsg := message.NewPrecommit(currentRound, currentHeight.Uint64(), common.Hash{}, clientSigner)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		t.Log("curRound", currentRound, "curHeight", currentHeight, "committeeSizeAndMaxRound", committeeSizeAndMaxRound)

		backendMock := interfaces.NewMockBackend(ctrl)
		c := New(backendMock, nil, clientAddr, log.Root())
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setCommitteeSet(committeeSet)
		c.SetStep(context.Background(), Prevote)

		backendMock.EXPECT().Broadcast(committeeSet.Committee(), precommitMsg)
		backendMock.EXPECT().Sign(gomock.Any()).DoAndReturn(clientSigner)

		c.handleTimeoutPrevote(context.Background(), timeoutE)
		assert.Equal(t, currentHeight, c.Height())
		assert.Equal(t, currentRound, c.Round())
		assert.Equal(t, Precommit, c.step)
	})
}

// The following tests aim to test lines 34 - 43 of Tendermint Algorithm described on page 6 of
// https://arxiv.org/pdf/1807.04938.pdf.
func TestQuorumPrevote(t *testing.T) {
	committeeSizeAndMaxRound := rand.Intn(maxSize-minSize) + minSize
	committeeSet, privateKeys := prepareCommittee(t, committeeSizeAndMaxRound)
	members := committeeSet.Committee()
	clientAddr := members[0].Address
	clientSigner := makeSigner(privateKeys[clientAddr], clientAddr)
	signer := func(index int64) message.Signer {
		return makeSigner(privateKeys[members[index].Address], members[index].Address)
	}

	t.Run("receive quroum prevote for proposal block when in step >= prevote", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		//randomly choose prevote or precommit step
		currentStep := Step(rand.Intn(2) + 1)                                                                                              //nolint:gosec
		proposal := generateBlockProposal(currentRound, currentHeight, int64(rand.Intn(int(currentRound+1))), false, signer(currentRound)) //nolint:gosec
		prevoteMsg := message.NewPrevote(currentRound, currentHeight.Uint64(), proposal.Block().Hash(), signer(currentRound)).MustVerify(stubVerifier)
		precommitMsg := message.NewPrecommit(currentRound, currentHeight.Uint64(), proposal.Block().Hash(), clientSigner)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().DoAndReturn(clientSigner)
		t.Log("curStep", currentStep, "curRound", currentRound, "curHeight", currentHeight, "committeeSizeAndMaxRound", committeeSizeAndMaxRound)

		c := New(backendMock, nil, clientAddr, log.Root())
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setCommitteeSet(committeeSet)
		c.SetStep(context.Background(), currentStep)
		c.curRoundMessages.SetProposal(proposal, true)
		fakePrevote := message.Fake{
			FakeValue:  proposal.Block().Hash(),
			FakeRound:  currentRound,
			FakeHeight: currentHeight.Uint64(),
			FakeSender: members[int(currentRound+1)%len(members)].Address,
			FakePower:  new(big.Int).Sub(c.CommitteeSet().Quorum(), common.Big1),
		}
		c.curRoundMessages.AddPrevote(message.NewFakePrevote(fakePrevote))

		if currentStep == Prevote {
			backendMock.EXPECT().Broadcast(committeeSet.Committee(), precommitMsg)
			err := c.handleValidMsg(context.Background(), prevoteMsg)
			assert.NoError(t, err)
			assert.Equal(t, proposal.Block(), c.lockedValue)
			assert.Equal(t, currentRound, c.lockedRound)
			assert.Equal(t, Precommit, c.step)

		} else if currentStep == Precommit {
			err := c.handleValidMsg(context.Background(), prevoteMsg)
			assert.NoError(t, err)
			assert.Equal(t, proposal.Block(), c.validValue)
			assert.Equal(t, currentRound, c.validRound)
		}
		assert.Equal(t, currentHeight, c.Height())
		assert.Equal(t, currentRound, c.Round())

	})

	t.Run("receive more than quorum prevote for proposal block when in step >= prevote", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		//randomly choose prevote or precommit step
		currentStep := Step(rand.Intn(2) + 1) //nolint:gosec
		proposal := generateBlockProposal(currentRound, currentHeight, currentRound-1, false, signer(currentRound)).MustVerify(stubVerifier)

		prevoteMsg1 := message.NewPrevote(currentRound, currentHeight.Uint64(), proposal.Block().Hash(), signer(1)).MustVerify(stubVerifier)
		prevoteMsg2 := message.NewPrevote(currentRound, currentHeight.Uint64(), proposal.Block().Hash(), signer(2)).MustVerify(stubVerifier)
		precommitMsg := message.NewPrecommit(currentRound, currentHeight.Uint64(), proposal.Block().Hash(), clientSigner)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().DoAndReturn(clientSigner)
		t.Log("curStep", currentStep, "curRound", currentRound, "curHeight", currentHeight, "committeeSizeAndMaxRound", committeeSizeAndMaxRound)

		c := New(backendMock, nil, clientAddr, log.Root())
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setCommitteeSet(committeeSet)
		c.SetStep(context.Background(), currentStep)
		c.curRoundMessages.SetProposal(proposal, true)
		fakePrevote := message.Fake{
			FakeValue:  proposal.Block().Hash(),
			FakeSender: members[3].Address,
			FakePower:  new(big.Int).Sub(c.CommitteeSet().Quorum(), common.Big1),
		}
		c.curRoundMessages.AddPrevote(message.NewFakePrevote(fakePrevote))

		// receive first prevote to increase the total to quorum
		if currentStep == Prevote {

			backendMock.EXPECT().Broadcast(committeeSet.Committee(), precommitMsg)

			err := c.handleValidMsg(context.Background(), prevoteMsg1)
			assert.NoError(t, err)

			assert.Equal(t, proposal.Block(), c.lockedValue)
			assert.Equal(t, currentRound, c.lockedRound)
			assert.Equal(t, Precommit, c.step)

		} else if currentStep == Precommit {
			err := c.handleValidMsg(context.Background(), prevoteMsg1)
			assert.NoError(t, err)

			assert.Equal(t, proposal.Block(), c.validValue)
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
	clientSigner := makeSigner(privateKeys[clientAddr], clientAddr)
	currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
	currentRound := int64(rand.Intn(committeeSizeAndMaxRound))

	prevoteMsg := message.NewPrevote(currentRound, currentHeight.Uint64(), common.Hash{}, makeSigner(privateKeys[members[1].Address], members[1].Address)).MustVerify(stubVerifier)
	precommitMsg := message.NewPrecommit(currentRound, currentHeight.Uint64(), common.Hash{}, clientSigner)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	backendMock := interfaces.NewMockBackend(ctrl)
	backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().DoAndReturn(clientSigner)
	t.Log("curRound", currentRound, "curHeight", currentHeight, "committeeSizeAndMaxRound", committeeSizeAndMaxRound)

	c := New(backendMock, nil, clientAddr, log.Root())
	c.setHeight(currentHeight)
	c.setRound(currentRound)
	c.setCommitteeSet(committeeSet)
	c.SetStep(context.Background(), Prevote)
	fakePrevote := message.Fake{
		FakeValue:  common.Hash{},
		FakeSender: members[2].Address,
		FakePower:  new(big.Int).Sub(c.CommitteeSet().Quorum(), common.Big1),
	}
	c.curRoundMessages.AddPrevote(message.NewFakePrevote(fakePrevote))

	backendMock.EXPECT().Broadcast(committeeSet.Committee(), precommitMsg)

	err := c.handleValidMsg(context.Background(), prevoteMsg)
	assert.NoError(t, err)

	assert.Equal(t, currentHeight, c.Height())
	assert.Equal(t, currentRound, c.Round())
	assert.Equal(t, Precommit, c.step)
}

// The following tests aim to test lines 47 - 48 & 65 - 67 of Tendermint Algorithm described on page 6 of
// https://arxiv.org/pdf/1807.04938.pdf.
func TestPrecommitTimeout(t *testing.T) {
	committeeSizeAndMaxRound := rand.Intn(maxSize-minSize) + minSize
	committeeSet, privateKeys := prepareCommittee(t, committeeSizeAndMaxRound)
	members := committeeSet.Committee()
	clientAddr := members[0].Address

	t.Run("at propose step, precommit Timeout started after quorum of precommits with different hashes", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))

		precommit := message.NewPrecommit(
			currentRound,
			currentHeight.Uint64(),
			generateBlock(currentHeight).Hash(),
			makeSigner(privateKeys[members[1].Address], members[1].Address),
		).MustVerify(stubVerifier)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		t.Log("curRound", currentRound, "curHeight", currentHeight, "committeeSizeAndMaxRound", committeeSizeAndMaxRound)

		c := New(backendMock, nil, clientAddr, log.Root())
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.SetStep(context.Background(), Propose)
		c.setCommitteeSet(committeeSet)
		// create quorum precommit messages however there is no quorum on a specific hash
		fakePrecommit1 := message.Fake{
			FakeValue:  common.Hash{},
			FakeSender: members[2].Address,
			FakePower:  new(big.Int).Sub(c.CommitteeSet().Quorum(), common.Big2),
		}
		c.curRoundMessages.AddPrecommit(message.NewFakePrecommit(fakePrecommit1))
		fakePrecommit2 := message.Fake{
			FakeValue:  generateBlock(currentHeight).Hash(),
			FakeSender: members[3].Address,
			FakePower:  common.Big1,
		}
		c.curRoundMessages.AddPrecommit(message.NewFakePrecommit(fakePrecommit2))

		assert.False(t, c.precommitTimeout.TimerStarted())
		err := c.handleValidMsg(context.Background(), precommit)
		assert.NoError(t, err)
		assert.True(t, c.precommitTimeout.TimerStarted())

		// stop the timer to clean up
		err = c.precommitTimeout.StopTimer()
		assert.NoError(t, err)

		assert.Equal(t, currentHeight, c.Height())
		assert.Equal(t, currentRound, c.Round())
		assert.Equal(t, Propose, c.step)
	})

	t.Run("at vote step, precommit Timeout started after quorum of precommits with different hashes", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))

		precommit := message.NewPrecommit(
			currentRound,
			currentHeight.Uint64(),
			generateBlock(currentHeight).Hash(),
			makeSigner(privateKeys[members[1].Address], members[1].Address),
		).MustVerify(stubVerifier)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		t.Log("curRound", currentRound, "curHeight", currentHeight, "committeeSizeAndMaxRound", committeeSizeAndMaxRound)

		c := New(backendMock, nil, clientAddr, log.Root())
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setCommitteeSet(committeeSet)
		c.SetStep(context.Background(), Precommit)
		// create quorum precommit messages however there is no quorum on a specific hash
		fakePrecommit1 := message.Fake{
			FakeValue:  common.Hash{},
			FakeSender: members[2].Address,
			FakePower:  new(big.Int).Sub(c.CommitteeSet().Quorum(), common.Big2),
		}
		c.curRoundMessages.AddPrecommit(message.NewFakePrecommit(fakePrecommit1))
		fakePrecommit2 := message.Fake{
			FakeValue:  generateBlock(currentHeight).Hash(),
			FakeSender: members[3].Address,
			FakePower:  common.Big1,
		}
		c.curRoundMessages.AddPrecommit(message.NewFakePrecommit(fakePrecommit2))

		assert.False(t, c.precommitTimeout.TimerStarted())
		err := c.handleValidMsg(context.Background(), precommit)
		assert.NoError(t, err)
		assert.True(t, c.precommitTimeout.TimerStarted())

		// stop the timer to clean up
		err = c.precommitTimeout.StopTimer()
		assert.NoError(t, err)

		assert.Equal(t, currentHeight, c.Height())
		assert.Equal(t, currentRound, c.Round())
		assert.Equal(t, Precommit, c.step)
	})
	t.Run("precommit Timeout is not started multiple times", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))

		precommitFrom1 := message.NewPrecommit(currentRound,
			currentHeight.Uint64(),
			generateBlock(currentHeight).Hash(),
			makeSigner(privateKeys[members[1].Address], members[1].Address),
		).MustVerify(stubVerifier)

		precommitFrom2 := message.NewPrecommit(
			currentRound,
			currentHeight.Uint64(),
			generateBlock(currentHeight).Hash(),
			makeSigner(privateKeys[members[2].Address], members[2].Address),
		).MustVerify(stubVerifier)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		c := New(backendMock, nil, clientAddr, log.Root())
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setCommitteeSet(committeeSet)
		currentStep := Step(rand.Intn(3))
		c.SetStep(context.Background(), currentStep)
		// create quorum prevote messages however there is no quorum on a specific hash
		fakePrecommit1 := message.Fake{
			FakeValue:  common.Hash{},
			FakeSender: members[3].Address,
			FakePower:  new(big.Int).Sub(c.CommitteeSet().Quorum(), common.Big2),
		}
		c.curRoundMessages.AddPrecommit(message.NewFakePrecommit(fakePrecommit1))
		fakePrecommit2 := message.Fake{
			FakeValue:  generateBlock(currentHeight).Hash(),
			FakeSender: members[0].Address,
			FakePower:  common.Big1,
		}
		c.curRoundMessages.AddPrecommit(message.NewFakePrecommit(fakePrecommit2))

		assert.False(t, c.precommitTimeout.TimerStarted())
		t.Log("curStep", currentStep, "curRound", currentRound, "curHeight", currentHeight, "committeeSizeAndMaxRound", committeeSizeAndMaxRound)

		err := c.handleValidMsg(context.Background(), precommitFrom1)
		assert.NoError(t, err)
		assert.True(t, c.precommitTimeout.TimerStarted())

		timeNow := time.Now()

		err = c.handleValidMsg(context.Background(), precommitFrom2)
		assert.NoError(t, err)
		assert.True(t, c.precommitTimeout.TimerStarted())
		assert.True(t, c.precommitTimeout.Start.Before(timeNow))

		// stop the timer to clean up
		err = c.precommitTimeout.StopTimer()
		assert.NoError(t, err)

		assert.Equal(t, currentHeight, c.Height())
		assert.Equal(t, currentRound, c.Round())
		assert.Equal(t, currentStep, c.step)
	})
	t.Run("at precommit Timeout expiry Timeout event is sent", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)

		c := New(backendMock, nil, clientAddr, log.Root())
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setCommitteeSet(committeeSet)
		currentStep := Step(rand.Intn(3))
		c.SetStep(context.Background(), currentStep)

		assert.False(t, c.precommitTimeout.TimerStarted())
		backendMock.EXPECT().Post(TimeoutEvent{RoundWhenCalled: currentRound, HeightWhenCalled: currentHeight, Step: Precommit})
		t.Log("curStep", currentStep, "curRound", currentRound, "curHeight", currentHeight, "committeeSizeAndMaxRound", committeeSizeAndMaxRound)
		c.precommitTimeout.ScheduleTimeout(timeoutDuration, c.Round(), c.Height(), c.onTimeoutPrecommit)
		assert.True(t, c.precommitTimeout.TimerStarted())
		assert.Equal(t, currentHeight, c.Height())
		assert.Equal(t, currentRound, c.Round())
		assert.Equal(t, currentStep, c.step)
		time.Sleep(sleepDuration)
	})
	t.Run("at reception of precommit Timeout event next round will be started", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		// ensure client is not the proposer for next round
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		for (currentRound+1)%int64(len(members)) == 0 {
			currentRound = int64(rand.Intn(committeeSizeAndMaxRound))
		}
		timeoutE := TimeoutEvent{RoundWhenCalled: currentRound, HeightWhenCalled: currentHeight, Step: Precommit}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)

		c := New(backendMock, nil, clientAddr, log.Root())
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setCommitteeSet(committeeSet)
		currentStep := Step(rand.Intn(3))
		c.SetStep(context.Background(), currentStep)
		t.Log("curStep", currentStep, "curRound", currentRound, "curHeight", currentHeight, "committeeSizeAndMaxRound", committeeSizeAndMaxRound)

		c.handleTimeoutPrecommit(context.Background(), timeoutE)

		assert.Equal(t, currentHeight, c.Height())
		assert.Equal(t, currentRound+1, c.Round())
		assert.Equal(t, Propose, c.step)

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
	nextProposalMsg := generateBlockProposal(0, big.NewInt(int64(nextHeight)), int64(-1), false, makeSigner(privateKeys[clientAddr], clientAddr))

	currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
	proposal := generateBlockProposal(currentRound, currentHeight, currentRound, false, makeSigner(privateKeys[members[currentRound].Address], members[currentRound].Address)) //nolint:gosec
	sender := 1
	precommit := message.NewPrecommit(currentRound, currentHeight.Uint64(), proposal.Block().Hash(), makeSigner(privateKeys[members[sender].Address], members[sender].Address)).MustVerify(stubVerifier)
	setCommitteeAndSealOnBlock(t, proposal.Block(), committeeSet, privateKeys, 1)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	backendMock := interfaces.NewMockBackend(ctrl)
	c := New(backendMock, nil, clientAddr, log.Root())
	c.setHeight(currentHeight)
	c.setRound(currentRound)
	c.setCommitteeSet(committeeSet)
	c.SetStep(context.Background(), Precommit)
	c.curRoundMessages.SetProposal(proposal, true)
	quorumPrecommitMsg := message.Fake{
		FakeValue:  proposal.Block().Hash(),
		FakePower:  new(big.Int).Sub(c.CommitteeSet().Quorum(), common.Big1),
		FakeSender: members[2].Address,
	}
	c.curRoundMessages.AddPrecommit(message.NewFakePrecommit(quorumPrecommitMsg))

	// The committed seal order is unpredictable, therefore, using gomock.Any()
	// TODO: investigate what order should be on committed seals
	backendMock.EXPECT().Commit(proposal.Block(), currentRound, gomock.Any())
	// In case of Timeout propose
	backendMock.EXPECT().Post(gomock.Any()).AnyTimes()
	t.Log("curRound", currentRound, "curHeight", currentHeight, "committeeSizeAndMaxRound", committeeSizeAndMaxRound)

	err := c.handleValidMsg(context.Background(), precommit)
	assert.NoError(t, err)

	newCommitteeSet, err := tdmcommittee.NewRoundRobinSet(committeeSet.Committee(), members[currentRound].Address)
	c.committee = newCommitteeSet
	assert.NoError(t, err)
	backendMock.EXPECT().HeadBlock().Return(proposal.Block()).MaxTimes(2)
	backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().DoAndReturn(makeSigner(privateKeys[clientAddr], clientAddr))
	// if the client is the next proposer
	if newCommitteeSet.GetProposer(0).Address == clientAddr {
		t.Log("is proposer")
		c.pendingCandidateBlocks[nextHeight] = nextProposalMsg.Block()
		backendMock.EXPECT().SetProposedBlockHash(nextProposalMsg.Block().Hash())
		backendMock.EXPECT().Broadcast(committeeSet.Committee(), nextProposalMsg)
	}

	// It is hard to control tendermint's state machine if we construct the full backend since it overwrites the
	// state we simulated on this test context again and again. So we assume the CommitEvent is sent from miner/worker
	// thread via backend's interface, and it is handled to start new round here:
	c.precommiter.HandleCommit(context.Background())

	assert.Equal(t, big.NewInt(int64(nextHeight)), c.Height())
	assert.Equal(t, int64(0), c.Round())
	assert.Equal(t, Propose, c.step)
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
		currentStep := Step(rand.Intn(3)) //nolint:gosec
		// create random prevote or precommit from 2 different
		msg1 := message.NewPrevote(currentRound+1, currentHeight.Uint64(), common.Hash{}, makeSigner(privateKeys[sender1.Address], sender1.Address)).
			MustVerify(func(address common.Address) *types.CommitteeMember {
				return &types.CommitteeMember{
					Address:     address,
					VotingPower: new(big.Int).Sub(roundChangeThreshold, common.Big1),
				}
			})
		msg2 := message.NewPrevote(currentRound+1, currentHeight.Uint64(), common.Hash{}, makeSigner(privateKeys[sender2.Address], sender2.Address)).
			MustVerify(func(address common.Address) *types.CommitteeMember {
				return &types.CommitteeMember{
					Address:     address,
					VotingPower: new(big.Int).Sub(roundChangeThreshold, common.Big1),
				}
			})

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(gomock.Any()).AnyTimes()
		t.Log("curStep", currentStep, "curRound", currentRound, "curHeight", currentHeight, "committeeSizeAndMaxRound", committeeSizeAndMaxRound)

		c := New(backendMock, nil, clientAddr, log.Root())
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setCommitteeSet(committeeSet)
		c.SetStep(context.Background(), currentStep)

		err := c.handleValidMsg(context.Background(), msg1)
		assert.Equal(t, constants.ErrFutureRoundMessage, err)

		err = c.handleValidMsg(context.Background(), msg2)
		assert.Equal(t, constants.ErrFutureRoundMessage, err)

		assert.Equal(t, currentHeight, c.Height())
		assert.Equal(t, currentRound+1, c.Round())
		assert.Equal(t, Propose, c.step)
		assert.Equal(t, 0, len(c.backlogs[sender1.Address])+len(c.backlogs[sender2.Address]))
	})

	t.Run("different messages from the same sender cannot cause round change", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		currentStep := Step(rand.Intn(3)) //nolint:gosec
		// The collective power of the 2 messages  is more than roundChangeThreshold
		prevoteMsg := message.NewPrevote(currentRound+1, currentHeight.Uint64(), common.Hash{}, makeSigner(privateKeys[sender1.Address], sender1.Address)).
			MustVerify(func(address common.Address) *types.CommitteeMember {
				return &types.CommitteeMember{
					Address:     address,
					VotingPower: new(big.Int).Sub(roundChangeThreshold, common.Big1),
				}
			})
		precommitMsg := message.NewPrecommit(currentRound+1, currentHeight.Uint64(), common.Hash{}, makeSigner(privateKeys[sender1.Address], sender1.Address)).
			MustVerify(func(address common.Address) *types.CommitteeMember {
				return &types.CommitteeMember{
					Address:     address,
					VotingPower: new(big.Int).Sub(roundChangeThreshold, common.Big1),
				}
			})

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		t.Log("curStep", currentStep, "curRound", currentRound, "curHeight", currentHeight, "committeeSizeAndMaxRound", committeeSizeAndMaxRound)

		c := New(backendMock, nil, clientAddr, log.Root())
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setCommitteeSet(committeeSet)
		c.SetStep(context.Background(), currentStep)

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
		msg := message.NewPrevote(1, prevHeight.Uint64(), common.Hash{}, makeSigner(key2, key2PubAddr))

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		t.Log("preHeight", prevHeight)

		core := New(backendMock, nil, key1PubAddr, log.Root())
		core.height = new(big.Int).Add(prevHeight, common.Big1)
		core.setCommitteeSet(committeeSet)
		core.setLastHeader(prevBlock.Header())
		err = core.handleMsg(context.Background(), msg)

		assert.Error(t, err, "unauthorised sender, sender is not is committees set")
	})

	t.Run("malicious sender sends incorrect signature", func(t *testing.T) {
		prevHeight := big.NewInt(int64(rand.Intn(100) + 1))
		prevBlock := generateBlock(prevHeight)
		msg := message.NewPrevote(1, prevHeight.Uint64()+1, common.Hash{}, func(_ common.Hash) ([]byte, common.Address) {
			out, err := crypto.Sign(crypto.Keccak256([]byte("random bytes")), testKey)
			assert.NoError(t, err)
			return out, testAddr
		})
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		t.Log("preHeight", prevHeight)

		core := New(backendMock, nil, key1PubAddr, log.Root())
		core.setCommitteeSet(committeeSet)
		core.setHeight(new(big.Int).Add(prevBlock.Header().Number, common.Big1))
		core.setLastHeader(prevBlock.Header())
		err = core.handleMsg(context.Background(), msg)

		assert.Error(t, err, "malicious sender sends different signature to signature of message")
	})
}

func generateBlockProposal(r int64, h *big.Int, vr int64, invalid bool, signer message.Signer) *message.Propose {
	var block *types.Block
	if invalid {
		header := &types.Header{Number: h}
		header.Difficulty = nil
		block = types.NewBlock(header, nil, nil, nil, new(trie.Trie))
	} else {
		block = generateBlock(h)
	}
	return message.NewPropose(r, h.Uint64(), vr, block, signer)
}

// Committee will be ordered such that the proposer for round(n) == committeeSet.members[n % len(committeeSet.members)]
func prepareCommittee(t *testing.T, cSize int) (interfaces.Committee, AddressKeyMap) {
	committeeMembers, privateKeys := GenerateCommittee(cSize)
	committeeSet, err := tdmcommittee.NewRoundRobinSet(committeeMembers, committeeMembers[len(committeeMembers)-1].Address)
	assert.NoError(t, err)
	return committeeSet, privateKeys
}

func generateBlock(height *big.Int) *types.Block {
	// use random nonce to create different blocks
	var nonce types.BlockNonce
	for i := 0; i < len(nonce); i++ {
		nonce[i] = byte(rand.Intn(256))
	}
	header := &types.Header{Number: height, Nonce: nonce}
	block := types.NewBlockWithHeader(header)
	return block
}

func setCommitteeAndSealOnBlock(t *testing.T, b *types.Block, c interfaces.Committee, keys map[common.Address]*ecdsa.PrivateKey, signerIndex int) {
	h := b.Header()
	h.Committee = c.Committee()
	hashData := types.SigHash(h)
	signature, err := crypto.Sign(hashData[:], keys[c.Committee()[signerIndex].Address])
	require.NoError(t, err)
	err = types.WriteSeal(h, signature)
	require.NoError(t, err)
	*b = *b.WithSeal(h)
}

type ConsensusENV struct {
	ConsensusView
	TendermintState

	clientAddress common.Address
	clientKey     *ecdsa.PrivateKey
	clientSigner  message.Signer

	curBlock *types.Block

	curProposal *message.Propose

	core *Core
}

type ConsensusView struct {
	previousHeight *big.Int
	previousValue  *types.Block
	committeeSize  int
	committee      interfaces.Committee
	keys           AddressKeyMap
}

type TendermintState struct {
	curHeight   *big.Int
	curRound    int64
	step        Step
	lockedValue *types.Block
	lockedRound int64
	validValue  *types.Block
	validRound  int64
}

func NewConsensusEnv(t *testing.T, customize func(*ConsensusENV)) *ConsensusENV {
	env := &ConsensusENV{}

	// setup view
	env.previousHeight = big.NewInt(int64(rand.Intn(100) + 1))
	env.previousValue = generateBlock(env.previousHeight)
	env.committeeSize = rand.Intn(maxSize-minSize) + minSize
	env.committee, env.keys = prepareCommittee(t, env.committeeSize)

	// setup initial state
	env.curHeight = big.NewInt(env.previousHeight.Int64() + 1)
	env.curRound = 0
	env.step = PrecommitDone
	env.lockedValue = nil
	env.validValue = nil
	env.lockedRound = -1
	env.validRound = -1

	setCommitteeAndSealOnBlock(t, env.previousValue, env.committee, env.keys, 0)

	env.clientAddress = env.committee.Committee()[0].Address
	env.clientKey = env.keys[env.clientAddress]
	env.clientSigner = makeSigner(env.clientKey, env.clientAddress)

	env.curBlock = generateBlock(env.curHeight)

	if customize != nil {
		customize(env)
	}

	return env
}

func (e *ConsensusENV) logging(t *testing.T) {
	t.Log("curHeight", e.curHeight, "curRound", e.curRound, "committeeSize", e.committeeSize)
}

func (e *ConsensusENV) setupCore(backend interfaces.Backend, address common.Address) {
	e.core = New(backend, nil, address, log.Root())
	e.core.setCommitteeSet(e.committee)
	e.core.setHeight(e.curHeight)
	e.core.setRound(e.curRound)
	e.core.SetValidRound(e.validRound)
	e.core.SetLockedRound(e.lockedRound)
	e.core.SetValidValue(e.validValue)
	e.core.SetLockedValue(e.lockedValue)
	e.core.SetStep(context.Background(), e.step)
}

func (e *ConsensusENV) checkState(t *testing.T, h *big.Int, r int64, s Step, lv *types.Block, lr int64, vv *types.Block, vr int64) { //nolint
	assert.Equal(t, h, e.core.Height())
	assert.Equal(t, r, e.core.Round())
	assert.Equal(t, s, e.core.step)
	assert.Equal(t, lv, e.core.lockedValue)
	assert.Equal(t, lr, e.core.lockedRound)
	assert.Equal(t, vv, e.core.validValue)
	assert.Equal(t, vr, e.core.validRound)
}
