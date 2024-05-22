package accountability

import (
	"crypto/ecdsa"
	"github.com/autonity/autonity/consensus/ethash"
	"github.com/autonity/autonity/core/rawdb"
	"github.com/autonity/autonity/core/vm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"math/big"
	"math/rand"
	"testing"

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
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/params/generated"
)

var (
	committee, keys, nodeKeys = generateCommittee()
	cSize                     = len(committee)
	proposerIdx               = 0
	proposer                  = committee[proposerIdx].Address
	proposerKey               = keys[proposerIdx]
	proposerNodeKey           = nodeKeys[proposerIdx]
	signer                    = makeSigner(proposerKey)
	self                      = &committee[proposerIdx]

	remotePeerIdx = 1
	remote        = &committee[remotePeerIdx]
	remotePeer    = committee[remotePeerIdx].Address
	remoteSigner  = makeSigner(keys[remotePeerIdx])
)

func generateCommittee() (types.Committee, []blst.SecretKey, []*ecdsa.PrivateKey) {
	n := 5
	validators := make(types.Committee, n)
	pkeys := make([]*ecdsa.PrivateKey, n)
	var consensusKeys []blst.SecretKey
	for i := 0; i < n; i++ {
		privateKey, _ := crypto.GenerateKey()
		consensusKey, _ := blst.RandKey()
		committeeMember := types.CommitteeMember{
			Address:           crypto.PubkeyToAddress(privateKey.PublicKey),
			VotingPower:       new(big.Int).SetUint64(1),
			ConsensusKey:      consensusKey.PublicKey(),
			ConsensusKeyBytes: consensusKey.PublicKey().Marshal(),
			Index:             uint64(i),
		}
		validators[i] = committeeMember
		pkeys[i] = privateKey
		consensusKeys = append(consensusKeys, consensusKey)
	}
	return validators, consensusKeys, pkeys
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
func newValidatedProposalMessage(h uint64, r int64, vr int64, signer message.Signer, committee types.Committee, withValue *types.Block, idx int) *message.Propose {
	block := withValue
	if withValue == nil {
		header := newBlockHeader(h, committee)
		block = types.NewBlockWithHeader(header)
	}
	p := message.NewPropose(r, h, vr, block, signer, &committee[idx])
	lastHeader := newBlockHeader(h-1, committee)
	if err := p.PreValidate(lastHeader); err != nil {
		panic(err)
	}
	if err := p.Validate(); err != nil {
		panic(err)
	}
	return p
}

func TestSameVote(t *testing.T) {
	height := uint64(100)
	r1 := int64(0)
	r2 := int64(1)
	validRound := int64(1)
	proposal := newValidatedProposalMessage(height, r1, validRound, signer, committee, nil, proposerIdx)
	proposal2 := newValidatedProposalMessage(height, r2, validRound, signer, committee, nil, proposerIdx)
	require.Equal(t, false, proposal.Hash() == proposal2.Hash())
}

func TestSubmitMisbehaviour(t *testing.T) {
	height := uint64(100)
	round := int64(0)
	lastHeader := newBlockHeader(height-1, committee)
	// submit a equivocation proofs.
	proposal := newValidatedLightProposal(t, height, round, -1, signer, committee, lastHeader, nil, proposerIdx)
	proposal2 := newValidatedLightProposal(t, height, round, -1, signer, committee, lastHeader, nil, proposerIdx)
	var proofs []message.Msg
	proofs = append(proofs, proposal2)

	fd := &FaultDetector{
		misbehaviourProofCh: make(chan *autonity.AccountabilityEvent, 100),
		logger:              log.New("FaultDetector", nil),
	}

	fd.submitMisbehavior(proposal, proofs, errEquivocation, proposer, proposerIdx)
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

		fd := NewFaultDetector(chainMock, fdAddr, nil, core.NewMsgStore(), nil, nil, proposerNodeKey, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())
		// store a msg before check point height in case of node is start from reset.
		msgBeforeCheckPointHeight := newValidatedProposalMessage(checkPointHeight-1, 0, -1, makeSigner(keys[1]), committee, nil, 1)
		fd.msgStore.Save(msgBeforeCheckPointHeight, committee)

		// simulate there was a maliciousProposal at init round 0, and save to msg store.
		initProposal := newValidatedProposalMessage(checkPointHeight, 0, -1, makeSigner(keys[1]), committee, nil, 1)
		fd.msgStore.Save(initProposal, committee)
		// simulate there were quorum preVotes for initProposal at init round 0, and save them.
		for i := 0; i < len(committee); i++ {
			preVote := newValidatedPrevote(t, 0, checkPointHeight, initProposal.Value(), makeSigner(keys[i]), &committee[i], cSize, lastHeader)
			fd.msgStore.Save(preVote, committee)
		}

		// Node preCommit for init Proposal at init round 0 since there were quorum preVotes for it, and save it.
		preCommit := newValidatedPrecommit(t, 0, checkPointHeight, initProposal.Value(), signer, self, cSize, lastHeader)
		fd.msgStore.Save(preCommit, committee)

		// While Node propose a new malicious Proposal at new round with VR as -1 which is malicious, should be addressed by rule PN.
		maliciousProposal := newValidatedProposalMessage(checkPointHeight, round, -1, signer, committee, nil, proposerIdx)
		fd.msgStore.Save(maliciousProposal, committee)

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

func TestGenerateOnChainProof(t *testing.T) {
	height := uint64(100)
	round := int64(3)
	lastHeader := newBlockHeader(height-1, committee)

	proposal := newValidatedLightProposal(t, height, round, -1, signer, committee, lastHeader, nil, proposerIdx)
	equivocatedProposal := newValidatedLightProposal(t, height, round, -1, signer, committee, lastHeader, nil, proposerIdx)
	var evidence []message.Msg
	evidence = append(evidence, equivocatedProposal)

	p := Proof{
		Offender:      proposer,
		OffenderIndex: proposerIdx,
		Type:          autonity.Misbehaviour,
		Rule:          autonity.Equivocation,
		Message:       proposal,
		Evidences:     evidence,
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
//
//	in such context, the accusation is considered as useless, it should be dropped.
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

		fd := NewFaultDetector(chainMock, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}), core.NewMsgStore(), nil, nil, proposerNodeKey, &autonity.ProtocolContracts{Accountability: bindings}, log.Root())
		// simulate a proposal message with an old value and a valid round.
		proposal := newValidatedProposalMessage(height, round, validRound, signer, committee, nil, proposerIdx)
		fd.msgStore.Save(proposal, committee)

		// simulate at least quorum num of preVotes for a value at a validRound.
		for i := 0; i < len(committee); i++ {
			preVote := newValidatedPrevote(t, validRound, height, proposal.Value(), makeSigner(keys[i]), &committee[i], cSize, lastHeader)
			fd.msgStore.Save(preVote, committee)
		}

		var accusation = Proof{
			Offender:      proposer,
			OffenderIndex: proposerIdx,
			Type:          autonity.Accusation,
			Rule:          autonity.PO,
			Message:       proposal.ToLight(),
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
		fd := NewFaultDetector(chainMock, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}), core.NewMsgStore(), nil, nil, proposerNodeKey, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())
		// simulate a proposal message with an old value and a valid round.
		proposal := newValidatedProposalMessage(height, round, validRound, signer, committee, nil, proposerIdx)
		fd.msgStore.Save(proposal, committee)

		// simulate less than quorum num of preVotes for a value at a validRound.
		preVote := newValidatedPrevote(t, validRound, height, proposal.Value(), signer, self, cSize, lastHeader)
		fd.msgStore.Save(preVote, committee)

		var accusation = Proof{
			Type:          autonity.Accusation,
			Rule:          autonity.PO,
			Offender:      proposer,
			OffenderIndex: proposerIdx,
			Message:       message.NewLightProposal(proposal),
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
		proposal := newValidatedProposalMessage(height, round, -1, signer, committee, nil, proposerIdx)
		fd.msgStore.Save(proposal, committee)

		// simulate at least quorum num of preVotes for a value at a validRound.
		for i := 0; i < len(committee); i++ {
			pv := newValidatedPrevote(t, round, height, proposal.Value(), makeSigner(keys[i]), &committee[i], cSize, lastHeader)
			fd.msgStore.Save(pv, committee)
		}

		preVote := newValidatedPrevote(t, round, height, proposal.Value(), signer, self, cSize, lastHeader)

		var accusation = Proof{
			OffenderIndex: proposerIdx,
			Offender:      proposer,
			Type:          autonity.Accusation,
			Rule:          autonity.PVN,
			Message:       preVote,
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

		preVote := newValidatedPrevote(t, round, height, noneNilValue, signer, self, cSize, lastHeader)
		fd.msgStore.Save(preVote, committee)

		var accusation = Proof{
			Offender:      proposer,
			OffenderIndex: proposerIdx,
			Type:          autonity.Accusation,
			Rule:          autonity.PVN,
			Message:       preVote,
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

		fd := NewFaultDetector(chainMock, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}), core.NewMsgStore(), nil, nil, proposerNodeKey, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())

		var p Proof
		p.Rule = autonity.PVO
		p.Offender = proposer
		p.OffenderIndex = proposerIdx
		oldProposal := newValidatedProposalMessage(height, 1, 0, signer, committee, nil, proposerIdx)
		preVote := newValidatedPrevote(t, 1, height, oldProposal.Value(), signer, self, cSize, lastHeader)
		p.Message = preVote
		p.Evidences = append(p.Evidences, oldProposal.ToLight())

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

		fd := NewFaultDetector(chainMock, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}), core.NewMsgStore(), nil, nil, proposerNodeKey, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())
		var p Proof
		p.Rule = autonity.PVO
		p.OffenderIndex = proposerIdx
		p.Offender = proposer
		validRound := int64(0)
		oldProposal := newValidatedProposalMessage(height, 1, validRound, signer, committee, nil, proposerIdx)
		preVote := newValidatedPrevote(t, 1, height, oldProposal.Value(), signer, self, cSize, lastHeader)
		p.Message = preVote
		p.Evidences = append(p.Evidences, message.NewLightProposal(oldProposal))

		// prepare quorum preVotes at msg store.
		for i := 0; i < len(committee); i++ {
			pv := newValidatedPrevote(t, validRound, height, oldProposal.Value(), makeSigner(keys[i]), &committee[i], cSize, lastHeader)
			fd.msgStore.Save(pv, committee)
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
		fd := NewFaultDetector(chainMock, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}), core.NewMsgStore(), nil, nil, proposerNodeKey, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())

		// simulate at least quorum num of preVotes for a value at a validRound.
		for i := 0; i < len(committee); i++ {
			preVote := newValidatedPrevote(t, round, height, noneNilValue, makeSigner(keys[i]), &committee[i], cSize, lastHeader)
			fd.msgStore.Save(preVote, committee)
		}

		preCommit := newValidatedPrecommit(t, round, height, noneNilValue, signer, self, cSize, lastHeader)
		fd.msgStore.Save(preCommit, committee)

		var accusation = Proof{
			Offender:      proposer,
			OffenderIndex: proposerIdx,
			Type:          autonity.Accusation,
			Rule:          autonity.C1,
			Message:       preCommit,
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

		fd := NewFaultDetector(chainMock, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}), core.NewMsgStore(), nil, nil, proposerNodeKey, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())

		preCommit := newValidatedPrecommit(t, round, height, noneNilValue, signer, self, cSize, lastHeader)
		fd.msgStore.Save(preCommit, committee)

		var accusation = Proof{
			OffenderIndex: proposerIdx,
			Offender:      proposer,
			Type:          autonity.Accusation,
			Rule:          autonity.C1,
			Message:       preCommit,
		}

		_, err := fd.innocenceProofC1(&accusation)
		assert.Equal(t, errNoEvidenceForC1, err)
	})

	t.Run("Test error to rule mapping", func(t *testing.T) {
		rule := errorToRule(errEquivocation)
		assert.Equal(t, autonity.Equivocation, rule)

		rule = errorToRule(errProposer)
		assert.Equal(t, autonity.InvalidProposer, rule)
	})
}

