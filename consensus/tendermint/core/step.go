package core

type Step uint64

const (
	Propose Step = iota
	Prevote
	Precommit
	PrecommitDone
)

func (s Step) String() string {
	switch s {
	case Propose:
		return "propose"
	case Prevote:
		return "prevote"
	case Precommit:
		return "precommit"
	case PrecommitDone:
		return "precommitDone"
	default:
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
