package autonity

// Todo(youssef): improve abigen to generate that automatically

type Rule uint8

const (
	PN Rule = iota
	PO
	PVN
	PVO
	PVO12
	C
	C1

	InvalidProposal // The value proposed by proposer cannot pass the blockchain's validation.
	InvalidProposer // A proposal sent from none proposer nodes of the committee.
	Equivocation    // Multiple distinguish votes(proposal, prevote, precommit) sent by validator.
)

type AccountabilityEventType uint8

const (
	Misbehaviour AccountabilityEventType = iota
	Accusation
	Innocence
)

// human understandable explanation of accountability rules
func (r Rule) Explanation() string {
	paper := "page 6 of \"The latest gossip on BFT consensus\" (https://arxiv.org/pdf/1807.04938.pdf) by Ethan Buchman, Jae Kwon, and Zarko Milosevic"
	var explanation string
	switch r {
	case PN:
		explanation = "Validator broadcasted a proposal with vr == -1, but had precommitted in a previous round.\n" +
			"Reference: line 15 and 42 in " + paper //nolint
	case PO:
		explanation = "Validator broadcasted a proposal with vr != 1, but he either:\n" +
			"1. had precommitted in vr for a different value.\n" +
			"2. had precommitted for a non-nil value between vr and currentRound.\n" +
			"3. there is a quorum of prevotes for a value different than the proposed one at vr.\n" +
			"4. If accusation, the local validator was not able to provide a quorum of prevotes at vr for the proposed value.\n" +
			"Reference: line 15, 36 and 40 in " + paper //nolint
	case PVN:
		explanation = "Validator has broadcasted a prevote for a proposal with vr == -1, but he either:\n" +
			"1. was locked on a different value.\n" +
			"2. If accusation, the local validator was not able to provide the proposal corresponding to his prevote.\n" +
			"NOTE: this could happen also for a prevote for a proposal with vr != -1.\n" +
			"Reference: line 22 in " + paper //nolint
	case PVO:
		explanation = "Validator has broadcasted a prevote for a proposal with vr != -1, but either:\n" +
			"1. there was a quorum of prevotes for another value in vr.\n" +
			"2. If accusation, the local validator was not able to provide the quorum of prevotes at vr which justify his prevote\n" +
			"Reference: line 28 in " + paper //nolint
	case PVO12:
		explanation = "Validator has broadcasted a prevote for a proposal with vr != -1, but he was locked on a different value at a round > vr.\n" +
			"Reference: line 28 in " + paper //nolint
	case C:
		explanation = "Validator has broadcasted a precommit, but there is a prevote quorum for another value in the same round.\n" +
			"Reference: line 36 in " + paper //nolint
	case C1:
		explanation = "Validator has broadcasted a precommit, and failed to provide supporting quorum of prevotes as innocence proof.\n" +
			"Reference: line 36 in " + paper //nolint
	case InvalidProposal:
		explanation = "Proposer broadcasted an invalid proposal."
	case InvalidProposer:
		explanation = "Proposer broadcasted a proposal when he was not the elected proposer.\n" +
			"Reference: line 14 in " + paper //nolint
	case Equivocation:
		explanation = "Validator broadcasted multiple messages during the same (height,round,phase) tuple.\n" +
			"Example: broadcast of 2 proposals, 2 prevotes or 2 precommits during the same round."
	default:
		explanation = "invalid rule" //nolint
	}
	return r.String() + " explanation: " + explanation
}

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
	default:
		return "invalid rule" //nolint
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
