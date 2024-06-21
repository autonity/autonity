package accountability

import (
	"crypto/ecdsa"
	"math/big"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/autonity/autonity/consensus/ethash"
	"github.com/autonity/autonity/core/rawdb"
	"github.com/autonity/autonity/core/vm"

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
	cSize                     = committee.Len()
	proposerIdx               = 0
	prevoterIdx               = 2
	proposer                  = committee.Members[proposerIdx].Address
	proposerKey               = keys[proposerIdx]
	proposerNodeKey           = nodeKeys[proposerIdx]
	signer                    = makeSigner(proposerKey)
	self                      = &committee.Members[proposerIdx]

	remotePeerIdx = 1
	remote        = &committee.Members[remotePeerIdx]
	remotePeer    = committee.Members[remotePeerIdx].Address
	remoteSigner  = makeSigner(keys[remotePeerIdx])
)

func generateCommittee() (*types.Committee, []blst.SecretKey, []*ecdsa.PrivateKey) {
	n := 5
	c := new(types.Committee)
	pkeys := make([]*ecdsa.PrivateKey, n)
	consensusKeys := make([]blst.SecretKey, n)
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
		c.Members = append(c.Members, committeeMember)
		pkeys[i] = privateKey
		consensusKeys[i] = consensusKey
	}
	return c, consensusKeys, pkeys
}

func newBlockHeader(height uint64, committee *types.Committee) *types.Header {
	// use random nonce to create different blocks
	var nonce types.BlockNonce
	for i := 0; i < len(nonce); i++ {
		nonce[i] = byte(rand.Intn(256)) //nolint
	}

	var epoch types.Epoch
	epoch.Committee = committee
	epoch.PreviousEpochBlock = common.Big0
	epoch.NextEpochBlock = new(big.Int).SetUint64(height + 30)

	header := &types.Header{
		Number: new(big.Int).SetUint64(height),
		Nonce:  nonce,
	}

	header.Epoch = &epoch

	if err := header.Epoch.Committee.Enrich(); err != nil {
		panic(err)
	}

	return header
}

// new proposal with metadata, if the withValue is not nil, it will use the value as proposal, otherwise a
// random block will be used as the value for proposal.
func newValidatedProposalMessage(h uint64, r int64, vr int64, signer message.Signer, committee *types.Committee, withValue *types.Block, idx int) *message.Propose {
	block := withValue
	if withValue == nil {
		header := newBlockHeader(h, committee)
		block = types.NewBlockWithHeader(header)
	}
	return message.NewPropose(r, h, vr, block, signer, &committee.Members[idx])
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
	// submit a equivocation proofs.
	proposal := newValidatedLightProposal(height, round, -1, signer, committee, nil, proposerIdx)
	proposal2 := newValidatedLightProposal(height, round, -1, signer, committee, nil, proposerIdx)
	var proofs []message.Msg
	proofs = append(proofs, proposal2)

	fd := &FaultDetector{
		misbehaviourProofCh: make(chan *autonity.AccountabilityEvent, 100),
		logger:              log.New("FaultDetector", nil),
	}

	fd.submitMisbehavior(proposal, proofs, errEquivocation, proposerIdx, proposer)
	p := <-fd.misbehaviourProofCh

	require.Equal(t, uint8(autonity.Misbehaviour), p.EventType)
	require.Equal(t, proposer, p.Offender)
}

func TestRunRuleEngine(t *testing.T) {
	round := int64(3)
	t.Run("test run rules with malicious behaviour should be detected", func(t *testing.T) {
		chainHead := uint64(100)
		checkPointHeight := chainHead - uint64(DeltaBlocks)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chainMock := NewMockChainContext(ctrl)
		chainMock.EXPECT().CommitteeOfHeight(checkPointHeight).Return(committee, nil)
		chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})
		var blockSub event.Subscription
		chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
		fdAddr := committee.Members[1].Address
		accountability, _ := autonity.NewAccountability(proposer, backends.NewSimulatedBackend(ccore.GenesisAlloc{fdAddr: {Balance: big.NewInt(params.Ether)}}, 10000000))

		fd := NewFaultDetector(chainMock, fdAddr, nil, core.NewMsgStore(), nil, nil, proposerNodeKey, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())
		// store a msg before check point height in case of node is start from reset.
		msgBeforeCheckPointHeight := newValidatedProposalMessage(checkPointHeight-1, 0, -1, makeSigner(keys[1]), committee, nil, 1)
		fd.msgStore.Save(msgBeforeCheckPointHeight)

		// simulate there was a maliciousProposal at init round 0, and save to msg store.
		initProposal := newValidatedProposalMessage(checkPointHeight, 0, -1, makeSigner(keys[1]), committee, nil, 1)
		fd.msgStore.Save(initProposal)

		aggregatedVotes := aggregatedPreVote(committee.Len(), checkPointHeight, 0, initProposal.Value(), keys, committee)
		fd.msgStore.Save(aggregatedVotes)

		// Node preCommit for init Proposal at init round 0 since there were quorum preVotes for it, and save it.
		preCommit := newValidatedPrecommit(0, checkPointHeight, initProposal.Value(), signer, self, cSize)
		fd.msgStore.Save(preCommit)

		// While Node propose a new malicious Proposal at new round with VR as -1 which is malicious, should be addressed by rule PN.
		maliciousProposal := newValidatedProposalMessage(checkPointHeight, round, -1, signer, committee, nil, proposerIdx)
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

func TestGenerateOnChainProof(t *testing.T) {
	height := uint64(100)
	round := int64(3)

	proposal := newValidatedLightProposal(height, round, -1, signer, committee, nil, proposerIdx)
	equivocatedProposal := newValidatedLightProposal(height, round, -1, signer, committee, nil, proposerIdx)
	var evidence []message.Msg
	evidence = append(evidence, equivocatedProposal)

	p := Proof{
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

	onChainEvent := fd.eventFromProof(&p, proposer)

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
		methodName := "handleMisbehaviour"
		_, err := generated.AccountabilityAbi.Pack(methodName, onChainEvent)
		require.NoError(t, err)
	})
}

