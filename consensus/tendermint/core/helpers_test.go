package core

import (
	"bytes"
	"crypto/ecdsa"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/autonity/autonity/common"
	tdmcommittee "github.com/autonity/autonity/consensus/tendermint/core/committee"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/rlp"
	"github.com/autonity/autonity/trie"
)

// this file only contains helper function and helper variables for the core tests

const minSize, maxSize = 4, 100

var (
	testKey, _          = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	testConsensusKey, _ = blst.SecretKeyFromHex("667e85b8b64622c4b8deadf59964e4c6ae38768a54dbbbc8bbd926777b896584")
	testAddr            = crypto.PubkeyToAddress(testKey.PublicKey)
	testSignatureBytes  = common.Hex2Bytes("8ff38c5915e56029ace231f12e6911587fac4b5618077f3dfe8068138ff1dc7a7ea45a5e0d6a51747cc5f4d990c9d4de1242f4efa93d8165936bfe111f86aaafeea5eda0c38fa3dc2f854576dde63214d7438ea398e48072bc6a0c8e6c2830ef")
	testSignature, _    = blst.SignatureFromBytes(testSignatureBytes)
	testCommitteeMember = &types.CommitteeMember{
		Address:           testAddr,
		ConsensusKey:      testConsensusKey.PublicKey(),
		ConsensusKeyBytes: testConsensusKey.PublicKey().Marshal(),
		VotingPower:       common.Big1,
		Index:             0,
	}
)

func makeSigner(key blst.SecretKey) message.Signer {
	return func(hash common.Hash) blst.Signature {
		return key.Sign(hash[:])
	}
}

func defaultSigner(h common.Hash) blst.Signature {
	return testConsensusKey.Sign(h[:])
}

// creates a signers data structure that carries the requested power
func signersWithPower(index uint64, committeeSize int, requestedPower *big.Int) *types.Signers {
	signers := types.NewSigners(committeeSize)
	fakeMember := &types.CommitteeMember{Index: index, VotingPower: requestedPower}
	signers.Increment(fakeMember)
	return signers
}

func waitForExpects(ctrl *gomock.Controller) {
	for i := 0; i < 200; i++ { // after ten second of waiting, give up and call finish() even if not satisfied
		time.Sleep(50 * time.Millisecond)
		if ctrl.Satisfied() {
			break
		}
	}
	ctrl.Finish()
}

type AddressKeyMap map[common.Address]Keys

type Keys struct {
	consensus blst.SecretKey
	node      *ecdsa.PrivateKey
}

func GenerateCommittee(n int) (*types.Committee, AddressKeyMap) {
	committee := new(types.Committee)
	keymap := make(AddressKeyMap)
	for i := 0; i < n; i++ {
		privateKey, _ := crypto.GenerateKey()
		consensusKey, _ := blst.RandKey()
		committeeMember := types.CommitteeMember{
			Address:           crypto.PubkeyToAddress(privateKey.PublicKey),
			VotingPower:       new(big.Int).SetUint64(1),
			ConsensusKey:      consensusKey.PublicKey(),
			ConsensusKeyBytes: consensusKey.PublicKey().Marshal(),
		}
		committee.Members = append(committee.Members, committeeMember)
		keymap[committeeMember.Address] = Keys{consensus: consensusKey, node: privateKey}
	}
	types.SortCommitteeMembers(committee.Members)
	// assign indexes
	for i := 0; i < n; i++ {
		committee.Members[i].Index = uint64(i)
	}
	return committee, keymap
}

func NewTestCommitteeSet(n int) interfaces.Committee {
	validators, _ := GenerateCommittee(n)
	set, _ := tdmcommittee.NewRoundRobinSet(validators, validators.Members[0].Address)
	return set
}

func NewTestCommitteeSetWithKeys(n int) (interfaces.Committee, AddressKeyMap) {
	validators, keyMap := GenerateCommittee(n)
	set, _ := tdmcommittee.NewRoundRobinSet(validators, validators.Members[0].Address)
	return set, keyMap
}

func generateBlockProposal(r int64, h *big.Int, vr int64, invalid bool, signer message.Signer, member *types.CommitteeMember, parentHeader *types.Header) *message.Propose {
	var block *types.Block
	if invalid {
		header := &types.Header{Number: h, ParentHash: parentHeader.Hash()}
		header.Difficulty = nil
		block = types.NewBlock(header, nil, nil, nil, new(trie.Trie))
	} else {
		block = generateBlock(h, parentHeader)
	}
	return message.NewPropose(r, h.Uint64(), vr, block, signer, member)
}

func randomProposal(t *testing.T) *message.Propose {
	currentHeight := big.NewInt(int64(rand.Intn(100) + 1))
	currentRound := int64(rand.Intn(100) + 1)
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	addr := crypto.PubkeyToAddress(key.PublicKey)
	consensusKey, err := blst.RandKey()
	require.NoError(t, err)
	consensusPubkey := consensusKey.PublicKey()
	committeeMember := &types.CommitteeMember{Address: addr, ConsensusKey: consensusPubkey, ConsensusKeyBytes: consensusPubkey.Marshal(), VotingPower: common.Big1, Index: 0}
	lastHeader := &types.Header{
		Number: big.NewInt(currentHeight.Int64()).Sub(currentHeight, common.Big1),
	}
	return generateBlockProposal(currentRound, currentHeight, currentRound-1, false, makeSigner(consensusKey), committeeMember, lastHeader)
}

// Committee will be ordered such that the proposer for round(n) == committeeSet.members[n % len(committeeSet.members)]
func prepareCommittee(t *testing.T, cSize int) (interfaces.Committee, AddressKeyMap) {
	committee, privateKeys := GenerateCommittee(cSize)
	committeeSet, err := tdmcommittee.NewRoundRobinSet(committee, committee.Members[committee.Len()-1].Address)
	assert.NoError(t, err)
	return committeeSet, privateKeys
}

func generateBlock(height *big.Int, parentHeader *types.Header) *types.Block {
	// use random nonce to create different blocks
	var nonce types.BlockNonce
	for i := 0; i < len(nonce); i++ {
		nonce[i] = byte(rand.Intn(256))
	}
	header := &types.Header{Number: height, Nonce: nonce, ParentHash: parentHeader.Hash()}
	block := types.NewBlockWithHeader(header)
	return block
}

func newUnverifiedPrecommit(r int64, h uint64, value common.Hash, signer message.Signer, self *types.CommitteeMember, csize int) *message.Precommit {
	precommit := message.NewPrecommit(r, h, value, signer, self, csize)
	unverifiedPrecommit := &message.Precommit{}
	reader := bytes.NewReader(precommit.Payload())
	if err := rlp.Decode(reader, unverifiedPrecommit); err != nil {
		panic("cannot decode precommit: " + err.Error())
	}
	return unverifiedPrecommit
}

func newUnverifiedPropose(r int64, h uint64, vr int64, block *types.Block, signer message.Signer, self *types.CommitteeMember) *message.Propose {
	propose := message.NewPropose(r, h, vr, block, signer, self)
	unverifiedPropose := &message.Propose{}
	reader := bytes.NewReader(propose.Payload())
	if err := rlp.Decode(reader, unverifiedPropose); err != nil {
		panic("cannot decode propose: " + err.Error())
	}
	return unverifiedPropose
}
