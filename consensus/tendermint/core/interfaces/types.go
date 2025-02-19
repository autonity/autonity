package interfaces

import (
	"github.com/autonity/autonity/consensus/tendermint/latency"
	"github.com/autonity/autonity/consensus/tendermint/latency/ping"
)

type Services struct {
	Broadcaster func(c Core) Broadcaster
	Prevoter    func(c Core) Prevoter
	Proposer    func(c Core) Proposer
	Precommiter func(c Core) Precommiter
	Gossiper    func(b Backend) Gossiper
	Selector    latency.PeerSelector
	Pinger      ping.Pinger
}
