package interfaces

import (
	"github.com/autonity/autonity/consensus/tendermint/core/message"
)

type Broadcaster interface {
	// Broadcast sends a message to all validators (include self)
	Broadcast(msg message.Msg)
}
