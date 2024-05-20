package core

import (
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
	"github.com/autonity/autonity/trie"
)

// this file only contains helper function and helper variables for the core tests

var (
	testKey, _          = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	testConsensusKey, _ = blst.SecretKeyFromHex("667e85b8b64622c4b8deadf59964e4c6ae38768a54dbbbc8bbd926777b896584")
	testAddr            = crypto.PubkeyToAddress(testKey.PublicKey)
)

func makeSigner(key blst.SecretKey) message.Signer {
	return func(hash common.Hash) blst.Signature {
		return key.Sign(hash[:])
	}
}

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
}

func generateBlockProposal(r int64, h *big.Int, vr int64, invalid bool, signer message.Signer) *message.Propose {
	var block *types.Block
	if invalid {
		header := &types.Header{Number: h}
		header.Difficulty = nil
		block = types.NewBlock(header, nil, nil, nil, new(trie.Trie))
	} else {
		block = generateBlock(h)
	}
	return message.NewPropose(r, h.Uint64(), vr, block, signer)
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
