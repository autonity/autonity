package accountability

import (
	ethereum "github.com/autonity/autonity"
	"github.com/autonity/autonity/accounts/abi/bind/backends"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/backend"
	"github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	ccore "github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/rlp"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"math/big"
	"testing"
	"time"
)

func TestNewOffChainAccusationRateLimiter(t *testing.T) {
	msgSender := common.Address{}
	msgHash1 := common.Hash{0x1}
	msgHash2 := common.Hash{0x2}
	t.Run("test rate limit with a 1st accusation", func(t *testing.T) {
		rl := NewAccusationRateLimiter()
		err := rl.validAccusationRate(msgSender)
		require.NoError(t, err)
		require.Equal(t, 1, rl.accusationRates[msgSender])

		rl.resetRateLimiter()
		require.Equal(t, 0, len(rl.accusationRates))
	})

	t.Run("test rate limit with limited rate", func(t *testing.T) {
		rl := NewAccusationRateLimiter()
		for i := 0; i < maxAccusationRatePerHeight*2; i++ {
			err := rl.validAccusationRate(msgSender)
			require.NoError(t, err)
		}
		err := rl.validAccusationRate(msgSender)
		require.Error(t, errAccusationRateMalicious, err)

		rl.resetRateLimiter()
		require.Equal(t, 0, len(rl.accusationRates))
	})

	t.Run("test duplicated accusation", func(t *testing.T) {
		rl := NewAccusationRateLimiter()
		err := rl.checkPeerDuplicatedAccusation(msgSender, msgHash1)
		require.NoError(t, err)
		_, ok := rl.peerProcessedAccusations[msgSender][msgHash1]
		require.Equal(t, true, ok)
		err = rl.checkPeerDuplicatedAccusation(msgSender, msgHash1)
		require.Error(t, errPeerDuplicatedAccusation, err)
		err = rl.checkPeerDuplicatedAccusation(msgSender, msgHash2)
		require.NoError(t, err)

		rl.resetPeerJustifiedAccusations()
		_, ok = rl.peerProcessedAccusations[msgSender][msgHash1]
		require.Equal(t, false, ok)
		_, ok = rl.peerProcessedAccusations[msgSender][msgHash2]
		require.Equal(t, false, ok)
	})

	t.Run("test accusation rate limit over a height", func(t *testing.T) {
		rl := NewAccusationRateLimiter()

		for h := uint64(0); h < uint64(99); h++ {
			for i := 0; i < maxAccusationRatePerHeight; i++ {
				err := rl.checkHeightAccusationRate(msgSender, h)
				require.NoError(t, err)
			}
			err := rl.checkHeightAccusationRate(msgSender, h)
			require.Error(t, errAccusationRateMalicious, err)

			rl.resetHeightRateLimiter()
			err = rl.checkHeightAccusationRate(msgSender, h)
			require.NoError(t, err)
		}
	})
}

func TestNewInnocenceProofBuffer(t *testing.T) {
	t.Run("cache and get innocence proof", func(t *testing.T) {
		c := NewInnocenceProofBuffer()
		rawPayload := make([]byte, 128)
		hash := types.RLPHash(rawPayload)
		c.cacheInnocenceProof(hash, rawPayload)
		ret := c.getInnocenceProofFromCache(hash)
		require.Equal(t, rawPayload, ret)
		require.Equal(t, 1, len(c.proofs))
		require.Equal(t, 1, len(c.accusationList))
		require.Equal(t, hash, c.accusationList[0])
	})

	t.Run("cache innocence proof with LRU swap", func(t *testing.T) {
		c := NewInnocenceProofBuffer()
		for i := 0; i < maxNumOfInnocenceProofCached*4; i++ {
			rawPayload := make([]byte, i+1)
			hash := types.RLPHash(rawPayload)
			c.cacheInnocenceProof(hash, rawPayload)
			ret := c.getInnocenceProofFromCache(hash)
			require.Equal(t, rawPayload, ret)
		}

		// the swap out one should no longer in the cache.
		swapOut := make([]byte, 1)
		swapHash := types.RLPHash(swapOut)
		ret := c.getInnocenceProofFromCache(swapHash)
		require.Equal(t, []byte(nil), ret)

		require.Equal(t, maxNumOfInnocenceProofCached, len(c.proofs))
		require.Equal(t, maxNumOfInnocenceProofCached, len(c.accusationList))
	})
}

