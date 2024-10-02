package interfaces

import (
	"context"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
)

type Committer interface {
	Commit(ctx context.Context, round int64, messages *message.RoundMessages)
}
