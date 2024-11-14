package core

import (
	"context"
	"errors"
	"math/big"
	"math/rand"
	"sync"
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
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/log"
)

const timeoutDuration, sleepDuration = 1 * time.Microsecond, 1 * time.Millisecond

var testSender = common.HexToAddress("0x8605cdbbdb6d264aa742e77020dcbc58fcdce182")

var signer = func(e *ConsensusENV, index int64) message.Signer {
	return makeSigner(e.keys[e.committee.Committee().Members[index].Address].consensus)
}
var member = func(e *ConsensusENV, index int64) *types.CommitteeMember {
	return &e.committee.Committee().Members[index]
}

// The following tests aim to test lines 1 - 21 & 57 - 60 of Tendermint Algorithm described on page 6 of
// https://arxiv.org/pdf/1807.04938.pdf.
func TestStartRoundVariables(t *testing.T) {
	t.Run("ensure round 0 state variables are set correctly", func(t *testing.T) {
		env := NewConsensusEnv(t, nil)
		ctrl := gomock.NewController(t)
		defer waitForExpects(ctrl)

		backendMock := interfaces.NewMockBackend(ctrl)
		env.setupCore(backendMock, env.clientAddress)
		backendMock.EXPECT().EpochByHeight(env.core.Height().Uint64()).Return(env.LatestEpoch(), nil)
		backendMock.EXPECT().HeadBlock().Return(env.previousValue)
		backendMock.EXPECT().Post(gomock.Any()).Times(1)
		backendMock.EXPECT().ProcessFutureMsgs(env.previousHeight.Uint64() + 1).Times(1)

		env.core.StartRound(context.Background(), env.curRound)

		// Check the initial consensus state
		env.checkState(t, env.curHeight, env.curRound, Propose, nil, int64(-1), nil, int64(-1))

		// stop the timer to clean up
		err := env.core.proposeTimeout.StopTimer()
		assert.NoError(t, err)
	})
	t.Run("ensure round x state variables are updated correctly", func(t *testing.T) {
		env := NewConsensusEnv(t, nil)
		ctrl := gomock.NewController(t)
		defer waitForExpects(ctrl)

		// In this test we are interested in making sure that that values which change in the current round that may
		// have an impact on the actions performed in the following round (in case of round change) are persisted
		// through to the subsequent round.
		backendMock := interfaces.NewMockBackend(ctrl)
		env.setupCore(backendMock, env.clientAddress)
		backendMock.EXPECT().EpochByHeight(env.core.Height().Uint64()).Return(env.LatestEpoch(), nil)
		backendMock.EXPECT().HeadBlock().Return(env.previousValue).MaxTimes(2)
		backendMock.EXPECT().Post(gomock.Any()).Times(3)
		backendMock.EXPECT().ProcessFutureMsgs(env.previousHeight.Uint64() + 1).Times(1)

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
		currentBlock2 := generateBlock(env.curHeight, env.previousValue.Header())
		env.core.SetValidValue(currentBlock2)
		env.core.SetValidRound(env.curRound + 1)

		// Move to next round and check the expected state
		env.core.StartRound(context.Background(), env.curRound+2)

		env.checkState(t, env.curHeight, env.curRound+2, Propose, env.curBlock, env.curRound, currentBlock2, env.curRound+1)

		// stop the timer to clean up
		err := env.core.proposeTimeout.StopTimer()
		assert.NoError(t, err)
	})
}

func TestStartRound(t *testing.T) {

	t.Run("client is the proposer and valid value is nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer waitForExpects(ctrl)

		e := NewConsensusEnv(t, nil)
		proposal := generateBlockProposal(e.curRound, e.curHeight, -1, false, e.clientSigner, e.clientMember, e.previousValue.Header())

		backendMock := interfaces.NewMockBackend(ctrl)
		e.setupCore(backendMock, e.clientAddress)
		backendMock.EXPECT().EpochByHeight(e.core.Height().Uint64()).Return(e.LatestEpoch(), nil)
		backendMock.EXPECT().Sign(gomock.Any()).DoAndReturn(e.clientSigner)
		backendMock.EXPECT().SetProposedBlockHash(proposal.Block().Hash())
		backendMock.EXPECT().Broadcast(e.committee.Committee(), proposal)
		backendMock.EXPECT().HeadBlock().Return(e.previousValue).Times(2)
		backendMock.EXPECT().Post(gomock.Any()).Times(1)
		backendMock.EXPECT().ProcessFutureMsgs(e.previousHeight.Uint64() + 1).Times(1)
		e.core.pendingCandidateBlocks[e.curHeight.Uint64()] = proposal.Block()

		e.core.StartRound(context.Background(), e.curRound)
		e.checkState(t, e.curHeight, e.curRound, Propose, nil, int64(-1), nil, int64(-1))
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
		e := NewConsensusEnv(t, customizer)
		proposal := generateBlockProposal(e.curRound, e.curHeight, e.validRound, false, e.clientSigner, e.clientMember, e.previousValue.Header())

		ctrl := gomock.NewController(t)
		defer waitForExpects(ctrl)

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().DoAndReturn(e.clientSigner)
		backendMock.EXPECT().SetProposedBlockHash(proposal.Block().Hash())
		backendMock.EXPECT().Broadcast(e.committee.Committee(), proposal)
		backendMock.EXPECT().Post(gomock.Any()).Times(1)
		backendMock.EXPECT().HeadBlock().Return(e.previousValue)

		e.setupCore(backendMock, e.clientAddress)
		e.core.validValue = proposal.Block()
		e.core.StartRound(context.Background(), e.curRound)
		e.checkState(t, e.curHeight, e.curRound, Propose, nil, int64(-1), proposal.Block(), e.validRound)
	})
	t.Run("client is not the proposer", func(t *testing.T) {
		customizer := func(e *ConsensusENV) {
			clientIndex := e.committeeSize - 1
			e.clientAddress = e.committee.Committee().Members[clientIndex].Address
			e.clientKey = e.keys[e.clientAddress].consensus
			e.clientSigner = signer(e, int64(clientIndex))
			e.clientMember = member(e, int64(clientIndex))
			// ensure the client is not the proposer for current round
			currentRound := int64(rand.Intn(e.committeeSize))
			for currentRound%int64(clientIndex) == 0 {
				currentRound = int64(rand.Intn(e.committeeSize))
			}
			e.curRound = currentRound
		}
		e := NewConsensusEnv(t, customizer)

		ctrl := gomock.NewController(t)
		defer waitForExpects(ctrl)

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().DoAndReturn(e.clientSigner)
		backendMock.EXPECT().Post(gomock.Any()).Times(1)

		e.setupCore(backendMock, e.clientAddress)
		e.core.setCommitteeSet(e.committee)

		if e.curRound == 0 {
			backendMock.EXPECT().HeadBlock().Return(e.previousValue)
			backendMock.EXPECT().ProcessFutureMsgs(e.previousHeight.Uint64() + 1).Times(1)
		}

		e.core.StartRound(context.Background(), e.curRound)
		assert.Equal(t, e.curRound, e.core.Round())
		assert.True(t, e.core.proposeTimeout.TimerStarted())

		// stop the timer to clean up
		err := e.core.proposeTimeout.StopTimer()
		assert.NoError(t, err)
		e.checkState(t, e.curHeight, e.curRound, Propose, nil, int64(-1), nil, int64(-1))
	})

	t.Run("at proposal Timeout expiry Timeout event is sent", func(t *testing.T) {
		e := NewConsensusEnv(t, nil)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().DoAndReturn(e.clientSigner)
		backendMock.EXPECT().Post(TimeoutEvent{RoundWhenCalled: e.curRound, HeightWhenCalled: e.curHeight, Step: Propose})
		e.setupCore(backendMock, e.clientAddress)
		assert.False(t, e.core.proposeTimeout.TimerStarted())
		e.core.prevoteTimeout.ScheduleTimeout(timeoutDuration, e.core.Round(), e.core.Height(), e.core.onTimeoutPropose)
		assert.True(t, e.core.prevoteTimeout.TimerStarted())
		time.Sleep(sleepDuration)
		e.checkState(t, e.curHeight, e.curRound, PrecommitDone, nil, int64(-1), nil, int64(-1))
	})
	t.Run("at reception of proposal Timeout event prevote nil is sent", func(t *testing.T) {
		customizer := func(e *ConsensusENV) {
			e.step = Propose
		}
		e := NewConsensusEnv(t, customizer)
		timeoutE := TimeoutEvent{RoundWhenCalled: e.curRound, HeightWhenCalled: e.curHeight, Step: Propose}
		prevoteMsg := message.NewPrevote(e.curRound, e.curHeight.Uint64(), common.Hash{}, e.clientSigner, e.clientMember, e.committeeSize)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().DoAndReturn(e.clientSigner)
		backendMock.EXPECT().Broadcast(e.committee.Committee(), prevoteMsg)

		e.setupCore(backendMock, e.clientAddress)
		e.core.handleTimeoutPropose(context.Background(), timeoutE)
		e.checkState(t, e.curHeight, e.curRound, Prevote, nil, int64(-1), nil, int64(-1))
	})
}

