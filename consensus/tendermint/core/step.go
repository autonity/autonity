package core

type Step uint64

const (
	propose Step = iota
	prevote
	StepPrevoteDone
	StepPrecommitDone
)

func (s Step) String() string {
	if s == propose {
		return "propose"
	} else if s == prevote {
		return "prevote"
	} else if s == StepPrevoteDone {
		return "StepPrevoteDone"
	} else if s == StepPrecommitDone {
		return "StepPrecommitDone"
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
