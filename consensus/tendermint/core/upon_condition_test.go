package core

import (
	"context"
	"errors"
	"math/big"
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/committee"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const minSize, maxSize = 4, 100

// The following tests aim to test lines 1 - 21 of Tendermint Algorithm described on page 6 of
// https://arxiv.org/pdf/1807.04938.pdf.
func TestStartRoundVariables(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	backendMock := NewMockBackend(ctrl)

	prevHeight := big.NewInt(int64(rand.Intn(100) + 1))
	prevBlock := generateBlock(prevHeight)
	currentHeight := big.NewInt(prevHeight.Int64() + 1)
	currentBlock := generateBlock(currentHeight)
	currentRound := int64(0)

	// We don't care who is the next proposer so for simplicity we ensure that clientAddress is not the next
	// proposer by setting clientAddress to be the last proposer. This will ensure that the test will not run the
	// broadcast method.
	committeeSizeAndMaxRound := rand.Intn(maxSize-minSize) + minSize
	committeeSet := prepareCommittee(t, committeeSizeAndMaxRound)
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
	t.Run("ensure round 0 state variables are set correctly", func(t *testing.T) {
		backendMock.EXPECT().Address().Return(clientAddress)
		backendMock.EXPECT().LastCommittedProposal().Return(prevBlock, clientAddress)
		backendMock.EXPECT().Committee(currentHeight.Uint64()).Return(committeeSet, nil)

		core := New(backendMock)
		overrideDefaultCoreValues(core)
		core.startRound(context.Background(), currentRound)

		// Check the initial consensus state
		assert.Equal(t, currentHeight, core.Height())
		assert.Equal(t, currentRound, core.Round())
		assert.Equal(t, propose, core.step)
		assert.Nil(t, core.lockedValue)
		assert.Equal(t, int64(-1), core.lockedRound)
		assert.Nil(t, core.validValue)
		assert.Equal(t, int64(-1), core.validRound)
	})
	t.Run("ensure round x state variables are updated correctly", func(t *testing.T) {
		// In this test we are interested in making sure that that values which change in the current round that may
		// have an impact on the actions performed in the following round (in case of round change) are persisted
		// through to the subsequent round.
		backendMock.EXPECT().Address().Return(clientAddress)
		backendMock.EXPECT().LastCommittedProposal().Return(prevBlock, clientAddress).MaxTimes(2)
		backendMock.EXPECT().Committee(currentHeight.Uint64()).Return(committeeSet, nil).MaxTimes(2)

		core := New(backendMock)
		overrideDefaultCoreValues(core)
		core.startRound(context.Background(), currentRound)

		// Check the initial consensus state
		assert.Equal(t, currentHeight, core.Height())
		assert.Equal(t, currentRound, core.Round())
		assert.Equal(t, propose, core.step)
		assert.Nil(t, core.lockedValue)
		assert.Equal(t, int64(-1), core.lockedRound)
		assert.Nil(t, core.validValue)
		assert.Equal(t, int64(-1), core.validRound)

		// Update locked and valid Value (if locked value changes then valid value also changes, ie quorum(prevotes)
		// delivered in prevote step)
		core.lockedValue = currentBlock
		core.lockedRound = currentRound
		core.validValue = currentBlock
		core.validRound = currentRound

		// Move to next round and check the expected state
		core.startRound(context.Background(), currentRound+1)

		assert.Equal(t, core.Height(), currentHeight)
		assert.Equal(t, currentRound+1, core.Round())
		assert.Equal(t, propose, core.step)
		assert.Equal(t, currentBlock, core.lockedValue)
		assert.Equal(t, currentRound, core.lockedRound)
		assert.Equal(t, currentBlock, core.validValue)
		assert.Equal(t, currentRound, core.validRound)

		// Update valid value (we didn't receive quorum prevote in prevote step, also the block changed, ie, locked
		// value and valid value are different)
		currentBlock2 := generateBlock(currentHeight)
		core.validValue = currentBlock2
		core.validRound = currentRound + 1

		// Move to next round and check the expected state
		core.startRound(context.Background(), currentRound+2)

		assert.Equal(t, core.Height(), currentHeight)
		assert.Equal(t, currentRound+2, core.Round())
		assert.Equal(t, propose, core.step)
		assert.Equal(t, currentBlock, core.lockedValue)
		assert.Equal(t, currentRound, core.lockedRound)
		assert.Equal(t, currentBlock2, core.validValue)
		assert.Equal(t, currentRound+1, core.validRound)
	})
}

func TestStartRound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Committee will be ordered such that the proposer for round(n) == committeeSet.members[n % len(committeeSet.members)]
	committeeSizeAndMaxRound := rand.Intn(maxSize-minSize) + minSize
	committeeSet := prepareCommittee(t, committeeSizeAndMaxRound)
	members := committeeSet.Committee()
	clientAddress := members[0].Address

	t.Run("client is the proposer and valid value is nil", func(t *testing.T) {
		lastBlockProposer := members[len(members)-1].Address
		prevHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		prevBlock := generateBlock(prevHeight)
		proposalHeight := big.NewInt(prevHeight.Int64() + 1)
		proposalBlock := generateBlock(proposalHeight)
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		for currentRound%int64(len(members)) != 0 {
			currentRound = int64(rand.Intn(committeeSizeAndMaxRound))
		}

		proposalMsg, proposalMsgRLPNoSig, proposalMsgRLPWithSig := prepareProposal(t, currentRound, proposalHeight, int64(-1), proposalBlock, clientAddress)

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddress)

		core := New(backendMock)
		// We assume that round 0 can only happen when we move to a new height, therefore, height is
		// incremented by 1 in start round when round = 0, and the committee set is updated. However, in test case where
		// round is more than 0, then we need to explicitly update the committee set and height.
		if currentRound > 0 {
			core.committeeSet = committeeSet
			core.height = proposalHeight
		}
		core.pendingUnminedBlocks[proposalHeight.Uint64()] = proposalBlock

		if currentRound == 0 {
			// We expect the following extra calls when round = 0
			backendMock.EXPECT().LastCommittedProposal().Return(prevBlock, lastBlockProposer)
			backendMock.EXPECT().Committee(proposalHeight.Uint64()).Return(committeeSet, nil)
		}
		backendMock.EXPECT().SetProposedBlockHash(proposalBlock.Hash())
		backendMock.EXPECT().Sign(proposalMsgRLPNoSig).Return(proposalMsg.Signature, nil)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, proposalMsgRLPWithSig).Return(nil)

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

		proposalMsg, proposalMsgRLPNoSig, proposalMsgRLPWithSig := prepareProposal(t, currentRound, proposalHeight, validR, proposalBlock, clientAddress)

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddress)

		core := New(backendMock)
		core.committeeSet = committeeSet
		core.height = proposalHeight
		core.validRound = validR
		core.validValue = proposalBlock

		backendMock.EXPECT().SetProposedBlockHash(proposalBlock.Hash())
		backendMock.EXPECT().Sign(proposalMsgRLPNoSig).Return(proposalMsg.Signature, nil)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, proposalMsgRLPWithSig).Return(nil)

		core.startRound(context.Background(), currentRound)
	})
	t.Run("client is not the proposer", func(t *testing.T) {
		clientIndex := len(members) - 1
		clientAddress := members[clientIndex].Address

		prevHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		prevBlock := generateBlock(prevHeight)
		// ensure the client is not the proposer for current round
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		for currentRound%int64(clientIndex) == 0 {
			currentRound = int64(rand.Intn(committeeSizeAndMaxRound))
		}

		backendMock := NewMockBackend(ctrl)
		backendMock.EXPECT().Address().Return(clientAddress)

		core := New(backendMock)

		if currentRound > 0 {
			core.committeeSet = committeeSet
		}

		if currentRound == 0 {
			backendMock.EXPECT().LastCommittedProposal().Return(prevBlock, clientAddress)
			backendMock.EXPECT().Committee(big.NewInt(prevHeight.Int64()+1)).Return(committeeSet, nil)
		}

		core.startRound(context.Background(), currentRound)

		assert.Equal(t, currentRound, core.Round())
		assert.True(t, core.proposeTimeout.timerStarted())
	})
}

