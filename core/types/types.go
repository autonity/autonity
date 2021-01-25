package types

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/rlp"
	"io"
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

func (p *Proof) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{p.ParentHash, p.Rule, p.Message, p.Evidence})
}

func (p *Proof) DecodeRLP(s *rlp.Stream) error {
	var proof struct {
		ParentHash common.Hash
		Rule uint8
		Message ConsensusMessage
		Evidence []ConsensusMessage
	}
	if err := s.Decode(&proof); err != nil {
		return err
	}

	p.ParentHash, p.Rule, p.Message, p.Evidence = proof.ParentHash, proof.Rule, proof.Message, proof.Evidence
	return nil
}