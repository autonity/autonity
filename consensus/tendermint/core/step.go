package core

type Step uint64

const (
	Propose Step = iota
	Prevote
	Precommit
	PrecommitDone
)

func (s Step) String() string {
	if s == Propose {
		return "propose"
	} else if s == Prevote {
		return "prevote"
	} else if s == Precommit {
		return "precommit"
	} else if s == PrecommitDone {
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