func TestFaultDetector_sendOffChainInnocenceProof(t *testing.T) {
	clientAddr := common.Address{}
	sk, _ := crypto.GenerateKey()
	remotePeer := common.Address{0x1}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chainMock := NewMockChainContext(ctrl)
	var blockSub event.Subscription
	chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
	chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})
	accountability, _ := autonity.NewAccountability(clientAddr, backends.NewSimulatedBackend(ccore.GenesisAlloc{clientAddr: {Balance: big.NewInt(params.Ether)}}, 10000000))

	fd := NewFaultDetector(chainMock, clientAddr, nil, nil, nil, nil, sk, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())
	broadcasterMock := consensus.NewMockBroadcaster(ctrl)
	fd.SetBroadcaster(broadcasterMock)

	payload := make([]byte, 128)

	targets := make(map[common.Address]struct{})
	targets[remotePeer] = struct{}{}
	mockedPeer := ethereum.NewMockPeer(ctrl)
	mockedPeer.EXPECT().Send(uint64(backend.AccountabilityMsg), payload).MaxTimes(1)
	peers := make(map[common.Address]ethereum.Peer)
	peers[remotePeer] = mockedPeer
	broadcasterMock.EXPECT().FindPeers(targets).Return(peers)
	fd.sendOffChainInnocenceProof(remotePeer, payload)
	// wait for msg send routine to be terminated.
	<-time.NewTimer(2 * time.Second).C
}

func TestFaultDetector_sendOffChainAccusationMsg(t *testing.T) {
	clientAddr := common.Address{}
	remotePeer := common.Address{0x1}
	sk, _ := crypto.GenerateKey()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chainMock := NewMockChainContext(ctrl)
	var blockSub event.Subscription
	chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
	chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})
	accountability, _ := autonity.NewAccountability(clientAddr, backends.NewSimulatedBackend(ccore.GenesisAlloc{clientAddr: {Balance: big.NewInt(params.Ether)}}, 10000000))

	fd := NewFaultDetector(chainMock, clientAddr, nil, nil, nil, nil, sk, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())

	broadcasterMock := consensus.NewMockBroadcaster(ctrl)
	fd.SetBroadcaster(broadcasterMock)

	var proposal = message.Message{
		Address: remotePeer,
	}
	var accusation = Proof{
		Type:      autonity.Accusation,
		Rule:      autonity.PO,
		Message:   &proposal,
		Evidences: nil,
	}
	payload, err := rlp.EncodeToBytes(&accusation)
	require.NoError(t, err)

	targets := make(map[common.Address]struct{})
	targets[remotePeer] = struct{}{}
	mockedPeer := ethereum.NewMockPeer(ctrl)
	mockedPeer.EXPECT().Send(uint64(backend.AccountabilityMsg), payload).MaxTimes(1)
	peers := make(map[common.Address]ethereum.Peer)
	peers[remotePeer] = mockedPeer
	broadcasterMock.EXPECT().FindPeers(targets).Return(peers)
	fd.sendOffChainAccusationMsg(&accusation)
	// wait for msg send routine to be terminated.
	<-time.NewTimer(2 * time.Second).C
}

