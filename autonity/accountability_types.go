package autonity

import "github.com/clearmatics/autonity/common"

// OnChainProof to be stored by autonity contract for on-chain proof management.
type OnChainProof struct {
	Sender  common.Address `abi:"sender"`
	Msghash common.Hash    `abi:"msghash"`
	// rlp enoded bytes for struct Rawproof object.
	Rawproof []byte `abi:"rawproof"`
}

