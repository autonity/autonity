package interfaces

import (
	"context"

	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/core/types"
)

type Gossiper interface {
	// Gossip consensus message to the other committee members
	Gossip(ctx context.Context, committee types.Committee, payload []byte)
	SetBroadcaster(broadcaster consensus.Broadcaster)
}
