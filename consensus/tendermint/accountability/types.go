package accountability

import (
	"fmt"
	"io"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/rlp"
)

type Proof struct {
	Type               autonity.AccountabilityEventType // Accountability event types: Misbehaviour, Accusation, Innocence.
	Rule               autonity.Rule                    // Rule ID defined in AFD rule engine.
	Message            message.Msg                      // the consensus message which is accountable.
	Evidences          []message.Msg                    // the proofs of the accountability event.
	DistinctPrecommits HighlyAggregatedPrecommit        // the DistinctPrecommits contains highly aggregated precommits.
	OffenderIndex      int                              // the offender index.
}

type encodedProof struct {
	Type               autonity.AccountabilityEventType
	Rule               autonity.Rule
	OffenderIndex      uint
	Message            message.TypedMessage
	Evidences          []message.TypedMessage
	DistinctPrecommits HighlyAggregatedPrecommit
}

func (p *Proof) EncodeRLP(w io.Writer) error {
	encoded := encodedProof{
		Type:               p.Type,
		Rule:               p.Rule,
		OffenderIndex:      uint(p.OffenderIndex),
		DistinctPrecommits: p.DistinctPrecommits,
	}
	encoded.Message = message.TypedMessage{Msg: p.Message}
	for _, m := range p.Evidences {
		encoded.Evidences = append(encoded.Evidences, message.TypedMessage{Msg: m})
	}
	return rlp.Encode(w, &encoded)
}

func (p *Proof) DecodeRLP(stream *rlp.Stream) error {
	encoded := encodedProof{}
	if err := stream.Decode(&encoded); err != nil {
		return fmt.Errorf("could not decode encoded proof %w", err)
	}
	p.Type = encoded.Type
	p.Rule = encoded.Rule
	p.OffenderIndex = int(encoded.OffenderIndex)
	p.Message = encoded.Message.Msg
	p.DistinctPrecommits = encoded.DistinctPrecommits

	p.Evidences = make([]message.Msg, len(encoded.Evidences))
	for i := range encoded.Evidences {
		p.Evidences[i] = encoded.Evidences[i].Msg
	}
	return nil
}
