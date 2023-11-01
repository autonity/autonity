package interfaces

import (
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/core/types"
)

type Gossiper interface {
	// Gossip consensus message to the other committee members
	Gossip(committee types.Committee, payload []byte)
	SetBroadcaster(broadcaster consensus.Broadcaster)
}