// The following tests aim to test lines 22 - 27 of Tendermint Algorithm described on page 6 of
// https://arxiv.org/pdf/1807.04938.pdf.
func TestTendermintNewProposal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	backendMock := NewMockBackend(ctrl)

	committeeSizeAndMaxRound := rand.Intn(maxSize-minSize) + minSize
	committeeSet := prepareCommittee(t, committeeSizeAndMaxRound)
	members := committeeSet.Committee()
	clientAddr := members[0].Address

	backendMock.EXPECT().Address().Return(clientAddr)
	c := New(backendMock)
	c.setCommitteeSet(committeeSet)

	t.Run("receive invalid proposal for current round", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))

		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setStep(propose)

		var invalidProposal Proposal
		// members[currentRound] means that the sender is the proposer for the current round
		// assume that the message is from a member of committee set and the signature is signing the contents, however,
		// the proposal block inside the message is invalid
		invalidMsg := generateBlockProposal(t, currentRound, currentHeight, members[currentRound].Address, true)
		err := invalidMsg.Decode(&invalidProposal)
		assert.Nil(t, err)

		// prepare prevote nil
		prevoteMsg, prevoteMsgRLPNoSig, prevoteMsgRLPWithSig := preparePrevote(t, currentRound, currentHeight, common.Hash{}, clientAddr)

		backendMock.EXPECT().VerifyProposal(*invalidProposal.ProposalBlock).Return(time.Duration(1), errors.New("invalid proposal"))
		backendMock.EXPECT().Sign(prevoteMsgRLPNoSig).Return(prevoteMsg.Signature, nil)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, prevoteMsgRLPWithSig).Return(nil)

		err = c.handleCheckedMsg(context.Background(), invalidMsg, members[currentRound])
		assert.Error(t, err, "expected an error for invalid proposal")
		assert.Equal(t, prevote, c.step)
	})
	t.Run("receive proposal with validRound = -1 and client's lockedRound = -1", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		clientLockedRound := int64(-1)

		var proposal Proposal
		proposalMsg := generateBlockProposal(t, currentRound, currentHeight, members[currentRound].Address, false)
		err := proposalMsg.Decode(&proposal) // we have to do this because encoding and decoding changes some default values
		assert.Nil(t, err)

		prevoteMsg, prevoteMsgRLPNoSig, prevoteMsgRLPWithSig := preparePrevote(t, currentRound, currentHeight, proposal.ProposalBlock.Hash(), clientAddr)

		// if lockedRround = - 1 then lockedValue = nil
		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setStep(propose)
		c.lockedRound = -1
		c.lockedValue = nil

		backendMock.EXPECT().VerifyProposal(*proposal.ProposalBlock).Return(time.Duration(1), nil)
		backendMock.EXPECT().Sign(prevoteMsgRLPNoSig).Return(prevoteMsg.Signature, nil)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, prevoteMsgRLPWithSig).Return(nil)

		err = c.handleCheckedMsg(context.Background(), proposalMsg, members[currentRound])
		assert.Nil(t, err)
		assert.Equal(t, prevote, c.step)
		assert.Nil(t, c.lockedValue)
		assert.Equal(t, clientLockedRound, c.lockedRound)
	})
	t.Run("receive proposal with validRound = -1 and client's lockedValue is same as proposal block", func(t *testing.T) {
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1))
		currentRound := int64(rand.Intn(committeeSizeAndMaxRound))
		clientLockedRound := int64(0)

		var proposal Proposal
		proposalMsg := generateBlockProposal(t, currentRound, currentHeight, members[currentRound].Address, false)
		// we have to do this because encoding and decoding changes some default values and thus same blocks are no longer equal
		err := proposalMsg.Decode(&proposal)
		assert.Nil(t, err)

		prevoteMsg, prevoteMsgRLPNoSig, prevoteMsgRLPWithSig := preparePrevote(t, currentRound, currentHeight, proposal.ProposalBlock.Hash(), clientAddr)

		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setStep(propose)
		c.lockedRound = clientLockedRound
		c.lockedValue = proposal.ProposalBlock
		c.validRound = clientLockedRound
		c.validValue = proposal.ProposalBlock

		backendMock.EXPECT().VerifyProposal(*proposal.ProposalBlock).Return(time.Duration(1), nil)
		backendMock.EXPECT().Sign(prevoteMsgRLPNoSig).Return(prevoteMsg.Signature, nil)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, prevoteMsgRLPWithSig).Return(nil)

		err = c.handleCheckedMsg(context.Background(), proposalMsg, members[currentRound])
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

		var proposal Proposal
		proposalMsg := generateBlockProposal(t, currentRound, currentHeight, members[currentRound].Address, false)
		// we have to do this because encoding and decoding changes some default values and thus same blocks are no longer equal
		err := proposalMsg.Decode(&proposal)
		assert.Nil(t, err)

		prevoteMsg, prevoteMsgRLPNoSig, prevoteMsgRLPWithSig := preparePrevote(t, currentRound, currentHeight, common.Hash{}, clientAddr)

		c.setHeight(currentHeight)
		c.setRound(currentRound)
		c.setStep(propose)
		c.lockedRound = clientLockedRound
		c.lockedValue = clientLockedValue
		c.validRound = clientLockedRound
		c.validValue = clientLockedValue

		backendMock.EXPECT().VerifyProposal(*proposal.ProposalBlock).Return(time.Duration(1), nil)
		backendMock.EXPECT().Sign(prevoteMsgRLPNoSig).Return(prevoteMsg.Signature, nil)
		backendMock.EXPECT().Broadcast(context.Background(), committeeSet, prevoteMsgRLPWithSig).Return(nil)

		err = c.handleCheckedMsg(context.Background(), proposalMsg, members[currentRound])
		assert.Nil(t, err)
		assert.Equal(t, prevote, c.step)
		assert.Equal(t, clientLockedValue, c.lockedValue)
		assert.Equal(t, clientLockedRound, c.lockedRound)
		assert.Equal(t, clientLockedValue, c.validValue)
		assert.Equal(t, clientLockedRound, c.validRound)
	})
}

