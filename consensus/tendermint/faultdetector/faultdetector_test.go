package faultdetector

import (
	"github.com/clearmatics/autonity/autonity"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"math/big"
	"sort"
	"testing"
)

func TestSameVote(t *testing.T) {
	height := uint64(100)
	committee, keys := generateCommittee(5)
	proposer := committee[0].Address
	proposerKey := keys[proposer]
	r1 := int64(0)
	r2 := int64(1)
	validRound := int64(1)
	proposal := newProposalMessage(height, r1, validRound, proposerKey, committee, nil)
	proposal2 := newProposalMessage(height, r2, validRound, proposerKey, committee, nil)
	require.Equal(t, true, sameVote(proposal, proposal))
	require.Equal(t, false, sameVote(proposal, proposal2))
}

func TestIsProposerMsg(t *testing.T) {
	// test get proposer on height 1 since the parent block is genesis block, it elect proposer by round robin.
	height := uint64(1)
	lastHeight := height - 1
	round := int64(0)
	committee, keys := generateCommittee(5)
	parentHeader := newBlockHeader(lastHeight, committee)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	chainMock := NewMockBlockChainContext(ctrl)
	chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(parentHeader)
	sort.Sort(parentHeader.Committee)
	proposerAddr := parentHeader.Committee[round%int64(len(parentHeader.Committee))].Address
	proposal := newProposalMessage(height, round, -1, keys[proposerAddr], committee, nil)

	require.Equal(t, true, isProposerMsg(chainMock, proposal))
}

func TestDeCodeVote(t *testing.T) {
	height := uint64(100)
	round := int64(0)
	committee, keys := generateCommittee(5)
	proposal := newProposalMessage(height, round, -1, keys[committee[0].Address], committee, nil)
	vote := newVoteMsg(height, round, msgPrevote, keys[committee[0].Address], proposal.Value(), committee)
	require.NoError(t, decodeVote(vote))
}

func TestCheckMsgSignature(t *testing.T) {
	height := uint64(100)
	lastHeight := height - 1
	round := int64(0)
	committee, keys := generateCommittee(5)

	t.Run("normal case, proposal msg is checked correctly", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		currentHeader := newBlockHeader(lastHeight, committee)
		proposal := newProposalMessage(height, round, -1, keys[committee[0].Address], committee, nil)
		chainMock := NewMockBlockChainContext(ctrl)
		chainMock.EXPECT().CurrentHeader().Return(currentHeader)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(currentHeader)
		require.Nil(t, checkMsgSignature(chainMock, proposal))
	})

	t.Run("a future msg is received, expect an error of errFutureMsg", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		futureHeight := height + 1
		currentHeader := newBlockHeader(lastHeight, committee)
		proposal := newProposalMessage(futureHeight, round, -1, keys[committee[0].Address], committee, nil)
		chainMock := NewMockBlockChainContext(ctrl)
		chainMock.EXPECT().CurrentHeader().Return(currentHeader)
		require.Equal(t, errFutureMsg, checkMsgSignature(chainMock, proposal))
	})

	t.Run("chain cannot provide the last header of the height that msg votes on, expect an error of errFutureMsg", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		currentHeader := newBlockHeader(lastHeight, committee)
		proposal := newProposalMessage(height, round, -1, keys[committee[0].Address], committee, nil)
		chainMock := NewMockBlockChainContext(ctrl)
		chainMock.EXPECT().CurrentHeader().Return(currentHeader)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(nil)
		require.Equal(t, errFutureMsg, checkMsgSignature(chainMock, proposal))
	})

	t.Run("abnormal case, msg is not signed by committee", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		wrongCommitte, ks := generateCommittee(5)
		currentHeader := newBlockHeader(lastHeight, committee)
		proposal := newProposalMessage(height, round, -1, ks[wrongCommitte[0].Address], wrongCommitte, nil)
		chainMock := NewMockBlockChainContext(ctrl)
		chainMock.EXPECT().CurrentHeader().Return(currentHeader)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(currentHeader)
		require.Equal(t, errNotCommitteeMsg, checkMsgSignature(chainMock, proposal))
	})
}

