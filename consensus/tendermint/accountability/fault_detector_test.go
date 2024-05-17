package accountability

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/autonity/autonity/consensus/ethash"
	"github.com/autonity/autonity/core/rawdb"
	"github.com/autonity/autonity/core/vm"
	"math/big"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/autonity/autonity/accounts/abi/bind/backends"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	"github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/events"
	ccore "github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/params/generated"
)

var (
	committee, keys = generateCommittee()
	proposer        = committee[0].Address
	proposerKey     = keys[0]
	signer          = makeSigner(proposerKey, committee[0])
	remotePeer      = committee[1].Address
	remoteSigner    = makeSigner(keys[1], committee[1])
)

func generateCommittee() (types.Committee, []*ecdsa.PrivateKey) {
	n := 5
	validators := make(types.Committee, n)
	pkeys := make([]*ecdsa.PrivateKey, n)
	for i := 0; i < n; i++ {
		privateKey, _ := crypto.GenerateKey()
		committeeMember := types.CommitteeMember{
			Address:     crypto.PubkeyToAddress(privateKey.PublicKey),
			VotingPower: new(big.Int).SetUint64(1),
		}
		validators[i] = committeeMember
		pkeys[i] = privateKey
	}
	return validators, pkeys
}

func newBlockHeader(height uint64, committee types.Committee) *types.Header {
	// use random nonce to create different blocks
	var nonce types.BlockNonce
	for i := 0; i < len(nonce); i++ {
		nonce[i] = byte(rand.Intn(256)) //nolint
	}
	return &types.Header{
		Number:    new(big.Int).SetUint64(height),
		Nonce:     nonce,
		Committee: committee,
	}
}

// new proposal with metadata, if the withValue is not nil, it will use the value as proposal, otherwise a
// random block will be used as the value for proposal.
func newProposalMessage(h uint64, r int64, vr int64, signer message.Signer, committee types.Committee, withValue *types.Block) *message.Propose {
	block := withValue
	if withValue == nil {
		header := newBlockHeader(h, committee)
		block = types.NewBlockWithHeader(header)
	}
	return message.NewPropose(r, h, vr, block, signer)
}

func TestSameVote(t *testing.T) {
	height := uint64(100)
	r1 := int64(0)
	r2 := int64(1)
	validRound := int64(1)
	proposal := newProposalMessage(height, r1, validRound, signer, committee, nil)
	proposal2 := newProposalMessage(height, r2, validRound, signer, committee, nil)
	require.Equal(t, false, proposal.Hash() == proposal2.Hash())
}

func TestSubmitMisbehaviour(t *testing.T) {
	height := uint64(100)
	round := int64(0)
	// submit a equivocation proofs.
	proposal := newProposalMessage(height, round, -1, signer, committee, nil).MustVerify(stubVerifier).ToLight()
	proposal2 := newProposalMessage(height, round, -1, signer, committee, nil).MustVerify(stubVerifier).ToLight()
	var proofs []message.Msg
	proofs = append(proofs, proposal2)

	fd := &FaultDetector{
		misbehaviourProofCh: make(chan *autonity.AccountabilityEvent, 100),
		logger:              log.New("FaultDetector", nil),
	}

	fd.submitMisbehavior(proposal, proofs, errEquivocation)
	p := <-fd.misbehaviourProofCh

	require.Equal(t, uint8(autonity.Misbehaviour), p.EventType)
	require.Equal(t, proposer, p.Offender)
}

func TestRunRuleEngine(t *testing.T) {
	round := int64(3)
	t.Run("test run rules with malicious behaviour should be detected", func(t *testing.T) {
		chainHead := uint64(100)
		checkPointHeight := chainHead - uint64(DeltaBlocks)
		lastHeader := &types.Header{Number: new(big.Int).SetUint64(checkPointHeight - 1), Committee: committee}
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chainMock := NewMockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(checkPointHeight - 1).Return(lastHeader)
		chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})
		var blockSub event.Subscription
		chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
		fdAddr := committee[1].Address
		accountability, _ := autonity.NewAccountability(proposer, backends.NewSimulatedBackend(ccore.GenesisAlloc{fdAddr: {Balance: big.NewInt(params.Ether)}}, 10000000))

		fd := NewFaultDetector(chainMock, fdAddr, nil, core.NewMsgStore(), nil, nil, proposerKey, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())
		// store a msg before check point height in case of node is start from reset.
		msgBeforeCheckPointHeight := newProposalMessage(checkPointHeight-1, 0, -1, makeSigner(keys[1], committee[1]), committee, nil).MustVerify(stubVerifier)
		fd.msgStore.Save(msgBeforeCheckPointHeight)

		// simulate there was a maliciousProposal at init round 0, and save to msg store.
		initProposal := newProposalMessage(checkPointHeight, 0, -1, makeSigner(keys[1], committee[1]), committee, nil).MustVerify(stubVerifier)
		fd.msgStore.Save(initProposal)
		// simulate there were quorum preVotes for initProposal at init round 0, and save them.
		for i := 0; i < len(committee); i++ {
			preVote := message.NewPrevote(0, checkPointHeight, initProposal.Value(), makeSigner(keys[i], committee[i])).MustVerify(stubVerifier)
			fd.msgStore.Save(preVote)
		}

		// Node preCommit for init Proposal at init round 0 since there were quorum preVotes for it, and save it.
		preCommit := message.NewPrecommit(0, checkPointHeight, initProposal.Value(), signer).MustVerify(stubVerifier)
		fd.msgStore.Save(preCommit)

		// While Node propose a new malicious Proposal at new round with VR as -1 which is malicious, should be addressed by rule PN.
		maliciousProposal := newProposalMessage(checkPointHeight, round, -1, signer, committee, nil).MustVerify(stubVerifier)
		fd.msgStore.Save(maliciousProposal)

		// Run rule engine over msg store on current height.
		onChainProofs := fd.runRuleEngine(checkPointHeight)
		require.Equal(t, 1, len(onChainProofs))
		require.Equal(t, uint8(autonity.Misbehaviour), onChainProofs[0].EventType)
		require.Equal(t, proposer, onChainProofs[0].Offender)
		proof, err := decodeRawProof(onChainProofs[0].RawProof)
		require.NoError(t, err)
		expected := message.NewLightProposal(maliciousProposal)
		require.Equal(t, expected.Code(), proof.Message.Code())
		require.Equal(t, expected.Value(), proof.Message.Value())
		require.Equal(t, expected.R(), proof.Message.R())
		require.Equal(t, expected.H(), proof.Message.H())
	})
}

func TestProcessMsg(t *testing.T) {
	futureHeight := uint64(110)
	round := int64(3)
	t.Run("test process future msg, msg should be buffered", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chainMock := NewMockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(futureHeight - 1).Return(nil)
		chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})
		bindings, _ := autonity.NewAccountability(proposer, backends.NewSimulatedBackend(ccore.GenesisAlloc{proposer: {Balance: big.NewInt(params.Ether)}}, 10000000))

		proposal := newProposalMessage(futureHeight, round, -1, signer, committee, nil)
		var blockSub event.Subscription
		chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
		fd := NewFaultDetector(chainMock, proposer, nil, core.NewMsgStore(), nil, nil, proposerKey, &autonity.ProtocolContracts{Accountability: bindings}, log.Root())
		require.Equal(t, errFutureMsg, fd.processMsg(proposal))
		require.Equal(t, proposal, fd.futureMessages[futureHeight][0])
	})
}

