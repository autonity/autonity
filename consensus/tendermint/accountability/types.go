package accountability

import (
	"errors"
	"fmt"
	"github.com/autonity/autonity/common"
	"io"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/rlp"
)

var (
	errUnexpectedCode = errors.New("unexpected message code")
	errTypeMissing    = errors.New("message code is missing")
	errFullProposal   = errors.New("only light proposals are supported")
)

type typedMessage struct {
	message.Msg
}

func (t *typedMessage) EncodeRLP(w io.Writer) error {
	if t.Msg.Code() == message.ProposalCode {
		return errFullProposal
	}
	return rlp.Encode(w, []any{t.Code(), t.Msg})
}

func (t *typedMessage) DecodeRLP(stream *rlp.Stream) error {
	// Getting back the original payload directly from this stream object
	// is currently not supported.
	if _, err := stream.List(); err != nil {
		return err
	}
	b, err := stream.Bytes()
	if err != nil {
		return err
	}
	if len(b) == 0 {
		return errTypeMissing
	}
	// Note that a nil element is not accepted.
	var p message.Msg
	switch b[0] {
	case message.PrevoteCode:
		p = &message.Prevote{}
	case message.PrecommitCode:
		p = &message.Precommit{}
	case message.LightProposalCode:
		p = &message.LightProposal{}
	default:
		return errUnexpectedCode
	}
	if err := stream.Decode(p); err != nil {
		return fmt.Errorf("could not decode proof's typed message %w", err)
	}
	t.Msg = p
	return stream.ListEnd()
}

type Proof struct {
	Type      autonity.AccountabilityEventType // Accountability event types: Misbehaviour, Accusation, Innocence.
	Rule      autonity.Rule                    // Rule ID defined in AFD rule engine.
	Message   message.Msg                      // the consensus message which is accountable.
	Evidences []message.Msg                    // the proofs of the accountability event.
	Offender  common.Address                   // the offender who break the rule.
}

type encodedProof struct {
	Type      autonity.AccountabilityEventType
	Rule      autonity.Rule
	Offender  common.Address
	Message   typedMessage
	Evidences []typedMessage
}

func (p *Proof) EncodeRLP(w io.Writer) error {
	encoded := encodedProof{
		Type:     p.Type,
		Rule:     p.Rule,
		Offender: p.Offender,
	}
	encoded.Message = typedMessage{p.Message}
	for _, m := range p.Evidences {
		encoded.Evidences = append(encoded.Evidences, typedMessage{m})
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
	p.Offender = encoded.Offender
	p.Message = encoded.Message.Msg

	p.Evidences = make([]message.Msg, len(encoded.Evidences))
	for i := range encoded.Evidences {
		p.Evidences[i] = encoded.Evidences[i].Msg
	}
	return nil
}
