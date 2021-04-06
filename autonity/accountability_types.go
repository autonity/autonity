package autonity

import (
	"github.com/clearmatics/autonity/common"
	"math/big"
)

// OnChainProof to be stored by autonity contract for on-chain proof management.
type OnChainProof struct {
	Type     *big.Int       `abi:"t"        json:"Type"     gencodec:"required"` // Misbehaviour, Accusation, Innocence to dispatch proof to precompiled contract.
	Sender   common.Address `abi:"sender"   json:"Sender"   gencodec:"required"`
	Msghash  common.Hash    `abi:"msghash"  json:"Msghash"  gencodec:"required"`
	Rawproof []byte         `abi:"rawproof" json:"Rawproof" gencodec:"required"` // rlp enoded bytes for struct Rawproof object.
}
