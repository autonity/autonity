package faultdetector

import (
	"github.com/clearmatics/autonity/autonity"
	"github.com/clearmatics/autonity/consensus/tendermint/core"
	"github.com/clearmatics/autonity/rlp"
	"io"
)

const (
	msgProposal uint64 = iota
	msgPrevote
	msgPrecommit
)

// The proof used by accountability precompiled contract to validate the proof of innocent or misbehavior.
// Since precompiled contract take raw bytes as input, so it should be rlp encoded into bytes before client send the
// proof into autonity contract.
type RawProof struct {
	Rule     autonity.Rule // rule id.
	Message  []byte        // the raw rlp encoded msg to be considered as suspicious one
	Evidence [][]byte      // the raw rlp encoded msgs as proof of innocent or misbehavior.
}

func (p *RawProof) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{p.Rule, p.Message, p.Evidence})
}

func (p *RawProof) DecodeRLP(s *rlp.Stream) error {
	var proof struct {
		Rule     autonity.Rule
		Message  []byte
		Evidence [][]byte
	}
	if err := s.Decode(&proof); err != nil {
		return err
	}

	p.Rule, p.Message, p.Evidence = proof.Rule, proof.Message, proof.Evidence
	return nil
}

// Proof is what to prove that one is misbehaving, one should be slashed when a valid proof is rise.
type Proof struct {
	Type     autonity.ProofType // Misbehaviour, Accusation, Innocence.
	Rule     autonity.Rule
	Message  *core.Message   // the msg to be considered as suspicious or misbehaved one
	Evidence []*core.Message // the proofs of innocence or misbehaviour.
}

// event to submit proofs via standard transaction.
type AccountabilityEvent struct {
	Proofs []autonity.OnChainProof
}
