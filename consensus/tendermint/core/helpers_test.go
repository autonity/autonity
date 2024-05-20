package core

import (
	"bytes"
	"crypto/ecdsa"
	"math/big"
	"math/rand"
	"sort"

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

var (
	testKey, _          = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	testConsensusKey, _ = blst.SecretKeyFromHex("667e85b8b64622c4b8deadf59964e4c6ae38768a54dbbbc8bbd926777b896584")
	testAddr            = crypto.PubkeyToAddress(testKey.PublicKey)
	testSignatureBytes  = common.Hex2Bytes("8ff38c5915e56029ace231f12e6911587fac4b5618077f3dfe8068138ff1dc7a7ea45a5e0d6a51747cc5f4d990c9d4de1242f4efa93d8165936bfe111f86aaafeea5eda0c38fa3dc2f854576dde63214d7438ea398e48072bc6a0c8e6c2830ef")
	testSignature, _    = blst.SignatureFromBytes(testSignatureBytes)
)

func makeSigner(key blst.SecretKey) message.Signer {
	return func(hash common.Hash) blst.Signature {
		return key.Sign(hash[:])
	}
}

/*
func makeCommitteeMember(address common.Address, key blst.PublicKey, power int64, index uint64) *types.CommitteeMember {
	return &types.CommitteeMember{
		Address:           address,
		VotingPower:       big.NewInt(power),
		ConsensusKey:      key,
		ConsensusKeyBytes: key.Marshal(),
		Index:             index,
	}
}*/

func defaultSigner(h common.Hash) blst.Signature {
	return testConsensusKey.Sign(h[:])
}

/*
func defaultVerifier(address common.Address) *types.CommitteeMember {
	return &types.CommitteeMember{
		Address:      address,
		VotingPower:  common.Big1,
		ConsensusKey: testConsensusKey.PublicKey(),
	}
}*/

type AddressKeyMap map[common.Address]Keys

type Keys struct {
	consensus blst.SecretKey
	node      *ecdsa.PrivateKey
}

func GenerateCommittee(n int) (types.Committee, AddressKeyMap) {
	validators := make(types.Committee, 0)
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
		validators = append(validators, committeeMember)
		keymap[committeeMember.Address] = Keys{consensus: consensusKey, node: privateKey}
	}
	sort.Sort(validators)
	// assign indexes
	for i := 0; i < n; i++ {
		validators[i].Index = uint64(i)
	}
	return validators, keymap
}

func NewTestCommitteeSet(n int) interfaces.Committee {
	validators, _ := GenerateCommittee(n)
	set, _ := tdmcommittee.NewRoundRobinSet(validators, validators[0].Address)
	return set
}

func NewTestCommitteeSetWithKeys(n int) (interfaces.Committee, AddressKeyMap) {
	validators, keyMap := GenerateCommittee(n)
	set, _ := tdmcommittee.NewRoundRobinSet(validators, validators[0].Address)
	return set, keyMap
}

/*
func stubVerifier(consensusKey blst.PublicKey) func(address common.Address) *types.CommitteeMember {
	return func(address common.Address) *types.CommitteeMember {
		return &types.CommitteeMember{
			Address:      address,
			VotingPower:  common.Big1,
			ConsensusKey: consensusKey,
		}
	}
}

func stubVerifierWithPower(consensusKey blst.PublicKey, power int64) func(address common.Address) *types.CommitteeMember {
	return func(address common.Address) *types.CommitteeMember {
		return &types.CommitteeMember{
			Address:      address,
			VotingPower:  big.NewInt(power),
			ConsensusKey: consensusKey,
		}
	}
}*/

func generateBlockProposal(r int64, h *big.Int, vr int64, invalid bool, signer message.Signer, member *types.CommitteeMember) *message.Propose {
	var block *types.Block
	if invalid {
		header := &types.Header{Number: h}
		header.Difficulty = nil
		block = types.NewBlock(header, nil, nil, nil, new(trie.Trie))
	} else {
		block = generateBlock(h)
	}
	return message.NewPropose(r, h.Uint64(), vr, block, signer, member)
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

// locally created messages are considered as verified, we decode it to simulate a msgs arriving from the wire
func newUnverifiedPrevote(r int64, h uint64, value common.Hash, signer message.Signer, self *types.CommitteeMember, csize int) *message.Prevote {
	prevote := message.NewPrevote(r, h, value, signer, self, csize)
	unverifiedPrevote := &message.Prevote{}
	reader := bytes.NewReader(prevote.Payload())
	if err := rlp.Decode(reader, unverifiedPrevote); err != nil {
		panic("cannot decode prevote: " + err.Error())
	}
	return unverifiedPrevote
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