// The following tests aim to test lines 22 - 27 of Tendermint Algorithm described on page 6 of
// https://arxiv.org/pdf/1807.04938.pdf.
func TestNewProposal(t *testing.T) {

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
		invalidProposal := generateBlockProposal(e.curRound, e.curHeight, -1, true, signer(e, e.curRound), member(e, e.curRound), e.previousValue.Header())
		// prepare prevote nil and target the malicious proposer and the corresponding value.
		prevoteMsg := message.NewPrevote(e.curRound, e.curHeight.Uint64(), common.Hash{}, e.clientSigner, e.clientMember, e.committeeSize)

		backendMock.EXPECT().ProposedBlockHash().Return(common.Hash{})
		backendMock.EXPECT().IsProposalStateCached(invalidProposal.Block().Hash()).Return(false)
		backendMock.EXPECT().VerifyProposal(invalidProposal.Block()).Return(time.Duration(1), errors.New("invalid proposal"))
		backendMock.EXPECT().Sign(gomock.Any()).DoAndReturn(e.clientSigner)
		backendMock.EXPECT().Broadcast(e.committee.Committee(), prevoteMsg)

		err := e.core.handleMsg(context.Background(), invalidProposal)
		assert.Error(t, err, "expected an error for invalid proposal")
		e.checkState(t, e.curHeight, e.curRound, Prevote, nil, int64(-1), nil, int64(-1))
	})
	t.Run("receive proposal with validRound = -1 and client's lockedRound = -1", func(t *testing.T) {

		customizer := func(e *ConsensusENV) {
			e.step = Propose
		}
		e := NewConsensusEnv(t, customizer)

		proposal := generateBlockProposal(e.curRound, e.curHeight, -1, false, signer(e, e.curRound), member(e, e.curRound), e.previousValue.Header())
		prevoteMsg := message.NewPrevote(e.curRound, e.curHeight.Uint64(), proposal.Block().Hash(), e.clientSigner, e.clientMember, e.committeeSize)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		wg := sync.WaitGroup{}
		wg.Add(1)
		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).DoAndReturn(e.clientSigner)
		backendMock.EXPECT().ProposedBlockHash().Return(common.Hash{})
		backendMock.EXPECT().IsProposalStateCached(proposal.Block().Hash()).Return(false)
		backendMock.EXPECT().ProposalVerified(proposal.Block()).Do(func(i any) { wg.Done() })
		backendMock.EXPECT().VerifyProposal(proposal.Block()).Return(time.Duration(1), nil)
		backendMock.EXPECT().Broadcast(e.committee.Committee(), prevoteMsg)

		e.setupCore(backendMock, e.clientAddress)
		err := e.core.handleMsg(context.Background(), proposal)
		wg.Wait()
		assert.NoError(t, err)
		e.checkState(t, e.curHeight, e.curRound, Prevote, nil, e.lockedRound, nil, int64(-1))
	})
	t.Run("receive proposal with validRound = -1 and client's lockedValue is same as proposal block", func(t *testing.T) {

		customizer := func(e *ConsensusENV) {
			e.step = Propose
			e.lockedRound = 0
			e.validRound = 0

			e.curProposal = generateBlockProposal(e.curRound, e.curHeight, -1, false, signer(e, e.curRound), member(e, e.curRound), e.previousValue.Header())
			e.lockedValue = e.curProposal.Block()
			e.validValue = e.curProposal.Block()
		}
		e := NewConsensusEnv(t, customizer)
		prevoteMsg := message.NewPrevote(e.curRound, e.curHeight.Uint64(), e.curProposal.Block().Hash(), e.clientSigner, e.clientMember, e.committeeSize)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		wg := sync.WaitGroup{}
		wg.Add(1)
		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).DoAndReturn(e.clientSigner)
		backendMock.EXPECT().ProposedBlockHash().Return(common.Hash{})
		backendMock.EXPECT().IsProposalStateCached(e.curProposal.Block().Hash()).Return(false)
		backendMock.EXPECT().ProposalVerified(e.curProposal.Block()).Do(func(i any) { wg.Done() })
		backendMock.EXPECT().VerifyProposal(e.curProposal.Block()).Return(time.Duration(1), nil)
		backendMock.EXPECT().Broadcast(e.committee.Committee(), prevoteMsg)

		e.setupCore(backendMock, e.clientAddress)
		err := e.core.handleMsg(context.Background(), e.curProposal)
		wg.Wait()
		assert.NoError(t, err)
		e.checkState(t, e.curHeight, e.curRound, Prevote, e.curProposal.Block(), e.lockedRound, e.curProposal.Block(), e.validRound)
	})
	t.Run("receive proposal with validRound = -1 and client's lockedValue is different from proposal block", func(t *testing.T) {

		customizer := func(e *ConsensusENV) {
			e.step = Propose
			e.lockedRound = 0
			e.validRound = 0
			e.lockedValue = generateBlock(e.curHeight, e.previousValue.Header())
			e.validValue = e.lockedValue
			e.curProposal = generateBlockProposal(e.curRound, e.curHeight, -1, false, signer(e, e.curRound), member(e, e.curRound), e.previousValue.Header())
		}
		e := NewConsensusEnv(t, customizer)

		prevoteMsg := message.NewPrevote(e.curRound, e.curHeight.Uint64(), common.Hash{}, e.clientSigner, e.clientMember, e.committeeSize)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).DoAndReturn(e.clientSigner)

		e.setupCore(backendMock, e.clientAddress)
		wg := sync.WaitGroup{}
		wg.Add(1)

		backendMock.EXPECT().ProposedBlockHash().Return(common.Hash{})
		backendMock.EXPECT().IsProposalStateCached(e.curProposal.Block().Hash()).Return(false)
		backendMock.EXPECT().ProposalVerified(e.curProposal.Block()).Do(func(i any) { wg.Done() })
		backendMock.EXPECT().VerifyProposal(e.curProposal.Block()).Return(time.Duration(1), nil)
		backendMock.EXPECT().Broadcast(e.committee.Committee(), prevoteMsg)

		err := e.core.handleMsg(context.Background(), e.curProposal)
		wg.Wait()
		assert.NoError(t, err)
		e.checkState(t, e.curHeight, e.curRound, Prevote, e.lockedValue, e.lockedRound, e.validValue, e.validRound)
	})
}

