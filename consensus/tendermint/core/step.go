package core

type Step uint64

const (
	propose Step = iota
	prevote
	precommit
	precommitDone
)

func (s Step) String() string {
	if s == propose {
		return "propose"
	} else if s == prevote {
		return "prevote"
	} else if s == precommit {
		return "precommit"
	} else if s == precommitDone {
		return "precommitDone"
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
