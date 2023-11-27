package helpers

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/binary"
	"github.com/autonity/autonity/common"
	tdmcommittee "github.com/autonity/autonity/consensus/tendermint/core/committee"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/crypto/bls"
	"math/big"
)

type AddressKeyMap map[common.Address]*ecdsa.PrivateKey

// PrepareCommittedSeal returns a committed seal for the given hashbytes
func PrepareCommittedSeal(hash common.Hash, round int64, height *big.Int) []byte {
	var buf bytes.Buffer
	roundBytes := make([]byte, 8)
	// todo(youssef): endianness seems wrong and the buffer length for the height should be invariant
	binary.LittleEndian.PutUint64(roundBytes, uint64(round))
	buf.Write(roundBytes)
	buf.Write(height.Bytes())
	buf.Write(hash.Bytes())
	return buf.Bytes()
}

// GenerateCommittee is a helper function to generate committee for testing.
func GenerateCommittee(n int) (*types.Committee, AddressKeyMap) {
	// generate committee members data and keymap
	var vals []*types.CommitteeMember
	vals = make([]*types.CommitteeMember, 0)
	keymap := make(AddressKeyMap)
	for i := 0; i < n; i++ {
		privateKey, _ := crypto.GenerateKey()
		blsKey, _ := bls.SecretKeyFromECDSAKey(privateKey)
		committeeMember := &types.CommitteeMember{
			Address:      crypto.PubkeyToAddress(privateKey.PublicKey),
			VotingPower:  new(big.Int).SetUint64(1),
			ValidatorKey: blsKey.PublicKey().Marshal(),
		}
		vals = append(vals, committeeMember)
		keymap[committeeMember.Address] = privateKey
	}
	// wrap members in Committee
	committee := &types.Committee{Members: make([]*types.CommitteeMember, len(vals))}
	for i, m := range vals {
		committee.Members[i] = &types.CommitteeMember{
			Address:      m.Address,
			VotingPower:  new(big.Int).Set(m.VotingPower),
			ValidatorKey: m.ValidatorKey,
		}
	}
	committee.Sort()
	return committee, keymap
}

func NewTestCommitteeSet(n int) interfaces.Committee {
	committee, _ := GenerateCommittee(n)
	set, _ := tdmcommittee.NewRoundRobinSet(committee, committee.Members[0].Address)
	return set
}

func NewTestCommitteeSetWithKeys(n int) (interfaces.Committee, AddressKeyMap) {
	validators, keyMap := GenerateCommittee(n)
	set, _ := tdmcommittee.NewRoundRobinSet(validators, validators.Members[0].Address)
	return set, keyMap
}

func GeneratePrivateKey() (*ecdsa.PrivateKey, error) {
	key := "bb047e5940b6d83354d9432db7c449ac8fca2248008aaa7271369880f9f11cc1"
	return crypto.HexToECDSA(key)
}

func GetAddress() common.Address {
	return common.HexToAddress("0x70524d664ffe731100208a0154e556f9bb679ae6")
}