// The following tests aim to test lines 28 - 33 of Tendermint Algorithm described on page 6 of
// https://arxiv.org/pdf/1807.04938.pdf.
func TestOldProposal(t *testing.T) {
	t.Run("receive proposal with vr >= 0 and client's lockedRound <= vr", func(t *testing.T) {
		customizer := func(e *ConsensusENV) {
			e.step = Propose
			currentRound := int64(rand.Intn(e.committeeSize-1)) + 1
			if currentRound == int64(e.committeeSize) {
				currentRound--
			}
			e.curRound = currentRound

			// vr >= 0 && vr < round_p
			proposalValidRound := int64(rand.Intn(int(e.curRound)))
			// -1 <= c.lockedRound <= vr
			clientLockedRound := int64(rand.Intn(int(proposalValidRound+2) - 1))
			e.lockedRound = clientLockedRound
			e.validRound = clientLockedRound
			e.curProposal = generateBlockProposal(e.curRound, e.curHeight, proposalValidRound, false, signer(e, e.curRound), member(e, e.curRound), e.previousValue.Header())
		}
		e := NewConsensusEnv(t, customizer)

		prevoteMsg := message.NewPrevote(e.curRound, e.curHeight.Uint64(), e.curProposal.Block().Hash(), e.clientSigner, e.clientMember, e.committeeSize)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		wg := sync.WaitGroup{}
		wg.Add(1)

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().DoAndReturn(e.clientSigner)
		backendMock.EXPECT().VerifyProposal(e.curProposal.Block()).Return(time.Duration(1), nil)
		backendMock.EXPECT().ProposedBlockHash().Return(common.Hash{})
		backendMock.EXPECT().IsProposalStateCached(e.curProposal.Block().Hash()).Return(false)
		backendMock.EXPECT().ProposalVerified(e.curProposal.Block()).Do(func(i any) { wg.Done() })
		backendMock.EXPECT().Broadcast(e.committee.Committee(), prevoteMsg)

		e.setupCore(backendMock, e.clientAddress)
		e.core.curRoundMessages = e.core.messages.GetOrCreate(e.curRound)
		fakePrevote := message.Fake{
			FakeValue:   e.curProposal.Block().Hash(),
			FakeSigners: signersWithPower(0, e.committeeSize, e.core.CommitteeSet().Quorum()),
		}
		e.core.messages.GetOrCreate(e.curProposal.ValidRound()).AddPrevote(message.NewFakePrevote(fakePrevote))

		err := e.core.handleMsg(context.Background(), e.curProposal)
		wg.Wait()
		assert.NoError(t, err)
		e.checkState(t, e.curHeight, e.curRound, Prevote, nil, e.lockedRound, nil, e.validRound)
	})
	t.Run("receive proposal with vr >= 0 and client's lockedValue is same as proposal block", func(t *testing.T) {
		customizer := func(e *ConsensusENV) {
			e.step = Propose
			currentRound := int64(rand.Intn(e.committeeSize-1)) + 1
			if currentRound == int64(e.committeeSize) {
				currentRound--
			}
			e.curRound = currentRound

			// vr >= 0 && vr < round_p
			proposalValidRound := int64(rand.Intn(int(e.curRound)))
			// -1 <= c.lockedRound <= vr
			e.lockedRound = proposalValidRound + 1
			e.validRound = proposalValidRound + 1

			e.curProposal = generateBlockProposal(e.curRound, e.curHeight, proposalValidRound, false, signer(e, e.curRound), member(e, e.curRound), e.previousValue.Header())
			e.lockedValue = e.curProposal.Block()
			e.validValue = e.curProposal.Block()
		}
		e := NewConsensusEnv(t, customizer)

		prevoteMsg := message.NewPrevote(e.curRound, e.curHeight.Uint64(), e.curProposal.Block().Hash(), e.clientSigner, e.clientMember, e.committeeSize)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		wg := sync.WaitGroup{}
		wg.Add(1)
		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).DoAndReturn(e.clientSigner)
		backendMock.EXPECT().VerifyProposal(e.curProposal.Block()).Return(time.Duration(1), nil)
		backendMock.EXPECT().Broadcast(e.committee.Committee(), prevoteMsg)
		backendMock.EXPECT().ProposedBlockHash().Return(common.Hash{})
		backendMock.EXPECT().IsProposalStateCached(e.curProposal.Block().Hash()).Return(false)
		backendMock.EXPECT().ProposalVerified(e.curProposal.Block()).Do(func(i any) { wg.Done() })
		e.setupCore(backendMock, e.clientAddress)
		fakePrevote := message.Fake{
			FakeValue:   e.curProposal.Block().Hash(),
			FakeSigners: signersWithPower(0, e.committeeSize, e.core.CommitteeSet().Quorum()),
		}
		e.core.messages.GetOrCreate(e.curProposal.ValidRound()).AddPrevote(message.NewFakePrevote(fakePrevote))

		err := e.core.handleMsg(context.Background(), e.curProposal)
		wg.Wait()
		assert.NoError(t, err)
		e.checkState(t, e.curHeight, e.curRound, Prevote, e.curProposal.Block(), e.curProposal.ValidRound()+1, e.curProposal.Block(), e.curProposal.ValidRound()+1)
	})
	t.Run("receive proposal with vr >= 0 and client has lockedRound > vr and lockedValue != proposal", func(t *testing.T) {
		customizer := func(e *ConsensusENV) {
			e.step = Propose
			currentRound := int64(rand.Intn(e.committeeSize-1)) + 1
			if currentRound == int64(e.committeeSize) {
				currentRound--
			}
			e.curRound = currentRound
			e.lockedValue = generateBlock(e.curHeight, e.previousValue.Header())
			e.validValue = e.lockedValue
			// vr >= 0 && vr < round_p
			proposalValidRound := int64(rand.Intn(int(e.curRound)))
			e.curProposal = generateBlockProposal(e.curRound, e.curHeight, proposalValidRound, false, signer(e, e.curRound), member(e, e.curRound), e.previousValue.Header())
			e.lockedRound = proposalValidRound + 1
			e.validRound = proposalValidRound + 1
		}
		e := NewConsensusEnv(t, customizer)

		prevoteMsg := message.NewPrevote(e.curRound, e.curHeight.Uint64(), common.Hash{}, e.clientSigner, e.clientMember, e.committeeSize)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).DoAndReturn(e.clientSigner)

		e.setupCore(backendMock, e.clientAddress)
		e.core.curRoundMessages = e.core.messages.GetOrCreate(e.curRound)

		fakePrevote := message.NewFakePrevote(message.Fake{FakeSigners: signersWithPower(0, e.committeeSize, e.core.CommitteeSet().Quorum()), FakeValue: e.curProposal.Block().Hash()})
		e.core.messages.GetOrCreate(e.curProposal.ValidRound()).AddPrevote(fakePrevote)
		wg := sync.WaitGroup{}
		wg.Add(1)

		backendMock.EXPECT().VerifyProposal(e.curProposal.Block()).Return(time.Duration(0), nil)
		backendMock.EXPECT().Broadcast(e.committee.Committee(), prevoteMsg)
		backendMock.EXPECT().ProposedBlockHash().Return(common.Hash{})
		backendMock.EXPECT().IsProposalStateCached(e.curProposal.Block().Hash()).Return(false)
		backendMock.EXPECT().ProposalVerified(e.curProposal.Block()).Do(func(i any) { wg.Done() })

		err := e.core.handleMsg(context.Background(), e.curProposal)
		wg.Wait()
		assert.NoError(t, err)
		e.checkState(t, e.curHeight, e.curRound, Prevote, e.lockedValue, e.curProposal.ValidRound()+1, e.validValue, e.curProposal.ValidRound()+1)
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
		customizer := func(e *ConsensusENV) {
			clientIndex := e.committeeSize - 1
			clientAddr := e.committee.Committee().Members[clientIndex].Address
			e.clientAddress = clientAddr
			e.clientKey = e.keys[clientAddr].consensus
			e.clientSigner = signer(e, int64(clientIndex))
			e.clientMember = member(e, int64(clientIndex))

			// ensure the client is not the proposer for current round
			currentRound := int64(rand.Intn(e.committeeSize-1)) + 1
			for currentRound%int64(clientIndex) == 0 {
				currentRound = int64(rand.Intn(e.committeeSize-1)) + 1
			}

			if currentRound == int64(e.committeeSize) {
				currentRound--
			}

			// vr >= 0 && vr < round_p
			proposalValidRound := int64(0)
			if currentRound > 0 {
				proposalValidRound = int64(rand.Intn(int(currentRound)))
			}

			// -1 <= c.lockedRound < vr, if the client lockedValue = vr then the client had received the prevotes in a
			// timely manner thus there are no old prevote yet to arrive
			clientLockedRound := int64(-1) // -1
			if proposalValidRound > 0 {
				clientLockedRound = int64(rand.Intn(int(proposalValidRound)) - 1)
			}

			e.step = Propose
			e.curRound = currentRound
			e.lockedRound = clientLockedRound
			e.validRound = clientLockedRound
			// the new round proposal
			e.curProposal = generateBlockProposal(e.curRound, e.curHeight, proposalValidRound, false, signer(e, e.curRound), member(e, e.curRound), e.previousValue.Header())

			// old proposal some random block
			e.lockedValue = generateBlock(e.curHeight, e.previousValue.Header())
			e.validValue = e.lockedValue
		}
		e := NewConsensusEnv(t, customizer)

		// the old round prevote msg to be handled to get the full quorum prevote on old round vr with value v.
		prevoteMsg := message.NewPrevote(e.curProposal.ValidRound(), e.curHeight.Uint64(), e.curProposal.Block().Hash(), e.clientSigner, e.clientMember, e.committeeSize)

		// the expected prevote msg to be broadcast for the new round with <currentHeight, currentRound, proposal.Block().Hash()>
		prevoteMsgToBroadcast := message.NewPrevote(e.curRound, e.curHeight.Uint64(), e.curProposal.Block().Hash(), e.clientSigner, e.clientMember, e.committeeSize)

		ctrl := gomock.NewController(t)
		defer waitForExpects(ctrl)

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(gomock.Any()).Times(1)
		e.setupCore(backendMock, e.clientAddress)

		// construct round state with: old round's quorum-1 prevote for v on valid round.
		fakePrevote := message.Fake{
			FakeRound:     uint64(e.curProposal.ValidRound()),
			FakeSigners:   signersWithPower(1, e.committeeSize, new(big.Int).Sub(e.core.CommitteeSet().Quorum(), common.Big1)),
			FakeSignerKey: testConsensusKey.PublicKey(), // whatever key is fine
			FakeSignature: testSignature,                // whatever signature is fine
			FakeValue:     e.curProposal.Block().Hash(),
		}
		e.core.messages.GetOrCreate(e.curProposal.ValidRound()).AddPrevote(message.NewFakePrevote(fakePrevote))

		//schedule the proposer Timeout since the client is not the proposer for this round
		e.core.proposeTimeout.ScheduleTimeout(1*time.Second, e.core.Round(), e.core.Height(), e.core.onTimeoutPropose)
		wg := sync.WaitGroup{}
		wg.Add(1)

		backendMock.EXPECT().VerifyProposal(e.curProposal.Block()).Return(time.Duration(1), nil)
		backendMock.EXPECT().ProposedBlockHash().Return(common.Hash{})
		backendMock.EXPECT().IsProposalStateCached(e.curProposal.Block().Hash()).Return(false)
		backendMock.EXPECT().ProposalVerified(e.curProposal.Block()).Do(func(i any) { wg.Done() })
		backendMock.EXPECT().Broadcast(e.committee.Committee(), prevoteMsgToBroadcast)
		backendMock.EXPECT().Sign(gomock.Any()).DoAndReturn(e.clientSigner)

		// now we handle new round's proposal with round_p > vr on value v.
		err := e.core.handleMsg(context.Background(), e.curProposal)
		wg.Wait()
		assert.NoError(t, err)

		// check that the propose timeout is still started, as the proposal did not cause a step change
		assert.True(t, e.core.proposeTimeout.TimerStarted())

		// now we receive the last old round's prevote MSG to get quorum prevote on vr for value v.
		// the old round's prevote is accepted into the round state which now have the line 28 condition satisfied.
		// now to take the action of line 28 which was not align with pseudo code before.

		err = e.core.handleMsg(context.Background(), prevoteMsg)
		if !errors.Is(err, constants.ErrOldRoundMessage) {
			t.Fatalf("Expected %v, got %v", constants.ErrOldRoundMessage, err)
		}
		assert.True(t, e.core.sentPrevote)
		e.checkState(t, e.curHeight, e.curRound, Prevote, e.lockedValue, e.lockedRound, e.validValue, e.validRound)
		// now the propose timeout should be stopped, since we moved to prevote step
		assert.False(t, e.core.proposeTimeout.TimerStarted())
	})
}

