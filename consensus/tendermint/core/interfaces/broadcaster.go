package interfaces

import (
	"context"
	"github.com/autonity/autonity/consensus/tendermint/core/messageutils"
)

type Broadcaster interface {
	// Broadcast sends a message to all validators (include self)
	Broadcast(ctx context.Context, msg *messageutils.Message)
}