func TestCheckEquivocation(t *testing.T) {
	height := uint64(100)
	round := int64(0)
	committee, keys := generateCommittee(5)

	t.Run("check equivocation with valid proof of equivocation", func(t *testing.T) {
		nilValue := common.Hash{}
		proposal := newProposalMessage(height, round, -1, keys[committee[0].Address], committee, nil)
		vote1 := newVoteMsg(height, round, msgPrevote, keys[committee[0].Address], proposal.Value(), committee)
		vote2 := newVoteMsg(height, round, msgPrevote, keys[committee[0].Address], nilValue, committee)
		var proofs []*core.Message
		proofs = append(proofs, vote2)
		require.Equal(t, errEquivocation, checkEquivocation(nil, vote1, proofs))
	})

	t.Run("check equivocation with invalid proof of equivocation", func(t *testing.T) {
		proposal := newProposalMessage(height, round, -1, keys[committee[0].Address], committee, nil)
		vote1 := newVoteMsg(height, round, msgPrevote, keys[committee[0].Address], proposal.Value(), committee)
		var proofs []*core.Message
		proofs = append(proofs, vote1)
		require.Nil(t, checkEquivocation(nil, vote1, proofs))
	})
}

func TestSubmitMisbehaviour(t *testing.T) {
	height := uint64(100)
	round := int64(0)
	committee, keys := generateCommittee(5)
	proposer := committee[0].Address
	// submit a equivocation proofs.
	proposal := newProposalMessage(height, round, -1, keys[proposer], committee, nil)
	proposal2 := newProposalMessage(height, round, -1, keys[proposer], committee, nil)
	var proofs []*core.Message
	proofs = append(proofs, proposal2)

	fd := NewFaultDetector(nil, proposer, nil)
	fd.submitMisbehavior(proposal, proofs, errEquivocation)

	require.Equal(t, 1, len(fd.onChainProofsBuffer))
	require.Equal(t, autonity.Misbehaviour, fd.onChainProofsBuffer[0].Type)
	require.Equal(t, proposer, fd.onChainProofsBuffer[0].Sender)
	require.Equal(t, proposal.MsgHash(), fd.onChainProofsBuffer[0].Msghash)
}

func TestRunRuleEngine(t *testing.T) {
	committee, keys := generateCommittee(5)
	proposer := committee[0].Address
	proposerKey := keys[proposer]
	round := int64(3)
	t.Run("test run rules with chain height less than delta height", func(t *testing.T) {
		height := uint64(deltaToWaitForAccountability - 1)
		fd := NewFaultDetector(nil, common.Address{}, nil)
		require.Equal(t, 0, len(fd.runRuleEngine(height)))
	})

	t.Run("test run rules with malicious behaviour should be detected", func(t *testing.T) {
		chainHead := uint64(100)
		checkPointHeight := chainHead - uint64(deltaToWaitForAccountability)
		lastHeader := &types.Header{Number: new(big.Int).SetUint64(checkPointHeight - 1), Committee: committee}
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chainMock := NewMockBlockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(checkPointHeight - 1).Return(lastHeader)
		fd := NewFaultDetector(chainMock, proposer, nil)

		// simulate there was a maliciousProposal at init round 0, and save to msg store.
		initProposal := newProposalMessage(checkPointHeight, 0, -1, keys[committee[1].Address], committee, nil)
		_, err := fd.msgStore.Save(initProposal)
		require.NoError(t, err)
		// simulate there were quorum preVotes for initProposal at init round 0, and save them.
		for i := 0; i < len(committee); i++ {
			preVote := newVoteMsg(checkPointHeight, 0, msgPrevote, keys[committee[i].Address], initProposal.Value(), committee)
			_, err = fd.msgStore.Save(preVote)
			require.NoError(t, err)
		}

		// Node preCommit for init Proposal at init round 0 since there were quorum preVotes for it, and save it.
		preCommit := newVoteMsg(checkPointHeight, 0, msgPrecommit, proposerKey, initProposal.Value(), committee)
		_, err = fd.msgStore.Save(preCommit)
		require.NoError(t, err)

		// While Node propose a new malicious Proposal at new round with VR as -1 which is malicious, should be addressed by rule PN.
		maliciousProposal := newProposalMessage(checkPointHeight, round, -1, proposerKey, committee, nil)
		_, err = fd.msgStore.Save(maliciousProposal)
		require.NoError(t, err)

		// Run rule engine over msg store on current height.
		onChainProofs := fd.runRuleEngine(chainHead)
		require.Equal(t, 1, len(onChainProofs))
		require.Equal(t, autonity.Misbehaviour, onChainProofs[0].Type)
		require.Equal(t, proposer, onChainProofs[0].Sender)
		require.Equal(t, maliciousProposal.MsgHash(), onChainProofs[0].Msghash)
	})
}