func TestOffChainAccusationManagement(t *testing.T) {
	clientAddr := common.Address{}
	sk, _ := crypto.GenerateKey()
	remotePeer := common.Address{0x1}
	t.Run("Add off chain accusation", func(t *testing.T) {
		var proposal = message.Message{
			Address: remotePeer,
		}
		var accusation = Proof{
			Type:      autonity.Accusation,
			Rule:      autonity.PO,
			Message:   &proposal,
			Evidences: nil,
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		chainMock := NewMockChainContext(ctrl)
		var blockSub event.Subscription
		chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
		chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})
		accountability, _ := autonity.NewAccountability(clientAddr, backends.NewSimulatedBackend(ccore.GenesisAlloc{clientAddr: {Balance: big.NewInt(params.Ether)}}, 10000000))

		fd := NewFaultDetector(chainMock, clientAddr, nil, nil, nil, nil, sk, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())

		fd.addOffChainAccusation(&accusation)
		require.Equal(t, 1, len(fd.offChainAccusations))
	})

	t.Run("remove off chain accusation", func(t *testing.T) {
		var proposal = message.Message{
			Address: remotePeer,
		}
		var accusationPO = Proof{
			Type:      autonity.Accusation,
			Rule:      autonity.PO,
			Message:   &proposal,
			Evidences: nil,
		}

		var preCommit = message.Message{
			Address: remotePeer,
		}
		var accusationC1 = Proof{
			Type:      autonity.Accusation,
			Rule:      autonity.C1,
			Message:   &preCommit,
			Evidences: nil,
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		chainMock := NewMockChainContext(ctrl)
		var blockSub event.Subscription
		chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
		chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})
		accountability, _ := autonity.NewAccountability(clientAddr, backends.NewSimulatedBackend(ccore.GenesisAlloc{clientAddr: {Balance: big.NewInt(params.Ether)}}, 10000000))

		fd := NewFaultDetector(chainMock, clientAddr, nil, nil, nil, nil, sk, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())

		fd.addOffChainAccusation(&accusationPO)
		fd.addOffChainAccusation(&accusationC1)

		require.Equal(t, 2, len(fd.offChainAccusations))

		var innocenceProof = Proof{
			Type:    autonity.Accusation,
			Rule:    autonity.PO,
			Message: &proposal,
		}
		fd.removeOffChainAccusation(&innocenceProof)
		require.Equal(t, 1, len(fd.offChainAccusations))
	})

	t.Run("get expired off chain accusation", func(t *testing.T) {
		currentHeight := uint64(31)
		msgHeight := uint64(10)
		msgRound := int64(1)
		validRound := int64(0)

		committee, keys := generateCommittee()
		proposer := committee[0].Address
		proposerKey := keys[proposer]

		proposal := newProposalMessage(msgHeight, msgRound, validRound, proposerKey, committee, nil)
		var accusationPO = Proof{
			Type:      autonity.Accusation,
			Rule:      autonity.PO,
			Message:   proposal,
			Evidences: nil,
		}

		preCommit := newVoteMsg(msgHeight, msgRound, consensus.MsgPrecommit, proposerKey, nilValue, committee)
		var accusationC1 = Proof{
			Type:      autonity.Accusation,
			Rule:      autonity.C1,
			Message:   preCommit,
			Evidences: nil,
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		chainMock := NewMockChainContext(ctrl)
		var blockSub event.Subscription
		chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
		chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})
		accountability, _ := autonity.NewAccountability(clientAddr, backends.NewSimulatedBackend(ccore.GenesisAlloc{clientAddr: {Balance: big.NewInt(params.Ether)}}, 10000000))

		fd := NewFaultDetector(chainMock, clientAddr, nil, nil, nil, nil, sk, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())

		fd.addOffChainAccusation(&accusationPO)
		fd.addOffChainAccusation(&accusationC1)
		expires := fd.getExpiredOffChainAccusation(currentHeight)
		require.Equal(t, 2, len(expires))
		require.Equal(t, 2, len(fd.offChainAccusations))
	})

	t.Run("escalateExpiredAccusations", func(t *testing.T) {
		currentHeight := uint64(31)
		msgHeight := uint64(10)
		msgRound := int64(1)
		validRound := int64(0)

		committee, keys := generateCommittee()
		proposer := committee[0].Address
		proposerKey := keys[proposer]

		proposal := newProposalMessage(msgHeight, msgRound, validRound, proposerKey, committee, nil)
		var accusationPO = Proof{
			Type:      autonity.Accusation,
			Rule:      autonity.PO,
			Message:   proposal,
			Evidences: nil,
		}

		preCommit := newVoteMsg(msgHeight, msgRound, consensus.MsgPrecommit, proposerKey, nilValue, committee)
		var accusationC1 = Proof{
			Type:      autonity.Accusation,
			Rule:      autonity.C1,
			Message:   preCommit,
			Evidences: nil,
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		chainMock := NewMockChainContext(ctrl)
		var blockSub event.Subscription
		chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
		chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})
		accountability, _ := autonity.NewAccountability(clientAddr, backends.NewSimulatedBackend(ccore.GenesisAlloc{clientAddr: {Balance: big.NewInt(params.Ether)}}, 10000000))

		fd := NewFaultDetector(chainMock, clientAddr, nil, nil, nil, nil, sk, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())

		fd.addOffChainAccusation(&accusationPO)
		fd.addOffChainAccusation(&accusationC1)

		fd.escalateExpiredAccusations(currentHeight)
		require.Equal(t, 0, len(fd.offChainAccusations))
		require.Equal(t, 2, len(fd.pendingEvents))
	})
}