func TestProposeTimeout(t *testing.T) {
	t.Run("propose Timeout is not stopped if the proposal does not cause a step change", func(t *testing.T) {
		customizer := func(e *ConsensusENV) {
			e.step = Propose
		}
		e := NewConsensusEnv(t, customizer)

		// proposal with vr > r
		proposal := generateBlockProposal(e.curRound, e.curHeight, e.curRound+1, false, signer(e, e.curRound), member(e, e.curRound), e.previousValue.Header())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		wg := sync.WaitGroup{}
		wg.Add(1)

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().ProposedBlockHash().Return(common.Hash{})
		backendMock.EXPECT().IsProposalStateCached(proposal.Block().Hash()).Return(false)
		backendMock.EXPECT().VerifyProposal(proposal.Block()).Return(time.Duration(1), nil)
		backendMock.EXPECT().ProposalVerified(proposal.Block()).Do(func(i any) { wg.Done() })
		e.setupCore(backendMock, e.clientAddress)

		// propose timer should be started
		e.core.proposeTimeout.ScheduleTimeout(5*time.Second, e.core.Round(), e.core.Height(), e.core.onTimeoutPropose)
		assert.True(t, e.core.proposeTimeout.TimerStarted())
		err := e.core.handleMsg(context.Background(), proposal)
		wg.Wait()
		assert.NoError(t, err)
		// propose timer should still be running
		assert.True(t, e.core.proposeTimeout.TimerStarted())
		e.checkState(t, e.curHeight, e.curRound, Propose, e.lockedValue, e.lockedRound, e.validValue, e.validRound)
		assert.False(t, e.core.sentPrevote)
	})
}