// todo: (Jason) add test to cover an accusation over a committed block scenario,
//
//	in such context, the accusation is considered as useless, it should be dropped.
func TestAccusationProvers(t *testing.T) {
	height := uint64(100)
	round := int64(3)
	validRound := int64(1)
	noneNilValue := common.Hash{0x1}

	t.Run("innocenceProofPO have quorum preVotes", func(t *testing.T) {

		// PO: node propose an old value with an validRound, innocent Proof of it should be:
		// there were quorum num of preVote for that value at the validRound.
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chainMock := NewMockChainContext(ctrl)
		var blockSub event.Subscription
		chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
		chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})
		bindings, _ := autonity.NewAccountability(proposer, backends.NewSimulatedBackend(ccore.GenesisAlloc{proposer: {Balance: big.NewInt(params.Ether)}}, 10000000))

		fd := NewFaultDetector(chainMock, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}), core.NewMsgStore(), nil, nil, proposerNodeKey, &autonity.ProtocolContracts{Accountability: bindings}, log.Root())
		// simulate a proposal message with an old value and a valid round.
		proposal := newValidatedProposalMessage(height, round, validRound, signer, committee, nil, proposerIdx)
		fd.msgStore.Save(proposal)

		// simulate at least quorum num of preVotes for a value at a validRound.
		aggregatedVote := aggregatedPreVote(committee.Len(), height, validRound, proposal.Value(), keys, committee)
		fd.msgStore.Save(aggregatedVote)

		var accusation = Proof{
			OffenderIndex: proposerIdx,
			Type:          autonity.Accusation,
			Rule:          autonity.PO,
			Message:       proposal.ToLight(),
		}

		proof, err := fd.innocenceProofPO(&accusation, committee)
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
		chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})
		accountability, _ := autonity.NewAccountability(proposer, backends.NewSimulatedBackend(ccore.GenesisAlloc{proposer: {Balance: big.NewInt(params.Ether)}}, 10000000))
		var blockSub event.Subscription
		chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
		fd := NewFaultDetector(chainMock, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}), core.NewMsgStore(), nil, nil, proposerNodeKey, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())
		// simulate a proposal message with an old value and a valid round.
		proposal := newValidatedProposalMessage(height, round, validRound, signer, committee, nil, proposerIdx)
		fd.msgStore.Save(proposal)

		// simulate less than quorum num of preVotes for a value at a validRound.
		preVote := newValidatedPrevote(validRound, height, proposal.Value(), signer, self, cSize)
		fd.msgStore.Save(preVote)

		var accusation = Proof{
			Type:          autonity.Accusation,
			Rule:          autonity.PO,
			OffenderIndex: proposerIdx,
			Message:       message.NewLightProposal(proposal),
		}

		_, err := fd.innocenceProofPO(&accusation, committee)
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
		fd.msgStore.Save(proposal)

		// simulate at least quorum num of preVotes for a value at a validRound.
		aggregatedVote := aggregatedPreVote(committee.Len(), height, round, proposal.Value(), keys, committee)
		fd.msgStore.Save(aggregatedVote)

		preVote := newValidatedPrevote(round, height, proposal.Value(), signer, self, cSize)

		var accusation = Proof{
			OffenderIndex: proposerIdx,
			Type:          autonity.Accusation,
			Rule:          autonity.PVN,
			Message:       preVote,
		}

		proof, err := fd.innocenceProofPVN(&accusation, committee)
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

		preVote := newValidatedPrevote(round, height, noneNilValue, signer, self, cSize)
		fd.msgStore.Save(preVote)

		var accusation = Proof{
			OffenderIndex: proposerIdx,
			Type:          autonity.Accusation,
			Rule:          autonity.PVN,
			Message:       preVote,
		}

		_, err := fd.innocenceProofPVN(&accusation, committee)
		assert.Equal(t, errNoEvidenceForPVN, err)
	})

	t.Run("getInnocentProofofPVO have no quorum preVotes", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chainMock := NewMockChainContext(ctrl)
		var blockSub event.Subscription
		chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
		chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})
		accountability, _ := autonity.NewAccountability(proposer, backends.NewSimulatedBackend(ccore.GenesisAlloc{proposer: {Balance: big.NewInt(params.Ether)}}, 10000000))

		fd := NewFaultDetector(chainMock, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}), core.NewMsgStore(), nil, nil, proposerNodeKey, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())

		var p Proof
		p.Rule = autonity.PVO
		p.OffenderIndex = proposerIdx
		oldProposal := newValidatedProposalMessage(height, 1, 0, signer, committee, nil, proposerIdx)
		preVote := newValidatedPrevote(1, height, oldProposal.Value(), signer, self, cSize)
		p.Message = preVote
		p.Evidences = append(p.Evidences, oldProposal.ToLight())

		_, err := fd.innocenceProofPVO(&p, committee)
		assert.Equal(t, err, errNoEvidenceForPVO)
	})

	t.Run("getInnocentProofofPVO have quorum preVotes", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chainMock := NewMockChainContext(ctrl)
		var blockSub event.Subscription
		chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
		chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})
		accountability, _ := autonity.NewAccountability(proposer, backends.NewSimulatedBackend(ccore.GenesisAlloc{proposer: {Balance: big.NewInt(params.Ether)}}, 10000000))

		fd := NewFaultDetector(chainMock, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}), core.NewMsgStore(), nil, nil, proposerNodeKey, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())
		var p Proof
		p.Rule = autonity.PVO
		p.OffenderIndex = proposerIdx
		validRound := int64(0)
		oldProposal := newValidatedProposalMessage(height, 1, validRound, signer, committee, nil, proposerIdx)
		preVote := newValidatedPrevote(1, height, oldProposal.Value(), signer, self, cSize)
		p.Message = preVote
		p.Evidences = append(p.Evidences, message.NewLightProposal(oldProposal))

		// prepare quorum preVotes at msg store.
		aggregatedVote := aggregatedPreVote(committee.Len(), height, validRound, oldProposal.Value(), keys, committee)
		fd.msgStore.Save(aggregatedVote)

		onChainProof, err := fd.innocenceProofPVO(&p, committee)
		assert.NoError(t, err)
		assert.Equal(t, uint8(autonity.Innocence), onChainProof.EventType)
		assert.Equal(t, proposer, onChainProof.Reporter)

	})

	t.Run("innocenceProofC1 have quorum preVotes", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chainMock := NewMockChainContext(ctrl)
		var blockSub event.Subscription
		chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
		chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})
		accountability, _ := autonity.NewAccountability(proposer, backends.NewSimulatedBackend(ccore.GenesisAlloc{proposer: {Balance: big.NewInt(params.Ether)}}, 10000000))

		// C1: node preCommit at a none nil value, there must be quorum corresponding preVotes with same value and round.
		fd := NewFaultDetector(chainMock, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}), core.NewMsgStore(), nil, nil, proposerNodeKey, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())

		// simulate at least quorum num of preVotes for a value at a validRound.
		aggregatedVote := aggregatedPreVote(committee.Len(), height, round, noneNilValue, keys, committee)
		fd.msgStore.Save(aggregatedVote)

		preCommit := newValidatedPrecommit(round, height, noneNilValue, signer, self, cSize)
		fd.msgStore.Save(preCommit)

		var accusation = Proof{
			OffenderIndex: proposerIdx,
			Type:          autonity.Accusation,
			Rule:          autonity.C1,
			Message:       preCommit,
		}

		proof, err := fd.innocenceProofC1(&accusation, committee)
		assert.NoError(t, err)
		assert.Equal(t, uint8(autonity.Innocence), proof.EventType)
		assert.Equal(t, proposer, proof.Reporter)

	})

	t.Run("innocenceProofC1 have no quorum preVotes", func(t *testing.T) {

		// C1: node preCommit at a none nil value, there must be quorum corresponding preVotes with same value and round.
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chainMock := NewMockChainContext(ctrl)
		var blockSub event.Subscription
		chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
		chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})
		accountability, _ := autonity.NewAccountability(proposer, backends.NewSimulatedBackend(ccore.GenesisAlloc{proposer: {Balance: big.NewInt(params.Ether)}}, 10000000))

		fd := NewFaultDetector(chainMock, proposer, new(event.TypeMux).Subscribe(events.MessageEvent{}), core.NewMsgStore(), nil, nil, proposerNodeKey, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())

		preCommit := newValidatedPrecommit(round, height, noneNilValue, signer, self, cSize)
		fd.msgStore.Save(preCommit)

		var accusation = Proof{
			OffenderIndex: proposerIdx,
			Type:          autonity.Accusation,
			Rule:          autonity.C1,
			Message:       preCommit,
		}

		_, err := fd.innocenceProofC1(&accusation, committee)
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
	newProposal0 := newValidatedProposalMessage(height, 3, -1, signer, committee, nil, proposerIdx)
	nonNilPrecommit0 := newValidatedPrecommit(1, height, common.BytesToHash([]byte("test")), signer, self, cSize)
	nilPrecommit0 := newValidatedPrecommit(1, height, common.Hash{}, signer, self, cSize)

	newProposal1 := newValidatedProposalMessage(height, 5, -1, signer, committee, nil, proposerIdx)
	nilPrecommit1 := newValidatedPrecommit(3, height, common.Hash{}, signer, self, cSize)

	newProposal0E := newValidatedProposalMessage(height, 3, 1, signer, committee, nil, proposerIdx)

	t.Run("misbehaviour when pi has sent a non-nil precommit in a previous round", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposal0)
		fd.msgStore.Save(nonNilPrecommit0)

		expectedProof := &Proof{
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
			OffenderIndex: proposerIdx,
			Type:          autonity.Misbehaviour,
			Rule:          autonity.PN,
			Evidences:     []message.Msg{nonNilPrecommit0},
			Message:       newProposal0,
		}

		expectedProof1 := &Proof{
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

	header := newBlockHeader(height, committee)
	block := types.NewBlockWithHeader(header)
	header1 := newBlockHeader(height, committee)
	block1 := types.NewBlockWithHeader(header1)

	oldProposal0 := newValidatedProposalMessage(height, 3, 0, signer, committee, block, proposerIdx)
	oldProposal5 := newValidatedProposalMessage(height, 5, 2, signer, committee, block, proposerIdx)
	oldProposal0E := newValidatedProposalMessage(height, 3, 2, signer, committee, block1, proposerIdx)
	oldProposal0E2 := newValidatedProposalMessage(height, 3, 0, signer, committee, block1, proposerIdx)

	nonNilPrecommit0V := newValidatedPrecommit(0, height, block.Hash(), signer, self, cSize)
	nonNilPrecommit0VPrime := newValidatedPrecommit(0, height, block1.Hash(), signer, self, cSize)
	nonNilPrecommit2VPrime := newValidatedPrecommit(2, height, block1.Hash(), signer, self, cSize)
	nonNilPrecommit1 := newValidatedPrecommit(1, height, block.Hash(), signer, self, cSize)

	nilPrecommit0 := newValidatedPrecommit(0, height, nilValue, signer, self, cSize)
	quorumPrevotes0VPrime := aggregatedPreVote(int(quorum.Int64()), height, 0, block1.Hash(), keys, committee)
	quorumPrevotes0V := aggregatedPreVote(int(quorum.Int64()), height, 0, block.Hash(), keys, committee)
	lessThanQurorumPrevotes := aggregatedPreVote(int(quorum.Int64())-1, height, 0, block.Hash(), keys, committee)

	var precommiteNilAfterVR []message.Msg
	for i := 1; i < 3; i++ {
		precommit := newValidatedPrecommit(int64(i), height, nilValue, signer, self, cSize)
		precommiteNilAfterVR = append(precommiteNilAfterVR, precommit)
	}

	t.Run("misbehaviour when pi precommited for a different value in valid round than in the old proposal", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0)
		fd.msgStore.Save(nonNilPrecommit0VPrime)

		expectedProof := &Proof{
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
		fd.msgStore.Save(oldProposal0)
		fd.msgStore.Save(nonNilPrecommit2VPrime)

		expectedProof := &Proof{
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
		fd.msgStore.Save(oldProposal0)
		fd.msgStore.Save(nonNilPrecommit1)

		expectedProof := &Proof{
			OffenderIndex: proposerIdx,
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
		fd.msgStore.Save(oldProposal0)
		fd.msgStore.Save(quorumPrevotes0VPrime)

		expectedProof := &Proof{
			OffenderIndex: proposerIdx,
			Type:          autonity.Misbehaviour,
			Rule:          autonity.PO,
			Evidences:     []message.Msg{quorumPrevotes0VPrime},
			Message:       message.NewLightProposal(oldProposal0),
		}

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 1, len(proofs))
		actualProof := proofs[0]
		require.Equal(t, expectedProof.Type, actualProof.Type)
		require.Equal(t, expectedProof.Rule, actualProof.Rule)
		require.Equal(t, expectedProof.Message, actualProof.Message)
		require.Equal(t, expectedProof.Evidences[0], actualProof.Evidences[0])
	})

	t.Run("accusation when no prevotes for proposal value in valid round", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0)

		expectedProof := &Proof{
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

	t.Run("accusation when less than quorum prevotes for proposal value in valid round", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0)
		fd.msgStore.Save(lessThanQurorumPrevotes)

		expectedProof := &Proof{
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
		fd.msgStore.Save(quorumPrevotes0V)
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
		fd.msgStore.Save(quorumPrevotes0V)
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
		fd.msgStore.Save(quorumPrevotes0V)
		fd.msgStore.Save(nonNilPrecommit0V)

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when quorum of prevotes for V in vr, precommit nil from pi in vr", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0)
		fd.msgStore.Save(quorumPrevotes0V)
		fd.msgStore.Save(nilPrecommit0)

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when quorum of prevotes for V in vr", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0)
		fd.msgStore.Save(quorumPrevotes0V)

		proofs := fd.oldProposalsAccountabilityCheck(height, quorum)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("multiple proofs from different rounds", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposal0)
		fd.msgStore.Save(nonNilPrecommit0VPrime)

		expectedMisbehaviour := &Proof{
			OffenderIndex: proposerIdx,
			Type:          autonity.Misbehaviour,
			Rule:          autonity.PO,
			Evidences:     []message.Msg{nonNilPrecommit0VPrime},
			Message:       message.NewLightProposal(oldProposal0),
		}

		fd.msgStore.Save(oldProposal5)
		expectedAccusation := &Proof{
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
	height := uint64(1)
	parent := newBlockHeader(height-1, committee)
	header := newBlockHeader(height, committee)
	block := types.NewBlockWithHeader(header)
	header1 := newBlockHeader(height, committee)
	block1 := types.NewBlockWithHeader(header1)

	newProposalForB := newValidatedProposalMessage(height, 5, -1, signer, committee, block, proposerIdx)

	prevoteForB := newValidatedPrevote(5, height, block.Hash(), signer, self, cSize)
	prevoteForB1 := newValidatedPrevote(5, height, block1.Hash(), signer, self, cSize)

	otherPrevoteForB := newValidatedPrevote(prevoteForB.R(), prevoteForB.H(), prevoteForB.Value(),
		makeSigner(keys[prevoterIdx]), &committee.Members[prevoterIdx], cSize)
	otherPrevoteForB1 := newValidatedPrevote(prevoteForB1.R(), prevoteForB1.H(), prevoteForB1.Value(),
		makeSigner(keys[prevoterIdx]), &committee.Members[prevoterIdx], cSize)

	aggregatedPrevoteForB := message.AggregatePrevotes([]message.Vote{prevoteForB, otherPrevoteForB})
	aggregatedPrevoteForB1 := message.AggregatePrevotes([]message.Vote{prevoteForB1, otherPrevoteForB1})

	precommitForB := newValidatedPrecommit(3, height, block.Hash(), signer, self, cSize)
	otherPrecommitForB := newValidatedPrecommit(precommitForB.R(), precommitForB.H(), precommitForB.Value(),
		makeSigner(keys[prevoterIdx]), &committee.Members[prevoterIdx], cSize)
	aggregatedPrecommitForB := message.AggregatePrecommits([]message.Vote{precommitForB, otherPrecommitForB})

	precommitForB1 := newValidatedPrecommit(4, height, block1.Hash(), signer, self, cSize)
	otherPrecommitForB1 := newValidatedPrecommit(precommitForB1.R(), precommitForB1.H(), precommitForB1.Hash(),
		makeSigner(keys[prevoterIdx]), &committee.Members[prevoterIdx], cSize)
	aggregatedPrecommitForB1 := message.AggregatePrecommits([]message.Vote{precommitForB1, otherPrecommitForB1})

	precommitForB1In0 := newValidatedPrecommit(0, height, block1.Hash(), signer, self, cSize)
	precommitForB1In1 := newValidatedPrecommit(1, height, block1.Hash(), signer, self, cSize)
	precommitForBIn0 := newValidatedPrecommit(0, height, block.Hash(), signer, self, cSize)
	precommitForBIn4 := newValidatedPrecommit(4, height, block.Hash(), signer, self, cSize)

	signerBis := makeSigner(keys[1])
	oldProposalB10 := newValidatedProposalMessage(height, 10, 5, signerBis, committee, block, 1)
	newProposalB1In5 := newValidatedProposalMessage(height, 5, -1, signerBis, committee, block1, 1)
	newProposalBIn5 := newValidatedProposalMessage(height, 5, -1, signerBis, committee, block, 1)

	prevoteForOldB10 := newValidatedPrevote(10, height, block.Hash(), signer, self, cSize)
	otherPrevoteForOldB10 := newValidatedPrevote(prevoteForOldB10.R(), prevoteForOldB10.H(), prevoteForOldB10.Value(),
		makeSigner(keys[prevoterIdx]), &committee.Members[prevoterIdx], cSize)
	aggregatedPrevoteForOldB10 := message.AggregatePrevotes([]message.Vote{prevoteForOldB10, otherPrevoteForOldB10})

	precommitForB1In8 := newValidatedPrecommit(8, height, block1.Hash(), signer, self, cSize)
	otherPrecommitForB1In8 := newValidatedPrecommit(8, height, block1.Hash(), makeSigner(keys[prevoterIdx]),
		&committee.Members[prevoterIdx], cSize)
	aggregatedPrecommitForB1In8 := message.AggregatePrecommits([]message.Vote{precommitForB1In8, otherPrecommitForB1In8})

	precommitForBIn7 := newValidatedPrecommit(7, height, block.Hash(), signer, self, cSize)

	t.Run("accusation when there are no corresponding proposals", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(prevoteForB)
		expectedAccusation := &Proof{
			OffenderIndex: proposerIdx,
			Type:          autonity.Accusation,
			Rule:          autonity.PVN,
			Message:       prevoteForB,
		}
		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 1, len(proofs))
		require.Contains(t, proofs, expectedAccusation)
	})

	t.Run("accusation of aggregated prevotes when there are no corresponding proposals", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(aggregatedPrevoteForB)
		expectedAccusation1 := &Proof{
			OffenderIndex: proposerIdx,
			Type:          autonity.Accusation,
			Rule:          autonity.PVN,
			Message:       aggregatedPrevoteForB,
		}
		expectedAccusation2 := &Proof{
			OffenderIndex: prevoterIdx,
			Type:          autonity.Accusation,
			Rule:          autonity.PVN,
			Message:       aggregatedPrevoteForB,
		}
		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 2, len(proofs))
		require.Contains(t, proofs, expectedAccusation1)
		require.Contains(t, proofs, expectedAccusation2)
	})

	// Testcases for PVN
	t.Run("misbehaviour when pi precommited for a different value in a previous round than the prevoted value", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(aggregatedPrevoteForB)
		fd.msgStore.Save(aggregatedPrecommitForB1)

		expectedMisbehaviour1 := &Proof{
			OffenderIndex: proposerIdx,
			Type:          autonity.Misbehaviour,
			Rule:          autonity.PVN,
			Evidences:     []message.Msg{message.NewLightProposal(newProposalForB), aggregatedPrecommitForB1},
			Message:       aggregatedPrevoteForB,
		}

		expectedMisbehaviour2 := &Proof{
			OffenderIndex: prevoterIdx,
			Type:          autonity.Misbehaviour,
			Rule:          autonity.PVN,
			Evidences:     []message.Msg{message.NewLightProposal(newProposalForB), aggregatedPrecommitForB1},
			Message:       aggregatedPrevoteForB,
		}

		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 2, len(proofs))
		require.Contains(t, proofs, expectedMisbehaviour1)
		require.Contains(t, proofs, expectedMisbehaviour2)
	})

	t.Run("misbehaviour when pi precommited for a different value in a previous round than the prevoted value", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(aggregatedPrevoteForB)
		fd.msgStore.Save(aggregatedPrecommitForB1)
		fd.msgStore.Save(aggregatedPrecommitForB) // this is not required for PVN detection.

		expectedMisbehaviour1 := &Proof{
			OffenderIndex: proposerIdx,
			Type:          autonity.Misbehaviour,
			Rule:          autonity.PVN,
			Evidences:     []message.Msg{message.NewLightProposal(newProposalForB), aggregatedPrecommitForB1},
			Message:       aggregatedPrevoteForB,
		}

		expectedMisbehaviour2 := &Proof{
			OffenderIndex: prevoterIdx,
			Type:          autonity.Misbehaviour,
			Rule:          autonity.PVN,
			Evidences:     []message.Msg{message.NewLightProposal(newProposalForB), aggregatedPrecommitForB1},
			Message:       aggregatedPrevoteForB,
		}

		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 2, len(proofs))
		require.Contains(t, proofs, expectedMisbehaviour1)
		require.Contains(t, proofs, expectedMisbehaviour2)
	})

	t.Run("misbehaviour when pi precommited for a different value in a previous round than the prevoted value while precommit nils in middle rounds", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(aggregatedPrevoteForB)
		fd.msgStore.Save(precommitForB1In0)

		aggregatablePrecomits := make([]*message.Precommit, 5)
		aggregatablePrecomits[0] = precommitForB1In0
		for i := 1; i < 5; i++ {
			precommitNil := newValidatedPrecommit(int64(i), height, nilValue, signer, self, cSize)
			fd.msgStore.Save(precommitNil)
			aggregatablePrecomits[i] = precommitNil
		}

		aggPrecomit := AggregateDistinctPrecommits(aggregatablePrecomits)

		expectedMisbehaviour := &Proof{
			OffenderIndex:      proposerIdx,
			Type:               autonity.Misbehaviour,
			Rule:               autonity.PVN,
			DistinctPrecommits: aggPrecomit,
			Message:            aggregatedPrevoteForB,
		}

		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 1, len(proofs))
		actualProof := proofs[0]
		require.Equal(t, expectedMisbehaviour.Type, actualProof.Type)
		require.Equal(t, expectedMisbehaviour.Rule, actualProof.Rule)
		require.Equal(t, expectedMisbehaviour.Message, actualProof.Message)
		require.Equal(t, expectedMisbehaviour.DistinctPrecommits, actualProof.DistinctPrecommits)
		err := actualProof.DistinctPrecommits.PreValidate(parent.Epoch.Committee, height)
		require.NoError(t, err)
		err = actualProof.DistinctPrecommits.Validate()
		require.NoError(t, err)
	})

	t.Run("misbehaviour when pi precommited for a different value in a previous round than the prevoted value, after a flip flop, while precommit nils in middle rounds", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(aggregatedPrevoteForB)
		fd.msgStore.Save(precommitForBIn0)
		fd.msgStore.Save(precommitForB1In1)

		var precommits []*message.Precommit
		precommits = append(precommits, precommitForB1In1)

		for i := 2; i < 5; i++ {
			precommitNil := newValidatedPrecommit(int64(i), height, nilValue, signer, self, cSize)
			precommits = append(precommits, precommitNil)
			fd.msgStore.Save(precommitNil)
		}
		aggPrecomits := AggregateDistinctPrecommits(precommits)

		expectedMisbehaviour := &Proof{
			OffenderIndex:      proposerIdx,
			Type:               autonity.Misbehaviour,
			Rule:               autonity.PVN,
			DistinctPrecommits: aggPrecomits,
			Message:            aggregatedPrevoteForB,
		}

		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 1, len(proofs))
		actualProof := proofs[0]
		require.Equal(t, expectedMisbehaviour.Type, actualProof.Type)
		require.Equal(t, expectedMisbehaviour.Rule, actualProof.Rule)
		require.Equal(t, expectedMisbehaviour.Message, actualProof.Message)
		for _, m := range expectedMisbehaviour.Evidences {
			require.Contains(t, actualProof.Evidences, m)
		}
		require.Equal(t, expectedMisbehaviour.DistinctPrecommits, actualProof.DistinctPrecommits)
		err := actualProof.DistinctPrecommits.PreValidate(parent.Epoch.Committee, height)
		require.NoError(t, err)
		err = actualProof.DistinctPrecommits.Validate()
		require.NoError(t, err)
	})

	t.Run("no proof when pi precommited for the same value as the prevoted value in a previous round", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(aggregatedPrevoteForB)
		fd.msgStore.Save(precommitForBIn4)

		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when pi precommited for the same value as the prevoted value in a previous round with missing precommits in middle rounds", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(prevoteForB)
		fd.msgStore.Save(precommitForBIn0)

		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when pi precommited for the same value as the prevoted value in a previous round with some missing precommits and precommit nils in middle rounds", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(aggregatedPrevoteForB)
		fd.msgStore.Save(precommitForBIn0)
		fd.msgStore.Save(newValidatedPrecommit(3, height, nilValue, signer, self, cSize))

		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when pi precommited for the same value as the prevoted value in a previous round with no missing precommits in middle rounds", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(aggregatedPrevoteForB)
		fd.msgStore.Save(precommitForBIn0)
		for i := 1; i < 5; i++ {
			precommitNil := newValidatedPrecommit(int64(i), height, nilValue, signer, self, cSize)
			fd.msgStore.Save(precommitNil)
		}

		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when pi precommited for {B1,nil*,B} and then prevoted B", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(aggregatedPrevoteForB)

		fd.msgStore.Save(precommitForB1In0)

		// fill gaps with nil
		for i := 1; i < 4; i++ {
			precommitNil := newValidatedPrecommit(int64(i), height, nilValue, signer, self, cSize)
			fd.msgStore.Save(precommitNil)
		}

		fd.msgStore.Save(precommitForBIn4)

		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 0, len(proofs))
	})

	// Testcases for PVO
	t.Run("accusation when there is no quorum for the prevote value in the valid round", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposalB10)
		fd.msgStore.Save(aggregatedPrevoteForOldB10)

		expectedAccusation1 := &Proof{
			OffenderIndex: proposerIdx,
			Type:          autonity.Accusation,
			Rule:          autonity.PVO,
			Message:       aggregatedPrevoteForOldB10,
			Evidences:     []message.Msg{message.NewLightProposal(oldProposalB10)},
		}
		expectedAccusation2 := &Proof{
			OffenderIndex: prevoterIdx,
			Type:          autonity.Accusation,
			Rule:          autonity.PVO,
			Message:       aggregatedPrevoteForOldB10,
			Evidences:     []message.Msg{message.NewLightProposal(oldProposalB10)},
		}

		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 2, len(proofs))
		require.Contains(t, proofs, expectedAccusation1)
		require.Contains(t, proofs, expectedAccusation2)
	})

	t.Run("misbehaviour when pi prevotes for an old proposal while in the valid round there is quorum for different value", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposalB10)
		// Need to add this new proposal in valid round so that unwanted accusation are not returned by the prevotes
		// accountability check method. Since we are adding a quorum of prevotes in round 6 we also need to add a new
		// proposal in round 6 to allow for those prevotes to not return accusations.
		fd.msgStore.Save(newProposalB1In5)
		fd.msgStore.Save(aggregatedPrevoteForOldB10)
		// quorum of prevotes for B1 in vr = 6

		var vr5Prevotes []message.Vote
		for i := uint64(0); i < quorum.Uint64(); i++ {
			vr5Prevote := newValidatedPrevote(5, height, block1.Hash(), makeSigner(keys[i]), &committee.Members[i], cSize)
			vr5Prevotes = append(vr5Prevotes, vr5Prevote)
			fd.msgStore.Save(vr5Prevote)
		}
		aggVr5Votes := message.AggregatePrevotes(vr5Prevotes)

		expectedMisbehaviour1 := &Proof{
			OffenderIndex: proposerIdx,
			Type:          autonity.Misbehaviour,
			Rule:          autonity.PVO,
			Message:       aggregatedPrevoteForOldB10,
		}
		expectedMisbehaviour1.Evidences = append(expectedMisbehaviour1.Evidences, message.NewLightProposal(oldProposalB10))
		expectedMisbehaviour1.Evidences = append(expectedMisbehaviour1.Evidences, aggVr5Votes)

		expectedMisbehaviour2 := &Proof{
			OffenderIndex: prevoterIdx,
			Type:          autonity.Misbehaviour,
			Rule:          autonity.PVO,
			Message:       aggregatedPrevoteForOldB10,
		}
		expectedMisbehaviour2.Evidences = append(expectedMisbehaviour2.Evidences, message.NewLightProposal(oldProposalB10))
		expectedMisbehaviour2.Evidences = append(expectedMisbehaviour2.Evidences, aggVr5Votes)

		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 2, len(proofs))
		actualProof1 := proofs[0]
		require.Equal(t, expectedMisbehaviour1.Type, actualProof1.Type)
		require.Equal(t, expectedMisbehaviour1.OffenderIndex, actualProof1.OffenderIndex)
		require.Equal(t, expectedMisbehaviour1.Rule, actualProof1.Rule)
		require.Equal(t, expectedMisbehaviour1.Message, actualProof1.Message)
		for _, m := range expectedMisbehaviour1.Evidences {
			require.Contains(t, actualProof1.Evidences, m)
		}

		actualProof2 := proofs[1]
		require.Equal(t, expectedMisbehaviour2.Type, actualProof2.Type)
		require.Equal(t, expectedMisbehaviour2.OffenderIndex, actualProof2.OffenderIndex)
		require.Equal(t, expectedMisbehaviour2.Rule, actualProof2.Rule)
		require.Equal(t, expectedMisbehaviour2.Message, actualProof2.Message)
		for _, m := range expectedMisbehaviour2.Evidences {
			require.Contains(t, actualProof2.Evidences, m)
		}
	})

	t.Run("misbehaviour when pi has precommited for V in a previous round however the latest precommit from pi is not for V yet pi still prevoted for V in the current round", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposalB10)
		fd.msgStore.Save(aggregatedPrevoteForOldB10)
		fd.msgStore.Save(newProposalBIn5)

		aggVotes := aggregatedPreVote(committee.Len(), height, 5, oldProposalB10.Value(), keys, committee)
		fd.msgStore.Save(aggVotes)

		var precommitsFromPiAfterLatestPrecommitForB []*message.Precommit
		// create precomits in between the valid round and the current only for proposer node, thus this event is only
		// accountable for propser node. Missing precomits for the other voter, making the event is not accountable for it.
		for i := newProposalBIn5.R(); i < precommitForBIn7.R(); i++ {
			pc := newValidatedPrecommit(i, height, nilValue, signer, self, cSize)
			fd.msgStore.Save(pc)
			if i > oldProposalB10.ValidRound() {
				precommitsFromPiAfterLatestPrecommitForB = append(precommitsFromPiAfterLatestPrecommitForB, pc)
			}
		}

		fd.msgStore.Save(precommitForBIn7)
		precommitsFromPiAfterLatestPrecommitForB = append(precommitsFromPiAfterLatestPrecommitForB, precommitForBIn7)
		fd.msgStore.Save(precommitForB1In8)
		precommitsFromPiAfterLatestPrecommitForB = append(precommitsFromPiAfterLatestPrecommitForB, precommitForB1In8)
		p := newValidatedPrecommit(precommitForB1In8.R()+1, height, nilValue, signer, self, cSize)
		fd.msgStore.Save(p)
		precommitsFromPiAfterLatestPrecommitForB = append(precommitsFromPiAfterLatestPrecommitForB, p)

		aggPrecomits := AggregateDistinctPrecommits(precommitsFromPiAfterLatestPrecommitForB)

		// only the proposer node is accounted for the PVO12 event since the other node does not have the precommits in
		// between the valid round and current round.
		expectedMisbehaviour := &Proof{
			OffenderIndex:      proposerIdx,
			Type:               autonity.Misbehaviour,
			Rule:               autonity.PVO12,
			Message:            aggregatedPrevoteForOldB10,
			DistinctPrecommits: aggPrecomits,
		}
		expectedMisbehaviour.Evidences = append(expectedMisbehaviour.Evidences, message.NewLightProposal(oldProposalB10))

		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 1, len(proofs))
		actualProof := proofs[0]
		require.Equal(t, expectedMisbehaviour.Type, actualProof.Type)
		require.Equal(t, expectedMisbehaviour.Rule, actualProof.Rule)
		require.Equal(t, expectedMisbehaviour.Message, actualProof.Message)
		require.Equal(t, expectedMisbehaviour.DistinctPrecommits, actualProof.DistinctPrecommits)
		err := actualProof.DistinctPrecommits.PreValidate(parent.Epoch.Committee, height)
		require.NoError(t, err)
		err = actualProof.DistinctPrecommits.Validate()
		require.NoError(t, err)

		for _, m := range expectedMisbehaviour.Evidences {
			require.Contains(t, actualProof.Evidences, m)
		}
	})

	t.Run("no proof when pi has precommited for V in a previous round and precommit nils afterwards", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposalB10)
		fd.msgStore.Save(aggregatedPrevoteForOldB10)
		fd.msgStore.Save(newProposalBIn5)

		aggVotes := aggregatedPreVote(committee.Len(), height, 5, block.Hash(), keys, committee)
		fd.msgStore.Save(aggVotes)
		fd.msgStore.Save(precommitForBIn7)
		for i := precommitForBIn7.R() + 1; i < oldProposalB10.R(); i++ {
			v := newValidatedPrecommit(i, height, nilValue, signer, self, cSize)
			fd.msgStore.Save(v)
		}

		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 0, len(proofs))

	})

	t.Run("no proof when pi has precommited for V in a previous round however the latest precommit from pi is not for V yet pi still prevoted for V in the current round"+
		" but there are missing message after latest precommit for V", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposalB10)
		fd.msgStore.Save(aggregatedPrevoteForOldB10)
		fd.msgStore.Save(newProposalBIn5)

		aggVotes := aggregatedPreVote(committee.Len(), height, 5, block.Hash(), keys, committee)
		fd.msgStore.Save(aggVotes)

		fd.msgStore.Save(precommitForBIn7)
		fd.msgStore.Save(precommitForB1In8)

		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 0, len(proofs))

	})

	t.Run("misbehaviour when pi has never precommited for V in a previous round however pi prevoted for V which is being reproposed", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposalB10)
		fd.msgStore.Save(aggregatedPrevoteForOldB10)
		fd.msgStore.Save(newProposalBIn5)
		for i := 0; i < committee.Len(); i++ {
			fd.msgStore.Save(newValidatedPrevote(5, height, block.Hash(), makeSigner(keys[i]), &committee.Members[i], cSize))
		}

		var precomitsFromPiAfterVR1 []*message.Precommit
		for i := newProposalBIn5.R() + 1; i < aggregatedPrecommitForB1In8.R(); i++ {
			p := newValidatedPrecommit(i, height, nilValue, signer, self, cSize)
			fd.msgStore.Save(p)
			precomitsFromPiAfterVR1 = append(precomitsFromPiAfterVR1, p)
		}

		fd.msgStore.Save(aggregatedPrecommitForB1In8)
		precomitsFromPiAfterVR1 = append(precomitsFromPiAfterVR1, aggregatedPrecommitForB1In8)

		p := newValidatedPrecommit(aggregatedPrecommitForB1In8.R()+1, height, nilValue, signer, self, cSize)
		fd.msgStore.Save(p)
		precomitsFromPiAfterVR1 = append(precomitsFromPiAfterVR1, p)

		aggPrecomitsVR1 := AggregateDistinctPrecommits(precomitsFromPiAfterVR1)

		var precommitsFromPiAfterVR2 []*message.Precommit
		for i := newProposalBIn5.R() + 1; i < aggregatedPrecommitForB1In8.R(); i++ {
			p = newValidatedPrecommit(i, height, nilValue, makeSigner(keys[prevoterIdx]), &committee.Members[prevoterIdx], cSize)
			fd.msgStore.Save(p)
			precommitsFromPiAfterVR2 = append(precommitsFromPiAfterVR2, p)
		}

		precommitsFromPiAfterVR2 = append(precommitsFromPiAfterVR2, aggregatedPrecommitForB1In8)

		p = newValidatedPrecommit(aggregatedPrecommitForB1In8.R()+1, height, nilValue, makeSigner(keys[prevoterIdx]),
			&committee.Members[prevoterIdx], cSize)
		fd.msgStore.Save(p)
		precommitsFromPiAfterVR2 = append(precommitsFromPiAfterVR2, p)

		aggPrecommitsVR2 := AggregateDistinctPrecommits(precommitsFromPiAfterVR2)

		expectedMisbehaviour1 := &Proof{
			OffenderIndex:      proposerIdx,
			Type:               autonity.Misbehaviour,
			Rule:               autonity.PVO12,
			DistinctPrecommits: aggPrecomitsVR1,
			Message:            aggregatedPrevoteForOldB10,
		}
		expectedMisbehaviour1.Evidences = append(expectedMisbehaviour1.Evidences, message.NewLightProposal(oldProposalB10))

		expectedMisbehaviour2 := &Proof{
			OffenderIndex:      prevoterIdx,
			Type:               autonity.Misbehaviour,
			Rule:               autonity.PVO12,
			Message:            aggregatedPrevoteForOldB10,
			DistinctPrecommits: aggPrecommitsVR2,
		}
		expectedMisbehaviour2.Evidences = append(expectedMisbehaviour2.Evidences, message.NewLightProposal(oldProposalB10))

		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 2, len(proofs))
		actualProof := proofs[0]
		require.Equal(t, expectedMisbehaviour1.Type, actualProof.Type)
		require.Equal(t, expectedMisbehaviour1.Rule, actualProof.Rule)
		require.Equal(t, expectedMisbehaviour1.Message, actualProof.Message)
		require.Equal(t, expectedMisbehaviour1.DistinctPrecommits, actualProof.DistinctPrecommits)
		for _, m := range expectedMisbehaviour1.Evidences {
			require.Contains(t, actualProof.Evidences, m)
		}
		err := actualProof.DistinctPrecommits.PreValidate(parent.Epoch.Committee, height)
		require.NoError(t, err)
		err = actualProof.DistinctPrecommits.Validate()
		require.NoError(t, err)

		actualProof = proofs[1]
		require.Equal(t, expectedMisbehaviour2.Type, actualProof.Type)
		require.Equal(t, expectedMisbehaviour2.Rule, actualProof.Rule)
		require.Equal(t, expectedMisbehaviour2.Message, actualProof.Message)
		require.Equal(t, expectedMisbehaviour2.DistinctPrecommits, actualProof.DistinctPrecommits)
		for _, m := range expectedMisbehaviour2.Evidences {
			require.Contains(t, actualProof.Evidences, m)
		}
		err = actualProof.DistinctPrecommits.PreValidate(parent.Epoch.Committee, height)
		require.NoError(t, err)
		err = actualProof.DistinctPrecommits.Validate()
		require.NoError(t, err)
	})

	t.Run("no proof when pi has never precommited for V in a previous round however has precommitted nil after VR", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposalB10)
		fd.msgStore.Save(aggregatedPrevoteForOldB10)
		fd.msgStore.Save(newProposalBIn5)

		aggVotes := aggregatedPreVote(committee.Len(), height, 5, block.Hash(), keys, committee)
		fd.msgStore.Save(aggVotes)

		for i := newProposalBIn5.R() + 1; i < oldProposalB10.R(); i++ {
			fd.msgStore.Save(newValidatedPrecommit(i, height, nilValue, signer, self, cSize))
		}

		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when pi has never precommited for V in a previous round however pi prevoted for V while it has precommited for V' but there are missing precommit before precommit for V'", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposalB10)
		fd.msgStore.Save(aggregatedPrevoteForOldB10)
		fd.msgStore.Save(newProposalBIn5)

		aggVotes := aggregatedPreVote(committee.Len(), height, 5, block.Hash(), keys, committee)
		fd.msgStore.Save(aggVotes)

		fd.msgStore.Save(precommitForB1In8)

		p := newValidatedPrecommit(precommitForB1In8.R()+1, height, nilValue, signer, self, cSize)
		fd.msgStore.Save(p)

		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when pi has never precommited for V in a previous round however pi prevoted for V while it has precommited for V' but there are missing precommit after precommit for V'", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(oldProposalB10)

		fd.msgStore.Save(aggregatedPrevoteForOldB10)
		fd.msgStore.Save(newProposalBIn5)
		aggVotes := aggregatedPreVote(committee.Len(), height, 5, block.Hash(), keys, committee)
		fd.msgStore.Save(aggVotes)

		for i := newProposalBIn5.R() + 1; i < precommitForB1In8.R(); i++ {
			fd.msgStore.Save(newValidatedPrecommit(i, height, nilValue, signer, self, cSize))
		}
		fd.msgStore.Save(precommitForB1In8)

		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("prevotes accountability check can return multiple proofs", func(t *testing.T) {
		fd := testFD()

		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(aggregatedPrevoteForB)
		fd.msgStore.Save(aggregatedPrecommitForB1)
		fd.msgStore.Save(aggregatedPrecommitForB)

		fd.msgStore.Save(oldProposalB10)
		fd.msgStore.Save(aggregatedPrevoteForOldB10)

		for i := 0; i < committee.Len(); i++ {
			fd.msgStore.Save(newValidatedPrevote(6, height, block1.Hash(), makeSigner(keys[i]), &committee.Members[i], cSize))
		}

		// Misbehaviour of PVN and Accusation of PVO shall rise to both two nodes, thus we will expect 4 proofs.
		proofs := fd.prevotesAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 4, len(proofs))
	})

	t.Run("no proof when prevote is equivocated with different values", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(aggregatedPrevoteForB)
		fd.msgStore.Save(aggregatedPrevoteForB1)

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

	precommitForB := newValidatedPrecommit(2, height, block.Hash(), signer, self, cSize)
	otherPrecommitForB := newValidatedPrecommit(precommitForB.R(), precommitForB.H(), precommitForB.Value(),
		makeSigner(keys[prevoterIdx]), &committee.Members[prevoterIdx], cSize)
	aggregatedPrecommitForB := message.AggregatePrecommits([]message.Vote{precommitForB, otherPrecommitForB})

	precommitForB1 := newValidatedPrecommit(2, height, block1.Hash(), signer, self, cSize)
	otherPrecommitForB1 := newValidatedPrecommit(2, height, block1.Hash(), makeSigner(keys[prevoterIdx]),
		&committee.Members[prevoterIdx], cSize)
	aggregatedPrecommitForB1 := message.AggregatePrecommits([]message.Vote{precommitForB1, otherPrecommitForB1})

	precommitForB1In3 := newValidatedPrecommit(3, height, block1.Hash(), signer, self, cSize)

	t.Run("accusation when prevotes is less than quorum", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(aggregatedPrecommitForB)

		for i := int64(0); i < quorum.Int64()-1; i++ {
			fd.msgStore.Save(newValidatedPrevote(2, height, block.Hash(), makeSigner(keys[i]), &committee.Members[i], cSize))
		}

		expectedAccusation1 := &Proof{
			OffenderIndex: proposerIdx,
			Type:          autonity.Accusation,
			Rule:          autonity.C1,
			Message:       aggregatedPrecommitForB,
		}

		expectedAccusation2 := &Proof{
			OffenderIndex: prevoterIdx,
			Type:          autonity.Accusation,
			Rule:          autonity.C1,
			Message:       aggregatedPrecommitForB,
		}

		proofs := fd.precommitsAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 2, len(proofs))
		require.Equal(t, expectedAccusation1, proofs[0])
		require.Equal(t, expectedAccusation2, proofs[1])
	})

	t.Run("misbehaviour when there is a quorum for V' than what pi precommitted for", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(aggregatedPrecommitForB)

		preVotesForB1 := aggregatedPreVote(int(quorum.Int64()), height, 2, block1.Hash(), keys, committee)
		fd.msgStore.Save(preVotesForB1)

		expectedMisbehaviour1 := &Proof{
			OffenderIndex: proposerIdx,
			Type:          autonity.Misbehaviour,
			Rule:          autonity.C,
			Message:       aggregatedPrecommitForB,
		}
		expectedMisbehaviour1.Evidences = append(expectedMisbehaviour1.Evidences, preVotesForB1)

		expectedMisbehaviour2 := &Proof{
			OffenderIndex: prevoterIdx,
			Type:          autonity.Misbehaviour,
			Rule:          autonity.C,
			Message:       aggregatedPrecommitForB,
		}
		expectedMisbehaviour2.Evidences = append(expectedMisbehaviour2.Evidences, preVotesForB1)

		proofs := fd.precommitsAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 2, len(proofs))
		require.Equal(t, expectedMisbehaviour1, proofs[0])
		require.Equal(t, expectedMisbehaviour2, proofs[1])
	})

	t.Run("multiple proofs can be returned from precommits accountability check", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(precommitForB1In3)

		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(precommitForB)

		votesForB1 := make([]message.Vote, quorum.Int64())
		for i := int64(0); i < quorum.Int64(); i++ {
			p := newValidatedPrevote(2, height, block1.Hash(), makeSigner(keys[i]), &committee.Members[i], cSize)
			fd.msgStore.Save(p)
			votesForB1[i] = p
		}

		aggVoteForB1 := message.AggregatePrevotes(votesForB1)

		expectedProof0 := &Proof{
			OffenderIndex: proposerIdx,
			Type:          autonity.Misbehaviour,
			Rule:          autonity.C,
			Evidences:     []message.Msg{aggVoteForB1},
			Message:       precommitForB,
		}

		expectedProof1 := &Proof{
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
		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(aggregatedPrecommitForB)

		for i := 0; i < committee.Len(); i++ {
			fd.msgStore.Save(newValidatedPrevote(2, height, block.Hash(), makeSigner(keys[i]), &committee.Members[i], cSize))
		}

		proofs := fd.precommitsAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when there is more than quorum prevotes ", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(newProposalForB)
		fd.msgStore.Save(aggregatedPrecommitForB)

		for i := 0; i < committee.Len(); i++ {
			fd.msgStore.Save(newValidatedPrevote(2, height, block.Hash(), makeSigner(keys[i]), &committee.Members[i], cSize))
		}

		proofs := fd.precommitsAccountabilityCheck(height, quorum, committee)
		require.Equal(t, 0, len(proofs))
	})

	t.Run("no proof when precommit is equivocated with different values", func(t *testing.T) {
		fd := testFD()
		fd.msgStore.Save(aggregatedPrecommitForB)
		fd.msgStore.Save(aggregatedPrecommitForB1)

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

func aggregatedPreVote(numOfSigners int, h uint64, r int64, v common.Hash, keys []blst.SecretKey,
	committee *types.Committee) *message.Prevote {
	var votes []message.Vote
	for i := 0; i < numOfSigners && i < committee.Len(); i++ {
		preVote := newValidatedPrevote(r, h, v, makeSigner(keys[i]), &committee.Members[i], committee.Len())
		votes = append(votes, preVote)
	}
	aggregatedVote := message.AggregatePrevotes(votes)
	return aggregatedVote
}