// The following tests are not specific to proposal messages but rather apply to all messages
func TestHandleMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	key1, err := crypto.GenerateKey()
	assert.Nil(t, err)
	key2, err := crypto.GenerateKey()
	assert.Nil(t, err)

	key1PubAddr := crypto.PubkeyToAddress(key1.PublicKey)
	key2PubAddr := crypto.PubkeyToAddress(key2.PublicKey)

	committeeSet, err := committee.NewSet(types.Committee{types.CommitteeMember{
		Address:     key1PubAddr,
		VotingPower: big.NewInt(1),
	}}, key1PubAddr)
	assert.Nil(t, err)

	backendMock := NewMockBackend(ctrl)
	backendMock.EXPECT().Address().Return(key1PubAddr)
	core := New(backendMock)

	t.Run("message sender is not in the committee set", func(t *testing.T) {
		// Prepare message
		msg := &Message{Address: key2PubAddr, Code: uint64(rand.Intn(3)), Msg: []byte("random message1")}

		msgRlpNoSig, err := msg.PayloadNoSig()
		assert.Nil(t, err)

		msg.Signature, err = crypto.Sign(crypto.Keccak256(msgRlpNoSig), key2)
		assert.Nil(t, err)

		msgRlpWithSig, err := msg.Payload()
		assert.Nil(t, err)

		core.setCommitteeSet(committeeSet)
		err = core.handleMsg(context.Background(), msgRlpWithSig)

		assert.Error(t, err, "unauthorised sender, sender is not is committees set")
	})

	t.Run("message sender is not the message siger", func(t *testing.T) {
		msg := &Message{Address: key1PubAddr, Code: uint64(rand.Intn(3)), Msg: []byte("random message2")}

		msgRlpNoSig, err := msg.PayloadNoSig()
		assert.Nil(t, err)

		msg.Signature, err = crypto.Sign(crypto.Keccak256(msgRlpNoSig), key1)
		assert.Nil(t, err)

		msgRlpWithSig, err := msg.Payload()
		assert.Nil(t, err)

		core.setCommitteeSet(committeeSet)
		err = core.handleMsg(context.Background(), msgRlpWithSig)

		assert.Error(t, err, "unauthorised sender, sender is not the signer of the message")
	})

	t.Run("malicious sender sends incorrect signature", func(t *testing.T) {
		sig, err := crypto.Sign(crypto.Keccak256([]byte("random bytes")), key1)
		assert.Nil(t, err)

		msg := &Message{Address: key1PubAddr, Code: uint64(rand.Intn(3)), Msg: []byte("random message2"), Signature: sig}
		msgRlpWithSig, err := msg.Payload()
		assert.Nil(t, err)

		core.setCommitteeSet(committeeSet)
		err = core.handleMsg(context.Background(), msgRlpWithSig)

		assert.Error(t, err, "malicious sender sends different signature to signature of message")
	})
}

