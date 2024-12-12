package interfaces

import (
	"github.com/autonity/autonity/consensus/tendermint/router"
	"github.com/autonity/autonity/consensus/tendermint/router/ping"
)

type Services struct {
	Broadcaster func(c Core) Broadcaster
	Prevoter    func(c Core) Prevoter
	Proposer    func(c Core) Proposer
	Precommiter func(c Core) Precommiter
	Gossiper func(b Backend) Gossiper
	Selector router.PeerSelector
	Pinger   ping.Pinger
}