// The following tests aim to test lines 34 - 35 & 61 - 64 of Tendermint Algorithm described on page 6 of
// https://arxiv.org/pdf/1807.04938.pdf.
func TestPrevoteTimeout(t *testing.T) {
	t.Run("prevote Timeout started after quorum of prevotes with different hashes", func(t *testing.T) {
		customizer := func(e *ConsensusENV) {
			e.step = Prevote
		}
		e := NewConsensusEnv(t, customizer)
		lastHeader := &types.Header{Number: big.NewInt(e.curHeight.Int64()).Sub(e.curHeight, common.Big1)}
		prevoteMsg := message.NewPrevote(e.curRound, e.curHeight.Uint64(), generateBlock(e.curHeight, lastHeader).Hash(), signer(e, 1), member(e, 1), e.committeeSize)

		ctrl := gomock.NewController(t)
		defer waitForExpects(ctrl)
		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(gomock.Any()).Times(1)
		e.setupCore(backendMock, e.clientAddress)

		// create quorum prevote messages however there is no quorum on a specific hash
		prevote1 := message.Fake{
			FakeValue:   common.Hash{},
			FakeSigners: signersWithPower(2, e.committeeSize, new(big.Int).Sub(e.core.CommitteeSet().Quorum(), common.Big2)),
		}
		e.core.curRoundMessages.AddPrevote(message.NewFakePrevote(prevote1))
		prevote2 := message.Fake{
			FakeValue:   generateBlock(e.curHeight, lastHeader).Hash(),
			FakeSigners: signersWithPower(3, e.committeeSize, common.Big1),
		}
		e.core.curRoundMessages.AddPrevote(message.NewFakePrevote(prevote2))

		assert.False(t, e.core.prevoteTimeout.TimerStarted())
		err := e.core.handleMsg(context.Background(), prevoteMsg)
		assert.NoError(t, err)
		assert.True(t, e.core.prevoteTimeout.TimerStarted())

		// stop the timer to clean up
		err = e.core.prevoteTimeout.StopTimer()
		assert.NoError(t, err)
		e.checkState(t, e.curHeight, e.curRound, Prevote, e.lockedValue, e.lockedRound, e.validValue, e.validRound)
	})
	t.Run("prevote Timeout is not started multiple times", func(t *testing.T) {
		customizer := func(e *ConsensusENV) {
			e.step = Prevote
		}
		e := NewConsensusEnv(t, customizer)

		lastHeader := &types.Header{Number: big.NewInt(e.curHeight.Int64()).Sub(e.curHeight, common.Big1)}
		prevote1Msg := message.NewPrevote(e.curRound, e.curHeight.Uint64(), generateBlock(e.curHeight, lastHeader).Hash(), signer(e, 1), member(e, 1), e.committeeSize)
		prevote2Msg := message.NewPrevote(e.curRound, e.curHeight.Uint64(), generateBlock(e.curHeight, lastHeader).Hash(), signer(e, 2), member(e, 2), e.committeeSize)

		ctrl := gomock.NewController(t)
		defer waitForExpects(ctrl)

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(gomock.Any()).Times(2)
		e.setupCore(backendMock, e.clientAddress)
		// create quorum prevote messages however there is no quorum on a specific hash
		prevote1 := message.Fake{
			FakeValue:   common.Hash{},
			FakeSigners: signersWithPower(3, e.committeeSize, new(big.Int).Sub(e.core.CommitteeSet().Quorum(), common.Big2)),
		}
		e.core.curRoundMessages.AddPrevote(message.NewFakePrevote(prevote1))

		prevote2 := message.Fake{
			FakeValue:   generateBlock(e.curHeight, lastHeader).Hash(),
			FakeSigners: signersWithPower(0, e.committeeSize, common.Big1),
		}
		e.core.curRoundMessages.AddPrevote(message.NewFakePrevote(prevote2))

		assert.False(t, e.core.prevoteTimeout.TimerStarted())

		err := e.core.handleMsg(context.Background(), prevote1Msg)
		assert.NoError(t, err)
		assert.True(t, e.core.prevoteTimeout.TimerStarted())

		timeNow := time.Now()

		err = e.core.handleMsg(context.Background(), prevote2Msg)
		assert.NoError(t, err)
		assert.True(t, e.core.prevoteTimeout.TimerStarted())
		assert.True(t, e.core.prevoteTimeout.Start.Before(timeNow))

		// stop the timer to clean up
		err = e.core.prevoteTimeout.StopTimer()
		assert.NoError(t, err)
		e.checkState(t, e.curHeight, e.curRound, Prevote, e.lockedValue, e.lockedRound, e.validValue, e.validRound)
	})
	t.Run("at prevote Timeout expiry Timeout event is sent", func(t *testing.T) {
		customizer := func(e *ConsensusENV) {
			e.step = Prevote
		}
		e := NewConsensusEnv(t, customizer)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		e.setupCore(backendMock, e.clientAddress)

		assert.False(t, e.core.prevoteTimeout.TimerStarted())
		backendMock.EXPECT().Post(TimeoutEvent{RoundWhenCalled: e.curRound, HeightWhenCalled: e.curHeight, Step: Prevote})
		e.core.prevoteTimeout.ScheduleTimeout(timeoutDuration, e.core.Round(), e.core.Height(), e.core.onTimeoutPrevote)
		assert.True(t, e.core.prevoteTimeout.TimerStarted())
		e.checkState(t, e.curHeight, e.curRound, Prevote, e.lockedValue, e.lockedRound, e.validValue, e.validRound)
		time.Sleep(sleepDuration)
	})
	t.Run("at reception of prevote Timeout event precommit nil is sent", func(t *testing.T) {
		customizer := func(e *ConsensusENV) {
			e.step = Prevote
		}
		e := NewConsensusEnv(t, customizer)

		timeoutE := TimeoutEvent{RoundWhenCalled: e.curRound, HeightWhenCalled: e.curHeight, Step: Prevote}
		precommitMsg := message.NewPrecommit(e.curRound, e.curHeight.Uint64(), common.Hash{}, e.clientSigner, e.clientMember, e.committeeSize)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		backendMock := interfaces.NewMockBackend(ctrl)
		e.setupCore(backendMock, e.clientAddress)

		backendMock.EXPECT().Broadcast(e.committee.Committee(), precommitMsg)
		backendMock.EXPECT().Sign(gomock.Any()).DoAndReturn(e.clientSigner)

		e.core.handleTimeoutPrevote(context.Background(), timeoutE)
		e.checkState(t, e.curHeight, e.curRound, Precommit, e.lockedValue, e.lockedRound, e.validValue, e.validRound)
	})
}

// The following tests aim to test lines 34 - 43 of Tendermint Algorithm described on page 6 of
// https://arxiv.org/pdf/1807.04938.pdf.
func TestQuorumPrevote(t *testing.T) {
	t.Run("receive quorum prevote for proposal block when in step >= prevote", func(t *testing.T) {
		customizer := func(e *ConsensusENV) {
			//randomly choose prevote or precommit step
			e.step = Step(rand.Intn(2) + 1) //nolint:gosec
			e.curProposal = generateBlockProposal(e.curRound, e.curHeight, int64(rand.Intn(int(e.curRound+1))), false, signer(e, e.curRound), member(e, e.curRound), e.previousValue.Header())
		}
		e := NewConsensusEnv(t, customizer)

		prevoteMsg := message.NewPrevote(e.curRound, e.curHeight.Uint64(), e.curProposal.Block().Hash(), signer(e, e.curRound), member(e, e.curRound), e.committeeSize)
		precommitMsg := message.NewPrecommit(e.curRound, e.curHeight.Uint64(), e.curProposal.Block().Hash(), e.clientSigner, e.clientMember, e.committeeSize)

		ctrl := gomock.NewController(t)
		defer waitForExpects(ctrl)

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().DoAndReturn(e.clientSigner)
		backendMock.EXPECT().Post(gomock.Any()).Times(1)
		e.setupCore(backendMock, e.clientAddress)
		e.core.curRoundMessages.SetProposal(e.curProposal, true)

		fakePrevote := message.Fake{
			FakeValue:     e.curProposal.Block().Hash(),
			FakeRound:     uint64(e.curRound),
			FakeHeight:    e.curHeight.Uint64(),
			FakeSigners:   signersWithPower(uint64(int(e.curRound+1)%e.committeeSize), e.committeeSize, new(big.Int).Sub(e.core.CommitteeSet().Quorum(), common.Big1)),
			FakeSignerKey: testConsensusKey.PublicKey(), // whatever key is fine
			FakeSignature: testSignature,                // whatever signature is fine
		}
		e.core.curRoundMessages.AddPrevote(message.NewFakePrevote(fakePrevote))

		if e.step == Prevote {
			backendMock.EXPECT().Broadcast(e.committee.Committee(), precommitMsg)
			err := e.core.handleMsg(context.Background(), prevoteMsg)
			assert.NoError(t, err)
			assert.Equal(t, e.curProposal.Block(), e.core.lockedValue)
			assert.Equal(t, e.curRound, e.core.lockedRound)
			assert.Equal(t, Precommit, e.core.step)

		} else if e.step == Precommit {
			err := e.core.handleMsg(context.Background(), prevoteMsg)
			assert.NoError(t, err)
			assert.Equal(t, e.curProposal.Block(), e.core.validValue)
			assert.Equal(t, e.curRound, e.core.validRound)
		}
	})

	t.Run("receive more than quorum prevote for proposal block when in step >= prevote", func(t *testing.T) {
		customizer := func(e *ConsensusENV) {
			//randomly choose prevote or precommit step
			e.step = Step(rand.Intn(2) + 1) //nolint:gosec
			e.curProposal = generateBlockProposal(e.curRound, e.curHeight, e.curRound-1, false, signer(e, e.curRound), member(e, e.curRound), e.previousValue.Header())
		}
		e := NewConsensusEnv(t, customizer)

		prevoteMsg1 := message.NewPrevote(e.curRound, e.curHeight.Uint64(), e.curProposal.Block().Hash(), signer(e, 1), member(e, 1), e.committeeSize)
		prevoteMsg2 := message.NewPrevote(e.curRound, e.curHeight.Uint64(), e.curProposal.Block().Hash(), signer(e, 2), member(e, 2), e.committeeSize)
		precommitMsg := message.NewPrecommit(e.curRound, e.curHeight.Uint64(), e.curProposal.Block().Hash(), e.clientSigner, e.clientMember, e.committeeSize)

		ctrl := gomock.NewController(t)
		defer waitForExpects(ctrl)

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().DoAndReturn(e.clientSigner)
		backendMock.EXPECT().Post(gomock.Any()).Times(2)
		e.setupCore(backendMock, e.clientAddress)
		e.core.curRoundMessages.SetProposal(e.curProposal, true)

		fakePrevote := message.Fake{
			FakeValue:     e.curProposal.Block().Hash(),
			FakeSigners:   signersWithPower(3, e.committeeSize, new(big.Int).Sub(e.core.CommitteeSet().Quorum(), common.Big1)),
			FakeSignerKey: testConsensusKey.PublicKey(), // whatever key is fine
			FakeSignature: testSignature,                // whatever signature is fine
		}
		e.core.curRoundMessages.AddPrevote(message.NewFakePrevote(fakePrevote))

		// receive first prevote to increase the total to quorum
		if e.step == Prevote {
			backendMock.EXPECT().Broadcast(e.committee.Committee(), precommitMsg)
			err := e.core.handleMsg(context.Background(), prevoteMsg1)
			assert.NoError(t, err)
			assert.Equal(t, e.curProposal.Block(), e.core.lockedValue)
			assert.Equal(t, e.curRound, e.core.lockedRound)
			assert.Equal(t, Precommit, e.core.step)

		} else if e.step == Precommit {
			err := e.core.handleMsg(context.Background(), prevoteMsg1)
			assert.NoError(t, err)
			assert.Equal(t, e.curProposal.Block(), e.core.validValue)
			assert.Equal(t, e.curRound, e.core.validRound)
		}

		// receive second prevote to increase the total to more than quorum
		lockedValueBefore := e.core.lockedValue
		validValueBefore := e.core.validValue
		lockedRoundBefore := e.core.lockedRound
		validRoundBefore := e.core.validRound

		err := e.core.handleMsg(context.Background(), prevoteMsg2)
		assert.NoError(t, err)
		e.checkState(t, e.curHeight, e.curRound, Precommit, lockedValueBefore, lockedRoundBefore, validValueBefore, validRoundBefore)
	})
}

