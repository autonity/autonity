package autonity

import (
	"github.com/clearmatics/autonity/common"
	"math/big"
)

type ProofType uint64

const (
	Misbehaviour ProofType = iota
	Accusation
	Innocence
)

type Rule uint8

const (
	PN Rule = iota
	PO
	PVN
	PVO
	C
	C1

	GarbageMessage  // message was signed by valid member, but it cannot be decoded.
	InvalidProposal // The value proposed by proposer cannot pass the blockchain's validation.
	InvalidProposer // A proposal sent from none proposer nodes of the committee.
	Equivocation    // Multiple distinguish votes(proposal, prevote, precommit) sent by validator.
	UnknownRule
)

// OnChainProof to be stored by autonity contract for on-chain proof management.
type OnChainProof struct {
	Type     *big.Int       `abi:"t"` // Misbehaviour, Accusation, Innocence to dispatch proof to precompiled contract.
	Sender   common.Address `abi:"sender"`
	Msghash  common.Hash    `abi:"msghash"`
	Rawproof []byte         `abi:"rawproof"` // rlp enoded bytes for struct Rawproof object.
}