func TestGenerateOnChainProof(t *testing.T) {
	height := uint64(100)
	round := int64(3)

	proposal := newProposalMessage(height, round, -1, signer, committee, nil).MustVerify(stubVerifier).ToLight()
	equivocatedProposal := newProposalMessage(height, round, -1, signer, committee, nil).MustVerify(stubVerifier).ToLight()
	var evidence []message.Msg
	evidence = append(evidence, equivocatedProposal)

	p := Proof{
		Type:      autonity.Misbehaviour,
		Rule:      autonity.Equivocation,
		Message:   proposal,
		Evidences: evidence,
	}

	fd := FaultDetector{
		address: proposer,
		logger:  log.New("FaultDetector", nil),
	}

	onChainEvent := fd.eventFromProof(&p)

	t.Run("on chain event generation", func(t *testing.T) {
		require.Equal(t, uint8(autonity.Misbehaviour), onChainEvent.EventType)
		require.Equal(t, proposer, onChainEvent.Reporter)

		decodedProof, err := decodeRawProof(onChainEvent.RawProof)
		require.NoError(t, err)
		require.Equal(t, p.Type, decodedProof.Type)
		require.Equal(t, p.Rule, decodedProof.Rule)
		require.Equal(t, p.Message.Value(), decodedProof.Message.Value())
		require.Equal(t, proposal.H(), decodedProof.Message.H())
		require.Equal(t, proposal.R(), decodedProof.Message.R())
		require.Equal(t, equivocatedProposal.H(), decodedProof.Evidences[0].H())
		require.Equal(t, equivocatedProposal.R(), decodedProof.Evidences[0].R())
	})

	t.Run("Test abi packing of onChainProof", func(t *testing.T) {
		methodName := "handleEvent"
		_, err := generated.AccountabilityAbi.Pack(methodName, onChainEvent)
		require.NoError(t, err)
	})
}

// todo: (Jason) add test to cover an accusation over a committed block scenario,
//  in such context, the accusation is considered as useless, it should be dropped.

