package types

import (
	"github.com/clearmatics/autonity/common"
)

// The proof used by accountability precompiled contract to validate the proof of innocent or misbehavior.
// Since precompiled contract take raw bytes as input, so it should be rlp encoded into bytes before client send the
// proof into autonity contract.
type Proof struct {
	ParentHash common.Hash        // use by precompiled contract to get committee from chain db.
	Rule       uint8              // rule id.
	Message    ConsensusMessage   // the message to be considered as suspicious one
	Evidence   []ConsensusMessage // the proof of innocent or misbehavior.
}

// todo: implement the RLP encode and decode for Proof.