func prepareProposal(t *testing.T, currentRound int64, proposalHeight *big.Int, validR int64, proposalBlock *types.Block, clientAddress common.Address) (*Message, []byte, []byte) {
	// prepare the proposal message
	proposalRLP, err := Encode(NewProposal(currentRound, proposalHeight, validR, proposalBlock))
	assert.Nil(t, err)
	proposalMsg := &Message{Code: msgProposal, Msg: proposalRLP, Address: clientAddress, Signature: []byte("proposal signature")}
	proposalMsgRLPNoSig, err := proposalMsg.PayloadNoSig()
	assert.Nil(t, err)
	proposalMsgRLPWithSig, err := proposalMsg.Payload()
	assert.Nil(t, err)
	return proposalMsg, proposalMsgRLPNoSig, proposalMsgRLPWithSig
}

func preparePrevote(t *testing.T, round int64, height *big.Int, blockHash common.Hash, clientAddr common.Address) (*Message, []byte, []byte) {
	// prepare the proposal message
	voteRLP, err := Encode(&Vote{Round: round, Height: height, ProposedBlockHash: blockHash})
	assert.Nil(t, err)
	prevoteMsg := &Message{Code: msgPrevote, Msg: voteRLP, Address: clientAddr, Signature: []byte("prevote signature")}
	prevoteMsgRLPNoSig, err := prevoteMsg.PayloadNoSig()
	assert.Nil(t, err)
	prevoteMsgRLPWithSig, err := prevoteMsg.Payload()
	assert.Nil(t, err)
	return prevoteMsg, prevoteMsgRLPNoSig, prevoteMsgRLPWithSig
}

