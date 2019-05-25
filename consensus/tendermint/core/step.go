package core

type Step uint64

const (
	StepAcceptProposal Step = iota
	StepProposeDone
	StepPrevoteDone
	StepPrecommitDone
)

func (s Step) String() string {
	if s == StepAcceptProposal {
		return "Accepting proposal"
	} else if s == StepProposeDone {
		return "Proposal"
	} else if s == StepPrevoteDone {
		return "Prevoted"
	} else if s == StepPrecommitDone {
		return "Precommitted"
	} else {
		return "Unknown"
	}
}

func (s Step) Cmp(y Step) int {
	if uint64(s) < uint64(y) {
		return -1
	}
	if uint64(s) > uint64(y) {
		return 1
	}
	return 0
}
