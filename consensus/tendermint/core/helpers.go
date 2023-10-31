package core

import (
	"crypto/ecdsa"
	"github.com/autonity/autonity/common"
	tdmcommittee "github.com/autonity/autonity/consensus/tendermint/core/committee"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"math/big"
	"sort"
)

// OverQuorumVotes compute voting power out from a set of prevotes or precommits of a certain round and height, the caller
// should make sure that the votes belong to a certain round and height, it returns a set of votes that the corresponding
// voting power is over quorum, otherwise it returns nil.
func OverQuorumVotes(msgs []message.Message, quorum *big.Int) (overQuorumVotes []message.Message) {
	votingPower := new(big.Int)
	counted := make(map[common.Address]struct{})
	for _, v := range msgs {
		if _, ok := counted[v.Sender()]; ok {
			continue
		}
		counted[v.Sender()] = struct{}{}
		votingPower = votingPower.Add(votingPower, v.Power())
		overQuorumVotes = append(overQuorumVotes, v)
		if votingPower.Cmp(quorum) >= 0 {
			return overQuorumVotes
		}
	}
	return nil
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

func GeneratePrivateKey() (*ecdsa.PrivateKey, error) {
	key := "bb047e5940b6d83354d9432db7c449ac8fca2248008aaa7271369880f9f11cc1"
	return crypto.HexToECDSA(key)
}

func GetAddress() common.Address {
	return common.HexToAddress("0x70524d664ffe731100208a0154e556f9bb679ae6")
}
