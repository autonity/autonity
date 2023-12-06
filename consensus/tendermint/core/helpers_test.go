package core

import (
	"crypto/ecdsa"
	"math/big"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	tdmcommittee "github.com/autonity/autonity/consensus/tendermint/core/committee"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
)

func makeSigner(key *ecdsa.PrivateKey, addr common.Address) message.Signer {
	return func(hash common.Hash) ([]byte, common.Address, *big.Int) {
		out, _ := crypto.Sign(hash[:], key)
		//TODO(lorenzo) fine to return fixed power?
		return out, addr, big.NewInt(1000)
	}
}

func defaultSigner(h common.Hash) ([]byte, common.Address, *big.Int) {
	out, _ := crypto.Sign(h[:], testKey)
	return out, testAddr, testPower
}

type AddressKeyMap map[common.Address]*ecdsa.PrivateKey

func GenerateCommittee(n int) (types.Committee, AddressKeyMap) {
	validators := make(types.Committee, 0)
	keymap := make(AddressKeyMap)
	for i := 0; i < n; i++ {
		privateKey, _ := crypto.GenerateKey()
		committeeMember := types.CommitteeMember{
			Address:     crypto.PubkeyToAddress(privateKey.PublicKey),
			VotingPower: new(big.Int).SetUint64(1),
		}
		validators = append(validators, committeeMember)
		keymap[committeeMember.Address] = privateKey
	}
	sort.Sort(validators)
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
func GeneratePrivateKey() (*ecdsa.PrivateKey, error) {
	key := "bb047e5940b6d83354d9432db7c449ac8fca2248008aaa7271369880f9f11cc1"
	return crypto.HexToECDSA(key)
}

func GetAddress() common.Address {
	return common.HexToAddress("0x70524d664ffe731100208a0154e556f9bb679ae6")
}*/

func TestOverQuorumVotes(t *testing.T) {
	t.Run("with duplicated votes, it returns none duplicated votes of just quorum ones", func(t *testing.T) {
		seats := 10
		committee, _ := GenerateCommittee(seats)
		quorum := bft.Quorum(big.NewInt(int64(seats)))
		height := uint64(1)
		round := int64(0)
		notNilValue := common.Hash{0x1}
		var preVotes []message.Msg
		for i := 0; i < len(committee); i++ {
			preVote := message.NewFakePrevote(message.Fake{
				FakeSender: committee[i].Address,
				FakeRound:  round,
				FakeHeight: height,
				FakeValue:  notNilValue,
				FakePower:  common.Big1,
			})
			preVotes = append(preVotes, preVote)
		}

		// let duplicated msg happens, the counting should skip duplicated ones.
		preVotes = append(preVotes, preVotes...)

		overQuorumVotes := OverQuorumVotes(preVotes, quorum)
		require.Equal(t, quorum.Uint64(), uint64(len(overQuorumVotes)))
	})

	t.Run("with less quorum votes, it returns no votes", func(t *testing.T) {
		seats := 10
		committee, _ := GenerateCommittee(seats)
		quorum := bft.Quorum(new(big.Int).SetInt64(int64(seats)))
		height := uint64(1)
		round := int64(0)
		noneNilValue := common.Hash{0x1}
		var preVotes []message.Msg
		for i := 0; i < int(quorum.Uint64()-1); i++ {
			preVote := message.NewFakePrevote(message.Fake{
				FakeRound:  round,
				FakeHeight: height,
				FakeValue:  noneNilValue,
				FakeSender: committee[i].Address,
				FakePower:  common.Big1,
			})
			preVotes = append(preVotes, preVote)
		}

		// let duplicated msg happens, the counting should skip duplicated ones.
		preVotes = append(preVotes, preVotes...)

		overQuorumVotes := OverQuorumVotes(preVotes, quorum)
		require.Nil(t, overQuorumVotes)
	})
}