// The following tests aim to test lines 44 - 46 of Tendermint Algorithm described on page 6 of
// https://arxiv.org/pdf/1807.04938.pdf.
func TestQuorumPrevoteNil(t *testing.T) {

	customizer := func(e *ConsensusENV) {
		e.step = Prevote
	}
	e := NewConsensusEnv(t, customizer)

	prevoteMsg := message.NewPrevote(e.curRound, e.curHeight.Uint64(), common.Hash{}, signer(e, 1), member(e, 1), e.committeeSize)
	precommitMsg := message.NewPrecommit(e.curRound, e.curHeight.Uint64(), common.Hash{}, e.clientSigner, e.clientMember, e.committeeSize)

	ctrl := gomock.NewController(t)
	defer waitForExpects(ctrl)

	backendMock := interfaces.NewMockBackend(ctrl)
	backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().DoAndReturn(e.clientSigner)
	backendMock.EXPECT().Post(gomock.Any()).Times(1)
	e.setupCore(backendMock, e.clientAddress)

	fakePrevote := message.Fake{
		FakeValue:     common.Hash{},
		FakeSigners:   signersWithPower(2, e.committeeSize, new(big.Int).Sub(e.core.CommitteeSet().Quorum(), common.Big1)),
		FakeSignerKey: testConsensusKey.PublicKey(), // whatever key is fine
		FakeSignature: testSignature,                // whatever signature is fine
	}
	e.core.curRoundMessages.AddPrevote(message.NewFakePrevote(fakePrevote))
	backendMock.EXPECT().Broadcast(e.committee.Committee(), precommitMsg)

	err := e.core.handleMsg(context.Background(), prevoteMsg)
	assert.NoError(t, err)
	e.checkState(t, e.curHeight, e.curRound, Precommit, e.lockedValue, e.lockedRound, e.validValue, e.validRound)
}