func TestProcessMsg(t *testing.T) {
	height := uint64(100)
	futureHeight := uint64(110)
	committee, keys := generateCommittee(5)
	proposer := committee[0].Address
	proposerKey := keys[proposer]
	round := int64(3)
	lastHeader := &types.Header{Number: new(big.Int).SetUint64(height - 1), Committee: committee}
	t.Run("test process future msg, msg should be buffered", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chainMock := NewMockBlockChainContext(ctrl)
		chainMock.EXPECT().CurrentHeader().Return(lastHeader)
		proposal := newProposalMessage(futureHeight, round, -1, proposerKey, committee, nil)

		fd := NewFaultDetector(chainMock, proposer, nil)
		require.Equal(t, errFutureMsg, fd.processMsg(proposal))
		require.Equal(t, proposal, fd.futureHeightMsg[futureHeight][0])
	})

	t.Run("test process msg, msg should be stored at msg store once verified", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chainMock := NewMockBlockChainContext(ctrl)
		chainMock.EXPECT().CurrentHeader().Return(lastHeader)
		chainMock.EXPECT().GetHeaderByNumber(height - 1).Return(lastHeader)
		proposal := newProposalMessage(height, round, -1, proposerKey, committee, nil)
		vote := newVoteMsg(height, round, msgPrevote, proposerKey, proposal.Value(), committee)

		fd := NewFaultDetector(chainMock, proposer, nil)
		require.Equal(t, nil, fd.processMsg(vote))
		require.Equal(t, vote, fd.msgStore.messages[height][round][msgPrevote][proposer][0])
	})

	t.Run("test process msg, msg should be stored at msg store once verified", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chainMock := NewMockBlockChainContext(ctrl)
		chainMock.EXPECT().CurrentHeader().AnyTimes().Return(lastHeader)
		chainMock.EXPECT().GetHeaderByNumber(height - 1).AnyTimes().Return(lastHeader)
		proposal := newProposalMessage(height, round, -1, proposerKey, committee, nil)
		vote := newVoteMsg(height, round, msgPrevote, proposerKey, proposal.Value(), committee)
		equivocatedVote := newVoteMsg(height, round, msgPrevote, proposerKey, common.Hash{}, committee)
		fd := NewFaultDetector(chainMock, proposer, nil)
		require.Equal(t, nil, fd.processMsg(vote))
		require.Equal(t, errEquivocation, fd.processMsg(equivocatedVote))
		require.Equal(t, 1, len(fd.onChainProofsBuffer))
		require.Equal(t, autonity.Misbehaviour, fd.onChainProofsBuffer[0].Type)
		require.Equal(t, proposer, fd.onChainProofsBuffer[0].Sender)
		require.Equal(t, equivocatedVote.MsgHash(), fd.onChainProofsBuffer[0].Msghash)
	})

	t.Run("test process buffered msg, msg should be removed from buffer and stored at msg store once verified", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chainMock := NewMockBlockChainContext(ctrl)
		chainMock.EXPECT().CurrentHeader().Return(lastHeader)
		chainMock.EXPECT().GetHeaderByNumber(height - 1).Return(lastHeader)
		proposal := newProposalMessage(height, round, -1, proposerKey, committee, nil)
		vote := newVoteMsg(height, round, msgPrevote, proposerKey, proposal.Value(), committee)

		fd := NewFaultDetector(chainMock, proposer, nil)
		// buffer the vote msg at afd.
		fd.futureHeightMsg[height] = append(fd.futureHeightMsg[height], vote)
		fd.processBufferedMsgs(height)
		require.Nil(t, fd.futureHeightMsg[height])
		require.Equal(t, vote, fd.msgStore.messages[height][round][msgPrevote][proposer][0])
	})
}

func TestGenerateOnChainProof(t *testing.T) {
	height := uint64(100)
	committee, keys := generateCommittee(5)
	proposer := committee[0].Address
	proposerKey := keys[proposer]
	round := int64(3)

	proposal := newProposalMessage(height, round, -1, proposerKey, committee, nil)
	equivocatedProposal := newProposalMessage(height, round, -1, proposerKey, committee, nil)
	var evidence []*core.Message
	evidence = append(evidence, equivocatedProposal)

	proof := proof{
		Type:     autonity.Misbehaviour,
		Rule:     Equivocation,
		Message:  proposal,
		Evidence: evidence,
	}

	fd := NewFaultDetector(nil, proposer, nil)

	onChainProof, err := fd.generateOnChainProof(&proof)
	require.NoError(t, err)
	require.Equal(t, autonity.Misbehaviour, onChainProof.Type)
	require.Equal(t, proposer, onChainProof.Sender)
	require.Equal(t, proposal.MsgHash(), onChainProof.Msghash)

	decodedProof, err := decodeRawProof(onChainProof.Rawproof)
	require.NoError(t, err)
	require.Equal(t, proof.Type, decodedProof.Type)
	require.Equal(t, proof.Rule, decodedProof.Rule)
	require.Equal(t, proof.Message.MsgHash(), decodedProof.Message.MsgHash())
}