func TestAccusationProvers(t *testing.T) {
	height := uint64(100)
	lastHeight := height - 1
	round := int64(3)
	validRound := int64(1)
	noneNilValue := common.Hash{0x1}
	lastHeader := &types.Header{Number: new(big.Int).SetUint64(lastHeight), Committee: committee}
	t.Run("innocenceProof with unprovable rule id", func(t *testing.T) {
		fd := FaultDetector{}
		var input = Proof{
			Rule: autonity.PVO12,
		}
		_, err := fd.innocenceProof(&input)
		assert.NotNil(t, err)
	})

	t.Run("innocenceProofPO have quorum preVotes", func(t *testing.T) {

		// PO: node propose an old value with an validRound, innocent Proof of it should be:
		// there were quorum num of preVote for that value at the validRound.
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chainMock := NewMockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(lastHeader)
		var blockSub event.Subscription
		chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
		chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})
		bindings, _ := autonity.NewAccountability(proposer, backends.NewSimulatedBackend(ccore.GenesisAlloc{proposer: {Balance: big.NewInt(params.Ether)}}, 10000000))

		fd := NewFaultDetector(chainMock, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}), core.NewMsgStore(), nil, nil, proposerKey, &autonity.ProtocolContracts{Accountability: bindings}, log.Root())
		// simulate a proposal message with an old value and a valid round.
		proposal := newProposalMessage(height, round, validRound, signer, committee, nil).MustVerify(stubVerifier)
		fd.msgStore.Save(proposal)

		// simulate at least quorum num of preVotes for a value at a validRound.
		for i := 0; i < len(committee); i++ {
			preVote := message.NewPrevote(validRound, height, proposal.Value(), makeSigner(keys[i], committee[i])).MustVerify(stubVerifier)
			fd.msgStore.Save(preVote)
		}

		var accusation = Proof{
			Type:    autonity.Accusation,
			Rule:    autonity.PO,
			Message: message.NewLightProposal(proposal),
		}

		proof, err := fd.innocenceProofPO(&accusation)
		assert.NoError(t, err)
		assert.Equal(t, uint8(autonity.Innocence), proof.EventType)
		assert.Equal(t, proposer, proof.Reporter)

	})

	t.Run("innocenceProofPO no quorum preVotes", func(t *testing.T) {

		// PO: node propose an old value with an validRound, innocent Proof of it should be:
		// there were quorum num of preVote for that value at the validRound.
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chainMock := NewMockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(lastHeader)
		chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})
		accountability, _ := autonity.NewAccountability(proposer, backends.NewSimulatedBackend(ccore.GenesisAlloc{proposer: {Balance: big.NewInt(params.Ether)}}, 10000000))
		var blockSub event.Subscription
		chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
		fd := NewFaultDetector(chainMock, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}), core.NewMsgStore(), nil, nil, proposerKey, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())
		// simulate a proposal message with an old value and a valid round.
		proposal := newProposalMessage(height, round, validRound, signer, committee, nil).MustVerify(stubVerifier)
		fd.msgStore.Save(proposal)

		// simulate less than quorum num of preVotes for a value at a validRound.
		preVote := message.NewPrevote(validRound, height, proposal.Value(), signer).MustVerify(stubVerifier)
		fd.msgStore.Save(preVote)

		var accusation = Proof{
			Type:    autonity.Accusation,
			Rule:    autonity.PO,
			Message: message.NewLightProposal(proposal),
		}

		_, err := fd.innocenceProofPO(&accusation)
		assert.Equal(t, errNoEvidenceForPO, err)
	})

	t.Run("innocenceProofPVN have quorum prevotes", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		// PVN: node prevote for a none nil value, then there must be a corresponding proposal.
		chainMock := NewMockChainContext(ctrl)
		chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})

		fd := FaultDetector{
			blockchain: chainMock,
			address:    proposer,
			msgStore:   core.NewMsgStore(),
			logger:     log.New("FaultDetector", nil),
		}
		// simulate a proposal message with an old value and a valid round.
		proposal := newProposalMessage(height, round, -1, signer, committee, nil).MustVerify(stubVerifier)
		fd.msgStore.Save(proposal)

		// simulate at least quorum num of preVotes for a value at a validRound.
		for i := 0; i < len(committee); i++ {
			pv := message.NewPrevote(round, height, proposal.Value(), makeSigner(keys[i], committee[i])).MustVerify(stubVerifier)
			fd.msgStore.Save(pv)
		}

		preVote := message.NewPrevote(round, height, proposal.Value(), signer).MustVerify(stubVerifier)

		var accusation = Proof{
			Type:    autonity.Accusation,
			Rule:    autonity.PVN,
			Message: preVote,
		}

		proof, err := fd.innocenceProofPVN(&accusation)
		assert.NoError(t, err)
		assert.Equal(t, uint8(autonity.Innocence), proof.EventType)
		assert.Equal(t, proposer, proof.Reporter)

	})

	t.Run("innocenceProofPVN have no corresponding proposal", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chainMock := NewMockChainContext(ctrl)
		// PVN: node prevote for a none nil value, then there must be a corresponding proposal.
		fd := FaultDetector{blockchain: chainMock, address: proposer, msgStore: core.NewMsgStore()}

		preVote := message.NewPrevote(round, height, noneNilValue, signer).MustVerify(stubVerifier)
		fd.msgStore.Save(preVote)

		var accusation = Proof{
			Type:    autonity.Accusation,
			Rule:    autonity.PVN,
			Message: preVote,
		}

		_, err := fd.innocenceProofPVN(&accusation)
		assert.Equal(t, errNoEvidenceForPVN, err)
	})

	t.Run("getInnocentProofofPVO have no quorum preVotes", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chainMock := NewMockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(lastHeader)
		var blockSub event.Subscription
		chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
		chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})
		accountability, _ := autonity.NewAccountability(proposer, backends.NewSimulatedBackend(ccore.GenesisAlloc{proposer: {Balance: big.NewInt(params.Ether)}}, 10000000))

		fd := NewFaultDetector(chainMock, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}), core.NewMsgStore(), nil, nil, proposerKey, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())

		var p Proof
		p.Rule = autonity.PVO
		oldProposal := newProposalMessage(height, 1, 0, signer, committee, nil).MustVerify(stubVerifier)
		preVote := message.NewPrevote(1, height, oldProposal.Value(), signer).MustVerify(stubVerifier)
		p.Message = preVote
		p.Evidences = append(p.Evidences, message.NewLightProposal(oldProposal))

		_, err := fd.innocenceProofPVO(&p)
		assert.Equal(t, err, errNoEvidenceForPVO)
	})

	t.Run("getInnocentProofofPVO have quorum preVotes", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chainMock := NewMockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(lastHeader)
		var blockSub event.Subscription
		chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
		chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})
		accountability, _ := autonity.NewAccountability(proposer, backends.NewSimulatedBackend(ccore.GenesisAlloc{proposer: {Balance: big.NewInt(params.Ether)}}, 10000000))

		fd := NewFaultDetector(chainMock, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}), core.NewMsgStore(), nil, nil, proposerKey, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())
		var p Proof
		p.Rule = autonity.PVO
		validRound := int64(0)
		oldProposal := newProposalMessage(height, 1, validRound, signer, committee, nil).MustVerify(stubVerifier)
		preVote := message.NewPrevote(1, height, oldProposal.Value(), signer).MustVerify(stubVerifier)
		p.Message = preVote
		p.Evidences = append(p.Evidences, message.NewLightProposal(oldProposal))

		// prepare quorum preVotes at msg store.
		for i := 0; i < len(committee); i++ {
			pv := message.NewPrevote(validRound, height, oldProposal.Value(), makeSigner(keys[i], committee[i])).MustVerify(stubVerifier)
			fd.msgStore.Save(pv)
		}

		onChainProof, err := fd.innocenceProofPVO(&p)
		assert.NoError(t, err)
		assert.Equal(t, uint8(autonity.Innocence), onChainProof.EventType)
		assert.Equal(t, proposer, onChainProof.Reporter)

	})

	t.Run("innocenceProofC1 have quorum preVotes", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chainMock := NewMockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(lastHeader)
		var blockSub event.Subscription
		chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
		chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})
		accountability, _ := autonity.NewAccountability(proposer, backends.NewSimulatedBackend(ccore.GenesisAlloc{proposer: {Balance: big.NewInt(params.Ether)}}, 10000000))

		// C1: node preCommit at a none nil value, there must be quorum corresponding preVotes with same value and round.
		fd := NewFaultDetector(chainMock, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}), core.NewMsgStore(), nil, nil, proposerKey, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())

		// simulate at least quorum num of preVotes for a value at a validRound.
		for i := 0; i < len(committee); i++ {
			preVote := message.NewPrevote(round, height, noneNilValue, makeSigner(keys[i], committee[i])).MustVerify(stubVerifier)
			fd.msgStore.Save(preVote)
		}

		preCommit := message.NewPrecommit(round, height, noneNilValue, signer).MustVerify(stubVerifier)
		fd.msgStore.Save(preCommit)

		var accusation = Proof{
			Type:    autonity.Accusation,
			Rule:    autonity.C1,
			Message: preCommit,
		}

		proof, err := fd.innocenceProofC1(&accusation)
		assert.NoError(t, err)
		assert.Equal(t, uint8(autonity.Innocence), proof.EventType)
		assert.Equal(t, proposer, proof.Reporter)

	})

	t.Run("innocenceProofC1 have no quorum preVotes", func(t *testing.T) {

		// C1: node preCommit at a none nil value, there must be quorum corresponding preVotes with same value and round.
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chainMock := NewMockChainContext(ctrl)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(lastHeader)
		var blockSub event.Subscription
		chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
		chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})
		accountability, _ := autonity.NewAccountability(proposer, backends.NewSimulatedBackend(ccore.GenesisAlloc{proposer: {Balance: big.NewInt(params.Ether)}}, 10000000))

		fd := NewFaultDetector(chainMock, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}), core.NewMsgStore(), nil, nil, proposerKey, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())

		preCommit := message.NewPrecommit(round, height, noneNilValue, signer).MustVerify(stubVerifier)
		fd.msgStore.Save(preCommit)

		var accusation = Proof{
			Type:    autonity.Accusation,
			Rule:    autonity.C1,
			Message: preCommit,
		}

		_, err := fd.innocenceProofC1(&accusation)
		assert.Equal(t, errNoEvidenceForC1, err)
	})

	t.Run("Test error to rule mapping", func(t *testing.T) {
		rule, err := errorToRule(errEquivocation)
		assert.NoError(t, err)
		assert.Equal(t, autonity.Equivocation, rule)

		rule, err = errorToRule(errProposer)
		assert.NoError(t, err)
		assert.Equal(t, autonity.InvalidProposer, rule)

		_, err = errorToRule(fmt.Errorf("unknown err"))
		assert.Error(t, err)
	})
}