// The following tests aim to test lines 47 - 48 & 65 - 67 of Tendermint Algorithm described on page 6 of
// https://arxiv.org/pdf/1807.04938.pdf.
func TestPrecommitTimeout(t *testing.T) {
	t.Run("at propose step, precommit Timeout started after quorum of precommits with different hashes", func(t *testing.T) {

		customizer := func(e *ConsensusENV) {
			e.step = Propose
		}
		e := NewConsensusEnv(t, customizer)
		lastHeader := &types.Header{Number: big.NewInt(e.curHeight.Int64()).Sub(e.curHeight, common.Big1)}
		precommit := message.NewPrecommit(e.curRound, e.curHeight.Uint64(), generateBlock(e.curHeight, lastHeader).Hash(), signer(e, 1), member(e, 1), e.committeeSize)
		ctrl := gomock.NewController(t)
		defer waitForExpects(ctrl)

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(gomock.Any()).Times(1)
		e.setupCore(backendMock, e.clientAddress)

		// create quorum precommit messages however there is no quorum on a specific hash
		fakePrecommit1 := message.Fake{
			FakeValue:     common.Hash{},
			FakeSigners:   signersWithPower(2, e.committeeSize, new(big.Int).Sub(e.core.CommitteeSet().Quorum(), common.Big2)),
			FakeSignerKey: testConsensusKey.PublicKey(), // whatever key is fine
			FakeSignature: testSignature,                // whatever signature is fine
		}
		e.core.curRoundMessages.AddPrecommit(message.NewFakePrecommit(fakePrecommit1))
		fakePrecommit2 := message.Fake{
			FakeValue:     generateBlock(e.curHeight, lastHeader).Hash(),
			FakeSigners:   signersWithPower(3, e.committeeSize, common.Big1),
			FakeSignerKey: testConsensusKey.PublicKey(), // whatever key is fine
			FakeSignature: testSignature,                // whatever signature is fine
		}
		e.core.curRoundMessages.AddPrecommit(message.NewFakePrecommit(fakePrecommit2))

		assert.False(t, e.core.precommitTimeout.TimerStarted())
		err := e.core.handleMsg(context.Background(), precommit)
		assert.NoError(t, err)
		assert.True(t, e.core.precommitTimeout.TimerStarted())

		// stop the timer to clean up
		err = e.core.precommitTimeout.StopTimer()
		assert.NoError(t, err)

		e.checkState(t, e.curHeight, e.curRound, Propose, e.lockedValue, e.lockedRound, e.validValue, e.validRound)
	})
	t.Run("at vote step, precommit Timeout started after quorum of precommits with different hashes", func(t *testing.T) {

		customizer := func(e *ConsensusENV) {
			e.step = Precommit
		}
		e := NewConsensusEnv(t, customizer)
		lastHeader := &types.Header{Number: big.NewInt(e.curHeight.Int64()).Sub(e.curHeight, common.Big1)}
		precommit := message.NewPrecommit(e.curRound, e.curHeight.Uint64(), generateBlock(e.curHeight, lastHeader).Hash(), signer(e, 1), member(e, 1), e.committeeSize)
		ctrl := gomock.NewController(t)
		defer waitForExpects(ctrl)

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(gomock.Any()).Times(1)
		e.setupCore(backendMock, e.clientAddress)

		// create quorum precommit messages however there is no quorum on a specific hash
		fakePrecommit1 := message.Fake{
			FakeValue:     common.Hash{},
			FakeSigners:   signersWithPower(2, e.committeeSize, new(big.Int).Sub(e.core.CommitteeSet().Quorum(), common.Big2)),
			FakeSignerKey: testConsensusKey.PublicKey(), // whatever key is fine
			FakeSignature: testSignature,                // whatever signature is fine
		}
		e.core.curRoundMessages.AddPrecommit(message.NewFakePrecommit(fakePrecommit1))
		fakePrecommit2 := message.Fake{
			FakeValue:     generateBlock(e.curHeight, lastHeader).Hash(),
			FakeSigners:   signersWithPower(3, e.committeeSize, common.Big1),
			FakeSignerKey: testConsensusKey.PublicKey(), // whatever key is fine
			FakeSignature: testSignature,                // whatever signature is fine
		}
		e.core.curRoundMessages.AddPrecommit(message.NewFakePrecommit(fakePrecommit2))

		assert.False(t, e.core.precommitTimeout.TimerStarted())
		err := e.core.handleMsg(context.Background(), precommit)
		assert.NoError(t, err)
		assert.True(t, e.core.precommitTimeout.TimerStarted())

		// stop the timer to clean up
		err = e.core.precommitTimeout.StopTimer()
		assert.NoError(t, err)

		e.checkState(t, e.curHeight, e.curRound, Precommit, e.lockedValue, e.lockedRound, e.validValue, e.validRound)
	})
	t.Run("precommit Timeout is not started multiple times", func(t *testing.T) {

		customizer := func(e *ConsensusENV) {
			e.step = Step(rand.Intn(3))
		}
		e := NewConsensusEnv(t, customizer)

		lastHeader := &types.Header{Number: big.NewInt(e.curHeight.Int64()).Sub(e.curHeight, common.Big1)}
		precommitFrom1 := message.NewPrecommit(e.curRound, e.curHeight.Uint64(), generateBlock(e.curHeight, lastHeader).Hash(), signer(e, 1), member(e, 1), e.committeeSize)
		precommitFrom2 := message.NewPrecommit(e.curRound, e.curHeight.Uint64(), generateBlock(e.curHeight, lastHeader).Hash(), signer(e, 2), member(e, 2), e.committeeSize)

		ctrl := gomock.NewController(t)
		defer waitForExpects(ctrl)

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(gomock.Any()).Times(2)
		e.setupCore(backendMock, e.clientAddress)

		// create quorum prevote messages however there is no quorum on a specific hash
		fakePrecommit1 := message.Fake{
			FakeValue:     common.Hash{},
			FakeSigners:   signersWithPower(3, e.committeeSize, new(big.Int).Sub(e.core.CommitteeSet().Quorum(), common.Big2)),
			FakeSignerKey: testConsensusKey.PublicKey(), // whatever key is fine
			FakeSignature: testSignature,                // whatever signature is fine
		}
		e.core.curRoundMessages.AddPrecommit(message.NewFakePrecommit(fakePrecommit1))
		fakePrecommit2 := message.Fake{
			FakeValue:     generateBlock(e.curHeight, lastHeader).Hash(),
			FakeSigners:   signersWithPower(0, e.committeeSize, common.Big1),
			FakeSignerKey: testConsensusKey.PublicKey(), // whatever key is fine
			FakeSignature: testSignature,                // whatever signature is fine
		}
		e.core.curRoundMessages.AddPrecommit(message.NewFakePrecommit(fakePrecommit2))

		assert.False(t, e.core.precommitTimeout.TimerStarted())

		err := e.core.handleMsg(context.Background(), precommitFrom1)
		assert.NoError(t, err)
		assert.True(t, e.core.precommitTimeout.TimerStarted())

		timeNow := time.Now()

		err = e.core.handleMsg(context.Background(), precommitFrom2)
		assert.NoError(t, err)
		assert.True(t, e.core.precommitTimeout.TimerStarted())
		assert.True(t, e.core.precommitTimeout.Start.Before(timeNow))

		// stop the timer to clean up
		err = e.core.precommitTimeout.StopTimer()
		assert.NoError(t, err)
		e.checkState(t, e.curHeight, e.curRound, e.step, e.lockedValue, e.lockedRound, e.validValue, e.validRound)
	})
	t.Run("at precommit Timeout expiry Timeout event is sent", func(t *testing.T) {
		customizer := func(e *ConsensusENV) {
			e.step = Step(rand.Intn(3))
		}
		e := NewConsensusEnv(t, customizer)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		e.setupCore(backendMock, e.clientAddress)
		assert.False(t, e.core.precommitTimeout.TimerStarted())
		backendMock.EXPECT().Post(TimeoutEvent{RoundWhenCalled: e.curRound, HeightWhenCalled: e.curHeight, Step: Precommit})
		e.core.precommitTimeout.ScheduleTimeout(timeoutDuration, e.core.Round(), e.core.Height(), e.core.onTimeoutPrecommit)
		assert.True(t, e.core.precommitTimeout.TimerStarted())
		e.checkState(t, e.curHeight, e.curRound, e.step, e.lockedValue, e.lockedRound, e.validValue, e.validRound)

		time.Sleep(sleepDuration)
	})
	t.Run("at reception of precommit Timeout event next round will be started", func(t *testing.T) {
		customizer := func(e *ConsensusENV) {
			e.step = Step(rand.Intn(3))
			// ensure client is not the proposer for next round
			currentRound := int64(rand.Intn(e.committeeSize))
			for (currentRound+1)%int64(e.committeeSize) == 0 {
				currentRound = int64(rand.Intn(e.committeeSize))
			}
			e.curRound = currentRound
		}
		e := NewConsensusEnv(t, customizer)
		timeoutE := TimeoutEvent{RoundWhenCalled: e.curRound, HeightWhenCalled: e.curHeight, Step: Precommit}
		ctrl := gomock.NewController(t)
		defer waitForExpects(ctrl)

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(gomock.Any()).Times(1)
		e.setupCore(backendMock, e.clientAddress)
		e.core.handleTimeoutPrecommit(context.Background(), timeoutE)
		e.checkState(t, e.curHeight, e.curRound+1, Propose, e.lockedValue, e.lockedRound, e.validValue, e.validRound)

		// stop the timer to clean up, since start round can start propose Timeout
		err := e.core.proposeTimeout.StopTimer()
		assert.NoError(t, err)
	})
}

// The following tests aim to test lines 49 - 54 of Tendermint Algorithm described on page 6 of
// https://arxiv.org/pdf/1807.04938.pdf.
func TestQuorumPrecommit(t *testing.T) {

	customizer := func(e *ConsensusENV) {
		e.step = Precommit
	}
	e := NewConsensusEnv(t, customizer)

	nextHeight := e.curHeight.Uint64() + 1
	lastHeader := &types.Header{Number: big.NewInt(int64(nextHeight - 1))}
	nextProposalMsg := generateBlockProposal(0, big.NewInt(int64(nextHeight)), int64(-1), false, signer(e, 0), member(e, 0), lastHeader)
	lastHeader = &types.Header{Number: big.NewInt(e.curHeight.Int64()).Sub(e.curHeight, common.Big1)}
	proposal := generateBlockProposal(e.curRound, e.curHeight, e.curRound, false, signer(e, e.curRound), member(e, e.curRound), lastHeader)
	precommit := message.NewPrecommit(e.curRound, e.curHeight.Uint64(), proposal.Block().Hash(), signer(e, 1), member(e, 1), e.committeeSize)
	sealProposal(t, proposal.Block(), e.committee, e.keys, 1)

	ctrl := gomock.NewController(t)
	defer waitForExpects(ctrl)

	backendMock := interfaces.NewMockBackend(ctrl)
	e.setupCore(backendMock, e.clientAddress)
	backendMock.EXPECT().EpochByHeight(e.core.Height().Uint64()+1).Return(e.LatestEpoch(), nil)
	e.core.curRoundMessages.SetProposal(proposal, true)

	quorumPrecommitMsgFake := message.Fake{
		FakeValue:     proposal.Block().Hash(),
		FakeSigners:   signersWithPower(2, e.committeeSize, new(big.Int).Sub(e.core.CommitteeSet().Quorum(), common.Big1)),
		FakeSignerKey: testConsensusKey.PublicKey(), // whatever key is fine
		FakeSignature: testSignature,                // whatever signature is fine
	}
	quorumPrecommitMsg := message.NewFakePrecommit(quorumPrecommitMsgFake)
	e.core.curRoundMessages.AddPrecommit(quorumPrecommitMsg)

	quorumCertificateSigners := quorumPrecommitMsg.Signers().Copy()
	quorumCertificateSigners.Merge(precommit.Signers())
	quorumCertificateSignature := blst.AggregateSignatures([]blst.Signature{quorumPrecommitMsg.Signature(), precommit.Signature()})
	backendMock.EXPECT().Commit(proposal.Block(), e.curRound, gomock.Any()).Do(
		func(proposalBlock *types.Block, round int64, quorumCertificate *types.AggregateSignature) {
			if quorumCertificateSignature.Hex() != quorumCertificate.Signature.Hex() {
				t.Fatal("quorum certificate has wrong signature")
			}
			if quorumCertificateSigners.String() != quorumCertificate.Signers.String() {
				t.Fatal("Commit called with wrong signers information")
			}
		})
	backendMock.EXPECT().ProcessFutureMsgs(nextHeight).Times(1)
	backendMock.EXPECT().Post(gomock.Any()).MaxTimes(2)
	backendMock.EXPECT().Post(TimeoutEvent{
		RoundWhenCalled:  0,
		HeightWhenCalled: new(big.Int).SetUint64(nextHeight),
		Step:             Propose,
	}).MaxTimes(1)

	err := e.core.handleMsg(context.Background(), precommit)
	assert.NoError(t, err)

	newCommitteeSet, err := tdmcommittee.NewRoundRobinSet(e.committee.Committee(), e.committee.Committee().Members[e.curRound].Address)
	e.core.committee = newCommitteeSet
	assert.NoError(t, err)
	backendMock.EXPECT().HeadBlock().Return(proposal.Block()).MaxTimes(2)

	backendMock.EXPECT().Sign(gomock.Any()).AnyTimes().DoAndReturn(signer(e, 0))
	// if the client is the next proposer
	if newCommitteeSet.GetProposer(0).Address == e.clientAddress {
		t.Log("is proposer")
		e.core.pendingCandidateBlocks[nextHeight] = nextProposalMsg.Block()
		backendMock.EXPECT().SetProposedBlockHash(nextProposalMsg.Block().Hash())
		backendMock.EXPECT().Broadcast(e.committee.Committee(), nextProposalMsg)
	}

	// It is hard to control tendermint's state machine if we construct the full backend since it overwrites the
	// state we simulated on this test context again and again. So we assume the CommitEvent is sent from miner/worker
	// thread via backend's interface, and it is handled to start new round here:
	e.core.precommiter.HandleCommit(context.Background())
	e.checkState(t, big.NewInt(int64(nextHeight)), int64(0), Propose, e.lockedValue, e.lockedRound, e.validValue, e.validRound)
}

