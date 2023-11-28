package interfaces

type Services struct {
	Broadcaster func(c Core) Broadcaster
	Prevoter    func(c Core) Prevoter
	Proposer    func(c Core) Proposer
	Precommiter func(c Core) Precommiter
	Gossiper    func(b Backend) Gossiper
}