func TestHandleOffChainAccountabilityEvent(t *testing.T) {
	committee, keys := generateCommittee()
	proposer := committee[0].Address
	sender := committee[1].Address
	proposerKey := keys[proposer]
	height := uint64(100)
	round := int64(1)
	validRound := int64(0)
	lastHeight := height - 1
	t.Run("malicious accusation with duplicated msg", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chainMock := NewMockChainContext(ctrl)
		var blockSub event.Subscription
		chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
		ms := core.NewMsgStore()
		chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})
		accountability, _ := autonity.NewAccountability(proposer, backends.NewSimulatedBackend(ccore.GenesisAlloc{proposer: {Balance: big.NewInt(params.Ether)}}, 10000000))

		fd := NewFaultDetector(chainMock, proposer, nil, ms, nil, nil, proposerKey, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())

		proposal := newProposalMessage(height, round, validRound, proposerKey, committee, nil)
		var accusationPO = Proof{
			Type:      autonity.Accusation,
			Rule:      autonity.PO,
			Message:   proposal.ToLightProposal(),
			Evidences: nil,
		}

		payLoad, err := rlp.EncodeToBytes(&accusationPO)
		require.NoError(t, err)

		lastHeader := newBlockHeader(lastHeight, committee)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(lastHeader).AnyTimes()
		chainMock.EXPECT().CurrentHeader().Return(lastHeader).AnyTimes()

		for _, c := range committee {
			preVote := newVoteMsg(height, validRound, consensus.MsgPrevote, keys[c.Address], proposal.Value(), committee)
			ms.Save(preVote)
		}

		//hash := types.RLPHash(payLoad)
		for i := 0; i < 200; i++ {
			err = fd.handleOffChainAccountabilityEvent(payLoad, sender)
			if err != nil {
				break
			}
		}
		require.Equal(t, 1, len(fd.innocenceProofBuff.accusationList))
		require.Equal(t, errPeerDuplicatedAccusation, err)
	})

	t.Run("accusation is not from committee member", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chainMock := NewMockChainContext(ctrl)
		var blockSub event.Subscription
		chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
		ms := core.NewMsgStore()
		chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})
		accountability, _ := autonity.NewAccountability(sender, backends.NewSimulatedBackend(ccore.GenesisAlloc{sender: {Balance: big.NewInt(params.Ether)}}, 10000000))

		fd := NewFaultDetector(chainMock, sender, nil, core.NewMsgStore(), nil, nil, proposerKey, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())

		proposal := newProposalMessage(height, round, validRound, proposerKey, committee, nil)
		var accusationPO = Proof{
			Type:      autonity.Accusation,
			Rule:      autonity.PO,
			Message:   proposal,
			Evidences: nil,
		}

		payLoad, err := rlp.EncodeToBytes(&accusationPO)
		require.NoError(t, err)

		lastHeader := newBlockHeader(lastHeight, committee)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(lastHeader).AnyTimes()
		chainMock.EXPECT().CurrentHeader().Return(lastHeader).AnyTimes()

		for _, c := range committee {
			preVote := newVoteMsg(height, validRound, consensus.MsgPrevote, keys[c.Address], proposal.Value(), committee)
			ms.Save(preVote)
		}

		maliciousSender := common.Address{}
		err = fd.handleOffChainAccountabilityEvent(payLoad, maliciousSender)
		require.Equal(t, errAccusationFromNoneValidator, err)
	})
}

func TestHandleOffChainAccusation(t *testing.T) {
	committee, keys := generateCommittee()
	proposer := committee[0].Address
	sender := committee[1].Address
	proposerKey := keys[proposer]
	height := uint64(100)
	round := int64(1)
	validRound := int64(0)
	lastHeight := height - 1
	t.Run("accusation have invalid proof of wrong signature", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		invalidCommittee, iKeys := generateCommittee()
		chainMock := NewMockChainContext(ctrl)
		var blockSub event.Subscription
		chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
		chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})
		accountability, _ := autonity.NewAccountability(sender, backends.NewSimulatedBackend(ccore.GenesisAlloc{sender: {Balance: big.NewInt(params.Ether)}}, 10000000))

		fd := NewFaultDetector(chainMock, sender, nil, core.NewMsgStore(), nil, nil, proposerKey, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())

		var p Proof
		p.Rule = autonity.PO
		p.Type = autonity.Accusation
		invalidProposal := newProposalMessage(height, 1, 0, iKeys[invalidCommittee[0].Address], invalidCommittee, nil)
		p.Message = invalidProposal.ToLightProposal()
		payload, err := rlp.EncodeToBytes(p)
		require.NoError(t, err)
		hash := types.RLPHash(payload)

		err = fd.handleOffChainAccusation(&p, common.Address{}, hash)
		require.Equal(t, errInvalidAccusation, err)
	})

	t.Run("happy case with innocence proof collected from msg store", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		proposal := newProposalMessage(height, round, validRound, proposerKey, committee, nil)
		var accusationPO = Proof{
			Type:      autonity.Accusation,
			Rule:      autonity.PO,
			Message:   proposal.ToLightProposal(),
			Evidences: nil,
		}

		payLoad, err := rlp.EncodeToBytes(&accusationPO)
		require.NoError(t, err)
		hash := types.RLPHash(payLoad)

		chainMock := NewMockChainContext(ctrl)
		var blockSub event.Subscription
		chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
		mStore := core.NewMsgStore()
		chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})
		accountability, _ := autonity.NewAccountability(proposer, backends.NewSimulatedBackend(ccore.GenesisAlloc{proposer: {Balance: big.NewInt(params.Ether)}}, 10000000))

		fd := NewFaultDetector(chainMock, proposer, nil, mStore, nil, nil, proposerKey, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())

		// save corresponding prevotes in msg store.
		for _, c := range committee {
			preVote := newVoteMsg(height, validRound, consensus.MsgPrevote, keys[c.Address], proposal.Value(), committee)
			mStore.Save(preVote)
		}

		lastHeader := newBlockHeader(lastHeight, committee)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(lastHeader).AnyTimes()
		chainMock.EXPECT().CurrentHeader().Return(lastHeader).AnyTimes()

		err = fd.handleOffChainAccusation(&accusationPO, sender, hash)
		require.NoError(t, err)

		require.Equal(t, 1, len(fd.innocenceProofBuff.accusationList))
	})
}