// Please refer to the rules in the rule engine for each step of tendermint to understand the test context.
// TestNewProposalAccountabilityCheck, it tests the accountability events over a new proposal sent by a proposer.
func TestNewProposalAccountabilityCheck(t *testing.T) {
	height := uint64(0)
	lastHeader := newBlockHeader(height-1, committee)
	newProposal0 := newValidatedProposalMessage(height, 3, -1, signer, committee, nil, proposerIdx)
	nonNilPrecommit0 := newValidatedPrecommit(t, 1, height, common.BytesToHash([]byte("test")), signer, self, cSize, lastHeader)
	nilPrecommit0 := newValidatedPrecommit(t, 1, height, common.Hash{}, signer, self, cSize, lastHeader)

	newProposal1 := newValidatedProposalMessage(height, 5, -1, signer, committee, nil, proposerIdx)
	nilPrecommit1 := newValidatedPrecommit(t, 3, height, common.Hash{}, signer, self, cSize, lastHeader)

	newProposal0E := newValidatedProposalMessage(height, 3, 1, signer, committee, nil, proposerIdx)

	t.Run("misbehaviour when pi has sent a non-nil precommit in a previous round", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposal0, committee)
		fd.msgStore.Save(nonNilPrecommit0, committee)

		expectedProof := &Proof{
			Offender:      proposer,
			OffenderIndex: proposerIdx,
			Type:          autonity.Misbehaviour,
			Rule:          autonity.PN,
			Evidences:     []message.Msg{nonNilPrecommit0},
			Message:       message.NewLightProposal(newProposal0),
		}

		proofs := fd.newProposalsAccountabilityCheck(0)
		require.Equal(t, 1, len(proofs))
		actualProof := proofs[0]
		require.Equal(t, expectedProof, actualProof)
	})

	t.Run("no proof is returned when proposal is equivocated", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposal0, committee)
		fd.msgStore.Save(nonNilPrecommit0, committee)
		fd.msgStore.Save(newProposal0E, committee)

		proofs := fd.newProposalsAccountabilityCheck(0)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when pi proposes a new proposal and no precommit has been sent", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposal0, committee)
		fd.msgStore.Save(newProposal1, committee)

		proofs := fd.newProposalsAccountabilityCheck(0)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when pi proposes a new proposal and has sent nil precommits in previous rounds", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposal0, committee)
		fd.msgStore.Save(nilPrecommit0, committee)
		fd.msgStore.Save(newProposal1, committee)
		fd.msgStore.Save(nilPrecommit1, committee)

		proofs := fd.newProposalsAccountabilityCheck(0)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("multiple proof of misbehaviours when pi has sent non-nil precommits in previous rounds for multiple proposals", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposal0, committee)
		fd.msgStore.Save(nonNilPrecommit0, committee)
		fd.msgStore.Save(newProposal1, committee)

		expectedProof0 := &Proof{
			OffenderIndex: proposerIdx,
			Offender:      proposer,
			Type:          autonity.Misbehaviour,
			Rule:          autonity.PN,
			Evidences:     []message.Msg{nonNilPrecommit0},
			Message:       newProposal0,
		}

		expectedProof1 := &Proof{
			Offender:      proposer,
			OffenderIndex: proposerIdx,
			Type:          autonity.Misbehaviour,
			Rule:          autonity.PN,
			Evidences:     []message.Msg{nonNilPrecommit0},
			Message:       newProposal1,
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

	lastHeader := newBlockHeader(height-1, committee)
	header := newBlockHeader(height, committee)
	block := types.NewBlockWithHeader(header)
	header1 := newBlockHeader(height, committee)
	block1 := types.NewBlockWithHeader(header1)

	oldProposal0 := newValidatedProposalMessage(height, 3, 0, signer, committee, block, proposerIdx)
	oldProposal5 := newValidatedProposalMessage(height, 5, 2, signer, committee, block, proposerIdx)
	oldProposal0E := newValidatedProposalMessage(height, 3, 2, signer, committee, block1, proposerIdx)
	oldProposal0E2 := newValidatedProposalMessage(height, 3, 0, signer, committee, block1, proposerIdx)

	nonNilPrecommit0V := newValidatedPrecommit(t, 0, height, block.Hash(), signer, self, cSize, lastHeader)
	nonNilPrecommit0VPrime := newValidatedPrecommit(t, 0, height, block1.Hash(), signer, self, cSize, lastHeader)
	nonNilPrecommit2VPrime := newValidatedPrecommit(t, 2, height, block1.Hash(), signer, self, cSize, lastHeader)
	nonNilPrecommit1 := newValidatedPrecommit(t, 1, height, block.Hash(), signer, self, cSize, lastHeader)

	nilPrecommit0 := newValidatedPrecommit(t, 0, height, nilValue, signer, self, cSize, lastHeader)

	var quorumPrevotes0VPrime []message.Msg
	for i := int64(0); i < quorum.Int64(); i++ {
		prevote := newValidatedPrevote(t, 0, height, block1.Hash(), makeSigner(keys[i]), &committee[i], cSize, lastHeader)
		quorumPrevotes0VPrime = append(quorumPrevotes0VPrime, prevote)
	}

	var quorumPrevotes0V []message.Msg
	for i := int64(0); i < quorum.Int64(); i++ {
		prevote := newValidatedPrevote(t, 0, height, block.Hash(), makeSigner(keys[i]), &committee[i], cSize, lastHeader)
		quorumPrevotes0V = append(quorumPrevotes0V, prevote)
	}

	var precommiteNilAfterVR []message.Msg
	for i := 1; i < 3; i++ {
		precommit := newValidatedPrecommit(t, int64(i), height, nilValue, signer, self, cSize, lastHeader)
		precommiteNilAfterVR = append(precommiteNilAfterVR, precommit)
	}

	t.Run("misbehaviour when pi precommited for a different value in valid round than in the old proposal", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0, committee)
		fd.msgStore.Save(nonNilPrecommit0VPrime, committee)

		expectedProof := &Proof{
			Offender:      proposer,
			OffenderIndex: proposerIdx,
			Type:          autonity.Misbehaviour,
			Rule:          autonity.PO,
			Evidences:     []message.Msg{nonNilPrecommit0VPrime},
			Message:       message.NewLightProposal(oldProposal0),
		}

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 1, len(proofs))
		actualProof := proofs[0]
		require.Equal(t, expectedProof, actualProof)
	})

	t.Run("misbehaviour when pi incorrectly set the valid round with a different value than the proposal", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0, committee)
		fd.msgStore.Save(nonNilPrecommit2VPrime, committee)

		expectedProof := &Proof{
			Offender:      proposer,
			OffenderIndex: proposerIdx,
			Type:          autonity.Misbehaviour,
			Rule:          autonity.PO,
			Evidences:     []message.Msg{nonNilPrecommit2VPrime},
			Message:       message.NewLightProposal(oldProposal0),
		}

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 1, len(proofs))
		actualProof := proofs[0]
		require.Equal(t, expectedProof, actualProof)
	})

	t.Run("misbehaviour when pi incorrectly set the valid round with the same value as the proposal", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0, committee)
		fd.msgStore.Save(nonNilPrecommit1, committee)

		expectedProof := &Proof{
			OffenderIndex: proposerIdx,
			Offender:      proposer,
			Type:          autonity.Misbehaviour,
			Rule:          autonity.PO,
			Evidences:     []message.Msg{nonNilPrecommit1},
			Message:       message.NewLightProposal(oldProposal0),
		}

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 1, len(proofs))
		actualProof := proofs[0]
		require.Equal(t, expectedProof, actualProof)
	})

	t.Run("misbehaviour when in valid round there is a quorum of prevotes for a value different than old proposal", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0, committee)
		for _, m := range quorumPrevotes0VPrime {
			fd.msgStore.Save(m, committee)
		}

		expectedProof := &Proof{
			Offender:      proposer,
			OffenderIndex: proposerIdx,
			Type:          autonity.Misbehaviour,
			Rule:          autonity.PO,
			Evidences:     quorumPrevotes0VPrime,
			Message:       message.NewLightProposal(oldProposal0),
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
		fd.msgStore.Save(oldProposal0, committee)

		expectedProof := &Proof{
			OffenderIndex: proposerIdx,
			Offender:      proposer,
			Type:          autonity.Accusation,
			Rule:          autonity.PO,
			Message:       message.NewLightProposal(oldProposal0),
		}

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 1, len(proofs))
		actualProof := proofs[0]
		require.Equal(t, expectedProof, actualProof)
	})

	t.Run("accusation when less than quorum prevotes for proposal value in valid round", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0, committee)
		lessThanQurorumPrevotes := quorumPrevotes0V[:len(quorumPrevotes0V)-2]
		for _, m := range lessThanQurorumPrevotes {
			fd.msgStore.Save(m, committee)
		}

		expectedProof := &Proof{
			Offender:      proposer,
			OffenderIndex: proposerIdx,
			Type:          autonity.Accusation,
			Rule:          autonity.PO,
			Message:       message.NewLightProposal(oldProposal0),
		}

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 1, len(proofs))
		actualProof := proofs[0]
		require.Equal(t, expectedProof, actualProof)
	})

	t.Run("no proof for equivocated proposal with different valid round", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0, committee)
		fd.msgStore.Save(oldProposal0E, committee)

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof for equivocated proposal with same valid round however different block value", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0, committee)
		fd.msgStore.Save(oldProposal0E2, committee)

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when quorum of prevotes for V in vr, precommit for V from pi in vr, and precommit nils from pi from vr+1 to r", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0, committee)
		for _, m := range quorumPrevotes0V {
			fd.msgStore.Save(m, committee)
		}
		fd.msgStore.Save(nonNilPrecommit0V, committee)
		for _, m := range precommiteNilAfterVR {
			fd.msgStore.Save(m, committee)
		}

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when quorum of prevotes for V in vr, precommit for V from pi in vr, and some precommit nils from pi from vr+1 to r", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0, committee)
		for _, m := range quorumPrevotes0V {
			fd.msgStore.Save(m, committee)
		}
		fd.msgStore.Save(nonNilPrecommit0V, committee)
		somePrecommits := precommiteNilAfterVR[:len(precommiteNilAfterVR)-2]
		for _, m := range somePrecommits {
			fd.msgStore.Save(m, committee)
		}

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when quorum of prevotes for V in vr, precommit for V from pi in vr", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0, committee)
		for _, m := range quorumPrevotes0V {
			fd.msgStore.Save(m, committee)
		}
		fd.msgStore.Save(nonNilPrecommit0V, committee)

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when quorum of prevotes for V in vr, precommit nil from pi in vr", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0, committee)
		for _, m := range quorumPrevotes0V {
			fd.msgStore.Save(m, committee)
		}
		fd.msgStore.Save(nilPrecommit0, committee)

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when quorum of prevotes for V in vr", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0, committee)
		for _, m := range quorumPrevotes0V {
			fd.msgStore.Save(m, committee)
		}

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("multiple proofs from different rounds", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0, committee)
		fd.msgStore.Save(nonNilPrecommit0VPrime, committee)

		expectedMisbehaviour := &Proof{
			Offender:      proposer,
			OffenderIndex: proposerIdx,
			Type:          autonity.Misbehaviour,
			Rule:          autonity.PO,
			Evidences:     []message.Msg{nonNilPrecommit0VPrime},
			Message:       message.NewLightProposal(oldProposal0),
		}

		fd.msgStore.Save(oldProposal5, committee)
		expectedAccusation := &Proof{
			Offender:      proposer,
			OffenderIndex: proposerIdx,
			Type:          autonity.Accusation,
			Rule:          autonity.PO,
			Message:       message.NewLightProposal(oldProposal5),
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
	lastHeader := newBlockHeader(height-1, committee)
	header := newBlockHeader(height, committee)
	block := types.NewBlockWithHeader(header)
	header1 := newBlockHeader(height, committee)
	block1 := types.NewBlockWithHeader(header1)

	newProposalForB := newValidatedProposalMessage(height, 5, -1, signer, committee, block, proposerIdx)

	prevoteForB := newValidatedPrevote(t, 5, height, block.Hash(), signer, self, cSize, lastHeader)
	prevoteForB1 := newValidatedPrevote(t, 5, height, block1.Hash(), signer, self, cSize, lastHeader)

	precommitForB := newValidatedPrecommit(t, 3, height, block.Hash(), signer, self, cSize, lastHeader)
	precommitForB1 := newValidatedPrecommit(t, 4, height, block1.Hash(), signer, self, cSize, lastHeader)
	precommitForB1In0 := newValidatedPrecommit(t, 0, height, block1.Hash(), signer, self, cSize, lastHeader)
	precommitForB1In1 := newValidatedPrecommit(t, 1, height, block1.Hash(), signer, self, cSize, lastHeader)
	precommitForBIn0 := newValidatedPrecommit(t, 0, height, block.Hash(), signer, self, cSize, lastHeader)
	precommitForBIn4 := newValidatedPrecommit(t, 4, height, block.Hash(), signer, self, cSize, lastHeader)

	signerBis := makeSigner(keys[1])
	oldProposalB10 := newValidatedProposalMessage(height, 10, 5, signerBis, committee, block, 1)
	newProposalB1In5 := newValidatedProposalMessage(height, 5, -1, signerBis, committee, block1, 1)
	newProposalBIn5 := newValidatedProposalMessage(height, 5, -1, signerBis, committee, block, 1)

	prevoteForOldB10 := newValidatedPrevote(t, 10, height, block.Hash(), signer, self, cSize, lastHeader)
	precommitForB1In8 := newValidatedPrecommit(t, 8, height, block1.Hash(), signer, self, cSize, lastHeader)
	precommitForBIn7 := newValidatedPrecommit(t, 7, height, block.Hash(), signer, self, cSize, lastHeader)

	t.Run("accusation when there are no corresponding proposals", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(prevoteForB, committee)

		expectedAccusation := &Proof{
			Offender:      proposer,
			OffenderIndex: proposerIdx,
			Type:          autonity.Accusation,
			Rule:          autonity.PVN,
			Message:       prevoteForB,
		}
		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 1, len(proofs))
		require.Contains(t, proofs, expectedAccusation)
	})

	// Testcases for PVN
	t.Run("misbehaviour when pi precommited for a different value in a previous round than the prevoted value", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB, committee)
		fd.msgStore.Save(prevoteForB, committee)
		fd.msgStore.Save(precommitForB1, committee)

		expectedMisbehaviour := &Proof{
			Offender:      proposer,
			OffenderIndex: proposerIdx,
			Type:          autonity.Misbehaviour,
			Rule:          autonity.PVN,
			Evidences:     []message.Msg{message.NewLightProposal(newProposalForB), precommitForB1},
			Message:       prevoteForB,
		}

		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 1, len(proofs))
		require.Equal(t, expectedMisbehaviour, proofs[0])
	})

	t.Run("misbehaviour when pi precommited for a different value in a previous round than the prevoted value, after a flip flop", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB, committee)
		fd.msgStore.Save(prevoteForB, committee)
		fd.msgStore.Save(precommitForB1, committee)
		fd.msgStore.Save(precommitForB, committee)

		expectedMisbehaviour := &Proof{
			OffenderIndex: proposerIdx,
			Offender:      proposer,
			Type:          autonity.Misbehaviour,
			Rule:          autonity.PVN,
			Evidences:     []message.Msg{message.NewLightProposal(newProposalForB), precommitForB1},
			Message:       prevoteForB,
		}

		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 1, len(proofs))
		require.Equal(t, expectedMisbehaviour, proofs[0])
	})

	t.Run("misbehaviour when pi precommited for a different value in a previous round than the prevoted value while precommit nils in middle rounds", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB, committee)
		fd.msgStore.Save(prevoteForB, committee)
		fd.msgStore.Save(precommitForB1In0, committee)

		var precommitNilsAfter0 []message.Msg
		for i := 1; i < 5; i++ {
			precommitNil := newValidatedPrecommit(t, int64(i), height, nilValue, signer, self, cSize, lastHeader)
			precommitNilsAfter0 = append(precommitNilsAfter0, precommitNil)
			fd.msgStore.Save(precommitNil, committee)
		}

		expectedMisbehaviour := &Proof{
			Offender:      proposer,
			OffenderIndex: proposerIdx,
			Type:          autonity.Misbehaviour,
			Rule:          autonity.PVN,
			Evidences:     []message.Msg{precommitForB1In0},
			Message:       prevoteForB,
		}
		expectedMisbehaviour.Evidences = append(expectedMisbehaviour.Evidences, precommitNilsAfter0...)

		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
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
		fd.msgStore.Save(newProposalForB, committee)
		fd.msgStore.Save(prevoteForB, committee)
		fd.msgStore.Save(precommitForBIn0, committee)
		fd.msgStore.Save(precommitForB1In1, committee)

		var precommitNilsAfter1 []message.Msg
		for i := 2; i < 5; i++ {
			precommitNil := newValidatedPrecommit(t, int64(i), height, nilValue, signer, self, cSize, lastHeader)
			precommitNilsAfter1 = append(precommitNilsAfter1, precommitNil)
			fd.msgStore.Save(precommitNil, committee)
		}

		expectedMisbehaviour := &Proof{
			Offender:      proposer,
			OffenderIndex: proposerIdx,
			Type:          autonity.Misbehaviour,
			Rule:          autonity.PVN,
			Evidences:     []message.Msg{precommitForB1In1},
			Message:       prevoteForB,
		}
		expectedMisbehaviour.Evidences = append(expectedMisbehaviour.Evidences, precommitNilsAfter1...)

		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
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
		fd.msgStore.Save(newProposalForB, committee)
		fd.msgStore.Save(prevoteForB, committee)
		fd.msgStore.Save(precommitForBIn4, committee)

		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when pi precommited for the same value as the prevoted value in a previous round with missing precommits in middle rounds", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB, committee)
		fd.msgStore.Save(prevoteForB, committee)
		fd.msgStore.Save(precommitForBIn0, committee)

		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when pi precommited for the same value as the prevoted value in a previous round with some missing precommits and precommit nils in middle rounds", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB, committee)
		fd.msgStore.Save(prevoteForB, committee)
		fd.msgStore.Save(precommitForBIn0, committee)
		fd.msgStore.Save(newValidatedPrecommit(t, 3, height, nilValue, signer, self, cSize, lastHeader), committee)

		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when pi precommited for the same value as the prevoted value in a previous round with no missing precommits in middle rounds", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB, committee)
		fd.msgStore.Save(prevoteForB, committee)
		fd.msgStore.Save(precommitForBIn0, committee)
		for i := 1; i < 5; i++ {
			precommitNil := newValidatedPrecommit(t, int64(i), height, nilValue, signer, self, cSize, lastHeader)
			fd.msgStore.Save(precommitNil, committee)
		}

		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when pi precommited for {B1,nil*,B} and then prevoted B", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB, committee)
		fd.msgStore.Save(prevoteForB, committee)

		fd.msgStore.Save(precommitForB1In0, committee)

		// fill gaps with nil
		for i := 1; i < 4; i++ {
			precommitNil := newValidatedPrecommit(t, int64(i), height, nilValue, signer, self, cSize, lastHeader)
			fd.msgStore.Save(precommitNil, committee)
		}

		fd.msgStore.Save(precommitForBIn4, committee)

		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 0, len(proofs))
	})

	// Testcases for PVO
	t.Run("accusation when there is no quorum for the prevote value in the valid round", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposalB10, committee)
		fd.msgStore.Save(prevoteForOldB10, committee)

		expectedAccusation := &Proof{
			Offender:      proposer,
			OffenderIndex: proposerIdx,
			Type:          autonity.Accusation,
			Rule:          autonity.PVO,
			Message:       prevoteForOldB10,
			Evidences:     []message.Msg{message.NewLightProposal(oldProposalB10)},
		}

		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 1, len(proofs))
		acutalProof := proofs[0]
		require.Equal(t, expectedAccusation, acutalProof)
	})

	t.Run("misbehaviour when pi prevotes for an old proposal while in the valid round there is quorum for different value", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposalB10, committee)
		// Need to add this new proposal in valid round so that unwanted accusation are not returned by the prevotes
		// accountability check method. Since we are adding a quorum of prevotes in round 6 we also need to add a new
		// proposal in round 6 to allow for those prevotes to not return accusations.
		fd.msgStore.Save(newProposalB1In5, committee)
		fd.msgStore.Save(prevoteForOldB10, committee)
		// quorum of prevotes for B1 in vr = 6
		var vr5Prevotes []message.Msg
		for i := uint64(0); i < quorum.Uint64(); i++ {
			vr6Prevote := newValidatedPrevote(t, 5, height, block1.Hash(), makeSigner(keys[i]), &committee[i], cSize, lastHeader)
			vr5Prevotes = append(vr5Prevotes, vr6Prevote)
			fd.msgStore.Save(vr6Prevote, committee)
		}

		expectedMisbehaviour := &Proof{
			Offender:      proposer,
			OffenderIndex: proposerIdx,
			Type:          autonity.Misbehaviour,
			Rule:          autonity.PVO,
			Message:       prevoteForOldB10,
		}
		expectedMisbehaviour.Evidences = append(expectedMisbehaviour.Evidences, message.NewLightProposal(oldProposalB10))
		expectedMisbehaviour.Evidences = append(expectedMisbehaviour.Evidences, vr5Prevotes...)

		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
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
		fd.msgStore.Save(oldProposalB10, committee)
		fd.msgStore.Save(prevoteForOldB10, committee)
		fd.msgStore.Save(newProposalBIn5, committee)
		for i := 0; i < len(committee); i++ {
			fd.msgStore.Save(newValidatedPrevote(t, 5, height, block.Hash(), makeSigner(keys[i]), &committee[i], cSize, lastHeader), committee)
		}
		for i := newProposalBIn5.R(); i < precommitForBIn7.R(); i++ {
			fd.msgStore.Save(newValidatedPrecommit(t, i, height, nilValue, signer, self, cSize, lastHeader), committee)
		}
		var precommitsFromPiAfterLatestPrecommitForB []message.Msg
		fd.msgStore.Save(precommitForBIn7, committee)

		precommitsFromPiAfterLatestPrecommitForB = append(precommitsFromPiAfterLatestPrecommitForB, precommitForBIn7)
		fd.msgStore.Save(precommitForB1In8, committee)
		precommitsFromPiAfterLatestPrecommitForB = append(precommitsFromPiAfterLatestPrecommitForB, precommitForB1In8)
		p := newValidatedPrecommit(t, precommitForB1In8.R()+1, height, nilValue, signer, self, cSize, lastHeader)
		fd.msgStore.Save(p, committee)
		precommitsFromPiAfterLatestPrecommitForB = append(precommitsFromPiAfterLatestPrecommitForB, p)

		expectedMisbehaviour := &Proof{
			Offender:      proposer,
			OffenderIndex: proposerIdx,
			Type:          autonity.Misbehaviour,
			Rule:          autonity.PVO12,
			Message:       prevoteForOldB10,
		}
		expectedMisbehaviour.Evidences = append(expectedMisbehaviour.Evidences, message.NewLightProposal(oldProposalB10))
		expectedMisbehaviour.Evidences = append(expectedMisbehaviour.Evidences, precommitsFromPiAfterLatestPrecommitForB...)

		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
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
		fd.msgStore.Save(oldProposalB10, committee)
		fd.msgStore.Save(prevoteForOldB10, committee)
		fd.msgStore.Save(newProposalBIn5, committee)
		for i := 0; i < len(committee); i++ {
			v := newValidatedPrevote(t, 5, height, block.Hash(), makeSigner(keys[i]), &committee[i], cSize, lastHeader)
			fd.msgStore.Save(v, committee)
		}
		fd.msgStore.Save(precommitForBIn7, committee)
		for i := precommitForBIn7.R() + 1; i < oldProposalB10.R(); i++ {
			v := newValidatedPrecommit(t, i, height, nilValue, signer, self, cSize, lastHeader)
			fd.msgStore.Save(v, committee)
		}

		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 0, len(proofs))

	})

	t.Run("no proof when pi has precommited for V in a previous round however the latest precommit from pi is not for V yet pi still prevoted for V in the current round"+
		" but there are missing message after latest precommit for V", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposalB10, committee)
		fd.msgStore.Save(prevoteForOldB10, committee)
		fd.msgStore.Save(newProposalBIn5, committee)
		for i := 0; i < len(committee); i++ {
			fd.msgStore.Save(newValidatedPrevote(t, 5, height, block.Hash(), makeSigner(keys[i]), &committee[i], cSize, lastHeader), committee)
		}
		fd.msgStore.Save(precommitForBIn7, committee)
		fd.msgStore.Save(precommitForB1In8, committee)

		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 0, len(proofs))

	})

	t.Run("misbehaviour when pi has never precommited for V in a previous round however pi prevoted for V which is being reproposed", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposalB10, committee)
		fd.msgStore.Save(prevoteForOldB10, committee)
		fd.msgStore.Save(newProposalBIn5, committee)
		for i := 0; i < len(committee); i++ {
			fd.msgStore.Save(newValidatedPrevote(t, 5, height, block.Hash(), makeSigner(keys[i]), &committee[i], cSize, lastHeader), committee)
		}

		var precommitsFromPiAfterVR []message.Msg
		for i := newProposalBIn5.R() + 1; i < precommitForB1In8.R(); i++ {
			p := newValidatedPrecommit(t, i, height, nilValue, signer, self, cSize, lastHeader)
			fd.msgStore.Save(p, committee)
			precommitsFromPiAfterVR = append(precommitsFromPiAfterVR, p)
		}
		fd.msgStore.Save(precommitForB1In8, committee)
		precommitsFromPiAfterVR = append(precommitsFromPiAfterVR, precommitForB1In8)
		p := newValidatedPrecommit(t, precommitForB1In8.R()+1, height, nilValue, signer, self, cSize, lastHeader)
		fd.msgStore.Save(p, committee)
		precommitsFromPiAfterVR = append(precommitsFromPiAfterVR, p)

		expectedMisbehaviour := &Proof{
			Offender:      proposer,
			OffenderIndex: proposerIdx,
			Type:          autonity.Misbehaviour,
			Rule:          autonity.PVO12,
			Message:       prevoteForOldB10,
		}
		expectedMisbehaviour.Evidences = append(expectedMisbehaviour.Evidences, message.NewLightProposal(oldProposalB10))
		expectedMisbehaviour.Evidences = append(expectedMisbehaviour.Evidences, precommitsFromPiAfterVR...)

		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
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
		fd.msgStore.Save(oldProposalB10, committee)
		fd.msgStore.Save(prevoteForOldB10, committee)
		fd.msgStore.Save(newProposalBIn5, committee)
		for i := 0; i < len(committee); i++ {
			fd.msgStore.Save(newValidatedPrevote(t, 5, height, block.Hash(), makeSigner(keys[i]), &committee[i], cSize, lastHeader), committee)
		}

		for i := newProposalBIn5.R() + 1; i < oldProposalB10.R(); i++ {
			fd.msgStore.Save(newValidatedPrecommit(t, i, height, nilValue, signer, self, cSize, lastHeader), committee)
		}

		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when pi has never precommited for V in a previous round however pi prevoted for V while it has precommited for V' but there are missing precommit before precommit for V'", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposalB10, committee)
		fd.msgStore.Save(prevoteForOldB10, committee)
		fd.msgStore.Save(newProposalBIn5, committee)
		for i := 0; i < len(committee); i++ {
			fd.msgStore.Save(newValidatedPrevote(t, 5, height, block.Hash(), makeSigner(keys[i]), &committee[i], cSize, lastHeader), committee)
		}

		fd.msgStore.Save(precommitForB1In8, committee)

		p := newValidatedPrecommit(t, precommitForB1In8.R()+1, height, nilValue, signer, self, cSize, lastHeader)
		fd.msgStore.Save(p, committee)

		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when pi has never precommited for V in a previous round however pi prevoted for V while it has precommited for V' but there are missing precommit after precommit for V'", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposalB10, committee)

		fd.msgStore.Save(prevoteForOldB10, committee)
		fd.msgStore.Save(newProposalBIn5, committee)
		for i := 0; i < len(committee); i++ {
			fd.msgStore.Save(newValidatedPrevote(t, 5, height, block.Hash(), makeSigner(keys[i]), &committee[i], cSize, lastHeader), committee)
		}

		for i := newProposalBIn5.R() + 1; i < precommitForB1In8.R(); i++ {
			fd.msgStore.Save(newValidatedPrecommit(t, i, height, nilValue, signer, self, cSize, lastHeader), committee)
		}
		fd.msgStore.Save(precommitForB1In8, committee)

		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("prevotes accountability check can return multiple proofs", func(t *testing.T) {
		fd := testFD()

		fd.msgStore.Save(newProposalForB, committee)
		fd.msgStore.Save(prevoteForB, committee)
		fd.msgStore.Save(precommitForB1, committee)
		fd.msgStore.Save(precommitForB, committee)

		fd.msgStore.Save(oldProposalB10, committee)
		fd.msgStore.Save(prevoteForOldB10, committee)
		for i := 0; i < len(committee); i++ {
			fd.msgStore.Save(newValidatedPrevote(t, 6, height, block1.Hash(), makeSigner(keys[i]), &committee[i], cSize, lastHeader), committee)
		}

		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 2, len(proofs))
	})

	t.Run("no proof when prevote is equivocated with different values", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(prevoteForB, committee)
		fd.msgStore.Save(prevoteForB1, committee)

		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
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

	newProposalForB := newValidatedProposalMessage(height, 2, -1, makeSigner(keys[1]), committee, block, 1)

	precommitForB := newValidatedPrecommit(t, 2, height, block.Hash(), signer, self, cSize, header)
	precommitForB1 := newValidatedPrecommit(t, 2, height, block1.Hash(), signer, self, cSize, header)
	precommitForB1In3 := newValidatedPrecommit(t, 3, height, block1.Hash(), signer, self, cSize, header)

	t.Run("accusation when prevotes is less than quorum", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB, committee)
		fd.msgStore.Save(precommitForB, committee)

		for i := int64(0); i < quorum.Int64()-1; i++ {
			fd.msgStore.Save(newValidatedPrevote(t, 2, height, block.Hash(), makeSigner(keys[i]), &committee[i], cSize, header), committee)
		}

		expectedAccusation := &Proof{
			Offender:      proposer,
			OffenderIndex: proposerIdx,
			Type:          autonity.Accusation,
			Rule:          autonity.C1,
			Message:       precommitForB,
		}
		proofs := fd.precommitsAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 1, len(proofs))
		require.Equal(t, expectedAccusation, proofs[0])
	})

	t.Run("misbehaviour when there is a quorum for V' than what pi precommitted for", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB, committee)
		fd.msgStore.Save(precommitForB, committee)

		var prevotesForB1 []message.Msg
		for i := int64(0); i < quorum.Int64(); i++ {
			p := newValidatedPrevote(t, 2, height, block1.Hash(), makeSigner(keys[i]), &committee[i], cSize, header)
			fd.msgStore.Save(p, committee)
			prevotesForB1 = append(prevotesForB1, p)
		}

		expectedMisbehaviour := &Proof{
			Offender:      proposer,
			OffenderIndex: proposerIdx,
			Type:          autonity.Misbehaviour,
			Rule:          autonity.C,
			Evidences:     prevotesForB1,
			Message:       precommitForB,
		}
		proofs := fd.precommitsAccountabilityCheck(height, quorum, committee)
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
		fd.msgStore.Save(precommitForB1In3, committee)

		fd.msgStore.Save(newProposalForB, committee)
		fd.msgStore.Save(precommitForB, committee)

		var prevotesForB1 []message.Msg
		for i := int64(0); i < quorum.Int64(); i++ {
			p := newValidatedPrevote(t, 2, height, block1.Hash(), makeSigner(keys[i]), &committee[i], cSize, header)
			fd.msgStore.Save(p, committee)
			prevotesForB1 = append(prevotesForB1, p)
		}

		expectedProof0 := &Proof{
			Offender:      proposer,
			OffenderIndex: proposerIdx,
			Type:          autonity.Misbehaviour,
			Rule:          autonity.C,
			Evidences:     prevotesForB1,
			Message:       precommitForB,
		}

		expectedProof1 := &Proof{
			Offender:      proposer,
			OffenderIndex: proposerIdx,
			Type:          autonity.Accusation,
			Rule:          autonity.C1,
			Message:       precommitForB1In3,
		}
		proofs := fd.precommitsAccountabilityCheck(height, quorum, committee)
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
		fd.msgStore.Save(newProposalForB, committee)
		fd.msgStore.Save(precommitForB, committee)

		for i := 0; i < len(committee); i++ {
			fd.msgStore.Save(newValidatedPrevote(t, 2, height, block.Hash(), makeSigner(keys[i]), &committee[i], cSize, header), committee)
		}

		proofs := fd.precommitsAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when there is more than quorum prevotes ", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB, committee)
		fd.msgStore.Save(precommitForB, committee)

		for i := 0; i < len(committee); i++ {
			fd.msgStore.Save(newValidatedPrevote(t, 2, height, block.Hash(), makeSigner(keys[i]), &committee[i], cSize, header), committee)
		}

		proofs := fd.precommitsAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when precommit is equivocated with different values", func(t *testing.T) {
		//t.Skip("not stable in CI, but work in local.")
		fd := testFD()
		fd.msgStore.Save(precommitForB, committee)
		fd.msgStore.Save(precommitForB1, committee)

		proofs := fd.precommitsAccountabilityCheck(height, quorum, committee)
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