func generateBlockProposal(t *testing.T, r int64, h *big.Int, src common.Address, invalid bool) *Message {
	var block *types.Block
	if invalid {
		header := &types.Header{Number: h}
		header.Difficulty = nil
		block = types.NewBlock(header, nil, nil, nil)
	} else {
		block = generateBlock(h)
	}
	proposal := NewProposal(r, h, -1, block)
	proposalRlp, err := Encode(proposal)
	assert.Nil(t, err)
	return &Message{
		Code:    msgProposal,
		Msg:     proposalRlp,
		Address: src,
	}
}

// Committee will be ordered such that the proposer for round(n) == committeeSet.members[n % len(committeeSet.members)]
func prepareCommittee(t *testing.T, cSize int) *committee.Set {
	c := types.Committee{}
	for i := 1; i <= cSize; i++ {
		key, err := crypto.GenerateKey()
		assert.Nil(t, err)
		member := types.CommitteeMember{
			Address:     crypto.PubkeyToAddress(key.PublicKey),
			VotingPower: new(big.Int).SetInt64(1),
		}
		c = append(c, member)
	}

	sort.Sort(c)
	committeeSet, err := committee.NewSet(c, c[len(c)-1].Address)
	assert.Nil(t, err)
	return committeeSet
}

func generateBlock(height *big.Int) *types.Block {
	// use random nonce to create different blocks
	var nonce types.BlockNonce
	for i := 0; i < len(nonce); i++ {
		nonce[0] = byte(rand.Intn(256))
	}
	header := &types.Header{Number: height, Nonce: nonce}
	block := types.NewBlock(header, nil, nil, nil)
	return block
}
