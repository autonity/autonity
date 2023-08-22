package interfaces

import (
	"context"

	"github.com/autonity/autonity/consensus/tendermint/core/message"
)

type Broadcaster interface {
	// Broadcast sends a message to all validators (include self)
	SignAndBroadcast(ctx context.Context, msg *message.Message)
}
