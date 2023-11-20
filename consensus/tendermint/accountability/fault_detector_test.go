package accountability

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/autonity/autonity/accounts/abi/bind/backends"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	proto "github.com/autonity/autonity/consensus"
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
	vals := make(types.Committee, n)
	keys := make([]*ecdsa.PrivateKey, n)
	for i := 0; i < n; i++ {
		privateKey, _ := crypto.GenerateKey()
		committeeMember := types.CommitteeMember{
			Address:     crypto.PubkeyToAddress(privateKey.PublicKey),
			VotingPower: new(big.Int).SetUint64(1),
		}
		vals[i] = committeeMember
		keys[i] = privateKey
	}
	return vals, keys
}

func newBlockHeader(height uint64, committee types.Committee) *types.Header {
	// use random nonce to create different blocks
	var nonce types.BlockNonce
	for i := 0; i < len(nonce); i++ {
		nonce[0] = byte(rand.Intn(256)) //nolint
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
		misbehaviourProofsCh: make(chan *autonity.AccountabilityEvent, 100),
		logger:               log.New("FaultDetector", nil),
	}

	fd.submitMisbehavior(proposal, proofs, errEquivocation)
	p := <-fd.misbehaviourProofsCh

	require.Equal(t, uint8(autonity.Misbehaviour), p.EventType)
	require.Equal(t, proposer, p.Offender)
}

