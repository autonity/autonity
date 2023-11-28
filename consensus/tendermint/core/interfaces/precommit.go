package interfaces

import (
	"context"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
)

type Precommiter interface {
	SendPrecommit(ctx context.Context, isNil bool)
	HandlePrecommit(ctx context.Context, msg *message.Precommit) error
	HandleCommit(ctx context.Context)
	LogPrecommitMessageEvent(message string, precommit *message.Precommit, from, to string)
}
