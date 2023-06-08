package autonity

import (
	"github.com/autonity/autonity/common"
)

type Rule uint8

const (
	_ Rule = iota
	PN
	PO
	PVN
	PVO
	PVO12 // rule that encompasses both PVO1 rule and PVO2 from D3 paper
	PVO3
	C
	C1

	InvalidProposal // The value proposed by proposer cannot pass the blockchain's validation.
	InvalidProposer // A proposal sent from none proposer nodes of the committee.
	Equivocation    // Multiple distinguish votes(proposal, prevote, precommit) sent by validator.

	InvalidRound              // message contains invalid round
	WrongValidRound           // proposal contains a wrong valid round number.
	AccountableGarbageMessage // message was signed by sender, but it cannot be decoded.
)

type AccountabilityEventType uint8

const (
	Misbehaviour AccountabilityEventType = iota
	Accusation
	Innocence
)

// AccountabilityEvent to be handled by autonity contract for on-chain accountability management.
type AccountabilityEvent struct {
	Chunks   uint8          `abi:"Chunks"`   // Counter of number of chunks for oversize accountability event
	ChunkID  uint8          `abi:"ChunkID"`  // Chunk index to construct the oversize accountability event
	Type     uint8          `abi:"Type"`     // Accountability event types: Misbehaviour, Accusation, Innocence.
	Rule     uint8          `abi:"Rule"`     // Rule ID defined in AFD rule engine.
	Reporter common.Address `abi:"Reporter"` // The node address of the validator who report this event, for incentive protocol.
	Sender   common.Address `abi:"Sender"`   // The corresponding node address of this accountability event.
	MsgHash  common.Hash    `abi:"MsgHash"`  // The corresponding consensus msg's hash of this accountability event.
	RawProof []byte         `abi:"RawProof"` // rlp encoded bytes of Proof object.
}