// Please refer to the rules in the rule engine for each step of tendermint to understand the test context.
// TestNewProposalAccountabilityCheck, it tests the accountability events over a new proposal sent by a proposer.
func TestNewProposalAccountabilityCheck(t *testing.T) {
	height := uint64(0)
	newProposal0 := newProposalMessage(height, 3, -1, signer, committee, nil).MustVerify(stubVerifier)
	nonNilPrecommit0 := message.NewPrecommit(1, height, common.BytesToHash([]byte("test")), signer).MustVerify(stubVerifier)
	nilPrecommit0 := message.NewPrecommit(1, height, common.Hash{}, signer).MustVerify(stubVerifier)

	newProposal1 := newProposalMessage(height, 5, -1, signer, committee, nil).MustVerify(stubVerifier)
	nilPrecommit1 := message.NewPrecommit(3, height, common.Hash{}, signer).MustVerify(stubVerifier)

	newProposal0E := newProposalMessage(height, 3, 1, signer, committee, nil).MustVerify(stubVerifier)

	t.Run("misbehaviour when pi has sent a non-nil precommit in a previous round", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposal0)
		fd.msgStore.Save(nonNilPrecommit0)

		expectedProof := &Proof{
			Type:      autonity.Misbehaviour,
			Rule:      autonity.PN,
			Evidences: []message.Msg{nonNilPrecommit0},
			Message:   message.NewLightProposal(newProposal0),
		}

		proofs := fd.newProposalsAccountabilityCheck(0)
		require.Equal(t, 1, len(proofs))
		actualProof := proofs[0]
		require.Equal(t, expectedProof, actualProof)
	})

	t.Run("no proof is returned when proposal is equivocated", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposal0)
		fd.msgStore.Save(nonNilPrecommit0)
		fd.msgStore.Save(newProposal0E)

		proofs := fd.newProposalsAccountabilityCheck(0)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when pi proposes a new proposal and no precommit has been sent", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposal0)
		fd.msgStore.Save(newProposal1)

		proofs := fd.newProposalsAccountabilityCheck(0)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when pi proposes a new proposal and has sent nil precommits in previous rounds", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposal0)
		fd.msgStore.Save(nilPrecommit0)
		fd.msgStore.Save(newProposal1)
		fd.msgStore.Save(nilPrecommit1)

		proofs := fd.newProposalsAccountabilityCheck(0)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("multiple proof of misbehaviours when pi has sent non-nil precommits in previous rounds for multiple proposals", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposal0)
		fd.msgStore.Save(nonNilPrecommit0)
		fd.msgStore.Save(newProposal1)

		expectedProof0 := &Proof{
			Type:      autonity.Misbehaviour,
			Rule:      autonity.PN,
			Evidences: []message.Msg{nonNilPrecommit0},
			Message:   newProposal0,
		}

		expectedProof1 := &Proof{
			Type:      autonity.Misbehaviour,
			Rule:      autonity.PN,
			Evidences: []message.Msg{nonNilPrecommit0},
			Message:   newProposal1,
		}

		proofs := fd.newProposalsAccountabilityCheck(0)
		require.Equal(t, 2, len(proofs))

		// The order of proofs is non know apriori
		for _, p := range proofs {
			if p.Message == expectedProof0.Message {
				require.Equal(t, expectedProof0, p)
			}

			if p.Message == expectedProof1.Message {
				// The Evidences list elements can be returned in any order therefore when we have evidence which includes
				// multiple messages we need to check that each message is present separately
				require.Equal(t, expectedProof1.Type, p.Type)
				require.Equal(t, expectedProof1.Rule, p.Rule)
				require.Equal(t, expectedProof1.Message, p.Message)
				require.Contains(t, p.Evidences, nonNilPrecommit0)
			}
		}
	})
}

// TestOldProposalsAccountabilityCheck, it tests the accountability events over a proposal that was validated at previous round
func TestOldProposalsAccountabilityCheck(t *testing.T) {
	quorum := bft.Quorum(committee.TotalVotingPower())
	height := uint64(0)

	header := newBlockHeader(height, committee)
	block := types.NewBlockWithHeader(header)
	header1 := newBlockHeader(height, committee)
	block1 := types.NewBlockWithHeader(header1)

	oldProposal0 := newProposalMessage(height, 3, 0, signer, committee, block).MustVerify(stubVerifier)
	oldProposal5 := newProposalMessage(height, 5, 2, signer, committee, block).MustVerify(stubVerifier)
	oldProposal0E := newProposalMessage(height, 3, 2, signer, committee, block1).MustVerify(stubVerifier)
	oldProposal0E2 := newProposalMessage(height, 3, 0, signer, committee, block1).MustVerify(stubVerifier)

	nonNilPrecommit0V := message.NewPrecommit(0, height, block.Hash(), signer).MustVerify(stubVerifier)
	nonNilPrecommit0VPrime := message.NewPrecommit(0, height, block1.Hash(), signer).MustVerify(stubVerifier)
	nonNilPrecommit2VPrime := message.NewPrecommit(2, height, block1.Hash(), signer).MustVerify(stubVerifier)
	nonNilPrecommit1 := message.NewPrecommit(1, height, block.Hash(), signer).MustVerify(stubVerifier)

	nilPrecommit0 := message.NewPrecommit(0, height, nilValue, signer).MustVerify(stubVerifier)

	var quorumPrevotes0VPrime []message.Msg
	for i := int64(0); i < quorum.Int64(); i++ {
		prevote := message.NewPrevote(0, height, block1.Hash(), makeSigner(keys[i], committee[i])).MustVerify(stubVerifier)
		quorumPrevotes0VPrime = append(quorumPrevotes0VPrime, prevote)
	}

	var quorumPrevotes0V []message.Msg
	for i := int64(0); i < quorum.Int64(); i++ {
		prevote := message.NewPrevote(0, height, block.Hash(), makeSigner(keys[i], committee[i])).MustVerify(stubVerifier)
		quorumPrevotes0V = append(quorumPrevotes0V, prevote)
	}

	var precommiteNilAfterVR []message.Msg
	for i := 1; i < 3; i++ {
		precommit := message.NewPrecommit(int64(i), height, nilValue, signer).MustVerify(stubVerifier)
		precommiteNilAfterVR = append(precommiteNilAfterVR, precommit)
	}

	t.Run("misbehaviour when pi precommited for a different value in valid round than in the old proposal", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0)
		fd.msgStore.Save(nonNilPrecommit0VPrime)

		expectedProof := &Proof{
			Type:      autonity.Misbehaviour,
			Rule:      autonity.PO,
			Evidences: []message.Msg{nonNilPrecommit0VPrime},
			Message:   message.NewLightProposal(oldProposal0),
		}

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 1, len(proofs))
		actualProof := proofs[0]
		require.Equal(t, expectedProof, actualProof)
	})

	t.Run("misbehaviour when pi incorrectly set the valid round with a different value than the proposal", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0)
		fd.msgStore.Save(nonNilPrecommit2VPrime)

		expectedProof := &Proof{
			Type:      autonity.Misbehaviour,
			Rule:      autonity.PO,
			Evidences: []message.Msg{nonNilPrecommit2VPrime},
			Message:   message.NewLightProposal(oldProposal0),
		}

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 1, len(proofs))
		actualProof := proofs[0]
		require.Equal(t, expectedProof, actualProof)
	})

	t.Run("misbehaviour when pi incorrectly set the valid round with the same value as the proposal", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0)
		fd.msgStore.Save(nonNilPrecommit1)

		expectedProof := &Proof{
			Type:      autonity.Misbehaviour,
			Rule:      autonity.PO,
			Evidences: []message.Msg{nonNilPrecommit1},
			Message:   message.NewLightProposal(oldProposal0),
		}

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 1, len(proofs))
		actualProof := proofs[0]
		require.Equal(t, expectedProof, actualProof)
	})

	t.Run("misbehaviour when in valid round there is a quorum of prevotes for a value different than old proposal", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0)
		for _, m := range quorumPrevotes0VPrime {
			fd.msgStore.Save(m)
		}

		expectedProof := &Proof{
			Type:      autonity.Misbehaviour,
			Rule:      autonity.PO,
			Evidences: quorumPrevotes0VPrime,
			Message:   message.NewLightProposal(oldProposal0),
		}

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 1, len(proofs))
		actualProof := proofs[0]
		require.Equal(t, expectedProof.Type, actualProof.Type)
		require.Equal(t, expectedProof.Rule, actualProof.Rule)
		require.Equal(t, expectedProof.Message, actualProof.Message)
		// The order of the evidence is not known apriori
		for _, m := range expectedProof.Evidences {
			require.Contains(t, actualProof.Evidences, m)
		}

	})

	t.Run("accusation when no prevotes for proposal value in valid round", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0)

		expectedProof := &Proof{
			Type:    autonity.Accusation,
			Rule:    autonity.PO,
			Message: message.NewLightProposal(oldProposal0),
		}

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 1, len(proofs))
		actualProof := proofs[0]
		require.Equal(t, expectedProof, actualProof)
	})

	t.Run("accusation when less than quorum prevotes for proposal value in valid round", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0)
		lessThanQurorumPrevotes := quorumPrevotes0V[:len(quorumPrevotes0V)-2]
		for _, m := range lessThanQurorumPrevotes {
			fd.msgStore.Save(m)
		}

		expectedProof := &Proof{
			Type:    autonity.Accusation,
			Rule:    autonity.PO,
			Message: message.NewLightProposal(oldProposal0),
		}

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 1, len(proofs))
		actualProof := proofs[0]
		require.Equal(t, expectedProof, actualProof)
	})

	t.Run("no proof for equivocated proposal with different valid round", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0)
		fd.msgStore.Save(oldProposal0E)

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof for equivocated proposal with same valid round however different block value", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0)
		fd.msgStore.Save(oldProposal0E2)

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when quorum of prevotes for V in vr, precommit for V from pi in vr, and precommit nils from pi from vr+1 to r", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0)
		for _, m := range quorumPrevotes0V {
			fd.msgStore.Save(m)
		}
		fd.msgStore.Save(nonNilPrecommit0V)
		for _, m := range precommiteNilAfterVR {
			fd.msgStore.Save(m)
		}

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when quorum of prevotes for V in vr, precommit for V from pi in vr, and some precommit nils from pi from vr+1 to r", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0)
		for _, m := range quorumPrevotes0V {
			fd.msgStore.Save(m)
		}
		fd.msgStore.Save(nonNilPrecommit0V)
		somePrecommits := precommiteNilAfterVR[:len(precommiteNilAfterVR)-2]
		for _, m := range somePrecommits {
			fd.msgStore.Save(m)
		}

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when quorum of prevotes for V in vr, precommit for V from pi in vr", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0)
		for _, m := range quorumPrevotes0V {
			fd.msgStore.Save(m)
		}
		fd.msgStore.Save(nonNilPrecommit0V)

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when quorum of prevotes for V in vr, precommit nil from pi in vr", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0)
		for _, m := range quorumPrevotes0V {
			fd.msgStore.Save(m)
		}
		fd.msgStore.Save(nilPrecommit0)

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when quorum of prevotes for V in vr", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0)
		for _, m := range quorumPrevotes0V {
			fd.msgStore.Save(m)
		}

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("multiple proofs from different rounds", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0)
		fd.msgStore.Save(nonNilPrecommit0VPrime)

		expectedMisbehaviour := &Proof{
			Type:      autonity.Misbehaviour,
			Rule:      autonity.PO,
			Evidences: []message.Msg{nonNilPrecommit0VPrime},
			Message:   message.NewLightProposal(oldProposal0),
		}

		fd.msgStore.Save(oldProposal5)
		expectedAccusation := &Proof{
			Type:    autonity.Accusation,
			Rule:    autonity.PO,
			Message: message.NewLightProposal(oldProposal5),
		}

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 2, len(proofs))
		require.Contains(t, proofs, expectedMisbehaviour)
		require.Contains(t, proofs, expectedAccusation)
	})
}

