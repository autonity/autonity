package autonity

// Todo(youssef): improve abigen to generate that automatically

type Rule uint8

const (
	PN Rule = iota
	PO
	PVN
	PVO
	PVO12
	PVO3
	C
	C1

	InvalidProposal // The value proposed by proposer cannot pass the blockchain's validation.
	InvalidProposer // A proposal sent from none proposer nodes of the committee.
	Equivocation    // Multiple distinguish votes(proposal, prevote, precommit) sent by validator.

	InvalidRound    // message contains invalid round
	WrongValidRound // proposal contains a wrong valid round number.
	GarbageMessage  // message was signed by sender, but it cannot be decoded.
)

type AccountabilityEventType uint8

const (
	Misbehaviour AccountabilityEventType = iota
	Accusation
	Innocence
)

func (r Rule) String() string {
	switch r {
	case PN:
		return "PN"
	case PO:
		return "PO"
	case PVN:
		return "PVN"
	case PVO:
		return "PVO"
	case PVO12:
		return "PVO12"
	case PVO3:
		return "PVO3"
	case C:
		return "C"
	case C1:
		return "C1"
	case InvalidProposal:
		return "Invalid Proposal"
	case InvalidProposer:
		return "Invalid Proposer"
	case Equivocation:
		return "Equivocation"
	case InvalidRound:
		return "Invalid Round"
	case GarbageMessage:
		return "Garbage Message"
	case WrongValidRound:
		return "Wrong Valid Round"
	default:
		return "invalid rule"
	}
}

func (e AccountabilityEventType) String() string {
	switch e {
	case Misbehaviour:
		return "Fault Proof"
	case Accusation:
		return "Accusation"
	case Innocence:
		return "Innocence"
	}
	return "invalid"
}
