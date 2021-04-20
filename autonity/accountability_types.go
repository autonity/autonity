package autonity

import (
	"github.com/clearmatics/autonity/common"
)

const (
	Misbehaviour uint8 = iota
	Accusation
	Innocence
)

// OnChainProof to be stored by autonity contract for on-chain proof management.
type OnChainProof struct {
	Type     uint8          `abi:"t"` // Misbehaviour, Accusation, Innocence to dispatch proof to precompiled contract.
	Sender   common.Address `abi:"sender"`
	Msghash  common.Hash    `abi:"msghash"`
	Rawproof []byte         `abi:"rawproof"` // rlp encoded bytes for struct Proof object.
}