// The following tests aim to test lines 49 - 54 of Tendermint Algorithm described on page 6 of
// https://arxiv.org/pdf/1807.04938.pdf.
func TestFutureRoundChange(t *testing.T) {
	t.Run("move to future round after receiving more than F voting power messages", func(t *testing.T) {
		customizer := func(e *ConsensusENV) {
			e.step = Step(rand.Intn(3))
			e.curRound = int64(50)
		}
		e := NewConsensusEnv(t, customizer)
		futureRound := e.curRound + 1

		// create random prevote or precommit from 2 different
		fakePrevote := message.Fake{
			FakeRound:     uint64(futureRound),
			FakeHeight:    e.curHeight.Uint64(),
			FakeSigners:   signersWithPower(1, e.committeeSize, new(big.Int).Set(e.committee.F())),
			FakeSignerKey: testConsensusKey.PublicKey(), // whatever key is fine
			FakeSignature: testSignature,                // whatever signature is fine
			FakeValue:     common.Hash{},
		}
		msg1 := message.NewFakePrevote(fakePrevote)

		msg2 := message.NewPrevote(futureRound, e.curHeight.Uint64(), common.Hash{}, signer(e, 2), member(e, 2), e.committeeSize)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(gomock.Any()).AnyTimes()
		e.setupCore(backendMock, e.clientAddress)

		err := e.core.handleMsg(context.Background(), msg1)
		assert.Equal(t, constants.ErrFutureRoundMessage, err)
		err = e.core.handleMsg(context.Background(), msg2)
		assert.Equal(t, constants.ErrFutureRoundMessage, err)

		e.checkState(t, e.curHeight, futureRound, Propose, e.lockedValue, e.lockedRound, e.validValue, e.validRound)
		assert.Equal(t, 0, len(e.core.futureRound[futureRound]))
	})

	t.Run("different messages from the same sender cannot cause round change", func(t *testing.T) {
		customizer := func(e *ConsensusENV) {
			e.step = Step(rand.Intn(3))
		}
		e := NewConsensusEnv(t, customizer)
		futureRound := e.curRound + 1

		// The collective power of the 2 messages  is more than roundChangeThreshold
		fakePrevote := message.Fake{
			FakeRound:     uint64(futureRound),
			FakeHeight:    e.curHeight.Uint64(),
			FakeSigners:   signersWithPower(1, e.committeeSize, new(big.Int).Set(e.committee.F())),
			FakeSignerKey: testConsensusKey.PublicKey(), // whatever key is fine
			FakeSignature: testSignature,                // whatever signature is fine
			FakeValue:     common.Hash{},
		}
		prevoteMsg := message.NewFakePrevote(fakePrevote)

		precommitMsg := message.NewPrecommit(futureRound, e.curHeight.Uint64(), common.Hash{}, signer(e, 1), member(e, 1), e.committeeSize)

		ctrl := gomock.NewController(t)
		defer waitForExpects(ctrl)

		backendMock := interfaces.NewMockBackend(ctrl)
		backendMock.EXPECT().Post(gomock.Any()).Times(2)
		e.setupCore(backendMock, e.clientAddress)
		err := e.core.handleMsg(context.Background(), prevoteMsg)
		assert.Equal(t, constants.ErrFutureRoundMessage, err)

		err = e.core.handleMsg(context.Background(), precommitMsg)
		assert.Equal(t, constants.ErrFutureRoundMessage, err)
		e.checkState(t, e.curHeight, e.curRound, e.step, e.lockedValue, e.lockedRound, e.validValue, e.validRound)
		assert.Equal(t, 2, len(e.core.futureRound[futureRound]))
	})
}

func sealProposal(t *testing.T, b *types.Block, c interfaces.Committee, keys AddressKeyMap, signerIndex int) {
	h := b.Header()
	hashData := types.SigHash(h)
	signature, err := crypto.Sign(hashData[:], keys[c.Committee().Members[signerIndex].Address].node)
	require.NoError(t, err)
	h.ProposerSeal = signature
	*b = *b.WithSeal(h)
}

type ConsensusENV struct {
	ConsensusView
	TendermintState
	clientAddress common.Address
	clientKey     blst.SecretKey
	clientSigner  message.Signer
	clientMember  *types.CommitteeMember
	curBlock      *types.Block
	curProposal   *message.Propose
	core          *Core
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
	env.committeeSize = rand.Intn(maxSize-minSize) + minSize
	env.committee, env.keys = prepareCommittee(t, env.committeeSize)
	lastHeader := &types.Header{Number: big.NewInt(env.previousHeight.Int64()).Sub(env.previousHeight, common.Big1)}
	env.previousValue = generateBlock(env.previousHeight, lastHeader)

	// setup initial state
	env.curHeight = big.NewInt(env.previousHeight.Int64() + 1)
	env.curRound = 0
	env.step = PrecommitDone
	env.lockedValue = nil
	env.validValue = nil
	env.lockedRound = -1
	env.validRound = -1

	sealProposal(t, env.previousValue, env.committee, env.keys, 0)

	env.clientAddress = env.committee.Committee().Members[0].Address
	env.clientMember = &env.committee.Committee().Members[0]
	env.clientKey = env.keys[env.clientAddress].consensus
	env.clientSigner = makeSigner(env.clientKey)
	env.curBlock = generateBlock(env.curHeight, env.previousValue.Header())

	if customize != nil {
		customize(env)
	}

	t.Log("curHeight", env.curHeight, "curRound", env.curRound, "committeeSize", env.committeeSize)
	return env
}

func (e *ConsensusENV) setupCore(backend interfaces.Backend, address common.Address) {
	e.core = New(backend, nil, address, log.Root(), false)

	e.core.epoch = &types.EpochInfo{
		EpochBlock: common.Big0,
		Epoch: types.Epoch{
			PreviousEpochBlock: common.Big0,
			NextEpochBlock:     new(big.Int).Add(e.curHeight, common.Big256),
			Committee:          e.committee.Committee(),
		},
	}
	e.core.setCommitteeSet(e.committee)
	e.core.setHeight(e.curHeight)
	e.core.setRound(e.curRound)
	e.core.SetValidRound(e.validRound)
	e.core.SetLockedRound(e.lockedRound)
	e.core.SetValidValue(e.validValue)
	e.core.SetLockedValue(e.lockedValue)
	e.core.SetStep(context.Background(), e.step)
}

func (e *ConsensusENV) LatestEpoch() *types.EpochInfo {
	return e.core.epoch
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