func TestRunRuleEngine(t *testing.T) {
	round := int64(3)
	t.Run("test run rules with malicious behaviour should be detected", func(t *testing.T) {
		chainHead := uint64(100)
		checkPointHeight := chainHead - uint64(proto.DeltaBlocks)
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
		require.Equal(t, proposal, fd.futureHeightMsgs[futureHeight][0])
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

func TestRuleEngine(t *testing.T) {
	height := uint64(100)
	lastHeight := height - 1
	round := int64(3)
	validRound := int64(1)
	totalPower := big.NewInt(int64(len(committee)))
	noneNilValue := common.Hash{0x1}
	lastHeader := &types.Header{Number: new(big.Int).SetUint64(lastHeight), Committee: committee}
	maliciousSigner := makeSigner(keys[1], committee[1])
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

	t.Run("RunRule address the misbehaviour of PN rule", func(t *testing.T) {
		// ------------New Proposal------------
		// PN:  (Mr′<r,P C|pi)∗ <--- (Mr,P|pi)
		// PN1: [nil ∨ ⊥] <--- [V]
		// If one send a maliciousProposal for a new V, then all preCommits for previous rounds from this sender are nil.
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chainMock := NewMockChainContext(ctrl)
		var blockSub event.Subscription
		chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
		chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})
		accountability, _ := autonity.NewAccountability(proposer, backends.NewSimulatedBackend(ccore.GenesisAlloc{proposer: {Balance: big.NewInt(params.Ether)}}, 10000000))

		fd := NewFaultDetector(chainMock, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}), core.NewMsgStore(), nil, nil, proposerKey, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())
		quorum := bft.Quorum(totalPower)

		// simulate there was a maliciousProposal at init round 0, and save to msg store.
		initProposal := newProposalMessage(height, 0, -1, makeSigner(keys[1], committee[1]), committee, nil).MustVerify(stubVerifier)
		fd.msgStore.Save(initProposal)
		// simulate there were quorum preVotes for initProposal at init round 0, and save them.
		for i := 0; i < len(committee); i++ {
			preVote := message.NewPrevote(0, height, initProposal.Value(), makeSigner(keys[i], committee[i])).MustVerify(stubVerifier)
			fd.msgStore.Save(preVote)
		}

		// Node preCommit for init Proposal at init round 0 since there were quorum preVotes for it, and save it.
		preCommit := message.NewPrecommit(0, height, initProposal.Value(), signer).MustVerify(stubVerifier)
		fd.msgStore.Save(preCommit)

		// While Node propose a new malicious Proposal at new round with VR as -1 which is malicious, should be addressed by rule PN.
		maliciousProposal := newProposalMessage(height, round, -1, signer, committee, nil).MustVerify(stubVerifier)
		fd.msgStore.Save(maliciousProposal)

		// Run rule engine over msg store on height.
		onChainProofs := fd.runRulesOverHeight(height, quorum)
		assert.Equal(t, 1, len(onChainProofs))
		assert.Equal(t, autonity.Misbehaviour, onChainProofs[0].Type)
		assert.Equal(t, autonity.PN, onChainProofs[0].Rule)
		assert.Equal(t, message.NewLightProposal(maliciousProposal).Signature(), onChainProofs[0].Message.Signature())
		assert.Equal(t, preCommit.Signature(), onChainProofs[0].Evidences[0].Signature())
	})

	t.Run("RunRule address the misbehaviour of PO rule, the old value proposed is not locked", func(t *testing.T) {
		// ------------Old Proposal------------
		// PO: (Mr′<r,PV) ∧ (Mr′,PC|pi) ∧ (Mr′<r′′<r,P C|pi)∗ <--- (Mr,P|pi)
		// PO1: [#(Mr′,PV|V) ≥ 2f+ 1] ∧ [nil ∨ V ∨ ⊥] ∧ [nil ∨ ⊥] <--- [V]

		// to address below scenario:
		// Is there a precommit for a value other than nil or the proposed value
		// by the current proposer in the valid round? If there is the proposer
		// has proposed a value for which it is not locked on, thus a Proof of
		// misbehaviour can be generated.
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chainMock := NewMockChainContext(ctrl)
		var blockSub event.Subscription
		chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
		chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})
		accountability, _ := autonity.NewAccountability(proposer, backends.NewSimulatedBackend(ccore.GenesisAlloc{proposer: {Balance: big.NewInt(params.Ether)}}, 10000000))

		fd := NewFaultDetector(chainMock, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}), core.NewMsgStore(), nil, nil, proposerKey, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())
		quorum := bft.Quorum(totalPower)

		// simulate an init proposal at r: 0, with v1.
		initProposal := newProposalMessage(height, 0, -1, makeSigner(keys[1], committee[1]), committee, nil).MustVerify(stubVerifier)
		fd.msgStore.Save(initProposal)

		// simulate quorum preVotes at r: 0 for v1.
		for i := 0; i < len(committee); i++ {
			preVote := message.NewPrevote(0, height, initProposal.Value(), makeSigner(keys[i], committee[i])).MustVerify(stubVerifier)
			fd.msgStore.Save(preVote)
		}

		// simulate a preCommit at r: 0 for v1 for the node who is going to be addressed as
		// malicious on rule PO for proposing an old value which was not locked at all.
		preCommit := message.NewPrecommit(0, height, initProposal.Value(), signer).MustVerify(stubVerifier)
		fd.msgStore.Save(preCommit)

		// simulate malicious proposal at r: 1, with v2 which was not locked at all.
		// simulate an init proposal at r: 0, with v1.
		maliciousProposal := newProposalMessage(height, 1, 0, signer, committee, nil).MustVerify(stubVerifier)
		fd.msgStore.Save(maliciousProposal)

		// run rule engine.
		onChainProofs := fd.runRulesOverHeight(height, quorum)
		assert.Equal(t, 1, len(onChainProofs))
		assert.Equal(t, autonity.Misbehaviour, onChainProofs[0].Type)
		assert.Equal(t, autonity.PO, onChainProofs[0].Rule)
		assert.Equal(t, message.NewLightProposal(maliciousProposal).Signature(), onChainProofs[0].Message.Signature())
		assert.Equal(t, preCommit.Signature(), onChainProofs[0].Evidences[0].Signature())
	})

	t.Run("RunRule address the misbehaviour of PO rule, the valid round proposed is not correct", func(t *testing.T) {
		// ------------Old Proposal------------
		// PO: (Mr′<r,PV) ∧ (Mr′,PC|pi) ∧ (Mr′<r′′<r,P C|pi)∗ <--- (Mr,P|pi)
		// PO1: [#(Mr′,PV|V) ≥ 2f+ 1] ∧ [nil ∨ V ∨ ⊥] ∧ [nil ∨ ⊥] <--- [V]

		// To address below scenario:
		// Is there a precommit for anything other than nil from the proposer
		// between the valid round and the round of the proposal? If there is
		// then that implies the proposer saw 2f+1 prevotes in that round and
		// hence it should have set that round as the valid round.

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chainMock := NewMockChainContext(ctrl)
		var blockSub event.Subscription
		chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
		chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})
		accountability, _ := autonity.NewAccountability(proposer, backends.NewSimulatedBackend(ccore.GenesisAlloc{proposer: {Balance: big.NewInt(params.Ether)}}, 10000000))

		fd := NewFaultDetector(chainMock, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}), core.NewMsgStore(), nil, nil, proposerKey, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())
		quorum := bft.Quorum(totalPower)
		signerBis := makeSigner(keys[1], committee[1])
		maliciousSigner := makeSigner(keys[2], committee[2])

		header := newBlockHeader(height, committee)
		block := types.NewBlockWithHeader(header)

		// simulate an init proposal at r: 0, with v.
		initProposal := newProposalMessage(height, 0, -1, signer, committee, block).MustVerify(stubVerifier)
		fd.msgStore.Save(initProposal)

		// simulate quorum preVotes for init proposal
		for i := 0; i < len(committee); i++ {
			preVote := message.NewPrevote(0, height, initProposal.Value(), makeSigner(keys[i], committee[i])).MustVerify(stubVerifier)
			fd.msgStore.Save(preVote)
		}

		// simulate a preCommit for init proposal of proposer1, now valid round == 0.
		preCommit1 := message.NewPrecommit(0, height, initProposal.Value(), signerBis).MustVerify(stubVerifier)
		fd.msgStore.Save(preCommit1)

		// assume round changes happens, now proposer1 propose old value with vr: 0.
		// simulate a new proposal at r: 3, with v.
		proposal1 := newProposalMessage(height, 3, 0, signerBis, committee, block).MustVerify(stubVerifier)
		fd.msgStore.Save(proposal1)

		// now quorum preVotes for proposal1, it makes valid round change to 3.
		for i := 0; i < len(committee); i++ {
			preVote := message.NewPrevote(3, height, initProposal.Value(), makeSigner(keys[i], committee[i])).MustVerify(stubVerifier)
			fd.msgStore.Save(preVote)
		}

		// the malicious proposer did preCommit on this round, make its valid round == 3
		preCommit := message.NewPrecommit(3, height, initProposal.Value(), maliciousSigner).MustVerify(stubVerifier)
		fd.msgStore.Save(preCommit)

		// malicious proposer propose at r: 5, with v and a vr: 0 which was not correct anymore.
		maliciousProposal := newProposalMessage(height, 5, 0, maliciousSigner, committee, block).MustVerify(stubVerifier)
		fd.msgStore.Save(maliciousProposal)

		// run rule engine.
		onChainProofs := fd.runRulesOverHeight(height, quorum)
		assert.Equal(t, 1, len(onChainProofs))
		assert.Equal(t, autonity.Misbehaviour, onChainProofs[0].Type)
		assert.Equal(t, autonity.PO, onChainProofs[0].Rule)
		assert.Equal(t, message.NewLightProposal(maliciousProposal).Signature(), onChainProofs[0].Message.Signature())
		assert.Equal(t, preCommit.Signature(), onChainProofs[0].Evidences[0].Signature())
	})

	t.Run("RunRule address the Accusation of PO rule, no quorum preVotes presented on the valid round", func(t *testing.T) {
		// ------------Old Proposal------------
		// PO: (Mr′<r,PV) ∧ (Mr′,PC|pi) ∧ (Mr′<r′′<r,P C|pi)∗ <--- (Mr,P|pi)
		// PO1: [#(Mr′,PV|V) ≥ 2f+ 1] ∧ [nil ∨ V ∨ ⊥] ∧ [nil ∨ ⊥] <--- [V]

		// To address below accusation scenario:
		// If proposer rise an old proposal, then there must be a quorum preVotes on the valid round.
		// Do we see a quorum of preVotes in the valid round, if not we can raise an accusation, since we cannot be sure
		// that these preVotes don't exist

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chainMock := NewMockChainContext(ctrl)
		var blockSub event.Subscription
		chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
		chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})
		accountability, _ := autonity.NewAccountability(proposer, backends.NewSimulatedBackend(ccore.GenesisAlloc{proposer: {Balance: big.NewInt(params.Ether)}}, 10000000))

		fd := NewFaultDetector(chainMock, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}), core.NewMsgStore(), nil, nil, proposerKey, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())
		quorum := bft.Quorum(totalPower)

		header := newBlockHeader(height, committee)
		block := types.NewBlockWithHeader(header)

		// simulate an old proposal at r: 2, with v and vr: 0.
		oldProposal := newProposalMessage(height, 2, 0, signer, committee, block).MustVerify(stubVerifier)
		fd.msgStore.Save(oldProposal)

		// run rule engine.
		onChainProofs := fd.runRulesOverHeight(height, quorum)
		assert.Equal(t, 1, len(onChainProofs))
		assert.Equal(t, autonity.Accusation, onChainProofs[0].Type)
		assert.Equal(t, autonity.PO, onChainProofs[0].Rule)
		assert.Equal(t, message.NewLightProposal(oldProposal).Signature(), onChainProofs[0].Message.Signature())
	})

	t.Run("RunRule address the accusation of PVN, no corresponding proposal of preVote", func(t *testing.T) {
		// PVN: (Mr′<r,PC|pi)∧(Mr′<r′′<r,PC|pi)* ∧ (Mr,P|proposer(r)) <--- (Mr,PV|pi)
		// To address below accusation scenario:
		// If there exist a preVote for a non nil value, then there must be a corresponding proposal at the same round,
		// otherwise an accusation of PVN should rise.
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chainMock := NewMockChainContext(ctrl)
		var blockSub event.Subscription
		chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
		chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})
		accountability, _ := autonity.NewAccountability(proposer, backends.NewSimulatedBackend(ccore.GenesisAlloc{proposer: {Balance: big.NewInt(params.Ether)}}, 10000000))

		fd := NewFaultDetector(chainMock, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}), core.NewMsgStore(), nil, nil, proposerKey, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())
		quorum := bft.Quorum(totalPower)

		// simulate a preVote for v at round, let's make the corresponding proposal missing.
		preVote := message.NewPrevote(round, height, noneNilValue, makeSigner(keys[1], committee[1])).MustVerify(stubVerifier)
		fd.msgStore.Save(preVote)

		// run rule engine.
		onChainProofs := fd.runRulesOverHeight(height, quorum)
		assert.Equal(t, 1, len(onChainProofs))
		assert.Equal(t, autonity.Accusation, onChainProofs[0].Type)
		assert.Equal(t, autonity.PVN, onChainProofs[0].Rule)
		assert.Equal(t, preVote.Signature(), onChainProofs[0].Message.Signature())
	})

	t.Run("RunRule address the misbehaviour of PVN, node preVote for value V1 while it preCommitted another value at previous round", func(t *testing.T) {
		//t.Skip("skip this case from CI jobs, it works in local environment.")
		// PVN: (Mr′<r,PC|pi)∧(Mr′<r′′<r,PC|pi)* ∧ (Mr,P|proposer(r)) <--- (Mr,PV|pi)
		// PVN2: [nil ∨ ⊥] ∧ [nil ∨ ⊥] ∧ [V:Valid(V)] <--- [V]: r′= 0,∀r′′< r:Mr′′,PC|pi=nil
		// PVN2, If there is a valid proposal V at round r, and pi never
		// ever precommit(locked a value) before, then pi should prevote
		// for V or a nil in case of timeout at this round.

		// To address below misbehaviour scenario:
		// Node preCommitted at v1 at R_x, while it preVote for v2 at R_x + n.
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chainMock := NewMockChainContext(ctrl)
		var blockSub event.Subscription
		chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
		chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})
		accountability, _ := autonity.NewAccountability(proposer, backends.NewSimulatedBackend(ccore.GenesisAlloc{proposer: {Balance: big.NewInt(params.Ether)}}, 10000000))

		fd := NewFaultDetector(chainMock, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}), core.NewMsgStore(), nil, nil, proposerKey, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())
		quorum := bft.Quorum(totalPower)

		newProposer := makeSigner(keys[2], committee[2])

		header := newBlockHeader(height, committee)
		block := types.NewBlockWithHeader(header)

		// simulate an init proposal at r: 0, with v.
		initProposal := newProposalMessage(height, 0, -1, signer, committee, block).MustVerify(stubVerifier)
		fd.msgStore.Save(initProposal)

		// simulate quorum preVotes for init proposal
		for i := 0; i < len(committee); i++ {
			preVote := message.NewPrevote(0, height, initProposal.Value(), makeSigner(keys[i], committee[i])).MustVerify(stubVerifier)
			fd.msgStore.Save(preVote)
		}

		// the malicious node did preCommit for v1 on round 0
		preCommit := message.NewPrecommit(0, height, initProposal.Value(), maliciousSigner).MustVerify(stubVerifier)
		fd.msgStore.Save(preCommit)

		// the malicious node did preCommit for nil at round 1, and round 2 to fill the round gaps.
		for i := 1; i < 3; i++ {
			pc := message.NewPrecommit(int64(i), height, nilValue, maliciousSigner).MustVerify(stubVerifier)
			fd.msgStore.Save(pc)
		}

		// assume round changes, someone propose V2 at round 3, and malicious Node now it preVote for this V2.
		newProposal := newProposalMessage(height, 3, -1, newProposer, committee, nil).MustVerify(stubVerifier)
		fd.msgStore.Save(newProposal)

		// now the malicious node preVote for a new value V2 at higher round 3.
		preVote := message.NewPrevote(3, height, newProposal.Value(), maliciousSigner).MustVerify(stubVerifier)
		fd.msgStore.Save(preVote)

		onChainProofs := fd.runRulesOverHeight(height, quorum)

		assert.Equal(t, 1, len(onChainProofs))
		assert.Equal(t, autonity.Misbehaviour, onChainProofs[0].Type)
		assert.Equal(t, autonity.PVN, onChainProofs[0].Rule)
		assert.Equal(t, 4, len(onChainProofs[0].Evidences))
		assert.Equal(t, preVote.Signature(), onChainProofs[0].Message.Signature())
		assert.Equal(t, message.NewLightProposal(newProposal).Signature(), onChainProofs[0].Evidences[0].Signature())
		assert.Equal(t, preCommit.Signature(), onChainProofs[0].Evidences[1].Signature())
	})

	t.Run("RunRule address the misbehaviour of PVN, with gaps of preCommits, the PVN is not provable", func(t *testing.T) {
		// PVN: (Mr′<r,PC|pi)∧(Mr′<r′′<r,PC|pi)* ∧ (Mr,P|proposer(r)) <--- (Mr,PV|pi)
		// PVN2: [nil ∨ ⊥] ∧ [nil ∨ ⊥] ∧ [V:Valid(V)] <--- [V]: r′= 0,∀r′′< r:Mr′′,PC|pi=nil
		// PVN2, If there is a valid proposal V at round r, and pi never
		// ever precommit(locked a value) before, then pi should prevote
		// for V or a nil in case of timeout at this round.

		// To address below misbehaviour scenario:
		// Node preCommitted at v1 at R_x, while it preVote for v2 at R_x + n.
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chainMock := NewMockChainContext(ctrl)
		var blockSub event.Subscription
		chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
		chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})
		accountability, _ := autonity.NewAccountability(proposer, backends.NewSimulatedBackend(ccore.GenesisAlloc{proposer: {Balance: big.NewInt(params.Ether)}}, 10000000))

		fd := NewFaultDetector(chainMock, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}), core.NewMsgStore(), nil, nil, proposerKey, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())
		quorum := bft.Quorum(totalPower)

		newProposer := makeSigner(keys[2], committee[2])

		header := newBlockHeader(height, committee)
		block := types.NewBlockWithHeader(header)

		// simulate an init proposal at r: 0, with v.
		initProposal := newProposalMessage(height, 0, -1, signer, committee, block).MustVerify(stubVerifier)
		fd.msgStore.Save(initProposal)

		// simulate quorum preVotes for init proposal
		for i := 0; i < len(committee); i++ {
			preVote := message.NewPrevote(0, height, initProposal.Value(), makeSigner(keys[i], committee[i])).MustVerify(stubVerifier)
			fd.msgStore.Save(preVote)
		}

		// the malicious node did preCommit for v1 on round 0
		preCommit := message.NewPrecommit(0, height, initProposal.Value(), maliciousSigner).MustVerify(stubVerifier)
		fd.msgStore.Save(preCommit)

		// the malicious node did preCommit for nil at round 1, let no preCommit at round 2 to form the gap.
		preCommitR1 := message.NewPrecommit(int64(1), height, nilValue, maliciousSigner).MustVerify(stubVerifier)
		fd.msgStore.Save(preCommitR1)

		// assume round changes, someone propose V2 at round 3, and malicious Node now it preVote for this V2.
		newProposal := newProposalMessage(height, 3, -1, newProposer, committee, nil).MustVerify(stubVerifier)
		fd.msgStore.Save(newProposal)

		// now the malicious node preVote for a new value V2 at higher round 3.
		preVote := message.NewPrevote(3, height, newProposal.Value(), maliciousSigner).MustVerify(stubVerifier)
		fd.msgStore.Save(preVote)

		onChainProofs := fd.runRulesOverHeight(height, quorum)

		assert.Equal(t, 0, len(onChainProofs))
	})

	t.Run("RunRule to address Accusation of rule PVO, no quorum preVotes for valid round", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chainMock := NewMockChainContext(ctrl)
		var blockSub event.Subscription
		chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
		chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})
		accountability, _ := autonity.NewAccountability(proposer, backends.NewSimulatedBackend(ccore.GenesisAlloc{proposer: {Balance: big.NewInt(params.Ether)}}, 10000000))

		fd := NewFaultDetector(chainMock, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}), core.NewMsgStore(), nil, nil, proposerKey, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())
		quorum := bft.Quorum(totalPower)

		header := newBlockHeader(height, committee)
		block := types.NewBlockWithHeader(header)

		// simulate a proposal at r: 3, and vr: 1, with v.
		oldProposal := newProposalMessage(height, 3, 1, signer, committee, block).MustVerify(stubVerifier)
		fd.msgStore.Save(oldProposal)

		// simulate a preVote at r: 3 for value v.
		preVote := message.NewPrevote(3, height, oldProposal.Value(), maliciousSigner).MustVerify(stubVerifier)
		fd.msgStore.Save(preVote)

		onChainProofs := fd.runRulesOverHeight(height, quorum)

		assert.Equal(t, 2, len(onChainProofs))
		assert.Equal(t, autonity.Accusation, onChainProofs[0].Type)
		assert.Equal(t, autonity.PO, onChainProofs[0].Rule)
		assert.Equal(t, message.NewLightProposal(oldProposal).Signature(), onChainProofs[0].Message.Signature())

		assert.Equal(t, autonity.Accusation, onChainProofs[1].Type)
		assert.Equal(t, autonity.PVO, onChainProofs[1].Rule)
		assert.Equal(t, preVote.Signature(), onChainProofs[1].Message.Signature())
	})

	t.Run("RunRule to address misbehaviour of rule PVO, there were quorum prevote for not V at valid round", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chainMock := NewMockChainContext(ctrl)
		var blockSub event.Subscription
		chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
		chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})
		accountability, _ := autonity.NewAccountability(proposer, backends.NewSimulatedBackend(ccore.GenesisAlloc{proposer: {Balance: big.NewInt(params.Ether)}}, 10000000))

		fd := NewFaultDetector(chainMock, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}), core.NewMsgStore(), nil, nil, proposerKey, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())
		quorum := bft.Quorum(totalPower)

		header := newBlockHeader(height, committee)
		block := types.NewBlockWithHeader(header)

		// simulate a proposal at r: 3, and vr: 0, with v.
		oldProposal := newProposalMessage(height, 3, 0, signer, committee, block).MustVerify(stubVerifier)
		fd.msgStore.Save(oldProposal)

		// simulate quorum prevotes for not v at vr.
		for i := 0; i < len(committee); i++ {
			preVote := message.NewPrevote(0, height, noneNilValue, makeSigner(keys[i], committee[i])).MustVerify(stubVerifier)
			fd.msgStore.Save(preVote)
		}
		// simulate a preVote at r: 3 for value v, thus it is a misbehaviour.
		preVote := message.NewPrevote(3, height, oldProposal.Value(), maliciousSigner).MustVerify(stubVerifier)
		fd.msgStore.Save(preVote)
		onChainProofs := fd.runRulesOverHeight(height, quorum)
		presentPVO := false
		for _, p := range onChainProofs {
			if p.Type == autonity.Misbehaviour && p.Rule == autonity.PVO {
				presentPVO = true
				assert.Equal(t, message.PrevoteCode, p.Message.Code())
				assert.Equal(t, preVote.Signature(), p.Message.Signature())
			}
		}
		assert.Equal(t, true, presentPVO)
	})

	t.Run("RunRule to address misbehaviour of rule PVO1, node last precommited at a value of not v", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chainMock := NewMockChainContext(ctrl)
		var blockSub event.Subscription
		chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
		chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})
		accountability, _ := autonity.NewAccountability(proposer, backends.NewSimulatedBackend(ccore.GenesisAlloc{proposer: {Balance: big.NewInt(params.Ether)}}, 10000000))

		fd := NewFaultDetector(chainMock, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}), core.NewMsgStore(), nil, nil, proposerKey, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())
		quorum := bft.Quorum(totalPower)

		header := newBlockHeader(height, committee)
		block := types.NewBlockWithHeader(header)

		// simulate a proposal at r: 3, and vr: 0, with v.
		oldProposal := newProposalMessage(height, 3, 0, signer, committee, block).MustVerify(stubVerifier)
		fd.msgStore.Save(oldProposal)

		// simulate quorum prevotes for v at vr.
		for i := 0; i < len(committee); i++ {
			preVote := message.NewPrevote(0, height, oldProposal.Value(), makeSigner(keys[i], committee[i])).MustVerify(stubVerifier)
			fd.msgStore.Save(preVote)
		}

		// simulate a precommit at r: 0 with value v.
		pcValidRound := message.NewPrecommit(0, height, oldProposal.Value(), maliciousSigner).MustVerify(stubVerifier)
		fd.msgStore.Save(pcValidRound)

		// simulate a precommit at r: 1 with value v.
		preCommitForV := message.NewPrecommit(1, height, oldProposal.Value(), maliciousSigner).MustVerify(stubVerifier)
		fd.msgStore.Save(preCommitForV)

		// simulate a precommit at r: 2 with value not v.
		preCommitForNotV := message.NewPrecommit(2, height, noneNilValue, maliciousSigner).MustVerify(stubVerifier)
		fd.msgStore.Save(preCommitForNotV)

		// simulate a preVote at r: 3 for value v.
		preVote := message.NewPrevote(3, height, oldProposal.Value(), maliciousSigner).MustVerify(stubVerifier)
		fd.msgStore.Save(preVote)

		onChainProofs := fd.runRulesOverHeight(height, quorum)
		presentPVO1 := false
		for _, p := range onChainProofs {
			if p.Type == autonity.Misbehaviour && p.Rule == autonity.PVO12 {
				presentPVO1 = true
				assert.Equal(t, message.PrevoteCode, p.Message.Code())
				assert.Equal(t, preVote.Signature(), p.Message.Signature())
			}
		}
		assert.Equal(t, true, presentPVO1)
	})

	t.Run("RunRule to address misbehaviour of rule PVO2, node did precommited at a value of not v between valid "+
		"round and current round", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chainMock := NewMockChainContext(ctrl)
		var blockSub event.Subscription
		chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
		chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})
		accountability, _ := autonity.NewAccountability(proposer, backends.NewSimulatedBackend(ccore.GenesisAlloc{proposer: {Balance: big.NewInt(params.Ether)}}, 10000000))

		fd := NewFaultDetector(chainMock, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}), core.NewMsgStore(), nil, nil, proposerKey, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())
		quorum := bft.Quorum(totalPower)
		header := newBlockHeader(height, committee)
		block := types.NewBlockWithHeader(header)

		// simulate a proposal at r: 3, and vr: 0, with v.
		oldProposal := newProposalMessage(height, 3, 0, signer, committee, block).MustVerify(stubVerifier)
		fd.msgStore.Save(oldProposal)

		// simulate quorum prevotes for v at vr.
		for i := 0; i < len(committee); i++ {
			preVote := message.NewPrevote(0, height, oldProposal.Value(), makeSigner(keys[i], committee[i])).MustVerify(stubVerifier)
			fd.msgStore.Save(preVote)
		}

		// simulate a precommit at r: 0 with value not v.
		pcValidRound := message.NewPrecommit(0, height, noneNilValue, maliciousSigner).MustVerify(stubVerifier)
		fd.msgStore.Save(pcValidRound)

		// simulate a precommit at r: 1 with value not v.
		preCommitForV := message.NewPrecommit(1, height, noneNilValue, maliciousSigner).MustVerify(stubVerifier)
		fd.msgStore.Save(preCommitForV)

		// simulate a precommit at r: 2 with value not v.
		preCommitForNotV := message.NewPrecommit(2, height, noneNilValue, maliciousSigner).MustVerify(stubVerifier)
		fd.msgStore.Save(preCommitForNotV)

		// simulate a preVote at r: 3 for value v.
		preVote := message.NewPrevote(3, height, oldProposal.Value(), maliciousSigner).MustVerify(stubVerifier)
		fd.msgStore.Save(preVote)

		onChainProofs := fd.runRulesOverHeight(height, quorum)
		presentPVO1 := false
		for _, p := range onChainProofs {
			if p.Type == autonity.Misbehaviour && p.Rule == autonity.PVO12 {
				presentPVO1 = true
				assert.Equal(t, message.PrevoteCode, p.Message.Code())
				assert.Equal(t, preVote.Signature(), p.Message.Signature())
			}
		}
		assert.Equal(t, true, presentPVO1)
	})

	t.Run("RunRule address Accusation of rule C1, no corresponding quorum prevotes for a preCommit msg", func(t *testing.T) {
		// C: [Mr,P|proposer(r)] ∧ [Mr,PV] <--- [Mr,PC|pi]
		// C1: [V:Valid(V)] ∧ [#(V) ≥ 2f+ 1] <--- [V]

		// To address below accusation scenario:
		// Node preCommit for a V at round R, but we cannot see the corresponding quorum preVotes that propose the value
		// at the same round of that preCommit msg.

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chainMock := NewMockChainContext(ctrl)
		var blockSub event.Subscription
		chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
		chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})
		accountability, _ := autonity.NewAccountability(proposer, backends.NewSimulatedBackend(ccore.GenesisAlloc{proposer: {Balance: big.NewInt(params.Ether)}}, 10000000))

		fd := NewFaultDetector(chainMock, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}), core.NewMsgStore(), nil, nil, proposerKey, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())
		quorum := bft.Quorum(totalPower)

		preCommit := message.NewPrecommit(0, height, noneNilValue, signer).MustVerify(stubVerifier)
		fd.msgStore.Save(preCommit)

		onChainProofs := fd.runRulesOverHeight(height, quorum)
		assert.Equal(t, 1, len(onChainProofs))
		assert.Equal(t, autonity.Accusation, onChainProofs[0].Type)
		assert.Equal(t, autonity.C1, onChainProofs[0].Rule)
		assert.Equal(t, preCommit.Signature(), onChainProofs[0].Message.Signature())
	})

	t.Run("RunRule address accusation of rule C1, no present of quorum preVotes of V to justify the preCommit of V", func(t *testing.T) {
		// ------------precommits------------
		// C: [Mr,P|proposer(r)] ∧ [Mr,PV] <--- [Mr,PC|pi]
		// C1: [V:Valid(V)] ∧ [#(V) ≥ 2f+ 1] <--- [V]

		// To address below accusation scenario:
		// Node preCommit for a value V, but observer haven't seen quorum preVotes for V at the round, an accusation shall
		// rise.
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chainMock := NewMockChainContext(ctrl)
		var blockSub event.Subscription
		chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
		chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})
		accountability, _ := autonity.NewAccountability(proposer, backends.NewSimulatedBackend(ccore.GenesisAlloc{proposer: {Balance: big.NewInt(params.Ether)}}, 10000000))

		fd := NewFaultDetector(chainMock, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}), core.NewMsgStore(), nil, nil, proposerKey, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())
		quorum := bft.Quorum(totalPower)
		header := newBlockHeader(height, committee)
		block := types.NewBlockWithHeader(header)

		// simulate an init proposal at r: 0, with v.
		initProposal := newProposalMessage(height, 0, -1, signer, committee, block).MustVerify(stubVerifier)
		fd.msgStore.Save(initProposal)

		// malicious node preCommit to v even through there was no quorum preVotes for v.
		preCommit := message.NewPrecommit(0, height, initProposal.Value(), maliciousSigner).MustVerify(stubVerifier)
		fd.msgStore.Save(preCommit)

		onChainProofs := fd.runRulesOverHeight(height, quorum)
		assert.Equal(t, 1, len(onChainProofs))
		assert.Equal(t, autonity.Accusation, onChainProofs[0].Type)
		assert.Equal(t, autonity.C1, onChainProofs[0].Rule)
		assert.Equal(t, preCommit.Signature(), onChainProofs[0].Message.Signature())
	})
}