// TestPrevotesAccountabilityCheck, it tests the accountability events over prevotes messages sent by pi.
func TestPrevotesAccountabilityCheck(t *testing.T) {
	quorum := bft.Quorum(committee.TotalVotingPower())
	height := uint64(0)

	header := newBlockHeader(height, committee)
	block := types.NewBlockWithHeader(header)
	header1 := newBlockHeader(height, committee)
	block1 := types.NewBlockWithHeader(header1)

	newProposalForB := newProposalMessage(height, 5, -1, signer, committee, block).MustVerify(stubVerifier)

	prevoteForB := message.NewPrevote(5, height, block.Hash(), signer).MustVerify(stubVerifier)
	prevoteForB1 := message.NewPrevote(5, height, block1.Hash(), signer).MustVerify(stubVerifier)

	precommitForB := message.NewPrecommit(3, height, block.Hash(), signer).MustVerify(stubVerifier)
	precommitForB1 := message.NewPrecommit(4, height, block1.Hash(), signer).MustVerify(stubVerifier)
	precommitForB1In0 := message.NewPrecommit(0, height, block1.Hash(), signer).MustVerify(stubVerifier)
	precommitForB1In1 := message.NewPrecommit(1, height, block1.Hash(), signer).MustVerify(stubVerifier)
	precommitForBIn0 := message.NewPrecommit(0, height, block.Hash(), signer).MustVerify(stubVerifier)
	precommitForBIn4 := message.NewPrecommit(4, height, block.Hash(), signer).MustVerify(stubVerifier)

	signerBis := makeSigner(keys[1], committee[1])
	oldProposalB10 := newProposalMessage(height, 10, 5, signerBis, committee, block).MustVerify(stubVerifier)
	newProposalB1In5 := newProposalMessage(height, 5, -1, signerBis, committee, block1).MustVerify(stubVerifier)
	newProposalBIn5 := newProposalMessage(height, 5, -1, signerBis, committee, block).MustVerify(stubVerifier)

	prevoteForOldB10 := message.NewPrevote(10, height, block.Hash(), signer).MustVerify(stubVerifier)
	precommitForB1In8 := message.NewPrecommit(8, height, block1.Hash(), signer).MustVerify(stubVerifier)
	precommitForBIn7 := message.NewPrecommit(7, height, block.Hash(), signer).MustVerify(stubVerifier)

	t.Run("accusation when there are no corresponding proposals", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(prevoteForB)

		expectedAccusation := &Proof{
			Type:    autonity.Accusation,
			Rule:    autonity.PVN,
			Message: prevoteForB,
		}
		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 1, len(proofs))
		require.Contains(t, proofs, expectedAccusation)
	})

	// Testcases for PVN
	t.Run("misbehaviour when pi precommited for a different value in a previous round than the prevoted value", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(prevoteForB)
		fd.msgStore.Save(precommitForB1)

		expectedMisbehaviour := &Proof{
			Type:      autonity.Misbehaviour,
			Rule:      autonity.PVN,
			Evidences: []message.Msg{message.NewLightProposal(newProposalForB), precommitForB1},
			Message:   prevoteForB,
		}

		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 1, len(proofs))
		require.Equal(t, expectedMisbehaviour, proofs[0])
	})

	t.Run("misbehaviour when pi precommited for a different value in a previous round than the prevoted value, after a flip flop", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(prevoteForB)
		fd.msgStore.Save(precommitForB1)
		fd.msgStore.Save(precommitForB)

		expectedMisbehaviour := &Proof{
			Type:      autonity.Misbehaviour,
			Rule:      autonity.PVN,
			Evidences: []message.Msg{message.NewLightProposal(newProposalForB), precommitForB1},
			Message:   prevoteForB,
		}

		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 1, len(proofs))
		require.Equal(t, expectedMisbehaviour, proofs[0])
	})

	t.Run("misbehaviour when pi precommited for a different value in a previous round than the prevoted value while precommit nils in middle rounds", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(prevoteForB)
		fd.msgStore.Save(precommitForB1In0)

		var precommitNilsAfter0 []message.Msg
		for i := 1; i < 5; i++ {
			precommitNil := message.NewPrecommit(int64(i), height, nilValue, signer).MustVerify(stubVerifier)
			precommitNilsAfter0 = append(precommitNilsAfter0, precommitNil)
			fd.msgStore.Save(precommitNil)
		}

		expectedMisbehaviour := &Proof{
			Type:      autonity.Misbehaviour,
			Rule:      autonity.PVN,
			Evidences: []message.Msg{precommitForB1In0},
			Message:   prevoteForB,
		}
		expectedMisbehaviour.Evidences = append(expectedMisbehaviour.Evidences, precommitNilsAfter0...)

		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 1, len(proofs))
		actualProof := proofs[0]
		require.Equal(t, expectedMisbehaviour.Type, actualProof.Type)
		require.Equal(t, expectedMisbehaviour.Rule, actualProof.Rule)
		require.Equal(t, expectedMisbehaviour.Message, actualProof.Message)
		for _, m := range expectedMisbehaviour.Evidences {
			require.Contains(t, actualProof.Evidences, m)
		}
	})

	t.Run("misbehaviour when pi precommited for a different value in a previous round than the prevoted value, after a flip flop, while precommit nils in middle rounds", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(prevoteForB)
		fd.msgStore.Save(precommitForBIn0)
		fd.msgStore.Save(precommitForB1In1)

		var precommitNilsAfter1 []message.Msg
		for i := 2; i < 5; i++ {
			precommitNil := message.NewPrecommit(int64(i), height, nilValue, signer).MustVerify(stubVerifier)
			precommitNilsAfter1 = append(precommitNilsAfter1, precommitNil)
			fd.msgStore.Save(precommitNil)
		}

		expectedMisbehaviour := &Proof{
			Type:      autonity.Misbehaviour,
			Rule:      autonity.PVN,
			Evidences: []message.Msg{precommitForB1In1},
			Message:   prevoteForB,
		}
		expectedMisbehaviour.Evidences = append(expectedMisbehaviour.Evidences, precommitNilsAfter1...)

		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 1, len(proofs))
		actualProof := proofs[0]
		require.Equal(t, expectedMisbehaviour.Type, actualProof.Type)
		require.Equal(t, expectedMisbehaviour.Rule, actualProof.Rule)
		require.Equal(t, expectedMisbehaviour.Message, actualProof.Message)
		for _, m := range expectedMisbehaviour.Evidences {
			require.Contains(t, actualProof.Evidences, m)
		}
	})

	t.Run("no proof when pi precommited for the same value as the prevoted value in a previous round", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(prevoteForB)
		fd.msgStore.Save(precommitForBIn4)

		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when pi precommited for the same value as the prevoted value in a previous round with missing precommits in middle rounds", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(prevoteForB)
		fd.msgStore.Save(precommitForBIn0)

		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when pi precommited for the same value as the prevoted value in a previous round with some missing precommits and precommit nils in middle rounds", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(prevoteForB)
		fd.msgStore.Save(precommitForBIn0)
		fd.msgStore.Save(message.NewPrecommit(3, height, nilValue, signer).MustVerify(stubVerifier))

		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when pi precommited for the same value as the prevoted value in a previous round with no missing precommits in middle rounds", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(prevoteForB)
		fd.msgStore.Save(precommitForBIn0)
		for i := 1; i < 5; i++ {
			precommitNil := message.NewPrecommit(int64(i), height, nilValue, signer).MustVerify(stubVerifier)
			fd.msgStore.Save(precommitNil)
		}

		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when pi precommited for {B1,nil*,B} and then prevoted B", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(prevoteForB)

		fd.msgStore.Save(precommitForB1In0)

		// fill gaps with nil
		for i := 1; i < 4; i++ {
			precommitNil := message.NewPrecommit(int64(i), height, nilValue, signer).MustVerify(stubVerifier)
			fd.msgStore.Save(precommitNil)
		}

		fd.msgStore.Save(precommitForBIn4)

		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	// Testcases for PVO
	t.Run("accusation when there is no quorum for the prevote value in the valid round", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposalB10)
		fd.msgStore.Save(prevoteForOldB10)

		expectedAccusation := &Proof{
			Type:      autonity.Accusation,
			Rule:      autonity.PVO,
			Message:   prevoteForOldB10,
			Evidences: []message.Msg{message.NewLightProposal(oldProposalB10)},
		}

		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 1, len(proofs))
		acutalProof := proofs[0]
		require.Equal(t, expectedAccusation, acutalProof)
	})

	t.Run("misbehaviour when pi prevotes for an old proposal while in the valid round there is quorum for different value", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposalB10)
		// Need to add this new proposal in valid round so that unwanted accusation are not returned by the prevotes
		// accountability check method. Since we are adding a quorum of prevotes in round 6 we also need to add a new
		// proposal in round 6 to allow for those prevotes to not return accusations.
		fd.msgStore.Save(newProposalB1In5)
		fd.msgStore.Save(prevoteForOldB10)
		// quorum of prevotes for B1 in vr = 6
		var vr5Prevotes []message.Msg
		for i := uint64(0); i < quorum.Uint64(); i++ {
			vr6Prevote := message.NewPrevote(5, height, block1.Hash(), makeSigner(keys[i], committee[i])).MustVerify(stubVerifier)
			vr5Prevotes = append(vr5Prevotes, vr6Prevote)
			fd.msgStore.Save(vr6Prevote)
		}

		expectedMisbehaviour := &Proof{
			Type:    autonity.Misbehaviour,
			Rule:    autonity.PVO,
			Message: prevoteForOldB10,
		}
		expectedMisbehaviour.Evidences = append(expectedMisbehaviour.Evidences, message.NewLightProposal(oldProposalB10))
		expectedMisbehaviour.Evidences = append(expectedMisbehaviour.Evidences, vr5Prevotes...)

		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 1, len(proofs))
		actualProof := proofs[0]
		require.Equal(t, expectedMisbehaviour.Type, actualProof.Type)
		require.Equal(t, expectedMisbehaviour.Rule, actualProof.Rule)
		require.Equal(t, expectedMisbehaviour.Message, actualProof.Message)
		for _, m := range expectedMisbehaviour.Evidences {
			require.Contains(t, actualProof.Evidences, m)
		}

	})

	t.Run("misbehaviour when pi has precommited for V in a previous round however the latest precommit from pi is not for V yet pi still prevoted for V in the current round", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposalB10)
		fd.msgStore.Save(prevoteForOldB10)
		fd.msgStore.Save(newProposalBIn5)
		for i := 0; i < len(committee); i++ {
			fd.msgStore.Save(message.NewPrevote(5, height, block.Hash(), makeSigner(keys[i], committee[i])).MustVerify(stubVerifier))
		}
		for i := newProposalBIn5.R(); i < precommitForBIn7.R(); i++ {
			fd.msgStore.Save(message.NewPrecommit(i, height, nilValue, signer).MustVerify(stubVerifier))
		}
		var precommitsFromPiAfterLatestPrecommitForB []message.Msg
		fd.msgStore.Save(precommitForBIn7)

		precommitsFromPiAfterLatestPrecommitForB = append(precommitsFromPiAfterLatestPrecommitForB, precommitForBIn7)
		fd.msgStore.Save(precommitForB1In8)
		precommitsFromPiAfterLatestPrecommitForB = append(precommitsFromPiAfterLatestPrecommitForB, precommitForB1In8)
		p := message.NewPrecommit(precommitForB1In8.R()+1, height, nilValue, signer).MustVerify(stubVerifier)
		fd.msgStore.Save(p)
		precommitsFromPiAfterLatestPrecommitForB = append(precommitsFromPiAfterLatestPrecommitForB, p)

		expectedMisbehaviour := &Proof{
			Type:    autonity.Misbehaviour,
			Rule:    autonity.PVO12,
			Message: prevoteForOldB10,
		}
		expectedMisbehaviour.Evidences = append(expectedMisbehaviour.Evidences, message.NewLightProposal(oldProposalB10))
		expectedMisbehaviour.Evidences = append(expectedMisbehaviour.Evidences, precommitsFromPiAfterLatestPrecommitForB...)

		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 1, len(proofs))
		actualProof := proofs[0]
		require.Equal(t, expectedMisbehaviour.Type, actualProof.Type)
		require.Equal(t, expectedMisbehaviour.Rule, actualProof.Rule)
		require.Equal(t, expectedMisbehaviour.Message, actualProof.Message)
		for _, m := range expectedMisbehaviour.Evidences {
			require.Contains(t, actualProof.Evidences, m)
		}
	})

	t.Run("no proof when pi has precommited for V in a previous round and precommit nils afterwards", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposalB10)
		fd.msgStore.Save(prevoteForOldB10)
		fd.msgStore.Save(newProposalBIn5)
		for i := 0; i < len(committee); i++ {
			fd.msgStore.Save(message.NewPrevote(5, height, block.Hash(), makeSigner(keys[i], committee[i])).MustVerify(stubVerifier))
		}
		fd.msgStore.Save(precommitForBIn7)
		for i := precommitForBIn7.R() + 1; i < oldProposalB10.R(); i++ {
			fd.msgStore.Save(message.NewPrecommit(i, height, nilValue, signer).MustVerify(stubVerifier))
		}

		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))

	})

	t.Run("no proof when pi has precommited for V in a previous round however the latest precommit from pi is not for V yet pi still prevoted for V in the current round"+
		" but there are missing message after latest precommit for V", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposalB10)
		fd.msgStore.Save(prevoteForOldB10)
		fd.msgStore.Save(newProposalBIn5)
		for i := 0; i < len(committee); i++ {
			fd.msgStore.Save(message.NewPrevote(5, height, block.Hash(), makeSigner(keys[i], committee[i])).MustVerify(stubVerifier))
		}
		fd.msgStore.Save(precommitForBIn7)
		fd.msgStore.Save(precommitForB1In8)

		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))

	})

	t.Run("misbehaviour when pi has never precommited for V in a previous round however pi prevoted for V which is being reproposed", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposalB10)
		fd.msgStore.Save(prevoteForOldB10)
		fd.msgStore.Save(newProposalBIn5)
		for i := 0; i < len(committee); i++ {
			fd.msgStore.Save(message.NewPrevote(5, height, block.Hash(), makeSigner(keys[i], committee[i])).MustVerify(stubVerifier))
		}

		var precommitsFromPiAfterVR []message.Msg
		for i := newProposalBIn5.R() + 1; i < precommitForB1In8.R(); i++ {
			p := message.NewPrecommit(i, height, nilValue, signer).MustVerify(stubVerifier)
			fd.msgStore.Save(p)
			precommitsFromPiAfterVR = append(precommitsFromPiAfterVR, p)
		}
		fd.msgStore.Save(precommitForB1In8)
		precommitsFromPiAfterVR = append(precommitsFromPiAfterVR, precommitForB1In8)
		p := message.NewPrecommit(precommitForB1In8.R()+1, height, nilValue, signer).MustVerify(stubVerifier)
		fd.msgStore.Save(p)
		precommitsFromPiAfterVR = append(precommitsFromPiAfterVR, p)

		expectedMisbehaviour := &Proof{
			Type:    autonity.Misbehaviour,
			Rule:    autonity.PVO12,
			Message: prevoteForOldB10,
		}
		expectedMisbehaviour.Evidences = append(expectedMisbehaviour.Evidences, message.NewLightProposal(oldProposalB10))
		expectedMisbehaviour.Evidences = append(expectedMisbehaviour.Evidences, precommitsFromPiAfterVR...)

		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 1, len(proofs))
		actualProof := proofs[0]
		require.Equal(t, expectedMisbehaviour.Type, actualProof.Type)
		require.Equal(t, expectedMisbehaviour.Rule, actualProof.Rule)
		require.Equal(t, expectedMisbehaviour.Message, actualProof.Message)
		for _, m := range expectedMisbehaviour.Evidences {
			require.Contains(t, actualProof.Evidences, m)
		}
	})

	t.Run("no proof when pi has never precommited for V in a previous round however has precommitted nil after VR", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposalB10)
		fd.msgStore.Save(prevoteForOldB10)
		fd.msgStore.Save(newProposalBIn5)
		for i := 0; i < len(committee); i++ {
			fd.msgStore.Save(message.NewPrevote(5, height, block.Hash(), makeSigner(keys[i], committee[i])).MustVerify(stubVerifier))
		}

		for i := newProposalBIn5.R() + 1; i < oldProposalB10.R(); i++ {
			fd.msgStore.Save(message.NewPrecommit(i, height, nilValue, signer).MustVerify(stubVerifier))
		}

		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when pi has never precommited for V in a previous round however pi prevoted for V while it has precommited for V' but there are missing precommit before precommit for V'", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposalB10)
		fd.msgStore.Save(prevoteForOldB10)
		fd.msgStore.Save(newProposalBIn5)
		for i := 0; i < len(committee); i++ {
			fd.msgStore.Save(message.NewPrevote(5, height, block.Hash(), makeSigner(keys[i], committee[i])).MustVerify(stubVerifier))
		}

		fd.msgStore.Save(precommitForB1In8)

		p := message.NewPrecommit(precommitForB1In8.R()+1, height, nilValue, signer).MustVerify(stubVerifier)
		fd.msgStore.Save(p)

		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when pi has never precommited for V in a previous round however pi prevoted for V while it has precommited for V' but there are missing precommit after precommit for V'", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposalB10)

		fd.msgStore.Save(prevoteForOldB10)
		fd.msgStore.Save(newProposalBIn5)
		for i := 0; i < len(committee); i++ {
			fd.msgStore.Save(message.NewPrevote(5, height, block.Hash(), makeSigner(keys[i], committee[i])).MustVerify(stubVerifier))
		}

		for i := newProposalBIn5.R() + 1; i < precommitForB1In8.R(); i++ {
			fd.msgStore.Save(message.NewPrecommit(i, height, nilValue, signer).MustVerify(stubVerifier))
		}
		fd.msgStore.Save(precommitForB1In8)

		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("prevotes accountability check can return multiple proofs", func(t *testing.T) {
		fd := testFD()

		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(prevoteForB)
		fd.msgStore.Save(precommitForB1)
		fd.msgStore.Save(precommitForB)

		fd.msgStore.Save(oldProposalB10)
		fd.msgStore.Save(prevoteForOldB10)
		for i := 0; i < len(committee); i++ {
			fd.msgStore.Save(message.NewPrevote(6, height, block1.Hash(), makeSigner(keys[i], committee[i])).MustVerify(stubVerifier))
		}

		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 2, len(proofs))
	})

	t.Run("no proof when prevote is equivocated with different values", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(prevoteForB)
		fd.msgStore.Save(prevoteForB1)

		proofs := fd.prevotesAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})
}

