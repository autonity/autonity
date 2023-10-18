package interfaces

import (
	"context"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
)

type Prevoter interface {
	SendPrevote(ctx context.Context, isNil bool)
	HandlePrevote(ctx context.Context, msg *message.Prevote) error
	LogPrevoteMessageEvent(message string, prevote *message.Vote, from, to string)
}
