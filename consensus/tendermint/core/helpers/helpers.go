package helpers

import (
	"crypto/ecdsa"
	"github.com/autonity/autonity/common"
	tdmcommittee "github.com/autonity/autonity/consensus/tendermint/core/committee"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/rlp"
	"math/big"
	"sort"
)

type AddressKeyMap map[common.Address]*ecdsa.PrivateKey

// PrepareCommittedSeal returns a committed seal for the given hashbytes
func PrepareCommittedSeal(hash common.Hash, round int64, height *big.Int) []byte {
	// this is matching the signature input that we get from the committed messages.
	buf, _ := rlp.EncodeToBytes([]any{message.PrecommitCode, uint64(round), height.Uint64(), hash})
	return buf
}

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