// TestPrecommitsAccountabilityCheck, it tests the accountability events over precommit messages sent by pi.
func TestPrecommitsAccountabilityCheck(t *testing.T) {
	quorum := bft.Quorum(committee.TotalVotingPower())
	height := uint64(0)

	header := newBlockHeader(height, committee)
	block := types.NewBlockWithHeader(header)
	header1 := newBlockHeader(height, committee)
	block1 := types.NewBlockWithHeader(header1)

	newProposalForB := newProposalMessage(height, 2, -1, makeSigner(keys[1], committee[1]), committee, block).MustVerify(stubVerifier)

	precommitForB := message.NewPrecommit(2, height, block.Hash(), signer).MustVerify(stubVerifier)
	precommitForB1 := message.NewPrecommit(2, height, block1.Hash(), signer).MustVerify(stubVerifier)
	precommitForB1In3 := message.NewPrecommit(3, height, block1.Hash(), signer).MustVerify(stubVerifier)

	t.Run("accusation when prevotes is less than quorum", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(precommitForB)

		for i := int64(0); i < quorum.Int64()-1; i++ {
			fd.msgStore.Save(message.NewPrevote(2, height, block.Hash(), makeSigner(keys[i], committee[i])).MustVerify(stubVerifier))
		}

		expectedAccusation := &Proof{
			Type:    autonity.Accusation,
			Rule:    autonity.C1,
			Message: precommitForB,
		}
		proofs := fd.precommitsAccountabilityCheck(height, quorum)
		require.Equal(t, 1, len(proofs))
		require.Equal(t, expectedAccusation, proofs[0])
	})

	t.Run("misbehaviour when there is a quorum for V' than what pi precommitted for", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(precommitForB)

		var prevotesForB1 []message.Msg
		for i := int64(0); i < quorum.Int64(); i++ {
			p := message.NewPrevote(2, height, block1.Hash(), makeSigner(keys[i], committee[i])).MustVerify(stubVerifier)
			fd.msgStore.Save(p)
			prevotesForB1 = append(prevotesForB1, p)
		}

		expectedMisbehaviour := &Proof{
			Type:      autonity.Misbehaviour,
			Rule:      autonity.C,
			Evidences: prevotesForB1,
			Message:   precommitForB,
		}
		proofs := fd.precommitsAccountabilityCheck(height, quorum)
		require.Equal(t, 1, len(proofs))
		actualProof := proofs[0]
		require.Equal(t, expectedMisbehaviour.Type, actualProof.Type)
		require.Equal(t, expectedMisbehaviour.Rule, actualProof.Rule)
		require.Equal(t, expectedMisbehaviour.Message, actualProof.Message)
		for _, m := range expectedMisbehaviour.Evidences {
			require.Contains(t, actualProof.Evidences, m)
		}
	})

	t.Run("multiple proofs can be returned from precommits accountability check", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(precommitForB1In3)

		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(precommitForB)

		var prevotesForB1 []message.Msg
		for i := int64(0); i < quorum.Int64(); i++ {
			p := message.NewPrevote(2, height, block1.Hash(), makeSigner(keys[i], committee[i])).MustVerify(stubVerifier)
			fd.msgStore.Save(p)
			prevotesForB1 = append(prevotesForB1, p)
		}

		expectedProof0 := &Proof{
			Type:      autonity.Misbehaviour,
			Rule:      autonity.C,
			Evidences: prevotesForB1,
			Message:   precommitForB,
		}

		expectedProof1 := &Proof{
			Type:    autonity.Accusation,
			Rule:    autonity.C1,
			Message: precommitForB1In3,
		}
		proofs := fd.precommitsAccountabilityCheck(height, quorum)
		require.Equal(t, 2, len(proofs))

		for _, p := range proofs {
			if p.Message == expectedProof1.Message {
				require.Equal(t, expectedProof1, p)
			}

			if p.Message == expectedProof0.Message {
				// The Evidences list elements can be returned in any order therefore when we have evidence which includes
				// multiple messages we need to check that each message is present separately
				require.Equal(t, expectedProof0.Type, p.Type)
				require.Equal(t, expectedProof0.Rule, p.Rule)
				require.Equal(t, expectedProof0.Message, p.Message)

				for _, m := range expectedProof0.Evidences {
					require.Contains(t, p.Evidences, m)
				}
			}
		}
	})

	t.Run("no proof when there is enough prevotes to form a quorum", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(precommitForB)

		for i := 0; i < len(committee); i++ {
			fd.msgStore.Save(message.NewPrevote(2, height, block.Hash(), makeSigner(keys[i], committee[i])).MustVerify(stubVerifier))
		}

		proofs := fd.precommitsAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when there is more than quorum prevotes ", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(precommitForB)

		for i := 0; i < len(committee); i++ {
			fd.msgStore.Save(message.NewPrevote(2, height, block.Hash(), makeSigner(keys[i], committee[i])).MustVerify(stubVerifier))
		}

		proofs := fd.precommitsAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when precommit is equivocated with different values", func(t *testing.T) {
		//t.Skip("not stable in CI, but work in local.")
		fd := testFD()
		fd.msgStore.Save(precommitForB)
		fd.msgStore.Save(precommitForB1)

		proofs := fd.precommitsAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})
}

func testFD() *FaultDetector {
	return &FaultDetector{
		msgStore:   core.NewMsgStore(),
		logger:     log.Root(),
		blockchain: newTestBlockchain(),
	}
}

func newTestBlockchain() *ccore.BlockChain {
	db := rawdb.NewMemoryDatabase()
	ccore.GenesisBlockForTesting(db, common.Address{}, common.Big0)

	chain, err := ccore.NewBlockChain(db, nil, params.TestChainConfig, ethash.NewFaker(), vm.Config{}, nil, &ccore.TxSenderCacher{}, nil, backends.NewInternalBackend(nil), log.Root())
	if err != nil {
		panic(err)
	}
	return chain
}