func TestHandleOffChainProofOfInnocence(t *testing.T) {
	committee, keys := generateCommittee()
	proposer := committee[0].Address
	proposerKey := keys[proposer]
	height := uint64(100)
	round := int64(1)
	validRound := int64(0)
	lastHeight := height - 1
	t.Run("innocence proof is invalid without any evidence", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		chainMock := NewMockChainContext(ctrl)
		var blockSub event.Subscription
		chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
		chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})
		accountability, _ := autonity.NewAccountability(proposer, backends.NewSimulatedBackend(ccore.GenesisAlloc{proposer: {Balance: big.NewInt(params.Ether)}}, 10000000))

		fd := NewFaultDetector(chainMock, proposer, nil, core.NewMsgStore(), nil, nil, proposerKey, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())
		var p Proof
		p.Rule = autonity.PO
		p.Type = autonity.Innocence
		invalidCommittee, iKeys := generateCommittee()
		invalidProposal := newProposalMessage(height, 1, 0, iKeys[invalidCommittee[0].Address], invalidCommittee, nil)
		p.Message = invalidProposal

		lastHeader := newBlockHeader(lastHeight, committee)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(lastHeader)

		err := fd.handleOffChainProofOfInnocence(&p, invalidCommittee[0].Address)
		require.Equal(t, errInvalidInnocenceProof, err)
	})

	t.Run("happy case", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		proposal := newProposalMessage(height, round, validRound, proposerKey, committee, nil)
		var accusationPO = Proof{
			Type:      autonity.Accusation,
			Rule:      autonity.PO,
			Message:   proposal.ToLightProposal(),
			Evidences: nil,
		}

		chainMock := NewMockChainContext(ctrl)
		var blockSub event.Subscription
		chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
		chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})
		accountability, _ := autonity.NewAccountability(proposer, backends.NewSimulatedBackend(ccore.GenesisAlloc{proposer: {Balance: big.NewInt(params.Ether)}}, 10000000))

		fd := NewFaultDetector(chainMock, proposer, nil, core.NewMsgStore(), nil, nil, proposerKey, &autonity.ProtocolContracts{Accountability: accountability}, log.Root())
		// add accusation in fd first.
		fd.addOffChainAccusation(&accusationPO)

		var proofPO = Proof{
			Type:    autonity.Innocence,
			Rule:    autonity.PO,
			Message: proposal.ToLightProposal(),
		}

		// handle a valid innocence proof then.
		for _, c := range committee {
			preVote := newVoteMsg(height, validRound, consensus.MsgPrevote, keys[c.Address], proposal.Value(), committee)
			proofPO.Evidences = append(proofPO.Evidences, preVote)
		}

		lastHeader := newBlockHeader(lastHeight, committee)
		chainMock.EXPECT().GetHeaderByNumber(lastHeight).Return(lastHeader).AnyTimes()

		err := fd.handleOffChainProofOfInnocence(&proofPO, proposer)

		require.NoError(t, err)
		require.Equal(t, 0, len(fd.offChainAccusations))
	})
}
